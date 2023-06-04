package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.sr.ht/~jamesponddotco/accio127/internal/server/middleware"
	"github.com/julienschmidt/httprouter"
)

func TestPrivacyPolicy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		uri  string
	}{
		{
			name: "Test privacy policy header",
			uri:  "http://example.com/privacy",
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
			router.GET("/", middleware.PrivacyPolicy(tt.uri, func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				w.Write([]byte("OK"))
			}))

			router.ServeHTTP(recorder, req)

			if got := recorder.Header().Get("Privacy-Policy"); got != tt.uri {
				t.Errorf("Expected Privacy-Policy header %s, but got %s", tt.uri, got)
			}
		})
	}
}

func TestSecureHeader(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodGet, "http://localhost/", http.NoBody)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	recorder := httptest.NewRecorder()

	router := httprouter.New()
	router.GET("/", middleware.SecureHeader(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Write([]byte("OK"))
	}))

	router.ServeHTTP(recorder, req)

	expectedHeaders := map[string]string{
		"Strict-Transport-Security": "max-age=63072000; includeSubDomains",
		"Content-Security-Policy":   "default-src 'none'",
		"X-Frame-Options":           "deny",
		"X-Content-Type-Options":    "nosniff",
	}

	for header, expected := range expectedHeaders {
		if got := recorder.Header().Get(header); got != expected {
			t.Errorf("Expected %s header %s, but got %s", header, expected, got)
		}
	}
}

func TestCORS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		method string
	}{
		{
			name:   "get_method",
			method: http.MethodGet,
		},
		{
			name:   "options_method",
			method: http.MethodOptions,
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
			router.Handle(tt.method, "/", middleware.CORS(func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
				w.Write([]byte("OK"))
			}))

			router.ServeHTTP(recorder, req)

			expectedHeaders := map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, HEAD, OPTIONS",
				"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding",
			}

			for header, expected := range expectedHeaders {
				if got := recorder.Header().Get(header); got != expected {
					t.Errorf("Expected %s header %s, but got %s", header, expected, got)
				}
			}
		})
	}
}
