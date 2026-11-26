package market_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

func TestPenceCanonical(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   market.Pence
		want money.Pence
	}{
		{"12000 legacy pence", 12000, money.FromLSD(120, 0, 0)},
		{"48000 legacy pence", 48000, money.FromLSD(480, 0, 0)},
		{"200 legacy pence (Ch. 3 circuit seed: £2 = 200 pence)", 200, money.FromLSD(2, 0, 0)},
		{"zero", 0, 0},
		{"1 rounds down (2.4 -> 2)", 1, 2},
		{"3 rounds down (7.2 -> 7)", 3, 7},
		{"2 rounds up (4.8 -> 5)", 2, 5},
		{"negative rounds away from zero (-4.8 -> -5)", -2, -5},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.in.Canonical(); got != tc.want {
				t.Fatalf("Pence(%d).Canonical() = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

// 240/100 = 12/5, so the fractional part is always a multiple of 0.2 and an
// exact half never occurs; assert the sign symmetry Rescale guarantees instead.
func TestPenceCanonicalSignSymmetry(t *testing.T) {
	t.Parallel()
	for _, n := range []market.Pence{1, 2, 3, 7, 11, 13, 99, 12345} {
		if pos, neg := n.Canonical(), (-n).Canonical(); pos != -neg {
			t.Fatalf("asymmetric rounding for %d: %d vs %d", n, pos, neg)
		}
	}
}
