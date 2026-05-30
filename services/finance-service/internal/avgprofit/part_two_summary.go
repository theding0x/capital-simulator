package avgprofit

// This file carries Ch. 12's third remark, which closes Part II: across the
// whole social capital, total price of production equals total value, and for
// the average-composition capital the surplus-value produced equals the average
// profit received (the conservation law established in Ch. 9). The summary is a
// read-only snapshot aggregated from the spheres, prices of production, and
// wage-effect analyses recorded across Part II. It is computed on demand and is
// NOT persisted.

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// PartIISummaryID identifies a Part II summary snapshot.
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type PartIISummaryID string

// NewPartIISummaryID returns a fresh random PartIISummaryID.
func NewPartIISummaryID() PartIISummaryID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("avgprofit: rand: %w", err))
	}
	return PartIISummaryID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty PartIISummaryID.
func (id PartIISummaryID) IsZero() bool { return id == "" }

// AverageCompositionCheck records the equality, for the social aggregate, of the
// surplus-value produced and the average profit drawn at the general rate. When
// value is conserved (Σ prices == Σ values), the two are equal.
type AverageCompositionCheck struct {
	TotalSurplusValue               SurplusValue `json:"total_surplus_value"`
	TotalAverageProfit              Profit       `json:"total_average_profit"`
	SurplusValueEqualsAverageProfit bool         `json:"surplus_value_equals_average_profit"`
}

// PartIISummary is the read-only Part II snapshot: the general rate, the counts
// of the records gathered, the conservation of value (Σ prices == Σ values), and
// the average-composition check. Computed by ComputePartIISummary.
type PartIISummary struct {
	ID                      PartIISummaryID         `json:"id"`
	GeneralRate             SphereProfitRate        `json:"general_rate"`
	SphereCount             int                     `json:"sphere_count"`
	PricesOfProductionCount int                     `json:"prices_of_production_count"`
	WageEffectCount         int                     `json:"wage_effect_count"`
	TotalValues             CommodityValue          `json:"total_values"`
	TotalPricesOfProduction ProductionPrice         `json:"total_prices_of_production"`
	ValueConserved          bool                    `json:"value_conserved"`
	AverageCompositionCheck AverageCompositionCheck `json:"average_composition_check"`
	CreatedAt               time.Time               `json:"created_at"`
}

// ComputePartIISummary aggregates the Part II records into a conservation
// snapshot:
//
//	TotalValues             = Σ prices[i].CommodityValue
//	TotalPricesOfProduction = Σ prices[i].Price
//	ValueConserved          = TotalValues == TotalPricesOfProduction
//	GeneralRate             = prices[0].GeneralRate if any price exists,
//	                          else ComputeGeneralProfitRate(spheres).Rate
//	TotalSurplusValue       = Σ spheres[i].SurplusValue
//	TotalAverageProfit      = Σ ( spheres[i].C * GeneralRate / 10000 )
//	SurplusValueEqualsAverageProfit = ValueConserved
//
// It assigns no ID or timestamp — the handler/store do that as needed.
func ComputePartIISummary(spheres []ProductionSphere, prices []PriceOfProduction, wageAnalyses []WageEffectAnalysis) PartIISummary {
	var totalValues CommodityValue
	var totalPrices ProductionPrice
	for _, p := range prices {
		totalValues += p.CommodityValue
		totalPrices += p.Price
	}
	valueConserved := int64(totalValues) == int64(totalPrices)

	var generalRate SphereProfitRate
	if len(prices) > 0 {
		generalRate = prices[0].GeneralRate
	} else {
		generalRate = ComputeGeneralProfitRate(spheres).Rate
	}

	var totalSurplus SurplusValue
	var totalAvgProfit Profit
	for _, s := range spheres {
		totalSurplus += s.SurplusValue
		totalAvgProfit += Profit(int64(s.C) * int64(generalRate) / basisPoints)
	}

	return PartIISummary{
		GeneralRate:             generalRate,
		SphereCount:             len(spheres),
		PricesOfProductionCount: len(prices),
		WageEffectCount:         len(wageAnalyses),
		TotalValues:             totalValues,
		TotalPricesOfProduction: totalPrices,
		ValueConserved:          valueConserved,
		AverageCompositionCheck: AverageCompositionCheck{
			TotalSurplusValue:               totalSurplus,
			TotalAverageProfit:              totalAvgProfit,
			SurplusValueEqualsAverageProfit: valueConserved,
		},
	}
}
