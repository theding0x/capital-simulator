package merchant

import (
	"encoding/hex"
	"errors"
	"testing"
)

// --- MoneyDealingCapitalID ---

// TestNewMoneyDealingCapitalID_HexUnique asserts IDs are non-empty, 24 hex chars
// (96-bit), unique, and that IsZero works correctly.
func TestNewMoneyDealingCapitalID_HexUnique(t *testing.T) {
	t.Parallel()

	seen := make(map[MoneyDealingCapitalID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewMoneyDealingCapitalID()
		if id.IsZero() {
			t.Fatal("NewMoneyDealingCapitalID returned the empty id")
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

// TestMoneyDealingCapitalID_IsZero asserts IsZero returns true for the empty
// string and false for a freshly generated id.
func TestMoneyDealingCapitalID_IsZero(t *testing.T) {
	t.Parallel()

	var empty MoneyDealingCapitalID
	if !empty.IsZero() {
		t.Error("empty MoneyDealingCapitalID.IsZero() = false, want true")
	}
	id := NewMoneyDealingCapitalID()
	if id.IsZero() {
		t.Error("freshly generated MoneyDealingCapitalID.IsZero() = true, want false")
	}
}

// --- MoneyOperationKind.IsValid ---

// TestMoneyOperationKind_IsValid asserts that 1–5 are valid and 0 / 6 are not.
func TestMoneyOperationKind_IsValid(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		kind  MoneyOperationKind
		valid bool
	}{
		{0, false},
		{KindReceipt, true},
		{KindPayment, true},
		{KindExchange, true},
		{KindSafekeeping, true},
		{KindBookkeeping, true},
		{6, false},
	} {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			if got := tc.kind.IsValid(); got != tc.valid {
				t.Errorf("MoneyOperationKind(%d).IsValid() = %v, want %v", tc.kind, got, tc.valid)
			}
		})
	}
}

// --- NewMoneyDealingCapital ---

// TestNewMoneyDealingCapital_Seed1901 covers the generic-dealer seed fixture
// (Vol. III Ch. 19, seed 1901):
//
//	moneyAdvanced=500000, generalRateBP=1500, operationKinds=[1,2,3]
//	profit = roundHalfUp(500000*1500, 10000) = 75000
//	Invariant: NewValueCreated == 0
func TestNewMoneyDealingCapital_Seed1901(t *testing.T) {
	t.Parallel()

	kinds := []MoneyOperationKind{KindReceipt, KindPayment, KindExchange}
	md, err := NewMoneyDealingCapital(500000, 1500, kinds)
	if err != nil {
		t.Fatalf("NewMoneyDealingCapital: %v", err)
	}
	if md.MoneyDealingProfit != 75000 {
		t.Errorf("MoneyDealingProfit = %d, want 75000", md.MoneyDealingProfit)
	}
	// Invariant 1: money-dealing capital creates no new value.
	if md.NewValueCreated != 0 {
		t.Errorf("NewValueCreated = %d, want 0 (money-dealing creates no value)", md.NewValueCreated)
	}
	if len(md.OperationKinds) != 3 {
		t.Errorf("len(OperationKinds) = %d, want 3", len(md.OperationKinds))
	}
	// Constructor assigns no ID or timestamp — the store does that.
	if !md.ID.IsZero() {
		t.Errorf("ID = %q, want empty (store assigns it)", md.ID)
	}
	if !md.CreatedAt.IsZero() {
		t.Errorf("CreatedAt = %v, want zero (store assigns it)", md.CreatedAt)
	}
}

// TestNewMoneyDealingCapital_Seed1902 covers seed 1902:
//
//	moneyAdvanced=200000, generalRateBP=1500, operationKinds=[4]
//	profit = roundHalfUp(200000*1500, 10000) = 30000
func TestNewMoneyDealingCapital_Seed1902(t *testing.T) {
	t.Parallel()

	md, err := NewMoneyDealingCapital(200000, 1500, []MoneyOperationKind{KindSafekeeping})
	if err != nil {
		t.Fatalf("NewMoneyDealingCapital: %v", err)
	}
	if md.MoneyDealingProfit != 30000 {
		t.Errorf("MoneyDealingProfit = %d, want 30000", md.MoneyDealingProfit)
	}
	if md.NewValueCreated != 0 {
		t.Errorf("NewValueCreated = %d, want 0", md.NewValueCreated)
	}
}

// TestNewMoneyDealingCapital_Seed1903 covers seed 1903:
//
//	moneyAdvanced=300000, generalRateBP=1500, operationKinds=[3,5]
//	profit = roundHalfUp(300000*1500, 10000) = 45000
func TestNewMoneyDealingCapital_Seed1903(t *testing.T) {
	t.Parallel()

	md, err := NewMoneyDealingCapital(300000, 1500, []MoneyOperationKind{KindExchange, KindBookkeeping})
	if err != nil {
		t.Fatalf("NewMoneyDealingCapital: %v", err)
	}
	if md.MoneyDealingProfit != 45000 {
		t.Errorf("MoneyDealingProfit = %d, want 45000", md.MoneyDealingProfit)
	}
	if md.NewValueCreated != 0 {
		t.Errorf("NewValueCreated = %d, want 0", md.NewValueCreated)
	}
}

// TestNewMoneyDealingCapital_EmptyOperationKinds verifies that an empty
// operationKinds slice is accepted (the implementation iterates over it
// without rejecting it).
func TestNewMoneyDealingCapital_EmptyOperationKinds(t *testing.T) {
	t.Parallel()

	md, err := NewMoneyDealingCapital(500000, 1500, []MoneyOperationKind{})
	if err != nil {
		t.Fatalf("NewMoneyDealingCapital with empty kinds: %v", err)
	}
	if len(md.OperationKinds) != 0 {
		t.Errorf("len(OperationKinds) = %d, want 0", len(md.OperationKinds))
	}
}

// TestNewMoneyDealingCapital_Boundaries checks every sentinel error path.
func TestNewMoneyDealingCapital_Boundaries(t *testing.T) {
	t.Parallel()

	validKinds := []MoneyOperationKind{KindReceipt}

	for _, tc := range []struct {
		name    string
		money   int64
		rateBP  int64
		kinds   []MoneyOperationKind
		wantErr error
	}{
		{"moneyAdvanced=0", 0, 1500, validKinds, ErrNonPositiveMoney},
		{"moneyAdvanced<0", -1, 1500, validKinds, ErrNonPositiveMoney},
		{"generalRateBP=0", 500000, 0, validKinds, ErrNonPositiveCapital},
		{"generalRateBP<0", 500000, -1, validKinds, ErrNonPositiveCapital},
		{"kind=0 invalid", 500000, 1500, []MoneyOperationKind{0}, ErrInvalidOperationKind},
		{"kind=6 invalid", 500000, 1500, []MoneyOperationKind{6}, ErrInvalidOperationKind},
		{"mixed valid+invalid", 500000, 1500, []MoneyOperationKind{KindReceipt, 0}, ErrInvalidOperationKind},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewMoneyDealingCapital(tc.money, tc.rateBP, tc.kinds)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
			}
		})
	}
}

