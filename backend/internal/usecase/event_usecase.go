package usecase

import (
	"context"
	"strings"
	"time"

	"backend/internal/domain"
	"backend/internal/interface/repository"
	"github.com/google/uuid"
)

type EventUsecase interface {
	CreateEvent(ctx context.Context, creatorID string, title, description string, startTime string) (*domain.Event, error)
	GetMyEvents(ctx context.Context, creatorID string) ([]domain.Event, error)
}

type eventUsecase struct {
	eventRepo repository.EventRepository
}

func NewEventUsecase(eventRepo repository.EventRepository) EventUsecase {
	return &eventUsecase{eventRepo: eventRepo}
}

func (u *eventUsecase) CreateEvent(ctx context.Context, creatorID string, title, description string, startTime string) (*domain.Event, error) {
	// DisplayID の自動生成 (UUIDをベースにした短い文字列)
	shortID := strings.ReplaceAll(uuid.New().String(), "-", "")[:12]
	
	// 時間のパース (RFC3339 形式を想定)
	parsedTime, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		parsedTime = time.Now() // フォールバック
	}
	
	event := &domain.Event{
		CreatorID:   creatorID,
		Title:       title,
		Description: description,
		StartTime:   parsedTime,
		DisplayID:   shortID,
		Status:      "draft", // 初期状態は下書き
	}
	
	err = u.eventRepo.Create(ctx, event)
	return event, err
}

func (u *eventUsecase) GetMyEvents(ctx context.Context, creatorID string) ([]domain.Event, error) {
	return u.eventRepo.FindByCreatorID(ctx, creatorID)
}
