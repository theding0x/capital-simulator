// Package circulation models the time dimension of capital's circuit —
// the period during which capital moves through the selling and buying phases
// without producing surplus-value. Vol. II Ch. 5.
package circulation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Pence is the canonical monetary unit in the simulation (1 GBP = 240 pence).
type Pence int64

// Sentinel errors.
var (
	ErrConcurrentProductionAndCirculation = errors.New("circulation: production time and circulation time are mutually exclusive")
	ErrSellingPhaseAlreadyOpen            = errors.New("circulation: selling phase already open for this turnover time")
	ErrBuyingPhaseAlreadyOpen             = errors.New("circulation: buying phase already open for this turnover time")
	ErrNoOpenSellingPhase                 = errors.New("circulation: no open selling phase to close")
	ErrNoOpenBuyingPhase                  = errors.New("circulation: no open buying phase to close")
)

// --- ID types ---------------------------------------------------------------

// TurnoverTimeID is a 96-bit hex identifier for a TurnoverTime record.
type TurnoverTimeID string

func (id TurnoverTimeID) IsZero() bool { return id == "" }

// NewTurnoverTimeID generates a 96-bit hex identifier via crypto/rand.
func NewTurnoverTimeID() TurnoverTimeID {
	return TurnoverTimeID(newHexID())
}

// IndustrialCapitalID is an opaque reference to an IndustrialCapital record
// managed by simulation-engine (Vol. II Ch. 4).
type IndustrialCapitalID string

func (id IndustrialCapitalID) IsZero() bool { return id == "" }

// CommodityCircuitID is an opaque reference to a CommodityCircuit in
// simulation-engine (Vol. II Ch. 3).
type CommodityCircuitID string

// MoneyCircuitID is an opaque reference to a MoneyCircuit in
// simulation-engine (Vol. II Ch. 1).
type MoneyCircuitID string

// CommodityID is an opaque reference to a Commodity in commodity-service.
type CommodityID string

// --- CirculationPhase -------------------------------------------------------

// CirculationPhase names the two metamorphosis-phases within the time of
// circulation.
type CirculationPhase string

const (
	PhaseSelling CirculationPhase = "selling" // C—M
	PhaseBuying  CirculationPhase = "buying"  // M—C
)

// --- ProductionTime ---------------------------------------------------------

// ProductionTime decomposes the production sphere into four mutually exclusive
// sub-spans. Their sum equals the total production time for one circuit.
// §: "the time of production comprises: 1) labour-time proper, 2) interruptions,
// 3) time held in readiness as latent productive capital, 4) natural-process spans."
type ProductionTime struct {
	LabourTimeNanos         int64 `json:"labour_time_nanos"`
	LabourInterruptionNanos int64 `json:"labour_interruption_nanos"`
	LatentNanos             int64 `json:"latent_nanos"`
	NaturalProcessNanos     int64 `json:"natural_process_nanos"`
}

// TotalNanos returns the sum of all four sub-spans.
func (p ProductionTime) TotalNanos() int64 {
	return p.LabourTimeNanos + p.LabourInterruptionNanos + p.LatentNanos + p.NaturalProcessNanos
}

// --- CirculationTime --------------------------------------------------------

// CirculationTime decomposes the circulation sphere into selling-time (C—M)
// and buying-time (M—C).
// §: "time of circulation is divided into C—M and M—C."
type CirculationTime struct {
	SellingTimeNanos int64 `json:"selling_time_nanos"`
	BuyingTimeNanos  int64 `json:"buying_time_nanos"`
}

// TotalNanos returns SellingTimeNanos + BuyingTimeNanos.
func (c CirculationTime) TotalNanos() int64 {
	return c.SellingTimeNanos + c.BuyingTimeNanos
}

// --- TurnoverTime -----------------------------------------------------------

// TurnoverTime is the total circuit duration for one industrial capital:
// ProductionTime + CirculationTime. It does not compute turnover-number or
// annual surplus-value rates (deferred to Vol. II Ch. 7 and Ch. 16).
// §: "total time = time of production + time of circulation."
type TurnoverTime struct {
	ID                  TurnoverTimeID      `json:"id"`
	IndustrialCapitalID IndustrialCapitalID `json:"industrial_capital_id"`
	Production          ProductionTime      `json:"production"`
	Circulation         CirculationTime     `json:"circulation"`
	// CirculationOpen is true when a selling or buying phase is currently open.
	// While CirculationOpen is true, adding production sub-spans is forbidden.
	CirculationOpen bool      `json:"circulation_open"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TotalNanos returns the sum of production time and circulation time.
func (t TurnoverTime) TotalNanos() int64 {
	return t.Production.TotalNanos() + t.Circulation.TotalNanos()
}

// ActiveFractionBasisPoints returns 10000 × ProductionTime / TotalNanos.
// Returns 0 when TotalNanos() == 0 to avoid division by zero.
// §: "circulation time limits production time and hence surplus-value production."
func (t TurnoverTime) ActiveFractionBasisPoints() int64 {
	total := t.TotalNanos()
	if total == 0 {
		return 0
	}
	return 10000 * t.Production.TotalNanos() / total
}

// Validate checks required fields.
func (t TurnoverTime) Validate() error {
	if t.IndustrialCapitalID.IsZero() {
		return fmt.Errorf("circulation: industrial_capital_id is required")
	}
	return nil
}

// --- helpers ----------------------------------------------------------------

func newHexID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(b)
}
