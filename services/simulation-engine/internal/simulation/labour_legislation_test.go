package simulation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// Statute of Labourers 1349 — Marx's paradigm of state wage-fixing:
// maximum wage set, minimum zero, enforced by imprisonment.
var statuteOfLabourers = simulation.WageStatute{
	Period:             "1349",
	Jurisdiction:       "England",
	MaxWagePence:       12,
	MinWagePence:       0,
	EnforcementPenalty: "imprisonment",
}

// 8 George II 1740 — forbidden a higher day's wage than 2s. 7¾d for
// journeymen tailors in and around London (31 halfpennies = 2s 7½d).
var georgeIITailors = simulation.WageStatute{
	Period:             "1740s",
	Jurisdiction:       "London tailors",
	MaxWagePence:       31,
	MinWagePence:       0,
	EnforcementPenalty: "fine and imprisonment",
}

// Henry VIII vagrancy act 1530 — sturdy vagabonds whipped and
// imprisoned, tied to the cart-tail, sworn to return to birthplace.
var henryVIIIVagrancy = simulation.VagrancyLaw{
	Period:           "1530",
	Jurisdiction:     "England",
	Punishment:       "whipping+imprisonment",
	TargetPopulation: "sturdy vagabonds",
}

func TestWageStatute_ValidStatuteOfLabourers(t *testing.T) {
	t.Parallel()
	if err := statuteOfLabourers.Validate(); err != nil {
		t.Fatalf("Statute of Labourers should be valid: %v", err)
	}
}

func TestWageStatute_ValidGeorgeII(t *testing.T) {
	t.Parallel()
	if err := georgeIITailors.Validate(); err != nil {
		t.Fatalf("George II tailors statute should be valid: %v", err)
	}
}

func TestWageStatute_MissingPeriod(t *testing.T) {
	t.Parallel()
	w := statuteOfLabourers
	w.Period = ""
	if err := w.Validate(); err != simulation.ErrWageStatutePeriodRequired {
		t.Fatalf("want ErrWageStatutePeriodRequired, got %v", err)
	}
}

func TestWageStatute_MissingJurisdiction(t *testing.T) {
	t.Parallel()
	w := statuteOfLabourers
	w.Jurisdiction = ""
	if err := w.Validate(); err != simulation.ErrWageStatuteJurisdictionRequired {
		t.Fatalf("want ErrWageStatuteJurisdictionRequired, got %v", err)
	}
}

func TestWageStatute_MissingPenalty(t *testing.T) {
	t.Parallel()
	w := statuteOfLabourers
	w.EnforcementPenalty = ""
	if err := w.Validate(); err != simulation.ErrWageStatutePenaltyRequired {
		t.Fatalf("want ErrWageStatutePenaltyRequired, got %v", err)
	}
}

func TestWageStatute_NegativeMaxWage(t *testing.T) {
	t.Parallel()
	w := statuteOfLabourers
	w.MaxWagePence = -1
	if err := w.Validate(); err != simulation.ErrWageStatuteNegativePence {
		t.Fatalf("want ErrWageStatuteNegativePence, got %v", err)
	}
}

func TestWageStatute_NegativeMinWage(t *testing.T) {
	t.Parallel()
	w := statuteOfLabourers
	w.MinWagePence = -1
	if err := w.Validate(); err != simulation.ErrWageStatuteNegativePence {
		t.Fatalf("want ErrWageStatuteNegativePence, got %v", err)
	}
}

func TestWageStatute_MaxLessThanMin(t *testing.T) {
	t.Parallel()
	w := simulation.WageStatute{
		Period:             "1600",
		Jurisdiction:       "England",
		MaxWagePence:       5,
		MinWagePence:       10,
		EnforcementPenalty: "fine",
	}
	if err := w.Validate(); err != simulation.ErrWageStatuteBoundsInconsistent {
		t.Fatalf("want ErrWageStatuteBoundsInconsistent, got %v", err)
	}
}

func TestWageStatute_EqualMaxMin(t *testing.T) {
	t.Parallel()
	// A fixed wage (max == min) is valid — it is both a ceiling and a floor.
	w := simulation.WageStatute{
		Period:             "1750",
		Jurisdiction:       "England",
		MaxWagePence:       20,
		MinWagePence:       20,
		EnforcementPenalty: "fine",
	}
	if err := w.Validate(); err != nil {
		t.Fatalf("equal max/min should be valid: %v", err)
	}
}

func TestWageStatuteID_Uniqueness(t *testing.T) {
	t.Parallel()
	a := simulation.NewWageStatuteID()
	b := simulation.NewWageStatuteID()
	if a == b {
		t.Fatal("two successive IDs must differ")
	}
	if len(string(a)) != 24 {
		t.Fatalf("ID length should be 24 hex chars, got %d", len(string(a)))
	}
}

func TestVagrancyLaw_Valid(t *testing.T) {
	t.Parallel()
	if err := henryVIIIVagrancy.Validate(); err != nil {
		t.Fatalf("Henry VIII vagrancy law should be valid: %v", err)
	}
}

