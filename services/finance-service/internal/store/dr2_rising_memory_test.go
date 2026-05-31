// Package store — in-memory tests for Vol. III Ch. 43 (DR II Rising Price)
// store methods. Fixtures from Engels Table XI: base=0s, increment=12s, grades=5.
package store

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// ── helpers ───────────────────────────────────────────────────────────────────

func ch43RisingOutcome() rent.RisingPriceDR2Outcome {
	return rent.RisingPriceDR2Outcome{
		Variant:                 rent.RisingPriceEventualityAVariant1,
		TotalCapitalShillings:   500,
		TotalRentShillings:      120,
		TotalOutputBushels:      1500,
		PricePerBushelShillings: 3,
	}
}

func ch43RentSequence() rent.RentSequence {
	return rent.RentSequence{
		BaseRent:  0,
		Increment: 12,
		NumGrades: 5,
	}
}

func ch43EngelsEntry() rent.EngelsTableEntry {
	return rent.EngelsTableEntry{
		TableNumber:      11,
		SoilGradeLabel:   "A",
		OutputBushels:    300,
		RentShillings:    0,
		CapitalShillings: 100,
	}
}

// ── RisingPriceDR2Outcome ─────────────────────────────────────────────────────

func TestMemory_DR2Rising_CreateRisingPriceDR2Outcome_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateRisingPriceDR2Outcome(ctx, ch43RisingOutcome())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.Variant != rent.RisingPriceEventualityAVariant1 {
		t.Errorf("Variant = %d, want RisingPriceEventualityAVariant1 (%d)",
			saved.Variant, rent.RisingPriceEventualityAVariant1)
	}
}

func TestMemory_DR2Rising_ListRisingPriceDR2Outcomes_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListRisingPriceDR2Outcomes(context.Background())
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

func TestMemory_DR2Rising_CreateRisingPriceDR2Outcome_VariantPreserved(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	o := ch43RisingOutcome()
	o.Variant = rent.RisingPriceEventualityBVariant3
	saved, err := m.CreateRisingPriceDR2Outcome(ctx, o)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListRisingPriceDR2Outcomes(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.Variant != rent.RisingPriceEventualityBVariant3 {
				t.Errorf("list Variant = %d, want RisingPriceEventualityBVariant3 (%d)",
					item.Variant, rent.RisingPriceEventualityBVariant3)
			}
			break
		}
	}
	if !found {
		t.Errorf("created outcome ID %q not found in list", saved.ID)
	}
}

func TestMemory_DR2Rising_CreateRisingPriceDR2Outcome_DuplicateID(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateRisingPriceDR2Outcome(ctx, ch43RisingOutcome())
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	dup := saved
	_, err = m.CreateRisingPriceDR2Outcome(ctx, dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}

// ── RentSequence ──────────────────────────────────────────────────────────────

func TestMemory_DR2Rising_CreateRentSequence_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateRentSequence(ctx, ch43RentSequence())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.BaseRent != 0 {
		t.Errorf("BaseRent = %d, want 0", saved.BaseRent)
	}
	if saved.Increment != 12 {
		t.Errorf("Increment = %d, want 12", saved.Increment)
	}
	if saved.NumGrades != 5 {
		t.Errorf("NumGrades = %d, want 5", saved.NumGrades)
	}
}

func TestMemory_DR2Rising_ListRentSequences_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListRentSequences(context.Background())
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

func TestMemory_DR2Rising_CreateRentSequence_FieldsPreserved(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateRentSequence(ctx, ch43RentSequence())
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListRentSequences(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.BaseRent != 0 {
				t.Errorf("list BaseRent = %d, want 0", item.BaseRent)
			}
			if item.Increment != 12 {
				t.Errorf("list Increment = %d, want 12", item.Increment)
			}
			if item.NumGrades != 5 {
				t.Errorf("list NumGrades = %d, want 5", item.NumGrades)
			}
			break
		}
	}
	if !found {
		t.Errorf("created rent-sequence ID %q not found in list", saved.ID)
	}
}

func TestMemory_DR2Rising_CreateRentSequence_DuplicateID(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateRentSequence(ctx, ch43RentSequence())
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	dup := saved
	_, err = m.CreateRentSequence(ctx, dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}

// ── EngelsTableEntry ──────────────────────────────────────────────────────────

func TestMemory_DR2Rising_CreateEngelsTableEntry_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateEngelsTableEntry(ctx, ch43EngelsEntry())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.SoilGradeLabel != "A" {
		t.Errorf("SoilGradeLabel = %q, want %q", saved.SoilGradeLabel, "A")
	}
	if saved.TableNumber != 11 {
		t.Errorf("TableNumber = %d, want 11", saved.TableNumber)
	}
	if saved.OutputBushels != 300 {
		t.Errorf("OutputBushels = %d, want 300", saved.OutputBushels)
	}
}

func TestMemory_DR2Rising_ListEngelsTableEntries_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListEngelsTableEntries(context.Background())
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

func TestMemory_DR2Rising_CreateEngelsTableEntry_SoilGradeLabelPreserved(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	e := ch43EngelsEntry()
	e.SoilGradeLabel = "a"
	saved, err := m.CreateEngelsTableEntry(ctx, e)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListEngelsTableEntries(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.SoilGradeLabel != "a" {
				t.Errorf("list SoilGradeLabel = %q, want %q", item.SoilGradeLabel, "a")
			}
			break
		}
	}
	if !found {
		t.Errorf("created Engels-table entry ID %q not found in list", saved.ID)
	}
}

func TestMemory_DR2Rising_CreateEngelsTableEntry_DuplicateID(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateEngelsTableEntry(ctx, ch43EngelsEntry())
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	dup := saved
	_, err = m.CreateEngelsTableEntry(ctx, dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}
