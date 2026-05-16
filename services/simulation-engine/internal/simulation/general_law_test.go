package simulation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// §2: ComputeOrganicComposition(ValueComposition{8000,2000}).Ratio == 0.8
func TestComputeOrganicComposition_InitialComposition(t *testing.T) {
	t.Parallel()
	vc := simulation.ValueComposition{ConstantCapital: 8000, VariableCapital: 2000}
	got := simulation.ComputeOrganicComposition(vc)
	if got.Ratio != 0.8 {
		t.Errorf("ratio = %v, want 0.8", got.Ratio)
	}
}

// §2 rising OC: ComputeOrganicComposition(ValueComposition{9000,1000}).Ratio == 0.9 > 0.8
func TestComputeOrganicComposition_RisingComposition(t *testing.T) {
	t.Parallel()
	low := simulation.ComputeOrganicComposition(simulation.ValueComposition{ConstantCapital: 8000, VariableCapital: 2000})
	high := simulation.ComputeOrganicComposition(simulation.ValueComposition{ConstantCapital: 9000, VariableCapital: 1000})
	if high.Ratio != 0.9 {
		t.Errorf("ratio = %v, want 0.9", high.Ratio)
	}
	if high.Ratio <= low.Ratio {
		t.Errorf("rising OC invariant violated: high %v <= low %v", high.Ratio, low.Ratio)
	}
}

// Zero total capital returns Ratio=0, no division by zero.
func TestComputeOrganicComposition_ZeroTotal(t *testing.T) {
	t.Parallel()
	got := simulation.ComputeOrganicComposition(simulation.ValueComposition{})
	if got.Ratio != 0 {
		t.Errorf("ratio = %v, want 0", got.Ratio)
	}
}

// §1 unchanged composition: doubling total capital doubles labour demand.
func TestComputeLabourDemand_DoublingCapitalDoublesWorkers(t *testing.T) {
	t.Parallel()
	vc := simulation.ValueComposition{ConstantCapital: 8000, VariableCapital: 2000}
	oc := simulation.ComputeOrganicComposition(vc)
	wagePence := int64(200)

	demandBefore := simulation.ComputeLabourDemand(10000, oc, wagePence)
	demandAfter := simulation.ComputeLabourDemand(20000, oc, wagePence)

	if demandAfter.Workers != 2*demandBefore.Workers {
		t.Errorf("doubling capital: workers = %d, want %d", demandAfter.Workers, 2*demandBefore.Workers)
	}
}

// §2: with the same total capital, higher OC means fewer workers demanded.
func TestComputeLabourDemand_RisingOCReducesDemand(t *testing.T) {
	t.Parallel()
	total := simulation.Pence(10000)
	wagePence := int64(200)

	lowOC := simulation.ComputeOrganicComposition(simulation.ValueComposition{ConstantCapital: 8000, VariableCapital: 2000})
	highOC := simulation.ComputeOrganicComposition(simulation.ValueComposition{ConstantCapital: 9000, VariableCapital: 1000})

	demandLow := simulation.ComputeLabourDemand(total, lowOC, wagePence)
	demandHigh := simulation.ComputeLabourDemand(total, highOC, wagePence)

	if demandHigh.Workers >= demandLow.Workers {
		t.Errorf("rising OC should reduce demand: high=%d low=%d", demandHigh.Workers, demandLow.Workers)
	}
}

// Zero wage or zero capital returns zero workers demanded.
func TestComputeLabourDemand_ZeroInputsReturnZero(t *testing.T) {
	t.Parallel()
	oc := simulation.OrganicComposition{Ratio: 0.8}
	if simulation.ComputeLabourDemand(0, oc, 200).Workers != 0 {
		t.Error("zero capital should return zero workers")
	}
	if simulation.ComputeLabourDemand(10000, oc, 0).Workers != 0 {
		t.Error("zero wage should return zero workers")
	}
}

// §3: reserve army grows as capital accumulates with rising OC.
func TestComputeReserveArmy_GrowsWithRisingOC(t *testing.T) {
	t.Parallel()
	supply := int64(50)
	wagePence := int64(200)

	oc1 := simulation.ComputeOrganicComposition(simulation.ValueComposition{ConstantCapital: 8000, VariableCapital: 2000})
	demand1 := simulation.ComputeLabourDemand(10000, oc1, wagePence)
	army1 := simulation.ComputeReserveArmy(supply, demand1)

	oc2 := simulation.ComputeOrganicComposition(simulation.ValueComposition{ConstantCapital: 9000, VariableCapital: 1000})
	demand2 := simulation.ComputeLabourDemand(12000, oc2, wagePence)
	army2 := simulation.ComputeReserveArmy(supply, demand2)

	if army2.Size <= army1.Size {
		t.Errorf("reserve army should grow with rising OC: period1=%d period2=%d", army1.Size, army2.Size)
	}
}

// Demand exceeding supply clamps reserve army to zero.
func TestComputeReserveArmy_ZeroFloor(t *testing.T) {
	t.Parallel()
	supply := int64(5)
	demand := simulation.LabourDemand{Workers: 100}
	army := simulation.ComputeReserveArmy(supply, demand)
	if army.Size != 0 {
		t.Errorf("size = %d, want 0 when demand exceeds supply", army.Size)
	}
}

