package simulation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Ch. 26 — The Secret of Primitive Accumulation.
//
// Primitive accumulation is the historical pre-condition of capitalist
// production: a process that separates producers from the means of
// production, generating on one side a free wage-labourer class and on
// the other a concentrated ownership of capital. The types below let the
// simulation seed a starting CapitalStock and a free-proletarian supply
// from a named historical stage rather than treating them as given.

// HistoricalStageID is a 96-bit hex identifier for a persisted stage.
type HistoricalStageID string

// IsZero reports whether the ID is unset.
func (id HistoricalStageID) IsZero() bool { return id == "" }

// NewHistoricalStageID generates a 96-bit hex identifier via crypto/rand.
func NewHistoricalStageID() HistoricalStageID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return HistoricalStageID(hex.EncodeToString(b))
}

// PrimitiveAccumulation records a single historical episode: a period,
// the method by which producers were separated from their means
// (enclosure, expropriation, colonial seizure), the count of labourers
// affected, and the resulting concentration of capital.
type PrimitiveAccumulation struct {
	Period                string `json:"period"`
	Method                string `json:"method"`
	LabourersExpropriated int64  `json:"labourers_expropriated"`
	CapitalFormed         Pence  `json:"capital_formed"`
}

// ErrCapitalFormedNonPositive — primitive accumulation always results in
// a positive capital stock, however it was obtained.
var ErrCapitalFormedNonPositive = errors.New("simulation: primitive accumulation must form positive capital")

// ErrLabourersNegative — labourers_expropriated cannot be negative.
var ErrLabourersNegative = errors.New("simulation: labourers_expropriated cannot be negative")

// ErrSeparationMismatch — the count of dispossessed must be conserved:
// every displaced producer reappears as a free proletarian.
var ErrSeparationMismatch = errors.New("simulation: free_proletarians must equal displaced_workers")

// ErrStageEmpty — a historical stage must include at least one episode.
var ErrStageEmpty = errors.New("simulation: historical stage must include at least one primitive accumulation")

// Validate enforces invariants on a single primitive-accumulation episode.
func (p PrimitiveAccumulation) Validate() error {
	if p.CapitalFormed <= 0 {
		return ErrCapitalFormedNonPositive
	}
	if p.LabourersExpropriated < 0 {
		return ErrLabourersNegative
	}
	return nil
}

// ProducerSeparation records the structural divide on which wage-labour
// rests: a pre-capitalist population of working producers is reduced by
// dispossession, and every displaced producer reappears on the labour
// market as a free proletarian — the invariant
// FreeProletarians == DisplacedWorkers is the conservation law of
// primitive accumulation.
type ProducerSeparation struct {
	PreCapitalistWorkers int64 `json:"pre_capitalist_workers"`
	DisplacedWorkers     int64 `json:"displaced_workers"`
	FreeProletarians     int64 `json:"free_proletarians"`
}

// Validate enforces the conservation invariant.
func (s ProducerSeparation) Validate() error {
	if s.PreCapitalistWorkers < 0 || s.DisplacedWorkers < 0 || s.FreeProletarians < 0 {
		return ErrLabourersNegative
	}
	if s.FreeProletarians != s.DisplacedWorkers {
		return ErrSeparationMismatch
	}
	return nil
}

// HistoricalStage is a named epoch grouping primitive-accumulation
// episodes. England 15th–18th c. is the canonical example: enclosure,
// statutory wage-fixing, and colonial plunder together build up the
// initial capital and free-proletarian supply that the Ch. 23/24
// reproduction loop then takes as its starting point.
type HistoricalStage struct {
	ID                     HistoricalStageID       `json:"id"`
	Name                   string                  `json:"name"`
	Description            string                  `json:"description"`
	PrimitiveAccumulations []PrimitiveAccumulation `json:"primitive_accumulations"`
	CreatedAt              time.Time               `json:"created_at"`
}

// Validate enforces invariants on the stage and every episode it contains.
func (h HistoricalStage) Validate() error {
	if h.Name == "" {
		return errors.New("simulation: historical stage name is required")
	}
	if len(h.PrimitiveAccumulations) == 0 {
		return ErrStageEmpty
	}
	for _, p := range h.PrimitiveAccumulations {
		if err := p.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// TotalCapitalFormed sums the capital formed across every episode in the stage.
func (h HistoricalStage) TotalCapitalFormed() Pence {
	var total Pence
	for _, p := range h.PrimitiveAccumulations {
		total += p.CapitalFormed
	}
	return total
}

// TotalLabourersExpropriated sums the labourers dispossessed across every
// episode in the stage — these become the supply side of the labour market.
func (h HistoricalStage) TotalLabourersExpropriated() int64 {
	var total int64
	for _, p := range h.PrimitiveAccumulations {
		total += p.LabourersExpropriated
	}
	return total
}

// SeparationFromStage projects a HistoricalStage onto a ProducerSeparation,
// given the pre-capitalist working population the stage acts upon.
// FreeProletarians is set equal to DisplacedWorkers so the invariant holds
// by construction; if the stage dispossesses more workers than exist in
// the pre-capitalist population, the count is clamped to the population.
func SeparationFromStage(h HistoricalStage, preCapitalistWorkers int64) ProducerSeparation {
	displaced := h.TotalLabourersExpropriated()
	if displaced > preCapitalistWorkers {
		displaced = preCapitalistWorkers
	}
	return ProducerSeparation{
		PreCapitalistWorkers: preCapitalistWorkers,
		DisplacedWorkers:     displaced,
		FreeProletarians:     displaced,
	}
}

// SeedCapitalStock derives a starting CapitalStock from a historical stage,
// given the organic composition ratio applied to the aggregate capital
// formed. Constant capital takes ratio * total; variable capital takes the
// complement, so constant + variable == total in all cases.
func SeedCapitalStock(h HistoricalStage, compositionRatio float64) CapitalStock {
	total := h.TotalCapitalFormed()
	if total <= 0 {
		return CapitalStock{}
	}
	if compositionRatio < 0 {
		compositionRatio = 0
	}
	if compositionRatio >= 1 {
		compositionRatio = 0.99
	}
	constant := Pence(int64(float64(total) * compositionRatio))
	variable := total - constant
	return CapitalStock{ConstantCapital: constant, VariableCapital: variable}
}
