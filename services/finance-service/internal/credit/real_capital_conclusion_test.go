package credit

import (
	"testing"
)

// ── Fixtures from Capital Vol. III Ch. 32 ────────────────────────────────────
//
// The conclusion of the trilogy: loanable money-capital always overstates real
// accumulation. Capital released by falling input prices, merchants' idle money,
// the rentier class, and revenue passing through money form all swell loan
// capital with no new productive capital behind them. In crisis the demand for
// loan capital is a demand for means of payment (convertibility), not productive
// capital. All fixtures here are ch32-prefixed because the whole chapter
// test-suite shares one package.

// ch32FixtureRelease returns the raw-material-price-fall release: capital freed by
// a fall in input prices (Cause=CauseInputPriceFall), not reinvested — it merely
// overstates real accumulation. Overstatement 2500 bp (25%).
func ch32FixtureRelease(t *testing.T) CapitalRelease {
	t.Helper()
	over := LoanCapitalOverstatement{OverstatementPct: 2500, Description: "freed money lent out, no new productive capital"}
	c, err := NewCapitalRelease(
		"Capital released by raw-material price fall",
		5_000_000,
		CauseInputPriceFall,
		false,
		over,
		"raw-material prices fall; the same reproduction ties up less money; the difference is lent out",
	)
	if err != nil {
		t.Fatalf("ch32FixtureRelease: %v", err)
	}
	return c
}

// ch32FixtureReleaseReinvested returns the same release once the freed capital is
// actually reinvested — it then reflects real accumulation.
func ch32FixtureReleaseReinvested(t *testing.T) CapitalRelease {
	t.Helper()
	over := LoanCapitalOverstatement{OverstatementPct: 0, Description: "freed money reinvested in production"}
	c, err := NewCapitalRelease(
		"Capital released by raw-material price fall (reinvested)",
		5_000_000,
		CauseInputPriceFall,
		true,
		over,
		"the freed capital is reinvested in production and reflects real accumulation",
	)
	if err != nil {
		t.Fatalf("ch32FixtureReleaseReinvested: %v", err)
	}
	return c
}

// ── CapitalRelease constructor ────────────────────────────────────────────────

func TestNewCapitalRelease_PriceFall(t *testing.T) {
	t.Parallel()
	c := ch32FixtureRelease(t)
	if c.ReleasedAmount != 5_000_000 {
		t.Errorf("ReleasedAmount: got %d, want 5000000", c.ReleasedAmount)
	}
	if c.Cause != CauseInputPriceFall {
		t.Errorf("Cause: got %d, want %d", c.Cause, CauseInputPriceFall)
	}
	if c.Reinvested {
		t.Error("Reinvested must be false for the unreinvested release")
	}
	if c.Overstatement.OverstatementPct != 2500 {
		t.Errorf("Overstatement.OverstatementPct: got %d, want 2500", c.Overstatement.OverstatementPct)
	}
	// ID and CreatedAt must be zero — the store assigns them.
	if !c.ID.IsZero() {
		t.Error("ID must be zero before store persists")
	}
	if !c.CreatedAt.IsZero() {
		t.Error("CreatedAt must be zero before store persists")
	}
}

func TestNewCapitalRelease_ZeroAmountAllowed(t *testing.T) {
	t.Parallel()
	c, err := NewCapitalRelease("zero", 0, CauseMerchantIdleMoney, false, LoanCapitalOverstatement{}, "")
	if err != nil {
		t.Fatalf("zero released amount must be allowed: %v", err)
	}
	if c.ReleasedAmount != 0 {
		t.Errorf("ReleasedAmount: got %d, want 0", c.ReleasedAmount)
	}
}

func TestNewCapitalRelease_NegativeAmount(t *testing.T) {
	t.Parallel()
	if _, err := NewCapitalRelease("bad", -1, CauseInputPriceFall, false, LoanCapitalOverstatement{}, ""); err != ErrNegativeReleasedAmount {
		t.Errorf("got %v, want ErrNegativeReleasedAmount", err)
	}
}

