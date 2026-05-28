package store_test

import (
	"context"
	"errors"
	"testing"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
)

// newCh18Store returns a fresh in-memory store for Ch. 18 tests.
// Each test that calls this gets an isolated state.
func newCh18Store(t *testing.T) *store.Memory {
	t.Helper()
	return store.NewMemory()
}

// --- MoneySupplyApportionment ---

// TestMemoryCreateGetMoneySupplyApportionment verifies the round-trip create/get
// for a MoneySupplyApportionment and checks ErrNotFound for missing IDs.
func TestMemoryCreateGetMoneySupplyApportionment(t *testing.T) {
	t.Parallel()
	mem := newCh18Store(t)
	ctx := context.Background()

	input := repro.MoneySupplyApportionment{
		ID:                    repro.NewMoneySupplyApportionmentID(),
		Period:                "1865",
		TotalSocialMoneyPence: 1_000_000,
		DeptIReservePence:     400_000,
		DeptIIReservePence:    350_000,
		WageRotationPence:     200_000,
		IdleHoardPence:        50_000,
	}

	created, err := mem.CreateMoneySupplyApportionment(ctx, input)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty ID after create")
	}

	got, err := mem.GetMoneySupplyApportionment(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Period != input.Period {
		t.Fatalf("Period: want %q, got %q", input.Period, got.Period)
	}
	if got.TotalSocialMoneyPence != input.TotalSocialMoneyPence {
		t.Fatalf("TotalSocialMoneyPence: want %d, got %d", input.TotalSocialMoneyPence, got.TotalSocialMoneyPence)
	}
	if got.DeptIReservePence != input.DeptIReservePence {
		t.Fatalf("DeptIReservePence: want %d, got %d", input.DeptIReservePence, got.DeptIReservePence)
	}
	if got.DeptIIReservePence != input.DeptIIReservePence {
		t.Fatalf("DeptIIReservePence: want %d, got %d", input.DeptIIReservePence, got.DeptIIReservePence)
	}
	if got.WageRotationPence != input.WageRotationPence {
		t.Fatalf("WageRotationPence: want %d, got %d", input.WageRotationPence, got.WageRotationPence)
	}
	if got.IdleHoardPence != input.IdleHoardPence {
		t.Fatalf("IdleHoardPence: want %d, got %d", input.IdleHoardPence, got.IdleHoardPence)
	}

	// Missing ID must return ErrNotFound.
	_, err = mem.GetMoneySupplyApportionment(ctx, "nonexistent-apportionment-id")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for missing ID, got %v", err)
	}
}

// TestMemoryListMoneySupplyApportionments verifies that list returns all created
// records and that an empty store returns a non-nil slice.
func TestMemoryListMoneySupplyApportionments(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("empty store returns non-nil slice", func(t *testing.T) {
		t.Parallel()
		mem := newCh18Store(t)
		list, err := mem.ListMoneySupplyApportionments(ctx)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if list == nil {
			t.Fatal("expected non-nil empty slice, got nil")
		}
	})

	t.Run("two records are both returned", func(t *testing.T) {
		t.Parallel()
		mem := newCh18Store(t)

		for i := 0; i < 2; i++ {
			_, err := mem.CreateMoneySupplyApportionment(ctx, repro.MoneySupplyApportionment{
				ID:                    repro.NewMoneySupplyApportionmentID(),
				Period:                "1865",
				TotalSocialMoneyPence: 1_000_000,
				DeptIReservePence:     400_000,
				DeptIIReservePence:    350_000,
				WageRotationPence:     200_000,
				IdleHoardPence:        50_000,
			})
			if err != nil {
				t.Fatalf("create %d: %v", i, err)
			}
		}

		list, err := mem.ListMoneySupplyApportionments(ctx)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(list) != 2 {
			t.Fatalf("expected 2 apportionments, got %d", len(list))
		}
	})
}

// --- DepartmentMoneyReserve ---

