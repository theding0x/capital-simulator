// Package store is the persistence layer for simulation-engine. As of Ch. 15
// it persists Machine and Factory records and the tick log produced by
// advancing a Factory one period at a time. Ch. 25 adds GeneralLawScenario.
// Vol. II Ch. 2 adds ProductiveCircuit. Vol. II Ch. 16 adds AnnualSurplusRate.
// Vol. II Ch. 17 adds SurplusCirculation. Vol. II Ch. 20 adds
// SimpleReproductionScheme. Vol. II Ch. 21 adds ExtendedReproductionScheme.
package store

import (
	"context"
	"errors"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	comp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/composition"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	tv "github.com/theding0x/capital-simulator/services/simulation-engine/internal/turnover"
	val "github.com/theding0x/capital-simulator/services/simulation-engine/internal/valorisation"
	wp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/workperiod"
)

var (
	ErrNotFound      = errors.New("machinery: not found")
	ErrAlreadyExists = errors.New("machinery: already exists")
)

// MachineUpdate is a partial-update payload. Non-nil fields are applied.
type MachineUpdate struct {
	AccumulatedWear         *machinery.MaterialWear
	AccumulatedDepreciation *machinery.MoralDepreciation
}

func (u MachineUpdate) IsEmpty() bool {
	return u.AccumulatedWear == nil && u.AccumulatedDepreciation == nil
}

// MachineStore is the persistence contract for Machine records.
type MachineStore interface {
	CreateMachine(ctx context.Context, m machinery.Machine) (machinery.Machine, error)
	GetMachine(ctx context.Context, id machinery.MachineID) (machinery.Machine, error)
	ListMachines(ctx context.Context) ([]machinery.Machine, error)
	UpdateMachine(ctx context.Context, id machinery.MachineID, u MachineUpdate) (machinery.Machine, error)
}

// FactoryStore is the persistence contract for Factory records.
// CreateFactory persists the factory and links each machine_id. AdvanceTick
// atomically: re-reads the factory's machines, computes a tick result, writes
// it to factory_ticks, and updates each machine's accumulated wear.
// ListTicks returns the persisted tick history for a factory in
// ascending-sequence order; the simulation can replay the history or
// surface it in the UI.
type FactoryStore interface {
	CreateFactory(ctx context.Context, f machinery.Factory) (machinery.Factory, error)
	GetFactory(ctx context.Context, id machinery.FactoryID) (machinery.Factory, error)
	ListFactories(ctx context.Context) ([]machinery.Factory, error)
	AdvanceTick(ctx context.Context, id machinery.FactoryID) (machinery.Factory, engine.Tick, error)
	ListTicks(ctx context.Context, id machinery.FactoryID, limit int) ([]engine.Tick, error)
}

// GeneralLawStore is the persistence contract for Ch. 25 general-law scenarios.
type GeneralLawStore interface {
	CreateGeneralLawScenario(ctx context.Context, s simulation.GeneralLawScenario) (simulation.GeneralLawScenario, error)
	GetGeneralLawScenario(ctx context.Context, id simulation.GeneralLawScenarioID) (simulation.GeneralLawScenario, error)
}

// HistoricalStageStore is the persistence contract for Ch. 26 historical
// stages and their primitive-accumulation episodes.
type HistoricalStageStore interface {
	CreateHistoricalStage(ctx context.Context, h simulation.HistoricalStage) (simulation.HistoricalStage, error)
	GetHistoricalStage(ctx context.Context, id simulation.HistoricalStageID) (simulation.HistoricalStage, error)
	ListHistoricalStages(ctx context.Context) ([]simulation.HistoricalStage, error)
}

// EnclosureEventStore is the persistence contract for Ch. 27 enclosure events.
type EnclosureEventStore interface {
	CreateEnclosureEvent(ctx context.Context, e simulation.EnclosureEvent) (simulation.EnclosureEvent, error)
	ListEnclosureEvents(ctx context.Context) ([]simulation.EnclosureEvent, error)
}

// WageStatuteStore is the persistence contract for Ch. 28 wage statutes.
type WageStatuteStore interface {
	CreateWageStatute(ctx context.Context, w simulation.WageStatute) (simulation.WageStatute, error)
	ListWageStatutesByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.WageStatute, error)
}

// VagrancyLawStore is the persistence contract for Ch. 28 vagrancy laws.
type VagrancyLawStore interface {
	CreateVagrancyLaw(ctx context.Context, v simulation.VagrancyLaw) (simulation.VagrancyLaw, error)
	ListVagrancyLawsByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.VagrancyLaw, error)
}

// FarmTenureStore is the persistence contract for Ch. 29 farm tenure records.
type FarmTenureStore interface {
	CreateFarmTenure(ctx context.Context, f simulation.FarmTenure) (simulation.FarmTenure, error)
	ListFarmTenuresByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.FarmTenure, error)
}

