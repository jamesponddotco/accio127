package handler

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"git.sr.ht/~jamesponddotco/accio127/internal/config"
	"git.sr.ht/~jamesponddotco/accio127/internal/database"
	"git.sr.ht/~jamesponddotco/accio127/internal/errors"
	"git.sr.ht/~jamesponddotco/accio127/internal/server/model"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

// IPHandler is an HTTP handler for the /ip endpoint.
type IPHandler struct {
	cfg    *config.Config
	db     *database.DB
	logger *zap.Logger
}

// NewIPHandler creates a new IPHandler instance.
func NewIPHandler(cfg *config.Config, db *database.DB, logger *zap.Logger) *IPHandler {
	return &IPHandler{
		cfg:    cfg,
		db:     db,
		logger: logger,
	}
}

// ServeHTTP serves the /ip endpoint.
func (h *IPHandler) Handle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ip, err := ClientIP(r, h.cfg.Proxy)
	if err != nil {
		h.logger.Error("Failed to get client IP address", zap.Error(err))

		errors.JSON(w, h.logger, errors.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get IP address. Please try again later.",
		})

		return
	}

	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		w.Header().Set("Content-Type", "application/json")

		ipModel := model.NewIP(ip)

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

		_, err := w.Write([]byte(ip))
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

// ClientIP returns the client's IP address from the request headers or RemoteAddr.
func ClientIP(r *http.Request, proxy string) (string, error) {
	var (
		headers = []string{
			"CF-Connecting-IP",
			"True-Client-IP",
			"X-Real-IP",
			"X-Forwarded-For",
		}
		trustedProxy = proxy
	)

	// Check the remote IP. If it's not the trusted proxy, we can't trust the headers.
	// https://adam-p.ca/blog/2022/03/x-forwarded-for/
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", fmt.Errorf("failed to split remote address: %w", err)
	}

	remoteIP = strings.Trim(remoteIP, "[]")

	if remoteIP != trustedProxy {
		return remoteIP, nil
	}

	// If the request comes from a trusted proxy, check the headers.
	for _, header := range headers {
		ip := r.Header.Get(header)
		if ip != "" {
			ips := strings.Split(ip, ", ")

			// If we find multiple IPs in the header, we want to take the
			// left-most one since we're going from client to server, left to
			// right.
			return ips[0], nil
		}
	}

	return remoteIP, nil
}
