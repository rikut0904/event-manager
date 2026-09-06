package web

import (
	"net/http"

	"backend/internal/infrastructure/commonid"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(commonID *commonid.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if commonID == nil {
				return echo.NewHTTPError(http.StatusServiceUnavailable, "Common IDが設定されていません")
			}
			cookie, err := c.Cookie("common_id_session")
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "ログインが必要です")
			}
			user, err := commonID.CheckSession(c.Request().Context(), cookie.Value)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "ログインセッションが無効です")
			}
			c.Set("userID", user.CommonUserID)
			return next(c)
		}
	}
}

func OptionalAuthMiddleware(commonID *commonid.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if commonID == nil {
				return next(c)
			}
			cookie, err := c.Cookie("common_id_session")
			if err == nil {
				if user, checkErr := commonID.CheckSession(c.Request().Context(), cookie.Value); checkErr == nil {
					c.Set("userID", user.CommonUserID)
				}
			}
			return next(c)
		}
	}
}