func TestVagrancyLaw_MissingPeriod(t *testing.T) {
	t.Parallel()
	v := henryVIIIVagrancy
	v.Period = ""
	if err := v.Validate(); err != simulation.ErrVagrancyPeriodRequired {
		t.Fatalf("want ErrVagrancyPeriodRequired, got %v", err)
	}
}

func TestVagrancyLaw_MissingJurisdiction(t *testing.T) {
	t.Parallel()
	v := henryVIIIVagrancy
	v.Jurisdiction = ""
	if err := v.Validate(); err != simulation.ErrVagrancyJurisdictionRequired {
		t.Fatalf("want ErrVagrancyJurisdictionRequired, got %v", err)
	}
}

func TestVagrancyLaw_MissingPunishment(t *testing.T) {
	t.Parallel()
	v := henryVIIIVagrancy
	v.Punishment = ""
	if err := v.Validate(); err != simulation.ErrVagrancyPunishmentRequired {
		t.Fatalf("want ErrVagrancyPunishmentRequired, got %v", err)
	}
}

func TestVagrancyLaw_MissingTarget(t *testing.T) {
	t.Parallel()
	v := henryVIIIVagrancy
	v.TargetPopulation = ""
	if err := v.Validate(); err != simulation.ErrVagrancyTargetRequired {
		t.Fatalf("want ErrVagrancyTargetRequired, got %v", err)
	}
}

func TestVagrancyLawID_Uniqueness(t *testing.T) {
	t.Parallel()
	a := simulation.NewVagrancyLawID()
	b := simulation.NewVagrancyLawID()
	if a == b {
		t.Fatal("two successive IDs must differ")
	}
	if len(string(a)) != 24 {
		t.Fatalf("ID length should be 24 hex chars, got %d", len(string(a)))
	}
}

func TestComputeStatutoryWage_WageSuppression(t *testing.T) {
	t.Parallel()
	// George II cap: 31 halfpennies max when market rate is 35 → suppression
	sw, err := simulation.ComputeStatutoryWage(31, 35)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sw.ActedWagePence != 31 {
		t.Errorf("acted: want 31, got %d", sw.ActedWagePence)
	}
	if sw.MarketWagePence != 35 {
		t.Errorf("market: want 35, got %d", sw.MarketWagePence)
	}
	if sw.Deviation >= 0 {
		t.Errorf("deviation should be negative (wage suppression), got %f", sw.Deviation)
	}
	want := float64(31-35) / float64(35)
	if sw.Deviation != want {
		t.Errorf("deviation: want %f, got %f", want, sw.Deviation)
	}
}

func TestComputeStatutoryWage_PostRepeal(t *testing.T) {
	t.Parallel()
	// After 1813 repeal: acted == market → zero deviation
	sw, err := simulation.ComputeStatutoryWage(40, 40)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sw.Deviation != 0.0 {
		t.Errorf("post-repeal deviation should be 0, got %f", sw.Deviation)
	}
}

func TestComputeStatutoryWage_ZeroMarket(t *testing.T) {
	t.Parallel()
	if _, err := simulation.ComputeStatutoryWage(10, 0); err != simulation.ErrMarketWageZero {
		t.Fatalf("want ErrMarketWageZero, got %v", err)
	}
}

func TestComputeStatutoryWage_NegativeMarket(t *testing.T) {
	t.Parallel()
	if _, err := simulation.ComputeStatutoryWage(10, -5); err != simulation.ErrMarketWageZero {
		t.Fatalf("want ErrMarketWageZero, got %v", err)
	}
}

func TestAssembleRegime_DedupMechanisms(t *testing.T) {
	t.Parallel()
	ws := []simulation.WageStatute{
		{EnforcementPenalty: "imprisonment"},
		{EnforcementPenalty: "imprisonment"}, // duplicate
	}
	vl := []simulation.VagrancyLaw{
		{Punishment: "whipping"},
		{Punishment: "imprisonment"}, // duplicate across types
	}
	regime := simulation.AssembleRegime("England 15th–18th c.", ws, vl)
	if len(regime.Mechanisms) != 2 {
		t.Errorf("want 2 deduped mechanisms, got %d: %v", len(regime.Mechanisms), regime.Mechanisms)
	}
}

func TestAssembleRegime_EmptyIsPostCapitalist(t *testing.T) {
	t.Parallel()
	// An empty regime represents post-1825 market compulsion.
	regime := simulation.AssembleRegime("post-1825", nil, nil)
	if len(regime.Mechanisms) != 0 {
		t.Errorf("empty regime should have no mechanisms, got %v", regime.Mechanisms)
	}
	if regime.WageStatutes == nil {
		t.Error("WageStatutes should be non-nil empty slice, not nil")
	}
	if regime.VagrancyLaws == nil {
		t.Error("VagrancyLaws should be non-nil empty slice, not nil")
	}
}