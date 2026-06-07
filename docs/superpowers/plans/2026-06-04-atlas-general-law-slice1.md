# Atlas General Law — Slice 1 (Real Growth + Corrected Motion) Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the existing Atlas field genuinely move and grow — capitals accumulate surplus and spiral outward, and each orbit's value-dot travels its static ring at the real turnover rate, lingering in production.

**Architecture:** A new `AccumulationTicker` in `simulation-engine` capitalises a share (α) of each `IndustrialCapital`'s surplus into `total_pence` every scheduler pass (rescaling its M/P/C stage distribution — the spiral). The snapshot carries a per-capital `turnover_number` (new column). `Orbit.tsx` is rewritten to keep the ring static, drive dot travel with a requestAnimationFrame pacing function that laps once per turnover and lingers in P, and ease the orbit's scale toward the growing magnitude.

**Tech Stack:** Go 1.25 (`engine.Ticker`, `database/sql` tx), goose migration, React 18 + TS, SVG + requestAnimationFrame + CSS transform. No new dependencies.

**Spec:** `docs/superpowers/specs/2026-06-04-atlas-general-law-design.md` (Slice 1). Branch: `feature/atlas-observatory` (continues Phase 1).

**Scope notes (tightened from the spec for YAGNI):** Slice 1 carries only `turnover_number` in the snapshot (not per-capital `v`/`constant` or a stage-timing field — those land in Slice 2 where s/v is rendered). Growth is **linear** (α·s of the seeded surplus reinvested each pass); the *compounding* feedback (s = v·s′, rising composition) is Slice 2's General Law. Dot pacing derives from the existing M/P/C arc split (no new backend timing field).

---

## File Structure

**Backend (`services/simulation-engine/`):**
- Create `internal/store/migrations/00065_atlas_turnover_number.sql` — add `turnover_number` column + seed values.
- Modify `internal/circulation/field.go` — add `TurnoverNumber` to `FieldCapital`.
- Modify `internal/store/store.go` — add `AccumulateCapital` to `IndustrialCapitalStore`.
- Modify `internal/store/memory.go` — `FieldSnapshot` sets `TurnoverNumber`; add `AccumulateCapital`.
- Modify `internal/store/mysql.go` — `FieldSnapshot` reads the column; add `AccumulateCapital` (tx).
- Modify `internal/transport/httpapi/observatory_handler.go` — add `turnover_number` to the DTO.
- Create `internal/engine/accumulation_ticker.go` — the ticker.
- Create `internal/engine/accumulation_ticker_test.go` — ticker test.
- Modify `cmd/simulation-engine/main.go` — register the ticker + `accumulationRateBP()`.

**Frontend (`web/src/`):**
- Modify `atlas/animation.ts` — add `lapRateFor`, `pacedAngle`, `targetScale`.
- Modify `atlas/Orbit.tsx` — static ring, rAF dot pacing (linger in P), CSS-scale growth.
- Modify `types.ts` — add `turnover_number` to `FieldCapital`.

---

# GROUP A — Backend

## Task A1: Per-capital `turnover_number` through the snapshot

**Files:**
- Create: `services/simulation-engine/internal/store/migrations/00065_atlas_turnover_number.sql`
- Modify: `services/simulation-engine/internal/circulation/field.go`
- Modify: `services/simulation-engine/internal/store/memory.go` (`FieldSnapshot`)
- Modify: `services/simulation-engine/internal/store/mysql.go` (`FieldSnapshot`)
- Modify: `services/simulation-engine/internal/transport/httpapi/observatory_handler.go`
- Test: `services/simulation-engine/internal/transport/httpapi/observatory_handler_test.go` (extend)

- [ ] **Step 1: Add `TurnoverNumber` to `FieldCapital`**

In `internal/circulation/field.go`, add the field to the struct:

```go
	Status          IndustrialCapitalStatus
	TurnoverNumber  int64 // laps per period; orbit dot lap-rate (Vol. II Ch. 7)
}
```

- [ ] **Step 2: Write the migration**

Create `internal/store/migrations/00065_atlas_turnover_number.sql`:

