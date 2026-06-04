package avgprofit

// This file carries Ch. 9's argument: competition of capitals equalises the
// divergent individual rates of profit (Ch. 8) into a single general rate
// p̄′ = ΣS / ΣC — total social surplus-value divided by total social capital,
// capital-weighted. Each capital then draws average profit = its capital × p̄′,
// regardless of how much surplus-value it itself produced. The commodity no
// longer sells at its value (c+v+s) but at its price of production = k + p̄·p̄′.
//
// Equalisation is zero-sum across society: higher-composition branches (more c,
// less v → produce LESS surplus-value than average) receive more profit than
// they produced → price of production > value (deviation positive). Lower-
// composition branches (more v → produce MORE) are net donors: price < value
// (deviation negative). Conservation: Σ prices = Σ values; Σ avg profits = Σ s.

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// CostPrice is k = c + v consumed in a single period of production.
type CostPrice LabourMinutes

// ProductionPrice is the price of production = k + average profit.
type ProductionPrice LabourMinutes

// CommodityValue is the value of the commodity = c + v + s.
type CommodityValue LabourMinutes

// Profit is a profit magnitude (average profit or actual profit).
type Profit LabourMinutes

// ValueDeviation is signed: ProductionPrice − CommodityValue.
// Positive for higher-composition branches (receivers), negative for lower
// (donors), zero for average.
type ValueDeviation int64

// GeneralProfitRateID identifies a persisted general-rate computation.
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type GeneralProfitRateID string

// NewGeneralProfitRateID returns a fresh random GeneralProfitRateID.
func NewGeneralProfitRateID() GeneralProfitRateID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("avgprofit: rand: %w", err))
	}
	return GeneralProfitRateID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty GeneralProfitRateID.
func (id GeneralProfitRateID) IsZero() bool { return id == "" }

// generalProfitRate computes p̄′ = round_half_up(ΣS·10000 / ΣC).
// Mirrors individualProfitRate: if sumC<=0 returns 0;
// otherwise (sumS*10000 + sumC/2) / sumC.
func generalProfitRate(sumS SurplusValue, sumC TotalCapital) SphereProfitRate {
	if sumC <= 0 {
		return 0
	}
	return SphereProfitRate((int64(sumS)*basisPoints + int64(sumC)/2) / int64(sumC))
}

// GeneralProfitRate is the persisted capital-weighted general rate over a set
// of production spheres. Computed by ComputeGeneralProfitRate; persisted by
// the store.
type GeneralProfitRate struct {
	ID                     GeneralProfitRateID `json:"id"`
	Spheres                []ProductionSphere  `json:"spheres"`
	Rate                   SphereProfitRate    `json:"rate"`                     // p̄′ in bp
	SumSurplusValues       SurplusValue        `json:"sum_surplus_values"`
	SumTotalCapitals       TotalCapital        `json:"sum_total_capitals"`
	AverageVariablePercent int64               `json:"average_variable_percent"` // ΣV·100/ΣC
	CreatedAt              time.Time           `json:"created_at"`
}

