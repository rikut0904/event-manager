// Command migrate-users transfers the application ownership IDs from Firebase
// UIDs to Common ID user IDs. It is intentionally dry-run by default.
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"backend/internal/domain"
	"backend/internal/infrastructure/firebase"
)

type migrationUser struct {
	SourceUserID  string `json:"source_user_id"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

type migrationResult struct {
	SourceUserID string `json:"source_user_id"`
	Email        string `json:"email"`
	CommonUserID string `json:"common_user_id"`
	Status       string `json:"status"`
	Error        string `json:"error"`
}

type migrationResponse struct {
	Results []migrationResult `json:"results"`
}

type userMapping struct {
	SourceUserID string `json:"source_user_id"`
	CommonUserID string `json:"common_user_id"`
	Status       string `json:"status"`
}

func main() {
	loadDotEnv()
	endpoint := flag.String("endpoint", envOrDefault("COMMON_ID_API_ORIGIN", "http://localhost:18080"), "Common ID API URL")
	apiKey := flag.String("api-key", "", "Common ID application API key (defaults to COMMON_ID_API_KEY)")
	clientID := flag.String("client-id", "", "registered Common ID client ID (defaults to COMMON_ID_CLIENT_ID)")
	inputPath := flag.String("input", "", "optional JSONL input; when omitted, users are read from Firebase Auth")
	outputPath := flag.String("output", "", "optional JSON mapping output path")
	dryRun := flag.Bool("dry-run", true, "validate and show the migration without changing the app DB")
	updateDB := flag.Bool("update-db", false, "update users.id and events.creator_id after a successful migration")
	allowUnverified := flag.Bool("allow-unverified", false, "allow migration of unverified source emails (also sends allow_unverified=true)")
	flag.Parse()

	apiKeyValue := firstNonEmpty(*apiKey, os.Getenv("COMMON_ID_API_KEY"))
	clientIDValue := firstNonEmpty(*clientID, os.Getenv("COMMON_ID_CLIENT_ID"))
	if clientIDValue == "" || apiKeyValue == "" {
		log.Fatal("-client-id and -api-key (or COMMON_ID_CLIENT_ID and COMMON_ID_API_KEY) are required")
	}
	if *updateDB && *dryRun {
		log.Fatal("-update-db=true requires -dry-run=false")
	}

	ctx := context.Background()
	db, err := openDB()
	if err != nil {
		log.Fatal(err)
	}
	users, err := loadUsers(ctx, db, *inputPath)
	if err != nil {
		log.Fatal(err)
	}
	if len(users) == 0 {
		log.Fatal("移管対象のユーザーがありません")
	}

	results, err := migrate(ctx, *endpoint, apiKeyValue, clientIDValue, *dryRun, *allowUnverified, users)
	if err != nil {
		log.Fatal(err)
	}
	if *dryRun {
		for _, result := range results {
			log.Printf("source=%s status=%s email=%s error=%s", result.SourceUserID, result.Status, result.Email, result.Error)
		}
		log.Printf("dry-run完了: %d users（DBは変更していません）", len(results))
		return
	}
	mappings, err := buildMappings(results)
	if err != nil {
		log.Fatal(err)
	}
	if *outputPath != "" {
		if err := writeMappings(*outputPath, mappings); err != nil {
			log.Fatal(err)
		}
	}

	if *updateDB {
		if err := applyMappings(ctx, db, mappings); err != nil {
			log.Fatal(err)
		}
		log.Printf("DB移管完了: %d users", len(mappings))
	} else {
		log.Printf("dry-run完了: %d users（DBは変更していません）", len(mappings))
	}
}

func loadDotEnv() {
	for _, path := range []string{".env", "../.env"} {
		if err := godotenv.Load(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Printf("WARNING: %s を読み込めませんでした: %v", path, err)
		}
	}
}

func openDB() (*gorm.DB, error) {
	dsn := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dsn == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func loadUsers(ctx context.Context, db *gorm.DB, inputPath string) ([]migrationUser, error) {
	if inputPath != "" {
		return readJSONL(inputPath)
	}
	var stored []domain.User
	if err := db.WithContext(ctx).Order("id").Find(&stored).Error; err != nil {
		return nil, fmt.Errorf("usersの読み込みに失敗しました: %w", err)
	}
	fb := firebase.NewClient(ctx)
	users := make([]migrationUser, 0, len(stored))
	for _, user := range stored {
		record, err := fb.Auth.GetUser(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("Firebase user %s の読み込みに失敗しました: %w", user.ID, err)
		}
		users = append(users, migrationUser{SourceUserID: record.UID, Email: record.Email, EmailVerified: record.EmailVerified})
	}
	return users, nil
}

func readJSONL(path string) ([]migrationUser, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	users := make([]migrationUser, 0)
	for line := 1; scanner.Scan(); line++ {
		if strings.TrimSpace(scanner.Text()) == "" {
			continue
		}
		var user migrationUser
		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, fmt.Errorf("%s:%d: invalid JSON: %w", path, line, err)
		}
		users = append(users, user)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func migrate(ctx context.Context, endpoint, apiKey, clientID string, dryRun, allowUnverified bool, users []migrationUser) ([]migrationResult, error) {
	all := make([]migrationResult, 0, len(users))
	httpClient := &http.Client{Timeout: 30 * time.Second}
	for start := 0; start < len(users); start += 1000 {
		end := start + 1000
		if end > len(users) {
			end = len(users)
		}
		body, err := json.Marshal(map[string]any{
			"client_id":        clientID,
			"dry_run":          dryRun,
			"allow_unverified": allowUnverified,
			"users":            users[start:end],
		})
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(endpoint, "/")+"/v1/migrations/users", bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-API-Key", apiKey)
		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("Common ID移管APIへの接続に失敗しました: %w", err)
		}
		responseBody, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			return nil, readErr
		}
		if resp.StatusCode >= http.StatusMultipleChoices {
			return nil, fmt.Errorf("Common ID移管APIが%sを返しました: %s", resp.Status, strings.TrimSpace(string(responseBody)))
		}
		var result migrationResponse
		if err := json.Unmarshal(responseBody, &result); err != nil {
			return nil, fmt.Errorf("Common ID移管APIの応答を解釈できません: %w", err)
		}
		all = append(all, result.Results...)
	}
	return all, nil
}

func buildMappings(results []migrationResult) ([]userMapping, error) {
	mappings := make([]userMapping, 0, len(results))
	seenSource := make(map[string]struct{}, len(results))
	seenTarget := make(map[string]string, len(results))
	for _, result := range results {
		if result.SourceUserID == "" || result.CommonUserID == "" {
			return nil, fmt.Errorf("移管結果にIDがありません: source=%q status=%s error=%s", result.SourceUserID, result.Status, result.Error)
		}
		if result.Status == "failed" || result.Status == "email_pending" {
			return nil, fmt.Errorf("ユーザー %s の移管に失敗しました: %s", result.SourceUserID, result.Error)
		}
		if _, ok := seenSource[result.SourceUserID]; ok {
			return nil, fmt.Errorf("source_user_idが重複しています: %s", result.SourceUserID)
		}
		if previous, ok := seenTarget[result.CommonUserID]; ok && previous != result.SourceUserID {
			return nil, fmt.Errorf("common_user_idが複数ユーザーに割り当てられています: %s", result.CommonUserID)
		}
		seenSource[result.SourceUserID] = struct{}{}
		seenTarget[result.CommonUserID] = result.SourceUserID
		mappings = append(mappings, userMapping{SourceUserID: result.SourceUserID, CommonUserID: result.CommonUserID, Status: result.Status})
	}
	return mappings, nil
}

func writeMappings(path string, mappings []userMapping) error {
	data, err := json.MarshalIndent(mappings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0600)
}

func applyMappings(ctx context.Context, db *gorm.DB, mappings []userMapping) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, mapping := range mappings {
			var count int64
			if err := tx.Model(&domain.User{}).Where("id = ?", mapping.CommonUserID).Count(&count).Error; err != nil {
				return err
			}
			if count > 0 && mapping.SourceUserID != mapping.CommonUserID {
				return fmt.Errorf("移管先IDが既存ユーザーと衝突しています: %s", mapping.CommonUserID)
			}
		}
		for _, mapping := range mappings {
			temporaryID := "migration-" + uuid.NewString()
			result := tx.Model(&domain.User{}).Where("id = ?", mapping.SourceUserID).Update("id", temporaryID)
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected != 1 {
				return fmt.Errorf("移管元ユーザーが見つかりません: %s", mapping.SourceUserID)
			}
			if err := tx.Model(&domain.Event{}).Where("creator_id = ?", mapping.SourceUserID).Update("creator_id", temporaryID).Error; err != nil {
				return err
			}
			if err := tx.Model(&domain.User{}).Where("id = ?", temporaryID).Update("id", mapping.CommonUserID).Error; err != nil {
				return err
			}
			if err := tx.Model(&domain.Event{}).Where("creator_id = ?", temporaryID).Update("creator_id", mapping.CommonUserID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func envOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func firstNonEmpty(value, fallback string) string {
	if strings.TrimSpace(value) != "" {
		return value
	}
	return strings.TrimSpace(fallback)
}
