package agent_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// § 12-hour day, 24 pieces: daily wage 3s. = 144 farthings → piece price = 144/24 = 6 farthings (1½d.)
func TestComputePiecePrice_24Pieces(t *testing.T) {
	t.Parallel()

	got := agent.ComputePiecePrice(144, 24)
	if got != 6 {
		t.Errorf("ComputePiecePrice(144, 24) = %d, want 6", got)
	}
	// invariant: price × normalOutput == dailyWage
	if got*24 != 144 {
		t.Errorf("invariant broken: %d×24 = %d, want 144", got, got*24)
	}
}

// § doubled productiveness, 48 pieces: same daily wage 144 farthings → piece price = 144/48 = 3 farthings (¾d.)
func TestComputePiecePrice_48Pieces(t *testing.T) {
	t.Parallel()

	got := agent.ComputePiecePrice(144, 48)
	if got != 3 {
		t.Errorf("ComputePiecePrice(144, 48) = %d, want 3", got)
	}
	if got*48 != 144 {
		t.Errorf("invariant broken: %d×48 = %d, want 144", got, got*48)
	}
}

// § piece value > piece price: day value product = 6s. = 288 farthings, 24 pieces → piece value = 12 farthings (3d.)
func TestComputePieceValue_ExceedsPiecePrice(t *testing.T) {
	t.Parallel()

	piecePrice := agent.ComputePiecePrice(144, 24) // 6 farthings — wage side
	pieceValue := agent.ComputePieceValue(288, 24) // 12 farthings — value side
	if pieceValue <= piecePrice {
		t.Errorf("piece value %d must exceed piece price %d", pieceValue, piecePrice)
	}
	if pieceValue != 12 {
		t.Errorf("ComputePieceValue(288, 24) = %d, want 12", pieceValue)
	}
}

// § accepted session: 24 pieces × 6 farthings piece price = 144 farthings earnings
func TestComputePieceEarnings_Accepted(t *testing.T) {
	t.Parallel()

	pw := agent.PieceWage{PricePence: 6, NormalOutput: 24}
	s := agent.PieceSession{PiecesProduced: 24, QualityOutcome: agent.QualityAccepted}
	got := agent.ComputePieceEarnings(s, pw)
	if got != 144 {
		t.Errorf("ComputePieceEarnings accepted = %d, want 144", got)
	}
}

// § rejected session earns nothing — quality control enforced by the wage form itself
func TestComputePieceEarnings_Rejected(t *testing.T) {
	t.Parallel()

	pw := agent.PieceWage{PricePence: 6, NormalOutput: 24}
	s := agent.PieceSession{PiecesProduced: 24, QualityOutcome: agent.QualityRejected}
	got := agent.ComputePieceEarnings(s, pw)
	if got != 0 {
		t.Errorf("ComputePieceEarnings rejected = %d, want 0", got)
	}
}

// § sub-contract sweating: head keeps spread = piece_rate - assistant_rate
func TestSubContractSpread(t *testing.T) {
	t.Parallel()

	sc := agent.SubContract{
		HeadLabourerID:     "head-x",
		AssistantIDs:       []agent.AgentID{"assist-a", "assist-b"},
		PieceRatePence:     6,
		AssistantRatePence: 4,
	}
	got := agent.SubContractSpread(sc)
	if got != 2 {
		t.Errorf("SubContractSpread = %d, want 2", got)
	}
}

// Validate: normal_output zero is rejected
func TestPieceWage_Validate_ZeroOutput(t *testing.T) {
	t.Parallel()

	pw := agent.PieceWage{AgentID: "agent-1", PricePence: 6, NormalOutput: 0}
	if err := pw.Validate(); err == nil {
		t.Error("expected error for zero normal_output, got nil")
	}
}

// Validate: SubContract spread must be positive
func TestSubContract_Validate_InvertedSpread(t *testing.T) {
	t.Parallel()

	sc := agent.SubContract{
		HeadLabourerID:     "head-x",
		PieceRatePence:     3,
		AssistantRatePence: 6,
	}
	if err := sc.Validate(); err == nil {
		t.Error("expected error for inverted spread, got nil")
	}
}
