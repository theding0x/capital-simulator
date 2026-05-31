package credit

import (
	"testing"
)

// ── Fixtures from Capital Vol. III Ch. 31 ────────────────────────────────────
//
// Marx distinguishes two sources of loanable money-capital: the simple
// transformation of idle money into loan capital (inverse to real accumulation),
// and the conversion of capital/revenue into money that becomes loan capital
// (accompanies real accumulation). Bill-brokers redistribute surplus rural
// deposits to manufacturing. His arithmetical illustration: the same £20 lent
// five times performs £100 of loan-capital service. All fixtures here are
// ch31-prefixed because the whole chapter test-suite shares one package.

// ch31FixtureVelocity20x5 returns the £20-lent-5-times velocity: money=2000 pence,
// timesLent=5 → loanCapitalCreated=10000 pence (£100).
func ch31FixtureVelocity20x5(t *testing.T) LoanCapitalVelocity {
	t.Helper()
	v, err := NewLoanCapitalVelocity(2000, 5)
	if err != nil {
		t.Fatalf("ch31FixtureVelocity20x5: %v", err)
	}
	return v
}

// ch31FixtureIpswich returns the Ipswich-farmers parcel: surplus rural deposits
// redistributed to manufacturing — Source=SourceCapitalRevenueConversion, which
// accompanies real accumulation.
func ch31FixtureIpswich(t *testing.T) FloatingCapital {
	t.Helper()
	saving := CirculationReserveSaving{Amount: 50_000, Description: "country-bank till economies"}
	f, err := NewFloatingCapital(
		"Ipswich farmers' surplus deposits",
		2_000_000,
		SourceCapitalRevenueConversion,
		ch31FixtureVelocity20x5(t),
		saving,
		"rural surplus rediscounted to the manufacturing centres",
	)
	if err != nil {
		t.Fatalf("ch31FixtureIpswich: %v", err)
	}
	return f
}

// ch31FixturePostCrisisIdle returns the post-crisis idle parcel:
// Source=SourceSimpleTransformation, inverse to real accumulation.
func ch31FixturePostCrisisIdle(t *testing.T) FloatingCapital {
	t.Helper()
	saving := CirculationReserveSaving{Amount: 0, Description: "no economies recorded"}
	v, err := NewLoanCapitalVelocity(100_000, 1)
	if err != nil {
		t.Fatalf("ch31FixturePostCrisisIdle velocity: %v", err)
	}
	f, err := NewFloatingCapital(
		"Post-crisis idle industrial capital",
		3_000_000,
		SourceSimpleTransformation,
		v,
		saving,
		"money lying fallow after a crisis, transformed into loan capital",
	)
	if err != nil {
		t.Fatalf("ch31FixturePostCrisisIdle: %v", err)
	}
	return f
}

// ── LoanCapitalVelocity ───────────────────────────────────────────────────────

func TestNewLoanCapitalVelocity_20x5(t *testing.T) {
	t.Parallel()
	v := ch31FixtureVelocity20x5(t)
	if v.MoneyAmount != 2000 {
		t.Errorf("MoneyAmount: got %d, want 2000", v.MoneyAmount)
	}
	if v.TimesLent != 5 {
		t.Errorf("TimesLent: got %d, want 5", v.TimesLent)
	}
	// £20 × 5 = £100 of loan capital — the invariant must hold exactly.
	if v.LoanCapitalCreated != 10000 {
		t.Errorf("LoanCapitalCreated: got %d, want 10000", v.LoanCapitalCreated)
	}
	if v.LoanCapitalCreated != v.MoneyAmount*v.TimesLent {
		t.Errorf("invariant LoanCapitalCreated == MoneyAmount × TimesLent violated: %d != %d", v.LoanCapitalCreated, v.MoneyAmount*v.TimesLent)
	}
}

