package commodity

import (
	"errors"
	"fmt"
)

// ValueForm models the four forms of value Marx develops in §3:
//
//	A. Simple, Isolated, or Accidental Form:    x A = y B
//	B. Total or Expanded Form:                  x A = y B = z C = ...
//	C. General Form:                            x A, y B, z C, ... = w D
//	D. Money Form:                              x A, y B, z C, ... = w gold
//
// The simulation builds these as derived views over a population of
// Commodities and a chosen relative or equivalent.
//
// In the Relative form, the commodity whose value is being expressed appears
// on the left; in the Equivalent form, the commodity expressing that value
// appears on the right and counts as immediately exchangeable for any other
// commodity.

// SimpleForm is "x of A = y of B": one commodity (the relative) has its
// value expressed in another (the equivalent).
type SimpleForm struct {
	Relative   Commodity `json:"relative"`
	Equivalent Commodity `json:"equivalent"`
	RelativeQty Quantity `json:"relative_qty"`
	EquivalentQty Quantity `json:"equivalent_qty"`
	CommonValue LabourMinutes `json:"common_value"`
}

// SimpleFormOf builds the simple value-form for one unit of relative,
// expressed in equivalent.
func SimpleFormOf(relative, equivalent Commodity) (SimpleForm, error) {
	r, err := ExchangeRatioBetween(relative, equivalent, 1)
	if err != nil {
		return SimpleForm{}, err
	}
	return SimpleForm{
		Relative:      r.Base,
		Equivalent:    r.Quote,
		RelativeQty:   r.BaseQty,
		EquivalentQty: r.QuoteQty,
		CommonValue:   r.CommonValue,
	}, nil
}

// ExpandedForm is "x A = y B = z C = ...": the value of one relative
// commodity expressed in many equivalents.
//
// Marx (§3 B): "By this form ... the commodity now appears to be in social
// relation, no longer with only one other kind of commodity, but with the
// whole world of commodities."
type ExpandedForm struct {
	Relative    Commodity     `json:"relative"`
	RelativeQty Quantity      `json:"relative_qty"`
	Equivalents []SimpleForm  `json:"equivalents"`
	CommonValue LabourMinutes `json:"common_value"`
}

// ExpandedFormOf builds the expanded value-form for one unit of relative
// against every other commodity in the population.
func ExpandedFormOf(relative Commodity, population []Commodity) (ExpandedForm, error) {
	common := relative.Value(1)
	out := ExpandedForm{
		Relative:    relative,
		RelativeQty: 1,
		CommonValue: common,
	}
	for _, c := range population {
		if c.ID == relative.ID || c.ID.IsZero() && c.Name == relative.Name {
			continue
		}
		if c.SNLTPerUnit <= 0 {
			continue
		}
		sf, err := SimpleFormOf(relative, c)
		if err != nil {
			return ExpandedForm{}, err
		}
		out.Equivalents = append(out.Equivalents, sf)
	}
	return out, nil
}

// GeneralForm inverts the expanded form: many relatives, one equivalent.
// Marx (§3 C): "all commodities now express their value (1) in an elementary
// form, because in a single commodity; (2) with unity, because in one and the
// same commodity. This form of value is elementary and the same for all,
// therefore general."
type GeneralForm struct {
	Equivalent Commodity     `json:"equivalent"`
	Relatives  []SimpleForm  `json:"relatives"`
}

// GeneralFormOf builds the general value-form, expressing the value of every
// other commodity in the population in terms of the chosen general equivalent.
func GeneralFormOf(equivalent Commodity, population []Commodity) (GeneralForm, error) {
	if equivalent.SNLTPerUnit <= 0 {
		return GeneralForm{}, fmt.Errorf("commodity: equivalent %q has no value", equivalent.Name)
	}
	out := GeneralForm{Equivalent: equivalent}
	for _, c := range population {
		if c.ID == equivalent.ID || (c.ID.IsZero() && c.Name == equivalent.Name) {
			continue
		}
		if c.SNLTPerUnit <= 0 {
			continue
		}
		sf, err := SimpleFormOf(c, equivalent)
		if err != nil {
			return GeneralForm{}, err
		}
		out.Relatives = append(out.Relatives, sf)
	}
	return out, nil
}

// MoneyForm is the General form with a designated money-commodity. Marx
// (§3 D): "The specific kind of commodity with whose natural form the
// equivalent form is socially identified now becomes the money-commodity,
// or serves as money."
//
// Modeling-wise it is a GeneralForm whose Equivalent has been crowned the
// money-commodity. The crown is a piece of social information rather than a
// computational difference - which is precisely Marx's point.
type MoneyForm struct {
	GeneralForm
	MoneyCommodity Commodity `json:"money_commodity"`
}

// MoneyFormOf builds the money-form: pick a commodity from the population to
// serve as money, and express every other commodity's value in it.
func MoneyFormOf(money Commodity, population []Commodity) (MoneyForm, error) {
	g, err := GeneralFormOf(money, population)
	if err != nil {
		return MoneyForm{}, err
	}
	return MoneyForm{GeneralForm: g, MoneyCommodity: money}, nil
}

// ValueFormKind identifies which of the four forms a caller wants rendered.
type ValueFormKind string

const (
	KindSimple   ValueFormKind = "simple"
	KindExpanded ValueFormKind = "expanded"
	KindGeneral  ValueFormKind = "general"
	KindMoney    ValueFormKind = "money"
)

// Validate reports whether k is one of the recognized value-form kinds.
func (k ValueFormKind) Validate() error {
	switch k {
	case KindSimple, KindExpanded, KindGeneral, KindMoney:
		return nil
	default:
		return errors.New("commodity: unknown value-form kind")
	}
}
