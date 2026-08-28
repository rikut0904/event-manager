package usecase

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"backend/internal/domain"
)

type eventRepositoryMock struct {
	created   *domain.Event
	updated   *domain.Event
	found     *domain.Event
	findErr   error
	updateErr error
}

func (m *eventRepositoryMock) Create(_ context.Context, event *domain.Event) error {
	m.created = event
	return nil
}

func (m *eventRepositoryMock) Update(_ context.Context, event *domain.Event) error {
	m.updated = event
	return m.updateErr
}

func (m *eventRepositoryMock) Delete(context.Context, string) error { return nil }

func (m *eventRepositoryMock) FindByCreatorID(context.Context, string) ([]domain.Event, error) {
	return nil, nil
}

func (m *eventRepositoryMock) FindByID(context.Context, string) (*domain.Event, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	return m.found, nil
}

func (m *eventRepositoryMock) FindByDisplayID(context.Context, string) (*domain.Event, error) {
	return nil, nil
}

func (m *eventRepositoryMock) FindAllPublished(context.Context) ([]domain.Event, error) {
	return nil, nil
}

func (m *eventRepositoryMock) UpdateExpiredStatus(context.Context) error { return nil }

func TestEventUsecase_CreateEventValidation(t *testing.T) {
	now := time.Now().Add(2 * time.Hour).Truncate(time.Second)
	end := now.Add(time.Hour)

	tests := []struct {
		name      string
		startTime string
		endTime   string
		status    string
		sourceURL string
		wantError string
	}{
		{
			name:      "開始日時が過去",
			startTime: time.Now().Add(-time.Hour).Format(time.RFC3339),
			wantError: "開始時間は現在時刻より後",
		},
		{
			name:      "終了日時が開始日時と同じ",
			startTime: now.Format(time.RFC3339),
			endTime:   now.Format(time.RFC3339),
			wantError: "終了時間は開始時間より後",
		},
		{
			name:      "ステータスが不正",
			startTime: now.Format(time.RFC3339),
			endTime:   end.Format(time.RFC3339),
			status:    "archived",
			wantError: "ステータスが不正",
		},
		{
			name:      "外部リンクがhttp(s)ではない",
			startTime: now.Format(time.RFC3339),
			endTime:   end.Format(time.RFC3339),
			sourceURL: "javascript:alert(1)",
			wantError: "httpまたはhttps",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &eventRepositoryMock{}
			uc := NewEventUsecase(repo)

			_, err := uc.CreateEvent(context.Background(), "user-1", "テストイベント", "説明", tt.startTime, tt.endTime, "会場", tt.status, false, 10, tt.sourceURL, "")
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("error = %v, want error containing %q", err, tt.wantError)
			}
			if repo.created != nil {
				t.Fatal("invalid event was passed to repository")
			}
		})
	}
}

func TestEventUsecase_CreateEvent_PreservesPublishedStatus(t *testing.T) {
	repo := &eventRepositoryMock{}
	uc := NewEventUsecase(repo)
	start := time.Now().Add(2 * time.Hour).Truncate(time.Second)

	event, err := uc.CreateEvent(context.Background(), "user-1", "公開イベント", "説明", start.Format(time.RFC3339), "", "会場", "published", false, 20, "https://example.com/apply", "")
	if err != nil {
		t.Fatalf("CreateEvent() error = %v", err)
	}
	if event.Status != "published" || repo.created == nil || repo.created.Status != "published" {
		t.Fatalf("event status = %q, want published", event.Status)
	}
}

func TestEventUsecase_UpdateEvent_RejectsInvalidAndFinishedStatus(t *testing.T) {
	start := time.Now().Add(2 * time.Hour).Truncate(time.Second)

	tests := []struct {
		name      string
		current   string
		newStatus string
		wantError string
	}{
		{name: "任意のステータス", current: "draft", newStatus: "archived", wantError: "ステータスが不正"},
		{name: "終了済みからの変更", current: "finished", newStatus: "draft", wantError: "終了済みイベントのステータスは変更できません"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &eventRepositoryMock{found: &domain.Event{ID: "event-1", CreatorID: "user-1", Status: tt.current}}
			uc := NewEventUsecase(repo)

			_, err := uc.UpdateEvent(context.Background(), "user-1", "event-1", "更新イベント", "説明", start.Format(time.RFC3339), "", "会場", tt.newStatus, nil, nil, nil, nil)
			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("error = %v, want error containing %q", err, tt.wantError)
			}
			if repo.updated != nil {
				t.Fatal("invalid status update was passed to repository")
			}
		})
	}
}

func TestEventUsecase_UpdateEvent_AllowsValidTransition(t *testing.T) {
	start := time.Now().Add(2 * time.Hour).Truncate(time.Second)
	online := true
	capacity := 50
	sourceURL := "https://example.com/event"
	thumbnailURL := "https://example.com/image.png"
	repo := &eventRepositoryMock{found: &domain.Event{ID: "event-1", CreatorID: "user-1", Status: "draft"}}
	uc := NewEventUsecase(repo)

	event, err := uc.UpdateEvent(context.Background(), "user-1", "event-1", "更新イベント", "説明", start.Format(time.RFC3339), "", "会場", "published", &online, &capacity, &sourceURL, &thumbnailURL)
	if err != nil {
		t.Fatalf("UpdateEvent() error = %v", err)
	}
	if repo.updated == nil {
		t.Fatal("valid update was not passed to repository")
	}
	if event.Status != "published" || !event.IsOnline || event.Capacity != capacity || event.SourceURL != sourceURL || event.ThumbnailURL != thumbnailURL {
		t.Fatalf("updated event was not preserved: %+v", event)
	}
}

func TestEventUsecase_UpdateEvent_PropagatesRepositoryError(t *testing.T) {
	start := time.Now().Add(2 * time.Hour).Truncate(time.Second)
	repo := &eventRepositoryMock{
		found:     &domain.Event{ID: "event-1", CreatorID: "user-1", Status: "draft"},
		updateErr: errors.New("保存失敗"),
	}
	uc := NewEventUsecase(repo)

	_, err := uc.UpdateEvent(context.Background(), "user-1", "event-1", "更新イベント", "説明", start.Format(time.RFC3339), "", "会場", "", nil, nil, nil, nil)
	if !errors.Is(err, repo.updateErr) {
		t.Fatalf("error = %v, want repository error", err)
	}
}
