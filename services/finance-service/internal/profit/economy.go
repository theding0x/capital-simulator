package profit

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// This file carries Ch. 5's argument: a saving in constant capital raises the
// rate of profit without adding a farthing of surplus-value. The rate of profit
// p' = s / (c + v) carries the constant capital in its denominator alone — leave
// s and v untouched and every economy in c shrinks the denominator and lifts the
// rate. Marx surveys the forms this economy takes: the same fixed capital made to
// serve a longer working-day (its charge spread over more hours of surplus-
// labour); the conditions of production — buildings, power, heating, lighting —
// used in common and amortised across many labourers; and the waste of raw
// material reclaimed and returned to production. None of them touches the
// surplus-value; all of them cheapen the constant capital against which that
// surplus is measured, and so raise the rate of profit.
//
// Two enterprises make the point: with the same s = £1,000 on the same v =
// £1,000, the one carrying c = £9,000 shows p' = 10%, the one carrying c =
// £11,000 only 8 1/3%. The £2,000 of constant capital economised is the whole
// difference between them.

// EconomyKind names the form a saving in constant capital takes. Marx runs
// through them in Ch. 5: the working-day lengthened over the same fixed capital,
// the conditions of production shared in common, waste reclaimed, machinery
// improved, and raw material used more sparingly.
type EconomyKind string

const (
	// KindSharedConditions: buildings, power, heating and lighting used
	// communally and amortised across many labourers or enterprises (§III).
	KindSharedConditions EconomyKind = "shared_conditions"
	// KindWasteReduction: the excretions of production reclaimed and returned to
	// production, cheapening the raw material (§IV).
	KindWasteReduction EconomyKind = "waste_reduction"
	// KindWorkdayExtension: the same fixed capital made to serve a longer
	// working-day, its charge spread over more hours of surplus-labour (§I).
	KindWorkdayExtension EconomyKind = "workday_extension"
	// KindImprovedMachinery: improvements that let existing machinery work more
	// cheaply and effectively, reducing the wear charged to the product (§I).
	KindImprovedMachinery EconomyKind = "improved_machinery"
	// KindRawMaterialSaving: a smaller outlay on raw and auxiliary materials for
	// the same product, through better material or less spoilage (§I).
	KindRawMaterialSaving EconomyKind = "raw_material_saving"
)

// Valid reports whether k is one of the recognised economy kinds.
func (k EconomyKind) Valid() bool {
	switch k {
	case KindSharedConditions, KindWasteReduction, KindWorkdayExtension,
		KindImprovedMachinery, KindRawMaterialSaving:
		return true
	default:
		return false
	}
}

// ConstantCapitalEconomy is a saving in constant capital that raises the rate of
// profit without changing surplus-value or variable capital. It carries the
// capital the saving applies to — the original constant capital c, the variable
// capital v and the surplus-value s — together with the economy realised in c.
// Because s and v are untouched and the denominator c + v shrinks by the saving,
// the rate of profit s / (c + v) rises: the whole burden of the chapter.
type ConstantCapitalEconomy struct {
	Kind            EconomyKind     `json:"kind"`
	ConstantCapital ConstantCapital `json:"constant_capital"`
	VariableCapital VariableCapital `json:"variable_capital"`
	SurplusValue    SurplusValue    `json:"surplus_value"`
	Saving          ConstantCapital `json:"saving"`
}

// Validate enforces the magnitude invariants: non-negative inputs, a recognised
// economy kind, and a saving that cannot exceed the constant capital it economises.
func (e ConstantCapitalEconomy) Validate() error {
	if e.ConstantCapital < 0 || e.VariableCapital < 0 || e.SurplusValue < 0 || e.Saving < 0 {
		return ErrNegativeMagnitude
	}
	if !e.Kind.Valid() {
		return ErrUnknownEconomyKind
	}
	if e.Saving > e.ConstantCapital {
		return ErrSavingExceedsConstant
	}
	return nil
}

// OldProfitRate is the rate of profit before the economy, s / (c + v), in basis
// points (100% = 10000 bp).
func (e ConstantCapitalEconomy) OldProfitRate() int64 {
	return profitRateBP(e.SurplusValue, e.ConstantCapital, e.VariableCapital)
}

// NewProfitRate is the rate of profit after the economy, s / (c − saving + v), in
// basis points. With the same s and v over a smaller constant capital it is the
// higher rate.
func (e ConstantCapitalEconomy) NewProfitRate() int64 {
	return profitRateBP(e.SurplusValue, e.ConstantCapital-e.Saving, e.VariableCapital)
}

// ProfitRateGain is the rise in the rate of profit the economy produces, in basis
// points: NewProfitRate − OldProfitRate. It is positive whenever the saving is
// positive (s and v unchanged) and zero when nothing is economised.
func (e ConstantCapitalEconomy) ProfitRateGain() int64 {
	return e.NewProfitRate() - e.OldProfitRate()
}

// profitRateBP forms s / (c + v) in basis points, guarding a zero denominator:
// with nothing advanced there is no rate of self-expansion to express.
func profitRateBP(s SurplusValue, c ConstantCapital, v VariableCapital) int64 {
	total := int64(c) + int64(v)
	if total <= 0 {
		return 0
	}
	return int64(s) * basisPoints / total
}

