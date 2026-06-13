package simulation

import "testing"

func TestNewAbodeStateReadout(t *testing.T) {
	t.Parallel()
	r := NewAbodeState().Readout()
	if r.TotalVariablePence != 300000 {
		t.Errorf("v = %d, want 300000", r.TotalVariablePence)
	}
	// 120 employed (v/baseWage = 300000/2500), 30 in the reserve of 150.
	if r.EmployedCount != 120 {
		t.Errorf("employed = %d, want 120", r.EmployedCount)
	}
	if r.ReserveArmyCount != 30 {
		t.Errorf("reserve = %d, want 30", r.ReserveArmyCount)
	}
	// pressure = 30/150 = 2000 bp; s' = 10000*(1+0.20) = 12000 bp.
	if r.RateOfExploitationBP != 12000 {
		t.Errorf("s/v = %d, want 12000", r.RateOfExploitationBP)
	}
	// c/v = 600000/300000 = 2.0 → 20000 bp.
	if r.OrganicCompositionBP != 20000 {
		t.Errorf("c/v = %d, want 20000", r.OrganicCompositionBP)
	}
	// surplus = v * s'/10000 = 300000 * 1.2 = 360000.
	if r.TotalSurplusPence != 360000 {
		t.Errorf("surplus = %d, want 360000", r.TotalSurplusPence)
	}
	// working day 600' splits necessary = 600*10000/22000 = 272, surplus = 328.
	if r.NecessaryLabourMinutes+r.SurplusLabourMinutes != SocialWorkingDayMinutes {
		t.Errorf("working day = %d, want %d",
			r.NecessaryLabourMinutes+r.SurplusLabourMinutes, SocialWorkingDayMinutes)
	}
	if r.SurplusLabourMinutes <= r.NecessaryLabourMinutes {
		t.Errorf("surplus labour %d should exceed necessary %d at s'=120%%",
			r.SurplusLabourMinutes, r.NecessaryLabourMinutes)
	}
	// paid wage compressed below the value of labour-power (2500).
	if r.WagePence >= 2500 {
		t.Errorf("paid wage = %d, want < 2500 (value of labour-power)", r.WagePence)
	}
}

func TestAdvanceGeneralLawImmiserates(t *testing.T) {
	t.Parallel()
	s := NewAbodeState()
	first := s.Readout()
	var last GeneralLawPeriod
	for i := 0; i < 12; i++ {
		var p GeneralLawPeriod
		s, p = AdvanceGeneralLaw(s)
		last = p
	}
	// Over twelve periods the law runs its course: the reserve army grows, the
	// organic composition rises, the rate of exploitation rises, the wage falls.
	if last.ReserveArmyCount <= first.ReserveArmyCount {
		t.Errorf("reserve army %d did not grow from %d", last.ReserveArmyCount, first.ReserveArmyCount)
	}
	if last.OrganicCompositionBP <= first.OrganicCompositionBP {
		t.Errorf("composition %d did not rise from %d", last.OrganicCompositionBP, first.OrganicCompositionBP)
	}
	if last.RateOfExploitationBP <= first.RateOfExploitationBP {
		t.Errorf("s/v %d did not rise from %d", last.RateOfExploitationBP, first.RateOfExploitationBP)
	}
	if last.WagePence >= first.WagePence {
		t.Errorf("wage %d did not fall from %d", last.WagePence, first.WagePence)
	}
	if s.Period != 12 {
		t.Errorf("period = %d, want 12", s.Period)
	}
}

func TestAdvanceGeneralLawConserves(t *testing.T) {
	t.Parallel()
	s := NewAbodeState()
	next, _ := AdvanceGeneralLaw(s)
	// Accumulation + displacement only move value between c and v and add the
	// capitalised surplus; neither v nor c goes negative.
	if next.VariablePence < 0 || next.ConstantPence < 0 {
		t.Fatalf("negative capital: c=%d v=%d", next.ConstantPence, next.VariablePence)
	}
	// Total social capital grows (α·s re-accumulated) — the spiral.
	if next.ConstantPence+next.VariablePence <= s.ConstantPence+s.VariablePence {
		t.Errorf("total capital did not grow: %d -> %d",
			s.ConstantPence+s.VariablePence, next.ConstantPence+next.VariablePence)
	}
}

