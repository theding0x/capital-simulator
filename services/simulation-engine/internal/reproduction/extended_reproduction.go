// Package reproduction — Vol. II Ch. 21: Accumulation and Reproduction on an Extended Scale.
//
// Marx's two-department extended reproduction scheme shows how accumulation
// of part of the surplus-value enables society to reproduce itself on a growing
// scale. Unlike simple reproduction (Ch. 20), part of s is converted to
// additional c and v, expanding the means of production in Dept I faster than
// in Dept II.
//
// The keystone equation for extended reproduction is:
//
//	I(v + Δv + consumed_s_I) == II(c + Δc)
//
// Marx's canonical fixture (period "Year I"):
//
//	Dept I:  4000c + 1000v + 1000s, 50% accumulation rate → ΔcI=400, ΔvI=100
//	Dept II: 1500c +  750v +  750s
//
// Balance check: I.v(1000) + ΔvI(100) + consumed_I(500) = 1600 = II.c(1500) + ΔcII(100). Balanced.
//
// All monetary values are in pence (1 GBP = 240 pence).
// Accumulation rates are in basis points (10000 bps = 100%).
package reproduction

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

// --- Sentinel errors --------------------------------------------------------

var (
	// ErrExtendedImbalance is returned when the extended reproduction equation
	// is not satisfied.
	ErrExtendedImbalance = errors.New("reproduction: extended reproduction equation not satisfied")

	// ErrNegativeAccumulationRate is returned when an accumulation rate < 0 is
	// supplied.
	ErrNegativeAccumulationRate = errors.New("reproduction: accumulation rate must be non-negative")

	// ErrAccumulationRateExceeds100Pct is returned when an accumulation rate
	// exceeds 10000 bps (100%).
	ErrAccumulationRateExceeds100Pct = errors.New("reproduction: accumulation rate exceeds 100% (10000 bps)")
)

// --- ID types ---------------------------------------------------------------

// ExtendedReproductionSchemeID identifies an ExtendedReproductionScheme record.
type ExtendedReproductionSchemeID string

// NewExtendedReproductionSchemeID returns a 96-bit random hex identifier.
func NewExtendedReproductionSchemeID() ExtendedReproductionSchemeID {
	return ExtendedReproductionSchemeID(newExtendedReproHexID())
}

// ReinvestmentID identifies a Reinvestment record.
type ReinvestmentID string

// NewReinvestmentID returns a 96-bit random hex identifier.
func NewReinvestmentID() ReinvestmentID {
	return ReinvestmentID(newExtendedReproHexID())
}

// MultiPeriodSchemeID identifies a MultiPeriodScheme record.
type MultiPeriodSchemeID string

// NewMultiPeriodSchemeID returns a 96-bit random hex identifier.
func NewMultiPeriodSchemeID() MultiPeriodSchemeID {
	return MultiPeriodSchemeID(newExtendedReproHexID())
}

// CompositionShiftID identifies a CompositionShift record.
type CompositionShiftID string

// NewCompositionShiftID returns a 96-bit random hex identifier.
func NewCompositionShiftID() CompositionShiftID {
	return CompositionShiftID(newExtendedReproHexID())
}

// ExtendedMoneyRequirementID identifies an ExtendedMoneyRequirement record.
type ExtendedMoneyRequirementID string

// NewExtendedMoneyRequirementID returns a 96-bit random hex identifier.
func NewExtendedMoneyRequirementID() ExtendedMoneyRequirementID {
	return ExtendedMoneyRequirementID(newExtendedReproHexID())
}

// DepartmentIGrowthLeadID identifies a DepartmentIGrowthLead record.
type DepartmentIGrowthLeadID string

// NewDepartmentIGrowthLeadID returns a 96-bit random hex identifier.
func NewDepartmentIGrowthLeadID() DepartmentIGrowthLeadID {
	return DepartmentIGrowthLeadID(newExtendedReproHexID())
}

// --- Structs ----------------------------------------------------------------

// ExtendedReproductionScheme is the root of the extended reproduction model.
// It holds references to both departments and tracks the accumulation rates
// that drive the expansion of production in each department.
//
// IsBalanced is true when:
//
//	I(v + Δv + consumed_s_I) == II(c + Δc)
//
// TickCount advances each time TickExtendedScheme is called.
type ExtendedReproductionScheme struct {
	ID                    ExtendedReproductionSchemeID `json:"id"`
	Period                string                       `json:"period"`
	DepartmentI           *DepartmentalCapital         `json:"department_i,omitempty"`
	DepartmentII          *DepartmentalCapital         `json:"department_ii,omitempty"`
	AccumulationRateIBps  int64                        `json:"accumulation_rate_i_bps"`
	AccumulationRateIIBps int64                        `json:"accumulation_rate_ii_bps"`
	IsBalanced            bool                         `json:"is_balanced"`
	TickCount             int64                        `json:"tick_count"`
}

