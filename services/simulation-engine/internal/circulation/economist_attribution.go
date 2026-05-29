// Vol. II Ch. 10–11, 19 — Theories of Fixed and Circulating Capital and
// Former Presentations of Reproduction:
// The Physiocrats, Adam Smith, and Ricardo.
//
// This file is documentation-only: no behavioural logic. The types surface
// pre-Marxian economists' fixed/circulating terminology and reproduction schemes,
// their anticipations of the correct distinctions (Ch. 8, 17–18), and the
// specific theoretical errors that Marx's analysis corrects.
package circulation

import (
	"crypto/rand"
	"encoding/hex"
)

// KnownEconomistError is an enumeration of documented theoretical errors in
// pre-Marxian political economy identified and refuted in Ch. 10.
type KnownEconomistError string

const (
	// ErrorSmithConflation: Smith conflates the fixed/circulating distinction
	// with the constant/variable distinction. Fixed capital ≠ constant capital;
	// circulating capital ≠ variable capital.
	ErrorSmithConflation KnownEconomistError = "error_smith_conflation"

	// ErrorSmithCirculationCapitalConflation: Smith treats money-capital and
	// commodity-capital (both circulation-forms) as sub-categories of
	// "circulating capital" — a category that belongs to productive capital
	// alone. The error collapses the circuit-form distinction into the
	// productive-capital distinction.
	ErrorSmithCirculationCapitalConflation KnownEconomistError = "error_smith_circulation_capital_conflation"

	// ErrorSmithRevenueInCapital: Smith imports revenue flows (wages, profit,
	// rent) back into the capital circuit, treating them as portions of capital
	// itself. This violates the RevenueCircuit invariant established in Ch. 2:
	// revenue exits from C′—M′ and does not re-enter as M—C(Lp+Mp).
	ErrorSmithRevenueInCapital KnownEconomistError = "error_smith_revenue_in_capital"

	// ErrorRicardoDurabilityCollapse: Ricardo reduces the fixed/circulating
	// distinction to mere physical durability — instruments of labour last many
	// periods (fixed); wage-goods and raw materials are consumed at once
	// (circulating). This collapses the form-of-circulation invariant
	// established in Ch. 8: the kind of a capital component is derived from its
	// role in the labour-process, not from the material durability of the object.
	ErrorRicardoDurabilityCollapse KnownEconomistError = "error_ricardo_durability_collapse"

	// ErrorRicardoConflation: Ricardo repeats Smith's fixed/circulating ↔
	// constant/variable conflation. Fixed capital (instruments of labour) ≠
	// constant capital; circulating capital (as Ricardo uses the term) ≠
	// variable capital. Refuted by Ch. 8 KindForRole invariant.
	ErrorRicardoConflation KnownEconomistError = "error_ricardo_conflation"

	// ErrorRicardoFixedCapitalPriceExplanation: Ricardo correctly observes that
	// differences in fixed-capital intensity affect relative prices when wages
	// change. He mis-locates the cause in the fixed/circulating ratio rather
	// than in the constant/variable ratio and the extraction of surplus-value.
	// The full resolution (prices of production, average rate of profit) is
	// deferred to Vol. III Part Two.
	ErrorRicardoFixedCapitalPriceExplanation KnownEconomistError = "error_ricardo_fixed_capital_price_explanation"

	// ErrorRicardoNoAggregateTurnover: Ricardo lacks the Ch. 9
	// aggregate-turnover formula. Without it he cannot derive the mean-average
	// number of turnovers across heterogeneous capital components and
	// therefore mis-states the periodicity of capital returns.
	ErrorRicardoNoAggregateTurnover KnownEconomistError = "error_ricardo_no_aggregate_turnover"

	// ErrorRicardoNoValueRevolution: Ricardo treats value as static within a
	// single turnover. He cannot account for a ValueRevolutionEvent (Ch. 4) —
	// a change in the value of means of production mid-turnover that alters
	// MoneyCapitalTiedUp without altering the physical volume of production.
	ErrorRicardoNoValueRevolution KnownEconomistError = "error_ricardo_no_value_revolution"

	// Ch. 19 — Former Presentations of the Subject
	// (Quesnay, Smith, and Ricardo on total social reproduction)

	// ErrorQuesnayProductiveSterileDivision: Quesnay divides social labour into
	// a "productive class" (agriculture) and a "sterile class" (manufacture and
	// trade). Only agricultural labour produces a net product (produit net). Marx
	// shows this is the Physiocratic form of the surplus-value concept — correct
	// in grasping that surplus arises in production, but mistaken in limiting
	// production to agriculture. Manufacture creates value too; the productive/
	// sterile divide must be replaced by the constant/variable decomposition.
	ErrorQuesnayProductiveSterileDivision KnownEconomistError = "error_quesnay_productive_sterile_division"

	// ErrorQuesnayMissingValueTheory: Quesnay's Tableau Économique maps the
	// circuit of the annual product without a labour-theory-of-value foundation.
	// He cannot derive why the net product arises, only that it arises in
	// agriculture. The Tableau's flows prefigure the two-department reproduction
	// scheme (Chs. 20–21) but cannot be made rigorous without the value categories
	// of Vol. I (constant capital, variable capital, surplus-value).
	ErrorQuesnayMissingValueTheory KnownEconomistError = "error_quesnay_missing_value_theory"

	// ErrorSmithRevenueDogma: Smith resolves the entire annual product into wages,
	// profit, and rent — the "revenue dogma" (V. II, Ch. 19). He omits constant
	// capital entirely: the component of value that merely replaces the means of
	// production consumed. Without a constant-capital term there is no balancing
	// equation for the two departments (I(v+s) = II(c) cannot be derived).
	ErrorSmithRevenueDogma KnownEconomistError = "error_smith_revenue_dogma"

	// ErrorSmithMissingConstantReplacement: Consequence of the revenue dogma. If
	// the entire product resolves into revenue, there is no fund for replacing the
	// means of production consumed during the year. Smith cannot explain how
	// constant capital reproduces itself — it appears to dissolve into air. Marx
	// traces the error to Smith's conflation of individual capital with total
	// social capital: for an individual capitalist, bought inputs are the seller's
	// revenue, but for total social capital, constant capital must be reproduced
	// in kind from Department I output.
	ErrorSmithMissingConstantReplacement KnownEconomistError = "error_smith_missing_constant_replacement"

	// ErrorSmithLabourRegress: Smith attempts to escape the revenue-dogma impasse
	// by reducing the value of every means of production to the wages, profit, and
	// rent paid to produce it — and those inputs in turn to earlier wages, profit,
	// and rent, ad infinitum. The regress never terminates. Marx shows this makes
	// no logical or historical sense: constant capital cannot be eliminated by an
	// infinite historical decomposition; it is a present, synchronic component of
	// every act of production.
	ErrorSmithLabourRegress KnownEconomistError = "error_smith_labour_regress"

	// ErrorRicardoIncompleteReproduction: Ricardo perceives Smith's revenue dogma
	// and its contradictions but fails to develop a workable reproduction scheme.
	// He corrects Smith on the individual-capital level (replacement of worn
	// machinery is part of cost, not revenue) yet never applies this correction at
	// the level of total social capital. He therefore cannot derive the two-
	// department balance condition I(v+s) = II(c) that Chs. 20–21 will establish.
	ErrorRicardoIncompleteReproduction KnownEconomistError = "error_ricardo_incomplete_reproduction"
)

