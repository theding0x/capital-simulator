package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newProductionTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	h := New(nil, nil, nil, nil, nil, nil, nil)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/production/working-day", h.RecordWorkingDay)
	mux.HandleFunc("POST /v1/production/working-day/shorten", h.ShortenWorkingDay)
	mux.HandleFunc("GET /v1/production/rate-of-surplus-value", h.GetProductionRate)
	mux.HandleFunc("POST /v1/production/extra-surplus-value", h.ComputeExtraSurplusValue)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// §1 fixture: total=720, lpv=600 → WorkingDay{720, 600, 120}
func TestRecordWorkingDay(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	body := `{"total":720,"labour_power_value":600}`
	res, err := http.Post(ts.URL+"/v1/production/working-day", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	var got struct {
		Total           int64 `json:"total"`
		NecessaryLabour int64 `json:"necessary_labour"`
		SurplusLabour   int64 `json:"surplus_labour"`
	}
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Total != 720 || got.NecessaryLabour != 600 || got.SurplusLabour != 120 {
		t.Fatalf("unexpected response: %+v", got)
	}
}

// validation: labour_power_value >= total must be rejected
func TestRecordWorkingDay_InvalidLPV(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	body := `{"total":720,"labour_power_value":720}`
	res, err := http.Post(ts.URL+"/v1/production/working-day", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

// §1 fixture: shorten(wd={720,600,120}, newLPV=540) → {720,540,180}, delta=60
func TestShortenWorkingDay(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	body := `{"total":720,"necessary_labour":600,"surplus_labour":120,"new_labour_power_value":540}`
	res, err := http.Post(ts.URL+"/v1/production/working-day/shorten", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	var got struct {
		WorkingDay struct {
			Total           int64 `json:"total"`
			NecessaryLabour int64 `json:"necessary_labour"`
			SurplusLabour   int64 `json:"surplus_labour"`
		} `json:"working_day"`
		RelativeSurplusValue int64 `json:"relative_surplus_value"`
	}
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.WorkingDay.NecessaryLabour != 540 {
		t.Fatalf("expected NecessaryLabour=540, got %d", got.WorkingDay.NecessaryLabour)
	}
	if got.WorkingDay.SurplusLabour != 180 {
		t.Fatalf("expected SurplusLabour=180, got %d", got.WorkingDay.SurplusLabour)
	}
	if got.RelativeSurplusValue != 60 {
		t.Fatalf("expected RelativeSurplusValue=60, got %d", got.RelativeSurplusValue)
	}
}

// validation: new_labour_power_value >= total must be rejected
func TestShortenWorkingDay_InvalidNewLPV(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	body := `{"total":720,"necessary_labour":600,"surplus_labour":120,"new_labour_power_value":720}`
	res, err := http.Post(ts.URL+"/v1/production/working-day/shorten", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

// §1 rate: necessary=600, surplus=120 → rate=0.2
func TestGetProductionRate(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	res, err := http.Get(ts.URL + "/v1/production/rate-of-surplus-value?necessary=600&surplus=120")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	var got struct {
		Rate            float64 `json:"rate"`
		SurplusLabour   int64   `json:"surplus_labour"`
		NecessaryLabour int64   `json:"necessary_labour"`
	}
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Rate != 0.2 {
		t.Fatalf("expected rate=0.2, got %f", got.Rate)
	}
	if got.SurplusLabour != 120 || got.NecessaryLabour != 600 {
		t.Fatalf("unexpected echoed fields: %+v", got)
	}
}

// validation: missing necessary query param → 400
func TestGetProductionRate_MissingParam(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	res, err := http.Get(ts.URL + "/v1/production/rate-of-surplus-value?surplus=120")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", res.StatusCode)
	}
}

// §1 fixture: IndividualValue=30, SocialValue=60, Quantity=24 → extra=720, per_unit=30
func TestComputeExtraSurplusValue(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	body := `{"individual_value":30,"social_value":60,"quantity":24}`
	res, err := http.Post(ts.URL+"/v1/production/extra-surplus-value", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	var got struct {
		ExtraSurplusValue int64 `json:"extra_surplus_value"`
		PerUnit           int64 `json:"per_unit"`
		IndividualValue   int64 `json:"individual_value"`
		SocialValue       int64 `json:"social_value"`
		Quantity          int64 `json:"quantity"`
	}
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ExtraSurplusValue != 720 {
		t.Fatalf("expected extra_surplus_value=720, got %d", got.ExtraSurplusValue)
	}
	if got.PerUnit != 30 {
		t.Fatalf("expected per_unit=30, got %d", got.PerUnit)
	}
}

// §1 invariant: extra=0 when iv==sv
func TestComputeExtraSurplusValue_ZeroWhenEqual(t *testing.T) {
	t.Parallel()
	ts := newProductionTestServer(t)

	body := `{"individual_value":60,"social_value":60,"quantity":24}`
	res, err := http.Post(ts.URL+"/v1/production/extra-surplus-value", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", res.StatusCode)
	}
	var got struct {
		ExtraSurplusValue int64 `json:"extra_surplus_value"`
	}
	json.NewDecoder(res.Body).Decode(&got) //nolint
	if got.ExtraSurplusValue != 0 {
		t.Fatalf("expected 0 when iv==sv, got %d", got.ExtraSurplusValue)
	}
}
