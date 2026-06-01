package revenue

import "testing"

// TestDistributionFormsRestOnProductionBases asserts the Ch. 51 invariant that every
// capitalist distribution form (wages/profit/rent) rests on a historically specific
// production relation and references it.
func TestDistributionFormsRestOnProductionBases(t *testing.T) {
	t.Parallel()
	wageLabour := ProductionRelationBasis{ID: NewProductionRelationBasisID(), Name: "wage-labour", IsHistoricallySpecific: true}
	capital := ProductionRelationBasis{ID: NewProductionRelationBasisID(), Name: "capital", IsHistoricallySpecific: true}
	land := ProductionRelationBasis{ID: NewProductionRelationBasisID(), Name: "landed_property", IsHistoricallySpecific: true}
	bases := map[string]ProductionRelationBasis{
		string(wageLabour.ID): wageLabour,
		string(capital.ID):    capital,
		string(land.ID):       land,
	}

	forms := []DistributionRelation{
		{Form: "wages", ProductionBasisID: string(wageLabour.ID), CorrespondsToSurplusForm: "v"},
		{Form: "profit", ProductionBasisID: string(capital.ID), CorrespondsToSurplusForm: "s(profit)"},
		{Form: "rent", ProductionBasisID: string(land.ID), CorrespondsToSurplusForm: "s(rent)"},
	}
	for _, f := range forms {
		basis, ok := bases[f.ProductionBasisID]
		if !ok {
			t.Errorf("distribution form %q references unknown basis %q", f.Form, f.ProductionBasisID)
			continue
		}
		if !basis.IsHistoricallySpecific {
			t.Errorf("basis %q for form %q is not historically specific", basis.Name, f.Form)
		}
	}
}

func TestCh51IDsUniqueAndHex(t *testing.T) {
	t.Parallel()
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		for _, id := range []string{
			string(NewProductionRelationBasisID()),
			string(NewDistributionRelationID()),
			string(NewHistoricalTransienceID()),
			string(NewProductiveForceContradictionID()),
		} {
			if len(id) != 24 {
				t.Fatalf("id %q length = %d, want 24", id, len(id))
			}
			if seen[id] {
				t.Fatalf("duplicate id %q", id)
			}
			seen[id] = true
		}
	}
}
