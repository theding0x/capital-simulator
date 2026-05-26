package workperiod_test

import (
	"strings"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/workperiod"
)

// --- WorkingPeriod.Multiplier -----------------------------------------------

func TestMultiplier_Discrete(t *testing.T) {
	t.Parallel()
	// SpinningMill1871: 7-day discrete period; each day yields a sellable unit.
	// Multiplier must be 1 regardless of WorkingDayCount.
	wp := workperiod.WorkingPeriod{
		ID:                      workperiod.NewWorkingPeriodID(),
		IndustrialCapitalID:     "mill-1871",
		CommodityID:             "yarn",
		WorkingDayCount:         7,
		Mode:                    workperiod.ModeDiscrete,
		WorkingDayLengthMinutes: 600,
	}
	if got := wp.Multiplier(); got != 1 {
		t.Fatalf("discrete Multiplier = %d, want 1", got)
	}
}

func TestMultiplier_Connected(t *testing.T) {
	t.Parallel()
	// LocomotiveWorks: 84 working-days before a single locomotive exists.
	// Multiplier must equal WorkingDayCount.
	wp := workperiod.WorkingPeriod{
		ID:                      workperiod.NewWorkingPeriodID(),
		IndustrialCapitalID:     "locomotive-works",
		CommodityID:             "locomotive",
		WorkingDayCount:         84,
		Mode:                    workperiod.ModeConnected,
		WorkingDayLengthMinutes: 600,
	}
	if got := wp.Multiplier(); got != 84 {
		t.Fatalf("connected Multiplier = %d, want 84", got)
	}
}

func TestMultiplier_DiscreteIgnoresCount(t *testing.T) {
	t.Parallel()
	// Discrete processes with varying WorkingDayCount always return 1.
	for _, n := range []int64{1, 7, 30, 365} {
		wp := workperiod.WorkingPeriod{
			IndustrialCapitalID:     "mill",
			CommodityID:             "yarn",
			WorkingDayCount:         n,
			Mode:                    workperiod.ModeDiscrete,
			WorkingDayLengthMinutes: 600,
		}
		if got := wp.Multiplier(); got != 1 {
			t.Fatalf("discrete Multiplier(count=%d) = %d, want 1", n, got)
		}
	}
}

// --- WorkingPeriod.Validate -------------------------------------------------

func TestWorkingPeriod_Validate_Happy(t *testing.T) {
	t.Parallel()
	wp := workperiod.WorkingPeriod{
		IndustrialCapitalID:     "mill-1871",
		CommodityID:             "yarn",
		WorkingDayCount:         1,
		Mode:                    workperiod.ModeDiscrete,
		WorkingDayLengthMinutes: 600,
	}
	if err := wp.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestWorkingPeriod_Validate_MissingFields(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		wp   workperiod.WorkingPeriod
		want string
	}{
		{
			name: "missing industrial_capital_id",
			wp: workperiod.WorkingPeriod{
				CommodityID: "yarn", WorkingDayCount: 1,
				Mode: workperiod.ModeDiscrete, WorkingDayLengthMinutes: 600,
			},
			want: "industrial_capital_id",
		},
		{
			name: "missing commodity_id",
			wp: workperiod.WorkingPeriod{
				IndustrialCapitalID: "mill", WorkingDayCount: 1,
				Mode: workperiod.ModeDiscrete, WorkingDayLengthMinutes: 600,
			},
			want: "commodity_id",
		},
		{
			name: "zero working_day_count",
			wp: workperiod.WorkingPeriod{
				IndustrialCapitalID: "mill", CommodityID: "yarn",
				WorkingDayCount: 0, Mode: workperiod.ModeDiscrete,
				WorkingDayLengthMinutes: 600,
			},
			want: "working_day_count",
		},
		{
			name: "invalid mode",
			wp: workperiod.WorkingPeriod{
				IndustrialCapitalID: "mill", CommodityID: "yarn",
				WorkingDayCount: 1, Mode: "batch",
				WorkingDayLengthMinutes: 600,
			},
			want: "mode",
		},
		{
			name: "zero working_day_length_minutes",
			wp: workperiod.WorkingPeriod{
				IndustrialCapitalID: "mill", CommodityID: "yarn",
				WorkingDayCount: 1, Mode: workperiod.ModeDiscrete,
				WorkingDayLengthMinutes: 0,
			},
			want: "working_day_length_minutes",
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			err := c.wp.Validate()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Fatalf("error %q does not mention %q", err.Error(), c.want)
			}
		})
	}
}

