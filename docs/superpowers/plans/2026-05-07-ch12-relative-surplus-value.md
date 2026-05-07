# Ch. 12 — Relative Surplus-Value — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the Ch. 12 domain types, HTTP endpoints, and React panel for relative surplus-value, wiring them into the simulation-engine service and api-gateway.

**Architecture:** A new `production` package under `simulation-engine/internal/production` holds all pure domain types and functions (no persistence needed). Four stateless HTTP endpoints are added to a new `production_handler.go` file in the existing `httpapi` package. The api-gateway gets proxy routes for `/v1/production/*`. The React panel adds three interaction sections: working-day recording, shortening (relative surplus-value), and extra-surplus-value probe.

**Tech Stack:** Go 1.25 (standard library only for domain), `net/http` + `httptest` for handler tests, React 18 + TypeScript for the frontend panel.

---

## File Map

### Created
- `services/simulation-engine/internal/production/production.go` — domain types and pure functions
- `services/simulation-engine/internal/production/production_test.go` — domain unit tests
- `services/simulation-engine/internal/transport/httpapi/production_handler.go` — HTTP handler methods
- `services/simulation-engine/internal/transport/httpapi/production_handler_test.go` — HTTP handler tests
- `web/src/chapters/Ch12RelativeSurplusValue.tsx` — React panel

### Modified
- `services/simulation-engine/internal/transport/httpapi/routes.go` — register production routes
- `services/simulation-engine/cmd/simulation-engine/main.go` — update status comment
- `services/api-gateway/cmd/api-gateway/main.go` — add `/v1/production/*` proxy
- `web/src/types.ts` — add Ch. 12 TypeScript types
- `web/src/api.ts` — add Ch. 12 API client methods
- `web/src/components/ChapterShell.tsx` — wire Ch12 panel
- `web/src/chapters/registry.ts` — mark `ch12` as `"done"`
- `docs/architecture.md` — update roadmap row

---

## Task 1: Write failing tests for the production domain

**Files:**
- Create: `services/simulation-engine/internal/production/production_test.go`

- [ ] **Step 1: Create the test file**

```go
package production_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/production"
)

// §1 "the working day of 12 hours; the portion a b 10 hours necessary, b c 2 hours surplus"
func TestWorkingDayInvariant(t *testing.T) {
	t.Parallel()
	wd := production.WorkingDay{
		Total:           720,
		NecessaryLabour: 600,
		SurplusLabour:   120,
	}
	if production.LabourMinutes(wd.NecessaryLabour)+production.LabourMinutes(wd.SurplusLabour) != wd.Total {
		t.Fatalf("NecessaryLabour + SurplusLabour must equal Total: %d + %d ≠ %d",
			wd.NecessaryLabour, wd.SurplusLabour, wd.Total)
	}
}

// §1 "surplus-labour increases by one half, from 2 hours to 3 hours, although the working day remains at 12 hours"
// ShortenNecessaryLabour(wd, 540) → WorkingDay{Total:720, NecessaryLabour:540, SurplusLabour:180}
func TestShortenNecessaryLabour_Fixture(t *testing.T) {
	t.Parallel()
	wd := production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	got := production.ShortenNecessaryLabour(wd, production.LabourPowerValue(540))
	if got.Total != 720 {
		t.Fatalf("expected Total=720, got %d", got.Total)
	}
	if got.NecessaryLabour != 540 {
		t.Fatalf("expected NecessaryLabour=540, got %d", got.NecessaryLabour)
	}
	if got.SurplusLabour != 180 {
		t.Fatalf("expected SurplusLabour=180, got %d", got.SurplusLabour)
	}
}

// §1 invariant: ShortenNecessaryLabour increases SurplusLabour when newLPV < wd.NecessaryLabour
func TestShortenNecessaryLabour_SurplusIncreases(t *testing.T) {
	t.Parallel()
	wd := production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	got := production.ShortenNecessaryLabour(wd, production.LabourPowerValue(360))
	if got.SurplusLabour <= wd.SurplusLabour {
		t.Fatalf("surplus must increase: before=%d, after=%d", wd.SurplusLabour, got.SurplusLabour)
	}
}

// §1 invariant: result of ShortenNecessaryLabour still satisfies NL + SL = Total
func TestShortenNecessaryLabour_Invariant(t *testing.T) {
	t.Parallel()
	wd := production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	got := production.ShortenNecessaryLabour(wd, production.LabourPowerValue(360))
	if production.LabourMinutes(got.NecessaryLabour)+production.LabourMinutes(got.SurplusLabour) != got.Total {
		t.Fatalf("invariant broken after shorten: %d + %d ≠ %d",
			got.NecessaryLabour, got.SurplusLabour, got.Total)
	}
}

// §1 "value of labour-power reduced from five shillings to three, surplus-value increases from one to three"
// After reducing NecessaryLabour from 600 to 360: SurplusLabour = 720 - 360 = 360
func TestShortenNecessaryLabour_LPVReduction(t *testing.T) {
	t.Parallel()
	wd := production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	// new LabourPowerValue = 360 min (representing the reduced subsistence value in LabourMinutes)
	got := production.ShortenNecessaryLabour(wd, production.LabourPowerValue(360))
	if got.NecessaryLabour != 360 {
		t.Fatalf("expected NecessaryLabour=360, got %d", got.NecessaryLabour)
	}
	if got.SurplusLabour != 360 {
		t.Fatalf("expected SurplusLabour=360, got %d", got.SurplusLabour)
	}
}

// §1 RateOfSurplusValue: s/v = 120/600 = 0.2
func TestRateOfSurplusValue(t *testing.T) {
	t.Parallel()
	rate := production.RateOfSurplusValue(120, 600)
	if rate != 0.2 {
		t.Fatalf("expected 0.2, got %f", rate)
	}
}

// §1 invariant: RateOfSurplusValue strictly increases as ProductivityFactor rises (NL shrinks, total fixed)
func TestRateOfSurplusValue_IncreasesWithProductivity(t *testing.T) {
	t.Parallel()
	wd := production.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	rateOld := production.RateOfSurplusValue(wd.SurplusLabour, wd.NecessaryLabour)
	wdNew := production.ShortenNecessaryLabour(wd, production.LabourPowerValue(360))
	rateNew := production.RateOfSurplusValue(wdNew.SurplusLabour, wdNew.NecessaryLabour)
	if rateNew <= rateOld {
		t.Fatalf("rate must increase with productivity rise: old=%f, new=%f", rateOld, rateNew)
	}
}

// §1 RateOfSurplusValue: zero necessary labour returns 0 (guard)
func TestRateOfSurplusValue_ZeroNL(t *testing.T) {
	t.Parallel()
	if got := production.RateOfSurplusValue(120, 0); got != 0 {
		t.Fatalf("expected 0 for zero NL, got %f", got)
	}
}

// §1 "double the productiveness ... individual value below social value"
// ExtraSurplusValue(IndividualValue(30), SocialValue(60), Quantity(24)) = (60-30)*24 = 720
func TestExtraSurplusValue_Fixture(t *testing.T) {
	t.Parallel()
	got := production.ExtraSurplusValue(
		production.IndividualValue(30),
		production.SocialValue(60),
		production.Quantity(24),
	)
	if got != 720 {
		t.Fatalf("expected 720, got %d", got)
	}
}

// §1 per-article extra = SocialValue − IndividualValue
func TestExtraSurplusValue_PerUnit(t *testing.T) {
	t.Parallel()
	got := production.ExtraSurplusValue(
		production.IndividualValue(30),
		production.SocialValue(60),
		production.Quantity(1),
	)
	if got != 30 {
		t.Fatalf("expected per-unit=30, got %d", got)
	}
}

// §1 invariant: ExtraSurplusValue == 0 when iv == sv
func TestExtraSurplusValue_ZeroWhenEqual(t *testing.T) {
	t.Parallel()
	got := production.ExtraSurplusValue(
		production.IndividualValue(60),
		production.SocialValue(60),
		production.Quantity(24),
	)
	if got != 0 {
		t.Fatalf("expected 0 when iv==sv, got %d", got)
	}
}

// §1 invariant: ExtraSurplusValue == 0 when iv > sv (no extra gain)
func TestExtraSurplusValue_ZeroWhenAbove(t *testing.T) {
	t.Parallel()
	got := production.ExtraSurplusValue(
		production.IndividualValue(70),
		production.SocialValue(60),
		production.Quantity(24),
	)
	if got != 0 {
		t.Fatalf("expected 0 when iv>sv, got %d", got)
	}
}

// value ∝ 1/productivity: doubling productivity halves SNLT
func TestApplyProductivityToSNLT_Double(t *testing.T) {
	t.Parallel()
	got := production.ApplyProductivityToSNLT(60, production.ProductivityFactor(2.0))
	if got != 30 {
		t.Fatalf("expected 30, got %d", got)
	}
}

// value ∝ 1/productivity: halving productivity doubles SNLT
func TestApplyProductivityToSNLT_Half(t *testing.T) {
	t.Parallel()
	got := production.ApplyProductivityToSNLT(60, production.ProductivityFactor(0.5))
	if got != 120 {
		t.Fatalf("expected 120, got %d", got)
	}
}

// guard: non-positive factor returns original snlt unchanged
func TestApplyProductivityToSNLT_ZeroFactor(t *testing.T) {
	t.Parallel()
	got := production.ApplyProductivityToSNLT(60, production.ProductivityFactor(0))
	if got != 60 {
		t.Fatalf("expected original 60 for zero factor, got %d", got)
	}
}
```

