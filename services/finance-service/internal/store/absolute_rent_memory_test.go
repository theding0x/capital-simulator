// Package store — in-memory tests for Vol. III Ch. 45 (Absolute Ground-Rent)
// store methods.
package store

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// ── helpers ───────────────────────────────────────────────────────────────────

func ch45AgriculturalComposition() rent.AgriculturalComposition {
	return rent.AgriculturalComposition{
		ConstantCapitalBP:    60,
		VariableCapitalBP:    40,
		CompositionRatioBP:   15000,
		SocialAverageRatioBP: 40000,
		IsBelowAverage:       true,
	}
}

func ch45ValuePriceGap() rent.ValuePriceGap {
	return rent.ValuePriceGap{
		ProductID:                      "agricultural-output",
		ValueLabourMinutes:             140,
		PriceOfProductionLabourMinutes: 120,
		GapLabourMinutes:               20,
	}
}

func ch45AbsoluteRent() rent.AbsoluteRent {
	return rent.AbsoluteRent{
		ParcelID:        "parcel-alpha",
		ValuePriceGapID: "vg-001",
		SurplusValueBP:  40,
		AverageProfitBP: 20,
		AbsoluteRentBP:  20,
	}
}

func ch45AbsoluteRentLimit() rent.AbsoluteRentLimit {
	return rent.AbsoluteRentLimit{
		MaxRentBP:    20,
		ActualRentBP: 20,
		IsAtLimit:    true,
	}
}

// ── AgriculturalComposition ───────────────────────────────────────────────────

func TestMemory_Ch45_CreateAgriculturalComposition_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateAgriculturalComposition(ctx, ch45AgriculturalComposition())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if !saved.IsBelowAverage {
		t.Error("IsBelowAverage must be preserved as true")
	}
	if saved.CompositionRatioBP != 15000 {
		t.Errorf("CompositionRatioBP = %d, want 15000", saved.CompositionRatioBP)
	}
}

func TestMemory_Ch45_ListAgriculturalCompositions_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListAgriculturalCompositions(context.Background())
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

func TestMemory_Ch45_CreateAgriculturalComposition_IsBelowAveragePreserved(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// manufacturing capital: above-average composition
	above := rent.AgriculturalComposition{
		ConstantCapitalBP:    80,
		VariableCapitalBP:    20,
		CompositionRatioBP:   40000,
		SocialAverageRatioBP: 40000,
		IsBelowAverage:       false,
	}
	saved, err := m.CreateAgriculturalComposition(ctx, above)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListAgriculturalCompositions(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.IsBelowAverage {
				t.Error("IsBelowAverage must be false for above-average composition")
			}
			break
		}
	}
	if !found {
		t.Errorf("created agricultural-composition ID %q not found in list", saved.ID)
	}
}

func TestMemory_Ch45_CreateAgriculturalComposition_DuplicateID(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateAgriculturalComposition(ctx, ch45AgriculturalComposition())
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	dup := saved
	_, err = m.CreateAgriculturalComposition(ctx, dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}

// ── ValuePriceGap ─────────────────────────────────────────────────────────────

func TestMemory_Ch45_CreateValuePriceGap_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateValuePriceGap(ctx, ch45ValuePriceGap())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.ProductID != "agricultural-output" {
		t.Errorf("ProductID = %q, want agricultural-output", saved.ProductID)
	}
	if saved.GapLabourMinutes != 20 {
		t.Errorf("GapLabourMinutes = %d, want 20", saved.GapLabourMinutes)
	}
}

func TestMemory_Ch45_ListValuePriceGaps_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListValuePriceGaps(context.Background())
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

func TestMemory_Ch45_CreateValuePriceGap_FieldsPreserved(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateValuePriceGap(ctx, ch45ValuePriceGap())
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListValuePriceGaps(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.ValueLabourMinutes != 140 {
				t.Errorf("list ValueLabourMinutes = %d, want 140", item.ValueLabourMinutes)
			}
			if item.PriceOfProductionLabourMinutes != 120 {
				t.Errorf("list PriceOfProductionLabourMinutes = %d, want 120", item.PriceOfProductionLabourMinutes)
			}
			break
		}
	}
	if !found {
		t.Errorf("created value-price-gap ID %q not found in list", saved.ID)
	}
}

// ── AbsoluteRent ──────────────────────────────────────────────────────────────

func TestMemory_Ch45_CreateAbsoluteRent_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateAbsoluteRent(ctx, ch45AbsoluteRent())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.ParcelID != "parcel-alpha" {
		t.Errorf("ParcelID = %q, want parcel-alpha", saved.ParcelID)
	}
	if saved.AbsoluteRentBP != 20 {
		t.Errorf("AbsoluteRentBP = %d, want 20", saved.AbsoluteRentBP)
	}
	if saved.SurplusValueBP != 40 {
		t.Errorf("SurplusValueBP = %d, want 40", saved.SurplusValueBP)
	}
	if saved.AverageProfitBP != 20 {
		t.Errorf("AverageProfitBP = %d, want 20", saved.AverageProfitBP)
	}
}

func TestMemory_Ch45_GetAbsoluteRent_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateAbsoluteRent(ctx, ch45AbsoluteRent())
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := m.GetAbsoluteRent(ctx, saved.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != saved.ID {
		t.Errorf("get ID = %q, want %q", got.ID, saved.ID)
	}
	if got.AbsoluteRentBP != saved.AbsoluteRentBP {
		t.Errorf("get AbsoluteRentBP = %d, want %d", got.AbsoluteRentBP, saved.AbsoluteRentBP)
	}
}

func TestMemory_Ch45_GetAbsoluteRent_NotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	_, err := m.GetAbsoluteRent(context.Background(), rent.AbsoluteRentID("nope"))
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound for unknown ID", err)
	}
}

func TestMemory_Ch45_ListAbsoluteRents_ContainsCreated(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateAbsoluteRent(ctx, ch45AbsoluteRent())
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListAbsoluteRents(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("created absolute-rent ID %q not found in list", saved.ID)
	}
}

func TestMemory_Ch45_ListAbsoluteRents_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListAbsoluteRents(context.Background())
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

// ── AbsoluteRentLimit ─────────────────────────────────────────────────────────

func TestMemory_Ch45_CreateAbsoluteRentLimit_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateAbsoluteRentLimit(ctx, ch45AbsoluteRentLimit())
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if !saved.IsAtLimit {
		t.Error("IsAtLimit must be preserved as true when actual >= max")
	}
	if saved.MaxRentBP != 20 {
		t.Errorf("MaxRentBP = %d, want 20", saved.MaxRentBP)
	}
	if saved.ActualRentBP != 20 {
		t.Errorf("ActualRentBP = %d, want 20", saved.ActualRentBP)
	}
}

func TestMemory_Ch45_ListAbsoluteRentLimits_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListAbsoluteRentLimits(context.Background())
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

func TestMemory_Ch45_CreateAbsoluteRentLimit_IsAtLimitPreserved(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// below limit: actual < max
	below := rent.AbsoluteRentLimit{
		MaxRentBP:    30,
		ActualRentBP: 15,
		IsAtLimit:    false,
	}
	saved, err := m.CreateAbsoluteRentLimit(ctx, below)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListAbsoluteRentLimits(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.IsAtLimit {
				t.Error("IsAtLimit must be false when actual < max")
			}
			break
		}
	}
	if !found {
		t.Errorf("created absolute-rent-limit ID %q not found in list", saved.ID)
	}
}
