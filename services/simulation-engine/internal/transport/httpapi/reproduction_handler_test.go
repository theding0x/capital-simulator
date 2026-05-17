package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

func TestRunSimpleReproduction_CoreExample(t *testing.T) {
	t.Parallel()
	// "capital of £1,000 begets yearly a surplus-value of £200"; 5 periods
	body := `{"constant_capital":800,"variable_capital":200,"surplus_rate":1.0,"periods":5}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/simple", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	h.RunSimpleReproduction(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Cycles []struct {
			SurplusTotal int64 `json:"surplus_total"`
			Revenue      int64 `json:"revenue"`
			Accumulated  int64 `json:"accumulated"`
		} `json:"cycles"`
		RepaymentPeriod int64 `json:"repayment_period"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Cycles) != 5 {
		t.Fatalf("want 5 cycles, got %d", len(resp.Cycles))
	}
	for i, c := range resp.Cycles {
		if c.SurplusTotal != 200 {
			t.Errorf("cycle %d: surplus_total = %d, want 200", i+1, c.SurplusTotal)
		}
		if c.Revenue != 200 {
			t.Errorf("cycle %d: revenue = %d, want 200", i+1, c.Revenue)
		}
		if c.Accumulated != 0 {
			t.Errorf("cycle %d: accumulated = %d, want 0", i+1, c.Accumulated)
		}
	}
	if resp.RepaymentPeriod != 5 {
		t.Errorf("repayment_period = %d, want 5", resp.RepaymentPeriod)
	}
}

func TestRunSimpleReproduction_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/simple", bytes.NewBufferString(`{bad json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	h.RunSimpleReproduction(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestRunSimpleReproduction_BadRequest_ZeroVariableCapital(t *testing.T) {
	t.Parallel()
	body := `{"constant_capital":800,"variable_capital":0,"surplus_rate":1.0,"periods":5}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/simple", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	h.RunSimpleReproduction(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestComputeRepaymentPeriod_CoreExample(t *testing.T) {
	t.Parallel()
	// "capital of £1,000 / £200 annual revenue = 5 periods"
	body := `{"capital":1000,"annual_revenue":200}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/repayment-period", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	h.ComputeRepaymentPeriod(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		RepaymentPeriod int64 `json:"repayment_period"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.RepaymentPeriod != 5 {
		t.Errorf("repayment_period = %d, want 5", resp.RepaymentPeriod)
	}
}

func TestComputeRepaymentPeriod_HalfConsumed(t *testing.T) {
	t.Parallel()
	// "if only one half were consumed ... the same result at the end of 10 years"
	body := `{"capital":1000,"annual_revenue":100}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/repayment-period", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	h.ComputeRepaymentPeriod(w, req)

	var resp struct {
		RepaymentPeriod int64 `json:"repayment_period"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.RepaymentPeriod != 10 {
		t.Errorf("repayment_period = %d, want 10", resp.RepaymentPeriod)
	}
}

func TestComputeRepaymentPeriod_BadRequest_NegativeCapital(t *testing.T) {
	t.Parallel()
	body := `{"capital":-1,"annual_revenue":200}`
	req := httptest.NewRequest(http.MethodPost, "/v1/reproductions/repayment-period", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h := httpapi.New(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	h.ComputeRepaymentPeriod(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}