// Validate returns an error if any sphere in the set is invalid.
// An empty sphere slice is valid (rate will be 0).
func (g GeneralProfitRate) Validate() error {
	for _, s := range g.Spheres {
		if err := s.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// ComputeGeneralProfitRate sums SurplusValue and TotalCapital over all
// spheres (assumes each sphere is already Computed, i.e. SurplusValue filled),
// derives the capital-weighted general rate, and calculates AverageVariablePercent.
func ComputeGeneralProfitRate(spheres []ProductionSphere) GeneralProfitRate {
	var sumS SurplusValue
	var sumC TotalCapital
	var sumV int64
	for _, s := range spheres {
		sumS += s.SurplusValue
		sumC += s.C
		sumV += int64(s.V)
	}
	var avgVarPct int64
	if sumC > 0 {
		avgVarPct = sumV * 100 / int64(sumC)
	}
	return GeneralProfitRate{
		Spheres:                spheres,
		Rate:                   generalProfitRate(sumS, sumC),
		SumSurplusValues:       sumS,
		SumTotalCapitals:       sumC,
		AverageVariablePercent: avgVarPct,
	}
}

// AverageProfit is the average profit drawn by a capital at the general rate.
// It is pure and not persisted: Amount = capital × rate / 10000.
type AverageProfit struct {
	CapitalAdvanced TotalCapital     `json:"capital_advanced"`
	GeneralRate     SphereProfitRate `json:"general_rate"`
	Amount          Profit           `json:"amount"`
}

// ComputeAverageProfit returns the average profit for a capital advanced at
// the general rate: Amount = c * rate / 10000 (integer truncation).
func ComputeAverageProfit(c TotalCapital, rate SphereProfitRate) AverageProfit {
	return AverageProfit{
		CapitalAdvanced: c,
		GeneralRate:     rate,
		Amount:          Profit(int64(c) * int64(rate) / basisPoints),
	}
}

// SocialCapitalAggregate is a compute-only (not persisted) aggregate over a
// set of production spheres. It derives the weighted general rate and each
// sphere's price of production, and checks the conservation laws.
type SocialCapitalAggregate struct {
	Spheres                 []ProductionSphere  `json:"spheres"`
	WeightedGeneralRate     SphereProfitRate    `json:"weighted_general_rate"`
	SumSurplusValues        SurplusValue        `json:"sum_surplus_values"`
	SumTotalCapitals        TotalCapital        `json:"sum_total_capitals"`
	AverageVariablePercent  int64               `json:"average_variable_percent"`
	PricesOfProduction      []PriceOfProduction `json:"prices_of_production"`
	TotalValues             CommodityValue      `json:"total_values"`
	TotalPricesOfProduction ProductionPrice     `json:"total_prices_of_production"`
	SumAverageProfits       Profit              `json:"sum_average_profits"`
	ValueConserved          bool                `json:"value_conserved"`
	ProfitConserved         bool                `json:"profit_conserved"`
}

// ComputeSocialCapitalAggregate derives the weighted general rate over all
// spheres and computes per-sphere prices of production, then checks both halves
// of L-dM6 conservation: ValueConserved (Σprice = Σvalue) and ProfitConserved
// (Σprofit = Σsurplus-value).
// For each sphere, k = C (= c+v), value = C + s,
// price = ComputePriceOfProduction("", CostPrice(k), weightedRate, CommodityValue(value)).
func ComputeSocialCapitalAggregate(spheres []ProductionSphere) SocialCapitalAggregate {
	g := ComputeGeneralProfitRate(spheres)
	rate := g.Rate

	pops := make([]PriceOfProduction, 0, len(spheres))
	var totalValues CommodityValue
	var totalPrices ProductionPrice
	var sumAvgProfits Profit

	for _, s := range spheres {
		k := CostPrice(s.C)
		value := CommodityValue(int64(s.C) + int64(s.SurplusValue))
		pop := ComputePriceOfProduction("", k, rate, value)
		pops = append(pops, pop)
		totalValues += value
		totalPrices += pop.Price
		sumAvgProfits += pop.AverageProfit
	}

	return SocialCapitalAggregate{
		Spheres:                 spheres,
		WeightedGeneralRate:     rate,
		SumSurplusValues:        g.SumSurplusValues,
		SumTotalCapitals:        g.SumTotalCapitals,
		AverageVariablePercent:  g.AverageVariablePercent,
		PricesOfProduction:      pops,
		TotalValues:             totalValues,
		TotalPricesOfProduction: totalPrices,
		SumAverageProfits:       sumAvgProfits,
		ValueConserved:          totalValues == CommodityValue(totalPrices),
		ProfitConserved:         sumAvgProfits == Profit(g.SumSurplusValues),
	}
}
