package repository

import (
	"context"
	"backend/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	UpdateConnpassID(ctx context.Context, id string, connpassID string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Save(ctx context.Context, user *domain.User) error {
	// Upsert (存在すれば更新、なければ挿入)
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"email", "name", "username", "avatar_url", "bio", "updated_at"}),
	}).Create(user).Error
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateConnpassID(ctx context.Context, id string, connpassID string) error {
	return r.db.WithContext(ctx).Model(&domain.User{}).Where("id = ?", id).Update("connpass_id", connpassID).Error
}
