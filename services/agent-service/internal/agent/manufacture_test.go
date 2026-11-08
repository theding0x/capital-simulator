package agent_test

import (
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

func carriageManufacture() agent.Manufacture {
	// §1 two-fold origin — combination: a carriage workshop assembles diverse
	// handicrafts under one capitalist.
	return agent.Manufacture{
		CapitalistID:                agent.NewAgentID(),
		Form:                        agent.HeterogeneousManufacture,
		Origin:                      agent.OriginCombination,
		IndividualWorkingDayMinutes: 720,
		Roles: []agent.DetailRole{
			{Name: "wheelwright", SkillLevel: agent.SkilledLevel, OutputRatePerHour: 4, HeadCount: 2, ToolName: "spokeshave"},
			{Name: "harness-maker", SkillLevel: agent.SkilledLevel, OutputRatePerHour: 4, HeadCount: 2, ToolName: "awl"},
			{Name: "locksmith", SkillLevel: agent.SkilledLevel, OutputRatePerHour: 4, HeadCount: 1, ToolName: "file"},
			{Name: "fitter", SkillLevel: agent.UnskilledLevel, OutputRatePerHour: 4, HeadCount: 3, ToolName: ""},
		},
	}
}

// §1 two-fold origin — combination: a carriage manufacture assembles
// wheelwrights, harness-makers, locksmiths, etc. under one capitalist.
func TestManufacture_Origin_Combination(t *testing.T) {
	t.Parallel()
	m := carriageManufacture()
	if err := m.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
	if m.Origin != agent.OriginCombination {
		t.Fatalf("Origin = %q, want %q", m.Origin, agent.OriginCombination)
	}
	if len(m.Roles) < 2 {
		t.Fatalf("len(Roles) = %d, want >= 2", len(m.Roles))
	}
}

// §1 two-fold origin — splitting: a needle manufacture subdivides one
// handicraft into 20+ sequential operations.
func TestManufacture_Origin_Splitting(t *testing.T) {
	t.Parallel()
	roles := make([]agent.DetailRole, 20)
	for i := range roles {
		roles[i] = agent.DetailRole{
			Name:              "op-" + string(rune('a'+i)),
			SkillLevel:        agent.UnskilledLevel,
			OutputRatePerHour: 60,
			HeadCount:         1,
		}
	}
	m := agent.Manufacture{
		CapitalistID:                agent.NewAgentID(),
		Form:                        agent.SerialManufacture,
		Origin:                      agent.OriginSplitting,
		IndividualWorkingDayMinutes: 720,
		Roles:                       roles,
	}
	if err := m.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
	if len(m.Roles) != 20 {
		t.Fatalf("len(Roles) = %d, want 20", len(m.Roles))
	}
	if m.Origin != agent.OriginSplitting {
		t.Fatalf("Origin = %q, want %q", m.Origin, agent.OriginSplitting)
	}
}

// §2 type-manufacture proportionality: "there are four founders and two
// breakers to one rubber: the founder casts 2,000 type an hour, the breaker
// breaks up 4,000, and the rubber polishes 8,000."
func TestProportionalGroupSize_Para2_TypeManufacture(t *testing.T) {
	t.Parallel()
	roles := []agent.DetailRole{
		{Name: "founder", SkillLevel: agent.SkilledLevel, OutputRatePerHour: 2000, HeadCount: 1},
		{Name: "breaker", SkillLevel: agent.UnskilledLevel, OutputRatePerHour: 4000, HeadCount: 1},
		{Name: "rubber", SkillLevel: agent.UnskilledLevel, OutputRatePerHour: 8000, HeadCount: 1},
	}
	got, err := agent.ProportionalGroupSize(roles, 8000)
	if err != nil {
		t.Fatalf("ProportionalGroupSize() = %v, want nil", err)
	}
	want := map[string]int{"founder": 4, "breaker": 2, "rubber": 1}
	for name, hc := range want {
		if got[name] != hc {
			t.Fatalf("ProportionalGroupSize()[%q] = %d, want %d", name, got[name], hc)
		}
	}
	for _, r := range roles {
		if int64(got[r.Name])*r.OutputRatePerHour != 8000 {
			t.Fatalf("rate × headcount for %q = %d, want 8000", r.Name,
				int64(got[r.Name])*r.OutputRatePerHour)
		}
	}
}

// §2 bottleneck: target output rate that doesn't divide evenly into a
// role's output rate has no integer solution.
func TestProportionalGroupSize_Bottleneck(t *testing.T) {
	t.Parallel()
	roles := []agent.DetailRole{
		{Name: "a", OutputRatePerHour: 3, HeadCount: 1, SkillLevel: agent.SkilledLevel},
	}
	if _, err := agent.ProportionalGroupSize(roles, 100); !errors.Is(err, agent.ErrBottleneck) {
		t.Fatalf("ProportionalGroupSize bottleneck = %v, want ErrBottleneck", err)
	}
}

// §2 serial manufacture — needle wire: 72 detail workmen. Marx gives no
// formula for the gain, so nothing is asserted from the 72 roles alone: with no
// factor recorded the output is exactly the 72 x 720 baseline. A recorded
// observation of 1.3 scales it.
func TestManufactureProductivePower_Para2_72Workmen(t *testing.T) {
	t.Parallel()
	roles := make([]agent.DetailRole, 72)
	for i := range roles {
		roles[i] = agent.DetailRole{
			Name:              "op",
			SkillLevel:        agent.UnskilledLevel,
			OutputRatePerHour: 100,
			HeadCount:         1,
		}
	}
	m := agent.Manufacture{
		CapitalistID:                agent.NewAgentID(),
		Form:                        agent.SerialManufacture,
		Origin:                      agent.OriginSplitting,
		IndividualWorkingDayMinutes: 720,
		Roles:                       roles,
	}
	if len(m.Roles) != 72 {
		t.Fatalf("len(Roles) = %d, want 72", len(m.Roles))
	}
	baseline := agent.CollectiveWorkingDay(agent.CooperationSize(m.TotalWorkers()), m.IndividualWorkingDayMinutes)
	if got := agent.ManufactureProductivePower(m, 720); got != baseline {
		t.Fatalf("no factor recorded: power = %d, want exactly the baseline %d", got, baseline)
	}
	if got := agent.ManufactureProductivePowerFactor(m); got != 1.0 {
		t.Fatalf("no factor recorded: factor = %v, want 1.0", got)
	}
	m.ProductivePower = 1.3
	want := agent.LabourMinutes(51840 * 13 / 10) // 72 x 720 x 1.3
	if got := agent.ManufactureProductivePower(m, 720); got != want {
		t.Fatalf("recorded 1.3: power = %d, want %d", got, want)
	}
}

// A recorded factor below 1 is rejected: the collective labourer never produces
// less than the same workers isolated. 0 (unrecorded) and 1 or more are valid.
func TestManufacture_Validate_ProductivePower(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		factor float64
		want   error
	}{
		{0, nil}, {1, nil}, {1.4, nil},
		{0.9, agent.ErrManufactureProductivePower},
		{-2, agent.ErrManufactureProductivePower},
	} {
		m := carriageManufacture()
		m.ProductivePower = tc.factor
		if err := m.Validate(); !errors.Is(err, tc.want) {
			t.Errorf("factor %v: Validate() = %v, want %v", tc.factor, err, tc.want)
		}
	}
}

