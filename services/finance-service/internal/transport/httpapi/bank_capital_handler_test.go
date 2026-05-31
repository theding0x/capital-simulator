package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// ── BankCapital handler tests ─────────────────────────────────────────────────

func TestCreateBankCapital_Scotland(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{
		"name": "Scottish Bank",
		"cash_amount": 300000000,
		"securities_amount": 2400000000,
		"components": [
			{"amount": 300000000, "kind": 1, "description": "gold and notes"},
			{"amount": 1400000000, "kind": 2, "description": "discounted bills"},
			{"amount": 1000000000, "kind": 3, "description": "consols"}
		],
		"description": "cash a small fraction of the whole"
	}`
	res, err := http.Post(ts.URL+"/v1/credit/bank-capital", "application/json", strings.NewReader(body))
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

	var bc credit.BankCapital
	if err := json.NewDecoder(res.Body).Decode(&bc); err != nil {
		t.Fatal(err)
	}
	if bc.TotalCapital != 2700000000 {
		t.Errorf("TotalCapital=%d, want 2700000000", bc.TotalCapital)
	}
	if bc.TotalCapital != bc.CashAmount+bc.SecuritiesAmount {
		t.Errorf("invariant violated: %d != %d + %d", bc.TotalCapital, bc.CashAmount, bc.SecuritiesAmount)
	}
	if len(bc.Components) != 3 {
		t.Errorf("Components length=%d, want 3", len(bc.Components))
	}
}

// TestGetBankCapital_ComponentRoundTrip proves the components_json column
// round-trips a component through create → get with full fidelity.
func TestGetBankCapital_ComponentRoundTrip(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{
		"name": "Round-Trip Bank",
		"cash_amount": 100,
		"securities_amount": 900,
		"components": [
			{"amount": 555, "kind": 4, "description": "railway shares"}
		],
		"description": "rt"
	}`
	res, err := http.Post(ts.URL+"/v1/credit/bank-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var created credit.BankCapital
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	getRes, err := http.Get(ts.URL + "/v1/credit/bank-capital/" + string(created.ID))
	if err != nil {
		t.Fatal(err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRes.StatusCode)
	}

	var fetched credit.BankCapital
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatal(err)
	}
	if len(fetched.Components) != 1 {
		t.Fatalf("Components length=%d, want 1", len(fetched.Components))
	}
	c := fetched.Components[0]
	if c.Amount != 555 || c.Kind != credit.KindStock || c.Description != "railway shares" {
		t.Errorf("component round-trip mismatch: %+v", c)
	}
	if fetched.TotalCapital != 1000 {
		t.Errorf("TotalCapital round-trip: got %d, want 1000", fetched.TotalCapital)
	}
}

func TestGetBankCapital_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/bank-capital/doesnotexist")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", res.StatusCode)
	}
}

func TestCreateBankCapital_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/bank-capital", "application/json", strings.NewReader("{bad json"))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestCreateBankCapital_InvalidComponentKind(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{
		"name": "Bad Bank",
		"cash_amount": 100,
		"securities_amount": 0,
		"components": [{"amount": 10, "kind": 6, "description": "bad"}],
		"description": ""
	}`
	res, err := http.Post(ts.URL+"/v1/credit/bank-capital", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", res.StatusCode)
	}
}

func TestListBankCapitals_NonNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/bank-capital")
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

// ── FictitiousCapitalValuation handler tests ─────────────────────────────────

func TestCreateFictitiousCapitalValuation_Bond(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{
		"name": "Government Bond @5%",
		"nominal_value": 10000,
		"annual_income": 500,
		"interest_rate_bp": 500,
		"dividend_rate_bp": 0,
		"description": "consol"
	}`
	res, err := http.Post(ts.URL+"/v1/credit/fictitious-capital-valuation", "application/json", strings.NewReader(body))
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

	var v credit.FictitiousCapitalValuation
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		t.Fatal(err)
	}
	if v.MarketValue != 10000 {
		t.Errorf("MarketValue=%d, want 10000", v.MarketValue)
	}
}

