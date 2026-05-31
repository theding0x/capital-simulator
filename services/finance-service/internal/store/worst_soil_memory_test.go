// Package store — in-memory tests for Vol. III Ch. 44 (Worst Soil Rent)
// store methods.
package store

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// ── helpers ───────────────────────────────────────────────────────────────────

func ch44RentOnWorstSoil() rent.RentOnWorstSoil {
	return rent.RentOnWorstSoil{
		ParcelID:               "parcel-a",
		Mechanism:              rent.MechanismBetterSoilRegulates,
		OldPriceOfProductionBP: 30000,
		NewPriceOfProductionBP: 35000,
		RentPerAcreBP:          5000,
	}
}

func ch44SoilImprovementCapital() rent.SoilImprovementCapital {
	return rent.SoilImprovementCapital{
		ParcelID:           "parcel-b",
		CapitalInvestedBP:  50000,
		ProductivityGainBP: 12000,
		AmortisationYears:  15,
		AmortisedYears:     7,
	}
}

func ch44RegulatingPriceShift() rent.RegulatingPriceShift {
	return rent.RegulatingPriceShift{
		FromGrade:  rent.SoilGradeA,
		ToGrade:    rent.SoilGradeB,
		OldPriceBP: 30000,
		NewPriceBP: 35000,
	}
}

// ── RentOnWorstSoil ───────────────────────────────────────────────────────────

func TestMemory_Ch44_CreateRentOnWorstSoil_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateRentOnWorstSoil(ctx, ch44RentOnWorstSoil())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.Mechanism != rent.MechanismBetterSoilRegulates {
		t.Errorf("Mechanism = %d, want MechanismBetterSoilRegulates (%d)",
			saved.Mechanism, rent.MechanismBetterSoilRegulates)
	}
	if saved.RentPerAcreBP != 5000 {
		t.Errorf("RentPerAcreBP = %d, want 5000", saved.RentPerAcreBP)
	}
}

func TestMemory_Ch44_ListRentOnWorstSoils_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListRentOnWorstSoils(context.Background())
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

func TestMemory_Ch44_CreateRentOnWorstSoil_MechanismPreserved(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	r := ch44RentOnWorstSoil()
	r.Mechanism = rent.MechanismDecreasingProductivity
	saved, err := m.CreateRentOnWorstSoil(ctx, r)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListRentOnWorstSoils(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.Mechanism != rent.MechanismDecreasingProductivity {
				t.Errorf("list Mechanism = %d, want MechanismDecreasingProductivity (%d)",
					item.Mechanism, rent.MechanismDecreasingProductivity)
			}
			break
		}
	}
	if !found {
		t.Errorf("created rent-on-worst-soil ID %q not found in list", saved.ID)
	}
}

func TestMemory_Ch44_CreateRentOnWorstSoil_DuplicateID(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateRentOnWorstSoil(ctx, ch44RentOnWorstSoil())
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	dup := saved
	_, err = m.CreateRentOnWorstSoil(ctx, dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}

// ── SoilImprovementCapital ────────────────────────────────────────────────────

func TestMemory_Ch44_CreateSoilImprovementCapital_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateSoilImprovementCapital(ctx, ch44SoilImprovementCapital())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.CapitalInvestedBP != 50000 {
		t.Errorf("CapitalInvestedBP = %d, want 50000", saved.CapitalInvestedBP)
	}
	if saved.AmortisationYears != 15 {
		t.Errorf("AmortisationYears = %d, want 15", saved.AmortisationYears)
	}
}

func TestMemory_Ch44_ListSoilImprovementCapitals_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListSoilImprovementCapitals(context.Background())
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

func TestMemory_Ch44_CreateSoilImprovementCapital_FieldsPreserved(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateSoilImprovementCapital(ctx, ch44SoilImprovementCapital())
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListSoilImprovementCapitals(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.ProductivityGainBP != 12000 {
				t.Errorf("list ProductivityGainBP = %d, want 12000", item.ProductivityGainBP)
			}
			if item.AmortisedYears != 7 {
				t.Errorf("list AmortisedYears = %d, want 7", item.AmortisedYears)
			}
			break
		}
	}
	if !found {
		t.Errorf("created soil-improvement-capital ID %q not found in list", saved.ID)
	}
}

func TestMemory_Ch44_CreateSoilImprovementCapital_DuplicateID(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateSoilImprovementCapital(ctx, ch44SoilImprovementCapital())
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	dup := saved
	_, err = m.CreateSoilImprovementCapital(ctx, dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}

// ── RegulatingPriceShift ──────────────────────────────────────────────────────

func TestMemory_Ch44_CreateRegulatingPriceShift_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateRegulatingPriceShift(ctx, ch44RegulatingPriceShift())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.FromGrade != rent.SoilGradeA {
		t.Errorf("FromGrade = %d, want SoilGradeA (%d)", saved.FromGrade, rent.SoilGradeA)
	}
	if saved.ToGrade != rent.SoilGradeB {
		t.Errorf("ToGrade = %d, want SoilGradeB (%d)", saved.ToGrade, rent.SoilGradeB)
	}
}

func TestMemory_Ch44_ListRegulatingPriceShifts_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListRegulatingPriceShifts(context.Background())
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

func TestMemory_Ch44_CreateRegulatingPriceShift_GradesPreserved(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	s := ch44RegulatingPriceShift()
	s.FromGrade = rent.SoilGradeB
	s.ToGrade = rent.SoilGradeC
	saved, err := m.CreateRegulatingPriceShift(ctx, s)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListRegulatingPriceShifts(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.FromGrade != rent.SoilGradeB {
				t.Errorf("list FromGrade = %d, want SoilGradeB (%d)", item.FromGrade, rent.SoilGradeB)
			}
			if item.ToGrade != rent.SoilGradeC {
				t.Errorf("list ToGrade = %d, want SoilGradeC (%d)", item.ToGrade, rent.SoilGradeC)
			}
			break
		}
	}
	if !found {
		t.Errorf("created regulating-price-shift ID %q not found in list", saved.ID)
	}
}

func TestMemory_Ch44_CreateRegulatingPriceShift_DuplicateID(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateRegulatingPriceShift(ctx, ch44RegulatingPriceShift())
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	dup := saved
	_, err = m.CreateRegulatingPriceShift(ctx, dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}
