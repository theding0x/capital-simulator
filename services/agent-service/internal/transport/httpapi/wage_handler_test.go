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

// § core fixture from Ch. 19:
// daily_pence=36 (3s), working_day_minutes=720, lpv_daily_pence=36, necessary_minutes=360
// → hourly_wage=3, appearance={paid:12,unpaid:0}, decomposition={paid:360,unpaid:360}

func TestCreateWageForm(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)

	body := map[string]any{
		"agent_id":          "abc123",
		"daily_pence":       36,
		"working_day_minutes": 720,
		"lpv_daily_pence":   36,
		"necessary_minutes": 360,
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/wage-forms", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateWageForm(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body: %s", rr.Code, http.StatusCreated, rr.Body.String())
	}

	var resp struct {
		agent.WageForm
		HourlyWage    agent.Pence               `json:"hourly_wage"`
		Appearance    agent.WageAppearance       `json:"appearance"`
		Decomposition agent.LabourDecomposition  `json:"decomposition"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.HourlyWage != 3 {
		t.Errorf("hourly_wage = %d, want 3", resp.HourlyWage)
	}
	if resp.Appearance.PaidHours != 12 || resp.Appearance.UnpaidHours != 0 {
		t.Errorf("appearance = {%d,%d}, want {12,0}", resp.Appearance.PaidHours, resp.Appearance.UnpaidHours)
	}
	if resp.Decomposition.PaidMinutes != 360 || resp.Decomposition.UnpaidMinutes != 360 {
		t.Errorf("decomposition = {%d,%d}, want {360,360}",
			resp.Decomposition.PaidMinutes, resp.Decomposition.UnpaidMinutes)
	}
}

func TestGetWageForm(t *testing.T) {
	t.Parallel()

	mem := store.NewMemory()
	h := httpapi.New(mem, nil)

	// seed
	wf := agent.WageForm{
		AgentID: "xyz789",
		Wage:    agent.Wage{DailyPence: 24, WorkingDayMinutes: 720},
		LabourPowerValue: agent.WageLabourValue{
			DailyPence:       24,
			NecessaryMinutes: 360,
		},
	}
	if _, err := mem.CreateWageForm(t.Context(), wf); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/wage-forms/xyz789", nil)
	req.SetPathValue("agentID", "xyz789")
	rr := httptest.NewRecorder()
	h.GetWageForm(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", rr.Code, http.StatusOK, rr.Body.String())
	}

	var resp struct {
		HourlyWage agent.Pence `json:"hourly_wage"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.HourlyWage != 2 {
		t.Errorf("hourly_wage = %d, want 2 (2s/12h)", resp.HourlyWage)
	}
}

func TestGetWageForm_NotFound(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/wage-forms/nobody", nil)
	req.SetPathValue("agentID", "nobody")
	rr := httptest.NewRecorder()
	h.GetWageForm(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
}

func TestCreateWageForm_InvalidRequest(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)

	t.Run("missing agent_id", func(t *testing.T) {
		t.Parallel()
		body := map[string]any{
			"daily_pence":       36,
			"working_day_minutes": 720,
			"lpv_daily_pence":   36,
			"necessary_minutes": 360,
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/v1/wage-forms", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.CreateWageForm(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})

	t.Run("zero daily_pence", func(t *testing.T) {
		t.Parallel()
		body := map[string]any{
			"agent_id":          "abc",
			"daily_pence":       0,
			"working_day_minutes": 720,
			"lpv_daily_pence":   36,
			"necessary_minutes": 360,
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/v1/wage-forms", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.CreateWageForm(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})
}
