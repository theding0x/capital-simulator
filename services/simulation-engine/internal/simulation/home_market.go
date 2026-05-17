package simulation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Ch. 30 — Reaction of the Agricultural Revolution on Industry.
//
// The displacement of agricultural labourers and the destruction of rural
// domestic industry freed both labour and purchasing power, creating the home
// market that industrial capital required. Formerly self-produced subsistence
// goods became commodities; former peasants became wage workers who must now
// buy what they once made.

// DomesticIndustryID is a 96-bit hex identifier for a persisted domestic industry.
type DomesticIndustryID string

// IsZero reports whether the ID is unset.
func (id DomesticIndustryID) IsZero() bool { return id == "" }

// NewDomesticIndustryID generates a 96-bit hex identifier via crypto/rand.
func NewDomesticIndustryID() DomesticIndustryID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return DomesticIndustryID(hex.EncodeToString(b))
}

var (
	// ErrDomesticIndustryNameRequired — the industry must be named.
	ErrDomesticIndustryNameRequired = errors.New("simulation: domestic industry name is required")

	// ErrHouseholdsNonPositive — there must be at least one household engaged.
	ErrHouseholdsNonPositive = errors.New("simulation: households_engaged must be positive")

	// ErrAnnualOutputNonPositive — the annual output value must be positive.
	ErrAnnualOutputNonPositive = errors.New("simulation: annual_output_pence must be positive")
)

// DomesticIndustry records a pre-capitalist household production unit —
// spinning, weaving, flax work, or other rural putting-out labour.
// When proletarianisation sweeps a region these households are dissolved;
// their output becomes a commodity and their members join the wage-labour market.
type DomesticIndustry struct {
	ID                DomesticIndustryID `json:"id"`
	HistoricalStageID HistoricalStageID  `json:"historical_stage_id"`
	Name              string             `json:"name"`
	HouseholdsEngaged int64              `json:"households_engaged"`
	AnnualOutputPence Pence              `json:"annual_output_pence"`
	CreatedAt         time.Time          `json:"created_at"`
}

// Validate enforces invariants on a DomesticIndustry record.
func (d DomesticIndustry) Validate() error {
	if d.Name == "" {
		return ErrDomesticIndustryNameRequired
	}
	if d.HouseholdsEngaged <= 0 {
		return ErrHouseholdsNonPositive
	}
	if d.AnnualOutputPence <= 0 {
		return ErrAnnualOutputNonPositive
	}
	return nil
}

// ManufactureReunie is Mirabeau's "grand manufactory" — hundreds of workers
// concentrated under a single direction. Concentration disciplines labour and
// captures the full productive advantages of co-operation.
type ManufactureReunie struct {
	Workers     int64 `json:"workers"`
	OutputPence Pence `json:"output_pence"`
}

// ManufactureSeparee is Mirabeau's preferred "discrete workshop" — scattered
// independent producers, each working on their own account. Politically stable
// for petty property but economically inferior to concentrated manufacture.
type ManufactureSeparee struct {
	IndependentProducers int64 `json:"independent_producers"`
	OutputPence          Pence `json:"output_pence"`
}

// HomeMarketFormation records the quantitative result of proletarianisation:
// domestic producers become wage workers, and their formerly self-consumed
// output enters the commodity market.
type HomeMarketFormation struct {
	Period                    string `json:"period"`
	LabourersProletarianised  int64  `json:"labourers_proletarianised"`
	DomesticIndustryDestroyed bool   `json:"domestic_industry_destroyed"`
	NewWageLabourers          int64  `json:"new_wage_labourers"`
	MarketSizePence           Pence  `json:"market_size_pence"`
}

// HomeMarketSize is the aggregate home market created by proletarianisation
// within a historical stage.
type HomeMarketSize struct {
	CommoditisedOutputPence Pence `json:"commoditised_output_pence"`
	WageLabourers           int64 `json:"wage_labourers"`
}

// ComputeMarketFormation converts domestic producers into wage workers and
// subsistence goods into commodities. Every expropriated household reappears
// as a new wage worker; their annual output becomes the commoditised market.
// DomesticIndustryDestroyed is true when expropriated >= engaged households.
//
// Invariant: NewWageLabourers == LabourersProletarianised (conservation law).
func ComputeMarketFormation(di DomesticIndustry, expropriated int64) HomeMarketFormation {
	if expropriated < 0 {
		expropriated = 0
	}
	return HomeMarketFormation{
		Period:                    di.Name,
		LabourersProletarianised:  expropriated,
		DomesticIndustryDestroyed: expropriated >= di.HouseholdsEngaged,
		NewWageLabourers:          expropriated,
		MarketSizePence:           di.AnnualOutputPence,
	}
}

// AggregateHomeMarket sums the home market contribution of all domestic
// industries in a historical stage. Each household becomes a potential wage
// worker; each unit of annual output becomes a commoditised good.
func AggregateHomeMarket(industries []DomesticIndustry) HomeMarketSize {
	var total Pence
	var workers int64
	for _, d := range industries {
		total += d.AnnualOutputPence
		workers += d.HouseholdsEngaged
	}
	return HomeMarketSize{
		CommoditisedOutputPence: total,
		WageLabourers:           workers,
	}
}
