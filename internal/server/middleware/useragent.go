package middleware

import "net/http"

// UserAgent ensures that the request has a valid user agent.
func UserAgent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.UserAgent() == "" {
			http.Error(w, "User agent is missing.", http.StatusBadRequest)

			return
		}

		next.ServeHTTP(w, r)
	})
}
