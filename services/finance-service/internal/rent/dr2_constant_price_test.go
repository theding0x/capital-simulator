// Package rent — unit tests for Vol. III Ch. 41 Differential Rent II,
// First Case (Constant Price of Production). Fixtures drawn from Marx
// Vol. III Ch. 41 Tables II & III: Soil D, two successive capital tranches
// each £2.50 (250bp) with proportional (900bp each) or decreasing (900/450bp)
// surplus-profit.
package rent

import (
	"testing"
)

// ── CapitalProductivityCase.IsValid ──────────────────────────────────────────

func TestCapitalProductivityCaseIsValid(t *testing.T) {
	t.Parallel()

	valid := []CapitalProductivityCase{
		CaseAverageOnly,
		CaseProportional,
		CaseDecreasingRate,
		CaseIncreasingReturn,
	}
	for _, c := range valid {
		if !c.IsValid() {
			t.Errorf("case %d should be valid", c)
		}
	}

	invalid := []CapitalProductivityCase{0, 5}
	for _, c := range invalid {
		if c.IsValid() {
			t.Errorf("case %d should be invalid", c)
		}
	}
}

// ── ComputeDR2Outcome — Case II (Proportional) ────────────────────────────────

// TestComputeDR2Outcome_CaseII mirrors Marx Table II, Soil D: two tranches
// each 250bp capital, 900bp surplus-profit → TotalRentalBP=1800,
// AdditionalCapitalBP=250 (tranche 2 only), RateOfSurplusProfitBP=36000.
func TestComputeDR2Outcome_CaseII(t *testing.T) {
	t.Parallel()

	invs := []SuccessiveInvestment{
		{TrancheNumber: 1, CapitalBP: 250, SurplusProfitBP: 900, ParcelID: "PD"},
		{TrancheNumber: 2, CapitalBP: 250, SurplusProfitBP: 900, ParcelID: "PD"},
	}
	out := ComputeDR2Outcome(CaseProportional, invs)

	if out.TotalRentalBP != 1800 {
		t.Errorf("TotalRentalBP = %d, want 1800", out.TotalRentalBP)
	}
	if out.AdditionalCapitalBP != 250 {
		t.Errorf("AdditionalCapitalBP = %d, want 250 (tranche-2 only)", out.AdditionalCapitalBP)
	}
	if out.NewRentPerAcreBP != 1800 {
		t.Errorf("NewRentPerAcreBP = %d, want 1800", out.NewRentPerAcreBP)
	}
	if out.RateOfSurplusProfitBP != 36000 {
		t.Errorf("RateOfSurplusProfitBP = %d, want 36000", out.RateOfSurplusProfitBP)
	}
	if out.ParcelID != "PD" {
		t.Errorf("ParcelID = %q, want %q", out.ParcelID, "PD")
	}
	if out.Case != CaseProportional {
		t.Errorf("Case = %d, want CaseProportional (%d)", out.Case, CaseProportional)
	}
}

// ── ComputeDR2Outcome — Case III (Decreasing Rate) ────────────────────────────

// TestComputeDR2Outcome_CaseIII_DecreasingRate mirrors Marx Table III, Soil D:
// tranche 1 surplus=900bp, tranche 2 surplus=450bp.
// TotalRentalBP=1350, RateOfSurplusProfitBP=27000.
// Absolute rent rose (1350 > 900) though the rate fell (27000 < 36000).
func TestComputeDR2Outcome_CaseIII_DecreasingRate(t *testing.T) {
	t.Parallel()

	invs := []SuccessiveInvestment{
		{TrancheNumber: 1, CapitalBP: 250, SurplusProfitBP: 900, ParcelID: "PD"},
		{TrancheNumber: 2, CapitalBP: 250, SurplusProfitBP: 450, ParcelID: "PD"},
	}
	out := ComputeDR2Outcome(CaseDecreasingRate, invs)

	if out.TotalRentalBP != 1350 {
		t.Errorf("TotalRentalBP = %d, want 1350", out.TotalRentalBP)
	}
	if out.RateOfSurplusProfitBP != 27000 {
		t.Errorf("RateOfSurplusProfitBP = %d, want 27000", out.RateOfSurplusProfitBP)
	}
	// Absolute rent rose even though rate fell.
	if !(out.RateOfSurplusProfitBP < 36000) {
		t.Errorf("rate %d should be < 36000 (Case II rate)", out.RateOfSurplusProfitBP)
	}
	if !(out.TotalRentalBP > 900) {
		t.Errorf("total rental %d should be > 900 (single-tranche base)", out.TotalRentalBP)
	}
}

