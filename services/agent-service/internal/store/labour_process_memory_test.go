package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

func makeLabourProcess() agent.LabourProcess {
	return agent.LabourProcess{
		WorkerID:               agent.NewAgentID(),
		CapitalistID:           agent.NewAgentID(),
		Duration:               720,
		NecessaryLabourMinutes: 360,
		Means: agent.MeansOfProduction{
			RawMaterials: []agent.RawMaterial{
				{CommodityID: "cotton", Quantity: 10, SNLTPerUnit: 120},
			},
			Instruments: []agent.Instrument{
				{CommodityID: "spindle", WearPerRun: 240},
			},
		},
		ProductKind:     "yarn",
		ProductQuantity: 10,
	}
}

func TestMemory_CreateGetLabourProcess(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()
	lp := makeLabourProcess()
	created, err := m.CreateLabourProcess(ctx, lp)
	if err != nil {
		t.Fatalf("CreateLabourProcess: %v", err)
	}
	if created.ID.IsZero() {
		t.Error("CreateLabourProcess should assign an ID")
	}
	got, err := m.GetLabourProcess(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetLabourProcess: %v", err)
	}
	if got.Duration != lp.Duration {
		t.Errorf("duration mismatch: want %d, got %d", lp.Duration, got.Duration)
	}
	if got.NecessaryLabourMinutes != lp.NecessaryLabourMinutes {
		t.Errorf("necessary_labour mismatch: want %d, got %d",
			lp.NecessaryLabourMinutes, got.NecessaryLabourMinutes)
	}
	if got.ProductKind != lp.ProductKind {
		t.Errorf("product_kind mismatch: want %q, got %q", lp.ProductKind, got.ProductKind)
	}
}

func TestMemory_GetLabourProcess_NotFound(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.GetLabourProcess(context.Background(), agent.NewLabourProcessID())
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestMemory_CreateLabourProcess_InvalidDuration(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	lp := makeLabourProcess()
	lp.Duration = 0
	_, err := m.CreateLabourProcess(context.Background(), lp)
	if err == nil {
		t.Error("zero-duration LabourProcess should fail validation")
	}
}
