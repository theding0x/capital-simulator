package market

import (
	"errors"
	"testing"
)

// weaverCircuit is the canonical Ch. 3 §2 fixture: the weaver sells 20 yards of
// linen for £2 (200 pence), then lays out that same £2 on a family Bible.
func weaverCircuit() Circuit {
	return Circuit{
		SaleLeg:     CircuitLeg{Kind: KindSale, CommodityID: "linen", MoneyID: "gold", OwnerID: "weaver", Value: 200},
		PurchaseLeg: CircuitLeg{Kind: KindPurchase, CommodityID: "bible", MoneyID: "gold", OwnerID: "weaver", Value: 200},
	}
}

// TestCircuitValidate_Equivalent asserts the load-bearing invariant: a circuit
// whose sale and purchase legs carry the same value-magnitude is valid. "The
// same value-magnitude leaves in commodity-form and returns in commodity-form."
func TestCircuitValidate_Equivalent(t *testing.T) {
	t.Parallel()
	if err := weaverCircuit().Validate(); err != nil {
		t.Fatalf("equivalent circuit should validate, got %v", err)
	}
}

// TestCircuitValidate_NonEquivalentRejected asserts that unequal legs are
// rejected with ErrCircuitNotEquivalent — the difference would be surplus-value,
// which belongs to the M—C—M′ capital circuit (Ch. 4), not simple circulation.
func TestCircuitValidate_NonEquivalentRejected(t *testing.T) {
	t.Parallel()
	c := weaverCircuit()
	c.PurchaseLeg.Value = 250 // bought a £2½ Bible with £2 of sale-money
	if err := c.Validate(); !errors.Is(err, ErrCircuitNotEquivalent) {
		t.Fatalf("want ErrCircuitNotEquivalent, got %v", err)
	}
}

// TestCircuitValidate_WrongLegKinds rejects a circuit whose legs are not a sale
// followed by a purchase.
func TestCircuitValidate_WrongLegKinds(t *testing.T) {
	t.Parallel()
	c := weaverCircuit()
	c.SaleLeg.Kind = KindPurchase
	if err := c.Validate(); err == nil {
		t.Fatal("expected error when sale_leg kind is not a sale")
	}
}

// TestCircuitValidate_LegMissingFields exercises the per-leg field validation.
func TestCircuitValidate_LegMissingFields(t *testing.T) {
	t.Parallel()
	cases := map[string]func(*Circuit){
		"missing sale owner":     func(c *Circuit) { c.SaleLeg.OwnerID = "" },
		"missing purchase money": func(c *Circuit) { c.PurchaseLeg.MoneyID = "" },
		"zero sale value":        func(c *Circuit) { c.SaleLeg.Value = 0; c.PurchaseLeg.Value = 0 },
	}
	for name, mutate := range cases {
		mutate := mutate
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c := weaverCircuit()
			mutate(&c)
			if err := c.Validate(); err == nil {
				t.Fatalf("%s: expected validation error", name)
			}
		})
	}
}
