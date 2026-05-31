package credit

import (
	"testing"
)

// TestNewProfitDivision_BorrowedCapitalFixture covers Marx's canonical Ch. 23 fixture:
// functioning capitalist borrows at 5% (500 bp); total profit = 20% (2000 bp);
// profit of enterprise = 1500 bp.
func TestNewProfitDivision_BorrowedCapitalFixture(t *testing.T) {
	t.Parallel()

	pd, err := NewProfitDivision(2000, 500, FormSeparated, WagesOfSuperintendence{LabourMinutes: 0}, 0)
	if err != nil {
		t.Fatalf("NewProfitDivision: %v", err)
	}
	if pd.TotalProfitBP != 2000 {
		t.Errorf("TotalProfitBP = %d, want 2000", pd.TotalProfitBP)
	}
	if pd.InterestBP != 500 {
		t.Errorf("InterestBP = %d, want 500", pd.InterestBP)
	}
	if pd.ProfitOfEnterpriseBP != 1500 {
		t.Errorf("ProfitOfEnterpriseBP = %d, want 1500", pd.ProfitOfEnterpriseBP)
	}
	// Partition invariant.
	if pd.TotalProfitBP != pd.InterestBP+pd.ProfitOfEnterpriseBP {
		t.Errorf("partition violated: %d != %d + %d",
			pd.TotalProfitBP, pd.InterestBP, pd.ProfitOfEnterpriseBP)
	}
	if pd.OwnershipForm != FormSeparated {
		t.Errorf("OwnershipForm = %d, want FormSeparated", pd.OwnershipForm)
	}
}

// TestNewProfitDivision_OwnCapitalFixture covers the own-capital case:
// capitalist employs own capital; interest = 0; profit of enterprise = total.
func TestNewProfitDivision_OwnCapitalFixture(t *testing.T) {
	t.Parallel()

	pd, err := NewProfitDivision(2000, 0, FormOwnerOperator, WagesOfSuperintendence{LabourMinutes: 480}, 0)
	if err != nil {
		t.Fatalf("NewProfitDivision: %v", err)
	}
	if pd.ProfitOfEnterpriseBP != 2000 {
		t.Errorf("ProfitOfEnterpriseBP = %d, want 2000", pd.ProfitOfEnterpriseBP)
	}
	if pd.InterestBP != 0 {
		t.Errorf("InterestBP = %d, want 0", pd.InterestBP)
	}
	if pd.TotalProfitBP != pd.InterestBP+pd.ProfitOfEnterpriseBP {
		t.Errorf("partition violated: %d != %d + %d",
			pd.TotalProfitBP, pd.InterestBP, pd.ProfitOfEnterpriseBP)
	}
	if pd.OwnershipForm != FormOwnerOperator {
		t.Errorf("OwnershipForm = %d, want FormOwnerOperator", pd.OwnershipForm)
	}
}

// TestNewProfitDivision_NonPositiveTotal asserts ErrNonPositiveTotalProfit
// when total_profit_bp <= 0.
func TestNewProfitDivision_NonPositiveTotal(t *testing.T) {
	t.Parallel()

	_, err := NewProfitDivision(0, 0, FormOwnerOperator, WagesOfSuperintendence{}, 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestNewProfitDivision_NegativeInterest asserts ErrNegativeInterest.
func TestNewProfitDivision_NegativeInterest(t *testing.T) {
	t.Parallel()

	_, err := NewProfitDivision(2000, -1, FormSeparated, WagesOfSuperintendence{}, 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestNewProfitDivision_InterestEqualsTotal asserts ErrInterestNotLessThanTotal
// when interest_bp >= total_profit_bp.
func TestNewProfitDivision_InterestEqualsTotal(t *testing.T) {
	t.Parallel()

	_, err := NewProfitDivision(2000, 2000, FormSeparated, WagesOfSuperintendence{}, 0)
	if err == nil {
		t.Fatal("expected error when interest == total, got nil")
	}
}

// TestNewProfitDivision_InterestExceedsTotal asserts error when interest > total.
func TestNewProfitDivision_InterestExceedsTotal(t *testing.T) {
	t.Parallel()

	_, err := NewProfitDivision(2000, 2001, FormSeparated, WagesOfSuperintendence{}, 0)
	if err == nil {
		t.Fatal("expected error when interest > total, got nil")
	}
}

// TestNewProfitDivision_InvalidOwnershipForm asserts ErrInvalidOwnershipForm.
func TestNewProfitDivision_InvalidOwnershipForm(t *testing.T) {
	t.Parallel()

	_, err := NewProfitDivision(2000, 500, CapitalOwnershipForm(99), WagesOfSuperintendence{}, 0)
	if err == nil {
		t.Fatal("expected error for invalid ownership form, got nil")
	}
}

// TestNewProfitOfEnterprise_Derived verifies the value-object derivation.
func TestNewProfitOfEnterprise_Derived(t *testing.T) {
	t.Parallel()

	poe, err := NewProfitOfEnterprise(2000, 500)
	if err != nil {
		t.Fatalf("NewProfitOfEnterprise: %v", err)
	}
	if poe.ProfitOfEnterpriseBP != 1500 {
		t.Errorf("ProfitOfEnterpriseBP = %d, want 1500", poe.ProfitOfEnterpriseBP)
	}
}

// TestNewProfitOfEnterprise_ZeroInterest verifies own-capital case.
func TestNewProfitOfEnterprise_ZeroInterest(t *testing.T) {
	t.Parallel()

	poe, err := NewProfitOfEnterprise(2000, 0)
	if err != nil {
		t.Fatalf("NewProfitOfEnterprise: %v", err)
	}
	if poe.ProfitOfEnterpriseBP != 2000 {
		t.Errorf("ProfitOfEnterpriseBP = %d, want 2000", poe.ProfitOfEnterpriseBP)
	}
}

// TestCapitalOwnershipForm_IsValid verifies both constants are valid and others are not.
func TestCapitalOwnershipForm_IsValid(t *testing.T) {
	t.Parallel()

	if !FormOwnerOperator.IsValid() {
		t.Error("FormOwnerOperator.IsValid() = false, want true")
	}
	if !FormSeparated.IsValid() {
		t.Error("FormSeparated.IsValid() = false, want true")
	}
	if CapitalOwnershipForm(0).IsValid() {
		t.Error("CapitalOwnershipForm(0).IsValid() = true, want false")
	}
	if CapitalOwnershipForm(3).IsValid() {
		t.Error("CapitalOwnershipForm(3).IsValid() = true, want false")
	}
}

// TestProfitDivisionID_IsZero verifies the zero-value sentinel.
func TestProfitDivisionID_IsZero(t *testing.T) {
	t.Parallel()

	var id ProfitDivisionID
	if !id.IsZero() {
		t.Error("empty ProfitDivisionID.IsZero() = false, want true")
	}
	id = NewProfitDivisionID()
	if id.IsZero() {
		t.Error("NewProfitDivisionID().IsZero() = true, want false")
	}
}
