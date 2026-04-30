// Package commodity models the commodity-form analyzed in
// Karl Marx's Capital Vol. I, Ch. 1.
//
// The commodity has two factors:
//
//  1. Use-value, the qualitative utility of the thing - bound to its physical
//     properties and only realized in consumption.
//  2. Value, an abstraction made when we abstract from use-value. The substance
//     of value is human labour in the abstract; its magnitude is socially
//     necessary labour-time (SNLT).
//
// This file defines the core types Commodity and UseValue. Labour, value
// computation, value-forms, and the social-relations view that makes
// commodity fetishism explicit are split across labour.go, value.go,
// valueform.go, and fetishism.go.
package commodity

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ID uniquely identifies a Commodity.
type ID string

// NewID returns a new random hex ID. 96 bits is plenty for our universe.
func NewID() ID {
	var b [12]byte
	if _, err := rand.Read(b[:]); err != nil {
		// crypto/rand should never fail; if it does, panic is appropriate.
		panic(fmt.Errorf("commodity: rand: %w", err))
	}
	return ID(hex.EncodeToString(b[:]))
}

// IsZero reports whether id is the empty ID.
func (id ID) IsZero() bool { return id == "" }

// Commodity is "an object outside us, a thing that by its properties satisfies
// human wants of some sort or another." (Capital I, Ch. 1, §1)
//
// A Commodity has two factors: a use-value (its qualitative utility) and a
// value (the magnitude of socially necessary labour-time congealed in it).
// Concrete useful labour produces the use-value; the same labour, considered
// as homogeneous human labour-power, constitutes the value.
type Commodity struct {
	ID             ID             `json:"id"`
	Name           string         `json:"name"`
	UseValue       UseValue       `json:"use_value"`
	ConcreteLabour ConcreteLabour `json:"concrete_labour"`

	// SNLTPerUnit is the socially necessary labour-time required to produce
	// one Unit of this commodity under the conditions of production normal for
	// the society and with the average degree of skill and intensity of labour
	// prevalent in it (Capital I, Ch. 1, §1).
	SNLTPerUnit LabourMinutes `json:"snlt_per_unit"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UseValue captures the qualitative side of a commodity: what want it
// satisfies and the unit in which it is socially measured. "The utility of a
// thing makes it a use value" (Capital I, Ch. 1, §1).
type UseValue struct {
	// Description is a free-text explanation of the want this commodity
	// satisfies (e.g. "linen for clothing").
	Description string `json:"description"`
	// Unit is the socially-recognized standard of measure for quantities
	// of this commodity (e.g. "yards", "cwt", "qtr"). Diverse measures arise
	// "partly in the diverse nature of the objects to be measured, partly in
	// convention" (Capital I, Ch. 1, §1).
	Unit string `json:"unit"`
}

// Validate checks that c is well-formed enough to persist.
func (c Commodity) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return errors.New("commodity: name is required")
	}
	if strings.TrimSpace(c.UseValue.Description) == "" {
		return errors.New("commodity: use_value.description is required")
	}
	if strings.TrimSpace(c.UseValue.Unit) == "" {
		return errors.New("commodity: use_value.unit is required")
	}
	if c.SNLTPerUnit < 0 {
		return errors.New("commodity: snlt_per_unit must be non-negative")
	}
	if err := c.ConcreteLabour.Validate(); err != nil {
		return err
	}
	return nil
}
