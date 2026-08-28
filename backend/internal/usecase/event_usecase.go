package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"backend/internal/domain"
	"backend/internal/interface/repository"
	"github.com/google/uuid"
)

type EventUsecase interface {
	CreateEvent(ctx context.Context, creatorID string, title, description string, startTime string, endTime string, location string, status string, isOnline bool, capacity int, sourceURL string, thumbnailURL string) (*domain.Event, error)
	UpdateEvent(ctx context.Context, creatorID string, eventID string, title, description string, startTime string, endTime string, location string, status string, isOnline *bool, capacity *int, sourceURL *string, thumbnailURL *string) (*domain.Event, error)
	DeleteEvent(ctx context.Context, creatorID string, eventID string) error
	GetMyEvents(ctx context.Context, creatorID string) ([]domain.Event, error)
	GetEventByID(ctx context.Context, requesterID string, eventID string) (*domain.Event, error)
	GetPublicEventByDisplayID(ctx context.Context, displayID string) (*domain.Event, error)
	GetAllPublishedEvents(ctx context.Context) ([]domain.Event, error)
	UpdateExpiredEvents(ctx context.Context) error
}

type eventUsecase struct {
	eventRepo repository.EventRepository
}

func NewEventUsecase(eventRepo repository.EventRepository) EventUsecase {
	return &eventUsecase{eventRepo: eventRepo}
}

func (u *eventUsecase) validateEvent(title string, startTime, endTime time.Time) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("イベント名は必須です")
	}
	if !endTime.IsZero() && !endTime.After(startTime) {
		return errors.New("終了時間は開始時間より後である必要があります")
	}
	return nil
}

func (u *eventUsecase) CreateEvent(ctx context.Context, creatorID string, title, description string, startTime string, endTime string, location string, status string, isOnline bool, capacity int, sourceURL string, thumbnailURL string) (*domain.Event, error) {
	shortID := strings.ReplaceAll(uuid.New().String(), "-", "")[:12]
	parsedStart, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, errors.New("開始日時の形式が正しくありません")
	}
	var parsedEnd time.Time
	if endTime != "" && endTime != "null" {
		parsedEnd, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			return nil, errors.New("終了日時の形式が正しくありません")
		}
	}

	if parsedStart.Before(time.Now()) {
		return nil, errors.New("開始時間は現在時刻より後である必要があります")
	}

	if err := u.validateEvent(title, parsedStart, parsedEnd); err != nil {
		return nil, err
	}
	if status == "" {
		status = "draft"
	}
	if status != "draft" && status != "published" {
		return nil, errors.New("ステータスが不正です")
	}
	event := &domain.Event{
		CreatorID: creatorID, Title: title, Description: description,
		StartTime: parsedStart, EndTime: parsedEnd, Location: location,
		DisplayID: shortID, Status: status,
		IsOnline: isOnline, Capacity: capacity, SourceURL: sourceURL, ThumbnailURL: thumbnailURL,
	}
	err = u.eventRepo.Create(ctx, event)
	return event, err
}

func (u *eventUsecase) UpdateEvent(ctx context.Context, creatorID string, eventID string, title, description string, startTime string, endTime string, location string, status string, isOnline *bool, capacity *int, sourceURL *string, thumbnailURL *string) (*domain.Event, error) {
	event, err := u.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, errors.New("event not found")
	}
	if event.CreatorID != creatorID {
		return nil, errors.New("unauthorized")
	}

	parsedStart, parseErr := time.Parse(time.RFC3339, startTime)
	if parseErr != nil {
		return nil, errors.New("開始日時の形式が正しくありません")
	}
	var parsedEnd time.Time
	if endTime != "" && endTime != "null" {
		parsedEnd, parseErr = time.Parse(time.RFC3339, endTime)
		if parseErr != nil {
			return nil, errors.New("終了日時の形式が正しくありません")
		}
	}

	if err := u.validateEvent(title, parsedStart, parsedEnd); err != nil {
		return nil, err
	}

	event.Title = title
	event.Description = description
	event.StartTime = parsedStart
	event.EndTime = parsedEnd
	event.Location = location

	if isOnline != nil {
		event.IsOnline = *isOnline
	}
	if capacity != nil {
		event.Capacity = *capacity
	}
	if sourceURL != nil {
		event.SourceURL = *sourceURL
	}
	if thumbnailURL != nil {
		event.ThumbnailURL = *thumbnailURL
	}

	if status != "" {
		if event.Status == "published" && status == "draft" {
			return nil, errors.New("once published, event cannot be moved back to draft")
		}
		event.Status = status
	}

	err = u.eventRepo.Update(ctx, event)
	return event, err
}

func (u *eventUsecase) DeleteEvent(ctx context.Context, creatorID string, eventID string) error {
	event, err := u.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return errors.New("event not found")
	}
	if event.CreatorID != creatorID {
		return errors.New("unauthorized")
	}
	return u.eventRepo.Delete(ctx, eventID)
}

func (u *eventUsecase) GetMyEvents(ctx context.Context, creatorID string) ([]domain.Event, error) {
	return u.eventRepo.FindByCreatorID(ctx, creatorID)
}

func (u *eventUsecase) GetEventByID(ctx context.Context, requesterID string, eventID string) (*domain.Event, error) {
	event, err := u.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if event.Status == "published" || event.Status == "finished" {
		return event, nil
	}
	if event.CreatorID == requesterID && requesterID != "" {
		return event, nil
	}
	return nil, errors.New("event not found")
}

func (u *eventUsecase) GetPublicEventByDisplayID(ctx context.Context, displayID string) (*domain.Event, error) {
	if strings.TrimSpace(displayID) == "" {
		return nil, errors.New("event not found")
	}

	return u.eventRepo.FindByDisplayID(ctx, displayID)
}

func (u *eventUsecase) GetAllPublishedEvents(ctx context.Context) ([]domain.Event, error) {
	return u.eventRepo.FindAllPublished(ctx)
}

func (u *eventUsecase) UpdateExpiredEvents(ctx context.Context) error {
	return u.eventRepo.UpdateExpiredStatus(ctx)
}
