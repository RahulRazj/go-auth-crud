package router

import (
	"os"
	"path/filepath"

	"github.com/RahulRazj/go-crud-auth/internal/handler"

	"github.com/labstack/echo/v4"
)

func registerSystemRoutes(r *echo.Echo, h *handler.Handlers) {
	r.GET("/status", h.Health.CheckHealth)

	r.Static("/static", staticFilePath())

	r.GET("/docs", h.OpenAPI.ServeOpenAPIUI)
}

func staticFilePath() string {
	wd, _ := os.Getwd()
	return filepath.Join(wd, "internal", "static")
}
