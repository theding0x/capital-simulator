package store_test

import (
	"context"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// TestMemoryUpsertAggregateShares_UpdatesSeededRow verifies that writing shares
// onto the seeded 1871 aggregate updates the share columns while preserving the
// other magnitudes (issue #215).
func TestMemoryUpsertAggregateShares_UpdatesSeededRow(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	before, err := mem.SnapshotAggregate(ctx, "1871")
	if err != nil {
		t.Fatalf("SnapshotAggregate(1871): %v", err)
	}
	if before.DepartmentIShareBP != 0 || before.DepartmentIIShareBP != 0 {
		t.Fatalf("seeded shares should start at 0/0, got %d/%d",
			before.DepartmentIShareBP, before.DepartmentIIShareBP)
	}

	if _, err := mem.UpsertAggregateShares(ctx, "1871", 6667, 3333); err != nil {
		t.Fatalf("UpsertAggregateShares: %v", err)
	}

	after, err := mem.SnapshotAggregate(ctx, "1871")
	if err != nil {
		t.Fatalf("SnapshotAggregate(1871) after upsert: %v", err)
	}
	if after.DepartmentIShareBP != 6667 || after.DepartmentIIShareBP != 3333 {
		t.Fatalf("shares not updated: got %d/%d, want 6667/3333",
			after.DepartmentIShareBP, after.DepartmentIIShareBP)
	}
	// The share write must not clobber the other seeded magnitudes.
	if after.TotalAdvancedPence != before.TotalAdvancedPence {
		t.Fatalf("TotalAdvancedPence clobbered: got %d, want %d",
			after.TotalAdvancedPence, before.TotalAdvancedPence)
	}
}

// TestMemoryUpsertAggregateShares_CreatesMissingRow verifies that projecting
// shares for a period with no aggregate row creates one and is readable back
// via SnapshotAggregate (issue #215).
func TestMemoryUpsertAggregateShares_CreatesMissingRow(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	ctx := context.Background()

	const period = "1865"
	got, err := mem.UpsertAggregateShares(ctx, period, 6667, 3333)
	if err != nil {
		t.Fatalf("UpsertAggregateShares(%s): %v", period, err)
	}
	if got.Period != period || got.DepartmentIShareBP != 6667 || got.DepartmentIIShareBP != 3333 {
		t.Fatalf("created row mismatch: %+v", got)
	}

	round, err := mem.SnapshotAggregate(ctx, period)
	if err != nil {
		t.Fatalf("SnapshotAggregate(%s): %v", period, err)
	}
	if round.DepartmentIShareBP != 6667 || round.DepartmentIIShareBP != 3333 {
		t.Fatalf("round-trip shares: got %d/%d, want 6667/3333",
			round.DepartmentIShareBP, round.DepartmentIIShareBP)
	}
}
