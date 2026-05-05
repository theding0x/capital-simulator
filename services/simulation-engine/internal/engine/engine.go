// Package engine holds the time-step orchestrator that advances the
// simulated economy one period at a time. The full tick scheduler is
// deferred to Ch. 10+; this package currently exposes domain types
// that later chapters will wire into the scheduler.
package engine

// ProductionRun records a ValorizationProcess result for a given simulation
// tick. LabourProcessID references a record in agent-service. Introduced in
// Ch. 7; the tick loop and persistence are added in Ch. 10+.
type ProductionRun struct {
	ID              string `json:"id"`
	TickID          string `json:"tick_id"`
	LabourProcessID string `json:"labour_process_id"`
	SurplusValue    int64  `json:"surplus_value"` // LabourMinutes
}
