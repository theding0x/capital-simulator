# Ch. 05 — Contradictions in the General Formula of Capital: Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Prove mathematically that circulation alone cannot produce surplus-value, by adding exchange simulation types and endpoints to agent-service and wiring them into the Ch. 05 React panel.

**Architecture:** All new logic lives in agent-service. Two new pure-computation HTTP endpoints (`POST /v1/circuits`, `POST /v1/exchange-simulations`) return proofs of value-conservation without persisting results. A new `exchange.go` file holds all Ch. 05 domain types. The existing `Agent` struct gains an `Owner` class and a `labour_minutes` field (added via a goose migration). The React panel has two forms: a circuit probe and an exchange simulation.

**Tech Stack:** Go 1.25, `database/sql` + MySQL 8, `pressly/goose/v3`, React 18 + TypeScript. Module: `github.com/theding0x/capital-simulator`.

---

## Context for implementers

Ch. 04 already built `agent-service` with `Agent`, `CapitalCircuit`, `Store`, `CircuitStore`, MySQL persistence, and nine HTTP endpoints. Ch. 05 **adds onto** that implementation — do not rename or replace existing types.

**Naming reconciliation** (spec used pre-Ch.04 names; map them to what already exists):

| Spec name | Existing identifier | Action |
|---|---|---|
| `AgentID` | `agent.ID` | Reuse — same thing |
| `AgentKind` | `agent.Class` | Reuse — same thing |
| `AgentKindCapitalist` | `agent.Capitalist` | Exists |
| `AgentKindOwner` | `agent.Owner` | **Add new Class constant** |
| `MoneyAmount` | `agent.Pence` | Reuse — same thing |
| `LabourMinutes` | new `int64` field on `Agent` | **Add to struct + DB** |
| `Circuit` (M-C-M') | `agent.CapitalCircuit` | Reuse — same thing |
| `MerchantsCapital` | new type | **Add in `exchange.go`** |
| `UsurersCapital` | new type | **Add in `exchange.go`** |
| `ExchangeResult` | new type | **Add in `exchange.go`** |

**Route conflict:** The spec proposes `POST /v1/exchanges`, but market-service already owns that path. Use `POST /v1/exchange-simulations` instead — same semantics, no gateway conflict.

---

## File Map

**New files:**
- `services/agent-service/internal/agent/exchange.go` — ExchangeResult, MerchantsCapital, UsurersCapital, ExchangeEquivalents, ExchangeNonEquivalents, TotalValue, Role constants
- `services/agent-service/internal/agent/exchange_test.go` — all six spec fixtures as tests
- `services/agent-service/internal/store/migrations/00004_ch05_labour_minutes.sql` — ALTER TABLE adds `labour_minutes` column
- `web/src/chapters/Ch05Contradictions.tsx` — Ch. 05 UI panel

**Modified files:**
- `services/agent-service/internal/agent/agent.go` — add `Owner` Class constant, add `LabourMinutes int64` to `Agent` struct, update `Validate()`
- `services/agent-service/internal/store/store.go` — add `LabourMinutes *int64` to `Update` struct, update `IsEmpty()` and `Apply()`
- `services/agent-service/internal/store/mysql.go` — add `labour_minutes` to all agent SQL queries and scan helpers
- `services/agent-service/internal/transport/httpapi/handler.go` — add `ComputeCircuit` and `ComputeExchange` handlers; add `labour_minutes` to create/update request structs; extend `CreateCircuit` to block Owner from M-C-M'
- `services/agent-service/internal/transport/httpapi/routes.go` — register new routes
- `services/api-gateway/cmd/api-gateway/main.go` — proxy `/v1/circuits` and `/v1/exchange-simulations` to agent-service
- `web/src/types.ts` — add `owner` to Agent class union, add `labour_minutes` field, add CircuitProof, ExchangeSimulation, ComputeCircuitInput, SimulateExchangeInput
- `web/src/api.ts` — add `computeCircuit` and `simulateExchange` API methods
- `web/src/components/ChapterShell.tsx` — import Ch05Contradictions, add ch05 quote, wire render branch
- `web/src/chapters/registry.ts` — mark ch05 as `"done"`
- `docs/architecture.md` — update chapter status table

---

## Task 1: Exchange domain types and functions

**Files:**
- Create: `services/agent-service/internal/agent/exchange.go`
- Create: `services/agent-service/internal/agent/exchange_test.go`

- [ ] **Step 1: Write the failing tests**

Create `services/agent-service/internal/agent/exchange_test.go`:

```go
package agent_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// §1: Exchange wine=50 for corn=50; both break even; total value unchanged.
func TestExchangeEquivalents_WineForCorn(t *testing.T) {
	t.Parallel()
	r := agent.ExchangeEquivalents(50, 50)
	if r.AAfter != 50 {
		t.Errorf("A after: want 50, got %d", r.AAfter)
	}
	if r.BAfter != 50 {
		t.Errorf("B after: want 50, got %d", r.BAfter)
	}
	if r.TotalBefore() != r.TotalAfter() {
		t.Errorf("total not conserved: before=%d after=%d", r.TotalBefore(), r.TotalAfter())
	}
	if r.AAfter-r.ABefore != 0 {
		t.Errorf("surplus for A: want 0, got %d", r.AAfter-r.ABefore)
	}
}

// §2: Seller of worth-100 sells at 110; gains 10, buyer loses 10, total 210.
func TestExchangeNonEquivalents_SellerAboveValue(t *testing.T) {
	t.Parallel()
	r := agent.ExchangeNonEquivalents(100, 110)
	sellerGain := r.AAfter - r.ABefore
	buyerGain := r.BAfter - r.BBefore
	if sellerGain != 10 {
		t.Errorf("seller gain: want 10, got %d", sellerGain)
	}
	if buyerGain != -10 {
		t.Errorf("buyer gain: want -10, got %d", buyerGain)
	}
	if r.TotalBefore() != r.TotalAfter() {
		t.Errorf("total not conserved: before=%d after=%d", r.TotalBefore(), r.TotalAfter())
	}
	if r.TotalBefore() != 210 {
		t.Errorf("total before: want 210, got %d", r.TotalBefore())
	}
}

// §3: A=40, B=50, total 90; after swap A=50, B=40, total still 90.
func TestExchangeNonEquivalents_WineForCorn(t *testing.T) {
	t.Parallel()
	r := agent.ExchangeNonEquivalents(40, 50)
	if r.AAfter != 50 {
		t.Errorf("A after: want 50, got %d", r.AAfter)
	}
	if r.BAfter != 40 {
		t.Errorf("B after: want 40, got %d", r.BAfter)
	}
	if r.TotalBefore() != 90 || r.TotalAfter() != 90 {
		t.Errorf("total not conserved: before=%d after=%d", r.TotalBefore(), r.TotalAfter())
	}
}

// §4: Property — for any (x, y), TotalBefore == TotalAfter.
func TestExchange_ValueConservation_Property(t *testing.T) {
	t.Parallel()
	cases := [][2]agent.Pence{
		{0, 0},
		{100, 100},
		{100, 200},
		{50, 150},
		{9999, 1},
	}
	for _, c := range cases {
		r := agent.ExchangeNonEquivalents(c[0], c[1])
		if r.TotalBefore() != r.TotalAfter() {
			t.Errorf("(%d,%d) non-equiv: total not conserved: before=%d after=%d",
				c[0], c[1], r.TotalBefore(), r.TotalAfter())
		}
		r2 := agent.ExchangeEquivalents(c[0], c[0])
		if r2.TotalBefore() != r2.TotalAfter() {
			t.Errorf("equiv(%d): total not conserved", c[0])
		}
	}
}

// §5: Scaling all values by k leaves SurplusValue==0 and ratios unchanged.
func TestMerchantsCapital_ScalingInvariant(t *testing.T) {
	t.Parallel()
	mc := agent.MerchantsCapital{M: 100, CommodityID: "cotton", MPrime: 100}
	if mc.SurplusValue() != 0 {
		t.Errorf("surplus: want 0, got %d", mc.SurplusValue())
	}
	scaled := agent.MerchantsCapital{M: 300, CommodityID: "cotton", MPrime: 300}
	if scaled.SurplusValue() != 0 {
		t.Errorf("scaled surplus: want 0, got %d", scaled.SurplusValue())
	}
}

// §6: M-M' (UsurersCapital) — SurplusValue=10; no commodity field to locate
// the source.
func TestUsurersCapital_SurplusValue(t *testing.T) {
	t.Parallel()
	uc := agent.UsurersCapital{M: 100, MPrime: 110}
	if uc.SurplusValue() != 10 {
		t.Errorf("surplus: want 10, got %d", uc.SurplusValue())
	}
}

// MerchantsCapital.Origin() returns "equivalent" / "redistribution".
func TestMerchantsCapital_Origin(t *testing.T) {
	t.Parallel()
	equiv := agent.MerchantsCapital{M: 100, CommodityID: "cotton", MPrime: 100}
	if equiv.Origin() != "equivalent" {
		t.Errorf("zero-surplus origin: want 'equivalent', got %q", equiv.Origin())
	}
	nonEquiv := agent.MerchantsCapital{M: 100, CommodityID: "cotton", MPrime: 110}
	if nonEquiv.Origin() != "redistribution" {
		t.Errorf("non-zero-surplus origin: want 'redistribution', got %q", nonEquiv.Origin())
	}
}

// TotalValue sums MoneyBalance across agents; invariant across exchange.
func TestTotalValue_Conservation(t *testing.T) {
	t.Parallel()
	a := agent.Agent{MoneyBalance: 100}
	b := agent.Agent{MoneyBalance: 200}
	before := agent.TotalValue([]agent.Agent{a, b})
	if before != 300 {
		t.Errorf("want 300, got %d", before)
	}
	a.MoneyBalance = 200
	b.MoneyBalance = 100
	after := agent.TotalValue([]agent.Agent{a, b})
	if after != before {
		t.Errorf("TotalValue not conserved: before=%d after=%d", before, after)
	}
}

// ExchangeNonEquivalents: zero-sum property — A's gain + B's gain == 0.
func TestExchangeNonEquivalents_ZeroSum(t *testing.T) {
	t.Parallel()
	r := agent.ExchangeNonEquivalents(100, 110)
	aGain := r.AAfter - r.ABefore
	bGain := r.BAfter - r.BBefore
	if aGain+bGain != 0 {
		t.Errorf("zero-sum violated: aGain=%d bGain=%d", aGain, bGain)
	}
}
```

- [ ] **Step 2: Run tests — expect compilation failure**

```bash
cd /path/to/capital-simulator && go test ./services/agent-service/internal/agent/...
```

Expected: `undefined: agent.ExchangeEquivalents` (and similar). This is correct — the types don't exist yet.

- [ ] **Step 3: Write exchange.go**

Create `services/agent-service/internal/agent/exchange.go`:

```go
package agent

// Role identifies a party's position in a bilateral exchange.
type Role string

const (
	RoleSeller Role = "seller"
	RoleBuyer  Role = "buyer"
)

// ExchangeResult records the outcome of a bilateral exchange between parties A
// and B. The invariant TotalBefore() == TotalAfter() must always hold.
type ExchangeResult struct {
	ABefore Pence  `json:"a_before"`
	BBefore Pence  `json:"b_before"`
	AAfter  Pence  `json:"a_after"`
	BAfter  Pence  `json:"b_after"`
	Origin  string `json:"origin"` // "equivalent" or "redistribution"
}

// TotalBefore returns the combined value of both parties before exchange.
func (r ExchangeResult) TotalBefore() Pence { return r.ABefore + r.BBefore }

// TotalAfter returns the combined value of both parties after exchange.
func (r ExchangeResult) TotalAfter() Pence { return r.AAfter + r.BAfter }

// MerchantsCapital represents M-C-M' operating purely within the sphere of
// circulation. Any surplus-value arises from redistribution, not creation.
type MerchantsCapital struct {
	M           Pence  `json:"m"`
	CommodityID string `json:"commodity_id"`
	MPrime      Pence  `json:"m_prime"`
}

// SurplusValue returns MPrime - M.
func (mc MerchantsCapital) SurplusValue() Pence { return mc.MPrime - mc.M }

// Origin returns "equivalent" when SurplusValue is zero, "redistribution"
// otherwise. It never returns "creation" — circulation cannot create value.
func (mc MerchantsCapital) Origin() string {
	if mc.SurplusValue() == 0 {
		return "equivalent"
	}
	return "redistribution"
}

// UsurersCapital is the degenerate M-M' circuit: money exchanged for more
// money without a commodity intermediary. The source of the surplus cannot be
// located within the circuit.
type UsurersCapital struct {
	M      Pence `json:"m"`
	MPrime Pence `json:"m_prime"`
}

// SurplusValue returns MPrime - M.
func (uc UsurersCapital) SurplusValue() Pence { return uc.MPrime - uc.M }

// ExchangeEquivalents models a bilateral swap of equal values: A's commodity
// worth aValue trades for B's commodity worth bValue. Neither party gains
// surplus-value.
func ExchangeEquivalents(aValue, bValue Pence) ExchangeResult {
	return ExchangeResult{
		ABefore: aValue,
		BBefore: bValue,
		AAfter:  bValue,
		BAfter:  aValue,
		Origin:  "equivalent",
	}
}

// ExchangeNonEquivalents models a seller (A) obtaining price above commodity
// value. A gains (price − sellerValue); B loses the same amount. Total social
// value is conserved.
func ExchangeNonEquivalents(sellerValue, price Pence) ExchangeResult {
	return ExchangeResult{
		ABefore: sellerValue,
		BBefore: price,
		AAfter:  price,
		BAfter:  sellerValue,
		Origin:  "redistribution",
	}
}

// TotalValue returns the sum of MoneyBalance across all agents. This is the
// social total that exchange cannot increase.
func TotalValue(agents []Agent) Pence {
	var total Pence
	for _, a := range agents {
		total += a.MoneyBalance
	}
	return total
}
```

- [ ] **Step 4: Run tests — expect pass**

```bash
go test ./services/agent-service/internal/agent/...
```

Expected: all tests pass. Output example:
```
ok  github.com/theding0x/capital-simulator/services/agent-service/internal/agent  0.XXXs
```

- [ ] **Step 5: Commit**

```bash
git add services/agent-service/internal/agent/exchange.go \
        services/agent-service/internal/agent/exchange_test.go
git commit -m "feat(agent): add Ch.05 exchange domain types and value-conservation functions

ExchangeResult, MerchantsCapital, UsurersCapital, ExchangeEquivalents,
ExchangeNonEquivalents, TotalValue — proves circulation cannot produce
surplus-value (Capital Vol. I, Ch. 5).

All six spec fixtures pass as tests (§1–§6)."
```

---

## Task 2: Add Owner class, LabourMinutes field, and migration

**Files:**
- Modify: `services/agent-service/internal/agent/agent.go`
- Modify: `services/agent-service/internal/store/store.go`
- Modify: `services/agent-service/internal/store/mysql.go`
- Create: `services/agent-service/internal/store/migrations/00004_ch05_labour_minutes.sql`

- [ ] **Step 1: Update agent.go — add Owner class and LabourMinutes**

In `services/agent-service/internal/agent/agent.go`, make two changes:

**a) Add `Owner` constant alongside the existing Class constants:**

