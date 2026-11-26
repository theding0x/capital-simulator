package agent_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

func TestPenceCanonical(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   agent.Pence
		want money.Pence
	}{
		{"Watt's £521 (seed 00003)", 52100, money.FromLSD(521, 0, 0)},
		{"£100 advance", 10000, money.FromLSD(100, 0, 0)},
		{"zero", 0, 0},
		{"odd value rounds down (1 -> 2.4 -> 2)", 1, 2},
		{"odd value rounds up (3 -> 7.2 -> 7)", 3, 7},
		{"odd value rounds up (2 -> 4.8 -> 5)", 2, 5},
		{"exact multiple (5 -> 12)", 5, 12},
		{"negative mirrors positive (-1 -> -2)", -1, -2},
		{"negative rounds away from zero (-2 -> -5)", -2, -5},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.in.Canonical(); got != tc.want {
				t.Errorf("Pence(%d).Canonical() = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

func TestPenceCanonicalWattIs125040(t *testing.T) {
	t.Parallel()
	if got := agent.Pence(52100).Canonical(); got != 125040 {
		t.Errorf("Watt's £521 = %d, want 125040", got)
	}
}
