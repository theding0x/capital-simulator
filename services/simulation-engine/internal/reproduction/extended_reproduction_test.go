// Package reproduction — Vol. II Ch. 21: Accumulation and Reproduction on an Extended Scale.
//
// Tests for pure domain functions and ID constructors introduced in Ch. 21.
// Marx's canonical Ch. 21 fixture (in units of 1 000 pence):
//
//	Dept I:  4 000c + 1 000v + 1 000s = 6 000  total  (50% accumulation → ΔcI=400, ΔvI=100)
//	Dept II: 1 500c +   750v +   750s = 3 000  total
//
// Keystone equation: I(v + ΔvI + consumed_sI) == II(c + ΔcII)
//
//	1000 + 100 + 500 = 1600 == 1500 + 100 = 1600 ✓
package reproduction

import (
	"errors"
	"testing"
)

// ch21 canonical fixture (pence — scaled 1000× from Marx's £ values)
const (
	ch21DeptIConstant  int64 = 4_000_000
	ch21DeptIVariable  int64 = 1_000_000
	ch21DeptISurplus   int64 = 1_000_000
	ch21DeptITotal     int64 = 6_000_000
	ch21DeptIIConstant int64 = 1_500_000
	ch21DeptIIVariable int64 = 750_000
	ch21DeptIISurplus  int64 = 750_000
	ch21DeptIITotal    int64 = 3_000_000
)

func newExtDeptI(schemeID ExtendedReproductionSchemeID) DepartmentalCapital {
	return DepartmentalCapital{
		ID:            NewDepartmentalCapitalID(),
		SchemeID:      SimpleReproductionSchemeID(schemeID),
		Department:    DepartmentI,
		ConstantPence: ch21DeptIConstant,
		VariablePence: ch21DeptIVariable,
		SurplusPence:  ch21DeptISurplus,
		TotalPence:    ch21DeptITotal,
	}
}

func newExtDeptII(schemeID ExtendedReproductionSchemeID) DepartmentalCapital {
	return DepartmentalCapital{
		ID:            NewDepartmentalCapitalID(),
		SchemeID:      SimpleReproductionSchemeID(schemeID),
		Department:    DepartmentII,
		ConstantPence: ch21DeptIIConstant,
		VariablePence: ch21DeptIIVariable,
		SurplusPence:  ch21DeptIISurplus,
		TotalPence:    ch21DeptIITotal,
	}
}

// --- ComputeExtendedBalance -------------------------------------------------

// TestComputeExtendedBalance_Balanced verifies the extended keystone equation
// for Marx's canonical Year I fixture:
//
//	I.v(1_000_000) + ΔvI(100_000) + consumed_sI(500_000) = 1_600_000
//	II.c(1_500_000) + ΔcII(100_000) = 1_600_000  ✓
func TestComputeExtendedBalance_Balanced(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dI := newExtDeptI(sid)
	dII := newExtDeptII(sid)

	// Dept I: 50% accumulation, c/(c+v) split → ΔcI=400_000, ΔvI=100_000, consumed=500_000
	rI, err := ComputeReinvestment(dI, 5000)
	if err != nil {
		t.Fatalf("ComputeReinvestment(Dept I, 50%%): %v", err)
	}
	// For balance: lhs = dI.v + rI.ΔvI + rI.consumed = dI.v + rI.surplus
	// rhs = dII.c + rII.ΔcII
	// We need rII.ΔcII = lhs - dII.c
	lhs := dI.VariablePence + rI.DeltaVariablePence + rI.ConsumedSurplusPence
	rhs := dII.ConstantPence
	rII := Reinvestment{
		SchemeID:             sid,
		Department:           DepartmentII,
		DeltaConstantPence:   lhs - rhs, // sets rII.ΔcII so equation is satisfied
		DeltaVariablePence:   0,
		ConsumedSurplusPence: dII.SurplusPence - (lhs - rhs),
	}

	if !ComputeExtendedBalance(dI, dII, rI, rII) {
		t.Fatalf("expected balanced: lhs=%d, rhs=%d",
			dI.VariablePence+rI.DeltaVariablePence+rI.ConsumedSurplusPence,
			dII.ConstantPence+rII.DeltaConstantPence)
	}
}

