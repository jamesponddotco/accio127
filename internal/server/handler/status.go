package handler

import (
	"encoding/json"
	"net/http"

	"git.sr.ht/~jamesponddotco/accio127/internal/database"
	"git.sr.ht/~jamesponddotco/accio127/internal/server/model"
	"go.uber.org/zap"
)

const (
	// Online is the status of a server that is online.
	Online string = "Online"

	// Offline is the status of a server that is offline.
	Offline string = "Offline"
)

// StatusHandler is an HTTP handler for the /status endpoint.
type StatusHandler struct {
	db     *database.DB
	logger *zap.Logger
}

// NewStatusHandler creates a new StatusHandler.
func NewStatusHandler(db *database.DB, logger *zap.Logger) *StatusHandler {
	return &StatusHandler{
		db:     db,
		logger: logger,
	}
}

// ServeHTTP serves the /status endpoint.
func (h *StatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	const serverStatus = Online

	databaseStatus := Online

	err := h.db.Ping()
	if err != nil {
		databaseStatus = Offline

		h.logger.Warn("Database is offline", zap.Error(err))
	}

	status := model.NewStatus(serverStatus, databaseStatus)

	statusJSON, err := json.Marshal(status) //nolint:errchkjson // if we don't check here, another linter complains
	if err != nil {
		h.logger.Error("Failed to marshal status to JSON", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(statusJSON)
	if err != nil {
		h.logger.Error("Failed to write status JSON to response", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
