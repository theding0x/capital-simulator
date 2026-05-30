package credit

import (
	"encoding/hex"
	"errors"
	"testing"
)

// --- InterestBearingCapitalID ---

// TestNewInterestBearingCapitalID_HexUnique asserts IDs are non-empty, 24 hex
// chars (96-bit), unique, and that IsZero works correctly.
func TestNewInterestBearingCapitalID_HexUnique(t *testing.T) {
	t.Parallel()

	seen := make(map[InterestBearingCapitalID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewInterestBearingCapitalID()
		if id.IsZero() {
			t.Fatal("NewInterestBearingCapitalID returned the empty id")
		}
		if len(string(id)) != 24 {
			t.Fatalf("id %q has len %d, want 24 hex chars", id, len(string(id)))
		}
		if _, err := hex.DecodeString(string(id)); err != nil {
			t.Fatalf("id %q is not valid hex: %v", id, err)
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id generated: %q", id)
		}
		seen[id] = struct{}{}
	}
}

// TestInterestBearingCapitalID_IsZero asserts IsZero returns true for the
// empty string and false for a freshly generated id.
func TestInterestBearingCapitalID_IsZero(t *testing.T) {
	t.Parallel()

	var empty InterestBearingCapitalID
	if !empty.IsZero() {
		t.Error("empty InterestBearingCapitalID.IsZero() = false, want true")
	}
	id := NewInterestBearingCapitalID()
	if id.IsZero() {
		t.Error("freshly generated InterestBearingCapitalID.IsZero() = true, want false")
	}
}

// --- NewInterestBearingCapital ---

// TestNewInterestBearingCapital_Marx1000Pounds covers the canonical Marx
// fixture (Vol. III Ch. 21):
//
//	£1,000 at 5% for one year:
//	MoneyAdvanced=100000 pence, InterestRateBP=500, LoanTermDays=365
//	InterestEarned = roundHalfUp(100000*500, 10000) = 5000
//	MoneyReturned  = 100000 + 5000 = 105000
//	NewValueCreated = 0
func TestNewInterestBearingCapital_Marx1000Pounds(t *testing.T) {
	t.Parallel()

	ibc, err := NewInterestBearingCapital(100000, 500, 365)
	if err != nil {
		t.Fatalf("NewInterestBearingCapital: %v", err)
	}
	if ibc.InterestEarned != 5000 {
		t.Errorf("InterestEarned = %d, want 5000", ibc.InterestEarned)
	}
	if ibc.MoneyReturned != 105000 {
		t.Errorf("MoneyReturned = %d, want 105000", ibc.MoneyReturned)
	}
	// Invariant 5: interest-bearing capital creates no new value.
	if ibc.NewValueCreated != 0 {
		t.Errorf("NewValueCreated = %d, want 0", ibc.NewValueCreated)
	}
	// Invariant 2: MoneyReturned == MoneyAdvanced + InterestEarned.
	if ibc.MoneyReturned != ibc.MoneyAdvanced+ibc.InterestEarned {
		t.Errorf("MoneyReturned (%d) != MoneyAdvanced (%d) + InterestEarned (%d)",
			ibc.MoneyReturned, ibc.MoneyAdvanced, ibc.InterestEarned)
	}
	// Constructor assigns no ID or timestamp — the store does that.
	if !ibc.ID.IsZero() {
		t.Errorf("ID = %q, want empty (store assigns it)", ibc.ID)
	}
	if !ibc.CreatedAt.IsZero() {
		t.Errorf("CreatedAt = %v, want zero (store assigns it)", ibc.CreatedAt)
	}
}

// TestNewInterestBearingCapital_Seed2101 covers seed 2101:
//
//	£5,000 at 15% for 180 days:
//	moneyAdvanced=500000, interestRateBP=1500, loanTermDays=180
//	InterestEarned = roundHalfUp(500000*1500, 10000)
//	             = roundHalfUp(750000000, 10000) = 75000
//	MoneyReturned  = 500000 + 75000 = 575000
func TestNewInterestBearingCapital_Seed2101(t *testing.T) {
	t.Parallel()

	ibc, err := NewInterestBearingCapital(500000, 1500, 180)
	if err != nil {
		t.Fatalf("NewInterestBearingCapital: %v", err)
	}
	if ibc.InterestEarned != 75000 {
		t.Errorf("InterestEarned = %d, want 75000", ibc.InterestEarned)
	}
	if ibc.MoneyReturned != 575000 {
		t.Errorf("MoneyReturned = %d, want 575000", ibc.MoneyReturned)
	}
	if ibc.NewValueCreated != 0 {
		t.Errorf("NewValueCreated = %d, want 0", ibc.NewValueCreated)
	}
}

// TestNewInterestBearingCapital_Invariant1 verifies InterestEarned =
// roundHalfUp(MoneyAdvanced*InterestRateBP, 10000).
func TestNewInterestBearingCapital_Invariant1(t *testing.T) {
	t.Parallel()

	cases := []struct {
		money, rateBP, term int64
		wantInterest        int64
	}{
		{100000, 500, 365, 5000},  // £1,000 @ 5%
		{200000, 1000, 90, 20000}, // £2,000 @ 10%
		{300000, 333, 365, 9990},  // 3.33% — tests rounding: 300000*333=99900000; (99900000+5000)/10000=9990
		{100000, 1, 30, 10},       // minimum rate, 1bp: 100000*1=100000; roundHalfUp(100000, 10000)=(100000+5000)/10000=10
	}
	for _, tc := range cases {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			got, err := NewInterestBearingCapital(tc.money, tc.rateBP, tc.term)
			if err != nil {
				t.Fatalf("NewInterestBearingCapital: %v", err)
			}
			if got.InterestEarned != tc.wantInterest {
				t.Errorf("money=%d rateBP=%d: InterestEarned=%d, want %d",
					tc.money, tc.rateBP, got.InterestEarned, tc.wantInterest)
			}
		})
	}
}

