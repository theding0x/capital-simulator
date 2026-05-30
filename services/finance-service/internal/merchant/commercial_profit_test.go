package merchant

import (
	"encoding/hex"
	"errors"
	"testing"
)

// --- CommercialProfitID ---

// TestNewCommercialProfitID_HexUnique asserts IDs are non-empty, 24 hex chars
// (96-bit), unique, and that IsZero works correctly.
func TestNewCommercialProfitID_HexUnique(t *testing.T) {
	t.Parallel()

	seen := make(map[CommercialProfitID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewCommercialProfitID()
		if id.IsZero() {
			t.Fatal("NewCommercialProfitID returned the empty id")
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

// TestCommercialProfitID_IsZero asserts IsZero returns true for the empty string
// and false for a freshly generated id.
func TestCommercialProfitID_IsZero(t *testing.T) {
	t.Parallel()

	var empty CommercialProfitID
	if !empty.IsZero() {
		t.Error("empty CommercialProfitID.IsZero() = false, want true")
	}
	id := NewCommercialProfitID()
	if id.IsZero() {
		t.Error("freshly generated CommercialProfitID.IsZero() = true, want false")
	}
}

// --- NewCommercialProfit ---

// TestNewCommercialProfit_MarxFixture covers the 720+180 fixture from Vol. III
// Ch. 17: productive=720000, commercial=180000, surplus=144000.
// unadjusted = round_half_up(144000*10000, 720000) = 2000 bp (20%)
// adjusted   = round_half_up(144000*10000, 900000) = 1600 bp (16%)
// Invariant 2: AdjustedRateBP < UnadjustedRateBP.
// NewValueCreated must be 0.
func TestNewCommercialProfit_MarxFixture(t *testing.T) {
	t.Parallel()

	cp, err := NewCommercialProfit(720000, 180000, 144000)
	if err != nil {
		t.Fatalf("NewCommercialProfit: %v", err)
	}
	if cp.UnadjustedRateBP != 2000 {
		t.Errorf("UnadjustedRateBP = %d, want 2000 (20%%)", cp.UnadjustedRateBP)
	}
	if cp.AdjustedRateBP != 1600 {
		t.Errorf("AdjustedRateBP = %d, want 1600 (16%%)", cp.AdjustedRateBP)
	}
	// Invariant 2: commercial capital dilutes the general rate.
	if cp.AdjustedRateBP >= cp.UnadjustedRateBP {
		t.Errorf("invariant 2 violated: AdjustedRateBP %d >= UnadjustedRateBP %d",
			cp.AdjustedRateBP, cp.UnadjustedRateBP)
	}
	if cp.NewValueCreated != 0 {
		t.Errorf("NewValueCreated = %d, want 0 (commercial activity creates no new value)", cp.NewValueCreated)
	}
	if cp.ProductiveCapital != 720000 {
		t.Errorf("ProductiveCapital = %d, want 720000", cp.ProductiveCapital)
	}
	if cp.CommercialCapital != 180000 {
		t.Errorf("CommercialCapital = %d, want 180000", cp.CommercialCapital)
	}
	if cp.SurplusValue != 144000 {
		t.Errorf("SurplusValue = %d, want 144000", cp.SurplusValue)
	}
	// Constructor assigns no ID or timestamp — the store does that.
	if !cp.ID.IsZero() {
		t.Errorf("ID = %q, want empty (store assigns it)", cp.ID)
	}
	if !cp.CreatedAt.IsZero() {
		t.Errorf("CreatedAt = %v, want zero (store assigns it)", cp.CreatedAt)
	}
}

// TestNewCommercialProfit_RoundHalfUp exercises the round-half-up path using
// productive=800000, commercial=100000, surplus=160000.
// adjusted = round_half_up(160000*10000, 900000) = round_half_up(1600000000, 900000)
// = (1600000000 + 450000) / 900000 = 1600450000 / 900000 = 1778 bp
func TestNewCommercialProfit_RoundHalfUp(t *testing.T) {
	t.Parallel()

	cp, err := NewCommercialProfit(800000, 100000, 160000)
	if err != nil {
		t.Fatalf("NewCommercialProfit: %v", err)
	}
	if cp.AdjustedRateBP != 1778 {
		t.Errorf("AdjustedRateBP = %d, want 1778 (round-half-up)", cp.AdjustedRateBP)
	}
	// unadjusted = round_half_up(160000*10000, 800000) = 2000 bp
	if cp.UnadjustedRateBP != 2000 {
		t.Errorf("UnadjustedRateBP = %d, want 2000", cp.UnadjustedRateBP)
	}
}

// TestNewCommercialProfit_ZeroSurplus verifies that surplus_value == 0 is valid:
// both rates are 0, no error.
func TestNewCommercialProfit_ZeroSurplus(t *testing.T) {
	t.Parallel()

	cp, err := NewCommercialProfit(720000, 180000, 0)
	if err != nil {
		t.Fatalf("NewCommercialProfit(surplus=0): %v", err)
	}
	if cp.UnadjustedRateBP != 0 {
		t.Errorf("UnadjustedRateBP = %d, want 0", cp.UnadjustedRateBP)
	}
	if cp.AdjustedRateBP != 0 {
		t.Errorf("AdjustedRateBP = %d, want 0", cp.AdjustedRateBP)
	}
}

// TestNewCommercialProfit_ErrNonPositiveCapital verifies that productiveCapital == 0
// and commercialCapital == 0 both return ErrNonPositiveCapital.
func TestNewCommercialProfit_ErrNonPositiveCapital(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name      string
		prod, com int64
	}{
		{"productiveCapital=0", 0, 180000},
		{"commercialCapital=0", 720000, 0},
		{"both zero", 0, 0},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := NewCommercialProfit(tc.prod, tc.com, 144000); !errors.Is(err, ErrNonPositiveCapital) {
				t.Errorf("%s: err = %v, want ErrNonPositiveCapital", tc.name, err)
			}
		})
	}
}

// TestNewCommercialProfit_ErrNegativeSurplusValue verifies that surplus_value < 0
// returns ErrNegativeSurplusValue.
func TestNewCommercialProfit_ErrNegativeSurplusValue(t *testing.T) {
	t.Parallel()

	if _, err := NewCommercialProfit(720000, 180000, -1); !errors.Is(err, ErrNegativeSurplusValue) {
		t.Errorf("surplus=-1: err = %v, want ErrNegativeSurplusValue", err)
	}
}

// --- NewAdjustedGeneralRate ---

// TestNewAdjustedGeneralRate_MarxFixture mirrors the 720+180 fixture and
// asserts invariant 2: AdjustedRateBP < UnadjustedRateBP.
func TestNewAdjustedGeneralRate_MarxFixture(t *testing.T) {
	t.Parallel()

	agr, err := NewAdjustedGeneralRate(720000, 180000, 144000)
	if err != nil {
		t.Fatalf("NewAdjustedGeneralRate: %v", err)
	}
	if agr.UnadjustedRateBP != 2000 {
		t.Errorf("UnadjustedRateBP = %d, want 2000", agr.UnadjustedRateBP)
	}
	if agr.AdjustedRateBP != 1600 {
		t.Errorf("AdjustedRateBP = %d, want 1600", agr.AdjustedRateBP)
	}
	// Invariant 2.
	if agr.AdjustedRateBP >= agr.UnadjustedRateBP {
		t.Errorf("invariant 2 violated: AdjustedRateBP %d >= UnadjustedRateBP %d",
			agr.AdjustedRateBP, agr.UnadjustedRateBP)
	}
}

// TestNewAdjustedGeneralRate_ErrNonPositiveCapital verifies sentinel errors.
func TestNewAdjustedGeneralRate_ErrNonPositiveCapital(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name      string
		prod, com int64
	}{
		{"productiveCapital=0", 0, 180000},
		{"commercialCapital=0", 720000, 0},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := NewAdjustedGeneralRate(tc.prod, tc.com, 144000); !errors.Is(err, ErrNonPositiveCapital) {
				t.Errorf("%s: err = %v, want ErrNonPositiveCapital", tc.name, err)
			}
		})
	}
}

