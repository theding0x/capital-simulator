package valorisation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
	"github.com/theding0x/capital-simulator/pkg/surplusrate"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/valorisation"
)

// marxYearMinutes is Ch. 16's year: Marx reckons "a year of 50 weeks" (sec. I; Ch. 15
// likewise), 50 x 7 x 1,440 = 504,000 minutes. ComputeVariableCapitalReproduction only
// echoes its year argument (its totals come from n), so passing this changes no result;
// the year decides the figures through n, which is a client-supplied input.
const marxYearMinutes int64 = 50 * 7 * 24 * 60

// turnoversBP is n in basis points for a turnover period on Marx's 50-week year,
// rounded half up like every other scaled-integer rate.
func turnoversBP(periodMinutes int64) int64 {
	return surplusrate.RoundHalfUpDiv(10_000*marxYearMinutes, periodMinutes)
}

// weekMinutes is one week of seven 1,440-minute days.
const weekMinutes int64 = 7 * 24 * 60

// --- VariableCapitalAdvance.Validate -----------------------------------------

func TestVariableCapitalAdvance_Validate_Happy(t *testing.T) {
	t.Parallel()
	// Capitalist A: £500 advance, 5-week turnover (Marx Ch. 16 §I).
	adv := valorisation.VariableCapitalAdvance{
		ID:                          valorisation.NewVariableCapitalAdvanceID(),
		IndustrialCapitalID:         "test-ic-a",
		Pence:                       120000, // £500
		TurnoverPeriodMinutes:       50400,  // 5 weeks × 7 days × 1440 min
		TurnoversPerYearBasisPoints: 100000, // n=10
	}
	if err := adv.Validate(); err != nil {
		t.Fatalf("expected valid advance, got: %v", err)
	}
}

func TestVariableCapitalAdvance_Validate_ZeroPence(t *testing.T) {
	t.Parallel()
	adv := valorisation.VariableCapitalAdvance{
		Pence:                 0,
		TurnoverPeriodMinutes: 50400,
	}
	if err := adv.Validate(); err == nil {
		t.Fatal("expected error for zero pence")
	}
}

func TestVariableCapitalAdvance_Validate_ZeroTurnoverPeriod(t *testing.T) {
	t.Parallel()
	adv := valorisation.VariableCapitalAdvance{
		Pence:                 120000,
		TurnoverPeriodMinutes: 0,
	}
	if err := adv.Validate(); err == nil {
		t.Fatal("expected error for zero turnover period")
	}
}

// --- ComputePerCircuitRate ---------------------------------------------------

func TestComputePerCircuitRate_CapitalistA(t *testing.T) {
	t.Parallel()
	// Marx Ch. 16 §I: s/v = 100% → BasisPoints = 10000.
	id := valorisation.NewVariableCapitalAdvanceID()
	rate, err := valorisation.ComputePerCircuitRate(id, 120000, 120000)
	if err != nil {
		t.Fatal(err)
	}
	if rate.BasisPoints != 10000 {
		t.Errorf("expected BasisPoints=10000, got %d", rate.BasisPoints)
	}
	if rate.VariableCapitalAdvanceID != id {
		t.Error("advance ID mismatch")
	}
}

func TestComputePerCircuitRate_SpinningMill1871(t *testing.T) {
	t.Parallel()
	// The spinning mill is Vol. II Ch. 1's fixture (not Ch. 16): 50 labourers on a weekly
	// wage of £50 (v), £78 of surplus-value on the week's yarn (s), "156 per cent".
	// 18,720p / 12,000p = 156% exactly → BasisPoints = 15600. Over Ch. 16's 50-week
	// year the annual rate is 7,800% (see the annual-rate test below).
	id := valorisation.NewVariableCapitalAdvanceID()
	rate, err := valorisation.ComputePerCircuitRate(id, money.FromLSD(78, 0, 0), money.FromLSD(50, 0, 0))
	if err != nil {
		t.Fatal(err)
	}
	if rate.BasisPoints != 15600 {
		t.Errorf("expected BasisPoints=15600, got %d", rate.BasisPoints)
	}
	if rate.SurplusPence != 18720 || rate.VariablePence != 12000 {
		t.Errorf("expected s=18720p, v=12000p, got s=%d, v=%d", rate.SurplusPence, rate.VariablePence)
	}
}

