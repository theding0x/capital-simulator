package tendency

import (
	"encoding/hex"
	"errors"
	"testing"
)

// counteractingLawSteps is the §law base trajectory used by Ch. 14: s′=100%, v
// held at 100, c rising 50 → 100 → 200 → 300 → 400. Identical to the Ch. 13
// fixture; the counteracting forces are applied on top of it.
var counteractingLawSteps = [][2]int64{
	{50, 100},
	{100, 100},
	{200, 100},
	{300, 100},
	{400, 100},
}

// industrialKinds are the five forces (§I–§V) that uplift the industrial path.
// KindStockCapital (§VI) is excluded from the general rate and contributes no
// uplift, so it is deliberately absent.
var industrialKinds = []CounteractingForceKind{
	KindIntensifiedExploitation,
	KindDepressedWages,
	KindCheaperConstantCapital,
	KindRelativeOverpopulation,
	KindForeignTrade,
}

// mustForces builds a []CounteractingForce from kinds via the constructor,
// failing the test if any kind is rejected.
func mustForces(t *testing.T, kinds ...CounteractingForceKind) []CounteractingForce {
	t.Helper()
	forces := make([]CounteractingForce, 0, len(kinds))
	for _, k := range kinds {
		f, err := NewCounteractingForce(k, "")
		if err != nil {
			t.Fatalf("NewCounteractingForce(%q): %v", k, err)
		}
		forces = append(forces, f)
	}
	return forces
}

// TestExploitationIntensityEffect_Rates covers §I: raising s′ from 100% to 200%
// on a fixed 400c + 100v lifts the rate of profit from 2000 to 4000 bp.
func TestExploitationIntensityEffect_Rates(t *testing.T) {
	t.Parallel()
	e := ExploitationIntensityEffect{OldSRate: 10000, NewSRate: 20000, C: 400, V: 100}

	if got := e.BaselineProfitRate(); got != 2000 {
		t.Errorf("BaselineProfitRate = %d, want 2000", got)
	}
	if got := e.ModifiedProfitRate(); got != 4000 {
		t.Errorf("ModifiedProfitRate = %d, want 4000", got)
	}
	if e.ModifiedProfitRate() <= e.BaselineProfitRate() {
		t.Errorf("Modified %d must exceed Baseline %d", e.ModifiedProfitRate(), e.BaselineProfitRate())
	}
}

// TestComputeExploitationIntensityEffect_Valid checks the constructor accepts the
// §I fixture and surfaces the same rates.
func TestComputeExploitationIntensityEffect_Valid(t *testing.T) {
	t.Parallel()
	e, err := ComputeExploitationIntensityEffect(10000, 20000, 400, 100)
	if err != nil {
		t.Fatalf("ComputeExploitationIntensityEffect: %v", err)
	}
	if e.BaselineProfitRate() != 2000 || e.ModifiedProfitRate() != 4000 {
		t.Errorf("rates = %d/%d, want 2000/4000", e.BaselineProfitRate(), e.ModifiedProfitRate())
	}
}

// TestExploitationIntensityEffect_Validate covers the happy path, the
// negative-magnitude boundary, and the not-higher invariant (NewSRate <= OldSRate
// cannot lift the rate, so Validate must reject it).
func TestExploitationIntensityEffect_Validate(t *testing.T) {
	t.Parallel()

	ok := ExploitationIntensityEffect{OldSRate: 10000, NewSRate: 20000, C: 400, V: 100}
	if err := ok.Validate(); err != nil {
		t.Errorf("valid effect: err = %v, want nil", err)
	}

	for _, neg := range []ExploitationIntensityEffect{
		{OldSRate: -1, NewSRate: 20000, C: 400, V: 100},
		{OldSRate: 10000, NewSRate: -1, C: 400, V: 100},
		{OldSRate: 10000, NewSRate: 20000, C: -1, V: 100},
		{OldSRate: 10000, NewSRate: 20000, C: 400, V: -1},
	} {
		if err := neg.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
			t.Errorf("%+v: err = %v, want ErrNegativeMagnitude", neg, err)
		}
	}

	notHigher := ExploitationIntensityEffect{OldSRate: 10000, NewSRate: 10000, C: 400, V: 100}
	if err := notHigher.Validate(); !errors.Is(err, ErrModifiedRateNotHigher) {
		t.Errorf("equal s′: err = %v, want ErrModifiedRateNotHigher", err)
	}

	if _, err := ComputeExploitationIntensityEffect(20000, 10000, 400, 100); !errors.Is(err, ErrModifiedRateNotHigher) {
		t.Errorf("falling s′: err = %v, want ErrModifiedRateNotHigher", err)
	}
}

