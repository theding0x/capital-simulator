package circulation_test

import (
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

// Fixtures drawn from Vol. II Ch. 1 §I–§III (Marx's Spinning Mill 1871 examples).

// TestPurchasePhaseValidate_SpinningMill1871 checks the proportional law from §I.
// "weekly wage of its 50 labourers amounts to £50, £372 must be spent for means
// of production" — 3,000 total labour-hours with 3,000 capacity-hours in MP.
func TestPurchasePhaseValidate_SpinningMill1871(t *testing.T) {
	t.Parallel()
	p := circulation.PurchasePhase{
		LabourLeg: circulation.LabourLeg{Amount: 5000, LabourPowerHours: 3000},
		MeansLeg:  circulation.MeansLeg{Amount: 37200, MeansCapacityHours: 3000},
	}
	if err := p.Validate(); err != nil {
		t.Fatalf("valid purchase phase rejected: %v", err)
	}
}

// TestPurchasePhaseValidate_InsufficientMeans checks §I:
// "If the means of production at hand were insufficient, the excess labour at
// the disposal of the purchaser could not be utilised."
func TestPurchasePhaseValidate_InsufficientMeans(t *testing.T) {
	t.Parallel()
	p := circulation.PurchasePhase{
		LabourLeg: circulation.LabourLeg{Amount: 5000, LabourPowerHours: 3000},
		MeansLeg:  circulation.MeansLeg{Amount: 10000, MeansCapacityHours: 1000},
	}
	if err := p.Validate(); !errors.Is(err, circulation.ErrLabourMeansMismatch) {
		t.Fatalf("expected ErrLabourMeansMismatch, got %v", err)
	}
}

// TestPurchasePhaseValidate_ExcessMeans checks §I:
// "If there were more means of production than available labour, they would not
// be saturated with labour, would not be transformed into products."
func TestPurchasePhaseValidate_ExcessMeans(t *testing.T) {
	t.Parallel()
	p := circulation.PurchasePhase{
		LabourLeg: circulation.LabourLeg{Amount: 5000, LabourPowerHours: 3000},
		MeansLeg:  circulation.MeansLeg{Amount: 80000, MeansCapacityHours: 8000},
	}
	if err := p.Validate(); !errors.Is(err, circulation.ErrLabourMeansMismatch) {
		t.Fatalf("expected ErrLabourMeansMismatch, got %v", err)
	}
}

// TestPurchasePhaseTotal verifies LabourLeg.Amount + MeansLeg.Amount = advance total.
func TestPurchasePhaseTotal(t *testing.T) {
	t.Parallel()
	p := circulation.PurchasePhase{
		LabourLeg: circulation.LabourLeg{Amount: 5000, LabourPowerHours: 3000},
		MeansLeg:  circulation.MeansLeg{Amount: 37200, MeansCapacityHours: 3000},
	}
	if got := p.Total(); got != 42200 {
		t.Fatalf("expected 42200 (£422), got %d", got)
	}
}

// TestProductiveStateTotalAdvance checks form-change magnitude preservation §I.
// "M is the same capital-value as P, only it has a different mode of existence."
func TestProductiveStateTotalAdvance(t *testing.T) {
	t.Parallel()
	ps := circulation.ProductiveState{
		ConstantPence: 37200, // £372 means of production
		VariablePence: 5000,  // £50 wages
	}
	if got := ps.TotalAdvance(); got != 42200 {
		t.Fatalf("expected 42200 (£422), got %d", got)
	}
}

// TestCommodityCapitalTotal checks §II:
// "The 10,000 lbs. of yarn therefore represent a value of £500."
// C′ = C (£422 advance) + c (£78 surplus) = £500.
func TestCommodityCapitalTotal(t *testing.T) {
	t.Parallel()
	cc := circulation.CommodityCapital{
		ValueOriginal: 42200, // £422 advance
		ValueSurplus:  7800,  // £78 surplus
	}
	if got := cc.Total(); got != 50000 {
		t.Fatalf("expected 50000 (£500), got %d", got)
	}
}

// TestMoneyCircuitSurplusRealised_156pct checks §II: Spinning Mill 1871 at 156%.
// Advance £422, realised £500 → surplus = £78. Rate = 78/50 = 156%.
func TestMoneyCircuitSurplusRealised_156pct(t *testing.T) {
	t.Parallel()
	mc := circulation.MoneyCircuit{
		Advance:     circulation.MoneyAdvance{Amount: 42200},
		Realisation: &circulation.Realisation{RealisedPence: 50000},
	}
	if got := mc.SurplusRealised(); got != 7800 {
		t.Fatalf("expected 7800 (£78), got %d", got)
	}
}

// TestMoneyCircuitSurplusRealised_100pct checks §III: Spinning Mill 1871 at 100%.
// "If he had paid his labourers £64 in wages instead of £50 his surplus-value
// would only be £64 instead of £78." Advance = £372 + £64 = £436.
func TestMoneyCircuitSurplusRealised_100pct(t *testing.T) {
	t.Parallel()
	mc := circulation.MoneyCircuit{
		Advance:     circulation.MoneyAdvance{Amount: 43600}, // £372 MP + £64 wages
		Realisation: &circulation.Realisation{RealisedPence: 50000},
	}
	if got := mc.SurplusRealised(); got != 6400 {
		t.Fatalf("expected 6400 (£64), got %d", got)
	}
}

// TestMoneyCircuitSurplusRealised_BeforeRealisation verifies 0 before M′.
func TestMoneyCircuitSurplusRealised_BeforeRealisation(t *testing.T) {
	t.Parallel()
	mc := circulation.MoneyCircuit{
		Advance: circulation.MoneyAdvance{Amount: 42200},
	}
	if got := mc.SurplusRealised(); got != 0 {
		t.Fatalf("expected 0 before realisation, got %d", got)
	}
}

// TestMoneyCircuitSurplusRealised_PartialRealisation checks §III partial cases.
// "If the capitalist succeeds in selling only 7,440 lbs. at their value of £372,
// he has replaced only the value of his constant capital."
// RealisedPence=37200 → surplus = 37200 - 42200 = -5000 (loss).
func TestMoneyCircuitSurplusRealised_PartialRealisation(t *testing.T) {
	t.Parallel()
	// Only constant capital recovered; variable and surplus both lost.
	mc := circulation.MoneyCircuit{
		Advance:     circulation.MoneyAdvance{Amount: 42200},
		Realisation: &circulation.Realisation{RealisedPence: 37200},
	}
	if got := mc.SurplusRealised(); got != -5000 {
		t.Fatalf("expected -5000 (loss), got %d", got)
	}
}

// TestMoneyCircuitSurplusRealised_BreakEven checks break-even (simple reproduction limit).
// RealisedPence == advance → surplus = 0.
func TestMoneyCircuitSurplusRealised_BreakEven(t *testing.T) {
	t.Parallel()
	mc := circulation.MoneyCircuit{
		Advance:     circulation.MoneyAdvance{Amount: 42200},
		Realisation: &circulation.Realisation{RealisedPence: 42200},
	}
	if got := mc.SurplusRealised(); got != 0 {
		t.Fatalf("expected 0 (break-even), got %d", got)
	}
}

// TestMoneyCircuitCanTransitionTo_HappyPath checks the full linear sequence M→M′.
func TestMoneyCircuitCanTransitionTo_HappyPath(t *testing.T) {
	t.Parallel()
	pairs := []struct {
		from circulation.CircuitMoment
		to   circulation.CircuitMoment
	}{
		{circulation.MomentM, circulation.MomentMC},
		{circulation.MomentMC, circulation.MomentP},
		{circulation.MomentP, circulation.MomentCPrime},
		{circulation.MomentCPrime, circulation.MomentCMPrime},
		{circulation.MomentCMPrime, circulation.MomentMPrime},
	}
	for _, tc := range pairs {
		mc := circulation.MoneyCircuit{Moment: tc.from}
		if err := mc.CanTransitionTo(tc.to); err != nil {
			t.Errorf("expected valid transition %s→%s, got %v", tc.from, tc.to, err)
		}
	}
}

// TestMoneyCircuitCanTransitionTo_SkipP verifies ErrInvalidPhaseTransition on skip.
func TestMoneyCircuitCanTransitionTo_SkipP(t *testing.T) {
	t.Parallel()
	mc := circulation.MoneyCircuit{Moment: circulation.MomentMC}
	// Attempting to jump directly to C-prime without going through P.
	if err := mc.CanTransitionTo(circulation.MomentCPrime); !errors.Is(err, circulation.ErrInvalidPhaseTransition) {
		t.Fatalf("expected ErrInvalidPhaseTransition, got %v", err)
	}
}

// TestMoneyCircuitCanTransitionTo_ReEntry verifies ErrInvalidPhaseTransition on re-entry.
func TestMoneyCircuitCanTransitionTo_ReEntry(t *testing.T) {
	t.Parallel()
	mc := circulation.MoneyCircuit{Moment: circulation.MomentP}
	if err := mc.CanTransitionTo(circulation.MomentMC); !errors.Is(err, circulation.ErrInvalidPhaseTransition) {
		t.Fatalf("expected ErrInvalidPhaseTransition, got %v", err)
	}
}

// TestNewMoneyCircuitID checks uniqueness and length.
func TestNewMoneyCircuitID(t *testing.T) {
	t.Parallel()
	a := circulation.NewMoneyCircuitID()
	b := circulation.NewMoneyCircuitID()
	if a == b {
		t.Fatal("two IDs must not be equal")
	}
	if len(string(a)) != 24 {
		t.Fatalf("expected 24 hex chars (96 bits), got %d", len(string(a)))
	}
}
