package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// ── StockCompany handler tests ────────────────────────────────────────────────

func TestCreateStockCompany_UnitedAlkali(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	// United Alkali Trust: fixed=500_000_000, floating=100_000_000, total=600_000_000.
	body := `{
		"name": "United Alkali Trust",
		"fixed_capital": 500000000,
		"floating_capital": 100000000,
		"total_capital": 600000000,
		"manager_name": "Salaried Director",
		"manager_title": "Managing Director",
		"manager_wage_minutes": 480,
		"description": "Joint-stock amalgamation of alkali firms"
	}`
	res, err := http.Post(ts.URL+"/v1/credit/stock-companies", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}
	if res.Header.Get("Location") == "" {
		t.Error("expected Location header")
	}

	var sc credit.StockCompany
	if err := json.NewDecoder(res.Body).Decode(&sc); err != nil {
		t.Fatal(err)
	}
	if sc.TotalCapital != sc.FixedCapital+sc.FloatingCapital {
		t.Errorf("TotalCapital invariant: %d != %d + %d", sc.TotalCapital, sc.FixedCapital, sc.FloatingCapital)
	}
	if sc.Manager.WageMinutes != 480 {
		t.Errorf("expected manager_wage_minutes=480, got %d", sc.Manager.WageMinutes)
	}
	if sc.Manager.Name != "Salaried Director" {
		t.Errorf("expected manager_name='Salaried Director', got %q", sc.Manager.Name)
	}
}

func TestGetStockCompany_RoundTrip(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{
		"name": "Round-Trip Co",
		"fixed_capital": 200000000,
		"floating_capital": 50000000,
		"total_capital": 250000000,
		"manager_name": "Test Manager",
		"manager_title": "GM",
		"manager_wage_minutes": 120,
		"description": "round-trip test"
	}`
	res, err := http.Post(ts.URL+"/v1/credit/stock-companies", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var created credit.StockCompany
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	// GET the created entity.
	getRes, err := http.Get(ts.URL + "/v1/credit/stock-companies/" + string(created.ID))
	if err != nil {
		t.Fatal(err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRes.StatusCode)
	}

	var fetched credit.StockCompany
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatal(err)
	}
	// Verify manager fields round-trip (parity check).
	if fetched.Manager.Name != created.Manager.Name {
		t.Errorf("manager_name mismatch: got %q, want %q", fetched.Manager.Name, created.Manager.Name)
	}
	if fetched.Manager.Title != created.Manager.Title {
		t.Errorf("manager_title mismatch: got %q, want %q", fetched.Manager.Title, created.Manager.Title)
	}
	if fetched.Manager.WageMinutes != created.Manager.WageMinutes {
		t.Errorf("manager_wage_minutes mismatch: got %d, want %d", fetched.Manager.WageMinutes, created.Manager.WageMinutes)
	}
	if fetched.FixedCapital != created.FixedCapital {
		t.Errorf("fixed_capital mismatch: got %d, want %d", fetched.FixedCapital, created.FixedCapital)
	}
}

func TestGetStockCompany_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/stock-companies/doesnotexist")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.StatusCode)
	}
}

func TestCreateStockCompany_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/stock-companies", "application/json", strings.NewReader("{bad json"))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestCreateStockCompany_CapitalMismatch(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	// total_capital does not equal fixed + floating.
	body := `{
		"name": "Bad Co",
		"fixed_capital": 500000000,
		"floating_capital": 100000000,
		"total_capital": 500000000,
		"manager_name": "Dir",
		"manager_title": "MD",
		"manager_wage_minutes": 60,
		"description": ""
	}`
	res, err := http.Post(ts.URL+"/v1/credit/stock-companies", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", res.StatusCode)
	}
}

func TestCreateStockCompany_ZeroWage(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{
		"name": "Zero Wage Co",
		"fixed_capital": 500000000,
		"floating_capital": 100000000,
		"total_capital": 600000000,
		"manager_name": "Dir",
		"manager_title": "MD",
		"manager_wage_minutes": 0,
		"description": ""
	}`
	res, err := http.Post(ts.URL+"/v1/credit/stock-companies", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", res.StatusCode)
	}
}

func TestListStockCompanies_NonNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/stock-companies")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	var envelope map[string]json.RawMessage
	if err := json.NewDecoder(res.Body).Decode(&envelope); err != nil {
		t.Fatal(err)
	}
	if string(envelope["items"]) == "null" {
		t.Error("items must not be null on empty list")
	}
}

// ── CooperativeFactory handler tests ─────────────────────────────────────────

func TestCreateCooperativeFactory_Rochdale(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{
		"name": "Rochdale Pioneers",
		"worker_owned": true,
		"worker_count": 28,
		"capital_advanced": 2800,
		"description": "First modern workers cooperative (1844)"
	}`
	res, err := http.Post(ts.URL+"/v1/credit/cooperative-factories", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", res.StatusCode)
	}
	if res.Header.Get("Location") == "" {
		t.Error("expected Location header")
	}

	var cf credit.CooperativeFactory
	if err := json.NewDecoder(res.Body).Decode(&cf); err != nil {
		t.Fatal(err)
	}
	if !cf.WorkerOwned {
		t.Error("expected worker_owned=true")
	}
	if cf.WorkerCount != 28 {
		t.Errorf("expected worker_count=28, got %d", cf.WorkerCount)
	}
}

func TestCreateCooperativeFactory_NotWorkerOwned(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{
		"name": "Bad Factory",
		"worker_owned": false,
		"worker_count": 10,
		"capital_advanced": 5000,
		"description": ""
	}`
	res, err := http.Post(ts.URL+"/v1/credit/cooperative-factories", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", res.StatusCode)
	}
}

