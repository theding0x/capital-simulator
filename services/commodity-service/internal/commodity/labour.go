package commodity

import (
	"errors"
	"fmt"
	"strings"
)

// LabourMinutes measures labour-time in whole minutes. It is the canonical
// unit for value-magnitude in this simulation: every Commodity records
// SNLTPerUnit as LabourMinutes, and Value is always expressed in
// LabourMinutes regardless of whose concrete labour produced it.
//
// Marx, Ch. 1 §1: "The quantity of labour, however, is measured by its
// duration, and labour-time in its turn finds its standard in weeks, days
// and hours." We choose minutes as the simulation's atomic unit.
type LabourMinutes int64

// Hours expresses the labour-time as fractional hours. Useful in summaries
// where minutes are too granular to read.
func (lm LabourMinutes) Hours() float64 { return float64(lm) / 60.0 }

// String renders the labour-time as "Xh Ym" or "Ym" for sub-hour values.
func (lm LabourMinutes) String() string {
	if lm < 60 {
		return fmt.Sprintf("%dm", lm)
	}
	h := lm / 60
	m := lm % 60
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	return fmt.Sprintf("%dh %dm", h, m)
}

// ConcreteLabour describes a particular kind of useful labour - "tailoring,
// weaving, spinning, etc." (Capital I, Ch. 1, §2). Concrete labour produces
// the use-value of a commodity. Within the simulation, every commodity
// records the concrete labour that produces it; for the purposes of value
// magnitude, however, all such concrete labours are reduced to homogeneous
// human labour in the abstract.
type ConcreteLabour struct {
	// Kind is a short identifier for the type of useful labour
	// (e.g. "weaving", "tailoring", "smelting", "agriculture").
	Kind string `json:"kind"`

	// Description is a human-readable explanation of the labour process.
	Description string `json:"description,omitempty"`
}

// Validate ensures the concrete-labour record carries at least a kind.
func (c ConcreteLabour) Validate() error {
	if strings.TrimSpace(c.Kind) == "" {
		return errors.New("commodity: concrete_labour.kind is required")
	}
	return nil
}

// AsAbstractLabour reduces concrete labour to abstract human labour. This is
// the operation that Marx flags as the pivotal move in §2 - "all labour is,
// speaking physiologically, an expenditure of human labour-power, and in its
// character of identical abstract human labour, it creates and forms the
// value of commodities."
//
// Operationally the reduction is the identity on labour-time: a minute of
// weaving and a minute of tailoring count as a minute of abstract human
// labour. The function exists so that the model makes the reduction explicit
// at every value-computation site instead of smuggling it in implicitly.
func AsAbstractLabour(_ ConcreteLabour, t LabourMinutes) LabourMinutes {
	return t
}
