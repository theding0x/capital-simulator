package surplus_test

import (
	"math"
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/labour"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/production"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/surplus"
)

// Fixtures. Every day below is Marx's own, cited to the chapter it comes from;
// anything else in this file is labelled illustrative where it is used.
var (
	// Ch. 12 (opening, line a–b–c): 12 h = 10 h necessary + 2 h surplus;
	// labour-power 5s, 1 h = 6d.
	ch12Day = surplus.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	// Ch. 10 §1 (line A–B): 6 h necessary labour, nothing beyond it.
	ch10Necessary = surplus.WorkingDay{Total: 360, NecessaryLabour: 360, SurplusLabour: 0}
	// Ch. 16 (eastern bread-cutter): his own wants, 12 h of necessary labour =
	// "one working day".
	breadCutterOwn = surplus.WorkingDay{Total: 720, NecessaryLabour: 720, SurplusLabour: 0}
)

// rateTolerance bounds float comparison of a rate against a ratio of integers
// that is not exactly representable (1/6, 1/3).
const rateTolerance = 1e-9

type absoluteCase struct {
	name     string
	from     surplus.WorkingDay
	extra    labour.LabourMinutes
	want     surplus.WorkingDay
	wantASV  labour.LabourMinutes
	wantRate float64
}

// absoluteCases are Ch. 10 §1's three prolongations of the 6 h line A–B by 1, 3
// and 6 h into days of 7, 9 and 12 h, where Marx gives s′ as "16 2/3, 50 and 100
// per cent". The last row prolongs day II into day III, the one row with a
// non-zero surplus baseline: ASV is the growth in surplus labour (180), not the
// new surplus (360). That row is derived: it uses Marx's two days, but the +3 h
// step is not a prolongation he performs.
var absoluteCases = []absoluteCase{
	{"Ch. 10 day I (+1 h)", ch10Necessary, 60,
		surplus.WorkingDay{Total: 420, NecessaryLabour: 360, SurplusLabour: 60}, 60, 1.0 / 6.0},
	{"Ch. 10 day II (+3 h)", ch10Necessary, 180,
		surplus.WorkingDay{Total: 540, NecessaryLabour: 360, SurplusLabour: 180}, 180, 0.5},
	{"Ch. 10 day III (+6 h)", ch10Necessary, 360,
		surplus.WorkingDay{Total: 720, NecessaryLabour: 360, SurplusLabour: 360}, 360, 1.0},
	{"Ch. 10 day II → III (derived step)", surplus.WorkingDay{Total: 540, NecessaryLabour: 360, SurplusLabour: 180}, 180,
		surplus.WorkingDay{Total: 720, NecessaryLabour: 360, SurplusLabour: 360}, 180, 1.0},
}

type relativeCase struct {
	name     string
	factor   surplus.ProductivityFactor
	want     surplus.WorkingDay
	wantRSV  labour.LabourMinutes
	wantRate float64
}

// relativeCases shorten the Ch. 12 day (s′ = 0.2). Ch. 12 gives no social
// productivity factor for these figures; 10/9 and 5/3 are derived from his
// before/after figures as old necessary labour over new.
var relativeCases = []relativeCase{
	{"Ch. 12 b moves to b′ (10 h → 9 h)", 10.0 / 9.0,
		surplus.WorkingDay{Total: 720, NecessaryLabour: 540, SurplusLabour: 180}, 60, 1.0 / 3.0},
	{"Ch. 12 labour-power 5s → 3s (10 h → 6 h)", 5.0 / 3.0,
		surplus.WorkingDay{Total: 720, NecessaryLabour: 360, SurplusLabour: 360}, 240, 1.0},
	// The 7½ h row is Marx's individual-capitalist case: productiveness doubled
	// (24 articles instead of 12), the 12 h create 8s where average social labour
	// of the same kind creates 6s ("exceptionally productive labour operates as
	// intensified labour"), labour-power is still paid 5s, so necessary labour
	// falls to 7½ h and surplus rises by 2½ h (1s → 3s; necessary:surplus 5:1
	// becomes 5:3). The extra surplus-value vanishes once the method is general.
	// 4/3 = 8s/6s = old/new necessary labour; ReduceNecessaryLabour models it as a
	// fall in necessary labour.
	{"Ch. 12 individual capitalist, 8s in 12 h (10 h → 7½ h)", 4.0 / 3.0,
		surplus.WorkingDay{Total: 720, NecessaryLabour: 450, SurplusLabour: 270}, 150, 0.6},
}

