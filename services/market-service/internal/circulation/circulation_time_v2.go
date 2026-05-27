// Vol. II Ch. 14 — The Time of Circulation (refined).
//
// This file deepens Ch. 5's two-span model by decomposing selling-time and
// buying-time into sub-spans (buyer-search, bargaining, communication-lag,
// legal settlement), adding distance-lag relations, and introducing
// circulation-speed-improvement events (telegraph, steamship, railway,
// telephone). The England–India fixture shows how the 1870 telegraph collapsed
// a 12-month communication lag to 1 day.
package circulation

import (
	"errors"
	"fmt"
	"time"
)

// Sentinel errors for Ch. 14.
var (
	// ErrLagExceedsMaximum is returned when CommunicationLagMinutes exceeds
	// the bound implied by the DistanceLagRelation for the market.
	ErrLagExceedsMaximum = errors.New("circulation: communication lag exceeds maximum for distance")

	// ErrSubSpanMismatch is returned when the sub-span sum does not equal
	// the recorded total phase duration.
	ErrSubSpanMismatch = errors.New("circulation: sub-spans do not sum to total phase duration")
)

// MarketID is an opaque reference to a market entity.
type MarketID string

// --- SellingPhaseDetailed ---------------------------------------------------

// SellingPhaseDetailedID is a 96-bit hex identifier.
type SellingPhaseDetailedID string

func (id SellingPhaseDetailedID) IsZero() bool { return id == "" }

// NewSellingPhaseDetailedID generates a 96-bit hex identifier via crypto/rand.
func NewSellingPhaseDetailedID() SellingPhaseDetailedID {
	return SellingPhaseDetailedID(newHexID())
}

// SellingPhaseDetailed decomposes a Ch. 5 SellingPhase into four sub-spans.
// The sub-spans must sum to TotalMinutes (the elapsed duration of the parent
// SellingPhase).
// §: "the time of sale divides into: search for buyers, bargaining,
//
//	communication-lag, legal settlement."
type SellingPhaseDetailed struct {
	ID                      SellingPhaseDetailedID `json:"id"`
	SellingPhaseID          SellingPhaseID         `json:"selling_phase_id"`
	TurnoverTimeID          TurnoverTimeID         `json:"turnover_time_id"`
	IndustrialCapitalID     IndustrialCapitalID    `json:"industrial_capital_id"`
	MarketID                MarketID               `json:"market_id"`
	MarketDistanceMeters    int64                  `json:"market_distance_meters"`
	CommunicationLagMinutes int64                  `json:"communication_lag_minutes"`
	BuyerSearchMinutes      int64                  `json:"buyer_search_minutes"`
	BargainingMinutes       int64                  `json:"bargaining_minutes"`
	LegalSettlementMinutes  int64                  `json:"legal_settlement_minutes"`
	// TotalMinutes is the elapsed duration of the parent SellingPhase.
	// Sub-spans must sum to this value.
	TotalMinutes int64     `json:"total_minutes"`
	CreatedAt    time.Time `json:"created_at"`
}

// ValidateSubSpans returns ErrSubSpanMismatch when the sub-spans do not sum
// to TotalMinutes.
func (s SellingPhaseDetailed) ValidateSubSpans() error {
	sum := s.BuyerSearchMinutes + s.BargainingMinutes +
		s.CommunicationLagMinutes + s.LegalSettlementMinutes
	if sum != s.TotalMinutes {
		return fmt.Errorf("%w: got %d min, want %d min", ErrSubSpanMismatch, sum, s.TotalMinutes)
	}
	return nil
}

// ValidateLagBound returns ErrLagExceedsMaximum when CommunicationLagMinutes
// exceeds lagMinutesPerKm × (MarketDistanceMeters / 1000).
func (s SellingPhaseDetailed) ValidateLagBound(lagMinutesPerKm int64) error {
	maxLag := lagMinutesPerKm * s.MarketDistanceMeters / 1000
	if s.CommunicationLagMinutes > maxLag {
		return fmt.Errorf("%w: lag %d min > max %d min (%.0f km × %d min/km)",
			ErrLagExceedsMaximum,
			s.CommunicationLagMinutes,
			maxLag,
			float64(s.MarketDistanceMeters)/1000,
			lagMinutesPerKm,
		)
	}
	return nil
}

