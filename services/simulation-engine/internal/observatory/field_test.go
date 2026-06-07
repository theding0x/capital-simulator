package observatory

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

func TestAdvanceFieldGrowsTotalByAlphaSurplus(t *testing.T) {
	t.Parallel()
	in := []circulation.FieldCapital{{
		ID:              "5eed0000000000000004ic1",
		TotalPence:      300000,
		MoneyPence:      100000,
		ProductionPence: 100000,
		CommodityPence:  100000,
		CostPricePence:  240000,
		SurplusPence:    60000,
	}}
	out := advanceField(in, 5000) // capitalise 50% of 60000 surplus = 30000

	if out[0].TotalPence != 330000 {
		t.Fatalf("TotalPence = %d, want 330000", out[0].TotalPence)
	}
	if sum := out[0].MoneyPence + out[0].ProductionPence + out[0].CommodityPence; sum != out[0].TotalPence {
		t.Fatalf("M+P+C = %d, want %d (parts must sum to total)", sum, out[0].TotalPence)
	}
	if out[0].SurplusPence != 60000 || out[0].CostPricePence != 240000 {
		t.Fatalf("surplus/cost mutated: %d/%d", out[0].SurplusPence, out[0].CostPricePence)
	}
	if in[0].TotalPence != 300000 {
		t.Fatalf("input mutated (not pure): TotalPence = %d, want 300000", in[0].TotalPence)
	}
}

func TestAdvanceFieldZeroSurplusIsNoOp(t *testing.T) {
	t.Parallel()
	in := []circulation.FieldCapital{{ID: "x", TotalPence: 1000, MoneyPence: 1000, SurplusPence: 0}}
	out := advanceField(in, 5000)
	if out[0].TotalPence != 1000 {
		t.Fatalf("TotalPence = %d, want 1000 (no surplus → no growth)", out[0].TotalPence)
	}
}

func TestAdvanceFieldZeroAlphaIsNoOp(t *testing.T) {
	t.Parallel()
	in := []circulation.FieldCapital{{ID: "x", TotalPence: 1000, MoneyPence: 1000, SurplusPence: 500}}
	out := advanceField(in, 0)
	if out[0].TotalPence != 1000 {
		t.Fatalf("TotalPence = %d, want 1000 (alpha 0 → no growth)", out[0].TotalPence)
	}
}
