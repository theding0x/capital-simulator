// Package reproduction — Vol. II Ch. 17: The Circulation of Surplus-Value.
//
// Marx poses the question at the social level: where does the money come from
// to realise the aggregate surplus-value produced by all capitalists in a
// year? Individual-capital analysis cannot answer this — the answer only
// becomes possible when we treat the total social capital as a whole.
//
// The chapter identifies five possible sources through which surplus-value
// can be realised:
//   - capitalist consumption (revenue expenditure)
//   - capitalised surplus (re-investment)
//   - gold production (the only source that adds to the social money-stock)
//   - foreign trade (abstracted away in Vol. II)
//   - credit (defers realisation; creates a matching claim)
//
// The puzzle (formally resolved in Chs. 20-21): ΣSurplus == ΣRealisation
// requires that the sources sum correctly at the aggregate level, yet no
// single capitalist can identify the source of their own surplus-value money.
//
// All monetary values are in pence (1 GBP = 240 pence).
package reproduction

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

// --- Sentinel errors --------------------------------------------------------

var (
	// ErrNegativeSurplus is returned when total_surplus_pence is negative.
	ErrNegativeSurplus = errors.New("reproduction: total_surplus_pence cannot be negative")

	// ErrSourceMixMismatch is returned when the realisation-source mix does not
	// sum to realised_surplus_pence.
	ErrSourceMixMismatch = errors.New("reproduction: realisation source mix does not sum to realised surplus")

	// ErrUnknownSource is returned when an unrecognised RealisationSource
	// constant is presented.
	ErrUnknownSource = errors.New("reproduction: unknown realisation source")

	// ErrNegativePence is returned when a pence value is negative.
	ErrNegativePence = errors.New("reproduction: pence value cannot be negative")

	// ErrRealisationExceedsTotal is returned when adding a source entry would
	// push realised surplus above total surplus.
	ErrRealisationExceedsTotal = errors.New("reproduction: realised surplus would exceed total")
)

// --- Pence ------------------------------------------------------------------

// Pence is the canonical monetary unit (1 GBP = 240 pence).
type Pence int64

// --- ID types ---------------------------------------------------------------

// SurplusCirculationID identifies a SurplusCirculation record.
type SurplusCirculationID string

// IsZero reports whether the ID is the zero value.
func (id SurplusCirculationID) IsZero() bool { return id == "" }

// NewSurplusCirculationID returns a 96-bit random hex identifier.
func NewSurplusCirculationID() SurplusCirculationID {
	return SurplusCirculationID(newHexID())
}

// --- RealisationSource ------------------------------------------------------

// RealisationSource identifies the economic mechanism through which
// surplus-value is converted into money.
type RealisationSource string

const (
	// SourceCapitalistConsumption — capitalists spend revenue on personal
	// consumption; this withdraws money from circulation to realise surplus.
	SourceCapitalistConsumption RealisationSource = "capitalist_consumption"

	// SourceCapitalisedSurplus — surplus is re-invested rather than consumed;
	// new means of production or labour-power are purchased.
	SourceCapitalisedSurplus RealisationSource = "capitalised_surplus"

	// SourceForeignTrade — exports absorb domestic surplus-value through
	// foreign-exchange inflows. Abstracted away in Vol. II (always 0).
	SourceForeignTrade RealisationSource = "foreign_trade"

	// SourceGoldProduction — production of gold (the money commodity) is the
	// only mechanism that expands the social money-stock directly.
	SourceGoldProduction RealisationSource = "gold_production"

	// SourceCredit — credit-money advances payment; creates a deferred
	// realisation claim (DeferredRealisationPence > 0).
	SourceCredit RealisationSource = "credit"
)

// validSources is the canonical set of allowed RealisationSource values.
var validSources = map[RealisationSource]bool{
	SourceCapitalistConsumption: true,
	SourceCapitalisedSurplus:    true,
	SourceForeignTrade:          true,
	SourceGoldProduction:        true,
	SourceCredit:                true,
}

// Validate returns ErrUnknownSource if the source is not a recognised constant.
func (s RealisationSource) Validate() error {
	if !validSources[s] {
		return ErrUnknownSource
	}
	return nil
}

// --- RealisationSourceEntry -------------------------------------------------

// RealisationSourceEntry records how much of a SurplusCirculation's surplus
// was realised through a particular mechanism.
//
// Invariants:
//   - Pence >= 0
//   - If Source == SourceCredit, DeferredRealisationPence == Pence
//   - If Source == SourceForeignTrade, Pence == 0 for Vol. II fixtures
type RealisationSourceEntry struct {
	SurplusCirculationID     SurplusCirculationID `json:"surplus_circulation_id"`
	Source                   RealisationSource    `json:"source"`
	Pence                    Pence                `json:"pence"`
	DeferredRealisationPence Pence                `json:"deferred_realisation_pence"`
}

// --- SurplusCirculation -----------------------------------------------------

