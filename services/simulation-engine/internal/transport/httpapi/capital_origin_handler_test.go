package httpapi_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/transport/httpapi"
)

// § Bank of England founding — 8% to private bankers.
const bankOfEnglandDebtBody = `{
  "amount_pence": 2400000,
  "interest_rate_bps": 800,
  "creditor_class": "private-bankers"
}`

// § Liverpool slave trade, 1730 — 15 ships.
const liverpoolSlaveTransferBody = `{
  "from": "West Africa",
  "to": "England",
  "value_pence": 15000,
  "method": "slave-trade"
}`

// § Colonial plunder from the Americas.
const colonialPlunderOriginBody = `{
  "source": "colonial-plunder",
  "amount_pence": 200000,
  "period": "1500-1800"
}`

func seedStageForGenesis(t *testing.T, st *store.Memory) simulation.HistoricalStageID {
	t.Helper()
	stage, err := st.CreateHistoricalStage(t.Context(), simulation.HistoricalStage{
		Name:        "England 15th–18th c.",
		Description: "Primitive accumulation",
		PrimitiveAccumulations: []simulation.PrimitiveAccumulation{
			{Period: "1450–1640", Method: "enclosure", LabourersExpropriated: 40000, CapitalFormed: 80000},
		},
	})
	if err != nil {
		t.Fatalf("seedStageForGenesis: %v", err)
	}
	return stage.ID
}

// --- CreateCapitalOrigin ---

func TestCreateCapitalOrigin_Success(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForGenesis(t, st)
	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, CapitalOrigins: st})

	req := httptest.NewRequest(http.MethodPost,
		"/v1/historical-stages/"+string(stageID)+"/capital-origins",
		strings.NewReader(colonialPlunderOriginBody))
	req.SetPathValue("id", string(stageID))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateCapitalOrigin(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); !strings.Contains(loc, "capital-origins") {
		t.Errorf("missing Location header: %q", loc)
	}
	var resp struct {
		ID          string `json:"id"`
		Source      string `json:"source"`
		AmountPence int64  `json:"amount_pence"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID == "" {
		t.Error("response missing id")
	}
	if resp.Source != "colonial-plunder" {
		t.Errorf("source = %q, want colonial-plunder", resp.Source)
	}
	if resp.AmountPence != 200000 {
		t.Errorf("amount_pence = %d, want 200000", resp.AmountPence)
	}
}

func TestCreateCapitalOrigin_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForGenesis(t, st)
	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, CapitalOrigins: st})

	req := httptest.NewRequest(http.MethodPost,
		"/v1/historical-stages/"+string(stageID)+"/capital-origins",
		strings.NewReader(`{nope`))
	req.SetPathValue("id", string(stageID))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateCapitalOrigin(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateCapitalOrigin_BadRequest_MissingSource(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForGenesis(t, st)
	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, CapitalOrigins: st})

	body := `{"amount_pence": 1000}`
	req := httptest.NewRequest(http.MethodPost,
		"/v1/historical-stages/"+string(stageID)+"/capital-origins",
		strings.NewReader(body))
	req.SetPathValue("id", string(stageID))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateCapitalOrigin(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateCapitalOrigin_NotFound_MissingStage(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, CapitalOrigins: st})

	req := httptest.NewRequest(http.MethodPost,
		"/v1/historical-stages/ghost/capital-origins",
		strings.NewReader(colonialPlunderOriginBody))
	req.SetPathValue("id", "ghost")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateCapitalOrigin(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d: %s", w.Code, w.Body.String())
	}
}

// --- CreateColonialTransfer ---

func TestCreateColonialTransfer_Success(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForGenesis(t, st)
	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, ColonialTransfers: st})

	req := httptest.NewRequest(http.MethodPost,
		"/v1/historical-stages/"+string(stageID)+"/colonial-transfers",
		strings.NewReader(liverpoolSlaveTransferBody))
	req.SetPathValue("id", string(stageID))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateColonialTransfer(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		ID         string `json:"id"`
		ValuePence int64  `json:"value_pence"`
		Method     string `json:"method"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID == "" {
		t.Error("response missing id")
	}
	if resp.ValuePence != 15000 {
		t.Errorf("value_pence = %d, want 15000", resp.ValuePence)
	}
	if resp.Method != "slave-trade" {
		t.Errorf("method = %q, want slave-trade", resp.Method)
	}
}

func TestCreateColonialTransfer_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForGenesis(t, st)
	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, ColonialTransfers: st})

	req := httptest.NewRequest(http.MethodPost,
		"/v1/historical-stages/"+string(stageID)+"/colonial-transfers",
		strings.NewReader(`{bad`))
	req.SetPathValue("id", string(stageID))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateColonialTransfer(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateColonialTransfer_BadRequest_ZeroValue(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForGenesis(t, st)
	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, ColonialTransfers: st})

	body := `{"from":"Americas","to":"England","value_pence":0,"method":"colonial-plunder"}`
	req := httptest.NewRequest(http.MethodPost,
		"/v1/historical-stages/"+string(stageID)+"/colonial-transfers",
		strings.NewReader(body))
	req.SetPathValue("id", string(stageID))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateColonialTransfer(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateColonialTransfer_NotFound_MissingStage(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, ColonialTransfers: st})

	req := httptest.NewRequest(http.MethodPost,
		"/v1/historical-stages/ghost/colonial-transfers",
		strings.NewReader(liverpoolSlaveTransferBody))
	req.SetPathValue("id", "ghost")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateColonialTransfer(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d: %s", w.Code, w.Body.String())
	}
}

