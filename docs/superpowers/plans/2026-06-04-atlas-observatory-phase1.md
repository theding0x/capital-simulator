# Atlas Observatory — Phase 1 (Living Field) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Ship a watch-only "Atlas" landing page that renders the field of industrial capitals as faithful orbits in continuous motion, driven by a new `GET /v1/observatory/snapshot` read-model and the real engine scheduler.

**Architecture:** A new projection endpoint in `simulation-engine` aggregates each `IndustrialCapital`'s latest `StageDistribution` (the orbit arcs) plus its latest `SupplyDemandImbalance` (cost-price c+v and surplus s), computing the aggregate average rate of profit p̄′ = ΣS/ΣC server-side. The React app gains an `atlas/` module that polls the snapshot every 2s, animates orbits at 60fps via CSS, and controls the engine (start/stop/status) with a tick-heartbeat ECG. A minimal hash router makes Atlas `/` and demotes the existing dashboard to `#/chapters`.

**Tech Stack:** Go 1.25 (`net/http` 1.22 mux, `database/sql`), goose migrations, React 18 + Vite + TS, SVG/CSS animation. No new dependencies.

**Design spec:** `docs/superpowers/specs/2026-06-04-atlas-observatory-design.md`. Branch: `feature/atlas-observatory` (already created).

