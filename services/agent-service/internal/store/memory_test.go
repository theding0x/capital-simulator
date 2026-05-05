package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

func makeAgent(class agent.Class, balance agent.Pence) agent.Agent {
	return agent.Agent{Name: "Test Agent", Class: class, MoneyBalance: balance}
}

func TestMemory_CreateGet(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	created, err := m.Create(ctx, makeAgent(agent.CapitalistClass, 10000))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID.IsZero() {
		t.Error("Create should assign an ID")
	}
	got, err := m.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.Name != created.Name {
		t.Errorf("want name %q, got %q", created.Name, got.Name)
	}
}

func TestMemory_Get_NotFound(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.Get(context.Background(), agent.NewID())
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestMemory_ListByClass(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	if _, err := m.Create(ctx, makeAgent(agent.CapitalistClass, 10000)); err != nil {
		t.Fatalf("Create capitalist: %v", err)
	}
	if _, err := m.Create(ctx, makeAgent(agent.WorkerClass, 500)); err != nil {
		t.Fatalf("Create worker: %v", err)
	}
	if _, err := m.Create(ctx, agent.Agent{Name: "Worker2", Class: agent.WorkerClass, MoneyBalance: 300}); err != nil {
		t.Fatalf("Create worker2: %v", err)
	}
	caps, err := m.ListByClass(ctx, agent.CapitalistClass)
	if err != nil {
		t.Fatalf("ListByClass: %v", err)
	}
	if len(caps) != 1 {
		t.Errorf("want 1 capitalist, got %d", len(caps))
	}
	workers, err := m.ListByClass(ctx, agent.WorkerClass)
	if err != nil {
		t.Fatalf("ListByClass: %v", err)
	}
	if len(workers) != 2 {
		t.Errorf("want 2 workers, got %d", len(workers))
	}
}

func TestMemory_Update(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	a, err := m.Create(ctx, makeAgent(agent.CapitalistClass, 10000))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	newName := "Updated"
	updated, err := m.Update(ctx, a.ID, store.Update{Name: &newName})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != newName {
		t.Errorf("want name %q, got %q", newName, updated.Name)
	}
}

func TestMemory_Delete(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	a, err := m.Create(ctx, makeAgent(agent.CapitalistClass, 10000))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if err := m.Delete(ctx, a.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := m.Delete(ctx, a.ID); !errors.Is(err, store.ErrNotFound) {
		t.Errorf("second Delete want ErrNotFound, got %v", err)
	}
}

// §8: CreateCircuit computes balance update from SurplusValue.
func TestMemory_CreateCircuit_UpdatesBalance(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	a, err := m.Create(ctx, makeAgent(agent.CapitalistClass, 10000))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	c := agent.CapitalCircuit{
		AgentID:      a.ID,
		MAdvanced:    10000,
		CommodityID:  "cotton",
		MReturned:    11000,
		SurplusValue: 1000,
		CircuitType:  agent.CircuitMCM,
	}
	saved, err := m.CreateCircuit(ctx, c)
	if err != nil {
		t.Fatalf("CreateCircuit: %v", err)
	}
	if saved.ID.IsZero() {
		t.Error("CreateCircuit should assign an ID")
	}
	updated, err := m.Get(ctx, a.ID)
	if err != nil {
		t.Fatalf("Get after circuit: %v", err)
	}
	if updated.MoneyBalance != 11000 {
		t.Errorf("want balance 11000 after circuit, got %d", updated.MoneyBalance)
	}
}

// Invariant: CreateCircuit returns ErrInsufficientFunds when balance < MAdvanced.
func TestMemory_CreateCircuit_InsufficientFunds(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	a, err := m.Create(ctx, makeAgent(agent.CapitalistClass, 5000))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	c := agent.CapitalCircuit{
		AgentID:      a.ID,
		MAdvanced:    10000,
		CommodityID:  "cotton",
		MReturned:    11000,
		SurplusValue: 1000,
		CircuitType:  agent.CircuitMCM,
	}
	_, err = m.CreateCircuit(ctx, c)
	if !errors.Is(err, agent.ErrInsufficientFunds) {
		t.Errorf("want ErrInsufficientFunds, got %v", err)
	}
}

func TestMemory_ListCircuits(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	a, err := m.Create(ctx, makeAgent(agent.CapitalistClass, 30000))
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	for i := 0; i < 3; i++ {
		_, err := m.CreateCircuit(ctx, agent.CapitalCircuit{
			AgentID:      a.ID,
			MAdvanced:    1000,
			CommodityID:  "cotton",
			MReturned:    1100,
			SurplusValue: 100,
			CircuitType:  agent.CircuitMCM,
		})
		if err != nil {
			t.Fatalf("CreateCircuit %d: %v", i, err)
		}
	}
	cs, err := m.ListCircuits(ctx, a.ID)
	if err != nil {
		t.Fatalf("ListCircuits: %v", err)
	}
	if len(cs) != 3 {
		t.Errorf("want 3 circuits, got %d", len(cs))
	}
}
