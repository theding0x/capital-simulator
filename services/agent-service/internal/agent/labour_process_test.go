package agent_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// Shilling conversion: 3 shillings = 6 hours = 360 minutes → 1 shilling = 120 LabourMinutes (§2.a).
const (
	necessaryLabour = agent.LabourMinutes(360) // §2.a: 6 hours = 3 shillings
	workingDay      = agent.LabourMinutes(720) // §2.b: 12-hour working day
)

// §1.d: zero-duration run must fail validation
func TestLabourProcess_Validate_ZeroDuration(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Duration:     0,
	}
	if err := lp.Validate(); err == nil {
		t.Error("zero-duration LabourProcess should fail validation")
	}
}

// §1.d: missing worker_id must fail validation
func TestLabourProcess_Validate_MissingWorker(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		CapitalistID: agent.NewAgentID(),
		Duration:     360,
	}
	if err := lp.Validate(); err == nil {
		t.Error("LabourProcess without worker_id should fail validation")
	}
}

// §1.d: missing capitalist_id must fail validation
func TestLabourProcess_Validate_MissingCapitalist(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		WorkerID: agent.NewAgentID(),
		Duration: 360,
	}
	if err := lp.Validate(); err == nil {
		t.Error("LabourProcess without capitalist_id should fail validation")
	}
}

// §1.c: valid process with all three elements passes validation
func TestLabourProcess_Validate_Valid(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Means: agent.MeansOfProduction{
			RawMaterials: []agent.RawMaterial{
				{CommodityID: "cotton", Quantity: 10, SNLTPerUnit: 120},
			},
			Instruments: []agent.Instrument{
				{CommodityID: "spindle", WearPerRun: 240},
			},
		},
		Duration: workingDay,
	}
	if err := lp.Validate(); err != nil {
		t.Errorf("valid LabourProcess should pass: %v", err)
	}
}

// §2.b: SurplusLabour(720, 360) == 360
func TestSurplusLabour_TwelveHourDay(t *testing.T) {
	t.Parallel()
	got := agent.SurplusLabour(workingDay, necessaryLabour)
	const want = agent.LabourMinutes(360)
	if got != want {
		t.Errorf("SurplusLabour(720, 360): want %d, got %d", want, got)
	}
}

// Invariant: SurplusLabour(wd, nl) == wd - nl
func TestSurplusLabour_Invariant(t *testing.T) {
	t.Parallel()
	cases := [][2]agent.LabourMinutes{{720, 360}, {480, 240}, {600, 300}, {360, 360}}
	for _, c := range cases {
		wd, nl := c[0], c[1]
		got := agent.SurplusLabour(wd, nl)
		want := wd - nl
		if got != want {
			t.Errorf("SurplusLabour(%d, %d): want %d, got %d", wd, nl, want, got)
		}
	}
}

// §2.c: TransferredValue = cotton(10*120) + spindle(240) = 1440
func TestTransferredValue_FifteenShillingsRun(t *testing.T) {
	t.Parallel()
	mp := agent.MeansOfProduction{
		RawMaterials: []agent.RawMaterial{{CommodityID: "cotton", Quantity: 10, SNLTPerUnit: 120}},
		Instruments:  []agent.Instrument{{CommodityID: "spindle", WearPerRun: 240}},
	}
	got := agent.TransferredValue(mp)
	const want = agent.LabourMinutes(1440) // 10*120 + 240
	if got != want {
		t.Errorf("TransferredValue: want %d, got %d", want, got)
	}
}

// Invariant: TransferredValue >= 0 for empty means
func TestTransferredValue_EmptyMeans(t *testing.T) {
	t.Parallel()
	if got := agent.TransferredValue(agent.MeansOfProduction{}); got < 0 {
		t.Errorf("TransferredValue of empty means should be >= 0, got %d", got)
	}
}

