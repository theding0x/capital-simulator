// Package httpapi — handler tests for Vol. III Ch. 45 (Absolute Ground-Rent)
// endpoints.
//
// Fixtures drawn from Marx, Capital Vol. III Ch. 45:
//   - Agricultural composition C=60, V=40 → composition_ratio_bp=15000 < social
//     average 40000 → is_below_average=true.
//   - Value of agricultural output = 140 labour-minutes, price of production =
//     120 → gap = 20.
//   - Surplus-value = 40 bp, average profit = 20 bp → absolute rent = 20 bp.
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

// newAbsoluteRentTestServer builds a minimal httptest.Server that registers
// only the Vol. III Ch. 45 absolute-rent routes.
func newAbsoluteRentTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/rent/compositions", h.CreateAgriculturalComposition)
	mux.HandleFunc("GET /v1/rent/compositions", h.ListAgriculturalCompositions)
	mux.HandleFunc("POST /v1/rent/value-gaps", h.CreateValuePriceGap)
	mux.HandleFunc("GET /v1/rent/value-gaps", h.ListValuePriceGaps)
	mux.HandleFunc("POST /v1/rent/absolute-rents", h.CreateAbsoluteRent)
	mux.HandleFunc("GET /v1/rent/absolute-rents", h.ListAbsoluteRents)
	mux.HandleFunc("GET /v1/rent/absolute-rents/{id}", h.GetAbsoluteRent)
	mux.HandleFunc("POST /v1/rent/rent-limits", h.CreateAbsoluteRentLimit)
	mux.HandleFunc("GET /v1/rent/rent-limits", h.ListAbsoluteRentLimits)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

// ── POST /v1/rent/compositions ───────────────────────────────────────────────

// TestCreateAgriculturalComposition_HappyPath uses Marx Ch. 45 agricultural
// capital C=60, V=40, social average=40000 → composition_ratio_bp=15000,
// is_below_average=true.
func TestCreateAgriculturalComposition_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	body := `{"constant_capital_bp":60,"variable_capital_bp":40,"social_average_ratio_bp":40000}`
	res, err := http.Post(ts.URL+"/v1/rent/compositions", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.AgriculturalComposition
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	if saved.CompositionRatioBP != 15000 {
		t.Errorf("composition_ratio_bp = %d, want 15000", saved.CompositionRatioBP)
	}
	if !saved.IsBelowAverage {
		t.Error("is_below_average must be true when ratio 15000 < social average 40000")
	}
}

func TestCreateAgriculturalComposition_ZeroVariableCapital_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	body := `{"constant_capital_bp":60,"variable_capital_bp":0,"social_average_ratio_bp":40000}`
	res, err := http.Post(ts.URL+"/v1/rent/compositions", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for variable_capital_bp=0", res.StatusCode)
	}
}

// ── POST /v1/rent/value-gaps ─────────────────────────────────────────────────

// TestCreateValuePriceGap_HappyPath: agricultural output value=140, price of
// production=120 → gap_labour_minutes=20.
func TestCreateValuePriceGap_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	body := `{"product_id":"agricultural-output","value_labour_minutes":140,"price_of_production_labour_minutes":120}`
	res, err := http.Post(ts.URL+"/v1/rent/value-gaps", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.ValuePriceGap
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	if saved.GapLabourMinutes != 20 {
		t.Errorf("gap_labour_minutes = %d, want 20", saved.GapLabourMinutes)
	}
}

func TestCreateValuePriceGap_EmptyProductID_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	body := `{"product_id":"","value_labour_minutes":140,"price_of_production_labour_minutes":120}`
	res, err := http.Post(ts.URL+"/v1/rent/value-gaps", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty product_id", res.StatusCode)
	}
}

// ── POST /v1/rent/absolute-rents ─────────────────────────────────────────────

// TestCreateAbsoluteRent_HappyPath: surplus_value=40, average_profit=20
// → absolute_rent_labour_minutes=20. Captures the returned ID for subsequent GET tests.
func TestCreateAbsoluteRent_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	body := `{"parcel_id":"P1","value_price_gap_id":"VG1","surplus_value_labour_minutes":40,"average_profit_labour_minutes":20}`
	res, err := http.Post(ts.URL+"/v1/rent/absolute-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.AbsoluteRent
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	if saved.AbsoluteRentLabourMinutes != 20 {
		t.Errorf("absolute_rent_labour_minutes = %d, want 20", saved.AbsoluteRentLabourMinutes)
	}
}

// TestGetAbsoluteRent_RoundTrip creates a record then GETs it by the returned ID.
func TestGetAbsoluteRent_RoundTrip(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	// create
	body := `{"parcel_id":"P1","value_price_gap_id":"VG1","surplus_value_labour_minutes":40,"average_profit_labour_minutes":20}`
	createRes, err := http.Post(ts.URL+"/v1/rent/absolute-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("create post: %v", err)
	}
	defer createRes.Body.Close()
	if createRes.StatusCode != http.StatusCreated {
		t.Fatalf("create status = %d, want 201", createRes.StatusCode)
	}

	var created rent.AbsoluteRent
	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	// get by id
	getRes, err := http.Get(ts.URL + "/v1/rent/absolute-rents/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", getRes.StatusCode)
	}

	var got rent.AbsoluteRent
	if err := json.NewDecoder(getRes.Body).Decode(&got); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("get id = %q, want %q", got.ID, created.ID)
	}
	if got.AbsoluteRentLabourMinutes != 20 {
		t.Errorf("get absolute_rent_labour_minutes = %d, want 20", got.AbsoluteRentLabourMinutes)
	}
}

func TestGetAbsoluteRent_NotFound_Returns404(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/absolute-rents/nope")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404 for unknown id", res.StatusCode)
	}
}

func TestCreateAbsoluteRent_EmptyParcelID_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	body := `{"parcel_id":"","value_price_gap_id":"VG1","surplus_value_labour_minutes":40,"average_profit_labour_minutes":20}`
	res, err := http.Post(ts.URL+"/v1/rent/absolute-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty parcel_id", res.StatusCode)
	}
}

func TestCreateAbsoluteRent_MalformedJSON_Returns400(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/absolute-rents", "application/json",
		strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for malformed JSON", res.StatusCode)
	}
}

// ── POST /v1/rent/rent-limits ────────────────────────────────────────────────

// TestCreateAbsoluteRentLimit_HappyPath: max=20, actual=20 → is_at_limit=true.
func TestCreateAbsoluteRentLimit_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	body := `{"max_rent_bp":20,"actual_rent_bp":20}`
	res, err := http.Post(ts.URL+"/v1/rent/rent-limits", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.AbsoluteRentLimit
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	if !saved.IsAtLimit {
		t.Error("is_at_limit must be true when actual_rent_bp == max_rent_bp")
	}
}

// ── GET list endpoints — never null ──────────────────────────────────────────

func TestListAgriculturalCompositions_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/compositions")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.AgriculturalComposition `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

func TestListValuePriceGaps_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/value-gaps")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.ValuePriceGap `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

func TestListAbsoluteRents_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/absolute-rents")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.AbsoluteRent `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

func TestListAbsoluteRentLimits_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newAbsoluteRentTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/rent-limits")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.AbsoluteRentLimit `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}
