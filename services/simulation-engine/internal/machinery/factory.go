package machinery

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

// FactoryID is a 96-bit hex identifier for stored Factory records.
type FactoryID string

func (id FactoryID) IsZero() bool { return id == "" }

func NewFactoryID() FactoryID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return FactoryID(hex.EncodeToString(b))
}

// Factory is "an organised system of machines, to which motion is
// communicated by the transmitting mechanism from a central automaton" (§1c).
// All machines in a Factory share a single PrimeMover. IntensityFactor
// captures the §3C labour-density multiplier: when the working day is
// shortened under machinery, the remaining hours fill more densely. The
// factor is bounded below at 1.0 by NewIntensityFactor; default zero
// means "unset" and is treated as 1.0 by the tick path.
//
// Attendants are the factory's living labour (§2: a machine's productivity
// is measured by the labour it costs — the value it transfers plus the
// living labour of those who tend it). Only living labour creates new value
// (W = c + v + s). Each attendant works LivingLabourPerAttendant per tick, of
// which NecessaryLabourPerAttendant reproduces their labour-power (v) and the
// remainder is surplus-labour (s). §3C intensity densifies that living labour.
type Factory struct {
	ID              FactoryID       `json:"id"`
	Name            string          `json:"name"`
	PrimeMover      PrimeMover      `json:"prime_mover"`
	Machines        []Machine       `json:"machines"`
	TickCount       int64           `json:"tick_count"`
	IntensityFactor IntensityFactor `json:"intensity_factor"`

	Attendants                  int64         `json:"attendants"`
	LivingLabourPerAttendant    LabourMinutes `json:"living_labour_per_attendant"`
	NecessaryLabourPerAttendant LabourMinutes `json:"necessary_labour_per_attendant"`

	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrFactoryNoMachines = errors.New("factory: at least one machine is required")
	ErrFactoryNoMover    = errors.New("factory: prime_mover kind is required")
	ErrFactoryLabour     = errors.New("factory: attendant labour must be non-negative and necessary labour cannot exceed living labour")
)

func (f Factory) Validate() error {
	if f.PrimeMover.Kind == "" {
		return ErrFactoryNoMover
	}
	if len(f.Machines) == 0 {
		return ErrFactoryNoMachines
	}
	if f.Attendants < 0 || f.LivingLabourPerAttendant < 0 || f.NecessaryLabourPerAttendant < 0 ||
		f.NecessaryLabourPerAttendant > f.LivingLabourPerAttendant {
		return ErrFactoryLabour
	}
	for _, m := range f.Machines {
		if err := m.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// TotalProductivePower returns the sum of ProductivePower across all machines.
// The Factory's daily use-value output scales linearly with the number of
// machines a single prime mover can drive (§4).
func (f Factory) TotalProductivePower() ProductivePower {
	var total int64
	for _, m := range f.Machines {
		total += int64(m.ProductivePower)
	}
	return ProductivePower(total)
}

// DailyValueTransfer is the aggregate LabourMinutes transferred to product
// per working day across every machine in the factory, each capped at its
// remaining MachineValue (see NextWear).
func (f Factory) DailyValueTransfer() LabourMinutes {
	var total LabourMinutes
	for _, m := range f.Machines {
		total += NextWear(m)
	}
	return total
}

// LivingLabour is the new value the attendants create per tick: their nominal
// labour densified by the §3C intensity factor (an unset factor counts as
// 1.0). Machines transfer value but create none; this is the only source of
// v + s.
func (f Factory) LivingLabour() LabourMinutes {
	if f.Attendants <= 0 || f.LivingLabourPerAttendant <= 0 {
		return 0
	}
	return EffectiveLabour(LabourMinutes(f.Attendants)*f.LivingLabourPerAttendant, f.IntensityFactor)
}

// VariableCapital is the value of the attendants' labour-power per tick (v),
// bounded above by the living labour they actually perform.
func (f Factory) VariableCapital() LabourMinutes {
	if f.Attendants <= 0 || f.NecessaryLabourPerAttendant <= 0 {
		return 0
	}
	v := LabourMinutes(f.Attendants) * f.NecessaryLabourPerAttendant
	return min(v, f.LivingLabour())
}

// SurplusValue is the unpaid part of the living labour per tick (s).
func (f Factory) SurplusValue() LabourMinutes {
	return f.LivingLabour() - f.VariableCapital()
}

// TickResult is one period's worth of value transfer and output for a Factory.
// ValueTransferred (c) + VariableCapital (v) + SurplusValue (s) is the value
// of the period's product; LivingLabour = VariableCapital + SurplusValue.
type TickResult struct {
	ValueTransferred LabourMinutes `json:"value_transferred"`
	UnitsProduced    int64         `json:"units_produced"`
	HandLabourSaved  LabourMinutes `json:"hand_labour_saved"`
	LivingLabour     LabourMinutes `json:"living_labour"`
	VariableCapital  LabourMinutes `json:"variable_capital"`
	SurplusValue     LabourMinutes `json:"surplus_value"`
}

// RunTick computes one period's contribution without mutating the receiver.
// The store layer is responsible for persisting the accumulation.
// HandLabourSaved is the hand labour the machines displace (§2). Machines
// transfer value by wear (c); only the attendants' living labour, densified
// by the factory's IntensityFactor (§3C), creates new value (v + s).
func (f Factory) RunTick() TickResult {
	var saved LabourMinutes
	var units int64
	for _, m := range f.Machines {
		units += int64(m.ProductivePower)
		saved += LabourDisplaced(m.HandLabourPerUnit, int64(m.ProductivePower))
	}
	return TickResult{
		ValueTransferred: f.DailyValueTransfer(),
		UnitsProduced:    units,
		HandLabourSaved:  saved,
		LivingLabour:     f.LivingLabour(),
		VariableCapital:  f.VariableCapital(),
		SurplusValue:     f.SurplusValue(),
	}
}

// FactoryProductivityFactor returns the factor by which a machinery-driven
// factory's labour-time productivity exceeds an equivalent hand-labour
// process. It is the §2 ratio:
//
//	(hand labour displaced) / (labour the machinery costs)
//
// where the labour the machinery costs is the value it transfers to product
// (DailyValueTransfer) plus the living labour of its attendants. When the
// factory displaces no hand labour, or the machinery costs as much as it
// replaces (outside the §2 social limit), the factor is 1.0. This number is
// the Ch. 15 contribution to the same productivity-factor vocabulary used by
// Ch. 13 and Ch. 14 — Marx's three forms of relative-SV mechanism.
func FactoryProductivityFactor(f Factory) float64 {
	tr := f.RunTick()
	cost := tr.ValueTransferred + tr.LivingLabour
	if tr.HandLabourSaved <= 0 || cost <= 0 {
		return 1.0
	}
	return max(1.0, float64(tr.HandLabourSaved)/float64(cost))
}
