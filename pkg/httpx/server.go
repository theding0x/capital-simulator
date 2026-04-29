// Package httpx contains shared HTTP server scaffolding used by every Go
// microservice in capital-simulator: a Server with sensible timeouts, graceful
// shutdown, and built-in /healthz and /readyz endpoints.
package httpx

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// Config controls the Server lifecycle.
type Config struct {
	Addr            string        // ":8080" if empty
	ReadTimeout     time.Duration // 10s if zero
	WriteTimeout    time.Duration // 15s if zero
	IdleTimeout     time.Duration // 60s if zero
	ShutdownTimeout time.Duration // 15s if zero
}

// Server wraps *http.Server with health/ready state and graceful shutdown.
type Server struct {
	cfg    Config
	mux    *http.ServeMux
	srv    *http.Server
	logger *slog.Logger
	ready  atomic.Bool
}

// New constructs a Server. The returned Server already has /healthz and
// /readyz registered; call Handle to add additional routes.
func New(cfg Config, logger *slog.Logger) *Server {
	if cfg.Addr == "" {
		cfg.Addr = ":8080"
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 10 * time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 15 * time.Second
	}
	if cfg.IdleTimeout == 0 {
		cfg.IdleTimeout = 60 * time.Second
	}
	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 15 * time.Second
	}

	mux := http.NewServeMux()
	s := &Server{
		cfg:    cfg,
		mux:    mux,
		logger: logger,
	}

	mux.HandleFunc("/healthz", s.healthz)
	mux.HandleFunc("/readyz", s.readyz)

	s.srv = &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}
	return s
}

// Handle registers a route on the server's mux.
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.mux.Handle(pattern, handler)
}

// HandleFunc registers a route handler on the server's mux.
func (s *Server) HandleFunc(pattern string, handler http.HandlerFunc) {
	s.mux.HandleFunc(pattern, handler)
}

// MarkReady flips the readiness flag returned by /readyz.
func (s *Server) MarkReady(ready bool) {
	s.ready.Store(ready)
}

// Run starts the HTTP server and blocks until the process receives SIGINT or
// SIGTERM, at which point it performs a graceful shutdown.
func (s *Server) Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("http server listening", slog.String("addr", s.cfg.Addr))
		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case <-ctx.Done():
		s.logger.Info("shutdown signal received")
	case err := <-errCh:
		if err != nil {
			return err
		}
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
	defer cancel()

	s.MarkReady(false)
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	s.logger.Info("http server stopped")
	return nil
}

func (s *Server) healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) readyz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.ready.Load() {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ready"}`))
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
	_, _ = w.Write([]byte(`{"status":"not_ready"}`))
}
