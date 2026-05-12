package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

// ManufactureID is a 96-bit hex identifier for stored Manufacture records.
type ManufactureID string

func (id ManufactureID) IsZero() bool { return id == "" }

func NewManufactureID() ManufactureID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return ManufactureID(hex.EncodeToString(b))
}

// ManufactureForm enumerates Marx's two organisational shapes for manufacture
// (Capital Vol. I, Ch. 14 §2): heterogeneous (assembly of independently
// produced partial products) and serial (a single product passing through
// connected sequential stages).
type ManufactureForm string

const (
	HeterogeneousManufacture ManufactureForm = "heterogeneous"
	SerialManufacture        ManufactureForm = "serial"
)

func (f ManufactureForm) Valid() bool {
	switch f {
	case HeterogeneousManufacture, SerialManufacture:
		return true
	}
	return false
}

// ManufactureOrigin captures the two-fold historical origin Marx identifies
// in §1: combination (diverse handicrafts united under one capitalist) and
// splitting (one handicraft subdivided into detail operations).
type ManufactureOrigin string

const (
	OriginCombination ManufactureOrigin = "combination"
	OriginSplitting   ManufactureOrigin = "splitting"
)

func (o ManufactureOrigin) Valid() bool {
	switch o {
	case OriginCombination, OriginSplitting:
		return true
	}
	return false
}

// SkillLevel is the binary skilled / unskilled distinction manufacture
// produces in the labour force (Ch. 14 §5).
type SkillLevel string

const (
	SkilledLevel   SkillLevel = "skilled"
	UnskilledLevel SkillLevel = "unskilled"
)

func (s SkillLevel) Valid() bool {
	switch s {
	case SkilledLevel, UnskilledLevel:
		return true
	}
	return false
}

// SocialDivisionOfLabour is a marker constant used to distinguish in
// commentary and tests the market-mediated social division of labour
// (anarchic, a posteriori) from the workshop's planned, a priori
// division of labour. Ch. 14 §4: the two are governed by opposite laws
// and must never be conflated in domain logic.
const SocialDivisionOfLabour = "social"

// Default daily wages used by the minimum-capital probe. Marx's standard
// reference value is 6 s. = 72 pence; unskilled labour-power is reproduced
// without the apprenticeship cost and so commands a lower wage.
const (
	SkilledDailyWagePence   Pence = 72
	UnskilledDailyWagePence Pence = 48
)

// DefaultApprenticeshipMinutes is the SNLT embedded in skilled labour-power
// by Marx's "training, education, apprenticeship". Used as the default
// gap between skilled and unskilled `RoleLabourPowerValue` when callers do
// not supply one explicitly.
const DefaultApprenticeshipMinutes LabourMinutes = 60 * 8 * 200 // 200 working days at 8h

// PartialProduct is the fractional output of one DetailRole per period.
// It is deliberately a value-magnitude alias and carries no exchange-value:
// only the finished output of the CollectiveLabourer enters circulation as
// a commodity. (Ch. 14 §4, §1: "the products of his partial labours are
// not commodities" — they exist only inside the workshop.)
type PartialProduct LabourMinutes

// DetailRole is a named fractional function within a Manufacture. Each
// worker is "annexed" (§1) to exactly one role. OutputRatePerHour is the
// role's nominal throughput; HeadCount is the number of workers performing
// it; ToolName names the specialised instrument (Ch. 14 §1, Birmingham
// 500-hammer footnote) that the role uses.
type DetailRole struct {
	Name              string     `json:"name"`
	SkillLevel        SkillLevel `json:"skill_level"`
	OutputRatePerHour int64      `json:"output_rate_per_hour"`
	HeadCount         int        `json:"head_count"`
	ToolName          string     `json:"tool_name"`
}

