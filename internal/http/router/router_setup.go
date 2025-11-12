package router

import (
	"net/http"

	bookhttp "github.com/bkiran6398/library/internal/books/http"
	"github.com/gorilla/mux"
)

// registerHealthEndpoint registers the health check endpoint.
func registerHealthEndpoint(router *mux.Router) {
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)
}

// registerBookRoutes registers all book-related API routes.
func registerBookRoutes(apiRouter *mux.Router, bookHandler bookhttp.Handler) {
	apiRouter.HandleFunc("/books", bookHandler.List).Methods(http.MethodGet)
	apiRouter.HandleFunc("/books", bookHandler.Create).Methods(http.MethodPost)
	apiRouter.HandleFunc("/books/{id}", bookHandler.Get).Methods(http.MethodGet)
	apiRouter.HandleFunc("/books/{id}", bookHandler.Update).Methods(http.MethodPut)
	apiRouter.HandleFunc("/books/{id}", bookHandler.Delete).Methods(http.MethodDelete)
}
