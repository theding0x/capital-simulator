# Ch. 10 — The Working-Day Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement Ch. 10 "The Working-Day" domain types, persistence, HTTP API, and React UI in agent-service.

**Architecture:** New domain types (`WorkingDay`, `RelaySchedule`, etc.) land in `agent/working_day.go`; the store interface gains a `WorkingDayStore`; Memory and MySQL implementations get the new methods; a new `working_day_handler.go` handles four HTTP endpoints; the React panel lives in `Ch10WorkingDay.tsx` and is wired into `ChapterShell`.

**Tech Stack:** Go 1.25 (agent-service), goose SQL migrations, React 18 + TypeScript (web/src)

---

## File Map

| File | Action | Responsibility |
|------|--------|---------------|
| `services/agent-service/internal/agent/working_day.go` | Create | All Ch.10 domain types, methods, sentinel errors |
| `services/agent-service/internal/agent/working_day_test.go` | Create | Unit tests for domain types (TDD) |
| `services/agent-service/internal/store/store.go` | Modify | Add `WorkingDayStore` interface |
| `services/agent-service/internal/store/memory.go` | Modify | Implement `WorkingDayStore` in `Memory` |
| `services/agent-service/internal/store/working_day_memory_test.go` | Create | Store-layer tests for Memory |
| `services/agent-service/internal/store/migrations/00008_ch10_working_day.sql` | Create | `working_days` and `relay_schedules` DDL |
| `services/agent-service/internal/store/mysql.go` | Modify | Implement `WorkingDayStore` in `MySQL` |
| `services/agent-service/internal/transport/httpapi/handler.go` | Modify | Add `WorkingDayStore` field to `Handler` |
| `services/agent-service/internal/transport/httpapi/working_day_handler.go` | Create | HTTP handler for working-day + relay-schedule endpoints |
| `services/agent-service/internal/transport/httpapi/routes.go` | Modify | Register new routes |
| `services/agent-service/cmd/agent-service/main.go` | Modify | Embed `WorkingDayStore` in `agentStore` interface |
| `services/api-gateway/cmd/api-gateway/main.go` | Modify | Add proxy routes for `/v1/working-days` and `/v1/relay-schedules` |
| `web/src/types.ts` | Modify | Add Ch.10 TypeScript types |
| `web/src/api.ts` | Modify | Add Ch.10 API methods |
| `web/src/chapters/Ch10WorkingDay.tsx` | Create | React panel: working-day builder + relay-schedule visualiser |
| `web/src/components/ChapterShell.tsx` | Modify | Add `ch10` branch |
| `web/src/chapters/registry.ts` | Modify | Mark `ch10` status `"done"` |
| `docs/architecture.md` | Modify | Update chapter table for Ch.10 |

---

## Task 1: Domain Types — Core Working-Day types

**Files:**
- Create: `services/agent-service/internal/agent/working_day.go`

- [ ] **Step 1: Write the failing test**

Create `services/agent-service/internal/agent/working_day_test.go`:

```go
package agent_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

func TestWorkingDay_TotalMinutes(t *testing.T) {
	t.Parallel()
	wd := agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 60}
	if got := wd.TotalMinutes(); got != 420 {
		t.Fatalf("TotalMinutes = %d, want 420", got)
	}
}

func TestWorkingDay_RateOfSurplusValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		nl, sl int64
		want   float64
	}{
		{360, 60, 1.0 / 6.0},   // §1 fixture I: 6h necessary, 1h surplus → 16.67%
		{360, 180, 0.50},        // §1 fixture II: 6h necessary, 3h surplus → 50%
		{360, 360, 1.00},        // §1 fixture III: 6h necessary, 6h surplus → 100%
		{480, 480, 1.00},        // §1: rate=100% does not fix day length
		{600, 600, 1.00},
		{720, 720, 1.00},
	}
	for _, tt := range tests {
		wd := agent.WorkingDay{
			NecessaryLabourMinutes: agent.NecessaryLabourMinutes(tt.nl),
			SurplusLabourMinutes:   agent.SurplusLabourMinutes(tt.sl),
		}
		got := wd.RateOfSurplusValue()
		if diff := got - tt.want; diff > 1e-9 || diff < -1e-9 {
			t.Errorf("nl=%d sl=%d: RateOfSurplusValue()=%f, want %f", tt.nl, tt.sl, got, tt.want)
		}
	}
}

func TestWorkingDay_Validate_PhysicalMax(t *testing.T) {
	t.Parallel()
	// 24h exactly is valid
	wd := agent.WorkingDay{NecessaryLabourMinutes: 720, SurplusLabourMinutes: 720}
	if err := wd.Validate(); err != nil {
		t.Fatalf("expected no error for 1440 min, got %v", err)
	}
	// Over 24h is invalid
	over := agent.WorkingDay{NecessaryLabourMinutes: 721, SurplusLabourMinutes: 720}
	if err := over.Validate(); err != agent.ErrWorkingDayExceedsPhysicalMax {
		t.Fatalf("expected ErrWorkingDayExceedsPhysicalMax, got %v", err)
	}
}

func TestWorkingDay_Validate_ZeroNecessary(t *testing.T) {
	t.Parallel()
	wd := agent.WorkingDay{NecessaryLabourMinutes: 0, SurplusLabourMinutes: 60}
	if err := wd.Validate(); err == nil {
		t.Fatal("expected error for zero necessary labour, got nil")
	}
}

func TestValidateConstraint(t *testing.T) {
	t.Parallel()
	c := agent.WorkingDayConstraint{Label: "Factory Act 1850", StatutoryLimitMinutes: 630}
	ok := agent.WorkingDay{NecessaryLabourMinutes: 315, SurplusLabourMinutes: 315} // 630 min exactly
	if err := agent.ValidateConstraint(ok, c); err != nil {
		t.Fatalf("expected nil for total==limit, got %v", err)
	}
	over := agent.WorkingDay{NecessaryLabourMinutes: 316, SurplusLabourMinutes: 315} // 631 min
	if err := agent.ValidateConstraint(over, c); err != agent.ErrWorkingDayExceedsStatutoryLimit {
		t.Fatalf("expected ErrWorkingDayExceedsStatutoryLimit, got %v", err)
	}
}

func TestOverwork_AnnualMinutes(t *testing.T) {
	t.Parallel()
	// §3 fixture: Factory Inspector: 5 min/day × 300 working days = 1500 min per year
	o := agent.Overwork{MinutesPerDay: 5}
	if got := o.AnnualMinutes(300); got != 1500 {
		t.Fatalf("AnnualMinutes(300) = %d, want 1500", got)
	}
}

func TestWeeklyFixtures(t *testing.T) {
	t.Parallel()
	// §2: 6h necessary + 6h surplus = 36h surplus per week (6 days × 6h)
	wd := agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360}
	if got := wd.TotalMinutes() * 6 / 60; got != 72 {
		t.Errorf("weekly total hours = %d, want 72", got)
	}
	if got := int64(wd.SurplusLabourMinutes) * 6 / 60; got != 36 {
		t.Errorf("weekly surplus hours = %d, want 36", got)
	}
}

func TestWallachianCorvee(t *testing.T) {
	t.Parallel()
	// §2: Wallachian corvée: 56 corvée days / 140 working days → rate 56/84 ≈ 66.67%
	wd := agent.WorkingDay{
		NecessaryLabourMinutes: agent.NecessaryLabourMinutes(84 * 480),
		SurplusLabourMinutes:   agent.SurplusLabourMinutes(56 * 480),
	}
	got := wd.RateOfSurplusValue()
	want := float64(56) / float64(84)
	if diff := got - want; diff > 1e-6 || diff < -1e-6 {
		t.Errorf("Wallachian rate = %f, want %f", got, want)
	}
}

func TestRelaySchedule_Validate(t *testing.T) {
	t.Parallel()
	rs := agent.RelaySchedule{
		Sets: [2]agent.RelaySet{
			{ShiftKind: agent.ShiftDay, WorkingDay: agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360}},
			{ShiftKind: agent.ShiftNight, WorkingDay: agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 0}},
		},
	}
	// 720 + 360 = 1080 <= 1440
	if err := rs.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRelaySchedule_Validate_ExceedsPhysical(t *testing.T) {
	t.Parallel()
	rs := agent.RelaySchedule{
		Sets: [2]agent.RelaySet{
			{ShiftKind: agent.ShiftDay, WorkingDay: agent.WorkingDay{NecessaryLabourMinutes: 720, SurplusLabourMinutes: 360}},
			{ShiftKind: agent.ShiftNight, WorkingDay: agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 1}},
		},
	}
	if err := rs.Validate(); err == nil {
		t.Fatal("expected error for combined > 1440 min")
	}
}

func TestFactoryActs(t *testing.T) {
	t.Parallel()
	fa1833 := agent.FactoryAct{Year: 1833, ChildLimitMinutes: 480, YoungPersonLimitMinutes: 720}
	if fa1833.Year != 1833 || fa1833.ChildLimitMinutes != 480 {
		t.Errorf("FactoryAct 1833 wrong: %+v", fa1833)
	}
	fa1847 := agent.FactoryAct{Year: 1847, YoungPersonLimitMinutes: 600}
	if fa1847.YoungPersonLimitMinutes != 600 {
		t.Errorf("FactoryAct 1847 wrong: %+v", fa1847)
	}
}
```

