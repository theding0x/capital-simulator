package credit_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
)

// TestNewFictitiousCapital_ConsolFixture tests Marx's government-bond fixture:
// AnnualIncome=500 pence (£5) at 5% (500 bp) → CapitalisedValue=10000 pence (£100).
// The original £100 loan is consumed — real_capital_exists=false.
func TestNewFictitiousCapital_ConsolFixture(t *testing.T) {
	t.Parallel()

	fc, err := credit.NewFictitiousCapital(credit.KindPublicDebt, "UK Consol", 500, 500, "Government bond, original capital consumed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fc.CapitalisedValue != 10000 {
		t.Errorf("CapitalisedValue: got %d, want 10000 (500 income @ 5%%)", fc.CapitalisedValue)
	}
	if fc.RealCapitalExists {
		t.Error("RealCapitalExists must be false for KindPublicDebt")
	}
	if fc.AnnualIncome != 500 {
		t.Errorf("AnnualIncome: got %d, want 500", fc.AnnualIncome)
	}
	if fc.RateBP != 500 {
		t.Errorf("RateBP: got %d, want 500", fc.RateBP)
	}
	if fc.Kind != credit.KindPublicDebt {
		t.Errorf("Kind: got %d, want KindPublicDebt", fc.Kind)
	}
}

// TestNewFictitiousCapital_ShareFixture tests a stock (KindShare) where
// real capital exists — the productive investment is still in place.
func TestNewFictitiousCapital_ShareFixture(t *testing.T) {
	t.Parallel()

	fc, err := credit.NewFictitiousCapital(credit.KindShare, "Railway Co.", 1000, 500, "Railway stock, real capital in track and locomotives")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 1000 * 10000 / 500 = 20000
	if fc.CapitalisedValue != 20000 {
		t.Errorf("CapitalisedValue: got %d, want 20000", fc.CapitalisedValue)
	}
	if !fc.RealCapitalExists {
		t.Error("RealCapitalExists must be true for KindShare")
	}
}

// TestNewFictitiousCapital_CapitalisedValueInvariant verifies the formula
// CapitalisedValue == roundHalfUp(AnnualIncome*10000, RateBP).
func TestNewFictitiousCapital_CapitalisedValueInvariant(t *testing.T) {
	t.Parallel()

	cases := []struct {
		income   int64
		rateBP   int64
		wantCV   int64
	}{
		{500, 500, 10000},    // £5 @ 5% → £100
		{1000, 1000, 10000},  // £10 @ 10% → £100
		{300, 500, 6000},     // 300 @ 5% → 6000
		{1, 10000, 1},        // 1 @ 100% → 1 (roundHalfUp(10000,10000)=1)
	}
	for _, tc := range cases {
		tc := tc
		fc, err := credit.NewFictitiousCapital(credit.KindShare, "test", tc.income, tc.rateBP, "")
		if err != nil {
			t.Fatalf("income=%d rate=%d: unexpected error: %v", tc.income, tc.rateBP, err)
		}
		if fc.CapitalisedValue != tc.wantCV {
			t.Errorf("income=%d rate=%d: CapitalisedValue got %d, want %d", tc.income, tc.rateBP, fc.CapitalisedValue, tc.wantCV)
		}
	}
}

// TestNewFictitiousCapital_ValidationErrors tests invariant enforcement.
func TestNewFictitiousCapital_ValidationErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name    string
		income  int64
		rateBP  int64
		wantErr error
	}{
		{"zero income", 0, 500, credit.ErrNonPositiveAnnualIncome},
		{"negative income", -1, 500, credit.ErrNonPositiveAnnualIncome},
		{"zero rate", 500, 0, credit.ErrNonPositiveRateBPFict},
		{"negative rate", 500, -1, credit.ErrNonPositiveRateBPFict},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := credit.NewFictitiousCapital(credit.KindPublicDebt, "x", tc.income, tc.rateBP, "")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err != tc.wantErr {
				t.Errorf("got %v, want %v", err, tc.wantErr)
			}
		})
	}
}

// TestNewBillOfExchange_1839UKBills tests Marx's 1839 UK bills fixture.
// £132,123,460 face value, 90-day maturity (standard commercial bill).
func TestNewBillOfExchange_1839UKBills(t *testing.T) {
	t.Parallel()

	// 1839 UK bills outstanding: £132,123,460 represented as pence×100.
	bill, err := credit.NewBillOfExchange("Manchester Spinner", "London Merchant", 13212346000, 90, "1839 UK bill of exchange")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bill.FaceValue != 13212346000 {
		t.Errorf("FaceValue: got %d, want 13212346000", bill.FaceValue)
	}
	if bill.MaturityDays != 90 {
		t.Errorf("MaturityDays: got %d, want 90", bill.MaturityDays)
	}
	if bill.Drawer != "Manchester Spinner" {
		t.Errorf("Drawer: got %q", bill.Drawer)
	}
}

// TestNewBillOfExchange_ValidationErrors tests invariant enforcement.
func TestNewBillOfExchange_ValidationErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		faceValue    int64
		maturityDays int64
		wantErr      error
	}{
		{"zero face value", 0, 90, credit.ErrNonPositiveFaceValue},
		{"negative face value", -1, 90, credit.ErrNonPositiveFaceValue},
		{"zero maturity", 1000, 0, credit.ErrNonPositiveMaturityDays},
		{"negative maturity", 1000, -1, credit.ErrNonPositiveMaturityDays},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := credit.NewBillOfExchange("A", "B", tc.faceValue, tc.maturityDays, "")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err != tc.wantErr {
				t.Errorf("got %v, want %v", err, tc.wantErr)
			}
		})
	}
}

// TestNewBillOfExchangeID_Format verifies IDs are 24-char hex strings.
func TestNewBillOfExchangeID_Format(t *testing.T) {
	t.Parallel()

	id := credit.NewBillOfExchangeID()
	if len(string(id)) != 24 {
		t.Errorf("BillOfExchangeID length: got %d, want 24", len(string(id)))
	}
	for _, c := range string(id) {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("non-hex character %q in ID %s", c, id)
		}
	}
}

// TestNewFictitiousCapitalID_Format verifies IDs are 24-char hex strings.
func TestNewFictitiousCapitalID_Format(t *testing.T) {
	t.Parallel()

	id := credit.NewFictitiousCapitalID()
	if len(string(id)) != 24 {
		t.Errorf("FictitiousCapitalID length: got %d, want 24", len(string(id)))
	}
	for _, c := range string(id) {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("non-hex character %q in ID %s", c, id)
		}
	}
}
