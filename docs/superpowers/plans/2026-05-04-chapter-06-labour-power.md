# Chapter 06 — The Buying and Selling of Labour-Power: Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement domain types, persistence, HTTP endpoints, and React UI for Capital Vol. I Ch. 6 — the buying and selling of labour-power as a commodity.

**Architecture:** All new types live in the existing `agent` package alongside Ch. 4/5 types. The `Memory` and `MySQL` store impls gain a new `LabourPowerStore` interface. Eight new HTTP endpoints and a `Ch06LabourPower` React panel complete the chapter. Two existing `Class` constants (`Worker`, `Capitalist`) conflict with the spec's new type names and must be renamed first.

**Tech Stack:** Go 1.25, `database/sql` + goose migrations, React 18 + TypeScript, existing `pkg/httpx` and `pkg/mysql` patterns.

---

## Naming constraint

The existing `agent` package has:
```go
const Worker     Class = "worker"
const Capitalist Class = "capitalist"
```
The spec wants `type Worker struct` and `type Capitalist struct`. Go does not allow a constant and a type to share an identifier. **Task 1** renames the constants; Tasks 2+ use the freed names. The JSON wire value (`"worker"`, `"capitalist"`) is unchanged — only the Go identifiers change.

The spec's `Agent` base struct conflicts with the existing `agent.Agent`. The base is implemented as `LabourAgent` (embedded by `Worker` and `Capitalist`). All other spec identifiers are used verbatim.

---

## File map

| Action  | Path |
|---------|------|
| Modify  | `services/agent-service/internal/agent/agent.go` |
| Modify  | `services/agent-service/internal/agent/agent_test.go` |
| Create  | `services/agent-service/internal/agent/labour_power.go` |
| Create  | `services/agent-service/internal/agent/labour_power_test.go` |
| Modify  | `services/agent-service/internal/store/store.go` |
| Modify  | `services/agent-service/internal/store/memory.go` |
| Modify  | `services/agent-service/internal/store/memory_test.go` |
| Create  | `services/agent-service/internal/store/labour_power_memory_test.go` |
| Create  | `services/agent-service/internal/store/migrations/00005_ch06_labour_power.sql` |
| Modify  | `services/agent-service/internal/store/mysql.go` |
| Modify  | `services/agent-service/internal/transport/httpapi/handler.go` |
| Create  | `services/agent-service/internal/transport/httpapi/labour_power_handler.go` |
| Modify  | `services/agent-service/internal/transport/httpapi/routes.go` |
| Modify  | `services/agent-service/cmd/agent-service/main.go` |
| Modify  | `services/api-gateway/cmd/api-gateway/main.go` |
| Modify  | `web/src/types.ts` |
| Modify  | `web/src/api.ts` |
| Create  | `web/src/chapters/Ch06LabourPower.tsx` |
| Modify  | `web/src/components/ChapterShell.tsx` |
| Modify  | `web/src/chapters/registry.ts` |
| Modify  | `docs/architecture.md` |

---

## Task 1: Rename conflicting `Class` constants

**Files:**
- Modify: `services/agent-service/internal/agent/agent.go`
- Modify: `services/agent-service/internal/agent/agent_test.go`
- Modify: `services/agent-service/internal/store/memory_test.go`
- Modify: `services/agent-service/internal/transport/httpapi/handler.go`

- [ ] **Step 1: In `agent/agent.go` rename `Worker` and `Capitalist` constants**

Replace this block in `agent.go`:
```go
const (
	Capitalist Class = "capitalist"
	Worker     Class = "worker"
	Miser      Class = "miser"
	Owner      Class = "owner"
)
```
With:
```go
const (
	CapitalistClass Class = "capitalist"
	WorkerClass     Class = "worker"
	Miser           Class = "miser"
	Owner           Class = "owner"
)
```

- [ ] **Step 2: Update `agent/agent.go` method that references the old constant names**

In `agent.go`, the `Validate()` method references nothing (agent class validation uses a switch), and `Advance()` references `Miser` (unchanged). Check for any usage of the old `Worker` or `Capitalist` constant names — there are none in the methods themselves (the constants appear only in tests and handler). No changes needed inside `agent.go` methods beyond the const rename above.

- [ ] **Step 3: Update `agent_test.go` — replace constant references**

In `agent_test.go`, update every `agent.Worker` → `agent.WorkerClass` and `agent.Capitalist` → `agent.CapitalistClass`:

```go
// Change this helper:
func newCapitalist(balance agent.Pence) agent.Agent {
	return agent.Agent{ID: agent.NewID(), Name: "Capitalist", Class: agent.CapitalistClass, MoneyBalance: balance}
}

func newWorker(balance agent.Pence) agent.Agent {
	return agent.Agent{ID: agent.NewID(), Name: "Worker", Class: agent.WorkerClass, MoneyBalance: balance}
}
```

- [ ] **Step 4: Update `store/memory_test.go` — replace constant references**

Replace all `agent.Capitalist` → `agent.CapitalistClass` and `agent.Worker` → `agent.WorkerClass`:

```go
// Line ~20: makeAgent usage
func makeAgent(class agent.Class, balance agent.Pence) agent.Agent {
	return agent.Agent{Name: "Test Agent", Class: class, MoneyBalance: balance}
}

// Line ~20: TestMemory_CreateGet
created, err := m.Create(ctx, makeAgent(agent.CapitalistClass, 10000))

// TestMemory_ListByClass
if _, err := m.Create(ctx, makeAgent(agent.CapitalistClass, 10000)); err != nil {
if _, err := m.Create(ctx, makeAgent(agent.WorkerClass, 500)); err != nil {
if _, err := m.Create(ctx, agent.Agent{Name: "Worker2", Class: agent.WorkerClass, MoneyBalance: 300}); err != nil {
caps, err := m.ListByClass(ctx, agent.CapitalistClass)
// and so on — replace all occurrences
```

- [ ] **Step 5: Update `transport/httpapi/handler.go` — replace constant references**

In `handler.go` the `CreateCircuit` handler contains:
```go
if (a.Class == agent.Worker || a.Class == agent.Owner) && req.CircuitType == agent.CircuitMCM {
```
Change to:
```go
if (a.Class == agent.WorkerClass || a.Class == agent.Owner) && req.CircuitType == agent.CircuitMCM {
```

- [ ] **Step 6: Ask user to verify compilation**

Run: `make vet test build`
Expected: all tests pass, binary builds with no errors.

---

## Task 2: New domain types

**Files:**
- Create: `services/agent-service/internal/agent/labour_power.go`
- Create: `services/agent-service/internal/agent/labour_power_test.go`

- [ ] **Step 1: Write failing tests in `agent/labour_power_test.go`**

