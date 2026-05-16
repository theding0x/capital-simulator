package simulation

import (
	"crypto/rand"
	"encoding/hex"
	"math"
	"time"
)

// ValueComposition is the monetary ratio of constant to variable capital (c:v).
type ValueComposition struct {
	ConstantCapital Pence `json:"constant_capital"`
	VariableCapital Pence `json:"variable_capital"`
}

// TechnicalComposition is the material ratio of means of production to labour-power.
type TechnicalComposition struct {
	MeansOfProductionUnits int64 `json:"means_of_production_units"`
	LabourPowerUnits        int64 `json:"labour_power_units"`
}

// OrganicComposition is the value composition determined by the technical
// composition: Ratio = c/(c+v), always in [0,1).
type OrganicComposition struct {
	Ratio float64 `json:"ratio"`
}

// Concentration records the average capital per firm as accumulation proceeds.
type Concentration struct {
	TotalCapital Pence `json:"total_capital"`
	Firms        int64 `json:"firms"`
}

// AverageCapitalPerFirm returns TotalCapital / Firms, or 0 if Firms is zero.
func (c Concentration) AverageCapitalPerFirm() Pence {
	if c.Firms <= 0 {
		return 0
	}
	return Pence(int64(c.TotalCapital) / c.Firms)
}

// Centralisation records capital absorbed by merger or takeover rather than by
// surplus-value accumulation.
type Centralisation struct {
	AcquiredCapital Pence `json:"acquired_capital"`
	FirmsAbsorbed   int64 `json:"firms_absorbed"`
}

// RelativeSurplusPopulation is the reserve army divided into its three strata:
// floating (factory workers displaced by machinery), latent (agricultural
// workers ready to enter industry), and stagnant (irregular casual labour).
type RelativeSurplusPopulation struct {
	Floating int64 `json:"floating"`
	Latent   int64 `json:"latent"`
	Stagnant int64 `json:"stagnant"`
}

// Total returns the sum of the three strata.
func (r RelativeSurplusPopulation) Total() int64 {
	return r.Floating + r.Latent + r.Stagnant
}

// IndustrialReserveArmy is the aggregate reserve of labour and its proportion
// relative to the active army.
type IndustrialReserveArmy struct {
	Size               int64   `json:"size"`
	RelativeProportion float64 `json:"relative_proportion"`
}

// LabourDemand is the number of workers demanded by capital at a given
// composition and wage.
type LabourDemand struct {
	Workers int64 `json:"workers"`
}

// ComputeOrganicComposition derives the organic composition from the value
// composition. Ratio = c/(c+v); always in [0,1).
func ComputeOrganicComposition(vc ValueComposition) OrganicComposition {
	total := float64(vc.ConstantCapital) + float64(vc.VariableCapital)
	if total <= 0 {
		return OrganicComposition{Ratio: 0}
	}
	return OrganicComposition{Ratio: float64(vc.ConstantCapital) / total}
}

// ComputeLabourDemand returns the number of workers demanded given total
// capital, organic composition, and the wage per worker in pence.
// Workers = total_capital * (1 - OC.Ratio) / wagePence.
func ComputeLabourDemand(totalCapital Pence, oc OrganicComposition, wagePence int64) LabourDemand {
	if wagePence <= 0 || totalCapital <= 0 {
		return LabourDemand{Workers: 0}
	}
	variablePortion := float64(totalCapital) * (1 - oc.Ratio)
	return LabourDemand{Workers: int64(math.Round(variablePortion / float64(wagePence)))}
}

// ComputeReserveArmy returns the industrial reserve army given total labour
// supply and current demand. Size is clamped to zero from below — workers do
// not disappear into negative employment.
func ComputeReserveArmy(supply int64, demand LabourDemand) IndustrialReserveArmy {
	size := supply - demand.Workers
	if size < 0 {
		size = 0
	}
	var proportion float64
	if demand.Workers > 0 {
		proportion = float64(size) / float64(demand.Workers)
	}
	return IndustrialReserveArmy{Size: size, RelativeProportion: proportion}
}

