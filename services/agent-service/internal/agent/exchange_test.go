package agent_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// §1: Exchange wine=50 for corn=50; both break even; total value unchanged.
func TestExchangeEquivalents_WineForCorn(t *testing.T) {
	t.Parallel()
	r := agent.ExchangeEquivalents(50, 50)
	if r.AAfter != 50 {
		t.Errorf("A after: want 50, got %d", r.AAfter)
	}
	if r.BAfter != 50 {
		t.Errorf("B after: want 50, got %d", r.BAfter)
	}
	if r.TotalBefore() != r.TotalAfter() {
		t.Errorf("total not conserved: before=%d after=%d", r.TotalBefore(), r.TotalAfter())
	}
	if r.AAfter-r.ABefore != 0 {
		t.Errorf("surplus for A: want 0, got %d", r.AAfter-r.ABefore)
	}
}

// §2: Seller of worth-100 sells at 110; gains 10, buyer loses 10, total 210.
func TestExchangeNonEquivalents_SellerAboveValue(t *testing.T) {
	t.Parallel()
	r := agent.ExchangeNonEquivalents(100, 110)
	sellerGain := r.AAfter - r.ABefore
	buyerGain := r.BAfter - r.BBefore
	if sellerGain != 10 {
		t.Errorf("seller gain: want 10, got %d", sellerGain)
	}
	if buyerGain != -10 {
		t.Errorf("buyer gain: want -10, got %d", buyerGain)
	}
	if r.TotalBefore() != r.TotalAfter() {
		t.Errorf("total not conserved: before=%d after=%d", r.TotalBefore(), r.TotalAfter())
	}
	if r.TotalBefore() != 210 {
		t.Errorf("total before: want 210, got %d", r.TotalBefore())
	}
}

// §3: A=40, B=50, total 90; after swap A=50, B=40, total still 90.
func TestExchangeNonEquivalents_WineForCorn(t *testing.T) {
	t.Parallel()
	r := agent.ExchangeNonEquivalents(40, 50)
	if r.AAfter != 50 {
		t.Errorf("A after: want 50, got %d", r.AAfter)
	}
	if r.BAfter != 40 {
		t.Errorf("B after: want 40, got %d", r.BAfter)
	}
	if r.TotalBefore() != 90 || r.TotalAfter() != 90 {
		t.Errorf("total not conserved: before=%d after=%d", r.TotalBefore(), r.TotalAfter())
	}
}

// §4: Property — for any (x, y), TotalBefore == TotalAfter.
func TestExchange_ValueConservation_Property(t *testing.T) {
	t.Parallel()
	cases := [][2]agent.Pence{
		{0, 0},
		{100, 100},
		{100, 200},
		{50, 150},
		{9999, 1},
	}
	for _, c := range cases {
		r := agent.ExchangeNonEquivalents(c[0], c[1])
		if r.TotalBefore() != r.TotalAfter() {
			t.Errorf("(%d,%d) non-equiv: total not conserved: before=%d after=%d",
				c[0], c[1], r.TotalBefore(), r.TotalAfter())
		}
		r2 := agent.ExchangeEquivalents(c[0], c[0])
		if r2.TotalBefore() != r2.TotalAfter() {
			t.Errorf("equiv(%d): total not conserved", c[0])
		}
	}
}

// §5: Scaling all values by k leaves SurplusValue==0 and ratios unchanged.
func TestMerchantsCapital_ScalingInvariant(t *testing.T) {
	t.Parallel()
	mc := agent.MerchantsCapital{M: 100, CommodityID: "cotton", MPrime: 100}
	if mc.SurplusValue() != 0 {
		t.Errorf("surplus: want 0, got %d", mc.SurplusValue())
	}
	scaled := agent.MerchantsCapital{M: 300, CommodityID: "cotton", MPrime: 300}
	if scaled.SurplusValue() != 0 {
		t.Errorf("scaled surplus: want 0, got %d", scaled.SurplusValue())
	}
}

// §6: M-M' (UsurersCapital) — SurplusValue=10; no commodity field to locate
// the source.
func TestUsurersCapital_SurplusValue(t *testing.T) {
	t.Parallel()
	uc := agent.UsurersCapital{M: 100, MPrime: 110}
	if uc.SurplusValue() != 10 {
		t.Errorf("surplus: want 10, got %d", uc.SurplusValue())
	}
}

// MerchantsCapital.Origin() returns "equivalent" / "redistribution".
func TestMerchantsCapital_Origin(t *testing.T) {
	t.Parallel()
	equiv := agent.MerchantsCapital{M: 100, CommodityID: "cotton", MPrime: 100}
	if equiv.Origin() != "equivalent" {
		t.Errorf("zero-surplus origin: want 'equivalent', got %q", equiv.Origin())
	}
	nonEquiv := agent.MerchantsCapital{M: 100, CommodityID: "cotton", MPrime: 110}
	if nonEquiv.Origin() != "redistribution" {
		t.Errorf("non-zero-surplus origin: want 'redistribution', got %q", nonEquiv.Origin())
	}
}

// TotalValue sums MoneyBalance across agents; invariant across exchange.
func TestTotalValue_Conservation(t *testing.T) {
	t.Parallel()
	a := agent.Agent{MoneyBalance: 100}
	b := agent.Agent{MoneyBalance: 200}
	before := agent.TotalValue([]agent.Agent{a, b})
	if before != 300 {
		t.Errorf("want 300, got %d", before)
	}
	a.MoneyBalance = 200
	b.MoneyBalance = 100
	after := agent.TotalValue([]agent.Agent{a, b})
	if after != before {
		t.Errorf("TotalValue not conserved: before=%d after=%d", before, after)
	}
}

// ExchangeNonEquivalents: zero-sum property — A's gain + B's gain == 0.
func TestExchangeNonEquivalents_ZeroSum(t *testing.T) {
	t.Parallel()
	r := agent.ExchangeNonEquivalents(100, 110)
	aGain := r.AAfter - r.ABefore
	bGain := r.BAfter - r.BBefore
	if aGain+bGain != 0 {
		t.Errorf("zero-sum violated: aGain=%d bGain=%d", aGain, bGain)
	}
}
