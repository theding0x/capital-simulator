// Package credit implements Capital Vol. III, Part V.
// This file covers Ch. 27 ("The Role of Credit in Capitalist Production").
//
// Credit plays three roles:
//
//	(I)  It equalises the rate of profit by enabling capital to flow across
//	     spheres without the friction of cash settlement.
//	(II) It reduces circulation costs by substituting credit-money (bills,
//	     bank-notes, deposits) for metallic coin.
//	(III) It enables joint-stock companies — the abolition of private capital
//	     within capitalism itself. The functioning capitalist becomes a salaried
//	     manager; the owner becomes a pure money-capitalist drawing interest.
//	     Co-operative factories are the positive, worker-owned form of the
//	     transition beyond capital.
//
// Key invariants:
//   - StockCompany.TotalCapital == FixedCapital + FloatingCapital.
//   - Manager.WageMinutes > 0 (wages of skilled superintendence, not profit).
//   - CooperativeFactory.WorkerOwned == true.
//   - CreditRole ∈ {1, 2, 3}.
package credit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Ch. 27 sentinel errors.
var (
	// ErrCapitalMismatch is returned when TotalCapital != FixedCapital + FloatingCapital.
	ErrCapitalMismatch = errors.New("credit: total_capital must equal fixed_capital + floating_capital")

	// ErrNonPositiveWage is returned when Manager.WageMinutes <= 0.
	ErrNonPositiveWage = errors.New("credit: manager.wage_minutes must be positive")

	// ErrNotWorkerOwned is returned when CooperativeFactory.WorkerOwned is false.
	ErrNotWorkerOwned = errors.New("credit: cooperative_factory must be worker_owned")

	// ErrInvalidCreditRole is returned when a CreditRole is outside {1, 2, 3}.
	ErrInvalidCreditRole = errors.New("credit: credit_role must be 1, 2, or 3")
)

// CreditRole enumerates the three roles credit plays in capitalist production
// (Vol. III Ch. 27).
type CreditRole int

const (
	// RoleEqualiseProfitRate: credit lets capital flow across spheres, equalising p'.
	RoleEqualiseProfitRate CreditRole = 1
	// RoleReduceCirculationCosts: credit substitutes for metallic money, cutting costs.
	RoleReduceCirculationCosts CreditRole = 2
	// RoleStockCompanies: credit enables joint-stock companies and co-operatives.
	RoleStockCompanies CreditRole = 3
)

// IsValid reports whether r is a defined CreditRole.
func (r CreditRole) IsValid() bool {
	return r == RoleEqualiseProfitRate || r == RoleReduceCirculationCosts || r == RoleStockCompanies
}

// ── Manager ───────────────────────────────────────────────────────────────────

// Manager is the employed functionary of the joint-stock company — a salaried
// worker of superintendence, not a profit-drawing owner.
//
// Invariant: WageMinutes > 0.
type Manager struct {
	Name        string `json:"name"`
	Title       string `json:"title"`
	WageMinutes int64  `json:"wage_minutes"` // > 0; wages of skilled labour-power
}

// ── StockCompanyID ────────────────────────────────────────────────────────────

// StockCompanyID identifies a persisted StockCompany.
type StockCompanyID string

// NewStockCompanyID returns a fresh random StockCompanyID.
func NewStockCompanyID() StockCompanyID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return StockCompanyID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty StockCompanyID.
func (id StockCompanyID) IsZero() bool { return id == "" }

// ── StockCompany ──────────────────────────────────────────────────────────────

// StockCompany is a joint-stock company (Vol. III Ch. 27) — the social form
// in which private capital is abolished within capitalism. The owner becomes a
// mere money-capitalist; the Manager runs the firm for wages.
//
// Invariants:
//
//	TotalCapital == FixedCapital + FloatingCapital
//	Manager.WageMinutes > 0
type StockCompany struct {
	ID              StockCompanyID `json:"id"`
	Name            string         `json:"name"`
	FixedCapital    int64          `json:"fixed_capital"`    // pence
	FloatingCapital int64          `json:"floating_capital"` // pence
	TotalCapital    int64          `json:"total_capital"`    // pence; MUST == Fixed + Floating
	Manager         Manager        `json:"manager"`
	Description     string         `json:"description"`
	CreatedAt       time.Time      `json:"created_at"`
}

