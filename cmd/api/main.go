package main

import (
	"github.com/RahulRazj/go-crud-auth/internal/config"
	"github.com/RahulRazj/go-crud-auth/internal/logger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	loggerService, err := logger.NewLoggerService(cfg)
	if err != nil {
		panic("failed to create logger service: " + err.Error())
	}
	defer loggerService.Shutdown()

	loggerService.Info("Config loaded successfully", "port", cfg.Server.Port)
}
