# Chapter 7 — Labour-Process and Valorization Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement Ch. 7 domain types (LabourProcess, ValorizationProcess, MeansOfProduction, Product), store, HTTP endpoints, simulation-engine ProductionRun, api-gateway proxy, and React UI panel.

**Architecture:** New domain types live in `agent-service/internal/agent/labour_process.go`. A new `LabourProcessStore` interface extends the store layer with `CreateLabourProcess` and `GetLabourProcess`. The `Worker` struct gains `LabourPowerValueMinutes` (the daily reproduction cost, snapshotted into labour process records). The simulation-engine gets a `ProductionRun` type. React adds `Ch07LabourProcess.tsx`.

**Tech Stack:** Go 1.25, MySQL 8 (goose migration), React 18 + TypeScript, Vite.

---

## File Map

**Create:**
- `services/agent-service/internal/agent/labour_process.go` — domain types and pure functions
- `services/agent-service/internal/agent/labour_process_test.go` — Marx fixture tests
- `services/agent-service/internal/store/migrations/00006_ch07_labour_process.sql` — new tables + column
- `services/agent-service/internal/transport/httpapi/labour_process_handler.go` — HTTP handlers
- `web/src/chapters/Ch07LabourProcess.tsx` — React UI panel

**Modify:**
- `services/agent-service/internal/agent/labour_power.go` — add `LabourPowerValueMinutes` to `Worker`
- `services/agent-service/internal/store/store.go` — add `LabourProcessStore` interface
- `services/agent-service/internal/store/memory.go` — implement `LabourProcessStore` + worker field
- `services/agent-service/internal/store/mysql.go` — implement `LabourProcessStore` + worker field
- `services/agent-service/internal/transport/httpapi/handler.go` — add `LabourProcessStore` field + update `New`
- `services/agent-service/internal/transport/httpapi/labour_power_handler.go` — add `labour_power_value_minutes` to `createWorkerRequest`
- `services/agent-service/internal/transport/httpapi/routes.go` — add labour-process routes
- `services/agent-service/cmd/agent-service/main.go` — add `LabourProcessStore` to `agentStore` interface
- `services/simulation-engine/internal/engine/engine.go` — add `ProductionRun` type
- `services/api-gateway/cmd/api-gateway/main.go` — proxy `/v1/labour-processes/*`
- `web/src/types.ts` — add Ch. 7 types + extend `LabourWorker`
- `web/src/api.ts` — add Ch. 7 API calls
- `web/src/chapters/registry.ts` — set ch07 status to "done"
- `web/src/components/ChapterShell.tsx` — wire Ch07 component + add quote
- `docs/architecture.md` — update chapter table

---

## Task 1: Domain types and pure functions

**Files:**
- Create: `services/agent-service/internal/agent/labour_process.go`

- [ ] **Step 1: Write the failing test** (see Task 2 — write tests and implementation together in this task)

- [ ] **Step 2: Write `labour_process.go`**

```go
package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

// LabourProcessID is a 96-bit hex ID for LabourProcess records.
type LabourProcessID string

func (id LabourProcessID) IsZero() bool { return id == "" }

func NewLabourProcessID() LabourProcessID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return LabourProcessID(hex.EncodeToString(b))
}

// RawMaterial is an input commodity consumed in one production run [§1.c].
type RawMaterial struct {
	CommodityID string        `json:"commodity_id"`
	Quantity    int64         `json:"quantity"`
	SNLTPerUnit LabourMinutes `json:"snlt_per_unit"`
}

// Instrument is a tool or machine that transfers value to the product gradually [§1.c].
type Instrument struct {
	CommodityID string        `json:"commodity_id"`
	WearPerRun  LabourMinutes `json:"wear_per_run"`
}

// MeansOfProduction bundles raw materials and instruments consumed in one run [§1.c].
type MeansOfProduction struct {
	RawMaterials []RawMaterial `json:"raw_materials"`
	Instruments  []Instrument  `json:"instruments"`
}

// LabourProcess is one purposeful act of production [§1].
// NecessaryLabourMinutes, ProductKind, and ProductQuantity are snapshotted at
// run-time so GET /v1/labour-processes/{id} can reconstruct the full result.
type LabourProcess struct {
	ID                     LabourProcessID   `json:"id"`
	WorkerID               AgentID           `json:"worker_id"`
	CapitalistID           AgentID           `json:"capitalist_id"`
	Means                  MeansOfProduction `json:"means"`
	Duration               LabourMinutes     `json:"duration"`
	NecessaryLabourMinutes LabourMinutes     `json:"necessary_labour_minutes"`
	ProductKind            string            `json:"product_kind"`
	ProductQuantity        int64             `json:"product_quantity"`
	CreatedAt              time.Time         `json:"created_at"`
}

// Validate rejects zero-duration runs and nil required parties [§1.d].
func (lp LabourProcess) Validate() error {
	if lp.WorkerID.IsZero() {
		return errors.New("agent: worker_id is required")
	}
	if lp.CapitalistID.IsZero() {
		return errors.New("agent: capitalist_id is required")
	}
	if lp.Duration <= 0 {
		return errors.New("agent: duration must be positive")
	}
	for _, rm := range lp.Means.RawMaterials {
		if rm.CommodityID == "" {
			return errors.New("agent: raw_material commodity_id is required")
		}
		if rm.SNLTPerUnit < 0 {
			return errors.New("agent: raw_material snlt_per_unit cannot be negative")
		}
	}
	return nil
}

// Product is the output use-value of a LabourProcess [§1.b].
type Product struct {
	CommodityKind string        `json:"commodity_kind"`
	Quantity      int64         `json:"quantity"`
	TotalValue    LabourMinutes `json:"total_value"`
}

// ValorizationProcess wraps a LabourProcess and computes value-magnitude
// quantities from it [§2].
type ValorizationProcess struct {
	Process                LabourProcess `json:"process"`
	NecessaryLabourMinutes LabourMinutes `json:"necessary_labour_minutes"`
}

// NecessaryLabour returns the labour-time needed to reproduce the worker [§2.a].
func (vp ValorizationProcess) NecessaryLabour() LabourMinutes { return vp.NecessaryLabourMinutes }

// SurplusLabour returns the unpaid portion of the working day [§2.b].
func (vp ValorizationProcess) SurplusLabour() LabourMinutes {
	return SurplusLabour(vp.Process.Duration, vp.NecessaryLabourMinutes)
}

// SurplusValue returns the value produced above the value of labour-power [§2.f].
func (vp ValorizationProcess) SurplusValue() LabourMinutes { return vp.SurplusLabour() }

// ProductValue returns the total value of the product [§2.c].
func (vp ValorizationProcess) ProductValue() LabourMinutes {
	return TransferredValue(vp.Process.Means) + ValueAdded(vp.Process.Duration)
}

// TransferredValue sums the SNLT of all raw materials and instrument wear [§2.c].
// Constant capital transfers value to the product but creates no new value.
func TransferredValue(mp MeansOfProduction) LabourMinutes {
	var total LabourMinutes
	for _, rm := range mp.RawMaterials {
		total += rm.SNLTPerUnit * LabourMinutes(rm.Quantity)
	}
	for _, inst := range mp.Instruments {
		total += inst.WearPerRun
	}
	return total
}

// ValueAdded returns the new value created by living labour during duration [§2].
// For uniform skill level = 1 the reduction to abstract labour is the identity.
func ValueAdded(duration LabourMinutes) LabourMinutes { return duration }

// SurplusLabour is the unpaid portion of the working day: wd - nl [§2].
func SurplusLabour(workingDay, necessaryLabour LabourMinutes) LabourMinutes {
	return workingDay - necessaryLabour
}
```

