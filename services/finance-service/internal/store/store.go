// Package store defines the persistence boundary for finance-service.
// The Store interface is the only seam between the HTTP/domain layer and
// the underlying database (MySQL in production, in-memory in tests).
//
// finance-service is the home for Capital Vol. III — profit, the average
// rate of profit, commercial capital, interest-bearing capital, credit,
// rent, and the revenue distribution. Each Vol. III chapter PR adds the
// methods it needs to this interface and provides the matching Memory +
// MySQL implementations.
package store

import (
	"context"
	"errors"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// Sentinel errors that callers can branch on with errors.Is.
var (
	ErrNotFound      = errors.New("finance: not found")
	ErrAlreadyExists = errors.New("finance: already exists")
)

// Store is the persistence contract for finance-service.
//
// Vol. III Ch. 1 ("Cost-Price and Profit") opens the interface with cost-price
// persistence. Later chapters append their own methods here and the matching
// implementations in memory.go + mysql.go, following the per-chapter pattern
// documented in CLAUDE.md.
type Store interface {
	// CreateCostPrice persists a computed cost-price, assigning an ID and
	// created-at timestamp when absent. It returns ErrAlreadyExists if the ID
	// collides.
	CreateCostPrice(ctx context.Context, cp profit.CostPrice) (profit.CostPrice, error)
	// GetCostPrice returns the cost-price with the given ID, or ErrNotFound.
	GetCostPrice(ctx context.Context, id profit.CostPriceID) (profit.CostPrice, error)
	// ListCostPrices returns all stored cost-prices, newest first.
	ListCostPrices(ctx context.Context) ([]profit.CostPrice, error)

	// CreateProfitRate persists a rate-of-profit analysis (Vol. III Ch. 2),
	// assigning an ID and created-at timestamp when absent. It returns
	// ErrAlreadyExists if the ID collides.
	CreateProfitRate(ctx context.Context, a profit.ProfitRateAnalysis) (profit.ProfitRateAnalysis, error)
	// GetProfitRate returns the analysis with the given ID, or ErrNotFound.
	GetProfitRate(ctx context.Context, id profit.ProfitRateID) (profit.ProfitRateAnalysis, error)
	// ListProfitRates returns all stored analyses, newest first.
	ListProfitRates(ctx context.Context) ([]profit.ProfitRateAnalysis, error)

	// CreateVariation persists a variation analysis (Vol. III Ch. 3) — the
	// movement of the rate of profit between two decompositions — assigning an
	// ID and created-at timestamp when absent. It returns ErrAlreadyExists if
	// the ID collides.
	CreateVariation(ctx context.Context, a profit.VariationAnalysis) (profit.VariationAnalysis, error)
	// GetVariation returns the variation analysis with the given ID, or ErrNotFound.
	GetVariation(ctx context.Context, id profit.VariationAnalysisID) (profit.VariationAnalysis, error)
	// ListVariations returns all stored variation analyses, newest first.
	ListVariations(ctx context.Context) ([]profit.VariationAnalysis, error)

	// CreateTurnoverAnalysis persists a turnover analysis (Vol. III Ch. 4) — the
	// annual rate of profit p' = s'·n·(v/C) corrected for the number of turnovers
	// — assigning an ID and created-at timestamp when absent. It returns
	// ErrAlreadyExists if the ID collides.
	CreateTurnoverAnalysis(ctx context.Context, a profit.TurnoverAnalysis) (profit.TurnoverAnalysis, error)
	// GetTurnoverAnalysis returns the turnover analysis with the given ID, or ErrNotFound.
	GetTurnoverAnalysis(ctx context.Context, id profit.TurnoverAnalysisID) (profit.TurnoverAnalysis, error)
	// ListTurnoverAnalyses returns all stored turnover analyses, newest first.
	ListTurnoverAnalyses(ctx context.Context) ([]profit.TurnoverAnalysis, error)

	// CreateEconomyAnalysis persists an economy analysis (Vol. III Ch. 5) — a
	// saving in constant capital and the rise it produces in the rate of profit
	// s/(c+v) — assigning an ID and created-at timestamp when absent. It returns
	// ErrAlreadyExists if the ID collides.
	CreateEconomyAnalysis(ctx context.Context, a profit.EconomyAnalysis) (profit.EconomyAnalysis, error)
	// GetEconomyAnalysis returns the economy analysis with the given ID, or ErrNotFound.
	GetEconomyAnalysis(ctx context.Context, id profit.EconomyAnalysisID) (profit.EconomyAnalysis, error)
	// ListEconomyAnalyses returns all stored economy analyses, newest first.
	ListEconomyAnalyses(ctx context.Context) ([]profit.EconomyAnalysis, error)

	// CreatePriceFluctuationAnalysis persists a price-fluctuation analysis (Vol.
	// III Ch. 6) — the movement in the rate of profit produced by a re-pricing of
	// the raw-material element of constant capital — assigning an ID and
	// created-at timestamp when absent. It returns ErrAlreadyExists if the ID
	// collides.
	CreatePriceFluctuationAnalysis(ctx context.Context, a profit.PriceFluctuationAnalysis) (profit.PriceFluctuationAnalysis, error)
	// GetPriceFluctuationAnalysis returns the analysis with the given ID, or ErrNotFound.
	GetPriceFluctuationAnalysis(ctx context.Context, id profit.PriceFluctuationAnalysisID) (profit.PriceFluctuationAnalysis, error)
	// ListPriceFluctuationAnalyses returns all stored analyses, newest first.
	ListPriceFluctuationAnalyses(ctx context.Context) ([]profit.PriceFluctuationAnalysis, error)

	// CreateCompositionEffect persists an organic-composition comparison (Vol. III
	// Ch. 7) — two capitals of equal s and v but different c, and the gap in their
	// rates of profit — assigning an ID and created-at timestamp when absent. It
	// returns ErrAlreadyExists if the ID collides.
	CreateCompositionEffect(ctx context.Context, a profit.CompositionEffectAnalysis) (profit.CompositionEffectAnalysis, error)
	// GetCompositionEffect returns the comparison with the given ID, or ErrNotFound.
	GetCompositionEffect(ctx context.Context, id profit.CompositionEffectID) (profit.CompositionEffectAnalysis, error)
	// ListCompositionEffects returns all stored comparisons, newest first.
	ListCompositionEffects(ctx context.Context) ([]profit.CompositionEffectAnalysis, error)

	// CreateMagnitudeChange persists a capital-magnitude change (Vol. III Ch. 7) —
	// a change in the size of capital and what becomes of its rate of profit —
	// assigning an ID and created-at timestamp when absent. It returns
	// ErrAlreadyExists if the ID collides.
	CreateMagnitudeChange(ctx context.Context, a profit.MagnitudeChangeAnalysis) (profit.MagnitudeChangeAnalysis, error)
	// GetMagnitudeChange returns the change with the given ID, or ErrNotFound.
	GetMagnitudeChange(ctx context.Context, id profit.MagnitudeChangeID) (profit.MagnitudeChangeAnalysis, error)
	// ListMagnitudeChanges returns all stored changes, newest first.
	ListMagnitudeChanges(ctx context.Context) ([]profit.MagnitudeChangeAnalysis, error)

	// CreateProductionSphere persists a production-sphere analysis (Vol. III Ch. 8)
	// — a branch of production with its organic composition and individual rate of
	// profit — assigning an ID and created-at timestamp when absent. It returns
	// ErrAlreadyExists if the ID collides.
	CreateProductionSphere(ctx context.Context, s avgprofit.ProductionSphere) (avgprofit.ProductionSphere, error)
	// GetProductionSphere returns the sphere with the given ID, or ErrNotFound.
	GetProductionSphere(ctx context.Context, id avgprofit.ProductionSphereID) (avgprofit.ProductionSphere, error)
	// ListProductionSpheres returns all stored spheres, newest first.
	ListProductionSpheres(ctx context.Context) ([]avgprofit.ProductionSphere, error)
}
