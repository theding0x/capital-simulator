package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func newMachineryTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	st := store.NewMemory()
	h := New(nil, Deps{Machines: st, Factories: st})
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/machines", h.CreateMachine)
	mux.HandleFunc("GET /v1/machines", h.ListMachines)
	mux.HandleFunc("GET /v1/machines/{id}", h.GetMachine)
	mux.HandleFunc("GET /v1/machines/{id}/wear", h.GetMachineWear)
	mux.HandleFunc("POST /v1/factories", h.CreateFactory)
	mux.HandleFunc("GET /v1/factories/{id}", h.GetFactory)
	mux.HandleFunc("POST /v1/factories/{id}/tick", h.TickFactory)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

// §2 fixture: machine_value=600000, lifespan_days=1000 → daily_wear_and_tear=600
func TestCreateMachine_FixtureDailyWearAndTear(t *testing.T) {
	t.Parallel()
	ts := newMachineryTestServer(t)
	body := `{"name":"power-loom","machine_value":600000,"lifespan_days":1000,"productive_power":145000}`
	res, err := http.Post(ts.URL+"/v1/machines", "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", res.StatusCode)
	}
	var got machineDTO
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.DailyWearAndTear != 600 {
		t.Fatalf("DailyWearAndTear = %d, want 600", got.DailyWearAndTear)
	}
	if got.ID == "" {
		t.Fatal("expected server-assigned id")
	}
}

func TestCreateMachine_RejectsZeroLifespan(t *testing.T) {
	t.Parallel()
	ts := newMachineryTestServer(t)
	body := `{"name":"bad","machine_value":1000,"lifespan_days":0,"productive_power":10}`
	res, err := http.Post(ts.URL+"/v1/machines", "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}

// §8A fixture: one needle-machine = 145,000 units/day; four under one
// supervisor produce ~580,000 units/day. A factory tick yields the aggregate.
func TestTickFactory_Para8A_NeedleFactory(t *testing.T) {
	t.Parallel()
	ts := newMachineryTestServer(t)

	// Create the factory with four embedded needle machines.
	body := `{
		"name": "Needle Works",
		"prime_mover": {"kind":"steam","horsepower":30},
		"machines": [
			{"name":"needle-1","machine_value":600000,"lifespan_days":1000,"productive_power":145000},
			{"name":"needle-2","machine_value":600000,"lifespan_days":1000,"productive_power":145000},
			{"name":"needle-3","machine_value":600000,"lifespan_days":1000,"productive_power":145000},
			{"name":"needle-4","machine_value":600000,"lifespan_days":1000,"productive_power":145000}
		]
	}`
	res, err := http.Post(ts.URL+"/v1/factories", "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatalf("POST factory: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("create factory status = %d", res.StatusCode)
	}
	var f factoryDTO
	if err := json.NewDecoder(res.Body).Decode(&f); err != nil {
		t.Fatalf("decode factory: %v", err)
	}
	if f.TotalProductivePower != 580_000 {
		t.Fatalf("TotalProductivePower = %d, want 580000", f.TotalProductivePower)
	}

	// Advance one tick.
	tickRes, err := http.Post(ts.URL+"/v1/factories/"+f.ID+"/tick", "application/json", nil)
	if err != nil {
		t.Fatalf("POST tick: %v", err)
	}
	defer tickRes.Body.Close()
	if tickRes.StatusCode != http.StatusOK {
		t.Fatalf("tick status = %d, want 200", tickRes.StatusCode)
	}
	var tr tickResponse
	if err := json.NewDecoder(tickRes.Body).Decode(&tr); err != nil {
		t.Fatalf("decode tick: %v", err)
	}
	if tr.UnitsProduced != 580_000 {
		t.Fatalf("UnitsProduced = %d, want 580000", tr.UnitsProduced)
	}
	if tr.ValueTransferred != 2400 { // 4 × 600
		t.Fatalf("ValueTransferred = %d, want 2400", tr.ValueTransferred)
	}
	if tr.Factory.TickCount != 1 {
		t.Fatalf("TickCount = %d, want 1", tr.Factory.TickCount)
	}
}

func TestGetMachineWear_ReturnsRemainingValue(t *testing.T) {
	t.Parallel()
	ts := newMachineryTestServer(t)
	body := `{"name":"power-loom","machine_value":600000,"lifespan_days":1000,"productive_power":145000}`
	res, err := http.Post(ts.URL+"/v1/machines", "application/json", bytes.NewReader([]byte(body)))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer res.Body.Close()
	var m machineDTO
	json.NewDecoder(res.Body).Decode(&m) //nolint
	wearRes, err := http.Get(ts.URL + "/v1/machines/" + m.ID + "/wear")
	if err != nil {
		t.Fatalf("GET wear: %v", err)
	}
	defer wearRes.Body.Close()
	var wr machineWearResponse
	if err := json.NewDecoder(wearRes.Body).Decode(&wr); err != nil {
		t.Fatalf("decode wear: %v", err)
	}
	if wr.RemainingValue != 600000 {
		t.Fatalf("RemainingValue = %d, want 600000 (fresh machine)", wr.RemainingValue)
	}
}

func TestGetFactory_NotFound(t *testing.T) {
	t.Parallel()
	ts := newMachineryTestServer(t)
	res, err := http.Get(ts.URL + "/v1/factories/does-not-exist")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", res.StatusCode)
	}
}