```go
const (
	Capitalist Class = "capitalist"
	Worker     Class = "worker"
	Miser      Class = "miser"
	Owner      Class = "owner"
)
```

**b) Add `LabourMinutes int64` to the Agent struct (between MoneyBalance and Hoarding):**

```go
type Agent struct {
	ID            ID        `json:"id"`
	Name          string    `json:"name"`
	Class         Class     `json:"class"`
	MoneyBalance  Pence     `json:"money_balance"`
	LabourMinutes int64     `json:"labour_minutes"`
	Hoarding      bool      `json:"hoarding"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
```

**c) Update Validate() to accept Owner:**

```go
func (a Agent) Validate() error {
	if a.Name == "" {
		return errors.New("agent: name is required")
	}
	switch a.Class {
	case Capitalist, Worker, Miser, Owner:
	default:
		return errors.New("agent: unknown class")
	}
	if a.MoneyBalance < 0 {
		return errors.New("agent: money_balance cannot be negative")
	}
	return nil
}
```

- [ ] **Step 2: Run existing tests — expect pass**

```bash
go test ./services/agent-service/internal/agent/...
```

Expected: all existing tests still pass (LabourMinutes defaults to zero; Owner is accepted by Validate).

- [ ] **Step 3: Update store.go — add LabourMinutes to Update**

Replace the `Update` struct and its methods in `services/agent-service/internal/store/store.go`:

```go
type Update struct {
	Name          *string
	MoneyBalance  *agent.Pence
	Hoarding      *bool
	LabourMinutes *int64
}

