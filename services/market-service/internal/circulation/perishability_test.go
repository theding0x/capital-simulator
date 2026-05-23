package circulation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/market-service/internal/circulation"
)

// Perishability window constants.
// §: "The more perishable a commodity is … the less is it suited to be an
// object of capitalist production."
const (
	twoDayWindowNanos   = int64(2 * 24 * 60 * 60 * 1_000_000_000)
	sevenDayWindowNanos = int64(7 * 24 * 60 * 60 * 1_000_000_000)
	threeDayNanos       = int64(3 * 24 * 60 * 60 * 1_000_000_000)
	oneDayNanos         = int64(1 * 24 * 60 * 60 * 1_000_000_000)
)

func TestPerishability_Spoiled_MilkExpired(t *testing.T) {
	t.Parallel()
	// Milk: 2-day window. 3 days elapsed → spoiled.
	p := circulation.Perishability{
		ID:          circulation.NewPerishabilityID(),
		CommodityID: "milk-commodity-id",
		WindowNanos: twoDayWindowNanos,
	}
	if !p.Spoiled(threeDayNanos) {
		t.Errorf("milk elapsed %d > window %d must be spoiled", threeDayNanos, twoDayWindowNanos)
	}
}

func TestPerishability_Spoiled_MilkFresh(t *testing.T) {
	t.Parallel()
	// Milk: 2-day window. 1 day elapsed → not spoiled.
	p := circulation.Perishability{
		ID:          circulation.NewPerishabilityID(),
		CommodityID: "milk-commodity-id",
		WindowNanos: twoDayWindowNanos,
	}
	if p.Spoiled(oneDayNanos) {
		t.Errorf("milk elapsed %d < window %d must not be spoiled", oneDayNanos, twoDayWindowNanos)
	}
}

func TestPerishability_Spoiled_NoConstraint(t *testing.T) {
	t.Parallel()
	// WindowNanos == 0 means no perishability constraint.
	p := circulation.Perishability{
		ID:          circulation.NewPerishabilityID(),
		CommodityID: "iron-commodity-id",
		WindowNanos: 0,
	}
	if p.Spoiled(int64(1e18)) {
		t.Error("unconstrained commodity (WindowNanos=0) must never spoil")
	}
}

func TestPerishability_Spoiled_AtExactBoundary(t *testing.T) {
	t.Parallel()
	// Elapsed exactly equals WindowNanos — must NOT be spoiled (strictly >).
	p := circulation.Perishability{
		WindowNanos: twoDayWindowNanos,
	}
	if p.Spoiled(twoDayWindowNanos) {
		t.Error("elapsed == window must not spoil (spoilage requires strictly >)")
	}
}

func TestPerishability_IsConstrained_WithWindow(t *testing.T) {
	t.Parallel()
	// Beer: 7-day window.
	p := circulation.Perishability{WindowNanos: sevenDayWindowNanos}
	if !p.IsConstrained() {
		t.Error("WindowNanos > 0 must be constrained")
	}
}

func TestPerishability_IsConstrained_NoWindow(t *testing.T) {
	t.Parallel()
	p := circulation.Perishability{WindowNanos: 0}
	if p.IsConstrained() {
		t.Error("WindowNanos == 0 must not be constrained")
	}
}

func TestMarketSeparation_IsSeparated_DifferentMarkets(t *testing.T) {
	t.Parallel()
	// §: "C—M and M—C may be separate not only in time but also in space."
	ms := circulation.MarketSeparation{
		ID:                  circulation.NewMarketSeparationID(),
		IndustrialCapitalID: "5eed000000000000000502",
		SellingMarketID:     "manchester",
		BuyingMarketID:      "liverpool",
	}
	if !ms.IsSeparated() {
		t.Error("different selling and buying markets must be separated")
	}
}

func TestMarketSeparation_IsSeparated_SameMarket(t *testing.T) {
	t.Parallel()
	ms := circulation.MarketSeparation{
		SellingMarketID: "manchester",
		BuyingMarketID:  "manchester",
	}
	if ms.IsSeparated() {
		t.Error("same selling and buying market must not be separated")
	}
}

func TestCirculationCompressionTarget_IsPaymentOnDelivery_Zero(t *testing.T) {
	t.Parallel()
	// §: "if payment is made in his own means of production,
	// the time of circulation approaches zero."
	cct := circulation.CirculationCompressionTarget{
		ID:                     circulation.NewCirculationCompressionTargetID(),
		IndustrialCapitalID:    "5eed000000000000000502",
		TargetCirculationNanos: 0,
	}
	if !cct.IsPaymentOnDelivery() {
		t.Error("TargetCirculationNanos == 0 must be payment-on-delivery")
	}
}

func TestCirculationCompressionTarget_IsPaymentOnDelivery_Nonzero(t *testing.T) {
	t.Parallel()
	cct := circulation.CirculationCompressionTarget{
		TargetCirculationNanos: oneDayNanos,
	}
	if cct.IsPaymentOnDelivery() {
		t.Error("TargetCirculationNanos > 0 must not be payment-on-delivery")
	}
}

func TestNewPerishabilityID_Unique(t *testing.T) {
	t.Parallel()
	a := circulation.NewPerishabilityID()
	b := circulation.NewPerishabilityID()
	if a == b {
		t.Errorf("duplicate PerishabilityID: %q", a)
	}
}

func TestNewMarketSeparationID_Unique(t *testing.T) {
	t.Parallel()
	a := circulation.NewMarketSeparationID()
	b := circulation.NewMarketSeparationID()
	if a == b {
		t.Errorf("duplicate MarketSeparationID: %q", a)
	}
}
