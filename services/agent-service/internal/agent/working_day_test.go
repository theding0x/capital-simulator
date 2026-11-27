package agent_test

import (
	"errors"
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

func TestWorkingDay_Validate_NegativeSurplus(t *testing.T) {
	t.Parallel()
	wd := agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: -60}
	if err := wd.Validate(); err != agent.ErrSurplusLabourNegative {
		t.Fatalf("expected ErrSurplusLabourNegative, got %v", err)
	}
}

func TestWorkingDay_Validate_ZeroSurplus(t *testing.T) {
	t.Parallel()
	wd := agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 0}
	if err := wd.Validate(); err != nil {
		t.Fatalf("zero surplus is the limiting case and must be admissible, got %v", err)
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

// Marx's Factory Act fixtures (Capital Vol. I Ch. 10 §6): the daily limits in
// force once each Act was fully operative. The adult limit is 0 ("no statutory
// limit") in all three because none of the Acts limited the working day of the
// male worker over 18.
var (
	// 1833: children 8 hours, young persons 12 hours.
	factoryAct1833 = agent.FactoryAct{Year: 1833, ChildLimitMinutes: 480, YoungPersonLimitMinutes: 720}
	// 1844: children cut to 6½ hours; women and young persons stay at 12.
	factoryAct1844 = agent.FactoryAct{Year: 1844, ChildLimitMinutes: 390, YoungPersonLimitMinutes: 720}
	// 1847: the Ten Hours' Act, 10 hours for young persons, definitive from
	// 1 May 1848. It only amended and completed the Acts of 1833 and 1844, so
	// children stay on the 1844 limit of 6½ hours.
	factoryAct1847 = agent.FactoryAct{Year: 1847, ChildLimitMinutes: 390, YoungPersonLimitMinutes: 600}
)

// requireRejected asserts that Constraint refused with want (and not with
// other) and handed back the zero constraint rather than a usable one.
func requireRejected(t *testing.T, c agent.WorkingDayConstraint, err, want, other error) {
	t.Helper()
	if !errors.Is(err, want) {
		t.Fatalf("Constraint error = %v, want %v", err, want)
	}
	if errors.Is(err, other) {
		t.Fatalf("Constraint error %v must not be %v", err, other)
	}
	if c != (agent.WorkingDayConstraint{}) {
		t.Fatalf("Constraint on error = %+v, want the zero constraint", c)
	}
}

// Constraint derives the limit and label for each Act and age band.
func TestFactoryAct_Constraint_LimitAndLabel(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name      string
		act       agent.FactoryAct
		band      agent.AgeBand
		wantLimit agent.StatutoryLimitMinutes
		wantLabel string
	}{
		{"1833 child", factoryAct1833, agent.AgeBandChild, 480, "Factory Act 1833 (child)"},
		{"1833 young person", factoryAct1833, agent.AgeBandYoungPerson, 720, "Factory Act 1833 (young person)"},
		{"1844 child", factoryAct1844, agent.AgeBandChild, 390, "Factory Act 1844 (child)"},
		{"1844 young person", factoryAct1844, agent.AgeBandYoungPerson, 720, "Factory Act 1844 (young person)"},
		{"1847 child inherits 1844", factoryAct1847, agent.AgeBandChild, 390, "Factory Act 1847 (child)"},
		{"1847 young person", factoryAct1847, agent.AgeBandYoungPerson, 600, "Factory Act 1847 (young person)"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := tc.act.Constraint(tc.band)
			if err != nil {
				t.Fatalf("Constraint(%q) error = %v, want nil", tc.band, err)
			}
			want := agent.WorkingDayConstraint{Label: tc.wantLabel, StatutoryLimitMinutes: tc.wantLimit}
			if got != want {
				t.Fatalf("Constraint(%q) = %+v, want %+v", tc.band, got, want)
			}
		})
	}
}

// None of the three Acts limited the adult male's working day, so the adult
// band has no constraint to derive: an error, never a zero-minute cap.
func TestFactoryAct_Constraint_AdultHasNoStatutoryLimit(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		act  agent.FactoryAct
	}{
		{"1833", factoryAct1833},
		{"1844", factoryAct1844},
		{"1847", factoryAct1847},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c, err := tc.act.Constraint(agent.AgeBandAdult)
			requireRejected(t, c, err, agent.ErrNoStatutoryLimit, agent.ErrUnknownAgeBand)
		})
	}
}