```sql
-- +goose Up
-- Atlas: per-capital turnover number drives the orbit dot's lap rate (Vol. II Ch.7).
ALTER TABLE industrial_capitals ADD COLUMN turnover_number INT NOT NULL DEFAULT 1;
UPDATE industrial_capitals SET turnover_number = 5 WHERE id = '5eed000000000000000401'; -- Spinning Mill (5×/yr)
UPDATE industrial_capitals SET turnover_number = 8 WHERE id = '5eed000000000000009001'; -- Tannery (light, fast)
UPDATE industrial_capitals SET turnover_number = 2 WHERE id = '5eed000000000000009011'; -- Steelworks (heavy fixed capital, slow)
UPDATE industrial_capitals SET turnover_number = 4 WHERE id = '5eed000000000000009021'; -- Textile

-- +goose Down
ALTER TABLE industrial_capitals DROP COLUMN turnover_number;
```

- [ ] **Step 3: Memory `FieldSnapshot` sets a default turnover**

In `internal/store/memory.go`, in `FieldSnapshot`, set `TurnoverNumber` on the constructed `FieldCapital` (memory has no column → default 1):

```go
		fc := circulation.FieldCapital{
			ID:             id,
			TotalPence:     ic.TotalPence,
			MoneyPence:     ic.TotalPence, // default when no distribution recorded
			Status:         ic.Status,
			TurnoverNumber: 1,
		}
```

- [ ] **Step 4: MySQL `FieldSnapshot` reads the column**

In `internal/store/mysql.go`, in `FieldSnapshot`: add `ic.turnover_number` to the SELECT (first line of the column list) and scan it. The SELECT becomes:

```go
	const q = `
SELECT ic.id, ic.total_pence, ic.status, ic.turnover_number,
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
```

and update the scan target block to include `turnover`:

```go
		var (
			id, status        string
			total, turnover   int64
			money, prod, comm sql.NullInt64
			demand, excess    sql.NullInt64
		)
		if err := rows.Scan(&id, &total, &status, &turnover, &money, &prod, &comm, &demand, &excess); err != nil {
			return nil, err
		}
		fc := circulation.FieldCapital{
			ID:             circulation.IndustrialCapitalID(id),
			TotalPence:     circulation.Pence(total),
			Status:         circulation.IndustrialCapitalStatus(status),
			MoneyPence:     circulation.Pence(total),
			TurnoverNumber: turnover,
		}
```

(The rest of the loop — the `money.Valid` / `demand.Valid` blocks — is unchanged.)

- [ ] **Step 5: Add `turnover_number` to the snapshot DTO**

In `internal/transport/httpapi/observatory_handler.go`, add the field to `fieldCapitalDTO` and map it:

```go
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
```

In the `for i, fc := range field` loop, add to the struct literal:

```go
			Status:          string(fc.Status),
			TurnoverNumber:  fc.TurnoverNumber,
		}
```

- [ ] **Step 6: Extend the handler test to assert the field is present**

In `internal/transport/httpapi/observatory_handler_test.go`, inside `TestGetObservatorySnapshot`, after the existing assertions add:

```go
	for _, c := range resp.Capitals {
		if c.TurnoverNumber < 1 {
			t.Errorf("capital %s turnover_number = %d, want >= 1 (memory default)", c.ID, c.TurnoverNumber)
		}
	}
```

- [ ] **Step 7: Run tests**

Run: `cd services/simulation-engine && go test ./internal/store/ ./internal/transport/httpapi/ -run 'FieldSnapshot|ObservatorySnapshot' -v`
Expected: PASS (memory path; turnover defaults to 1).

- [ ] **Step 8: Commit**

```bash
git add services/simulation-engine/internal/circulation/field.go \
        services/simulation-engine/internal/store/migrations/00065_atlas_turnover_number.sql \
        services/simulation-engine/internal/store/memory.go \
        services/simulation-engine/internal/store/mysql.go \
        services/simulation-engine/internal/transport/httpapi/observatory_handler.go \
        services/simulation-engine/internal/transport/httpapi/observatory_handler_test.go
