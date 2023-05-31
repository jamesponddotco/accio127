// Package errors provides custom error types for the application.
package errors

import (
	"encoding/json"
	"net/http"

	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"go.uber.org/zap"
)

const (
	// ErrNilLogger is returned when a nil logger is passed to a function.
	ErrNilLogger xerrors.Error = "logger cannot be nil"

	// ErrEmptyDSN is returned when an empty DSN is passed to a function.
	ErrEmptyDSN xerrors.Error = "dsn cannot be empty"
)

// ErrorResponse is the response returned by the API when an error occurs.
type ErrorResponse struct {
	// Message is a human-readable message describing the error.
	Message string `json:"message"`

	// Code is a machine-readable code describing the error.
	Code uint `json:"code"`
}

// JSON sends an ErrorResponse to the HTTP response writer as JSON.
func JSON(w http.ResponseWriter, logger *zap.Logger, response ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(response.Code))

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error("failed to encode error response", zap.Error(err))
	}
}
