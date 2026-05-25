package turnover_test

import (
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/turnover"
)

// YearReferenceMinutes constant must be fixed at 525600 (365 × 24 × 60).
func TestYearReferenceMinutes(t *testing.T) {
	t.Parallel()
	if turnover.YearReferenceMinutes != 525_600 {
		t.Errorf("YearReferenceMinutes = %d, want 525600", turnover.YearReferenceMinutes)
	}
}

// --- ValidateLens -----------------------------------------------------------

// LensCommodity (Form III) is rejected for turnover analysis.
// Marx: "But it is not to be used in connection with the turnover of capital,
// which always begins with the advance of capital-value."
func TestValidateLens_CommodityRejected(t *testing.T) {
	t.Parallel()
	err := turnover.ValidateLens(turnover.LensCommodity)
	if err == nil {
		t.Fatal("expected ErrLensNotApplicableForTurnover, got nil")
	}
	if err != turnover.ErrLensNotApplicableForTurnover {
		t.Errorf("err = %v, want ErrLensNotApplicableForTurnover", err)
	}
}

func TestValidateLens_MoneyAccepted(t *testing.T) {
	t.Parallel()
	if err := turnover.ValidateLens(turnover.LensMoney); err != nil {
		t.Errorf("LensMoney rejected: %v", err)
	}
}

func TestValidateLens_ProductiveAccepted(t *testing.T) {
	t.Parallel()
	if err := turnover.ValidateLens(turnover.LensProductive); err != nil {
		t.Errorf("LensProductive rejected: %v", err)
	}
}

func TestValidateLens_UnknownRejected(t *testing.T) {
	t.Parallel()
	if err := turnover.ValidateLens("fictional"); err == nil {
		t.Error("expected error for unknown lens, got nil")
	}
}

// --- ComputeTurnoverNumber --------------------------------------------------

// "If we designate the year as the unit of measure of the turnover time by T,
// the time of turnover of a given capital by t, and the number of its turnovers
// by n, then n = T/t."
// BasisPoints invariant: BasisPoints == 10000 * YearReferenceMinutes / t
func TestComputeTurnoverNumber_FormulaIdentity(t *testing.T) {
	t.Parallel()
	cases := []int64{10080, 129600, 777600, 2102400}
	for _, tt := range cases {
		tn := turnover.ComputeTurnoverNumber("id-test", tt)
		want := 10_000 * turnover.YearReferenceMinutes / tt
		if tn.BasisPoints != want {
			t.Errorf("t=%d: BasisPoints=%d, want %d", tt, tn.BasisPoints, want)
		}
		if tn.TurnoverTimeMinutes != tt {
			t.Errorf("t=%d: TurnoverTimeMinutes=%d, want %d", tt, tn.TurnoverTimeMinutes, tt)
		}
		if tn.YearReferenceMinutes != turnover.YearReferenceMinutes {
			t.Errorf("YearReferenceMinutes mutable")
		}
	}
}

// SpinningMill: weekly turnover ≈ 52 times/year.
// t = 7 × 24 × 60 = 10080 minutes. BasisPoints ≈ 521428 (52.14 × 10000).
func TestComputeTurnoverNumber_SpinningMill_Weekly(t *testing.T) {
	t.Parallel()
	const weekMinutes int64 = 7 * 24 * 60 // 10080
	tn := turnover.ComputeTurnoverNumber("spinning-mill", weekMinutes)
	want := int64(10_000 * 525_600 / 10_080) // 521428
	if tn.BasisPoints != want {
		t.Errorf("SpinningMill BasisPoints = %d, want %d (≈52.14 turnovers/yr)", tn.BasisPoints, want)
	}
	// Sub-year turnovers: BasisPoints > 10000
	if tn.BasisPoints <= 10_000 {
		t.Errorf("weekly turnover should have BasisPoints > 10000 (sub-year), got %d", tn.BasisPoints)
	}
}