func (u Update) IsEmpty() bool {
	return u.Name == nil && u.MoneyBalance == nil && u.Hoarding == nil && u.LabourMinutes == nil
}

func (u Update) Apply(a agent.Agent) agent.Agent {
	out := a
	if u.Name != nil {
		out.Name = *u.Name
	}
	if u.MoneyBalance != nil {
		out.MoneyBalance = *u.MoneyBalance
	}
	if u.Hoarding != nil {
		out.Hoarding = *u.Hoarding
	}
	if u.LabourMinutes != nil {
		out.LabourMinutes = *u.LabourMinutes
	}
	return out
}
```

- [ ] **Step 4: Create migration 00004**

Create `services/agent-service/internal/store/migrations/00004_ch05_labour_minutes.sql`:

```sql
-- +goose Up
ALTER TABLE agents ADD COLUMN labour_minutes BIGINT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE agents DROP COLUMN labour_minutes;
```

- [ ] **Step 5: Update mysql.go — add labour_minutes to all agent queries**

In `services/agent-service/internal/store/mysql.go`, make the following changes:

**a) `Create` — add column and value:**

```go
func (m *MySQL) Create(ctx context.Context, a agent.Agent) (agent.Agent, error) {
	if err := a.Validate(); err != nil {
		return agent.Agent{}, err
	}
	if a.ID.IsZero() {
		a.ID = agent.NewID()
	}
	now := m.now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now
	const q = `INSERT INTO agents (id, name, class, money_balance, hoarding, labour_minutes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID), a.Name, string(a.Class), int64(a.MoneyBalance), a.Hoarding,
		a.LabourMinutes, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return agent.Agent{}, err
	}
	return a, nil
}
```

**b) `Get` — add column to SELECT:**

