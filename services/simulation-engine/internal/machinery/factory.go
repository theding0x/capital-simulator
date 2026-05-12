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
// All machines in a Factory share a single PrimeMover.
type Factory struct {
	ID         FactoryID  `json:"id"`
	Name       string     `json:"name"`
	PrimeMover PrimeMover `json:"prime_mover"`
	Machines   []Machine  `json:"machines"`
	TickCount  int64      `json:"tick_count"`
	CreatedAt  time.Time  `json:"created_at"`
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
// The store layer is responsible for persisting the accumulation.
func (f Factory) RunTick() TickResult {
	var saved LabourMinutes
	var units int64
	for _, m := range f.Machines {
		units += int64(m.ProductivePower)
		if m.HandLabourPerUnit > 0 {
			saved += LabourDisplaced(m.HandLabourPerUnit, int64(m.ProductivePower))
		}
	}
	return TickResult{
		ValueTransferred: f.DailyValueTransfer(),
		UnitsProduced:    units,
		HandLabourSaved:  saved,
	}
}
