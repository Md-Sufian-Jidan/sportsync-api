package middleware

import (
	"fmt"
	"net/http"
	"sportsync-api/internal/auth"
	"sportsync-api/internal/httpResponse"
	"strings"

	"github.com/labstack/echo/v4"
)

func CheckAuth(jwtService auth.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, httpResponse.Error{
					Success: false,
					Message: "Unauthorized",
					Errors:  "Missing authorization header",
				})
			}

			var token string
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			} else if len(parts) == 1 {
				token = parts[0]
			} else {
				return c.JSON(http.StatusUnauthorized, httpResponse.Error{
					Success: false,
					Message: "Unauthorized",
					Errors:  "Invalid authorization header format",
				})
			}
			claims, err := jwtService.ValidateToken(token)

			fmt.Println("claims", claims)
			fmt.Println("err", err)

			if err != nil {
				return c.JSON(http.StatusUnauthorized, httpResponse.Error{
					Success: false,
					Message: "Unauthorized",
					Errors:  "Invalid or expired token",
				})
			}

			c.Set("user_id", claims.UserID)
			c.Set("user_name", claims.Name)
			c.Set("user_email", claims.Email)
			c.Set("user_role", claims.Role)

			return next(c)
		}
	}
}

func RequireRole(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("user_role").(string)
			if !ok || userRole != role {
				return c.JSON(http.StatusForbidden, httpResponse.Error{
					Success: false,
					Message: "Forbidden",
					Errors:  "Forbidden: insufficient permissions",
				})
			}
			return next(c)
		}
	}
}
