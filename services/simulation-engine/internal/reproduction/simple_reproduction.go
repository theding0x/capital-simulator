// Package reproduction — Vol. II Ch. 20: Simple Reproduction.
//
// Marx derives the two-department reproduction scheme for total social capital:
// means of production (Dept I) and articles of consumption (Dept II) must
// satisfy the keystone equation
//
//	I(v + s) == II(c)
//
// for society to reproduce itself on the same scale each year.
//
// Chapter 20 resolves the realisation puzzle posed in Ch. 17 and provides
// the constructive refutation of Smith's revenue dogma (Ch. 19). The three
// exchange flows that constitute simple reproduction are modelled explicitly:
//  1. Intra-Dept-I replacement — Dept I buys means of production from itself
//     (value = I.c); no money crosses the boundary.
//  2. I(v+s) → II — Dept I workers and capitalists purchase articles of
//     consumption from Dept II; money flows I→II, commodities II→I.
//  3. II(c) → I — Dept II purchases means of production from Dept I;
//     money flows II→I, commodities I→II.
//
// The three flows constitute a closed money loop: money advanced by Dept II
// returns to its point of origin by the end of the year.
//
// All monetary values are in pence (1 GBP = 240 pence).
package reproduction

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
)

// --- Sentinel errors --------------------------------------------------------

var (
	// ErrDepartmentalImbalance is returned when the keystone equation
	// I(v+s) != II(c) is violated and a balanced scheme is required.
	ErrDepartmentalImbalance = errors.New("reproduction: I(v+s) != II(c) — scheme is not balanced")

	// ErrDepartmentTotalMismatch is returned when a DepartmentalCapital's
	// TotalPence != ConstantPence + VariablePence + SurplusPence.
	ErrDepartmentTotalMismatch = errors.New("reproduction: total_pence != c + v + s")

	// ErrNegativeTurnover is returned when a DepartmentalCapital carries a
	// negative TurnoverNumberBP. Zero is allowed and treated as annual (10000).
	ErrNegativeTurnover = errors.New("reproduction: turnover_number_bp must not be negative")

	// ErrInvalidExchangeKind is returned when an ExchangeKind is not recognised.
	ErrInvalidExchangeKind = errors.New("reproduction: invalid exchange_kind")

	// ErrSchemeNotBalanced is returned when an operation that requires a
	// balanced scheme is attempted on an imbalanced one.
	ErrSchemeNotBalanced = errors.New("reproduction: scheme is not balanced — add both departments first")

	// ErrDepartmentAlreadySet is returned when a department has already been
	// registered on the scheme.
	ErrDepartmentAlreadySet = errors.New("reproduction: department already registered on this scheme")
)

// --- Enum: ExchangeKind -----------------------------------------------------

// ExchangeKind identifies the direction of value flow in an inter-department
// exchange: money flowing toward a department, or goods flowing from it.
type ExchangeKind string

const (
	// ExchangeMoneyForGoods — the initiating department advances money and
	// receives goods in return (e.g. Dept II advances money → receives MP from Dept I).
	ExchangeMoneyForGoods ExchangeKind = "money_for_goods"

	// ExchangeGoodsForMoney — the initiating department ships goods and
	// receives money in return (e.g. Dept I ships MP → receives money from Dept II).
	ExchangeGoodsForMoney ExchangeKind = "goods_for_money"
)

// validExchangeKinds is the canonical set of allowed ExchangeKind values.
var validExchangeKinds = map[ExchangeKind]bool{
	ExchangeMoneyForGoods: true,
	ExchangeGoodsForMoney: true,
}

// Validate returns ErrInvalidExchangeKind if the kind is not recognised.
func (k ExchangeKind) Validate() error {
	if !validExchangeKinds[k] {
		return ErrInvalidExchangeKind
	}
	return nil
}

// --- ID types ---------------------------------------------------------------

// SimpleReproductionSchemeID identifies a SimpleReproductionScheme record.
type SimpleReproductionSchemeID string

// NewSimpleReproductionSchemeID returns a 96-bit random hex identifier.
func NewSimpleReproductionSchemeID() SimpleReproductionSchemeID {
	return SimpleReproductionSchemeID(newSimpleReproHexID())
}

// DepartmentalCapitalID identifies a DepartmentalCapital record.
type DepartmentalCapitalID string

