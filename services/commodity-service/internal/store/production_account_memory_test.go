package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
	"github.com/theding0x/capital-simulator/services/commodity-service/internal/store"
)

func makeProductionAccount() commodity.ProductionAccount {
	return commodity.ProductionAccount{
		Constant: 410,
		Variable: 90,
		Surplus:  commodity.SurplusValue(90),
	}
}

func TestMemory_CreateGetProductionAccount(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	a := makeProductionAccount()
	created, err := m.CreateProductionAccount(ctx, a)
	if err != nil {
		t.Fatalf("CreateProductionAccount: %v", err)
	}
	if created.ID.IsZero() {
		t.Error("CreateProductionAccount should assign an ID")
	}
	got, err := m.GetProductionAccount(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetProductionAccount: %v", err)
	}
	if got.Constant != a.Constant || got.Variable != a.Variable || got.Surplus != a.Surplus {
		t.Errorf("round-trip mismatch: want %+v, got %+v", a, got)
	}
}

func TestMemory_GetProductionAccount_NotFound(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.GetProductionAccount(context.Background(), commodity.NewProductionAccountID())
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestMemory_ListProductionAccounts(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	for range 3 {
		_, err := m.CreateProductionAccount(ctx, makeProductionAccount())
		if err != nil {
			t.Fatalf("CreateProductionAccount: %v", err)
		}
	}
	all, err := m.ListProductionAccounts(ctx)
	if err != nil {
		t.Fatalf("ListProductionAccounts: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("want 3 accounts, got %d", len(all))
	}
}
