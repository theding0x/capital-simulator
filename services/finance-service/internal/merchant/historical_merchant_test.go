package merchant

import (
	"encoding/hex"
	"errors"
	"testing"
)

// --- HistoricalMerchantCapitalID ---

// TestNewHistoricalMerchantCapitalID_HexUnique asserts IDs are non-empty, 24
// hex chars (96-bit), unique, and that IsZero works correctly.
func TestNewHistoricalMerchantCapitalID_HexUnique(t *testing.T) {
	t.Parallel()

	seen := make(map[HistoricalMerchantCapitalID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewHistoricalMerchantCapitalID()
		if id.IsZero() {
			t.Fatal("NewHistoricalMerchantCapitalID returned the empty id")
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

// TestHistoricalMerchantCapitalID_IsZero asserts IsZero returns true for the
// empty string and false for a freshly generated id.
func TestHistoricalMerchantCapitalID_IsZero(t *testing.T) {
	t.Parallel()

	var empty HistoricalMerchantCapitalID
	if !empty.IsZero() {
		t.Error("empty HistoricalMerchantCapitalID.IsZero() = false, want true")
	}
	id := NewHistoricalMerchantCapitalID()
	if id.IsZero() {
		t.Error("freshly generated HistoricalMerchantCapitalID.IsZero() = true, want false")
	}
}

// --- DevelopmentStage.IsValid ---

// TestDevelopmentStage_IsValid asserts that 1–3 are valid and 0 / 4 are not.
func TestDevelopmentStage_IsValid(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		stage DevelopmentStage
		valid bool
	}{
		{0, false},
		{StagePreCapitalist, true},
		{StageTransition, true},
		{StageSubordinated, true},
		{4, false},
	} {
		tc := tc
		t.Run("", func(t *testing.T) {
			t.Parallel()
			if got := tc.stage.IsValid(); got != tc.valid {
				t.Errorf("DevelopmentStage(%d).IsValid() = %v, want %v", tc.stage, got, tc.valid)
			}
		})
	}
}

// --- NewHistoricalMerchantCapital ---

// TestNewHistoricalMerchantCapital_Venice covers the Venice seed fixture
// (Vol. III Ch. 20, seed 2001):
//
//	name="Venice carrying trade", stage=StagePreCapitalist,
//	profitSource="unequal_exchange", wageLabour=false, subordinationIndex=500
//
// All fields must round-trip; no ID or timestamp assigned.
func TestNewHistoricalMerchantCapital_Venice(t *testing.T) {
	t.Parallel()

	hm, err := NewHistoricalMerchantCapital(
		"Venice carrying trade",
		StagePreCapitalist,
		"unequal_exchange",
		false,
		500,
	)
	if err != nil {
		t.Fatalf("NewHistoricalMerchantCapital: %v", err)
	}
	if hm.Name != "Venice carrying trade" {
		t.Errorf("Name = %q, want %q", hm.Name, "Venice carrying trade")
	}
	if hm.Stage != StagePreCapitalist {
		t.Errorf("Stage = %d, want %d", hm.Stage, StagePreCapitalist)
	}
	if hm.ProfitSource != "unequal_exchange" {
		t.Errorf("ProfitSource = %q, want %q", hm.ProfitSource, "unequal_exchange")
	}
	if hm.WageLabour {
		t.Error("WageLabour = true, want false for pre-capitalist")
	}
	if hm.SubordinationIndex != 500 {
		t.Errorf("SubordinationIndex = %d, want 500", hm.SubordinationIndex)
	}
	// Constructor assigns no ID or timestamp — the store does that.
	if !hm.ID.IsZero() {
		t.Errorf("ID = %q, want empty (store assigns it)", hm.ID)
	}
	if !hm.CreatedAt.IsZero() {
		t.Errorf("CreatedAt = %v, want zero (store assigns it)", hm.CreatedAt)
	}
}

// TestNewHistoricalMerchantCapital_PreCapitalistWageLabour asserts that
// stage=StagePreCapitalist with wageLabour=true returns
// ErrWageLabourInPreCapitalist (invariant 1).
func TestNewHistoricalMerchantCapital_PreCapitalistWageLabour(t *testing.T) {
	t.Parallel()

	_, err := NewHistoricalMerchantCapital(
		"Hypothetical pre-cap",
		StagePreCapitalist,
		"unequal_exchange",
		true,  // forbidden
		500,
	)
	if !errors.Is(err, ErrWageLabourInPreCapitalist) {
		t.Errorf("err = %v, want ErrWageLabourInPreCapitalist", err)
	}
}

// TestNewHistoricalMerchantCapital_SubordinationBoundaries checks the
// [0, 10000] boundary conditions for SubordinationIndex.
func TestNewHistoricalMerchantCapital_SubordinationBoundaries(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		index   int64
		wantErr error
	}{
		{"min boundary 0", 0, nil},
		{"max boundary 10000", 10000, nil},
		{"below min -1", -1, ErrSubordinationOutOfRange},
		{"above max 10001", 10001, ErrSubordinationOutOfRange},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewHistoricalMerchantCapital(
				"test",
				StageTransition,
				"carrying_monopoly",
				false,
				tc.index,
			)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("%s: err = %v, want %v", tc.name, err, tc.wantErr)
			}
		})
	}
}