func TestComputePerCircuitRate_NegativeSurplus(t *testing.T) {
	t.Parallel()
	id := valorisation.NewVariableCapitalAdvanceID()
	_, err := valorisation.ComputePerCircuitRate(id, -1, 10000)
	if err == nil {
		t.Fatal("expected error for negative surplus")
	}
}

func TestComputePerCircuitRate_ZeroVariable(t *testing.T) {
	t.Parallel()
	id := valorisation.NewVariableCapitalAdvanceID()
	_, err := valorisation.ComputePerCircuitRate(id, 5000, 0)
	if err == nil {
		t.Fatal("expected error for zero variable")
	}
}

func TestComputePerCircuitRate_ZeroSurplus(t *testing.T) {
	t.Parallel()
	// Zero surplus is valid (s/v = 0%).
	id := valorisation.NewVariableCapitalAdvanceID()
	rate, err := valorisation.ComputePerCircuitRate(id, 0, 10000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rate.BasisPoints != 0 {
		t.Errorf("expected BasisPoints=0, got %d", rate.BasisPoints)
	}
}

// --- ComputeAnnualSurplusRate ------------------------------------------------

func TestComputeAnnualSurplusRate_CapitalistA(t *testing.T) {
	t.Parallel()
	// A: n=10, s/v=100% → annual rate = 1000% (100000 bp).
	id := valorisation.NewVariableCapitalAdvanceID()
	rate := valorisation.ComputeAnnualSurplusRate(id, 100000, 10000, "1871")
	if rate.AnnualBasisPoints != 100000 {
		t.Errorf("Capitalist A: expected AnnualBasisPoints=100000, got %d", rate.AnnualBasisPoints)
	}
	if rate.Period != "1871" {
		t.Errorf("expected period=1871, got %s", rate.Period)
	}
}

func TestComputeAnnualSurplusRate_CapitalistB(t *testing.T) {
	t.Parallel()
	// B: n=1, s/v=100% → annual rate = 100% (10000 bp).
	id := valorisation.NewVariableCapitalAdvanceID()
	rate := valorisation.ComputeAnnualSurplusRate(id, 10000, 10000, "1871")
	if rate.AnnualBasisPoints != 10000 {
		t.Errorf("Capitalist B: expected AnnualBasisPoints=10000, got %d", rate.AnnualBasisPoints)
	}
}

func TestComputeAnnualSurplusRate_SpinningMill1871(t *testing.T) {
	t.Parallel()
	// Spinning Mill 1871: a one-week turnover (the week as its period is assumed, I) on
	// Ch. 16's 50-week year is n=50, and s/v=156% → annual = 50 × 156% = 7,800%. Marx gives
	// the mill no annual rate: this is derived, and it is 50 turnovers, not the 52 an
	// earlier seed used.
	// n_bp = 500000, per_circuit_bp = 15600: annual = 500000 × 15600 / 10000 = 780000.
	id := valorisation.NewVariableCapitalAdvanceID()
	if n := turnoversBP(weekMinutes); n != 500000 {
		t.Fatalf("one week on a 50-week year: n = %d bp, want 500000", n)
	}
	rate := valorisation.ComputeAnnualSurplusRate(id, 500000, 15600, "1871")
	if rate.AnnualBasisPoints != 780000 {
		t.Errorf("Spinning Mill 1871: expected AnnualBasisPoints=780000 (7800%%), got %d", rate.AnnualBasisPoints)
	}
}

func TestComputeAnnualSurplusRate_Quadrupling(t *testing.T) {
	t.Parallel()
	// Invariant: doubling both n and s/v quadruples the annual rate.
	id := valorisation.NewVariableCapitalAdvanceID()
	base := valorisation.ComputeAnnualSurplusRate(id, 10000, 10000, "test")  // n=1, s/v=100% → 10000
	double := valorisation.ComputeAnnualSurplusRate(id, 20000, 20000, "test") // n=2, s/v=200% → 40000
	if double.AnnualBasisPoints != 4*base.AnnualBasisPoints {
		t.Errorf("expected quadrupling: base=%d, doubled=%d", base.AnnualBasisPoints, double.AnnualBasisPoints)
	}
}

func TestComputeAnnualSurplusRate_ZeroN(t *testing.T) {
	t.Parallel()
	id := valorisation.NewVariableCapitalAdvanceID()
	rate := valorisation.ComputeAnnualSurplusRate(id, 0, 10000, "test")
	if rate.AnnualBasisPoints != 0 {
		t.Errorf("expected AnnualBasisPoints=0 when n=0, got %d", rate.AnnualBasisPoints)
	}
}

// --- ComputeAnnualSurplusRateContrast ----------------------------------------

func TestComputeAnnualSurplusRateContrast_CanonicalAB(t *testing.T) {
	t.Parallel()
	// A: 1000% (100000 bp), B: 100% (10000 bp) → ratio = 10× (100000 bp).
	aID := valorisation.NewVariableCapitalAdvanceID()
	bID := valorisation.NewVariableCapitalAdvanceID()
	contrast := valorisation.ComputeAnnualSurplusRateContrast(aID, bID, 100000, 10000)
	if contrast.RatioBasisPoints != 100000 {
		t.Errorf("expected RatioBasisPoints=100000, got %d", contrast.RatioBasisPoints)
	}
	if contrast.CapitalAID != aID {
		t.Error("CapitalAID mismatch")
	}
	if contrast.CapitalBID != bID {
		t.Error("CapitalBID mismatch")
	}
}

func TestComputeAnnualSurplusRateContrast_ZeroDenominator(t *testing.T) {
	t.Parallel()
	// Should not panic; returns zero ratio.
	aID := valorisation.NewVariableCapitalAdvanceID()
	bID := valorisation.NewVariableCapitalAdvanceID()
	contrast := valorisation.ComputeAnnualSurplusRateContrast(aID, bID, 100000, 0)
	if contrast.RatioBasisPoints != 0 {
		t.Errorf("expected RatioBasisPoints=0 for zero denominator, got %d", contrast.RatioBasisPoints)
	}
}

// --- ComputeVariableCapitalReproduction --------------------------------------

func TestComputeVariableCapitalReproduction_CapitalistA(t *testing.T) {
	t.Parallel()
	// A: £500 × n=10 = £5000 annual variable; s/v=100% → £5000 annual surplus.
	adv := valorisation.VariableCapitalAdvance{
		ID:                          valorisation.NewVariableCapitalAdvanceID(),
		IndustrialCapitalID:         "test-ic-a",
		Pence:                       120000, // £500
		TurnoverPeriodMinutes:       50400,
		TurnoversPerYearBasisPoints: 100000, // n=10
	}
	repro := valorisation.ComputeVariableCapitalReproduction(adv, 10000, marxYearMinutes)
	// TotalVariable = 120000 × 100000 / 10000 = 1200000p (£5000).
	if repro.TotalVariablePaidAnnualPence != 1200000 {
		t.Errorf("expected TotalVariable=1200000, got %d", repro.TotalVariablePaidAnnualPence)
	}
	// TotalSurplus = 1200000 × 10000 / 10000 = 1200000p (£5000).
	if repro.TotalSurplusAnnualPence != 1200000 {
		t.Errorf("expected TotalSurplus=1200000, got %d", repro.TotalSurplusAnnualPence)
	}
}

func TestComputeVariableCapitalReproduction_CapitalistB(t *testing.T) {
	t.Parallel()
	// B: £5000 × n=1 = £5000 annual variable; s/v=100% → £5000 annual surplus.
	// Same total annual surplus as A — the key invariant.
	adv := valorisation.VariableCapitalAdvance{
		ID:                          valorisation.NewVariableCapitalAdvanceID(),
		IndustrialCapitalID:         "test-ic-b",
		Pence:                       1200000, // £5000
		TurnoverPeriodMinutes:       504000,
		TurnoversPerYearBasisPoints: 10000, // n=1
	}
	repro := valorisation.ComputeVariableCapitalReproduction(adv, 10000, marxYearMinutes)
	if repro.TotalVariablePaidAnnualPence != 1200000 {
		t.Errorf("expected TotalVariable=1200000, got %d", repro.TotalVariablePaidAnnualPence)
	}
	if repro.TotalSurplusAnnualPence != 1200000 {
		t.Errorf("expected TotalSurplus=1200000, got %d", repro.TotalSurplusAnnualPence)
	}
}

func TestComputeVariableCapitalReproduction_SameAnnualSurplusRegardlessOfN(t *testing.T) {
	t.Parallel()
	// Invariant: same weekly labour → same total annual surplus regardless of n.
	// A and B both arrive at £5000 annual surplus (Ch. 16 §I, M); the weekly capital
	// that lays out the same £100 a week with a one-week turnover (n=50) does too (D).
	advA := valorisation.VariableCapitalAdvance{
		ID:                          valorisation.NewVariableCapitalAdvanceID(),
		Pence:                       120000,  // £500
		TurnoverPeriodMinutes:       50400,
		TurnoversPerYearBasisPoints: 100000, // n=10
	}
	advB := valorisation.VariableCapitalAdvance{
		ID:                          valorisation.NewVariableCapitalAdvanceID(),
		Pence:                       1200000, // £5000
		TurnoverPeriodMinutes:       504000,
		TurnoversPerYearBasisPoints: 10000, // n=1
	}
	advWeekly := valorisation.VariableCapitalAdvance{
		ID:                          valorisation.NewVariableCapitalAdvanceID(),
		Pence:                       money.FromLSD(100, 0, 0), // £100 a week
		TurnoverPeriodMinutes:       weekMinutes,
		TurnoversPerYearBasisPoints: turnoversBP(weekMinutes), // n=50
	}
	reproA := valorisation.ComputeVariableCapitalReproduction(advA, 10000, marxYearMinutes)
	reproB := valorisation.ComputeVariableCapitalReproduction(advB, 10000, marxYearMinutes)
	reproW := valorisation.ComputeVariableCapitalReproduction(advWeekly, 10000, marxYearMinutes)
	if reproA.TotalSurplusAnnualPence != reproB.TotalSurplusAnnualPence {
		t.Errorf(
			"same weekly labour must produce same annual surplus: A=%d B=%d",
			reproA.TotalSurplusAnnualPence, reproB.TotalSurplusAnnualPence,
		)
	}
	if reproW.TotalSurplusAnnualPence != reproA.TotalSurplusAnnualPence {
		t.Errorf(
			"a weekly turnover of the same outlay must match A and B: weekly=%d A=%d",
			reproW.TotalSurplusAnnualPence, reproA.TotalSurplusAnnualPence,
		)
	}
}

func TestComputeVariableCapitalReproduction_AtFullRate(t *testing.T) {
	t.Parallel()
	// At PerCircuitBasisPoints == 10000 (100%):
	// TotalSurplusAnnualPence == TotalVariablePaidAnnualPence.
	// Marx's weekly outlay in Ch. 16 §I is £100 (Capitals A and B), turned over
	// every week (n=50 on a 50-week year): £5,000 a year of each.
	adv := valorisation.VariableCapitalAdvance{
		ID:                          valorisation.NewVariableCapitalAdvanceID(),
		Pence:                       money.FromLSD(100, 0, 0),
		TurnoverPeriodMinutes:       weekMinutes,
		TurnoversPerYearBasisPoints: 500000, // n=50 (weekly)
	}
	repro := valorisation.ComputeVariableCapitalReproduction(adv, 10000, marxYearMinutes)
	if repro.TotalSurplusAnnualPence != repro.TotalVariablePaidAnnualPence {
		t.Errorf(
			"at 100%% per-circuit rate, surplus should equal variable: surplus=%d variable=%d",
			repro.TotalSurplusAnnualPence, repro.TotalVariablePaidAnnualPence,
		)
	}
	if want := money.FromLSD(5000, 0, 0); repro.TotalVariablePaidAnnualPence != want {
		t.Errorf("expected £5000 (%dp) a year, got %d", want, repro.TotalVariablePaidAnnualPence)
	}
}

// Vol. II Ch. 2 → Ch. 16: latent money-capital is hoarded surplus-value that
// is not yet functioning capital, so it never enters the annual rate. The
// Spinning Mill 1871 circuit (Vol. II Ch. 1's £50 wages and £78 surplus, so
// s/v = 156%) at n=1 has the same
// annual rate whatever hoard sits beside it; only capitalisation into
// variable capital changes the denominator and so the rate.
func TestAnnualSurplusRate_LatentHoardExcluded_OnlyCapitalisedVariableMoves(t *testing.T) {
	t.Parallel()
	// Magnitudes are Marx's real pounds at 240 pence to the pound (money.Pence):
	// v = £50 = 12,000p, s = £78 = 18,720p, so s/v = 156%.
	var (
		variable = money.FromLSD(50, 0, 0) // v
		surplus  = money.FromLSD(78, 0, 0) // s (s/v = 156%)
		hoard    = surplus                 // the whole surplus held latent
	)
	const n int64 = 10000
	annual := func(v valorisation.Pence) int64 {
		per, err := valorisation.ComputePerCircuitRate("a", surplus, v)
		if err != nil {
			t.Fatalf("ComputePerCircuitRate: %v", err)
		}
		return valorisation.ComputeAnnualSurplusRate("a", n, per.BasisPoints, "1871").AnnualBasisPoints
	}

	base := annual(variable)
	if base != 15600 {
		t.Fatalf("base annual rate = %d bp, want 15600 (156%%)", base)
	}
	// A latent hoard is not an input to the rate at all (v and s are the only
	// inputs), so it cannot move it; capitalising it into constant capital
	// leaves v unchanged for the same reason.
	// Capitalising the hoard into variable capital raises v and lowers s/v.
	if got := annual(variable + hoard); got >= base {
		t.Errorf("rate after capitalising hoard into v = %d bp, want < %d", got, base)
	}
}

// --- Round-half-up of scaled-integer rates (#449) -----------------------------

// The project convention is round-half-up for every scaled-integer magnitude
// derived from a ratio (pkg/surplusrate.RoundHalfUpDiv); truncating understated
// each rate by up to one basis point.
func TestComputePerCircuitRate_RoundsHalfUp(t *testing.T) {
	t.Parallel()
	// s/v = 2/3 = 6666.67 bp → 6667, not the truncated 6666.
	rate, err := valorisation.ComputePerCircuitRate(valorisation.NewVariableCapitalAdvanceID(), 2, 3)
	if err != nil {
		t.Fatal(err)
	}
	if rate.BasisPoints != 6667 {
		t.Errorf("2/3: got %d bp, want 6667", rate.BasisPoints)
	}
	// 1/3 = 3333.33 bp rounds down to 3333.
	rate, err = valorisation.ComputePerCircuitRate(valorisation.NewVariableCapitalAdvanceID(), 1, 3)
	if err != nil {
		t.Fatal(err)
	}
	if rate.BasisPoints != 3333 {
		t.Errorf("1/3: got %d bp, want 3333", rate.BasisPoints)
	}
}

func TestComputeAnnualSurplusRate_RoundsHalfUp(t *testing.T) {
	t.Parallel()
	// n = 1.5 turnovers, s/v = 33.33%: 15000 × 3333 / 10000 = 4999.5 → 5000.
	rate := valorisation.ComputeAnnualSurplusRate(valorisation.NewVariableCapitalAdvanceID(), 15000, 3333, "x")
	if rate.AnnualBasisPoints != 5000 {
		t.Errorf("got %d bp, want 5000", rate.AnnualBasisPoints)
	}
}

func TestComputeAnnualSurplusRateContrast_RoundsHalfUp(t *testing.T) {
	t.Parallel()
	// A:B = 20000:30000 = 2/3 → 6667 bp, not the truncated 6666.
	contrast := valorisation.ComputeAnnualSurplusRateContrast(
		valorisation.NewVariableCapitalAdvanceID(), valorisation.NewVariableCapitalAdvanceID(), 20000, 30000)
	if contrast.RatioBasisPoints != 6667 {
		t.Errorf("got %d bp, want 6667", contrast.RatioBasisPoints)
	}
}

// --- Marx's fixtures over his 50-week year (issue #450) ----------------------

// within reports |got-want| <= tol.
func within(got, want, tol int64) bool {
	d := got - want
	if d < 0 {
		d = -d
	}
	return d <= tol
}

// n is the number of turnovers in a year, here of 50 weeks: 50 / period-in-weeks, in
// basis points. M = Marx's own figure, D = derived from his by this arithmetic, I =
// assumed.
func TestTurnoversFromMarxFiftyWeekYear(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		weeks int64
		want  int64
	}{
		{"Capital A and the opening capital, 5 weeks (M: n = 10)", 5, 100_000},
		{"Capital B, 50 weeks (M: n = 1)", 50, 10_000},
		{"the 55-week case (Marx does not name it C), n = 50/55 = 10/11 (M)", 55, 9_091},
		{"capital I, 4 working weeks (M: turns over 12 1/2 times)", 4, 125_000},
		{"spinning mill, 1 week (D; the week as its period is I)", 1, 500_000},
		{"locomotive works, 12 weeks (D: 41,666.67 rounded half up)", 12, 41_667},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := turnoversBP(tc.weeks * weekMinutes); got != tc.want {
				t.Fatalf("n for %d weeks = %d bp, want %d", tc.weeks, got, tc.want)
			}
		})
	}
}

