package profit

import "testing"

// Capital III, Ch. 1, §2: C = k + p where the profit p equals the surplus-value
// s = £100. The form mystifies the origin of surplus-value because it attributes
// it to the whole capital rather than to variable capital alone.
func TestComputeProfitForm_MystifiesOrigin(t *testing.T) {
	t.Parallel()
	form := ComputeProfitForm(500, 100)
	if form.Profit != 100 {
		t.Errorf("profit = %d, want 100", form.Profit)
	}
	if LabourMinutes(form.SurplusValue) != form.Profit {
		t.Errorf("profit %d != surplus-value %d (p must equal s in this chapter)", form.Profit, form.SurplusValue)
	}
	if form.CommodityValue != 600 {
		t.Errorf("commodity value = %d, want 600", form.CommodityValue)
	}
	if !form.MystifiesOrigin {
		t.Error("the profit-form must mystify the origin of surplus-value")
	}
}

// §1: selling at £510 a commodity worth £600 whose cost-price is £500 still
// realises a £10 profit, even though it is sold £90 below its value.
func TestComputeProfit_BelowValueStillProfits(t *testing.T) {
	t.Parallel()
	p := ComputeProfit(510, 500, 600)
	if p.Amount != 10 {
		t.Errorf("amount = %d, want 10", p.Amount)
	}
	if !p.BelowValue {
		t.Error("510 sits below value 600 but above cost-price 500 — BelowValue should be true")
	}
}

// Sold at value (£600): the profit equals the entire surplus-value (£100) and is
// not flagged as a below-value sale.
func TestComputeProfit_SoldAtValueRealisesFullSurplus(t *testing.T) {
	t.Parallel()
	p := ComputeProfit(600, 500, 600)
	if p.Amount != 100 {
		t.Errorf("amount = %d, want 100", p.Amount)
	}
	if p.BelowValue {
		t.Error("a sale at value is not a below-value sale")
	}
}

// The minimal limit of the selling-price is the cost-price: at k, profit is 0
// and the sale is not below value.
func TestComputeProfit_AtCostPriceNoProfit(t *testing.T) {
	t.Parallel()
	p := ComputeProfit(500, 500, 600)
	if p.Amount != 0 {
		t.Errorf("amount = %d, want 0", p.Amount)
	}
	if p.BelowValue {
		t.Error("selling at the cost-price is not a below-value profit")
	}
}

// Invariant: Amount > 0 whenever SellingPrice > CostPrice, regardless of the
// relation between selling-price and value. Marx's own ladder: 510, 520, 530,
// 560, 590 (below value) and 600 (at value).
func TestComputeProfit_AboveCostAlwaysProfits(t *testing.T) {
	t.Parallel()
	for _, sell := range []LabourMinutes{510, 520, 530, 560, 590, 600, 650} {
		p := ComputeProfit(sell, 500, 600)
		if p.Amount <= 0 {
			t.Errorf("selling at %d above cost-price 500 should yield a profit, got %d", sell, p.Amount)
		}
		if want := sell - 500; p.Amount != want {
			t.Errorf("selling at %d: amount = %d, want %d", sell, p.Amount, want)
		}
	}
}
