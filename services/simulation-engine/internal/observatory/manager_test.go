package observatory

import (
	"log/slog"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

func testManager() *Manager {
	return NewManager(
		simulation.NewAbodeState(),
		[]circulation.FieldCapital{{ID: "ic1", TotalPence: 300000, MoneyPence: 300000, SurplusPence: 60000}},
		slog.Default(),
	)
}

func TestManagerGetOrCreateIsolatesSessions(t *testing.T) {
	t.Parallel()
	m := testManager()
	a := m.GetOrCreate("A")
	b := m.GetOrCreate("B")
	if a == b {
		t.Fatal("distinct sessions returned the same Run")
	}
	a.Advance(5)
	if m.GetOrCreate("B").Snapshot().Tick != 0 {
		t.Fatal("session B advanced when only A was advanced")
	}
}

func TestManagerGetOrCreateSameIDReturnsSameRun(t *testing.T) {
	t.Parallel()
	m := testManager()
	m.GetOrCreate("s").Advance(2)
	if m.GetOrCreate("s").Snapshot().Tick != 2 {
		t.Fatal("same id returned a different run")
	}
}

func TestManagerEmptyIDIsTransient(t *testing.T) {
	t.Parallel()
	m := testManager()
	m.GetOrCreate("").Advance(1)
	if m.Len() != 0 {
		t.Fatalf("empty id populated the session map: len=%d", m.Len())
	}
}

func TestManagerSeedTemplateNotMutatedByRuns(t *testing.T) {
	t.Parallel()
	m := testManager()
	m.GetOrCreate("s").Advance(5)
	fresh := m.GetOrCreate("other")
	if fresh.Snapshot().Field[0].TotalPence != 300000 {
		t.Fatalf("seed template mutated: %d", fresh.Snapshot().Field[0].TotalPence)
	}
}

func TestManagerSweepEvictsIdle(t *testing.T) {
	t.Parallel()
	m := testManager()
	clock := time.Now()
	m.now = func() time.Time { return clock }
	m.ttl = 10 * time.Minute
	m.GetOrCreate("s")
	clock = clock.Add(11 * time.Minute)
	m.sweep()
	if m.Len() != 0 {
		t.Fatalf("idle session not evicted: len=%d", m.Len())
	}
}

func TestManagerCapEvictsOldest(t *testing.T) {
	t.Parallel()
	m := testManager()
	clock := time.Now()
	m.now = func() time.Time { return clock }
	m.maxSessions = 2
	m.GetOrCreate("a")
	clock = clock.Add(time.Minute)
	m.GetOrCreate("b")
	clock = clock.Add(time.Minute)
	m.GetOrCreate("c") // exceeds cap → evicts oldest ("a")
	if m.Len() != 2 {
		t.Fatalf("cap not enforced: len=%d", m.Len())
	}
}