func TestNewLoanCapitalVelocity_BoundaryOne(t *testing.T) {
	t.Parallel()
	v, err := NewLoanCapitalVelocity(2000, 1)
	if err != nil {
		t.Fatalf("times_lent=1 must be allowed: %v", err)
	}
	if v.LoanCapitalCreated != 2000 {
		t.Errorf("LoanCapitalCreated: got %d, want 2000", v.LoanCapitalCreated)
	}
}

func TestNewLoanCapitalVelocity_TimesLentZero(t *testing.T) {
	t.Parallel()
	if _, err := NewLoanCapitalVelocity(2000, 0); err != ErrTimesLentTooLow {
		t.Errorf("times_lent=0: got %v, want ErrTimesLentTooLow", err)
	}
	if _, err := NewLoanCapitalVelocity(2000, -3); err != ErrTimesLentTooLow {
		t.Errorf("times_lent=-3: got %v, want ErrTimesLentTooLow", err)
	}
}

// ── FloatingCapital constructor ───────────────────────────────────────────────

func TestNewFloatingCapital_Ipswich(t *testing.T) {
	t.Parallel()
	f := ch31FixtureIpswich(t)
	if f.Amount != 2_000_000 {
		t.Errorf("Amount: got %d, want 2000000", f.Amount)
	}
	if f.Source != SourceCapitalRevenueConversion {
		t.Errorf("Source: got %d, want %d", f.Source, SourceCapitalRevenueConversion)
	}
	if f.Velocity.LoanCapitalCreated != 10000 {
		t.Errorf("Velocity.LoanCapitalCreated: got %d, want 10000", f.Velocity.LoanCapitalCreated)
	}
	if f.ReserveSaving.Amount != 50_000 {
		t.Errorf("ReserveSaving.Amount: got %d, want 50000", f.ReserveSaving.Amount)
	}
	// ID and CreatedAt must be zero — the store assigns them.
	if !f.ID.IsZero() {
		t.Error("ID must be zero before store persists")
	}
	if !f.CreatedAt.IsZero() {
		t.Error("CreatedAt must be zero before store persists")
	}
}

func TestNewFloatingCapital_ZeroAmountAllowed(t *testing.T) {
	t.Parallel()
	v := ch31FixtureVelocity20x5(t)
	f, err := NewFloatingCapital("zero", 0, SourceSimpleTransformation, v, CirculationReserveSaving{}, "")
	if err != nil {
		t.Fatalf("zero amount must be allowed: %v", err)
	}
	if f.Amount != 0 {
		t.Errorf("Amount: got %d, want 0", f.Amount)
	}
}

func TestNewFloatingCapital_NegativeAmount(t *testing.T) {
	t.Parallel()
	v := ch31FixtureVelocity20x5(t)
	if _, err := NewFloatingCapital("bad", -1, SourceSimpleTransformation, v, CirculationReserveSaving{}, ""); err != ErrNegativeFloatingAmount {
		t.Errorf("got %v, want ErrNegativeFloatingAmount", err)
	}
}

func TestNewFloatingCapital_InvalidSource(t *testing.T) {
	t.Parallel()
	v := ch31FixtureVelocity20x5(t)
	if _, err := NewFloatingCapital("bad", 100, 0, v, CirculationReserveSaving{}, ""); err != ErrInvalidFloatingSource {
		t.Errorf("source=0: got %v, want ErrInvalidFloatingSource", err)
	}
	if _, err := NewFloatingCapital("bad", 100, 3, v, CirculationReserveSaving{}, ""); err != ErrInvalidFloatingSource {
		t.Errorf("source=3: got %v, want ErrInvalidFloatingSource", err)
	}
}

func TestNewFloatingCapital_NegativeReserveSaving(t *testing.T) {
	t.Parallel()
	v := ch31FixtureVelocity20x5(t)
	saving := CirculationReserveSaving{Amount: -1}
	if _, err := NewFloatingCapital("bad", 100, SourceSimpleTransformation, v, saving, ""); err != ErrNegativeReserveSaving {
		t.Errorf("got %v, want ErrNegativeReserveSaving", err)
	}
}