func rateOf(wd surplus.WorkingDay) float64 {
	return surplus.RateSurplusValue(labour.LabourMinutes(wd.SurplusLabour), labour.LabourMinutes(wd.NecessaryLabour))
}

// Ch. 10 §1: prolonging the line A–B (6 h necessary labour) by 1, 3 and 6 h
// leaves necessary labour fixed and lands the whole extension in surplus
// labour, an AbsoluteSurplusValue tagged "absolute".
func TestProlongWorkingDay_Ch10_MarxsWorkingDays(t *testing.T) {
	t.Parallel()
	for _, tc := range absoluteCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, asv := surplus.ProlongWorkingDay(tc.from, tc.extra)
			if got != tc.want {
				t.Fatalf("day = %+v, want %+v", got, tc.want)
			}
			if got.NecessaryLabour != tc.from.NecessaryLabour {
				t.Fatalf("NecessaryLabour = %d, want %d (fixed)", got.NecessaryLabour, tc.from.NecessaryLabour)
			}
			if !got.Partitioned() {
				t.Fatalf("partition violated: %+v", got)
			}
			if asv.LabourMinutes != tc.wantASV || asv.Origin != surplus.OriginAbsolute {
				t.Fatalf("ASV = %+v, want {%d %q}", asv.SurplusValue, tc.wantASV, surplus.OriginAbsolute)
			}
			if asv.LabourMinutes != labour.LabourMinutes(got.SurplusLabour-tc.from.SurplusLabour) {
				t.Fatalf("ASV = %d, want newSL − fromSL = %d", asv.LabourMinutes, got.SurplusLabour-tc.from.SurplusLabour)
			}
			if rate := rateOf(got); math.Abs(rate-tc.wantRate) > rateTolerance {
				t.Fatalf("s′ = %v, want %v", rate, tc.wantRate)
			}
		})
	}
}

// Ch. 12: shortening necessary labour from 10 h to 9 h (b → b′), 6 h
// (labour-power 5s → 3s) and 7½ h (the individual capitalist's 8s in 12 h)
// holds the 12 h day fixed and moves the freed minutes into surplus labour, a
// RelativeSurplusValue tagged "relative".
func TestReduceNecessaryLabour_Ch12_MarxsFigures(t *testing.T) {
	t.Parallel()
	if before := rateOf(ch12Day); math.Abs(before-0.2) > rateTolerance {
		t.Fatalf("Ch. 12 day s′ = %v, want 0.2 (necessary:surplus 5:1)", before)
	}
	for _, tc := range relativeCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, rsv := surplus.ReduceNecessaryLabour(ch12Day, tc.factor)
			if got != tc.want {
				t.Fatalf("day = %+v, want %+v", got, tc.want)
			}
			if got.Total != 720 {
				t.Fatalf("Total = %d, want 720 (fixed)", got.Total)
			}
			if !got.Partitioned() {
				t.Fatalf("partition violated: %+v", got)
			}
			if rsv.LabourMinutes != tc.wantRSV || rsv.Origin != surplus.OriginRelative {
				t.Fatalf("RSV = %+v, want {%d %q}", rsv.SurplusValue, tc.wantRSV, surplus.OriginRelative)
			}
			if rate := rateOf(got); math.Abs(rate-tc.wantRate) > rateTolerance {
				t.Fatalf("s′ = %v, want %v", rate, tc.wantRate)
			}
		})
	}
}

// Ch. 16, eastern bread-cutter: he needs 12 working hours a week for his own
// wants, which Marx equates with "one working day"; under capital he would work
// six days a week, five of them surplus labour. s′ = 5 (500%) is derived, five
// days to one; Marx states no percentage.
//
// Units: WorkingDay is period-agnostic and here Total is a week's work (six
// 12-hour days), not a 168 h week. 6 × 12 h × 60 = 4,320 min is that week's
// work, 720 is one 12-hour day, 3,600 is five of them.
func TestRateSurplusValue_Ch16_EasternBreadCutter(t *testing.T) {
	t.Parallel()
	got, asv := surplus.ProlongWorkingDay(breadCutterOwn, 5*720)
	if want := (surplus.WorkingDay{Total: 4320, NecessaryLabour: 720, SurplusLabour: 3600}); got != want {
		t.Fatalf("week = %+v, want %+v", got, want)
	}
	if want := (surplus.SurplusValue{LabourMinutes: 3600, Origin: surplus.OriginAbsolute}); asv.SurplusValue != want {
		t.Fatalf("ASV = %+v, want %+v", asv.SurplusValue, want)
	}
	if rate := surplus.RateSurplusValue(3600, 720); rate != 5.0 {
		t.Fatalf("RateSurplusValue = %v, want 5.0 (500%%)", rate)
	}
}

