package engine_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
)

// Capital Vol. I, Ch. 15, §7:
// "the life of modern industry becomes a series of periods of moderate
//  activity, prosperity, over-production, crisis and stagnation."
func TestCyclePhase_Para7_FourCanonicalPhases(t *testing.T) {
	t.Parallel()
	for _, p := range []engine.CyclePhase{
		engine.CyclePhaseProsperity,
		engine.CyclePhaseOverproduction,
		engine.CyclePhaseCrisis,
		engine.CyclePhaseStagnation,
	} {
		if !p.IsValid() {
			t.Fatalf("CyclePhase %q must be valid", p)
		}
	}
	if engine.CyclePhase("boom").IsValid() {
		t.Fatal("non-canonical phase should be invalid")
	}
}
