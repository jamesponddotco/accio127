package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// AcceptRequests returns a 405 Method Not Allowed if the request method is not
// GET, HEAD, or OPTIONS.
func AcceptRequests(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r.Method != "GET" && r.Method != "HEAD" && r.Method != "OPTIONS" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		next(w, r, ps)
	}
}
