# Ch. 11 — Rate and Mass of Surplus Value: Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement `simulation-engine/internal/surplus` domain types and two HTTP endpoints (`POST /v1/surplus/mass`, `GET /v1/surplus/limits`), wire them through the API gateway, and build the Ch. 11 React panel illustrating the three laws of surplus-value mass.

**Architecture:** A new `surplus` package in `simulation-engine` provides pure-computation domain types and functions (no store needed). The simulation-engine main.go wires a new transport handler. The API gateway proxies `/v1/surplus/...` to simulation-engine (port 8084). The React frontend adds a `Ch11RateAndMassOfSurplusValue.tsx` panel with fixture presets.

**Tech Stack:** Go 1.25, `pkg/httpx`, React 18 + Vite + TypeScript

---

## ⚠️ Spec Note — Fixture #3

The spec fixture `MassByRate(SurplusValueRate{SurplusLabour: 12, NecessaryLabour: 6}, VariableCapital(150)) == SurplusValueMass(150)` is incorrect. The formula `S = (s/v) × V = SurplusLabour × V / NecessaryLabour` yields `12 × 150 / 6 = 300`, which is also confirmed by Marx's text (doubling the rate while halving V keeps S constant at 300s, per the compensation law). The correct fixture value is `SurplusValueMass(300)`. Implement the correct formula; note the discrepancy in the test.

---

## File Map

| Action | Path | Responsibility |
|--------|------|----------------|
| Create | `services/simulation-engine/internal/surplus/surplus.go` | Domain types, const, functions |
| Create | `services/simulation-engine/internal/surplus/surplus_test.go` | Marx textual fixture tests + invariants |
| Create | `services/simulation-engine/internal/transport/httpapi/handler.go` | HTTP handlers for mass and limits |
| Create | `services/simulation-engine/internal/transport/httpapi/routes.go` | Route registration helper |
| Modify | `services/simulation-engine/cmd/simulation-engine/main.go` | Wire transport, register routes |
| Modify | `services/api-gateway/cmd/api-gateway/main.go` | Proxy `/v1/surplus/...` to simulation-engine |
| Modify | `web/src/types.ts` | Add Ch. 11 TypeScript types |
| Modify | `web/src/api.ts` | Add Ch. 11 API methods |
| Create | `web/src/chapters/Ch11RateAndMassOfSurplusValue.tsx` | React panel |
| Modify | `web/src/components/ChapterShell.tsx` | Wire ch11 case |
| Modify | `web/src/chapters/registry.ts` | Mark ch11 as done |
| Modify | `docs/architecture.md` | Update roadmap row |

---

## Task 1: Domain Types and Functions

**Files:**
- Create: `services/simulation-engine/internal/surplus/surplus.go`

- [ ] **Step 1: Create the domain file**

