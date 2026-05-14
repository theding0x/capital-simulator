package commodity

import (
	"strings"
	"testing"
)

// TestValue_MarxLinenAndCoat exercises Marx's classic example from §1:
// 20 yards of linen = 1 coat. Linen is 30 minutes/yard; a coat is 10 hours.
// Both bundles congeal exactly 600 minutes of abstract human labour.
func TestValue_MarxLinenAndCoat(t *testing.T) {
	t.Parallel()
	linen := makeLinen()
	coat := makeCoat()

	if got := linen.Value(20); got != 600 {
		t.Errorf("Value(20 yards linen) = %d, want 600", got)
	}
	if got := coat.Value(1); got != 600 {
		t.Errorf("Value(1 coat) = %d, want 600", got)
	}
}

// TestValue_NegativeQuantity guards against value going negative when a
// caller passes nonsense quantities. Negative quantities are clamped to zero.
func TestValue_NegativeQuantity(t *testing.T) {
	t.Parallel()
	if got := makeLinen().Value(-5); got != 0 {
		t.Errorf("Value(-5) = %d, want 0", got)
	}
}

// TestValue_DualCharacterOfLabour - §2: equal SNLT yields equal value
// regardless of which concrete labour produced it. Weaving and tailoring,
// considered as expenditures of homogeneous human labour-power, count
// identically.
func TestValue_DualCharacterOfLabour(t *testing.T) {
	t.Parallel()
	weaver := Commodity{
		Name:           "linen",
		UseValue:       UseValue{Description: "linen", Unit: "yards"},
		ConcreteLabour: ConcreteLabour{Kind: "weaving"},
		SNLTPerUnit:    60,
	}
	tailor := Commodity{
		Name:           "coat-half",
		UseValue:       UseValue{Description: "half a coat", Unit: "halves"},
		ConcreteLabour: ConcreteLabour{Kind: "tailoring"},
		SNLTPerUnit:    60,
	}
	if weaver.Value(1) != tailor.Value(1) {
		t.Errorf("equal SNLT should yield equal value: %d vs %d",
			weaver.Value(1), tailor.Value(1))
	}
}

func TestExchangeRatioBetween_LinenForCoat(t *testing.T) {
	t.Parallel()
	linen := makeLinen()
	coat := makeCoat()

	r, err := ExchangeRatioBetween(linen, coat, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.CommonValue != 600 {
		t.Errorf("CommonValue = %d, want 600", r.CommonValue)
	}
	if r.QuoteQty != 1 {
		t.Errorf("20 yards linen should equal 1 coat; got %v", r.QuoteQty)
	}
}

func TestExchangeRatioBetween_CoatForLinen(t *testing.T) {
	t.Parallel()
	r, err := ExchangeRatioBetween(makeCoat(), makeLinen(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.QuoteQty != 20 {
		t.Errorf("1 coat should equal 20 yards linen; got %v", r.QuoteQty)
	}
}

func TestExchangeRatioBetween_RejectsNonpositiveBase(t *testing.T) {
	t.Parallel()
	_, err := ExchangeRatioBetween(makeLinen(), makeCoat(), 0)
	if err == nil {
		t.Fatal("expected error for zero base_qty, got nil")
	}
}

func TestExchangeRatioBetween_RejectsZeroValueQuote(t *testing.T) {
	t.Parallel()
	free := makeCoat()
	free.SNLTPerUnit = 0
	_, err := ExchangeRatioBetween(makeLinen(), free, 1)
	if err == nil {
		t.Fatal("expected error for zero-SNLT quote, got nil")
	}
	if !strings.Contains(err.Error(), "no value") {
		t.Errorf("error should mention zero value: %v", err)
	}
}

// TestProductivity_InverseLaw - §1: doubled productivity halves value.
func TestProductivity_InverseLaw(t *testing.T) {
	t.Parallel()
	linen := makeLinen() // 30 min/yard
	doubled, err := ProductivityChange{Factor: 2}.Apply(linen)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if doubled.SNLTPerUnit != 15 {
		t.Errorf("doubled productivity should halve SNLT: got %d, want 15",
			doubled.SNLTPerUnit)
	}
	// And the value of any quantity should also halve.
	if doubled.Value(20) != 300 {
		t.Errorf("Value(20 yards) at doubled productivity = %d, want 300",
			doubled.Value(20))
	}
}

func TestProductivity_HalvedRaisesValue(t *testing.T) {
	t.Parallel()
	linen := makeLinen() // 30 min/yard
	halved, err := ProductivityChange{Factor: 0.5}.Apply(linen)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if halved.SNLTPerUnit != 60 {
		t.Errorf("halved productivity should double SNLT: got %d, want 60",
			halved.SNLTPerUnit)
	}
}

func TestProductivity_RejectsNonPositiveFactor(t *testing.T) {
	t.Parallel()
	for _, f := range []float64{0, -1, -0.5} {
		_, err := ProductivityChange{Factor: f}.Apply(makeLinen())
		if err == nil {
			t.Errorf("expected error for factor %v, got nil", f)
		}
	}
}

func TestAsAbstractLabour_Identity(t *testing.T) {
	t.Parallel()
	// The reduction is the identity on labour-time: a minute of weaving
	// counts the same as a minute of tailoring (Capital I, §2).
	for _, kind := range []string{"weaving", "tailoring", "smelting"} {
		got := AsAbstractLabour(ConcreteLabour{Kind: kind}, 42)
		if got != 42 {
			t.Errorf("AsAbstractLabour(%q, 42) = %d, want 42", kind, got)
		}
	}
}
