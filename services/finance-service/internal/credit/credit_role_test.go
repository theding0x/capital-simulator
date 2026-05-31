package credit

import (
	"strings"
	"testing"
)

// ── CreditRole.IsValid ────────────────────────────────────────────────────────

func TestCreditRole_IsValid(t *testing.T) {
	t.Parallel()

	valid := []CreditRole{RoleEqualiseProfitRate, RoleReduceCirculationCosts, RoleStockCompanies}
	for _, r := range valid {
		if !r.IsValid() {
			t.Errorf("expected role %d to be valid", r)
		}
	}

	invalid := []CreditRole{0, 4, 99, -1}
	for _, r := range invalid {
		if r.IsValid() {
			t.Errorf("expected role %d to be invalid", r)
		}
	}
}

// ── StockCompany — United Alkali Trust fixture ────────────────────────────────
//
// Marx's fixture (Vol. III Ch. 27): United Alkali Trust.
//   Fixed capital  £5,000,000 → 500,000,000 pence
//   Floating capital £1,000,000 → 100,000,000 pence
//   Total capital  £6,000,000 → 600,000,000 pence

func TestNewStockCompany_UnitedAlkaliTrust(t *testing.T) {
	t.Parallel()

	mgr := Manager{Name: "Salaried Director", Title: "Managing Director", WageMinutes: 120}
	sc, err := NewStockCompany(
		"United Alkali Trust",
		500_000_000, // fixed
		100_000_000, // floating
		600_000_000, // total = fixed + floating
		mgr,
		"Joint-stock amalgamation of alkali firms; private ownership socialised within capitalism",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sc.TotalCapital != sc.FixedCapital+sc.FloatingCapital {
		t.Errorf("TotalCapital invariant broken: got %d != %d + %d",
			sc.TotalCapital, sc.FixedCapital, sc.FloatingCapital)
	}
	if sc.Manager.WageMinutes != 120 {
		t.Errorf("expected WageMinutes=120, got %d", sc.Manager.WageMinutes)
	}
	if sc.Name != "United Alkali Trust" {
		t.Errorf("unexpected name: %q", sc.Name)
	}
}

func TestNewStockCompany_CapitalMismatch(t *testing.T) {
	t.Parallel()

	mgr := Manager{Name: "Director", Title: "MD", WageMinutes: 60}
	_, err := NewStockCompany("Bad Co", 500_000_000, 100_000_000, 500_000_000, mgr, "")
	if err == nil {
		t.Fatal("expected ErrCapitalMismatch, got nil")
	}
	if err != ErrCapitalMismatch {
		t.Errorf("expected ErrCapitalMismatch, got %v", err)
	}
}

func TestNewStockCompany_ZeroWage(t *testing.T) {
	t.Parallel()

	mgr := Manager{Name: "Director", Title: "MD", WageMinutes: 0}
	_, err := NewStockCompany("Bad Co", 500_000_000, 100_000_000, 600_000_000, mgr, "")
	if err == nil {
		t.Fatal("expected ErrNonPositiveWage, got nil")
	}
	if err != ErrNonPositiveWage {
		t.Errorf("expected ErrNonPositiveWage, got %v", err)
	}
}

func TestNewStockCompany_NegativeWage(t *testing.T) {
	t.Parallel()

	mgr := Manager{Name: "Director", Title: "MD", WageMinutes: -1}
	_, err := NewStockCompany("Bad Co", 500_000_000, 100_000_000, 600_000_000, mgr, "")
	if err != ErrNonPositiveWage {
		t.Errorf("expected ErrNonPositiveWage, got %v", err)
	}
}

// ── CooperativeFactory ────────────────────────────────────────────────────────

func TestNewCooperativeFactory_Happy(t *testing.T) {
	t.Parallel()

	cf, err := NewCooperativeFactory(
		"Rochdale Pioneers",
		true,  // workerOwned
		28,    // workerCount
		10000, // capitalAdvanced (pence)
		"First worker cooperative; positive sprout of new mode of production",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cf.WorkerOwned {
		t.Error("expected WorkerOwned=true")
	}
	if cf.WorkerCount != 28 {
		t.Errorf("expected WorkerCount=28, got %d", cf.WorkerCount)
	}
}

func TestNewCooperativeFactory_NotWorkerOwned(t *testing.T) {
	t.Parallel()

	_, err := NewCooperativeFactory("Bad Factory", false, 10, 5000, "")
	if err == nil {
		t.Fatal("expected ErrNotWorkerOwned, got nil")
	}
	if err != ErrNotWorkerOwned {
		t.Errorf("expected ErrNotWorkerOwned, got %v", err)
	}
}

// ── CirculationCostSaving — Scottish banking fixture ──────────────────────────
//
// Marx's fixture (Vol. III Ch. 27): Scotland.
//   £27m of deposits circulate via £3m of currency — saving = £24m.
//   In pence: 2,700,000,000 - 300,000,000 = 2,400,000,000.

func TestComputeCirculationCostSaving_Scottish(t *testing.T) {
	t.Parallel()

	const (
		metallicEq = int64(2_700_000_000) // £27m in pence
		actual     = int64(300_000_000)   // £3m in pence
		wantSaving = int64(2_400_000_000) // £24m in pence
	)

	saving := ComputeCirculationCostSaving(metallicEq, actual, "Scottish banking 1826")
	if saving.SavingAmount != wantSaving {
		t.Errorf("SavingAmount: got %d, want %d", saving.SavingAmount, wantSaving)
	}
	if saving.MetallicEquivalent != metallicEq {
		t.Errorf("MetallicEquivalent: got %d, want %d", saving.MetallicEquivalent, metallicEq)
	}
	if saving.ActualCurrency != actual {
		t.Errorf("ActualCurrency: got %d, want %d", saving.ActualCurrency, actual)
	}
}

func TestComputeCirculationCostSaving_Zero(t *testing.T) {
	t.Parallel()

	saving := ComputeCirculationCostSaving(0, 0, "")
	if saving.SavingAmount != 0 {
		t.Errorf("expected 0, got %d", saving.SavingAmount)
	}
}

// ── ID format and uniqueness ──────────────────────────────────────────────────

func TestNewStockCompanyID_Format(t *testing.T) {
	t.Parallel()

	id := NewStockCompanyID()
	if len(string(id)) != 24 {
		t.Errorf("expected 24 hex chars, got %d: %q", len(string(id)), id)
	}
	if strings.TrimLeft(string(id), "0123456789abcdef") != "" {
		t.Errorf("expected lowercase hex, got %q", id)
	}
}

func TestNewStockCompanyID_Unique(t *testing.T) {
	t.Parallel()

	a, b := NewStockCompanyID(), NewStockCompanyID()
	if a == b {
		t.Error("two IDs should not be equal")
	}
}

func TestNewCooperativeFactoryID_Format(t *testing.T) {
	t.Parallel()

	id := NewCooperativeFactoryID()
	if len(string(id)) != 24 {
		t.Errorf("expected 24 hex chars, got %d: %q", len(string(id)), id)
	}
}

func TestNewCooperativeFactoryID_Unique(t *testing.T) {
	t.Parallel()

	a, b := NewCooperativeFactoryID(), NewCooperativeFactoryID()
	if a == b {
		t.Error("two IDs should not be equal")
	}
}

// ── StockCompanyID.IsZero ─────────────────────────────────────────────────────

func TestStockCompanyID_IsZero(t *testing.T) {
	t.Parallel()

	if !StockCompanyID("").IsZero() {
		t.Error("empty should be zero")
	}
	if NewStockCompanyID().IsZero() {
		t.Error("fresh ID should not be zero")
	}
}

// ── CooperativeFactoryID.IsZero ───────────────────────────────────────────────

func TestCooperativeFactoryID_IsZero(t *testing.T) {
	t.Parallel()

	if !CooperativeFactoryID("").IsZero() {
		t.Error("empty should be zero")
	}
	if NewCooperativeFactoryID().IsZero() {
		t.Error("fresh ID should not be zero")
	}
}
