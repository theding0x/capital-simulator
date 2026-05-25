// Package store is the persistence layer for simulation-engine. As of Ch. 15
// it persists Machine and Factory records and the tick log produced by
// advancing a Factory one period at a time. Ch. 25 adds GeneralLawScenario.
// Vol. II Ch. 2 adds ProductiveCircuit.
package store

import (
	"context"
	"errors"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
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