```go
// Package surplus implements the rate and mass of surplus-value from
// Capital Vol. I, Ch. 11. All functions are pure; no persistence is needed.
package surplus

// LabourMinutes is the canonical value-magnitude unit, reused from agent-service
// semantics (identical definition, separate package to avoid import cycles).
type LabourMinutes int64

// AbsoluteWorkdayLimit is the physical ceiling of any working day: 24 h × 60 min.
const AbsoluteWorkdayLimit LabourMinutes = 24 * 60

// SurplusValueMass is the total surplus-value extracted from all simultaneously
// employed workers. Expressed in the same abstract unit as VariableCapital (§1).
type SurplusValueMass int64

// SurplusValueRate is the ratio of surplus-labour to necessary-labour, both in
// LabourMinutes. Corresponds to Marx's s/v or a′/a notation [§1].
type SurplusValueRate struct {
	SurplusLabour   int64 `json:"surplus_labour"`
	NecessaryLabour int64 `json:"necessary_labour"`
}

// Rate returns the dimensionless ratio s/v as a float64.
func (r SurplusValueRate) Rate() float64 {
	if r.NecessaryLabour == 0 {
		return 0
	}
	return float64(r.SurplusLabour) / float64(r.NecessaryLabour)
}

// VariableCapital is the total money-value advanced for labour-power across all
// simultaneously employed workers [§1]. Expressed in the same integer unit as
// SurplusValueMass so that S = rate × V holds exactly.
type VariableCapital int64

// LabourPowerValue is the daily cost of reproducing a single worker [§1].
type LabourPowerValue int64

// WorkerCount is the number of simultaneously employed workers [§1].
type WorkerCount int

// SurplusValueSnapshot carries the aggregate calculation result returned by
// POST /v1/surplus/mass. Both formula results are included for cross-validation.
type SurplusValueSnapshot struct {
	Rate          SurplusValueRate `json:"rate"`
	VariableCapital VariableCapital `json:"variable_capital"`
	WorkerCount   WorkerCount      `json:"worker_count"`
	Mass          SurplusValueMass `json:"mass"`
	MassByRate    SurplusValueMass `json:"mass_by_rate"`
	MassByWorkers SurplusValueMass `json:"mass_by_workers"`
}

// IndividualSurplus returns the surplus-value produced by a single worker
// given the rate and the value of their labour-power [§1].
// s_individual = v × (a′/a) = v × SurplusLabour / NecessaryLabour
func IndividualSurplus(rate SurplusValueRate, v LabourPowerValue) SurplusValueMass {
	if rate.NecessaryLabour == 0 {
		return 0
	}
	return SurplusValueMass(int64(v) * rate.SurplusLabour / rate.NecessaryLabour)
}

// MassByRate computes S = (s/v) × V — total surplus-value from total variable
// capital and the rate of surplus-value [§1, formula I].
func MassByRate(rate SurplusValueRate, totalVariableCapital VariableCapital) SurplusValueMass {
	if rate.NecessaryLabour == 0 {
		return 0
	}
	return SurplusValueMass(rate.SurplusLabour * int64(totalVariableCapital) / rate.NecessaryLabour)
}

// MassByWorkers computes S = P × (a′/a) × n — total surplus-value from per-worker
// labour-power value, the rate, and the number of workers [§1, formula II].
func MassByWorkers(v LabourPowerValue, rate SurplusValueRate, n WorkerCount) SurplusValueMass {
	if rate.NecessaryLabour == 0 {
		return 0
	}
	return SurplusValueMass(int64(v) * rate.SurplusLabour * int64(n) / rate.NecessaryLabour)
}

// MinimumCapital returns the minimum variable capital required to employ n workers
// each at daily reproduction cost v [§1].
// V_min = v × n
func MinimumCapital(v LabourPowerValue, n WorkerCount) VariableCapital {
	return VariableCapital(int64(v) * int64(n))
}
```

- [ ] **Step 2: Confirm the file compiles (no build tool available in sandbox)**

No tool to run. Ask the user to run `make vet` after all Go files are written.

---

## Task 2: Domain Tests

**Files:**
- Create: `services/simulation-engine/internal/surplus/surplus_test.go`

- [ ] **Step 1: Write the test file**

