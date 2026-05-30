package merchant

import (
	"encoding/hex"
	"errors"
	"testing"
)

// --- MerchantTurnoverID ---

// TestNewMerchantTurnoverID_HexUnique asserts IDs are non-empty, 24 hex chars
// (96-bit), unique, and that IsZero works correctly.
func TestNewMerchantTurnoverID_HexUnique(t *testing.T) {
	t.Parallel()

	seen := make(map[MerchantTurnoverID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewMerchantTurnoverID()
		if id.IsZero() {
			t.Fatal("NewMerchantTurnoverID returned the empty id")
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

// TestMerchantTurnoverID_IsZero asserts IsZero returns true for the empty string
// and false for a freshly generated id.
func TestMerchantTurnoverID_IsZero(t *testing.T) {
	t.Parallel()

	var empty MerchantTurnoverID
	if !empty.IsZero() {
		t.Error("empty MerchantTurnoverID.IsZero() = false, want true")
	}
	id := NewMerchantTurnoverID()
	if id.IsZero() {
		t.Error("freshly generated MerchantTurnoverID.IsZero() = true, want false")
	}
}

// --- NewMerchantTurnover ---

// TestNewMerchantTurnover_1Turnover covers the sugar-merchant 1-turnover fixture
// (Vol. III Ch. 18): moneyAdvanced=100000p, generalRateBP=1500, turnoverCount=1,
// unitsPerTurnover=300.
//
//	annual = roundHalfUp(100000*1500, 10000) = roundHalfUp(150000000, 10000) = 15000
//	markup = 15000 / (1 * 300) = 50
func TestNewMerchantTurnover_1Turnover(t *testing.T) {
	t.Parallel()

	mt, err := NewMerchantTurnover(100000, 1500, 1, 300)
	if err != nil {
		t.Fatalf("NewMerchantTurnover: %v", err)
	}
	if mt.AnnualMerchantProfit != 15000 {
		t.Errorf("AnnualMerchantProfit = %d, want 15000", mt.AnnualMerchantProfit)
	}
	if mt.MarkupPerUnit != 50 {
		t.Errorf("MarkupPerUnit = %d, want 50", mt.MarkupPerUnit)
	}
	// Constructor assigns no ID or timestamp — the store does that.
	if !mt.ID.IsZero() {
		t.Errorf("ID = %q, want empty (store assigns it)", mt.ID)
	}
	if !mt.CreatedAt.IsZero() {
		t.Errorf("CreatedAt = %v, want zero (store assigns it)", mt.CreatedAt)
	}
}

// TestNewMerchantTurnover_5Turnover covers the 5-turnover variant.
//
//	annual = 15000 (UNCHANGED — invariant 1: profit tied to capital not units)
//	markup = 15000 / (5 * 300) = 10
//
// Invariant 2: 1-turnover markup (50) == 5 × 5-turnover markup (10).
func TestNewMerchantTurnover_5Turnover(t *testing.T) {
	t.Parallel()

	mt1, err := NewMerchantTurnover(100000, 1500, 1, 300)
	if err != nil {
		t.Fatalf("NewMerchantTurnover (1-turn): %v", err)
	}
	mt5, err := NewMerchantTurnover(100000, 1500, 5, 300)
	if err != nil {
		t.Fatalf("NewMerchantTurnover (5-turn): %v", err)
	}

	// Invariant 1: annual profit is invariant to turnover count.
	if mt5.AnnualMerchantProfit != 15000 {
		t.Errorf("AnnualMerchantProfit (5-turn) = %d, want 15000 (profit invariant to turnover)", mt5.AnnualMerchantProfit)
	}
	if mt5.AnnualMerchantProfit != mt1.AnnualMerchantProfit {
		t.Errorf("invariant 1 violated: 1-turn profit %d != 5-turn profit %d",
			mt1.AnnualMerchantProfit, mt5.AnnualMerchantProfit)
	}
	if mt5.MarkupPerUnit != 10 {
		t.Errorf("MarkupPerUnit (5-turn) = %d, want 10", mt5.MarkupPerUnit)
	}
	// Invariant 2: markup is inversely proportional to turnover count.
	if mt1.MarkupPerUnit != 5*mt5.MarkupPerUnit {
		t.Errorf("invariant 2 violated: 1-turn markup %d != 5 × 5-turn markup %d",
			mt1.MarkupPerUnit, 5*mt5.MarkupPerUnit)
	}
}

// TestNewMerchantTurnover_3Turnover verifies integer division truncation:
//
//	markup = 15000 / (3 * 300) = 15000 / 900 = 16 (remainder 600 discarded).
func TestNewMerchantTurnover_3Turnover(t *testing.T) {
	t.Parallel()

	mt, err := NewMerchantTurnover(100000, 1500, 3, 300)
	if err != nil {
		t.Fatalf("NewMerchantTurnover (3-turn): %v", err)
	}
	// Lock the truncation behaviour: must be exactly 16, not 17.
	if mt.MarkupPerUnit != 16 {
		t.Errorf("MarkupPerUnit (3-turn) = %d, want 16 (integer div truncation)", mt.MarkupPerUnit)
	}
	if mt.AnnualMerchantProfit != 15000 {
		t.Errorf("AnnualMerchantProfit (3-turn) = %d, want 15000", mt.AnnualMerchantProfit)
	}
}

// TestNewMerchantTurnover_Boundaries checks that every sentinel error is returned
// for its corresponding invalid input.
func TestNewMerchantTurnover_Boundaries(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name                        string
		money, rateBP, count, units int64
		wantErr                     error
	}{
		{"moneyAdvanced=0", 0, 1500, 1, 300, ErrNonPositiveCapital},
		{"generalRateBP=0", 100000, 0, 1, 300, ErrNonPositiveCapital},
		{"turnoverCount=0", 100000, 1500, 0, 300, ErrInvalidTurnoverCount},
		{"unitsPerTurnover=0", 100000, 1500, 1, 0, ErrNoUnits},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewMerchantTurnover(tc.money, tc.rateBP, tc.count, tc.units)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
			}
		})
	}
}

