package handler

import (
	"encoding/json"
	"net"
	"net/http"

	"git.sr.ht/~jamesponddotco/accio127/internal/config"
	"git.sr.ht/~jamesponddotco/accio127/internal/database"
	"git.sr.ht/~jamesponddotco/accio127/internal/errors"
	"git.sr.ht/~jamesponddotco/accio127/internal/server/model"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// AnonymizedIPHandler is an HTTP handler for the /ip/anonymized endpoint.
type AnonymizedIPHandler struct {
	cfg    *config.Config
	db     *database.DB
	logger *zap.Logger
}

// NewAnonymizedIPHandler creates a new AnonymizedIPHandler instance.
func NewAnonymizedIPHandler(cfg *config.Config, db *database.DB, logger *zap.Logger) *AnonymizedIPHandler {
	return &AnonymizedIPHandler{
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}

// ServeHTTP serves the /ip/anonymized endpoint.
func (h *AnonymizedIPHandler) Handle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ip, err := ClientIP(r, h.cfg.Proxy)
	if err != nil {
		h.logger.Error("Failed to get client IP address", zap.Error(err))

		errors.JSON(w, h.logger, errors.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get IP address. Please try again later.",
		})

		return
	}

	anonymizedIP := AnonymizeIP(ip)
	if anonymizedIP == "" {
		h.logger.Error("Failed to anonymize client IP address", zap.Error(err))

		errors.JSON(w, h.logger, errors.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to anonymize IP address. Please try again later.",
		})

		return
	}

	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		w.Header().Set("Content-Type", "application/json")

		ipModel := model.NewIP(anonymizedIP)

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
		w.Header().Set("Content-Type", "text/plain")

		_, err := w.Write([]byte(anonymizedIP))
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

// AnonymizeIP anonymizes the last two octets of an IPv4 address or the last 80
// bits of an IPv6 address.
func AnonymizeIP(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return ""
	}

	if ip.To4() != nil {
		return ip.Mask(net.CIDRMask(16, 32)).String()
	}

	return ip.Mask(net.CIDRMask(48, 128)).String()
}
