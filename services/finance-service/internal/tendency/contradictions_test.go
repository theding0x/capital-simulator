package tendency

import (
	"encoding/hex"
	"errors"
	"testing"
)

// --- §I InternalContradiction / ContradictionKind ---

// TestContradictionKind_IsValid asserts all four named kinds are valid, an
// unknown string is invalid, and each const's string value is exactly the
// snake_case the spec names.
func TestContradictionKind_IsValid(t *testing.T) {
	t.Parallel()

	valid := map[ContradictionKind]string{
		KindFallingRateRisingMass:                  "falling_rate_rising_mass",
		KindOverproductionUnderConsumption:         "overproduction_underconsumption",
		KindOveraccumulationRelativeOverpopulation: "overaccumulation_relative_overpopulation",
		KindConcentrationCentralisation:            "concentration_centralisation",
	}
	for k, want := range valid {
		if !k.IsValid() {
			t.Errorf("%q.IsValid() = false, want true", k)
		}
		if string(k) != want {
			t.Errorf("kind string = %q, want %q", string(k), want)
		}
		if k.Description() == "" {
			t.Errorf("%q.Description() is empty, want a gloss", k)
		}
	}

	for _, bad := range []ContradictionKind{"", "usury", "Falling_Rate_Rising_Mass", "crisis"} {
		if bad.IsValid() {
			t.Errorf("%q.IsValid() = true, want false", bad)
		}
		if bad.Description() != "" {
			t.Errorf("%q.Description() = %q, want empty", bad, bad.Description())
		}
	}
}

// TestNewInternalContradiction_Coexistent covers §I: a contradiction built with
// KindFallingRateRisingMass is coexistent, and so are all four valid kinds.
func TestNewInternalContradiction_Coexistent(t *testing.T) {
	t.Parallel()

	c, err := NewInternalContradiction(KindFallingRateRisingMass, "Marx §I")
	if err != nil {
		t.Fatalf("NewInternalContradiction(falling_rate_rising_mass): %v", err)
	}
	if !c.IsCoexistent {
		t.Error("IsCoexistent = false, want true for falling_rate_rising_mass (§I)")
	}
	if c.Kind != KindFallingRateRisingMass {
		t.Errorf("Kind = %q, want %q", c.Kind, KindFallingRateRisingMass)
	}
	if c.Note != "Marx §I" {
		t.Errorf("Note = %q, want %q", c.Note, "Marx §I")
	}
	// Constructor assigns no ID or timestamp — the store does that.
	if !c.ID.IsZero() {
		t.Errorf("ID = %q, want empty (store assigns it)", c.ID)
	}
	if !c.CreatedAt.IsZero() {
		t.Errorf("CreatedAt = %v, want zero (store assigns it)", c.CreatedAt)
	}

	for _, k := range []ContradictionKind{
		KindFallingRateRisingMass,
		KindOverproductionUnderConsumption,
		KindOveraccumulationRelativeOverpopulation,
		KindConcentrationCentralisation,
	} {
		got, err := NewInternalContradiction(k, "")
		if err != nil {
			t.Fatalf("NewInternalContradiction(%q): %v", k, err)
		}
		if !got.IsCoexistent {
			t.Errorf("%q: IsCoexistent = false, want true (all four kinds coexist)", k)
		}
	}
}

// TestNewInternalContradiction_InvalidKind returns ErrInvalidContradictionKind
// for an unknown kind.
func TestNewInternalContradiction_InvalidKind(t *testing.T) {
	t.Parallel()
	if _, err := NewInternalContradiction(ContradictionKind("usury"), ""); !errors.Is(err, ErrInvalidContradictionKind) {
		t.Errorf("unknown kind: err = %v, want ErrInvalidContradictionKind", err)
	}
	if _, err := NewInternalContradiction(ContradictionKind(""), ""); !errors.Is(err, ErrInvalidContradictionKind) {
		t.Errorf("empty kind: err = %v, want ErrInvalidContradictionKind", err)
	}
}

// --- §II Overproduction ---

// TestOverproduction_UnrealisedValue covers the §II fixture: excess commodity
// value 50000 against realisable value 40000 leaves 10000 unrealised, plus a
// zero/negative boundary.
func TestOverproduction_UnrealisedValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		excess  int64
		realisn int64
		want    int64
	}{
		{"50000/40000", 50000, 40000, 10000},
		{"all realised", 40000, 40000, 0},
		{"over-realised (negative)", 40000, 50000, -10000},
		{"zero", 0, 0, 0},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			o := Overproduction{ExcessCommodityValue: tc.excess, RealisableValue: tc.realisn}
			if got := o.UnrealisedValue(); got != tc.want {
				t.Errorf("UnrealisedValue() = %d, want %d", got, tc.want)
			}
		})
	}
}