- [ ] **Step 2: Run test to confirm it fails**

```bash
cd /mnt/c/Users/AaronHulse/IdeaProjects/capital-simulator
go test ./services/agent-service/internal/agent/... 2>&1 | head -30
```
Expected: compile error — `working_day.go` not yet created.

- [ ] **Step 3: Write the domain types**

Create `services/agent-service/internal/agent/working_day.go`:

```go
package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

// PhysicalMaxMinutes is the hard ceiling for any working day: 24 h × 60 min.
const PhysicalMaxMinutes int64 = 1440

var (
	ErrWorkingDayExceedsPhysicalMax   = errors.New("agent: working day exceeds physical maximum of 24 hours")
	ErrWorkingDayExceedsStatutoryLimit = errors.New("agent: working day exceeds statutory limit")
)

// Named value types for the two segments of the working day.
type NecessaryLabourMinutes int64
type SurplusLabourMinutes int64
type StatutoryLimitMinutes int64

// WorkingDayID is a 96-bit hex identifier for stored WorkingDay records.
type WorkingDayID string

func (id WorkingDayID) IsZero() bool { return id == "" }

func NewWorkingDayID() WorkingDayID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return WorkingDayID(hex.EncodeToString(b))
}

// WorkingDay encodes the two segments of the working day: AB (necessary) and BC (surplus).
type WorkingDay struct {
	ID                     WorkingDayID           `json:"id"`
	NecessaryLabourMinutes NecessaryLabourMinutes `json:"necessary_labour_minutes"`
	SurplusLabourMinutes   SurplusLabourMinutes   `json:"surplus_labour_minutes"`
	CreatedAt              time.Time              `json:"created_at"`
}

// TotalMinutes returns necessary + surplus labour time in minutes.
func (wd WorkingDay) TotalMinutes() int64 {
	return int64(wd.NecessaryLabourMinutes) + int64(wd.SurplusLabourMinutes)
}

// RateOfSurplusValue returns s/v — the degree of exploitation.
func (wd WorkingDay) RateOfSurplusValue() float64 {
	return float64(wd.SurplusLabourMinutes) / float64(wd.NecessaryLabourMinutes)
}

// Validate checks that necessary labour is positive and the total does not exceed
// the physical maximum of 24 hours.
func (wd WorkingDay) Validate() error {
	if wd.NecessaryLabourMinutes <= 0 {
		return errors.New("agent: necessary_labour_minutes must be positive")
	}
	if wd.TotalMinutes() > PhysicalMaxMinutes {
		return ErrWorkingDayExceedsPhysicalMax
	}
	return nil
}

// WorkingDayConstraint pairs a statutory limit with the jurisdiction or epoch it applies to.
type WorkingDayConstraint struct {
	Label                 string                `json:"label"`
	StatutoryLimitMinutes StatutoryLimitMinutes `json:"statutory_limit_minutes"`
}

// ValidateConstraint returns ErrWorkingDayExceedsStatutoryLimit when the working day
// total exceeds the constraint's statutory limit.
func ValidateConstraint(wd WorkingDay, c WorkingDayConstraint) error {
	if wd.TotalMinutes() > int64(c.StatutoryLimitMinutes) {
		return ErrWorkingDayExceedsStatutoryLimit
	}
	return nil
}

// NormalWorkingDay is a WorkingDay validated against a statutory limit.
type NormalWorkingDay struct {
	WorkingDay
	Constraint WorkingDayConstraint `json:"constraint"`
}

// Validate checks the working day invariants and the statutory limit.
func (n NormalWorkingDay) Validate() error {
	if err := n.WorkingDay.Validate(); err != nil {
		return err
	}
	return ValidateConstraint(n.WorkingDay, n.Constraint)
}

// Overwork represents minutes stolen above the statutory limit per working day.
type Overwork struct {
	MinutesPerDay int64 `json:"minutes_per_day"`
}

// AnnualMinutes returns total overwork minutes across workingDays working days per year.
func (o Overwork) AnnualMinutes(workingDays int) int64 {
	return o.MinutesPerDay * int64(workingDays)
}

// FactoryAct records a historical limit on working time.
type FactoryAct struct {
	Year                    int   `json:"year"`
	ChildLimitMinutes       int64 `json:"child_limit_minutes"`
	YoungPersonLimitMinutes int64 `json:"young_person_limit_minutes"`
	AdultLimitMinutes       int64 `json:"adult_limit_minutes"`
}

// ShiftKind distinguishes day and night shifts in a relay system.
type ShiftKind string

const (
	ShiftDay   ShiftKind = "day"
	ShiftNight ShiftKind = "night"
)

// TimerClass maps to adult vs. child workers in the relay system.
type TimerClass string

const (
	TimerFull TimerClass = "full"
	TimerHalf TimerClass = "half"
)

// RelaySet is one half of a relay system: a shift kind, a working day, and the set of workers.
type RelaySet struct {
	ShiftKind  ShiftKind  `json:"shift_kind"`
	WorkingDay WorkingDay `json:"working_day"`
	WorkerIDs  []AgentID  `json:"worker_ids"`
}

// RelayScheduleID is a 96-bit hex identifier for stored RelaySchedule records.
type RelayScheduleID string

func (id RelayScheduleID) IsZero() bool { return id == "" }

func NewRelayScheduleID() RelayScheduleID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return RelayScheduleID(hex.EncodeToString(b))
}

// RelaySchedule is two alternating worker sets covering complementary shifts.
type RelaySchedule struct {
	ID        RelayScheduleID `json:"id"`
	Sets      [2]RelaySet     `json:"sets"`
	CreatedAt time.Time       `json:"created_at"`
}

// Validate ensures each set's working day is individually valid and their combined
// total does not exceed the physical maximum of 24 hours.
func (rs RelaySchedule) Validate() error {
	if err := rs.Sets[0].WorkingDay.Validate(); err != nil {
		return err
	}
	if err := rs.Sets[1].WorkingDay.Validate(); err != nil {
		return err
	}
	combined := rs.Sets[0].WorkingDay.TotalMinutes() + rs.Sets[1].WorkingDay.TotalMinutes()
	if combined > PhysicalMaxMinutes {
		return ErrWorkingDayExceedsPhysicalMax
	}
	return nil
}
```

