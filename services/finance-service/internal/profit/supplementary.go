package profit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// This file carries Ch. 7's supplementary remarks, which close Part I by drawing
// together the determinants of the rate of profit and refuting Rodbertus's error.
//
// Two threads run through the chapter:
//
//   - The organic composition channel. With the same surplus-value and the same
//     variable capital, two capitals of different organic composition show
//     different rates of profit: the one carrying more constant capital shows the
//     lower rate. Enterprise A (s = £1,000, v = £1,000, c = £9,000) shows
//     p' = 10%; Enterprise B (same s and v, c = £11,000) only 8 1/3%. The whole
//     difference is the composition.
//
//   - Rodbertus's error. A mere change in the money-value of capital — gold
//     falling by half, say, so that the same capital is reckoned at twice the
//     figure — changes nothing real: capital and profit are re-valued in the same
//     proportion, and the rate of profit is unmoved. So too a proportional change
//     in the magnitude of capital at an unchanged composition. Only a change that
//     alters the organic composition (or the rate of surplus-value) moves the rate
//     of profit. Rodbertus mistook the first, nominal kind of change for the last.
//
// The profit-rate's surface appearance — that profit springs from the capitalist's
// own enterprise and from the total capital advanced, not from the unpaid labour
// contained in the variable capital alone — is the mystification Part I has been
// unwinding; ProfitSurfaceAppearance records it explicitly. PartISummary gathers
// the rate-of-profit analyses of Chs. 1–6 and names the determinants the chapter
// has established.

// ErrUnknownMagnitudeKind is returned for a capital-magnitude change whose kind is
// not one of the three recognised forms.
var ErrUnknownMagnitudeKind = errors.New("profit: unknown magnitude change kind")

// absInt64 returns the absolute value of n. The rate-of-profit difference between
// two enterprises is reported as a magnitude, irrespective of which is the larger.
func absInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

// OrganicCompositionEffect sets two capitals side by side that carry the same
// surplus-value on the same variable capital but differ in their constant
// capital — and so in their organic composition. It is the demonstration that the
// composition alone, with s and v held equal, sets the rate of profit: the capital
// with the heavier constant component shows the lower rate.
type OrganicCompositionEffect struct {
	SA SurplusValue    `json:"s_a"`
	VA VariableCapital `json:"v_a"`
	CA ConstantCapital `json:"c_a"`
	SB SurplusValue    `json:"s_b"`
	VB VariableCapital `json:"v_b"`
	CB ConstantCapital `json:"c_b"`
}

// Validate enforces non-negative magnitudes on both enterprises.
func (e OrganicCompositionEffect) Validate() error {
	if e.SA < 0 || e.VA < 0 || e.CA < 0 || e.SB < 0 || e.VB < 0 || e.CB < 0 {
		return ErrNegativeMagnitude
	}
	return nil
}

// ProfitRateA is enterprise A's rate of profit, s/(c+v), in basis points.
func (e OrganicCompositionEffect) ProfitRateA() int64 {
	return profitRateBP(e.SA, e.CA, e.VA)
}

// ProfitRateB is enterprise B's rate of profit, s/(c+v), in basis points.
func (e OrganicCompositionEffect) ProfitRateB() int64 {
	return profitRateBP(e.SB, e.CB, e.VB)
}

// RateDifference is the gap between the two rates of profit, in basis points,
// reported as a magnitude. With equal s and v it measures the effect of the
// difference in organic composition alone.
func (e OrganicCompositionEffect) RateDifference() int64 {
	return absInt64(e.ProfitRateA() - e.ProfitRateB())
}

// CompositionEffectID identifies a stored organic-composition comparison.
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type CompositionEffectID string

// NewCompositionEffectID returns a fresh random CompositionEffectID.
func NewCompositionEffectID() CompositionEffectID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; if it does, panic is appropriate.
		panic(fmt.Errorf("profit: rand: %w", err))
	}
	return CompositionEffectID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CompositionEffectID.
func (id CompositionEffectID) IsZero() bool { return id == "" }

// CompositionEffectAnalysis records one organic-composition comparison and the two
// rates of profit it gives rise to, along with the gap between them, all in basis
// points.
type CompositionEffectAnalysis struct {
	ID             CompositionEffectID      `json:"id"`
	Effect         OrganicCompositionEffect `json:"effect"`
	ProfitRateA    int64                    `json:"profit_rate_a"`
	ProfitRateB    int64                    `json:"profit_rate_b"`
	RateDifference int64                    `json:"rate_difference"`
	CreatedAt      time.Time                `json:"created_at"`
}

// Validate defers to the underlying effect's invariants.
func (a CompositionEffectAnalysis) Validate() error {
	return a.Effect.Validate()
}

