package api

import (
	"encoding/json"
	"net/http"
)

// OKResponse writes a JSON response with status 200 OK
func OKResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

// ErrorResponse writes a JSON response with a given status code and error message
func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := map[string]any{
		"error": message,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
