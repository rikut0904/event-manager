package web

import (
	"backend/internal/infrastructure/firebase"
	"backend/internal/interface/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(
	healthHandler *handler.HealthHandler,
	authHandler *handler.AuthHandler,
	fbClient *firebase.Client,
) *echo.Echo {
	e := echo.New()

	// Global Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Public routes
	e.GET("/health", healthHandler.HealthCheck)
	e.POST("/auth/signup", authHandler.SignUp)
	e.POST("/auth/login", authHandler.Login)

	// Protected routes
	r := e.Group("")
	r.Use(AuthMiddleware(fbClient))
	r.POST("/auth/link-connpass", authHandler.LinkConnpass)

	return e
}
