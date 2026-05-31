// Package httpapi — handler tests for Vol. III Ch. 42 DR II Falling Price
// endpoints. Fixtures drawn from Marx Vol. III Ch. 42: Soil D excluded,
// new regulating grade 2, new price of production 540bp; two tranches
// each 250bp capital, 400qtrs output, 125bp surplus-profit →
// TotalGrainRentQtrs=250, TotalMoneyRentalBP=1350.
package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"log/slog"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

// newDR2FallingTestServer builds a minimal httptest.Server that registers
// only the Vol. III Ch. 42 DR II falling-price routes.
func newDR2FallingTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/rent/exclusion-events", h.CreateSoilExclusionEvent)
	mux.HandleFunc("GET /v1/rent/exclusion-events", h.ListSoilExclusionEvents)
	mux.HandleFunc("POST /v1/rent/dr2-falling-outcomes", h.CreateFallingPriceDR2Outcome)
	mux.HandleFunc("GET /v1/rent/dr2-falling-outcomes", h.ListFallingPriceDR2Outcomes)
	mux.HandleFunc("POST /v1/rent/normal-capital", h.CreateNormalCapitalPerAcre)
	mux.HandleFunc("GET /v1/rent/normal-capital", h.ListNormalCapitalPerAcre)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

// ── POST /v1/rent/exclusion-events — happy path ───────────────────────────────

func TestCreateSoilExclusionEvent_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	body := `{"excluded_grade":1,"new_regulating_grade":2,"new_price_of_production_bp":540}`
	res, err := http.Post(ts.URL+"/v1/rent/exclusion-events", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/rent/exclusion-events") {
		t.Errorf("Location = %q, want prefix /v1/rent/exclusion-events", loc)
	}

	var saved rent.SoilExclusionEvent
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
}

// ── POST /v1/rent/exclusion-events — validation failures ─────────────────────

func TestCreateSoilExclusionEvent_RegulatingLEExcluded_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	// new_regulating_grade (1) <= excluded_grade (2): invalid
	body := `{"excluded_grade":2,"new_regulating_grade":1,"new_price_of_production_bp":540}`
	res, err := http.Post(ts.URL+"/v1/rent/exclusion-events", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 when regulating <= excluded", res.StatusCode)
	}
}

func TestCreateSoilExclusionEvent_InvalidGrade_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	body := `{"excluded_grade":9,"new_regulating_grade":2,"new_price_of_production_bp":540}`
	res, err := http.Post(ts.URL+"/v1/rent/exclusion-events", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for excluded_grade=9", res.StatusCode)
	}
}

func TestCreateSoilExclusionEvent_ZeroPrice_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	body := `{"excluded_grade":1,"new_regulating_grade":2,"new_price_of_production_bp":0}`
	res, err := http.Post(ts.URL+"/v1/rent/exclusion-events", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for new_price_of_production_bp=0", res.StatusCode)
	}
}

func TestCreateSoilExclusionEvent_MalformedJSON_Returns400(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/exclusion-events", "application/json",
		strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for malformed JSON", res.StatusCode)
	}
}

// ── GET /v1/rent/exclusion-events ─────────────────────────────────────────────

func TestListSoilExclusionEvents_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/exclusion-events")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.SoilExclusionEvent `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

// ── POST /v1/rent/dr2-falling-outcomes — happy path ──────────────────────────

