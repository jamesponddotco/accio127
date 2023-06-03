package handler

import (
	"encoding/json"
	"net/http"

	"git.sr.ht/~jamesponddotco/accio127/internal/build"
	"git.sr.ht/~jamesponddotco/accio127/internal/database"
	"git.sr.ht/~jamesponddotco/accio127/internal/errors"
	"git.sr.ht/~jamesponddotco/accio127/internal/server/model"
	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

const (
	// Online is the status of a service that is online.
	Online string = "Online"

	// Offline is the status of a service that is offline.
	Offline string = "Offline"
)

// HealthHandler is an HTTP handler for the /health endpoint.
type HealthHandler struct {
	db     *database.DB
	logger *zap.Logger
}

// NewHealthHandler creates a new HealthHandler instance.
func NewHealthHandler(db *database.DB, logger *zap.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		logger: logger,
	}
}

// ServeHTTP serves the /health endpoint.
func (h *HealthHandler) Handle(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	databaseStatus := Online

	err := h.db.Ping()
	if err != nil {
		databaseStatus = Offline

		h.logger.Warn("Database is offline", zap.Error(err))
	}

	var (
		dependencies = []model.Dependency{
			{
				Service: "sqlite",
				Status:  databaseStatus,
			},
		}
		status = model.NewHealth(build.Name, build.Version, dependencies)
	)

	statusJSON, err := json.Marshal(status) //nolint:errchkjson // if we don't check here, another linter complains
	if err != nil {
		h.logger.Error("Failed to marshal status to JSON", zap.Error(err))

		errors.JSON(w, h.logger, errors.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to marshal status to JSON.",
		})

		return
	}

	w.Header().Set(xhttp.ContentType, xhttp.ApplicationJSON)

	_, err = w.Write(statusJSON)
	if err != nil {
		h.logger.Error("Failed to write status JSON to response", zap.Error(err))

		errors.JSON(w, h.logger, errors.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to write status JSON to response.",
		})

		return
	}
}