// TestNewInterestBearingCapital_ValidationErrors checks every sentinel error
// path.
func TestNewInterestBearingCapital_ValidationErrors(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		money   int64
		rateBP  int64
		term    int64
		wantErr error
	}{
		{"moneyAdvanced=0", 0, 500, 365, ErrNonPositiveMoney},
		{"moneyAdvanced<0", -1, 500, 365, ErrNonPositiveMoney},
		{"interestRateBP=0", 100000, 0, 365, ErrNonPositiveRate},
		{"interestRateBP<0", 100000, -1, 365, ErrNonPositiveRate},
		{"loanTermDays=0", 100000, 500, 0, ErrNonPositiveTerm},
		{"loanTermDays<0", 100000, 500, -1, ErrNonPositiveTerm},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewInterestBearingCapital(tc.money, tc.rateBP, tc.term)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
			}
		})
	}
}

// --- NewMoneyCapitalCircuit ---

// TestNewMoneyCapitalCircuit_Marx1000Pounds covers the M—M′ circuit
// round-trip for the canonical Marx fixture.
func TestNewMoneyCapitalCircuit_Marx1000Pounds(t *testing.T) {
	t.Parallel()

	c, err := NewMoneyCapitalCircuit(100000, 5000, 105000)
	if err != nil {
		t.Fatalf("NewMoneyCapitalCircuit: %v", err)
	}
	if c.MoneyAdvanced != 100000 {
		t.Errorf("MoneyAdvanced = %d, want 100000", c.MoneyAdvanced)
	}
	if c.InterestEarned != 5000 {
		t.Errorf("InterestEarned = %d, want 5000", c.InterestEarned)
	}
	if c.MoneyReturned != 105000 {
		t.Errorf("MoneyReturned = %d, want 105000", c.MoneyReturned)
	}
}

// TestNewMoneyCapitalCircuit_InvariantViolation asserts that MoneyReturned ≠
// MoneyAdvanced+InterestEarned is rejected.
func TestNewMoneyCapitalCircuit_InvariantViolation(t *testing.T) {
	t.Parallel()

	_, err := NewMoneyCapitalCircuit(100000, 5000, 104999) // off by 1
	if err == nil {
		t.Error("expected error for MoneyReturned != MoneyAdvanced+InterestEarned, got nil")
	}
}

// TestNewMoneyCapitalCircuit_ZeroMoneyAdvanced asserts ErrNonPositiveMoney.
func TestNewMoneyCapitalCircuit_ZeroMoneyAdvanced(t *testing.T) {
	t.Parallel()

	_, err := NewMoneyCapitalCircuit(0, 0, 0)
	if !errors.Is(err, ErrNonPositiveMoney) {
		t.Errorf("err = %v, want ErrNonPositiveMoney", err)
	}
}

// --- NewInterest ---

// TestNewInterest_Valid asserts interest < average profit is accepted.
func TestNewInterest_Valid(t *testing.T) {
	t.Parallel()

	i, err := NewInterest(5000, 10000)
	if err != nil {
		t.Fatalf("NewInterest: %v", err)
	}
	if i.Value != 5000 {
		t.Errorf("Value = %d, want 5000", i.Value)
	}
}