---

## Task 2: Domain tests using Marx's fixtures

**Files:**
- Create: `services/agent-service/internal/agent/labour_process_test.go`

- [ ] **Step 1: Write tests**

```go
package agent_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// Shilling conversion used across §2 fixtures: 3 shillings = 6 hours = 360 minutes
// → 1 shilling = 120 LabourMinutes (established in §2.a).
const (
	necessaryLabour = agent.LabourMinutes(360) // §2.a: 6 hours = 3 shillings
	workingDay      = agent.LabourMinutes(720) // §2.b: 12-hour working day
)

// §1.d: zero-duration run must fail validation
func TestLabourProcess_Validate_ZeroDuration(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Duration:     0,
	}
	if err := lp.Validate(); err == nil {
		t.Error("zero-duration LabourProcess should fail validation")
	}
}

// §1.d: missing worker_id must fail validation
func TestLabourProcess_Validate_MissingWorker(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		CapitalistID: agent.NewAgentID(),
		Duration:     360,
	}
	if err := lp.Validate(); err == nil {
		t.Error("LabourProcess without worker_id should fail validation")
	}
}

// §1.d: missing capitalist_id must fail validation
func TestLabourProcess_Validate_MissingCapitalist(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		WorkerID: agent.NewAgentID(),
		Duration: 360,
	}
	if err := lp.Validate(); err == nil {
		t.Error("LabourProcess without capitalist_id should fail validation")
	}
}

// §1.c: valid process with all three elements passes validation
func TestLabourProcess_Validate_Valid(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Means: agent.MeansOfProduction{
			RawMaterials: []agent.RawMaterial{
				{CommodityID: "cotton", Quantity: 10, SNLTPerUnit: 120},
			},
			Instruments: []agent.Instrument{
				{CommodityID: "spindle", WearPerRun: 240},
			},
		},
		Duration: workingDay,
	}
	if err := lp.Validate(); err != nil {
		t.Errorf("valid LabourProcess should pass: %v", err)
	}
}

// §2.b: SurplusLabour(720, 360) == 360
func TestSurplusLabour_TwelveHourDay(t *testing.T) {
	t.Parallel()
	got := agent.SurplusLabour(workingDay, necessaryLabour)
	const want = agent.LabourMinutes(360)
	if got != want {
		t.Errorf("SurplusLabour(720, 360): want %d, got %d", want, got)
	}
}

// Invariant: SurplusLabour(wd, nl) == wd - nl
func TestSurplusLabour_Invariant(t *testing.T) {
	t.Parallel()
	cases := [][2]agent.LabourMinutes{{720, 360}, {480, 240}, {600, 300}, {360, 360}}
	for _, c := range cases {
		wd, nl := c[0], c[1]
		got := agent.SurplusLabour(wd, nl)
		want := wd - nl
		if got != want {
			t.Errorf("SurplusLabour(%d, %d): want %d, got %d", wd, nl, want, got)
		}
	}
}

// §2.c: TransferredValue = cotton(10*120) + spindle(240) = 1440
func TestTransferredValue_FifteenShillingsRun(t *testing.T) {
	t.Parallel()
	mp := agent.MeansOfProduction{
		RawMaterials: []agent.RawMaterial{{CommodityID: "cotton", Quantity: 10, SNLTPerUnit: 120}},
		Instruments:  []agent.Instrument{{CommodityID: "spindle", WearPerRun: 240}},
	}
	got := agent.TransferredValue(mp)
	const want = agent.LabourMinutes(1440) // 10*120 + 240
	if got != want {
		t.Errorf("TransferredValue: want %d, got %d", want, got)
	}
}

// Invariant: TransferredValue >= 0 for empty means
func TestTransferredValue_EmptyMeans(t *testing.T) {
	t.Parallel()
	if got := agent.TransferredValue(agent.MeansOfProduction{}); got < 0 {
		t.Errorf("TransferredValue of empty means should be >= 0, got %d", got)
	}
}

// §2.c: ProductValue = TransferredValue(1440) + ValueAdded(360) = 1800 (= 15 shillings * 120)
func TestProductValue_FifteenShillings(t *testing.T) {
	t.Parallel()
	mp := agent.MeansOfProduction{
		RawMaterials: []agent.RawMaterial{{CommodityID: "cotton", Quantity: 10, SNLTPerUnit: 120}},
		Instruments:  []agent.Instrument{{CommodityID: "spindle", WearPerRun: 240}},
	}
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Means:        mp,
		Duration:     necessaryLabour, // 6-hour day: only necessary labour (§2.c scenario)
	}
	vp := agent.ValorizationProcess{Process: lp, NecessaryLabourMinutes: necessaryLabour}
	got := vp.ProductValue()
	const want = agent.LabourMinutes(1800) // 1440 + 360
	if got != want {
		t.Errorf("ProductValue: want %d, got %d", want, got)
	}
}

// Invariant: ProductValue == TransferredValue + ValueAdded
func TestProductValue_Invariant(t *testing.T) {
	t.Parallel()
	mp := agent.MeansOfProduction{
		RawMaterials: []agent.RawMaterial{{CommodityID: "cotton", Quantity: 10, SNLTPerUnit: 120}},
		Instruments:  []agent.Instrument{{CommodityID: "spindle", WearPerRun: 240}},
	}
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Means:        mp,
		Duration:     workingDay,
	}
	vp := agent.ValorizationProcess{Process: lp, NecessaryLabourMinutes: necessaryLabour}
	want := agent.TransferredValue(mp) + agent.ValueAdded(workingDay)
	if vp.ProductValue() != want {
		t.Errorf("ProductValue invariant: want %d, got %d", want, vp.ProductValue())
	}
}

// Invariant: ValueAdded(d) == d for uniform skill level = 1
func TestValueAdded_Identity(t *testing.T) {
	t.Parallel()
	for _, d := range []agent.LabourMinutes{0, 60, 360, 720} {
		if got := agent.ValueAdded(d); got != d {
			t.Errorf("ValueAdded(%d): want %d, got %d", d, d, got)
		}
	}
}

// §2.e: ValueAdded(workingDay) > NecessaryLabour when WorkingDay > NecessaryLabour
func TestValueAdded_ExceedsNecessaryLabour(t *testing.T) {
	t.Parallel()
	added := agent.ValueAdded(workingDay)
	if added <= necessaryLabour {
		t.Errorf("ValueAdded(%d)=%d should exceed NecessaryLabour %d", workingDay, added, necessaryLabour)
	}
}

// Invariant: SurplusValue == 0 iff WorkingDay == NecessaryLabour
func TestValorizationProcess_NoSurplusWhenEqual(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Duration:     necessaryLabour,
	}
	vp := agent.ValorizationProcess{Process: lp, NecessaryLabourMinutes: necessaryLabour}
	if sv := vp.SurplusValue(); sv != 0 {
		t.Errorf("SurplusValue when wd==nl: want 0, got %d", sv)
	}
}

// §2.f: SurplusValue > 0 whenever WorkingDay > NecessaryLabour
func TestValorizationProcess_SurplusValuePositive(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Duration:     workingDay,
	}
	vp := agent.ValorizationProcess{Process: lp, NecessaryLabourMinutes: necessaryLabour}
	if vp.SurplusValue() <= 0 {
		t.Errorf("SurplusValue: want > 0, got %d", vp.SurplusValue())
	}
}

// Invariant: NecessaryLabour + SurplusLabour == WorkingDay
func TestValorizationProcess_WorkingDayPartition(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Duration:     workingDay,
	}
	vp := agent.ValorizationProcess{Process: lp, NecessaryLabourMinutes: necessaryLabour}
	if vp.NecessaryLabour()+vp.SurplusLabour() != workingDay {
		t.Errorf("WorkingDay partition invariant: %d + %d != %d",
			vp.NecessaryLabour(), vp.SurplusLabour(), workingDay)
	}
}

// NewLabourProcessID produces a 24-char hex string
func TestNewLabourProcessID(t *testing.T) {
	t.Parallel()
	id := agent.NewLabourProcessID()
	if id.IsZero() {
		t.Error("NewLabourProcessID should not be zero")
	}
	if len(string(id)) != 24 {
		t.Errorf("LabourProcessID: want 24 chars, got %d", len(string(id)))
	}
}
```