// Marx's ordering: the annual rate exceeds the per-circuit rate when n > 1, equals it
// when n = 1 and falls below it when n < 1. Each row's annual rate is n × s/v from the
// domain, with n taken from the period on the 50-week year.
func TestComputeAnnualSurplusRate_MarxFixturesOverFiftyWeeks(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		weeks        int64
		perCircuitBP int64
		want         int64
	}{
		{"Capital A: n = 10, s' = 100% → 1,000% (M)", 5, 10_000, 100_000},
		{"Capital B: n = 1, s' = 100% → 100% (M)", 50, 10_000, 10_000},
		{"the 55-week case: n = 10/11, s' = 100% → 90 10/11% (M)", 55, 10_000, 9_091},
		{"spinning mill: n = 50, s' = 156% → 7,800% (D)", 1, 15_600, 780_000},
		{"locomotive works: n = 50/12, s' = 156% → 650% (D)", 12, 15_600, 65_001},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			n := turnoversBP(tc.weeks * weekMinutes)
			got := valorisation.ComputeAnnualSurplusRate(valorisation.NewVariableCapitalAdvanceID(), n, tc.perCircuitBP, "1871")
			if got.AnnualBasisPoints != tc.want {
				t.Fatalf("annual = %d bp, want %d", got.AnnualBasisPoints, tc.want)
			}
			switch {
			case n > 10_000 && got.AnnualBasisPoints <= tc.perCircuitBP:
				t.Errorf("n > 1 must raise the rate above s': %d <= %d", got.AnnualBasisPoints, tc.perCircuitBP)
			case n == 10_000 && got.AnnualBasisPoints != tc.perCircuitBP:
				t.Errorf("n = 1 must equal s': %d != %d", got.AnnualBasisPoints, tc.perCircuitBP)
			case n < 10_000 && got.AnnualBasisPoints >= tc.perCircuitBP:
				t.Errorf("n < 1 must lower the rate below s': %d >= %d", got.AnnualBasisPoints, tc.perCircuitBP)
			}
		})
	}
}

