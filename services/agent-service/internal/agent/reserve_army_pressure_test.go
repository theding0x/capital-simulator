package agent

import (
	"errors"
	"testing"
)

// Marx, Ch. 25 §3: the relative surplus-population (the industrial reserve
// army) regulates the general movement of wages. The fixtures below take a
// labour force of 1,000,000 and vary the reserve army measured against it.

func TestReserveArmyPressure_PressureBP(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		magnitude int64
		workforce int64
		wantBP    int64
	}{
		{"full employment (no reserve)", 0, 1_000_000, 0},
		{"unset workforce is inactive", 200_000, 0, 0},
		{"one fifth in reserve", 200_000, 1_000_000, 2000},
		{"one quarter in reserve", 250_000, 1_000_000, 2500},
		{"half in reserve", 500_000, 1_000_000, 5000},
		{"reserve as large as the whole force", 1_000_000, 1_000_000, 10000},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p := ReserveArmyPressure{Magnitude: tc.magnitude, WorkforceTotal: tc.workforce}
			if got := p.PressureBP(); got != tc.wantBP {
				t.Fatalf("PressureBP() = %d, want %d", got, tc.wantBP)
			}
		})
	}
}

func TestReserveArmyPressure_CompressPence(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		magnitude int64
		workforce int64
		nominal   int64
		want      int64
	}{
		{"full employment leaves the wage at its value", 0, 1_000_000, 30, 30},
		{"a fifth in reserve shaves a fifth off the wage", 200_000, 1_000_000, 30, 24},
		{"half in reserve halves the money-wage", 500_000, 1_000_000, 30, 15},
		{"a reserve the size of the force drives the wage to zero", 1_000_000, 1_000_000, 30, 0},
		{"a non-positive wage is untouched", 500_000, 1_000_000, 0, 0},
		{"half-pence rounds up (1.5 -> 2)", 250_000, 1_000_000, 2, 2},
		{"half-pence rounds up (0.5 -> 1)", 750_000, 1_000_000, 2, 1},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p := ReserveArmyPressure{Magnitude: tc.magnitude, WorkforceTotal: tc.workforce}
			if got := p.CompressPence(tc.nominal); got != tc.want {
				t.Fatalf("CompressPence(%d) = %d, want %d", tc.nominal, got, tc.want)
			}
		})
	}
}

// CompressMinutes shows the price of labour pushed below the value of
// labour-power: a 300-minute (5 h) value of labour-power compresses to 225
// minutes once a quarter of the workforce sits in the reserve army.
func TestReserveArmyPressure_CompressMinutes(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		magnitude int64
		workforce int64
		value     LabourMinutes
		want      LabourMinutes
	}{
		{"value of labour-power held when fully employed", 0, 1_000_000, 300, 300},
		{"a quarter in reserve pushes 300 below its value to 225", 250_000, 1_000_000, 300, 225},
		{"half in reserve halves an 8-hour value", 500_000, 1_000_000, 480, 240},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p := ReserveArmyPressure{Magnitude: tc.magnitude, WorkforceTotal: tc.workforce}
			got := p.CompressMinutes(tc.value)
			if got != tc.want {
				t.Fatalf("CompressMinutes(%d) = %d, want %d", tc.value, got, tc.want)
			}
			if got > tc.value {
				t.Fatalf("compression must not raise the wage: got %d > value %d", got, tc.value)
			}
		})
	}
}

func TestReserveArmyPressure_Validate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name      string
		magnitude int64
		workforce int64
		wantErr   error
	}{
		{"valid partition", 200_000, 1_000_000, nil},
		{"empty pressure is valid", 0, 0, nil},
		{"reserve exceeds workforce", 1_200_000, 1_000_000, ErrReserveArmyExceedsWorkforce},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := ReserveArmyPressure{Magnitude: tc.magnitude, WorkforceTotal: tc.workforce}.Validate()
			if tc.wantErr == nil && err != nil {
				t.Fatalf("Validate() = %v, want nil", err)
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Fatalf("Validate() = %v, want %v", err, tc.wantErr)
			}
		})
	}
}

func TestReserveArmyPressure_IsActive(t *testing.T) {
	t.Parallel()
	if (ReserveArmyPressure{Magnitude: 0, WorkforceTotal: 1_000_000}).IsActive() {
		t.Fatal("full employment must be inactive")
	}
	if (ReserveArmyPressure{Magnitude: 100, WorkforceTotal: 0}).IsActive() {
		t.Fatal("an unset workforce must be inactive")
	}
	if !(ReserveArmyPressure{Magnitude: 100, WorkforceTotal: 1_000}).IsActive() {
		t.Fatal("a real reserve against a real workforce must be active")
	}
}
