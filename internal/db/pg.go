package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

type Pool = pgxpool.Pool

const (
	maxBackoffDuration = 5 * time.Second
	pingTimeout        = 3 * time.Second
	maxRetryAttempts   = 10
)

// buildConnectionString constructs a PostgreSQL connection string (DSN).
func buildConnectionString(host string, port int, user, password, databaseName, sslMode string) string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=%s", user, password, host, port, databaseName, sslMode)
}

// Connect connects to Postgres using a pgxpool.
func Connect(ctx context.Context, host string, port int, user, password, databaseName, sslMode string, maxConns, minConns int32) (*pgxpool.Pool, error) {
	connectionString := buildConnectionString(host, port, user, password, databaseName, sslMode)

	poolConfig, err := createPoolConfig(connectionString, maxConns, minConns)
	if err != nil {
		return nil, err
	}

	databasePool, err := createDatabasePool(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := waitForDatabaseReady(ctx, databasePool); err != nil {
		databasePool.Close()
		return nil, err
	}

	return databasePool, nil
}

// createPoolConfig creates a pgxpool configuration with connection limits.
func createPoolConfig(connectionString string, maxConns, minConns int32) (*pgxpool.Config, error) {
	poolConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse database config: %w", err)
	}
	poolConfig.MaxConns = maxConns
	poolConfig.MinConns = minConns
	return poolConfig, nil
}

// createDatabasePool creates a new database connection pool.
func createDatabasePool(ctx context.Context, poolConfig *pgxpool.Config) (*pgxpool.Pool, error) {
	databasePool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create database pool: %w", err)
	}
	return databasePool, nil
}

// waitForDatabaseReady waits for the database to be ready with exponential backoff.
func waitForDatabaseReady(ctx context.Context, databasePool *pgxpool.Pool) error {
	backoffDuration := time.Second
	var lastError error

	for range maxRetryAttempts {
		pingContext, cancel := context.WithTimeout(ctx, pingTimeout)
		lastError = databasePool.Ping(pingContext)
		cancel()

		if lastError == nil {
			return nil
		}

		time.Sleep(backoffDuration)
		if backoffDuration < maxBackoffDuration {
			backoffDuration *= 2
		}
	}

	return fmt.Errorf("database not ready after %d attempts: %w", maxRetryAttempts, lastError)
}
