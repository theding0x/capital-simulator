package engine_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
)

var errFactoryBoom = errors.New("advance boom")

type fakeFactoryStore struct {
	factories []machinery.Factory
	advanced  map[machinery.FactoryID]int
	failOn    machinery.FactoryID
	listErr   error
}

func (f *fakeFactoryStore) ListFactories(context.Context) ([]machinery.Factory, error) {
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.factories, nil
}

func (f *fakeFactoryStore) AdvanceTick(_ context.Context, id machinery.FactoryID) (machinery.Factory, engine.Tick, error) {
	if id == f.failOn {
		return machinery.Factory{}, engine.Tick{}, errFactoryBoom
	}
	if f.advanced == nil {
		f.advanced = map[machinery.FactoryID]int{}
	}
	f.advanced[id]++
	return machinery.Factory{ID: id}, engine.Tick{FactoryID: string(id)}, nil
}

func TestFactoryTicker_AdvancesEveryFactory(t *testing.T) {
	t.Parallel()
	fs := &fakeFactoryStore{
		factories: []machinery.Factory{
			{ID: "needle-works"}, {ID: "spinning-mill"}, {ID: "power-loom"},
		},
		advanced: map[machinery.FactoryID]int{},
	}
	tk := engine.NewFactoryTicker(fs)
	if tk.Name() != "factory" {
		t.Fatalf("Name = %q, want factory", tk.Name())
	}
	advanced, err := tk.Tick(context.Background())
	if err != nil {
		t.Fatalf("Tick: %v", err)
	}
	if advanced != 3 {
		t.Fatalf("advanced = %d, want 3", advanced)
	}
	for _, id := range []machinery.FactoryID{"needle-works", "spinning-mill", "power-loom"} {
		if fs.advanced[id] != 1 {
			t.Errorf("factory %s advanced %d times, want 1", id, fs.advanced[id])
		}
	}
}

func TestFactoryTicker_ContinuesPastFailure(t *testing.T) {
	t.Parallel()
	fs := &fakeFactoryStore{
		factories: []machinery.Factory{
			{ID: "good-1"}, {ID: "bad"}, {ID: "good-2"},
		},
		advanced: map[machinery.FactoryID]int{},
		failOn:   "bad",
	}
	tk := engine.NewFactoryTicker(fs)
	advanced, err := tk.Tick(context.Background())
	if err == nil {
		t.Fatal("expected a joined error for the failing factory")
	}
	if !errors.Is(err, errFactoryBoom) {
		t.Fatalf("error = %v, want it to wrap errFactoryBoom", err)
	}
	if advanced != 2 {
		t.Fatalf("advanced = %d, want 2 (good-1, good-2)", advanced)
	}
	if fs.advanced["good-1"] != 1 || fs.advanced["good-2"] != 1 {
		t.Fatalf("non-failing factories not advanced: %+v", fs.advanced)
	}
}

func TestFactoryTicker_ListError(t *testing.T) {
	t.Parallel()
	fs := &fakeFactoryStore{listErr: errFactoryBoom}
	tk := engine.NewFactoryTicker(fs)
	advanced, err := tk.Tick(context.Background())
	if err == nil {
		t.Fatal("expected the ListFactories error to propagate")
	}
	if advanced != 0 {
		t.Fatalf("advanced = %d, want 0", advanced)
	}
}
