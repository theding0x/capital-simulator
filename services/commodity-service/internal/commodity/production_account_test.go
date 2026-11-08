package commodity_test

import (
	"errors"
	"math"
	"testing"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)

func TestComputeRate_HundredPercent(t *testing.T) {
	t.Parallel()
	rate, err := commodity.ComputeRate(commodity.SurplusValue(90), commodity.LabourMinutes(90))
	if err != nil {
		t.Fatalf("ComputeRate: %v", err)
	}
	if math.Abs(float64(rate)-1.0) > 1e-9 {
		t.Errorf("rate: want 1.0, got %v", rate)
	}
}

func TestComputeRate_NotTotalCapital(t *testing.T) {
	t.Parallel()
	rateOnTotal, _ := commodity.ComputeRate(commodity.SurplusValue(90), commodity.LabourMinutes(500))
	rateOnVar, _ := commodity.ComputeRate(commodity.SurplusValue(90), commodity.LabourMinutes(90))
	if math.Abs(float64(rateOnTotal)-0.18) > 0.01 {
		t.Errorf("s/C: want ~0.18, got %v", rateOnTotal)
	}
	if math.Abs(float64(rateOnVar)-1.0) > 1e-9 {
		t.Errorf("s/v: want 1.0, got %v", rateOnVar)
	}
}

func TestComputeRate_SpinningMill1871(t *testing.T) {
	t.Parallel()
	rate, err := commodity.ComputeRate(commodity.SurplusValue(80), commodity.LabourMinutes(52))
	if err != nil {
		t.Fatalf("ComputeRate: %v", err)
	}
	want := 80.0 / 52.0
	if math.Abs(float64(rate)-want) > 1e-9 {
		t.Errorf("spinning mill rate: want %v, got %v", want, rate)
	}
}

func TestComputeRate_ZeroVariable(t *testing.T) {
	t.Parallel()
	_, err := commodity.ComputeRate(commodity.SurplusValue(90), commodity.LabourMinutes(0))
	if !errors.Is(err, commodity.ErrDivisionByZero) {
		t.Errorf("want ErrDivisionByZero, got %v", err)
	}
}

func TestComputeRate_EqualsLabourTimeRatio(t *testing.T) {
	t.Parallel()
	nl, sl := commodity.LabourMinutes(360), commodity.LabourMinutes(360)
	rateByValue, _ := commodity.ComputeRate(commodity.SurplusValue(sl), nl)
	rateByTime := float64(sl) / float64(nl)
	if math.Abs(float64(rateByValue)-rateByTime) > 1e-9 {
		t.Errorf("rate mismatch: by-value=%v, by-time=%v", rateByValue, rateByTime)
	}
}

func TestProductionAccount_ValueProduct(t *testing.T) {
	t.Parallel()
	a := commodity.ProductionAccount{Constant: 410, Variable: 90, Surplus: commodity.SurplusValue(90)}
	got := a.ValueProduct()
	want := commodity.ValueProduct(180)
	if got != want {
		t.Errorf("ValueProduct: want %d, got %d", want, got)
	}
}

func TestProductionAccount_ExpandedCapital(t *testing.T) {
	t.Parallel()
	a := commodity.ProductionAccount{Constant: 410, Variable: 90, Surplus: commodity.SurplusValue(90)}
	got := a.ExpandedCapital()
	want := commodity.LabourMinutes(590)
	if got != want {
		t.Errorf("ExpandedCapital: want %d, got %d", want, got)
	}
}

func TestSurplusProduceRatio(t *testing.T) {
	t.Parallel()
	// Marx Ch.9 §4: 2 lbs surplus-produce : 2 lbs necessary-produce = 100%.
	got := commodity.SurplusProduceRatio(360, 360)
	if math.Abs(got-1.0) > 1e-9 {
		t.Errorf("SurplusProduceRatio(360,360): want 1.0, got %v", got)
	}
}

func TestSurplusProduceRatio_EqualsRate(t *testing.T) {
	t.Parallel()
	// 1871 Spinning Mill: v=90, s=90 -> s'=100%; ratio must match Rate().
	a := commodity.ProductionAccount{Constant: 410, Variable: 90, Surplus: 90}
	rate, err := a.Rate()
	if err != nil {
		t.Fatalf("Rate: %v", err)
	}
	got := commodity.SurplusProduceRatio(commodity.LabourMinutes(a.Surplus), a.Variable)
	if math.Abs(got-float64(rate)) > 1e-9 {
		t.Errorf("SurplusProduceRatio = %v, want Rate() = %v", got, float64(rate))
	}
}

