package commodity_test

import (
	"math"
	"testing"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)

// minutesPerPound converts Marx's £ to labour-minutes at the seed scale:
// 1s = 120 minutes, 20s = £1.
const minutesPerPound = 2400

// Ch.9 §1: "C = £410 const. + £90 var.; surplus-value = £90" — c:v ≈ 4.56.
func TestProductionAccount_Ch9Worked410(t *testing.T) {
	t.Parallel()
	a := commodity.ProductionAccount{
		Constant: 410 * minutesPerPound,
		Variable: 90 * minutesPerPound,
		Surplus:  commodity.SurplusValue(90 * minutesPerPound),
	}
	if got := float64(a.Constant) / float64(a.Variable); math.Abs(got-410.0/90.0) > 1e-9 {
		t.Errorf("c:v = %v, want %v", got, 410.0/90.0)
	}
	rate, err := a.Rate()
	if err != nil {
		t.Fatalf("Rate: %v", err)
	}
	if math.Abs(float64(rate)-1.0) > 1e-9 {
		t.Errorf("s' = %v, want 1.0", rate)
	}
}

// Ch.9 §1, spinning mill 1871: c = £378, v = £52, s = £80; s' = 153 11/13%.
func TestProductionAccount_SpinningMill1871(t *testing.T) {
	t.Parallel()
	a := commodity.ProductionAccount{
		Constant: 378 * minutesPerPound,
		Variable: 52 * minutesPerPound,
		Surplus:  commodity.SurplusValue(80 * minutesPerPound),
	}
	// Marx: price of 10,000 lbs yarn £510 = c + v + s = 378 + 52 + 80.
	total := a.Constant + a.Variable + commodity.LabourMinutes(a.Surplus)
	if want := commodity.LabourMinutes(510 * minutesPerPound); total != want {
		t.Errorf("c+v+s = %d, want %d (£510)", total, want)
	}
	// 153 11/13 % = 15384.6 bp, rounded half-up.
	bp, err := commodity.ComputeRateBasisPoints(a.Surplus, a.Variable)
	if err != nil {
		t.Fatalf("ComputeRateBasisPoints: %v", err)
	}
	if bp != 15385 {
		t.Errorf("s' = %d bp, want 15385 (153 11/13%%)", bp)
	}
}

// The same mill on a 10-hour day (600 min): necessary labour is
// v/(v+s) = 52/132 of the day, about 236 min, not the spec's 231.
func TestSpinningMill1871_TenHourDaySplit(t *testing.T) {
	t.Parallel()
	const day = 600.0
	v, s := 52.0, 80.0
	necessary := day * v / (v + s)
	surplus := day * s / (v + s)
	if math.Round(necessary) != 236 {
		t.Errorf("necessary = %.2f min, want ~236", necessary)
	}
	if math.Abs(necessary+surplus-day) > 1e-9 {
		t.Errorf("necessary+surplus = %v, want %v", necessary+surplus, day)
	}
}

// Ch.9 §1, Jacob's wheat (1815): wages £3 10s, surplus £3 11s; Marx puts
// seed and manure (£3 19s) at zero. s' = 71/70, "more than 100%".
func TestProductionAccount_JacobWheat1815(t *testing.T) {
	t.Parallel()
	const minutesPerShilling = 120
	a := commodity.ProductionAccount{
		Variable: 70 * minutesPerShilling,
		Surplus:  commodity.SurplusValue(71 * minutesPerShilling),
	}
	rate, err := a.Rate()
	if err != nil {
		t.Fatalf("Rate: %v", err)
	}
	if math.Abs(float64(rate)-71.0/70.0) > 1e-9 {
		t.Errorf("s' = %v, want 71/70", rate)
	}
	if rate <= 1 {
		t.Errorf("s' = %v, Marx: more than 100%%", rate)
	}
}
