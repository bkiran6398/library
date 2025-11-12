package response

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func JSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func Error(w http.ResponseWriter, status int, code, message string, details interface{}) {
	JSON(w, status, map[string]interface{}{
		"error": ErrorBody{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