// TestCreateFallingPriceDR2Outcome_HappyPath uses Marx Ch. 42 values:
// Soil D excluded (grade 1 → grade 2), new price 540bp.
// Two tranches 250bp capital, 400qtrs output, 125bp surplus-profit each.
// TotalGrainRentQtrs = 125+125 = 250; TotalMoneyRentalBP = 250×540/100 = 1350.
func TestCreateFallingPriceDR2Outcome_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	// First create an exclusion event to obtain a real ID.
	evBody := `{"excluded_grade":1,"new_regulating_grade":2,"new_price_of_production_bp":540}`
	evRes, err := http.Post(ts.URL+"/v1/rent/exclusion-events", "application/json",
		strings.NewReader(evBody))
	if err != nil {
		t.Fatalf("create exclusion event: %v", err)
	}
	defer evRes.Body.Close()
	if evRes.StatusCode != http.StatusCreated {
		t.Fatalf("exclusion event status = %d, want 201", evRes.StatusCode)
	}
	var ev rent.SoilExclusionEvent
	if err := json.NewDecoder(evRes.Body).Decode(&ev); err != nil {
		t.Fatalf("decode exclusion event: %v", err)
	}

	outBody := `{` +
		`"variant":1,` +
		`"exclusion_event_id":"` + string(ev.ID) + `",` +
		`"investments":[` +
		`{"parcel_id":"PD","tranche_number":1,"capital_bp":250,"output_quarters":400,"surplus_profit_bp":125},` +
		`{"parcel_id":"PD","tranche_number":2,"capital_bp":250,"output_quarters":400,"surplus_profit_bp":125}` +
		`]}`
	outRes, err := http.Post(ts.URL+"/v1/rent/dr2-falling-outcomes", "application/json",
		strings.NewReader(outBody))
	if err != nil {
		t.Fatalf("post dr2-falling-outcomes: %v", err)
	}
	defer outRes.Body.Close()

	if outRes.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", outRes.StatusCode)
	}

	var saved rent.FallingPriceDR2Outcome
	if err := json.NewDecoder(outRes.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	if saved.TotalGrainRentQtrs != 250 {
		t.Errorf("total_grain_rent_qtrs = %d, want 250", saved.TotalGrainRentQtrs)
	}
	if saved.TotalMoneyRentalBP != 1350 {
		t.Errorf("total_money_rental_bp = %d, want 1350", saved.TotalMoneyRentalBP)
	}
}

// ── POST /v1/rent/dr2-falling-outcomes — 404 on missing exclusion event ──────

func TestCreateFallingPriceDR2Outcome_MissingExclusionEvent_Returns404(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	body := `{` +
		`"variant":1,` +
		`"exclusion_event_id":"does-not-exist",` +
		`"investments":[` +
		`{"parcel_id":"PD","tranche_number":1,"capital_bp":250,"output_quarters":400,"surplus_profit_bp":125}` +
		`]}`
	res, err := http.Post(ts.URL+"/v1/rent/dr2-falling-outcomes", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404 for missing exclusion_event_id", res.StatusCode)
	}
}

// ── POST /v1/rent/dr2-falling-outcomes — validation failures ─────────────────

func TestCreateFallingPriceDR2Outcome_InvalidVariant_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	body := `{` +
		`"variant":9,` +
		`"exclusion_event_id":"some-id",` +
		`"investments":[` +
		`{"parcel_id":"PD","tranche_number":1,"capital_bp":250,"output_quarters":400,"surplus_profit_bp":125}` +
		`]}`
	res, err := http.Post(ts.URL+"/v1/rent/dr2-falling-outcomes", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for variant=9", res.StatusCode)
	}
}

func TestCreateFallingPriceDR2Outcome_EmptyInvestments_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	body := `{"variant":1,"exclusion_event_id":"some-id","investments":[]}`
	res, err := http.Post(ts.URL+"/v1/rent/dr2-falling-outcomes", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty investments", res.StatusCode)
	}
}

func TestCreateFallingPriceDR2Outcome_MalformedJSON_Returns400(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/dr2-falling-outcomes", "application/json",
		strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for malformed JSON", res.StatusCode)
	}
}

// ── GET /v1/rent/dr2-falling-outcomes ────────────────────────────────────────

func TestListFallingPriceDR2Outcomes_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/dr2-falling-outcomes")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.FallingPriceDR2Outcome `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

// ── POST /v1/rent/normal-capital — happy path ─────────────────────────────────

func TestCreateNormalCapitalPerAcre_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	body := `{"period_label":"England pre-1848","capital_per_acre_bp":800}`
	res, err := http.Post(ts.URL+"/v1/rent/normal-capital", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.NormalCapitalPerAcre
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
}

// ── POST /v1/rent/normal-capital — validation failures ───────────────────────

func TestCreateNormalCapitalPerAcre_EmptyPeriodLabel_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	body := `{"period_label":"","capital_per_acre_bp":800}`
	res, err := http.Post(ts.URL+"/v1/rent/normal-capital", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty period_label", res.StatusCode)
	}
}

func TestCreateNormalCapitalPerAcre_ZeroCapital_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	body := `{"period_label":"England pre-1848","capital_per_acre_bp":0}`
	res, err := http.Post(ts.URL+"/v1/rent/normal-capital", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for capital_per_acre_bp=0", res.StatusCode)
	}
}

// ── GET /v1/rent/normal-capital ───────────────────────────────────────────────

func TestListNormalCapitalPerAcre_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2FallingTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/normal-capital")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.NormalCapitalPerAcre `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}
