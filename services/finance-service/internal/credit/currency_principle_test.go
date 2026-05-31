package credit

import (
	"strings"
	"testing"
)

// ── Ch. 34 fixtures ──────────────────────────────────────────────────────────
//
// Bank Charter Act 1844: fiduciary (unbacked) limit £14m, gold backing £14m
// → MaxNotes £28m. Crisis of 1857: Act suspended; £48,883 in excess notes
// issued. All identifiers carry the ch34 prefix to avoid collisions with
// ch31/ch32/ch33 fixtures in the shared package.

// ch34FixtureBankAct1844 returns the founding Bank Act 1844 NoteIssueConstraint.
// goldBacking=1 400 000 000, unbackedLimit=1 400 000 000 → MaxNotes=2 800 000 000.
func ch34FixtureBankAct1844(t *testing.T) NoteIssueConstraint {
	t.Helper()
	n, err := NewNoteIssueConstraint(
		"Bank Act 1844 — founding regime",
		RegimeAct1844,
		1_400_000_000,
		1_400_000_000,
		0,
		"Currency School doctrine institutionalised.",
	)
	if err != nil {
		t.Fatalf("ch34FixtureBankAct1844: %v", err)
	}
	return n
}

// ch34FixtureCrisis1857 returns the 1857-crisis suspended-Act constraint.
// notesBeyondMax=48 883 000 represents £48,883 in excess notes.
func ch34FixtureCrisis1857(t *testing.T) NoteIssueConstraint {
	t.Helper()
	n, err := NewNoteIssueConstraint(
		"Crisis of 1857 — Act suspended",
		RegimeSuspended,
		1_400_000_000,
		1_400_000_000,
		48_883_000,
		"Treasury letter suspends the Bank Act; notes issued beyond £28m maximum.",
	)
	if err != nil {
		t.Fatalf("ch34FixtureCrisis1857: %v", err)
	}
	return n
}

// ── BankActRegime.IsValid ────────────────────────────────────────────────────

func TestBankActRegime_IsValid_ValidValues(t *testing.T) {
	t.Parallel()
	cases := []struct {
		regime BankActRegime
		name   string
	}{
		{RegimePreAct1844, "RegimePreAct1844 (1)"},
		{RegimeAct1844, "RegimeAct1844 (2)"},
		{RegimeSuspended, "RegimeSuspended (3)"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if !tc.regime.IsValid() {
				t.Errorf("IsValid() = false, want true for %s", tc.name)
			}
		})
	}
}

func TestBankActRegime_IsValid_InvalidValues(t *testing.T) {
	t.Parallel()
	cases := []struct {
		regime BankActRegime
		name   string
	}{
		{BankActRegime(0), "zero"},
		{BankActRegime(4), "four"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.regime.IsValid() {
				t.Errorf("IsValid() = true, want false for regime %d", int(tc.regime))
			}
		})
	}
}

// ── BankDepartment.IsValid ───────────────────────────────────────────────────

func TestBankDepartment_IsValid_ValidValues(t *testing.T) {
	t.Parallel()
	cases := []struct {
		dept BankDepartment
		name string
	}{
		{DeptIssue, "DeptIssue (1)"},
		{DeptBanking, "DeptBanking (2)"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if !tc.dept.IsValid() {
				t.Errorf("IsValid() = false, want true for %s", tc.name)
			}
		})
	}
}

func TestBankDepartment_IsValid_InvalidValues(t *testing.T) {
	t.Parallel()
	cases := []struct {
		dept BankDepartment
		name string
	}{
		{BankDepartment(0), "zero"},
		{BankDepartment(3), "three"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.dept.IsValid() {
				t.Errorf("IsValid() = true, want false for dept %d", int(tc.dept))
			}
		})
	}
}

// ── NewNoteIssueConstraint ───────────────────────────────────────────────────

