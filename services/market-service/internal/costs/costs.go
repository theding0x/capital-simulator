// Package costs models the costs of circulation from Capital Vol. II, Ch. 6.
// Circulation costs are sub-divided into:
//
//   - Five "faux frais" (spurious costs): purchase/sale time, book-keeping,
//     money-as-faux-frais, supply formation, and commodity-supply holding.
//     These are deductions from surplus-value and do NOT add to commodity-value.
//
//   - Transportation: the one productive cost of circulation that DOES add
//     value to the transported commodity (labour + wear of transport-MP +
//     surplus appropriated by the transport-capitalist).
package costs

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Pence is the canonical monetary unit (1 GBP = 240 pence).
type Pence int64

// Sentinel errors.
var (
	ErrInvalidCirculationCost = errors.New("costs: invalid circulation cost")
	ErrInvalidTransportLeg    = errors.New("costs: invalid transport leg")
	ErrInvalidCommoditySupply = errors.New("costs: invalid commodity supply")
	ErrInvalidTransportTariff = errors.New("costs: invalid transport tariff")
)

// --- ID types ---------------------------------------------------------------

// CirculationCostID is a 96-bit hex identifier for a CirculationCost record.
type CirculationCostID string

// NewCirculationCostID generates a 96-bit hex identifier via crypto/rand.
func NewCirculationCostID() CirculationCostID { return CirculationCostID(newHexID()) }

// CirculationAgentID is a 96-bit hex identifier for a CirculationAgent record.
type CirculationAgentID string

// NewCirculationAgentID generates a 96-bit hex identifier via crypto/rand.
func NewCirculationAgentID() CirculationAgentID { return CirculationAgentID(newHexID()) }

// CommoditySupplyID is a 96-bit hex identifier for a CommoditySupply record.
type CommoditySupplyID string

// NewCommoditySupplyID generates a 96-bit hex identifier via crypto/rand.
func NewCommoditySupplyID() CommoditySupplyID { return CommoditySupplyID(newHexID()) }

// StorageCostID is a 96-bit hex identifier for a StorageCost record.
type StorageCostID string

// NewStorageCostID generates a 96-bit hex identifier via crypto/rand.
func NewStorageCostID() StorageCostID { return StorageCostID(newHexID()) }

// TransportLegID is a 96-bit hex identifier for a TransportLeg record.
type TransportLegID string

// NewTransportLegID generates a 96-bit hex identifier via crypto/rand.
func NewTransportLegID() TransportLegID { return TransportLegID(newHexID()) }

// TransportTariffID is a 96-bit hex identifier for a TransportTariff record.
type TransportTariffID string

// NewTransportTariffID generates a 96-bit hex identifier via crypto/rand.
func NewTransportTariffID() TransportTariffID { return TransportTariffID(newHexID()) }

// MoneyAsFauxFraisID is a 96-bit hex identifier for a MoneyAsFauxFrais record.
type MoneyAsFauxFraisID string

// NewMoneyAsFauxFraisID generates a 96-bit hex identifier via crypto/rand.
func NewMoneyAsFauxFraisID() MoneyAsFauxFraisID { return MoneyAsFauxFraisID(newHexID()) }

// --- Cross-service reference types ------------------------------------------

// IndustrialCapitalID is an opaque reference to an IndustrialCapital record
// managed by simulation-engine (Vol. II Ch. 4).
type IndustrialCapitalID string

// AgentID is an opaque reference to an agent record in agent-service.
type AgentID string

// CommodityID is an opaque reference to a Commodity in commodity-service.
type CommodityID string

// --- Circulation cost kind and nature ---------------------------------------

// CirculationCostKind names one of Marx's six categories of circulation cost.
// Only CostTransportation is value-adding; the other five are faux frais.
type CirculationCostKind string

