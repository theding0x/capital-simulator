// Vol. II Ch. 10 — Theories of Fixed and Circulating Capital:
// The Physiocrats and Adam Smith.
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
)

// AllKnownEconomistErrors lists every defined error constant. Tests use this
// to verify that every value has at least one EconomistAttribution row citing it.
var AllKnownEconomistErrors = []KnownEconomistError{
	ErrorSmithConflation,
	ErrorSmithCirculationCapitalConflation,
	ErrorSmithRevenueInCapital,
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
