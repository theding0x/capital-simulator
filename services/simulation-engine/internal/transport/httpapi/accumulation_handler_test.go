package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

// § spinner example §1: £10,000 capital, s'=100%, full accumulation, 3 periods.
// Abraham sequence: surplus 2000 → 400 → 80.
func TestRunExtendedReproduction_AbrahamSequence(t *testing.T) {
	t.Parallel()
	body := `{"constant_capital":8000,"variable_capital":2000,"surplus_rate":1.0,"accum_rate":1.0,"composition_ratio":0.8,"periods":3}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/extended", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil)
	h.RunExtendedReproduction(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Cycles []struct {
			SurplusProduced int64 `json:"surplus_produced"`
			Revenue         int64 `json:"revenue"`
		} `json:"cycles"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Cycles) != 3 {
		t.Fatalf("want 3 cycles, got %d", len(resp.Cycles))
	}

	wantSurplus := []int64{2000, 400, 80}
	for i, c := range resp.Cycles {
		if c.SurplusProduced != wantSurplus[i] {
			t.Errorf("cycle %d: surplus_produced = %d, want %d", i+1, c.SurplusProduced, wantSurplus[i])
		}
		if c.Revenue != 0 {
			t.Errorf("cycle %d: revenue = %d, want 0 (full accumulation)", i+1, c.Revenue)
		}
	}
}

// § partial accumulation §3: revenue = half the surplus.
func TestRunExtendedReproduction_PartialAccumulation(t *testing.T) {
	t.Parallel()
	body := `{"constant_capital":8000,"variable_capital":2000,"surplus_rate":1.0,"accum_rate":0.5,"composition_ratio":0.8,"periods":1}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/extended", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil)
	h.RunExtendedReproduction(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Cycles []struct {
			NewConstant int64 `json:"new_constant"`
			NewVariable int64 `json:"new_variable"`
			Revenue     int64 `json:"revenue"`
		} `json:"cycles"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	c := resp.Cycles[0]
	if c.NewConstant != 800 {
		t.Errorf("new_constant = %d, want 800", c.NewConstant)
	}
	if c.NewVariable != 200 {
		t.Errorf("new_variable = %d, want 200", c.NewVariable)
	}
	if c.Revenue != 1000 {
		t.Errorf("revenue = %d, want 1000", c.Revenue)
	}
}

func TestRunExtendedReproduction_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/extended", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil)
	h.RunExtendedReproduction(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestRunExtendedReproduction_BadRequest_ZeroVariableCapital(t *testing.T) {
	t.Parallel()
	body := `{"constant_capital":8000,"variable_capital":0,"surplus_rate":1.0,"accum_rate":1.0,"composition_ratio":0.8,"periods":3}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/extended", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil)
	h.RunExtendedReproduction(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestRunExtendedReproduction_BadRequest_AccumRateOutOfRange(t *testing.T) {
	t.Parallel()
	body := `{"constant_capital":8000,"variable_capital":2000,"surplus_rate":1.0,"accum_rate":1.5,"composition_ratio":0.8,"periods":3}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/extended", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil)
	h.RunExtendedReproduction(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

// § split-surplus: spinner example — SplitSurplus(2000, 1.0, 0.8).
func TestSplitSurplusValue_SpinnerExample(t *testing.T) {
	t.Parallel()
	body := `{"surplus":2000,"accum_rate":1.0,"composition_ratio":0.8}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/split-surplus", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil)
	h.SplitSurplusValue(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		NewConstant int64 `json:"new_constant"`
		NewVariable int64 `json:"new_variable"`
		Revenue     int64 `json:"revenue"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.NewConstant != 1600 {
		t.Errorf("new_constant = %d, want 1600", resp.NewConstant)
	}
	if resp.NewVariable != 400 {
		t.Errorf("new_variable = %d, want 400", resp.NewVariable)
	}
	if resp.Revenue != 0 {
		t.Errorf("revenue = %d, want 0", resp.Revenue)
	}
}

// § partial accumulation §3: SplitSurplus(2000, 0.5, 0.8); Revenue == 1000.
func TestSplitSurplusValue_PartialAccumulation(t *testing.T) {
	t.Parallel()
	body := `{"surplus":2000,"accum_rate":0.5,"composition_ratio":0.8}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/split-surplus", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil)
	h.SplitSurplusValue(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		NewConstant int64 `json:"new_constant"`
		NewVariable int64 `json:"new_variable"`
		Revenue     int64 `json:"revenue"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.NewConstant != 800 {
		t.Errorf("new_constant = %d, want 800", resp.NewConstant)
	}
	if resp.NewVariable != 200 {
		t.Errorf("new_variable = %d, want 200", resp.NewVariable)
	}
	if resp.Revenue != 1000 {
		t.Errorf("revenue = %d, want 1000", resp.Revenue)
	}
}

func TestSplitSurplusValue_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/split-surplus", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil)
	h.SplitSurplusValue(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestSplitSurplusValue_BadRequest_NegativeSurplus(t *testing.T) {
	t.Parallel()
	body := `{"surplus":-1,"accum_rate":1.0,"composition_ratio":0.8}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/split-surplus", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil)
	h.SplitSurplusValue(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}
