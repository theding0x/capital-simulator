package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

var ErrInvalidProcess = errors.New("agent: invalid process")

// LabourProcessID is a 96-bit hex ID for LabourProcess records.
type LabourProcessID string

func (id LabourProcessID) IsZero() bool { return id == "" }

func NewLabourProcessID() LabourProcessID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return LabourProcessID(hex.EncodeToString(b))
}

// RawMaterial is an input commodity consumed in one production run [§1.c].
type RawMaterial struct {
	CommodityID string        `json:"commodity_id"`
	Quantity    int64         `json:"quantity"`
	SNLTPerUnit LabourMinutes `json:"snlt_per_unit"`
}

// Instrument is a tool or machine that transfers value to the product gradually [§1.c].
type Instrument struct {
	CommodityID string        `json:"commodity_id"`
	WearPerRun  LabourMinutes `json:"wear_per_run"`
}

// MeansOfProduction bundles raw materials and instruments consumed in one run [§1.c].
type MeansOfProduction struct {
	RawMaterials []RawMaterial `json:"raw_materials"`
	Instruments  []Instrument  `json:"instruments"`
}

// LabourProcess is one purposeful act of production [§1].
// NecessaryLabourMinutes, ProductKind, and ProductQuantity are snapshotted at
// run-time so GET /v1/labour-processes/{id} can reconstruct the full result.
type LabourProcess struct {
	ID                     LabourProcessID   `json:"id"`
	WorkerID               AgentID           `json:"worker_id"`
	CapitalistID           AgentID           `json:"capitalist_id"`
	Means                  MeansOfProduction `json:"means"`
	Duration               LabourMinutes     `json:"duration"`
	NecessaryLabourMinutes LabourMinutes     `json:"necessary_labour_minutes"`
	ProductKind            string            `json:"product_kind"`
	ProductQuantity        int64             `json:"product_quantity"`
	CreatedAt              time.Time         `json:"created_at"`
}

// Validate rejects zero-duration runs and nil required parties [§1.d].
func (lp LabourProcess) Validate() error {
	if lp.WorkerID.IsZero() {
		return errors.New("agent: worker_id is required")
	}
	if lp.CapitalistID.IsZero() {
		return errors.New("agent: capitalist_id is required")
	}
	if lp.Duration <= 0 {
		return errors.New("agent: duration must be positive")
	}
	for _, rm := range lp.Means.RawMaterials {
		if rm.CommodityID == "" {
			return errors.New("agent: raw_material commodity_id is required")
		}
		if rm.SNLTPerUnit < 0 {
			return errors.New("agent: raw_material snlt_per_unit cannot be negative")
		}
	}
	return nil
}

// Product is the output use-value of a LabourProcess [§1.b].
type Product struct {
	CommodityKind string        `json:"commodity_kind"`
	Quantity      int64         `json:"quantity"`
	TotalValue    LabourMinutes `json:"total_value"`
}

// ValorizationProcess wraps a LabourProcess and computes value-magnitude
// quantities from it [§2].
type ValorizationProcess struct {
	Process                LabourProcess `json:"process"`
	NecessaryLabourMinutes LabourMinutes `json:"necessary_labour_minutes"`
}

// NecessaryLabour returns the labour-time needed to reproduce the worker [§2.a].
func (vp ValorizationProcess) NecessaryLabour() LabourMinutes { return vp.NecessaryLabourMinutes }

// SurplusLabour returns the unpaid portion of the working day [§2.b].
func (vp ValorizationProcess) SurplusLabour() LabourMinutes {
	return SurplusLabour(vp.Process.Duration, vp.NecessaryLabourMinutes)
}

// SurplusValue returns the value produced above the value of labour-power [§2.f].
func (vp ValorizationProcess) SurplusValue() LabourMinutes { return vp.SurplusLabour() }

// ProductValue returns the total value of the product [§2.c].
func (vp ValorizationProcess) ProductValue() LabourMinutes {
	return TransferredValue(vp.Process.Means) + ValueAdded(vp.Process.Duration)
}

// TransferredValue sums the SNLT of all raw materials and instrument wear [§2.c].
// Constant capital transfers value to the product but creates no new value.
func TransferredValue(mp MeansOfProduction) LabourMinutes {
	var total LabourMinutes
	for _, rm := range mp.RawMaterials {
		total += rm.SNLTPerUnit * LabourMinutes(rm.Quantity)
	}
	for _, inst := range mp.Instruments {
		total += inst.WearPerRun
	}
	return total
}

// ValueAdded returns the new value created by living labour during duration [§2].
// For uniform skill level = 1 the reduction to abstract labour is the identity.
func ValueAdded(duration LabourMinutes) LabourMinutes { return duration }

// SurplusLabour is the unpaid portion of the working day: wd - nl [§2].
func SurplusLabour(workingDay, necessaryLabour LabourMinutes) LabourMinutes {
	return workingDay - necessaryLabour
}