func TestNewCapitalRelease_InvalidCause(t *testing.T) {
	t.Parallel()
	if _, err := NewCapitalRelease("bad", 100, 0, false, LoanCapitalOverstatement{}, ""); err != ErrInvalidReleaseCause {
		t.Errorf("cause=0: got %v, want ErrInvalidReleaseCause", err)
	}
	if _, err := NewCapitalRelease("bad", 100, 5, false, LoanCapitalOverstatement{}, ""); err != ErrInvalidReleaseCause {
		t.Errorf("cause=5: got %v, want ErrInvalidReleaseCause", err)
	}
}

func TestNewCapitalRelease_NegativeOverstatement(t *testing.T) {
	t.Parallel()
	over := LoanCapitalOverstatement{OverstatementPct: -1}
	if _, err := NewCapitalRelease("bad", 100, CauseInputPriceFall, false, over, ""); err != ErrNegativeOverstatement {
		t.Errorf("got %v, want ErrNegativeOverstatement", err)
	}
}

// ── CapitalReleaseCause.IsValid ───────────────────────────────────────────────

func TestCapitalReleaseCause_IsValid(t *testing.T) {
	t.Parallel()
	for _, c := range []CapitalReleaseCause{CauseInputPriceFall, CauseMerchantIdleMoney, CauseRentierSavings, CauseRevenuePassthrough} {
		if !c.IsValid() {
			t.Errorf("cause %d must be valid", c)
		}
	}
	if CapitalReleaseCause(0).IsValid() {
		t.Error("cause 0 must be invalid")
	}
	if CapitalReleaseCause(5).IsValid() {
		t.Error("cause 5 must be invalid")
	}
}

// ── AccumulationVerdict ───────────────────────────────────────────────────────

