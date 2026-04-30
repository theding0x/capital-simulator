// Package market models exchange and the emergence of the money-form,
// as analyzed in Karl Marx's Capital Vol. I, Ch. 2.
//
// Exchange presupposes owners: persons who stand in relation to each other as
// "guardians" of their commodities, recognising in each other the rights of
// private proprietors. When no single commodity acts as universal equivalent,
// the social act of all others sets one apart — and through repeated exchange,
// that commodity crystallises into money.
package market

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Sentinel errors.
var (
	ErrSelfExchange  = errors.New("market: giver and receiver cannot be the same owner")
	ErrOfferInvalid  = errors.New("market: offer is invalid: owner cannot seek the same commodity they offer")
	ErrValueMismatch = errors.New("market: realised values do not match")
)

// --- ID types ---------------------------------------------------------------

// OwnerID uniquely identifies an Owner.
type OwnerID string

// NewOwnerID returns a new random 96-bit hex owner ID.
func NewOwnerID() OwnerID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("market: rand: %w", err))
	}
	return OwnerID(hex.EncodeToString(b[:]))
}

func (id OwnerID) IsZero() bool { return id == "" }

// OfferID uniquely identifies an Offer.
type OfferID string

// NewOfferID returns a new random 96-bit hex offer ID.
func NewOfferID() OfferID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("market: rand: %w", err))
	}
	return OfferID(hex.EncodeToString(b[:]))
}

func (id OfferID) IsZero() bool { return id == "" }

// ExchangeID uniquely identifies a completed Exchange.
type ExchangeID string

// NewExchangeID returns a new random 96-bit hex exchange ID.
func NewExchangeID() ExchangeID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("market: rand: %w", err))
	}
	return ExchangeID(hex.EncodeToString(b[:]))
}

func (id ExchangeID) IsZero() bool { return id == "" }

// CommodityID references a commodity in commodity-service. Shares the same
// 96-bit hex format as commodity.ID to allow round-tripping across service boundaries.
type CommodityID string

func (id CommodityID) IsZero() bool { return id == "" }

// --- Value types ------------------------------------------------------------

// RealisedValue is the labour-time magnitude confirmed by the act of exchange,
// measured in minutes. It mirrors commodity.LabourMinutes but is defined here
// to avoid importing across service-internal package boundaries.
//
// "The act of exchange gives to the commodity converted into money, not its
// value, but its specific value-form." (Ch. 2)
type RealisedValue int64

// PriceAmount is the quantity of money-commodity units expressing a price.
// Integer arithmetic avoids floating-point drift when comparing proportions.
type PriceAmount int64

// LegKind distinguishes the direction of a C-M-C circuit step.
type LegKind string

const (
	KindSale     LegKind = "sale"     // C → M: selling commodity for money
	KindPurchase LegKind = "purchase" // M → C: spending money on commodity
)

// --- Domain types -----------------------------------------------------------

