// Package httpapi — handler tests for Vol. III Ch. 44 Differential Rent on
// the Worst Cultivated Soil endpoints.
//
// Fixture: old_price_of_production_bp=30000 (£3.00), new_price_of_production_bp=35000 (£3.50)
// → rent_per_acre_bp=5000 (£0.50). Source: ComputeWorstSoilRentPerAcre(30000,35000)==5000.
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

// newWorstSoilTestServer builds a minimal httptest.Server that registers
// only the Vol. III Ch. 44 worst-soil-rent routes.
func newWorstSoilTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/rent/worst-soil-rents", h.CreateRentOnWorstSoil)
	mux.HandleFunc("GET /v1/rent/worst-soil-rents", h.ListRentOnWorstSoils)
	mux.HandleFunc("POST /v1/rent/soil-improvements", h.CreateSoilImprovementCapital)
	mux.HandleFunc("GET /v1/rent/soil-improvements", h.ListSoilImprovementCapitals)
	mux.HandleFunc("POST /v1/rent/regulating-price-shifts", h.CreateRegulatingPriceShift)
	mux.HandleFunc("GET /v1/rent/regulating-price-shifts", h.ListRegulatingPriceShifts)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

// ── POST /v1/rent/worst-soil-rents — happy path ──────────────────────────────

// TestCreateRentOnWorstSoil_HappyPath uses Marx Ch. 44 mechanism 1:
// old_price_of_production_bp=30000 (£3.00), new_price_of_production_bp=35000 (£3.50)
// → server computes rent_per_acre_bp=5000 (£0.50).
func TestCreateRentOnWorstSoil_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	body := `{"parcel_id":"P1","mechanism":1,"old_price_of_production_bp":30000,"new_price_of_production_bp":35000}`
	res, err := http.Post(ts.URL+"/v1/rent/worst-soil-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/rent/worst-soil-rents") {
		t.Errorf("Location = %q, want prefix /v1/rent/worst-soil-rents", loc)
	}

	var saved rent.RentOnWorstSoil
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	if saved.RentPerAcreBP != 5000 {
		t.Errorf("rent_per_acre_bp = %d, want 5000", saved.RentPerAcreBP)
	}
}

// ── POST /v1/rent/worst-soil-rents — validation failures ────────────────────

func TestCreateRentOnWorstSoil_MechanismZero_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	body := `{"parcel_id":"P1","mechanism":0,"old_price_of_production_bp":30000,"new_price_of_production_bp":35000}`
	res, err := http.Post(ts.URL+"/v1/rent/worst-soil-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for mechanism=0", res.StatusCode)
	}
}

func TestCreateRentOnWorstSoil_MechanismFour_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	body := `{"parcel_id":"P1","mechanism":4,"old_price_of_production_bp":30000,"new_price_of_production_bp":35000}`
	res, err := http.Post(ts.URL+"/v1/rent/worst-soil-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for mechanism=4", res.StatusCode)
	}
}

func TestCreateRentOnWorstSoil_EmptyParcelID_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	body := `{"parcel_id":"","mechanism":1,"old_price_of_production_bp":30000,"new_price_of_production_bp":35000}`
	res, err := http.Post(ts.URL+"/v1/rent/worst-soil-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty parcel_id", res.StatusCode)
	}
}

func TestCreateRentOnWorstSoil_ZeroOldPrice_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	body := `{"parcel_id":"P1","mechanism":1,"old_price_of_production_bp":0,"new_price_of_production_bp":35000}`
	res, err := http.Post(ts.URL+"/v1/rent/worst-soil-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for old_price_of_production_bp=0", res.StatusCode)
	}
}

func TestCreateRentOnWorstSoil_MalformedJSON_Returns400(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/worst-soil-rents", "application/json",
		strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for malformed JSON", res.StatusCode)
	}
}

// ── GET /v1/rent/worst-soil-rents ────────────────────────────────────────────

func TestListRentOnWorstSoils_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/worst-soil-rents")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.RentOnWorstSoil `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

// ── POST /v1/rent/soil-improvements — happy path ─────────────────────────────

// TestCreateSoilImprovementCapital_HappyPath models a drainage investment:
// capital_invested_bp=3000, productivity_gain_bp=500, amortisation_years=15,
// amortised_years=0 (newly recorded).
func TestCreateSoilImprovementCapital_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	body := `{"parcel_id":"P1","capital_invested_bp":3000,"productivity_gain_bp":500,"amortisation_years":15,"amortised_years":0}`
	res, err := http.Post(ts.URL+"/v1/rent/soil-improvements", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.SoilImprovementCapital
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
}

// ── POST /v1/rent/soil-improvements — validation failures ────────────────────

func TestCreateSoilImprovementCapital_EmptyParcelID_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	body := `{"parcel_id":"","capital_invested_bp":3000,"productivity_gain_bp":500,"amortisation_years":15,"amortised_years":0}`
	res, err := http.Post(ts.URL+"/v1/rent/soil-improvements", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty parcel_id", res.StatusCode)
	}
}

func TestCreateSoilImprovementCapital_ZeroAmortisationYears_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	body := `{"parcel_id":"P1","capital_invested_bp":3000,"productivity_gain_bp":500,"amortisation_years":0,"amortised_years":0}`
	res, err := http.Post(ts.URL+"/v1/rent/soil-improvements", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for amortisation_years=0", res.StatusCode)
	}
}

// ── GET /v1/rent/soil-improvements ──────────────────────────────────────────

func TestListSoilImprovementCapitals_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/soil-improvements")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.SoilImprovementCapital `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

// ── POST /v1/rent/regulating-price-shifts — happy path ───────────────────────

// TestCreateRegulatingPriceShift_HappyPath models Marx Ch. 44: regulating grade
// moves from A (grade 1) to B (grade 2), old_price_bp=30000, new_price_bp=35000.
func TestCreateRegulatingPriceShift_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	body := `{"from_grade":1,"to_grade":2,"old_price_bp":30000,"new_price_bp":35000}`
	res, err := http.Post(ts.URL+"/v1/rent/regulating-price-shifts", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.RegulatingPriceShift
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
}

// ── POST /v1/rent/regulating-price-shifts — validation failures ──────────────

func TestCreateRegulatingPriceShift_InvalidFromGrade_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	// grade 9 is outside the valid range 1–4
	body := `{"from_grade":9,"to_grade":2,"old_price_bp":30000,"new_price_bp":35000}`
	res, err := http.Post(ts.URL+"/v1/rent/regulating-price-shifts", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for from_grade=9", res.StatusCode)
	}
}

func TestCreateRegulatingPriceShift_ZeroOldPrice_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	body := `{"from_grade":1,"to_grade":2,"old_price_bp":0,"new_price_bp":35000}`
	res, err := http.Post(ts.URL+"/v1/rent/regulating-price-shifts", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for old_price_bp=0", res.StatusCode)
	}
}

// ── GET /v1/rent/regulating-price-shifts ─────────────────────────────────────

func TestListRegulatingPriceShifts_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newWorstSoilTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/regulating-price-shifts")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.RegulatingPriceShift `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}
