// Package httpapi — handler tests for Vol. III Ch. 46 (Building-Site Rent,
// Rent in Mining, Price of Land) endpoints.
//
// Fixtures drawn from Marx, Capital Vol. III Ch. 46:
//   - Building-site rent: London Mayfair, location_grade=5, lease=99 years.
//   - Mining rent: Cornwall tin mine, ore_grade=5, surplus profit 1200 bp.
//   - Monopoly-price rent: Burgundy grand cru, value=200 LM, sale=400 LM
//     → monopoly_rent_labour_minutes=200 (server-computed).
//   - Land-price scenario (kind=1, falling interest): rent=10 bp, rate falls
//     from 500 bp to 250 bp → new_land_price_bp=400 (server-computed via
//     CapitalisedLandPrice(10, 250) = 400).
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

// newLandPriceTestServer builds a minimal httptest.Server that registers
// only the Vol. III Ch. 46 land-price / rent routes.
func newLandPriceTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/rent/building-rents", h.CreateBuildingSiteRent)
	mux.HandleFunc("GET /v1/rent/building-rents", h.ListBuildingSiteRents)
	mux.HandleFunc("POST /v1/rent/mining-rents", h.CreateMiningRent)
	mux.HandleFunc("GET /v1/rent/mining-rents", h.ListMiningRents)
	mux.HandleFunc("POST /v1/rent/monopoly-rents", h.CreateMonopolyPriceRent)
	mux.HandleFunc("GET /v1/rent/monopoly-rents", h.ListMonopolyPriceRents)
	mux.HandleFunc("POST /v1/rent/land-price-scenarios", h.CreateLandPriceScenario)
	mux.HandleFunc("GET /v1/rent/land-price-scenarios", h.ListLandPriceScenarios)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

// ── POST /v1/rent/building-rents ─────────────────────────────────────────────

// TestCreateBuildingSiteRent_HappyPath: London Mayfair, grade 5, lease 99 → 201.
func TestCreateBuildingSiteRent_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	body := `{"parcel_id":"P1","location_grade":5,"annual_rent_labour_minutes":300,"lease_term_years":99}`
	res, err := http.Post(ts.URL+"/v1/rent/building-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.BuildingSiteRent
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	if saved.LocationGrade != 5 {
		t.Errorf("location_grade = %d, want 5", saved.LocationGrade)
	}
}

func TestCreateBuildingSiteRent_EmptyParcelID_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	body := `{"parcel_id":"","location_grade":5,"annual_rent_labour_minutes":300,"lease_term_years":99}`
	res, err := http.Post(ts.URL+"/v1/rent/building-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty parcel_id", res.StatusCode)
	}
}

func TestCreateBuildingSiteRent_ZeroLocationGrade_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	body := `{"parcel_id":"P1","location_grade":0,"annual_rent_labour_minutes":300,"lease_term_years":99}`
	res, err := http.Post(ts.URL+"/v1/rent/building-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for location_grade=0", res.StatusCode)
	}
}

// ── POST /v1/rent/mining-rents ────────────────────────────────────────────────

// TestCreateMiningRent_HappyPath: Cornwall tin mine, ore grade 1, zero rent → 201.
func TestCreateMiningRent_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	body := `{"parcel_id":"P1","ore_grade":1,"annual_rent_labour_minutes":0,"surplus_profit_bp":0}`
	res, err := http.Post(ts.URL+"/v1/rent/mining-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.MiningRent
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	if saved.OreGrade != 1 {
		t.Errorf("ore_grade = %d, want 1", saved.OreGrade)
	}
}

func TestCreateMiningRent_ZeroOreGrade_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	body := `{"parcel_id":"P1","ore_grade":0,"annual_rent_labour_minutes":0,"surplus_profit_bp":0}`
	res, err := http.Post(ts.URL+"/v1/rent/mining-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for ore_grade=0", res.StatusCode)
	}
}

// ── POST /v1/rent/monopoly-rents ──────────────────────────────────────────────

// TestCreateMonopolyPriceRent_HappyPath: burgundy grand cru, value=200, sale=400
// → server computes monopoly_rent_labour_minutes=200.
func TestCreateMonopolyPriceRent_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	body := `{"parcel_id":"P1","product_kind":"burgundy_grand_cru","value_labour_minutes":200,"monopoly_sale_price_labour_minutes":400}`
	res, err := http.Post(ts.URL+"/v1/rent/monopoly-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.MonopolyPriceRent
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	if saved.MonopolyRentLabourMinutes != 200 {
		t.Errorf("monopoly_rent_labour_minutes = %d, want 200 (server-computed)",
			saved.MonopolyRentLabourMinutes)
	}
}

func TestCreateMonopolyPriceRent_EmptyProductKind_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	body := `{"parcel_id":"P1","product_kind":"","value_labour_minutes":200,"monopoly_sale_price_labour_minutes":400}`
	res, err := http.Post(ts.URL+"/v1/rent/monopoly-rents", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty product_kind", res.StatusCode)
	}
}

func TestCreateMonopolyPriceRent_MalformedJSON_Returns400(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/monopoly-rents", "application/json",
		strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for malformed JSON", res.StatusCode)
	}
}

// ── POST /v1/rent/land-price-scenarios ────────────────────────────────────────

// TestCreateLandPriceScenario_HappyPath: kind=1 (falling interest), rent=10 bp,
// old rate=500 bp, new rate=250 bp → server computes new_land_price_bp=400.
func TestCreateLandPriceScenario_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	body := `{"kind":1,"parcel_id":"P1","old_annual_rent_bp":10,"new_annual_rent_bp":10,"old_interest_rate_bp":500,"new_interest_rate_bp":250,"old_land_price_bp":200}`
	res, err := http.Post(ts.URL+"/v1/rent/land-price-scenarios", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.LandPriceScenario
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	if saved.NewLandPriceBP != 400 {
		t.Errorf("new_land_price_bp = %d, want 400 (server-computed via CapitalisedLandPrice(10,250))",
			saved.NewLandPriceBP)
	}
}

func TestCreateLandPriceScenario_InvalidKind_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	body := `{"kind":9,"parcel_id":"P1","old_annual_rent_bp":10,"new_annual_rent_bp":10,"old_interest_rate_bp":500,"new_interest_rate_bp":250,"old_land_price_bp":200}`
	res, err := http.Post(ts.URL+"/v1/rent/land-price-scenarios", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for kind=9 (out of range)", res.StatusCode)
	}
}

func TestCreateLandPriceScenario_ZeroNewInterestRate_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	body := `{"kind":1,"parcel_id":"P1","old_annual_rent_bp":10,"new_annual_rent_bp":10,"old_interest_rate_bp":500,"new_interest_rate_bp":0,"old_land_price_bp":200}`
	res, err := http.Post(ts.URL+"/v1/rent/land-price-scenarios", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for new_interest_rate_bp=0", res.StatusCode)
	}
}

// ── GET list endpoints — never null ──────────────────────────────────────────

func TestListBuildingSiteRents_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/building-rents")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.BuildingSiteRent `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

func TestListMiningRents_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/mining-rents")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.MiningRent `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

func TestListMonopolyPriceRents_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/monopoly-rents")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.MonopolyPriceRent `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

func TestListLandPriceScenarios_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newLandPriceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/land-price-scenarios")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.LandPriceScenario `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}
