// Package reproduction — Vol. II Ch. 18: The Role of Money-Capital in Reproduction.
//
// Marx poses the question: what role does money-capital play in the annual
// reproduction process? The chapter shows that money-capital does not
// produce value, but it is indispensable for starting the circuit and
// keeping it in motion. The two departments of social production (Dept I —
// means of production; Dept II — articles of consumption) must exchange
// with one another, and money is the medium through which that exchange is
// effected at the social level.
//
// Key concepts modelled here:
//   - MoneySupplyApportionment — how the total circulating money-stock is
//     apportioned among the departments, the wage-rotation fund, and idle
//     hoards.
//   - DepartmentMoneyReserve — the money reserve held by each department
//     for a specific operational purpose (wages, MP purchase, etc.).
//   - CirculatingMoneyMass — M × V = effective circulating value; velocity
//     is expressed in basis points per year (10000 BP = 1 turnover p.a.).
//   - WageRotationFund — the portion of money capital that must be advanced
//     in advance of the sale of output, owing to the weekly/monthly wage cycle.
//   - InterDepartmentSettlement — the bilateral exchange flows through which
//     Dept I's variable capital and surplus flow to Dept II for articles of
//     consumption, and Dept II's constant-capital portion flows back to Dept I
//     for means of production.
//
// Monetary values are integer money units local to the reproduction package;
// the £-scale is package-local and not interchangeable with other packages'
// pence without conversion. See docs/architecture.md → "Money scale across
// packages". (The Ch. 20/21 reproduction-scheme fixtures scale Marx's figures
// ×1000, e.g. £1000 → 1_000_000.)
package reproduction

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

// --- Sentinel errors --------------------------------------------------------

var (
	// ErrNegativeMoneyStock is returned when money_stock_pence is negative.
	ErrNegativeMoneyStock = errors.New("reproduction: money_stock_pence cannot be negative")

	// ErrApportionmentMismatch is returned when the parts of a
	// MoneySupplyApportionment do not sum to TotalSocialMoneyPence.
	ErrApportionmentMismatch = errors.New("reproduction: apportionment parts do not sum to total")

	// ErrInvalidDepartment is returned when a CapitalDepartment value is
	// not a recognised constant.
	ErrInvalidDepartment = errors.New("reproduction: invalid capital department")

	// ErrInvalidReservePurpose is returned when a ReservePurpose value is
	// not a recognised constant.
	ErrInvalidReservePurpose = errors.New("reproduction: invalid reserve purpose")

	// ErrInvalidInterDeptPurpose is returned when an InterDepartmentPurpose
	// value is not a recognised constant.
	ErrInvalidInterDeptPurpose = errors.New("reproduction: invalid inter-department purpose")

	// ErrZeroVelocity is returned when velocity_per_year_basis_points is zero
	// or negative, making effective circulating value undefined.
	ErrZeroVelocity = errors.New("reproduction: velocity_per_year_basis_points must be positive")
)

// --- Enum: CapitalDepartment ------------------------------------------------

// CapitalDepartment identifies one of the two departments of social production
// as defined by Marx in Vol. II, Part III:
//   - DepartmentI  — means of production (machinery, raw materials, etc.)
//   - DepartmentII — articles of consumption (wage goods and luxury goods)
type CapitalDepartment string

const (
	// DepartmentI produces means of production.
	DepartmentI CapitalDepartment = "department_i"

	// DepartmentII produces articles of consumption.
	DepartmentII CapitalDepartment = "department_ii"
)

// validDepartments is the canonical set of allowed CapitalDepartment values.
var validDepartments = map[CapitalDepartment]bool{
	DepartmentI:  true,
	DepartmentII: true,
}

// Validate returns ErrInvalidDepartment if the department is not recognised.
func (d CapitalDepartment) Validate() error {
	if !validDepartments[d] {
		return ErrInvalidDepartment
	}
	return nil
}

// --- Enum: ReservePurpose ---------------------------------------------------

// ReservePurpose identifies the reason a department holds a money reserve.
type ReservePurpose string

const (
	// PurposeMeansOfProduction — the reserve covers purchase of means of
	// production from Dept I.
	PurposeMeansOfProduction ReservePurpose = "means_of_production"

	// PurposeArticlesOfConsumption — the reserve covers purchase of articles
	// of consumption from Dept II.
	PurposeArticlesOfConsumption ReservePurpose = "articles_of_consumption"

	// PurposeWagePayment — the reserve is held to pay wages before output
	// is sold; necessary because the wage cycle precedes realisation.
	PurposeWagePayment ReservePurpose = "wage_payment"

	// PurposeIdleHoard — latent money capital held as an idle hoard,
	// not yet re-entered into productive circulation.
	PurposeIdleHoard ReservePurpose = "idle_hoard"
)

