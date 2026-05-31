package credit

import (
	"strings"
	"testing"
)

// ── Ch. 36 fixtures ───────────────────────────────────────────────────────────
//
// Ancient Rome (StageAntiquity): rates up to 48% p.a. (4800 bp), borrowers are
// plebeians, small producers, nobles. No wage labour — surplus extracted via
// slavery and tribute. Medieval (StageMediaeval): Church-prohibited but
// ubiquitous; montes pietatis charged ~24% (2400 bp). Subordinated
// (StageSubordinated): banking system disciplines usury under industrial_capital
// at ~8% (800 bp); wage labour is universal.
//
// All identifiers carry the ch36 prefix to prevent collisions with other
// chapter fixtures in the shared package.

func ch36FixtureAncientRome(t *testing.T) UsurersCapital {
	t.Helper()
	u, err := NewUsurersCapital(
		"Ancient Roman Usury",
		StageAntiquity,
		false,
		[]string{"plebeians", "small producers", "nobles"},
		4800,
		"",
		"Usury in ancient Rome at rates of 48% p.a.",
	)
	if err != nil {
		t.Fatalf("ch36FixtureAncientRome: %v", err)
	}
	return u
}

func ch36FixtureMedieval(t *testing.T) UsurersCapital {
	t.Helper()
	u, err := NewUsurersCapital(
		"Medieval Church Usury",
		StageMediaeval,
		false,
		[]string{"peasants", "petty producers", "feudal lords"},
		2400,
		"",
		"Medieval usury under feudalism at ~24% p.a.",
	)
	if err != nil {
		t.Fatalf("ch36FixtureMedieval: %v", err)
	}
	return u
}

func ch36FixtureSubordinated(t *testing.T) UsurersCapital {
	t.Helper()
	u, err := NewUsurersCapital(
		"Modern Subordinated Usury",
		StageSubordinated,
		true,
		[]string{"industrial capitalists", "commercial capitalists"},
		800,
		"industrial_capital",
		"Usury subordinated under industrial capital via banking system.",
	)
	if err != nil {
		t.Fatalf("ch36FixtureSubordinated: %v", err)
	}
	return u
}

// ── UsuryDevelopmentStage.IsValid ─────────────────────────────────────────────

func TestCh36UsuryDevelopmentStage_IsValid_ValidValues(t *testing.T) {
	t.Parallel()
	cases := []struct {
		s    UsuryDevelopmentStage
		name string
	}{
		{StageAntiquity, "StageAntiquity (1)"},
		{StageMediaeval, "StageMediaeval (2)"},
		{StageTransition, "StageTransition (3)"},
		{StageSubordinated, "StageSubordinated (4)"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if !tc.s.IsValid() {
				t.Errorf("IsValid() = false, want true for %s", tc.name)
			}
		})
	}
}

func TestCh36UsuryDevelopmentStage_IsValid_InvalidValues(t *testing.T) {
	t.Parallel()
	cases := []struct {
		s    UsuryDevelopmentStage
		name string
	}{
		{UsuryDevelopmentStage(0), "zero"},
		{UsuryDevelopmentStage(5), "five"},
		{UsuryDevelopmentStage(-1), "negative one"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.s.IsValid() {
				t.Errorf("IsValid() = true, want false for %d", int(tc.s))
			}
		})
	}
}

// ── NewUsurersCapital: happy paths ────────────────────────────────────────────

func TestCh36NewUsurersCapital_AncientRome_Ok(t *testing.T) {
	t.Parallel()
	u := ch36FixtureAncientRome(t)
	if u.Name != "Ancient Roman Usury" {
		t.Errorf("Name: got %q, want %q", u.Name, "Ancient Roman Usury")
	}
	if u.Stage != StageAntiquity {
		t.Errorf("Stage: got %d, want StageAntiquity (%d)", u.Stage, StageAntiquity)
	}
	if u.WageLabour {
		t.Error("WageLabour: got true, want false for StageAntiquity")
	}
	if len(u.Borrowers) != 3 {
		t.Errorf("Borrowers length: got %d, want 3", len(u.Borrowers))
	}
	if u.InterestRateBP != 4800 {
		t.Errorf("InterestRateBP: got %d, want 4800", u.InterestRateBP)
	}
	if u.SubordinatedTo != "" {
		t.Errorf("SubordinatedTo: got %q, want empty", u.SubordinatedTo)
	}
	// ID and CreatedAt are assigned by the store, not the constructor.
	if !u.ID.IsZero() {
		t.Error("ID must be zero before store persists")
	}
	if !u.CreatedAt.IsZero() {
		t.Error("CreatedAt must be zero before store persists")
	}
}

