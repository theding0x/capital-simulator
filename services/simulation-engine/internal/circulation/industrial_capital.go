// Package circulation — Vol. II Ch. 4: The Three Formulas of the Circuit.
// IndustrialCapital is the synthesis that shows every industrial capital
// occupying all three stages (M, P, C) simultaneously — Marx's Nebeneinander.
package circulation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// newHexID returns a 96-bit (12-byte) random identifier as a hex string.
func newHexID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(b)
}

// --- Sentinel errors --------------------------------------------------------

var (
	ErrNonCapitalCounterparty    = errors.New("circulation: metamorphosis interlock requires both sides to be capital-form")
	ErrCreditWithoutMoneyEconomy = errors.New("circulation: cannot transition to EconomyCredit without first establishing EconomyMoney")
	ErrSinkingFundMismatch       = errors.New("circulation: annual_payment_pence * lifetime_years must equal fixed_capital_pence")
	ErrNonZeroDirectionPence     = errors.New("circulation: direction_pence must be non-zero for a value-revolution event")
)

// --- CapitalStage -----------------------------------------------------------

// CapitalStage names the three functional forms of industrial capital.
type CapitalStage string

const (
	StageMoney      CapitalStage = "money"
	StageProduction CapitalStage = "production"
	StageCommodity  CapitalStage = "commodity"
)

// --- IndustrialCapitalStatus ------------------------------------------------

// IndustrialCapitalStatus tracks whether the capital is running or halted by
// a persistent stagnation block.
type IndustrialCapitalStatus string

const (
	StatusActive IndustrialCapitalStatus = "active"
	StatusHalted IndustrialCapitalStatus = "halted"
)

// --- EconomyMode ------------------------------------------------------------

// EconomyMode names the mode of the economy in which this industrial capital
// operates. EconomyCredit presupposes EconomyMoney.
type EconomyMode string

const (
	EconomyNatural EconomyMode = "natural"
	EconomyMoney   EconomyMode = "money"
	EconomyCredit  EconomyMode = "credit"
)

// --- StageBlockReason -------------------------------------------------------

// StageBlockReason names the cause of a stagnation block at one stage.
type StageBlockReason string

const (
	BlockUnsoldCommodity          StageBlockReason = "unsold_commodity"
	BlockMissingMeansOfProduction StageBlockReason = "missing_means_of_production"
	BlockDelayedMoneyReturn       StageBlockReason = "delayed_money_return"
)

// --- MPSourceOrigin ---------------------------------------------------------

// MPSourceOrigin names the origin of means of production entering a circuit.
// Extends Ch. 3's MeansOfProductionSourceKind to include non-capitalist
// sources. Only OriginCapitalist and OriginMerchant may appear in a
// MetamorphosisInterlock.
type MPSourceOrigin string

const (
	OriginCapitalist            MPSourceOrigin = "capitalist"
	OriginNonCapitalistProducer MPSourceOrigin = "non_capitalist_producer"
	OriginMerchant              MPSourceOrigin = "merchant"
	OriginIndependentLabourer   MPSourceOrigin = "independent_labourer"
	OriginImport                MPSourceOrigin = "import"
)

// isCapitalForm reports whether the origin is a capital-form seller.
func (o MPSourceOrigin) isCapitalForm() bool {
	return o == OriginCapitalist || o == OriginMerchant
}

// --- ValueRevolutionAffecting -----------------------------------------------

// ValueRevolutionAffecting names which element of the circuit is repriced.
type ValueRevolutionAffecting string

const (
	AffectingMeansOfProduction ValueRevolutionAffecting = "means_of_production"
	AffectingLabourPower       ValueRevolutionAffecting = "labour_power"
	AffectingOutput            ValueRevolutionAffecting = "output"
)

// --- ID types ---------------------------------------------------------------

// IndustrialCapitalID is a 96-bit hex identifier.
type IndustrialCapitalID string

func (id IndustrialCapitalID) IsZero() bool { return id == "" }

// NewIndustrialCapitalID generates a 96-bit hex identifier via crypto/rand.
func NewIndustrialCapitalID() IndustrialCapitalID {
	return IndustrialCapitalID(newHexID())
}

// CapitalPartID is a 96-bit hex identifier for a CapitalPart.
type CapitalPartID string

func (id CapitalPartID) IsZero() bool { return id == "" }

