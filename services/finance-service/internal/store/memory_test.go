package store

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

func TestMemory_CostPriceRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	cp := profit.ComputeCostPrice(profit.CapitalOutlay{
		Constant: 400, Variable: 100, FixedWearAndTear: 20, FixedAdvanced: 1200,
	})
	created, err := m.CreateCostPrice(ctx, cp)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}

	got, err := m.GetCostPrice(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetCostPriceNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetCostPrice(context.Background(), profit.CostPriceID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateCostPriceDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	cp := profit.ComputeCostPrice(profit.CapitalOutlay{Constant: 400, Variable: 100})
	cp.ID = profit.CostPriceID("fixed-id-1234")

	if _, err := m.CreateCostPrice(ctx, cp); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateCostPrice(ctx, cp); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListCostPricesNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListCostPrices(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

func TestMemory_ProfitRateRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	a := profit.ComputeProfitRateAnalysis(400, 100, 100)
	created, err := m.CreateProfitRate(ctx, a)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}

	got, err := m.GetProfitRate(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetProfitRateNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetProfitRate(context.Background(), profit.ProfitRateID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateProfitRateDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	a := profit.ComputeProfitRateAnalysis(400, 100, 100)
	a.ID = profit.ProfitRateID("fixed-id-5678")

	if _, err := m.CreateProfitRate(ctx, a); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateProfitRate(ctx, a); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListProfitRatesNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListProfitRates(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

func TestMemory_VariationRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	a := profit.ComputeVariationAnalysis(
		profit.ProfitRateFormula{C: 80, V: 20, SRate: 10000},
		profit.ProfitRateFormula{C: 100, V: 20, SRate: 10000},
	)
	created, err := m.CreateVariation(ctx, a)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}

	got, err := m.GetVariation(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetVariationNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetVariation(context.Background(), profit.VariationAnalysisID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateVariationDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	a := profit.ComputeVariationAnalysis(
		profit.ProfitRateFormula{C: 80, V: 20, SRate: 10000},
		profit.ProfitRateFormula{C: 60, V: 20, SRate: 10000},
	)
	a.ID = profit.VariationAnalysisID("fixed-id-9012")

	if _, err := m.CreateVariation(ctx, a); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateVariation(ctx, a); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListVariationsNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListVariations(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

func TestMemory_TurnoverAnalysisRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Marx's Capital A: 100C, 20v, s' = 100%, n = 2 → annual p' = 40%.
	a := profit.ComputeTurnoverAnalysis(100, 20, 10000, 2)
	created, err := m.CreateTurnoverAnalysis(ctx, a)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}

	got, err := m.GetTurnoverAnalysis(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetTurnoverAnalysisNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetTurnoverAnalysis(context.Background(), profit.TurnoverAnalysisID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateTurnoverAnalysisDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	a := profit.ComputeTurnoverAnalysis(100, 20, 10000, 2)
	a.ID = profit.TurnoverAnalysisID("fixed-id-3456")

	if _, err := m.CreateTurnoverAnalysis(ctx, a); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateTurnoverAnalysis(ctx, a); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListTurnoverAnalysesNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListTurnoverAnalyses(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

func TestMemory_EconomyAnalysisRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Marx's §I two-enterprise illustration: economising £2,000 of constant
	// capital lifts p' from 8 1/3% to 10%.
	a := profit.ComputeEconomyAnalysis(profit.ConstantCapitalEconomy{
		Kind: profit.KindRawMaterialSaving, ConstantCapital: 11000, VariableCapital: 1000, SurplusValue: 1000, Saving: 2000,
	})
	created, err := m.CreateEconomyAnalysis(ctx, a)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}

	got, err := m.GetEconomyAnalysis(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetEconomyAnalysisNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetEconomyAnalysis(context.Background(), profit.EconomyAnalysisID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateEconomyAnalysisDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	a := profit.ComputeEconomyAnalysis(profit.ConstantCapitalEconomy{
		Kind: profit.KindRawMaterialSaving, ConstantCapital: 11000, VariableCapital: 1000, SurplusValue: 1000, Saving: 2000,
	})
	a.ID = profit.EconomyAnalysisID("fixed-id-7890")

	if _, err := m.CreateEconomyAnalysis(ctx, a); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateEconomyAnalysis(ctx, a); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListEconomyAnalysesNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListEconomyAnalyses(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

// --- ProductionSphere store tests (Ch. 8) ---

func TestMemory_ProductionSphereRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Sphere A: 600c + 100v, s' = 100%.
	sp := avgprofit.ComputeProductionSphere("Sphere A (600c+100v)", 700, 100, 10000)
	created, err := m.CreateProductionSphere(ctx, sp)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}

	got, err := m.GetProductionSphere(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetProductionSphereNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetProductionSphere(context.Background(), avgprofit.ProductionSphereID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateProductionSphereDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	sp := avgprofit.ComputeProductionSphere("Sphere B (100c+600v)", 700, 600, 10000)
	sp.ID = avgprofit.ProductionSphereID("fixed-sphere-id-ch08")

	if _, err := m.CreateProductionSphere(ctx, sp); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateProductionSphere(ctx, sp); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListProductionSpheresNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListProductionSpheres(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

func TestMemory_ListProductionSpheres_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Create Sphere A then Sphere B; list should return B first (newest).
	spA := avgprofit.ComputeProductionSphere("Sphere A", 700, 100, 10000)
	spB := avgprofit.ComputeProductionSphere("Sphere B", 700, 600, 10000)

	createdA, err := m.CreateProductionSphere(ctx, spA)
	if err != nil {
		t.Fatalf("create A: %v", err)
	}
	createdB, err := m.CreateProductionSphere(ctx, spB)
	if err != nil {
		t.Fatalf("create B: %v", err)
	}

	out, err := m.ListProductionSpheres(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("list len = %d, want 2", len(out))
	}
	// Newest first: B was created after A, so B should come first (unless same
	// nanosecond tick — in that case ordering is allowed to be either way).
	if !createdB.CreatedAt.After(createdA.CreatedAt) {
		// Same nanosecond: skip ordering assertion.
		return
	}
	if out[0].ID != createdB.ID {
		t.Errorf("first item id = %s, want B (%s)", out[0].ID, createdB.ID)
	}
}
