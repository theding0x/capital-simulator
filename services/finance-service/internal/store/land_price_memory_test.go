// Package store — in-memory tests for Vol. III Ch. 46 (Building Site Rent,
// Rent in Mining, Price of Land) store methods.
package store

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// ── BuildingSiteRent ──────────────────────────────────────────────────────────

func TestMemory_Ch46_CreateBuildingSiteRent_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	rec := rent.BuildingSiteRent{
		ParcelID:                "london-mayfair",
		LocationGrade:           5,
		AnnualRentLabourMinutes: 300,
		LeaseTermYears:          99,
	}
	saved, err := m.CreateBuildingSiteRent(ctx, rec)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.ParcelID != "london-mayfair" {
		t.Errorf("ParcelID = %q, want london-mayfair", saved.ParcelID)
	}
	if saved.LocationGrade != 5 {
		t.Errorf("LocationGrade = %d, want 5", saved.LocationGrade)
	}
	if saved.AnnualRentLabourMinutes != 300 {
		t.Errorf("AnnualRentLabourMinutes = %d, want 300", saved.AnnualRentLabourMinutes)
	}
	if saved.LeaseTermYears != 99 {
		t.Errorf("LeaseTermYears = %d, want 99", saved.LeaseTermYears)
	}
}

func TestMemory_Ch46_ListBuildingSiteRents_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListBuildingSiteRents(context.Background())
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

func TestMemory_Ch46_CreateBuildingSiteRent_ContainsInList(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateBuildingSiteRent(ctx, rent.BuildingSiteRent{
		ParcelID:       "city-of-london",
		LocationGrade:  3,
		LeaseTermYears: 21,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListBuildingSiteRents(ctx)
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
		t.Errorf("created building-site-rent ID %q not found in list", saved.ID)
	}
}

func TestMemory_Ch46_CreateBuildingSiteRent_DuplicateID_ErrAlreadyExists(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateBuildingSiteRent(ctx, rent.BuildingSiteRent{
		ParcelID:       "dup-parcel",
		LocationGrade:  2,
		LeaseTermYears: 10,
	})
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	_, err = m.CreateBuildingSiteRent(ctx, saved)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}

// ── MiningRent ────────────────────────────────────────────────────────────────

func TestMemory_Ch46_CreateMiningRent_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Marx Ch. 46: best mine (ore grade 5) yields surplus profit above worst.
	rec := rent.MiningRent{
		ParcelID:                "cornwall-tin-mine",
		OreGrade:                5,
		AnnualRentLabourMinutes: 480,
		SurplusProfitBP:         1200,
	}
	saved, err := m.CreateMiningRent(ctx, rec)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.OreGrade != 5 {
		t.Errorf("OreGrade = %d, want 5", saved.OreGrade)
	}
	if saved.SurplusProfitBP != 1200 {
		t.Errorf("SurplusProfitBP = %d, want 1200", saved.SurplusProfitBP)
	}
}

func TestMemory_Ch46_ListMiningRents_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListMiningRents(context.Background())
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

func TestMemory_Ch46_CreateMiningRent_ContainsInList(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateMiningRent(ctx, rent.MiningRent{
		ParcelID: "durham-coal",
		OreGrade: 3,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListMiningRents(ctx)
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
		t.Errorf("created mining-rent ID %q not found in list", saved.ID)
	}
}

// ── MonopolyPriceRent ─────────────────────────────────────────────────────────

func TestMemory_Ch46_CreateMonopolyPriceRent_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Marx Ch. 46: burgundy grand cru, value=200, sale=400, rent=200.
	rec := rent.MonopolyPriceRent{
		ParcelID:                       "burgundy-clos-de-vougeot",
		ProductKind:                    "burgundy_grand_cru",
		ValueLabourMinutes:             200,
		MonopolySalePriceLabourMinutes: 400,
		MonopolyRentLabourMinutes:      200,
	}
	saved, err := m.CreateMonopolyPriceRent(ctx, rec)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.ProductKind != "burgundy_grand_cru" {
		t.Errorf("ProductKind = %q, want burgundy_grand_cru", saved.ProductKind)
	}
	if saved.MonopolyRentLabourMinutes != 200 {
		t.Errorf("MonopolyRentLabourMinutes = %d, want 200", saved.MonopolyRentLabourMinutes)
	}
}

func TestMemory_Ch46_ListMonopolyPriceRents_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListMonopolyPriceRents(context.Background())
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

func TestMemory_Ch46_CreateMonopolyPriceRent_ContainsInList(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateMonopolyPriceRent(ctx, rent.MonopolyPriceRent{
		ParcelID:    "rhine-riesling",
		ProductKind: "rhine_vineyard",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListMonopolyPriceRents(ctx)
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
		t.Errorf("created monopoly-price-rent ID %q not found in list", saved.ID)
	}
}

// ── LandPriceScenario ─────────────────────────────────────────────────────────

func TestMemory_Ch46_CreateLandPriceScenario_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Marx Ch. 46 scenario: falling interest (kind 1), rent=10 bp, rate
	// falls from 500 bp to 250 bp → new land price = 400 bp.
	rec := rent.LandPriceScenario{
		Kind:              rent.LandPriceFallingInterestRate,
		ParcelID:          "english-farm-1865",
		OldAnnualRentBP:   10,
		NewAnnualRentBP:   10,
		OldInterestRateBP: 500,
		NewInterestRateBP: 250,
		OldLandPriceBP:    200,
		NewLandPriceBP:    400,
	}
	saved, err := m.CreateLandPriceScenario(ctx, rec)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.Kind != rent.LandPriceFallingInterestRate {
		t.Errorf("Kind = %d, want %d (LandPriceFallingInterestRate)",
			saved.Kind, rent.LandPriceFallingInterestRate)
	}
	if saved.NewLandPriceBP != 400 {
		t.Errorf("NewLandPriceBP = %d, want 400", saved.NewLandPriceBP)
	}
}

func TestMemory_Ch46_ListLandPriceScenarios_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListLandPriceScenarios(context.Background())
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

func TestMemory_Ch46_CreateLandPriceScenario_KindPreservedInList(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Kind 3: product price rises → rent rises → price rises.
	rec := rent.LandPriceScenario{
		Kind:              rent.LandPriceRentRisesWithProduct,
		ParcelID:          "wheat-farm-1870",
		OldInterestRateBP: 500,
		NewInterestRateBP: 500,
		OldAnnualRentBP:   10,
		NewAnnualRentBP:   20,
		OldLandPriceBP:    200,
		NewLandPriceBP:    400,
	}
	saved, err := m.CreateLandPriceScenario(ctx, rec)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListLandPriceScenarios(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.Kind != rent.LandPriceRentRisesWithProduct {
				t.Errorf("list Kind = %d, want %d (LandPriceRentRisesWithProduct)",
					item.Kind, rent.LandPriceRentRisesWithProduct)
			}
			break
		}
	}
	if !found {
		t.Errorf("created land-price-scenario ID %q not found in list", saved.ID)
	}
}
