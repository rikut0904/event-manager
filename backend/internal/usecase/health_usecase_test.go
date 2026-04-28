package usecase

import (
	"testing"
	"time"
)

func TestHealthInteractor_Execute(t *testing.T) {
	uc := NewHealthUseCase()
	result := uc.Execute()

	if result.Status != "OK" {
		t.Errorf("Expected status 'OK', got '%s'", result.Status)
	}

	if result.Timestamp == "" {
		t.Error("Expected timestamp to be present, got empty string")
	}

	_, err := time.Parse("2006-01-02T15:04:05.000Z", result.Timestamp)
	if err != nil {
		t.Errorf("Timestamp format is invalid: %v", err)
	}
}
