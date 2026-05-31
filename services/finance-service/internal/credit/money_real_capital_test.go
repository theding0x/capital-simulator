package credit

import (
	"testing"
)

// ── Fixtures from Capital Vol. III Ch. 30 ────────────────────────────────────
//
// Marx contrasts the accumulation of money-capital with the accumulation of real
// capital. The 1847 crisis is the paradigm: commercial credit collapses not from
// a shortage of money but because returns are blocked (commodities will not sell)
// and mutual confidence evaporates. The revival case is the mirror image, where
// the two accumulations advance together. The cotton→yarn→cloth chain of bills is
// his illustration of commercial credit as the foundation of the credit system.

// fixtureCrisis1847 returns the 1847 crisis observation: commercial credit
// contracts, the cause is blocked returns + lost confidence (not money scarcity),
// and money-capital grows while real capital shrinks (no correspondence).
func fixtureCrisis1847(t *testing.T) RealCapitalAccumulation {
	t.Helper()
	corr := AccumulationCorrespondence{
		MoneyCapitalGrowthBP: 800,
		RealCapitalGrowthBP:  -600,
		Corresponds:          false,
		Explanation:          "loanable money piles up in banks while productive capital and stocks contract",
	}
	limit := CreditLimit{ReserveCapital: 500_000, Description: "Bank of England banking reserve"}
	a, err := NewRealCapitalAccumulation(
		"Crisis 1847",
		PhaseCrisis,
		corr,
		limit,
		true,  // commercial credit contracts
		false, // cause is NOT money scarcity
		"blocked returns and lost confidence contract commercial credit",
	)
	if err != nil {
		t.Fatalf("fixtureCrisis1847: %v", err)
	}
	return a
}

// fixtureCh30Revival returns a revival observation where the two accumulations align.
func fixtureCh30Revival(t *testing.T) RealCapitalAccumulation {
	t.Helper()
	corr := AccumulationCorrespondence{
		MoneyCapitalGrowthBP: 400,
		RealCapitalGrowthBP:  450,
		Corresponds:          true,
		Explanation:          "loanable funds and productive investment expand together",
	}
	limit := CreditLimit{ReserveCapital: 3_000_000, Description: "ample reserve in the revival phase"}
	a, err := NewRealCapitalAccumulation(
		"Revival",
		PhaseRevival,
		corr,
		limit,
		false, // commercial credit not contracting
		false,
		"money- and real-capital accumulation advance together",
	)
	if err != nil {
		t.Fatalf("fixtureRevival: %v", err)
	}
	return a
}

// ── NewRealCapitalAccumulation happy paths ────────────────────────────────────

func TestNewRealCapitalAccumulation_Crisis1847(t *testing.T) {
	t.Parallel()
	a := fixtureCrisis1847(t)
	if a.Phase != PhaseCrisis {
		t.Errorf("Phase: got %d, want %d", a.Phase, PhaseCrisis)
	}
	if !a.CommercialCreditContracts {
		t.Error("CommercialCreditContracts must be true in 1847")
	}
	if a.CauseIsMoneyScarcity {
		t.Error("CauseIsMoneyScarcity must be false — cause is blocked returns + lost confidence")
	}
	if a.Correspondence.Corresponds {
		t.Error("Correspondence.Corresponds must be false when money-/real-capital growth diverge")
	}
	if a.Correspondence.MoneyCapitalGrowthBP != 800 {
		t.Errorf("MoneyCapitalGrowthBP: got %d, want 800", a.Correspondence.MoneyCapitalGrowthBP)
	}
	if a.Correspondence.RealCapitalGrowthBP != -600 {
		t.Errorf("RealCapitalGrowthBP: got %d, want -600", a.Correspondence.RealCapitalGrowthBP)
	}
	if a.CreditLimit.ReserveCapital != 500_000 {
		t.Errorf("ReserveCapital: got %d, want 500000", a.CreditLimit.ReserveCapital)
	}
	// ID and CreatedAt must be zero — store assigns them.
	if !a.ID.IsZero() {
		t.Error("ID must be zero before store persists")
	}
}

