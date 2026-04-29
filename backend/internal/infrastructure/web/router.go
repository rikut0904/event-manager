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
	allowedOrigins := []string{"http://localhost:3000"} // Default
	if allowedOriginsStr != "" {
		allowedOrigins = strings.Split(allowedOriginsStr, ",")
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	e.GET("/health", healthHandler.HealthCheck)
	e.POST("/auth/signup", authHandler.SignUp)
	e.POST("/auth/login", authHandler.Login)

	// Protected routes
	r := e.Group("/api/v1")
	r.Use(AuthMiddleware(fbClient))
	
	// Event routes
	r.POST("/events", eventHandler.Create)
	r.GET("/events/me", eventHandler.GetMyEvents)

	return e
}
