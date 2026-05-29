package store

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)

// Memory is an in-memory Store. Used by unit tests and as a fallback when
// MongoDB is unavailable in local development.
type Memory struct {
	mu       sync.RWMutex
	items    map[commodity.ID]commodity.Commodity
	accounts map[commodity.ProductionAccountID]commodity.ProductionAccount
	now      func() time.Time
}

// NewMemory returns an empty Memory store.
func NewMemory() *Memory {
	return &Memory{
		items:    make(map[commodity.ID]commodity.Commodity),
		accounts: make(map[commodity.ProductionAccountID]commodity.ProductionAccount),
		now:      time.Now,
	}
}

// Create inserts c. Name is treated as a uniqueness key (case-insensitive
// after trimming) so the simulation has one canonical record per commodity.
func (m *Memory) Create(_ context.Context, c commodity.Commodity) (commodity.Commodity, error) {
	if err := c.Validate(); err != nil {
		return commodity.Commodity{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	wantName := normalizeName(c.Name)
	for _, existing := range m.items {
		if normalizeName(existing.Name) == wantName {
			return commodity.Commodity{}, ErrAlreadyExists
		}
	}

	if c.ID.IsZero() {
		c.ID = commodity.NewID()
	}
	now := m.now()
	c.CreatedAt = now
	c.UpdatedAt = now
	m.items[c.ID] = c
	return c, nil
}

func (m *Memory) Get(_ context.Context, id commodity.ID) (commodity.Commodity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.items[id]
	if !ok {
		return commodity.Commodity{}, ErrNotFound
	}
	return c, nil
}

func (m *Memory) List(_ context.Context) ([]commodity.Commodity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]commodity.Commodity, 0, len(m.items))
	for _, c := range m.items {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

func (m *Memory) Update(_ context.Context, id commodity.ID, u Update) (commodity.Commodity, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cur, ok := m.items[id]
	if !ok {
		return commodity.Commodity{}, ErrNotFound
	}
	next := u.Apply(cur)
	if err := next.Validate(); err != nil {
		return commodity.Commodity{}, err
	}
	// Name uniqueness is preserved on rename.
	if u.Name != nil && !strings.EqualFold(*u.Name, cur.Name) {
		want := normalizeName(*u.Name)
		for otherID, other := range m.items {
			if otherID == id {
				continue
			}
			if normalizeName(other.Name) == want {
				return commodity.Commodity{}, ErrAlreadyExists
			}
		}
	}
	next.UpdatedAt = m.now()
	m.items[id] = next
	return next, nil
}

func (m *Memory) Delete(_ context.Context, id commodity.ID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.items[id]; !ok {
		return ErrNotFound
	}
	delete(m.items, id)
	return nil
}

func normalizeName(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

func (m *Memory) CreateProductionAccount(_ context.Context, a commodity.ProductionAccount) (commodity.ProductionAccount, error) {
	if err := a.Validate(); err != nil {
		return commodity.ProductionAccount{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if a.ID.IsZero() {
		a.ID = commodity.NewProductionAccountID()
	}
	a.CreatedAt = m.now()
	m.accounts[a.ID] = a
	return a, nil
}

func (m *Memory) GetProductionAccount(_ context.Context, id commodity.ProductionAccountID) (commodity.ProductionAccount, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.accounts[id]
	if !ok {
		return commodity.ProductionAccount{}, ErrNotFound
	}
	return a, nil
}

func (m *Memory) ListProductionAccounts(_ context.Context) ([]commodity.ProductionAccount, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]commodity.ProductionAccount, 0, len(m.accounts))
	for _, a := range m.accounts {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}
