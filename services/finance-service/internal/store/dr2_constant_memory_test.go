// Package store — in-memory tests for Vol. III Ch. 41 (DR II Constant Price)
// store methods. Fixtures from Marx Vol. III Ch. 41 Table II: Soil D,
// TotalRentalBP=1800, RateOfSurplusProfitBP=36000.
package store

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// ch41Outcome returns a DifferentialRentIIOutcome with Marx Ch. 41 Table II values.
func ch41Outcome() rent.DifferentialRentIIOutcome {
	return rent.DifferentialRentIIOutcome{
		ParcelID:              "soil-D",
		Case:                  rent.CaseProportional,
		AdditionalCapitalBP:   250,
		NewRentPerAcreBP:      1800,
		TotalRentalBP:         1800,
		RateOfSurplusProfitBP: 36000,
	}
}

// ── CreateDifferentialRentIIOutcome / ListDifferentialRentIIOutcomes ──────────

func TestMemory_DR2Constant_CreateOutcome_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	o := ch41Outcome()
	saved, err := m.CreateDifferentialRentIIOutcome(ctx, o)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.Case != rent.CaseProportional {
		t.Errorf("Case = %d, want CaseProportional (%d)", saved.Case, rent.CaseProportional)
	}
	if saved.TotalRentalBP != 1800 {
		t.Errorf("TotalRentalBP = %d, want 1800", saved.TotalRentalBP)
	}
	if saved.RateOfSurplusProfitBP != 36000 {
		t.Errorf("RateOfSurplusProfitBP = %d, want 36000", saved.RateOfSurplusProfitBP)
	}
	if saved.ParcelID != "soil-D" {
		t.Errorf("ParcelID = %q, want %q", saved.ParcelID, "soil-D")
	}
}

func TestMemory_DR2Constant_ListOutcomes_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListDifferentialRentIIOutcomes(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if items == nil {
		t.Error("list must return non-nil slice on empty store")
	}
	if len(items) != 0 {
		t.Errorf("list len = %d, want 0 on empty store", len(items))
	}
}

func TestMemory_DR2Constant_ListOutcomes_ContainsCreated(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateDifferentialRentIIOutcome(ctx, ch41Outcome())
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListDifferentialRentIIOutcomes(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) < 1 {
		t.Fatal("list returned 0 items, want at least 1")
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created outcome ID %q not found in list", saved.ID)
	}
}

func TestMemory_DR2Constant_CreateOutcome_DuplicateID(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateDifferentialRentIIOutcome(ctx, ch41Outcome())
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	dup := saved
	_, err = m.CreateDifferentialRentIIOutcome(ctx, dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}

// ── CreateIntensiveExtensiveComparison / ListIntensiveExtensiveComparisons ────

// ch41Comparison returns an IntensiveExtensiveComparison where intensive
// rent per acre (1800bp) > extensive rent per acre (900bp), same total rental.
func ch41Comparison() rent.IntensiveExtensiveComparison {
	return rent.IntensiveExtensiveComparison{
		IntensiveRentPerAcreBP: 1800,
		ExtensiveRentPerAcreBP: 900,
		TotalCapitalBP:         500,
		TotalRentalSameBP:      1800,
	}
}

func TestMemory_DR2Constant_CreateComparison_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	c := ch41Comparison()
	saved, err := m.CreateIntensiveExtensiveComparison(ctx, c)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.IntensiveRentPerAcreBP != 1800 {
		t.Errorf("IntensiveRentPerAcreBP = %d, want 1800", saved.IntensiveRentPerAcreBP)
	}
	if saved.ExtensiveRentPerAcreBP != 900 {
		t.Errorf("ExtensiveRentPerAcreBP = %d, want 900", saved.ExtensiveRentPerAcreBP)
	}
	if saved.TotalCapitalBP != 500 {
		t.Errorf("TotalCapitalBP = %d, want 500", saved.TotalCapitalBP)
	}
	if saved.TotalRentalSameBP != 1800 {
		t.Errorf("TotalRentalSameBP = %d, want 1800", saved.TotalRentalSameBP)
	}
}

func TestMemory_DR2Constant_ListComparisons_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListIntensiveExtensiveComparisons(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if items == nil {
		t.Error("list must return non-nil slice on empty store")
	}
	if len(items) != 0 {
		t.Errorf("list len = %d, want 0 on empty store", len(items))
	}
}

func TestMemory_DR2Constant_ListComparisons_ContainsCreated(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateIntensiveExtensiveComparison(ctx, ch41Comparison())
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListIntensiveExtensiveComparisons(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) < 1 {
		t.Fatal("list returned 0 items, want at least 1")
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created comparison ID %q not found in list", saved.ID)
	}
}

func TestMemory_DR2Constant_CreateComparison_DuplicateID(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateIntensiveExtensiveComparison(ctx, ch41Comparison())
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	dup := saved
	_, err = m.CreateIntensiveExtensiveComparison(ctx, dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}