```go
func (m *MySQL) Get(ctx context.Context, id agent.ID) (agent.Agent, error) {
	const q = `SELECT id, name, class, money_balance, hoarding, labour_minutes, created_at, updated_at
		FROM agents WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanAgent(row)
}
```

**c) `List` — add column to SELECT:**

```go
func (m *MySQL) List(ctx context.Context) ([]agent.Agent, error) {
	const q = `SELECT id, name, class, money_balance, hoarding, labour_minutes, created_at, updated_at
		FROM agents ORDER BY name ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAgents(rows)
}
```

**d) `ListByClass` — add column to SELECT:**

```go
func (m *MySQL) ListByClass(ctx context.Context, class agent.Class) ([]agent.Agent, error) {
	const q = `SELECT id, name, class, money_balance, hoarding, labour_minutes, created_at, updated_at
		FROM agents WHERE class = ? ORDER BY name ASC`
	rows, err := m.db.QueryContext(ctx, q, string(class))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAgents(rows)
}
```

**e) `Update` — add labour_minutes to SET clause:**

```go
func (m *MySQL) Update(ctx context.Context, id agent.ID, u Update) (agent.Agent, error) {
	if u.IsEmpty() {
		return m.Get(ctx, id)
	}
	cur, err := m.Get(ctx, id)
	if err != nil {
		return agent.Agent{}, err
	}
	next := u.Apply(cur)
	if err := next.Validate(); err != nil {
		return agent.Agent{}, err
	}
	next.UpdatedAt = m.now().UTC()
	const q = `UPDATE agents SET name = ?, class = ?, money_balance = ?, hoarding = ?, labour_minutes = ?, updated_at = ?
		WHERE id = ?`
	res, err := m.db.ExecContext(ctx, q,
		next.Name, string(next.Class), int64(next.MoneyBalance), next.Hoarding,
		next.LabourMinutes, next.UpdatedAt, string(id),
	)
	if err != nil {
		return agent.Agent{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return agent.Agent{}, ErrNotFound
	}
	return next, nil
}
```

**f) `scanAgent` — add labour_minutes scan:**

```go
func scanAgent(row *sql.Row) (agent.Agent, error) {
	var a agent.Agent
	var id, class string
	var balance int64
	var hoarding bool
	err := row.Scan(&id, &a.Name, &class, &balance, &hoarding, &a.LabourMinutes, &a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.Agent{}, ErrNotFound
	}
	if err != nil {
		return agent.Agent{}, err
	}
	a.ID = agent.ID(id)
	a.Class = agent.Class(class)
	a.MoneyBalance = agent.Pence(balance)
	a.Hoarding = hoarding
	return a, nil
}
```

**g) `scanAgents` — add labour_minutes scan:**

```go
func scanAgents(rows *sql.Rows) ([]agent.Agent, error) {
	var out []agent.Agent
	for rows.Next() {
		var a agent.Agent
		var id, class string
		var balance int64
		var hoarding bool
		if err := rows.Scan(&id, &a.Name, &class, &balance, &hoarding, &a.LabourMinutes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		a.ID = agent.ID(id)
		a.Class = agent.Class(class)
		a.MoneyBalance = agent.Pence(balance)
		a.Hoarding = hoarding
		out = append(out, a)
	}
	return out, rows.Err()
}
```

- [ ] **Step 6: Run all store tests — expect pass**

```bash
go test ./services/agent-service/...
```

Expected: all tests pass (memory store uses the Agent struct directly; no SQL changes needed there).

- [ ] **Step 7: Commit**

```bash
git add services/agent-service/internal/agent/agent.go \
        services/agent-service/internal/store/store.go \
        services/agent-service/internal/store/mysql.go \
        services/agent-service/internal/store/migrations/00004_ch05_labour_minutes.sql
git commit -m "feat(agent): add Owner class, LabourMinutes field, and DB migration

Owner is a simple commodity owner who can buy and sell equivalents but
cannot originate surplus-value (Capital Vol. I, Ch. 5).

labour_minutes tracks the abstract-labour magnitude an agent commands,
carried over from Ch. 1's LabourMinutes unit."
```

---

