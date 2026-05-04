package agent

// ExchangeResult records the outcome of a bilateral exchange between parties A
// and B. The invariant TotalBefore() == TotalAfter() must always hold.
type ExchangeResult struct {
	ABefore Pence  `json:"a_before"`
	BBefore Pence  `json:"b_before"`
	AAfter  Pence  `json:"a_after"`
	BAfter  Pence  `json:"b_after"`
	Origin  string `json:"origin"` // "equivalent" or "redistribution"
}

// TotalBefore returns the combined value of both parties before exchange.
func (r ExchangeResult) TotalBefore() Pence { return r.ABefore + r.BBefore }

// TotalAfter returns the combined value of both parties after exchange.
func (r ExchangeResult) TotalAfter() Pence { return r.AAfter + r.BAfter }

// MerchantsCapital represents M-C-M' operating purely within the sphere of
// circulation. Any surplus-value arises from redistribution, not creation.
type MerchantsCapital struct {
	M           Pence  `json:"m"`
	CommodityID string `json:"commodity_id"`
	MPrime      Pence  `json:"m_prime"`
}

// SurplusValue returns MPrime - M.
func (mc MerchantsCapital) SurplusValue() Pence { return mc.MPrime - mc.M }

// Origin returns "equivalent" when SurplusValue is zero, "redistribution"
// otherwise. It never returns "creation" — circulation cannot create value.
func (mc MerchantsCapital) Origin() string {
	if mc.SurplusValue() == 0 {
		return "equivalent"
	}
	return "redistribution"
}

// UsurersCapital is the degenerate M-M' circuit: money exchanged for more
// money without a commodity intermediary. The source of the surplus cannot be
// located within the circuit.
type UsurersCapital struct {
	M      Pence `json:"m"`
	MPrime Pence `json:"m_prime"`
}

// SurplusValue returns MPrime - M.
func (uc UsurersCapital) SurplusValue() Pence { return uc.MPrime - uc.M }

// Origin returns "equivalent" when SurplusValue is zero, "redistribution" otherwise.
func (uc UsurersCapital) Origin() string {
	if uc.SurplusValue() == 0 {
		return "equivalent"
	}
	return "redistribution"
}

// ExchangeEquivalents models a bilateral swap of equal values: A's commodity
// worth aValue trades for B's commodity worth bValue. Neither party gains
// surplus-value.
func ExchangeEquivalents(aValue, bValue Pence) ExchangeResult {
	return ExchangeResult{
		ABefore: aValue,
		BBefore: bValue,
		AAfter:  bValue,
		BAfter:  aValue,
		Origin:  "equivalent",
	}
}

// ExchangeNonEquivalents models a seller (A) obtaining price above commodity
// value. A gains (price − sellerValue); B loses the same amount. Total social
// value is conserved.
func ExchangeNonEquivalents(sellerValue, price Pence) ExchangeResult {
	return ExchangeResult{
		ABefore: sellerValue,
		BBefore: price,
		AAfter:  price,
		BAfter:  sellerValue,
		Origin:  "redistribution",
	}
}

// TotalValue returns the sum of MoneyBalance across all agents. This is the
// social total that exchange cannot increase.
func TotalValue(agents []Agent) Pence {
	var total Pence
	for _, a := range agents {
		total += a.MoneyBalance
	}
	return total
}