// TestNewHistoricalMerchantCapital_InvalidStage asserts that stage 0 and
// stage 4 return ErrInvalidDevelopmentStage.
func TestNewHistoricalMerchantCapital_InvalidStage(t *testing.T) {
	t.Parallel()

	for _, bad := range []DevelopmentStage{0, 4} {
		bad := bad
		t.Run("", func(t *testing.T) {
			t.Parallel()
			_, err := NewHistoricalMerchantCapital("test", bad, "unequal_exchange", false, 500)
			if !errors.Is(err, ErrInvalidDevelopmentStage) {
				t.Errorf("stage=%d: err = %v, want ErrInvalidDevelopmentStage", bad, err)
			}
		})
	}
}

// TestNewHistoricalMerchantCapital_SubordinatedWageLabour asserts that
// StageSubordinated with wageLabour=true is permitted (only StagePreCapitalist
// forbids wage-labour).
func TestNewHistoricalMerchantCapital_SubordinatedWageLabour(t *testing.T) {
	t.Parallel()

	hm, err := NewHistoricalMerchantCapital(
		"Subordination under industrial capital",
		StageSubordinated,
		"industrial_subordination",
		true, // allowed for StageSubordinated
		7000,
	)
	if err != nil {
		t.Fatalf("NewHistoricalMerchantCapital StageSubordinated+wageLabour=true: %v", err)
	}
	if !hm.WageLabour {
		t.Error("WageLabour = false, want true for subordinated stage")
	}
}

// TestNewHistoricalMerchantCapital_TransitionWageLabour asserts that
// StageTransition with wageLabour=true is also permitted.
func TestNewHistoricalMerchantCapital_TransitionWageLabour(t *testing.T) {
	t.Parallel()

	_, err := NewHistoricalMerchantCapital(
		"Dutch carrying trade",
		StageTransition,
		"carrying_monopoly",
		true,
		3000,
	)
	if err != nil {
		t.Fatalf("NewHistoricalMerchantCapital StageTransition+wageLabour=true: %v", err)
	}
}

// --- InverseDevelopmentRelation ---

// TestNewInverseDevelopmentRelation_Monotonic builds stage points from Marx's
// seed fixtures with strictly rising subordination indices and asserts
// Monotonic==true (invariant 2).
//
//	Venice (stage 1, index 500) → Dutch (stage 2, index 3000) → Subordinated (stage 3, index 7000)
func TestNewInverseDevelopmentRelation_Monotonic(t *testing.T) {
	t.Parallel()

	points := []StagePoint{
		{Stage: StagePreCapitalist, Index: 500},
		{Stage: StageTransition, Index: 3000},
		{Stage: StageSubordinated, Index: 7000},
	}
	rel, err := NewInverseDevelopmentRelation(points)
	if err != nil {
		t.Fatalf("NewInverseDevelopmentRelation: %v", err)
	}
	if !rel.Monotonic {
		t.Error("Monotonic = false, want true for strictly rising indices")
	}
	if len(rel.Stages) != 3 {
		t.Errorf("len(Stages) = %d, want 3", len(rel.Stages))
	}
	if len(rel.SubordinationIndices) != 3 {
		t.Errorf("len(SubordinationIndices) = %d, want 3", len(rel.SubordinationIndices))
	}
	// Sorted by stage ascending.
	if rel.Stages[0] != StagePreCapitalist || rel.Stages[1] != StageTransition || rel.Stages[2] != StageSubordinated {
		t.Errorf("stages order unexpected: %v", rel.Stages)
	}
	if rel.SubordinationIndices[0] != 500 || rel.SubordinationIndices[1] != 3000 || rel.SubordinationIndices[2] != 7000 {
		t.Errorf("indices order unexpected: %v", rel.SubordinationIndices)
	}
}

// TestNewInverseDevelopmentRelation_NonMonotonic builds a sequence where the
// subordination index does not strictly increase and asserts Monotonic==false.
func TestNewInverseDevelopmentRelation_NonMonotonic(t *testing.T) {
	t.Parallel()

	// Stage 1→3000, Stage 2→1000 (falls) — not monotonically increasing.
	points := []StagePoint{
		{Stage: StagePreCapitalist, Index: 3000},
		{Stage: StageTransition, Index: 1000},
	}
	rel, err := NewInverseDevelopmentRelation(points)
	if err != nil {
		t.Fatalf("NewInverseDevelopmentRelation: %v", err)
	}
	if rel.Monotonic {
		t.Error("Monotonic = true, want false for non-increasing sequence")
	}
}

