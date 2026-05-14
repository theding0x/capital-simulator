// Package machinery models Capital Vol. I, Ch. 15 — Machinery and Modern
// Industry. A Machine is Marx's three-part instrument (motor mechanism,
// transmitting mechanism, working tool); a Factory is an aggregate of
// machines driven by one prime mover. The package exposes the pure
// value-transfer functions (wear and tear, moral depreciation, labour
// displaced) and the capital-composition pair that mechanisation reshapes.
//
// All functions are deterministic; persistence lives in internal/store.
package machinery

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/labour"
)

// LabourMinutes is the canonical value-magnitude unit. The canonical
// definition lives in the labour package; this alias keeps existing
// machinery call-sites compiling while letting Ch. 15 values be passed
// to Ch. 12 production / Ch. 11 surplus functions without explicit casts.
type LabourMinutes = labour.LabourMinutes

// MachineValue is the total labour crystallised in the machine at the moment
// it leaves its place of production. Expressed in LabourMinutes.
type MachineValue LabourMinutes

// LifespanDays is the number of standard working days over which a machine
// transfers its value to product. A "standard" working day is
// StandardWorkingDayHours hours long.
type LifespanDays int64

// StandardWorkingDayHours is the calibration constant for LifespanDays. Marx's
// numerical examples in §2 implicitly assume an 8-hour reference day.
const StandardWorkingDayHours int64 = 8

// ProductivePower is the use-value output a machine produces per working day,
// measured in discrete units of product (needles, yards of yarn, etc.).
type ProductivePower int64

// PrimeMoverKind enumerates the source of motive power that drives a Factory.
type PrimeMoverKind string

const (
	PrimeMoverSteam    PrimeMoverKind = "steam"
	PrimeMoverWater    PrimeMoverKind = "water"
	PrimeMoverElectric PrimeMoverKind = "electric"
	PrimeMoverAnimal   PrimeMoverKind = "animal"
)

// PrimeMover is the source of motive power. Each Factory shares one prime
// mover across all its constituent machines.
type PrimeMover struct {
	Kind       PrimeMoverKind `json:"kind"`
	Horsepower int64          `json:"horsepower"`
}

// MachineID is a 96-bit hex identifier for stored Machine records.
type MachineID string

func (id MachineID) IsZero() bool { return id == "" }

// NewMachineID returns a fresh 96-bit hex identifier. Mirrors commodity.NewID.
func NewMachineID() MachineID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return MachineID(hex.EncodeToString(b))
}

// Machine is Marx's three-part instrument of production (Ch. 15, §1).
// The motor mechanism delivers force, the transmitting mechanism distributes
// it, the working tool acts upon the material. A Machine carries a finite
// MachineValue, a LifespanDays, and a ProductivePower — together these
// determine the rate at which the machine surrenders its own value to
// product, and the daily volume of use-values produced.
type Machine struct {
	ID                     MachineID       `json:"id"`
	Name                   string          `json:"name"`
	MotorMechanism         string          `json:"motor_mechanism"`
	TransmittingMechanism  string          `json:"transmitting_mechanism"`
	WorkingTool            string          `json:"working_tool"`
	MachineValue           MachineValue    `json:"machine_value"`
	LifespanDays           LifespanDays    `json:"lifespan_days"`
	ProductivePower        ProductivePower `json:"productive_power"`
	HandLabourPerUnit      LabourMinutes   `json:"hand_labour_per_unit"`
	AccumulatedWear        MaterialWear    `json:"accumulated_wear"`
	AccumulatedDepreciation MoralDepreciation `json:"accumulated_depreciation"`
	CreatedAt              time.Time       `json:"created_at"`
}

// MaterialWear is the labour-value lost to physical degradation in use.
// Accumulates as the machine is operated; bounded above by MachineValue.
type MaterialWear struct {
	Value LabourMinutes `json:"value"`
}

// MoralDepreciation is the labour-value lost when cheaper or better machines
// enter the market — independent of physical wear. "It loses exchange-value,
// either by machines of the same sort being produced cheaper than it, or by
// better machines entering into competition with it." (§3B)
type MoralDepreciation struct {
	Value LabourMinutes `json:"value"`
}

// Domain errors.
var (
	ErrMachineValue    = errors.New("machine: machine_value must be positive")
	ErrLifespanDays    = errors.New("machine: lifespan_days must be positive")
	ErrProductivePower = errors.New("machine: productive_power must be positive")
	ErrPrimeMoverKind  = errors.New("machine: prime_mover kind is required")
)

