package workperiod_test

import (
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/workperiod"
)

// --- ProductionLabourGap tests ------------------------------------------

// TestProductionLabourGap_CornFixture verifies Marx's corn example:
// ~10 labour-days across ~280 production-days (§ "grain ripens in summer").
// Using nanoseconds: 1 day = 86_400_000_000_000 ns.
func TestProductionLabourGap_CornFixture(t *testing.T) {
	t.Parallel()
	const day = int64(86_400_000_000_000) // 1 day in nanoseconds
	labour := 10 * day
	production := 280 * day
	g := workperiod.ProductionLabourGap{
		ID:                  workperiod.NewProductionLabourGapID(),
		WorkingPeriodID:     "wp-corn-test",
		ProductionTimeNanos: production,
		LabourTimeNanos:     labour,
		GapNanos:            production - labour,
		CreatedAt:           time.Now(),
	}
	if err := g.Validate(); err != nil {
		t.Fatalf("corn fixture: unexpected validation error: %v", err)
	}
	if g.GapNanos != production-labour {
		t.Errorf("GapNanos: got %d, want %d", g.GapNanos, production-labour)
	}
}

// TestProductionLabourGap_ForestryFixture verifies the forestry extreme:
// 30 labour-days across 36500 production-days (100 years).
func TestProductionLabourGap_ForestryFixture(t *testing.T) {
	t.Parallel()
	const day = int64(86_400_000_000_000)
	labour := 30 * day
	production := 36500 * day
	g := workperiod.ProductionLabourGap{
		ID:                  workperiod.NewProductionLabourGapID(),
		WorkingPeriodID:     "wp-timber-test",
		ProductionTimeNanos: production,
		LabourTimeNanos:     labour,
		GapNanos:            production - labour,
		CreatedAt:           time.Now(),
	}
	if err := g.Validate(); err != nil {
		t.Fatalf("forestry fixture: unexpected validation error: %v", err)
	}
}

// TestProductionLabourGap_InvariantViolation ensures ProductionTime < LabourTime is rejected.
func TestProductionLabourGap_InvariantViolation(t *testing.T) {
	t.Parallel()
	g := workperiod.ProductionLabourGap{
		ID:                  workperiod.NewProductionLabourGapID(),
		WorkingPeriodID:     "wp-test",
		ProductionTimeNanos: 100,
		LabourTimeNanos:     200, // labour > production — invalid
		GapNanos:            -100,
	}
	if err := g.Validate(); err == nil {
		t.Error("expected validation error for labour > production, got nil")
	}
}

// TestProductionLabourGap_WrongGap ensures gap != production - labour is rejected.
func TestProductionLabourGap_WrongGap(t *testing.T) {
	t.Parallel()
	g := workperiod.ProductionLabourGap{
		ID:                  workperiod.NewProductionLabourGapID(),
		WorkingPeriodID:     "wp-test",
		ProductionTimeNanos: 1000,
		LabourTimeNanos:     100,
		GapNanos:            999, // wrong: should be 900
	}
	if err := g.Validate(); err == nil {
		t.Error("expected validation error for incorrect gap, got nil")
	}
}

// TestProductionLabourGap_ZeroLabour ensures zero labour_time_nanos is rejected.
func TestProductionLabourGap_ZeroLabour(t *testing.T) {
	t.Parallel()
	g := workperiod.ProductionLabourGap{
		WorkingPeriodID:     "wp-test",
		ProductionTimeNanos: 1000,
		LabourTimeNanos:     0,
		GapNanos:            1000,
	}
	if err := g.Validate(); err == nil {
		t.Error("expected validation error for zero labour_time_nanos, got nil")
	}
}

// TestProductionLabourGap_MissingWorkingPeriod ensures empty WorkingPeriodID is rejected.
func TestProductionLabourGap_MissingWorkingPeriod(t *testing.T) {
	t.Parallel()
	g := workperiod.ProductionLabourGap{
		ProductionTimeNanos: 1000,
		LabourTimeNanos:     100,
		GapNanos:            900,
	}
	if err := g.Validate(); err == nil {
		t.Error("expected validation error for missing working_period_id, got nil")
	}
}

// --- NaturalProcessActivation tests ------------------------------------