// millAdvance holds Marx's Ch. 16 figures as money: £400 means of production,
// £100 wages, s′ = 20% so s = £20 (= 4,800d).
func millAdvance(t *testing.T) (c, v, s money.Pence) {
	t.Helper()
	c, v = money.FromLSD(400, 0, 0), money.FromLSD(100, 0, 0)
	s = v * 20 / 100
	if s != money.FromLSD(20, 0, 0) {
		t.Fatalf("s = %dd, want £20 = 4800d", s)
	}
	return c, v, s
}

// value crosses from money to labour-minutes at the MELT (£1 = 2,400).
func value(p money.Pence) labour.LabourMinutes {
	return labour.LabourMinutes(p.ToLabourMinutes())
}

// Ch. 16, Mill's error: £400 c + £100 v, s′ 20%, so profit is 20:500 = 4%, not
// 20%. money.Pence cannot be passed unconverted to surplus.RateOfProfit;
// value() applies the MELT.
func TestRateOfProfit_Ch16_MillCritique(t *testing.T) {
	t.Parallel()
	c, v, s := millAdvance(t)
	for name, got := range map[string][2]labour.LabourMinutes{
		"c £400": {value(c), 960_000}, "v £100": {value(v), 240_000},
		"s £20": {value(s), 48_000}, "c+v £500": {value(c + v), 500 * money.LabourMinutesPerPound},
	} {
		if got[0] != got[1] {
			t.Fatalf("%s = %d labour-minutes, want %d", name, got[0], got[1])
		}
	}
	rsv := surplus.RateSurplusValue(value(s), value(v))
	rop := surplus.RateOfProfit(value(s), value(c+v))
	if rsv != 0.20 {
		t.Fatalf("RateSurplusValue = %v, want 0.20 (48,000 / 240,000)", rsv)
	}
	if rop != 0.04 {
		t.Fatalf("RateOfProfit = %v, want 0.04 (48,000 / 1,200,000)", rop)
	}
	// p′ = s′·v/(c+v); not ==, since 0.2*0.2 is 0.04000000000000001.
	if id := rsv * float64(value(v)) / float64(value(c+v)); math.Abs(id-rop) > 1e-12 {
		t.Fatalf("p′ = %v, want s′·v/(c+v) = %v", rop, id)
	}
	if !(rop < rsv) {
		t.Fatalf("rates must differ: RateOfProfit (%v) must be < RateSurplusValue (%v)", rop, rsv)
	}
}

// Invariant: the partition holds on both paths and on every input day.
func TestInvariant_PartitionHolds(t *testing.T) {
	t.Parallel()
	for _, tc := range absoluteCases {
		t.Run("absolute "+tc.name, func(t *testing.T) {
			t.Parallel()
			got, _ := surplus.ProlongWorkingDay(tc.from, tc.extra)
			if !tc.from.Partitioned() || !got.Partitioned() {
				t.Fatalf("partition violated: from %+v, to %+v", tc.from, got)
			}
		})
	}
	for _, tc := range relativeCases {
		t.Run("relative "+tc.name, func(t *testing.T) {
			t.Parallel()
			got, _ := surplus.ReduceNecessaryLabour(ch12Day, tc.factor)
			if !ch12Day.Partitioned() || !got.Partitioned() {
				t.Fatalf("partition violated: from %+v, to %+v", ch12Day, got)
			}
		})
	}
	t.Run("bread-cutter's own day", func(t *testing.T) {
		t.Parallel()
		if !breadCutterOwn.Partitioned() {
			t.Fatalf("partition violated: %+v", breadCutterOwn)
		}
	})
}

