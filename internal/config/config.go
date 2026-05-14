package config

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Primary  Primary        `koanf:"primary" validate:"required"`
	Server   ServerConfig   `koanf:"server" validate:"required"`
	Database DatabaseConfig `koanf:"database" validate:"required"`
	Logging  LoggingConfig  `koanf:"logging" validate:"required"`
}

type Primary struct {
	Env string `koanf:"env" validate:"required"`
}

type LoggingConfig struct {
	Level string `koanf:"level" validate:"required"`
	File  string `koanf:"file"`
}

type ServerConfig struct {
	Port               string   `koanf:"port" validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout" validate:"required"`
	WriteTimeout       int      `koanf:"write_timeout" validate:"required"`
	IdleTimeout        int      `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	Host            string `koanf:"host" validate:"required"`
	Port            int    `koanf:"port" validate:"required"`
	User            string `koanf:"user" validate:"required"`
	Password        string `koanf:"password"`
	Name            string `koanf:"name" validate:"required"`
	SSLMode         string `koanf:"ssl_mode" validate:"required"`
	MaxOpenConns    int    `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns    int    `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time" validate:"required"`
}

func LoadConfig() (*Config, error) {
	k := koanf.New(".")

	err := k.Load(env.Provider("GO_AUTH_CRUD_", ".", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, "GO_AUTH_CRUD_"))
	}), nil)

	if err != nil {
		return nil, err
	}

	mainConfig := &Config{}

	if err = k.Unmarshal("", mainConfig); err != nil {
		return nil, err
	}

	validate := validator.New()

	if err := validate.Struct(mainConfig); err != nil {
		return nil, err
	}

	mainConfig.Server.CORSAllowedOrigins = parseCORSOrigins(mainConfig.Server.CORSAllowedOrigins)

	return mainConfig, nil
}

func parseCORSOrigins(origins []string) []string {
	if len(origins) == 1 {
		parts := strings.Split(origins[0], ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(parts[i])
		}
		return parts
	}

	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	return origins
}