// validReservePurposes is the canonical set of allowed ReservePurpose values.
var validReservePurposes = map[ReservePurpose]bool{
	PurposeMeansOfProduction:     true,
	PurposeArticlesOfConsumption: true,
	PurposeWagePayment:           true,
	PurposeIdleHoard:             true,
}

// Validate returns ErrInvalidReservePurpose if the purpose is not recognised.
func (p ReservePurpose) Validate() error {
	if !validReservePurposes[p] {
		return ErrInvalidReservePurpose
	}
	return nil
}

// --- Enum: InterDepartmentPurpose -------------------------------------------

// InterDepartmentPurpose identifies the direction and nature of an
// inter-department settlement flow.
type InterDepartmentPurpose string

const (
	// PurposeIVplusStoIIC — Dept I's variable capital (v) and surplus-value
	// (s) flow to Dept II to purchase articles of consumption. This is the
	// dominant flow in Marx's two-department scheme: I(v+s) → II(c).
	PurposeIVplusStoIIC InterDepartmentPurpose = "i_v_plus_s_to_ii_c"

	// PurposeIICtoIForMP — Dept II's constant-capital portion (c) flows back
	// to Dept I to purchase the means of production it needs. This is the
	// return flow: II(c) → I.
	PurposeIICtoIForMP InterDepartmentPurpose = "ii_c_to_i_for_mp"
)

// validInterDeptPurposes is the canonical set of allowed InterDepartmentPurpose values.
var validInterDeptPurposes = map[InterDepartmentPurpose]bool{
	PurposeIVplusStoIIC: true,
	PurposeIICtoIForMP:  true,
}

// Validate returns ErrInvalidInterDeptPurpose if the purpose is not recognised.
func (p InterDepartmentPurpose) Validate() error {
	if !validInterDeptPurposes[p] {
		return ErrInvalidInterDeptPurpose
	}
	return nil
}

// --- ID types ---------------------------------------------------------------

// MoneySupplyApportionmentID identifies a MoneySupplyApportionment record.
type MoneySupplyApportionmentID string

// NewMoneySupplyApportionmentID returns a 96-bit random hex identifier.
func NewMoneySupplyApportionmentID() MoneySupplyApportionmentID {
	return MoneySupplyApportionmentID(newMoneySupplyHexID())
}

// DepartmentMoneyReserveID identifies a DepartmentMoneyReserve record.
type DepartmentMoneyReserveID string

// NewDepartmentMoneyReserveID returns a 96-bit random hex identifier.
func NewDepartmentMoneyReserveID() DepartmentMoneyReserveID {
	return DepartmentMoneyReserveID(newMoneySupplyHexID())
}

// CirculatingMoneyMassID identifies a CirculatingMoneyMass record.
type CirculatingMoneyMassID string

// NewCirculatingMoneyMassID returns a 96-bit random hex identifier.
func NewCirculatingMoneyMassID() CirculatingMoneyMassID {
	return CirculatingMoneyMassID(newMoneySupplyHexID())
}

// WageRotationFundID identifies a WageRotationFund record.
type WageRotationFundID string

// NewWageRotationFundID returns a 96-bit random hex identifier.
func NewWageRotationFundID() WageRotationFundID {
	return WageRotationFundID(newMoneySupplyHexID())
}

// InterDepartmentSettlementID identifies an InterDepartmentSettlement record.
type InterDepartmentSettlementID string

// NewInterDepartmentSettlementID returns a 96-bit random hex identifier.
func NewInterDepartmentSettlementID() InterDepartmentSettlementID {
	return InterDepartmentSettlementID(newMoneySupplyHexID())
}

// --- MoneySupplyApportionment -----------------------------------------------

// MoneySupplyApportionment records how the total circulating money-stock is
// divided among the two departments, the wage-rotation fund, and idle hoards.
//
// Invariant (enforced by Validate):
//
//	TotalSocialMoneyPence ==
//	  DeptIReservePence + DeptIIReservePence +
//	  WageRotationPence + IdleHoardPence
type MoneySupplyApportionment struct {
	ID                    MoneySupplyApportionmentID `json:"id"`
	TotalSocialMoneyPence int64                      `json:"total_social_money_pence"`
	DeptIReservePence     int64                      `json:"dept_i_reserve_pence"`
	DeptIIReservePence    int64                      `json:"dept_ii_reserve_pence"`
	WageRotationPence     int64                      `json:"wage_rotation_pence"`
	IdleHoardPence        int64                      `json:"idle_hoard_pence"`
	Period                string                     `json:"period"`
}

