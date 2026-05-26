// Vol. II Ch. 10–11 — Theories of Fixed and Circulating Capital:
// The Physiocrats, Adam Smith, and Ricardo.
//
// This file is documentation-only: no behavioural logic. The types surface
// pre-Marxian economists' fixed/circulating terminology, their anticipations
// of the correct distinction (Ch. 8), and the specific theoretical errors
// that Marx's analysis corrects.
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
)

// AllKnownEconomistErrors lists every defined error constant. Tests use this
// to verify that every value has at least one EconomistAttribution row citing it.
var AllKnownEconomistErrors = []KnownEconomistError{
	ErrorSmithConflation,
	ErrorSmithCirculationCapitalConflation,
	ErrorSmithRevenueInCapital,
	ErrorRicardoDurabilityCollapse,
	ErrorRicardoConflation,
	ErrorRicardoFixedCapitalPriceExplanation,
	ErrorRicardoNoAggregateTurnover,
	ErrorRicardoNoValueRevolution,
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
