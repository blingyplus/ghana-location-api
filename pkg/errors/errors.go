package errors

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound    = errors.New("resource not found")
	ErrInvalidSlug = errors.New("invalid slug format")
)

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error":"` + message + `"}`))
}
