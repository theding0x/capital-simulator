package store

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/tendency"
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

// --- GeneralProfitRate store tests (Ch. 9) ---

func TestMemory_GeneralProfitRateRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Five-sphere fixture: ΣS=110, ΣC=500 → Rate=2200.
	spheres := []avgprofit.ProductionSphere{
		avgprofit.ComputeProductionSphere("Sphere I (80c+20v)", 100, 20, 10000),
		avgprofit.ComputeProductionSphere("Sphere II (70c+30v)", 100, 30, 10000),
		avgprofit.ComputeProductionSphere("Sphere III (60c+40v)", 100, 40, 10000),
		avgprofit.ComputeProductionSphere("Sphere IV (85c+15v)", 100, 15, 10000),
		avgprofit.ComputeProductionSphere("Sphere V (95c+5v)", 100, 5, 10000),
	}
	g := avgprofit.ComputeGeneralProfitRate(spheres)
	created, err := m.CreateGeneralProfitRate(ctx, g)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.Rate != 2200 {
		t.Errorf("Rate = %d, want 2200", created.Rate)
	}

	got, err := m.GetGeneralProfitRate(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID || got.Rate != created.Rate {
		t.Errorf("round-trip mismatch: got ID=%s Rate=%d, want ID=%s Rate=%d",
			got.ID, got.Rate, created.ID, created.Rate)
	}
}

