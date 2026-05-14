package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RahulRazj/go-crud-auth/internal/config"
	"github.com/rs/zerolog"
)

// Environment constants for logger
const (
	Development Environment = iota
	Production
	Staging
)

type Environment int

type LoggerService struct {
	logger zerolog.Logger
	file   *os.File
}

// New creates and returns a new LoggerService instance using the app config.
func NewLoggerService(cfg *config.Config) (*LoggerService, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	env := EnvironmentFromString(cfg.Primary.Env)
	logFilePath := cfg.Logging.File
	logLevel := cfg.Logging.Level

	var writer io.Writer = os.Stderr
	var logFile *os.File

	if env == Development {
		writer = zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		}
	} else if strings.TrimSpace(logFilePath) != "" {
		if err := os.MkdirAll(filepath.Dir(logFilePath), 0o755); err != nil {
			return nil, err
		}

		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, err
		}

		logFile = file
		writer = file
	}

	logger := zerolog.New(writer).
		With().
		Timestamp().
		Logger()

	logger = logger.Level(levelFromConfig(env, logLevel))

	return &LoggerService{logger: logger, file: logFile}, nil
}

func (ls *LoggerService) Shutdown() error {
	if ls.file != nil {
		return ls.file.Close()
	}
	return nil
}

func levelFromConfig(env Environment, level string) zerolog.Level {
	raw := strings.TrimSpace(strings.ToLower(level))
	if raw == "" {
		if env == Development {
			return zerolog.DebugLevel
		}
		return zerolog.InfoLevel
	}

	switch raw {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

// Info logs an informational message
func (ls *LoggerService) Info(msg string, fields ...interface{}) {
	ls.logger.Info().Fields(fieldsToMap(fields)).Msg(msg)
}

// Error logs an error message
func (ls *LoggerService) Error(msg string, err error, fields ...interface{}) {
	ctx := ls.logger.Error()
	if err != nil {
		ctx = ctx.Err(err)
	}
	ctx.Fields(fieldsToMap(fields)).Msg(msg)
}

// Debug logs a debug message
func (ls *LoggerService) Debug(msg string, fields ...interface{}) {
	ls.logger.Debug().Fields(fieldsToMap(fields)).Msg(msg)
}

// Warn logs a warning message
func (ls *LoggerService) Warn(msg string, fields ...interface{}) {
	ls.logger.Warn().Fields(fieldsToMap(fields)).Msg(msg)
}

// Fatal logs a fatal message and exits
func (ls *LoggerService) Fatal(msg string, err error, fields ...interface{}) {
	ctx := ls.logger.Fatal()
	if err != nil {
		ctx = ctx.Err(err)
	}
	ctx.Fields(fieldsToMap(fields)).Msg(msg)
}

func EnvironmentFromString(raw string) Environment {
	switch strings.TrimSpace(strings.ToLower(raw)) {
	case "prod", "production":
		return Production
	case "stage", "staging":
		return Staging
	case "dev", "development":
		return Development
	default:
		return Development
	}
}

// fieldsToMap converts variadic fields to a map for structured logging
func fieldsToMap(fields []interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for i := 0; i < len(fields); i += 2 {
		if i+1 < len(fields) {
			key, ok := fields[i].(string)
			if ok {
				result[key] = fields[i+1]
			}
		}
	}
	return result
}
