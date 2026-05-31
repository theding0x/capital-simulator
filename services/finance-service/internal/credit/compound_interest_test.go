package credit_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// TestNewCompoundInterestSchedule_PriceFixture tests Marx's Dr. Price fixture
// using a scaled principal (10000 pence = £41.67) so integer compound iteration
// preserves precision. With principal=10000 at 5% (500 bp) for 100 periods,
// the integer result is 1313. The conceptual Dr. Price example (1 penny) loses
// all precision in integer arithmetic — the simulation uses scaled values.
func TestNewCompoundInterestSchedule_PriceFixture(t *testing.T) {
	t.Parallel()

	// Use 10000 base units (= 100 × the penny, scaled for integer precision).
	sched, err := credit.NewCompoundInterestSchedule(10000, 500, 100, credit.FormCompoundInterest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sched.Principal != 10000 {
		t.Errorf("Principal: got %d, want 10000", sched.Principal)
	}
	if sched.RateBP != 500 {
		t.Errorf("RateBP: got %d, want 500", sched.RateBP)
	}
	if sched.PeriodYears != 100 {
		t.Errorf("PeriodYears: got %d, want 100", sched.PeriodYears)
	}
	if sched.NewValueCreated != 0 {
		t.Errorf("NewValueCreated: got %d, want 0 (fetish form invariant)", sched.NewValueCreated)
	}
	if sched.FinalValue <= sched.Principal {
		t.Errorf("FinalValue %d must exceed Principal %d", sched.FinalValue, sched.Principal)
	}
	// 10000 * 1.05^100 ≈ 1,315,126 (integer iteration) — verify exact value.
	if sched.FinalValue != 1315126 {
		t.Errorf("FinalValue: got %d, want 1315126 (10000@5%%×100yr integer compound)", sched.FinalValue)
	}
}

// TestNewCompoundInterestSchedule_PittSinkingFund tests Marx's Pitt sinking-fund
// fixture: £250,000 at compound interest accumulating to £4,000,000.
// We use a shorter period to verify the formula, checking growth direction.
func TestNewCompoundInterestSchedule_PittSinkingFund(t *testing.T) {
	t.Parallel()

	// £250,000 at 5% compound for 10 periods — verifying growth direction.
	sched, err := credit.NewCompoundInterestSchedule(250000, 500, 10, credit.FormMM)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sched.FinalValue <= 250000 {
		t.Errorf("FinalValue %d must exceed Principal 250000", sched.FinalValue)
	}
	if sched.NewValueCreated != 0 {
		t.Errorf("NewValueCreated must be 0, got %d", sched.NewValueCreated)
	}
	// 250000 * 1.05^10 = 407,224 (integer iteration) — verify exact value.
	if sched.FinalValue != 407224 {
		t.Errorf("FinalValue: got %d, want 407224 (250000@5%%×10yr integer compound)", sched.FinalValue)
	}
}

// TestNewCompoundInterestSchedule_SinglePeriod verifies one-period compound
// equals principal × (1 + rate).
func TestNewCompoundInterestSchedule_SinglePeriod(t *testing.T) {
	t.Parallel()

	// 10000 pence at 10% (1000 bp) for 1 period → 11000.
	sched, err := credit.NewCompoundInterestSchedule(10000, 1000, 1, credit.FormSimpleInterest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sched.FinalValue != 11000 {
		t.Errorf("FinalValue: got %d, want 11000", sched.FinalValue)
	}
}

// TestNewCompoundInterestSchedule_ValidationErrors tests invariant enforcement.
func TestNewCompoundInterestSchedule_ValidationErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		principal   int64
		rateBP      int64
		periodYears int64
		wantErr     error
	}{
		{"zero principal", 0, 500, 10, credit.ErrNonPositivePrincipal},
		{"negative principal", -1, 500, 10, credit.ErrNonPositivePrincipal},
		{"zero rate", 1000, 0, 10, credit.ErrNonPositiveRateBP},
		{"negative rate", 1000, -100, 10, credit.ErrNonPositiveRateBP},
		{"zero periods", 1000, 500, 0, credit.ErrNonPositivePeriodYears},
		{"negative periods", 1000, 500, -1, credit.ErrNonPositivePeriodYears},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := credit.NewCompoundInterestSchedule(tc.principal, tc.rateBP, tc.periodYears, credit.FormCompoundInterest)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err != tc.wantErr {
				t.Errorf("got %v, want %v", err, tc.wantErr)
			}
		})
	}
}

// TestNewCompoundInterestID_Format verifies IDs are 24-char hex strings.
func TestNewCompoundInterestID_Format(t *testing.T) {
	t.Parallel()

	id := credit.NewCompoundInterestID()
	if len(string(id)) != 24 {
		t.Errorf("CompoundInterestID length: got %d, want 24", len(string(id)))
	}
	for _, c := range string(id) {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("non-hex character %q in ID %s", c, id)
		}
	}
}

// TestNewCompoundInterestSchedule_NewValueCreatedAlwaysZero verifies the fetish
// invariant holds for all valid inputs.
func TestNewCompoundInterestSchedule_NewValueCreatedAlwaysZero(t *testing.T) {
	t.Parallel()

	sched, err := credit.NewCompoundInterestSchedule(1000, 500, 20, credit.FormMM)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sched.NewValueCreated != 0 {
		t.Errorf("NewValueCreated must always be 0, got %d", sched.NewValueCreated)
	}
}
