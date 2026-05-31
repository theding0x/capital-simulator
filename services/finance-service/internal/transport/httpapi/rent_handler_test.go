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

// newRentTestServer builds a minimal httptest.Server that only registers the
// Vol. III Ch. 37 rent routes. This avoids re-registering every prior chapter
// while still exercising the real Handler methods and Memory store.
func newRentTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/rent/parcels", h.CreateLandParcel)
	mux.HandleFunc("GET /v1/rent/parcels", h.ListLandParcels)
	mux.HandleFunc("GET /v1/rent/parcels/{id}", h.GetLandParcel)
	mux.HandleFunc("POST /v1/rent/rents", h.CreateGroundRent)
	mux.HandleFunc("GET /v1/rent/rents", h.ListGroundRents)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

// validParcelBody returns a minimal valid create-parcel JSON body.
// Fixture: fertility grade 2 Highland parcel owned by the Duke of Sutherland.
func validParcelBody() string {
	return `{"owner_id":"duke-of-sutherland-0001","fertility_grade":2,"area_hectares":50,"location":"Sutherland, Scotland"}`
}

// validRentBody returns a minimal valid create-rent JSON body.
// Fixture: tenant farmer pays 2400 LabourMinutes differential rent to the Duke.
func validRentBody(ownerID, capitalistID string) string {
	return `{"parcel_id":"parcel-0001","capitalist_id":"` + capitalistID + `","land_owner_id":"` + ownerID + `","form":1,"amount_labour_minutes":2400,"period_years":7}`
}

// ── POST /v1/rent/parcels ─────────────────────────────────────────────────────

func TestCreateLandParcel_Valid(t *testing.T) {
	t.Parallel()
	ts, _ := newRentTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/parcels", "application/json", strings.NewReader(validParcelBody()))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if !strings.HasPrefix(loc, "/v1/rent/parcels/") {
		t.Errorf("Location = %q, want /v1/rent/parcels/ prefix", loc)
	}
	var p rent.LandParcel
	if err := json.NewDecoder(res.Body).Decode(&p); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if p.ID == "" {
		t.Error("expected a non-empty id in response body")
	}
	// Location suffix must match the assigned ID.
	wantLoc := "/v1/rent/parcels/" + string(p.ID)
	if loc != wantLoc {
		t.Errorf("Location = %q, want %q", loc, wantLoc)
	}
}

func TestCreateLandParcel_FertilityGradeOutOfRange(t *testing.T) {
	t.Parallel()
	ts, _ := newRentTestServer(t)

	// fertility_grade=5 is outside [1,4]
	body := `{"owner_id":"duke-of-sutherland-0001","fertility_grade":5,"area_hectares":50,"location":"Sutherland"}`
	res, err := http.Post(ts.URL+"/v1/rent/parcels", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestCreateLandParcel_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newRentTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/parcels", "application/json", strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// ── GET /v1/rent/parcels ──────────────────────────────────────────────────────

func TestListLandParcels_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newRentTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/parcels")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.LandParcel `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

// ── GET /v1/rent/parcels/{id} ─────────────────────────────────────────────────

func TestGetLandParcel_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newRentTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/parcels/nonexistent-id")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// ── Round-trip: POST then GET by returned ID ──────────────────────────────────

func TestCreateThenGetLandParcel_RoundTrip(t *testing.T) {
	t.Parallel()
	ts, _ := newRentTestServer(t)

	// Create
	res, err := http.Post(ts.URL+"/v1/rent/parcels", "application/json", strings.NewReader(validParcelBody()))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create status = %d, want 201", res.StatusCode)
	}
	var created rent.LandParcel
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	// Fetch by ID
	got, err := http.Get(ts.URL + "/v1/rent/parcels/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched rent.LandParcel
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("id mismatch: got %s, want %s", fetched.ID, created.ID)
	}
	if fetched.FertilityGrade != created.FertilityGrade {
		t.Errorf("fertility_grade: got %d, want %d", fetched.FertilityGrade, created.FertilityGrade)
	}
	if fetched.OwnerID != created.OwnerID {
		t.Errorf("owner_id: got %q, want %q", fetched.OwnerID, created.OwnerID)
	}
}

// ── POST /v1/rent/rents ───────────────────────────────────────────────────────

func TestCreateGroundRent_Valid(t *testing.T) {
	t.Parallel()
	ts, _ := newRentTestServer(t)

	body := validRentBody("duke-of-sutherland-0001", "highland-tenant-farmer-01")
	res, err := http.Post(ts.URL+"/v1/rent/rents", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	loc := res.Header.Get("Location")
	if !strings.HasPrefix(loc, "/v1/rent/rents/") {
		t.Errorf("Location = %q, want /v1/rent/rents/ prefix", loc)
	}
}

func TestCreateGroundRent_SameOwnerOperator(t *testing.T) {
	t.Parallel()
	ts, _ := newRentTestServer(t)

	// LandOwnerID == CapitalistID triggers ErrSameOwnerOperator → 422
	body := validRentBody("same-person-0001", "same-person-0001")
	res, err := http.Post(ts.URL+"/v1/rent/rents", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestCreateGroundRent_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newRentTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/rents", "application/json", strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// ── GET /v1/rent/rents ────────────────────────────────────────────────────────

func TestListGroundRents_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newRentTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/rents")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.GroundRent `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}
