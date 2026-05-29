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
}
