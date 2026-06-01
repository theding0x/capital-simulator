package valorisation_test

import (
	"testing"

	val "github.com/theding0x/capital-simulator/services/simulation-engine/internal/valorisation"
)

// TestComputeAdjustedAnnualRate covers the Ch.15→16 bridge (issue #222): a price
// revolution on output re-prices the value added (v+s); the whole swing accrues
// to surplus against the original advance, then re-annualises by n. Base
// fixture: v=1000, s=1000 (per-circuit 100%), n=1× → annual 100% (10000 bp).
func TestComputeAdjustedAnnualRate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name               string
		priceBP            int64
		wantAdjustedSurp   int64
		wantAdjustedAnnual int64
	}{
		{"no change", 0, 1000, 10000},
		{"output price +10%", 1000, 1200, 12000}, // +10% × (v+s=2000) = +200 → s=1200
		{"output price -10%", -1000, 800, 8000},  // -200 → s=800
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := val.ComputeAdjustedAnnualRate(1000, 1000, 10000, tc.priceBP)
			if int64(got.AdjustedSurplusPence) != tc.wantAdjustedSurp {
				t.Fatalf("adjusted surplus: got %d, want %d", got.AdjustedSurplusPence, tc.wantAdjustedSurp)
			}
			if got.AdjustedAnnualBasisPoints != tc.wantAdjustedAnnual {
				t.Fatalf("adjusted annual bp: got %d, want %d", got.AdjustedAnnualBasisPoints, tc.wantAdjustedAnnual)
			}
			if got.BaseAnnualBasisPoints != 10000 {
				t.Fatalf("base annual bp: got %d, want 10000", got.BaseAnnualBasisPoints)
			}
		})
	}
}

// TestComputeAdjustedAnnualRate_TurnoverMultiplies: doubling n doubles the
// annual rate (Marx's n × s/v).
func TestComputeAdjustedAnnualRate_TurnoverMultiplies(t *testing.T) {
	t.Parallel()
	got := val.ComputeAdjustedAnnualRate(1000, 1000, 20000, 0) // n = 2×
	if got.AdjustedAnnualBasisPoints != 20000 {
		t.Fatalf("annual bp at n=2: got %d, want 20000", got.AdjustedAnnualBasisPoints)
	}
}

// TestComputeAdjustedAnnualRate_ZeroVariable yields zero rates.
func TestComputeAdjustedAnnualRate_ZeroVariable(t *testing.T) {
	t.Parallel()
	got := val.ComputeAdjustedAnnualRate(0, 0, 10000, 1000)
	if got.BasePerCircuitBasisPoints != 0 || got.AdjustedAnnualBasisPoints != 0 {
		t.Fatalf("zero variable should yield zero rates, got base=%d annual=%d",
			got.BasePerCircuitBasisPoints, got.AdjustedAnnualBasisPoints)
	}
}
