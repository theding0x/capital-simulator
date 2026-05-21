package httpapi

import "github.com/theding0x/capital-simulator/pkg/httpx"

// Register wires finance-service routes onto s.
//
// The function is intentionally empty in foundation Phase 3 — pkg/httpx
// already exposes /healthz and /readyz, which is all that's needed for
// the service to be probed by k8s and proxied by api-gateway. Each Vol. III
// chapter PR appends `s.HandleFunc("METHOD /v1/<resource>", h.Whatever)`
// lines here, following the pattern in agent-service's routes.go.
func Register(_ *httpx.Server, _ *Handler) {}