// Ch. 16's 55-week case (Marx names only A and B, not this one): the premises are kept
// and the turnover lengthened to 55 weeks, so £100 a week lays out £5,500, which turns
// over 50/55 = 10/11 of a time in the year. Its annual rate is 90 10/11%, below its 100%
// per-circuit rate, and the year still lays out only £5,000 of wages.
func TestComputeAnnualSurplusRate_FiftyFiveWeekCase(t *testing.T) {
	t.Parallel()
	advance := money.FromLSD(5500, 0, 0) // 55 × £100 (M)
	if advance != 1_320_000 {
		t.Fatalf("£5,500 = %dp, want 1,320,000", advance)
	}
	id := valorisation.NewVariableCapitalAdvanceID()
	rate, err := valorisation.ComputePerCircuitRate(id, advance, advance) // s' = 100%
	if err != nil {
		t.Fatal(err)
	}
	n := turnoversBP(55 * weekMinutes)
	annual := valorisation.ComputeAnnualSurplusRate(id, n, rate.BasisPoints, "1871")
	if annual.AnnualBasisPoints != 9_091 || annual.AnnualBasisPoints >= rate.BasisPoints {
		t.Errorf("annual = %d bp, want 9,091 (90 10/11%%), below s' = %d", annual.AnnualBasisPoints, rate.BasisPoints)
	}
	if !within(11*annual.AnnualBasisPoints, 100_000, 11) { // 10/11 of 10,000 bp, to 1 bp
		t.Errorf("annual = %d bp, want 10/11 of 10,000 to within 1 bp", annual.AnnualBasisPoints)
	}
	adv := valorisation.VariableCapitalAdvance{ID: id, Pence: advance, TurnoverPeriodMinutes: 55 * weekMinutes, TurnoversPerYearBasisPoints: n}
	repro := valorisation.ComputeVariableCapitalReproduction(adv, rate.BasisPoints, marxYearMinutes)
	// n = 9,090.909 rounded up to 9,091 puts each total 12p over the exact £5,000
	// (1,200,000p). Pinned exactly: reproduction truncates (see #449), though
	// 1,320,000 × 9,091 / 10,000 divides without remainder here.
	if repro.TotalVariablePaidAnnualPence != 1_200_012 || repro.TotalSurplusAnnualPence != 1_200_012 {
		t.Errorf("reproduces v=%d s=%d, want 1,200,012 each", repro.TotalVariablePaidAnnualPence, repro.TotalSurplusAnnualPence)
	}
}