// SurplusCirculation records the aggregate surplus-value produced in one
// period and the mix of mechanisms through which it was (or was not) realised.
//
// Invariant:
//
//	TotalSurplusPence == RealisedSurplusPence + UnrealisedSurplusPence
//	RealisedSurplusPence == Σ RealisationSources[i].Pence
type SurplusCirculation struct {
	ID                     SurplusCirculationID     `json:"id"`
	Period                 string                   `json:"period"`
	TotalSurplusPence      Pence                    `json:"total_surplus_pence"`
	RealisedSurplusPence   Pence                    `json:"realised_surplus_pence"`
	UnrealisedSurplusPence Pence                    `json:"unrealised_surplus_pence"`
	RealisationSources     []RealisationSourceEntry `json:"realisation_sources"`
}

// Validate checks the fundamental surplus-conservation invariant.
func (sc SurplusCirculation) Validate() error {
	if sc.TotalSurplusPence < 0 {
		return ErrNegativeSurplus
	}
	if sc.RealisedSurplusPence < 0 || sc.UnrealisedSurplusPence < 0 {
		return ErrNegativePence
	}
	if sc.TotalSurplusPence != sc.RealisedSurplusPence+sc.UnrealisedSurplusPence {
		return ErrSourceMixMismatch
	}
	var mixSum Pence
	for _, e := range sc.RealisationSources {
		mixSum += e.Pence
	}
	if mixSum != sc.RealisedSurplusPence {
		return ErrSourceMixMismatch
	}
	return nil
}

// AddRealisationSource applies a new source entry to a SurplusCirculation,
// updating RealisedSurplusPence and UnrealisedSurplusPence accordingly.
// Returns ErrRealisationExceedsTotal if the entry would push realised surplus
// above TotalSurplusPence.
func AddRealisationSource(sc SurplusCirculation, entry RealisationSourceEntry) (SurplusCirculation, error) {
	if entry.Pence < 0 {
		return sc, ErrNegativePence
	}
	if err := entry.Source.Validate(); err != nil {
		return sc, err
	}
	entry.SurplusCirculationID = sc.ID
	if entry.Source == SourceCredit {
		entry.DeferredRealisationPence = entry.Pence
	}
	newRealised := sc.RealisedSurplusPence + entry.Pence
	if newRealised > sc.TotalSurplusPence {
		return sc, ErrRealisationExceedsTotal
	}
	sc.RealisedSurplusPence = newRealised
	sc.UnrealisedSurplusPence = sc.TotalSurplusPence - newRealised
	sc.RealisationSources = append(sc.RealisationSources, entry)
	return sc, nil
}

// --- SocialCapitalAggregate -------------------------------------------------

// SocialCapitalAggregate is a period snapshot of the total social capital,
// aggregating all individual capitals in the economy.
//
// DepartmentIShareBP and DepartmentIIShareBP are placeholders here; they are
// populated fully in Chs. 20-21 (simple and extended reproduction schemes).
// Expressed in basis points (100% = 10000).
type SocialCapitalAggregate struct {
	Period                  string `json:"period"`
	TotalAdvancedPence      Pence  `json:"total_advanced_pence"`
	TotalAnnualOutputPence  Pence  `json:"total_annual_output_pence"`
	TotalAnnualSurplusPence Pence  `json:"total_annual_surplus_pence"`
	// DepartmentIShareBP is the fraction of total advanced capital in Dept. I
	// (means of production). Placeholder — populated in Ch. 20.
	DepartmentIShareBP int64 `json:"department_i_share_bp"`
	// DepartmentIIShareBP is the fraction in Dept. II (articles of consumption).
	// Placeholder — populated in Ch. 20.
	DepartmentIIShareBP int64 `json:"department_ii_share_bp"`
}

// ComputeDepartmentSharesBP returns the basis-point shares of total advanced
// capital (c + v) held by Department I and Department II. The two shares always
// sum to exactly 10000 — Department II's share is taken as the residual so that
// rounding can never break the partition. A non-positive total yields (0, 0).
//
// This is the writer Ch. 17 reserved DepartmentIShareBP / DepartmentIIShareBP
// for: the Ch. 20 (simple) and Ch. 21 (extended) reproduction ticks project
// their two departments' advanced capital back onto the period's
// SocialCapitalAggregate (issue #215).
func ComputeDepartmentSharesBP(deptIAdvancedPence, deptIIAdvancedPence int64) (deptIShareBP, deptIIShareBP int64) {
	total := deptIAdvancedPence + deptIIAdvancedPence
	if total <= 0 {
		return 0, 0
	}
	deptIShareBP = (deptIAdvancedPence*10000 + total/2) / total
	deptIIShareBP = 10000 - deptIShareBP
	return deptIShareBP, deptIIShareBP
}

// --- RealisationPuzzle ------------------------------------------------------

// RealisationPuzzle captures the canonical formulation of Marx's
// "where does the money come from?" question. One canonical puzzle row
// exists per chapter fixture; it is inserted by the seed migration.
type RealisationPuzzle struct {
	SurplusCirculationID SurplusCirculationID `json:"surplus_circulation_id"`
	PuzzleStatement      string               `json:"puzzle_statement"`
	ResolvedInChapter    int64                `json:"resolved_in_chapter"`
}

// --- helpers ----------------------------------------------------------------

func newHexID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(b)
}
