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
	"fmt"
	"time"
)

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

// newID returns a fresh 96-bit hex string from crypto/rand.
func newID() string {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("revenue: rand: %w", err))
	}
	return hex.EncodeToString(b[:])
}
