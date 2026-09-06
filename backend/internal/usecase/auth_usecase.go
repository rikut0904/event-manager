package usecase

import (
	"context"
	"errors"
	"strings"

	"backend/internal/domain"
	"backend/internal/interface/repository"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	LinkConnpass(ctx context.Context, userID, connpassID string) error
	SyncCommonUser(ctx context.Context, commonUserID, email string) (*domain.User, error)
	GetUser(ctx context.Context, userID string) (*domain.User, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
}

func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
	return &authUsecase{userRepo: userRepo}
}

func (u *authUsecase) LinkConnpass(ctx context.Context, userID, connpassID string) error {
	connpassID = strings.TrimSpace(connpassID)
	if connpassID == "" {
		return errors.New("connpass IDは必須です")
	}

	return u.userRepo.UpdateConnpassID(ctx, userID, connpassID)
}

func (u *authUsecase) SyncCommonUser(ctx context.Context, commonUserID, email string) (*domain.User, error) {
	if strings.TrimSpace(commonUserID) == "" || strings.TrimSpace(email) == "" {
		return nil, errors.New("Common IDのユーザー情報が不正です")
	}
	user, err := u.userRepo.FindByID(ctx, commonUserID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		user = &domain.User{ID: commonUserID, Email: email}
		if err := u.userRepo.Save(ctx, user); err != nil {
			return nil, err
		}
		return user, nil
	}
	if user.Email != email {
		user.Email = email
		if err := u.userRepo.Save(ctx, user); err != nil {
			return nil, err
		}
	}
	return user, nil
}

func (u *authUsecase) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	return u.userRepo.FindByID(ctx, userID)
}
