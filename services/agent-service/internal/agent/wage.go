package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

// WageFormID is a 96-bit hex identifier for WageForm records.
type WageFormID string

func (id WageFormID) IsZero() bool { return id == "" }

func NewWageFormID() WageFormID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return WageFormID(hex.EncodeToString(b))
}

// WageLabourValue is the Ch. 19 money-form of labour-power value: the daily
// pence equivalent and the necessary labour-time it represents.
type WageLabourValue struct {
	DailyPence       Pence         `json:"daily_pence"`
	NecessaryMinutes LabourMinutes `json:"necessary_minutes"`
}

// Wage is the money form paid per working day — the price of labour as it
// appears to workers, capitalists, and economists.
type Wage struct {
	DailyPence      Pence `json:"daily_pence"`
	WorkingDayHours int64 `json:"working_day_hours"`
}

// WageAppearance models the ideological inversion Ch. 19 critiques: wages
// appear to pay for all hours worked, concealing that only necessary labour
// is compensated.
type WageAppearance struct {
	PaidHours   int64 `json:"paid_hours"`
	UnpaidHours int64 `json:"unpaid_hours"`
}

// LabourDecomposition reveals the true paid/unpaid split beneath the wage form.
type LabourDecomposition struct {
	PaidMinutes   LabourMinutes `json:"paid_minutes"`
	UnpaidMinutes LabourMinutes `json:"unpaid_minutes"`
}

// WageForm is the persisted record linking an agent to their wage and
// labour-power value at a point in time.
type WageForm struct {
	ID               WageFormID      `json:"id"`
	AgentID          AgentID         `json:"agent_id"`
	Wage             Wage            `json:"wage"`
	LabourPowerValue WageLabourValue `json:"labour_power_value"`
	CreatedAt        time.Time       `json:"created_at"`
}

// HourlyWage returns the apparent price of one hour of labour — the form in
// which workers and economists reason about the cost of labour.
func HourlyWage(w Wage) Pence {
	if w.WorkingDayHours == 0 {
		return 0
	}
	return Pence(int64(w.DailyPence) / w.WorkingDayHours)
}

// Decompose reveals the true paid/unpaid split beneath the wage form.
// PaidMinutes equals the necessary labour reproducing labour-power value;
// UnpaidMinutes is the surplus extracted by capital.
func Decompose(w Wage, lpv WageLabourValue) LabourDecomposition {
	total := LabourMinutes(w.WorkingDayHours * 60)
	return LabourDecomposition{
		PaidMinutes:   lpv.NecessaryMinutes,
		UnpaidMinutes: total - lpv.NecessaryMinutes,
	}
}

// Appearance models the ideological inversion: the wage form presents all
// hours as paid, making surplus labour invisible.
func Appearance(w Wage) WageAppearance {
	return WageAppearance{
		PaidHours:   w.WorkingDayHours,
		UnpaidHours: 0,
	}
}

func (wf WageForm) Validate() error {
	if wf.AgentID.IsZero() {
		return errors.New("wage_form: agent_id is required")
	}
	if wf.Wage.DailyPence <= 0 {
		return errors.New("wage_form: daily_pence must be positive")
	}
	if wf.Wage.WorkingDayHours <= 0 {
		return errors.New("wage_form: working_day_hours must be positive")
	}
	if wf.LabourPowerValue.DailyPence <= 0 {
		return errors.New("wage_form: labour_power_value daily_pence must be positive")
	}
	if wf.LabourPowerValue.NecessaryMinutes <= 0 {
		return errors.New("wage_form: necessary_minutes must be positive")
	}
	return nil
}
