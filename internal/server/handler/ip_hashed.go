package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"git.sr.ht/~jamesponddotco/accio127/internal/config"
	"git.sr.ht/~jamesponddotco/accio127/internal/database"
	"git.sr.ht/~jamesponddotco/accio127/internal/errors"
	"git.sr.ht/~jamesponddotco/accio127/internal/server/model"
	"git.sr.ht/~jamesponddotco/xstd-go/xnet/xhttp"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// HashedIPHandler is an HTTP handler for the /ip/hashed endpoint.
type HashedIPHandler struct {
	cfg    *config.Config
	db     *database.DB
	logger *zap.Logger
}

// NewHashedIPHandler returns a new HashedIPHandler instance.
func NewHashedIPHandler(cfg *config.Config, db *database.DB, logger *zap.Logger) *HashedIPHandler {
	return &HashedIPHandler{
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}

// ServeHTTP serves the /ip/hashed endpoint.
func (h *HashedIPHandler) Handle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ip, err := ClientIP(r, h.cfg.Proxy)
	if err != nil {
		h.logger.Error("Failed to get client IP address", zap.Error(err))

		errors.JSON(w, h.logger, errors.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get IP address. Please try again later.",
		})

		return
	}

	var (
		hashedIP    = HashIP(ip)
		contentType = r.Header.Get(xhttp.ContentType)
	)

	if contentType == xhttp.ApplicationJSON {
		w.Header().Set(xhttp.ContentType, xhttp.ApplicationJSON)

		ipModel := model.NewIP(hashedIP)

		ipJSON, err := json.Marshal(ipModel) //nolint:errchkjson // if we don't check here, another linter complains
		if err != nil {
			h.logger.Error("Failed to marshal IP address to JSON", zap.Error(err))

			errors.JSON(w, h.logger, errors.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to marshal IP address to JSON. Please try again later.",
			})

			return
		}

		_, err = w.Write(ipJSON)
		if err != nil {
			h.logger.Error("Failed to write IP address JSON to response", zap.Error(err))

			errors.JSON(w, h.logger, errors.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to write IP address JSON to response. Please try again later.",
			})

			return
		}
	} else {
		w.Header().Set(xhttp.ContentType, xhttp.TextPlain)

		_, err := w.Write([]byte(hashedIP))
		if err != nil {
			h.logger.Error("Failed to write IP address to response", zap.Error(err))

			errors.JSON(w, h.logger, errors.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Failed to write IP address to response. Please try again later.",
			})

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
