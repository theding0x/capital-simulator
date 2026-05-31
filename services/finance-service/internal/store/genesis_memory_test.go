// Package store — in-memory tests for Vol. III Ch. 47 (Genesis of Capitalist
// Ground-Rent) store methods.
//
// Fixtures drawn from Marx, Capital Vol. III Ch. 47:
//   - RentHistoricalForm stage 1: labour-rent, Normandy 10th–13th c.
//   - HistoricalRentTransition: labour-rent → product-rent (stage 1→2).
//   - SmallPeasantProduction: conflated income; TotalIncomeBP = 2+1+3 = 6 bp.
package store

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// ── RentHistoricalForm ────────────────────────────────────────────────────────

func TestMemory_Ch47_CreateRentHistoricalForm_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Marx Ch. 47: labour-rent — the direct serf works part of the week on the
	// lord's land (Stage=1, Normandy, 10th–13th c).
	rec := rent.RentHistoricalForm{
		Stage:         rent.RentFormLabour,
		Name:          "labour-rent",
		CoercionKind:  "extra-economic",
		SurplusForm:   "direct_labour",
		ExampleRegion: "Normandy",
		ExamplePeriod: "10th-13th c",
	}
	saved, err := m.CreateRentHistoricalForm(ctx, rec)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.Stage != rent.RentFormLabour {
		t.Errorf("Stage = %d, want %d (RentFormLabour)", saved.Stage, rent.RentFormLabour)
	}
	if saved.Name != "labour-rent" {
		t.Errorf("Name = %q, want labour-rent", saved.Name)
	}
}

func TestMemory_Ch47_ListRentHistoricalForms_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListRentHistoricalForms(context.Background())
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

func TestMemory_Ch47_CreateRentHistoricalForm_ContainsInList(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateRentHistoricalForm(ctx, rent.RentHistoricalForm{
		Stage:         rent.RentFormProduct,
		Name:          "product-rent",
		CoercionKind:  "economic",
		SurplusForm:   "surplus_product",
		ExampleRegion: "medieval-France",
		ExamplePeriod: "13th-15th c",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListRentHistoricalForms(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.Stage != rent.RentFormProduct {
				t.Errorf("list Stage = %d, want %d (RentFormProduct)",
					item.Stage, rent.RentFormProduct)
			}
			if item.Name != "product-rent" {
				t.Errorf("list Name = %q, want product-rent", item.Name)
			}
			break
		}
	}
	if !found {
		t.Errorf("created rent-historical-form ID %q not found in list", saved.ID)
	}
}

func TestMemory_Ch47_CreateRentHistoricalForm_DuplicateID_ErrAlreadyExists(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateRentHistoricalForm(ctx, rent.RentHistoricalForm{
		Stage:         rent.RentFormMoney,
		Name:          "money-rent",
		CoercionKind:  "economic",
		SurplusForm:   "money_payment",
		ExampleRegion: "England",
		ExamplePeriod: "15th-16th c",
	})
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	_, err = m.CreateRentHistoricalForm(ctx, saved)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}

// ── HistoricalRentTransition ──────────────────────────────────────────────────

func TestMemory_Ch47_CreateHistoricalRentTransition_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Marx Ch. 47: transition from labour-rent (1) to product-rent (2) as
	// commodity exchange penetrates the village economy.
	rec := rent.HistoricalRentTransition{
		FromStage:     rent.RentFormLabour,
		ToStage:       rent.RentFormProduct,
		Preconditions: "development of commodity exchange",
		DrivingForce:  "rising productivity in peasant household",
	}
	saved, err := m.CreateHistoricalRentTransition(ctx, rec)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.FromStage != rent.RentFormLabour {
		t.Errorf("FromStage = %d, want %d (RentFormLabour)", saved.FromStage, rent.RentFormLabour)
	}
	if saved.ToStage != rent.RentFormProduct {
		t.Errorf("ToStage = %d, want %d (RentFormProduct)", saved.ToStage, rent.RentFormProduct)
	}
}

