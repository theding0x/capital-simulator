package agent

import (
	"errors"
	"testing"
)

func TestLabourProductivity_Validate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		lp      LabourProductivity
		wantErr error
	}{
		{"normal", LabourProductivity{Factor: 1.0}, nil},
		{"rise", LabourProductivity{Factor: 2.0}, nil},
		{"fall", LabourProductivity{Factor: 0.75}, nil},
		{"zero rejected", LabourProductivity{Factor: 0}, ErrProductivityFactorNonPositive},
		{"negative rejected", LabourProductivity{Factor: -1}, ErrProductivityFactorNonPositive},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got := c.lp.Validate()
			if !errors.Is(got, c.wantErr) {
				t.Fatalf("Validate(%v) = %v, want %v", c.lp, got, c.wantErr)
			}
		})
	}
}

// Marx Ch.17 §1: "the value of the labour-power falls from 4 shillings
// to 3, i.e. by 1/4 or 25%, the surplus-value rises from 2 shillings to
// 3, i.e. by 1/2 or 50%."
//
// In labour-time terms with a 12-hour working day: necessary falls from
// 480 min (4 sh) to 360 min (3 sh) when productivity rises by a factor of
// 4/3. Verifies inverse proportionality with rounding tolerance.
func TestLabourProductivity_ApplyToNecessaryLabour_MarxExample(t *testing.T) {
	t.Parallel()
	const necessaryBefore LabourMinutes = 480 // 4 sh
	productivity := LabourProductivity{Factor: 4.0 / 3.0}
	got := productivity.ApplyToNecessaryLabour(necessaryBefore)
	// 480 / (4/3) = 360
	if got != 360 {
		t.Fatalf("necessary after productivity rise: got %d, want 360", got)
	}
}

func TestLabourProductivity_ApplyToNecessaryLabour_NoChange(t *testing.T) {
	t.Parallel()
	got := NormalLabourProductivity.ApplyToNecessaryLabour(360)
	if got != 360 {
		t.Fatalf("normal productivity must be identity: got %d, want 360", got)
	}
}