// Reinvestment records how surplus-value was split into Δc, Δv, and consumed
// surplus for one department in one accumulation cycle.
type Reinvestment struct {
	ID                   ReinvestmentID               `json:"id"`
	SchemeID             ExtendedReproductionSchemeID `json:"scheme_id"`
	Department           CapitalDepartment            `json:"department"`
	DeltaConstantPence   int64                        `json:"delta_constant_pence"`
	DeltaVariablePence   int64                        `json:"delta_variable_pence"`
	ConsumedSurplusPence int64                        `json:"consumed_surplus_pence"`
	CycleNumber          int64                        `json:"cycle_number"`
}

// MultiPeriodScheme links a sequence of ExtendedReproductionSchemes representing
// successive accumulation periods, enabling cross-period analysis of organic
// composition and growth leads.
type MultiPeriodScheme struct {
	ID        MultiPeriodSchemeID            `json:"id"`
	Label     string                         `json:"label"`
	SchemeIDs []ExtendedReproductionSchemeID `json:"scheme_ids"`
}

// CompositionShift records the organic composition c/v (in basis points) for
// both departments and the social average across two consecutive periods.
type CompositionShift struct {
	ID                   CompositionShiftID  `json:"id"`
	MultiPeriodID        MultiPeriodSchemeID `json:"multi_period_id"`
	CycleNumber          int64               `json:"cycle_number"`
	DeptICompositionBps  int64               `json:"dept_i_composition_bps"`
	DeptIICompositionBps int64               `json:"dept_ii_composition_bps"`
	SocialCompositionBps int64               `json:"social_composition_bps"`
}

// ExtendedMoneyRequirement captures the additional money needed to circulate
// a growing commodity mass during an accumulation cycle.
type ExtendedMoneyRequirement struct {
	ID                   ExtendedMoneyRequirementID   `json:"id"`
	SchemeID             ExtendedReproductionSchemeID `json:"scheme_id"`
	BaseMoneyPence       int64                        `json:"base_money_pence"`
	AdditionalMoneyPence int64                        `json:"additional_money_pence"`
	TotalMoneyPence      int64                        `json:"total_money_pence"`
}

// DepartmentIGrowthLead records whether Dept I grew faster than Dept II
// between two consecutive periods, confirming Marx's thesis that means-of-
// production production must outpace consumption-goods production for
// accumulation to proceed.
type DepartmentIGrowthLead struct {
	ID              DepartmentIGrowthLeadID `json:"id"`
	MultiPeriodID   MultiPeriodSchemeID     `json:"multi_period_id"`
	CycleNumber     int64                   `json:"cycle_number"`
	DeptIGrowthBps  int64                   `json:"dept_i_growth_bps"`
	DeptIIGrowthBps int64                   `json:"dept_ii_growth_bps"`
	LeadConfirmed   bool                    `json:"lead_confirmed"`
}

// --- Pure domain functions --------------------------------------------------

// ComputeExtendedBalance returns true iff the extended reproduction keystone
// equation holds:
//
//	I.v + reinvestI.DeltaVariablePence + reinvestI.ConsumedSurplusPence
//	  == II.c + reinvestII.DeltaConstantPence
//
// This ensures Dept II can absorb all of Dept I's wage and consumption
// output while replacing its worn-out constant capital.
func ComputeExtendedBalance(deptI, deptII DepartmentalCapital, reinvestI, reinvestII Reinvestment) bool {
	lhs := deptI.VariablePence + reinvestI.DeltaVariablePence + reinvestI.ConsumedSurplusPence
	rhs := deptII.ConstantPence + reinvestII.DeltaConstantPence
	return lhs == rhs
}