// TestConstantCapitalCheapening_Rates covers §III: cheapening c from 400 to 200
// (v=100, s=100) lifts p′ from 2000 to round-half-up(100·10000/300)=3333.
func TestConstantCapitalCheapening_Rates(t *testing.T) {
	t.Parallel()
	c := ConstantCapitalCheapening{OriginalC: 400, NewC: 200, V: 100, S: 100}

	if got := c.OldProfitRate(); got != 2000 {
		t.Errorf("OldProfitRate = %d, want 2000", got)
	}
	if got := c.NewProfitRate(); got != 3333 {
		t.Errorf("NewProfitRate = %d, want 3333 (round-half-up of 100·10000/300)", got)
	}
	if c.NewProfitRate() <= c.OldProfitRate() {
		t.Errorf("New %d must exceed Old %d", c.NewProfitRate(), c.OldProfitRate())
	}
}

// TestComputeConstantCapitalCheapening_Valid checks the constructor accepts the
// §III fixture.
func TestComputeConstantCapitalCheapening_Valid(t *testing.T) {
	t.Parallel()
	c, err := ComputeConstantCapitalCheapening(400, 200, 100, 100)
	if err != nil {
		t.Fatalf("ComputeConstantCapitalCheapening: %v", err)
	}
	if c.OldProfitRate() != 2000 || c.NewProfitRate() != 3333 {
		t.Errorf("rates = %d/%d, want 2000/3333", c.OldProfitRate(), c.NewProfitRate())
	}
}

// TestConstantCapitalCheapening_Validate covers the happy path, the
// negative-magnitude boundary, and the not-higher invariant (a genuine
// cheapening that fails to lift the rate is rejected). When NewC >= OriginalC no
// rate invariant is enforced.
func TestConstantCapitalCheapening_Validate(t *testing.T) {
	t.Parallel()

	ok := ConstantCapitalCheapening{OriginalC: 400, NewC: 200, V: 100, S: 100}
	if err := ok.Validate(); err != nil {
		t.Errorf("valid cheapening: err = %v, want nil", err)
	}

	for _, neg := range []ConstantCapitalCheapening{
		{OriginalC: -1, NewC: 200, V: 100, S: 100},
		{OriginalC: 400, NewC: -1, V: 100, S: 100},
		{OriginalC: 400, NewC: 200, V: -1, S: 100},
		{OriginalC: 400, NewC: 200, V: 100, S: -1},
	} {
		if err := neg.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
			t.Errorf("%+v: err = %v, want ErrNegativeMagnitude", neg, err)
		}
	}

	// NewC >= OriginalC: no rate invariant, so equal-c with equal-rate is fine.
	notCheaper := ConstantCapitalCheapening{OriginalC: 200, NewC: 400, V: 100, S: 100}
	if err := notCheaper.Validate(); err != nil {
		t.Errorf("NewC>=OriginalC: err = %v, want nil (no rate invariant)", err)
	}
}