// TestNewAdjustedGeneralRate_ErrNegativeSurplusValue verifies the surplus_value sentinel.
func TestNewAdjustedGeneralRate_ErrNegativeSurplusValue(t *testing.T) {
	t.Parallel()

	if _, err := NewAdjustedGeneralRate(720000, 180000, -1); !errors.Is(err, ErrNegativeSurplusValue) {
		t.Errorf("surplus=-1: err = %v, want ErrNegativeSurplusValue", err)
	}
}

// --- NewSurplusValueDeduction ---

// TestNewSurplusValueDeduction_MarxFixture covers the 720+180 split:
// total=144000, productive=720000, commercial=180000, total_capital=900000
// MerchantShare = round_half_up(144000*180000, 900000) = round_half_up(25920000000, 900000) = 28800
// ProductiveCapitalistShare = 144000 - 28800 = 115200
// Invariant 1: conservation (sum == total).
func TestNewSurplusValueDeduction_MarxFixture(t *testing.T) {
	t.Parallel()

	d, err := NewSurplusValueDeduction(144000, 720000, 180000)
	if err != nil {
		t.Fatalf("NewSurplusValueDeduction: %v", err)
	}
	if d.MerchantShare != 28800 {
		t.Errorf("MerchantShare = %d, want 28800", d.MerchantShare)
	}
	if d.ProductiveCapitalistShare != 115200 {
		t.Errorf("ProductiveCapitalistShare = %d, want 115200", d.ProductiveCapitalistShare)
	}
	// Invariant 1: conservation.
	if d.MerchantShare+d.ProductiveCapitalistShare != d.TotalSurplusValue {
		t.Errorf("invariant 1 violated: %d + %d = %d != TotalSurplusValue %d",
			d.MerchantShare, d.ProductiveCapitalistShare,
			d.MerchantShare+d.ProductiveCapitalistShare, d.TotalSurplusValue)
	}
}