// NewCapitalPartID generates a 96-bit hex identifier via crypto/rand.
func NewCapitalPartID() CapitalPartID { return CapitalPartID(newHexID()) }

// StageDistributionID is a 96-bit hex identifier.
type StageDistributionID string

func (id StageDistributionID) IsZero() bool { return id == "" }

// NewStageDistributionID generates a 96-bit hex identifier via crypto/rand.
func NewStageDistributionID() StageDistributionID { return StageDistributionID(newHexID()) }

// StageBlockID is a 96-bit hex identifier.
type StageBlockID string

func (id StageBlockID) IsZero() bool { return id == "" }

// NewStageBlockID generates a 96-bit hex identifier via crypto/rand.
func NewStageBlockID() StageBlockID { return StageBlockID(newHexID()) }

// ValueRevolutionEventID is a 96-bit hex identifier.
type ValueRevolutionEventID string

func (id ValueRevolutionEventID) IsZero() bool { return id == "" }

// NewValueRevolutionEventID generates a 96-bit hex identifier via crypto/rand.
func NewValueRevolutionEventID() ValueRevolutionEventID {
	return ValueRevolutionEventID(newHexID())
}

// MoneyCapitalSetFreeID is a 96-bit hex identifier.
type MoneyCapitalSetFreeID string

func (id MoneyCapitalSetFreeID) IsZero() bool { return id == "" }

// NewMoneyCapitalSetFreeID generates a 96-bit hex identifier via crypto/rand.
func NewMoneyCapitalSetFreeID() MoneyCapitalSetFreeID {
	return MoneyCapitalSetFreeID(newHexID())
}

// MoneyCapitalTiedUpID is a 96-bit hex identifier.
type MoneyCapitalTiedUpID string

func (id MoneyCapitalTiedUpID) IsZero() bool { return id == "" }

// NewMoneyCapitalTiedUpID generates a 96-bit hex identifier via crypto/rand.
func NewMoneyCapitalTiedUpID() MoneyCapitalTiedUpID {
	return MoneyCapitalTiedUpID(newHexID())
}

// MetamorphosisInterlockID is a 96-bit hex identifier.
type MetamorphosisInterlockID string

func (id MetamorphosisInterlockID) IsZero() bool { return id == "" }

// NewMetamorphosisInterlockID generates a 96-bit hex identifier via crypto/rand.
func NewMetamorphosisInterlockID() MetamorphosisInterlockID {
	return MetamorphosisInterlockID(newHexID())
}

// SupplyDemandImbalanceID is a 96-bit hex identifier.
type SupplyDemandImbalanceID string

func (id SupplyDemandImbalanceID) IsZero() bool { return id == "" }

// NewSupplyDemandImbalanceID generates a 96-bit hex identifier via crypto/rand.
func NewSupplyDemandImbalanceID() SupplyDemandImbalanceID {
	return SupplyDemandImbalanceID(newHexID())
}

// SinkingFundID is a 96-bit hex identifier.
type SinkingFundID string

func (id SinkingFundID) IsZero() bool { return id == "" }

// NewSinkingFundID generates a 96-bit hex identifier via crypto/rand.
func NewSinkingFundID() SinkingFundID { return SinkingFundID(newHexID()) }

// --- IndustrialCapital ------------------------------------------------------

