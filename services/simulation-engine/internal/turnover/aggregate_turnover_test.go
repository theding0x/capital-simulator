package turnover_test

import (
	"testing"

	comp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/composition"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/turnover"
)

// TestRecomputeAggregate_MarxCentralExample verifies the §3 canonical fixture:
// £80,000 fixed capital (10yr) + £20,000 circulating (5×/yr) → £108,000 annual
// turnover, aggregate n = 10800 bp (= 1 + 2/25 of £100,000 advanced).
func TestRecomputeAggregate_MarxCentralExample(t *testing.T) {
	t.Parallel()

	id := turnover.NewAggregateTurnoverID()
	fixed := turnover.ComponentTurnoverContribution{
		ID:                        turnover.NewComponentTurnoverContributionID(),
		AggregateTurnoverID:       id,
		Kind:                      comp.KindFixed,
		AdvancedPence:             8_000_000,
		TurnoverNumberBasisPoints: 1_000, // 1/10 per year (10-year life)
		AnnualReturnPence:         800_000,
		DifferenceKind:            turnover.DifferenceReal,
	}
	circulating := turnover.ComponentTurnoverContribution{
		ID:                        turnover.NewComponentTurnoverContributionID(),
		AggregateTurnoverID:       id,
		Kind:                      comp.KindCirculating,
		AdvancedPence:             2_000_000,
		TurnoverNumberBasisPoints: 50_000, // 5 times per year
		AnnualReturnPence:         10_000_000,
		DifferenceKind:            turnover.DifferenceReal,
	}

	at := turnover.AggregateTurnover{
		ID:                  id,
		IndustrialCapitalID: "mill-1867",
		Lens:                turnover.LensMoney,
		AdvancedTotalPence:  10_000_000,
	}
	result := turnover.RecomputeAggregate(at, []turnover.ComponentTurnoverContribution{fixed, circulating})

	if result.AdvancedTotalPence != 10_000_000 {
		t.Errorf("AdvancedTotalPence: want 10000000, got %d", result.AdvancedTotalPence)
	}
	// £800k + £10,000k = £10,800k
	if result.AnnualTurnedOverPence != 10_800_000 {
		t.Errorf("AnnualTurnedOverPence: want 10800000, got %d", result.AnnualTurnedOverPence)
	}
	// 10000 × 10800000 / 10000000 = 10800
	if result.AggregateNumberBasisPoints != 10_800 {
		t.Errorf("AggregateNumberBasisPoints: want 10800, got %d", result.AggregateNumberBasisPoints)
	}
}

// TestMeanTerm_ScropeThreeComponentExample verifies §5: $50,000 capital with three
// components producing $33,750 annually → mean term ≈ 778,667 min ≈ 18 months.
func TestMeanTerm_ScropeThreeComponentExample(t *testing.T) {
	t.Parallel()

	at := turnover.AggregateTurnover{
		Lens:                  turnover.LensMoney,
		AdvancedTotalPence:    5_000_000, // $50,000
		AnnualTurnedOverPence: 3_375_000, // $33,750
	}

	// MeanTerm = 5000000 × 525600 / 3375000 = 778666... (integer division)
	got := at.MeanTerm()
	want := int64(5_000_000) * turnover.YearReferenceMinutes / 3_375_000
	if got != want {
		t.Errorf("MeanTerm: want %d, got %d", want, got)
	}
	// Should be approximately 18 months (18 × 525600/12 = 788400 min) — within 2%.
	eighteenMonths := int64(18) * turnover.YearReferenceMinutes / 12
	diff := got - eighteenMonths
	if diff < 0 {
		diff = -diff
	}
	if diff > eighteenMonths/50 { // within 2%
		t.Errorf("MeanTerm %d is not within 2%% of 18 months (%d)", got, eighteenMonths)
	}
}

// TestMeanTermWithProfit_ShorterThanWithout verifies the Marx §5 observation that
// the profit-inclusive mean term is shorter than the profit-discounted term.
func TestMeanTermWithProfit_ShorterThanWithout(t *testing.T) {
	t.Parallel()

	at := turnover.AggregateTurnover{
		Lens:                  turnover.LensMoney,
		AdvancedTotalPence:    5_000_000,
		AnnualTurnedOverPence: 3_375_000,
	}
	surplusValue := turnover.Pence(500_000)

	withoutProfit := at.MeanTerm()
	withProfit := at.MeanTermWithProfit(surplusValue)

	if withProfit >= withoutProfit {
		t.Errorf("profit-inclusive mean term (%d) should be shorter than profit-discounted (%d)", withProfit, withoutProfit)
	}
}

