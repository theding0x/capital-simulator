package store

import (
	"sync"
	"time"
)

// Memory is an in-memory Store for unit tests and local development.
// The fields below are scaffolding for Vol. III chapter PRs to populate:
// add `something map[somethingID]Something` entries here and matching
// CRUD methods on *Memory.
type Memory struct {
	mu  sync.RWMutex
	now func() time.Time
}

// NewMemory returns an empty in-memory store.
func NewMemory() *Memory {
	return &Memory{now: time.Now}
}
