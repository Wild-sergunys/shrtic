package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Wild-sergunys/shrtic/internal/model"
)

func writeError(w http.ResponseWriter, status int, errorType, message string, details any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(model.ErrorResponse{
		Error:   errorType,
		Message: message,
		Details: details,
	})
}