// TestNewInterest_EqualToProfit rejects interest == average profit.
func TestNewInterest_EqualToProfit(t *testing.T) {
	t.Parallel()

	_, err := NewInterest(10000, 10000)
	if !errors.Is(err, ErrInterestExceedsProfit) {
		t.Errorf("err = %v, want ErrInterestExceedsProfit", err)
	}
}

// TestNewInterest_ExceedsProfit rejects interest > average profit.
func TestNewInterest_ExceedsProfit(t *testing.T) {
	t.Parallel()

	_, err := NewInterest(10001, 10000)
	if !errors.Is(err, ErrInterestExceedsProfit) {
		t.Errorf("err = %v, want ErrInterestExceedsProfit", err)
	}
}

// --- NewLoanableCapital ---

// TestNewLoanableCapital_Valid asserts positive amount is accepted.
func TestNewLoanableCapital_Valid(t *testing.T) {
	t.Parallel()

	lc, err := NewLoanableCapital(1000000)
	if err != nil {
		t.Fatalf("NewLoanableCapital: %v", err)
	}
	if lc.Amount != 1000000 {
		t.Errorf("Amount = %d, want 1000000", lc.Amount)
	}
}

// TestNewLoanableCapital_NonPositive asserts ≤ 0 returns ErrNonPositiveMoney.
func TestNewLoanableCapital_NonPositive(t *testing.T) {
	t.Parallel()

	for _, v := range []int64{0, -1} {
		v := v
		t.Run("", func(t *testing.T) {
			t.Parallel()
			_, err := NewLoanableCapital(v)
			if !errors.Is(err, ErrNonPositiveMoney) {
				t.Errorf("amount=%d: err = %v, want ErrNonPositiveMoney", v, err)
			}
		})
	}
}

// --- NewLoanTerm ---

// TestNewLoanTerm_Valid asserts positive days are accepted.
func TestNewLoanTerm_Valid(t *testing.T) {
	t.Parallel()

	lt, err := NewLoanTerm(365)
	if err != nil {
		t.Fatalf("NewLoanTerm: %v", err)
	}
	if lt.Days != 365 {
		t.Errorf("Days = %d, want 365", lt.Days)
	}
}

// TestNewLoanTerm_NonPositive asserts ≤ 0 returns ErrNonPositiveTerm.
func TestNewLoanTerm_NonPositive(t *testing.T) {
	t.Parallel()

	for _, v := range []int64{0, -1} {
		v := v
		t.Run("", func(t *testing.T) {
			t.Parallel()
			_, err := NewLoanTerm(v)
			if !errors.Is(err, ErrNonPositiveTerm) {
				t.Errorf("days=%d: err = %v, want ErrNonPositiveTerm", v, err)
			}
		})
	}
}

// --- MoneyCapitalist / FunctioningCapitalist ---

// TestMoneyCapitalist_EarnsInterest asserts a money-capitalist earns interest,
// not profit of enterprise.
func TestMoneyCapitalist_EarnsInterest(t *testing.T) {
	t.Parallel()

	mc := MoneyCapitalist{Role: "lender", EarnsInterest: true}
	if !mc.EarnsInterest {
		t.Error("MoneyCapitalist.EarnsInterest = false, want true")
	}
}

// TestFunctioningCapitalist_DoesNotEarnInterest asserts a functioning
// capitalist earns profit of enterprise, not interest.
func TestFunctioningCapitalist_DoesNotEarnInterest(t *testing.T) {
	t.Parallel()

	fc := FunctioningCapitalist{Role: "industrial capitalist", EarnsInterest: false}
	if fc.EarnsInterest {
		t.Error("FunctioningCapitalist.EarnsInterest = true, want false")
	}
}

// TestNewValueCreated_AlwaysZero asserts NewValueCreated == 0 for all valid
// inputs, asserting the core Marxian invariant.
func TestNewValueCreated_AlwaysZero(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		money, rateBP, term int64
	}{
		{100000, 500, 365},
		{500000, 1500, 180},
		{1, 1, 1},
		{999999999, 10000, 1000},
	} {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			ibc, err := NewInterestBearingCapital(tc.money, tc.rateBP, tc.term)
			if err != nil {
				t.Fatalf("NewInterestBearingCapital: %v", err)
			}
			if ibc.NewValueCreated != 0 {
				t.Errorf("NewValueCreated = %d, want 0 (interest-bearing capital creates no value)",
					ibc.NewValueCreated)
			}
		})
	}
}
