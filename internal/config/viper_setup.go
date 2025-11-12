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
}

// setupEnvironmentOverrides configures Viper to read from environment variables.
func setupEnvironmentOverrides(viperInstance *viper.Viper) {
	viperInstance.AutomaticEnv()
	viperInstance.SetEnvPrefix("LIB")
	viperInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
}
