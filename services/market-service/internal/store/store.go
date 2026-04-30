// Package store defines the persistence boundary for market-service.
// The Store interface is the only seam between the HTTP/domain layer and
// the underlying database (MySQL in production, in-memory in tests).
package store

import (
	"context"
	"errors"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

// Sentinel errors that callers can branch on with errors.Is.
var (
	ErrNotFound      = errors.New("market: not found")
	ErrAlreadyExists = errors.New("market: already exists")
)

// Store is the persistence contract for market-service domain records.
type Store interface {
	// Owner operations.
	CreateOwner(ctx context.Context, o market.Owner) (market.Owner, error)
	GetOwner(ctx context.Context, id market.OwnerID) (market.Owner, error)
	ListOwners(ctx context.Context) ([]market.Owner, error)

	// Offer operations.
	CreateOffer(ctx context.Context, o market.Offer) (market.Offer, error)
	GetOffer(ctx context.Context, id market.OfferID) (market.Offer, error)
	ListOffers(ctx context.Context) ([]market.Offer, error)
	DeleteOffer(ctx context.Context, id market.OfferID) error

	// Exchange operations.
	CreateExchange(ctx context.Context, e market.Exchange) (market.Exchange, error)
	GetExchange(ctx context.Context, id market.ExchangeID) (market.Exchange, error)
	ListExchanges(ctx context.Context) ([]market.Exchange, error)

	// UniversalEquivalent is a singleton: the commodity currently elected by
	// social act to serve as universal equivalent. SetUniversalEquivalent is
	// idempotent when called with the same commodity ID.
	SetUniversalEquivalent(ctx context.Context, ue market.UniversalEquivalent) (market.UniversalEquivalent, error)
	GetUniversalEquivalent(ctx context.Context) (market.UniversalEquivalent, error)

	// MoneyCommodity is a singleton: the crystallised universal equivalent.
	SetMoneyCommodity(ctx context.Context, mc market.MoneyCommodity) (market.MoneyCommodity, error)
	GetMoneyCommodity(ctx context.Context) (market.MoneyCommodity, error)

	// Price operations. GetPrice is keyed by commodity_id.
	SetPrice(ctx context.Context, p market.Price) (market.Price, error)
	GetPrice(ctx context.Context, commodityID market.CommodityID) (market.Price, error)
	ListPrices(ctx context.Context) ([]market.Price, error)
}
