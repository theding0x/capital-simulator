# Atlas Ephemeral Per-Session In-Memory Runs — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the Atlas Observatory simulation reset on every page reload by running each browser session's economy entirely in server RAM (seeded once from MySQL, advanced on poll, never written back), while UI preferences persist in `localStorage`.

**Architecture:** A new `internal/observatory` package owns an in-memory seed template and a `map[sessionID]*Run`. The browser mints an ephemeral `X-Atlas-Session` id per page load; `GET /v1/observatory/snapshot?advance=N` advances that session's run N periods (using the tested `simulation.AdvanceGeneralLaw` domain) and returns it. The `AccumulationTicker` + `GeneralLawTicker` are removed from the global scheduler, so no per-tick Atlas writes hit MySQL. No new routes, no gateway change, no migrations.

**Tech Stack:** Go 1.25 (stdlib only — `sync`, `time`, `context`, `log/slog`), React 18 + TS + Vite, the existing `simulation`/`circulation` domains.

**Spec:** `docs/superpowers/specs/2026-06-05-atlas-ephemeral-session-runs-design.md`

---

## Conventions for this plan

- Go import groups: stdlib / blank line / third-party / blank line / local. (No third-party here.)
- JSON tags `snake_case`; value-magnitudes are `circulation.Pence` (int64) / `int64`.
- Every test fn calls `t.Parallel()` first.
- Run Go commands from the repo root. Module path: `github.com/theding0x/capital-simulator`.
- Commit after each task (the repo signs commits via `commit.gpgsign=true`).

---

## Task 1: `observatory.advanceField` — pure orrery accumulation

**Files:**
- Create: `services/simulation-engine/internal/observatory/field.go`
- Test: `services/simulation-engine/internal/observatory/field_test.go`

- [ ] **Step 1: Write the failing test**

Create `services/simulation-engine/internal/observatory/field_test.go`:

```go
package observatory

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

func TestAdvanceFieldGrowsTotalByAlphaSurplus(t *testing.T) {
	t.Parallel()
	in := []circulation.FieldCapital{{
		ID:              "5eed0000000000000004ic1",
		TotalPence:      300000,
		MoneyPence:      100000,
		ProductionPence: 100000,
		CommodityPence:  100000,
		CostPricePence:  240000,
		SurplusPence:    60000,
	}}
	out := advanceField(in, 5000) // capitalise 50% of 60000 surplus = 30000

	if out[0].TotalPence != 330000 {
		t.Fatalf("TotalPence = %d, want 330000", out[0].TotalPence)
	}
	if sum := out[0].MoneyPence + out[0].ProductionPence + out[0].CommodityPence; sum != out[0].TotalPence {
		t.Fatalf("M+P+C = %d, want %d (parts must sum to total)", sum, out[0].TotalPence)
	}
	if out[0].SurplusPence != 60000 || out[0].CostPricePence != 240000 {
		t.Fatalf("surplus/cost mutated: %d/%d", out[0].SurplusPence, out[0].CostPricePence)
	}
	if in[0].TotalPence != 300000 {
		t.Fatalf("input mutated (not pure): TotalPence = %d, want 300000", in[0].TotalPence)
	}
}

func TestAdvanceFieldZeroSurplusIsNoOp(t *testing.T) {
	t.Parallel()
	in := []circulation.FieldCapital{{ID: "x", TotalPence: 1000, MoneyPence: 1000, SurplusPence: 0}}
	out := advanceField(in, 5000)
	if out[0].TotalPence != 1000 {
		t.Fatalf("TotalPence = %d, want 1000 (no surplus → no growth)", out[0].TotalPence)
	}
}

func TestAdvanceFieldZeroAlphaIsNoOp(t *testing.T) {
	t.Parallel()
	in := []circulation.FieldCapital{{ID: "x", TotalPence: 1000, MoneyPence: 1000, SurplusPence: 500}}
	out := advanceField(in, 0)
	if out[0].TotalPence != 1000 {
		t.Fatalf("TotalPence = %d, want 1000 (alpha 0 → no growth)", out[0].TotalPence)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./services/simulation-engine/internal/observatory/ -run TestAdvanceField -v`