// ComputeReinvestment splits dept.SurplusPence into Δc, Δv, and consumed
// surplus according to accumulationRateBps (basis points; 5000 = 50%).
//
// The organic composition of the department determines the c:v split of the
// reinvested amount:
//
//	reinvestedTotal = SurplusPence * accumulationRateBps / 10000
//	DeltaConstantPence = reinvestedTotal * ConstantPence / TotalPence  (integer division)
//	DeltaVariablePence = reinvestedTotal - DeltaConstantPence
//	ConsumedSurplusPence = SurplusPence - reinvestedTotal
//
// Returns ErrNegativeAccumulationRate if rate < 0.
// Returns ErrAccumulationRateExceeds100Pct if rate > 10000.
func ComputeReinvestment(dept DepartmentalCapital, accumulationRateBps int64) (Reinvestment, error) {
	if accumulationRateBps < 0 {
		return Reinvestment{}, ErrNegativeAccumulationRate
	}
	if accumulationRateBps > 10000 {
		return Reinvestment{}, ErrAccumulationRateExceeds100Pct
	}
	reinvestedTotal := dept.SurplusPence * accumulationRateBps / 10000
	var deltaC int64
	if cv := dept.ConstantPence + dept.VariablePence; cv > 0 {
		deltaC = reinvestedTotal * dept.ConstantPence / cv
	}
	deltaV := reinvestedTotal - deltaC
	consumed := dept.SurplusPence - reinvestedTotal
	return Reinvestment{
		ID:                   NewReinvestmentID(),
		Department:           dept.Department,
		DeltaConstantPence:   deltaC,
		DeltaVariablePence:   deltaV,
		ConsumedSurplusPence: consumed,
	}, nil
}

// ComputeNextPeriodCapital advances dept by one accumulation cycle using r.
//
//	New constant = ConstantPence + DeltaConstantPence
//	New variable = VariablePence + DeltaVariablePence
//	New surplus at same rate: surplus_rate_bps = SurplusPence * 10000 / VariablePence
//	  NewSurplus = NewVariable * surplus_rate_bps / 10000
//	New total = NewConstant + NewVariable + NewSurplus
//
// The ID and SchemeID are preserved from dept.
func ComputeNextPeriodCapital(dept DepartmentalCapital, r Reinvestment) DepartmentalCapital {
	newC := dept.ConstantPence + r.DeltaConstantPence
	newV := dept.VariablePence + r.DeltaVariablePence
	var surplusRateBps int64
	if dept.VariablePence > 0 {
		surplusRateBps = dept.SurplusPence * 10000 / dept.VariablePence
	}
	newS := newV * surplusRateBps / 10000
	newTotal := newC + newV + newS
	return DepartmentalCapital{
		ID:            dept.ID,
		SchemeID:      dept.SchemeID,
		Department:    dept.Department,
		ConstantPence: newC,
		VariablePence: newV,
		SurplusPence:  newS,
		TotalPence:    newTotal,
	}
}

// ComputeCompositionBps returns the organic composition c/v in basis points.
// Returns 0 if VariablePence == 0.
func ComputeCompositionBps(dept DepartmentalCapital) int64 {
	if dept.VariablePence == 0 {
		return 0
	}
	return dept.ConstantPence * 10000 / dept.VariablePence
}

// ComputeGrowthBps returns (curr.TotalPence - prev.TotalPence) * 10000 / prev.TotalPence
// in basis points. Returns 0 if prev.TotalPence == 0.
func ComputeGrowthBps(prev, curr DepartmentalCapital) int64 {
	if prev.TotalPence == 0 {
		return 0
	}
	return (curr.TotalPence - prev.TotalPence) * 10000 / prev.TotalPence
}

// ComputeExtendedMoneyRequirement computes the money-stock needed to circulate
// the commodity mass in an accumulation cycle.
//
//	AdditionalMoney = (DeptI.TotalPence + DeptII.TotalPence) * (AccumulationRateIBps + AccumulationRateIIBps) / 20000
//	TotalMoney = baseMoneyPence + AdditionalMoney
//
// Returns ErrExtendedImbalance if s.DepartmentI or s.DepartmentII is nil.
func ComputeExtendedMoneyRequirement(s ExtendedReproductionScheme, baseMoneyPence int64) (ExtendedMoneyRequirement, error) {
	if s.DepartmentI == nil || s.DepartmentII == nil {
		return ExtendedMoneyRequirement{}, ErrExtendedImbalance
	}
	totalCapital := s.DepartmentI.TotalPence + s.DepartmentII.TotalPence
	additionalMoney := totalCapital * (s.AccumulationRateIBps + s.AccumulationRateIIBps) / 20000
	return ExtendedMoneyRequirement{
		ID:                   NewExtendedMoneyRequirementID(),
		SchemeID:             s.ID,
		BaseMoneyPence:       baseMoneyPence,
		AdditionalMoneyPence: additionalMoney,
		TotalMoneyPence:      baseMoneyPence + additionalMoney,
	}, nil
}

// --- helpers ----------------------------------------------------------------

// newExtendedReproHexID returns a 24-character lowercase hex string (96 bits).
func newExtendedReproHexID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(b)
}
