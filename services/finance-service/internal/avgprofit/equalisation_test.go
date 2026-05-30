package avgprofit

import (
	"errors"
	"testing"
)

// TestComputeCapitalFlow_PriceBasedOutflow verifies the spec fixture:
// ComputeCapitalFlow("", 0, 0, 95, 100).Direction == KindOutflow
// because generalRate==0 so the price path fires: price 95 < value 100 → outflow.
func TestComputeCapitalFlow_PriceBasedOutflow(t *testing.T) {
	t.Parallel()
	cf := ComputeCapitalFlow("", 0, 0, 95, 100)

	if cf.Direction != KindOutflow {
		t.Errorf("Direction = %q, want %q (price 95 < value 100)", cf.Direction, KindOutflow)
	}
}

// TestComputeCapitalFlow_RateBasedInflow verifies the spec fixture:
// ComputeCapitalFlow("cotton", 3000, 2200, 0, 0).Direction == KindInflow
// because currentRate 3000 > generalRate 2200 → capital flows in.
func TestComputeCapitalFlow_RateBasedInflow(t *testing.T) {
	t.Parallel()
	cf := ComputeCapitalFlow("cotton", 3000, 2200, 0, 0)

	if cf.Direction != KindInflow {
		t.Errorf("Direction = %q, want %q (rate 3000 > general 2200)", cf.Direction, KindInflow)
	}
	if cf.SphereName != "cotton" {
		t.Errorf("SphereName = %q, want %q", cf.SphereName, "cotton")
	}
}

// TestComputeCapitalFlow_EqualRatesStable verifies the spec fixture:
// ComputeCapitalFlow("x", 2200, 2200, 0, 0).Direction == KindStable.
func TestComputeCapitalFlow_EqualRatesStable(t *testing.T) {
	t.Parallel()
	cf := ComputeCapitalFlow("x", 2200, 2200, 0, 0)

	if cf.Direction != KindStable {
		t.Errorf("Direction = %q, want %q (rates equal)", cf.Direction, KindStable)
	}
}

// TestComputeCapitalFlow_RateBasedOutflow verifies that a sphere rate below the general
// rate produces KindOutflow (rate path).
func TestComputeCapitalFlow_RateBasedOutflow(t *testing.T) {
	t.Parallel()
	cf := ComputeCapitalFlow("linen", 1500, 2200, 0, 0)

	if cf.Direction != KindOutflow {
		t.Errorf("Direction = %q, want %q (rate 1500 < general 2200)", cf.Direction, KindOutflow)
	}
}

// TestComputeCapitalFlow_PriceBasedInflow verifies price > value → KindInflow (price path).
func TestComputeCapitalFlow_PriceBasedInflow(t *testing.T) {
	t.Parallel()
	cf := ComputeCapitalFlow("coal", 0, 0, 110, 100)

	if cf.Direction != KindInflow {
		t.Errorf("Direction = %q, want %q (price 110 > value 100)", cf.Direction, KindInflow)
	}
}

// TestComputeCapitalFlow_AllZeros verifies that when all inputs are zero, KindStable is returned.
func TestComputeCapitalFlow_AllZeros(t *testing.T) {
	t.Parallel()
	cf := ComputeCapitalFlow("empty", 0, 0, 0, 0)

	if cf.Direction != KindStable {
		t.Errorf("Direction = %q, want %q (all zero)", cf.Direction, KindStable)
	}
}

// TestComputeCapitalFlow_RatePriorityOverPrice verifies that when generalRate != 0,
// the rate comparison wins over the market price/value.
func TestComputeCapitalFlow_RatePriorityOverPrice(t *testing.T) {
	t.Parallel()
	// generalRate != 0 → rate path: currentRate 1500 < generalRate 2200 → outflow,
	// even though price 110 > value 100 would signal inflow on the price path.
	cf := ComputeCapitalFlow("mixed", 1500, 2200, 110, 100)

	if cf.Direction != KindOutflow {
		t.Errorf("Direction = %q, want %q (rate path takes priority)", cf.Direction, KindOutflow)
	}
}