// DomesticIndustryStore is the persistence contract for Ch. 30 domestic industries.
type DomesticIndustryStore interface {
	CreateDomesticIndustry(ctx context.Context, d simulation.DomesticIndustry) (simulation.DomesticIndustry, error)
	ListDomesticIndustriesByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.DomesticIndustry, error)
}

// CapitalOriginStore is the persistence contract for Ch. 31 capital origin records.
type CapitalOriginStore interface {
	CreateCapitalOrigin(ctx context.Context, c simulation.CapitalOrigin) (simulation.CapitalOrigin, error)
	ListCapitalOriginsByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.CapitalOrigin, error)
}

// ColonialTransferStore is the persistence contract for Ch. 31 colonial transfer records.
type ColonialTransferStore interface {
	CreateColonialTransfer(ctx context.Context, t simulation.ColonialTransfer) (simulation.ColonialTransfer, error)
	ListColonialTransfersByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.ColonialTransfer, error)
}

// NationalDebtStore is the persistence contract for Ch. 31 national debt records.
type NationalDebtStore interface {
	CreateNationalDebt(ctx context.Context, d simulation.NationalDebt) (simulation.NationalDebt, error)
	ListNationalDebtsByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.NationalDebt, error)
}

// ProtectionSystemStore is the persistence contract for Ch. 31 protection system records.
type ProtectionSystemStore interface {
	ListProtectionSystemsByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.ProtectionSystem, error)
}

// AccumulationTrajectoryStore is the persistence contract for Ch. 32
// long-run centralisation trajectories. A trajectory header carries the
// initial/final aggregates; its steps are stored alongside and are
// returned in step_index order.
type AccumulationTrajectoryStore interface {
	CreateAccumulationTrajectory(ctx context.Context, t simulation.AccumulationTrajectory) (simulation.AccumulationTrajectory, error)
	GetAccumulationTrajectory(ctx context.Context, id simulation.AccumulationTrajectoryID) (simulation.AccumulationTrajectory, error)
	ListAccumulationTrajectories(ctx context.Context) ([]simulation.AccumulationTrajectory, error)
}

// ColonialLabourMarketUpdate is the partial-update payload for
// regulating a Ch. 33 colonial labour market. Non-nil fields are
// applied; the rest are preserved.
type ColonialLabourMarketUpdate struct {
	WakefieldSchemeApplied   *bool
	IndependenceYears        *int64
	SurplusLabourExtractable *bool
}

// IsEmpty reports whether the update would mutate any field.
func (u ColonialLabourMarketUpdate) IsEmpty() bool {
	return u.WakefieldSchemeApplied == nil &&
		u.IndependenceYears == nil &&
		u.SurplusLabourExtractable == nil
}

// ColonialLabourMarketStore is the persistence contract for Ch. 33
// colonial labour markets. Colony names are unique (case-insensitive).
// Update applies a partial regulation payload in place.
// RegulateColonialLabourMarket reads the locked row, applies
// simulation.ColonialLabourRegulation against it, and writes the
// regulated state back in a single transaction — concurrent
// /regulate calls cannot race a stale baseline through to a clobber.
type ColonialLabourMarketStore interface {
	CreateColonialLabourMarket(ctx context.Context, m simulation.ColonialLabourMarket) (simulation.ColonialLabourMarket, error)
	GetColonialLabourMarket(ctx context.Context, id simulation.ColonialLabourMarketID) (simulation.ColonialLabourMarket, error)
	ListColonialLabourMarkets(ctx context.Context) ([]simulation.ColonialLabourMarket, error)
	UpdateColonialLabourMarket(ctx context.Context, id simulation.ColonialLabourMarketID, u ColonialLabourMarketUpdate) (simulation.ColonialLabourMarket, error)
	RegulateColonialLabourMarket(ctx context.Context, id simulation.ColonialLabourMarketID, scheme simulation.SystematicColonisation) (simulation.ColonialLabourMarket, error)
}

// ProductiveCircuitStore is the persistence contract for Vol. II Ch. 2
// productive-circuit records. The full state (latent capital, reserve fund,
// revenue exits, capitalisation steps, reserve draws) is assembled on Get.
type ProductiveCircuitStore interface {
	CreateProductiveCircuit(ctx context.Context, pc circulation.ProductiveCircuit) (circulation.ProductiveCircuit, error)
	GetProductiveCircuit(ctx context.Context, id circulation.ProductiveCircuitID) (circulation.ProductiveCircuit, error)
	ListProductiveCircuits(ctx context.Context) ([]circulation.ProductiveCircuit, error)
	RecordRevenue(ctx context.Context, id circulation.ProductiveCircuitID, rc circulation.RevenueCircuit) (circulation.RevenueCircuit, error)
	Accumulate(ctx context.Context, id circulation.ProductiveCircuitID, amount circulation.Pence) (circulation.LatentMoneyCapital, error)
	Capitalise(ctx context.Context, id circulation.ProductiveCircuitID, amount, dc, dv circulation.Pence) (circulation.CapitalisationStep, error)
	DepositReserve(ctx context.Context, id circulation.ProductiveCircuitID, amount circulation.Pence) (circulation.ReserveFund, error)
	WithdrawReserve(ctx context.Context, id circulation.ProductiveCircuitID, amount circulation.Pence, reason circulation.ReserveDrawReason) (circulation.ReserveDraw, error)
}

