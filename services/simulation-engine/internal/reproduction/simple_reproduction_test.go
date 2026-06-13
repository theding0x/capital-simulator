package reproduction_test

import (
	"testing"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// Marx's canonical Vol. II Ch. 20 fixture (in units of 1 000 pence):
//   Dept I:  4 000c + 1 000v + 1 000s = 6 000  total
//   Dept II: 2 000c +   500v +   500s = 3 000  total
//   Keystone: I(v+s) = 2 000 == II(c) = 2 000  ✓

const (
	marxDeptIConstant  int64 = 4_000_000
	marxDeptIVariable  int64 = 1_000_000
	marxDeptISurplus   int64 = 1_000_000
	marxDeptITotal     int64 = 6_000_000
	marxDeptIIConstant int64 = 2_000_000
	marxDeptIIVariable int64 = 500_000
	marxDeptIISurplus  int64 = 500_000
	marxDeptIITotal    int64 = 3_000_000
)

func newMarxDeptI(schemeID repro.SimpleReproductionSchemeID) repro.DepartmentalCapital {
	return repro.DepartmentalCapital{
		ID:            repro.NewDepartmentalCapitalID(),
		SchemeID:      schemeID,
		Department:    repro.DepartmentI,
		ConstantPence: marxDeptIConstant,
		VariablePence: marxDeptIVariable,
		SurplusPence:  marxDeptISurplus,
		TotalPence:    marxDeptITotal,
	}
}

func newMarxDeptII(schemeID repro.SimpleReproductionSchemeID) repro.DepartmentalCapital {
	return repro.DepartmentalCapital{
		ID:            repro.NewDepartmentalCapitalID(),
		SchemeID:      schemeID,
		Department:    repro.DepartmentII,
		ConstantPence: marxDeptIIConstant,
		VariablePence: marxDeptIIVariable,
		SurplusPence:  marxDeptIISurplus,
		TotalPence:    marxDeptIITotal,
	}
}

// --- DepartmentalCapital.Validate -------------------------------------------

// TestDepartmentalCapitalValidateHappy verifies Marx's canonical fixtures pass.
func TestDepartmentalCapitalValidateHappy(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	if err := newMarxDeptI(sid).Validate(); err != nil {
		t.Fatalf("Dept I validate: unexpected error: %v", err)
	}
	if err := newMarxDeptII(sid).Validate(); err != nil {
		t.Fatalf("Dept II validate: unexpected error: %v", err)
	}
}

// TestDepartmentalCapitalValidateTotalMismatch verifies that a bad total is
// rejected with ErrDepartmentTotalMismatch.
func TestDepartmentalCapitalValidateTotalMismatch(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	d := newMarxDeptI(sid)
	d.TotalPence = d.ConstantPence + d.VariablePence // surplus omitted intentionally
	if err := d.Validate(); err == nil {
		t.Fatal("expected ErrDepartmentTotalMismatch, got nil")
	}
}

// TestDepartmentalCapitalVplusSPence verifies I(v+s) = 2 000 000 for the
// canonical Marx fixture.
func TestDepartmentalCapitalVplusSPence(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	d := newMarxDeptI(sid)
	want := int64(2_000_000) // 1 000 000v + 1 000 000s
	if got := d.VplusSPence(); got != want {
		t.Fatalf("VplusSPence: want %d, got %d", want, got)
	}
}

// TestDepartmentalCapitalValidateZero verifies a zero-value department passes
// validation (zero surplus is allowable in theory).
func TestDepartmentalCapitalValidateZero(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	d := repro.DepartmentalCapital{
		ID:            repro.NewDepartmentalCapitalID(),
		SchemeID:      sid,
		Department:    repro.DepartmentII,
		ConstantPence: 0,
		VariablePence: 0,
		SurplusPence:  0,
		TotalPence:    0,
	}
	if err := d.Validate(); err != nil {
		t.Fatalf("zero-value dept: unexpected error: %v", err)
	}
}

// --- ComputeIsBalanced ------------------------------------------------------

// TestComputeIsBalancedTrue verifies the keystone equation holds for Marx's
// canonical fixture: I(v+s) = 2 000 000 == II(c) = 2 000 000.
func TestComputeIsBalancedTrue(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	dI := newMarxDeptI(sid)
	dII := newMarxDeptII(sid)
	if !repro.ComputeIsBalanced(dI, dII) {
		t.Fatalf("expected balanced (I(v+s)=%d, II(c)=%d)", dI.VplusSPence(), dII.ConstantPence)
	}
}

// TestComputeIsBalancedFalse verifies the equation fails when II(c) != I(v+s).
func TestComputeIsBalancedFalse(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	dI := newMarxDeptI(sid)
	dII := newMarxDeptII(sid)
	dII.ConstantPence = 1_500_000 // deliberately wrong
	if repro.ComputeIsBalanced(dI, dII) {
		t.Fatal("expected imbalanced, got balanced")
	}
}

// --- ComputeDeficit ---------------------------------------------------------

// TestComputeDeficitZero verifies that a balanced scheme has zero deficit.
func TestComputeDeficitZero(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	dI := newMarxDeptI(sid)
	dII := newMarxDeptII(sid)
	if got := repro.ComputeDeficit(dI, dII); got != 0 {
		t.Fatalf("deficit: want 0, got %d", got)
	}
}

// TestComputeDeficitPositive verifies a positive deficit when II(c) < I(v+s).
func TestComputeDeficitPositive(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	dI := newMarxDeptI(sid)
	dII := newMarxDeptII(sid)
	dII.ConstantPence = 1_500_000 // deficit of 500 000
	want := int64(500_000)
	if got := repro.ComputeDeficit(dI, dII); got != want {
		t.Fatalf("deficit: want %d, got %d", want, got)
	}
}

// TestComputeDeficitNegative verifies a negative deficit (oversupply) when
// II(c) > I(v+s).
func TestComputeDeficitNegative(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	dI := newMarxDeptI(sid)
	dII := newMarxDeptII(sid)
	dII.ConstantPence = 2_500_000 // oversupply
	want := int64(-500_000)
	if got := repro.ComputeDeficit(dI, dII); got != want {
		t.Fatalf("deficit: want %d, got %d", want, got)
	}
}

// --- SimpleReproductionScheme.CheckBalance ----------------------------------

// TestSchemeCheckBalanceWithBothDepartments verifies CheckBalance returns true
// when both departments are present and balanced.
func TestSchemeCheckBalanceWithBothDepartments(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	dI := newMarxDeptI(sid)
	dII := newMarxDeptII(sid)
	scheme := repro.SimpleReproductionScheme{
		ID:           sid,
		Period:       "1865",
		DepartmentI:  &dI,
		DepartmentII: &dII,
	}
	if !scheme.CheckBalance() {
		t.Fatal("expected balanced scheme")
	}
}

// TestSchemeCheckBalanceMissingDept verifies CheckBalance returns false when
// only one department is present.
func TestSchemeCheckBalanceMissingDept(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	dI := newMarxDeptI(sid)
	scheme := repro.SimpleReproductionScheme{
		ID:          sid,
		Period:      "1865",
		DepartmentI: &dI,
	}
	if scheme.CheckBalance() {
		t.Fatal("expected false when Dept II is missing")
	}
}

// TestSchemeCheckBalanceImbalanced verifies CheckBalance returns false for an
// imbalanced two-department scheme.
func TestSchemeCheckBalanceImbalanced(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	dI := newMarxDeptI(sid)
	dII := newMarxDeptII(sid)
	dII.ConstantPence = 2_500_000 // II(c) > I(v+s)
	scheme := repro.SimpleReproductionScheme{
		ID:           sid,
		Period:       "1865",
		DepartmentI:  &dI,
		DepartmentII: &dII,
	}
	if scheme.CheckBalance() {
		t.Fatal("expected imbalanced scheme")
	}
}

// --- ExchangeKind.Validate --------------------------------------------------

// TestExchangeKindValidateHappy verifies known kinds pass.
func TestExchangeKindValidateHappy(t *testing.T) {
	t.Parallel()
	for _, k := range []repro.ExchangeKind{
		repro.ExchangeMoneyForGoods,
		repro.ExchangeGoodsForMoney,
	} {
		if err := k.Validate(); err != nil {
			t.Fatalf("kind %q: unexpected error: %v", k, err)
		}
	}
}

// TestExchangeKindValidateUnknown verifies an unrecognised kind fails.
func TestExchangeKindValidateUnknown(t *testing.T) {
	t.Parallel()
	k := repro.ExchangeKind("barter")
	if err := k.Validate(); err == nil {
		t.Fatal("expected ErrInvalidExchangeKind for unknown kind")
	}
}

// --- InterDepartmentExchange.Validate ---------------------------------------

// TestInterDepartmentExchangeValidateHappy verifies a canonical I→II exchange.
func TestInterDepartmentExchangeValidateHappy(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	ex := repro.InterDepartmentExchange{
		ID:             repro.NewInterDepartmentExchangeID(),
		SchemeID:       sid,
		FromDepartment: repro.DepartmentI,
		ToDepartment:   repro.DepartmentII,
		Pence:          2_000_000,
		Kind:           repro.ExchangeGoodsForMoney,
		Description:    "I(v+s) → II",
	}
	if err := ex.Validate(); err != nil {
		t.Fatalf("exchange validate: unexpected error: %v", err)
	}
}

// TestInterDepartmentExchangeValidateInvalidKind verifies a bad kind is rejected.
func TestInterDepartmentExchangeValidateInvalidKind(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	ex := repro.InterDepartmentExchange{
		ID:             repro.NewInterDepartmentExchangeID(),
		SchemeID:       sid,
		FromDepartment: repro.DepartmentI,
		ToDepartment:   repro.DepartmentII,
		Pence:          1_000_000,
		Kind:           repro.ExchangeKind("swap"),
	}
	if err := ex.Validate(); err == nil {
		t.Fatal("expected error for invalid exchange kind")
	}
}

// --- ID generators ----------------------------------------------------------

// TestNewSimpleReproductionSchemeID verifies unique non-empty IDs.
func TestNewSimpleReproductionSchemeID(t *testing.T) {
	t.Parallel()
	a := repro.NewSimpleReproductionSchemeID()
	b := repro.NewSimpleReproductionSchemeID()
	if a == "" || b == "" {
		t.Fatal("expected non-empty IDs")
	}
	if a == b {
		t.Fatal("expected distinct IDs")
	}
}

// TestNewDepartmentalCapitalID verifies unique non-empty IDs.
func TestNewDepartmentalCapitalID(t *testing.T) {
	t.Parallel()
	a := repro.NewDepartmentalCapitalID()
	b := repro.NewDepartmentalCapitalID()
	if a == "" || b == "" {
		t.Fatal("expected non-empty IDs")
	}
	if a == b {
		t.Fatal("expected distinct IDs")
	}
}

// --- IntraDepartmentIReplacement (FIX 2) ------------------------------------

// TestIntraDepartmentIReplacementInvariantHolds verifies that an
// IntraDepartmentIReplacement constructed from Dept I carries Pence equal
// to Dept I's ConstantPence (4 000 000 for the Marx fixture).
func TestIntraDepartmentIReplacementInvariantHolds(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	dI := newMarxDeptI(sid)
	r := repro.NewIntraDepartmentIReplacement(dI)
	if r.Pence != marxDeptIConstant {
		t.Fatalf("Pence: want %d (I.c), got %d", marxDeptIConstant, r.Pence)
	}
	if err := r.Validate(dI); err != nil {
		t.Fatalf("Validate: unexpected error for matching pence: %v", err)
	}
}

// TestIntraDepartmentIReplacementMismatchRejected verifies that a mismatched
// Pence value is rejected with ErrIntraDeptIMismatch.
func TestIntraDepartmentIReplacementMismatchRejected(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	dI := newMarxDeptI(sid)
	r := repro.IntraDepartmentIReplacement{Pence: marxDeptIConstant + 1}
	if err := r.Validate(dI); err == nil {
		t.Fatal("expected ErrIntraDeptIMismatch for mismatched pence, got nil")
	}
}

// --- ComputeMoneyClosedLoop (FIX 3) -----------------------------------------

// TestComputeMoneyClosedLoopClosed verifies that three Marx inter-department
// exchanges sum to zero net flow, yielding IsClosed=true.
//
// Marx's three flows (Ch. 20 §III):
//  1. Intra-Dept-I:   I→I  4 000 000 (intra; conventionally treated as zero net cross-flow)
//  2. I(v+s) → II:   I→II  2 000 000
//  3. II(c) → I:     II→I  2 000 000
//
// The net inter-department balance is: +2 000 000 (I→II) - 2 000 000 (II→I) = 0.
func TestComputeMoneyClosedLoopClosed(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	exchanges := []repro.InterDepartmentExchange{
		{
			ID:             repro.NewInterDepartmentExchangeID(),
			SchemeID:       sid,
			FromDepartment: repro.DepartmentI,
			ToDepartment:   repro.DepartmentII,
			Pence:          2_000_000,
			Kind:           repro.ExchangeGoodsForMoney,
			Description:    "I(v+s) → II",
		},
		{
			ID:             repro.NewInterDepartmentExchangeID(),
			SchemeID:       sid,
			FromDepartment: repro.DepartmentII,
			ToDepartment:   repro.DepartmentI,
			Pence:          2_000_000,
			Kind:           repro.ExchangeMoneyForGoods,
			Description:    "II(c) → I",
		},
	}
	loop := repro.ComputeMoneyClosedLoop(exchanges)
	if loop.NetFlowPence != 0 {
		t.Fatalf("NetFlowPence: want 0, got %d", loop.NetFlowPence)
	}
	if !loop.IsClosed {
		t.Fatal("IsClosed: want true, got false")
	}
}

// TestComputeMoneyClosedLoopNotClosed verifies that perturbing one exchange
// yields a non-zero net flow and IsClosed=false.
func TestComputeMoneyClosedLoopNotClosed(t *testing.T) {
	t.Parallel()
	sid := repro.NewSimpleReproductionSchemeID()
	exchanges := []repro.InterDepartmentExchange{
		{
			ID:             repro.NewInterDepartmentExchangeID(),
			SchemeID:       sid,
			FromDepartment: repro.DepartmentI,
			ToDepartment:   repro.DepartmentII,
			Pence:          2_500_000, // perturbed: surplus leaks
			Kind:           repro.ExchangeGoodsForMoney,
			Description:    "I(v+s) → II (perturbed)",
		},
		{
			ID:             repro.NewInterDepartmentExchangeID(),
			SchemeID:       sid,
			FromDepartment: repro.DepartmentII,
			ToDepartment:   repro.DepartmentI,
			Pence:          2_000_000,
			Kind:           repro.ExchangeMoneyForGoods,
			Description:    "II(c) → I",
		},
	}
	loop := repro.ComputeMoneyClosedLoop(exchanges)
	if loop.IsClosed {
		t.Fatalf("expected IsClosed=false for perturbed exchanges, NetFlowPence=%d", loop.NetFlowPence)
	}
	if loop.NetFlowPence == 0 {
		t.Fatal("expected non-zero NetFlowPence for perturbed exchanges")
	}
}
