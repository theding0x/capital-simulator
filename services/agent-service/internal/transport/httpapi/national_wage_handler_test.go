package httpapi_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/transport/httpapi"
)

func newNationalWageHandler(t *testing.T) *httpapi.Handler {
	t.Helper()
	return httpapi.New(store.NewMemory(), nil)
}

func TestRegisterIntensity(t *testing.T) {
	t.Parallel()
	h := newNationalWageHandler(t)

	body, _ := json.Marshal(map[string]any{
		"country_code": "GB",
		"factor":       2.0,
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/intensities", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.RegisterIntensity(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", w.Code)
	}
	var got agent.NationalIntensity
	json.NewDecoder(w.Body).Decode(&got)
	if got.CountryCode != "GB" || got.Factor != 2.0 {
		t.Fatalf("response: got %+v", got)
	}
}

func TestRegisterIntensityBadFactor(t *testing.T) {
	t.Parallel()
	h := newNationalWageHandler(t)

	body, _ := json.Marshal(map[string]any{"country_code": "GB", "factor": 0})
	req := httptest.NewRequest(http.MethodPost, "/v1/intensities", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.RegisterIntensity(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want 400", w.Code)
	}
}

func TestRegisterDayWage(t *testing.T) {
	t.Parallel()
	h := newNationalWageHandler(t)

	body, _ := json.Marshal(map[string]any{
		"country_code":        "GB",
		"nominal_pence":       48,
		"working_day_minutes": 600,
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/wages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.RegisterDayWage(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want 201", w.Code)
	}
}

func TestRegisterDayWageAlreadyExists(t *testing.T) {
	t.Parallel()
	h := newNationalWageHandler(t)

	body, _ := json.Marshal(map[string]any{
		"country_code": "GB", "nominal_pence": 48, "working_day_minutes": 600,
	})
	for i := range 2 {
		req := httptest.NewRequest(http.MethodPost, "/v1/wages", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		h.RegisterDayWage(w, req)
		if i == 1 && w.Code != http.StatusConflict {
			t.Fatalf("second registration: got %d, want 409", w.Code)
		}
	}
}

func TestGetStandardisedWage(t *testing.T) {
	t.Parallel()
	h := newNationalWageHandler(t)

	// Register FR: 36 pence / 720 min
	body, _ := json.Marshal(map[string]any{
		"country_code": "FR", "nominal_pence": 36, "working_day_minutes": 720,
	})
	req := httptest.NewRequest(http.MethodPost, "/v1/wages", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.RegisterDayWage(rr, req)

	// Standardise to 600 min ref day → 36 * 600 / 720 = 30
	req2 := httptest.NewRequest(http.MethodGet, "/v1/wages/FR/standardised?reference_day_minutes=600", nil)
	req2.SetPathValue("country", "FR")
	w2 := httptest.NewRecorder()
	h.GetStandardisedWage(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", w2.Code)
	}
	var got agent.StandardisedWage
	json.NewDecoder(w2.Body).Decode(&got)
	if got.Amount != 30 {
		t.Fatalf("amount: got %d, want 30", got.Amount)
	}
}

func TestGetStandardisedWageNotFound(t *testing.T) {
	t.Parallel()
	h := newNationalWageHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/wages/XX/standardised", nil)
	req.SetPathValue("country", "XX")
	w := httptest.NewRecorder()
	h.GetStandardisedWage(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status: got %d, want 404", w.Code)
	}
}

func TestGetWageComparisonEmpty(t *testing.T) {
	t.Parallel()
	h := newNationalWageHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/v1/comparisons", nil)
	w := httptest.NewRecorder()
	h.GetWageComparison(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", w.Code)
	}
	var got agent.WageComparison
	json.NewDecoder(w.Body).Decode(&got)
	if got.DayWages == nil {
		t.Fatal("day_wages must be non-nil slice")
	}
}

func TestGetWageComparisonWithData(t *testing.T) {
	t.Parallel()
	h := newNationalWageHandler(t)

	seedNationalWageData(t, h)

	req := httptest.NewRequest(http.MethodGet, "/v1/comparisons", nil)
	w := httptest.NewRecorder()
	h.GetWageComparison(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", w.Code)
	}
	var cmp agent.WageComparison
	json.NewDecoder(w.Body).Decode(&cmp)
	if len(cmp.DayWages) != 2 {
		t.Fatalf("day_wages: got %d, want 2", len(cmp.DayWages))
	}
	if len(cmp.RelativePrices) != 2 {
		t.Fatalf("relative_prices: got %d, want 2", len(cmp.RelativePrices))
	}
}

func seedNationalWageData(t *testing.T, h *httpapi.Handler) {
	t.Helper()
	for _, tc := range []struct {
		cc     string
		pence  int
		mins   int
		factor float64
	}{
		{"GB", 48, 600, 2.0},
		{"DE", 24, 720, 1.0},
	} {
		wBody, _ := json.Marshal(map[string]any{
			"country_code": tc.cc, "nominal_pence": tc.pence, "working_day_minutes": tc.mins,
		})
		wr := httptest.NewRequest(http.MethodPost, "/v1/wages", bytes.NewReader(wBody))
		wr.Header.Set("Content-Type", "application/json")
		h.RegisterDayWage(httptest.NewRecorder(), wr)

		iBody, _ := json.Marshal(map[string]any{"country_code": tc.cc, "factor": tc.factor})
		ir := httptest.NewRequest(http.MethodPost, "/v1/intensities", bytes.NewReader(iBody))
		ir.Header.Set("Content-Type", "application/json")
		h.RegisterIntensity(httptest.NewRecorder(), ir)
	}
}
