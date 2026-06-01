package store

import (
	"context"
	"sort"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// CreateSocialClass stores c, assigning an ID and timestamp when absent (Vol. III Ch. 52).
func (m *Memory) CreateSocialClass(_ context.Context, c revenue.SocialClass) (revenue.SocialClass, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if c.ID == "" {
		c.ID = revenue.NewSocialClassID()
	}
	if _, exists := m.socialClasses[c.ID]; exists {
		return revenue.SocialClass{}, ErrAlreadyExists
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
	m.socialClasses[c.ID] = c
	return c, nil
}

// GetSocialClass returns the class with the given ID, or ErrNotFound.
func (m *Memory) GetSocialClass(_ context.Context, id revenue.SocialClassID) (revenue.SocialClass, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	c, ok := m.socialClasses[id]
	if !ok {
		return revenue.SocialClass{}, ErrNotFound
	}
	return c, nil
}

// ListSocialClasses returns all stored classes, newest first. Never returns nil.
func (m *Memory) ListSocialClasses(_ context.Context) ([]revenue.SocialClass, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.SocialClass, 0, len(m.socialClasses))
	for _, x := range m.socialClasses {
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

// CreateClassIncomeSource stores s, assigning an ID and timestamp when absent (Vol. III Ch. 52).
func (m *Memory) CreateClassIncomeSource(_ context.Context, s revenue.ClassIncomeSource) (revenue.ClassIncomeSource, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if s.ID == "" {
		s.ID = revenue.NewClassIncomeSourceID()
	}
	if _, exists := m.classIncomeSources[s.ID]; exists {
		return revenue.ClassIncomeSource{}, ErrAlreadyExists
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now().UTC()
	}
	m.classIncomeSources[s.ID] = s
	return s, nil
}

// ListClassIncomeSources returns all stored income sources, newest first. Never returns nil.
func (m *Memory) ListClassIncomeSources(_ context.Context) ([]revenue.ClassIncomeSource, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.ClassIncomeSource, 0, len(m.classIncomeSources))
	for _, x := range m.classIncomeSources {
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

// CreateClassTendency stores t, assigning an ID and timestamp when absent (Vol. III Ch. 52).
func (m *Memory) CreateClassTendency(_ context.Context, t revenue.ClassTendency) (revenue.ClassTendency, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if t.ID == "" {
		t.ID = revenue.NewClassTendencyID()
	}
	if _, exists := m.classTendencies[t.ID]; exists {
		return revenue.ClassTendency{}, ErrAlreadyExists
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now().UTC()
	}
	m.classTendencies[t.ID] = t
	return t, nil
}

// ListClassTendencies returns all stored tendencies, newest first. Never returns nil.
func (m *Memory) ListClassTendencies(_ context.Context) ([]revenue.ClassTendency, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.ClassTendency, 0, len(m.classTendencies))
	for _, x := range m.classTendencies {
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

// CreateClassAmbiguity stores a, assigning an ID and timestamp when absent (Vol. III Ch. 52).
func (m *Memory) CreateClassAmbiguity(_ context.Context, a revenue.ClassAmbiguity) (revenue.ClassAmbiguity, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if a.ID == "" {
		a.ID = revenue.NewClassAmbiguityID()
	}
	if _, exists := m.classAmbiguities[a.ID]; exists {
		return revenue.ClassAmbiguity{}, ErrAlreadyExists
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now().UTC()
	}
	m.classAmbiguities[a.ID] = a
	return a, nil
}

// ListClassAmbiguities returns all stored ambiguities, newest first. Never returns nil.
func (m *Memory) ListClassAmbiguities(_ context.Context) ([]revenue.ClassAmbiguity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := make([]revenue.ClassAmbiguity, 0, len(m.classAmbiguities))
	for _, x := range m.classAmbiguities {
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
