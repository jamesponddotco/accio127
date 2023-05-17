package middleware

import "net/http"

// AcceptRequests returns a 405 Method Not Allowed if the request method is not
// GET, HEAD, or OPTIONS.
func AcceptRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" && r.Method != "HEAD" && r.Method != "OPTIONS" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		next.ServeHTTP(w, r)
	})
}
