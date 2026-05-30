package avgprofit

// Tests for Ch. 11 — Effects of General Wage Fluctuations on Prices of Production.
// All fixtures drawn directly from design-ch11.md (verified hand calculations).
// Prices are centi-units (×100); composition inputs c,v are whole units.

import "testing"

// --- helper: sphereWagePrice ---

func TestSphereWagePrice_OldPrice_50c50v(t *testing.T) {
	t.Parallel()
	// 50c+50v at factor 100 (old), rate 2000: k_centi=5000+5000=10000; price=10000*12000/10000=12000
	got := sphereWagePrice(50, 50, 100, 2000)
	if got != 12000 {
		t.Errorf("sphereWagePrice(50,50,100,2000) = %d, want 12000", got)
	}
}

func TestSphereWagePrice_NewPrice_50c50v_factor125(t *testing.T) {
	t.Parallel()
	// 50c+50v at factor 125, rate 1429: k_centi=5000+6250=11250; price=11250*11429/10000=12857 (trunc)
	got := sphereWagePrice(50, 50, 125, 1429)
	if got != 12857 {
		t.Errorf("sphereWagePrice(50,50,125,1429) = %d, want 12857", got)
	}
}

func TestSphereWagePrice_NewPrice_92c8v_factor125(t *testing.T) {
	t.Parallel()
	// 92c+8v at factor 125, rate 1429: k_centi=9200+1000=10200; price=10200*11429/10000=11657 (trunc)
	got := sphereWagePrice(92, 8, 125, 1429)
	if got != 11657 {
		t.Errorf("sphereWagePrice(92,8,125,1429) = %d, want 11657", got)
	}
}

// --- helper: baseGeneralRate ---

func TestBaseGeneralRate_80c20v_sRate10000(t *testing.T) {
	t.Parallel()
	// sv = 20*10000/10000 = 20; rate = round_half_up(20*10000/100) = 2000
	got := baseGeneralRate(80, 20, 10000)
	if got != 2000 {
		t.Errorf("baseGeneralRate(80,20,10000) = %d, want 2000", got)
	}
}

// --- helper: wageAdjustedGeneralRate ---

func TestWageAdjustedGeneralRate_factor125(t *testing.T) {
	t.Parallel()
	// livingLabour = 20 + 20*10000/10000 = 40; newV=20*125/100=25; newS=40-25=15; newC=80+25=105
	// rate = round_half_up(15*10000/105) = round_half_up(150000/105) = round_half_up(1428.57) = 1429
	got := wageAdjustedGeneralRate(80, 20, 10000, 125)
	if got != 1429 {
		t.Errorf("wageAdjustedGeneralRate(80,20,10000,125) = %d, want 1429", got)
	}
}

func TestWageAdjustedGeneralRate_factor75(t *testing.T) {
	t.Parallel()
	// livingLabour=40; newV=20*75/100=15; newS=40-15=25; newC=80+15=95
	// rate = round_half_up(25*10000/95) = round_half_up(2631.578) = 2632
	got := wageAdjustedGeneralRate(80, 20, 10000, 75)
	if got != 2632 {
		t.Errorf("wageAdjustedGeneralRate(80,20,10000,75) = %d, want 2632", got)
	}
}

// --- ComputeWageFluctuation ---

func TestComputeWageFluctuation_Rise(t *testing.T) {
	t.Parallel()
	f := ComputeWageFluctuation(125)
	if f.Kind != KindWageRise {
		t.Errorf("Kind = %q, want %q", f.Kind, KindWageRise)
	}
	if f.Factor != 125 {
		t.Errorf("Factor = %d, want 125", f.Factor)
	}
}

func TestComputeWageFluctuation_Fall(t *testing.T) {
	t.Parallel()
	f := ComputeWageFluctuation(75)
	if f.Kind != KindWageFall {
		t.Errorf("Kind = %q, want %q", f.Kind, KindWageFall)
	}
}

func TestComputeWageFluctuation_NoChange(t *testing.T) {
	t.Parallel()
	f := ComputeWageFluctuation(100)
	if f.Kind != "" {
		t.Errorf("Kind = %q, want empty string for factor==100", f.Kind)
	}
}

// --- composition classification ---