// NewDepartmentalCapitalID returns a 96-bit random hex identifier.
func NewDepartmentalCapitalID() DepartmentalCapitalID {
	return DepartmentalCapitalID(newSimpleReproHexID())
}

// InterDepartmentExchangeID identifies an InterDepartmentExchange record.
type InterDepartmentExchangeID string

// NewInterDepartmentExchangeID returns a 96-bit random hex identifier.
func NewInterDepartmentExchangeID() InterDepartmentExchangeID {
	return InterDepartmentExchangeID(newSimpleReproHexID())
}

// ReproductionTickID identifies a ReproductionTick record.
type ReproductionTickID string

// NewReproductionTickID returns a 96-bit random hex identifier.
func NewReproductionTickID() ReproductionTickID {
	return ReproductionTickID(newSimpleReproHexID())
}

// MoneyClosedLoopID identifies a MoneyClosedLoop record.
type MoneyClosedLoopID string

// NewMoneyClosedLoopID returns a 96-bit random hex identifier.
func NewMoneyClosedLoopID() MoneyClosedLoopID {
	return MoneyClosedLoopID(newSimpleReproHexID())
}

// --- DepartmentalCapital ----------------------------------------------------

// DepartmentalCapital records the c+v+s composition for one department in
// one reproduction period.
//
// Invariant (enforced by Validate):
//
//	TotalPence == ConstantPence + VariablePence + SurplusPence
//
// Marx's canonical fixture (Vol. II Ch. 20):
//
//	Dept I:  4 000c + 1 000v + 1 000s = 6 000
//	Dept II: 2 000c +   500v +   500s = 3 000
type DepartmentalCapital struct {
	ID            DepartmentalCapitalID      `json:"id"`
	SchemeID      SimpleReproductionSchemeID `json:"scheme_id"`
	Department    CapitalDepartment          `json:"department"`
	ConstantPence int64                      `json:"constant_pence"`
	VariablePence int64                      `json:"variable_pence"`
	SurplusPence  int64                      `json:"surplus_pence"`
	TotalPence    int64                      `json:"total_pence"`
	// TurnoverNumberBP is the department's turnover number in basis points
	// (10000 = one turnover per year = annual). When Dept I and Dept II turn
	// over at different speeds the keystone exchange settles on annualised
	// flows rather than raw per-period magnitudes (issue #221). A zero value
	// is treated as annual.
	TurnoverNumberBP int64 `json:"turnover_number_bp"`
}

// Validate checks the total invariant.
func (d DepartmentalCapital) Validate() error {
	if err := d.Department.Validate(); err != nil {
		return err
	}
	if d.TotalPence != d.ConstantPence+d.VariablePence+d.SurplusPence {
		return ErrDepartmentTotalMismatch
	}
	if d.TurnoverNumberBP < 0 {
		return ErrNegativeTurnover
	}
	return nil
}

// VplusSPence returns v + s — the value that must cross the departmental
// boundary into Dept II for articles of consumption.
func (d DepartmentalCapital) VplusSPence() int64 {
	return d.VariablePence + d.SurplusPence
}

// AdvancedPence returns c + v — the capital actually advanced (laid out) by the
// department in the period, excluding the surplus-value it produces. This is
// the magnitude the SocialCapitalAggregate department shares partition
// (issue #215).
func (d DepartmentalCapital) AdvancedPence() int64 {
	return d.ConstantPence + d.VariablePence
}

// turnoverNumberOrAnnual returns the turnover number in basis points, treating
// a non-positive (unset) value as annual (10000 = one turnover per year).
func (d DepartmentalCapital) turnoverNumberOrAnnual() int64 {
	if d.TurnoverNumberBP <= 0 {
		return 10000
	}
	return d.TurnoverNumberBP
}

// AnnualisedVplusSPence returns (v + s) scaled by the department's turnover
// number — the annual flow of variable capital plus surplus that the keystone
// exchange settles against when departments turn over at different speeds
// (issue #221).
func (d DepartmentalCapital) AnnualisedVplusSPence() int64 {
	return d.VplusSPence() * d.turnoverNumberOrAnnual() / 10000
}

// AnnualisedConstantPence returns constant capital scaled by the department's
// turnover number — the annual demand for means of production.
func (d DepartmentalCapital) AnnualisedConstantPence() int64 {
	return d.ConstantPence * d.turnoverNumberOrAnnual() / 10000
}

// --- SimpleReproductionScheme -----------------------------------------------

