package commodity_test

import (
	"errors"
	"math"
	"testing"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)

// Ch.9 §3, Senior's "last hour": an 11½-hour day in which, on the
// manufacturers' own figures, wages and net profit are equal, so the workman
// "produces his wages in 5¾ hours, and your net profit in the other 5¾ hours".
const (
	seniorNecessary commodity.NecessaryLabour = 345 // 5¾ h
	seniorDay       commodity.LabourMinutes   = 690 // 11½ h
)

func TestReviseWorkingDay_Ch9SeniorFixtures(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name        string
		day         commodity.LabourMinutes
		wantSurplus commodity.SurplusLabour
		wantRate    float64
		wantBP      int64
	}{
		// "the rate of surplus-value ... 100%"
		{"11½ h, Senior's day", seniorDay, 345, 1.0, 10000},
		// "the latter will be increased from 5¾ hours' labour to 7¼ hours'
		// labour, and the rate of surplus-value from 100% to 126 2/23%"
		{"13 h, lengthened", 780, 435, 29.0 / 23.0, 12609},
		// "the surplus-labour will fall from 5¾ hours to 4¾ hours ... namely
		// 82 14/23%": a one-hour cut, to 10½ h (Marx's text says 10, but only
		// 10½ h gives 4¾ h).
		{"10½ h, cut by one hour", 630, 285, 19.0 / 23.0, 8261},
		// What a cut to the 10 h the text names really gives.
		{"10 h, cut by 1½ hours", 600, 255, 17.0 / 23.0, 7391},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := commodity.ReviseWorkingDay(seniorNecessary, tc.day)
			if err != nil {
				t.Fatalf("ReviseWorkingDay: %v", err)
			}
			if got.WorkingDay != tc.day {
				t.Errorf("WorkingDay = %d, want %d", got.WorkingDay, tc.day)
			}
			if got.NecessaryLabour != seniorNecessary {
				t.Errorf("NecessaryLabour = %d, want %d (held fixed)", got.NecessaryLabour, seniorNecessary)
			}
			if got.SurplusLabour != tc.wantSurplus {
				t.Errorf("SurplusLabour = %d, want %d", got.SurplusLabour, tc.wantSurplus)
			}
			if math.Abs(float64(got.Rate)-tc.wantRate) > 1e-9 {
				t.Errorf("Rate = %v, want %v", got.Rate, tc.wantRate)
			}
			if got.RateBasisPoints != tc.wantBP {
				t.Errorf("RateBasisPoints = %d, want %d", got.RateBasisPoints, tc.wantBP)
			}
		})
	}
}

// Marx's rebuttal: lengthening the day by 1½ hours lifts the rate by about a
// quarter, nowhere near Senior's "more than doubled", and shortening it by an
// hour leaves a "very profitable" rate, so no "last hour" holds the profit.
func TestReviseWorkingDay_SeniorsLastHourRefuted(t *testing.T) {
	t.Parallel()
	base, err := commodity.ReviseWorkingDay(seniorNecessary, seniorDay)
	if err != nil {
		t.Fatalf("base: %v", err)
	}
	longer, err := commodity.ReviseWorkingDay(seniorNecessary, 780)
	if err != nil {
		t.Fatalf("longer: %v", err)
	}
	shorter, err := commodity.ReviseWorkingDay(seniorNecessary, 630)
	if err != nil {
		t.Fatalf("shorter: %v", err)
	}
	if float64(longer.Rate) >= 2*float64(base.Rate) {
		t.Errorf("rate on the 13 h day = %v, which doubles the %v Senior expected", longer.Rate, base.Rate)
	}
	if shorter.SurplusLabour <= 0 {
		t.Errorf("surplus labour on the 10½ h day = %d; Senior said the net profit would be destroyed", shorter.SurplusLabour)
	}
}