// TestComputeEqualisation_CottonConverging verifies the spec fixture:
// ComputeEqualisation("cotton", 3000, 2200): IsConverging true, Direction KindInflow,
// InitialRate 3000, TargetRate 2200.
func TestComputeEqualisation_CottonConverging(t *testing.T) {
	t.Parallel()
	eq := ComputeEqualisation("cotton", 3000, 2200)

	if !eq.IsConverging {
		t.Error("IsConverging = false, want true (3000 != 2200)")
	}
	if eq.Direction != KindInflow {
		t.Errorf("Direction = %q, want %q (3000 > 2200)", eq.Direction, KindInflow)
	}
	if eq.InitialRate != 3000 {
		t.Errorf("InitialRate = %d, want 3000", eq.InitialRate)
	}
	if eq.TargetRate != 2200 {
		t.Errorf("TargetRate = %d, want 2200", eq.TargetRate)
	}
	if eq.SphereName != "cotton" {
		t.Errorf("SphereName = %q, want %q", eq.SphereName, "cotton")
	}
}

// TestComputeEqualisation_EqualRatesBoundary verifies the boundary:
// equal rates ⇒ IsConverging false + KindStable.
func TestComputeEqualisation_EqualRatesBoundary(t *testing.T) {
	t.Parallel()
	eq := ComputeEqualisation("linen", 2200, 2200)

	if eq.IsConverging {
		t.Error("IsConverging = true, want false (equal rates)")
	}
	if eq.Direction != KindStable {
		t.Errorf("Direction = %q, want %q (equal rates)", eq.Direction, KindStable)
	}
}

// TestComputeEqualisation_BelowAverageOutflow verifies that a sphere with rate below
// general rate is converging outflow.
func TestComputeEqualisation_BelowAverageOutflow(t *testing.T) {
	t.Parallel()
	eq := ComputeEqualisation("iron", 1500, 2200)

	if !eq.IsConverging {
		t.Error("IsConverging = false, want true (1500 != 2200)")
	}
	if eq.Direction != KindOutflow {
		t.Errorf("Direction = %q, want %q (1500 < 2200)", eq.Direction, KindOutflow)
	}
}

// TestComputeEqualisation_FlowFieldPopulated verifies the embedded Flow matches
// the ComputeCapitalFlow result for the rate-only path.
func TestComputeEqualisation_FlowFieldPopulated(t *testing.T) {
	t.Parallel()
	eq := ComputeEqualisation("cotton", 3000, 2200)

	if eq.Flow.Direction != KindInflow {
		t.Errorf("Flow.Direction = %q, want %q", eq.Flow.Direction, KindInflow)
	}
	if eq.Flow.CurrentSphereRate != 3000 {
		t.Errorf("Flow.CurrentSphereRate = %d, want 3000", eq.Flow.CurrentSphereRate)
	}
	if eq.Flow.GeneralRate != 2200 {
		t.Errorf("Flow.GeneralRate = %d, want 2200", eq.Flow.GeneralRate)
	}
}

// TestEqualisation_Validate_NegativeRates verifies that negative InitialRate or
// TargetRate returns ErrNegativeMagnitude.
func TestEqualisation_Validate_NegativeRates(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name        string
		initialRate SphereProfitRate
		targetRate  SphereProfitRate
	}{
		{"negative initial rate", -1, 2200},
		{"negative target rate", 3000, -1},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			eq := Equalisation{InitialRate: tc.initialRate, TargetRate: tc.targetRate}
			if err := eq.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
				t.Errorf("Validate() = %v, want ErrNegativeMagnitude", err)
			}
		})
	}
}

// TestEqualisation_Validate_Valid verifies a correctly constructed record passes.
func TestEqualisation_Validate_Valid(t *testing.T) {
	t.Parallel()
	eq := ComputeEqualisation("cotton", 3000, 2200)
	if err := eq.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

// TestNewEqualisationID_Format verifies that NewEqualisationID returns a 24-character
// hex string (96-bit = 12 bytes).
func TestNewEqualisationID_Format(t *testing.T) {
	t.Parallel()
	id := NewEqualisationID()
	if id.IsZero() {
		t.Error("NewEqualisationID() returned a zero ID")
	}
	if len(string(id)) != 24 {
		t.Errorf("id length = %d, want 24 hex chars", len(string(id)))
	}
}
