package simulation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// §main monotone inequality: of the peasants an episode dispossesses, only some
// become wage-labourers — NewProletarians ≤ PeasantsDispossessed. Sutherland
// clearances 1814–1820: 15,000 driven off (Marx's footnote figure).
func TestExpropriationEpisode_NewProletariansBounded(t *testing.T) {
	t.Parallel()

	t.Run("fewer become proletarians than were dispossessed", func(t *testing.T) {
		t.Parallel()
		e := simulation.ExpropriationEpisode{Period: "1814–1820", PeasantsDispossessed: 15000, NewProletarians: 12000}
		if err := e.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("all dispossessed become proletarians (boundary)", func(t *testing.T) {
		t.Parallel()
		e := simulation.ExpropriationEpisode{Period: "1814–1820", PeasantsDispossessed: 15000, NewProletarians: 15000}
		if err := e.Validate(); err != nil {
			t.Fatalf("equality is allowed, got %v", err)
		}
	})

	t.Run("more proletarians than dispossessed violates the inequality", func(t *testing.T) {
		t.Parallel()
		e := simulation.ExpropriationEpisode{Period: "1814–1820", PeasantsDispossessed: 15000, NewProletarians: 15001}
		if err := e.Validate(); err != simulation.ErrNewProletariansExceedDispossessed {
			t.Fatalf("want ErrNewProletariansExceedDispossessed, got %v", err)
		}
	})

	t.Run("negative counts rejected", func(t *testing.T) {
		t.Parallel()
		e := simulation.ExpropriationEpisode{Period: "1814–1820", PeasantsDispossessed: -1}
		if err := e.Validate(); err != simulation.ErrDispossessedNegative {
			t.Fatalf("want ErrDispossessedNegative, got %v", err)
		}
	})

	t.Run("missing period rejected", func(t *testing.T) {
		t.Parallel()
		e := simulation.ExpropriationEpisode{PeasantsDispossessed: 15000, NewProletarians: 12000}
		if err := e.Validate(); err != simulation.ErrPeriodRequired {
			t.Fatalf("want ErrPeriodRequired, got %v", err)
		}
	})
}

// Acre conservation: the commons enclosed reappear acre-for-acre as private
// land. Sutherland: 794,000 acres turned from clan land into sheep-walk.
func TestEnclosure_AcreConservation(t *testing.T) {
	t.Parallel()

	t.Run("common acres become private acre-for-acre", func(t *testing.T) {
		t.Parallel()
		e := simulation.Enclosure{Period: "1814–1820", CommonAcres: 794000, PrivateAcres: 794000}
		if err := e.Validate(); err != nil {
			t.Fatalf("conserved enclosure should be valid, got %v", err)
		}
	})

	t.Run("acreage not conserved is rejected", func(t *testing.T) {
		t.Parallel()
		e := simulation.Enclosure{Period: "1814–1820", CommonAcres: 794000, PrivateAcres: 700000}
		if err := e.Validate(); err != simulation.ErrAcresNotConserved {
			t.Fatalf("want ErrAcresNotConserved, got %v", err)
		}
	})

	t.Run("zero acres enclosed is rejected", func(t *testing.T) {
		t.Parallel()
		e := simulation.Enclosure{Period: "1450–1540", CommonAcres: 0, PrivateAcres: 0}
		if err := e.Validate(); err != simulation.ErrAcresNonPositive {
			t.Fatalf("want ErrAcresNonPositive, got %v", err)
		}
	})
}

// Aggregation invariant: the proletariat the transition produces is the sum of
// the proletarians each episode produced. Three canonical English waves: 15th-c.
// gentry enclosures, the Parliamentary enclosures 1760–1820, and Sutherland.
func TestAgrarianTransition_TotalIsSumOfEpisodes(t *testing.T) {
	t.Parallel()

	episodes := []simulation.ExpropriationEpisode{
		{Period: "1450–1540", PeasantsDispossessed: 40000, NewProletarians: 30000},
		{Period: "1760–1820", PeasantsDispossessed: 150000, NewProletarians: 120000},
		{Period: "1814–1820", PeasantsDispossessed: 15000, NewProletarians: 15000},
	}

	t.Run("total equal to the sum is valid", func(t *testing.T) {
		t.Parallel()
		tr := simulation.AgrarianTransition{Episodes: episodes, TotalProletarians: 165000}
		if got := tr.SumNewProletarians(); got != 165000 {
			t.Fatalf("SumNewProletarians() = %d, want 165000", got)
		}
		if err := tr.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("total not equal to the sum is rejected", func(t *testing.T) {
		t.Parallel()
		tr := simulation.AgrarianTransition{Episodes: episodes, TotalProletarians: 160000}
		if err := tr.Validate(); err != simulation.ErrProletarianSumMismatch {
			t.Fatalf("want ErrProletarianSumMismatch, got %v", err)
		}
	})

	t.Run("a bad episode invalidates the whole transition", func(t *testing.T) {
		t.Parallel()
		bad := []simulation.ExpropriationEpisode{
			{Period: "1814–1820", PeasantsDispossessed: 15000, NewProletarians: 20000},
		}
		tr := simulation.AgrarianTransition{Episodes: bad, TotalProletarians: 20000}
		if err := tr.Validate(); err != simulation.ErrNewProletariansExceedDispossessed {
			t.Fatalf("want the episode error to propagate, got %v", err)
		}
	})
}