// --- WorkingPeriodInterruption.Validate -------------------------------------

func TestInterruption_Discrete_DisallowsDeterioration(t *testing.T) {
	t.Parallel()
	// A discrete interruption (SpinningMill stop) must carry zero deterioration.
	intr := workperiod.WorkingPeriodInterruption{
		ID:                 workperiod.NewWorkingPeriodInterruptionID(),
		WorkingPeriodID:    "some-id",
		InterruptedAt:      time.Now(),
		DeteriorationPence: 100,
	}
	if err := intr.Validate(workperiod.ModeDiscrete); err == nil {
		t.Fatal("expected error for discrete interruption with deterioration_pence, got nil")
	}
}

func TestInterruption_Discrete_ZeroDeteriorationOK(t *testing.T) {
	t.Parallel()
	intr := workperiod.WorkingPeriodInterruption{
		ID:                 workperiod.NewWorkingPeriodInterruptionID(),
		WorkingPeriodID:    "some-id",
		InterruptedAt:      time.Now(),
		DeteriorationPence: 0,
	}
	if err := intr.Validate(workperiod.ModeDiscrete); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInterruption_Connected_AllowsDeterioration(t *testing.T) {
	t.Parallel()
	// A locomotive left half-finished during a stoppage accumulates deterioration.
	intr := workperiod.WorkingPeriodInterruption{
		ID:                            workperiod.NewWorkingPeriodInterruptionID(),
		WorkingPeriodID:               "loco-id",
		InterruptedAt:                 time.Now(),
		DeteriorationPence:            500,
		WasteOfMeansOfProductionPence: 200,
	}
	if err := intr.Validate(workperiod.ModeConnected); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInterruption_NegativeDeteriorationRejected(t *testing.T) {
	t.Parallel()
	intr := workperiod.WorkingPeriodInterruption{
		DeteriorationPence: -1,
	}
	if err := intr.Validate(workperiod.ModeConnected); err == nil {
		t.Fatal("expected error for negative deterioration_pence")
	}
}

func TestInterruption_NegativeWasteRejected(t *testing.T) {
	t.Parallel()
	intr := workperiod.WorkingPeriodInterruption{
		WasteOfMeansOfProductionPence: -1,
	}
	if err := intr.Validate(workperiod.ModeConnected); err == nil {
		t.Fatal("expected error for negative waste_of_means_of_production_pence")
	}
}

// --- WorkingPeriodShortening.Validate ---------------------------------------

func TestShortening_Machinery_RequiresFixedCapital(t *testing.T) {
	t.Parallel()
	s := workperiod.WorkingPeriodShortening{
		Kind:               workperiod.ShorteningMachinery,
		OldWorkingDayCount: 84,
		NewWorkingDayCount: 42,
		// AdditionalFixedCapitalPence deliberately zero
	}
	if err := s.Validate(); err == nil {
		t.Fatal("expected error: machinery shortening requires additional_fixed_capital_pence > 0")
	}
}

func TestShortening_Machinery_Happy(t *testing.T) {
	t.Parallel()
	s := workperiod.WorkingPeriodShortening{
		Kind:                        workperiod.ShorteningMachinery,
		OldWorkingDayCount:          84,
		NewWorkingDayCount:          42,
		AdditionalFixedCapitalPence: 5000,
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShortening_ParallelLabour_RequiresLabourCount(t *testing.T) {
	t.Parallel()
	s := workperiod.WorkingPeriodShortening{
		Kind:               workperiod.ShorteningParallelLabour,
		OldWorkingDayCount: 84,
		NewWorkingDayCount: 42,
		// AdditionalLabourCount deliberately zero
	}
	if err := s.Validate(); err == nil {
		t.Fatal("expected error: parallel-labour shortening requires additional_labour_count > 0")
	}
}

func TestShortening_BakewellBreeding_NoExtraCostsRequired(t *testing.T) {
	t.Parallel()
	// Bakewell's selective breeding (cattle maturation 3 y → 1.5 y)
	// shortens the period with no additional fixed capital or labour beyond
	// the breeding programme itself.
	s := workperiod.WorkingPeriodShortening{
		Kind:               workperiod.ShorteningBakewellBreeding,
		OldWorkingDayCount: 1095, // 3 years
		NewWorkingDayCount: 547,  // ~1.5 years
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestShortening_NewCountMustBeLessOld(t *testing.T) {
	t.Parallel()
	// Equal counts are also rejected.
	s := workperiod.WorkingPeriodShortening{
		Kind:               workperiod.ShorteningCooperation,
		OldWorkingDayCount: 10,
		NewWorkingDayCount: 10,
	}
	if err := s.Validate(); err == nil {
		t.Fatal("expected error: new must be less than old")
	}
}

func TestShortening_ZeroCountsRejected(t *testing.T) {
	t.Parallel()
	s := workperiod.WorkingPeriodShortening{
		Kind:               workperiod.ShorteningCooperation,
		OldWorkingDayCount: 0,
		NewWorkingDayCount: 0,
	}
	if err := s.Validate(); err == nil {
		t.Fatal("expected error for zero working day counts")
	}
}

// --- NaturalWorkingPeriodConstraint.Validate --------------------------------

func TestNaturalConstraint_Validate_Happy(t *testing.T) {
	t.Parallel()
	nc := workperiod.NaturalWorkingPeriodConstraint{
		CommodityID:    "corn",
		MinWorkingDays: 365,
		Reason:         "annual crop — grain cannot be harvested before one growing season",
	}
	if err := nc.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNaturalConstraint_Validate_MissingFields(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		nc   workperiod.NaturalWorkingPeriodConstraint
	}{
		{"missing commodity_id", workperiod.NaturalWorkingPeriodConstraint{MinWorkingDays: 365, Reason: "r"}},
		{"zero min_working_days", workperiod.NaturalWorkingPeriodConstraint{CommodityID: "corn", Reason: "r"}},
		{"missing reason", workperiod.NaturalWorkingPeriodConstraint{CommodityID: "corn", MinWorkingDays: 365}},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			if err := c.nc.Validate(); err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}

// --- WorkingPeriodCreditFinancing.Validate ----------------------------------

func TestCreditFinancing_Validate_Happy(t *testing.T) {
	t.Parallel()
	// London builder: own capital £100 + borrowed £2,000.
	cf := workperiod.WorkingPeriodCreditFinancing{
		OwnCapitalPence:      24000,  // £100
		BorrowedCapitalPence: 480000, // £2,000
	}
	if err := cf.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreditFinancing_Validate_NegativeOwnRejected(t *testing.T) {
	t.Parallel()
	cf := workperiod.WorkingPeriodCreditFinancing{
		OwnCapitalPence:      -1,
		BorrowedCapitalPence: 1000,
	}
	if err := cf.Validate(); err == nil {
		t.Fatal("expected error for negative own_capital_pence")
	}
}

func TestCreditFinancing_Validate_ZeroTotalRejected(t *testing.T) {
	t.Parallel()
	cf := workperiod.WorkingPeriodCreditFinancing{
		OwnCapitalPence:      0,
		BorrowedCapitalPence: 0,
	}
	if err := cf.Validate(); err == nil {
		t.Fatal("expected error for zero total financing")
	}
}

// --- ID constructors --------------------------------------------------------

func TestNewWorkingPeriodID_Unique(t *testing.T) {
	t.Parallel()
	a := workperiod.NewWorkingPeriodID()
	b := workperiod.NewWorkingPeriodID()
	if a == b {
		t.Fatal("expected distinct IDs")
	}
	if len(a) != 24 {
		t.Fatalf("expected 24-char hex ID, got %d", len(a))
	}
}

func TestNewWorkingPeriodID_SeedPatternDistinct(t *testing.T) {
	t.Parallel()
	// Runtime IDs must not start with the seed prefix 5eed000000000000001200.
	id := workperiod.NewWorkingPeriodID()
	if strings.HasPrefix(string(id), "5eed") {
		t.Fatalf("runtime ID %q looks like a seed ID", id)
	}
}
