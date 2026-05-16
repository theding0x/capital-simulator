package simulation

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// Ch. 28 — Bloody Legislation Against the Expropriated.
//
// After the enclosures stripped producers from the land, the bourgeois state
// applied two complementary instruments of legal coercion: vagrancy laws that
// forced the newly landless to accept any wage on pain of whipping,
// branding, or execution; and wage statutes that kept the price of labour-power
// below its market rate by fixing legal ceilings. Together they constitute the
// "extra-economic" foundation that primitive accumulation required before the
// "dull compulsion of economic relations" could operate on its own.

// WageStatuteID is a 96-bit hex identifier for a persisted wage statute.
type WageStatuteID string

// IsZero reports whether the ID is unset.
func (id WageStatuteID) IsZero() bool { return id == "" }

// NewWageStatuteID generates a 96-bit hex identifier via crypto/rand.
func NewWageStatuteID() WageStatuteID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return WageStatuteID(hex.EncodeToString(b))
}

// VagrancyLawID is a 96-bit hex identifier for a persisted vagrancy law.
type VagrancyLawID string

// IsZero reports whether the ID is unset.
func (id VagrancyLawID) IsZero() bool { return id == "" }

// NewVagrancyLawID generates a 96-bit hex identifier via crypto/rand.
func NewVagrancyLawID() VagrancyLawID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return VagrancyLawID(hex.EncodeToString(b))
}

var (
	// WageStatute errors
	ErrWageStatutePeriodRequired       = errors.New("simulation: wage statute period is required")
	ErrWageStatuteJurisdictionRequired = errors.New("simulation: wage statute jurisdiction is required")
	ErrWageStatutePenaltyRequired      = errors.New("simulation: wage statute enforcement_penalty is required")
	ErrWageStatuteNegativePence        = errors.New("simulation: wage pence values cannot be negative")
	ErrWageStatuteBoundsInconsistent   = errors.New("simulation: max_wage_pence must be >= min_wage_pence")

	// VagrancyLaw errors
	ErrVagrancyPeriodRequired       = errors.New("simulation: vagrancy law period is required")
	ErrVagrancyJurisdictionRequired = errors.New("simulation: vagrancy law jurisdiction is required")
	ErrVagrancyPunishmentRequired   = errors.New("simulation: vagrancy law punishment is required")
	ErrVagrancyTargetRequired       = errors.New("simulation: vagrancy law target_population is required")

	// StatutoryWage errors
	ErrMarketWageZero = errors.New("simulation: market_wage_pence must be positive to compute deviation")
)

// WageStatute is a legal maximum or minimum wage fixed by Act of Parliament.
// Until the 1813 repeal the state always set a ceiling (max) and no floor;
// max_wage_pence is the legal cap, min_wage_pence is the legal floor (0 = none).
type WageStatute struct {
	ID                 WageStatuteID     `json:"id"`
	HistoricalStageID  HistoricalStageID `json:"historical_stage_id"`
	Period             string            `json:"period"`
	Jurisdiction       string            `json:"jurisdiction"`
	MaxWagePence       int64             `json:"max_wage_pence"`
	MinWagePence       int64             `json:"min_wage_pence"`
	EnforcementPenalty string            `json:"enforcement_penalty"`
	CreatedAt          time.Time         `json:"created_at"`
}

// Validate enforces invariants on a wage statute.
func (w WageStatute) Validate() error {
	if w.Period == "" {
		return ErrWageStatutePeriodRequired
	}
	if w.Jurisdiction == "" {
		return ErrWageStatuteJurisdictionRequired
	}
	if w.EnforcementPenalty == "" {
		return ErrWageStatutePenaltyRequired
	}
	if w.MaxWagePence < 0 || w.MinWagePence < 0 {
		return ErrWageStatuteNegativePence
	}
	if w.MaxWagePence < w.MinWagePence {
		return ErrWageStatuteBoundsInconsistent
	}
	return nil
}

