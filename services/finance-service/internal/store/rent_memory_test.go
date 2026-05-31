package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
)

// Fixtures drawn from Marx, Vol. III Ch. 37:
//
//   - Landowner  : Duke of Sutherland (historically cleared the Highland crofters)
//   - Capitalist : the tenant-farmer who advances capital and pays rent
//   - Parcel     : Highland grade-B soil (fertility 2 of 4), 50 ha
//   - GroundRent : 2400 LabourMinutes differential rent, period 7 years

func ch37Landowner(name string) rent.Landowner {
	return rent.Landowner{Name: name}
}

func ch37AgriculturalCapitalist(parcelID string) rent.AgriculturalCapitalist {
	return rent.AgriculturalCapitalist{
		CapitalAdvanced: 4800,
		LeaseParcelID:   parcelID,
	}
}

func ch37LandParcel(ownerID string) rent.LandParcel {
	return rent.LandParcel{
		OwnerID:        ownerID,
		FertilityGrade: 2,
		AreaHectares:   50,
		Location:       "Sutherland, Scotland",
	}
}

func ch37GroundRent(parcelID, capitalistID, ownerID string) rent.GroundRent {
	return rent.GroundRent{
		ParcelID:            parcelID,
		CapitalistID:        capitalistID,
		LandOwnerID:         ownerID,
		Form:                rent.RentFormDifferential,
		AmountLabourMinutes: 2400,
		PeriodYears:         7,
	}
}

// ── LandParcel: Create → Get round-trip ──────────────────────────────────────

func TestMemory_LandParcelRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	p := ch37LandParcel("duke-owner-001")

	created, err := m.CreateLandParcel(ctx, p)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.FertilityGrade != 2 {
		t.Errorf("FertilityGrade: got %d, want 2", created.FertilityGrade)
	}
	if created.AreaHectares != 50 {
		t.Errorf("AreaHectares: got %d, want 50", created.AreaHectares)
	}

	got, err := m.GetLandParcel(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

// ── LandParcel: Get with unknown ID → ErrNotFound ────────────────────────────

func TestMemory_GetLandParcelNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetLandParcel(context.Background(), rent.LandParcelID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

// ── LandParcel: Duplicate create → ErrAlreadyExists ──────────────────────────

func TestMemory_CreateLandParcelDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	p := ch37LandParcel("duke-owner-002")
	p.ID = rent.LandParcelID("5eed000000000000003701")

	if _, err := m.CreateLandParcel(ctx, p); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateLandParcel(ctx, p); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

// ── LandParcel: List on empty store → non-nil empty slice ────────────────────

func TestMemory_ListLandParcelsNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListLandParcels(context.Background())
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

// ── LandParcel: List after 2 creates → newest-first ──────────────────────────

func TestMemory_ListLandParcels_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	first := ch37LandParcel("duke-owner-A")
	second := ch37LandParcel("duke-owner-B")

	createdFirst, err := m.CreateLandParcel(ctx, first)
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	time.Sleep(time.Nanosecond)
	createdSecond, err := m.CreateLandParcel(ctx, second)
	if err != nil {
		t.Fatalf("create second: %v", err)
	}

	out, err := m.ListLandParcels(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("list len = %d, want 2", len(out))
	}
	if !createdSecond.CreatedAt.After(createdFirst.CreatedAt) {
		return // clocks tied; ordering assertion not possible
	}
	if out[0].ID != createdSecond.ID {
		t.Errorf("first item id = %s, want newest (%s)", out[0].ID, createdSecond.ID)
	}
}

// ── Landowner: Create → Get round-trip ───────────────────────────────────────

func TestMemory_LandownerRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	lo := ch37Landowner("Duke of Sutherland")

	created, err := m.CreateLandowner(ctx, lo)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.Name != "Duke of Sutherland" {
		t.Errorf("Name: got %q, want %q", created.Name, "Duke of Sutherland")
	}

	got, err := m.GetLandowner(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

// ── Landowner: Get with unknown ID → ErrNotFound ─────────────────────────────

func TestMemory_GetLandownerNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetLandowner(context.Background(), rent.LandownerID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

// ── AgriculturalCapitalist: Create → Get round-trip ──────────────────────────

func TestMemory_AgriculturalCapitalistRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	ac := ch37AgriculturalCapitalist("parcel-abc-001")

	created, err := m.CreateAgriculturalCapitalist(ctx, ac)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.CapitalAdvanced != 4800 {
		t.Errorf("CapitalAdvanced: got %d, want 4800", created.CapitalAdvanced)
	}

	got, err := m.GetAgriculturalCapitalist(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

// ── AgriculturalCapitalist: Get with unknown ID → ErrNotFound ────────────────

func TestMemory_GetAgriculturalCapitalistNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetAgriculturalCapitalist(context.Background(), rent.AgriculturalCapitalistID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

// ── GroundRent: Create then ListGroundRents returns the row ──────────────────

func TestMemory_GroundRentCreateAndList(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	g := ch37GroundRent("parcel-x", "tenant-farmer-01", "duke-sutherland-01")

	created, err := m.CreateGroundRent(ctx, g)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.AmountLabourMinutes != 2400 {
		t.Errorf("AmountLabourMinutes: got %d, want 2400", created.AmountLabourMinutes)
	}

	all, err := m.ListGroundRents(ctx)
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

// ── GroundRent: ListGroundRents on empty store → non-nil empty slice ──────────

func TestMemory_ListGroundRentsNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListGroundRents(context.Background())
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
