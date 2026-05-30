package avgprofit

// This file carries Ch. 11's argument: a general rise in wages does NOT raise
// all prices of production — it lowers the general rate of profit and shifts
// prices in opposite directions by organic composition.
//
//   - average-composition: price unchanged.
//   - lower-composition (more v): price rises.
//   - higher-composition (more c): price falls.
//
// A wage fall reverses all three directions and raises the general rate.
// Living labour (v+s) per sphere is fixed; a wage change only shifts the v/s
// split, leaving total value unchanged — so Σ price changes = 0.
//
// Units: all price/cost magnitudes are integer centi-units (×100). For example
// £120.00 is stored as 12000. Composition inputs c,v remain whole units.
// The web panel divides price fields by 100 for display.

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// WageFluctuationKind indicates whether wages have risen, fallen, or are unchanged.
type WageFluctuationKind string

const (
	KindWageRise WageFluctuationKind = "rise"
	KindWageFall WageFluctuationKind = "fall"
)

// WageOutcomeDirection encodes how the price of production moves for a sphere.
type WageOutcomeDirection string

const (
	KindPriceRises     WageOutcomeDirection = "rises"
	KindPriceFalls     WageOutcomeDirection = "falls"
	KindPriceUnchanged WageOutcomeDirection = "unchanged"
)

// WageFluctuation is a value object encoding the direction of a wage change.
type WageFluctuation struct {
	Factor int64               `json:"factor"` // percent: 100 = no change, 125 = +25%, 75 = -25%
	Kind   WageFluctuationKind `json:"kind"`   // "rise" | "fall" | ""
}

// ComputeWageFluctuation derives a WageFluctuation from a percent factor.
// factor > 100 -> KindWageRise; factor < 100 -> KindWageFall; == 100 -> "".
func ComputeWageFluctuation(factor int64) WageFluctuation {
	var kind WageFluctuationKind
	switch {
	case factor > 100:
		kind = KindWageRise
	case factor < 100:
		kind = KindWageFall
	}
	return WageFluctuation{Factor: factor, Kind: kind}
}

// SphereWageOutcome records the price-of-production shift for one sphere under
// a given wage change. OldPrice, NewPrice, and PriceChange are in centi-units.
type SphereWageOutcome struct {
	SphereName      string               `json:"sphere_name"`
	ConstantCapital ConstantCapital      `json:"constant_capital"` // c, whole units
	VariableCapital VariableCapital      `json:"variable_capital"` // v, whole units
	WageFactor      int64                `json:"wage_factor"`
	Composition     CompositionKind      `json:"composition"`  // vs base: higher/average/lower
	OldPrice        ProductionPrice      `json:"old_price"`    // centi-units, at old general rate, factor 100
	NewPrice        ProductionPrice      `json:"new_price"`    // centi-units, at new general rate, wage factor
	PriceChange     ProductionPrice      `json:"price_change"` // NewPrice - OldPrice (signed)
	Direction       WageOutcomeDirection `json:"direction"`
}

// WageEffectAnalysisID identifies a persisted wage-effect analysis.
// 96-bit hex from crypto/rand.
type WageEffectAnalysisID string

// NewWageEffectAnalysisID returns a fresh random WageEffectAnalysisID.
func NewWageEffectAnalysisID() WageEffectAnalysisID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("avgprofit: rand: %w", err))
	}
	return WageEffectAnalysisID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty WageEffectAnalysisID.
func (id WageEffectAnalysisID) IsZero() bool { return id == "" }

// WageEffectAnalysis is the persisted entity: a full analysis of how a general
// wage fluctuation shifts prices of production across production spheres.
type WageEffectAnalysis struct {
	ID             WageEffectAnalysisID `json:"id"`
	BaseConstant   ConstantCapital      `json:"base_constant"`    // average capital c
	BaseVariable   VariableCapital      `json:"base_variable"`    // average capital v
	SRate          RateOfSurplusValue   `json:"s_rate"`           // base s'
	WageFactor     int64                `json:"wage_factor"`
	Kind           WageFluctuationKind  `json:"kind"`
	OldGeneralRate SphereProfitRate     `json:"old_general_rate"`
	NewGeneralRate SphereProfitRate     `json:"new_general_rate"`
	AverageOutcome SphereWageOutcome    `json:"average_outcome"` // base sphere's outcome (always unchanged)
	Outcomes       []SphereWageOutcome  `json:"outcomes"`        // all provided spheres
	CreatedAt      time.Time            `json:"created_at"`
}

// Validate returns ErrNegativeMagnitude if any input magnitude is negative.
func (w WageEffectAnalysis) Validate() error {
	if w.BaseConstant < 0 || w.BaseVariable < 0 || w.SRate < 0 || w.WageFactor < 0 {
		return ErrNegativeMagnitude
	}
	return nil
}

// SphereInput is a helper input type for ComputeWageEffectAnalysis.
// It is not persisted; only the computed WageEffectAnalysis is stored.
type SphereInput struct {
	Name string
	C    ConstantCapital
	V    VariableCapital
}