- [ ] **Step 4: Run tests to confirm pass**

```bash
go test ./services/agent-service/internal/agent/... -run TestWorkingDay -v 2>&1
go test ./services/agent-service/internal/agent/... -run TestValidateConstraint -v 2>&1
go test ./services/agent-service/internal/agent/... -run TestOverwork -v 2>&1
go test ./services/agent-service/internal/agent/... -run TestRelay -v 2>&1
go test ./services/agent-service/internal/agent/... -run TestWeekly -v 2>&1
go test ./services/agent-service/internal/agent/... -run TestWallachian -v 2>&1
go test ./services/agent-service/internal/agent/... -run TestFactory -v 2>&1
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add services/agent-service/internal/agent/working_day.go \
        services/agent-service/internal/agent/working_day_test.go
git commit -m "feat(ch10): domain types for WorkingDay, RelaySchedule, FactoryAct, Overwork"
```

---

## Task 2: Store Interface

**Files:**
- Modify: `services/agent-service/internal/store/store.go`

- [ ] **Step 1: Add the WorkingDayStore interface**

In `services/agent-service/internal/store/store.go`, append after the `LabourProcessStore` interface:

```go
// WorkingDayStore is the persistence contract for Ch. 10 working-day records.
type WorkingDayStore interface {
	CreateWorkingDay(ctx context.Context, wd agent.WorkingDay) (agent.WorkingDay, error)
	GetWorkingDay(ctx context.Context, id agent.WorkingDayID) (agent.WorkingDay, error)
	CreateRelaySchedule(ctx context.Context, rs agent.RelaySchedule) (agent.RelaySchedule, error)
	GetRelaySchedule(ctx context.Context, id agent.RelayScheduleID) (agent.RelaySchedule, error)
}
```

- [ ] **Step 2: Verify the file compiles**

```bash
go build ./services/agent-service/... 2>&1
```
Expected: compile error in `main.go` — `agentStore` doesn't yet embed `WorkingDayStore`. Ignore that for now.

---

## Task 3: Memory Store Implementation + Tests

**Files:**
- Modify: `services/agent-service/internal/store/memory.go`
- Create: `services/agent-service/internal/store/working_day_memory_test.go`

- [ ] **Step 1: Write the failing tests**

Create `services/agent-service/internal/store/working_day_memory_test.go`:

```go
package store_test

import (
	"context"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

func TestMemory_WorkingDay_RoundTrip(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()

	wd := agent.WorkingDay{
		NecessaryLabourMinutes: 360,
		SurplusLabourMinutes:   360,
	}
	created, err := m.CreateWorkingDay(ctx, wd)
	if err != nil {
		t.Fatalf("CreateWorkingDay: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected non-zero ID")
	}
	if created.CreatedAt.IsZero() {
		t.Fatal("expected non-zero CreatedAt")
	}
	if created.NecessaryLabourMinutes != 360 || created.SurplusLabourMinutes != 360 {
		t.Errorf("fields not preserved: %+v", created)
	}

	got, err := m.GetWorkingDay(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetWorkingDay: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID mismatch: got %q want %q", got.ID, created.ID)
	}
}

func TestMemory_WorkingDay_NotFound(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.GetWorkingDay(context.Background(), "nonexistent")
	if err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMemory_WorkingDay_InvalidRejected(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	// NecessaryLabourMinutes = 0 is invalid
	_, err := m.CreateWorkingDay(context.Background(), agent.WorkingDay{SurplusLabourMinutes: 60})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestMemory_RelaySchedule_RoundTrip(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()

	rs := agent.RelaySchedule{
		Sets: [2]agent.RelaySet{
			{
				ShiftKind:  agent.ShiftDay,
				WorkingDay: agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360},
				WorkerIDs:  []agent.AgentID{"worker-a"},
			},
			{
				ShiftKind:  agent.ShiftNight,
				WorkingDay: agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 0},
				WorkerIDs:  []agent.AgentID{"worker-b"},
			},
		},
	}
	created, err := m.CreateRelaySchedule(ctx, rs)
	if err != nil {
		t.Fatalf("CreateRelaySchedule: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected non-zero ID")
	}

	got, err := m.GetRelaySchedule(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetRelaySchedule: %v", err)
	}
	if got.Sets[0].ShiftKind != agent.ShiftDay {
		t.Errorf("Sets[0].ShiftKind = %q, want %q", got.Sets[0].ShiftKind, agent.ShiftDay)
	}
	if len(got.Sets[0].WorkerIDs) != 1 || got.Sets[0].WorkerIDs[0] != "worker-a" {
		t.Errorf("Sets[0].WorkerIDs = %v, want [worker-a]", got.Sets[0].WorkerIDs)
	}
}

func TestMemory_RelaySchedule_NotFound(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.GetRelaySchedule(context.Background(), "nonexistent")
	if err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
```

- [ ] **Step 2: Run tests to confirm they fail**

```bash
go test ./services/agent-service/internal/store/... -run TestMemory_WorkingDay -v 2>&1 | head -20
```
Expected: compile error — `CreateWorkingDay` not on `Memory`.

- [ ] **Step 3: Add maps and methods to Memory**

In `services/agent-service/internal/store/memory.go`:

**Add fields to the `Memory` struct** (after `labourProcesses` field):
```go
workingDays     map[agent.WorkingDayID]agent.WorkingDay
relaySchedules  map[agent.RelayScheduleID]agent.RelaySchedule
```

**Initialize in `NewMemory()`** (add after `labourProcesses: make(...)`):
```go
workingDays:    make(map[agent.WorkingDayID]agent.WorkingDay),
relaySchedules: make(map[agent.RelayScheduleID]agent.RelaySchedule),
```