Expected: FAIL — `undefined: advanceField` (package doesn't compile yet).

- [ ] **Step 3: Write the implementation**

Create `services/simulation-engine/internal/observatory/field.go`:

```go
// Package observatory runs the Atlas Observatory as ephemeral, per-session,
// in-memory simulation runs. A Manager seeds each browser session's Run from an
// immutable template (loaded once from the store at boot) and advances it on
// poll; nothing is written back to the store, so a page reload (new session)
// starts a clean run from seed.
package observatory

import "github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"

// advanceField capitalises alphaBP basis points of each capital's surplus back
// into it — the spiral of accumulation (Vol. I Ch. 24) applied to the orrery
// read-model. TotalPence grows and the M/P/C arcs rescale to the new total
// (CommodityPence absorbs integer rounding so the parts still sum to the total).
// CostPrice, Surplus, Status and TurnoverNumber are unchanged. Pure: returns a
// new slice and never mutates its input. Mirrors store.Memory.AccumulateCapital.
func advanceField(in []circulation.FieldCapital, alphaBP int64) []circulation.FieldCapital {
	out := make([]circulation.FieldCapital, len(in))
	copy(out, in)
	if alphaBP <= 0 {
		return out
	}
	for i := range out {
		fc := out[i]
		if fc.SurplusPence <= 0 {
			continue
		}
		delta := circulation.Pence(int64(fc.SurplusPence) * alphaBP / 10000)
		if delta <= 0 {
			continue
		}
		oldTotal := fc.TotalPence
		newTotal := oldTotal + delta
		if oldTotal > 0 {
			money := fc.MoneyPence * newTotal / oldTotal
			prod := fc.ProductionPence * newTotal / oldTotal
			fc.MoneyPence = money
			fc.ProductionPence = prod
			fc.CommodityPence = newTotal - money - prod
		} else {
			fc.MoneyPence = newTotal
			fc.ProductionPence = 0
			fc.CommodityPence = 0
		}
		fc.TotalPence = newTotal
		out[i] = fc
	}
	return out
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./services/simulation-engine/internal/observatory/ -run TestAdvanceField -v`
Expected: PASS (3 tests).

- [ ] **Step 5: Commit**

```bash
git add services/simulation-engine/internal/observatory/field.go services/simulation-engine/internal/observatory/field_test.go
git commit -m "feat(atlas): pure in-memory field accumulation for observatory runs"
```

---

## Task 2: `observatory.Run` — one session's in-memory run

**Files:**
- Create: `services/simulation-engine/internal/observatory/run.go`
- Test: `services/simulation-engine/internal/observatory/run_test.go`

- [ ] **Step 1: Write the failing test**

Create `services/simulation-engine/internal/observatory/run_test.go`:

```go
package observatory

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

func seedRun() *Run {
	return &Run{
		abode: simulation.NewAbodeState(),
		field: []circulation.FieldCapital{{
			ID: "ic1", TotalPence: 300000, MoneyPence: 300000, SurplusPence: 60000, CostPricePence: 240000,
		}},
	}
}

func TestRunAdvanceIncrementsTickGrowsFieldAndSeries(t *testing.T) {
	t.Parallel()
	r := seedRun()
	r.Advance(3)
	snap := r.Snapshot()
	if snap.Tick != 3 {
		t.Fatalf("Tick = %d, want 3", snap.Tick)
	}
	if len(snap.Periods) != 3 {
		t.Fatalf("len(Periods) = %d, want 3", len(snap.Periods))
	}
	if snap.Field[0].TotalPence <= 300000 {
		t.Fatalf("field did not accumulate: %d", snap.Field[0].TotalPence)
	}
}

func TestRunAdvanceZeroIsNoOp(t *testing.T) {
	t.Parallel()
	r := seedRun()
	r.Advance(0)
	snap := r.Snapshot()
	if snap.Tick != 0 || len(snap.Periods) != 0 {
		t.Fatalf("Advance(0) changed state: tick=%d periods=%d", snap.Tick, len(snap.Periods))
	}
}

func TestRunAdvanceClampsAndCapsSeries(t *testing.T) {
	t.Parallel()
	r := seedRun()
	for i := 0; i < 11; i++ {
		r.Advance(1000) // each call clamps to maxAdvancePerPoll
	}
	snap := r.Snapshot()
	if len(snap.Periods) > maxSeries {
		t.Fatalf("len(Periods) = %d, want <= %d", len(snap.Periods), maxSeries)
	}
	if snap.Tick != int64(11*maxAdvancePerPoll) {
		t.Fatalf("Tick = %d, want %d", snap.Tick, 11*maxAdvancePerPoll)
	}
}

func TestRunApplyLeversClampsAndPersistsInSession(t *testing.T) {
	t.Parallel()
	r := seedRun()
	over := int64(999999)
	abode := r.ApplyLevers(simulation.LeverUpdate{SurplusRateBaseBP: &over})
	if abode.SurplusRateBaseBP != 100000 {
		t.Fatalf("returned SurplusRateBaseBP = %d, want clamped 100000", abode.SurplusRateBaseBP)
	}
	if got := r.Snapshot().Abode.SurplusRateBaseBP; got != 100000 {
		t.Fatalf("snapshot SurplusRateBaseBP = %d, want 100000", got)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./services/simulation-engine/internal/observatory/ -run TestRun -v`
Expected: FAIL — `undefined: Run`, `undefined: maxSeries`, `undefined: maxAdvancePerPoll`.

- [ ] **Step 3: Write the implementation**

Create `services/simulation-engine/internal/observatory/run.go`:

```go
package observatory

import (
	"sync"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

const (
	// maxSeries caps the immiseration time-series retained per run (the abode
	// sparkline shows the most recent 60 periods).
	maxSeries = 60
	// maxAdvancePerPoll bounds how many General-Law periods a single poll may
	// advance, so a crafted advance count cannot pin the CPU.
	maxAdvancePerPoll = 10
)

// Run is one in-memory Atlas simulation run, owned by a single browser session.
// It is seeded from the Manager's template and advanced on poll; it is never
// persisted. All access is guarded by mu.
type Run struct {
	mu      sync.Mutex
	abode   simulation.AbodeState         // carries the live levers too
	field   []circulation.FieldCapital    // the orrery
	periods []simulation.GeneralLawPeriod // capped to the last maxSeries
	tick    int64
}

// RunSnapshot is the read-model projection of a Run at one instant; the
// transport layer maps it to the observatory snapshot response.
type RunSnapshot struct {
	Tick    int64
	Abode   simulation.AbodeState
	Readout simulation.AbodeReadout
	Field   []circulation.FieldCapital
	Periods []simulation.GeneralLawPeriod
}

// Advance runs n periods of the General Law (clamped to [0, maxAdvancePerPoll]),
// growing the field by accumulation and appending each period to the capped
// immiseration series.
func (r *Run) Advance(n int) {
	if n < 0 {
		n = 0
	}
	if n > maxAdvancePerPoll {
		n = maxAdvancePerPoll
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := 0; i < n; i++ {
		next, period := simulation.AdvanceGeneralLaw(r.abode)
		r.abode = next
		r.periods = append(r.periods, period)
		if len(r.periods) > maxSeries {
			r.periods = r.periods[len(r.periods)-maxSeries:]
		}
		r.field = advanceField(r.field, r.abode.AccumulationRateBP)
		r.tick++
	}
}

// ApplyLevers perturbs the run's live abode (clamped) and returns the new abode.
func (r *Run) ApplyLevers(u simulation.LeverUpdate) simulation.AbodeState {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.abode = r.abode.ApplyLevers(u)
	return r.abode
}

// Snapshot returns a race-free copy of the run's current state.
func (r *Run) Snapshot() RunSnapshot {
	r.mu.Lock()
	defer r.mu.Unlock()
	field := make([]circulation.FieldCapital, len(r.field))
	copy(field, r.field)
	periods := make([]simulation.GeneralLawPeriod, len(r.periods))
	copy(periods, r.periods)
	return RunSnapshot{
		Tick:    r.tick,
		Abode:   r.abode,
		Readout: r.abode.Readout(),
		Field:   field,
		Periods: periods,
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./services/simulation-engine/internal/observatory/ -run TestRun -v`
Expected: PASS (4 tests).

- [ ] **Step 5: Commit**

```bash
git add services/simulation-engine/internal/observatory/run.go services/simulation-engine/internal/observatory/run_test.go
git commit -m "feat(atlas): per-session in-memory Run (advance-on-poll, levers, snapshot)"
```

---

## Task 3: `observatory.Manager` — seed template, sessions, eviction

**Files:**
- Create: `services/simulation-engine/internal/observatory/manager.go`
- Test: `services/simulation-engine/internal/observatory/manager_test.go`

- [ ] **Step 1: Write the failing test**

Create `services/simulation-engine/internal/observatory/manager_test.go`:

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./services/simulation-engine/internal/observatory/ -run TestManager -v`
Expected: FAIL — `undefined: Manager`, `undefined: NewManager`.

- [ ] **Step 3: Write the implementation**

Create `services/simulation-engine/internal/observatory/manager.go`:

```go
package observatory

import (
	"context"
	"log/slog"
	"sort"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// session is one tracked run plus the time it was last accessed (for eviction).
type session struct {
	run        *Run
	lastAccess time.Time
}

// Manager owns the immutable seed template and the live, in-memory sessions. It
// is safe for concurrent use. The seed is read once at construction; runs never
// write back to any store.
type Manager struct {
	mu       sync.Mutex
	sessions map[string]*session

	seedAbode simulation.AbodeState
	seedField []circulation.FieldCapital

	ttl           time.Duration
	maxSessions   int
	sweepInterval time.Duration
	now           func() time.Time
	logger        *slog.Logger
}

// NewManager builds a Manager from the seed abode and seed field. The seed field
// is copied and sorted by ID so every run starts from a deterministic order.
func NewManager(seedAbode simulation.AbodeState, seedField []circulation.FieldCapital, logger *slog.Logger) *Manager {
	if logger == nil {
		logger = slog.Default()
	}
	field := make([]circulation.FieldCapital, len(seedField))
	copy(field, seedField)
	sort.Slice(field, func(i, j int) bool { return field[i].ID < field[j].ID })
	return &Manager{
		sessions:      make(map[string]*session),
		seedAbode:     seedAbode,
		seedField:     field,
		ttl:           15 * time.Minute,
		maxSessions:   500,
		sweepInterval: time.Minute,
		now:           time.Now,
		logger:        logger,
	}
}

// GetOrCreate returns the Run for sessionID, creating it from the seed template
// if absent. An empty sessionID returns a fresh, unstored (transient) run, so
// header-less callers get a clean seed snapshot without populating the map.
func (m *Manager) GetOrCreate(sessionID string) *Run {
	if sessionID == "" {
		return m.newRun()
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.sessions[sessionID]; ok {
		s.lastAccess = m.now()
		return s.run
	}
	if len(m.sessions) >= m.maxSessions {
		m.evictOldestLocked()
	}
	run := m.newRun()
	m.sessions[sessionID] = &session{run: run, lastAccess: m.now()}
	return run
}

// Len reports the number of live sessions.
func (m *Manager) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.sessions)
}

// StartSweeper launches a background goroutine that evicts idle sessions every
// sweepInterval until ctx is cancelled.
func (m *Manager) StartSweeper(ctx context.Context) {
	go func() {
		t := time.NewTicker(m.sweepInterval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				m.sweep()
			}
		}
	}()
}