## Task 3: HTTP handlers for circuit probe and exchange simulation

**Files:**
- Modify: `services/agent-service/internal/transport/httpapi/handler.go`
- Modify: `services/agent-service/internal/transport/httpapi/routes.go`

- [ ] **Step 1: Update handler.go**

In `services/agent-service/internal/transport/httpapi/handler.go`, make three changes:

**a) Add `labour_minutes` to `createAgentRequest` and `updateAgentRequest`:**

```go
type createAgentRequest struct {
	Name          string      `json:"name"`
	Class         agent.Class `json:"class"`
	MoneyBalance  agent.Pence `json:"money_balance"`
	LabourMinutes int64       `json:"labour_minutes"`
}

type updateAgentRequest struct {
	Name          *string      `json:"name,omitempty"`
	MoneyBalance  *agent.Pence `json:"money_balance,omitempty"`
	LabourMinutes *int64       `json:"labour_minutes,omitempty"`
}
```

**b) Wire `LabourMinutes` in `Create` and `Update` handlers:**

In `Create`, update the `Agent` literal:
```go
a := agent.Agent{
	Name:          strings.TrimSpace(req.Name),
	Class:         req.Class,
	MoneyBalance:  req.MoneyBalance,
	LabourMinutes: req.LabourMinutes,
}
```

In `Update`, add the field to the `store.Update`:
```go
updated, err := h.Store.Update(r.Context(), id, store.Update{
	Name:          req.Name,
	MoneyBalance:  req.MoneyBalance,
	LabourMinutes: req.LabourMinutes,
})
```

**c) Extend `CreateCircuit` to block Owner from M-C-M':**

Find this existing block:
```go
if a.Class == agent.Worker && req.CircuitType == agent.CircuitMCM {
    writeError(w, http.StatusBadRequest, agent.ErrWrongClass.Error())
    return
}
```

Replace with:
```go
if (a.Class == agent.Worker || a.Class == agent.Owner) && req.CircuitType == agent.CircuitMCM {
	writeError(w, http.StatusBadRequest, agent.ErrWrongClass.Error())
	return
}
```

**d) Add two new request/response structs and two new handlers at the bottom of handler.go (before `decodeJSON`):**

```go
type computeCircuitRequest struct {
	M           agent.Pence `json:"m"`
	CommodityID string      `json:"commodity_id"`
	MPrime      agent.Pence `json:"m_prime"`
}

type computeCircuitResponse struct {
	M            agent.Pence `json:"m"`
	CommodityID  string      `json:"commodity_id,omitempty"`
	MPrime       agent.Pence `json:"m_prime"`
	SurplusValue agent.Pence `json:"surplus_value"`
	Origin       string      `json:"origin"`
}

type computeExchangeRequest struct {
	AValue agent.Pence `json:"a_value"`
	BValue agent.Pence `json:"b_value"`
}

// ComputeCircuit is a stateless endpoint that computes surplus-value and
// origin for a given circuit. If commodity_id is blank, interprets as
// UsurersCapital (M-M'); otherwise MerchantsCapital (M-C-M').
func (h *Handler) ComputeCircuit(w http.ResponseWriter, r *http.Request) {
	var req computeCircuitRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.M <= 0 {
		writeError(w, http.StatusBadRequest, "m must be positive")
		return
	}
	if req.MPrime < 0 {
		writeError(w, http.StatusBadRequest, "m_prime cannot be negative")
		return
	}
	var surplusValue agent.Pence
	var origin string
	if req.CommodityID != "" {
		mc := agent.MerchantsCapital{M: req.M, CommodityID: req.CommodityID, MPrime: req.MPrime}
		surplusValue = mc.SurplusValue()
		origin = mc.Origin()
	} else {
		uc := agent.UsurersCapital{M: req.M, MPrime: req.MPrime}
		surplusValue = uc.SurplusValue()
		origin = "redistribution"
	}
	writeJSON(w, http.StatusOK, computeCircuitResponse{
		M:            req.M,
		CommodityID:  req.CommodityID,
		MPrime:       req.MPrime,
		SurplusValue: surplusValue,
		Origin:       origin,
	})
}

// ComputeExchange is a stateless endpoint that simulates a bilateral exchange
// and returns an ExchangeResult proving value conservation.
func (h *Handler) ComputeExchange(w http.ResponseWriter, r *http.Request) {
	var req computeExchangeRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.AValue < 0 || req.BValue < 0 {
		writeError(w, http.StatusBadRequest, "values cannot be negative")
		return
	}
	var result agent.ExchangeResult
	if req.AValue == req.BValue {
		result = agent.ExchangeEquivalents(req.AValue, req.BValue)
	} else {
		result = agent.ExchangeNonEquivalents(req.AValue, req.BValue)
	}
	writeJSON(w, http.StatusOK, result)
}
```

- [ ] **Step 2: Update routes.go — register new routes**

Replace the full contents of `services/agent-service/internal/transport/httpapi/routes.go`:

```go
package httpapi

import "github.com/theding0x/capital-simulator/pkg/httpx"

func Register(s *httpx.Server, h *Handler) {
	s.HandleFunc("POST /v1/agents", h.Create)
	s.HandleFunc("GET /v1/agents", h.List)
	s.HandleFunc("GET /v1/agents/{id}", h.Get)
	s.HandleFunc("PATCH /v1/agents/{id}", h.Update)
	s.HandleFunc("DELETE /v1/agents/{id}", h.Delete)
	s.HandleFunc("POST /v1/agents/{id}/circuits", h.CreateCircuit)
	s.HandleFunc("GET /v1/agents/{id}/circuits", h.ListCircuits)
	s.HandleFunc("POST /v1/agents/{id}/reinvest", h.Reinvest)
	s.HandleFunc("POST /v1/agents/{id}/hoard", h.Hoard)
	s.HandleFunc("POST /v1/circuits", h.ComputeCircuit)
	s.HandleFunc("POST /v1/exchange-simulations", h.ComputeExchange)
}
```

