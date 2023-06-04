package errors_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.sr.ht/~jamesponddotco/accio127/internal/errors"
	"go.uber.org/zap"
)

func TestJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		response errors.ErrorResponse
		wantBody string
	}{
		{
			name: "valid_error_response",
			response: errors.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "User agent is missing. Please provide a valid user agent.",
			},
			wantBody: `{"message":"User agent is missing. Please provide a valid user agent.","code":400}`,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			recorder := httptest.NewRecorder()

			logger, err := zap.NewProduction()
			if err != nil {
				t.Fatalf("Failed to create logger: %v", err)
			}

			errors.JSON(recorder, logger, tt.response)

			if got := recorder.Body.String(); got != tt.wantBody+"\n" {
				t.Errorf("Expected body %q, but got %q", tt.wantBody, got)
			}

			if got, want := recorder.Header().Get("Content-Type"), "application/json"; got != want {
				t.Errorf("Expected Content-Type %s, but got %s", want, got)
			}

			if got, want := recorder.Code, int(tt.response.Code); got != want {
				t.Errorf("Expected status code %d, but got %d", want, got)
			}
		})
	}
}
