// Package store defines the persistence boundary for market-service.
// The Store interface is the only seam between the HTTP/domain layer and
// the underlying database (MySQL in production, in-memory in tests).
package store

import (
	"context"
	"errors"

	"github.com/theding0x/capital-simulator/services/market-service/internal/circulation"
	"github.com/theding0x/capital-simulator/services/market-service/internal/costs"
	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

// Sentinel errors that callers can branch on with errors.Is.
var (
	ErrNotFound      = errors.New("market: not found")
	ErrAlreadyExists = errors.New("market: already exists")
)

// TurnoverTimeStore is the persistence contract for Vol. II Ch. 5 circulation
// time records. It is embedded in Store.
type TurnoverTimeStore interface {
	// TurnoverTime operations.
	CreateTurnoverTime(ctx context.Context, t circulation.TurnoverTime) (circulation.TurnoverTime, error)
	GetTurnoverTime(ctx context.Context, id circulation.TurnoverTimeID) (circulation.TurnoverTime, error)
	ListTurnoverTimes(ctx context.Context) ([]circulation.TurnoverTime, error)
	AddLabourTime(ctx context.Context, id circulation.TurnoverTimeID, nanos int64) (circulation.TurnoverTime, error)
	AddLabourInterruption(ctx context.Context, id circulation.TurnoverTimeID, nanos int64) (circulation.TurnoverTime, error)

	// Sub-span records.
	RecordLatentMP(ctx context.Context, lpc circulation.LatentProductiveCapital) (circulation.LatentProductiveCapital, error)
	RecordNaturalProcess(ctx context.Context, nps circulation.NaturalProcessSpan) (circulation.NaturalProcessSpan, error)

	// Selling and buying phases. Open = ClosedAt nil; Close = ClosedAt set.
	OpenSellingPhase(ctx context.Context, sp circulation.SellingPhase) (circulation.SellingPhase, error)
	CloseSellingPhase(ctx context.Context, id circulation.TurnoverTimeID, spID circulation.SellingPhaseID, outcome circulation.SellingOutcome) (circulation.SellingPhase, error)
	OpenBuyingPhase(ctx context.Context, bp circulation.BuyingPhase) (circulation.BuyingPhase, error)
	CloseBuyingPhase(ctx context.Context, id circulation.TurnoverTimeID, bpID circulation.BuyingPhaseID) (circulation.BuyingPhase, error)

	// Perishability and market separation.
	SetPerishability(ctx context.Context, p circulation.Perishability) (circulation.Perishability, error)
	GetPerishability(ctx context.Context, commodityID circulation.CommodityID) (circulation.Perishability, error)
	SetMarketSeparation(ctx context.Context, ms circulation.MarketSeparation) (circulation.MarketSeparation, error)
}

// CirculationTimeV2Store is the persistence contract for Vol. II Ch. 14
// circulation-time refinements: phase sub-spans, distance-lag relations,
// speed-improvement events, and the annual-surplus-penalty stub.
type CirculationTimeV2Store interface {
	// Selling phase details — decomposes a SellingPhase into sub-spans.
	CreateSellingPhaseDetailed(ctx context.Context, d circulation.SellingPhaseDetailed) (circulation.SellingPhaseDetailed, error)
	GetSellingPhaseDetailed(ctx context.Context, id circulation.SellingPhaseDetailedID) (circulation.SellingPhaseDetailed, error)

	// Buying phase details — decomposes a BuyingPhase into sub-spans.
	CreateBuyingPhaseDetailed(ctx context.Context, d circulation.BuyingPhaseDetailed) (circulation.BuyingPhaseDetailed, error)
	GetBuyingPhaseDetailed(ctx context.Context, id circulation.BuyingPhaseDetailedID) (circulation.BuyingPhaseDetailed, error)

	// Distance-lag relations — median lag per km for a market route.
	CreateDistanceLagRelation(ctx context.Context, r circulation.DistanceLagRelation) (circulation.DistanceLagRelation, error)
	GetDistanceLagRelation(ctx context.Context, id circulation.DistanceLagRelationID) (circulation.DistanceLagRelation, error)
	ListDistanceLagRelations(ctx context.Context) ([]circulation.DistanceLagRelation, error)

	// Circulation speed improvements — events that changed lag per km.
	CreateCirculationSpeedImprovement(ctx context.Context, csi circulation.CirculationSpeedImprovement) (circulation.CirculationSpeedImprovement, error)
	ListCirculationSpeedImprovements(ctx context.Context, marketID circulation.MarketID) ([]circulation.CirculationSpeedImprovement, error)

	// Annual surplus penalty — stub, full calculation deferred to Ch. 16.
	CreateAnnualSurplusPenalty(ctx context.Context, asp circulation.AnnualSurplusPenalty) (circulation.AnnualSurplusPenalty, error)
	GetAnnualSurplusPenalty(ctx context.Context, icID circulation.IndustrialCapitalID, period string) (circulation.AnnualSurplusPenalty, error)
}

// CirculationCostStore is the persistence contract for Vol. II Ch. 6
// costs-of-circulation records. It is embedded in Store.
type CirculationCostStore interface {
	CreateCirculationCost(ctx context.Context, c costs.CirculationCost) (costs.CirculationCost, error)
	GetCirculationCost(ctx context.Context, id costs.CirculationCostID) (costs.CirculationCost, error)
	ListCirculationCosts(ctx context.Context, icID costs.IndustrialCapitalID, kind costs.CirculationCostKind, nature costs.CostNature) ([]costs.CirculationCost, error)

	CreateCirculationAgent(ctx context.Context, a costs.CirculationAgent) (costs.CirculationAgent, error)
	ListCirculationAgents(ctx context.Context, icID costs.IndustrialCapitalID) ([]costs.CirculationAgent, error)

	CreateMoneyAsFauxFrais(ctx context.Context, m costs.MoneyAsFauxFrais) (costs.MoneyAsFauxFrais, error)
	ListMoneyAsFauxFrais(ctx context.Context) ([]costs.MoneyAsFauxFrais, error)

	CreateCommoditySupply(ctx context.Context, s costs.CommoditySupply) (costs.CommoditySupply, error)
	GetCommoditySupply(ctx context.Context, id costs.CommoditySupplyID) (costs.CommoditySupply, error)
	ListCommoditySupplies(ctx context.Context, icID costs.IndustrialCapitalID) ([]costs.CommoditySupply, error)

	CreateStorageCost(ctx context.Context, sc costs.StorageCost) (costs.StorageCost, error)
	ListStorageCosts(ctx context.Context, supplyID costs.CommoditySupplyID) ([]costs.StorageCost, error)

	SetTransportTariff(ctx context.Context, t costs.TransportTariff) (costs.TransportTariff, error)
	GetTransportTariff(ctx context.Context, commodityID costs.CommodityID) (costs.TransportTariff, error)

	CreateTransportLeg(ctx context.Context, leg costs.TransportLeg) (costs.TransportLeg, error)
	ListTransportLegs(ctx context.Context, commodityID costs.CommodityID) ([]costs.TransportLeg, error)

	AggregateForCapital(ctx context.Context, icID costs.IndustrialCapitalID, period string) (costs.AggregateResult, error)
	SystemFauxFrais(ctx context.Context, period string) (costs.SystemFauxFraisResult, error)
}

// MoneyCh3Store is the persistence contract for Vol. I Ch. 3 money-form
// records: hoards (money withdrawn from circulation), payment obligations
// (means of payment / credit), and world-money transfers (cross-border
// settlement). It is embedded in Store.
type MoneyCh3Store interface {
	// Hoard operations.
	CreateHoard(ctx context.Context, h market.Hoard) (market.Hoard, error)
	GetHoard(ctx context.Context, id market.HoardID) (market.Hoard, error)
	ListHoards(ctx context.Context) ([]market.Hoard, error)

	// PaymentObligation operations. SettlePaymentObligation marks an
	// obligation paid and returns ErrNotFound when the id is unknown.
	CreatePaymentObligation(ctx context.Context, p market.PaymentObligation) (market.PaymentObligation, error)
	GetPaymentObligation(ctx context.Context, id market.PaymentObligationID) (market.PaymentObligation, error)
	ListPaymentObligations(ctx context.Context) ([]market.PaymentObligation, error)
	SettlePaymentObligation(ctx context.Context, id market.PaymentObligationID) (market.PaymentObligation, error)

	// WorldMoneyTransfer operations.
	CreateWorldMoneyTransfer(ctx context.Context, t market.WorldMoneyTransfer) (market.WorldMoneyTransfer, error)
	ListWorldMoneyTransfers(ctx context.Context) ([]market.WorldMoneyTransfer, error)
}

// Store is the persistence contract for market-service domain records.
type Store interface {
	TurnoverTimeStore
	CirculationTimeV2Store
	CirculationCostStore
	MoneyCh3Store
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
