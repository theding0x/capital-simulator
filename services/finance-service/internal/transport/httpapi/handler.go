// Package httpapi exposes finance-service over HTTP. The Handler is the
// composition root for the service's chapter handlers; each Vol. III
// chapter PR adds methods on *Handler and registers them in routes.go.
package httpapi

import (
	"log/slog"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

// Handler bundles the dependencies for the HTTP layer.
type Handler struct {
	Store  store.Store
	Logger *slog.Logger
}

// New constructs a Handler. Vol. III chapter PRs may extend the signature
// to accept additional store-shaped dependencies (cf. agent-service's
// AgentStore composite once finance-service has multiple stores).
func New(s store.Store, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{Store: s, Logger: logger}
}
