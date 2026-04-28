package web

import (
	"net/http"
	"strings"

	"backend/internal/infrastructure/firebase"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(fbClient *firebase.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
			}

			idToken := strings.Replace(authHeader, "Bearer ", "", 1)
			token, err := fbClient.Auth.VerifyIDToken(c.Request().Context(), idToken)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
			}

			// ContextにuserIDをセット
			c.Set("userID", token.UID)
			return next(c)
		}
	}
}