// newRun deep-copies the seed template into a fresh Run (AbodeState is all-scalar
// so it copies by value; the field slice is copied; the series starts empty).
func (m *Manager) newRun() *Run {
	field := make([]circulation.FieldCapital, len(m.seedField))
	copy(field, m.seedField)
	return &Run{abode: m.seedAbode, field: field}
}

// sweep removes every session idle longer than ttl.
func (m *Manager) sweep() {
	m.mu.Lock()
	defer m.mu.Unlock()
	cutoff := m.now().Add(-m.ttl)
	for id, s := range m.sessions {
		if s.lastAccess.Before(cutoff) {
			delete(m.sessions, id)
		}
	}
}

// evictOldestLocked removes the least-recently-accessed session. Caller holds mu.
func (m *Manager) evictOldestLocked() {
	var oldestID string
	var oldest time.Time
	first := true
	for id, s := range m.sessions {
		if first || s.lastAccess.Before(oldest) {
			oldestID, oldest, first = id, s.lastAccess, false
		}
	}
	if !first {
		delete(m.sessions, oldestID)
	}
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./services/simulation-engine/internal/observatory/ -v`
Expected: PASS (all field/run/manager tests).

- [ ] **Step 5: Commit**

```bash
git add services/simulation-engine/internal/observatory/manager.go services/simulation-engine/internal/observatory/manager_test.go
git commit -m "feat(atlas): observatory session Manager with seed template + TTL/cap eviction"
```

---

## Task 4: rewire the transport handlers to the Manager

**Files:**
- Modify: `services/simulation-engine/internal/transport/httpapi/handler.go` (add `Observatory` to `Handler`, `Deps`, `New`)
- Replace: `services/simulation-engine/internal/transport/httpapi/observatory_handler.go`
- Replace: `services/simulation-engine/internal/transport/httpapi/observatory_levers_handler.go`
- Replace: `services/simulation-engine/internal/transport/httpapi/observatory_handler_test.go`

- [ ] **Step 1: Write the failing test (rewrite the handler test)**

Replace the entire contents of `services/simulation-engine/internal/transport/httpapi/observatory_handler_test.go` with:

```go
package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/observatory"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// newTestObservatory seeds a memory store with two Marx-faithful capitals, then
// builds a Manager from the seed snapshot (as main.go does at boot).
func newTestObservatory(t *testing.T) *observatory.Manager {
	t.Helper()
	ctx := context.Background()
	m := store.NewMemory()
	a, _ := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{TotalPence: 500000, EconomyMode: circulation.EconomyMoney})
	b, _ := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{TotalPence: 300000, EconomyMode: circulation.EconomyMoney})
	_, _ = m.Snapshot(ctx, a.ID, circulation.StageDistribution{MoneyPence: 100000, ProductionPence: 300000, CommodityPence: 100000})
	_, _ = m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{IndustrialCapitalID: a.ID, Period: "1871", DemandPence: 400000, SupplyPence: 480000, ExcessPence: 80000})
	_, _ = m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{IndustrialCapitalID: b.ID, Period: "1871", DemandPence: 250000, SupplyPence: 275000, ExcessPence: 25000})
	seedAbode, _ := m.GetAbodeState(ctx)
	seedField, _ := m.FieldSnapshot(ctx)
	return observatory.NewManager(seedAbode, seedField, slog.Default())
}

func snapshot(t *testing.T, h *Handler, sess string, advance int) observatorySnapshotResponse {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/observatory/snapshot?advance=%d", advance), nil)
	if sess != "" {
		req.Header.Set("X-Atlas-Session", sess)
	}
	rec := httptest.NewRecorder()
	h.GetObservatorySnapshot(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("snapshot status = %d; body=%s", rec.Code, rec.Body.String())
	}
	var resp observatorySnapshotResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return resp
}

func TestGetObservatorySnapshotSeed(t *testing.T) {
	t.Parallel()
	h := New(nil, Deps{Observatory: newTestObservatory(t)})
	resp := snapshot(t, h, "seed", 0) // advance=0 → pristine seed

	if resp.Capitals == nil {
		t.Fatal("capitals must never be null")
	}
	if len(resp.Capitals) != 2 {
		t.Fatalf("capitals = %d, want 2", len(resp.Capitals))
	}
	// ΣS = 105000, ΣC = 650000 → p̄′ = round(10000*105000/650000) = 1615 bp.
	if got := resp.Aggregate.AvgRateOfProfitBP; got != 1615 {
		t.Errorf("avg_rate_of_profit_bp = %d, want 1615", got)
	}
	if resp.Aggregate.TotalSocialCapitalPence != 800000 {
		t.Errorf("total social capital = %d, want 800000", resp.Aggregate.TotalSocialCapitalPence)
	}
	if resp.Aggregate.SurplusPence != 105000 || resp.Aggregate.CostPricePence != 650000 {
		t.Errorf("aggregate cost/surplus = %d/%d", resp.Aggregate.CostPricePence, resp.Aggregate.SurplusPence)
	}
	if resp.Abode.TotalVariablePence <= 0 {
		t.Errorf("abode Σv = %d, want > 0", resp.Abode.TotalVariablePence)
	}
	if resp.Abode.LawSeries == nil {
		t.Error("abode.law_series must be non-nil (never null)")
	}
	if resp.Abode.AccumulationRateBP == 0 && resp.Abode.BaseWagePence == 0 {
		t.Error("abode block missing base lever values")
	}
}

func TestGetObservatorySnapshotAdvanceAccumulates(t *testing.T) {
	t.Parallel()
	h := New(nil, Deps{Observatory: newTestObservatory(t)})
	base := snapshot(t, h, "grow", 0)
	after := snapshot(t, h, "grow", 3) // same session advances 3 periods
	if after.Tick != 3 {
		t.Fatalf("tick = %d, want 3", after.Tick)
	}
	if after.Aggregate.TotalSocialCapitalPence <= base.Aggregate.TotalSocialCapitalPence {
		t.Fatalf("capital did not accumulate: %d <= %d", after.Aggregate.TotalSocialCapitalPence, base.Aggregate.TotalSocialCapitalPence)
	}
}

