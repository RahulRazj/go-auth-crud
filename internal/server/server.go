package server

import (
	"fmt"

	"github.com/RahulRazj/go-crud-auth/internal/config"
	"github.com/RahulRazj/go-crud-auth/internal/database"
	"github.com/RahulRazj/go-crud-auth/internal/logger"
)

type Server struct {
	Config *config.Config
	Logger *logger.LoggerService
	DB *database.Database
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