// LocomotiveMaker: 18-month turnover n = 2/3.
// "If t = 18 months, then n = 12/18 = 2/3"
// t = 18 × 30 × 24 × 60 = 777600 minutes. BasisPoints ≈ 6759.
func TestComputeTurnoverNumber_LocomotiveMaker_18Months(t *testing.T) {
	t.Parallel()
	const t18months int64 = 18 * 30 * 24 * 60 // 777600
	tn := turnover.ComputeTurnoverNumber("locomotive-maker", t18months)
	want := int64(10_000 * 525_600 / 777_600) // 6759
	if tn.BasisPoints != want {
		t.Errorf("LocomotiveMaker BasisPoints = %d, want %d (≈2/3)", tn.BasisPoints, want)
	}
	// Multi-year turnovers: BasisPoints < 10000
	if tn.BasisPoints >= 10_000 {
		t.Errorf("18-month turnover should have BasisPoints < 10000 (multi-year), got %d", tn.BasisPoints)
	}
}

// WineMaturer: 4-year turnover n = 1/4.
// t = 4 × 525600 = 2102400 minutes. BasisPoints = 2500 exactly.
func TestComputeTurnoverNumber_WineMaturer_4Years(t *testing.T) {
	t.Parallel()
	const t4years int64 = 4 * 525_600 // 2102400
	tn := turnover.ComputeTurnoverNumber("wine-maturer", t4years)
	const want int64 = 2_500 // 10000 / 4
	if tn.BasisPoints != want {
		t.Errorf("WineMaturer BasisPoints = %d, want %d (= 1/4 × 10000)", tn.BasisPoints, want)
	}
}

// Zero turnover-time must not panic; BasisPoints is 0.
func TestComputeTurnoverNumber_ZeroTime(t *testing.T) {
	t.Parallel()
	tn := turnover.ComputeTurnoverNumber("zero", 0)
	if tn.BasisPoints != 0 {
		t.Errorf("BasisPoints for zero time = %d, want 0", tn.BasisPoints)
	}
}

// --- TurnoverCycle.Validate --------------------------------------------------

func TestTurnoverCycle_Validate_Complete(t *testing.T) {
	t.Parallel()
	c := turnover.TurnoverCycle{
		ID:                 turnover.NewTurnoverCycleID(),
		TurnoverID:         "t-001",
		StartedAt:          time.Now().Add(-10080 * time.Minute),
		EndedAt:            time.Now(),
		AdvancePence:       1000,
		ReturnedPence:      1200,
		ProductionMinutes:  8000,
		CirculationMinutes: 2080,
	}
	if err := c.Validate(); err != nil {
		t.Errorf("complete cycle invalid: %v", err)
	}
}

func TestTurnoverCycle_Validate_Incomplete(t *testing.T) {
	t.Parallel()
	c := turnover.TurnoverCycle{
		AdvancePence:  1000,
		ReturnedPence: 800,
	}
	if err := c.Validate(); err == nil {
		t.Error("expected ErrCycleNotComplete, got nil")
	}
}

// --- TurnoverCycle.TotalMinutes ---------------------------------------------

func TestTurnoverCycle_TotalMinutes(t *testing.T) {
	t.Parallel()
	c := turnover.TurnoverCycle{ProductionMinutes: 8000, CirculationMinutes: 2080}
	if c.TotalMinutes() != 10080 {
		t.Errorf("TotalMinutes = %d, want 10080", c.TotalMinutes())
	}
}

// --- AverageTurnoverTime ----------------------------------------------------

func TestAverageTurnoverTime_TwoCycles(t *testing.T) {
	t.Parallel()
	cycles := []turnover.TurnoverCycle{
		{ProductionMinutes: 8000, CirculationMinutes: 2000},
		{ProductionMinutes: 8080, CirculationMinutes: 2000},
	}
	avg := turnover.AverageTurnoverTime(cycles)
	if avg != 10040 {
		t.Errorf("avg = %d, want 10040", avg)
	}
}

func TestAverageTurnoverTime_Empty(t *testing.T) {
	t.Parallel()
	if avg := turnover.AverageTurnoverTime(nil); avg != 0 {
		t.Errorf("avg for empty = %d, want 0", avg)
	}
}

// --- ID uniqueness -----------------------------------------------------------

func TestNewTurnoverID_Unique(t *testing.T) {
	t.Parallel()
	a, b := turnover.NewTurnoverID(), turnover.NewTurnoverID()
	if a == b {
		t.Error("NewTurnoverID returned duplicate")
	}
}

func TestNewTurnoverCycleID_Unique(t *testing.T) {
	t.Parallel()
	a, b := turnover.NewTurnoverCycleID(), turnover.NewTurnoverCycleID()
	if a == b {
		t.Error("NewTurnoverCycleID returned duplicate")
	}
}
