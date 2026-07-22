package middleware

import (
	"strings"

	"github.com/RahulRazj/go-crud-auth/internal/errs"
	"github.com/RahulRazj/go-crud-auth/internal/service"
	"github.com/labstack/echo/v4"
)

const (
	ContextUserIDKey = "user_id"
	ContextEmailKey  = "user_email"
)

func JWTAuth(jwtService *service.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return errs.NewUnauthorizedError("Missing Authorization header", false)
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return errs.NewUnauthorizedError("Invalid Authorization header format", false)
			}

			tokenStr := strings.TrimSpace(parts[1])
			if tokenStr == "" {
				return errs.NewUnauthorizedError("Token is required", false)
			}

			claims, err := jwtService.ValidateToken(tokenStr, service.TokenTypeAccess)
			if err != nil {
				return errs.NewUnauthorizedError("Invalid or expired access token", false)
			}

			c.Set(ContextUserIDKey, claims.UserID)
			c.Set(ContextEmailKey, claims.Email)

			return next(c)
		}
	}
}