```go
package surplus_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/surplus"
)

func TestAbsoluteWorkdayLimit(t *testing.T) {
	t.Parallel()
	if surplus.AbsoluteWorkdayLimit != 1440 {
		t.Fatalf("expected 1440, got %d", surplus.AbsoluteWorkdayLimit)
	}
}

func TestSurplusValueRate_Rate(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	if got := rate.Rate(); got != 1.0 {
		t.Fatalf("expected 1.0, got %f", got)
	}
	rate200 := surplus.SurplusValueRate{SurplusLabour: 12, NecessaryLabour: 6}
	if got := rate200.Rate(); got != 2.0 {
		t.Fatalf("expected 2.0, got %f", got)
	}
}

// §1 fixture: "if the rate of surplus-value be = 100%, this variable capital of
// 3s. produces a mass of surplus-value of 3s."
func TestMassByRate_100pct(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	got := surplus.MassByRate(rate, surplus.VariableCapital(3))
	if got != surplus.SurplusValueMass(3) {
		t.Fatalf("expected 3, got %d", got)
	}
}

// §1 fixture: 100 labourers at rate 100%, V=300s → S=300s
func TestMassByWorkers_100workers(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	got := surplus.MassByWorkers(surplus.LabourPowerValue(3), rate, surplus.WorkerCount(100))
	if got != surplus.SurplusValueMass(300) {
		t.Fatalf("expected 300, got %d", got)
	}
}

// §1 compensation law: rate doubles (100%→200%), V halves (300s→150s), workers halve (100→50).
// S stays constant at 300 (not 150 — the spec fixture has a typo; the formula is correct).
func TestMassByRate_CompensationLaw(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 12, NecessaryLabour: 6}
	got := surplus.MassByRate(rate, surplus.VariableCapital(150))
	// S = (12/6) × 150 = 2 × 150 = 300
	if got != surplus.SurplusValueMass(300) {
		t.Fatalf("expected 300, got %d", got)
	}
}

// §1 fixture: V=1500s, rate=100% → S=1500s
func TestMassByRate_Large(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	got := surplus.MassByRate(rate, surplus.VariableCapital(1500))
	if got != surplus.SurplusValueMass(1500) {
		t.Fatalf("expected 1500, got %d", got)
	}
}

// §1 fixture: "A capital of 300s. that employs 100 labourers a day with a rate of
// surplus-value of 200% … produces only a mass of surplus-value of 600s."
func TestMassByRate_200pct(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 12, NecessaryLabour: 6}
	got := surplus.MassByRate(rate, surplus.VariableCapital(300))
	if got != surplus.SurplusValueMass(600) {
		t.Fatalf("expected 600, got %d", got)
	}
}

// §1 minimum capital: "he would have to employ two labourers in order to live …"
func TestMinimumCapital(t *testing.T) {
	t.Parallel()
	got := surplus.MinimumCapital(surplus.LabourPowerValue(3), surplus.WorkerCount(2))
	if got != surplus.VariableCapital(6) {
		t.Fatalf("expected 6, got %d", got)
	}
}

// Invariant: MassByRate and MassByWorkers agree when V = v × n [§1]
func TestInvariant_FormulaAgreement(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	v := surplus.LabourPowerValue(3)
	n := surplus.WorkerCount(100)
	V := surplus.VariableCapital(int64(v) * int64(n)) // 300
	byRate := surplus.MassByRate(rate, V)
	byWorkers := surplus.MassByWorkers(v, rate, n)
	if byRate != byWorkers {
		t.Fatalf("formula disagreement: MassByRate=%d MassByWorkers=%d", byRate, byWorkers)
	}
}

// Invariant: IndividualSurplus × n == MassByWorkers [§1]
func TestInvariant_IndividualScaled(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	v := surplus.LabourPowerValue(3)
	n := surplus.WorkerCount(100)
	individual := surplus.IndividualSurplus(rate, v)
	mass := surplus.MassByWorkers(v, rate, n)
	if surplus.SurplusValueMass(int64(individual)*int64(n)) != mass {
		t.Fatalf("scale invariant failed: %d × %d ≠ %d", individual, n, mass)
	}
}

// Invariant: working day (necessary + surplus) < AbsoluteWorkdayLimit [§1]
func TestInvariant_WorkingDayBelowAbsoluteLimit(t *testing.T) {
	t.Parallel()
	rate := surplus.SurplusValueRate{SurplusLabour: 6, NecessaryLabour: 6}
	total := surplus.LabourMinutes(rate.NecessaryLabour + rate.SurplusLabour)
	if total >= surplus.AbsoluteWorkdayLimit {
		t.Fatalf("working day %d must be < %d", total, surplus.AbsoluteWorkdayLimit)
	}
}
```

- [ ] **Step 2: Commit the domain layer**

```bash
git add services/simulation-engine/internal/surplus/
git commit -m "feat(ch11): surplus domain types and Marx textual fixture tests"
```

---

## Task 3: HTTP Transport

**Files:**
- Create: `services/simulation-engine/internal/transport/httpapi/handler.go`
- Create: `services/simulation-engine/internal/transport/httpapi/routes.go`

- [ ] **Step 1: Create handler.go**

