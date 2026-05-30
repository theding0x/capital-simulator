package credit

import (
	"encoding/hex"
	"errors"
	"testing"
)

// --- RateOfInterestID ---

// TestNewRateOfInterestID_HexUnique asserts IDs are non-empty, 24 hex chars
// (96-bit), unique, and that IsZero works correctly.
func TestNewRateOfInterestID_HexUnique(t *testing.T) {
	t.Parallel()

	seen := make(map[RateOfInterestID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewRateOfInterestID()
		if id.IsZero() {
			t.Fatal("NewRateOfInterestID returned the empty id")
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

// TestRateOfInterestID_IsZero asserts IsZero returns true for the empty string
// and false for a freshly generated id.
func TestRateOfInterestID_IsZero(t *testing.T) {
	t.Parallel()

	var empty RateOfInterestID
	if !empty.IsZero() {
		t.Error("empty RateOfInterestID.IsZero() = false, want true")
	}
	id := NewRateOfInterestID()
	if id.IsZero() {
		t.Error("freshly generated RateOfInterestID.IsZero() = true, want false")
	}
}

// --- IndustrialCyclePhase ---

// TestIndustrialCyclePhase_IsValid asserts all five constants are valid and
// values outside 1..5 are not.
func TestIndustrialCyclePhase_IsValid(t *testing.T) {
	t.Parallel()

	valid := []IndustrialCyclePhase{
		PhaseInactivity, PhaseRevival, PhaseProsperity, PhaseOverproduction, PhaseCrisis,
	}
	for _, p := range valid {
		p := p
		t.Run("", func(t *testing.T) {
			t.Parallel()
			if !p.IsValid() {
				t.Errorf("phase %d IsValid() = false, want true", p)
			}
		})
	}

	invalid := []IndustrialCyclePhase{0, 6, -1, 100}
	for _, p := range invalid {
		p := p
		t.Run("", func(t *testing.T) {
			t.Parallel()
			if p.IsValid() {
				t.Errorf("phase %d IsValid() = true, want false", p)
			}
		})
	}
}

// --- NewRateOfInterest ---

// TestNewRateOfInterest_Prosperity1843 covers the Marx fixture for Sept. 1843:
// interest fell to 150 bp (1.5%) during prosperity; avg profit 2000 bp (20%).
func TestNewRateOfInterest_Prosperity1843(t *testing.T) {
	t.Parallel()

	r, err := NewRateOfInterest(150, 2000, PhaseProsperity, "1843")
	if err != nil {
		t.Fatalf("NewRateOfInterest: %v", err)
	}
	if r.RateBP != 150 {
		t.Errorf("RateBP = %d, want 150", r.RateBP)
	}
	if r.AverageProfitRateBP != 2000 {
		t.Errorf("AverageProfitRateBP = %d, want 2000", r.AverageProfitRateBP)
	}
	if r.CyclePhase != PhaseProsperity {
		t.Errorf("CyclePhase = %d, want PhaseProsperity (%d)", r.CyclePhase, PhaseProsperity)
	}
	if r.Period != "1843" {
		t.Errorf("Period = %q, want 1843", r.Period)
	}
	// Constructor assigns no ID or timestamp — the store does that.
	if !r.ID.IsZero() {
		t.Errorf("ID = %q, want empty (store assigns it)", r.ID)
	}
	if !r.CreatedAt.IsZero() {
		t.Errorf("CreatedAt = %v, want zero (store assigns it)", r.CreatedAt)
	}
}

// TestNewRateOfInterest_Crisis1847 covers the Marx fixture for the 1847 crisis:
// interest rose to 800 bp (8%+); avg profit 2000 bp (20%).
func TestNewRateOfInterest_Crisis1847(t *testing.T) {
	t.Parallel()

	r, err := NewRateOfInterest(800, 2000, PhaseCrisis, "1847 crisis")
	if err != nil {
		t.Fatalf("NewRateOfInterest: %v", err)
	}
	if r.RateBP != 800 {
		t.Errorf("RateBP = %d, want 800", r.RateBP)
	}
	if r.CyclePhase != PhaseCrisis {
		t.Errorf("CyclePhase = %d, want PhaseCrisis (%d)", r.CyclePhase, PhaseCrisis)
	}
}

// TestNewRateOfInterest_ValidationErrors checks every sentinel error path.
func TestNewRateOfInterest_ValidationErrors(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		rateBP  int64
		avgBP   int64
		phase   IndustrialCyclePhase
		wantErr error
	}{
		{"rateBP<0", -1, 2000, PhaseProsperity, ErrNegativeRate},
		{"avgBP=0", 150, 0, PhaseProsperity, ErrNonPositiveAverageProfit},
		{"avgBP<0", 150, -1, PhaseProsperity, ErrNonPositiveAverageProfit},
		{"rateBP==avgBP", 2000, 2000, PhaseProsperity, ErrRateExceedsAverageProfit},
		{"rateBP>avgBP", 2001, 2000, PhaseProsperity, ErrRateExceedsAverageProfit},
		{"invalidPhase=0", 150, 2000, IndustrialCyclePhase(0), ErrInvalidCyclePhase},
		{"invalidPhase=6", 150, 2000, IndustrialCyclePhase(6), ErrInvalidCyclePhase},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewRateOfInterest(tc.rateBP, tc.avgBP, tc.phase, "test")
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
			}
		})
	}
}