func TestCreateCooperativeFactory_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/cooperative-factories", "application/json", strings.NewReader("{broken"))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestListCooperativeFactories_NonNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/cooperative-factories")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	var envelope map[string]json.RawMessage
	if err := json.NewDecoder(res.Body).Decode(&envelope); err != nil {
		t.Fatal(err)
	}
	if string(envelope["items"]) == "null" {
		t.Error("items must not be null on empty list")
	}
}

// The Location header set by CreateCooperativeFactory must resolve: POST then
// GET the returned id and confirm the factory round-trips.
func TestGetCooperativeFactory_RoundTrip(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{
		"name": "Rochdale Pioneers",
		"worker_owned": true,
		"worker_count": 28,
		"capital_advanced": 2800,
		"description": "First modern workers cooperative (1844)"
	}`
	res, err := http.Post(ts.URL+"/v1/credit/cooperative-factories", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	location := res.Header.Get("Location")
	if location == "" {
		t.Fatal("expected Location header")
	}

	var created credit.CooperativeFactory
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	// Follow the advertised Location header directly — it must resolve.
	getRes, err := http.Get(ts.URL + location)
	if err != nil {
		t.Fatal(err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 from Location %q, got %d", location, getRes.StatusCode)
	}

	var fetched credit.CooperativeFactory
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatal(err)
	}
	if fetched.ID != created.ID {
		t.Errorf("id mismatch: got %q, want %q", fetched.ID, created.ID)
	}
	if fetched.WorkerCount != created.WorkerCount {
		t.Errorf("worker_count mismatch: got %d, want %d", fetched.WorkerCount, created.WorkerCount)
	}
}

func TestGetCooperativeFactory_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/cooperative-factories/doesnotexist")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.StatusCode)
	}
}

// ── Store (Memory) tests ──────────────────────────────────────────────────────

func TestStockCompanyStore_ErrNotFound(t *testing.T) {
	t.Parallel()
	ts, st := newFinanceTestServer(t)
	_ = ts

	_, err := st.GetStockCompany(t.Context(), credit.StockCompanyID("doesnotexist"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestStockCompanyStore_CreateGetRoundTrip(t *testing.T) {
	t.Parallel()
	_, st := newFinanceTestServer(t)

	sc, err := credit.NewStockCompany(
		"United Alkali Trust",
		500_000_000, 100_000_000, 600_000_000,
		credit.Manager{Name: "Dir", Title: "MD", WageMinutes: 480},
		"test",
	)
	if err != nil {
		t.Fatal(err)
	}

	created, err := st.CreateStockCompany(t.Context(), sc)
	if err != nil {
		t.Fatal(err)
	}
	fetched, err := st.GetStockCompany(t.Context(), created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if fetched.Manager.WageMinutes != 480 {
		t.Errorf("WageMinutes round-trip: got %d, want 480", fetched.Manager.WageMinutes)
	}
	if fetched.Manager.Name != "Dir" {
		t.Errorf("Manager.Name round-trip: got %q, want Dir", fetched.Manager.Name)
	}
	if fetched.TotalCapital != 600_000_000 {
		t.Errorf("TotalCapital round-trip: got %d", fetched.TotalCapital)
	}
}

func TestStockCompanyStore_ListNonNil(t *testing.T) {
	t.Parallel()
	_, st := newFinanceTestServer(t)

	items, err := st.ListStockCompanies(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if items == nil {
		t.Error("ListStockCompanies must return non-nil on empty store")
	}
}

func TestCooperativeFactoryStore_ErrNotFound(t *testing.T) {
	t.Parallel()
	_, st := newFinanceTestServer(t)

	_, err := st.GetCooperativeFactory(t.Context(), credit.CooperativeFactoryID("missing"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCooperativeFactoryStore_CreateGetRoundTrip(t *testing.T) {
	t.Parallel()
	_, st := newFinanceTestServer(t)

	cf, err := credit.NewCooperativeFactory("Rochdale Pioneers", true, 28, 2800, "test")
	if err != nil {
		t.Fatal(err)
	}
	created, err := st.CreateCooperativeFactory(t.Context(), cf)
	if err != nil {
		t.Fatal(err)
	}
	fetched, err := st.GetCooperativeFactory(t.Context(), created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !fetched.WorkerOwned {
		t.Error("WorkerOwned round-trip: expected true")
	}
	if fetched.WorkerCount != 28 {
		t.Errorf("WorkerCount round-trip: got %d, want 28", fetched.WorkerCount)
	}
}

func TestCooperativeFactoryStore_ListNonNil(t *testing.T) {
	t.Parallel()
	_, st := newFinanceTestServer(t)

	items, err := st.ListCooperativeFactories(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if items == nil {
		t.Error("ListCooperativeFactories must return non-nil on empty store")
	}
}
