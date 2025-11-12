package response

import (
	"errors"
	"net/http"

	intErr "github.com/bkiran6398/library/internal/errors"
)

// MapServiceErrorToHTTP maps service layer errors to appropriate HTTP status codes and responses.
// This is a generic utility that can be used across all transport layers.
func MapServiceErrorToHTTP(w http.ResponseWriter, serviceError error) {
	if serviceError == nil {
		return
	}

	switch {
	case errors.Is(serviceError, intErr.ErrNotFound):
		Error(w, http.StatusNotFound, "not_found", "Resource not found", nil)
	case errors.Is(serviceError, intErr.ErrConflict):
		Error(w, http.StatusConflict, "conflict", "Resource conflict", nil)
	case errors.Is(serviceError, intErr.ErrBadRequest):
		Error(w, http.StatusBadRequest, "bad_request", serviceError.Error(), nil)
	default:
		// For unknown errors, return generic internal server error
		Error(w, http.StatusInternalServerError, "internal_error", "An internal error occurred", nil)
	}
}
