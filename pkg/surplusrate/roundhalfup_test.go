package surplusrate

import "testing"

func TestRoundHalfUpDiv(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		num, den int64
		want     int64
	}{
		{"exact", 300, 100, 3},
		{"just under half rounds down", 1, 3, 0},
		{"exactly half rounds up", 3, 2, 2},
		{"over half rounds up", 2, 3, 1},
		{"zero numerator", 0, 7, 0},
		{"zero denominator", 5, 0, 0},
		{"negative denominator", 5, -2, 0},
		{"does not overflow near int64 max", 1<<62 + 1, 1 << 62, 1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := RoundHalfUpDiv(tc.num, tc.den); got != tc.want {
				t.Errorf("RoundHalfUpDiv(%d,%d) = %d, want %d", tc.num, tc.den, got, tc.want)
			}
		})
	}
}

// BasisPoints is round-half-up division of s*10000 by v, so the two agree.
func TestRoundHalfUpDivAgreesWithBasisPoints(t *testing.T) {
	t.Parallel()

	for _, c := range []struct{ s, v int64 }{{1, 3}, {2, 3}, {3, 8}, {7, 9}, {300, 100}} {
		bp, _ := BasisPoints(c.s, c.v)
		if got := RoundHalfUpDiv(c.s*BasisPointsPerUnit, c.v); got != bp {
			t.Errorf("s=%d v=%d: RoundHalfUpDiv = %d, BasisPoints = %d", c.s, c.v, got, bp)
		}
	}
}
