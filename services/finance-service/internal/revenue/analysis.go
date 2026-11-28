// analysis.go — Vol. III Ch. 49 types: Concerning the Analysis of the Process
// of Production. The value of the annual product resolves into c + v + s.
// Only v + s is newly created value; the constant capital c is transferred and
// must be reproduced, not consumed as revenue. Adam Smith's error is to resolve
// the whole product into wages + profit + rent, omitting c. The surplus-value s
// is the limit from which profit (of enterprise + interest) and rent are drawn.
package revenue

import (
	"errors"
	"fmt"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/value"
)

// AnnualValueCompositionID identifies an annual-value-composition record.
type AnnualValueCompositionID string

// NewAnnualValueCompositionID returns a fresh 96-bit hex ID.
func NewAnnualValueCompositionID() AnnualValueCompositionID {
	return AnnualValueCompositionID(newID())
}

// AnnualValueComposition decomposes the value of the total annual product into
// constant capital (transferred), variable capital (reproduced as wages), and
// surplus-value (newly created). TotalValue = c + v + s; NewValue = v + s.
type AnnualValueComposition struct {
	ID                           AnnualValueCompositionID `json:"id"`
	ConstantCapitalLabourMinutes int64                    `json:"constant_capital_labour_minutes"`
	VariableCapitalLabourMinutes int64                    `json:"variable_capital_labour_minutes"`
	SurplusValueLabourMinutes    int64                    `json:"surplus_value_labour_minutes"`
	TotalValueLabourMinutes      int64                    `json:"total_value_labour_minutes"`
	NewValueLabourMinutes        int64                    `json:"new_value_labour_minutes"`
	CreatedAt                    time.Time                `json:"created_at"`
}

// ErrCompositionNotConserved is returned when a composition violates L-P2:
// its components are negative, or W != c+v+s, or New != v+s.
var ErrCompositionNotConserved = errors.New("revenue: annual value composition does not conserve value (want W = c+v+s and New = v+s, all parts >= 0)")

// Validate asserts the Ch. 49 identities W = c + v + s and New = v + s, with
// c, v, s >= 0. It guards both construction and read, so a stored row that has
// drifted from the identities is rejected rather than served.
func (c AnnualValueComposition) Validate() error {
	if c.ConstantCapitalLabourMinutes < 0 || c.VariableCapitalLabourMinutes < 0 || c.SurplusValueLabourMinutes < 0 {
		return fmt.Errorf("%w: c, v, and s must be >= 0", ErrCompositionNotConserved)
	}
	if want := int64(ComputeTotalValue(
		value.ConstantCapital(c.ConstantCapitalLabourMinutes),
		value.VariableCapital(c.VariableCapitalLabourMinutes),
		value.SurplusValue(c.SurplusValueLabourMinutes),
	)); c.TotalValueLabourMinutes != want {
		return fmt.Errorf("%w: total value %d, want c+v+s = %d", ErrCompositionNotConserved, c.TotalValueLabourMinutes, want)
	}
	if want := int64(ComputeNewValue(
		value.VariableCapital(c.VariableCapitalLabourMinutes),
		value.SurplusValue(c.SurplusValueLabourMinutes),
	)); c.NewValueLabourMinutes != want {
		return fmt.Errorf("%w: new value %d, want v+s = %d", ErrCompositionNotConserved, c.NewValueLabourMinutes, want)
	}
	return nil
}

// RevenueLimitID identifies a revenue-limit record.
type RevenueLimitID string

// NewRevenueLimitID returns a fresh 96-bit hex ID.
func NewRevenueLimitID() RevenueLimitID {
	return RevenueLimitID(newID())
}

// RevenueLimit records the ceiling on distributable revenue implied by a
// composition: revenue cannot exceed the new value (v + s), and the constant
// capital must be reproduced, not consumed.
type RevenueLimit struct {
	ID                         RevenueLimitID `json:"id"`
	CompositionID              string         `json:"composition_id"`
	MaxDistributableRevenueLM  int64          `json:"max_distributable_revenue_lm"`
	ConstantCapitalReplacement int64          `json:"constant_capital_replacement"`
	CreatedAt                  time.Time      `json:"created_at"`
}

// SurplusValuePartitionID identifies a surplus-value-partition record.
type SurplusValuePartitionID string

// NewSurplusValuePartitionID returns a fresh 96-bit hex ID.
func NewSurplusValuePartitionID() SurplusValuePartitionID {
	return SurplusValuePartitionID(newID())
}

// SurplusValuePartition divides the surplus-value s into profit of enterprise,
// interest, ground-rent, and any uncaptured residual. The parts sum to the total,
// and the total may not exceed the surplus-value the composition produced.
type SurplusValuePartition struct {
	ID                  SurplusValuePartitionID `json:"id"`
	CompositionID       string                  `json:"composition_id"`
	TotalSurplusValueLM int64                   `json:"total_surplus_value_lm"`
	ProfitEnterpriseLM  int64                   `json:"profit_enterprise_lm"`
	InterestLM          int64                   `json:"interest_lm"`
	GroundRentLM        int64                   `json:"ground_rent_lm"`
	ResidualSurplusLM   int64                   `json:"residual_surplus_lm"`
	CreatedAt           time.Time               `json:"created_at"`
}

