package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/RahulRazj/go-crud-auth/internal/config"
	"github.com/RahulRazj/go-crud-auth/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Pool   *pgxpool.Pool
	logger *logger.LoggerService
}

const DatabasePingTimeout = 10

func New(config *config.Config, logger *logger.LoggerService) (*Database, error) {
	if logger == nil {
		return nil, fmt.Errorf("logger must not be nil")
	}

	hostPort := net.JoinHostPort(config.Database.Host, strconv.Itoa(config.Database.Port))

	// URL-encode the password
	encodedPassword := url.QueryEscape(config.Database.Password)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s",
		config.Database.User,
		encodedPassword,
		hostPort,
		config.Database.Name,
		config.Database.SSLMode,
	)

	pgxPoolConfig, err := pgxpool.ParseConfig(dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx pool config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxPoolConfig)

	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	database := &Database{
		Pool:   pool,
		logger: logger,
	}

	ctx, cancel := context.WithTimeout(context.Background(), DatabasePingTimeout*time.Second)
	defer cancel()

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("connected to the database")

	return database, nil
}

func (db *Database) Close() error {
	db.logger.Info("closing database connection")

	db.Pool.Close()
	return nil
}