func (r DetailRole) Validate() error {
	if r.Name == "" {
		return ErrManufactureRoleName
	}
	if !r.SkillLevel.Valid() {
		return ErrManufactureSkillLevel
	}
	if r.OutputRatePerHour <= 0 {
		return ErrManufactureOutputRate
	}
	if r.HeadCount <= 0 {
		return ErrManufactureHeadCount
	}
	return nil
}

// SpecialisedTool is the tool adapted to a single DetailRole. Tool
// differentiation is the material precondition manufacture bequeaths to
// machinery (Ch. 15); here it lives as the role's nominated instrument.
type SpecialisedTool struct {
	Name    string `json:"name"`
	RoleRef string `json:"role_ref"`
}

// LabourHierarchy is the ordered slice of DetailRole sorted by skill level
// (Ch. 14 §5: "manufacture develops a hierarchy of labour-powers"). The
// ordering descends from skilled to unskilled so that callers can read off
// the wage gradient directly.
type LabourHierarchy []DetailRole

func (h LabourHierarchy) Len() int           { return len(h) }
func (h LabourHierarchy) At(i int) DetailRole { return h[i] }

// Unskilled returns the unskilled tail of the hierarchy — the class
// manufacture begets "for whom the cost of apprenticeship vanishes" (§5).
func (h LabourHierarchy) Unskilled() []DetailRole {
	out := make([]DetailRole, 0, len(h))
	for _, r := range h {
		if r.SkillLevel == UnskilledLevel {
			out = append(out, r)
		}
	}
	return out
}

// CollectiveLabourer is the organic whole formed by the DetailRole workers
// of a Manufacture (Ch. 14 §4). Its productive power exceeds the arithmetic
// sum of isolated workers; if any single role is unstaffed the whole is
// "paralysed" (§3, glass-furnace fixture).
type CollectiveLabourer struct {
	Roles []DetailRole `json:"roles"`
}

func (cl CollectiveLabourer) TotalWorkers() int {
	n := 0
	for _, r := range cl.Roles {
		n += r.HeadCount
	}
	return n
}

// IsParalysed reports whether the absence of `absentRoles` roles halts
// production. In a serial manufacture (§3) any missing detail role breaks
// the chain; we model the threshold as "one or more absent roles".
func (cl CollectiveLabourer) IsParalysed(absentRoles int) bool {
	return absentRoles > 0
}

// OutputPerPeriod returns the bottleneck-bounded use-value output of the
// collective labourer over periodMinutes. This is the §2 proportionality
// law in action: the slowest group caps the whole.
func (cl CollectiveLabourer) OutputPerPeriod(periodMinutes int64) LabourMinutes {
	if len(cl.Roles) == 0 || periodMinutes <= 0 {
		return 0
	}
	var minCap int64 = -1
	for _, r := range cl.Roles {
		cap := r.OutputRatePerHour * int64(r.HeadCount) * periodMinutes / 60
		if minCap == -1 || cap < minCap {
			minCap = cap
		}
	}
	if minCap < 0 {
		return 0
	}
	return LabourMinutes(minCap)
}

// Manufacture is a Cooperation (Ch. 13) reorganised around a fixed division
// of detail labour. It is owned by a Capitalist and holds a set of
// DetailRoles. (Ch. 14 §1: "Manufacture is that co-operation in which the
// division of labour is carried through in the workshop.")
type Manufacture struct {
	ID                          ManufactureID     `json:"id"`
	Name                        string            `json:"name"`
	CapitalistID                AgentID           `json:"capitalist_id"`
	Form                        ManufactureForm   `json:"form"`
	Origin                      ManufactureOrigin `json:"origin"`
	IndividualWorkingDayMinutes LabourMinutes     `json:"individual_working_day_minutes"`
	Roles                       []DetailRole      `json:"roles"`
	CreatedAt                   time.Time         `json:"created_at"`
}

