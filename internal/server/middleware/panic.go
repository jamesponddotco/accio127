package middleware

import (
	"net/http"

	"git.sr.ht/~jamesponddotco/accio127/internal/errors"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// PanicRecovery tries to recover from panics and returns a 500 error if there
// was one.
func PanicRecovery(logger *zap.Logger, next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered", zap.Any("error", err))

				errors.JSON(w, logger, errors.ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "Internal server error. Please try again later.",
				})
			}
		}()

		next(w, r, ps)
	}
}