// Invariant: RelativeProportion == Size / demand.Workers.
func TestComputeReserveArmy_RelativeProportion(t *testing.T) {
	t.Parallel()
	supply := int64(50)
	demand := simulation.LabourDemand{Workers: 40}
	army := simulation.ComputeReserveArmy(supply, demand)
	if army.Size != 10 {
		t.Errorf("size = %d, want 10", army.Size)
	}
	want := float64(10) / float64(40)
	if army.RelativeProportion != want {
		t.Errorf("proportion = %v, want %v", army.RelativeProportion, want)
	}
}

// §1 fixture: unchanged composition — OC stays flat across all periods.
func TestRunGeneralLaw_UnchangedComposition(t *testing.T) {
	t.Parallel()
	s := simulation.GeneralLawScenario{
		ConstantCapital:    8000,
		VariableCapital:    2000,
		SurplusRate:        1.0,
		AccumulationRate:   1.0,
		ProductivityGrowth: 0.0,
		WagePence:          200,
		WorkerSupply:       50,
		Periods:            3,
	}
	snaps := simulation.RunGeneralLaw(s)
	if len(snaps) != 3 {
		t.Fatalf("len = %d, want 3", len(snaps))
	}
	firstOC := snaps[0].OrganicComposition
	for i, sn := range snaps {
		if sn.OrganicComposition != firstOC {
			t.Errorf("period %d: OC = %v, want %v (unchanged composition)", i+1, sn.OrganicComposition, firstOC)
		}
	}
}

// §2 fixture: rising composition — OC strictly increases each period.
func TestRunGeneralLaw_RisingComposition(t *testing.T) {
	t.Parallel()
	s := simulation.GeneralLawScenario{
		ConstantCapital:    8000,
		VariableCapital:    2000,
		SurplusRate:        1.0,
		AccumulationRate:   1.0,
		ProductivityGrowth: 0.05,
		WagePence:          200,
		WorkerSupply:       50,
		Periods:            5,
	}
	snaps := simulation.RunGeneralLaw(s)
	if len(snaps) != 5 {
		t.Fatalf("len = %d, want 5", len(snaps))
	}
	for i := 1; i < len(snaps); i++ {
		if snaps[i].OrganicComposition <= snaps[i-1].OrganicComposition {
			t.Errorf("period %d: OC %v should exceed period %d OC %v",
				i+1, snaps[i].OrganicComposition, i, snaps[i-1].OrganicComposition)
		}
	}
}

// §3: with rising OC, the worker-to-capital ratio falls — labour demand grows
// slower than total capital. This is the measurable content of the general law
// in a model where V always grows (newOCRatio < 1) so absolute army size shrinks
// as accumulated capital absorbs workers from the pre-existing reserve.
func TestRunGeneralLaw_LabourShareFalls(t *testing.T) {
	t.Parallel()
	s := simulation.GeneralLawScenario{
		ConstantCapital:    8000,
		VariableCapital:    2000,
		SurplusRate:        1.0,
		AccumulationRate:   1.0,
		ProductivityGrowth: 0.05,
		WagePence:          200,
		WorkerSupply:       50,
		Periods:            5,
	}
	snaps := simulation.RunGeneralLaw(s)
	first := snaps[0]
	last := snaps[len(snaps)-1]
	firstTotal := int64(first.ConstantCapital) + int64(first.VariableCapital)
	lastTotal := int64(last.ConstantCapital) + int64(last.VariableCapital)
	firstRatio := float64(first.Workers) / float64(firstTotal)
	lastRatio := float64(last.Workers) / float64(lastTotal)
	if lastRatio >= firstRatio {
		t.Errorf("labour share should fall as OC rises: first=%.6f last=%.6f", firstRatio, lastRatio)
	}
}

// Zero periods returns an empty slice, never nil.
func TestRunGeneralLaw_ZeroPeriods(t *testing.T) {
	t.Parallel()
	s := simulation.GeneralLawScenario{
		ConstantCapital:    8000,
		VariableCapital:    2000,
		SurplusRate:        1.0,
		AccumulationRate:   1.0,
		ProductivityGrowth: 0.05,
		WagePence:          200,
		WorkerSupply:       50,
		Periods:            0,
	}
	snaps := simulation.RunGeneralLaw(s)
	if snaps == nil {
		t.Error("want empty slice, got nil")
	}
	if len(snaps) != 0 {
		t.Errorf("len = %d, want 0", len(snaps))
	}
}

// Invariant: OrganicComposition == c/(c+v) for each period snapshot.
func TestRunGeneralLaw_OCInvariant(t *testing.T) {
	t.Parallel()
	s := simulation.GeneralLawScenario{
		ConstantCapital:    8000,
		VariableCapital:    2000,
		SurplusRate:        0.5,
		AccumulationRate:   0.5,
		ProductivityGrowth: 0.02,
		WagePence:          200,
		WorkerSupply:       100,
		Periods:            4,
	}
	snaps := simulation.RunGeneralLaw(s)
	for i, sn := range snaps {
		total := float64(sn.ConstantCapital) + float64(sn.VariableCapital)
		want := float64(sn.ConstantCapital) / total
		if sn.OrganicComposition != want {
			t.Errorf("period %d: OC = %v, want %v", i+1, sn.OrganicComposition, want)
		}
	}
}
