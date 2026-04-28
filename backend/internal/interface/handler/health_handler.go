package handler

import (
	"net/http"
	"backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type HealthHandler struct {
	useCase usecase.HealthUsecase
}

func NewHealthHandler(uc usecase.HealthUsecase) *HealthHandler {
	return &HealthHandler{useCase: uc}
}

func (h *HealthHandler) HealthCheck(c echo.Context) error {
	status := h.useCase.Execute()
	return c.JSON(http.StatusOK, status)
}
