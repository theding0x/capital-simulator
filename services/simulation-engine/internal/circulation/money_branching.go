package circulation

import "errors"

// Vol. II Ch. 3 — the money-branching at M′ that defines Form III.
//
// "Not until its transformation into money has been accomplished does the
// movement branch out into movements of capital and of revenue." The fork
// occurs after C′—M′ realisation and before the M—C purchase.

var (
	ErrBranchingExceedsSurplus = errors.New("circulation: revenue cannot exceed the realised surplus-value")
	ErrBranchingModeMismatch   = errors.New("circulation: revenue branch contradicts the reproduction mode")
	ErrBranchingUnbalanced     = errors.New("circulation: reproduction fund + revenue + capitalised must equal the realised money")
	ErrBranchingNegative       = errors.New("circulation: branching components cannot be negative")
)

// MoneyBranching is the fork of the realised M′ into its three destinations.
//
//   - ReproductionFundPence is the capital-value (c + v) that returns to the
//     circuit as the next M—C(L+MP).
//   - RevenuePence leaves the circuit as individual consumption via c—m—c.
//   - CapitalisedPence is surplus-value re-advanced as additional capital
//     (extended reproduction).
//
// Invariant: ReproductionFundPence + RevenuePence + CapitalisedPence == RealisedPence.
type MoneyBranching struct {
	CommodityCircuitID    CommodityCircuitID `json:"commodity_circuit_id"`
	RealisedPence         Pence              `json:"realised_pence"`
	ReproductionFundPence Pence              `json:"reproduction_fund_pence"`
	RevenuePence          Pence              `json:"revenue_pence"`
	CapitalisedPence      Pence              `json:"capitalised_pence"`
}

// Validate checks that no component is negative and the three branches
// exhaust the realised money exactly.
func (b MoneyBranching) Validate() error {
	if b.RealisedPence < 0 || b.ReproductionFundPence < 0 || b.RevenuePence < 0 || b.CapitalisedPence < 0 {
		return ErrBranchingNegative
	}
	if b.ReproductionFundPence+b.RevenuePence+b.CapitalisedPence != b.RealisedPence {
		return ErrBranchingUnbalanced
	}
	return nil
}

// BranchMoney forks the fully realised C′ into reproduction fund, revenue and
// capitalised surplus. revenuePence is the share of surplus-value consumed as
// revenue; the remainder of the surplus is capitalised. The reproduction fund
// is the original capital-value c + v. Simple reproduction requires all
// surplus to go to revenue, extended reproduction none of it; mixed allows any
// split.
func (c CommodityCircuit) BranchMoney(revenuePence Pence) (MoneyBranching, error) {
	surplus := c.Initial.SurplusPence
	if revenuePence < 0 {
		return MoneyBranching{}, ErrBranchingNegative
	}
	if revenuePence > surplus {
		return MoneyBranching{}, ErrBranchingExceedsSurplus
	}
	switch c.Mode {
	case ReproductionSimple:
		if revenuePence != surplus {
			return MoneyBranching{}, ErrBranchingModeMismatch
		}
	case ReproductionExtended:
		if revenuePence != 0 {
			return MoneyBranching{}, ErrBranchingModeMismatch
		}
	}
	b := MoneyBranching{
		CommodityCircuitID:    c.ID,
		RealisedPence:         c.Initial.Total(),
		ReproductionFundPence: c.Initial.ValueOriginalPence(),
		RevenuePence:          revenuePence,
		CapitalisedPence:      surplus - revenuePence,
	}
	return b, b.Validate()
}
