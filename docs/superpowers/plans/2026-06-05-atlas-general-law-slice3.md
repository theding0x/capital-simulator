# Atlas General Law — Slice 3 (The Levers) Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Put the General Law under the operator's hand — three levers (the **working day** = the rate of surplus-value, the **wage** = the value of labour-power, and the **accumulation rate** α) that perturb the live `AbodeState`, so you can drag a control and watch the law respond over the next scheduler passes (reserve army, s/v, composition, wage all shift).

**Architecture:** A new `POST /v1/observatory/levers` applies a partial `LeverUpdate` to the persisted singleton `AbodeState` (the three existing parameter fields `SurplusRateBaseBP`, `BaseWagePence`, `AccumulationRateBP`), clamped to safe ranges — the verified `AdvanceGeneralLaw` step is untouched and simply reads the new parameters each pass. The snapshot's `abode` block carries the three base lever values so the UI sliders reflect live state. A `Levers` panel in the hidden abode POSTs on change.

**Tech Stack:** Go 1.25 (pure `ApplyLevers` + `database/sql` tx read-modify-write), React 18 + TS range inputs. No migration (reuses the Slice-2 `abode_state` row), no new dependencies.

**Spec:** `docs/superpowers/specs/2026-06-04-atlas-general-law-design.md` (Slice 3 = §10 item 3; levers named in §2 and §4). Branch: `feature/atlas-observatory` (continues Slice 2).