```go
package agent_test

import (
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// §2: LabourPower must have positive capacity to be valid
func TestLabourPower_ZeroCapacity(t *testing.T) {
	t.Parallel()
	w := agent.Worker{
		OwnsLabourPower:       true,
		OwnsCommoditiesToSell: false,
		LabourPower:           agent.LabourPower{CapacityMinutesPerDay: 0},
	}
	if err := w.Validate(); err == nil {
		t.Error("worker with zero capacity should fail validation")
	}
}

// §3: IsFreeLabourer returns false when worker does not own their labour-power
func TestIsFreeLabourer_NotOwner(t *testing.T) {
	t.Parallel()
	w := agent.Worker{OwnsLabourPower: false, OwnsCommoditiesToSell: false}
	if agent.IsFreeLabourer(w) {
		t.Error("worker who does not own labour-power is not a free labourer")
	}
}

// §4: IsFreeLabourer returns false when worker owns other commodities to sell
func TestIsFreeLabourer_OwnsCommodities(t *testing.T) {
	t.Parallel()
	w := agent.Worker{OwnsLabourPower: true, OwnsCommoditiesToSell: true}
	if agent.IsFreeLabourer(w) {
		t.Error("worker who owns commodities is not a free labourer in Marx's double sense")
	}
}

// Invariant: IsFreeLabourer == (OwnsLabourPower && !OwnsCommoditiesToSell)
func TestIsFreeLabourer_DoubleFreedom(t *testing.T) {
	t.Parallel()
	w := agent.Worker{OwnsLabourPower: true, OwnsCommoditiesToSell: false}
	if !agent.IsFreeLabourer(w) {
		t.Error("worker with double freedom should be a free labourer")
	}
}

// §3: ContractDays <= 0 returns ErrInvalidContract
func TestLabourPowerOffering_ZeroContract(t *testing.T) {
	t.Parallel()
	o := agent.LabourPowerOffering{
		OwnerID:               agent.NewAgentID(),
		CapacityMinutesPerDay: 480,
		ContractDays:          0,
		AskingWage:            240,
	}
	if err := o.Validate(); !errors.Is(err, agent.ErrInvalidContract) {
		t.Errorf("ContractDays=0: want ErrInvalidContract, got %v", err)
	}
}

// §3: Negative ContractDays also returns ErrInvalidContract
func TestLabourPowerOffering_NegativeContract(t *testing.T) {
	t.Parallel()
	o := agent.LabourPowerOffering{
		OwnerID:               agent.NewAgentID(),
		CapacityMinutesPerDay: 480,
		ContractDays:          -5,
		AskingWage:            240,
	}
	if err := o.Validate(); !errors.Is(err, agent.ErrInvalidContract) {
		t.Errorf("ContractDays=-5: want ErrInvalidContract, got %v", err)
	}
}

// §5: DailyValue == SubsistenceBasket.TotalSNLT()
func TestLabourPowerValue_DailyValueEqualsBasketTotal(t *testing.T) {
	t.Parallel()
	basket := agent.SubsistenceBasket{
		{Name: "food", SNLTMinutes: 120, Essential: true},
		{Name: "shelter", SNLTMinutes: 120, Essential: false},
	}
	lpv := agent.LabourPowerValue{Basket: basket}
	if lpv.DailyValue() != basket.TotalSNLT() {
		t.Errorf("DailyValue (%d) must equal basket.TotalSNLT() (%d)",
			lpv.DailyValue(), basket.TotalSNLT())
	}
}

// §5: Half-day subsistence cost → DailyValue = 240 labour-minutes
func TestLabourPowerValue_HalfDaySubsistence(t *testing.T) {
	t.Parallel()
	basket := agent.SubsistenceBasket{
		{Name: "subsistence", SNLTMinutes: 240, Essential: true},
	}
	lpv := agent.LabourPowerValue{Basket: basket}
	if got := lpv.DailyValue(); got != agent.LabourMinutes(240) {
		t.Errorf("DailyValue: want 240, got %d", got)
	}
}

// §5: MinimumValue equals SNLT of essential items only
func TestLabourPowerValue_MinimumValue(t *testing.T) {
	t.Parallel()
	basket := agent.SubsistenceBasket{
		{Name: "food", SNLTMinutes: 120, Essential: true},
		{Name: "clothing", SNLTMinutes: 60, Essential: true},
		{Name: "leisure", SNLTMinutes: 60, Essential: false},
	}
	lpv := agent.LabourPowerValue{Basket: basket}
	if got := lpv.MinimumValue(); got != agent.LabourMinutes(180) {
		t.Errorf("MinimumValue: want 180, got %d", got)
	}
}

// Invariant: MinimumValue() <= DailyValue()
func TestLabourPowerValue_MinimumLEDailyValue(t *testing.T) {
	t.Parallel()
	basket := agent.SubsistenceBasket{
		{Name: "food", SNLTMinutes: 120, Essential: true},
		{Name: "tobacco", SNLTMinutes: 60, Essential: false},
	}
	lpv := agent.LabourPowerValue{Basket: basket}
	if lpv.MinimumValue() > lpv.DailyValue() {
		t.Errorf("MinimumValue (%d) must be <= DailyValue (%d)",
			lpv.MinimumValue(), lpv.DailyValue())
	}
}

// ReproductionCost == DailyValue under normal conditions [§5]
func TestLabourPowerValue_ReproductionCost(t *testing.T) {
	t.Parallel()
	basket := agent.SubsistenceBasket{
		{Name: "food", SNLTMinutes: 240, Essential: true},
	}
	lpv := agent.LabourPowerValue{Basket: basket}
	if lpv.ReproductionCost() != lpv.DailyValue() {
		t.Errorf("ReproductionCost (%d) must equal DailyValue (%d)",
			lpv.ReproductionCost(), lpv.DailyValue())
	}
}

// §1: WageMinutes == DailyValue is an exchange of equivalents — purchase is valid
func TestLabourPowerPurchase_FairPrice_Valid(t *testing.T) {
	t.Parallel()
	basket := agent.SubsistenceBasket{
		{Name: "subsistence", SNLTMinutes: 240, Essential: true},
	}
	lpv := agent.LabourPowerValue{Basket: basket}
	p := agent.LabourPowerPurchase{
		SellerID:     agent.NewAgentID(),
		BuyerID:      agent.NewAgentID(),
		WageMinutes:  lpv.DailyValue(),
		ContractDays: 1,
	}
	if err := p.Validate(); err != nil {
		t.Errorf("fair-price purchase should be valid: %v", err)
	}
}

// Invariant: WageMinutes >= 0
func TestLabourPowerPurchase_NegativeWage_Invalid(t *testing.T) {
	t.Parallel()
	p := agent.LabourPowerPurchase{
		SellerID:     agent.NewAgentID(),
		BuyerID:      agent.NewAgentID(),
		WageMinutes:  -1,
		ContractDays: 1,
	}
	if err := p.Validate(); err == nil {
		t.Error("negative wage should fail validation")
	}
}

// NewAgentID produces a non-empty 24-char hex string
func TestNewAgentID(t *testing.T) {
	t.Parallel()
	id := agent.NewAgentID()
	if id.IsZero() {
		t.Error("NewAgentID should not produce zero ID")
	}
	if len(string(id)) != 24 {
		t.Errorf("AgentID should be 24 chars, got %d", len(string(id)))
	}
}

// NewPurchaseID produces a non-empty 24-char hex string
func TestNewPurchaseID(t *testing.T) {
	t.Parallel()
	id := agent.NewPurchaseID()
	if id.IsZero() {
		t.Error("NewPurchaseID should not produce zero ID")
	}
	if len(string(id)) != 24 {
		t.Errorf("PurchaseID should be 24 chars, got %d", len(string(id)))
	}
}
```

- [ ] **Step 2: Ask user to verify tests fail**

Run: `make vet test build`
Expected: compilation errors — `agent.LabourPower`, `agent.Worker` (type), etc. are undefined.

- [ ] **Step 3: Create `agent/labour_power.go`**

```go
package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

// LabourMinutes is the canonical value-magnitude unit: socially-necessary labour-time
// expressed in minutes.
type LabourMinutes int64

// AgentID is a 96-bit hex identifier for ch6 Worker and Capitalist agents.
type AgentID string

func (id AgentID) IsZero() bool { return id == "" }

func NewAgentID() AgentID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return AgentID(hex.EncodeToString(b))
}

// PurchaseID is a 96-bit hex identifier for LabourPowerPurchase records.
type PurchaseID string

func (id PurchaseID) IsZero() bool { return id == "" }

func NewPurchaseID() PurchaseID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return PurchaseID(hex.EncodeToString(b))
}

// AgentKind distinguishes the two sides of the labour market.
type AgentKind string

const (
	AgentKindWorker     AgentKind = "worker"
	AgentKindCapitalist AgentKind = "capitalist"
)

// ErrInvalidContract is returned when ContractDays is not positive and finite.
var ErrInvalidContract = errors.New("agent: contract duration must be positive and finite")

// LabourAgent holds the shared identity fields embedded by Worker and Capitalist.
// (The spec calls this "Agent"; the name LabourAgent avoids a conflict with the
// existing agent.Agent type introduced in Ch. 4.)
type LabourAgent struct {
	ID        AgentID   `json:"id"`
	Kind      AgentKind `json:"kind"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LabourPower is the aggregate of mental and physical capabilities existing in a
// human being that is put to work during the contracted period [§2].
type LabourPower struct {
	CapacityMinutesPerDay LabourMinutes `json:"capacity_minutes_per_day"`
}

// Worker is the agent who owns their labour-power and sells it as a commodity.
// They are "free" in the double sense: free to sell, and freed from other means
// of subsistence [§3, §4].
type Worker struct {
	LabourAgent
	OwnsLabourPower       bool        `json:"owns_labour_power"`
	OwnsCommoditiesToSell bool        `json:"owns_commodities_to_sell"`
	LabourPower           LabourPower `json:"labour_power"`
}

func (w Worker) Validate() error {
	if w.LabourPower.CapacityMinutesPerDay <= 0 {
		return errors.New("agent: labour_power capacity_minutes_per_day must be positive")
	}
	return nil
}

// Capitalist is the agent who possesses money-capital and purchases labour-power.
type Capitalist struct {
	LabourAgent
	MoneyCapital LabourMinutes `json:"money_capital"`
}

func (c Capitalist) Validate() error {
	if c.MoneyCapital < 0 {
		return errors.New("agent: money_capital cannot be negative")
	}
	return nil
}

// IsFreeLabourer returns true when the worker satisfies Marx's "double freedom":
// owns their labour-power AND lacks other commodities to sell [§3, §4].
func IsFreeLabourer(w Worker) bool {
	return w.OwnsLabourPower && !w.OwnsCommoditiesToSell
}

// SubsistenceItem is a single commodity in the worker's subsistence basket.
// Essential marks physically indispensable items used by MinimumValue [§5].
type SubsistenceItem struct {
	Name        string        `json:"name"`
	SNLTMinutes LabourMinutes `json:"snlt_minutes"`
	Essential   bool          `json:"essential"`
}

// SubsistenceBasket is the set of commodities required to maintain and reproduce
// the worker [§5].
type SubsistenceBasket []SubsistenceItem

// TotalSNLT returns the sum of SNLT across all items in the basket.
func (b SubsistenceBasket) TotalSNLT() LabourMinutes {
	var total LabourMinutes
	for _, item := range b {
		total += item.SNLTMinutes
	}
	return total
}

// LabourPowerValue computes the value of labour-power from a SubsistenceBasket [§5].
type LabourPowerValue struct {
	Basket SubsistenceBasket `json:"basket"`
}

