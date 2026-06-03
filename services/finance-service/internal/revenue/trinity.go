// Package revenue — Vol. III Part VII types: the revenues and their sources.
//
// Chapter 48 ("The Trinity Formula") models the ideological surface of
// bourgeois economics: Capital → interest, Land → ground-rent, Labour →
// wages. The formula presents three independent, co-ordinate sources of
// revenue. Marx shows the appearance is a fetish: all surplus-value (profit,
// interest, rent) is pumped from wage-labour and merely partitioned among the
// three claimants, while wages recover only the variable capital. RevenueStream
// records one branch (source → revenue) with both its fetishised apparent yield
// and the portion of surplus-value it really represents; TrinityFormula binds
// the three streams and the undivided surplus-value behind them; RevenueFetishForm
// contrasts the mystified surface relation against the real one.
package revenue

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// ErrSurplusNotConserved is returned when a trinity's streams do not partition a
// single surplus-value: the apparent-revenue overshoot above newly-created value
// (v + s) must equal the ground-rent — that overshoot is the trinity's one
// fetish (rent appears as income with no value behind it), not extra wealth.
var ErrSurplusNotConserved = errors.New("revenue: trinity does not conserve surplus-value (apparent-revenue overshoot must equal ground-rent)")

// RevenueSourceKind enumerates the three apparent sources of the trinity formula.
type RevenueSourceKind int

const (
	RevenueSourceCapital RevenueSourceKind = 1 // yields interest + profit of enterprise
	RevenueSourceLand    RevenueSourceKind = 2 // yields ground-rent
	RevenueSourceLabour  RevenueSourceKind = 3 // yields wages
)

// IsValid reports whether k is a defined RevenueSourceKind.
func (k RevenueSourceKind) IsValid() bool {
	return k >= RevenueSourceCapital && k <= RevenueSourceLabour
}

// RevenueStreamID identifies a revenue-stream record.
type RevenueStreamID string

// NewRevenueStreamID returns a fresh 96-bit hex ID.
func NewRevenueStreamID() RevenueStreamID {
	return RevenueStreamID(newID())
}

// RevenueStream is one branch of the trinity formula: a source and its revenue.
// ApparentRevenueBP is the fetishised yield as it appears on the surface;
// ActualSourceBP is the portion of surplus-value the revenue really represents
// (zero for wages, which recover variable capital, not surplus-value).
type RevenueStream struct {
	ID                RevenueStreamID   `json:"id"`
	Source            RevenueSourceKind `json:"source"`
	ApparentRevenueBP int64             `json:"apparent_revenue_bp"`
	ActualSourceBP    int64             `json:"actual_source_bp"`
	IsFetishised      bool              `json:"is_fetishised"`
	CreatedAt         time.Time         `json:"created_at"`
}

// TrinityFormulaID identifies a trinity-formula record.
type TrinityFormulaID string

// NewTrinityFormulaID returns a fresh 96-bit hex ID.
func NewTrinityFormulaID() TrinityFormulaID {
	return TrinityFormulaID(newID())
}

// TrinityFormula binds the three revenue streams and the undivided surplus-value
// that actually lies behind them. TotalApparentRevenueBP is the sum of the three
// streams' apparent yields; TotalSurplusValueBP is the sum of the portions of
// surplus-value they really represent — the formula partitions surplus-value, it
// does not create it anew.
type TrinityFormula struct {
	ID                     TrinityFormulaID `json:"id"`
	CapitalStreamID        string           `json:"capital_stream_id"`
	LandStreamID           string           `json:"land_stream_id"`
	LabourStreamID         string           `json:"labour_stream_id"`
	TotalSurplusValueBP    int64            `json:"total_surplus_value_bp"`
	TotalApparentRevenueBP int64            `json:"total_apparent_revenue_bp"`
	CreatedAt              time.Time        `json:"created_at"`
}

// RevenueFetishFormID identifies a revenue-fetish-form record.
type RevenueFetishFormID string