// ComputeCompositionEffectAnalysis derives the two rates of profit and the gap
// between them from a composition comparison. It is a pure function: it assigns no
// ID or timestamp — the store does that on persistence.
func ComputeCompositionEffectAnalysis(e OrganicCompositionEffect) CompositionEffectAnalysis {
	return CompositionEffectAnalysis{
		Effect:         e,
		ProfitRateA:    e.ProfitRateA(),
		ProfitRateB:    e.ProfitRateB(),
		RateDifference: e.RateDifference(),
	}
}

// MagnitudeChangeKind names the form a change in the magnitude of capital takes.
// Marx runs through them in Ch. 7 to refute Rodbertus.
type MagnitudeChangeKind string

const (
	// KindMoneyValueChange: a change in the money-value of capital only (gold
	// rises or falls), re-valuing capital and profit in the same proportion. The
	// rate of profit is unmoved — the change is purely nominal.
	KindMoneyValueChange MagnitudeChangeKind = "money_value_change"
	// KindProportionalChange: the magnitude of capital changes but the organic
	// composition (and the rate of surplus-value) is unchanged, so profit scales
	// with capital and the rate of profit is unmoved.
	KindProportionalChange MagnitudeChangeKind = "proportional_change"
	// KindOrganicChange: the change alters the organic composition, so profit does
	// not scale with capital and the rate of profit moves. This is the only one of
	// the three that touches the rate.
	KindOrganicChange MagnitudeChangeKind = "organic_change"
)

// Valid reports whether k is one of the three recognised magnitude-change kinds.
func (k MagnitudeChangeKind) Valid() bool {
	switch k {
	case KindMoneyValueChange, KindProportionalChange, KindOrganicChange:
		return true
	default:
		return false
	}
}

// magnitudeFactorBase is the scale for CapitalMagnitudeChange.Factor: it is a
// percentage of the original magnitude, so 100 = unchanged, 200 = doubled.
const magnitudeFactorBase = 100

// CapitalMagnitudeChange carries a change in the absolute magnitude of a capital
// and asks what becomes of its rate of profit. OriginalCapital is the capital C
// and OriginalProfit the profit p it yielded; Factor is the new magnitude as a
// percentage of the old (100 = unchanged, 200 = doubled). Whether the rate moves
// depends entirely on the Kind: a nominal re-valuation or a proportional change
// at constant composition leaves it untouched; only a change in the organic
// composition moves it.
type CapitalMagnitudeChange struct {
	Kind            MagnitudeChangeKind `json:"kind"`
	OriginalCapital TotalCapital        `json:"original_capital"`
	OriginalProfit  SurplusValue        `json:"original_profit"`
	Factor          int64               `json:"factor"`
}

// Validate enforces non-negative magnitudes, a strictly positive factor, and a
// recognised kind.
func (m CapitalMagnitudeChange) Validate() error {
	if m.OriginalCapital < 0 || m.OriginalProfit < 0 {
		return ErrNegativeMagnitude
	}
	if m.Factor <= 0 {
		return ErrNonPositivePriceFactor
	}
	if !m.Kind.Valid() {
		return ErrUnknownMagnitudeKind
	}
	return nil
}

// OldRate is the rate of profit before the change, p / C, in basis points.
func (m CapitalMagnitudeChange) OldRate() int64 {
	return magnitudeRateBP(int64(m.OriginalProfit), int64(m.OriginalCapital))
}

// NewRate is the rate of profit after the change, in basis points. For a nominal
// money-value change or a proportional change at constant composition, capital and
// profit move together and the rate is unchanged. For an organic change, only the
// capital is scaled (the composition shifts), so the rate moves.
func (m CapitalMagnitudeChange) NewRate() int64 {
	if m.Kind == KindOrganicChange {
		newCapital := int64(m.OriginalCapital) * m.Factor / magnitudeFactorBase
		return magnitudeRateBP(int64(m.OriginalProfit), newCapital)
	}
	return m.OldRate()
}

// RateUnchanged reports whether the change leaves the rate of profit untouched —
// true for the two nominal kinds, false in general for an organic change. This is
// the kernel of the refutation of Rodbertus.
func (m CapitalMagnitudeChange) RateUnchanged() bool {
	return m.NewRate() == m.OldRate()
}

// magnitudeRateBP forms profit / capital in basis points, rounded half-up (the
// Vol. III house convention), guarding a zero denominator: with nothing advanced
// there is no rate to express.
func magnitudeRateBP(profit, capital int64) int64 {
	if capital <= 0 {
		return 0
	}
	return roundHalfUp(profit*basisPoints, capital)
}

// MagnitudeChangeID identifies a stored capital-magnitude change.
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type MagnitudeChangeID string