// DailyValue returns the value of one day's labour-power: the sum of SNLT of all
// subsistence items. Invariant: DailyValue() == Basket.TotalSNLT() [§5].
func (v LabourPowerValue) DailyValue() LabourMinutes {
	return v.Basket.TotalSNLT()
}

// MinimumValue returns the physical floor of labour-power value: SNLT of
// physically indispensable items only. Invariant: MinimumValue() <= DailyValue() [§5].
func (v LabourPowerValue) MinimumValue() LabourMinutes {
	var total LabourMinutes
	for _, item := range v.Basket {
		if item.Essential {
			total += item.SNLTMinutes
		}
	}
	return total
}

// ReproductionCost returns the total SNLT to reproduce labour-power for another
// day. Equals DailyValue() under normal conditions [§5].
func (v LabourPowerValue) ReproductionCost() LabourMinutes {
	return v.DailyValue()
}

// LabourPowerOffering is a worker's labour-power posted for sale.
// ID is added for persistence; it is not in the spec's field list but required
// for CRUD operations.
type LabourPowerOffering struct {
	ID                    AgentID       `json:"id"`
	OwnerID               AgentID       `json:"owner_id"`
	CapacityMinutesPerDay LabourMinutes `json:"capacity_minutes_per_day"`
	ContractDays          int64         `json:"contract_days"`
	AskingWage            LabourMinutes `json:"asking_wage"`
	CreatedAt             time.Time     `json:"created_at"`
}

func (o LabourPowerOffering) Validate() error {
	if o.OwnerID.IsZero() {
		return errors.New("agent: offering owner_id is required")
	}
	if o.CapacityMinutesPerDay <= 0 {
		return errors.New("agent: capacity_minutes_per_day must be positive")
	}
	if o.ContractDays <= 0 {
		return ErrInvalidContract
	}
	if o.AskingWage < 0 {
		return errors.New("agent: asking_wage cannot be negative")
	}
	return nil
}

// LabourPowerPurchase is the completed transaction of buying labour-power.
type LabourPowerPurchase struct {
	ID           PurchaseID    `json:"id"`
	SellerID     AgentID       `json:"seller_id"`
	BuyerID      AgentID       `json:"buyer_id"`
	WageMinutes  LabourMinutes `json:"wage_minutes"`
	ContractDays int64         `json:"contract_days"`
	CreatedAt    time.Time     `json:"created_at"`
}

func (p LabourPowerPurchase) Validate() error {
	if p.SellerID.IsZero() {
		return errors.New("agent: purchase seller_id is required")
	}
	if p.BuyerID.IsZero() {
		return errors.New("agent: purchase buyer_id is required")
	}
	if p.WageMinutes < 0 {
		return errors.New("agent: wage_minutes cannot be negative")
	}
	if p.ContractDays <= 0 {
		return ErrInvalidContract
	}
	return nil
}
```

- [ ] **Step 4: Ask user to run tests**

Run: `make vet test build`
Expected: all tests in `agent/labour_power_test.go` pass; full binary builds.

---

## Task 3: Store interface + Memory implementation

**Files:**
- Modify: `services/agent-service/internal/store/store.go`
- Modify: `services/agent-service/internal/store/memory.go`
- Create: `services/agent-service/internal/store/labour_power_memory_test.go`

- [ ] **Step 1: Write failing store tests in `store/labour_power_memory_test.go`**

```go
package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

func makeWorker() agent.Worker {
	return agent.Worker{
		OwnsLabourPower:       true,
		OwnsCommoditiesToSell: false,
		LabourPower:           agent.LabourPower{CapacityMinutesPerDay: 480},
	}
}

func makeCapitalist() agent.Capitalist {
	return agent.Capitalist{
		MoneyCapital: 1000,
	}
}

func makeOffering(ownerID agent.AgentID) agent.LabourPowerOffering {
	return agent.LabourPowerOffering{
		OwnerID:               ownerID,
		CapacityMinutesPerDay: 480,
		ContractDays:          5,
		AskingWage:            240,
	}
}

func TestMemory_CreateGetWorker(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	created, err := m.CreateWorker(ctx, makeWorker())
	if err != nil {
		t.Fatalf("CreateWorker: %v", err)
	}
	if created.ID.IsZero() {
		t.Error("CreateWorker should assign an ID")
	}
	if created.Kind != agent.AgentKindWorker {
		t.Errorf("want Kind %q, got %q", agent.AgentKindWorker, created.Kind)
	}
	got, err := m.GetWorker(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetWorker: %v", err)
	}
	if got.LabourPower.CapacityMinutesPerDay != created.LabourPower.CapacityMinutesPerDay {
		t.Errorf("capacity mismatch: want %d, got %d",
			created.LabourPower.CapacityMinutesPerDay,
			got.LabourPower.CapacityMinutesPerDay)
	}
}

