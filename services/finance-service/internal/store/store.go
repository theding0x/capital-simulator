// Package store defines the persistence boundary for finance-service.
// The Store interface is the only seam between the HTTP/domain layer and
// the underlying database (MySQL in production, in-memory in tests).
//
// finance-service is the home for Capital Vol. III — profit, the average
// rate of profit, commercial capital, interest-bearing capital, credit,
// rent, and the revenue distribution. Each Vol. III chapter PR adds the
// methods it needs to this interface and provides the matching Memory +
// MySQL implementations.
package store

import "errors"

// Sentinel errors that callers can branch on with errors.Is.
var (
	ErrNotFound      = errors.New("finance: not found")
	ErrAlreadyExists = errors.New("finance: already exists")
)

// Store is the persistence contract for finance-service.
//
// The interface is intentionally empty in foundation Phase 3 — finance-service
// is scaffolded ahead of any Vol. III chapter so the api-gateway can route to
// it and `make build` produces a binary. Each Vol. III chapter PR adds methods
// here and corresponding implementations in memory.go + mysql.go, following
// the per-chapter pattern documented in CLAUDE.md.
type Store interface{}
