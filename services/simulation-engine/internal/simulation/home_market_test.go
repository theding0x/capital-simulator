package simulation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// § Westphalian flax (Ch. 30): peasants who all spun flax are forcibly
// expropriated; large flax-spinning and weaving establishments arise in
// their place. Their output, once self-consumed, enters the commodity market.
var westphalianFlax = simulation.DomesticIndustry{
	ID:                simulation.DomesticIndustryID("5eed000000000000003001"),
	HistoricalStageID: simulation.HistoricalStageID("5eed000000000000002601"),
	Name:              "Westphalian flax spinning",
	HouseholdsEngaged: 2000,
	AnnualOutputPence: 48000,
}

func TestDomesticIndustryValidate(t *testing.T) {
	t.Parallel()

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		if err := westphalianFlax.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("name required", func(t *testing.T) {
		t.Parallel()
		d := westphalianFlax
		d.Name = ""
		if err := d.Validate(); err != simulation.ErrDomesticIndustryNameRequired {
			t.Errorf("want ErrDomesticIndustryNameRequired, got %v", err)
		}
	})

	t.Run("households non-positive", func(t *testing.T) {
		t.Parallel()
		d := westphalianFlax
		d.HouseholdsEngaged = 0
		if err := d.Validate(); err != simulation.ErrHouseholdsNonPositive {
			t.Errorf("want ErrHouseholdsNonPositive, got %v", err)
		}
	})

	t.Run("annual output non-positive", func(t *testing.T) {
		t.Parallel()
		d := westphalianFlax
		d.AnnualOutputPence = 0
		if err := d.Validate(); err != simulation.ErrAnnualOutputNonPositive {
			t.Errorf("want ErrAnnualOutputNonPositive, got %v", err)
		}
	})
}

func TestComputeMarketFormation(t *testing.T) {
	t.Parallel()

	// § peasant → commodity buyer: expropriation of 1500 Westphalian households
	// creates 1500 new wage workers who must now purchase what they once made.
	t.Run("Westphalian flax partial expropriation", func(t *testing.T) {
		t.Parallel()
		f := simulation.ComputeMarketFormation(westphalianFlax, 1500)

		if f.NewWageLabourers != f.LabourersProletarianised {
			t.Errorf("invariant: NewWageLabourers(%d) != LabourersProletarianised(%d)",
				f.NewWageLabourers, f.LabourersProletarianised)
		}
		if f.NewWageLabourers != 1500 {
			t.Errorf("want 1500 new wage labourers, got %d", f.NewWageLabourers)
		}
		if f.MarketSizePence != westphalianFlax.AnnualOutputPence {
			t.Errorf("market size should equal domestic output %d, got %d",
				westphalianFlax.AnnualOutputPence, f.MarketSizePence)
		}
		if f.DomesticIndustryDestroyed {
			t.Error("partial expropriation should not mark industry as destroyed")
		}
		if f.Period != westphalianFlax.Name {
			t.Errorf("period should be %q, got %q", westphalianFlax.Name, f.Period)
		}
	})

	t.Run("total expropriation destroys domestic industry", func(t *testing.T) {
		t.Parallel()
		f := simulation.ComputeMarketFormation(westphalianFlax, 2000)
		if !f.DomesticIndustryDestroyed {
			t.Error("full expropriation should mark domestic industry as destroyed")
		}
	})

	t.Run("zero expropriation", func(t *testing.T) {
		t.Parallel()
		f := simulation.ComputeMarketFormation(westphalianFlax, 0)
		if f.NewWageLabourers != 0 {
			t.Errorf("want 0 new wage labourers, got %d", f.NewWageLabourers)
		}
		if f.DomesticIndustryDestroyed {
			t.Error("zero expropriation should not destroy industry")
		}
		if f.MarketSizePence != westphalianFlax.AnnualOutputPence {
			t.Errorf("market size should equal domestic output even at zero expropriation")
		}
	})

	t.Run("negative expropriation clamped to zero", func(t *testing.T) {
		t.Parallel()
		f := simulation.ComputeMarketFormation(westphalianFlax, -100)
		if f.NewWageLabourers != 0 {
			t.Errorf("negative expropriated should clamp to 0, got %d", f.NewWageLabourers)
		}
	})
}

// § Mirabeau's distinction (Ch. 30): "grand manufactories ... hundreds of men
// under a director" vs. "discrete workshop ... no one will become rich, but many
// labourers will be comfortable". With equal headcount, concentration produces more.
func TestMirabelaDistinction(t *testing.T) {
	t.Parallel()

	reunie := simulation.ManufactureReunie{Workers: 200, OutputPence: 12000}
	separee := simulation.ManufactureSeparee{IndependentProducers: 200, OutputPence: 8000}

	if reunie.OutputPence <= separee.OutputPence {
		t.Errorf("ManufactureReunie(%d) must exceed ManufactureSeparee(%d) at equal headcount",
			reunie.OutputPence, separee.OutputPence)
	}
}

func TestAggregateHomeMarket(t *testing.T) {
	t.Parallel()

	industries := []simulation.DomesticIndustry{
		westphalianFlax,
		{
			Name:              "Westphalian linen weaving",
			HouseholdsEngaged: 1000,
			AnnualOutputPence: 24000,
		},
	}

	size := simulation.AggregateHomeMarket(industries)

	wantOutput := simulation.Pence(72000)
	if size.CommoditisedOutputPence != wantOutput {
		t.Errorf("want CommoditisedOutputPence %d, got %d", wantOutput, size.CommoditisedOutputPence)
	}
	if size.WageLabourers != 3000 {
		t.Errorf("want 3000 wage labourers, got %d", size.WageLabourers)
	}

	t.Run("empty slice returns zero", func(t *testing.T) {
		t.Parallel()
		empty := simulation.AggregateHomeMarket(nil)
		if empty.CommoditisedOutputPence != 0 || empty.WageLabourers != 0 {
			t.Errorf("empty industries should yield zero HomeMarketSize, got %+v", empty)
		}
	})
}
