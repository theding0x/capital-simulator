package simulation

import (
	"errors"
	"testing"
)

// Ch. 32 — The Historical Tendency of Capitalist Accumulation.
// Fixtures drawn from Marx's text: the three dialectical stages and the
// asymptotic centralisation that begets socialised property.

func TestNegationOfNegation_HappyPath(t *testing.T) {
	t.Parallel()
	got, err := NegationOfNegation(DialecticalSequence())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Stage != NegationStageSocialisedProperty {
		t.Fatalf("want socialised-property, got %q", got.Stage)
	}
	if got.Description == "" {
		t.Fatal("third stage description must not be empty")
	}
}

func TestNegationOfNegation_WrongCount(t *testing.T) {
	t.Parallel()
	cases := [][]Negation{
		{},
		{{Stage: NegationStagePettyProperty}},
		{
			{Stage: NegationStagePettyProperty},
			{Stage: NegationStageCapitalistExpropriation},
		},
		{
			{Stage: NegationStagePettyProperty},
			{Stage: NegationStageCapitalistExpropriation},
			{Stage: NegationStageSocialisedProperty},
			{Stage: "communist-fourth-stage"},
		},
	}
	for i, stages := range cases {
		if _, err := NegationOfNegation(stages); !errors.Is(err, ErrNegationStagesCount) {
			t.Errorf("case %d: want ErrNegationStagesCount, got %v", i, err)
		}
	}
}

func TestNegationOfNegation_WrongOrder(t *testing.T) {
	t.Parallel()
	out := []Negation{
		{Stage: NegationStageCapitalistExpropriation},
		{Stage: NegationStagePettyProperty},
		{Stage: NegationStageSocialisedProperty},
	}
	if _, err := NegationOfNegation(out); !errors.Is(err, ErrNegationStagesOrder) {
		t.Fatalf("want ErrNegationStagesOrder, got %v", err)
	}
}

func TestDialecticalSequence_HasThreeStages(t *testing.T) {
	t.Parallel()
	seq := DialecticalSequence()
	if len(seq) != 3 {
		t.Fatalf("want 3 stages, got %d", len(seq))
	}
	if seq[0].Stage != NegationStagePettyProperty {
		t.Errorf("first stage = %q, want %q", seq[0].Stage, NegationStagePettyProperty)
	}
	if seq[1].Stage != NegationStageCapitalistExpropriation {
		t.Errorf("second stage = %q, want %q", seq[1].Stage, NegationStageCapitalistExpropriation)
	}
	if seq[2].Stage != NegationStageSocialisedProperty {
		t.Errorf("third stage = %q, want %q", seq[2].Stage, NegationStageSocialisedProperty)
	}
}

func TestRunCentralisation_LancashireCotton(t *testing.T) {
	t.Parallel()
	// "Hand in hand with this centralisation, or this expropriation of
	// many capitalists by few, develop, on an ever-extending scale, the
	// cooperative form of the labour process."
	// Lancashire cotton 1820: roughly 1,000 small mills, around £2,000
	// each, around 50,000 wage-labourers. By 1880 the trade consolidates.
	initial := CapitalistPrivateProperty{
		Firms:             1000,
		TotalCapitalPence: 2_000_000 * 240, // £2M in pence
		WageLabourers:     50_000,
	}
	traj := RunCentralisation(initial, 6, 0.2, 0.05)

	if traj.FinalFirms >= initial.Firms {
		t.Fatalf("centralisation must reduce firms: initial %d, final %d", initial.Firms, traj.FinalFirms)
	}
	if traj.FinalCapitalPence <= 0 {
		t.Fatalf("final capital must remain positive, got %d", traj.FinalCapitalPence)
	}
	if len(traj.Steps) == 0 {
		t.Fatal("trajectory must contain at least one step")
	}
	for i, s := range traj.Steps {
		if s.FirmsAbsorbed <= 0 {
			t.Errorf("step %d: firms_absorbed must be > 0, got %d", i, s.FirmsAbsorbed)
		}
		if s.CapitalConcentratedPence <= 0 {
			t.Errorf("step %d: capital_concentrated_pence must be > 0, got %d", i, s.CapitalConcentratedPence)
		}
		if s.StepIndex != int64(i)+1 {
			t.Errorf("step %d: step_index = %d, want %d", i, s.StepIndex, i+1)
		}
	}
	if traj.ReserveArmySize <= 0 {
		t.Errorf("centralisation should release labourers into the reserve army, got %d", traj.ReserveArmySize)
	}
}