// TestNewSurplusValueDeduction_ErrNonPositiveCapital verifies sentinel errors.
func TestNewSurplusValueDeduction_ErrNonPositiveCapital(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name      string
		prod, com int64
	}{
		{"productiveCapital=0", 0, 180000},
		{"commercialCapital=0", 720000, 0},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := NewSurplusValueDeduction(144000, tc.prod, tc.com); !errors.Is(err, ErrNonPositiveCapital) {
				t.Errorf("%s: err = %v, want ErrNonPositiveCapital", tc.name, err)
			}
		})
	}
}

// TestNewSurplusValueDeduction_ErrNegativeSurplusValue verifies the surplus_value sentinel.
func TestNewSurplusValueDeduction_ErrNegativeSurplusValue(t *testing.T) {
	t.Parallel()

	if _, err := NewSurplusValueDeduction(-1, 720000, 180000); !errors.Is(err, ErrNegativeSurplusValue) {
		t.Errorf("total=-1: err = %v, want ErrNegativeSurplusValue", err)
	}
}

// --- NewCommercialWorker ---

// TestNewCommercialWorker_ClerkFixture covers the clerk fixture from Vol. III
// Ch. 17: necessary=14400, workday=28800 → UnpaidLabourMinutes=14400,
// NewValueCreated=0.
func TestNewCommercialWorker_ClerkFixture(t *testing.T) {
	t.Parallel()

	w, err := NewCommercialWorker(14400, 28800)
	if err != nil {
		t.Fatalf("NewCommercialWorker: %v", err)
	}
	if w.UnpaidLabourMinutes != 14400 {
		t.Errorf("UnpaidLabourMinutes = %d, want 14400", w.UnpaidLabourMinutes)
	}
	// Invariant: commercial worker creates no new value.
	if w.NewValueCreated != 0 {
		t.Errorf("NewValueCreated = %d, want 0", w.NewValueCreated)
	}
	if w.NecessaryLabourMinutes != 14400 {
		t.Errorf("NecessaryLabourMinutes = %d, want 14400", w.NecessaryLabourMinutes)
	}
	if w.WorkdayMinutes != 28800 {
		t.Errorf("WorkdayMinutes = %d, want 28800", w.WorkdayMinutes)
	}
}