func TestGetObservatorySnapshotSessionsIndependent(t *testing.T) {
	t.Parallel()
	h := New(nil, Deps{Observatory: newTestObservatory(t)})
	_ = snapshot(t, h, "A", 5)
	b := snapshot(t, h, "B", 0)
	if b.Tick != 0 {
		t.Fatalf("session B tick = %d, want 0 (independent of A)", b.Tick)
	}
}

func TestGetObservatorySnapshotEmptyFieldNeverNull(t *testing.T) {
	t.Parallel()
	empty := observatory.NewManager(simulation.NewAbodeState(), nil, slog.Default())
	h := New(nil, Deps{Observatory: empty})
	resp := snapshot(t, h, "e", 0)
	if resp.Capitals == nil {
		t.Fatal("empty field: capitals must be [] not null")
	}
	if resp.Aggregate.AvgRateOfProfitBP != 0 {
		t.Errorf("empty field p̄′ = %d, want 0 (no divide-by-zero)", resp.Aggregate.AvgRateOfProfitBP)
	}
}

func TestSetObservatoryLevers(t *testing.T) {
	t.Parallel()
	h := New(nil, Deps{Observatory: newTestObservatory(t)})

	body := `{"accumulation_rate_bp": 0, "surplus_rate_base_bp": 30000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/observatory/levers", strings.NewReader(body))
	req.Header.Set("X-Atlas-Session", "lev")
	rec := httptest.NewRecorder()
	h.SetObservatoryLevers(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}
	var resp struct {
		SurplusRateBaseBP  int64 `json:"surplus_rate_base_bp"`
		BaseWagePence      int64 `json:"base_wage_pence"`
		AccumulationRateBP int64 `json:"accumulation_rate_bp"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.AccumulationRateBP != 0 || resp.SurplusRateBaseBP != 30000 {
		t.Errorf("levers not applied: %+v", resp)
	}
	if resp.BaseWagePence != 2500 {
		t.Errorf("untouched wage = %d, want 2500 (seed default)", resp.BaseWagePence)
	}

	// reflected in the SAME session's snapshot
	same := snapshot(t, h, "lev", 0)
	if same.Abode.SurplusRateBaseBP != 30000 || same.Abode.AccumulationRateBP != 0 {
		t.Errorf("levers not reflected in session snapshot: %+v", same.Abode)
	}

	// a DIFFERENT session still sees the seed default (10000)
	other := snapshot(t, h, "other", 0)
	if other.Abode.SurplusRateBaseBP == 30000 {
		t.Error("levers leaked across sessions")
	}

	// malformed body → 400
	bad := httptest.NewRequest(http.MethodPost, "/v1/observatory/levers", strings.NewReader("{"))
	bad.Header.Set("X-Atlas-Session", "lev")
	badRec := httptest.NewRecorder()
	h.SetObservatoryLevers(badRec, bad)
	if badRec.Code != http.StatusBadRequest {
		t.Errorf("malformed code = %d, want 400", badRec.Code)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./services/simulation-engine/internal/transport/httpapi/ -run Observatory -v`
Expected: FAIL — `Deps has no field Observatory` / the handlers still read the store.

- [ ] **Step 3a: Add the `Observatory` dependency to `handler.go`**

In `services/simulation-engine/internal/transport/httpapi/handler.go`, add the import (local group) — change:

```go
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/surplus"
```

to:

```go
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/observatory"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/surplus"
```

In the `Handler` struct, replace:

```go
	Scheduler             *engine.Scheduler
	EngineTicks           store.EngineTickStore
}
```

with:

```go
	Scheduler             *engine.Scheduler
	Observatory           *observatory.Manager
	EngineTicks           store.EngineTickStore
}
```

In the `Deps` struct, replace the identical two lines:

```go
	Scheduler             *engine.Scheduler
	EngineTicks           store.EngineTickStore
}
```

with:

```go
	Scheduler             *engine.Scheduler
	Observatory           *observatory.Manager
	EngineTicks           store.EngineTickStore
}
```

In `New`, replace:

```go
		Scheduler:             d.Scheduler,
		EngineTicks:           d.EngineTicks,
	}
}
```

with:

```go
		Scheduler:             d.Scheduler,
		Observatory:           d.Observatory,
		EngineTicks:           d.EngineTicks,
	}
}
```

- [ ] **Step 3b: Replace `observatory_handler.go`**

Replace the entire contents of `services/simulation-engine/internal/transport/httpapi/observatory_handler.go` with:

```go
package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/observatory"
)

// snapshotIntervalMS is the advisory client poll interval echoed in the snapshot.
// Advancement is driven by the client's `advance` query param, not the server.
const snapshotIntervalMS = 2000

// observatorySnapshotResponse is the GET /v1/observatory/snapshot body: the whole
// field of industrial capitals, the aggregate vital-signs, and the hidden abode,
// for one session's in-memory run. Consumed by the Atlas page.
type observatorySnapshotResponse struct {
	Tick       int64              `json:"tick"`
	Running    bool               `json:"running"`
	IntervalMS int64              `json:"interval_ms"`
	Capitals   []fieldCapitalDTO  `json:"capitals"`
	Aggregate  aggregateVitalsDTO `json:"aggregate"`
	Abode      abodeDTO           `json:"abode"`
}

type fieldCapitalDTO struct {
	ID              string `json:"id"`
	TotalPence      int64  `json:"total_pence"`
	MoneyPence      int64  `json:"money_pence"`
	ProductionPence int64  `json:"production_pence"`
	CommodityPence  int64  `json:"commodity_pence"`
	CostPricePence  int64  `json:"cost_price_pence"`
	SurplusPence    int64  `json:"surplus_pence"`
	Status          string `json:"status"`
	TurnoverNumber  int64  `json:"turnover_number"`
}

type aggregateVitalsDTO struct {
	TotalSocialCapitalPence int64 `json:"total_social_capital_pence"`
	CostPricePence          int64 `json:"cost_price_pence"`
	SurplusPence            int64 `json:"surplus_pence"`
	AvgRateOfProfitBP       int64 `json:"avg_rate_of_profit_bp"`
}

type abodeDTO struct {
	TotalVariablePence     int64                 `json:"total_variable_pence"`
	TotalSurplusPence      int64                 `json:"total_surplus_pence"`
	RateOfExploitationBP   int64                 `json:"rate_of_exploitation_bp"`
	NecessaryLabourMinutes int64                 `json:"necessary_labour_minutes"`
	SurplusLabourMinutes   int64                 `json:"surplus_labour_minutes"`
	OrganicCompositionBP   int64                 `json:"organic_composition_bp"`
	ReserveArmyCount       int64                 `json:"reserve_army_count"`
	ReserveArmyPressureBP  int64                 `json:"reserve_army_pressure_bp"`
	EmployedCount          int64                 `json:"employed_count"`
	WagePence              int64                 `json:"wage_pence"`
	SurplusRateBaseBP      int64                 `json:"surplus_rate_base_bp"`
	BaseWagePence          int64                 `json:"base_wage_pence"`
	AccumulationRateBP     int64                 `json:"accumulation_rate_bp"`
	LawSeries              []generalLawPeriodDTO `json:"law_series"`
}

