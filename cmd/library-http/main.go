package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	bookhttp "github.com/bkiran6398/library/internal/books/http"
	bookrepo "github.com/bkiran6398/library/internal/books/repository"
	booksvc "github.com/bkiran6398/library/internal/books/service"
	"github.com/bkiran6398/library/internal/config"
	"github.com/bkiran6398/library/internal/db"
	"github.com/bkiran6398/library/internal/http/router"
	"github.com/bkiran6398/library/internal/logger"
	"github.com/rs/zerolog"
)

const shutdownTimeout = 10 * time.Second

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

	bookHandler := initializeBookHandler(databasePool)

	routeHandler := initializeHTTPRouter(loggerInstance, configuration.Server.CORSAllowedOrigins, bookHandler)

	server := startHTTPServer(
		loggerInstance,
		routeHandler,
		configuration.Server.Port,
	)

	<-waitForShutdownSignal()
	shutdownServer(loggerInstance, server)
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

// initializeBookHandler creates and wires up the book handler with its dependencies.
func initializeBookHandler(databasePool *db.Pool) bookhttp.Handler {
	bookRepo := bookrepo.NewPgRepository(databasePool)
	bookSvc := booksvc.NewService(bookRepo)
	return bookhttp.NewHandler(bookSvc)
}

// initializeHTTPRouter creates and configures the HTTP router with all routes and middleware.
func initializeHTTPRouter(loggerInstance zerolog.Logger, allowedOrigins []string, bookHandler bookhttp.Handler) http.Handler {
	return router.NewRouter(
		loggerInstance,
		router.CORSConfig{AllowedOrigins: allowedOrigins},
		bookHandler,
	)
}

// startHTTPServer starts the HTTP server in a goroutine and returns it.
func startHTTPServer(logger zerolog.Logger, handler http.Handler, port int) *http.Server {
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: handler,
	}

	go func() {
		logger.Info().Str("address", server.Addr).Msg("HTTP server listening")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("HTTP server error")
		}
	}()

	return server
}

// waitForShutdownSignal waits for OS shutdown signals (SIGINT or SIGTERM).
func waitForShutdownSignal() <-chan os.Signal {
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)
	return shutdownSignal
}

// shutdownServer gracefully shuts down the HTTP server.
func shutdownServer(logger zerolog.Logger, server *http.Server) {
	logger.Info().Msg("shutting down server...")
	shutdownContext, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownContext); err != nil {
		logger.Error().Err(err).Msg("server shutdown error")
		return
	}
	logger.Info().Msg("server stopped")
}