// CommodityCircuitStore is the persistence contract for Vol. II Ch. 3
// commodity-circuit records. The full state (partial sales, MP sources, terminal)
// is assembled on Get and included on List.
type CommodityCircuitStore interface {
	CreateCommodityCircuit(ctx context.Context, cc circulation.CommodityCircuit) (circulation.CommodityCircuit, error)
	GetCommodityCircuit(ctx context.Context, id circulation.CommodityCircuitID) (circulation.CommodityCircuit, error)
	ListCommodityCircuits(ctx context.Context, agentID string) ([]circulation.CommodityCircuit, error)
	RecordPartialSale(ctx context.Context, id circulation.CommodityCircuitID, sale circulation.SuccessivePartialSale) (circulation.SuccessivePartialSale, error)
	LinkMPSource(ctx context.Context, id circulation.CommodityCircuitID, source circulation.MeansOfProductionSource) (circulation.MeansOfProductionSource, error)
	CloseCommodityCircuit(ctx context.Context, id circulation.CommodityCircuitID, aug circulation.CommodityAugmented) (circulation.CommodityCircuit, error)
}

// MoneyCircuitStore is the persistence contract for Vol. II Ch. 1
// money-circuit records. Each phase (M—C, P, C′, C′—M′) is recorded
// sequentially; the store enforces the linear phase-transition invariant.
type MoneyCircuitStore interface {
	CreateMoneyCircuit(ctx context.Context, mc circulation.MoneyCircuit) (circulation.MoneyCircuit, error)
	GetMoneyCircuit(ctx context.Context, id circulation.MoneyCircuitID) (circulation.MoneyCircuit, error)
	ListMoneyCircuits(ctx context.Context, agentID string, moment circulation.CircuitMoment) ([]circulation.MoneyCircuit, error)
	RecordPurchase(ctx context.Context, id circulation.MoneyCircuitID, p circulation.PurchasePhase) (circulation.MoneyCircuit, error)
	RecordProductive(ctx context.Context, id circulation.MoneyCircuitID, ps circulation.ProductiveState) (circulation.MoneyCircuit, error)
	RecordCommodity(ctx context.Context, id circulation.MoneyCircuitID, cc circulation.CommodityCapital) (circulation.MoneyCircuit, error)
	RecordRealisation(ctx context.Context, id circulation.MoneyCircuitID, r circulation.Realisation) (circulation.MoneyCircuit, error)
}

// IndustrialCapitalStore is the persistence contract for Vol. II Ch. 4
// IndustrialCapital records and their satellite tables.
type IndustrialCapitalStore interface {
	CreateIndustrialCapital(ctx context.Context, ic circulation.IndustrialCapital) (circulation.IndustrialCapital, error)
	GetIndustrialCapital(ctx context.Context, id circulation.IndustrialCapitalID) (circulation.IndustrialCapital, error)
	ListIndustrialCapitals(ctx context.Context, agentID, status, economyMode string) ([]circulation.IndustrialCapital, error)
	RecordCapitalPart(ctx context.Context, id circulation.IndustrialCapitalID, part circulation.CapitalPart) (circulation.CapitalPart, error)
	Snapshot(ctx context.Context, id circulation.IndustrialCapitalID, sd circulation.StageDistribution) (circulation.StageDistribution, error)
	OpenBlock(ctx context.Context, id circulation.IndustrialCapitalID, b circulation.StageBlock) (circulation.StageBlock, error)
	CloseBlock(ctx context.Context, id circulation.IndustrialCapitalID, blockID circulation.StageBlockID) (circulation.StageBlock, error)
	RecordValueRevolution(ctx context.Context, res circulation.ValueRevolutionResult) (circulation.ValueRevolutionResult, error)
	RecordInterlock(ctx context.Context, id circulation.IndustrialCapitalID, mi circulation.MetamorphosisInterlock) (circulation.MetamorphosisInterlock, error)
	RecordSupplyDemand(ctx context.Context, sdi circulation.SupplyDemandImbalance) (circulation.SupplyDemandImbalance, error)
	GetSupplyDemand(ctx context.Context, id circulation.IndustrialCapitalID, period string) (circulation.SupplyDemandImbalance, error)
	AggregateSupplyDemand(ctx context.Context, period string) (circulation.AggregateSupplyDemandImbalance, error)
	SetSinkingFund(ctx context.Context, id circulation.IndustrialCapitalID, sf circulation.SinkingFund) (circulation.SinkingFund, error)
	TickSinkingFund(ctx context.Context, id circulation.IndustrialCapitalID) (circulation.SinkingFund, error)
}

