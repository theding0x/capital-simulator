package credit

import (
	"strings"
	"testing"
)

// ── Ch. 35 fixtures ───────────────────────────────────────────────────────────
//
// BoE gold reserve Nov 1857: £800 000 — the reserve that triggered suspension
// of the Bank Act 1844. London-Hamburg 1847: −150 bp below par (adverse
// balance). All identifiers carry the ch35 prefix to avoid collisions with
// ch34/ch33/ch32/… fixtures in the shared package.

// ch35FixtureBoE1857 returns a GoldReserve for the BoE crisis fixture.
func ch35FixtureBoE1857(t *testing.T) GoldReserve {
	t.Helper()
	g, err := NewGoldReserve("BoE Nov 1857", 800_000_000, "Bank of England gold reserve during the 1857 crisis.")
	if err != nil {
		t.Fatalf("ch35FixtureBoE1857: %v", err)
	}
	return g
}

// ch35FixtureLondonHamburg1847 returns a RateOfExchange for the London-Hamburg
// 1847 adverse rate: −150 basis points below par.
func ch35FixtureLondonHamburg1847() RateOfExchange {
	return NewRateOfExchange("London-Hamburg 1847", -150, "Adverse balance from corn imports.")
}

// ── WorldMoneyFunction.IsValid ────────────────────────────────────────────────

func TestCh35WorldMoneyFunction_IsValid_ValidValues(t *testing.T) {
	t.Parallel()
	cases := []struct {
		f    WorldMoneyFunction
		name string
	}{
		{FunctionWorldMoney, "FunctionWorldMoney (1)"},
		{FunctionReserveFund, "FunctionReserveFund (2)"},
		{FunctionAdjustment, "FunctionAdjustment (3)"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if !tc.f.IsValid() {
				t.Errorf("IsValid() = false, want true for %s", tc.name)
			}
		})
	}
}

func TestCh35WorldMoneyFunction_IsValid_InvalidValues(t *testing.T) {
	t.Parallel()
	cases := []struct {
		f    WorldMoneyFunction
		name string
	}{
		{WorldMoneyFunction(0), "zero"},
		{WorldMoneyFunction(4), "four"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.f.IsValid() {
				t.Errorf("IsValid() = true, want false for %d", int(tc.f))
			}
		})
	}
}

// ── GoldFlowDirection.IsValid ─────────────────────────────────────────────────

func TestCh35GoldFlowDirection_IsValid_ValidValues(t *testing.T) {
	t.Parallel()
	cases := []struct {
		d    GoldFlowDirection
		name string
	}{
		{FlowInward, "FlowInward (1)"},
		{FlowOutward, "FlowOutward (2)"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if !tc.d.IsValid() {
				t.Errorf("IsValid() = false, want true for %s", tc.name)
			}
		})
	}
}

func TestCh35GoldFlowDirection_IsValid_InvalidValues(t *testing.T) {
	t.Parallel()
	cases := []struct {
		d    GoldFlowDirection
		name string
	}{
		{GoldFlowDirection(0), "zero"},
		{GoldFlowDirection(3), "three"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.d.IsValid() {
				t.Errorf("IsValid() = true, want false for %d", int(tc.d))
			}
		})
	}
}

// ── NewGoldReserve ────────────────────────────────────────────────────────────

func TestCh35NewGoldReserve_BoE1857_Ok(t *testing.T) {
	t.Parallel()
	g, err := NewGoldReserve("BoE Nov 1857", 800_000_000, "Crisis fixture.")
	if err != nil {
		t.Fatalf("amount=800000000 must be allowed: %v", err)
	}
	if g.Amount != 800_000_000 {
		t.Errorf("Amount: got %d, want 800000000", g.Amount)
	}
	// ID and CreatedAt are assigned by the store, not the constructor.
	if !g.ID.IsZero() {
		t.Error("ID must be zero before store persists")
	}
	if !g.CreatedAt.IsZero() {
		t.Error("CreatedAt must be zero before store persists")
	}
}

func TestCh35NewGoldReserve_ZeroAmount_Ok(t *testing.T) {
	t.Parallel()
	g, err := NewGoldReserve("empty reserve", 0, "")
	if err != nil {
		t.Fatalf("amount=0 must be allowed: %v", err)
	}
	if g.Amount != 0 {
		t.Errorf("Amount: got %d, want 0", g.Amount)
	}
}

