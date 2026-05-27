package store

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/circulation"
	"github.com/theding0x/capital-simulator/services/market-service/internal/costs"
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

	// Vol. II Ch. 14 — circulation-time refinements.
	sellingPhaseDetails          map[circulation.SellingPhaseDetailedID]circulation.SellingPhaseDetailed
	buyingPhaseDetails           map[circulation.BuyingPhaseDetailedID]circulation.BuyingPhaseDetailed
	distanceLagRelations         map[circulation.DistanceLagRelationID]circulation.DistanceLagRelation
	circulationSpeedImprovements map[circulation.CirculationSpeedImprovementID]circulation.CirculationSpeedImprovement
	// annualSurplusPenalties keyed by IndustrialCapitalID+":"+Period.
	annualSurplusPenalties map[string]circulation.AnnualSurplusPenalty

	// Vol. II Ch. 6 — costs of circulation.
	circulationCosts    map[costs.CirculationCostID]costs.CirculationCost
	circulationAgents   map[costs.CirculationAgentID]costs.CirculationAgent
	moneyAsFauxFrais    map[costs.MoneyAsFauxFraisID]costs.MoneyAsFauxFrais
	commoditySupplies   map[costs.CommoditySupplyID]costs.CommoditySupply
	storageCosts        map[costs.StorageCostID]costs.StorageCost
	transportTariffs    map[costs.CommodityID]costs.TransportTariff
	transportLegs       map[costs.TransportLegID]costs.TransportLeg
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
		// Vol. II Ch. 14
		sellingPhaseDetails:          make(map[circulation.SellingPhaseDetailedID]circulation.SellingPhaseDetailed),
		buyingPhaseDetails:           make(map[circulation.BuyingPhaseDetailedID]circulation.BuyingPhaseDetailed),
		distanceLagRelations:         make(map[circulation.DistanceLagRelationID]circulation.DistanceLagRelation),
		circulationSpeedImprovements: make(map[circulation.CirculationSpeedImprovementID]circulation.CirculationSpeedImprovement),
		annualSurplusPenalties:       make(map[string]circulation.AnnualSurplusPenalty),
		// Vol. II Ch. 6
		circulationCosts:  make(map[costs.CirculationCostID]costs.CirculationCost),
		circulationAgents: make(map[costs.CirculationAgentID]costs.CirculationAgent),
		moneyAsFauxFrais:  make(map[costs.MoneyAsFauxFraisID]costs.MoneyAsFauxFrais),
		commoditySupplies: make(map[costs.CommoditySupplyID]costs.CommoditySupply),
		storageCosts:      make(map[costs.StorageCostID]costs.StorageCost),
		transportTariffs:  make(map[costs.CommodityID]costs.TransportTariff),
		transportLegs:     make(map[costs.TransportLegID]costs.TransportLeg),
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

// --- Vol. II Ch. 6 — CirculationCostStore -----------------------------------

func (m *Memory) CreateCirculationCost(_ context.Context, c costs.CirculationCost) (costs.CirculationCost, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if c.ID == "" {
		c.ID = costs.NewCirculationCostID()
	}
	if c.IncurredAt.IsZero() {
		c.IncurredAt = m.now()
	}
	c.Nature = costs.NatureOf(c.Kind)
	m.circulationCosts[c.ID] = c
	return c, nil
}

func (m *Memory) GetCirculationCost(_ context.Context, id costs.CirculationCostID) (costs.CirculationCost, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.circulationCosts[id]
	if !ok {
		return costs.CirculationCost{}, ErrNotFound
	}
	return c, nil
}

func (m *Memory) ListCirculationCosts(_ context.Context, icID costs.IndustrialCapitalID, kind costs.CirculationCostKind, nature costs.CostNature) ([]costs.CirculationCost, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []costs.CirculationCost
	for _, c := range m.circulationCosts {
		if icID != "" && c.IndustrialCapitalID != icID {
			continue
		}
		if kind != "" && c.Kind != kind {
			continue
		}
		if nature != "" && c.Nature != nature {
			continue
		}
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].IncurredAt.Before(out[j].IncurredAt) })
	return out, nil
}

func (m *Memory) CreateCirculationAgent(_ context.Context, a costs.CirculationAgent) (costs.CirculationAgent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if a.ID == "" {
		a.ID = costs.NewCirculationAgentID()
	}
	m.circulationAgents[a.ID] = a
	return a, nil
}

