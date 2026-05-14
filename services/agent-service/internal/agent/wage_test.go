package agent_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// § core example from Ch. 19:
// "daily value of labour-power is 3s., the value of the product of 6
// working-hours; working-day is 12 hours"
// → DailyPence = 36 (3 shillings), NecessaryMinutes = 360 (6 h),
//   WorkingDayMinutes = 720.

func TestHourlyWage(t *testing.T) {
	t.Parallel()

	t.Run("standard fixture 3s/12h = 3p per hour", func(t *testing.T) {
		t.Parallel()
		w := agent.Wage{DailyPence: 36, WorkingDayMinutes: 720}
		if got := agent.HourlyWage(w); got != 3 {
			t.Fatalf("HourlyWage = %d, want 3", got)
		}
	})

	t.Run("price falls: 2s/12h = 2p per hour", func(t *testing.T) {
		t.Parallel()
		// "value of labour-power may vary from 3 to 2 shillings"
		w := agent.Wage{DailyPence: 24, WorkingDayMinutes: 720}
		if got := agent.HourlyWage(w); got != 2 {
			t.Fatalf("HourlyWage = %d, want 2", got)
		}
	})

	t.Run("zero working hours returns 0", func(t *testing.T) {
		t.Parallel()
		w := agent.Wage{DailyPence: 36, WorkingDayMinutes: 0}
		if got := agent.HourlyWage(w); got != 0 {
			t.Fatalf("HourlyWage = %d, want 0", got)
		}
	})
}

func TestDecompose(t *testing.T) {
	t.Parallel()

	lpv := agent.WageLabourValue{DailyPence: 36, NecessaryMinutes: 360}
	w := agent.Wage{DailyPence: 36, WorkingDayMinutes: 720}

	t.Run("paid minutes equal necessary labour", func(t *testing.T) {
		t.Parallel()
		d := agent.Decompose(w, lpv)
		if d.PaidMinutes != lpv.NecessaryMinutes {
			t.Fatalf("PaidMinutes = %d, want %d", d.PaidMinutes, lpv.NecessaryMinutes)
		}
	})

	t.Run("decomposition covers full working day", func(t *testing.T) {
		t.Parallel()
		d := agent.Decompose(w, lpv)
		total := w.WorkingDayMinutes
		if d.PaidMinutes+d.UnpaidMinutes != total {
			t.Fatalf("PaidMinutes(%d) + UnpaidMinutes(%d) = %d, want %d",
				d.PaidMinutes, d.UnpaidMinutes, d.PaidMinutes+d.UnpaidMinutes, total)
		}
	})

	t.Run("standard fixture 6h paid 6h unpaid", func(t *testing.T) {
		t.Parallel()
		d := agent.Decompose(w, lpv)
		if d.PaidMinutes != 360 || d.UnpaidMinutes != 360 {
			t.Fatalf("Decompose = {%d, %d}, want {360, 360}", d.PaidMinutes, d.UnpaidMinutes)
		}
	})

	t.Run("wage rises lpv unchanged decomposition unchanged", func(t *testing.T) {
		t.Parallel()
		wRise := agent.Wage{DailyPence: 48, WorkingDayMinutes: 720}
		d := agent.Decompose(wRise, lpv)
		if d.PaidMinutes != 360 {
			t.Fatalf("PaidMinutes = %d, want 360", d.PaidMinutes)
		}
	})
}

func TestAppearance(t *testing.T) {
	t.Parallel()

	w := agent.Wage{DailyPence: 36, WorkingDayMinutes: 720}

	t.Run("all hours appear paid zero unpaid", func(t *testing.T) {
		t.Parallel()
		a := agent.Appearance(w)
		if a.PaidHours != 12 || a.UnpaidHours != 0 {
			t.Fatalf("Appearance = {%d, %d}, want {12, 0}", a.PaidHours, a.UnpaidHours)
		}
	})
}

func TestWageForm_Validate(t *testing.T) {
	t.Parallel()

	valid := agent.WageForm{
		AgentID:          agent.AgentID("abc123"),
		Wage:             agent.Wage{DailyPence: 36, WorkingDayMinutes: 720},
		LabourPowerValue: agent.WageLabourValue{DailyPence: 36, NecessaryMinutes: 360},
	}

	t.Run("valid wage form", func(t *testing.T) {
		t.Parallel()
		if err := valid.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("missing agent_id", func(t *testing.T) {
		t.Parallel()
		wf := valid
		wf.AgentID = ""
		if err := wf.Validate(); err == nil {
			t.Fatal("expected error for empty agent_id")
		}
	})

	t.Run("zero daily_pence", func(t *testing.T) {
		t.Parallel()
		wf := valid
		wf.Wage.DailyPence = 0
		if err := wf.Validate(); err == nil {
			t.Fatal("expected error for zero daily_pence")
		}
	})

	t.Run("zero working_day_minutes", func(t *testing.T) {
		t.Parallel()
		wf := valid
		wf.Wage.WorkingDayMinutes = 0
		if err := wf.Validate(); err == nil {
			t.Fatal("expected error for zero working_day_minutes")
		}
	})

	t.Run("zero necessary_minutes", func(t *testing.T) {
		t.Parallel()
		wf := valid
		wf.LabourPowerValue.NecessaryMinutes = 0
		if err := wf.Validate(); err == nil {
			t.Fatal("expected error for zero necessary_minutes")
		}
	})
}