// TestNewRateOfInterest_ZeroRateBP asserts rateBP=0 is valid (interest can
// approach but not quite reach zero in Marx's analysis).
func TestNewRateOfInterest_ZeroRateBP(t *testing.T) {
	t.Parallel()

	r, err := NewRateOfInterest(0, 2000, PhaseInactivity, "inactivity")
	if err != nil {
		t.Fatalf("NewRateOfInterest(0, 2000, ...): unexpected error: %v", err)
	}
	if r.RateBP != 0 {
		t.Errorf("RateBP = %d, want 0", r.RateBP)
	}
}

// --- NewInterestRateAnalysis ---

// TestNewInterestRateAnalysis_MarxQuarterOfProfit covers the canonical Ch. 22
// fixture: avg 2000 bp, interest = ¼ of 2000 = 500 bp → enterprise 1500 bp.
func TestNewInterestRateAnalysis_MarxQuarterOfProfit(t *testing.T) {
	t.Parallel()

	a, err := NewInterestRateAnalysis(2000, 500)
	if err != nil {
		t.Fatalf("NewInterestRateAnalysis: %v", err)
	}
	if a.AverageProfitRateBP != 2000 {
		t.Errorf("AverageProfitRateBP = %d, want 2000", a.AverageProfitRateBP)
	}
	if a.InterestRateBP != 500 {
		t.Errorf("InterestRateBP = %d, want 500", a.InterestRateBP)
	}
	// Invariant: ProfitOfEnterpriseBP == AverageProfitRateBP - InterestRateBP
	if a.ProfitOfEnterpriseBP != 1500 {
		t.Errorf("ProfitOfEnterpriseBP = %d, want 1500", a.ProfitOfEnterpriseBP)
	}
	if a.ProfitOfEnterpriseBP != a.AverageProfitRateBP-a.InterestRateBP {
		t.Errorf("invariant violated: %d != %d - %d", a.ProfitOfEnterpriseBP, a.AverageProfitRateBP, a.InterestRateBP)
	}
}

// TestNewInterestRateAnalysis_OneFifthTable covers Marx's 1/5 table from Ch. 22.
func TestNewInterestRateAnalysis_OneFifthTable(t *testing.T) {
	t.Parallel()

	cases := []struct {
		avgBP int64
		intBP int64 // 1/5 of avg
		entBP int64 // 4/5 of avg
	}{
		{1000, 200, 800},  // profit 10%, interest 2%, enterprise 8%
		{2000, 400, 1600}, // profit 20%, interest 4%, enterprise 16%
		{2500, 500, 2000}, // profit 25%, interest 5%, enterprise 20%
		{3000, 600, 2400}, // profit 30%, interest 6%, enterprise 24%
		{3500, 700, 2800}, // profit 35%, interest 7%, enterprise 28%
	}
	for _, tc := range cases {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			a, err := NewInterestRateAnalysis(tc.avgBP, tc.intBP)
			if err != nil {
				t.Fatalf("avg=%d int=%d: %v", tc.avgBP, tc.intBP, err)
			}
			if a.ProfitOfEnterpriseBP != tc.entBP {
				t.Errorf("avg=%d int=%d: ProfitOfEnterpriseBP=%d, want %d",
					tc.avgBP, tc.intBP, a.ProfitOfEnterpriseBP, tc.entBP)
			}
		})
	}
}

