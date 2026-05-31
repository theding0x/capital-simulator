package credit

import (
	"testing"
)

// ── Fixtures from Capital Vol. III Ch. 26 ────────────────────────────────────
//
// Marx observes the accumulation of money-capital across the industrial cycle:
// post-1847 inactivity (stagnation, abundant loanable funds, low rate), the
// revival phase (demand rising, rate moderate), and the 1857 crisis peak
// (acute demand for means of payment, rate spiking).

// fixtureInactivity returns the post-1847 inactivity observation.
func fixtureInactivity(t *testing.T) MoneyCapitalAccumulation {
	t.Helper()
	obs := InterestCycleObservation{
		Phase:          PhaseInactivity,
		InterestRateBP: 100,
		Period:         "post-1847 inactivity",
	}
	supply := LoanableCapitalSupply{Amount: 5_000_000, Abundant: true}
	idle := IdleIndustrialCapital{Amount: 2_000_000, Description: "idle cotton-mill capital"}
	gap := AccumulationGap{GapBP: 500, Description: "money-capital accumulates faster than productive"}
	mca, err := NewMoneyCapitalAccumulation("Post-1847 Inactivity", obs, supply, idle, gap, "post-crisis stagnation; loanable money piles up")
	if err != nil {
		t.Fatalf("fixtureInactivity: %v", err)
	}
	return mca
}

// fixtureRevival returns the revival-phase observation.
func fixtureRevival(t *testing.T) MoneyCapitalAccumulation {
	t.Helper()
	obs := InterestCycleObservation{
		Phase:          PhaseRevival,
		InterestRateBP: 300,
		Period:         "revival 1850s",
	}
	supply := LoanableCapitalSupply{Amount: 4_000_000, Abundant: true}
	idle := IdleIndustrialCapital{Amount: 500_000, Description: "partially re-engaged capital"}
	gap := AccumulationGap{GapBP: 100, Description: "gap narrows as productive capital re-activates"}
	mca, err := NewMoneyCapitalAccumulation("Revival Phase", obs, supply, idle, gap, "rising demand, moderate rate")
	if err != nil {
		t.Fatalf("fixtureRevival: %v", err)
	}
	return mca
}

// fixtureCrisis returns the 1857 crisis peak observation.
func fixtureCrisis(t *testing.T) MoneyCapitalAccumulation {
	t.Helper()
	obs := InterestCycleObservation{
		Phase:          PhaseCrisis,
		InterestRateBP: 1000,
		Period:         "1857 crisis",
	}
	supply := LoanableCapitalSupply{Amount: 500_000, Abundant: false}
	idle := IdleIndustrialCapital{Amount: 0, Description: "no idle capital — all scrambling for means of payment"}
	gap := AccumulationGap{GapBP: -800, Description: "productive capital outpaces loanable supply"}
	mca, err := NewMoneyCapitalAccumulation("Crisis Peak 1857", obs, supply, idle, gap, "acute demand; rate spikes to 10%")
	if err != nil {
		t.Fatalf("fixtureCrisis: %v", err)
	}
	return mca
}

// ── NewMoneyCapitalAccumulation happy paths ───────────────────────────────────

func TestNewMoneyCapitalAccumulation_Inactivity(t *testing.T) {
	t.Parallel()
	mca := fixtureInactivity(t)
	if mca.CycleObservation.Phase != PhaseInactivity {
		t.Errorf("Phase: got %d, want %d", mca.CycleObservation.Phase, PhaseInactivity)
	}
	if mca.CycleObservation.InterestRateBP != 100 {
		t.Errorf("InterestRateBP: got %d, want 100", mca.CycleObservation.InterestRateBP)
	}
	if !mca.LoanableSupply.Abundant {
		t.Error("LoanableSupply.Abundant must be true in inactivity")
	}
	if mca.IdleCapital.Amount != 2_000_000 {
		t.Errorf("IdleCapital.Amount: got %d, want 2000000", mca.IdleCapital.Amount)
	}
	if mca.Gap.GapBP != 500 {
		t.Errorf("Gap.GapBP: got %d, want 500", mca.Gap.GapBP)
	}
	// ID and CreatedAt must be zero — store assigns them.
	if !mca.ID.IsZero() {
		t.Error("ID must be zero before store persists")
	}
}

func TestNewMoneyCapitalAccumulation_Revival(t *testing.T) {
	t.Parallel()
	mca := fixtureRevival(t)
	if mca.CycleObservation.Phase != PhaseRevival {
		t.Errorf("Phase: got %d, want %d", mca.CycleObservation.Phase, PhaseRevival)
	}
	if mca.CycleObservation.InterestRateBP != 300 {
		t.Errorf("InterestRateBP: got %d, want 300", mca.CycleObservation.InterestRateBP)
	}
	if !mca.LoanableSupply.Abundant {
		t.Error("LoanableSupply.Abundant must be true in revival")
	}
}

func TestNewMoneyCapitalAccumulation_Crisis(t *testing.T) {
	t.Parallel()
	mca := fixtureCrisis(t)
	if mca.CycleObservation.Phase != PhaseCrisis {
		t.Errorf("Phase: got %d, want %d", mca.CycleObservation.Phase, PhaseCrisis)
	}
	if mca.CycleObservation.InterestRateBP != 1000 {
		t.Errorf("InterestRateBP: got %d, want 1000", mca.CycleObservation.InterestRateBP)
	}
	if mca.LoanableSupply.Abundant {
		t.Error("LoanableSupply.Abundant must be false in crisis")
	}
	if mca.IdleCapital.Amount != 0 {
		t.Errorf("IdleCapital.Amount: got %d, want 0", mca.IdleCapital.Amount)
	}
	if mca.Gap.GapBP != -800 {
		t.Errorf("Gap.GapBP: got %d, want -800", mca.Gap.GapBP)
	}
}