type generalLawPeriodDTO struct {
	Period               int64 `json:"period"`
	WagePence            int64 `json:"wage_pence"`
	RateOfExploitationBP int64 `json:"rate_of_exploitation_bp"`
	ReserveArmyCount     int64 `json:"reserve_army_count"`
	OrganicCompositionBP int64 `json:"organic_composition_bp"`
}

// GetObservatorySnapshot handles GET /v1/observatory/snapshot?advance=N. It reads
// the X-Atlas-Session header, advances that session's in-memory run by N periods
// (default 1), and returns the projection. No store I/O; nothing is persisted.
func (h *Handler) GetObservatorySnapshot(w http.ResponseWriter, r *http.Request) {
	if h.Observatory == nil {
		h.writeServerError(w, errors.New("observatory not configured"))
		return
	}
	advance := parseAdvance(r.URL.Query().Get("advance"))
	run := h.Observatory.GetOrCreate(r.Header.Get("X-Atlas-Session"))
	run.Advance(advance)
	writeJSON(w, http.StatusOK, buildSnapshotResponse(run.Snapshot(), advance))
}

// parseAdvance returns the requested advance count: empty or invalid → 1, "0" →
// 0 (paused). Run.Advance clamps the upper bound.
func parseAdvance(q string) int {
	if q == "" {
		return 1
	}
	n, err := strconv.Atoi(q)
	if err != nil || n < 0 {
		return 1
	}
	return n
}

// buildSnapshotResponse maps a RunSnapshot to the wire DTO. Slices are always
// non-nil so the client never sees `null`.
func buildSnapshotResponse(snap observatory.RunSnapshot, advance int) observatorySnapshotResponse {
	resp := observatorySnapshotResponse{
		Tick:       snap.Tick,
		Running:    advance > 0,
		IntervalMS: snapshotIntervalMS,
		Capitals:   make([]fieldCapitalDTO, len(snap.Field)),
	}
	var sumTotal, sumCost, sumSurplus int64
	for i, fc := range snap.Field {
		resp.Capitals[i] = fieldCapitalDTO{
			ID:              string(fc.ID),
			TotalPence:      int64(fc.TotalPence),
			MoneyPence:      int64(fc.MoneyPence),
			ProductionPence: int64(fc.ProductionPence),
			CommodityPence:  int64(fc.CommodityPence),
			CostPricePence:  int64(fc.CostPricePence),
			SurplusPence:    int64(fc.SurplusPence),
			Status:          string(fc.Status),
			TurnoverNumber:  fc.TurnoverNumber,
		}
		sumTotal += int64(fc.TotalPence)
		sumCost += int64(fc.CostPricePence)
		sumSurplus += int64(fc.SurplusPence)
	}
	resp.Aggregate = aggregateVitalsDTO{
		TotalSocialCapitalPence: sumTotal,
		CostPricePence:          sumCost,
		SurplusPence:            sumSurplus,
		AvgRateOfProfitBP:       rateBP(sumSurplus, sumCost),
	}
	ar := snap.Readout
	law := make([]generalLawPeriodDTO, len(snap.Periods))
	for i, p := range snap.Periods {
		law[i] = generalLawPeriodDTO{
			Period:               p.Period,
			WagePence:            p.WagePence,
			RateOfExploitationBP: p.RateOfExploitationBP,
			ReserveArmyCount:     p.ReserveArmyCount,
			OrganicCompositionBP: p.OrganicCompositionBP,
		}
	}
	resp.Abode = abodeDTO{
		TotalVariablePence:     ar.TotalVariablePence,
		TotalSurplusPence:      ar.TotalSurplusPence,
		RateOfExploitationBP:   ar.RateOfExploitationBP,
		NecessaryLabourMinutes: ar.NecessaryLabourMinutes,
		SurplusLabourMinutes:   ar.SurplusLabourMinutes,
		OrganicCompositionBP:   ar.OrganicCompositionBP,
		ReserveArmyCount:       ar.ReserveArmyCount,
		ReserveArmyPressureBP:  ar.ReserveArmyPressureBP,
		EmployedCount:          ar.EmployedCount,
		WagePence:              ar.WagePence,
		SurplusRateBaseBP:      snap.Abode.SurplusRateBaseBP,
		BaseWagePence:          snap.Abode.BaseWagePence,
		AccumulationRateBP:     snap.Abode.AccumulationRateBP,
		LawSeries:              law,
	}
	return resp
}

// rateBP returns round-half-up(10000 * num / den) basis points; 0 when den <= 0.
func rateBP(num, den int64) int64 {
	if den <= 0 {
		return 0
	}
	return (num*10000 + den/2) / den
}
```

- [ ] **Step 3c: Replace `observatory_levers_handler.go`**

Replace the entire contents of `services/simulation-engine/internal/transport/httpapi/observatory_levers_handler.go` with:

```go
package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// leversRequest is the POST /v1/observatory/levers body. Each field is optional;
// omitted fields leave that lever unchanged.
type leversRequest struct {
	SurplusRateBaseBP  *int64 `json:"surplus_rate_base_bp"`
	BaseWagePence      *int64 `json:"base_wage_pence"`
	AccumulationRateBP *int64 `json:"accumulation_rate_bp"`
}

// leversResponse echoes the applied (clamped) lever values.
type leversResponse struct {
	SurplusRateBaseBP  int64 `json:"surplus_rate_base_bp"`
	BaseWagePence      int64 `json:"base_wage_pence"`
	AccumulationRateBP int64 `json:"accumulation_rate_bp"`
}

// SetObservatoryLevers handles POST /v1/observatory/levers. It applies a partial
// perturbation of the abode's law parameters — the working day (s′), the wage
// (the value of labour-power), and the accumulation rate α — to the caller's
// in-memory session run (keyed by X-Atlas-Session). Returns the applied
// (clamped) values. No store I/O.
func (h *Handler) SetObservatoryLevers(w http.ResponseWriter, r *http.Request) {
	if h.Observatory == nil {
		h.writeServerError(w, errors.New("observatory not configured"))
		return
	}
	var body leversRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "malformed JSON body")
		return
	}
	run := h.Observatory.GetOrCreate(r.Header.Get("X-Atlas-Session"))
	abode := run.ApplyLevers(simulation.LeverUpdate{
		SurplusRateBaseBP:  body.SurplusRateBaseBP,
		BaseWagePence:      body.BaseWagePence,
		AccumulationRateBP: body.AccumulationRateBP,
	})
	writeJSON(w, http.StatusOK, leversResponse{
		SurplusRateBaseBP:  abode.SurplusRateBaseBP,
		BaseWagePence:      abode.BaseWagePence,
		AccumulationRateBP: abode.AccumulationRateBP,
	})
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./services/simulation-engine/internal/transport/httpapi/ -run Observatory -v`
Expected: PASS (TestGetObservatorySnapshotSeed, …AdvanceAccumulates, …SessionsIndependent, …EmptyFieldNeverNull, TestSetObservatoryLevers).

- [ ] **Step 5: Commit**

```bash
git add services/simulation-engine/internal/transport/httpapi/handler.go \
  services/simulation-engine/internal/transport/httpapi/observatory_handler.go \
  services/simulation-engine/internal/transport/httpapi/observatory_levers_handler.go \
  services/simulation-engine/internal/transport/httpapi/observatory_handler_test.go