func TestCh35NewGoldReserve_NegativeAmount_Errors(t *testing.T) {
	t.Parallel()
	if _, err := NewGoldReserve("bad", -1, ""); err != ErrNegativeGoldReserveAmount {
		t.Errorf("amount=-1: got %v, want ErrNegativeGoldReserveAmount", err)
	}
}

// ── NewRateOfExchange ─────────────────────────────────────────────────────────

func TestCh35NewRateOfExchange_NegativeDeviation_Ok(t *testing.T) {
	t.Parallel()
	// London-Hamburg 1847: adverse balance → −150 bp below par.
	re := NewRateOfExchange("London-Hamburg 1847", -150, "Adverse balance from corn imports.")
	if re.DeviationBP != -150 {
		t.Errorf("DeviationBP: got %d, want -150", re.DeviationBP)
	}
}

func TestCh35NewRateOfExchange_PositiveDeviation_Ok(t *testing.T) {
	t.Parallel()
	re := NewRateOfExchange("favourable rate", 200, "Favourable balance of payments.")
	if re.DeviationBP != 200 {
		t.Errorf("DeviationBP: got %d, want 200", re.DeviationBP)
	}
}

func TestCh35NewRateOfExchange_ZeroDeviation_Ok(t *testing.T) {
	t.Parallel()
	re := NewRateOfExchange("at par", 0, "")
	if re.DeviationBP != 0 {
		t.Errorf("DeviationBP: got %d, want 0", re.DeviationBP)
	}
}

// ── NewBalanceOfPayments ──────────────────────────────────────────────────────

func TestCh35NewBalanceOfPayments_NegativeNet_Unfavourable(t *testing.T) {
	t.Parallel()
	// Adverse balance: net deficit → Favourable==false.
	b := NewBalanceOfPayments(-5_000_000, "1847 corn-import deficit.")
	if b.Favourable {
		t.Error("Favourable: got true, want false for negative NetBalance")
	}
	if b.NetBalance != -5_000_000 {
		t.Errorf("NetBalance: got %d, want -5000000", b.NetBalance)
	}
	if b.Favourable != (b.NetBalance > 0) {
		t.Errorf("Favourable invariant violated: Favourable=%v but NetBalance=%d", b.Favourable, b.NetBalance)
	}
}

func TestCh35NewBalanceOfPayments_PositiveNet_Favourable(t *testing.T) {
	t.Parallel()
	// Surplus balance: net credit → Favourable==true.
	b := NewBalanceOfPayments(3_000_000, "Export surplus.")
	if !b.Favourable {
		t.Error("Favourable: got false, want true for positive NetBalance")
	}
	if b.NetBalance != 3_000_000 {
		t.Errorf("NetBalance: got %d, want 3000000", b.NetBalance)
	}
	if b.Favourable != (b.NetBalance > 0) {
		t.Errorf("Favourable invariant violated: Favourable=%v but NetBalance=%d", b.Favourable, b.NetBalance)
	}
}

func TestCh35NewBalanceOfPayments_ZeroNet_Unfavourable(t *testing.T) {
	t.Parallel()
	// Zero is not strictly positive — Favourable must be false.
	b := NewBalanceOfPayments(0, "balanced payments")
	if b.Favourable {
		t.Error("Favourable: got true, want false for zero NetBalance (strict >0 required)")
	}
	if b.Favourable != (b.NetBalance > 0) {
		t.Errorf("Favourable invariant violated: Favourable=%v but NetBalance=%d", b.Favourable, b.NetBalance)
	}
}

// ── NewGoldFlow ───────────────────────────────────────────────────────────────

func TestCh35NewGoldFlow_Outward_Ok(t *testing.T) {
	t.Parallel()
	// Adverse balance drains gold outward.
	gf, err := NewGoldFlow(FlowOutward, 200_000_000, "1847 gold drain to Hamburg.")
	if err != nil {
		t.Fatalf("FlowOutward, amount=200000000 must be allowed: %v", err)
	}
	if gf.Direction != FlowOutward {
		t.Errorf("Direction: got %d, want FlowOutward (%d)", gf.Direction, FlowOutward)
	}
	if gf.Amount != 200_000_000 {
		t.Errorf("Amount: got %d, want 200000000", gf.Amount)
	}
}

