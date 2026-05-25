// Vol. II Ch. 8 — Fixed Capital and Circulating Capital.
// Package composition models the fixed/circulating distinction within productive
// capital: instruments of labour (fixed) vs raw, auxiliary, and labour-power
// (circulating). The distinction applies ONLY to productive capital, never to
// money-capital or commodity-capital (Marx, §I).
package composition

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Pence is the canonical monetary unit.
type Pence int64

// CapitalComponentID is the unique identifier for a CapitalComponent.
type CapitalComponentID string

// FixedCapitalItemID is the unique identifier for a FixedCapitalItem.
type FixedCapitalItemID string

// FixedCapitalSubcomponentID is the unique identifier for a FixedCapitalSubcomponent.
type FixedCapitalSubcomponentID string

// WearAndTearID is the unique identifier for a WearAndTear record.
type WearAndTearID string

// SinkingFundForItemID is the unique identifier for a SinkingFundForItem.
type SinkingFundForItemID string

// RepairID is the unique identifier for a Repair record.
type RepairID string

// CirculatingCycleID is the unique identifier for a CirculatingCycle.
type CirculatingCycleID string

// SinkingFundReinvestmentID is the unique identifier for a SinkingFundReinvestment.
type SinkingFundReinvestmentID string

// CapitalKind distinguishes fixed from circulating productive-capital components.
// The distinction is internal to productive capital only.
type CapitalKind string

const (
	KindFixed       CapitalKind = "fixed"
	KindCirculating CapitalKind = "circulating"
)

// RoleInProcess is the specific role a component plays in the labour-process.
// The role uniquely determines CapitalKind via KindForRole.
type RoleInProcess string

const (
	// RoleInstrumentOfLabour → KindFixed (machines, buildings, beasts of toil)
	RoleInstrumentOfLabour RoleInProcess = "instrument_of_labour"
	// RoleRawMaterial → KindCirculating (enters product bodily)
	RoleRawMaterial RoleInProcess = "raw_material"
	// RoleAuxiliaryMaterial → KindCirculating (Marx rejects Ramsay's misclassification
	// of auxiliaries as fixed even when they don't pass bodily into the product)
	RoleAuxiliaryMaterial RoleInProcess = "auxiliary_material"
	// RoleLabourPower → KindCirculating (variable capital always circulates)
	RoleLabourPower RoleInProcess = "labour_power"
)

// KindForRole returns the CapitalKind that the given role entails.
// RoleInstrumentOfLabour → KindFixed; all others → KindCirculating.
func KindForRole(r RoleInProcess) CapitalKind {
	if r == RoleInstrumentOfLabour {
		return KindFixed
	}
	return KindCirculating
}

// WearCause classifies why a fixed-capital item loses value.
type WearCause string

const (
	WearUse               WearCause = "use"                // proportional to use-time
	WearAtmospheric       WearCause = "atmospheric"        // rot, weather, corrosion
	WearMoralDepreciation WearCause = "moral_depreciation" // cheaper equivalents on market
)

// WearModel selects the per-period wear formula.
type WearModel string

const (
	WearLinear    WearModel = "linear"    // pence = purchased × period / life (default)
	WearQuadratic WearModel = "quadratic" // rail-with-speed: ∝ (speed/baseline)²
)

// RepairKind distinguishes running maintenance from capitalised large repairs.
type RepairKind string

const (
	RepairRunning RepairKind = "running" // circulating-capital cost; no life extension
	RepairLarge   RepairKind = "large"   // capitalised; extends ServiceLifeMinutes
)

// ReinvestmentKind classifies how sinking-fund money is reinvested in new instruments.
type ReinvestmentKind string

const (
	ReinvestmentIntensive ReinvestmentKind = "intensive" // more productive means of production
	ReinvestmentExtensive ReinvestmentKind = "extensive" // more means of production
)

// Sentinel errors.
var (
	ErrRoleKindMismatch               = errors.New("role and kind are inconsistent: use KindForRole to derive kind from role")
	ErrRunningRepairMustNotExtendLife  = errors.New("running repair must not extend service life (service_life_extension_minutes must be 0)")
	ErrLargeRepairMustExtendLife      = errors.New("large repair must extend service life (service_life_extension_minutes must be > 0)")
	ErrZeroServiceLife                = errors.New("service_life_minutes must be positive")
	ErrInvalidRole                    = errors.New("role is not one of the valid RoleInProcess constants")
)

