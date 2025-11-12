package main

import (
	"fmt"

	"github.com/bkiran6398/library/internal/config"
	"github.com/bkiran6398/library/internal/logger"
)

func main() {
	configuration, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load configuration: %w", err))
	}

	loggerInstance, err := logger.New(configuration.Log.Level)
	if err != nil {
		panic(fmt.Errorf("failed to create logger: %w", err))
	}
	loggerInstance.Info().Msg("starting library API")
}
