package store_test

import (
	"context"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/store"
)

func TestMemory_WorkingDay_RoundTrip(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()

	wd := agent.WorkingDay{
		NecessaryLabourMinutes: 360,
		SurplusLabourMinutes:   360,
	}
	created, err := m.CreateWorkingDay(ctx, wd)
	if err != nil {
		t.Fatalf("CreateWorkingDay: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected non-zero ID")
	}
	if created.CreatedAt.IsZero() {
		t.Fatal("expected non-zero CreatedAt")
	}
	if created.NecessaryLabourMinutes != 360 || created.SurplusLabourMinutes != 360 {
		t.Errorf("fields not preserved: %+v", created)
	}

	got, err := m.GetWorkingDay(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetWorkingDay: %v", err)
	}
	if got.ID != created.ID {
		t.Errorf("ID mismatch: got %q want %q", got.ID, created.ID)
	}
}

func TestMemory_WorkingDay_NotFound(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.GetWorkingDay(context.Background(), "nonexistent")
	if err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMemory_WorkingDay_InvalidRejected(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	// NecessaryLabourMinutes = 0 is invalid
	_, err := m.CreateWorkingDay(context.Background(), agent.WorkingDay{SurplusLabourMinutes: 60})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestMemory_RelaySchedule_RoundTrip(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	ctx := context.Background()

	rs := agent.RelaySchedule{
		Sets: [2]agent.RelaySet{
			{
				ShiftKind:  agent.ShiftDay,
				WorkingDay: agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 360},
				WorkerIDs:  []agent.AgentID{"worker-a"},
			},
			{
				ShiftKind:  agent.ShiftNight,
				WorkingDay: agent.WorkingDay{NecessaryLabourMinutes: 360, SurplusLabourMinutes: 0},
				WorkerIDs:  []agent.AgentID{"worker-b"},
			},
		},
	}
	created, err := m.CreateRelaySchedule(ctx, rs)
	if err != nil {
		t.Fatalf("CreateRelaySchedule: %v", err)
	}
	if created.ID.IsZero() {
		t.Fatal("expected non-zero ID")
	}

	got, err := m.GetRelaySchedule(ctx, created.ID)
	if err != nil {
		t.Fatalf("GetRelaySchedule: %v", err)
	}
	if got.Sets[0].ShiftKind != agent.ShiftDay {
		t.Errorf("Sets[0].ShiftKind = %q, want %q", got.Sets[0].ShiftKind, agent.ShiftDay)
	}
	if len(got.Sets[0].WorkerIDs) != 1 || got.Sets[0].WorkerIDs[0] != "worker-a" {
		t.Errorf("Sets[0].WorkerIDs = %v, want [worker-a]", got.Sets[0].WorkerIDs)
	}
}

func TestMemory_RelaySchedule_NotFound(t *testing.T) {
	t.Parallel()
	m := store.NewMemory()
	_, err := m.GetRelaySchedule(context.Background(), "nonexistent")
	if err != store.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