func TestCompositionClassification_Lower(t *testing.T) {
	t.Parallel()
	// base var% = 20/(80+20)*100 = 20; sphere 50c+50v var% = 50 > 20 → KindLower
	baseVarPct := int64(20)
	out := ComputeSphereWageOutcome("lower", 50, 50, 100, 2000, 2000, baseVarPct)
	if out.Composition != KindLower {
		t.Errorf("Composition = %q, want KindLower", out.Composition)
	}
}

func TestCompositionClassification_Higher(t *testing.T) {
	t.Parallel()
	// sphere 92c+8v var% = 8/(92+8)*100 = 8 < 20 → KindHigher
	baseVarPct := int64(20)
	out := ComputeSphereWageOutcome("higher", 92, 8, 100, 2000, 2000, baseVarPct)
	if out.Composition != KindHigher {
		t.Errorf("Composition = %q, want KindHigher", out.Composition)
	}
}

func TestCompositionClassification_Average(t *testing.T) {
	t.Parallel()
	// sphere 80c+20v var% = 20 == 20 → KindAverage
	baseVarPct := int64(20)
	out := ComputeSphereWageOutcome("average", 80, 20, 100, 2000, 2000, baseVarPct)
	if out.Composition != KindAverage {
		t.Errorf("Composition = %q, want KindAverage", out.Composition)
	}
}

// --- §main: ComputeWageEffectAnalysis factor 125 (wage rise) ---

func TestComputeWageEffectAnalysis_MainFixture(t *testing.T) {
	t.Parallel()
	// Base 80c+20v, s'=100% (sRate 10000), factor 125.
	// OldGeneralRate=2000, NewGeneralRate=1429.
	// AverageOutcome: KindPriceUnchanged (80c+20v new price = 10500*11429/10000 = 12000 trunc).
	// Lower 50c+50v: new price 12857 → KindPriceRises.
	// Higher 92c+8v: new price 11657 → KindPriceFalls.
	spheres := []SphereInput{
		{Name: "lower", C: 50, V: 50},
		{Name: "higher", C: 92, V: 8},
	}
	a := ComputeWageEffectAnalysis(80, 20, 10000, 125, spheres)

	if a.OldGeneralRate != 2000 {
		t.Errorf("OldGeneralRate = %d, want 2000", a.OldGeneralRate)
	}
	if a.NewGeneralRate != 1429 {
		t.Errorf("NewGeneralRate = %d, want 1429", a.NewGeneralRate)
	}
	if a.NewGeneralRate >= a.OldGeneralRate {
		t.Errorf("wage rise must lower general rate: new=%d old=%d", a.NewGeneralRate, a.OldGeneralRate)
	}
	if a.Kind != KindWageRise {
		t.Errorf("Kind = %q, want KindWageRise", a.Kind)
	}
	// AverageOutcome
	if a.AverageOutcome.Direction != KindPriceUnchanged {
		t.Errorf("AverageOutcome.Direction = %q, want KindPriceUnchanged", a.AverageOutcome.Direction)
	}
	if a.AverageOutcome.OldPrice != 12000 {
		t.Errorf("AverageOutcome.OldPrice = %d, want 12000", a.AverageOutcome.OldPrice)
	}
	if a.AverageOutcome.NewPrice != 12000 {
		t.Errorf("AverageOutcome.NewPrice = %d, want 12000", a.AverageOutcome.NewPrice)
	}
	// Lower sphere (index 0)
	if a.Outcomes[0].NewPrice != 12857 {
		t.Errorf("lower NewPrice = %d, want 12857", a.Outcomes[0].NewPrice)
	}
	if a.Outcomes[0].Direction != KindPriceRises {
		t.Errorf("lower Direction = %q, want KindPriceRises", a.Outcomes[0].Direction)
	}
	// Higher sphere (index 1)
	if a.Outcomes[1].NewPrice != 11657 {
		t.Errorf("higher NewPrice = %d, want 11657", a.Outcomes[1].NewPrice)
	}
	if a.Outcomes[1].Direction != KindPriceFalls {
		t.Errorf("higher Direction = %q, want KindPriceFalls", a.Outcomes[1].Direction)
	}
}

// --- §reverse: factor 75 (wage fall) ---