- [ ] **Step 2: Ask the user to run `make vet test build`**

Expected: compilation fails with "cannot find package ...production"

---

## Task 2: Implement the production domain package

**Files:**
- Create: `services/simulation-engine/internal/production/production.go`

- [ ] **Step 1: Create the domain file**

```go
// Package production implements relative surplus-value concepts from
// Capital Vol. I, Ch. 12. All functions are pure; no persistence is needed.
package production

import "math"

// LabourMinutes is the canonical value-magnitude unit.
type LabourMinutes int64

// Named LabourMinutes types for distinct domain roles (Ch. 12, §1).
type AbsoluteSurplusValue LabourMinutes
type RelativeSurplusValue LabourMinutes
type NecessaryLabour LabourMinutes
type SurplusLabour LabourMinutes
type LabourPowerValue LabourMinutes
type IndividualValue LabourMinutes
type SocialValue LabourMinutes

// ProductivityFactor is a multiplicative scalar on SNLT: value ∝ 1/productivity [§1].
type ProductivityFactor float64

// Quantity is the number of articles produced or sold.
type Quantity int64

// WorkingDay holds the tripartite split of the total working day [§1].
// Invariant: NecessaryLabour + SurplusLabour == Total.
type WorkingDay struct {
	Total           LabourMinutes   `json:"total"`
	NecessaryLabour NecessaryLabour `json:"necessary_labour"`
	SurplusLabour   SurplusLabour   `json:"surplus_labour"`
}

// ShortenNecessaryLabour recomputes the working-day split after a productivity
// rise reduces the value of labour-power to newLPV. Total day is unchanged [§1].
func ShortenNecessaryLabour(wd WorkingDay, newLPV LabourPowerValue) WorkingDay {
	nl := NecessaryLabour(newLPV)
	sl := SurplusLabour(wd.Total - LabourMinutes(nl))
	return WorkingDay{Total: wd.Total, NecessaryLabour: nl, SurplusLabour: sl}
}

// RateOfSurplusValue returns s/v as a float64. Returns 0 when nl == 0 [§1].
func RateOfSurplusValue(sl SurplusLabour, nl NecessaryLabour) float64 {
	if nl == 0 {
		return 0
	}
	return float64(sl) / float64(nl)
}

// ExtraSurplusValue is the total extra gain captured by a capitalist whose
// individual value is below the social value, selling at social value.
// Returns 0 when iv >= sv [§1].
func ExtraSurplusValue(iv IndividualValue, sv SocialValue, qty Quantity) LabourMinutes {
	if iv >= IndividualValue(sv) {
		return 0
	}
	return LabourMinutes(sv-SocialValue(iv)) * LabourMinutes(qty)
}

// ApplyProductivityToSNLT implements the inverse law: value ∝ 1/productivity.
// Returns snlt unchanged for non-positive factor [§1].
func ApplyProductivityToSNLT(snlt LabourMinutes, pf ProductivityFactor) LabourMinutes {
	if pf <= 0 {
		return snlt
	}
	return LabourMinutes(math.Round(float64(snlt) / float64(pf)))
}
```