// NewRevenueFetishFormID returns a fresh 96-bit hex ID.
func NewRevenueFetishFormID() RevenueFetishFormID {
	return RevenueFetishFormID(newID())
}

// RevenueFetishForm contrasts the mystified surface relation of a revenue source
// (SurfaceFormula, e.g. "Capital — Interest") with the real underlying relation
// (RealRelation) and names the kind of mystification at work.
type RevenueFetishForm struct {
	ID                RevenueFetishFormID `json:"id"`
	Source            RevenueSourceKind   `json:"source"`
	SurfaceFormula    string              `json:"surface_formula"`
	RealRelation      string              `json:"real_relation"`
	MystificationKind string              `json:"mystification_kind"`
	CreatedAt         time.Time           `json:"created_at"`
}

// SumApparentRevenueBP returns the total apparent (fetishised) revenue across streams.
func SumApparentRevenueBP(streams []RevenueStream) int64 {
	var sum int64
	for _, s := range streams {
		sum += s.ApparentRevenueBP
	}
	return sum
}

// SumActualSourceBP returns the total surplus-value the streams really represent.
func SumActualSourceBP(streams []RevenueStream) int64 {
	var sum int64
	for _, s := range streams {
		sum += s.ActualSourceBP
	}
	return sum
}

// VariableCapitalBP returns the variable capital recovered by the labour streams —
// the wages, which are not surplus-value. Returns 0 when no labour stream is present.
func VariableCapitalBP(streams []RevenueStream) int64 {
	var v int64
	for _, s := range streams {
		if s.Source == RevenueSourceLabour {
			v += s.ApparentRevenueBP
		}
	}
	return v
}

// ProfitBP returns the capital streams' actual share of surplus-value — the
// profit of enterprise plus interest the capitalist retains *after* ground-rent
// is deducted. It is a part of the single social surplus-value, not an addition.
func ProfitBP(streams []RevenueStream) int64 {
	var p int64
	for _, s := range streams {
		if s.Source == RevenueSourceCapital {
			p += s.ActualSourceBP
		}
	}
	return p
}

// RentBP returns the land streams' actual share of surplus-value — the portion
// the landed-property barrier redirects to the landlord. Rent is a deduction
// from the single social surplus-value, not value newly created by the land.
func RentBP(streams []RevenueStream) int64 {
	var r int64
	for _, s := range streams {
		if s.Source == RevenueSourceLand {
			r += s.ActualSourceBP
		}
	}
	return r
}

// NewValueBP returns the newly-created value v + s: the variable capital the
// wages recover plus the undivided surplus-value the streams partition
// (profit + rent). This is all the living labour added — nothing more exists to
// be distributed.
func NewValueBP(streams []RevenueStream) int64 {
	return VariableCapitalBP(streams) + SumActualSourceBP(streams)
}

// FetishOvershootBP returns the amount by which the trinity's apparent revenues
// exceed the value actually created (v + s). Marx's point: this overshoot is not
// extra wealth. It is the ground-rent counted twice — once inside the
// capitalist's average profit and again as the landlord's "independent" income
// from land. Land creates no value; the overshoot is the measure of the fetish.
// A faithful trinity has FetishOvershootBP == RentBP.
func FetishOvershootBP(streams []RevenueStream) int64 {
	return SumApparentRevenueBP(streams) - NewValueBP(streams)
}

// CheckConservation reports whether the streams form a trinity that conserves
// surplus-value: wages carry no surplus-value and the apparent-revenue overshoot
// above v + s equals the ground-rent (profit + rent partition one s). Returns
// ErrSurplusNotConserved otherwise.
func CheckConservation(streams []RevenueStream) error {
	if FetishOvershootBP(streams) != RentBP(streams) {
		return ErrSurplusNotConserved
	}
	return nil
}

// newID returns a fresh 96-bit hex string from crypto/rand.
func newID() string {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("revenue: rand: %w", err))
	}
	return hex.EncodeToString(b[:])
}
