package profit

// This file completes Ch. 1's second movement: surplus-value's appearance as
// profit. Two related forms:
//
//   - ProfitForm regroups the value formula C = k + s as C = k + p, where the
//     profit p equals the surplus-value s but is now attributed to the total
//     capital advanced rather than to variable capital alone. That re-attribution
//     is the source of the "mystification": with the cost-price erasing the
//     distinction between constant and variable capital, surplus-value appears
//     to spring from the whole outlay.
//   - Profit is the realised excess of selling-price over cost-price. Because
//     the cost-price (k) is below the value (C), a commodity can sell below its
//     value and still yield a profit — the cornerstone of the later theory of
//     prices of production.

// ProfitForm expresses the transformation C = k + s → C = k + p.
//
// Profit equals SurplusValue in magnitude; the two differ only in their
// conceptual attribution. MystifiesOrigin is always true: the form conceals
// that surplus-value originates solely in variable capital.
type ProfitForm struct {
	CostPrice       LabourMinutes `json:"cost_price"`
	SurplusValue    SurplusValue  `json:"surplus_value"`
	Profit          LabourMinutes `json:"profit"`
	CommodityValue  LabourMinutes `json:"commodity_value"`
	MystifiesOrigin bool          `json:"mystifies_origin"`
}

// ComputeProfitForm builds C = k + p from a cost-price k and surplus-value s.
// p = s by magnitude, and the commodity-value C = k + p.
func ComputeProfitForm(k LabourMinutes, s SurplusValue) ProfitForm {
	return ProfitForm{
		CostPrice:       k,
		SurplusValue:    s,
		Profit:          LabourMinutes(s),
		CommodityValue:  k + LabourMinutes(s),
		MystifiesOrigin: true,
	}
}

// Profit is the profit realised when a commodity is sold at a given price.
//
// Amount is SellingPrice − CostPrice. BelowValue records the case Marx
// stresses: a selling-price above the cost-price but below the value still
// realises a profit (a portion of the surplus-value), even though the commodity
// is sold under its value.
type Profit struct {
	SellingPrice   LabourMinutes `json:"selling_price"`
	CostPrice      LabourMinutes `json:"cost_price"`
	CommodityValue LabourMinutes `json:"commodity_value"`
	Amount         LabourMinutes `json:"amount"`
	BelowValue     bool          `json:"below_value"`
}

// ComputeProfit computes the profit realised by selling at sellingPrice a
// commodity whose cost-price is costPrice and whose value is commodityValue.
//
// Profit is realised on any sale above the cost-price (the minimal limit of the
// selling-price). When the selling-price also sits below the commodity's value,
// BelowValue is true: the capitalist has sold under value yet still pocketed a
// profit — only a fraction of the surplus-value, but a profit nonetheless.
func ComputeProfit(sellingPrice, costPrice, commodityValue LabourMinutes) Profit {
	return Profit{
		SellingPrice:   sellingPrice,
		CostPrice:      costPrice,
		CommodityValue: commodityValue,
		Amount:         sellingPrice - costPrice,
		BelowValue:     sellingPrice > costPrice && sellingPrice < commodityValue,
	}
}