- [ ] **Step 2: Ask the user to run `make vet test build`**

Expected: `ok  github.com/theding0x/capital-simulator/services/simulation-engine/internal/production`

---

## Task 3: Write failing HTTP handler tests

**Files:**
- Create: `services/simulation-engine/internal/transport/httpapi/production_handler_test.go`

- [ ] **Step 1: Create the test file**

```go
package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newProductionTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	h := New(nil)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/production/working-day", h.RecordWorkingDay)
	mux.HandleFunc("POST /v1/production/working-day/shorten", h.ShortenWorkingDay)
	mux.HandleFunc("GET /v1/production/rate-of-surplus-value", h.GetProductionRate)
	mux.HandleFunc("POST /v1/production/extra-surplus-value", h.ComputeExtraSurplusValue)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// §1 fixture: total=720, lpv=600 → WorkingDay{720, 600, 120}
func TestRecordWorkingDay(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	body := `{"total":720,"labour_power_value":600}`
	res, err := http.Post(ts.URL+"/v1/production/working-day", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	var got struct {
		Total           int64 `json:"total"`
		NecessaryLabour int64 `json:"necessary_labour"`
		SurplusLabour   int64 `json:"surplus_labour"`
	}
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Total != 720 || got.NecessaryLabour != 600 || got.SurplusLabour != 120 {
		t.Fatalf("unexpected response: %+v", got)
	}
}

// validation: labour_power_value >= total must be rejected
func TestRecordWorkingDay_InvalidLPV(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	body := `{"total":720,"labour_power_value":720}`
	res, _ := http.Post(ts.URL+"/v1/production/working-day", "application/json", strings.NewReader(body))
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

// §1 fixture: shorten(wd={720,600,120}, newLPV=540) → {720,540,180}, delta=60
func TestShortenWorkingDay(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	body := `{"total":720,"necessary_labour":600,"surplus_labour":120,"new_labour_power_value":540}`
	res, err := http.Post(ts.URL+"/v1/production/working-day/shorten", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	var got struct {
		WorkingDay struct {
			Total           int64 `json:"total"`
			NecessaryLabour int64 `json:"necessary_labour"`
			SurplusLabour   int64 `json:"surplus_labour"`
		} `json:"working_day"`
		RelativeSurplusValue int64 `json:"relative_surplus_value"`
	}
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.WorkingDay.NecessaryLabour != 540 {
		t.Fatalf("expected NecessaryLabour=540, got %d", got.WorkingDay.NecessaryLabour)
	}
	if got.WorkingDay.SurplusLabour != 180 {
		t.Fatalf("expected SurplusLabour=180, got %d", got.WorkingDay.SurplusLabour)
	}
	if got.RelativeSurplusValue != 60 {
		t.Fatalf("expected RelativeSurplusValue=60, got %d", got.RelativeSurplusValue)
	}
}

// validation: new_labour_power_value >= total must be rejected
func TestShortenWorkingDay_InvalidNewLPV(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	body := `{"total":720,"necessary_labour":600,"surplus_labour":120,"new_labour_power_value":720}`
	res, _ := http.Post(ts.URL+"/v1/production/working-day/shorten", "application/json", strings.NewReader(body))
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

// §1 rate: necessary=600, surplus=120 → rate=0.2
func TestGetProductionRate(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	res, err := http.Get(ts.URL + "/v1/production/rate-of-surplus-value?necessary=600&surplus=120")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	var got struct {
		Rate            float64 `json:"rate"`
		SurplusLabour   int64   `json:"surplus_labour"`
		NecessaryLabour int64   `json:"necessary_labour"`
	}
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Rate != 0.2 {
		t.Fatalf("expected rate=0.2, got %f", got.Rate)
	}
	if got.SurplusLabour != 120 || got.NecessaryLabour != 600 {
		t.Fatalf("unexpected echoed fields: %+v", got)
	}
}

// validation: missing necessary query param → 400
func TestGetProductionRate_MissingParam(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	res, _ := http.Get(ts.URL + "/v1/production/rate-of-surplus-value?surplus=120")
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

// §1 fixture: IndividualValue=30, SocialValue=60, Quantity=24 → extra=720, per_unit=30
func TestComputeExtraSurplusValue(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	body := `{"individual_value":30,"social_value":60,"quantity":24}`
	res, err := http.Post(ts.URL+"/v1/production/extra-surplus-value", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	var got struct {
		ExtraSurplusValue int64 `json:"extra_surplus_value"`
		PerUnit           int64 `json:"per_unit"`
		IndividualValue   int64 `json:"individual_value"`
		SocialValue       int64 `json:"social_value"`
		Quantity          int64 `json:"quantity"`
	}
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ExtraSurplusValue != 720 {
		t.Fatalf("expected extra_surplus_value=720, got %d", got.ExtraSurplusValue)
	}
	if got.PerUnit != 30 {
		t.Fatalf("expected per_unit=30, got %d", got.PerUnit)
	}
}

// §1 invariant: extra=0 when iv==sv
func TestComputeExtraSurplusValue_ZeroWhenEqual(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	body := `{"individual_value":60,"social_value":60,"quantity":24}`
	res, _ := http.Post(ts.URL+"/v1/production/extra-surplus-value", "application/json", strings.NewReader(body))
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	var got struct {
		ExtraSurplusValue int64 `json:"extra_surplus_value"`
	}
	json.NewDecoder(res.Body).Decode(&got) //nolint
	if got.ExtraSurplusValue != 0 {
		t.Fatalf("expected 0 when iv==sv, got %d", got.ExtraSurplusValue)
	}
}
```

