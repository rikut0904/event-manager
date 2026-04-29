package handler

import (
	"net/http"
	"backend/internal/usecase"
	"github.com/labstack/echo/v4"
)

type EventHandler struct {
	eventUsecase usecase.EventUsecase
}

func NewEventHandler(u usecase.EventUsecase) *EventHandler {
	return &EventHandler{eventUsecase: u}
}

func (h *EventHandler) Create(c echo.Context) error {
	userID := c.Get("userID").(string)

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		StartTime   string `json:"start_time"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	event, err := h.eventUsecase.CreateEvent(c.Request().Context(), userID, req.Title, req.Description, req.StartTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, event)
}

func (h *EventHandler) GetMyEvents(c echo.Context) error {
	userID := c.Get("userID").(string)

	events, err := h.eventUsecase.GetMyEvents(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, events)
}
