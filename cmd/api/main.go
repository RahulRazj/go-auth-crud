package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/RahulRazj/go-crud-auth/internal/config"
	"github.com/RahulRazj/go-crud-auth/internal/database"
	"github.com/RahulRazj/go-crud-auth/internal/handler"
	"github.com/RahulRazj/go-crud-auth/internal/logger"
	"github.com/RahulRazj/go-crud-auth/internal/repository"
	"github.com/RahulRazj/go-crud-auth/internal/router"
	"github.com/RahulRazj/go-crud-auth/internal/server"
)

const DefaultContextTimeout = 30

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	logger, err := logger.NewLoggerService(cfg)
	if err != nil {
		panic("failed to create logger service: " + err.Error())
	}
	defer logger.Shutdown()

	if cfg.Primary.Env != "local" {
		if err := database.Migrate(context.Background(), logger, cfg); err != nil {
			logger.Fatal("failed to migrate database: %w", err)
		}
	}

	// Initialize server
	srv, err := server.New(cfg, logger)
	if err != nil {
		logger.Fatal("failed to initialize server: %w", err)
	}

	baseRepository := repository.NewBaseRepository(srv.DB)
	userRepository := repository.NewUserRepository(baseRepository)
	handlers := handler.NewHandlers(srv, userRepository)

	// Initialize router
	r := router.NewRouter(srv, handlers)

	// Setup HTTP server
	srv.SetupHTTPServer(r)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	// Start server
	go func() {
		if err = srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("failed to start server: %w", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)

	if err = srv.Shutdown(ctx); err != nil {
		logger.Fatal("server forced to shutdown: %w", err)
	}

	stop()
	cancel()

	logger.Info("server exited properly")
}
