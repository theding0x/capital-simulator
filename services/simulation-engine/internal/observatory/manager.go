package observatory

import (
	"context"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// session is one tracked run plus the time it was last accessed (for eviction).
type session struct {
	run        *Run
	lastAccess time.Time
}

// Manager owns the immutable seed template and the live, in-memory sessions. It
// is safe for concurrent use. The seed is read once at construction; runs never
// write back to any store.
type Manager struct {
	mu       sync.Mutex
	sessions map[string]*session

	seedAbode simulation.AbodeState
	seedField []circulation.FieldCapital

	ttl           time.Duration
	maxSessions   int
	sweepInterval time.Duration
	now           func() time.Time
	logger        *slog.Logger
}

// NewManager builds a Manager from the seed abode and seed field. The seed field
// is copied and sorted by ID so every run starts from a deterministic order.
func NewManager(seedAbode simulation.AbodeState, seedField []circulation.FieldCapital, logger *slog.Logger) *Manager {
	if logger == nil {
		logger = slog.Default()
	}
	field := make([]circulation.FieldCapital, len(seedField))
	copy(field, seedField)
	sort.Slice(field, func(i, j int) bool { return field[i].ID < field[j].ID })
	return &Manager{
		sessions:      make(map[string]*session),
		seedAbode:     seedAbode,
		seedField:     field,
		ttl:           15 * time.Minute,
		maxSessions:   500,
		sweepInterval: time.Minute,
		now:           time.Now,
		logger:        logger,
	}
}

// GetOrCreate returns the Run for sessionID, creating it from the seed template
// if absent. An empty sessionID returns a fresh, unstored (transient) run, so
// header-less callers get a clean seed snapshot without populating the map.
func (m *Manager) GetOrCreate(sessionID string) *Run {
	if sessionID == "" {
		return m.newRun()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[sessionID]; ok {
		s.lastAccess = m.now()
		return s.run
	}
	if len(m.sessions) >= m.maxSessions {
		m.evictOldestLocked()
	}
	run := m.newRun()
	m.sessions[sessionID] = &session{run: run, lastAccess: m.now()}
	return run
}

// Len reports the number of live sessions.
func (m *Manager) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.sessions)
}

// StartSweeper launches a background goroutine that evicts idle sessions every
// sweepInterval until ctx is cancelled.
func (m *Manager) StartSweeper(ctx context.Context) {
	go func() {
		t := time.NewTicker(m.sweepInterval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				m.sweep()
			}
		}
	}()
}

// newRun deep-copies the seed template into a fresh Run (AbodeState is all-scalar
// so it copies by value; the field slice is copied; the series starts empty).
func (m *Manager) newRun() *Run {
	field := make([]circulation.FieldCapital, len(m.seedField))
	copy(field, m.seedField)
	return &Run{abode: m.seedAbode, field: field}
}

// sweep removes every session idle longer than ttl.
func (m *Manager) sweep() {
	m.mu.Lock()
	defer m.mu.Unlock()
	cutoff := m.now().Add(-m.ttl)
	for id, s := range m.sessions {
		if s.lastAccess.Before(cutoff) {
			delete(m.sessions, id)
		}
	}
}

// evictOldestLocked removes the least-recently-accessed session. Caller holds mu.
func (m *Manager) evictOldestLocked() {
	var oldestID string
	var oldest time.Time
	first := true
	for id, s := range m.sessions {
		if first || s.lastAccess.Before(oldest) {
			oldestID, oldest, first = id, s.lastAccess, false
		}
	}
	if !first {
		delete(m.sessions, oldestID)
		m.logger.Warn("observatory: session cap reached, evicted oldest session",
			"max_sessions", m.maxSessions)
	}
}