```go
package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/surplus"
)

// Handler holds the logger; all surplus endpoints are stateless.
type Handler struct {
	Logger *slog.Logger
}

func New(logger *slog.Logger) *Handler {
	if logger == nil {
		logger = slog.Default()
	}
	return &Handler{Logger: logger}
}

// massRequest accepts either {rate, variable_capital} or
// {rate, labour_power_value, worker_count}, or all fields for cross-validation.
type massRequest struct {
	SurplusLabour    int64  `json:"surplus_labour"`
	NecessaryLabour  int64  `json:"necessary_labour"`
	VariableCapital  *int64 `json:"variable_capital,omitempty"`
	LabourPowerValue *int64 `json:"labour_power_value,omitempty"`
	WorkerCount      *int   `json:"worker_count,omitempty"`
}

// ComputeMass handles POST /v1/surplus/mass.
func (h *Handler) ComputeMass(w http.ResponseWriter, r *http.Request) {
	var req massRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.NecessaryLabour <= 0 {
		writeError(w, http.StatusBadRequest, "necessary_labour must be positive")
		return
	}
	if req.SurplusLabour < 0 {
		writeError(w, http.StatusBadRequest, "surplus_labour cannot be negative")
		return
	}
	if req.VariableCapital == nil && (req.LabourPowerValue == nil || req.WorkerCount == nil) {
		writeError(w, http.StatusBadRequest, "provide variable_capital or both labour_power_value and worker_count")
		return
	}

	rate := surplus.SurplusValueRate{
		SurplusLabour:   req.SurplusLabour,
		NecessaryLabour: req.NecessaryLabour,
	}

	var snap surplus.SurplusValueSnapshot
	snap.Rate = rate

	if req.VariableCapital != nil {
		vc := surplus.VariableCapital(*req.VariableCapital)
		snap.VariableCapital = vc
		snap.MassByRate = surplus.MassByRate(rate, vc)
	}
	if req.LabourPowerValue != nil && req.WorkerCount != nil {
		v := surplus.LabourPowerValue(*req.LabourPowerValue)
		n := surplus.WorkerCount(*req.WorkerCount)
		snap.WorkerCount = n
		snap.MassByWorkers = surplus.MassByWorkers(v, rate, n)
		// Derive variable_capital if not supplied
		if req.VariableCapital == nil {
			vc := surplus.MinimumCapital(v, n)
			snap.VariableCapital = vc
			snap.MassByRate = surplus.MassByRate(rate, vc)
		}
	}

	// Primary mass: prefer worker-count formula when both available, else rate formula.
	if req.LabourPowerValue != nil && req.WorkerCount != nil {
		snap.Mass = snap.MassByWorkers
	} else {
		snap.Mass = snap.MassByRate
	}

	writeJSON(w, http.StatusOK, snap)
}

// limitsResponse is the GET /v1/surplus/limits response shape.
type limitsResponse struct {
	AbsoluteWorkdayLimit int64  `json:"absolute_workday_limit"`
	MinimumCapital       *int64 `json:"minimum_capital,omitempty"`
	LabourPowerValue     *int64 `json:"labour_power_value,omitempty"`
	WorkerCount          *int   `json:"worker_count,omitempty"`
}

// GetLimits handles GET /v1/surplus/limits.
func (h *Handler) GetLimits(w http.ResponseWriter, r *http.Request) {
	limit := int64(surplus.AbsoluteWorkdayLimit)
	resp := limitsResponse{AbsoluteWorkdayLimit: limit}

	lpvStr := r.URL.Query().Get("labour_power_value")
	wcStr := r.URL.Query().Get("worker_count")
	if lpvStr != "" && wcStr != "" {
		lpv, err1 := strconv.ParseInt(lpvStr, 10, 64)
		wc, err2 := strconv.Atoi(wcStr)
		if err1 != nil || err2 != nil || lpv <= 0 || wc <= 0 {
			writeError(w, http.StatusBadRequest, "labour_power_value and worker_count must be positive integers")
			return
		}
		mc := int64(surplus.MinimumCapital(surplus.LabourPowerValue(lpv), surplus.WorkerCount(wc)))
		resp.MinimumCapital = &mc
		resp.LabourPowerValue = &lpv
		resp.WorkerCount = &wc
	}

	writeJSON(w, http.StatusOK, resp)
}

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
```

- [ ] **Step 2: Create routes.go**

```go
package httpapi

import "github.com/theding0x/capital-simulator/pkg/httpx"

func Register(s *httpx.Server, h *Handler) {
	s.HandleFunc("POST /v1/surplus/mass", h.ComputeMass)
	s.HandleFunc("GET /v1/surplus/limits", h.GetLimits)
}
```

- [ ] **Step 3: Commit the transport layer**

```bash
git add services/simulation-engine/internal/transport/
git commit -m "feat(ch11): HTTP transport for surplus mass and limits endpoints"
```

---

## Task 4: Wire Simulation-Engine main.go

**Files:**
- Modify: `services/simulation-engine/cmd/simulation-engine/main.go`