git commit -m "feat(atlas): observatory handlers serve per-session in-memory runs"
```

---

## Task 5: wire the Manager into `main.go` and drop the Atlas tickers

**Files:**
- Modify: `services/simulation-engine/cmd/simulation-engine/main.go`

- [ ] **Step 1: Add the observatory import**

Change the local import group — replace:

```go
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/piecewage"
```

with:

```go
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/observatory"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/piecewage"
```

- [ ] **Step 2: Remove the Atlas tickers from the scheduler**

Replace:

```go
	scheduler := engine.NewScheduler(tickInterval(), []engine.Ticker{
		engine.NewFactoryTicker(st),
		engine.NewReproductionTicker(st),
		engine.NewPiecePriceTicker(engine.NewFactoryProductivitySource(st), repricer),
		engine.NewAccumulationTicker(st, accumulationRateBP()),
		engine.NewGeneralLawTicker(st),
	}, st, logger)
```

with:

```go
	// The Atlas run (field accumulation + General Law) is no longer scheduler-driven:
	// it runs per-session in memory (internal/observatory), advanced on poll, with no
	// MySQL writes. The scheduler keeps only the MySQL-backed chapter tickers.
	scheduler := engine.NewScheduler(tickInterval(), []engine.Ticker{
		engine.NewFactoryTicker(st),
		engine.NewReproductionTicker(st),
		engine.NewPiecePriceTicker(engine.NewFactoryProductivitySource(st), repricer),
	}, st, logger)
```

- [ ] **Step 3: Build the seed template + Manager before constructing the handler**

Immediately before the `h := httpapi.New(logger, httpapi.Deps{` line, insert:

```go
	// Atlas Observatory — load the seed once and hand it to the session Manager.
	// Each browser session gets its own in-memory run; nothing is written back.
	seedAbode, err := st.GetAbodeState(ctx)
	if err != nil {
		logger.Error("could not load seed abode", "err", err)
		os.Exit(1)
	}
	seedField, err := st.FieldSnapshot(ctx)
	if err != nil {
		logger.Error("could not load seed field", "err", err)
		os.Exit(1)
	}
	obsMgr := observatory.NewManager(seedAbode, seedField, logger)
	obsMgr.StartSweeper(ctx)

```

- [ ] **Step 4: Pass the Manager in `Deps`**

In the `httpapi.Deps{ … }` literal, replace:

```go
		Scheduler:             scheduler,
		EngineTicks:           st,
	})
```

with:

```go
		Scheduler:             scheduler,
		Observatory:           obsMgr,
		EngineTicks:           st,
	})
