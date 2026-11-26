package circulation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

// TestPence_Canonical bridges the legacy 100-per-£ scale to money.Pence
// (240 per £). Fixtures are Vol. II Ch. 2 and 3, Spinning Mill 1871:
// c=£372, v=£50, s=£78, advance £422, C' £500.
func TestPence_Canonical(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		legacy circulation.Pence
		want   money.Pence
	}{
		{"zero", 0, 0},
		{"constant c £372", 37_200, money.FromLSD(372, 0, 0)},
		{"variable v £50", 5_000, money.FromLSD(50, 0, 0)},
		{"surplus s £78", 7_800, money.FromLSD(78, 0, 0)},
		{"advance £422", 42_200, money.FromLSD(422, 0, 0)},
		{"C' £500", 50_000, money.FromLSD(500, 0, 0)},
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
