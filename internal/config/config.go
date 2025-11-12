package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type LogConfig struct {
	Level string
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
	MaxConns int32
	MinConns int32
}

type Config struct {
	Log LogConfig
	DB     DBConfig
}

// Load loads configuration from config/config.yaml, allowing environment variables to override values.
// Environment variables are prefixed with LIB_ and use _ instead of dots (e.g., LIB_DB_HOST).
func Load() (*Config, error) {
	viperInstance := setupViperConfig()
	setDefaultValues(viperInstance)
	setupEnvironmentOverrides(viperInstance)

	if err := viperInstance.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	return unmarshalConfig(viperInstance)
}

// unmarshalConfig unmarshals Viper configuration into Config struct.
func unmarshalConfig(viperInstance *viper.Viper) (*Config, error) {
	var configuration Config
	if err := viperInstance.Unmarshal(&configuration); err != nil {
		return nil, fmt.Errorf("failed unmarshalling config: %w", err)
	}
	return &configuration, nil
}
