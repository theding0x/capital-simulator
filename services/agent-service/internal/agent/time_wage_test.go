package agent_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// § unit measure: daily value 3s. (36p), 12-hour day → price of 1 hour = 36/12 = 3d.
func TestComputeHourlyPrice_12Hour(t *testing.T) {
	t.Parallel()

	hp := agent.ComputeHourlyPrice(
		agent.DailyLabourPowerValue{Pence: 36},
		agent.WorkingDayHours{Hours: 12},
	)
	if hp.Numerator != 36 {
		t.Errorf("numerator = %d, want 36", hp.Numerator)
	}
	if hp.Denominator != 12 {
		t.Errorf("denominator = %d, want 12", hp.Denominator)
	}
	if hp.AsFloat() != 3.0 {
		t.Errorf("as_float = %f, want 3.0", hp.AsFloat())
	}
}

// § 10-hour day: daily value 3s. (36p), 10 hours → price = 36/10 = 3.6d (stored exact).
func TestComputeHourlyPrice_10Hour(t *testing.T) {
	t.Parallel()

	hp := agent.ComputeHourlyPrice(
		agent.DailyLabourPowerValue{Pence: 36},
		agent.WorkingDayHours{Hours: 10},
	)
	if hp.Numerator != 36 || hp.Denominator != 10 {
		t.Errorf("fraction = %d/%d, want 36/10", hp.Numerator, hp.Denominator)
	}
	const want = 3.6
	if got := hp.AsFloat(); got != want {
		t.Errorf("as_float = %f, want %f", got, want)
	}
}

// § 15-hour day: daily value 3s. (36p), 15 hours → price = 36/15 = 2.4d.
func TestComputeHourlyPrice_15Hour(t *testing.T) {
	t.Parallel()

	hp := agent.ComputeHourlyPrice(
		agent.DailyLabourPowerValue{Pence: 36},
		agent.WorkingDayHours{Hours: 15},
	)
	if hp.Numerator != 36 || hp.Denominator != 15 {
		t.Errorf("fraction = %d/%d, want 36/15", hp.Numerator, hp.Denominator)
	}
	const want = 2.4
	if got := hp.AsFloat(); got != want {
		t.Errorf("as_float = %f, want %f", got, want)
	}
}

// Invariant: hourly price falls when the working day lengthens with unchanged daily value.
func TestComputeHourlyPrice_LongerDayLowersPrice(t *testing.T) {
	t.Parallel()

	v := agent.DailyLabourPowerValue{Pence: 36}
	short := agent.ComputeHourlyPrice(v, agent.WorkingDayHours{Hours: 10})
	long := agent.ComputeHourlyPrice(v, agent.WorkingDayHours{Hours: 12})

	if short.AsFloat() <= long.AsFloat() {
		t.Errorf("10h price (%f) should exceed 12h price (%f)", short.AsFloat(), long.AsFloat())
	}
}

// Invariant: fraction is stored exactly — numerator = daily pence, denominator = hours.
func TestComputeHourlyPrice_FractionIsExact(t *testing.T) {
	t.Parallel()

	v := agent.DailyLabourPowerValue{Pence: 36}
	h := agent.WorkingDayHours{Hours: 12}
	hp := agent.ComputeHourlyPrice(v, h)

	if hp.Numerator != v.Pence {
		t.Errorf("numerator = %d, want %d", hp.Numerator, v.Pence)
	}
	if hp.Denominator != h.Hours {
		t.Errorf("denominator = %d, want %d", hp.Denominator, h.Hours)
	}
}

// § nominal wage stable: 12h at 36p → total = 36p (no overtime).
func TestComputeSessionWage_NoOvertime(t *testing.T) {
	t.Parallel()

	s := agent.WorkingSession{
		AgentID:               "worker1",
		DailyLabourPowerValue: agent.DailyLabourPowerValue{Pence: 36},
		WorkingDayHours:       agent.WorkingDayHours{Hours: 12},
		OvertimeHours:         agent.OvertimeHours{Hours: 0},
		OvertimeRatePence:     agent.OvertimeRatePence{Pence: 0},
		WagePeriod:            agent.WagePeriodDaily,
	}
	wage := agent.ComputeSessionWage(s)
	if wage.Pence != 36 {
		t.Errorf("nominal_wage = %d, want 36", wage.Pence)
	}
}