// TestOverproduction_Validate covers the happy path and the negative-magnitude
// boundary, and confirms the constructor surfaces the same.
func TestOverproduction_Validate(t *testing.T) {
	t.Parallel()

	ok := Overproduction{ExcessCommodityValue: 50000, RealisableValue: 40000}
	if err := ok.Validate(); err != nil {
		t.Errorf("valid overproduction: err = %v, want nil", err)
	}

	for _, neg := range []Overproduction{
		{ExcessCommodityValue: -1, RealisableValue: 40000},
		{ExcessCommodityValue: 50000, RealisableValue: -1},
	} {
		if err := neg.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
			t.Errorf("%+v: err = %v, want ErrNegativeMagnitude", neg, err)
		}
	}

	o, err := ComputeOverproduction(50000, 40000)
	if err != nil {
		t.Fatalf("ComputeOverproduction: %v", err)
	}
	if o.UnrealisedValue() != 10000 {
		t.Errorf("UnrealisedValue() = %d, want 10000", o.UnrealisedValue())
	}
	if _, err := ComputeOverproduction(-1, 40000); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("ComputeOverproduction(-1,...): err = %v, want ErrNegativeMagnitude", err)
	}
}

// --- §III Overaccumulation ---

// TestOveraccumulation_AbsoluteAndShortfall covers §III: IsAbsolute is true when
// the actual rate on the excess capital is ≤ 0, and ShortfallBP is the signed
// gap between the average and actual rates.
func TestOveraccumulation_AbsoluteAndShortfall(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		excess       int64
		average      int64
		actual       int64
		wantAbsolute bool
		wantShortBP  int64
	}{
		{"relative (actual 1000 < average 2000)", 100000, 2000, 1000, false, 1000},
		{"absolute (actual 0)", 100000, 2000, 0, true, 2000},
		{"at par (actual == average)", 100000, 2000, 2000, false, 0},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			o := Overaccumulation{ExcessCapital: tc.excess, AverageProfitRate: tc.average, ActualProfitRate: tc.actual}
			if got := o.IsAbsolute(); got != tc.wantAbsolute {
				t.Errorf("IsAbsolute() = %v, want %v", got, tc.wantAbsolute)
			}
			if got := o.ShortfallBP(); got != tc.wantShortBP {
				t.Errorf("ShortfallBP() = %d, want %d", got, tc.wantShortBP)
			}
		})
	}
}

// TestOveraccumulation_Validate covers the happy path and the negative-magnitude
// boundary, and confirms the constructor surfaces the same.
func TestOveraccumulation_Validate(t *testing.T) {
	t.Parallel()

	ok := Overaccumulation{ExcessCapital: 100000, AverageProfitRate: 2000, ActualProfitRate: 1000}
	if err := ok.Validate(); err != nil {
		t.Errorf("valid overaccumulation: err = %v, want nil", err)
	}

	for _, neg := range []Overaccumulation{
		{ExcessCapital: -1, AverageProfitRate: 2000, ActualProfitRate: 1000},
		{ExcessCapital: 100000, AverageProfitRate: -1, ActualProfitRate: 1000},
		{ExcessCapital: 100000, AverageProfitRate: 2000, ActualProfitRate: -1},
	} {
		if err := neg.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
			t.Errorf("%+v: err = %v, want ErrNegativeMagnitude", neg, err)
		}
	}

	o, err := ComputeOveraccumulation(100000, 2000, 1000)
	if err != nil {
		t.Fatalf("ComputeOveraccumulation: %v", err)
	}
	if o.ShortfallBP() != 1000 || o.IsAbsolute() {
		t.Errorf("computed = shortfall %d absolute %v, want 1000/false", o.ShortfallBP(), o.IsAbsolute())
	}
	if _, err := ComputeOveraccumulation(100000, 2000, -1); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("ComputeOveraccumulation(neg): err = %v, want ErrNegativeMagnitude", err)
	}
}

// --- §IV CapitalConcentration ---

