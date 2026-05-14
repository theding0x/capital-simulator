// Package engine holds the time-step orchestrator that advances the
// simulated economy one period at a time. Domain types here are wired into
// the tick scheduler in later chapters.
package engine

import "time"

// ProductionRun records a ValorizationProcess result for a given simulation
// tick. LabourProcessID references a record in agent-service. Introduced in
// Ch. 7; the tick loop and persistence are added in Ch. 10+.
type ProductionRun struct {
	ID              string `json:"id"`
	TickID          string `json:"tick_id"`
	LabourProcessID string `json:"labour_process_id"`
	SurplusValue    int64  `json:"surplus_value"` // LabourMinutes
}

// Tick is one simulation time-step. Introduced in Ch. 15 to drive the
// machinery / factory loop: each Tick advances wear, value transfer, and
// output by one working day. FactoryID names the factory that produced
// the period; Sequence is monotonically increasing within that factory;
// OccurredAt is server clock at the moment the tick was committed.
type Tick struct {
	FactoryID        string    `json:"factory_id"`
	Sequence         int64     `json:"sequence"`
	ValueTransferred int64     `json:"value_transferred"`
	UnitsProduced    int64     `json:"units_produced"`
	HandLabourSaved  int64     `json:"hand_labour_saved"`
	OccurredAt       time.Time `json:"occurred_at"`
}
