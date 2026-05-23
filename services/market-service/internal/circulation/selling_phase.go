package circulation

import "time"

// SellingOutcome records the result of a C—M metamorphosis.
type SellingOutcome string

const (
	OutcomeSold    SellingOutcome = "sold"
	OutcomeSpoiled SellingOutcome = "spoiled"
	OutcomePartial SellingOutcome = "partial"
	OutcomePending SellingOutcome = "pending"
)

// SellingPhaseID is a 96-bit hex identifier for a SellingPhase.
type SellingPhaseID string

func (id SellingPhaseID) IsZero() bool { return id == "" }

// NewSellingPhaseID generates a 96-bit hex identifier via crypto/rand.
func NewSellingPhaseID() SellingPhaseID { return SellingPhaseID(newHexID()) }

// SellingPhase records the C—M half of the circulation sphere.
// §: "C—M, the sale, is the most difficult part of its metamorphosis."
type SellingPhase struct {
	ID                  SellingPhaseID      `json:"id"`
	TurnoverTimeID      TurnoverTimeID      `json:"turnover_time_id"`
	IndustrialCapitalID IndustrialCapitalID `json:"industrial_capital_id"`
	CommodityCircuitID  CommodityCircuitID  `json:"commodity_circuit_id"`
	OpenedAt            time.Time           `json:"opened_at"`
	ClosedAt            *time.Time          `json:"closed_at,omitempty"`
	Outcome             SellingOutcome      `json:"outcome"`
}

// IsOpen reports whether the selling phase has not yet been closed.
func (s SellingPhase) IsOpen() bool { return s.ClosedAt == nil }

// ElapsedNanos returns the duration of the selling phase in nanoseconds.
// Returns 0 if ClosedAt is nil.
func (s SellingPhase) ElapsedNanos() int64 {
	if s.ClosedAt == nil {
		return 0
	}
	return s.ClosedAt.Sub(s.OpenedAt).Nanoseconds()
}

// BuyingPhaseID is a 96-bit hex identifier for a BuyingPhase.
type BuyingPhaseID string

func (id BuyingPhaseID) IsZero() bool { return id == "" }

// NewBuyingPhaseID generates a 96-bit hex identifier via crypto/rand.
func NewBuyingPhaseID() BuyingPhaseID { return BuyingPhaseID(newHexID()) }

// BuyingPhase records the M—C half of the circulation sphere.
// §: "M—C, the purchase, is the easier part."
type BuyingPhase struct {
	ID                  BuyingPhaseID       `json:"id"`
	TurnoverTimeID      TurnoverTimeID      `json:"turnover_time_id"`
	IndustrialCapitalID IndustrialCapitalID `json:"industrial_capital_id"`
	MoneyCircuitID      MoneyCircuitID      `json:"money_circuit_id"`
	OpenedAt            time.Time           `json:"opened_at"`
	ClosedAt            *time.Time          `json:"closed_at,omitempty"`
	MarketLocation      string              `json:"market_location"`
}

// IsOpen reports whether the buying phase has not yet been closed.
func (b BuyingPhase) IsOpen() bool { return b.ClosedAt == nil }

// ElapsedNanos returns the duration of the buying phase in nanoseconds.
// Returns 0 if ClosedAt is nil.
func (b BuyingPhase) ElapsedNanos() int64 {
	if b.ClosedAt == nil {
		return 0
	}
	return b.ClosedAt.Sub(b.OpenedAt).Nanoseconds()
}