// Validate enforces the basic positivity constraints required for any of the
// downstream value-transfer formulas to be meaningful.
func (m Machine) Validate() error {
	if m.MachineValue <= 0 {
		return ErrMachineValue
	}
	if m.LifespanDays <= 0 {
		return ErrLifespanDays
	}
	if m.ProductivePower <= 0 {
		return ErrProductivePower
	}
	return nil
}

// DailyWearAndTear returns the LabourMinutes a machine transfers to product
// each working day. From §2: "the machine ... enters into the value-begetting
// process only by bits. It never adds more value than it loses, on an
// average, by wear and tear."
//
//	dwt = MachineValue / LifespanDays
func DailyWearAndTear(m Machine) LabourMinutes {
	if m.LifespanDays <= 0 {
		return 0
	}
	return LabourMinutes(int64(m.MachineValue) / int64(m.LifespanDays))
}

// ValueTransferredPerUnit returns the LabourMinutes per unit of product
// contributed by machine wear. Doubling output halves per-unit value (§2).
//
//	vtu = DailyWearAndTear(m) / unitsPerDay
func ValueTransferredPerUnit(m Machine, unitsPerDay int64) LabourMinutes {
	if unitsPerDay <= 0 {
		return 0
	}
	return DailyWearAndTear(m) / LabourMinutes(unitsPerDay)
}

// TotalValueTransferred returns the cumulative LabourMinutes transferred by a
// machine that has run for `days` calendar days at `hoursPerDay` hours per
// day, where the machine's LifespanDays is calibrated against an 8-hour
// reference day. The total is capped at MachineValue: a machine never adds
// more value over its life than the value with which it began.
//
// §2 lifecycle equivalence: 16 h/day for 7½ years and 8 h/day for 15 years
// both exhaust the full MachineValue, yet the 16-hour case reproduces it
// twice as quickly per calendar year (see RateOfValueReproduction).
func TotalValueTransferred(m Machine, hoursPerDay, days int64) LabourMinutes {
	if hoursPerDay <= 0 || days <= 0 || m.LifespanDays <= 0 {
		return 0
	}
	lifeHours := int64(m.LifespanDays) * StandardWorkingDayHours
	worked := hoursPerDay * days
	if worked >= lifeHours {
		return LabourMinutes(m.MachineValue)
	}
	return LabourMinutes(int64(m.MachineValue) * worked / lifeHours)
}

// RateOfValueReproduction returns the LabourMinutes per calendar day at which
// a machine surrenders its value when operated `hoursPerDay` hours per day.
// Doubling the hours doubles the rate; total value transferred over the full
// life is unchanged. §2: "in the first case the value would be reproduced
// twice as quickly."
func RateOfValueReproduction(m Machine, hoursPerDay int64) float64 {
	if m.LifespanDays <= 0 {
		return 0
	}
	return float64(m.MachineValue) * float64(hoursPerDay) /
		(float64(m.LifespanDays) * float64(StandardWorkingDayHours))
}

// LabourSavingRatio returns how many times more hand labour the machine path
// replaces. From §2: "366 lbs. of cotton absorb only 150 hours' labour" by
// machine vs. 27,000 hours by hand — a 180× saving.
func LabourSavingRatio(handMinutes, machineMinutes LabourMinutes) int64 {
	if machineMinutes <= 0 {
		return 0
	}
	return int64(handMinutes) / int64(machineMinutes)
}

// LabourDisplaced returns the LabourMinutes of hand labour the machine
// replaces per period for `unitsProduced` units of product. The machine's
// productive advantage is realised only when this exceeds the labour
// embodied in (and amortised across) the machine itself (§2 invariant).
func LabourDisplaced(handLabourPerUnit LabourMinutes, unitsProduced int64) LabourMinutes {
	if unitsProduced <= 0 || handLabourPerUnit <= 0 {
		return 0
	}
	return handLabourPerUnit * LabourMinutes(unitsProduced)
}

// ComputeMoralDepreciation returns the value an existing machine loses when
// a rival with the same ProductivePower is produced more cheaply (§3B). The
// result is clamped at zero: competition from cheaper successors never
// raises a machine's value.
func ComputeMoralDepreciation(existing, rival Machine) MoralDepreciation {
	if rival.MachineValue >= existing.MachineValue {
		return MoralDepreciation{Value: 0}
	}
	return MoralDepreciation{
		Value: LabourMinutes(existing.MachineValue - rival.MachineValue),
	}
}
