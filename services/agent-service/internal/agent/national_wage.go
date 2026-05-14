package agent

import "errors"

// CountryCode is an ISO 3166-1 alpha-2 country identifier.
type CountryCode string

// NationalIntensity is the intensity of labour in a country relative to the
// international average (1.0 = average). Greater intensity produces more value
// in the same duration.
type NationalIntensity struct {
	CountryCode CountryCode `json:"country_code"`
	Factor      float64     `json:"factor"`
}

// DayWage is the nominal day wage as paid in the domestic market.
type DayWage struct {
	CountryCode       CountryCode `json:"country_code"`
	NominalPence      int64       `json:"nominal_pence"`
	WorkingDayMinutes int64       `json:"working_day_minutes"`
}

// StandardisedWage is a day wage reduced to a common working-day length for
// cross-national comparison.
type StandardisedWage struct {
	CountryCode CountryCode `json:"country_code"`
	Amount      int64       `json:"amount"`
}

// RelativeLabourPrice is the price of labour relative to the value it creates.
// Inversely related to national intensity: higher intensity → lower relative
// price despite potentially higher nominal wage.
type RelativeLabourPrice struct {
	CountryCode CountryCode `json:"country_code"`
	Ratio       float64     `json:"ratio"`
}

// SpindleRatio is the number of spindles per worker — a productivity proxy
// from Redgrave's 1866 factory inspectorate data.
type SpindleRatio struct {
	CountryCode       CountryCode `json:"country_code"`
	SpindlesPerWorker int64       `json:"spindles_per_worker"`
}

// WageComparison aggregates all per-country wage metrics for cross-national
// comparison.
type WageComparison struct {
	Countries         []CountryCode         `json:"countries"`
	DayWages          []DayWage             `json:"day_wages"`
	StandardisedWages []StandardisedWage    `json:"standardised_wages"`
	RelativePrices    []RelativeLabourPrice `json:"relative_prices"`
	SpindleRatios     []SpindleRatio        `json:"spindle_ratios"`
}

var (
	ErrNationalIntensityFactor       = errors.New("national_intensity: factor must be positive")
	ErrDayWageNominalPence           = errors.New("day_wage: nominal_pence must be positive")
	ErrDayWageWorkingDayMinutes      = errors.New("day_wage: working_day_minutes must be positive")
	ErrSpindleRatioSpindlesPerWorker = errors.New("spindle_ratio: spindles_per_worker must be positive")
)

func (ni NationalIntensity) Validate() error {
	if ni.Factor <= 0 {
		return ErrNationalIntensityFactor
	}
	return nil
}

func (w DayWage) Validate() error {
	if w.NominalPence <= 0 {
		return ErrDayWageNominalPence
	}
	if w.WorkingDayMinutes <= 0 {
		return ErrDayWageWorkingDayMinutes
	}
	return nil
}

func (sr SpindleRatio) Validate() error {
	if sr.SpindlesPerWorker <= 0 {
		return ErrSpindleRatioSpindlesPerWorker
	}
	return nil
}

// StandardiseWage reduces a day wage to a uniform reference working-day length.
// Invariant: Amount == NominalPence * referenceDayMinutes / WorkingDayMinutes.
func StandardiseWage(w DayWage, referenceDayMinutes int64) StandardisedWage {
	return StandardisedWage{
		CountryCode: w.CountryCode,
		Amount:      w.NominalPence * referenceDayMinutes / w.WorkingDayMinutes,
	}
}

// ComputeRelativePrice returns the price of labour relative to the value it
// produces: nominal wage per unit of intensity (NominalPence / Factor).
// Higher intensity reduces the ratio even when nominal wages are higher —
// encoding Marx's observation that "wages are virtually lower to the capitalist,
// though higher to the operative" in high-intensity nations (Cowell, 1833).
func ComputeRelativePrice(w DayWage, ni NationalIntensity) RelativeLabourPrice {
	return RelativeLabourPrice{
		CountryCode: w.CountryCode,
		Ratio:       float64(w.NominalPence) / ni.Factor,
	}
}
