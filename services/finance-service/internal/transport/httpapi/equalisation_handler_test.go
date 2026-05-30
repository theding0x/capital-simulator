package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
)

// TestCreateCapitalFlow_CottonInflow verifies 201 + Location for the cotton spec
// fixture: currentSphereRate=3000, generalRate=2200 → direction=inflow, IsConverging=true.
// The persisted Equalisation is then retrievable via GET /equalisation/{id}.
func TestCreateCapitalFlow_CottonInflow(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"cotton","current_sphere_rate":3000,"general_rate":2200,"market_price":0,"market_value":0}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/capital-flow", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	if loc := res.Header.Get("Location"); !strings.HasPrefix(loc, "/v1/avgprofit/equalisation/") {
		t.Errorf("Location = %q, want /v1/avgprofit/equalisation/ prefix", loc)
	}

	var eq avgprofit.Equalisation
	if err := json.NewDecoder(res.Body).Decode(&eq); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if eq.ID.IsZero() {
		t.Error("expected an assigned id")
	}
	if eq.Direction != avgprofit.KindInflow {
		t.Errorf("Direction = %q, want %q (3000 > 2200)", eq.Direction, avgprofit.KindInflow)
	}
	if !eq.IsConverging {
		t.Error("IsConverging = false, want true (3000 != 2200)")
	}
	if eq.InitialRate != 3000 {
		t.Errorf("InitialRate = %d, want 3000", eq.InitialRate)
	}
	if eq.TargetRate != 2200 {
		t.Errorf("TargetRate = %d, want 2200", eq.TargetRate)
	}
}

// TestCreateCapitalFlow_PersistAndGet verifies that the Equalisation created by
// POST /capital-flow is retrievable via GET /equalisation/{id}.
func TestCreateCapitalFlow_PersistAndGet(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"cotton","current_sphere_rate":3000,"general_rate":2200,"market_price":0,"market_value":0}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/capital-flow", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	var created avgprofit.Equalisation
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("decode create: %v", err)
	}

	got, err := http.Get(ts.URL + "/v1/avgprofit/equalisation/" + string(created.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	if got.StatusCode != http.StatusOK {
		t.Fatalf("get status = %d, want 200", got.StatusCode)
	}

	var fetched avgprofit.Equalisation
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("fetched.ID = %s, want %s", fetched.ID, created.ID)
	}
	if fetched.Direction != avgprofit.KindInflow {
		t.Errorf("fetched.Direction = %q, want %q", fetched.Direction, avgprofit.KindInflow)
	}
}

// TestCreateCapitalFlow_BadJSON returns 400 on malformed body.
func TestCreateCapitalFlow_BadJSON(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Post(ts.URL+"/v1/avgprofit/capital-flow", "application/json", strings.NewReader("{bad"))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", res.StatusCode)
	}
}

// TestCreateCapitalFlow_NegativeInitialRate returns 422 on negative current_sphere_rate.
func TestCreateCapitalFlow_NegativeInitialRate(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"cotton","current_sphere_rate":-1,"general_rate":2200,"market_price":0,"market_value":0}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/capital-flow", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestCreateCapitalFlow_NegativeGeneralRate returns 422 on negative general_rate.
func TestCreateCapitalFlow_NegativeGeneralRate(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"cotton","current_sphere_rate":3000,"general_rate":-1,"market_price":0,"market_value":0}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/capital-flow", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want 422", res.StatusCode)
	}
}

// TestGetEqualisation_NotFound returns 404 for an unknown id.
func TestGetEqualisation_NotFound(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	res, err := http.Get(ts.URL + "/v1/avgprofit/equalisation/nonexistent")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Errorf("status = %d, want 404", res.StatusCode)
	}
}

// TestCreateCapitalFlow_InflowForCurrentAboveGeneral verifies the invariant:
// Direction == KindInflow ⟺ CurrentSphereRate > GeneralRate (rate-based path).
func TestCreateCapitalFlow_InflowForCurrentAboveGeneral(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"iron","current_sphere_rate":4000,"general_rate":2200,"market_price":0,"market_value":0}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/capital-flow", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()

	var eq avgprofit.Equalisation
	if err := json.NewDecoder(res.Body).Decode(&eq); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if eq.Direction != avgprofit.KindInflow {
		t.Errorf("Direction = %q, want %q (4000 > 2200)", eq.Direction, avgprofit.KindInflow)
	}
}

// TestCreateCapitalFlow_PricePathDirectionConsistent guards a review finding:
// when general_rate is 0 the direction comes from the market-price-vs-value path,
// and the persisted Equalisation.Direction must equal Flow.Direction and survive a
// round-trip through GET (price 95 < value 100 ⇒ outflow).
func TestCreateCapitalFlow_PricePathDirectionConsistent(t *testing.T) {
	t.Parallel()
	ts, _ := newFinanceTestServer(t)

	body := `{"sphere_name":"linen","current_sphere_rate":0,"general_rate":0,"market_price":95,"market_value":100}`
	res, err := http.Post(ts.URL+"/v1/avgprofit/capital-flow", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("post: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}

	var eq avgprofit.Equalisation
	if err := json.NewDecoder(res.Body).Decode(&eq); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if eq.Direction != avgprofit.KindOutflow {
		t.Errorf("Direction = %q, want %q (price 95 < value 100)", eq.Direction, avgprofit.KindOutflow)
	}
	if eq.Flow.Direction != eq.Direction {
		t.Errorf("Flow.Direction %q != Direction %q — must be consistent", eq.Flow.Direction, eq.Direction)
	}

	got, err := http.Get(ts.URL + "/v1/avgprofit/equalisation/" + string(eq.ID))
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer got.Body.Close()
	var fetched avgprofit.Equalisation
	if err := json.NewDecoder(got.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get: %v", err)
	}
	if fetched.Direction != avgprofit.KindOutflow || fetched.Flow.Direction != avgprofit.KindOutflow {
		t.Errorf("after GET: Direction=%q Flow.Direction=%q, want both outflow", fetched.Direction, fetched.Flow.Direction)
	}
}
