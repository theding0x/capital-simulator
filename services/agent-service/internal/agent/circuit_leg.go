package agent

import (
	"errors"
	"time"
)

// CircuitLegKind distinguishes the two complementary halves of a coordinated
// multi-agent exchange.
type CircuitLegKind string

const (
	CircuitLegSale     CircuitLegKind = "sale"     // C → M: parted with a commodity, realised its value
	CircuitLegPurchase CircuitLegKind = "purchase" // M → C: spent value, obtained a commodity
)

// Valid reports whether k is one of the two recognised leg kinds.
func (k CircuitLegKind) Valid() bool {
	return k == CircuitLegSale || k == CircuitLegPurchase
}

// CircuitLeg is one closed half of a coordinated multi-agent exchange. Ch. 4
// recorded the M—C—M′ / C—M—C circuit on a single agent in isolation; a
// market-service Exchange between two owners now lands a matching leg on each
// owner's circuit through the cross-service settlement webhook (#216). The
// seller closes a sale leg (C → M); the buyer closes a purchase leg (M → C).
//
// The leg is a loosely-coupled record keyed by the owner/agent id — like the
// market service's reference to a commodity id, it does not require the
// counterparty to exist as an agent-service Agent, so it does not mutate any
// balance. Value is carried in LabourMinutes, the canonical value-magnitude
// unit, matching the market Exchange's realised value.
type CircuitLeg struct {
	ID             ID             `json:"id"`
	AgentID        ID             `json:"agent_id"`
	Kind           CircuitLegKind `json:"kind"`
	CounterpartyID ID             `json:"counterparty_id"`
	CommodityID    string         `json:"commodity_id"`
	ValueMinutes   LabourMinutes  `json:"value_minutes"`
	ExchangeID     string         `json:"exchange_id"`
	CreatedAt      time.Time      `json:"created_at"`
}

// ErrCircuitLegSelfExchange is returned when a leg names the same owner as both
// party and counterparty — a single owner cannot exchange with themselves.
var ErrCircuitLegSelfExchange = errors.New("circuit_leg: agent and counterparty cannot be the same owner")

// Validate checks the leg is structurally sound.
func (l CircuitLeg) Validate() error {
	if l.AgentID.IsZero() {
		return errors.New("circuit_leg: agent_id is required")
	}
	if l.CounterpartyID.IsZero() {
		return errors.New("circuit_leg: counterparty_id is required")
	}
	if l.AgentID == l.CounterpartyID {
		return ErrCircuitLegSelfExchange
	}
	if !l.Kind.Valid() {
		return errors.New("circuit_leg: kind must be 'sale' or 'purchase'")
	}
	if l.CommodityID == "" {
		return errors.New("circuit_leg: commodity_id is required")
	}
	if l.ValueMinutes <= 0 {
		return errors.New("circuit_leg: value_minutes must be positive")
	}
	if l.ExchangeID == "" {
		return errors.New("circuit_leg: exchange_id is required")
	}
	return nil
}

// ExchangeSettlement is the cross-service event a completed market Exchange
// emits to agent-service: it names both owners, the commodity that changed
// hands, and the realised value. agent-service turns it into the two
// complementary circuit legs.
type ExchangeSettlement struct {
	ExchangeID          string        `json:"exchange_id"`
	GiverID             ID            `json:"giver_id"`
	ReceiverID          ID            `json:"receiver_id"`
	GiverCommodityID    string        `json:"giver_commodity_id"`
	ReceiverCommodityID string        `json:"receiver_commodity_id"`
	RealisedValue       LabourMinutes `json:"realised_value"`
}

// Legs derives the two complementary circuit legs from a settlement. The giver
// (seller) closes a sale of the commodity they parted with; the receiver
// (buyer) closes a purchase of that same commodity. Both legs carry the
// realised value and the originating exchange id.
func (s ExchangeSettlement) Legs() (giver, receiver CircuitLeg) {
	giver = CircuitLeg{
		AgentID:        s.GiverID,
		Kind:           CircuitLegSale,
		CounterpartyID: s.ReceiverID,
		CommodityID:    s.GiverCommodityID,
		ValueMinutes:   s.RealisedValue,
		ExchangeID:     s.ExchangeID,
	}
	receiver = CircuitLeg{
		AgentID:        s.ReceiverID,
		Kind:           CircuitLegPurchase,
		CounterpartyID: s.GiverID,
		CommodityID:    s.GiverCommodityID,
		ValueMinutes:   s.RealisedValue,
		ExchangeID:     s.ExchangeID,
	}
	return giver, receiver
}

// Validate checks that both derived legs are structurally sound.
func (s ExchangeSettlement) Validate() error {
	giver, receiver := s.Legs()
	if err := giver.Validate(); err != nil {
		return err
	}
	return receiver.Validate()
}
