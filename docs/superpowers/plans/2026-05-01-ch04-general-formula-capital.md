# Ch. 04 — The General Formula for Capital: Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement agent-service (domain, store, HTTP) for Capital Vol. I Ch. 4 — the M—C—M′ circuit, class positions (Capitalist/Worker/Miser), and surplus-value observation.

**Architecture:** A new fully-functional `agent-service` replaces the stub. It owns `Agent` and `CapitalCircuit` records with MySQL storage and an in-memory fallback. The API gateway gains proxy routes for `/v1/agents`. The React UI gains a Ch. 04 panel.

**Tech Stack:** Go 1.25, `database/sql` + MySQL 8, `pressly/goose/v3` migrations, React 18 + TypeScript. Module: `github.com/theding0x/capital-simulator`.

---

## File Map

**New files:**
- `services/agent-service/internal/agent/agent.go` — domain types, methods, errors
- `services/agent-service/internal/agent/agent_test.go` — domain unit tests
- `services/agent-service/internal/store/store.go` — Store + CircuitStore interfaces, Update, sentinels
- `services/agent-service/internal/store/memory.go` — in-memory implementation
- `services/agent-service/internal/store/memory_test.go` — store tests (via Memory)
- `services/agent-service/internal/store/migrations/00001_ch04_agents.sql`
- `services/agent-service/internal/store/migrations/00002_ch04_circuits.sql`
- `services/agent-service/internal/store/mysql.go` — MySQL implementation
- `services/agent-service/internal/transport/httpapi/handler.go` — HTTP handlers
- `services/agent-service/internal/transport/httpapi/routes.go` — route registration
- `web/src/chapters/Ch04Capital.tsx` — chapter 04 UI panel

**Modified files:**
- `services/agent-service/internal/agent/agent.go` — replace placeholder comment
- `services/agent-service/cmd/agent-service/main.go` — replace stub with full initialization
- `services/api-gateway/cmd/api-gateway/main.go` — add agent proxy routes
- `web/src/types.ts` — add Agent, CapitalCircuit, input types
- `web/src/api.ts` — add agent-service API methods
- `web/src/components/ChapterShell.tsx` — wire Ch04Capital
- `web/src/chapters/registry.ts` — mark ch04 as "done"
- `docs/architecture.md` — update chapter status table

---

## Task 1: Domain types, methods, and errors

**Files:**
- Create: `services/agent-service/internal/agent/agent.go`

- [ ] **Step 1: Write the domain package**

Replace the placeholder `services/agent-service/internal/agent/agent.go` with:

```go
package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Named scalar types — never use raw int64/string for money or IDs.
type ID string
type Class string
type Pence int64
type CircuitType string

const (
	Capitalist Class = "capitalist"
	Worker     Class = "worker"
	Miser      Class = "miser"
)

const (
	CircuitCMC CircuitType = "C-M-C"
	CircuitMCM CircuitType = "M-C-M-prime"
)

var (
	ErrInsufficientFunds = errors.New("agent: insufficient funds")
	ErrNotCapitalist     = errors.New("agent: operation requires capitalist class")
	ErrWrongClass        = errors.New("agent: operation not permitted for this class")
)

// Agent is the bearer of a class relation with a money balance.
type Agent struct {
	ID           ID        `json:"id"`
	Name         string    `json:"name"`
	Class        Class     `json:"class"`
	MoneyBalance Pence     `json:"money_balance"`
	Hoarding     bool      `json:"hoarding"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CapitalCircuit records a single execution of M—C—M′.
// SurplusValue is always computed as MReturned - MAdvanced; never stored independently.
type CapitalCircuit struct {
	ID           ID          `json:"id"`
	AgentID      ID          `json:"agent_id"`
	MAdvanced    Pence       `json:"m_advanced"`
	CommodityID  string      `json:"commodity_id"`
	MReturned    Pence       `json:"m_returned"`
	SurplusValue Pence       `json:"surplus_value"`
	CircuitType  CircuitType `json:"circuit_type"`
	CreatedAt    time.Time   `json:"created_at"`
}

func (id ID) IsZero() bool { return id == "" }

// NewID returns a 96-bit random hex ID (24 chars).
func NewID() ID {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return ID(hex.EncodeToString(b))
}

func (a Agent) Validate() error {
	if a.Name == "" {
		return errors.New("agent: name is required")
	}
	switch a.Class {
	case Capitalist, Worker, Miser:
	default:
		return errors.New("agent: unknown class")
	}
	if a.MoneyBalance < 0 {
		return errors.New("agent: money_balance cannot be negative")
	}
	return nil
}

func (c CapitalCircuit) Validate() error {
	if c.AgentID.IsZero() {
		return errors.New("circuit: agent_id is required")
	}
	if c.MAdvanced <= 0 {
		return errors.New("circuit: m_advanced must be positive")
	}
	if c.CommodityID == "" {
		return errors.New("circuit: commodity_id is required")
	}
	switch c.CircuitType {
	case CircuitCMC, CircuitMCM:
	default:
		return errors.New("circuit: unknown circuit_type")
	}
	if c.SurplusValue != c.MReturned-c.MAdvanced {
		return errors.New("circuit: surplus_value must equal m_returned - m_advanced")
	}
	return nil
}

// Advance deducts mAdvanced from MoneyBalance. Returns ErrNotCapitalist for
// Miser agents and ErrInsufficientFunds if mAdvanced > MoneyBalance.
func (a Agent) Advance(mAdvanced Pence) (Agent, error) {
	if a.Class == Miser {
		return Agent{}, ErrNotCapitalist
	}
	if mAdvanced > a.MoneyBalance {
		return Agent{}, ErrInsufficientFunds
	}
	out := a
	out.MoneyBalance -= mAdvanced
	return out, nil
}

// Realise credits circuit.MReturned to MoneyBalance.
func (a Agent) Realise(circuit CapitalCircuit) Agent {
	out := a
	out.MoneyBalance += circuit.MReturned
	return out
}

// Reinvest creates a new M-C-M' circuit using the agent's full balance as
// MAdvanced. Valid only for Capitalist agents.
func (a Agent) Reinvest(commodityID string, mReturned Pence) (CapitalCircuit, Agent, error) {
	if a.Class != Capitalist {
		return CapitalCircuit{}, Agent{}, ErrNotCapitalist
	}
	if a.MoneyBalance <= 0 {
		return CapitalCircuit{}, Agent{}, ErrInsufficientFunds
	}
	circuit := CapitalCircuit{
		ID:           NewID(),
		AgentID:      a.ID,
		MAdvanced:    a.MoneyBalance,
		CommodityID:  commodityID,
		MReturned:    mReturned,
		SurplusValue: mReturned - a.MoneyBalance,
		CircuitType:  CircuitMCM,
		CreatedAt:    time.Now().UTC(),
	}
	if err := circuit.Validate(); err != nil {
		return CapitalCircuit{}, Agent{}, err
	}
	advanced, err := a.Advance(circuit.MAdvanced)
	if err != nil {
		return CapitalCircuit{}, Agent{}, err
	}
	realised := advanced.Realise(circuit)
	return circuit, realised, nil
}

// Hoard sets the Hoarding flag on a Miser agent. Returns ErrNotCapitalist for
// non-Miser agents. Balance is unaffected (idempotent no-op on balance).
func (a Agent) Hoard() (Agent, error) {
	if a.Class != Miser {
		return Agent{}, ErrNotCapitalist
	}
	out := a
	out.Hoarding = true
	return out, nil
}
```

- [ ] **Step 2: Write the domain tests**

Create `services/agent-service/internal/agent/agent_test.go`:

```go
package agent_test