**Key modeling decisions:**
- **The three levers map to the three existing `AbodeState` parameter fields** — no change to the verified `AdvanceGeneralLaw` mechanism. The working-day lever sets `SurplusRateBaseBP` (the working day's necessary↔surplus division *is* the rate of surplus-value, Ch. 9–10); the wage lever sets `BaseWagePence` (the value of labour-power, which sets employment `v/wage` and thus the reserve army); the α lever sets `AccumulationRateBP` (the share of surplus re-accumulated).
- **`AdvanceAbode` already carries the parameters forward unchanged**, so a lever set between ticks is read by the very next pass. The (microsecond) read-modify-write race inside one 5-second pass is acceptable for a live viz; `SetAbodeLevers` still uses a `FOR UPDATE` tx on MySQL (per the repo's balance-mutation convention).
- **The abode block gains `surplus_rate_base_bp`, `base_wage_pence`, `accumulation_rate_bp`** (the *base* lever values, distinct from the derived `rate_of_exploitation_bp`/`wage_pence`) so the sliders initialise to live positions.

---

## File Structure

**Backend (`services/simulation-engine/`):**
- Modify `internal/simulation/abode.go` — add `LeverUpdate`, `IsEmpty`, `AbodeState.ApplyLevers`.
- Modify `internal/simulation/abode_test.go` — `ApplyLevers` tests.
- Modify `internal/store/store.go` — add `SetAbodeLevers` to `AbodeStateStore`.
- Modify `internal/store/memory.go` — `SetAbodeLevers`.
- Modify `internal/store/mysql.go` — `SetAbodeLevers` (tx + `FOR UPDATE`).
- Modify `internal/store/abode_test.go` — memory lever test.
- Create `internal/transport/httpapi/observatory_levers_handler.go` — `POST /v1/observatory/levers`.
- Modify `internal/transport/httpapi/observatory_handler.go` — three base lever fields on `abodeDTO`.
- Modify `internal/transport/httpapi/observatory_handler_test.go` — levers handler test + snapshot lever-field assertion.
- Modify `internal/transport/httpapi/routes.go` — register the POST route.
- Modify `services/api-gateway/cmd/api-gateway/main.go` — proxy `/v1/observatory/levers`.

**Frontend (`web/src/`):**
- Modify `types.ts` — `LeverUpdate`, `LeverState`, three fields on `AbodeReadout`.
- Modify `api.ts` — `setObservatoryLevers`.
- Create `atlas/Levers.tsx` — the three lever controls.
- Modify `atlas/Abode.tsx` — render `<Levers>`.
- Modify `atlas/atlas.css` — lever styles.

---

# GROUP A — Backend domain (apply levers)

## Task A1: `LeverUpdate` + `ApplyLevers`

**Files:**
- Modify: `services/simulation-engine/internal/simulation/abode.go`
- Test: `services/simulation-engine/internal/simulation/abode_test.go`

- [ ] **Step 1: Write the failing test**

Append to `services/simulation-engine/internal/simulation/abode_test.go`:

```go
func TestApplyLevers(t *testing.T) {
	t.Parallel()
	s := NewAbodeState() // surplus_rate_base=10000, base_wage=2500, accum=5000

	sr, wage, ac := int64(20000), int64(4000), int64(8000)
	got := s.ApplyLevers(LeverUpdate{SurplusRateBaseBP: &sr, BaseWagePence: &wage, AccumulationRateBP: &ac})
	if got.SurplusRateBaseBP != 20000 || got.BaseWagePence != 4000 || got.AccumulationRateBP != 8000 {
		t.Errorf("levers not applied: %+v", got)
	}

	// A partial update leaves the other parameters untouched.
	only := s.ApplyLevers(LeverUpdate{AccumulationRateBP: &ac})
	if only.AccumulationRateBP != 8000 || only.BaseWagePence != s.BaseWagePence || only.SurplusRateBaseBP != s.SurplusRateBaseBP {
		t.Errorf("partial update bled into other fields: %+v", only)
	}

	// Clamps: α to [0,10000]; wage floored at 1; base surplus rate to [0,100000].
	big, zero, huge := int64(99999), int64(0), int64(500000)
	cl := s.ApplyLevers(LeverUpdate{AccumulationRateBP: &big, BaseWagePence: &zero, SurplusRateBaseBP: &huge})
	if cl.AccumulationRateBP != 10000 {
		t.Errorf("alpha not clamped: %d", cl.AccumulationRateBP)
	}
	if cl.BaseWagePence != 1 {
		t.Errorf("wage not floored: %d", cl.BaseWagePence)
	}
	if cl.SurplusRateBaseBP != 100000 {
		t.Errorf("surplus rate not clamped: %d", cl.SurplusRateBaseBP)
	}

	// An empty update changes nothing and reports empty.
	if !(LeverUpdate{}).IsEmpty() {
		t.Error("empty LeverUpdate should report IsEmpty")
	}
	if got := s.ApplyLevers(LeverUpdate{}); got != s {
		t.Error("empty update should be a no-op")
	}
}
```

- [ ] **Step 2: Run to confirm it fails**

Run: `cd services/simulation-engine && go test ./internal/simulation/ -run TestApplyLevers 2>&1 | head`
Expected: FAIL — `undefined: LeverUpdate` / `ApplyLevers`.

- [ ] **Step 3: Add `LeverUpdate` + `ApplyLevers` to `abode.go`**

Append to `services/simulation-engine/internal/simulation/abode.go`:

```go
// LeverUpdate is a partial perturbation of the live abode's law parameters
// (Slice 3 — the levers). Non-nil fields are applied; nil fields are left
// unchanged. The three levers are the working day (the rate of surplus-value),
// the wage (the value of labour-power), and the accumulation rate α.
type LeverUpdate struct {
	SurplusRateBaseBP  *int64 `json:"surplus_rate_base_bp,omitempty"`  // working day: necessary↔surplus
	BaseWagePence      *int64 `json:"base_wage_pence,omitempty"`       // value of labour-power
	AccumulationRateBP *int64 `json:"accumulation_rate_bp,omitempty"`  // α
}

// IsEmpty reports whether the update would change nothing.
func (u LeverUpdate) IsEmpty() bool {
	return u.SurplusRateBaseBP == nil && u.BaseWagePence == nil && u.AccumulationRateBP == nil
}

// ApplyLevers returns a copy of the state with the supplied levers applied, each
// clamped so the law cannot be driven to a degenerate state: the base rate of
// surplus-value to [0, 100000] (0–1000%), the wage to at least 1 pence (a
// positive value of labour-power, so employment v/wage stays finite), and α to
// [0, 10000] basis points.
func (a AbodeState) ApplyLevers(u LeverUpdate) AbodeState {
	next := a
	if u.SurplusRateBaseBP != nil {
		v := *u.SurplusRateBaseBP
		if v < 0 {
			v = 0
		}
		if v > 100000 {
			v = 100000
		}
		next.SurplusRateBaseBP = v
	}
	if u.BaseWagePence != nil {
		v := *u.BaseWagePence
		if v < 1 {
			v = 1
		}
		next.BaseWagePence = v
	}
	if u.AccumulationRateBP != nil {
		v := *u.AccumulationRateBP
		if v < 0 {
			v = 0
		}
		if v > 10000 {
			v = 10000
		}
		next.AccumulationRateBP = v
	}
	return next
}
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `cd services/simulation-engine && go test ./internal/simulation/ -run TestApplyLevers -v 2>&1 | tail`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add services/simulation-engine/internal/simulation/abode.go \
        services/simulation-engine/internal/simulation/abode_test.go
git commit --no-gpg-sign -m "feat(atlas): AbodeState.ApplyLevers — clamp + apply the three levers"
```

---

# GROUP B — Backend persistence

## Task B1: `SetAbodeLevers` (interface + memory + mysql)

**Files:**
- Modify: `services/simulation-engine/internal/store/store.go`
- Modify: `services/simulation-engine/internal/store/memory.go`
- Modify: `services/simulation-engine/internal/store/mysql.go`
- Test: `services/simulation-engine/internal/store/abode_test.go`

- [ ] **Step 1: Write the failing memory test**

Append to `services/simulation-engine/internal/store/abode_test.go`:

```go
func TestMemorySetAbodeLevers(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	m := store.NewMemory()

	zero := int64(0)
	got, err := m.SetAbodeLevers(ctx, simulation.LeverUpdate{AccumulationRateBP: &zero})
	if err != nil {
		t.Fatalf("set: %v", err)
	}
	if got.AccumulationRateBP != 0 {
		t.Errorf("alpha = %d, want 0", got.AccumulationRateBP)
	}
	// The other parameters are untouched (default base wage 2500).
	if got.BaseWagePence != 2500 {
		t.Errorf("base wage bled = %d, want 2500", got.BaseWagePence)
	}

	// Persisted: a later Get reflects the lever.
	st, _ := m.GetAbodeState(ctx)
	if st.AccumulationRateBP != 0 {
		t.Errorf("not persisted: %d", st.AccumulationRateBP)
	}

	// With α = 0 the law performs simple reproduction: no surplus is
	// capitalised, so total social capital (c+v) does not grow when advanced
	// (displacement only shifts value between c and v).
	next, _ := simulation.AdvanceGeneralLaw(st)
	if next.ConstantPence+next.VariablePence != st.ConstantPence+st.VariablePence {
		t.Errorf("alpha=0 should not grow total capital: %d -> %d",
			st.ConstantPence+st.VariablePence, next.ConstantPence+next.VariablePence)
	}
}
```

- [ ] **Step 2: Run to confirm it fails**

Run: `cd services/simulation-engine && go test ./internal/store/ -run TestMemorySetAbodeLevers 2>&1 | head`
Expected: FAIL — `m.SetAbodeLevers undefined`.

- [ ] **Step 3: Add the interface method**

In `internal/store/store.go`, inside the `AbodeStateStore` interface (after `ListGeneralLawPeriods`), add:

```go
	// SetAbodeLevers applies a partial lever update (Slice 3) to the live abode
	// and persists it, returning the updated state. Unset fields are unchanged.
	SetAbodeLevers(ctx context.Context, u simulation.LeverUpdate) (simulation.AbodeState, error)
```

- [ ] **Step 4: Implement in `memory.go`**

Append to `internal/store/memory.go`:

```go
// SetAbodeLevers implements AbodeStateStore (Slice 3 — the levers).
func (m *Memory) SetAbodeLevers(_ context.Context, u simulation.LeverUpdate) (simulation.AbodeState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cur := simulation.NewAbodeState()
	if m.abodeState != nil {
		cur = *m.abodeState
	}
	next := cur.ApplyLevers(u)
	m.abodeState = &next
	return next, nil
}
```

- [ ] **Step 5: Implement in `mysql.go`**

Append to `internal/store/mysql.go` (reuses `abodeStateID` defined in Slice 2; `simulation`, `circulation`, `database/sql` are imported):

```go
// SetAbodeLevers implements AbodeStateStore (Slice 3 — the levers). It reads the
// singleton row FOR UPDATE, applies the lever clamp, writes back only the three
// lever columns, and returns the updated state.
func (m *MySQL) SetAbodeLevers(ctx context.Context, u simulation.LeverUpdate) (simulation.AbodeState, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return simulation.AbodeState{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var a simulation.AbodeState
	var c, v int64
	err = tx.QueryRowContext(ctx, `
SELECT period, constant_pence, variable_pence, base_wage_pence, worker_supply,
       surplus_rate_base_bp, accumulation_rate_bp, marginal_composition_bp,
       displacement_rate_bp, productivity_growth_bp, population_growth_bp
FROM abode_state WHERE id = ? FOR UPDATE`, abodeStateID).Scan(
		&a.Period, &c, &v, &a.BaseWagePence, &a.WorkerSupply,
		&a.SurplusRateBaseBP, &a.AccumulationRateBP, &a.MarginalCompositionBP,
		&a.DisplacementRateBP, &a.ProductivityGrowthBP, &a.PopulationGrowthBP)
	if err == sql.ErrNoRows {
		a = simulation.NewAbodeState()
	} else if err != nil {
		return simulation.AbodeState{}, err
	} else {
		a.ConstantPence = simulation.Pence(c)
		a.VariablePence = simulation.Pence(v)
	}

	next := a.ApplyLevers(u)
	if _, err = tx.ExecContext(ctx, `
UPDATE abode_state SET surplus_rate_base_bp = ?, base_wage_pence = ?, accumulation_rate_bp = ?
WHERE id = ?`,
		next.SurplusRateBaseBP, next.BaseWagePence, next.AccumulationRateBP, abodeStateID); err != nil {
		return simulation.AbodeState{}, err
	}
	if err = tx.Commit(); err != nil {
		return simulation.AbodeState{}, err
	}
	return next, nil
}
```

- [ ] **Step 6: Run the memory test + build**

Run: `cd services/simulation-engine && go test ./internal/store/ -run 'Abode' 2>&1 | tail && go build ./internal/store/ 2>&1 | head`
Expected: PASS + clean build (the `var _ AbodeStateStore` assertions from Slice 2 now require the new method on both stores).

- [ ] **Step 7: Commit**

```bash
git add services/simulation-engine/internal/store/store.go \
        services/simulation-engine/internal/store/memory.go \
        services/simulation-engine/internal/store/mysql.go \
        services/simulation-engine/internal/store/abode_test.go
git commit --no-gpg-sign -m "feat(atlas): SetAbodeLevers store method (memory + mysql FOR UPDATE)"
```

---

# GROUP C — Backend transport

## Task C1: `POST /v1/observatory/levers` handler + route + gateway

**Files:**
- Create: `services/simulation-engine/internal/transport/httpapi/observatory_levers_handler.go`
- Modify: `services/simulation-engine/internal/transport/httpapi/observatory_handler.go`
- Modify: `services/simulation-engine/internal/transport/httpapi/routes.go`
- Modify: `services/api-gateway/cmd/api-gateway/main.go`
- Test: `services/simulation-engine/internal/transport/httpapi/observatory_handler_test.go`

- [ ] **Step 1: Add the three base lever fields to the snapshot's abode block**

In `internal/transport/httpapi/observatory_handler.go`, add three fields to `abodeDTO` (after `WagePence`, before `LawSeries`):

```go
	WagePence              int64                 `json:"wage_pence"`
	SurplusRateBaseBP      int64                 `json:"surplus_rate_base_bp"`
	BaseWagePence          int64                 `json:"base_wage_pence"`
	AccumulationRateBP     int64                 `json:"accumulation_rate_bp"`
	LawSeries              []generalLawPeriodDTO `json:"law_series"`
```

Then, in `GetObservatorySnapshot`, in the `resp.Abode = abodeDTO{...}` literal (inside the `if h.AbodeStates != nil` block), add the three values sourced from `state` (not the readout):

```go
			WagePence:              ar.WagePence,
			SurplusRateBaseBP:      state.SurplusRateBaseBP,
			BaseWagePence:          state.BaseWagePence,
			AccumulationRateBP:     state.AccumulationRateBP,
			LawSeries:              law,
```

- [ ] **Step 2: Write the levers handler**

Create `internal/transport/httpapi/observatory_levers_handler.go`:

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

// SetObservatoryLevers handles POST /v1/observatory/levers (Slice 3 — the
// levers). It applies a partial perturbation of the abode's law parameters —
// the working day (the rate of surplus-value), the wage (the value of
// labour-power), and the accumulation rate α — to the live state; the effects
// appear in subsequent snapshots as the General Law responds. Returns the
// applied (clamped) lever values.
func (h *Handler) SetObservatoryLevers(w http.ResponseWriter, r *http.Request) {
	if h.AbodeStates == nil {
		h.writeServerError(w, errors.New("abode state store not configured"))
		return
	}
	var body leversRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "malformed JSON body")
		return
	}
	state, err := h.AbodeStates.SetAbodeLevers(r.Context(), simulation.LeverUpdate{
		SurplusRateBaseBP:  body.SurplusRateBaseBP,
		BaseWagePence:      body.BaseWagePence,
		AccumulationRateBP: body.AccumulationRateBP,
	})
	if err != nil {
		h.writeServerError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, leversResponse{
		SurplusRateBaseBP:  state.SurplusRateBaseBP,
		BaseWagePence:      state.BaseWagePence,
		AccumulationRateBP: state.AccumulationRateBP,
	})
}
```

- [ ] **Step 3: Register the route**

In `internal/transport/httpapi/routes.go`, after the observatory snapshot line (`s.HandleFunc("GET /v1/observatory/snapshot", h.GetObservatorySnapshot)`), add:

```go
	s.HandleFunc("POST /v1/observatory/levers", h.SetObservatoryLevers)