func (m *Memory) ListCirculationAgents(_ context.Context, icID costs.IndustrialCapitalID) ([]costs.CirculationAgent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []costs.CirculationAgent
	for _, a := range m.circulationAgents {
		if icID != "" && a.IndustrialCapitalID != icID {
			continue
		}
		out = append(out, a)
	}
	return out, nil
}

func (m *Memory) CreateMoneyAsFauxFrais(_ context.Context, mf costs.MoneyAsFauxFrais) (costs.MoneyAsFauxFrais, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if mf.ID == "" {
		mf.ID = costs.NewMoneyAsFauxFraisID()
	}
	m.moneyAsFauxFrais[mf.ID] = mf
	return mf, nil
}

func (m *Memory) ListMoneyAsFauxFrais(_ context.Context) ([]costs.MoneyAsFauxFrais, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]costs.MoneyAsFauxFrais, 0, len(m.moneyAsFauxFrais))
	for _, mf := range m.moneyAsFauxFrais {
		out = append(out, mf)
	}
	return out, nil
}

func (m *Memory) CreateCommoditySupply(_ context.Context, s costs.CommoditySupply) (costs.CommoditySupply, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s.ID == "" {
		s.ID = costs.NewCommoditySupplyID()
	}
	if s.HeldSince.IsZero() {
		s.HeldSince = m.now()
	}
	m.commoditySupplies[s.ID] = s
	return s, nil
}

func (m *Memory) GetCommoditySupply(_ context.Context, id costs.CommoditySupplyID) (costs.CommoditySupply, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.commoditySupplies[id]
	if !ok {
		return costs.CommoditySupply{}, ErrNotFound
	}
	return s, nil
}

func (m *Memory) ListCommoditySupplies(_ context.Context, icID costs.IndustrialCapitalID) ([]costs.CommoditySupply, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []costs.CommoditySupply
	for _, s := range m.commoditySupplies {
		if icID != "" && s.IndustrialCapitalID != icID {
			continue
		}
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].HeldSince.Before(out[j].HeldSince) })
	return out, nil
}

func (m *Memory) CreateStorageCost(_ context.Context, sc costs.StorageCost) (costs.StorageCost, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sc.ID == "" {
		sc.ID = costs.NewStorageCostID()
	}
	m.storageCosts[sc.ID] = sc
	return sc, nil
}

func (m *Memory) ListStorageCosts(_ context.Context, supplyID costs.CommoditySupplyID) ([]costs.StorageCost, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []costs.StorageCost
	for _, sc := range m.storageCosts {
		if sc.CommoditySupplyID == supplyID {
			out = append(out, sc)
		}
	}
	return out, nil
}

func (m *Memory) SetTransportTariff(_ context.Context, t costs.TransportTariff) (costs.TransportTariff, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if t.ID == "" {
		t.ID = costs.NewTransportTariffID()
	}
	m.transportTariffs[t.CommodityID] = t
	return t, nil
}

func (m *Memory) GetTransportTariff(_ context.Context, commodityID costs.CommodityID) (costs.TransportTariff, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.transportTariffs[commodityID]
	if !ok {
		return costs.TransportTariff{}, ErrNotFound
	}
	return t, nil
}

func (m *Memory) CreateTransportLeg(_ context.Context, leg costs.TransportLeg) (costs.TransportLeg, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if leg.ID == "" {
		leg.ID = costs.NewTransportLegID()
	}
	m.transportLegs[leg.ID] = leg
	return leg, nil
}

func (m *Memory) ListTransportLegs(_ context.Context, commodityID costs.CommodityID) ([]costs.TransportLeg, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []costs.TransportLeg
	for _, leg := range m.transportLegs {
		if commodityID != "" && leg.CommodityID != commodityID {
			continue
		}
		out = append(out, leg)
	}
	return out, nil
}

func (m *Memory) AggregateForCapital(_ context.Context, icID costs.IndustrialCapitalID, period string) (costs.AggregateResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := costs.AggregateResult{
		IndustrialCapitalID: icID,
		Period:              period,
		Items:               []costs.CirculationCost{},
	}
	for _, c := range m.circulationCosts {
		if c.IndustrialCapitalID != icID {
			continue
		}
		res.Items = append(res.Items, c)
		if c.Nature == costs.NatureValueAdding {
			res.TotalTransportPence += c.Pence
		} else {
			res.TotalFauxFraisPence += c.Pence
		}
	}
	sort.Slice(res.Items, func(i, j int) bool { return res.Items[i].IncurredAt.Before(res.Items[j].IncurredAt) })
	return res, nil
}