// Owner is a commodity-owner brought into relation with others by exchange.
// "Commodities cannot go to market and make exchanges of their own account.
// We must, therefore, have recourse to their guardians, who are also their
// owners." (Ch. 2)
type Owner struct {
	ID        OwnerID   `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Validate checks that o is well-formed.
func (o Owner) Validate() error {
	if strings.TrimSpace(o.Name) == "" {
		return errors.New("market: owner name is required")
	}
	return nil
}

// Offer is a trade intention: an owner brings a commodity to market.
// SeeksCommodityID is optional; when set, it names the specific commodity
// the owner wants in return.
type Offer struct {
	ID               OfferID     `json:"id" bson:"_id"`
	OwnerID          OwnerID     `json:"owner_id" bson:"owner_id"`
	CommodityID      CommodityID `json:"commodity_id" bson:"commodity_id"`
	Quantity         float64     `json:"quantity" bson:"quantity"`
	SeeksKind        string      `json:"seeks_kind" bson:"seeks_kind"`
	SeeksCommodityID CommodityID `json:"seeks_commodity_id,omitempty" bson:"seeks_commodity_id,omitempty"`
	CreatedAt        time.Time   `json:"created_at" bson:"created_at"`
}

// Validate checks that o is well-formed. An owner cannot seek the same
// commodity they offer — "His commodity possesses for himself no immediate
// use-value. Otherwise, he would not bring it to the market." (Ch. 2)
func (o Offer) Validate() error {
	if o.OwnerID.IsZero() {
		return errors.New("market: offer owner_id is required")
	}
	if o.CommodityID.IsZero() {
		return errors.New("market: offer commodity_id is required")
	}
	if o.Quantity <= 0 {
		return errors.New("market: offer quantity must be positive")
	}
	if !o.SeeksCommodityID.IsZero() && o.SeeksCommodityID == o.CommodityID {
		return ErrOfferInvalid
	}
	return nil
}

// BarterRatio captures the direct proportion: x use-value A = y use-value B.
// "The form of direct barter is x use-value A = y use-value B." (Ch. 2)
type BarterRatio struct {
	CommodityA CommodityID `json:"commodity_a" bson:"commodity_a"`
	QtyA       float64     `json:"qty_a" bson:"qty_a"`
	CommodityB CommodityID `json:"commodity_b" bson:"commodity_b"`
	QtyB       float64     `json:"qty_b" bson:"qty_b"`
}

// Validate checks structural integrity of the ratio.
func (r BarterRatio) Validate() error {
	if r.CommodityA.IsZero() || r.CommodityB.IsZero() {
		return errors.New("market: barter_ratio requires two commodity IDs")
	}
	if r.CommodityA == r.CommodityB {
		return errors.New("market: barter_ratio requires two different commodities")
	}
	if r.QtyA <= 0 || r.QtyB <= 0 {
		return errors.New("market: barter_ratio quantities must be positive")
	}
	return nil
}

// Exchange is a completed bilateral transfer of commodities between two owners.
// "All commodities are non-use-values for their owners, and use-values for
// their non-owners. Consequently, they must all change hands." (Ch. 2)
type Exchange struct {
	ID                  ExchangeID    `json:"id" bson:"_id"`
	GiverID             OwnerID       `json:"giver_id" bson:"giver_id"`
	ReceiverID          OwnerID       `json:"receiver_id" bson:"receiver_id"`
	GiverCommodityID    CommodityID   `json:"giver_commodity_id" bson:"giver_commodity_id"`
	GiverQty            float64       `json:"giver_qty" bson:"giver_qty"`
	ReceiverCommodityID CommodityID   `json:"receiver_commodity_id" bson:"receiver_commodity_id"`
	ReceiverQty         float64       `json:"receiver_qty" bson:"receiver_qty"`
	RealisedValue       RealisedValue `json:"realised_value" bson:"realised_value"`
	CreatedAt           time.Time     `json:"created_at" bson:"created_at"`
}

// Validate checks the exchange is structurally sound. It enforces that giver
// and receiver are distinct owners — commodity owners must "mutually recognise
// in each other the rights of private proprietors" (Ch. 2), which presupposes
// they are different persons.
func (e Exchange) Validate() error {
	if e.GiverID.IsZero() || e.ReceiverID.IsZero() {
		return errors.New("market: exchange requires giver_id and receiver_id")
	}
	if e.GiverID == e.ReceiverID {
		return ErrSelfExchange
	}
	if e.GiverCommodityID.IsZero() || e.ReceiverCommodityID.IsZero() {
		return errors.New("market: exchange requires both commodity IDs")
	}
	if e.GiverQty <= 0 || e.ReceiverQty <= 0 {
		return errors.New("market: exchange quantities must be positive")
	}
	if e.RealisedValue <= 0 {
		return errors.New("market: exchange realised_value must be positive")
	}
	return nil
}

// CircuitLeg is one step in the C-M-C circuit: either a sale (C → M) or a
// purchase (M → C).
type CircuitLeg struct {
	Kind        LegKind       `json:"kind" bson:"kind"`
	CommodityID CommodityID   `json:"commodity_id" bson:"commodity_id"`
	MoneyID     CommodityID   `json:"money_id" bson:"money_id"`
	OwnerID     OwnerID       `json:"owner_id" bson:"owner_id"`
	Value       RealisedValue `json:"value" bson:"value"`
}

// UniversalEquivalent is the commodity set apart by social act to serve as the
// form in which all others express their values.
// "The social action therefore of all other commodities, sets apart the
// particular commodity in which they all represent their values." (Ch. 2)
type UniversalEquivalent struct {
	CommodityID CommodityID `json:"commodity_id" bson:"commodity_id"`
	SetAt       time.Time   `json:"set_at" bson:"set_at"`
}

// MoneyCommodity is the universal equivalent once crystallised into the money-form.
// "Money is a crystal formed of necessity in the course of the exchanges." (Ch. 2)
type MoneyCommodity struct {
	CommodityID CommodityID `json:"commodity_id" bson:"commodity_id"`
	CreatedAt   time.Time   `json:"created_at" bson:"created_at"`
}

// Price is the value of a commodity expressed as a quantity of the money-commodity.
type Price struct {
	CommodityID      CommodityID `json:"commodity_id" bson:"commodity_id"`
	MoneyCommodityID CommodityID `json:"money_commodity_id" bson:"money_commodity_id"`
	Amount           PriceAmount `json:"amount" bson:"amount"`
	UpdatedAt        time.Time   `json:"updated_at" bson:"updated_at"`
}

// ComputePrice computes the price of one unit of a commodity expressed in
// money-commodity units:
//
//	amount = (commodity_snlt * unit_qty) / money_snlt
//
// This encodes the Petty law (Capital I, Ch. 2, fn. 12): "if a man can procure
// two ounces of silver as easily as he formerly did one, the corn will be as
// cheap at ten shillings the bushel as it was before at five shillings" —
// halving the money commodity's SNLT (doubling silver productivity) doubles
// the price of corn expressed in silver.
func ComputePrice(commoditySNLT, moneySNLT, unitQty int64) (PriceAmount, error) {
	if moneySNLT <= 0 {
		return 0, errors.New("market: money_snlt must be positive")
	}
	if commoditySNLT < 0 {
		return 0, errors.New("market: commodity_snlt must be non-negative")
	}
	if unitQty <= 0 {
		return 0, errors.New("market: unit_qty must be positive")
	}
	return PriceAmount((commoditySNLT * unitQty) / moneySNLT), nil
}