func TestMemory_Ch47_ListHistoricalRentTransitions_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListHistoricalRentTransitions(context.Background())
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

func TestMemory_Ch47_CreateHistoricalRentTransition_ContainsInList(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Transition product-rent (2) → money-rent (3): money-economy matures.
	saved, err := m.CreateHistoricalRentTransition(ctx, rent.HistoricalRentTransition{
		FromStage:     rent.RentFormProduct,
		ToStage:       rent.RentFormMoney,
		Preconditions: "money economy matures",
		DrivingForce:  "generalisation of commodity production",
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListHistoricalRentTransitions(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.FromStage != rent.RentFormProduct {
				t.Errorf("list FromStage = %d, want %d (RentFormProduct)",
					item.FromStage, rent.RentFormProduct)
			}
			if item.ToStage != rent.RentFormMoney {
				t.Errorf("list ToStage = %d, want %d (RentFormMoney)",
					item.ToStage, rent.RentFormMoney)
			}
			break
		}
	}
	if !found {
		t.Errorf("created historical-rent-transition ID %q not found in list", saved.ID)
	}
}

func TestMemory_Ch47_CreateHistoricalRentTransition_DuplicateID_ErrAlreadyExists(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateHistoricalRentTransition(ctx, rent.HistoricalRentTransition{
		FromStage:     rent.RentFormMoney,
		ToStage:       rent.RentFormCapitalist,
		Preconditions: "capitalist mode of production established",
		DrivingForce:  "separation of direct producers from land",
	})
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	_, err = m.CreateHistoricalRentTransition(ctx, saved)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}

// ── SmallPeasantProduction ────────────────────────────────────────────────────

func TestMemory_Ch47_CreateSmallPeasantProduction_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Marx Ch. 47: the small peasant proprietor who is simultaneously labourer,
	// capitalist, and landowner — rent=2, profit=1, wages=3 → total=6 bp.
	rec := rent.SmallPeasantProduction{
		ParcelID:          "normandy-peasant-plot",
		RentComponentBP:   2,
		ProfitComponentBP: 1,
		WagesComponentBP:  3,
		TotalIncomeBP:     6,
		IsConflated:       true,
	}
	saved, err := m.CreateSmallPeasantProduction(ctx, rec)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.IsConflated != true {
		t.Errorf("IsConflated = %v, want true", saved.IsConflated)
	}
	if saved.TotalIncomeBP != 6 {
		t.Errorf("TotalIncomeBP = %d, want 6", saved.TotalIncomeBP)
	}
}

func TestMemory_Ch47_ListSmallPeasantProductions_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()

	items, err := m.ListSmallPeasantProductions(context.Background())
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

func TestMemory_Ch47_CreateSmallPeasantProduction_ContainsInList(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateSmallPeasantProduction(ctx, rent.SmallPeasantProduction{
		ParcelID:          "english-copyhold-1600",
		RentComponentBP:   5,
		ProfitComponentBP: 3,
		WagesComponentBP:  2,
		TotalIncomeBP:     10,
		IsConflated:       true,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListSmallPeasantProductions(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	found := false
	for _, item := range items {
		if item.ID == saved.ID {
			found = true
			if item.IsConflated != true {
				t.Errorf("list IsConflated = %v, want true", item.IsConflated)
			}
			if item.TotalIncomeBP != 10 {
				t.Errorf("list TotalIncomeBP = %d, want 10", item.TotalIncomeBP)
			}
			break
		}
	}
	if !found {
		t.Errorf("created small-peasant-production ID %q not found in list", saved.ID)
	}
}

func TestMemory_Ch47_CreateSmallPeasantProduction_DuplicateID_ErrAlreadyExists(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateSmallPeasantProduction(ctx, rent.SmallPeasantProduction{
		ParcelID:          "dup-parcel",
		RentComponentBP:   1,
		ProfitComponentBP: 1,
		WagesComponentBP:  1,
		TotalIncomeBP:     3,
		IsConflated:       false,
	})
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	_, err = m.CreateSmallPeasantProduction(ctx, saved)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}
