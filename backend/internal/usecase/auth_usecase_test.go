package usecase

import (
	"context"
	"errors"
	"testing"

	"backend/internal/domain"
	"gorm.io/gorm"
)

type authRepositoryMock struct {
	user    *domain.User
	findErr error
	saved   *domain.User
}

func (m *authRepositoryMock) Save(_ context.Context, user *domain.User) error {
	m.saved = user
	return nil
}

func (m *authRepositoryMock) FindByID(context.Context, string) (*domain.User, error) {
	return m.user, m.findErr
}

func (m *authRepositoryMock) UpdateConnpassID(context.Context, string, string) error {
	return nil
}

func TestAuthUsecase_SyncCommonUser_CreatesWhenNotFound(t *testing.T) {
	repo := &authRepositoryMock{findErr: gorm.ErrRecordNotFound}
	uc := NewAuthUsecase(repo)

	user, err := uc.SyncCommonUser(context.Background(), "common-user-1", "user@example.com")
	if err != nil {
		t.Fatalf("SyncCommonUser() error = %v", err)
	}
	if user == nil || repo.saved == nil || repo.saved.ID != "common-user-1" {
		t.Fatalf("user was not created: user=%+v saved=%+v", user, repo.saved)
	}
}

func TestAuthUsecase_SyncCommonUser_ReturnsUnexpectedFindError(t *testing.T) {
	findErr := errors.New("database connection failed")
	repo := &authRepositoryMock{findErr: findErr}
	uc := NewAuthUsecase(repo)

	user, err := uc.SyncCommonUser(context.Background(), "common-user-1", "user@example.com")
	if !errors.Is(err, findErr) {
		t.Fatalf("error = %v, want %v", err, findErr)
	}
	if user != nil || repo.saved != nil {
		t.Fatalf("unexpected upsert after repository error: user=%+v saved=%+v", user, repo.saved)
	}
}
