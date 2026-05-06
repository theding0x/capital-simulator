package agent_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

func TestWorkingDay_TotalMinutes(t *testing.T) {
	t.Parallel()
	wd := agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 60}
	if got := wd.TotalMinutes(); got != 420 {
		t.Fatalf("TotalMinutes = %d, want 420", got)
	}
}

func TestWorkingDay_RateOfSurplusValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		nl, sl int64
		want   float64
	}{
		{360, 60, 1.0 / 6.0},
		{360, 180, 0.50},
		{360, 360, 1.00},
		{480, 480, 1.00},
		{600, 600, 1.00},
		{720, 720, 1.00},
	}
	for _, tt := range tests {
		wd := agent.WorkingDay{
			NecessaryLabourMinutes: agent.NecessaryLabourMinutes(tt.nl),
			SurplusLabourMinutes:   agent.SurplusLabourMinutes(tt.sl),
		}
		got := wd.RateOfSurplusValue()
		if diff := got - tt.want; diff > 1e-9 || diff < -1e-9 {
			t.Errorf("nl=%d sl=%d: RateOfSurplusValue()=%f, want %f", tt.nl, tt.sl, got, tt.want)
		}
	}
}

func TestWorkingDay_Validate_PhysicalMax(t *testing.T) {
	t.Parallel()
	wd := agent.WorkingDay{NecessaryLabourMinutes: 720, SurplusLabourMinutes: 720}
	if err := wd.Validate(); err != nil {
		t.Fatalf("expected no error for 1440 min, got %v", err)
	}
	over := agent.WorkingDay{NecessaryLabourMinutes: 721, SurplusLabourMinutes: 720}
	if err := over.Validate(); err != agent.ErrWorkingDayExceedsPhysicalMax {
		t.Fatalf("expected ErrWorkingDayExceedsPhysicalMax, got %v", err)
	}
}

func TestWorkingDay_Validate_ZeroNecessary(t *testing.T) {
	t.Parallel()
	wd := agent.WorkingDay{NecessaryLabourMinutes: 0, SurplusLabourMinutes: 60}
	if err := wd.Validate(); err == nil {
		t.Fatal("expected error for zero necessary labour, got nil")
	}
}

func TestValidateConstraint(t *testing.T) {
	t.Parallel()
	c := agent.WorkingDayConstraint{Label: "Factory Act 1850", StatutoryLimitMinutes: 630}
	ok := agent.WorkingDay{NecessaryLabourMinutes: 315, SurplusLabourMinutes: 315}
	if err := agent.ValidateConstraint(ok, c); err != nil {
		t.Fatalf("expected nil for total==limit, got %v", err)
	}
	over := agent.WorkingDay{NecessaryLabourMinutes: 316, SurplusLabourMinutes: 315}
	if err := agent.ValidateConstraint(over, c); err != agent.ErrWorkingDayExceedsStatutoryLimit {
		t.Fatalf("expected ErrWorkingDayExceedsStatutoryLimit, got %v", err)
	}
}

func TestOverwork_AnnualMinutes(t *testing.T) {
	t.Parallel()
	o := agent.Overwork{MinutesPerDay: 5}
	if got := o.AnnualMinutes(300); got != 1500 {
		t.Fatalf("AnnualMinutes(300) = %d, want 1500", got)
	}
}

func TestWeeklyFixtures(t *testing.T) {
	t.Parallel()
	wd := agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360}
	if got := wd.TotalMinutes() * 6 / 60; got != 72 {
		t.Errorf("weekly total hours = %d, want 72", got)
	}
	if got := int64(wd.SurplusLabourMinutes) * 6 / 60; got != 36 {
		t.Errorf("weekly surplus hours = %d, want 36", got)
	}
}

func TestWallachianCorvee(t *testing.T) {
	t.Parallel()
	wd := agent.WorkingDay{
		NecessaryLabourMinutes: agent.NecessaryLabourMinutes(84 * 480),
		SurplusLabourMinutes:   agent.SurplusLabourMinutes(56 * 480),
	}
	got := wd.RateOfSurplusValue()
	want := float64(56) / float64(84)
	if diff := got - want; diff > 1e-6 || diff < -1e-6 {
		t.Errorf("Wallachian rate = %f, want %f", got, want)
	}
}

func TestRelaySchedule_Validate(t *testing.T) {
	t.Parallel()
	rs := agent.RelaySchedule{
		Sets: [2]agent.RelaySet{
			{ShiftKind: agent.ShiftDay, WorkingDay: agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360}},
			{ShiftKind: agent.ShiftNight, WorkingDay: agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 0}},
		},
	}
	if err := rs.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRelaySchedule_Validate_ExceedsPhysical(t *testing.T) {
	t.Parallel()
	rs := agent.RelaySchedule{
		Sets: [2]agent.RelaySet{
			{ShiftKind: agent.ShiftDay, WorkingDay: agent.WorkingDay{NecessaryLabourMinutes: 720, SurplusLabourMinutes: 360}},
			{ShiftKind: agent.ShiftNight, WorkingDay: agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 1}},
		},
	}
	if err := rs.Validate(); err == nil {
		t.Fatal("expected error for combined > 1440 min")
	}
}

func TestFactoryActs(t *testing.T) {
	t.Parallel()
	fa1833 := agent.FactoryAct{Year: 1833, ChildLimitMinutes: 480, YoungPersonLimitMinutes: 720}
	if fa1833.Year != 1833 || fa1833.ChildLimitMinutes != 480 {
		t.Errorf("FactoryAct 1833 wrong: %+v", fa1833)
	}
	fa1847 := agent.FactoryAct{Year: 1847, YoungPersonLimitMinutes: 600}
	if fa1847.YoungPersonLimitMinutes != 600 {
		t.Errorf("FactoryAct 1847 wrong: %+v", fa1847)
	}
}
