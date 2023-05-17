package handler

import (
	"encoding/json"
	"net"
	"net/http"

	"git.sr.ht/~jamesponddotco/accio127/internal/database"
	"git.sr.ht/~jamesponddotco/accio127/internal/server/model"
	"go.uber.org/zap"
)

// AnonymizedIPHandler is an HTTP handler for the /ip/anonymized endpoint.
type AnonymizedIPHandler struct {
	db     *database.DB
	logger *zap.Logger
}

// NewAnonymizedIPHandler creates a new AnonymizedIPHandler instance.
func NewAnonymizedIPHandler(db *database.DB, logger *zap.Logger) *AnonymizedIPHandler {
	return &AnonymizedIPHandler{
		db:     db,
		logger: logger,
	}
}

// ServeHTTP serves the /ip/anonymized endpoint.
func (h *AnonymizedIPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip, err := ClientIP(r)
	if err != nil {
		h.logger.Error("Failed to get client IP address", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	anonymizedIP := AnonymizeIP(ip)
	if anonymizedIP == "" {
		h.logger.Error("Failed to anonymize client IP address", zap.Error(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		w.Header().Set("Content-Type", "application/json")

		ipModel := model.NewIP(anonymizedIP)

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

		_, err := w.Write([]byte(anonymizedIP))
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
