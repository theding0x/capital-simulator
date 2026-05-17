package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newSurplusValueTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	h := New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/surplus-value/absolute", h.ComputeAbsoluteSurplusValue)
	mux.HandleFunc("POST /v1/surplus-value/relative", h.ComputeRelativeSurplusValue)
	mux.HandleFunc("GET /v1/surplus-value/rate", h.GetSurplusValueRate)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// §1 absolute fixture: extending a 720/600/120 day by 60 minutes lands all
// the extension in surplus-labour.
func TestComputeAbsoluteSurplusValue_Para1(t *testing.T) {
	t.Parallel()
	ts := newSurplusValueTestServer(t)
	body := `{"working_day_minutes":720,"necessary_labour_minutes":600,"extension_minutes":60}`
	res, err := http.Post(ts.URL+"/v1/surplus-value/absolute", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got absoluteSurplusValueResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.NewWorkingDay.Total != 780 {
		t.Fatalf("new Total = %d, want 780", got.NewWorkingDay.Total)
	}
	if got.NewWorkingDay.NecessaryLabour != 600 {
		t.Fatalf("new NecessaryLabour = %d, want 600 (fixed)", got.NewWorkingDay.NecessaryLabour)
	}
	if got.AbsoluteSurplusValue.LabourMinutes != 60 {
		t.Fatalf("ASV = %d, want 60", got.AbsoluteSurplusValue.LabourMinutes)
	}
	if got.AbsoluteSurplusValue.Origin != "absolute" {
		t.Fatalf("Origin = %q, want absolute", got.AbsoluteSurplusValue.Origin)
	}
}

func TestComputeAbsoluteSurplusValue_RejectsNonPositiveExtension(t *testing.T) {
	t.Parallel()
	ts := newSurplusValueTestServer(t)
	body := `{"working_day_minutes":720,"necessary_labour_minutes":600,"extension_minutes":0}`
	res, err := http.Post(ts.URL+"/v1/surplus-value/absolute", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}

// §1 relative fixture: factor 2.0 halves necessary labour.
func TestComputeRelativeSurplusValue_Para1(t *testing.T) {
	t.Parallel()
	ts := newSurplusValueTestServer(t)
	body := `{"working_day_minutes":720,"necessary_labour_minutes":600,"productivity_factor":2.0}`
	res, err := http.Post(ts.URL+"/v1/surplus-value/relative", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got relativeSurplusValueResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.NewWorkingDay.NecessaryLabour != 300 {
		t.Fatalf("new NecessaryLabour = %d, want 300", got.NewWorkingDay.NecessaryLabour)
	}
	if got.NewWorkingDay.SurplusLabour != 420 {
		t.Fatalf("new SurplusLabour = %d, want 420", got.NewWorkingDay.SurplusLabour)
	}
	if got.RelativeSurplusValue.LabourMinutes != 300 {
		t.Fatalf("RSV = %d, want 300", got.RelativeSurplusValue.LabourMinutes)
	}
	if got.RelativeSurplusValue.Origin != "relative" {
		t.Fatalf("Origin = %q, want relative", got.RelativeSurplusValue.Origin)
	}
}

// §1 Mill critique: rate of profit (4%) ≠ rate of surplus-value (20%)
// for £400c + £100v + £20s. The endpoint surfaces both and asserts the
// Mill critique held.
func TestGetSurplusValueRate_Para1_MillCritique(t *testing.T) {
	t.Parallel()
	ts := newSurplusValueTestServer(t)
	res, err := http.Get(ts.URL + "/v1/surplus-value/rate?surplus_labour=20&necessary_labour=100&total_capital=500")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got surplusValueRateResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.RateOfSurplusValue != 0.20 {
		t.Fatalf("RateOfSurplusValue = %v, want 0.20", got.RateOfSurplusValue)
	}
	if got.RateOfProfit != 0.04 {
		t.Fatalf("RateOfProfit = %v, want 0.04", got.RateOfProfit)
	}
	if !got.MillCritiqueHolds {
		t.Fatal("MillCritiqueHolds must be true when total_capital > necessary_labour")
	}
}

func TestGetSurplusValueRate_MissingParam(t *testing.T) {
	t.Parallel()
	ts := newSurplusValueTestServer(t)
	res, err := http.Get(ts.URL + "/v1/surplus-value/rate?surplus_labour=20&necessary_labour=100")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}
