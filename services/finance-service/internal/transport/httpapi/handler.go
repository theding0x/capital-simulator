// Package httpapi exposes finance-service over HTTP. The Handler is the
// composition root for the service's chapter handlers; each Vol. III
// chapter PR adds methods on *Handler and registers them in routes.go.
package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

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

// --- helpers ----------------------------------------------------------------

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return errors.New("invalid json: " + err.Error())
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeStoreError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, store.ErrAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}
