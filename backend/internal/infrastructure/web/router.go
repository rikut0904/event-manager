package web

import (
	"strings"

	"backend/internal/infrastructure/commonid"
	"backend/internal/infrastructure/session"
	"backend/internal/interface/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	eventHandler *handler.EventHandler,
	commonID *commonid.Client,
	appOrigin string,
	appSession *session.Manager,
) *echo.Echo {
	e := echo.New()

	// The same canonical app origin is used for the frontend redirect and CORS.
	allowedOrigins := []string{strings.TrimSpace(appOrigin)}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.GET("/health", healthHandler.HealthCheck)
	e.GET("/auth/callback", authHandler.Callback)
	e.GET("/auth/logout", authHandler.BeginLogout)
	e.GET("/auth/logout/callback", authHandler.LogoutCallback)
	e.GET("/auth/:intent", authHandler.Begin)
	e.POST("/auth/link-connpass", authHandler.LinkConnpass, AuthMiddleware(appSession))

	// 公開閲覧用
	e.GET("/api/v1/events/published", eventHandler.GetPublished)
	e.GET("/api/v1/events/view/:id", eventHandler.GetByID, OptionalAuthMiddleware(appSession))
	e.GET("/api/v1/events/public/:display_id", eventHandler.GetPublicByDisplayID)

	// Protected routes
	r := e.Group("/api/v1")
	r.Use(AuthMiddleware(appSession))
	r.GET("/users/me", authHandler.CurrentUser)

	// Event routes
	r.POST("/events", eventHandler.Create)
	r.PATCH("/events/:id", eventHandler.Update)
	r.DELETE("/events/:id", eventHandler.Delete)
	r.GET("/events/me", eventHandler.GetMyEvents)
	r.GET("/events/:id", eventHandler.GetByID)

	return e
}