// TestNaturalProcessActivation_CornGrowth verifies the agricultural-growth activation.
func TestNaturalProcessActivation_CornGrowth(t *testing.T) {
	t.Parallel()
	started := time.Date(1871, 4, 1, 0, 0, 0, 0, time.UTC)
	ended := time.Date(1871, 10, 1, 0, 0, 0, 0, time.UTC)
	a := workperiod.NaturalProcessActivation{
		ID:                        workperiod.NewNaturalProcessActivationID(),
		WorkingPeriodID:           "wp-corn",
		Kind:                      workperiod.NaturalAgriculturalGrowth,
		StartedAt:                 started,
		EndedAt:                   &ended,
		RequiresLabourMaintenance: false,
	}
	if err := a.Validate(); err != nil {
		t.Fatalf("corn activation: unexpected error: %v", err)
	}
}

// TestNaturalProcessActivation_MaintenanceRequired verifies that
// RequiresLabourMaintenance=true is valid (cattle require periodic labour).
func TestNaturalProcessActivation_MaintenanceRequired(t *testing.T) {
	t.Parallel()
	started := time.Date(1871, 1, 1, 0, 0, 0, 0, time.UTC)
	a := workperiod.NaturalProcessActivation{
		ID:                        workperiod.NewNaturalProcessActivationID(),
		WorkingPeriodID:           "wp-cattle",
		Kind:                      workperiod.NaturalAnimalMaturation,
		StartedAt:                 started,
		RequiresLabourMaintenance: true,
	}
	if err := a.Validate(); err != nil {
		t.Fatalf("cattle maintenance activation: unexpected error: %v", err)
	}
}

// TestNaturalProcessActivation_UnknownKind ensures an unknown kind is rejected.
func TestNaturalProcessActivation_UnknownKind(t *testing.T) {
	t.Parallel()
	a := workperiod.NaturalProcessActivation{
		ID:              workperiod.NewNaturalProcessActivationID(),
		WorkingPeriodID: "wp-test",
		Kind:            workperiod.NaturalProcessKind("unknown_kind"),
		StartedAt:       time.Now(),
	}
	if err := a.Validate(); err == nil {
		t.Error("expected validation error for unknown kind, got nil")
	}
}

// TestNaturalProcessActivation_EndedBeforeStarted ensures reversed times are rejected.
func TestNaturalProcessActivation_EndedBeforeStarted(t *testing.T) {
	t.Parallel()
	started := time.Date(1871, 6, 1, 0, 0, 0, 0, time.UTC)
	ended := time.Date(1871, 1, 1, 0, 0, 0, 0, time.UTC) // before started
	a := workperiod.NaturalProcessActivation{
		ID:              workperiod.NewNaturalProcessActivationID(),
		WorkingPeriodID: "wp-test",
		Kind:            workperiod.NaturalFermentation,
		StartedAt:       started,
		EndedAt:         &ended,
	}
	if err := a.Validate(); err == nil {
		t.Error("expected validation error for ended_at before started_at, got nil")
	}
}

// TestNaturalProcessActivation_AllKindsValid verifies every NaturalProcessKind constant passes.
func TestNaturalProcessActivation_AllKindsValid(t *testing.T) {
	t.Parallel()
	kinds := []workperiod.NaturalProcessKind{
		workperiod.NaturalRipening,
		workperiod.NaturalFermentation,
		workperiod.NaturalTanning,
		workperiod.NaturalDrying,
		workperiod.NaturalOther,
		workperiod.NaturalAgriculturalGrowth,
		workperiod.NaturalForestGrowth,
		workperiod.NaturalAnimalMaturation,
		workperiod.NaturalChemical,
		workperiod.NaturalSeasonalRest,
	}
	for _, k := range kinds {
		k := k
		t.Run(string(k), func(t *testing.T) {
			t.Parallel()
			a := workperiod.NaturalProcessActivation{
				ID:              workperiod.NewNaturalProcessActivationID(),
				WorkingPeriodID: "wp-test",
				Kind:            k,
				StartedAt:       time.Now(),
			}
			if err := a.Validate(); err != nil {
				t.Errorf("kind %q: unexpected error: %v", k, err)
			}
		})
	}
}

// --- NaturalSubject tests -----------------------------------------------

// TestNaturalSubject_Living verifies a living subject (cattle) is valid.
func TestNaturalSubject_Living(t *testing.T) {
	t.Parallel()
	s := workperiod.NaturalSubject{
		ID:              workperiod.NewNaturalSubjectID(),
		WorkingPeriodID: "wp-cattle",
		CommodityID:     "cattle",
		Description:     "Cattle raised on Bakewell farm; living subject requiring periodic feeding",
		IsLiving:        true,
		CreatedAt:       time.Now(),
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("living subject: unexpected error: %v", err)
	}
}

