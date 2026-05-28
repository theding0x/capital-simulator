package store

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	tv "github.com/theding0x/capital-simulator/services/simulation-engine/internal/turnover"
)

// Memory is the in-memory implementation of all store interfaces.
type Memory struct {
	mu                 sync.RWMutex
	machines           map[machinery.MachineID]machinery.Machine
	factories          map[machinery.FactoryID]machinery.Factory
	ticks              map[machinery.FactoryID][]engine.Tick
	generalLaw         map[simulation.GeneralLawScenarioID]simulation.GeneralLawScenario
	historicalStages   map[simulation.HistoricalStageID]simulation.HistoricalStage
	stageNames         map[string]simulation.HistoricalStageID
	enclosureEvents    []simulation.EnclosureEvent
	wageStatutes       []simulation.WageStatute
	vagrancyLaws       []simulation.VagrancyLaw
	farmTenures        []simulation.FarmTenure
	domesticIndustries []simulation.DomesticIndustry
	capitalOrigins     []simulation.CapitalOrigin
	colonialTransfers  []simulation.ColonialTransfer
	nationalDebts      []simulation.NationalDebt
	protectionSystems  []simulation.ProtectionSystem
	trajectories       map[simulation.AccumulationTrajectoryID]simulation.AccumulationTrajectory
	colonialMarkets    map[simulation.ColonialLabourMarketID]simulation.ColonialLabourMarket
	colonyNames        map[string]simulation.ColonialLabourMarketID
	productiveCircuits  map[circulation.ProductiveCircuitID]circulation.ProductiveCircuit
	commodityCircuits   map[circulation.CommodityCircuitID]circulation.CommodityCircuit
	moneyCircuits       map[circulation.MoneyCircuitID]circulation.MoneyCircuit
	industrialCapitals  map[circulation.IndustrialCapitalID]circulation.IndustrialCapital
	capitalParts        map[circulation.IndustrialCapitalID][]circulation.CapitalPart
	stageDistributions  map[circulation.IndustrialCapitalID][]circulation.StageDistribution
	stageBlocks         map[circulation.IndustrialCapitalID][]circulation.StageBlock
	valueRevolutions    map[circulation.IndustrialCapitalID][]circulation.ValueRevolutionResult
	interlocks          map[circulation.IndustrialCapitalID][]circulation.MetamorphosisInterlock
	supplyDemand        map[circulation.IndustrialCapitalID][]circulation.SupplyDemandImbalance
	sinkingFunds        map[circulation.IndustrialCapitalID]circulation.SinkingFund
	turnovers            map[tv.TurnoverID]tv.Turnover
	turnoverCycles       map[tv.TurnoverID][]tv.TurnoverCycle
	turnoverNumbers      map[tv.TurnoverID]tv.TurnoverNumber
	composition           *memoryComposition
	aggregateTurnovers    *memoryAggregateTurnover
	economistAttributions []circulation.EconomistAttribution
	workingPeriod         *memoryWorkingPeriod
	productionTime        *memoryProductionTime
	priceRevolution       *memoryPriceRevolution
	valorisation          *memoryValorisation
	reproduction          *memoryReproduction
	moneyCapital          *memoryMoneyCapital
	simpleRepro           *memorySimpleReproduction
	extendedRepro         *memoryExtendedReproduction
	now                   func() time.Time
}

func NewMemory() *Memory {
	return &Memory{
		machines:           make(map[machinery.MachineID]machinery.Machine),
		factories:          make(map[machinery.FactoryID]machinery.Factory),
		ticks:              make(map[machinery.FactoryID][]engine.Tick),
		generalLaw:         make(map[simulation.GeneralLawScenarioID]simulation.GeneralLawScenario),
		historicalStages:   make(map[simulation.HistoricalStageID]simulation.HistoricalStage),
		stageNames:         make(map[string]simulation.HistoricalStageID),
		trajectories:       make(map[simulation.AccumulationTrajectoryID]simulation.AccumulationTrajectory),
		colonialMarkets:    make(map[simulation.ColonialLabourMarketID]simulation.ColonialLabourMarket),
		colonyNames:        make(map[string]simulation.ColonialLabourMarketID),
		productiveCircuits: make(map[circulation.ProductiveCircuitID]circulation.ProductiveCircuit),
		commodityCircuits:  make(map[circulation.CommodityCircuitID]circulation.CommodityCircuit),
		moneyCircuits:      make(map[circulation.MoneyCircuitID]circulation.MoneyCircuit),
		industrialCapitals: make(map[circulation.IndustrialCapitalID]circulation.IndustrialCapital),
		capitalParts:       make(map[circulation.IndustrialCapitalID][]circulation.CapitalPart),
		stageDistributions: make(map[circulation.IndustrialCapitalID][]circulation.StageDistribution),
		stageBlocks:        make(map[circulation.IndustrialCapitalID][]circulation.StageBlock),
		valueRevolutions:   make(map[circulation.IndustrialCapitalID][]circulation.ValueRevolutionResult),
		interlocks:         make(map[circulation.IndustrialCapitalID][]circulation.MetamorphosisInterlock),
		supplyDemand:       make(map[circulation.IndustrialCapitalID][]circulation.SupplyDemandImbalance),
		sinkingFunds:       make(map[circulation.IndustrialCapitalID]circulation.SinkingFund),
		turnovers:           make(map[tv.TurnoverID]tv.Turnover),
		turnoverCycles:      make(map[tv.TurnoverID][]tv.TurnoverCycle),
		turnoverNumbers:     make(map[tv.TurnoverID]tv.TurnoverNumber),
		composition:        newMemoryComposition(),
		aggregateTurnovers: newMemoryAggregateTurnover(),
		workingPeriod:      newMemoryWorkingPeriod(),
		productionTime:     newMemoryProductionTime(),
		priceRevolution:    newMemoryPriceRevolution(),
		valorisation:       newMemoryValorisation(),
		moneyCapital:  newMemoryMoneyCapital(),
		simpleRepro:   newMemorySimpleReproduction(),
		extendedRepro: newMemoryExtendedReproduction(),
		reproduction: func() *memoryReproduction {
			r := newMemoryReproduction()
			// Seed the canonical Ch. 17 realisation puzzle and 100-spinning-mill fixture.
			seedID := repro.SurplusCirculationID("5eed0000000000001701")
			r.circulations[seedID] = repro.SurplusCirculation{
				ID:                     seedID,
				Period:                 "1871",
				TotalSurplusPence:      50000000,
				RealisedSurplusPence:   50000000,
				UnrealisedSurplusPence: 0,
				RealisationSources:     []repro.RealisationSourceEntry{},
			}
			r.createdAt[seedID] = time.Now().UTC()
			r.realisationSources[seedID] = []repro.RealisationSourceEntry{
				{SurplusCirculationID: seedID, Source: repro.SourceCapitalistConsumption, Pence: 40000000},
				{SurplusCirculationID: seedID, Source: repro.SourceCapitalisedSurplus, Pence: 9500000},
				{SurplusCirculationID: seedID, Source: repro.SourceGoldProduction, Pence: 500000},
			}
			r.aggregates = []repro.SocialCapitalAggregate{
				{
					Period:                  "1871",
					TotalAdvancedPence:      500000000,
					TotalAnnualOutputPence:  550000000,
					TotalAnnualSurplusPence: 50000000,
					DepartmentIShareBP:      0,
					DepartmentIIShareBP:     0,
				},
			}
			r.realisationPuzzles = []repro.RealisationPuzzle{
				{
					SurplusCirculationID: seedID,
					PuzzleStatement:      "Where does the money come from to realise the aggregate surplus-value of all capitalists? Each individual capitalist can identify buyers for their own output, but at the social level the total money required to purchase C' seems to exceed what was advanced as M. The puzzle dissolves only when we treat total social capital as a whole — formally resolved in the reproduction schemes of Chs. 20-21.",
					ResolvedInChapter:    20,
				},
			}
			return r
		}(),
		economistAttributions: []circulation.EconomistAttribution{
			{
				ID:          "5eed000000000000001001",
				Concept:     "avances primitives",
				Theorist:    "Quesnay",
				EditionYear: 1758,
				Anticipates: "fixed_capital_item",
				Errors:      nil,
			},
			{
				ID:          "5eed000000000000001002",
				Concept:     "avances annuelles",
				Theorist:    "Quesnay",
				EditionYear: 1758,
				Anticipates: "circulating",
				Errors:      nil,
			},
			{
				ID:          "5eed000000000000001003",
				Concept:     "fixed and circulating stock",
				Theorist:    "Smith",
				EditionYear: 1776,
				Anticipates: "fixed_capital_item",
				Errors: []circulation.KnownEconomistError{
					circulation.ErrorSmithConflation,
					circulation.ErrorSmithCirculationCapitalConflation,
					circulation.ErrorSmithRevenueInCapital,
				},
			},
		},
		now: time.Now,
	}
}

