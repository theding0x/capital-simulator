package agent_test

import (
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

// Capital Vol. I, Ch. 13, § (para 2):
// "If a working-day of 12 hours be embodied in six shillings,
//
//	1,200 such days will be embodied in 1,200 times 6 shillings."
//
// 1,200 workers × 720 min/day = 864,000 minutes — value scales linearly with
// worker count when no qualitative cooperation bonus is applied.
func TestCollectiveWorkingDay_Para2_1200Workers(t *testing.T) {
	t.Parallel()
	got := agent.CollectiveWorkingDay(1200, agent.LabourMinutes(720))
	if got != agent.LabourMinutes(864_000) {
		t.Fatalf("CollectiveWorkingDay(1200, 720) = %d, want 864000", got)
	}
}

// Capital Vol. I, Ch. 13, § (para 3):
// "the collective working-day of 12 men simultaneously employed,
//
//	consists of 144 hours; and although the labour of each may deviate
//	more or less from average social labour ... it possesses the qualities
//	of an average social working-day."
//
// 12 men × 12 h = 144 h = 8640 min. Reducing back via AverageSocialLabour
// returns the per-worker average of 720 min (12 h).
func TestAverageSocialLabour_Para3_12Men144Hours(t *testing.T) {
	t.Parallel()
	got := agent.AverageSocialLabour(agent.LabourMinutes(8640), 12)
	if got != agent.LabourMinutes(720) {
		t.Fatalf("AverageSocialLabour(8640, 12) = %d, want 720", got)
	}
}

// Capital Vol. I, Ch. 13, § (para 3, Burke fn. 1):
// "any given five men will, in their total, afford a proportion of labour
//
//	equal to any other five."
//
// Minimum platoon size for deviation cancellation is 5.
func TestAverageSocialLabour_BurkeFootnote_FiveMen(t *testing.T) {
	t.Parallel()
	got := agent.AverageSocialLabour(agent.LabourMinutes(3600), 5)
	if got != agent.LabourMinutes(720) {
		t.Fatalf("AverageSocialLabour(3600, 5) = %d, want 720", got)
	}
	if agent.CooperationMinSize != 5 {
		t.Fatalf("CooperationMinSize = %d, want 5", agent.CooperationMinSize)
	}
}

// Capital Vol. I, Ch. 13, para 5: a dozen working together "produce far more
// than twelve isolated men". Marx gives no formula, so the gain is recorded on
// the cooperation, never derived from its size: 12 masons with nothing recorded
// assert no gain, and a recorded observation is reported as given.
func TestCooperation_PowerFactorIsRecordedNotDerived(t *testing.T) {
	t.Parallel()
	members := make([]agent.CooperationMember, 12)
	for i := range members {
		members[i] = agent.CooperationMember{WorkerID: agent.NewAgentID(), WorkingDayMinutes: 720}
	}
	unrecorded := agent.Cooperation{CapitalistID: agent.NewAgentID(), Members: members}
	if got := unrecorded.PowerFactor(); got != 1.0 {
		t.Fatalf("12 masons, nothing recorded: PowerFactor = %v, want exactly 1.0", got)
	}
	observed := unrecorded
	observed.ProductivePower = 1.2
	if got := observed.PowerFactor(); got != 1.2 {
		t.Fatalf("recorded 1.2: PowerFactor = %v, want 1.2", got)
	}
	if got := agent.CooperationProductivityFactor(observed); got != 1.2 {
		t.Fatalf("CooperationProductivityFactor = %v, want 1.2", got)
	}
}

// Co-operation never shrinks productive power: a recorded factor below 1 is
// rejected; 0 (unrecorded) and 1 or more are valid.
func TestCooperation_Validate_ProductivePower(t *testing.T) {
	t.Parallel()
	base := agent.Cooperation{
		CapitalistID: agent.NewAgentID(),
		Members:      []agent.CooperationMember{{WorkerID: agent.NewAgentID(), WorkingDayMinutes: 720}},
	}
	for _, tc := range []struct {
		factor float64
		want   error
	}{
		{0, nil}, {1, nil}, {1.2, nil},
		{0.5, agent.ErrCooperationProductivePower},
		{-1, agent.ErrCooperationProductivePower},
	} {
		c := base
		c.ProductivePower = tc.factor
		if err := c.Validate(); !errors.Is(err, tc.want) {
			t.Errorf("factor %v: Validate() = %v, want %v", tc.factor, err, tc.want)
		}
	}
}

// Capital Vol. I, Ch. 13, § (para 7):
// "100 men co-operating extend the working-day to 1,200 hours."
//
// 100 × 720 min = 72,000 min — time-critical tasks addressed by massing
// labour at the decisive moment.
func TestCollectiveWorkingDay_Para7_100Men(t *testing.T) {
	t.Parallel()
	got := agent.CollectiveWorkingDay(100, agent.LabourMinutes(720))
	if got != agent.LabourMinutes(72_000) {
		t.Fatalf("CollectiveWorkingDay(100, 720) = %d, want 72000", got)
	}
}

// Capital Vol. I, Ch. 13, § (para 8):
// "The payment of 300 workmen at once, even if only for a single day, requires
//
//	a greater outlay of capital than does the payment of a smaller number of
//	men, week by week."
//
// MinimumCapital(300) / MinimumCapital(10) = 30, matching Marx's ratio.
func TestMinimumCapital_Para8_300vs10Ratio(t *testing.T) {
	t.Parallel()
	wage := agent.Pence(720) // 6 shillings = 72 pence; pence-scaling is irrelevant for the ratio.
	mc300 := agent.MinimumCapital(300, wage)
	mc10 := agent.MinimumCapital(10, wage)
	if mc300 != agent.Pence(300)*wage {
		t.Fatalf("MinimumCapital(300, %d) = %d, want %d", wage, mc300, agent.Pence(300)*wage)
	}
	if mc10 != agent.Pence(10)*wage {
		t.Fatalf("MinimumCapital(10, %d) = %d, want %d", wage, mc10, agent.Pence(10)*wage)
	}
	if mc300/mc10 != 30 {
		t.Fatalf("ratio mc300/mc10 = %d, want 30", mc300/mc10)
	}
}

// Invariant: CollectiveWorkingDay(n, d) == n * d for any n, d.
// Para 2: value addition is strictly additive from the standpoint of value.
func TestInvariant_CollectiveWorkingDayIsLinear(t *testing.T) {
	t.Parallel()
	cases := []struct {
		n    agent.CooperationSize
		d    agent.LabourMinutes
		want agent.LabourMinutes
	}{
		{1, 720, 720},
		{12, 720, 8640},
		{100, 720, 72_000},
		{1200, 720, 864_000},
	}
	for _, tc := range cases {
		got := agent.CollectiveWorkingDay(tc.n, tc.d)
		if got != tc.want {
			t.Fatalf("CollectiveWorkingDay(%d, %d) = %d, want %d", tc.n, tc.d, got, tc.want)
		}
	}
}

// Invariant: AverageSocialLabour(CollectiveWorkingDay(n, d), n) == d for any n >= 1.
// Para 3: the social law holds for any n >= 1.
func TestInvariant_AverageSocialLabourReversesCollective(t *testing.T) {
	t.Parallel()
	for _, n := range []agent.CooperationSize{1, 5, 12, 100, 1200} {
		d := agent.LabourMinutes(720)
		collective := agent.CollectiveWorkingDay(n, d)
		got := agent.AverageSocialLabour(collective, n)
		if got != d {
			t.Fatalf("AverageSocialLabour(CollectiveWorkingDay(%d, %d), %d) = %d, want %d",
				n, d, n, got, d)
		}
	}
}

// Invariant: MinimumCapital(n, w) == n * w.
// Para 8: the capitalist must have at least n × dailyWage in pocket.
func TestInvariant_MinimumCapitalIsLinear(t *testing.T) {
	t.Parallel()
	wage := agent.Pence(50)
	for _, n := range []agent.CooperationSize{1, 10, 300, 1000} {
		got := agent.MinimumCapital(n, wage)
		want := agent.Pence(int64(n)) * wage
		if got != want {
			t.Fatalf("MinimumCapital(%d, %d) = %d, want %d", n, wage, got, want)
		}
	}
}

// CooperativeOrigin: collective power appears as a property of capital, not labour.
func TestCooperativeOrigin_AppearsAsCapital(t *testing.T) {
	t.Parallel()
	if got := agent.CooperativeOrigin(); got != "capital" {
		t.Fatalf("CooperativeOrigin() = %q, want %q", got, "capital")
	}
}

func TestCooperation_Validate(t *testing.T) {
	t.Parallel()
	capID := agent.NewAgentID()
	w1 := agent.NewAgentID()
	w2 := agent.NewAgentID()

	good := agent.Cooperation{
		Name:         "spinning floor",
		CapitalistID: capID,
		Members: []agent.CooperationMember{
			{WorkerID: w1, WorkingDayMinutes: agent.LabourMinutes(720)},
			{WorkerID: w2, WorkingDayMinutes: agent.LabourMinutes(720), Supervisory: true},
		},
	}
	if err := good.Validate(); err != nil {
		t.Fatalf("Validate() returned %v, want nil", err)
	}
	if good.Size() != 2 {
		t.Fatalf("Size() = %d, want 2", good.Size())
	}
	if sups := good.Supervisors(); len(sups) != 1 || sups[0].WorkerID != w2 {
		t.Fatalf("Supervisors() = %+v, want one supervisor with id %s", sups, w2)
	}

	tests := []struct {
		name string
		coop agent.Cooperation
		err  error
	}{
		{
			name: "no capitalist",
			coop: agent.Cooperation{
				Members: []agent.CooperationMember{{WorkerID: w1, WorkingDayMinutes: 720}},
			},
			err: agent.ErrCooperationCapitalistID,
		},
		{
			name: "no members",
			coop: agent.Cooperation{CapitalistID: capID, Members: nil},
			err:  agent.ErrCooperationNoMembers,
		},
		{
			name: "member missing worker id",
			coop: agent.Cooperation{
				CapitalistID: capID,
				Members:      []agent.CooperationMember{{WorkingDayMinutes: 720}},
			},
			err: agent.ErrCooperationMemberInvalid,
		},
		{
			name: "non-positive working day",
			coop: agent.Cooperation{
				CapitalistID: capID,
				Members:      []agent.CooperationMember{{WorkerID: w1, WorkingDayMinutes: 0}},
			},
			err: agent.ErrCooperationWorkingDay,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.coop.Validate(); err != tc.err {
				t.Fatalf("Validate() = %v, want %v", err, tc.err)
			}
		})
	}
}

func TestCooperation_AverageIndividualWorkingDay(t *testing.T) {
	t.Parallel()
	capID := agent.NewAgentID()
	coop := agent.Cooperation{
		CapitalistID: capID,
		Members: []agent.CooperationMember{
			{WorkerID: agent.NewAgentID(), WorkingDayMinutes: 720},
			{WorkerID: agent.NewAgentID(), WorkingDayMinutes: 720},
			{WorkerID: agent.NewAgentID(), WorkingDayMinutes: 720},
		},
	}
	if got := coop.AverageIndividualWorkingDay(); got != 720 {
		t.Fatalf("AverageIndividualWorkingDay() = %d, want 720", got)
	}
}
