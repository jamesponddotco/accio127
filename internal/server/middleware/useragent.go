package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// UserAgent ensures that the request has a valid user agent.
func UserAgent(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r.UserAgent() == "" {
			http.Error(w, "User agent is missing.", http.StatusBadRequest)
			return
		}

		next(w, r, ps)
	}
}
