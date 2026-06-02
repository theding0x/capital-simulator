package reproduction

// from_annual_rates.go — Vol. II Ch. 16 → Ch. 20/21 bridge (issue #222).
//
// Ch. 16 produces an annual rate of surplus-value per capital; Ch. 20/21
// redistribute the resulting surplus across the two departments. This builder
// derives a SimpleReproductionScheme directly from the per-department annual
// rates so the reproduction layer cannot drift from the annual-rate layer: the
// surplus is computed from the rate rather than supplied independently.

// AnnualRateDepartment supplies one department's advanced capital (c, v) plus
// the annual rate of surplus-value (basis points) computed in Ch. 16. The
// AnnualSurplusBP is valorisation.AnnualSurplusRate.AnnualBasisPoints; it is
// taken as an int so this package need not import valorisation.
type AnnualRateDepartment struct {
	Department      CapitalDepartment `json:"department"`
	ConstantPence   int64             `json:"constant_pence"`
	VariablePence   int64             `json:"variable_pence"`
	AnnualSurplusBP int64             `json:"annual_surplus_bp"`
}

// BuildReproductionFromAnnualRates assembles a SimpleReproductionScheme from per
// department annual rates. Each department's surplus is v × AnnualSurplusBP /
// 10000 and its total is c + v + s; departments default to annual turnover
// (TurnoverNumberBP = 10000). IsBalanced reflects the keystone I(v+s) == II(c).
// Duplicate departments are not merged — the last one wins.
func BuildReproductionFromAnnualRates(period string, depts []AnnualRateDepartment) SimpleReproductionScheme {
	s := SimpleReproductionScheme{
		ID:     NewSimpleReproductionSchemeID(),
		Period: period,
	}
	for _, d := range depts {
		surplus := d.VariablePence * d.AnnualSurplusBP / 10000
		dc := DepartmentalCapital{
			ID:               NewDepartmentalCapitalID(),
			SchemeID:         s.ID,
			Department:       d.Department,
			ConstantPence:    d.ConstantPence,
			VariablePence:    d.VariablePence,
			SurplusPence:     surplus,
			TotalPence:       d.ConstantPence + d.VariablePence + surplus,
			TurnoverNumberBP: 10000,
		}
		switch d.Department {
		case DepartmentI:
			cp := dc
			s.DepartmentI = &cp
		case DepartmentII:
			cp := dc
			s.DepartmentII = &cp
		}
	}
	s.IsBalanced = s.CheckBalance()
	return s
}