func TestAdvanceGeneralLawDoesNotCollapse(t *testing.T) {
	t.Parallel()
	s := NewAbodeState()
	seedV := s.VariablePence
	seedReserve := s.Readout().ReserveArmyCount
	// Run the law well past where the old machinery-destroys-v mechanism stranded
	// variable capital at a few pence — the degenerate collapse the user hit:
	// Σv = 0, wage pressure pegged at 100%, the whole workforce idle in the reserve.
	for i := 0; i < 60; i++ {
		s, _ = AdvanceGeneralLaw(s)
	}
	r := s.Readout()
	// Variable capital (living labour) GROWS in absolute magnitude — never destroyed
	// (Vol. I Ch. 25 §2: "the variable part... increases, but in a diminishing proportion").
	if s.VariablePence <= seedV {
		t.Errorf("variable capital did not grow over 60 periods: %d (seed %d) — living labour collapsed", s.VariablePence, seedV)
	}
	if r.EmployedCount <= 0 {
		t.Errorf("employment collapsed to %d — no living labour employed", r.EmployedCount)
	}
	if r.TotalVariablePence <= 0 {
		t.Errorf("Σv collapsed to %d", r.TotalVariablePence)
	}
	// The active army always persists, so pressure never fully pegs at 100%.
	if r.ReserveArmyPressureBP >= 10000 {
		t.Errorf("wage pressure pegged at %d bp (100%%) — the active army vanished", r.ReserveArmyPressureBP)
	}
	// The reserve army still grows (relative surplus population) — immiseration,
	// not degeneration.
	if r.ReserveArmyCount <= seedReserve {
		t.Errorf("reserve army did not grow: %d (seed %d)", r.ReserveArmyCount, seedReserve)
	}
}

func TestWorkerSupplyCappedAtSociety(t *testing.T) {
	t.Parallel()
	// The labouring population cannot exceed society as a whole (London 1850).
	s := NewAbodeState()
	// Start the supply just under the ceiling so a single period of growth would
	// otherwise overshoot it.
	s.WorkerSupply = LondonWorkingPopulation - 10
	next, _ := AdvanceGeneralLaw(s)
	if next.WorkerSupply > LondonWorkingPopulation {
		t.Errorf("worker supply %d exceeded society of %d — more workers at the gate than exist",
			next.WorkerSupply, LondonWorkingPopulation)
	}
	// Already at the ceiling: the supply holds, it does not keep climbing.
	s.WorkerSupply = LondonWorkingPopulation
	next, _ = AdvanceGeneralLaw(s)
	if next.WorkerSupply != LondonWorkingPopulation {
		t.Errorf("worker supply moved off the ceiling: %d, want %d", next.WorkerSupply, LondonWorkingPopulation)
	}
	// The reserve army at the gate is then bounded by the whole society too.
	if r := next.Readout(); r.ReserveArmyCount > r.TotalPopulation {
		t.Errorf("reserve army %d exceeds society %d", r.ReserveArmyCount, r.TotalPopulation)
	}
}

func TestEmploymentNeverExceedsSociety(t *testing.T) {
	t.Parallel()
	// Capital's demand for labour can run away as it accumulates, but no more
	// workers can be employed than exist in society, and employed + at-the-gate can
	// never exceed the whole population. Drive the law hard past the point where
	// demand outstrips the capped supply.
	s := NewAbodeState()
	for i := 0; i < 2000; i++ {
		s, _ = AdvanceGeneralLaw(s)
	}
	r := s.Readout()
	if r.EmployedCount > r.TotalPopulation {
		t.Errorf("employed %d exceeds society of %d — more workers employed than exist",
			r.EmployedCount, r.TotalPopulation)
	}
	if r.EmployedCount+r.ReserveArmyCount > r.TotalPopulation {
		t.Errorf("employed %d + at the gate %d exceeds society of %d",
			r.EmployedCount, r.ReserveArmyCount, r.TotalPopulation)
	}
}

func TestApplyLevers(t *testing.T) {
	t.Parallel()
	s := NewAbodeState() // surplus_rate_base=10000, base_wage=2500, accum=5000

	sr, wage, ac := int64(20000), int64(4000), int64(8000)
	got := s.ApplyLevers(LeverUpdate{SurplusRateBaseBP: &sr, BaseWagePence: &wage, AccumulationRateBP: &ac})
	if got.SurplusRateBaseBP != 20000 || got.BaseWagePence != 4000 || got.AccumulationRateBP != 8000 {
		t.Errorf("levers not applied: %+v", got)
	}

	// A partial update leaves the other parameters untouched.
	only := s.ApplyLevers(LeverUpdate{AccumulationRateBP: &ac})
	if only.AccumulationRateBP != 8000 || only.BaseWagePence != s.BaseWagePence || only.SurplusRateBaseBP != s.SurplusRateBaseBP {
		t.Errorf("partial update bled into other fields: %+v", only)
	}

	// Clamps: α to [0,10000]; wage floored at 1; base surplus rate to [0,100000].
	big, zero, huge := int64(99999), int64(0), int64(500000)
	cl := s.ApplyLevers(LeverUpdate{AccumulationRateBP: &big, BaseWagePence: &zero, SurplusRateBaseBP: &huge})
	if cl.AccumulationRateBP != 10000 {
		t.Errorf("alpha not clamped: %d", cl.AccumulationRateBP)
	}
	if cl.BaseWagePence != 1 {
		t.Errorf("wage not floored: %d", cl.BaseWagePence)
	}
	if cl.SurplusRateBaseBP != 100000 {
		t.Errorf("surplus rate not clamped: %d", cl.SurplusRateBaseBP)
	}

	// An empty update changes nothing and reports empty.
	if !(LeverUpdate{}).IsEmpty() {
		t.Error("empty LeverUpdate should report IsEmpty")
	}
	if got := s.ApplyLevers(LeverUpdate{}); got != s {
		t.Error("empty update should be a no-op")
	}
}
