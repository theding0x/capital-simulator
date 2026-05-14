package agent_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

func TestStandardiseWage(t *testing.T) {
	t.Parallel()
	// Spec fixture: FR nominal 36d, works 720 min; standardised to 600 min ref day
	// Expected: 36 * 600 / 720 = 30
	w := agent.DayWage{CountryCode: "FR", NominalPence: 36, WorkingDayMinutes: 720}
	got := agent.StandardiseWage(w, 600)
	if got.CountryCode != "FR" {
		t.Fatalf("country_code: got %s, want FR", got.CountryCode)
	}
	if got.Amount != 30 {
		t.Fatalf("amount: got %d, want 30", got.Amount)
	}
}

func TestStandardiseWageIdentityWhenSameDay(t *testing.T) {
	t.Parallel()
	w := agent.DayWage{CountryCode: "GB", NominalPence: 48, WorkingDayMinutes: 600}
	got := agent.StandardiseWage(w, 600)
	if got.Amount != 48 {
		t.Fatalf("amount: got %d, want 48", got.Amount)
	}
}

func TestComputeRelativePriceInversion(t *testing.T) {
	t.Parallel()
	// Marx's inversion: GB higher nominal wage (60d) + higher intensity (4×) → lower cost
	// to capitalist than DE lower nominal wage (24d) + lower intensity (1×).
	// GB: 60 / 4.0 = 15.0 < DE: 24 / 1.0 = 24.0 ✓
	gbWage := agent.DayWage{CountryCode: "GB", NominalPence: 60, WorkingDayMinutes: 600}
	gbIntensity := agent.NationalIntensity{CountryCode: "GB", Factor: 4.0}
	deWage := agent.DayWage{CountryCode: "DE", NominalPence: 24, WorkingDayMinutes: 720}
	deIntensity := agent.NationalIntensity{CountryCode: "DE", Factor: 1.0}

	gbPrice := agent.ComputeRelativePrice(gbWage, gbIntensity)
	dePrice := agent.ComputeRelativePrice(deWage, deIntensity)

	if gbPrice.Ratio >= dePrice.Ratio {
		t.Fatalf("GB ratio %f should be < DE ratio %f (Marx's inversion)", gbPrice.Ratio, dePrice.Ratio)
	}
}

func TestComputeRelativePriceRatio(t *testing.T) {
	t.Parallel()
	// GB: 48d nominal / 2.0 intensity = 24.0 (NominalPence participates in the ratio)
	w := agent.DayWage{CountryCode: "GB", NominalPence: 48, WorkingDayMinutes: 600}
	ni := agent.NationalIntensity{CountryCode: "GB", Factor: 2.0}
	got := agent.ComputeRelativePrice(w, ni)
	if got.Ratio != 24.0 {
		t.Fatalf("ratio: got %f, want 24.0 (48/2.0)", got.Ratio)
	}
}

func TestComputeRelativePriceSameIntensityDifferentWage(t *testing.T) {
	t.Parallel()
	// Two countries with the same intensity but different wages must differ — the old
	// 1/factor formula made them identical, hiding the wage difference.
	low := agent.DayWage{CountryCode: "A", NominalPence: 20, WorkingDayMinutes: 600}
	high := agent.DayWage{CountryCode: "B", NominalPence: 40, WorkingDayMinutes: 600}
	intensity := agent.NationalIntensity{Factor: 2.0}

	pA := agent.ComputeRelativePrice(low, intensity)
	pB := agent.ComputeRelativePrice(high, intensity)

	if pA.Ratio == pB.Ratio {
		t.Fatalf("same-intensity countries with different wages must have different ratios (got %f == %f)", pA.Ratio, pB.Ratio)
	}
}

func TestNationalIntensityValidate(t *testing.T) {
	t.Parallel()
	if err := (agent.NationalIntensity{CountryCode: "GB", Factor: 0}).Validate(); err == nil {
		t.Fatal("expected error for zero factor")
	}
	if err := (agent.NationalIntensity{CountryCode: "GB", Factor: -1}).Validate(); err == nil {
		t.Fatal("expected error for negative factor")
	}
	if err := (agent.NationalIntensity{CountryCode: "GB", Factor: 2.0}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDayWageValidate(t *testing.T) {
	t.Parallel()
	if err := (agent.DayWage{CountryCode: "GB", NominalPence: 0, WorkingDayMinutes: 600}).Validate(); err == nil {
		t.Fatal("expected error for zero nominal_pence")
	}
	if err := (agent.DayWage{CountryCode: "GB", NominalPence: 48, WorkingDayMinutes: 0}).Validate(); err == nil {
		t.Fatal("expected error for zero working_day_minutes")
	}
	if err := (agent.DayWage{CountryCode: "GB", NominalPence: 48, WorkingDayMinutes: 600}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSpindleRatioValidate(t *testing.T) {
	t.Parallel()
	if err := (agent.SpindleRatio{CountryCode: "GB", SpindlesPerWorker: 0}).Validate(); err == nil {
		t.Fatal("expected error for zero spindles_per_worker")
	}
	if err := (agent.SpindleRatio{CountryCode: "GB", SpindlesPerWorker: 74}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
