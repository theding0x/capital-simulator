package simulation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

func TestEnclosureEventValidate(t *testing.T) {
	t.Parallel()

	t.Run("valid Sutherland clearances", func(t *testing.T) {
		t.Parallel()
		e := simulation.EnclosureEvent{
			Period:              "1814–1820",
			AcresEnclosed:       794000,
			PopulationDisplaced: 15000,
			Beneficiary:         "Countess of Sutherland",
		}
		if err := e.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("missing period", func(t *testing.T) {
		t.Parallel()
		e := simulation.EnclosureEvent{
			AcresEnclosed:       500000,
			PopulationDisplaced: 40000,
			Beneficiary:         "landed gentry",
		}
		if err := e.Validate(); err != simulation.ErrPeriodRequired {
			t.Fatalf("want ErrPeriodRequired, got %v", err)
		}
	})

	t.Run("zero acres", func(t *testing.T) {
		t.Parallel()
		e := simulation.EnclosureEvent{
			Period:              "1450–1540",
			AcresEnclosed:       0,
			PopulationDisplaced: 40000,
			Beneficiary:         "landed gentry",
		}
		if err := e.Validate(); err != simulation.ErrAcresNonPositive {
			t.Fatalf("want ErrAcresNonPositive, got %v", err)
		}
	})

	t.Run("negative acres", func(t *testing.T) {
		t.Parallel()
		e := simulation.EnclosureEvent{
			Period:              "1450–1540",
			AcresEnclosed:       -1,
			PopulationDisplaced: 40000,
			Beneficiary:         "landed gentry",
		}
		if err := e.Validate(); err != simulation.ErrAcresNonPositive {
			t.Fatalf("want ErrAcresNonPositive, got %v", err)
		}
	})

	t.Run("negative population", func(t *testing.T) {
		t.Parallel()
		e := simulation.EnclosureEvent{
			Period:              "1450–1540",
			AcresEnclosed:       500000,
			PopulationDisplaced: -1,
			Beneficiary:         "landed gentry",
		}
		if err := e.Validate(); err != simulation.ErrPopulationNegative {
			t.Fatalf("want ErrPopulationNegative, got %v", err)
		}
	})

	t.Run("zero population is allowed", func(t *testing.T) {
		t.Parallel()
		e := simulation.EnclosureEvent{
			Period:              "1760–1800",
			AcresEnclosed:       3000000,
			PopulationDisplaced: 0,
			Beneficiary:         "improving landlords",
		}
		if err := e.Validate(); err != nil {
			t.Fatalf("unexpected error for zero population: %v", err)
		}
	})

	t.Run("missing beneficiary", func(t *testing.T) {
		t.Parallel()
		e := simulation.EnclosureEvent{
			Period:              "1450–1540",
			AcresEnclosed:       500000,
			PopulationDisplaced: 40000,
		}
		if err := e.Validate(); err != simulation.ErrBeneficiaryRequired {
			t.Fatalf("want ErrBeneficiaryRequired, got %v", err)
		}
	})
}

func TestNewEnclosureEventID(t *testing.T) {
	t.Parallel()
	id1 := simulation.NewEnclosureEventID()
	id2 := simulation.NewEnclosureEventID()
	if id1 == id2 {
		t.Fatal("two consecutive IDs must differ")
	}
	if len(string(id1)) != 24 {
		t.Fatalf("expected 24-char hex ID, got %d chars", len(string(id1)))
	}
}