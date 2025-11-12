package main

import (
	"context"
	"fmt"

	"github.com/bkiran6398/library/internal/config"
	"github.com/bkiran6398/library/internal/db"
	"github.com/bkiran6398/library/internal/logger"
	"github.com/rs/zerolog"
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

	databasePool, err := initializeDatabase(context.Background(), loggerInstance, configuration.DB)
	if err != nil {
		loggerInstance.Fatal().Err(err).Msg("failed to initialize database")
	}
	defer databasePool.Close()
}

// initializeDatabase connects to the database and runs migrations.
func initializeDatabase(ctx context.Context, loggerInstance zerolog.Logger, dbConfig config.DBConfig) (*db.Pool, error) {
	databasePool, err := db.ConnectAndMigrate(
		ctx,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Name,
		dbConfig.SSLMode,
		dbConfig.MaxConns,
		dbConfig.MinConns,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	loggerInstance.Info().Msg("database connected and migrations completed")
	return databasePool, nil
}
