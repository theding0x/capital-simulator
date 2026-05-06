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
	got := commodity.SurplusProduceRatio(360, 360)
	if math.Abs(got-0.5) > 1e-9 {
		t.Errorf("SurplusProduceRatio(360,360): want 0.5, got %v", got)
	}
}

func TestSurplusProduceRatio_ZeroInputs(t *testing.T) {
	t.Parallel()
	if got := commodity.SurplusProduceRatio(0, 0); got != 0 {
		t.Errorf("SurplusProduceRatio(0,0): want 0, got %v", got)
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
