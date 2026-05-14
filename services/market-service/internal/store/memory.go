package store

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

// Memory is an in-memory Store for unit tests and local development.
type Memory struct {
	mu                  sync.RWMutex
	owners              map[market.OwnerID]market.Owner
	offers              map[market.OfferID]market.Offer
	exchanges           map[market.ExchangeID]market.Exchange
	universalEquivalent *market.UniversalEquivalent
	moneyCommodity      *market.MoneyCommodity
	prices              map[market.CommodityID]market.Price
	now                 func() time.Time
}

// NewMemory returns an empty in-memory store.
func NewMemory() *Memory {
	return &Memory{
		owners:    make(map[market.OwnerID]market.Owner),
		offers:    make(map[market.OfferID]market.Offer),
		exchanges: make(map[market.ExchangeID]market.Exchange),
		prices:    make(map[market.CommodityID]market.Price),
		now:       time.Now,
	}
}

// --- Owner ------------------------------------------------------------------

func (m *Memory) CreateOwner(_ context.Context, o market.Owner) (market.Owner, error) {
	if err := o.Validate(); err != nil {
		return market.Owner{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if o.ID.IsZero() {
		o.ID = market.NewOwnerID()
	}
	now := m.now()
	o.CreatedAt = now
	o.UpdatedAt = now
	m.owners[o.ID] = o
	return o, nil
}

func (m *Memory) GetOwner(_ context.Context, id market.OwnerID) (market.Owner, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	o, ok := m.owners[id]
	if !ok {
		return market.Owner{}, ErrNotFound
	}
	return o, nil
}

func (m *Memory) ListOwners(_ context.Context) ([]market.Owner, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]market.Owner, 0, len(m.owners))
	for _, o := range m.owners {
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

// --- Offer ------------------------------------------------------------------

func (m *Memory) CreateOffer(_ context.Context, o market.Offer) (market.Offer, error) {
	if err := o.Validate(); err != nil {
		return market.Offer{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if o.ID.IsZero() {
		o.ID = market.NewOfferID()
	}
	o.CreatedAt = m.now()
	m.offers[o.ID] = o
	return o, nil
}

func (m *Memory) GetOffer(_ context.Context, id market.OfferID) (market.Offer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	o, ok := m.offers[id]
	if !ok {
		return market.Offer{}, ErrNotFound
	}
	return o, nil
}

func (m *Memory) ListOffers(_ context.Context) ([]market.Offer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]market.Offer, 0, len(m.offers))
	for _, o := range m.offers {
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

func (m *Memory) DeleteOffer(_ context.Context, id market.OfferID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.offers[id]; !ok {
		return ErrNotFound
	}
	delete(m.offers, id)
	return nil
}

// --- Exchange ---------------------------------------------------------------

func (m *Memory) CreateExchange(_ context.Context, e market.Exchange) (market.Exchange, error) {
	if err := e.Validate(); err != nil {
		return market.Exchange{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if e.ID.IsZero() {
		e.ID = market.NewExchangeID()
	}
	e.CreatedAt = m.now()
	m.exchanges[e.ID] = e
	return e, nil
}

func (m *Memory) GetExchange(_ context.Context, id market.ExchangeID) (market.Exchange, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.exchanges[id]
	if !ok {
		return market.Exchange{}, ErrNotFound
	}
	return e, nil
}

func (m *Memory) ListExchanges(_ context.Context) ([]market.Exchange, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]market.Exchange, 0, len(m.exchanges))
	for _, e := range m.exchanges {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

// --- UniversalEquivalent (singleton) ----------------------------------------

func (m *Memory) SetUniversalEquivalent(_ context.Context, ue market.UniversalEquivalent) (market.UniversalEquivalent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Idempotent: if already set to the same commodity, return as-is.
	if m.universalEquivalent != nil && m.universalEquivalent.CommodityID == ue.CommodityID {
		return *m.universalEquivalent, nil
	}
	if ue.SetAt.IsZero() {
		ue.SetAt = m.now()
	}
	m.universalEquivalent = &ue
	return ue, nil
}

func (m *Memory) GetUniversalEquivalent(_ context.Context) (market.UniversalEquivalent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.universalEquivalent == nil {
		return market.UniversalEquivalent{}, ErrNotFound
	}
	return *m.universalEquivalent, nil
}

// --- MoneyCommodity (singleton) ---------------------------------------------

func (m *Memory) SetMoneyCommodity(_ context.Context, mc market.MoneyCommodity) (market.MoneyCommodity, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if mc.CreatedAt.IsZero() {
		mc.CreatedAt = m.now()
	}
	m.moneyCommodity = &mc
	return mc, nil
}

func (m *Memory) GetMoneyCommodity(_ context.Context) (market.MoneyCommodity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.moneyCommodity == nil {
		return market.MoneyCommodity{}, ErrNotFound
	}
	return *m.moneyCommodity, nil
}

// --- Price ------------------------------------------------------------------

func (m *Memory) SetPrice(_ context.Context, p market.Price) (market.Price, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = m.now()
	}
	m.prices[p.CommodityID] = p
	return p, nil
}

func (m *Memory) GetPrice(_ context.Context, commodityID market.CommodityID) (market.Price, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.prices[commodityID]
	if !ok {
		return market.Price{}, ErrNotFound
	}
	return p, nil
}

func (m *Memory) ListPrices(_ context.Context) ([]market.Price, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]market.Price, 0, len(m.prices))
	for _, p := range m.prices {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		return string(out[i].CommodityID) < string(out[j].CommodityID)
	})
	return out, nil
}
