package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

func newEngineTestServer(t *testing.T, d Deps) *httptest.Server {
	t.Helper()
	h := New(nil, d)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/engine/start", h.StartEngine)
	mux.HandleFunc("POST /v1/engine/stop", h.StopEngine)
	mux.HandleFunc("GET /v1/engine/status", h.EngineStatus)
	mux.HandleFunc("GET /v1/engine/ticks", h.ListEngineTicks)
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)
	return ts
}

func TestEngineStatus_BeforeStart(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	// A long interval keeps the loop from firing during the test.
	sched := engine.NewScheduler(time.Hour, []engine.Ticker{
		engine.NewFactoryTicker(st), engine.NewReproductionTicker(st),
	}, st, nil)
	ts := newEngineTestServer(t, Deps{Scheduler: sched, EngineTicks: st})

	res, err := http.Get(ts.URL + "/v1/engine/status")
	if err != nil {
		t.Fatalf("GET status: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var got engine.SchedulerStatus
	if err := json.NewDecoder(res.Body).Decode(&got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Running {
		t.Fatal("running should be false before start")
	}
	if len(got.Tickers) != 2 {
		t.Fatalf("tickers = %v, want 2 registered", got.Tickers)
	}
}

func TestEngineStartStop(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	sched := engine.NewScheduler(time.Hour, nil, st, nil)
	ts := newEngineTestServer(t, Deps{Scheduler: sched, EngineTicks: st})
	t.Cleanup(func() { sched.Stop() })

	res, err := http.Post(ts.URL+"/v1/engine/start", "application/json", nil)
	if err != nil {
		t.Fatalf("POST start: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("start status = %d, want 200", res.StatusCode)
	}
	var startResp engineControlResponse
	if err := json.NewDecoder(res.Body).Decode(&startResp); err != nil {
		t.Fatalf("decode start: %v", err)
	}
	if !startResp.Started {
		t.Fatal("started should be true on first start")
	}
	if !startResp.Status.Running {
		t.Fatal("status.running should be true after start")
	}

	res2, err := http.Post(ts.URL+"/v1/engine/stop", "application/json", nil)
	if err != nil {
		t.Fatalf("POST stop: %v", err)
	}
	defer res2.Body.Close()
	var stopResp engineControlResponse
	if err := json.NewDecoder(res2.Body).Decode(&stopResp); err != nil {
		t.Fatalf("decode stop: %v", err)
	}
	if !stopResp.Stopped {
		t.Fatal("stopped should be true")
	}
	if stopResp.Status.Running {
		t.Fatal("status.running should be false after stop")
	}
}

func TestEngineStatus_NilSchedulerIs500(t *testing.T) {
	t.Parallel()
	ts := newEngineTestServer(t, Deps{}) // no scheduler wired
	res, err := http.Get(ts.URL + "/v1/engine/status")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", res.StatusCode)
	}
}

func TestListEngineTicks_BadLimit(t *testing.T) {
	t.Parallel()
	ts := newEngineTestServer(t, Deps{EngineTicks: store.NewMemory()})
	res, err := http.Get(ts.URL + "/v1/engine/ticks?limit=-1")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", res.StatusCode)
	}
}

func TestListEngineTicks_EmptyArrayNotNull(t *testing.T) {
	t.Parallel()
	ts := newEngineTestServer(t, Deps{EngineTicks: store.NewMemory()})
	res, err := http.Get(ts.URL + "/v1/engine/ticks")
	if err != nil {
		t.Fatalf("GET: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", res.StatusCode)
	}
	var ticks []engine.ScheduledTick
	if err := json.NewDecoder(res.Body).Decode(&ticks); err != nil {
		t.Fatalf("decode (should be [] not null): %v", err)
	}
	if ticks == nil {
		t.Fatal("expected an empty [] array, got null")
	}
}
