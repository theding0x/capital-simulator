package merchant

import (
	"encoding/hex"
	"errors"
	"testing"
)

// --- CommercialCapitalFunction constants ---

// TestCommercialCapitalFunction_Values asserts the four function consts are 0/1/2/3
// exactly as the chapter spec requires.
func TestCommercialCapitalFunction_Values(t *testing.T) {
	t.Parallel()

	if FunctionRealisation != 0 {
		t.Errorf("FunctionRealisation = %d, want 0", FunctionRealisation)
	}
	if FunctionStorage != 1 {
		t.Errorf("FunctionStorage = %d, want 1", FunctionStorage)
	}
	if FunctionTransport != 2 {
		t.Errorf("FunctionTransport = %d, want 2", FunctionTransport)
	}
	if FunctionDistribution != 3 {
		t.Errorf("FunctionDistribution = %d, want 3", FunctionDistribution)
	}
}

// --- CommercialCapitalID ---

// TestNewCommercialCapitalID_HexUnique asserts IDs are non-empty, 24 hex chars
// (96-bit), unique, and that IsZero works correctly.
func TestNewCommercialCapitalID_HexUnique(t *testing.T) {
	t.Parallel()

	seen := make(map[CommercialCapitalID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewCommercialCapitalID()
		if id.IsZero() {
			t.Fatal("NewCommercialCapitalID returned the empty id")
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

// TestCommercialCapitalID_IsZero asserts IsZero returns true for the empty string
// and false for a freshly generated id.
func TestCommercialCapitalID_IsZero(t *testing.T) {
	t.Parallel()

	var empty CommercialCapitalID
	if !empty.IsZero() {
		t.Error("empty CommercialCapitalID.IsZero() = false, want true")
	}
	id := NewCommercialCapitalID()
	if id.IsZero() {
		t.Error("freshly generated CommercialCapitalID.IsZero() = true, want false")
	}
}

// --- NewCommercialCapital ---

// TestNewCommercialCapital_LinenMerchant covers §1: Marx's linen merchant advances
// £3,000 (300,000 pence). The entity must carry SurplusValueProduced == 0 (the
// central invariant of commercial capital). The constructor must NOT assign an ID
// or timestamp — the store does that.
func TestNewCommercialCapital_LinenMerchant(t *testing.T) {
	t.Parallel()

	lm := LinenMerchant()
	cc, err := NewCommercialCapital(lm.MoneyAdvanced, lm.Description, FunctionRealisation)
	if err != nil {
		t.Fatalf("NewCommercialCapital: %v", err)
	}
	if cc.MoneyAdvanced != 300000 {
		t.Errorf("MoneyAdvanced = %d, want 300000 (£3,000)", cc.MoneyAdvanced)
	}
	// Invariant: commercial capital creates no surplus-value of its own.
	if cc.SurplusValueProduced != 0 {
		t.Errorf("SurplusValueProduced = %d, want 0 (commercial capital creates no new value)", cc.SurplusValueProduced)
	}
	if cc.Function != FunctionRealisation {
		t.Errorf("Function = %d, want FunctionRealisation (0)", cc.Function)
	}
	// Constructor assigns no ID or timestamp — the store does that.
	if !cc.ID.IsZero() {
		t.Errorf("ID = %q, want empty (store assigns it)", cc.ID)
	}
	if !cc.CreatedAt.IsZero() {
		t.Errorf("CreatedAt = %v, want zero (store assigns it)", cc.CreatedAt)
	}
}

// TestNewCommercialCapital_AllFunctions verifies that all four CommercialCapitalFunction
// values produce SurplusValueProduced == 0, since none of the four functions
// creates new value (§1–§4 invariant).
func TestNewCommercialCapital_AllFunctions(t *testing.T) {
	t.Parallel()

	for _, fn := range []CommercialCapitalFunction{
		FunctionRealisation,
		FunctionStorage,
		FunctionTransport,
		FunctionDistribution,
	} {
		fn := fn
		t.Run("", func(t *testing.T) {
			t.Parallel()
			cc, err := NewCommercialCapital(300000, "linen", fn)
			if err != nil {
				t.Fatalf("NewCommercialCapital(fn=%d): %v", fn, err)
			}
			if cc.SurplusValueProduced != 0 {
				t.Errorf("fn=%d: SurplusValueProduced = %d, want 0", fn, cc.SurplusValueProduced)
			}
		})
	}
}

// TestNewCommercialCapital_NonPositiveMoney verifies that moneyAdvanced == 0 and
// moneyAdvanced < 0 both return ErrNonPositiveMoney.
func TestNewCommercialCapital_NonPositiveMoney(t *testing.T) {
	t.Parallel()

	for _, bad := range []int64{0, -1, -300000} {
		bad := bad
		t.Run("", func(t *testing.T) {
			t.Parallel()
			if _, err := NewCommercialCapital(bad, "linen", FunctionRealisation); !errors.Is(err, ErrNonPositiveMoney) {
				t.Errorf("moneyAdvanced=%d: err = %v, want ErrNonPositiveMoney", bad, err)
			}
		})
	}
}

// --- NewMerchantCircuit ---

// TestNewMerchantCircuit_LinenMerchant covers §1: £3,000 advanced, £3,000 commodity
// value, £3,300 returned → Profit() == 30,000 pence (£300). The profit is a claim
// on industrial surplus-value, not newly created value.
func TestNewMerchantCircuit_LinenMerchant(t *testing.T) {
	t.Parallel()

	mc, err := NewMerchantCircuit(300000, 300000, 330000)
	if err != nil {
		t.Fatalf("NewMerchantCircuit: %v", err)
	}
	if mc.Profit() != 30000 {
		t.Errorf("Profit() = %d, want 30000 (£300)", mc.Profit())
	}
	if mc.MoneyAdvanced != 300000 {
		t.Errorf("MoneyAdvanced = %d, want 300000", mc.MoneyAdvanced)
	}
	if mc.CommodityValue != 300000 {
		t.Errorf("CommodityValue = %d, want 300000", mc.CommodityValue)
	}
	if mc.MoneyReturned != 330000 {
		t.Errorf("MoneyReturned = %d, want 330000", mc.MoneyReturned)
	}
}

// TestNewMerchantCircuit_NoProfit verifies that MoneyReturned <= MoneyAdvanced
// returns ErrNoProfit. The merchant's circuit must extract profit from the
// industrial circuit — equal return and zero return both fail.
func TestNewMerchantCircuit_NoProfit(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name          string
		advanced      int64
		commodityVal  int64
		returned      int64
	}{
		{"equal (no gain)", 300000, 300000, 300000},
		{"loss (returned < advanced)", 300000, 300000, 299999},
		{"zero returned", 300000, 300000, 0},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := NewMerchantCircuit(tc.advanced, tc.commodityVal, tc.returned); !errors.Is(err, ErrNoProfit) {
				t.Errorf("%s: err = %v, want ErrNoProfit", tc.name, err)
			}
		})
	}
}

// TestNewMerchantCircuit_NonPositiveMoney verifies that MoneyAdvanced <= 0 returns
// ErrNonPositiveMoney.
func TestNewMerchantCircuit_NonPositiveMoney(t *testing.T) {
	t.Parallel()

	for _, bad := range []int64{0, -1} {
		bad := bad
		t.Run("", func(t *testing.T) {
			t.Parallel()
			if _, err := NewMerchantCircuit(bad, 300000, 330000); !errors.Is(err, ErrNonPositiveMoney) {
				t.Errorf("moneyAdvanced=%d: err = %v, want ErrNonPositiveMoney", bad, err)
			}
		})
	}
}

// --- NewCommodityCapitalForm ---

// TestNewCommodityCapitalForm_IndustrialHandoff covers §1: industrial handoff
// with c=40000 (£400), v=10000 (£100), s=10000 (£100) → TotalValue()==60000 (£600),
// SourceCapital=="industrial". This is invariant 4: commercial capital can only
// handle what industrial capital has already produced.
func TestNewCommodityCapitalForm_IndustrialHandoff(t *testing.T) {
	t.Parallel()

	f, err := NewCommodityCapitalForm(40000, 10000, 10000)
	if err != nil {
		t.Fatalf("NewCommodityCapitalForm: %v", err)
	}
	if f.TotalValue() != 60000 {
		t.Errorf("TotalValue() = %d, want 60000 (c+v+s=£400+£100+£100)", f.TotalValue())
	}
	if f.SourceCapital != "industrial" {
		t.Errorf("SourceCapital = %q, want %q", f.SourceCapital, "industrial")
	}
	if f.ConstantCapital != 40000 {
		t.Errorf("ConstantCapital = %d, want 40000", f.ConstantCapital)
	}
	if f.VariableCapital != 10000 {
		t.Errorf("VariableCapital = %d, want 10000", f.VariableCapital)
	}
	if f.SurplusValue != 10000 {
		t.Errorf("SurplusValue = %d, want 10000", f.SurplusValue)
	}
}

// TestNewCommodityCapitalForm_ZeroFields verifies that all-zero fields are valid
// (commercial capital may handle commodity-capital where surplus has been
// temporarily zero in an accounting period).
func TestNewCommodityCapitalForm_ZeroFields(t *testing.T) {
	t.Parallel()

	f, err := NewCommodityCapitalForm(0, 0, 0)
	if err != nil {
		t.Fatalf("NewCommodityCapitalForm(0,0,0): %v", err)
	}
	if f.TotalValue() != 0 {
		t.Errorf("TotalValue() = %d, want 0", f.TotalValue())
	}
	if f.SourceCapital != "industrial" {
		t.Errorf("SourceCapital = %q, want %q", f.SourceCapital, "industrial")
	}
}

// TestNewCommodityCapitalForm_NegativeMagnitude verifies that a negative c, v, or s
// returns ErrNegativeMagnitude.
func TestNewCommodityCapitalForm_NegativeMagnitude(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name string
		c, v, s int64
	}{
		{"negative c", -1, 10000, 10000},
		{"negative v", 40000, -1, 10000},
		{"negative s", 40000, 10000, -1},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := NewCommodityCapitalForm(tc.c, tc.v, tc.s); !errors.Is(err, ErrNegativeMagnitude) {
				t.Errorf("%s: err = %v, want ErrNegativeMagnitude", tc.name, err)
			}
		})
	}
}

