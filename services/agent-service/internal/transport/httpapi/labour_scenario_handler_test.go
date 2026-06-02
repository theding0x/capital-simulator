package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

// newScenarioTestServer wires only the Ch. 17 endpoint. The compute is
// stateless, so the store is unused — but the Handler constructor still
// requires it.
func newScenarioTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	h := New(store.NewMemory(), slog.Default())
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/labour-scenarios", h.ComputeLabourScenario)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// Capital Vol. I, Ch. 17 §1: working day 12 h (720 min), necessary
// labour 6 h (360 min), normal intensity → daily value 720 min, surplus
// 360 min, rate of surplus-value 100%.
func TestComputeLabourScenario_S1NormalIntensity(t *testing.T) {
	t.Parallel()
	ts := newScenarioTestServer(t)
	body := `{
		"working_day_minutes": 720,
		"necessary_labour_minutes": 360,
		"intensity_factor": 1.0,
		"productivity_factor": 1.0
	}`
	res, err := http.Post(ts.URL+"/v1/labour-scenarios", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got labourScenarioResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.DailyValueMinutes != 720 {
		t.Errorf("daily value: got %d, want 720", got.DailyValueMinutes)
	}
	if got.SurplusLabourMinutes != 360 {
		t.Errorf("surplus: got %d, want 360", got.SurplusLabourMinutes)
	}
	if got.RateOfSurplusValue != 1.0 {
		t.Errorf("rate: got %v, want 1.0", got.RateOfSurplusValue)
	}
	if !got.LawConstantDailyValue {
		t.Error("§1 Law 1 must hold at normal intensity")
	}
}

// Ch. 17 §2: intensified labour with factor 4/3 packs 960 min of value
// into the same 720-min day.
func TestComputeLabourScenario_S2IntensityScalesValue(t *testing.T) {
	t.Parallel()
	ts := newScenarioTestServer(t)
	body := `{
		"working_day_minutes": 720,
		"necessary_labour_minutes": 360,
		"intensity_factor": 1.333333333,
		"productivity_factor": 1.0
	}`
	res, err := http.Post(ts.URL+"/v1/labour-scenarios", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	var got labourScenarioResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.DailyValueMinutes < 950 || got.DailyValueMinutes > 970 {
		t.Errorf("daily value at 4/3 intensity: got %d, want ~960", got.DailyValueMinutes)
	}
	if got.LawConstantDailyValue {
		t.Error("§1 Law 1 must not hold when intensity != 1.0")
	}
}

func TestComputeLabourScenario_RejectsNecessaryExceedingDay(t *testing.T) {
	t.Parallel()
	ts := newScenarioTestServer(t)
	body := `{"working_day_minutes":720,"necessary_labour_minutes":800,"intensity_factor":1.0,"productivity_factor":1.0}`
	res, err := http.Post(ts.URL+"/v1/labour-scenarios", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}

func TestComputeLabourScenario_DefaultsFactorsToOne(t *testing.T) {
	t.Parallel()
	ts := newScenarioTestServer(t)
	body := `{"working_day_minutes":720,"necessary_labour_minutes":360}` // omit factors
	res, err := http.Post(ts.URL+"/v1/labour-scenarios", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got labourScenarioResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.DailyValueMinutes != 720 {
		t.Errorf("default factors should yield daily value 720, got %d", got.DailyValueMinutes)
	}
}

// Ch. 25 → Ch. 17: with a quarter of the workforce in the reserve army, the
// wage is bid down to 270 min, below the 360-min value of labour-power.
func TestComputeLabourScenario_ReserveArmyCompression(t *testing.T) {
	t.Parallel()
	ts := newScenarioTestServer(t)
	body := `{
		"working_day_minutes": 720,
		"necessary_labour_minutes": 360,
		"intensity_factor": 1.0,
		"productivity_factor": 1.0,
		"reserve_army_magnitude": 250000,
		"reserve_army_workforce_total": 1000000
	}`
	res, err := http.Post(ts.URL+"/v1/labour-scenarios", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got labourScenarioResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.LabourPowerValueMinutes != 360 {
		t.Errorf("value of labour-power: got %d, want 360", got.LabourPowerValueMinutes)
	}
	if got.ReserveArmyPressureBP != 2500 {
		t.Errorf("pressure: got %d bp, want 2500", got.ReserveArmyPressureBP)
	}
	if got.CompressedWageMinutes != 270 {
		t.Errorf("compressed wage: got %d, want 270", got.CompressedWageMinutes)
	}
	if got.CompressedWageMinutes >= got.LabourPowerValueMinutes {
		t.Error("reserve army must push the wage below the value of labour-power")
	}
}

// An empty (no overlay) scenario must not report any compression.
func TestComputeLabourScenario_NoReserveArmyLeavesWageAtValue(t *testing.T) {
	t.Parallel()
	ts := newScenarioTestServer(t)
	body := `{"working_day_minutes":720,"necessary_labour_minutes":360}`
	res, err := http.Post(ts.URL+"/v1/labour-scenarios", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	var got labourScenarioResponse
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.ReserveArmyPressureBP != 0 || got.CompressedWageMinutes != 0 {
		t.Errorf("no overlay should leave compression zero, got bp=%d wage=%d", got.ReserveArmyPressureBP, got.CompressedWageMinutes)
	}
}

func TestComputeLabourScenario_RejectsReserveArmyExceedingWorkforce(t *testing.T) {
	t.Parallel()
	ts := newScenarioTestServer(t)
	body := `{
		"working_day_minutes": 720,
		"necessary_labour_minutes": 360,
		"reserve_army_magnitude": 1200000,
		"reserve_army_workforce_total": 1000000
	}`
	res, err := http.Post(ts.URL+"/v1/labour-scenarios", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}
