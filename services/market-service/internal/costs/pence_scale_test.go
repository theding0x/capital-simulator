package costs_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
	"github.com/theding0x/capital-simulator/services/market-service/internal/costs"
)

// Scale evidence for Ch. 6: Marx's Birmingham glass crates (£11 historic,
// £2 modern) are held as 2640 and 480, and the 24 d ton-mile tariff is literal
// pence, i.e. 240 pence to the pound, identical to money.Pence.
func TestPence_IsCanonicalMoneyPence(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		want money.Pence
		got  costs.Pence
	}{
		{"historic glass crate £11", money.FromLSD(11, 0, 0), 2640},
		{"modern glass crate £2", money.FromLSD(2, 0, 0), 480},
		{"24 d per ton-mile tariff = 2s", money.FromLSD(0, 2, 0), 24},
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

func TestTransportTariff_EffectiveRate_IsMoneyPence(t *testing.T) {
	t.Parallel()
	tt := costs.TransportTariff{BasePencePerTonMile: 24, BreakageRiskMultiplierBasisPoints: 30000}
	if got, want := tt.EffectiveRatePencePerTonMile(), money.FromLSD(0, 6, 0); got != want {
		t.Errorf("3x breakage tariff = %d, want %d", got, want)
	}
}