// NewStockCompany constructs a StockCompany, enforcing invariants.
//
//	fixedCapital + floatingCapital != totalCapital → ErrCapitalMismatch
//	manager.WageMinutes <= 0                       → ErrNonPositiveWage
func NewStockCompany(
	name string,
	fixedCapital, floatingCapital, totalCapital int64,
	manager Manager,
	description string,
) (StockCompany, error) {
	if fixedCapital+floatingCapital != totalCapital {
		return StockCompany{}, ErrCapitalMismatch
	}
	if manager.WageMinutes <= 0 {
		return StockCompany{}, ErrNonPositiveWage
	}
	return StockCompany{
		Name:            name,
		FixedCapital:    fixedCapital,
		FloatingCapital: floatingCapital,
		TotalCapital:    totalCapital,
		Manager:         manager,
		Description:     description,
	}, nil
}

// ── CooperativeFactoryID ──────────────────────────────────────────────────────

// CooperativeFactoryID identifies a persisted CooperativeFactory.
type CooperativeFactoryID string

// NewCooperativeFactoryID returns a fresh random CooperativeFactoryID.
func NewCooperativeFactoryID() CooperativeFactoryID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return CooperativeFactoryID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CooperativeFactoryID.
func (id CooperativeFactoryID) IsZero() bool { return id == "" }

// ── CooperativeFactory ────────────────────────────────────────────────────────

// CooperativeFactory is a worker-owned productive enterprise — the positive
// sprout of the new mode of production within the old (Vol. III Ch. 27).
//
// Invariant: WorkerOwned == true.
type CooperativeFactory struct {
	ID              CooperativeFactoryID `json:"id"`
	Name            string               `json:"name"`
	WorkerOwned     bool                 `json:"worker_owned"`     // INVARIANT: must be true
	WorkerCount     int64                `json:"worker_count"`     // number of worker-owners
	CapitalAdvanced int64                `json:"capital_advanced"` // pence; total advanced capital
	Description     string               `json:"description"`
	CreatedAt       time.Time            `json:"created_at"`
}

// NewCooperativeFactory constructs a CooperativeFactory, enforcing invariants.
//
//	workerOwned == false → ErrNotWorkerOwned
func NewCooperativeFactory(
	name string,
	workerOwned bool,
	workerCount, capitalAdvanced int64,
	description string,
) (CooperativeFactory, error) {
	if !workerOwned {
		return CooperativeFactory{}, ErrNotWorkerOwned
	}
	return CooperativeFactory{
		Name:            name,
		WorkerOwned:     workerOwned,
		WorkerCount:     workerCount,
		CapitalAdvanced: capitalAdvanced,
		Description:     description,
	}, nil
}

// ── CreditFunctionAnalysisID ──────────────────────────────────────────────────

// CreditFunctionAnalysisID identifies a CreditFunctionAnalysis value.
type CreditFunctionAnalysisID string

// NewCreditFunctionAnalysisID returns a fresh random CreditFunctionAnalysisID.
func NewCreditFunctionAnalysisID() CreditFunctionAnalysisID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("credit: rand: %w", err))
	}
	return CreditFunctionAnalysisID(hex.EncodeToString(b[:]))
}

// ── CirculationCostSaving ─────────────────────────────────────────────────────

// CirculationCostSaving captures the pence saved when credit-money substitutes
// for metallic coin in circulation (Vol. III Ch. 27, role II).
//
// SavingAmount = MetallicEquivalent - ActualCurrency.
type CirculationCostSaving struct {
	MetallicEquivalent int64  `json:"metallic_equivalent"` // pence; what metallic circulation would cost
	ActualCurrency     int64  `json:"actual_currency"`     // pence; what credit-money circulation costs
	SavingAmount       int64  `json:"saving_amount"`       // DERIVED: MetallicEquivalent - ActualCurrency
	Description        string `json:"description"`
}

// ComputeCirculationCostSaving derives SavingAmount from metallic and actual amounts.
func ComputeCirculationCostSaving(metallicEquivalent, actualCurrency int64, description string) CirculationCostSaving {
	return CirculationCostSaving{
		MetallicEquivalent: metallicEquivalent,
		ActualCurrency:     actualCurrency,
		SavingAmount:       metallicEquivalent - actualCurrency,
		Description:        description,
	}
}

// ── CreditFunctionAnalysis ────────────────────────────────────────────────────

// CreditFunctionAnalysis is a compute-only (not persisted) value type that
// quantifies credit's aggregate effect across all three roles (Vol. III Ch. 27).
type CreditFunctionAnalysis struct {
	ID                CreditFunctionAnalysisID `json:"id"`
	Role              CreditRole               `json:"role"` // which of the three roles
	CirculationSaving CirculationCostSaving    `json:"circulation_saving"`
	StockCompanyCount int64                    `json:"stock_company_count"`
	CooperativeCount  int64                    `json:"cooperative_count"`
	Description       string                   `json:"description"`
}
