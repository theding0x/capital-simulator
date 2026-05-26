// Vol. II Ch. 9 — The Aggregate Turnover of Advanced Capital. Cycles of Turnover.
//
// Marx establishes that the turnover of a mixed (fixed + circulating) capital can
// only be measured homogeneously in the money-form (M…M'). The aggregate number n
// equals 10,000 × AnnualTurnedOverPence / AdvancedTotalPence (basis-points).
// For the canonical example: £80k fixed (10yr) + £20k circulating (5×/yr) →
// n = 10800 bp (1 + 2/25). The "mean term" is AdvancedTotalPence × T / AnnualTurnedOverPence.
package turnover

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	comp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/composition"
)

// ErrLensNotHomogeneous is returned when an AggregateTurnover is created with a
// lens other than LensMoney. Marx: "The qualitative identity does not come about
// if we take as our starting-point P...P. However the form M...M' undoubtedly
// yields this identity of turnover." (§2)
var ErrLensNotHomogeneous = errors.New("aggregate turnover requires money-form (LensMoney) for homogeneous comparison")

// --- ID types ----------------------------------------------------------------

// AggregateTurnoverID identifies an AggregateTurnover.
type AggregateTurnoverID string

// NewAggregateTurnoverID returns a 96-bit random hex identifier.
func NewAggregateTurnoverID() AggregateTurnoverID {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return AggregateTurnoverID(hex.EncodeToString(b))
}

// ComponentTurnoverContributionID identifies a per-component contribution row.
type ComponentTurnoverContributionID string

// NewComponentTurnoverContributionID returns a 96-bit random hex identifier.
func NewComponentTurnoverContributionID() ComponentTurnoverContributionID {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return ComponentTurnoverContributionID(hex.EncodeToString(b))
}

// LifetimeCycleID identifies a LifetimeCycle.
type LifetimeCycleID string

// NewLifetimeCycleID returns a 96-bit random hex identifier.
func NewLifetimeCycleID() LifetimeCycleID {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return LifetimeCycleID(hex.EncodeToString(b))
}

// CrisisCyclePhaseRecordID identifies a persisted CrisisCyclePhaseRecord.
type CrisisCyclePhaseRecordID string

// NewCrisisCyclePhaseRecordID returns a 96-bit random hex identifier.
func NewCrisisCyclePhaseRecordID() CrisisCyclePhaseRecordID {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return CrisisCyclePhaseRecordID(hex.EncodeToString(b))
}

// --- Enum types --------------------------------------------------------------

// CrisisCyclePhase is one of the four phases of the long-run crisis cycle that
// Marx associates with fixed-capital reproduction. Marx explicitly names these
// four phases in §4: "depression, medium activity, precipitancy, crisis."
type CrisisCyclePhase string

const (
	PhaseDepression     CrisisCyclePhase = "depression"
	PhaseMediumActivity CrisisCyclePhase = "medium_activity"
	PhasePrecipitancy   CrisisCyclePhase = "precipitancy"
	PhaseCrisis         CrisisCyclePhase = "crisis"
)

// TurnoverDifferenceKind tags a ComponentTurnoverContribution to distinguish
// real differences in turnover (arising from the nature of the capital) from
// apparent differences (driven by credit or payment-term lengths). Both are
// aggregated identically; the tag is for display and Vol. III follow-ups. (§6)
type TurnoverDifferenceKind string

const (
	// DifferenceReal — difference arises from the intrinsic nature of the capital.
	DifferenceReal TurnoverDifferenceKind = "real"
	// DifferenceApparent — difference is an artefact of credit / payment-term lengths.
	DifferenceApparent TurnoverDifferenceKind = "apparent"
)

// --- Aggregate types ---------------------------------------------------------

// AggregateTurnover is the whole-capital (fixed + circulating) turnover for a
// given IndustrialCapital. Lens is forced to LensMoney: only the money-form
// circuit yields the homogeneous comparison required for aggregation. (§2)
//
// AggregateNumberBasisPoints = 10000 × AnnualTurnedOverPence / AdvancedTotalPence.
// Marx's central example: £108,000 / £100,000 = 1.08 → 10800 bp (= 1 + 2/25).
type AggregateTurnover struct {
	ID                         AggregateTurnoverID             `json:"id"`
	IndustrialCapitalID        string                          `json:"industrial_capital_id"`
	Lens                       CircuitLens                     `json:"lens"`
	AdvancedTotalPence         Pence                           `json:"advanced_total_pence"`
	AnnualTurnedOverPence      Pence                           `json:"annual_turned_over_pence"`
	AggregateNumberBasisPoints int64                           `json:"aggregate_number_basis_points"`
	Contributions              []ComponentTurnoverContribution `json:"contributions"`
	LifetimeCycle              *LifetimeCycle                  `json:"lifetime_cycle,omitempty"`
	LatestCrisisPhase          *CrisisCyclePhaseRecord         `json:"latest_crisis_phase,omitempty"`
}

// MeanTerm returns the aggregate mean term of turnover in minutes without
// surplus-value (the Marx-corrected / profit-discounted form).
//
//	MeanTerm = AdvancedTotalPence × YearReferenceMinutes / AnnualTurnedOverPence
//
// Marx's §5 example: $50,000 advanced / $33,750 annual × 525,600 min ≈ 778,667 min ≈ 18 months.
func (a AggregateTurnover) MeanTerm() int64 {
	if a.AnnualTurnedOverPence == 0 {
		return 0
	}
	return int64(a.AdvancedTotalPence) * YearReferenceMinutes / int64(a.AnnualTurnedOverPence)
}

