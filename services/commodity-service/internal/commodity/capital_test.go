package commodity_test

import (
	"math"
	"testing"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)

func TestWearFractionFor_Instrument(t *testing.T) {
	t.Parallel()
	c := commodity.ConstantCapital{
		OriginalValue:   600,
		Kind:            commodity.KindInstrument,
		ServiceLifeDays: 6,
	}
	got := commodity.WearFractionFor(c)
	want := commodity.WearFraction(1.0 / 6.0)
	if math.Abs(float64(got-want)) > 1e-9 {
		t.Errorf("WearFractionFor(6-day instrument): want %v, got %v", want, got)
	}
}

func TestWearFractionFor_RawMaterial(t *testing.T) {
	t.Parallel()
	c := commodity.ConstantCapital{
		OriginalValue:   120,
		Kind:            commodity.KindRawMaterial,
		ServiceLifeDays: 0,
	}
	got := commodity.WearFractionFor(c)
	if float64(got) != 1.0 {
		t.Errorf("WearFractionFor(raw_material): want 1.0, got %v", got)
	}
}

func TestTransferredValue_Machine(t *testing.T) {
	t.Parallel()
	machine := commodity.ConstantCapital{
		OriginalValue:   100_000,
		Kind:            commodity.KindInstrument,
		ServiceLifeDays: 1000,
	}
	fraction := commodity.WearFractionFor(machine)
	got := commodity.TransferredValue(machine, fraction)
	const want = commodity.LabourMinutes(100)
	if got != want {
		t.Errorf("TransferredValue(1000-day machine): want %d, got %d", want, got)
	}
}

func TestTransferredValue_NeverExceedsOriginal(t *testing.T) {
	t.Parallel()
	c := commodity.ConstantCapital{
		OriginalValue:   360,
		Kind:            commodity.KindRawMaterial,
		ServiceLifeDays: 0,
	}
	fraction := commodity.WearFractionFor(c)
	got := commodity.TransferredValue(c, fraction)
	if got > c.OriginalValue {
		t.Errorf("TransferredValue(%d) > OriginalValue(%d)", got, c.OriginalValue)
	}
}

func TestNewValue_Identity(t *testing.T) {
	t.Parallel()
	for _, d := range []commodity.LabourMinutes{0, 60, 360, 720} {
		if got := commodity.NewValue(d); got != d {
			t.Errorf("NewValue(%d): want %d, got %d", d, d, got)
		}
	}
}

func TestDecomposeProductValue_CottonSpinning(t *testing.T) {
	t.Parallel()
	constants := []commodity.ConstantCapital{
		{OriginalValue: 1200, Kind: commodity.KindRawMaterial, ServiceLifeDays: 0},
		{OriginalValue: 240, Kind: commodity.KindInstrument, ServiceLifeDays: 1},
	}
	v := commodity.VariableCapital{WageValue: 360, WorkingDay: 720}
	pv := commodity.DecomposeProductValue(constants, v)
	if pv.Constant != 1440 {
		t.Errorf("Constant: want 1440, got %d", pv.Constant)
	}
	if pv.Variable != 360 {
		t.Errorf("Variable: want 360, got %d", pv.Variable)
	}
	if pv.Surplus != 360 {
		t.Errorf("Surplus: want 360, got %d", pv.Surplus)
	}
	if pv.Total() != 2160 {
		t.Errorf("Total: want 2160, got %d", pv.Total())
	}
}

func TestProductValue_TotalInvariant(t *testing.T) {
	t.Parallel()
	pv := commodity.ProductValue{Constant: 500, Variable: 90, Surplus: 90}
	want := commodity.LabourMinutes(680)
	if pv.Total() != want {
		t.Errorf("Total invariant: want %d, got %d", want, pv.Total())
	}
}

func TestCapitalComposition_Ratio(t *testing.T) {
	t.Parallel()
	cc := commodity.CapitalComposition{Constant: 410, Variable: 90}
	got := cc.Ratio()
	want := 410.0 / 90.0
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("Ratio: want %v, got %v", want, got)
	}
}

func TestCapitalComposition_RatioZeroVariable(t *testing.T) {
	t.Parallel()
	cc := commodity.CapitalComposition{Constant: 100, Variable: 0}
	if got := cc.Ratio(); got != 0 {
		t.Errorf("Ratio with zero variable: want 0, got %v", got)
	}
}

func TestVariableCapital_Validate_WageExceedsDay(t *testing.T) {
	t.Parallel()
	v := commodity.VariableCapital{WageValue: 720, WorkingDay: 360}
	if err := v.Validate(); err == nil {
		t.Error("expected error when wage_value > working_day")
	}
}