func TestCh35NewGoldFlow_InwardZeroAmount_Ok(t *testing.T) {
	t.Parallel()
	gf, err := NewGoldFlow(FlowInward, 0, "no movement")
	if err != nil {
		t.Fatalf("FlowInward, amount=0 must be allowed: %v", err)
	}
	if gf.Amount != 0 {
		t.Errorf("Amount: got %d, want 0", gf.Amount)
	}
}

func TestCh35NewGoldFlow_NegativeAmount_Errors(t *testing.T) {
	t.Parallel()
	if _, err := NewGoldFlow(FlowOutward, -1, ""); err != ErrNegativeGoldFlowAmount {
		t.Errorf("amount=-1: got %v, want ErrNegativeGoldFlowAmount", err)
	}
}

func TestCh35NewGoldFlow_InvalidDirection_Errors(t *testing.T) {
	t.Parallel()
	if _, err := NewGoldFlow(GoldFlowDirection(0), 100, ""); err != ErrInvalidGoldFlowDirection {
		t.Errorf("direction=0: got %v, want ErrInvalidGoldFlowDirection", err)
	}
}

// ── ID format and uniqueness ──────────────────────────────────────────────────

func TestCh35NewGoldReserveID_Format(t *testing.T) {
	t.Parallel()
	id := NewGoldReserveID()
	if len(string(id)) != 24 {
		t.Errorf("GoldReserveID length: got %d, want 24 hex chars", len(string(id)))
	}
	if id.IsZero() {
		t.Error("fresh GoldReserveID must not be zero")
	}
}

func TestCh35NewGoldReserveID_Unique(t *testing.T) {
	t.Parallel()
	a := NewGoldReserveID()
	b := NewGoldReserveID()
	if a == b {
		t.Error("two GoldReserveIDs must not be equal")
	}
}

func TestCh35GoldReserveID_IsZero(t *testing.T) {
	t.Parallel()
	var id GoldReserveID
	if !id.IsZero() {
		t.Error("zero-value GoldReserveID must report IsZero==true")
	}
}

func TestCh35NewRateOfExchangeID_Format(t *testing.T) {
	t.Parallel()
	id := NewRateOfExchangeID()
	if len(string(id)) != 24 {
		t.Errorf("RateOfExchangeID length: got %d, want 24 hex chars", len(string(id)))
	}
	if id.IsZero() {
		t.Error("fresh RateOfExchangeID must not be zero")
	}
}

func TestCh35NewRateOfExchangeID_Unique(t *testing.T) {
	t.Parallel()
	a := NewRateOfExchangeID()
	b := NewRateOfExchangeID()
	if a == b {
		t.Error("two RateOfExchangeIDs must not be equal")
	}
}

func TestCh35RateOfExchangeID_IsZero(t *testing.T) {
	t.Parallel()
	var id RateOfExchangeID
	if !id.IsZero() {
		t.Error("zero-value RateOfExchangeID must report IsZero==true")
	}
}

// ── Ch. 35 seed-ID pattern ────────────────────────────────────────────────────

func TestCh35SeedIDs(t *testing.T) {
	t.Parallel()
	ids := []string{
		"5eed000000000000003500",
		"5eed000000000000003501",
		"5eed000000000000003510",
		"5eed000000000000003511",
	}
	const prefix = "5eed00000000000000"
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if !strings.HasPrefix(id, prefix) {
			t.Errorf("seed id %q lacks prefix %q", id, prefix)
		}
		// The two chars immediately after the prefix are the chapter token.
		if cc := id[len(prefix) : len(prefix)+2]; cc != "35" {
			t.Errorf("seed id %q chapter token = %q, want 35", id, cc)
		}
		if len(id) != 22 {
			t.Errorf("seed id %q length = %d, want 22", id, len(id))
		}
		if _, dup := seen[id]; dup {
			t.Errorf("duplicate seed id %q", id)
		}
		seen[id] = struct{}{}
	}
}
