package router

import (
	"github.com/RahulRazj/go-crud-auth/internal/handler"
	"github.com/RahulRazj/go-crud-auth/internal/middleware"
	"github.com/RahulRazj/go-crud-auth/internal/server"
	"github.com/labstack/echo/v4"
)

func NewRouter(s *server.Server, h *handler.Handlers) *echo.Echo {
	middlewares := middleware.NewMiddlewares(s)
	router := echo.New()

	router.HTTPErrorHandler = middlewares.Global.GlobalErrorHandler

	// global middlewares
	router.Use(
		middlewares.Global.CORS(),
		middlewares.Global.Secure(),
		middleware.RequestID(),
		middlewares.Global.RequestLogger(),
		middlewares.Global.Recover(),
	)

	registerSystemRoutes(router, h)
	registerAuthRoutes(router, h)
	return router
}
