package middleware

import (
	"github.com/RahulRazj/go-crud-auth/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

const (
	UserIDKey   = "user_id"
	UserRoleKey = "user_role"
	LoggerKey   = "logger"
)

type ContextEnhancer struct {
	server *server.Server
}

func NewContextEnhancer(s *server.Server) *ContextEnhancer {
	return &ContextEnhancer{server: s}
}

func (ce *ContextEnhancer) extractUserID(c echo.Context) string {
	// Check if user_id was already set by auth middleware (Clerk)
	if userID, ok := c.Get("user_id").(string); ok && userID != "" {
		return userID
	}
	return ""
}

func GetUserID(c echo.Context) string {
	if userID, ok := c.Get(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

func GetLogger(c echo.Context) *zerolog.Logger {
	if logger, ok := c.Get(LoggerKey).(*zerolog.Logger); ok {
		return logger
	}
	// Fallback to a basic logger if not found
	logger := zerolog.Nop()
	return &logger
}