func TestCh36NewUsurersCapital_Medieval_Ok(t *testing.T) {
	t.Parallel()
	u := ch36FixtureMedieval(t)
	if u.Stage != StageMediaeval {
		t.Errorf("Stage: got %d, want StageMediaeval (%d)", u.Stage, StageMediaeval)
	}
	if u.WageLabour {
		t.Error("WageLabour: got true, want false for StageMediaeval")
	}
	if u.InterestRateBP != 2400 {
		t.Errorf("InterestRateBP: got %d, want 2400", u.InterestRateBP)
	}
}

func TestCh36NewUsurersCapital_Subordinated_Ok(t *testing.T) {
	t.Parallel()
	u := ch36FixtureSubordinated(t)
	if u.Stage != StageSubordinated {
		t.Errorf("Stage: got %d, want StageSubordinated (%d)", u.Stage, StageSubordinated)
	}
	if !u.WageLabour {
		t.Error("WageLabour: got false, want true for StageSubordinated")
	}
	if u.InterestRateBP != 800 {
		t.Errorf("InterestRateBP: got %d, want 800", u.InterestRateBP)
	}
	if u.SubordinatedTo != "industrial_capital" {
		t.Errorf("SubordinatedTo: got %q, want \"industrial_capital\"", u.SubordinatedTo)
	}
}

// ── NewUsurersCapital: error paths ────────────────────────────────────────────

func TestCh36NewUsurersCapital_InvalidStage_Errors(t *testing.T) {
	t.Parallel()
	_, err := NewUsurersCapital("bad", UsuryDevelopmentStage(0), false, nil, 100, "", "")
	if err != ErrInvalidUsuryStage {
		t.Errorf("stage=0: got %v, want ErrInvalidUsuryStage", err)
	}
}

func TestCh36NewUsurersCapital_AntiquityWageLabour_Errors(t *testing.T) {
	t.Parallel()
	_, err := NewUsurersCapital("bad", StageAntiquity, true, nil, 100, "", "")
	if err != ErrWageLabourInPreCapitalist {
		t.Errorf("StageAntiquity+wageLabour=true: got %v, want ErrWageLabourInPreCapitalist", err)
	}
}

func TestCh36NewUsurersCapital_MediaevalWageLabour_Errors(t *testing.T) {
	t.Parallel()
	_, err := NewUsurersCapital("bad", StageMediaeval, true, nil, 100, "", "")
	if err != ErrWageLabourInPreCapitalist {
		t.Errorf("StageMediaeval+wageLabour=true: got %v, want ErrWageLabourInPreCapitalist", err)
	}
}

func TestCh36NewUsurersCapital_TransitionWageLabour_Ok(t *testing.T) {
	t.Parallel()
	// StageTransition allows wage labour — capitalism is emerging.
	_, err := NewUsurersCapital("transition", StageTransition, true, nil, 500, "", "")
	if err != nil {
		t.Errorf("StageTransition+wageLabour=true must be allowed: %v", err)
	}
}

func TestCh36NewUsurersCapital_NegativeRate_Errors(t *testing.T) {
	t.Parallel()
	_, err := NewUsurersCapital("bad", StageAntiquity, false, nil, -1, "", "")
	if err != ErrNegativeHistoricalRate {
		t.Errorf("rate=-1: got %v, want ErrNegativeHistoricalRate", err)
	}
}

func TestCh36NewUsurersCapital_SubordinatedEmptySubordinatedTo_Errors(t *testing.T) {
	t.Parallel()
	_, err := NewUsurersCapital("bad", StageSubordinated, true, nil, 800, "", "")
	if err != ErrMissingSubordination {
		t.Errorf("StageSubordinated+subordinated_to=\"\": got %v, want ErrMissingSubordination", err)
	}
}

func TestCh36NewUsurersCapital_SubordinatedMerchantCapital_Errors(t *testing.T) {
	t.Parallel()
	// Only "industrial_capital" is the admissible value.
	_, err := NewUsurersCapital("bad", StageSubordinated, true, nil, 800, "merchant_capital", "")
	if err != ErrMissingSubordination {
		t.Errorf("StageSubordinated+subordinated_to=\"merchant_capital\": got %v, want ErrMissingSubordination", err)
	}
}

func TestCh36NewUsurersCapital_AntiquityWithSubordinatedTo_Ok(t *testing.T) {
	t.Parallel()
	// subordinated_to is ignored (not checked) for early stages.
	_, err := NewUsurersCapital("ancient", StageAntiquity, false, nil, 4800, "anything", "")
	if err != nil {
		t.Errorf("StageAntiquity+subordinated_to=\"anything\" must be allowed: %v", err)
	}
}

