package handler

import (
	"github.com/RahulRazj/go-crud-auth/internal/repository"
	"github.com/RahulRazj/go-crud-auth/internal/server"
)

type Handler struct {
	server *server.Server
}

func NewHandler(s *server.Server) Handler {
	return Handler{
		server: s,
	}
}

type Handlers struct {
	Health  *HealthHandler
	OpenAPI *OpenAPIHandler
	Auth    *AuthHandler
}

func NewHandlers(s *server.Server, users *repository.UserRepository) *Handlers {
	return &Handlers{
		Health:  NewHealthHandler(s),
		OpenAPI: NewOpenAPIHandler(s),
		Auth:    NewAuthHandler(s, users),
	}
}
