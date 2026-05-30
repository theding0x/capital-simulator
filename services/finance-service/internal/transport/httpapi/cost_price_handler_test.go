package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/store"
)

func newFinanceTestServer(t *testing.T) (*httptest.Server, *store.Memory) {
	t.Helper()
	st := store.NewMemory()
	h := New(st, slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/profit/cost-price", h.CreateCostPrice)
	mux.HandleFunc("GET /v1/profit/cost-price", h.ListCostPrices)
	mux.HandleFunc("GET /v1/profit/cost-price/{id}", h.GetCostPrice)
	mux.HandleFunc("POST /v1/profit/profit-form", h.ComputeProfitForm)
	mux.HandleFunc("POST /v1/profit/rate", h.CreateProfitRate)
	mux.HandleFunc("GET /v1/profit/rate", h.ListProfitRates)
	mux.HandleFunc("GET /v1/profit/rate/{id}", h.GetProfitRate)
	mux.HandleFunc("POST /v1/profit/variation", h.CreateVariation)
	mux.HandleFunc("GET /v1/profit/variation", h.ListVariations)
	mux.HandleFunc("GET /v1/profit/variation/{id}", h.GetVariation)
	mux.HandleFunc("POST /v1/profit/compare", h.CompareProfitRates)
	mux.HandleFunc("POST /v1/profit/turnover-analysis", h.CreateTurnoverAnalysis)
	mux.HandleFunc("GET /v1/profit/turnover-analysis", h.ListTurnoverAnalyses)
	mux.HandleFunc("GET /v1/profit/turnover-analysis/{id}", h.GetTurnoverAnalysis)
	mux.HandleFunc("POST /v1/profit/economy", h.CreateEconomyAnalysis)
	mux.HandleFunc("GET /v1/profit/economy", h.ListEconomyAnalyses)
	mux.HandleFunc("GET /v1/profit/economy/{id}", h.GetEconomyAnalysis)
	mux.HandleFunc("POST /v1/profit/price-fluctuation", h.CreatePriceFluctuationAnalysis)
	mux.HandleFunc("GET /v1/profit/price-fluctuation", h.ListPriceFluctuationAnalyses)
	mux.HandleFunc("GET /v1/profit/price-fluctuation/{id}", h.GetPriceFluctuationAnalysis)
	mux.HandleFunc("POST /v1/profit/composition-effect", h.CreateCompositionEffect)
	mux.HandleFunc("GET /v1/profit/composition-effect", h.ListCompositionEffects)
	mux.HandleFunc("GET /v1/profit/composition-effect/{id}", h.GetCompositionEffect)
	mux.HandleFunc("POST /v1/profit/magnitude-change", h.CreateMagnitudeChange)
	mux.HandleFunc("GET /v1/profit/magnitude-change", h.ListMagnitudeChanges)
	mux.HandleFunc("GET /v1/profit/magnitude-change/{id}", h.GetMagnitudeChange)
	mux.HandleFunc("GET /v1/profit/summary", h.GetPartISummary)
	// Vol. III Ch. 8 — Different Compositions of Capitals in Different Branches
	mux.HandleFunc("POST /v1/avgprofit/spheres", h.CreateProductionSphere)
	mux.HandleFunc("GET /v1/avgprofit/spheres", h.ListProductionSpheres)
	mux.HandleFunc("GET /v1/avgprofit/spheres/{id}", h.GetProductionSphere)
	// Vol. III Ch. 9 — Formation of a General Rate of Profit
	mux.HandleFunc("POST /v1/avgprofit/general-rate", h.CreateGeneralProfitRate)
	mux.HandleFunc("GET /v1/avgprofit/general-rate/{id}", h.GetGeneralProfitRate)
	mux.HandleFunc("POST /v1/avgprofit/price-of-production", h.CreatePriceOfProduction)
	mux.HandleFunc("GET /v1/avgprofit/price-of-production", h.ListPricesOfProduction)
	mux.HandleFunc("GET /v1/avgprofit/price-of-production/{id}", h.GetPriceOfProduction)
	mux.HandleFunc("POST /v1/avgprofit/social-aggregate", h.ComputeSocialAggregate)
	// Vol. III Ch. 10 — Equalisation of the General Rate of Profit Through Competition
	mux.HandleFunc("POST /v1/avgprofit/market-value", h.CreateMarketValue)
	mux.HandleFunc("GET /v1/avgprofit/market-value/{id}", h.GetMarketValue)
	mux.HandleFunc("POST /v1/avgprofit/surplus-profit", h.CreateSurplusProfit)
	mux.HandleFunc("GET /v1/avgprofit/surplus-profit/{id}", h.GetSurplusProfit)
	mux.HandleFunc("POST /v1/avgprofit/capital-flow", h.CreateCapitalFlow)
	mux.HandleFunc("GET /v1/avgprofit/equalisation/{id}", h.GetEqualisation)
	// Vol. III Ch. 11 — Effects of General Wage Fluctuations on Prices of Production
	mux.HandleFunc("POST /v1/avgprofit/wage-effect", h.CreateWageEffectAnalysis)
	mux.HandleFunc("GET /v1/avgprofit/wage-effect", h.ListWageEffectAnalyses)
	mux.HandleFunc("GET /v1/avgprofit/wage-effect/{id}", h.GetWageEffectAnalysis)
	// Vol. III Ch. 12 — Supplementary Remarks (on prices of production)
	mux.HandleFunc("POST /v1/avgprofit/price-change", h.CreatePriceOfProductionChange)
	mux.HandleFunc("GET /v1/avgprofit/price-change/{id}", h.GetPriceOfProductionChange)
	mux.HandleFunc("POST /v1/avgprofit/compensation-ground", h.ComputeCompensationGround)
	mux.HandleFunc("GET /v1/avgprofit/summary", h.GetPartIISummary)
	// Vol. III Ch. 13 — The Law As Such (Tendential Fall in the Rate of Profit)
	mux.HandleFunc("POST /v1/tendency/trajectory", h.CreateCompositionTrajectory)
	mux.HandleFunc("GET /v1/tendency/trajectory", h.ListCompositionTrajectories)
	mux.HandleFunc("GET /v1/tendency/trajectory/{id}", h.GetCompositionTrajectory)
	mux.HandleFunc("POST /v1/tendency/rate-mass", h.CreateRateMassContradiction)
	mux.HandleFunc("GET /v1/tendency/rate-mass/{id}", h.GetRateMassContradiction)
	// Vol. III Ch. 14 — Counteracting Influences
	mux.HandleFunc("POST /v1/tendency/counteracting-force", h.CreateCounteractingForce)
	mux.HandleFunc("GET /v1/tendency/counteracting-force/{id}", h.GetCounteractingForce)
	mux.HandleFunc("POST /v1/tendency/scenario", h.CreateCounteractingScenario)
	mux.HandleFunc("GET /v1/tendency/scenario", h.ListCounteractingScenarios)
	mux.HandleFunc("GET /v1/tendency/scenario/{id}", h.GetCounteractingScenario)
	// Vol. III Ch. 15 — Exposition of the Internal Contradictions of the Law
	mux.HandleFunc("POST /v1/tendency/crisis", h.CreateCrisis)
	mux.HandleFunc("GET /v1/tendency/crisis", h.ListCrises)
	mux.HandleFunc("GET /v1/tendency/crisis/{id}", h.GetCrisis)
	mux.HandleFunc("POST /v1/tendency/contradiction", h.CreateInternalContradiction)
	mux.HandleFunc("GET /v1/tendency/contradiction/{id}", h.GetInternalContradiction)
	mux.HandleFunc("GET /v1/tendency/summary", h.GetPartIIISummary)
	// Vol. III Ch. 16 — Commercial Capital (opens Part IV)
	mux.HandleFunc("POST /v1/merchant/commercial-capital", h.CreateCommercialCapital)
	mux.HandleFunc("GET /v1/merchant/commercial-capital", h.ListCommercialCapitals)
	mux.HandleFunc("GET /v1/merchant/commercial-capital/{id}", h.GetCommercialCapital)
	// Vol. III Ch. 17 — Commercial Profit
	mux.HandleFunc("POST /v1/merchant/commercial-profit", h.CreateCommercialProfit)
	mux.HandleFunc("GET /v1/merchant/commercial-profit", h.ListCommercialProfits)
	mux.HandleFunc("GET /v1/merchant/commercial-profit/{id}", h.GetCommercialProfit)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts, st
}

func TestCreateCostPrice_Marx(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	body := `{"constant":400,"variable":100,"fixed_wear_and_tear":20,"fixed_advanced":1200}`
	res, err := http.Post(ts.URL+"/v1/profit/cost-price", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/profit/cost-price/") {
		t.Errorf("Location = %q, want /v1/profit/cost-price/ prefix", loc)
	}
	var cp profit.CostPrice
	if err := json.NewDecoder(res.Body).Decode(&cp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if cp.K != 500 || cp.FixedComponent != 20 || cp.CirculatingComponent != 480 {
		t.Errorf("got k=%d fixed=%d circ=%d, want 500/20/480", cp.K, cp.FixedComponent, cp.CirculatingComponent)
	}
	if cp.ID.IsZero() {
		t.Error("expected an assigned id")
	}
}

func TestCreateCostPrice_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/cost-price", "application/json", strings.NewReader("{not json"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

func TestCreateCostPrice_InvalidOutlay(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	// fixed wear-and-tear exceeds consumed constant capital → 422
	body := `{"constant":10,"variable":100,"fixed_wear_and_tear":20}`
	res, err := http.Post(ts.URL+"/v1/profit/cost-price", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

func TestGetCostPrice_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/cost-price/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

func TestCreateThenGetCostPrice(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/cost-price", "application/json",
		strings.NewReader(`{"constant":400,"variable":100}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	var created profit.CostPrice
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode: %v", err)
	}

	got, err := http.Get(ts.URL + "/v1/profit/cost-price/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}
	var fetched profit.CostPrice
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID || fetched.K != 500 {
		t.Errorf("fetched = %+v, want id %s and k 500", fetched, created.ID)
	}
}

func TestListCostPrices_NeverNull(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Get(ts.URL + "/v1/profit/cost-price")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var body struct {
		Items []profit.CostPrice `json:"items"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Items == nil {
		t.Error("items should be a non-nil array, not null")
	}
}

func TestComputeProfitForm_Marx(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)
	res, err := http.Post(ts.URL+"/v1/profit/profit-form", "application/json",
		strings.NewReader(`{"cost_price":500,"surplus_value":100}`))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got profitFormResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Profit != 100 || got.CommodityValue != 600 || !got.MystifiesOrigin {
		t.Errorf("profit-form = %+v, want profit 100, C 600, mystifies", got.ProfitForm)
	}
	if len(got.SellingPriceScenarios) != 5 {
		t.Fatalf("scenarios = %d, want 5", len(got.SellingPriceScenarios))
	}
	last := got.SellingPriceScenarios[len(got.SellingPriceScenarios)-1]
	if last.SellingPrice != 600 || last.Amount != 100 || last.BelowValue {
		t.Errorf("final scenario = %+v, want sell 600, amount 100, not below value", last)
	}
}
