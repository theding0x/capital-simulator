// Package merchant implements Capital Vol. III, Part IV, Ch. 16
// ("Commercial Capital"). Part IV opens here.
//
// Commercial capital (merchant's capital) is a form of the metamorphosis of
// commodity-capital and money-capital — M—C—M′ — that is separated off from
// the industrial circuit and becomes an independent mode of capital. Its
// specific function is the realisation of surplus-value (converting C′ into
// M′); it does not produce surplus-value of its own. The four functions
// distinguished here correspond to Marx's §1–§4 in the chapter:
//
//	§1 Realisation — the merchant buys C at its value and sells at value+profit.
//	§2 Storage     — the merchant holds stock until market conditions improve.
//	§3 Transport   — value is added by carriage, but only use-value to the buyer;
//	                 no new surplus-value from the merchant's standpoint.
//	§4 Distribution — breaking bulk, packaging, etc. — again no new s.
//
// The central invariant across all four functions is that
// SurplusValueProduced == 0: commercial capital is parasitic on the
// surplus-value produced in the industrial circuit.
package merchant

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// Sentinel errors for invariant violations in Ch. 16.
var (
	ErrNonPositiveMoney     = errors.New("merchant: money_advanced must be positive")
	ErrNoProfit             = errors.New("merchant: money_returned must exceed money_advanced")
	ErrNonIndustrialSource  = errors.New("merchant: source_capital must be industrial")
	ErrNegativeMagnitude    = errors.New("merchant: negative magnitude")
)

// CommercialCapitalFunction enumerates the four functions of commercial capital
// distinguished in Vol. III Ch. 16. None produces surplus-value.
type CommercialCapitalFunction int

const (
	// FunctionRealisation — §1: the merchant converts C′ into M′ for the
	// industrial capitalist, realising the surplus-value already produced.
	FunctionRealisation CommercialCapitalFunction = 0
	// FunctionStorage — §2: the merchant holds stock of commodities, bearing
	// the costs of warehousing while awaiting a favourable price.
	FunctionStorage CommercialCapitalFunction = 1
	// FunctionTransport — §3: the merchant moves commodities from the site of
	// production to the site of consumption. Transport adds use-value to the
	// buyer but creates no surplus-value for the merchant.
	FunctionTransport CommercialCapitalFunction = 2
	// FunctionDistribution — §4: breaking bulk, re-packaging and allied
	// activities. Again, no surplus-value is generated.
	FunctionDistribution CommercialCapitalFunction = 3
)

// CommercialCapitalID identifies a persisted commercial-capital record.
// 96-bit hex from crypto/rand, like every other domain ID in the simulator.
type CommercialCapitalID string

// NewCommercialCapitalID returns a fresh random CommercialCapitalID.
func NewCommercialCapitalID() CommercialCapitalID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(fmt.Errorf("merchant: rand: %w", err))
	}
	return CommercialCapitalID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty CommercialCapitalID.
func (id CommercialCapitalID) IsZero() bool { return id == "" }

// CommercialCapital is the persisted entity representing a merchant capital
// advance (Vol. III Ch. 16). The central invariant is SurplusValueProduced == 0:
// commercial capital appropriates a share of the surplus-value produced in the
// industrial circuit — it does not produce any itself.
type CommercialCapital struct {
	ID                   CommercialCapitalID       `json:"id"`
	MoneyAdvanced        int64                     `json:"money_advanced"`        // pence; >0
	CommodityDescription string                    `json:"commodity_description"` // textual description of the commodity handled
	Function             CommercialCapitalFunction `json:"function"`              // realisation/storage/transport/distribution
	SurplusValueProduced int64                     `json:"surplus_value_produced"` // always 0 — invariant of commercial capital
	CreatedAt            time.Time                 `json:"created_at"`
}

// NewCommercialCapital constructs a CommercialCapital from the money advanced, a
// commodity description, and the merchant's function. It returns ErrNonPositiveMoney
// if moneyAdvanced ≤ 0. SurplusValueProduced is always 0. No ID or timestamp is
// assigned — the store does that on persistence.
func NewCommercialCapital(moneyAdvanced int64, commodityDescription string, fn CommercialCapitalFunction) (CommercialCapital, error) {
	if moneyAdvanced <= 0 {
		return CommercialCapital{}, ErrNonPositiveMoney
	}
	return CommercialCapital{
		MoneyAdvanced:        moneyAdvanced,
		CommodityDescription: commodityDescription,
		Function:             fn,
		SurplusValueProduced: 0,
	}, nil
}