func TestNewRealCapitalAccumulation_Revival(t *testing.T) {
	t.Parallel()
	a := fixtureCh30Revival(t)
	if a.Phase != PhaseRevival {
		t.Errorf("Phase: got %d, want %d", a.Phase, PhaseRevival)
	}
	if a.CommercialCreditContracts {
		t.Error("CommercialCreditContracts must be false in revival")
	}
	if !a.Correspondence.Corresponds {
		t.Error("Correspondence.Corresponds must be true when the two accumulations align")
	}
}

// AccumulationCorrespondence allows the two growth rates to diverge with no
// equality invariant — confirm the constructor does not reject a divergence.
func TestNewRealCapitalAccumulation_DivergenceAllowed(t *testing.T) {
	t.Parallel()
	corr := AccumulationCorrespondence{
		MoneyCapitalGrowthBP: 1000,
		RealCapitalGrowthBP:  -1000,
		Corresponds:          false,
		Explanation:          "extreme divergence",
	}
	limit := CreditLimit{ReserveCapital: 0}
	a, err := NewRealCapitalAccumulation("divergent", PhaseOverproduction, corr, limit, true, false, "")
	if err != nil {
		t.Fatalf("divergence must be allowed: %v", err)
	}
	if a.Correspondence.MoneyCapitalGrowthBP == a.Correspondence.RealCapitalGrowthBP {
		t.Error("growth rates should remain distinct")
	}
}

// ── Invariant failures ────────────────────────────────────────────────────────

func TestNewRealCapitalAccumulation_InvalidPhase(t *testing.T) {
	t.Parallel()
	corr := AccumulationCorrespondence{}
	limit := CreditLimit{ReserveCapital: 0}
	_, err := NewRealCapitalAccumulation("bad", 99, corr, limit, false, false, "")
	if err != ErrInvalidCyclePhase {
		t.Errorf("got %v, want ErrInvalidCyclePhase", err)
	}
}

func TestNewRealCapitalAccumulation_InvalidPhaseZero(t *testing.T) {
	t.Parallel()
	corr := AccumulationCorrespondence{}
	limit := CreditLimit{ReserveCapital: 0}
	_, err := NewRealCapitalAccumulation("bad", 0, corr, limit, false, false, "")
	if err != ErrInvalidCyclePhase {
		t.Errorf("got %v, want ErrInvalidCyclePhase", err)
	}
}

func TestNewRealCapitalAccumulation_NegativeReserve(t *testing.T) {
	t.Parallel()
	corr := AccumulationCorrespondence{}
	limit := CreditLimit{ReserveCapital: -1}
	_, err := NewRealCapitalAccumulation("bad", PhaseInactivity, corr, limit, false, false, "")
	if err != ErrNegativeReserveCapital {
		t.Errorf("got %v, want ErrNegativeReserveCapital", err)
	}
}

// ── CommercialCreditCircuit ───────────────────────────────────────────────────

func TestNewCommercialCreditCircuit_CottonChain(t *testing.T) {
	t.Parallel()
	stages := []string{"spinner", "weaver", "manufacturer", "exporter"}
	c, err := NewCommercialCreditCircuit(stages, 1_000_000, "cotton→yarn→cloth chain of bills")
	if err != nil {
		t.Fatalf("NewCommercialCreditCircuit: %v", err)
	}
	if c.TotalBillsAmount != 1_000_000 {
		t.Errorf("TotalBillsAmount: got %d, want 1000000", c.TotalBillsAmount)
	}
	if c.StageCount() != 4 {
		t.Errorf("StageCount: got %d, want 4", c.StageCount())
	}
}

func TestNewCommercialCreditCircuit_NonPositiveBills(t *testing.T) {
	t.Parallel()
	stages := []string{"spinner", "weaver"}
	if _, err := NewCommercialCreditCircuit(stages, 0, ""); err != ErrNonPositiveBills {
		t.Errorf("zero bills: got %v, want ErrNonPositiveBills", err)
	}
	if _, err := NewCommercialCreditCircuit(stages, -5, ""); err != ErrNonPositiveBills {
		t.Errorf("negative bills: got %v, want ErrNonPositiveBills", err)
	}
}