- [ ] **Step 2: Ask the user to run** `make vet test build` and confirm the new tests pass.

- [ ] **Step 3: Commit**

```bash
git add services/agent-service/internal/agent/labour_process.go \
        services/agent-service/internal/agent/labour_process_test.go
git commit -m "feat(ch07): add LabourProcess domain types and pure functions

Adds RawMaterial, Instrument, MeansOfProduction, LabourProcess,
Product, and ValorizationProcess types to the agent package.
Implements TransferredValue, ValueAdded, and SurplusLabour as pure
functions. All seven invariants from the Ch. 7 spec are covered by
tests using Marx's 15-shillings and 12-hour-working-day fixtures.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
```

---

## Task 3: Extend Worker with LabourPowerValueMinutes

**Files:**
- Modify: `services/agent-service/internal/agent/labour_power.go`
- Modify: `services/agent-service/internal/transport/httpapi/labour_power_handler.go`

- [ ] **Step 1: Add field to `Worker` struct in `labour_power.go`**

In `labour_power.go`, find the `Worker` struct and add the new field:

```go
// Worker is the agent who owns their labour-power and sells it as a commodity.
type Worker struct {
	LabourAgent
	OwnsLabourPower          bool          `json:"owns_labour_power"`
	OwnsCommoditiesToSell    bool          `json:"owns_commodities_to_sell"`
	LabourPower              LabourPower   `json:"labour_power"`
	LabourPowerValueMinutes  LabourMinutes `json:"labour_power_value_minutes"`
}
```

`Validate()` does not need to require `LabourPowerValueMinutes > 0` — zero is valid (means "not set"; Ch. 6 tests use zero and must keep passing).

- [ ] **Step 2: Add field to `createWorkerRequest` in `labour_power_handler.go`**

Find `createWorkerRequest` and add the new field:

```go
type createWorkerRequest struct {
	OwnsLabourPower         bool                `json:"owns_labour_power"`
	OwnsCommoditiesToSell   bool                `json:"owns_commodities_to_sell"`
	CapacityMinutesPerDay   agent.LabourMinutes `json:"capacity_minutes_per_day"`
	LabourPowerValueMinutes agent.LabourMinutes `json:"labour_power_value_minutes"`
}
```

- [ ] **Step 3: Wire the new field in `CreateWorker` handler**

In `CreateWorker` in `labour_power_handler.go`, find the block that builds `worker` and add the field:

```go
worker := agent.Worker{
	OwnsLabourPower:         req.OwnsLabourPower,
	OwnsCommoditiesToSell:   req.OwnsCommoditiesToSell,
	LabourPower:             agent.LabourPower{CapacityMinutesPerDay: req.CapacityMinutesPerDay},
	LabourPowerValueMinutes: req.LabourPowerValueMinutes,
}
```

- [ ] **Step 4: Ask the user to run** `make vet test build` — all existing Ch. 6 tests should still pass.

---

## Task 4: Store interface — add LabourProcessStore

**Files:**
- Modify: `services/agent-service/internal/store/store.go`

- [ ] **Step 1: Add `LabourProcessStore` interface**

Append to `store.go` after the existing `LabourPowerStore`:

```go
// LabourProcessStore is the persistence contract for Ch. 7 labour process records.
type LabourProcessStore interface {
	CreateLabourProcess(ctx context.Context, lp agent.LabourProcess) (agent.LabourProcess, error)
	GetLabourProcess(ctx context.Context, id agent.LabourProcessID) (agent.LabourProcess, error)
}
```

---

## Task 5: Migration — new tables and column

**Files:**
- Create: `services/agent-service/internal/store/migrations/00006_ch07_labour_process.sql`

- [ ] **Step 1: Write migration**

```sql
-- +goose Up
ALTER TABLE labour_workers
    ADD COLUMN labour_power_value_minutes BIGINT NOT NULL DEFAULT 0;

CREATE TABLE IF NOT EXISTS labour_processes (
    id                        VARCHAR(24)   NOT NULL PRIMARY KEY,
    worker_id                 VARCHAR(24)   NOT NULL,
    capitalist_id             VARCHAR(24)   NOT NULL,
    duration                  BIGINT        NOT NULL,
    necessary_labour_minutes  BIGINT        NOT NULL DEFAULT 0,
    means_json                JSON          NOT NULL,
    product_kind              VARCHAR(255)  NOT NULL DEFAULT '',
    product_quantity          BIGINT        NOT NULL DEFAULT 0,
    created_at                DATETIME(6)   NOT NULL,
    INDEX idx_lp_worker_id    (worker_id),
    INDEX idx_lp_capitalist   (capitalist_id)
);

-- +goose Down
DROP TABLE IF EXISTS labour_processes;
ALTER TABLE labour_workers DROP COLUMN labour_power_value_minutes;
```

---

## Task 6: Memory store — implement LabourProcessStore + Worker field

**Files:**
- Modify: `services/agent-service/internal/store/memory.go`

- [ ] **Step 1: Add `labourProcesses` map to `Memory` struct**

Find the `Memory` struct and add a new field:

```go
type Memory struct {
	mu                sync.RWMutex
	agents            map[agent.ID]agent.Agent
	circuits          map[agent.ID]agent.CapitalCircuit
	labourWorkers     map[agent.AgentID]agent.Worker
	labourCapitalists map[agent.AgentID]agent.Capitalist
	offerings         map[agent.AgentID]agent.LabourPowerOffering
	purchases         map[agent.PurchaseID]agent.LabourPowerPurchase
	labourProcesses   map[agent.LabourProcessID]agent.LabourProcess
	now               func() time.Time
}
```

- [ ] **Step 2: Initialize the new map in `NewMemory`**

```go
func NewMemory() *Memory {
	return &Memory{
		agents:            make(map[agent.ID]agent.Agent),
		circuits:          make(map[agent.ID]agent.CapitalCircuit),
		labourWorkers:     make(map[agent.AgentID]agent.Worker),
		labourCapitalists: make(map[agent.AgentID]agent.Capitalist),
		offerings:         make(map[agent.AgentID]agent.LabourPowerOffering),
		purchases:         make(map[agent.PurchaseID]agent.LabourPowerPurchase),
		labourProcesses:   make(map[agent.LabourProcessID]agent.LabourProcess),
		now:               time.Now,
	}
}
```

- [ ] **Step 3: Add `LabourPowerValueMinutes` to `CreateWorker`**

In the existing `CreateWorker` method, the worker is stored after assignment. Because Go copies the struct, `LabourPowerValueMinutes` is preserved automatically — no code change needed in `CreateWorker`. The field is already on the struct.