- [ ] **Step 3: Run all tests — expect pass**

```bash
go test ./services/agent-service/...
```

Expected: all tests pass.

- [ ] **Step 4: Commit**

```bash
git add services/agent-service/internal/transport/httpapi/handler.go \
        services/agent-service/internal/transport/httpapi/routes.go
git commit -m "feat(agent): add Ch.05 circuit-probe and exchange-simulation endpoints

POST /v1/circuits — stateless MerchantsCapital / UsurersCapital probe.
POST /v1/exchange-simulations — stateless bilateral exchange with value-
conservation proof.

Also extends CreateCircuit to block Owner agents from M-C-M' circuits."
```

---

## Task 4: API gateway routes

**Files:**
- Modify: `services/api-gateway/cmd/api-gateway/main.go`

- [ ] **Step 1: Add /v1/circuits and /v1/exchange-simulations to agent proxy**

In `services/api-gateway/cmd/api-gateway/main.go`, find the existing agent proxy block:

```go
srv.Handle("/v1/agents", agentProxy)
srv.Handle("/v1/agents/{rest...}", agentProxy)
```

Add two more lines immediately after:

```go
srv.Handle("/v1/agents", agentProxy)
srv.Handle("/v1/agents/{rest...}", agentProxy)
srv.Handle("/v1/circuits", agentProxy)
srv.Handle("/v1/exchange-simulations", agentProxy)
```

- [ ] **Step 2: Run vet**

```bash
go vet ./services/api-gateway/...
```

Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add services/api-gateway/cmd/api-gateway/main.go
git commit -m "feat(gateway): proxy /v1/circuits and /v1/exchange-simulations to agent-service"
```

---

## Task 5: TypeScript types and API client

**Files:**
- Modify: `web/src/types.ts`
- Modify: `web/src/api.ts`

- [ ] **Step 1: Update types.ts**

**a) Add `"owner"` to the Agent class union and add `labour_minutes` field.**

Find the existing `Agent` interface and replace it:

```typescript
export interface Agent {
  id: string;
  name: string;
  class: "capitalist" | "worker" | "miser" | "owner";
  money_balance: number; // Pence (pennies); divide by 100 for £
  labour_minutes: number;
  hoarding: boolean;
  created_at: string;
  updated_at: string;
}
```

**b) Add `"owner"` to `CreateAgentInput.class` and add optional `labour_minutes`:**

```typescript
export interface CreateAgentInput {
  name: string;
  class: "capitalist" | "worker" | "miser" | "owner";
  money_balance: number;
  labour_minutes?: number;
}
```

**c) Append new Ch. 05 types at the end of the file (after the existing Ch. 4 block):**

```typescript
// --- agent-service types (Ch. 5: Contradictions in the General Formula) ----

export interface CircuitProof {
  m: number; // Pence
  commodity_id?: string;
  m_prime: number; // Pence
  surplus_value: number; // Pence
  origin: "equivalent" | "redistribution";
}

export interface ExchangeSimulation {
  a_before: number; // Pence
  b_before: number;
  a_after: number;
  b_after: number;
  origin: "equivalent" | "redistribution";
}

export interface ComputeCircuitInput {
  m: number;
  commodity_id?: string;
  m_prime: number;
}

export interface SimulateExchangeInput {
  a_value: number;
  b_value: number;
}
```

- [ ] **Step 2: Update api.ts**

**a) Add imports at the top:**

```typescript
import type {
  ...existing imports...
  CircuitProof,
  ComputeCircuitInput,
  ExchangeSimulation,
  SimulateExchangeInput,
} from "./types";
```

**b) Add two methods at the end of the `api` object (inside the closing `}`), after `hoardAgent`:**

```typescript
  computeCircuit: (input: ComputeCircuitInput) =>
    http<CircuitProof>("/v1/circuits", { method: "POST", body: JSON.stringify(input) }),

  simulateExchange: (input: SimulateExchangeInput) =>
    http<ExchangeSimulation>("/v1/exchange-simulations", {
      method: "POST",
      body: JSON.stringify(input),
    }),
```

- [ ] **Step 3: Run TypeScript check**

```bash
cd web && npm run lint
```

Expected: no type errors.

- [ ] **Step 4: Commit**

```bash
git add web/src/types.ts web/src/api.ts
git commit -m "feat(web/types): add Ch.05 CircuitProof, ExchangeSimulation types and API methods

Also extends Agent with labour_minutes and owner class."
```

---

## Task 6: Ch05 React panel

**Files:**
- Create: `web/src/chapters/Ch05Contradictions.tsx`

- [ ] **Step 1: Write Ch05Contradictions.tsx**

Create `web/src/chapters/Ch05Contradictions.tsx`:

```tsx
import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type { CircuitProof, ExchangeSimulation } from "../types";

interface Ch05Props {
  onSharedChanged: () => void;
}

function penceToGBP(pence: number): string {
  return `£${(pence / 100).toFixed(2)}`;
}

export function Ch05Contradictions({ onSharedChanged: _onSharedChanged }: Ch05Props) {
  return (
    <>
      <CircuitProbePanel />
      <ExchangeSimulationPanel />
    </>
  );
}

