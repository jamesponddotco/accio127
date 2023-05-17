package handler

import (
	"encoding/json"
	"net/http"

	"git.sr.ht/~jamesponddotco/accio127/internal/database"
	"git.sr.ht/~jamesponddotco/accio127/internal/server/model"
	"go.uber.org/zap"
)

// MetricsHandler is an HTTP handler for the /metrics endpoint.
type MetricsHandler struct {
	db     *database.DB
	logger *zap.Logger
}

// NewMetricsHandler creates a new MetricsHandler instance.
func NewMetricsHandler(db *database.DB, logger *zap.Logger) *MetricsHandler {
	return &MetricsHandler{
		db:     db,
		logger: logger,
	}
}

// ServeHTTP serves the /metrics endpoint.
func (h *MetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		count   = h.db.Count()
		counter = model.NewCounter(count)
	)

	counterJSON, err := json.Marshal(counter) //nolint:errchkjson // if we don't check here, another linter complains
	if err != nil {
		h.logger.Error("Failed to marshal access counter to JSON", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(counterJSON)
	if err != nil {
		h.logger.Error("Failed to write access counter JSON to response", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}
}