import (
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

func newCapitalist(balance agent.Pence) agent.Agent {
	return agent.Agent{ID: agent.NewID(), Name: "Capitalist", Class: agent.Capitalist, MoneyBalance: balance}
}

func newMiser(balance agent.Pence) agent.Agent {
	return agent.Agent{ID: agent.NewID(), Name: "Miser", Class: agent.Miser, MoneyBalance: balance}
}

func newWorker(balance agent.Pence) agent.Agent {
	return agent.Agent{ID: agent.NewID(), Name: "Worker", Class: agent.Worker, MoneyBalance: balance}
}

// §1: Advance deducts balance to zero.
func TestAdvance_§1_DeductsToZero(t *testing.T) {
	t.Parallel()
	a := newCapitalist(10000)
	got, err := a.Advance(10000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.MoneyBalance != 0 {
		t.Errorf("want balance 0, got %d", got.MoneyBalance)
	}
}

// §8: SurplusValue == MReturned - MAdvanced.
func TestCapitalCircuit_§8_SurplusValue(t *testing.T) {
	t.Parallel()
	c := agent.CapitalCircuit{
		ID:           agent.NewID(),
		AgentID:      agent.NewID(),
		MAdvanced:    10000,
		CommodityID:  "cotton",
		MReturned:    11000,
		SurplusValue: 1000,
		CircuitType:  agent.CircuitMCM,
	}
	if c.SurplusValue != c.MReturned-c.MAdvanced {
		t.Errorf("want surplus %d, got %d", c.MReturned-c.MAdvanced, c.SurplusValue)
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("valid circuit rejected: %v", err)
	}
}

// §10: Zero surplus circuit is valid but signals no valorisation.
func TestCapitalCircuit_§10_ZeroSurplus(t *testing.T) {
	t.Parallel()
	c := agent.CapitalCircuit{
		ID:           agent.NewID(),
		AgentID:      agent.NewID(),
		MAdvanced:    10000,
		CommodityID:  "cotton",
		MReturned:    10000,
		SurplusValue: 0,
		CircuitType:  agent.CircuitMCM,
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("zero-surplus circuit should be valid: %v", err)
	}
}

// §14: Miser.Hoard() succeeds; balance unchanged; Hoarding flag set.
func TestHoard_§14_MiserSucceeds(t *testing.T) {
	t.Parallel()
	a := newMiser(10000)
	got, err := a.Hoard()
	if err != nil {
		t.Fatalf("Miser.Hoard() error: %v", err)
	}
	if !got.Hoarding {
		t.Error("want Hoarding=true")
	}
	if got.MoneyBalance != a.MoneyBalance {
		t.Errorf("Hoard must not change balance; want %d got %d", a.MoneyBalance, got.MoneyBalance)
	}
}

// §14: Miser.Reinvest() → ErrNotCapitalist.
func TestReinvest_§14_MiserFails(t *testing.T) {
	t.Parallel()
	a := newMiser(10000)
	_, _, err := a.Reinvest("cotton", 11000)
	if !errors.Is(err, agent.ErrNotCapitalist) {
		t.Errorf("Miser.Reinvest() want ErrNotCapitalist, got %v", err)
	}
}

// §14: Capitalist.Reinvest() succeeds; circuit.MAdvanced = old balance.
func TestReinvest_§14_CapitalistSucceeds(t *testing.T) {
	t.Parallel()
	a := newCapitalist(10000)
	circuit, updated, err := a.Reinvest("cotton", 11000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if circuit.MAdvanced != 10000 {
		t.Errorf("want MAdvanced 10000, got %d", circuit.MAdvanced)
	}
	if circuit.SurplusValue != 1000 {
		t.Errorf("want SurplusValue 1000, got %d", circuit.SurplusValue)
	}
	if updated.MoneyBalance != 11000 {
		t.Errorf("want balance 11000 after reinvest, got %d", updated.MoneyBalance)
	}
}

// §14: Capitalist.Hoard() → ErrNotCapitalist.
func TestHoard_§14_CapitalistFails(t *testing.T) {
	t.Parallel()
	a := newCapitalist(10000)
	_, err := a.Hoard()
	if !errors.Is(err, agent.ErrNotCapitalist) {
		t.Errorf("Capitalist.Hoard() want ErrNotCapitalist, got %v", err)
	}
}

// §15: After Realise, second circuit uses full new balance as MAdvanced.
func TestRealise_§15_ExpandingCircuit(t *testing.T) {
	t.Parallel()
	a := newCapitalist(10000)
	c1 := agent.CapitalCircuit{
		ID: agent.NewID(), AgentID: a.ID,
		MAdvanced: 10000, CommodityID: "cotton", MReturned: 11000,
		SurplusValue: 1000, CircuitType: agent.CircuitMCM,
	}
	advanced, _ := a.Advance(c1.MAdvanced)
	realised := advanced.Realise(c1)
	if realised.MoneyBalance != 11000 {
		t.Fatalf("want balance 11000 after Realise, got %d", realised.MoneyBalance)
	}
	c2, _, err := realised.Reinvest("cotton", 12100)
	if err != nil {
		t.Fatalf("second Reinvest error: %v", err)
	}
	if c2.MAdvanced != 11000 {
		t.Errorf("second circuit want MAdvanced 11000 (full new balance), got %d", c2.MAdvanced)
	}
}

// Invariant: Advance returns ErrInsufficientFunds when mAdvanced > balance.
func TestAdvance_InsufficientFunds(t *testing.T) {
	t.Parallel()
	a := newCapitalist(5000)
	_, err := a.Advance(10000)
	if !errors.Is(err, agent.ErrInsufficientFunds) {
		t.Errorf("want ErrInsufficientFunds, got %v", err)
	}
}

// Invariant: Miser cannot Advance.
func TestAdvance_MiserCannotAdvance(t *testing.T) {
	t.Parallel()
	a := newMiser(10000)
	_, err := a.Advance(1000)
	if !errors.Is(err, agent.ErrNotCapitalist) {
		t.Errorf("Miser.Advance() want ErrNotCapitalist, got %v", err)
	}
}

// §6: CircuitCMC is valid for any agent.
func TestCircuitValidate_§6_CMCIsValid(t *testing.T) {
	t.Parallel()
	c := agent.CapitalCircuit{
		ID: agent.NewID(), AgentID: agent.NewID(),
		MAdvanced: 5000, CommodityID: "bread", MReturned: 5000,
		SurplusValue: 0, CircuitType: agent.CircuitCMC,
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("C-M-C circuit should be valid: %v", err)
	}
}

// Invariant: circuit.MAdvanced <= 0 is invalid.
func TestCircuitValidate_ZeroAdvanced(t *testing.T) {
	t.Parallel()
	c := agent.CapitalCircuit{
		ID: agent.NewID(), AgentID: agent.NewID(),
		MAdvanced: 0, CommodityID: "cotton", MReturned: 0,
		SurplusValue: 0, CircuitType: agent.CircuitMCM,
	}
	if err := c.Validate(); err == nil {
		t.Error("circuit with MAdvanced=0 should fail validation")
	}
}
```

- [ ] **Step 3: Ask user to run tests**

```bash
make vet test build
```
Expected: all agent package tests pass.

---

## Task 2: Store interfaces and errors

**Files:**
- Create: `services/agent-service/internal/store/store.go`

- [ ] **Step 1: Write store.go**

```go
package store

import (
	"context"
	"errors"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

var (
	ErrNotFound      = errors.New("agent: not found")
	ErrAlreadyExists = errors.New("agent: already exists")
)

// Update is a partial-update payload for PATCH /v1/agents/{id}.
// Non-nil fields are applied; nil fields are left untouched.
type Update struct {
	Name         *string
	MoneyBalance *agent.Pence
	Hoarding     *bool
}

func (u Update) IsEmpty() bool {
	return u.Name == nil && u.MoneyBalance == nil && u.Hoarding == nil
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
	return out
}

// Store is the persistence contract for Agent records.
type Store interface {
	Create(ctx context.Context, a agent.Agent) (agent.Agent, error)
	Get(ctx context.Context, id agent.ID) (agent.Agent, error)
	List(ctx context.Context) ([]agent.Agent, error)
	ListByClass(ctx context.Context, class agent.Class) ([]agent.Agent, error)
	Update(ctx context.Context, id agent.ID, u Update) (agent.Agent, error)
	Delete(ctx context.Context, id agent.ID) error
}

// CircuitStore is the persistence contract for CapitalCircuit records.
// CreateCircuit atomically inserts the circuit and updates the agent balance
// by circuit.SurplusValue.
type CircuitStore interface {
	CreateCircuit(ctx context.Context, c agent.CapitalCircuit) (agent.CapitalCircuit, error)
	GetCircuit(ctx context.Context, id agent.ID) (agent.CapitalCircuit, error)
	ListCircuits(ctx context.Context, agentID agent.ID) ([]agent.CapitalCircuit, error)
}
```

---

## Task 3: In-memory store implementation and tests

**Files:**
- Create: `services/agent-service/internal/store/memory.go`
- Create: `services/agent-service/internal/store/memory_test.go`

- [ ] **Step 1: Write memory.go**

```go
package store

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// Memory implements Store and CircuitStore for tests and local dev.
type Memory struct {
	mu       sync.RWMutex
	agents   map[agent.ID]agent.Agent
	circuits map[agent.ID]agent.CapitalCircuit
	now      func() time.Time
}

func NewMemory() *Memory {
	return &Memory{
		agents:   make(map[agent.ID]agent.Agent),
		circuits: make(map[agent.ID]agent.CapitalCircuit),
		now:      time.Now,
	}
}

func (m *Memory) Create(_ context.Context, a agent.Agent) (agent.Agent, error) {
	if err := a.Validate(); err != nil {
		return agent.Agent{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if a.ID.IsZero() {
		a.ID = agent.NewID()
	}
	now := m.now()
	a.CreatedAt = now
	a.UpdatedAt = now
	m.agents[a.ID] = a
	return a, nil
}

func (m *Memory) Get(_ context.Context, id agent.ID) (agent.Agent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	a, ok := m.agents[id]
	if !ok {
		return agent.Agent{}, ErrNotFound
	}
	return a, nil
}

func (m *Memory) List(_ context.Context) ([]agent.Agent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]agent.Agent, 0, len(m.agents))
	for _, a := range m.agents {
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

func (m *Memory) ListByClass(_ context.Context, class agent.Class) ([]agent.Agent, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []agent.Agent
	for _, a := range m.agents {
		if a.Class == class {
			out = append(out, a)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

func (m *Memory) Update(_ context.Context, id agent.ID, u Update) (agent.Agent, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cur, ok := m.agents[id]
	if !ok {
		return agent.Agent{}, ErrNotFound
	}
	next := u.Apply(cur)
	if err := next.Validate(); err != nil {
		return agent.Agent{}, err
	}
	next.UpdatedAt = m.now()
	m.agents[id] = next
	return next, nil
}

func (m *Memory) Delete(_ context.Context, id agent.ID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.agents[id]; !ok {
		return ErrNotFound
	}
	delete(m.agents, id)
	return nil
}

// CreateCircuit atomically inserts the circuit and updates the agent's
// money_balance by circuit.SurplusValue. Returns ErrNotFound if the agent
// doesn't exist; returns agent.ErrInsufficientFunds if balance < MAdvanced.
func (m *Memory) CreateCircuit(_ context.Context, c agent.CapitalCircuit) (agent.CapitalCircuit, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	a, ok := m.agents[c.AgentID]
	if !ok {
		return agent.CapitalCircuit{}, ErrNotFound
	}
	if c.MAdvanced > a.MoneyBalance {
		return agent.CapitalCircuit{}, agent.ErrInsufficientFunds
	}
	if c.ID.IsZero() {
		c.ID = agent.NewID()
	}
	c.CreatedAt = m.now()
	m.circuits[c.ID] = c
	a.MoneyBalance += c.SurplusValue
	a.UpdatedAt = m.now()
	m.agents[a.ID] = a
	return c, nil
}

func (m *Memory) GetCircuit(_ context.Context, id agent.ID) (agent.CapitalCircuit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	c, ok := m.circuits[id]
	if !ok {
		return agent.CapitalCircuit{}, ErrNotFound
	}
	return c, nil
}

func (m *Memory) ListCircuits(_ context.Context, agentID agent.ID) ([]agent.CapitalCircuit, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var out []agent.CapitalCircuit
	for _, c := range m.circuits {
		if c.AgentID == agentID {
			out = append(out, c)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}
```

- [ ] **Step 2: Write memory_test.go**

```go
package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

func makeAgent(class agent.Class, balance agent.Pence) agent.Agent {
	return agent.Agent{Name: "Test Agent", Class: class, MoneyBalance: balance}
}

func TestMemory_CreateGet(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	created, err := m.Create(ctx, makeAgent(agent.Capitalist, 10000))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID.IsZero() {
		t.Error("Create should assign an ID")
	}
	got, err := m.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != created.Name {
		t.Errorf("want name %q, got %q", created.Name, got.Name)
	}
}

func TestMemory_Get_NotFound(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.Get(context.Background(), agent.NewID())
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestMemory_ListByClass(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	_, _ = m.Create(ctx, makeAgent(agent.Capitalist, 10000))
	_, _ = m.Create(ctx, makeAgent(agent.Worker, 500))
	_, _ = m.Create(ctx, agent.Agent{Name: "Worker2", Class: agent.Worker, MoneyBalance: 300})
	caps, _ := m.ListByClass(ctx, agent.Capitalist)
	if len(caps) != 1 {
		t.Errorf("want 1 capitalist, got %d", len(caps))
	}
	workers, _ := m.ListByClass(ctx, agent.Worker)
	if len(workers) != 2 {
		t.Errorf("want 2 workers, got %d", len(workers))
	}
}

func TestMemory_Update(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	a, _ := m.Create(ctx, makeAgent(agent.Capitalist, 10000))
	newName := "Updated"
	updated, err := m.Update(ctx, a.ID, store.Update{Name: &newName})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != newName {
		t.Errorf("want name %q, got %q", newName, updated.Name)
	}
}

func TestMemory_Delete(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	a, _ := m.Create(ctx, makeAgent(agent.Capitalist, 10000))
	if err := m.Delete(ctx, a.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := m.Delete(ctx, a.ID); !errors.Is(err, store.ErrNotFound) {
		t.Errorf("second Delete want ErrNotFound, got %v", err)
	}
}

// §8: CreateCircuit computes balance update from SurplusValue.
func TestMemory_CreateCircuit_UpdatesBalance(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	a, _ := m.Create(ctx, makeAgent(agent.Capitalist, 10000))
	c := agent.CapitalCircuit{
		AgentID:      a.ID,
		MAdvanced:    10000,
		CommodityID:  "cotton",
		MReturned:    11000,
		SurplusValue: 1000,
		CircuitType:  agent.CircuitMCM,
	}
	saved, err := m.CreateCircuit(ctx, c)
	if err != nil {
		t.Fatalf("CreateCircuit: %v", err)
	}
	if saved.ID.IsZero() {
		t.Error("CreateCircuit should assign an ID")
	}
	updated, _ := m.Get(ctx, a.ID)
	if updated.MoneyBalance != 11000 {
		t.Errorf("want balance 11000 after circuit, got %d", updated.MoneyBalance)
	}
}

// Invariant: CreateCircuit returns ErrInsufficientFunds when balance < MAdvanced.
func TestMemory_CreateCircuit_InsufficientFunds(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	a, _ := m.Create(ctx, makeAgent(agent.Capitalist, 5000))
	c := agent.CapitalCircuit{
		AgentID: a.ID, MAdvanced: 10000, CommodityID: "cotton",
		MReturned: 11000, SurplusValue: 1000, CircuitType: agent.CircuitMCM,
	}
	_, err := m.CreateCircuit(ctx, c)
	if !errors.Is(err, agent.ErrInsufficientFunds) {
		t.Errorf("want ErrInsufficientFunds, got %v", err)
	}
}

func TestMemory_ListCircuits(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	a, _ := m.Create(ctx, makeAgent(agent.Capitalist, 30000))
	for i := 0; i < 3; i++ {
		_, _ = m.CreateCircuit(ctx, agent.CapitalCircuit{
			AgentID: a.ID, MAdvanced: 1000, CommodityID: "cotton",
			MReturned: 1100, SurplusValue: 100, CircuitType: agent.CircuitMCM,
		})
	}
	cs, err := m.ListCircuits(ctx, a.ID)
	if err != nil {
		t.Fatalf("ListCircuits: %v", err)
	}
	if len(cs) != 3 {
		t.Errorf("want 3 circuits, got %d", len(cs))
	}
}
```

- [ ] **Step 3: Ask user to run tests**

```bash
make vet test build
```
Expected: all store tests pass.

---

## Task 4: SQL migrations

**Files:**
- Create: `services/agent-service/internal/store/migrations/00001_ch04_agents.sql`
- Create: `services/agent-service/internal/store/migrations/00002_ch04_circuits.sql`

- [ ] **Step 1: Write agents migration**

`services/agent-service/internal/store/migrations/00001_ch04_agents.sql`:

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS agents (
    id            VARCHAR(24)  NOT NULL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    class         VARCHAR(50)  NOT NULL,
    money_balance BIGINT       NOT NULL DEFAULT 0,
    hoarding      TINYINT(1)   NOT NULL DEFAULT 0,
    created_at    DATETIME(6)  NOT NULL,
    updated_at    DATETIME(6)  NOT NULL,
    INDEX idx_class (class)
);

-- +goose Down
DROP TABLE IF EXISTS agents;
```

- [ ] **Step 2: Write circuits migration**

`services/agent-service/internal/store/migrations/00002_ch04_circuits.sql`:

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS capital_circuits (
    id            VARCHAR(24)  NOT NULL PRIMARY KEY,
    agent_id      VARCHAR(24)  NOT NULL,
    m_advanced    BIGINT       NOT NULL,
    commodity_id  VARCHAR(255) NOT NULL,
    m_returned    BIGINT       NOT NULL DEFAULT 0,
    surplus_value BIGINT       NOT NULL DEFAULT 0,
    circuit_type  VARCHAR(20)  NOT NULL,
    created_at    DATETIME(6)  NOT NULL,
    INDEX idx_agent_id (agent_id)
);

-- +goose Down
DROP TABLE IF EXISTS capital_circuits;
```

---

## Task 5: MySQL store implementation

**Files:**
- Create: `services/agent-service/internal/store/mysql.go`

- [ ] **Step 1: Write mysql.go**

```go
package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"io/fs"
	"strings"
	"time"

	pkgmysql "github.com/theding0x/capital-simulator/pkg/mysql"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

//go:embed migrations
var migrationsFS embed.FS

// MySQL implements Store and CircuitStore backed by MySQL.
type MySQL struct {
	db  *sql.DB
	now func() time.Time
}

func NewMySQL(ctx context.Context, db *sql.DB) (*MySQL, error) {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return nil, err
	}
	if err := pkgmysql.Migrate(ctx, db, sub); err != nil {
		return nil, err
	}
	return &MySQL{db: db, now: time.Now}, nil
}

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
	const q = `INSERT INTO agents (id, name, class, money_balance, hoarding, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID), a.Name, string(a.Class), int64(a.MoneyBalance), a.Hoarding,
		a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return agent.Agent{}, err
	}
	return a, nil
}

func (m *MySQL) Get(ctx context.Context, id agent.ID) (agent.Agent, error) {
	const q = `SELECT id, name, class, money_balance, hoarding, created_at, updated_at
		FROM agents WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanAgent(row)
}

func (m *MySQL) List(ctx context.Context) ([]agent.Agent, error) {
	const q = `SELECT id, name, class, money_balance, hoarding, created_at, updated_at
		FROM agents ORDER BY name ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAgents(rows)
}

func (m *MySQL) ListByClass(ctx context.Context, class agent.Class) ([]agent.Agent, error) {
	const q = `SELECT id, name, class, money_balance, hoarding, created_at, updated_at
		FROM agents WHERE class = ? ORDER BY name ASC`
	rows, err := m.db.QueryContext(ctx, q, string(class))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAgents(rows)
}

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
	const q = `UPDATE agents SET name = ?, class = ?, money_balance = ?, hoarding = ?, updated_at = ?
		WHERE id = ?`
	res, err := m.db.ExecContext(ctx, q,
		next.Name, string(next.Class), int64(next.MoneyBalance), next.Hoarding, next.UpdatedAt,
		string(id),
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

func (m *MySQL) Delete(ctx context.Context, id agent.ID) error {
	res, err := m.db.ExecContext(ctx, `DELETE FROM agents WHERE id = ?`, string(id))
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// CreateCircuit atomically inserts the circuit and updates the agent's
// money_balance by circuit.SurplusValue in a single transaction.
func (m *MySQL) CreateCircuit(ctx context.Context, c agent.CapitalCircuit) (agent.CapitalCircuit, error) {
	if c.ID.IsZero() {
		c.ID = agent.NewID()
	}
	c.CreatedAt = m.now().UTC()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return agent.CapitalCircuit{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var bal int64
	err = tx.QueryRowContext(ctx,
		`SELECT money_balance FROM agents WHERE id = ? FOR UPDATE`, string(c.AgentID),
	).Scan(&bal)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.CapitalCircuit{}, ErrNotFound
	}
	if err != nil {
		return agent.CapitalCircuit{}, err
	}
	if agent.Pence(bal) < c.MAdvanced {
		return agent.CapitalCircuit{}, agent.ErrInsufficientFunds
	}

	const insertQ = `INSERT INTO capital_circuits
		(id, agent_id, m_advanced, commodity_id, m_returned, surplus_value, circuit_type, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, insertQ,
		string(c.ID), string(c.AgentID), int64(c.MAdvanced),
		c.CommodityID, int64(c.MReturned), int64(c.SurplusValue),
		string(c.CircuitType), c.CreatedAt,
	)
	if err != nil {
		return agent.CapitalCircuit{}, err
	}

	now := m.now().UTC()
	_, err = tx.ExecContext(ctx,
		`UPDATE agents SET money_balance = money_balance + ?, updated_at = ? WHERE id = ?`,
		int64(c.SurplusValue), now, string(c.AgentID),
	)
	if err != nil {
		return agent.CapitalCircuit{}, err
	}
	if err := tx.Commit(); err != nil {
		return agent.CapitalCircuit{}, err
	}
	return c, nil
}

func (m *MySQL) GetCircuit(ctx context.Context, id agent.ID) (agent.CapitalCircuit, error) {
	const q = `SELECT id, agent_id, m_advanced, commodity_id, m_returned, surplus_value, circuit_type, created_at
		FROM capital_circuits WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanCircuit(row)
}

func (m *MySQL) ListCircuits(ctx context.Context, agentID agent.ID) ([]agent.CapitalCircuit, error) {
	const q = `SELECT id, agent_id, m_advanced, commodity_id, m_returned, surplus_value, circuit_type, created_at
		FROM capital_circuits WHERE agent_id = ? ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(agentID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []agent.CapitalCircuit
	for rows.Next() {
		c, err := scanCircuitRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func scanAgent(row *sql.Row) (agent.Agent, error) {
	var a agent.Agent
	var id, class string
	var balance int64
	var hoarding bool
	err := row.Scan(&id, &a.Name, &class, &balance, &hoarding, &a.CreatedAt, &a.UpdatedAt)
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

func scanAgents(rows *sql.Rows) ([]agent.Agent, error) {
	var out []agent.Agent
	for rows.Next() {
		var a agent.Agent
		var id, class string
		var balance int64
		var hoarding bool
		if err := rows.Scan(&id, &a.Name, &class, &balance, &hoarding, &a.CreatedAt, &a.UpdatedAt); err != nil {
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

func scanCircuit(row *sql.Row) (agent.CapitalCircuit, error) {
	var c agent.CapitalCircuit
	var id, agentID, circuitType string
	var mAdv, mRet, sv int64
	err := row.Scan(&id, &agentID, &mAdv, &c.CommodityID, &mRet, &sv, &circuitType, &c.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.CapitalCircuit{}, ErrNotFound
	}
	if err != nil {
		return agent.CapitalCircuit{}, err
	}
	c.ID = agent.ID(id)
	c.AgentID = agent.ID(agentID)
	c.MAdvanced = agent.Pence(mAdv)
	c.MReturned = agent.Pence(mRet)
	c.SurplusValue = agent.Pence(sv)
	c.CircuitType = agent.CircuitType(circuitType)
	return c, nil
}

func scanCircuitRow(rows *sql.Rows) (agent.CapitalCircuit, error) {
	var c agent.CapitalCircuit
	var id, agentID, circuitType string
	var mAdv, mRet, sv int64
	if err := rows.Scan(&id, &agentID, &mAdv, &c.CommodityID, &mRet, &sv, &circuitType, &c.CreatedAt); err != nil {
		return agent.CapitalCircuit{}, err
	}
	c.ID = agent.ID(id)
	c.AgentID = agent.ID(agentID)
	c.MAdvanced = agent.Pence(mAdv)
	c.MReturned = agent.Pence(mRet)
	c.SurplusValue = agent.Pence(sv)
	c.CircuitType = agent.CircuitType(circuitType)
	return c, nil
}

func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "1062") || strings.Contains(s, "Duplicate entry")
}
```

- [ ] **Step 2: Ask user to run tests**

```bash
make vet test build
```
Expected: vet and build pass; store tests pass.

---

## Task 6: HTTP handler and routes

**Files:**
- Create: `services/agent-service/internal/transport/httpapi/handler.go`
- Create: `services/agent-service/internal/transport/httpapi/routes.go`

- [ ] **Step 1: Write handler.go**

```go
package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

type Handler struct {
	Store        store.Store
	CircuitStore store.CircuitStore
	Logger       *slog.Logger
}

func New(s store.Store, cs store.CircuitStore, logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{Store: s, CircuitStore: cs, Logger: logger}
}

// --- request / response types ---

type createAgentRequest struct {
	Name         string      `json:"name"`
	Class        agent.Class `json:"class"`
	MoneyBalance agent.Pence `json:"money_balance"`
}

type updateAgentRequest struct {
	Name         *string      `json:"name,omitempty"`
	MoneyBalance *agent.Pence `json:"money_balance,omitempty"`
}

type createCircuitRequest struct {
	MAdvanced   agent.Pence       `json:"m_advanced"`
	CommodityID string            `json:"commodity_id"`
	MReturned   agent.Pence       `json:"m_returned"`
	CircuitType agent.CircuitType `json:"circuit_type"`
}

type reinvestRequest struct {
	CommodityID string      `json:"commodity_id"`
	MReturned   agent.Pence `json:"m_returned"`
}

// --- handlers ---

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAgentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a := agent.Agent{
		Name:         strings.TrimSpace(req.Name),
		Class:        req.Class,
		MoneyBalance: req.MoneyBalance,
	}
	if err := a.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.Store.Create(r.Context(), a)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	classParam := r.URL.Query().Get("class")
	var (
		agents []agent.Agent
		err    error
	)
	if classParam != "" {
		agents, err = h.Store.ListByClass(ctx, agent.Class(classParam))
	} else {
		agents, err = h.Store.List(ctx)
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if agents == nil {
		agents = []agent.Agent{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": agents})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := agent.ID(r.PathValue("id"))
	a, err := h.Store.Get(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, a)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := agent.ID(r.PathValue("id"))
	var req updateAgentRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	updated, err := h.Store.Update(r.Context(), id, store.Update{
		Name:         req.Name,
		MoneyBalance: req.MoneyBalance,
	})
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := agent.ID(r.PathValue("id"))
	if err := h.Store.Delete(r.Context(), id); err != nil {
		writeAppError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateCircuit(w http.ResponseWriter, r *http.Request) {
	agentID := agent.ID(r.PathValue("id"))
	var req createCircuitRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a, err := h.Store.Get(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	if a.Class == agent.Worker && req.CircuitType == agent.CircuitMCM {
		writeError(w, http.StatusBadRequest, agent.ErrWrongClass.Error())
		return
	}
	c := agent.CapitalCircuit{
		AgentID:      agentID,
		MAdvanced:    req.MAdvanced,
		CommodityID:  req.CommodityID,
		MReturned:    req.MReturned,
		SurplusValue: req.MReturned - req.MAdvanced,
		CircuitType:  req.CircuitType,
	}
	if err := c.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.CircuitStore.CreateCircuit(r.Context(), c)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, saved)
}

func (h *Handler) ListCircuits(w http.ResponseWriter, r *http.Request) {
	agentID := agent.ID(r.PathValue("id"))
	circuits, err := h.CircuitStore.ListCircuits(r.Context(), agentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if circuits == nil {
		circuits = []agent.CapitalCircuit{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": circuits})
}

func (h *Handler) Reinvest(w http.ResponseWriter, r *http.Request) {
	agentID := agent.ID(r.PathValue("id"))
	var req reinvestRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	a, err := h.Store.Get(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	circuit, _, err := a.Reinvest(req.CommodityID, req.MReturned)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.CircuitStore.CreateCircuit(r.Context(), circuit)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, saved)
}

func (h *Handler) Hoard(w http.ResponseWriter, r *http.Request) {
	agentID := agent.ID(r.PathValue("id"))
	a, err := h.Store.Get(r.Context(), agentID)
	if err != nil {
		writeAppError(w, err)
		return
	}
	updated, err := a.Hoard()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	hoarding := true
	saved, err := h.Store.Update(r.Context(), agentID, store.Update{Hoarding: &hoarding})
	if err != nil {
		writeAppError(w, err)
		return
	}
	_ = updated
	writeJSON(w, http.StatusOK, saved)
}

// --- helpers ---

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return errors.New("invalid json: " + err.Error())
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeAppError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, store.ErrAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, agent.ErrInsufficientFunds),
		errors.Is(err, agent.ErrNotCapitalist),
		errors.Is(err, agent.ErrWrongClass):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}
```

- [ ] **Step 2: Write routes.go**

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
}
```

- [ ] **Step 3: Ask user to run tests**

```bash
make vet test build
```
Expected: vet and build pass.

---

## Task 7: Replace agent-service main.go

**Files:**
- Modify: `services/agent-service/cmd/agent-service/main.go`

- [ ] **Step 1: Replace the stub main.go**

```go
package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/theding0x/capital-simulator/pkg/httpx"
	applog "github.com/theding0x/capital-simulator/pkg/log"
	pmysql "github.com/theding0x/capital-simulator/pkg/mysql"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/transport/httpapi"
)

const serviceName = "agent-service"

func main() {
	logger := applog.New(serviceName)
	applog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st, mysqlDB, err := openStore(ctx, logger)
	if err != nil {
		logger.Error("could not open any store", "err", err)
		os.Exit(1)
	}
	if mysqlDB != nil {
		defer func() { _ = mysqlDB.Close() }()
	}

	addr := getenv("SERVICE_ADDR", ":8082")
	srv := httpx.New(httpx.Config{Addr: addr}, logger)

	httpapi.Register(srv, httpapi.New(st, st, logger))
	srv.MarkReady(true)

	if err := srv.Run(ctx); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func openStore(ctx context.Context, logger *slog.Logger) (*store.MySQL, *pmysql.DB, error) {
	if strings.EqualFold(os.Getenv("MYSQL_DISABLED"), "true") {
		logger.Warn("MYSQL_DISABLED=true; using in-memory store")
		return nil, nil, nil
	}
	cfg := pmysql.ConfigFromEnv(serviceName)
	dialCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()
	cli, err := pmysql.Connect(dialCtx, cfg)
	if err != nil {
		if strings.EqualFold(os.Getenv("FALLBACK_MEMORY"), "true") {
			logger.Warn("mysql connect failed; falling back to in-memory store", "err", err)
			return nil, nil, nil
		}
		return nil, nil, err
	}
	initCtx, initCancel := context.WithTimeout(ctx, 30*time.Second)
	defer initCancel()
	mstore, err := store.NewMySQL(initCtx, cli.SQL)
	if err != nil {
		_ = cli.Close()
		return nil, nil, err
	}
	logger.Info("mysql store ready", "dsn_prefix", cfg.DSN[:min(len(cfg.DSN), 30)])
	return mstore, cli, nil
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
```

Note: `openStore` returns `*store.MySQL` because it implements both `Store` and `CircuitStore`. When MySQL is disabled/failed, the main crashes unless there's a fallback. Add a fallback path to use `*store.Memory`:

Replace the full `openStore` function to handle the memory fallback properly:

```go
func openStore(ctx context.Context, logger *slog.Logger) (interface {
	store.Store
	store.CircuitStore
}, *pmysql.DB, error) {
	if strings.EqualFold(os.Getenv("MYSQL_DISABLED"), "true") {
		logger.Warn("MYSQL_DISABLED=true; using in-memory store")
		return store.NewMemory(), nil, nil
	}
	cfg := pmysql.ConfigFromEnv(serviceName)
	dialCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()
	cli, err := pmysql.Connect(dialCtx, cfg)
	if err != nil {
		if strings.EqualFold(os.Getenv("FALLBACK_MEMORY"), "true") {
			logger.Warn("mysql connect failed; falling back to in-memory store", "err", err)
			return store.NewMemory(), nil, nil
		}
		return nil, nil, err
	}
	initCtx, initCancel := context.WithTimeout(ctx, 30*time.Second)
	defer initCancel()
	mstore, err := store.NewMySQL(initCtx, cli.SQL)
	if err != nil {
		_ = cli.Close()
		return nil, nil, err
	}
	logger.Info("mysql store ready", "dsn_prefix", cfg.DSN[:min(len(cfg.DSN), 30)])
	return mstore, cli, nil
}
```

And update the call to `httpapi.New` accordingly:
```go
httpapi.Register(srv, httpapi.New(st, st, logger))
```

Full correct main.go:

```go
package main

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/theding0x/capital-simulator/pkg/httpx"
	applog "github.com/theding0x/capital-simulator/pkg/log"
	pmysql "github.com/theding0x/capital-simulator/pkg/mysql"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/transport/httpapi"
)

const serviceName = "agent-service"

type agentStore interface {
	store.Store
	store.CircuitStore
}

func main() {
	logger := applog.New(serviceName)
	applog.SetDefault(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st, mysqlDB, err := openStore(ctx, logger)
	if err != nil {
		logger.Error("could not open any store", "err", err)
		os.Exit(1)
	}
	if mysqlDB != nil {
		defer func() { _ = mysqlDB.Close() }()
	}

	addr := getenv("SERVICE_ADDR", ":8082")
	srv := httpx.New(httpx.Config{Addr: addr}, logger)

	httpapi.Register(srv, httpapi.New(st, st, logger))
	srv.MarkReady(true)

	if err := srv.Run(ctx); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func openStore(ctx context.Context, logger *slog.Logger) (agentStore, *pmysql.DB, error) {
	if strings.EqualFold(os.Getenv("MYSQL_DISABLED"), "true") {
		logger.Warn("MYSQL_DISABLED=true; using in-memory store")
		return store.NewMemory(), nil, nil
	}
	cfg := pmysql.ConfigFromEnv(serviceName)
	dialCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()
	cli, err := pmysql.Connect(dialCtx, cfg)
	if err != nil {
		if strings.EqualFold(os.Getenv("FALLBACK_MEMORY"), "true") {
			logger.Warn("mysql connect failed; falling back to in-memory store", "err", err)
			return store.NewMemory(), nil, nil
		}
		return nil, nil, err
	}
	initCtx, initCancel := context.WithTimeout(ctx, 30*time.Second)
	defer initCancel()
	mstore, err := store.NewMySQL(initCtx, cli.SQL)
	if err != nil {
		_ = cli.Close()
		return nil, nil, err
	}
	logger.Info("mysql store ready", "dsn_prefix", cfg.DSN[:min(len(cfg.DSN), 30)])
	return mstore, cli, nil
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
```

- [ ] **Step 2: Ask user to run tests**

```bash
make vet test build
```
Expected: builds successfully.

---

## Task 8: API Gateway — add agent proxy routes

**Files:**
- Modify: `services/api-gateway/cmd/api-gateway/main.go`

- [ ] **Step 1: Add agent proxy in main.go**

After the commodityProxy block and before `srv.MarkReady(true)`, add:

```go
	// Reverse-proxy routes to agent-service.
	agentURL := getenv("AGENT_SERVICE_URL", "http://agent-service:8082")
	agentProxy, err := proxy.New(agentURL, logger)
	if err != nil {
		logger.Error("failed to build agent proxy", "err", err)
		os.Exit(1)
	}
	srv.Handle("/v1/agents", agentProxy)
	srv.Handle("/v1/agents/{rest...}", agentProxy)
```

Also update the `handleInfo` response body: change `"status": "ch-2-exchange"` to `"status": "ch-4-capital"` and update `"chapter"` to `"Capital Vol. I, Ch. 4 - The General Formula for Capital"`.

- [ ] **Step 2: Ask user to run build**

```bash
make vet test build
```
Expected: vet and build pass.

---

## Task 9: TypeScript types and API client

**Files:**
- Modify: `web/src/types.ts`
- Modify: `web/src/api.ts`

- [ ] **Step 1: Add agent-service types to web/src/types.ts**

Append after the last `CreateWorldMoneyTransferInput` interface:

```typescript
// --- agent-service types (Ch. 4: The General Formula for Capital) -----------

export interface Agent {
  id: string;
  name: string;
  class: "capitalist" | "worker" | "miser";
  money_balance: number; // Pence (pennies); divide by 100 for £
  hoarding: boolean;
  created_at: string;
  updated_at: string;
}

export interface CapitalCircuit {
  id: string;
  agent_id: string;
  m_advanced: number; // Pence
  commodity_id: string;
  m_returned: number; // Pence
  surplus_value: number; // Pence; = m_returned - m_advanced
  circuit_type: "C-M-C" | "M-C-M-prime";
  created_at: string;
}

export interface CreateAgentInput {
  name: string;
  class: "capitalist" | "worker" | "miser";
  money_balance: number;
}

export interface UpdateAgentInput {
  name?: string;
  money_balance?: number;
}

export interface CreateCircuitInput {
  m_advanced: number;
  commodity_id: string;
  m_returned: number;
  circuit_type: "C-M-C" | "M-C-M-prime";
}

export interface ReinvestInput {
  commodity_id: string;
  m_returned: number;
}
```

- [ ] **Step 2: Add agent-service API methods to web/src/api.ts**

Add imports at the top of the import block:
```typescript
import type {
  // ... existing imports ...
  Agent,
  CapitalCircuit,
  CreateAgentInput,
  UpdateAgentInput,
  CreateCircuitInput as CreateCircuitCh04Input,
  ReinvestInput,
} from "./types";
```

Note: `CreateCircuitInput` is already used for the Ch.03 market-service circuit type. Name the Ch04 version `CreateCircuitCh04Input` in the import, but expose it as the correct api method names.

Actually, since the import alias is internal, just add these to the api object:

```typescript
  // --- agent-service (Ch. 4) ---

  listAgents: (classFilter?: string) =>
    http<{ items: Agent[] }>(
      classFilter ? `/v1/agents?class=${encodeURIComponent(classFilter)}` : "/v1/agents"
    ).then((r) => r.items),

  createAgent: (input: CreateAgentInput) =>
    http<Agent>("/v1/agents", { method: "POST", body: JSON.stringify(input) }),

  getAgent: (id: string) => http<Agent>(`/v1/agents/${id}`),

  updateAgent: (id: string, input: UpdateAgentInput) =>
    http<Agent>(`/v1/agents/${id}`, { method: "PATCH", body: JSON.stringify(input) }),

  deleteAgent: (id: string) =>
    http<void>(`/v1/agents/${id}`, { method: "DELETE" }),

  createAgentCircuit: (agentId: string, input: {
    m_advanced: number;
    commodity_id: string;
    m_returned: number;
    circuit_type: "C-M-C" | "M-C-M-prime";
  }) =>
    http<CapitalCircuit>(`/v1/agents/${agentId}/circuits`, {
      method: "POST",
      body: JSON.stringify(input),
    }),

  listAgentCircuits: (agentId: string) =>
    http<{ items: CapitalCircuit[] }>(`/v1/agents/${agentId}/circuits`).then((r) => r.items),

  reinvestAgent: (agentId: string, commodityId: string, mReturned: number) =>
    http<CapitalCircuit>(`/v1/agents/${agentId}/reinvest`, {
      method: "POST",
      body: JSON.stringify({ commodity_id: commodityId, m_returned: mReturned }),
    }),

  hoardAgent: (agentId: string) =>
    http<Agent>(`/v1/agents/${agentId}/hoard`, { method: "POST", body: "{}" }),
```

Also add the new types to the import block at the top of api.ts:
```typescript
  Agent,
  CapitalCircuit,
  CreateAgentInput,
  UpdateAgentInput,
```

- [ ] **Step 3: Ask user to run lint**

```bash
cd web && npm run lint
```
Expected: no TypeScript errors.

---

## Task 10: Ch04Capital frontend panel

**Files:**
- Create: `web/src/chapters/Ch04Capital.tsx`

- [ ] **Step 1: Write Ch04Capital.tsx**

```tsx
import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type { Agent, CapitalCircuit } from "../types";

interface Ch04Props {
  onSharedChanged: () => void;
}

function penceToGBP(pence: number): string {
  return `£${(pence / 100).toFixed(2)}`;
}

export function Ch04Capital({ onSharedChanged: _onSharedChanged }: Ch04Props) {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [error, setError] = useState<string | null>(null);

  async function refreshAgents() {
    try {
      const list = await api.listAgents();
      setAgents(list);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    }
  }

  useEffect(() => {
    refreshAgents();
  }, []);

  const capitalists = agents.filter((a) => a.class === "capitalist");
  const workers = agents.filter((a) => a.class === "worker");
  const misers = agents.filter((a) => a.class === "miser");

  return (
    <>
      <CreateAgentPanel onCreated={refreshAgents} />
      {error && <p className="error">{error}</p>}
      <AgentClassSection title="Capitalists" agents={capitalists} onChanged={refreshAgents} />
      <AgentClassSection title="Workers" agents={workers} onChanged={refreshAgents} />
      <AgentClassSection title="Misers" agents={misers} onChanged={refreshAgents} />
    </>
  );
}

function CreateAgentPanel({ onCreated }: { onCreated: () => void }) {
  const [name, setName] = useState("");
  const [cls, setCls] = useState<"capitalist" | "worker" | "miser">("capitalist");
  const [balance, setBalance] = useState(10000);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await api.createAgent({ name, class: cls, money_balance: balance });
      setName("");
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Create Agent</h2>
      <form className="form-grid" onSubmit={submit}>
        <label>
          Name
          <input value={name} onChange={(e) => setName(e.target.value)} required />
        </label>
        <label>
          Class
          <select value={cls} onChange={(e) => setCls(e.target.value as typeof cls)}>
            <option value="capitalist">Capitalist</option>
            <option value="worker">Worker</option>
            <option value="miser">Miser</option>
          </select>
        </label>
        <label>
          Initial balance (pence)
          <input
            type="number"
            value={balance}
            onChange={(e) => setBalance(Number(e.target.value))}
            min={0}
          />
        </label>
        <button type="submit">Create</button>
        {err && <span className="error">{err}</span>}
      </form>
    </section>
  );
}

function AgentClassSection({
  title,
  agents,
  onChanged,
}: {
  title: string;
  agents: Agent[];
  onChanged: () => void;
}) {
  if (agents.length === 0) return null;
  return (
    <section className="card">
      <h2>{title}</h2>
      <div className="item-list">
        {agents.map((a) => (
          <AgentCard key={a.id} agent={a} onChanged={onChanged} />
        ))}
      </div>
    </section>
  );
}

function AgentCard({ agent: a, onChanged }: { agent: Agent; onChanged: () => void }) {
  const [circuits, setCircuits] = useState<CapitalCircuit[]>([]);
  const [showCircuits, setShowCircuits] = useState(false);
  const [err, setErr] = useState<string | null>(null);

  async function loadCircuits() {
    try {
      const list = await api.listAgentCircuits(a.id);
      setCircuits(list);
      setShowCircuits(true);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  async function hoard() {
    setErr(null);
    try {
      await api.hoardAgent(a.id);
      onChanged();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <div className="item-card">
      <div className="item-header">
        <span className="item-name">{a.name}</span>
        <span className="item-meta">{penceToGBP(a.money_balance)}</span>
        {a.hoarding && <span className="item-tag">hoarding</span>}
      </div>
      <div className="item-actions">
        <button onClick={loadCircuits} type="button">
          {showCircuits ? "Refresh circuits" : "Show circuits"}
        </button>
        {a.class === "miser" && !a.hoarding && (
          <button onClick={hoard} type="button">
            Hoard
          </button>
        )}
      </div>
      {err && <span className="error">{err}</span>}
      {showCircuits && (
        <>
          <CircuitTable circuits={circuits} />
          {a.class !== "miser" && (
            <CreateCircuitForm agentId={a.id} agentClass={a.class} onCreated={() => { loadCircuits(); onChanged(); }} />
          )}
        </>
      )}
    </div>
  );
}

function CircuitTable({ circuits }: { circuits: CapitalCircuit[] }) {
  if (circuits.length === 0) return <p className="muted small">No circuits yet.</p>;
  return (
    <table className="data-table">
      <thead>
        <tr>
          <th>Type</th>
          <th>M (advanced)</th>
          <th>C (commodity)</th>
          <th>M′ (returned)</th>
          <th>∆M (surplus)</th>
        </tr>
      </thead>
      <tbody>
        {circuits.map((c) => (
          <tr key={c.id}>
            <td>{c.circuit_type}</td>
            <td>{penceToGBP(c.m_advanced)}</td>
            <td className="monospace small">{c.commodity_id.slice(0, 8)}&hellip;</td>
            <td>{penceToGBP(c.m_returned)}</td>
            <td className={c.surplus_value > 0 ? "positive" : c.surplus_value < 0 ? "negative" : ""}>
              {c.surplus_value > 0 ? "+" : ""}{penceToGBP(c.surplus_value)}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

function CreateCircuitForm({
  agentId,
  agentClass,
  onCreated,
}: {
  agentId: string;
  agentClass: "capitalist" | "worker" | "miser";
  onCreated: () => void;
}) {
  const [commodityId, setCommodityId] = useState("");
  const [mAdvanced, setMAdvanced] = useState(10000);
  const [mReturned, setMReturned] = useState(11000);
  const [circuitType, setCircuitType] = useState<"C-M-C" | "M-C-M-prime">(
    agentClass === "worker" ? "C-M-C" : "M-C-M-prime"
  );
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      await api.createAgentCircuit(agentId, {
        m_advanced: mAdvanced,
        commodity_id: commodityId,
        m_returned: mReturned,
        circuit_type: circuitType,
      });
      onCreated();
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <form className="form-grid" onSubmit={submit}>
      <label>
        Commodity ID
        <input value={commodityId} onChange={(e) => setCommodityId(e.target.value)} required />
      </label>
      <label>
        M advanced (pence)
        <input
          type="number"
          value={mAdvanced}
          onChange={(e) => setMAdvanced(Number(e.target.value))}
          min={1}
        />
      </label>
      <label>
        M′ returned (pence)
        <input
          type="number"
          value={mReturned}
          onChange={(e) => setMReturned(Number(e.target.value))}
          min={0}
        />
      </label>
      {agentClass !== "worker" && (
        <label>
          Circuit type
          <select
            value={circuitType}
            onChange={(e) => setCircuitType(e.target.value as typeof circuitType)}
          >
            <option value="M-C-M-prime">M—C—M′ (capital)</option>
            <option value="C-M-C">C—M—C (worker)</option>
          </select>
        </label>
      )}
      <button type="submit">Record circuit</button>
      {err && <span className="error">{err}</span>}
    </form>
  );
}
```

---

## Task 11: Wire Ch04 into the UI shell

**Files:**
- Modify: `web/src/components/ChapterShell.tsx`
- Modify: `web/src/chapters/registry.ts`

- [ ] **Step 1: Add ch04 quote and wire component in ChapterShell.tsx**

In the `QUOTES` object, add:
```typescript
  ch04: "The circulation of commodities is the starting-point of capital.",
```

In the imports, add:
```typescript
import { Ch04Capital } from "../chapters/Ch04Capital";
```

In the JSX chain that renders `ch01`/`ch02`/`ch03`, add:
```tsx
        ) : activeChapterId === "ch04" ? (
          <Ch04Capital onSharedChanged={onSharedChanged} />
```

- [ ] **Step 2: Mark ch04 as "done" in registry.ts**

In `web/src/chapters/registry.ts`, change:
```typescript
  { id: "ch04", number: 4, title: "The General Formula for Capital", ..., status: "pending" },
```
to:
```typescript
  { id: "ch04", number: 4, title: "The General Formula for Capital", ..., status: "done" },
```

- [ ] **Step 3: Ask user to run lint and build**

```bash
cd web && npm run lint && npm run build
```
Expected: no TypeScript errors, build succeeds.

---

## Task 12: Update architecture docs

**Files:**
- Modify: `docs/architecture.md`

- [ ] **Step 1: Update chapter status table**

Find the chapter status table in `docs/architecture.md`. Add or update the Ch. 4 entry to show it as complete. The exact location depends on the table format; typically add a row like:

```
| Ch. 4  | The General Formula for Capital | agent-service | Agents, CapitalCircuit, M—C—M′, class positions |
```

---

## Self-Review Checklist

After all tasks are implemented, verify spec coverage:

**Concepts → types:**
- [ ] `Agent`, `ID`, `Class`, `Pence`, `CircuitType`, `CapitalCircuit`, `Update` defined ✓
- [ ] Constants `Capitalist`, `Worker`, `Miser`, `CircuitCMC`, `CircuitMCM` defined ✓
- [ ] `Store` interface: Create/Get/List/Update/Delete + ListByClass ✓
- [ ] `CircuitStore` interface: CreateCircuit/GetCircuit/ListCircuits ✓
- [ ] `Memory` and `MySQL` store implementations ✓

**Fixtures:**
- [ ] §1: Advance(10000) on 10000-balance agent → balance=0 ✓ (agent_test.go)
- [ ] §8: SurplusValue = 11000-10000 = 1000 ✓ (agent_test.go)
- [ ] §10: Zero-surplus circuit is valid ✓ (agent_test.go)
- [ ] §14: Miser.Hoard succeeds; Miser.Reinvest → ErrNotCapitalist; Capitalist.Reinvest succeeds; Capitalist.Hoard → ErrNotCapitalist ✓ (agent_test.go)
- [ ] §15: After Realise, second circuit uses full new balance ✓ (agent_test.go)
- [ ] §6: Worker can only CircuitCMC; CircuitMCM → ErrWrongClass ✓ (handler.go CreateCircuit)

**Invariants:**
- [ ] SurplusValue == MReturned - MAdvanced enforced in Validate ✓
- [ ] MoneyBalance >= 0: Advance returns ErrInsufficientFunds ✓
- [ ] MAdvanced > 0: circuit.Validate() ✓
- [ ] After Realise, balance = prevBalance + MReturned ✓
- [ ] Miser's balance unchanged: Advance/Reinvest → ErrNotCapitalist ✓

**HTTP endpoints:**
- [ ] POST /v1/agents ✓
- [ ] GET /v1/agents (+ ?class=) ✓
- [ ] GET /v1/agents/{id} ✓
- [ ] PATCH /v1/agents/{id} ✓
- [ ] DELETE /v1/agents/{id} ✓
- [ ] POST /v1/agents/{id}/circuits ✓
- [ ] GET /v1/agents/{id}/circuits ✓
- [ ] POST /v1/agents/{id}/reinvest ✓
- [ ] POST /v1/agents/{id}/hoard ✓

**Frontend:**
- [ ] Agent and CapitalCircuit in types.ts ✓
- [ ] API methods in api.ts ✓
- [ ] Ch04 panel: agents by class, £ balance, M→C→M′ table with ∆M ✓
- [ ] ch04 status = "done" in registry ✓

**Gateway + Docs:**
- [ ] Agent proxy routes in api-gateway ✓
- [ ] architecture.md updated ✓
