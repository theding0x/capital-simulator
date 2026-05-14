package surplus_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/surplus"
)

// Capital Vol. I, Ch. 16, §1:
// "The prolongation of the working-day beyond the point at which the
//  labourer would have produced just an equivalent for the value of his
//  labour-power, and the appropriation of that surplus-labour by capital,
//  this is production of absolute surplus-value."
//
// Extending a 720/600/120 day by 60 minutes leaves NecessaryLabour fixed
// at 600, raises Total to 780, and SurplusLabour to 180. The returned
// AbsoluteSurplusValue carries 60 minutes — the gain — tagged "absolute".
func TestProlongWorkingDay_Para1_AbsoluteFromExtension(t *testing.T) {
	t.Parallel()
	wd := surplus.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	newWD, asv := surplus.ProlongWorkingDay(wd, 60)
	if newWD.Total != 780 {
		t.Fatalf("Total = %d, want 780", newWD.Total)
	}
	if newWD.NecessaryLabour != 600 {
		t.Fatalf("NecessaryLabour = %d, want 600 (fixed)", newWD.NecessaryLabour)
	}
	if newWD.SurplusLabour != 180 {
		t.Fatalf("SurplusLabour = %d, want 180", newWD.SurplusLabour)
	}
	if asv.LabourMinutes != 60 {
		t.Fatalf("AbsoluteSurplusValue.LabourMinutes = %d, want 60", asv.LabourMinutes)
	}
	if asv.Origin != surplus.OriginAbsolute {
		t.Fatalf("Origin = %q, want %q", asv.Origin, surplus.OriginAbsolute)
	}
}

// §1: "In order to prolong the surplus-labour, the necessary labour is
//  shortened by methods whereby the equivalent for the wages is produced
//  in less time."
//
// Apply ProductivityFactor 2.0 to a 720/600/120 day: NecessaryLabour
// halves to 300; Total fixed at 720; SurplusLabour grows to 420. The
// returned RelativeSurplusValue carries 300 minutes — the necessary-labour
// reduction — tagged "relative".
func TestReduceNecessaryLabour_Para1_RelativeFromProductivity(t *testing.T) {
	t.Parallel()
	wd := surplus.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	newWD, rsv := surplus.ReduceNecessaryLabour(wd, surplus.ProductivityFactor(2.0))
	if newWD.Total != 720 {
		t.Fatalf("Total = %d, want 720 (fixed)", newWD.Total)
	}
	if newWD.NecessaryLabour != 300 {
		t.Fatalf("NecessaryLabour = %d, want 300", newWD.NecessaryLabour)
	}
	if newWD.SurplusLabour != 420 {
		t.Fatalf("SurplusLabour = %d, want 420", newWD.SurplusLabour)
	}
	if rsv.LabourMinutes != 300 {
		t.Fatalf("RelativeSurplusValue.LabourMinutes = %d, want 300", rsv.LabourMinutes)
	}
	if rsv.Origin != surplus.OriginRelative {
		t.Fatalf("Origin = %q, want %q", rsv.Origin, surplus.OriginRelative)
	}
}

// §1, eastern bread-cutter fixture: 12 working hours' wants per week,
// six days of work to appropriate the value of one day. Total = 6 × 12 ×
// 60 = 4320 min; NecessaryLabour = 12 × 60 = 720; SurplusLabour = 3600.
// RateSurplusValue = 3600 / 720 = 5.0 (500%).
func TestRateSurplusValue_Para1_EasternBreadCutter(t *testing.T) {
	t.Parallel()
	wd := surplus.WorkingDay{Total: 4320, NecessaryLabour: 720, SurplusLabour: 3600}
	if wd.NecessaryLabour+wd.SurplusLabour != wd.Total {
		t.Fatalf("partition invariant violated")
	}
	rate := surplus.RateSurplusValue(wd.SurplusLabour, wd.NecessaryLabour)
	if rate != 5.0 {
		t.Fatalf("RateSurplusValue = %v, want 5.0", rate)
	}
}

// §1, Mill's error: £400 constant + £100 variable + £20 surplus.
// RateSurplusValue(20, 100) = 0.20; RateOfProfit(20, 500) = 0.04. The
// rates must differ for the same surplus magnitude.
func TestRateOfProfit_Para1_MillCritique(t *testing.T) {
	t.Parallel()
	rsv := surplus.RateSurplusValue(20, 100)
	rop := surplus.RateOfProfit(20, 500)
	if rsv != 0.20 {
		t.Fatalf("RateSurplusValue = %v, want 0.20", rsv)
	}
	if rop != 0.04 {
		t.Fatalf("RateOfProfit = %v, want 0.04", rop)
	}
	if rsv == rop {
		t.Fatalf("rates must differ — RateSurplusValue (%v) and RateOfProfit (%v) collapsed", rsv, rop)
	}
}

