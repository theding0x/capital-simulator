package simulation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"time"
)

// Ch. 32 — The Historical Tendency of Capitalist Accumulation.
//
// Marx names the dialectical telos of Part VIII. Three stages:
//
//   1. Petty property — the labourer owns the means of labour; production
//      is dispersed, small-scale, and excludes co-operation.
//   2. Capitalist private property — many producers are expropriated by
//      few capitalists; means of production are socialised in fact while
//      remaining private in form.
//   3. Socialised property — the second stage, by its own laws of
//      centralisation, prepares its own negation: "the expropriators are
//      expropriated."
//
// The simulation here models the second-to-third transition's recurring
// event: a CentralisationStep in which `FirmsAbsorbed` capitalists are
// swallowed by a smaller number, with `CapitalConcentratedPence` value
// redistributed (not destroyed). Aggregating steps gives an
// AccumulationTrajectory whose endpoint is the asymptotic limit Marx
// describes — `FinalFirms == 1` is the logical terminus, never reached.

const (
	// NegationStagePettyProperty is the dialectical first stage: producers
	// own their means of labour.
	NegationStagePettyProperty = "petty-property"
	// NegationStageCapitalistExpropriation is the second stage: the many
	// are expropriated by the few.
	NegationStageCapitalistExpropriation = "capitalist-expropriation"
	// NegationStageSocialisedProperty is the third stage: the few are
	// expropriated by the many.
	NegationStageSocialisedProperty = "socialised-property"
)

// PettyProperty is the first stage: self-earned property based on the
// labourer's own use of the means of labour. Producers is the count of
// independent producers; CapitalPerProducerPence is the typical stock
// each holds.
type PettyProperty struct {
	Producers               int64 `json:"producers"`
	CapitalPerProducerPence Pence `json:"capital_per_producer_pence"`
}

// CapitalistPrivateProperty is the second stage: property that rests on
// the exploitation of others' unpaid labour. Firms is the number of
// capitalist enterprises; TotalCapitalPence is their aggregate stock;
// WageLabourers is the proletariat dependent on them.
type CapitalistPrivateProperty struct {
	Firms             int64 `json:"firms"`
	TotalCapitalPence Pence `json:"total_capital_pence"`
	WageLabourers     int64 `json:"wage_labourers"`
}

// CentralisationStep records one event in the absorption of many
// capitalists by few — a merger, takeover, or competitive elimination.
// FirmsAbsorbed is how many capitalists were swallowed; the survivor
// gains CapitalConcentratedPence of their stock.
type CentralisationStep struct {
	StepIndex                int64 `json:"step_index"`
	FirmsAbsorbed            int64 `json:"firms_absorbed"`
	CapitalConcentratedPence Pence `json:"capital_concentrated_pence"`
}

// Negation names one of the three dialectical stages. Stage is one of
// NegationStage* constants; Description is the textual gloss.
type Negation struct {
	Stage       string `json:"stage"`
	Description string `json:"description"`
}

// AccumulationTrajectoryID is a 96-bit hex identifier for a persisted
// long-run trajectory.
type AccumulationTrajectoryID string

// IsZero reports whether the ID is unset.
func (id AccumulationTrajectoryID) IsZero() bool { return id == "" }

// NewAccumulationTrajectoryID generates a 96-bit hex identifier via
// crypto/rand.
func NewAccumulationTrajectoryID() AccumulationTrajectoryID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return AccumulationTrajectoryID(hex.EncodeToString(b))
}

// CentralisationStepID is a 96-bit hex identifier for a single
// centralisation_steps row within a trajectory.
type CentralisationStepID string

// NewCentralisationStepID generates a 96-bit hex identifier via
// crypto/rand.
func NewCentralisationStepID() CentralisationStepID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return CentralisationStepID(hex.EncodeToString(b))
}

// AccumulationTrajectory is the long-run series of centralisation steps
// from an initial CapitalistPrivateProperty to a final state with fewer
// firms holding a larger capital. ReserveArmySize is the displaced
// labour stratum generated alongside the centralisation; it links the
// chapter back to Ch. 25's general law.
type AccumulationTrajectory struct {
	ID                  AccumulationTrajectoryID `json:"id"`
	Name                string                   `json:"name"`
	InitialFirms        int64                    `json:"initial_firms"`
	InitialCapitalPence Pence                    `json:"initial_capital_pence"`
	Steps               []CentralisationStep     `json:"steps"`
	FinalFirms          int64                    `json:"final_firms"`
	FinalCapitalPence   Pence                    `json:"final_capital_pence"`
	ReserveArmySize     int64                    `json:"reserve_army_size"`
	CreatedAt           time.Time                `json:"created_at"`
}

var (
	// ErrNegationStagesCount is returned when NegationOfNegation receives
	// fewer or more than three stages.
	ErrNegationStagesCount = errors.New("simulation: negation-of-negation requires exactly three stages")
	// ErrNegationStagesOrder is returned when the three stages are not in
	// the historical order petty -> capitalist -> socialised.
	ErrNegationStagesOrder = errors.New("simulation: stages must be ordered petty-property, capitalist-expropriation, socialised-property")
	// ErrTrajectoryNameRequired is returned when a persisted trajectory
	// lacks a name.
	ErrTrajectoryNameRequired = errors.New("simulation: accumulation trajectory name is required")
	// ErrTrajectoryFirmsPositive is returned when initial firms is not
	// strictly positive.
	ErrTrajectoryFirmsPositive = errors.New("simulation: initial firms must be positive")
	// ErrCentralisationStepsPositive is returned when the requested step
	// count is not strictly positive.
	ErrCentralisationStepsPositive = errors.New("simulation: steps must be positive")
	// ErrCentralisationRateRange is returned when the absorption rate per
	// step is outside (0, 1].
	ErrCentralisationRateRange = errors.New("simulation: absorption_rate must be in (0, 1]")
	// ErrCentralisationGrowthRange is returned when capital growth per
	// step is negative.
	ErrCentralisationGrowthRange = errors.New("simulation: capital_growth_rate must be >= 0")
)