// Invariant: from a zero-surplus baseline the ASV equals the new surplus labour
// (the Ch. 10 days and the bread-cutter's own day).
func TestInvariant_AbsoluteEqualsSurplusLabour(t *testing.T) {
	t.Parallel()
	type row struct {
		name  string
		from  surplus.WorkingDay
		extra labour.LabourMinutes
	}
	var rows []row
	for _, tc := range absoluteCases {
		if tc.from.SurplusLabour == 0 {
			rows = append(rows, row{tc.name, tc.from, tc.extra})
		}
	}
	rows = append(rows, row{"Ch. 16 bread-cutter", breadCutterOwn, 5 * 720})
	for _, r := range rows {
		t.Run(r.name, func(t *testing.T) {
			t.Parallel()
			newWD, asv := surplus.ProlongWorkingDay(r.from, r.extra)
			if asv.LabourMinutes != labour.LabourMinutes(newWD.SurplusLabour) {
				t.Fatalf("ASV = %d, want SurplusLabour = %d", asv.LabourMinutes, newWD.SurplusLabour)
			}
		})
	}
}

// Invariant: RSV == oldNecessary − newNecessary == newSurplus − oldSurplus.
func TestInvariant_RelativeEqualsNecessaryReduction(t *testing.T) {
	t.Parallel()
	for _, tc := range relativeCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			newWD, rsv := surplus.ReduceNecessaryLabour(ch12Day, tc.factor)
			if got := labour.LabourMinutes(ch12Day.NecessaryLabour - newWD.NecessaryLabour); rsv.LabourMinutes != got {
				t.Fatalf("RSV = %d, want oldN − newN = %d", rsv.LabourMinutes, got)
			}
			if got := labour.LabourMinutes(newWD.SurplusLabour - ch12Day.SurplusLabour); rsv.LabourMinutes != got {
				t.Fatalf("RSV = %d, want newSL − oldSL = %d", rsv.LabourMinutes, got)
			}
		})
	}
}

// Invariant: p′ < s′ whenever constant capital is positive (Mill critique).
func TestInvariant_ProfitRateBelowSurplusRate(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		c      money.Pence
		strict bool
	}{
		{"Mill: c = £400, p′ strictly below s′", money.FromLSD(400, 0, 0), true},
		// Derived limiting case, not a Ch. 16 figure; cf. the finance-service Ch. 2
		// limiting case (finance-service/internal/profit/rate_of_profit_test.go).
		{"limiting case: c = £0, p′ = s′", 0, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, v, s := millAdvance(t)
			rsv := surplus.RateSurplusValue(value(s), value(v))
			rop := surplus.RateOfProfit(value(s), value(tc.c+v))
			switch {
			case tc.strict && !(rop < rsv):
				t.Fatalf("RateOfProfit (%v) must be < RateSurplusValue (%v) when c > 0", rop, rsv)
			case !tc.strict && math.Abs(rop-rsv) > 1e-12:
				t.Fatalf("RateOfProfit (%v) must equal RateSurplusValue (%v) when c = 0", rop, rsv)
			}
		})
	}
}

// Both surplus-value variants share one magnitude type; the distinction is the
// Origin tag (Ch. 16: "from one standpoint, any distinction between absolute
// and relative surplus-value appears illusory"). The 60 min is
// Marx's b′b, "half of b c" = 1 h (Ch. 12); the 13-hour day is derived, not
// Marx's.
func TestSurplusValue_OriginTagDistinguishes(t *testing.T) {
	t.Parallel()
	absDay, asv := surplus.ProlongWorkingDay(ch12Day, 60)
	relDay, rsv := surplus.ReduceNecessaryLabour(ch12Day, surplus.ProductivityFactor(10.0/9.0))
	if want := (surplus.WorkingDay{Total: 780, NecessaryLabour: 600, SurplusLabour: 180}); absDay != want {
		t.Fatalf("absolute day = %+v, want %+v", absDay, want)
	}
	if want := (surplus.WorkingDay{Total: 720, NecessaryLabour: 540, SurplusLabour: 180}); relDay != want {
		t.Fatalf("relative day = %+v, want %+v", relDay, want)
	}
	if asv.LabourMinutes != 60 || rsv.LabourMinutes != 60 {
		t.Fatalf("same magnitude expected: ASV %d, RSV %d, want 60 both", asv.LabourMinutes, rsv.LabourMinutes)
	}
	if asv.Origin == rsv.Origin {
		t.Fatalf("Origins must differ — got %q on both", asv.Origin)
	}
	if absDay.Total == relDay.Total || absDay.NecessaryLabour == relDay.NecessaryLabour {
		t.Fatalf("the two mechanisms must reach different days: %+v vs %+v", absDay, relDay)
	}
}

