// Package store — in-memory tests for Vol. III Ch. 42 (DR II Falling Price)
// store methods. Fixtures from Marx Vol. III Ch. 42 Table IV: Grade A excluded,
// Grade B becomes regulator at 45s (540 pence) per quarter.
package store

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// ── helpers ───────────────────────────────────────────────────────────────────

// ch42ExclusionEvent returns a Ch. 42 Table IV soil-exclusion event fixture.
func ch42ExclusionEvent() rent.SoilExclusionEvent {
	return rent.SoilExclusionEvent{
		ExcludedGrade:          rent.SoilGradeA,
		NewRegulatingGrade:     rent.SoilGradeB,
		NewPriceOfProductionBP: 540,
	}
}

// ch42FallingOutcome returns a Ch. 42 Table IV falling-price DR II outcome.
func ch42FallingOutcome() rent.FallingPriceDR2Outcome {
	return rent.FallingPriceDR2Outcome{
		Variant:             rent.FallingPriceConstantProductivity,
		ExclusionEventID:    "test-exclusion-01",
		TotalGrainRentQtrs:  250,
		TotalMoneyRentalBP:  1350,
		TotalOutputQuarters: 2400,
		TotalCapitalBP:      1500,
	}
}

// ch42NormalCapital returns a pre-1848 England normal-capital-per-acre fixture
// (Marx: £8/acre = 800 pence).
func ch42NormalCapital() rent.NormalCapitalPerAcre {
	return rent.NormalCapitalPerAcre{
		PeriodLabel:      "pre-1848 England",
		CapitalPerAcreBP: 800,
	}
}

// ── SoilExclusionEvent ────────────────────────────────────────────────────────

func TestMemory_DR2Falling_CreateSoilExclusionEvent_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateSoilExclusionEvent(ctx, ch42ExclusionEvent())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.ExcludedGrade != rent.SoilGradeA {
		t.Errorf("ExcludedGrade = %v, want SoilGradeA", saved.ExcludedGrade)
	}
	if saved.NewRegulatingGrade != rent.SoilGradeB {
		t.Errorf("NewRegulatingGrade = %v, want SoilGradeB", saved.NewRegulatingGrade)
	}
	if saved.NewPriceOfProductionBP != 540 {
		t.Errorf("NewPriceOfProductionBP = %d, want 540", saved.NewPriceOfProductionBP)
	}

	got, err := m.GetSoilExclusionEvent(ctx, saved.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != saved.ID {
		t.Errorf("get ID = %q, want %q", got.ID, saved.ID)
	}
	if got.ExcludedGrade != saved.ExcludedGrade {
		t.Errorf("get ExcludedGrade = %v, want %v", got.ExcludedGrade, saved.ExcludedGrade)
	}
	if got.NewPriceOfProductionBP != saved.NewPriceOfProductionBP {
		t.Errorf("get NewPriceOfProductionBP = %d, want %d",
			got.NewPriceOfProductionBP, saved.NewPriceOfProductionBP)
	}
}

func TestMemory_DR2Falling_GetSoilExclusionEvent_NotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	_, err := m.GetSoilExclusionEvent(context.Background(), "unknown-id")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_DR2Falling_ListSoilExclusionEvents_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListSoilExclusionEvents(context.Background())
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

func TestMemory_DR2Falling_ListSoilExclusionEvents_ContainsCreated(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateSoilExclusionEvent(ctx, ch42ExclusionEvent())
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListSoilExclusionEvents(ctx)
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
		t.Errorf("created soil-exclusion event ID %q not found in list", saved.ID)
	}
}

// ── FallingPriceDR2Outcome ────────────────────────────────────────────────────

func TestMemory_DR2Falling_CreateFallingPriceDR2Outcome_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateFallingPriceDR2Outcome(ctx, ch42FallingOutcome())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.Variant != rent.FallingPriceConstantProductivity {
		t.Errorf("Variant = %d, want FallingPriceConstantProductivity (%d)",
			saved.Variant, rent.FallingPriceConstantProductivity)
	}
	if saved.TotalMoneyRentalBP != 1350 {
		t.Errorf("TotalMoneyRentalBP = %d, want 1350", saved.TotalMoneyRentalBP)
	}
	if saved.TotalGrainRentQtrs != 250 {
		t.Errorf("TotalGrainRentQtrs = %d, want 250", saved.TotalGrainRentQtrs)
	}
}

func TestMemory_DR2Falling_ListFallingPriceDR2Outcomes_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListFallingPriceDR2Outcomes(context.Background())
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

func TestMemory_DR2Falling_CreateFallingPriceDR2Outcome_VariantPreserved(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	o := ch42FallingOutcome()
	o.Variant = rent.FallingPriceRisingProductivity
	saved, err := m.CreateFallingPriceDR2Outcome(ctx, o)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListFallingPriceDR2Outcomes(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.Variant != rent.FallingPriceRisingProductivity {
				t.Errorf("list Variant = %d, want FallingPriceRisingProductivity (%d)",
					item.Variant, rent.FallingPriceRisingProductivity)
			}
			break
		}
	}
	if !found {
		t.Errorf("created outcome ID %q not found in list", saved.ID)
	}
}

func TestMemory_DR2Falling_CreateFallingPriceDR2Outcome_DuplicateID(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateFallingPriceDR2Outcome(ctx, ch42FallingOutcome())
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	dup := saved
	_, err = m.CreateFallingPriceDR2Outcome(ctx, dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}

// ── NormalCapitalPerAcre ──────────────────────────────────────────────────────

func TestMemory_DR2Falling_CreateNormalCapitalPerAcre_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateNormalCapitalPerAcre(ctx, ch42NormalCapital())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.PeriodLabel != "pre-1848 England" {
		t.Errorf("PeriodLabel = %q, want %q", saved.PeriodLabel, "pre-1848 England")
	}
	if saved.CapitalPerAcreBP != 800 {
		t.Errorf("CapitalPerAcreBP = %d, want 800", saved.CapitalPerAcreBP)
	}

	items, err := m.ListNormalCapitalPerAcre(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.CapitalPerAcreBP != 800 {
				t.Errorf("list CapitalPerAcreBP = %d, want 800", item.CapitalPerAcreBP)
			}
			break
		}
	}
	if !found {
		t.Errorf("created normal-capital ID %q not found in list", saved.ID)
	}
}

func TestMemory_DR2Falling_ListNormalCapitalPerAcre_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListNormalCapitalPerAcre(context.Background())
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
