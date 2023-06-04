package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.sr.ht/~jamesponddotco/accio127/internal/server/middleware"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func TestUserAgent(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	tests := []struct {
		name        string
		userAgent   string
		expectError bool
	}{
		{
			name:        "request_with_user_agent",
			userAgent:   "test-agent",
			expectError: false,
		},
		{
			name:        "request_without_user_agent",
			userAgent:   "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(http.MethodGet, "http://localhost/", http.NoBody)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			req.Header.Set("User-Agent", tt.userAgent)

			recorder := httptest.NewRecorder()

			router := httprouter.New()
			router.GET("/", middleware.UserAgent(logger, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				w.Write([]byte("OK"))
			}))

			router.ServeHTTP(recorder, req)

			if tt.expectError && recorder.Code == http.StatusOK {
				t.Errorf("Expected error but got success response")
			}

			if !tt.expectError && recorder.Code != http.StatusOK {
				t.Errorf("Expected success response but got error, status code: %d", recorder.Code)
			}
		})
	}
}