// --- NewTurnoverEffectAnalysis ---

// TestNewTurnoverEffectAnalysis_MarxFixture covers the [1,5]-turnover analysis.
//
//	moneyAdvanced=100000, generalRateBP=1500, unitsPerTurnover=300, counts=[1,5]
//	annual = 15000
//	Points[0]: count=1, total=300,  markup=50
//	Points[1]: count=5, total=1500, markup=10
func TestNewTurnoverEffectAnalysis_MarxFixture(t *testing.T) {
	t.Parallel()

	a, err := NewTurnoverEffectAnalysis(100000, 1500, 300, []int64{1, 5})
	if err != nil {
		t.Fatalf("NewTurnoverEffectAnalysis: %v", err)
	}
	if a.AnnualMerchantProfit != 15000 {
		t.Errorf("AnnualMerchantProfit = %d, want 15000", a.AnnualMerchantProfit)
	}
	if a.MoneyAdvanced != 100000 {
		t.Errorf("MoneyAdvanced = %d, want 100000", a.MoneyAdvanced)
	}
	if a.GeneralRateBP != 1500 {
		t.Errorf("GeneralRateBP = %d, want 1500", a.GeneralRateBP)
	}
	if len(a.Points) != 2 {
		t.Fatalf("len(Points) = %d, want 2", len(a.Points))
	}
	// Points[0]: 1-turnover.
	p0 := a.Points[0]
	if p0.TurnoverCount != 1 {
		t.Errorf("Points[0].TurnoverCount = %d, want 1", p0.TurnoverCount)
	}
	if p0.TotalUnits != 300 {
		t.Errorf("Points[0].TotalUnits = %d, want 300", p0.TotalUnits)
	}
	if p0.MarkupPerUnit != 50 {
		t.Errorf("Points[0].MarkupPerUnit = %d, want 50", p0.MarkupPerUnit)
	}
	// Points[1]: 5-turnover.
	p1 := a.Points[1]
	if p1.TurnoverCount != 5 {
		t.Errorf("Points[1].TurnoverCount = %d, want 5", p1.TurnoverCount)
	}
	if p1.TotalUnits != 1500 {
		t.Errorf("Points[1].TotalUnits = %d, want 1500", p1.TotalUnits)
	}
	if p1.MarkupPerUnit != 10 {
		t.Errorf("Points[1].MarkupPerUnit = %d, want 10", p1.MarkupPerUnit)
	}
}