func (m *Memory) CreateMachine(_ context.Context, mc machinery.Machine) (machinery.Machine, error) {
	if err := mc.Validate(); err != nil {
		return machinery.Machine{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if mc.ID.IsZero() {
		mc.ID = machinery.NewMachineID()
	}
	if _, ok := m.machines[mc.ID]; ok {
		return machinery.Machine{}, ErrAlreadyExists
	}
	mc.CreatedAt = m.now().UTC()
	m.machines[mc.ID] = mc
	return mc, nil
}

func (m *Memory) GetMachine(_ context.Context, id machinery.MachineID) (machinery.Machine, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mc, ok := m.machines[id]
	if !ok {
		return machinery.Machine{}, ErrNotFound
	}
	return mc, nil
}

func (m *Memory) ListMachines(_ context.Context) ([]machinery.Machine, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]machinery.Machine, 0, len(m.machines))
	for _, mc := range m.machines {
		out = append(out, mc)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (m *Memory) UpdateMachine(_ context.Context, id machinery.MachineID, u MachineUpdate) (machinery.Machine, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cur, ok := m.machines[id]
	if !ok {
		return machinery.Machine{}, ErrNotFound
	}
	if u.AccumulatedWear != nil {
		cur.AccumulatedWear = *u.AccumulatedWear
	}
	if u.AccumulatedDepreciation != nil {
		cur.AccumulatedDepreciation = *u.AccumulatedDepreciation
	}
	m.machines[id] = cur
	return cur, nil
}

func (m *Memory) CreateFactory(ctx context.Context, f machinery.Factory) (machinery.Factory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if f.ID.IsZero() {
		f.ID = machinery.NewFactoryID()
	}
	if _, ok := m.factories[f.ID]; ok {
		return machinery.Factory{}, ErrAlreadyExists
	}
	// Resolve machine references against the store. Factory may be created
	// either with embedded Machine records or with bare IDs; we always
	// persist the embedded form.
	resolved := make([]machinery.Machine, 0, len(f.Machines))
	for _, mc := range f.Machines {
		if !mc.ID.IsZero() {
			if existing, ok := m.machines[mc.ID]; ok {
				resolved = append(resolved, existing)
				continue
			}
		}
		// Inline machine: persist it transparently.
		if err := mc.Validate(); err != nil {
			return machinery.Factory{}, err
		}
		if mc.ID.IsZero() {
			mc.ID = machinery.NewMachineID()
		}
		mc.CreatedAt = m.now().UTC()
		m.machines[mc.ID] = mc
		resolved = append(resolved, mc)
	}
	f.Machines = resolved
	if err := f.Validate(); err != nil {
		return machinery.Factory{}, err
	}
	f.CreatedAt = m.now().UTC()
	m.factories[f.ID] = f
	return f, nil
}

func (m *Memory) GetFactory(_ context.Context, id machinery.FactoryID) (machinery.Factory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.factories[id]
	if !ok {
		return machinery.Factory{}, ErrNotFound
	}
	// Refresh embedded machines from the canonical store map so accumulated
	// wear reads come back current.
	refreshed := make([]machinery.Machine, 0, len(f.Machines))
	for _, mc := range f.Machines {
		if cur, ok := m.machines[mc.ID]; ok {
			refreshed = append(refreshed, cur)
		} else {
			refreshed = append(refreshed, mc)
		}
	}
	f.Machines = refreshed
	return f, nil
}

func (m *Memory) ListFactories(_ context.Context) ([]machinery.Factory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]machinery.Factory, 0, len(m.factories))
	for _, f := range m.factories {
		out = append(out, f)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (m *Memory) AdvanceTick(_ context.Context, id machinery.FactoryID) (machinery.Factory, engine.Tick, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f, ok := m.factories[id]
	if !ok {
		return machinery.Factory{}, engine.Tick{}, ErrNotFound
	}
	// Refresh embedded machines.
	live := make([]machinery.Machine, 0, len(f.Machines))
	for _, mc := range f.Machines {
		if cur, ok := m.machines[mc.ID]; ok {
			live = append(live, cur)
		} else {
			live = append(live, mc)
		}
	}
	f.Machines = live
	result := f.RunTick()
	// Accumulate wear per machine.
	for i, mc := range f.Machines {
		dwt := machinery.DailyWearAndTear(mc)
		mc.AccumulatedWear.Value += dwt
		f.Machines[i] = mc
		m.machines[mc.ID] = mc
	}
	f.TickCount++
	m.factories[id] = f
	tick := engine.Tick{
		FactoryID:        string(id),
		Sequence:         f.TickCount,
		ValueTransferred: int64(result.ValueTransferred),
		UnitsProduced:    result.UnitsProduced,
		HandLabourSaved:  int64(result.HandLabourSaved),
		OccurredAt:       m.now().UTC(),
	}
	m.ticks[id] = append(m.ticks[id], tick)
	return f, tick, nil
}

func (m *Memory) CreateGeneralLawScenario(_ context.Context, s simulation.GeneralLawScenario) (simulation.GeneralLawScenario, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s.ID.IsZero() {
		s.ID = simulation.NewGeneralLawScenarioID()
	}
	if _, ok := m.generalLaw[s.ID]; ok {
		return simulation.GeneralLawScenario{}, ErrAlreadyExists
	}
	s.CreatedAt = m.now().UTC()
	m.generalLaw[s.ID] = s
	return s, nil
}

func (m *Memory) GetGeneralLawScenario(_ context.Context, id simulation.GeneralLawScenarioID) (simulation.GeneralLawScenario, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.generalLaw[id]
	if !ok {
		return simulation.GeneralLawScenario{}, ErrNotFound
	}
	return s, nil
}

func (m *Memory) CreateHistoricalStage(_ context.Context, h simulation.HistoricalStage) (simulation.HistoricalStage, error) {
	if err := h.Validate(); err != nil {
		return simulation.HistoricalStage{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.stageNames[strings.ToLower(h.Name)]; ok {
		return simulation.HistoricalStage{}, ErrAlreadyExists
	}
	if h.ID.IsZero() {
		h.ID = simulation.NewHistoricalStageID()
	}
	if _, ok := m.historicalStages[h.ID]; ok {
		return simulation.HistoricalStage{}, ErrAlreadyExists
	}
	h.CreatedAt = m.now().UTC()
	stored := simulation.HistoricalStage{
		ID:                     h.ID,
		Name:                   h.Name,
		Description:            h.Description,
		PrimitiveAccumulations: append([]simulation.PrimitiveAccumulation(nil), h.PrimitiveAccumulations...),
		CreatedAt:              h.CreatedAt,
	}
	m.historicalStages[h.ID] = stored
	m.stageNames[strings.ToLower(h.Name)] = h.ID
	return stored, nil
}

func (m *Memory) GetHistoricalStage(_ context.Context, id simulation.HistoricalStageID) (simulation.HistoricalStage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	h, ok := m.historicalStages[id]
	if !ok {
		return simulation.HistoricalStage{}, ErrNotFound
	}
	cp := h
	cp.PrimitiveAccumulations = append([]simulation.PrimitiveAccumulation(nil), h.PrimitiveAccumulations...)
	return cp, nil
}

func (m *Memory) ListHistoricalStages(_ context.Context) ([]simulation.HistoricalStage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.HistoricalStage, 0, len(m.historicalStages))
	for _, h := range m.historicalStages {
		cp := h
		cp.PrimitiveAccumulations = append([]simulation.PrimitiveAccumulation(nil), h.PrimitiveAccumulations...)
		out = append(out, cp)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (m *Memory) ListTicks(_ context.Context, id machinery.FactoryID, limit int) ([]engine.Tick, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	all := m.ticks[id]
	if limit <= 0 || limit > len(all) {
		limit = len(all)
	}
	// Return the latest `limit` ticks in ascending-sequence order.
	start := len(all) - limit
	out := make([]engine.Tick, limit)
	copy(out, all[start:])
	return out, nil
}

func (m *Memory) CreateEnclosureEvent(_ context.Context, e simulation.EnclosureEvent) (simulation.EnclosureEvent, error) {
	if err := e.Validate(); err != nil {
		return simulation.EnclosureEvent{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if e.ID.IsZero() {
		e.ID = simulation.NewEnclosureEventID()
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = m.now()
	}
	m.enclosureEvents = append(m.enclosureEvents, e)
	return e, nil
}

func (m *Memory) ListEnclosureEvents(_ context.Context) ([]simulation.EnclosureEvent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.EnclosureEvent, len(m.enclosureEvents))
	copy(out, m.enclosureEvents)
	return out, nil
}

func (m *Memory) CreateWageStatute(_ context.Context, w simulation.WageStatute) (simulation.WageStatute, error) {
	if err := w.Validate(); err != nil {
		return simulation.WageStatute{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if w.ID.IsZero() {
		w.ID = simulation.NewWageStatuteID()
	}
	if w.CreatedAt.IsZero() {
		w.CreatedAt = m.now()
	}
	m.wageStatutes = append(m.wageStatutes, w)
	return w, nil
}

func (m *Memory) ListWageStatutesByStage(_ context.Context, stageID simulation.HistoricalStageID) ([]simulation.WageStatute, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.WageStatute, 0)
	for _, w := range m.wageStatutes {
		if w.HistoricalStageID == stageID {
			out = append(out, w)
		}
	}
	return out, nil
}

func (m *Memory) CreateVagrancyLaw(_ context.Context, v simulation.VagrancyLaw) (simulation.VagrancyLaw, error) {
	if err := v.Validate(); err != nil {
		return simulation.VagrancyLaw{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if v.ID.IsZero() {
		v.ID = simulation.NewVagrancyLawID()
	}
	if v.CreatedAt.IsZero() {
		v.CreatedAt = m.now()
	}
	m.vagrancyLaws = append(m.vagrancyLaws, v)
	return v, nil
}

func (m *Memory) ListVagrancyLawsByStage(_ context.Context, stageID simulation.HistoricalStageID) ([]simulation.VagrancyLaw, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.VagrancyLaw, 0)
	for _, v := range m.vagrancyLaws {
		if v.HistoricalStageID == stageID {
			out = append(out, v)
		}
	}
	return out, nil
}

func (m *Memory) CreateFarmTenure(_ context.Context, f simulation.FarmTenure) (simulation.FarmTenure, error) {
	if err := f.Validate(); err != nil {
		return simulation.FarmTenure{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if f.ID.IsZero() {
		f.ID = simulation.NewFarmTenureID()
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = m.now()
	}
	m.farmTenures = append(m.farmTenures, f)
	return f, nil
}

func (m *Memory) ListFarmTenuresByStage(_ context.Context, stageID simulation.HistoricalStageID) ([]simulation.FarmTenure, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.FarmTenure, 0)
	for _, f := range m.farmTenures {
		if f.HistoricalStageID == stageID {
			out = append(out, f)
		}
	}
	return out, nil
}

func (m *Memory) CreateDomesticIndustry(_ context.Context, d simulation.DomesticIndustry) (simulation.DomesticIndustry, error) {
	if err := d.Validate(); err != nil {
		return simulation.DomesticIndustry{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if d.ID.IsZero() {
		d.ID = simulation.NewDomesticIndustryID()
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = m.now()
	}
	m.domesticIndustries = append(m.domesticIndustries, d)
	return d, nil
}

func (m *Memory) ListDomesticIndustriesByStage(_ context.Context, stageID simulation.HistoricalStageID) ([]simulation.DomesticIndustry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.DomesticIndustry, 0)
	for _, d := range m.domesticIndustries {
		if d.HistoricalStageID == stageID {
			out = append(out, d)
		}
	}
	return out, nil
}

func (m *Memory) CreateCapitalOrigin(_ context.Context, c simulation.CapitalOrigin) (simulation.CapitalOrigin, error) {
	if err := c.Validate(); err != nil {
		return simulation.CapitalOrigin{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if c.ID.IsZero() {
		c.ID = simulation.NewCapitalOriginID()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now()
	}
	m.capitalOrigins = append(m.capitalOrigins, c)
	return c, nil
}

func (m *Memory) ListCapitalOriginsByStage(_ context.Context, stageID simulation.HistoricalStageID) ([]simulation.CapitalOrigin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.CapitalOrigin, 0)
	for _, c := range m.capitalOrigins {
		if c.HistoricalStageID == stageID {
			out = append(out, c)
		}
	}
	return out, nil
}

func (m *Memory) CreateColonialTransfer(_ context.Context, t simulation.ColonialTransfer) (simulation.ColonialTransfer, error) {
	if err := t.Validate(); err != nil {
		return simulation.ColonialTransfer{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if t.ID.IsZero() {
		t.ID = simulation.NewColonialTransferID()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = m.now()
	}
	m.colonialTransfers = append(m.colonialTransfers, t)
	return t, nil
}

func (m *Memory) ListColonialTransfersByStage(_ context.Context, stageID simulation.HistoricalStageID) ([]simulation.ColonialTransfer, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.ColonialTransfer, 0)
	for _, t := range m.colonialTransfers {
		if t.HistoricalStageID == stageID {
			out = append(out, t)
		}
	}
	return out, nil
}

func (m *Memory) CreateNationalDebt(_ context.Context, d simulation.NationalDebt) (simulation.NationalDebt, error) {
	if err := d.Validate(); err != nil {
		return simulation.NationalDebt{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if d.ID.IsZero() {
		d.ID = simulation.NewNationalDebtID()
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = m.now()
	}
	m.nationalDebts = append(m.nationalDebts, d)
	return d, nil
}

func (m *Memory) ListNationalDebtsByStage(_ context.Context, stageID simulation.HistoricalStageID) ([]simulation.NationalDebt, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.NationalDebt, 0)
	for _, d := range m.nationalDebts {
		if d.HistoricalStageID == stageID {
			out = append(out, d)
		}
	}
	return out, nil
}

func (m *Memory) ListProtectionSystemsByStage(_ context.Context, stageID simulation.HistoricalStageID) ([]simulation.ProtectionSystem, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.ProtectionSystem, 0)
	for _, s := range m.protectionSystems {
		if s.HistoricalStageID == stageID {
			out = append(out, s)
		}
	}
	return out, nil
}

func (m *Memory) CreateAccumulationTrajectory(_ context.Context, t simulation.AccumulationTrajectory) (simulation.AccumulationTrajectory, error) {
	if err := t.Validate(); err != nil {
		return simulation.AccumulationTrajectory{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if t.ID.IsZero() {
		t.ID = simulation.NewAccumulationTrajectoryID()
	}
	if _, ok := m.trajectories[t.ID]; ok {
		return simulation.AccumulationTrajectory{}, ErrAlreadyExists
	}
	// Mirror the MySQL uq_accumulation_trajectories_name unique constraint.
	wantName := strings.ToLower(t.Name)
	for _, existing := range m.trajectories {
		if strings.ToLower(existing.Name) == wantName {
			return simulation.AccumulationTrajectory{}, ErrAlreadyExists
		}
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = m.now().UTC()
	}
	if t.Steps == nil {
		t.Steps = []simulation.CentralisationStep{}
	}
	stepsCopy := make([]simulation.CentralisationStep, len(t.Steps))
	copy(stepsCopy, t.Steps)
	t.Steps = stepsCopy
	m.trajectories[t.ID] = t
	return t, nil
}

func (m *Memory) GetAccumulationTrajectory(_ context.Context, id simulation.AccumulationTrajectoryID) (simulation.AccumulationTrajectory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.trajectories[id]
	if !ok {
		return simulation.AccumulationTrajectory{}, ErrNotFound
	}
	stepsCopy := make([]simulation.CentralisationStep, len(t.Steps))
	copy(stepsCopy, t.Steps)
	t.Steps = stepsCopy
	return t, nil
}

func (m *Memory) ListAccumulationTrajectories(_ context.Context) ([]simulation.AccumulationTrajectory, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.AccumulationTrajectory, 0, len(m.trajectories))
	for _, t := range m.trajectories {
		stepsCopy := make([]simulation.CentralisationStep, len(t.Steps))
		copy(stepsCopy, t.Steps)
		t.Steps = stepsCopy
		out = append(out, t)
	}
	sort.Slice(out, func(i, j int) bool {
		if !out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].CreatedAt.Before(out[j].CreatedAt)
		}
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

// CreateColonialLabourMarket persists a Ch. 33 colonial labour market.
// Colony names are unique case-insensitively.
func (m *Memory) CreateColonialLabourMarket(_ context.Context, market simulation.ColonialLabourMarket) (simulation.ColonialLabourMarket, error) {
	if err := market.Validate(); err != nil {
		return simulation.ColonialLabourMarket{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	key := strings.ToLower(market.Colony)
	if _, ok := m.colonyNames[key]; ok {
		return simulation.ColonialLabourMarket{}, ErrAlreadyExists
	}
	if market.ID.IsZero() {
		market.ID = simulation.NewColonialLabourMarketID()
	}
	if _, ok := m.colonialMarkets[market.ID]; ok {
		return simulation.ColonialLabourMarket{}, ErrAlreadyExists
	}
	if market.CreatedAt.IsZero() {
		market.CreatedAt = m.now().UTC()
	}
	m.colonialMarkets[market.ID] = market
	m.colonyNames[key] = market.ID
	return market, nil
}

// GetColonialLabourMarket returns the persisted market or ErrNotFound.
func (m *Memory) GetColonialLabourMarket(_ context.Context, id simulation.ColonialLabourMarketID) (simulation.ColonialLabourMarket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	market, ok := m.colonialMarkets[id]
	if !ok {
		return simulation.ColonialLabourMarket{}, ErrNotFound
	}
	return market, nil
}

// ListColonialLabourMarkets returns markets ordered by created_at
// ascending then by colony name.
func (m *Memory) ListColonialLabourMarkets(_ context.Context) ([]simulation.ColonialLabourMarket, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.ColonialLabourMarket, 0, len(m.colonialMarkets))
	for _, market := range m.colonialMarkets {
		out = append(out, market)
	}
	sort.Slice(out, func(i, j int) bool {
		if !out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].CreatedAt.Before(out[j].CreatedAt)
		}
		return strings.ToLower(out[i].Colony) < strings.ToLower(out[j].Colony)
	})
	return out, nil
}

// UpdateColonialLabourMarket applies a partial regulation update.
func (m *Memory) UpdateColonialLabourMarket(_ context.Context, id simulation.ColonialLabourMarketID, u ColonialLabourMarketUpdate) (simulation.ColonialLabourMarket, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cur, ok := m.colonialMarkets[id]
	if !ok {
		return simulation.ColonialLabourMarket{}, ErrNotFound
	}
	if u.WakefieldSchemeApplied != nil {
		cur.WakefieldSchemeApplied = *u.WakefieldSchemeApplied
	}
	if u.IndependenceYears != nil {
		cur.IndependenceYears = *u.IndependenceYears
	}
	if u.SurplusLabourExtractable != nil {
		cur.SurplusLabourExtractable = *u.SurplusLabourExtractable
	}
	m.colonialMarkets[id] = cur
	return cur, nil
}

// RegulateColonialLabourMarket reads the market under the Memory
// mutex, runs the pure simulation.ColonialLabourRegulation against
// it, and writes the regulated state back atomically. Two concurrent
// /regulate callers see serialized reads-and-writes.
func (m *Memory) RegulateColonialLabourMarket(_ context.Context, id simulation.ColonialLabourMarketID, scheme simulation.SystematicColonisation) (simulation.ColonialLabourMarket, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cur, ok := m.colonialMarkets[id]
	if !ok {
		return simulation.ColonialLabourMarket{}, ErrNotFound
	}
	regulated := simulation.ColonialLabourRegulation(cur, scheme)
	m.colonialMarkets[id] = regulated
	return regulated, nil
}

// --- Vol. II Ch. 2 — ProductiveCircuitStore ---

func (m *Memory) CreateProductiveCircuit(_ context.Context, pc circulation.ProductiveCircuit) (circulation.ProductiveCircuit, error) {
	if err := pc.Validate(); err != nil {
		return circulation.ProductiveCircuit{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if pc.ID.IsZero() {
		pc.ID = circulation.NewProductiveCircuitID()
	}
	if _, ok := m.productiveCircuits[pc.ID]; ok {
		return circulation.ProductiveCircuit{}, ErrAlreadyExists
	}
	if pc.RevenueCircuits == nil {
		pc.RevenueCircuits = []circulation.RevenueCircuit{}
	}
	if pc.CapitalisationSteps == nil {
		pc.CapitalisationSteps = []circulation.CapitalisationStep{}
	}
	if pc.ReserveDraws == nil {
		pc.ReserveDraws = []circulation.ReserveDraw{}
	}
	pc.LatentMoneyCapital.ProductiveCircuitID = pc.ID
	pc.LatentMoneyCapital.Threshold = pc.MinCapitalisationIncrement
	pc.ReserveFund.ProductiveCircuitID = pc.ID
	pc.CreatedAt = m.now().UTC()
	m.productiveCircuits[pc.ID] = pc
	return pc, nil
}

func (m *Memory) GetProductiveCircuit(_ context.Context, id circulation.ProductiveCircuitID) (circulation.ProductiveCircuit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	pc, ok := m.productiveCircuits[id]
	if !ok {
		return circulation.ProductiveCircuit{}, ErrNotFound
	}
	pc.RevenueCircuits = append([]circulation.RevenueCircuit(nil), pc.RevenueCircuits...)
	pc.CapitalisationSteps = append([]circulation.CapitalisationStep(nil), pc.CapitalisationSteps...)
	pc.ReserveDraws = append([]circulation.ReserveDraw(nil), pc.ReserveDraws...)
	return pc, nil
}

func (m *Memory) ListProductiveCircuits(_ context.Context) ([]circulation.ProductiveCircuit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]circulation.ProductiveCircuit, 0, len(m.productiveCircuits))
	for _, pc := range m.productiveCircuits {
		pc.RevenueCircuits = append([]circulation.RevenueCircuit(nil), pc.RevenueCircuits...)
		pc.CapitalisationSteps = append([]circulation.CapitalisationStep(nil), pc.CapitalisationSteps...)
		pc.ReserveDraws = append([]circulation.ReserveDraw(nil), pc.ReserveDraws...)
		out = append(out, pc)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (m *Memory) RecordRevenue(_ context.Context, id circulation.ProductiveCircuitID, rc circulation.RevenueCircuit) (circulation.RevenueCircuit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	pc, ok := m.productiveCircuits[id]
	if !ok {
		return circulation.RevenueCircuit{}, ErrNotFound
	}
	rc.ProductiveCircuitID = id
	if rc.SpentAt.IsZero() {
		rc.SpentAt = m.now().UTC()
	}
	pc.RevenueCircuits = append(pc.RevenueCircuits, rc)
	m.productiveCircuits[id] = pc
	return rc, nil
}

func (m *Memory) Accumulate(_ context.Context, id circulation.ProductiveCircuitID, amount circulation.Pence) (circulation.LatentMoneyCapital, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	pc, ok := m.productiveCircuits[id]
	if !ok {
		return circulation.LatentMoneyCapital{}, ErrNotFound
	}
	pc.LatentMoneyCapital.Accumulated += amount
	m.productiveCircuits[id] = pc
	return pc.LatentMoneyCapital, nil
}

func (m *Memory) Capitalise(_ context.Context, id circulation.ProductiveCircuitID, amount, dc, dv circulation.Pence) (circulation.CapitalisationStep, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	pc, ok := m.productiveCircuits[id]
	if !ok {
		return circulation.CapitalisationStep{}, ErrNotFound
	}
	if amount > 0 && amount < pc.MinCapitalisationIncrement {
		return circulation.CapitalisationStep{}, circulation.ErrIndivisibleProductionElement
	}
	step := circulation.CapitalisationStep{
		ProductiveCircuitID: id,
		AmountInjected:      amount,
		DeltaConstantPence:  dc,
		DeltaVariablePence:  dv,
		OccurredAt:          m.now().UTC(),
	}
	pc.ConstantPence += dc
	pc.VariablePence += dv
	pc.LatentMoneyCapital.Accumulated -= amount
	pc.CapitalisationSteps = append(pc.CapitalisationSteps, step)
	m.productiveCircuits[id] = pc
	return step, nil
}

func (m *Memory) DepositReserve(_ context.Context, id circulation.ProductiveCircuitID, amount circulation.Pence) (circulation.ReserveFund, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	pc, ok := m.productiveCircuits[id]
	if !ok {
		return circulation.ReserveFund{}, ErrNotFound
	}
	pc.ReserveFund.Balance += amount
	m.productiveCircuits[id] = pc
	return pc.ReserveFund, nil
}

func (m *Memory) WithdrawReserve(_ context.Context, id circulation.ProductiveCircuitID, amount circulation.Pence, reason circulation.ReserveDrawReason) (circulation.ReserveDraw, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	pc, ok := m.productiveCircuits[id]
	if !ok {
		return circulation.ReserveDraw{}, ErrNotFound
	}
	if amount > pc.ReserveFund.Balance {
		return circulation.ReserveDraw{}, circulation.ErrInsufficientReserve
	}
	draw := circulation.ReserveDraw{
		ProductiveCircuitID: id,
		DrawnPence:          amount,
		Reason:              reason,
		OccurredAt:          m.now().UTC(),
	}
	pc.ReserveFund.Balance -= amount
	pc.ReserveDraws = append(pc.ReserveDraws, draw)
	m.productiveCircuits[id] = pc
	return draw, nil
}

// Vol. II Ch. 3 — CommodityCircuitStore implementation.

func (m *Memory) CreateCommodityCircuit(_ context.Context, cc circulation.CommodityCircuit) (circulation.CommodityCircuit, error) {
	if err := cc.Validate(); err != nil {
		return circulation.CommodityCircuit{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if cc.ID.IsZero() {
		cc.ID = circulation.NewCommodityCircuitID()
	}
	if _, ok := m.commodityCircuits[cc.ID]; ok {
		return circulation.CommodityCircuit{}, ErrAlreadyExists
	}
	cc.CreatedAt = m.now().UTC()
	if cc.PartialSales == nil {
		cc.PartialSales = []circulation.SuccessivePartialSale{}
	}
	if cc.MPSources == nil {
		cc.MPSources = []circulation.MeansOfProductionSource{}
	}
	m.commodityCircuits[cc.ID] = cc
	return cc, nil
}

func (m *Memory) GetCommodityCircuit(_ context.Context, id circulation.CommodityCircuitID) (circulation.CommodityCircuit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cc, ok := m.commodityCircuits[id]
	if !ok {
		return circulation.CommodityCircuit{}, ErrNotFound
	}
	return cc, nil
}

func (m *Memory) ListCommodityCircuits(_ context.Context, agentID string) ([]circulation.CommodityCircuit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]circulation.CommodityCircuit, 0, len(m.commodityCircuits))
	for _, cc := range m.commodityCircuits {
		if agentID != "" && cc.AgentID != agentID {
			continue
		}
		out = append(out, cc)
	}
	return out, nil
}

func (m *Memory) RecordPartialSale(_ context.Context, id circulation.CommodityCircuitID, sale circulation.SuccessivePartialSale) (circulation.SuccessivePartialSale, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cc, ok := m.commodityCircuits[id]
	if !ok {
		return circulation.SuccessivePartialSale{}, ErrNotFound
	}
	sale.CommodityCircuitID = id
	if sale.SoldAt.IsZero() {
		sale.SoldAt = m.now().UTC()
	}
	cc.PartialSales = append(cc.PartialSales, sale)
	m.commodityCircuits[id] = cc
	return sale, nil
}

func (m *Memory) LinkMPSource(_ context.Context, id circulation.CommodityCircuitID, source circulation.MeansOfProductionSource) (circulation.MeansOfProductionSource, error) {
	if err := source.Validate(); err != nil {
		return circulation.MeansOfProductionSource{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	cc, ok := m.commodityCircuits[id]
	if !ok {
		return circulation.MeansOfProductionSource{}, ErrNotFound
	}
	source.CommodityCircuitID = id
	cc.MPSources = append(cc.MPSources, source)
	m.commodityCircuits[id] = cc
	return source, nil
}

func (m *Memory) CloseCommodityCircuit(_ context.Context, id circulation.CommodityCircuitID, aug circulation.CommodityAugmented) (circulation.CommodityCircuit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cc, ok := m.commodityCircuits[id]
	if !ok {
		return circulation.CommodityCircuit{}, ErrNotFound
	}
	aug.CommodityCircuitID = id
	if aug.ClosedAt.IsZero() {
		aug.ClosedAt = m.now().UTC()
	}
	cc.Terminal = &aug
	m.commodityCircuits[id] = cc
	return cc, nil
}

// --- Vol. II Ch. 1 — MoneyCircuitStore ---

func (m *Memory) CreateMoneyCircuit(_ context.Context, mc circulation.MoneyCircuit) (circulation.MoneyCircuit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if mc.ID.IsZero() {
		mc.ID = circulation.NewMoneyCircuitID()
	}
	if _, ok := m.moneyCircuits[mc.ID]; ok {
		return circulation.MoneyCircuit{}, ErrAlreadyExists
	}
	mc.Moment = circulation.MomentM
	mc.CreatedAt = m.now().UTC()
	if mc.Advance.AdvancedAt.IsZero() {
		mc.Advance.AdvancedAt = mc.CreatedAt
	}
	m.moneyCircuits[mc.ID] = mc
	return mc, nil
}

func (m *Memory) GetMoneyCircuit(_ context.Context, id circulation.MoneyCircuitID) (circulation.MoneyCircuit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	mc, ok := m.moneyCircuits[id]
	if !ok {
		return circulation.MoneyCircuit{}, ErrNotFound
	}
	return mc, nil
}

func (m *Memory) ListMoneyCircuits(_ context.Context, agentID string, moment circulation.CircuitMoment) ([]circulation.MoneyCircuit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]circulation.MoneyCircuit, 0, len(m.moneyCircuits))
	for _, mc := range m.moneyCircuits {
		if agentID != "" && mc.Advance.AgentID != agentID {
			continue
		}
		if moment != "" && mc.Moment != moment {
			continue
		}
		out = append(out, mc)
	}
	return out, nil
}

func (m *Memory) RecordPurchase(_ context.Context, id circulation.MoneyCircuitID, p circulation.PurchasePhase) (circulation.MoneyCircuit, error) {
	if err := p.Validate(); err != nil {
		return circulation.MoneyCircuit{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	mc, ok := m.moneyCircuits[id]
	if !ok {
		return circulation.MoneyCircuit{}, ErrNotFound
	}
	if err := mc.CanTransitionTo(circulation.MomentMC); err != nil {
		return circulation.MoneyCircuit{}, err
	}
	p.MoneyCircuitID = id
	mc.Purchase = &p
	mc.Moment = circulation.MomentMC
	m.moneyCircuits[id] = mc
	return mc, nil
}

func (m *Memory) RecordProductive(_ context.Context, id circulation.MoneyCircuitID, ps circulation.ProductiveState) (circulation.MoneyCircuit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mc, ok := m.moneyCircuits[id]
	if !ok {
		return circulation.MoneyCircuit{}, ErrNotFound
	}
	if err := mc.CanTransitionTo(circulation.MomentP); err != nil {
		return circulation.MoneyCircuit{}, err
	}
	if ps.TotalAdvance() != mc.Advance.Amount {
		return circulation.MoneyCircuit{}, circulation.ErrMagnitudeNotPreserved
	}
	ps.MoneyCircuitID = id
	if ps.EnteredAt.IsZero() {
		ps.EnteredAt = m.now().UTC()
	}
	mc.Productive = &ps
	mc.Moment = circulation.MomentP
	m.moneyCircuits[id] = mc
	return mc, nil
}

func (m *Memory) RecordCommodity(_ context.Context, id circulation.MoneyCircuitID, cc circulation.CommodityCapital) (circulation.MoneyCircuit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mc, ok := m.moneyCircuits[id]
	if !ok {
		return circulation.MoneyCircuit{}, ErrNotFound
	}
	if err := mc.CanTransitionTo(circulation.MomentCPrime); err != nil {
		return circulation.MoneyCircuit{}, err
	}
	if mc.Productive == nil {
		return circulation.MoneyCircuit{}, circulation.ErrNoProductionPhase
	}
	cc.MoneyCircuitID = id
	mc.Commodity = &cc
	mc.Moment = circulation.MomentCPrime
	m.moneyCircuits[id] = mc
	return mc, nil
}

func (m *Memory) RecordRealisation(_ context.Context, id circulation.MoneyCircuitID, r circulation.Realisation) (circulation.MoneyCircuit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	mc, ok := m.moneyCircuits[id]
	if !ok {
		return circulation.MoneyCircuit{}, ErrNotFound
	}
	if err := mc.CanTransitionTo(circulation.MomentCMPrime); err != nil {
		return circulation.MoneyCircuit{}, err
	}
	r.MoneyCircuitID = id
	if r.SoldAt.IsZero() {
		r.SoldAt = m.now().UTC()
	}
	mc.Realisation = &r
	mc.Moment = circulation.MomentMPrime
	m.moneyCircuits[id] = mc
	return mc, nil
}

// --- IndustrialCapitalStore (Vol. II Ch. 4) ---------------------------------

func (m *Memory) CreateIndustrialCapital(_ context.Context, ic circulation.IndustrialCapital) (circulation.IndustrialCapital, error) {
	if err := ic.Validate(); err != nil {
		return circulation.IndustrialCapital{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if ic.ID.IsZero() {
		ic.ID = circulation.NewIndustrialCapitalID()
	}
	if _, ok := m.industrialCapitals[ic.ID]; ok {
		return circulation.IndustrialCapital{}, ErrAlreadyExists
	}
	now := m.now().UTC()
	ic.CreatedAt = now
	ic.UpdatedAt = now
	if ic.Status == "" {
		ic.Status = circulation.StatusActive
	}
	m.industrialCapitals[ic.ID] = ic
	return ic, nil
}

func (m *Memory) GetIndustrialCapital(_ context.Context, id circulation.IndustrialCapitalID) (circulation.IndustrialCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	ic, ok := m.industrialCapitals[id]
	if !ok {
		return circulation.IndustrialCapital{}, ErrNotFound
	}
	if sds, ok2 := m.stageDistributions[id]; ok2 && len(sds) > 0 {
		latest := sds[len(sds)-1]
		ic.Latest = &latest
	}
	ic.OpenBlocks = nil
	for _, b := range m.stageBlocks[id] {
		if b.IsOpen() {
			ic.OpenBlocks = append(ic.OpenBlocks, b)
		}
	}
	if ic.OpenBlocks == nil {
		ic.OpenBlocks = []circulation.StageBlock{}
	}
	return ic, nil
}

func (m *Memory) ListIndustrialCapitals(_ context.Context, agentID, status, economyMode string) ([]circulation.IndustrialCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []circulation.IndustrialCapital{}
	for _, ic := range m.industrialCapitals {
		if agentID != "" && ic.AgentID != agentID {
			continue
		}
		if status != "" && string(ic.Status) != status {
			continue
		}
		if economyMode != "" && string(ic.EconomyMode) != economyMode {
			continue
		}
		out = append(out, ic)
	}
	return out, nil
}

func (m *Memory) RecordCapitalPart(_ context.Context, id circulation.IndustrialCapitalID, part circulation.CapitalPart) (circulation.CapitalPart, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.industrialCapitals[id]; !ok {
		return circulation.CapitalPart{}, ErrNotFound
	}
	if part.ID.IsZero() {
		part.ID = circulation.NewCapitalPartID()
	}
	part.IndustrialCapitalID = id
	if part.EnteredStageAt.IsZero() {
		part.EnteredStageAt = m.now().UTC()
	}
	m.capitalParts[id] = append(m.capitalParts[id], part)
	return part, nil
}

func (m *Memory) Snapshot(_ context.Context, id circulation.IndustrialCapitalID, sd circulation.StageDistribution) (circulation.StageDistribution, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ic, ok := m.industrialCapitals[id]
	if !ok {
		return circulation.StageDistribution{}, ErrNotFound
	}
	if err := sd.Validate(ic.TotalPence); err != nil {
		return circulation.StageDistribution{}, err
	}
	if sd.ID.IsZero() {
		sd.ID = circulation.NewStageDistributionID()
	}
	sd.IndustrialCapitalID = id
	if sd.At.IsZero() {
		sd.At = m.now().UTC()
	}
	m.stageDistributions[id] = append(m.stageDistributions[id], sd)
	return sd, nil
}

func (m *Memory) OpenBlock(_ context.Context, id circulation.IndustrialCapitalID, b circulation.StageBlock) (circulation.StageBlock, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ic, ok := m.industrialCapitals[id]
	if !ok {
		return circulation.StageBlock{}, ErrNotFound
	}
	if b.ID.IsZero() {
		b.ID = circulation.NewStageBlockID()
	}
	b.IndustrialCapitalID = id
	if b.OpenedAt.IsZero() {
		b.OpenedAt = m.now().UTC()
	}
	m.stageBlocks[id] = append(m.stageBlocks[id], b)
	openCount := int64(0)
	for _, sb := range m.stageBlocks[id] {
		if sb.IsOpen() {
			openCount++
		}
	}
	if ic.StagnationToleranceTicks > 0 && openCount > ic.StagnationToleranceTicks {
		ic.Status = circulation.StatusHalted
		m.industrialCapitals[id] = ic
	}
	return b, nil
}

func (m *Memory) CloseBlock(_ context.Context, id circulation.IndustrialCapitalID, blockID circulation.StageBlockID) (circulation.StageBlock, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.industrialCapitals[id]; !ok {
		return circulation.StageBlock{}, ErrNotFound
	}
	blocks := m.stageBlocks[id]
	for i, b := range blocks {
		if b.ID == blockID {
			if !b.IsOpen() {
				return b, nil
			}
			now := m.now().UTC()
			blocks[i].ClosedAt = &now
			m.stageBlocks[id] = blocks
			return blocks[i], nil
		}
	}
	return circulation.StageBlock{}, ErrNotFound
}

func (m *Memory) RecordValueRevolution(_ context.Context, res circulation.ValueRevolutionResult) (circulation.ValueRevolutionResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	id := res.Event.IndustrialCapitalID
	if _, ok := m.industrialCapitals[id]; !ok {
		return circulation.ValueRevolutionResult{}, ErrNotFound
	}
	if res.Event.ID.IsZero() {
		res.Event.ID = circulation.NewValueRevolutionEventID()
	}
	if res.Event.OccurredAt.IsZero() {
		res.Event.OccurredAt = m.now().UTC()
	}
	m.valueRevolutions[id] = append(m.valueRevolutions[id], res)
	return res, nil
}

func (m *Memory) RecordInterlock(_ context.Context, id circulation.IndustrialCapitalID, mi circulation.MetamorphosisInterlock) (circulation.MetamorphosisInterlock, error) {
	if err := mi.Validate(); err != nil {
		return circulation.MetamorphosisInterlock{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.industrialCapitals[id]; !ok {
		return circulation.MetamorphosisInterlock{}, ErrNotFound
	}
	if mi.ID.IsZero() {
		mi.ID = circulation.NewMetamorphosisInterlockID()
	}
	mi.BuyerIndustrialCapitalID = id
	if mi.OccurredAt.IsZero() {
		mi.OccurredAt = m.now().UTC()
	}
	m.interlocks[id] = append(m.interlocks[id], mi)
	return mi, nil
}

func (m *Memory) RecordSupplyDemand(_ context.Context, sdi circulation.SupplyDemandImbalance) (circulation.SupplyDemandImbalance, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.industrialCapitals[sdi.IndustrialCapitalID]; !ok {
		return circulation.SupplyDemandImbalance{}, ErrNotFound
	}
	if sdi.ID.IsZero() {
		sdi.ID = circulation.NewSupplyDemandImbalanceID()
	}
	m.supplyDemand[sdi.IndustrialCapitalID] = append(m.supplyDemand[sdi.IndustrialCapitalID], sdi)
	return sdi, nil
}

func (m *Memory) GetSupplyDemand(_ context.Context, id circulation.IndustrialCapitalID, period string) (circulation.SupplyDemandImbalance, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, sdi := range m.supplyDemand[id] {
		if sdi.Period == period {
			return sdi, nil
		}
	}
	return circulation.SupplyDemandImbalance{}, ErrNotFound
}

func (m *Memory) AggregateSupplyDemand(_ context.Context, period string) (circulation.AggregateSupplyDemandImbalance, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	agg := circulation.AggregateSupplyDemandImbalance{Period: period}
	for _, sdis := range m.supplyDemand {
		for _, sdi := range sdis {
			if sdi.Period == period {
				agg.DemandPence += sdi.DemandPence
				agg.SupplyPence += sdi.SupplyPence
				agg.ExcessPence += sdi.ExcessPence
			}
		}
	}
	return agg, nil
}

func (m *Memory) SetSinkingFund(_ context.Context, id circulation.IndustrialCapitalID, sf circulation.SinkingFund) (circulation.SinkingFund, error) {
	if err := sf.Validate(); err != nil {
		return circulation.SinkingFund{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.industrialCapitals[id]; !ok {
		return circulation.SinkingFund{}, ErrNotFound
	}
	if sf.ID.IsZero() {
		sf.ID = circulation.NewSinkingFundID()
	}
	sf.IndustrialCapitalID = id
	m.sinkingFunds[id] = sf
	return sf, nil
}

func (m *Memory) TickSinkingFund(_ context.Context, id circulation.IndustrialCapitalID) (circulation.SinkingFund, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	sf, ok := m.sinkingFunds[id]
	if !ok {
		return circulation.SinkingFund{}, ErrNotFound
	}
	sf = sf.Tick()
	m.sinkingFunds[id] = sf
	return sf, nil
}

// Vol. II Ch. 10 — Theories of Fixed and Circulating Capital.

// ListEconomistAttributions returns the pre-seeded attribution records,
// optionally filtered by theorist (case-insensitive). An empty theorist
// string returns all records.
func (m *Memory) ListEconomistAttributions(_ context.Context, theorist string) ([]circulation.EconomistAttribution, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if theorist == "" {
		out := make([]circulation.EconomistAttribution, len(m.economistAttributions))
		copy(out, m.economistAttributions)
		return out, nil
	}
	upper := strings.ToUpper(theorist)
	var out []circulation.EconomistAttribution
	for _, a := range m.economistAttributions {
		if strings.ToUpper(a.Theorist) == upper {
			out = append(out, a)
		}
	}
	if out == nil {
		out = []circulation.EconomistAttribution{}
	}
	return out, nil
}
