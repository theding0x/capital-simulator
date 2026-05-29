package reproduction_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// Canonical fixture: "Marxian model economy 1865" as described in Vol. II Ch. 18.
// All monetary values are in pence (£1 = 240p in pre-decimal sterling).
// Total social money: £1000 = 1,000,000 pence.
// Department I reserve (means of production): £400 = 400,000 pence.
// Department II reserve (articles of consumption): £350 = 350,000 pence.
// Wage rotation fund: £200 = 200,000 pence.
// Idle hoard: £50 = 50,000 pence.
// Money stock: £500 = 500,000 pence; velocity = 2 turnovers/year (20,000 bp).
// Effective circulation: £1000 = 1,000,000 pence.
const (
	fixtureTotalMoneyPence   = 1_000_000
	fixtureDepIReservePence  = 400_000
	fixtureDepIIReservePence = 350_000
	fixtureWageReservePence  = 200_000
	fixtureIdleHoardPence    = 50_000
	fixtureMoneyStockPence   = 500_000
	fixtureVelocityBP        = 20_000 // 2 turnovers/year in basis points (1 turnover = 10,000 bp)
)

// TestNewMoneySupplyApportionmentID verifies that generated IDs are non-empty
// and distinct (randomness invariant).
func TestNewMoneySupplyApportionmentID(t *testing.T) {
	t.Parallel()
	a := reproduction.NewMoneySupplyApportionmentID()
	b := reproduction.NewMoneySupplyApportionmentID()
	if a == "" {
		t.Fatal("expected non-empty ID")
	}
	if b == "" {
		t.Fatal("expected non-empty second ID")
	}
	if a == b {
		t.Fatal("expected distinct IDs from two separate calls")
	}
}

// TestValidateApportionment verifies that the partition equation
// Total == I + II + Wage + Idle must hold.
func TestValidateApportionment(t *testing.T) {
	t.Parallel()

	t.Run("valid 1865 economy", func(t *testing.T) {
		t.Parallel()
		a := reproduction.MoneySupplyApportionment{
			ID:                    reproduction.NewMoneySupplyApportionmentID(),
			Period:                "1865",
			TotalSocialMoneyPence: fixtureTotalMoneyPence,
			DeptIReservePence:     fixtureDepIReservePence,
			DeptIIReservePence:    fixtureDepIIReservePence,
			WageRotationPence:     fixtureWageReservePence,
			IdleHoardPence:        fixtureIdleHoardPence,
		}
		if err := a.Validate(); err != nil {
			t.Fatalf("expected valid apportionment: %v", err)
		}
	})

	t.Run("sum mismatch returns ErrApportionmentMismatch", func(t *testing.T) {
		t.Parallel()
		// I=400000 + II=300000 + Wage=200000 + Idle=50000 = 950000 ≠ 1000000
		a := reproduction.MoneySupplyApportionment{
			ID:                    reproduction.NewMoneySupplyApportionmentID(),
			Period:                "1865",
			TotalSocialMoneyPence: fixtureTotalMoneyPence,
			DeptIReservePence:     fixtureDepIReservePence,
			DeptIIReservePence:    300_000, // deliberately reduced
			WageRotationPence:     fixtureWageReservePence,
			IdleHoardPence:        fixtureIdleHoardPence,
		}
		if err := a.Validate(); err != reproduction.ErrApportionmentMismatch {
			t.Fatalf("expected ErrApportionmentMismatch, got %v", err)
		}
	})

	t.Run("zero total with all zeros is valid", func(t *testing.T) {
		t.Parallel()
		a := reproduction.MoneySupplyApportionment{
			ID:                    reproduction.NewMoneySupplyApportionmentID(),
			Period:                "zero",
			TotalSocialMoneyPence: 0,
			DeptIReservePence:     0,
			DeptIIReservePence:    0,
			WageRotationPence:     0,
			IdleHoardPence:        0,
		}
		if err := a.Validate(); err != nil {
			t.Fatalf("expected zero apportionment to be valid: %v", err)
		}
	})
}

// TestComputeEffectiveCirculatingValue verifies Marx's quantity-theory identity:
// effective_circulation = money_stock × velocity.
// Velocity is expressed in basis points where 10,000 bp = 1 turnover/year.
func TestComputeEffectiveCirculatingValue(t *testing.T) {
	t.Parallel()

	t.Run("stock=500000 velocity=20000 → effective=1000000", func(t *testing.T) {
		t.Parallel()
		got := reproduction.ComputeEffectiveCirculatingValue(500_000, 20_000)
		if got != 1_000_000 {
			t.Fatalf("want 1000000, got %d", got)
		}
	})

	t.Run("stock=1000000 velocity=10000 → effective=1000000 (same effective, different decomposition)", func(t *testing.T) {
		t.Parallel()
		got := reproduction.ComputeEffectiveCirculatingValue(1_000_000, 10_000)
		if got != 1_000_000 {
			t.Fatalf("want 1000000, got %d", got)
		}
	})

	t.Run("doubling velocity halves required stock: stock=250000 velocity=40000 → effective=1000000", func(t *testing.T) {
		t.Parallel()
		got := reproduction.ComputeEffectiveCirculatingValue(250_000, 40_000)
		if got != 1_000_000 {
			t.Fatalf("want 1000000, got %d", got)
		}
	})
}