func TestCh36NewUsurersCapital_ZeroRate_Ok(t *testing.T) {
	t.Parallel()
	_, err := NewUsurersCapital("zero-rate", StageAntiquity, false, nil, 0, "", "")
	if err != nil {
		t.Errorf("rate=0 must be allowed: %v", err)
	}
}

// ── NewHistoricalInterestRate ─────────────────────────────────────────────────

func TestCh36NewHistoricalInterestRate_RomeFourtyEightPct_Ok(t *testing.T) {
	t.Parallel()
	h, err := NewHistoricalInterestRate("Ancient Rome", 4800, "48% per annum usury in ancient Rome.")
	if err != nil {
		t.Fatalf("rate=4800 must be allowed: %v", err)
	}
	if h.RateBP != 4800 {
		t.Errorf("RateBP: got %d, want 4800", h.RateBP)
	}
	if h.Period != "Ancient Rome" {
		t.Errorf("Period: got %q, want %q", h.Period, "Ancient Rome")
	}
}

func TestCh36NewHistoricalInterestRate_NegativeRate_Errors(t *testing.T) {
	t.Parallel()
	_, err := NewHistoricalInterestRate("bad", -1, "")
	if err != ErrNegativeHistoricalRate {
		t.Errorf("rate=-1: got %v, want ErrNegativeHistoricalRate", err)
	}
}

func TestCh36NewHistoricalInterestRate_ZeroRate_Ok(t *testing.T) {
	t.Parallel()
	h, err := NewHistoricalInterestRate("at par", 0, "")
	if err != nil {
		t.Fatalf("rate=0 must be allowed: %v", err)
	}
	if h.RateBP != 0 {
		t.Errorf("RateBP: got %d, want 0", h.RateBP)
	}
}

// ── NewPreCapitalistCreditForm ────────────────────────────────────────────────

func TestCh36NewPreCapitalistCreditForm_TempleLoan_Ok(t *testing.T) {
	t.Parallel()
	f := NewPreCapitalistCreditForm("Temple Loan", "Loans extended by temple treasuries in antiquity.")
	if f.ID == "" {
		t.Error("ID must be non-empty after construction")
	}
	if f.Name != "Temple Loan" {
		t.Errorf("Name: got %q, want %q", f.Name, "Temple Loan")
	}
	if f.Description != "Loans extended by temple treasuries in antiquity." {
		t.Errorf("Description mismatch: %q", f.Description)
	}
}

// ── UsurersCapitalID format and uniqueness ────────────────────────────────────

func TestCh36NewUsurersCapitalID_Format(t *testing.T) {
	t.Parallel()
	id := NewUsurersCapitalID()
	if len(string(id)) != 24 {
		t.Errorf("UsurersCapitalID length: got %d, want 24 hex chars", len(string(id)))
	}
	if id.IsZero() {
		t.Error("fresh UsurersCapitalID must not be zero")
	}
}

func TestCh36NewUsurersCapitalID_Unique(t *testing.T) {
	t.Parallel()
	a := NewUsurersCapitalID()
	b := NewUsurersCapitalID()
	if a == b {
		t.Error("two UsurersCapitalIDs must not be equal")
	}
}

func TestCh36UsurersCapitalID_IsZero(t *testing.T) {
	t.Parallel()
	var id UsurersCapitalID
	if !id.IsZero() {
		t.Error("zero-value UsurersCapitalID must report IsZero==true")
	}
}

// ── Ch. 36 seed-ID pattern ────────────────────────────────────────────────────

func TestCh36SeedIDs(t *testing.T) {
	t.Parallel()
	ids := []string{
		"5eed000000000000003600",
		"5eed000000000000003601",
		"5eed000000000000003602",
	}
	const prefix = "5eed00000000000000"
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if !strings.HasPrefix(id, prefix) {
			t.Errorf("seed id %q lacks prefix %q", id, prefix)
		}
		// The two chars immediately after the prefix are the chapter token.
		if cc := id[len(prefix) : len(prefix)+2]; cc != "36" {
			t.Errorf("seed id %q chapter token = %q, want 36", id, cc)
		}
		if len(id) != 22 {
			t.Errorf("seed id %q length = %d, want 22", id, len(id))
		}
		if _, dup := seen[id]; dup {
			t.Errorf("duplicate seed id %q", id)
		}
		seen[id] = struct{}{}
	}
}