- [ ] **Step 2: Ask the user to run `make vet test build`**

Expected: compilation error — `h.RecordWorkingDay` etc. undefined

---

## Task 4: Implement the HTTP production handler

**Files:**
- Create: `services/simulation-engine/internal/transport/httpapi/production_handler.go`

- [ ] **Step 1: Create the handler file**

```go
package httpapi

import (
	"net/http"
	"strconv"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/production"
)

// workingDayDTO is the shared JSON shape for a WorkingDay.
type workingDayDTO struct {
	Total           int64 `json:"total"`
	NecessaryLabour int64 `json:"necessary_labour"`
	SurplusLabour   int64 `json:"surplus_labour"`
}

func dtoFromWorkingDay(wd production.WorkingDay) workingDayDTO {
	return workingDayDTO{
		Total:           int64(wd.Total),
		NecessaryLabour: int64(wd.NecessaryLabour),
		SurplusLabour:   int64(wd.SurplusLabour),
	}
}

// recordWorkingDayRequest is the POST /v1/production/working-day body.
type recordWorkingDayRequest struct {
	Total            int64 `json:"total"`
	LabourPowerValue int64 `json:"labour_power_value"`
}

// RecordWorkingDay handles POST /v1/production/working-day.
// Builds a WorkingDay from total and labour-power value (both in LabourMinutes).
func (h *Handler) RecordWorkingDay(w http.ResponseWriter, r *http.Request) {
	var req recordWorkingDayRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Total <= 0 {
		writeError(w, http.StatusBadRequest, "total must be positive")
		return
	}
	if req.LabourPowerValue <= 0 {
		writeError(w, http.StatusBadRequest, "labour_power_value must be positive")
		return
	}
	if req.LabourPowerValue >= req.Total {
		writeError(w, http.StatusBadRequest, "labour_power_value must be less than total")
		return
	}
	wd := production.WorkingDay{
		Total:           production.LabourMinutes(req.Total),
		NecessaryLabour: production.NecessaryLabour(req.LabourPowerValue),
		SurplusLabour:   production.SurplusLabour(req.Total - req.LabourPowerValue),
	}
	writeJSON(w, http.StatusOK, dtoFromWorkingDay(wd))
}

// shortenWorkingDayRequest is the POST /v1/production/working-day/shorten body.
type shortenWorkingDayRequest struct {
	Total               int64 `json:"total"`
	NecessaryLabour     int64 `json:"necessary_labour"`
	SurplusLabour       int64 `json:"surplus_labour"`
	NewLabourPowerValue int64 `json:"new_labour_power_value"`
}

// shortenWorkingDayResponse is the POST /v1/production/working-day/shorten response.
type shortenWorkingDayResponse struct {
	WorkingDay           workingDayDTO `json:"working_day"`
	RelativeSurplusValue int64         `json:"relative_surplus_value"`
}

// ShortenWorkingDay handles POST /v1/production/working-day/shorten.
// Returns the updated WorkingDay and the RelativeSurplusValue delta (new SL − old SL).
func (h *Handler) ShortenWorkingDay(w http.ResponseWriter, r *http.Request) {
	var req shortenWorkingDayRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Total <= 0 {
		writeError(w, http.StatusBadRequest, "total must be positive")
		return
	}
	if req.NecessaryLabour <= 0 || req.SurplusLabour < 0 {
		writeError(w, http.StatusBadRequest, "necessary_labour must be positive; surplus_labour cannot be negative")
		return
	}
	if req.NecessaryLabour+req.SurplusLabour != req.Total {
		writeError(w, http.StatusBadRequest, "necessary_labour + surplus_labour must equal total")
		return
	}
	if req.NewLabourPowerValue <= 0 {
		writeError(w, http.StatusBadRequest, "new_labour_power_value must be positive")
		return
	}
	if req.NewLabourPowerValue >= req.Total {
		writeError(w, http.StatusBadRequest, "new_labour_power_value must be less than total")
		return
	}
	wd := production.WorkingDay{
		Total:           production.LabourMinutes(req.Total),
		NecessaryLabour: production.NecessaryLabour(req.NecessaryLabour),
		SurplusLabour:   production.SurplusLabour(req.SurplusLabour),
	}
	newWD := production.ShortenNecessaryLabour(wd, production.LabourPowerValue(req.NewLabourPowerValue))
	delta := int64(newWD.SurplusLabour) - int64(wd.SurplusLabour)
	writeJSON(w, http.StatusOK, shortenWorkingDayResponse{
		WorkingDay:           dtoFromWorkingDay(newWD),
		RelativeSurplusValue: delta,
	})
}

// rateResponse is the GET /v1/production/rate-of-surplus-value response.
type rateResponse struct {
	Rate            float64 `json:"rate"`
	SurplusLabour   int64   `json:"surplus_labour"`
	NecessaryLabour int64   `json:"necessary_labour"`
}

// GetProductionRate handles GET /v1/production/rate-of-surplus-value?necessary=&surplus=.
func (h *Handler) GetProductionRate(w http.ResponseWriter, r *http.Request) {
	necessaryStr := r.URL.Query().Get("necessary")
	surplusStr := r.URL.Query().Get("surplus")
	if necessaryStr == "" || surplusStr == "" {
		writeError(w, http.StatusBadRequest, "necessary and surplus query params required")
		return
	}
	necessary, err1 := strconv.ParseInt(necessaryStr, 10, 64)
	surplus, err2 := strconv.ParseInt(surplusStr, 10, 64)
	if err1 != nil || err2 != nil || necessary <= 0 || surplus < 0 {
		writeError(w, http.StatusBadRequest, "necessary must be a positive integer; surplus must be a non-negative integer")
		return
	}
	rate := production.RateOfSurplusValue(
		production.SurplusLabour(surplus),
		production.NecessaryLabour(necessary),
	)
	writeJSON(w, http.StatusOK, rateResponse{
		Rate:            rate,
		SurplusLabour:   surplus,
		NecessaryLabour: necessary,
	})
}

// extraSurplusRequest is the POST /v1/production/extra-surplus-value body.
type extraSurplusRequest struct {
	IndividualValue int64 `json:"individual_value"`
	SocialValue     int64 `json:"social_value"`
	Quantity        int64 `json:"quantity"`
}

// extraSurplusResponse is the POST /v1/production/extra-surplus-value response.
type extraSurplusResponse struct {
	ExtraSurplusValue int64 `json:"extra_surplus_value"`
	PerUnit           int64 `json:"per_unit"`
	IndividualValue   int64 `json:"individual_value"`
	SocialValue       int64 `json:"social_value"`
	Quantity          int64 `json:"quantity"`
}

// ComputeExtraSurplusValue handles POST /v1/production/extra-surplus-value.
func (h *Handler) ComputeExtraSurplusValue(w http.ResponseWriter, r *http.Request) {
	var req extraSurplusRequest
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.IndividualValue <= 0 {
		writeError(w, http.StatusBadRequest, "individual_value must be positive")
		return
	}
	if req.SocialValue <= 0 {
		writeError(w, http.StatusBadRequest, "social_value must be positive")
		return
	}
	if req.Quantity <= 0 {
		writeError(w, http.StatusBadRequest, "quantity must be positive")
		return
	}
	total := production.ExtraSurplusValue(
		production.IndividualValue(req.IndividualValue),
		production.SocialValue(req.SocialValue),
		production.Quantity(req.Quantity),
	)
	perUnit := int64(0)
	if req.SocialValue > req.IndividualValue {
		perUnit = req.SocialValue - req.IndividualValue
	}
	writeJSON(w, http.StatusOK, extraSurplusResponse{
		ExtraSurplusValue: int64(total),
		PerUnit:           perUnit,
		IndividualValue:   req.IndividualValue,
		SocialValue:       req.SocialValue,
		Quantity:          req.Quantity,
	})
}
```