func TestMemory_GetWorker_NotFound(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.GetWorker(context.Background(), agent.NewAgentID())
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestMemory_ListWorkers(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	if _, err := m.CreateWorker(ctx, makeWorker()); err != nil {
		t.Fatalf("CreateWorker: %v", err)
	}
	if _, err := m.CreateWorker(ctx, makeWorker()); err != nil {
		t.Fatalf("CreateWorker 2: %v", err)
	}
	workers, err := m.ListWorkers(ctx)
	if err != nil {
		t.Fatalf("ListWorkers: %v", err)
	}
	if len(workers) != 2 {
		t.Errorf("want 2 workers, got %d", len(workers))
	}
}

func TestMemory_CreateGetCapitalist(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	created, err := m.CreateCapitalist(ctx, makeCapitalist())
	if err != nil {
		t.Fatalf("CreateCapitalist: %v", err)
	}
	if created.ID.IsZero() {
		t.Error("CreateCapitalist should assign an ID")
	}
	if created.Kind != agent.AgentKindCapitalist {
		t.Errorf("want Kind %q, got %q", agent.AgentKindCapitalist, created.Kind)
	}
	got, err := m.GetCapitalist(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetCapitalist: %v", err)
	}
	if got.MoneyCapital != created.MoneyCapital {
		t.Errorf("money_capital mismatch: want %d, got %d", created.MoneyCapital, got.MoneyCapital)
	}
}

func TestMemory_CreateGetOffering(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	ownerID := agent.NewAgentID()
	created, err := m.CreateOffering(ctx, makeOffering(ownerID))
	if err != nil {
		t.Fatalf("CreateOffering: %v", err)
	}
	if created.ID.IsZero() {
		t.Error("CreateOffering should assign an ID")
	}
	list, err := m.ListOfferings(ctx)
	if err != nil {
		t.Fatalf("ListOfferings: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("want 1 offering, got %d", len(list))
	}
	if list[0].ID != created.ID {
		t.Errorf("listing ID mismatch: want %q, got %q", created.ID, list[0].ID)
	}
}

func TestMemory_CreateGetPurchase(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	p := agent.LabourPowerPurchase{
		SellerID:     agent.NewAgentID(),
		BuyerID:      agent.NewAgentID(),
		WageMinutes:  240,
		ContractDays: 5,
	}
	created, err := m.CreatePurchase(ctx, p)
	if err != nil {
		t.Fatalf("CreatePurchase: %v", err)
	}
	if created.ID.IsZero() {
		t.Error("CreatePurchase should assign an ID")
	}
	got, err := m.GetPurchase(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetPurchase: %v", err)
	}
	if got.WageMinutes != p.WageMinutes {
		t.Errorf("wage_minutes mismatch: want %d, got %d", p.WageMinutes, got.WageMinutes)
	}
}

func TestMemory_GetPurchase_NotFound(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.GetPurchase(context.Background(), agent.NewPurchaseID())
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestMemory_ListPurchases(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	p := agent.LabourPowerPurchase{
		SellerID: agent.NewAgentID(), BuyerID: agent.NewAgentID(),
		WageMinutes: 240, ContractDays: 1,
	}
	if _, err := m.CreatePurchase(ctx, p); err != nil {
		t.Fatalf("CreatePurchase: %v", err)
	}
	list, err := m.ListPurchases(ctx)
	if err != nil {
		t.Fatalf("ListPurchases: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("want 1 purchase, got %d", len(list))
	}
}
```

- [ ] **Step 2: Add `LabourPowerStore` interface to `store/store.go`**

Append to the end of `store/store.go`:
```go
// LabourPowerStore is the persistence contract for Ch. 6 labour-power domain objects.
type LabourPowerStore interface {
	CreateWorker(ctx context.Context, w agent.Worker) (agent.Worker, error)
	GetWorker(ctx context.Context, id agent.AgentID) (agent.Worker, error)
	ListWorkers(ctx context.Context) ([]agent.Worker, error)
	CreateCapitalist(ctx context.Context, c agent.Capitalist) (agent.Capitalist, error)
	GetCapitalist(ctx context.Context, id agent.AgentID) (agent.Capitalist, error)
	ListCapitalists(ctx context.Context) ([]agent.Capitalist, error)
	CreateOffering(ctx context.Context, o agent.LabourPowerOffering) (agent.LabourPowerOffering, error)
	ListOfferings(ctx context.Context) ([]agent.LabourPowerOffering, error)
	CreatePurchase(ctx context.Context, p agent.LabourPowerPurchase) (agent.LabourPowerPurchase, error)
	GetPurchase(ctx context.Context, id agent.PurchaseID) (agent.LabourPowerPurchase, error)
	ListPurchases(ctx context.Context) ([]agent.LabourPowerPurchase, error)
}
```

- [ ] **Step 3: Add ch6 maps to `Memory` struct and `NewMemory()` in `store/memory.go`**

In the `Memory` struct, add four new map fields:
```go
type Memory struct {
	mu                sync.RWMutex
	agents            map[agent.ID]agent.Agent
	circuits          map[agent.ID]agent.CapitalCircuit
	labourWorkers     map[agent.AgentID]agent.Worker
	labourCapitalists map[agent.AgentID]agent.Capitalist
	offerings         map[agent.AgentID]agent.LabourPowerOffering
	purchases         map[agent.PurchaseID]agent.LabourPowerPurchase
	now               func() time.Time
}
```

In `NewMemory()`, initialise the four new maps:
```go
func NewMemory() *Memory {
	return &Memory{
		agents:            make(map[agent.ID]agent.Agent),
		circuits:          make(map[agent.ID]agent.CapitalCircuit),
		labourWorkers:     make(map[agent.AgentID]agent.Worker),
		labourCapitalists: make(map[agent.AgentID]agent.Capitalist),
		offerings:         make(map[agent.AgentID]agent.LabourPowerOffering),
		purchases:         make(map[agent.PurchaseID]agent.LabourPowerPurchase),
		now:               time.Now,
	}
}
```

- [ ] **Step 4: Append `LabourPowerStore` methods to `store/memory.go`**

Add the following at the end of `memory.go`:
```go
func (m *Memory) CreateWorker(_ context.Context, w agent.Worker) (agent.Worker, error) {
	if err := w.Validate(); err != nil {
		return agent.Worker{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if w.ID.IsZero() {
		w.ID = agent.NewAgentID()
	}
	now := m.now()
	w.CreatedAt = now
	w.UpdatedAt = now
	w.Kind = agent.AgentKindWorker
	m.labourWorkers[w.ID] = w
	return w, nil
}

func (m *Memory) GetWorker(_ context.Context, id agent.AgentID) (agent.Worker, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	w, ok := m.labourWorkers[id]
	if !ok {
		return agent.Worker{}, ErrNotFound
	}
	return w, nil
}

func (m *Memory) ListWorkers(_ context.Context) ([]agent.Worker, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.Worker, 0, len(m.labourWorkers))
	for _, w := range m.labourWorkers {
		out = append(out, w)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (m *Memory) CreateCapitalist(_ context.Context, c agent.Capitalist) (agent.Capitalist, error) {
	if err := c.Validate(); err != nil {
		return agent.Capitalist{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if c.ID.IsZero() {
		c.ID = agent.NewAgentID()
	}
	now := m.now()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.Kind = agent.AgentKindCapitalist
	m.labourCapitalists[c.ID] = c
	return c, nil
}

func (m *Memory) GetCapitalist(_ context.Context, id agent.AgentID) (agent.Capitalist, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.labourCapitalists[id]
	if !ok {
		return agent.Capitalist{}, ErrNotFound
	}
	return c, nil
}

func (m *Memory) ListCapitalists(_ context.Context) ([]agent.Capitalist, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.Capitalist, 0, len(m.labourCapitalists))
	for _, c := range m.labourCapitalists {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (m *Memory) CreateOffering(_ context.Context, o agent.LabourPowerOffering) (agent.LabourPowerOffering, error) {
	if err := o.Validate(); err != nil {
		return agent.LabourPowerOffering{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if o.ID.IsZero() {
		o.ID = agent.NewAgentID()
	}
	o.CreatedAt = m.now()
	m.offerings[o.ID] = o
	return o, nil
}

func (m *Memory) ListOfferings(_ context.Context) ([]agent.LabourPowerOffering, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.LabourPowerOffering, 0, len(m.offerings))
	for _, o := range m.offerings {
		out = append(out, o)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (m *Memory) CreatePurchase(_ context.Context, p agent.LabourPowerPurchase) (agent.LabourPowerPurchase, error) {
	if err := p.Validate(); err != nil {
		return agent.LabourPowerPurchase{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if p.ID.IsZero() {
		p.ID = agent.NewPurchaseID()
	}
	p.CreatedAt = m.now()
	m.purchases[p.ID] = p
	return p, nil
}

func (m *Memory) GetPurchase(_ context.Context, id agent.PurchaseID) (agent.LabourPowerPurchase, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.purchases[id]
	if !ok {
		return agent.LabourPowerPurchase{}, ErrNotFound
	}
	return p, nil
}

func (m *Memory) ListPurchases(_ context.Context) ([]agent.LabourPowerPurchase, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.LabourPowerPurchase, 0, len(m.purchases))
	for _, p := range m.purchases {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}
```

- [ ] **Step 5: Ask user to run tests**

Run: `make vet test build`
Expected: all store tests pass, including the new `labour_power_memory_test.go`.

---

## Task 4: MySQL migration + MySQL store implementation

**Files:**
- Create: `services/agent-service/internal/store/migrations/00005_ch06_labour_power.sql`
- Modify: `services/agent-service/internal/store/mysql.go`

- [ ] **Step 1: Create migration file `00005_ch06_labour_power.sql`**

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS labour_workers (
    id                       VARCHAR(24)  NOT NULL PRIMARY KEY,
    kind                     VARCHAR(50)  NOT NULL DEFAULT 'worker',
    owns_labour_power        TINYINT(1)   NOT NULL DEFAULT 1,
    owns_commodities_to_sell TINYINT(1)   NOT NULL DEFAULT 0,
    capacity_minutes_per_day BIGINT       NOT NULL DEFAULT 0,
    created_at               DATETIME(6)  NOT NULL,
    updated_at               DATETIME(6)  NOT NULL
);

CREATE TABLE IF NOT EXISTS labour_capitalists (
    id            VARCHAR(24)  NOT NULL PRIMARY KEY,
    kind          VARCHAR(50)  NOT NULL DEFAULT 'capitalist',
    money_capital BIGINT       NOT NULL DEFAULT 0,
    created_at    DATETIME(6)  NOT NULL,
    updated_at    DATETIME(6)  NOT NULL
);

CREATE TABLE IF NOT EXISTS labour_power_offerings (
    id                       VARCHAR(24)  NOT NULL PRIMARY KEY,
    owner_id                 VARCHAR(24)  NOT NULL,
    capacity_minutes_per_day BIGINT       NOT NULL,
    contract_days            BIGINT       NOT NULL,
    asking_wage              BIGINT       NOT NULL,
    created_at               DATETIME(6)  NOT NULL,
    INDEX idx_owner_id (owner_id)
);

CREATE TABLE IF NOT EXISTS labour_power_purchases (
    id            VARCHAR(24)  NOT NULL PRIMARY KEY,
    seller_id     VARCHAR(24)  NOT NULL,
    buyer_id      VARCHAR(24)  NOT NULL,
    wage_minutes  BIGINT       NOT NULL,
    contract_days BIGINT       NOT NULL,
    created_at    DATETIME(6)  NOT NULL,
    INDEX idx_seller_id (seller_id),
    INDEX idx_buyer_id  (buyer_id)
);

-- +goose Down
DROP TABLE IF EXISTS labour_power_purchases;
DROP TABLE IF EXISTS labour_power_offerings;
DROP TABLE IF EXISTS labour_capitalists;
DROP TABLE IF EXISTS labour_workers;
```

- [ ] **Step 2: Append `LabourPowerStore` methods to `store/mysql.go`**

Add the following to the end of `mysql.go`:
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
		(id, kind, owns_labour_power, owns_commodities_to_sell, capacity_minutes_per_day, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(w.ID), string(w.Kind),
		w.OwnsLabourPower, w.OwnsCommoditiesToSell,
		int64(w.LabourPower.CapacityMinutesPerDay),
		w.CreatedAt, w.UpdatedAt,
	)
	if err != nil {
		return agent.Worker{}, err
	}
	return w, nil
}

func (m *MySQL) GetWorker(ctx context.Context, id agent.AgentID) (agent.Worker, error) {
	const q = `SELECT id, kind, owns_labour_power, owns_commodities_to_sell,
		capacity_minutes_per_day, created_at, updated_at
		FROM labour_workers WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanWorker(row)
}

func (m *MySQL) ListWorkers(ctx context.Context) ([]agent.Worker, error) {
	const q = `SELECT id, kind, owns_labour_power, owns_commodities_to_sell,
		capacity_minutes_per_day, created_at, updated_at
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

func (m *MySQL) CreateCapitalist(ctx context.Context, c agent.Capitalist) (agent.Capitalist, error) {
	if err := c.Validate(); err != nil {
		return agent.Capitalist{}, err
	}
	if c.ID.IsZero() {
		c.ID = agent.NewAgentID()
	}
	now := m.now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.Kind = agent.AgentKindCapitalist
	const q = `INSERT INTO labour_capitalists
		(id, kind, money_capital, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID), string(c.Kind),
		int64(c.MoneyCapital),
		c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return agent.Capitalist{}, err
	}
	return c, nil
}

func (m *MySQL) GetCapitalist(ctx context.Context, id agent.AgentID) (agent.Capitalist, error) {
	const q = `SELECT id, kind, money_capital, created_at, updated_at
		FROM labour_capitalists WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanCapitalist(row)
}

func (m *MySQL) ListCapitalists(ctx context.Context) ([]agent.Capitalist, error) {
	const q = `SELECT id, kind, money_capital, created_at, updated_at
		FROM labour_capitalists ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []agent.Capitalist
	for rows.Next() {
		c, err := scanCapitalistRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (m *MySQL) CreateOffering(ctx context.Context, o agent.LabourPowerOffering) (agent.LabourPowerOffering, error) {
	if err := o.Validate(); err != nil {
		return agent.LabourPowerOffering{}, err
	}
	if o.ID.IsZero() {
		o.ID = agent.NewAgentID()
	}
	o.CreatedAt = m.now().UTC()
	const q = `INSERT INTO labour_power_offerings
		(id, owner_id, capacity_minutes_per_day, contract_days, asking_wage, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(o.ID), string(o.OwnerID),
		int64(o.CapacityMinutesPerDay), o.ContractDays,
		int64(o.AskingWage), o.CreatedAt,
	)
	if err != nil {
		return agent.LabourPowerOffering{}, err
	}
	return o, nil
}

func (m *MySQL) ListOfferings(ctx context.Context) ([]agent.LabourPowerOffering, error) {
	const q = `SELECT id, owner_id, capacity_minutes_per_day, contract_days, asking_wage, created_at
		FROM labour_power_offerings ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []agent.LabourPowerOffering
	for rows.Next() {
		var o agent.LabourPowerOffering
		var id, ownerID string
		var cap, wage int64
		if err := rows.Scan(&id, &ownerID, &cap, &o.ContractDays, &wage, &o.CreatedAt); err != nil {
			return nil, err
		}
		o.ID = agent.AgentID(id)
		o.OwnerID = agent.AgentID(ownerID)
		o.CapacityMinutesPerDay = agent.LabourMinutes(cap)
		o.AskingWage = agent.LabourMinutes(wage)
		out = append(out, o)
	}
	return out, rows.Err()
}

func (m *MySQL) CreatePurchase(ctx context.Context, p agent.LabourPowerPurchase) (agent.LabourPowerPurchase, error) {
	if err := p.Validate(); err != nil {
		return agent.LabourPowerPurchase{}, err
	}
	if p.ID.IsZero() {
		p.ID = agent.NewPurchaseID()
	}
	p.CreatedAt = m.now().UTC()
	const q = `INSERT INTO labour_power_purchases
		(id, seller_id, buyer_id, wage_minutes, contract_days, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(p.ID), string(p.SellerID), string(p.BuyerID),
		int64(p.WageMinutes), p.ContractDays, p.CreatedAt,
	)
	if err != nil {
		return agent.LabourPowerPurchase{}, err
	}
	return p, nil
}

func (m *MySQL) GetPurchase(ctx context.Context, id agent.PurchaseID) (agent.LabourPowerPurchase, error) {
	const q = `SELECT id, seller_id, buyer_id, wage_minutes, contract_days, created_at
		FROM labour_power_purchases WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	var p agent.LabourPowerPurchase
	var pid, sellerID, buyerID string
	var wage int64
	err := row.Scan(&pid, &sellerID, &buyerID, &wage, &p.ContractDays, &p.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.LabourPowerPurchase{}, ErrNotFound
	}
	if err != nil {
		return agent.LabourPowerPurchase{}, err
	}
	p.ID = agent.PurchaseID(pid)
	p.SellerID = agent.AgentID(sellerID)
	p.BuyerID = agent.AgentID(buyerID)
	p.WageMinutes = agent.LabourMinutes(wage)
	return p, nil
}

func (m *MySQL) ListPurchases(ctx context.Context) ([]agent.LabourPowerPurchase, error) {
	const q = `SELECT id, seller_id, buyer_id, wage_minutes, contract_days, created_at
		FROM labour_power_purchases ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []agent.LabourPowerPurchase
	for rows.Next() {
		var p agent.LabourPowerPurchase
		var pid, sellerID, buyerID string
		var wage int64
		if err := rows.Scan(&pid, &sellerID, &buyerID, &wage, &p.ContractDays, &p.CreatedAt); err != nil {
			return nil, err
		}
		p.ID = agent.PurchaseID(pid)
		p.SellerID = agent.AgentID(sellerID)
		p.BuyerID = agent.AgentID(buyerID)
		p.WageMinutes = agent.LabourMinutes(wage)
		out = append(out, p)
	}
	return out, rows.Err()
}

// scan helpers for labour-power domain objects

func scanWorker(row *sql.Row) (agent.Worker, error) {
	var w agent.Worker
	var id, kind string
	var cap int64
	err := row.Scan(&id, &kind, &w.OwnsLabourPower, &w.OwnsCommoditiesToSell, &cap, &w.CreatedAt, &w.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.Worker{}, ErrNotFound
	}
	if err != nil {
		return agent.Worker{}, err
	}
	w.ID = agent.AgentID(id)
	w.Kind = agent.AgentKind(kind)
	w.LabourPower.CapacityMinutesPerDay = agent.LabourMinutes(cap)
	return w, nil
}

func scanWorkerRow(rows *sql.Rows) (agent.Worker, error) {
	var w agent.Worker
	var id, kind string
	var cap int64
	if err := rows.Scan(&id, &kind, &w.OwnsLabourPower, &w.OwnsCommoditiesToSell, &cap, &w.CreatedAt, &w.UpdatedAt); err != nil {
		return agent.Worker{}, err
	}
	w.ID = agent.AgentID(id)
	w.Kind = agent.AgentKind(kind)
	w.LabourPower.CapacityMinutesPerDay = agent.LabourMinutes(cap)
	return w, nil
}

func scanCapitalist(row *sql.Row) (agent.Capitalist, error) {
	var c agent.Capitalist
	var id, kind string
	var mc int64
	err := row.Scan(&id, &kind, &mc, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.Capitalist{}, ErrNotFound
	}
	if err != nil {
		return agent.Capitalist{}, err
	}
	c.ID = agent.AgentID(id)
	c.Kind = agent.AgentKind(kind)
	c.MoneyCapital = agent.LabourMinutes(mc)
	return c, nil
}

func scanCapitalistRow(rows *sql.Rows) (agent.Capitalist, error) {
	var c agent.Capitalist
	var id, kind string
	var mc int64
	if err := rows.Scan(&id, &kind, &mc, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return agent.Capitalist{}, err
	}
	c.ID = agent.AgentID(id)
	c.Kind = agent.AgentKind(kind)
	c.MoneyCapital = agent.LabourMinutes(mc)
	return c, nil
}
```

- [ ] **Step 3: Ask user to run tests**

Run: `make vet test build`
Expected: compilation succeeds, existing tests pass (MySQL store methods are not exercised by unit tests — they are integration-tested at deploy time).

---

## Task 5: HTTP handlers, routes, and main.go wiring

**Files:**
- Modify: `services/agent-service/internal/transport/httpapi/handler.go`
- Create: `services/agent-service/internal/transport/httpapi/labour_power_handler.go`
- Modify: `services/agent-service/internal/transport/httpapi/routes.go`
- Modify: `services/agent-service/cmd/agent-service/main.go`

- [ ] **Step 1: Add `LabourPowerStore` field to `Handler` and update `New()`**

In `handler.go`, change the struct and constructor:
```go
type Handler struct {
	Store            store.Store
	CircuitStore     store.CircuitStore
	LabourPowerStore store.LabourPowerStore
	Logger           *slog.Logger
}

func New(s store.Store, cs store.CircuitStore, lps store.LabourPowerStore, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{Store: s, CircuitStore: cs, LabourPowerStore: lps, Logger: logger}
}
```

- [ ] **Step 2: Create `transport/httpapi/labour_power_handler.go`**

```go
package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

type createWorkerRequest struct {
	OwnsLabourPower       bool                `json:"owns_labour_power"`
	OwnsCommoditiesToSell bool                `json:"owns_commodities_to_sell"`
	CapacityMinutesPerDay agent.LabourMinutes `json:"capacity_minutes_per_day"`
}

type createCapitalistRequest struct {
	MoneyCapital agent.LabourMinutes `json:"money_capital"`
}

type createOfferingRequest struct {
	OwnerID               agent.AgentID       `json:"owner_id"`
	CapacityMinutesPerDay agent.LabourMinutes `json:"capacity_minutes_per_day"`
	ContractDays          int64               `json:"contract_days"`
	AskingWage            agent.LabourMinutes `json:"asking_wage"`
}

type createPurchaseRequest struct {
	SellerID     agent.AgentID       `json:"seller_id"`
	BuyerID      agent.AgentID       `json:"buyer_id"`
	WageMinutes  agent.LabourMinutes `json:"wage_minutes"`
	ContractDays int64               `json:"contract_days"`
}

func (h *Handler) CreateWorker(w http.ResponseWriter, r *http.Request) {
	var req createWorkerRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	worker := agent.Worker{
		OwnsLabourPower:       req.OwnsLabourPower,
		OwnsCommoditiesToSell: req.OwnsCommoditiesToSell,
		LabourPower:           agent.LabourPower{CapacityMinutesPerDay: req.CapacityMinutesPerDay},
	}
	if err := worker.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.LabourPowerStore.CreateWorker(r.Context(), worker)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) GetWorker(w http.ResponseWriter, r *http.Request) {
	id := agent.AgentID(r.PathValue("id"))
	worker, err := h.LabourPowerStore.GetWorker(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, worker)
}

func (h *Handler) ListWorkers(w http.ResponseWriter, r *http.Request) {
	workers, err := h.LabourPowerStore.ListWorkers(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}
	if workers == nil {
		workers = []agent.Worker{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": workers})
}

func (h *Handler) CreateCapitalist(w http.ResponseWriter, r *http.Request) {
	var req createCapitalistRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	cap := agent.Capitalist{
		MoneyCapital: req.MoneyCapital,
	}
	if err := cap.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.LabourPowerStore.CreateCapitalist(r.Context(), cap)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) GetCapitalist(w http.ResponseWriter, r *http.Request) {
	id := agent.AgentID(r.PathValue("id"))
	cap, err := h.LabourPowerStore.GetCapitalist(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, cap)
}

func (h *Handler) ListCapitalists(w http.ResponseWriter, r *http.Request) {
	caps, err := h.LabourPowerStore.ListCapitalists(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}
	if caps == nil {
		caps = []agent.Capitalist{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": caps})
}

func (h *Handler) CreateOffering(w http.ResponseWriter, r *http.Request) {
	var req createOfferingRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	offering := agent.LabourPowerOffering{
		OwnerID:               req.OwnerID,
		CapacityMinutesPerDay: req.CapacityMinutesPerDay,
		ContractDays:          req.ContractDays,
		AskingWage:            req.AskingWage,
	}
	if err := offering.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.LabourPowerStore.CreateOffering(r.Context(), offering)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) ListOfferings(w http.ResponseWriter, r *http.Request) {
	offerings, err := h.LabourPowerStore.ListOfferings(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}
	if offerings == nil {
		offerings = []agent.LabourPowerOffering{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": offerings})
}

func (h *Handler) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	var req createPurchaseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	purchase := agent.LabourPowerPurchase{
		SellerID:     req.SellerID,
		BuyerID:      req.BuyerID,
		WageMinutes:  req.WageMinutes,
		ContractDays: req.ContractDays,
	}
	if err := purchase.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.validatePurchaseParties(r, req.SellerID, req.BuyerID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.LabourPowerStore.CreatePurchase(r.Context(), purchase)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

// validatePurchaseParties checks that seller is a known worker and buyer is a known capitalist.
func (h *Handler) validatePurchaseParties(r *http.Request, sellerID, buyerID agent.AgentID) error {
	if _, err := h.LabourPowerStore.GetWorker(r.Context(), sellerID); err != nil {
		if err == store.ErrNotFound {
			return agent.ErrInvalidContract
		}
		return err
	}
	if _, err := h.LabourPowerStore.GetCapitalist(r.Context(), buyerID); err != nil {
		if err == store.ErrNotFound {
			return agent.ErrInvalidContract
		}
		return err
	}
	return nil
}

func (h *Handler) GetPurchase(w http.ResponseWriter, r *http.Request) {
	id := agent.PurchaseID(r.PathValue("id"))
	purchase, err := h.LabourPowerStore.GetPurchase(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, purchase)
}

func (h *Handler) ListPurchases(w http.ResponseWriter, r *http.Request) {
	purchases, err := h.LabourPowerStore.ListPurchases(r.Context())
	if err != nil {
		writeAppError(w, err)
		return
	}
	if purchases == nil {
		purchases = []agent.LabourPowerPurchase{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": purchases})
}
```

- [ ] **Step 3: Register new routes in `transport/httpapi/routes.go`**

Append the Ch. 6 routes:
```go
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
	s.HandleFunc("POST /v1/circuit-probes", h.ComputeCircuit)
	s.HandleFunc("POST /v1/exchange-simulations", h.ComputeExchange)
	// Ch. 6 — The Buying and Selling of Labour-Power
	s.HandleFunc("POST /v1/workers", h.CreateWorker)
	s.HandleFunc("GET /v1/workers", h.ListWorkers)
	s.HandleFunc("GET /v1/workers/{id}", h.GetWorker)
	s.HandleFunc("POST /v1/capitalists", h.CreateCapitalist)
	s.HandleFunc("GET /v1/capitalists", h.ListCapitalists)
	s.HandleFunc("GET /v1/capitalists/{id}", h.GetCapitalist)
	s.HandleFunc("POST /v1/labour-power/offerings", h.CreateOffering)
	s.HandleFunc("GET /v1/labour-power/offerings", h.ListOfferings)
	s.HandleFunc("POST /v1/labour-power/purchases", h.CreatePurchase)
	s.HandleFunc("GET /v1/labour-power/purchases", h.ListPurchases)
	s.HandleFunc("GET /v1/labour-power/purchases/{id}", h.GetPurchase)
}
```

- [ ] **Step 4: Update `cmd/agent-service/main.go` — add `LabourPowerStore` to `agentStore` and `New()` call**

Change the `agentStore` interface:
```go
type agentStore interface {
	store.Store
	store.CircuitStore
	store.LabourPowerStore
}
```

Change the `httpapi.New(...)` call:
```go
httpapi.Register(srv, httpapi.New(st, st, st, logger))
```

- [ ] **Step 5: Ask user to run tests**

Run: `make vet test build`
Expected: all tests pass, binary compiles.

---

## Task 6: API gateway — add Ch. 6 proxy routes

**Files:**
- Modify: `services/api-gateway/cmd/api-gateway/main.go`

- [ ] **Step 1: Add Ch. 6 proxy routes to the agent section and update info handler**

In the agent proxy section, append three route pairs:
```go
// existing agent routes
srv.Handle("/v1/agents", agentProxy)
srv.Handle("/v1/agents/{rest...}", agentProxy)
srv.Handle("/v1/circuit-probes", agentProxy)
srv.Handle("/v1/exchange-simulations", agentProxy)
// Ch. 6 additions
srv.Handle("/v1/workers", agentProxy)
srv.Handle("/v1/workers/{rest...}", agentProxy)
srv.Handle("/v1/capitalists", agentProxy)
srv.Handle("/v1/capitalists/{rest...}", agentProxy)
srv.Handle("/v1/labour-power", agentProxy)
srv.Handle("/v1/labour-power/{rest...}", agentProxy)
```

Update the `handleInfo` response:
```go
"status":  "ch-6-labour-power",
"description": "External entrypoint. ...; labour-power routes (workers/capitalists/offerings/purchases) to agent-service.",
"chapter": "Capital Vol. I, Ch. 6 - The Buying and Selling of Labour-Power",
```

- [ ] **Step 2: Ask user to build**

Run: `make vet test build`
Expected: gateway builds with no errors.

---

## Task 7: TypeScript types and API calls

**Files:**
- Modify: `web/src/types.ts`
- Modify: `web/src/api.ts`

- [ ] **Step 1: Append Ch. 6 types to `web/src/types.ts`**

Add at the end of `types.ts`:
```typescript
// --- agent-service types (Ch. 6: The Buying and Selling of Labour-Power) -----

export interface LabourPower {
  capacity_minutes_per_day: number; // labour-minutes
}

export interface SubsistenceItem {
  name: string;
  snlt_minutes: number; // labour-minutes
  essential: boolean;
}

export interface SubsistenceBasket {
  items: SubsistenceItem[];
  total_snlt: number; // labour-minutes
}

export interface LabourWorker {
  id: string;
  kind: "worker";
  owns_labour_power: boolean;
  owns_commodities_to_sell: boolean;
  labour_power: LabourPower;
  created_at: string;
  updated_at: string;
}

export interface LabourCapitalist {
  id: string;
  kind: "capitalist";
  money_capital: number; // labour-minutes (capital in value form)
  created_at: string;
  updated_at: string;
}

export interface LabourPowerOffering {
  id: string;
  owner_id: string;
  capacity_minutes_per_day: number; // labour-minutes
  contract_days: number;
  asking_wage: number; // labour-minutes
  created_at: string;
}

export interface LabourPowerPurchase {
  id: string;
  seller_id: string;
  buyer_id: string;
  wage_minutes: number; // labour-minutes
  contract_days: number;
  created_at: string;
}

export interface CreateLabourWorkerInput {
  owns_labour_power: boolean;
  owns_commodities_to_sell: boolean;
  capacity_minutes_per_day: number;
}

export interface CreateLabourCapitalistInput {
  money_capital: number;
}

export interface CreateLabourPowerOfferingInput {
  owner_id: string;
  capacity_minutes_per_day: number;
  contract_days: number;
  asking_wage: number;
}

export interface CreateLabourPowerPurchaseInput {
  seller_id: string;
  buyer_id: string;
  wage_minutes: number;
  contract_days: number;
}
```

- [ ] **Step 2: Add the new imports and API methods to `web/src/api.ts`**

Add to the import list at the top of `api.ts`:
```typescript
import type {
  // ... existing imports ...
  CreateLabourCapitalistInput,
  CreateLabourPowerOfferingInput,
  CreateLabourPowerPurchaseInput,
  CreateLabourWorkerInput,
  LabourCapitalist,
  LabourPowerOffering,
  LabourPowerPurchase,
  LabourWorker,
} from "./types";
```

Append to the `api` object (before the closing `}`):
```typescript
  // --- agent-service (Ch. 6: Labour-Power) ---

  listLabourWorkers: () =>
    http<{ items: LabourWorker[] }>("/v1/workers").then((r) => r.items),

  createLabourWorker: (input: CreateLabourWorkerInput) =>
    http<LabourWorker>("/v1/workers", { method: "POST", body: JSON.stringify(input) }),

  getLabourWorker: (id: string) => http<LabourWorker>(`/v1/workers/${id}`),

  listLabourCapitalists: () =>
    http<{ items: LabourCapitalist[] }>("/v1/capitalists").then((r) => r.items),

  createLabourCapitalist: (input: CreateLabourCapitalistInput) =>
    http<LabourCapitalist>("/v1/capitalists", { method: "POST", body: JSON.stringify(input) }),

  getLabourCapitalist: (id: string) => http<LabourCapitalist>(`/v1/capitalists/${id}`),

  listOfferings: () =>
    http<{ items: LabourPowerOffering[] }>("/v1/labour-power/offerings").then((r) => r.items),

  createOffering: (input: CreateLabourPowerOfferingInput) =>
    http<LabourPowerOffering>("/v1/labour-power/offerings", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  listLabourPurchases: () =>
    http<{ items: LabourPowerPurchase[] }>("/v1/labour-power/purchases").then((r) => r.items),

  createLabourPurchase: (input: CreateLabourPowerPurchaseInput) =>
    http<LabourPowerPurchase>("/v1/labour-power/purchases", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getLabourPurchase: (id: string) =>
    http<LabourPowerPurchase>(`/v1/labour-power/purchases/${id}`),
```

- [ ] **Step 3: Ask user to typecheck**

Run: `cd web && npm run lint`
Expected: no TypeScript errors.

---

## Task 8: React panel component

**Files:**
- Create: `web/src/chapters/Ch06LabourPower.tsx`
- Modify: `web/src/components/ChapterShell.tsx`
- Modify: `web/src/chapters/registry.ts`

- [ ] **Step 1: Create `web/src/chapters/Ch06LabourPower.tsx`**

```tsx
import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type {
  LabourWorker,
  LabourCapitalist,
  LabourPowerOffering,
  LabourPowerPurchase,
} from "../types";

interface Ch06Props {
  onSharedChanged: () => void;
}

function minutesToHours(m: number): string {
  const h = Math.floor(m / 60);
  const min = m % 60;
  return min === 0 ? `${h}h` : `${h}h ${min}m`;
}

export function Ch06LabourPower({ onSharedChanged: _onSharedChanged }: Ch06Props) {
  const [workers, setWorkers] = useState<LabourWorker[]>([]);
  const [capitalists, setCapitalists] = useState<LabourCapitalist[]>([]);
  const [offerings, setOfferings] = useState<LabourPowerOffering[]>([]);
  const [purchases, setPurchases] = useState<LabourPowerPurchase[]>([]);
  const [loadErr, setLoadErr] = useState<string | null>(null);

  async function refresh() {
    try {
      const [ws, cs, os, ps] = await Promise.all([
        api.listLabourWorkers(),
        api.listLabourCapitalists(),
        api.listOfferings(),
        api.listLabourPurchases(),
      ]);
      setWorkers(ws);
      setCapitalists(cs);
      setOfferings(os);
      setPurchases(ps);
    } catch (e) {
      setLoadErr(e instanceof Error ? e.message : String(e));
    }
  }

  useEffect(() => { refresh(); }, []);

  return (
    <>
      {loadErr && <p className="error">{loadErr}</p>}
      <RegisterWorkerPanel onCreated={refresh} />
      <RegisterCapitalistPanel onCreated={refresh} />
      <WorkerListPanel workers={workers} />
      <CapitalistListPanel capitalists={capitalists} />
      <PostOfferingPanel workers={workers} onCreated={refresh} />
      <ActiveOfferingsPanel offerings={offerings} workers={workers} />
      <PurchasePanel
        workers={workers}
        capitalists={capitalists}
        onCreated={refresh}
      />
      <PurchaseHistoryPanel purchases={purchases} workers={workers} capitalists={capitalists} />
    </>
  );
}

function RegisterWorkerPanel({ onCreated }: { onCreated: () => void }) {
  const [capacity, setCapacity] = useState(480);
  const [ownsLP, setOwnsLP] = useState(true);
  const [ownsCommodities, setOwnsCommodities] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await api.createLabourWorker({
        owns_labour_power: ownsLP,
        owns_commodities_to_sell: ownsCommodities,
        capacity_minutes_per_day: capacity,
      });
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Register Worker</h2>
      <p className="description">
        A labourer who owns their capacity for labour and is obliged to sell it.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Capacity (minutes/day)</span>
          <input
            type="number"
            value={capacity}
            onChange={(e) => setCapacity(Number(e.target.value))}
            min={1}
          />
        </label>
        <label>
          <span>Owns labour-power</span>
          <input type="checkbox" checked={ownsLP} onChange={(e) => setOwnsLP(e.target.checked)} />
        </label>
        <label>
          <span>Owns commodities to sell</span>
          <input
            type="checkbox"
            checked={ownsCommodities}
            onChange={(e) => setOwnsCommodities(e.target.checked)}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Register</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
    </section>
  );
}

function RegisterCapitalistPanel({ onCreated }: { onCreated: () => void }) {
  const [moneyCapital, setMoneyCapital] = useState(480);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await api.createLabourCapitalist({ money_capital: moneyCapital });
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Register Capitalist</h2>
      <p className="description">
        The owner of money-capital who purchases labour-power as a commodity.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Money capital (labour-minutes)</span>
          <input
            type="number"
            value={moneyCapital}
            onChange={(e) => setMoneyCapital(Number(e.target.value))}
            min={0}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Register</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
    </section>
  );
}

function WorkerListPanel({ workers }: { workers: LabourWorker[] }) {
  if (workers.length === 0) return null;
  return (
    <section className="card">
      <h2>Workers</h2>
      <div className="item-list">
        {workers.map((w) => (
          <div key={w.id} className="item-card">
            <div className="item-header">
              <span className="item-name monospace small">{w.id.slice(0, 8)}&hellip;</span>
              <span className="item-meta">
                {minutesToHours(w.labour_power.capacity_minutes_per_day)} / day
              </span>
              {w.owns_labour_power && !w.owns_commodities_to_sell && (
                <span className="item-tag">free labourer</span>
              )}
            </div>
            <p className="small muted">
              Owns LP: {w.owns_labour_power ? "yes" : "no"} &middot;
              Owns commodities: {w.owns_commodities_to_sell ? "yes" : "no"}
            </p>
          </div>
        ))}
      </div>
    </section>
  );
}

function CapitalistListPanel({ capitalists }: { capitalists: LabourCapitalist[] }) {
  if (capitalists.length === 0) return null;
  return (
    <section className="card">
      <h2>Capitalists</h2>
      <div className="item-list">
        {capitalists.map((c) => (
          <div key={c.id} className="item-card">
            <div className="item-header">
              <span className="item-name monospace small">{c.id.slice(0, 8)}&hellip;</span>
              <span className="item-meta">{minutesToHours(c.money_capital)} capital</span>
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}

function PostOfferingPanel({
  workers,
  onCreated,
}: {
  workers: LabourWorker[];
  onCreated: () => void;
}) {
  const [ownerID, setOwnerID] = useState("");
  const [capacity, setCapacity] = useState(480);
  const [contractDays, setContractDays] = useState(5);
  const [askingWage, setAskingWage] = useState(240);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await api.createOffering({
        owner_id: ownerID,
        capacity_minutes_per_day: capacity,
        contract_days: contractDays,
        asking_wage: askingWage,
      });
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Post Labour-Power for Sale</h2>
      <p className="description">
        A worker offers their capacity for labour for a finite period at an asking wage.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Worker</span>
          <select value={ownerID} onChange={(e) => setOwnerID(e.target.value)} required>
            <option value="">Select a worker…</option>
            {workers.map((w) => (
              <option key={w.id} value={w.id}>
                {w.id.slice(0, 8)}… ({minutesToHours(w.labour_power.capacity_minutes_per_day)}/day)
              </option>
            ))}
          </select>
        </label>
        <label>
          <span>Capacity (min/day)</span>
          <input
            type="number"
            value={capacity}
            onChange={(e) => setCapacity(Number(e.target.value))}
            min={1}
          />
        </label>
        <label>
          <span>Contract (days)</span>
          <input
            type="number"
            value={contractDays}
            onChange={(e) => setContractDays(Number(e.target.value))}
            min={1}
          />
        </label>
        <label>
          <span>Asking wage (labour-min/day)</span>
          <input
            type="number"
            value={askingWage}
            onChange={(e) => setAskingWage(Number(e.target.value))}
            min={0}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Post offering</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
    </section>
  );
}

function ActiveOfferingsPanel({
  offerings,
  workers,
}: {
  offerings: LabourPowerOffering[];
  workers: LabourWorker[];
}) {
  if (offerings.length === 0) return null;
  const workerMap = new Map(workers.map((w) => [w.id, w]));
  return (
    <section className="card">
      <h2>Active Offerings</h2>
      <table className="data-table">
        <thead>
          <tr>
            <th>Worker</th>
            <th>Capacity</th>
            <th>Days</th>
            <th>Asking wage / day</th>
          </tr>
        </thead>
        <tbody>
          {offerings.map((o) => {
            const w = workerMap.get(o.owner_id);
            return (
              <tr key={o.id}>
                <td className="monospace small">
                  {w ? `${o.owner_id.slice(0, 8)}…` : o.owner_id.slice(0, 8) + "…"}
                </td>
                <td>{minutesToHours(o.capacity_minutes_per_day)}</td>
                <td>{o.contract_days}</td>
                <td>{minutesToHours(o.asking_wage)}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </section>
  );
}

function PurchasePanel({
  workers,
  capitalists,
  onCreated,
}: {
  workers: LabourWorker[];
  capitalists: LabourCapitalist[];
  onCreated: () => void;
}) {
  const [sellerID, setSellerID] = useState("");
  const [buyerID, setBuyerID] = useState("");
  const [wageMinutes, setWageMinutes] = useState(240);
  const [contractDays, setContractDays] = useState(5);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await api.createLabourPurchase({
        seller_id: sellerID,
        buyer_id: buyerID,
        wage_minutes: wageMinutes,
        contract_days: contractDays,
      });
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Purchase Labour-Power</h2>
      <p className="description">
        The capitalist buys labour-power as a commodity. The price equals its value when
        wage = daily subsistence cost; surplus arises only in production (deferred to Ch. 7).
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Seller (worker)</span>
          <select value={sellerID} onChange={(e) => setSellerID(e.target.value)} required>
            <option value="">Select a worker…</option>
            {workers.map((w) => (
              <option key={w.id} value={w.id}>{w.id.slice(0, 8)}…</option>
            ))}
          </select>
        </label>
        <label>
          <span>Buyer (capitalist)</span>
          <select value={buyerID} onChange={(e) => setBuyerID(e.target.value)} required>
            <option value="">Select a capitalist…</option>
            {capitalists.map((c) => (
              <option key={c.id} value={c.id}>{c.id.slice(0, 8)}…</option>
            ))}
          </select>
        </label>
        <label>
          <span>Wage (labour-min/day)</span>
          <input
            type="number"
            value={wageMinutes}
            onChange={(e) => setWageMinutes(Number(e.target.value))}
            min={0}
          />
        </label>
        <label>
          <span>Contract (days)</span>
          <input
            type="number"
            value={contractDays}
            onChange={(e) => setContractDays(Number(e.target.value))}
            min={1}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Purchase</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
    </section>
  );
}

function PurchaseHistoryPanel({
  purchases,
  workers,
  capitalists,
}: {
  purchases: LabourPowerPurchase[];
  workers: LabourWorker[];
  capitalists: LabourCapitalist[];
}) {
  if (purchases.length === 0) return null;
  const workerMap = new Map(workers.map((w) => [w.id, w]));
  const capitalistMap = new Map(capitalists.map((c) => [c.id, c]));

  return (
    <section className="card">
      <h2>Purchase History</h2>
      <table className="data-table">
        <thead>
          <tr>
            <th>Seller</th>
            <th>Buyer</th>
            <th>Wage / day</th>
            <th>Days</th>
            <th>Total wage</th>
          </tr>
        </thead>
        <tbody>
          {purchases.map((p) => {
            const seller = workerMap.get(p.seller_id);
            const buyer = capitalistMap.get(p.buyer_id);
            return (
              <tr key={p.id}>
                <td className="monospace small">
                  {seller ? p.seller_id.slice(0, 8) + "…" : p.seller_id.slice(0, 8) + "…"}
                </td>
                <td className="monospace small">
                  {buyer ? p.buyer_id.slice(0, 8) + "…" : p.buyer_id.slice(0, 8) + "…"}
                </td>
                <td>{minutesToHours(p.wage_minutes)}</td>
                <td>{p.contract_days}</td>
                <td>{minutesToHours(p.wage_minutes * p.contract_days)}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </section>
  );
}
```

- [ ] **Step 2: Wire Ch06 into `web/src/components/ChapterShell.tsx`**

Add the import at the top:
```typescript
import { Ch06LabourPower } from "../chapters/Ch06LabourPower";
```

Add to the `QUOTES` map:
```typescript
ch06: "The owner of labour-power is mortal. If his appearance in the market is to be continuous, the seller of labour-power must perpetuate himself.",
```

Add the render case in `ChapterShell`:
```tsx
) : activeChapterId === "ch06" ? (
  <Ch06LabourPower onSharedChanged={onSharedChanged} />
) : null}
```

The full updated conditional chain in `ChapterShell`:
```tsx
{chapter.status === "pending" ? (
  <div className="chapter-placeholder">...</div>
) : activeChapterId === "ch01" ? (
  <Ch01Commodity commodities={commodities} onSharedChanged={onSharedChanged} />
) : activeChapterId === "ch02" ? (
  <Ch02Exchange commodities={commodities} owners={owners} onSharedChanged={onSharedChanged} />
) : activeChapterId === "ch03" ? (
  <Ch03Money owners={owners} onSharedChanged={onSharedChanged} />
) : activeChapterId === "ch04" ? (
  <Ch04Capital onSharedChanged={onSharedChanged} />
) : activeChapterId === "ch05" ? (
  <Ch05Contradictions onSharedChanged={onSharedChanged} />
) : activeChapterId === "ch06" ? (
  <Ch06LabourPower onSharedChanged={onSharedChanged} />
) : null}
```

- [ ] **Step 3: Update `registry.ts` — mark ch06 as done**

Change `ch06` status from `"pending"` to `"done"`:
```typescript
{ id: "ch06", number: 6, title: "The Sale and Purchase of Labour-Power", part: "Part II — The Transformation of Money into Capital", status: "done" },
```

- [ ] **Step 4: Ask user to typecheck and build**

Run: `cd web && npm run lint && npm run build`
Expected: no TypeScript errors, Vite production build succeeds.

---

## Task 9: Documentation

**Files:**
- Modify: `docs/architecture.md`

- [ ] **Step 1: Update the chapter table in `docs/architecture.md`**

In the roadmap table, change the Ch. 6 row from pending to done and split Ch. 6-7:
```markdown
| Ch. 6     | ✅ Done     | Labour-power as commodity; workers, capitalists, labour-power value, wage, subsistence basket; buying and selling of labour-power | agent-service |
| Ch. 7     | Pending     | Labour-process, valorization, surplus-value production | agent-service, simulation-eng |
```

Add a `### Ch. 6 — what was built` section after the Ch. 5 section with a concise summary of what was implemented.

- [ ] **Step 2: Final verification**

Run: `make vet test build && cd web && npm run lint && npm run build`
Expected: everything green.

---

## Self-review against spec

| Spec requirement | Task |
|---|---|
| `LabourMinutes` type | Task 2 |
| `AgentID` + `NewAgentID()` | Task 2 |
| `PurchaseID` + `NewPurchaseID()` | Task 2 |
| `AgentKind` enum + `AgentKindWorker`, `AgentKindCapitalist` | Task 2 |
| `LabourAgent` base (spec: `Agent`) | Task 2 — name changed due to conflict |
| `Worker` struct + `Validate()` | Task 2 |
| `Capitalist` struct + `Validate()` | Task 2 |
| `LabourPower` struct | Task 2 |
| `SubsistenceBasket` + `SubsistenceItem` + `TotalSNLT()` | Task 2 |
| `LabourPowerValue` + `DailyValue()` + `MinimumValue()` + `ReproductionCost()` | Task 2 |
| `LabourPowerOffering` + `Validate()` | Task 2 |
| `LabourPowerPurchase` + `Validate()` | Task 2 |
| `IsFreeLabourer(w Worker) bool` | Task 2 |
| `ErrInvalidContract` | Task 2 |
| All §1–§5 fixtures as tests | Task 2 |
| All invariants as tests | Task 2 |
| Store interface + Memory impl | Task 3 |
| MySQL migration + impl | Task 4 |
| `POST /v1/workers` | Task 5 |
| `GET /v1/workers/{id}` | Task 5 |
| `POST /v1/capitalists` | Task 5 |
| `GET /v1/capitalists/{id}` | Task 5 |
| `POST /v1/labour-power/offerings` | Task 5 |
| `GET /v1/labour-power/offerings` | Task 5 |
| `POST /v1/labour-power/purchases` | Task 5 |
| `GET /v1/labour-power/purchases/{id}` | Task 5 |
| API gateway proxy routes | Task 6 |
| `AgentPanel` / `Ch06LabourPower` React component | Task 8 |
| Worker + Capitalist in `types.ts` + `api.ts` | Task 7 |
| `LabourPowerOffering` + `LabourPowerPurchase` in `types.ts` + `api.ts` | Task 7 |
| `docs/architecture.md` updated | Task 9 |

Extra endpoints not in spec but added for UI: `GET /v1/workers`, `GET /v1/capitalists`, `GET /v1/labour-power/purchases` (list). These are natural complements needed to render the React panels.