func TestNewNoteIssueConstraint_BankAct1844_MaxNotes(t *testing.T) {
	t.Parallel()
	n := ch34FixtureBankAct1844(t)
	if n.MaxNotes != 2_800_000_000 {
		t.Errorf("MaxNotes: got %d, want 2800000000", n.MaxNotes)
	}
	if n.MaxNotes != n.GoldBacking+n.UnbackedLimit {
		t.Errorf("MaxNotes invariant violated: %d != %d + %d",
			n.MaxNotes, n.GoldBacking, n.UnbackedLimit)
	}
	// ID and CreatedAt are assigned by the store, not the constructor.
	if !n.ID.IsZero() {
		t.Error("ID must be zero before store persists")
	}
	if !n.CreatedAt.IsZero() {
		t.Error("CreatedAt must be zero before store persists")
	}
}

func TestNewNoteIssueConstraint_InvalidRegime_Zero(t *testing.T) {
	t.Parallel()
	_, err := NewNoteIssueConstraint("bad", BankActRegime(0), 1_400_000_000, 1_400_000_000, 0, "")
	if err != ErrInvalidBankActRegime {
		t.Errorf("got %v, want ErrInvalidBankActRegime", err)
	}
}

func TestNewNoteIssueConstraint_InvalidRegime_Four(t *testing.T) {
	t.Parallel()
	_, err := NewNoteIssueConstraint("bad", BankActRegime(4), 1_400_000_000, 1_400_000_000, 0, "")
	if err != ErrInvalidBankActRegime {
		t.Errorf("got %v, want ErrInvalidBankActRegime", err)
	}
}

func TestNewNoteIssueConstraint_ZeroGoldBacking(t *testing.T) {
	t.Parallel()
	// goldBacking=0 → MaxNotes == unbackedLimit (fiduciary-only regime).
	n, err := NewNoteIssueConstraint("fiduciary only", RegimePreAct1844, 0, 1_400_000_000, 0, "")
	if err != nil {
		t.Fatalf("zero gold backing must be allowed: %v", err)
	}
	if n.MaxNotes != 1_400_000_000 {
		t.Errorf("MaxNotes: got %d, want 1400000000", n.MaxNotes)
	}
}

func TestNewNoteIssueConstraint_SuspendedWithNotesBeyondMax(t *testing.T) {
	t.Parallel()
	n := ch34FixtureCrisis1857(t)
	if n.Regime != RegimeSuspended {
		t.Errorf("Regime: got %d, want RegimeSuspended (%d)", n.Regime, RegimeSuspended)
	}
	if n.NotesBeyondMax != 48_883_000 {
		t.Errorf("NotesBeyondMax: got %d, want 48883000", n.NotesBeyondMax)
	}
	if n.MaxNotes != n.GoldBacking+n.UnbackedLimit {
		t.Errorf("MaxNotes invariant violated: %d != %d + %d",
			n.MaxNotes, n.GoldBacking, n.UnbackedLimit)
	}
}

// ── NewGoldDrain ─────────────────────────────────────────────────────────────

func TestNewGoldDrain_ExportOk(t *testing.T) {
	t.Parallel()
	gd, err := NewGoldDrain(500_000, true, "Gold exported to France — reserve falls.")
	if err != nil {
		t.Fatalf("amount=500000 export must be allowed: %v", err)
	}
	if gd.Amount != 500_000 {
		t.Errorf("Amount: got %d, want 500000", gd.Amount)
	}
	if !gd.IsExport {
		t.Error("IsExport: got false, want true")
	}
}

func TestNewGoldDrain_ZeroAmount(t *testing.T) {
	t.Parallel()
	gd, err := NewGoldDrain(0, false, "no movement")
	if err != nil {
		t.Fatalf("amount=0 must be allowed: %v", err)
	}
	if gd.Amount != 0 {
		t.Errorf("Amount: got %d, want 0", gd.Amount)
	}
}

func TestNewGoldDrain_NegativeAmount(t *testing.T) {
	t.Parallel()
	if _, err := NewGoldDrain(-1, true, ""); err != ErrNegativeGoldDrain {
		t.Errorf("amount=-1: got %v, want ErrNegativeGoldDrain", err)
	}
}