git commit --no-gpg-sign -m "feat(atlas): carry per-capital turnover_number in the snapshot"
```

---

## Task A2: `AccumulateCapital` store method (the spiral)

**Files:**
- Modify: `services/simulation-engine/internal/store/store.go` (`IndustrialCapitalStore`)
- Modify: `services/simulation-engine/internal/store/memory.go` (append method)
- Modify: `services/simulation-engine/internal/store/mysql.go` (append method)
- Test: `services/simulation-engine/internal/store/accumulate_test.go`

- [ ] **Step 1: Write the failing memory test**

Create `services/simulation-engine/internal/store/accumulate_test.go`:

```go
package store_test

import (
	"context"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func TestMemoryAccumulateCapital(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	ic, err := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 500000, EconomyMode: circulation.EconomyMoney,
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := m.Snapshot(ctx, ic.ID, circulation.StageDistribution{
		MoneyPence: 100000, ProductionPence: 300000, CommodityPence: 100000,
	}); err != nil {
		t.Fatalf("snapshot: %v", err)
	}

	// Accumulate £1,000 (100000 pence) into the capital.
	grown, err := m.AccumulateCapital(ctx, ic.ID, 100000)
	if err != nil {
		t.Fatalf("accumulate: %v", err)
	}
	if grown.TotalPence != 600000 {
		t.Errorf("total = %d, want 600000", grown.TotalPence)
	}

	// The field snapshot's latest distribution rescales, preserving proportions
	// (money 20%, production 60%, commodity 20% of 600000) and summing to total.
	field, _ := m.FieldSnapshot(ctx)
	var fc circulation.FieldCapital
	for _, f := range field {
		if f.ID == ic.ID {
			fc = f
		}
	}
	if fc.TotalPence != 600000 {
		t.Fatalf("field total = %d, want 600000", fc.TotalPence)
	}
	if got := fc.MoneyPence + fc.ProductionPence + fc.CommodityPence; got != 600000 {
		t.Errorf("distribution sum = %d, want 600000", got)
	}
	if fc.ProductionPence != 360000 {
		t.Errorf("production = %d, want 360000 (60%% of 600000)", fc.ProductionPence)
	}

	// Missing capital → ErrNotFound.
	if _, err := m.AccumulateCapital(ctx, circulation.IndustrialCapitalID("nope"), 1000); err != store.ErrNotFound {
		t.Errorf("missing capital err = %v, want ErrNotFound", err)
	}
}
```

- [ ] **Step 2: Run to confirm it fails**

Run: `cd services/simulation-engine && go test ./internal/store/ -run TestMemoryAccumulateCapital`
Expected: FAIL — `m.AccumulateCapital undefined`.

- [ ] **Step 3: Add the interface method**

In `internal/store/store.go`, inside `type IndustrialCapitalStore interface { ... }` (after `FieldSnapshot`), add:

```go
	// AccumulateCapital capitalises deltaPence into the capital: total_pence grows
	// and a new StageDistribution is appended, rescaled to the new total while
	// preserving the latest M/P/C proportions (the spiral of accumulation). A
	// non-positive delta is a no-op. ErrNotFound if the capital does not exist.
	AccumulateCapital(ctx context.Context, id circulation.IndustrialCapitalID, deltaPence circulation.Pence) (circulation.IndustrialCapital, error)
```

- [ ] **Step 4: Implement in `memory.go`**

Append to `internal/store/memory.go` (IndustrialCapital section):

```go
// AccumulateCapital implements IndustrialCapitalStore (the spiral of accumulation).
func (m *Memory) AccumulateCapital(_ context.Context, id circulation.IndustrialCapitalID, delta circulation.Pence) (circulation.IndustrialCapital, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ic, ok := m.industrialCapitals[id]
	if !ok {
		return circulation.IndustrialCapital{}, ErrNotFound
	}
	if delta <= 0 {
		return ic, nil
	}
	oldTotal := ic.TotalPence
	newTotal := oldTotal + delta

	var money, prod, comm circulation.Pence
	if sds := m.stageDistributions[id]; len(sds) > 0 && oldTotal > 0 {
		last := sds[len(sds)-1]
		money = last.MoneyPence * newTotal / oldTotal
		prod = last.ProductionPence * newTotal / oldTotal
		comm = newTotal - money - prod // absorb integer rounding so the sum == newTotal
	} else {
		money = newTotal // default: all value sits as money
	}

	ic.TotalPence = newTotal
	ic.UpdatedAt = m.now().UTC()
	m.industrialCapitals[id] = ic

	m.stageDistributions[id] = append(m.stageDistributions[id], circulation.StageDistribution{
		ID:                  circulation.NewStageDistributionID(),
		IndustrialCapitalID: id,
		At:                  m.now().UTC(),
		MoneyPence:          money,
		ProductionPence:     prod,
		CommodityPence:      comm,
	})
	return ic, nil
}
```

- [ ] **Step 5: Implement in `mysql.go`**

Append to `internal/store/mysql.go` (receiver `m *MySQL`, handle `m.db`; `database/sql` and `time` are imported; ID constructor `circulation.NewStageDistributionID()`):

```go
// AccumulateCapital implements IndustrialCapitalStore (the spiral of accumulation).
func (m *MySQL) AccumulateCapital(ctx context.Context, id circulation.IndustrialCapitalID, delta circulation.Pence) (circulation.IndustrialCapital, error) {
	if delta <= 0 {
		return m.GetIndustrialCapital(ctx, id)
	}
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return circulation.IndustrialCapital{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var oldTotal int64
	err = tx.QueryRowContext(ctx, `SELECT total_pence FROM industrial_capitals WHERE id = ? FOR UPDATE`, string(id)).Scan(&oldTotal)
	if err == sql.ErrNoRows {
		return circulation.IndustrialCapital{}, ErrNotFound
	}
	if err != nil {
		return circulation.IndustrialCapital{}, err
	}
	newTotal := oldTotal + int64(delta)

	var money, prod, comm int64
	var lm, lp, lc sql.NullInt64
	err = tx.QueryRowContext(ctx,
		`SELECT money_pence, production_pence, commodity_pence FROM stage_distributions
		 WHERE industrial_capital_id = ? ORDER BY at_time DESC LIMIT 1`, string(id)).Scan(&lm, &lp, &lc)
	if err != nil && err != sql.ErrNoRows {
		return circulation.IndustrialCapital{}, err
	}
	if lm.Valid && oldTotal > 0 {
		money = lm.Int64 * newTotal / oldTotal
		prod = lp.Int64 * newTotal / oldTotal
		comm = newTotal - money - prod
	} else {
		money = newTotal
	}

	now := time.Now().UTC()
	if _, err = tx.ExecContext(ctx,
		`UPDATE industrial_capitals SET total_pence = ?, updated_at = ? WHERE id = ?`,
		newTotal, now, string(id)); err != nil {
		return circulation.IndustrialCapital{}, err
	}
	if _, err = tx.ExecContext(ctx,
		`INSERT INTO stage_distributions (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		string(circulation.NewStageDistributionID()), string(id), now, money, prod, comm); err != nil {
		return circulation.IndustrialCapital{}, err
	}
	if err = tx.Commit(); err != nil {
		return circulation.IndustrialCapital{}, err
	}
	return m.GetIndustrialCapital(ctx, id)
}
```

- [ ] **Step 6: Run the test to verify it passes**

Run: `cd services/simulation-engine && go test ./internal/store/ -run TestMemoryAccumulateCapital -v`
Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add services/simulation-engine/internal/store/store.go \
        services/simulation-engine/internal/store/memory.go \
        services/simulation-engine/internal/store/mysql.go \
        services/simulation-engine/internal/store/accumulate_test.go
git commit --no-gpg-sign -m "feat(atlas): AccumulateCapital store method (the spiral)"
```

---

## Task A3: The `AccumulationTicker` + wiring

**Files:**
- Create: `services/simulation-engine/internal/engine/accumulation_ticker.go`
- Test: `services/simulation-engine/internal/engine/accumulation_ticker_test.go`
- Modify: `services/simulation-engine/cmd/simulation-engine/main.go`

- [ ] **Step 1: Write the failing ticker test**

Create `services/simulation-engine/internal/engine/accumulation_ticker_test.go`:

```go
package engine_test

import (
	"context"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func TestAccumulationTickerGrowsCapital(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	ic, _ := m.CreateIndustrialCapital(ctx, circulation.IndustrialCapital{
		TotalPence: 500000, EconomyMode: circulation.EconomyMoney,
	})
	_, _ = m.Snapshot(ctx, ic.ID, circulation.StageDistribution{
		MoneyPence: 100000, ProductionPence: 300000, CommodityPence: 100000,
	})
	_, _ = m.RecordSupplyDemand(ctx, circulation.SupplyDemandImbalance{
		IndustrialCapitalID: ic.ID, Period: "1871",
		DemandPence: 400000, SupplyPence: 480000, ExcessPence: 80000, // s = 80000
	})

	// α = 5000 bp (50% of surplus reinvested per pass).
	tk := engine.NewAccumulationTicker(m, 5000)
	if tk.Name() != "accumulation" {
		t.Fatalf("name = %q", tk.Name())
	}

	advanced, err := tk.Tick(ctx)
	if err != nil {
		t.Fatalf("tick: %v", err)
	}
	if advanced != 1 {
		t.Fatalf("advanced = %d, want 1", advanced)
	}

	// total grew by α·s = 0.5 * 80000 = 40000 → 540000.
	got, _ := m.GetIndustrialCapital(ctx, ic.ID)
	if got.TotalPence != 540000 {
		t.Errorf("total after one pass = %d, want 540000", got.TotalPence)
	}

	// A second pass grows it again (linear in Slice 1).
	_, _ = tk.Tick(ctx)
	got, _ = m.GetIndustrialCapital(ctx, ic.ID)
	if got.TotalPence != 580000 {
		t.Errorf("total after two passes = %d, want 580000", got.TotalPence)
	}
}
```

- [ ] **Step 2: Run to confirm it fails**

Run: `cd services/simulation-engine && go test ./internal/engine/ -run TestAccumulationTicker`
Expected: FAIL — `engine.NewAccumulationTicker undefined`.

- [ ] **Step 3: Write the ticker**

Create `internal/engine/accumulation_ticker.go`:

```go
package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

// CapitalAccumulator is the store subset the AccumulationTicker needs. It is
// satisfied structurally by *store.Memory and *store.MySQL.
type CapitalAccumulator interface {
	FieldSnapshot(ctx context.Context) ([]circulation.FieldCapital, error)
	AccumulateCapital(ctx context.Context, id circulation.IndustrialCapitalID, deltaPence circulation.Pence) (circulation.IndustrialCapital, error)
}

// AccumulationTicker capitalises a share (alphaBP basis points) of each
// industrial capital's surplus back into it every scheduler pass — the spiral
// of accumulation (Vol. I Ch. 24). It implements Ticker.
type AccumulationTicker struct {
	store   CapitalAccumulator
	alphaBP int64
}

// NewAccumulationTicker constructs the ticker. alphaBP is the accumulation rate
// in basis points (5000 = reinvest 50% of surplus per pass); negative values
// are treated as 0 (no accumulation).
func NewAccumulationTicker(s CapitalAccumulator, alphaBP int64) *AccumulationTicker {
	if alphaBP < 0 {
		alphaBP = 0
	}
	return &AccumulationTicker{store: s, alphaBP: alphaBP}
}

// Name identifies this ticker in status and the audit log.
func (t *AccumulationTicker) Name() string { return "accumulation" }

// Tick capitalises alphaBP·surplus into every capital. It continues past a
// single capital's failure and returns the joined error.
func (t *AccumulationTicker) Tick(ctx context.Context) (int, error) {
	field, err := t.store.FieldSnapshot(ctx)
	if err != nil {
		return 0, err
	}
	var advanced int
	var errs []error
	for _, fc := range field {
		if fc.SurplusPence <= 0 || t.alphaBP <= 0 {
			continue
		}
		delta := circulation.Pence(int64(fc.SurplusPence) * t.alphaBP / 10000)
		if delta <= 0 {
			continue
		}
		if _, err := t.store.AccumulateCapital(ctx, fc.ID, delta); err != nil {
			errs = append(errs, fmt.Errorf("capital %s: %w", fc.ID, err))
			continue
		}
		advanced++
	}
	return advanced, errors.Join(errs...)
}
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `cd services/simulation-engine && go test ./internal/engine/ -run TestAccumulationTicker -v`
Expected: PASS.

- [ ] **Step 5: Register the ticker + accumulation-rate helper**

In `cmd/simulation-engine/main.go`, add `engine.NewAccumulationTicker(st, accumulationRateBP())` to the scheduler's ticker slice:

```go
	scheduler := engine.NewScheduler(tickInterval(), []engine.Ticker{
		engine.NewFactoryTicker(st),
		engine.NewReproductionTicker(st),
		engine.NewPiecePriceTicker(engine.NewFactoryProductivitySource(st), repricer),
		engine.NewAccumulationTicker(st, accumulationRateBP()),
	}, st, logger)
```

Add the helper near `tickInterval()` (find `func tickInterval()` and add beside it):

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

(If `strconv` is not yet imported in `main.go`, add it to the import block.)

- [ ] **Step 6: Build + full check**

Run: `cd services/simulation-engine && go build ./... && cd ../.. && make vet test build`
Expected: PASS — all packages, all six binaries.

- [ ] **Step 7: Commit**

```bash
git add services/simulation-engine/internal/engine/accumulation_ticker.go \
        services/simulation-engine/internal/engine/accumulation_ticker_test.go \
        services/simulation-engine/cmd/simulation-engine/main.go
git commit --no-gpg-sign -m "feat(atlas): AccumulationTicker — capitals spiral each pass"
```

---

# GROUP B — Frontend

## Task B1: Pacing + growth math

**Files:**
- Modify: `web/src/atlas/animation.ts` (append)
- Modify: `web/src/types.ts` (add `turnover_number`)

- [ ] **Step 1: Add `turnover_number` to the TS type**

In `web/src/types.ts`, in `interface FieldCapital`, add:

```ts
  status: string;
  turnover_number: number;
}
```

- [ ] **Step 2: Append pacing/growth helpers to `animation.ts`**

Add to the end of `web/src/atlas/animation.ts`:

```ts
/** Laps per second for the dot: turnover_number scaled to a watchable base tempo. */
export function lapRateFor(turnoverNumber: number, speed: number): number {
  const n = turnoverNumber > 0 ? turnoverNumber : 1;
  // One turnover ≈ 8s at speed 1; higher turnover = proportionally faster.
  return (n / 8) * Math.max(1, speed);
}

/**
 * Map lap-progress p∈[0,1) to an angle (radians, 0 at top, clockwise) over arcs
 * sized [fm, fp, fc] (summing to 1), spending extra time in production (P) by
 * `lingerP`× so the dot visibly slows there. Pure + deterministic.
 */
export function pacedAngle(
  p: number,
  fm: number,
  fp: number,
  fc: number,
  lingerP = 2.4
): number {
  // Time weights: production over-weighted, then normalised to 1.
  const wm = fm;
  const wp = fp * lingerP;
  const wc = fc;
  const wsum = wm + wp + wc || 1;
  const tm = wm / wsum;
  const tp = wp / wsum;
  // Cumulative ANGLE boundaries (fraction of circle) for M, P, C.
  const aM = fm;
  const aP = fm + fp;
  let frac: number; // fraction around the circle, 0..1
  const q = ((p % 1) + 1) % 1;
  if (q < tm) {
    frac = (q / tm) * aM; // through M
  } else if (q < tm + tp) {
    frac = aM + ((q - tm) / tp) * (aP - aM); // through P (slow)
  } else {
    frac = aP + ((q - tm - tp) / (1 - tm - tp)) * (1 - aP); // through C
  }
  return frac * 2 * Math.PI; // 0 at top, clockwise (caller offsets)
}

/** Target CSS scale for an orbit of `totalPence` relative to the field max. */
export function targetScale(totalPence: number, maxTotal: number): number {
  if (maxTotal <= 0) return 1;
  return Math.max(0.4, Math.sqrt(Math.max(0, totalPence) / maxTotal));
}
```

- [ ] **Step 3: Typecheck**

Run: `cd web && npm run lint`
Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add web/src/atlas/animation.ts web/src/types.ts
git commit --no-gpg-sign -m "feat(atlas): pacing + growth math (linger in P, turnover lap rate)"
```

---

## Task B2: Rewrite `Orbit.tsx` (static ring · paced dots · growth)

**Files:**
- Modify: `web/src/atlas/Orbit.tsx` (full replacement)

- [ ] **Step 1: Replace `Orbit.tsx`**

Replace the entire contents of `web/src/atlas/Orbit.tsx` with:

```tsx
import { useEffect, useRef } from "react";
import type { FieldCapital } from "../types";
import { arcFractions, lapRateFor, pacedAngle, targetScale } from "./animation";

interface OrbitProps {
  capital: FieldCapital;
  maxTotal: number;
  /** Animation tempo multiplier (1 = base). */
  speed: number;
}

const GOLD = "#c8a240";
const RED = "#c0392b";
const LEAD = "#4a5a8a";
const BASE_R = 60; // base ring radius in svg units; growth applied via CSS scale
const PAD = 12;
const DOT_COUNT = 3;

/** One faithful orbit: a STILL ring of M/P/C arcs; value-dots travel it, lapping
 *  once per turnover and lingering in production; the whole orbit scales toward
 *  its growing magnitude. */
export function Orbit({ capital, maxTotal, speed }: OrbitProps) {
  const size = (BASE_R + PAD) * 2;
  const c = size / 2;
  const sw = 11;
  const [fm, fp] = arcFractions(capital);
  const fc = 1 - fm - fp;
  const gap = 1.5;
  const lm = Math.max(0, fm * 100 - gap);
  const lp = Math.max(0, fp * 100 - gap);
  const lc = Math.max(0, fc * 100 - gap);

  const wrapRef = useRef<HTMLDivElement>(null);
  const dotRefs = useRef<(SVGCircleElement | null)[]>([]);
  // Keep latest values in a ref so the rAF loop reads fresh data without restart.
  const live = useRef({ fm, fp, fc, lap: lapRateFor(capital.turnover_number, speed) });
  live.current = { fm, fp, fc, lap: lapRateFor(capital.turnover_number, speed) };

  // Growth: set the wrapper's CSS scale toward the target (CSS transition eases it).
  useEffect(() => {
    if (wrapRef.current) {
      wrapRef.current.style.transform = `scale(${targetScale(capital.total_pence, maxTotal).toFixed(3)})`;
    }
  }, [capital.total_pence, maxTotal]);

  // Dots travel the ring; rAF updates positions imperatively (no re-render).
  useEffect(() => {
    let raf = 0;
    let prev = performance.now();
    const phase = [0, 1 / DOT_COUNT, 2 / DOT_COUNT]; // evenly spaced
    const tick = (now: number) => {
      const dt = Math.min(0.05, (now - prev) / 1000);
      prev = now;
      const cur = live.current;
      for (let i = 0; i < DOT_COUNT; i++) {
        phase[i] = (phase[i] + dt * cur.lap) % 1;
        const a = pacedAngle(phase[i], cur.fm, cur.fp, cur.fc) - Math.PI / 2; // 0 at top
        const el = dotRefs.current[i];
        if (el) {
          el.setAttribute("cx", (c + BASE_R * Math.cos(a)).toFixed(2));
          el.setAttribute("cy", (c + BASE_R * Math.sin(a)).toFixed(2));
        }
      }
      raf = requestAnimationFrame(tick);
    };
    raf = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(raf);
  }, [c]);

  const halted = capital.status === "halted";

  return (
    <div
      ref={wrapRef}
      className={"atlas-orbit" + (halted ? " halted" : "")}
      style={{ transition: "transform 1.6s ease-out", willChange: "transform" }}
      title={`${capital.id.slice(0, 8)} · total ${capital.total_pence} · turnover ${capital.turnover_number}× · ${capital.status}`}
    >
      <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
        <circle cx={c} cy={c} r={BASE_R} fill="none" stroke="#161922" strokeWidth={sw} />
        <g transform={`rotate(-90 ${c} ${c})`}>
          <circle cx={c} cy={c} r={BASE_R} fill="none" stroke={GOLD} strokeWidth={sw}
            pathLength={100} strokeDasharray={`${lm} ${100 - lm}`} strokeDashoffset={0} />
          <circle cx={c} cy={c} r={BASE_R} fill="none" stroke={RED} strokeWidth={sw}
            pathLength={100} strokeDasharray={`${lp} ${100 - lp}`} strokeDashoffset={-(fm * 100)} />
          <circle cx={c} cy={c} r={BASE_R} fill="none" stroke={LEAD} strokeWidth={sw}
            pathLength={100} strokeDasharray={`${lc} ${100 - lc}`} strokeDashoffset={-((fm + fp) * 100)} />
        </g>
        {Array.from({ length: DOT_COUNT }).map((_, i) => (
          <circle key={i} ref={(el) => { dotRefs.current[i] = el; }}
            cx={c} cy={c - BASE_R} r={3.2} fill="#f4ecd8"
            style={{ filter: "drop-shadow(0 0 3px rgba(244,236,216,.8))" }} />
        ))}
      </svg>
    </div>
  );
}
```

- [ ] **Step 2: Typecheck + build**

Run: `cd web && npm run lint && npm run build`
Expected: success.

- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/Orbit.tsx
git commit --no-gpg-sign -m "feat(atlas): orbit — still ring, paced dots (linger in P), growth scale"
```

---

## Task B3: End-to-end acceptance on the booted stack (Playwright MCP)

**Files:** none (per CLAUDE.md "always boot stack for panel E2E").

- [ ] **Step 1: Boot fresh (new migration 00065 must apply)**

Run:
```bash
docker compose down -v
docker compose up --build -d
```
Wait ~3 min for MySQL + migrations. Confirm the snapshot carries turnover and that totals climb (start the engine first if it is not auto-running):
```bash
curl -s -X POST http://localhost:8080/v1/engine/start >/dev/null
curl -s http://localhost:8080/v1/observatory/snapshot | head -c 400; echo
# wait, then read again
curl -s http://localhost:8080/v1/observatory/snapshot | head -c 400; echo
```
Expected: `turnover_number` present per capital; `total_pence` for `5eed...0401` is **larger** in the second read (accumulation is running).

- [ ] **Step 2: Drive the page with Playwright MCP**

1. `browser_navigate` → `http://localhost:5173/`.
2. If transport shows ▶, click it (start the engine) so accumulation runs.
3. `browser_evaluate` → read the first `.atlas-orbit` wrapper's computed `transform` scale, wait ~20s, re-read → assert at least one orbit's scale **increased** (growth).
4. `browser_evaluate` over ~3s → sample a dot's `cy` attribute at intervals and assert it changes **non-uniformly** (slower while in the lower/P arc); at minimum that dots move and higher-turnover capitals lap more often.
5. `browser_take_screenshot` → `atlas-growth.png` for the PR.

Expected: orbits visibly grow; dots travel the still rings; higher-turnover capitals lap faster.

- [ ] **Step 3: Tear down + final check**

```bash
docker compose down
make vet test build
cd web && npm run lint && npm run build
```
Expected: all PASS.

---

# Done criteria (Slice 1)

- Capitals' `total_pence` grows each pass; orbits visibly spiral outward over ~20s.
- Rings are still; value-dots travel the circumference, lapping once per turnover and lingering in production; higher-turnover capitals lap faster.
- `make vet test build` + `npm run lint && npm run build` pass; Playwright shows growth + paced motion.

# Deferred (Slice 2 / 3)
- The hidden abode + the full General-Law feedback (composition → reserve army → wage pressure → s/v); per-capital `v`; compounding growth.
- The levers (working day / wage / accumulation rate).