// § overtime: 12-hour normal day at 3d/h; 2 overtime hours at 4d → total = 36 + 8 = 44p.
func TestComputeSessionWage_WithOvertime(t *testing.T) {
	t.Parallel()

	s := agent.WorkingSession{
		AgentID:               "worker1",
		DailyLabourPowerValue: agent.DailyLabourPowerValue{Pence: 36},
		WorkingDayHours:       agent.WorkingDayHours{Hours: 12},
		OvertimeHours:         agent.OvertimeHours{Hours: 2},
		OvertimeRatePence:     agent.OvertimeRatePence{Pence: 4},
		WagePeriod:            agent.WagePeriodDaily,
	}
	wage := agent.ComputeSessionWage(s)
	if wage.Pence != 44 {
		t.Errorf("nominal_wage = %d, want 44 (36 + 2×4)", wage.Pence)
	}
}

// Invariant: the wage formula matches the spec's integer-arithmetic definition.
func TestComputeSessionWage_IntegerArithmetic(t *testing.T) {
	t.Parallel()

	s := agent.WorkingSession{
		AgentID:               "worker1",
		DailyLabourPowerValue: agent.DailyLabourPowerValue{Pence: 36},
		WorkingDayHours:       agent.WorkingDayHours{Hours: 12},
		OvertimeHours:         agent.OvertimeHours{Hours: 3},
		OvertimeRatePence:     agent.OvertimeRatePence{Pence: 4},
		WagePeriod:            agent.WagePeriodDaily,
	}
	want := (s.DailyLabourPowerValue.Pence/s.WorkingDayHours.Hours)*s.WorkingDayHours.Hours +
		s.OvertimeHours.Hours*s.OvertimeRatePence.Pence
	got := agent.ComputeSessionWage(s).Pence
	if got != want {
		t.Errorf("nominal_wage = %d, want %d", got, want)
	}
}

// Validate rejects zero agent_id.
func TestWorkingSession_Validate_MissingAgentID(t *testing.T) {
	t.Parallel()

	s := agent.WorkingSession{
		DailyLabourPowerValue: agent.DailyLabourPowerValue{Pence: 36},
		WorkingDayHours:       agent.WorkingDayHours{Hours: 12},
		WagePeriod:            agent.WagePeriodDaily,
	}
	if err := s.Validate(); err == nil {
		t.Error("expected validation error for empty agent_id")
	}
}

// Validate rejects invalid wage_period.
func TestWorkingSession_Validate_InvalidPeriod(t *testing.T) {
	t.Parallel()

	s := agent.WorkingSession{
		AgentID:               "worker1",
		DailyLabourPowerValue: agent.DailyLabourPowerValue{Pence: 36},
		WorkingDayHours:       agent.WorkingDayHours{Hours: 12},
		WagePeriod:            agent.WagePeriod("monthly"),
	}
	if err := s.Validate(); err == nil {
		t.Error("expected validation error for invalid wage_period")
	}
}

// Validate accepts weekly period.
func TestWorkingSession_Validate_WeeklyPeriod(t *testing.T) {
	t.Parallel()

	s := agent.WorkingSession{
		AgentID:               "worker1",
		DailyLabourPowerValue: agent.DailyLabourPowerValue{Pence: 36},
		WorkingDayHours:       agent.WorkingDayHours{Hours: 12},
		OvertimeHours:         agent.OvertimeHours{Hours: 0},
		OvertimeRatePence:     agent.OvertimeRatePence{Pence: 0},
		WagePeriod:            agent.WagePeriodWeekly,
	}
	if err := s.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

// HourlyPriceOfLabour.AsFloat returns 0 for zero denominator.
func TestHourlyPriceOfLabour_AsFloat_ZeroDenominator(t *testing.T) {
	t.Parallel()

	hp := agent.HourlyPriceOfLabour{Numerator: 36, Denominator: 0}
	if hp.AsFloat() != 0 {
		t.Errorf("as_float = %f, want 0 for zero denominator", hp.AsFloat())
	}
}
