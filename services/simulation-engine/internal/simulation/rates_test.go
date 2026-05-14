package simulation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/labour"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

func TestFormulaI_SixOverSix_100Percent(t *testing.T) {
	t.Parallel()
	// §formula I: "6 hours surplus-labour / 6 hours necessary labour = 100%"
	got := simulation.FormulaI(simulation.SurplusLabour{Minutes: 360}, simulation.NecessaryLabour{Minutes: 360})
	if got != 1.0 {
		t.Fatalf("FormulaI: want 1.0, got %f", got)
	}
}

func TestFormulaII_SixOverTwelve_50Percent(t *testing.T) {
	t.Parallel()
	// §formula II: "6 hours surplus-labour / 12-hour working-day = 50%"
	got := simulation.FormulaII(simulation.SurplusLabour{Minutes: 360}, simulation.WorkingDay{TotalMinutes: 720})
	if got != 0.5 {
		t.Fatalf("FormulaII: want 0.5, got %f", got)
	}
}

func TestFormulaII_BoundedBelow1(t *testing.T) {
	t.Parallel()
	// Formula II can never reach or exceed 1.0
	got := simulation.FormulaII(simulation.SurplusLabour{Minutes: 719}, simulation.WorkingDay{TotalMinutes: 720})
	if got >= 1.0 {
		t.Fatalf("FormulaII: want < 1.0, got %f", got)
	}
}

func TestFormulaI_EnglishAgriculturalLabourer_300Percent(t *testing.T) {
	t.Parallel()
	// "the labourer gets only 1/4, the capitalist ... 3/4 ... rate of exploitation of 300%"
	got := simulation.FormulaI(simulation.SurplusLabour{Minutes: 540}, simulation.NecessaryLabour{Minutes: 180})
	if got != 3.0 {
		t.Fatalf("FormulaI (agricultural labourer): want 3.0, got %f", got)
	}
}


func TestComputeRates_Invariants(t *testing.T) {
	t.Parallel()
	sc := simulation.RateScenario{
		NecessaryLabour: simulation.NecessaryLabour{Minutes: 360},
		SurplusLabour:   simulation.SurplusLabour{Minutes: 360},
	}
	r := simulation.ComputeRates(sc)

	// FormulaI == FormulaIII (same ratio, renamed labels)
	if r.FormulaI != r.FormulaIII {
		t.Errorf("FormulaI != FormulaIII: %f vs %f", r.FormulaI, r.FormulaIII)
	}
	// FormulaII < 1.0 whenever surplus > 0
	if r.FormulaII >= 1.0 {
		t.Errorf("FormulaII must be < 1.0, got %f", r.FormulaII)
	}
	// FormulaI > FormulaII whenever necessary < total working day
	if r.FormulaI <= r.FormulaII {
		t.Errorf("FormulaI must exceed FormulaII, got %f vs %f", r.FormulaI, r.FormulaII)
	}
}

func TestComputeRates_WorkingDayDerivedFromPartitions(t *testing.T) {
	t.Parallel()
	sc := simulation.RateScenario{
		NecessaryLabour: simulation.NecessaryLabour{Minutes: labour.LabourMinutes(480)},
		SurplusLabour:   simulation.SurplusLabour{Minutes: labour.LabourMinutes(240)},
	}
	r := simulation.ComputeRates(sc)
	// FormulaII = 240/(480+240) = 240/720 = 1/3
	const want = float64(240) / float64(720)
	if r.FormulaII != want {
		t.Errorf("FormulaII: want %f, got %f", want, r.FormulaII)
	}
}
