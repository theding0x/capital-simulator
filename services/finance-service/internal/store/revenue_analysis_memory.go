package store

import (
	"context"
	"sort"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// CreateAnnualValueComposition stores c, assigning an ID and timestamp when absent (Vol. III Ch. 49).
func (m *Memory) CreateAnnualValueComposition(_ context.Context, c revenue.AnnualValueComposition) (revenue.AnnualValueComposition, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c.ID == "" {
		c.ID = revenue.NewAnnualValueCompositionID()
	}
	if _, exists := m.annualValueCompositions[c.ID]; exists {
		return revenue.AnnualValueComposition{}, ErrAlreadyExists
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
	m.annualValueCompositions[c.ID] = c
	return c, nil
}

// GetAnnualValueComposition returns the composition with the given ID, or ErrNotFound.
func (m *Memory) GetAnnualValueComposition(_ context.Context, id revenue.AnnualValueCompositionID) (revenue.AnnualValueComposition, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, ok := m.annualValueCompositions[id]
	if !ok {
		return revenue.AnnualValueComposition{}, ErrNotFound
	}
	return c, nil
}

// CreateRevenueLimit stores l, assigning an ID and timestamp when absent (Vol. III Ch. 49).
func (m *Memory) CreateRevenueLimit(_ context.Context, l revenue.RevenueLimit) (revenue.RevenueLimit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if l.ID == "" {
		l.ID = revenue.NewRevenueLimitID()
	}
	if _, exists := m.revenueLimits[l.ID]; exists {
		return revenue.RevenueLimit{}, ErrAlreadyExists
	}
	if l.CreatedAt.IsZero() {
		l.CreatedAt = time.Now().UTC()
	}
	m.revenueLimits[l.ID] = l
	return l, nil
}

// GetRevenueLimitByComposition returns the revenue limit for the given composition ID, or ErrNotFound.
func (m *Memory) GetRevenueLimitByComposition(_ context.Context, compositionID string) (revenue.RevenueLimit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, l := range m.revenueLimits {
		if l.CompositionID == compositionID {
			return l, nil
		}
	}
	return revenue.RevenueLimit{}, ErrNotFound
}

// CreateSurplusValuePartition stores p, assigning an ID and timestamp when absent (Vol. III Ch. 49).
func (m *Memory) CreateSurplusValuePartition(_ context.Context, p revenue.SurplusValuePartition) (revenue.SurplusValuePartition, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if p.ID == "" {
		p.ID = revenue.NewSurplusValuePartitionID()
	}
	if _, exists := m.surplusValuePartitions[p.ID]; exists {
		return revenue.SurplusValuePartition{}, ErrAlreadyExists
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now().UTC()
	}
	m.surplusValuePartitions[p.ID] = p
	return p, nil
}

// ListSurplusValuePartitions returns all stored partitions, newest first. Never returns nil.
func (m *Memory) ListSurplusValuePartitions(_ context.Context) ([]revenue.SurplusValuePartition, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.SurplusValuePartition, 0, len(m.surplusValuePartitions))
	for _, p := range m.surplusValuePartitions {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return string(out[i].ID) < string(out[j].ID)
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}
