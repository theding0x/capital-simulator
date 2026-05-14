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
	h := New(nil, nil, nil, fetcher)
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
	body := `{"working_day_total":720,"current_lpv":600,"productivity_factor":2.0}`
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
	body := `{"working_day_total":720,"current_lpv":600,"source":"cooperation","source_id":"abc"}`
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