// CapitalComponent is one line-item of productive capital with an identified role
// and a derived kind. The kind is always derived from the role via KindForRole.
type CapitalComponent struct {
	ID                  CapitalComponentID `json:"id"`
	IndustrialCapitalID string             `json:"industrial_capital_id"`
	Kind                CapitalKind        `json:"kind"`
	Role                RoleInProcess      `json:"role"`
	Pence               Pence              `json:"pence"`
	EnteredAt           time.Time          `json:"entered_at"`
}

// Validate confirms that Kind == KindForRole(Role) and that Role is a known constant.
func (c CapitalComponent) Validate() error {
	switch c.Role {
	case RoleInstrumentOfLabour, RoleRawMaterial, RoleAuxiliaryMaterial, RoleLabourPower:
		// valid
	default:
		return ErrInvalidRole
	}
	if c.Kind != KindForRole(c.Role) {
		return ErrRoleKindMismatch
	}
	return nil
}

// FixedCapitalItem is the bodily form of a fixed-capital component: a machine,
// building, beast of toil, or any instrument of labour. Its value transfers to
// the product piecemeal through WearAndTear over its ServiceLifeMinutes.
type FixedCapitalItem struct {
	ID                 FixedCapitalItemID `json:"id"`
	CapitalComponentID CapitalComponentID `json:"capital_component_id"`
	Description        string             `json:"description"`
	PencePurchased     Pence              `json:"pence_purchased"`
	ServiceLifeMinutes int64              `json:"service_life_minutes"`
	WearModel          WearModel          `json:"wear_model"`
	PurchasedAt        time.Time          `json:"purchased_at"`
	RetiredAt          *time.Time         `json:"retired_at,omitempty"`
}

// Validate returns an error if the item's service life is non-positive.
func (f FixedCapitalItem) Validate() error {
	if f.ServiceLifeMinutes <= 0 {
		return ErrZeroServiceLife
	}
	return nil
}

// ComputeLinearWear calculates the WearAndTear value transferred to the product
// during periodMinutes under the linear (default) model:
//
//	pence = PencePurchased × periodMinutes / ServiceLifeMinutes
//
// Integer division; the sum over the full lifetime equals PencePurchased.
func ComputeLinearWear(item FixedCapitalItem, periodMinutes int64, cause WearCause) WearAndTear {
	var p Pence
	if item.ServiceLifeMinutes > 0 {
		p = Pence(int64(item.PencePurchased) * periodMinutes / item.ServiceLifeMinutes)
	}
	return WearAndTear{
		FixedCapitalItemID:   item.ID,
		Pence:                p,
		PeriodCoveredMinutes: periodMinutes,
		CalculatedAt:         time.Now().UTC(),
		Cause:                cause,
	}
}

// FixedCapitalSubcomponent is a homogeneous constituent of a FixedCapitalItem
// that wears at its own rate and may be replaced independently (e.g., railway
// rails vs earthworks vs sleepers in a British railway).
type FixedCapitalSubcomponent struct {
	ID                 FixedCapitalSubcomponentID `json:"id"`
	FixedCapitalItemID FixedCapitalItemID         `json:"fixed_capital_item_id"`
	Name               string                     `json:"name"`
	Pence              Pence                      `json:"pence"`
	ServiceLifeMinutes int64                      `json:"service_life_minutes"`
	WearModel          WearModel                  `json:"wear_model"`
	LastReplacedAt     time.Time                  `json:"last_replaced_at"`
}

// WearAndTear records the value transferred from a FixedCapitalItem to the
// product in one period. Many WearAndTear rows exist over the item's life.
type WearAndTear struct {
	ID                   WearAndTearID      `json:"id"`
	FixedCapitalItemID   FixedCapitalItemID `json:"fixed_capital_item_id"`
	Pence                Pence              `json:"pence"`
	PeriodCoveredMinutes int64              `json:"period_covered_minutes"`
	CalculatedAt         time.Time          `json:"calculated_at"`
	Cause                WearCause          `json:"cause"`
}

