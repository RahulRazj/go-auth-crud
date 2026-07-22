package router

import (
	"github.com/RahulRazj/go-crud-auth/internal/handler"
	"github.com/RahulRazj/go-crud-auth/internal/middleware"
	"github.com/RahulRazj/go-crud-auth/internal/service"
	"github.com/labstack/echo/v4"
)

func registerAuthRoutes(r *echo.Echo, h *handler.Handlers, jwt *service.JWTService) {
	r.POST("/auth/signup", h.Auth.Signup)
	r.POST("/auth/login", h.Auth.Login)
	r.POST("/auth/refresh", h.Auth.RefreshToken)

	// Protected routes
	authMiddleware := middleware.JWTAuth(jwt)
	r.GET("/auth/me", h.Auth.Me, authMiddleware)
}

