package profit

import (
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
)

// Provenance of the fixtures below (Ch. 5 of Vol. III). The two enterprises
// (s = v = £1,000, c = £9,000 or £11,000) are Ch. 7's illustration, borrowed
// because it isolates c. Ch. 5's own figures are the 12, 16 and 18 hour working
// days, the single shaft that saves 60 per cent (in some cases 80) of the
// shafting, woollen rags at £4 4s. then £44 a ton, and the limit c = 0.
// Magnitudes are LabourMinutes-typed and read as £1 per unit (Part I convention).

// Ch. 7's two-enterprise illustration (borrowed by Ch. 5): equal s = £1,000 on
// equal v = £1,000, the rate of profit set only by the constant capital each
// carries.
func TestConstantCapitalEconomyRates(t *testing.T) {
	t.Parallel()

	// Enterprise A: c = £9,000 → p' = 1000/10000 = 10% = 1000 bp.
	a := ConstantCapitalEconomy{Kind: KindRawMaterialSaving, ConstantCapital: 9000, VariableCapital: 1000, SurplusValue: 1000}
	if got := a.OldProfitRate(); got != 1000 {
		t.Fatalf("enterprise A old p': got %d bp, want 1000", got)
	}

	// Enterprise B: c = £11,000 → p' = 1000/12000 = 8 1/3% = 833.33 bp, rounded half-up.
	b := ConstantCapitalEconomy{Kind: KindRawMaterialSaving, ConstantCapital: 11000, VariableCapital: 1000, SurplusValue: 1000}
	if got := b.OldProfitRate(); got != 833 {
		t.Fatalf("enterprise B old p': got %d bp, want 833", got)
	}
}

// Ch. 7's pair again: economising the £2,000 of constant capital that separates B
// from A lifts the rate of profit from 8 1/3% to 10% without touching s or v.
func TestComputeEconomyAnalysisTwoEnterprises(t *testing.T) {
	t.Parallel()

	e := ConstantCapitalEconomy{
		Kind:            KindRawMaterialSaving,
		ConstantCapital: 11000,
		VariableCapital: 1000,
		SurplusValue:    1000,
		Saving:          2000,
	}
	a := ComputeEconomyAnalysis(e)
	if a.OldProfitRate != 833 {
		t.Errorf("old p': got %d bp, want 833", a.OldProfitRate)
	}
	if a.NewProfitRate != 1000 {
		t.Errorf("new p': got %d bp, want 1000", a.NewProfitRate)
	}
	if a.ProfitRateGain != 167 {
		t.Errorf("gain: got %d bp, want 167", a.ProfitRateGain)
	}
	// Invariant: a positive saving (s, v unchanged) raises the rate of profit.
	if a.NewProfitRate <= a.OldProfitRate {
		t.Errorf("new p' (%d) must exceed old p' (%d) when saving > 0", a.NewProfitRate, a.OldProfitRate)
	}
	if a.ID != "" || !a.CreatedAt.IsZero() {
		t.Errorf("ComputeEconomyAnalysis must not assign ID/timestamp: id=%q createdAt=%v", a.ID, a.CreatedAt)
	}
}

// With nothing economised the rate is unchanged; with nothing advanced it is zero.
// Illustrative boundary cases (no Marx figure); 800c + 200v + 200s is a free pick.
func TestConstantCapitalEconomyBoundaries(t *testing.T) {
	t.Parallel()

	noSaving := ConstantCapitalEconomy{Kind: KindWorkdayExtension, ConstantCapital: 800, VariableCapital: 200, SurplusValue: 200}
	if got := noSaving.ProfitRateGain(); got != 0 {
		t.Errorf("no saving: gain got %d bp, want 0", got)
	}

	empty := ConstantCapitalEconomy{Kind: KindWorkdayExtension}
	if got := empty.OldProfitRate(); got != 0 {
		t.Errorf("zero capital: old p' got %d bp, want 0", got)
	}
	if got := empty.NewProfitRate(); got != 0 {
		t.Errorf("zero capital: new p' got %d bp, want 0", got)
	}
}

func TestConstantCapitalEconomyValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		e    ConstantCapitalEconomy
		want error
	}{
		{
			name: "valid",
			e:    ConstantCapitalEconomy{Kind: KindRawMaterialSaving, ConstantCapital: 11000, VariableCapital: 1000, SurplusValue: 1000, Saving: 2000},
			want: nil,
		},
		{
			name: "negative magnitude",
			e:    ConstantCapitalEconomy{Kind: KindRawMaterialSaving, ConstantCapital: -1},
			want: ErrNegativeMagnitude,
		},
		{
			name: "unknown kind",
			e:    ConstantCapitalEconomy{Kind: "bogus", ConstantCapital: 100, VariableCapital: 50, SurplusValue: 50},
			want: ErrUnknownEconomyKind,
		},
		{
			name: "saving exceeds constant",
			e:    ConstantCapitalEconomy{Kind: KindRawMaterialSaving, ConstantCapital: 100, VariableCapital: 50, SurplusValue: 50, Saving: 200},
			want: ErrSavingExceedsConstant,
		},
		{
			// Marx's limit c = 0 (Ch. 5 §I): the whole constant capital economised.
			name: "saving equals constant",
			e:    ConstantCapitalEconomy{Kind: KindRawMaterialSaving, ConstantCapital: 100, VariableCapital: 50, SurplusValue: 50, Saving: 100},
			want: nil,
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

func TestEconomyKindValid(t *testing.T) {
	t.Parallel()

	for _, k := range []EconomyKind{
		KindSharedConditions, KindWasteReduction, KindWorkdayExtension,
		KindImprovedMachinery, KindRawMaterialSaving,
	} {
		if !k.Valid() {
			t.Errorf("kind %q should be valid", k)
		}
	}
	if EconomyKind("not-a-kind").Valid() {
		t.Error(`"not-a-kind" should be invalid`)
	}
}

// §III — the single shaft. Marx's cotton spinner drives 29,000 doubling spindles
// by one shaft, a saving in shafting of 60 per cent against other firms, in some
// cases 80. The £500 of shafting the others need is illustrative (Marx gives no
// sum); the percentages are his, so the saving must be exactly 60 (80) per cent.
func TestSharedConditionSavingPerParticipant(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		shared      SharedCondition
		want        ConstantCapital
		wantPerCent int64 // Marx's per cent of SoloCost saved; 0 where he gives none
	}{
		{
			name:        "single shaft saves 60 per cent",
			shared:      SharedCondition{Kind: KindSharedConditions, SoloCost: 500, SharedCost: 200, ParticipantCount: 1},
			want:        300,
			wantPerCent: 60,
		},
		{
			name:        "in some cases 80 per cent",
			shared:      SharedCondition{Kind: KindSharedConditions, SoloCost: 500, SharedCost: 100, ParticipantCount: 1},
			want:        400,
			wantPerCent: 80,
		},
		{
			// I: the division by participants, boundary only (no Marx figure).
			name:   "illustrative division over twelve",
			shared: SharedCondition{Kind: KindSharedConditions, SoloCost: 2400, SharedCost: 1200, ParticipantCount: 12},
			want:   100,
		},
		{
			name:   "zero participants",
			shared: SharedCondition{Kind: KindSharedConditions, SoloCost: 2400, SharedCost: 1200, ParticipantCount: 0},
			want:   0,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.shared.SavingPerParticipant()
			if got != tc.want {
				t.Fatalf("saving/participant: got %d, want %d", got, tc.want)
			}
			if tc.wantPerCent != 0 && int64(got)*100 != tc.wantPerCent*int64(tc.shared.SoloCost) {
				t.Errorf("saving %d of %d is not %d per cent", got, tc.shared.SoloCost, tc.wantPerCent)
			}
		})
	}
}

// §IV — woollen rags, which Marx prices per ton: about £4 4s. thirty years ago,
// £44 now. The 10 tons reclaimed are illustrative (the price is his), so
// reclaiming all 10 cheapens the constant capital by 10 × £44 = £440.
func TestWasteReductionConstantCapitalSaving(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		waste WasteReduction
		want  ConstantCapital
	}{
		{"all ten tons reclaimed", WasteReduction{Kind: KindWasteReduction, OriginalWaste: 10, RecoveredWaste: 10, PricePerUnit: 44}, 440},
		// Partial recovery: the saving follows RecoveredWaste, not OriginalWaste.
		{"four of ten tons reclaimed", WasteReduction{Kind: KindWasteReduction, OriginalWaste: 10, RecoveredWaste: 4, PricePerUnit: 44}, 176},
		{"nothing reclaimed", WasteReduction{Kind: KindWasteReduction, OriginalWaste: 10, RecoveredWaste: 0, PricePerUnit: 44}, 0},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.waste.ConstantCapitalSaving(); got != tc.want {
				t.Errorf("waste saving: got %d, want %d", got, tc.want)
			}
		})
	}
}

