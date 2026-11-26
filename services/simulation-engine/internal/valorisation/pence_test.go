package valorisation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/pkg/money"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/valorisation"
)

// The valorisation scale is already canonical (240 pence to the pound), so the
// alias to money.Pence must leave every recorded magnitude unchanged. Comparing
// valorisation.Pence to money.Pence without a conversion only compiles while
// Pence is an alias, so this test fails to build without it.
func TestPenceIsCanonicalMoney(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		got         valorisation.Pence
		want        money.Pence
		wantMinutes int64
	}{
		{"Marx's 500 pound capital (seed 00051)", valorisation.Pence(120000), money.FromLSD(500, 0, 0), 500 * money.LabourMinutesPerPound},
		{"100 pounds a week variable advance", valorisation.Pence(24000), money.FromLSD(100, 0, 0), 100 * money.LabourMinutesPerPound},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.got != tc.want {
				t.Fatalf("got %d pence, want %d", tc.got, tc.want)
			}
			if m := tc.got.ToLabourMinutes(); m != tc.wantMinutes {
				t.Fatalf("labour minutes: got %d, want %d", m, tc.wantMinutes)
			}
		})
	}
}