**Append new methods at the bottom of the file:**

```go
func (m *Memory) CreateWorkingDay(_ context.Context, wd agent.WorkingDay) (agent.WorkingDay, error) {
	if err := wd.Validate(); err != nil {
		return agent.WorkingDay{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if wd.ID.IsZero() {
		wd.ID = agent.NewWorkingDayID()
	}
	wd.CreatedAt = m.now()
	m.workingDays[wd.ID] = wd
	return wd, nil
}

func (m *Memory) GetWorkingDay(_ context.Context, id agent.WorkingDayID) (agent.WorkingDay, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	wd, ok := m.workingDays[id]
	if !ok {
		return agent.WorkingDay{}, ErrNotFound
	}
	return wd, nil
}

func (m *Memory) CreateRelaySchedule(_ context.Context, rs agent.RelaySchedule) (agent.RelaySchedule, error) {
	if err := rs.Validate(); err != nil {
		return agent.RelaySchedule{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if rs.ID.IsZero() {
		rs.ID = agent.NewRelayScheduleID()
	}
	rs.CreatedAt = m.now()
	m.relaySchedules[rs.ID] = rs
	return rs, nil
}

func (m *Memory) GetRelaySchedule(_ context.Context, id agent.RelayScheduleID) (agent.RelaySchedule, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	rs, ok := m.relaySchedules[id]
	if !ok {
		return agent.RelaySchedule{}, ErrNotFound
	}
	return rs, nil
}
```

- [ ] **Step 4: Run tests to confirm they pass**

```bash
go test ./services/agent-service/internal/store/... -run TestMemory_WorkingDay -v 2>&1
go test ./services/agent-service/internal/store/... -run TestMemory_RelaySchedule -v 2>&1
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add services/agent-service/internal/store/store.go \
        services/agent-service/internal/store/memory.go \
        services/agent-service/internal/store/working_day_memory_test.go
git commit -m "feat(ch10): WorkingDayStore interface + Memory implementation"
```

---

## Task 4: MySQL Migration

**Files:**
- Create: `services/agent-service/internal/store/migrations/00008_ch10_working_day.sql`

- [ ] **Step 1: Write the migration**

```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS working_days (
    id                       VARCHAR(24)  NOT NULL PRIMARY KEY,
    necessary_labour_minutes BIGINT       NOT NULL,
    surplus_labour_minutes   BIGINT       NOT NULL,
    created_at               DATETIME(6)  NOT NULL
);

CREATE TABLE IF NOT EXISTS relay_schedules (
    id              VARCHAR(24)  NOT NULL PRIMARY KEY,
    shift_kind_0    VARCHAR(10)  NOT NULL,
    nl_minutes_0    BIGINT       NOT NULL,
    sl_minutes_0    BIGINT       NOT NULL,
    worker_ids_0    JSON         NOT NULL,
    shift_kind_1    VARCHAR(10)  NOT NULL,
    nl_minutes_1    BIGINT       NOT NULL,
    sl_minutes_1    BIGINT       NOT NULL,
    worker_ids_1    JSON         NOT NULL,
    created_at      DATETIME(6)  NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS relay_schedules;
DROP TABLE IF EXISTS working_days;
```

- [ ] **Step 2: Commit**

```bash
git add services/agent-service/internal/store/migrations/00008_ch10_working_day.sql
git commit -m "feat(ch10): add working_days and relay_schedules MySQL migration"
```

---

## Task 5: MySQL Store Implementation

**Files:**
- Modify: `services/agent-service/internal/store/mysql.go`

- [ ] **Step 1: Append MySQL methods at the bottom of `mysql.go`**

```go
func (m *MySQL) CreateWorkingDay(ctx context.Context, wd agent.WorkingDay) (agent.WorkingDay, error) {
	if err := wd.Validate(); err != nil {
		return agent.WorkingDay{}, err
	}
	if wd.ID.IsZero() {
		wd.ID = agent.NewWorkingDayID()
	}
	wd.CreatedAt = m.now().UTC()
	const q = `INSERT INTO working_days (id, necessary_labour_minutes, surplus_labour_minutes, created_at)
		VALUES (?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(wd.ID),
		int64(wd.NecessaryLabourMinutes),
		int64(wd.SurplusLabourMinutes),
		wd.CreatedAt,
	)
	if err != nil {
		return agent.WorkingDay{}, err
	}
	return wd, nil
}