- [ ] **Step 1: Replace main.go with wired version**

```go
// simulation-engine is the time-step orchestrator. It advances the simulated
// economy one period at a time, telling the domain services to produce,
// exchange, and accumulate. Without it, the rest of the world is static.
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/theding0x/capital-simulator/pkg/httpx"
	applog "github.com/theding0x/capital-simulator/pkg/log"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

const serviceName = "simulation-engine"

func main() {
	logger := applog.New(serviceName)
	applog.SetDefault(logger)

	addr := getenv("SERVICE_ADDR", ":8084")
	srv := httpx.New(httpx.Config{Addr: addr}, logger)

	srv.HandleFunc("/v1/sim/status", handleStatus)

	// Ch. 11 — Rate and Mass of Surplus-Value
	h := httpapi.New(logger)
	httpapi.Register(srv, h)

	srv.MarkReady(true)

	if err := srv.Run(context.Background()); err != nil {
		logger.Error("server exited with error", "err", err)
		os.Exit(1)
	}
}

func handleStatus(w http.ResponseWriter, _ *http.Request) {
	resp := map[string]any{
		"service":     serviceName,
		"status":      "ch-11-rate-mass-surplus-value",
		"description": "Drives the simulated economy forward one period at a time.",
		"tick":        0,
		"running":     false,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
```

- [ ] **Step 2: Commit**

```bash
git add services/simulation-engine/cmd/simulation-engine/main.go
git commit -m "feat(ch11): wire surplus handler into simulation-engine main"
```

---

## Task 5: API Gateway Proxy

**Files:**
- Modify: `services/api-gateway/cmd/api-gateway/main.go`

- [ ] **Step 1: Add simulation-engine proxy and surplus routes**

In `main()`, after the agentProxy block (around line 93), add:

```go
	// Reverse-proxy routes to simulation-engine.
	simURL := getenv("SIM_ENGINE_URL", "http://simulation-engine:8084")
	simProxy, err := proxy.New(simURL, logger)
	if err != nil {
		logger.Error("failed to build simulation-engine proxy", "err", err)
		os.Exit(1)
	}
	srv.Handle("/v1/sim/status", simProxy)

	// Ch. 11 — Rate and Mass of Surplus-Value → simulation-engine
	srv.Handle("/v1/surplus/mass", simProxy)
	srv.Handle("/v1/surplus/limits", simProxy)
```

Also update the `handleInfo` response:
```go
"status": "ch-11-rate-mass-surplus-value",
"chapter": "Capital Vol. I, Ch. 11 - The Rate and Mass of Surplus-Value",
```

- [ ] **Step 2: Commit**

```bash
git add services/api-gateway/cmd/api-gateway/main.go
git commit -m "feat(ch11): proxy /v1/surplus/... to simulation-engine in api-gateway"
```

---

## Task 6: TypeScript Types

**Files:**
- Modify: `web/src/types.ts`

- [ ] **Step 1: Append Ch. 11 types to the bottom of types.ts**

```typescript
// --- simulation-engine types (Ch. 11: Rate and Mass of Surplus-Value) ---------

export interface SurplusValueRate {
  surplus_labour: number;   // LabourMinutes
  necessary_labour: number; // LabourMinutes
}

export interface SurplusValueSnapshot {
  rate: SurplusValueRate;
  variable_capital: number;  // integer unit (same as SurplusValueMass)
  worker_count: number;
  mass: number;              // primary result
  mass_by_rate: number;      // formula I: (s/v) × V
  mass_by_workers: number;   // formula II: v × (s/v) × n
}

export interface ComputeMassInput {
  surplus_labour: number;
  necessary_labour: number;
  variable_capital?: number;
  labour_power_value?: number;
  worker_count?: number;
}

export interface SurplusLimitsResponse {
  absolute_workday_limit: number; // 1440
  minimum_capital?: number;
  labour_power_value?: number;
  worker_count?: number;
}
```

- [ ] **Step 2: Commit**

```bash
git add web/src/types.ts
git commit -m "feat(ch11): add SurplusValueSnapshot and related TypeScript types"
```

---

## Task 7: API Client Methods

**Files:**
- Modify: `web/src/api.ts`

- [ ] **Step 1: Add imports for new types at the top of api.ts**

