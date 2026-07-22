package router

import (
	"github.com/RahulRazj/go-crud-auth/internal/handler"
	"github.com/labstack/echo/v4"
)

func registerAuthRoutes(r *echo.Echo, h *handler.Handlers) {
	r.POST("/auth/signup", h.Auth.Signup)
	r.POST("/auth/login", h.Auth.Login)
}
