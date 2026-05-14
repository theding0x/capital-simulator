package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

// WorkingSessionID is a 96-bit hex identifier for WorkingSession records.
type WorkingSessionID string

func (id WorkingSessionID) IsZero() bool { return id == "" }

func NewWorkingSessionID() WorkingSessionID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return WorkingSessionID(hex.EncodeToString(b))
}

// DailyLabourPowerValue is the money expression of the daily value of labour-power
// — the denominator of the Ch. 20 hourly-price fraction.
type DailyLabourPowerValue struct {
	Pence int64 `json:"pence"`
}

// HourlyPriceOfLabour is the price of one hour of labour expressed as an exact
// rational fraction (DailyLabourPowerValue / WorkingDayMinutes÷60). Stored without
// rounding so the inverse relationship between day length and hourly price is exact.
type HourlyPriceOfLabour struct {
	Numerator   int64 `json:"numerator"`
	Denominator int64 `json:"denominator"`
}

// AsFloat returns the decimal approximation for display.
func (h HourlyPriceOfLabour) AsFloat() float64 {
	if h.Denominator == 0 {
		return 0
	}
	return float64(h.Numerator) / float64(h.Denominator)
}

// WorkingDayMinutes is the total duration of a time-wage working session in minutes.
type WorkingDayMinutes struct {
	Minutes LabourMinutes `json:"minutes"`
}

// OvertimeHours is the number of hours worked beyond the normal working day.
type OvertimeHours struct {
	Hours int64 `json:"hours"`
}

// OvertimeRatePence is the extra pay per overtime hour (often above the normal
// hourly rate, though Marx shows the capitalist still extracts unpaid labour).
type OvertimeRatePence struct {
	Pence int64 `json:"pence"`
}

// WagePeriod indicates whether the base daily value is expressed as a daily or
// weekly wage.
type WagePeriod string

const (
	WagePeriodDaily  WagePeriod = "daily"
	WagePeriodWeekly WagePeriod = "weekly"
)

// NominalWage is the total money received for a working session (normal pay plus
// overtime). It is the magnitude that appears to the worker as their wage.
type NominalWage struct {
	Pence int64 `json:"pence"`
}

// WorkingSession is the persisted record grouping all inputs for a time-wage
// calculation tied to an agent for a single period.
type WorkingSession struct {
	ID                    WorkingSessionID      `json:"id"`
	AgentID               AgentID               `json:"agent_id"`
	DailyLabourPowerValue DailyLabourPowerValue `json:"daily_labour_power_value"`
	WorkingDayMinutes     WorkingDayMinutes     `json:"working_day_minutes"`
	OvertimeHours         OvertimeHours         `json:"overtime_hours"`
	OvertimeRatePence     OvertimeRatePence     `json:"overtime_rate_pence"`
	WagePeriod            WagePeriod            `json:"wage_period"`
	CreatedAt             time.Time             `json:"created_at"`
}

// ErrTimeWageInvalidPeriod is returned when wage_period is neither "daily" nor "weekly".
var ErrTimeWageInvalidPeriod = errors.New("time_wage: wage_period must be 'daily' or 'weekly'")

func (s WorkingSession) Validate() error {
	if s.AgentID.IsZero() {
		return errors.New("time_wage: agent_id is required")
	}
	if s.DailyLabourPowerValue.Pence <= 0 {
		return errors.New("time_wage: daily_labour_power_value.pence must be positive")
	}
	if s.WorkingDayMinutes.Minutes <= 0 {
		return errors.New("time_wage: working_day_minutes.minutes must be positive")
	}
	if s.OvertimeHours.Hours < 0 {
		return errors.New("time_wage: overtime_hours.hours cannot be negative")
	}
	if s.OvertimeRatePence.Pence < 0 {
		return errors.New("time_wage: overtime_rate_pence.pence cannot be negative")
	}
	if s.WagePeriod != WagePeriodDaily && s.WagePeriod != WagePeriodWeekly {
		return ErrTimeWageInvalidPeriod
	}
	return nil
}

// ComputeHourlyPrice returns the price of one hour of labour as an exact rational
// fraction: Numerator = daily value in pence, Denominator = working-day hours.
// Never rounds — the fraction is the invariant form Marx uses to show that
// lengthening the day lowers the hourly price even when nominal pay is unchanged.
func ComputeHourlyPrice(v DailyLabourPowerValue, m WorkingDayMinutes) HourlyPriceOfLabour {
	return HourlyPriceOfLabour{Numerator: v.Pence, Denominator: int64(m.Minutes) / 60}
}

// ComputeSessionWage returns the total nominal money paid for the session:
// (floor hourly rate × normal hours) + (overtime hours × overtime rate).
// Uses integer arithmetic throughout — no rounding at this layer.
func ComputeSessionWage(s WorkingSession) NominalWage {
	hours := int64(s.WorkingDayMinutes.Minutes) / 60
	hourlyRate := s.DailyLabourPowerValue.Pence / hours
	normalPay := hourlyRate * hours
	overtimePay := s.OvertimeHours.Hours * s.OvertimeRatePence.Pence
	return NominalWage{Pence: normalPay + overtimePay}
}
