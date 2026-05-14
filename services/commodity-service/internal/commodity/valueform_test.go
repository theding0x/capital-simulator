package commodity

import "testing"

func makeIron() Commodity {
	return Commodity{
		ID:             "iron-id",
		Name:           "iron",
		UseValue:       UseValue{Description: "pig iron for tools", Unit: "cwt"},
		ConcreteLabour: ConcreteLabour{Kind: "smelting"},
		SNLTPerUnit:    120, // 2h per cwt
	}
}

func makeCorn() Commodity {
	return Commodity{
		ID:             "corn-id",
		Name:           "corn",
		UseValue:       UseValue{Description: "corn for food", Unit: "qtr"},
		ConcreteLabour: ConcreteLabour{Kind: "agriculture"},
		SNLTPerUnit:    240, // 4h per qtr
	}
}

func makeGold() Commodity {
	return Commodity{
		ID:             "gold-id",
		Name:           "gold",
		UseValue:       UseValue{Description: "gold as universal equivalent", Unit: "oz"},
		ConcreteLabour: ConcreteLabour{Kind: "mining"},
		SNLTPerUnit:    60, // 1h per oz
	}
}

func TestSimpleFormOf(t *testing.T) {
	t.Parallel()
	linen := makeLinen()
	linen.ID = "linen-id"
	coat := makeCoat()
	coat.ID = "coat-id"

	sf, err := SimpleFormOf(linen, coat)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sf.Relative.ID != "linen-id" || sf.Equivalent.ID != "coat-id" {
		t.Errorf("relative/equivalent mis-assigned: %+v", sf)
	}
	// 1 yard linen (30m) = 1/20 coat (600m * 1/20 = 30m)
	if sf.RelativeQty != 1 {
		t.Errorf("RelativeQty = %v, want 1", sf.RelativeQty)
	}
	if sf.EquivalentQty != 0.05 {
		t.Errorf("EquivalentQty = %v, want 0.05", sf.EquivalentQty)
	}
	if sf.CommonValue != 30 {
		t.Errorf("CommonValue = %d, want 30", sf.CommonValue)
	}
}

func TestExpandedFormOf(t *testing.T) {
	t.Parallel()
	linen := makeLinen()
	linen.ID = "linen-id"
	coat := makeCoat()
	coat.ID = "coat-id"
	iron := makeIron()
	corn := makeCorn()

	pop := []Commodity{linen, coat, iron, corn}
	ef, err := ExpandedFormOf(linen, pop)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ef.Relative.ID != "linen-id" {
		t.Errorf("Relative = %v, want linen-id", ef.Relative.ID)
	}
	// Should equate linen to every commodity except itself: coat, iron, corn.
	if got := len(ef.Equivalents); got != 3 {
		t.Fatalf("Equivalents count = %d, want 3", got)
	}
	if ef.CommonValue != 30 {
		t.Errorf("CommonValue = %d, want 30 (one yard linen)", ef.CommonValue)
	}
	for _, sf := range ef.Equivalents {
		if sf.Equivalent.ID == "linen-id" {
			t.Errorf("expanded form should not include self as equivalent")
		}
		if sf.CommonValue != 30 {
			t.Errorf("each equivalent should equal the relative's value: got %d, want 30",
				sf.CommonValue)
		}
	}
}

func TestExpandedFormOf_SkipsZeroValueCommodities(t *testing.T) {
	t.Parallel()
	linen := makeLinen()
	linen.ID = "linen-id"
	coat := makeCoat()
	coat.ID = "coat-id"
	free := makeIron()
	free.SNLTPerUnit = 0

	ef, err := ExpandedFormOf(linen, []Commodity{linen, coat, free})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ef.Equivalents) != 1 {
		t.Errorf("Equivalents = %d, want 1 (coat only)", len(ef.Equivalents))
	}
}

func TestGeneralFormOf(t *testing.T) {
	t.Parallel()
	gold := makeGold()
	linen := makeLinen()
	linen.ID = "linen-id"
	coat := makeCoat()
	coat.ID = "coat-id"
	iron := makeIron()

	gf, err := GeneralFormOf(gold, []Commodity{gold, linen, coat, iron})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gf.Equivalent.ID != "gold-id" {
		t.Errorf("Equivalent = %v, want gold-id", gf.Equivalent.ID)
	}
	if got := len(gf.Relatives); got != 3 {
		t.Fatalf("Relatives count = %d, want 3", got)
	}
	for _, sf := range gf.Relatives {
		if sf.Equivalent.ID != "gold-id" {
			t.Errorf("each relative should be expressed in gold; got equivalent=%v",
				sf.Equivalent.ID)
		}
	}
}

func TestGeneralFormOf_RejectsZeroValueEquivalent(t *testing.T) {
	t.Parallel()
	dud := makeGold()
	dud.SNLTPerUnit = 0
	_, err := GeneralFormOf(dud, []Commodity{makeLinen(), makeCoat()})
	if err == nil {
		t.Fatal("expected error when equivalent has no value, got nil")
	}
}

func TestMoneyFormOf(t *testing.T) {
	t.Parallel()
	gold := makeGold()
	linen := makeLinen()
	coat := makeCoat()

	mf, err := MoneyFormOf(gold, []Commodity{gold, linen, coat})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mf.MoneyCommodity.ID != "gold-id" {
		t.Errorf("MoneyCommodity = %v, want gold-id", mf.MoneyCommodity.ID)
	}
	if mf.Equivalent.ID != mf.MoneyCommodity.ID {
		t.Errorf("MoneyForm.Equivalent should be the money-commodity")
	}
	// 1 yard linen (30m) priced in gold (60m/oz) = 0.5 oz
	for _, sf := range mf.Relatives {
		if sf.Relative.Name == "linen" && sf.EquivalentQty != 0.5 {
			t.Errorf("1 yard linen in gold = %v oz, want 0.5", sf.EquivalentQty)
		}
		if sf.Relative.Name == "coat" && sf.EquivalentQty != 10 {
			t.Errorf("1 coat in gold = %v oz, want 10", sf.EquivalentQty)
		}
	}
}

func TestValueFormKind_Validate(t *testing.T) {
	t.Parallel()
	for _, k := range []ValueFormKind{KindSimple, KindExpanded, KindGeneral, KindMoney} {
		if err := k.Validate(); err != nil {
			t.Errorf("Validate(%q) = %v, want nil", k, err)
		}
	}
	if err := ValueFormKind("nonsense").Validate(); err == nil {
		t.Error("expected error for unknown kind, got nil")
	}
}
