package circulation_test

import (
	"errors"
	"math"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

// Spinning Mill 1871 (Vol. II Ch. 2 §I): C' is 10,000 lbs of yarn valued at
// £500, of which £422 (8,440 lbs) replaces the productive capital and £78
// (1,560 lbs) is surplus-value. In legacy pence: 42,200 + 7,800 = 50,000.
func spinningMillExit() *circulation.ExitCommodityCapital {
	return &circulation.ExitCommodityCapital{
		ProductiveCircuitID: "5eed000000000000020a",
		ValueOriginal:       42200,
		ValueSurplus:        7800,
	}
}

func TestExitCommodityCapital_Total(t *testing.T) {
	t.Parallel()
	if got := spinningMillExit().Total(); got != 50000 {
		t.Errorf("Total = %d, want 50000 (£500)", got)
	}
}

func TestExitCommodityCapital_Validate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		original circulation.Pence
		surplus  circulation.Pence
		wantErr  bool
	}{
		{"Spinning Mill 1871", 42200, 7800, false},
		{"no surplus (zero surplus-value is allowed)", 42200, 0, false},
		{"zero original", 0, 7800, true},
		{"negative original", -1, 7800, true},
		{"negative surplus", 42200, -1, true},
		// Total() must stay representable: a wrapped (negative) total would let a
		// one-penny sale pass the "sold in its entirety" guard.
		{"total exactly MaxInt64", math.MaxInt64 - 1, 1, false},
		{"total overflows int64", math.MaxInt64, 1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := circulation.ExitCommodityCapital{ValueOriginal: tt.original, ValueSurplus: tt.surplus}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// "C' must be sold at its value and in its entirety."
func TestValidateRealisation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		commodity *circulation.ExitCommodityCapital
		realised  circulation.Pence
		wantErr   error
	}{
		{"no C' recorded", nil, 50000, circulation.ErrNoCommodityCapital},
		{"capital-value yarn sold, surplus-value yarn unsold", spinningMillExit(), 42200, circulation.ErrPartialRealisation},
		{"one penny short", spinningMillExit(), 49999, circulation.ErrPartialRealisation},
		{"nothing sold", spinningMillExit(), 0, circulation.ErrPartialRealisation},
		{"sold at value in its entirety", spinningMillExit(), 50000, nil},
		{"sold above value (Ch. 15 concern, accepted here)", spinningMillExit(), 52000, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := circulation.ValidateRealisation(tt.commodity, tt.realised); !errors.Is(err, tt.wantErr) {
				t.Errorf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestProductiveCircuit_ValidateRealisation_MatchesFunctionForm(t *testing.T) {
	t.Parallel()
	for _, realised := range []circulation.Pence{0, 42200, 49999, 50000, 52000} {
		pc := circulation.ProductiveCircuit{ExitCommodity: spinningMillExit()}
		want := circulation.ValidateRealisation(pc.ExitCommodity, realised)
		if got := pc.ValidateRealisation(realised); !errors.Is(got, want) {
			t.Errorf("realised=%d: method = %v, function = %v", realised, got, want)
		}
	}
}

func TestProductiveCircuit_ValidateRealisation_NoCommodity(t *testing.T) {
	t.Parallel()
	pc := circulation.ProductiveCircuit{}
	if err := pc.ValidateRealisation(50000); !errors.Is(err, circulation.ErrNoCommodityCapital) {
		t.Errorf("got %v, want ErrNoCommodityCapital", err)
	}
}