// SinkingFundForItem accumulates the money-form of value shed by a fixed-capital
// item as WearAndTear. When AccumulatedPence reaches TargetPence (== PencePurchased)
// the item may be replaced.
type SinkingFundForItem struct {
	ID                 SinkingFundForItemID `json:"id"`
	FixedCapitalItemID FixedCapitalItemID   `json:"fixed_capital_item_id"`
	AccumulatedPence   Pence                `json:"accumulated_pence"`
	TargetPence        Pence                `json:"target_pence"` // == FixedCapitalItem.PencePurchased
}

// Repair records maintenance expenditure on a fixed-capital item. Running repairs
// are part of circulating-capital costs and must not extend service life; large
// repairs are capitalised and must extend service life.
type Repair struct {
	ID                          RepairID           `json:"id"`
	FixedCapitalItemID          FixedCapitalItemID `json:"fixed_capital_item_id"`
	Kind                        RepairKind         `json:"kind"`
	Pence                       Pence              `json:"pence"`
	OccurredAt                  time.Time          `json:"occurred_at"`
	ServiceLifeExtensionMinutes int64              `json:"service_life_extension_minutes"`
}

// Validate enforces the running/large-repair life-extension invariants.
func (r Repair) Validate() error {
	switch r.Kind {
	case RepairRunning:
		if r.ServiceLifeExtensionMinutes != 0 {
			return ErrRunningRepairMustNotExtendLife
		}
	case RepairLarge:
		if r.ServiceLifeExtensionMinutes <= 0 {
			return ErrLargeRepairMustExtendLife
		}
	}
	return nil
}

// CirculatingCycle records one full replacement cycle of a circulating-capital
// component. Many cycles occur within one FixedCapitalItem.ServiceLifeMinutes.
type CirculatingCycle struct {
	ID                 CirculatingCycleID `json:"id"`
	CapitalComponentID CapitalComponentID `json:"capital_component_id"`
	StartedAt          time.Time          `json:"started_at"`
	EndedAt            time.Time          `json:"ended_at"`
	Pence              Pence              `json:"pence"`
}

// SinkingFundReinvestment records a reinvestment of sinking-fund money into
// improved or additional instruments of labour. This is reproduction on an
// extended scale funded by released fixed-capital value — not from accumulated
// surplus-value (Ch. 2's LatentMoneyCapital).
type SinkingFundReinvestment struct {
	ID                 SinkingFundReinvestmentID `json:"id"`
	FixedCapitalItemID FixedCapitalItemID         `json:"fixed_capital_item_id"`
	Pence              Pence                      `json:"pence"`
	Kind               ReinvestmentKind           `json:"kind"`
	OccurredAt         time.Time                  `json:"occurred_at"`
}

// ID generators — 96-bit hex via crypto/rand.

func newID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(b)
}

// NewCapitalComponentID returns a new unique CapitalComponentID.
func NewCapitalComponentID() CapitalComponentID { return CapitalComponentID(newID()) }

// NewFixedCapitalItemID returns a new unique FixedCapitalItemID.
func NewFixedCapitalItemID() FixedCapitalItemID { return FixedCapitalItemID(newID()) }

// NewFixedCapitalSubcomponentID returns a new unique FixedCapitalSubcomponentID.
func NewFixedCapitalSubcomponentID() FixedCapitalSubcomponentID {
	return FixedCapitalSubcomponentID(newID())
}

// NewWearAndTearID returns a new unique WearAndTearID.
func NewWearAndTearID() WearAndTearID { return WearAndTearID(newID()) }

// NewSinkingFundForItemID returns a new unique SinkingFundForItemID.
func NewSinkingFundForItemID() SinkingFundForItemID { return SinkingFundForItemID(newID()) }

// NewRepairID returns a new unique RepairID.
func NewRepairID() RepairID { return RepairID(newID()) }

// NewCirculatingCycleID returns a new unique CirculatingCycleID.
func NewCirculatingCycleID() CirculatingCycleID { return CirculatingCycleID(newID()) }

// NewSinkingFundReinvestmentID returns a new unique SinkingFundReinvestmentID.
func NewSinkingFundReinvestmentID() SinkingFundReinvestmentID {
	return SinkingFundReinvestmentID(newID())
}