// ── FloatingCapitalSource.IsValid ─────────────────────────────────────────────

func TestFloatingCapitalSource_IsValid(t *testing.T) {
	t.Parallel()
	if !SourceSimpleTransformation.IsValid() {
		t.Error("source 1 must be valid")
	}
	if !SourceCapitalRevenueConversion.IsValid() {
		t.Error("source 2 must be valid")
	}
	if FloatingCapitalSource(0).IsValid() {
		t.Error("source 0 must be invalid")
	}
	if FloatingCapitalSource(3).IsValid() {
		t.Error("source 3 must be invalid")
	}
}

// ── RelationToAccumulation ────────────────────────────────────────────────────

func TestRelationToAccumulation_Inverse(t *testing.T) {
	t.Parallel()
	f := ch31FixturePostCrisisIdle(t)
	got := f.RelationToAccumulation()
	want := "inverse to real accumulation — idle money transformed into loan capital swells when production stagnates"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRelationToAccumulation_Accompanies(t *testing.T) {
	t.Parallel()
	f := ch31FixtureIpswich(t)
	got := f.RelationToAccumulation()
	want := "accompanies real accumulation — capital and revenue converted into money on the way to fresh investment"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// ── FloatingCapitalID ─────────────────────────────────────────────────────────

func TestNewFloatingCapitalID_Format(t *testing.T) {
	t.Parallel()
	id := NewFloatingCapitalID()
	if len(string(id)) != 24 {
		t.Errorf("ID length: got %d, want 24 hex chars", len(string(id)))
	}
}

func TestNewFloatingCapitalID_Unique(t *testing.T) {
	t.Parallel()
	a := NewFloatingCapitalID()
	b := NewFloatingCapitalID()
	if a == b {
		t.Error("two IDs must not be equal")
	}
}

func TestFloatingCapitalID_IsZero(t *testing.T) {
	t.Parallel()
	var id FloatingCapitalID
	if !id.IsZero() {
		t.Error("zero value must be zero")
	}
	id = NewFloatingCapitalID()
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}

// ── LoanCapitalVelocityID / BillBrokerID ──────────────────────────────────────

func TestNewLoanCapitalVelocityID_Format(t *testing.T) {
	t.Parallel()
	id := NewLoanCapitalVelocityID()
	if len(string(id)) != 24 {
		t.Errorf("ID length: got %d, want 24 hex chars", len(string(id)))
	}
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}

func TestNewBillBrokerID_Format(t *testing.T) {
	t.Parallel()
	id := NewBillBrokerID()
	if len(string(id)) != 24 {
		t.Errorf("ID length: got %d, want 24 hex chars", len(string(id)))
	}
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}

// ── BillBroker / Rediscount value types ───────────────────────────────────────

func TestBillBroker_LombardStreet(t *testing.T) {
	t.Parallel()
	b := BillBroker{
		ID:               NewBillBrokerID(),
		Name:             "Overend, Gurney & Co.",
		RediscountVolume: 5_000_000,
		Description:      "Lombard Street bill-broker redistributing rural deposits to manufacturing",
	}
	if b.RediscountVolume != 5_000_000 {
		t.Errorf("RediscountVolume: got %d, want 5000000", b.RediscountVolume)
	}
	if b.ID.IsZero() {
		t.Error("broker ID must not be zero")
	}
}

func TestRediscount_RuralToCentre(t *testing.T) {
	t.Parallel()
	r := Rediscount{
		OriginalDiscounter: "Ipswich country bank",
		Broker:             "Overend, Gurney & Co.",
		BillAmount:         2000,
	}
	if r.BillAmount != 2000 {
		t.Errorf("BillAmount: got %d, want 2000", r.BillAmount)
	}
}