In the import block, add:
```typescript
  ComputeMassInput,
  SurplusLimitsResponse,
  SurplusValueSnapshot,
```

- [ ] **Step 2: Append methods to the `api` object (before the closing `}`)**

```typescript
  // --- simulation-engine (Ch. 11: Rate and Mass of Surplus-Value) ---

  computeSurplusMass: (input: ComputeMassInput) =>
    http<SurplusValueSnapshot>("/v1/surplus/mass", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getSurplusLimits: (labourPowerValue?: number, workerCount?: number) => {
    const params = new URLSearchParams();
    if (labourPowerValue !== undefined) params.set("labour_power_value", String(labourPowerValue));
    if (workerCount !== undefined) params.set("worker_count", String(workerCount));
    const qs = params.toString();
    return http<SurplusLimitsResponse>(`/v1/surplus/limits${qs ? `?${qs}` : ""}`);
  },
```

- [ ] **Step 3: Commit**

```bash
git add web/src/api.ts
git commit -m "feat(ch11): add computeSurplusMass and getSurplusLimits to API client"
```

---

## Task 8: React Panel

**Files:**
- Create: `web/src/chapters/Ch11RateAndMassOfSurplusValue.tsx`

- [ ] **Step 1: Create the panel file**

```tsx
import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type { SurplusValueSnapshot, SurplusLimitsResponse } from "../types";

interface Ch11Props {
  onSharedChanged: () => void;
}

export function Ch11RateAndMassOfSurplusValue({ onSharedChanged: _unused }: Ch11Props) {
  return (
    <>
      <MassCalculatorPanel />
      <LimitsPanel />
    </>
  );
}

type Fixture = {
  label: string;
  surplusLabour: number;
  necessaryLabour: number;
  variableCapital?: number;
  labourPowerValue?: number;
  workerCount?: number;
};

const FIXTURES: Fixture[] = [
  {
    label: "§1 Rate 100%, V=3s (1 worker)",
    surplusLabour: 6,
    necessaryLabour: 6,
    variableCapital: 3,
  },
  {
    label: "§1 100 workers, v=3s, rate 100%",
    surplusLabour: 6,
    necessaryLabour: 6,
    labourPowerValue: 3,
    workerCount: 100,
  },
  {
    label: "§1 Compensation: rate 200%, V=150s",
    surplusLabour: 12,
    necessaryLabour: 6,
    variableCapital: 150,
  },
  {
    label: "§1 V=1500s, rate 100%",
    surplusLabour: 6,
    necessaryLabour: 6,
    variableCapital: 1500,
  },
  {
    label: "§1 V=300s, rate 200%",
    surplusLabour: 12,
    necessaryLabour: 6,
    variableCapital: 300,
  },
];

function MassCalculatorPanel() {
  const [sl, setSl] = useState(6);
  const [nl, setNl] = useState(6);
  const [vc, setVc] = useState<number | "">("");
  const [lpv, setLpv] = useState<number | "">("");
  const [wc, setWc] = useState<number | "">("");
  const [result, setResult] = useState<SurplusValueSnapshot | null>(null);
  const [err, setErr] = useState<string | null>(null);

  function loadFixture(f: Fixture) {
    setSl(f.surplusLabour);
    setNl(f.necessaryLabour);
    setVc(f.variableCapital ?? "");
    setLpv(f.labourPowerValue ?? "");
    setWc(f.workerCount ?? "");
    setResult(null);
    setErr(null);
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const input = {
        surplus_labour: sl,
        necessary_labour: nl,
        ...(vc !== "" ? { variable_capital: Number(vc) } : {}),
        ...(lpv !== "" ? { labour_power_value: Number(lpv) } : {}),
        ...(wc !== "" ? { worker_count: Number(wc) } : {}),
      };
      const r = await api.computeSurplusMass(input);
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  const rate = nl > 0 ? ((sl / nl) * 100).toFixed(1) : "—";

  return (
    <section className="card">
      <h2>Mass of Surplus-Value</h2>
      <p className="description">
        Computes S via two formulas: S = (s/v) × V [rate formula] and S = v × (s/v) × n
        [worker formula]. When V = v × n both must agree — the compensation law shows that
        raising the rate can offset a fall in the number of workers.
      </p>
      <div style={{ marginBottom: "0.75rem", display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
        {FIXTURES.map((f) => (
          <button key={f.label} type="button" onClick={() => loadFixture(f)}>
            {f.label}
          </button>
        ))}
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Surplus Labour a′ (min)</span>
          <input type="number" min={0} value={sl} onChange={(e) => setSl(Number(e.target.value))} />
        </label>
        <label>
          <span>Necessary Labour a (min)</span>
          <input type="number" min={1} value={nl} onChange={(e) => setNl(Number(e.target.value))} />
        </label>
        <label>
          <span>Rate s/v</span>
          <input type="text" readOnly value={`${rate}%`} />
        </label>
        <label>
          <span>Variable Capital V (optional)</span>
          <input
            type="number"
            min={0}
            value={vc}
            placeholder="e.g. 300"
            onChange={(e) => setVc(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <label>
          <span>Labour-Power Value v per worker (optional)</span>
          <input
            type="number"
            min={0}
            value={lpv}
            placeholder="e.g. 3"
            onChange={(e) => setLpv(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <label>
          <span>Worker Count n (optional)</span>
          <input
            type="number"
            min={1}
            value={wc}
            placeholder="e.g. 100"
            onChange={(e) => setWc(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Compute</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <table className="data-table">
          <tbody>
            <tr>
              <td>Rate (a′/a)</td>
              <td><strong>{(result.rate.surplus_labour / result.rate.necessary_labour * 100).toFixed(1)}%</strong></td>
            </tr>
            <tr><td>Variable Capital V</td><td>{result.variable_capital}</td></tr>
            <tr><td>Worker Count n</td><td>{result.worker_count || "—"}</td></tr>
            <tr>
              <td>Mass by Rate formula  S = (s/v)×V</td>
              <td>{result.mass_by_rate}</td>
            </tr>
            <tr>
              <td>Mass by Workers formula  S = v×(s/v)×n</td>
              <td>{result.mass_by_workers || "—"}</td>
            </tr>
            <tr>
              <td><strong>Total Surplus-Value S</strong></td>
              <td><strong>{result.mass}</strong></td>
            </tr>
            {result.mass_by_rate > 0 && result.mass_by_workers > 0 && (
              <tr>
                <td>Formulas agree?</td>
                <td>
                  {result.mass_by_rate === result.mass_by_workers
                    ? "Yes ✓"
                    : <span className="error">No — inputs inconsistent (V ≠ v×n)</span>}
                </td>
              </tr>
            )}
          </tbody>
        </table>
      )}
    </section>
  );
}

function LimitsPanel() {
  const [lpv, setLpv] = useState<number | "">("");
  const [wc, setWc] = useState<number | "">("");
  const [result, setResult] = useState<SurplusLimitsResponse | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.getSurplusLimits(
        lpv !== "" ? Number(lpv) : undefined,
        wc !== "" ? Number(wc) : undefined,
      );
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Surplus-Value Limits</h2>
      <p className="description">
        Returns the absolute physical limit of the working day (24 h = 1440 min) and,
        optionally, the minimum variable capital required to employ n workers at daily
        reproduction cost v.
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Labour-Power Value v (optional)</span>
          <input
            type="number"
            min={1}
            value={lpv}
            placeholder="e.g. 3"
            onChange={(e) => setLpv(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <label>
          <span>Worker Count n (optional)</span>
          <input
            type="number"
            min={1}
            value={wc}
            placeholder="e.g. 2"
            onChange={(e) => setWc(e.target.value === "" ? "" : Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Get Limits</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <table className="data-table">
          <tbody>
            <tr>
              <td>Absolute Workday Limit</td>
              <td>{result.absolute_workday_limit} min ({result.absolute_workday_limit / 60} h)</td>
            </tr>
            {result.minimum_capital !== undefined && (
              <>
                <tr>
                  <td>Labour-Power Value v</td>
                  <td>{result.labour_power_value}</td>
                </tr>
                <tr>
                  <td>Worker Count n</td>
                  <td>{result.worker_count}</td>
                </tr>
                <tr>
                  <td><strong>Minimum Capital V_min = v × n</strong></td>
                  <td><strong>{result.minimum_capital}</strong></td>
                </tr>
              </>
            )}
          </tbody>
        </table>
      )}
    </section>
  );
}
```

