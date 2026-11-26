package composition_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/composition"
)

// TestPence_Canonical bridges the legacy 100-per-£ scale to money.Pence
// (240 per £). Fixtures are Vol. II Ch. 8: the £10,000 ten-year machine and
// its £1,000 annual wear.
func TestPence_Canonical(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		legacy composition.Pence
		want   money.Pence
	}{
		{"zero", 0, 0},
		{"machine £10,000", 1_000_000, money.FromLSD(10000, 0, 0)},
		{"annual wear £1,000", 100_000, money.FromLSD(1000, 0, 0)},
		{"railway stock £4,000", 400_000, money.FromLSD(4000, 0, 0)},
		{"one legacy penny rounds down (2.4 -> 2)", 1, 2},
		{"odd value rounds down (3 -> 7.2 -> 7)", 3, 7},
		{"odd value rounds up (7 -> 16.8 -> 17)", 7, 17},
		{"negative rounds away from zero (-7 -> -17)", -7, -17},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.legacy.Canonical(); got != tt.want {
				t.Errorf("Pence(%d).Canonical() = %d, want %d", tt.legacy, got, tt.want)
			}
		})
	}
}
