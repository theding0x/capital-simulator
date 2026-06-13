package profit

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// This file carries Ch. 6's argument: the rate of profit p' = s / (c + v) moves
// with the price of the raw materials that make up the circulating part of the
// constant capital, even though that price-movement adds not a farthing of
// surplus-value. The surplus-value s is wrung from living labour in production
// and is wholly indifferent to what the cotton cost; but it is *measured* against
// the whole capital advanced, and a rise in the price of cotton swells the
// denominator c while leaving s untouched — so the rate of profit falls. A fall
// in the price of cotton shrinks c and the rate of profit rises. The fixed
// capital (machinery, buildings) is not re-priced by the day's cotton market, so
// only the circulating raw-material component fluctuates.
//
// Marx's §I illustration runs both ways from a capital of c = £400, v = £100,
// s = £100 (p' = 20%): when the price of the raw material doubles the circulating
// component swells and the rate of profit falls to 17.24%; when it halves the
// component shrinks and the rate rises to 32.26%. The two movements are taken
// from capitals with different fixed/circulating splits of the same £400 — the
// rate of profit is the more sensitive to price the larger the raw-material share
// of the constant capital.
//
// §II turns to the second face of the same movement: the capital tied up in a
// stock of raw material. When prices fall, the capital that was advanced to hold
// a given stock is *released* — set free for other employment; when they rise,
// further capital must be *tied up* to hold the same stock. Release and tie-up
// are the same magnitude with opposite sign.

// ErrNonPositivePriceFactor is returned when a price factor is not strictly
// positive. The factor is a percentage of the original price (100 = unchanged),
// so a zero or negative factor describes no real price at all.
var ErrNonPositivePriceFactor = errors.New("profit: price factor must be positive")

// priceFactorBase is the scale for ConstantCapitalPriceEffect.PriceFactor: it is
// a percentage of the original raw-material price, so 100 = unchanged, 200 = the
// price has doubled, 50 = the price has halved.
const priceFactorBase = 100

// PriceFluctuationKind names the direction of a price movement in the raw
// material that makes up the circulating constant capital.
type PriceFluctuationKind string

const (
	// KindRise: the raw-material price has risen, swelling c and depressing p'.
	KindRise PriceFluctuationKind = "rise"
	// KindFall: the raw-material price has fallen, shrinking c and lifting p'.
	KindFall PriceFluctuationKind = "fall"
)

// Valid reports whether k is one of the recognised fluctuation directions.
func (k PriceFluctuationKind) Valid() bool {
	switch k {
	case KindRise, KindFall:
		return true
	default:
		return false
	}
}

// RawMaterialPriceFluctuation is a change in the price of one unit of the raw
// material that enters the circulating constant capital. It carries the price
// before and after the movement; the direction and the proportional factor are
// derived from the two.
type RawMaterialPriceFluctuation struct {
	OriginalPricePerUnit LabourMinutes `json:"original_price_per_unit"`
	NewPricePerUnit      LabourMinutes `json:"new_price_per_unit"`
}

// Kind reports whether the price rose or fell. A price that does not move is not
// a fluctuation; Kind treats the no-change case as a fall (the degenerate, p'-
// neutral boundary), and callers that care distinguish it by Factor() == 100.
func (f RawMaterialPriceFluctuation) Kind() PriceFluctuationKind {
	if f.NewPricePerUnit > f.OriginalPricePerUnit {
		return KindRise
	}
	return KindFall
}

// Factor returns the new price as a percentage of the original (100 = unchanged,
// 200 = doubled, 50 = halved). With no original price there is no ratio to form,
// so it returns 0.
func (f RawMaterialPriceFluctuation) Factor() int64 {
	if f.OriginalPricePerUnit <= 0 {
		return 0
	}
	return int64(f.NewPricePerUnit) * priceFactorBase / int64(f.OriginalPricePerUnit)
}

