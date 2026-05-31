package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

// newDiffRentTestServer builds a minimal httptest.Server that registers only
// the Vol. III Ch. 38 differential-rent routes.
func newDiffRentTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/rent/surplus-profits", h.CreatePoPSurplusProfit)
	mux.HandleFunc("GET /v1/rent/surplus-profits", h.ListPoPSurplusProfits)
	mux.HandleFunc("POST /v1/rent/monopolised-forces", h.CreateMonopolisedNaturalForce)
	mux.HandleFunc("GET /v1/rent/monopolised-forces", h.ListMonopolisedNaturalForces)
	mux.HandleFunc("POST /v1/rent/differential-rents", h.CreateDifferentialRent)
	mux.HandleFunc("GET /v1/rent/differential-rents", h.ListDifferentialRents)
	mux.HandleFunc("POST /v1/rent/capitalised-prices", h.CreateCapitalisedRentPrice)
	mux.HandleFunc("GET /v1/rent/capitalised-prices", h.ListCapitalisedRentPrices)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

// ── POST /v1/rent/surplus-profits ────────────────────────────────────────────

func TestCreatePoPSurplusProfit_Valid(t *testing.T) {
	t.Parallel()
	ts, _ := newDiffRentTestServer(t)

	body := `{"capital_id":"waterfall-capital-001","general_price_of_production_bp":11500,"individual_price_of_production_bp":9000}`
	res, err := http.Post(ts.URL+"/v1/rent/surplus-profits", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if !strings.HasPrefix(loc, "/v1/rent/surplus-profits/") {
		t.Errorf("Location = %q, want /v1/rent/surplus-profits/ prefix", loc)
	}

	var s rent.PriceOfProductionSurplusProfit
	if err := json.NewDecoder(res.Body).Decode(&s); err != nil {
		t.Fatalf("decode: %v", err)
	}
	// Server must compute the surplus-profit regardless of what the client sent.
	if s.SurplusProfitBP != 2500 {
		t.Errorf("surplus_profit_bp = %d, want 2500 (server-computed)", s.SurplusProfitBP)
	}
	if s.ID == "" {
		t.Error("expected a non-empty id in response body")
	}
	wantLoc := "/v1/rent/surplus-profits/" + string(s.ID)
	if loc != wantLoc {
		t.Errorf("Location = %q, want %q", loc, wantLoc)
	}
}

func TestCreatePoPSurplusProfit_EmptyCapitalID(t *testing.T) {
	t.Parallel()
	ts, _ := newDiffRentTestServer(t)

	body := `{"capital_id":"","general_price_of_production_bp":11500,"individual_price_of_production_bp":9000}`
	res, err := http.Post(ts.URL+"/v1/rent/surplus-profits", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestCreatePoPSurplusProfit_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newDiffRentTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/surplus-profits", "application/json", strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// ── GET /v1/rent/surplus-profits ─────────────────────────────────────────────

func TestListPoPSurplusProfits_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDiffRentTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/surplus-profits")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.PriceOfProductionSurplusProfit `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

// ── POST /v1/rent/capitalised-prices ─────────────────────────────────────────

func TestCreateCapitalisedRentPrice_Valid(t *testing.T) {
	t.Parallel()
	ts, _ := newDiffRentTestServer(t)

	// Marx fixture: £10 annual rent (4800 LM) at 5% (500 bp) → 96000 LM.
	body := `{"parcel_id":"highland-parcel-001","annual_rent_labour_minutes":4800,"interest_rate_bp":500}`
	res, err := http.Post(ts.URL+"/v1/rent/capitalised-prices", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if !strings.HasPrefix(loc, "/v1/rent/capitalised-prices/") {
		t.Errorf("Location = %q, want /v1/rent/capitalised-prices/ prefix", loc)
	}

	var crp rent.CapitalisedRentPrice
	if err := json.NewDecoder(res.Body).Decode(&crp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	// Server must compute the capitalised price regardless of any client value.
	if crp.CapitalisedPriceLabourMinutes != 96000 {
		t.Errorf("capitalised_price_labour_minutes = %d, want 96000 (server-computed)", crp.CapitalisedPriceLabourMinutes)
	}
}

func TestCreateCapitalisedRentPrice_ZeroInterestRate(t *testing.T) {
	t.Parallel()
	ts, _ := newDiffRentTestServer(t)

	body := `{"parcel_id":"highland-parcel-002","annual_rent_labour_minutes":4800,"interest_rate_bp":0}`
	res, err := http.Post(ts.URL+"/v1/rent/capitalised-prices", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// ── GET /v1/rent/capitalised-prices ──────────────────────────────────────────

func TestListCapitalisedRentPrices_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDiffRentTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/capitalised-prices")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.CapitalisedRentPrice `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

// ── POST /v1/rent/monopolised-forces ─────────────────────────────────────────

func TestCreateMonopolisedNaturalForce_Valid(t *testing.T) {
	t.Parallel()
	ts, _ := newDiffRentTestServer(t)

	body := `{"kind":"waterfall","is_monopolisable":true,"parcel_id":"highland-parcel-003"}`
	res, err := http.Post(ts.URL+"/v1/rent/monopolised-forces", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if !strings.HasPrefix(loc, "/v1/rent/monopolised-forces/") {
		t.Errorf("Location = %q, want /v1/rent/monopolised-forces/ prefix", loc)
	}
}

func TestCreateMonopolisedNaturalForce_EmptyKind(t *testing.T) {
	t.Parallel()
	ts, _ := newDiffRentTestServer(t)

	body := `{"kind":"","is_monopolisable":true,"parcel_id":"highland-parcel-004"}`
	res, err := http.Post(ts.URL+"/v1/rent/monopolised-forces", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// ── GET /v1/rent/monopolised-forces ──────────────────────────────────────────

func TestListMonopolisedNaturalForces_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDiffRentTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/monopolised-forces")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.MonopolisedNaturalForce `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

// ── POST /v1/rent/differential-rents ─────────────────────────────────────────

func TestCreateDifferentialRent_Valid(t *testing.T) {
	t.Parallel()
	ts, _ := newDiffRentTestServer(t)

	body := `{"parcel_id":"highland-parcel-005","surplus_profit_bp":2500,"annual_rent_labour_minutes":4800}`
	res, err := http.Post(ts.URL+"/v1/rent/differential-rents", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if !strings.HasPrefix(loc, "/v1/rent/differential-rents/") {
		t.Errorf("Location = %q, want /v1/rent/differential-rents/ prefix", loc)
	}
}

func TestCreateDifferentialRent_EmptyParcelID(t *testing.T) {
	t.Parallel()
	ts, _ := newDiffRentTestServer(t)

	body := `{"parcel_id":"","surplus_profit_bp":2500,"annual_rent_labour_minutes":4800}`
	res, err := http.Post(ts.URL+"/v1/rent/differential-rents", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// ── GET /v1/rent/differential-rents ──────────────────────────────────────────

func TestListDifferentialRents_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDiffRentTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/differential-rents")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.DifferentialRent `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}