func (m *MySQL) GetWorkingDay(ctx context.Context, id agent.WorkingDayID) (agent.WorkingDay, error) {
	const q = `SELECT id, necessary_labour_minutes, surplus_labour_minutes, created_at
		FROM working_days WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	var wd agent.WorkingDay
	var wid string
	var nl, sl int64
	err := row.Scan(&wid, &nl, &sl, &wd.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.WorkingDay{}, ErrNotFound
	}
	if err != nil {
		return agent.WorkingDay{}, err
	}
	wd.ID = agent.WorkingDayID(wid)
	wd.NecessaryLabourMinutes = agent.NecessaryLabourMinutes(nl)
	wd.SurplusLabourMinutes = agent.SurplusLabourMinutes(sl)
	return wd, nil
}

func (m *MySQL) CreateRelaySchedule(ctx context.Context, rs agent.RelaySchedule) (agent.RelaySchedule, error) {
	if err := rs.Validate(); err != nil {
		return agent.RelaySchedule{}, err
	}
	if rs.ID.IsZero() {
		rs.ID = agent.NewRelayScheduleID()
	}
	rs.CreatedAt = m.now().UTC()
	wids0, err := json.Marshal(rs.Sets[0].WorkerIDs)
	if err != nil {
		return agent.RelaySchedule{}, err
	}
	wids1, err := json.Marshal(rs.Sets[1].WorkerIDs)
	if err != nil {
		return agent.RelaySchedule{}, err
	}
	const q = `INSERT INTO relay_schedules
		(id, shift_kind_0, nl_minutes_0, sl_minutes_0, worker_ids_0,
		     shift_kind_1, nl_minutes_1, sl_minutes_1, worker_ids_1, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = m.db.ExecContext(ctx, q,
		string(rs.ID),
		string(rs.Sets[0].ShiftKind),
		int64(rs.Sets[0].WorkingDay.NecessaryLabourMinutes),
		int64(rs.Sets[0].WorkingDay.SurplusLabourMinutes),
		string(wids0),
		string(rs.Sets[1].ShiftKind),
		int64(rs.Sets[1].WorkingDay.NecessaryLabourMinutes),
		int64(rs.Sets[1].WorkingDay.SurplusLabourMinutes),
		string(wids1),
		rs.CreatedAt,
	)
	if err != nil {
		return agent.RelaySchedule{}, err
	}
	return rs, nil
}

func (m *MySQL) GetRelaySchedule(ctx context.Context, id agent.RelayScheduleID) (agent.RelaySchedule, error) {
	const q = `SELECT id, shift_kind_0, nl_minutes_0, sl_minutes_0, worker_ids_0,
		shift_kind_1, nl_minutes_1, sl_minutes_1, worker_ids_1, created_at
		FROM relay_schedules WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	var rs agent.RelaySchedule
	var rid, sk0, sk1 string
	var nl0, sl0, nl1, sl1 int64
	var wids0raw, wids1raw string
	err := row.Scan(&rid, &sk0, &nl0, &sl0, &wids0raw, &sk1, &nl1, &sl1, &wids1raw, &rs.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.RelaySchedule{}, ErrNotFound
	}
	if err != nil {
		return agent.RelaySchedule{}, err
	}
	rs.ID = agent.RelayScheduleID(rid)
	rs.Sets[0].ShiftKind = agent.ShiftKind(sk0)
	rs.Sets[0].WorkingDay.NecessaryLabourMinutes = agent.NecessaryLabourMinutes(nl0)
	rs.Sets[0].WorkingDay.SurplusLabourMinutes = agent.SurplusLabourMinutes(sl0)
	rs.Sets[1].ShiftKind = agent.ShiftKind(sk1)
	rs.Sets[1].WorkingDay.NecessaryLabourMinutes = agent.NecessaryLabourMinutes(nl1)
	rs.Sets[1].WorkingDay.SurplusLabourMinutes = agent.SurplusLabourMinutes(sl1)
	if err := json.Unmarshal([]byte(wids0raw), &rs.Sets[0].WorkerIDs); err != nil {
		return agent.RelaySchedule{}, err
	}
	if err := json.Unmarshal([]byte(wids1raw), &rs.Sets[1].WorkerIDs); err != nil {
		return agent.RelaySchedule{}, err
	}
	return rs, nil
}
```

- [ ] **Step 2: Verify compile**

```bash
go build ./services/agent-service/... 2>&1
```
Expected: compile error only from `main.go` agentStore interface (not yet updated). All other packages should compile.

- [ ] **Step 3: Commit**

```bash
git add services/agent-service/internal/store/mysql.go
git commit -m "feat(ch10): MySQL implementation of WorkingDayStore"
```

---

## Task 6: HTTP Handler

**Files:**
- Modify: `services/agent-service/internal/transport/httpapi/handler.go`
- Create: `services/agent-service/internal/transport/httpapi/working_day_handler.go`

- [ ] **Step 1: Add WorkingDayStore to Handler**

In `handler.go`, update the `Handler` struct and `New` function:

Change the struct from:
```go
type Handler struct {
	Store              store.Store
	CircuitStore       store.CircuitStore
	LabourPowerStore   store.LabourPowerStore
	LabourProcessStore store.LabourProcessStore
	Logger             *slog.Logger
}
```
to:
```go
type Handler struct {
	Store              store.Store
	CircuitStore       store.CircuitStore
	LabourPowerStore   store.LabourPowerStore
	LabourProcessStore store.LabourProcessStore
	WorkingDayStore    store.WorkingDayStore
	Logger             *slog.Logger
}
```

Change `New` signature from:
```go
func New(s store.Store, cs store.CircuitStore, lps store.LabourPowerStore, lproc store.LabourProcessStore, logger *slog.Logger) *Handler {
```
to:
```go
func New(s store.Store, cs store.CircuitStore, lps store.LabourPowerStore, lproc store.LabourProcessStore, wds store.WorkingDayStore, logger *slog.Logger) *Handler {
```

Change the return statement to:
```go
return &Handler{Store: s, CircuitStore: cs, LabourPowerStore: lps, LabourProcessStore: lproc, WorkingDayStore: wds, Logger: logger}
```

Also update `writeAppError` to include the new sentinel errors. After `errors.Is(err, agent.ErrInvalidContract)`:
```go
errors.Is(err, agent.ErrWorkingDayExceedsPhysicalMax),
errors.Is(err, agent.ErrWorkingDayExceedsStatutoryLimit):
    writeError(w, http.StatusBadRequest, err.Error())
```

The full updated `writeAppError`:
```go
func writeAppError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, store.ErrAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, agent.ErrInsufficientFunds),
		errors.Is(err, agent.ErrNotCapitalist),
		errors.Is(err, agent.ErrWrongClass),
		errors.Is(err, agent.ErrInvalidProcess),
		errors.Is(err, agent.ErrInvalidContract),
		errors.Is(err, agent.ErrWorkingDayExceedsPhysicalMax),
		errors.Is(err, agent.ErrWorkingDayExceedsStatutoryLimit):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, err.Error())
	}
}
```

- [ ] **Step 2: Create working_day_handler.go**

Create `services/agent-service/internal/transport/httpapi/working_day_handler.go`:

```go
package httpapi

import (
	"net/http"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

type createWorkingDayRequest struct {
	NecessaryLabourMinutes int64  `json:"necessary_labour_minutes"`
	SurplusLabourMinutes   int64  `json:"surplus_labour_minutes"`
	StatutoryLimitMinutes  *int64 `json:"statutory_limit_minutes,omitempty"`
}

type workingDayResponse struct {
	WorkingDay          agent.WorkingDay `json:"working_day"`
	TotalMinutes        int64            `json:"total_minutes"`
	RateOfSurplusValue  float64          `json:"rate_of_surplus_value"`
	ExceedsStatutory    bool             `json:"exceeds_statutory,omitempty"`
}

func buildWorkingDayResponse(wd agent.WorkingDay, limit *int64) workingDayResponse {
	resp := workingDayResponse{
		WorkingDay:         wd,
		TotalMinutes:       wd.TotalMinutes(),
		RateOfSurplusValue: wd.RateOfSurplusValue(),
	}
	if limit != nil {
		resp.ExceedsStatutory = wd.TotalMinutes() > *limit
	}
	return resp
}

func (h *Handler) CreateWorkingDay(w http.ResponseWriter, r *http.Request) {
	var req createWorkingDayRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	wd := agent.WorkingDay{
		NecessaryLabourMinutes: agent.NecessaryLabourMinutes(req.NecessaryLabourMinutes),
		SurplusLabourMinutes:   agent.SurplusLabourMinutes(req.SurplusLabourMinutes),
	}
	if err := wd.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.StatutoryLimitMinutes != nil {
		c := agent.WorkingDayConstraint{StatutoryLimitMinutes: agent.StatutoryLimitMinutes(*req.StatutoryLimitMinutes)}
		if err := agent.ValidateConstraint(wd, c); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	saved, err := h.WorkingDayStore.CreateWorkingDay(r.Context(), wd)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, buildWorkingDayResponse(saved, req.StatutoryLimitMinutes))
}

func (h *Handler) GetWorkingDay(w http.ResponseWriter, r *http.Request) {
	id := agent.WorkingDayID(r.PathValue("id"))
	wd, err := h.WorkingDayStore.GetWorkingDay(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, buildWorkingDayResponse(wd, nil))
}

type validateWorkingDayRequest struct {
	NecessaryLabourMinutes int64  `json:"necessary_labour_minutes"`
	SurplusLabourMinutes   int64  `json:"surplus_labour_minutes"`
	StatutoryLimitMinutes  *int64 `json:"statutory_limit_minutes,omitempty"`
}

type validateWorkingDayResponse struct {
	TotalMinutes       int64   `json:"total_minutes"`
	RateOfSurplusValue float64 `json:"rate_of_surplus_value"`
	Valid              bool    `json:"valid"`
	Error              string  `json:"error,omitempty"`
}

func (h *Handler) ValidateWorkingDay(w http.ResponseWriter, r *http.Request) {
	var req validateWorkingDayRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	wd := agent.WorkingDay{
		NecessaryLabourMinutes: agent.NecessaryLabourMinutes(req.NecessaryLabourMinutes),
		SurplusLabourMinutes:   agent.SurplusLabourMinutes(req.SurplusLabourMinutes),
	}
	resp := validateWorkingDayResponse{
		TotalMinutes: wd.TotalMinutes(),
	}
	if wd.NecessaryLabourMinutes > 0 {
		resp.RateOfSurplusValue = wd.RateOfSurplusValue()
	}
	if err := wd.Validate(); err != nil {
		resp.Error = err.Error()
		writeJSON(w, http.StatusOK, resp)
		return
	}
	if req.StatutoryLimitMinutes != nil {
		c := agent.WorkingDayConstraint{StatutoryLimitMinutes: agent.StatutoryLimitMinutes(*req.StatutoryLimitMinutes)}
		if err := agent.ValidateConstraint(wd, c); err != nil {
			resp.Error = err.Error()
			writeJSON(w, http.StatusOK, resp)
			return
		}
	}
	resp.Valid = true
	writeJSON(w, http.StatusOK, resp)
}

type relaySetInput struct {
	ShiftKind              string    `json:"shift_kind"`
	NecessaryLabourMinutes int64     `json:"necessary_labour_minutes"`
	SurplusLabourMinutes   int64     `json:"surplus_labour_minutes"`
	WorkerIDs              []string  `json:"worker_ids"`
}

type createRelayScheduleRequest struct {
	Sets [2]relaySetInput `json:"sets"`
}

func (h *Handler) CreateRelaySchedule(w http.ResponseWriter, r *http.Request) {
	var req createRelayScheduleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	rs := agent.RelaySchedule{}
	for i, s := range req.Sets {
		wids := make([]agent.AgentID, len(s.WorkerIDs))
		for j, id := range s.WorkerIDs {
			wids[j] = agent.AgentID(id)
		}
		rs.Sets[i] = agent.RelaySet{
			ShiftKind: agent.ShiftKind(s.ShiftKind),
			WorkingDay: agent.WorkingDay{
				NecessaryLabourMinutes: agent.NecessaryLabourMinutes(s.NecessaryLabourMinutes),
				SurplusLabourMinutes:   agent.SurplusLabourMinutes(s.SurplusLabourMinutes),
			},
			WorkerIDs: wids,
		}
	}
	if err := rs.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	saved, err := h.WorkingDayStore.CreateRelaySchedule(r.Context(), rs)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, saved)
}

func (h *Handler) GetRelaySchedule(w http.ResponseWriter, r *http.Request) {
	id := agent.RelayScheduleID(r.PathValue("id"))
	rs, err := h.WorkingDayStore.GetRelaySchedule(r.Context(), id)
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, rs)
}
```

- [ ] **Step 3: Verify compile**

```bash
go build ./services/agent-service/internal/transport/... 2>&1
```
Expected: compile error because `main.go` still calls old `New(...)` signature. We'll fix that in Task 8.

---

## Task 7: Register Routes

**Files:**
- Modify: `services/agent-service/internal/transport/httpapi/routes.go`

- [ ] **Step 1: Add Ch. 10 routes**

Append to `Register` in `routes.go`:

```go
	// Ch. 10 — The Working-Day
	s.HandleFunc("POST /v1/working-days", h.CreateWorkingDay)
	s.HandleFunc("GET /v1/working-days/{id}", h.GetWorkingDay)
	s.HandleFunc("POST /v1/working-days/validate", h.ValidateWorkingDay)
	s.HandleFunc("POST /v1/relay-schedules", h.CreateRelaySchedule)
	s.HandleFunc("GET /v1/relay-schedules/{id}", h.GetRelaySchedule)
```

---

## Task 8: Wire main.go and API Gateway

**Files:**
- Modify: `services/agent-service/cmd/agent-service/main.go`
- Modify: `services/api-gateway/cmd/api-gateway/main.go`

- [ ] **Step 1: Update agentStore interface in main.go**

Change the `agentStore` interface to embed `WorkingDayStore`:

```go
type agentStore interface {
	store.Store
	store.CircuitStore
	store.LabourPowerStore
	store.LabourProcessStore
	store.WorkingDayStore
}
```

Update the `httpapi.Register` call to pass the new argument:

```go
httpapi.Register(srv, httpapi.New(st, st, st, st, st, logger))
```

- [ ] **Step 2: Update API gateway to add Ch.10 proxy routes**

In `services/api-gateway/cmd/api-gateway/main.go`:

After the Ch. 7 comment block, add:

```go
	// Ch. 10 — working-day routes proxy to agent-service
	srv.Handle("/v1/working-days", agentProxy)
	srv.Handle("/v1/working-days/{rest...}", agentProxy)
	srv.Handle("/v1/relay-schedules", agentProxy)
	srv.Handle("/v1/relay-schedules/{rest...}", agentProxy)
```

Also update `handleInfo` status string to `"ch-10-working-day"` and chapter description to `"Capital Vol. I, Ch. 10 - The Working-Day"`.

- [ ] **Step 3: Verify the whole module compiles**

```bash
go build ./... 2>&1
```
Expected: no errors.

- [ ] **Step 4: Run all agent-service tests**

```bash
go test ./services/agent-service/... -count=1 2>&1
```
Expected: all pass.

- [ ] **Step 5: Commit**

```bash
git add services/agent-service/internal/transport/httpapi/handler.go \
        services/agent-service/internal/transport/httpapi/working_day_handler.go \
        services/agent-service/internal/transport/httpapi/routes.go \
        services/agent-service/cmd/agent-service/main.go \
        services/api-gateway/cmd/api-gateway/main.go
git commit -m "feat(ch10): HTTP endpoints for working-days and relay-schedules"
```

- [ ] **Step 6: Full build + vet**

Ask user to run:
```bash
make vet test build
```
Expected: all pass.

---

## Task 9: React Frontend

**Files:**
- Modify: `web/src/types.ts`
- Modify: `web/src/api.ts`
- Create: `web/src/chapters/Ch10WorkingDay.tsx`
- Modify: `web/src/components/ChapterShell.tsx`
- Modify: `web/src/chapters/registry.ts`

- [ ] **Step 1: Add Ch.10 types to types.ts**

Append after the Ch.9 section at the end of `web/src/types.ts`:

```typescript
// --- agent-service types (Ch. 10: The Working-Day) ---------------------------

export interface WorkingDay {
  id: string;
  necessary_labour_minutes: number;
  surplus_labour_minutes: number;
  created_at: string;
}

export interface WorkingDayResponse {
  working_day: WorkingDay;
  total_minutes: number;
  rate_of_surplus_value: number;
  exceeds_statutory?: boolean;
}

export interface ValidateWorkingDayResponse {
  total_minutes: number;
  rate_of_surplus_value: number;
  valid: boolean;
  error?: string;
}

export interface RelaySet {
  shift_kind: "day" | "night";
  working_day: WorkingDay;
  worker_ids: string[];
}

export interface RelaySchedule {
  id: string;
  sets: [RelaySet, RelaySet];
  created_at: string;
}

export interface CreateWorkingDayInput {
  necessary_labour_minutes: number;
  surplus_labour_minutes: number;
  statutory_limit_minutes?: number;
}

export interface ValidateWorkingDayInput {
  necessary_labour_minutes: number;
  surplus_labour_minutes: number;
  statutory_limit_minutes?: number;
}

export interface RelaySetInput {
  shift_kind: "day" | "night";
  necessary_labour_minutes: number;
  surplus_labour_minutes: number;
  worker_ids: string[];
}

export interface CreateRelayScheduleInput {
  sets: [RelaySetInput, RelaySetInput];
}
```

- [ ] **Step 2: Add Ch.10 API methods to api.ts**

Add to the `import type` block in `api.ts` (add to the existing import list):
```typescript
  CreateWorkingDayInput,
  ValidateWorkingDayInput,
  ValidateWorkingDayResponse,
  WorkingDayResponse,
  CreateRelayScheduleInput,
  RelaySchedule,
```

Append to the `api` object after `computeRateOfSurplusValue`:

```typescript
  // --- agent-service (Ch. 10: The Working-Day) ---

  createWorkingDay: (input: CreateWorkingDayInput) =>
    http<WorkingDayResponse>("/v1/working-days", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getWorkingDay: (id: string) =>
    http<WorkingDayResponse>(`/v1/working-days/${id}`),

  validateWorkingDay: (input: ValidateWorkingDayInput) =>
    http<ValidateWorkingDayResponse>("/v1/working-days/validate", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  createRelaySchedule: (input: CreateRelayScheduleInput) =>
    http<RelaySchedule>("/v1/relay-schedules", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getRelaySchedule: (id: string) =>
    http<RelaySchedule>(`/v1/relay-schedules/${id}`),
```

- [ ] **Step 3: Create Ch10WorkingDay.tsx**

Create `web/src/chapters/Ch10WorkingDay.tsx`:

```tsx
import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type {
  WorkingDayResponse,
  ValidateWorkingDayResponse,
  RelaySchedule,
} from "../types";

interface Ch10Props {
  onSharedChanged: () => void;
}

function minutesToHours(m: number): string {
  const h = Math.floor(m / 60);
  const min = m % 60;
  return min === 0 ? `${h}h` : `${h}h ${min}m`;
}

export function Ch10WorkingDay({ onSharedChanged: _onSharedChanged }: Ch10Props) {
  return (
    <>
      <ValidatePanel />
      <CreatePanel />
      <RelaySchedulePanel />
    </>
  );
}

function ValidatePanel() {
  const [nl, setNl] = useState(360);
  const [sl, setSl] = useState(360);
  const [limit, setLimit] = useState<number | "">("");
  const [result, setResult] = useState<ValidateWorkingDayResponse | null>(null);
  const [err, setErr] = useState<string | null>(null);

  function loadFixture(necessary: number, surplus: number, lim?: number) {
    setNl(necessary);
    setSl(surplus);
    setLimit(lim ?? "");
    setResult(null);
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const input = {
        necessary_labour_minutes: nl,
        surplus_labour_minutes: sl,
        ...(limit !== "" ? { statutory_limit_minutes: Number(limit) } : {}),
      };
      const r = await api.validateWorkingDay(input);
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Working-Day Validator</h2>
      <p className="description">
        Validate a working day against the physical maximum (24 h) and an optional statutory
        limit. Computes rate of surplus-value from the A–B / B–C segments.
      </p>
      <div style={{ marginBottom: "0.75rem", display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
        <button type="button" onClick={() => loadFixture(360, 60)}>
          §1 WD I (6h+1h)
        </button>
        <button type="button" onClick={() => loadFixture(360, 180)}>
          §1 WD II (6h+3h)
        </button>
        <button type="button" onClick={() => loadFixture(360, 360)}>
          §1 WD III (6h+6h)
        </button>
        <button type="button" onClick={() => loadFixture(315, 315, 630)}>
          Factory Act 1850 (10.5 h)
        </button>
        <button type="button" onClick={() => loadFixture(84 * 480, 56 * 480)}>
          Wallachian Corvée
        </button>
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Necessary Labour A–B (min)</span>
          <input type="number" min={1} value={nl} onChange={(e) => setNl(Number(e.target.value))} />
        </label>
        <label>
          <span>Surplus Labour B–C (min)</span>
          <input type="number" min={0} value={sl} onChange={(e) => setSl(Number(e.target.value))} />
        </label>
        <label>
          <span>Statutory Limit (min, optional)</span>
          <input
            type="number"
            min={1}
            value={limit}
            placeholder="e.g. 630"
            onChange={(e) => setLimit(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Validate</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <table className="data-table">
          <tbody>
            <tr><td>Total</td><td>{minutesToHours(result.total_minutes)}</td></tr>
            <tr>
              <td>Rate of Surplus-Value (s/v)</td>
              <td><strong>{(result.rate_of_surplus_value * 100).toFixed(2)}%</strong></td>
            </tr>
            <tr>
              <td>Valid</td>
              <td>{result.valid ? "Yes" : <span className="error">No — {result.error}</span>}</td>
            </tr>
          </tbody>
        </table>
      )}
    </section>
  );
}

function CreatePanel() {
  const [nl, setNl] = useState(360);
  const [sl, setSl] = useState(360);
  const [limit, setLimit] = useState<number | "">("");
  const [result, setResult] = useState<WorkingDayResponse | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const input = {
        necessary_labour_minutes: nl,
        surplus_labour_minutes: sl,
        ...(limit !== "" ? { statutory_limit_minutes: Number(limit) } : {}),
      };
      const r = await api.createWorkingDay(input);
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Create Working-Day Record</h2>
      <p className="description">
        Persist a working-day record. Returns the stored ID alongside computed metrics.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Necessary Labour A–B (min)</span>
          <input type="number" min={1} value={nl} onChange={(e) => setNl(Number(e.target.value))} />
        </label>
        <label>
          <span>Surplus Labour B–C (min)</span>
          <input type="number" min={0} value={sl} onChange={(e) => setSl(Number(e.target.value))} />
        </label>
        <label>
          <span>Statutory Limit (min, optional)</span>
          <input
            type="number"
            min={1}
            value={limit}
            placeholder="e.g. 630"
            onChange={(e) => setLimit(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Create</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <div>
          <p className="small muted">Saved: <span className="monospace">{result.working_day.id}</span></p>
          <table className="data-table">
            <tbody>
              <tr><td>Total</td><td>{minutesToHours(result.total_minutes)}</td></tr>
              <tr>
                <td>Rate of Surplus-Value (s/v)</td>
                <td><strong>{(result.rate_of_surplus_value * 100).toFixed(2)}%</strong></td>
              </tr>
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}

function RelaySchedulePanel() {
  const [nl0, setNl0] = useState(360);
  const [sl0, setSl0] = useState(360);
  const [nl1, setNl1] = useState(360);
  const [sl1, setSl1] = useState(0);
  const [result, setResult] = useState<RelaySchedule | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.createRelaySchedule({
        sets: [
          { shift_kind: "day", necessary_labour_minutes: nl0, surplus_labour_minutes: sl0, worker_ids: [] },
          { shift_kind: "night", necessary_labour_minutes: nl1, surplus_labour_minutes: sl1, worker_ids: [] },
        ],
      });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  const combinedMin = nl0 + sl0 + nl1 + sl1;

  return (
    <section className="card">
      <h2>Relay Schedule</h2>
      <p className="description">
        A relay system alternates two worker sets to keep machinery running. Combined shifts must
        not exceed 24 h (1440 min).
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Day shift — Necessary (min)</span>
          <input type="number" min={1} value={nl0} onChange={(e) => setNl0(Number(e.target.value))} />
        </label>
        <label>
          <span>Day shift — Surplus (min)</span>
          <input type="number" min={0} value={sl0} onChange={(e) => setSl0(Number(e.target.value))} />
        </label>
        <label>
          <span>Night shift — Necessary (min)</span>
          <input type="number" min={1} value={nl1} onChange={(e) => setNl1(Number(e.target.value))} />
        </label>
        <label>
          <span>Night shift — Surplus (min)</span>
          <input type="number" min={0} value={sl1} onChange={(e) => setSl1(Number(e.target.value))} />
        </label>
        <div className="form-actions span2">
          <span className="small muted">
            Combined: {minutesToHours(combinedMin)} / 24h max
            {combinedMin > 1440 && <span className="error"> — exceeds physical max</span>}
          </span>
          <button type="submit" className="primary">Create Relay Schedule</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <div>
          <p className="small muted">Saved: <span className="monospace">{result.id}</span></p>
          <table className="data-table">
            <thead>
              <tr><th>Shift</th><th>Necessary</th><th>Surplus</th><th>Total</th></tr>
            </thead>
            <tbody>
              {result.sets.map((s, i) => (
                <tr key={i}>
                  <td>{s.shift_kind}</td>
                  <td>{minutesToHours(s.working_day.necessary_labour_minutes)}</td>
                  <td>{minutesToHours(s.working_day.surplus_labour_minutes)}</td>
                  <td>
                    {minutesToHours(
                      s.working_day.necessary_labour_minutes + s.working_day.surplus_labour_minutes
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </section>
  );
}
```

- [ ] **Step 4: Wire Ch10 into ChapterShell.tsx**

In `web/src/components/ChapterShell.tsx`:

Add import at top with other chapter imports:
```tsx
import { Ch10WorkingDay } from "../chapters/Ch10WorkingDay";
```

Add to `QUOTES`:
```tsx
ch10: "The capitalist maintains his rights as a purchaser when he tries to make the working-day as long as possible.",
```

Change `} : null}` at the end to:
```tsx
        } : activeChapterId === "ch10" ? (
          <Ch10WorkingDay onSharedChanged={onSharedChanged} />
        ) : null}
```

- [ ] **Step 5: Update registry.ts**

In `web/src/chapters/registry.ts`, change ch10 status from `"pending"` to `"done"`:

```typescript
  { id: "ch10", number: 10, title: "The Working Day",  part: "Part III — The Production of Absolute Surplus-Value",  status: "done" },
```

- [ ] **Step 6: Verify TypeScript**

Ask user to run:
```bash
cd web && npm run lint
```
Expected: no errors.

- [ ] **Step 7: Commit**

```bash
git add web/src/types.ts \
        web/src/api.ts \
        web/src/chapters/Ch10WorkingDay.tsx \
        web/src/components/ChapterShell.tsx \
        web/src/chapters/registry.ts
git commit -m "feat(ch10): React panel — working-day builder and relay schedule visualiser"
```

---

## Task 10: Update Architecture Docs

**Files:**
- Modify: `docs/architecture.md`

- [ ] **Step 1: Update the chapter table in architecture.md**

Find the ch10 row and update its status from `pending` to `done` with the service and description. The exact content of `docs/architecture.md` needs to be read first — update the row for Chapter 10 to reflect `agent-service`, endpoints `POST/GET /v1/working-days`, `POST /v1/working-days/validate`, `POST/GET /v1/relay-schedules`.

- [ ] **Step 2: Commit**

```bash
git add docs/architecture.md
git commit -m "docs: mark Ch.10 done in architecture table"
```

---

## Self-Review

### Spec coverage check

| Spec requirement | Covered by |
|---|---|
| `WorkingDay`, `NecessaryLabourMinutes`, `SurplusLabourMinutes` types | Task 1 |
| `TotalMinutes()`, `RateOfSurplusValue()`, `Validate()` methods | Task 1 |
| `PhysicalMaxMinutes = 1440`, `ErrWorkingDayExceedsPhysicalMax` | Task 1 |
| `StatutoryLimitMinutes`, `WorkingDayConstraint`, `ValidateConstraint` | Task 1 |
| `ErrWorkingDayExceedsStatutoryLimit` | Task 1 |
| `NormalWorkingDay` | Task 1 |
| `Overwork.AnnualMinutes()` | Task 1 |
| `FactoryAct{Year, ChildLimit, YoungPersonLimit}` | Task 1 |
| `ShiftKind` ("day"/"night"), `TimerClass` ("full"/"half") | Task 1 |
| `RelaySet`, `RelaySchedule`, `RelaySchedule.Validate()` | Task 1 |
| All §1-§6 fixtures as tests | Task 1 |
| Store interface `WorkingDayStore` | Task 2 |
| Memory store + tests | Task 3 |
| MySQL migration `00008_ch10_working_day.sql` | Task 4 |
| MySQL implementation | Task 5 |
| `POST /v1/working-days` | Task 6 |
| `GET /v1/working-days/{id}` | Task 6 |
| `POST /v1/working-days/validate` | Task 6 |
| `POST /v1/relay-schedules` | Task 6 |
| `GET /v1/relay-schedules/{id}` | Task 6 |
| API gateway routing | Task 8 |
| React panel: working-day builder + relay visualiser | Task 9 |
| `docs/architecture.md` update | Task 10 |

All spec requirements are covered. No gaps found.

### Placeholder scan

No TBDs, TODOs, or "similar to Task N" patterns present.

### Type consistency

- `NecessaryLabourMinutes` and `SurplusLabourMinutes` named types used consistently in domain, store, and handler layers.
- `WorkingDayID` / `RelayScheduleID` / `AgentID` typed IDs used consistently.
- `ShiftDay` / `ShiftNight` constants used in tests and domain.
- `ValidateConstraint` function name consistent across all tasks.
- `buildWorkingDayResponse` helper used in both `CreateWorkingDay` and `GetWorkingDay`.
