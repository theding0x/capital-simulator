// Package store is the persistence layer for simulation-engine. As of Ch. 15
// it persists Machine and Factory records and the tick log produced by
// advancing a Factory one period at a time.
package store

import (
	"context"
	"errors"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
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
