package commodity

import (
	"errors"
	"math"
)

// ConstantKind classifies the means of production.
type ConstantKind string

const (
	KindInstrument  ConstantKind = "instrument"
	KindRawMaterial ConstantKind = "raw_material"
	KindAuxiliary   ConstantKind = "auxiliary"
)

// ConstantCapital is a means of production whose value is preserved and
// transferred to the product — not created anew [§2].
type ConstantCapital struct {
	OriginalValue   LabourMinutes `json:"original_value"`
	Kind            ConstantKind  `json:"kind"`
	ServiceLifeDays int64         `json:"service_life_days"`
}

// Validate reports if c is well-formed.
func (c ConstantCapital) Validate() error {
	if c.OriginalValue < 0 {
		return errors.New("commodity: constant_capital original_value must be non-negative")
	}
	if c.Kind == "" {
		return errors.New("commodity: constant_capital kind is required")
	}
	if c.ServiceLifeDays < 0 {
		return errors.New("commodity: constant_capital service_life_days must be non-negative")
	}
	return nil
}

// WearFraction is the portion of an instrument's value transferred per
// production period.
type WearFraction float64

// WearFractionFor returns the fraction of c's value transferred in one run.
// For raw materials and auxiliaries (ServiceLifeDays == 0) the full value transfers.
func WearFractionFor(c ConstantCapital) WearFraction {
	if c.ServiceLifeDays <= 0 {
		return WearFraction(1.0)
	}
	return WearFraction(1.0 / float64(c.ServiceLifeDays))
}

// TransferredValue returns the LabourMinutes of c's value transferred to the
// product in one period, given the wear fraction.
func TransferredValue(c ConstantCapital, fraction WearFraction) LabourMinutes {
	return LabourMinutes(math.Round(float64(c.OriginalValue) * float64(fraction)))
}

// NewValue returns the new value created by living labour in labourTime.
// At uniform skill level = 1 (abstract labour), the reduction is the identity.
func NewValue(labourTime LabourMinutes) LabourMinutes { return labourTime }

// VariableCapital is the portion of capital representing purchased labour-power.
type VariableCapital struct {
	WageValue  LabourMinutes `json:"wage_value"`
	WorkingDay LabourMinutes `json:"working_day"`
}

// Validate reports if v is well-formed.
func (v VariableCapital) Validate() error {
	if v.WageValue < 0 {
		return errors.New("commodity: variable_capital wage_value must be non-negative")
	}
	if v.WorkingDay <= 0 {
		return errors.New("commodity: variable_capital working_day must be positive")
	}
	if v.WageValue > v.WorkingDay {
		return errors.New("commodity: variable_capital wage_value cannot exceed working_day")
	}
	return nil
}

// SurplusLabourFrom returns the unpaid portion of the working day.
func (v VariableCapital) SurplusLabourFrom() LabourMinutes { return v.WorkingDay - v.WageValue }

// ProductValue decomposes the total product value into c + v + s.
type ProductValue struct {
	Constant LabourMinutes `json:"constant"`
	Variable LabourMinutes `json:"variable"`
	Surplus  LabourMinutes `json:"surplus"`
}

// Total returns c + v + s.
func (pv ProductValue) Total() LabourMinutes {
	return pv.Constant + pv.Variable + pv.Surplus
}

// DecomposeProductValue builds the ProductValue from constant-capital inputs and
// one variable-capital input.
func DecomposeProductValue(constants []ConstantCapital, v VariableCapital) ProductValue {
	var constantTotal LabourMinutes
	for _, c := range constants {
		fraction := WearFractionFor(c)
		constantTotal += TransferredValue(c, fraction)
	}
	surplus := v.SurplusLabourFrom()
	return ProductValue{
		Constant: constantTotal,
		Variable: v.WageValue,
		Surplus:  surplus,
	}
}

// CapitalComposition records the ratio of constant to variable capital.
type CapitalComposition struct {
	Constant LabourMinutes `json:"constant"`
	Variable LabourMinutes `json:"variable"`
}

// Ratio returns c/v. Returns 0 if Variable is 0.
func (cc CapitalComposition) Ratio() float64 {
	if cc.Variable == 0 {
		return 0
	}
	return float64(cc.Constant) / float64(cc.Variable)
}
