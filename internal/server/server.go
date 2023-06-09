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
	apierror "git.sr.ht/~jamesponddotco/accio127/internal/errors"
	"git.sr.ht/~jamesponddotco/accio127/internal/server/handler"
	"git.sr.ht/~jamesponddotco/accio127/internal/server/middleware"
	"git.sr.ht/~jamesponddotco/xstd-go/xcrypto/xtls"
	"github.com/julienschmidt/httprouter"
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

	var tlsConfig *tls.Config

	if cfg.MinTLSVersion == "TLS13" {
		tlsConfig = xtls.ModernServerConfig()
	}

	if cfg.MinTLSVersion == "TLS12" {
		tlsConfig = xtls.IntermediateServerConfig()
	}

	tlsConfig.Certificates = []tls.Certificate{cert}

	middlewares := []func(httprouter.Handle) httprouter.Handle{
		func(h httprouter.Handle) httprouter.Handle { return middleware.PanicRecovery(logger, h) },
		func(h httprouter.Handle) httprouter.Handle { return middleware.UserAgent(logger, h) },
		func(h httprouter.Handle) httprouter.Handle { return middleware.AcceptRequests(logger, h) },
		func(h httprouter.Handle) httprouter.Handle { return middleware.PrivacyPolicy(cfg.PrivacyPolicy, h) },
		middleware.SecureHeader,
		middleware.CORS,
	}

	var (
		ipHandler           = handler.NewIPHandler(cfg, db, logger)
		anonymizedIPHandler = handler.NewAnonymizedIPHandler(cfg, db, logger)
		hashedIPHandler     = handler.NewHashedIPHandler(cfg, db, logger)
		metricsHandler      = handler.NewMetricsHandler(db, logger)
		healthHandler       = handler.NewHealthHandler(db, logger)
		heartbeatHandler    = handler.NewHeartbeatHandler(logger)
	)

	mux := httprouter.New()
	mux.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apierror.JSON(w, logger, apierror.ErrorResponse{
			Code:    http.StatusNotFound,
			Message: "Page not found. Check the URL and try again.",
		})
	})

	mux.GET(endpoint.IP, middleware.Chain(ipHandler.Handle, middlewares...))
	mux.GET(endpoint.IPAnonymize, middleware.Chain(anonymizedIPHandler.Handle, middlewares...))
	mux.GET(endpoint.IPHashed, middleware.Chain(hashedIPHandler.Handle, middlewares...))
	mux.GET(endpoint.Metrics, middleware.Chain(metricsHandler.Handle, middlewares...))
	mux.GET(endpoint.Health, middleware.Chain(healthHandler.Handle, middlewares...))
	mux.GET(endpoint.Ping, middleware.Chain(heartbeatHandler.Handle, middlewares...))

	httpServer := &http.Server{
		Addr:         cfg.Address,
		Handler:      mux,
		TLSConfig:    tlsConfig,
		ReadTimeout:  time.Duration(cfg.ReadTimeout),
		WriteTimeout: time.Duration(cfg.WriteTimeout),
		IdleTimeout:  time.Duration(cfg.IdleTimeout),
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