- [ ] **Step 2: Ask the user to run `make vet test build`**

Expected: all tests pass, including the new production handler tests.

---

## Task 5: Register production routes in httpapi/routes.go

**Files:**
- Modify: `services/simulation-engine/internal/transport/httpapi/routes.go`

- [ ] **Step 1: Replace the file contents**

```go
package httpapi

import "github.com/theding0x/capital-simulator/pkg/httpx"

func Register(s *httpx.Server, h *Handler) {
	// Ch. 11 — Rate and Mass of Surplus-Value
	s.HandleFunc("POST /v1/surplus/mass", h.ComputeMass)
	s.HandleFunc("GET /v1/surplus/limits", h.GetLimits)

	// Ch. 12 — Relative Surplus-Value
	s.HandleFunc("POST /v1/production/working-day", h.RecordWorkingDay)
	s.HandleFunc("POST /v1/production/working-day/shorten", h.ShortenWorkingDay)
	s.HandleFunc("GET /v1/production/rate-of-surplus-value", h.GetProductionRate)
	s.HandleFunc("POST /v1/production/extra-surplus-value", h.ComputeExtraSurplusValue)
}
```

---

## Task 6: Update simulation-engine main.go status comment

**Files:**
- Modify: `services/simulation-engine/cmd/simulation-engine/main.go`

- [ ] **Step 1: Update the status string in main.go**

In `handleStatus`, change:
```go
"status": "ch-11-rate-mass-surplus-value",
```
to:
```go
"status": "ch-12-relative-surplus-value",
```

And the chapter description comment:
```go
// Ch. 11 — Rate and Mass of Surplus-Value
h := httpapi.New(logger)
httpapi.Register(srv, h)
```
to:
```go
// Ch. 11–12 — Surplus-Value (rate/mass + relative)
h := httpapi.New(logger)
httpapi.Register(srv, h)
```

---

## Task 7: Add api-gateway proxy routes for /v1/production/*

**Files:**
- Modify: `services/api-gateway/cmd/api-gateway/main.go`

- [ ] **Step 1: Add proxy routes after the existing surplus routes**

After the block:
```go
// Ch. 11 — Rate and Mass of Surplus-Value → simulation-engine
srv.Handle("/v1/surplus/mass", simProxy)
srv.Handle("/v1/surplus/limits", simProxy)
```

Add:
```go
// Ch. 12 — Relative Surplus-Value → simulation-engine
srv.Handle("/v1/production/working-day", simProxy)
srv.Handle("/v1/production/working-day/shorten", simProxy)
srv.Handle("/v1/production/rate-of-surplus-value", simProxy)
srv.Handle("/v1/production/extra-surplus-value", simProxy)
```

Also update the `handleInfo` status string:
```go
"status": "ch-12-relative-surplus-value",
```

And the chapter description:
```go
"chapter": "Capital Vol. I, Ch. 12 - The Concept of Relative Surplus-Value",
```

- [ ] **Step 2: Ask the user to run `make vet test build`**

Expected: clean build for all services.

---

## Task 8: Add TypeScript types for Ch. 12

**Files:**
- Modify: `web/src/types.ts`

- [ ] **Step 1: Append Ch. 12 section at the end of types.ts**