function CircuitProbePanel() {
  const [m, setM] = useState(10000);
  const [commodityId, setCommodityId] = useState("cotton");
  const [mPrime, setMPrime] = useState(10000);
  const [result, setResult] = useState<CircuitProof | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    setResult(null);
    try {
      const r = await api.computeCircuit({
        m,
        commodity_id: commodityId || undefined,
        m_prime: mPrime,
      });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Circuit Probe</h2>
      <p className="muted small">
        Test whether M—C—M′ can produce surplus-value through circulation alone.
        Leave commodity blank for M—M′ (usurer&rsquo;s capital).
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          M advanced (pence)
          <input
            type="number"
            value={m}
            onChange={(e) => setM(Number(e.target.value))}
            min={1}
          />
        </label>
        <label>
          Commodity ID (blank = M—M′)
          <input value={commodityId} onChange={(e) => setCommodityId(e.target.value)} />
        </label>
        <label>
          M′ returned (pence)
          <input
            type="number"
            value={mPrime}
            onChange={(e) => setMPrime(Number(e.target.value))}
            min={0}
          />
        </label>
        <button type="submit">Compute</button>
        {err && <span className="error">{err}</span>}
      </form>
      {result && (
        <div className="item-card">
          <div className="item-header">
            <span className="item-name">
              {penceToGBP(result.m)}
              {result.commodity_id ? ` → C (${result.commodity_id.slice(0, 8)}…) →` : " →"}
              {" "}{penceToGBP(result.m_prime)}
            </span>
            <span className={`item-tag${result.origin === "redistribution" ? " negative" : ""}`}>
              {result.origin}
            </span>
          </div>
          <p className="small muted">
            ∆M = {result.surplus_value >= 0 ? "+" : ""}
            {penceToGBP(result.surplus_value)}.{" "}
            {result.origin === "redistribution"
              ? "Value was redistributed between parties, not created."
              : "Exchange at value: no surplus arose from circulation."}
          </p>
        </div>
      )}
    </section>
  );
}

function ExchangeSimulationPanel() {
  const [aValue, setAValue] = useState(5000);
  const [bValue, setBValue] = useState(5000);
  const [result, setResult] = useState<ExchangeSimulation | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    setResult(null);
    try {
      const r = await api.simulateExchange({ a_value: aValue, b_value: bValue });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  const totalBefore = result ? result.a_before + result.b_before : 0;
  const totalAfter = result ? result.a_after + result.b_after : 0;

  return (
    <section className="card">
      <h2>Exchange Simulation</h2>
      <p className="muted small">
        Prove that bilateral exchange conserves total social value.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          A holds (pence)
          <input
            type="number"
            value={aValue}
            onChange={(e) => setAValue(Number(e.target.value))}
            min={0}
          />
        </label>
        <label>
          B holds (pence)
          <input
            type="number"
            value={bValue}
            onChange={(e) => setBValue(Number(e.target.value))}
            min={0}
          />
        </label>
        <button type="submit">Simulate</button>
        {err && <span className="error">{err}</span>}
      </form>
      {result && (
        <>
          <table className="data-table">
            <thead>
              <tr>
                <th>Party</th>
                <th>Before</th>
                <th>After</th>
                <th>∆</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>A</td>
                <td>{penceToGBP(result.a_before)}</td>
                <td>{penceToGBP(result.a_after)}</td>
                <td
                  className={
                    result.a_after > result.a_before
                      ? "positive"
                      : result.a_after < result.a_before
                      ? "negative"
                      : ""
                  }
                >
                  {result.a_after >= result.a_before ? "+" : ""}
                  {penceToGBP(result.a_after - result.a_before)}
                </td>
              </tr>
              <tr>
                <td>B</td>
                <td>{penceToGBP(result.b_before)}</td>
                <td>{penceToGBP(result.b_after)}</td>
                <td
                  className={
                    result.b_after > result.b_before
                      ? "positive"
                      : result.b_after < result.b_before
                      ? "negative"
                      : ""
                  }
                >
                  {result.b_after >= result.b_before ? "+" : ""}
                  {penceToGBP(result.b_after - result.b_before)}
                </td>
              </tr>
              <tr>
                <td>
                  <strong>Total</strong>
                </td>
                <td>{penceToGBP(totalBefore)}</td>
                <td>{penceToGBP(totalAfter)}</td>
                <td>{penceToGBP(totalAfter - totalBefore)}</td>
              </tr>
            </tbody>
          </table>
          <p className="small muted">
            Origin: <strong>{result.origin}</strong>.{" "}
            {result.origin === "redistribution"
              ? "A's gain is exactly B's loss. No new value was created."
              : "Values are equal: exchange conserves value perfectly."}
          </p>
        </>
      )}
    </section>
  );
}
```

- [ ] **Step 2: Run TypeScript check**

```bash
cd web && npm run lint
```

Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add web/src/chapters/Ch05Contradictions.tsx
git commit -m "feat(web): add Ch.05 React panel

Circuit Probe form: input M / commodity / M', compute surplus-value and
origin tag (equivalent vs redistribution).

Exchange Simulation form: input two party values, display before/after
table proving TotalValue is conserved."
```

---

## Task 7: Wire Ch05 into the UI shell and mark done

**Files:**
- Modify: `web/src/components/ChapterShell.tsx`
- Modify: `web/src/chapters/registry.ts`

- [ ] **Step 1: Update ChapterShell.tsx**

**a) Add import after the Ch04 import line:**

```typescript
import { Ch05Contradictions } from "../chapters/Ch05Contradictions";
```

**b) Add ch05 quote to the QUOTES object:**

```typescript
  ch05: "Circulation, or the exchange of commodities, begets no value.",
```

**c) Add render branch after the ch04 branch:**

```tsx
        ) : activeChapterId === "ch05" ? (
          <Ch05Contradictions onSharedChanged={onSharedChanged} />
```

The full updated ternary chain (from ch01 to ch05) should look like:

