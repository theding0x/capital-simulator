package profit

import (
	"errors"
	"testing"
)

// Marx's §I illustration, rising branch: a capital of £320 fixed + £80 raw
// material + £100 v + £100 s shows p' = 20%. Doubling the price of the raw
// material swells the circulating component from £80 to £160, so c rises from
// £400 to £480 and C' = £580; the rate of profit falls to 100/580 = 17.24%.
func TestConstantCapitalPriceEffectDoubling(t *testing.T) {
	t.Parallel()

	e := ConstantCapitalPriceEffect{
		FixedCapital:            320,
		OriginalMaterialCapital: 80,
		PriceFactor:             200, // the price has doubled
		VariableCapital:         100,
		SurplusValue:            100,
	}
	if got := e.OriginalConstantCapital(); got != 400 {
		t.Errorf("original c: got %d, want 400", got)
	}
	if got := e.NewConstantCapital(); got != 480 {
		t.Errorf("new c: got %d, want 480", got)
	}
	if got := e.OldProfitRate(); got != 2000 {
		t.Errorf("old p': got %d bp, want 2000", got)
	}
	if got := e.NewProfitRate(); got != 1724 {
		t.Errorf("new p': got %d bp, want 1724", got)
	}
	// Invariant [§I]: a price rise lowers the rate of profit (s, v unchanged).
	if e.Kind() != KindRise {
		t.Errorf("kind: got %q, want rise", e.Kind())
	}
	if e.NewProfitRate() >= e.OldProfitRate() {
		t.Errorf("price rise must lower p': new %d >= old %d", e.NewProfitRate(), e.OldProfitRate())
	}
}

// Marx's §I illustration, falling branch: a capital of £20 fixed + £380 raw
// material + £100 v + £100 s also shows p' = 20%. Halving the price of the raw
// material shrinks the circulating component from £380 to £190, so c falls from
// £400 to £210 and C' = £310; the rate of profit rises to 100/310 = 32.26%
// (3226 bp, rounded half-up — the Vol. III house convention).
func TestConstantCapitalPriceEffectHalving(t *testing.T) {
	t.Parallel()

	e := ConstantCapitalPriceEffect{
		FixedCapital:            20,
		OriginalMaterialCapital: 380,
		PriceFactor:             50, // the price has halved
		VariableCapital:         100,
		SurplusValue:            100,
	}
	if got := e.NewConstantCapital(); got != 210 {
		t.Errorf("new c: got %d, want 210", got)
	}
	if got := e.OldProfitRate(); got != 2000 {
		t.Errorf("old p': got %d bp, want 2000", got)
	}
	if got := e.NewProfitRate(); got != 3226 {
		t.Errorf("new p': got %d bp, want 3226", got)
	}
	// Invariant [§I]: a price fall raises the rate of profit (s, v unchanged).
	if e.Kind() != KindFall {
		t.Errorf("kind: got %q, want fall", e.Kind())
	}
	if e.NewProfitRate() <= e.OldProfitRate() {
		t.Errorf("price fall must raise p': new %d <= old %d", e.NewProfitRate(), e.OldProfitRate())
	}
}

// A price factor of 100 is no fluctuation at all: c, and hence the rate of
// profit, is unchanged. With no fixed or material capital there is no rate.
func TestConstantCapitalPriceEffectBoundaries(t *testing.T) {
	t.Parallel()

	noChange := ConstantCapitalPriceEffect{FixedCapital: 320, OriginalMaterialCapital: 80, PriceFactor: 100, VariableCapital: 100, SurplusValue: 100}
	if noChange.NewProfitRate() != noChange.OldProfitRate() {
		t.Errorf("factor 100: new %d != old %d", noChange.NewProfitRate(), noChange.OldProfitRate())
	}

	empty := ConstantCapitalPriceEffect{PriceFactor: 200}
	if got := empty.OldProfitRate(); got != 0 {
		t.Errorf("zero capital: old p' got %d, want 0", got)
	}
	if got := empty.NewProfitRate(); got != 0 {
		t.Errorf("zero capital: new p' got %d, want 0", got)
	}
}

// ComputePriceFluctuationAnalysis derives the rates and direction but assigns no
// ID or timestamp — the store does that on persistence.
func TestComputePriceFluctuationAnalysis(t *testing.T) {
	t.Parallel()

	a := ComputePriceFluctuationAnalysis(ConstantCapitalPriceEffect{
		FixedCapital:            320,
		OriginalMaterialCapital: 80,
		PriceFactor:             200,
		VariableCapital:         100,
		SurplusValue:            100,
	})
	if a.OldProfitRate != 2000 || a.NewProfitRate != 1724 {
		t.Errorf("got old=%d new=%d, want 2000/1724", a.OldProfitRate, a.NewProfitRate)
	}
	if a.Kind != KindRise {
		t.Errorf("kind: got %q, want rise", a.Kind)
	}
	if a.ID != "" || !a.CreatedAt.IsZero() {
		t.Errorf("Compute must not assign ID/timestamp: id=%q createdAt=%v", a.ID, a.CreatedAt)
	}
}

func TestConstantCapitalPriceEffectValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    ConstantCapitalPriceEffect
		want error
	}{
		{
			name: "valid",
			e:    ConstantCapitalPriceEffect{FixedCapital: 320, OriginalMaterialCapital: 80, PriceFactor: 200, VariableCapital: 100, SurplusValue: 100},
			want: nil,
		},
		{
			name: "negative magnitude",
			e:    ConstantCapitalPriceEffect{OriginalMaterialCapital: -1, PriceFactor: 200},
			want: ErrNegativeMagnitude,
		},
		{
			name: "zero price factor",
			e:    ConstantCapitalPriceEffect{FixedCapital: 320, OriginalMaterialCapital: 80, PriceFactor: 0},
			want: ErrNonPositivePriceFactor,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.e.Validate(); !errors.Is(err, tc.want) {
				t.Errorf("Validate() = %v, want %v", err, tc.want)
			}
		})
	}
}

// §II — a 3-month stock of 1,000 units a month. When the price falls from £10 to
// £7, the capital advanced to hold the stock is released: 3 × 1,000 × (10 − 7) =
// £9,000. When it rises from £7 to £10, the same magnitude must be tied up.
func TestCapitalReleaseAndTieUp(t *testing.T) {
	t.Parallel()

	rel := CapitalRelease{StockMonths: 3, StockUnits: 1000, OriginalPricePerUnit: 10, NewPricePerUnit: 7}
	if got := rel.Released(); got != 9000 {
		t.Errorf("released: got %d, want 9000", got)
	}

	tie := CapitalTieUp{StockMonths: 3, StockUnits: 1000, OriginalPricePerUnit: 7, NewPricePerUnit: 10}
	if got := tie.TiedUp(); got != 9000 {
		t.Errorf("tied up: got %d, want 9000", got)
	}

	// Invariant [§II]: for one and the same price movement, release and tie-up
	// are the same magnitude with opposite sign — they sum to zero.
	mvRelease := CapitalRelease{StockMonths: 3, StockUnits: 1000, OriginalPricePerUnit: 10, NewPricePerUnit: 7}
	mvTieUp := CapitalTieUp{StockMonths: 3, StockUnits: 1000, OriginalPricePerUnit: 10, NewPricePerUnit: 7}
	if mvRelease.Released()+mvTieUp.TiedUp() != 0 {
		t.Errorf("release %d + tie-up %d should be 0", mvRelease.Released(), mvTieUp.TiedUp())
	}
}

func TestRawMaterialPriceFluctuation(t *testing.T) {
	t.Parallel()

	rise := RawMaterialPriceFluctuation{OriginalPricePerUnit: 7, NewPricePerUnit: 10}
	if rise.Kind() != KindRise {
		t.Errorf("7→10 kind: got %q, want rise", rise.Kind())
	}
	if got := rise.Factor(); got != 142 { // 10 × 100 / 7 = 142 (truncated)
		t.Errorf("7→10 factor: got %d, want 142", got)
	}

	fall := RawMaterialPriceFluctuation{OriginalPricePerUnit: 10, NewPricePerUnit: 7}
	if fall.Kind() != KindFall {
		t.Errorf("10→7 kind: got %q, want fall", fall.Kind())
	}
	if got := fall.Factor(); got != 70 { // 7 × 100 / 10 = 70
		t.Errorf("10→7 factor: got %d, want 70", got)
	}

	if got := (RawMaterialPriceFluctuation{}).Factor(); got != 0 {
		t.Errorf("zero original price: factor got %d, want 0", got)
	}
}

func TestSupplyRegulationAction(t *testing.T) {
	t.Parallel()

	fall := SupplyRegulation{Fluctuation: RawMaterialPriceFluctuation{OriginalPricePerUnit: 10, NewPricePerUnit: 7}}
	if fall.Action() != ActionExpandOutput {
		t.Errorf("price fall: got %q, want expand_output", fall.Action())
	}

	rise := SupplyRegulation{Fluctuation: RawMaterialPriceFluctuation{OriginalPricePerUnit: 7, NewPricePerUnit: 10}}
	if rise.Action() != ActionContractOutput {
		t.Errorf("price rise: got %q, want contract_output", rise.Action())
	}
}

func TestPriceFluctuationKindValid(t *testing.T) {
	t.Parallel()

	for _, k := range []PriceFluctuationKind{KindRise, KindFall} {
		if !k.Valid() {
			t.Errorf("kind %q should be valid", k)
		}
	}
	if PriceFluctuationKind("sideways").Valid() {
		t.Error(`"sideways" should be invalid`)
	}
}
