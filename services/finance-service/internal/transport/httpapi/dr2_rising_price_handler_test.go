// Package httpapi — handler tests for Vol. III Ch. 43 DR II Rising Price
// endpoints. Fixtures drawn from Engels Table XI: base_rent=0, increment=12,
// num_grades=5 → sequence=[0,12,24,36,48]; total_rent=120, total_capital=250,
// total_output=70 bushels, price_per_bushel=6 shillings (variant 3, inferior
// soil 'a' regulates with constant productivity).
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

// newDR2RisingTestServer builds a minimal httptest.Server that registers
// only the Vol. III Ch. 43 DR II rising-price routes.
func newDR2RisingTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/rent/rising-price-outcomes", h.CreateRisingPriceDR2Outcome)
	mux.HandleFunc("GET /v1/rent/rising-price-outcomes", h.ListRisingPriceDR2Outcomes)
	mux.HandleFunc("POST /v1/rent/rent-sequences", h.CreateRentSequence)
	mux.HandleFunc("GET /v1/rent/rent-sequences", h.ListRentSequences)
	mux.HandleFunc("POST /v1/rent/engels-tables", h.CreateEngelsTableEntry)
	mux.HandleFunc("GET /v1/rent/engels-tables", h.ListEngelsTableEntries)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

// ── POST /v1/rent/rising-price-outcomes — happy path ─────────────────────────

// TestCreateRisingPriceDR2Outcome_HappyPath uses Engels Ch. 43 variant 3
// (inferior soil 'a' regulates, constant productivity): total_capital=250s,
// total_rent=120s, total_output=70 bushels, price_per_bushel=6s.
func TestCreateRisingPriceDR2Outcome_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	body := `{"variant":3,"total_capital_shillings":250,"total_rent_shillings":120,"total_output_bushels":70,"price_per_bushel_shillings":6}`
	res, err := http.Post(ts.URL+"/v1/rent/rising-price-outcomes", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/rent/rising-price-outcomes") {
		t.Errorf("Location = %q, want prefix /v1/rent/rising-price-outcomes", loc)
	}

	var saved rent.RisingPriceDR2Outcome
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
}

// ── POST /v1/rent/rising-price-outcomes — validation failures ────────────────

func TestCreateRisingPriceDR2Outcome_VariantZero_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	body := `{"variant":0,"total_capital_shillings":250,"total_rent_shillings":120,"total_output_bushels":70,"price_per_bushel_shillings":6}`
	res, err := http.Post(ts.URL+"/v1/rent/rising-price-outcomes", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for variant=0", res.StatusCode)
	}
}

func TestCreateRisingPriceDR2Outcome_VariantSix_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	body := `{"variant":6,"total_capital_shillings":250,"total_rent_shillings":120,"total_output_bushels":70,"price_per_bushel_shillings":6}`
	res, err := http.Post(ts.URL+"/v1/rent/rising-price-outcomes", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for variant=6", res.StatusCode)
	}
}

func TestCreateRisingPriceDR2Outcome_ZeroCapital_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	body := `{"variant":3,"total_capital_shillings":0,"total_rent_shillings":120,"total_output_bushels":70,"price_per_bushel_shillings":6}`
	res, err := http.Post(ts.URL+"/v1/rent/rising-price-outcomes", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for total_capital_shillings=0", res.StatusCode)
	}
}

func TestCreateRisingPriceDR2Outcome_ZeroPrice_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	body := `{"variant":3,"total_capital_shillings":250,"total_rent_shillings":120,"total_output_bushels":70,"price_per_bushel_shillings":0}`
	res, err := http.Post(ts.URL+"/v1/rent/rising-price-outcomes", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for price_per_bushel_shillings=0", res.StatusCode)
	}
}

func TestCreateRisingPriceDR2Outcome_MalformedJSON_Returns400(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	res, err := http.Post(ts.URL+"/v1/rent/rising-price-outcomes", "application/json",
		strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400 for malformed JSON", res.StatusCode)
	}
}

// ── GET /v1/rent/rising-price-outcomes ───────────────────────────────────────