// TestNewInverseDevelopmentRelation_EqualIndices asserts that equal adjacent
// indices also give Monotonic==false (strict increase required).
func TestNewInverseDevelopmentRelation_EqualIndices(t *testing.T) {
	t.Parallel()

	points := []StagePoint{
		{Stage: StagePreCapitalist, Index: 3000},
		{Stage: StageTransition, Index: 3000}, // equal, not strictly greater
	}
	rel, err := NewInverseDevelopmentRelation(points)
	if err != nil {
		t.Fatalf("NewInverseDevelopmentRelation: %v", err)
	}
	if rel.Monotonic {
		t.Error("Monotonic = true, want false for equal adjacent indices")
	}
}

// TestNewInverseDevelopmentRelation_InvalidStage asserts ErrInvalidDevelopmentStage
// is returned when a point carries stage 0.
func TestNewInverseDevelopmentRelation_InvalidStage(t *testing.T) {
	t.Parallel()

	points := []StagePoint{
		{Stage: 0, Index: 500},
	}
	_, err := NewInverseDevelopmentRelation(points)
	if !errors.Is(err, ErrInvalidDevelopmentStage) {
		t.Errorf("err = %v, want ErrInvalidDevelopmentStage", err)
	}
}

// TestNewInverseDevelopmentRelation_IndexOutOfRange asserts
// ErrSubordinationOutOfRange is returned for an index above 10000.
func TestNewInverseDevelopmentRelation_IndexOutOfRange(t *testing.T) {
	t.Parallel()

	points := []StagePoint{
		{Stage: StagePreCapitalist, Index: 10001},
	}
	_, err := NewInverseDevelopmentRelation(points)
	if !errors.Is(err, ErrSubordinationOutOfRange) {
		t.Errorf("err = %v, want ErrSubordinationOutOfRange", err)
	}
}

// --- HistoricalTradingSystem ---

// TestNewHistoricalTradingSystem_Venice builds the Venice trading system value
// object and checks fields round-trip.
func TestNewHistoricalTradingSystem_Venice(t *testing.T) {
	t.Parallel()

	ts, err := NewHistoricalTradingSystem(
		"Venice carrying trade",
		StagePreCapitalist,
		"unequal_exchange",
		false,
	)
	if err != nil {
		t.Fatalf("NewHistoricalTradingSystem: %v", err)
	}
	if ts.Name != "Venice carrying trade" {
		t.Errorf("Name = %q, want Venice carrying trade", ts.Name)
	}
	if ts.Stage != StagePreCapitalist {
		t.Errorf("Stage = %d, want %d", ts.Stage, StagePreCapitalist)
	}
	if ts.WageLabour {
		t.Error("WageLabour = true, want false")
	}
}

// TestNewHistoricalTradingSystem_InvalidStage returns ErrInvalidDevelopmentStage
// for stage 0.
func TestNewHistoricalTradingSystem_InvalidStage(t *testing.T) {
	t.Parallel()

	_, err := NewHistoricalTradingSystem("bad", 0, "unequal_exchange", false)
	if !errors.Is(err, ErrInvalidDevelopmentStage) {
		t.Errorf("err = %v, want ErrInvalidDevelopmentStage", err)
	}
}

// TestNewHistoricalTradingSystem_PreCapitalistWageLabour returns
// ErrWageLabourInPreCapitalist when stage=1 and wageLabour=true.
func TestNewHistoricalTradingSystem_PreCapitalistWageLabour(t *testing.T) {
	t.Parallel()

	_, err := NewHistoricalTradingSystem(
		"bad",
		StagePreCapitalist,
		"unequal_exchange",
		true,
	)
	if !errors.Is(err, ErrWageLabourInPreCapitalist) {
		t.Errorf("err = %v, want ErrWageLabourInPreCapitalist", err)
	}
}

// --- PreCapitalistForm ---

// TestNewPreCapitalistForm_ExploitsWageLabourAlwaysFalse asserts that
// ExploitsWageLabour is always false regardless of caller intent (structural
// enforcement of invariant 1).
func TestNewPreCapitalistForm_ExploitsWageLabourAlwaysFalse(t *testing.T) {
	t.Parallel()

	pcf, err := NewPreCapitalistForm("unequal_exchange")
	if err != nil {
		t.Fatalf("NewPreCapitalistForm: %v", err)
	}
	if pcf.ExploitsWageLabour {
		t.Error("ExploitsWageLabour = true, want always false for pre-capitalist form")
	}
	if pcf.ProfitSource != "unequal_exchange" {
		t.Errorf("ProfitSource = %q, want %q", pcf.ProfitSource, "unequal_exchange")
	}
}

// TestNewPreCapitalistForm_MonopolyTrade asserts monopoly_trade profit source
// also produces ExploitsWageLabour=false.
func TestNewPreCapitalistForm_MonopolyTrade(t *testing.T) {
	t.Parallel()

	pcf, err := NewPreCapitalistForm("monopoly_trade")
	if err != nil {
		t.Fatalf("NewPreCapitalistForm: %v", err)
	}
	if pcf.ExploitsWageLabour {
		t.Error("ExploitsWageLabour = true, want false")
	}
}
