package config

import (
	"strings"

	"github.com/spf13/viper"
)

// setupViperConfig initializes Viper with configuration file settings.
func setupViperConfig() *viper.Viper {
	viperInstance := viper.New()
	viperInstance.SetConfigName("config")
	viperInstance.SetConfigType("yaml")
	viperInstance.AddConfigPath("config")

	// Custom config file path via LIB_CONFIG
	if customConfigPath := strings.TrimSpace(viper.GetString("LIB_CONFIG")); customConfigPath != "" {
		viperInstance.SetConfigFile(customConfigPath)
	}

	return viperInstance
}

// setDefaultValues sets default configuration values.
func setDefaultValues(viperInstance *viper.Viper) {
	// Log defaults
	viperInstance.SetDefault("log.level", "info")

	// Database defaults
	viperInstance.SetDefault("db.host", "localhost")
	viperInstance.SetDefault("db.port", 5432)
	viperInstance.SetDefault("db.user", "library")
	viperInstance.SetDefault("db.password", "secret")
	viperInstance.SetDefault("db.name", "library")
	viperInstance.SetDefault("db.sslmode", "disable")
	viperInstance.SetDefault("db.max_conns", 10)
	viperInstance.SetDefault("db.min_conns", 1)

	// Server defaults
	viperInstance.SetDefault("server.port", 8080)
	viperInstance.SetDefault("server.cors_allowed_origins", []string{"*"})
}

// setupEnvironmentOverrides configures Viper to read from environment variables.
func setupEnvironmentOverrides(viperInstance *viper.Viper) {
	viperInstance.AutomaticEnv()
	viperInstance.SetEnvPrefix("LIB")
	viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
}
