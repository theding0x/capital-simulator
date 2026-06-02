package agent

import (
	"errors"
	"testing"
)

// Ch. 2 fixture: 20 yards of linen = 1 coat. The weaver sells the linen
// (a sale leg) to the tailor, who buys it (a purchase leg); the realised
// value is 1200 labour-minutes.
const (
	weaverID = "5eed00000000000000000003"
	tailorID = "5eed00000000000000000004"
	linenID  = "000000000000000000000001"
)

func TestExchangeSettlement_Legs(t *testing.T) {
	t.Parallel()
	s := ExchangeSettlement{
		ExchangeID:          "exch-001",
		GiverID:             weaverID,
		ReceiverID:          tailorID,
		GiverCommodityID:    linenID,
		ReceiverCommodityID: "000000000000000000000099", // money
		RealisedValue:       1200,
	}
	giver, receiver := s.Legs()

	if giver.Kind != CircuitLegSale {
		t.Errorf("giver kind = %q, want sale", giver.Kind)
	}
	if receiver.Kind != CircuitLegPurchase {
		t.Errorf("receiver kind = %q, want purchase", receiver.Kind)
	}
	if giver.AgentID != weaverID || giver.CounterpartyID != tailorID {
		t.Errorf("giver leg parties = %s↔%s, want %s↔%s", giver.AgentID, giver.CounterpartyID, weaverID, tailorID)
	}
	if receiver.AgentID != tailorID || receiver.CounterpartyID != weaverID {
		t.Errorf("receiver leg parties = %s↔%s, want %s↔%s", receiver.AgentID, receiver.CounterpartyID, tailorID, weaverID)
	}
	if giver.CommodityID != linenID || receiver.CommodityID != linenID {
		t.Error("both legs should reference the commodity that changed hands (the linen)")
	}
	if giver.ValueMinutes != 1200 || receiver.ValueMinutes != 1200 {
		t.Errorf("both legs should carry the realised value 1200, got %d/%d", giver.ValueMinutes, receiver.ValueMinutes)
	}
	if giver.ExchangeID != "exch-001" || receiver.ExchangeID != "exch-001" {
		t.Error("both legs should carry the originating exchange id")
	}
}

func TestExchangeSettlement_Validate(t *testing.T) {
	t.Parallel()
	valid := ExchangeSettlement{
		ExchangeID:       "exch-001",
		GiverID:          weaverID,
		ReceiverID:       tailorID,
		GiverCommodityID: linenID,
		RealisedValue:    1200,
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid settlement rejected: %v", err)
	}

	self := valid
	self.ReceiverID = weaverID
	if err := self.Validate(); !errors.Is(err, ErrCircuitLegSelfExchange) {
		t.Fatalf("self-exchange Validate() = %v, want ErrCircuitLegSelfExchange", err)
	}

	noValue := valid
	noValue.RealisedValue = 0
	if err := noValue.Validate(); err == nil {
		t.Fatal("zero realised value must be rejected")
	}
}

func TestCircuitLeg_Validate(t *testing.T) {
	t.Parallel()
	base := CircuitLeg{
		AgentID:        weaverID,
		Kind:           CircuitLegSale,
		CounterpartyID: tailorID,
		CommodityID:    linenID,
		ValueMinutes:   1200,
		ExchangeID:     "exch-001",
	}
	if err := base.Validate(); err != nil {
		t.Fatalf("valid leg rejected: %v", err)
	}

	cases := []struct {
		name   string
		mutate func(CircuitLeg) CircuitLeg
	}{
		{"missing agent", func(l CircuitLeg) CircuitLeg { l.AgentID = ""; return l }},
		{"missing counterparty", func(l CircuitLeg) CircuitLeg { l.CounterpartyID = ""; return l }},
		{"self exchange", func(l CircuitLeg) CircuitLeg { l.CounterpartyID = l.AgentID; return l }},
		{"bad kind", func(l CircuitLeg) CircuitLeg { l.Kind = "barter"; return l }},
		{"missing commodity", func(l CircuitLeg) CircuitLeg { l.CommodityID = ""; return l }},
		{"non-positive value", func(l CircuitLeg) CircuitLeg { l.ValueMinutes = 0; return l }},
		{"missing exchange", func(l CircuitLeg) CircuitLeg { l.ExchangeID = ""; return l }},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if err := tc.mutate(base).Validate(); err == nil {
				t.Fatalf("%s should be rejected", tc.name)
			}
		})
	}
}
