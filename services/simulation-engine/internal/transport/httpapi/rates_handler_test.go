package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newRatesTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	h := New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/surplus-value/rates", h.ComputeRatesOfSurplusValue)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// §formula I fixture: "6 hours surplus-labour / 6 hours necessary labour = 100%"
// §formula II fixture: "6 hours surplus-labour / 12-hour working-day = 50%"
func TestComputeRatesOfSurplusValue_SixAndSix(t *testing.T) {
	t.Parallel()
	ts := newRatesTestServer(t)
	body := `{"necessary_labour_minutes":360,"surplus_labour_minutes":360}`
	res, err := http.Post(ts.URL+"/v1/surplus-value/rates", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got ratesResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.FormulaI != 1.0 {
		t.Errorf("FormulaI = %f, want 1.0", got.FormulaI)
	}
	if got.FormulaII != 0.5 {
		t.Errorf("FormulaII = %f, want 0.5", got.FormulaII)
	}
	if got.FormulaI != got.FormulaIII {
		t.Errorf("FormulaI (%f) != FormulaIII (%f)", got.FormulaI, got.FormulaIII)
	}
	if got.WorkingDayMinutes != 720 {
		t.Errorf("WorkingDayMinutes = %d, want 720", got.WorkingDayMinutes)
	}
}

// "rate of exploitation of 300%": NL=3h=180min, SL=9h=540min
func TestComputeRatesOfSurplusValue_AgriculturalLabourer(t *testing.T) {
	t.Parallel()
	ts := newRatesTestServer(t)
	body := `{"necessary_labour_minutes":180,"surplus_labour_minutes":540}`
	res, err := http.Post(ts.URL+"/v1/surplus-value/rates", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got ratesResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.FormulaI != 3.0 {
		t.Errorf("FormulaI = %f, want 3.0 (300%%)", got.FormulaI)
	}
	if got.FormulaII >= 1.0 {
		t.Errorf("FormulaII = %f, must be < 1.0", got.FormulaII)
	}
}

func TestComputeRatesOfSurplusValue_RejectsZeroNecessaryLabour(t *testing.T) {
	t.Parallel()
	ts := newRatesTestServer(t)
	body := `{"necessary_labour_minutes":0,"surplus_labour_minutes":360}`
	res, err := http.Post(ts.URL+"/v1/surplus-value/rates", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}

func TestComputeRatesOfSurplusValue_RejectsNegativeSurplusLabour(t *testing.T) {
	t.Parallel()
	ts := newRatesTestServer(t)
	body := `{"necessary_labour_minutes":360,"surplus_labour_minutes":-1}`
	res, err := http.Post(ts.URL+"/v1/surplus-value/rates", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}
