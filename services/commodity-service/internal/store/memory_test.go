package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)

// tickingClock returns a now func that advances by one second on each call,
// starting from the given base. Used to remove time.Now() resolution flake
// from tests that assert on timestamp ordering.
func tickingClock(base time.Time) func() time.Time {
	var n int64
	return func() time.Time {
		n++
		return base.Add(time.Duration(n) * time.Second)
	}
}

func makeCommodity(name, kind string, snlt commodity.LabourMinutes) commodity.Commodity {
	return commodity.Commodity{
		Name: name,
		UseValue: commodity.UseValue{
			Description: name + " for use",
			Unit:        "units",
		},
		ConcreteLabour: commodity.ConcreteLabour{Kind: kind},
		SNLTPerUnit:    snlt,
	}
}

func TestMemory_CreateGet(t *testing.T) {
	t.Parallel()
	s := NewMemory()
	ctx := context.Background()

	created, err := s.Create(ctx, makeCommodity("Linen", "weaving", 30))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID.IsZero() {
		t.Errorf("Create should populate ID")
	}
	if created.CreatedAt.IsZero() || created.UpdatedAt.IsZero() {
		t.Errorf("Create should populate timestamps")
	}

	got, err := s.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != "Linen" {
		t.Errorf("Get returned wrong name: %v", got.Name)
	}
}

func TestMemory_Get_NotFound(t *testing.T) {
	t.Parallel()
	s := NewMemory()
	_, err := s.Get(context.Background(), "nope")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Get(nope) = %v, want ErrNotFound", err)
	}
}

func TestMemory_Create_Duplicate(t *testing.T) {
	t.Parallel()
	s := NewMemory()
	ctx := context.Background()
	if _, err := s.Create(ctx, makeCommodity("linen", "weaving", 30)); err != nil {
		t.Fatalf("Create: %v", err)
	}
	_, err := s.Create(ctx, makeCommodity("LINEN", "weaving", 30))
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("Create duplicate (case-insensitive) = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_List_Sorted(t *testing.T) {
	t.Parallel()
	s := NewMemory()
	ctx := context.Background()
	for _, n := range []string{"iron", "Linen", "corn", "coat"} {
		if _, err := s.Create(ctx, makeCommodity(n, "x", 30)); err != nil {
			t.Fatalf("Create %s: %v", n, err)
		}
	}
	list, err := s.List(ctx)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 4 {
		t.Fatalf("List len = %d, want 4", len(list))
	}
	want := []string{"coat", "corn", "iron", "Linen"}
	for i, c := range list {
		if c.Name != want[i] {
			t.Errorf("List[%d] = %q, want %q", i, c.Name, want[i])
		}
	}
}

func TestMemory_Update_AppliesPartialFields(t *testing.T) {
	t.Parallel()
	s := NewMemory()
	s.now = tickingClock(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	ctx := context.Background()
	created, err := s.Create(ctx, makeCommodity("linen", "weaving", 30))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	newSNLT := commodity.LabourMinutes(15)
	updated, err := s.Update(ctx, created.ID, Update{SNLTPerUnit: &newSNLT})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.SNLTPerUnit != 15 {
		t.Errorf("SNLTPerUnit = %d, want 15", updated.SNLTPerUnit)
	}
	if updated.Name != "linen" {
		t.Errorf("Name should be unchanged, got %v", updated.Name)
	}
	if !updated.UpdatedAt.After(created.UpdatedAt) {
		t.Errorf("UpdatedAt should advance on Update (created=%s, updated=%s)",
			created.UpdatedAt, updated.UpdatedAt)
	}
}

func TestMemory_Update_RejectsRenameIntoExisting(t *testing.T) {
	t.Parallel()
	s := NewMemory()
	ctx := context.Background()
	if _, err := s.Create(ctx, makeCommodity("linen", "weaving", 30)); err != nil {
		t.Fatal(err)
	}
	coat, err := s.Create(ctx, makeCommodity("coat", "tailoring", 600))
	if err != nil {
		t.Fatal(err)
	}
	rename := "LINEN"
	_, err = s.Update(ctx, coat.ID, Update{Name: &rename})
	if !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("rename collision = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_Delete(t *testing.T) {
	t.Parallel()
	s := NewMemory()
	ctx := context.Background()
	c, err := s.Create(ctx, makeCommodity("linen", "weaving", 30))
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(ctx, c.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err = s.Get(ctx, c.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Get after delete = %v, want ErrNotFound", err)
	}
	if err := s.Delete(ctx, c.ID); !errors.Is(err, ErrNotFound) {
		t.Errorf("Delete twice = %v, want ErrNotFound", err)
	}
}
