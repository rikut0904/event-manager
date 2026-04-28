package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"backend/internal/domain"
	"backend/internal/infrastructure/firebase"
	"firebase.google.com/go/v4/auth"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (*domain.AuthResponse, error)
	SignUp(ctx context.Context, email, password string) (*domain.User, error)
	LinkConnpass(ctx context.Context, userID, connpassID string) error
}

type authUsecase struct {
	fbClient *firebase.Client
}

func NewAuthUsecase(fbClient *firebase.Client) AuthUsecase {
	return &authUsecase{fbClient: fbClient}
}

// Login は Firebase Auth REST API を使用してサインインし、IDトークンを取得します
func (u *authUsecase) Login(ctx context.Context, email, password string) (*domain.AuthResponse, error) {
	apiKey := os.Getenv("FIREBASE_API_KEY")
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", apiKey)

	payload := map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to login: status %d", resp.StatusCode)
	}

	var result struct {
		IDToken string `json:"idToken"`
		LocalID string `json:"localId"`
		Email   string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// ユーザー情報の取得（カスタムクレーム等を含む可能性を考慮）
	fbUser, err := u.fbClient.Auth.GetUser(ctx, result.LocalID)
	if err != nil {
		return nil, err
	}

	connpassID := ""
	if cid, ok := fbUser.CustomClaims["connpass_id"].(string); ok {
		connpassID = cid
	}

	return &domain.AuthResponse{
		Token: result.IDToken,
		User: &domain.User{
			ID:         result.LocalID,
			Email:      result.Email,
			ConnpassID: connpassID,
		},
	}, nil
}

func (u *authUsecase) SignUp(ctx context.Context, email, password string) (*domain.User, error) {
	params := (&auth.UserToCreate{}).
		Email(email).
		Password(password)
	
	fbUser, err := u.fbClient.Auth.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:    fbUser.UID,
		Email: fbUser.Email,
	}, nil
}

func (u *authUsecase) LinkConnpass(ctx context.Context, userID, connpassID string) error {
	// Connpass API でユーザーが存在するか確認
	// 実際には /event/ API などで nickname を指定して確認
	url := fmt.Sprintf("https://connpass.com/api/v1/event/?nickname=%s", connpassID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("connpass api error: status %d", resp.StatusCode)
	}

	// Firebaseのカスタムクレームに connpass_id を保存
	claims := map[string]interface{}{"connpass_id": connpassID}
	err = u.fbClient.Auth.SetCustomUserClaims(ctx, userID, claims)
	if err != nil {
		return err
	}

	return nil
}
