package simulation_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

func TestCapitalOriginValidate(t *testing.T) {
	t.Parallel()

	t.Run("valid commerce origin", func(t *testing.T) {
		t.Parallel()
		o := simulation.CapitalOrigin{
			Source:      "commerce",
			AmountPence: 50000,
			Period:      "15th-17th c.",
		}
		if err := o.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("source required", func(t *testing.T) {
		t.Parallel()
		o := simulation.CapitalOrigin{AmountPence: 1000}
		if err := o.Validate(); err == nil {
			t.Fatal("expected error for missing source")
		}
	})

	t.Run("amount must be positive", func(t *testing.T) {
		t.Parallel()
		o := simulation.CapitalOrigin{Source: "usury", AmountPence: 0}
		if err := o.Validate(); err == nil {
			t.Fatal("expected error for zero amount")
		}
	})

	t.Run("all valid sources accepted", func(t *testing.T) {
		t.Parallel()
		sources := []string{
			"usury", "commerce", "colonial-plunder",
			"national-debt", "taxation", "guild-master-accumulation",
		}
		for _, s := range sources {
			o := simulation.CapitalOrigin{Source: s, AmountPence: 1000}
			if err := o.Validate(); err != nil {
				t.Errorf("source %q should be valid: %v", s, err)
			}
		}
	})
}

func TestColonialTransferValidate(t *testing.T) {
	t.Parallel()

	t.Run("valid slave trade transfer — Liverpool 1730", func(t *testing.T) {
		t.Parallel()
		// Liverpool employed in the slave-trade, in 1730, 15 ships
		tr := simulation.ColonialTransfer{
			From:       "West Africa",
			To:         "England",
			ValuePence: 120000,
			Method:     "slave-trade",
		}
		if err := tr.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("from required", func(t *testing.T) {
		t.Parallel()
		tr := simulation.ColonialTransfer{To: "England", ValuePence: 1000, Method: "colonial-plunder"}
		if err := tr.Validate(); err == nil {
			t.Fatal("expected error for missing from")
		}
	})

	t.Run("to required", func(t *testing.T) {
		t.Parallel()
		tr := simulation.ColonialTransfer{From: "Americas", ValuePence: 1000, Method: "colonial-plunder"}
		if err := tr.Validate(); err == nil {
			t.Fatal("expected error for missing to")
		}
	})

	t.Run("value must be positive", func(t *testing.T) {
		t.Parallel()
		tr := simulation.ColonialTransfer{From: "Americas", To: "England", ValuePence: 0, Method: "colonial-plunder"}
		if err := tr.Validate(); err == nil {
			t.Fatal("expected error for zero value")
		}
	})

	t.Run("method required", func(t *testing.T) {
		t.Parallel()
		tr := simulation.ColonialTransfer{From: "Americas", To: "England", ValuePence: 1000}
		if err := tr.Validate(); err == nil {
			t.Fatal("expected error for missing method")
		}
	})
}

func TestNationalDebtValidate(t *testing.T) {
	t.Parallel()

	t.Run("Bank of England founding — 8% to private bankers", func(t *testing.T) {
		t.Parallel()
		// Bank of England began with lending its money to the Government at 8%
		d := simulation.NationalDebt{
			AmountPence:     2400000,
			InterestRateBps: 800, // 8% = 800 basis points
			CreditorClass:   "private-bankers",
		}
		if err := d.Validate(); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("amount must be positive", func(t *testing.T) {
		t.Parallel()
		d := simulation.NationalDebt{AmountPence: 0, InterestRateBps: 800, CreditorClass: "private-bankers"}
		if err := d.Validate(); err == nil {
			t.Fatal("expected error for zero amount")
		}
	})

	t.Run("interest rate must be positive", func(t *testing.T) {
		t.Parallel()
		d := simulation.NationalDebt{AmountPence: 1000000, InterestRateBps: 0, CreditorClass: "private-bankers"}
		if err := d.Validate(); err == nil {
			t.Fatal("expected error for zero interest rate")
		}
	})

	t.Run("creditor class required", func(t *testing.T) {
		t.Parallel()
		d := simulation.NationalDebt{AmountPence: 1000000, InterestRateBps: 500}
		if err := d.Validate(); err == nil {
			t.Fatal("expected error for missing creditor class")
		}
	})
}

func TestComputeGenesis(t *testing.T) {
	t.Parallel()

	stageID := simulation.HistoricalStageID("5eed000000000000002601")

	t.Run("total capital is sum of origins and transfers", func(t *testing.T) {
		t.Parallel()
		origins := []simulation.CapitalOrigin{
			{Source: "usury", AmountPence: 30000},
			{Source: "commerce", AmountPence: 20000},
		}
		transfers := []simulation.ColonialTransfer{
			{From: "Americas", To: "England", ValuePence: 50000, Method: "colonial-plunder"},
		}
		debts := []simulation.NationalDebt{
			{AmountPence: 2400000, InterestRateBps: 800, CreditorClass: "private-bankers"},
		}
		systems := []simulation.ProtectionSystem{
			{TariffRateBps: 3000, Beneficiary: "English manufacturers", PeriodStart: "17th c", PeriodEnd: "19th c"},
		}

		g := simulation.ComputeGenesis(stageID, origins, transfers, debts, systems)

		want := simulation.Pence(100000) // 30000 + 20000 + 50000
		if g.TotalCapitalFormedPence != want {
			t.Errorf("TotalCapitalFormedPence = %d, want %d", g.TotalCapitalFormedPence, want)
		}
	})

	t.Run("national debts excluded from capital total", func(t *testing.T) {
		t.Parallel()
		debts := []simulation.NationalDebt{
			{AmountPence: 9999999, InterestRateBps: 800, CreditorClass: "private-bankers"},
		}
		g := simulation.ComputeGenesis(stageID, nil, nil, debts, nil)
		if g.TotalCapitalFormedPence != 0 {
			t.Errorf("national debts should not add to capital total; got %d", g.TotalCapitalFormedPence)
		}
	})

	t.Run("nil slices coerced to empty", func(t *testing.T) {
		t.Parallel()
		g := simulation.ComputeGenesis(stageID, nil, nil, nil, nil)
		if g.Origins == nil {
			t.Error("Origins should not be nil")
		}
		if g.ColonialTransfers == nil {
			t.Error("ColonialTransfers should not be nil")
		}
		if g.NationalDebts == nil {
			t.Error("NationalDebts should not be nil")
		}
		if g.ProtectionSystems == nil {
			t.Error("ProtectionSystems should not be nil")
		}
	})

	t.Run("stage ID preserved", func(t *testing.T) {
		t.Parallel()
		g := simulation.ComputeGenesis(stageID, nil, nil, nil, nil)
		if g.HistoricalStageID != stageID {
			t.Errorf("HistoricalStageID = %q, want %q", g.HistoricalStageID, stageID)
		}
	})

	t.Run("Liverpool slave trade series — 1730 through 1792", func(t *testing.T) {
		t.Parallel()
		// Liverpool employed in the slave-trade, in 1730, 15 ships; 1751, 53; 1760, 74; 1770, 96; 1792, 132
		transfers := []simulation.ColonialTransfer{
			{From: "West Africa", To: "England", ValuePence: 15000, Method: "slave-trade"},
			{From: "West Africa", To: "England", ValuePence: 53000, Method: "slave-trade"},
			{From: "West Africa", To: "England", ValuePence: 74000, Method: "slave-trade"},
			{From: "West Africa", To: "England", ValuePence: 96000, Method: "slave-trade"},
			{From: "West Africa", To: "England", ValuePence: 132000, Method: "slave-trade"},
		}
		g := simulation.ComputeGenesis(stageID, nil, transfers, nil, nil)
		want := simulation.Pence(15000 + 53000 + 74000 + 96000 + 132000)
		if g.TotalCapitalFormedPence != want {
			t.Errorf("TotalCapitalFormedPence = %d, want %d", g.TotalCapitalFormedPence, want)
		}
		if len(g.ColonialTransfers) != 5 {
			t.Errorf("ColonialTransfers len = %d, want 5", len(g.ColonialTransfers))
		}
	})
}

func TestNewCapitalOriginID(t *testing.T) {
	t.Parallel()
	id1 := simulation.NewCapitalOriginID()
	id2 := simulation.NewCapitalOriginID()
	if id1 == id2 {
		t.Error("NewCapitalOriginID should produce unique IDs")
	}
	if len(string(id1)) != 24 {
		t.Errorf("CapitalOriginID length = %d, want 24", len(string(id1)))
	}
}

func TestNewColonialTransferID(t *testing.T) {
	t.Parallel()
	id1 := simulation.NewColonialTransferID()
	id2 := simulation.NewColonialTransferID()
	if id1 == id2 {
		t.Error("NewColonialTransferID should produce unique IDs")
	}
}

func TestNewNationalDebtID(t *testing.T) {
	t.Parallel()
	id1 := simulation.NewNationalDebtID()
	id2 := simulation.NewNationalDebtID()
	if id1 == id2 {
		t.Error("NewNationalDebtID should produce unique IDs")
	}
}
