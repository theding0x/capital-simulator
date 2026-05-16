package simulation_test

import (
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// England 15th–18th c. — the canonical fixture used throughout the tests
// below. Magnitudes are pedagogical and mirror 00011_ch26_seed.sql.
func englandStage() simulation.HistoricalStage {
	return simulation.HistoricalStage{
		Name:        "England 15th-18th c.",
		Description: "agrarian expropriation creating free proletariat",
		PrimitiveAccumulations: []simulation.PrimitiveAccumulation{
			{Period: "1450-1640", Method: "enclosure of common lands", LabourersExpropriated: 40000, CapitalFormed: 80000},
			{Period: "1500-1640", Method: "expropriation of church property", LabourersExpropriated: 15000, CapitalFormed: 30000},
			{Period: "1500-1800", Method: "colonial plunder and slave trade", LabourersExpropriated: 5000, CapitalFormed: 90000},
		},
	}
}

func TestPrimitiveAccumulation_ValidateRequiresPositiveCapital(t *testing.T) {
	t.Parallel()
	p := simulation.PrimitiveAccumulation{Period: "X", Method: "y", LabourersExpropriated: 1, CapitalFormed: 0}
	if err := p.Validate(); !errors.Is(err, simulation.ErrCapitalFormedNonPositive) {
		t.Fatalf("want ErrCapitalFormedNonPositive, got %v", err)
	}
}

func TestPrimitiveAccumulation_ValidateRejectsNegativeLabourers(t *testing.T) {
	t.Parallel()
	p := simulation.PrimitiveAccumulation{Period: "X", Method: "y", LabourersExpropriated: -1, CapitalFormed: 100}
	if err := p.Validate(); !errors.Is(err, simulation.ErrLabourersNegative) {
		t.Fatalf("want ErrLabourersNegative, got %v", err)
	}
}

func TestPrimitiveAccumulation_ValidateAcceptsCanonical(t *testing.T) {
	t.Parallel()
	for _, p := range englandStage().PrimitiveAccumulations {
		if err := p.Validate(); err != nil {
			t.Fatalf("England episode %q failed validation: %v", p.Method, err)
		}
	}
}

func TestProducerSeparation_FreeProletariansEqualDisplaced(t *testing.T) {
	t.Parallel()
	good := simulation.ProducerSeparation{
		PreCapitalistWorkers: 100000,
		DisplacedWorkers:     40000,
		FreeProletarians:     40000,
	}
	if err := good.Validate(); err != nil {
		t.Fatalf("conservation case rejected: %v", err)
	}

	bad := simulation.ProducerSeparation{
		PreCapitalistWorkers: 100000,
		DisplacedWorkers:     40000,
		FreeProletarians:     30000,
	}
	if err := bad.Validate(); !errors.Is(err, simulation.ErrSeparationMismatch) {
		t.Fatalf("want ErrSeparationMismatch, got %v", err)
	}
}

func TestProducerSeparation_RejectsNegative(t *testing.T) {
	t.Parallel()
	s := simulation.ProducerSeparation{PreCapitalistWorkers: -1, DisplacedWorkers: 0, FreeProletarians: 0}
	if err := s.Validate(); !errors.Is(err, simulation.ErrLabourersNegative) {
		t.Fatalf("want ErrLabourersNegative, got %v", err)
	}
}

func TestHistoricalStage_ValidateRequiresName(t *testing.T) {
	t.Parallel()
	stage := englandStage()
	stage.Name = ""
	if err := stage.Validate(); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestHistoricalStage_ValidateRequiresAtLeastOneEpisode(t *testing.T) {
	t.Parallel()
	stage := englandStage()
	stage.PrimitiveAccumulations = nil
	if err := stage.Validate(); !errors.Is(err, simulation.ErrStageEmpty) {
		t.Fatalf("want ErrStageEmpty, got %v", err)
	}
}

func TestHistoricalStage_ValidateBubblesEpisodeError(t *testing.T) {
	t.Parallel()
	stage := englandStage()
	stage.PrimitiveAccumulations[1].CapitalFormed = 0
	if err := stage.Validate(); !errors.Is(err, simulation.ErrCapitalFormedNonPositive) {
		t.Fatalf("want ErrCapitalFormedNonPositive, got %v", err)
	}
}

func TestHistoricalStage_Totals(t *testing.T) {
	t.Parallel()
	stage := englandStage()
	if got, want := stage.TotalCapitalFormed(), simulation.Pence(200000); got != want {
		t.Errorf("TotalCapitalFormed = %d, want %d", got, want)
	}
	if got, want := stage.TotalLabourersExpropriated(), int64(60000); got != want {
		t.Errorf("TotalLabourersExpropriated = %d, want %d", got, want)
	}
}

func TestSeparationFromStage_ConservesDispossessed(t *testing.T) {
	t.Parallel()
	stage := englandStage() // 60_000 total dispossessed
	sep := simulation.SeparationFromStage(stage, 100000)
	if sep.DisplacedWorkers != 60000 {
		t.Errorf("DisplacedWorkers = %d, want 60000", sep.DisplacedWorkers)
	}
	if sep.FreeProletarians != sep.DisplacedWorkers {
		t.Errorf("FreeProletarians = %d, want == DisplacedWorkers (%d)", sep.FreeProletarians, sep.DisplacedWorkers)
	}
	if err := sep.Validate(); err != nil {
		t.Errorf("returned separation fails Validate: %v", err)
	}
}

func TestSeparationFromStage_ClampsToPopulation(t *testing.T) {
	t.Parallel()
	stage := englandStage() // 60_000 total dispossessed
	sep := simulation.SeparationFromStage(stage, 50000)
	if sep.DisplacedWorkers != 50000 {
		t.Errorf("DisplacedWorkers = %d, want 50000 (clamped)", sep.DisplacedWorkers)
	}
	if sep.FreeProletarians != sep.DisplacedWorkers {
		t.Errorf("invariant broken after clamping: FP=%d, DW=%d", sep.FreeProletarians, sep.DisplacedWorkers)
	}
}

func TestSeedCapitalStock_PartitionsTotal(t *testing.T) {
	t.Parallel()
	stage := englandStage() // total = 200_000
	cs := simulation.SeedCapitalStock(stage, 0.8)
	total := cs.ConstantCapital + cs.VariableCapital
	if total != simulation.Pence(200000) {
		t.Errorf("partition broken: c+v = %d, want 200000", total)
	}
	if cs.ConstantCapital != simulation.Pence(160000) {
		t.Errorf("ConstantCapital = %d, want 160000 (0.8 * 200000)", cs.ConstantCapital)
	}
}

func TestSeedCapitalStock_ClampsRatio(t *testing.T) {
	t.Parallel()
	stage := englandStage()
	cs := simulation.SeedCapitalStock(stage, 1.5)
	if cs.ConstantCapital+cs.VariableCapital != stage.TotalCapitalFormed() {
		t.Errorf("partition broken after clamp: got c=%d v=%d", cs.ConstantCapital, cs.VariableCapital)
	}
	if cs.VariableCapital <= 0 {
		t.Error("clamping must leave positive variable capital")
	}
}

func TestSeedCapitalStock_EmptyTotalYieldsZero(t *testing.T) {
	t.Parallel()
	stage := simulation.HistoricalStage{Name: "empty"}
	cs := simulation.SeedCapitalStock(stage, 0.5)
	if cs.ConstantCapital != 0 || cs.VariableCapital != 0 {
		t.Errorf("empty stage must yield zero capital, got c=%d v=%d", cs.ConstantCapital, cs.VariableCapital)
	}
}

func TestNewHistoricalStageID_Unique(t *testing.T) {
	t.Parallel()
	seen := make(map[simulation.HistoricalStageID]struct{})
	for i := 0; i < 64; i++ {
		id := simulation.NewHistoricalStageID()
		if id.IsZero() {
			t.Fatal("generated zero id")
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id generated: %s", id)
		}
		seen[id] = struct{}{}
		if len(string(id)) != 24 {
			t.Fatalf("id length = %d, want 24 hex chars", len(string(id)))
		}
	}
}
