// Package rent — unit tests for Vol. III Ch. 44 worst-soil rent domain.
package rent

import (
	"testing"
)

func TestComputeWorstSoilRentPerAcre(t *testing.T) {
	t.Parallel()

	cases := []struct {
		old  int64
		new  int64
		want int64
	}{
		{30000, 35000, 5000},
		{35000, 30000, 0},
		{30000, 30000, 0},
	}
	for _, c := range cases {
		got := ComputeWorstSoilRentPerAcre(c.old, c.new)
		if got != c.want {
			t.Errorf("ComputeWorstSoilRentPerAcre(%d, %d) = %d, want %d",
				c.old, c.new, got, c.want)
		}
	}
}

func TestRentEmergenceMechanismIsValid(t *testing.T) {
	t.Parallel()

	valid := []RentEmergenceMechanism{
		MechanismBetterSoilRegulates,
		MechanismIncreasingProductivity,
		MechanismDecreasingProductivity,
	}
	for _, m := range valid {
		if !m.IsValid() {
			t.Errorf("mechanism %d should be valid", m)
		}
	}

	invalid := []RentEmergenceMechanism{0, 4}
	for _, m := range invalid {
		if m.IsValid() {
			t.Errorf("mechanism %d should be invalid", m)
		}
	}
}

func TestSoilImprovementCapitalIsAmortised(t *testing.T) {
	t.Parallel()

	cases := []struct {
		amortised int
		years     int
		want      bool
	}{
		{15, 15, true},
		{14, 15, false},
		{16, 15, true},
	}
	for _, c := range cases {
		s := SoilImprovementCapital{
			AmortisedYears:    c.amortised,
			AmortisationYears: c.years,
		}
		got := s.IsAmortised()
		if got != c.want {
			t.Errorf("AmortisedYears=%d AmortisationYears=%d: IsAmortised() = %v, want %v",
				c.amortised, c.years, got, c.want)
		}
	}
}

func TestCh44NewIDsDistinct(t *testing.T) {
	t.Parallel()

	a1, a2 := NewRentOnWorstSoilID(), NewRentOnWorstSoilID()
	if len(string(a1)) != 24 {
		t.Errorf("NewRentOnWorstSoilID len = %d, want 24", len(string(a1)))
	}
	if a1 == a2 {
		t.Error("two NewRentOnWorstSoilID calls must differ")
	}

	b1, b2 := NewSoilImprovementCapitalID(), NewSoilImprovementCapitalID()
	if len(string(b1)) != 24 {
		t.Errorf("NewSoilImprovementCapitalID len = %d, want 24", len(string(b1)))
	}
	if b1 == b2 {
		t.Error("two NewSoilImprovementCapitalID calls must differ")
	}

	c1, c2 := NewRegulatingPriceShiftID(), NewRegulatingPriceShiftID()
	if len(string(c1)) != 24 {
		t.Errorf("NewRegulatingPriceShiftID len = %d, want 24", len(string(c1)))
	}
	if c1 == c2 {
		t.Error("two NewRegulatingPriceShiftID calls must differ")
	}
}