- [ ] **Step 2: Commit**

```bash
git add web/src/chapters/Ch11RateAndMassOfSurplusValue.tsx
git commit -m "feat(ch11): React panel for mass and limits of surplus-value"
```

---

## Task 9: Wire the Panel (ChapterShell + registry)

**Files:**
- Modify: `web/src/components/ChapterShell.tsx`
- Modify: `web/src/chapters/registry.ts`

- [ ] **Step 1: Add import to ChapterShell.tsx**

After the `Ch10WorkingDay` import line, add:
```tsx
import { Ch11RateAndMassOfSurplusValue } from "../chapters/Ch11RateAndMassOfSurplusValue";
```

- [ ] **Step 2: Add the quote for ch11**

In the `QUOTES` object, add:
```tsx
  ch11: "The mass of the surplus-value produced is equal to the amount of the variable capital advanced, multiplied by the rate of surplus-value.",
```

- [ ] **Step 3: Add the ch11 case to the render switch**

After the `ch10` case and before the `null` fallback:
```tsx
        ) : activeChapterId === "ch11" ? (
          <Ch11RateAndMassOfSurplusValue onSharedChanged={onSharedChanged} />
```

- [ ] **Step 4: Mark ch11 as done in registry.ts**

Change:
```ts
  { id: "ch11", number: 11, title: "The Rate and Mass of Surplus-Value", part: "Part III — The Production of Absolute Surplus-Value", status: "pending" },
```
to:
```ts
  { id: "ch11", number: 11, title: "The Rate and Mass of Surplus-Value", part: "Part III — The Production of Absolute Surplus-Value", status: "done" },
```