// MeanTermWithProfit returns the Scrope-style profit-inclusive mean term in minutes
// (shorter than MeanTerm because the surplus-value inflates the denominator).
// Marx §5: 18 months without profit, 16 months with profit for the Scrope $50k example.
// surplusValuePence is the annual surplus-value produced by the capital.
func (a AggregateTurnover) MeanTermWithProfit(surplusValuePence Pence) int64 {
	denom := int64(a.AnnualTurnedOverPence) + int64(surplusValuePence)
	if denom == 0 {
		return 0
	}
	return int64(a.AdvancedTotalPence) * YearReferenceMinutes / denom
}

// Validate checks the money-form homogeneity invariant and basic field sanity.
func (a AggregateTurnover) Validate() error {
	if a.Lens != LensMoney {
		return ErrLensNotHomogeneous
	}
	if a.AdvancedTotalPence <= 0 {
		return errors.New("advanced_total_pence must be positive")
	}
	return nil
}

// ComponentTurnoverContribution records how one component of a mixed capital
// contributes to the aggregate annual turnover.
//
//	AnnualReturnPence = AdvancedPence × TurnoverNumberBasisPoints / 10000
//
// Fixed capital at 1/10th per year: bp = 1000, annual_return = advanced / 10.
// Circulating capital at 5×/year: bp = 50000, annual_return = advanced × 5.
type ComponentTurnoverContribution struct {
	ID                        ComponentTurnoverContributionID `json:"id"`
	AggregateTurnoverID       AggregateTurnoverID             `json:"aggregate_turnover_id"`
	CapitalComponentID        comp.CapitalComponentID         `json:"capital_component_id"`
	Kind                      comp.CapitalKind                `json:"kind"`
	AdvancedPence             Pence                           `json:"advanced_pence"`
	TurnoverNumberBasisPoints int64                           `json:"turnover_number_basis_points"`
	AnnualReturnPence         Pence                           `json:"annual_return_pence"`
	DifferenceKind            TurnoverDifferenceKind          `json:"difference_kind"`
}

// ComputeAnnualReturn derives AnnualReturnPence from the other fields.
func (c ComponentTurnoverContribution) ComputeAnnualReturn() Pence {
	return Pence(int64(c.AdvancedPence) * c.TurnoverNumberBasisPoints / 10_000)
}

// LifetimeCycle is the long-run fixed-capital reproduction cycle whose duration
// is set by the longest fixed-capital service-life in the composition.
// "This cycle is determined by the life, hence the reproduction or turnover time,
// of the applied fixed capital." (Marx, §4). For modern industry: ~10 years.
type LifetimeCycle struct {
	ID                  LifetimeCycleID          `json:"id"`
	AggregateTurnoverID AggregateTurnoverID      `json:"aggregate_turnover_id"`
	DurationMinutes     int64                    `json:"duration_minutes"`
	StartedAt           time.Time                `json:"started_at"`
	CompletedCycles     int64                    `json:"completed_cycles"`
	CrisisPhases        []CrisisCyclePhaseRecord `json:"crisis_phases,omitempty"`
}

// CrisisCyclePhaseRecord is an observational annotation of the crisis-cycle phase
// at a specific moment within a LifetimeCycle. Not a deterministic state machine:
// "True, periods in which capital is invested differ greatly and far from coincide
// in time. But a crisis always forms the starting-point of large new investments."
// (Marx, §4)
type CrisisCyclePhaseRecord struct {
	ID              CrisisCyclePhaseRecordID `json:"id"`
	LifetimeCycleID LifetimeCycleID          `json:"lifetime_cycle_id"`
	Phase           CrisisCyclePhase         `json:"phase"`
	ObservedAt      time.Time                `json:"observed_at"`
	Notes           string                   `json:"notes,omitempty"`
}

// --- Pure functions ----------------------------------------------------------

// RecomputeAggregate rebuilds the AggregateTurnover totals from a fresh set of
// contributions. It sets AdvancedTotalPence, AnnualTurnedOverPence, and
// AggregateNumberBasisPoints on a copy of a and returns the updated value.
// The contributions slice is attached to the returned aggregate.
//
// Invariant: AnnualTurnedOverPence == Σ AnnualReturnPence across contributions.
// Invariant: AggregateNumberBasisPoints = 10000 × AnnualTurnedOverPence / AdvancedTotalPence.
func RecomputeAggregate(a AggregateTurnover, contributions []ComponentTurnoverContribution) AggregateTurnover {
	var totalAdvanced, totalAnnual int64
	for _, c := range contributions {
		totalAdvanced += int64(c.AdvancedPence)
		totalAnnual += int64(c.AnnualReturnPence)
	}
	a.AdvancedTotalPence = Pence(totalAdvanced)
	a.AnnualTurnedOverPence = Pence(totalAnnual)
	if totalAdvanced > 0 {
		a.AggregateNumberBasisPoints = 10_000 * totalAnnual / totalAdvanced
	}
	a.Contributions = contributions
	return a
}
