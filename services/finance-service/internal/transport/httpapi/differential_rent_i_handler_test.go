// Package httpapi — handler tests for Vol. III Ch. 39 Differential Rent I
// endpoints. Fixtures drawn from Marx Vol. III Ch. 39, Table I.
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

// newDR1TestServer builds a minimal httptest.Server that registers only
// the Vol. III Ch. 39 DR I and location-rent-factor routes.
func newDR1TestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/rent/dr1-tables", h.CreateDifferentialRentITable)
	mux.HandleFunc("GET /v1/rent/dr1-tables", h.ListDifferentialRentITables)
	mux.HandleFunc("GET /v1/rent/dr1-tables/{id}", h.GetDifferentialRentITable)
	mux.HandleFunc("GET /v1/rent/dr1-tables/{id}/entries", h.ListDifferentialRentIEntries)
	mux.HandleFunc("POST /v1/rent/location-rent-factors", h.CreateLocationRentFactor)
	mux.HandleFunc("GET /v1/rent/location-rent-factors", h.ListLocationRentFactors)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

// marxTableIBody returns the JSON body for Marx Table I: 4 soil grades,
// capital £2.50 (250p), price £3/qtr (300p), output A=100, B=200, C=300, D=400.
const marxTableIBody = `{
	"description": "Marx Table I — Ch. 39: Differential Rent I",
	"entries": [
		{"grade":1,"capital_advanced_bp":250,"output_quarters":100,"price_of_production_bp":300},
		{"grade":2,"capital_advanced_bp":250,"output_quarters":200,"price_of_production_bp":300},
		{"grade":3,"capital_advanced_bp":250,"output_quarters":300,"price_of_production_bp":300},
		{"grade":4,"capital_advanced_bp":250,"output_quarters":400,"price_of_production_bp":300}
	]
}`

// ── POST /v1/rent/dr1-tables — happy path ─────────────────────────────────────

