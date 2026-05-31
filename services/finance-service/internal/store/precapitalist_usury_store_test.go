package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// ch36UsurersCapital returns a valid UsurersCapital using the Ancient Rome
// fixture (Marx, Vol. III Ch. 36): StageAntiquity, no wage labour,
// borrowers = ["plebeians","small producers","nobles"], rate = 4800 bp.
func ch36UsurersCapital(t *testing.T, name string) credit.UsurersCapital {
	t.Helper()
	u, err := credit.NewUsurersCapital(
		name,
		credit.StageAntiquity,
		false,
		[]string{"plebeians", "small producers", "nobles"},
		4800,
		"",
		"Usury in ancient Rome at rates of 48% p.a.",
	)
	if err != nil {
		t.Fatalf("ch36UsurersCapital: %v", err)
	}
	return u
}

// ── Create → Get round-trip ───────────────────────────────────────────────────

func TestMemory_UsurersCapitalRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	u := ch36UsurersCapital(t, "Ancient Roman Usury")

	created, err := m.CreateUsurersCapital(ctx, u)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.Name != "Ancient Roman Usury" {
		t.Errorf("Name: got %q, want %q", created.Name, "Ancient Roman Usury")
	}
	if created.Stage != credit.StageAntiquity {
		t.Errorf("Stage: got %d, want StageAntiquity", created.Stage)
	}
	if created.WageLabour {
		t.Error("WageLabour: got true, want false")
	}
	if created.InterestRateBP != 4800 {
		t.Errorf("InterestRateBP: got %d, want 4800", created.InterestRateBP)
	}
	if created.SubordinatedTo != "" {
		t.Errorf("SubordinatedTo: got %q, want empty", created.SubordinatedTo)
	}

	got, err := m.GetUsurersCapital(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID mismatch: got %s, want %s", got.ID, created.ID)
	}
	if got.Name != created.Name {
		t.Errorf("Name mismatch: got %q, want %q", got.Name, created.Name)
	}
	if got.Stage != created.Stage {
		t.Errorf("Stage mismatch: got %d, want %d", got.Stage, created.Stage)
	}
	if got.WageLabour != created.WageLabour {
		t.Errorf("WageLabour mismatch: got %v, want %v", got.WageLabour, created.WageLabour)
	}
	if got.InterestRateBP != created.InterestRateBP {
		t.Errorf("InterestRateBP mismatch: got %d, want %d", got.InterestRateBP, created.InterestRateBP)
	}
	if got.SubordinatedTo != created.SubordinatedTo {
		t.Errorf("SubordinatedTo mismatch: got %q, want %q", got.SubordinatedTo, created.SubordinatedTo)
	}
	if got.Description != created.Description {
		t.Errorf("Description mismatch: got %q, want %q", got.Description, created.Description)
	}
	if !got.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("CreatedAt mismatch: got %v, want %v", got.CreatedAt, created.CreatedAt)
	}
}

// ── Borrowers round-trip ──────────────────────────────────────────────────────

func TestMemory_UsurersCapital_BorrowersRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	u := ch36UsurersCapital(t, "Rome borrowers test")
	created, err := m.CreateUsurersCapital(ctx, u)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := m.GetUsurersCapital(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(got.Borrowers) != len(created.Borrowers) {
		t.Fatalf("Borrowers length: got %d, want %d", len(got.Borrowers), len(created.Borrowers))
	}
	for i, b := range created.Borrowers {
		if got.Borrowers[i] != b {
			t.Errorf("Borrowers[%d]: got %q, want %q", i, got.Borrowers[i], b)
		}
	}
}

func TestMemory_UsurersCapital_NilBorrowersRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Create with nil Borrowers — store must return non-nil empty slice.
	u, err := credit.NewUsurersCapital(
		"nil-borrowers",
		credit.StageAntiquity,
		false,
		nil,
		100,
		"",
		"",
	)
	if err != nil {
		t.Fatalf("NewUsurersCapital: %v", err)
	}

	created, err := m.CreateUsurersCapital(ctx, u)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := m.GetUsurersCapital(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Borrowers == nil {
		t.Error("Borrowers: got nil, want non-nil empty []string{}")
	}
	if len(got.Borrowers) != 0 {
		t.Errorf("Borrowers length: got %d, want 0", len(got.Borrowers))
	}
}

// ── Get unknown ID → ErrNotFound ──────────────────────────────────────────────

func TestMemory_GetUsurersCapitalNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetUsurersCapital(context.Background(), credit.UsurersCapitalID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

// ── Duplicate create → ErrAlreadyExists ──────────────────────────────────────

func TestMemory_CreateUsurersCapitalDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	u := ch36UsurersCapital(t, "Rome duplicate test")
	u.ID = credit.UsurersCapitalID("5eed000000000000003600")

	if _, err := m.CreateUsurersCapital(ctx, u); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateUsurersCapital(ctx, u); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

// ── List fresh → non-nil len 0 ───────────────────────────────────────────────

func TestMemory_ListUsurersCapitalsNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListUsurersCapitals(context.Background())
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

// ── List after 2 creates → len 2, newest-first ───────────────────────────────

func TestMemory_ListUsurersCapitals_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	first := ch36UsurersCapital(t, "Ancient Roman Usury")
	second := ch36UsurersCapital(t, "Medieval Church Usury")

	createdFirst, err := m.CreateUsurersCapital(ctx, first)
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	// Ensure a distinct timestamp for the ordering assertion.
	time.Sleep(time.Nanosecond)
	createdSecond, err := m.CreateUsurersCapital(ctx, second)
	if err != nil {
		t.Fatalf("create second: %v", err)
	}

	out, err := m.ListUsurersCapitals(ctx)
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
