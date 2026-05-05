package agent_test

import (
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// §2: LabourPower must have positive capacity to be valid
func TestLabourPower_ZeroCapacity(t *testing.T) {
	t.Parallel()
	w := agent.Worker{
		OwnsLabourPower:       true,
		OwnsCommoditiesToSell: false,
		LabourPower:           agent.LabourPower{CapacityMinutesPerDay: 0},
	}
	if err := w.Validate(); err == nil {
		t.Error("worker with zero capacity should fail validation")
	}
}

// §3: IsFreeLabourer returns false when worker does not own their labour-power
func TestIsFreeLabourer_NotOwner(t *testing.T) {
	t.Parallel()
	w := agent.Worker{OwnsLabourPower: false, OwnsCommoditiesToSell: false}
	if agent.IsFreeLabourer(w) {
		t.Error("worker who does not own labour-power is not a free labourer")
	}
}

// §4: IsFreeLabourer returns false when worker owns other commodities to sell
func TestIsFreeLabourer_OwnsCommodities(t *testing.T) {
	t.Parallel()
	w := agent.Worker{OwnsLabourPower: true, OwnsCommoditiesToSell: true}
	if agent.IsFreeLabourer(w) {
		t.Error("worker who owns commodities is not a free labourer in Marx's double sense")
	}
}

// Invariant: IsFreeLabourer == (OwnsLabourPower && !OwnsCommoditiesToSell)
func TestIsFreeLabourer_DoubleFreedom(t *testing.T) {
	t.Parallel()
	w := agent.Worker{OwnsLabourPower: true, OwnsCommoditiesToSell: false}
	if !agent.IsFreeLabourer(w) {
		t.Error("worker with double freedom should be a free labourer")
	}
}

// §3: ContractDays <= 0 returns ErrInvalidContract
func TestLabourPowerOffering_ZeroContract(t *testing.T) {
	t.Parallel()
	o := agent.LabourPowerOffering{
		OwnerID:               agent.NewAgentID(),
		CapacityMinutesPerDay: 480,
		ContractDays:          0,
		AskingWage:            240,
	}
	if err := o.Validate(); !errors.Is(err, agent.ErrInvalidContract) {
		t.Errorf("ContractDays=0: want ErrInvalidContract, got %v", err)
	}
}

// §3: Negative ContractDays also returns ErrInvalidContract
func TestLabourPowerOffering_NegativeContract(t *testing.T) {
	t.Parallel()
	o := agent.LabourPowerOffering{
		OwnerID:               agent.NewAgentID(),
		CapacityMinutesPerDay: 480,
		ContractDays:          -5,
		AskingWage:            240,
	}
	if err := o.Validate(); !errors.Is(err, agent.ErrInvalidContract) {
		t.Errorf("ContractDays=-5: want ErrInvalidContract, got %v", err)
	}
}

// §5: DailyValue == SubsistenceBasket.TotalSNLT()
func TestLabourPowerValue_DailyValueEqualsBasketTotal(t *testing.T) {
	t.Parallel()
	basket := agent.SubsistenceBasket{
		{Name: "food", SNLTMinutes: 120, Essential: true},
		{Name: "shelter", SNLTMinutes: 120, Essential: false},
	}
	lpv := agent.LabourPowerValue{Basket: basket}
	if lpv.DailyValue() != basket.TotalSNLT() {
		t.Errorf("DailyValue (%d) must equal basket.TotalSNLT() (%d)",
			lpv.DailyValue(), basket.TotalSNLT())
	}
}

// §5: Half-day subsistence cost → DailyValue = 240 labour-minutes
func TestLabourPowerValue_HalfDaySubsistence(t *testing.T) {
	t.Parallel()
	basket := agent.SubsistenceBasket{
		{Name: "subsistence", SNLTMinutes: 240, Essential: true},
	}
	lpv := agent.LabourPowerValue{Basket: basket}
	if got := lpv.DailyValue(); got != agent.LabourMinutes(240) {
		t.Errorf("DailyValue: want 240, got %d", got)
	}
}

// §5: MinimumValue equals SNLT of essential items only
func TestLabourPowerValue_MinimumValue(t *testing.T) {
	t.Parallel()
	basket := agent.SubsistenceBasket{
		{Name: "food", SNLTMinutes: 120, Essential: true},
		{Name: "clothing", SNLTMinutes: 60, Essential: true},
		{Name: "leisure", SNLTMinutes: 60, Essential: false},
	}
	lpv := agent.LabourPowerValue{Basket: basket}
	if got := lpv.MinimumValue(); got != agent.LabourMinutes(180) {
		t.Errorf("MinimumValue: want 180, got %d", got)
	}
}

// Invariant: MinimumValue() <= DailyValue()
func TestLabourPowerValue_MinimumLEDailyValue(t *testing.T) {
	t.Parallel()
	basket := agent.SubsistenceBasket{
		{Name: "food", SNLTMinutes: 120, Essential: true},
		{Name: "tobacco", SNLTMinutes: 60, Essential: false},
	}
	lpv := agent.LabourPowerValue{Basket: basket}
	if lpv.MinimumValue() > lpv.DailyValue() {
		t.Errorf("MinimumValue (%d) must be <= DailyValue (%d)",
			lpv.MinimumValue(), lpv.DailyValue())
	}
}

// ReproductionCost == DailyValue under normal conditions [§5]
func TestLabourPowerValue_ReproductionCost(t *testing.T) {
	t.Parallel()
	basket := agent.SubsistenceBasket{
		{Name: "food", SNLTMinutes: 240, Essential: true},
	}
	lpv := agent.LabourPowerValue{Basket: basket}
	if lpv.ReproductionCost() != lpv.DailyValue() {
		t.Errorf("ReproductionCost (%d) must equal DailyValue (%d)",
			lpv.ReproductionCost(), lpv.DailyValue())
	}
}

// §1: WageMinutes == DailyValue is an exchange of equivalents — purchase is valid
func TestLabourPowerPurchase_FairPrice_Valid(t *testing.T) {
	t.Parallel()
	basket := agent.SubsistenceBasket{
		{Name: "subsistence", SNLTMinutes: 240, Essential: true},
	}
	lpv := agent.LabourPowerValue{Basket: basket}
	p := agent.LabourPowerPurchase{
		SellerID:     agent.NewAgentID(),
		BuyerID:      agent.NewAgentID(),
		WageMinutes:  lpv.DailyValue(),
		ContractDays: 1,
	}
	if err := p.Validate(); err != nil {
		t.Errorf("fair-price purchase should be valid: %v", err)
	}
}

// Invariant: WageMinutes >= 0
func TestLabourPowerPurchase_NegativeWage_Invalid(t *testing.T) {
	t.Parallel()
	p := agent.LabourPowerPurchase{
		SellerID:     agent.NewAgentID(),
		BuyerID:      agent.NewAgentID(),
		WageMinutes:  -1,
		ContractDays: 1,
	}
	if err := p.Validate(); err == nil {
		t.Error("negative wage should fail validation")
	}
}

// NewAgentID produces a non-empty 24-char hex string
func TestNewAgentID(t *testing.T) {
	t.Parallel()
	id := agent.NewAgentID()
	if id.IsZero() {
		t.Error("NewAgentID should not produce zero ID")
	}
	if len(string(id)) != 24 {
		t.Errorf("AgentID should be 24 chars, got %d", len(string(id)))
	}
}

// NewPurchaseID produces a non-empty 24-char hex string
func TestNewPurchaseID(t *testing.T) {
	t.Parallel()
	id := agent.NewPurchaseID()
	if id.IsZero() {
		t.Error("NewPurchaseID should not produce zero ID")
	}
	if len(string(id)) != 24 {
		t.Errorf("PurchaseID should be 24 chars, got %d", len(string(id)))
	}
}
