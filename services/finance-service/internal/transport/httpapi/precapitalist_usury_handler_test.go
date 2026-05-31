package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// ── POST /v1/credit/usurers-capital ──────────────────────────────────────────

// TestCreateUsurersCapital_AncientRome verifies the Ancient Rome fixture
// (StageAntiquity, no wage labour, rate=4800 bp) → 201, Location set,
// borrowers is a JSON array of length 2.
func TestCreateUsurersCapital_AncientRome(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Ancient Roman Usury","stage":1,"wage_labour":false,"borrowers":["plebeians","nobles"],"interest_rate_bp":4800,"subordinated_to":"","description":"Usury in ancient Rome."}`
	res, err := http.Post(ts.URL+"/v1/credit/usurers-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status: got %d, want 201", res.StatusCode)
	}
	if res.Header.Get("Location") == "" {
		t.Error("expected Location header")
	}

	var u credit.UsurersCapital
	if err := json.NewDecoder(res.Body).Decode(&u); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if u.ID.IsZero() {
		t.Error("ID must be set")
	}
	if len(u.Borrowers) != 2 {
		t.Errorf("borrowers length: got %d, want 2", len(u.Borrowers))
	}
	if u.Borrowers == nil {
		t.Error("borrowers must not be nil")
	}
}

// TestCreateUsurersCapital_BadJSON tests 400 on malformed JSON.
func TestCreateUsurersCapital_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/usurers-capital", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", res.StatusCode)
	}
}

// TestCreateUsurersCapital_InvalidStage tests 422 on stage=0.
func TestCreateUsurersCapital_InvalidStage(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"bad","stage":0,"wage_labour":false,"borrowers":[],"interest_rate_bp":100,"subordinated_to":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/usurers-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestCreateUsurersCapital_AntiquityWageLabour tests 422 on
// stage=1 + wage_labour=true (pre-capitalist stage cannot have wage labour).
func TestCreateUsurersCapital_AntiquityWageLabour(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"bad","stage":1,"wage_labour":true,"borrowers":[],"interest_rate_bp":100,"subordinated_to":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/usurers-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestCreateUsurersCapital_NegativeRate tests 422 on interest_rate_bp=-1.
func TestCreateUsurersCapital_NegativeRate(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"bad","stage":1,"wage_labour":false,"borrowers":[],"interest_rate_bp":-1,"subordinated_to":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/usurers-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// TestCreateUsurersCapital_SubordinatedEmptyTo tests 422 on
// stage=4 + subordinated_to="" (StageSubordinated requires "industrial_capital").
func TestCreateUsurersCapital_SubordinatedEmptyTo(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"bad","stage":4,"wage_labour":true,"borrowers":[],"interest_rate_bp":800,"subordinated_to":"","description":""}`
	res, err := http.Post(ts.URL+"/v1/credit/usurers-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status: got %d, want 422", res.StatusCode)
	}
}

// ── GET /v1/credit/usurers-capital ───────────────────────────────────────────

// TestListUsurersCapitals_NeverNull verifies the list endpoint returns
// {"items":[]} not null when the store is empty.
func TestListUsurersCapitals_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/usurers-capital")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status: got %d, want 200", res.StatusCode)
	}

	var payload struct {
		Items []credit.UsurersCapital `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Items == nil {
		t.Error("items must not be nil")
	}
}

// TestListUsurersCapitals_PostThenList verifies that a created record appears
// in the list with a non-null borrowers array.
func TestListUsurersCapitals_PostThenList(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Ancient Roman Usury","stage":1,"wage_labour":false,"borrowers":["plebeians","nobles"],"interest_rate_bp":4800,"subordinated_to":"","description":"Rome."}`
	createRes, err := http.Post(ts.URL+"/v1/credit/usurers-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	createRes.Body.Close()
	if createRes.StatusCode != http.StatusCreated {
		t.Fatalf("create status: got %d, want 201", createRes.StatusCode)
	}

	listRes, err := http.Get(ts.URL + "/v1/credit/usurers-capital")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	defer listRes.Body.Close()
	if listRes.StatusCode != http.StatusOK {
		t.Fatalf("list status: got %d, want 200", listRes.StatusCode)
	}

	var payload struct {
		Items []credit.UsurersCapital `json:"items"`
	}
	if err := json.NewDecoder(listRes.Body).Decode(&payload); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(payload.Items) != 1 {
		t.Fatalf("items length: got %d, want 1", len(payload.Items))
	}
	if payload.Items[0].Borrowers == nil {
		t.Error("borrowers in list item must not be nil")
	}
}

// ── GET /v1/credit/usurers-capital/{id} ──────────────────────────────────────

// TestGetUsurersCapital_PostThenGetByID verifies that a created record can be
// retrieved via its ID with all fields intact.
func TestGetUsurersCapital_PostThenGetByID(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"name":"Medieval Church Usury","stage":2,"wage_labour":false,"borrowers":["peasants","feudal lords"],"interest_rate_bp":2400,"subordinated_to":"","description":"Medieval usury."}`
	createRes, err := http.Post(ts.URL+"/v1/credit/usurers-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer createRes.Body.Close()
	if createRes.StatusCode != http.StatusCreated {
		t.Fatalf("create status: got %d, want 201", createRes.StatusCode)
	}

	var created credit.UsurersCapital
	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	getRes, err := http.Get(ts.URL + "/v1/credit/usurers-capital/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("get status: got %d, want 200", getRes.StatusCode)
	}

	var fetched credit.UsurersCapital
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("ID mismatch: got %s, want %s", fetched.ID, created.ID)
	}
	if fetched.Stage != created.Stage {
		t.Errorf("Stage mismatch: got %d, want %d", fetched.Stage, created.Stage)
	}
	if fetched.InterestRateBP != 2400 {
		t.Errorf("InterestRateBP: got %d, want 2400", fetched.InterestRateBP)
	}
}

// TestGetUsurersCapital_NotFound tests 404 for an unknown ID.
func TestGetUsurersCapital_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/usurers-capital/doesnotexist")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", res.StatusCode)
	}
}
