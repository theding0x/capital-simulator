package agent

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"time"
)

// LabourMinutes is the canonical value-magnitude unit: socially-necessary labour-time
// expressed in minutes.
type LabourMinutes int64

// AgentID is a 96-bit hex identifier for ch6 Worker and Capitalist agents.
type AgentID string

func (id AgentID) IsZero() bool { return id == "" }

func NewAgentID() AgentID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return AgentID(hex.EncodeToString(b))
}

// PurchaseID is a 96-bit hex identifier for LabourPowerPurchase records.
type PurchaseID string

func (id PurchaseID) IsZero() bool { return id == "" }

func NewPurchaseID() PurchaseID {
	b := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err)
	}
	return PurchaseID(hex.EncodeToString(b))
}

// AgentKind distinguishes the two sides of the labour market.
type AgentKind string

const (
	AgentKindWorker     AgentKind = "worker"
	AgentKindCapitalist AgentKind = "capitalist"
)

// ErrInvalidContract is returned when ContractDays is not positive and finite.
var ErrInvalidContract = errors.New("agent: contract duration must be positive and finite")

// LabourAgent holds the shared identity fields embedded by Worker and Capitalist.
// (The spec calls this "Agent"; the name LabourAgent avoids a conflict with the
// existing agent.Agent type introduced in Ch. 4.)
type LabourAgent struct {
	ID        AgentID   `json:"id"`
	Kind      AgentKind `json:"kind"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LabourPower is the aggregate of mental and physical capabilities existing in a
// human being that is put to work during the contracted period [§2].
type LabourPower struct {
	CapacityMinutesPerDay LabourMinutes `json:"capacity_minutes_per_day"`
}

// Worker is the agent who owns their labour-power and sells it as a commodity.
// They are "free" in the double sense: free to sell, and freed from other means
// of subsistence [§3, §4].
type Worker struct {
	LabourAgent
	OwnsLabourPower       bool        `json:"owns_labour_power"`
	OwnsCommoditiesToSell bool        `json:"owns_commodities_to_sell"`
	LabourPower           LabourPower `json:"labour_power"`
}

func (w Worker) Validate() error {
	if w.LabourPower.CapacityMinutesPerDay <= 0 {
		return errors.New("agent: labour_power capacity_minutes_per_day must be positive")
	}
	return nil
}

// Capitalist is the agent who possesses money-capital and purchases labour-power.
type Capitalist struct {
	LabourAgent
	MoneyCapital LabourMinutes `json:"money_capital"`
}

func (c Capitalist) Validate() error {
	if c.MoneyCapital < 0 {
		return errors.New("agent: money_capital cannot be negative")
	}
	return nil
}

// IsFreeLabourer returns true when the worker satisfies Marx's "double freedom":
// owns their labour-power AND lacks other commodities to sell [§3, §4].
func IsFreeLabourer(w Worker) bool {
	return w.OwnsLabourPower && !w.OwnsCommoditiesToSell
}

// SubsistenceItem is a single commodity in the worker's subsistence basket.
// Essential marks physically indispensable items used by MinimumValue [§5].
type SubsistenceItem struct {
	Name        string        `json:"name"`
	SNLTMinutes LabourMinutes `json:"snlt_minutes"`
	Essential   bool          `json:"essential"`
}

// SubsistenceBasket is the set of commodities required to maintain and reproduce
// the worker [§5].
type SubsistenceBasket []SubsistenceItem

// TotalSNLT returns the sum of SNLT across all items in the basket.
func (b SubsistenceBasket) TotalSNLT() LabourMinutes {
	var total LabourMinutes
	for _, item := range b {
		total += item.SNLTMinutes
	}
	return total
}

// LabourPowerValue computes the value of labour-power from a SubsistenceBasket [§5].
type LabourPowerValue struct {
	Basket SubsistenceBasket `json:"basket"`
}

// DailyValue returns the value of one day's labour-power: the sum of SNLT of all
// subsistence items. Invariant: DailyValue() == Basket.TotalSNLT() [§5].
func (v LabourPowerValue) DailyValue() LabourMinutes {
	return v.Basket.TotalSNLT()
}

// MinimumValue returns the physical floor of labour-power value: SNLT of
// physically indispensable items only. Invariant: MinimumValue() <= DailyValue() [§5].
func (v LabourPowerValue) MinimumValue() LabourMinutes {
	var total LabourMinutes
	for _, item := range v.Basket {
		if item.Essential {
			total += item.SNLTMinutes
		}
	}
	return total
}

// ReproductionCost returns the total SNLT to reproduce labour-power for another
// day. Equals DailyValue() under normal conditions [§5].
func (v LabourPowerValue) ReproductionCost() LabourMinutes {
	return v.DailyValue()
}

// LabourPowerOffering is a worker's labour-power posted for sale.
// ID is added for persistence; it is not in the spec's field list but required
// for CRUD operations.
type LabourPowerOffering struct {
	ID                    AgentID       `json:"id"`
	OwnerID               AgentID       `json:"owner_id"`
	CapacityMinutesPerDay LabourMinutes `json:"capacity_minutes_per_day"`
	ContractDays          int64         `json:"contract_days"`
	AskingWage            LabourMinutes `json:"asking_wage"`
	CreatedAt             time.Time     `json:"created_at"`
}

func (o LabourPowerOffering) Validate() error {
	if o.OwnerID.IsZero() {
		return errors.New("agent: offering owner_id is required")
	}
	if o.CapacityMinutesPerDay <= 0 {
		return errors.New("agent: capacity_minutes_per_day must be positive")
	}
	if o.ContractDays <= 0 {
		return ErrInvalidContract
	}
	if o.AskingWage < 0 {
		return errors.New("agent: asking_wage cannot be negative")
	}
	return nil
}

// LabourPowerPurchase is the completed transaction of buying labour-power.
type LabourPowerPurchase struct {
	ID           PurchaseID    `json:"id"`
	SellerID     AgentID       `json:"seller_id"`
	BuyerID      AgentID       `json:"buyer_id"`
	WageMinutes  LabourMinutes `json:"wage_minutes"`
	ContractDays int64         `json:"contract_days"`
	CreatedAt    time.Time     `json:"created_at"`
}

func (p LabourPowerPurchase) Validate() error {
	if p.SellerID.IsZero() {
		return errors.New("agent: purchase seller_id is required")
	}
	if p.BuyerID.IsZero() {
		return errors.New("agent: purchase buyer_id is required")
	}
	if p.WageMinutes < 0 {
		return errors.New("agent: wage_minutes cannot be negative")
	}
	if p.ContractDays <= 0 {
		return ErrInvalidContract
	}
	return nil
}
