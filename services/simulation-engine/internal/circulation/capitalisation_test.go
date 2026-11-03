package circulation_test

import (
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

// Spinning Mill 1871 (Vol. II Ch. 2): £78 of surplus-value hoarded for
// capitalisation, injected as 6500 constant / 1300 variable pence.
const (
	fixtureHoard circulation.Pence = 7800
	fixtureMin   circulation.Pence = 100
)

func TestValidateCapitalisation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		mode    circulation.ReproductionMode
		hoard   circulation.Pence
		amount  circulation.Pence
		dc, dv  circulation.Pence
		wantErr error
	}{
		{"whole hoard, extended", circulation.ReproductionExtended, fixtureHoard, 7800, 6500, 1300, nil},
		{"part of hoard, mixed", circulation.ReproductionMixed, fixtureHoard, 3000, 2500, 500, nil},
		{"exactly at hoard", circulation.ReproductionExtended, 100, 100, 100, 0, nil},
		{"exceeds hoard", circulation.ReproductionExtended, fixtureHoard, 7900, 7900, 0, circulation.ErrCapitalisationExceedsHoard},
		{"empty hoard", circulation.ReproductionExtended, 0, 100, 100, 0, circulation.ErrCapitalisationExceedsHoard},
		{"simple reproduction", circulation.ReproductionSimple, fixtureHoard, 100, 100, 0, circulation.ErrSimpleReproductionNoGrowth},
		{"not conserved", circulation.ReproductionExtended, fixtureHoard, 200, 100, 50, circulation.ErrCapitalisationNotConserved},
		{"below increment", circulation.ReproductionExtended, fixtureHoard, 50, 50, 0, circulation.ErrIndivisibleProductionElement},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pc := circulation.ProductiveCircuit{
				Mode:                       tt.mode,
				MinCapitalisationIncrement: fixtureMin,
				LatentMoneyCapital:         circulation.LatentMoneyCapital{Accumulated: tt.hoard},
			}
			if err := pc.ValidateCapitalisation(tt.amount, tt.dc, tt.dv); !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}
