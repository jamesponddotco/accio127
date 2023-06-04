package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.sr.ht/~jamesponddotco/accio127/internal/server/middleware"
	"github.com/julienschmidt/httprouter"
)

func TestChain(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest(http.MethodGet, "http://localhost/", http.NoBody)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	var (
		recorder = httptest.NewRecorder()
		router   = httprouter.New()
		handler  = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
			w.Write([]byte("OK"))
		}
	)

	router.GET("/", middleware.Chain(handler,
		func(h httprouter.Handle) httprouter.Handle {
			return middleware.PrivacyPolicy("http://example.com/privacy", h)
		},
		middleware.SecureHeader,
	))

	router.ServeHTTP(recorder, req)

	expectedHeaders := map[string]string{
		"Privacy-Policy":            "http://example.com/privacy",
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
