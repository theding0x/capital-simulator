package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// ch34NoteIssueConstraint returns a valid NoteIssueConstraint built from the
// Bank Act 1844 founding fixture (Marx, Vol. III Ch. 34):
// goldBacking=1 400 000 000, unbackedLimit=1 400 000 000 → MaxNotes=2 800 000 000.
func ch34NoteIssueConstraint(t *testing.T, name string) credit.NoteIssueConstraint {
	t.Helper()
	n, err := credit.NewNoteIssueConstraint(
		name,
		credit.RegimeAct1844,
		1_400_000_000,
		1_400_000_000,
		0,
		"Currency School doctrine institutionalised.",
	)
	if err != nil {
		t.Fatalf("ch34NoteIssueConstraint: %v", err)
	}
	return n
}

// ── Create → Get round-trip ──────────────────────────────────────────────────

func TestMemory_NoteIssueConstraintRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	n := ch34NoteIssueConstraint(t, "Bank Act 1844 — founding regime")

	created, err := m.CreateNoteIssueConstraint(ctx, n)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	// Domain invariants survive persistence.
	if created.Regime != credit.RegimeAct1844 {
		t.Errorf("Regime: got %d, want RegimeAct1844 (%d)", created.Regime, credit.RegimeAct1844)
	}
	if created.MaxNotes != 2_800_000_000 {
		t.Errorf("MaxNotes: got %d, want 2800000000", created.MaxNotes)
	}
	if created.NotesBeyondMax != 0 {
		t.Errorf("NotesBeyondMax: got %d, want 0", created.NotesBeyondMax)
	}

	got, err := m.GetNoteIssueConstraint(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

// ── Get with unknown ID → ErrNotFound ───────────────────────────────────────

func TestMemory_GetNoteIssueConstraintNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetNoteIssueConstraint(context.Background(), credit.NoteIssueConstraintID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

// ── Duplicate create → ErrAlreadyExists ─────────────────────────────────────

func TestMemory_CreateNoteIssueConstraintDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	n := ch34NoteIssueConstraint(t, "Bank Act 1844 duplicate test")
	n.ID = credit.NoteIssueConstraintID("5eed000000000000003400")

	if _, err := m.CreateNoteIssueConstraint(ctx, n); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateNoteIssueConstraint(ctx, n); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

// ── List fresh → non-nil len 0 ───────────────────────────────────────────────

func TestMemory_ListNoteIssueConstraintsNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListNoteIssueConstraints(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
	if len(out) != 0 {
		t.Errorf("list len = %d, want 0 on empty store", len(out))
	}
}

// ── List after 2 creates → len 2, newest-first ───────────────────────────────

func TestMemory_ListNoteIssueConstraints_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	first := ch34NoteIssueConstraint(t, "Bank Act 1844 — founding regime")
	second, err := credit.NewNoteIssueConstraint(
		"Crisis of 1857 — Act suspended",
		credit.RegimeSuspended,
		1_400_000_000,
		1_400_000_000,
		48_883_000,
		"Treasury letter suspends the Bank Act.",
	)
	if err != nil {
		t.Fatalf("second fixture: %v", err)
	}

	createdFirst, err := m.CreateNoteIssueConstraint(ctx, first)
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	// Ensure a distinct timestamp for the ordering assertion.
	time.Sleep(time.Nanosecond)
	createdSecond, err := m.CreateNoteIssueConstraint(ctx, second)
	if err != nil {
		t.Fatalf("create second: %v", err)
	}

	out, err := m.ListNoteIssueConstraints(ctx)
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
