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
	Name          *string
	MoneyBalance  *agent.Pence
	Hoarding      *bool
	LabourMinutes *int64
}

func (u Update) IsEmpty() bool {
	return u.Name == nil && u.MoneyBalance == nil && u.Hoarding == nil && u.LabourMinutes == nil
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
	if u.LabourMinutes != nil {
		out.LabourMinutes = *u.LabourMinutes
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

// LabourPowerStore is the persistence contract for Ch. 6 labour-power domain objects.
type LabourPowerStore interface {
	CreateWorker(ctx context.Context, w agent.Worker) (agent.Worker, error)
	GetWorker(ctx context.Context, id agent.AgentID) (agent.Worker, error)
	ListWorkers(ctx context.Context) ([]agent.Worker, error)
	CreateCapitalist(ctx context.Context, c agent.Capitalist) (agent.Capitalist, error)
	GetCapitalist(ctx context.Context, id agent.AgentID) (agent.Capitalist, error)
	ListCapitalists(ctx context.Context) ([]agent.Capitalist, error)
	CreateOffering(ctx context.Context, o agent.LabourPowerOffering) (agent.LabourPowerOffering, error)
	ListOfferings(ctx context.Context) ([]agent.LabourPowerOffering, error)
	CreatePurchase(ctx context.Context, p agent.LabourPowerPurchase) (agent.LabourPowerPurchase, error)
	GetPurchase(ctx context.Context, id agent.PurchaseID) (agent.LabourPowerPurchase, error)
	ListPurchases(ctx context.Context) ([]agent.LabourPowerPurchase, error)
}

// LabourProcessStore is the persistence contract for Ch. 7 labour process records.
type LabourProcessStore interface {
	CreateLabourProcess(ctx context.Context, lp agent.LabourProcess) (agent.LabourProcess, error)
	GetLabourProcess(ctx context.Context, id agent.LabourProcessID) (agent.LabourProcess, error)
}

// WorkingDayStore is the persistence contract for Ch. 10 working-day records.
type WorkingDayStore interface {
	CreateWorkingDay(ctx context.Context, wd agent.WorkingDay) (agent.WorkingDay, error)
	GetWorkingDay(ctx context.Context, id agent.WorkingDayID) (agent.WorkingDay, error)
	CreateRelaySchedule(ctx context.Context, rs agent.RelaySchedule) (agent.RelaySchedule, error)
	GetRelaySchedule(ctx context.Context, id agent.RelayScheduleID) (agent.RelaySchedule, error)
}

// CooperationStore is the persistence contract for Ch. 13 cooperation records.
type CooperationStore interface {
	CreateCooperation(ctx context.Context, c agent.Cooperation) (agent.Cooperation, error)
	GetCooperation(ctx context.Context, id agent.CooperationID) (agent.Cooperation, error)
	ListCooperations(ctx context.Context) ([]agent.Cooperation, error)
	ListCooperationsByCapitalist(ctx context.Context, capitalistID agent.AgentID) ([]agent.Cooperation, error)
}

// ManufactureStore is the persistence contract for Ch. 14 manufacture records.
type ManufactureStore interface {
	CreateManufacture(ctx context.Context, m agent.Manufacture) (agent.Manufacture, error)
	GetManufacture(ctx context.Context, id agent.ManufactureID) (agent.Manufacture, error)
	ListManufactures(ctx context.Context) ([]agent.Manufacture, error)
	ListManufacturesByCapitalist(ctx context.Context, capitalistID agent.AgentID) ([]agent.Manufacture, error)
	ListManufacturesByForm(ctx context.Context, form agent.ManufactureForm) ([]agent.Manufacture, error)
}

// WageFormStore is the persistence contract for Ch. 19 wage-form records.
type WageFormStore interface {
	CreateWageForm(ctx context.Context, wf agent.WageForm) (agent.WageForm, error)
	GetWageForm(ctx context.Context, agentID agent.AgentID) (agent.WageForm, error)
}

// TimeWageStore is the persistence contract for Ch. 20 time-wage session records.
type TimeWageStore interface {
	CreateWorkingSession(ctx context.Context, s agent.WorkingSession) (agent.WorkingSession, error)
	GetWorkingSession(ctx context.Context, id agent.WorkingSessionID) (agent.WorkingSession, error)
}

// PieceWageStore is the persistence contract for Ch. 21 piece-wage records.
type PieceWageStore interface {
	CreatePieceWage(ctx context.Context, pw agent.PieceWage) (agent.PieceWage, error)
	GetPieceWage(ctx context.Context, agentID agent.AgentID) (agent.PieceWage, error)
	CreateSubContract(ctx context.Context, sc agent.SubContract) (agent.SubContract, error)
	GetSubContract(ctx context.Context, id agent.SubContractID) (agent.SubContract, error)
}
