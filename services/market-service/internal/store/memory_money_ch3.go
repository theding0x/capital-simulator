package store

import (
	"context"
	"sort"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

// --- Hoard ------------------------------------------------------------------

func (m *Memory) CreateHoard(_ context.Context, h market.Hoard) (market.Hoard, error) {
	if err := h.Validate(); err != nil {
		return market.Hoard{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if h.ID.IsZero() {
		h.ID = market.NewHoardID()
	}
	if _, exists := m.hoards[h.ID]; exists {
		return market.Hoard{}, ErrAlreadyExists
	}
	if h.CreatedAt.IsZero() {
		h.CreatedAt = m.now().UTC()
	}
	m.hoards[h.ID] = h
	return h, nil
}

func (m *Memory) GetHoard(_ context.Context, id market.HoardID) (market.Hoard, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	h, ok := m.hoards[id]
	if !ok {
		return market.Hoard{}, ErrNotFound
	}
	return h, nil
}

func (m *Memory) ListHoards(_ context.Context) ([]market.Hoard, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]market.Hoard, 0, len(m.hoards))
	for _, h := range m.hoards {
		out = append(out, h)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

// --- PaymentObligation -----------------------------------------------------

func (m *Memory) CreatePaymentObligation(_ context.Context, p market.PaymentObligation) (market.PaymentObligation, error) {
	if err := p.Validate(); err != nil {
		return market.PaymentObligation{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if p.ID.IsZero() {
		p.ID = market.NewPaymentObligationID()
	}
	if _, exists := m.paymentObligations[p.ID]; exists {
		return market.PaymentObligation{}, ErrAlreadyExists
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = m.now().UTC()
	}
	m.paymentObligations[p.ID] = p
	return p, nil
}

func (m *Memory) GetPaymentObligation(_ context.Context, id market.PaymentObligationID) (market.PaymentObligation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.paymentObligations[id]
	if !ok {
		return market.PaymentObligation{}, ErrNotFound
	}
	return p, nil
}

func (m *Memory) ListPaymentObligations(_ context.Context) ([]market.PaymentObligation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]market.PaymentObligation, 0, len(m.paymentObligations))
	for _, p := range m.paymentObligations {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

func (m *Memory) SettlePaymentObligation(_ context.Context, id market.PaymentObligationID) (market.PaymentObligation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.paymentObligations[id]
	if !ok {
		return market.PaymentObligation{}, ErrNotFound
	}
	if p.PaidAt == nil {
		now := m.now().UTC()
		p.PaidAt = &now
		m.paymentObligations[id] = p
	}
	return p, nil
}

// --- WorldMoneyTransfer ----------------------------------------------------

func (m *Memory) CreateWorldMoneyTransfer(_ context.Context, t market.WorldMoneyTransfer) (market.WorldMoneyTransfer, error) {
	if err := t.Validate(); err != nil {
		return market.WorldMoneyTransfer{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if t.ID.IsZero() {
		t.ID = market.NewWorldMoneyTransferID()
	}
	if _, exists := m.worldMoneyTransfers[t.ID]; exists {
		return market.WorldMoneyTransfer{}, ErrAlreadyExists
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = m.now().UTC()
	}
	m.worldMoneyTransfers[t.ID] = t
	return t, nil
}

func (m *Memory) ListWorldMoneyTransfers(_ context.Context) ([]market.WorldMoneyTransfer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]market.WorldMoneyTransfer, 0, len(m.worldMoneyTransfers))
	for _, t := range m.worldMoneyTransfers {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}