// SimpleReproductionScheme is the aggregating root that links both departments
// and records the balance status.
//
// IsBalanced is true when both departments have been registered and the
// keystone equation holds: DepartmentI.VplusSPence() == DepartmentII.ConstantPence.
//
// TickCount advances each time AdvanceReproductionTick is called.
type SimpleReproductionScheme struct {
	ID           SimpleReproductionSchemeID `json:"id"`
	Period       string                     `json:"period"`
	DepartmentI  *DepartmentalCapital       `json:"department_i,omitempty"`
	DepartmentII *DepartmentalCapital       `json:"department_ii,omitempty"`
	IsBalanced   bool                       `json:"is_balanced"`
	TickCount    int64                      `json:"tick_count"`
}

// CheckBalance returns true iff both departments are present and
// DepartmentI.VplusSPence() == DepartmentII.ConstantPence.
func (s SimpleReproductionScheme) CheckBalance() bool {
	if s.DepartmentI == nil || s.DepartmentII == nil {
		return false
	}
	return s.DepartmentI.VplusSPence() == s.DepartmentII.ConstantPence
}

// WageBillPence is the scheme's total variable capital (Dept I.v + Dept II.v) —
// the wage bill advanced to the working class each period (issue #220).
func (s SimpleReproductionScheme) WageBillPence() int64 {
	var v int64
	if s.DepartmentI != nil {
		v += s.DepartmentI.VariablePence
	}
	if s.DepartmentII != nil {
		v += s.DepartmentII.VariablePence
	}
	return v
}

// --- BalanceCheckResult -----------------------------------------------------

// BalanceCheckResult is the response payload for the balance-check endpoint.
// It surfaces the keystone equation components for diagnostic purposes.
type BalanceCheckResult struct {
	SchemeID            string `json:"scheme_id"`
	IsBalanced          bool   `json:"is_balanced"`
	DeptIVplusSPence    int64  `json:"dept_i_v_plus_s_pence"`
	DeptIIConstantPence int64  `json:"dept_ii_constant_pence"`
	DeficitPence        int64  `json:"deficit_pence"`
	// Annualised settlement (issue #221): the keystone flows scaled by each
	// department's turnover number. When both departments turn over at the same
	// speed these equal the raw figures above; when they differ, a scheme that
	// is balanced per period can be unbalanced on an annual basis.
	DeptITurnoverBP               int64 `json:"dept_i_turnover_bp"`
	DeptIITurnoverBP              int64 `json:"dept_ii_turnover_bp"`
	DeptIVplusSAnnualisedPence    int64 `json:"dept_i_v_plus_s_annualised_pence"`
	DeptIIConstantAnnualisedPence int64 `json:"dept_ii_constant_annualised_pence"`
	AnnualisedDeficitPence        int64 `json:"annualised_deficit_pence"`
	AnnualisedBalanced            bool  `json:"annualised_balanced"`
}

// ApplyAnnualisedSettlement fills the annualised settlement fields of res from
// the two departments' turnover numbers (issue #221). Either department may be
// nil (an incomplete scheme); a missing department contributes zero.
func ApplyAnnualisedSettlement(res *BalanceCheckResult, deptI, deptII *DepartmentalCapital) {
	if deptI != nil {
		res.DeptITurnoverBP = deptI.turnoverNumberOrAnnual()
		res.DeptIVplusSAnnualisedPence = deptI.AnnualisedVplusSPence()
	}
	if deptII != nil {
		res.DeptIITurnoverBP = deptII.turnoverNumberOrAnnual()
		res.DeptIIConstantAnnualisedPence = deptII.AnnualisedConstantPence()
	}
	res.AnnualisedDeficitPence = res.DeptIVplusSAnnualisedPence - res.DeptIIConstantAnnualisedPence
	res.AnnualisedBalanced = deptI != nil && deptII != nil && res.AnnualisedDeficitPence == 0
}

// --- InterDepartmentExchange ------------------------------------------------

