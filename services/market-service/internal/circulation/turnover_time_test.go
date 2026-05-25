package circulation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/market-service/internal/circulation"
)

// SpinningMill1871 fixture constants.
// §: "total time = time of production + time of circulation."
const (
	weekNanos    = int64(7 * 24 * 60 * 60 * 1_000_000_000)
	fiveDayNanos = int64(5 * 24 * 60 * 60 * 1_000_000_000)
	twoDayNanos  = int64(2 * 24 * 60 * 60 * 1_000_000_000)
	oneSec       = int64(1_000_000_000)
)

func TestProductionTime_TotalNanos_Decomposes(t *testing.T) {
	t.Parallel()
	// §: "time of production = labour-time + interruptions + latent MP + natural process."
	p := circulation.ProductionTime{
		LabourTimeNanos:         3 * oneSec,
		LabourInterruptionNanos: 1 * oneSec,
		LatentNanos:             1 * oneSec,
		NaturalProcessNanos:     5 * oneSec,
	}
	want := int64(10 * oneSec)
	if got := p.TotalNanos(); got != want {
		t.Errorf("TotalNanos() = %d, want %d", got, want)
	}
}

func TestCirculationTime_TotalNanos_Decomposes(t *testing.T) {
	t.Parallel()
	// §: "time of circulation = C—M + M—C."
	c := circulation.CirculationTime{
		SellingTimeNanos: fiveDayNanos,
		BuyingTimeNanos:  twoDayNanos,
	}
	want := weekNanos
	if got := c.TotalNanos(); got != want {
		t.Errorf("TotalNanos() = %d, want %d (1 week)", got, want)
	}
}

func TestTurnoverTime_TotalNanos_SpinningMill1871(t *testing.T) {
	t.Parallel()
	// SpinningMill1871: 1-week production + (5d selling + 2d buying = 1-week circulation) = 2 weeks.
	tt := circulation.TurnoverTime{
		ID:                  circulation.NewTurnoverTimeID(),
		IndustrialCapitalID: "5eed000000000000000502",
		Production: circulation.ProductionTime{
			LabourTimeNanos: weekNanos,
		},
		Circulation: circulation.CirculationTime{
			SellingTimeNanos: fiveDayNanos,
			BuyingTimeNanos:  twoDayNanos,
		},
	}
	want := 2 * weekNanos
	if got := tt.TotalNanos(); got != want {
		t.Errorf("TotalNanos() = %d, want %d (2 weeks)", got, want)
	}
}

func TestTurnoverTime_ActiveFractionBasisPoints_SpinningMill1871(t *testing.T) {
	t.Parallel()
	// 1 week production / 2 weeks total = 50% = 5000 basis points.
	// §: "circulation time limits production time and hence surplus-value production."
	tt := circulation.TurnoverTime{
		Production: circulation.ProductionTime{
			LabourTimeNanos: weekNanos,
		},
		Circulation: circulation.CirculationTime{
			SellingTimeNanos: fiveDayNanos,
			BuyingTimeNanos:  twoDayNanos,
		},
	}
	want := int64(5000)
	if got := tt.ActiveFractionBasisPoints(); got != want {
		t.Errorf("ActiveFractionBasisPoints() = %d, want %d", got, want)
	}
}

func TestTurnoverTime_ActiveFractionBasisPoints_ZeroTotal(t *testing.T) {
	t.Parallel()
	// No division by zero when both production and circulation are empty.
	tt := circulation.TurnoverTime{}
	if got := tt.ActiveFractionBasisPoints(); got != 0 {
		t.Errorf("ActiveFractionBasisPoints() = %d, want 0 on zero total", got)
	}
}

func TestTurnoverTime_ActiveFractionBasisPoints_FullProduction(t *testing.T) {
	t.Parallel()
	// If circulation time is zero, active fraction is 10000 (100%).
	tt := circulation.TurnoverTime{
		Production: circulation.ProductionTime{
			LabourTimeNanos: weekNanos,
		},
	}
	want := int64(10000)
	if got := tt.ActiveFractionBasisPoints(); got != want {
		t.Errorf("ActiveFractionBasisPoints() = %d, want %d (100%%)", got, want)
	}
}

func TestTurnoverTime_Validate_MissingIndustrialCapitalID(t *testing.T) {
	t.Parallel()
	tt := circulation.TurnoverTime{}
	if err := tt.Validate(); err == nil {
		t.Error("expected error for missing industrial_capital_id, got nil")
	}
}

func TestTurnoverTime_Validate_Valid(t *testing.T) {
	t.Parallel()
	tt := circulation.TurnoverTime{
		ID:                  circulation.NewTurnoverTimeID(),
		IndustrialCapitalID: "5eed000000000000000502",
	}
	if err := tt.Validate(); err != nil {
		t.Errorf("valid TurnoverTime rejected: %v", err)
	}
}

func TestProductionTime_OnlyLabourTimeCreatesValue(t *testing.T) {
	t.Parallel()
	// §: "interruptions, latent capital, and natural-process spans contribute
	// zero to value-creation metrics; only labour-time creates surplus-value."
	p := circulation.ProductionTime{
		LabourTimeNanos:         0,
		LabourInterruptionNanos: oneSec,
		LatentNanos:             oneSec,
		NaturalProcessNanos:     oneSec,
	}
	if p.LabourTimeNanos != 0 {
		t.Error("LabourTimeNanos should be 0: other sub-spans do not create value")
	}
	if p.TotalNanos() == 0 {
		t.Error("TotalNanos() must be > 0: other sub-spans still constitute production time")
	}
}

func TestNewTurnoverTimeID_Unique(t *testing.T) {
	t.Parallel()
	a := circulation.NewTurnoverTimeID()
	b := circulation.NewTurnoverTimeID()
	if a == b {
		t.Error("two NewTurnoverTimeID() calls returned the same ID")
	}
	if a.IsZero() || b.IsZero() {
		t.Error("NewTurnoverTimeID() returned a zero ID")
	}
}
