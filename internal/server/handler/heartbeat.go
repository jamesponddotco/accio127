package handler

import (
	"net/http"

	"git.sr.ht/~jamesponddotco/accio127/internal/errors"
	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// HeartbeatHandler is an HTTP handler for the /heartbeat endpoint.
type HeartbeatHandler struct {
	logger *zap.Logger
}

// NewHeartbeatHandler creates a new HeartbeatHandler instance.
func NewHeartbeatHandler(logger *zap.Logger) *HeartbeatHandler {
	return &HeartbeatHandler{
		logger: logger,
	}
}

// ServeHTTP serves the /heartbeat endpoint.
func (h *HeartbeatHandler) Handle(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Header().Set(xhttp.ContentType, xhttp.TextPlain)

	_, err := w.Write([]byte("pong"))
	if err != nil {
		h.logger.Error("Failed to write response", zap.Error(err))

		errors.JSON(w, h.logger, errors.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to write response. Please try again later.",
		})

		return
	}
}