// Invariant: NecessaryLabour > 0. ReduceNecessaryLabour must not drive it to
// zero (the wage relation collapses); the floor is one minute. Factor 1e9 is a
// deliberate boundary probe, not a Marx figure.
func TestReduceNecessaryLabour_GuardsAgainstZero(t *testing.T) {
	t.Parallel()
	newWD, rsv := surplus.ReduceNecessaryLabour(ch12Day, surplus.ProductivityFactor(1e9))
	if newWD.NecessaryLabour <= 0 {
		t.Fatalf("NecessaryLabour must stay > 0 under extreme productivity; got %d", newWD.NecessaryLabour)
	}
	if newWD.Total != ch12Day.Total || !newWD.Partitioned() {
		t.Fatalf("working-day partition broken: %+v from %+v", newWD, ch12Day)
	}
	if want := (surplus.WorkingDay{Total: 720, NecessaryLabour: 1, SurplusLabour: 719}); newWD != want {
		t.Fatalf("floored day = %+v, want %+v", newWD, want)
	}
	if rsv.LabourMinutes != 599 {
		t.Fatalf("RelativeSurplusValue = %d, want 599 (600 − floor of 1)", rsv.LabourMinutes)
	}
}

// Both rates are ratios with a guarded denominator: a zero (or negative)
// necessary labour has no wage relation, so s′ is 0, and a zero capital
// advance has no profit rate, so p′ is 0.
func TestRates_ZeroDenominatorIsZero(t *testing.T) {
	t.Parallel()
	if got := surplus.RateSurplusValue(120, 0); got != 0 {
		t.Fatalf("RateSurplusValue(120, 0) = %v, want 0", got)
	}
	if got := surplus.RateSurplusValue(120, -1); got != 0 {
		t.Fatalf("RateSurplusValue(120, -1) = %v, want 0", got)
	}
	if got := surplus.RateOfProfit(120, 0); got != 0 {
		t.Fatalf("RateOfProfit(120, 0) = %v, want 0", got)
	}
}

// The Ch. 12 bridge (production.*) and the Ch. 16 synthesis (surplus.*) must
// yield the same relative surplus-value for the same day and productivity
// factor. The first two rows are derived from Ch. 12's before/after figures (it
// gives no social factor); the third is Marx's individual-capitalist case,
// 4/3 = 8s/6s = old/new necessary labour. The last three exercise the shared
// round-to-nearest and are illustrative, not Marx's.
func TestReduceNecessaryLabour_AgreesWithCh12Bridge(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		factor surplus.ProductivityFactor
		want   int64
	}{
		{"Ch. 12 b → b′ (10/9)", 10.0 / 9.0, 60},
		{"Ch. 12 labour-power 5s → 3s (5/3)", 5.0 / 3.0, 240},
		{"Ch. 12 individual capitalist, 8s in 12 h (4/3)", 4.0 / 3.0, 150},
		// Illustrative, not Marx's: 600/7 = 85.71 rounds to 86 (round-to-nearest, not floor).
		{"factor 7 rounds 85.71 to 86", 7, 514},
		// Illustrative, not Marx's: 600/13 = 46.15 rounds down to 46 (round-to-nearest,
		// not ceil): necessary 46, surplus 674.
		{"factor 13 rounds 46.15 to 46", 13, 554},
		// Illustrative, not Marx's: 600/48 = 12.5 exactly; half away from zero gives 13
		// (floor and half-even give 12).
		{"factor 48 rounds 12.5 half away to 13", 48, 587},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, rsv := surplus.ReduceNecessaryLabour(ch12Day, tc.factor)

			newLPV := production.ApplyProductivityToSNLT(labour.LabourMinutes(ch12Day.NecessaryLabour), tc.factor)
			ch12, err := production.ShortenNecessaryLabour(ch12Day, production.LabourPowerValue(newLPV))
			if err != nil {
				t.Fatalf("ShortenNecessaryLabour: %v", err)
			}
			bridgeDelta := int64(ch12.SurplusLabour - ch12Day.SurplusLabour)

			if int64(rsv.LabourMinutes) != bridgeDelta {
				t.Fatalf("Ch. 16 relative SV = %d, Ch. 12 bridge = %d; paths diverge", rsv.LabourMinutes, bridgeDelta)
			}
			if int64(rsv.LabourMinutes) != tc.want {
				t.Errorf("relative SV = %d, want %d", rsv.LabourMinutes, tc.want)
			}
			if rsv.Origin != surplus.OriginRelative {
				t.Errorf("origin = %q, want relative", rsv.Origin)
			}
		})
	}
}