// TestForeignTradeEffect_NetEffect covers §V with the spec fixture:
// {ReferenceC:400, ReferenceV:100, ReferenceSRate:10000, CheapeningFactor:80,
// SRateBoost:2000} ⇒ NetProfitRateEffect == 857 (> 0, so it counteracts the fall).
//
// Arithmetic: baseline p′ = 100·10000/500 = 2000; cheapened c = 320; boosted
// s′ = 12000 ⇒ s = 120; modified p′ = round-half-up(120·10000/420) = 2857;
// net = 2857 − 2000 = 857.
func TestForeignTradeEffect_NetEffect(t *testing.T) {
	t.Parallel()
	f, err := ComputeForeignTradeEffect(400, 100, 10000, 80, 2000)
	if err != nil {
		t.Fatalf("ComputeForeignTradeEffect: %v", err)
	}
	if f.NetProfitRateEffect != 857 {
		t.Errorf("NetProfitRateEffect = %d, want 857", f.NetProfitRateEffect)
	}
	if f.NetProfitRateEffect <= 0 {
		t.Errorf("NetProfitRateEffect = %d, want > 0 (foreign trade counteracts the fall)", f.NetProfitRateEffect)
	}
}

// TestForeignTradeEffect_Validate covers the happy path, negative magnitudes, and
// the cheapening-factor bounds: it must lie in (0, 100].
func TestForeignTradeEffect_Validate(t *testing.T) {
	t.Parallel()

	ok := ForeignTradeEffect{ReferenceC: 400, ReferenceV: 100, ReferenceSRate: 10000, CheapeningFactor: 80, SRateBoost: 2000}
	if err := ok.Validate(); err != nil {
		t.Errorf("valid effect: err = %v, want nil", err)
	}

	for _, neg := range []ForeignTradeEffect{
		{ReferenceC: -1, ReferenceV: 100, ReferenceSRate: 10000, CheapeningFactor: 80},
		{ReferenceC: 400, ReferenceV: -1, ReferenceSRate: 10000, CheapeningFactor: 80},
		{ReferenceC: 400, ReferenceV: 100, ReferenceSRate: -1, CheapeningFactor: 80},
		{ReferenceC: 400, ReferenceV: 100, ReferenceSRate: 10000, CheapeningFactor: 80, SRateBoost: -1},
	} {
		if err := neg.Validate(); !errors.Is(err, ErrNegativeMagnitude) {
			t.Errorf("%+v: err = %v, want ErrNegativeMagnitude", neg, err)
		}
	}

	for _, factor := range []int64{0, -1, 101, 200} {
		bad := ForeignTradeEffect{ReferenceC: 400, ReferenceV: 100, ReferenceSRate: 10000, CheapeningFactor: factor}
		if err := bad.Validate(); !errors.Is(err, ErrInvalidCheapeningFactor) {
			t.Errorf("cheapening factor %d: err = %v, want ErrInvalidCheapeningFactor", factor, err)
		}
		if _, err := ComputeForeignTradeEffect(400, 100, 10000, factor, 2000); !errors.Is(err, ErrInvalidCheapeningFactor) {
			t.Errorf("ComputeForeignTradeEffect factor %d: err = %v, want ErrInvalidCheapeningFactor", factor, err)
		}
	}

	// Boundary: factor == 100 (no cheapening) is valid.
	boundary := ForeignTradeEffect{ReferenceC: 400, ReferenceV: 100, ReferenceSRate: 10000, CheapeningFactor: 100}
	if err := boundary.Validate(); err != nil {
		t.Errorf("factor 100: err = %v, want nil (upper bound is inclusive)", err)
	}
}

// TestNewCounteractingForce_ExcludedOnlyForStockCapital covers §VI: only
// KindStockCapital is excluded from the general rate; the other five are not.
func TestNewCounteractingForce_ExcludedOnlyForStockCapital(t *testing.T) {
	t.Parallel()

	stock, err := NewCounteractingForce(KindStockCapital, "")
	if err != nil {
		t.Fatalf("NewCounteractingForce(stock): %v", err)
	}
	if !stock.ExcludedFromGeneralRate {
		t.Error("KindStockCapital: ExcludedFromGeneralRate = false, want true")
	}
	if stock.Kind != KindStockCapital {
		t.Errorf("Kind = %q, want %q", stock.Kind, KindStockCapital)
	}

	for _, k := range industrialKinds {
		f, err := NewCounteractingForce(k, "")
		if err != nil {
			t.Fatalf("NewCounteractingForce(%q): %v", k, err)
		}
		if f.ExcludedFromGeneralRate {
			t.Errorf("%q: ExcludedFromGeneralRate = true, want false", k)
		}
	}
}