// The necessary labour is fixed, so every minute added to or removed from the
// day is added to or removed from the surplus labour alone.
func TestReviseWorkingDay_PartitionInvariant(t *testing.T) {
	t.Parallel()
	for day := commodity.LabourMinutes(seniorNecessary); day <= 24*60; day += 15 {
		got, err := commodity.ReviseWorkingDay(seniorNecessary, day)
		if err != nil {
			t.Fatalf("day %d: %v", day, err)
		}
		if commodity.LabourMinutes(got.NecessaryLabour)+commodity.LabourMinutes(got.SurplusLabour) != got.WorkingDay {
			t.Fatalf("day %d: necessary %d + surplus %d != working day %d",
				day, got.NecessaryLabour, got.SurplusLabour, got.WorkingDay)
		}
		if got.SurplusLabour < 0 {
			t.Fatalf("day %d: negative surplus labour %d", day, got.SurplusLabour)
		}
		if want := float64(got.SurplusLabour) / float64(got.NecessaryLabour); math.Abs(float64(got.Rate)-want) > 1e-12 {
			t.Fatalf("day %d: rate %v != surplus/necessary %v", day, got.Rate, want)
		}
	}
}

// A day exactly as long as the necessary labour leaves no surplus: a zero
// rate, not an error.
func TestReviseWorkingDay_DayEqualToNecessaryLabour(t *testing.T) {
	t.Parallel()
	got, err := commodity.ReviseWorkingDay(seniorNecessary, commodity.LabourMinutes(seniorNecessary))
	if err != nil {
		t.Fatalf("ReviseWorkingDay: %v", err)
	}
	if got.SurplusLabour != 0 || got.Rate != 0 || got.RateBasisPoints != 0 {
		t.Errorf("got %+v, want zero surplus labour and rate", got)
	}
}

func TestReviseWorkingDay_Rejects(t *testing.T) {
	t.Parallel()
	t.Run("day shorter than necessary labour", func(t *testing.T) {
		t.Parallel()
		_, err := commodity.ReviseWorkingDay(seniorNecessary, 344)
		if !errors.Is(err, commodity.ErrWorkingDayBelowNecessary) {
			t.Errorf("err = %v, want ErrWorkingDayBelowNecessary", err)
		}
	})
	t.Run("negative day", func(t *testing.T) {
		t.Parallel()
		_, err := commodity.ReviseWorkingDay(seniorNecessary, -60)
		if !errors.Is(err, commodity.ErrWorkingDayBelowNecessary) {
			t.Errorf("err = %v, want ErrWorkingDayBelowNecessary", err)
		}
	})
	t.Run("no necessary labour", func(t *testing.T) {
		t.Parallel()
		_, err := commodity.ReviseWorkingDay(0, 690)
		if !errors.Is(err, commodity.ErrDivisionByZero) {
			t.Errorf("err = %v, want ErrDivisionByZero", err)
		}
	})
	t.Run("negative necessary labour", func(t *testing.T) {
		t.Parallel()
		_, err := commodity.ReviseWorkingDay(-1, 690)
		if !errors.Is(err, commodity.ErrDivisionByZero) {
			t.Errorf("err = %v, want ErrDivisionByZero", err)
		}
	})
}

// A working-day cannot exceed the 24 hours of a natural day (Ch. 10). Beyond
// it the rate is meaningless, and a day near MaxInt64 would overflow the
// basis-point arithmetic into a plausible-looking wrong number.
func TestReviseWorkingDay_NaturalDayBound(t *testing.T) {
	t.Parallel()
	t.Run("24 hours is the limit and is accepted", func(t *testing.T) {
		t.Parallel()
		got, err := commodity.ReviseWorkingDay(seniorNecessary, 24*60)
		if err != nil {
			t.Fatalf("ReviseWorkingDay: %v", err)
		}
		if got.SurplusLabour != 24*60-345 {
			t.Errorf("SurplusLabour = %d, want %d", got.SurplusLabour, 24*60-345)
		}
	})
	for _, day := range []commodity.LabourMinutes{24*60 + 1, 1_000_000_000_000_000_000, math.MaxInt64} {
		_, err := commodity.ReviseWorkingDay(seniorNecessary, day)
		if !errors.Is(err, commodity.ErrWorkingDayExceedsNaturalDay) {
			t.Errorf("day %d: err = %v, want ErrWorkingDayExceedsNaturalDay", day, err)
		}
	}
	// A necessary labour that is itself beyond a natural day cannot be right.
	if _, err := commodity.ReviseWorkingDay(2000, 1500); !errors.Is(err, commodity.ErrWorkingDayExceedsNaturalDay) {
		t.Errorf("err = %v, want ErrWorkingDayExceedsNaturalDay", err)
	}
}