// TestNewInterestRateAnalysis_ValidationErrors checks every sentinel error path.
func TestNewInterestRateAnalysis_ValidationErrors(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		avgBP   int64
		intBP   int64
		wantErr error
	}{
		{"intBP<0", 2000, -1, ErrNegativeRate},
		{"avgBP=0", 0, 500, ErrNonPositiveAverageProfit},
		{"avgBP<0", -1, 500, ErrNonPositiveAverageProfit},
		{"intBP==avgBP", 2000, 2000, ErrRateExceedsAverageProfit},
		{"intBP>avgBP", 2000, 2001, ErrRateExceedsAverageProfit},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewInterestRateAnalysis(tc.avgBP, tc.intBP)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
			}
		})
	}
}

// TestNewInterestRateAnalysis_ZeroInterest asserts interest=0 is valid.
func TestNewInterestRateAnalysis_ZeroInterest(t *testing.T) {
	t.Parallel()

	a, err := NewInterestRateAnalysis(2000, 0)
	if err != nil {
		t.Fatalf("NewInterestRateAnalysis(2000, 0): unexpected error: %v", err)
	}
	if a.ProfitOfEnterpriseBP != 2000 {
		t.Errorf("ProfitOfEnterpriseBP = %d, want 2000", a.ProfitOfEnterpriseBP)
	}
}

// --- NewInterestRateBound ---

// TestNewInterestRateBound_Valid asserts a valid bound is accepted.
func TestNewInterestRateBound_Valid(t *testing.T) {
	t.Parallel()

	b, err := NewInterestRateBound(2000, 0)
	if err != nil {
		t.Fatalf("NewInterestRateBound: %v", err)
	}
	if b.MaxBP != 2000 {
		t.Errorf("MaxBP = %d, want 2000", b.MaxBP)
	}
	if b.MinBP != 0 {
		t.Errorf("MinBP = %d, want 0", b.MinBP)
	}
}

// TestNewInterestRateBound_ValidationErrors checks sentinel error paths.
func TestNewInterestRateBound_ValidationErrors(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		maxBP   int64
		minBP   int64
		wantErr error
	}{
		{"maxBP=0", 0, 0, ErrNonPositiveAverageProfit},
		{"maxBP<0", -1, 0, ErrNonPositiveAverageProfit},
		{"minBP<0", 2000, -1, ErrNegativeRate},
		{"minBP>maxBP", 2000, 2001, ErrRateExceedsAverageProfit},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewInterestRateBound(tc.maxBP, tc.minBP)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
			}
		})
	}
}

// --- NewMarketInterestRate ---

// TestNewMarketInterestRate_Valid asserts a non-negative rate is accepted.
func TestNewMarketInterestRate_Valid(t *testing.T) {
	t.Parallel()

	m, err := NewMarketInterestRate(150)
	if err != nil {
		t.Fatalf("NewMarketInterestRate: %v", err)
	}
	if m.RateBP != 150 {
		t.Errorf("RateBP = %d, want 150", m.RateBP)
	}
}

// TestNewMarketInterestRate_ZeroValid asserts rate=0 is accepted.
func TestNewMarketInterestRate_ZeroValid(t *testing.T) {
	t.Parallel()

	_, err := NewMarketInterestRate(0)
	if err != nil {
		t.Fatalf("NewMarketInterestRate(0): unexpected error: %v", err)
	}
}

// TestNewMarketInterestRate_Negative asserts rate<0 returns ErrNegativeRate.
func TestNewMarketInterestRate_Negative(t *testing.T) {
	t.Parallel()

	_, err := NewMarketInterestRate(-1)
	if !errors.Is(err, ErrNegativeRate) {
		t.Errorf("err = %v, want ErrNegativeRate", err)
	}
}
