package surplus_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/surplus"
)

// Scaled-integer magnitudes round half-up; they never truncate. Rate 150:300
// (Marx's 50%, Ch. 11) on V=3 gives 1.5, which truncation under-reported as 1.
func TestMassByRate_RoundsHalfUp(t *testing.T) {
	t.Parallel()
	half := surplus.SurplusValueRate{SurplusLabour: 150, NecessaryLabour: 300}
	third := surplus.SurplusValueRate{SurplusLabour: 100, NecessaryLabour: 300}
	twoThirds := surplus.SurplusValueRate{SurplusLabour: 200, NecessaryLabour: 300}

	for _, tc := range []struct {
		name string
		rate surplus.SurplusValueRate
		v    int64
		want surplus.SurplusValueMass
	}{
		{"1.5 rounds up to 2", half, 3, 2},
		{"0.33 rounds down to 0", third, 1, 0},
		{"0.67 rounds up to 1 (truncation gave 0)", twoThirds, 1, 1},
		{"exact 150 stays 150", half, 300, 150},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := surplus.MassByRate(tc.rate, surplus.VariableCapital(tc.v)); got != tc.want {
				t.Errorf("MassByRate = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestMassByWorkers_RoundsHalfUpOnceOverTheMass(t *testing.T) {
	t.Parallel()
	twoThirds := surplus.SurplusValueRate{SurplusLabour: 200, NecessaryLabour: 300}
	// 100 workers at v=1: 66.67 total, rounded once to 67 (truncation: 66).
	got := surplus.MassByWorkers(surplus.LabourPowerValue(1), twoThirds, surplus.WorkerCount(100))
	if got != 67 {
		t.Errorf("MassByWorkers = %d, want 67", got)
	}
}

func TestIndividualSurplus_RoundsHalfUp(t *testing.T) {
	t.Parallel()
	half := surplus.SurplusValueRate{SurplusLabour: 150, NecessaryLabour: 300}
	if got := surplus.IndividualSurplus(half, surplus.LabourPowerValue(3)); got != 2 {
		t.Errorf("IndividualSurplus = %d, want 2 (1.5 rounds up)", got)
	}
}

func TestMassByRate_NegativeNecessaryLabourIsZero(t *testing.T) {
	t.Parallel()
	bad := surplus.SurplusValueRate{SurplusLabour: 100, NecessaryLabour: -300}
	if got := surplus.MassByRate(bad, surplus.VariableCapital(3)); got != 0 {
		t.Errorf("MassByRate = %d, want 0 for non-positive necessary labour", got)
	}
}