// GeneralLawScenarioID is a 96-bit hex identifier for a persisted scenario.
type GeneralLawScenarioID string

// IsZero reports whether the ID is unset.
func (s GeneralLawScenarioID) IsZero() bool { return s == "" }

// NewGeneralLawScenarioID generates a 96-bit hex identifier via crypto/rand.
func NewGeneralLawScenarioID() GeneralLawScenarioID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return GeneralLawScenarioID(hex.EncodeToString(b))
}

// GeneralLawScenario is a persisted named simulation of the general law:
// a starting capital composition that accumulates over time while the organic
// composition rises due to productivity growth, depressing labour demand and
// expanding the industrial reserve army.
type GeneralLawScenario struct {
	ID                 GeneralLawScenarioID `json:"id"`
	Name               string               `json:"name"`
	ConstantCapital    Pence                `json:"constant_capital"`
	VariableCapital    Pence                `json:"variable_capital"`
	SurplusRate        float64              `json:"surplus_rate"`
	AccumulationRate   float64              `json:"accumulation_rate"`
	ProductivityGrowth float64              `json:"productivity_growth"`
	WagePence          int64                `json:"wage_pence"`
	WorkerSupply       int64                `json:"worker_supply"`
	Periods            int64                `json:"periods"`
	CreatedAt          time.Time            `json:"created_at"`
}

// GeneralLawSnapshot is the state of one period of a general-law simulation.
type GeneralLawSnapshot struct {
	Period             int64   `json:"period"`
	ConstantCapital    Pence   `json:"constant_capital"`
	VariableCapital    Pence   `json:"variable_capital"`
	OrganicComposition float64 `json:"organic_composition"`
	Workers            int64   `json:"workers"`
	ReserveArmySize    int64   `json:"reserve_army_size"`
	RelativeProportion float64 `json:"relative_proportion"`
}

// RunGeneralLaw simulates the general law of capitalist accumulation.
// Each period: the organic composition is computed from the current stock;
// labour demand and the reserve army are derived from it; then surplus-value
// is produced and reinvested. The productivity_growth field raises the
// composition ratio used when splitting the next period's surplus, so that
// new constant capital (machinery) is added at an ever-higher proportion —
// modelling how mechanisation steadily displaces variable capital.
func RunGeneralLaw(s GeneralLawScenario) []GeneralLawSnapshot {
	if s.Periods <= 0 {
		return []GeneralLawSnapshot{}
	}
	result := make([]GeneralLawSnapshot, s.Periods)
	c := s.ConstantCapital
	v := s.VariableCapital
	for i := int64(0); i < s.Periods; i++ {
		vc := ValueComposition{ConstantCapital: c, VariableCapital: v}
		oc := ComputeOrganicComposition(vc)
		total := c + v
		demand := ComputeLabourDemand(total, oc, s.WagePence)
		army := ComputeReserveArmy(s.WorkerSupply, demand)

		result[i] = GeneralLawSnapshot{
			Period:             i + 1,
			ConstantCapital:    c,
			VariableCapital:    v,
			OrganicComposition: oc.Ratio,
			Workers:            demand.Workers,
			ReserveArmySize:    army.Size,
			RelativeProportion: army.RelativeProportion,
		}

		// New capital's composition ratio rises by ProductivityGrowth each period,
		// modelling the displacement of labour by machinery.
		sv := Pence(math.Round(float64(v) * s.SurplusRate))
		newOCRatio := oc.Ratio + s.ProductivityGrowth
		if newOCRatio >= 1.0 {
			newOCRatio = 0.99
		}
		additional := SplitSurplus(sv, s.AccumulationRate, newOCRatio)
		c += additional.Constant
		v += additional.Variable
	}
	return result
}