```

- [ ] **Step 4: Proxy the new path through the gateway**

In `services/api-gateway/cmd/api-gateway/main.go`, after the line `srv.Handle("/v1/observatory/snapshot", simProxy)`, add:

```go
	srv.Handle("/v1/observatory/levers", simProxy)
```

(Without this the path 502s through `:8080` and CI never catches it — see the gateway per-path memory.)

- [ ] **Step 5: Write the handler test**

Append to `internal/transport/httpapi/observatory_handler_test.go`. Construct the `Handler` exactly as `TestGetObservatorySnapshot` does (an in-memory store `m` wired to `AbodeStates`); reuse that construction pattern:

```go
func TestSetObservatoryLevers(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	h := &Handler{AbodeStates: m} // mirror TestGetObservatorySnapshot's handler construction

	body := `{"accumulation_rate_bp": 0, "surplus_rate_base_bp": 30000}`
	req := httptest.NewRequest(http.MethodPost, "/v1/observatory/levers", strings.NewReader(body))
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
		t.Errorf("untouched wage = %d, want 2500 (default)", resp.BaseWagePence)
	}

	// Persisted on the live abode.
	st, _ := m.GetAbodeState(req.Context())
	if st.AccumulationRateBP != 0 || st.SurplusRateBaseBP != 30000 {
		t.Errorf("not persisted: %+v", st)
	}

	// Malformed body → 400.
	bad := httptest.NewRequest(http.MethodPost, "/v1/observatory/levers", strings.NewReader("{"))
	badRec := httptest.NewRecorder()
	h.SetObservatoryLevers(badRec, bad)
	if badRec.Code != http.StatusBadRequest {
		t.Errorf("malformed code = %d, want 400", badRec.Code)
	}
}
```

> If `TestGetObservatorySnapshot` constructs the handler via `Deps`/`New(...)` rather than a `&Handler{...}` literal, mirror that here instead (build the `Deps` with `AbodeStates: m` → `New(deps)`), and ensure `store`, `strings`, `encoding/json`, `net/http`, `net/http/httptest` are imported in the test file (add any missing import).

Also extend the existing `TestGetObservatorySnapshot` snapshot assertions with one line confirming the base lever fields are exposed:

```go
	if resp.Abode.AccumulationRateBP == 0 && resp.Abode.BaseWagePence == 0 {
		t.Error("abode block missing base lever values (accumulation_rate_bp / base_wage_pence)")
	}
