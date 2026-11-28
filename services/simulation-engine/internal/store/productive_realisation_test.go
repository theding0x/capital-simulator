package store_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Spinning Mill 1871 (Vol. II Ch. 2 §I): C' = 10,000 lbs of yarn = £500,
// £422 capital-value + £78 surplus-value (legacy pence 42,200 + 7,800).
const (
	millOriginal circulation.Pence = 42200
	millSurplus  circulation.Pence = 7800
	millTotal    circulation.Pence = 50000
)

func newMillWithExit(t *testing.T) (*store.Memory, circulation.ProductiveCircuitID) {
	t.Helper()
	m, id := newCapitaliseCircuit(t, circulation.ReproductionSimple)
	if _, err := m.RecordExitCommodity(context.Background(), id, millOriginal, millSurplus); err != nil {
		t.Fatalf("record C': %v", err)
	}
	return m, id
}

func TestMemory_RecordExitCommodity_UnknownCircuit(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.RecordExitCommodity(context.Background(), "nope", millOriginal, millSurplus)
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestMemory_RecordExitCommodity_RoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m, id := newCapitaliseCircuit(t, circulation.ReproductionSimple)

	got, err := m.RecordExitCommodity(ctx, id, millOriginal, millSurplus)
	if err != nil {
		t.Fatalf("record C': %v", err)
	}
	if got.ProductiveCircuitID != id || got.ValueOriginal != millOriginal || got.ValueSurplus != millSurplus {
		t.Errorf("returned %+v", got)
	}
	pc, err := m.GetProductiveCircuit(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if pc.ExitCommodity == nil || *pc.ExitCommodity != got {
		t.Fatalf("ExitCommodity = %+v, want %+v", pc.ExitCommodity, got)
	}
	if pc.ExitCommodity.Total() != millTotal {
		t.Errorf("Total = %d, want %d", pc.ExitCommodity.Total(), millTotal)
	}
	if pc.ExitRealisation != nil {
		t.Errorf("ExitRealisation = %+v, want nil before any sale", pc.ExitRealisation)
	}
}

func TestMemory_RecordExitCommodity_SecondIsAlreadyExists(t *testing.T) {
	t.Parallel()
	m, id := newMillWithExit(t)
	_, err := m.RecordExitCommodity(context.Background(), id, millOriginal, millSurplus)
	if !errors.Is(err, store.ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_RecordExitCommodity_RejectsInvalidValues(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m, id := newCapitaliseCircuit(t, circulation.ReproductionSimple)
	if _, err := m.RecordExitCommodity(ctx, id, 0, millSurplus); err == nil {
		t.Fatal("zero value_original accepted")
	}
	if _, err := m.RecordExitCommodity(ctx, id, millOriginal, -1); err == nil {
		t.Fatal("negative value_surplus accepted")
	}
	pc, _ := m.GetProductiveCircuit(ctx, id)
	if pc.ExitCommodity != nil {
		t.Errorf("invalid C' was recorded: %+v", pc.ExitCommodity)
	}
}

func TestMemory_RecordExitRealisation_UnknownCircuit(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.RecordExitRealisation(context.Background(), "nope", millTotal, time.Time{})
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("got %v, want ErrNotFound", err)
	}
}

func TestMemory_RecordExitRealisation_WithoutCommodity(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m, id := newCapitaliseCircuit(t, circulation.ReproductionSimple)
	_, err := m.RecordExitRealisation(ctx, id, millTotal, time.Time{})
	if !errors.Is(err, circulation.ErrNoCommodityCapital) {
		t.Fatalf("got %v, want ErrNoCommodityCapital", err)
	}
	pc, _ := m.GetProductiveCircuit(ctx, id)
	if pc.ExitRealisation != nil {
		t.Errorf("rejected realisation was recorded: %+v", pc.ExitRealisation)
	}
}

// The 8,440 lbs that replace the capital sell, the 1,560 lbs of surplus yarn do
// not: the circuit stays at C' and no M' exists.
func TestMemory_RecordExitRealisation_PartialIsRejectedAndNotRecorded(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m, id := newMillWithExit(t)

	for _, realised := range []circulation.Pence{millOriginal, millTotal - 1} {
		_, err := m.RecordExitRealisation(ctx, id, realised, time.Time{})
		if !errors.Is(err, circulation.ErrPartialRealisation) {
			t.Fatalf("realised %d: got %v, want ErrPartialRealisation", realised, err)
		}
	}
	pc, err := m.GetProductiveCircuit(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if pc.ExitRealisation != nil {
		t.Errorf("ExitRealisation = %+v, want nil after rejected sales", pc.ExitRealisation)
	}
	// A rejected sale must not block the genuine one.
	if _, err := m.RecordExitRealisation(ctx, id, millTotal, time.Time{}); err != nil {
		t.Fatalf("full realisation after rejected partials: %v", err)
	}
}

func TestMemory_RecordExitRealisation_FullSaleRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m, id := newMillWithExit(t)
	soldAt := time.Date(1871, 1, 6, 0, 0, 0, 0, time.UTC)

	got, err := m.RecordExitRealisation(ctx, id, millTotal, soldAt)
	if err != nil {
		t.Fatalf("realise: %v", err)
	}
	want := circulation.ExitRealisation{ProductiveCircuitID: id, SoldAt: soldAt, RealisedPence: millTotal}
	if got != want {
		t.Errorf("returned %+v, want %+v", got, want)
	}
	pc, err := m.GetProductiveCircuit(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if pc.ExitRealisation == nil || *pc.ExitRealisation != want {
		t.Errorf("ExitRealisation = %+v, want %+v", pc.ExitRealisation, want)
	}
}

func TestMemory_RecordExitRealisation_DefaultsSoldAt(t *testing.T) {
	t.Parallel()
	m, id := newMillWithExit(t)
	got, err := m.RecordExitRealisation(context.Background(), id, millTotal, time.Time{})
	if err != nil {
		t.Fatalf("realise: %v", err)
	}
	if got.SoldAt.IsZero() {
		t.Error("SoldAt left zero, want defaulted to now")
	}
}

func TestMemory_RecordExitRealisation_OverRealisationAccepted(t *testing.T) {
	t.Parallel()
	m, id := newMillWithExit(t)
	got, err := m.RecordExitRealisation(context.Background(), id, 52000, time.Time{})
	if err != nil {
		t.Fatalf("realise above value: %v", err)
	}
	if got.RealisedPence != 52000 {
		t.Errorf("RealisedPence = %d, want 52000", got.RealisedPence)
	}
}

func TestMemory_RecordExitRealisation_SecondIsAlreadyExists(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m, id := newMillWithExit(t)
	if _, err := m.RecordExitRealisation(ctx, id, millTotal, time.Time{}); err != nil {
		t.Fatalf("first realisation: %v", err)
	}
	_, err := m.RecordExitRealisation(ctx, id, millTotal, time.Time{})
	if !errors.Is(err, store.ErrAlreadyExists) {
		t.Fatalf("got %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListProductiveCircuits_IncludesExitFields(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m, id := newMillWithExit(t)
	if _, err := m.RecordExitRealisation(ctx, id, millTotal, time.Time{}); err != nil {
		t.Fatalf("realise: %v", err)
	}
	list, err := m.ListProductiveCircuits(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("len = %d, want 1", len(list))
	}
	if list[0].ExitCommodity == nil || list[0].ExitCommodity.Total() != millTotal {
		t.Errorf("list ExitCommodity = %+v", list[0].ExitCommodity)
	}
	if list[0].ExitRealisation == nil || list[0].ExitRealisation.RealisedPence != millTotal {
		t.Errorf("list ExitRealisation = %+v", list[0].ExitRealisation)
	}
}

// Get must hand back a copy: mutating it cannot rewrite the stored C'.
func TestMemory_GetProductiveCircuit_ExitFieldsAreCopies(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m, id := newMillWithExit(t)
	pc, _ := m.GetProductiveCircuit(ctx, id)
	pc.ExitCommodity.ValueSurplus = 0

	again, _ := m.GetProductiveCircuit(ctx, id)
	if again.ExitCommodity.ValueSurplus != millSurplus {
		t.Errorf("stored ValueSurplus = %d after caller mutation, want %d", again.ExitCommodity.ValueSurplus, millSurplus)
	}
}
