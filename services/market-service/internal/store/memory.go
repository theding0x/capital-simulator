package store

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/circulation"
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

	// Vol. II Ch. 5 — circulation time.
	turnoverTimes       map[circulation.TurnoverTimeID]circulation.TurnoverTime
	sellingPhases       map[circulation.SellingPhaseID]circulation.SellingPhase
	buyingPhases        map[circulation.BuyingPhaseID]circulation.BuyingPhase
	naturalProcessSpans map[circulation.NaturalProcessSpanID]circulation.NaturalProcessSpan
	latentMPs           map[circulation.LatentProductiveCapitalID]circulation.LatentProductiveCapital
	perishabilities     map[circulation.CommodityID]circulation.Perishability
	marketSeparations   map[circulation.IndustrialCapitalID]circulation.MarketSeparation
	// openSellingPhase tracks the currently open selling phase per TurnoverTimeID.
	openSellingPhase map[circulation.TurnoverTimeID]circulation.SellingPhaseID
	// openBuyingPhase tracks the currently open buying phase per TurnoverTimeID.
	openBuyingPhase map[circulation.TurnoverTimeID]circulation.BuyingPhaseID
}

// NewMemory returns an empty in-memory store.
func NewMemory() *Memory {
	return &Memory{
		owners:              make(map[market.OwnerID]market.Owner),
		offers:              make(map[market.OfferID]market.Offer),
		exchanges:           make(map[market.ExchangeID]market.Exchange),
		prices:              make(map[market.CommodityID]market.Price),
		turnoverTimes:       make(map[circulation.TurnoverTimeID]circulation.TurnoverTime),
		sellingPhases:       make(map[circulation.SellingPhaseID]circulation.SellingPhase),
		buyingPhases:        make(map[circulation.BuyingPhaseID]circulation.BuyingPhase),
		naturalProcessSpans: make(map[circulation.NaturalProcessSpanID]circulation.NaturalProcessSpan),
		latentMPs:           make(map[circulation.LatentProductiveCapitalID]circulation.LatentProductiveCapital),
		perishabilities:     make(map[circulation.CommodityID]circulation.Perishability),
		marketSeparations:   make(map[circulation.IndustrialCapitalID]circulation.MarketSeparation),
		openSellingPhase:    make(map[circulation.TurnoverTimeID]circulation.SellingPhaseID),
		openBuyingPhase:     make(map[circulation.TurnoverTimeID]circulation.BuyingPhaseID),
		now:                 time.Now,
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

// --- Vol. II Ch. 5 — TurnoverTimeStore -------------------------------------

func (m *Memory) CreateTurnoverTime(_ context.Context, t circulation.TurnoverTime) (circulation.TurnoverTime, error) {
	if err := t.Validate(); err != nil {
		return circulation.TurnoverTime{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if t.ID.IsZero() {
		t.ID = circulation.NewTurnoverTimeID()
	}
	now := m.now()
	t.CreatedAt = now
	t.UpdatedAt = now
	m.turnoverTimes[t.ID] = t
	return t, nil
}

func (m *Memory) GetTurnoverTime(_ context.Context, id circulation.TurnoverTimeID) (circulation.TurnoverTime, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.turnoverTimes[id]
	if !ok {
		return circulation.TurnoverTime{}, ErrNotFound
	}
	return t, nil
}

func (m *Memory) ListTurnoverTimes(_ context.Context) ([]circulation.TurnoverTime, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]circulation.TurnoverTime, 0, len(m.turnoverTimes))
	for _, t := range m.turnoverTimes {
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

func (m *Memory) AddLabourTime(_ context.Context, id circulation.TurnoverTimeID, nanos int64) (circulation.TurnoverTime, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.turnoverTimes[id]
	if !ok {
		return circulation.TurnoverTime{}, ErrNotFound
	}
	if t.CirculationOpen {
		return circulation.TurnoverTime{}, circulation.ErrConcurrentProductionAndCirculation
	}
	t.Production.LabourTimeNanos += nanos
	t.UpdatedAt = m.now()
	m.turnoverTimes[id] = t
	return t, nil
}

func (m *Memory) AddLabourInterruption(_ context.Context, id circulation.TurnoverTimeID, nanos int64) (circulation.TurnoverTime, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.turnoverTimes[id]
	if !ok {
		return circulation.TurnoverTime{}, ErrNotFound
	}
	if t.CirculationOpen {
		return circulation.TurnoverTime{}, circulation.ErrConcurrentProductionAndCirculation
	}
	t.Production.LabourInterruptionNanos += nanos
	t.UpdatedAt = m.now()
	m.turnoverTimes[id] = t
	return t, nil
}

func (m *Memory) RecordLatentMP(_ context.Context, lpc circulation.LatentProductiveCapital) (circulation.LatentProductiveCapital, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.turnoverTimes[lpc.TurnoverTimeID]; !ok {
		return circulation.LatentProductiveCapital{}, ErrNotFound
	}
	if lpc.ID.IsZero() {
		lpc.ID = circulation.NewLatentProductiveCapitalID()
	}
	m.latentMPs[lpc.ID] = lpc

	// Accumulate latent nanos in the parent TurnoverTime using HeldAt as 1-tick.
	// Latent capital contributes to ProductionTime.LatentNanos; we add 1 tick
	// per registration so the TurnoverTime reflects its presence.
	t := m.turnoverTimes[lpc.TurnoverTimeID]
	t.Production.LatentNanos += int64(time.Second) // 1-second symbolic tick
	t.UpdatedAt = m.now()
	m.turnoverTimes[t.ID] = t
	return lpc, nil
}

func (m *Memory) RecordNaturalProcess(_ context.Context, nps circulation.NaturalProcessSpan) (circulation.NaturalProcessSpan, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.turnoverTimes[nps.TurnoverTimeID]; !ok {
		return circulation.NaturalProcessSpan{}, ErrNotFound
	}
	if nps.ID.IsZero() {
		nps.ID = circulation.NewNaturalProcessSpanID()
	}
	m.naturalProcessSpans[nps.ID] = nps

	t := m.turnoverTimes[nps.TurnoverTimeID]
	t.Production.NaturalProcessNanos += nps.DurationNanos
	t.UpdatedAt = m.now()
	m.turnoverTimes[t.ID] = t
	return nps, nil
}

func (m *Memory) OpenSellingPhase(_ context.Context, sp circulation.SellingPhase) (circulation.SellingPhase, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.turnoverTimes[sp.TurnoverTimeID]
	if !ok {
		return circulation.SellingPhase{}, ErrNotFound
	}
	if _, exists := m.openSellingPhase[sp.TurnoverTimeID]; exists {
		return circulation.SellingPhase{}, circulation.ErrSellingPhaseAlreadyOpen
	}
	if sp.ID.IsZero() {
		sp.ID = circulation.NewSellingPhaseID()
	}
	if sp.Outcome == "" {
		sp.Outcome = circulation.OutcomePending
	}
	m.sellingPhases[sp.ID] = sp
	m.openSellingPhase[sp.TurnoverTimeID] = sp.ID

	t.CirculationOpen = true
	t.UpdatedAt = m.now()
	m.turnoverTimes[t.ID] = t
	return sp, nil
}

func (m *Memory) CloseSellingPhase(_ context.Context, id circulation.TurnoverTimeID, spID circulation.SellingPhaseID, outcome circulation.SellingOutcome) (circulation.SellingPhase, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	openID, hasOpen := m.openSellingPhase[id]
	if !hasOpen {
		return circulation.SellingPhase{}, circulation.ErrNoOpenSellingPhase
	}
	if openID != spID {
		return circulation.SellingPhase{}, circulation.ErrNoOpenSellingPhase
	}
	sp := m.sellingPhases[spID]
	now := m.now()
	sp.ClosedAt = &now
	sp.Outcome = outcome
	m.sellingPhases[spID] = sp
	delete(m.openSellingPhase, id)

	t := m.turnoverTimes[id]
	t.Circulation.SellingTimeNanos += sp.ElapsedNanos()
	// Circulation remains open until the buying phase is also closed.
	t.UpdatedAt = now
	m.turnoverTimes[id] = t
	return sp, nil
}

func (m *Memory) OpenBuyingPhase(_ context.Context, bp circulation.BuyingPhase) (circulation.BuyingPhase, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.turnoverTimes[bp.TurnoverTimeID]; !ok {
		return circulation.BuyingPhase{}, ErrNotFound
	}
	if _, exists := m.openBuyingPhase[bp.TurnoverTimeID]; exists {
		return circulation.BuyingPhase{}, circulation.ErrBuyingPhaseAlreadyOpen
	}
	if bp.ID.IsZero() {
		bp.ID = circulation.NewBuyingPhaseID()
	}
	m.buyingPhases[bp.ID] = bp
	m.openBuyingPhase[bp.TurnoverTimeID] = bp.ID

	t := m.turnoverTimes[bp.TurnoverTimeID]
	t.CirculationOpen = true
	t.UpdatedAt = m.now()
	m.turnoverTimes[t.ID] = t
	return bp, nil
}

func (m *Memory) CloseBuyingPhase(_ context.Context, id circulation.TurnoverTimeID, bpID circulation.BuyingPhaseID) (circulation.BuyingPhase, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	openID, hasOpen := m.openBuyingPhase[id]
	if !hasOpen {
		return circulation.BuyingPhase{}, circulation.ErrNoOpenBuyingPhase
	}
	if openID != bpID {
		return circulation.BuyingPhase{}, circulation.ErrNoOpenBuyingPhase
	}
	bp := m.buyingPhases[bpID]
	now := m.now()
	bp.ClosedAt = &now
	m.buyingPhases[bpID] = bp
	delete(m.openBuyingPhase, id)

	t := m.turnoverTimes[id]
	t.Circulation.BuyingTimeNanos += bp.ElapsedNanos()
	// All phases closed — circulation is no longer open.
	if _, stillOpen := m.openSellingPhase[id]; !stillOpen {
		t.CirculationOpen = false
	}
	t.UpdatedAt = now
	m.turnoverTimes[id] = t
	return bp, nil
}

func (m *Memory) SetPerishability(_ context.Context, p circulation.Perishability) (circulation.Perishability, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p.ID.IsZero() {
		p.ID = circulation.NewPerishabilityID()
	}
	m.perishabilities[p.CommodityID] = p
	return p, nil
}

func (m *Memory) GetPerishability(_ context.Context, commodityID circulation.CommodityID) (circulation.Perishability, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.perishabilities[commodityID]
	if !ok {
		return circulation.Perishability{}, ErrNotFound
	}
	return p, nil
}

func (m *Memory) SetMarketSeparation(_ context.Context, ms circulation.MarketSeparation) (circulation.MarketSeparation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if ms.ID.IsZero() {
		ms.ID = circulation.NewMarketSeparationID()
	}
	m.marketSeparations[ms.IndustrialCapitalID] = ms
	return ms, nil
}
