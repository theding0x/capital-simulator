package circulation

import (
	"errors"
	"testing"
)

// Vol. II Ch. 3 — money-branching at M′. Fixture: Spinning Mill 1871,
// C′ = £500 (c=£372, v=£50, s=£78).

func spinningMill1871Circuit(mode ReproductionMode) CommodityCircuit {
	return CommodityCircuit{
		ID: "mill-1871",
		Initial: OpeningCommodityCapital{
			ConstantPence: 37200,
			VariablePence: 5000,
			SurplusPence:  7800,
			PoundsTotal:   10000,
		},
		Mode: mode,
	}
}

func TestVol2Ch03_BranchMoney_SimpleReproduction_AllSurplusToRevenue(t *testing.T) {
	t.Parallel()
	b, err := spinningMill1871Circuit(ReproductionSimple).BranchMoney(7800)
	if err != nil {
		t.Fatalf("BranchMoney() = %v, want nil", err)
	}
	if b.ReproductionFundPence != 42200 || b.RevenuePence != 7800 || b.CapitalisedPence != 0 {
		t.Errorf("got fund=%d revenue=%d capitalised=%d, want 42200/7800/0",
			b.ReproductionFundPence, b.RevenuePence, b.CapitalisedPence)
	}
	if b.RealisedPence != 50000 {
		t.Errorf("RealisedPence = %d, want 50000", b.RealisedPence)
	}
}

func TestVol2Ch03_BranchMoney_ExtendedReproduction_AllSurplusCapitalised(t *testing.T) {
	t.Parallel()
	b, err := spinningMill1871Circuit(ReproductionExtended).BranchMoney(0)
	if err != nil {
		t.Fatalf("BranchMoney() = %v, want nil", err)
	}
	if b.ReproductionFundPence != 42200 || b.RevenuePence != 0 || b.CapitalisedPence != 7800 {
		t.Errorf("got fund=%d revenue=%d capitalised=%d, want 42200/0/7800",
			b.ReproductionFundPence, b.RevenuePence, b.CapitalisedPence)
	}
}

func TestVol2Ch03_BranchMoney_Mixed_SplitsSurplusAndConserves(t *testing.T) {
	t.Parallel()
	b, err := spinningMill1871Circuit(ReproductionMixed).BranchMoney(3000)
	if err != nil {
		t.Fatalf("BranchMoney() = %v, want nil", err)
	}
	if b.RevenuePence != 3000 || b.CapitalisedPence != 4800 {
		t.Errorf("got revenue=%d capitalised=%d, want 3000/4800", b.RevenuePence, b.CapitalisedPence)
	}
	if sum := b.ReproductionFundPence + b.RevenuePence + b.CapitalisedPence; sum != b.RealisedPence {
		t.Errorf("branches sum to %d, want realised %d", sum, b.RealisedPence)
	}
}

func TestVol2Ch03_BranchMoney_Errors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		mode    ReproductionMode
		revenue Pence
		want    error
	}{
		{"revenue exceeds surplus", ReproductionMixed, 7801, ErrBranchingExceedsSurplus},
		{"negative revenue", ReproductionMixed, -1, ErrBranchingNegative},
		{"simple must consume all surplus", ReproductionSimple, 3000, ErrBranchingModeMismatch},
		{"extended must consume none", ReproductionExtended, 100, ErrBranchingModeMismatch},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if _, err := spinningMill1871Circuit(tt.mode).BranchMoney(tt.revenue); !errors.Is(err, tt.want) {
				t.Errorf("BranchMoney(%d) = %v, want %v", tt.revenue, err, tt.want)
			}
		})
	}
}

func TestVol2Ch03_MoneyBranching_Validate(t *testing.T) {
	t.Parallel()
	if err := (MoneyBranching{RealisedPence: 50000, ReproductionFundPence: 42200, RevenuePence: 7000}).Validate(); !errors.Is(err, ErrBranchingUnbalanced) {
		t.Errorf("unbalanced: got %v, want ErrBranchingUnbalanced", err)
	}
	if err := (MoneyBranching{RealisedPence: 0, ReproductionFundPence: 100, RevenuePence: -100}).Validate(); !errors.Is(err, ErrBranchingNegative) {
		t.Errorf("negative: got %v, want ErrBranchingNegative", err)
	}
}