// ConstantCapitalPriceEffect carries how a re-pricing of the raw-material element
// of the constant capital ripples into the rate of profit. The constant capital
// is split into a FixedCapital component the price movement does not touch
// (machinery, buildings) and an OriginalMaterialCapital component — the
// circulating outlay on raw material — that scales with the price. PriceFactor is
// that price as a percentage of its original level (100 = unchanged). The
// variable capital and surplus-value are held fixed, so the whole movement of the
// rate of profit comes from the re-pricing of the material.
type ConstantCapitalPriceEffect struct {
	FixedCapital            ConstantCapital `json:"fixed_capital"`
	OriginalMaterialCapital ConstantCapital `json:"original_material_capital"`
	PriceFactor             int64           `json:"price_factor"`
	VariableCapital         VariableCapital `json:"variable_capital"`
	SurplusValue            SurplusValue    `json:"surplus_value"`
}

// Validate enforces non-negative magnitudes and a strictly positive price factor.
func (e ConstantCapitalPriceEffect) Validate() error {
	if e.FixedCapital < 0 || e.OriginalMaterialCapital < 0 || e.VariableCapital < 0 || e.SurplusValue < 0 {
		return ErrNegativeMagnitude
	}
	if e.PriceFactor <= 0 {
		return ErrNonPositivePriceFactor
	}
	return nil
}

// OriginalConstantCapital is c before the price movement: the fixed component
// plus the raw-material component at its original price.
func (e ConstantCapitalPriceEffect) OriginalConstantCapital() ConstantCapital {
	return e.FixedCapital + e.OriginalMaterialCapital
}

// NewMaterialCapital is the circulating raw-material outlay after the price
// movement: the original material capital scaled by the price factor.
func (e ConstantCapitalPriceEffect) NewMaterialCapital() ConstantCapital {
	return e.OriginalMaterialCapital * ConstantCapital(e.PriceFactor) / priceFactorBase
}

// NewConstantCapital is c after the price movement: the untouched fixed component
// plus the re-priced raw-material component. Only the circulating raw-material
// part changes.
func (e ConstantCapitalPriceEffect) NewConstantCapital() ConstantCapital {
	return e.FixedCapital + e.NewMaterialCapital()
}

// OldProfitRate is p' = s / (c + v) before the price movement, in basis points.
func (e ConstantCapitalPriceEffect) OldProfitRate() int64 {
	return profitRateBP(e.SurplusValue, e.OriginalConstantCapital(), e.VariableCapital)
}

// NewProfitRate is p' = s / (c' + v) after the price movement, in basis points,
// where c' carries the re-priced raw material. With s and v unchanged, a price
// rise lowers it and a price fall raises it.
func (e ConstantCapitalPriceEffect) NewProfitRate() int64 {
	return profitRateBP(e.SurplusValue, e.NewConstantCapital(), e.VariableCapital)
}

// Kind reports the direction of the price movement: a rise when the factor lifts
// the price above its original level, a fall otherwise.
func (e ConstantCapitalPriceEffect) Kind() PriceFluctuationKind {
	if e.PriceFactor > priceFactorBase {
		return KindRise
	}
	return KindFall
}

// PriceFluctuationAnalysisID identifies a stored price-fluctuation analysis.
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type PriceFluctuationAnalysisID string

// NewPriceFluctuationAnalysisID returns a fresh random PriceFluctuationAnalysisID.
func NewPriceFluctuationAnalysisID() PriceFluctuationAnalysisID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; if it does, panic is appropriate.
		panic(fmt.Errorf("profit: rand: %w", err))
	}
	return PriceFluctuationAnalysisID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty PriceFluctuationAnalysisID.
func (id PriceFluctuationAnalysisID) IsZero() bool { return id == "" }

// PriceFluctuationAnalysis records one re-pricing of the raw-material element of
// constant capital and the movement it produces in the rate of profit. It stores
// the effect (the capital split, the price factor, v and s) alongside the
// direction of the movement and the old and new rates of profit, both in basis
// points (100% = 10000 bp).
type PriceFluctuationAnalysis struct {
	ID            PriceFluctuationAnalysisID `json:"id"`
	Kind          PriceFluctuationKind       `json:"kind"`
	Effect        ConstantCapitalPriceEffect `json:"effect"`
	OldProfitRate int64                      `json:"old_profit_rate"`
	NewProfitRate int64                      `json:"new_profit_rate"`
	CreatedAt     time.Time                  `json:"created_at"`
}

