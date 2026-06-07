package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

// CapitalAccumulator is the store subset the AccumulationTicker needs. It is
// satisfied structurally by *store.Memory and *store.MySQL.
type CapitalAccumulator interface {
	FieldSnapshot(ctx context.Context) ([]circulation.FieldCapital, error)
	AccumulateCapital(ctx context.Context, id circulation.IndustrialCapitalID, deltaPence circulation.Pence) (circulation.IndustrialCapital, error)
}

// AccumulationTicker capitalises a share (alphaBP basis points) of each
// industrial capital's surplus back into it every scheduler pass — the spiral
// of accumulation (Vol. I Ch. 24). It implements Ticker.
type AccumulationTicker struct {
	store   CapitalAccumulator
	alphaBP int64
}

// NewAccumulationTicker constructs the ticker. alphaBP is the accumulation rate
// in basis points (5000 = reinvest 50% of surplus per pass); negative values
// are treated as 0 (no accumulation).
func NewAccumulationTicker(s CapitalAccumulator, alphaBP int64) *AccumulationTicker {
	if alphaBP < 0 {
		alphaBP = 0
	}
	return &AccumulationTicker{store: s, alphaBP: alphaBP}
}

// Name identifies this ticker in status and the audit log.
func (t *AccumulationTicker) Name() string { return "accumulation" }

// Tick capitalises alphaBP·surplus into every capital. It continues past a
// single capital's failure and returns the joined error.
func (t *AccumulationTicker) Tick(ctx context.Context) (int, error) {
	field, err := t.store.FieldSnapshot(ctx)
	if err != nil {
		return 0, err
	}
	var advanced int
	var errs []error
	for _, fc := range field {
		if fc.SurplusPence <= 0 || t.alphaBP <= 0 {
			continue
		}
		delta := circulation.Pence(int64(fc.SurplusPence) * t.alphaBP / 10000)
		if delta <= 0 {
			continue
		}
		if _, err := t.store.AccumulateCapital(ctx, fc.ID, delta); err != nil {
			errs = append(errs, fmt.Errorf("capital %s: %w", fc.ID, err))
			continue
		}
		advanced++
	}
	return advanced, errors.Join(errs...)
}
