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
// output by one working day.
type Tick struct {
	ID        string    `json:"id"`
	Sequence  int64     `json:"sequence"`
	OccurredAt time.Time `json:"occurred_at"`
}

// CyclePhase is the industrial-cycle enum from Ch. 15, §7: "the life of
// modern industry becomes a series of periods of moderate activity,
// prosperity, over-production, crisis and stagnation."
type CyclePhase string

const (
	CyclePhaseProsperity     CyclePhase = "prosperity"
	CyclePhaseOverproduction CyclePhase = "overproduction"
	CyclePhaseCrisis         CyclePhase = "crisis"
	CyclePhaseStagnation     CyclePhase = "stagnation"
)

// IsValid reports whether the phase value is one of the four canonical
// cycle phases.
func (p CyclePhase) IsValid() bool {
	switch p {
	case CyclePhaseProsperity, CyclePhaseOverproduction, CyclePhaseCrisis, CyclePhaseStagnation:
		return true
	}
	return false
}
