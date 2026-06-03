package engine

import (
	"context"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
)

// PieceWageRepricer pushes the current productivity factor to agent-service,
// which reprices every piece-wage from its baseline and appends a price point
// when the factor changes (Vol. I Ch. 21, #219). It returns the number of
// piece-wages whose price actually changed this pass. Implemented by
// *piecewage.Repricer; faked in tests.
type PieceWageRepricer interface {
	Reprice(ctx context.Context, productivityFactor float64) (advanced int, err error)
}

// ProductivitySource reports the current leading productivity factor driving
// piece-price reduction (Ch. 15 factory mechanisation → Ch. 21 piece-wages).
// Faked in tests; satisfied by *FactoryProductivitySource in production.
type ProductivitySource interface {
	LeadingProductivityFactor(ctx context.Context) (float64, error)
}

// FactoryProductivityReader lists factories so the leading productivity factor
// can be derived (Ch. 15). It is satisfied structurally by *store.Memory and
// *store.MySQL.
type FactoryProductivityReader interface {
	ListFactories(ctx context.Context) ([]machinery.Factory, error)
}

// FactoryProductivitySource derives the leading factory productivity factor —
// the most-mechanised factory sets the social pace — never below 1.0, since a
// factory at the hand-labour baseline does not push the piece price down.
type FactoryProductivitySource struct {
	factories FactoryProductivityReader
}

// NewFactoryProductivitySource constructs a FactoryProductivitySource.
func NewFactoryProductivitySource(factories FactoryProductivityReader) *FactoryProductivitySource {
	return &FactoryProductivitySource{factories: factories}
}

// LeadingProductivityFactor returns the highest Ch. 15 productivity factor over
// every persisted factory, floored at 1.0.
func (s *FactoryProductivitySource) LeadingProductivityFactor(ctx context.Context) (float64, error) {
	factories, err := s.factories.ListFactories(ctx)
	if err != nil {
		return 0, err
	}
	leading := 1.0
	for _, f := range factories {
		if pf := machinery.FactoryProductivityFactor(f); pf > leading {
			leading = pf
		}
	}
	return leading, nil
}

// PiecePriceTicker translates rising Ch. 15 factory productivity into falling
// Ch. 21 piece prices each scheduler pass: it reads the leading productivity
// factor and pushes it to agent-service, which reprices every piece-wage from
// its baseline and records a price point when the factor changes. This is the
// cross-chapter bridge issue #219 deferred to the tick scheduler. It implements
// Ticker.
type PiecePriceTicker struct {
	productivity ProductivitySource
	repricer     PieceWageRepricer
}

// NewPiecePriceTicker constructs a PiecePriceTicker over a productivity source
// and a repricer.
func NewPiecePriceTicker(productivity ProductivitySource, repricer PieceWageRepricer) *PiecePriceTicker {
	return &PiecePriceTicker{productivity: productivity, repricer: repricer}
}

// Name identifies this ticker in status and the audit log.
func (t *PiecePriceTicker) Name() string { return "piece-price" }

// Tick reads the leading productivity factor and reprices every piece-wage. It
// returns the number of piece-wages whose price actually changed this pass.
func (t *PiecePriceTicker) Tick(ctx context.Context) (int, error) {
	factor, err := t.productivity.LeadingProductivityFactor(ctx)
	if err != nil {
		return 0, err
	}
	return t.repricer.Reprice(ctx, factor)
}