// AllKnownEconomistErrors lists every defined error constant. Tests use this
// to verify that every value has at least one EconomistAttribution row citing it.
var AllKnownEconomistErrors = []KnownEconomistError{
	// Ch. 10–11: fixed/circulating capital theories
	ErrorSmithConflation,
	ErrorSmithCirculationCapitalConflation,
	ErrorSmithRevenueInCapital,
	ErrorRicardoDurabilityCollapse,
	ErrorRicardoConflation,
	ErrorRicardoFixedCapitalPriceExplanation,
	ErrorRicardoNoAggregateTurnover,
	ErrorRicardoNoValueRevolution,
	// Ch. 19: former presentations of reproduction
	ErrorQuesnayProductiveSterileDivision,
	ErrorQuesnayMissingValueTheory,
	ErrorSmithRevenueDogma,
	ErrorSmithMissingConstantReplacement,
	ErrorSmithLabourRegress,
	ErrorRicardoIncompleteReproduction,
}

// EconomistAttributionID is a 96-bit hex identifier for an EconomistAttribution.
type EconomistAttributionID string

// IsZero reports whether the ID is the zero value.
func (id EconomistAttributionID) IsZero() bool { return id == "" }

// NewEconomistAttributionID generates a 96-bit hex identifier via crypto/rand.
func NewEconomistAttributionID() EconomistAttributionID {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		panic("crypto/rand unavailable: " + err.Error())
	}
	return EconomistAttributionID(hex.EncodeToString(b))
}

// EconomistAttribution documents a pre-Marxian political economist's
// fixed/circulating terminology, its anticipations of Marx's own distinction,
// and any known theoretical errors.
//
// Records are static reference data seeded at migration time; no behaviour is
// defined on this type — it is used for the Ch. 10 React display panel only.
type EconomistAttribution struct {
	ID          EconomistAttributionID `json:"id"`
	Concept     string                 `json:"concept"`
	Theorist    string                 `json:"theorist"`
	EditionYear int64                  `json:"edition_year"`
	Anticipates string                 `json:"anticipates"`
	Errors      []KnownEconomistError  `json:"errors"`
}