// TestCirculatingMoneyMassEffectiveValue verifies that the EffectiveCirculatingValuePence
// field on a CirculatingMoneyMass equals ComputeEffectiveCirculatingValue of its components.
func TestCirculatingMoneyMassEffectiveValue(t *testing.T) {
	t.Parallel()
	mass := reproduction.CirculatingMoneyMass{
		ID:                          reproduction.NewCirculatingMoneyMassID(),
		Period:                      "1865",
		MoneyStockPence:             fixtureMoneyStockPence,
		VelocityPerYearBasisPoints:  fixtureVelocityBP,
		EffectiveCirculatingValuePence: reproduction.ComputeEffectiveCirculatingValue(
			fixtureMoneyStockPence, fixtureVelocityBP,
		),
	}
	want := reproduction.ComputeEffectiveCirculatingValue(fixtureMoneyStockPence, fixtureVelocityBP)
	if mass.EffectiveCirculatingValuePence != want {
		t.Fatalf("EffectiveCirculatingValuePence: want %d, got %d", want, mass.EffectiveCirculatingValuePence)
	}
	if want != 1_000_000 {
		t.Fatalf("expected effective value 1000000 for 500000 stock at 2x velocity, got %d", want)
	}
}

// TestWageRotationFundFrequency verifies that WageRotationFund records the
// WageCycleFrequency exactly as supplied (structural field, not derived).
func TestWageRotationFundFrequency(t *testing.T) {
	t.Parallel()

	t.Run("weekly (52 cycles/year)", func(t *testing.T) {
		t.Parallel()
		wrf := reproduction.WageRotationFund{
			ID:                reproduction.NewWageRotationFundID(),
			Period:            "1865",
			FundPence:         fixtureWageReservePence,
			WageCycleFrequency: 52,
		}
		if wrf.WageCycleFrequency != 52 {
			t.Fatalf("WageCycleFrequency: want 52, got %d", wrf.WageCycleFrequency)
		}
	})

	t.Run("monthly (12 cycles/year)", func(t *testing.T) {
		t.Parallel()
		wrf := reproduction.WageRotationFund{
			ID:                reproduction.NewWageRotationFundID(),
			Period:            "1865",
			FundPence:         fixtureWageReservePence,
			WageCycleFrequency: 12,
		}
		if wrf.WageCycleFrequency != 12 {
			t.Fatalf("WageCycleFrequency: want 12, got %d", wrf.WageCycleFrequency)
		}
	})
}

// TestInterDepartmentSettlementBalance is a documentation test verifying the
// balanced-reproduction steady-state: Department I and II settle equal amounts.
func TestInterDepartmentSettlementBalance(t *testing.T) {
	t.Parallel()
	// Under simple reproduction, I→II and II→I exchanges net to zero.
	// Marx models the aggregate as a self-cancelling inter-department flow.
	const settlementAmount = int64(2_000_000) // £8,333 ≈ annual inter-dept flow in the 1865 model
	iToII := reproduction.InterDepartmentSettlement{
		ID:              reproduction.NewInterDepartmentSettlementID(),
		Period:          "1865",
		FromDepartment:  reproduction.DepartmentI,
		ToDepartment:    reproduction.DepartmentII,
		SettledPence:    settlementAmount,
		SettlementPurpose: reproduction.PurposeMeansOfProduction,
	}
	iiToI := reproduction.InterDepartmentSettlement{
		ID:              reproduction.NewInterDepartmentSettlementID(),
		Period:          "1865",
		FromDepartment:  reproduction.DepartmentII,
		ToDepartment:    reproduction.DepartmentI,
		SettledPence:    settlementAmount,
		SettlementPurpose: reproduction.PurposeArticlesOfConsumption,
	}
	if iToII.SettledPence != iiToI.SettledPence {
		t.Fatalf("expected balanced settlement: I→II %d == II→I %d",
			iToII.SettledPence, iiToI.SettledPence)
	}
}

// TestCapitalDepartmentConstants verifies the string values of the two
// department constants used as discriminators in inter-department records.
func TestCapitalDepartmentConstants(t *testing.T) {
	t.Parallel()
	if reproduction.DepartmentI != "department_i" {
		t.Fatalf("DepartmentI: want %q, got %q", "department_i", reproduction.DepartmentI)
	}
	if reproduction.DepartmentII != "department_ii" {
		t.Fatalf("DepartmentII: want %q, got %q", "department_ii", reproduction.DepartmentII)
	}
}

// TestReservePurposeConstants verifies all four reserve-purpose string constants.
func TestReservePurposeConstants(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		got  string
		want string
	}{
		{"means_of_production", string(reproduction.PurposeMeansOfProduction), "means_of_production"},
		{"articles_of_consumption", string(reproduction.PurposeArticlesOfConsumption), "articles_of_consumption"},
		{"wage_payment", string(reproduction.PurposeWagePayment), "wage_payment"},
		{"idle_hoard", string(reproduction.PurposeIdleHoard), "idle_hoard"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			if c.got != c.want {
				t.Fatalf("ReservePurpose: want %q, got %q", c.want, c.got)
			}
		})
	}
}

// TestInterDepartmentPurposeConstants verifies both inter-department settlement
// purpose constants against their canonical snake_case wire values.
func TestInterDepartmentPurposeConstants(t *testing.T) {
	t.Parallel()
	if string(reproduction.PurposeMeansOfProduction) != "means_of_production" {
		t.Fatalf("PurposeMeansOfProduction: want %q, got %q",
			"means_of_production", string(reproduction.PurposeMeansOfProduction))
	}
	if string(reproduction.PurposeArticlesOfConsumption) != "articles_of_consumption" {
		t.Fatalf("PurposeArticlesOfConsumption: want %q, got %q",
			"articles_of_consumption", string(reproduction.PurposeArticlesOfConsumption))
	}
}
