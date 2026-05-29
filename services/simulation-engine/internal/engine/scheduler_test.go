package engine_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

var errTickerBoom = errors.New("ticker boom")

// recordingTicker is a test Ticker that records the order in which it was
// invoked, optionally signals a channel each pass, and returns a fixed result.
type recordingTicker struct {
	name     string
	advanced int
	err      error
	order    *[]string
	mu       *sync.Mutex
	fired    chan struct{}
}

func (t *recordingTicker) Name() string { return t.name }

func (t *recordingTicker) Tick(_ context.Context) (int, error) {
	if t.order != nil {
		t.mu.Lock()
		*t.order = append(*t.order, t.name)
		t.mu.Unlock()
	}
	if t.fired != nil {
		select {
		case t.fired <- struct{}{}:
		default:
		}
	}
	return t.advanced, t.err
}

// fakeRecorder is an in-memory engine.TickRecorder.
type fakeRecorder struct {
	mu    sync.Mutex
	ticks []engine.ScheduledTick
}

func (r *fakeRecorder) RecordEngineTick(_ context.Context, t engine.ScheduledTick) (engine.ScheduledTick, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if t.ID == "" {
		t.ID = engine.NewScheduledTickID()
	}
	r.ticks = append(r.ticks, t)
	return t, nil
}

func (r *fakeRecorder) count() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.ticks)
}

func TestScheduler_RunOnce_InvokesTickersInOrderAndAggregates(t *testing.T) {
	t.Parallel()
	var order []string
	var mu sync.Mutex
	tickers := []engine.Ticker{
		&recordingTicker{name: "factory", advanced: 3, order: &order, mu: &mu},
		&recordingTicker{name: "reproduction", advanced: 2, order: &order, mu: &mu},
	}
	rec := &fakeRecorder{}
	s := engine.NewScheduler(time.Hour, tickers, rec, nil)

	got, err := s.RunOnce(context.Background())
	if err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	if want := []string{"factory", "reproduction"}; !equalStrings(order, want) {
		t.Fatalf("invocation order = %v, want %v", order, want)
	}
	if got.Sequence != 1 {
		t.Fatalf("sequence = %d, want 1", got.Sequence)
	}
	if got.TickersRun != 2 {
		t.Fatalf("tickers_run = %d, want 2", got.TickersRun)
	}
	if got.EntitiesAdvanced != 5 {
		t.Fatalf("entities_advanced = %d, want 5", got.EntitiesAdvanced)
	}
	if got.ErrorCount != 0 {
		t.Fatalf("error_count = %d, want 0", got.ErrorCount)
	}
	if rec.count() != 1 {
		t.Fatalf("recorded ticks = %d, want 1", rec.count())
	}
}

func TestScheduler_RunOnce_CountsTickerErrors(t *testing.T) {
	t.Parallel()
	tickers := []engine.Ticker{
		&recordingTicker{name: "ok", advanced: 1},
		&recordingTicker{name: "bad", advanced: 0, err: errTickerBoom},
	}
	rec := &fakeRecorder{}
	s := engine.NewScheduler(time.Hour, tickers, rec, nil)

	got, err := s.RunOnce(context.Background())
	if err != nil {
		t.Fatalf("RunOnce should not surface a ticker error as its own error: %v", err)
	}
	if got.ErrorCount != 1 {
		t.Fatalf("error_count = %d, want 1", got.ErrorCount)
	}
	if got.EntitiesAdvanced != 1 {
		t.Fatalf("entities_advanced = %d, want 1", got.EntitiesAdvanced)
	}
}

func TestScheduler_StartStop_Idempotent(t *testing.T) {
	t.Parallel()
	s := engine.NewScheduler(time.Hour, nil, &fakeRecorder{}, nil)

	if !s.Start() {
		t.Fatal("first Start should return true")
	}
	if s.Start() {
		t.Fatal("second Start should return false (already running)")
	}
	if !s.Stop() {
		t.Fatal("first Stop should return true")
	}
	if s.Stop() {
		t.Fatal("second Stop should return false (already stopped)")
	}
}

func TestScheduler_Loop_TicksThenStopsGracefully(t *testing.T) {
	t.Parallel()
	fired := make(chan struct{}, 64)
	tickers := []engine.Ticker{
		&recordingTicker{name: "factory", advanced: 1, fired: fired},
	}
	rec := &fakeRecorder{}
	s := engine.NewScheduler(2*time.Millisecond, tickers, rec, nil)

	if !s.Start() {
		t.Fatal("Start should return true")
	}
	select {
	case <-fired:
	case <-time.After(2 * time.Second):
		t.Fatal("scheduler did not fire a tick within 2s")
	}
	if !s.Stop() {
		t.Fatal("Stop should return true")
	}

	// After Stop returns, the loop goroutine has exited — no further ticks.
	settled := rec.count()
	if settled == 0 {
		t.Fatal("expected at least one recorded tick")
	}
	time.Sleep(20 * time.Millisecond)
	if grew := rec.count(); grew != settled {
		t.Fatalf("ticks recorded after Stop: %d -> %d", settled, grew)
	}
}

func TestScheduler_Status(t *testing.T) {
	t.Parallel()
	tickers := []engine.Ticker{
		&recordingTicker{name: "factory", advanced: 1},
		&recordingTicker{name: "reproduction", advanced: 1},
	}
	s := engine.NewScheduler(250*time.Millisecond, tickers, &fakeRecorder{}, nil)

	st := s.Status()
	if st.Running {
		t.Fatal("Running should be false before Start")
	}
	if st.Tick != 0 {
		t.Fatalf("Tick = %d, want 0", st.Tick)
	}
	if st.LastTickAt != nil {
		t.Fatal("LastTickAt should be nil before any pass")
	}
	if st.IntervalMS != 250 {
		t.Fatalf("IntervalMS = %d, want 250", st.IntervalMS)
	}
	if want := []string{"factory", "reproduction"}; !equalStrings(st.Tickers, want) {
		t.Fatalf("Tickers = %v, want %v", st.Tickers, want)
	}

	if _, err := s.RunOnce(context.Background()); err != nil {
		t.Fatalf("RunOnce: %v", err)
	}
	st = s.Status()
	if st.Tick != 1 {
		t.Fatalf("Tick after one pass = %d, want 1", st.Tick)
	}
	if st.LastTickAt == nil {
		t.Fatal("LastTickAt should be set after a pass")
	}
}

func TestScheduler_DefaultInterval(t *testing.T) {
	t.Parallel()
	s := engine.NewScheduler(0, nil, &fakeRecorder{}, nil)
	if got, want := s.Status().IntervalMS, engine.DefaultTickInterval.Milliseconds(); got != want {
		t.Fatalf("IntervalMS = %d, want default %d", got, want)
	}
}

// TestMemorySatisfiesSchedulerContracts guards the structural contract that the
// real store satisfies the narrow interfaces the scheduler and tickers depend
// on. If a method signature drifts, this fails at compile time.
func TestMemorySatisfiesSchedulerContracts(t *testing.T) {
	t.Parallel()
	var st any = store.NewMemory()
	if _, ok := st.(engine.TickRecorder); !ok {
		t.Error("*store.Memory does not satisfy engine.TickRecorder")
	}
	if _, ok := st.(engine.FactoryAdvancer); !ok {
		t.Error("*store.Memory does not satisfy engine.FactoryAdvancer")
	}
	if _, ok := st.(engine.ReproductionAdvancer); !ok {
		t.Error("*store.Memory does not satisfy engine.ReproductionAdvancer")
	}
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
