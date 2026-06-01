// distribution.go — Vol. III Ch. 51 types: Distribution Relations and Production
// Relations. Distribution relations (wages, profit, rent) are not a separate sphere
// but the "other side" of production relations, sharing their historically transitory
// character. Each distribution form rests on a specific production relation (wage-
// labour, capital, landed property) and corresponds to a component of new value
// (v, s-as-profit, s-as-rent). A crisis arises when the development of the productive
// forces contradicts the existing social form.
package revenue

import "time"

// ProductionRelationBasisID identifies a production-relation-basis record.
type ProductionRelationBasisID string

// NewProductionRelationBasisID returns a fresh 96-bit hex ID.
func NewProductionRelationBasisID() ProductionRelationBasisID {
	return ProductionRelationBasisID(newID())
}

// ProductionRelationBasis is an underlying social form (wage-labour, capital,
// landed property) that determines a distribution relation. Under capitalism it is
// always historically specific, not natural or eternal.
type ProductionRelationBasis struct {
	ID                     ProductionRelationBasisID `json:"id"`
	Name                   string                    `json:"name"`
	CharacteristicFeature  string                    `json:"characteristic_feature"`
	IsHistoricallySpecific bool                      `json:"is_historically_specific"`
	CreatedAt              time.Time                 `json:"created_at"`
}

// DistributionRelationID identifies a distribution-relation record.
type DistributionRelationID string

// NewDistributionRelationID returns a fresh 96-bit hex ID.
func NewDistributionRelationID() DistributionRelationID {
	return DistributionRelationID(newID())
}

// DistributionRelation is a specific form (wages/profit/rent) with its production-
// relation basis, the character it appears to have (natural, eternal) versus its real
// character (historically determined, transitory), and the surplus-form it corresponds to.
type DistributionRelation struct {
	ID                       DistributionRelationID `json:"id"`
	Form                     string                 `json:"form"`
	ProductionBasisID        string                 `json:"production_basis_id"`
	ApparentCharacter        string                 `json:"apparent_character"`
	RealCharacter            string                 `json:"real_character"`
	CorrespondsToSurplusForm string                 `json:"corresponds_to_surplus_form"`
	CreatedAt                time.Time              `json:"created_at"`
}

// HistoricalTransienceID identifies a historical-transience record.
type HistoricalTransienceID string

// NewHistoricalTransienceID returns a fresh 96-bit hex ID.
func NewHistoricalTransienceID() HistoricalTransienceID {
	return HistoricalTransienceID(newID())
}

// HistoricalTransience documents that a distribution form is historically specific:
// the mode of production it belongs to, the precondition that established it, and the
// form that succeeds it when those conditions change.
type HistoricalTransience struct {
	ID                 HistoricalTransienceID `json:"id"`
	DistributionFormID string                 `json:"distribution_form_id"`
	ModeOfProduction   string                 `json:"mode_of_production"`
	Precondition       string                 `json:"precondition"`
	SucceedingForm     string                 `json:"succeeding_form"`
	CreatedAt          time.Time              `json:"created_at"`
}

// ProductiveForceContradictionID identifies a productive-force-contradiction record.
type ProductiveForceContradictionID string

// NewProductiveForceContradictionID returns a fresh 96-bit hex ID.
func NewProductiveForceContradictionID() ProductiveForceContradictionID {
	return ProductiveForceContradictionID(newID())
}

// ProductiveForceContradiction records the tension between the developing productive
// forces and the existing social form — the contradiction that signals the historical
// limit of the capitalist mode of production. It is a qualitative record, not a value
// calculation.
type ProductiveForceContradiction struct {
	ID                  ProductiveForceContradictionID `json:"id"`
	ProductiveForce     string                         `json:"productive_force"`
	SocialForm          string                         `json:"social_form"`
	ContradictionNature string                         `json:"contradiction_nature"`
	CrisisSignal        string                         `json:"crisis_signal"`
	CreatedAt           time.Time                      `json:"created_at"`
}
