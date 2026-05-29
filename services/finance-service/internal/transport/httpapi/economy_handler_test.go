package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// Marx's §I two-enterprise illustration: 11,000c + 1,000v + 1,000s economising
// £2,000 of constant capital lifts p' from 8 1/3% (833 bp) to 10% (1000 bp), a
// gain of 167 bp.
func TestCreateEconomyAnalysis_Marx(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"kind":"raw_material_saving","constant":11000,"variable":1000,"surplus_value":1000,"saving":2000}`
	res, err := http.Post(ts.URL+"/v1/profit/economy", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/profit/economy/") {
		t.Errorf("Location = %q, want /v1/profit/economy/ prefix", loc)
	}
	var a profit.EconomyAnalysis
	if err := json.NewDecoder(res.Body).Decode(&a); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if a.OldProfitRate != 833 || a.NewProfitRate != 1000 || a.ProfitRateGain != 167 {
		t.Errorf("got old=%d new=%d gain=%d, want 833/1000/167",
			a.OldProfitRate, a.NewProfitRate, a.ProfitRateGain)
	}
	if a.Economy.Kind != profit.KindRawMaterialSaving {
		t.Errorf("kind = %q, want raw_material_saving", a.Economy.Kind)
	}
	if a.ID.IsZero() {
		t.Error("expected an assigned id")
	}
}

func TestCreateEconomyAnalysis_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/economy", "application/json", strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

func TestCreateEconomyAnalysis_UnknownKind(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"kind":"bogus","constant":11000,"variable":1000,"surplus_value":1000,"saving":2000}`
	res, err := http.Post(ts.URL+"/v1/profit/economy", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestCreateEconomyAnalysis_Negative(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"kind":"waste_reduction","constant":2000,"variable":1000,"surplus_value":1000,"saving":-1}`
	res, err := http.Post(ts.URL+"/v1/profit/economy", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestCreateEconomyAnalysis_SavingExceedsConstant(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"kind":"raw_material_saving","constant":100,"variable":50,"surplus_value":50,"saving":200}`
	res, err := http.Post(ts.URL+"/v1/profit/economy", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestGetEconomyAnalysis_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/economy/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

func TestCreateThenGetEconomyAnalysis(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	// A spinning mill reclaiming waste cotton: 2000c + 1000v + 1000s, £200 saved.
	res, err := http.Post(ts.URL+"/v1/profit/economy", "application/json",
		strings.NewReader(`{"kind":"waste_reduction","constant":2000,"variable":1000,"surplus_value":1000,"saving":200}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created profit.EconomyAnalysis
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if created.OldProfitRate != 3333 || created.NewProfitRate != 3571 {
		t.Errorf("got old=%d new=%d, want 3333/3571", created.OldProfitRate, created.NewProfitRate)
	}

	got, err := http.Get(ts.URL + "/v1/profit/economy/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched profit.EconomyAnalysis
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.NewProfitRate != created.NewProfitRate {
		t.Errorf("fetched = %+v, want id %s", fetched, created.ID)
	}
}

func TestListEconomyAnalyses_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/economy")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []profit.EconomyAnalysis `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items should be a non-nil array, not null")
	}
}