// Validate checks that the four parts sum to the total.
// Returns ErrApportionmentMismatch if the invariant is violated.
func (a MoneySupplyApportionment) Validate() error {
	parts := a.DeptIReservePence +
		a.DeptIIReservePence +
		a.WageRotationPence +
		a.IdleHoardPence
	if parts != a.TotalSocialMoneyPence {
		return ErrApportionmentMismatch
	}
	return nil
}

// --- DepartmentMoneyReserve -------------------------------------------------

// DepartmentMoneyReserve records the money reserve held by one department
// for a specific operational purpose during a given period.
type DepartmentMoneyReserve struct {
	ID             DepartmentMoneyReserveID `json:"id"`
	Department     CapitalDepartment        `json:"department"`
	ReservePurpose ReservePurpose           `json:"reserve_purpose"`
	ReservePence   int64                    `json:"reserve_pence"`
	Period         string                   `json:"period"`
}

// --- CirculatingMoneyMass ---------------------------------------------------

// CirculatingMoneyMass captures the money stock and its velocity, yielding
// the effective circulating value M×V.
//
// VelocityPerYearBasisPoints is expressed in basis points (10000 = 1 turnover
// per year; 20000 = 2 turnovers per year). This avoids float64 in the domain.
//
// Invariant: EffectiveCirculatingValuePence == MoneyStockPence * VelocityPerYearBasisPoints / 10000
type CirculatingMoneyMass struct {
	ID                             CirculatingMoneyMassID `json:"id"`
	MoneyStockPence                int64                  `json:"money_stock_pence"`
	VelocityPerYearBasisPoints     int64                  `json:"velocity_per_year_basis_points"`
	EffectiveCirculatingValuePence int64                  `json:"effective_circulating_value_pence"`
	Period                         string                 `json:"period"`
}

// ComputeEffectiveCirculatingValue returns stock * velocityBP / 10000.
// Returns 0 if stock is 0. Panics if velocityBP <= 0 (caller must validate).
func ComputeEffectiveCirculatingValue(stock, velocityBP int64) int64 {
	return stock * velocityBP / 10000
}

// --- WageRotationFund -------------------------------------------------------

// WageRotationFund records the portion of money capital a department must
// hold in advance of output sales, owing to the wage-payment cycle.
//
// WageCycleFrequency is the number of wage payments per year
// (e.g. 52 for weekly, 12 for monthly, 1 for annual).
type WageRotationFund struct {
	ID                 WageRotationFundID `json:"id"`
	FundPence          int64              `json:"fund_pence"`
	WageCycleFrequency int64              `json:"wage_cycle_frequency"`
	Department         CapitalDepartment  `json:"department"`
	Period             string             `json:"period"`
}

// TurnoverPerYearBasisPoints converts WageCycleFrequency to basis points
// (1 turnover/year = 10 000 bp), matching the basis-point convention used
// throughout the reproduction package. Weekly (52) → 520 000; monthly (12)
// → 120 000. A faster wage cycle produces a higher basis-point value.
func (w WageRotationFund) TurnoverPerYearBasisPoints() int64 {
	return w.WageCycleFrequency * 10000
}

// --- InterDepartmentSettlement ----------------------------------------------

// InterDepartmentSettlement records one bilateral exchange flow between the
// two departments of social production. These flows are the mechanism through
// which simple reproduction is maintained at the social level.
//
// In the canonical two-department scheme:
//   - I(v+s) → II: Dept I workers and capitalists buy articles of
//     consumption from Dept II; money flows I → II, commodities II → I.
//   - II(c) → I:   Dept II purchases means of production from Dept I;
//     money flows II → I, commodities I → II.
type InterDepartmentSettlement struct {
	ID                InterDepartmentSettlementID `json:"id"`
	FromDepartment    CapitalDepartment           `json:"from_department"`
	ToDepartment      CapitalDepartment           `json:"to_department"`
	SettledPence      int64                       `json:"settled_pence"`
	SettlementPurpose ReservePurpose              `json:"settlement_purpose"`
	Period            string                      `json:"period"`
}

// --- helpers ----------------------------------------------------------------

// newMoneySupplyHexID returns a 24-character lowercase hex string (96 bits).
func newMoneySupplyHexID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(b)
}
