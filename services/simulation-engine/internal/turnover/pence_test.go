package turnover_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/turnover"
)

// TestPenceCanonical converts the legacy 100-per-pound turnover scale to the
// canonical 240-per-pound money.Pence.
func TestPenceCanonical(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		legacy turnover.Pence
		want   money.Pence
	}{
		{"Marx Ch. 9 cotton mill fixed capital £80,000", 8_000_000, money.FromLSD(80000, 0, 0)},
		{"Marx Ch. 9 cotton mill total advance £100,000", 10_000_000, money.FromLSD(100000, 0, 0)},
		{"zero", 0, 0},
		{"one legacy penny 2.4 rounds down to 2", 1, 2},
		{"five legacy pence is exactly 12", 5, 12},
		{"odd value 3 -> 7.2 rounds down to 7", 3, 7},
		{"odd value 7 -> 16.8 rounds up to 17", 7, 17},
		{"negative mirrors positive", -7, -17},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.legacy.Canonical(); got != tc.want {
				t.Errorf("Pence(%d).Canonical() = %d, want %d", tc.legacy, got, tc.want)
			}
		})
	}

	if got := turnover.Pence(8_000_000).Canonical(); got != 19_200_000 {
		t.Errorf("£80,000 canonical = %d, want 19200000", got)
	}
}