// TurnoverStore is the persistence contract for Vol. II Ch. 7 —
// The Turnover Time and the Number of Turnovers.
type TurnoverStore interface {
	CreateTurnover(ctx context.Context, t tv.Turnover) (tv.Turnover, error)
	GetTurnover(ctx context.Context, id tv.TurnoverID) (tv.Turnover, error)
	ListTurnovers(ctx context.Context, industrialCapitalID string, lens tv.CircuitLens) ([]tv.Turnover, error)
	RecordCycle(ctx context.Context, id tv.TurnoverID, c tv.TurnoverCycle) (tv.TurnoverCycle, error)
	RecomputeNumber(ctx context.Context, id tv.TurnoverID) (tv.TurnoverNumber, error)
	GetNumber(ctx context.Context, id tv.TurnoverID) (tv.TurnoverNumber, error)
}

// AggregateTurnoverStore is the persistence contract for Vol. II Ch. 9
// aggregate-turnover records and their satellite tables (contributions,
// lifetime cycle, crisis-phase annotations).
type AggregateTurnoverStore interface {
	CreateAggregateTurnover(ctx context.Context, at tv.AggregateTurnover) (tv.AggregateTurnover, error)
	GetAggregateTurnover(ctx context.Context, id tv.AggregateTurnoverID) (tv.AggregateTurnover, error)
	ListAggregateTurnovers(ctx context.Context, industrialCapitalID string) ([]tv.AggregateTurnover, error)
	RecordContribution(ctx context.Context, id tv.AggregateTurnoverID, c tv.ComponentTurnoverContribution) (tv.ComponentTurnoverContribution, error)
	RecomputeAggregate(ctx context.Context, id tv.AggregateTurnoverID) (tv.AggregateTurnover, error)
	RecordLifetimeCycle(ctx context.Context, id tv.AggregateTurnoverID, lc tv.LifetimeCycle) (tv.LifetimeCycle, error)
	RecordCrisisPhase(ctx context.Context, cycleID tv.LifetimeCycleID, phase tv.CrisisCyclePhaseRecord) (tv.CrisisCyclePhaseRecord, error)
}

// EconomistAttributionStore is the persistence contract for Vol. II Ch. 10
// economist-attribution reference records (Quesnay and Smith). Records are
// static reference data; the store only exposes List with an optional theorist
// filter.
type EconomistAttributionStore interface {
	ListEconomistAttributions(ctx context.Context, theorist string) ([]circulation.EconomistAttribution, error)
}

// WorkingPeriodStore is the persistence contract for Vol. II Ch. 12 —
// The Working Period. CreateWorkingPeriod checks any NaturalWorkingPeriodConstraint
// for the commodity and returns wp.ErrPrematureDelivery if the count is too low.
type WorkingPeriodStore interface {
	CreateWorkingPeriod(ctx context.Context, w wp.WorkingPeriod) (wp.WorkingPeriod, error)
	GetWorkingPeriod(ctx context.Context, id wp.WorkingPeriodID) (wp.WorkingPeriod, error)
	ListWorkingPeriods(ctx context.Context, industrialCapitalID, mode, commodityID string) ([]wp.WorkingPeriod, error)
	RecordInterruption(ctx context.Context, id wp.WorkingPeriodID, intr wp.WorkingPeriodInterruption) (wp.WorkingPeriodInterruption, error)
	RecordShortening(ctx context.Context, id wp.WorkingPeriodID, s wp.WorkingPeriodShortening) (wp.WorkingPeriodShortening, error)
	RecordCreditFinancing(ctx context.Context, id wp.WorkingPeriodID, cf wp.WorkingPeriodCreditFinancing) (wp.WorkingPeriodCreditFinancing, error)
	CreateNaturalConstraint(ctx context.Context, nc wp.NaturalWorkingPeriodConstraint) (wp.NaturalWorkingPeriodConstraint, error)
	ListNaturalConstraints(ctx context.Context) ([]wp.NaturalWorkingPeriodConstraint, error)
}