// §3 glass furnace group: a "hole" of 5 detail workers is paralysed if any
// one is absent.
func TestCollectiveLabourer_Para3_GlassFurnace(t *testing.T) {
	t.Parallel()
	cl := agent.CollectiveLabourer{Roles: []agent.DetailRole{
		{Name: "bottlemaker", OutputRatePerHour: 10, HeadCount: 1, SkillLevel: agent.SkilledLevel},
		{Name: "blower", OutputRatePerHour: 10, HeadCount: 1, SkillLevel: agent.SkilledLevel},
		{Name: "gatherer", OutputRatePerHour: 10, HeadCount: 1, SkillLevel: agent.SkilledLevel},
		{Name: "putter-up", OutputRatePerHour: 10, HeadCount: 1, SkillLevel: agent.UnskilledLevel},
		{Name: "taker-in", OutputRatePerHour: 10, HeadCount: 1, SkillLevel: agent.UnskilledLevel},
	}}
	if !cl.IsParalysed(1) {
		t.Fatalf("IsParalysed(1) = false, want true")
	}
	if cl.IsParalysed(0) {
		t.Fatalf("IsParalysed(0) = true, want false")
	}
	if cl.OutputPerPeriod(60) <= 0 {
		t.Fatalf("OutputPerPeriod(60) = %d, want > 0", cl.OutputPerPeriod(60))
	}
}

