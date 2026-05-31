package store

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// Fixtures drawn from Marx, Vol. III Ch. 38:
//
//   - Waterfall capital: general PoP 11500 bp, individual 9000 bp → surplus 2500 bp
//   - Capitalised rent: £10 annual rent (4800 LM) at 5% (500 bp) → 96000 LM land price

// ── PriceOfProductionSurplusProfit ───────────────────────────────────────────

func ch38PoPSurplusProfit(capitalID string) rent.PriceOfProductionSurplusProfit {
	return rent.PriceOfProductionSurplusProfit{
		CapitalID:                     capitalID,
		GeneralPriceOfProductionBP:    11500,
		IndividualPriceOfProductionBP: 9000,
		SurplusProfitBP:               2500,
	}
}

func TestMemory_PoPSurplusProfitCreateAndList(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	s := ch38PoPSurplusProfit("waterfall-capital-001")

	created, err := m.CreatePoPSurplusProfit(ctx, s)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.CapitalID != "waterfall-capital-001" {
		t.Errorf("CapitalID: got %q, want %q", created.CapitalID, "waterfall-capital-001")
	}
	if created.GeneralPriceOfProductionBP != 11500 {
		t.Errorf("GeneralPriceOfProductionBP: got %d, want 11500", created.GeneralPriceOfProductionBP)
	}
	if created.IndividualPriceOfProductionBP != 9000 {
		t.Errorf("IndividualPriceOfProductionBP: got %d, want 9000", created.IndividualPriceOfProductionBP)
	}
	if created.SurplusProfitBP != 2500 {
		t.Errorf("SurplusProfitBP: got %d, want 2500", created.SurplusProfitBP)
	}

	all, err := m.ListPoPSurplusProfits(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("list len = %d, want 1", len(all))
	}
	if all[0] != created {
		t.Errorf("listed row mismatch:\n got  %+v\n want %+v", all[0], created)
	}
}

func TestMemory_ListPoPSurplusProfitsNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListPoPSurplusProfits(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list must return a non-nil slice even when empty")
	}
	if len(out) != 0 {
		t.Errorf("list len = %d, want 0 on empty store", len(out))
	}
}

func TestMemory_CreatePoPSurplusProfitDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	s := ch38PoPSurplusProfit("waterfall-capital-002")
	s.ID = rent.PriceOfProductionSurplusProfitID("5eed000000000000003801")

	if _, err := m.CreatePoPSurplusProfit(ctx, s); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreatePoPSurplusProfit(ctx, s); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

// ── MonopolisedNaturalForce ───────────────────────────────────────────────────

func ch38MonopolisedForce(parcelID string) rent.MonopolisedNaturalForce {
	return rent.MonopolisedNaturalForce{
		Kind:            "waterfall",
		IsMonopolisable: true,
		ParcelID:        parcelID,
	}
}

func TestMemory_MonopolisedNaturalForceCreateAndList(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	f := ch38MonopolisedForce("highland-parcel-001")

	created, err := m.CreateMonopolisedNaturalForce(ctx, f)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.Kind != "waterfall" {
		t.Errorf("Kind: got %q, want %q", created.Kind, "waterfall")
	}
	if !created.IsMonopolisable {
		t.Error("IsMonopolisable: got false, want true")
	}
	if created.ParcelID != "highland-parcel-001" {
		t.Errorf("ParcelID: got %q, want %q", created.ParcelID, "highland-parcel-001")
	}

	all, err := m.ListMonopolisedNaturalForces(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("list len = %d, want 1", len(all))
	}
	if all[0] != created {
		t.Errorf("listed row mismatch:\n got  %+v\n want %+v", all[0], created)
	}
}

func TestMemory_ListMonopolisedNaturalForcesNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListMonopolisedNaturalForces(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list must return a non-nil slice even when empty")
	}
	if len(out) != 0 {
		t.Errorf("list len = %d, want 0 on empty store", len(out))
	}
}

// ── DifferentialRent ─────────────────────────────────────────────────────────

func ch38DifferentialRent(parcelID string) rent.DifferentialRent {
	return rent.DifferentialRent{
		ParcelID:                parcelID,
		SurplusProfitBP:         2500,
		AnnualRentLabourMinutes: 4800,
	}
}

func TestMemory_DifferentialRentCreateAndList(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	dr := ch38DifferentialRent("highland-parcel-002")

	created, err := m.CreateDifferentialRent(ctx, dr)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.SurplusProfitBP != 2500 {
		t.Errorf("SurplusProfitBP: got %d, want 2500", created.SurplusProfitBP)
	}
	if created.AnnualRentLabourMinutes != 4800 {
		t.Errorf("AnnualRentLabourMinutes: got %d, want 4800", created.AnnualRentLabourMinutes)
	}

	all, err := m.ListDifferentialRents(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("list len = %d, want 1", len(all))
	}
	if all[0] != created {
		t.Errorf("listed row mismatch:\n got  %+v\n want %+v", all[0], created)
	}
}

func TestMemory_ListDifferentialRentsNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListDifferentialRents(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list must return a non-nil slice even when empty")
	}
	if len(out) != 0 {
		t.Errorf("list len = %d, want 0 on empty store", len(out))
	}
}

// ── CapitalisedRentPrice ──────────────────────────────────────────────────────

func ch38CapitalisedRentPrice(parcelID string) rent.CapitalisedRentPrice {
	return rent.CapitalisedRentPrice{
		ParcelID:                      parcelID,
		AnnualRentLabourMinutes:       4800,
		InterestRateBP:                500,
		CapitalisedPriceLabourMinutes: 96000,
	}
}

func TestMemory_CapitalisedRentPriceCreateAndList(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	crp := ch38CapitalisedRentPrice("highland-parcel-003")

	created, err := m.CreateCapitalisedRentPrice(ctx, crp)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.AnnualRentLabourMinutes != 4800 {
		t.Errorf("AnnualRentLabourMinutes: got %d, want 4800", created.AnnualRentLabourMinutes)
	}
	if created.InterestRateBP != 500 {
		t.Errorf("InterestRateBP: got %d, want 500", created.InterestRateBP)
	}
	if created.CapitalisedPriceLabourMinutes != 96000 {
		t.Errorf("CapitalisedPriceLabourMinutes: got %d, want 96000", created.CapitalisedPriceLabourMinutes)
	}

	all, err := m.ListCapitalisedRentPrices(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("list len = %d, want 1", len(all))
	}
	if all[0] != created {
		t.Errorf("listed row mismatch:\n got  %+v\n want %+v", all[0], created)
	}
}

func TestMemory_ListCapitalisedRentPricesNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListCapitalisedRentPrices(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list must return a non-nil slice even when empty")
	}
	if len(out) != 0 {
		t.Errorf("list len = %d, want 0 on empty store", len(out))
	}
}