func TestListRisingPriceDR2Outcomes_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/rising-price-outcomes")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.RisingPriceDR2Outcome `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

// ── POST /v1/rent/rent-sequences — happy path ────────────────────────────────

// TestCreateRentSequence_HappyPath uses Engels Table XI:
// base_rent=0, increment=12, num_grades=5 → sequence=[0,12,24,36,48], sum=120.
func TestCreateRentSequence_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	body := `{"base_rent":0,"increment":12,"num_grades":5}`
	res, err := http.Post(ts.URL+"/v1/rent/rent-sequences", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var resp struct {
		rent.RentSequence
		Sequence []int64 `json:"sequence"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID == "" {
		t.Error("response id must be non-empty")
	}

	wantSeq := []int64{0, 12, 24, 36, 48}
	if len(resp.Sequence) != len(wantSeq) {
		t.Fatalf("sequence length = %d, want %d", len(resp.Sequence), len(wantSeq))
	}
	for i, v := range wantSeq {
		if resp.Sequence[i] != v {
			t.Errorf("sequence[%d] = %d, want %d", i, resp.Sequence[i], v)
		}
	}
}

// ── POST /v1/rent/rent-sequences — validation failures ───────────────────────

func TestCreateRentSequence_ZeroNumGrades_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	body := `{"base_rent":0,"increment":12,"num_grades":0}`
	res, err := http.Post(ts.URL+"/v1/rent/rent-sequences", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for num_grades=0", res.StatusCode)
	}
}

func TestCreateRentSequence_NegativeIncrement_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	body := `{"base_rent":0,"increment":-1,"num_grades":5}`
	res, err := http.Post(ts.URL+"/v1/rent/rent-sequences", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for increment=-1", res.StatusCode)
	}
}

// ── GET /v1/rent/rent-sequences ──────────────────────────────────────────────

func TestListRentSequences_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/rent-sequences")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.RentSequence `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}

// ── POST /v1/rent/engels-tables — happy path ─────────────────────────────────

// TestCreateEngelsTableEntry_HappyPath uses Engels Table XI row A:
// table_number=11, soil_grade_label="A", output_bushels=10, rent_shillings=0,
// capital_shillings=50.
func TestCreateEngelsTableEntry_HappyPath(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	body := `{"table_number":11,"soil_grade_label":"A","output_bushels":10,"rent_shillings":0,"capital_shillings":50}`
	res, err := http.Post(ts.URL+"/v1/rent/engels-tables", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var saved rent.EngelsTableEntry
	if err := json.NewDecoder(res.Body).Decode(&saved); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if saved.ID == "" {
		t.Error("response id must be non-empty")
	}
}

// ── POST /v1/rent/engels-tables — validation failures ────────────────────────

func TestCreateEngelsTableEntry_EmptySoilGradeLabel_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	body := `{"table_number":11,"soil_grade_label":"","output_bushels":10,"rent_shillings":0,"capital_shillings":50}`
	res, err := http.Post(ts.URL+"/v1/rent/engels-tables", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for empty soil_grade_label", res.StatusCode)
	}
}

func TestCreateEngelsTableEntry_ZeroOutputBushels_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	body := `{"table_number":11,"soil_grade_label":"A","output_bushels":0,"rent_shillings":0,"capital_shillings":50}`
	res, err := http.Post(ts.URL+"/v1/rent/engels-tables", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for output_bushels=0", res.StatusCode)
	}
}

func TestCreateEngelsTableEntry_ZeroTableNumber_Returns422(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	body := `{"table_number":0,"soil_grade_label":"A","output_bushels":10,"rent_shillings":0,"capital_shillings":50}`
	res, err := http.Post(ts.URL+"/v1/rent/engels-tables", "application/json",
		strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422 for table_number=0", res.StatusCode)
	}
}

// ── GET /v1/rent/engels-tables ───────────────────────────────────────────────

func TestListEngelsTableEntries_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newDR2RisingTestServer(t)

	res, err := http.Get(ts.URL + "/v1/rent/engels-tables")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []rent.EngelsTableEntry `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items must be a non-nil array, not null")
	}
}
