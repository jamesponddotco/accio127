// Package server is the main server for the application.
package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"git.sr.ht/~jamesponddotco/accio127/internal/config"
	"git.sr.ht/~jamesponddotco/accio127/internal/database"
	"git.sr.ht/~jamesponddotco/accio127/internal/endpoint"
	"git.sr.ht/~jamesponddotco/accio127/internal/server/handler"
	"git.sr.ht/~jamesponddotco/accio127/internal/server/middleware"
	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xtls"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func New(cfg *config.Config, db *database.DB, logger *zap.Logger) (*Server, error) {
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.CertKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}

	tlsConfig := xtls.DefaultServerConfig()
	tlsConfig.Certificates = []tls.Certificate{cert}

	middlewares := []func(http.Handler) http.Handler{
		func(h http.Handler) http.Handler { return middleware.PanicRecovery(logger, h) },
		middleware.StrictSlash,
		middleware.UserAgent,
		middleware.AcceptRequests,
		func(h http.Handler) http.Handler { return middleware.PrivacyPolicy(cfg.PrivacyPolicy, h) },
		middleware.SecureHeader,
		middleware.CORS,
	}

	var (
		ipHandler           = handler.NewIPHandler(db, logger)
		anonymizedIPHandler = handler.NewAnonymizedIPHandler(db, logger)
		hashedIPHandler     = handler.NewHashedIPHandler(db, logger)
		metricsHandler      = handler.NewMetricsHandler(db, logger)
		statusHandler       = handler.NewStatusHandler(db, logger)
	)

	mux := http.NewServeMux()
	mux.Handle(endpoint.IP, middleware.Chain(ipHandler, middlewares...))
	mux.Handle(endpoint.IPAnonymize, middleware.Chain(anonymizedIPHandler, middlewares...))
	mux.Handle(endpoint.IPHashed, middleware.Chain(hashedIPHandler, middlewares...))
	mux.Handle(endpoint.Metrics, middleware.Chain(metricsHandler, middlewares...))
	mux.Handle(endpoint.Status, middleware.Chain(statusHandler, middlewares...))

	httpServer := &http.Server{
		Addr:         cfg.Address,
		Handler:      mux,
		TLSConfig:    tlsConfig,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		httpServer: httpServer,
		logger:     logger,
	}, nil
}

func (s *Server) Start() error {
	var (
		sigint            = make(chan os.Signal, 1)
		shutdownCompleted = make(chan struct{})
	)

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigint

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.Error("HTTP server Shutdown:", zap.Error(err))
		}

		close(shutdownCompleted)
	}()

	if err := s.httpServer.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	<-shutdownCompleted

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}

func (s *Server) Addr() string {
	return s.httpServer.Addr
}
