package credit

import "testing"

// ── FictitiousCapitalValuation ────────────────────────────────────────────────

func TestNewFictitiousCapitalValuation_GovernmentBondThreeRates(t *testing.T) {
	t.Parallel()

	// £100 government bond @ 5%: AnnualIncome=500, InterestRateBP=500 → 10000.
	at5pct, err := NewFictitiousCapitalValuation("Consol @5%", 10000, 500, 500, 0, "")
	if err != nil {
		t.Fatal(err)
	}
	if at5pct.MarketValue != 10000 {
		t.Errorf("@5%%: MarketValue=%d, want 10000", at5pct.MarketValue)
	}

	// Same income at 2.5% (250 bp) → 20000 (rate halved, value doubled).
	at2pct, err := NewFictitiousCapitalValuation("Consol @2.5%", 10000, 500, 250, 0, "")
	if err != nil {
		t.Fatal(err)
	}
	if at2pct.MarketValue != 20000 {
		t.Errorf("@2.5%%: MarketValue=%d, want 20000", at2pct.MarketValue)
	}

	// Same income at 10% (1000 bp) → 5000 (rate doubled, value halved).
	at10pct, err := NewFictitiousCapitalValuation("Consol @10%", 10000, 500, 1000, 0, "")
	if err != nil {
		t.Fatal(err)
	}
	if at10pct.MarketValue != 5000 {
		t.Errorf("@10%%: MarketValue=%d, want 5000", at10pct.MarketValue)
	}

	// Inverse relation: as the rate rises, market value falls strictly.
	if !(at2pct.MarketValue > at5pct.MarketValue && at5pct.MarketValue > at10pct.MarketValue) {
		t.Errorf("inverse relation violated: 250bp=%d, 500bp=%d, 1000bp=%d",
			at2pct.MarketValue, at5pct.MarketValue, at10pct.MarketValue)
	}
}

func TestNewFictitiousCapitalValuation_Stock(t *testing.T) {
	t.Parallel()

	// Stock yielding 10% on nominal 10000: income = 10000×1000/10000 = 1000.
	// Capitalised at the prevailing 5% (500 bp): 1000×10000/500 = 20000.
	income := int64(10000) * 1000 / basisPoints
	v, err := NewFictitiousCapitalValuation("Joint-Stock Share", 10000, income, 500, 1000, "")
	if err != nil {
		t.Fatal(err)
	}
	if v.MarketValue != 20000 {
		t.Errorf("stock MarketValue=%d, want 20000", v.MarketValue)
	}
	if v.DividendRateBP != 1000 {
		t.Errorf("DividendRateBP=%d, want 1000", v.DividendRateBP)
	}
}

func TestNewFictitiousCapitalValuation_NonPositiveRate(t *testing.T) {
	t.Parallel()

	if _, err := NewFictitiousCapitalValuation("zero", 10000, 500, 0, 0, ""); err != ErrNonPositiveInterestRate {
		t.Errorf("rate=0: err=%v, want ErrNonPositiveInterestRate", err)
	}
	if _, err := NewFictitiousCapitalValuation("neg", 10000, 500, -100, 0, ""); err != ErrNonPositiveInterestRate {
		t.Errorf("rate<0: err=%v, want ErrNonPositiveInterestRate", err)
	}
}

func TestNewFictitiousCapitalValuation_DerivedInvariant(t *testing.T) {
	t.Parallel()

	v, err := NewFictitiousCapitalValuation("Invariant", 0, 737, 333, 0, "")
	if err != nil {
		t.Fatal(err)
	}
	if want := roundHalfUp(737*basisPoints, 333); v.MarketValue != want {
		t.Errorf("MarketValue=%d, want roundHalfUp result %d", v.MarketValue, want)
	}
}

// ── BankCapital ───────────────────────────────────────────────────────────────

