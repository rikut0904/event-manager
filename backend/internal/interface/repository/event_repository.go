package repository

import (
	"context"
	"time"

	"backend/internal/domain"
	"gorm.io/gorm"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) error
	Update(ctx context.Context, event *domain.Event) error
	Delete(ctx context.Context, id string) error
	FindByCreatorID(ctx context.Context, creatorID string) ([]domain.Event, error)
	FindByID(ctx context.Context, id string) (*domain.Event, error)
	FindByDisplayID(ctx context.Context, displayID string) (*domain.Event, error)
	FindAllPublished(ctx context.Context) ([]domain.Event, error)
	UpdateExpiredStatus(ctx context.Context) error
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

func (r *eventRepository) Update(ctx context.Context, event *domain.Event) error {
	return r.db.WithContext(ctx).Model(event).Updates(map[string]interface{}{
		"title":         event.Title,
		"description":   event.Description,
		"start_time":    event.StartTime,
		"end_time":      event.EndTime,
		"location":      event.Location,
		"status":        event.Status,
		"is_online":     event.IsOnline,
		"capacity":      event.Capacity,
		"source_url":    event.SourceURL,
		"thumbnail_url": event.ThumbnailURL,
		"updated_at":    time.Now(),
	}).Error
}

func (r *eventRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Event{}, "id = ?", id).Error
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

func (r *eventRepository) FindByDisplayID(ctx context.Context, displayID string) (*domain.Event, error) {
	var event domain.Event
	err := r.db.WithContext(ctx).Where("display_id = ? AND status = ?", displayID, "published").First(&event).Error
	return &event, err
}

func (r *eventRepository) FindAllPublished(ctx context.Context) ([]domain.Event, error) {
	var events []domain.Event
	err := r.db.WithContext(ctx).Where("status = ?", "published").Order("start_time asc").Find(&events).Error
	return events, err
}

func (r *eventRepository) UpdateExpiredStatus(ctx context.Context) error {
	// 終了時刻が設定され、終了時刻を過ぎていて、かつステータスが published のものを finished に更新
	return r.db.WithContext(ctx).Model(&domain.Event{}).
		Where("status = ? AND end_time IS NOT NULL AND end_time > ? AND end_time < ?", "published", time.Time{}, time.Now()).
		Update("status", "finished").Error
}
