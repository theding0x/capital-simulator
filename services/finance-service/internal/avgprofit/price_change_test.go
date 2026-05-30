package avgprofit

import (
	"errors"
	"testing"
)

// TestComputePriceOfProductionChange_GeneralRateChange: when only the general
// rate of profit moves, RateChanged fires, ValueChanged does not, and the price
// is reported as changed.
func TestComputePriceOfProductionChange_GeneralRateChange(t *testing.T) {
	t.Parallel()
	c := ComputePriceOfProductionChange("cotton", CauseGeneralRateChange, 122, 130)
	if !c.PriceChanged {
		t.Error("PriceChanged = false, want true")
	}
	if !c.RateChanged {
		t.Error("RateChanged = false, want true")
	}
	if c.ValueChanged {
		t.Error("ValueChanged = true, want false")
	}
	if c.SphereName != "cotton" {
		t.Errorf("SphereName = %q, want cotton", c.SphereName)
	}
	if c.OldPrice != 122 || c.NewPrice != 130 {
		t.Errorf("prices = old %d new %d, want 122/130", c.OldPrice, c.NewPrice)
	}
}

// TestComputePriceOfProductionChange_IndividualValueChange: when only the
// commodity's own value moves, ValueChanged fires, RateChanged does not, and the
// price is reported as changed.
func TestComputePriceOfProductionChange_IndividualValueChange(t *testing.T) {
	t.Parallel()
	c := ComputePriceOfProductionChange("linen", CauseIndividualValueChange, 122, 110)
	if !c.PriceChanged {
		t.Error("PriceChanged = false, want true")
	}
	if !c.ValueChanged {
		t.Error("ValueChanged = false, want true")
	}
	if c.RateChanged {
		t.Error("RateChanged = true, want false")
	}
}

// TestComputePriceOfProductionChange_Both: when both causes are at work, both
// flags fire and the price is reported as changed.
func TestComputePriceOfProductionChange_Both(t *testing.T) {
	t.Parallel()
	c := ComputePriceOfProductionChange("iron", CauseBoth, 100, 150)
	if !c.RateChanged {
		t.Error("RateChanged = false, want true")
	}
	if !c.ValueChanged {
		t.Error("ValueChanged = false, want true")
	}
	if !c.PriceChanged {
		t.Error("PriceChanged = false, want true")
	}
}

// TestComputePriceOfProductionChange_NoCause: an empty cause means nothing
// changed — no flags fire, including PriceChanged.
func TestComputePriceOfProductionChange_NoCause(t *testing.T) {
	t.Parallel()
	c := ComputePriceOfProductionChange("wheat", "", 100, 100)
	if c.PriceChanged {
		t.Error("PriceChanged = true, want false for empty cause")
	}
	if c.RateChanged {
		t.Error("RateChanged = true, want false for empty cause")
	}
	if c.ValueChanged {
		t.Error("ValueChanged = true, want false for empty cause")
	}
}

// TestPriceChanged_EqualsCauseNonEmpty asserts the invariant PriceChanged ==
// (Cause != "") across every cause value.
func TestPriceChanged_EqualsCauseNonEmpty(t *testing.T) {
	t.Parallel()
	causes := []PriceChangeCause{
		CauseGeneralRateChange,
		CauseIndividualValueChange,
		CauseBoth,
		"",
	}
	for _, cause := range causes {
		c := ComputePriceOfProductionChange("x", cause, 0, 0)
		want := cause != ""
		if c.PriceChanged != want {
			t.Errorf("cause %q: PriceChanged = %v, want %v", cause, c.PriceChanged, want)
		}
	}
}

// TestPriceOfProductionChange_Validate exercises the negative-magnitude guard.
func TestPriceOfProductionChange_Validate(t *testing.T) {
	t.Parallel()

	if err := (PriceOfProductionChange{OldPrice: 122, NewPrice: 130}).Validate(); err != nil {
		t.Errorf("valid change: err = %v, want nil", err)
	}
	if err := (PriceOfProductionChange{OldPrice: -1, NewPrice: 130}).Validate(); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("negative OldPrice: err = %v, want ErrNegativeMagnitude", err)
	}
	if err := (PriceOfProductionChange{OldPrice: 122, NewPrice: -5}).Validate(); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("negative NewPrice: err = %v, want ErrNegativeMagnitude", err)
	}
}

// TestComputeCompensationGround asserts Ch. 12's second remark: a claimed ground
// for a higher price is never value creation, only a share in the aggregate
// surplus-value.
func TestComputeCompensationGround(t *testing.T) {
	t.Parallel()
	g := ComputeCompensationGround("slow turnover", "extra value creation")
	if g.ActuallyIs != "share in aggregate surplus-value" {
		t.Errorf("ActuallyIs = %q, want %q", g.ActuallyIs, "share in aggregate surplus-value")
	}
	if g.CreatesValue {
		t.Error("CreatesValue = true, want false")
	}
	if g.Reason != "slow turnover" {
		t.Errorf("Reason = %q, want %q", g.Reason, "slow turnover")
	}
	if g.AppearsToBe != "extra value creation" {
		t.Errorf("AppearsToBe = %q, want %q", g.AppearsToBe, "extra value creation")
	}
}

// TestComputeCompensationGround_InvariantForAnyReason asserts that whatever the
// reason or appearance, ActuallyIs is the constant and CreatesValue is false.
func TestComputeCompensationGround_InvariantForAnyReason(t *testing.T) {
	t.Parallel()
	cases := []struct {
		reason  string
		appears string
	}{
		{"risk premium", "reward for bearing risk"},
		{"superior skill", "fruit of the entrepreneur's talent"},
		{"", ""},
	}
	for _, tc := range cases {
		g := ComputeCompensationGround(tc.reason, tc.appears)
		if g.ActuallyIs != "share in aggregate surplus-value" {
			t.Errorf("reason %q: ActuallyIs = %q, want constant", tc.reason, g.ActuallyIs)
		}
		if g.CreatesValue {
			t.Errorf("reason %q: CreatesValue = true, want false", tc.reason)
		}
	}
}
