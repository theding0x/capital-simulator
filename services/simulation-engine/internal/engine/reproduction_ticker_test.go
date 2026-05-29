package engine_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

var errReproBoom = errors.New("tick boom")

type fakeReproStore struct {
	simple       []repro.SimpleReproductionScheme
	extended     []repro.ExtendedReproductionScheme
	simpleTicked map[repro.SimpleReproductionSchemeID]int
	extTicked    map[repro.ExtendedReproductionSchemeID]int
	failSimple   repro.SimpleReproductionSchemeID
	failExtended repro.ExtendedReproductionSchemeID
}

func (f *fakeReproStore) ListSimpleReproductionSchemes(context.Context, string) ([]repro.SimpleReproductionScheme, error) {
	return f.simple, nil
}

func (f *fakeReproStore) AdvanceReproductionTick(_ context.Context, id repro.SimpleReproductionSchemeID) (repro.ReproductionTick, error) {
	if id == f.failSimple {
		return repro.ReproductionTick{}, errReproBoom
	}
	if f.simpleTicked == nil {
		f.simpleTicked = map[repro.SimpleReproductionSchemeID]int{}
	}
	f.simpleTicked[id]++
	return repro.ReproductionTick{}, nil
}

func (f *fakeReproStore) ListExtendedReproductionSchemes(context.Context, string) ([]repro.ExtendedReproductionScheme, error) {
	return f.extended, nil
}

func (f *fakeReproStore) TickExtendedScheme(_ context.Context, id repro.ExtendedReproductionSchemeID) (repro.ExtendedReproductionScheme, error) {
	if id == f.failExtended {
		return repro.ExtendedReproductionScheme{}, errReproBoom
	}
	if f.extTicked == nil {
		f.extTicked = map[repro.ExtendedReproductionSchemeID]int{}
	}
	f.extTicked[id]++
	return repro.ExtendedReproductionScheme{ID: id}, nil
}

func TestReproductionTicker_AdvancesSimpleAndExtended(t *testing.T) {
	t.Parallel()
	fs := &fakeReproStore{
		simple:   []repro.SimpleReproductionScheme{{ID: "s1"}, {ID: "s2"}},
		extended: []repro.ExtendedReproductionScheme{{ID: "e1"}},
	}
	tk := engine.NewReproductionTicker(fs)
	if tk.Name() != "reproduction" {
		t.Fatalf("Name = %q, want reproduction", tk.Name())
	}
	advanced, err := tk.Tick(context.Background())
	if err != nil {
		t.Fatalf("Tick: %v", err)
	}
	if advanced != 3 {
		t.Fatalf("advanced = %d, want 3", advanced)
	}
	if fs.simpleTicked["s1"] != 1 || fs.simpleTicked["s2"] != 1 || fs.extTicked["e1"] != 1 {
		t.Fatalf("not all schemes advanced: simple=%+v extended=%+v", fs.simpleTicked, fs.extTicked)
	}
}

func TestReproductionTicker_ContinuesPastFailure(t *testing.T) {
	t.Parallel()
	fs := &fakeReproStore{
		simple:     []repro.SimpleReproductionScheme{{ID: "s1"}, {ID: "bad"}},
		extended:   []repro.ExtendedReproductionScheme{{ID: "e1"}},
		failSimple: "bad",
	}
	tk := engine.NewReproductionTicker(fs)
	advanced, err := tk.Tick(context.Background())
	if err == nil {
		t.Fatal("expected a joined error for the failing scheme")
	}
	if !errors.Is(err, errReproBoom) {
		t.Fatalf("error = %v, want it to wrap errReproBoom", err)
	}
	if advanced != 2 {
		t.Fatalf("advanced = %d, want 2 (s1, e1)", advanced)
	}
	if fs.simpleTicked["s1"] != 1 || fs.extTicked["e1"] != 1 {
		t.Fatalf("non-failing schemes not advanced: simple=%+v extended=%+v", fs.simpleTicked, fs.extTicked)
	}
}