func TestCreateDR1Table_MarxTableI(t *testing.T) {
	t.Parallel()
	ts, _ := newDR1TestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/dr1-tables", "application/json", strings.NewReader(marxTableIBody))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if !strings.HasPrefix(loc, "/v1/rent/dr1-tables/") {
		t.Errorf("Location = %q, want /v1/rent/dr1-tables/ prefix", loc)
	}

	var body struct {
		Table   rent.DifferentialRentITable   `json:"table"`
		Entries []rent.DifferentialRentIEntry `json:"entries"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if body.Table.TotalRentalBP != 1800 {
		t.Errorf("table.total_rental_bp = %d, want 1800", body.Table.TotalRentalBP)
	}
	if body.Table.RegulatingGrade != rent.SoilGradeA {
		t.Errorf("table.regulating_grade = %d, want %d (SoilGradeA)", body.Table.RegulatingGrade, rent.SoilGradeA)
	}
	if len(body.Entries) != 4 {
		t.Fatalf("entries len = %d, want 4", len(body.Entries))
	}

	// Verify Location points to the returned table ID.
	if body.Table.ID == "" {
		t.Fatal("table.id must not be empty")
	}
	wantLoc := "/v1/rent/dr1-tables/" + string(body.Table.ID)
	if loc != wantLoc {
		t.Errorf("Location = %q, want %q", loc, wantLoc)
	}

	// Grade D must have money_rent_bp == 900.
	var dEntry rent.DifferentialRentIEntry
	for _, e := range body.Entries {
		if e.Grade == rent.SoilGradeD {
			dEntry = e
			break
		}
	}
	if dEntry.Grade != rent.SoilGradeD {
		t.Fatal("no grade-D entry in response")
	}
	if dEntry.MoneyRentBP != 900 {
		t.Errorf("grade-D money_rent_bp = %d, want 900", dEntry.MoneyRentBP)
	}
}

// ── POST /v1/rent/dr1-tables — validation failures ────────────────────────────

func TestCreateDR1Table_EmptyDescription(t *testing.T) {
	t.Parallel()
	ts, _ := newDR1TestServer(t)

	body := `{"description":"","entries":[{"grade":1,"capital_advanced_bp":250,"output_quarters":100,"price_of_production_bp":300}]}`
	res, err := http.Post(ts.URL+"/v1/rent/dr1-tables", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestCreateDR1Table_EmptyEntries(t *testing.T) {
	t.Parallel()
	ts, _ := newDR1TestServer(t)

	body := `{"description":"test","entries":[]}`
	res, err := http.Post(ts.URL+"/v1/rent/dr1-tables", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestCreateDR1Table_InvalidGrade(t *testing.T) {
	t.Parallel()
	ts, _ := newDR1TestServer(t)

	body := `{"description":"test","entries":[{"grade":9,"capital_advanced_bp":250,"output_quarters":100,"price_of_production_bp":300}]}`
	res, err := http.Post(ts.URL+"/v1/rent/dr1-tables", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for grade=9", res.StatusCode)
	}
}

func TestCreateDR1Table_MalformedJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newDR1TestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/dr1-tables", "application/json", strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// ── GET /v1/rent/dr1-tables ───────────────────────────────────────────────────

func TestListDR1Tables_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDR1TestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/dr1-tables")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.DifferentialRentITable `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

// ── Round-trip: POST then GET /v1/rent/dr1-tables/{id} ───────────────────────

func TestDR1TableRoundTrip(t *testing.T) {
	t.Parallel()
	ts, _ := newDR1TestServer(t)

	// Create.
	postRes, err := http.Post(ts.URL+"/v1/rent/dr1-tables", "application/json", strings.NewReader(marxTableIBody))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer postRes.Body.Close()
	if postRes.StatusCode != http.StatusCreated {
		t.Fatalf("post status = %d, want 201", postRes.StatusCode)
	}
	var created struct {
		Table   rent.DifferentialRentITable   `json:"table"`
		Entries []rent.DifferentialRentIEntry `json:"entries"`
	}
	if err := json.NewDecoder(postRes.Body).Decode(&created); err != nil {
		t.Fatalf("decode post response: %v", err)
	}
	tableID := string(created.Table.ID)

	// Get.
	getRes, err := http.Get(ts.URL + "/v1/rent/dr1-tables/" + tableID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", getRes.StatusCode)
	}
	var got struct {
		Table   rent.DifferentialRentITable   `json:"table"`
		Entries []rent.DifferentialRentIEntry `json:"entries"`
	}
	if err := json.NewDecoder(getRes.Body).Decode(&got); err != nil {
		t.Fatalf("decode get response: %v", err)
	}
	if string(got.Table.ID) != tableID {
		t.Errorf("get table.id = %q, want %q", got.Table.ID, tableID)
	}
	if len(got.Entries) != len(created.Entries) {
		t.Errorf("get entries len = %d, want %d", len(got.Entries), len(created.Entries))
	}
}

func TestGetDR1Table_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newDR1TestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/dr1-tables/nonexistent-id-0000")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// ── GET /v1/rent/dr1-tables/{id}/entries ─────────────────────────────────────

func TestListDR1Entries_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDR1TestServer(t)

	// Create a table first.
	postRes, err := http.Post(ts.URL+"/v1/rent/dr1-tables", "application/json", strings.NewReader(marxTableIBody))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer postRes.Body.Close()
	var created struct {
		Table rent.DifferentialRentITable `json:"table"`
	}
	if err := json.NewDecoder(postRes.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}

	res, err := http.Get(ts.URL + "/v1/rent/dr1-tables/" + string(created.Table.ID) + "/entries")
	if err != nil {
		t.Fatalf("get entries: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.DifferentialRentIEntry `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

func TestListDR1Entries_UnknownTable(t *testing.T) {
	t.Parallel()
	ts, _ := newDR1TestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/dr1-tables/no-such-table/entries")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404 for unknown table", res.StatusCode)
	}
}

// ── POST /v1/rent/location-rent-factors ──────────────────────────────────────

func TestCreateLocationRentFactor_Valid(t *testing.T) {
	t.Parallel()
	ts, _ := newDR1TestServer(t)

	body := `{"parcel_id":"suffolk-parcel-001","transport_cost_bp":50,"rent_equivalent_bp":50}`
	res, err := http.Post(ts.URL+"/v1/rent/location-rent-factors", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if !strings.HasPrefix(loc, "/v1/rent/location-rent-factors/") {
		t.Errorf("Location = %q, want /v1/rent/location-rent-factors/ prefix", loc)
	}

	var saved rent.LocationRentFactor
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("expected non-empty id in response")
	}
	if saved.ParcelID != "suffolk-parcel-001" {
		t.Errorf("parcel_id = %q, want %q", saved.ParcelID, "suffolk-parcel-001")
	}
	wantLoc := "/v1/rent/location-rent-factors/" + string(saved.ID)
	if loc != wantLoc {
		t.Errorf("Location = %q, want %q", loc, wantLoc)
	}
}

func TestCreateLocationRentFactor_EmptyParcelID(t *testing.T) {
	t.Parallel()
	ts, _ := newDR1TestServer(t)

	body := `{"parcel_id":"","transport_cost_bp":50,"rent_equivalent_bp":50}`
	res, err := http.Post(ts.URL+"/v1/rent/location-rent-factors", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// ── GET /v1/rent/location-rent-factors ───────────────────────────────────────

func TestListLocationRentFactors_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDR1TestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/location-rent-factors")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.LocationRentFactor `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}
