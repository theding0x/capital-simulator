package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

func makeWorker() agent.Worker {
	return agent.Worker{
		OwnsLabourPower:       true,
		OwnsCommoditiesToSell: false,
		LabourPower:           agent.LabourPower{CapacityMinutesPerDay: 480},
	}
}

func makeCapitalist() agent.Capitalist {
	return agent.Capitalist{
		MoneyCapital: 1000,
	}
}

func makeOffering(ownerID agent.AgentID) agent.LabourPowerOffering {
	return agent.LabourPowerOffering{
		OwnerID:               ownerID,
		CapacityMinutesPerDay: 480,
		ContractDays:          5,
		AskingWage:            240,
	}
}

func TestMemory_CreateGetWorker(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	created, err := m.CreateWorker(ctx, makeWorker())
	if err != nil {
		t.Fatalf("CreateWorker: %v", err)
	}
	if created.ID.IsZero() {
		t.Error("CreateWorker should assign an ID")
	}
	if created.Kind != agent.AgentKindWorker {
		t.Errorf("want Kind %q, got %q", agent.AgentKindWorker, created.Kind)
	}
	got, err := m.GetWorker(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetWorker: %v", err)
	}
	if got.LabourPower.CapacityMinutesPerDay != created.LabourPower.CapacityMinutesPerDay {
		t.Errorf("capacity mismatch: want %d, got %d",
			created.LabourPower.CapacityMinutesPerDay,
			got.LabourPower.CapacityMinutesPerDay)
	}
}

func TestMemory_GetWorker_NotFound(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.GetWorker(context.Background(), agent.NewAgentID())
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestMemory_ListWorkers(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	if _, err := m.CreateWorker(ctx, makeWorker()); err != nil {
		t.Fatalf("CreateWorker: %v", err)
	}
	if _, err := m.CreateWorker(ctx, makeWorker()); err != nil {
		t.Fatalf("CreateWorker 2: %v", err)
	}
	workers, err := m.ListWorkers(ctx)
	if err != nil {
		t.Fatalf("ListWorkers: %v", err)
	}
	if len(workers) != 2 {
		t.Errorf("want 2 workers, got %d", len(workers))
	}
}

func TestMemory_CreateGetCapitalist(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	created, err := m.CreateCapitalist(ctx, makeCapitalist())
	if err != nil {
		t.Fatalf("CreateCapitalist: %v", err)
	}
	if created.ID.IsZero() {
		t.Error("CreateCapitalist should assign an ID")
	}
	if created.Kind != agent.AgentKindCapitalist {
		t.Errorf("want Kind %q, got %q", agent.AgentKindCapitalist, created.Kind)
	}
	got, err := m.GetCapitalist(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetCapitalist: %v", err)
	}
	if got.MoneyCapital != created.MoneyCapital {
		t.Errorf("money_capital mismatch: want %d, got %d", created.MoneyCapital, got.MoneyCapital)
	}
}

func TestMemory_CreateGetOffering(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	ownerID := agent.NewAgentID()
	created, err := m.CreateOffering(ctx, makeOffering(ownerID))
	if err != nil {
		t.Fatalf("CreateOffering: %v", err)
	}
	if created.ID.IsZero() {
		t.Error("CreateOffering should assign an ID")
	}
	list, err := m.ListOfferings(ctx)
	if err != nil {
		t.Fatalf("ListOfferings: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("want 1 offering, got %d", len(list))
	}
	if list[0].ID != created.ID {
		t.Errorf("listing ID mismatch: want %q, got %q", created.ID, list[0].ID)
	}
}

func TestMemory_CreateGetPurchase(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	p := agent.LabourPowerPurchase{
		SellerID:     agent.NewAgentID(),
		BuyerID:      agent.NewAgentID(),
		WageMinutes:  240,
		ContractDays: 5,
	}
	created, err := m.CreatePurchase(ctx, p)
	if err != nil {
		t.Fatalf("CreatePurchase: %v", err)
	}
	if created.ID.IsZero() {
		t.Error("CreatePurchase should assign an ID")
	}
	got, err := m.GetPurchase(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetPurchase: %v", err)
	}
	if got.WageMinutes != p.WageMinutes {
		t.Errorf("wage_minutes mismatch: want %d, got %d", p.WageMinutes, got.WageMinutes)
	}
}

func TestMemory_GetPurchase_NotFound(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.GetPurchase(context.Background(), agent.NewPurchaseID())
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestMemory_ListPurchases(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	p := agent.LabourPowerPurchase{
		SellerID: agent.NewAgentID(), BuyerID: agent.NewAgentID(),
		WageMinutes: 240, ContractDays: 1,
	}
	if _, err := m.CreatePurchase(ctx, p); err != nil {
		t.Fatalf("CreatePurchase: %v", err)
	}
	list, err := m.ListPurchases(ctx)
	if err != nil {
		t.Fatalf("ListPurchases: %v", err)
	}
	if len(list) != 1 {
		t.Errorf("want 1 purchase, got %d", len(list))
	}
}