// A day exactly at the derived limit is admitted and one minute over is
// rejected, through both ValidateConstraint and NormalWorkingDay.Validate.
func TestFactoryAct_Constraint_BindsWorkingDay(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name               string
		act                agent.FactoryAct
		band               agent.AgeBand
		necessary, surplus int64
		wantErr            error
	}{
		{"1833 child, 8h day at limit", factoryAct1833, agent.AgeBandChild, 360, 120, nil},
		{"1833 child, 9h day over limit", factoryAct1833, agent.AgeBandChild, 360, 180, agent.ErrWorkingDayExceedsStatutoryLimit},
		{"1833 young person, 12h day at limit", factoryAct1833, agent.AgeBandYoungPerson, 360, 360, nil},
		{"1833 young person, one minute over 12h", factoryAct1833, agent.AgeBandYoungPerson, 360, 361, agent.ErrWorkingDayExceedsStatutoryLimit},
		{"1844 child, 6.5h day at limit", factoryAct1844, agent.AgeBandChild, 360, 30, nil},
		{"1844 child, one minute over 6.5h", factoryAct1844, agent.AgeBandChild, 360, 31, agent.ErrWorkingDayExceedsStatutoryLimit},
		{"1847 young person, 10h day at limit", factoryAct1847, agent.AgeBandYoungPerson, 360, 240, nil},
		{"1847 young person, one minute over 10h", factoryAct1847, agent.AgeBandYoungPerson, 360, 241, agent.ErrWorkingDayExceedsStatutoryLimit},
		{"1847 young person, 12h day over limit", factoryAct1847, agent.AgeBandYoungPerson, 360, 360, agent.ErrWorkingDayExceedsStatutoryLimit},
		{"1847 child, 6.5h day at inherited limit", factoryAct1847, agent.AgeBandChild, 360, 30, nil},
		{"1847 child, one minute over inherited 6.5h", factoryAct1847, agent.AgeBandChild, 360, 31, agent.ErrWorkingDayExceedsStatutoryLimit},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c, err := tc.act.Constraint(tc.band)
			if err != nil {
				t.Fatalf("Constraint(%q) error = %v, want nil", tc.band, err)
			}
			wd := agent.WorkingDay{
				NecessaryLabourMinutes: agent.NecessaryLabourMinutes(tc.necessary),
				SurplusLabourMinutes:   agent.SurplusLabourMinutes(tc.surplus),
			}
			if err := wd.Validate(); err != nil {
				t.Fatalf("fixture day %+v is not a valid WorkingDay: %v", wd, err)
			}
			if err := agent.ValidateConstraint(wd, c); !errors.Is(err, tc.wantErr) {
				t.Errorf("ValidateConstraint(%d min vs limit %d) = %v, want %v", wd.TotalMinutes(), c.StatutoryLimitMinutes, err, tc.wantErr)
			}
			n := agent.NormalWorkingDay{WorkingDay: wd, Constraint: c}
			if err := n.Validate(); !errors.Is(err, tc.wantErr) {
				t.Errorf("NormalWorkingDay.Validate() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

// A limit of zero (or less) means the Act sets no limit for that band; it is
// reported as ErrNoStatutoryLimit, not as an unknown band, so a missing limit
// can never turn into a zero-minute cap that rejects every day.
func TestFactoryAct_Constraint_NoLimitSet(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		act  agent.FactoryAct
		band agent.AgeBand
	}{
		{"zero act, child", agent.FactoryAct{}, agent.AgeBandChild},
		{"zero act, young person", agent.FactoryAct{}, agent.AgeBandYoungPerson},
		{"zero act, adult", agent.FactoryAct{}, agent.AgeBandAdult},
		{"negative child limit", agent.FactoryAct{ChildLimitMinutes: -1}, agent.AgeBandChild},
		{"negative young person limit", agent.FactoryAct{YoungPersonLimitMinutes: -1}, agent.AgeBandYoungPerson},
		{"negative adult limit", agent.FactoryAct{AdultLimitMinutes: -1}, agent.AgeBandAdult},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c, err := tc.act.Constraint(tc.band)
			requireRejected(t, c, err, agent.ErrNoStatutoryLimit, agent.ErrUnknownAgeBand)
		})
	}
}

// An unrecognised band is ErrUnknownAgeBand even on an Act that sets every
// limit or none at all, and is never confused with the no-limit case. The
// populated act is synthetic (Year 0): Marx gives no adult limit, so no
// historical Act is used for it.
func TestFactoryAct_Constraint_UnknownAgeBand(t *testing.T) {
	t.Parallel()
	populated := agent.FactoryAct{Year: 0, ChildLimitMinutes: 390, YoungPersonLimitMinutes: 600, AdultLimitMinutes: 720}
	for _, tc := range []struct {
		name string
		act  agent.FactoryAct
		band agent.AgeBand
	}{
		{"populated act, empty band", populated, agent.AgeBand("")},
		{"populated act, infant", populated, agent.AgeBand("infant")},
		{"zero act, empty band", agent.FactoryAct{}, agent.AgeBand("")},
		{"zero act, infant", agent.FactoryAct{}, agent.AgeBand("infant")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c, err := tc.act.Constraint(tc.band)
			requireRejected(t, c, err, agent.ErrUnknownAgeBand, agent.ErrNoStatutoryLimit)
		})
	}
}