const (
	// CostPurchaseSaleTime is the labour-time spent buying and selling
	// (excluding transportation). It is a faux frais. [§I(a)]
	CostPurchaseSaleTime CirculationCostKind = "purchase_sale_time"
	// CostBookKeeping covers book-keeping labour and stationery. [§I(b)]
	CostBookKeeping CirculationCostKind = "book_keeping"
	// CostMoneyAsFauxFrais covers the gold/silver mass tied up in circulation
	// and its annual replacement burden. [§I(c)]
	CostMoneyAsFauxFraisKind CirculationCostKind = "money_as_faux_frais"
	// CostSupplyFormation covers the costs of forming and maintaining a
	// commodity supply (storage, preservation). [§II]
	CostSupplyFormation CirculationCostKind = "supply_formation"
	// CostCommoditySupply covers the capital tied up as unsold stock. [§II]
	CostCommoditySupply CirculationCostKind = "commodity_supply"
	// CostTransportation is the ONE productive cost: it adds value. [§III]
	CostTransportation CirculationCostKind = "transportation"
)

// NatureOf returns the CostNature that always accompanies the given kind.
// Only transportation is value-adding; all others are value-preserving faux frais.
func NatureOf(kind CirculationCostKind) CostNature {
	if kind == CostTransportation {
		return NatureValueAdding
	}
	return NatureValuePreserving
}

// CostNature records whether a cost adds to commodity-value or merely
// deducts from surplus-value without adding to social wealth.
type CostNature string

const (
	// NatureValuePreserving — faux frais; deducted from surplus-value, not
	// added to commodity-value.
	NatureValuePreserving CostNature = "value_preserving"
	// NatureValueAdding — transportation labour + MP wear + transport-surplus;
	// added to the transported commodity's value.
	NatureValueAdding CostNature = "value_adding"
)

// --- CirculationCost --------------------------------------------------------

// CirculationCost records one incurred cost of circulation for an industrial
// capital. Its Nature is derived deterministically from its Kind.
type CirculationCost struct {
	ID                  CirculationCostID   `json:"id"`
	IndustrialCapitalID IndustrialCapitalID `json:"industrial_capital_id"`
	Kind                CirculationCostKind `json:"kind"`
	Nature              CostNature          `json:"nature"`
	Pence               Pence               `json:"pence"`
	IncurredAt          time.Time           `json:"incurred_at"`
}

// NewCirculationCost constructs a CirculationCost with Nature derived from Kind.
func NewCirculationCost(icID IndustrialCapitalID, kind CirculationCostKind, pence Pence, at time.Time) (CirculationCost, error) {
	if icID == "" {
		return CirculationCost{}, ErrInvalidCirculationCost
	}
	if pence < 0 {
		return CirculationCost{}, ErrInvalidCirculationCost
	}
	return CirculationCost{
		ID:                  NewCirculationCostID(),
		IndustrialCapitalID: icID,
		Kind:                kind,
		Nature:              NatureOf(kind),
		Pence:               pence,
		IncurredAt:          at,
	}, nil
}

// --- CirculationAgent -------------------------------------------------------

// CirculationAgentRole names the function a circulation labourer performs.
type CirculationAgentRole string

const (
	RoleBuyer       CirculationAgentRole = "buyer"
	RoleSeller      CirculationAgentRole = "seller"
	RoleBookKeeper  CirculationAgentRole = "book_keeper"
	RoleCashier     CirculationAgentRole = "cashier"
	RoleStorekeeper CirculationAgentRole = "storekeeper"
)

// CirculationAgent is a wage-labourer employed in buying, selling, or
// book-keeping functions. Neither their necessary nor their surplus labour
// adds value to commodities — both are deductions from social wealth.
//
// [§I(a): "in the first place, looking at it from the standpoint of society,
// labour-power is used up now as before for ten hours in a mere function of
// circulation. It cannot be used for anything else."]
type CirculationAgent struct {
	ID                     CirculationAgentID   `json:"id"`
	IndustrialCapitalID    IndustrialCapitalID  `json:"industrial_capital_id"`
	AgentID                AgentID              `json:"agent_id,omitempty"`
	Role                   CirculationAgentRole `json:"role"`
	WageRatePence          Pence                `json:"wage_rate_pence"`
	// LabourMinutesNecessary is the portion of the working day for which the
	// agent receives their wage equivalent. It creates no value.
	LabourMinutesNecessary int64 `json:"labour_minutes_necessary"`
	// LabourMinutesSurplus is the unpaid portion of the working day. It reduces
	// the capitalist's faux-frais burden but creates no new social value.
	LabourMinutesSurplus int64 `json:"labour_minutes_surplus"`
}