// --- NewMoneyOperation ---

// TestNewMoneyOperation_Valid asserts all five kinds construct without error.
func TestNewMoneyOperation_Valid(t *testing.T) {
	t.Parallel()

	for _, kind := range []MoneyOperationKind{
		KindReceipt, KindPayment, KindExchange, KindSafekeeping, KindBookkeeping,
	} {
		kind := kind
		t.Run("", func(t *testing.T) {
			t.Parallel()
			op, err := NewMoneyOperation(kind, "description")
			if err != nil {
				t.Fatalf("NewMoneyOperation(kind=%d): %v", kind, err)
			}
			if op.Kind != kind {
				t.Errorf("Kind = %d, want %d", op.Kind, kind)
			}
		})
	}
}

// TestNewMoneyOperation_InvalidKind asserts that kind 0 and kind 6 return
// ErrInvalidOperationKind.
func TestNewMoneyOperation_InvalidKind(t *testing.T) {
	t.Parallel()

	for _, bad := range []MoneyOperationKind{0, 6} {
		bad := bad
		t.Run("", func(t *testing.T) {
			t.Parallel()
			_, err := NewMoneyOperation(bad, "desc")
			if !errors.Is(err, ErrInvalidOperationKind) {
				t.Errorf("kind=%d: err = %v, want ErrInvalidOperationKind", bad, err)
			}
		})
	}
}

