// Package workperiod — Vol. II Ch. 13 additions: The Time of Production.
// Ch. 13 deepens Ch. 12 by separating the span during which capital is locked
// in natural processes (biological growth, fermentation, chemical change) from
// the labour-time that actually creates value. The key result: only
// LabourTimeNanos creates surplus-value; GapNanos is latent, value-free time.
package workperiod

import (
	"errors"
	"time"
)

// --- NaturalProcessKind -------------------------------------------------

// NaturalProcessKind names the kind of natural-process production span.
// The first five values mirror the market-service Ch. 5 enum; the Ch. 13
// additions are listed beneath them. Defined locally because Go's internal
// package rule prevents cross-service imports.
type NaturalProcessKind string

const (
	// Ch. 5 equivalents — kept here so simulation-engine types are self-contained.
	NaturalRipening     NaturalProcessKind = "ripening"
	NaturalFermentation NaturalProcessKind = "fermentation"
	NaturalTanning      NaturalProcessKind = "tanning"
	NaturalDrying       NaturalProcessKind = "drying"
	NaturalOther        NaturalProcessKind = "other"

	// Ch. 13 additions.
	NaturalAgriculturalGrowth NaturalProcessKind = "agricultural_growth" // corn: seeding → harvest
	NaturalForestGrowth       NaturalProcessKind = "forest_growth"       // timber: 12–100 years
	NaturalAnimalMaturation   NaturalProcessKind = "animal_maturation"   // cattle / horses / wine ageing
	NaturalChemical           NaturalProcessKind = "chemical"            // porcelain firing, bleaching
	NaturalSeasonalRest       NaturalProcessKind = "seasonal_rest"       // winter fallow
)

// validNaturalProcessKind reports whether k is a recognised constant.
func validNaturalProcessKind(k NaturalProcessKind) bool {
	switch k {
	case NaturalRipening, NaturalFermentation, NaturalTanning, NaturalDrying, NaturalOther,
		NaturalAgriculturalGrowth, NaturalForestGrowth, NaturalAnimalMaturation,
		NaturalChemical, NaturalSeasonalRest:
		return true
	}
	return false
}

// --- ProductionLabourGap ------------------------------------------------

// ProductionLabourGapID is a 96-bit hex identifier for a ProductionLabourGap.
type ProductionLabourGapID string

// NewProductionLabourGapID returns a new random ProductionLabourGapID.
func NewProductionLabourGapID() ProductionLabourGapID {
	return ProductionLabourGapID(newHexID())
}

// ProductionLabourGap records the split between the total production-time span
// and the labour-time contained within it for a given working period.
// GapNanos is the latent span during which capital is locked in a natural
// process and no surplus-value is created.
// Invariant: ProductionTimeNanos ≥ LabourTimeNanos; GapNanos == ProductionTimeNanos - LabourTimeNanos.
type ProductionLabourGap struct {
	ID                  ProductionLabourGapID `json:"id"`
	WorkingPeriodID     WorkingPeriodID       `json:"working_period_id"`
	ProductionTimeNanos int64                 `json:"production_time_nanos"`
	LabourTimeNanos     int64                 `json:"labour_time_nanos"`
	GapNanos            int64                 `json:"gap_nanos"`
	CreatedAt           time.Time             `json:"created_at"`
}

// Validate reports whether the gap is internally consistent.
func (g ProductionLabourGap) Validate() error {
	if g.WorkingPeriodID == "" {
		return errors.New("production_labour_gap: working_period_id is required")
	}
	if g.ProductionTimeNanos <= 0 {
		return errors.New("production_labour_gap: production_time_nanos must be positive")
	}
	if g.LabourTimeNanos <= 0 {
		return errors.New("production_labour_gap: labour_time_nanos must be positive")
	}
	if g.ProductionTimeNanos < g.LabourTimeNanos {
		return errors.New("production_labour_gap: production_time_nanos must be >= labour_time_nanos")
	}
	if g.GapNanos != g.ProductionTimeNanos-g.LabourTimeNanos {
		return errors.New("production_labour_gap: gap_nanos must equal production_time_nanos - labour_time_nanos")
	}
	return nil
}

// --- NaturalProcessActivation ------------------------------------------

// NaturalProcessActivationID is a 96-bit hex identifier for a NaturalProcessActivation.
type NaturalProcessActivationID string

// NewNaturalProcessActivationID returns a new random NaturalProcessActivationID.
func NewNaturalProcessActivationID() NaturalProcessActivationID {
	return NaturalProcessActivationID(newHexID())
}

