package store

import (
	"context"
	"sort"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// CreateRevenueStream stores s, assigning an ID and timestamp when absent (Vol. III Ch. 48).
func (m *Memory) CreateRevenueStream(_ context.Context, s revenue.RevenueStream) (revenue.RevenueStream, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s.ID == "" {
		s.ID = revenue.NewRevenueStreamID()
	}
	if _, exists := m.revenueStreams[s.ID]; exists {
		return revenue.RevenueStream{}, ErrAlreadyExists
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now().UTC()
	}
	m.revenueStreams[s.ID] = s
	return s, nil
}

// ListRevenueStreams returns all stored revenue streams, newest first. Never returns nil.
func (m *Memory) ListRevenueStreams(_ context.Context) ([]revenue.RevenueStream, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.RevenueStream, 0, len(m.revenueStreams))
	for _, s := range m.revenueStreams {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return string(out[i].ID) < string(out[j].ID)
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}

// CreateTrinityFormula stores tf, assigning an ID and timestamp when absent (Vol. III Ch. 48).
func (m *Memory) CreateTrinityFormula(_ context.Context, tf revenue.TrinityFormula) (revenue.TrinityFormula, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if tf.ID == "" {
		tf.ID = revenue.NewTrinityFormulaID()
	}
	if _, exists := m.trinityFormulas[tf.ID]; exists {
		return revenue.TrinityFormula{}, ErrAlreadyExists
	}
	if tf.CreatedAt.IsZero() {
		tf.CreatedAt = time.Now().UTC()
	}
	m.trinityFormulas[tf.ID] = tf
	return tf, nil
}

// GetTrinityFormula returns the trinity formula with the given ID, or ErrNotFound.
func (m *Memory) GetTrinityFormula(_ context.Context, id revenue.TrinityFormulaID) (revenue.TrinityFormula, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tf, ok := m.trinityFormulas[id]
	if !ok {
		return revenue.TrinityFormula{}, ErrNotFound
	}
	return tf, nil
}

// CreateRevenueFetishForm stores f, assigning an ID and timestamp when absent (Vol. III Ch. 48).
func (m *Memory) CreateRevenueFetishForm(_ context.Context, f revenue.RevenueFetishForm) (revenue.RevenueFetishForm, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if f.ID == "" {
		f.ID = revenue.NewRevenueFetishFormID()
	}
	if _, exists := m.revenueFetishForms[f.ID]; exists {
		return revenue.RevenueFetishForm{}, ErrAlreadyExists
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now().UTC()
	}
	m.revenueFetishForms[f.ID] = f
	return f, nil
}

// ListRevenueFetishForms returns all stored revenue fetish forms, newest first. Never returns nil.
func (m *Memory) ListRevenueFetishForms(_ context.Context) ([]revenue.RevenueFetishForm, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.RevenueFetishForm, 0, len(m.revenueFetishForms))
	for _, f := range m.revenueFetishForms {
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return string(out[i].ID) < string(out[j].ID)
		}
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}