// TestCommodityCapitalForm_Validate_NonIndustrialSource verifies that
// Validate() returns ErrNonIndustrialSource when SourceCapital is not "industrial".
// (NewCommodityCapitalForm always sets "industrial", but Validate() is the
// contract-enforcing gate.)
func TestCommodityCapitalForm_Validate_NonIndustrialSource(t *testing.T) {
	t.Parallel()

	bad := CommodityCapitalForm{
		ConstantCapital: 40000,
		VariableCapital: 10000,
		SurplusValue:    10000,
		SourceCapital:   "commercial", // wrong — must be industrial
	}
	if err := bad.Validate(); !errors.Is(err, ErrNonIndustrialSource) {
		t.Errorf("Validate() = %v, want ErrNonIndustrialSource", err)
	}
}

// TestLinenMerchant_Fixture cross-checks the LinenMerchant() helper so that
// MoneyAdvanced == YardsBought × PricePerYardPence.
func TestLinenMerchant_Fixture(t *testing.T) {
	t.Parallel()

	lm := LinenMerchant()
	if lm.MoneyAdvanced != lm.YardsBought*lm.PricePerYardPence {
		t.Errorf("MoneyAdvanced %d != YardsBought %d × PricePerYardPence %d = %d",
			lm.MoneyAdvanced, lm.YardsBought, lm.PricePerYardPence,
			lm.YardsBought*lm.PricePerYardPence)
	}
	if lm.MoneyAdvanced != 300000 {
		t.Errorf("MoneyAdvanced = %d, want 300000 (£3,000)", lm.MoneyAdvanced)
	}
}