// TestMemoryCreateGetDepartmentMoneyReserve verifies round-trip for a
// DepartmentMoneyReserve and ErrNotFound for a missing ID.
func TestMemoryCreateGetDepartmentMoneyReserve(t *testing.T) {
	t.Parallel()
	mem := newCh18Store(t)
	ctx := context.Background()

	input := repro.DepartmentMoneyReserve{
		ID:             repro.NewDepartmentMoneyReserveID(),
		Period:         "1865",
		Department:     repro.DepartmentI,
		ReservePurpose: repro.PurposeMeansOfProduction,
		ReservePence:   400_000,
	}

	created, err := mem.CreateDepartmentMoneyReserve(ctx, input)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := mem.GetDepartmentMoneyReserve(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Department != repro.DepartmentI {
		t.Fatalf("Department: want %q, got %q", repro.DepartmentI, got.Department)
	}
	if got.ReservePence != input.ReservePence {
		t.Fatalf("ReservePence: want %d, got %d", input.ReservePence, got.ReservePence)
	}

	_, err = mem.GetDepartmentMoneyReserve(ctx, "missing-reserve-id")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// TestMemoryListDepartmentMoneyReserves verifies list length and non-nil empty slice.
func TestMemoryListDepartmentMoneyReserves(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("empty store returns non-nil slice", func(t *testing.T) {
		t.Parallel()
		mem := newCh18Store(t)
		list, err := mem.ListDepartmentMoneyReserves(ctx)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if list == nil {
			t.Fatal("expected non-nil empty slice, got nil")
		}
	})

	t.Run("two reserves are both returned", func(t *testing.T) {
		t.Parallel()
		mem := newCh18Store(t)

		departments := []struct {
			dept    repro.CapitalDepartment
			purpose repro.ReservePurpose
			pence   int64
		}{
			{repro.DepartmentI, repro.PurposeMeansOfProduction, 400_000},
			{repro.DepartmentII, repro.PurposeArticlesOfConsumption, 350_000},
		}
		for _, d := range departments {
			_, err := mem.CreateDepartmentMoneyReserve(ctx, repro.DepartmentMoneyReserve{
				ID:             repro.NewDepartmentMoneyReserveID(),
				Period:         "1865",
				Department:     d.dept,
				ReservePurpose: d.purpose,
				ReservePence:   d.pence,
			})
			if err != nil {
				t.Fatalf("create %s: %v", d.dept, err)
			}
		}

		list, err := mem.ListDepartmentMoneyReserves(ctx)
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if len(list) != 2 {
			t.Fatalf("expected 2 reserves, got %d", len(list))
		}
	})
}

// --- CirculatingMoneyMass ---

// TestMemoryCreateGetCirculatingMoneyMass verifies round-trip and ErrNotFound.
func TestMemoryCreateGetCirculatingMoneyMass(t *testing.T) {
	t.Parallel()
	mem := newCh18Store(t)
	ctx := context.Background()

	effective := repro.ComputeEffectiveCirculatingValue(500_000, 20_000) // → 1,000,000
	input := repro.CirculatingMoneyMass{
		ID:                             repro.NewCirculatingMoneyMassID(),
		Period:                         "1865",
		MoneyStockPence:                500_000,
		VelocityPerYearBasisPoints:     20_000,
		EffectiveCirculatingValuePence: effective,
	}

	created, err := mem.CreateCirculatingMoneyMass(ctx, input)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := mem.GetCirculatingMoneyMass(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.MoneyStockPence != 500_000 {
		t.Fatalf("MoneyStockPence: want 500000, got %d", got.MoneyStockPence)
	}
	if got.VelocityPerYearBasisPoints != 20_000 {
		t.Fatalf("VelocityPerYearBasisPoints: want 20000, got %d", got.VelocityPerYearBasisPoints)
	}
	if got.EffectiveCirculatingValuePence != 1_000_000 {
		t.Fatalf("EffectiveCirculatingValuePence: want 1000000, got %d", got.EffectiveCirculatingValuePence)
	}

	_, err = mem.GetCirculatingMoneyMass(ctx, "missing-mass-id")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// --- WageRotationFund ---

// TestMemoryCreateGetWageRotationFund verifies round-trip and ErrNotFound.
func TestMemoryCreateGetWageRotationFund(t *testing.T) {
	t.Parallel()
	mem := newCh18Store(t)
	ctx := context.Background()

	input := repro.WageRotationFund{
		ID:                 repro.NewWageRotationFundID(),
		Period:             "1865",
		FundPence:          200_000,
		WageCycleFrequency: 52, // weekly payment cycle
	}

	created, err := mem.CreateWageRotationFund(ctx, input)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := mem.GetWageRotationFund(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.FundPence != 200_000 {
		t.Fatalf("FundPence: want 200000, got %d", got.FundPence)
	}
	if got.WageCycleFrequency != 52 {
		t.Fatalf("WageCycleFrequency: want 52, got %d", got.WageCycleFrequency)
	}

	_, err = mem.GetWageRotationFund(ctx, "missing-fund-id")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// --- InterDepartmentSettlement ---

// TestMemoryCreateGetInterDepartmentSettlement verifies round-trip and ErrNotFound.
func TestMemoryCreateGetInterDepartmentSettlement(t *testing.T) {
	t.Parallel()
	mem := newCh18Store(t)
	ctx := context.Background()

	// Balanced steady reproduction: I→II means-of-production flow.
	input := repro.InterDepartmentSettlement{
		ID:                repro.NewInterDepartmentSettlementID(),
		Period:            "1865",
		FromDepartment:    repro.DepartmentI,
		ToDepartment:      repro.DepartmentII,
		SettledPence:      2_000_000,
		SettlementPurpose: repro.PurposeMeansOfProduction,
	}

	created, err := mem.CreateInterDepartmentSettlement(ctx, input)
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := mem.GetInterDepartmentSettlement(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.FromDepartment != repro.DepartmentI {
		t.Fatalf("FromDepartment: want %q, got %q", repro.DepartmentI, got.FromDepartment)
	}
	if got.ToDepartment != repro.DepartmentII {
		t.Fatalf("ToDepartment: want %q, got %q", repro.DepartmentII, got.ToDepartment)
	}
	if got.SettledPence != 2_000_000 {
		t.Fatalf("SettledPence: want 2000000, got %d", got.SettledPence)
	}

	_, err = mem.GetInterDepartmentSettlement(ctx, "missing-settlement-id")
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// --- CheckApportionmentBalance ---

// TestMemoryCheckApportionmentBalance verifies balanced/unbalanced detection and
// ErrNotFound for missing IDs.
func TestMemoryCheckApportionmentBalance(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("balanced apportionment returns true", func(t *testing.T) {
		t.Parallel()
		mem := newCh18Store(t)
		// 400000 + 350000 + 200000 + 50000 = 1000000 == Total
		a := repro.MoneySupplyApportionment{
			ID:                    repro.NewMoneySupplyApportionmentID(),
			Period:                "1865",
			TotalSocialMoneyPence: 1_000_000,
			DeptIReservePence:     400_000,
			DeptIIReservePence:    350_000,
			WageRotationPence:     200_000,
			IdleHoardPence:        50_000,
		}
		created, err := mem.CreateMoneySupplyApportionment(ctx, a)
		if err != nil {
			t.Fatalf("create: %v", err)
		}

		balanced, err := mem.CheckApportionmentBalance(ctx, created.ID)
		if err != nil {
			t.Fatalf("check: %v", err)
		}
		if !balanced {
			t.Fatal("expected balanced=true for a correctly partitioned apportionment")
		}
	})

	t.Run("all-zero apportionment is balanced", func(t *testing.T) {
		t.Parallel()
		mem := newCh18Store(t)
		// A degenerate but valid apportionment: all fields zero, total zero.
		a := repro.MoneySupplyApportionment{
			ID:                    repro.NewMoneySupplyApportionmentID(),
			Period:                "zero",
			TotalSocialMoneyPence: 0,
			DeptIReservePence:     0,
			DeptIIReservePence:    0,
			WageRotationPence:     0,
			IdleHoardPence:        0,
		}
		created, err := mem.CreateMoneySupplyApportionment(ctx, a)
		if err != nil {
			t.Fatalf("create zero apportionment: %v", err)
		}

		balanced, err := mem.CheckApportionmentBalance(ctx, created.ID)
		if err != nil {
			t.Fatalf("check: %v", err)
		}
		if !balanced {
			t.Fatal("expected balanced=true for all-zero apportionment")
		}
	})

	t.Run("missing ID returns ErrNotFound", func(t *testing.T) {
		t.Parallel()
		mem := newCh18Store(t)
		_, err := mem.CheckApportionmentBalance(ctx, "missing-id")
		if !errors.Is(err, store.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})
}
