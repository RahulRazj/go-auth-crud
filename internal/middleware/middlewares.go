package middleware

import (
	"github.com/RahulRazj/go-crud-auth/internal/server"
)

type Middlewares struct {
	Global          *GlobalMiddlewares
	ContextEnhancer *ContextEnhancer
}

func NewMiddlewares(s *server.Server) *Middlewares {
	return &Middlewares{
		Global:          NewGlobalMiddlewares(s),
		ContextEnhancer: NewContextEnhancer(s),
	}
}