// IndustrialCapital is the synthesising aggregate for Vol. II Ch. 4. It
// references one MoneyCircuit (Ch. 1), one ProductiveCircuit (Ch. 2), and one
// CommodityCircuit (Ch. 3) as views of the same advanced capital. Parts of
// the capital occupy all three stages simultaneously (Nebeneinander).
// §: "every individual industrial capital is present simultaneously in all
// three circuits."
type IndustrialCapital struct {
	ID                       IndustrialCapitalID     `json:"id"`
	AgentID                  string                  `json:"agent_id"`
	MoneyCircuitID           MoneyCircuitID          `json:"money_circuit_id"`
	ProductiveCircuitID      ProductiveCircuitID     `json:"productive_circuit_id"`
	CommodityCircuitID       CommodityCircuitID      `json:"commodity_circuit_id"`
	TotalPence               Pence                   `json:"total_pence"`
	EconomyMode              EconomyMode             `json:"economy_mode"`
	StagnationToleranceTicks int64                   `json:"stagnation_tolerance_ticks"`
	Status                   IndustrialCapitalStatus `json:"status"`
	// Latest is the most recent StageDistribution, populated on Get.
	Latest *StageDistribution `json:"latest_distribution,omitempty"`
	// OpenBlocks are the currently open stagnation blocks, populated on Get.
	OpenBlocks []StageBlock `json:"open_blocks,omitempty"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
}

// Validate checks required fields.
func (ic IndustrialCapital) Validate() error {
	if ic.TotalPence <= 0 {
		return fmt.Errorf("circulation: total_pence must be positive")
	}
	switch ic.EconomyMode {
	case EconomyNatural, EconomyMoney, EconomyCredit:
	default:
		return fmt.Errorf("circulation: unknown economy_mode %q", ic.EconomyMode)
	}
	if ic.StagnationToleranceTicks < 0 {
		return fmt.Errorf("circulation: stagnation_tolerance_ticks must be non-negative")
	}
	return nil
}

// --- CapitalPart ------------------------------------------------------------

// CapitalPart is a divided share of the industrial capital currently present
// at one stage. Several parts of the same capital exist simultaneously, each
// at a different stage — Marx's Nebeneinander.
type CapitalPart struct {
	ID                  CapitalPartID       `json:"id"`
	IndustrialCapitalID IndustrialCapitalID `json:"industrial_capital_id"`
	Pence               Pence               `json:"pence"`
	Stage               CapitalStage        `json:"stage"`
	EnteredStageAt      time.Time           `json:"entered_stage_at"`
}

// --- StageDistribution ------------------------------------------------------

// StageDistribution is a snapshot of the capital's split across the three
// stages at a point in time.
// Invariant: MoneyPence + ProductionPence + CommodityPence == IndustrialCapital.TotalPence.
type StageDistribution struct {
	ID                  StageDistributionID `json:"id"`
	IndustrialCapitalID IndustrialCapitalID `json:"industrial_capital_id"`
	At                  time.Time           `json:"at"`
	MoneyPence          Pence               `json:"money_pence"`
	ProductionPence     Pence               `json:"production_pence"`
	CommodityPence      Pence               `json:"commodity_pence"`
}

// Total returns MoneyPence + ProductionPence + CommodityPence.
func (sd StageDistribution) Total() Pence {
	return sd.MoneyPence + sd.ProductionPence + sd.CommodityPence
}

// Validate checks that the three pence fields are non-negative and that their
// sum equals expectedTotal (when expectedTotal > 0).
func (sd StageDistribution) Validate(expectedTotal Pence) error {
	if sd.MoneyPence < 0 || sd.ProductionPence < 0 || sd.CommodityPence < 0 {
		return fmt.Errorf("circulation: stage distribution pence must be non-negative")
	}
	if expectedTotal > 0 && sd.Total() != expectedTotal {
		return fmt.Errorf("circulation: stage distribution total %d != industrial capital total %d",
			sd.Total(), expectedTotal)
	}
	return nil
}

// --- StageBlock -------------------------------------------------------------

// StageBlock records a stagnation at one stage. While ClosedAt is nil, new
// capital parts trying to enter Stage are deferred.
// §: "every stagnation in succession carries disorder into co-existence"
type StageBlock struct {
	ID                  StageBlockID        `json:"id"`
	IndustrialCapitalID IndustrialCapitalID `json:"industrial_capital_id"`
	Stage               CapitalStage        `json:"stage"`
	Reason              StageBlockReason    `json:"reason"`
	OpenedAt            time.Time           `json:"opened_at"`
	ClosedAt            *time.Time          `json:"closed_at,omitempty"`
}

// IsOpen reports whether the block has not been closed.
func (b StageBlock) IsOpen() bool { return b.ClosedAt == nil }

// --- ValueRevolutionEvent ---------------------------------------------------

// ValueRevolutionEvent records a price change on means of production, labour
// power, or output. DirectionPence > 0 means prices rose (money tied up);
// DirectionPence < 0 means prices fell (money set free).
//
// Ch. 15 extends the record with BasisPointsChange (the percentage change
// expressed in basis points, e.g., -2000 for a 20% fall) and
// RevaluedCapitalPence (the new advance required after the revolution, derived
// by DeriveRevaluedCapital and stored here for quick reference).
// §: "value (prices) fall — a part of the money-capital existing hitherto
// is set free … The opposite occurs if the value of the elements of
// replacement of a commodity-capital increases."
type ValueRevolutionEvent struct {
	ID                   ValueRevolutionEventID   `json:"id"`
	IndustrialCapitalID  IndustrialCapitalID      `json:"industrial_capital_id"`
	Affecting            ValueRevolutionAffecting `json:"affecting"`
	DirectionPence       Pence                    `json:"direction_pence"`
	OccurredAt           time.Time                `json:"occurred_at"`
	BasisPointsChange    int64                    `json:"basis_points_change"`
	RevaluedCapitalPence Pence                    `json:"revalued_capital_pence"`
}

// Validate checks that DirectionPence is non-zero.
func (e ValueRevolutionEvent) Validate() error {
	if e.DirectionPence == 0 {
		return ErrNonZeroDirectionPence
	}
	return nil
}

// MoneyCapitalSetFree records the money capital freed when input prices fall.
type MoneyCapitalSetFree struct {
	ID                     MoneyCapitalSetFreeID  `json:"id"`
	ValueRevolutionEventID ValueRevolutionEventID `json:"value_revolution_event_id"`
	IndustrialCapitalID    IndustrialCapitalID    `json:"industrial_capital_id"`
	AmountPence            Pence                  `json:"amount_pence"`
}

// MoneyCapitalTiedUp records the additional money capital required when input
// prices rise.
type MoneyCapitalTiedUp struct {
	ID                     MoneyCapitalTiedUpID   `json:"id"`
	ValueRevolutionEventID ValueRevolutionEventID `json:"value_revolution_event_id"`
	IndustrialCapitalID    IndustrialCapitalID    `json:"industrial_capital_id"`
	AmountPence            Pence                  `json:"amount_pence"`
}

// --- ValueRevolutionResult --------------------------------------------------

// ValueRevolutionResult bundles the event with the emitted record.
type ValueRevolutionResult struct {
	Event   ValueRevolutionEvent `json:"event"`
	SetFree *MoneyCapitalSetFree `json:"set_free,omitempty"`
	TiedUp  *MoneyCapitalTiedUp  `json:"tied_up,omitempty"`
}

// ApplyValueRevolution creates the correct SetFree or TiedUp record from an
// event and returns the result.
// DirectionPence > 0 → TiedUp; DirectionPence < 0 → SetFree.
func ApplyValueRevolution(evt ValueRevolutionEvent) (ValueRevolutionResult, error) {
	if err := evt.Validate(); err != nil {
		return ValueRevolutionResult{}, err
	}
	res := ValueRevolutionResult{Event: evt}
	if evt.DirectionPence > 0 {
		res.TiedUp = &MoneyCapitalTiedUp{
			ID:                     NewMoneyCapitalTiedUpID(),
			ValueRevolutionEventID: evt.ID,
			IndustrialCapitalID:    evt.IndustrialCapitalID,
			AmountPence:            evt.DirectionPence,
		}
	} else {
		res.SetFree = &MoneyCapitalSetFree{
			ID:                     NewMoneyCapitalSetFreeID(),
			ValueRevolutionEventID: evt.ID,
			IndustrialCapitalID:    evt.IndustrialCapitalID,
			AmountPence:            -evt.DirectionPence,
		}
	}
	return res, nil
}

// --- MetamorphosisInterlock -------------------------------------------------

// MetamorphosisInterlock records that one industrial capital's M—C is the
// counterpart of another's C—M, both sides being in capital-form.
// §: "it is always M—C on one side and C—M on the other, but not always an
// intermingling of metamorphoses of capitals."
type MetamorphosisInterlock struct {
	ID                        MetamorphosisInterlockID `json:"id"`
	BuyerIndustrialCapitalID  IndustrialCapitalID      `json:"buyer_industrial_capital_id"`
	SellerIndustrialCapitalID IndustrialCapitalID      `json:"seller_industrial_capital_id"`
	SellerMPSourceOrigin      MPSourceOrigin           `json:"seller_mp_source_origin"`
	Pence                     Pence                    `json:"pence"`
	OccurredAt                time.Time                `json:"occurred_at"`
}

// Validate returns ErrNonCapitalCounterparty if the seller's origin is not
// capital-form.
func (mi MetamorphosisInterlock) Validate() error {
	if !mi.SellerMPSourceOrigin.isCapitalForm() {
		return ErrNonCapitalCounterparty
	}
	if mi.Pence <= 0 {
		return fmt.Errorf("circulation: interlock pence must be positive")
	}
	return nil
}

// --- OrganicComposition -----------------------------------------------------

// OrganicComposition is the value-composition of capital: c (constant) and
// v (variable). RatioBasisPoints returns 10000 * c / v (integer, no float64).
type OrganicComposition struct {
	ConstantPence Pence `json:"constant_pence"`
	VariablePence Pence `json:"variable_pence"`
}

// Total returns ConstantPence + VariablePence.
func (oc OrganicComposition) Total() Pence { return oc.ConstantPence + oc.VariablePence }

// RatioBasisPoints returns 10000 * c / v. Returns 0 when VariablePence == 0.
func (oc OrganicComposition) RatioBasisPoints() int64 {
	if oc.VariablePence == 0 {
		return 0
	}
	return int64(10000 * oc.ConstantPence / oc.VariablePence)
}

// --- SupplyDemandImbalance --------------------------------------------------

// SupplyDemandImbalance captures the structural gap between supply and demand
// generated by surplus-value for one industrial capital over one period.
// §: "if the composition of his commodity-capital is 80c + 20v + 20s, his
// demand is equal to 80c + 20v, hence … one-fifth smaller than his supply."
type SupplyDemandImbalance struct {
	ID                  SupplyDemandImbalanceID `json:"id"`
	IndustrialCapitalID IndustrialCapitalID     `json:"industrial_capital_id"`
	Period              string                  `json:"period"`
	DemandPence         Pence                   `json:"demand_pence"`
	SupplyPence         Pence                   `json:"supply_pence"`
	ExcessPence         Pence                   `json:"excess_pence"`
}

// ComputeSupplyDemandImbalance derives a SupplyDemandImbalance from an
// OrganicComposition and a surplus-value magnitude.
// demand = c + v; supply = c + v + s; excess = s.
func ComputeSupplyDemandImbalance(ic IndustrialCapitalID, period string, oc OrganicComposition, surplusPence Pence) SupplyDemandImbalance {
	demand := oc.Total()
	supply := demand + surplusPence
	return SupplyDemandImbalance{
		ID:                  NewSupplyDemandImbalanceID(),
		IndustrialCapitalID: ic,
		Period:              period,
		DemandPence:         demand,
		SupplyPence:         supply,
		ExcessPence:         surplusPence,
	}
}

// AggregateSupplyDemandImbalance is the class-level total across all capitals.
type AggregateSupplyDemandImbalance struct {
	Period      string `json:"period"`
	DemandPence Pence  `json:"demand_pence"`
	SupplyPence Pence  `json:"supply_pence"`
	ExcessPence Pence  `json:"excess_pence"`
}

// --- SinkingFund ------------------------------------------------------------

// SinkingFund models the annual depreciation reserve for fixed capital.
// §: "he pays every year one-tenth, or £400, into a sinking fund."
// Invariant: AnnualPaymentPence * LifetimeYears == FixedCapitalPence.
type SinkingFund struct {
	ID                  SinkingFundID       `json:"id"`
	IndustrialCapitalID IndustrialCapitalID `json:"industrial_capital_id"`
	FixedCapitalPence   Pence               `json:"fixed_capital_pence"`
	LifetimeYears       int64               `json:"lifetime_years"`
	AnnualPaymentPence  Pence               `json:"annual_payment_pence"`
	AccumulatedPence    Pence               `json:"accumulated_pence"`
}

// Validate checks the sinking-fund invariant.
func (sf SinkingFund) Validate() error {
	if sf.FixedCapitalPence <= 0 {
		return fmt.Errorf("circulation: fixed_capital_pence must be positive")
	}
	if sf.LifetimeYears <= 0 {
		return fmt.Errorf("circulation: lifetime_years must be positive")
	}
	expected := sf.FixedCapitalPence / Pence(sf.LifetimeYears)
	if sf.AnnualPaymentPence != expected {
		return ErrSinkingFundMismatch
	}
	return nil
}

// Tick books one annual payment. Returns the updated fund.
func (sf SinkingFund) Tick() SinkingFund {
	sf.AccumulatedPence += sf.AnnualPaymentPence
	return sf
}
