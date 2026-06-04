package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

// Issue #412 — ProtectionSystem gained a Create path; a created record persists
// and surfaces in the stage's industrial-capital genesis.
func TestCreateProtectionSystem_Success(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForGenesis(t, st)
	h := &httpapi.Handler{HistoricalStages: st, ProtectionSystems: st}

	body := `{"tariff_rate_bps":3000,"beneficiary":"English manufacturers","period_start":"17th c","period_end":"19th c"}`
	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/"+string(stageID)+"/protection-systems", strings.NewReader(body))
	req.SetPathValue("id", string(stageID))
	w := httptest.NewRecorder()
	h.CreateProtectionSystem(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); !strings.Contains(loc, "protection-systems") {
		t.Errorf("missing Location header: %q", loc)
	}
	var resp struct {
		ID            string `json:"id"`
		TariffRateBps int64  `json:"tariff_rate_bps"`
		Beneficiary   string `json:"beneficiary"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID == "" || resp.TariffRateBps != 3000 || resp.Beneficiary != "English manufacturers" {
		t.Errorf("unexpected response: %+v", resp)
	}

	// It now appears in the stage's genesis (the gap #412 closed).
	hg := &httpapi.Handler{HistoricalStages: st, CapitalOrigins: st, ColonialTransfers: st, NationalDebts: st, ProtectionSystems: st}
	greq := httptest.NewRequest(http.MethodGet, "/v1/historical-stages/"+string(stageID)+"/genesis", nil)
	greq.SetPathValue("id", string(stageID))
	gw := httptest.NewRecorder()
	hg.GetIndustrialCapitalGenesis(gw, greq)
	if gw.Code != http.StatusOK {
		t.Fatalf("genesis want 200, got %d: %s", gw.Code, gw.Body.String())
	}
	var genesis struct {
		ProtectionSystems []struct {
			Beneficiary string `json:"beneficiary"`
		} `json:"protection_systems"`
	}
	if err := json.NewDecoder(gw.Body).Decode(&genesis); err != nil {
		t.Fatalf("decode genesis: %v", err)
	}
	if len(genesis.ProtectionSystems) != 1 || genesis.ProtectionSystems[0].Beneficiary != "English manufacturers" {
		t.Errorf("genesis protection_systems = %+v, want one English-manufacturers entry", genesis.ProtectionSystems)
	}
}

func TestCreateProtectionSystem_MissingBeneficiary(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForGenesis(t, st)
	h := &httpapi.Handler{HistoricalStages: st, ProtectionSystems: st}

	req := httptest.NewRequest(http.MethodPost, "/v1/historical-stages/"+string(stageID)+"/protection-systems", strings.NewReader(`{"tariff_rate_bps":3000,"beneficiary":""}`))
	req.SetPathValue("id", string(stageID))
	w := httptest.NewRecorder()
	h.CreateProtectionSystem(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400 for empty beneficiary, got %d: %s", w.Code, w.Body.String())
	}
}