```

(If the test decodes into a local struct, add `surplus_rate_base_bp`, `base_wage_pence`, `accumulation_rate_bp` to its `Abode` sub-struct.)

- [ ] **Step 6: Run the handler tests + full backend check**

Run: `cd services/simulation-engine && go test ./internal/transport/httpapi/ -run 'Observatory' -v 2>&1 | tail -20`
Then: `cd /mnt/c/Users/AaronHulse/IdeaProjects/capital-simulator && make vet test build 2>&1 | tail -15`
Expected: PASS — all packages, all six binaries.

- [ ] **Step 7: Commit**

```bash
git add services/simulation-engine/internal/transport/httpapi/observatory_levers_handler.go \
        services/simulation-engine/internal/transport/httpapi/observatory_handler.go \
        services/simulation-engine/internal/transport/httpapi/observatory_handler_test.go \
        services/simulation-engine/internal/transport/httpapi/routes.go \
        services/api-gateway/cmd/api-gateway/main.go
git commit --no-gpg-sign -m "feat(atlas): POST /v1/observatory/levers + base lever values on the snapshot"
```

---

# GROUP D — Frontend (the controls)

## Task D1: Wire types + api

**Files:**
- Modify: `web/src/types.ts`
- Modify: `web/src/api.ts`

- [ ] **Step 1: Add lever types + the three base fields on `AbodeReadout`**

In `web/src/types.ts`, in the Atlas Observatory block, add the three base lever fields to `AbodeReadout` (after `wage_pence`, before `law_series`):

```ts
  wage_pence: number;
  surplus_rate_base_bp: number;
  base_wage_pence: number;
  accumulation_rate_bp: number;
  law_series: GeneralLawTrendPoint[];
}
```

Then add the lever request/response types (immediately after the `AbodeReadout` interface):

```ts
export interface LeverUpdate {
  surplus_rate_base_bp?: number;
  base_wage_pence?: number;
  accumulation_rate_bp?: number;
}

