package costs_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
	"github.com/theding0x/capital-simulator/services/market-service/internal/costs"
)

// Compile-time assertion: costs.Pence is an alias of money.Pence (ADR 0001).
var _ money.Pence = costs.Pence(0)

// The 240/£ anchor for Ch. 6 is the in-package test costs_test.go (Pence(2640)
// // £11 x 240, Pence(480) // £2 x 240), not seed 00006. The 24 d tariff is
// scale-neutral; this test exercises real domain code: 24 d x 3 breakage
// multiplier = 72 d = 6s in canonical money.Pence.
func TestTransportTariff_EffectiveRate_IsMoneyPence(t *testing.T) {
	t.Parallel()
	tt := costs.TransportTariff{BasePencePerTonMile: 24, BreakageRiskMultiplierBasisPoints: 30000}
	if got, want := tt.EffectiveRatePencePerTonMile(), money.FromLSD(0, 6, 0); got != want {
		t.Errorf("3x breakage tariff = %d, want %d", got, want)
	}
}
