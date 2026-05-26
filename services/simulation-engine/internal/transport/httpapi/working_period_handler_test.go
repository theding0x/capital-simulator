package httpapi_test

// Vol. II Ch. 12 — The Working Period handler tests.
// Fixtures drawn from Marx's text: SpinningMill1871 (discrete, 1-day yarn),
// LocomotiveWorks1867 (connected, 84-day locomotive), and the London builder
// credit-financing example.

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
	wp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/workperiod"
)

func newWorkingPeriodHandler() *httpapi.Handler {
	return httpapi.New(nil, httpapi.Deps{WorkingPeriods: store.NewMemory()})
}

// --- CreateWorkingPeriod -----------------------------------------------------

// TestCreateWorkingPeriod_SpinningMill1871 verifies 201 for a discrete spinning mill.
// Marx: yarn is produced each day — multiplier = 1 regardless of day count.
func TestCreateWorkingPeriod_SpinningMill1871(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()
	body := map[string]any{
		"industrial_capital_id":      "spinning-mill-1871",
		"commodity_id":               "yarn-001",
		"working_day_count":          1,
		"mode":                       "discrete",
		"working_day_length_minutes": 600,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.CreateWorkingPeriod(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if loc := rec.Header().Get("Location"); loc == "" {
		t.Error("expected non-empty Location header")
	}
	var resp wp.WorkingPeriod
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID == "" {
		t.Error("expected non-empty id")
	}
	if resp.Mode != wp.ModeDiscrete {
		t.Errorf("mode: want discrete, got %s", resp.Mode)
	}
}

// TestCreateWorkingPeriod_LocomotiveWorks1867 verifies 201 for a connected locomotive works.
// Marx: 84 working-days required — multiplier = 84.
func TestCreateWorkingPeriod_LocomotiveWorks1867(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()
	body := map[string]any{
		"industrial_capital_id":      "locomotive-works-1867",
		"commodity_id":               "locomotive-001",
		"working_day_count":          84,
		"mode":                       "connected",
		"working_day_length_minutes": 600,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.CreateWorkingPeriod(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp wp.WorkingPeriod
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Multiplier() != 84 {
		t.Errorf("multiplier: want 84, got %d", resp.Multiplier())
	}
}

// TestCreateWorkingPeriod_InvalidMode verifies 400 for an unrecognised mode.
func TestCreateWorkingPeriod_InvalidMode(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()
	body := map[string]any{
		"industrial_capital_id":      "mill",
		"commodity_id":               "yarn",
		"working_day_count":          1,
		"mode":                       "batch",
		"working_day_length_minutes": 600,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateWorkingPeriod(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

// TestCreateWorkingPeriod_MissingFields verifies 400 when required fields are absent.
func TestCreateWorkingPeriod_MissingFields(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()
	body := map[string]any{"mode": "discrete"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateWorkingPeriod(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

// TestCreateWorkingPeriod_NaturalConstraint verifies 400 when a natural
// constraint blocks premature delivery. Register a >= 365-day crop constraint
// then try to create a 30-day period for the same commodity.
func TestCreateWorkingPeriod_NaturalConstraint(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()

	// Register a natural constraint: winter wheat takes >= 365 days.
	ncBody := map[string]any{
		"commodity_id":     "winter-wheat-001",
		"min_working_days": 365,
		"reason":           "annual crop; cannot be harvested before growing season completes",
	}
	b, _ := json.Marshal(ncBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/natural-constraints", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateNaturalConstraint(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create constraint: want 201, got %d: %s", rec.Code, rec.Body.String())
	}

	// Now try to create a 30-day working period for the same commodity.
	wpBody := map[string]any{
		"industrial_capital_id":      "wheat-farm-1867",
		"commodity_id":               "winter-wheat-001",
		"working_day_count":          30,
		"mode":                       "connected",
		"working_day_length_minutes": 600,
	}
	b, _ = json.Marshal(wpBody)
	req = httptest.NewRequest(http.MethodPost, "/v1/working-periods", bytes.NewReader(b))
	rec = httptest.NewRecorder()
	h.CreateWorkingPeriod(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400 (ErrPrematureDelivery), got %d: %s", rec.Code, rec.Body.String())
	}
}

// --- GetWorkingPeriod ---------------------------------------------------------

// TestGetWorkingPeriod_NotFound verifies 404 for an unknown id.
func TestGetWorkingPeriod_NotFound(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()
	req := httptest.NewRequest(http.MethodGet, "/v1/working-periods/doesnotexist", nil)
	req.SetPathValue("id", "doesnotexist")
	rec := httptest.NewRecorder()
	h.GetWorkingPeriod(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("want 404, got %d", rec.Code)
	}
}

// TestGetWorkingPeriod_RoundTrip verifies create -> get returns the same record.
func TestGetWorkingPeriod_RoundTrip(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()
	body := map[string]any{
		"industrial_capital_id":      "spinning-mill-rt",
		"commodity_id":               "yarn-rt",
		"working_day_count":          1,
		"mode":                       "discrete",
		"working_day_length_minutes": 480,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateWorkingPeriod(rec, req)
	var created wp.WorkingPeriod
	_ = json.NewDecoder(rec.Body).Decode(&created)

	req = httptest.NewRequest(http.MethodGet, "/v1/working-periods/"+string(created.ID), nil)
	req.SetPathValue("id", string(created.ID))
	rec = httptest.NewRecorder()
	h.GetWorkingPeriod(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var got wp.WorkingPeriod
	_ = json.NewDecoder(rec.Body).Decode(&got)
	if got.ID != created.ID {
		t.Errorf("id mismatch: want %s, got %s", created.ID, got.ID)
	}
}

// --- ListWorkingPeriods -------------------------------------------------------

// TestListWorkingPeriods_Empty verifies 200 with a non-null (empty) slice.
func TestListWorkingPeriods_Empty(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()
	req := httptest.NewRequest(http.MethodGet, "/v1/working-periods", nil)
	rec := httptest.NewRecorder()
	h.ListWorkingPeriods(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200, got %d", rec.Code)
	}
	var list []wp.WorkingPeriod
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if list == nil {
		t.Error("list must be non-null slice, not null")
	}
}

// --- RecordWorkingPeriodInterruption -----------------------------------------

// TestRecordInterruption_Connected verifies 201 for a connected-period interruption
// with deterioration (idle locomotive rusts).
func TestRecordInterruption_Connected(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()

	// Create a connected period.
	wpBody := map[string]any{
		"industrial_capital_id":      "loco-works-intr",
		"commodity_id":               "loco-002",
		"working_day_count":          84,
		"mode":                       "connected",
		"working_day_length_minutes": 600,
	}
	b, _ := json.Marshal(wpBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateWorkingPeriod(rec, req)
	var period wp.WorkingPeriod
	_ = json.NewDecoder(rec.Body).Decode(&period)

	// Record an interruption with deterioration.
	intrBody := map[string]any{
		"deterioration_pence":                500,
		"waste_of_means_of_production_pence": 200,
	}
	b, _ = json.Marshal(intrBody)
	req = httptest.NewRequest(http.MethodPost, "/v1/working-periods/"+string(period.ID)+"/interruptions", bytes.NewReader(b))
	req.SetPathValue("id", string(period.ID))
	rec = httptest.NewRecorder()
	h.RecordWorkingPeriodInterruption(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var intr wp.WorkingPeriodInterruption
	_ = json.NewDecoder(rec.Body).Decode(&intr)
	if intr.DeteriorationPence != 500 {
		t.Errorf("deterioration_pence: want 500, got %d", intr.DeteriorationPence)
	}
}

// TestRecordInterruption_DiscreteWithDeterioration verifies 400 when a discrete
// interruption carries deterioration_pence (invariant: discrete products do not sit idle).
func TestRecordInterruption_DiscreteWithDeterioration(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()

	wpBody := map[string]any{
		"industrial_capital_id":      "spinning-intr",
		"commodity_id":               "yarn-intr",
		"working_day_count":          1,
		"mode":                       "discrete",
		"working_day_length_minutes": 600,
	}
	b, _ := json.Marshal(wpBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateWorkingPeriod(rec, req)
	var period wp.WorkingPeriod
	_ = json.NewDecoder(rec.Body).Decode(&period)

	intrBody := map[string]any{
		"deterioration_pence":                100,
		"waste_of_means_of_production_pence": 0,
	}
	b, _ = json.Marshal(intrBody)
	req = httptest.NewRequest(http.MethodPost, "/v1/working-periods/"+string(period.ID)+"/interruptions", bytes.NewReader(b))
	req.SetPathValue("id", string(period.ID))
	rec = httptest.NewRecorder()
	h.RecordWorkingPeriodInterruption(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400 for discrete deterioration, got %d", rec.Code)
	}
}

// TestRecordInterruption_NotFound verifies 404 for an unknown working period.
func TestRecordInterruption_NotFound(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()
	body := map[string]any{
		"deterioration_pence":                0,
		"waste_of_means_of_production_pence": 0,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods/unknown/interruptions", bytes.NewReader(b))
	req.SetPathValue("id", "unknown")
	rec := httptest.NewRecorder()
	h.RecordWorkingPeriodInterruption(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("want 404, got %d", rec.Code)
	}
}

// --- RecordWorkingPeriodShortening -------------------------------------------

// TestRecordShortening_BakewellBreeding verifies 201 for a bakewell_breeding shortening
// (Marx's cattle maturation example -- no extra fixed capital required).
func TestRecordShortening_BakewellBreeding(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()

	wpBody := map[string]any{
		"industrial_capital_id":      "bakewell-farm",
		"commodity_id":               "cattle-001",
		"working_day_count":          365,
		"mode":                       "connected",
		"working_day_length_minutes": 600,
	}
	b, _ := json.Marshal(wpBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateWorkingPeriod(rec, req)
	var period wp.WorkingPeriod
	_ = json.NewDecoder(rec.Body).Decode(&period)

	shortBody := map[string]any{
		"kind":                           "bakewell_breeding",
		"old_working_day_count":          365,
		"new_working_day_count":          270,
		"additional_fixed_capital_pence": 0,
		"additional_labour_count":        0,
	}
	b, _ = json.Marshal(shortBody)
	req = httptest.NewRequest(http.MethodPost, "/v1/working-periods/"+string(period.ID)+"/shortenings", bytes.NewReader(b))
	req.SetPathValue("id", string(period.ID))
	rec = httptest.NewRecorder()
	h.RecordWorkingPeriodShortening(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var s wp.WorkingPeriodShortening
	_ = json.NewDecoder(rec.Body).Decode(&s)
	if s.Kind != wp.ShorteningBakewellBreeding {
		t.Errorf("kind: want bakewell_breeding, got %s", s.Kind)
	}
}

// TestRecordShortening_MachineryMissingFixedCapital verifies 400 when a machinery
// shortening omits additional_fixed_capital_pence.
func TestRecordShortening_MachineryMissingFixedCapital(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()

	wpBody := map[string]any{
		"industrial_capital_id":      "loco-mach",
		"commodity_id":               "loco-003",
		"working_day_count":          84,
		"mode":                       "connected",
		"working_day_length_minutes": 600,
	}
	b, _ := json.Marshal(wpBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateWorkingPeriod(rec, req)
	var period wp.WorkingPeriod
	_ = json.NewDecoder(rec.Body).Decode(&period)

	shortBody := map[string]any{
		"kind":                           "machinery",
		"old_working_day_count":          84,
		"new_working_day_count":          42,
		"additional_fixed_capital_pence": 0, // missing — machinery requires > 0
		"additional_labour_count":        0,
	}
	b, _ = json.Marshal(shortBody)
	req = httptest.NewRequest(http.MethodPost, "/v1/working-periods/"+string(period.ID)+"/shortenings", bytes.NewReader(b))
	req.SetPathValue("id", string(period.ID))
	rec = httptest.NewRecorder()
	h.RecordWorkingPeriodShortening(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

// --- RecordWorkingPeriodCreditFinancing --------------------------------------

// TestRecordCreditFinancing_LondonBuilder verifies 201 for the Marx London-builder
// credit-financing example (own capital = 0, borrowed = full value).
func TestRecordCreditFinancing_LondonBuilder(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()

	wpBody := map[string]any{
		"industrial_capital_id":      "london-builder",
		"commodity_id":               "house-001",
		"working_day_count":          200,
		"mode":                       "connected",
		"working_day_length_minutes": 600,
	}
	b, _ := json.Marshal(wpBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateWorkingPeriod(rec, req)
	var period wp.WorkingPeriod
	_ = json.NewDecoder(rec.Body).Decode(&period)

	cfBody := map[string]any{
		"own_capital_pence":      0,
		"borrowed_capital_pence": 1_000_000, // £10,000 as pence
		"mortgage_provider_id":   "bank-of-england-placeholder",
	}
	b, _ = json.Marshal(cfBody)
	req = httptest.NewRequest(http.MethodPost, "/v1/working-periods/"+string(period.ID)+"/credit-financing", bytes.NewReader(b))
	req.SetPathValue("id", string(period.ID))
	rec = httptest.NewRecorder()
	h.RecordWorkingPeriodCreditFinancing(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var cf wp.WorkingPeriodCreditFinancing
	_ = json.NewDecoder(rec.Body).Decode(&cf)
	if cf.BorrowedCapitalPence != 1_000_000 {
		t.Errorf("borrowed_capital_pence: want 1000000, got %d", cf.BorrowedCapitalPence)
	}
}

// TestRecordCreditFinancing_ZeroTotal verifies 400 when both own and borrowed are 0.
func TestRecordCreditFinancing_ZeroTotal(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()

	wpBody := map[string]any{
		"industrial_capital_id":      "builder-zero",
		"commodity_id":               "house-002",
		"working_day_count":          100,
		"mode":                       "connected",
		"working_day_length_minutes": 600,
	}
	b, _ := json.Marshal(wpBody)
	req := httptest.NewRequest(http.MethodPost, "/v1/working-periods", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateWorkingPeriod(rec, req)
	var period wp.WorkingPeriod
	_ = json.NewDecoder(rec.Body).Decode(&period)

	cfBody := map[string]any{
		"own_capital_pence":      0,
		"borrowed_capital_pence": 0,
	}
	b, _ = json.Marshal(cfBody)
	req = httptest.NewRequest(http.MethodPost, "/v1/working-periods/"+string(period.ID)+"/credit-financing", bytes.NewReader(b))
	req.SetPathValue("id", string(period.ID))
	rec = httptest.NewRecorder()
	h.RecordWorkingPeriodCreditFinancing(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

// --- CreateNaturalConstraint -------------------------------------------------

// TestCreateNaturalConstraint_WinterWheat verifies 201 for the annual-crop fixture.
func TestCreateNaturalConstraint_WinterWheat(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()
	body := map[string]any{
		"commodity_id":     "winter-wheat-nc",
		"min_working_days": 365,
		"reason":           "annual crop",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/natural-constraints", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.CreateNaturalConstraint(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	if loc := rec.Header().Get("Location"); loc == "" {
		t.Error("expected non-empty Location header")
	}
	var nc wp.NaturalWorkingPeriodConstraint
	_ = json.NewDecoder(rec.Body).Decode(&nc)
	if nc.MinWorkingDays != 365 {
		t.Errorf("min_working_days: want 365, got %d", nc.MinWorkingDays)
	}
}

// TestCreateNaturalConstraint_MissingReason verifies 400 when reason is omitted.
func TestCreateNaturalConstraint_MissingReason(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()
	body := map[string]any{
		"commodity_id":     "wheat-no-reason",
		"min_working_days": 365,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/natural-constraints", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateNaturalConstraint(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

// --- ListNaturalConstraints --------------------------------------------------

// TestListNaturalConstraints_Empty verifies 200 with a non-null (empty) slice.
func TestListNaturalConstraints_Empty(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()
	req := httptest.NewRequest(http.MethodGet, "/v1/natural-constraints", nil)
	rec := httptest.NewRecorder()
	h.ListNaturalConstraints(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200, got %d", rec.Code)
	}
	var list []wp.NaturalWorkingPeriodConstraint
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if list == nil {
		t.Error("list must be non-null slice, not null")
	}
}

// TestListNaturalConstraints_RoundTrip verifies that a created constraint appears in the list.
func TestListNaturalConstraints_RoundTrip(t *testing.T) {
	t.Parallel()

	h := newWorkingPeriodHandler()
	body := map[string]any{
		"commodity_id":     "timber-rt",
		"min_working_days": 1825, // 5-year timber
		"reason":           "multi-year forestry maturation",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/natural-constraints", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateNaturalConstraint(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: want 201, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/v1/natural-constraints", nil)
	rec = httptest.NewRecorder()
	h.ListNaturalConstraints(rec, req)

	var list []wp.NaturalWorkingPeriodConstraint
	_ = json.NewDecoder(rec.Body).Decode(&list)
	if len(list) != 1 {
		t.Errorf("want 1 constraint, got %d", len(list))
	}
	if list[0].MinWorkingDays != 1825 {
		t.Errorf("min_working_days: want 1825, got %d", list[0].MinWorkingDays)
	}
}
