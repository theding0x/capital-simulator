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

// § unit measure: daily 36p, 12h → numerator=36, denominator=12, as_float=3.0
func TestComputeHourlyPrice_OK(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)
	body := map[string]any{"daily_pence": 36, "working_day_hours": 12}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/time-wages/hourly-price", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.ComputeHourlyPrice(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rr.Code, rr.Body.String())
	}
	var resp struct {
		Numerator   int64   `json:"numerator"`
		Denominator int64   `json:"denominator"`
		AsFloat     float64 `json:"as_float"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Numerator != 36 || resp.Denominator != 12 {
		t.Errorf("fraction = %d/%d, want 36/12", resp.Numerator, resp.Denominator)
	}
	if resp.AsFloat != 3.0 {
		t.Errorf("as_float = %f, want 3.0", resp.AsFloat)
	}
}

func TestComputeHourlyPrice_InvalidRequest(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)

	t.Run("bad json", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodPost, "/v1/time-wages/hourly-price", bytes.NewReader([]byte("{bad}")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.ComputeHourlyPrice(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})

	t.Run("zero daily_pence", func(t *testing.T) {
		t.Parallel()
		body := map[string]any{"daily_pence": 0, "working_day_hours": 12}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/v1/time-wages/hourly-price", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.ComputeHourlyPrice(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})
}

// § 12-hour day, no overtime: nominal wage = 36p.
func TestCreateWorkingSession_OK(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)
	body := map[string]any{
		"agent_id":                 "worker-abc",
		"daily_labour_power_value": 36,
		"working_day_hours":        12,
		"overtime_hours":           0,
		"overtime_rate_pence":      0,
		"wage_period":              "daily",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/time-wages/sessions", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateWorkingSession(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", rr.Code, rr.Body.String())
	}
	if rr.Header().Get("Location") == "" {
		t.Error("Location header missing")
	}

	var resp struct {
		agent.WorkingSession
		HourlyPrice struct {
			Numerator   int64   `json:"numerator"`
			Denominator int64   `json:"denominator"`
			AsFloat     float64 `json:"as_float"`
		} `json:"hourly_price"`
		NominalWage struct {
			Pence int64 `json:"pence"`
		} `json:"nominal_wage"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.HourlyPrice.Numerator != 36 || resp.HourlyPrice.Denominator != 12 {
		t.Errorf("hourly fraction = %d/%d, want 36/12", resp.HourlyPrice.Numerator, resp.HourlyPrice.Denominator)
	}
	if resp.HourlyPrice.AsFloat != 3.0 {
		t.Errorf("hourly as_float = %f, want 3.0", resp.HourlyPrice.AsFloat)
	}
	if resp.NominalWage.Pence != 36 {
		t.Errorf("nominal_wage = %d, want 36", resp.NominalWage.Pence)
	}
}

// § overtime fixture: 12h + 2 OT at 4d → nominal = 44p.
func TestCreateWorkingSession_WithOvertime(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)
	body := map[string]any{
		"agent_id":                 "worker-abc",
		"daily_labour_power_value": 36,
		"working_day_hours":        12,
		"overtime_hours":           2,
		"overtime_rate_pence":      4,
		"wage_period":              "daily",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/time-wages/sessions", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	h.CreateWorkingSession(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body: %s", rr.Code, rr.Body.String())
	}
	var resp struct {
		NominalWage struct {
			Pence int64 `json:"pence"`
		} `json:"nominal_wage"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.NominalWage.Pence != 44 {
		t.Errorf("nominal_wage = %d, want 44 (36 + 2×4)", resp.NominalWage.Pence)
	}
}

func TestCreateWorkingSession_BadRequest(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)

	t.Run("bad json", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodPost, "/v1/time-wages/sessions", bytes.NewReader([]byte("{bad}")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.CreateWorkingSession(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})

	t.Run("missing agent_id", func(t *testing.T) {
		t.Parallel()
		body := map[string]any{
			"daily_labour_power_value": 36,
			"working_day_hours":        12,
			"wage_period":              "daily",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/v1/time-wages/sessions", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.CreateWorkingSession(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})

	t.Run("invalid wage_period", func(t *testing.T) {
		t.Parallel()
		body := map[string]any{
			"agent_id":                 "w1",
			"daily_labour_power_value": 36,
			"working_day_hours":        12,
			"wage_period":              "monthly",
		}
		b, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/v1/time-wages/sessions", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.CreateWorkingSession(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})
}

func TestGetWorkingSession_OK(t *testing.T) {
	t.Parallel()

	mem := store.NewMemory()
	h := httpapi.New(mem, nil)

	seed := agent.WorkingSession{
		AgentID:               "worker-xyz",
		DailyLabourPowerValue: agent.DailyLabourPowerValue{Pence: 36},
		WorkingDayHours:       agent.WorkingDayHours{Hours: 12},
		OvertimeHours:         agent.OvertimeHours{Hours: 0},
		OvertimeRatePence:     agent.OvertimeRatePence{Pence: 0},
		WagePeriod:            agent.WagePeriodDaily,
	}
	saved, err := mem.CreateWorkingSession(t.Context(), seed)
	if err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/time-wages/sessions/"+string(saved.ID), nil)
	req.SetPathValue("id", string(saved.ID))
	rr := httptest.NewRecorder()
	h.GetWorkingSession(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body: %s", rr.Code, rr.Body.String())
	}
	var resp struct {
		NominalWage struct {
			Pence int64 `json:"pence"`
		} `json:"nominal_wage"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.NominalWage.Pence != 36 {
		t.Errorf("nominal_wage = %d, want 36", resp.NominalWage.Pence)
	}
}

func TestGetWorkingSession_NotFound(t *testing.T) {
	t.Parallel()

	h := httpapi.New(store.NewMemory(), nil)
	req := httptest.NewRequest(http.MethodGet, "/v1/time-wages/sessions/nobody", nil)
	req.SetPathValue("id", "nobody")
	rr := httptest.NewRecorder()
	h.GetWorkingSession(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rr.Code)
	}
}