func (m Manufacture) Validate() error {
	if m.CapitalistID.IsZero() {
		return ErrManufactureCapitalistID
	}
	if !m.Form.Valid() {
		return ErrManufactureForm
	}
	if !m.Origin.Valid() {
		return ErrManufactureOrigin
	}
	if m.IndividualWorkingDayMinutes <= 0 {
		return ErrManufactureWorkingDay
	}
	if len(m.Roles) == 0 {
		return ErrManufactureNoRoles
	}
	for _, r := range m.Roles {
		if err := r.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (m Manufacture) TotalWorkers() int {
	return m.CollectiveLabourer().TotalWorkers()
}

func (m Manufacture) CollectiveLabourer() CollectiveLabourer {
	roles := make([]DetailRole, len(m.Roles))
	copy(roles, m.Roles)
	return CollectiveLabourer{Roles: roles}
}

// Hierarchy returns the role slice ordered by skill level (skilled before
// unskilled), preserving original order within each level.
func (m Manufacture) Hierarchy() LabourHierarchy {
	skilled := make([]DetailRole, 0, len(m.Roles))
	unskilled := make([]DetailRole, 0, len(m.Roles))
	for _, r := range m.Roles {
		if r.SkillLevel == UnskilledLevel {
			unskilled = append(unskilled, r)
		} else {
			skilled = append(skilled, r)
		}
	}
	return LabourHierarchy(append(skilled, unskilled...))
}

// Errors --------------------------------------------------------------------

var (
	ErrManufactureCapitalistID = errors.New("manufacture: capitalist_id is required")
	ErrManufactureForm         = errors.New("manufacture: form must be 'heterogeneous' or 'serial'")
	ErrManufactureOrigin       = errors.New("manufacture: origin must be 'combination' or 'splitting'")
	ErrManufactureWorkingDay   = errors.New("manufacture: individual_working_day_minutes must be positive")
	ErrManufactureNoRoles      = errors.New("manufacture: at least one detail role is required")
	ErrManufactureRoleName     = errors.New("manufacture: role name is required")
	ErrManufactureSkillLevel   = errors.New("manufacture: role skill_level must be 'skilled' or 'unskilled'")
	ErrManufactureOutputRate   = errors.New("manufacture: role output_rate_per_hour must be positive")
	ErrManufactureHeadCount    = errors.New("manufacture: role head_count must be positive")
	ErrBottleneck              = errors.New("manufacture: no integer proportional headcount satisfies the target output rate")
	ErrInvalidMultiplier       = errors.New("manufacture: multiplier must be >= 1")
)

// Pure functions ------------------------------------------------------------

// ProportionalGroupSize returns the per-role headcount required so that the
// product of each role's output rate × headcount equals targetOutputRate.
// This is Marx's "iron law" of manufacture (Ch. 14 §2): "the different
// processes carried on in a workshop must be performed in definite
// proportions to one another". Returns ErrBottleneck if no integer
// solution exists for any role.
func ProportionalGroupSize(roles []DetailRole, targetOutputRate int64) (map[string]int, error) {
	if targetOutputRate <= 0 {
		return nil, ErrManufactureOutputRate
	}
	out := make(map[string]int, len(roles))
	for _, r := range roles {
		if r.OutputRatePerHour <= 0 {
			return nil, ErrManufactureOutputRate
		}
		if targetOutputRate%r.OutputRatePerHour != 0 {
			return nil, ErrBottleneck
		}
		out[r.Name] = int(targetOutputRate / r.OutputRatePerHour)
	}
	return out, nil
}

// ManufactureProductivePowerFactor returns the multiplicative factor by
// which the manufacture's collective output of use-values exceeds the
// arithmetic sum of n isolated workers. It compounds the Ch. 13 cooperation
// bonus with a division-of-labour increment that grows with the number of
// detail roles. Bounded by construction; not normative.
func ManufactureProductivePowerFactor(m Manufacture) float64 {
	n := m.TotalWorkers()
	roleCount := len(m.Roles)
	if n <= 0 {
		return 1.0
	}
	coop := CollectiveProductivePower(CooperationSize(n), m.IndividualWorkingDayMinutes)
	if roleCount <= 1 {
		return coop
	}
	extra := 0.05 * float64(roleCount-1)
	if extra > 0.75 {
		extra = 0.75
	}
	return coop + extra
}

// ManufactureProductivePower returns the collective output of a Manufacture
// in LabourMinutes-of-use-value over periodMinutes. For len(Roles) > 1 it
// is strictly greater than the simple-cooperation baseline of n × period.
// (Ch. 14 §2, §3: "the collective labourer ... produces far more than the
// sum of his isolated parts".)
func ManufactureProductivePower(m Manufacture, periodMinutes int64) LabourMinutes {
	n := m.TotalWorkers()
	if n <= 0 || periodMinutes <= 0 {
		return 0
	}
	baseline := int64(n) * periodMinutes
	factor := ManufactureProductivePowerFactor(m)
	out := int64(float64(baseline) * factor)
	if len(m.Roles) > 1 && out <= baseline {
		out = baseline + 1
	}
	return LabourMinutes(out)
}

// DailyWageForSkill returns the default daily wage for a skill level. The
// gap between skilled and unskilled is Marx's apprenticeship-cost wedge
// (Ch. 14 §5) expressed in pence.
func DailyWageForSkill(s SkillLevel) Pence {
	if s == UnskilledLevel {
		return UnskilledDailyWagePence
	}
	return SkilledDailyWagePence
}

// RoleLabourPowerValue returns the SNLT in LabourMinutes required to
// reproduce a worker assigned to the given DetailRole. Skilled roles carry
// the additional apprenticeship cost; unskilled roles do not. (Ch. 14 §5:
// "the cost of apprenticeship vanishes for the unskilled labourers.")
func RoleLabourPowerValue(role DetailRole, apprenticeshipMinutes LabourMinutes) LabourMinutes {
	const subsistence LabourMinutes = 60 * 6 // 6 hours subsistence baseline
	if apprenticeshipMinutes < 0 {
		apprenticeshipMinutes = 0
	}
	if role.SkillLevel == UnskilledLevel {
		return subsistence
	}
	return subsistence + apprenticeshipMinutes
}

// ManufactureMinimumCapital returns the minimum Pence the capitalist must
// command before the manufacture can begin (Ch. 14 §5). Manufacture's
// minimum is always strictly greater than the simple-cooperation minimum
// of the same headcount: it must cover wages in role proportions plus raw
// material consumed at the faster collective rate.
//
// rawMaterialCostFactor is pence per LabourMinute of collective output;
// 0 disables the raw-material component but still produces a sum
// >= MinimumCapital(totalWorkers, averageDailyWage).
func ManufactureMinimumCapital(m Manufacture, rawMaterialCostFactor float64) Pence {
	if rawMaterialCostFactor < 0 {
		rawMaterialCostFactor = 0
	}
	var wages Pence
	for _, r := range m.Roles {
		wages += Pence(int64(r.HeadCount)) * DailyWageForSkill(r.SkillLevel)
	}
	period := int64(m.IndividualWorkingDayMinutes)
	if period <= 0 {
		period = 480
	}
	output := ManufactureProductivePower(m, period)
	raw := Pence(int64(float64(output) * rawMaterialCostFactor))
	return wages + raw
}

// ScaleManufacture returns a new Manufacture whose every DetailRole has
// its HeadCount multiplied by the given integer. (Ch. 14 §3: once the
// proportion is set, "that scale can be extended only by employing a
// multiple of each particular group".) Non-integer or sub-unit multipliers
// are invalid.
func ScaleManufacture(m Manufacture, multiplier int) (Manufacture, error) {
	if multiplier < 1 {
		return Manufacture{}, ErrInvalidMultiplier
	}
	out := m
	out.Roles = make([]DetailRole, len(m.Roles))
	for i, r := range m.Roles {
		r.HeadCount *= multiplier
		out.Roles[i] = r
	}
	return out, nil
}