// CompositionStore is the persistence contract for Vol. II Ch. 8 fixed and
// circulating capital composition records.
type CompositionStore interface {
	CreateComponent(ctx context.Context, c comp.CapitalComponent) (comp.CapitalComponent, error)
	GetComponent(ctx context.Context, id comp.CapitalComponentID) (comp.CapitalComponent, error)
	ListComponents(ctx context.Context, industrialCapitalID string, kind comp.CapitalKind, role comp.RoleInProcess) ([]comp.CapitalComponent, error)
	RegisterFixedItem(ctx context.Context, item comp.FixedCapitalItem) (comp.FixedCapitalItem, error)
	GetFixedItem(ctx context.Context, id comp.FixedCapitalItemID) (comp.FixedCapitalItem, error)
	RecordSubcomponent(ctx context.Context, sub comp.FixedCapitalSubcomponent) (comp.FixedCapitalSubcomponent, error)
	RecordWear(ctx context.Context, w comp.WearAndTear) (comp.WearAndTear, error)
	GetSinkingFund(ctx context.Context, itemID comp.FixedCapitalItemID) (comp.SinkingFundForItem, error)
	RecordRepair(ctx context.Context, r comp.Repair) (comp.Repair, error)
	RecordCirculatingCycle(ctx context.Context, c comp.CirculatingCycle) (comp.CirculatingCycle, error)
	RecordReinvestment(ctx context.Context, r comp.SinkingFundReinvestment) (comp.SinkingFundReinvestment, error)
}

// ProductionTimeStore is the persistence contract for Vol. II Ch. 13 —
// The Time of Production. Manages the production-labour gap, natural-process
// activations, natural subjects, and industry benchmarks tied to working periods.
type ProductionTimeStore interface {
	RecordProductionLabourGap(ctx context.Context, gap wp.ProductionLabourGap) (wp.ProductionLabourGap, error)
	GetProductionLabourGap(ctx context.Context, wpID wp.WorkingPeriodID) (wp.ProductionLabourGap, error)
	RecordNaturalProcessActivation(ctx context.Context, act wp.NaturalProcessActivation) (wp.NaturalProcessActivation, error)
	ListNaturalProcessActivations(ctx context.Context, wpID wp.WorkingPeriodID) ([]wp.NaturalProcessActivation, error)
	RecordNaturalSubject(ctx context.Context, sub wp.NaturalSubject) (wp.NaturalSubject, error)
	GetNaturalSubject(ctx context.Context, wpID wp.WorkingPeriodID) (wp.NaturalSubject, error)
	CreateIndustryBenchmark(ctx context.Context, ind wp.NaturalProcessIndustry) (wp.NaturalProcessIndustry, error)
	ListIndustryBenchmarks(ctx context.Context) ([]wp.NaturalProcessIndustry, error)
}

// PriceRevolutionStore is the persistence contract for Vol. II Ch. 15 —
// The Effects of a Change of Prices. It manages the derived satellite records
// for ValueRevolutionEvent: RevaluedCapital, RealisedValueDelta,
// InventoryRevaluation, SpeculativeHold, and CompoundCapitalAdjustment.
type PriceRevolutionStore interface {
	// CreatePriceRevolution persists a ValueRevolutionEvent together with its
	// derived RevaluedCapital (for MP/LP events) or RealisedValueDelta (for
	// output events). Returns the assembled record.
	CreatePriceRevolution(ctx context.Context, rec circulation.PriceRevolutionRecord) (circulation.PriceRevolutionRecord, error)
	// GetPriceRevolution returns a ValueRevolutionEvent with all derived rows.
	GetPriceRevolution(ctx context.Context, id circulation.ValueRevolutionEventID) (circulation.PriceRevolutionRecord, error)
	// RecordPriceCase attaches the operator-selected fall or rise case to an event.
	RecordPriceCase(ctx context.Context, pc circulation.PriceCaseRecord) (circulation.PriceCaseRecord, error)
	// RecordInventoryRevaluation persists a paper gain/loss on held stock.
	RecordInventoryRevaluation(ctx context.Context, inv circulation.InventoryRevaluation) (circulation.InventoryRevaluation, error)
	// RecordSpeculativeHold registers a commodity held speculatively.
	RecordSpeculativeHold(ctx context.Context, sh circulation.SpeculativeHold) (circulation.SpeculativeHold, error)
	// AggregateCompound computes the compound capital adjustment for one
	// industrial capital over one period, summing all delta streams.
	AggregateCompound(ctx context.Context, industrialCapitalID circulation.IndustrialCapitalID, period string) (circulation.CompoundCapitalAdjustment, error)
}