// ── Invariant failures ────────────────────────────────────────────────────────

func TestNewMoneyCapitalAccumulation_InvalidPhase(t *testing.T) {
	t.Parallel()
	obs := InterestCycleObservation{Phase: 99, InterestRateBP: 100}
	supply := LoanableCapitalSupply{Amount: 100}
	idle := IdleIndustrialCapital{Amount: 0}
	gap := AccumulationGap{GapBP: 0}
	_, err := NewMoneyCapitalAccumulation("bad", obs, supply, idle, gap, "")
	if err != ErrInvalidCyclePhase {
		t.Errorf("got %v, want ErrInvalidCyclePhase", err)
	}
}

func TestNewMoneyCapitalAccumulation_NegativeInterestRate(t *testing.T) {
	t.Parallel()
	obs := InterestCycleObservation{Phase: PhaseInactivity, InterestRateBP: -1}
	supply := LoanableCapitalSupply{Amount: 100}
	idle := IdleIndustrialCapital{Amount: 0}
	gap := AccumulationGap{GapBP: 0}
	_, err := NewMoneyCapitalAccumulation("bad", obs, supply, idle, gap, "")
	if err != ErrNegativeInterestRate {
		t.Errorf("got %v, want ErrNegativeInterestRate", err)
	}
}

func TestNewMoneyCapitalAccumulation_NegativeIdle(t *testing.T) {
	t.Parallel()
	obs := InterestCycleObservation{Phase: PhaseInactivity, InterestRateBP: 100}
	supply := LoanableCapitalSupply{Amount: 100}
	idle := IdleIndustrialCapital{Amount: -1}
	gap := AccumulationGap{GapBP: 0}
	_, err := NewMoneyCapitalAccumulation("bad", obs, supply, idle, gap, "")
	if err != ErrNegativeIdleAmount {
		t.Errorf("got %v, want ErrNegativeIdleAmount", err)
	}
}

func TestNewMoneyCapitalAccumulation_NegativeLoanableAmount(t *testing.T) {
	t.Parallel()
	obs := InterestCycleObservation{Phase: PhaseInactivity, InterestRateBP: 100}
	supply := LoanableCapitalSupply{Amount: -1}
	idle := IdleIndustrialCapital{Amount: 0}
	gap := AccumulationGap{GapBP: 0}
	_, err := NewMoneyCapitalAccumulation("bad", obs, supply, idle, gap, "")
	if err != ErrNegativeIdleAmount {
		t.Errorf("got %v, want ErrNegativeIdleAmount", err)
	}
}

// ── NewMoneyCapitalAccumulationID ─────────────────────────────────────────────

func TestNewMoneyCapitalAccumulationID_Format(t *testing.T) {
	t.Parallel()
	id := NewMoneyCapitalAccumulationID()
	if len(string(id)) != 24 {
		t.Errorf("ID length: got %d, want 24 hex chars", len(string(id)))
	}
}

func TestNewMoneyCapitalAccumulationID_Unique(t *testing.T) {
	t.Parallel()
	a := NewMoneyCapitalAccumulationID()
	b := NewMoneyCapitalAccumulationID()
	if a == b {
		t.Error("two IDs must not be equal")
	}
}

func TestMoneyCapitalAccumulationID_IsZero(t *testing.T) {
	t.Parallel()
	var id MoneyCapitalAccumulationID
	if !id.IsZero() {
		t.Error("zero value must be zero")
	}
	id = NewMoneyCapitalAccumulationID()
	if id.IsZero() {
		t.Error("fresh ID must not be zero")
	}
}

// ── ClassifyAccumulation ──────────────────────────────────────────────────────

func TestClassifyAccumulation_AllPhases(t *testing.T) {
	t.Parallel()
	cases := []struct {
		phase IndustrialCyclePhase
		want  string
	}{
		{PhaseInactivity, "post-crisis idle: loanable money abundant, low demand, depressed rate"},
		{PhaseRevival, "revival: loanable money still abundant, rising demand, moderate rate"},
		{PhaseProsperity, "prosperity: credit stretched, rate rising above average"},
		{PhaseOverproduction, "overproduction: credit fully extended, rate high"},
		{PhaseCrisis, "crisis: acute demand for means of payment, rate peaks"},
	}
	for _, tc := range cases {
		obs := InterestCycleObservation{Phase: tc.phase}
		supply := LoanableCapitalSupply{}
		got := ClassifyAccumulation(obs, supply)
		if got != tc.want {
			t.Errorf("phase %d: got %q, want %q", tc.phase, got, tc.want)
		}
	}
}

func TestClassifyAccumulation_InvalidPhase(t *testing.T) {
	t.Parallel()
	obs := InterestCycleObservation{Phase: 99}
	got := ClassifyAccumulation(obs, LoanableCapitalSupply{})
	if got != "unknown phase" {
		t.Errorf("got %q, want %q", got, "unknown phase")
	}
}