// NaturalProcessActivation records a span during which a natural process runs
// for a given working period. When RequiresLabourMaintenance is true the
// maintenance labour-time counts as value-creating; the natural-process span
// itself does NOT increment any value field.
type NaturalProcessActivation struct {
	ID                        NaturalProcessActivationID `json:"id"`
	WorkingPeriodID           WorkingPeriodID            `json:"working_period_id"`
	Kind                      NaturalProcessKind         `json:"kind"`
	StartedAt                 time.Time                  `json:"started_at"`
	EndedAt                   *time.Time                 `json:"ended_at,omitempty"`
	RequiresLabourMaintenance bool                       `json:"requires_labour_maintenance"`
}

// Validate reports whether the activation is well-formed.
func (a NaturalProcessActivation) Validate() error {
	if a.WorkingPeriodID == "" {
		return errors.New("natural_process_activation: working_period_id is required")
	}
	if !validNaturalProcessKind(a.Kind) {
		return errors.New("natural_process_activation: unknown kind: " + string(a.Kind))
	}
	if a.StartedAt.IsZero() {
		return errors.New("natural_process_activation: started_at is required")
	}
	if a.EndedAt != nil && !a.EndedAt.After(a.StartedAt) {
		return errors.New("natural_process_activation: ended_at must be after started_at")
	}
	return nil
}

// --- NaturalSubject -----------------------------------------------------

// NaturalSubjectID is a 96-bit hex identifier for a NaturalSubject.
type NaturalSubjectID string

// NewNaturalSubjectID returns a new random NaturalSubjectID.
func NewNaturalSubjectID() NaturalSubjectID {
	return NaturalSubjectID(newHexID())
}

// NaturalSubject describes the commodity or material undergoing a natural
// process within a working period. IsLiving distinguishes animate subjects
// (cattle, timber) from inert ones (wine, porcelain, cloth), though both
// use the same NaturalProcessActivation logic — it is a documentation tag only.
type NaturalSubject struct {
	ID              NaturalSubjectID `json:"id"`
	WorkingPeriodID WorkingPeriodID  `json:"working_period_id"`
	CommodityID     string           `json:"commodity_id"`
	Description     string           `json:"description"`
	IsLiving        bool             `json:"is_living"`
	CreatedAt       time.Time        `json:"created_at"`
}

// Validate reports whether the subject is well-formed.
func (s NaturalSubject) Validate() error {
	if s.WorkingPeriodID == "" {
		return errors.New("natural_subject: working_period_id is required")
	}
	if s.CommodityID == "" {
		return errors.New("natural_subject: commodity_id is required")
	}
	if s.Description == "" {
		return errors.New("natural_subject: description is required")
	}
	return nil
}

// --- NaturalProcessIndustry ---------------------------------------------

// NaturalProcessIndustryID is a 96-bit hex identifier for a NaturalProcessIndustry.
type NaturalProcessIndustryID string

// NewNaturalProcessIndustryID returns a new random NaturalProcessIndustryID.
func NewNaturalProcessIndustryID() NaturalProcessIndustryID {
	return NaturalProcessIndustryID(newHexID())
}

// NaturalProcessIndustry is a descriptive benchmark for an industry whose
// production-time is dominated by a natural process. RatioBasisPoints encodes
// the labour/production ratio as 10000 × TypicalLabourDaysCount /
// TypicalProductionDaysCount; a low value means most production-time is
// latent (no value is being created). Benchmarks are descriptive only — not
// validated against any specific WorkingPeriod row.
type NaturalProcessIndustry struct {
	ID                         NaturalProcessIndustryID `json:"id"`
	IndustryName               string                   `json:"industry_name"`
	TypicalProductionDaysCount int64                    `json:"typical_production_days_count"`
	TypicalLabourDaysCount     int64                    `json:"typical_labour_days_count"`
	RatioBasisPoints           int64                    `json:"ratio_basis_points"`
	CreatedAt                  time.Time                `json:"created_at"`
}

// Validate reports whether the benchmark is well-formed.
func (i NaturalProcessIndustry) Validate() error {
	if i.IndustryName == "" {
		return errors.New("natural_process_industry: industry_name is required")
	}
	if i.TypicalProductionDaysCount <= 0 {
		return errors.New("natural_process_industry: typical_production_days_count must be positive")
	}
	if i.TypicalLabourDaysCount <= 0 {
		return errors.New("natural_process_industry: typical_labour_days_count must be positive")
	}
	if i.TypicalProductionDaysCount < i.TypicalLabourDaysCount {
		return errors.New("natural_process_industry: typical_production_days_count must be >= typical_labour_days_count")
	}
	if i.RatioBasisPoints < 0 || i.RatioBasisPoints > 10000 {
		return errors.New("natural_process_industry: ratio_basis_points must be in [0, 10000]")
	}
	return nil
}