func TestGetFictitiousCapitalValuation_RoundTrip(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{
		"name": "Joint-Stock Share",
		"nominal_value": 10000,
		"annual_income": 1000,
		"interest_rate_bp": 500,
		"dividend_rate_bp": 1000,
		"description": "share"
	}`
	res, err := http.Post(ts.URL+"/v1/credit/fictitious-capital-valuation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	var created credit.FictitiousCapitalValuation
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	getRes, err := http.Get(ts.URL + "/v1/credit/fictitious-capital-valuation/" + string(created.ID))
	if err != nil {
		t.Fatal(err)
	}
	defer getRes.Body.Close()
	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRes.StatusCode)
	}

	var fetched credit.FictitiousCapitalValuation
	if err := json.NewDecoder(getRes.Body).Decode(&fetched); err != nil {
		t.Fatal(err)
	}
	if fetched.MarketValue != 20000 {
		t.Errorf("MarketValue round-trip: got %d, want 20000", fetched.MarketValue)
	}
	if fetched.DividendRateBP != 1000 {
		t.Errorf("DividendRateBP round-trip: got %d, want 1000", fetched.DividendRateBP)
	}
}

func TestCreateFictitiousCapitalValuation_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/credit/fictitious-capital-valuation", "application/json", strings.NewReader("{broken"))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

func TestCreateFictitiousCapitalValuation_NonPositiveRate(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{
		"name": "Bad Rate",
		"nominal_value": 10000,
		"annual_income": 500,
		"interest_rate_bp": 0,
		"dividend_rate_bp": 0,
		"description": ""
	}`
	res, err := http.Post(ts.URL+"/v1/credit/fictitious-capital-valuation", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", res.StatusCode)
	}
}

func TestListFictitiousCapitalValuations_NonNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/credit/fictitious-capital-valuation")
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

// ── Store (Memory) tests ──────────────────────────────────────────────────────

func TestBankCapitalStore_ErrNotFound(t *testing.T) {
	t.Parallel()
	_, st := newFinanceTestServer(t)

	_, err := st.GetBankCapital(t.Context(), credit.BankCapitalID("missing"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestBankCapitalStore_CreateGetRoundTrip(t *testing.T) {
	t.Parallel()
	_, st := newFinanceTestServer(t)

	bc, err := credit.NewBankCapital(
		"Scottish Bank",
		300_000_000, 2_400_000_000,
		[]credit.BankCapitalComponent{
			{Amount: 1_000_000_000, Kind: credit.KindGovernmentBond, Description: "consols"},
		},
		"test",
	)
	if err != nil {
		t.Fatal(err)
	}
	created, err := st.CreateBankCapital(t.Context(), bc)
	if err != nil {
		t.Fatal(err)
	}
	fetched, err := st.GetBankCapital(t.Context(), created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if fetched.TotalCapital != 2_700_000_000 {
		t.Errorf("TotalCapital round-trip: got %d, want 2700000000", fetched.TotalCapital)
	}
	if len(fetched.Components) != 1 || fetched.Components[0].Kind != credit.KindGovernmentBond {
		t.Errorf("Components round-trip mismatch: %+v", fetched.Components)
	}
}

func TestBankCapitalStore_ListNonNil(t *testing.T) {
	t.Parallel()
	_, st := newFinanceTestServer(t)

	items, err := st.ListBankCapitals(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if items == nil {
		t.Error("ListBankCapitals must return non-nil on empty store")
	}
}

func TestFictitiousValuationStore_ErrNotFound(t *testing.T) {
	t.Parallel()
	_, st := newFinanceTestServer(t)

	_, err := st.GetFictitiousCapitalValuation(t.Context(), credit.FictitiousCapitalValuationID("missing"))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFictitiousValuationStore_CreateGetRoundTrip(t *testing.T) {
	t.Parallel()
	_, st := newFinanceTestServer(t)

	v, err := credit.NewFictitiousCapitalValuation("Bond", 10000, 500, 500, 0, "test")
	if err != nil {
		t.Fatal(err)
	}
	created, err := st.CreateFictitiousCapitalValuation(t.Context(), v)
	if err != nil {
		t.Fatal(err)
	}
	fetched, err := st.GetFictitiousCapitalValuation(t.Context(), created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if fetched.MarketValue != 10000 {
		t.Errorf("MarketValue round-trip: got %d, want 10000", fetched.MarketValue)
	}
}

func TestFictitiousValuationStore_ListNonNil(t *testing.T) {
	t.Parallel()
	_, st := newFinanceTestServer(t)

	items, err := st.ListFictitiousCapitalValuations(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if items == nil {
		t.Error("ListFictitiousCapitalValuations must return non-nil on empty store")
	}
}