```typescript
// --- simulation-engine types (Ch. 12: Relative Surplus-Value) ----------------

export interface ProductionWorkingDay {
  total: number;             // LabourMinutes
  necessary_labour: number;  // LabourMinutes
  surplus_labour: number;    // LabourMinutes
}

export interface ShortenWorkingDayResponse {
  working_day: ProductionWorkingDay;
  relative_surplus_value: number; // LabourMinutes delta (new SL − old SL)
}

export interface ProductionRateResult {
  rate: number;              // s/v as float
  surplus_labour: number;    // LabourMinutes (echoed)
  necessary_labour: number;  // LabourMinutes (echoed)
}

export interface ExtraSurplusValueResult {
  extra_surplus_value: number; // LabourMinutes total
  per_unit: number;            // LabourMinutes per article
  individual_value: number;
  social_value: number;
  quantity: number;
}

export interface RecordWorkingDayInput {
  total: number;
  labour_power_value: number;
}

export interface ShortenWorkingDayInput {
  total: number;
  necessary_labour: number;
  surplus_labour: number;
  new_labour_power_value: number;
}

export interface ExtraSurplusValueInput {
  individual_value: number;
  social_value: number;
  quantity: number;
}
```

---

## Task 9: Add API client methods for Ch. 12

**Files:**
- Modify: `web/src/api.ts`

- [ ] **Step 1: Add imports for Ch. 12 types**

In the import block at the top of `api.ts`, add the new types:

```typescript
import type {
  // ... existing imports ...
  SurplusLimitsResponse,
  SurplusValueSnapshot,
  // add these:
  ProductionWorkingDay,
  ShortenWorkingDayResponse,
  ProductionRateResult,
  ExtraSurplusValueResult,
  RecordWorkingDayInput,
  ShortenWorkingDayInput,
  ExtraSurplusValueInput,
} from "./types";
```

- [ ] **Step 2: Add API methods after the Ch. 11 section in api.ts**

After the `getSurplusLimits` method, append:

```typescript
  // --- simulation-engine (Ch. 12: Relative Surplus-Value) ---

  recordWorkingDay: (input: RecordWorkingDayInput) =>
    http<ProductionWorkingDay>("/v1/production/working-day", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  shortenWorkingDay: (input: ShortenWorkingDayInput) =>
    http<ShortenWorkingDayResponse>("/v1/production/working-day/shorten", {
      method: "POST",
      body: JSON.stringify(input),
    }),

  getProductionRate: (necessary: number, surplus: number) =>
    http<ProductionRateResult>(
      `/v1/production/rate-of-surplus-value?necessary=${necessary}&surplus=${surplus}`
    ),

  computeExtraSurplusValue: (input: ExtraSurplusValueInput) =>
    http<ExtraSurplusValueResult>("/v1/production/extra-surplus-value", {
      method: "POST",
      body: JSON.stringify(input),
    }),
```

---

## Task 10: Create the React panel

**Files:**
- Create: `web/src/chapters/Ch12RelativeSurplusValue.tsx`

- [ ] **Step 1: Create the panel file**

