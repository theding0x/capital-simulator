package engine

import (
	"context"
	"errors"
	"fmt"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// ReproductionAdvancer is the subset of the store the ReproductionTicker needs:
// list and advance both the simple (Ch. 20) and extended (Ch. 21) two-department
// reproduction schemes by one cycle. It is satisfied structurally by
// *store.Memory and *store.MySQL.
type ReproductionAdvancer interface {
	ListSimpleReproductionSchemes(ctx context.Context, period string) ([]repro.SimpleReproductionScheme, error)
	AdvanceReproductionTick(ctx context.Context, id repro.SimpleReproductionSchemeID, workerPoolSize, subsistenceBasketPence int64) (repro.ReproductionTick, error)
	ListExtendedReproductionSchemes(ctx context.Context, period string) ([]repro.ExtendedReproductionScheme, error)
	TickExtendedScheme(ctx context.Context, id repro.ExtendedReproductionSchemeID) (repro.ExtendedReproductionScheme, error)
}

// ReproductionTicker advances every persisted reproduction scheme — simple
// (Ch. 20) and extended (Ch. 21) — by one cycle per scheduler pass. This is the
// multi-period advancement the Vol. II Ch. 21 spec deferred to the tick loop.
// It implements Ticker.
type ReproductionTicker struct {
	store ReproductionAdvancer
}

// NewReproductionTicker constructs a ReproductionTicker over the given store.
func NewReproductionTicker(s ReproductionAdvancer) *ReproductionTicker {
	return &ReproductionTicker{store: s}
}

// Name identifies this ticker in status and the audit log.
func (t *ReproductionTicker) Name() string { return "reproduction" }

// Tick advances every simple and extended scheme one cycle. It continues past a
// single scheme's failure (an unbalanced scheme cannot be ticked) and returns
// the joined error so one bad scheme does not starve the rest of the pass.
func (t *ReproductionTicker) Tick(ctx context.Context) (int, error) {
	var advanced int
	var errs []error

	simple, err := t.store.ListSimpleReproductionSchemes(ctx, "")
	if err != nil {
		errs = append(errs, fmt.Errorf("list simple schemes: %w", err))
	}
	for _, s := range simple {
		// Automated ticks do not carry labour-force inputs (issue #220).
		if _, err := t.store.AdvanceReproductionTick(ctx, s.ID, 0, 0); err != nil {
			errs = append(errs, fmt.Errorf("simple scheme %s: %w", s.ID, err))
			continue
		}
		advanced++
	}

	extended, err := t.store.ListExtendedReproductionSchemes(ctx, "")
	if err != nil {
		errs = append(errs, fmt.Errorf("list extended schemes: %w", err))
	}
	for _, s := range extended {
		if _, err := t.store.TickExtendedScheme(ctx, s.ID); err != nil {
			errs = append(errs, fmt.Errorf("extended scheme %s: %w", s.ID, err))
			continue
		}
		advanced++
	}

	return advanced, errors.Join(errs...)
}