func TestAccumulationVerdict_Overstates(t *testing.T) {
	t.Parallel()
	c := ch32FixtureRelease(t)
	got := c.AccumulationVerdict()
	want := "overstates real accumulation — the released capital swells loan capital but is not reinvested"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAccumulationVerdict_Reflects(t *testing.T) {
	t.Parallel()
	c := ch32FixtureReleaseReinvested(t)
	got := c.AccumulationVerdict()
	want := "reflects real accumulation — the released capital is reinvested in production"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// ── RevenueLoanCapitalConversion ──────────────────────────────────────────────

func TestNewRevenueLoanCapitalConversion_Positive(t *testing.T) {
	t.Parallel()
	rc, err := NewRevenueLoanCapitalConversion(120_000, "spinner's yarn proceeds set aside as revenue, deposited, become loan capital")
	if err != nil {
		t.Fatalf("positive revenue must be allowed: %v", err)
	}
	if rc.RevenueAmount != 120_000 {
		t.Errorf("RevenueAmount: got %d, want 120000", rc.RevenueAmount)
	}
	if !rc.ID.IsZero() {
		t.Error("ID must be zero until assigned")
	}
}

func TestNewRevenueLoanCapitalConversion_NonPositive(t *testing.T) {
	t.Parallel()
	if _, err := NewRevenueLoanCapitalConversion(0, ""); err != ErrNonPositiveRevenue {
		t.Errorf("revenue=0: got %v, want ErrNonPositiveRevenue", err)
	}
	if _, err := NewRevenueLoanCapitalConversion(-5, ""); err != ErrNonPositiveRevenue {
		t.Errorf("revenue=-5: got %v, want ErrNonPositiveRevenue", err)
	}
}

// ── LoanCapitalOverstatement ──────────────────────────────────────────────────

func TestLoanCapitalOverstatement_NonNegative(t *testing.T) {
	t.Parallel()
	over := LoanCapitalOverstatement{OverstatementPct: 0, Description: "boundary zero"}
	if over.OverstatementPct != 0 {
		t.Errorf("OverstatementPct: got %d, want 0", over.OverstatementPct)
	}
	c, err := NewCapitalRelease("ok", 100, CauseRentierSavings, false, over, "")
	if err != nil {
		t.Fatalf("zero overstatement must be allowed: %v", err)
	}
	if c.Overstatement.OverstatementPct != 0 {
		t.Errorf("round-trip OverstatementPct: got %d, want 0", c.Overstatement.OverstatementPct)
	}
}

// ── MeansOfPaymentDemand ──────────────────────────────────────────────────────

func TestNewMeansOfPaymentDemand_Crisis1857(t *testing.T) {
	t.Parallel()
	// Even if a caller asks for productive-capital demand, a crisis forces it to
	// means of payment only.
	d, err := NewMeansOfPaymentDemand(PhaseCrisis, 4_000_000, true, "1857 crisis: cry for convertibility into gold")
	if err != nil {
		t.Fatalf("crisis demand must construct: %v", err)
	}
	if d.IsProductiveCapitalDemand {
		t.Error("in PhaseCrisis IsProductiveCapitalDemand must be false")
	}
	if d.Phase != PhaseCrisis {
		t.Errorf("Phase: got %d, want %d", d.Phase, PhaseCrisis)
	}
}

func TestNewMeansOfPaymentDemand_NonCrisisKeepsFlag(t *testing.T) {
	t.Parallel()
	d, err := NewMeansOfPaymentDemand(PhaseProsperity, 1_000_000, true, "boom: genuine demand for productive capital")
	if err != nil {
		t.Fatalf("prosperity demand must construct: %v", err)
	}
	if !d.IsProductiveCapitalDemand {
		t.Error("outside a crisis the productive-capital flag must be preserved")
	}
}

func TestNewMeansOfPaymentDemand_InvalidPhase(t *testing.T) {
	t.Parallel()
	if _, err := NewMeansOfPaymentDemand(0, 100, false, ""); err != ErrInvalidCyclePhase {
		t.Errorf("phase=0: got %v, want ErrInvalidCyclePhase", err)
	}
	if _, err := NewMeansOfPaymentDemand(6, 100, false, ""); err != ErrInvalidCyclePhase {
		t.Errorf("phase=6: got %v, want ErrInvalidCyclePhase", err)
	}
}

// ── Rentier value type ────────────────────────────────────────────────────────

func TestRentier_GrowingClass(t *testing.T) {
	t.Parallel()
	r := Rentier{
		Name:           "Coupon-clipping annuitant",
		InterestIncome: 250_000,
		Description:    "lives off interest, withdraws savings from productive investment",
	}
	if r.InterestIncome != 250_000 {
		t.Errorf("InterestIncome: got %d, want 250000", r.InterestIncome)
	}
}

// ── ID format / uniqueness ────────────────────────────────────────────────────

func TestNewCapitalReleaseID_Format(t *testing.T) {
	t.Parallel()
	id := NewCapitalReleaseID()
	if len(string(id)) != 24 {
		t.Errorf("ID length: got %d, want 24 hex chars", len(string(id)))
	}
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}

func TestNewCapitalReleaseID_Unique(t *testing.T) {
	t.Parallel()
	a := NewCapitalReleaseID()
	b := NewCapitalReleaseID()
	if a == b {
		t.Error("two IDs must not be equal")
	}
}

func TestCapitalReleaseID_IsZero(t *testing.T) {
	t.Parallel()
	var id CapitalReleaseID
	if !id.IsZero() {
		t.Error("zero value must be zero")
	}
	id = NewCapitalReleaseID()
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}

func TestNewRevenueLoanCapitalConversionID_Format(t *testing.T) {
	t.Parallel()
	id := NewRevenueLoanCapitalConversionID()
	if len(string(id)) != 24 {
		t.Errorf("ID length: got %d, want 24 hex chars", len(string(id)))
	}
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}

func TestNewMeansOfPaymentDemandID_Format(t *testing.T) {
	t.Parallel()
	id := NewMeansOfPaymentDemandID()
	if len(string(id)) != 24 {
		t.Errorf("ID length: got %d, want 24 hex chars", len(string(id)))
	}
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}