func TestRunCentralisation_NeverReachesMonopoly(t *testing.T) {
	t.Parallel()
	// "The monopoly of capital becomes a fetter ... the knell of capitalist
	// private property sounds." Marx names the limit; the simulation must
	// approach it asymptotically rather than reach it.
	initial := CapitalistPrivateProperty{
		Firms:             100,
		TotalCapitalPence: 100_000,
		WageLabourers:     1_000,
	}
	traj := RunCentralisation(initial, 10_000, 0.5, 0.0)
	if traj.FinalFirms < 1 {
		t.Fatalf("simulation must not reach FinalFirms < 1, got %d", traj.FinalFirms)
	}
}

func TestRunCentralisation_ConcentrationDoesNotDestroyValue(t *testing.T) {
	t.Parallel()
	// Invariant from spec: trajectory.FinalCapitalPence >= sum(concentrated).
	// Concentration of ownership redistributes value; it does not annihilate it.
	initial := CapitalistPrivateProperty{
		Firms:             500,
		TotalCapitalPence: 500_000,
		WageLabourers:     10_000,
	}
	traj := RunCentralisation(initial, 8, 0.25, 0.0)
	var sumConcentrated int64
	for _, s := range traj.Steps {
		sumConcentrated += int64(s.CapitalConcentratedPence)
	}
	if int64(traj.FinalCapitalPence) < sumConcentrated {
		t.Fatalf("final capital %d < sum concentrated %d: value must not be destroyed", traj.FinalCapitalPence, sumConcentrated)
	}
}

func TestRunCentralisation_ZeroStepsLeavesInitialState(t *testing.T) {
	t.Parallel()
	initial := CapitalistPrivateProperty{
		Firms:             100,
		TotalCapitalPence: 10_000,
		WageLabourers:     500,
	}
	traj := RunCentralisation(initial, 0, 0.2, 0.0)
	if traj.FinalFirms != initial.Firms {
		t.Errorf("final firms = %d, want %d", traj.FinalFirms, initial.Firms)
	}
	if traj.FinalCapitalPence != initial.TotalCapitalPence {
		t.Errorf("final capital = %d, want %d", traj.FinalCapitalPence, initial.TotalCapitalPence)
	}
	if len(traj.Steps) != 0 {
		t.Errorf("want empty steps, got %d", len(traj.Steps))
	}
}

func TestRunCentralisation_SingleFirmIsTerminal(t *testing.T) {
	t.Parallel()
	initial := CapitalistPrivateProperty{
		Firms:             1,
		TotalCapitalPence: 10_000,
		WageLabourers:     50,
	}
	traj := RunCentralisation(initial, 10, 0.5, 0.0)
	if traj.FinalFirms != 1 {
		t.Fatalf("a monopoly cannot be further centralised, got final firms %d", traj.FinalFirms)
	}
	if len(traj.Steps) != 0 {
		t.Fatalf("want no steps when already at monopoly, got %d", len(traj.Steps))
	}
}

func TestAccumulationTrajectory_Validate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		traj    AccumulationTrajectory
		wantErr error
	}{
		{
			name:    "missing-name",
			traj:    AccumulationTrajectory{InitialFirms: 10},
			wantErr: ErrTrajectoryNameRequired,
		},
		{
			name:    "zero-firms",
			traj:    AccumulationTrajectory{Name: "x", InitialFirms: 0},
			wantErr: ErrTrajectoryFirmsPositive,
		},
		{
			name: "ok",
			traj: AccumulationTrajectory{Name: "Lancashire cotton, 1820-1880", InitialFirms: 1000},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.traj.Validate()
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("got %v, want %v", err, tc.wantErr)
			}
		})
	}
}
