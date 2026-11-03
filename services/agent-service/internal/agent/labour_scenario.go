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

	// Pressure is an optional industrial-reserve-army overlay (Ch. 25 → Ch. 17):
	// when set and active, ComputeOutcome reports the wage the reserve army's
	// competition pushes the worker down to, below the value of labour-power.
	// Nil means full employment — no compression.
	Pressure *ReserveArmyPressure `json:"reserve_army_pressure,omitempty"`
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

	// ReserveArmyPressureBP and CompressedWageMinutes are populated only when the
	// scenario carries an active ReserveArmyPressure. CompressedWageMinutes is the
	// value of labour-power after the reserve army's downward pull; it is
	// <= LabourPowerValueMinutes, the gap being the wage pushed below its value.
	ReserveArmyPressureBP int64         `json:"reserve_army_pressure_bp,omitempty"`
	CompressedWageMinutes LabourMinutes `json:"compressed_wage_minutes,omitempty"`
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
	if s.Pressure != nil {
		if err := s.Pressure.Validate(); err != nil {
			return err
		}
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
// necessary/surplus partition is the baseline at normal productivity,
// productivity then divides the necessary labour-time (Ch. 17 §1 Case I:
// the day is constant, so what necessary labour sheds becomes surplus
// labour), and intensity scales the daily value.
func ComputeOutcome(s LabourScenario) ScenarioOutcome {
	dayMinutes := s.WorkingDay.TotalMinutes()
	necessary := s.Productivity.ApplyToNecessaryLabour(LabourMinutes(s.WorkingDay.NecessaryLabourMinutes))
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
	outcome := ScenarioOutcome{
		DailyValueMinutes:       dailyValue,
		NecessaryLabourMinutes:  necessary,
		SurplusLabourMinutes:    surplus,
		LabourPowerValueMinutes: necessary,
		RateOfSurplusValue:      rate,
	}
	// Ch. 25 → Ch. 17: an active reserve army bids the price of labour down
	// below the value of labour-power (= necessary labour). The partition above
	// is untouched — only the wage actually paid is compressed.
	if s.Pressure != nil && s.Pressure.IsActive() {
		outcome.ReserveArmyPressureBP = s.Pressure.PressureBP()
		outcome.CompressedWageMinutes = s.Pressure.CompressMinutes(necessary)
	}
	return outcome
}

// LawConstantDailyValue encodes Ch. 17 §1 Law 1: at normal intensity,
// the daily value created by a working day of fixed length is
// independent of changes in productivity. (Productivity redistributes
// the partition between necessary and surplus labour without altering
// the total value created.)
//
// The check is derived from the domain: the scenario is evaluated at the
// given productivity and at normal productivity, and the daily values
// must agree and equal the working-day length.
func LawConstantDailyValue(day WorkingDay, intensity LabourIntensity, productivityFactor float64) bool {
	if !intensity.IsNormal() {
		return false
	}
	changed := ComputeOutcome(LabourScenario{WorkingDay: day, Intensity: intensity, Productivity: LabourProductivity{Factor: productivityFactor}})
	baseline := ComputeOutcome(LabourScenario{WorkingDay: day, Intensity: intensity, Productivity: NormalLabourProductivity})
	return changed.DailyValueMinutes == baseline.DailyValueMinutes &&
		changed.DailyValueMinutes == LabourMinutes(day.TotalMinutes())
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