// AnnualSurplusRateStore is the persistence contract for Vol. II Ch. 16 —
// The Turnover of Variable Capital. It manages VariableCapitalAdvance records
// and their derived rates (per-circuit and annual) and the canonical contrast.
type AnnualSurplusRateStore interface {
	// CreateAdvance persists a new VariableCapitalAdvance after validation.
	CreateAdvance(ctx context.Context, a val.VariableCapitalAdvance) (val.VariableCapitalAdvance, error)
	// GetAdvance returns a VariableCapitalAdvance by ID.
	GetAdvance(ctx context.Context, id val.VariableCapitalAdvanceID) (val.VariableCapitalAdvance, error)
	// ListAdvances returns all advances for an industrial capital.
	// An empty industrialCapitalID returns all advances.
	ListAdvances(ctx context.Context, industrialCapitalID string) ([]val.VariableCapitalAdvance, error)
	// RecordPerCircuitRate persists the per-circuit surplus rate for an advance.
	// A second call replaces the prior rate (one rate per advance).
	RecordPerCircuitRate(ctx context.Context, rate val.PerCircuitSurplusRate) (val.PerCircuitSurplusRate, error)
	// GetPerCircuitRate returns the per-circuit surplus rate for an advance.
	GetPerCircuitRate(ctx context.Context, advanceID val.VariableCapitalAdvanceID) (val.PerCircuitSurplusRate, error)
	// RecordAnnualRate persists a computed AnnualSurplusRate snapshot.
	RecordAnnualRate(ctx context.Context, rate val.AnnualSurplusRate) (val.AnnualSurplusRate, error)
	// ListAnnualRates returns annual-rate snapshots, optionally filtered by
	// industrialCapitalID (via advance lookup) and/or period string.
	ListAnnualRates(ctx context.Context, industrialCapitalID, period string) ([]val.AnnualSurplusRate, error)
	// RecordContrast persists a canonical AnnualSurplusRateContrast.
	RecordContrast(ctx context.Context, c val.AnnualSurplusRateContrast) (val.AnnualSurplusRateContrast, error)
}

// SurplusCirculationStore is the persistence contract for Vol. II Ch. 17 —
// The Circulation of Surplus-Value. It manages SurplusCirculation records,
// their realisation-source mix rows, social-capital aggregate snapshots, and
// the canonical realisation puzzle.
type SurplusCirculationStore interface {
	// CreateSurplusCirculation persists a new SurplusCirculation after
	// validation. The RealisationSources slice is ignored on creation; use
	// RecordRealisationSource to append entries.
	CreateSurplusCirculation(ctx context.Context, sc repro.SurplusCirculation) (repro.SurplusCirculation, error)
	// GetSurplusCirculation returns a SurplusCirculation with all
	// RealisationSourceEntry rows assembled.
	GetSurplusCirculation(ctx context.Context, id repro.SurplusCirculationID) (repro.SurplusCirculation, error)
	// ListSurplusCirculations returns all SurplusCirculation records ordered by
	// creation time ascending.
	ListSurplusCirculations(ctx context.Context) ([]repro.SurplusCirculation, error)
	// RecordRealisationSource appends a RealisationSourceEntry to the given
	// SurplusCirculation, updating realised/unrealised pence atomically.
	RecordRealisationSource(ctx context.Context, id repro.SurplusCirculationID, entry repro.RealisationSourceEntry) (repro.RealisationSourceEntry, error)
	// SnapshotAggregate returns a SocialCapitalAggregate for the given period.
	// If no explicit aggregate row exists, it derives total_annual_surplus_pence
	// by summing all SurplusCirculations whose period matches.
	SnapshotAggregate(ctx context.Context, period string) (repro.SocialCapitalAggregate, error)
	// UpsertAggregateShares writes the Department I / II share split (basis
	// points, summing to 10000) onto the SocialCapitalAggregate for the given
	// period, creating the row if absent and leaving the other columns intact.
	// Ch. 20/21 reproduction ticks call this to populate the placeholders Ch. 17
	// reserved (issue #215).
	UpsertAggregateShares(ctx context.Context, period string, deptIShareBP, deptIIShareBP int64) (repro.SocialCapitalAggregate, error)
	// ListRealisationPuzzles returns all canonical RealisationPuzzle rows.
	ListRealisationPuzzles(ctx context.Context) ([]repro.RealisationPuzzle, error)
}

