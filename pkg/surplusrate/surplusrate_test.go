package surplusrate

import (
	"errors"
	"testing"
)

func TestRatioMarxFixtures(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		s, v int64
		want float64
	}{
		{"100 percent: 6h surplus to 6h necessary", 360, 360, 1.0},
		{"English agricultural labourer 300 percent", 1080, 360, 3.0},
		{"no surplus labour", 0, 360, 0},
	}
	for _, tc := range cases {
		got, err := Ratio(tc.s, tc.v)
		if err != nil {
			t.Fatalf("%s: unexpected error %v", tc.name, err)
		}
		if got != tc.want {
			t.Errorf("%s: Ratio(%d,%d) = %v, want %v", tc.name, tc.s, tc.v, got, tc.want)
		}
	}
}

func TestBasisPointsRoundsHalfUp(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		s, v int64
		want int64
	}{
		{"100 percent", 360, 360, 10000},
		{"300 percent", 1080, 360, 30000},
		{"exact half rounds up: 1/20000 of a unit", 1, 20000, 1},
		{"just under half rounds down", 1, 20001, 0},
		{"two thirds: 6666.67 rounds to 6667", 2, 3, 6667},
		{"one third: 3333.33 rounds to 3333", 1, 3, 3333},
	}
	for _, tc := range cases {
		got, err := BasisPoints(tc.s, tc.v)
		if err != nil {
			t.Fatalf("%s: unexpected error %v", tc.name, err)
		}
		if got != tc.want {
			t.Errorf("%s: BasisPoints(%d,%d) = %d, want %d", tc.name, tc.s, tc.v, got, tc.want)
		}
	}
}

func TestZeroOrNegativeVariableIsRejected(t *testing.T) {
	t.Parallel()

	for _, v := range []int64{0, -1} {
		if _, err := Ratio(100, v); !errors.Is(err, ErrZeroVariable) {
			t.Errorf("Ratio(100,%d) err = %v, want ErrZeroVariable", v, err)
		}
		if _, err := BasisPoints(100, v); !errors.Is(err, ErrZeroVariable) {
			t.Errorf("BasisPoints(100,%d) err = %v, want ErrZeroVariable", v, err)
		}
	}
}