// Validate defers to the underlying effect's invariants.
func (a PriceFluctuationAnalysis) Validate() error {
	return a.Effect.Validate()
}

// ComputePriceFluctuationAnalysis derives the direction and the old and new rates
// of profit from a price effect on the constant capital. It is a pure function:
// it assigns no ID or timestamp — the store does that on persistence.
func ComputePriceFluctuationAnalysis(e ConstantCapitalPriceEffect) PriceFluctuationAnalysis {
	return PriceFluctuationAnalysis{
		Kind:          e.Kind(),
		Effect:        e,
		OldProfitRate: e.OldProfitRate(),
		NewProfitRate: e.NewProfitRate(),
	}
}

// StockCapitalRelease models §II's release of capital: when the price of the raw
// material falls, the capital that was advanced to hold a stock of it is set
// free. StockUnits is the quantity consumed per month and StockMonths the number
// of months' stock held, so the stock is StockMonths × StockUnits units; the
// capital released is that quantity valued at the fall in price per unit.
type StockCapitalRelease struct {
	StockMonths          int64         `json:"stock_months"`
	StockUnits           int64         `json:"stock_units"`
	OriginalPricePerUnit LabourMinutes `json:"original_price_per_unit"`
	NewPricePerUnit      LabourMinutes `json:"new_price_per_unit"`
}

// Released is the capital set free by the price fall: the whole stock valued at
// the difference between the old and the new price. It is positive when the price
// falls (capital freed), negative when it rises (capital tied up), and zero when
// the price holds.
func (r StockCapitalRelease) Released() LabourMinutes {
	stock := r.StockMonths * r.StockUnits
	return LabourMinutes(stock) * (r.OriginalPricePerUnit - r.NewPricePerUnit)
}

// CapitalTieUp models the converse of §II: when the price of the raw material
// rises, further capital must be advanced to hold the same stock. The fields
// mean the same as StockCapitalRelease.
type CapitalTieUp struct {
	StockMonths          int64         `json:"stock_months"`
	StockUnits           int64         `json:"stock_units"`
	OriginalPricePerUnit LabourMinutes `json:"original_price_per_unit"`
	NewPricePerUnit      LabourMinutes `json:"new_price_per_unit"`
}

// TiedUp is the additional capital absorbed by the price rise: the whole stock
// valued at the difference between the new and the old price. For a given price
// movement it is the exact negative of StockCapitalRelease.Released — release and
// tie-up are one magnitude with opposite sign.
func (t CapitalTieUp) TiedUp() LabourMinutes {
	stock := t.StockMonths * t.StockUnits
	return LabourMinutes(stock) * (t.NewPricePerUnit - t.OriginalPricePerUnit)
}

// SupplyAction names the capitalist's response to a raw-material price signal.
type SupplyAction string

const (
	// ActionExpandOutput: a fall in raw-material prices cheapens production and
	// invites the capitalist to expand output and build up stock.
	ActionExpandOutput SupplyAction = "expand_output"
	// ActionContractOutput: a rise in prices makes production dearer and presses
	// the capitalist to contract output and run stock down.
	ActionContractOutput SupplyAction = "contract_output"
)

// SupplyRegulation models §II's "experiments in corpore vili": the price of the
// raw material acts as a signal that regulates the scale of production. A fall
// invites expansion, a rise forces contraction — the blind regulation of supply
// by price that, Marx notes, is irreconcilable with any conscious control of
// production.
type SupplyRegulation struct {
	Fluctuation RawMaterialPriceFluctuation `json:"fluctuation"`
}

// Action returns the response the price signal invites: expand output when the
// price has fallen, contract it when the price has risen.
func (s SupplyRegulation) Action() SupplyAction {
	if s.Fluctuation.Kind() == KindFall {
		return ActionExpandOutput
	}
	return ActionContractOutput
}
