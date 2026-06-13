package circulation_test

import (
	"errors"
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

// --- OrganicComposition tests -----------------------------------------------

func TestOrganicComposition_Total(t *testing.T) {
	t.Parallel()
	// §Meeting of Demand and Supply: 80c + 20v
	oc := circulation.OrganicComposition{ConstantPence: 8000, VariablePence: 2000}
	if got := oc.Total(); got != 10000 {
		t.Fatalf("Total() = %d, want 10000", got)
	}
}

func TestOrganicComposition_RatioBasisPoints(t *testing.T) {
	t.Parallel()
	// 80c + 20v → ratio = 10000 * 8000 / 2000 = 40000 basis points (4:1)
	oc := circulation.OrganicComposition{ConstantPence: 8000, VariablePence: 2000}
	if got := oc.RatioBasisPoints(); got != 40000 {
		t.Fatalf("RatioBasisPoints() = %d, want 40000", got)
	}
}

func TestOrganicComposition_RatioBasisPoints_ZeroVariable(t *testing.T) {
	t.Parallel()
	oc := circulation.OrganicComposition{ConstantPence: 5000, VariablePence: 0}
	if got := oc.RatioBasisPoints(); got != 0 {
		t.Fatalf("RatioBasisPoints() = %d, want 0 when VariablePence == 0", got)
	}
}

// --- ComputeSupplyDemandImbalance tests -------------------------------------

func TestComputeSupplyDemandImbalance_CanonicalFixture(t *testing.T) {
	t.Parallel()
	// §Meeting of Demand and Supply: 80c + 20v + 20s → demand 10000, supply 12000, excess 2000
	ic := circulation.NewIndustrialCapitalID()
	oc := circulation.OrganicComposition{ConstantPence: 8000, VariablePence: 2000}
	imb := circulation.ComputeSupplyDemandImbalance(ic, "1871-W01", oc, 2000)

	if imb.DemandPence != 10000 {
		t.Errorf("DemandPence = %d, want 10000", imb.DemandPence)
	}
	if imb.SupplyPence != 12000 {
		t.Errorf("SupplyPence = %d, want 12000", imb.SupplyPence)
	}
	if imb.ExcessPence != 2000 {
		t.Errorf("ExcessPence = %d, want 2000", imb.ExcessPence)
	}
	// 5:6 ratio invariant: supply - demand == excess
	if imb.SupplyPence-imb.DemandPence != imb.ExcessPence {
		t.Error("invariant: supply - demand != excess")
	}
}

// --- StageDistribution tests ------------------------------------------------

func TestStageDistribution_Total(t *testing.T) {
	t.Parallel()
	// §intro: Spinning Mill 1871 — three weekly cohorts, total preserved
	sd := circulation.StageDistribution{
		MoneyPence:      3000,
		ProductionPence: 5000,
		CommodityPence:  2000,
	}
	if sd.Total() != 10000 {
		t.Fatalf("Total() = %d, want 10000", sd.Total())
	}
}

func TestStageDistribution_Validate_ConservationFails(t *testing.T) {
	t.Parallel()
	sd := circulation.StageDistribution{
		MoneyPence:      3000,
		ProductionPence: 5000,
		CommodityPence:  1000, // wrong: 9000 != 10000
	}
	if err := sd.Validate(10000); err == nil {
		t.Fatal("expected validation error when total doesn't match")
	}
}

func TestStageDistribution_Validate_NegativePence(t *testing.T) {
	t.Parallel()
	sd := circulation.StageDistribution{MoneyPence: -1}
	if err := sd.Validate(0); err == nil {
		t.Fatal("expected validation error for negative pence")
	}
}

// --- ApplyValueRevolution tests ---------------------------------------------

func TestApplyValueRevolution_TiedUp(t *testing.T) {
	t.Parallel()
	// §value revolution: DirectionPence > 0 → MoneyCapitalTiedUp
	ic := circulation.NewIndustrialCapitalID()
	evt := circulation.ValueRevolutionEvent{
		ID:                  circulation.NewValueRevolutionEventID(),
		IndustrialCapitalID: ic,
		Affecting:           circulation.AffectingMeansOfProduction,
		DirectionPence:      500, // prices rose by 500p
		OccurredAt:          time.Now(),
	}
	res, err := circulation.ApplyValueRevolution(evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.TiedUp == nil {
		t.Fatal("expected MoneyCapitalTiedUp, got nil")
	}
	if res.SetFree != nil {
		t.Fatal("expected SetFree to be nil")
	}
	if res.TiedUp.AmountPence != 500 {
		t.Errorf("AmountPence = %d, want 500", res.TiedUp.AmountPence)
	}
}

func TestApplyValueRevolution_SetFree(t *testing.T) {
	t.Parallel()
	// §value revolution: DirectionPence < 0 → MoneyCapitalSetFree
	ic := circulation.NewIndustrialCapitalID()
	evt := circulation.ValueRevolutionEvent{
		ID:                  circulation.NewValueRevolutionEventID(),
		IndustrialCapitalID: ic,
		Affecting:           circulation.AffectingMeansOfProduction,
		DirectionPence:      -300, // prices fell by 300p
		OccurredAt:          time.Now(),
	}
	res, err := circulation.ApplyValueRevolution(evt)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.SetFree == nil {
		t.Fatal("expected MoneyCapitalSetFree, got nil")
	}
	if res.TiedUp != nil {
		t.Fatal("expected TiedUp to be nil")
	}
	if res.SetFree.AmountPence != 300 {
		t.Errorf("AmountPence = %d, want 300", res.SetFree.AmountPence)
	}
}

func TestApplyValueRevolution_ZeroPence_Rejected(t *testing.T) {
	t.Parallel()
	evt := circulation.ValueRevolutionEvent{
		ID:             circulation.NewValueRevolutionEventID(),
		DirectionPence: 0,
	}
	if _, err := circulation.ApplyValueRevolution(evt); err == nil {
		t.Fatal("expected error for zero DirectionPence")
	}
}

// --- MetamorphosisInterlock tests -------------------------------------------

func TestMetamorphosisInterlock_Validate_CapitalistSeller(t *testing.T) {
	t.Parallel()
	// §middle: spinning mill buys coal from coal-mine capitalist — valid interlock
	mi := circulation.MetamorphosisInterlock{
		ID:                        circulation.NewMetamorphosisInterlockID(),
		BuyerIndustrialCapitalID:  circulation.NewIndustrialCapitalID(),
		SellerIndustrialCapitalID: circulation.NewIndustrialCapitalID(),
		SellerMPSourceOrigin:      circulation.OriginCapitalist,
		Pence:                     3700,
		OccurredAt:                time.Now(),
	}
	if err := mi.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestMetamorphosisInterlock_Validate_NonCapitalSeller_Rejected(t *testing.T) {
	t.Parallel()
	// §middle: spinning mill buys cotton from independent peasant — no interlock
	mi := circulation.MetamorphosisInterlock{
		BuyerIndustrialCapitalID:  circulation.NewIndustrialCapitalID(),
		SellerIndustrialCapitalID: circulation.NewIndustrialCapitalID(),
		SellerMPSourceOrigin:      circulation.OriginIndependentLabourer,
		Pence:                     1000,
	}
	if err := mi.Validate(); err == nil {
		t.Fatal("expected ErrNonCapitalCounterparty, got nil")
	}
}

func TestMetamorphosisInterlock_Validate_NonCapitalistProducer_Rejected(t *testing.T) {
	t.Parallel()
	mi := circulation.MetamorphosisInterlock{
		BuyerIndustrialCapitalID: circulation.NewIndustrialCapitalID(),
		SellerMPSourceOrigin:     circulation.OriginNonCapitalistProducer,
		Pence:                    1000,
	}
	if err := mi.Validate(); err == nil {
		t.Fatal("expected ErrNonCapitalCounterparty, got nil")
	}
}

// --- SinkingFund tests ------------------------------------------------------

func TestSinkingFund_Validate_Valid(t *testing.T) {
	t.Parallel()
	// §sinking fund: £4000 fixed, 10yr → £400/yr
	sf := circulation.SinkingFund{
		ID:                  circulation.NewSinkingFundID(),
		IndustrialCapitalID: circulation.NewIndustrialCapitalID(),
		FixedCapitalPence:   400000,
		LifetimeYears:       10,
		AnnualPaymentPence:  40000,
	}
	if err := sf.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestSinkingFund_Validate_MismatchRejected(t *testing.T) {
	t.Parallel()
	sf := circulation.SinkingFund{
		FixedCapitalPence:  400000,
		LifetimeYears:      10,
		AnnualPaymentPence: 50000, // wrong: 50000*10 != 400000
	}
	if err := sf.Validate(); err == nil {
		t.Fatal("expected ErrSinkingFundMismatch, got nil")
	}
}

func TestSinkingFund_Tick_AccumulatesCorrectly(t *testing.T) {
	t.Parallel()
	// After year 1: AccumulatedPence == 40000; FixedCapitalPence still 400000
	sf := circulation.SinkingFund{
		FixedCapitalPence:  400000,
		LifetimeYears:      10,
		AnnualPaymentPence: 40000,
		AccumulatedPence:   0,
	}
	sf = sf.Tick()
	if sf.AccumulatedPence != 40000 {
		t.Errorf("AccumulatedPence after 1 tick = %d, want 40000", sf.AccumulatedPence)
	}
	if sf.FixedCapitalPence != 400000 {
		t.Errorf("FixedCapitalPence changed unexpectedly to %d", sf.FixedCapitalPence)
	}
}

// --- IndustrialCapital.Validate tests ---------------------------------------

func TestIndustrialCapital_Validate_ZeroTotalPence(t *testing.T) {
	t.Parallel()
	ic := circulation.IndustrialCapital{
		TotalPence:  0,
		EconomyMode: circulation.EconomyMoney,
	}
	if err := ic.Validate(); err == nil {
		t.Fatal("expected error for zero TotalPence")
	}
}

func TestIndustrialCapital_Validate_InvalidEconomyMode(t *testing.T) {
	t.Parallel()
	ic := circulation.IndustrialCapital{
		TotalPence:  10000,
		EconomyMode: circulation.EconomyMode("barter"),
	}
	if err := ic.Validate(); err == nil {
		t.Fatal("expected error for invalid EconomyMode")
	}
}

// --- TransitionEconomy tests (FIX 1: issue #422) ----------------------------

// TestTransitionEconomy_NaturalToMoney checks the valid natural→money path.
// §Ch.4: simple commodity circulation precedes the money economy.
func TestTransitionEconomy_NaturalToMoney(t *testing.T) {
	t.Parallel()
	ic := circulation.IndustrialCapital{
		TotalPence:  42200,
		EconomyMode: circulation.EconomyNatural,
	}
	updated, err := ic.TransitionEconomy(circulation.EconomyMoney)
	if err != nil {
		t.Fatalf("unexpected error transitioning natural→money: %v", err)
	}
	if updated.EconomyMode != circulation.EconomyMoney {
		t.Errorf("EconomyMode = %q, want %q", updated.EconomyMode, circulation.EconomyMoney)
	}
}

// TestTransitionEconomy_MoneyToCredit checks the valid money→credit path.
// §Ch.4: credit economy presupposes the money economy.
func TestTransitionEconomy_MoneyToCredit(t *testing.T) {
	t.Parallel()
	// Spinning Mill 1871: £422 advance in money economy moves to credit economy.
	ic := circulation.IndustrialCapital{
		TotalPence:  42200,
		EconomyMode: circulation.EconomyMoney,
	}
	updated, err := ic.TransitionEconomy(circulation.EconomyCredit)
	if err != nil {
		t.Fatalf("unexpected error transitioning money→credit: %v", err)
	}
	if updated.EconomyMode != circulation.EconomyCredit {
		t.Errorf("EconomyMode = %q, want %q", updated.EconomyMode, circulation.EconomyCredit)
	}
}

// TestTransitionEconomy_NaturalToCredit_Rejected enforces ErrCreditWithoutMoneyEconomy.
// §Ch.4: "The credit economy, instead of abolishing these contradictions, only
// develops a new instrument for their propagation." — presupposes money economy.
func TestTransitionEconomy_NaturalToCredit_Rejected(t *testing.T) {
	t.Parallel()
	ic := circulation.IndustrialCapital{
		TotalPence:  42200,
		EconomyMode: circulation.EconomyNatural,
	}
	_, err := ic.TransitionEconomy(circulation.EconomyCredit)
	if !errors.Is(err, circulation.ErrCreditWithoutMoneyEconomy) {
		t.Fatalf("expected ErrCreditWithoutMoneyEconomy, got %v", err)
	}
}

// TestTransitionEconomy_ZeroModeToCredit_Rejected checks that an
// uninitialised IndustrialCapital (EconomyMode == "") also cannot jump to credit.
func TestTransitionEconomy_ZeroModeToCredit_Rejected(t *testing.T) {
	t.Parallel()
	ic := circulation.IndustrialCapital{TotalPence: 42200}
	_, err := ic.TransitionEconomy(circulation.EconomyCredit)
	if !errors.Is(err, circulation.ErrCreditWithoutMoneyEconomy) {
		t.Fatalf("expected ErrCreditWithoutMoneyEconomy for zero mode, got %v", err)
	}
}

// --- OrganicComposition.LabourDemandPence tests (FIX 2: issue #423 item 3) --

// TestOrganicComposition_LabourDemandPence_SpinningMill1871 checks the invariant
// "demand for labour == v exactly, not one iota greater."
// §Spinning Mill 1871: £372 MP (c) + £50 wages (v) = £422 advance;
// variable capital = 5000 pence; labour demand must equal that exactly.
func TestOrganicComposition_LabourDemandPence_SpinningMill1871(t *testing.T) {
	t.Parallel()
	oc := circulation.OrganicComposition{
		ConstantPence: 37200, // £372 means of production
		VariablePence: 5000,  // £50 wages
	}
	if got := oc.LabourDemandPence(); got != 5000 {
		t.Errorf("LabourDemandPence() = %d, want 5000 (£50 wages)", got)
	}
}

// TestOrganicComposition_LabourDemandPence_EqualToVariable checks that
// LabourDemandPence() == VariablePence for an arbitrary composition.
func TestOrganicComposition_LabourDemandPence_EqualToVariable(t *testing.T) {
	t.Parallel()
	// §Meeting of Demand and Supply: 80c + 20v
	oc := circulation.OrganicComposition{ConstantPence: 8000, VariablePence: 2000}
	if oc.LabourDemandPence() != oc.VariablePence {
		t.Errorf("LabourDemandPence() %d != VariablePence %d", oc.LabourDemandPence(), oc.VariablePence)
	}
}
