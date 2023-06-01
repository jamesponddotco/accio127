package middleware

import (
	"fmt"
	"net/http"

	"git.sr.ht/~jamesponddotco/accio127/internal/errors"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// AcceptRequests returns a 405 Method Not Allowed if the request method is not
// GET, HEAD, or OPTIONS.
func AcceptRequests(logger *zap.Logger, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r.Method != "GET" && r.Method != "HEAD" && r.Method != "OPTIONS" {
			errors.JSON(w, logger, errors.ErrorResponse{
				Code:    http.StatusMethodNotAllowed,
				Message: fmt.Sprintf("Method %s not allowed. Must be GET, HEAD, or OPTIONS.", r.Method),
			})

			return
		}

		next(w, r, ps)
	}
}
