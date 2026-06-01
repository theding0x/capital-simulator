package store

import (
	"context"
	"sort"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// CreateCompetitionIllusion stores i, assigning an ID and timestamp when absent (Vol. III Ch. 50).
func (m *Memory) CreateCompetitionIllusion(_ context.Context, i revenue.CompetitionIllusion) (revenue.CompetitionIllusion, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if i.ID == "" {
		i.ID = revenue.NewCompetitionIllusionID()
	}
	if _, exists := m.competitionIllusions[i.ID]; exists {
		return revenue.CompetitionIllusion{}, ErrAlreadyExists
	}
	if i.CreatedAt.IsZero() {
		i.CreatedAt = time.Now().UTC()
	}
	m.competitionIllusions[i.ID] = i
	return i, nil
}

// ListCompetitionIllusions returns all stored illusions, newest first. Never returns nil.
func (m *Memory) ListCompetitionIllusions(_ context.Context) ([]revenue.CompetitionIllusion, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.CompetitionIllusion, 0, len(m.competitionIllusions))
	for _, x := range m.competitionIllusions {
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

// CreateCostPriceFetish stores f, assigning an ID and timestamp when absent (Vol. III Ch. 50).
func (m *Memory) CreateCostPriceFetish(_ context.Context, f revenue.CostPriceFetish) (revenue.CostPriceFetish, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if f.ID == "" {
		f.ID = revenue.NewCostPriceFetishID()
	}
	if _, exists := m.costPriceFetishes[f.ID]; exists {
		return revenue.CostPriceFetish{}, ErrAlreadyExists
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now().UTC()
	}
	m.costPriceFetishes[f.ID] = f
	return f, nil
}

// ListCostPriceFetishes returns all stored cost-price fetishes, newest first. Never returns nil.
func (m *Memory) ListCostPriceFetishes(_ context.Context) ([]revenue.CostPriceFetish, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.CostPriceFetish, 0, len(m.costPriceFetishes))
	for _, x := range m.costPriceFetishes {
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

// CreateAnnualRevenueAccount stores a, assigning an ID and timestamp when absent (Vol. III Ch. 50).
func (m *Memory) CreateAnnualRevenueAccount(_ context.Context, a revenue.AnnualRevenueAccount) (revenue.AnnualRevenueAccount, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if a.ID == "" {
		a.ID = revenue.NewAnnualRevenueAccountID()
	}
	if _, exists := m.annualRevenueAccounts[a.ID]; exists {
		return revenue.AnnualRevenueAccount{}, ErrAlreadyExists
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now().UTC()
	}
	m.annualRevenueAccounts[a.ID] = a
	return a, nil
}

// ListAnnualRevenueAccounts returns all stored accounts, newest first. Never returns nil.
func (m *Memory) ListAnnualRevenueAccounts(_ context.Context) ([]revenue.AnnualRevenueAccount, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.AnnualRevenueAccount, 0, len(m.annualRevenueAccounts))
	for _, x := range m.annualRevenueAccounts {
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

// CreateValueAppearanceContrast stores c, assigning an ID and timestamp when absent (Vol. III Ch. 50).
func (m *Memory) CreateValueAppearanceContrast(_ context.Context, c revenue.ValueAppearanceContrast) (revenue.ValueAppearanceContrast, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c.ID == "" {
		c.ID = revenue.NewValueAppearanceContrastID()
	}
	if _, exists := m.valueAppearanceContrasts[c.ID]; exists {
		return revenue.ValueAppearanceContrast{}, ErrAlreadyExists
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
	m.valueAppearanceContrasts[c.ID] = c
	return c, nil
}

// ListValueAppearanceContrasts returns all stored contrasts, newest first. Never returns nil.
func (m *Memory) ListValueAppearanceContrasts(_ context.Context) ([]revenue.ValueAppearanceContrast, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.ValueAppearanceContrast, 0, len(m.valueAppearanceContrasts))
	for _, x := range m.valueAppearanceContrasts {
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