// Spinning mill (Vol. II Ch. 1: £50 wages, £78 surplus a week), a week's turnover on a
// 50-week year (the week as its turnover period is assumed): £2,500 of wages and £3,900
// of surplus a year, exactly.
func TestComputeVariableCapitalReproduction_SpinningMill_MarxWeekly(t *testing.T) {
	t.Parallel()
	v, s := money.FromLSD(50, 0, 0), money.FromLSD(78, 0, 0)
	id := valorisation.NewVariableCapitalAdvanceID()
	rate, err := valorisation.ComputePerCircuitRate(id, s, v)
	if err != nil {
		t.Fatal(err)
	}
	adv := valorisation.VariableCapitalAdvance{ID: id, Pence: v, TurnoverPeriodMinutes: weekMinutes, TurnoversPerYearBasisPoints: turnoversBP(weekMinutes)}
	repro := valorisation.ComputeVariableCapitalReproduction(adv, rate.BasisPoints, marxYearMinutes)
	if repro.TotalVariablePaidAnnualPence != money.FromLSD(2500, 0, 0) {
		t.Errorf("wages = %dp, want £2,500 = 600,000p", repro.TotalVariablePaidAnnualPence)
	}
	if repro.TotalSurplusAnnualPence != money.FromLSD(3900, 0, 0) {
		t.Errorf("surplus = %dp, want £3,900 = 936,000p", repro.TotalSurplusAnnualPence)
	}
}

