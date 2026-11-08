package production_test

import (
	"math"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/production"
)

// §1 "the working day of 12 hours; the portion a b 10 hours necessary, b c 2 hours surplus"
func TestWorkingDayInvariant(t *testing.T) {
	t.Parallel()
	wd := production.WorkingDay{
		Total:           720,
		NecessaryLabour: 600,
		SurplusLabour:   120,
	}
	if production.LabourMinutes(wd.NecessaryLabour)+production.LabourMinutes(wd.SurplusLabour) != wd.Total {
		t.Fatalf("NecessaryLabour + SurplusLabour must equal Total: %d + %d ≠ %d",
			wd.NecessaryLabour, wd.SurplusLabour, wd.Total)
	}
}

// §1 "surplus-labour increases by one half, from 2 hours to 3 hours, although the working day remains at 12 hours"
// ShortenNecessaryLabour(wd, 540) → WorkingDay{Total:720, NecessaryLabour:540, SurplusLabour:180}
func TestShortenNecessaryLabour_Fixture(t *testing.T) {
	t.Parallel()
	wd := production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	got := production.ShortenNecessaryLabour(wd, production.LabourPowerValue(540))
	if got.Total != 720 {
		t.Fatalf("expected Total=720, got %d", got.Total)
	}
	if got.NecessaryLabour != 540 {
		t.Fatalf("expected NecessaryLabour=540, got %d", got.NecessaryLabour)
	}
	if got.SurplusLabour != 180 {
		t.Fatalf("expected SurplusLabour=180, got %d", got.SurplusLabour)
	}
}

// §1 invariant: ShortenNecessaryLabour increases SurplusLabour when newLPV < wd.NecessaryLabour
func TestShortenNecessaryLabour_SurplusIncreases(t *testing.T) {
	t.Parallel()
	wd := production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	got := production.ShortenNecessaryLabour(wd, production.LabourPowerValue(360))
	if got.SurplusLabour <= wd.SurplusLabour {
		t.Fatalf("surplus must increase: before=%d, after=%d", wd.SurplusLabour, got.SurplusLabour)
	}
}

// §1 invariant: result of ShortenNecessaryLabour still satisfies NL + SL = Total
func TestShortenNecessaryLabour_Invariant(t *testing.T) {
	t.Parallel()
	wd := production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	got := production.ShortenNecessaryLabour(wd, production.LabourPowerValue(360))
	if production.LabourMinutes(got.NecessaryLabour)+production.LabourMinutes(got.SurplusLabour) != got.Total {
		t.Fatalf("invariant broken after shorten: %d + %d ≠ %d",
			got.NecessaryLabour, got.SurplusLabour, got.Total)
	}
}

// §1 "value of labour-power reduced from five shillings to three, surplus-value increases from one to three"
// After reducing NecessaryLabour from 600 to 360: SurplusLabour = 720 - 360 = 360
func TestShortenNecessaryLabour_LPVReduction(t *testing.T) {
	t.Parallel()
	wd := production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	// new LabourPowerValue = 360 min (representing the reduced subsistence value in LabourMinutes)
	got := production.ShortenNecessaryLabour(wd, production.LabourPowerValue(360))
	if got.NecessaryLabour != 360 {
		t.Fatalf("expected NecessaryLabour=360, got %d", got.NecessaryLabour)
	}
	if got.SurplusLabour != 360 {
		t.Fatalf("expected SurplusLabour=360, got %d", got.SurplusLabour)
	}
}

// §1 RateOfSurplusValue: s/v = 120/600 = 0.2
func TestRateOfSurplusValue(t *testing.T) {
	t.Parallel()
	rate := production.RateOfSurplusValue(120, 600)
	if rate != 0.2 {
		t.Fatalf("expected 0.2, got %f", rate)
	}
}

// §1 invariant: RateOfSurplusValue strictly increases as ProductivityFactor rises (NL shrinks, total fixed)
func TestRateOfSurplusValue_IncreasesWithProductivity(t *testing.T) {
	t.Parallel()
	wd := production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	rateOld := production.RateOfSurplusValue(wd.SurplusLabour, wd.NecessaryLabour)
	wdNew := production.ShortenNecessaryLabour(wd, production.LabourPowerValue(360))
	rateNew := production.RateOfSurplusValue(wdNew.SurplusLabour, wdNew.NecessaryLabour)
	if rateNew <= rateOld {
		t.Fatalf("rate must increase with productivity rise: old=%f, new=%f", rateOld, rateNew)
	}
}

// §1 RateOfSurplusValue: zero necessary labour returns 0 (guard)
func TestRateOfSurplusValue_ZeroNL(t *testing.T) {
	t.Parallel()
	if got := production.RateOfSurplusValue(120, 0); got != 0 {
		t.Fatalf("expected 0 for zero NL, got %f", got)
	}
}

// §1 "double the productiveness ... individual value below social value"
// ExtraSurplusValue(IndividualValue(30), SocialValue(60), Quantity(24)) = (60-30)*24 = 720
func TestExtraSurplusValue_Fixture(t *testing.T) {
	t.Parallel()
	got := production.ExtraSurplusValue(
		production.IndividualValue(30),
		production.SocialValue(60),
		production.Quantity(24),
	)
	if got != 720 {
		t.Fatalf("expected 720, got %d", got)
	}
}

