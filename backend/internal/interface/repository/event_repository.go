package repository

import (
	"context"
	"backend/internal/domain"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	FindByCreatorID(ctx context.Context, creatorID string) ([]domain.Event, error)
	FindByID(ctx context.Context, id string) (*domain.Event, error)
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) Create(ctx context.Context, event *domain.Event) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *eventRepository) FindByCreatorID(ctx context.Context, creatorID string) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.WithContext(ctx).Where("creator_id = ?", creatorID).Order("start_time desc").Find(&events).Error
	return events, err
}

func (r *eventRepository) FindByID(ctx context.Context, id string) (*domain.Event, error) {
	var event domain.Event
	err := r.db.WithContext(ctx).First(&event, "id = ?", id).Error
	return &event, err
}
