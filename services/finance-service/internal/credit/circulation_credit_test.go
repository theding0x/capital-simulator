package credit

import (
	"strings"
	"testing"
)

// ── Ch. 33 fixtures ──────────────────────────────────────────────────────────
//
// London Clearing House 1840s: bills worth £150 000 000 settled with only
// £1 000 000 in actual coin. Scotland velocity-9: a £1 note passes through
// 9 hands before returning to its bank. All identifiers carry the ch33 prefix
// so they don't collide with ch31/ch32 fixtures in the shared package.

// ch33FixtureLondonClearing returns the London Clearing House 1840s settlement.
// total_claims=15 000 000 000, money_used=1 000 000 000 → net_balance=14 000 000 000.
func ch33FixtureLondonClearing(t *testing.T) ClearingHouseSettlement {
	t.Helper()
	s, err := NewClearingHouseSettlement(
		"London Clearing House 1840s",
		"Bills of exchange cancel in the clearing house; only the net balance requires actual coin.",
		15_000_000_000,
		1_000_000_000,
	)
	if err != nil {
		t.Fatalf("ch33FixtureLondonClearing: %v", err)
	}
	return s
}

// ── NewClearingHouseSettlement ────────────────────────────────────────────────

func TestNewClearingHouseSettlement_London(t *testing.T) {
	t.Parallel()
	s := ch33FixtureLondonClearing(t)
	if s.TotalClaims != 15_000_000_000 {
		t.Errorf("TotalClaims: got %d, want 15000000000", s.TotalClaims)
	}
	if s.MoneyUsed != 1_000_000_000 {
		t.Errorf("MoneyUsed: got %d, want 1000000000", s.MoneyUsed)
	}
	if s.NetBalance != 14_000_000_000 {
		t.Errorf("NetBalance: got %d, want 14000000000", s.NetBalance)
	}
	// ID and CreatedAt are assigned by the store, not the constructor.
	if !s.ID.IsZero() {
		t.Error("ID must be zero before store persists")
	}
	if !s.CreatedAt.IsZero() {
		t.Error("CreatedAt must be zero before store persists")
	}
}

func TestNewClearingHouseSettlement_MoneyUsedExceedsClaims(t *testing.T) {
	t.Parallel()
	// money_used=1001 > total_claims=1000 → ErrMoneyUsedExceedsClaims.
	if _, err := NewClearingHouseSettlement("x", "", 1000, 1001); err != ErrMoneyUsedExceedsClaims {
		t.Errorf("got %v, want ErrMoneyUsedExceedsClaims", err)
	}
}

func TestNewClearingHouseSettlement_ZeroMoneyUsed(t *testing.T) {
	t.Parallel()
	// money_used=0 → NetBalance == TotalClaims.
	s, err := NewClearingHouseSettlement("zero money", "", 5_000_000, 0)
	if err != nil {
		t.Fatalf("zero money_used must be allowed: %v", err)
	}
	if s.NetBalance != s.TotalClaims {
		t.Errorf("NetBalance: got %d, want %d (== TotalClaims)", s.NetBalance, s.TotalClaims)
	}
}

func TestNewClearingHouseSettlement_MoneyUsedEqualsClaims(t *testing.T) {
	t.Parallel()
	// money_used == total_claims is the boundary: no clearing benefit, NetBalance==0.
	s, err := NewClearingHouseSettlement("full cash", "", 1_000_000, 1_000_000)
	if err != nil {
		t.Fatalf("money_used==total_claims must be allowed: %v", err)
	}
	if s.NetBalance != 0 {
		t.Errorf("NetBalance: got %d, want 0", s.NetBalance)
	}
}

// ── NewCurrencyVelocity ───────────────────────────────────────────────────────

func TestNewCurrencyVelocity_Five(t *testing.T) {
	t.Parallel()
	cv, err := NewCurrencyVelocity(5, "£500 in notes circulating 5×/day → effective supply £2500")
	if err != nil {
		t.Fatalf("velocity=5 must be allowed: %v", err)
	}
	if cv.Velocity != 5 {
		t.Errorf("Velocity: got %d, want 5", cv.Velocity)
	}
}

func TestNewCurrencyVelocity_ScotlandNine(t *testing.T) {
	t.Parallel()
	// Scotland fixture: £1 note passes through 9 hands.
	cv, err := NewCurrencyVelocity(9, "Scottish country-bank note; 9 payments in a week")
	if err != nil {
		t.Fatalf("velocity=9 must be allowed: %v", err)
	}
	if cv.Velocity != 9 {
		t.Errorf("Velocity: got %d, want 9", cv.Velocity)
	}
}

func TestNewCurrencyVelocity_Zero(t *testing.T) {
	t.Parallel()
	if _, err := NewCurrencyVelocity(0, ""); err != ErrVelocityTooLow {
		t.Errorf("velocity=0: got %v, want ErrVelocityTooLow", err)
	}
}

func TestNewCurrencyVelocity_One(t *testing.T) {
	t.Parallel()
	// velocity=1 is the boundary minimum: a note circulates exactly once.
	cv, err := NewCurrencyVelocity(1, "boundary: single circulation")
	if err != nil {
		t.Fatalf("velocity=1 must be allowed: %v", err)
	}
	if cv.Velocity != 1 {
		t.Errorf("Velocity: got %d, want 1", cv.Velocity)
	}
}