// TestNewTurnoverEffectAnalysis_EmptyCounts verifies that an empty turnoverCounts
// slice returns ErrInvalidTurnoverCount.
func TestNewTurnoverEffectAnalysis_EmptyCounts(t *testing.T) {
	t.Parallel()

	_, err := NewTurnoverEffectAnalysis(100000, 1500, 300, []int64{})
	if !errors.Is(err, ErrInvalidTurnoverCount) {
		t.Errorf("empty counts: err = %v, want ErrInvalidTurnoverCount", err)
	}
}

// TestNewTurnoverEffectAnalysis_ZeroCountInSlice verifies that a count of 0 in
// the slice returns ErrInvalidTurnoverCount.
func TestNewTurnoverEffectAnalysis_ZeroCountInSlice(t *testing.T) {
	t.Parallel()

	_, err := NewTurnoverEffectAnalysis(100000, 1500, 300, []int64{1, 0, 3})
	if !errors.Is(err, ErrInvalidTurnoverCount) {
		t.Errorf("zero count in slice: err = %v, want ErrInvalidTurnoverCount", err)
	}
}

// --- NewMerchantPriceOfProduction ---

// TestNewMerchantPriceOfProduction_Valid covers costPrice=100000, averageProfit=15000
// → PriceOfProduction=115000.
func TestNewMerchantPriceOfProduction_Valid(t *testing.T) {
	t.Parallel()

	p, err := NewMerchantPriceOfProduction(100000, 15000)
	if err != nil {
		t.Fatalf("NewMerchantPriceOfProduction: %v", err)
	}
	if p.PriceOfProduction != 115000 {
		t.Errorf("PriceOfProduction = %d, want 115000", p.PriceOfProduction)
	}
	if p.CostPrice != 100000 {
		t.Errorf("CostPrice = %d, want 100000", p.CostPrice)
	}
	if p.AverageProfit != 15000 {
		t.Errorf("AverageProfit = %d, want 15000", p.AverageProfit)
	}
}

// TestNewMerchantPriceOfProduction_ZeroCostPrice verifies that costPrice=0 with
// averageProfit>0 is valid (zero cost is non-negative).
func TestNewMerchantPriceOfProduction_ZeroCostPrice(t *testing.T) {
	t.Parallel()

	p, err := NewMerchantPriceOfProduction(0, 15000)
	if err != nil {
		t.Fatalf("NewMerchantPriceOfProduction(costPrice=0): %v", err)
	}
	if p.PriceOfProduction != 15000 {
		t.Errorf("PriceOfProduction = %d, want 15000", p.PriceOfProduction)
	}
}

// TestNewMerchantPriceOfProduction_NegativeInputs verifies that negative costPrice
// or negative averageProfit returns ErrNonPositiveCapital.
func TestNewMerchantPriceOfProduction_NegativeInputs(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name             string
		cost, avgProfit  int64
	}{
		{"negative costPrice", -1, 15000},
		{"negative averageProfit", 100000, -1},
		{"both negative", -1, -1},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewMerchantPriceOfProduction(tc.cost, tc.avgProfit)
			if !errors.Is(err, ErrNonPositiveCapital) {
				t.Errorf("%s: err = %v, want ErrNonPositiveCapital", tc.name, err)
			}
		})
	}
}