// TestComputeExtendedBalance_Imbalanced verifies the equation fails when
// ΔvI is perturbed so I(v+ΔvI+consumed_sI) != II(c+ΔcII).
func TestComputeExtendedBalance_Imbalanced(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dI := newExtDeptI(sid)
	dII := newExtDeptII(sid)

	// Break the balance: rI.consumed too low, rII.ΔcII=0
	rI := Reinvestment{
		SchemeID:             sid,
		Department:           DepartmentI,
		DeltaConstantPence:   0,
		DeltaVariablePence:   200_000,
		ConsumedSurplusPence: 100_000, // deliberately wrong — doesn't add up
	}
	rII := Reinvestment{
		SchemeID:             sid,
		Department:           DepartmentII,
		DeltaConstantPence:   0,
		DeltaVariablePence:   0,
		ConsumedSurplusPence: ch21DeptIISurplus,
	}

	if ComputeExtendedBalance(dI, dII, rI, rII) {
		t.Fatal("expected imbalanced")
	}
}

// --- ComputeReinvestment ----------------------------------------------------

// TestComputeReinvestment_HappyPath_50PctRate verifies a 50 % accumulation
// rate on Dept I (4000c+1000v+1000s) using Marx's c/(c+v) organic-composition split:
//
//	reinvestedTotal = 1_000_000 * 5000 / 10000 = 500_000
//	DeltaC = 500_000 * 4_000_000 / (4_000_000+1_000_000) = 400_000
//	DeltaV = 500_000 − 400_000 = 100_000
//	Consumed = 1_000_000 − 500_000 = 500_000
//
// Balance check: I.v(1_000_000) + ΔvI(100_000) + consumed_I(500_000) = 1_600_000
//
//	II.c(1_500_000) + ΔcII(100_000) = 1_600_000  ✓
func TestComputeReinvestment_HappyPath_50PctRate(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dept := newExtDeptI(sid)

	r, err := ComputeReinvestment(dept, 5000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.DeltaConstantPence != 400_000 {
		t.Fatalf("DeltaC: want 400_000, got %d", r.DeltaConstantPence)
	}
	if r.DeltaVariablePence != 100_000 {
		t.Fatalf("DeltaV: want 100_000, got %d", r.DeltaVariablePence)
	}
	if r.ConsumedSurplusPence != 500_000 {
		t.Fatalf("ConsumedSurplus: want 500_000, got %d", r.ConsumedSurplusPence)
	}
	// DeltaC + DeltaV + Consumed == SurplusPence
	total := r.DeltaConstantPence + r.DeltaVariablePence + r.ConsumedSurplusPence
	if total != dept.SurplusPence {
		t.Fatalf("partition invariant: DeltaC+DeltaV+Consumed=%d != SurplusPence=%d",
			total, dept.SurplusPence)
	}
}

// TestComputeReinvestment_ZeroRate verifies that a 0 bps rate produces
// no reinvestment — the simple-reproduction limiting case.
func TestComputeReinvestment_ZeroRate(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dept := newExtDeptI(sid)

	r, err := ComputeReinvestment(dept, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.DeltaConstantPence != 0 {
		t.Fatalf("DeltaConstant: want 0, got %d", r.DeltaConstantPence)
	}
	if r.DeltaVariablePence != 0 {
		t.Fatalf("DeltaVariable: want 0, got %d", r.DeltaVariablePence)
	}
	if r.ConsumedSurplusPence != dept.SurplusPence {
		t.Fatalf("ConsumedSurplus: want %d (all surplus), got %d",
			dept.SurplusPence, r.ConsumedSurplusPence)
	}
}

// TestComputeReinvestment_100PctRate verifies that a 10000 bps (100%) rate
// reinvests all surplus: ConsumedSurplus == 0.
func TestComputeReinvestment_100PctRate(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dept := newExtDeptI(sid)

	r, err := ComputeReinvestment(dept, 10000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.ConsumedSurplusPence != 0 {
		t.Fatalf("ConsumedSurplus: want 0, got %d", r.ConsumedSurplusPence)
	}
	gotTotal := r.DeltaConstantPence + r.DeltaVariablePence
	if gotTotal != dept.SurplusPence {
		t.Fatalf("DeltaC+DeltaV: want %d (all surplus), got %d",
			dept.SurplusPence, gotTotal)
	}
}

// TestComputeReinvestment_NegativeRate verifies that a negative rate returns
// ErrNegativeAccumulationRate.
func TestComputeReinvestment_NegativeRate(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dept := newExtDeptI(sid)

	_, err := ComputeReinvestment(dept, -1)
	if !errors.Is(err, ErrNegativeAccumulationRate) {
		t.Fatalf("expected ErrNegativeAccumulationRate, got %v", err)
	}
}

// TestComputeReinvestment_ExceedsMax verifies that a rate > 10000 bps returns
// ErrAccumulationRateExceeds100Pct.
func TestComputeReinvestment_ExceedsMax(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dept := newExtDeptI(sid)

	_, err := ComputeReinvestment(dept, 10001)
	if !errors.Is(err, ErrAccumulationRateExceeds100Pct) {
		t.Fatalf("expected ErrAccumulationRateExceeds100Pct, got %v", err)
	}
}

// --- ComputeNextPeriodCapital -----------------------------------------------

// TestComputeNextPeriodCapital_AdvancesCorrectly verifies that after applying a
// 50% accumulation reinvestment to Dept I, the new capital has grown correctly:
// new c = old c + ΔcI, new v = old v + ΔvI, TotalPence = c+v+s invariant holds.
func TestComputeNextPeriodCapital_AdvancesCorrectly(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dept := newExtDeptI(sid)

	r, err := ComputeReinvestment(dept, 5000)
	if err != nil {
		t.Fatalf("ComputeReinvestment: %v", err)
	}

	next := ComputeNextPeriodCapital(dept, r)
	if next.ConstantPence != dept.ConstantPence+r.DeltaConstantPence {
		t.Fatalf("ConstantPence: want %d, got %d",
			dept.ConstantPence+r.DeltaConstantPence, next.ConstantPence)
	}
	if next.VariablePence != dept.VariablePence+r.DeltaVariablePence {
		t.Fatalf("VariablePence: want %d, got %d",
			dept.VariablePence+r.DeltaVariablePence, next.VariablePence)
	}
	if next.TotalPence != next.ConstantPence+next.VariablePence+next.SurplusPence {
		t.Fatalf("TotalPence invariant violated: c+v+s=%d != total=%d",
			next.ConstantPence+next.VariablePence+next.SurplusPence, next.TotalPence)
	}
	// Capital must have grown
	if next.TotalPence <= dept.TotalPence {
		t.Fatalf("expected capital growth: next.Total=%d > prev.Total=%d",
			next.TotalPence, dept.TotalPence)
	}
}

// TestComputeNextPeriodCapital_ZeroReinvestmentIsIdentity verifies that a
// zero-delta reinvestment leaves constant and variable pence unchanged.
func TestComputeNextPeriodCapital_ZeroReinvestmentIsIdentity(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dept := newExtDeptI(sid)

	r := Reinvestment{
		SchemeID:             sid,
		Department:           DepartmentI,
		DeltaConstantPence:   0,
		DeltaVariablePence:   0,
		ConsumedSurplusPence: dept.SurplusPence,
	}

	next := ComputeNextPeriodCapital(dept, r)
	if next.ConstantPence != dept.ConstantPence {
		t.Fatalf("ConstantPence: want %d, got %d", dept.ConstantPence, next.ConstantPence)
	}
	if next.VariablePence != dept.VariablePence {
		t.Fatalf("VariablePence: want %d, got %d", dept.VariablePence, next.VariablePence)
	}
}

// --- ComputeCompositionBps --------------------------------------------------

// TestComputeCompositionBps_DeptI verifies the organic composition for
// Marx's Dept I: c/v = 4_000_000 / 1_000_000 = 4.0 → 40000 bps.
func TestComputeCompositionBps_DeptI(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dept := newExtDeptI(sid)

	want := int64(40000) // 4000/1000 = 4.0 → 40000 bps
	got := ComputeCompositionBps(dept)
	if got != want {
		t.Fatalf("CompositionBps: want %d, got %d", want, got)
	}
}

// TestComputeCompositionBps_DeptII verifies the organic composition for
// Marx's Dept II: c/v = 1_500_000 / 750_000 = 2.0 → 20000 bps.
func TestComputeCompositionBps_DeptII(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dept := newExtDeptII(sid)

	want := int64(20000) // 1500/750 = 2.0 → 20000 bps
	got := ComputeCompositionBps(dept)
	if got != want {
		t.Fatalf("CompositionBps Dept II: want %d, got %d", want, got)
	}
}

// TestComputeCompositionBps_ZeroVariable verifies that a zero variable
// capital returns 0 without panicking.
func TestComputeCompositionBps_ZeroVariable(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dept := DepartmentalCapital{
		ID:            NewDepartmentalCapitalID(),
		SchemeID:      SimpleReproductionSchemeID(sid),
		Department:    DepartmentI,
		ConstantPence: 4_000_000,
		VariablePence: 0,
		SurplusPence:  0,
		TotalPence:    4_000_000,
	}

	got := ComputeCompositionBps(dept)
	if got != 0 {
		t.Fatalf("CompositionBps with zero variable: want 0, got %d", got)
	}
}

// --- ComputeGrowthBps -------------------------------------------------------

// TestComputeGrowthBps_DeptIFasterThanDeptII verifies that with the standard
// Ch. 21 accumulation schedule Dept I (50% rate) grows proportionally faster
// than Dept II (20% rate).
func TestComputeGrowthBps_DeptIFasterThanDeptII(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()

	prevI := newExtDeptI(sid)
	rI, _ := ComputeReinvestment(prevI, 5000)
	nextI := ComputeNextPeriodCapital(prevI, rI)

	prevII := newExtDeptII(sid)
	rII, _ := ComputeReinvestment(prevII, 2000) // slower accumulation
	nextII := ComputeNextPeriodCapital(prevII, rII)

	growthI := ComputeGrowthBps(prevI, nextI)
	growthII := ComputeGrowthBps(prevII, nextII)

	if growthI <= growthII {
		t.Fatalf("expected Dept I growth (%d bps) > Dept II growth (%d bps)", growthI, growthII)
	}
}

// TestComputeGrowthBps_ZeroPrev verifies that a zero previous total returns 0.
func TestComputeGrowthBps_ZeroPrev(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	prev := DepartmentalCapital{
		ID:            NewDepartmentalCapitalID(),
		SchemeID:      SimpleReproductionSchemeID(sid),
		Department:    DepartmentI,
		ConstantPence: 0,
		VariablePence: 0,
		SurplusPence:  0,
		TotalPence:    0,
	}
	curr := newExtDeptI(sid)

	got := ComputeGrowthBps(prev, curr)
	if got != 0 {
		t.Fatalf("GrowthBps with zero prev: want 0, got %d", got)
	}
}

// TestComputeGrowthBps_NoGrowth verifies that identical periods produce 0 bps.
func TestComputeGrowthBps_NoGrowth(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dept := newExtDeptI(sid)

	got := ComputeGrowthBps(dept, dept)
	if got != 0 {
		t.Fatalf("GrowthBps for identical periods: want 0, got %d", got)
	}
}

// --- ComputeExtendedMoneyRequirement ----------------------------------------

// TestComputeExtendedMoneyRequirement_HappyPath verifies the money requirement
// calculation with Marx's fixture.
//
//	totalCapital = 6_000_000 + 3_000_000 = 9_000_000
//	additionalMoney = 9_000_000 * (5000 + 3000) / 20000 = 3_600_000
//	totalMoney = 1_000_000 + 3_600_000 = 4_600_000
func TestComputeExtendedMoneyRequirement_HappyPath(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dI := newExtDeptI(sid)
	dII := newExtDeptII(sid)

	scheme := ExtendedReproductionScheme{
		ID:                    sid,
		Period:                "1871",
		AccumulationRateIBps:  5000,
		AccumulationRateIIBps: 3000,
		DepartmentI:           &dI,
		DepartmentII:          &dII,
	}

	baseMoney := int64(1_000_000)
	req, err := ComputeExtendedMoneyRequirement(scheme, baseMoney)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.BaseMoneyPence != baseMoney {
		t.Fatalf("BaseMoneyPence: want %d, got %d", baseMoney, req.BaseMoneyPence)
	}
	if req.TotalMoneyPence != req.BaseMoneyPence+req.AdditionalMoneyPence {
		t.Fatalf("TotalMoney invariant: Base(%d) + Additional(%d) != Total(%d)",
			req.BaseMoneyPence, req.AdditionalMoneyPence, req.TotalMoneyPence)
	}
	if req.AdditionalMoneyPence < 0 {
		t.Fatalf("AdditionalMoneyPence must be >= 0, got %d", req.AdditionalMoneyPence)
	}
	// Computed: 9_000_000 * 8000 / 20000 = 3_600_000
	wantAdditional := int64(3_600_000)
	if req.AdditionalMoneyPence != wantAdditional {
		t.Fatalf("AdditionalMoneyPence: want %d, got %d", wantAdditional, req.AdditionalMoneyPence)
	}
	if req.SchemeID != sid {
		t.Fatalf("SchemeID: want %q, got %q", sid, req.SchemeID)
	}
}

// TestComputeExtendedMoneyRequirement_NilDepts verifies that a scheme with
// neither department returns ErrExtendedImbalance.
func TestComputeExtendedMoneyRequirement_NilDepts(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()

	scheme := ExtendedReproductionScheme{
		ID:                    sid,
		Period:                "1871",
		AccumulationRateIBps:  5000,
		AccumulationRateIIBps: 3000,
		DepartmentI:           nil,
		DepartmentII:          nil,
	}

	_, err := ComputeExtendedMoneyRequirement(scheme, 1_000_000)
	if !errors.Is(err, ErrExtendedImbalance) {
		t.Fatalf("expected ErrExtendedImbalance, got %v", err)
	}
}

// TestComputeExtendedMoneyRequirement_OnlyDeptI verifies that a scheme with
// only DepartmentI returns ErrExtendedImbalance.
func TestComputeExtendedMoneyRequirement_OnlyDeptI(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dI := newExtDeptI(sid)

	scheme := ExtendedReproductionScheme{
		ID:                    sid,
		Period:                "1871",
		AccumulationRateIBps:  5000,
		AccumulationRateIIBps: 3000,
		DepartmentI:           &dI,
		DepartmentII:          nil,
	}

	_, err := ComputeExtendedMoneyRequirement(scheme, 1_000_000)
	if !errors.Is(err, ErrExtendedImbalance) {
		t.Fatalf("expected ErrExtendedImbalance for missing Dept II, got %v", err)
	}
}

// TestComputeExtendedMoneyRequirement_ZeroRates verifies zero accumulation
// rates produce zero additional money.
func TestComputeExtendedMoneyRequirement_ZeroRates(t *testing.T) {
	t.Parallel()
	sid := NewExtendedReproductionSchemeID()
	dI := newExtDeptI(sid)
	dII := newExtDeptII(sid)

	scheme := ExtendedReproductionScheme{
		ID:                    sid,
		Period:                "1871",
		AccumulationRateIBps:  0,
		AccumulationRateIIBps: 0,
		DepartmentI:           &dI,
		DepartmentII:          &dII,
	}

	req, err := ComputeExtendedMoneyRequirement(scheme, 1_000_000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req.AdditionalMoneyPence != 0 {
		t.Fatalf("AdditionalMoney with zero rates: want 0, got %d", req.AdditionalMoneyPence)
	}
	if req.TotalMoneyPence != req.BaseMoneyPence {
		t.Fatalf("TotalMoney should equal BaseMoney when rates=0: %d != %d",
			req.TotalMoneyPence, req.BaseMoneyPence)
	}
}

// --- ID constructors --------------------------------------------------------

// TestNewExtendedReproductionSchemeID_Unique verifies that two generated IDs
// are distinct and non-empty.
func TestNewExtendedReproductionSchemeID_Unique(t *testing.T) {
	t.Parallel()
	a := NewExtendedReproductionSchemeID()
	b := NewExtendedReproductionSchemeID()
	if a == "" || b == "" {
		t.Fatal("expected non-empty IDs")
	}
	if a == b {
		t.Fatal("expected distinct IDs")
	}
}

// TestNewReinvestmentID_Unique verifies Reinvestment IDs are unique.
func TestNewReinvestmentID_Unique(t *testing.T) {
	t.Parallel()
	a := NewReinvestmentID()
	b := NewReinvestmentID()
	if a == "" || b == "" {
		t.Fatal("expected non-empty IDs")
	}
	if a == b {
		t.Fatal("expected distinct IDs")
	}
}

// TestNewMultiPeriodSchemeID_Unique verifies MultiPeriodScheme IDs are unique.
func TestNewMultiPeriodSchemeID_Unique(t *testing.T) {
	t.Parallel()
	a := NewMultiPeriodSchemeID()
	b := NewMultiPeriodSchemeID()
	if a == "" || b == "" {
		t.Fatal("expected non-empty IDs")
	}
	if a == b {
		t.Fatal("expected distinct IDs")
	}
}

// TestNewCompositionShiftID_Unique verifies CompositionShift IDs are unique.
func TestNewCompositionShiftID_Unique(t *testing.T) {
	t.Parallel()
	a := NewCompositionShiftID()
	b := NewCompositionShiftID()
	if a == "" || b == "" {
		t.Fatal("expected non-empty IDs")
	}
	if a == b {
		t.Fatal("expected distinct IDs")
	}
}

// TestNewDepartmentIGrowthLeadID_Unique verifies DepartmentIGrowthLead IDs are unique.
func TestNewDepartmentIGrowthLeadID_Unique(t *testing.T) {
	t.Parallel()
	a := NewDepartmentIGrowthLeadID()
	b := NewDepartmentIGrowthLeadID()
	if a == "" || b == "" {
		t.Fatal("expected non-empty IDs")
	}
	if a == b {
		t.Fatal("expected distinct IDs")
	}
}

// --- Sentinel errors --------------------------------------------------------

// TestSentinelErrors_Distinct verifies that the three Ch. 21 sentinel errors
// are distinct and non-nil.
func TestSentinelErrors_Distinct(t *testing.T) {
	t.Parallel()
	if ErrExtendedImbalance == nil {
		t.Fatal("ErrExtendedImbalance must be non-nil")
	}
	if ErrNegativeAccumulationRate == nil {
		t.Fatal("ErrNegativeAccumulationRate must be non-nil")
	}
	if ErrAccumulationRateExceeds100Pct == nil {
		t.Fatal("ErrAccumulationRateExceeds100Pct must be non-nil")
	}
	if errors.Is(ErrExtendedImbalance, ErrNegativeAccumulationRate) {
		t.Fatal("ErrExtendedImbalance and ErrNegativeAccumulationRate must be distinct")
	}
	if errors.Is(ErrNegativeAccumulationRate, ErrAccumulationRateExceeds100Pct) {
		t.Fatal("ErrNegativeAccumulationRate and ErrAccumulationRateExceeds100Pct must be distinct")
	}
}