// §IV at the MELT. Marx's two prices of a ton of woollen rags, £4 4s. and £44, are
// exact at pkg/money's scale (£1 = 2,400 LM) though not at the whole-£ scale of the
// Part I seeds. The price rose by a factor between ten and eleven, by £39 16s.
func TestWasteReduction_RagsAtTheMELT(t *testing.T) {
	t.Parallel()

	thenPrice := LabourMinutes(money.FromLSD(4, 4, 0).ToLabourMinutes())
	nowPrice := LabourMinutes(money.FromLSD(44, 0, 0).ToLabourMinutes())
	then := WasteReduction{Kind: KindWasteReduction, OriginalWaste: 1, RecoveredWaste: 1, PricePerUnit: thenPrice}
	now := WasteReduction{Kind: KindWasteReduction, OriginalWaste: 1, RecoveredWaste: 1, PricePerUnit: nowPrice}

	// £4 4s. = 1,008d = 10,080 LM; £44 = 10,560d = 105,600 LM.
	thenSaving, nowSaving := then.ConstantCapitalSaving(), now.ConstantCapitalSaving()
	if thenSaving != 10080 {
		t.Errorf("a ton at £4 4s.: got %d LM, want 10080", thenSaving)
	}
	if nowSaving != 105600 {
		t.Errorf("a ton at £44: got %d LM, want 105600", nowSaving)
	}
	if !(10*thenSaving < nowSaving && nowSaving < 11*thenSaving) {
		t.Errorf("the rag price should rise between ten and eleven fold: %d to %d", thenSaving, nowSaving)
	}
	rise := money.FromLabourMinutes(int64(nowSaving - thenSaving))
	if rise != money.FromLSD(39, 16, 0) {
		t.Errorf("rise: got %d d, want £39 16s. (%d d)", int64(rise), int64(money.FromLSD(39, 16, 0)))
	}
}

// The working-day economy's base (I, Marx gives no sum): 3,600 of fixed capital,
// all of c, over v = s = 1,200. 3,600 is chosen so that 12, 16 and 18 divide it.
const (
	ch5WorkdayFixed = ConstantCapital(3600)
	ch5WorkdayV     = VariableCapital(1200)
	ch5WorkdayS     = SurplusValue(1200)
)

// §I — Marx's working-days: the fixed capital is the same whether the
// labour-process serves 16 or 12 hours, and at 18 instead of 12 its value is
// reproduced in two thirds of the time. The per-hour charges are LabourMinutes
// per hour, not £ and not a rate.
func TestWorkdayExtensionEffect(t *testing.T) {
	t.Parallel()

	const fixed = ch5WorkdayFixed
	tests := []struct {
		name                     string
		original, extended       int64
		wantOriginal, wantNew    ConstantCapital
		wantGain                 ConstantCapital
		shareNumer, shareDenomin int64 // the extended charge is shareNumer/shareDenomin of the original
	}{
		{"12 to 18 hours: two thirds", 12, 18, 300, 200, 100, 2, 3},
		{"12 to 16 hours: three quarters", 12, 16, 300, 225, 75, 3, 4},
		{"12 to 12 hours: nothing gained", 12, 12, 300, 300, 0, 1, 1},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ext := WorkdayExtensionEffect{Kind: KindWorkdayExtension, FixedCapital: fixed, OriginalHours: tc.original, ExtendedHours: tc.extended}
			if got := ext.FixedChargePerHour(tc.original); got != tc.wantOriginal {
				t.Errorf("charge at %dh: got %d, want %d", tc.original, got, tc.wantOriginal)
			}
			if got := ext.FixedChargePerHour(tc.extended); got != tc.wantNew {
				t.Errorf("charge at %dh: got %d, want %d", tc.extended, got, tc.wantNew)
			}
			if got := ext.ProfitRateGain(); got != tc.wantGain {
				t.Errorf("per-hour saving: got %d LM/hour, want %d", got, tc.wantGain)
			}
			if ext.FixedChargePerHour(tc.extended)*ConstantCapital(tc.shareDenomin) != ext.FixedChargePerHour(tc.original)*ConstantCapital(tc.shareNumer) {
				t.Errorf("extended charge is not %d/%d of the original", tc.shareNumer, tc.shareDenomin)
			}
			// Invariant: lengthening the day over fixed capital > 0 always saves a charge.
			if tc.extended > tc.original && ext.ProfitRateGain() <= 0 {
				t.Errorf("a lengthened day must save a positive charge, got %d", ext.ProfitRateGain())
			}
		})
	}

	// No fixed capital to spread: no gain however long the day.
	noFixed := WorkdayExtensionEffect{Kind: KindWorkdayExtension, FixedCapital: 0, OriginalHours: 12, ExtendedHours: 18}
	if got := noFixed.ProfitRateGain(); got != 0 {
		t.Errorf("no fixed capital: gain got %d, want 0", got)
	}
	// A shortened day reverses the sign: the method is a signed per-hour charge saving.
	// This documents what the method does today, not a requirement of Marx's text.
	shortened := WorkdayExtensionEffect{Kind: KindWorkdayExtension, FixedCapital: fixed, OriginalHours: 18, ExtendedHours: 12}
	if got := shortened.ProfitRateGain(); got != -100 {
		t.Errorf("shortened day 18 to 12: got %d, want -100", got)
	}
}