Verify `scanWorkerRow` (MySQL) reads it; for Memory the struct copy handles it. ✓

- [ ] **Step 4: Add `CreateLabourProcess` and `GetLabourProcess` methods**

Append to `memory.go`:

```go
func (m *Memory) CreateLabourProcess(_ context.Context, lp agent.LabourProcess) (agent.LabourProcess, error) {
	if err := lp.Validate(); err != nil {
		return agent.LabourProcess{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if lp.ID.IsZero() {
		lp.ID = agent.NewLabourProcessID()
	}
	lp.CreatedAt = m.now()
	m.labourProcesses[lp.ID] = lp
	return lp, nil
}

func (m *Memory) GetLabourProcess(_ context.Context, id agent.LabourProcessID) (agent.LabourProcess, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	lp, ok := m.labourProcesses[id]
	if !ok {
		return agent.LabourProcess{}, ErrNotFound
	}
	return lp, nil
}
```

- [ ] **Step 5: Write memory store tests**

Add a new test file `services/agent-service/internal/store/labour_process_memory_test.go`:

```go
package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

func makeLabourProcess() agent.LabourProcess {
	return agent.LabourProcess{
		WorkerID:               agent.NewAgentID(),
		CapitalistID:           agent.NewAgentID(),
		Duration:               720,
		NecessaryLabourMinutes: 360,
		Means: agent.MeansOfProduction{
			RawMaterials: []agent.RawMaterial{
				{CommodityID: "cotton", Quantity: 10, SNLTPerUnit: 120},
			},
			Instruments: []agent.Instrument{
				{CommodityID: "spindle", WearPerRun: 240},
			},
		},
		ProductKind:     "yarn",
		ProductQuantity: 10,
	}
}

func TestMemory_CreateGetLabourProcess(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	lp := makeLabourProcess()
	created, err := m.CreateLabourProcess(ctx, lp)
	if err != nil {
		t.Fatalf("CreateLabourProcess: %v", err)
	}
	if created.ID.IsZero() {
		t.Error("CreateLabourProcess should assign an ID")
	}
	got, err := m.GetLabourProcess(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetLabourProcess: %v", err)
	}
	if got.Duration != lp.Duration {
		t.Errorf("duration mismatch: want %d, got %d", lp.Duration, got.Duration)
	}
	if got.NecessaryLabourMinutes != lp.NecessaryLabourMinutes {
		t.Errorf("necessary_labour mismatch: want %d, got %d",
			lp.NecessaryLabourMinutes, got.NecessaryLabourMinutes)
	}
	if got.ProductKind != lp.ProductKind {
		t.Errorf("product_kind mismatch: want %q, got %q", lp.ProductKind, got.ProductKind)
	}
}

func TestMemory_GetLabourProcess_NotFound(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.GetLabourProcess(context.Background(), agent.NewLabourProcessID())
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestMemory_CreateLabourProcess_InvalidDuration(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	lp := makeLabourProcess()
	lp.Duration = 0
	_, err := m.CreateLabourProcess(context.Background(), lp)
	if err == nil {
		t.Error("zero-duration LabourProcess should fail validation")
	}
}
```

- [ ] **Step 6: Ask user to run** `make vet test build`.

---

## Task 7: MySQL store — implement LabourProcessStore + Worker field

**Files:**
- Modify: `services/agent-service/internal/store/mysql.go`

- [ ] **Step 1: Add `encoding/json` import** (if not already present)

Add `"encoding/json"` to the imports in `mysql.go`.

- [ ] **Step 2: Update `CreateWorker` to write `labour_power_value_minutes`**

Find the `CreateWorker` method. Change the INSERT query and args:

```go
func (m *MySQL) CreateWorker(ctx context.Context, w agent.Worker) (agent.Worker, error) {
	if err := w.Validate(); err != nil {
		return agent.Worker{}, err
	}
	if w.ID.IsZero() {
		w.ID = agent.NewAgentID()
	}
	now := m.now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now
	w.Kind = agent.AgentKindWorker
	const q = `INSERT INTO labour_workers
		(id, kind, owns_labour_power, owns_commodities_to_sell,
		 capacity_minutes_per_day, labour_power_value_minutes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(w.ID), string(w.Kind),
		w.OwnsLabourPower, w.OwnsCommoditiesToSell,
		int64(w.LabourPower.CapacityMinutesPerDay),
		int64(w.LabourPowerValueMinutes),
		w.CreatedAt, w.UpdatedAt,
	)
	if err != nil {
		return agent.Worker{}, err
	}
	return w, nil
}
```

- [ ] **Step 3: Update `GetWorker` query and `scanWorker` to read new column**

Update `GetWorker`:
```go
func (m *MySQL) GetWorker(ctx context.Context, id agent.AgentID) (agent.Worker, error) {
	const q = `SELECT id, kind, owns_labour_power, owns_commodities_to_sell,
		capacity_minutes_per_day, labour_power_value_minutes, created_at, updated_at
		FROM labour_workers WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanWorker(row)
}
```

Update `ListWorkers`:
```go
func (m *MySQL) ListWorkers(ctx context.Context) ([]agent.Worker, error) {
	const q = `SELECT id, kind, owns_labour_power, owns_commodities_to_sell,
		capacity_minutes_per_day, labour_power_value_minutes, created_at, updated_at
		FROM labour_workers ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []agent.Worker
	for rows.Next() {
		w, err := scanWorkerRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}
```

Update `scanWorker` and `scanWorkerRow` to scan the new column:

```go
func scanWorker(row *sql.Row) (agent.Worker, error) {
	var w agent.Worker
	var id, kind string
	var cap, lpv int64
	err := row.Scan(&id, &kind, &w.OwnsLabourPower, &w.OwnsCommoditiesToSell,
		&cap, &lpv, &w.CreatedAt, &w.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.Worker{}, ErrNotFound
	}
	if err != nil {
		return agent.Worker{}, err
	}
	w.ID = agent.AgentID(id)
	w.Kind = agent.AgentKind(kind)
	w.LabourPower.CapacityMinutesPerDay = agent.LabourMinutes(cap)
	w.LabourPowerValueMinutes = agent.LabourMinutes(lpv)
	return w, nil
}

func scanWorkerRow(rows *sql.Rows) (agent.Worker, error) {
	var w agent.Worker
	var id, kind string
	var cap, lpv int64
	if err := rows.Scan(&id, &kind, &w.OwnsLabourPower, &w.OwnsCommoditiesToSell,
		&cap, &lpv, &w.CreatedAt, &w.UpdatedAt); err != nil {
		return agent.Worker{}, err
	}
	w.ID = agent.AgentID(id)
	w.Kind = agent.AgentKind(kind)
	w.LabourPower.CapacityMinutesPerDay = agent.LabourMinutes(cap)
	w.LabourPowerValueMinutes = agent.LabourMinutes(lpv)
	return w, nil
}
```

- [ ] **Step 4: Add `CreateLabourProcess` and `GetLabourProcess` to MySQL**

Append to `mysql.go`:

```go
func (m *MySQL) CreateLabourProcess(ctx context.Context, lp agent.LabourProcess) (agent.LabourProcess, error) {
	if err := lp.Validate(); err != nil {
		return agent.LabourProcess{}, err
	}
	if lp.ID.IsZero() {
		lp.ID = agent.NewLabourProcessID()
	}
	lp.CreatedAt = m.now().UTC()
	meansJSON, err := json.Marshal(lp.Means)
	if err != nil {
		return agent.LabourProcess{}, err
	}
	const q = `INSERT INTO labour_processes
		(id, worker_id, capitalist_id, duration, necessary_labour_minutes,
		 means_json, product_kind, product_quantity, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = m.db.ExecContext(ctx, q,
		string(lp.ID), string(lp.WorkerID), string(lp.CapitalistID),
		int64(lp.Duration), int64(lp.NecessaryLabourMinutes),
		string(meansJSON), lp.ProductKind, lp.ProductQuantity,
		lp.CreatedAt,
	)
	if err != nil {
		return agent.LabourProcess{}, err
	}
	return lp, nil
}

func (m *MySQL) GetLabourProcess(ctx context.Context, id agent.LabourProcessID) (agent.LabourProcess, error) {
	const q = `SELECT id, worker_id, capitalist_id, duration,
		necessary_labour_minutes, means_json, product_kind, product_quantity, created_at
		FROM labour_processes WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanLabourProcess(row)
}

func scanLabourProcess(row *sql.Row) (agent.LabourProcess, error) {
	var lp agent.LabourProcess
	var id, workerID, capitalistID, meansJSON string
	var dur, nl int64
	err := row.Scan(&id, &workerID, &capitalistID, &dur, &nl,
		&meansJSON, &lp.ProductKind, &lp.ProductQuantity, &lp.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.LabourProcess{}, ErrNotFound
	}
	if err != nil {
		return agent.LabourProcess{}, err
	}
	lp.ID = agent.LabourProcessID(id)
	lp.WorkerID = agent.AgentID(workerID)
	lp.CapitalistID = agent.AgentID(capitalistID)
	lp.Duration = agent.LabourMinutes(dur)
	lp.NecessaryLabourMinutes = agent.LabourMinutes(nl)
	if err := json.Unmarshal([]byte(meansJSON), &lp.Means); err != nil {
		return agent.LabourProcess{}, err
	}
	return lp, nil
}
```

- [ ] **Step 5: Ask user to run** `make vet test build`.

---

## Task 8: HTTP handler — labour process endpoint

**Files:**
- Create: `services/agent-service/internal/transport/httpapi/labour_process_handler.go`
- Modify: `services/agent-service/internal/transport/httpapi/handler.go`
- Modify: `services/agent-service/internal/transport/httpapi/routes.go`
- Modify: `services/agent-service/cmd/agent-service/main.go`

- [ ] **Step 1: Add `LabourProcessStore` to `Handler` in `handler.go`**

Find the `Handler` struct and `New` function and update them:

```go
type Handler struct {
	Store              store.Store
	CircuitStore       store.CircuitStore
	LabourPowerStore   store.LabourPowerStore
	LabourProcessStore store.LabourProcessStore
	Logger             *slog.Logger
}

func New(s store.Store, cs store.CircuitStore, lps store.LabourPowerStore, lproc store.LabourProcessStore, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{Store: s, CircuitStore: cs, LabourPowerStore: lps, LabourProcessStore: lproc, Logger: logger}
}
```

- [ ] **Step 2: Create `labour_process_handler.go`**

```go
package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

type runLabourProcessRequest struct {
	WorkerID          agent.AgentID           `json:"worker_id"`
	CapitalistID      agent.AgentID           `json:"capitalist_id"`
	MeansOfProduction agent.MeansOfProduction `json:"means_of_production"`
	DurationMinutes   agent.LabourMinutes     `json:"duration_minutes"`
	ProductKind       string                  `json:"product_kind"`
	ProductQuantity   int64                   `json:"product_quantity"`
}

type valorizationSummary struct {
	NecessaryLabour agent.LabourMinutes `json:"necessary_labour"`
	SurplusLabour   agent.LabourMinutes `json:"surplus_labour"`
	SurplusValue    agent.LabourMinutes `json:"surplus_value"`
	ProductValue    agent.LabourMinutes `json:"product_value"`
}

type labourProcessResponse struct {
	LabourProcess agent.LabourProcess `json:"labour_process"`
	Product       agent.Product       `json:"product"`
	Valorization  valorizationSummary `json:"valorization"`
}

func buildLabourProcessResponse(lp agent.LabourProcess) labourProcessResponse {
	vp := agent.ValorizationProcess{
		Process:                lp,
		NecessaryLabourMinutes: lp.NecessaryLabourMinutes,
	}
	product := agent.Product{
		CommodityKind: lp.ProductKind,
		Quantity:      lp.ProductQuantity,
		TotalValue:    vp.ProductValue(),
	}
	return labourProcessResponse{
		LabourProcess: lp,
		Product:       product,
		Valorization: valorizationSummary{
			NecessaryLabour: vp.NecessaryLabour(),
			SurplusLabour:   vp.SurplusLabour(),
			SurplusValue:    vp.SurplusValue(),
			ProductValue:    vp.ProductValue(),
		},
	}
}

func (h *Handler) RunLabourProcess(w http.ResponseWriter, r *http.Request) {
	var req runLabourProcessRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	worker, err := h.LabourPowerStore.GetWorker(r.Context(), req.WorkerID)
	if err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusBadRequest, "worker not found")
			return
		}
		writeAppError(w, err)
		return
	}

	if _, err := h.LabourPowerStore.GetCapitalist(r.Context(), req.CapitalistID); err != nil {
		if isNotFound(err) {
			writeError(w, http.StatusBadRequest, "capitalist not found")
			return
		}
		writeAppError(w, err)
		return
	}

	lp := agent.LabourProcess{
		WorkerID:               req.WorkerID,
		CapitalistID:           req.CapitalistID,
		Means:                  req.MeansOfProduction,
		Duration:               req.DurationMinutes,
		NecessaryLabourMinutes: worker.LabourPowerValueMinutes,
		ProductKind:            req.ProductKind,
		ProductQuantity:        req.ProductQuantity,
	}
	if err := lp.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	saved, err := h.LabourProcessStore.CreateLabourProcess(r.Context(), lp)
	if err != nil {
		writeAppError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, buildLabourProcessResponse(saved))
}

func (h *Handler) GetLabourProcessRecord(w http.ResponseWriter, r *http.Request) {
	id := agent.LabourProcessID(r.PathValue("id"))
	lp, err := h.LabourProcessStore.GetLabourProcess(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildLabourProcessResponse(lp))
}

func isNotFound(err error) bool {
	return err != nil && (err.Error() == store.ErrNotFound.Error() ||
		errors.Is(err, store.ErrNotFound))
}
```

Wait — `errors` needs to be imported. Add it to the imports:

```go
import (
	"errors"
	"net/http"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)
```

Also note: `isNotFound` uses `errors.Is` directly, so the helper is:

```go
func isNotFound(err error) bool {
	return errors.Is(err, store.ErrNotFound)
}
```

- [ ] **Step 3: Add routes in `routes.go`**

Append to `Register`:

```go
	// Ch. 7 — The Labour-Process and the Production of Surplus-Value
	s.HandleFunc("POST /v1/labour-processes", h.RunLabourProcess)
	s.HandleFunc("GET /v1/labour-processes/{id}", h.GetLabourProcessRecord)