func TestNewCommercialCreditCircuitID_Format(t *testing.T) {
	t.Parallel()
	id := NewCommercialCreditCircuitID()
	if len(string(id)) != 24 {
		t.Errorf("ID length: got %d, want 24 hex chars", len(string(id)))
	}
}

// ── ReturnFlow ────────────────────────────────────────────────────────────────

func TestNewReturnFlow_OK(t *testing.T) {
	t.Parallel()
	rf, err := NewReturnFlow(30, 30)
	if err != nil {
		t.Fatalf("NewReturnFlow: %v", err)
	}
	if rf.IsBlocked() {
		t.Error("on-time return must not be blocked")
	}
}

func TestNewReturnFlow_Blocked(t *testing.T) {
	t.Parallel()
	rf, err := NewReturnFlow(30, 90)
	if err != nil {
		t.Fatalf("NewReturnFlow: %v", err)
	}
	if !rf.IsBlocked() {
		t.Error("late return must be blocked")
	}
}

func TestNewReturnFlow_Negative(t *testing.T) {
	t.Parallel()
	if _, err := NewReturnFlow(-1, 0); err != ErrNegativeReturnDays {
		t.Errorf("negative expected: got %v, want ErrNegativeReturnDays", err)
	}
	if _, err := NewReturnFlow(0, -1); err != ErrNegativeReturnDays {
		t.Errorf("negative actual: got %v, want ErrNegativeReturnDays", err)
	}
}

// ── RealCapitalAccumulationID ─────────────────────────────────────────────────

func TestNewRealCapitalAccumulationID_Format(t *testing.T) {
	t.Parallel()
	id := NewRealCapitalAccumulationID()
	if len(string(id)) != 24 {
		t.Errorf("ID length: got %d, want 24 hex chars", len(string(id)))
	}
}

func TestNewRealCapitalAccumulationID_Unique(t *testing.T) {
	t.Parallel()
	a := NewRealCapitalAccumulationID()
	b := NewRealCapitalAccumulationID()
	if a == b {
		t.Error("two IDs must not be equal")
	}
}

func TestRealCapitalAccumulationID_IsZero(t *testing.T) {
	t.Parallel()
	var id RealCapitalAccumulationID
	if !id.IsZero() {
		t.Error("zero value must be zero")
	}
	id = NewRealCapitalAccumulationID()
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}

// ── CorrespondenceVerdict ─────────────────────────────────────────────────────

func TestCorrespondenceVerdict_CrisisBlockedReturns(t *testing.T) {
	t.Parallel()
	a := fixtureCrisis1847(t)
	got := a.CorrespondenceVerdict()
	want := "commercial credit contracting — cause is blocked returns and lost confidence, not a money shortage"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCorrespondenceVerdict_CrisisMoneyScarcity(t *testing.T) {
	t.Parallel()
	corr := AccumulationCorrespondence{Corresponds: false}
	limit := CreditLimit{ReserveCapital: 0}
	a, err := NewRealCapitalAccumulation("x", PhaseCrisis, corr, limit, true, true, "")
	if err != nil {
		t.Fatalf("NewRealCapitalAccumulation: %v", err)
	}
	got := a.CorrespondenceVerdict()
	want := "commercial credit contracting — attributed to a scarcity of money-capital"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCorrespondenceVerdict_Corresponds(t *testing.T) {
	t.Parallel()
	a := fixtureCh30Revival(t)
	got := a.CorrespondenceVerdict()
	want := "money-capital and real-capital accumulation correspond"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCorrespondenceVerdict_Diverge(t *testing.T) {
	t.Parallel()
	corr := AccumulationCorrespondence{Corresponds: false}
	limit := CreditLimit{ReserveCapital: 0}
	a, err := NewRealCapitalAccumulation("x", PhaseProsperity, corr, limit, false, false, "")
	if err != nil {
		t.Fatalf("NewRealCapitalAccumulation: %v", err)
	}
	got := a.CorrespondenceVerdict()
	want := "money-capital and real-capital accumulation diverge"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