// TestCapitalConcentration_FinalCAndPath covers §IV: 1,000,000 accumulating 20%
// of itself each period for 3 periods reaches 1,728,000, with the per-period
// path [1000000, 1200000, 1440000, 1728000]. FinalC must exceed InitialC.
func TestCapitalConcentration_FinalCAndPath(t *testing.T) {
	t.Parallel()
	c := CapitalConcentration{InitialC: 1000000, SurplusValueRate: 2000, Periods: 3}

	if got := c.FinalC(); got != 1728000 {
		t.Errorf("FinalC() = %d, want 1728000", got)
	}
	if c.FinalC() <= c.InitialC {
		t.Errorf("FinalC %d must exceed InitialC %d", c.FinalC(), c.InitialC)
	}

	wantPath := []int64{1000000, 1200000, 1440000, 1728000}
	path := c.Path()
	if len(path) != len(wantPath) {
		t.Fatalf("Path() len = %d, want %d", len(path), len(wantPath))
	}
	for i, w := range wantPath {
		if path[i] != w {
			t.Errorf("Path()[%d] = %d, want %d", i, path[i], w)
		}
	}
	// Path's last element equals FinalC.
	if path[len(path)-1] != c.FinalC() {
		t.Errorf("Path() tail = %d, want FinalC %d", path[len(path)-1], c.FinalC())
	}
}

// TestCapitalConcentration_Validate covers the happy path, the negative-magnitude
// boundary, and the Periods<1 invariant via both Validate and the constructor.
func TestCapitalConcentration_Validate(t *testing.T) {
	t.Parallel()

	ok := CapitalConcentration{InitialC: 1000000, SurplusValueRate: 2000, Periods: 3}
	if err := ok.Validate(); err != nil {
		t.Errorf("valid concentration: err = %v, want nil", err)
	}

	for _, neg := range []CapitalConcentration{
		{InitialC: -1, SurplusValueRate: 2000, Periods: 3},
		{InitialC: 1000000, SurplusValueRate: -1, Periods: 3},
		{InitialC: 1000000, SurplusValueRate: 2000, Periods: -1},
	} {
		if err := neg.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
			t.Errorf("%+v: err = %v, want ErrNegativeMagnitude", neg, err)
		}
	}

	zeroPeriods := CapitalConcentration{InitialC: 1000000, SurplusValueRate: 2000, Periods: 0}
	if err := zeroPeriods.Validate(); !errors.Is(err, ErrInvalidPeriods) {
		t.Errorf("zero periods: err = %v, want ErrInvalidPeriods", err)
	}

	if _, err := ComputeCapitalConcentration(1000000, 2000, 0); !errors.Is(err, ErrInvalidPeriods) {
		t.Errorf("ComputeCapitalConcentration(periods=0): err = %v, want ErrInvalidPeriods", err)
	}
	c, err := ComputeCapitalConcentration(1000000, 2000, 3)
	if err != nil {
		t.Fatalf("ComputeCapitalConcentration: %v", err)
	}
	if c.FinalC() != 1728000 {
		t.Errorf("FinalC() = %d, want 1728000", c.FinalC())
	}
}

// --- §IV CapitalCentralisation ---

// TestCapitalCentralisation_ResultingCapital covers §IV: absorbing 10 capitals
// of 100,000 each yields a centralised capital of 1,000,000.
func TestCapitalCentralisation_ResultingCapital(t *testing.T) {
	t.Parallel()
	c := CapitalCentralisation{AbsorbedCapitals: 10, AbsorbedSize: 100000}
	if got := c.ResultingCapital(); got != 1000000 {
		t.Errorf("ResultingCapital() = %d, want 1000000", got)
	}

	// Boundary: zero absorbed capitals → zero result.
	zero := CapitalCentralisation{AbsorbedCapitals: 0, AbsorbedSize: 100000}
	if got := zero.ResultingCapital(); got != 0 {
		t.Errorf("ResultingCapital() = %d, want 0 for zero absorbed", got)
	}
}

// TestCapitalCentralisation_Validate covers the happy path and the
// negative-magnitude boundary, and confirms the constructor surfaces the same.
func TestCapitalCentralisation_Validate(t *testing.T) {
	t.Parallel()

	ok := CapitalCentralisation{AbsorbedCapitals: 10, AbsorbedSize: 100000}
	if err := ok.Validate(); err != nil {
		t.Errorf("valid centralisation: err = %v, want nil", err)
	}

	for _, neg := range []CapitalCentralisation{
		{AbsorbedCapitals: -1, AbsorbedSize: 100000},
		{AbsorbedCapitals: 10, AbsorbedSize: -1},
	} {
		if err := neg.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
			t.Errorf("%+v: err = %v, want ErrNegativeMagnitude", neg, err)
		}
	}

	c, err := ComputeCapitalCentralisation(10, 100000)
	if err != nil {
		t.Fatalf("ComputeCapitalCentralisation: %v", err)
	}
	if c.ResultingCapital() != 1000000 {
		t.Errorf("ResultingCapital() = %d, want 1000000", c.ResultingCapital())
	}
	if _, err := ComputeCapitalCentralisation(-1, 100000); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("ComputeCapitalCentralisation(neg): err = %v, want ErrNegativeMagnitude", err)
	}
}

// --- Crisis ---