func TestMemory_GetGeneralProfitRateNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetGeneralProfitRate(context.Background(), avgprofit.GeneralProfitRateID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateGeneralProfitRateDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	g := avgprofit.ComputeGeneralProfitRate(nil)
	g.ID = avgprofit.GeneralProfitRateID("fixed-gpr-id-ch09")

	if _, err := m.CreateGeneralProfitRate(ctx, g); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateGeneralProfitRate(ctx, g); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

// --- PriceOfProduction store tests (Ch. 9) ---

func TestMemory_PriceOfProductionRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Sphere I (80c+20v): k=100, rate=2200, value=120, price=122, deviation=+2.
	pop := avgprofit.ComputePriceOfProduction("Sphere I (80c+20v)", 100, 2200, 120)
	created, err := m.CreatePriceOfProduction(ctx, pop)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.Price != 122 {
		t.Errorf("Price = %d, want 122", created.Price)
	}
	if created.Deviation != 2 {
		t.Errorf("Deviation = %d, want +2", created.Deviation)
	}

	got, err := m.GetPriceOfProduction(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID || got.Price != created.Price || got.Deviation != created.Deviation {
		t.Errorf("round-trip mismatch: got %+v, want ID=%s Price=%d Deviation=%d",
			got, created.ID, created.Price, created.Deviation)
	}
}

func TestMemory_GetPriceOfProductionNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetPriceOfProduction(context.Background(), avgprofit.PriceOfProductionID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreatePriceOfProductionDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	pop := avgprofit.ComputePriceOfProduction("Sphere I", 100, 2200, 120)
	pop.ID = avgprofit.PriceOfProductionID("fixed-pop-id-ch09")

	if _, err := m.CreatePriceOfProduction(ctx, pop); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreatePriceOfProduction(ctx, pop); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListPricesOfProductionNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListPricesOfProduction(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

func TestMemory_ListPricesOfProduction_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Create Sphere I then Sphere III; list should return Sphere III first (newest).
	popI := avgprofit.ComputePriceOfProduction("Sphere I (80c+20v)", 100, 2200, 120)
	popIII := avgprofit.ComputePriceOfProduction("Sphere III (60c+40v)", 100, 2200, 140)

	createdI, err := m.CreatePriceOfProduction(ctx, popI)
	if err != nil {
		t.Fatalf("create I: %v", err)
	}
	createdIII, err := m.CreatePriceOfProduction(ctx, popIII)
	if err != nil {
		t.Fatalf("create III: %v", err)
	}

	out, err := m.ListPricesOfProduction(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("list len = %d, want 2", len(out))
	}
	// Newest first: III was created after I, so III should come first (unless same
	// nanosecond tick — in that case ordering is allowed to be either way).
	if !createdIII.CreatedAt.After(createdI.CreatedAt) {
		return
	}
	if out[0].ID != createdIII.ID {
		t.Errorf("first item id = %s, want Sphere III (%s)", out[0].ID, createdIII.ID)
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

// --- MarketValue store tests (Ch. 10) ---

func TestMemory_MarketValueRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Cotton: bulk=100, best=80, worst=130 (spec fixture).
	mv := avgprofit.ComputeMarketValue("cotton", 100, 80, 130)
	created, err := m.CreateMarketValue(ctx, mv)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.Value != 100 {
		t.Errorf("Value = %d, want 100 (bulk condition)", created.Value)
	}

	got, err := m.GetMarketValue(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetMarketValueNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetMarketValue(context.Background(), avgprofit.MarketValueID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateMarketValueDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	mv := avgprofit.ComputeMarketValue("cotton", 100, 80, 130)
	mv.ID = avgprofit.MarketValueID("fixed-mv-id-0001")

	if _, err := m.CreateMarketValue(ctx, mv); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateMarketValue(ctx, mv); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

// --- SurplusProfit store tests (Ch. 10) ---

func TestMemory_SurplusProfitRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// High-productivity firm: indVal=60, marketVal=100, qty=1000 → amount=40000.
	sp := avgprofit.ComputeSurplusProfit("High-productivity firm", 60, 100, 1000, 2200)
	created, err := m.CreateSurplusProfit(ctx, sp)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.Amount != 40000 {
		t.Errorf("Amount = %d, want 40000", created.Amount)
	}

	got, err := m.GetSurplusProfit(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetSurplusProfitNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetSurplusProfit(context.Background(), avgprofit.SurplusProfitID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateSurplusProfitDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	sp := avgprofit.ComputeSurplusProfit("firm", 60, 100, 1000, 2200)
	sp.ID = avgprofit.SurplusProfitID("fixed-sp-id-0001")

	if _, err := m.CreateSurplusProfit(ctx, sp); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateSurplusProfit(ctx, sp); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

// --- Equalisation store tests (Ch. 10) ---

func TestMemory_EqualisationRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Cotton sphere: initialRate=3000, targetRate=2200 → inflow, converging.
	eq := avgprofit.ComputeEqualisation("cotton", 3000, 2200)
	created, err := m.CreateEqualisation(ctx, eq)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if !created.IsConverging {
		t.Error("IsConverging = false, want true")
	}
	if created.Direction != avgprofit.KindInflow {
		t.Errorf("Direction = %q, want %q", created.Direction, avgprofit.KindInflow)
	}

	got, err := m.GetEqualisation(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetEqualisationNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetEqualisation(context.Background(), avgprofit.EqualisationID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateEqualisationDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	eq := avgprofit.ComputeEqualisation("cotton", 3000, 2200)
	eq.ID = avgprofit.EqualisationID("fixed-eq-id-0001")

	if _, err := m.CreateEqualisation(ctx, eq); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateEqualisation(ctx, eq); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

// --- WageEffectAnalysis store tests (Ch. 11) ---

func TestMemory_WageEffectAnalysisRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// §main fixture: base 80c+20v, s'=100%, factor 125 (wage rise).
	spheres := []avgprofit.SphereInput{
		{Name: "lower", C: 50, V: 50},
		{Name: "higher", C: 92, V: 8},
	}
	a := avgprofit.ComputeWageEffectAnalysis(80, 20, 10000, 125, spheres)

	created, err := m.CreateWageEffectAnalysis(ctx, a)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.OldGeneralRate != 2000 {
		t.Errorf("OldGeneralRate = %d, want 2000", created.OldGeneralRate)
	}
	if created.NewGeneralRate != 1429 {
		t.Errorf("NewGeneralRate = %d, want 1429", created.NewGeneralRate)
	}
	if len(created.Outcomes) != 2 {
		t.Fatalf("Outcomes len = %d, want 2", len(created.Outcomes))
	}

	got, err := m.GetWageEffectAnalysis(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("round-trip ID mismatch: got %s, want %s", got.ID, created.ID)
	}
	if got.OldGeneralRate != created.OldGeneralRate || got.NewGeneralRate != created.NewGeneralRate {
		t.Errorf("round-trip rate mismatch: got old=%d new=%d, want old=%d new=%d",
			got.OldGeneralRate, got.NewGeneralRate, created.OldGeneralRate, created.NewGeneralRate)
	}
	if len(got.Outcomes) != len(created.Outcomes) {
		t.Errorf("Outcomes len mismatch: got %d, want %d", len(got.Outcomes), len(created.Outcomes))
	}
	for i, o := range got.Outcomes {
		if o != created.Outcomes[i] {
			t.Errorf("Outcomes[%d] mismatch:\n got  %+v\n want %+v", i, o, created.Outcomes[i])
		}
	}
	if got.AverageOutcome != created.AverageOutcome {
		t.Errorf("AverageOutcome mismatch:\n got  %+v\n want %+v", got.AverageOutcome, created.AverageOutcome)
	}
}

func TestMemory_GetWageEffectAnalysisNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetWageEffectAnalysis(context.Background(), avgprofit.WageEffectAnalysisID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateWageEffectAnalysisDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	a := avgprofit.ComputeWageEffectAnalysis(80, 20, 10000, 125, nil)
	a.ID = avgprofit.WageEffectAnalysisID("fixed-wage-id-ch11")

	if _, err := m.CreateWageEffectAnalysis(ctx, a); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateWageEffectAnalysis(ctx, a); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListWageEffectAnalysesNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListWageEffectAnalyses(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

// --- PriceOfProductionChange store tests (Ch. 12) ---

func TestMemory_PriceOfProductionChangeRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Cotton: general-rate change, old price 122 → new price 130.
	c := avgprofit.ComputePriceOfProductionChange("cotton", avgprofit.CauseGeneralRateChange, 122, 130)
	created, err := m.CreatePriceOfProductionChange(ctx, c)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}

	got, err := m.GetPriceOfProductionChange(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetPriceOfProductionChangeNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetPriceOfProductionChange(context.Background(), avgprofit.PriceOfProductionChangeID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreatePriceOfProductionChangeDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	c := avgprofit.ComputePriceOfProductionChange("cotton", avgprofit.CauseGeneralRateChange, 122, 130)
	c.ID = avgprofit.PriceOfProductionChangeID("fixed-ppc-id-ch12")

	if _, err := m.CreatePriceOfProductionChange(ctx, c); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreatePriceOfProductionChange(ctx, c); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListWageEffectAnalyses_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Create two analyses; the second should appear first in the list.
	aFirst := avgprofit.ComputeWageEffectAnalysis(80, 20, 10000, 125, nil)
	aSecond := avgprofit.ComputeWageEffectAnalysis(80, 20, 10000, 75, nil)

	createdFirst, err := m.CreateWageEffectAnalysis(ctx, aFirst)
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	createdSecond, err := m.CreateWageEffectAnalysis(ctx, aSecond)
	if err != nil {
		t.Fatalf("create second: %v", err)
	}

	out, err := m.ListWageEffectAnalyses(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("list len = %d, want 2", len(out))
	}
	// Skip ordering check if both have the same nanosecond timestamp.
	if !createdSecond.CreatedAt.After(createdFirst.CreatedAt) {
		return
	}
	if out[0].ID != createdSecond.ID {
		t.Errorf("first item id = %s, want newest (%s)", out[0].ID, createdSecond.ID)
	}
}

// --- CompositionTrajectory store tests (Ch. 13) ---

// ch13LawSteps mirrors the §law fixture: s′=100%, v=100, c rising 50→400.
var ch13LawSteps = [][2]int64{{50, 100}, {100, 100}, {200, 100}, {300, 100}, {400, 100}}

// ch13WantRates are the round-half-up rates the §law fixture produces.
var ch13WantRates = []int64{6667, 5000, 3333, 2500, 2000}

func TestMemory_CompositionTrajectoryRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	tr := tendency.BuildCompositionTrajectory("Marx §law s'=100%", 10000, ch13LawSteps)
	created, err := m.CreateCompositionTrajectory(ctx, tr)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if len(created.ProfitRates) != len(ch13WantRates) {
		t.Fatalf("ProfitRates len = %d, want %d", len(created.ProfitRates), len(ch13WantRates))
	}
	for i, want := range ch13WantRates {
		if created.ProfitRates[i] != want {
			t.Errorf("created ProfitRates[%d] = %d, want %d", i, created.ProfitRates[i], want)
		}
	}

	got, err := m.GetCompositionTrajectory(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID || got.Label != created.Label || got.SurplusValueRate != created.SurplusValueRate {
		t.Errorf("round-trip scalar mismatch:\n got  %+v\n want %+v", got, created)
	}
	if !got.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("round-trip CreatedAt mismatch: got %v, want %v", got.CreatedAt, created.CreatedAt)
	}
	if len(got.Periods) != len(created.Periods) {
		t.Fatalf("round-trip Periods len = %d, want %d", len(got.Periods), len(created.Periods))
	}
	for i := range created.Periods {
		if got.Periods[i] != created.Periods[i] {
			t.Errorf("round-trip Periods[%d] mismatch:\n got  %+v\n want %+v", i, got.Periods[i], created.Periods[i])
		}
	}
	// Derived ProfitRates survives the round-trip (re-derived on Get).
	if len(got.ProfitRates) != len(created.ProfitRates) {
		t.Fatalf("round-trip ProfitRates len = %d, want %d", len(got.ProfitRates), len(created.ProfitRates))
	}
	for i := range created.ProfitRates {
		if got.ProfitRates[i] != created.ProfitRates[i] {
			t.Errorf("round-trip ProfitRates[%d] = %d, want %d", i, got.ProfitRates[i], created.ProfitRates[i])
		}
	}
}

func TestMemory_GetCompositionTrajectoryNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetCompositionTrajectory(context.Background(), tendency.CompositionTrajectoryID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateCompositionTrajectoryDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	tr := tendency.BuildCompositionTrajectory("law", 10000, ch13LawSteps)
	tr.ID = tendency.CompositionTrajectoryID("fixed-traj-id-ch13")

	if _, err := m.CreateCompositionTrajectory(ctx, tr); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateCompositionTrajectory(ctx, tr); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListCompositionTrajectoriesNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListCompositionTrajectories(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

func TestMemory_ListCompositionTrajectories_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	first := tendency.BuildCompositionTrajectory("first", 10000, [][2]int64{{50, 100}})
	second := tendency.BuildCompositionTrajectory("second", 10000, [][2]int64{{400, 100}})

	createdFirst, err := m.CreateCompositionTrajectory(ctx, first)
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	createdSecond, err := m.CreateCompositionTrajectory(ctx, second)
	if err != nil {
		t.Fatalf("create second: %v", err)
	}

	out, err := m.ListCompositionTrajectories(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("list len = %d, want 2", len(out))
	}
	// Skip ordering check if both share the same nanosecond timestamp.
	if !createdSecond.CreatedAt.After(createdFirst.CreatedAt) {
		return
	}
	if out[0].ID != createdSecond.ID {
		t.Errorf("first item id = %s, want newest (%s)", out[0].ID, createdSecond.ID)
	}
}

// --- RateMassContradiction store tests (Ch. 13) ---

func TestMemory_RateMassContradictionRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Rate halves (20%→10%), capital doubles → mass preserved (MassChange 0).
	rm := tendency.ComputeRateMassContradiction(1_000_000, 2000, 2_000_000, 1000)
	created, err := m.CreateRateMassContradiction(ctx, rm)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.OldMass != 200000 || created.NewMass != 200000 || created.MassChange != 0 {
		t.Errorf("masses = old %d new %d change %d, want 200000/200000/0",
			created.OldMass, created.NewMass, created.MassChange)
	}

	got, err := m.GetRateMassContradiction(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetRateMassContradictionNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetRateMassContradiction(context.Background(), tendency.RateMassContradictionID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateRateMassContradictionDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	rm := tendency.ComputeRateMassContradiction(1_000_000, 2000, 3_000_000, 1000)
	rm.ID = tendency.RateMassContradictionID("fixed-ratemass-id-ch13")

	if _, err := m.CreateRateMassContradiction(ctx, rm); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateRateMassContradiction(ctx, rm); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListRateMassContradictionsNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListRateMassContradictions(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

func TestMemory_ListRateMassContradictions_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	first := tendency.ComputeRateMassContradiction(1_000_000, 2000, 2_000_000, 1000)
	second := tendency.ComputeRateMassContradiction(1_000_000, 2000, 3_000_000, 1000)

	createdFirst, err := m.CreateRateMassContradiction(ctx, first)
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	createdSecond, err := m.CreateRateMassContradiction(ctx, second)
	if err != nil {
		t.Fatalf("create second: %v", err)
	}

	out, err := m.ListRateMassContradictions(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("list len = %d, want 2", len(out))
	}
	if !createdSecond.CreatedAt.After(createdFirst.CreatedAt) {
		return
	}
	if out[0].ID != createdSecond.ID {
		t.Errorf("first item id = %s, want newest (%s)", out[0].ID, createdSecond.ID)
	}
}

// TestCh13SeedIDs asserts the Ch. 13 seed IDs follow the 5eed…<CC=13> convention
// and are unique (mirrors the seed migration 00026_v3_ch13_seed.sql).
func TestCh13SeedIDs(t *testing.T) {
	t.Parallel()
	ids := []string{
		"5eed000000000000001301",
		"5eed000000000000001302",
		"5eed000000000000001303",
	}
	const prefix = "5eed00000000000000"
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if len(id) <= len(prefix) || id[:len(prefix)] != prefix {
			t.Errorf("seed id %q lacks prefix %q", id, prefix)
		}
		if cc := id[len(prefix) : len(prefix)+2]; cc != "13" {
			t.Errorf("seed id %q chapter token = %q, want 13", id, cc)
		}
		if _, dup := seen[id]; dup {
			t.Errorf("duplicate seed id %q", id)
		}
		seen[id] = struct{}{}
	}
}

// --- CounteractingForce store tests (Ch. 14) ---

func TestMemory_CounteractingForceRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// §VI stock capital — the one kind excluded from the general rate.
	f, err := tendency.NewCounteractingForce(tendency.KindStockCapital, "Marx §VI")
	if err != nil {
		t.Fatalf("new force: %v", err)
	}
	created, err := m.CreateCounteractingForce(ctx, f)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if !created.ExcludedFromGeneralRate {
		t.Error("ExcludedFromGeneralRate = false, want true for stock_capital")
	}

	got, err := m.GetCounteractingForce(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetCounteractingForceNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetCounteractingForce(context.Background(), tendency.CounteractingForceID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateCounteractingForceDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	f, err := tendency.NewCounteractingForce(tendency.KindForeignTrade, "")
	if err != nil {
		t.Fatalf("new force: %v", err)
	}
	f.ID = tendency.CounteractingForceID("5eed000000000000001405")

	if _, err := m.CreateCounteractingForce(ctx, f); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateCounteractingForce(ctx, f); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

// --- CounteractingScenario store tests (Ch. 14) ---

// ch14LawSteps mirrors the §law base trajectory the scenario is built on.
var ch14LawSteps = [][2]int64{{50, 100}, {100, 100}, {200, 100}, {300, 100}, {400, 100}}

// ch14IndustrialKinds are the five forces (§I–§V) that uplift the path.
var ch14IndustrialKinds = []tendency.CounteractingForceKind{
	tendency.KindIntensifiedExploitation,
	tendency.KindDepressedWages,
	tendency.KindCheaperConstantCapital,
	tendency.KindRelativeOverpopulation,
	tendency.KindForeignTrade,
}

// ch14Forces builds the five-force slice, failing the test on a rejected kind.
func ch14Forces(t *testing.T) []tendency.CounteractingForce {
	t.Helper()
	forces := make([]tendency.CounteractingForce, 0, len(ch14IndustrialKinds))
	for _, k := range ch14IndustrialKinds {
		f, err := tendency.NewCounteractingForce(k, "")
		if err != nil {
			t.Fatalf("new force %q: %v", k, err)
		}
		forces = append(forces, f)
	}
	return forces
}

func TestMemory_CounteractingScenarioRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	trajectory := tendency.BuildCompositionTrajectory("law", 10000, ch14LawSteps)
	scenario, err := tendency.BuildCounteractingScenario("All five industrial forces", trajectory, ch14Forces(t))
	if err != nil {
		t.Fatalf("build scenario: %v", err)
	}

	created, err := m.CreateCounteractingScenario(ctx, scenario)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}

	got, err := m.GetCounteractingScenario(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != created.ID || got.Label != created.Label {
		t.Errorf("round-trip scalar mismatch:\n got  %+v\n want %+v", got, created)
	}
	if got.InitialProfitRate != created.InitialProfitRate || got.FinalProfitRate != created.FinalProfitRate {
		t.Errorf("round-trip rate mismatch: got init=%d final=%d, want init=%d final=%d",
			got.InitialProfitRate, got.FinalProfitRate, created.InitialProfitRate, created.FinalProfitRate)
	}
	if !got.CreatedAt.Equal(created.CreatedAt) {
		t.Errorf("round-trip CreatedAt mismatch: got %v, want %v", got.CreatedAt, created.CreatedAt)
	}
	// JSON-roundtripped Forces survive.
	if len(got.Forces) != len(created.Forces) {
		t.Fatalf("Forces len = %d, want %d", len(got.Forces), len(created.Forces))
	}
	for i := range created.Forces {
		if got.Forces[i] != created.Forces[i] {
			t.Errorf("Forces[%d] mismatch:\n got  %+v\n want %+v", i, got.Forces[i], created.Forces[i])
		}
	}
	// JSON-roundtripped ModifiedRates survive.
	if len(got.ModifiedRates) != len(created.ModifiedRates) {
		t.Fatalf("ModifiedRates len = %d, want %d", len(got.ModifiedRates), len(created.ModifiedRates))
	}
	for i := range created.ModifiedRates {
		if got.ModifiedRates[i] != created.ModifiedRates[i] {
			t.Errorf("ModifiedRates[%d] = %d, want %d", i, got.ModifiedRates[i], created.ModifiedRates[i])
		}
	}
	// JSON-roundtripped BaseTrajectory periods survive.
	if len(got.BaseTrajectory.Periods) != len(created.BaseTrajectory.Periods) {
		t.Fatalf("BaseTrajectory.Periods len = %d, want %d",
			len(got.BaseTrajectory.Periods), len(created.BaseTrajectory.Periods))
	}
	for i := range created.BaseTrajectory.Periods {
		if got.BaseTrajectory.Periods[i] != created.BaseTrajectory.Periods[i] {
			t.Errorf("BaseTrajectory.Periods[%d] mismatch:\n got  %+v\n want %+v",
				i, got.BaseTrajectory.Periods[i], created.BaseTrajectory.Periods[i])
		}
	}
}

func TestMemory_GetCounteractingScenarioNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetCounteractingScenario(context.Background(), tendency.CounteractingScenarioID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateCounteractingScenarioDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	trajectory := tendency.BuildCompositionTrajectory("law", 10000, ch14LawSteps)
	scenario, err := tendency.BuildCounteractingScenario("law", trajectory, ch14Forces(t))
	if err != nil {
		t.Fatalf("build scenario: %v", err)
	}
	scenario.ID = tendency.CounteractingScenarioID("5eed000000000000001407")

	if _, err := m.CreateCounteractingScenario(ctx, scenario); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateCounteractingScenario(ctx, scenario); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListCounteractingScenariosNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListCounteractingScenarios(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

func TestMemory_ListCounteractingScenarios_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	trajectory := tendency.BuildCompositionTrajectory("law", 10000, ch14LawSteps)
	first, err := tendency.BuildCounteractingScenario("first", trajectory, ch14Forces(t))
	if err != nil {
		t.Fatalf("build first: %v", err)
	}
	second, err := tendency.BuildCounteractingScenario("second", trajectory, ch14Forces(t))
	if err != nil {
		t.Fatalf("build second: %v", err)
	}

	createdFirst, err := m.CreateCounteractingScenario(ctx, first)
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	createdSecond, err := m.CreateCounteractingScenario(ctx, second)
	if err != nil {
		t.Fatalf("create second: %v", err)
	}

	out, err := m.ListCounteractingScenarios(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("list len = %d, want 2", len(out))
	}
	// Skip ordering check if both share the same nanosecond timestamp.
	if !createdSecond.CreatedAt.After(createdFirst.CreatedAt) {
		return
	}
	if out[0].ID != createdSecond.ID {
		t.Errorf("first item id = %s, want newest (%s)", out[0].ID, createdSecond.ID)
	}
}

// TestCh14SeedIDs asserts the Ch. 14 seed IDs (5eed…<CC=14>01–07) carry the 14
// token and are unique (mirrors the seed migration 00028_v3_ch14_seed.sql).
func TestCh14SeedIDs(t *testing.T) {
	t.Parallel()
	ids := []string{
		"5eed000000000000001401",
		"5eed000000000000001402",
		"5eed000000000000001403",
		"5eed000000000000001404",
		"5eed000000000000001405",
		"5eed000000000000001406",
		"5eed000000000000001407",
	}
	const prefix = "5eed00000000000000"
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if len(id) <= len(prefix) || id[:len(prefix)] != prefix {
			t.Errorf("seed id %q lacks prefix %q", id, prefix)
		}
		if cc := id[len(prefix) : len(prefix)+2]; cc != "14" {
			t.Errorf("seed id %q chapter token = %q, want 14", id, cc)
		}
		if _, dup := seen[id]; dup {
			t.Errorf("duplicate seed id %q", id)
		}
		seen[id] = struct{}{}
	}
}

// --- Crisis store tests (Ch. 15) ---

func TestMemory_CrisisRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// §crisis-A: 30% writedown on a 1000 bp pre-crisis rate restores it to 1429.
	c, err := tendency.NewCrisis(30, 1000)
	if err != nil {
		t.Fatalf("new crisis: %v", err)
	}
	created, err := m.CreateCrisis(ctx, c)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if created.PostCrisisProfitRate != 1429 {
		t.Errorf("PostCrisisProfitRate = %d, want 1429", created.PostCrisisProfitRate)
	}

	got, err := m.GetCrisis(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetCrisisNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetCrisis(context.Background(), tendency.CrisisID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateCrisisDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	c, err := tendency.NewCrisis(30, 1000)
	if err != nil {
		t.Fatalf("new crisis: %v", err)
	}
	c.ID = tendency.CrisisID("5eed000000000000001501")

	if _, err := m.CreateCrisis(ctx, c); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateCrisis(ctx, c); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListCrisesNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListCrises(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

func TestMemory_ListCrises_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	first, err := tendency.NewCrisis(20, 800)
	if err != nil {
		t.Fatalf("new first: %v", err)
	}
	second, err := tendency.NewCrisis(50, 600)
	if err != nil {
		t.Fatalf("new second: %v", err)
	}

	createdFirst, err := m.CreateCrisis(ctx, first)
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	createdSecond, err := m.CreateCrisis(ctx, second)
	if err != nil {
		t.Fatalf("create second: %v", err)
	}

	out, err := m.ListCrises(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("list len = %d, want 2", len(out))
	}
	// Skip ordering check if both share the same nanosecond timestamp.
	if !createdSecond.CreatedAt.After(createdFirst.CreatedAt) {
		return
	}
	if out[0].ID != createdSecond.ID {
		t.Errorf("first item id = %s, want newest (%s)", out[0].ID, createdSecond.ID)
	}
}

// --- InternalContradiction store tests (Ch. 15) ---

func TestMemory_InternalContradictionRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// §I falling rate / rising mass — coexistent by construction.
	c, err := tendency.NewInternalContradiction(tendency.KindFallingRateRisingMass, "Marx §I")
	if err != nil {
		t.Fatalf("new contradiction: %v", err)
	}
	created, err := m.CreateInternalContradiction(ctx, c)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	if !created.IsCoexistent {
		t.Error("IsCoexistent = false, want true")
	}

	got, err := m.GetInternalContradiction(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetInternalContradictionNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetInternalContradiction(context.Background(), tendency.InternalContradictionID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateInternalContradictionDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()
	c, err := tendency.NewInternalContradiction(tendency.KindConcentrationCentralisation, "")
	if err != nil {
		t.Fatalf("new contradiction: %v", err)
	}
	c.ID = tendency.InternalContradictionID("5eed000000000000001507")

	if _, err := m.CreateInternalContradiction(ctx, c); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateInternalContradiction(ctx, c); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListInternalContradictionsNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListInternalContradictions(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

func TestMemory_ListInternalContradictions_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	first, err := tendency.NewInternalContradiction(tendency.KindFallingRateRisingMass, "")
	if err != nil {
		t.Fatalf("new first: %v", err)
	}
	second, err := tendency.NewInternalContradiction(tendency.KindOverproductionUnderConsumption, "")
	if err != nil {
		t.Fatalf("new second: %v", err)
	}

	createdFirst, err := m.CreateInternalContradiction(ctx, first)
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	createdSecond, err := m.CreateInternalContradiction(ctx, second)
	if err != nil {
		t.Fatalf("create second: %v", err)
	}

	out, err := m.ListInternalContradictions(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("list len = %d, want 2", len(out))
	}
	if !createdSecond.CreatedAt.After(createdFirst.CreatedAt) {
		return
	}
	if out[0].ID != createdSecond.ID {
		t.Errorf("first item id = %s, want newest (%s)", out[0].ID, createdSecond.ID)
	}
}

// TestCh15SeedIDs asserts the Ch. 15 seed IDs (5eed…<CC=15>01–07) carry the 15
// token and are unique (mirrors the seed migration 00030_v3_ch15_seed.sql:
// crises 1501–1503, contradictions 1504–1507).
func TestCh15SeedIDs(t *testing.T) {
	t.Parallel()
	ids := []string{
		"5eed000000000000001501",
		"5eed000000000000001502",
		"5eed000000000000001503",
		"5eed000000000000001504",
		"5eed000000000000001505",
		"5eed000000000000001506",
		"5eed000000000000001507",
	}
	const prefix = "5eed00000000000000"
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if len(id) <= len(prefix) || id[:len(prefix)] != prefix {
			t.Errorf("seed id %q lacks prefix %q", id, prefix)
		}
		if cc := id[len(prefix) : len(prefix)+2]; cc != "15" {
			t.Errorf("seed id %q chapter token = %q, want 15", id, cc)
		}
		if _, dup := seen[id]; dup {
			t.Errorf("duplicate seed id %q", id)
		}
		seen[id] = struct{}{}
	}
}

// --- CommercialCapital store tests (Ch. 16) ---

func TestMemory_CommercialCapitalRoundTrip(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// Marx §1 linen-merchant fixture: £3,000 advanced, realisation function.
	cc, err := merchant.NewCommercialCapital(300000, "30,000 yards linen at 10p", merchant.FunctionRealisation)
	if err != nil {
		t.Fatalf("new commercial capital: %v", err)
	}
	created, err := m.CreateCommercialCapital(ctx, cc)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected an assigned ID")
	}
	if created.CreatedAt.IsZero() {
		t.Error("expected a created-at timestamp")
	}
	// Central invariant: commercial capital creates no surplus-value.
	if created.SurplusValueProduced != 0 {
		t.Errorf("SurplusValueProduced = %d, want 0", created.SurplusValueProduced)
	}

	got, err := m.GetCommercialCapital(ctx, created.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip mismatch:\n got  %+v\n want %+v", got, created)
	}
}

func TestMemory_GetCommercialCapitalNotFound(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	if _, err := m.GetCommercialCapital(context.Background(), merchant.CommercialCapitalID("missing")); !errors.Is(err, ErrNotFound) {
		t.Errorf("err = %v, want ErrNotFound", err)
	}
}

func TestMemory_CreateCommercialCapitalDuplicate(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	cc, err := merchant.NewCommercialCapital(300000, "linen", merchant.FunctionRealisation)
	if err != nil {
		t.Fatalf("new: %v", err)
	}
	cc.ID = merchant.CommercialCapitalID("5eed000000000000001601")

	if _, err := m.CreateCommercialCapital(ctx, cc); err != nil {
		t.Fatalf("first create: %v", err)
	}
	if _, err := m.CreateCommercialCapital(ctx, cc); !errors.Is(err, ErrAlreadyExists) {
		t.Errorf("second create err = %v, want ErrAlreadyExists", err)
	}
}

func TestMemory_ListCommercialCapitalsNeverNil(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	out, err := m.ListCommercialCapitals(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if out == nil {
		t.Error("list should return a non-nil slice even when empty")
	}
}

func TestMemory_ListCommercialCapitals_NewestFirst(t *testing.T) {
	t.Parallel()
	m := NewMemory()
	ctx := context.Background()

	// §1 realisation, §2 storage — second should appear first in the list.
	first, err := merchant.NewCommercialCapital(300000, "linen (realisation)", merchant.FunctionRealisation)
	if err != nil {
		t.Fatalf("new first: %v", err)
	}
	second, err := merchant.NewCommercialCapital(300000, "linen (storage)", merchant.FunctionStorage)
	if err != nil {
		t.Fatalf("new second: %v", err)
	}

	createdFirst, err := m.CreateCommercialCapital(ctx, first)
	if err != nil {
		t.Fatalf("create first: %v", err)
	}
	createdSecond, err := m.CreateCommercialCapital(ctx, second)
	if err != nil {
		t.Fatalf("create second: %v", err)
	}

	out, err := m.ListCommercialCapitals(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("list len = %d, want 2", len(out))
	}
	// Skip ordering check if both share the same nanosecond timestamp.
	if !createdSecond.CreatedAt.After(createdFirst.CreatedAt) {
		return
	}
	if out[0].ID != createdSecond.ID {
		t.Errorf("first item id = %s, want newest (%s)", out[0].ID, createdSecond.ID)
	}
}

// TestCh16SeedIDs asserts the Ch. 16 seed IDs (5eed…<CC=16>01–03) carry the 16
// token and are unique (mirrors the seed migration 00032_v3_ch16_seed.sql).
func TestCh16SeedIDs(t *testing.T) {
	t.Parallel()
	ids := []string{
		"5eed000000000000001601",
		"5eed000000000000001602",
		"5eed000000000000001603",
	}
	const prefix = "5eed00000000000000"
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if len(id) <= len(prefix) || id[:len(prefix)] != prefix {
			t.Errorf("seed id %q lacks prefix %q", id, prefix)
		}
		if cc := id[len(prefix) : len(prefix)+2]; cc != "16" {
			t.Errorf("seed id %q chapter token = %q, want 16", id, cc)
		}
		if _, dup := seen[id]; dup {
			t.Errorf("duplicate seed id %q", id)
		}
		seen[id] = struct{}{}
	}
}