// TestNaturalSubject_Inert verifies an inert subject (wine) is valid.
func TestNaturalSubject_Inert(t *testing.T) {
	t.Parallel()
	s := workperiod.NaturalSubject{
		ID:              workperiod.NewNaturalSubjectID(),
		WorkingPeriodID: "wp-wine",
		CommodityID:     "wine",
		Description:     "Wine fermenting in barrels; inert subject requiring no labour during fermentation",
		IsLiving:        false,
		CreatedAt:       time.Now(),
	}
	if err := s.Validate(); err != nil {
		t.Fatalf("inert subject: unexpected error: %v", err)
	}
}

// TestNaturalSubject_MissingFields ensures required fields are enforced.
func TestNaturalSubject_MissingFields(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		s    workperiod.NaturalSubject
	}{
		{"missing working_period_id", workperiod.NaturalSubject{CommodityID: "corn", Description: "d"}},
		{"missing commodity_id", workperiod.NaturalSubject{WorkingPeriodID: "wp-test", Description: "d"}},
		{"missing description", workperiod.NaturalSubject{WorkingPeriodID: "wp-test", CommodityID: "corn"}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.s.Validate(); err == nil {
				t.Errorf("%s: expected validation error, got nil", tc.name)
			}
		})
	}
}

// --- NaturalProcessIndustry tests ----------------------------------------

// TestNaturalProcessIndustry_Agriculture verifies Marx's agriculture benchmark:
// 280 production-days, 10 labour-days → 357 basis points.
func TestNaturalProcessIndustry_Agriculture(t *testing.T) {
	t.Parallel()
	i := workperiod.NaturalProcessIndustry{
		ID:                         workperiod.NewNaturalProcessIndustryID(),
		IndustryName:               "Agriculture (corn)",
		TypicalProductionDaysCount: 280,
		TypicalLabourDaysCount:     10,
		RatioBasisPoints:           357, // 10000 * 10 / 280 = 357
		CreatedAt:                  time.Now(),
	}
	if err := i.Validate(); err != nil {
		t.Fatalf("agriculture benchmark: unexpected error: %v", err)
	}
}

// TestNaturalProcessIndustry_Forestry verifies the forestry extreme:
// 36500 production-days, 30 labour-days → 8 basis points.
func TestNaturalProcessIndustry_Forestry(t *testing.T) {
	t.Parallel()
	i := workperiod.NaturalProcessIndustry{
		ID:                         workperiod.NewNaturalProcessIndustryID(),
		IndustryName:               "Forestry (100-year timber)",
		TypicalProductionDaysCount: 36500,
		TypicalLabourDaysCount:     30,
		RatioBasisPoints:           8, // 10000 * 30 / 36500 = 8
		CreatedAt:                  time.Now(),
	}
	if err := i.Validate(); err != nil {
		t.Fatalf("forestry benchmark: unexpected error: %v", err)
	}
}

// TestNaturalProcessIndustry_RatioBounds ensures ratio outside [0,10000] is rejected.
func TestNaturalProcessIndustry_RatioBounds(t *testing.T) {
	t.Parallel()
	tooHigh := workperiod.NaturalProcessIndustry{
		IndustryName:               "test",
		TypicalProductionDaysCount: 10,
		TypicalLabourDaysCount:     10,
		RatioBasisPoints:           10001,
	}
	if err := tooHigh.Validate(); err == nil {
		t.Error("expected error for ratio_basis_points > 10000, got nil")
	}
	negative := workperiod.NaturalProcessIndustry{
		IndustryName:               "test",
		TypicalProductionDaysCount: 10,
		TypicalLabourDaysCount:     10,
		RatioBasisPoints:           -1,
	}
	if err := negative.Validate(); err == nil {
		t.Error("expected error for negative ratio_basis_points, got nil")
	}
}

// TestNaturalProcessIndustry_LabourExceedsProduction ensures the constraint is enforced.
func TestNaturalProcessIndustry_LabourExceedsProduction(t *testing.T) {
	t.Parallel()
	i := workperiod.NaturalProcessIndustry{
		IndustryName:               "test",
		TypicalProductionDaysCount: 10,
		TypicalLabourDaysCount:     20, // labour > production: invalid
		RatioBasisPoints:           5000,
	}
	if err := i.Validate(); err == nil {
		t.Error("expected error for labour > production, got nil")
	}
}

// TestNaturalProcessIndustry_MissingName ensures industry_name is required.
func TestNaturalProcessIndustry_MissingName(t *testing.T) {
	t.Parallel()
	i := workperiod.NaturalProcessIndustry{
		TypicalProductionDaysCount: 280,
		TypicalLabourDaysCount:     10,
		RatioBasisPoints:           357,
	}
	if err := i.Validate(); err == nil {
		t.Error("expected error for missing industry_name, got nil")
	}
}
