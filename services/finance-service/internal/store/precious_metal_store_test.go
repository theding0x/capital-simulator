package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// ch35GoldReserve returns a valid GoldReserve for the BoE Nov 1857 crisis
// fixture (Marx, Vol. III Ch. 35): amount = 800 000 000 pence.
func ch35GoldReserve(t *testing.T, name string) credit.GoldReserve {
	t.Helper()
	g, err := credit.NewGoldReserve(name, 800_000_000, "BoE gold reserve during the 1857 crisis.")
	if err != nil {
		t.Fatalf("ch35GoldReserve: %v", err)
	}
	return g
}

// ch35RateOfExchange returns a valid RateOfExchange for the London-Hamburg 1847
// fixture (Marx, Vol. III Ch. 35): deviation = −150 bp below par.
func ch35RateOfExchange(name string) credit.RateOfExchange {
	return credit.NewRateOfExchange(name, -150, "Adverse balance from corn imports.")
}

// ── GoldReserve: Create → Get round-trip ─────────────────────────────────────

func TestMemory_GoldReserveRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	g := ch35GoldReserve(t, "BoE Nov 1857")

	created, err := m.CreateGoldReserve(ctx, g)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.Amount != 800_000_000 {
		t.Errorf("Amount: got %d, want 800000000", created.Amount)
	}
	if created.Name != "BoE Nov 1857" {
		t.Errorf("Name: got %q, want %q", created.Name, "BoE Nov 1857")
	}
	if created.Description != "BoE gold reserve during the 1857 crisis." {
		t.Errorf("Description mismatch: %q", created.Description)
	}

	got, err := m.GetGoldReserve(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

// ── GoldReserve: Get with unknown ID → ErrNotFound ───────────────────────────

func TestMemory_GetGoldReserveNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetGoldReserve(context.Background(), credit.GoldReserveID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

// ── GoldReserve: Duplicate create → ErrAlreadyExists ─────────────────────────

func TestMemory_CreateGoldReserveDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	g := ch35GoldReserve(t, "BoE Nov 1857 duplicate test")
	g.ID = credit.GoldReserveID("5eed000000000000003500")

	if _, err := m.CreateGoldReserve(ctx, g); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateGoldReserve(ctx, g); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

// ── GoldReserve: List fresh → non-nil len 0 ──────────────────────────────────

func TestMemory_ListGoldReservesNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListGoldReserves(context.Background())
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

// ── GoldReserve: List after 2 creates → len 2, newest-first ──────────────────

func TestMemory_ListGoldReserves_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	first := ch35GoldReserve(t, "BoE 1844 Baseline")
	second := ch35GoldReserve(t, "BoE Nov 1857")

	createdFirst, err := m.CreateGoldReserve(ctx, first)
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	// Ensure a distinct timestamp for the ordering assertion.
	time.Sleep(time.Nanosecond)
	createdSecond, err := m.CreateGoldReserve(ctx, second)
	if err != nil {
		t.Fatalf("create second: %v", err)
	}

	out, err := m.ListGoldReserves(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("list len = %d, want 2", len(out))
	}
	// Skip ordering assertion if both share the same nanosecond timestamp.
	if !createdSecond.CreatedAt.After(createdFirst.CreatedAt) {
		return
	}
	if out[0].ID != createdSecond.ID {
		t.Errorf("first item id = %s, want newest (%s)", out[0].ID, createdSecond.ID)
	}
}

// ── RateOfExchange: Create → Get round-trip ───────────────────────────────────

func TestMemory_RateOfExchangeRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	re := ch35RateOfExchange("London-Hamburg 1847")

	created, err := m.CreateRateOfExchange(ctx, re)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.DeviationBP != -150 {
		t.Errorf("DeviationBP: got %d, want -150", created.DeviationBP)
	}
	if created.Name != "London-Hamburg 1847" {
		t.Errorf("Name: got %q, want %q", created.Name, "London-Hamburg 1847")
	}
	if created.Description != "Adverse balance from corn imports." {
		t.Errorf("Description mismatch: %q", created.Description)
	}

	got, err := m.GetRateOfExchange(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

// ── RateOfExchange: Get with unknown ID → ErrNotFound ────────────────────────

func TestMemory_GetRateOfExchangeNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetRateOfExchange(context.Background(), credit.RateOfExchangeID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

// ── RateOfExchange: Duplicate create → ErrAlreadyExists ──────────────────────

func TestMemory_CreateRateOfExchangeDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	re := ch35RateOfExchange("London-Hamburg 1847 duplicate test")
	re.ID = credit.RateOfExchangeID("5eed000000000000003510")

	if _, err := m.CreateRateOfExchange(ctx, re); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateRateOfExchange(ctx, re); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

// ── RateOfExchange: List fresh → non-nil len 0 ───────────────────────────────

func TestMemory_ListRatesOfExchangeNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListRatesOfExchange(context.Background())
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

// ── RateOfExchange: List after 2 creates → len 2, newest-first ───────────────

func TestMemory_ListRatesOfExchange_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	first := ch35RateOfExchange("London-Hamburg 1847")
	second := credit.NewRateOfExchange("London-India 1857", -200, "Persistent adverse drain.")

	createdFirst, err := m.CreateRateOfExchange(ctx, first)
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	// Ensure a distinct timestamp for the ordering assertion.
	time.Sleep(time.Nanosecond)
	createdSecond, err := m.CreateRateOfExchange(ctx, second)
	if err != nil {
		t.Fatalf("create second: %v", err)
	}

	out, err := m.ListRatesOfExchange(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("list len = %d, want 2", len(out))
	}
	// Skip ordering assertion if both share the same nanosecond timestamp.
	if !createdSecond.CreatedAt.After(createdFirst.CreatedAt) {
		return
	}
	if out[0].ID != createdSecond.ID {
		t.Errorf("first item id = %s, want newest (%s)", out[0].ID, createdSecond.ID)
	}
}
