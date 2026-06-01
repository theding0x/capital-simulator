// classes.go — Vol. III Ch. 52 types: Classes. The final, incomplete chapter of
// Capital — the manuscript breaks off after identifying the three great classes of
// capitalist society (wage-labourers → wages, capitalists → profit, landowners →
// rent) and posing, but not answering, the question "what constitutes a class?".
// A class is defined by a shared relation to the means of production, not by the
// form of income alone — the answer Marx never wrote.
package revenue

import "time"

// SocialClassID identifies a social-class record.
type SocialClassID string

// NewSocialClassID returns a fresh 96-bit hex ID.
func NewSocialClassID() SocialClassID {
	return SocialClassID(newID())
}

// SocialClass is one of the three great classes, defined by its relation to the
// means of production and the revenue form that flows to it.
type SocialClass struct {
	ID                          SocialClassID `json:"id"`
	Name                        string        `json:"name"`
	RelationToMeansOfProduction string        `json:"relation_to_means_of_production"`
	PrimaryRevenueForm          string        `json:"primary_revenue_form"`
	HistoricalPrecondition      string        `json:"historical_precondition"`
	CreatedAt                   time.Time     `json:"created_at"`
}

// ClassIncomeSourceID identifies a class-income-source record.
type ClassIncomeSourceID string

// NewClassIncomeSourceID returns a fresh 96-bit hex ID.
func NewClassIncomeSourceID() ClassIncomeSourceID {
	return ClassIncomeSourceID(newID())
}

// ClassIncomeSource is the revenue form flowing to a class. Wages recover variable
// capital (SurplusValueComponentBP == 0); profit and rent are portions of surplus-
// value (SurplusValueComponentBP > 0).
type ClassIncomeSource struct {
	ID                      ClassIncomeSourceID `json:"id"`
	ClassID                 string              `json:"class_id"`
	RevenueForm             string              `json:"revenue_form"`
	SurplusValueComponentBP int64               `json:"surplus_value_component_bp"`
	CreatedAt               time.Time           `json:"created_at"`
}

// ClassTendencyID identifies a class-tendency record.
type ClassTendencyID string

// NewClassTendencyID returns a fresh 96-bit hex ID.
func NewClassTendencyID() ClassTendencyID {
	return ClassTendencyID(newID())
}

// ClassTendency records the observed direction of capitalist development for a class.
type ClassTendency struct {
	ID          ClassTendencyID `json:"id"`
	ClassID     string          `json:"class_id"`
	Description string          `json:"description"`
	Direction   string          `json:"direction"`
	CreatedAt   time.Time       `json:"created_at"`
}

// ClassAmbiguityID identifies a class-ambiguity record.
type ClassAmbiguityID string

// NewClassAmbiguityID returns a fresh 96-bit hex ID.
func NewClassAmbiguityID() ClassAmbiguityID {
	return ClassAmbiguityID(newID())
}

// ClassAmbiguity records an intermediate stratum that blurs the pure class lines.
// WhyNotAClass holds Marx's answer — seeded with "manuscript_breaks_off" to mark
// faithfully that the chapter (and Capital) ends before it is given.
type ClassAmbiguity struct {
	ID                ClassAmbiguityID `json:"id"`
	IntermediateGroup string           `json:"intermediate_group"`
	IncomeSource      string           `json:"income_source"`
	WhyNotAClass      string           `json:"why_not_a_class"`
	CreatedAt         time.Time        `json:"created_at"`
}

// ValidClassTendencyDirection reports whether d is a recognised direction of class
// development under capitalism.
func ValidClassTendencyDirection(d string) bool {
	switch d {
	case "expansion", "concentration", "dissolution", "subordination":
		return true
	default:
		return false
	}
}