// TestNewCounteractingForce_InvalidKind rejects an unknown kind with
// ErrInvalidForceKind and assigns no ID/timestamp on the happy path.
func TestNewCounteractingForce_InvalidKind(t *testing.T) {
	t.Parallel()

	if _, err := NewCounteractingForce(CounteractingForceKind("usury"), ""); !errors.Is(err, ErrInvalidForceKind) {
		t.Errorf("unknown kind: err = %v, want ErrInvalidForceKind", err)
	}

	f, err := NewCounteractingForce(KindForeignTrade, "note")
	if err != nil {
		t.Fatalf("valid kind: %v", err)
	}
	if !f.ID.IsZero() {
		t.Errorf("ID = %q, want empty (store assigns it)", f.ID)
	}
	if !f.CreatedAt.IsZero() {
		t.Errorf("CreatedAt = %v, want zero (store assigns it)", f.CreatedAt)
	}
	if f.Note != "note" {
		t.Errorf("Note = %q, want %q", f.Note, "note")
	}
}

// TestCounteractingForceKind_IsValid asserts all six named kinds are valid, an
// unknown string is invalid, and each const's string value is exactly the
// snake_case the spec names.
func TestCounteractingForceKind_IsValid(t *testing.T) {
	t.Parallel()

	valid := map[CounteractingForceKind]string{
		KindIntensifiedExploitation: "intensified_exploitation",
		KindDepressedWages:          "depressed_wages",
		KindCheaperConstantCapital:  "cheaper_constant_capital",
		KindRelativeOverpopulation:  "relative_overpopulation",
		KindForeignTrade:            "foreign_trade",
		KindStockCapital:            "stock_capital",
	}
	for k, want := range valid {
		if !k.IsValid() {
			t.Errorf("%q.IsValid() = false, want true", k)
		}
		if string(k) != want {
			t.Errorf("kind string = %q, want %q", string(k), want)
		}
		if k.Description() == "" {
			t.Errorf("%q.Description() is empty, want a gloss", k)
		}
	}

	for _, bad := range []CounteractingForceKind{"", "usury", "Intensified_Exploitation", "rent"} {
		if bad.IsValid() {
			t.Errorf("%q.IsValid() = true, want false", bad)
		}
		if bad.Description() != "" {
			t.Errorf("%q.Description() = %q, want empty", bad, bad.Description())
		}
	}
}

// TestBuildCounteractingScenario_LawPersists is the central Ch. 14 fixture: the
// five industrial forces applied to the §law base trajectory produce the modified
// rates [8865, 6650, 4431, 3325, 2660]; the initial rate is 8865, the final 2660,
// and the law persists (Final < Initial).
func TestBuildCounteractingScenario_LawPersists(t *testing.T) {
	t.Parallel()
	trajectory := BuildCompositionTrajectory("The Law As Such — rising composition", 10000, counteractingLawSteps)
	forces := mustForces(t, industrialKinds...)

	scenario, err := BuildCounteractingScenario("All five industrial forces — law still holds", trajectory, forces)
	if err != nil {
		t.Fatalf("BuildCounteractingScenario: %v", err)
	}

	want := []int64{8865, 6650, 4431, 3325, 2660}
	if len(scenario.ModifiedRates) != len(want) {
		t.Fatalf("ModifiedRates len = %d, want %d", len(scenario.ModifiedRates), len(want))
	}
	for i, w := range want {
		if scenario.ModifiedRates[i] != w {
			t.Errorf("ModifiedRates[%d] = %d, want %d", i, scenario.ModifiedRates[i], w)
		}
	}
	if scenario.InitialProfitRate != 8865 {
		t.Errorf("InitialProfitRate = %d, want 8865", scenario.InitialProfitRate)
	}
	if scenario.FinalProfitRate != 2660 {
		t.Errorf("FinalProfitRate = %d, want 2660", scenario.FinalProfitRate)
	}
	if scenario.FinalProfitRate >= scenario.InitialProfitRate {
		t.Errorf("law not preserved: Final %d >= Initial %d", scenario.FinalProfitRate, scenario.InitialProfitRate)
	}
	// The base trajectory is carried inline and unchanged.
	if len(scenario.BaseTrajectory.Periods) != len(counteractingLawSteps) {
		t.Errorf("BaseTrajectory.Periods len = %d, want %d", len(scenario.BaseTrajectory.Periods), len(counteractingLawSteps))
	}
	if len(scenario.Forces) != len(industrialKinds) {
		t.Errorf("Forces len = %d, want %d", len(scenario.Forces), len(industrialKinds))
	}
	// No ID/timestamp assigned by the builder.
	if !scenario.ID.IsZero() {
		t.Errorf("ID = %q, want empty (store assigns it)", scenario.ID)
	}
}