// §3 scaling: once the proportion is set, scale extends only by integer
// multiples of every group. ScaleManufacture(m, 3) triples every role's
// headcount and triples collective output.
func TestScaleManufacture_Para3_IntegerMultiples(t *testing.T) {
	t.Parallel()
	m := carriageManufacture()
	scaled, err := agent.ScaleManufacture(m, 3)
	if err != nil {
		t.Fatalf("ScaleManufacture(m, 3) = %v, want nil", err)
	}
	if scaled.TotalWorkers() != m.TotalWorkers()*3 {
		t.Fatalf("scaled.TotalWorkers() = %d, want %d", scaled.TotalWorkers(), m.TotalWorkers()*3)
	}
	for i, r := range scaled.Roles {
		if r.HeadCount != m.Roles[i].HeadCount*3 {
			t.Fatalf("role %q HeadCount = %d, want %d", r.Name, r.HeadCount, m.Roles[i].HeadCount*3)
		}
	}
	if _, err := agent.ScaleManufacture(m, 0); !errors.Is(err, agent.ErrInvalidMultiplier) {
		t.Fatalf("ScaleManufacture(m, 0) = %v, want ErrInvalidMultiplier", err)
	}
}

// §5 minimum capital grows with productive power; manufacture's minimum
// always exceeds the simple-cooperation minimum of the same headcount at
// the manufacture's own average wage (which is itself lower than the
// purely-skilled wage because manufacture deskills part of the workforce).
func TestManufactureMinimumCapital_Para5_ExceedsCooperation(t *testing.T) {
	t.Parallel()
	m := carriageManufacture()
	var totalWages agent.Pence
	for _, r := range m.Roles {
		totalWages += agent.Pence(int64(r.HeadCount)) * agent.DailyWageForSkill(r.SkillLevel)
	}
	manMin := agent.ManufactureMinimumCapital(m, 0.001)
	if manMin <= totalWages {
		t.Fatalf("ManufactureMinimumCapital = %d, want > totalWages %d (raw-material increment missing)",
			manMin, totalWages)
	}
}

// §5 labour-power value falls for unskilled roles by the amortised
// apprenticeship: the 96,000-min cost is spread over the working life, not
// added to a single day's value. The 9,600-day working life is a test input
// (Marx gives none), chosen so the share is a clean 10 min/day.
func TestRoleLabourPowerValue_Para5_UnskilledIsLower(t *testing.T) {
	t.Parallel()
	skilled := agent.DetailRole{Name: "engraver", SkillLevel: agent.SkilledLevel, OutputRatePerHour: 10, HeadCount: 1}
	unskilled := agent.DetailRole{Name: "stamp-presser", SkillLevel: agent.UnskilledLevel, OutputRatePerHour: 10, HeadCount: 1}
	apprenticeship := agent.LabourMinutes(60 * 8 * 200)
	const workingLifeDays = 9600
	skilledV, err := agent.RoleLabourPowerValue(skilled, apprenticeship, workingLifeDays)
	if err != nil {
		t.Fatalf("skilled: %v", err)
	}
	unskilledV, err := agent.RoleLabourPowerValue(unskilled, apprenticeship, workingLifeDays)
	if err != nil {
		t.Fatalf("unskilled: %v", err)
	}
	if unskilledV != 360 {
		t.Fatalf("unskilledV = %d, want 360 (the 6h / 3s baseline)", unskilledV)
	}
	if skilledV != 370 {
		t.Fatalf("skilledV = %d, want 370 (360 + 96000/9600)", skilledV)
	}
}

// The skilled premium is a daily share, never the whole apprenticeship: a day's
// labour-power must stay far below the 96,000-min training cost.
func TestRoleLabourPowerValue_AmortisedNotLumpSum(t *testing.T) {
	t.Parallel()
	skilled := agent.DetailRole{Name: "engraver", SkillLevel: agent.SkilledLevel, OutputRatePerHour: 10, HeadCount: 1}
	apprenticeship := agent.DefaultApprenticeshipMinutes
	v, err := agent.RoleLabourPowerValue(skilled, apprenticeship, 9600)
	if err != nil {
		t.Fatalf("RoleLabourPowerValue: %v", err)
	}
	if v >= apprenticeship {
		t.Fatalf("daily value %d must be far below the %d-min apprenticeship", v, apprenticeship)
	}
	// Rounds half-up: 5 min over 2 days is 2.5, so 3.
	got, _ := agent.RoleLabourPowerValue(skilled, 5, 2)
	if got != 363 {
		t.Fatalf("360 + round(5/2) = %d, want 363", got)
	}
}

// A skilled role cannot be amortised over a non-positive working life; an
// unskilled role never needs one.
func TestRoleLabourPowerValue_RequiresWorkingLifeForSkilled(t *testing.T) {
	t.Parallel()
	skilled := agent.DetailRole{Name: "engraver", SkillLevel: agent.SkilledLevel}
	unskilled := agent.DetailRole{Name: "stamp-presser", SkillLevel: agent.UnskilledLevel}
	if _, err := agent.RoleLabourPowerValue(skilled, 96000, 0); !errors.Is(err, agent.ErrManufactureWorkingLife) {
		t.Fatalf("skilled, 0 days: err = %v, want ErrManufactureWorkingLife", err)
	}
	if v, err := agent.RoleLabourPowerValue(unskilled, 96000, 0); err != nil || v != 360 {
		t.Fatalf("unskilled, 0 days: (%d, %v), want (360, nil)", v, err)
	}
}

