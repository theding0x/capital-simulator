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
type Factory struct {
	ID              FactoryID       `json:"id"`
	Name            string          `json:"name"`
	PrimeMover      PrimeMover      `json:"prime_mover"`
	Machines        []Machine       `json:"machines"`
	TickCount       int64           `json:"tick_count"`
	IntensityFactor IntensityFactor `json:"intensity_factor"`
	CreatedAt       time.Time       `json:"created_at"`
}

var (
	ErrFactoryNoMachines = errors.New("factory: at least one machine is required")
	ErrFactoryNoMover    = errors.New("factory: prime_mover kind is required")
)

func (f Factory) Validate() error {
	if f.PrimeMover.Kind == "" {
		return ErrFactoryNoMover
	}
	if len(f.Machines) == 0 {
		return ErrFactoryNoMachines
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
// per working day across every machine in the factory.
func (f Factory) DailyValueTransfer() LabourMinutes {
	var total LabourMinutes
	for _, m := range f.Machines {
		total += DailyWearAndTear(m)
	}
	return total
}

// TickResult is one period's worth of value transfer and output for a Factory.
type TickResult struct {
	ValueTransferred LabourMinutes `json:"value_transferred"`
	UnitsProduced    int64         `json:"units_produced"`
	HandLabourSaved  LabourMinutes `json:"hand_labour_saved"`
}

// RunTick computes one period's contribution without mutating the receiver.
// The store layer is responsible for persisting the accumulation. The
// factory's IntensityFactor (§3C) scales the saved hand-labour: a more
// densely-packed working day exposes more hand-labour displacement per
// calendar hour, even though MachineValue transfer remains tied to wear.
func (f Factory) RunTick() TickResult {
	var saved LabourMinutes
	var units int64
	intensity := f.IntensityFactor
	if intensity < 1.0 {
		intensity = 1.0
	}
	for _, m := range f.Machines {
		units += int64(m.ProductivePower)
		if m.HandLabourPerUnit > 0 {
			displaced := LabourDisplaced(m.HandLabourPerUnit, int64(m.ProductivePower))
			saved += EffectiveLabour(displaced, intensity)
		}
	}
	return TickResult{
		ValueTransferred: f.DailyValueTransfer(),
		UnitsProduced:    units,
		HandLabourSaved:  saved,
	}
}

// FactoryProductivityFactor returns the factor by which a machinery-driven
// factory's labour-time productivity exceeds an equivalent hand-labour
// process. It is the §2 ratio:
//
//	(hand labour required for the day's output) / (machine labour required)
//
// where the "machine labour required" approximates the value the machines
// transfer to product (DailyValueTransfer). When the factory carries no
// hand-labour baseline the factor falls back to 1.0. This number is the
// Ch. 15 contribution to the same productivity-factor vocabulary used by
// Ch. 13 and Ch. 14 — Marx's three forms of relative-SV mechanism.
func FactoryProductivityFactor(f Factory) float64 {
	tr := f.RunTick()
	if tr.HandLabourSaved <= 0 || tr.ValueTransferred <= 0 {
		return 1.0
	}
	hand := float64(tr.HandLabourSaved) + float64(tr.ValueTransferred)
	machine := float64(tr.ValueTransferred)
	return hand / machine
}