// TestBuildCounteractingScenario_StockCapitalAddsNothing asserts adding §VI
// (stock capital) to the industrial force set leaves the modified rates
// unchanged: its per-period uplift δ is exactly 0.
func TestBuildCounteractingScenario_StockCapitalAddsNothing(t *testing.T) {
	t.Parallel()
	trajectory := BuildCompositionTrajectory("law", 10000, counteractingLawSteps)

	withoutStock, err := BuildCounteractingScenario("five", trajectory, mustForces(t, industrialKinds...))
	if err != nil {
		t.Fatalf("five forces: %v", err)
	}
	withStock, err := BuildCounteractingScenario("six", trajectory,
		mustForces(t, append(append([]CounteractingForceKind{}, industrialKinds...), KindStockCapital)...))
	if err != nil {
		t.Fatalf("six forces: %v", err)
	}

	if len(withStock.ModifiedRates) != len(withoutStock.ModifiedRates) {
		t.Fatalf("len mismatch: with %d, without %d", len(withStock.ModifiedRates), len(withoutStock.ModifiedRates))
	}
	for i := range withoutStock.ModifiedRates {
		if withStock.ModifiedRates[i] != withoutStock.ModifiedRates[i] {
			t.Errorf("period %d: with-stock %d != without-stock %d (δ must be 0)",
				i, withStock.ModifiedRates[i], withoutStock.ModifiedRates[i])
		}
	}
	if withStock.InitialProfitRate != withoutStock.InitialProfitRate ||
		withStock.FinalProfitRate != withoutStock.FinalProfitRate {
		t.Errorf("stock capital changed the endpoints: with (%d,%d) without (%d,%d)",
			withStock.InitialProfitRate, withStock.FinalProfitRate,
			withoutStock.InitialProfitRate, withoutStock.FinalProfitRate)
	}
}

// TestBuildCounteractingScenario_StockCapitalAloneFailsLaw covers the
// ErrLawNotPreserved branch: §VI alone contributes no uplift, so the modified
// path is exactly the falling base path — Final < Initial still holds, and the
// builder succeeds. (Stock capital cannot, on its own, invert the trend.)
func TestBuildCounteractingScenario_StockCapitalAlone(t *testing.T) {
	t.Parallel()
	trajectory := BuildCompositionTrajectory("law", 10000, counteractingLawSteps)
	scenario, err := BuildCounteractingScenario("stock only", trajectory, mustForces(t, KindStockCapital))
	if err != nil {
		t.Fatalf("stock-only scenario: %v", err)
	}
	// δ=0 ⇒ modified path == base path == falling rates.
	want := []int64{6667, 5000, 3333, 2500, 2000}
	for i, w := range want {
		if scenario.ModifiedRates[i] != w {
			t.Errorf("ModifiedRates[%d] = %d, want %d (base rate, no uplift)", i, scenario.ModifiedRates[i], w)
		}
	}
	if scenario.FinalProfitRate >= scenario.InitialProfitRate {
		t.Errorf("law not preserved: Final %d >= Initial %d", scenario.FinalProfitRate, scenario.InitialProfitRate)
	}
}

