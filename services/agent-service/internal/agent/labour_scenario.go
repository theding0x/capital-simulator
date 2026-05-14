package agent

import "errors"

// LabourScenario groups the three variable factors Marx isolates in
// Capital Vol. I, Ch. 17: the length of the working day, the intensity
// of labour, and the productiveness of labour. The WorkingDay carries
// the partition into necessary and surplus labour (reused from Ch. 10);
// Intensity and Productivity are intensive/multiplicative magnitudes that
// scale value creation and shift the partition respectively.
//
// LabourScenario is purely an input shape — the persistent entity for the
// working day lives in WorkingDay itself. ComputeOutcome on this scenario
// is stateless and pure.
type LabourScenario struct {
	WorkingDay   WorkingDay         `json:"working_day"`
	Intensity    LabourIntensity    `json:"intensity"`
	Productivity LabourProductivity `json:"productivity"`
}

// ScenarioOutcome is the result of evaluating a LabourScenario: the daily
// value created, the partition between necessary and surplus labour, the
// implied value of labour-power (= necessary labour-time), and the rate
// of surplus-value s/v.
type ScenarioOutcome struct {
	DailyValueMinutes       LabourMinutes `json:"daily_value_minutes"`
	NecessaryLabourMinutes  LabourMinutes `json:"necessary_labour_minutes"`
	SurplusLabourMinutes    LabourMinutes `json:"surplus_labour_minutes"`
	LabourPowerValueMinutes LabourMinutes `json:"labour_power_value_minutes"`
	RateOfSurplusValue      float64       `json:"rate_of_surplus_value"`
}

var ErrScenarioNecessaryExceedsDay = errors.New("agent: necessary labour cannot exceed working day total")

func (s LabourScenario) Validate() error {
	if err := s.WorkingDay.Validate(); err != nil {
		return err
	}
	if err := s.Intensity.Validate(); err != nil {
		return err
	}
	if err := s.Productivity.Validate(); err != nil {
		return err
	}
	if int64(s.WorkingDay.NecessaryLabourMinutes) > s.WorkingDay.TotalMinutes() {
		return ErrScenarioNecessaryExceedsDay
	}
	return nil
}

// DailyValueCreated returns the labour-value embodied in a single
// working day, scaled by intensity. At normal intensity (Factor = 1.0)
// value-in-minutes equals labour-time-in-minutes; at higher intensity a
// minute of expended labour creates proportionally more value.
//
// This is Marx's §2 generalisation: the "magnitude of value created"
// depends jointly on the *duration* and the *intensity* of labour.
func DailyValueCreated(wd WorkingDay, li LabourIntensity) LabourMinutes {
	return LabourMinutes(float64(wd.TotalMinutes()) * li.Factor)
}

// ComputeOutcome partitions a LabourScenario into its outcome
// components. It is stateless and pure: the working day's
// necessary/surplus partition is taken as given, intensity scales the
// daily value, and productivity feeds the laws (not the partition
// directly — productivity changes have already been baked into the
// caller's choice of NecessaryLabourMinutes).
func ComputeOutcome(s LabourScenario) ScenarioOutcome {
	dayMinutes := s.WorkingDay.TotalMinutes()
	necessary := LabourMinutes(s.WorkingDay.NecessaryLabourMinutes)
	surplus := LabourMinutes(dayMinutes) - necessary

	intensity := s.Intensity
	if intensity.Factor <= 0 {
		intensity = NormalLabourIntensity
	}
	dailyValue := DailyValueCreated(s.WorkingDay, intensity)

	rate := 0.0
	if necessary > 0 {
		rate = float64(surplus) / float64(necessary)
	}
	return ScenarioOutcome{
		DailyValueMinutes:       dailyValue,
		NecessaryLabourMinutes:  necessary,
		SurplusLabourMinutes:    surplus,
		LabourPowerValueMinutes: necessary,
		RateOfSurplusValue:      rate,
	}
}

// LawConstantDailyValue encodes Ch. 17 §1 Law 1: at normal intensity,
// the daily value created by a working day of fixed length is
// independent of changes in productivity. (Productivity redistributes
// the partition between necessary and surplus labour without altering
// the total value created.)
//
// The productivityFactor argument is taken explicitly so the call site
// reads as a property check; it is used only to document the law's
// scope, since DailyValueCreated does not depend on productivity.
func LawConstantDailyValue(day WorkingDay, intensity LabourIntensity, productivityFactor float64) bool {
	if !intensity.IsNormal() {
		return false
	}
	_ = productivityFactor
	return DailyValueCreated(day, intensity) == LabourMinutes(day.TotalMinutes())
}

// LawInverseRelation encodes Ch. 17 §1 Law 2: when productivity changes
// (with working day and intensity held constant), the value of
// labour-power and the magnitude of surplus-value move in opposite
// directions. A no-change pair is trivially compatible with the law.
func LawInverseRelation(before, after ScenarioOutcome) bool {
	lpDelta := after.LabourPowerValueMinutes - before.LabourPowerValueMinutes
	svDelta := after.SurplusLabourMinutes - before.SurplusLabourMinutes
	if lpDelta == 0 && svDelta == 0 {
		return true
	}
	return (lpDelta > 0 && svDelta < 0) || (lpDelta < 0 && svDelta > 0)
}
