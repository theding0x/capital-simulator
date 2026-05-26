package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
	tv "github.com/theding0x/capital-simulator/services/simulation-engine/internal/turnover"
)

func newAggregateTurnoverHandler() *httpapi.Handler {
	return httpapi.New(nil, httpapi.Deps{AggregateTurnovers: store.NewMemory()})
}

// TestCreateAggregateTurnover_MarxMillExample verifies 201 for the §3 £80k+£20k fixture
// and checks that aggregate_number_basis_points = 10800 (= 1 + 2/25 of £100k advanced).
func TestCreateAggregateTurnover_MarxMillExample(t *testing.T) {
	t.Parallel()

	h := newAggregateTurnoverHandler()
	body := map[string]any{
		"industrial_capital_id": "cotton-mill-1867",
		"contributions": []map[string]any{
			{
				"capital_component_id":        "comp-fixed-001",
				"kind":                        "fixed",
				"advanced_pence":              8_000_000,
				"turnover_number_basis_points": 1_000,
				"difference_kind":             "real",
			},
			{
				"capital_component_id":        "comp-circ-001",
				"kind":                        "circulating",
				"advanced_pence":              2_000_000,
				"turnover_number_basis_points": 50_000,
				"difference_kind":             "real",
			},
		},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/aggregate-turnovers", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.CreateAggregateTurnover(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var resp tv.AggregateTurnover
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.ID == "" {
		t.Error("expected non-empty id")
	}
	if resp.AggregateNumberBasisPoints != 10_800 {
		t.Errorf("aggregate_number_basis_points: want 10800, got %d", resp.AggregateNumberBasisPoints)
	}
	if resp.AnnualTurnedOverPence != 10_800_000 {
		t.Errorf("annual_turned_over_pence: want 10800000, got %d", resp.AnnualTurnedOverPence)
	}
	if resp.Lens != tv.LensMoney {
		t.Errorf("lens: want money, got %s", resp.Lens)
	}
}

// TestCreateAggregateTurnover_MissingIndustrialCapitalID verifies 400 when the
// required field is omitted.
func TestCreateAggregateTurnover_MissingIndustrialCapitalID(t *testing.T) {
	t.Parallel()

	h := newAggregateTurnoverHandler()
	body := map[string]any{
		"contributions": []map[string]any{
			{"kind": "fixed", "advanced_pence": 1_000_000, "turnover_number_basis_points": 1_000},
		},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/aggregate-turnovers", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateAggregateTurnover(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

// TestCreateAggregateTurnover_EmptyContributions verifies 400 when contributions is empty.
func TestCreateAggregateTurnover_EmptyContributions(t *testing.T) {
	t.Parallel()

	h := newAggregateTurnoverHandler()
	body := map[string]any{
		"industrial_capital_id": "mill",
		"contributions":         []map[string]any{},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/aggregate-turnovers", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateAggregateTurnover(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}

// TestGetAggregateTurnover_NotFound verifies 404 for an unknown id.
func TestGetAggregateTurnover_NotFound(t *testing.T) {
	t.Parallel()

	h := newAggregateTurnoverHandler()
	req := httptest.NewRequest(http.MethodGet, "/v1/aggregate-turnovers/doesnotexist", nil)
	req.SetPathValue("id", "doesnotexist")
	rec := httptest.NewRecorder()
	h.GetAggregateTurnover(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("want 404, got %d", rec.Code)
	}
}

// TestListAggregateTurnovers_Empty verifies 200 with a non-null (empty) slice.
func TestListAggregateTurnovers_Empty(t *testing.T) {
	t.Parallel()

	h := newAggregateTurnoverHandler()
	req := httptest.NewRequest(http.MethodGet, "/v1/aggregate-turnovers", nil)
	rec := httptest.NewRecorder()
	h.ListAggregateTurnovers(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("want 200, got %d", rec.Code)
	}
	var list []tv.AggregateTurnover
	if err := json.NewDecoder(rec.Body).Decode(&list); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if list == nil {
		t.Error("list must be non-null slice, not null")
	}
}

// TestLifetimeCycle_RecordAndCrisisPhase verifies the full lifetime-cycle + crisis-phase
// roundtrip: register a 10-year cycle then annotate it with a depression phase.
func TestLifetimeCycle_RecordAndCrisisPhase(t *testing.T) {
	t.Parallel()

	h := newAggregateTurnoverHandler()

	// Create an aggregate turnover.
	body := map[string]any{
		"industrial_capital_id": "mill-cycle-test",
		"contributions": []map[string]any{
			{"kind": "fixed", "advanced_pence": 8_000_000, "turnover_number_basis_points": 1_000},
		},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/aggregate-turnovers", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateAggregateTurnover(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: want 201, got %d", rec.Code)
	}
	var at tv.AggregateTurnover
	_ = json.NewDecoder(rec.Body).Decode(&at)

	// Record a 10-year lifetime cycle (10 × 525,600 = 5,256,000 minutes).
	lcBody := map[string]any{"duration_minutes": 5_256_000}
	b, _ = json.Marshal(lcBody)
	req = httptest.NewRequest(http.MethodPost, "/v1/aggregate-turnovers/"+string(at.ID)+"/lifetime-cycle", bytes.NewReader(b))
	req.SetPathValue("id", string(at.ID))
	rec = httptest.NewRecorder()
	h.RecordAggregateTurnoverLifetimeCycle(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("lifetime-cycle: want 201, got %d: %s", rec.Code, rec.Body.String())
	}

	// Annotate with a depression phase (Marx references 1857, 1866, 1873).
	phBody := map[string]any{"phase": "depression", "notes": "1857 depression"}
	b, _ = json.Marshal(phBody)
	req = httptest.NewRequest(http.MethodPost, "/v1/aggregate-turnovers/"+string(at.ID)+"/crisis-phase", bytes.NewReader(b))
	req.SetPathValue("id", string(at.ID))
	rec = httptest.NewRecorder()
	h.RecordCrisisPhase(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("crisis-phase: want 201, got %d: %s", rec.Code, rec.Body.String())
	}
	var phResp tv.CrisisCyclePhaseRecord
	_ = json.NewDecoder(rec.Body).Decode(&phResp)
	if phResp.Phase != tv.PhaseDepression {
		t.Errorf("phase: want depression, got %s", phResp.Phase)
	}

	// GET — lifetime cycle and latest crisis phase must appear.
	req = httptest.NewRequest(http.MethodGet, "/v1/aggregate-turnovers/"+string(at.ID), nil)
	req.SetPathValue("id", string(at.ID))
	rec = httptest.NewRecorder()
	h.GetAggregateTurnover(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("get: want 200, got %d", rec.Code)
	}
	var got tv.AggregateTurnover
	_ = json.NewDecoder(rec.Body).Decode(&got)
	if got.LifetimeCycle == nil {
		t.Fatal("expected lifetime_cycle on get")
	}
	if got.LifetimeCycle.DurationMinutes != 5_256_000 {
		t.Errorf("duration_minutes: want 5256000, got %d", got.LifetimeCycle.DurationMinutes)
	}
	if got.LatestCrisisPhase == nil || got.LatestCrisisPhase.Phase != tv.PhaseDepression {
		t.Error("expected latest_crisis_phase = depression")
	}
}

// TestRecordCrisisPhase_InvalidPhase verifies 400 for an unknown phase string.
func TestRecordCrisisPhase_InvalidPhase(t *testing.T) {
	t.Parallel()

	h := newAggregateTurnoverHandler()

	body := map[string]any{
		"industrial_capital_id": "mill-phase-invalid",
		"contributions": []map[string]any{
			{"kind": "fixed", "advanced_pence": 1_000_000, "turnover_number_basis_points": 1_000},
		},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/aggregate-turnovers", bytes.NewReader(b))
	rec := httptest.NewRecorder()
	h.CreateAggregateTurnover(rec, req)
	var at tv.AggregateTurnover
	_ = json.NewDecoder(rec.Body).Decode(&at)

	phBody := map[string]any{"phase": "boom"} // invalid value
	b, _ = json.Marshal(phBody)
	req = httptest.NewRequest(http.MethodPost, "/v1/aggregate-turnovers/"+string(at.ID)+"/crisis-phase", bytes.NewReader(b))
	req.SetPathValue("id", string(at.ID))
	rec = httptest.NewRecorder()
	h.RecordCrisisPhase(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rec.Code)
	}
}