// §2.c: ProductValue = TransferredValue(1440) + ValueAdded(360) = 1800 (= 15 shillings * 120)
func TestProductValue_FifteenShillings(t *testing.T) {
	t.Parallel()
	mp := agent.MeansOfProduction{
		RawMaterials: []agent.RawMaterial{{CommodityID: "cotton", Quantity: 10, SNLTPerUnit: 120}},
		Instruments:  []agent.Instrument{{CommodityID: "spindle", WearPerRun: 240}},
	}
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Means:        mp,
		Duration:     necessaryLabour, // 6-hour day: only necessary labour (§2.c scenario)
	}
	vp := agent.ValorizationProcess{Process: lp, NecessaryLabourMinutes: necessaryLabour}
	got := vp.ProductValue()
	const want = agent.LabourMinutes(1800) // 1440 + 360
	if got != want {
		t.Errorf("ProductValue: want %d, got %d", want, got)
	}
}

// Invariant: ProductValue == TransferredValue + ValueAdded
func TestProductValue_Invariant(t *testing.T) {
	t.Parallel()
	mp := agent.MeansOfProduction{
		RawMaterials: []agent.RawMaterial{{CommodityID: "cotton", Quantity: 10, SNLTPerUnit: 120}},
		Instruments:  []agent.Instrument{{CommodityID: "spindle", WearPerRun: 240}},
	}
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Means:        mp,
		Duration:     workingDay,
	}
	vp := agent.ValorizationProcess{Process: lp, NecessaryLabourMinutes: necessaryLabour}
	want := agent.TransferredValue(mp) + agent.ValueAdded(workingDay)
	if vp.ProductValue() != want {
		t.Errorf("ProductValue invariant: want %d, got %d", want, vp.ProductValue())
	}
}

// Invariant: ValueAdded(d) == d for uniform skill level = 1
func TestValueAdded_Identity(t *testing.T) {
	t.Parallel()
	for _, d := range []agent.LabourMinutes{0, 60, 360, 720} {
		if got := agent.ValueAdded(d); got != d {
			t.Errorf("ValueAdded(%d): want %d, got %d", d, d, got)
		}
	}
}

// §2.e: ValueAdded(workingDay) > NecessaryLabour when WorkingDay > NecessaryLabour
func TestValueAdded_ExceedsNecessaryLabour(t *testing.T) {
	t.Parallel()
	added := agent.ValueAdded(workingDay)
	if added <= necessaryLabour {
		t.Errorf("ValueAdded(%d)=%d should exceed NecessaryLabour %d", workingDay, added, necessaryLabour)
	}
}

// Invariant: SurplusValue == 0 iff WorkingDay == NecessaryLabour
func TestValorizationProcess_NoSurplusWhenEqual(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Duration:     necessaryLabour,
	}
	vp := agent.ValorizationProcess{Process: lp, NecessaryLabourMinutes: necessaryLabour}
	if sv := vp.SurplusValue(); sv != 0 {
		t.Errorf("SurplusValue when wd==nl: want 0, got %d", sv)
	}
}

// §2.f: SurplusValue > 0 whenever WorkingDay > NecessaryLabour
func TestValorizationProcess_SurplusValuePositive(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Duration:     workingDay,
	}
	vp := agent.ValorizationProcess{Process: lp, NecessaryLabourMinutes: necessaryLabour}
	if vp.SurplusValue() <= 0 {
		t.Errorf("SurplusValue: want > 0, got %d", vp.SurplusValue())
	}
}

// Invariant: NecessaryLabour + SurplusLabour == WorkingDay
func TestValorizationProcess_WorkingDayPartition(t *testing.T) {
	t.Parallel()
	lp := agent.LabourProcess{
		WorkerID:     agent.NewAgentID(),
		CapitalistID: agent.NewAgentID(),
		Duration:     workingDay,
	}
	vp := agent.ValorizationProcess{Process: lp, NecessaryLabourMinutes: necessaryLabour}
	if vp.NecessaryLabour()+vp.SurplusLabour() != workingDay {
		t.Errorf("WorkingDay partition invariant: %d + %d != %d",
			vp.NecessaryLabour(), vp.SurplusLabour(), workingDay)
	}
}

// NewLabourProcessID produces a 24-char hex string
func TestNewLabourProcessID(t *testing.T) {
	t.Parallel()
	id := agent.NewLabourProcessID()
	if id.IsZero() {
		t.Error("NewLabourProcessID should not be zero")
	}
	if len(string(id)) != 24 {
		t.Errorf("LabourProcessID: want 24 chars, got %d", len(string(id)))
	}
}
