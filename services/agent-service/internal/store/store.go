package store

import (
	"context"
	"errors"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

var (
	ErrNotFound      = errors.New("agent: not found")
	ErrAlreadyExists = errors.New("agent: already exists")
)

// Update is a partial-update payload for PATCH /v1/agents/{id}.
// Non-nil fields are applied; nil fields are left untouched.
type Update struct {
	Name         *string
	MoneyBalance *agent.Pence
	Hoarding     *bool
}

func (u Update) IsEmpty() bool {
	return u.Name == nil && u.MoneyBalance == nil && u.Hoarding == nil
}

func (u Update) Apply(a agent.Agent) agent.Agent {
	out := a
	if u.Name != nil {
		out.Name = *u.Name
	}
	if u.MoneyBalance != nil {
		out.MoneyBalance = *u.MoneyBalance
	}
	if u.Hoarding != nil {
		out.Hoarding = *u.Hoarding
	}
	return out
}

// Store is the persistence contract for Agent records.
type Store interface {
	Create(ctx context.Context, a agent.Agent) (agent.Agent, error)
	Get(ctx context.Context, id agent.ID) (agent.Agent, error)
	List(ctx context.Context) ([]agent.Agent, error)
	ListByClass(ctx context.Context, class agent.Class) ([]agent.Agent, error)
	Update(ctx context.Context, id agent.ID, u Update) (agent.Agent, error)
	Delete(ctx context.Context, id agent.ID) error
}

// CircuitStore is the persistence contract for CapitalCircuit records.
// CreateCircuit atomically inserts the circuit and updates the agent balance
// by circuit.SurplusValue.
type CircuitStore interface {
	CreateCircuit(ctx context.Context, c agent.CapitalCircuit) (agent.CapitalCircuit, error)
	GetCircuit(ctx context.Context, id agent.ID) (agent.CapitalCircuit, error)
	ListCircuits(ctx context.Context, agentID agent.ID) ([]agent.CapitalCircuit, error)
}