// TestValidate_LensEnforcement verifies ErrLensNotHomogeneous is returned for
// non-money lenses. Only LensMoney is admissible for aggregate turnover.
func TestValidate_LensEnforcement(t *testing.T) {
	t.Parallel()

	cases := []struct {
		lens    turnover.CircuitLens
		wantErr bool
	}{
		{turnover.LensMoney, false},
		{turnover.LensProductive, true},
		{turnover.LensCommodity, true},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(string(tc.lens), func(t *testing.T) {
			t.Parallel()
			at := turnover.AggregateTurnover{
				Lens:               tc.lens,
				AdvancedTotalPence: 1_000_000,
			}
			err := at.Validate()
			if tc.wantErr && err == nil {
				t.Errorf("lens %q: want error, got nil", tc.lens)
			}
			if !tc.wantErr && err != nil {
				t.Errorf("lens %q: want no error, got %v", tc.lens, err)
			}
		})
	}
}

// TestComputeAnnualReturn verifies the per-component formula:
// AnnualReturn = AdvancedPence × TurnoverNumberBasisPoints / 10000.
func TestComputeAnnualReturn(t *testing.T) {
	t.Parallel()

	c := turnover.ComponentTurnoverContribution{
		AdvancedPence:             8_000_000,
		TurnoverNumberBasisPoints: 1_000, // 1/10 per year
	}
	if got := c.ComputeAnnualReturn(); got != 800_000 {
		t.Errorf("ComputeAnnualReturn: want 800000, got %d", got)
	}
}

// TestRecomputeAggregate_SingleCircuit verifies the §3 baseline: a single £4,000
// circulating capital turning over 5× yields AnnualTurnedOverPence = £20,000 and
// AggregateNumberBasisPoints = 50000.
func TestRecomputeAggregate_SingleCircuit(t *testing.T) {
	t.Parallel()

	id := turnover.NewAggregateTurnoverID()
	c := turnover.ComponentTurnoverContribution{
		AggregateTurnoverID:       id,
		Kind:                      comp.KindCirculating,
		AdvancedPence:             400_000, // £4,000
		TurnoverNumberBasisPoints: 50_000,  // 5 times per year
		AnnualReturnPence:         2_000_000,
		DifferenceKind:            turnover.DifferenceReal,
	}
	at := turnover.AggregateTurnover{ID: id, Lens: turnover.LensMoney, AdvancedTotalPence: 400_000}
	result := turnover.RecomputeAggregate(at, []turnover.ComponentTurnoverContribution{c})

	if result.AdvancedTotalPence != 400_000 {
		t.Errorf("AdvancedTotalPence: want 400000, got %d", result.AdvancedTotalPence)
	}
	if result.AnnualTurnedOverPence != 2_000_000 {
		t.Errorf("AnnualTurnedOverPence: want 2000000, got %d", result.AnnualTurnedOverPence)
	}
	if result.AggregateNumberBasisPoints != 50_000 {
		t.Errorf("AggregateNumberBasisPoints: want 50000, got %d", result.AggregateNumberBasisPoints)
	}
}

// TestRecomputeAggregate_AnnualTurnoverExceedsAdvance verifies the invariant from
// §3: "the capital-value turned over during the year may, on account of the
// repeated turnovers of the circulating capital within the same year, be larger
// than the aggregate value of the advanced capital."
func TestRecomputeAggregate_AnnualTurnoverExceedsAdvance(t *testing.T) {
	t.Parallel()

	id := turnover.NewAggregateTurnoverID()
	contributions := []turnover.ComponentTurnoverContribution{
		{
			AggregateTurnoverID:       id,
			Kind:                      comp.KindFixed,
			AdvancedPence:             8_000_000,
			TurnoverNumberBasisPoints: 1_000,
			AnnualReturnPence:         800_000,
		},
		{
			AggregateTurnoverID:       id,
			Kind:                      comp.KindCirculating,
			AdvancedPence:             2_000_000,
			TurnoverNumberBasisPoints: 50_000,
			AnnualReturnPence:         10_000_000,
		},
	}
	at := turnover.AggregateTurnover{ID: id, Lens: turnover.LensMoney}
	result := turnover.RecomputeAggregate(at, contributions)

	if result.AnnualTurnedOverPence <= result.AdvancedTotalPence {
		t.Errorf("annual turnover (%d) must exceed advanced total (%d) when circulating capital turns > 1×",
			result.AnnualTurnedOverPence, result.AdvancedTotalPence)
	}
}
