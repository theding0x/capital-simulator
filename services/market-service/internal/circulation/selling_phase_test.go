package circulation_test

import (
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/circulation"
)

func TestSellingPhase_IsOpen_NilClosedAt(t *testing.T) {
	t.Parallel()
	sp := circulation.SellingPhase{
		ID:      circulation.NewSellingPhaseID(),
		Outcome: circulation.OutcomePending,
	}
	if !sp.IsOpen() {
		t.Error("SellingPhase with nil ClosedAt must be open")
	}
}

func TestSellingPhase_IsOpen_AfterClose(t *testing.T) {
	t.Parallel()
	now := time.Now()
	sp := circulation.SellingPhase{
		ID:       circulation.NewSellingPhaseID(),
		ClosedAt: &now,
		Outcome:  circulation.OutcomeSold,
	}
	if sp.IsOpen() {
		t.Error("SellingPhase with non-nil ClosedAt must be closed")
	}
}

func TestSellingPhase_ElapsedNanos_WhileOpen(t *testing.T) {
	t.Parallel()
	// ElapsedNanos returns 0 while the phase is still open.
	sp := circulation.SellingPhase{
		ID:       circulation.NewSellingPhaseID(),
		OpenedAt: time.Now().Add(-time.Hour),
	}
	if got := sp.ElapsedNanos(); got != 0 {
		t.Errorf("ElapsedNanos() = %d while open, want 0", got)
	}
}

func TestSellingPhase_ElapsedNanos_AfterClose(t *testing.T) {
	t.Parallel()
	// SpinningMill1871: 5-day selling phase (C—M).
	// §: "C—M, the sale, is the most difficult part of its metamorphosis."
	opened := time.Date(1871, 1, 1, 0, 0, 0, 0, time.UTC)
	closed := opened.Add(5 * 24 * time.Hour)
	sp := circulation.SellingPhase{
		ID:       circulation.NewSellingPhaseID(),
		OpenedAt: opened,
		ClosedAt: &closed,
		Outcome:  circulation.OutcomeSold,
	}
	want := (5 * 24 * time.Hour).Nanoseconds()
	if got := sp.ElapsedNanos(); got != want {
		t.Errorf("ElapsedNanos() = %d, want %d (5 days)", got, want)
	}
}

func TestSellingOutcome_AllValuesDistinct(t *testing.T) {
	t.Parallel()
	// All four outcomes must be distinguishable strings.
	outcomes := []circulation.SellingOutcome{
		circulation.OutcomeSold,
		circulation.OutcomeSpoiled,
		circulation.OutcomePartial,
		circulation.OutcomePending,
	}
	seen := map[circulation.SellingOutcome]bool{}
	for _, o := range outcomes {
		if seen[o] {
			t.Errorf("duplicate SellingOutcome value: %q", o)
		}
		seen[o] = true
	}
}

func TestBuyingPhase_IsOpen_NilClosedAt(t *testing.T) {
	t.Parallel()
	bp := circulation.BuyingPhase{
		ID: circulation.NewBuyingPhaseID(),
	}
	if !bp.IsOpen() {
		t.Error("BuyingPhase with nil ClosedAt must be open")
	}
}

func TestBuyingPhase_ElapsedNanos_AfterClose(t *testing.T) {
	t.Parallel()
	// SpinningMill1871: 2-day buying phase (M—C).
	// §: "M—C, the purchase, is the easier part."
	opened := time.Date(1871, 1, 6, 0, 0, 0, 0, time.UTC)
	closed := opened.Add(2 * 24 * time.Hour)
	bp := circulation.BuyingPhase{
		ID:       circulation.NewBuyingPhaseID(),
		OpenedAt: opened,
		ClosedAt: &closed,
	}
	want := (2 * 24 * time.Hour).Nanoseconds()
	if got := bp.ElapsedNanos(); got != want {
		t.Errorf("ElapsedNanos() = %d, want %d (2 days)", got, want)
	}
}

func TestNewSellingPhaseID_Unique(t *testing.T) {
	t.Parallel()
	a := circulation.NewSellingPhaseID()
	b := circulation.NewSellingPhaseID()
	if a == b {
		t.Errorf("duplicate SellingPhaseID: %q", a)
	}
	if a.IsZero() {
		t.Error("NewSellingPhaseID() returned zero value")
	}
}

func TestNewBuyingPhaseID_Unique(t *testing.T) {
	t.Parallel()
	a := circulation.NewBuyingPhaseID()
	b := circulation.NewBuyingPhaseID()
	if a == b {
		t.Errorf("duplicate BuyingPhaseID: %q", a)
	}
}
