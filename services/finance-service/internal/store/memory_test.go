package store

import (
	"context"
	"errors"
	"testing"

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