// --- NewCashReserve ---

// TestNewCashReserve_Valid asserts positive and zero amounts are accepted.
func TestNewCashReserve_Valid(t *testing.T) {
	t.Parallel()

	cr, err := NewCashReserve(50000)
	if err != nil {
		t.Fatalf("NewCashReserve(50000): %v", err)
	}
	if cr.Amount != 50000 {
		t.Errorf("Amount = %d, want 50000", cr.Amount)
	}
}

// TestNewCashReserve_Zero asserts amount=0 is valid (invariant 4: amount >= 0).
func TestNewCashReserve_Zero(t *testing.T) {
	t.Parallel()

	cr, err := NewCashReserve(0)
	if err != nil {
		t.Fatalf("NewCashReserve(0): %v", err)
	}
	if cr.Amount != 0 {
		t.Errorf("Amount = %d, want 0", cr.Amount)
	}
}

// TestNewCashReserve_Negative asserts amount<0 returns ErrNegativeMagnitude.
func TestNewCashReserve_Negative(t *testing.T) {
	t.Parallel()

	_, err := NewCashReserve(-1)
	if !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("err = %v, want ErrNegativeMagnitude", err)
	}
}

// --- NewMoneyDealingCircuit ---

// TestNewMoneyDealingCircuit_Valid covers moneyAdvanced=500000,
// moneyReturned=575000 → Profit()==75000.
func TestNewMoneyDealingCircuit_Valid(t *testing.T) {
	t.Parallel()

	c, err := NewMoneyDealingCircuit(500000, 575000)
	if err != nil {
		t.Fatalf("NewMoneyDealingCircuit: %v", err)
	}
	if c.Profit() != 75000 {
		t.Errorf("Profit() = %d, want 75000", c.Profit())
	}
	if c.MoneyAdvanced != 500000 {
		t.Errorf("MoneyAdvanced = %d, want 500000", c.MoneyAdvanced)
	}
	if c.MoneyReturned != 575000 {
		t.Errorf("MoneyReturned = %d, want 575000", c.MoneyReturned)
	}
}

// TestNewMoneyDealingCircuit_NoProfit asserts moneyReturned<=moneyAdvanced
// returns ErrNoProfit.
func TestNewMoneyDealingCircuit_NoProfit(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name             string
		advanced, returned int64
	}{
		{"equal", 500000, 500000},
		{"returned less", 500000, 499999},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewMoneyDealingCircuit(tc.advanced, tc.returned)
			if !errors.Is(err, ErrNoProfit) {
				t.Errorf("%s: err = %v, want ErrNoProfit", tc.name, err)
			}
		})
	}
}

// TestNewMoneyDealingCircuit_NonPositiveMoney asserts moneyAdvanced<=0 returns
// ErrNonPositiveMoney.
func TestNewMoneyDealingCircuit_NonPositiveMoney(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		advanced int64
	}{
		{"zero", 0},
		{"negative", -1},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewMoneyDealingCircuit(tc.advanced, 575000)
			if !errors.Is(err, ErrNonPositiveMoney) {
				t.Errorf("%s: err = %v, want ErrNonPositiveMoney", tc.name, err)
			}
		})
	}
}
