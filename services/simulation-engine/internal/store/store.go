// Package store is the persistence layer for simulation-engine. As of Ch. 15
// it persists Machine and Factory records and the tick log produced by
// advancing a Factory one period at a time. Ch. 25 adds GeneralLawScenario.
package store

import (
	"context"
	"errors"

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
