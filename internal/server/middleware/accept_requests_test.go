package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.sr.ht/~jamesponddotco/accio127/internal/server/middleware"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func TestAcceptRequests(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "get_method",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "head_method",
			method:         http.MethodHead,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "options_method",
			method:         http.MethodOptions,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "post_method",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, err := http.NewRequest(tt.method, "http://localhost/", http.NoBody)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			recorder := httptest.NewRecorder()

			router := httprouter.New()
			router.Handle(tt.method, "/", middleware.AcceptRequests(logger, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				w.Write([]byte("OK"))
			}))

			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d but got %d", tt.expectedStatus, recorder.Code)
			}
		})
	}
}
