package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// TestGetPartIISummary_Empty returns 200 with a non-nil body even when the store
// holds no records: with no prices, Σ values == Σ prices == 0 so value is
// conserved.
func TestGetPartIISummary_Empty(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/avgprofit/summary")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var summary avgprofit.PartIISummary
	if err := json.NewDecoder(res.Body).Decode(&summary); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !summary.ValueConserved {
		t.Error("ValueConserved = false, want true (0 == 0)")
	}
	if summary.PricesOfProductionCount != 0 {
		t.Errorf("PricesOfProductionCount = %d, want 0", summary.PricesOfProductionCount)
	}
}

// TestGetPartIISummary_Conserved seeds the in-memory store with the five Ch. 9
// prices of production (each k=100, rate=2200, price 122; values 120/130/140/
// 115/105 → Σ values = Σ prices = 610) and asserts the summary reports value
// conserved with the general rate 2200 and five prices counted.
func TestGetPartIISummary_Conserved(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	bodies := []string{
		`{"sphere_name":"Sphere I (80c+20v)","cost_price":100,"general_rate":2200,"commodity_value":120}`,
		`{"sphere_name":"Sphere II (70c+30v)","cost_price":100,"general_rate":2200,"commodity_value":130}`,
		`{"sphere_name":"Sphere III (60c+40v)","cost_price":100,"general_rate":2200,"commodity_value":140}`,
		`{"sphere_name":"Sphere IV (85c+15v)","cost_price":100,"general_rate":2200,"commodity_value":115}`,
		`{"sphere_name":"Sphere V (95c+5v)","cost_price":100,"general_rate":2200,"commodity_value":105}`,
	}
	for i, b := range bodies {
		res, err := http.Post(ts.URL+"/v1/avgprofit/price-of-production", "application/json", strings.NewReader(b))
		if err != nil {
			t.Fatalf("seed price %d: %v", i, err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusCreated {
			t.Fatalf("seed price %d status = %d, want 201", i, res.StatusCode)
		}
	}

	res, err := http.Get(ts.URL + "/v1/avgprofit/summary")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var summary avgprofit.PartIISummary
	if err := json.NewDecoder(res.Body).Decode(&summary); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if summary.PricesOfProductionCount != 5 {
		t.Errorf("PricesOfProductionCount = %d, want 5", summary.PricesOfProductionCount)
	}
	if summary.TotalValues != 610 || summary.TotalPricesOfProduction != 610 {
		t.Errorf("totals = values %d prices %d, want 610/610",
			summary.TotalValues, summary.TotalPricesOfProduction)
	}
	if !summary.ValueConserved {
		t.Error("ValueConserved = false, want true")
	}
	if summary.GeneralRate != 2200 {
		t.Errorf("GeneralRate = %d, want 2200", summary.GeneralRate)
	}
}
