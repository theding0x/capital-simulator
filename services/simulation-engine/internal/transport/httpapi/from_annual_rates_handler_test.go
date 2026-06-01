package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// seedPriceRevolution inserts a value-revolution event with the given
// basis-points change and returns its id.
func seedPriceRevolution(t *testing.T, mem *store.Memory, bpChange int64) string {
	t.Helper()
	rec, err := mem.CreatePriceRevolution(context.Background(), circulation.PriceRevolutionRecord{
		Event: circulation.ValueRevolutionEvent{
			IndustrialCapitalID: circulation.IndustrialCapitalID("ic-test"),
			Affecting:           circulation.AffectingOutput,
			DirectionPence:      circulation.Pence(100),
			BasisPointsChange:   bpChange,
		},
	})
	if err != nil {
		t.Fatalf("seed price revolution: %v", err)
	}
	return string(rec.Event.ID)
}

func adjustedAnnualRate(t *testing.T, h *Handler, eventID, query string) map[string]any {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet,
		"/v1/price-revolutions/"+eventID+"/adjusted-annual-rate?"+query, nil)
	req.SetPathValue("id", eventID)
	w := httptest.NewRecorder()
	h.GetAdjustedAnnualRate(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("adjusted-annual-rate: expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode rate: %v", err)
	}
	return resp
}

// TestAdjustedRateThenBuildScheme is the issue #222 acceptance chain: feed a
// price-revolution event → compute the adjusted annual rate → build a
// reproduction scheme from annual rates → assert the keystone balances.
func TestAdjustedRateThenBuildScheme(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	h := New(nil, Deps{PriceRevolutions: mem, SimpleReproduction: mem})

	// (1) Event with no net price change: adjusted rate == base, so the
	// canonical balanced fixture stays balanced.
	eventID := seedPriceRevolution(t, mem, 0)

	// (2) Adjusted annual rate for Dept I's capital: v=1000, s=1000, n=1×.
	rate := adjustedAnnualRate(t, h, eventID, "variable_pence=1000&surplus_pence=1000&n_basis_points=10000")
	adjustedBP := int64(rate["adjusted_annual_basis_points"].(float64))
	if adjustedBP != 10000 {
		t.Fatalf("adjusted annual bp: got %d, want 10000 (no price change)", adjustedBP)
	}

	// (3) Build a scheme from annual rates; Dept I uses the adjusted rate.
	body := fmt.Sprintf(`{"period":"1865-from-rates","departments":[`+
		`{"department":"department_i","constant_pence":4000,"variable_pence":1000,"annual_surplus_bp":%d},`+
		`{"department":"department_ii","constant_pence":2000,"variable_pence":500,"annual_surplus_bp":10000}]}`, adjustedBP)
	br := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/simple/schemes/from-annual-rates", bytes.NewBufferString(body))
	br.Header.Set("Content-Type", "application/json")
	bw := httptest.NewRecorder()
	h.BuildReproductionFromAnnualRates(bw, br)
	if bw.Code != http.StatusOK {
		t.Fatalf("from-annual-rates: expected 200, got %d: %s", bw.Code, bw.Body.String())
	}
	var scheme struct {
		IsBalanced  bool `json:"is_balanced"`
		DepartmentI struct {
			SurplusPence int64 `json:"surplus_pence"`
		} `json:"department_i"`
	}
	if err := json.NewDecoder(bw.Body).Decode(&scheme); err != nil {
		t.Fatalf("decode scheme: %v", err)
	}
	// (4) Dept I surplus = 1000 (derived), so I(v+s)=2000 == II(c)=2000.
	if scheme.DepartmentI.SurplusPence != 1000 {
		t.Fatalf("dept I surplus: got %d, want 1000 (derived from annual rate)", scheme.DepartmentI.SurplusPence)
	}
	if !scheme.IsBalanced {
		t.Fatal("expected balanced scheme from the chained pipeline")
	}
}

// TestAdjustedAnnualRate_PriceRiseLiftsRate: a +10% output price revolution
// lifts the annual rate from 100% to 120% (the Ch.15→16 effect).
func TestAdjustedAnnualRate_PriceRiseLiftsRate(t *testing.T) {
	t.Parallel()
	mem := store.NewMemory()
	h := New(nil, Deps{PriceRevolutions: mem})
	eventID := seedPriceRevolution(t, mem, 1000) // +10%

	rate := adjustedAnnualRate(t, h, eventID, "variable_pence=1000&surplus_pence=1000&n_basis_points=10000")
	if base := int64(rate["base_annual_basis_points"].(float64)); base != 10000 {
		t.Fatalf("base annual bp: got %d, want 10000", base)
	}
	if adj := int64(rate["adjusted_annual_basis_points"].(float64)); adj != 12000 {
		t.Fatalf("adjusted annual bp: got %d, want 12000", adj)
	}
	if s := int64(rate["adjusted_surplus_pence"].(float64)); s != 1200 {
		t.Fatalf("adjusted surplus: got %d, want 1200", s)
	}
}

// TestAdjustedAnnualRate_NotFound returns 404 for a missing event.
func TestAdjustedAnnualRate_NotFound(t *testing.T) {
	t.Parallel()
	h := New(nil, Deps{PriceRevolutions: store.NewMemory()})
	req := httptest.NewRequest(http.MethodGet,
		"/v1/price-revolutions/ghost/adjusted-annual-rate?variable_pence=1000&surplus_pence=1000", nil)
	req.SetPathValue("id", "ghost")
	w := httptest.NewRecorder()
	h.GetAdjustedAnnualRate(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

// TestBuildFromAnnualRates_BadJSON returns 400.
func TestBuildFromAnnualRates_BadJSON(t *testing.T) {
	t.Parallel()
	h := New(nil, Deps{SimpleReproduction: store.NewMemory()})
	req := httptest.NewRequest(http.MethodPost,
		"/v1/reproduction/simple/schemes/from-annual-rates", bytes.NewBufferString(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.BuildReproductionFromAnnualRates(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