// --- BuyingPhaseDetailed ----------------------------------------------------

// BuyingPhaseDetailedID is a 96-bit hex identifier.
type BuyingPhaseDetailedID string

func (id BuyingPhaseDetailedID) IsZero() bool { return id == "" }

// NewBuyingPhaseDetailedID generates a 96-bit hex identifier via crypto/rand.
func NewBuyingPhaseDetailedID() BuyingPhaseDetailedID {
	return BuyingPhaseDetailedID(newHexID())
}

// BuyingPhaseDetailed decomposes a Ch. 5 BuyingPhase into four sub-spans.
// §: "the time of purchase divides into: supplier-search, order-processing,
//
//	communication-lag, legal settlement."
type BuyingPhaseDetailed struct {
	ID                      BuyingPhaseDetailedID `json:"id"`
	BuyingPhaseID           BuyingPhaseID         `json:"buying_phase_id"`
	TurnoverTimeID          TurnoverTimeID         `json:"turnover_time_id"`
	IndustrialCapitalID     IndustrialCapitalID    `json:"industrial_capital_id"`
	MarketID                MarketID               `json:"market_id"`
	MarketDistanceMeters    int64                  `json:"market_distance_meters"`
	CommunicationLagMinutes int64                  `json:"communication_lag_minutes"`
	SupplierSearchMinutes   int64                  `json:"supplier_search_minutes"`
	OrderProcessingMinutes  int64                  `json:"order_processing_minutes"`
	LegalSettlementMinutes  int64                  `json:"legal_settlement_minutes"`
	TotalMinutes            int64                  `json:"total_minutes"`
	CreatedAt               time.Time              `json:"created_at"`
}

// ValidateSubSpans returns ErrSubSpanMismatch when the sub-spans do not sum
// to TotalMinutes.
func (b BuyingPhaseDetailed) ValidateSubSpans() error {
	sum := b.SupplierSearchMinutes + b.OrderProcessingMinutes +
		b.CommunicationLagMinutes + b.LegalSettlementMinutes
	if sum != b.TotalMinutes {
		return fmt.Errorf("%w: got %d min, want %d min", ErrSubSpanMismatch, sum, b.TotalMinutes)
	}
	return nil
}

// ValidateLagBound returns ErrLagExceedsMaximum when CommunicationLagMinutes
// exceeds lagMinutesPerKm × (MarketDistanceMeters / 1000).
func (b BuyingPhaseDetailed) ValidateLagBound(lagMinutesPerKm int64) error {
	maxLag := lagMinutesPerKm * b.MarketDistanceMeters / 1000
	if b.CommunicationLagMinutes > maxLag {
		return fmt.Errorf("%w: lag %d min > max %d min (%.0f km × %d min/km)",
			ErrLagExceedsMaximum,
			b.CommunicationLagMinutes,
			maxLag,
			float64(b.MarketDistanceMeters)/1000,
			lagMinutesPerKm,
		)
	}
	return nil
}

// --- DistanceLagRelation ----------------------------------------------------

// DistanceLagRelationID is a 96-bit hex identifier.
type DistanceLagRelationID string

func (id DistanceLagRelationID) IsZero() bool { return id == "" }

// NewDistanceLagRelationID generates a 96-bit hex identifier via crypto/rand.
func NewDistanceLagRelationID() DistanceLagRelationID {
	return DistanceLagRelationID(newHexID())
}

// DistanceLagRelation records the median communication lag per kilometre for a
// market route, given the dominant transmission medium of the era.
// §: "the greater the distance the buyer's market is from the producer, the
//
//	longer is the time of circulation."
type DistanceLagRelation struct {
	ID              DistanceLagRelationID `json:"id"`
	MarketID        MarketID              `json:"market_id"`
	DistanceMeters  int64                 `json:"distance_meters"`
	LagMinutesPerKm int64                 `json:"lag_minutes_per_km"`
	CreatedAt       time.Time             `json:"created_at"`
}