// MoneyCapitalStore is the persistence contract for Vol. II Ch. 18 —
// The Role of Money-Capital in Reproduction. It manages the five entity
// types that model how money-capital is apportioned, reserved, and exchanged
// between the two departments of social production.
type MoneyCapitalStore interface {
	CreateMoneySupplyApportionment(ctx context.Context, a repro.MoneySupplyApportionment) (repro.MoneySupplyApportionment, error)
	GetMoneySupplyApportionment(ctx context.Context, id repro.MoneySupplyApportionmentID) (repro.MoneySupplyApportionment, error)
	ListMoneySupplyApportionments(ctx context.Context) ([]repro.MoneySupplyApportionment, error)
	CreateDepartmentMoneyReserve(ctx context.Context, r repro.DepartmentMoneyReserve) (repro.DepartmentMoneyReserve, error)
	GetDepartmentMoneyReserve(ctx context.Context, id repro.DepartmentMoneyReserveID) (repro.DepartmentMoneyReserve, error)
	ListDepartmentMoneyReserves(ctx context.Context) ([]repro.DepartmentMoneyReserve, error)
	CreateCirculatingMoneyMass(ctx context.Context, m repro.CirculatingMoneyMass) (repro.CirculatingMoneyMass, error)
	GetCirculatingMoneyMass(ctx context.Context, id repro.CirculatingMoneyMassID) (repro.CirculatingMoneyMass, error)
	ListCirculatingMoneyMasses(ctx context.Context) ([]repro.CirculatingMoneyMass, error)
	CreateWageRotationFund(ctx context.Context, w repro.WageRotationFund) (repro.WageRotationFund, error)
	GetWageRotationFund(ctx context.Context, id repro.WageRotationFundID) (repro.WageRotationFund, error)
	ListWageRotationFunds(ctx context.Context) ([]repro.WageRotationFund, error)
	CreateInterDepartmentSettlement(ctx context.Context, s repro.InterDepartmentSettlement) (repro.InterDepartmentSettlement, error)
	GetInterDepartmentSettlement(ctx context.Context, id repro.InterDepartmentSettlementID) (repro.InterDepartmentSettlement, error)
	ListInterDepartmentSettlements(ctx context.Context) ([]repro.InterDepartmentSettlement, error)
	// CheckApportionmentBalance returns true when the four parts of the given
	// apportionment sum exactly to TotalSocialMoneyPence.
	CheckApportionmentBalance(ctx context.Context, id repro.MoneySupplyApportionmentID) (bool, error)
}

// SimpleReproductionSchemeStore is the persistence contract for Vol. II Ch. 20
// — Simple Reproduction. It manages the five entity types that model the
// two-department reproduction scheme: SimpleReproductionScheme (root),
// DepartmentalCapital (c+v+s composition per department), InterDepartmentExchange
// (the three value flows), ReproductionTick (cycle advancement), and
// MoneyClosedLoop (closure verification).
type SimpleReproductionSchemeStore interface {
	// CreateSimpleReproductionScheme persists a new scheme for a period.
	CreateSimpleReproductionScheme(ctx context.Context, s repro.SimpleReproductionScheme) (repro.SimpleReproductionScheme, error)
	// GetSimpleReproductionScheme returns the scheme with both departments
	// assembled (when present) and IsBalanced recomputed.
	GetSimpleReproductionScheme(ctx context.Context, id repro.SimpleReproductionSchemeID) (repro.SimpleReproductionScheme, error)
	// ListSimpleReproductionSchemes returns all schemes. When period is non-empty
	// the list is filtered to that period string.
	ListSimpleReproductionSchemes(ctx context.Context, period string) ([]repro.SimpleReproductionScheme, error)
	// AddDepartmentToScheme registers a DepartmentalCapital on the scheme and
	// recomputes IsBalanced. Returns ErrDepartmentAlreadySet if the department
	// slot is already occupied.
	AddDepartmentToScheme(ctx context.Context, id repro.SimpleReproductionSchemeID, d repro.DepartmentalCapital) (repro.DepartmentalCapital, error)
	// RecordInterDepartmentExchange persists one exchange flow linked to the scheme.
	RecordInterDepartmentExchange(ctx context.Context, id repro.SimpleReproductionSchemeID, e repro.InterDepartmentExchange) (repro.InterDepartmentExchange, error)
	// ListInterDepartmentExchanges returns all exchange flows for a scheme.
	ListInterDepartmentExchanges(ctx context.Context, schemeID repro.SimpleReproductionSchemeID) ([]repro.InterDepartmentExchange, error)
	// AdvanceReproductionTick increments TickCount on the scheme and persists a
	// ReproductionTick. workerPoolSize and subsistenceBasketPence drive the
	// labour-force reproduction assertion (issue #220); pass 0 to skip it.
	AdvanceReproductionTick(ctx context.Context, id repro.SimpleReproductionSchemeID, workerPoolSize, subsistenceBasketPence int64) (repro.ReproductionTick, error)
	// CheckSchemeBalance returns a BalanceCheckResult for the scheme, surfacing
	// the keystone-equation components for diagnostic purposes.
	CheckSchemeBalance(ctx context.Context, id repro.SimpleReproductionSchemeID) (repro.BalanceCheckResult, error)
	// RecordMoneyClosedLoop persists a MoneyClosedLoop record for the scheme.
	RecordMoneyClosedLoop(ctx context.Context, id repro.SimpleReproductionSchemeID, loop repro.MoneyClosedLoop) (repro.MoneyClosedLoop, error)
}

