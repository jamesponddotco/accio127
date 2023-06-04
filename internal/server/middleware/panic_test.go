package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.sr.ht/~jamesponddotco/accio127/internal/server/middleware"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

func TestPanicRecovery(t *testing.T) {
	t.Parallel()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	tests := []struct {
		name           string
		handler        httprouter.Handle
		expectPanic    bool
		expectedStatus int
	}{
		{
			name: "Handler with panic",
			handler: func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				panic("test panic")
			},
			expectPanic:    true,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "Handler without panic",
			handler: func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				w.Write([]byte("OK"))
			},
			expectPanic:    false,
			expectedStatus: http.StatusOK,
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

			recorder := httptest.NewRecorder()

			router := httprouter.New()
			router.GET("/", middleware.PanicRecovery(logger, tt.handler))

			router.ServeHTTP(recorder, req)

			if tt.expectPanic && recorder.Code != http.StatusInternalServerError {
				t.Errorf("Expected internal server error but got status code: %d", recorder.Code)
			}

			if !tt.expectPanic && recorder.Code != http.StatusOK {
				t.Errorf("Expected success response but got status code: %d", recorder.Code)
			}
		})
	}
}