```tsx
import { useState } from "react";
import type { FormEvent } from "react";
import { api } from "../api";
import type {
  ProductionWorkingDay,
  ShortenWorkingDayResponse,
  ProductionRateResult,
  ExtraSurplusValueResult,
} from "../types";

interface Ch12Props {
  onSharedChanged: () => void;
}

export function Ch12RelativeSurplusValue({ onSharedChanged: _unused }: Ch12Props) {
  return (
    <>
      <WorkingDayPanel />
      <ShortenWorkingDayPanel />
      <ExtraSurplusValuePanel />
    </>
  );
}

// ── Working Day Panel ──────────────────────────────────────────────────────────

function WorkingDayPanel() {
  const [total, setTotal] = useState(720);
  const [lpv, setLpv] = useState(600);
  const [result, setResult] = useState<ProductionWorkingDay | null>(null);
  const [rateResult, setRateResult] = useState<ProductionRateResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const wd = await api.recordWorkingDay({ total, labour_power_value: lpv });
      setResult(wd);
      const rate = await api.getProductionRate(wd.necessary_labour, wd.surplus_labour);
      setRateResult(rate);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  const necessaryPct = result ? ((result.necessary_labour / result.total) * 100).toFixed(1) : null;
  const surplusPct = result ? ((result.surplus_labour / result.total) * 100).toFixed(1) : null;

  return (
    <section className="card">
      <h2>Working-Day Split</h2>
      <p className="description">
        Records the split of the working day into necessary labour (reproducing labour-power value)
        and surplus labour (producing surplus-value for capital). §1: "The working day is thus not
        a constant, but a variable quantity."
      </p>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Total Working Day (min)</span>
          <input
            type="number"
            min={1}
            value={total}
            onChange={(e) => setTotal(Number(e.target.value))}
          />
        </label>
        <label>
          <span>Labour-Power Value (min)</span>
          <input
            type="number"
            min={1}
            value={lpv}
            onChange={(e) => setLpv(Number(e.target.value))}
          />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Record</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <>
          <div style={{ marginTop: "1rem" }}>
            <div style={{ display: "flex", height: "1.5rem", borderRadius: "4px", overflow: "hidden", fontSize: "0.75rem" }}>
              <div
                style={{
                  width: `${necessaryPct}%`,
                  background: "var(--color-accent)",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  color: "#fff",
                  fontWeight: 600,
                }}
              >
                Necessary {necessaryPct}%
              </div>
              <div
                style={{
                  width: `${surplusPct}%`,
                  background: "var(--color-muted)",
                  display: "flex",
                  alignItems: "center",
                  justifyContent: "center",
                  color: "#fff",
                  fontWeight: 600,
                }}
              >
                Surplus {surplusPct}%
              </div>
            </div>
          </div>
          <table className="data-table" style={{ marginTop: "0.75rem" }}>
            <tbody>
              <tr><td>Total</td><td>{result.total} min ({(result.total / 60).toFixed(1)} h)</td></tr>
              <tr><td>Necessary Labour</td><td>{result.necessary_labour} min ({(result.necessary_labour / 60).toFixed(1)} h)</td></tr>
              <tr><td>Surplus Labour</td><td>{result.surplus_labour} min ({(result.surplus_labour / 60).toFixed(1)} h)</td></tr>
              {rateResult && (
                <tr>
                  <td><strong>Rate of Surplus-Value s/v</strong></td>
                  <td><strong>{(rateResult.rate * 100).toFixed(1)}%</strong></td>
                </tr>
              )}
            </tbody>
          </table>
        </>
      )}
    </section>
  );
}

// ── Shorten Necessary Labour Panel ────────────────────────────────────────────

const SHORTEN_FIXTURES = [
  { label: "§1 Surplus ½→⅓ (9h→8h NL)", total: 720, nl: 540, sl: 180, newLpv: 480 },
  { label: "§1 LPV 5s→3s (600→360 min NL)", total: 720, nl: 600, sl: 120, newLpv: 360 },
];

function ShortenWorkingDayPanel() {
  const [total, setTotal] = useState(720);
  const [nl, setNl] = useState(600);
  const [sl, setSl] = useState(120);
  const [newLpv, setNewLpv] = useState(540);
  const [result, setResult] = useState<ShortenWorkingDayResponse | null>(null);
  const [err, setErr] = useState<string | null>(null);

  function loadFixture(f: typeof SHORTEN_FIXTURES[number]) {
    setTotal(f.total); setNl(f.nl); setSl(f.sl); setNewLpv(f.newLpv);
    setResult(null); setErr(null);
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.shortenWorkingDay({
        total,
        necessary_labour: nl,
        surplus_labour: sl,
        new_labour_power_value: newLpv,
      });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Shorten Necessary Labour (Relative Surplus-Value)</h2>
      <p className="description">
        A productivity rise reduces the value of labour-power, contracting necessary labour and
        expanding surplus labour within the same working day. The gain in surplus labour is the
        relative surplus-value. §1: "I call that surplus-value which is produced by curtailment
        of the necessary labour-time … relative surplus-value."
      </p>
      <div style={{ marginBottom: "0.75rem", display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
        {SHORTEN_FIXTURES.map((f) => (
          <button key={f.label} type="button" onClick={() => loadFixture(f)}>{f.label}</button>
        ))}
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Total (min)</span>
          <input type="number" min={1} value={total} onChange={(e) => setTotal(Number(e.target.value))} />
        </label>
        <label>
          <span>Current Necessary Labour (min)</span>
          <input type="number" min={1} value={nl} onChange={(e) => setNl(Number(e.target.value))} />
        </label>
        <label>
          <span>Current Surplus Labour (min)</span>
          <input type="number" min={0} value={sl} onChange={(e) => setSl(Number(e.target.value))} />
        </label>
        <label>
          <span>New Labour-Power Value (min)</span>
          <input type="number" min={1} value={newLpv} onChange={(e) => setNewLpv(Number(e.target.value))} />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Shorten</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <table className="data-table">
          <tbody>
            <tr><td>New Total</td><td>{result.working_day.total} min</td></tr>
            <tr><td>New Necessary Labour</td><td>{result.working_day.necessary_labour} min</td></tr>
            <tr><td>New Surplus Labour</td><td>{result.working_day.surplus_labour} min</td></tr>
            <tr>
              <td><strong>Relative Surplus-Value (delta)</strong></td>
              <td><strong>+{result.relative_surplus_value} min</strong></td>
            </tr>
            <tr>
              <td>New Rate s/v</td>
              <td>
                {result.working_day.necessary_labour > 0
                  ? ((result.working_day.surplus_labour / result.working_day.necessary_labour) * 100).toFixed(1) + "%"
                  : "—"}
              </td>
            </tr>
          </tbody>
        </table>
      )}
    </section>
  );
}

// ── Extra Surplus-Value Panel ──────────────────────────────────────────────────

const EXTRA_FIXTURES = [
  { label: "§1 Double productivity (iv=30, sv=60, qty=24)", iv: 30, sv: 60, qty: 24 },
  { label: "§1 Per-article gain (qty=1)", iv: 30, sv: 60, qty: 1 },
];

function ExtraSurplusValuePanel() {
  const [iv, setIv] = useState(30);
  const [sv, setSv] = useState(60);
  const [qty, setQty] = useState(24);
  const [result, setResult] = useState<ExtraSurplusValueResult | null>(null);
  const [err, setErr] = useState<string | null>(null);

  function loadFixture(f: typeof EXTRA_FIXTURES[number]) {
    setIv(f.iv); setSv(f.sv); setQty(f.qty);
    setResult(null); setErr(null);
  }

  async function submit(e: FormEvent) {
    e.preventDefault();
    setErr(null);
    try {
      const r = await api.computeExtraSurplusValue({ individual_value: iv, social_value: sv, quantity: qty });
      setResult(r);
    } catch (e) {
      setErr(e instanceof Error ? e.message : String(e));
    }
  }

  return (
    <section className="card">
      <h2>Extra Surplus-Value Probe</h2>
      <p className="description">
        When a capitalist's individual value falls below the social value, they gain an extra
        surplus-value by selling at the higher social value. This temporary advantage drives
        generalisation of the new productivity. §1: "The capitalist … realises an extra
        surplus-value."
      </p>
      <div style={{ marginBottom: "0.75rem", display: "flex", gap: "0.5rem", flexWrap: "wrap" }}>
        {EXTRA_FIXTURES.map((f) => (
          <button key={f.label} type="button" onClick={() => loadFixture(f)}>{f.label}</button>
        ))}
      </div>
      <form className="form-grid" onSubmit={submit}>
        <label>
          <span>Individual Value (min per unit)</span>
          <input type="number" min={1} value={iv} onChange={(e) => setIv(Number(e.target.value))} />
        </label>
        <label>
          <span>Social Value (min per unit)</span>
          <input type="number" min={1} value={sv} onChange={(e) => setSv(Number(e.target.value))} />
        </label>
        <label>
          <span>Quantity produced</span>
          <input type="number" min={1} value={qty} onChange={(e) => setQty(Number(e.target.value))} />
        </label>
        <div className="form-actions span2">
          <button type="submit" className="primary">Compute</button>
          {err && <span className="error">{err}</span>}
        </div>
      </form>
      {result && (
        <table className="data-table">
          <tbody>
            <tr><td>Individual Value</td><td>{result.individual_value} min/unit</td></tr>
            <tr><td>Social Value</td><td>{result.social_value} min/unit</td></tr>
            <tr><td>Quantity</td><td>{result.quantity}</td></tr>
            <tr><td>Per-Unit Extra Gain</td><td>{result.per_unit} min</td></tr>
            <tr>
              <td><strong>Total Extra Surplus-Value</strong></td>
              <td><strong>{result.extra_surplus_value} min</strong></td>
            </tr>
            {result.extra_surplus_value === 0 && (
              <tr>
                <td colSpan={2} className="muted">
                  No extra surplus-value — individual value ≥ social value.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      )}
    </section>
  );
}
```

