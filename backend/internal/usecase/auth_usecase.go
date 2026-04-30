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
	"backend/internal/interface/repository"
	firebaseauth "firebase.google.com/go/v4/auth"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (*domain.AuthResponse, error)
	SignUp(ctx context.Context, email, password string) (*domain.User, error)
}

type authUsecase struct {
	fbClient *firebase.Client
	userRepo repository.UserRepository
}

func NewAuthUsecase(fbClient *firebase.Client, userRepo repository.UserRepository) AuthUsecase {
	return &authUsecase{fbClient: fbClient, userRepo: userRepo}
}

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

	// DBからユーザー情報を取得
	user, err := u.userRepo.FindByID(ctx, result.LocalID)
	if err != nil {
		// まだDBにユーザーがいない場合は作成（既存Firebaseユーザー対応）
		user = &domain.User{ID: result.LocalID, Email: result.Email}
		u.userRepo.Save(ctx, user)
	}

	return &domain.AuthResponse{
		Token: result.IDToken,
		User:  user,
	}, nil
}

func (u *authUsecase) SignUp(ctx context.Context, email, password string) (*domain.User, error) {
	params := (&firebaseauth.UserToCreate{}).
		Email(email).
		Password(password)

	firebaseUser, err := u.fbClient.Auth.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:    firebaseUser.UID,
		Email: firebaseUser.Email,
	}
	if err := u.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