```

- [ ] **Step 4: Update `agentStore` interface in `main.go`**

Find the `agentStore` interface and add the new store:

```go
type agentStore interface {
	store.Store
	store.CircuitStore
	store.LabourPowerStore
	store.LabourProcessStore
}
```

- [ ] **Step 5: Update `httpapi.New` call in `main.go`**

Find the call `httpapi.Register(srv, httpapi.New(st, st, st, logger))` and update:

```go
httpapi.Register(srv, httpapi.New(st, st, st, st, logger))
```

- [ ] **Step 6: Ask user to run** `make vet test build`.

- [ ] **Step 7: Commit**

```bash
git add services/agent-service/internal/agent/labour_power.go \
        services/agent-service/internal/store/store.go \
        services/agent-service/internal/store/memory.go \
        services/agent-service/internal/store/mysql.go \
        services/agent-service/internal/store/migrations/00006_ch07_labour_process.sql \
        services/agent-service/internal/store/labour_process_memory_test.go \
        services/agent-service/internal/transport/httpapi/handler.go \
        services/agent-service/internal/transport/httpapi/labour_power_handler.go \
        services/agent-service/internal/transport/httpapi/labour_process_handler.go \
        services/agent-service/internal/transport/httpapi/routes.go \
        services/agent-service/cmd/agent-service/main.go
git commit -m "feat(ch07): add LabourProcess store, HTTP endpoints, and Worker field

Extends agent-service with Ch. 7 production machinery:
- Worker gains labour_power_value_minutes (the daily reproduction cost
  snapshotted into each ValorizationProcess at run-time).
- LabourProcessStore interface with CreateLabourProcess / GetLabourProcess.
- Memory and MySQL implementations; migration 00006 adds the
  labour_processes table and the new column to labour_workers.
- POST /v1/labour-processes runs a labour process; GET /v1/labour-processes/{id}
  fetches it. Both responses include the full valorization summary.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
```

---

## Task 9: simulation-engine — add ProductionRun type

**Files:**
- Modify: `services/simulation-engine/internal/engine/engine.go`

- [ ] **Step 1: Add `ProductionRun` to `engine.go`**

Replace the file contents (currently just the package declaration) with:

```go
// Package engine holds the time-step orchestrator that advances the
// simulated economy one period at a time. The full tick scheduler is
// deferred to Ch. 10+; this package currently exposes domain types
// that later chapters will wire into the scheduler.
package engine

// ProductionRun records a ValorizationProcess result for a given simulation
// tick. LabourProcessID references a record in agent-service. Introduced in
// Ch. 7; the tick loop and persistence are added in Ch. 10+.
type ProductionRun struct {
	ID              string `json:"id"`
	TickID          string `json:"tick_id"`
	LabourProcessID string `json:"labour_process_id"`
	SurplusValue    int64  `json:"surplus_value"` // LabourMinutes
}
```

- [ ] **Step 2: Ask user to run** `make vet test build`.

---

## Task 10: api-gateway — proxy labour-processes

**Files:**
- Modify: `services/api-gateway/cmd/api-gateway/main.go`

- [ ] **Step 1: Add proxy route for `/v1/labour-processes`**

In `main.go`, after the existing Ch. 6 proxy block (the `srv.Handle("/v1/labour-power/...")` lines), add:

```go
	// Ch. 7 — labour-process routes proxy to agent-service
	srv.Handle("/v1/labour-processes", agentProxy)
	srv.Handle("/v1/labour-processes/{rest...}", agentProxy)
```

- [ ] **Step 2: Update `handleInfo` to reflect Ch. 7 status**

Change `"status": "ch-6-labour-power"` to `"status": "ch-7-labour-process"` and update the `chapter` string to `"Capital Vol. I, Ch. 7 - The Labour-Process and the Production of Surplus-Value"`.

- [ ] **Step 3: Ask user to run** `make vet test build`.

---

## Task 11: React — types.ts

**Files:**
- Modify: `web/src/types.ts`

- [ ] **Step 1: Extend `LabourWorker` with the new field**

Find the `LabourWorker` interface and add `labour_power_value_minutes`:

```typescript
export interface LabourWorker {
  id: string;
  kind: "worker";
  owns_labour_power: boolean;
  owns_commodities_to_sell: boolean;
  labour_power: LabourPower;
  labour_power_value_minutes: number; // LabourMinutes; daily reproduction cost
  created_at: string;
  updated_at: string;
}
```

Also extend `CreateLabourWorkerInput`:

```typescript
export interface CreateLabourWorkerInput {
  owns_labour_power: boolean;
  owns_commodities_to_sell: boolean;
  capacity_minutes_per_day: number;
  labour_power_value_minutes?: number; // optional; defaults to 0 on the server
}
```

- [ ] **Step 2: Add Ch. 7 types at the end of `types.ts`**

```typescript
// --- agent-service types (Ch. 7: The Labour-Process and Valorization) --------

export interface RawMaterial {
  commodity_id: string;
  quantity: number;
  snlt_per_unit: number; // LabourMinutes per unit
}

export interface Instrument {
  commodity_id: string;
  wear_per_run: number; // LabourMinutes transferred per run
}

export interface MeansOfProduction {
  raw_materials: RawMaterial[];
  instruments: Instrument[];
}

export interface LabourProcess {
  id: string;
  worker_id: string;
  capitalist_id: string;
  means: MeansOfProduction;
  duration: number; // LabourMinutes
  necessary_labour_minutes: number;
  product_kind: string;
  product_quantity: number;
  created_at: string;
}

export interface Product {
  commodity_kind: string;
  quantity: number;
  total_value: number; // LabourMinutes
}

export interface ValorizationSummary {
  necessary_labour: number;
  surplus_labour: number;
  surplus_value: number;
  product_value: number;
}

export interface RunLabourProcessResult {
  labour_process: LabourProcess;
  product: Product;
  valorization: ValorizationSummary;
}

export interface RunLabourProcessInput {
  worker_id: string;
  capitalist_id: string;
  means_of_production: MeansOfProduction;
  duration_minutes: number;
  product_kind: string;
  product_quantity: number;
}
```

---

## Task 12: React — api.ts

**Files:**
- Modify: `web/src/api.ts`

- [ ] **Step 1: Add Ch. 7 imports at top of `api.ts`**

In the import list from `"./types"`, add:

```typescript
  LabourProcess,
  RunLabourProcessInput,
  RunLabourProcessResult,
