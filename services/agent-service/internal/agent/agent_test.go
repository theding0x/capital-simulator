package agent_test

import (
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

func newCapitalist(balance agent.Pence) agent.Agent {
	return agent.Agent{ID: agent.NewID(), Name: "Capitalist", Class: agent.Capitalist, MoneyBalance: balance}
}

func newMiser(balance agent.Pence) agent.Agent {
	return agent.Agent{ID: agent.NewID(), Name: "Miser", Class: agent.Miser, MoneyBalance: balance}
}

func newWorker(balance agent.Pence) agent.Agent {
	return agent.Agent{ID: agent.NewID(), Name: "Worker", Class: agent.Worker, MoneyBalance: balance}
}

// §1: Advance deducts balance to zero.
func TestAdvance_DeductsToZero(t *testing.T) {
	t.Parallel()
	a := newCapitalist(10000)
	got, err := a.Advance(10000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.MoneyBalance != 0 {
		t.Errorf("want balance 0, got %d", got.MoneyBalance)
	}
}

// §8: SurplusValue == MReturned - MAdvanced.
func TestCapitalCircuit_SurplusValue(t *testing.T) {
	t.Parallel()
	c := agent.CapitalCircuit{
		ID:           agent.NewID(),
		AgentID:      agent.NewID(),
		MAdvanced:    10000,
		CommodityID:  "cotton",
		MReturned:    11000,
		SurplusValue: 1000,
		CircuitType:  agent.CircuitMCM,
	}
	if c.SurplusValue != c.MReturned-c.MAdvanced {
		t.Errorf("want surplus %d, got %d", c.MReturned-c.MAdvanced, c.SurplusValue)
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("valid circuit rejected: %v", err)
	}
}

// §10: Zero surplus circuit is valid but signals no valorisation.
func TestCapitalCircuit_ZeroSurplus(t *testing.T) {
	t.Parallel()
	c := agent.CapitalCircuit{
		ID:           agent.NewID(),
		AgentID:      agent.NewID(),
		MAdvanced:    10000,
		CommodityID:  "cotton",
		MReturned:    10000,
		SurplusValue: 0,
		CircuitType:  agent.CircuitMCM,
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("zero-surplus circuit should be valid: %v", err)
	}
}

// §14: Miser.Hoard() succeeds; balance unchanged; Hoarding flag set.
func TestHoard_MiserSucceeds(t *testing.T) {
	t.Parallel()
	a := newMiser(10000)
	got, err := a.Hoard()
	if err != nil {
		t.Fatalf("Miser.Hoard() error: %v", err)
	}
	if !got.Hoarding {
		t.Error("want Hoarding=true")
	}
	if got.MoneyBalance != a.MoneyBalance {
		t.Errorf("Hoard must not change balance; want %d got %d", a.MoneyBalance, got.MoneyBalance)
	}
}

// §14: Miser.Reinvest() → ErrNotCapitalist.
func TestReinvest_MiserFails(t *testing.T) {
	t.Parallel()
	a := newMiser(10000)
	_, _, err := a.Reinvest("cotton", 11000)
	if !errors.Is(err, agent.ErrNotCapitalist) {
		t.Errorf("Miser.Reinvest() want ErrNotCapitalist, got %v", err)
	}
}

// §14: Capitalist.Reinvest() succeeds; circuit.MAdvanced = old balance.
func TestReinvest_CapitalistSucceeds(t *testing.T) {
	t.Parallel()
	a := newCapitalist(10000)
	circuit, updated, err := a.Reinvest("cotton", 11000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if circuit.MAdvanced != 10000 {
		t.Errorf("want MAdvanced 10000, got %d", circuit.MAdvanced)
	}
	if circuit.SurplusValue != 1000 {
		t.Errorf("want SurplusValue 1000, got %d", circuit.SurplusValue)
	}
	if updated.MoneyBalance != 11000 {
		t.Errorf("want balance 11000 after reinvest, got %d", updated.MoneyBalance)
	}
}

// §14: Capitalist.Hoard() → ErrNotCapitalist.
func TestHoard_CapitalistFails(t *testing.T) {
	t.Parallel()
	a := newCapitalist(10000)
	_, err := a.Hoard()
	if !errors.Is(err, agent.ErrNotCapitalist) {
		t.Errorf("Capitalist.Hoard() want ErrNotCapitalist, got %v", err)
	}
}

// §15: After Realise, second circuit uses full new balance as MAdvanced.
func TestRealise_ExpandingCircuit(t *testing.T) {
	t.Parallel()
	a := newCapitalist(10000)
	c1 := agent.CapitalCircuit{
		ID: agent.NewID(), AgentID: a.ID,
		MAdvanced: 10000, CommodityID: "cotton", MReturned: 11000,
		SurplusValue: 1000, CircuitType: agent.CircuitMCM,
	}
	advanced, err := a.Advance(c1.MAdvanced)
	if err != nil {
		t.Fatalf("Advance: %v", err)
	}
	realised := advanced.Realise(c1)
	if realised.MoneyBalance != 11000 {
		t.Fatalf("want balance 11000 after Realise, got %d", realised.MoneyBalance)
	}
	c2, _, err := realised.Reinvest("cotton", 12100)
	if err != nil {
		t.Fatalf("second Reinvest error: %v", err)
	}
	if c2.MAdvanced != 11000 {
		t.Errorf("second circuit want MAdvanced 11000 (full new balance), got %d", c2.MAdvanced)
	}
}

// Invariant: Advance returns ErrInsufficientFunds when mAdvanced > balance.
func TestAdvance_InsufficientFunds(t *testing.T) {
	t.Parallel()
	a := newCapitalist(5000)
	_, err := a.Advance(10000)
	if !errors.Is(err, agent.ErrInsufficientFunds) {
		t.Errorf("want ErrInsufficientFunds, got %v", err)
	}
}

// Invariant: Miser cannot Advance.
func TestAdvance_MiserCannotAdvance(t *testing.T) {
	t.Parallel()
	a := newMiser(10000)
	_, err := a.Advance(1000)
	if !errors.Is(err, agent.ErrNotCapitalist) {
		t.Errorf("Miser.Advance() want ErrNotCapitalist, got %v", err)
	}
}

// §6: CircuitCMC is valid.
func TestCircuitValidate_CMCIsValid(t *testing.T) {
	t.Parallel()
	c := agent.CapitalCircuit{
		ID: agent.NewID(), AgentID: agent.NewID(),
		MAdvanced: 5000, CommodityID: "bread", MReturned: 5000,
		SurplusValue: 0, CircuitType: agent.CircuitCMC,
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("C-M-C circuit should be valid: %v", err)
	}
}

// Invariant: circuit.MAdvanced <= 0 is invalid.
func TestCircuitValidate_ZeroAdvanced(t *testing.T) {
	t.Parallel()
	c := agent.CapitalCircuit{
		ID: agent.NewID(), AgentID: agent.NewID(),
		MAdvanced: 0, CommodityID: "cotton", MReturned: 0,
		SurplusValue: 0, CircuitType: agent.CircuitMCM,
	}
	if err := c.Validate(); err == nil {
		t.Error("circuit with MAdvanced=0 should fail validation")
	}
}