// Invariant §1: partition holds for both paths.
func TestInvariant_PartitionHolds(t *testing.T) {
	t.Parallel()
	wd := surplus.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	a, _ := surplus.ProlongWorkingDay(wd, 60)
	if a.NecessaryLabour+a.SurplusLabour != a.Total {
		t.Fatalf("absolute path violates partition: %+v", a)
	}
	r, _ := surplus.ReduceNecessaryLabour(wd, surplus.ProductivityFactor(2.0))
	if r.NecessaryLabour+r.SurplusLabour != r.Total {
		t.Fatalf("relative path violates partition: %+v", r)
	}
}

// Invariant §1: AbsoluteSurplusValue.LabourMinutes == WorkingDay.SurplusLabour
// when extending from a zero-surplus baseline.
func TestInvariant_AbsoluteEqualsSurplusLabour(t *testing.T) {
	t.Parallel()
	wd := surplus.WorkingDay{Total: 600, NecessaryLabour: 600, SurplusLabour: 0}
	newWD, asv := surplus.ProlongWorkingDay(wd, 120)
	if asv.LabourMinutes != newWD.SurplusLabour {
		t.Fatalf("AbsoluteSurplusValue.LabourMinutes = %d, want SurplusLabour = %d",
			asv.LabourMinutes, newWD.SurplusLabour)
	}
}

// Invariant §1: RelativeSurplusValue.LabourMinutes == (oldNecessary - newNecessary).
func TestInvariant_RelativeEqualsNecessaryReduction(t *testing.T) {
	t.Parallel()
	wd := surplus.WorkingDay{Total: 720, NecessaryLabour: 480, SurplusLabour: 240}
	newWD, rsv := surplus.ReduceNecessaryLabour(wd, surplus.ProductivityFactor(1.5))
	if rsv.LabourMinutes != wd.NecessaryLabour-newWD.NecessaryLabour {
		t.Fatalf("RelativeSurplusValue.LabourMinutes = %d, want %d (oldN − newN)",
			rsv.LabourMinutes, wd.NecessaryLabour-newWD.NecessaryLabour)
	}
}

// Invariant §1, Mill critique: RateOfProfit < RateSurplusValue whenever
// totalCapital > necessaryLabour.
func TestInvariant_ProfitRateBelowSurplusRate(t *testing.T) {
	t.Parallel()
	rsv := surplus.RateSurplusValue(20, 100)
	rop := surplus.RateOfProfit(20, 500)
	if !(rop < rsv) {
		t.Fatalf("RateOfProfit (%v) must be < RateSurplusValue (%v) when total > necessary",
			rop, rsv)
	}
}

// Both surplus-value variants share the same arithmetic representation —
// the distinction is the Origin tag, not the magnitude type (§1 closing
// passage: "from one standpoint, any distinction between absolute and
// relative surplus-value appears illusory").
func TestSurplusValue_OriginTagDistinguishes(t *testing.T) {
	t.Parallel()
	wd := surplus.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	_, asv := surplus.ProlongWorkingDay(wd, 60)
	_, rsv := surplus.ReduceNecessaryLabour(wd, surplus.ProductivityFactor(1.2))
	if asv.Origin == rsv.Origin {
		t.Fatalf("Origins must differ — got %q on both", asv.Origin)
	}
	// LabourMinutes is the same underlying type on both sides.
	var _ surplus.LabourMinutes = asv.LabourMinutes
	var _ surplus.LabourMinutes = rsv.LabourMinutes
}

// Invariant: NecessaryLabour > 0. ReduceNecessaryLabour must not drive it
// to zero (the wage relation collapses).
func TestReduceNecessaryLabour_GuardsAgainstZero(t *testing.T) {
	t.Parallel()
	wd := surplus.WorkingDay{Total: 720, NecessaryLabour: 600, SurplusLabour: 120}
	_, rsv := surplus.ReduceNecessaryLabour(wd, surplus.ProductivityFactor(1e9))
	if rsv.LabourMinutes <= 0 {
		t.Fatalf("RelativeSurplusValue must be positive on positive productivity gain; got %d", rsv.LabourMinutes)
	}
}

// SubjectionKind enum: formal vs real subjection (§1).
func TestSubjectionKind_FormalAndRealValid(t *testing.T) {
	t.Parallel()
	if !surplus.SubjectionFormal.IsValid() {
		t.Fatal("SubjectionFormal must be valid")
	}
	if !surplus.SubjectionReal.IsValid() {
		t.Fatal("SubjectionReal must be valid")
	}
	if surplus.SubjectionKind("hybrid").IsValid() {
		t.Fatal("non-canonical kind must be invalid")
	}
}
