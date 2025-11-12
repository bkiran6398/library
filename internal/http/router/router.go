package router

import (
	"net/http"

	bookhttp "github.com/bkiran6398/library/internal/books/http"
	"github.com/bkiran6398/library/internal/http/middleware"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type CORSConfig struct {
	AllowedOrigins []string
}

func NewRouter(loggerInstance zerolog.Logger, corsConfig CORSConfig, bookHandler bookhttp.Handler) http.Handler {
	router := mux.NewRouter()

	// Apply global middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Recovery(loggerInstance))
	router.Use(middleware.Logging(loggerInstance))

	// Register routes
	registerHealthEndpoint(router)

	apiRouter := router.PathPrefix("/v1").Subrouter()
	registerBookRoutes(apiRouter, bookHandler)

	// Apply CORS
	corsHandler := configureCORS(corsConfig.AllowedOrigins)
	return corsHandler(router)
}

// configureCORS creates and returns a CORS handler with the specified configuration.
func configureCORS(allowedOrigins []string) func(http.Handler) http.Handler {
	return handlers.CORS(
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Request-ID"}),
	)
}
