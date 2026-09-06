package web

import (
	"net/http"

	"backend/internal/infrastructure/session"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(appSession *session.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if appSession == nil {
				return echo.NewHTTPError(http.StatusServiceUnavailable, "アプリセッションが設定されていません")
			}
			cookie, err := c.Cookie(session.CookieName)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "ログインが必要です")
			}
			userID, err := appSession.Verify(cookie.Value)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "ログインセッションが無効です")
			}
			c.Set("userID", userID)
			return next(c)
		}
	}
}

func OptionalAuthMiddleware(appSession *session.Manager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if appSession == nil {
				return next(c)
			}
			cookie, err := c.Cookie(session.CookieName)
			if err == nil {
				if userID, verifyErr := appSession.Verify(cookie.Value); verifyErr == nil {
					c.Set("userID", userID)
				}
			}
			return next(c)
		}
	}
}