// ── NewCurrencyPrincipleAnalysis ─────────────────────────────────────────────

func TestNewCurrencyPrincipleAnalysis_CorrelateFalse(t *testing.T) {
	t.Parallel()
	// Marx's empirical finding: gold drain correlates with high interest, not
	// with high commodity prices.
	a, err := NewCurrencyPrincipleAnalysis(-500_000, -200, 300, false,
		"1847 crisis: gold drain during falling prices but rising interest.")
	if err != nil {
		t.Fatalf("correlates=false must be allowed: %v", err)
	}
	if a.GoldDrainCorrelatesWithPrices {
		t.Error("GoldDrainCorrelatesWithPrices: got true, want false")
	}
	if a.GoldReserveChange != -500_000 {
		t.Errorf("GoldReserveChange: got %d, want -500000", a.GoldReserveChange)
	}
	if a.InterestRateChange != 300 {
		t.Errorf("InterestRateChange: got %d, want 300", a.InterestRateChange)
	}
}

func TestNewCurrencyPrincipleAnalysis_CorrelateTrue_Errors(t *testing.T) {
	t.Parallel()
	// The Currency Principle's claim is empirically refuted; setting true must error.
	if _, err := NewCurrencyPrincipleAnalysis(0, 0, 0, true, ""); err != ErrGoldDrainMustNotCorrelate {
		t.Errorf("correlates=true: got %v, want ErrGoldDrainMustNotCorrelate", err)
	}
}

func TestNewCurrencyPrincipleAnalysis_NegativeSignedFieldsAllowed(t *testing.T) {
	t.Parallel()
	// goldReserveChange and interestRateChange are signed — negative values are valid.
	a, err := NewCurrencyPrincipleAnalysis(-1_000_000, -500, -100, false, "contraction scenario")
	if err != nil {
		t.Fatalf("negative signed fields must be allowed: %v", err)
	}
	if a.GoldReserveChange != -1_000_000 {
		t.Errorf("GoldReserveChange: got %d, want -1000000", a.GoldReserveChange)
	}
	if a.CommodityPriceChange != -500 {
		t.Errorf("CommodityPriceChange: got %d, want -500", a.CommodityPriceChange)
	}
	if a.InterestRateChange != -100 {
		t.Errorf("InterestRateChange: got %d, want -100", a.InterestRateChange)
	}
}

// ── NoteIssueConstraintID format / uniqueness ────────────────────────────────

func TestNewNoteIssueConstraintID_Format(t *testing.T) {
	t.Parallel()
	id := NewNoteIssueConstraintID()
	if len(string(id)) != 24 {
		t.Errorf("ID length: got %d, want 24 hex chars", len(string(id)))
	}
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}

func TestNewNoteIssueConstraintID_Unique(t *testing.T) {
	t.Parallel()
	a := NewNoteIssueConstraintID()
	b := NewNoteIssueConstraintID()
	if a == b {
		t.Error("two IDs must not be equal")
	}
}

func TestNoteIssueConstraintID_IsZero(t *testing.T) {
	t.Parallel()
	var id NoteIssueConstraintID
	if !id.IsZero() {
		t.Error("zero value must report IsZero==true")
	}
	id = NewNoteIssueConstraintID()
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}

// ── Ch. 34 seed-ID pattern ────────────────────────────────────────────────────

func TestCh34SeedIDs(t *testing.T) {
	t.Parallel()
	ids := []string{
		"5eed000000000000003400",
		"5eed000000000000003401",
		"5eed000000000000003402",
	}
	const prefix = "5eed00000000000000"
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if !strings.HasPrefix(id, prefix) {
			t.Errorf("seed id %q lacks prefix %q", id, prefix)
		}
		if cc := id[len(prefix) : len(prefix)+2]; cc != "34" {
			t.Errorf("seed id %q chapter token = %q, want 34", id, cc)
		}
		if _, dup := seen[id]; dup {
			t.Errorf("duplicate seed id %q", id)
		}
		seen[id] = struct{}{}
	}
}