// sphereWagePrice computes the price of production in centi-units.
//
//	k_centi = c*100 + v*factor
//	price   = k_centi*(10000+rate)/10000  (TRUNCATING integer division)
//
// For the old price: pass factor=100 and rate=oldGeneralRate.
// For the new price: pass factor=wageFactor and rate=newGeneralRate.
func sphereWagePrice(c ConstantCapital, v VariableCapital, factor int64, rate SphereProfitRate) ProductionPrice {
	kCenti := int64(c)*100 + int64(v)*factor
	price := kCenti * (basisPoints + int64(rate)) / basisPoints
	return ProductionPrice(price)
}

// wageAdjustedGeneralRate computes the new general rate (round-half-up) after a
// wage change on the base (average) capital:
//
//	livingLabour = v + v*sRate/10000   (v+s; constant across wage changes)
//	newV         = v*factor/100
//	newS         = livingLabour - newV
//	newC         = c + newV            (total capital at new wage)
//	rate         = round_half_up(newS*10000 / newC)
//
// Returns 0 if newC <= 0.
func wageAdjustedGeneralRate(c ConstantCapital, v VariableCapital, sRate RateOfSurplusValue, factor int64) SphereProfitRate {
	livingLabour := int64(v) + int64(v)*int64(sRate)/basisPoints
	newV := int64(v) * factor / 100
	newS := livingLabour - newV
	newC := int64(c) + newV
	if newC <= 0 {
		return 0
	}
	return SphereProfitRate((newS*basisPoints + newC/2) / newC)
}

// baseGeneralRate computes the old general rate (round-half-up):
//
//	surplusValue = v * sRate / 10000
//	rate         = round_half_up(surplusValue*10000 / (c+v))
func baseGeneralRate(c ConstantCapital, v VariableCapital, sRate RateOfSurplusValue) SphereProfitRate {
	total := int64(c) + int64(v)
	if total <= 0 {
		return 0
	}
	sv := int64(v) * int64(sRate) / basisPoints
	return SphereProfitRate((sv*basisPoints + total/2) / total)
}

// ComputeSphereWageOutcome computes the price-of-production shift for one sphere.
//
//	OldPrice    = sphereWagePrice(c, v, 100, oldRate)
//	NewPrice    = sphereWagePrice(c, v, factor, newRate)
//	PriceChange = NewPrice - OldPrice
//	Direction   : >0 -> rises, <0 -> falls, ==0 -> unchanged
//	Composition : varPercent = v*100/(c+v); < base -> KindHigher; > -> KindLower; == -> KindAverage
func ComputeSphereWageOutcome(name string, c ConstantCapital, v VariableCapital, factor int64, oldRate, newRate SphereProfitRate, baseVarPercent int64) SphereWageOutcome {
	oldPrice := sphereWagePrice(c, v, 100, oldRate)
	newPrice := sphereWagePrice(c, v, factor, newRate)
	change := ProductionPrice(int64(newPrice) - int64(oldPrice))

	var dir WageOutcomeDirection
	switch {
	case change > 0:
		dir = KindPriceRises
	case change < 0:
		dir = KindPriceFalls
	default:
		dir = KindPriceUnchanged
	}

	var comp CompositionKind
	total := int64(c) + int64(v)
	if total > 0 {
		varPct := int64(v) * 100 / total
		switch {
		case varPct < baseVarPercent:
			comp = KindHigher
		case varPct > baseVarPercent:
			comp = KindLower
		default:
			comp = KindAverage
		}
	}

	return SphereWageOutcome{
		SphereName:      name,
		ConstantCapital: c,
		VariableCapital: v,
		WageFactor:      factor,
		Composition:     comp,
		OldPrice:        oldPrice,
		NewPrice:        newPrice,
		PriceChange:     change,
		Direction:       dir,
	}
}

// ComputeWageEffectAnalysis computes the full wage-effect analysis.
// baseC and baseV define the average-composition capital; sRate is the base s'.
// factor is the wage percent (100 = no change, 125 = +25%). spheres is the list
// of additional production spheres to compare.
func ComputeWageEffectAnalysis(baseC ConstantCapital, baseV VariableCapital, sRate RateOfSurplusValue, factor int64, spheres []SphereInput) WageEffectAnalysis {
	oldRate := baseGeneralRate(baseC, baseV, sRate)
	newRate := wageAdjustedGeneralRate(baseC, baseV, sRate, factor)
	fluctuation := ComputeWageFluctuation(factor)

	baseVarPercent := int64(0)
	if total := int64(baseC) + int64(baseV); total > 0 {
		baseVarPercent = int64(baseV) * 100 / total
	}

	avgOutcome := ComputeSphereWageOutcome("average", baseC, baseV, factor, oldRate, newRate, baseVarPercent)

	outcomes := make([]SphereWageOutcome, 0, len(spheres))
	for _, s := range spheres {
		outcomes = append(outcomes, ComputeSphereWageOutcome(s.Name, s.C, s.V, factor, oldRate, newRate, baseVarPercent))
	}

	return WageEffectAnalysis{
		BaseConstant:   baseC,
		BaseVariable:   baseV,
		SRate:          sRate,
		WageFactor:     factor,
		Kind:           fluctuation.Kind,
		OldGeneralRate: oldRate,
		NewGeneralRate: newRate,
		AverageOutcome: avgOutcome,
		Outcomes:       outcomes,
	}
}
