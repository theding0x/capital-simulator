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