// Marx's labour-power value is 3 s. = 36 pence a day, not the 6 s. = 72 pence a
// 12-hour day produces; both skill levels default to it (Ch. 14 gives no wage).
func TestDailyWageForSkill_IsValueOfLabourPower(t *testing.T) {
	t.Parallel()
	for _, s := range []agent.SkillLevel{agent.SkilledLevel, agent.UnskilledLevel} {
		if got := agent.DailyWageForSkill(s); got != 36 {
			t.Errorf("DailyWageForSkill(%s) = %d, want 36 (3s)", s, got)
		}
	}
}

// Invariant: every role satisfies rate × headcount == targetOutputRate
// after ProportionalGroupSize.
func TestInvariant_ProportionalRateHeadcount(t *testing.T) {
	t.Parallel()
	roles := []agent.DetailRole{
		{Name: "a", OutputRatePerHour: 100, HeadCount: 1, SkillLevel: agent.SkilledLevel},
		{Name: "b", OutputRatePerHour: 50, HeadCount: 1, SkillLevel: agent.UnskilledLevel},
		{Name: "c", OutputRatePerHour: 25, HeadCount: 1, SkillLevel: agent.UnskilledLevel},
	}
	target := int64(100)
	got, err := agent.ProportionalGroupSize(roles, target)
	if err != nil {
		t.Fatalf("ProportionalGroupSize() = %v, want nil", err)
	}
	for _, r := range roles {
		if int64(got[r.Name])*r.OutputRatePerHour != target {
			t.Fatalf("invariant violated for %q: %d × %d != %d",
				r.Name, got[r.Name], r.OutputRatePerHour, target)
		}
	}
}

// Invariant (Ch. 14 §3): scaling extends every detail group by the same
// multiple, so role proportions are preserved. With the factor held fixed,
// productive power is exactly k times the original: scaling invents no gain of
// its own, with or without a recorded factor.
func TestInvariant_ScaleManufactureIsExactlyLinear(t *testing.T) {
	t.Parallel()
	for _, factor := range []float64{0, 1.25} {
		m := carriageManufacture()
		m.ProductivePower = factor
		const period = int64(720)
		orig := agent.ManufactureProductivePower(m, period)
		for _, k := range []int{1, 2, 3, 5} {
			scaled, err := agent.ScaleManufacture(m, k)
			if err != nil {
				t.Fatalf("ScaleManufacture(m, %d) = %v", k, err)
			}
			for i, r := range scaled.Roles {
				if r.HeadCount != m.Roles[i].HeadCount*k {
					t.Fatalf("k=%d role %q headcount = %d, want %d", k, r.Name, r.HeadCount, m.Roles[i].HeadCount*k)
				}
			}
			if got, want := agent.ManufactureProductivePower(scaled, period), agent.LabourMinutes(int64(k))*orig; got != want {
				t.Fatalf("factor %v, k=%d: power = %d, want exactly %d", factor, k, got, want)
			}
		}
	}
}

// Invariant: PartialProduct has no exchange-value; it is a value-magnitude
// alias only. Confirm by type compatibility, not by an ExchangeValue field.
func TestInvariant_PartialProductNotACommodity(t *testing.T) {
	t.Parallel()
	var pp agent.PartialProduct = 100
	if agent.LabourMinutes(pp) != 100 {
		t.Fatalf("PartialProduct conversion = %d, want 100", agent.LabourMinutes(pp))
	}
}

// Invariant: with no factor recorded a manufacture equals the simple
// co-operation baseline; it exceeds it only by a recorded observation >= 1.
func TestInvariant_ManufactureExceedsCooperationOnlyWhenRecorded(t *testing.T) {
	t.Parallel()
	m := carriageManufacture()
	baseline := agent.CollectiveWorkingDay(agent.CooperationSize(m.TotalWorkers()), m.IndividualWorkingDayMinutes)
	if got := agent.ManufactureProductivePower(m, 720); got != baseline {
		t.Fatalf("unrecorded: power = %d, want baseline %d", got, baseline)
	}
	m.ProductivePower = 1.2
	if got := agent.ManufactureProductivePower(m, 720); got <= baseline {
		t.Fatalf("recorded 1.2: power = %d, want > baseline %d", got, baseline)
	}
}