// InterDepartmentExchange records one of the three money/goods flows that
// constitute a simple reproduction cycle.
//
// The three canonical exchanges for Marx's fixture are:
//  1. Intra-Dept-I — Dept I replaces its own constant capital internally
//     (I buys MP from itself); Pence = I.c = 4 000.
//  2. I(v+s) → II — Dept I's variable capital and surplus flow to Dept II;
//     Pence = I(v+s) = 2 000.
//  3. II(c) → I — Dept II's constant capital flows back to Dept I as
//     payment for means of production; Pence = II.c = 2 000.
type InterDepartmentExchange struct {
	ID             InterDepartmentExchangeID  `json:"id"`
	SchemeID       SimpleReproductionSchemeID `json:"scheme_id"`
	FromDepartment CapitalDepartment          `json:"from_department"`
	ToDepartment   CapitalDepartment          `json:"to_department"`
	Pence          int64                      `json:"pence"`
	Kind           ExchangeKind               `json:"kind"`
	Description    string                     `json:"description"`
}

// Validate checks department and kind fields.
func (e InterDepartmentExchange) Validate() error {
	if err := e.FromDepartment.Validate(); err != nil {
		return err
	}
	if err := e.ToDepartment.Validate(); err != nil {
		return err
	}
	if err := e.Kind.Validate(); err != nil {
		return err
	}
	return nil
}

// --- ReproductionTick -------------------------------------------------------

// ReproductionTick records one cycle advancement of a SimpleReproductionScheme.
// In simple reproduction the magnitudes are identical each cycle.
type ReproductionTick struct {
	ID         ReproductionTickID         `json:"id"`
	SchemeID   SimpleReproductionSchemeID `json:"scheme_id"`
	TickNumber int64                      `json:"tick_number"`
	Period     string                     `json:"period"`
	IsBalanced bool                       `json:"is_balanced"`
	// Labour-force reproduction (issue #220): each period's wage bill (the
	// scheme's total variable capital) is asserted to cover the worker pool's
	// subsistence. WorkerPoolSize == 0 means labour-force tracking was not
	// requested for this tick. Status is one of the TickStatus* values.
	WorkerPoolSize         int64  `json:"worker_pool_size"`
	WageBillPence          int64  `json:"wage_bill_pence"`
	SubsistenceBasketPence int64  `json:"subsistence_basket_pence"`
	SubsistenceCovered     bool   `json:"subsistence_covered"`
	Status                 string `json:"status"`
}

// Reproduction-tick status values (issue #220).
const (
	TickStatusReproduced      = "reproduced"
	TickStatusLabourShortfall = "labour-force-shortfall"
)

// ComputeLabourForceReproduction asserts that a period's wage bill covers the
// subsistence of the worker pool. workerPoolSize == 0 means tracking was not
// requested, in which case the tick is vacuously reproduced. Otherwise the wage
// bill must be at least workerPoolSize × subsistenceBasketPence (issue #220).
func ComputeLabourForceReproduction(wageBillPence, workerPoolSize, subsistenceBasketPence int64) (covered bool, status string) {
	if workerPoolSize <= 0 {
		return true, TickStatusReproduced
	}
	if wageBillPence >= workerPoolSize*subsistenceBasketPence {
		return true, TickStatusReproduced
	}
	return false, TickStatusLabourShortfall
}

// --- MoneyClosedLoop --------------------------------------------------------

// MoneyClosedLoop verifies that money paid out by one department returns to its
// starting distribution within the annual cycle — the money circuit closes.
//
// IsClosed is true when NetFlowPence == 0, confirming that money advanced by
// Dept II to purchase MP from Dept I flows back entirely via I(v+s)→II.
type MoneyClosedLoop struct {
	ID           MoneyClosedLoopID          `json:"id"`
	SchemeID     SimpleReproductionSchemeID `json:"scheme_id"`
	Period       string                     `json:"period"`
	IsClosed     bool                       `json:"is_closed"`
	NetFlowPence int64                      `json:"net_flow_pence"`
}

// --- Pure domain functions --------------------------------------------------

// ComputeIsBalanced returns whether I(v+s) == II(c) for two departments.
// This is a pure function used by both the store and domain tests.
func ComputeIsBalanced(deptI, deptII DepartmentalCapital) bool {
	return deptI.VplusSPence() == deptII.ConstantPence
}

// ComputeDeficit returns I(v+s) − II(c). A positive value means Dept II is
// undersupplied (the scheme overproduces MP); a negative value means it is
// oversupplied (the scheme underproduces MP).
func ComputeDeficit(deptI, deptII DepartmentalCapital) int64 {
	return deptI.VplusSPence() - deptII.ConstantPence
}

// --- helpers ----------------------------------------------------------------

// newSimpleReproHexID returns a 24-character lowercase hex string (96 bits).
func newSimpleReproHexID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return hex.EncodeToString(b)
}