// ── ComputeDR2Outcome — AdditionalCapital only from tranche ≥2 ───────────────

// TestComputeDR2Outcome_AdditionalCapitalOnlyTranche2Plus verifies that a
// single tranche-1 investment contributes zero AdditionalCapitalBP.
func TestComputeDR2Outcome_AdditionalCapitalOnlyTranche2Plus(t *testing.T) {
	t.Parallel()

	invs := []SuccessiveInvestment{
		{TrancheNumber: 1, CapitalBP: 250, SurplusProfitBP: 900, ParcelID: "PD"},
	}
	out := ComputeDR2Outcome(CaseProportional, invs)

	if out.AdditionalCapitalBP != 0 {
		t.Errorf("AdditionalCapitalBP = %d, want 0 for single tranche-1", out.AdditionalCapitalBP)
	}
}

// ── ComputeDR2Outcome — empty input ──────────────────────────────────────────

func TestComputeDR2Outcome_Empty(t *testing.T) {
	t.Parallel()

	out := ComputeDR2Outcome(CaseProportional, nil)

	if out.TotalRentalBP != 0 {
		t.Errorf("TotalRentalBP = %d, want 0 for empty input", out.TotalRentalBP)
	}
	if out.AdditionalCapitalBP != 0 {
		t.Errorf("AdditionalCapitalBP = %d, want 0 for empty input", out.AdditionalCapitalBP)
	}
	if out.NewRentPerAcreBP != 0 {
		t.Errorf("NewRentPerAcreBP = %d, want 0 for empty input", out.NewRentPerAcreBP)
	}
	if out.RateOfSurplusProfitBP != 0 {
		t.Errorf("RateOfSurplusProfitBP = %d, want 0 for empty input", out.RateOfSurplusProfitBP)
	}
	if out.Case != CaseProportional {
		t.Errorf("Case = %d, want CaseProportional (%d)", out.Case, CaseProportional)
	}
	if out.ParcelID != "" {
		t.Errorf("ParcelID = %q, want empty string for empty input", out.ParcelID)
	}

	// Empty slice is equivalent to nil.
	out2 := ComputeDR2Outcome(CaseDecreasingRate, []SuccessiveInvestment{})
	if out2.TotalRentalBP != 0 {
		t.Errorf("TotalRentalBP = %d, want 0 for empty slice", out2.TotalRentalBP)
	}
}

// ── ComputeDR2Outcome — divide-by-zero guard ─────────────────────────────────

// TestComputeDR2Outcome_RateGuard verifies that zero total capital yields
// RateOfSurplusProfitBP==0 (no panic).
func TestComputeDR2Outcome_RateGuard(t *testing.T) {
	t.Parallel()

	invs := []SuccessiveInvestment{
		{TrancheNumber: 1, CapitalBP: 0, SurplusProfitBP: 0, ParcelID: "PD"},
		{TrancheNumber: 2, CapitalBP: 0, SurplusProfitBP: 0, ParcelID: "PD"},
	}
	out := ComputeDR2Outcome(CaseProportional, invs)

	if out.RateOfSurplusProfitBP != 0 {
		t.Errorf("RateOfSurplusProfitBP = %d, want 0 when total capital is zero", out.RateOfSurplusProfitBP)
	}
}

// ── ID generators ─────────────────────────────────────────────────────────────

// TestCh41NewIDsDistinct verifies that both Ch. 41 ID generators return
// 24-character hex strings and that two successive calls produce distinct values.
func TestCh41NewIDsDistinct(t *testing.T) {
	t.Parallel()

	id1 := NewDifferentialRentIIOutcomeID()
	id2 := NewDifferentialRentIIOutcomeID()
	if len(string(id1)) != 24 {
		t.Errorf("DifferentialRentIIOutcomeID len = %d, want 24", len(string(id1)))
	}
	if id1 == id2 {
		t.Error("two DifferentialRentIIOutcomeIDs must be distinct")
	}

	cid1 := NewIntensiveExtensiveComparisonID()
	cid2 := NewIntensiveExtensiveComparisonID()
	if len(string(cid1)) != 24 {
		t.Errorf("IntensiveExtensiveComparisonID len = %d, want 24", len(string(cid1)))
	}
	if cid1 == cid2 {
		t.Error("two IntensiveExtensiveComparisonIDs must be distinct")
	}

	// IDs from different constructors must also be distinct.
	if DifferentialRentIIOutcomeID(cid1) == DifferentialRentIIOutcomeID(id1) {
		t.Error("IDs from different constructors must be distinct")
	}
}
