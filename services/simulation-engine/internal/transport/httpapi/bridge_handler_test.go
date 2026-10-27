package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// stubFetcher returns a canned factor for the requested source. Used to
// keep the bridge handler test hermetic — no HTTP calls to agent-service.
type stubFetcher struct {
	factors map[ProductivitySource]float64
}

func (s stubFetcher) Fetch(_ context.Context, source ProductivitySource, _ string) (float64, error) {
	v, ok := s.factors[source]
	if !ok {
		return 0, errors.New("stub: no factor")
	}
	return v, nil
}

func newBridgeTestServer(t *testing.T, fetcher ProductivityFetcher) *httptest.Server {
	t.Helper()
	h := New(nil, Deps{Productivity: fetcher})
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/production/relative-surplus", h.RelativeSurplus)
	mux.HandleFunc("POST /v1/production/relative-surplus-from-productivity", h.RelativeSurplusFromProductivity)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// Capital Vol. I, Ch. 12 §1: a productivity gain that halves the value of
// labour-power frees half the working day for surplus labour. Working day
// = 720 min; old LPV = 600 (10 h necessary, 2 h surplus); productivity
// factor = 2.0 → new LPV = 300; surplus labour grows from 120 → 420.
func TestRelativeSurplus_Para1_HalveLPV(t *testing.T) {
	t.Parallel()
	ts := newBridgeTestServer(t, nil)
	body := `{"working_day_total":720,"current_lpv":600,"productivity_factor":2.0,"scope":"wage_goods_branch"}`
	res, err := http.Post(ts.URL+"/v1/production/relative-surplus", "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got relativeSurplusResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.NewLabourPowerValue != 300 {
		t.Fatalf("NewLabourPowerValue = %d, want 300", got.NewLabourPowerValue)
	}
	if got.NewWorkingDay.SurplusLabour != 420 {
		t.Fatalf("new SurplusLabour = %d, want 420", got.NewWorkingDay.SurplusLabour)
	}
	if got.RelativeSurplusValue != 300 {
		t.Fatalf("RelativeSurplusValue = %d, want 300", got.RelativeSurplusValue)
	}
}

func TestRelativeSurplus_RejectsInvalidFactor(t *testing.T) {
	t.Parallel()
	ts := newBridgeTestServer(t, nil)
	body := `{"working_day_total":720,"current_lpv":600,"productivity_factor":0}`
	res, err := http.Post(ts.URL+"/v1/production/relative-surplus", "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}

// End-to-end: a Ch. 13 cooperation factor drives a Ch. 12 calculation.
// This is the acceptance-criterion test for the Part IV bridge.
func TestRelativeSurplusFromProductivity_CooperationDrivesCh12(t *testing.T) {
	t.Parallel()
	// A cooperation with a productivity factor of 1.2 (typical
	// CollectiveProductivePower for n=12).
	fetcher := stubFetcher{factors: map[ProductivitySource]float64{
		SourceCooperation: 1.2,
	}}
	ts := newBridgeTestServer(t, fetcher)
	body := `{"working_day_total":720,"current_lpv":600,"source":"cooperation","source_id":"abc","scope":"wage_goods_branch"}`
	res, err := http.Post(ts.URL+"/v1/production/relative-surplus-from-productivity", "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got relativeSurplusResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ProductivityFactor != 1.2 {
		t.Fatalf("ProductivityFactor = %v, want 1.2", got.ProductivityFactor)
	}
	// 600 / 1.2 = 500; surplus labour grows from 120 to 220.
	if got.NewLabourPowerValue != 500 {
		t.Fatalf("NewLabourPowerValue = %d, want 500", got.NewLabourPowerValue)
	}
	if got.RelativeSurplusValue != 100 {
		t.Fatalf("RelativeSurplusValue = %d, want 100", got.RelativeSurplusValue)
	}
	if got.Source == nil || got.Source.Source != SourceCooperation {
		t.Fatalf("Source = %+v, want cooperation echoed back", got.Source)
	}
}

func TestRelativeSurplusFromProductivity_RejectsUnknownSource(t *testing.T) {
	t.Parallel()
	fetcher := stubFetcher{factors: map[ProductivitySource]float64{}}
	ts := newBridgeTestServer(t, fetcher)
	body := `{"working_day_total":720,"current_lpv":600,"source":"galaxy_brain","source_id":"x"}`
	res, err := http.Post(ts.URL+"/v1/production/relative-surplus-from-productivity", "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}

func postRelativeSurplus(t *testing.T, ts *httptest.Server, body string) (*http.Response, relativeSurplusResponse) {
	t.Helper()
	res, err := http.Post(ts.URL+"/v1/production/relative-surplus", "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	t.Cleanup(func() { res.Body.Close() })
	var got relativeSurplusResponse
	if res.StatusCode == http.StatusOK {
		if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
			t.Fatalf("decode: %v", err)
		}
	}
	return res, got
}

// Ch. 12: one capitalist's gain is only temporary extra surplus-value. Social
// value 60 min/article, productivity factor 2.0 → individual value 30; selling
// 24 articles at the social value gains (60-30)*24 = 720. The value of
// labour-power and the working-day split do not move.
func TestRelativeSurplus_WorkplaceGainIsExtraSurplusValueOnly(t *testing.T) {
	t.Parallel()
	ts := newBridgeTestServer(t, nil)
	body := `{"working_day_total":720,"current_lpv":600,"productivity_factor":2.0,"social_value":60,"quantity":24}`
	res, got := postRelativeSurplus(t, ts, body)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	if got.Scope != ScopeWorkplace {
		t.Fatalf("Scope = %q, want workplace (default)", got.Scope)
	}
	if got.ExtraSurplusValue == nil || *got.ExtraSurplusValue != 720 {
		t.Fatalf("ExtraSurplusValue = %v, want 720", got.ExtraSurplusValue)
	}
	if got.NewLabourPowerValue != 600 {
		t.Fatalf("NewLabourPowerValue = %d, want unchanged 600", got.NewLabourPowerValue)
	}
	if got.RelativeSurplusValue != 0 {
		t.Fatalf("RelativeSurplusValue = %d, want 0 for a single workplace", got.RelativeSurplusValue)
	}
	if got.NewWorkingDay != got.OldWorkingDay {
		t.Fatalf("working day moved: %+v -> %+v", got.OldWorkingDay, got.NewWorkingDay)
	}
}

func TestRelativeSurplus_WageGoodsBranchMovesLPVWithoutExtraSurplusValue(t *testing.T) {
	t.Parallel()
	ts := newBridgeTestServer(t, nil)
	body := `{"working_day_total":720,"current_lpv":600,"productivity_factor":2.0,"scope":"wage_goods_branch"}`
	res, got := postRelativeSurplus(t, ts, body)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	if got.Scope != ScopeWageGoodsBranch {
		t.Fatalf("Scope = %q, want wage_goods_branch", got.Scope)
	}
	if got.ExtraSurplusValue != nil {
		t.Fatalf("ExtraSurplusValue = %d, want omitted", *got.ExtraSurplusValue)
	}
	if got.NewLabourPowerValue != 300 || got.RelativeSurplusValue != 300 {
		t.Fatalf("LPV=%d RSV=%d, want 300/300", got.NewLabourPowerValue, got.RelativeSurplusValue)
	}
}

func TestRelativeSurplus_WorkplaceRequiresSocialValueAndQuantity(t *testing.T) {
	t.Parallel()
	ts := newBridgeTestServer(t, nil)
	for name, body := range map[string]string{
		"no social value": `{"working_day_total":720,"current_lpv":600,"productivity_factor":2.0,"quantity":24}`,
		"no quantity":     `{"working_day_total":720,"current_lpv":600,"productivity_factor":2.0,"social_value":60}`,
	} {
		res, _ := postRelativeSurplus(t, ts, body)
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("%s: status = %d, want 400", name, res.StatusCode)
		}
	}
}

func TestRelativeSurplus_RejectsUnknownScope(t *testing.T) {
	t.Parallel()
	ts := newBridgeTestServer(t, nil)
	body := `{"working_day_total":720,"current_lpv":600,"productivity_factor":2.0,"scope":"national"}`
	res, _ := postRelativeSurplus(t, ts, body)
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}
