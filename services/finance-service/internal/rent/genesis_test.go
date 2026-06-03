// Package rent — unit tests for Vol. III Ch. 47 (Genesis of Capitalist Ground-Rent)
// domain types and pure functions.
//
// Fixtures drawn from Marx, Capital Vol. III Ch. 47:
//   - Historical rent stages 1–4: labour-rent, product-rent, money-rent,
//     capitalist ground-rent (Marx's own ordering).
//   - Small peasant proprietor (Normandy, 10th–13th c): conflates rent, profit,
//     and wages into one undifferentiated income.
//   - ComputeTotalIncome: (2+1+3)=6 basis points — the unit Marx uses when
//     decomposing the peasant's gross return.
package rent

import (
	"testing"
)

// ── RentFormStage.IsValid ─────────────────────────────────────────────────────

// TestRentFormStageIsValid: stages 1–4 are valid; 0 and 5 are out-of-range.
func TestRentFormStageIsValid(t *testing.T) {
	t.Parallel()

	valid := []RentFormStage{
		RentFormLabour,
		RentFormProduct,
		RentFormMoney,
		RentFormCapitalist,
	}
	for _, s := range valid {
		if !s.IsValid() {
			t.Errorf("RentFormStage(%d).IsValid() = false, want true", s)
		}
	}

	invalid := []RentFormStage{0, 5}
	for _, s := range invalid {
		if s.IsValid() {
			t.Errorf("RentFormStage(%d).IsValid() = true, want false", s)
		}
	}
}

// ── ComputeTotalIncome ────────────────────────────────────────────────────────

// TestComputeTotalIncome: peasant's undifferentiated income = rent+profit+wages.
func TestComputeTotalIncome(t *testing.T) {
	t.Parallel()

	cases := []struct {
		rent, profit, wages int64
		want                int64
	}{
		{2, 1, 3, 6},   // Normandy small-peasant, normalised to 6 bp
		{0, 0, 0, 0},   // degenerate: zero income
		{10, 20, 30, 60}, // proportional scaling
	}

	for _, tc := range cases {
		got := ComputeTotalIncome(tc.rent, tc.profit, tc.wages)
		if got != tc.want {
			t.Errorf("ComputeTotalIncome(%d,%d,%d) = %d, want %d",
				tc.rent, tc.profit, tc.wages, got, tc.want)
		}
	}
}

// ── Stage ordering & Validate guards ──────────────────────────────────────────

// Rent forms are historically ordered labour(1) < product(2) < money(3) <
// capitalist(4); the capitalist stage is terminal.
func TestRentFormStage_Ordering(t *testing.T) {
	t.Parallel()
	if !(RentFormLabour < RentFormProduct &&
		RentFormProduct < RentFormMoney &&
		RentFormMoney < RentFormCapitalist) {
		t.Error("rent form stages are not strictly ordered 1<2<3<4")
	}
	if RentFormCapitalist != 4 {
		t.Errorf("capitalist ground-rent stage = %d, want 4 (terminal)", RentFormCapitalist)
	}
}

// RentHistoricalForm.Validate accepts defined stages and rejects out-of-range ones.
func TestRentHistoricalForm_Validate(t *testing.T) {
	t.Parallel()
	if err := (RentHistoricalForm{Stage: RentFormMoney}).Validate(); err != nil {
		t.Errorf("money-rent form should validate, got %v", err)
	}
	if err := (RentHistoricalForm{Stage: 0}).Validate(); err != ErrInvalidRentFormStage {
		t.Errorf("stage 0 → %v, want ErrInvalidRentFormStage", err)
	}
}

// HistoricalRentTransition.Validate enforces a strictly progressive, well-formed
// transition: forward steps pass; stationary, regressive, and out-of-range fail.
func TestHistoricalRentTransition_Validate(t *testing.T) {
	t.Parallel()
	if err := (HistoricalRentTransition{FromStage: RentFormLabour, ToStage: RentFormMoney}).Validate(); err != nil {
		t.Errorf("progressive transition should validate, got %v", err)
	}
	if err := (HistoricalRentTransition{FromStage: RentFormMoney, ToStage: RentFormMoney}).Validate(); err != ErrRentTransitionNotProgressive {
		t.Errorf("stationary transition → %v, want ErrRentTransitionNotProgressive", err)
	}
	if err := (HistoricalRentTransition{FromStage: RentFormCapitalist, ToStage: RentFormLabour}).Validate(); err != ErrRentTransitionNotProgressive {
		t.Errorf("regressive transition → %v, want ErrRentTransitionNotProgressive", err)
	}
	if err := (HistoricalRentTransition{FromStage: 0, ToStage: RentFormMoney}).Validate(); err != ErrInvalidRentFormStage {
		t.Errorf("invalid from-stage → %v, want ErrInvalidRentFormStage", err)
	}
}

// ── ID constructors ───────────────────────────────────────────────────────────

// TestCh47NewIDsDistinct: each New*ID returns a 24-char hex string and two
// successive calls differ.
func TestCh47NewIDsDistinct(t *testing.T) {
	t.Parallel()

	// RentHistoricalFormID
	id1 := NewRentHistoricalFormID()
	id2 := NewRentHistoricalFormID()
	if len(string(id1)) != 24 {
		t.Errorf("NewRentHistoricalFormID: len = %d, want 24", len(string(id1)))
	}
	if id1 == id2 {
		t.Error("NewRentHistoricalFormID: two consecutive calls returned identical IDs")
	}

	// HistoricalRentTransitionID
	tid1 := NewHistoricalRentTransitionID()
	tid2 := NewHistoricalRentTransitionID()
	if len(string(tid1)) != 24 {
		t.Errorf("NewHistoricalRentTransitionID: len = %d, want 24", len(string(tid1)))
	}
	if tid1 == tid2 {
		t.Error("NewHistoricalRentTransitionID: two consecutive calls returned identical IDs")
	}

	// SmallPeasantProductionID
	pid1 := NewSmallPeasantProductionID()
	pid2 := NewSmallPeasantProductionID()
	if len(string(pid1)) != 24 {
		t.Errorf("NewSmallPeasantProductionID: len = %d, want 24", len(string(pid1)))
	}
	if pid1 == pid2 {
		t.Error("NewSmallPeasantProductionID: two consecutive calls returned identical IDs")
	}
}