// NewMagnitudeChangeID returns a fresh random MagnitudeChangeID.
func NewMagnitudeChangeID() MagnitudeChangeID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; if it does, panic is appropriate.
		panic(fmt.Errorf("profit: rand: %w", err))
	}
	return MagnitudeChangeID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty MagnitudeChangeID.
func (id MagnitudeChangeID) IsZero() bool { return id == "" }

// MagnitudeChangeAnalysis records one change in the magnitude of capital alongside
// the old and new rates of profit (in basis points) and whether the rate was left
// unchanged.
type MagnitudeChangeAnalysis struct {
	ID            MagnitudeChangeID      `json:"id"`
	Change        CapitalMagnitudeChange `json:"change"`
	OldRate       int64                  `json:"old_rate"`
	NewRate       int64                  `json:"new_rate"`
	RateUnchanged bool                   `json:"rate_unchanged"`
	CreatedAt     time.Time              `json:"created_at"`
}

// Validate defers to the underlying change's invariants.
func (a MagnitudeChangeAnalysis) Validate() error {
	return a.Change.Validate()
}

// ComputeMagnitudeChangeAnalysis derives the old and new rates and the unchanged
// flag from a magnitude change. It is a pure function: it assigns no ID or
// timestamp — the store does that on persistence.
func ComputeMagnitudeChangeAnalysis(m CapitalMagnitudeChange) MagnitudeChangeAnalysis {
	return MagnitudeChangeAnalysis{
		Change:        m,
		OldRate:       m.OldRate(),
		NewRate:       m.NewRate(),
		RateUnchanged: m.RateUnchanged(),
	}
}

// ProfitSurfaceAppearance records the inversion at the heart of the profit-form:
// to the capitalist the profit appears to spring from the total capital advanced
// and from his own enterprise, when its real source is the unpaid labour contained
// in the variable capital alone. It carries the rate of profit p' the capitalist
// sees and the rate of surplus-value s' that the profit-form conceals.
type ProfitSurfaceAppearance struct {
	ProfitRate       int64  `json:"profit_rate"`
	SurplusValueRate int64  `json:"surplus_value_rate"`
	ApparentSource   string `json:"apparent_source"`
	RealSource       string `json:"real_source"`
}

// NewProfitSurfaceAppearance records how the rate of profit p' presents itself
// against the rate of surplus-value s' it conceals.
func NewProfitSurfaceAppearance(profitRate, surplusValueRate int64) ProfitSurfaceAppearance {
	return ProfitSurfaceAppearance{
		ProfitRate:       profitRate,
		SurplusValueRate: surplusValueRate,
		ApparentSource:   "the total capital advanced and the capitalist's own enterprise",
		RealSource:       "the unpaid labour contained in the variable capital alone",
	}
}

// PartIDeterminants names the channels through which the rate of profit has been
// shown to vary across Part I (Chs. 1–6). The summary exposes the complete set so
// the determinants established chapter by chapter can be read at a glance.
func PartIDeterminants() []string {
	return []string{
		"rate of surplus-value s' (Chs. 2-3)",
		"organic composition of capital c:v (Chs. 3, 7)",
		"number of turnovers n (Ch. 4)",
		"economy in the employment of constant capital (Ch. 5)",
		"price fluctuation of raw materials (Ch. 6)",
	}
}

// PartISummaryID identifies a Part I summary snapshot.
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type PartISummaryID string

// NewPartISummaryID returns a fresh random PartISummaryID.
func NewPartISummaryID() PartISummaryID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; if it does, panic is appropriate.
		panic(fmt.Errorf("profit: rand: %w", err))
	}
	return PartISummaryID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty PartISummaryID.
func (id PartISummaryID) IsZero() bool { return id == "" }

// PartISummary gathers the rate-of-profit analyses accumulated across Part I and
// names the determinants the chapters have established. It is a read-only snapshot
// assembled on demand from the rate-of-profit store, not a persisted entity.
type PartISummary struct {
	ID           PartISummaryID       `json:"id"`
	Determinants []string             `json:"determinants"`
	Analyses     []ProfitRateAnalysis `json:"analyses"`
	CreatedAt    time.Time            `json:"created_at"`
}

// ComputePartISummary assembles a Part I summary from a set of rate-of-profit
// analyses, naming the determinants of the rate of profit established across the
// part. The analyses slice is normalised to non-nil so the summary always carries
// the determinants even when no analyses have been recorded.
func ComputePartISummary(analyses []ProfitRateAnalysis) PartISummary {
	if analyses == nil {
		analyses = []ProfitRateAnalysis{}
	}
	return PartISummary{
		ID:           NewPartISummaryID(),
		Determinants: PartIDeterminants(),
		Analyses:     analyses,
	}
}