export interface LeverState {
  surplus_rate_base_bp: number;
  base_wage_pence: number;
  accumulation_rate_bp: number;
}
```

- [ ] **Step 2: Add the api function**

In `web/src/api.ts`, in the Atlas Observatory section (after the `getObservatorySnapshot` entry), add:

```ts
  setObservatoryLevers: (u: import("./types").LeverUpdate) =>
    http<import("./types").LeverState>("/v1/observatory/levers", {
      method: "POST",
      body: JSON.stringify(u),
    }),
```

- [ ] **Step 3: Typecheck**

Run: `cd web && npm run lint 2>&1 | tail`
Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add web/src/types.ts web/src/api.ts
git commit --no-gpg-sign -m "feat(atlas): lever wire types + setObservatoryLevers api"
```

---

## Task D2: The `Levers` panel

**Files:**
- Create: `web/src/atlas/Levers.tsx`

- [ ] **Step 1: Write the component**

Create `web/src/atlas/Levers.tsx`:

```tsx
import { useState } from "react";
import type { AbodeReadout, LeverUpdate } from "../types";
import { api } from "../api";
import { formatBP, formatPence } from "./animation";

interface LeversProps {
  abode: AbodeReadout;
}

/** The levers — perturb the live General Law and watch it respond over the next
 *  passes. The working day sets the rate of surplus-value (necessary↔surplus);
 *  the wage is the value of labour-power; α is the share of surplus
 *  re-accumulated. Slider state is local (seeded once); the law's *response*
 *  shows in the abode readout above, not in the slider position. */
export function Levers({ abode }: LeversProps) {
  const [surplus, setSurplus] = useState(abode.surplus_rate_base_bp);
  const [wage, setWage] = useState(abode.base_wage_pence);
  const [accum, setAccum] = useState(abode.accumulation_rate_bp);

  const push = (u: LeverUpdate) => {
    void api.setObservatoryLevers(u).catch(() => {
      /* the 2s poll reflects the actual persisted state */
    });
  };

  return (
    <section className="abode-card levers" data-testid="levers">
      <div className="abode-card-k">The levers — perturb the law</div>

      <label className="lever">
        <span>Working day · rate of surplus-value <b>{formatBP(surplus)}</b></span>
        <input
          type="range" min={2000} max={40000} step={500} value={surplus}
          data-testid="lever-workingday"
          onChange={(e) => {
            const v = Number(e.target.value);
            setSurplus(v);
            push({ surplus_rate_base_bp: v });
          }}
        />
      </label>

      <label className="lever">
        <span>Wage · value of labour-power <b>{formatPence(wage)}</b></span>
        <input
          type="range" min={500} max={6000} step={100} value={wage}
          data-testid="lever-wage"
          onChange={(e) => {
            const v = Number(e.target.value);
            setWage(v);
            push({ base_wage_pence: v });
          }}
        />
      </label>

      <label className="lever">
        <span>Accumulation rate · α <b>{formatBP(accum)}</b></span>
        <input
          type="range" min={0} max={10000} step={250} value={accum}
          data-testid="lever-accumulation"
          onChange={(e) => {
            const v = Number(e.target.value);
            setAccum(v);
            push({ accumulation_rate_bp: v });
          }}
        />
      </label>
    </section>
  );
}
```