- [ ] **Step 5: Commit**

```bash
git add web/src/components/ChapterShell.tsx web/src/chapters/registry.ts
git commit -m "feat(ch11): wire Ch11 panel into ChapterShell and mark registry done"
```

---

## Task 10: Update docs/architecture.md

**Files:**
- Modify: `docs/architecture.md`

- [ ] **Step 1: Update the roadmap table**

Change:
```
| Ch. 11+   | Pending     | Cooperation, machinery, wages, accumulation              | all                           |
```
to:
```
| Ch. 11    | ✅ Done     | Rate and mass of surplus-value; SurplusValueRate, MassByRate, MassByWorkers, compensation law | simulation-engine |
| Ch. 12+   | Pending     | Relative surplus-value, cooperation, machinery, wages, accumulation | all            |
```

- [ ] **Step 2: Ask the user to run verification**

```bash
make vet test build
cd web && npm run lint && npm run build
```

- [ ] **Step 3: Commit**

```bash
git add docs/architecture.md
git commit -m "docs: mark Ch.11 done in architecture roadmap"
```

---

## Self-Review

### Spec Coverage

| Spec Item | Task |
|-----------|------|
| `SurplusValueMass` type | Task 1 |
| `SurplusValueRate` struct + `Rate()` method | Task 1 |
| `VariableCapital` type | Task 1 |
| `LabourPowerValue` type | Task 1 |
| `WorkerCount` type | Task 1 |
| `IndividualSurplus` func | Task 1 |
| `MassByRate` func | Task 1 |
| `MassByWorkers` func | Task 1 |
| `MinimumCapital` func | Task 1 |
| `AbsoluteWorkdayLimit` const | Task 1 |
| `SurplusValueSnapshot` struct | Task 1 |
| §1 fixtures (5 + minimum capital) | Task 2 |
| Invariants (formula agreement, individual×n, absolute limit, working day bound) | Task 2 |
| `POST /v1/surplus/mass` | Tasks 3–4 |
| `GET /v1/surplus/limits` | Tasks 3–4 |
| API gateway routing | Task 5 |
| React panel (rate inputs, mass display, compensation law fixtures) | Task 8 |
| Registry + ChapterShell | Task 9 |
| docs/architecture.md | Task 10 |

### Type Consistency Check

- `SurplusValueRate.SurplusLabour` / `.NecessaryLabour` are `int64` throughout Go and `number` in TypeScript.
- `SurplusValueMass` is `int64` in Go, `number` in TypeScript.
- Handler uses `surplus.VariableCapital`, `surplus.WorkerCount` etc consistently.
- React panel `result.mass_by_rate === result.mass_by_workers` comparison is valid for integer values.

### No Placeholders

All steps include complete code. No "TBD" or "implement later" present.