```tsx
        ) : activeChapterId === "ch01" ? (
          <Ch01Commodity commodities={commodities} onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch02" ? (
          <Ch02Exchange
            commodities={commodities}
            owners={owners}
            onSharedChanged={onSharedChanged}
          />
        ) : activeChapterId === "ch03" ? (
          <Ch03Money owners={owners} onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch04" ? (
          <Ch04Capital onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch05" ? (
          <Ch05Contradictions onSharedChanged={onSharedChanged} />
        ) : null}
```

- [ ] **Step 2: Update registry.ts — mark ch05 as done**

Find:
```typescript
  { id: "ch05", number: 5,  title: "Contradictions in the General Formula",  ..., status: "pending" },
```

Change `status: "pending"` to `status: "done"`.

- [ ] **Step 3: Run TypeScript check and build**

```bash
cd web && npm run lint && npm run build
```

Expected: no type errors, build succeeds.

- [ ] **Step 4: Commit**

```bash
git add web/src/components/ChapterShell.tsx web/src/chapters/registry.ts
git commit -m "feat(web): wire Ch.05 panel into shell and mark done in registry"
```

---

## Task 8: Update architecture docs

**Files:**
- Modify: `docs/architecture.md`

- [ ] **Step 1: Update the roadmap table**

Find the Ch. 5 row in the roadmap table:
```
| Ch. 5-7   | Pending     | Labour-process, valorization, surplus-value  ...
```

Replace with two rows to separate Ch. 5 (now done) from Ch. 6-7 (still pending):
```
| Ch. 5     | ✅ Done     | Contradictions in the general formula; value conservation proof | agent-service |
| Ch. 6-7   | Pending     | Labour-process, valorization, surplus-value  | agent-service, simulation-eng |
```

- [ ] **Step 2: Add "Ch. 5 — what was built" section**

Append after the existing `### Ch. 4 — what was built` section:

```markdown
### Ch. 5 — what was built

`agent-service` extends Ch. 4 to prove the central contradiction of Capital Vol. I, Ch. 5 — that circulation alone cannot be the source of surplus-value:

- **Owner class.** `Owner` added to `Class` enum — a simple commodity owner who buys and sells equivalents but cannot originate surplus-value (cannot do M—C—M′ circuits).
- **LabourMinutes.** `labour_minutes int64` field added to `Agent`, tracking the abstract-labour magnitude the agent commands; persisted via migration `00004_ch05_labour_minutes.sql`.
- **ExchangeResult.** Bilateral exchange outcome with before/after values for both parties and `origin` tag. `TotalBefore() == TotalAfter()` is the invariant.
- **MerchantsCapital.** M-C-M' operating purely in circulation; `Origin()` returns `"redistribution"` if `SurplusValue != 0`.
- **UsurersCapital.** Degenerate M-M' circuit (no commodity); the source of surplus cannot be located within the circuit.
- **ExchangeEquivalents / ExchangeNonEquivalents / TotalValue.** Pure functions proving value conservation for any bilateral exchange.
- **POST /v1/circuits** — stateless circuit probe; returns `surplus_value` and `origin` tag.
- **POST /v1/exchange-simulations** — stateless exchange simulator; returns full `ExchangeResult` with conservation proof.

The React UI adds a "Ch. 05 — Contradictions in the General Formula" panel with a circuit probe form and an exchange simulation table.
```

- [ ] **Step 3: Commit**

```bash
git add docs/architecture.md
git commit -m "docs: update architecture.md for Ch.05 completion"
```

---

## Self-Review Checklist

**Spec coverage:**

- [x] `AgentKindOwner` (`agent.Owner`) — Task 2
- [x] `Agent` carries `LabourMinutes` — Task 2
- [x] `ExchangeResult` type — Task 1
- [x] `MerchantsCapital` with `SurplusValue()` and `Origin()` — Task 1
- [x] `UsurersCapital` with `SurplusValue()` — Task 1
- [x] `ExchangeEquivalents` — Task 1
- [x] `ExchangeNonEquivalents` — Task 1
- [x] `TotalValue` — Task 1
- [x] `RoleSeller` / `RoleBuyer` constants — Task 1
- [x] §1 fixture: wine=50/corn=50 → SurplusValue=0, TotalValue=100 — exchange_test.go
- [x] §2 fixture: worth-100 sold for 110 → seller+10, buyer-10, total=210 — exchange_test.go
- [x] §3 fixture: A=40, B=50, total 90, after swap A=50, B=40, total 90 — exchange_test.go
- [x] §4 fixture: property test for any (x,y) — exchange_test.go
- [x] §5 fixture: scaling invariant — exchange_test.go
- [x] §6 fixture: M-M' UsurersCapital SurplusValue=10 — exchange_test.go
- [x] `POST /v1/circuits` endpoint — Task 3
- [x] `POST /v1/exchange-simulations` endpoint (renamed from `/v1/exchanges` to avoid market-service conflict) — Tasks 3-4
- [x] `GET/POST /v1/agents` already exist (Ch. 04) — no new work needed
- [x] React panel with circuit form and exchange simulation — Tasks 6-7
- [x] `ch05` wired into shell and marked done — Task 7
- [x] Docs updated — Task 8

**Type consistency check:**
- `agent.Pence` used consistently throughout Go code (not `int64` directly)
- `ExchangeResult` fields (`ABefore`, `BBefore`, `AAfter`, `BAfter`, `Origin`) consistent between exchange.go and handler response
- JSON tags in Go (`a_before`, `b_before`, `a_after`, `b_after`, `origin`) match TypeScript `ExchangeSimulation` interface
- `CircuitProof` JSON fields (`m`, `commodity_id`, `m_prime`, `surplus_value`, `origin`) match handler `computeCircuitResponse`
- `computeCircuitRequest` and `computeExchangeRequest` match `ComputeCircuitInput` and `SimulateExchangeInput` in TypeScript
