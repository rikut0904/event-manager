package handler

import (
	"backend/internal/usecase"
	"encoding/json"
	"net/http"
)

type HealthHandler struct {
	useCase usecase.HealthUseCase
}

func NewHealthHandler(uc usecase.HealthUseCase) *HealthHandler {
	return &HealthHandler{useCase: uc}
}

func (h *HealthHandler) Handle(w http.ResponseWriter, r *http.Request) {
	status := h.useCase.Execute()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}