// §1 per-article extra = SocialValue − IndividualValue
func TestExtraSurplusValue_PerUnit(t *testing.T) {
	t.Parallel()
	got := production.ExtraSurplusValue(
		production.IndividualValue(30),
		production.SocialValue(60),
		production.Quantity(1),
	)
	if got != 30 {
		t.Fatalf("expected per-unit=30, got %d", got)
	}
}

// §1 invariant: ExtraSurplusValue == 0 when iv == sv
func TestExtraSurplusValue_ZeroWhenEqual(t *testing.T) {
	t.Parallel()
	got := production.ExtraSurplusValue(
		production.IndividualValue(60),
		production.SocialValue(60),
		production.Quantity(24),
	)
	if got != 0 {
		t.Fatalf("expected 0 when iv==sv, got %d", got)
	}
}

// §1 invariant: ExtraSurplusValue == 0 when iv > sv (no extra gain)
func TestExtraSurplusValue_ZeroWhenAbove(t *testing.T) {
	t.Parallel()
	got := production.ExtraSurplusValue(
		production.IndividualValue(70),
		production.SocialValue(60),
		production.Quantity(24),
	)
	if got != 0 {
		t.Fatalf("expected 0 when iv>sv, got %d", got)
	}
}

// value ∝ 1/productivity: doubling productivity halves SNLT
func TestApplyProductivityToSNLT_Double(t *testing.T) {
	t.Parallel()
	got := production.ApplyProductivityToSNLT(60, production.ProductivityFactor(2.0))
	if got != 30 {
		t.Fatalf("expected 30, got %d", got)
	}
}

// value ∝ 1/productivity: halving productivity doubles SNLT
func TestApplyProductivityToSNLT_Half(t *testing.T) {
	t.Parallel()
	got := production.ApplyProductivityToSNLT(60, production.ProductivityFactor(0.5))
	if got != 120 {
		t.Fatalf("expected 120, got %d", got)
	}
}

// guard: non-positive factor returns original snlt unchanged
func TestApplyProductivityToSNLT_ZeroFactor(t *testing.T) {
	t.Parallel()
	got := production.ApplyProductivityToSNLT(60, production.ProductivityFactor(0))
	if got != 60 {
		t.Fatalf("expected original 60 for zero factor, got %d", got)
	}
}

// Partitioned is the single L-P1 check. Marx's 12-hour day: 10 h necessary +
// 2 h surplus = 720 min. A day that loses or invents a minute fails it.
func TestWorkingDayPartitioned(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		wd   production.WorkingDay
		want bool
	}{
		{"12h day, 10h + 2h", production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}, true},
		{"all necessary, no surplus", production.WorkingDay{Total: 360, NecessaryLabour: 360, SurplusLabour: 0}, true},
		{"minute lost", production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 119}, false},
		{"minute invented", production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 121}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.wd.Partitioned(); got != tc.want {
				t.Errorf("Partitioned() = %v, want %v", got, tc.want)
			}
		})
	}
}

// §1 Marx's worked case, in minutes (1s = 120 min, 1d = 10 min). Social day:
// 12h = 720 min, necessary 10h = 600 (LPV 5s), surplus 2h = 120 (s = 1s), so
// s' = 5:1. A capitalist doubles productivity: 24 articles at an individual
// value of 9d (90 min) instead of 12 at 12d, sold at 10d (100 min), not at the
// social value. Extra surplus-value = 24 x (100 - 90) = 240 min = 2s, so
// s rises from 1s to 3s.
func TestWorkedCase_SoldAtTenpence(t *testing.T) {
	t.Parallel()
	const (
		sold     = production.SocialValue(100) // 10d, the price he actually gets
		indiv    = production.IndividualValue(90)
		articles = production.Quantity(24)
	)
	extra := production.ExtraSurplusValue(indiv, sold, articles)
	if extra != 240 {
		t.Fatalf("extra surplus-value = %d min, want 240 (2s)", extra)
	}
	const s0 = 120 // 1s: surplus on the average social conditions
	if got := s0 + int64(extra); got != 360 {
		t.Fatalf("surplus-value = %d min, want 360 (3s)", got)
	}
}

// §1 Same case: the necessary labour shortens from 10h to 7.5h, so the ratio of
// necessary to surplus labour moves from 5:1 to 5:3 (12h = 450 + 270 min).
func TestWorkedCase_NecessaryLabourShortensToSevenAndAHalfHours(t *testing.T) {
	t.Parallel()
	before := production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	if got := production.RateOfSurplusValue(before.SurplusLabour, before.NecessaryLabour); math.Abs(got-0.2) > 1e-9 {
		t.Fatalf("s' before = %v, want 1/5", got)
	}
	// New value of 8s (960 min) per day vs LPV 5s (600): individual
	// necessary labour = 720 x 600 / 960 = 450 min = 7.5h.
	after := production.ShortenNecessaryLabour(before, 450)
	if after.NecessaryLabour != 450 || after.SurplusLabour != 270 {
		t.Fatalf("day = %d + %d, want 450 + 270 (7.5h + 4.5h)", after.NecessaryLabour, after.SurplusLabour)
	}
	if !after.Partitioned() {
		t.Fatal("day must partition: necessary + surplus == total")
	}
	if got := production.RateOfSurplusValue(after.SurplusLabour, after.NecessaryLabour); math.Abs(got-0.6) > 1e-9 {
		t.Fatalf("s' after = %v, want 3/5", got)
	}
}