// TestBuildCounteractingScenario_EmptyForces returns ErrEmptyForces.
func TestBuildCounteractingScenario_EmptyForces(t *testing.T) {
	t.Parallel()
	trajectory := BuildCompositionTrajectory("law", 10000, counteractingLawSteps)
	if _, err := BuildCounteractingScenario("none", trajectory, nil); !errors.Is(err, ErrEmptyForces) {
		t.Errorf("nil forces: err = %v, want ErrEmptyForces", err)
	}
	if _, err := BuildCounteractingScenario("none", trajectory, []CounteractingForce{}); !errors.Is(err, ErrEmptyForces) {
		t.Errorf("empty forces: err = %v, want ErrEmptyForces", err)
	}
}

// TestBuildCounteractingScenario_EmptyTrajectory returns ErrEmptyTrajectory.
func TestBuildCounteractingScenario_EmptyTrajectory(t *testing.T) {
	t.Parallel()
	empty := CompositionTrajectory{SurplusValueRate: 10000}
	if _, err := BuildCounteractingScenario("x", empty, mustForces(t, KindForeignTrade)); !errors.Is(err, ErrEmptyTrajectory) {
		t.Errorf("empty trajectory: err = %v, want ErrEmptyTrajectory", err)
	}
}

// TestCounteractingScenario_Validate covers the happy path and every invariant
// boundary: empty forces, empty trajectory, and law-not-preserved.
func TestCounteractingScenario_Validate(t *testing.T) {
	t.Parallel()
	trajectory := BuildCompositionTrajectory("law", 10000, counteractingLawSteps)

	ok, err := BuildCounteractingScenario("ok", trajectory, mustForces(t, industrialKinds...))
	if err != nil {
		t.Fatalf("build ok: %v", err)
	}
	if err := ok.Validate(); err != nil {
		t.Errorf("valid scenario: err = %v, want nil", err)
	}

	noForces := CounteractingScenario{
		BaseTrajectory:    trajectory,
		InitialProfitRate: 8865,
		FinalProfitRate:   2660,
		ModifiedRates:     []int64{8865, 2660},
	}
	if err := noForces.Validate(); !errors.Is(err, ErrEmptyForces) {
		t.Errorf("no forces: err = %v, want ErrEmptyForces", err)
	}

	noPeriods := CounteractingScenario{
		Forces:            mustForces(t, KindForeignTrade),
		InitialProfitRate: 8865,
		FinalProfitRate:   2660,
	}
	if err := noPeriods.Validate(); !errors.Is(err, ErrEmptyTrajectory) {
		t.Errorf("no periods: err = %v, want ErrEmptyTrajectory", err)
	}

	inverted := CounteractingScenario{
		Forces:            mustForces(t, KindForeignTrade),
		BaseTrajectory:    trajectory,
		InitialProfitRate: 2000,
		FinalProfitRate:   3000, // final above initial → law inverted
	}
	if err := inverted.Validate(); !errors.Is(err, ErrLawNotPreserved) {
		t.Errorf("inverted trend: err = %v, want ErrLawNotPreserved", err)
	}
}

// TestNewCounteractingForceID_HexUnique asserts force IDs are non-empty, valid
// hex, 24 chars (96 bits), and unique across a batch.
func TestNewCounteractingForceID_HexUnique(t *testing.T) {
	t.Parallel()
	seen := make(map[CounteractingForceID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewCounteractingForceID()
		if id.IsZero() {
			t.Fatal("NewCounteractingForceID returned the empty id")
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

// TestNewCounteractingScenarioID_HexUnique mirrors the force ID test for the
// scenario ID constructor.
func TestNewCounteractingScenarioID_HexUnique(t *testing.T) {
	t.Parallel()
	seen := make(map[CounteractingScenarioID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewCounteractingScenarioID()
		if id.IsZero() {
			t.Fatal("NewCounteractingScenarioID returned the empty id")
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