// MaxLagMinutes returns LagMinutesPerKm × (DistanceMeters / 1000).
func (d DistanceLagRelation) MaxLagMinutes() int64 {
	return d.LagMinutesPerKm * d.DistanceMeters / 1000
}

// --- CirculationSpeedImprovement --------------------------------------------

// CirculationSpeedImprovementID is a 96-bit hex identifier.
type CirculationSpeedImprovementID string

func (id CirculationSpeedImprovementID) IsZero() bool { return id == "" }

// NewCirculationSpeedImprovementID generates a 96-bit hex identifier via crypto/rand.
func NewCirculationSpeedImprovementID() CirculationSpeedImprovementID {
	return CirculationSpeedImprovementID(newHexID())
}

// CirculationSpeedReason names the means of transport or communication that
// caused a change in the lag per kilometre.
type CirculationSpeedReason string

const (
	ReasonTelegraph CirculationSpeedReason = "telegraph"
	ReasonSteamship CirculationSpeedReason = "steamship"
	ReasonTelephone CirculationSpeedReason = "telephone"
	ReasonRailway   CirculationSpeedReason = "railway"
	ReasonOther     CirculationSpeedReason = "other"
)

// CirculationSpeedImprovement records a historical event that changed the
// communication-lag rate for a market route.
// §: "with every improvement in the means of communication and transport,
//
//	the time of circulation diminishes."
type CirculationSpeedImprovement struct {
	ID                 CirculationSpeedImprovementID `json:"id"`
	MarketID           MarketID                      `json:"market_id"`
	OldLagMinutesPerKm int64                         `json:"old_lag_minutes_per_km"`
	NewLagMinutesPerKm int64                         `json:"new_lag_minutes_per_km"`
	Reason             CirculationSpeedReason         `json:"reason"`
	EffectiveAt        time.Time                     `json:"effective_at"`
	CreatedAt          time.Time                     `json:"created_at"`
}

// LagReduction returns OldLagMinutesPerKm - NewLagMinutesPerKm.
// A positive result means circulation time decreased (capital set free).
func (c CirculationSpeedImprovement) LagReduction() int64 {
	return c.OldLagMinutesPerKm - c.NewLagMinutesPerKm
}

// --- AnnualSurplusPenalty ---------------------------------------------------

// AnnualSurplusPenaltyID is a 96-bit hex identifier.
type AnnualSurplusPenaltyID string

func (id AnnualSurplusPenaltyID) IsZero() bool { return id == "" }

// NewAnnualSurplusPenaltyID generates a 96-bit hex identifier via crypto/rand.
func NewAnnualSurplusPenaltyID() AnnualSurplusPenaltyID {
	return AnnualSurplusPenaltyID(newHexID())
}

// AnnualSurplusPenalty records the penalty in surplus-value production caused
// by circulation time. IdealAnnualRateBasisPoints is the rate achievable with
// zero circulation time; ActualAnnualRateBasisPoints is the realised rate given
// actual circulation time. Full calculation is deferred to Vol. II Ch. 16.
// §: "the longer the time of circulation … the less surplus-value is produced."
type AnnualSurplusPenalty struct {
	ID                          AnnualSurplusPenaltyID `json:"id"`
	IndustrialCapitalID         IndustrialCapitalID    `json:"industrial_capital_id"`
	// Period is a label for the annual accounting period (e.g. "1871").
	Period                      string                 `json:"period"`
	IdealAnnualRateBasisPoints  int64                  `json:"ideal_annual_rate_basis_points"`
	ActualAnnualRateBasisPoints int64                  `json:"actual_annual_rate_basis_points"`
	PenaltyBasisPoints          int64                  `json:"penalty_basis_points"`
}

// ComputePenalty returns max(0, IdealAnnualRateBasisPoints - ActualAnnualRateBasisPoints).
// PenaltyBasisPoints is always ≥ 0. Full annual-rate calculation lives in Ch. 16.
func (a AnnualSurplusPenalty) ComputePenalty() int64 {
	p := a.IdealAnnualRateBasisPoints - a.ActualAnnualRateBasisPoints
	if p < 0 {
		return 0
	}
	return p
}
