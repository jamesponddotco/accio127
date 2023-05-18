package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// PrivacyPolicy adds a privacy policy header to the response.
func PrivacyPolicy(uri string, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Privacy-Policy", uri)

		next(w, r, ps)
	}
}

// SecureHeader adds basic security headers to the response.
func SecureHeader(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		w.Header().Set("Content-Security-Policy", "default-src 'none'")

		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-Content-Type-Options", "nosniff")

		next(w, r, ps)
	}
}

// CORS adds CORS headers to the response.
func CORS(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")

		if r.Method == http.MethodOptions {
			return
		}

		next(w, r, ps)
	}
}