func TestComputeWageEffectAnalysis_ReverseFactor75(t *testing.T) {
	t.Parallel()
	// Base 80c+20v, s'=100%, factor 75.
	// NewGeneralRate=2632 > OldGeneralRate=2000 (wage fall raises rate).
	// Lower 50c+50v: NewPrice=8750*12632/10000=11053 < 12000 → KindPriceFalls.
	// Higher 92c+8v: NewPrice=9800*12632/10000=12379 > 12000 → KindPriceRises.
	spheres := []SphereInput{
		{Name: "lower", C: 50, V: 50},
		{Name: "higher", C: 92, V: 8},
	}
	a := ComputeWageEffectAnalysis(80, 20, 10000, 75, spheres)

	if a.NewGeneralRate != 2632 {
		t.Errorf("NewGeneralRate = %d, want 2632", a.NewGeneralRate)
	}
	if a.NewGeneralRate <= a.OldGeneralRate {
		t.Errorf("wage fall must raise general rate: new=%d old=%d", a.NewGeneralRate, a.OldGeneralRate)
	}
	if a.Kind != KindWageFall {
		t.Errorf("Kind = %q, want KindWageFall", a.Kind)
	}
	// Lower 50c+50v direction reverses: falls under wage fall
	if a.Outcomes[0].Direction != KindPriceFalls {
		t.Errorf("lower Direction = %q, want KindPriceFalls", a.Outcomes[0].Direction)
	}
	// Higher 92c+8v direction reverses: rises under wage fall
	if a.Outcomes[1].Direction != KindPriceRises {
		t.Errorf("higher Direction = %q, want KindPriceRises", a.Outcomes[1].Direction)
	}
}

// --- §zero-sum: balanced sphere set, Σ PriceChange == 0 ---

func TestComputeWageEffectAnalysis_ZeroSum(t *testing.T) {
	t.Parallel()
	// Balanced: higher 90c+10v and lower 70c+30v (mean var% = 20 = base), factor 125.
	// PriceChange higher = 11714-12000 = -286; lower = 12286-12000 = +286; sum = 0.
	spheres := []SphereInput{
		{Name: "higher", C: 90, V: 10},
		{Name: "lower", C: 70, V: 30},
	}
	a := ComputeWageEffectAnalysis(80, 20, 10000, 125, spheres)

	var sum int64
	for _, o := range a.Outcomes {
		sum += int64(o.PriceChange)
	}
	if sum != 0 {
		t.Errorf("Σ PriceChange = %d, want 0 (value conservation)", sum)
	}
	// Verify individual values from design contract
	if a.Outcomes[0].PriceChange != -286 {
		t.Errorf("higher PriceChange = %d, want -286", a.Outcomes[0].PriceChange)
	}
	if a.Outcomes[1].PriceChange != 286 {
		t.Errorf("lower PriceChange = %d, want +286", a.Outcomes[1].PriceChange)
	}
}

// --- invariant: average composition always unchanged under any factor ---

func TestInvariant_AverageCompositionUnchanged_Rise(t *testing.T) {
	t.Parallel()
	a := ComputeWageEffectAnalysis(80, 20, 10000, 125, nil)
	if a.AverageOutcome.Direction != KindPriceUnchanged {
		t.Errorf("average composition under wage rise: Direction = %q, want unchanged", a.AverageOutcome.Direction)
	}
}

func TestInvariant_AverageCompositionUnchanged_Fall(t *testing.T) {
	t.Parallel()
	a := ComputeWageEffectAnalysis(80, 20, 10000, 75, nil)
	if a.AverageOutcome.Direction != KindPriceUnchanged {
		t.Errorf("average composition under wage fall: Direction = %q, want unchanged", a.AverageOutcome.Direction)
	}
}

// --- invariant: Kind == KindWageRise iff factor > 100 ---

func TestInvariant_KindRise_iff_FactorAbove100(t *testing.T) {
	t.Parallel()
	a := ComputeWageEffectAnalysis(80, 20, 10000, 125, nil)
	if a.Kind != KindWageRise {
		t.Errorf("factor 125: Kind = %q, want KindWageRise", a.Kind)
	}
}

func TestInvariant_KindFall_iff_FactorBelow100(t *testing.T) {
	t.Parallel()
	a := ComputeWageEffectAnalysis(80, 20, 10000, 75, nil)
	if a.Kind != KindWageFall {
		t.Errorf("factor 75: Kind = %q, want KindWageFall", a.Kind)
	}
}
