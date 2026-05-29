package engine_test

import (
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
)

// Tick is the §1c simulation period. Construction is plain — the test
// guards against accidental field-renaming.
func TestTick_Construction(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 5, 12, 9, 0, 0, 0, time.UTC)
	tk := engine.Tick{
		FactoryID:        "needle-works",
		Sequence:         1,
		ValueTransferred: 2400,
		UnitsProduced:    580_000,
		HandLabourSaved:  0,
		OccurredAt:       now,
	}
	if tk.Sequence != 1 {
		t.Fatalf("Sequence = %d, want 1", tk.Sequence)
	}
	if !tk.OccurredAt.Equal(now) {
		t.Fatalf("OccurredAt = %v, want %v", tk.OccurredAt, now)
	}
}
