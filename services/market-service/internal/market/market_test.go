package market_test

import (
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

func TestBarterRatio_Validate_LinenCoat(t *testing.T) {
	t.Parallel()
	// "x Commodity A = y Commodity B": 20 yards linen = 1 coat
	// Both sides must have different commodities and positive quantities.
	linenID := market.CommodityID("aaaaaaaaaaaaaaaaaaaaaaaaaa")
	coatID := market.CommodityID("bbbbbbbbbbbbbbbbbbbbbbbbbb")
	r := market.BarterRatio{CommodityA: linenID, QtyA: 20, CommodityB: coatID, QtyB: 1}
	if err := r.Validate(); err != nil {
		t.Fatalf("valid barter ratio rejected: %v", err)
	}
}

func TestBarterRatio_Validate_SameCommodity(t *testing.T) {
	t.Parallel()
	id := market.CommodityID("aaaaaaaaaaaaaaaaaaaaaaaaaa")
	r := market.BarterRatio{CommodityA: id, QtyA: 20, CommodityB: id, QtyB: 20}
	if err := r.Validate(); err == nil {
		t.Fatal("expected error for barter ratio with identical commodities")
	}
}

func TestExchange_Validate_SelfExchange(t *testing.T) {
	t.Parallel()
	// "commodities must be realised as values before they can be realised as
	// use-values" — but the giver and receiver must be different persons.
	ownerID := market.NewOwnerID()
	e := market.Exchange{
		GiverID:             ownerID,
		ReceiverID:          ownerID,
		GiverCommodityID:    market.CommodityID("aaaaaaaaaaaaaaaaaaaaaaaaaa"),
		GiverQty:            20,
		ReceiverCommodityID: market.CommodityID("bbbbbbbbbbbbbbbbbbbbbbbbbb"),
		ReceiverQty:         1,
		RealisedValue:       600,
	}
	if err := e.Validate(); !errors.Is(err, market.ErrSelfExchange) {
		t.Fatalf("expected ErrSelfExchange, got %v", err)
	}
}

func TestExchange_Validate_Valid(t *testing.T) {
	t.Parallel()
	giverID := market.NewOwnerID()
	receiverID := market.NewOwnerID()
	e := market.Exchange{
		GiverID:             giverID,
		ReceiverID:          receiverID,
		GiverCommodityID:    market.CommodityID("aaaaaaaaaaaaaaaaaaaaaaaaaa"),
		GiverQty:            20,
		ReceiverCommodityID: market.CommodityID("bbbbbbbbbbbbbbbbbbbbbbbbbb"),
		ReceiverQty:         1,
		RealisedValue:       600,
	}
	if err := e.Validate(); err != nil {
		t.Fatalf("valid exchange rejected: %v", err)
	}
}

func TestOffer_Validate_SeeksOwnCommodity(t *testing.T) {
	t.Parallel()
	// "All commodities are non-use-values for their owners" — an owner cannot
	// seek the same commodity they bring to market.
	ownerID := market.NewOwnerID()
	commodityID := market.CommodityID("aaaaaaaaaaaaaaaaaaaaaaaaaa")
	o := market.Offer{
		OwnerID:          ownerID,
		CommodityID:      commodityID,
		Quantity:         20,
		SeeksKind:        "linen",
		SeeksCommodityID: commodityID,
	}
	if err := o.Validate(); !errors.Is(err, market.ErrOfferInvalid) {
		t.Fatalf("expected ErrOfferInvalid, got %v", err)
	}
}

func TestOffer_Validate_Valid(t *testing.T) {
	t.Parallel()
	ownerID := market.NewOwnerID()
	o := market.Offer{
		OwnerID:     ownerID,
		CommodityID: market.CommodityID("aaaaaaaaaaaaaaaaaaaaaaaaaa"),
		Quantity:    20,
		SeeksKind:   "coat",
	}
	if err := o.Validate(); err != nil {
		t.Fatalf("valid offer rejected: %v", err)
	}
}

func TestComputePrice_PettyLaw(t *testing.T) {
	t.Parallel()
	// Petty (Capital Vol. I, Ch. 2, fn. 12):
	// "if a man can procure two ounces of silver as easily as he formerly did one,
	// the corn will be as cheap at ten shillings the bushel as it was before at
	// five shillings, caeteris paribus."
	// Halving silver SNLT (doubling productivity) must double corn's price in silver.
	cornSNLT := int64(300)
	silverSNLT := int64(300)
	p1, err := market.ComputePrice(cornSNLT, silverSNLT, 1)
	if err != nil {
		t.Fatalf("compute price: %v", err)
	}
	p2, err := market.ComputePrice(cornSNLT, silverSNLT/2, 1)
	if err != nil {
		t.Fatalf("compute price with halved silver SNLT: %v", err)
	}
	if p2 != p1*2 {
		t.Fatalf("Petty law violated: expected doubled price %d, got %d", p1*2, p2)
	}
}

func TestComputePrice_ZeroMoneySNLT(t *testing.T) {
	t.Parallel()
	_, err := market.ComputePrice(300, 0, 1)
	if err == nil {
		t.Fatal("expected error for zero money_snlt")
	}
}

func TestNewOwnerID_Unique(t *testing.T) {
	t.Parallel()
	a := market.NewOwnerID()
	b := market.NewOwnerID()
	if a == b {
		t.Fatal("NewOwnerID returned duplicate IDs")
	}
	if a.IsZero() || b.IsZero() {
		t.Fatal("NewOwnerID returned zero ID")
	}
}