// Locomotive works, derived from Ch. 12 (the spinner's week laid out for twelve): £600
// at 156% turned over 50/12 times. The exact year is £2,500 and £3,900, but n = 41,667
// bp is 41,666.67 rounded and the reproduction truncates, so the totals read 600,004p
// and 936,006p. Those 4p and 6p are n's rounding, inside half a basis point of n on the
// 144,000p advance: 7.2p for wages and, times 1.56, 11.2p for surplus. If reproduction
// later rounds half up (#449) the pinned pair changes deliberately.
func TestComputeVariableCapitalReproduction_LocomotiveWorks_WithinRounding(t *testing.T) {
	t.Parallel()
	v, s := 12*money.FromLSD(50, 0, 0), 12*money.FromLSD(78, 0, 0)
	id := valorisation.NewVariableCapitalAdvanceID()
	rate, err := valorisation.ComputePerCircuitRate(id, s, v)
	if err != nil || rate.BasisPoints != 15_600 {
		t.Fatalf("rate = %+v, %v; want 15,600 bp", rate, err)
	}
	adv := valorisation.VariableCapitalAdvance{ID: id, Pence: v, TurnoverPeriodMinutes: 12 * weekMinutes, TurnoversPerYearBasisPoints: turnoversBP(12 * weekMinutes)}
	repro := valorisation.ComputeVariableCapitalReproduction(adv, rate.BasisPoints, marxYearMinutes)
	if repro.TotalVariablePaidAnnualPence != 600_004 || repro.TotalSurplusAnnualPence != 936_006 {
		t.Errorf("reproduces v=%d s=%d, want 600,004 and 936,006", repro.TotalVariablePaidAnnualPence, repro.TotalSurplusAnnualPence)
	}
}

// Why Ch. 16 cannot reuse the 365-day year (525,600 min) the Ch. 7 turnover number
// uses: on it Capital A's n is 10.43, so its annual rate (n × 100%) reads 1,042.86%, not
// the 1,000% Marx reckons on his year of 50 weeks.
func TestAnnualRate_YearLengthChangesN(t *testing.T) {
	t.Parallel()
	const aPeriod = 5 * weekMinutes
	marxN := turnoversBP(aPeriod)
	calendarN := surplusrate.RoundHalfUpDiv(10_000*525_600, aPeriod)
	if marxN != 100_000 || calendarN != 104_286 {
		t.Fatalf("A's n: 50-week year %d bp, 365-day year %d bp; want 100,000 and 104,286", marxN, calendarN)
	}
	id := valorisation.NewVariableCapitalAdvanceID()
	for _, n := range []int64{marxN, calendarN} {
		if got := valorisation.ComputeAnnualSurplusRate(id, n, 10_000, "1871").AnnualBasisPoints; got != n {
			t.Errorf("annual rate at s' = 100%% and n = %d bp is %d bp, want n × 100%% = %d", n, got, n)
		}
	}
}
