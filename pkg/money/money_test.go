package money_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
)

// Pre-decimal arithmetic: £1 = 20s = 240d, so Marx's shilling and pence
// fixtures are exact integers.
func TestFromLSD_PreDecimal(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		l, s, d int64
		want    money.Pence
	}{
		{"one pound", 1, 0, 0, 240},
		{"3s (Ch. 17 value of labour-power)", 0, 3, 0, 36},
		{"3s 6d (Ch. 17 price of labour-power)", 0, 3, 6, 42},
		{"£410 (Ch. 9 constant capital)", 410, 0, 0, 98400},
	}
	for _, c := range cases {
		if got := money.FromLSD(c.l, c.s, c.d); got != c.want {
			t.Errorf("%s: got %d, want %d", c.name, got, c.want)
		}
	}
}

// The MELT is Marx's Ch. 17 ratio: a 12-hour day (720 min) creates 6s (72d).
func TestMELT_TwelveHourDayIsSixShillings(t *testing.T) {
	t.Parallel()
	if got := money.FromLabourMinutes(720); got != money.FromLSD(0, 6, 0) {
		t.Fatalf("720 labour-minutes: got %d d, want 72 d (6s)", got)
	}
	if got := money.FromLSD(0, 6, 0).ToLabourMinutes(); got != 720 {
		t.Fatalf("6s: got %d labour-minutes, want 720", got)
	}
	if money.LabourMinutesPerPound != 40*60 {
		t.Fatalf("£1 must be 40 hours, got %d min", money.LabourMinutesPerPound)
	}
}

func TestLabourMinutes_RoundTripIsExact(t *testing.T) {
	t.Parallel()
	for _, p := range []money.Pence{0, 1, 36, 42, 240, 98400} {
		if got := money.FromLabourMinutes(p.ToLabourMinutes()); got != p {
			t.Errorf("round trip of %d d: got %d", p, got)
		}
	}
}

func TestFromLabourMinutes_RoundsHalfUp(t *testing.T) {
	t.Parallel()
	// 1d = 10 min: 4 min = 0.4d -> 0; 5 min = 0.5d -> 1; 14 -> 1; 15 -> 2.
	for min, want := range map[int64]money.Pence{4: 0, 5: 1, 14: 1, 15: 2} {
		if got := money.FromLabourMinutes(min); got != want {
			t.Errorf("%d min: got %d, want %d", min, got, want)
		}
	}
}

func TestRescale_LegacyScales(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name         string
		amount, from int64
		want         money.Pence
	}{
		{"£1 at 100/£", 100, 100, 240},
		{"£1000 at the reproduction scale (1000/£)", 1_000_000, 1000, 240_000},
		{"already canonical", 42, 240, 42},
		{"1 at 100/£ = 2.4d rounds down", 1, 100, 2},
		{"3 at 100/£ = 7.2d rounds down", 3, 100, 7},
		{"7 at 100/£ = 16.8d rounds up", 7, 100, 17},
		{"tie 0.5d rounds up", 5, 480, 3}, // 5*240/480 = 2.5
		{"negative ties round away from zero", -5, 480, -3},
	}
	for _, c := range cases {
		if got := money.Rescale(c.amount, c.from); got != c.want {
			t.Errorf("%s: got %d, want %d", c.name, got, c.want)
		}
	}
}
