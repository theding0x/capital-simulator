package store

import (
	"context"
	"sort"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// CreateProductionRelationBasis stores b, assigning an ID and timestamp when absent (Vol. III Ch. 51).
func (m *Memory) CreateProductionRelationBasis(_ context.Context, b revenue.ProductionRelationBasis) (revenue.ProductionRelationBasis, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if b.ID == "" {
		b.ID = revenue.NewProductionRelationBasisID()
	}
	if _, exists := m.productionRelationBases[b.ID]; exists {
		return revenue.ProductionRelationBasis{}, ErrAlreadyExists
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = time.Now().UTC()
	}
	m.productionRelationBases[b.ID] = b
	return b, nil
}

// ListProductionRelationBases returns all stored bases, newest first. Never returns nil.
func (m *Memory) ListProductionRelationBases(_ context.Context) ([]revenue.ProductionRelationBasis, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.ProductionRelationBasis, 0, len(m.productionRelationBases))
	for _, x := range m.productionRelationBases {
		out = append(out, x)
	}
	sort.Slice(out, func(a, b int) bool {
		if out[a].CreatedAt.Equal(out[b].CreatedAt) {
			return string(out[a].ID) < string(out[b].ID)
		}
		return out[a].CreatedAt.After(out[b].CreatedAt)
	})
	return out, nil
}

// CreateDistributionRelation stores d, assigning an ID and timestamp when absent (Vol. III Ch. 51).
func (m *Memory) CreateDistributionRelation(_ context.Context, d revenue.DistributionRelation) (revenue.DistributionRelation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if d.ID == "" {
		d.ID = revenue.NewDistributionRelationID()
	}
	if _, exists := m.distributionRelations[d.ID]; exists {
		return revenue.DistributionRelation{}, ErrAlreadyExists
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = time.Now().UTC()
	}
	m.distributionRelations[d.ID] = d
	return d, nil
}

// ListDistributionRelations returns all stored relations, newest first. Never returns nil.
func (m *Memory) ListDistributionRelations(_ context.Context) ([]revenue.DistributionRelation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.DistributionRelation, 0, len(m.distributionRelations))
	for _, x := range m.distributionRelations {
		out = append(out, x)
	}
	sort.Slice(out, func(a, b int) bool {
		if out[a].CreatedAt.Equal(out[b].CreatedAt) {
			return string(out[a].ID) < string(out[b].ID)
		}
		return out[a].CreatedAt.After(out[b].CreatedAt)
	})
	return out, nil
}

// CreateHistoricalTransience stores t, assigning an ID and timestamp when absent (Vol. III Ch. 51).
func (m *Memory) CreateHistoricalTransience(_ context.Context, t revenue.HistoricalTransience) (revenue.HistoricalTransience, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if t.ID == "" {
		t.ID = revenue.NewHistoricalTransienceID()
	}
	if _, exists := m.historicalTransiences[t.ID]; exists {
		return revenue.HistoricalTransience{}, ErrAlreadyExists
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}
	m.historicalTransiences[t.ID] = t
	return t, nil
}

// ListHistoricalTransiences returns all stored records, newest first. Never returns nil.
func (m *Memory) ListHistoricalTransiences(_ context.Context) ([]revenue.HistoricalTransience, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.HistoricalTransience, 0, len(m.historicalTransiences))
	for _, x := range m.historicalTransiences {
		out = append(out, x)
	}
	sort.Slice(out, func(a, b int) bool {
		if out[a].CreatedAt.Equal(out[b].CreatedAt) {
			return string(out[a].ID) < string(out[b].ID)
		}
		return out[a].CreatedAt.After(out[b].CreatedAt)
	})
	return out, nil
}

// CreateProductiveForceContradiction stores c, assigning an ID and timestamp when absent (Vol. III Ch. 51).
func (m *Memory) CreateProductiveForceContradiction(_ context.Context, c revenue.ProductiveForceContradiction) (revenue.ProductiveForceContradiction, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c.ID == "" {
		c.ID = revenue.NewProductiveForceContradictionID()
	}
	if _, exists := m.productiveForceContradictions[c.ID]; exists {
		return revenue.ProductiveForceContradiction{}, ErrAlreadyExists
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
	m.productiveForceContradictions[c.ID] = c
	return c, nil
}

// ListProductiveForceContradictions returns all stored records, newest first. Never returns nil.
func (m *Memory) ListProductiveForceContradictions(_ context.Context) ([]revenue.ProductiveForceContradiction, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.ProductiveForceContradiction, 0, len(m.productiveForceContradictions))
	for _, x := range m.productiveForceContradictions {
		out = append(out, x)
	}
	sort.Slice(out, func(a, b int) bool {
		if out[a].CreatedAt.Equal(out[b].CreatedAt) {
			return string(out[a].ID) < string(out[b].ID)
		}
		return out[a].CreatedAt.After(out[b].CreatedAt)
	})
	return out, nil
}
