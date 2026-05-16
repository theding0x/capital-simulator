package simulation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

func TestComputeRealRent_HalvedCurrency(t *testing.T) {
	t.Parallel()
	// §depreciation factor: prices doubled → factor 0.5 → real rent halved.
	// "continuous rise in the price of corn ... swelled the money capital of the farmer"
	nominal := simulation.Pence(1200)
	md := simulation.MoneyDepreciation{DepreciationFactor: 0.5}
	got := simulation.ComputeRealRent(nominal, md)
	if got != 600 {
		t.Fatalf("want 600, got %d", got)
	}
}

func TestComputeRealRent_Invariant(t *testing.T) {
	t.Parallel()
	// Invariant: int64(ComputeRealRent) == int64(float64(nominal) * factor)
	nominal := simulation.Pence(1000)
	factor := 0.75
	md := simulation.MoneyDepreciation{DepreciationFactor: factor}
	got := simulation.ComputeRealRent(nominal, md)
	want := simulation.Pence(int64(float64(nominal) * factor))
	if got != want {
		t.Fatalf("want %d, got %d", want, got)
	}
}

func TestComputeRealRent_NoDepreciation(t *testing.T) {
	t.Parallel()
	// Factor 1.0 — currency unchanged, real rent equals nominal.
	nominal := simulation.Pence(800)
	md := simulation.MoneyDepreciation{DepreciationFactor: 1.0}
	got := simulation.ComputeRealRent(nominal, md)
	if got != nominal {
		t.Fatalf("want %d, got %d", nominal, got)
	}
}

func TestComputeFarmingSurplus_HarrisonFixture(t *testing.T) {
	t.Parallel()
	// Harrison's estimate: farmer accumulates fifty or a hundred pounds profit.
	// Pedagogical magnitudes: revenue 2000, rent 800, wages 400 → profit 800.
	s := simulation.ComputeFarmingSurplus(2000, 800, 400)
	if s.Profit != 800 {
		t.Fatalf("want profit 800, got %d", s.Profit)
	}
	// Invariant: Profit == Revenue - NominalRent - WageCosts
	if s.Profit != s.Revenue-s.NominalRent-s.WageCosts {
		t.Fatal("profit invariant violated")
	}
}

func TestComputeFarmingSurplus_ZeroRevenue(t *testing.T) {
	t.Parallel()
	s := simulation.ComputeFarmingSurplus(0, 0, 0)
	if s.Profit != 0 {
		t.Fatalf("want 0, got %d", s.Profit)
	}
}

func TestFarmTenure_Validate_Metayer(t *testing.T) {
	t.Parallel()
	// §15th-century metayer: advances one part of stock, product split 50/50.
	ft := simulation.FarmTenure{
		Form:                 simulation.TenantFormMetayer,
		LeasePeriodYears:     1,
		RentPence:            0,
		CapitalAdvancedPence: 500,
	}
	if err := ft.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFarmTenure_Validate_LongLease(t *testing.T) {
	t.Parallel()
	// §long-lease effect: 99-year lease, capitalist-farmer, fixed nominal rent.
	ft := simulation.FarmTenure{
		Form:             simulation.TenantFormCapitalistFarmer,
		LeasePeriodYears: 99,
		RentPence:        1200,
	}
	if err := ft.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFarmTenure_Validate_Bailiff(t *testing.T) {
	t.Parallel()
	ft := simulation.FarmTenure{
		Form:             simulation.TenantFormBailiff,
		LeasePeriodYears: 1,
	}
	if err := ft.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFarmTenure_Validate_MissingForm(t *testing.T) {
	t.Parallel()
	ft := simulation.FarmTenure{LeasePeriodYears: 10}
	if err := ft.Validate(); err == nil {
		t.Fatal("expected error for missing form")
	}
}

func TestFarmTenure_Validate_InvalidForm(t *testing.T) {
	t.Parallel()
	ft := simulation.FarmTenure{Form: "lord", LeasePeriodYears: 10}
	if err := ft.Validate(); err == nil {
		t.Fatal("expected error for invalid form")
	}
}

func TestFarmTenure_Validate_ZeroLease(t *testing.T) {
	t.Parallel()
	ft := simulation.FarmTenure{Form: simulation.TenantFormMetayer, LeasePeriodYears: 0}
	if err := ft.Validate(); err == nil {
		t.Fatal("expected error for zero lease_period_years")
	}
}

func TestFarmTenure_Validate_NegativeRent(t *testing.T) {
	t.Parallel()
	ft := simulation.FarmTenure{Form: simulation.TenantFormMetayer, LeasePeriodYears: 1, RentPence: -1}
	if err := ft.Validate(); err == nil {
		t.Fatal("expected error for negative rent_pence")
	}
}

func TestNewFarmTenureID_Unique(t *testing.T) {
	t.Parallel()
	id1 := simulation.NewFarmTenureID()
	id2 := simulation.NewFarmTenureID()
	if id1 == id2 {
		t.Fatal("NewFarmTenureID must return unique IDs")
	}
	if len(string(id1)) != 24 {
		t.Fatalf("expected 24 hex chars, got %d", len(string(id1)))
	}
}