// VagrancyLaw is a statute coercing the dispossessed into wage labour by
// criminalising vagrancy — i.e. being without employment and property.
type VagrancyLaw struct {
	ID                VagrancyLawID     `json:"id"`
	HistoricalStageID HistoricalStageID `json:"historical_stage_id"`
	Period            string            `json:"period"`
	Jurisdiction      string            `json:"jurisdiction"`
	Punishment        string            `json:"punishment"`
	TargetPopulation  string            `json:"target_population"`
	CreatedAt         time.Time         `json:"created_at"`
}

// Validate enforces invariants on a vagrancy law.
func (v VagrancyLaw) Validate() error {
	if v.Period == "" {
		return ErrVagrancyPeriodRequired
	}
	if v.Jurisdiction == "" {
		return ErrVagrancyJurisdictionRequired
	}
	if v.Punishment == "" {
		return ErrVagrancyPunishmentRequired
	}
	if v.TargetPopulation == "" {
		return ErrVagrancyTargetRequired
	}
	return nil
}

// StatutoryWage compares a legally imposed wage with the prevailing market
// rate. Deviation is negative when the statutory cap is below the market wage
// (wage suppression) and zero when no statute applies.
// Deviation = (ActedWagePence - MarketWagePence) / MarketWagePence.
type StatutoryWage struct {
	ActedWagePence  int64   `json:"acted_wage_pence"`
	MarketWagePence int64   `json:"market_wage_pence"`
	Deviation       float64 `json:"deviation"`
}

// ComputeStatutoryWage returns a StatutoryWage for the given acted and market
// wages. MarketWagePence must be positive; returns ErrMarketWageZero otherwise.
func ComputeStatutoryWage(acted, market int64) (StatutoryWage, error) {
	if market <= 0 {
		return StatutoryWage{}, ErrMarketWageZero
	}
	dev := float64(acted-market) / float64(market)
	return StatutoryWage{
		ActedWagePence:  acted,
		MarketWagePence: market,
		Deviation:       dev,
	}, nil
}

// LabourDisciplineRegime is the aggregate view of state coercion for a given
// historical stage: the full roster of wage statutes and vagrancy laws in force,
// plus a summary of enforcement mechanisms. When both slices are empty the
// regime represents the post-1825 norm where "dull compulsion of economic
// relations" alone disciplines labour.
type LabourDisciplineRegime struct {
	Period       string        `json:"period"`
	Mechanisms   []string      `json:"mechanisms"`
	WageStatutes []WageStatute `json:"wage_statutes"`
	VagrancyLaws []VagrancyLaw `json:"vagrancy_laws"`
}

// AssembleRegime builds a LabourDisciplineRegime from a stage name and the
// statutes and laws attached to it. Mechanisms is deduped from enforcement
// penalties and punishments across all records.
func AssembleRegime(stagePeriod string, ws []WageStatute, vl []VagrancyLaw) LabourDisciplineRegime {
	seen := make(map[string]struct{})
	var mechanisms []string
	for _, w := range ws {
		if _, ok := seen[w.EnforcementPenalty]; !ok && w.EnforcementPenalty != "" {
			seen[w.EnforcementPenalty] = struct{}{}
			mechanisms = append(mechanisms, w.EnforcementPenalty)
		}
	}
	for _, v := range vl {
		if _, ok := seen[v.Punishment]; !ok && v.Punishment != "" {
			seen[v.Punishment] = struct{}{}
			mechanisms = append(mechanisms, v.Punishment)
		}
	}
	if mechanisms == nil {
		mechanisms = []string{}
	}
	wsCopy := ws
	if wsCopy == nil {
		wsCopy = []WageStatute{}
	}
	vlCopy := vl
	if vlCopy == nil {
		vlCopy = []VagrancyLaw{}
	}
	return LabourDisciplineRegime{
		Period:       stagePeriod,
		Mechanisms:   mechanisms,
		WageStatutes: wsCopy,
		VagrancyLaws: vlCopy,
	}
}