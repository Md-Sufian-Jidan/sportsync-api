package middleware

import (
	"net/http"
	"sportsync-api/internal/auth"
	"strings"

	"github.com/labstack/echo/v5"
)

func CheckAuth(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c *echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Missing authorization header",
				})
			}

			parts := strings.Split(authHeader, "")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format",
				})
			}

			token := parts[1]
			claims, err := jwtService.ValidateToken(token)

			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"errro": "Invalid or expired token",
				})
			}

			c.Set("user_id", claims.ID)
			c.Set("user_name", claims.Name)
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}
