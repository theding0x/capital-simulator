package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

type ID string
type Class string
type Pence int64
type CircuitType string

const (
	Capitalist Class = "capitalist"
	Worker     Class = "worker"
	Miser      Class = "miser"
	Owner      Class = "owner"
)

const (
	CircuitCMC CircuitType = "C-M-C"
	CircuitMCM CircuitType = "M-C-M-prime"
)

var (
	ErrInsufficientFunds = errors.New("agent: insufficient funds")
	ErrNotCapitalist     = errors.New("agent: operation requires capitalist class")
	ErrWrongClass        = errors.New("agent: operation not permitted for this class")
)

type Agent struct {
	ID            ID        `json:"id"`
	Name          string    `json:"name"`
	Class         Class     `json:"class"`
	MoneyBalance  Pence     `json:"money_balance"`
	LabourMinutes int64     `json:"labour_minutes"`
	Hoarding      bool      `json:"hoarding"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// SurplusValue is always computed as MReturned - MAdvanced; never stored independently.
type CapitalCircuit struct {
	ID           ID          `json:"id"`
	AgentID      ID          `json:"agent_id"`
	MAdvanced    Pence       `json:"m_advanced"`
	CommodityID  string      `json:"commodity_id"`
	MReturned    Pence       `json:"m_returned"`
	SurplusValue Pence       `json:"surplus_value"`
	CircuitType  CircuitType `json:"circuit_type"`
	CreatedAt    time.Time   `json:"created_at"`
}

func (id ID) IsZero() bool { return id == "" }

func NewID() ID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return ID(hex.EncodeToString(b))
}

func (a Agent) Validate() error {
	if a.Name == "" {
		return errors.New("agent: name is required")
	}
	switch a.Class {
	case Capitalist, Worker, Miser, Owner:
	default:
		return errors.New("agent: unknown class")
	}
	if a.MoneyBalance < 0 {
		return errors.New("agent: money_balance cannot be negative")
	}
	return nil
}

func (c CapitalCircuit) Validate() error {
	if c.AgentID.IsZero() {
		return errors.New("circuit: agent_id is required")
	}
	if c.MAdvanced <= 0 {
		return errors.New("circuit: m_advanced must be positive")
	}
	if c.CommodityID == "" {
		return errors.New("circuit: commodity_id is required")
	}
	switch c.CircuitType {
	case CircuitCMC, CircuitMCM:
	default:
		return errors.New("circuit: unknown circuit_type")
	}
	if c.SurplusValue != c.MReturned-c.MAdvanced {
		return errors.New("circuit: surplus_value must equal m_returned - m_advanced")
	}
	return nil
}

// Advance deducts mAdvanced from MoneyBalance. Returns ErrNotCapitalist for
// Miser agents and ErrInsufficientFunds if mAdvanced > MoneyBalance.
func (a Agent) Advance(mAdvanced Pence) (Agent, error) {
	if a.Class == Miser {
		return Agent{}, ErrNotCapitalist
	}
	if mAdvanced > a.MoneyBalance {
		return Agent{}, ErrInsufficientFunds
	}
	out := a
	out.MoneyBalance -= mAdvanced
	return out, nil
}

func (a Agent) Realise(circuit CapitalCircuit) Agent {
	out := a
	out.MoneyBalance += circuit.MReturned
	return out
}

// Reinvest creates a new M-C-M' circuit using the agent's full balance as
// MAdvanced. Valid only for Capitalist agents.
func (a Agent) Reinvest(commodityID string, mReturned Pence) (CapitalCircuit, Agent, error) {
	if a.Class != Capitalist {
		return CapitalCircuit{}, Agent{}, ErrNotCapitalist
	}
	if a.MoneyBalance <= 0 {
		return CapitalCircuit{}, Agent{}, ErrInsufficientFunds
	}
	circuit := CapitalCircuit{
		ID:           NewID(),
		AgentID:      a.ID,
		MAdvanced:    a.MoneyBalance,
		CommodityID:  commodityID,
		MReturned:    mReturned,
		SurplusValue: mReturned - a.MoneyBalance,
		CircuitType:  CircuitMCM,
		CreatedAt:    time.Now().UTC(),
	}
	if err := circuit.Validate(); err != nil {
		return CapitalCircuit{}, Agent{}, err
	}
	advanced, err := a.Advance(circuit.MAdvanced)
	if err != nil {
		return CapitalCircuit{}, Agent{}, err
	}
	realised := advanced.Realise(circuit)
	return circuit, realised, nil
}

// Hoard sets the Hoarding flag on a Miser agent. Returns ErrNotCapitalist for
// non-Miser agents. Balance is unaffected (idempotent no-op on balance).
func (a Agent) Hoard() (Agent, error) {
	if a.Class != Miser {
		return Agent{}, ErrNotCapitalist
	}
	out := a
	out.Hoarding = true
	return out, nil
}
