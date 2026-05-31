// Package store — in-memory tests for Vol. III Ch. 40 (Differential Rent II)
// store methods. Fixtures drawn from Marx, Vol. III Ch. 40: Soil-D
// doubled-investment (capital 250p, output 400 qtrs, surplus-profit 900p per tranche).
package store

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// ch40Investment returns a SuccessiveInvestment with Marx Ch. 40 Soil-D values.
func ch40Investment(parcelID string, tranche int) rent.SuccessiveInvestment {
	return rent.SuccessiveInvestment{
		ParcelID:        parcelID,
		TrancheNumber:   tranche,
		CapitalBP:       250,
		OutputQuarters:  400,
		SurplusProfitBP: 900,
	}
}

// ── CreateSuccessiveInvestment / ListSuccessiveInvestments ────────────────────

func TestMemory_DR2CreateSuccessiveInvestment_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	inv := ch40Investment("soil-D", 1)
	saved, err := m.CreateSuccessiveInvestment(ctx, inv)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.ParcelID != "soil-D" {
		t.Errorf("ParcelID = %q, want %q", saved.ParcelID, "soil-D")
	}
	if saved.CapitalBP != 250 {
		t.Errorf("CapitalBP = %d, want 250", saved.CapitalBP)
	}
	if saved.SurplusProfitBP != 900 {
		t.Errorf("SurplusProfitBP = %d, want 900", saved.SurplusProfitBP)
	}
}

func TestMemory_DR2ListSuccessiveInvestments_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	items, err := m.ListSuccessiveInvestments(context.Background())
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

func TestMemory_DR2ListSuccessiveInvestments_ContainsCreated(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, err := m.CreateSuccessiveInvestment(ctx, ch40Investment("soil-D", 1))
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListSuccessiveInvestments(ctx)
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
		t.Errorf("created investment ID %q not found in list", saved.ID)
	}
}

func TestMemory_DR2CreateSuccessiveInvestment_DuplicateID(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	inv := ch40Investment("soil-D", 1)
	saved, err := m.CreateSuccessiveInvestment(ctx, inv)
	if err != nil {
		t.Fatalf("first create: %v", err)
	}

	// Re-submit with the same ID that the store assigned.
	dup := saved
	_, err = m.CreateSuccessiveInvestment(ctx, dup)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("err = %v, want ErrAlreadyExists on duplicate ID", err)
	}
}

// ── ListSuccessiveInvestmentsByParcel ─────────────────────────────────────────

func TestMemory_DR2ListByParcel_FiltersByParcelID(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Two investments on parcel P1, one on P2.
	p1a, err := m.CreateSuccessiveInvestment(ctx, ch40Investment("P1", 1))
	if err != nil {
		t.Fatalf("create P1 tranche 1: %v", err)
	}
	p1b, err := m.CreateSuccessiveInvestment(ctx, ch40Investment("P1", 2))
	if err != nil {
		t.Fatalf("create P1 tranche 2: %v", err)
	}
	_, err = m.CreateSuccessiveInvestment(ctx, ch40Investment("P2", 1))
	if err != nil {
		t.Fatalf("create P2 tranche 1: %v", err)
	}

	items, err := m.ListSuccessiveInvestmentsByParcel(ctx, "P1")
	if err != nil {
		t.Fatalf("list by parcel P1: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("list len = %d, want 2 for parcel P1", len(items))
	}
	ids := map[rent.SuccessiveInvestmentID]struct{}{
		p1a.ID: {},
		p1b.ID: {},
	}
	for _, item := range items {
		if _, ok := ids[item.ID]; !ok {
			t.Errorf("unexpected ID %q in parcel-P1 list", item.ID)
		}
		if item.ParcelID != "P1" {
			t.Errorf("item.ParcelID = %q, want %q", item.ParcelID, "P1")
		}
	}
}

func TestMemory_DR2ListByParcel_UnknownParcelReturnsEmptyNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	items, err := m.ListSuccessiveInvestmentsByParcel(ctx, "unknown")
	if err != nil {
		t.Fatalf("unexpected error for unknown parcel: %v", err)
	}
	if items == nil {
		t.Error("list must return non-nil slice for unknown parcel, not nil")
	}
	if len(items) != 0 {
		t.Errorf("list len = %d, want 0 for unknown parcel", len(items))
	}
}

// ── CreateLeaseTerm / ListLeaseTerms ──────────────────────────────────────────

func TestMemory_DR2CreateLeaseTerm_RoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	lt := rent.LeaseTerm{
		ParcelID:      "soil-D",
		StartYear:     1867,
		DurationYears: 21,
		FixedRentBP:   900,
	}
	saved, err := m.CreateLeaseTerm(ctx, lt)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.ParcelID != "soil-D" {
		t.Errorf("ParcelID = %q, want %q", saved.ParcelID, "soil-D")
	}
	if saved.StartYear != 1867 {
		t.Errorf("StartYear = %d, want 1867", saved.StartYear)
	}
	if saved.DurationYears != 21 {
		t.Errorf("DurationYears = %d, want 21", saved.DurationYears)
	}
	if saved.FixedRentBP != 900 {
		t.Errorf("FixedRentBP = %d, want 900", saved.FixedRentBP)
	}
}

func TestMemory_DR2ListLeaseTerms_EmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	items, err := m.ListLeaseTerms(context.Background())
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

func TestMemory_DR2ListLeaseTerms_ContainsCreated(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	lt := rent.LeaseTerm{
		ParcelID:      "soil-D",
		StartYear:     1867,
		DurationYears: 21,
		FixedRentBP:   900,
	}
	saved, err := m.CreateLeaseTerm(ctx, lt)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListLeaseTerms(ctx)
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
		t.Errorf("created lease-term ID %q not found in list", saved.ID)
	}
}
