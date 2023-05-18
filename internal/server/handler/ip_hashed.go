package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"git.sr.ht/~jamesponddotco/accio127/internal/database"
	"git.sr.ht/~jamesponddotco/accio127/internal/server/model"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// HashedIPHandler is an HTTP handler for the /ip/hashed endpoint.
type HashedIPHandler struct {
	db     *database.DB
	logger *zap.Logger
}

// NewHashedIPHandler returns a new HashedIPHandler instance.
func NewHashedIPHandler(db *database.DB, logger *zap.Logger) *HashedIPHandler {
	return &HashedIPHandler{
		db:     db,
		logger: logger,
	}
}

// ServeHTTP serves the /ip/hashed endpoint.
func (h *HashedIPHandler) Handle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ip, err := ClientIP(r)
	if err != nil {
		h.logger.Error("Failed to get client IP address", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	var (
		hashedIP    = HashIP(ip)
		contentType = r.Header.Get("Content-Type")
	)

	if contentType == "application/json" {
		w.Header().Set("Content-Type", "application/json")

		ipModel := model.NewIP(hashedIP)

		ipJSON, err := json.Marshal(ipModel) //nolint:errchkjson // if we don't check here, another linter complains
		if err != nil {
			h.logger.Error("Failed to marshal IP address to JSON", zap.Error(err))

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = w.Write(ipJSON)
		if err != nil {
			h.logger.Error("Failed to write IP address JSON to response", zap.Error(err))

			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	} else {
		w.Header().Set("Content-Type", "text/plain")

		_, err := w.Write([]byte(hashedIP))
		if err != nil {
			h.logger.Error("Failed to write IP address to response", zap.Error(err))

			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	go func() {
		_, err := h.db.Increment()
		if err != nil {
			h.logger.Error("Failed to increment access counter", zap.Error(err))
		}
	}()
}

// HashIP hashes an IP address using SHA256.
func HashIP(ip string) string {
	hash := sha256.Sum256([]byte(ip))

	return hex.EncodeToString(hash[:])
}