---

## Task 11: Wire the panel into ChapterShell, update registry and docs

**Files:**
- Modify: `web/src/components/ChapterShell.tsx`
- Modify: `web/src/chapters/registry.ts`
- Modify: `docs/architecture.md`

- [ ] **Step 1: Add Ch12 import and quote to ChapterShell.tsx**

Add import at top (after Ch11 import):
```tsx
import { Ch12RelativeSurplusValue } from "../chapters/Ch12RelativeSurplusValue";
```

Add quote entry to the `QUOTES` object:
```tsx
ch12: "The shortening of the working-day is, therefore, by no means what is aimed at, in capitalist production, when labour is economised by increasing its productiveness.",
```

Add the render branch after the `ch11` branch in the conditional render:
```tsx
) : activeChapterId === "ch12" ? (
  <Ch12RelativeSurplusValue onSharedChanged={onSharedChanged} />
```

- [ ] **Step 2: Mark ch12 as done in registry.ts**

Change:
```ts
{ id: "ch12", number: 12, title: "The Concept of Relative Surplus-Value", part: "Part IV — The Production of Relative Surplus-Value", status: "pending" },
```
to:
```ts
{ id: "ch12", number: 12, title: "The Concept of Relative Surplus-Value", part: "Part IV — The Production of Relative Surplus-Value", status: "done" },
```

- [ ] **Step 3: Update docs/architecture.md roadmap**

Change:
```
| Ch. 12+   | Pending     | Relative surplus-value, cooperation, machinery, wages, accumulation | all            |
```
to:
```
| Ch. 12    | ✅ Done     | Relative surplus-value; WorkingDay, ShortenNecessaryLabour, RateOfSurplusValue, ExtraSurplusValue, ApplyProductivityToSNLT | simulation-engine |
| Ch. 13+   | Pending     | Co-operation, machinery, wages, accumulation | all            |
```

---

## Task 12: Run web checks

- [ ] **Step 1: Ask the user to run `cd web && npm run lint && npm run build`**

Expected: TypeScript typecheck passes, Vite build succeeds with no errors.

---

## Self-Review

### Spec Coverage

| Spec requirement | Covered by |
|---|---|
| `WorkingDay` struct with `NecessaryLabour + SurplusLabour == Total` | Task 1 (test), Task 2 (impl) |
| `AbsoluteSurplusValue`, `RelativeSurplusValue`, `LabourPowerValue`, etc. as `LabourMinutes` aliases | Task 2 |
| `ShortenNecessaryLabour(WorkingDay, LabourPowerValue) WorkingDay` | Task 2 |
| `RateOfSurplusValue(SurplusLabour, NecessaryLabour) float64` | Task 2 |
| `ExtraSurplusValue(IndividualValue, SocialValue, Quantity) LabourMinutes` | Task 2 |
| `ApplyProductivityToSNLT` — inverse law | Task 2 |
| All §1 fixtures as tests | Tasks 1, 3 |
| All invariants as tests | Tasks 1, 3 |
| `POST /v1/production/working-day` | Tasks 3, 4, 5 |
| `POST /v1/production/working-day/shorten` | Tasks 3, 4, 5 |
| `GET /v1/production/rate-of-surplus-value` | Tasks 3, 4, 5 |
| `POST /v1/production/extra-surplus-value` | Tasks 3, 4, 5 |
| Gateway proxy for `/v1/production/*` | Task 7 |
| React panel with working-day split bar | Task 10 |
| React panel rate-of-surplus-value readout | Task 10 |
| React panel extra-surplus-value probe | Task 10 |
| Registry `ch12` → `"done"` | Task 11 |
| `docs/architecture.md` updated | Task 11 |

### Type Consistency Check

- `production.WorkingDay` fields: `Total LabourMinutes`, `NecessaryLabour NecessaryLabour`, `SurplusLabour SurplusLabour` — used consistently in handler via `int64` casts.
- `dtoFromWorkingDay` converts all fields to `int64` for JSON — used in both `RecordWorkingDay` and `ShortenWorkingDay` responses.
- `ShortenWorkingDayResponse.WorkingDay` field type is `workingDayDTO` — matches decode in test.
- TypeScript `ProductionWorkingDay` mirrors `workingDayDTO` snake_case fields.
- `ShortenWorkingDayResponse` mirrors `shortenWorkingDayResponse` Go struct.
- `ExtraSurplusValueResult` mirrors `extraSurplusResponse`.
- API methods return the correct TypeScript types.
- `Ch12RelativeSurplusValue` component imports from `../api` and `../types` correctly.
- ChapterShell imports `Ch12RelativeSurplusValue` from the correct path.
