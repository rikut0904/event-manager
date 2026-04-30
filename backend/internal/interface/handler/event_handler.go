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
		EndTime     string `json:"end_time"`
		Location    string `json:"location"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	event, err := h.eventUsecase.CreateEvent(c.Request().Context(), userID, req.Title, req.Description, req.StartTime, req.EndTime, req.Location)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, event)
}

func (h *EventHandler) Update(c echo.Context) error {
	userID := c.Get("userID").(string)
	eventID := c.Param("id")

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
		Location    string `json:"location"`
		Status      string `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	event, err := h.eventUsecase.UpdateEvent(c.Request().Context(), userID, eventID, req.Title, req.Description, req.StartTime, req.EndTime, req.Location, req.Status)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, event)
}

func (h *EventHandler) Delete(c echo.Context) error {
	userID := c.Get("userID").(string)
	eventID := c.Param("id")

	err := h.eventUsecase.DeleteEvent(c.Request().Context(), userID, eventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *EventHandler) GetMyEvents(c echo.Context) error {
	userID := c.Get("userID").(string)

	events, err := h.eventUsecase.GetMyEvents(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, events)
}

func (h *EventHandler) GetPublished(c echo.Context) error {
	events, err := h.eventUsecase.GetAllPublishedEvents(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, events)
}

func (h *EventHandler) GetByID(c echo.Context) error {
	// OptionalAuthMiddleware により、ログインしていれば userID が入る
	userID, _ := c.Get("userID").(string)
	eventID := c.Param("id")

	event, err := h.eventUsecase.GetEventByID(c.Request().Context(), userID, eventID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Event not found")
	}

	return c.JSON(http.StatusOK, event)
}

func (h *EventHandler) GetPublicByDisplayID(c echo.Context) error {
	displayID := c.Param("display_id")

	event, err := h.eventUsecase.GetPublicEventByDisplayID(c.Request().Context(), displayID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Event not found or not published")
	}

	return c.JSON(http.StatusOK, event)
}
