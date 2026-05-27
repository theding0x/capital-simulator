package circulation

import "time"

// NaturalProcessKind names the kinds of natural-process production span.
type NaturalProcessKind string

const (
	NaturalRipening     NaturalProcessKind = "ripening"
	NaturalFermentation NaturalProcessKind = "fermentation"
	NaturalTanning      NaturalProcessKind = "tanning"
	NaturalDrying       NaturalProcessKind = "drying"
	NaturalOther        NaturalProcessKind = "other"

	// Ch. 13 (Vol. II) additions — extend the enum with the five long-production-time kinds.
	NaturalAgriculturalGrowth NaturalProcessKind = "agricultural_growth" // corn: seeding → harvest (~280 days)
	NaturalForestGrowth       NaturalProcessKind = "forest_growth"       // timber: 12–100 years
	NaturalAnimalMaturation   NaturalProcessKind = "animal_maturation"   // cattle / horses / wine ageing
	NaturalChemical           NaturalProcessKind = "chemical"            // porcelain firing, bleaching
	NaturalSeasonalRest       NaturalProcessKind = "seasonal_rest"       // winter fallow
)

// NaturalProcessSpanID is a 96-bit hex identifier.
type NaturalProcessSpanID string

func (id NaturalProcessSpanID) IsZero() bool { return id == "" }

// NewNaturalProcessSpanID generates a 96-bit hex identifier via crypto/rand.
func NewNaturalProcessSpanID() NaturalProcessSpanID { return NaturalProcessSpanID(newHexID()) }

// NaturalProcessSpan represents a production sub-span during which the labour
// process has stopped but the natural process continues (grain ripening, wine
// fermenting, leather tanning). Labour is absent; no surplus-value is created.
// §: "grain, wine fermenting in the cellar, tanneries … exposed to chemical processes."
type NaturalProcessSpan struct {
	ID                  NaturalProcessSpanID `json:"id"`
	TurnoverTimeID      TurnoverTimeID       `json:"turnover_time_id"`
	IndustrialCapitalID IndustrialCapitalID  `json:"industrial_capital_id"`
	Process             NaturalProcessKind   `json:"process"`
	DurationNanos       int64                `json:"duration_nanos"`
	StartedAt           time.Time            `json:"started_at"`
}

// LatentProductiveCapitalID is a 96-bit hex identifier.
type LatentProductiveCapitalID string

func (id LatentProductiveCapitalID) IsZero() bool { return id == "" }

// NewLatentProductiveCapitalID generates a 96-bit hex identifier via crypto/rand.
func NewLatentProductiveCapitalID() LatentProductiveCapitalID {
	return LatentProductiveCapitalID(newHexID())
}

// LatentProductiveCapital is means-of-production held in readiness at the
// production site but not yet entering the production process. While
// EnteredProductionAt is nil, this capital creates neither product nor value.
// §: "That part of the latent productive capital held in readiness only as a
// requisite for the productive process acts as a creator of neither products
// nor value. It is fallow capital."
type LatentProductiveCapital struct {
	ID                  LatentProductiveCapitalID `json:"id"`
	TurnoverTimeID      TurnoverTimeID            `json:"turnover_time_id"`
	IndustrialCapitalID IndustrialCapitalID       `json:"industrial_capital_id"`
	Pence               Pence                     `json:"pence"`
	HeldAt              time.Time                 `json:"held_at"`
	// EnteredProductionAt is nil while the capital is still latent.
	EnteredProductionAt *time.Time `json:"entered_production_at,omitempty"`
}

// IsLatent reports whether this capital has not yet entered the production process.
func (l LatentProductiveCapital) IsLatent() bool { return l.EnteredProductionAt == nil }
