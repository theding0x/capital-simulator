package profit

import (
	"errors"
	"testing"
)

// Marx's §I two-enterprise illustration: equal s = £1,000 on equal v = £1,000,
// the rate of profit set only by the constant capital each carries.
func TestConstantCapitalEconomyRates(t *testing.T) {
	t.Parallel()

	// Enterprise A: c = £9,000 → p' = 1000/10000 = 10% = 1000 bp.
	a := ConstantCapitalEconomy{Kind: KindRawMaterialSaving, ConstantCapital: 9000, VariableCapital: 1000, SurplusValue: 1000}
	if got := a.OldProfitRate(); got != 1000 {
		t.Fatalf("enterprise A old p': got %d bp, want 1000", got)
	}

	// Enterprise B: c = £11,000 → p' = 1000/12000 = 8 1/3% = 833 bp (truncated).
	b := ConstantCapitalEconomy{Kind: KindRawMaterialSaving, ConstantCapital: 11000, VariableCapital: 1000, SurplusValue: 1000}
	if got := b.OldProfitRate(); got != 833 {
		t.Fatalf("enterprise B old p': got %d bp, want 833", got)
	}
}

// Economising the £2,000 of constant capital that separates B from A lifts the
// rate of profit from 8 1/3% to 10% without touching s or v.
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

// §IV — a communal boiler-house shared by 12 factories. If each factory bearing
// its own boiler would cost £2,400 in aggregate and the shared house costs
// £1,200, each factory economises (2400 − 1200) / 12 = £100.
func TestSharedConditionSavingPerParticipant(t *testing.T) {
	t.Parallel()

	boiler := SharedCondition{Kind: KindSharedConditions, SoloCost: 2400, SharedCost: 1200, ParticipantCount: 12}
	if got := boiler.SavingPerParticipant(); got != 100 {
		t.Errorf("boiler saving/participant: got %d, want 100", got)
	}

	none := SharedCondition{Kind: KindSharedConditions, SoloCost: 2400, SharedCost: 1200, ParticipantCount: 0}
	if got := none.SavingPerParticipant(); got != 0 {
		t.Errorf("zero participants: got %d, want 0", got)
	}
}

// §IV — reclaiming 40 lb of the 100 lb of waste cotton at £5 a unit cheapens the
// constant capital by 40 × 5 = £200.
func TestWasteReductionConstantCapitalSaving(t *testing.T) {
	t.Parallel()

	waste := WasteReduction{Kind: KindWasteReduction, OriginalWaste: 100, RecoveredWaste: 40, PricePerUnit: 5}
	if got := waste.ConstantCapitalSaving(); got != 200 {
		t.Errorf("waste saving: got %d, want 200", got)
	}
}

// §I — £1,200 of fixed capital serves a 10-hour or a 16-hour day alike. Extending
// the day spreads its charge from £120/hour to £75/hour: a per-hour economy of £45.
func TestWorkdayExtensionEffect(t *testing.T) {
	t.Parallel()

	ext := WorkdayExtensionEffect{Kind: KindWorkdayExtension, FixedCapital: 1200, OriginalHours: 10, ExtendedHours: 16}
	if got := ext.FixedChargePerHour(10); got != 120 {
		t.Errorf("charge at 10h: got %d, want 120", got)
	}
	if got := ext.FixedChargePerHour(16); got != 75 {
		t.Errorf("charge at 16h: got %d, want 75", got)
	}
	if got := ext.ProfitRateGain(); got != 45 {
		t.Errorf("gain: got %d, want 45", got)
	}
	if got := ext.ProfitRateGain(); got <= 0 {
		t.Errorf("gain must be positive when fixed capital > 0, got %d", got)
	}

	// No fixed capital to spread → no gain, however long the day.
	noFixed := WorkdayExtensionEffect{Kind: KindWorkdayExtension, FixedCapital: 0, OriginalHours: 10, ExtendedHours: 16}
	if got := noFixed.ProfitRateGain(); got != 0 {
		t.Errorf("no fixed capital: gain got %d, want 0", got)
	}
}