// SharedCondition models §III's economy from conditions of production used in
// common — a boiler-house, building, lighting or heating that many enterprises or
// labourers share instead of each providing its own. SoloCost is the aggregate
// the participants bear if each provides the condition for itself; SharedCost is
// the cost of the single shared facility serving them all. The economy is the
// difference between the two, spread across the participants.
type SharedCondition struct {
	Kind             EconomyKind     `json:"kind"`
	SoloCost         ConstantCapital `json:"solo_cost"`
	SharedCost       ConstantCapital `json:"shared_cost"`
	ParticipantCount int64           `json:"participant_count"`
}

// SavingPerParticipant returns the constant-capital economy each participant
// realises by sharing the condition rather than providing it alone:
// (SoloCost − SharedCost) / ParticipantCount. With no participants it is zero.
func (s SharedCondition) SavingPerParticipant() ConstantCapital {
	if s.ParticipantCount <= 0 {
		return 0
	}
	return (s.SoloCost - s.SharedCost) / ConstantCapital(s.ParticipantCount)
}

// WasteReduction models §IV's economy from reclaiming the excretions of
// production — waste cotton in spinning, for instance — and returning them to
// production. OriginalWaste is the quantity lost as waste before the economy;
// RecoveredWaste is the quantity reclaimed; PricePerUnit is the value of one unit
// of the raw material. Reclaiming waste cheapens the constant capital by the value
// of the material no longer bought afresh.
type WasteReduction struct {
	Kind           EconomyKind   `json:"kind"`
	OriginalWaste  int64         `json:"original_waste"`
	RecoveredWaste int64         `json:"recovered_waste"`
	PricePerUnit   LabourMinutes `json:"price_per_unit"`
}

// ConstantCapitalSaving returns the reduction in constant capital from reclaiming
// waste: RecoveredWaste × PricePerUnit.
func (w WasteReduction) ConstantCapitalSaving() ConstantCapital {
	return ConstantCapital(w.RecoveredWaste * int64(w.PricePerUnit))
}

// WorkdayExtensionEffect models §I's economy from lengthening the working-day: the
// same fixed capital serves the longer day without fresh outlay, so the
// fixed-capital charge each hour of labour bears falls. FixedCapital is the value
// of the instruments of labour; OriginalHours and ExtendedHours are the lengths of
// the working-day before and after the extension.
type WorkdayExtensionEffect struct {
	Kind          EconomyKind     `json:"kind"`
	FixedCapital  ConstantCapital `json:"fixed_capital"`
	OriginalHours int64           `json:"original_hours"`
	ExtendedHours int64           `json:"extended_hours"`
}

// FixedChargePerHour returns the share of the fixed capital borne by each hour of
// a working-day of the given length: FixedCapital / hours. With no hours it is
// zero.
func (e WorkdayExtensionEffect) FixedChargePerHour(hours int64) ConstantCapital {
	if hours <= 0 {
		return 0
	}
	return e.FixedCapital / ConstantCapital(hours)
}

// ProfitRateGain returns the economy in the per-hour fixed-capital charge won by
// lengthening the working-day from OriginalHours to ExtendedHours:
// FixedCapital/OriginalHours − FixedCapital/ExtendedHours. The same fixed capital
// serves the longer day for no fresh outlay, so each hour of labour bears a
// smaller charge; surplus-value unchanged, the rate of profit rises. It is zero
// when there is no fixed capital to spread (or the day is not lengthened) and
// positive otherwise.
func (e WorkdayExtensionEffect) ProfitRateGain() ConstantCapital {
	return e.FixedChargePerHour(e.OriginalHours) - e.FixedChargePerHour(e.ExtendedHours)
}

// EconomyAnalysisID identifies a stored economy analysis. 96-bit hex from
// crypto/rand, like every other domain ID in the simulator.
type EconomyAnalysisID string

// NewEconomyAnalysisID returns a fresh random EconomyAnalysisID.
func NewEconomyAnalysisID() EconomyAnalysisID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; if it does, panic is appropriate.
		panic(fmt.Errorf("profit: rand: %w", err))
	}
	return EconomyAnalysisID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty EconomyAnalysisID.
func (id EconomyAnalysisID) IsZero() bool { return id == "" }

// EconomyAnalysis records one economy in constant capital and the rise in the
// rate of profit it produces. It stores the economy itself (kind, the capital it
// applies to, and the saving) alongside the old rate of profit s/(c+v), the new
// rate s/(c−saving+v), and the gain between them — all in basis points.
type EconomyAnalysis struct {
	ID             EconomyAnalysisID      `json:"id"`
	Economy        ConstantCapitalEconomy `json:"economy"`
	OldProfitRate  int64                  `json:"old_profit_rate"`
	NewProfitRate  int64                  `json:"new_profit_rate"`
	ProfitRateGain int64                  `json:"profit_rate_gain"`
	CreatedAt      time.Time              `json:"created_at"`
}

// Validate defers to the underlying economy's invariants.
func (a EconomyAnalysis) Validate() error {
	return a.Economy.Validate()
}

// ComputeEconomyAnalysis derives the old and new rates of profit and the gain
// between them from an economy in constant capital. It is a pure function: it
// assigns no ID or timestamp — the store does that on persistence.
func ComputeEconomyAnalysis(e ConstantCapitalEconomy) EconomyAnalysis {
	return EconomyAnalysis{
		Economy:        e,
		OldProfitRate:  e.OldProfitRate(),
		NewProfitRate:  e.NewProfitRate(),
		ProfitRateGain: e.ProfitRateGain(),
	}
}
