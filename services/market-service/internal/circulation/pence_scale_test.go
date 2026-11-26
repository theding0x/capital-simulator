package circulation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
	"github.com/theding0x/capital-simulator/services/market-service/internal/circulation"
)

// The package Pence is money.Pence (ADR 0001): Marx's £-amounts built with
// money.FromLSD equal the integers the seeds already store.
func TestPence_IsCanonicalMoneyPence(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want money.Pence
		got  circulation.Pence
	}{
		{"£50 cotton stockpile, Spinning Mill 1871", money.FromLSD(50, 0, 0), 12000},
		{"£1", money.FromLSD(1, 0, 0), 240},
		{"3s 6d", money.FromLSD(0, 3, 6), 42},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.got != tc.want {
				t.Errorf("got %d, want %d", tc.got, tc.want)
			}
		})
	}
}

func TestLatentProductiveCapital_Pence_IsMoneyPence(t *testing.T) {
	t.Parallel()
	lpc := circulation.LatentProductiveCapital{Pence: money.FromLSD(50, 0, 0)}
	if lpc.Pence != 12000 {
		t.Errorf("£50 cotton = %d pence, want 12000", lpc.Pence)
	}
}
