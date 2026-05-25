package circulation_test

import (
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/circulation"
)

func TestNaturalProcessKind_AllValuesDistinct(t *testing.T) {
	t.Parallel()
	// §: "grain sown, wine fermenting in the cellar, tanneries, drying."
	kinds := []circulation.NaturalProcessKind{
		circulation.NaturalRipening,
		circulation.NaturalFermentation,
		circulation.NaturalTanning,
		circulation.NaturalDrying,
		circulation.NaturalOther,
	}
	seen := map[circulation.NaturalProcessKind]bool{}
	for _, k := range kinds {
		if seen[k] {
			t.Errorf("duplicate NaturalProcessKind: %q", k)
		}
		seen[k] = true
	}
}

func TestNaturalProcessSpan_GrainRipening(t *testing.T) {
	t.Parallel()
	// §: "grain sown … the product requires a relatively long period of time."
	ninetyDayNanos := int64(90 * 24 * 60 * 60 * 1_000_000_000)
	nps := circulation.NaturalProcessSpan{
		ID:            circulation.NewNaturalProcessSpanID(),
		Process:       circulation.NaturalRipening,
		DurationNanos: ninetyDayNanos,
		StartedAt:     time.Date(1871, 3, 1, 0, 0, 0, 0, time.UTC),
	}
	if nps.Process != circulation.NaturalRipening {
		t.Errorf("Process = %q, want %q", nps.Process, circulation.NaturalRipening)
	}
	if nps.DurationNanos != ninetyDayNanos {
		t.Errorf("DurationNanos = %d, want %d", nps.DurationNanos, ninetyDayNanos)
	}
}

func TestNaturalProcessSpan_WineFermentation(t *testing.T) {
	t.Parallel()
	// §: "wine fermenting in the cellar."
	fourteenDayNanos := int64(14 * 24 * 60 * 60 * 1_000_000_000)
	nps := circulation.NaturalProcessSpan{
		ID:            circulation.NewNaturalProcessSpanID(),
		Process:       circulation.NaturalFermentation,
		DurationNanos: fourteenDayNanos,
	}
	if nps.DurationNanos != fourteenDayNanos {
		t.Errorf("DurationNanos = %d, want %d", nps.DurationNanos, fourteenDayNanos)
	}
}

func TestNaturalProcessSpan_LeatherTanning(t *testing.T) {
	t.Parallel()
	// §: "tanneries exposed to chemical processes."
	thirtyDayNanos := int64(30 * 24 * 60 * 60 * 1_000_000_000)
	nps := circulation.NaturalProcessSpan{
		ID:            circulation.NewNaturalProcessSpanID(),
		Process:       circulation.NaturalTanning,
		DurationNanos: thirtyDayNanos,
	}
	if nps.Process != circulation.NaturalTanning {
		t.Errorf("Process = %q, want %q", nps.Process, circulation.NaturalTanning)
	}
}

func TestNewNaturalProcessSpanID_Unique(t *testing.T) {
	t.Parallel()
	a := circulation.NewNaturalProcessSpanID()
	b := circulation.NewNaturalProcessSpanID()
	if a == b {
		t.Errorf("duplicate NaturalProcessSpanID: %q", a)
	}
	if a.IsZero() {
		t.Error("NewNaturalProcessSpanID() returned zero value")
	}
}

func TestLatentProductiveCapital_IsLatent_BeforeProduction(t *testing.T) {
	t.Parallel()
	// §: "fallow capital … creates neither products nor value."
	// £50 cotton stockpile at SpinningMill1871 (£50 × 240 = 12000 pence).
	lpc := circulation.LatentProductiveCapital{
		ID:                  circulation.NewLatentProductiveCapitalID(),
		IndustrialCapitalID: "5eed000000000000000502",
		Pence:               12000,
		HeldAt:              time.Date(1871, 1, 1, 0, 0, 0, 0, time.UTC),
		EnteredProductionAt: nil,
	}
	if !lpc.IsLatent() {
		t.Error("LatentProductiveCapital with nil EnteredProductionAt must be latent")
	}
}

func TestLatentProductiveCapital_IsLatent_AfterProduction(t *testing.T) {
	t.Parallel()
	entered := time.Date(1871, 1, 8, 0, 0, 0, 0, time.UTC)
	lpc := circulation.LatentProductiveCapital{
		ID:                  circulation.NewLatentProductiveCapitalID(),
		EnteredProductionAt: &entered,
	}
	if lpc.IsLatent() {
		t.Error("LatentProductiveCapital with non-nil EnteredProductionAt must not be latent")
	}
}

func TestNewLatentProductiveCapitalID_Unique(t *testing.T) {
	t.Parallel()
	a := circulation.NewLatentProductiveCapitalID()
	b := circulation.NewLatentProductiveCapitalID()
	if a == b {
		t.Errorf("duplicate LatentProductiveCapitalID: %q", a)
	}
}
