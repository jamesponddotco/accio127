package middleware

import (
	"net/http"

	"git.sr.ht/~jamesponddotco/accio127/internal/errors"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// UserAgent ensures that the request has a valid user agent.
func UserAgent(logger *zap.Logger, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r.UserAgent() == "" {
			errors.JSON(w, logger, errors.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "User agent is missing. Please provide a valid user agent.",
			})

			return
		}

		next(w, r, ps)
	}
}