// ComputeTotalValue returns the value of the annual product: c + v + s. The
// parameters are the shared Vol. I types, so a Vol. I s cannot be swapped for v.
func ComputeTotalValue(constant value.ConstantCapital, variable value.VariableCapital, surplus value.SurplusValue) value.LabourMinutes {
	return value.LabourMinutes(constant) + value.LabourMinutes(variable) + value.LabourMinutes(surplus)
}

// ComputeNewValue returns the value newly added by living labour: v + s.
func ComputeNewValue(variable value.VariableCapital, surplus value.SurplusValue) value.LabourMinutes {
	return value.LabourMinutes(variable) + value.LabourMinutes(surplus)
}

// DeriveRevenueLimit builds the revenue limit implied by a composition: the
// maximum distributable revenue equals the new value (v + s), and the constant
// capital must be reproduced.
func DeriveRevenueLimit(comp AnnualValueComposition) RevenueLimit {
	return RevenueLimit{
		CompositionID:              string(comp.ID),
		MaxDistributableRevenueLM:  comp.NewValueLabourMinutes,
		ConstantCapitalReplacement: comp.ConstantCapitalLabourMinutes,
	}
}

// SumSurplusPartition returns the sum of the four parts of a surplus-value partition.
func SumSurplusPartition(p SurplusValuePartition) int64 {
	return p.ProfitEnterpriseLM + p.InterestLM + p.GroundRentLM + p.ResidualSurplusLM
}

// ValidateRevenueLimit checks that a surplus-value partition is consistent with the
// composition it divides: its parts must sum to its declared total, and that total
// may not exceed the surplus-value actually produced. Smith's error is precisely to
// resolve the whole product into revenue, omitting the constant capital.
func ValidateRevenueLimit(comp AnnualValueComposition, partition SurplusValuePartition) error {
	if got := SumSurplusPartition(partition); got != partition.TotalSurplusValueLM {
		return fmt.Errorf("surplus partition parts sum to %d, want declared total %d", got, partition.TotalSurplusValueLM)
	}
	if partition.TotalSurplusValueLM > comp.SurplusValueLabourMinutes {
		return fmt.Errorf("partitioned surplus-value %d exceeds produced surplus-value %d", partition.TotalSurplusValueLM, comp.SurplusValueLabourMinutes)
	}
	return nil
}

// ErrNegativeRevenueClaim is returned when a claimed revenue is negative: a
// claim that the product resolves into less than nothing is not Smith's error,
// it is malformed input.
var ErrNegativeRevenueClaim = errors.New("revenue: claimed revenue must be >= 0")

// SmithsAddingUpError records Adam Smith's "adding-up" error: resolving the value
// of the product into revenue (wages + profit + rent) without accounting for the
// constant capital c, which is transferred, not newly created, and must be
// replaced. Only the new value v + s can be revenue; anything claimed beyond it
// is constant capital treated as if it needed no replacement.
//
// On Marx's aggregate fixture (c = 6000, v = 1500, s = 1500), resolving the whole
// product, 9000, into revenue claims 9000 against a new value of 3000: an error
// of 6000, which is exactly c.
type SmithsAddingUpError struct {
	CompositionID string `json:"composition_id"`
	// ClaimedRevenueLM is the value Smith resolves into wages + profit + rent.
	ClaimedRevenueLM int64 `json:"claimed_revenue_lm"`
	// NewValueLM is v + s, the most that can be distributed as revenue.
	NewValueLM int64 `json:"new_value_lm"`
	// OmittedConstantCapitalLM is ClaimedRevenueLM - NewValueLM: the constant
	// capital consumed as revenue. For Smith's own claim, the whole product, it
	// equals c.
	OmittedConstantCapitalLM int64 `json:"omitted_constant_capital_lm"`
}

// Error implements the error interface.
func (e *SmithsAddingUpError) Error() string {
	return fmt.Sprintf(
		"revenue: Smith's adding-up error: %d claimed as revenue exceeds the new value v+s = %d; %d of constant capital omitted from replacement",
		e.ClaimedRevenueLM, e.NewValueLM, e.OmittedConstantCapitalLM,
	)
}

// CheckRevenueResolution tests a claim that claimedRevenueLM of the annual
// product resolves into revenue. Revenue cannot exceed the new value v + s, so a
// claim beyond it returns a *SmithsAddingUpError carrying the size of the
// omission. The composition is validated first, so the comparison is made
// against a new value that actually equals v + s.
func CheckRevenueResolution(comp AnnualValueComposition, claimedRevenueLM int64) error {
	if err := comp.Validate(); err != nil {
		return err
	}
	if claimedRevenueLM < 0 {
		return fmt.Errorf("%w: got %d", ErrNegativeRevenueClaim, claimedRevenueLM)
	}
	if claimedRevenueLM > comp.NewValueLabourMinutes {
		return &SmithsAddingUpError{
			CompositionID:            string(comp.ID),
			ClaimedRevenueLM:         claimedRevenueLM,
			NewValueLM:               comp.NewValueLabourMinutes,
			OmittedConstantCapitalLM: claimedRevenueLM - comp.NewValueLabourMinutes,
		}
	}
	return nil
}