func TestSurplusProduceRatio_ZeroNecessary(t *testing.T) {
	t.Parallel()
	if got := commodity.SurplusProduceRatio(0, 0); got != 0 {
		t.Errorf("SurplusProduceRatio(0,0): want 0, got %v", got)
	}
	if got := commodity.SurplusProduceRatio(10, 0); got != 0 {
		t.Errorf("SurplusProduceRatio(10,0): want 0, got %v", got)
	}
}

func TestNewProductionAccountID(t *testing.T) {
	t.Parallel()
	id := commodity.NewProductionAccountID()
	if id.IsZero() {
		t.Error("NewProductionAccountID should not be zero")
	}
	if len(string(id)) != 24 {
		t.Errorf("ProductionAccountID: want 24 chars, got %d", len(string(id)))
	}
}

func TestComputeRateBasisPoints_RoundsHalfUp(t *testing.T) {
	t.Parallel()

	// Spinning mill: s/v of 1/3 is 3333.33 bp, rounded to 3333; 2/3 is 6666.67, rounded to 6667.
	cases := []struct {
		s    commodity.SurplusValue
		v    commodity.LabourMinutes
		want int64
	}{
		{360, 360, 10000},
		{1080, 360, 30000},
		{1, 3, 3333},
		{2, 3, 6667},
	}
	for _, tc := range cases {
		got, err := commodity.ComputeRateBasisPoints(tc.s, tc.v)
		if err != nil {
			t.Fatalf("ComputeRateBasisPoints(%d,%d): %v", tc.s, tc.v, err)
		}
		if got != tc.want {
			t.Errorf("ComputeRateBasisPoints(%d,%d) = %d, want %d", tc.s, tc.v, got, tc.want)
		}
	}
}

func TestComputeRateBasisPoints_ZeroVariable(t *testing.T) {
	t.Parallel()

	if _, err := commodity.ComputeRateBasisPoints(10, 0); !errors.Is(err, commodity.ErrDivisionByZero) {
		t.Errorf("err = %v, want commodity.ErrDivisionByZero", err)
	}
}

// Ch. 8 (and spec 08 invariant): the value newly created by living labour is
// v + s and is independent of the constant capital c, which only re-appears
// in the product. Marx's Ch. 9 account is 410c + 90v + 90s; holding v and s
// fixed while c varies (including c = 0 and the 1/10th figure) must leave the
// value product and the rate untouched and move only the expanded capital.
func TestProductionAccount_NewValueIndependentOfConstantCapital(t *testing.T) {
	t.Parallel()
	base := commodity.ProductionAccount{Constant: 410, Variable: 90, Surplus: 90}
	wantProduct := commodity.ValueProduct(180)
	wantRate, err := base.Rate()
	if err != nil {
		t.Fatalf("Rate: %v", err)
	}

	for _, c := range []commodity.LabourMinutes{0, 41, 410, 4100} {
		a := base
		a.Constant = c

		if got := a.ValueProduct(); got != wantProduct {
			t.Errorf("c=%d: ValueProduct = %d, want %d (v+s, independent of c)", c, got, wantProduct)
		}
		if got := commodity.ComputeValueProduct(a.Variable, a.Surplus); got != wantProduct {
			t.Errorf("c=%d: ComputeValueProduct = %d, want %d", c, got, wantProduct)
		}
		if got, err := a.Rate(); err != nil || got != wantRate {
			t.Errorf("c=%d: Rate = %v, %v; want %v (s/v ignores c)", c, got, err, wantRate)
		}
		if got, want := a.ExpandedCapital(), c+90+90; got != want {
			t.Errorf("c=%d: ExpandedCapital = %d, want c+v+s = %d", c, got, want)
		}
		// c enters the product but contributes no new value.
		if got := commodity.LabourMinutes(a.ExpandedCapital()) - commodity.LabourMinutes(a.ValueProduct()); got != c {
			t.Errorf("c=%d: product minus new value = %d, want c", c, got)
		}
	}
}