// Validate checks that the trajectory has a name and a positive starting
// population of firms.
func (t AccumulationTrajectory) Validate() error {
	if t.Name == "" {
		return ErrTrajectoryNameRequired
	}
	if t.InitialFirms <= 0 {
		return ErrTrajectoryFirmsPositive
	}
	return nil
}

// NegationOfNegation returns the third dialectical stage — socialised
// property — given the prior two. It validates that exactly three
// Negation entries are provided in the historical order and that the
// final stage carries the socialised-property tag.
//
// The function exists to make the dialectic explicit in the type system,
// not to compute anything Marx hasn't already named.
func NegationOfNegation(stages []Negation) (Negation, error) {
	if len(stages) != 3 {
		return Negation{}, ErrNegationStagesCount
	}
	if stages[0].Stage != NegationStagePettyProperty ||
		stages[1].Stage != NegationStageCapitalistExpropriation ||
		stages[2].Stage != NegationStageSocialisedProperty {
		return Negation{}, ErrNegationStagesOrder
	}
	return stages[2], nil
}

// DialecticalSequence returns the canonical three-stage sequence with
// glosses drawn from Capital Vol. I, Ch. 32. The handler layer uses this
// for the stateless GET /v1/accumulation/negation-of-negation endpoint.
func DialecticalSequence() []Negation {
	return []Negation{
		{
			Stage:       NegationStagePettyProperty,
			Description: "Self-earned private property based on the fusing of the individual labourer with the conditions of his labour.",
		},
		{
			Stage:       NegationStageCapitalistExpropriation,
			Description: "Capitalist private property: many producers expropriated by few capitalists, resting on the exploitation of nominally free wage-labour.",
		},
		{
			Stage:       NegationStageSocialisedProperty,
			Description: "Individual property based on co-operation and possession in common of the land and means of production; the expropriators are expropriated.",
		},
	}
}

// RunCentralisation simulates progressive absorption of firms over the
// given number of steps. Each step a fraction `absorptionRate` of the
// surviving firms is swallowed by competitors; the absorbed capital is
// added to the survivors and the total capital grows by
// `capitalGrowthRate` (representing accumulation on top of centralisation).
//
// The function never reaches FinalFirms == 0; it floors absorption at one
// firm taken per step and stops short of complete monopoly. Marx's "knell
// of capitalist private property" is an asymptote, not a fixed point.
//
// FirmsAbsorbed and CapitalConcentratedPence per step are guaranteed to
// be strictly positive when the trajectory has any steps at all.
func RunCentralisation(initial CapitalistPrivateProperty, steps int64, absorptionRate float64, capitalGrowthRate float64) AccumulationTrajectory {
	traj := AccumulationTrajectory{
		InitialFirms:        initial.Firms,
		InitialCapitalPence: initial.TotalCapitalPence,
		Steps:               []CentralisationStep{},
		FinalFirms:          initial.Firms,
		FinalCapitalPence:   initial.TotalCapitalPence,
		ReserveArmySize:     0,
	}
	if steps <= 0 || initial.Firms <= 1 || absorptionRate <= 0 {
		return traj
	}

	firms := initial.Firms
	capital := initial.TotalCapitalPence
	labourers := initial.WageLabourers

	for i := int64(0); i < steps; i++ {
		if firms <= 1 {
			break
		}
		absorbed := int64(math.Floor(float64(firms) * absorptionRate))
		if absorbed < 1 {
			absorbed = 1
		}
		if absorbed >= firms {
			absorbed = firms - 1
		}
		// CapitalConcentratedPence is the share of the *initial* capital
		// pool held by firms absorbed at this step. Apportioning against
		// the initial pool keeps the invariant Σ(concentrated) ≤ initial
		// capital, so FinalCapitalPence (initial + growth) always exceeds
		// the sum — concentration redistributes value rather than destroying
		// it.
		concentrated := Pence(int64(math.Round(float64(initial.TotalCapitalPence) * (float64(absorbed) / float64(initial.Firms)))))
		if concentrated <= 0 {
			concentrated = 1
		}
		traj.Steps = append(traj.Steps, CentralisationStep{
			StepIndex:                i + 1,
			FirmsAbsorbed:            absorbed,
			CapitalConcentratedPence: concentrated,
		})
		firms -= absorbed
		if capitalGrowthRate > 0 {
			capital = Pence(int64(math.Round(float64(capital) * (1.0 + capitalGrowthRate))))
		}
	}

	traj.FinalFirms = firms
	traj.FinalCapitalPence = capital
	if initial.Firms > 0 {
		avgPerFirm := labourers / initial.Firms
		traj.ReserveArmySize = (initial.Firms - firms) * avgPerFirm
	}
	return traj
}
