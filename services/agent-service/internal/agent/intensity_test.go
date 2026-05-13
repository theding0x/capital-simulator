package agent

import (
	"errors"
	"testing"
)

func TestLabourIntensity_NormalIsFactorOne(t *testing.T) {
	t.Parallel()
	if !NormalLabourIntensity.IsNormal() {
		t.Fatal("NormalLabourIntensity must satisfy IsNormal()")
	}
	if NormalLabourIntensity.Factor != 1.0 {
		t.Fatalf("normal factor: got %v, want 1.0", NormalLabourIntensity.Factor)
	}
}

func TestLabourIntensity_Validate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		li      LabourIntensity
		wantErr error
	}{
		{"normal", LabourIntensity{Factor: 1.0}, nil},
		// Marx Ch.17 §2: "a working-day of more intense labour is embodied
		// in more products" — factors > 1.0 model intensified labour.
		{"intensified 4/3", LabourIntensity{Factor: 1.333}, nil},
		{"zero rejected", LabourIntensity{Factor: 0}, ErrIntensityFactorNonPositive},
		{"negative rejected", LabourIntensity{Factor: -0.5}, ErrIntensityFactorNonPositive},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got := c.li.Validate()
			if !errors.Is(got, c.wantErr) {
				t.Fatalf("Validate(%v) = %v, want %v", c.li, got, c.wantErr)
			}
		})
	}
}
