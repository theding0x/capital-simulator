package engine

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"log/slog"
	"sync"
	"time"
)

// DefaultTickInterval is the period between scheduler passes when
// SIM_TICK_INTERVAL is unset or invalid.
const DefaultTickInterval = 5 * time.Second

// tickTimeout bounds a single scheduler pass. It is detached from the Stop
// signal so an in-flight pass finishes its writes even when Stop is called
// mid-tick (graceful shutdown), while still capping a wedged ticker.
const tickTimeout = 30 * time.Second

// NewScheduledTickID returns a 96-bit random hex identifier for a ScheduledTick.
func NewScheduledTickID() string {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

// Ticker is one registered unit of work the scheduler advances every pass.
// Tick reports how many domain entities it advanced and any error; the
// scheduler invokes registered tickers in a deterministic (registration) order.
type Ticker interface {
	// Name is a stable identifier surfaced in status and the audit log.
	Name() string
	// Tick advances this ticker's domain by one period, returning the number
	// of entities advanced. A ticker should advance as many entities as it can
	// and report a joined error rather than aborting the whole pass.
	Tick(ctx context.Context) (advanced int, err error)
}

// ScheduledTick is one scheduler pass written to the engine_ticks audit log.
// It records how many tickers ran, how many domain entities they advanced, how
// long the pass took, and how many tickers reported an error.
type ScheduledTick struct {
	ID               string    `json:"id"`
	Sequence         int64     `json:"sequence"`
	OccurredAt       time.Time `json:"occurred_at"`
	DurationMS       int64     `json:"duration_ms"`
	TickersRun       int       `json:"tickers_run"`
	EntitiesAdvanced int       `json:"entities_advanced"`
	ErrorCount       int       `json:"error_count"`
}

// SchedulerStatus is the operator-facing snapshot returned by GET
// /v1/engine/status: whether the loop is running, the current pass number, when
// the last pass occurred, the configured interval, and the registered tickers.
type SchedulerStatus struct {
	Running    bool       `json:"running"`
	Tick       int64      `json:"tick"`
	LastTickAt *time.Time `json:"last_tick_at,omitempty"`
	IntervalMS int64      `json:"interval_ms"`
	Tickers    []string   `json:"tickers"`
}

// TickRecorder persists one ScheduledTick per scheduler pass. The store layer
// implements it; defining it here keeps the engine package free of any store
// import (the store already imports engine).
type TickRecorder interface {
	RecordEngineTick(ctx context.Context, t ScheduledTick) (ScheduledTick, error)
}

// Scheduler is the automatic time-step orchestrator (issue #214). It advances
// the simulated economy one period at a time on a fixed interval, invoking each
// registered Ticker in order and writing an audit row per pass. It is started
// and stopped by operator control (POST /v1/engine/start, /stop) and stops
// gracefully on process shutdown.
type Scheduler struct {
	interval time.Duration
	tickers  []Ticker
	recorder TickRecorder
	logger   *slog.Logger
	now      func() time.Time

	mu       sync.Mutex
	running  bool
	sequence int64
	lastTick time.Time
	cancel   context.CancelFunc
	done     chan struct{}
}

// NewScheduler builds a Scheduler. A non-positive interval falls back to
// DefaultTickInterval; a nil logger falls back to slog.Default.
func NewScheduler(interval time.Duration, tickers []Ticker, recorder TickRecorder, logger *slog.Logger) *Scheduler {
	if interval <= 0 {
		interval = DefaultTickInterval
	}
	if logger == nil {
		logger = slog.Default()
	}
	return &Scheduler{
		interval: interval,
		tickers:  tickers,
		recorder: recorder,
		logger:   logger,
		now:      time.Now,
	}
}

// Start launches the periodic loop in a goroutine. It returns true if the loop
// was newly started, false if it was already running (idempotent).
func (s *Scheduler) Start() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return false
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	s.cancel = cancel
	s.done = done
	s.running = true
	go s.loop(ctx, done)
	s.logger.Info("tick scheduler started",
		"interval", s.interval.String(), "tickers", len(s.tickers))
	return true
}

// Stop signals the loop to halt and blocks until the current pass (if any)
// finishes — a graceful shutdown. It returns true if the loop was running,
// false if it was already stopped (idempotent).
func (s *Scheduler) Stop() bool {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return false
	}
	cancel := s.cancel
	done := s.done
	s.mu.Unlock()

	cancel()
	<-done

	s.mu.Lock()
	s.running = false
	s.cancel = nil
	s.done = nil
	s.mu.Unlock()
	s.logger.Info("tick scheduler stopped")
	return true
}

// Status returns a snapshot of the scheduler for operator inspection.
func (s *Scheduler) Status() SchedulerStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	names := make([]string, len(s.tickers))
	for i, t := range s.tickers {
		names[i] = t.Name()
	}
	st := SchedulerStatus{
		Running:    s.running,
		Tick:       s.sequence,
		IntervalMS: s.interval.Milliseconds(),
		Tickers:    names,
	}
	if !s.lastTick.IsZero() {
		lt := s.lastTick.UTC()
		st.LastTickAt = &lt
	}
	return st
}

// RunOnce executes a single scheduler pass: it invokes every registered ticker
// in order, sums the advanced counts and errors, then records one ScheduledTick
// audit row. It is exported so the loop and tests share one code path.
func (s *Scheduler) RunOnce(ctx context.Context) (ScheduledTick, error) {
	start := s.now()

	var advanced, errCount int
	for _, t := range s.tickers {
		n, err := t.Tick(ctx)
		advanced += n
		if err != nil {
			errCount++
			s.logger.Error("ticker pass failed", "ticker", t.Name(), "err", err)
		}
	}

	s.mu.Lock()
	s.sequence++
	seq := s.sequence
	s.lastTick = start
	s.mu.Unlock()

	tick := ScheduledTick{
		ID:               NewScheduledTickID(),
		Sequence:         seq,
		OccurredAt:       start.UTC(),
		DurationMS:       s.now().Sub(start).Milliseconds(),
		TickersRun:       len(s.tickers),
		EntitiesAdvanced: advanced,
		ErrorCount:       errCount,
	}
	recorded, err := s.recorder.RecordEngineTick(ctx, tick)
	if err != nil {
		s.logger.Error("failed to record engine tick", "sequence", seq, "err", err)
		return tick, err
	}
	return recorded, nil
}

// loop runs RunOnce every interval until ctx is cancelled. Each pass runs on a
// context detached from ctx so a Stop mid-pass still lets the pass finish.
func (s *Scheduler) loop(ctx context.Context, done chan struct{}) {
	defer close(done)
	t := time.NewTicker(s.interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			tickCtx, cancel := context.WithTimeout(context.Background(), tickTimeout)
			if _, err := s.RunOnce(tickCtx); err != nil {
				s.logger.Error("scheduler pass error", "err", err)
			}
			cancel()
		}
	}
}