// WageEffectivePence returns the actual wage cost to the employing capitalist
// after the surplus portion of the agent's labour compresses the per-minute
// effective wage.
func (a CirculationAgent) WageEffectivePence() Pence {
	total := a.LabourMinutesNecessary + a.LabourMinutesSurplus
	if total == 0 {
		return 0
	}
	return Pence(int64(a.WageRatePence) * a.LabourMinutesNecessary / total)
}

// ValueContribution is always zero: neither necessary nor surplus labour of a
// circulation agent adds to commodity-value.
func (CirculationAgent) ValueContribution() Pence { return 0 }

// --- MoneyAsFauxFrais -------------------------------------------------------

// MoneyAsFauxFrais records the gold/silver mass tied up in circulation.
// This may be system-wide (IndustrialCapitalID empty) or per-capital.
// [§I(c): "Gold and silver as money-commodities mean circulation costs to
// society which arise solely out of the social form of production."]
type MoneyAsFauxFrais struct {
	ID                     MoneyAsFauxFraisID  `json:"id"`
	IndustrialCapitalID    IndustrialCapitalID `json:"industrial_capital_id,omitempty"`
	Pence                  Pence               `json:"pence"`
	AnnualReplacementPence Pence               `json:"annual_replacement_pence"`
}

// --- CommoditySupply --------------------------------------------------------

// CommoditySupply records capital held in commodity-form as unsold stock.
// When IsVoluntary=true, this is a normal buffer held for steady sale.
// When IsVoluntary=false, the supply exists because of sales resistance
// (stagnation) — both cases incur identical storage costs, but the latter
// is a pure loss.
//
// [§II(a): "Whether the formation of a supply is voluntary or involuntary…
// cannot effect the matter essentially."]
type CommoditySupply struct {
	ID                  CommoditySupplyID   `json:"id"`
	IndustrialCapitalID IndustrialCapitalID `json:"industrial_capital_id"`
	CommodityID         CommodityID         `json:"commodity_id,omitempty"`
	Pence               Pence               `json:"pence"`
	HeldSince           time.Time           `json:"held_since"`
	IsVoluntary         bool                `json:"is_voluntary"`
}

// SupplyDeteriorationKind names one of Marx's three modes of supply
// deterioration cost. [§II(a)]
type SupplyDeteriorationKind string

const (
	// DeteriorationQuantitative: mass loss (e.g. flour dust, evaporation).
	DeteriorationQuantitative SupplyDeteriorationKind = "quantitative"
	// DeteriorationQualitative: spoilage, quality degradation.
	DeteriorationQualitative SupplyDeteriorationKind = "qualitative"
	// DeteriorationLabour: preservation outlay.
	DeteriorationLabour SupplyDeteriorationKind = "labour"
)

// --- StorageCost ------------------------------------------------------------

// StorageCost records the three Marx-ian cost buckets of holding a
// CommoditySupply. All three are faux frais.
type StorageCost struct {
	ID                      StorageCostID     `json:"id"`
	CommoditySupplyID       CommoditySupplyID `json:"commodity_supply_id"`
	BuildingPence           Pence             `json:"building_pence"`
	LabourPence             Pence             `json:"labour_pence"`
	PreservationLabourPence Pence             `json:"preservation_labour_pence"`
}

// Total returns the aggregate storage cost.
func (sc StorageCost) Total() Pence {
	return sc.BuildingPence + sc.LabourPence + sc.PreservationLabourPence
}

// --- TransportTariff --------------------------------------------------------