// TestNewCrisis_RestoredRate covers the §crisis fixture: a 30% writedown of
// constant capital on a pre-crisis rate of 1000 bp restores the rate to 1429 bp
// (round-half-up of 1000·10000/7000), which exceeds the pre-crisis rate.
func TestNewCrisis_RestoredRate(t *testing.T) {
	t.Parallel()
	c, err := NewCrisis(30, 1000)
	if err != nil {
		t.Fatalf("NewCrisis(30, 1000): %v", err)
	}
	// (1000*10000 + 7000/2)/7000 = (10_000_000 + 3500)/7000 = 1429.
	if c.PostCrisisProfitRate != 1429 {
		t.Errorf("PostCrisisProfitRate = %d, want 1429 (round-half-up of 1000·10000/7000)", c.PostCrisisProfitRate)
	}
	if c.PostCrisisProfitRate <= c.PreCrisisProfitRate {
		t.Errorf("Post %d must exceed Pre %d", c.PostCrisisProfitRate, c.PreCrisisProfitRate)
	}
	if c.ConstantCapitalWritedown != 30 || c.PreCrisisProfitRate != 1000 {
		t.Errorf("inputs = writedown %d pre %d, want 30/1000", c.ConstantCapitalWritedown, c.PreCrisisProfitRate)
	}
	// Constructor assigns no ID or timestamp.
	if !c.ID.IsZero() {
		t.Errorf("ID = %q, want empty (store assigns it)", c.ID)
	}
	if !c.CreatedAt.IsZero() {
		t.Errorf("CreatedAt = %v, want zero (store assigns it)", c.CreatedAt)
	}
}

// TestNewCrisis_OtherFixtures cross-checks the two other seeded crises so the
// round-half-up restored rate is pinned at more than one point.
func TestNewCrisis_OtherFixtures(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		writedown int64
		pre       int64
		wantPost  int64
	}{
		{"writedown=20 pre=800", 20, 800, 1000},  // 8_000_000/8000 = 1000
		{"writedown=50 pre=600", 50, 600, 1200},  // 6_000_000/5000 = 1200
		{"writedown=30 pre=1000", 30, 1000, 1429}, // round-half-up
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			c, err := NewCrisis(tc.writedown, tc.pre)
			if err != nil {
				t.Fatalf("NewCrisis(%d, %d): %v", tc.writedown, tc.pre, err)
			}
			if c.PostCrisisProfitRate != tc.wantPost {
				t.Errorf("PostCrisisProfitRate = %d, want %d", c.PostCrisisProfitRate, tc.wantPost)
			}
			if c.PostCrisisProfitRate <= c.PreCrisisProfitRate {
				t.Errorf("Post %d must exceed Pre %d", c.PostCrisisProfitRate, c.PreCrisisProfitRate)
			}
		})
	}
}

// TestNewCrisis_InvalidWritedown rejects a writedown at or outside the (0,100)
// bounds with ErrInvalidWritedown, and a negative pre-crisis rate with
// ErrNegativeMagnitude.
func TestNewCrisis_InvalidWritedown(t *testing.T) {
	t.Parallel()

	for _, w := range []int64{0, 100, -1, 101, 200} {
		if _, err := NewCrisis(w, 1000); !errors.Is(err, ErrInvalidWritedown) {
			t.Errorf("writedown %d: err = %v, want ErrInvalidWritedown", w, err)
		}
	}

	// Boundary: writedown 1 and 99 are valid (inclusive interior).
	if _, err := NewCrisis(1, 1000); err != nil {
		t.Errorf("writedown 1: err = %v, want nil", err)
	}
	if _, err := NewCrisis(99, 1000); err != nil {
		t.Errorf("writedown 99: err = %v, want nil", err)
	}

	if _, err := NewCrisis(30, -1); !errors.Is(err, ErrNegativeMagnitude) {
		t.Errorf("negative pre-crisis rate: err = %v, want ErrNegativeMagnitude", err)
	}
}

// --- ID constructors ---

// TestNewCrisisID_HexUnique asserts crisis IDs are non-empty, valid hex, 24 chars
// (96 bits), and unique across a batch.
func TestNewCrisisID_HexUnique(t *testing.T) {
	t.Parallel()
	seen := make(map[CrisisID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewCrisisID()
		if id.IsZero() {
			t.Fatal("NewCrisisID returned the empty id")
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

// TestNewInternalContradictionID_HexUnique mirrors the crisis ID test for the
// contradiction ID constructor.
func TestNewInternalContradictionID_HexUnique(t *testing.T) {
	t.Parallel()
	seen := make(map[InternalContradictionID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewInternalContradictionID()
		if id.IsZero() {
			t.Fatal("NewInternalContradictionID returned the empty id")
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