// ── NewEffectiveMoneySupply ───────────────────────────────────────────────────

func TestNewEffectiveMoneySupply_FiveCirculations(t *testing.T) {
	t.Parallel()
	// Marx's fixture: £500 actual notes × 5 circulations = £2500 effective supply.
	ems, err := NewEffectiveMoneySupply(50_000, 5)
	if err != nil {
		t.Fatalf("velocity=5 must be allowed: %v", err)
	}
	if ems.Amount != 250_000 {
		t.Errorf("Amount: got %d, want 250000 (50000×5)", ems.Amount)
	}
	if ems.ActualNotesAmount != 50_000 {
		t.Errorf("ActualNotesAmount: got %d, want 50000", ems.ActualNotesAmount)
	}
	if ems.Velocity != 5 {
		t.Errorf("Velocity: got %d, want 5", ems.Velocity)
	}
}

func TestNewEffectiveMoneySupply_Arbitrary(t *testing.T) {
	t.Parallel()
	// 12345 × 7 = 86415.
	ems, err := NewEffectiveMoneySupply(12_345, 7)
	if err != nil {
		t.Fatalf("velocity=7 must be allowed: %v", err)
	}
	if ems.Amount != 86_415 {
		t.Errorf("Amount: got %d, want 86415 (12345×7)", ems.Amount)
	}
}

func TestNewEffectiveMoneySupply_ZeroVelocity(t *testing.T) {
	t.Parallel()
	if _, err := NewEffectiveMoneySupply(50_000, 0); err != ErrVelocityTooLow {
		t.Errorf("velocity=0: got %v, want ErrVelocityTooLow", err)
	}
}

// ── NewCreditContraction ──────────────────────────────────────────────────────

func TestNewCreditContraction_CrisisForces(t *testing.T) {
	t.Parallel()
	// PhaseCrisis with causes=false: constructor forces CausesPaymentCrisis=true.
	cc, err := NewCreditContraction(PhaseCrisis, false)
	if err != nil {
		t.Fatalf("crisis phase must construct: %v", err)
	}
	if !cc.CausesPaymentCrisis {
		t.Error("CausesPaymentCrisis must be forced true in PhaseCrisis")
	}
	if cc.Phase != PhaseCrisis {
		t.Errorf("Phase: got %d, want %d", cc.Phase, PhaseCrisis)
	}
}

func TestNewCreditContraction_ProsperityNoCrisis(t *testing.T) {
	t.Parallel()
	// PhaseProsperity with causes=false: valid, no crisis.
	cc, err := NewCreditContraction(PhaseProsperity, false)
	if err != nil {
		t.Fatalf("prosperity/no-crisis must construct: %v", err)
	}
	if cc.CausesPaymentCrisis {
		t.Error("CausesPaymentCrisis must remain false in PhaseProsperity")
	}
}

func TestNewCreditContraction_RevivalCausesError(t *testing.T) {
	t.Parallel()
	// PhaseRevival with causes=true: invalid — only PhaseCrisis may cause payment crisis.
	if _, err := NewCreditContraction(PhaseRevival, true); err != ErrCrisisMustCausePaymentCrisis {
		t.Errorf("revival+causes=true: got %v, want ErrCrisisMustCausePaymentCrisis", err)
	}
}

func TestNewCreditContraction_InvalidPhase(t *testing.T) {
	t.Parallel()
	if _, err := NewCreditContraction(IndustrialCyclePhase(99), false); err != ErrInvalidCyclePhase {
		t.Errorf("phase=99: got %v, want ErrInvalidCyclePhase", err)
	}
}

// ── ClearingHouseSettlementID format / uniqueness ─────────────────────────────

func TestNewClearingHouseSettlementID_Format(t *testing.T) {
	t.Parallel()
	id := NewClearingHouseSettlementID()
	if len(string(id)) != 24 {
		t.Errorf("ID length: got %d, want 24 hex chars", len(string(id)))
	}
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}

func TestNewClearingHouseSettlementID_Unique(t *testing.T) {
	t.Parallel()
	a := NewClearingHouseSettlementID()
	b := NewClearingHouseSettlementID()
	if a == b {
		t.Error("two IDs must not be equal")
	}
}

func TestClearingHouseSettlementID_IsZero(t *testing.T) {
	t.Parallel()
	var id ClearingHouseSettlementID
	if !id.IsZero() {
		t.Error("zero value must report IsZero==true")
	}
	id = NewClearingHouseSettlementID()
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}

// ── Ch. 33 seed-ID pattern ────────────────────────────────────────────────────

func TestCh33SeedIDs(t *testing.T) {
	t.Parallel()
	ids := []string{
		"5eed000000000000003300",
		"5eed000000000000003301",
	}
	const prefix = "5eed00000000000000"
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if !strings.HasPrefix(id, prefix) {
			t.Errorf("seed id %q lacks prefix %q", id, prefix)
		}
		if cc := id[len(prefix) : len(prefix)+2]; cc != "33" {
			t.Errorf("seed id %q chapter token = %q, want 33", id, cc)
		}
		if _, dup := seen[id]; dup {
			t.Errorf("duplicate seed id %q", id)
		}
		seen[id] = struct{}{}
	}
}
