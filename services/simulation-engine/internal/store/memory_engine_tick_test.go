package store

import (
	"context"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
)

func TestMemory_RecordAndListEngineTicks(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	for i := 1; i <= 3; i++ {
		if _, err := m.RecordEngineTick(ctx, engine.ScheduledTick{
			Sequence:         int64(i),
			OccurredAt:       time.Date(2026, 5, 29, 12, i, 0, 0, time.UTC),
			TickersRun:       2,
			EntitiesAdvanced: i,
		}); err != nil {
			t.Fatalf("RecordEngineTick %d: %v", i, err)
		}
	}

	all, err := m.ListEngineTicks(ctx, 0)
	if err != nil {
		t.Fatalf("ListEngineTicks: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("len = %d, want 3", len(all))
	}
	// Newest first.
	if all[0].Sequence != 3 || all[2].Sequence != 1 {
		t.Fatalf("order = [%d..%d], want newest-first 3..1", all[0].Sequence, all[2].Sequence)
	}
	// Each row received an auto-assigned ID.
	for _, tk := range all {
		if tk.ID == "" {
			t.Fatalf("tick seq %d has empty ID", tk.Sequence)
		}
	}
}

func TestMemory_ListEngineTicks_Limit(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	for i := 1; i <= 5; i++ {
		if _, err := m.RecordEngineTick(ctx, engine.ScheduledTick{Sequence: int64(i)}); err != nil {
			t.Fatalf("record: %v", err)
		}
	}
	got, err := m.ListEngineTicks(ctx, 2)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Sequence != 5 || got[1].Sequence != 4 {
		t.Fatalf("limited newest-first = [%d,%d], want [5,4]", got[0].Sequence, got[1].Sequence)
	}
}

func TestMemory_ListEngineTicks_EmptyNotNil(t *testing.T) {
	t.Parallel()
	got, err := NewMemory().ListEngineTicks(context.Background(), 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if got == nil {
		t.Fatal("ListEngineTicks returned nil, want a non-nil empty slice")
	}
	if len(got) != 0 {
		t.Fatalf("len = %d, want 0", len(got))
	}
}
