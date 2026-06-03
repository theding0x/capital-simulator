package engine_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
)

var _ engine.Ticker = (*engine.PiecePriceTicker)(nil)

type fakeProductivity struct {
	factor float64
	err    error
}

func (f fakeProductivity) LeadingProductivityFactor(context.Context) (float64, error) {
	return f.factor, f.err
}

type fakeRepricer struct {
	gotFactor float64
	advanced  int
	err       error
	calls     int
}

func (f *fakeRepricer) Reprice(_ context.Context, factor float64) (int, error) {
	f.calls++
	f.gotFactor = factor
	return f.advanced, f.err
}

func TestPiecePriceTicker_PushesFactorAndReturnsAdvanced(t *testing.T) {
	t.Parallel()
	rp := &fakeRepricer{advanced: 3}
	ticker := engine.NewPiecePriceTicker(fakeProductivity{factor: 2.0}, rp)

	if ticker.Name() != "piece-price" {
		t.Errorf("name = %q, want piece-price", ticker.Name())
	}
	got, err := ticker.Tick(context.Background())
	if err != nil {
		t.Fatalf("tick: %v", err)
	}
	if got != 3 {
		t.Errorf("advanced = %d, want 3", got)
	}
	if rp.gotFactor != 2.0 {
		t.Errorf("repricer received factor %v, want 2.0", rp.gotFactor)
	}
}

func TestPiecePriceTicker_ProductivityErrorShortCircuits(t *testing.T) {
	t.Parallel()
	rp := &fakeRepricer{advanced: 9}
	sentinel := errors.New("boom")
	ticker := engine.NewPiecePriceTicker(fakeProductivity{err: sentinel}, rp)

	got, err := ticker.Tick(context.Background())
	if !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want sentinel", err)
	}
	if got != 0 {
		t.Errorf("advanced = %d, want 0 when productivity read fails", got)
	}
	if rp.calls != 0 {
		t.Errorf("repricer called %d times, want 0 (short-circuit)", rp.calls)
	}
}

func TestPiecePriceTicker_RepricerErrorPropagates(t *testing.T) {
	t.Parallel()
	sentinel := errors.New("reprice failed")
	rp := &fakeRepricer{err: sentinel}
	ticker := engine.NewPiecePriceTicker(fakeProductivity{factor: 1.5}, rp)

	if _, err := ticker.Tick(context.Background()); !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want sentinel", err)
	}
}

type fakeFactories struct {
	factories []machinery.Factory
	err       error
}

func (f fakeFactories) ListFactories(context.Context) ([]machinery.Factory, error) {
	return f.factories, f.err
}

// With no factories — or only factories at the hand-labour baseline — the
// leading factor floors at 1.0, so the piece price is left unchanged.
func TestFactoryProductivitySource_FloorsAtOne(t *testing.T) {
	t.Parallel()
	src := engine.NewFactoryProductivitySource(fakeFactories{})
	got, err := src.LeadingProductivityFactor(context.Background())
	if err != nil {
		t.Fatalf("leading: %v", err)
	}
	if got != 1.0 {
		t.Errorf("factor = %v, want 1.0", got)
	}
}

func TestFactoryProductivitySource_ErrorPropagates(t *testing.T) {
	t.Parallel()
	sentinel := errors.New("store down")
	src := engine.NewFactoryProductivitySource(fakeFactories{err: sentinel})
	if _, err := src.LeadingProductivityFactor(context.Background()); !errors.Is(err, sentinel) {
		t.Fatalf("err = %v, want sentinel", err)
	}
}
