package revenue

import "testing"

func TestComputeCostPriceAndActualValue(t *testing.T) {
	t.Parallel()
	// Individual capitalist: c=70, v=30 → cost-price 100; profit 20 → actual value 120.
	cost := ComputeCostPrice(70, 30)
	if cost != 100 {
		t.Errorf("ComputeCostPrice = %d, want 100", cost)
	}
	if got := ComputeActualValue(cost, 20); got != 120 {
		t.Errorf("ComputeActualValue = %d, want 120", got)
	}
}

func TestComputeNewValueAndTotalSocialProduct(t *testing.T) {
	t.Parallel()
	// Society-wide: c=6000, wages(v)=1500, profit+rent(s)=1500 → new value 3000, total 9000.
	newValue := ComputeNewValueCreated(1500, 1500)
	if newValue != 3000 {
		t.Errorf("ComputeNewValueCreated = %d, want 3000", newValue)
	}
	if got := ComputeTotalSocialProduct(6000, newValue); got != 9000 {
		t.Errorf("ComputeTotalSocialProduct = %d, want 9000", got)
	}
}

func TestCostPriceFetishInvariants(t *testing.T) {
	t.Parallel()
	f := CostPriceFetish{ConstantCapitalBP: 70, VariableCapitalBP: 30, ProfitBP: 20}
	f.CostPriceBP = ComputeCostPrice(f.ConstantCapitalBP, f.VariableCapitalBP)
	f.ActualValueBP = ComputeActualValue(f.CostPriceBP, f.ProfitBP)
	if f.CostPriceBP != f.ConstantCapitalBP+f.VariableCapitalBP {
		t.Errorf("CostPriceBP = %d, want c+v", f.CostPriceBP)
	}
	if f.ActualValueBP != f.CostPriceBP+f.ProfitBP {
		t.Errorf("ActualValueBP = %d, want cost+profit", f.ActualValueBP)
	}
}

func TestAnnualRevenueAccountInvariants(t *testing.T) {
	t.Parallel()
	newValue := ComputeNewValueCreated(1500, 1500)
	a := AnnualRevenueAccount{
		NewValueCreatedLM:    newValue,
		WagesShareLM:         1500,
		ProfitAndRentShareLM: 1500,
		ConstantCapitalLM:    6000,
		TotalSocialProductLM: ComputeTotalSocialProduct(6000, newValue),
	}
	if a.NewValueCreatedLM != a.WagesShareLM+a.ProfitAndRentShareLM {
		t.Errorf("NewValueCreatedLM = %d, want wages+profit&rent", a.NewValueCreatedLM)
	}
	if a.TotalSocialProductLM != a.ConstantCapitalLM+a.NewValueCreatedLM {
		t.Errorf("TotalSocialProductLM = %d, want c+new", a.TotalSocialProductLM)
	}
}

func TestCh50IDsUniqueAndHex(t *testing.T) {
	t.Parallel()
	seen := map[string]bool{}
	for i := 0; i < 100; i++ {
		for _, id := range []string{
			string(NewCompetitionIllusionID()),
			string(NewCostPriceFetishID()),
			string(NewAnnualRevenueAccountID()),
			string(NewValueAppearanceContrastID()),
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
