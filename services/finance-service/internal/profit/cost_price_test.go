package profit

import (
	"errors"
	"testing"
)

// Capital III, Ch. 1, §1: a capital outlay of 400c + 100v, with the rate of
// surplus-value at 100%, yields a product worth 400c + 100v + 100s = £600. The
// cost-price is k = c + v = £500.
func TestComputeCostPrice_MarxSection1(t *testing.T) {
	t.Parallel()
	cp := ComputeCostPrice(CapitalOutlay{Constant: 400, Variable: 100})
	if cp.K != 500 {
		t.Fatalf("k = %d, want 500", cp.K)
	}
	if cp.CirculatingComponent != 500 {
		t.Errorf("circulating component = %d, want 500", cp.CirculatingComponent)
	}
	if cp.FixedComponent != 0 {
		t.Errorf("fixed component = %d, want 0 (no fixed capital declared)", cp.FixedComponent)
	}
}

// §1: of a £1,200 instrument of labour, only its £20 wear-and-tear enters the
// cost-price; alongside £380 circulating constant and £100 variable, the
// cost-price components are fixed £20 and circulating £480, k = £500.
func TestComputeCostPrice_FixedCirculatingSplit(t *testing.T) {
	t.Parallel()
	cp := ComputeCostPrice(CapitalOutlay{
		Constant:         400, // 20 fixed wear + 380 circulating constant
		Variable:         100,
		FixedWearAndTear: 20,
		FixedAdvanced:    1200,
	})
	if cp.FixedComponent != 20 {
		t.Errorf("fixed component = %d, want 20", cp.FixedComponent)
	}
	if cp.CirculatingComponent != 480 {
		t.Errorf("circulating component = %d, want 480", cp.CirculatingComponent)
	}
	if cp.K != 500 {
		t.Errorf("k = %d, want 500", cp.K)
	}
}

// Invariant: CostPrice.K == ConstantCapital + VariableCapital, and the fixed
// and circulating components always sum to K.
func TestComputeCostPrice_Invariants(t *testing.T) {
	t.Parallel()
	cases := []CapitalOutlay{
		{Constant: 0, Variable: 0},
		{Constant: 1, Variable: 0},
		{Constant: 400, Variable: 100, FixedWearAndTear: 20, FixedAdvanced: 1200},
		{Constant: 999, Variable: 333, FixedWearAndTear: 100, FixedAdvanced: 5000},
	}
	for _, o := range cases {
		cp := ComputeCostPrice(o)
		if want := LabourMinutes(o.Constant) + LabourMinutes(o.Variable); cp.K != want {
			t.Errorf("outlay %+v: k = %d, want c+v = %d", o, cp.K, want)
		}
		if cp.FixedComponent+cp.CirculatingComponent != cp.K {
			t.Errorf("outlay %+v: components %d + %d != k %d",
				o, cp.FixedComponent, cp.CirculatingComponent, cp.K)
		}
	}
}

func TestCapitalOutlay_Validate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		outlay  CapitalOutlay
		wantErr error
	}{
		{"valid simple", CapitalOutlay{Constant: 400, Variable: 100}, nil},
		{"valid split", CapitalOutlay{Constant: 400, Variable: 100, FixedWearAndTear: 20, FixedAdvanced: 1200}, nil},
		{"negative constant", CapitalOutlay{Constant: -1, Variable: 100}, ErrNegativeMagnitude},
		{"negative variable", CapitalOutlay{Constant: 400, Variable: -5}, ErrNegativeMagnitude},
		{"wear exceeds constant", CapitalOutlay{Constant: 10, Variable: 100, FixedWearAndTear: 20}, ErrFixedWearExceedsConstant},
		{"wear exceeds advanced", CapitalOutlay{Constant: 400, Variable: 100, FixedWearAndTear: 30, FixedAdvanced: 20}, ErrWearExceedsAdvanced},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.outlay.Validate(); !errors.Is(err, tc.wantErr) {
				t.Errorf("Validate() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

// Invariant: CommodityValueV3.Total == CostPrice.K + SurplusValue, and the
// cost-price is always strictly below the value while s > 0.
func TestComputeCommodityValue(t *testing.T) {
	t.Parallel()
	c := ComputeCommodityValue(500, 100)
	if c.Total != 600 {
		t.Errorf("total = %d, want 600", c.Total)
	}
	if c.Total != c.CostPrice+LabourMinutes(c.SurplusValue) {
		t.Errorf("total %d != k %d + s %d", c.Total, c.CostPrice, c.SurplusValue)
	}
	if c.CostPrice >= c.Total {
		t.Errorf("cost-price %d should be < total %d when s > 0", c.CostPrice, c.Total)
	}
}

func TestNewCostPriceID_UniqueAnd24Hex(t *testing.T) {
	t.Parallel()
	seen := make(map[CostPriceID]bool)
	for i := 0; i < 1000; i++ {
		id := NewCostPriceID()
		if len(string(id)) != 24 {
			t.Fatalf("id %q length = %d, want 24 hex chars", id, len(string(id)))
		}
		if seen[id] {
			t.Fatalf("duplicate id %q", id)
		}
		seen[id] = true
	}
}