// --- CreateNationalDebt ---

func TestCreateNationalDebt_Success(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForGenesis(t, st)
	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, NationalDebts: st})

	req := httptest.NewRequest(http.MethodPost,
		"/v1/historical-stages/"+string(stageID)+"/national-debts",
		strings.NewReader(bankOfEnglandDebtBody))
	req.SetPathValue("id", string(stageID))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateNationalDebt(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("want 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		ID              string `json:"id"`
		AmountPence     int64  `json:"amount_pence"`
		InterestRateBps int64  `json:"interest_rate_bps"`
		CreditorClass   string `json:"creditor_class"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID == "" {
		t.Error("response missing id")
	}
	if resp.InterestRateBps != 800 {
		t.Errorf("interest_rate_bps = %d, want 800", resp.InterestRateBps)
	}
	if resp.CreditorClass != "private-bankers" {
		t.Errorf("creditor_class = %q, want private-bankers", resp.CreditorClass)
	}
}

func TestCreateNationalDebt_BadRequest_MalformedJSON(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForGenesis(t, st)
	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, NationalDebts: st})

	req := httptest.NewRequest(http.MethodPost,
		"/v1/historical-stages/"+string(stageID)+"/national-debts",
		strings.NewReader(`{bad`))
	req.SetPathValue("id", string(stageID))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateNationalDebt(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateNationalDebt_BadRequest_ZeroInterestRate(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForGenesis(t, st)
	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, NationalDebts: st})

	body := `{"amount_pence":1000000,"interest_rate_bps":0,"creditor_class":"private-bankers"}`
	req := httptest.NewRequest(http.MethodPost,
		"/v1/historical-stages/"+string(stageID)+"/national-debts",
		strings.NewReader(body))
	req.SetPathValue("id", string(stageID))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateNationalDebt(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", w.Code)
	}
}

func TestCreateNationalDebt_NotFound_MissingStage(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, NationalDebts: st})

	req := httptest.NewRequest(http.MethodPost,
		"/v1/historical-stages/ghost/national-debts",
		strings.NewReader(bankOfEnglandDebtBody))
	req.SetPathValue("id", "ghost")
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.CreateNationalDebt(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d: %s", w.Code, w.Body.String())
	}
}

// --- GetIndustrialCapitalGenesis ---

func TestGetIndustrialCapitalGenesis_Success(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForGenesis(t, st)

	_, _ = st.CreateCapitalOrigin(t.Context(), simulation.CapitalOrigin{
		HistoricalStageID: stageID,
		Source:            "colonial-plunder",
		AmountPence:       200000,
		Period:            "1500-1800",
	})
	_, _ = st.CreateColonialTransfer(t.Context(), simulation.ColonialTransfer{
		HistoricalStageID: stageID,
		From:              "West Africa",
		To:                "England",
		ValuePence:        15000,
		Method:            "slave-trade",
	})
	_, _ = st.CreateNationalDebt(t.Context(), simulation.NationalDebt{
		HistoricalStageID: stageID,
		AmountPence:       2400000,
		InterestRateBps:   800,
		CreditorClass:     "private-bankers",
	})

	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, CapitalOrigins: st, ColonialTransfers: st, NationalDebts: st, ProtectionSystems: st})
	req := httptest.NewRequest(http.MethodGet, "/v1/historical-stages/"+string(stageID)+"/genesis", nil)
	req.SetPathValue("id", string(stageID))
	w := httptest.NewRecorder()
	h.GetIndustrialCapitalGenesis(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		TotalCapitalFormedPence int64 `json:"total_capital_formed_pence"`
		Origins                 []struct {
			Source      string `json:"source"`
			AmountPence int64  `json:"amount_pence"`
		} `json:"origins"`
		ColonialTransfers []struct {
			ValuePence int64 `json:"value_pence"`
		} `json:"colonial_transfers"`
		NationalDebts []struct {
			InterestRateBps int64 `json:"interest_rate_bps"`
		} `json:"national_debts"`
		ProtectionSystems []struct{} `json:"protection_systems"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	// Invariant: total = sum(origins) + sum(transfers); national debts excluded
	want := int64(200000 + 15000)
	if resp.TotalCapitalFormedPence != want {
		t.Errorf("total_capital_formed_pence = %d, want %d", resp.TotalCapitalFormedPence, want)
	}
	if len(resp.Origins) != 1 {
		t.Errorf("origins len = %d, want 1", len(resp.Origins))
	}
	if len(resp.ColonialTransfers) != 1 {
		t.Errorf("colonial_transfers len = %d, want 1", len(resp.ColonialTransfers))
	}
	if len(resp.NationalDebts) != 1 {
		t.Errorf("national_debts len = %d, want 1", len(resp.NationalDebts))
	}
	if resp.ProtectionSystems == nil {
		t.Error("protection_systems should not be null")
	}
}

func TestGetIndustrialCapitalGenesis_NotFound(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, CapitalOrigins: st, ColonialTransfers: st, NationalDebts: st, ProtectionSystems: st})

	req := httptest.NewRequest(http.MethodGet, "/v1/historical-stages/nonexistent/genesis", nil)
	req.SetPathValue("id", "nonexistent")
	w := httptest.NewRecorder()
	h.GetIndustrialCapitalGenesis(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("want 404, got %d", w.Code)
	}
}

func TestGetIndustrialCapitalGenesis_EmptyStage(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	stageID := seedStageForGenesis(t, st)

	h := httpapi.New(nil, httpapi.Deps{HistoricalStages: st, CapitalOrigins: st, ColonialTransfers: st, NationalDebts: st, ProtectionSystems: st})
	req := httptest.NewRequest(http.MethodGet, "/v1/historical-stages/"+string(stageID)+"/genesis", nil)
	req.SetPathValue("id", string(stageID))
	w := httptest.NewRecorder()
	h.GetIndustrialCapitalGenesis(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp struct {
		TotalCapitalFormedPence int64         `json:"total_capital_formed_pence"`
		Origins                 []interface{} `json:"origins"`
		ColonialTransfers       []interface{} `json:"colonial_transfers"`
		NationalDebts           []interface{} `json:"national_debts"`
		ProtectionSystems       []interface{} `json:"protection_systems"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.TotalCapitalFormedPence != 0 {
		t.Errorf("empty stage total = %d, want 0", resp.TotalCapitalFormedPence)
	}
	if resp.Origins == nil {
		t.Error("origins should not be null")
	}
	if resp.ColonialTransfers == nil {
		t.Error("colonial_transfers should not be null")
	}
}
