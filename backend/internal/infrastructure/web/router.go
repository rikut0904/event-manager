package web

import (
	"os"
	"strings"

	"backend/internal/infrastructure/firebase"
	"backend/internal/interface/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	eventHandler *handler.EventHandler,
	fbClient *firebase.Client,
) *echo.Echo {
	e := echo.New()

	// Get allowed origins from env
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	allowedOrigins := []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	if allowedOriginsStr != "" {
		envOrigins := strings.Split(allowedOriginsStr, ",")
		allowedOrigins = append(allowedOrigins, envOrigins...)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.PATCH, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.GET("/health", healthHandler.HealthCheck)
	e.POST("/auth/signup", authHandler.SignUp)
	e.POST("/auth/login", authHandler.Login)
	e.POST("/auth/link-connpass", authHandler.LinkConnpass, AuthMiddleware(fbClient))

	// 公開閲覧用
	e.GET("/api/v1/events/published", eventHandler.GetPublished)
	e.GET("/api/v1/events/view/:id", eventHandler.GetByID, OptionalAuthMiddleware(fbClient))
	e.GET("/api/v1/events/public/:display_id", eventHandler.GetPublicByDisplayID)

	// Protected routes
	r := e.Group("/api/v1")
	r.Use(AuthMiddleware(fbClient))

	// Event routes
	r.POST("/events", eventHandler.Create)
	r.PATCH("/events/:id", eventHandler.Update)
	r.DELETE("/events/:id", eventHandler.Delete)
	r.GET("/events/me", eventHandler.GetMyEvents)
	r.GET("/events/:id", eventHandler.GetByID)

	return e
}
