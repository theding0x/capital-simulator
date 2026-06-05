package engine

import (
	"context"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// AbodeAdvancer is the store subset the GeneralLawTicker needs. It is satisfied
// structurally by *store.Memory and *store.MySQL.
type AbodeAdvancer interface {
	GetAbodeState(ctx context.Context) (simulation.AbodeState, error)
	AdvanceAbode(ctx context.Context, next simulation.AbodeState, period simulation.GeneralLawPeriod) error
}

// GeneralLawTicker advances the hidden abode one period of the General Law of
// Capitalist Accumulation (Vol. I Ch. 25) every scheduler pass: rising
// composition and machinery repel labour into the reserve army, which presses
// the wage down and raises s/v, whose surplus re-accumulates. It implements
// Ticker. Registered after the AccumulationTicker so the surface (the field of
// orbits) and the abode (the class relation) advance in the same pass.
type GeneralLawTicker struct {
	store AbodeAdvancer
}

// NewGeneralLawTicker constructs the ticker.
func NewGeneralLawTicker(s AbodeAdvancer) *GeneralLawTicker {
	return &GeneralLawTicker{store: s}
}

// Name identifies this ticker in status and the audit log.
func (t *GeneralLawTicker) Name() string { return "general-law" }

// Tick advances the abode by one period and records it.
func (t *GeneralLawTicker) Tick(ctx context.Context) (int, error) {
	state, err := t.store.GetAbodeState(ctx)
	if err != nil {
		return 0, err
	}
	next, period := simulation.AdvanceGeneralLaw(state)
	if err := t.store.AdvanceAbode(ctx, next, period); err != nil {
		return 0, err
	}
	return 1, nil
}
