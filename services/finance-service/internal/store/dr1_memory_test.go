// Package store — in-memory tests for Vol. III Ch. 39 (Differential Rent I)
// store methods. Fixtures drawn from Marx, Vol. III Ch. 39, Table I.
package store

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// marxTableIInput returns raw (uncomputed) entries for Marx Table I:
// A=100, B=200, C=300, D=400 output-quarters; capital=250p; price=300p/qtr.
func marxTableIInput() []rent.DifferentialRentIEntry {
	return []rent.DifferentialRentIEntry{
		{Grade: rent.SoilGradeA, CapitalAdvancedBP: 250, OutputQuarters: 100, PriceOfProductionBP: 300},
		{Grade: rent.SoilGradeB, CapitalAdvancedBP: 250, OutputQuarters: 200, PriceOfProductionBP: 300},
		{Grade: rent.SoilGradeC, CapitalAdvancedBP: 250, OutputQuarters: 300, PriceOfProductionBP: 300},
		{Grade: rent.SoilGradeD, CapitalAdvancedBP: 250, OutputQuarters: 400, PriceOfProductionBP: 300},
	}
}

// marxTableIHeader returns a DifferentialRentITable header already computed
// for use in CreateDifferentialRentITable calls.
func marxTableIHeader(desc string) rent.DifferentialRentITable {
	return rent.DifferentialRentITable{
		Description:     desc,
		TotalRentalBP:   1800,
		RegulatingGrade: rent.SoilGradeA,
	}
}

// ── CreateDifferentialRentITable / GetDifferentialRentITable ─────────────────

func TestMemory_DR1CreateAndGet(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	tbl := marxTableIHeader("Marx Table I — Ch. 39")
	entries := marxTableIInput()

	savedTbl, savedEntries, err := m.CreateDifferentialRentITable(ctx, tbl, entries)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if savedTbl.ID == "" {
		t.Fatal("expected table ID to be assigned")
	}
	if savedTbl.CreatedAt.IsZero() {
		t.Error("expected table CreatedAt to be assigned")
	}
	if len(savedEntries) != 4 {
		t.Fatalf("len(savedEntries) = %d, want 4", len(savedEntries))
	}
	for i, e := range savedEntries {
		if e.ID == "" {
			t.Errorf("entries[%d].ID not assigned", i)
		}
		if e.CreatedAt.IsZero() {
			t.Errorf("entries[%d].CreatedAt not assigned", i)
		}
		if e.TableID != string(savedTbl.ID) {
			t.Errorf("entries[%d].TableID = %q, want %q", i, e.TableID, string(savedTbl.ID))
		}
	}

	// Round-trip via Get.
	gotTbl, gotEntries, err := m.GetDifferentialRentITable(ctx, savedTbl.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if gotTbl.ID != savedTbl.ID {
		t.Errorf("get table ID = %q, want %q", gotTbl.ID, savedTbl.ID)
	}
	if gotTbl.TotalRentalBP != 1800 {
		t.Errorf("TotalRentalBP = %d, want 1800", gotTbl.TotalRentalBP)
	}
	if gotTbl.RegulatingGrade != rent.SoilGradeA {
		t.Errorf("RegulatingGrade = %v, want SoilGradeA", gotTbl.RegulatingGrade)
	}
	if len(gotEntries) != len(savedEntries) {
		t.Fatalf("get entries len = %d, want %d", len(gotEntries), len(savedEntries))
	}
	// Verify entry ID preservation.
	savedIDs := make(map[rent.DifferentialRentIEntryID]struct{}, len(savedEntries))
	for _, e := range savedEntries {
		savedIDs[e.ID] = struct{}{}
	}
	for _, e := range gotEntries {
		if _, ok := savedIDs[e.ID]; !ok {
			t.Errorf("got entry ID %q not in saved entries", e.ID)
		}
	}
}

func TestMemory_DR1GetNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	_, _, err := m.GetDifferentialRentITable(ctx, rent.DifferentialRentITableID("unknown-id"))
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

// ── ListDifferentialRentITables ───────────────────────────────────────────────

func TestMemory_DR1ListTablesContainsCreated(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, _, err := m.CreateDifferentialRentITable(ctx, marxTableIHeader("list-test"), marxTableIInput())
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	items, err := m.ListDifferentialRentITables(ctx)
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
		t.Errorf("created table ID %q not found in list", saved.ID)
	}
}

func TestMemory_DR1ListTablesEmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	items, err := m.ListDifferentialRentITables(context.Background())
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

// ── ListDifferentialRentIEntries ──────────────────────────────────────────────

func TestMemory_DR1ListEntriesRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	saved, savedEntries, err := m.CreateDifferentialRentITable(ctx, marxTableIHeader("entries-test"), marxTableIInput())
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	listed, err := m.ListDifferentialRentIEntries(ctx, saved.ID)
	if err != nil {
		t.Fatalf("list entries: %v", err)
	}
	if len(listed) != len(savedEntries) {
		t.Fatalf("listed %d entries, want %d", len(listed), len(savedEntries))
	}
}

func TestMemory_DR1ListEntriesUnknownTable(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	_, err := m.ListDifferentialRentIEntries(context.Background(), rent.DifferentialRentITableID("no-such-table"))
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

// ── CreateLocationRentFactor / ListLocationRentFactors ────────────────────────

func TestMemory_LocationRentFactorCreateAndList(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	f := rent.LocationRentFactor{
		ParcelID:         "suffolk-parcel-001",
		TransportCostBP:  120,
		RentEquivalentBP: 80,
	}

	saved, err := m.CreateLocationRentFactor(ctx, f)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if saved.ID == "" {
		t.Fatal("expected ID to be assigned")
	}
	if saved.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be assigned")
	}
	if saved.ParcelID != "suffolk-parcel-001" {
		t.Errorf("ParcelID = %q, want %q", saved.ParcelID, "suffolk-parcel-001")
	}
	if saved.TransportCostBP != 120 {
		t.Errorf("TransportCostBP = %d, want 120", saved.TransportCostBP)
	}
	if saved.RentEquivalentBP != 80 {
		t.Errorf("RentEquivalentBP = %d, want 80", saved.RentEquivalentBP)
	}

	items, err := m.ListLocationRentFactors(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("list len = %d, want 1", len(items))
	}
	if items[0] != saved {
		t.Errorf("listed row mismatch:\n got  %+v\n want %+v", items[0], saved)
	}
}

func TestMemory_LocationRentFactorListEmptyIsNonNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	items, err := m.ListLocationRentFactors(context.Background())
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
