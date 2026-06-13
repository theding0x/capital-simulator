package store_test

// Vol. II Ch. 1 — negative-path tests for the money-circuit store guards.
//
// Guard 1 (ErrNoProductionPhase):  recording C′ after M—C without an
//   intervening P must be rejected.
// Guard 2 (ErrMagnitudeNotPreserved): recording P with a magnitude that
//   differs from the opening advance must be rejected.

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// newMoneyCircuitStore returns a fresh in-memory store for money-circuit tests.
func newMoneyCircuitStore(t *testing.T) *store.Memory {
	t.Helper()
	return store.NewMemory()
}

// openSpinningMillCircuit creates a MoneyCircuit whose advance is the Spinning
// Mill 1871 total: £372 MP + £50 wages = £422 = 42200p.
func openSpinningMillCircuit(t *testing.T, mem *store.Memory) circulation.MoneyCircuit {
	t.Helper()
	ctx := context.Background()
	mc, err := mem.CreateMoneyCircuit(ctx, circulation.MoneyCircuit{
		Advance: circulation.MoneyAdvance{Amount: 42200},
	})
	if err != nil {
		t.Fatalf("CreateMoneyCircuit: %v", err)
	}
	return mc
}

// advanceToPurchase records the M—C purchase phase on the given circuit.
func advanceToPurchase(t *testing.T, mem *store.Memory, id circulation.MoneyCircuitID) {
	t.Helper()
	ctx := context.Background()
	pp := circulation.PurchasePhase{
		LabourLeg: circulation.LabourLeg{Amount: 5000, LabourPowerHours: 1440},
		MeansLeg:  circulation.MeansLeg{Amount: 37200, MeansCapacityHours: 1440},
	}
	if _, err := mem.RecordPurchase(ctx, id, pp); err != nil {
		t.Fatalf("RecordPurchase: %v", err)
	}
}

// TestRecordCommodity_ErrNoProductionPhase drives guard 1:
// attempting to record C′ immediately after M—C (skipping P) must return
// ErrNoProductionPhase (Vol. II Ch. 1, §I: surplus-value arises only in P).
func TestRecordCommodity_ErrNoProductionPhase(t *testing.T) {
	t.Parallel()
	mem := newMoneyCircuitStore(t)
	ctx := context.Background()

	mc := openSpinningMillCircuit(t, mem)
	advanceToPurchase(t, mem, mc.ID)

	// Attempt to record C′ without an intervening P phase.
	// The circuit is at M-C; the next legal moment would be P, not C-prime.
	cc := circulation.CommodityCapital{
		ValueOriginal: 42200,
		ValueSurplus:  7800,
	}
	_, err := mem.RecordCommodity(ctx, mc.ID, cc)
	if !errors.Is(err, circulation.ErrNoProductionPhase) && !errors.Is(err, circulation.ErrInvalidPhaseTransition) {
		t.Fatalf("expected ErrNoProductionPhase or ErrInvalidPhaseTransition, got %v", err)
	}
}

// TestRecordProductive_ErrMagnitudeNotPreserved drives guard 2:
// recording P with c+v != opening advance must return ErrMagnitudeNotPreserved
// (Vol. II Ch. 1, §I: "M is the same capital-value as P, only it has a
// different mode of existence.").
func TestRecordProductive_ErrMagnitudeNotPreserved(t *testing.T) {
	t.Parallel()
	mem := newMoneyCircuitStore(t)
	ctx := context.Background()

	mc := openSpinningMillCircuit(t, mem)
	advanceToPurchase(t, mem, mc.ID)

	// Productive state with c+v = 40000 ≠ 42200 advance.
	ps := circulation.ProductiveState{
		ConstantPence: 35000, // wrong: total 35000+5000=40000 ≠ 42200
		VariablePence: 5000,
	}
	_, err := mem.RecordProductive(ctx, mc.ID, ps)
	if !errors.Is(err, circulation.ErrMagnitudeNotPreserved) {
		t.Fatalf("expected ErrMagnitudeNotPreserved, got %v", err)
	}
}

// TestRecordProductive_HappyPath verifies the guards do not fire for a valid P
// recording (regression guard: the negative tests above must not block the
// well-formed path).
func TestRecordProductive_HappyPath(t *testing.T) {
	t.Parallel()
	mem := newMoneyCircuitStore(t)
	ctx := context.Background()

	mc := openSpinningMillCircuit(t, mem)
	advanceToPurchase(t, mem, mc.ID)

	// c + v = 37200 + 5000 = 42200 == advance; magnitude preserved.
	ps := circulation.ProductiveState{
		ConstantPence: 37200,
		VariablePence: 5000,
	}
	updated, err := mem.RecordProductive(ctx, mc.ID, ps)
	if err != nil {
		t.Fatalf("RecordProductive valid: %v", err)
	}
	if updated.Moment != circulation.MomentP {
		t.Errorf("moment = %q, want %q", updated.Moment, circulation.MomentP)
	}
}