// §I as an economy in c. Spread over 18 hours instead of 12, the fixed charge each
// hour bears falls from 300 to 200, so the same 12 hours of labour carry 2,400
// instead of 3,600: a saving of a third of the fixed capital (Marx's two thirds
// of the time, seen as a charge). That saving, set against illustrative v = s =
// 1,200 with all of c taken as fixed (the upper bound), lifts p' from 25% to
// 33.33% without a farthing of new surplus-value.
func TestWorkdayExtensionEconomy_RaisesRateOfProfit(t *testing.T) {
	t.Parallel()

	ext := WorkdayExtensionEffect{Kind: KindWorkdayExtension, FixedCapital: ch5WorkdayFixed, OriginalHours: 12, ExtendedHours: 18}
	saving := ConstantCapital(ext.OriginalHours) * (ext.FixedChargePerHour(ext.OriginalHours) - ext.FixedChargePerHour(ext.ExtendedHours))
	if saving != 1200 || saving*3 != ext.FixedCapital {
		t.Fatalf("saving on the same 12 hours: got %d, want 1200 = fixed capital / 3", saving)
	}

	e := ConstantCapitalEconomy{
		Kind: KindWorkdayExtension, ConstantCapital: ext.FixedCapital,
		VariableCapital: ch5WorkdayV, SurplusValue: ch5WorkdayS, Saving: saving,
	}
	if err := e.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
	a := ComputeEconomyAnalysis(e)
	if a.OldProfitRate != 2500 || a.NewProfitRate != 3333 || a.ProfitRateGain != 833 {
		t.Errorf("got old=%d new=%d gain=%d bp, want 2500/3333/833", a.OldProfitRate, a.NewProfitRate, a.ProfitRateGain)
	}
}

// Marx's limit (Ch. 5 §I): with the value of c at zero, p' equals s' and the rate
// of profit is at its maximum. Economising all of c leaves p' equal to the rate of
// surplus-value, and no smaller saving can pass it. The 3,600, 1,200 and 1,200
// are illustrative; the limit is Marx's.
func TestConstantCapitalEconomy_NoConstantCapitalLeavesRateOfSurplusValue(t *testing.T) {
	t.Parallel()

	const c, v, s = ch5WorkdayFixed, ch5WorkdayV, ch5WorkdayS
	ceiling := int64(ComputeRateOfSurplusValue(s, v))

	whole := ConstantCapitalEconomy{Kind: KindWorkdayExtension, ConstantCapital: c, VariableCapital: v, SurplusValue: s, Saving: c}
	if err := whole.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
	if got := whole.NewProfitRate(); got != ceiling {
		t.Errorf("c = 0: p' got %d bp, want s' = %d", got, ceiling)
	}

	prev := int64(0)
	for saving := ConstantCapital(0); saving <= c; saving += 400 {
		e := ConstantCapitalEconomy{Kind: KindWorkdayExtension, ConstantCapital: c, VariableCapital: v, SurplusValue: s, Saving: saving}
		got := e.NewProfitRate()
		if got < prev || got > ceiling {
			t.Errorf("saving %d: p' = %d bp, want within [%d, %d]", saving, got, prev, ceiling)
		}
		prev = got
	}
}

// The two enterprises are Ch. 7's OrganicCompositionEffect pair; Ch. 5's economy
// of £2,000 turns B into A. Pinning one to the other stops them drifting apart.
func TestTwoEnterprises_AreCh7sOrganicCompositionPair(t *testing.T) {
	t.Parallel()

	pair := OrganicCompositionEffect{SA: 1000, VA: 1000, CA: 9000, SB: 1000, VB: 1000, CB: 11000}
	e := ConstantCapitalEconomy{
		Kind: KindRawMaterialSaving, ConstantCapital: pair.CB,
		VariableCapital: pair.VB, SurplusValue: pair.SB, Saving: pair.CB - pair.CA,
	}
	if got, want := e.OldProfitRate(), pair.ProfitRateB(); got != want {
		t.Errorf("old p': got %d, want Ch. 7's enterprise B %d", got, want)
	}
	if got, want := e.NewProfitRate(), pair.ProfitRateA(); got != want {
		t.Errorf("new p': got %d, want Ch. 7's enterprise A %d", got, want)
	}
	if got, want := e.ProfitRateGain(), pair.RateDifference(); got != want {
		t.Errorf("gain: got %d, want Ch. 7's rate difference %d", got, want)
	}
}