// TestNewCommercialWorker_EqualWorkday verifies that workday == necessary is
// valid (zero unpaid labour).
func TestNewCommercialWorker_EqualWorkday(t *testing.T) {
	t.Parallel()

	w, err := NewCommercialWorker(14400, 14400)
	if err != nil {
		t.Fatalf("NewCommercialWorker(equal): %v", err)
	}
	if w.UnpaidLabourMinutes != 0 {
		t.Errorf("UnpaidLabourMinutes = %d, want 0", w.UnpaidLabourMinutes)
	}
}

// TestNewCommercialWorker_ErrNegativeLabourMinutes verifies boundary errors:
// workday < necessary returns ErrNegativeLabourMinutes; necessary == 0 also errors.
func TestNewCommercialWorker_ErrNegativeLabourMinutes(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name            string
		necessary, work int64
	}{
		{"workday < necessary", 14400, 14399},
		{"necessary=0", 0, 28800},
		{"necessary<0", -1, 28800},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := NewCommercialWorker(tc.necessary, tc.work); !errors.Is(err, ErrNegativeLabourMinutes) {
				t.Errorf("%s: err = %v, want ErrNegativeLabourMinutes", tc.name, err)
			}
		})
	}
}

// --- NewUnpaidCommercialLabour ---

// TestNewUnpaidCommercialLabour_Valid verifies minutes=14400 is accepted and
// NewValueCreated==0 (invariant 3).
func TestNewUnpaidCommercialLabour_Valid(t *testing.T) {
	t.Parallel()

	u, err := NewUnpaidCommercialLabour(14400)
	if err != nil {
		t.Fatalf("NewUnpaidCommercialLabour: %v", err)
	}
	if u.Minutes != 14400 {
		t.Errorf("Minutes = %d, want 14400", u.Minutes)
	}
	if u.NewValueCreated != 0 {
		t.Errorf("NewValueCreated = %d, want 0", u.NewValueCreated)
	}
}

// TestNewUnpaidCommercialLabour_ZeroMinutes verifies minutes=0 is valid
// (invariant 3: minutes >= 0).
func TestNewUnpaidCommercialLabour_ZeroMinutes(t *testing.T) {
	t.Parallel()

	u, err := NewUnpaidCommercialLabour(0)
	if err != nil {
		t.Fatalf("NewUnpaidCommercialLabour(0): %v", err)
	}
	if u.Minutes != 0 {
		t.Errorf("Minutes = %d, want 0", u.Minutes)
	}
}

// TestNewUnpaidCommercialLabour_ErrNegativeLabourMinutes verifies minutes < 0
// returns ErrNegativeLabourMinutes.
func TestNewUnpaidCommercialLabour_ErrNegativeLabourMinutes(t *testing.T) {
	t.Parallel()

	if _, err := NewUnpaidCommercialLabour(-1); !errors.Is(err, ErrNegativeLabourMinutes) {
		t.Errorf("minutes=-1: err = %v, want ErrNegativeLabourMinutes", err)
	}
}

// --- NewWholesaleRetailSpread ---

// TestNewWholesaleRetailSpread_Valid verifies wholesale=11520, retail=14400
// → SpreadPence=2880.
func TestNewWholesaleRetailSpread_Valid(t *testing.T) {
	t.Parallel()

	s, err := NewWholesaleRetailSpread(11520, 14400)
	if err != nil {
		t.Fatalf("NewWholesaleRetailSpread: %v", err)
	}
	if s.SpreadPence != 2880 {
		t.Errorf("SpreadPence = %d, want 2880", s.SpreadPence)
	}
	if s.WholesalePrice != 11520 {
		t.Errorf("WholesalePrice = %d, want 11520", s.WholesalePrice)
	}
	if s.RetailPrice != 14400 {
		t.Errorf("RetailPrice = %d, want 14400", s.RetailPrice)
	}
}

// TestNewWholesaleRetailSpread_ErrRetailNotGreater verifies that retail <= wholesale
// returns an error.
func TestNewWholesaleRetailSpread_ErrRetailNotGreater(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name            string
		wholesale, retail int64
	}{
		{"retail == wholesale", 11520, 11520},
		{"retail < wholesale", 11520, 11519},
		{"wholesale=0", 0, 14400},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if _, err := NewWholesaleRetailSpread(tc.wholesale, tc.retail); err == nil {
				t.Errorf("%s: expected error, got nil", tc.name)
			}
		})
	}
}
