// Package httpapi — handler tests for Vol. III Ch. 47 (Genesis of Capitalist
// Ground-Rent) endpoints.
//
// Fixtures drawn from Marx, Capital Vol. III Ch. 47:
//   - RentHistoricalForm stage 1: labour-rent, Normandy, 10th–13th c.
//   - HistoricalRentTransition: labour-rent (1) → product-rent (2);
//     from_stage must be strictly less than to_stage.
//   - SmallPeasantProduction: rent=2, profit=1, wages=3 → server computes
//     total_income_bp=6 (the peasant's undifferentiated gross return).
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

// newGenesisTestServer builds a minimal httptest.Server that registers only
// the Vol. III Ch. 47 genesis rent routes.
func newGenesisTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/rent/historical-forms", h.CreateRentHistoricalForm)
	mux.HandleFunc("GET /v1/rent/historical-forms", h.ListRentHistoricalForms)
	mux.HandleFunc("POST /v1/rent/historical-transitions", h.CreateHistoricalRentTransition)
	mux.HandleFunc("GET /v1/rent/historical-transitions", h.ListHistoricalRentTransitions)
	mux.HandleFunc("POST /v1/rent/small-peasant-productions", h.CreateSmallPeasantProduction)
	mux.HandleFunc("GET /v1/rent/small-peasant-productions", h.ListSmallPeasantProductions)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

// ── POST /v1/rent/historical-forms ───────────────────────────────────────────

// TestCreateRentHistoricalForm_HappyPath: labour-rent stage 1, Normandy → 201.
func TestCreateRentHistoricalForm_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	body := `{"stage":1,"name":"labour-rent","coercion_kind":"extra-economic","surplus_form":"direct_labour","example_region":"Normandy","example_period":"10th-13th c"}`
	res, err := http.Post(ts.URL+"/v1/rent/historical-forms", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.RentHistoricalForm
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	if saved.Stage != rent.RentFormLabour {
		t.Errorf("stage = %d, want %d (RentFormLabour)", saved.Stage, rent.RentFormLabour)
	}
}

// TestCreateRentHistoricalForm_StageZero_Returns422: stage=0 is not a valid
// RentFormStage — IsValid() returns false.
func TestCreateRentHistoricalForm_StageZero_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	body := `{"stage":0,"name":"labour-rent","coercion_kind":"extra-economic","surplus_form":"direct_labour","example_region":"Normandy","example_period":"10th-13th c"}`
	res, err := http.Post(ts.URL+"/v1/rent/historical-forms", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for stage=0", res.StatusCode)
	}
}

// TestCreateRentHistoricalForm_EmptyName_Returns422: name is required.
func TestCreateRentHistoricalForm_EmptyName_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	body := `{"stage":1,"name":"","coercion_kind":"extra-economic","surplus_form":"direct_labour","example_region":"Normandy","example_period":"10th-13th c"}`
	res, err := http.Post(ts.URL+"/v1/rent/historical-forms", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty name", res.StatusCode)
	}
}

// TestCreateRentHistoricalForm_MalformedJSON_Returns400.
func TestCreateRentHistoricalForm_MalformedJSON_Returns400(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/historical-forms", "application/json",
		strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for malformed JSON", res.StatusCode)
	}
}

// ── POST /v1/rent/historical-transitions ─────────────────────────────────────

// TestCreateHistoricalRentTransition_HappyPath: labour-rent (1) → product-rent
// (2) → 201.
func TestCreateHistoricalRentTransition_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	body := `{"from_stage":1,"to_stage":2,"preconditions":"development of commodity exchange","driving_force":"rising peasant productivity"}`
	res, err := http.Post(ts.URL+"/v1/rent/historical-transitions", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.HistoricalRentTransition
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	if saved.FromStage != rent.RentFormLabour {
		t.Errorf("from_stage = %d, want %d (RentFormLabour)", saved.FromStage, rent.RentFormLabour)
	}
	if saved.ToStage != rent.RentFormProduct {
		t.Errorf("to_stage = %d, want %d (RentFormProduct)", saved.ToStage, rent.RentFormProduct)
	}
}

// TestCreateHistoricalRentTransition_FromGeToStage_Returns422: transitions are
// always progressive (from < to); from_stage=2, to_stage=1 must fail.
func TestCreateHistoricalRentTransition_FromGeToStage_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	body := `{"from_stage":2,"to_stage":1,"preconditions":"x","driving_force":"y"}`
	res, err := http.Post(ts.URL+"/v1/rent/historical-transitions", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for from_stage >= to_stage", res.StatusCode)
	}
}

// TestCreateHistoricalRentTransition_InvalidFromStage_Returns422: stage 9 is
// out-of-range.
func TestCreateHistoricalRentTransition_InvalidFromStage_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	body := `{"from_stage":9,"to_stage":1,"preconditions":"x","driving_force":"y"}`
	res, err := http.Post(ts.URL+"/v1/rent/historical-transitions", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for from_stage=9 (out of range)", res.StatusCode)
	}
}

// ── POST /v1/rent/small-peasant-productions ───────────────────────────────────

// TestCreateSmallPeasantProduction_HappyPath: rent=2, profit=1, wages=3 →
// 201; server computes total_income_bp=6.
func TestCreateSmallPeasantProduction_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	body := `{"parcel_id":"normandy-plot","rent_component_bp":2,"profit_component_bp":1,"wages_component_bp":3,"is_conflated":true}`
	res, err := http.Post(ts.URL+"/v1/rent/small-peasant-productions", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.SmallPeasantProduction
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	// Server must compute total_income_bp = rent + profit + wages = 6.
	if saved.TotalIncomeBP != 6 {
		t.Errorf("total_income_bp = %d, want 6 (server-computed 2+1+3)", saved.TotalIncomeBP)
	}
}

// TestCreateSmallPeasantProduction_EmptyParcelID_Returns422.
func TestCreateSmallPeasantProduction_EmptyParcelID_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	body := `{"parcel_id":"","rent_component_bp":2,"profit_component_bp":1,"wages_component_bp":3,"is_conflated":true}`
	res, err := http.Post(ts.URL+"/v1/rent/small-peasant-productions", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty parcel_id", res.StatusCode)
	}
}

// TestCreateSmallPeasantProduction_NegativeRentComponent_Returns422: component
// values must be non-negative.
func TestCreateSmallPeasantProduction_NegativeRentComponent_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	body := `{"parcel_id":"normandy-plot","rent_component_bp":-1,"profit_component_bp":1,"wages_component_bp":3,"is_conflated":true}`
	res, err := http.Post(ts.URL+"/v1/rent/small-peasant-productions", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for rent_component_bp=-1", res.StatusCode)
	}
}

// TestCreateSmallPeasantProduction_MalformedJSON_Returns400.
func TestCreateSmallPeasantProduction_MalformedJSON_Returns400(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/small-peasant-productions", "application/json",
		strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for malformed JSON", res.StatusCode)
	}
}

// ── GET list endpoints — never null ──────────────────────────────────────────

func TestListRentHistoricalForms_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/historical-forms")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.RentHistoricalForm `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

func TestListHistoricalRentTransitions_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/historical-transitions")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.HistoricalRentTransition `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

func TestListSmallPeasantProductions_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newGenesisTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/small-peasant-productions")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.SmallPeasantProduction `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}