```

- [ ] **Step 5: Remove the now-unused `accumulationRateBP` helper and `strconv` import**

Delete this entire function (and its doc comment):

```go
// accumulationRateBP returns the share of surplus reinvested per scheduler pass,
// in basis points, from SIM_ACCUMULATION_RATE_BP (default 5000 = 50%). Clamped
// to [0, 10000].
func accumulationRateBP() int64 {
	v := os.Getenv("SIM_ACCUMULATION_RATE_BP")
	if v == "" {
		return 5000
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil || n < 0 {
		return 5000
	}
	if n > 10000 {
		return 10000
	}
	return n
}

```

Then remove the now-unused `"strconv"` line from the stdlib import group.

- [ ] **Step 6: Build and vet**

Run: `make vet build`
Expected: builds clean; no "declared and not used", no unused-import errors.

- [ ] **Step 7: Smoke-test in-memory (no MySQL needed)**

Run:
```bash
MYSQL_DISABLED=true SERVICE_ADDR=:8084 go run ./services/simulation-engine/cmd/simulation-engine &
sleep 2
curl -s -H "X-Atlas-Session: smoke" "http://localhost:8084/v1/observatory/snapshot?advance=0" | head -c 400
echo
curl -s -H "X-Atlas-Session: smoke" "http://localhost:8084/v1/observatory/snapshot?advance=3" | head -c 200
echo
kill %1
```
Expected: first call returns a JSON snapshot with `"tick":0`; second returns `"tick":3` (the same session advanced). With `MYSQL_DISABLED=true` the field is empty, so `capitals` is `[]` and the abode block is seeded from `NewAbodeState()`.

- [ ] **Step 8: Commit**

```bash
git add services/simulation-engine/cmd/simulation-engine/main.go
git commit -m "feat(atlas): seed observatory Manager at boot; drop Atlas tickers from scheduler"
```

---

## Task 6: frontend — ephemeral session id + API wiring

**Files:**
- Create: `web/src/atlas/session.ts`
- Modify: `web/src/api.ts` (thread headers through `http`; send session header + `advance`)
- Modify: `web/src/atlas/useSnapshot.ts` (accept `advance`)

> Note: the frontend does not typecheck until Task 7 also lands (Atlas.tsx is updated there). Commit Tasks 6 + 7 together (Task 7 Step 8).

- [ ] **Step 1: Create the session module**

Create `web/src/atlas/session.ts`:

```ts
// A per-page-load Atlas session id. Held in memory only — never persisted — so a
// reload mints a new id and the server starts that session's run fresh from seed.
function newSessionId(): string {
  const c = globalThis.crypto;
  if (c && typeof c.randomUUID === "function") return c.randomUUID();
  if (c && typeof c.getRandomValues === "function") {
    const b = new Uint8Array(16);
    c.getRandomValues(b);
    return Array.from(b, (x) => x.toString(16).padStart(2, "0")).join("");
  }
  return `s-${Date.now()}-${Math.random().toString(16).slice(2)}`;
}

export const atlasSessionId: string = newSessionId();

/** Header to attach to every Atlas observatory request. */
export const atlasSessionHeader: Record<string, string> = {
  "X-Atlas-Session": atlasSessionId,
};
```

- [ ] **Step 2: Let `http()` accept custom headers**

In `web/src/api.ts`, replace the `http` helper's fetch call — change:

```ts
  const res = await fetch(`${BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...init,
  });
```

to (headers last so callers can add `X-Atlas-Session` without losing Content-Type):

```ts
  const res = await fetch(`${BASE}${path}`, {
    ...init,
    headers: { "Content-Type": "application/json", ...(init?.headers ?? {}) },
  });
```

- [ ] **Step 3: Send the session header + advance on the observatory calls**

At the top of `web/src/api.ts`, after the existing `} from "./types";` import block, add:

```ts
import { atlasSessionHeader } from "./atlas/session";
```

Replace the two observatory functions — change:

```ts
  getObservatorySnapshot: () =>
    http<ObservatorySnapshot>("/v1/observatory/snapshot"),

  setObservatoryLevers: (u: import("./types").LeverUpdate) =>
    http<import("./types").LeverState>("/v1/observatory/levers", {
      method: "POST",
      body: JSON.stringify(u),
    }),
```

to:

```ts
  getObservatorySnapshot: (advance = 1) =>
    http<ObservatorySnapshot>(`/v1/observatory/snapshot?advance=${advance}`, {
      headers: atlasSessionHeader,
    }),

  setObservatoryLevers: (u: import("./types").LeverUpdate) =>
    http<import("./types").LeverState>("/v1/observatory/levers", {
      method: "POST",
      headers: atlasSessionHeader,
      body: JSON.stringify(u),
    }),
```

- [ ] **Step 4: Pass `advance` through `useSnapshot`**

Replace the entire contents of `web/src/atlas/useSnapshot.ts` with:

```ts
import { useEffect, useRef, useState } from "react";
import { api } from "../api";
import type { ObservatorySnapshot } from "../types";

const POLL_MS = 2000;

export interface SnapshotState {
  snapshot: ObservatorySnapshot | null;
  error: string | null;
  stale: boolean;
}

/** Polls the observatory snapshot every 2s, advancing the session's in-memory run
 *  by `advance` periods per poll (0 = paused). Holds last-good on error. */
export function useSnapshot(advance: number): SnapshotState {
  const [snapshot, setSnapshot] = useState<ObservatorySnapshot | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [stale, setStale] = useState(false);
  const advanceRef = useRef(advance);
  advanceRef.current = advance;
  const timer = useRef<number | null>(null);

  useEffect(() => {
    let active = true;
    async function poll() {
      try {
        const snap = await api.getObservatorySnapshot(advanceRef.current);
        if (!active) return;
        setSnapshot(snap);
        setError(null);
        setStale(false);
      } catch (e) {
        if (!active) return;
        setError(e instanceof Error ? e.message : "snapshot failed");
        setStale(true);
      }
    }
    void poll();
    timer.current = window.setInterval(() => void poll(), POLL_MS);
    return () => {
      active = false;
      if (timer.current !== null) window.clearInterval(timer.current);
    };
  }, []);

  return { snapshot, error, stale };
}
```

---

## Task 7: frontend — run lifecycle + persisted UI prefs

**Files:**
- Create: `web/src/atlas/prefs.ts`
- Modify: `web/src/atlas/Atlas.tsx`
- Modify: `web/src/CurrencyContext.tsx`

- [ ] **Step 1: Create the prefs module**

Create `web/src/atlas/prefs.ts`:

```ts
// Durable UI preferences for the Atlas page, persisted in localStorage so they
// survive a reload (unlike the simulation run, which resets). Defensive: bad or
// absent values fall back to defaults, and storage failures are non-fatal.
const KEY = "atlas.prefs.v1";

export interface AtlasPrefs {
  speed: number;
  reduced: boolean;
}

const DEFAULTS: AtlasPrefs = { speed: 1, reduced: false };

export function loadPrefs(): AtlasPrefs {
  try {
    const raw = localStorage.getItem(KEY);
    if (!raw) return { ...DEFAULTS };
    const p = JSON.parse(raw) as Partial<AtlasPrefs>;
    return {
      speed: typeof p.speed === "number" ? p.speed : DEFAULTS.speed,
      reduced: typeof p.reduced === "boolean" ? p.reduced : DEFAULTS.reduced,
    };
  } catch {
    return { ...DEFAULTS };
  }
}

export function savePrefs(p: AtlasPrefs): void {
  try {
    localStorage.setItem(KEY, JSON.stringify(p));
  } catch {
    /* storage unavailable (e.g. private mode) — non-fatal */
  }
}
```

- [ ] **Step 2: Update `Atlas.tsx` — swap imports**

In `web/src/atlas/Atlas.tsx`, delete this line:

```ts
import { api } from "../api";
```

and add (next to the other `./` imports, e.g. after the `import { clamp, formatBP } from "./animation";` line):

```ts
import { loadPrefs, savePrefs } from "./prefs";
```

- [ ] **Step 3: Update `Atlas.tsx` — run/speed/reduced state + advance-on-poll**

Replace:

```ts
export default function Atlas() {
  const { snapshot } = useSnapshot();
  const prefersReduced =
    typeof window !== "undefined" &&
    window.matchMedia &&
    window.matchMedia("(prefers-reduced-motion: reduce)").matches;

  const [running, setRunning] = useState(false);
  const [speed, setSpeed] = useState(1);
  const [reduced, setReduced] = useState(!!prefersReduced);
  const [depth, setDepth] = useState(0);
```

with:

```ts
export default function Atlas() {
  const prefersReduced =
    typeof window !== "undefined" &&
    window.matchMedia &&
    window.matchMedia("(prefers-reduced-motion: reduce)").matches;

  const [running, setRunning] = useState(true);
  const [speed, setSpeed] = useState(() => loadPrefs().speed);
  const [reduced, setReduced] = useState(() => loadPrefs().reduced || !!prefersReduced);
  const [depth, setDepth] = useState(0);

  // Advance the session's in-memory run by `speed` periods per poll while running.
  const { snapshot } = useSnapshot(running ? speed : 0);

  // Persist UI preferences across reloads (the run itself does not persist).
  useEffect(() => {
    savePrefs({ speed, reduced });
  }, [speed, reduced]);
```

- [ ] **Step 4: Update `Atlas.tsx` — drop the server run-state sync**

Delete this effect entirely:

```ts
  // reflect the server's run state on first load
  useEffect(() => {
    if (snapshot) setRunning(snapshot.running);
  }, [snapshot?.running]); // eslint-disable-line react-hooks/exhaustive-deps
```

- [ ] **Step 5: Update `Atlas.tsx` — make `toggleRun` pure client state**

Replace:

```ts
  const toggleRun = () => {
    const next = !running;
    setRunning(next);
    surfaceRef.current?.setRunning(next);
    void (next ? api.startEngine() : api.stopEngine()).catch(() => {
      /* poll reflects truth */
    });
  };
```

with:

```ts
  const toggleRun = () => {
    const next = !running;
    setRunning(next);
    surfaceRef.current?.setRunning(next);
  };
```

- [ ] **Step 6: Persist the currency toggle**

Replace the entire contents of `web/src/CurrencyContext.tsx` with:

```tsx
import { createContext, useContext, useState } from "react";
import type { ReactNode } from "react";

import { fmtPounds, fmtPoundsModern } from "./format";

interface CurrencyContextShape {
  modern: boolean;
  toggle: () => void;
}

const STORE_KEY = "atlas.currency.modern";

function loadModern(): boolean {
  try {
    return localStorage.getItem(STORE_KEY) === "true";
  } catch {
    return false;
  }
}

const CurrencyContext = createContext<CurrencyContextShape>({
  modern: false,
  toggle: () => {},
});

export function CurrencyProvider({ children }: { children: ReactNode }) {
  const [modern, setModern] = useState(loadModern);
  const toggle = () =>
    setModern((m) => {
      const next = !m;
      try {
        localStorage.setItem(STORE_KEY, String(next));
      } catch {
        /* storage unavailable — non-fatal */
      }
      return next;
    });
  return (
    <CurrencyContext.Provider value={{ modern, toggle }}>
      {children}
    </CurrencyContext.Provider>
  );
}

export function useCurrency(): CurrencyContextShape {
  return useContext(CurrencyContext);
}

export function usePounds(): (pence: number) => string {
  const { modern } = useCurrency();
  return modern ? fmtPoundsModern : fmtPounds;
}

export function CurrencyToggle() {
  const { modern, toggle } = useCurrency();
  return (
    <button
      className={`currency-toggle${modern ? " currency-toggle--active" : ""}`}
      onClick={toggle}
      title={
        modern
          ? "Showing 2025 prices — click for 1860s historical"
          : "Showing 1860s prices — click for 2025 equivalent"
      }
    >
      {modern ? "2025 £" : "1860s £"}
    </button>
  );
}
```

- [ ] **Step 7: Typecheck and build**

Run: `cd web && npm run lint && npm run build`
Expected: PASS (tsc clean; Vite build succeeds).

- [ ] **Step 8: Commit Tasks 6 + 7 together**

```bash
git add web/src/atlas/session.ts web/src/atlas/prefs.ts web/src/api.ts \
  web/src/atlas/useSnapshot.ts web/src/atlas/Atlas.tsx web/src/CurrencyContext.tsx
git commit -m "feat(atlas): ephemeral session id, advance-on-poll, persisted UI prefs"
```

---

## Task 8: document the new runtime model

**Files:**
- Modify: `docs/architecture.md`

- [ ] **Step 1: Add a note in the Atlas/simulation-engine section**

Add this block to `docs/architecture.md` near the Atlas Observatory / simulation-engine description (search for "Atlas" or "observatory" and place it after the existing Atlas paragraph):

```markdown
### Atlas Observatory — ephemeral per-session runs

The Atlas Observatory run (the orrery field, the aggregate vitals, the hidden
abode, and the immiseration series) is **not persisted**. At boot the
simulation-engine reads the seed once (`abode_state` + the seeded
`industrial_capitals`) into an immutable template (`internal/observatory`). Each
browser session sends an ephemeral `X-Atlas-Session` header; the server keeps a
per-session in-memory `Run`, deep-copied from the template, advanced on poll via
`GET /v1/observatory/snapshot?advance=N` (`N = running ? speed : 0`). Levers
mutate the session's run only. A page reload mints a new session id ⇒ a clean run
from seed. UI preferences (currency, speed, reduced-motion) persist in the
browser's `localStorage`. The `general_law` and `accumulation` tickers were
removed from the scheduler, so the Atlas run performs **zero MySQL writes** at
runtime; idle sessions are evicted by TTL.
```

- [ ] **Step 2: Commit**

```bash
git add docs/architecture.md
git commit -m "docs(atlas): describe ephemeral per-session in-memory observatory runs"
```

---

## Task 9: full verification + end-to-end check

**Files:** none (verification only).

- [ ] **Step 1: Full Go check**

Run: `go mod tidy && make vet test build`
Expected: all packages build; all tests pass (including `internal/observatory` and the rewritten `httpapi` observatory tests).

- [ ] **Step 2: Full web check**

Run: `cd web && npm run lint && npm run build`
Expected: clean.

- [ ] **Step 3: Boot the full stack & confirm the gateway forwards header + query**

Run: `docker compose up --build -d`, wait for health, then:

```bash
curl -s -H "X-Atlas-Session: e2e" "http://localhost:8080/v1/observatory/snapshot?advance=0" | head -c 200
```
Expected: 200 with a JSON snapshot, `"tick":0`, populated `capitals` (the MySQL seed). (No gateway change was needed — `/v1/observatory/snapshot` is already proxied and Go's `ReverseProxy` forwards `X-Atlas-Session`.)

- [ ] **Step 4: Confirm NO MySQL writes during a run**

```bash
# baseline
docker compose exec -T mysql mysql -uroot -proot -N -e \
  "SELECT period FROM simulation.abode_state; SELECT COUNT(*) FROM simulation.general_law_periods;"
# drive the run
for i in $(seq 1 10); do curl -s -H "X-Atlas-Session: e2e" \
  "http://localhost:8080/v1/observatory/snapshot?advance=5" >/dev/null; done
# after
docker compose exec -T mysql mysql -uroot -proot -N -e \
  "SELECT period FROM simulation.abode_state; SELECT COUNT(*) FROM simulation.general_law_periods;"
```
Expected: `abode_state.period` and the `general_law_periods` row count are **identical** before and after — the run did not touch MySQL. (Adjust mysql credentials to match `deploy/mysql`/compose if different.)

- [ ] **Step 5: Playwright — reload resets, prefs persist**

With the web app served (`http://localhost:5173` dev, or the nginx prod image), drive it with Playwright:
1. Open the Atlas page; let it run; read the console "turn" counter — it climbs above 0.
2. Toggle currency to "2025 £"; set speed ×5.
3. Reload the page.
4. Assert: the "turn" counter restarts near 0 (fresh session ⇒ reset run); currency still shows "2025 £" and speed is still ×5 (localStorage prefs persisted).
5. Pause; confirm the "turn" counter stops climbing across several polls; play; confirm it resumes.
6. Move a lever (e.g. working day) and confirm the abode readout reacts; reload; confirm levers are back at seed defaults.

- [ ] **Step 6: Tear down**

Run: `docker compose down`

- [ ] **Step 7: Final commit (if any verification tweaks were needed)**

```bash
git add -A
git commit -m "test(atlas): verify ephemeral per-session runs end-to-end"
```

---

## Self-review notes (author)

- **Spec coverage:** seed-once template (Task 5) · per-session in-memory run (Tasks 2–3) · advance-on-poll, no new routes (Tasks 3b/4) · levers per session, no DB write (Task 4) · scheduler tickers dropped ⇒ no Atlas MySQL writes (Task 5) · ephemeral session id, reset on reload (Task 6) · UI prefs in localStorage, levers reset (Task 7) · TTL/cap eviction (Task 3) · docs (Task 8) · verification incl. no-MySQL-writes + Playwright reload (Task 9). All spec sections map to a task.
- **No new migration / no gateway change** — confirmed: `/v1/observatory/snapshot` and `/levers` are already proxied; advancement is a query param on the existing path; the custom header is forwarded by Go's `ReverseProxy`.
- **Type consistency:** `Run.Advance(int)`, `Run.ApplyLevers(simulation.LeverUpdate) simulation.AbodeState`, `Run.Snapshot() RunSnapshot`, `Manager.GetOrCreate(string) *Run`, `Manager.StartSweeper(context.Context)`, `Manager.Len() int` — referenced identically across Tasks 2–5. `advanceField([]circulation.FieldCapital, int64)` matches its caller in `run.go`. Wire DTO field names unchanged, so `web/src/types.ts` needs no edit.
- **Behaviour note:** the orrery's accumulation alpha now comes from the abode's `AccumulationRateBP` (the lever), not the old `SIM_ACCUMULATION_RATE_BP` env var. Both default to 5000, so default behaviour is unchanged; the accumulation lever now also drives orrery growth (an intended improvement).
```