package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
)

// FactoryAdvancer is the subset of the store the FactoryTicker needs: list every
// persisted Factory and advance one period of wear / value-transfer for one.
// It is satisfied structurally by *store.Memory and *store.MySQL.
type FactoryAdvancer interface {
	ListFactories(ctx context.Context) ([]machinery.Factory, error)
	AdvanceTick(ctx context.Context, id machinery.FactoryID) (machinery.Factory, Tick, error)
}

// FactoryTicker advances every persisted Factory (Ch. 15) by one period per
// scheduler pass — the machinery / factory loop the spec deferred to the tick
// scheduler. It implements Ticker.
type FactoryTicker struct {
	store FactoryAdvancer
}

// NewFactoryTicker constructs a FactoryTicker over the given store.
func NewFactoryTicker(s FactoryAdvancer) *FactoryTicker {
	return &FactoryTicker{store: s}
}

// Name identifies this ticker in status and the audit log.
func (t *FactoryTicker) Name() string { return "factory" }

// Tick advances every factory one period. It continues past a single factory's
// failure and returns the joined error, so one bad factory does not starve the
// rest of the pass.
func (t *FactoryTicker) Tick(ctx context.Context) (int, error) {
	factories, err := t.store.ListFactories(ctx)
	if err != nil {
		return 0, err
	}
	var advanced int
	var errs []error
	for _, f := range factories {
		if _, _, err := t.store.AdvanceTick(ctx, f.ID); err != nil {
			errs = append(errs, fmt.Errorf("factory %s: %w", f.ID, err))
			continue
		}
		advanced++
	}
	return advanced, errors.Join(errs...)
}