// Invariant: SocialDivisionOfLabour marker distinguishes market-mediated
// from workshop-planned division. (Ch. 14 §4.)
func TestInvariant_SocialDivisionOfLabour(t *testing.T) {
	t.Parallel()
	if agent.SocialDivisionOfLabour != "social" {
		t.Fatalf("SocialDivisionOfLabour = %q, want %q", agent.SocialDivisionOfLabour, "social")
	}
}

func TestManufacture_Validate_Errors(t *testing.T) {
	t.Parallel()
	capID := agent.NewAgentID()
	good := agent.DetailRole{Name: "a", SkillLevel: agent.SkilledLevel, OutputRatePerHour: 10, HeadCount: 1}
	cases := []struct {
		name string
		m    agent.Manufacture
		err  error
	}{
		{"no capitalist", agent.Manufacture{Form: agent.SerialManufacture, Origin: agent.OriginSplitting, IndividualWorkingDayMinutes: 720, Roles: []agent.DetailRole{good}}, agent.ErrManufactureCapitalistID},
		{"bad form", agent.Manufacture{CapitalistID: capID, Form: "weird", Origin: agent.OriginSplitting, IndividualWorkingDayMinutes: 720, Roles: []agent.DetailRole{good}}, agent.ErrManufactureForm},
		{"bad origin", agent.Manufacture{CapitalistID: capID, Form: agent.SerialManufacture, Origin: "weird", IndividualWorkingDayMinutes: 720, Roles: []agent.DetailRole{good}}, agent.ErrManufactureOrigin},
		{"zero working day", agent.Manufacture{CapitalistID: capID, Form: agent.SerialManufacture, Origin: agent.OriginSplitting, IndividualWorkingDayMinutes: 0, Roles: []agent.DetailRole{good}}, agent.ErrManufactureWorkingDay},
		{"no roles", agent.Manufacture{CapitalistID: capID, Form: agent.SerialManufacture, Origin: agent.OriginSplitting, IndividualWorkingDayMinutes: 720, Roles: nil}, agent.ErrManufactureNoRoles},
		{"role no name", agent.Manufacture{CapitalistID: capID, Form: agent.SerialManufacture, Origin: agent.OriginSplitting, IndividualWorkingDayMinutes: 720, Roles: []agent.DetailRole{{SkillLevel: agent.SkilledLevel, OutputRatePerHour: 1, HeadCount: 1}}}, agent.ErrManufactureRoleName},
		{"role bad skill", agent.Manufacture{CapitalistID: capID, Form: agent.SerialManufacture, Origin: agent.OriginSplitting, IndividualWorkingDayMinutes: 720, Roles: []agent.DetailRole{{Name: "x", SkillLevel: "weird", OutputRatePerHour: 1, HeadCount: 1}}}, agent.ErrManufactureSkillLevel},
		{"role zero rate", agent.Manufacture{CapitalistID: capID, Form: agent.SerialManufacture, Origin: agent.OriginSplitting, IndividualWorkingDayMinutes: 720, Roles: []agent.DetailRole{{Name: "x", SkillLevel: agent.SkilledLevel, OutputRatePerHour: 0, HeadCount: 1}}}, agent.ErrManufactureOutputRate},
		{"role zero head", agent.Manufacture{CapitalistID: capID, Form: agent.SerialManufacture, Origin: agent.OriginSplitting, IndividualWorkingDayMinutes: 720, Roles: []agent.DetailRole{{Name: "x", SkillLevel: agent.SkilledLevel, OutputRatePerHour: 1, HeadCount: 0}}}, agent.ErrManufactureHeadCount},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.m.Validate(); !errors.Is(err, tc.err) {
				t.Fatalf("Validate() = %v, want %v", err, tc.err)
			}
		})
	}
}

func TestLabourHierarchy_OrderingAndUnskilled(t *testing.T) {
	t.Parallel()
	m := carriageManufacture()
	h := m.Hierarchy()
	if h.Len() != len(m.Roles) {
		t.Fatalf("Hierarchy.Len() = %d, want %d", h.Len(), len(m.Roles))
	}
	// Skilled-first ordering.
	seenUnskilled := false
	for i := 0; i < h.Len(); i++ {
		r := h.At(i)
		if r.SkillLevel == agent.UnskilledLevel {
			seenUnskilled = true
		} else if seenUnskilled {
			t.Fatalf("skilled role %q after unskilled — hierarchy must be descending", r.Name)
		}
	}
	if len(h.Unskilled()) == 0 {
		t.Fatalf("Unskilled() empty; carriage manufacture has at least one unskilled role")
	}
}