- [ ] **Step 2: Typecheck**

Run: `cd web && npm run lint 2>&1 | tail`
Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add web/src/atlas/Levers.tsx
git commit --no-gpg-sign -m "feat(atlas): Levers panel — working day, wage, accumulation rate"
```

---

## Task D3: Render the levers in the abode + style

**Files:**
- Modify: `web/src/atlas/Abode.tsx`
- Modify: `web/src/atlas/atlas.css`

- [ ] **Step 1: Render `<Levers>` in `Abode.tsx`**

In `web/src/atlas/Abode.tsx`, add the import alongside the existing `GeneralLawTrend` import:

```tsx
import { Levers } from "./Levers";
```

Then render the panel as the last child of the `.abode` container — immediately after the closing `</section>` of the general-law-trend card and before the final `</div>` that closes `<div className="abode" …>`:

```tsx
      <Levers abode={abode} />
    </div>
  );
}
```

- [ ] **Step 2: Append the lever styles**

Append to `web/src/atlas/atlas.css`:

```css
/* --- The levers (Slice 3) ------------------------------------------------- */
.levers { display: flex; flex-direction: column; gap: 12px; }
.lever { display: flex; flex-direction: column; gap: 4px; }
.lever > span { color: #cdd3e0; font-size: 13px; }
.lever > span > b { color: #c8a240; font-weight: 600; }
.lever > input[type="range"] {
  width: 100%;
  accent-color: #c8a240;
  cursor: pointer;
}
```

- [ ] **Step 3: Typecheck + build**

Run: `cd web && npm run lint && npm run build 2>&1 | tail -20`
Expected: success.

- [ ] **Step 4: Commit**

```bash
git add web/src/atlas/Abode.tsx web/src/atlas/atlas.css
git commit --no-gpg-sign -m "feat(atlas): mount the levers in the hidden abode"
```

---

# GROUP E — End-to-end acceptance

## Task E1: Boot the stack + drive the levers (Playwright MCP)

**Files:** none.

- [ ] **Step 1: Boot fresh (no new migration, but rebuild for the new code)**

Run:
```bash
docker compose down -v
docker compose up --build -d
```
Wait ~3 min for MySQL + migrations. Confirm the levers endpoint round-trips and the law responds:
```bash
curl -s -X POST http://localhost:8080/v1/engine/start >/dev/null
# Push α to 0 (stop accumulation) and crank the working day up.
curl -s -X POST http://localhost:8080/v1/observatory/levers \
  -H 'Content-Type: application/json' \
  -d '{"accumulation_rate_bp": 0, "surplus_rate_base_bp": 30000}' ; echo
# The snapshot's abode now reports the new base levers and a higher s/v.
curl -s http://localhost:8080/v1/observatory/snapshot | python3 -c "
import sys,json
a=json.load(sys.stdin)['abode']
print('surplus_rate_base_bp:', a['surplus_rate_base_bp'])
print('accumulation_rate_bp:', a['accumulation_rate_bp'])
print('rate_of_exploitation_bp:', a['rate_of_exploitation_bp'])
"
```
Expected: the POST returns `{"surplus_rate_base_bp":30000,...,"accumulation_rate_bp":0}`; the snapshot's abode reflects `surplus_rate_base_bp: 30000`, `accumulation_rate_bp: 0`, and a `rate_of_exploitation_bp` ≥ 30000 (well above the un-levered ~12000).

- [ ] **Step 2: Drive the page with Playwright MCP**

1. `browser_navigate` → `http://localhost:5173/`.
2. If transport shows ▶, click it (start the engine).
3. `browser_click` the `[data-testid="threshold"]` control → the abode appears.
4. `browser_snapshot` → confirm the `[data-testid="levers"]` panel with the three sliders (`lever-workingday`, `lever-wage`, `lever-accumulation`) is present.
5. Read the working-day surplus width (`.workingday .wd-sur`), then drag the working-day lever up via `browser_evaluate` — set the range input's value through the native setter and dispatch an `input` event:
   ```js
   () => { const el = document.querySelector('[data-testid="lever-workingday"]');
           const set = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, 'value').set;
           set.call(el, '38000'); el.dispatchEvent(new Event('input', { bubbles: true })); return el.value; }
   ```
   Wait ~6s for the POST + a poll, then re-read `.workingday .wd-sur` width → assert the surplus share **grew** (the working day responded to the lever).
6. `browser_take_screenshot` → `atlas-levers.png` for the PR.

Expected: dragging the working-day lever widens the surplus portion of the working day within a couple of polls; the law visibly responds to the operator's hand.

- [ ] **Step 3: Tear down + final check**

```bash
docker compose down
cd /mnt/c/Users/AaronHulse/IdeaProjects/capital-simulator && make vet test build
cd web && npm run lint && npm run build
```
Expected: all PASS.

---

# Done criteria (Slice 3)

- `POST /v1/observatory/levers` applies a partial `{surplus_rate_base_bp, base_wage_pence, accumulation_rate_bp}` to the live abode (clamped), persists it, and proxies through the gateway; the snapshot's `abode` block exposes the three base lever values.
- Setting α = 0 makes the law perform simple reproduction (total social capital stops growing); raising the working-day lever raises s/v; both are visible in subsequent snapshots — verified by `ApplyLevers`/store tests and live over the booted stack.
- The hidden abode shows a `Levers` panel (working day · wage · accumulation rate); dragging a lever POSTs and the law responds within a poll or two.
- `make vet test build` + `npm run lint && npm run build` pass; Playwright shows the working day widening in response to the lever.

# Completes the Atlas "General Law in Motion" design
With Slice 3 the design's three slices are all shipped: **1** real growth + corrected orbit motion, **2** the hidden abode (the law in motion), **3** the levers (the law under the operator's hand). Remaining design §11/§12 deferrals (per-capital `v` in the `capitals` array, per-capital descent, server-push streaming) stay out of scope.