// MerchantCircuit is the value object for the merchant's M—C—M′ circuit (§1).
// MoneyAdvanced is the capital laid out; CommodityValue is the value of the
// commodity bought; MoneyReturned is the price realised on sale. The circuit
// is valid only when MoneyReturned > MoneyAdvanced — i.e. the merchant extracts
// a profit from the realisation of surplus-value already produced in the
// industrial circuit. Computed on demand, not persisted.
type MerchantCircuit struct {
	MoneyAdvanced   int64 `json:"money_advanced"`   // pence; >0
	CommodityValue  int64 `json:"commodity_value"`  // pence; value of the commodity
	MoneyReturned   int64 `json:"money_returned"`   // pence; price realised on sale
}

// Profit is the merchant's gain: MoneyReturned − MoneyAdvanced. A positive
// value represents a claim on the industrial surplus-value; it is never
// produced anew by the merchant.
func (c MerchantCircuit) Profit() int64 { return c.MoneyReturned - c.MoneyAdvanced }

// Validate returns ErrNonPositiveMoney if MoneyAdvanced ≤ 0, or ErrNoProfit if
// MoneyReturned ≤ MoneyAdvanced (the merchant must extract at least some profit
// from the circuit for it to be a valid commercial-capital movement).
func (c MerchantCircuit) Validate() error {
	if c.MoneyAdvanced <= 0 {
		return ErrNonPositiveMoney
	}
	if c.MoneyReturned <= c.MoneyAdvanced {
		return ErrNoProfit
	}
	return nil
}

// NewMerchantCircuit builds and validates a MerchantCircuit.
func NewMerchantCircuit(moneyAdvanced, commodityValue, moneyReturned int64) (MerchantCircuit, error) {
	c := MerchantCircuit{
		MoneyAdvanced:  moneyAdvanced,
		CommodityValue: commodityValue,
		MoneyReturned:  moneyReturned,
	}
	if err := c.Validate(); err != nil {
		return MerchantCircuit{}, err
	}
	return c, nil
}

// CommodityCapitalForm is the value object for the commodity-capital form that
// the merchant takes off the industrial capitalist's hands (§1). It records the
// triple decomposition (c+v+s) of the commodity-capital's value, and requires
// that SourceCapital == "industrial" — commercial capital can only handle what
// industrial capital has already produced. Computed on demand, not persisted.
type CommodityCapitalForm struct {
	ConstantCapital  int64  `json:"constant_capital"`  // pence; ≥0
	VariableCapital  int64  `json:"variable_capital"`  // pence; ≥0
	SurplusValue     int64  `json:"surplus_value"`     // pence; ≥0
	SourceCapital    string `json:"source_capital"`    // always "industrial"
}

// TotalValue is the sum c+v+s: the total value of the commodity-capital form.
func (f CommodityCapitalForm) TotalValue() int64 {
	return f.ConstantCapital + f.VariableCapital + f.SurplusValue
}

// Validate returns ErrNonIndustrialSource if SourceCapital != "industrial", or
// ErrNegativeMagnitude if any of c, v, s is negative.
func (f CommodityCapitalForm) Validate() error {
	if f.SourceCapital != "industrial" {
		return ErrNonIndustrialSource
	}
	if f.ConstantCapital < 0 || f.VariableCapital < 0 || f.SurplusValue < 0 {
		return ErrNegativeMagnitude
	}
	return nil
}

// NewCommodityCapitalForm builds and validates a CommodityCapitalForm,
// automatically setting SourceCapital = "industrial".
func NewCommodityCapitalForm(c, v, s int64) (CommodityCapitalForm, error) {
	f := CommodityCapitalForm{
		ConstantCapital: c,
		VariableCapital: v,
		SurplusValue:    s,
		SourceCapital:   "industrial",
	}
	if err := f.Validate(); err != nil {
		return CommodityCapitalForm{}, err
	}
	return f, nil
}

// LinenMerchantExample is the Marx §1 fixture: £3,000 advanced by a linen
// merchant who buys 30,000 yards at 10p/yard = 300,000 pence. This exemplar
// is used as the canonical test fixture for this chapter; it is not persisted.
type LinenMerchantExample struct {
	MoneyAdvanced      int64  `json:"money_advanced"`        // 300000 pence (£3,000)
	YardsBought        int64  `json:"yards_bought"`          // 30,000 yards
	PricePerYardPence  int64  `json:"price_per_yard_pence"`  // 10 pence/yard
	Description        string `json:"description"`
}

// LinenMerchant returns the Marx §1 linen-merchant fixture.
func LinenMerchant() LinenMerchantExample {
	return LinenMerchantExample{
		MoneyAdvanced:     300000,
		YardsBought:       30000,
		PricePerYardPence: 10,
		Description:       "Marx §1 (Ch. 16): £3,000 advanced, 30,000 yards at 10p/yard = 300,000 pence",
	}
}
