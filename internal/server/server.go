package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/RahulRazj/go-crud-auth/internal/config"
	"github.com/RahulRazj/go-crud-auth/internal/database"
	"github.com/RahulRazj/go-crud-auth/internal/logger"
)

type Server struct {
	Config     *config.Config
	Logger     *logger.LoggerService
	DB         *database.Database
	httpServer *http.Server
}

func New(cfg *config.Config, logger *logger.LoggerService) (*Server, error) {
	db, err := database.New(cfg, logger)

	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &Server{
		Config: cfg,
		Logger: logger,
		DB:     db,
	}, nil
}

func (s *Server) SetupHTTPServer(handler http.Handler) {
	s.httpServer = &http.Server{
		Addr:         ":" + s.Config.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Duration(s.Config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Config.Server.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Server.IdleTimeout) * time.Second,
	}
}

func (s *Server) Start() error {
	if s.httpServer == nil {
		return errors.New("HTTP server not initialized")
	}

	s.Logger.Info("starting server", "port", s.Config.Server.Port, "env", s.Config.Primary.Env)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	if err := s.DB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}

	return nil
}