**Contract note (resolved from spec's open questions):** per-capital `c/v/s` is *not* on `IndustrialCapital`; it lives in `SupplyDemandImbalance` (`demand_pence = c+v`, `excess_pence = s`). So the Phase-1 snapshot exposes `cost_price_pence` (c+v) and `surplus_pence` (s) per capital, and p̄′ = Σexcessᵢ/Σdemandᵢ. The richer vitals (s′, Dept I/II, accumulation) and `turnover_number` are **out of Phase 1** — orbit spin tempo is derived client-side from the capital id.

---

## File Structure

**Backend (`services/simulation-engine/`):**
- Create `internal/circulation/field.go` — `FieldCapital` read-model struct.
- Modify `internal/store/store.go` — add `FieldSnapshot` to `IndustrialCapitalStore`.
- Modify `internal/store/memory.go` — implement `FieldSnapshot`.
- Modify `internal/store/mysql.go` — implement `FieldSnapshot`.
- Create `internal/transport/httpapi/observatory_handler.go` — handler + DTOs + p̄′.
- Create `internal/transport/httpapi/observatory_handler_test.go` — handler test.
- Modify `internal/transport/httpapi/routes.go` — register the route.
- Create `internal/store/migrations/00064_atlas_field_seed.sql` — extra branches.

**Gateway:** Modify `services/api-gateway/cmd/api-gateway/main.go` — one proxy line.

**Frontend (`web/src/`):**
- Modify `types.ts` — snapshot/engine mirror types.
- Modify `api.ts` — snapshot + engine calls.
- Create `atlas/animation.ts` — pure render math.
- Create `atlas/useSnapshot.ts` — polling hook.
- Create `atlas/Orbit.tsx` — one faithful orbit.
- Create `atlas/CircuitField.tsx` — the field.
- Create `atlas/VitalSigns.tsx` — aggregate rail.
- Create `atlas/TickHeartbeat.tsx` — transport + ECG.
- Create `atlas/Atlas.tsx` — Observatory shell.
- Create `atlas/atlas.css` — styles (brand tokens).
- Create `Root.tsx` — minimal hash router.
- Modify `main.tsx` — render `<Root/>`.
- Modify `App.tsx` — add "← Atlas" link in topbar.

---

# GROUP A — Backend snapshot (Go, test-first)

## Task A1: `FieldCapital` read-model + store method (interface + memory + mysql)

**Files:**
- Create: `services/simulation-engine/internal/circulation/field.go`
- Modify: `services/simulation-engine/internal/store/store.go:226` (inside `IndustrialCapitalStore`, before its closing `}`)
- Modify: `services/simulation-engine/internal/store/memory.go` (append a method)
- Modify: `services/simulation-engine/internal/store/mysql.go` (append a method)
- Test: `services/simulation-engine/internal/store/field_snapshot_test.go`

- [ ] **Step 1: Write the failing memory test**

Create `services/simulation-engine/internal/store/field_snapshot_test.go`:

```go
package store_test

import (
	"context"
	"sort"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func TestMemoryFieldSnapshot(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	// Two capitals of differing organic composition (Vol. III Ch. 9 spheres).
	a, err := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 500000, EconomyMode: circulation.EconomyMoney,
	})
	if err != nil {
		t.Fatalf("create a: %v", err)
	}
	b, err := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 300000, EconomyMode: circulation.EconomyMoney,
	})
	if err != nil {
		t.Fatalf("create b: %v", err)
	}

	if _, err := m.Snapshot(ctx, a.ID, circulation.StageDistribution{
		MoneyPence: 100000, ProductionPence: 300000, CommodityPence: 100000,
	}); err != nil {
		t.Fatalf("snapshot a: %v", err)
	}
	if _, err := m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{
		IndustrialCapitalID: a.ID, Period: "1871",
		DemandPence: 400000, SupplyPence: 480000, ExcessPence: 80000, // p' = 20%
	}); err != nil {
		t.Fatalf("supply-demand a: %v", err)
	}
	if _, err := m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{
		IndustrialCapitalID: b.ID, Period: "1871",
		DemandPence: 250000, SupplyPence: 275000, ExcessPence: 25000, // p' = 10%
	}); err != nil {
		t.Fatalf("supply-demand b: %v", err)
	}

	field, err := m.FieldSnapshot(ctx)
	if err != nil {
		t.Fatalf("FieldSnapshot: %v", err)
	}
	if len(field) != 2 {
		t.Fatalf("want 2 capitals, got %d", len(field))
	}
	sort.Slice(field, func(i, j int) bool { return field[i].ID < field[j].ID })

	byID := map[circulation.IndustrialCapitalID]circulation.FieldCapital{}
	for _, fc := range field {
		byID[fc.ID] = fc
	}
	fa := byID[a.ID]
	if fa.MoneyPence != 100000 || fa.ProductionPence != 300000 || fa.CommodityPence != 100000 {
		t.Errorf("a distribution: %+v", fa)
	}
	if fa.CostPricePence != 400000 || fa.SurplusPence != 80000 {
		t.Errorf("a cost/surplus: cost=%d surplus=%d", fa.CostPricePence, fa.SurplusPence)
	}

	// Capital with no distribution recorded defaults to all-money.
	c, _ := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 120000, EconomyMode: circulation.EconomyMoney,
	})
	field, _ = m.FieldSnapshot(ctx)
	for _, fc := range field {
		if fc.ID == c.ID && fc.MoneyPence != 120000 {
			t.Errorf("default-money capital: money=%d want 120000", fc.MoneyPence)
		}
	}
}
```

- [ ] **Step 2: Run the test to confirm it fails to compile**

Run: `cd services/simulation-engine && go test ./internal/store/ -run TestMemoryFieldSnapshot`
Expected: FAIL — `m.FieldSnapshot undefined` and `circulation.FieldCapital` undefined.

- [ ] **Step 3: Add the `FieldCapital` struct**

Create `services/simulation-engine/internal/circulation/field.go`:

```go
package circulation

// FieldCapital is the per-capital projection consumed by the Atlas Observatory
// snapshot. It carries the latest StageDistribution (the orbit's M/P/C arcs),
// the capital's cost-price (c+v) and surplus (s) from its latest
// SupplyDemandImbalance, and its run status. Read-model only: no persistence,
// no ID constructor.
type FieldCapital struct {
	ID              IndustrialCapitalID
	TotalPence      Pence
	MoneyPence      Pence
	ProductionPence Pence
	CommodityPence  Pence
	CostPricePence  Pence // c + v (SupplyDemandImbalance.DemandPence)
	SurplusPence    Pence // s     (SupplyDemandImbalance.ExcessPence)
	Status          IndustrialCapitalStatus
}
```

- [ ] **Step 4: Add the interface method**

In `services/simulation-engine/internal/store/store.go`, inside `type IndustrialCapitalStore interface { ... }` (after the `TickSinkingFund` line, before the closing `}`), add:

```go
	// FieldSnapshot returns every industrial capital projected to a FieldCapital
	// for the Atlas Observatory: latest stage distribution + latest cost-price
	// and surplus. Capitals with no recorded distribution default to all-money.
	FieldSnapshot(ctx context.Context) ([]circulation.FieldCapital, error)
```

- [ ] **Step 5: Implement in `memory.go`**

Append to `services/simulation-engine/internal/store/memory.go` (after `TickSinkingFund`, in the IndustrialCapital section). The mutex field is `m.mu` (an `sync.RWMutex`):

```go
// FieldSnapshot implements IndustrialCapitalStore for the Atlas Observatory.
func (m *Memory) FieldSnapshot(_ context.Context) ([]circulation.FieldCapital, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []circulation.FieldCapital{}
	for id, ic := range m.industrialCapitals {
		fc := circulation.FieldCapital{
			ID:         id,
			TotalPence: ic.TotalPence,
			MoneyPence: ic.TotalPence, // default when no distribution recorded
			Status:     ic.Status,
		}
		if sds := m.stageDistributions[id]; len(sds) > 0 {
			latest := sds[len(sds)-1]
			fc.MoneyPence = latest.MoneyPence
			fc.ProductionPence = latest.ProductionPence
			fc.CommodityPence = latest.CommodityPence
		}
		if sdis := m.supplyDemand[id]; len(sdis) > 0 {
			latest := sdis[len(sdis)-1]
			fc.CostPricePence = latest.DemandPence
			fc.SurplusPence = latest.ExcessPence
		}
		out = append(out, fc)
	}
	return out, nil
}
```

- [ ] **Step 6: Implement in `mysql.go`**

Append to `services/simulation-engine/internal/store/mysql.go`. The receiver is `m *MySQL`, the handle is `m.db`, and `database/sql` is already imported:

```go
// FieldSnapshot implements IndustrialCapitalStore for the Atlas Observatory.
func (m *MySQL) FieldSnapshot(ctx context.Context) ([]circulation.FieldCapital, error) {
	const q = `
SELECT ic.id, ic.total_pence, ic.status,
       sd.money_pence, sd.production_pence, sd.commodity_pence,
       sdi.demand_pence, sdi.excess_pence
FROM industrial_capitals ic
LEFT JOIN (
    SELECT s.industrial_capital_id, s.money_pence, s.production_pence, s.commodity_pence
    FROM stage_distributions s
    JOIN (
        SELECT industrial_capital_id, MAX(at_time) AS mt
        FROM stage_distributions GROUP BY industrial_capital_id
    ) lt ON lt.industrial_capital_id = s.industrial_capital_id AND lt.mt = s.at_time
) sd ON sd.industrial_capital_id = ic.id
LEFT JOIN (
    SELECT d.industrial_capital_id, d.demand_pence, d.excess_pence
    FROM supply_demand_imbalances d
    JOIN (
        SELECT industrial_capital_id, MAX(period) AS mp
        FROM supply_demand_imbalances GROUP BY industrial_capital_id
    ) lp ON lp.industrial_capital_id = d.industrial_capital_id AND lp.mp = d.period
) sdi ON sdi.industrial_capital_id = ic.id
ORDER BY ic.id`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []circulation.FieldCapital{}
	for rows.Next() {
		var (
			id, status        string
			total             int64
			money, prod, comm sql.NullInt64
			demand, excess    sql.NullInt64
		)
		if err := rows.Scan(&id, &total, &status, &money, &prod, &comm, &demand, &excess); err != nil {
			return nil, err
		}
		fc := circulation.FieldCapital{
			ID:         circulation.IndustrialCapitalID(id),
			TotalPence: circulation.Pence(total),
			Status:     circulation.IndustrialCapitalStatus(status),
			MoneyPence: circulation.Pence(total), // default when no distribution
		}
		if money.Valid {
			fc.MoneyPence = circulation.Pence(money.Int64)
			fc.ProductionPence = circulation.Pence(prod.Int64)
			fc.CommodityPence = circulation.Pence(comm.Int64)
		}
		if demand.Valid {
			fc.CostPricePence = circulation.Pence(demand.Int64)
			fc.SurplusPence = circulation.Pence(excess.Int64)
		}
		out = append(out, fc)
	}
	return out, rows.Err()
}
```

- [ ] **Step 7: Run the test to verify it passes**

Run: `cd services/simulation-engine && go test ./internal/store/ -run TestMemoryFieldSnapshot -v`
Expected: PASS.

- [ ] **Step 8: Commit**

```bash
git add services/simulation-engine/internal/circulation/field.go \
        services/simulation-engine/internal/store/store.go \
        services/simulation-engine/internal/store/memory.go \
        services/simulation-engine/internal/store/mysql.go \
        services/simulation-engine/internal/store/field_snapshot_test.go
git commit --no-gpg-sign -m "feat(atlas): add FieldSnapshot store read-model"
```

---

## Task A2: Observatory handler + p̄′ aggregation (test-first)

**Files:**
- Create: `services/simulation-engine/internal/transport/httpapi/observatory_handler.go`
- Test: `services/simulation-engine/internal/transport/httpapi/observatory_handler_test.go`

- [ ] **Step 1: Write the failing handler test**

Create `services/simulation-engine/internal/transport/httpapi/observatory_handler_test.go`:

```go
package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func TestGetObservatorySnapshot(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	a, _ := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 500000, EconomyMode: circulation.EconomyMoney,
	})
	b, _ := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 300000, EconomyMode: circulation.EconomyMoney,
	})
	_, _ = m.Snapshot(ctx, a.ID, circulation.StageDistribution{
		MoneyPence: 100000, ProductionPence: 300000, CommodityPence: 100000,
	})
	_, _ = m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{
		IndustrialCapitalID: a.ID, Period: "1871",
		DemandPence: 400000, SupplyPence: 480000, ExcessPence: 80000,
	})
	_, _ = m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{
		IndustrialCapitalID: b.ID, Period: "1871",
		DemandPence: 250000, SupplyPence: 275000, ExcessPence: 25000,
	})

	h := New(nil, Deps{IndustrialCapitals: m}) // Scheduler nil → tick 0, not running

	req := httptest.NewRequest(http.MethodGet, "/v1/observatory/snapshot", nil)
	rec := httptest.NewRecorder()
	h.GetObservatorySnapshot(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	var resp observatorySnapshotResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
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
}

func TestGetObservatorySnapshotEmptyFieldNeverNull(t *testing.T) {
	t.Parallel()
	h := New(nil, Deps{IndustrialCapitals: store.NewMemory()})
	rec := httptest.NewRecorder()
	h.GetObservatorySnapshot(rec, httptest.NewRequest(http.MethodGet, "/v1/observatory/snapshot", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	var resp observatorySnapshotResponse
	_ = json.Unmarshal(rec.Body.Bytes(), &resp)
	if resp.Capitals == nil {
		t.Fatal("empty field: capitals must be [] not null")
	}
	if resp.Aggregate.AvgRateOfProfitBP != 0 {
		t.Errorf("empty field p̄′ = %d, want 0 (no divide-by-zero)", resp.Aggregate.AvgRateOfProfitBP)
	}
}
```

- [ ] **Step 2: Run to confirm it fails**

Run: `cd services/simulation-engine && go test ./internal/transport/httpapi/ -run TestGetObservatorySnapshot`
Expected: FAIL — `h.GetObservatorySnapshot` and `observatorySnapshotResponse` undefined.

- [ ] **Step 3: Write the handler**

Create `services/simulation-engine/internal/transport/httpapi/observatory_handler.go`:

```go
package httpapi

import (
	"errors"
	"net/http"
	"sort"
)

// observatorySnapshotResponse is the GET /v1/observatory/snapshot body: the
// whole field of industrial capitals, the aggregate vital-signs, and the
// engine's current tick/run state. Consumed by the Atlas page.
type observatorySnapshotResponse struct {
	Tick       int64              `json:"tick"`
	Running    bool               `json:"running"`
	IntervalMS int64              `json:"interval_ms"`
	Capitals   []fieldCapitalDTO  `json:"capitals"`
	Aggregate  aggregateVitalsDTO `json:"aggregate"`
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
}

type aggregateVitalsDTO struct {
	TotalSocialCapitalPence int64 `json:"total_social_capital_pence"`
	CostPricePence          int64 `json:"cost_price_pence"`
	SurplusPence            int64 `json:"surplus_pence"`
	AvgRateOfProfitBP       int64 `json:"avg_rate_of_profit_bp"`
}

// GetObservatorySnapshot handles GET /v1/observatory/snapshot. The capitals
// array is always non-null; p̄′ = ΣS/ΣC (round-half-up basis points), 0 when
// ΣC == 0. Engine tick/running come from the scheduler when configured.
func (h *Handler) GetObservatorySnapshot(w http.ResponseWriter, r *http.Request) {
	if h.IndustrialCapitals == nil {
		h.writeServerError(w, errors.New("industrial capital store not configured"))
		return
	}
	field, err := h.IndustrialCapitals.FieldSnapshot(r.Context())
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	sort.Slice(field, func(i, j int) bool { return field[i].ID < field[j].ID })

	resp := observatorySnapshotResponse{Capitals: make([]fieldCapitalDTO, len(field))}
	var sumTotal, sumCost, sumSurplus int64
	for i, fc := range field {
		resp.Capitals[i] = fieldCapitalDTO{
			ID:              string(fc.ID),
			TotalPence:      int64(fc.TotalPence),
			MoneyPence:      int64(fc.MoneyPence),
			ProductionPence: int64(fc.ProductionPence),
			CommodityPence:  int64(fc.CommodityPence),
			CostPricePence:  int64(fc.CostPricePence),
			SurplusPence:    int64(fc.SurplusPence),
			Status:          string(fc.Status),
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
	if h.Scheduler != nil {
		st := h.Scheduler.Status()
		resp.Tick = st.Tick
		resp.Running = st.Running
		resp.IntervalMS = st.IntervalMS
	}
	writeJSON(w, http.StatusOK, resp)
}

// rateBP returns round-half-up(10000 * num / den) basis points; 0 when den <= 0.
func rateBP(num, den int64) int64 {
	if den <= 0 {
		return 0
	}
	return (num*10000 + den/2) / den
}
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `cd services/simulation-engine && go test ./internal/transport/httpapi/ -run TestGetObservatorySnapshot -v`
Expected: PASS (both `TestGetObservatorySnapshot` and `TestGetObservatorySnapshotEmptyFieldNeverNull`).

- [ ] **Step 5: Commit**

```bash
git add services/simulation-engine/internal/transport/httpapi/observatory_handler.go \
        services/simulation-engine/internal/transport/httpapi/observatory_handler_test.go
git commit --no-gpg-sign -m "feat(atlas): observatory snapshot handler with p̄′ aggregation"
```

---

## Task A3: Register the route

**Files:**
- Modify: `services/simulation-engine/internal/transport/httpapi/routes.go:10` (after the engine block)

- [ ] **Step 1: Add the route**

In `routes.go`, immediately after the four `/v1/engine/...` lines (line 10), add:

```go
	// Atlas — the whole-circuit Observatory (field snapshot read-model)
	s.HandleFunc("GET /v1/observatory/snapshot", h.GetObservatorySnapshot)
```

- [ ] **Step 2: Build the service**

Run: `cd services/simulation-engine && go build ./...`
Expected: no output (success).

- [ ] **Step 3: Commit**

```bash
git add services/simulation-engine/internal/transport/httpapi/routes.go
git commit --no-gpg-sign -m "feat(atlas): route GET /v1/observatory/snapshot"
```

---

## Task A4: Gateway proxy line

**Files:**
- Modify: `services/api-gateway/cmd/api-gateway/main.go:174` (after `srv.Handle("/v1/engine/ticks", simProxy)`)

- [ ] **Step 1: Add the proxy line**

In `main.go`, right after the engine block (the line `srv.Handle("/v1/engine/ticks", simProxy)`), add:

```go
	// Atlas Observatory snapshot → simulation-engine
	srv.Handle("/v1/observatory/snapshot", simProxy)
```

- [ ] **Step 2: Build the gateway**

Run: `cd services/api-gateway && go build ./...`
Expected: success.

- [ ] **Step 3: Commit**

```bash
git add services/api-gateway/cmd/api-gateway/main.go
git commit --no-gpg-sign -m "feat(atlas): proxy /v1/observatory/snapshot at gateway"
```

---

## Task A5: Seed a richer field (extra branches)

**Files:**
- Create: `services/simulation-engine/internal/store/migrations/00064_atlas_field_seed.sql`

Three branches of differing organic composition join the existing Spinning Mill so the field is populated and p̄′ forms from real spread. Token `90` marks Atlas seeds (no chapter 90 exists, so IDs never collide).

- [ ] **Step 1: Write the seed migration**

Create `services/simulation-engine/internal/store/migrations/00064_atlas_field_seed.sql`:

```sql
-- +goose Up
-- Atlas Observatory field seed — extra industrial capitals of differing
-- organic composition (Vol. III Ch. 9 "spheres") so the average rate of
-- profit visibly forms from spread. Token 90 = Atlas field. IDs: 5eed...90XX.

-- Tannery (low composition, high surplus): total £2,000.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, created_at, updated_at)
VALUES
    ('5eed000000000000009001', 200000, 'money', 3, 'active',
     '1871-01-07 08:00:00.000000', '1871-01-07 08:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009002', '5eed000000000000009001',
     '1871-01-21 08:00:00.000000', 40000, 120000, 40000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009003', '5eed000000000000009001', '1871', 150000, 195000, 45000);

-- Steelworks (high composition, low surplus): total £6,000.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, created_at, updated_at)
VALUES
    ('5eed000000000000009011', 600000, 'money', 3, 'active',
     '1871-01-07 08:00:00.000000', '1871-01-07 08:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009012', '5eed000000000000009011',
     '1871-01-21 08:00:00.000000', 120000, 360000, 120000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009013', '5eed000000000000009011', '1871', 480000, 504000, 24000);

-- Textile (mid composition): total £3,500.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, created_at, updated_at)
VALUES
    ('5eed000000000000009021', 350000, 'money', 3, 'active',
     '1871-01-07 08:00:00.000000', '1871-01-07 08:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009022', '5eed000000000000009021',
     '1871-01-21 08:00:00.000000', 70000, 210000, 70000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009023', '5eed000000000000009021', '1871', 280000, 322000, 42000);

-- +goose Down
DELETE FROM supply_demand_imbalances WHERE id IN ('5eed000000000000009003','5eed000000000000009013','5eed000000000000009023');
DELETE FROM stage_distributions      WHERE id IN ('5eed000000000000009002','5eed000000000000009012','5eed000000000000009022');
DELETE FROM industrial_capitals      WHERE id IN ('5eed000000000000009001','5eed000000000000009011','5eed000000000000009021');
```

- [ ] **Step 2: Boot the stack and verify the snapshot is populated**

Run:
```bash
docker compose up --build -d mysql simulation-engine api-gateway
sleep 8
curl -s http://localhost:8080/v1/observatory/snapshot | head -c 1200; echo
```
Expected: JSON with a `capitals` array containing the Spinning Mill (`5eed...0401`) plus the three new branches, and an `aggregate.avg_rate_of_profit_bp` > 0.

> If the existing `mysql_data` volume is stale (goose "out-of-order" or missing tables), reset just this schema: `make mysql-bootstrap`, or drop/recreate the `simulation` schema, then restart `simulation-engine`.

- [ ] **Step 3: Commit**

```bash
git add services/simulation-engine/internal/store/migrations/00064_atlas_field_seed.sql
git commit --no-gpg-sign -m "feat(atlas): seed a field of differing-composition branches"
```

- [ ] **Step 4: Full backend check**

Run: `make vet test build`
Expected: PASS.

---

# GROUP B — Frontend Atlas (tsc + build + Playwright)

> The web app has no JS unit-test runner (verification = `npm run lint` (tsc) + `npm run build` + Playwright on the booted stack). Pure functions in `animation.ts` are written test-seam-ready for a future vitest setup, but are verified here via tsc and the Playwright acceptance task (B11). Do **not** add a test-runner dependency.

## Task B1: Wire types

**Files:**
- Modify: `web/src/types.ts` (append at end of file)

- [ ] **Step 1: Append the mirror types**

Add to the end of `web/src/types.ts`:

```ts
// --- Atlas Observatory (simulation-engine: GET /v1/observatory/snapshot) ---

export interface FieldCapital {
  id: string;
  total_pence: number;
  money_pence: number;
  production_pence: number;
  commodity_pence: number;
  cost_price_pence: number;
  surplus_pence: number;
  status: string;
}

export interface AggregateVitals {
  total_social_capital_pence: number;
  cost_price_pence: number;
  surplus_pence: number;
  avg_rate_of_profit_bp: number;
}

export interface ObservatorySnapshot {
  tick: number;
  running: boolean;
  interval_ms: number;
  capitals: FieldCapital[];
  aggregate: AggregateVitals;
}

export interface EngineStatus {
  running: boolean;
  tick: number;
  last_tick_at?: string;
  interval_ms: number;
  tickers: string[];
}

export interface EngineTick {
  id: string;
  sequence: number;
  occurred_at: string;
  duration_ms: number;
  tickers_run: number;
  entities_advanced: number;
  error_count: number;
}
```

- [ ] **Step 2: Typecheck**

Run: `cd web && npm run lint`
Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add web/src/types.ts
git commit --no-gpg-sign -m "feat(atlas): web wire types for snapshot and engine"
```

---

## Task B2: API client methods

**Files:**
- Modify: `web/src/api.ts` (import block near line 1–200, and the `api` object near line 226)

- [ ] **Step 1: Add type imports**

In the `import type { ... } from "./types";` block at the top of `web/src/api.ts`, add these names:

```ts
  ObservatorySnapshot,
  EngineStatus,
  EngineTick,
```

- [ ] **Step 2: Add the API methods**

Inside `export const api = {` (right after the opening line), add:

```ts
  // --- Atlas Observatory (simulation-engine) ---
  getObservatorySnapshot: () =>
    http<ObservatorySnapshot>("/v1/observatory/snapshot"),

  getEngineStatus: () => http<EngineStatus>("/v1/engine/status"),

  startEngine: () =>
    http<{ status: EngineStatus }>("/v1/engine/start", { method: "POST" }),

  stopEngine: () =>
    http<{ status: EngineStatus }>("/v1/engine/stop", { method: "POST" }),

  listEngineTicks: (limit = 60) =>
    http<EngineTick[]>(`/v1/engine/ticks?limit=${limit}`),
```

- [ ] **Step 3: Typecheck**

Run: `cd web && npm run lint`
Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add web/src/api.ts
git commit --no-gpg-sign -m "feat(atlas): api client for snapshot and engine control"
```

---

## Task B3: Pure render math

**Files:**
- Create: `web/src/atlas/animation.ts`

- [ ] **Step 1: Write the module**

Create `web/src/atlas/animation.ts`:

```ts
import type { FieldCapital } from "../types";

/** Linear interpolation of `current` toward `target` by fraction `alpha` (0..1). */
export function ease(current: number, target: number, alpha: number): number {
  return current + (target - current) * alpha;
}

/** Arc fractions [money, production, commodity] summing to 1; all-money fallback. */
export function arcFractions(
  c: Pick<FieldCapital, "money_pence" | "production_pence" | "commodity_pence">
): [number, number, number] {
  const sum = c.money_pence + c.production_pence + c.commodity_pence;
  if (sum <= 0) return [1, 0, 0];
  return [c.money_pence / sum, c.production_pence / sum, c.commodity_pence / sum];
}

/** Deterministic per-capital spin tempo (seconds) from its id — the visual turnover. */
export function spinSeconds(id: string): number {
  let h = 0;
  for (let i = 0; i < id.length; i++) h = (h * 31 + id.charCodeAt(i)) >>> 0;
  return 3 + (h % 60) / 10; // 3.0–9.0s
}

/** Orbit pixel radius from total relative to the field max (area-proportional). */
export function orbitRadius(totalPence: number, maxTotal: number, min = 26, max = 74): number {
  if (maxTotal <= 0) return min;
  const t = Math.sqrt(Math.max(0, totalPence) / maxTotal);
  return Math.round(min + (max - min) * t);
}

/** £ with thousands separators (pence → pounds). */
export function formatPence(pence: number): string {
  return "£" + Math.round(pence / 100).toLocaleString("en-GB");
}

/** Basis points → percent string, e.g. 1615 → "16.2%". */
export function formatBP(bp: number): string {
  return (bp / 100).toFixed(1) + "%";
}
```

- [ ] **Step 2: Typecheck**

Run: `cd web && npm run lint`
Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/animation.ts
git commit --no-gpg-sign -m "feat(atlas): pure render math (arcs, radius, formatting)"
```

---

## Task B4: Snapshot polling hook

**Files:**
- Create: `web/src/atlas/useSnapshot.ts`

- [ ] **Step 1: Write the hook**

Create `web/src/atlas/useSnapshot.ts`:

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

/** Polls the observatory snapshot every 2s; holds last-good on error (stale=true). */
export function useSnapshot(): SnapshotState {
  const [snapshot, setSnapshot] = useState<ObservatorySnapshot | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [stale, setStale] = useState(false);
  const timer = useRef<number | null>(null);

  useEffect(() => {
    let active = true;
    async function poll() {
      try {
        const snap = await api.getObservatorySnapshot();
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

- [ ] **Step 2: Typecheck**

Run: `cd web && npm run lint`
Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/useSnapshot.ts
git commit --no-gpg-sign -m "feat(atlas): useSnapshot polling hook"
```

---

## Task B5: Orbit component + styles

**Files:**
- Create: `web/src/atlas/Orbit.tsx`
- Create: `web/src/atlas/atlas.css`

- [ ] **Step 1: Write `atlas.css`**

Create `web/src/atlas/atlas.css`:

```css
/* Atlas Observatory — uses the brand tokens from index.css (:root). */
.atlas-shell { position: fixed; inset: 0; display: flex; flex-direction: column;
  background: radial-gradient(circle at 50% 42%, #12141b, var(--bg) 78%); color: var(--ink); }
.atlas-top { display: flex; align-items: center; justify-content: space-between;
  padding: 0 16px; height: var(--topbar-h); border-bottom: 1px solid var(--border); flex: 0 0 auto; }
.atlas-top .brand { color: var(--red); font-weight: 600; letter-spacing: .5px; }
.atlas-top .nav a { color: var(--ink-muted); text-decoration: none; margin-left: 14px; font-size: 13px; }
.atlas-top .nav a:hover { color: var(--ink); }
.atlas-body { flex: 1 1 auto; display: flex; min-height: 0; }
.atlas-rail { width: 220px; flex: 0 0 auto; border-right: 1px solid var(--border);
  padding: 16px; overflow-y: auto; background: rgba(13,15,20,.6); }
.atlas-field-wrap { flex: 1 1 auto; position: relative; overflow: hidden; }
.atlas-field { position: absolute; inset: 0; display: flex; flex-wrap: wrap;
  align-content: center; justify-content: center; gap: 22px; padding: 28px; }
.atlas-field-centre { position: absolute; left: 50%; top: 50%; transform: translate(-50%,-50%);
  text-align: center; pointer-events: none; opacity: .5; z-index: 0; }
.atlas-field-centre .pbar { font: 700 22px 'IBM Plex Mono', monospace; color: var(--gold-bright); }
.atlas-field-centre .lbl { font: 600 9px 'IBM Plex Mono', monospace; letter-spacing: 1px; color: var(--ink-muted); }
.atlas-orbit { position: relative; z-index: 1; }
.atlas-orbit svg { display: block; }
.atlas-orbit.halted { opacity: .4; filter: grayscale(.6); }
.atlas-flow { transform-box: view-box; transform-origin: center; animation: atlas-spin linear infinite; }
@keyframes atlas-spin { to { transform: rotate(360deg); } }
.atlas-bottom { flex: 0 0 auto; display: flex; align-items: center; gap: 14px;
  height: 56px; padding: 0 16px; border-top: 1px solid var(--border); background: rgba(16,18,24,.7); }
.atlas-btn { background: var(--surface-raised); color: var(--ink); border: 1px solid var(--border);
  border-radius: var(--radius-sm); padding: 6px 12px; cursor: pointer; font-size: 14px; }
.atlas-btn:hover { background: var(--surface-hover); }
.atlas-speed { color: var(--ink-muted); font: 600 11px 'IBM Plex Mono', monospace; cursor: pointer;
  padding: 3px 6px; border-radius: 3px; }
.atlas-speed.active { color: var(--gold-bright); background: var(--gold-bg); }
.atlas-ecg { flex: 1 1 auto; height: 28px; }
.atlas-turn { font: 600 12px 'IBM Plex Mono', monospace; color: var(--ink-muted); }
.atlas-stale { color: var(--red); font-size: 11px; }
.atlas-vital { margin-bottom: 14px; }
.atlas-vital .k { font: 600 8.5px 'IBM Plex Mono', monospace; letter-spacing: 1px; color: var(--ink-muted); text-transform: uppercase; }
.atlas-vital .v { font: 600 18px 'IBM Plex Mono', monospace; color: var(--ink); margin-top: 3px; }
.atlas-vital .v.gold { color: var(--gold-bright); }
```

- [ ] **Step 2: Write `Orbit.tsx`**

Create `web/src/atlas/Orbit.tsx`:

```tsx
import type { FieldCapital } from "../types";
import { arcFractions, orbitRadius, spinSeconds } from "./animation";

interface OrbitProps {
  capital: FieldCapital;
  maxTotal: number;
  /** Animation tempo multiplier (1 = base, 5 = 5× faster). */
  speed: number;
}

const GOLD = "#c8a240";
const RED = "#c0392b";
const LEAD = "#4a5a8a";

/** One faithful orbit: three coexisting arcs (M/P/C) with value-dots circulating. */
export function Orbit({ capital, maxTotal, speed }: OrbitProps) {
  const r = orbitRadius(capital.total_pence, maxTotal);
  const pad = 10;
  const size = (r + pad) * 2;
  const c = size / 2;
  const sw = Math.max(6, Math.round(r * 0.16));
  const [fm, fp] = arcFractions(capital);
  const fc = 1 - fm - fp;
  const gap = 1.5; // pathLength units between arcs

  // pathLength=100; lay arcs M, P, C clockwise from top with small gaps.
  const lm = Math.max(0, fm * 100 - gap);
  const lp = Math.max(0, fp * 100 - gap);
  const lc = Math.max(0, fc * 100 - gap);
  const offM = 0;
  const offP = -(fm * 100);
  const offC = -((fm + fp) * 100);

  const dur = Number((spinSeconds(capital.id) / Math.max(1, speed)).toFixed(2));
  const halted = capital.status === "halted";

  return (
    <div
      className={"atlas-orbit" + (halted ? " halted" : "")}
      title={`${capital.id.slice(0, 8)} · total ${capital.total_pence} · ${capital.status}`}
    >
      <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
        <circle cx={c} cy={c} r={r} fill="none" stroke="#161922" strokeWidth={sw} />
        <g transform={`rotate(-90 ${c} ${c})`}>
          <circle cx={c} cy={c} r={r} fill="none" stroke={GOLD} strokeWidth={sw}
            pathLength={100} strokeDasharray={`${lm} ${100 - lm}`} strokeDashoffset={offM} />
          <circle cx={c} cy={c} r={r} fill="none" stroke={RED} strokeWidth={sw}
            pathLength={100} strokeDasharray={`${lp} ${100 - lp}`} strokeDashoffset={offP} />
          <circle cx={c} cy={c} r={r} fill="none" stroke={LEAD} strokeWidth={sw}
            pathLength={100} strokeDasharray={`${lc} ${100 - lc}`} strokeDashoffset={offC} />
        </g>
        {[0, 1, 2, 3].map((i) => (
          <g key={i} className="atlas-flow"
            style={{ animationDuration: `${dur}s`, animationDelay: `-${(dur / 4) * i}s` }}>
            <circle cx={c} cy={c - r} r={Math.max(2, sw * 0.28)} fill="#f4ecd8" />
          </g>
        ))}
      </svg>
    </div>
  );
}
```

- [ ] **Step 3: Typecheck + build**

Run: `cd web && npm run lint && npm run build`
Expected: success.

- [ ] **Step 4: Commit**

```bash
git add web/src/atlas/Orbit.tsx web/src/atlas/atlas.css
git commit --no-gpg-sign -m "feat(atlas): Orbit component and Atlas styles"
```

---

## Task B6: The field

**Files:**
- Create: `web/src/atlas/CircuitField.tsx`

- [ ] **Step 1: Write `CircuitField.tsx`**

Create `web/src/atlas/CircuitField.tsx`:

```tsx
import type { ObservatorySnapshot } from "../types";
import { Orbit } from "./Orbit";
import { formatBP } from "./animation";

interface CircuitFieldProps {
  snapshot: ObservatorySnapshot;
  speed: number;
}

/** The field of orbits, with the average rate of profit as centre of gravity. */
export function CircuitField({ snapshot, speed }: CircuitFieldProps) {
  const maxTotal = snapshot.capitals.reduce((m, c) => Math.max(m, c.total_pence), 0);
  return (
    <div className="atlas-field-wrap">
      <div className="atlas-field-centre">
        <div className="pbar">{formatBP(snapshot.aggregate.avg_rate_of_profit_bp)}</div>
        <div className="lbl">p̄′ · centre of gravity</div>
      </div>
      <div className="atlas-field" data-testid="atlas-field">
        {snapshot.capitals.map((cap) => (
          <Orbit key={cap.id} capital={cap} maxTotal={maxTotal} speed={speed} />
        ))}
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Typecheck**

Run: `cd web && npm run lint`
Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/CircuitField.tsx
git commit --no-gpg-sign -m "feat(atlas): CircuitField with p̄′ centre of gravity"
```

---

## Task B7: Vital signs rail

**Files:**
- Create: `web/src/atlas/VitalSigns.tsx`

- [ ] **Step 1: Write `VitalSigns.tsx`**

Create `web/src/atlas/VitalSigns.tsx`:

```tsx
import type { AggregateVitals } from "../types";
import { formatBP, formatPence } from "./animation";

/** Aggregate readouts for the rail. */
export function VitalSigns({ vitals, capitalCount }: { vitals: AggregateVitals; capitalCount: number }) {
  return (
    <div>
      <div className="atlas-vital">
        <div className="k">Total social capital</div>
        <div className="v">{formatPence(vitals.total_social_capital_pence)}</div>
      </div>
      <div className="atlas-vital">
        <div className="k">Average rate of profit · p̄′</div>
        <div className="v gold">{formatBP(vitals.avg_rate_of_profit_bp)}</div>
      </div>
      <div className="atlas-vital">
        <div className="k">Surplus-value · ΣS</div>
        <div className="v">{formatPence(vitals.surplus_pence)}</div>
      </div>
      <div className="atlas-vital">
        <div className="k">Cost-price · ΣC (c+v)</div>
        <div className="v">{formatPence(vitals.cost_price_pence)}</div>
      </div>
      <div className="atlas-vital">
        <div className="k">Capitals in motion</div>
        <div className="v">{capitalCount}</div>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Typecheck**

Run: `cd web && npm run lint`
Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/VitalSigns.tsx
git commit --no-gpg-sign -m "feat(atlas): VitalSigns rail"
```

---

## Task B8: Tick heartbeat + transport

**Files:**
- Create: `web/src/atlas/TickHeartbeat.tsx`

- [ ] **Step 1: Write `TickHeartbeat.tsx`**

Create `web/src/atlas/TickHeartbeat.tsx`:

```tsx
import { useEffect, useRef, useState } from "react";
import { api } from "../api";
import type { EngineTick } from "../types";

interface TickHeartbeatProps {
  tick: number;
  running: boolean;
  stale: boolean;
  speed: number;
  onSpeed: (s: number) => void;
}

const SPEEDS = [1, 2, 5, 10];

/** Transport (play/pause/speed), tick counter, and an ECG of recent engine ticks. */
export function TickHeartbeat({ tick, running, stale, speed, onSpeed }: TickHeartbeatProps) {
  const [ticks, setTicks] = useState<EngineTick[]>([]);
  const [localRunning, setLocalRunning] = useState(running);
  const timer = useRef<number | null>(null);

  useEffect(() => setLocalRunning(running), [running]);

  useEffect(() => {
    let active = true;
    async function poll() {
      try {
        const t = await api.listEngineTicks(60);
        if (active) setTicks(t);
      } catch {
        /* heartbeat keeps last-good */
      }
    }
    void poll();
    timer.current = window.setInterval(() => void poll(), 2000);
    return () => {
      active = false;
      if (timer.current !== null) window.clearInterval(timer.current);
    };
  }, []);

  async function toggle() {
    const next = !localRunning;
    setLocalRunning(next); // optimistic
    try {
      if (next) await api.startEngine();
      else await api.stopEngine();
    } catch {
      setLocalRunning(!next); // revert on failure
    }
  }

  const points = ecgPoints(ticks);

  return (
    <>
      <button className="atlas-btn" onClick={() => void toggle()} aria-label={localRunning ? "Pause" : "Play"}>
        {localRunning ? "⏸" : "▶"}
      </button>
      <span role="group" aria-label="Animation speed">
        {SPEEDS.map((s) => (
          <span key={s} className={"atlas-speed" + (speed === s ? " active" : "")}
            onClick={() => onSpeed(s)}>×{s}</span>
        ))}
      </span>
      <svg className="atlas-ecg" viewBox="0 0 240 28" preserveAspectRatio="none" data-testid="atlas-ecg">
        <polyline points={points} fill="none" stroke="#c8a240" strokeWidth="1.4" opacity="0.85" />
      </svg>
      <span className="atlas-turn">turn {tick}</span>
      {stale && <span className="atlas-stale">⚠ reconnecting</span>}
    </>
  );
}

/** Map recent ticks' entities_advanced to an ECG polyline over a 240×28 box. */
function ecgPoints(ticks: EngineTick[]): string {
  if (ticks.length === 0) return "0,14 240,14";
  const ordered = [...ticks].reverse(); // oldest → newest
  const max = Math.max(1, ...ordered.map((t) => t.entities_advanced));
  const n = ordered.length;
  return ordered
    .map((t, i) => {
      const x = (i / Math.max(1, n - 1)) * 240;
      const y = 26 - (t.entities_advanced / max) * 24;
      return `${x.toFixed(1)},${y.toFixed(1)}`;
    })
    .join(" ");
}
```

- [ ] **Step 2: Typecheck**

Run: `cd web && npm run lint`
Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/TickHeartbeat.tsx
git commit --no-gpg-sign -m "feat(atlas): TickHeartbeat transport and ECG"
```

---

## Task B9: Atlas shell

**Files:**
- Create: `web/src/atlas/Atlas.tsx`

- [ ] **Step 1: Write `Atlas.tsx`**

Create `web/src/atlas/Atlas.tsx`:

```tsx
import { useState } from "react";
import "./atlas.css";
import { useSnapshot } from "./useSnapshot";
import { CircuitField } from "./CircuitField";
import { VitalSigns } from "./VitalSigns";
import { TickHeartbeat } from "./TickHeartbeat";

/** The Observatory: the whole circuit of capital, in motion. */
export default function Atlas() {
  const { snapshot, stale } = useSnapshot();
  const [speed, setSpeed] = useState(1);

  return (
    <div className="atlas-shell">
      <header className="atlas-top">
        <span className="brand">Capital Simulator — Atlas</span>
        <nav className="nav">
          <a href="#/">Atlas</a>
          <a href="#/chapters">Chapters</a>
        </nav>
      </header>

      <div className="atlas-body">
        <aside className="atlas-rail">
          {snapshot ? (
            <VitalSigns vitals={snapshot.aggregate} capitalCount={snapshot.capitals.length} />
          ) : (
            <p style={{ color: "var(--ink-muted)", fontSize: 13 }}>Loading the field…</p>
          )}
        </aside>

        {snapshot ? (
          <CircuitField snapshot={snapshot} speed={speed} />
        ) : (
          <div className="atlas-field-wrap" />
        )}
      </div>

      <footer className="atlas-bottom">
        <TickHeartbeat
          tick={snapshot?.tick ?? 0}
          running={snapshot?.running ?? false}
          stale={stale}
          speed={speed}
          onSpeed={setSpeed}
        />
      </footer>
    </div>
  );
}
```

- [ ] **Step 2: Typecheck + build**

Run: `cd web && npm run lint && npm run build`
Expected: success.

- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/Atlas.tsx
git commit --no-gpg-sign -m "feat(atlas): Observatory shell composing field, rail, transport"
```

---

## Task B10: Hash router split (Atlas = `/`, dashboard = `#/chapters`)

**Files:**
- Create: `web/src/Root.tsx`
- Modify: `web/src/main.tsx` (line 3 import, line ~13 render)
- Modify: `web/src/App.tsx` (topbar — add an Atlas link, ~line 46)

- [ ] **Step 1: Write `Root.tsx`**

Create `web/src/Root.tsx`:

```tsx
import { useEffect, useState } from "react";
import App from "./App";
import Atlas from "./atlas/Atlas";

/** Reads the current hash route, e.g. "#/chapters" → "/chapters". */
function useHashRoute(): string {
  const read = () => window.location.hash.replace(/^#/, "") || "/";
  const [route, setRoute] = useState(read);
  useEffect(() => {
    const onHash = () => setRoute(read());
    window.addEventListener("hashchange", onHash);
    return () => window.removeEventListener("hashchange", onHash);
  }, []);
  return route;
}

/** Minimal two-view router: Atlas at "/", the chapter dashboard at "/chapters". */
export default function Root() {
  const route = useHashRoute();
  if (route.startsWith("/chapters")) return <App />;
  return <Atlas />;
}
```

- [ ] **Step 2: Point `main.tsx` at `Root`**

In `web/src/main.tsx`, change the import on line 3 from `import App from "./App";` to:

```tsx
import Root from "./Root";
```

and change the render block to use `<Root />`:

```tsx
createRoot(rootEl).render(
  <StrictMode>
    <Root />
  </StrictMode>
);
```

- [ ] **Step 3: Add an "← Atlas" link to the dashboard topbar**

In `web/src/App.tsx`, inside the `<header className="topbar">` block, immediately after the `<span className="topbar-logo">Capital Simulator</span>` line, add:

```tsx
        <a href="#/" style={{ color: "var(--ink-muted)", textDecoration: "none", fontSize: 13, marginLeft: 12 }}>← Atlas</a>
```

- [ ] **Step 4: Typecheck + build**

Run: `cd web && npm run lint && npm run build`
Expected: success.

- [ ] **Step 5: Commit**

```bash
git add web/src/Root.tsx web/src/main.tsx web/src/App.tsx
git commit --no-gpg-sign -m "feat(atlas): hash router — Atlas landing, chapters drill-down"
```

---

## Task B11: End-to-end acceptance on the booted stack (Playwright MCP)

**Files:** none (manual verification per CLAUDE.md "always boot stack for panel E2E").

- [ ] **Step 1: Boot the full stack**

Run:
```bash
docker compose up --build -d
cd web && npm run dev   # serves http://localhost:5173 with /api → gateway
```

- [ ] **Step 2: Drive the page with Playwright MCP and assert**

Using the Playwright MCP browser tools:
1. `browser_navigate` → `http://localhost:5173/` (Atlas is the default route).
2. `browser_snapshot` → assert `[data-testid="atlas-field"]` exists and contains **4 orbits** (4 `.atlas-orbit` nodes: Spinning Mill + 3 seeded branches).
3. Assert the rail shows a non-empty "Average rate of profit · p̄′" value.
4. Click the transport ▶ button. `browser_network_requests` → assert `POST /api/v1/engine/start` returned 200.
5. Wait ~6s, `browser_snapshot` → assert the `turn N` counter advanced and `[data-testid="atlas-ecg"]` polyline has multiple points.
6. Click the "Chapters" nav link → assert the hash becomes `#/chapters` and the existing dashboard renders.
7. Click "← Atlas" → assert the field returns.

Expected: all assertions pass. Capture a screenshot (`browser_take_screenshot`) for the PR.

- [ ] **Step 3: Tear down**

Run: `docker compose down`

- [ ] **Step 4: Final full check + commit any fixes**

Run:
```bash
make vet test build
cd web && npm run lint && npm run build
```
Expected: all PASS. Commit any fixes with `--no-gpg-sign`.

---

# Done criteria (Phase 1)

- `GET /v1/observatory/snapshot` returns the field + aggregate p̄′ = ΣS/ΣC; gateway proxies it.
- Fresh MySQL volume comes up with ≥4 capitals of differing composition (no empty field).
- `/` renders the Atlas field of faithful orbits in continuous motion; transport starts/stops the real engine; the heartbeat ECG advances with ticks.
- `#/chapters` still reaches the full 106-panel dashboard; links bridge both ways.
- `make vet test build` and `npm run lint && npm run build` pass; Playwright acceptance (B11) passes on the booted stack.

# Deferred (later phases)
- **Phase 2 — Inspector:** click an orbit → single faithful orbit + c+v+s ledger + status + "open in chapters".
- **Phase 3 — Levers:** rate of surplus-value / accumulation / organic-composition drift wired to `POST /v1/accumulation/*`.
- Eased radius/arc tweening between snapshots; richer vitals (s′, Dept I/II, accumulation); backend tick-interval setter; SSE/WebSocket push.