// ExtendedReproductionStore — Vol. II Ch. 21: Accumulation and Reproduction on
// an Extended Scale. It manages ExtendedReproductionScheme (root),
// Reinvestment (per-cycle surplus split), MultiPeriodScheme (cross-period
// analysis), CompositionShift (rising c/v ratio), ExtendedMoneyRequirement
// (money needed for growing circulation), and DepartmentIGrowthLead
// (confirmation that Dept I grows faster than Dept II).
type ExtendedReproductionStore interface {
	// CreateExtendedReproductionScheme persists a new scheme with accumulation rates.
	CreateExtendedReproductionScheme(ctx context.Context, s repro.ExtendedReproductionScheme) (repro.ExtendedReproductionScheme, error)
	// GetExtendedReproductionScheme returns the scheme with both departments
	// assembled (when present) and IsBalanced recomputed.
	GetExtendedReproductionScheme(ctx context.Context, id repro.ExtendedReproductionSchemeID) (repro.ExtendedReproductionScheme, error)
	// ListExtendedReproductionSchemes returns all schemes. When period is non-empty
	// only schemes whose period field matches are returned. Never returns nil.
	ListExtendedReproductionSchemes(ctx context.Context, period string) ([]repro.ExtendedReproductionScheme, error)
	// AddExtendedDepartmentToScheme attaches a DepartmentalCapital to the scheme.
	// Returns ErrNotFound if scheme missing. Returns ErrAlreadyExists if that
	// department slot is already occupied.
	AddExtendedDepartmentToScheme(ctx context.Context, schemeID repro.ExtendedReproductionSchemeID, dept repro.DepartmentalCapital) (repro.DepartmentalCapital, error)
	// CreateReinvestment persists a Reinvestment linked to the scheme.
	// Returns ErrNotFound if scheme missing.
	CreateReinvestment(ctx context.Context, schemeID repro.ExtendedReproductionSchemeID, r repro.Reinvestment) (repro.Reinvestment, error)
	// ListReinvestments returns all Reinvestment records for a scheme in
	// ascending cycle_number order. Returns an empty slice (never nil) if none.
	ListReinvestments(ctx context.Context, schemeID repro.ExtendedReproductionSchemeID) ([]repro.Reinvestment, error)
	// TickExtendedScheme advances one accumulation cycle: computes Reinvestment
	// for each department, advances DepartmentalCapital values, updates IsBalanced
	// and increments TickCount atomically.
	// Returns ErrNotFound if scheme or departments are missing.
	TickExtendedScheme(ctx context.Context, schemeID repro.ExtendedReproductionSchemeID) (repro.ExtendedReproductionScheme, error)
	// GetExtendedMoneyRequirement computes the money requirement for the scheme
	// without persisting it. Returns ErrNotFound if scheme missing.
	GetExtendedMoneyRequirement(ctx context.Context, schemeID repro.ExtendedReproductionSchemeID, baseMoneyPence int64) (repro.ExtendedMoneyRequirement, error)
	// CreateMultiPeriodScheme persists a MultiPeriodScheme and eagerly computes
	// CompositionShift and DepartmentIGrowthLead rows for each consecutive pair
	// of scheme IDs. Returns ErrNotFound if any linked scheme is missing.
	CreateMultiPeriodScheme(ctx context.Context, mp repro.MultiPeriodScheme) (repro.MultiPeriodScheme, error)
	// GetMultiPeriodScheme returns a MultiPeriodScheme by ID.
	GetMultiPeriodScheme(ctx context.Context, id repro.MultiPeriodSchemeID) (repro.MultiPeriodScheme, error)
	// ListCompositionShifts returns all CompositionShift rows for a
	// MultiPeriodScheme sorted by cycle_number ascending. Never nil.
	ListCompositionShifts(ctx context.Context, mpID repro.MultiPeriodSchemeID) ([]repro.CompositionShift, error)
	// ListGrowthLeads returns all DepartmentIGrowthLead rows for a
	// MultiPeriodScheme sorted by cycle_number ascending. Never nil.
	ListGrowthLeads(ctx context.Context, mpID repro.MultiPeriodSchemeID) ([]repro.DepartmentIGrowthLead, error)
}

// EngineTickStore is the persistence contract for the automatic tick
// scheduler's audit log (issue #214). Each RecordEngineTick row is one
// scheduler pass: how many registered tickers ran, how many domain entities
// they advanced, how long the pass took, and how many tickers errored.
// ListEngineTicks returns recent passes newest-first; a non-positive limit
// returns every row. The log is append-only — no update or delete.
type EngineTickStore interface {
	RecordEngineTick(ctx context.Context, t engine.ScheduledTick) (engine.ScheduledTick, error)
	ListEngineTicks(ctx context.Context, limit int) ([]engine.ScheduledTick, error)
}