func (m *Memory) SystemFauxFrais(_ context.Context, period string) (costs.SystemFauxFraisResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := costs.SystemFauxFraisResult{
		Period: period,
		Items:  []costs.MoneyAsFauxFrais{},
	}
	for _, mf := range m.moneyAsFauxFrais {
		res.Items = append(res.Items, mf)
		res.TotalReservePence += mf.Pence
		res.TotalAnnualReplacementPence += mf.AnnualReplacementPence
	}
	return res, nil
}

// --- Vol. II Ch. 14 — circulation-time refinements -------------------------

func (m *Memory) CreateSellingPhaseDetailed(_ context.Context, d circulation.SellingPhaseDetailed) (circulation.SellingPhaseDetailed, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if d.ID.IsZero() {
		d.ID = circulation.NewSellingPhaseDetailedID()
	}
	d.CreatedAt = m.now()
	m.sellingPhaseDetails[d.ID] = d
	return d, nil
}

func (m *Memory) GetSellingPhaseDetailed(_ context.Context, id circulation.SellingPhaseDetailedID) (circulation.SellingPhaseDetailed, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.sellingPhaseDetails[id]
	if !ok {
		return circulation.SellingPhaseDetailed{}, ErrNotFound
	}
	return d, nil
}

func (m *Memory) CreateBuyingPhaseDetailed(_ context.Context, d circulation.BuyingPhaseDetailed) (circulation.BuyingPhaseDetailed, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if d.ID.IsZero() {
		d.ID = circulation.NewBuyingPhaseDetailedID()
	}
	d.CreatedAt = m.now()
	m.buyingPhaseDetails[d.ID] = d
	return d, nil
}

func (m *Memory) GetBuyingPhaseDetailed(_ context.Context, id circulation.BuyingPhaseDetailedID) (circulation.BuyingPhaseDetailed, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	d, ok := m.buyingPhaseDetails[id]
	if !ok {
		return circulation.BuyingPhaseDetailed{}, ErrNotFound
	}
	return d, nil
}

func (m *Memory) CreateDistanceLagRelation(_ context.Context, r circulation.DistanceLagRelation) (circulation.DistanceLagRelation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if r.ID.IsZero() {
		r.ID = circulation.NewDistanceLagRelationID()
	}
	r.CreatedAt = m.now()
	m.distanceLagRelations[r.ID] = r
	return r, nil
}

func (m *Memory) GetDistanceLagRelation(_ context.Context, id circulation.DistanceLagRelationID) (circulation.DistanceLagRelation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	r, ok := m.distanceLagRelations[id]
	if !ok {
		return circulation.DistanceLagRelation{}, ErrNotFound
	}
	return r, nil
}

func (m *Memory) ListDistanceLagRelations(_ context.Context) ([]circulation.DistanceLagRelation, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]circulation.DistanceLagRelation, 0, len(m.distanceLagRelations))
	for _, r := range m.distanceLagRelations {
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}

func (m *Memory) CreateCirculationSpeedImprovement(_ context.Context, csi circulation.CirculationSpeedImprovement) (circulation.CirculationSpeedImprovement, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if csi.ID.IsZero() {
		csi.ID = circulation.NewCirculationSpeedImprovementID()
	}
	csi.CreatedAt = m.now()
	m.circulationSpeedImprovements[csi.ID] = csi
	return csi, nil
}

func (m *Memory) ListCirculationSpeedImprovements(_ context.Context, marketID circulation.MarketID) ([]circulation.CirculationSpeedImprovement, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]circulation.CirculationSpeedImprovement, 0)
	for _, csi := range m.circulationSpeedImprovements {
		if marketID == "" || csi.MarketID == marketID {
			out = append(out, csi)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].EffectiveAt.Before(out[j].EffectiveAt) })
	return out, nil
}

func (m *Memory) CreateAnnualSurplusPenalty(_ context.Context, asp circulation.AnnualSurplusPenalty) (circulation.AnnualSurplusPenalty, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if asp.ID.IsZero() {
		asp.ID = circulation.NewAnnualSurplusPenaltyID()
	}
	asp.PenaltyBasisPoints = asp.ComputePenalty()
	key := string(asp.IndustrialCapitalID) + ":" + asp.Period
	m.annualSurplusPenalties[key] = asp
	return asp, nil
}

func (m *Memory) GetAnnualSurplusPenalty(_ context.Context, icID circulation.IndustrialCapitalID, period string) (circulation.AnnualSurplusPenalty, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	key := string(icID) + ":" + period
	asp, ok := m.annualSurplusPenalties[key]
	if !ok {
		return circulation.AnnualSurplusPenalty{}, ErrNotFound
	}
	return asp, nil
}
