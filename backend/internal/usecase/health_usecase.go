package usecase

import (
	"backend/internal/domain"
	"time"
)

type HealthUsecase interface {
	Execute() domain.HealthStatus
}

type healthInteractor struct{}

func NewHealthUsecase() HealthUsecase {
	return &healthInteractor{}
}

func (i *healthInteractor) Execute() domain.HealthStatus {
	return domain.HealthStatus{
		Status:    "OK",
		Timestamp: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
	}
}
