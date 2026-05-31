// Package httpapi — handler tests for Vol. III Ch. 40 Differential Rent II
// endpoints. Fixtures drawn from Marx Vol. III Ch. 40: two successive capital
// tranches on Soil D, each £2.50 (250bp), 4 qtr output (400), surplus-profit
// £9 (900bp) → total rent £18 (1800bp).
package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"log/slog"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/rent"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

// newDR2TestServer builds a minimal httptest.Server that registers only
// the Vol. III Ch. 40 DR II routes.
func newDR2TestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/rent/successive-investments", h.CreateSuccessiveInvestment)
	mux.HandleFunc("GET /v1/rent/successive-investments", h.ListSuccessiveInvestments)
	mux.HandleFunc("GET /v1/rent/dr2-tables/{parcel_id}", h.GetDR2Table)
	mux.HandleFunc("POST /v1/rent/lease-terms", h.CreateLeaseTerm)
	mux.HandleFunc("GET /v1/rent/lease-terms", h.ListLeaseTerms)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

// marxDR2InvestmentBody returns a valid successive-investment JSON body
// drawn from Marx Ch. 40 Soil D fixture.
const marxDR2InvestmentBody = `{"parcel_id":"P1","tranche_number":1,"capital_bp":250,"output_quarters":400,"surplus_profit_bp":900}`

// ── POST /v1/rent/successive-investments — happy path ────────────────────────

func TestCreateSuccessiveInvestment_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2TestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/successive-investments", "application/json",
		strings.NewReader(marxDR2InvestmentBody))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if !strings.HasPrefix(loc, "/v1/rent/successive-investments/") {
		t.Errorf("Location = %q, want /v1/rent/successive-investments/ prefix", loc)
	}

	var saved rent.SuccessiveInvestment
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	wantLoc := "/v1/rent/successive-investments/" + string(saved.ID)
	if loc != wantLoc {
		t.Errorf("Location = %q, want %q", loc, wantLoc)
	}
}

// ── POST /v1/rent/successive-investments — validation failures ───────────────

func TestCreateSuccessiveInvestment_EmptyParcelID(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2TestServer(t)

	body := `{"parcel_id":"","tranche_number":1,"capital_bp":250,"output_quarters":400,"surplus_profit_bp":900}`
	res, err := http.Post(ts.URL+"/v1/rent/successive-investments", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty parcel_id", res.StatusCode)
	}
}

func TestCreateSuccessiveInvestment_ZeroTranche(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2TestServer(t)

	body := `{"parcel_id":"P1","tranche_number":0,"capital_bp":250,"output_quarters":400,"surplus_profit_bp":900}`
	res, err := http.Post(ts.URL+"/v1/rent/successive-investments", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for tranche_number=0", res.StatusCode)
	}
}

func TestCreateSuccessiveInvestment_ZeroCapital(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2TestServer(t)

	body := `{"parcel_id":"P1","tranche_number":1,"capital_bp":0,"output_quarters":400,"surplus_profit_bp":900}`
	res, err := http.Post(ts.URL+"/v1/rent/successive-investments", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for capital_bp=0", res.StatusCode)
	}
}

func TestCreateSuccessiveInvestment_MalformedJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2TestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/successive-investments", "application/json",
		strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for malformed JSON", res.StatusCode)
	}
}

// ── GET /v1/rent/successive-investments ─────────────────────────────────────

func TestListSuccessiveInvestments_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2TestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/successive-investments")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.SuccessiveInvestment `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

// ── GET /v1/rent/dr2-tables/{parcel_id} ─────────────────────────────────────

// TestGetDR2Table_TwoTranches posts two Marx Ch. 40 successive investments
// (both Soil D: capital 250bp, output 400, surplus-profit 900bp each) for
// parcel "PX", then fetches the dr2-table. Expects total_rent_bp=1800 and
// two investments in the response.
func TestGetDR2Table_TwoTranches(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2TestServer(t)

	tranch1 := `{"parcel_id":"PX","tranche_number":1,"capital_bp":250,"output_quarters":400,"surplus_profit_bp":900}`
	tranch2 := `{"parcel_id":"PX","tranche_number":2,"capital_bp":250,"output_quarters":400,"surplus_profit_bp":900}`

	for i, b := range []string{tranch1, tranch2} {
		r, err := http.Post(ts.URL+"/v1/rent/successive-investments", "application/json",
			bytes.NewReader([]byte(b)))
		if err != nil {
			t.Fatalf("post tranche %d: %v", i+1, err)
		}
		r.Body.Close()
		if r.StatusCode != http.StatusCreated {
			t.Fatalf("post tranche %d status = %d, want 201", i+1, r.StatusCode)
		}
	}

	res, err := http.Get(ts.URL + "/v1/rent/dr2-tables/PX")
	if err != nil {
		t.Fatalf("get dr2-table: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}

	var body struct {
		Table       rent.DifferentialRentIITable `json:"table"`
		Investments []rent.SuccessiveInvestment   `json:"investments"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if body.Table.TotalRentBP != 1800 {
		t.Errorf("table.total_rent_bp = %d, want 1800", body.Table.TotalRentBP)
	}
	if len(body.Investments) != 2 {
		t.Errorf("investments len = %d, want 2", len(body.Investments))
	}
}

// TestGetDR2Table_UnknownParcel checks that an unknown parcel yields 200 with
// a zero-rent table and a non-nil empty investments array — no 404.
func TestGetDR2Table_UnknownParcel(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2TestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/dr2-tables/unknown-parcel")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200 (not 404) for unknown parcel", res.StatusCode)
	}

	var body struct {
		Table       rent.DifferentialRentIITable `json:"table"`
		Investments []rent.SuccessiveInvestment   `json:"investments"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Table.TotalRentBP != 0 {
		t.Errorf("table.total_rent_bp = %d, want 0 for unknown parcel", body.Table.TotalRentBP)
	}
	if body.Investments == nil {
		t.Error("investments must be a non-nil array, not null")
	}
	if len(body.Investments) != 0 {
		t.Errorf("investments len = %d, want 0 for unknown parcel", len(body.Investments))
	}
}

// ── POST /v1/rent/lease-terms ────────────────────────────────────────────────

func TestCreateLeaseTerm_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2TestServer(t)

	body := `{"parcel_id":"P1","start_year":1867,"duration_years":21,"fixed_rent_bp":900}`
	res, err := http.Post(ts.URL+"/v1/rent/lease-terms", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if !strings.HasPrefix(loc, "/v1/rent/lease-terms/") {
		t.Errorf("Location = %q, want /v1/rent/lease-terms/ prefix", loc)
	}

	var saved rent.LeaseTerm
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
	wantLoc := "/v1/rent/lease-terms/" + string(saved.ID)
	if loc != wantLoc {
		t.Errorf("Location = %q, want %q", loc, wantLoc)
	}
}

func TestCreateLeaseTerm_EmptyParcelID(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2TestServer(t)

	body := `{"parcel_id":"","start_year":1867,"duration_years":21,"fixed_rent_bp":900}`
	res, err := http.Post(ts.URL+"/v1/rent/lease-terms", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty parcel_id", res.StatusCode)
	}
}

// ── GET /v1/rent/lease-terms ─────────────────────────────────────────────────

func TestListLeaseTerms_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2TestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/lease-terms")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.LeaseTerm `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}
