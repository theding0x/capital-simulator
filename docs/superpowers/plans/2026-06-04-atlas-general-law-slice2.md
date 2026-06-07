# Atlas General Law — Slice 2 (The Hidden Abode) Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Cross the threshold from the value-circuit surface into the *hidden abode of production* and run the General Law of Capitalist Accumulation (Vol. I Ch. 25) in motion — rising organic composition and machinery repel living labour into the industrial reserve army, the swelling reserve presses the wage below the value of labour-power and raises the rate of exploitation s/v, and the heightened surplus re-accumulates at an ever-higher composition. The class antagonism becomes the engine, surfaced as a descend-able abode with a live working day, reserve army, and immiseration trend.

**Architecture:** A new evolving aggregate, `AbodeState` (total social `c`/`v`, the wages bill, the labouring population, and the law's parameters), is advanced one period per scheduler pass by a pure `simulation.AdvanceGeneralLaw` step and a new `GeneralLawTicker`. Each pass appends a `GeneralLawPeriod` to a persisted immiseration series. `GET /v1/observatory/snapshot` gains an `abode` block (the instantaneous class state + the recent series). The frontend adds a threshold control that descends from the field of orbits into an `Abode` view (working-day bar = s/v, living labour Σv, reserve army + wage pressure, surplus extraction, and a `GeneralLawTrend` sparkline).

**Tech Stack:** Go 1.25 (pure domain step + `engine.Ticker` + `database/sql` tx), goose migration, React 18 + TS, SVG sparkline + CSS transition. No new dependencies.

**Spec:** `docs/superpowers/specs/2026-06-04-atlas-general-law-design.md` (Slice 2 = §10 item 2). Branch: `feature/atlas-observatory` (continues Slice 1).

**Key modeling decisions (pinning the spec's §11 "open details"):**
- **The abode is an evolving aggregate `AbodeState`, not a live field re-aggregation.** The spec's §4/§5 ("advance … via RunGeneralLaw", "record one aggregate GeneralLawSnapshot-style row per period") and §9 (`s/v = round(10000·Σs/Σv)`) are *both* satisfied: Σv/Σs are the abode aggregate's `variable`/`surplus`, and the series is the time dimension. A live field re-aggregation can't move (Slice-1 accumulation rescales proportionally, so composition stays flat) — the law would not be *in motion*.
- **Immiseration driver:** rising marginal composition + a per-period `displacement_rate_bp` (machinery converting a share of the wages bill to constant capital — "the labouring population … produces … the means by which it is made relatively superfluous", Ch. 25) holds `v` down while `c` grows, so the reserve army grows by mechanism. Verified to immiserate monotonically from period 0 with the seeded parameters.
- **Working day:** `necessary`/`surplus` minutes are derived from the effective s/v against a fixed `SocialWorkingDayMinutes = 600` (Ch. 9), not a lever. Levers are Slice 3.
- **Deferred to Slice 3 (out of scope here):** per-capital `variable_pence`/`constant_pence` in the `capitals` array (the abode needs only the aggregate; no orbit component consumes per-capital `v` yet) and the working-day/wage/accumulation levers.

---

## File Structure

**Backend (`services/simulation-engine/`):**
- Create `internal/simulation/abode.go` — `AbodeState`, `GeneralLawPeriod`, `AbodeReadout`, `NewAbodeState`, `AdvanceGeneralLaw` (pure).
- Create `internal/simulation/abode_test.go` — direction-of-the-law tests.
- Create `internal/store/migrations/00066_atlas_abode.sql` — `abode_state` + `general_law_periods` tables + seed row.
- Modify `internal/store/store.go` — add `AbodeStateStore` interface.
- Modify `internal/store/memory.go` — `abodeState`/`generalLawPeriods` fields + three methods.
- Modify `internal/store/mysql.go` — three methods (tx for `AdvanceAbode`).
- Create `internal/store/abode_test.go` — memory round-trip test.
- Create `internal/engine/general_law_ticker.go` — the ticker.
- Create `internal/engine/general_law_ticker_test.go` — ticker test.
- Modify `internal/transport/httpapi/observatory_handler.go` — `abode` block on the snapshot.
- Modify `internal/transport/httpapi/observatory_handler_test.go` — abode assertions.
- Modify `internal/transport/httpapi/handler.go` — `AbodeStates` dependency.
- Modify `cmd/simulation-engine/main.go` — register ticker + wire `AbodeStates`.

**Frontend (`web/src/`):**
- Modify `types.ts` — `AbodeReadout`, `GeneralLawTrendPoint`, `abode` on `ObservatorySnapshot`.
- Create `atlas/Abode.tsx` — the hidden-abode view (working day, living labour, reserve army, surplus extraction).
- Create `atlas/GeneralLawTrend.tsx` — the immiseration sparkline.
- Modify `atlas/Atlas.tsx` — threshold toggle surface↕abode.
- Modify `atlas/animation.ts` — `formatMinutes` helper + a tiny sparkline path builder.
- Modify `atlas/atlas.css` — abode + threshold styles.

---

# GROUP A — Backend domain (the law)

## Task A1: `AbodeState` + `AdvanceGeneralLaw` (pure)

**Files:**
- Create: `services/simulation-engine/internal/simulation/abode.go`
- Test: `services/simulation-engine/internal/simulation/abode_test.go`

- [ ] **Step 1: Write the failing test**

Create `services/simulation-engine/internal/simulation/abode_test.go`:

```go
package simulation

import "testing"

func TestNewAbodeStateReadout(t *testing.T) {
	t.Parallel()
	r := NewAbodeState().Readout()
	if r.TotalVariablePence != 300000 {
		t.Errorf("v = %d, want 300000", r.TotalVariablePence)
	}
	// 120 employed (v/baseWage = 300000/2500), 30 in the reserve of 150.
	if r.EmployedCount != 120 {
		t.Errorf("employed = %d, want 120", r.EmployedCount)
	}
	if r.ReserveArmyCount != 30 {
		t.Errorf("reserve = %d, want 30", r.ReserveArmyCount)
	}
	// pressure = 30/150 = 2000 bp; s' = 10000*(1+0.20) = 12000 bp.
	if r.RateOfExploitationBP != 12000 {
		t.Errorf("s/v = %d, want 12000", r.RateOfExploitationBP)
	}
	// c/v = 600000/300000 = 2.0 → 20000 bp.
	if r.OrganicCompositionBP != 20000 {
		t.Errorf("c/v = %d, want 20000", r.OrganicCompositionBP)
	}
	// surplus = v * s'/10000 = 300000 * 1.2 = 360000.
	if r.TotalSurplusPence != 360000 {
		t.Errorf("surplus = %d, want 360000", r.TotalSurplusPence)
	}
	// working day 600' splits necessary = 600*10000/22000 = 272, surplus = 328.
	if r.NecessaryLabourMinutes+r.SurplusLabourMinutes != SocialWorkingDayMinutes {
		t.Errorf("working day = %d, want %d",
			r.NecessaryLabourMinutes+r.SurplusLabourMinutes, SocialWorkingDayMinutes)
	}
	if r.SurplusLabourMinutes <= r.NecessaryLabourMinutes {
		t.Errorf("surplus labour %d should exceed necessary %d at s'=120%%",
			r.SurplusLabourMinutes, r.NecessaryLabourMinutes)
	}
	// paid wage compressed below the value of labour-power (2500).
	if r.WagePence >= 2500 {
		t.Errorf("paid wage = %d, want < 2500 (value of labour-power)", r.WagePence)
	}
}

func TestAdvanceGeneralLawImmiserates(t *testing.T) {
	t.Parallel()
	s := NewAbodeState()
	first := s.Readout()
	var last GeneralLawPeriod
	for i := 0; i < 12; i++ {
		var p GeneralLawPeriod
		s, p = AdvanceGeneralLaw(s)
		last = p
	}
	// Over twelve periods the law runs its course: the reserve army grows, the
	// organic composition rises, the rate of exploitation rises, the wage falls.
	if last.ReserveArmyCount <= first.ReserveArmyCount {
		t.Errorf("reserve army %d did not grow from %d", last.ReserveArmyCount, first.ReserveArmyCount)
	}
	if last.OrganicCompositionBP <= first.OrganicCompositionBP {
		t.Errorf("composition %d did not rise from %d", last.OrganicCompositionBP, first.OrganicCompositionBP)
	}
	if last.RateOfExploitationBP <= first.RateOfExploitationBP {
		t.Errorf("s/v %d did not rise from %d", last.RateOfExploitationBP, first.RateOfExploitationBP)
	}
	if last.WagePence >= first.WagePence {
		t.Errorf("wage %d did not fall from %d", last.WagePence, first.WagePence)
	}
	if s.Period != 12 {
		t.Errorf("period = %d, want 12", s.Period)
	}
}

func TestAdvanceGeneralLawConserves(t *testing.T) {
	t.Parallel()
	s := NewAbodeState()
	next, _ := AdvanceGeneralLaw(s)
	// Accumulation + displacement only move value between c and v and add the
	// capitalised surplus; neither v nor c goes negative.
	if next.VariablePence < 0 || next.ConstantPence < 0 {
		t.Fatalf("negative capital: c=%d v=%d", next.ConstantPence, next.VariablePence)
	}
	// Total social capital grows (α·s re-accumulated) — the spiral.
	if next.ConstantPence+next.VariablePence <= s.ConstantPence+s.VariablePence {
		t.Errorf("total capital did not grow: %d -> %d",
			s.ConstantPence+s.VariablePence, next.ConstantPence+next.VariablePence)
	}
}
```

- [ ] **Step 2: Run to confirm it fails**

Run: `cd services/simulation-engine && go test ./internal/simulation/ -run 'AbodeState|GeneralLaw' 2>&1 | head`
Expected: FAIL — `undefined: NewAbodeState` / `AdvanceGeneralLaw` / `SocialWorkingDayMinutes`.

- [ ] **Step 3: Write `abode.go`**

Create `services/simulation-engine/internal/simulation/abode.go`:

```go
package simulation

// SocialWorkingDayMinutes is the length of the aggregate social working day
// (Vol. I Ch. 9) against which the rate of surplus-value splits paid (necessary)
// from unpaid (surplus) labour. Ten hours — the post-1850 normal working day.
const SocialWorkingDayMinutes = 600

// AbodeState is the evolving aggregate class relation beneath the field of
// capitals — "the hidden abode of production, on whose threshold there hangs the
// notice 'No admittance except on business'" (Vol. I Ch. 6 fin.). It is the
// total social capital partitioned into constant (c — dead labour) and variable
// (v — the wages bill) capital, the value of labour-power (BaseWagePence), the
// labouring population (WorkerSupply), and the parameters of the law. One
// scheduler pass advances it one period via AdvanceGeneralLaw — the General Law
// of Capitalist Accumulation in motion (Vol. I Ch. 25).
type AbodeState struct {
	Period                int64 `json:"period"`
	ConstantPence         Pence `json:"constant_pence"`          // c — dead labour
	VariablePence         Pence `json:"variable_pence"`          // v — the wages bill
	BaseWagePence         int64 `json:"base_wage_pence"`         // value of labour-power, per worker
	WorkerSupply          int64 `json:"worker_supply"`           // the labouring population
	SurplusRateBaseBP     int64 `json:"surplus_rate_base_bp"`    // s′ at full employment (10000 = 100%)
	AccumulationRateBP    int64 `json:"accumulation_rate_bp"`    // α — share of surplus capitalised
	MarginalCompositionBP int64 `json:"marginal_composition_bp"` // composition c/(c+v) of NEW capital
	DisplacementRateBP    int64 `json:"displacement_rate_bp"`    // share of v repelled to c each period
	ProductivityGrowthBP  int64 `json:"productivity_growth_bp"`  // gap to full automation closed per period
	PopulationGrowthBP    int64 `json:"population_growth_bp"`    // worker-supply growth per period
}

// NewAbodeState returns the seeded initial abode — a Marx-faithful aggregate
// (c:v = 2:1, s′ = 100%) with a small initial reserve army the law expands.
// These values are mirrored in migration 00066 so the MySQL-backed and
// memory-backed economies start identically.
func NewAbodeState() AbodeState {
	return AbodeState{
		Period:                0,
		ConstantPence:         600000, // £6,000 dead labour
		VariablePence:         300000, // £3,000 wages bill
		BaseWagePence:         2500,   // £25 value of labour-power → 120 employed
		WorkerSupply:          150,    // 30 already in the reserve army
		SurplusRateBaseBP:     10000,  // 100% rate of surplus-value
		AccumulationRateBP:    5000,   // reinvest 50% of surplus
		MarginalCompositionBP: 6667,   // new capital starts at the stock composition c/(c+v)
		DisplacementRateBP:    1800,   // machinery repels 18% of the wages bill per period
		ProductivityGrowthBP:  500,    // marginal composition closes 5% of its gap to full automation
		PopulationGrowthBP:    150,    // the labouring population grows 1.5% per period
	}
}

// GeneralLawPeriod is one recorded period of the general law — a point on the
// immiseration time-series surfaced in the abode.
type GeneralLawPeriod struct {
	Period               int64 `json:"period"`
	WagePence            int64 `json:"wage_pence"`              // paid wage (price of labour, compressed)
	RateOfExploitationBP int64 `json:"rate_of_exploitation_bp"` // s/v effective
	ReserveArmyCount     int64 `json:"reserve_army_count"`
	OrganicCompositionBP int64 `json:"organic_composition_bp"` // c/v
	EmployedCount        int64 `json:"employed_count"`
}

// AbodeReadout is the instantaneous class state derived from an AbodeState — the
// surface masks exactly this. Pence/minutes/basis-points integers throughout.
type AbodeReadout struct {
	TotalVariablePence     int64 // Σv — wages = paid labour
	TotalSurplusPence      int64 // Σs — unpaid labour
	RateOfExploitationBP   int64 // s/v
	NecessaryLabourMinutes int64
	SurplusLabourMinutes   int64
	OrganicCompositionBP   int64 // c/v
	ReserveArmyCount       int64
	ReserveArmyPressureBP  int64
	EmployedCount          int64
	WagePence              int64 // paid wage after reserve-army compression
}

// reservePressureBP returns reserve / workforce in basis points, clamped to
// [0, 10000]. Zero reserve or workforce is no pressure.
func reservePressureBP(reserve, workforce int64) int64 {
	if reserve <= 0 || workforce <= 0 {
		return 0
	}
	bp := reserve * 10000 / workforce
	if bp > 10000 {
		return 10000
	}
	return bp
}

// compressWage drives the price of labour below the value of labour-power as the
// reserve army grows (Vol. I Ch. 25 §3): nominal × (1 − pressure), rounded half
// up, floored at half the value of labour-power (subsistence).
func compressWage(baseWage, pressureBP int64) int64 {
	if baseWage <= 0 || pressureBP <= 0 {
		return baseWage
	}
	w := (baseWage*(10000-pressureBP) + 5000) / 10000
	if floor := baseWage / 2; w < floor {
		return floor
	}
	return w
}

// Readout projects the current AbodeState to the instantaneous class state.
func (a AbodeState) Readout() AbodeReadout {
	oc := ComputeOrganicComposition(CapitalStock{ConstantCapital: a.ConstantPence, VariableCapital: a.VariablePence})
	employed := ComputeLabourDemand(a.ConstantPence+a.VariablePence, oc, a.BaseWagePence)
	reserve := ComputeReserveArmy(a.WorkerSupply, employed)
	pressure := reservePressureBP(reserve.Size, a.WorkerSupply)
	effRate := a.SurplusRateBaseBP * (10000 + pressure) / 10000
	surplus := int64(a.VariablePence) * effRate / 10000
	var ocBP int64
	if a.VariablePence > 0 {
		ocBP = 10000 * int64(a.ConstantPence) / int64(a.VariablePence)
	}
	necessary := int64(SocialWorkingDayMinutes) * 10000 / (10000 + effRate)
	return AbodeReadout{
		TotalVariablePence:     int64(a.VariablePence),
		TotalSurplusPence:      surplus,
		RateOfExploitationBP:   effRate,
		NecessaryLabourMinutes: necessary,
		SurplusLabourMinutes:   int64(SocialWorkingDayMinutes) - necessary,
		OrganicCompositionBP:   ocBP,
		ReserveArmyCount:       reserve.Size,
		ReserveArmyPressureBP:  pressure,
		EmployedCount:          employed.Workers,
		WagePence:              compressWage(a.BaseWagePence, pressure),
	}
}

// AdvanceGeneralLaw runs one period of the general law, returning the next state
// and the GeneralLawPeriod recorded for the immiseration series. The loop
// (Vol. I Ch. 25): the heightened surplus is re-accumulated at the marginal
// (machine-heavier) composition; machinery then repels a share of the wages bill
// into constant capital, holding v down and swelling the reserve army; the
// labouring population grows and the marginal composition climbs toward — but
// never reaches — full automation. The circle closes.
func AdvanceGeneralLaw(s AbodeState) (AbodeState, GeneralLawPeriod) {
	r := s.Readout()
	period := GeneralLawPeriod{
		Period:               s.Period + 1,
		WagePence:            r.WagePence,
		RateOfExploitationBP: r.RateOfExploitationBP,
		ReserveArmyCount:     r.ReserveArmyCount,
		OrganicCompositionBP: r.OrganicCompositionBP,
		EmployedCount:        r.EmployedCount,
	}

	next := s
	next.Period = s.Period + 1

	// Re-accumulate α·s at the marginal composition (new capital is more
	// machine-heavy than the existing stock).
	capitalised := r.TotalSurplusPence * s.AccumulationRateBP / 10000
	dc := capitalised * s.MarginalCompositionBP / 10000
	dv := capitalised - dc
	next.ConstantPence = s.ConstantPence + Pence(dc)
	next.VariablePence = s.VariablePence + Pence(dv)

	// Machinery repels living labour: a share of the wages bill is converted to
	// constant capital — the working population "produces the means by which it
	// is made relatively superfluous" (§25.3). This holds v down so the reserve
	// army grows even as total capital accumulates.
	displaced := int64(next.VariablePence) * s.DisplacementRateBP / 10000
	next.VariablePence -= Pence(displaced)
	next.ConstantPence += Pence(displaced)

	// The labouring population grows; the marginal composition asymptotes below
	// full automation (no fictional ceiling — it closes a fixed fraction of the
	// remaining gap each period).
	next.WorkerSupply = s.WorkerSupply + s.WorkerSupply*s.PopulationGrowthBP/10000
	next.MarginalCompositionBP = s.MarginalCompositionBP +
		s.ProductivityGrowthBP*(10000-s.MarginalCompositionBP)/10000

	return next, period
}
```

- [ ] **Step 4: Run the tests to verify they pass**

Run: `cd services/simulation-engine && go test ./internal/simulation/ -run 'AbodeState|GeneralLaw' -v 2>&1 | tail -20`
Expected: PASS (all three).

- [ ] **Step 5: Commit**

```bash
git add services/simulation-engine/internal/simulation/abode.go \
        services/simulation-engine/internal/simulation/abode_test.go
git commit --no-gpg-sign -m "feat(atlas): AbodeState + AdvanceGeneralLaw — the general law in motion"
```

---

# GROUP B — Backend persistence

## Task B1: `AbodeStateStore` interface + migration

**Files:**
- Create: `services/simulation-engine/internal/store/migrations/00066_atlas_abode.sql`
- Modify: `services/simulation-engine/internal/store/store.go`

- [ ] **Step 1: Write the migration**

Create `services/simulation-engine/internal/store/migrations/00066_atlas_abode.sql`:

```sql
-- +goose Up
-- Atlas Slice 2: the hidden abode. abode_state is the single evolving aggregate
-- class relation (Vol. I Ch. 25); general_law_periods is the immiseration
-- time-series the General-Law ticker appends to each scheduler pass.
CREATE TABLE abode_state (
    id                      VARCHAR(64) NOT NULL,
    period                  BIGINT      NOT NULL,
    constant_pence          BIGINT      NOT NULL,
    variable_pence          BIGINT      NOT NULL,
    base_wage_pence         BIGINT      NOT NULL,
    worker_supply           BIGINT      NOT NULL,
    surplus_rate_base_bp    BIGINT      NOT NULL,
    accumulation_rate_bp    BIGINT      NOT NULL,
    marginal_composition_bp BIGINT      NOT NULL,
    displacement_rate_bp    BIGINT      NOT NULL,
    productivity_growth_bp  BIGINT      NOT NULL,
    population_growth_bp    BIGINT      NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE general_law_periods (
    id                      VARCHAR(64) NOT NULL,
    period                  BIGINT      NOT NULL,
    wage_pence              BIGINT      NOT NULL,
    rate_of_exploitation_bp BIGINT      NOT NULL,
    reserve_army_count      BIGINT      NOT NULL,
    organic_composition_bp  BIGINT      NOT NULL,
    employed_count          BIGINT      NOT NULL,
    created_at              DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    KEY idx_glp_period (period)
);

-- The singleton abode_state, mirroring simulation.NewAbodeState().
INSERT INTO abode_state
    (id, period, constant_pence, variable_pence, base_wage_pence, worker_supply,
     surplus_rate_base_bp, accumulation_rate_bp, marginal_composition_bp,
     displacement_rate_bp, productivity_growth_bp, population_growth_bp)
VALUES
    ('5eed000000000000abode1', 0, 600000, 300000, 2500, 150,
     10000, 5000, 6667, 1800, 500, 150);

-- +goose Down
DROP TABLE IF EXISTS general_law_periods;
DROP TABLE IF EXISTS abode_state;
```

(These two tables are created wholly by this migration, so the `Down` drops them — no other table's user data is touched.)

- [ ] **Step 2: Add the store interface**

In `internal/store/store.go`, after the `IndustrialCapitalStore` interface block (ends at the line with `AccumulateCapital(... )`), add:

```go
// AbodeStateStore is the persistence contract for the Atlas hidden abode
// (Slice 2). GetAbodeState returns the single evolving aggregate, defaulting to
// simulation.NewAbodeState() when none is persisted. AdvanceAbode atomically
// writes the next state and appends one GeneralLawPeriod. ListGeneralLawPeriods
// returns the most recent periods in ascending period order (oldest first), for
// the immiseration sparkline; a non-positive limit returns every row.
type AbodeStateStore interface {
	GetAbodeState(ctx context.Context) (simulation.AbodeState, error)
	AdvanceAbode(ctx context.Context, next simulation.AbodeState, period simulation.GeneralLawPeriod) error
	ListGeneralLawPeriods(ctx context.Context, limit int) ([]simulation.GeneralLawPeriod, error)
}
```

(`simulation` is already imported in `store.go`.)

- [ ] **Step 3: Build (interface only — no impls yet, so just vet the package compiles its own file)**

Run: `cd services/simulation-engine && go build ./internal/store/ 2>&1 | head`
Expected: BUILD passes — the interface is not yet assigned to Memory/MySQL anywhere, so adding the type alone compiles. (Memory/MySQL gain the methods in B2/B3, where the compile-time `var _` assertions are added.)

- [ ] **Step 4: Commit**

```bash
git add services/simulation-engine/internal/store/migrations/00066_atlas_abode.sql \
        services/simulation-engine/internal/store/store.go
git commit --no-gpg-sign -m "feat(atlas): abode_state + general_law_periods schema + AbodeStateStore"
```

---

## Task B2: Memory implementation + round-trip test

**Files:**
- Modify: `services/simulation-engine/internal/store/memory.go`
- Test: `services/simulation-engine/internal/store/abode_test.go`

- [ ] **Step 1: Write the failing test**

Create `services/simulation-engine/internal/store/abode_test.go`:

```go
package store_test

import (
	"context"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func TestMemoryAbodeRoundTrip(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	// Unseeded → the default initial abode.
	got, err := m.GetAbodeState(ctx)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Period != 0 || got.VariablePence != 300000 {
		t.Fatalf("default abode = %+v", got)
	}

	// Advance one period and persist it + its recorded series point.
	next, period := simulation.AdvanceGeneralLaw(got)
	if err := m.AdvanceAbode(ctx, next, period); err != nil {
		t.Fatalf("advance: %v", err)
	}
	got2, _ := m.GetAbodeState(ctx)
	if got2.Period != 1 {
		t.Errorf("period after advance = %d, want 1", got2.Period)
	}

	series, err := m.ListGeneralLawPeriods(ctx, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(series) != 1 {
		t.Fatalf("series len = %d, want 1", len(series))
	}
	if series[0].Period != 1 {
		t.Errorf("series[0].Period = %d, want 1", series[0].Period)
	}

	// Ascending order is preserved across several advances.
	cur := got2
	for i := 0; i < 5; i++ {
		n, p := simulation.AdvanceGeneralLaw(cur)
		_ = m.AdvanceAbode(ctx, n, p)
		cur = n
	}
	series, _ = m.ListGeneralLawPeriods(ctx, 3)
	if len(series) != 3 {
		t.Fatalf("limited series len = %d, want 3", len(series))
	}
	if series[0].Period >= series[2].Period {
		t.Errorf("series not ascending: %d..%d", series[0].Period, series[2].Period)
	}
}
```

- [ ] **Step 2: Run to confirm it fails**

Run: `cd services/simulation-engine && go test ./internal/store/ -run TestMemoryAbodeRoundTrip 2>&1 | head`
Expected: FAIL — `m.GetAbodeState undefined`.

- [ ] **Step 3: Add fields to the Memory struct**

In `internal/store/memory.go`, find the `Memory` struct definition (the block of `mu sync.RWMutex` and the various maps). Add two fields to it:

```go
	abodeState        *simulation.AbodeState
	generalLawPeriods []simulation.GeneralLawPeriod
```

(`simulation` is already imported in `memory.go`.)

- [ ] **Step 4: Implement the three methods**

Append to `internal/store/memory.go`:

```go
// GetAbodeState implements AbodeStateStore. Defaults to the seeded initial abode
// when none has been persisted yet (parity with the MySQL migration seed).
func (m *Memory) GetAbodeState(_ context.Context) (simulation.AbodeState, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.abodeState == nil {
		return simulation.NewAbodeState(), nil
	}
	return *m.abodeState, nil
}

// AdvanceAbode implements AbodeStateStore: replace the aggregate and append one
// period to the immiseration series.
func (m *Memory) AdvanceAbode(_ context.Context, next simulation.AbodeState, period simulation.GeneralLawPeriod) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	state := next
	m.abodeState = &state
	m.generalLawPeriods = append(m.generalLawPeriods, period)
	return nil
}

// ListGeneralLawPeriods implements AbodeStateStore: the most recent periods in
// ascending order. A non-positive limit returns the whole series.
func (m *Memory) ListGeneralLawPeriods(_ context.Context, limit int) ([]simulation.GeneralLawPeriod, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]simulation.GeneralLawPeriod, len(m.generalLawPeriods))
	copy(out, m.generalLawPeriods)
	if limit > 0 && len(out) > limit {
		out = out[len(out)-limit:]
	}
	return out, nil
}
```

- [ ] **Step 5: Run the test to verify it passes**

Run: `cd services/simulation-engine && go test ./internal/store/ -run TestMemoryAbodeRoundTrip -v 2>&1 | tail`
Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add services/simulation-engine/internal/store/memory.go \
        services/simulation-engine/internal/store/abode_test.go
git commit --no-gpg-sign -m "feat(atlas): Memory AbodeStateStore (default seed + ascending series)"
```

---

## Task B3: MySQL implementation

**Files:**
- Modify: `services/simulation-engine/internal/store/mysql.go`
- Modify: `services/simulation-engine/internal/store/store.go` (compile-time assertions)

- [ ] **Step 1: Implement the three methods**

Append to `internal/store/mysql.go` (receiver `m *MySQL`, handle `m.db`; `database/sql`, `context`, `time` are imported; the ID constructor `circulation.NewStageDistributionID()` is available for the series row id):

```go
// abodeStateID is the fixed primary key of the singleton abode_state row,
// seeded by migration 00066.
const abodeStateID = "5eed000000000000abode1"

// GetAbodeState implements AbodeStateStore. The row is seeded by the migration,
// so a missing row falls back to the in-code default rather than erroring.
func (m *MySQL) GetAbodeState(ctx context.Context) (simulation.AbodeState, error) {
	const q = `
SELECT period, constant_pence, variable_pence, base_wage_pence, worker_supply,
       surplus_rate_base_bp, accumulation_rate_bp, marginal_composition_bp,
       displacement_rate_bp, productivity_growth_bp, population_growth_bp
FROM abode_state WHERE id = ?`
	var a simulation.AbodeState
	var c, v int64
	err := m.db.QueryRowContext(ctx, q, abodeStateID).Scan(
		&a.Period, &c, &v, &a.BaseWagePence, &a.WorkerSupply,
		&a.SurplusRateBaseBP, &a.AccumulationRateBP, &a.MarginalCompositionBP,
		&a.DisplacementRateBP, &a.ProductivityGrowthBP, &a.PopulationGrowthBP)
	if err == sql.ErrNoRows {
		return simulation.NewAbodeState(), nil
	}
	if err != nil {
		return simulation.AbodeState{}, err
	}
	a.ConstantPence = circulation.Pence(c)
	a.VariablePence = circulation.Pence(v)
	return a, nil
}

// AdvanceAbode implements AbodeStateStore: update the singleton aggregate and
// append one immiseration-series row in a single transaction.
func (m *MySQL) AdvanceAbode(ctx context.Context, next simulation.AbodeState, period simulation.GeneralLawPeriod) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if _, err = tx.ExecContext(ctx, `
UPDATE abode_state SET period = ?, constant_pence = ?, variable_pence = ?,
       base_wage_pence = ?, worker_supply = ?, surplus_rate_base_bp = ?,
       accumulation_rate_bp = ?, marginal_composition_bp = ?,
       displacement_rate_bp = ?, productivity_growth_bp = ?, population_growth_bp = ?
WHERE id = ?`,
		next.Period, int64(next.ConstantPence), int64(next.VariablePence),
		next.BaseWagePence, next.WorkerSupply, next.SurplusRateBaseBP,
		next.AccumulationRateBP, next.MarginalCompositionBP, next.DisplacementRateBP,
		next.ProductivityGrowthBP, next.PopulationGrowthBP, abodeStateID); err != nil {
		return err
	}

	if _, err = tx.ExecContext(ctx, `
INSERT INTO general_law_periods
    (id, period, wage_pence, rate_of_exploitation_bp, reserve_army_count,
     organic_composition_bp, employed_count, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		string(circulation.NewStageDistributionID()), period.Period, period.WagePence,
		period.RateOfExploitationBP, period.ReserveArmyCount, period.OrganicCompositionBP,
		period.EmployedCount, time.Now().UTC()); err != nil {
		return err
	}
	return tx.Commit()
}

// ListGeneralLawPeriods implements AbodeStateStore: the most recent periods in
// ascending period order. A non-positive limit returns the whole series.
func (m *MySQL) ListGeneralLawPeriods(ctx context.Context, limit int) ([]simulation.GeneralLawPeriod, error) {
	q := `
SELECT period, wage_pence, rate_of_exploitation_bp, reserve_army_count,
       organic_composition_bp, employed_count
FROM general_law_periods ORDER BY period DESC`
	if limit > 0 {
		q += " LIMIT ?"
	}
	var rows *sql.Rows
	var err error
	if limit > 0 {
		rows, err = m.db.QueryContext(ctx, q, limit)
	} else {
		rows, err = m.db.QueryContext(ctx, q)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var desc []simulation.GeneralLawPeriod
	for rows.Next() {
		var p simulation.GeneralLawPeriod
		if err := rows.Scan(&p.Period, &p.WagePence, &p.RateOfExploitationBP,
			&p.ReserveArmyCount, &p.OrganicCompositionBP, &p.EmployedCount); err != nil {
			return nil, err
		}
		desc = append(desc, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Reverse to ascending (oldest first) for the sparkline.
	out := make([]simulation.GeneralLawPeriod, len(desc))
	for i, p := range desc {
		out[len(desc)-1-i] = p
	}
	return out, nil
}
```

(`mysql.go` already imports `circulation` and `time` — `FieldSnapshot`/`AccumulateCapital` use them.)

- [ ] **Step 2: Add compile-time interface assertions**

In `internal/store/store.go`, search for existing `var _ ` assertions (e.g. `var _ IndustrialCapitalStore = (*Memory)(nil)`). If a block exists, add these two lines beside it; otherwise add them immediately below the `AbodeStateStore` interface definition:

```go
var _ AbodeStateStore = (*Memory)(nil)
var _ AbodeStateStore = (*MySQL)(nil)
```

- [ ] **Step 3: Build the store package**

Run: `cd services/simulation-engine && go build ./internal/store/ 2>&1 | head`
Expected: PASS (both Memory and MySQL satisfy `AbodeStateStore`).

- [ ] **Step 4: Vet + test the store package**

Run: `cd services/simulation-engine && go test ./internal/store/ -run 'Abode|Accumulate|FieldSnapshot' 2>&1 | tail`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add services/simulation-engine/internal/store/mysql.go \
        services/simulation-engine/internal/store/store.go
git commit --no-gpg-sign -m "feat(atlas): MySQL AbodeStateStore (tx advance + ascending series)"
```

---

# GROUP C — Backend ticker + snapshot v2

## Task C1: The `GeneralLawTicker` + registration

**Files:**
- Create: `services/simulation-engine/internal/engine/general_law_ticker.go`
- Test: `services/simulation-engine/internal/engine/general_law_ticker_test.go`
- Modify: `services/simulation-engine/cmd/simulation-engine/main.go`

- [ ] **Step 1: Write the failing ticker test**

Create `services/simulation-engine/internal/engine/general_law_ticker_test.go`:

```go
package engine_test

import (
	"context"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func TestGeneralLawTickerAdvancesTheAbode(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	tk := engine.NewGeneralLawTicker(m)
	if tk.Name() != "general-law" {
		t.Fatalf("name = %q", tk.Name())
	}

	before, _ := m.GetAbodeState(ctx)
	advanced, err := tk.Tick(ctx)
	if err != nil {
		t.Fatalf("tick: %v", err)
	}
	if advanced != 1 {
		t.Fatalf("advanced = %d, want 1", advanced)
	}

	after, _ := m.GetAbodeState(ctx)
	if after.Period != before.Period+1 {
		t.Errorf("period %d -> %d, want +1", before.Period, after.Period)
	}

	// Three more passes; the immiseration series grows and the reserve army is
	// larger at the end than at the start (the law's direction).
	for i := 0; i < 3; i++ {
		if _, err := tk.Tick(ctx); err != nil {
			t.Fatalf("tick %d: %v", i, err)
		}
	}
	series, _ := m.ListGeneralLawPeriods(ctx, 0)
	if len(series) != 4 {
		t.Fatalf("series len = %d, want 4", len(series))
	}
	if series[len(series)-1].ReserveArmyCount <= series[0].ReserveArmyCount {
		t.Errorf("reserve army did not grow: %d -> %d",
			series[0].ReserveArmyCount, series[len(series)-1].ReserveArmyCount)
	}
}
```

- [ ] **Step 2: Run to confirm it fails**

Run: `cd services/simulation-engine && go test ./internal/engine/ -run TestGeneralLawTicker 2>&1 | head`
Expected: FAIL — `engine.NewGeneralLawTicker undefined`.

- [ ] **Step 3: Write the ticker**

Create `internal/engine/general_law_ticker.go`:

```go
package engine

import (
	"context"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// AbodeAdvancer is the store subset the GeneralLawTicker needs. It is satisfied
// structurally by *store.Memory and *store.MySQL.
type AbodeAdvancer interface {
	GetAbodeState(ctx context.Context) (simulation.AbodeState, error)
	AdvanceAbode(ctx context.Context, next simulation.AbodeState, period simulation.GeneralLawPeriod) error
}

// GeneralLawTicker advances the hidden abode one period of the General Law of
// Capitalist Accumulation (Vol. I Ch. 25) every scheduler pass: rising
// composition and machinery repel labour into the reserve army, which presses
// the wage down and raises s/v, whose surplus re-accumulates. It implements
// Ticker. Registered after the AccumulationTicker so the surface (the field of
// orbits) and the abode (the class relation) advance in the same pass.
type GeneralLawTicker struct {
	store AbodeAdvancer
}

// NewGeneralLawTicker constructs the ticker.
func NewGeneralLawTicker(s AbodeAdvancer) *GeneralLawTicker {
	return &GeneralLawTicker{store: s}
}

// Name identifies this ticker in status and the audit log.
func (t *GeneralLawTicker) Name() string { return "general-law" }

// Tick advances the abode by one period and records it.
func (t *GeneralLawTicker) Tick(ctx context.Context) (int, error) {
	state, err := t.store.GetAbodeState(ctx)
	if err != nil {
		return 0, err
	}
	next, period := simulation.AdvanceGeneralLaw(state)
	if err := t.store.AdvanceAbode(ctx, next, period); err != nil {
		return 0, err
	}
	return 1, nil
}
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `cd services/simulation-engine && go test ./internal/engine/ -run TestGeneralLawTicker -v 2>&1 | tail`
Expected: PASS.

- [ ] **Step 5: Register the ticker in `main.go`**

In `cmd/simulation-engine/main.go`, the scheduler is built around line 95. Add the general-law ticker as the last element of the ticker slice (keep the existing tickers exactly as they are):

```go
	scheduler := engine.NewScheduler(tickInterval(), []engine.Ticker{
		engine.NewFactoryTicker(st),
		engine.NewReproductionTicker(st),
		engine.NewPiecePriceTicker(engine.NewFactoryProductivitySource(st), repricer),
		engine.NewAccumulationTicker(st, accumulationRateBP()),
		engine.NewGeneralLawTicker(st),
	}, st, logger)
```

- [ ] **Step 6: Build**

Run: `cd services/simulation-engine && go build ./... 2>&1 | head`
Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add services/simulation-engine/internal/engine/general_law_ticker.go \
        services/simulation-engine/internal/engine/general_law_ticker_test.go \
        services/simulation-engine/cmd/simulation-engine/main.go
git commit --no-gpg-sign -m "feat(atlas): GeneralLawTicker advances the abode each scheduler pass"
```

---

## Task C2: Snapshot v2 — the `abode` block

**Files:**
- Modify: `services/simulation-engine/internal/transport/httpapi/handler.go`
- Modify: `services/simulation-engine/internal/transport/httpapi/observatory_handler.go`
- Modify: `services/simulation-engine/internal/transport/httpapi/observatory_handler_test.go`
- Modify: `services/simulation-engine/cmd/simulation-engine/main.go`

- [ ] **Step 1: Add the `AbodeStates` dependency to the Handler**

In `internal/transport/httpapi/handler.go`, find the `Handler` struct (the block of `Xxx store.XxxStore` fields — `IndustrialCapitals` is one of them) and add, next to `IndustrialCapitals`:

```go
	AbodeStates store.AbodeStateStore
```

(`store` is already imported here.)

- [ ] **Step 2: Extend the snapshot response + handler**

In `internal/transport/httpapi/observatory_handler.go`:

(a) Add `Abode abodeDTO` to the response struct:

```go
type observatorySnapshotResponse struct {
	Tick       int64              `json:"tick"`
	Running    bool               `json:"running"`
	IntervalMS int64              `json:"interval_ms"`
	Capitals   []fieldCapitalDTO  `json:"capitals"`
	Aggregate  aggregateVitalsDTO `json:"aggregate"`
	Abode      abodeDTO           `json:"abode"`
}
```

(b) Add the two DTOs (next to `aggregateVitalsDTO`):

```go
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
	LawSeries              []generalLawPeriodDTO `json:"law_series"`
}

type generalLawPeriodDTO struct {
	Period               int64 `json:"period"`
	WagePence            int64 `json:"wage_pence"`
	RateOfExploitationBP int64 `json:"rate_of_exploitation_bp"`
	ReserveArmyCount     int64 `json:"reserve_army_count"`
	OrganicCompositionBP int64 `json:"organic_composition_bp"`
}
```

(c) In `GetObservatorySnapshot`, after the `resp.Aggregate = ...` assignment and before the scheduler block, build the abode. The `law_series` is always a non-nil slice:

```go
	resp.Abode = abodeDTO{LawSeries: []generalLawPeriodDTO{}}
	if h.AbodeStates != nil {
		state, err := h.AbodeStates.GetAbodeState(r.Context())
		if err != nil {
			h.writeServerError(w, err)
			return
		}
		ar := state.Readout()
		series, err := h.AbodeStates.ListGeneralLawPeriods(r.Context(), 60)
		if err != nil {
			h.writeServerError(w, err)
			return
		}
		law := make([]generalLawPeriodDTO, len(series))
		for i, p := range series {
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
			LawSeries:              law,
		}
	}
```

- [ ] **Step 3: Wire the dependency in `main.go`**

In `cmd/simulation-engine/main.go`, in the `httpapi.Handler{...}` literal (around line 107–140, where `IndustrialCapitals: st,` is), add:

```go
		AbodeStates:           st,
```

- [ ] **Step 4: Extend the handler test**

In `internal/transport/httpapi/observatory_handler_test.go`, locate `TestGetObservatorySnapshot`. Find where the `Handler` under test is constructed (it sets `IndustrialCapitals:` to an in-memory store `m`). Add `AbodeStates: m,` to that same struct literal. Then, after the existing assertions, add:

```go
	// Snapshot v2 carries the hidden-abode block.
	ab := resp.Abode
	if ab.TotalVariablePence <= 0 {
		t.Errorf("abode Σv = %d, want > 0", ab.TotalVariablePence)
	}
	if ab.LawSeries == nil {
		t.Error("abode.law_series must be non-nil (never null)")
	}
	// s/v == round(10000 * Σs / Σv) — the rate of exploitation.
	wantSV := (ab.TotalSurplusPence*10000 + ab.TotalVariablePence/2) / ab.TotalVariablePence
	if ab.RateOfExploitationBP != wantSV {
		t.Errorf("s/v = %d, want %d (round 10000*Σs/Σv)", ab.RateOfExploitationBP, wantSV)
	}
	if ab.NecessaryLabourMinutes+ab.SurplusLabourMinutes == 0 {
		t.Error("working day minutes not populated")
	}
```

If the test decodes the response into a local struct (not the handler's unexported `observatorySnapshotResponse`), add a matching `Abode` field to that local struct with the JSON tags `total_variable_pence`, `total_surplus_pence`, `rate_of_exploitation_bp`, `necessary_labour_minutes`, `surplus_labour_minutes`, `organic_composition_bp`, `reserve_army_count`, `reserve_army_pressure_bp`, `employed_count`, `wage_pence`, and `law_series []struct{ Period, WagePence, RateOfExploitationBP, ReserveArmyCount, OrganicCompositionBP int64 with matching snake_case tags }`.

> Note on the s/v identity: `Readout` computes `surplus = v·effRate/10000` with integer truncation, so `round(10000·Σs/Σv)` recovers `effRate` exactly only when `v·effRate` is divisible by 10000. With the seeded `v = 300000` and `effRate = 12000`, `surplus = 360000` and `round(10000·360000/300000) = 12000` ✓. Keep the seed values so the identity holds at period 0 (the snapshot reads the un-advanced state in this test).

- [ ] **Step 5: Run the handler test**

Run: `cd services/simulation-engine && go test ./internal/transport/httpapi/ -run TestGetObservatorySnapshot -v 2>&1 | tail`
Expected: PASS.

- [ ] **Step 6: Full backend check**

Run: `cd /mnt/c/Users/AaronHulse/IdeaProjects/capital-simulator && make vet test build 2>&1 | tail -20`
Expected: PASS — all packages, all six binaries.

- [ ] **Step 7: Commit**

```bash
git add services/simulation-engine/internal/transport/httpapi/handler.go \
        services/simulation-engine/internal/transport/httpapi/observatory_handler.go \
        services/simulation-engine/internal/transport/httpapi/observatory_handler_test.go \
        services/simulation-engine/cmd/simulation-engine/main.go
git commit --no-gpg-sign -m "feat(atlas): snapshot v2 — the abode block (s/v, reserve army, immiseration series)"
```

---

# GROUP D — Frontend (the descent)

## Task D1: Wire types + animation helpers

**Files:**
- Modify: `web/src/types.ts`
- Modify: `web/src/atlas/animation.ts`

- [ ] **Step 1: Add the abode types**

In `web/src/types.ts`, in the Atlas Observatory block (after `AggregateVitals`, before `ObservatorySnapshot`), add:

```ts
export interface GeneralLawTrendPoint {
  period: number;
  wage_pence: number;
  rate_of_exploitation_bp: number;
  reserve_army_count: number;
  organic_composition_bp: number;
}

export interface AbodeReadout {
  total_variable_pence: number;
  total_surplus_pence: number;
  rate_of_exploitation_bp: number;
  necessary_labour_minutes: number;
  surplus_labour_minutes: number;
  organic_composition_bp: number;
  reserve_army_count: number;
  reserve_army_pressure_bp: number;
  employed_count: number;
  wage_pence: number;
  law_series: GeneralLawTrendPoint[];
}
```

Then add `abode` to `ObservatorySnapshot`:

```ts
export interface ObservatorySnapshot {
  tick: number;
  running: boolean;
  interval_ms: number;
  capitals: FieldCapital[];
  aggregate: AggregateVitals;
  abode: AbodeReadout;
}
```

- [ ] **Step 2: Add `formatMinutes` + `sparklinePath` helpers**

Append to `web/src/atlas/animation.ts`:

```ts
/** Minutes → "Hh Mm" working-day label, e.g. 272 → "4h 32m". */
export function formatMinutes(min: number): string {
  const h = Math.floor(min / 60);
  const m = Math.round(min % 60);
  return `${h}h ${m}m`;
}

/**
 * Build an SVG polyline path for a sparkline of `values` fitted to a `w`×`h`
 * box (with 2px padding). Auto-scales to the value range; a flat series draws a
 * mid-line. Pure + deterministic.
 */
export function sparklinePath(values: number[], w: number, h: number): string {
  if (values.length === 0) return "";
  const pad = 2;
  const lo = Math.min(...values);
  const hi = Math.max(...values);
  const span = hi - lo || 1;
  const innerW = w - pad * 2;
  const innerH = h - pad * 2;
  const step = values.length > 1 ? innerW / (values.length - 1) : 0;
  return values
    .map((v, i) => {
      const x = pad + i * step;
      const y = pad + innerH * (1 - (v - lo) / span);
      return `${i === 0 ? "M" : "L"}${x.toFixed(1)} ${y.toFixed(1)}`;
    })
    .join(" ");
}
```

- [ ] **Step 3: Typecheck**

Run: `cd web && npm run lint 2>&1 | tail`
Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add web/src/types.ts web/src/atlas/animation.ts
git commit --no-gpg-sign -m "feat(atlas): abode wire types + minutes/sparkline helpers"
```

---

## Task D2: `GeneralLawTrend` sparkline

**Files:**
- Create: `web/src/atlas/GeneralLawTrend.tsx`

- [ ] **Step 1: Write the component**

Create `web/src/atlas/GeneralLawTrend.tsx`:

```tsx
import type { GeneralLawTrendPoint } from "../types";
import { sparklinePath, formatBP } from "./animation";

interface TrendProps {
  series: GeneralLawTrendPoint[];
}

const W = 260;
const H = 64;
const GOLD = "#c8a240";
const RED = "#c0392b";

/** The immiseration trend: rate of exploitation (s/v) rising and the paid wage
 *  falling across the recorded periods of the general law. */
export function GeneralLawTrend({ series }: TrendProps) {
  if (series.length < 2) {
    return <p className="abode-hint">The law has not yet run — start the engine.</p>;
  }
  const exploitation = series.map((p) => p.rate_of_exploitation_bp);
  const wages = series.map((p) => p.wage_pence);
  const latest = series[series.length - 1];
  return (
    <div className="abode-trend">
      <svg width={W} height={H} viewBox={`0 0 ${W} ${H}`} role="img"
        aria-label="immiseration trend: rising exploitation, falling wage">
        <path d={sparklinePath(exploitation, W, H)} fill="none" stroke={GOLD} strokeWidth={2} />
        <path d={sparklinePath(wages, W, H)} fill="none" stroke={RED} strokeWidth={2}
          strokeDasharray="3 2" />
      </svg>
      <div className="abode-trend-legend">
        <span style={{ color: GOLD }}>s/v {formatBP(latest.rate_of_exploitation_bp)}</span>
        <span style={{ color: RED }}>wage ↓</span>
        <span className="abode-hint">period {latest.period}</span>
      </div>
    </div>
  );
}
```

- [ ] **Step 2: Typecheck**

Run: `cd web && npm run lint 2>&1 | tail`
Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/GeneralLawTrend.tsx
git commit --no-gpg-sign -m "feat(atlas): GeneralLawTrend — the immiseration sparkline"
```

---

## Task D3: `Abode` view

**Files:**
- Create: `web/src/atlas/Abode.tsx`

- [ ] **Step 1: Write the component**

Create `web/src/atlas/Abode.tsx`:

```tsx
import type { AbodeReadout } from "../types";
import { formatBP, formatPence, formatMinutes } from "./animation";
import { GeneralLawTrend } from "./GeneralLawTrend";

interface AbodeProps {
  abode: AbodeReadout;
}

/** The hidden abode of production. We have left "the noisy sphere, where
 *  everything takes place on the surface and in view of all men" for the place
 *  where the class relation — surplus as unpaid labour, the reserve army — is
 *  laid bare. The working day divides into necessary (paid, v) and surplus
 *  (unpaid, s) labour; their ratio is the rate of exploitation s/v. */
export function Abode({ abode }: AbodeProps) {
  const day = abode.necessary_labour_minutes + abode.surplus_labour_minutes || 1;
  const necPct = (abode.necessary_labour_minutes / day) * 100;
  const surPct = 100 - necPct;
  return (
    <div className="abode" data-testid="abode">
      <div className="abode-head">
        <h2>The hidden abode of production</h2>
        <p className="abode-hint">No admittance except on business.</p>
      </div>

      <section className="abode-card">
        <div className="abode-card-k">The social working day · s/v {formatBP(abode.rate_of_exploitation_bp)}</div>
        <div className="workingday" data-testid="workingday">
          <div className="wd-nec" style={{ width: `${necPct}%` }}>
            <span>necessary · {formatMinutes(abode.necessary_labour_minutes)}</span>
          </div>
          <div className="wd-sur" style={{ width: `${surPct}%` }}>
            <span>surplus · {formatMinutes(abode.surplus_labour_minutes)}</span>
          </div>
        </div>
      </section>

      <div className="abode-grid">
        <section className="abode-card">
          <div className="abode-card-k">Living labour · Σv (paid)</div>
          <div className="abode-card-v">{formatPence(abode.total_variable_pence)}</div>
          <div className="abode-hint">{abode.employed_count} employed</div>
        </section>
        <section className="abode-card gold">
          <div className="abode-card-k">Surplus extracted · Σs (unpaid)</div>
          <div className="abode-card-v">{formatPence(abode.total_surplus_pence)}</div>
          <div className="abode-hint">rises to the surface as capital</div>
        </section>
        <section className="abode-card">
          <div className="abode-card-k">Industrial reserve army</div>
          <div className="abode-card-v">{abode.reserve_army_count}</div>
          <div className="abode-hint">wage pressure {formatBP(abode.reserve_army_pressure_bp)} · paid wage {formatPence(abode.wage_pence)}</div>
        </section>
        <section className="abode-card">
          <div className="abode-card-k">Organic composition · c/v</div>
          <div className="abode-card-v">{formatBP(abode.organic_composition_bp)}</div>
          <div className="abode-hint">dead labour dominating living</div>
        </section>
      </div>

      <section className="abode-card">
        <div className="abode-card-k">The general law in motion</div>
        <GeneralLawTrend series={abode.law_series} />
      </section>
    </div>
  );
}
```

- [ ] **Step 2: Typecheck**

Run: `cd web && npm run lint 2>&1 | tail`
Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/Abode.tsx
git commit --no-gpg-sign -m "feat(atlas): Abode view — working day, living labour, reserve army, surplus"
```

---

## Task D4: Threshold toggle in `Atlas.tsx` + styles

**Files:**
- Modify: `web/src/atlas/Atlas.tsx`
- Modify: `web/src/atlas/atlas.css`

- [ ] **Step 1: Add the descend/ascend threshold to `Atlas.tsx`**

Replace the contents of `web/src/atlas/Atlas.tsx` with:

```tsx
import { useState } from "react";
import "./atlas.css";
import { useSnapshot } from "./useSnapshot";
import { CircuitField } from "./CircuitField";
import { Abode } from "./Abode";
import { VitalSigns } from "./VitalSigns";
import { TickHeartbeat } from "./TickHeartbeat";

/** The Observatory: the whole circuit of capital, in motion — and beneath it,
 *  the hidden abode of production. */
export default function Atlas() {
  const { snapshot, stale } = useSnapshot();
  const [speed, setSpeed] = useState(1);
  const [descended, setDescended] = useState(false);

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
          {snapshot && (
            <button
              className={"abode-threshold" + (descended ? " open" : "")}
              data-testid="threshold"
              onClick={() => setDescended((d) => !d)}
            >
              {descended ? "↑ Ascend to the surface" : "↓ Descend into production"}
            </button>
          )}
        </aside>

        {snapshot ? (
          descended ? (
            <div className="atlas-field-wrap abode-wrap">
              <Abode abode={snapshot.abode} />
            </div>
          ) : (
            <CircuitField snapshot={snapshot} speed={speed} />
          )
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

- [ ] **Step 2: Add the abode styles**

Append to `web/src/atlas/atlas.css` (note: the hover colour is the literal hex `#3a4256` — do not leave any placeholder text in it):

```css
/* --- The hidden abode (Slice 2) ------------------------------------------ */
.abode-threshold {
  margin-top: 18px;
  width: 100%;
  padding: 10px 12px;
  background: #161922;
  color: #f4ecd8;
  border: 1px solid #2a2f3e;
  border-radius: 6px;
  cursor: pointer;
  font: inherit;
  letter-spacing: 0.02em;
  transition: background 0.2s ease, border-color 0.2s ease;
}
.abode-threshold:hover { background: #1d2230; border-color: #3a4256; }
.abode-threshold.open { background: #2a1a16; border-color: #6e3b2f; }

.abode-wrap {
  overflow-y: auto;
  animation: abode-rise 0.5s ease-out;
}
@keyframes abode-rise {
  from { opacity: 0; transform: translateY(18px); }
  to { opacity: 1; transform: translateY(0); }
}

.abode { padding: 28px 32px; max-width: 820px; margin: 0 auto; }
.abode-head h2 { margin: 0 0 2px; color: #f4ecd8; font-size: 22px; }
.abode-hint { color: var(--ink-muted, #8a90a0); font-size: 12px; margin: 2px 0 0; }

.abode-card {
  background: #141821;
  border: 1px solid #242a38;
  border-radius: 8px;
  padding: 14px 16px;
  margin-top: 16px;
}
.abode-card.gold { border-color: #6b571f; background: #1c1810; }
.abode-card-k { color: var(--ink-muted, #8a90a0); font-size: 12px; text-transform: uppercase; letter-spacing: 0.04em; }
.abode-card-v { color: #f4ecd8; font-size: 24px; margin-top: 4px; }

.abode-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; margin-top: 16px; }
.abode-grid .abode-card { margin-top: 0; }

.workingday {
  display: flex;
  height: 38px;
  border-radius: 6px;
  overflow: hidden;
  margin-top: 8px;
  border: 1px solid #242a38;
}
.workingday .wd-nec,
.workingday .wd-sur {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  color: #0d0f14;
  white-space: nowrap;
  transition: width 0.8s ease;
}
.workingday .wd-nec { background: #4a5a8a; color: #e6ecff; }
.workingday .wd-sur { background: #c8a240; }

.abode-trend { display: flex; flex-direction: column; gap: 6px; margin-top: 6px; }
.abode-trend svg { background: #0d0f14; border-radius: 6px; }
.abode-trend-legend { display: flex; gap: 14px; font-size: 12px; }
```

- [ ] **Step 3: Typecheck + build**

Run: `cd web && npm run lint && npm run build 2>&1 | tail -20`
Expected: success.

- [ ] **Step 4: Commit**

```bash
git add web/src/atlas/Atlas.tsx web/src/atlas/atlas.css
git commit --no-gpg-sign -m "feat(atlas): threshold — descend from the surface into the hidden abode"
```

---

# GROUP E — End-to-end acceptance

## Task E1: Boot the stack + drive the descent (Playwright MCP)

**Files:** none (per CLAUDE.md "always boot stack for panel E2E").

- [ ] **Step 1: Boot fresh (migration 00066 must apply)**

Run:
```bash
docker compose down -v
docker compose up --build -d
```
Wait ~3 min for MySQL + migrations. Then start the engine and confirm the abode advances:
```bash
curl -s -X POST http://localhost:8080/v1/engine/start >/dev/null
curl -s http://localhost:8080/v1/observatory/snapshot | python3 -m json.tool | grep -A20 '"abode"'
# wait ~15s, read again — period climbs, reserve army grows, s/v rises
sleep 15
curl -s http://localhost:8080/v1/observatory/snapshot | python3 -m json.tool | grep -A20 '"abode"'
```
Expected: the `abode` block is present; `law_series` gains entries; across the two reads `reserve_army_count` and `rate_of_exploitation_bp` rise and `wage_pence` falls.

> If MySQL hits a goose "out-of-order migration" boot error from a stale `mysql_data` volume (see the dev-MySQL-drift feedback memory), `docker compose down -v` already dropped it; if you skipped `-v`, drop/recreate the `simulation` schema and restart the service.

- [ ] **Step 2: Drive the page with Playwright MCP**

1. `browser_navigate` → `http://localhost:5173/`.
2. If transport shows ▶, click it so the General-Law ticker runs.
3. `browser_snapshot` → confirm the field of orbits and the "↓ Descend into production" control are present.
4. `browser_click` the `[data-testid="threshold"]` control → assert `[data-testid="abode"]` appears with the working-day bar (`[data-testid="workingday"]`), a non-zero s/v, and the reserve-army card.
5. `browser_evaluate` → read the surplus arc width of `.workingday .wd-sur`, wait ~20s, re-read → assert the surplus share **grew** (s/v widening as the law runs).
6. `browser_take_screenshot` → `atlas-abode.png` for the PR.
7. `browser_click` the threshold again → assert it ascends back to the field.

Expected: the descent reveals a live abode; the working day widens toward surplus; the immiseration sparkline climbs.

- [ ] **Step 3: Tear down + final check**

```bash
docker compose down
cd /mnt/c/Users/AaronHulse/IdeaProjects/capital-simulator && make vet test build
cd web && npm run lint && npm run build
```
Expected: all PASS.

---

# Done criteria (Slice 2)

- A scheduler pass advances the hidden abode one period of the General Law: the reserve army grows, the organic composition rises, the rate of exploitation s/v rises, and the paid wage falls — verified by `AdvanceGeneralLaw` tests and live over the booted stack.
- `GET /v1/observatory/snapshot` carries an `abode` block (Σv, Σs, s/v, working-day minutes, c/v, reserve army + pressure, employed, paid wage, and the immiseration `law_series`), with `law_series` never null and `s/v == round(10000·Σs/Σv)`.
- The Atlas page descends from the field of orbits into the abode (threshold transition framed by "No admittance except on business"), showing the social working day (necessary|surplus), living labour, reserve army + wage pressure, surplus extraction, and the general-law trend.
- `make vet test build` + `npm run lint && npm run build` pass; Playwright shows the descent and the widening working day.

# Deferred (Slice 3)
- The levers — working day · wage level · accumulation rate (`α`) — POSTed to perturb the live `AbodeState` and watched responding (the design's §10 item 3).
- Per-capital `variable_pence`/`constant_pence` in the `capitals` array (orbit-level c/v colouring); not required by any Slice-2 abode component.
- Per-capital descent (into one factory) and server-push streaming (design §12 YAGNI).