// TransportTariff defines the per-ton-mile rate model for a commodity.
// Fragility and breakage-risk multipliers are in basis points (10000 bp = 1×).
// A BreakageRiskMultiplierBasisPoints of 30000 means 3× the base rate.
// [§III Birmingham glass: "the rate for carriage is the same as it was
// formerly, and higher than it was previously… the rate to cover risk of
// breakage is three times that amount."]
type TransportTariff struct {
	ID                                TransportTariffID `json:"id"`
	CommodityID                       CommodityID       `json:"commodity_id"`
	BasePencePerTonMile               Pence             `json:"base_pence_per_ton_mile"`
	FragilityMultiplierBasisPoints    int64             `json:"fragility_multiplier_basis_points"`
	BreakageRiskMultiplierBasisPoints int64             `json:"breakage_risk_multiplier_basis_points"`
}

// EffectiveRatePencePerTonMile applies the highest applicable multiplier to
// the base rate.
func (t TransportTariff) EffectiveRatePencePerTonMile() Pence {
	mul := t.FragilityMultiplierBasisPoints
	if t.BreakageRiskMultiplierBasisPoints > mul {
		mul = t.BreakageRiskMultiplierBasisPoints
	}
	if mul <= 0 {
		mul = 10000 // 1× — no multiplier
	}
	return Pence(int64(t.BasePencePerTonMile) * mul / 10000)
}

// --- TransportLeg -----------------------------------------------------------

// GramsPerTonne is the number of grams in one metric tonne.
const GramsPerTonne = 1_000_000

// MetersPerMile is the number of metres in one statute mile (approximate).
const MetersPerMile = 1609

// TransportLeg records one real spatial movement of a commodity from origin
// to destination. Unlike all other circulation costs, transport adds value.
//
// [§III: "The productive capital invested in this industry imparts value to
// the transported products, partly by transferring value from the means of
// transportation, partly by adding value through the labour performed in
// transport."]
type TransportLeg struct {
	ID                     TransportLegID `json:"id"`
	CommodityID            CommodityID    `json:"commodity_id"`
	Quantity               int64          `json:"quantity"`
	OriginID               string         `json:"origin_id"`
	DestinationID          string         `json:"destination_id"`
	DistanceMeters         int64          `json:"distance_meters"`
	WeightGrams            int64          `json:"weight_grams"`
	VolumeCubicMillimeters int64          `json:"volume_cubic_millimeters"`
	// LabourCostPence is the labour newly added during transport.
	LabourCostPence Pence `json:"labour_cost_pence"`
	// MeansOfTransportPence is the wear transferred from transport-MP.
	MeansOfTransportPence Pence `json:"means_of_transport_pence"`
	// SurplusPence is the surplus-value appropriated by the transport-capitalist.
	SurplusPence Pence `json:"surplus_pence"`
}

// AddedValue returns the total value added to the commodity by this transport
// leg. This is the only circulation cost that enters the commodity's value.
//
// Invariants (per §III):
//   - Monotonic in DistanceMeters: doubling distance doubles AddedValue().
//   - Antitonic in BasePencePerTonMile: halving productivity doubles AddedValue().
func (tl TransportLeg) AddedValue() Pence {
	return tl.LabourCostPence + tl.MeansOfTransportPence + tl.SurplusPence
}

// --- AggregateResult --------------------------------------------------------

// AggregateResult is the response body for GET /v1/circulation-costs/aggregate.
type AggregateResult struct {
	IndustrialCapitalID IndustrialCapitalID `json:"industrial_capital_id"`
	Period              string              `json:"period"`
	TotalFauxFraisPence Pence               `json:"total_faux_frais_pence"`
	TotalTransportPence Pence               `json:"total_transport_pence"`
	Items               []CirculationCost   `json:"items"`
}

// SystemFauxFraisResult is the response body for GET /v1/circulation-costs/system-faux-frais.
type SystemFauxFraisResult struct {
	Period                      string             `json:"period"`
	TotalReservePence           Pence              `json:"total_reserve_pence"`
	TotalAnnualReplacementPence Pence              `json:"total_annual_replacement_pence"`
	Items                       []MoneyAsFauxFrais `json:"items"`
}

// --- helpers ----------------------------------------------------------------

func newHexID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(b)
}
