package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// Spinning Mill 1871 (Vol. II Ch. 2): £372 constant + £50 variable advance,
// with the £78 surplus-value hoarded before being capitalised.
func newCapitaliseCircuit(t *testing.T, mode circulation.ReproductionMode) (*store.Memory, circulation.ProductiveCircuitID) {
	t.Helper()
	m := store.NewMemory()
	pc, err := m.CreateProductiveCircuit(context.Background(), circulation.ProductiveCircuit{
		ConstantPence: 37200, VariablePence: 5000, Mode: mode, MinCapitalisationIncrement: 100,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	return m, pc.ID
}

func TestMemory_Capitalise_DrawsDownHoardAndGrowsP(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m, id := newCapitaliseCircuit(t, circulation.ReproductionExtended)
	if _, err := m.Accumulate(ctx, id, 7800); err != nil {
		t.Fatalf("accumulate: %v", err)
	}
	if _, err := m.Capitalise(ctx, id, 7800, 6500, 1300); err != nil {
		t.Fatalf("capitalise: %v", err)
	}
	pc, err := m.GetProductiveCircuit(ctx, id)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if pc.LatentMoneyCapital.Accumulated != 0 {
		t.Errorf("hoard = %d, want 0", pc.LatentMoneyCapital.Accumulated)
	}
	if pc.ConstantPence != 43700 || pc.VariablePence != 6300 {
		t.Errorf("c/v = %d/%d, want 43700/6300", pc.ConstantPence, pc.VariablePence)
	}
}

func TestMemory_Capitalise_RejectsMoreThanHoard(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m, id := newCapitaliseCircuit(t, circulation.ReproductionExtended)
	if _, err := m.Accumulate(ctx, id, 7800); err != nil {
		t.Fatalf("accumulate: %v", err)
	}
	_, err := m.Capitalise(ctx, id, 7900, 7900, 0)
	if !errors.Is(err, circulation.ErrCapitalisationExceedsHoard) {
		t.Fatalf("got %v, want ErrCapitalisationExceedsHoard", err)
	}
	pc, _ := m.GetProductiveCircuit(ctx, id)
	if pc.LatentMoneyCapital.Accumulated != 7800 || pc.ConstantPence != 37200 {
		t.Errorf("state mutated on rejection: hoard=%d c=%d", pc.LatentMoneyCapital.Accumulated, pc.ConstantPence)
	}
}

func TestMemory_Capitalise_RejectsSimpleReproduction(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m, id := newCapitaliseCircuit(t, circulation.ReproductionSimple)
	if _, err := m.Accumulate(ctx, id, 7800); err != nil {
		t.Fatalf("accumulate: %v", err)
	}
	if _, err := m.Capitalise(ctx, id, 7800, 6500, 1300); !errors.Is(err, circulation.ErrSimpleReproductionNoGrowth) {
		t.Fatalf("got %v, want ErrSimpleReproductionNoGrowth", err)
	}
}

func TestMemory_Capitalise_RejectsNonConservedSplit(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m, id := newCapitaliseCircuit(t, circulation.ReproductionExtended)
	if _, err := m.Accumulate(ctx, id, 7800); err != nil {
		t.Fatalf("accumulate: %v", err)
	}
	if _, err := m.Capitalise(ctx, id, 7800, 6500, 1000); !errors.Is(err, circulation.ErrCapitalisationNotConserved) {
		t.Fatalf("got %v, want ErrCapitalisationNotConserved", err)
	}
}