func TestNewBankCapital_ScotlandTotal(t *testing.T) {
	t.Parallel()

	// Scottish bank: cash £3m, securities £24m → total £27m.
	components := []BankCapitalComponent{
		{Amount: 300000000, Kind: KindCash, Description: "gold and notes in till"},
		{Amount: 1400000000, Kind: KindBillOfExchange, Description: "discounted commercial paper"},
		{Amount: 1000000000, Kind: KindGovernmentBond, Description: "consols"},
	}
	bc, err := NewBankCapital("Scotland Bank", 300000000, 2400000000, components, "")
	if err != nil {
		t.Fatal(err)
	}
	if bc.TotalCapital != 2700000000 {
		t.Errorf("TotalCapital=%d, want 2700000000", bc.TotalCapital)
	}
	if bc.TotalCapital != bc.CashAmount+bc.SecuritiesAmount {
		t.Errorf("invariant violated: %d != %d + %d", bc.TotalCapital, bc.CashAmount, bc.SecuritiesAmount)
	}
	if len(bc.Components) != 3 {
		t.Errorf("Components length=%d, want 3", len(bc.Components))
	}
}

func TestNewBankCapital_ComputesTotalEvenWhenInconsistent(t *testing.T) {
	t.Parallel()

	// Constructor computes total; it never trusts a supplied figure.
	bc, err := NewBankCapital("Computed", 100, 250, nil, "")
	if err != nil {
		t.Fatal(err)
	}
	if bc.TotalCapital != 350 {
		t.Errorf("TotalCapital=%d, want 350", bc.TotalCapital)
	}
}

func TestNewBankCapital_InvalidComponentKind(t *testing.T) {
	t.Parallel()

	bad := []BankCapitalComponent{{Amount: 10, Kind: BankCapitalComponentKind(6), Description: "bad"}}
	if _, err := NewBankCapital("Bad", 10, 0, bad, ""); err != ErrInvalidComponentKind {
		t.Errorf("err=%v, want ErrInvalidComponentKind", err)
	}

	zero := []BankCapitalComponent{{Amount: 10, Kind: BankCapitalComponentKind(0), Description: "bad"}}
	if _, err := NewBankCapital("Zero", 10, 0, zero, ""); err != ErrInvalidComponentKind {
		t.Errorf("err=%v, want ErrInvalidComponentKind", err)
	}
}

// ── BankCapitalComponentKind ──────────────────────────────────────────────────

func TestBankCapitalComponentKind_IsValid(t *testing.T) {
	t.Parallel()

	for k := KindCash; k <= KindMortgage; k++ {
		if !k.IsValid() {
			t.Errorf("kind %d should be valid", k)
		}
	}
	for _, k := range []BankCapitalComponentKind{0, 6, -1, 99} {
		if k.IsValid() {
			t.Errorf("kind %d should be invalid", k)
		}
	}
}

// ── DepositEntry ──────────────────────────────────────────────────────────────

func TestNewDepositEntry_NonNegative(t *testing.T) {
	t.Parallel()

	d, err := NewDepositEntry(5000, "current account book entry")
	if err != nil {
		t.Fatal(err)
	}
	if d.Amount != 5000 {
		t.Errorf("Amount=%d, want 5000", d.Amount)
	}
	if d.ID == "" {
		t.Error("expected non-empty DepositEntryID")
	}

	zero, err := NewDepositEntry(0, "empty")
	if err != nil {
		t.Fatalf("zero amount should be allowed: %v", err)
	}
	if zero.Amount != 0 {
		t.Errorf("Amount=%d, want 0", zero.Amount)
	}

	if _, err := NewDepositEntry(-1, "negative"); err != ErrNegativeDeposit {
		t.Errorf("err=%v, want ErrNegativeDeposit", err)
	}
}

// ── ID format / uniqueness ────────────────────────────────────────────────────

func TestBankCapitalIDs_FormatAndUniqueness(t *testing.T) {
	t.Parallel()

	seen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id := NewBankCapitalID()
		s := string(id)
		if len(s) != 24 {
			t.Fatalf("BankCapitalID length=%d, want 24 (%q)", len(s), s)
		}
		if seen[s] {
			t.Fatalf("duplicate BankCapitalID %q", s)
		}
		seen[s] = true
	}
	if NewBankCapitalID().IsZero() {
		t.Error("fresh BankCapitalID must not be zero")
	}

	vseen := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		id := NewFictitiousCapitalValuationID()
		s := string(id)
		if len(s) != 24 {
			t.Fatalf("FictitiousCapitalValuationID length=%d, want 24 (%q)", len(s), s)
		}
		if vseen[s] {
			t.Fatalf("duplicate FictitiousCapitalValuationID %q", s)
		}
		vseen[s] = true
	}
}