```

- [ ] **Step 2: Add Ch. 7 API functions**

Append to the `api` object, after the last Ch. 6 function (`getLabourPurchase`):

```typescript
  // --- agent-service (Ch. 7: Labour-Process) ---

  runLabourProcess: (input: RunLabourProcessInput) =>
    http<RunLabourProcessResult>("/v1/labour-processes", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getLabourProcess: (id: string) =>
    http<RunLabourProcessResult>(`/v1/labour-processes/${id}`),
```

---

## Task 13: React — Ch07LabourProcess.tsx

**Files:**
- Create: `web/src/chapters/Ch07LabourProcess.tsx`

- [ ] **Step 1: Write the component**

```tsx
import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type {
  LabourWorker,
  LabourCapitalist,
  RunLabourProcessResult,
  RawMaterial,
  Instrument,
} from "../types";

interface Ch07Props {
  onSharedChanged: () => void;
}

function minutesToHours(m: number): string {
  const h = Math.floor(m / 60);
  const min = m % 60;
  return min === 0 ? `${h}h` : `${h}h ${min}m`;
}

export function Ch07LabourProcess({ onSharedChanged: _onSharedChanged }: Ch07Props) {
  const [workers, setWorkers] = useState<LabourWorker[]>([]);
  const [capitalists, setCapitalists] = useState<LabourCapitalist[]>([]);
  const [result, setResult] = useState<RunLabourProcessResult | null>(null);
  const [loadErr, setLoadErr] = useState<string | null>(null);

  async function refresh() {
    try {
      const [ws, cs] = await Promise.all([
        api.listLabourWorkers(),
        api.listLabourCapitalists(),
      ]);
      setWorkers(ws);
      setCapitalists(cs);
    } catch (e) {
      setLoadErr(e instanceof Error ? e.message : String(e));
    }
  }

  useEffect(() => { refresh(); }, []);

  return (
    <>
      {loadErr && <p className="error">{loadErr}</p>}
      <RunLabourProcessPanel
        workers={workers}
        capitalists={capitalists}
        onResult={setResult}
      />
      {result && <ValorizationResultPanel result={result} />}
    </>
  );
}

interface RunPanelProps {
  workers: LabourWorker[];
  capitalists: LabourCapitalist[];
  onResult: (r: RunLabourProcessResult) => void;
}

function RunLabourProcessPanel({ workers, capitalists, onResult }: RunPanelProps) {
  const [workerID, setWorkerID] = useState("");
  const [capitalistID, setCapitalistID] = useState("");
  const [duration, setDuration] = useState(720);
  const [productKind, setProductKind] = useState("yarn");
  const [productQty, setProductQty] = useState(10);
  const [rawMaterials, setRawMaterials] = useState<RawMaterial[]>([
    { commodity_id: "cotton", quantity: 10, snlt_per_unit: 120 },
  ]);
  const [instruments, setInstruments] = useState<Instrument[]>([
    { commodity_id: "spindle", wear_per_run: 240 },
  ]);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.runLabourProcess({
        worker_id: workerID,
        capitalist_id: capitalistID,
        means_of_production: { raw_materials: rawMaterials, instruments },
        duration_minutes: duration,
        product_kind: productKind,
        product_quantity: productQty,
      });
      onResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  function updateRawMaterial(i: number, field: keyof RawMaterial, value: string | number) {
    setRawMaterials((prev) =>
      prev.map((rm, idx) => (idx === i ? { ...rm, [field]: value } : rm))
    );
  }

  function updateInstrument(i: number, field: keyof Instrument, value: string | number) {
    setInstruments((prev) =>
      prev.map((inst, idx) => (idx === i ? { ...inst, [field]: value } : inst))
    );
  }

  return (
    <section className="panel">
      <h2>Run Labour Process</h2>
      <form onSubmit={submit}>
        <label>
          Worker
          <select value={workerID} onChange={(e) => setWorkerID(e.target.value)} required>
            <option value="">— select —</option>
            {workers.map((w) => (
              <option key={w.id} value={w.id}>
                {w.id.slice(0, 8)}… (capacity {minutesToHours(w.labour_power.capacity_minutes_per_day)},
                repro cost {minutesToHours(w.labour_power_value_minutes)})
              </option>
            ))}
          </select>
        </label>
        <label>
          Capitalist
          <select value={capitalistID} onChange={(e) => setCapitalistID(e.target.value)} required>
            <option value="">— select —</option>
            {capitalists.map((c) => (
              <option key={c.id} value={c.id}>
                {c.id.slice(0, 8)}…
              </option>
            ))}
          </select>
        </label>
        <label>
          Working Day Duration (minutes)
          <input
            type="number"
            min={1}
            value={duration}
            onChange={(e) => setDuration(Number(e.target.value))}
          />
        </label>

        <fieldset>
          <legend>Raw Materials</legend>
          {rawMaterials.map((rm, i) => (
            <div key={i} className="input-row">
              <input
                placeholder="commodity_id"
                value={rm.commodity_id}
                onChange={(e) => updateRawMaterial(i, "commodity_id", e.target.value)}
              />
              <input
                type="number"
                placeholder="qty"
                value={rm.quantity}
                onChange={(e) => updateRawMaterial(i, "quantity", Number(e.target.value))}
              />
              <input
                type="number"
                placeholder="snlt_per_unit (min)"
                value={rm.snlt_per_unit}
                onChange={(e) => updateRawMaterial(i, "snlt_per_unit", Number(e.target.value))}
              />
              <button
                type="button"
                onClick={() => setRawMaterials((prev) => prev.filter((_, idx) => idx !== i))}
              >
                ×
              </button>
            </div>
          ))}
          <button
            type="button"
            onClick={() =>
              setRawMaterials((prev) => [
                ...prev,
                { commodity_id: "", quantity: 1, snlt_per_unit: 0 },
              ])
            }
          >
            + Add Raw Material
          </button>
        </fieldset>

        <fieldset>
          <legend>Instruments of Labour</legend>
          {instruments.map((inst, i) => (
            <div key={i} className="input-row">
              <input
                placeholder="commodity_id"
                value={inst.commodity_id}
                onChange={(e) => updateInstrument(i, "commodity_id", e.target.value)}
              />
              <input
                type="number"
                placeholder="wear_per_run (min)"
                value={inst.wear_per_run}
                onChange={(e) => updateInstrument(i, "wear_per_run", Number(e.target.value))}
              />
              <button
                type="button"
                onClick={() => setInstruments((prev) => prev.filter((_, idx) => idx !== i))}
              >
                ×
              </button>
            </div>
          ))}
          <button
            type="button"
            onClick={() =>
              setInstruments((prev) => [...prev, { commodity_id: "", wear_per_run: 0 }])
            }
          >
            + Add Instrument
          </button>
        </fieldset>

        <label>
          Product Kind
          <input
            value={productKind}
            onChange={(e) => setProductKind(e.target.value)}
            placeholder="yarn"
          />
        </label>
        <label>
          Product Quantity
          <input
            type="number"
            min={1}
            value={productQty}
            onChange={(e) => setProductQty(Number(e.target.value))}
          />
        </label>

        {err && <p className="error">{err}</p>}
        <button type="submit">Run Labour Process</button>
      </form>
    </section>
  );
}

function ValorizationResultPanel({ result }: { result: RunLabourProcessResult }) {
  const { labour_process: lp, product, valorization: v } = result;
  const surplusRate =
    v.necessary_labour > 0
      ? ((v.surplus_value / v.necessary_labour) * 100).toFixed(1)
      : "∞";

  return (
    <section className="panel">
      <h2>Valorization Result</h2>
      <p>
        <strong>Process ID:</strong> <code>{lp.id}</code>
      </p>
      <table>
        <tbody>
          <tr>
            <td>Working Day</td>
            <td>{minutesToHours(lp.duration)}</td>
          </tr>
          <tr>
            <td>Necessary Labour</td>
            <td>{minutesToHours(v.necessary_labour)}</td>
          </tr>
          <tr>
            <td>Surplus Labour</td>
            <td>{minutesToHours(v.surplus_labour)}</td>
          </tr>
          <tr>
            <td>Transferred Value (constant capital)</td>
            <td>
              {minutesToHours(
                v.product_value - lp.duration
              )}
            </td>
          </tr>
          <tr>
            <td>Value Added (living labour)</td>
            <td>{minutesToHours(lp.duration)}</td>
          </tr>
          <tr>
            <td>Total Product Value</td>
            <td>
              <strong>{minutesToHours(v.product_value)}</strong>
            </td>
          </tr>
          <tr>
            <td>Surplus Value</td>
            <td>
              <strong>{minutesToHours(v.surplus_value)}</strong>
            </td>
          </tr>
          <tr>
            <td>Rate of Surplus Value (s/v)</td>
            <td>{surplusRate}%</td>
          </tr>
        </tbody>
      </table>
      <p>
        <strong>Product:</strong> {lp.product_quantity} × {product.commodity_kind} —
        total value {minutesToHours(product.total_value)}
      </p>
      <blockquote>
        "The secret of the self-expansion of capital resolves itself into having
        the disposal of a definite quantity of other people's unpaid labour."
      </blockquote>
    </section>
  );
}
```

---

## Task 14: React — wire Ch07 into shell and registry

**Files:**
- Modify: `web/src/chapters/registry.ts`
- Modify: `web/src/components/ChapterShell.tsx`

- [ ] **Step 1: Update `registry.ts` — set ch07 status to "done"**

Find the ch07 entry and change `status: "pending"` to `status: "done"`.

- [ ] **Step 2: Update `ChapterShell.tsx` — add import, quote, and render branch**

Add the import at the top:

```typescript
import { Ch07LabourProcess } from "../chapters/Ch07LabourProcess";
```

Add to `QUOTES`:

```typescript
  ch07: "The secret of the self-expansion of capital resolves itself into having the disposal of a definite quantity of other people's unpaid labour.",
```

Add a render branch after the ch06 branch:

```tsx
        ) : activeChapterId === "ch07" ? (
          <Ch07LabourProcess onSharedChanged={onSharedChanged} />
```

- [ ] **Step 3: Ask user to run** `cd web && npm run lint && npm run build` and verify zero errors.

---

## Task 15: Update docs/architecture.md

**Files:**
- Modify: `docs/architecture.md`

- [ ] **Step 1: Update chapter table**

Find the Ch. 7 row and change `Pending` to `✅ Done`:

```
| Ch. 7     | ✅ Done     | Labour-process, valorization, surplus-value production   | agent-service, simulation-eng |
```

- [ ] **Step 2: Add a "Ch. 7 — what was built" section** after the Ch. 6 section

```markdown
### Ch. 7 — what was built

`agent-service` and `simulation-engine` model Capital Vol. I, Ch. 7 — the labour-process as the unity of the labour-process proper and the valorization process:

- **LabourProcess.** `LabourProcess`, `LabourProcessID`, `NewLabourProcessID` — one purposeful act of production tying `WorkerID`, `CapitalistID`, `MeansOfProduction`, and `Duration` together. `Validate()` enforces §1.d: zero-duration runs and missing parties are rejected.
- **MeansOfProduction.** `RawMaterial` (commodity reference + quantity + SNLT per unit) and `Instrument` (commodity reference + wear per run) — Marx's three factors of §1.c.
- **ValorizationProcess.** Wraps a `LabourProcess`; exposes `NecessaryLabour()`, `SurplusLabour()`, `SurplusValue()`, and `ProductValue()`. All seven invariants from the spec are test-covered.
- **Pure functions.** `TransferredValue(MeansOfProduction)` (constant capital), `ValueAdded(duration)` (living labour, identity for uniform skill), `SurplusLabour(wd, nl) = wd - nl`.
- **Worker extension.** `Worker` gains `LabourPowerValueMinutes` — the daily reproduction cost of labour-power, snapshotted into each `LabourProcess` record at run-time.
- **Store.** `LabourProcessStore` interface; `Memory` and `MySQL` implementations. Migration `00006` adds `labour_processes` table (means stored as JSON) and `labour_power_value_minutes` column on `labour_workers`.
- **HTTP.** `POST /v1/labour-processes` (run a process; returns product + full valorization summary), `GET /v1/labour-processes/{id}` (fetch a recorded run). Proxied through api-gateway.
- **ProductionRun.** `simulation-engine/engine.ProductionRun` type introduced; full tick scheduler deferred to Ch. 10+.
- **React UI.** Ch. 07 panel: worker/capitalist picker, means-of-production builder (raw materials + instruments), working-day duration input, valorization result card (necessary / surplus labour breakdown, rate of surplus value, product total value).
```

---

## Task 16: Final integration commit

- [ ] **Step 1: Ask user to run the full check**

```bash
make vet test build && cd web && npm run lint && npm run build
```

- [ ] **Step 2: Commit remaining files**

```bash
git add services/simulation-engine/internal/engine/engine.go \
        services/api-gateway/cmd/api-gateway/main.go \
        web/src/types.ts \
        web/src/api.ts \
        web/src/chapters/Ch07LabourProcess.tsx \
        web/src/chapters/registry.ts \
        web/src/components/ChapterShell.tsx \
        docs/architecture.md
git commit -m "feat(ch07): wire simulation-engine ProductionRun, gateway proxy, and React UI

Completes Ch. 7 integration across all layers:
- simulation-engine/engine.ProductionRun type (tick loop deferred to Ch. 10+).
- api-gateway proxies /v1/labour-processes/* to agent-service.
- types.ts + api.ts extended with RawMaterial, Instrument, MeansOfProduction,
  LabourProcess, Product, ValorizationSummary, and RunLabourProcessResult.
- Ch07LabourProcess.tsx: worker/capitalist selector, means-of-production builder,
  valorization result card with necessary/surplus labour breakdown and
  rate-of-surplus-value display.
- registry.ts marks ch07 done; ChapterShell.tsx wires the component.
- docs/architecture.md records Ch. 7 as complete.

Co-Authored-By: Claude Sonnet 4.6 <noreply@anthropic.com>"
```

---

## Self-Review

**Spec coverage check:**

| Spec requirement | Task |
|---|---|
| `LabourProcessID`, `NewLabourProcessID` | Task 1 |
| `RawMaterial`, `Instrument`, `MeansOfProduction` | Task 1 |
| `LabourProcess` with `WorkerID`, `CapitalistID`, `Means`, `Duration` | Task 1 |
| `LabourProcess.Validate()` rejects zero-duration and nil parties [§1.d] | Task 1 |
| `Product` with `CommodityKind`, `Quantity`, `TotalValue` | Task 1 |
| `ValorizationProcess` with `NecessaryLabour()`, `SurplusLabour()`, `SurplusValue()`, `ProductValue()` | Task 1 |
| `TransferredValue`, `ValueAdded`, `SurplusLabour` functions | Task 1 |
| Worker carries `LabourPowerValue` (as `LabourPowerValueMinutes`) | Task 3 |
| Store: `CreateWorker`, `GetWorker`, `ListWorkers`, `CreateCapitalist`, `GetCapitalist` | Already in Ch. 6 |
| Store: `CreateLabourProcess`, `GetLabourProcess` | Tasks 4, 6, 7 |
| Migration 00006 | Task 5 |
| `POST /v1/labour-processes`, `GET /v1/labour-processes/{id}` | Task 8 |
| api-gateway proxy `/v1/labour-processes/*` | Task 10 |
| `ProductionRun` in simulation-engine | Task 9 |
| React types (RawMaterial, Instrument, MeansOfProduction, LabourProcess, Product, ValorizationSummary) | Task 11 |
| React api (runLabourProcess, getLabourProcess) | Task 12 |
| React Ch07 UI panel | Task 13 |
| Invariants tested (all seven from spec) | Task 2 |

**Placeholder scan:** No TODOs or TBDs. All code blocks are complete.

**Type consistency check:**
- `agent.LabourProcessID` used consistently in `labour_process.go`, store interface, memory.go, mysql.go, handler.
- `agent.MeansOfProduction` in the request struct matches the domain type exactly.
- `valorizationSummary` field names (`necessary_labour`, `surplus_labour`, `surplus_value`, `product_value`) match `ValorizationSummary` in TypeScript.
- `RunLabourProcessResult` in TypeScript matches `labourProcessResponse` in Go.
- `labour_power_value_minutes` added to both Go `Worker` struct and TypeScript `LabourWorker` interface.
