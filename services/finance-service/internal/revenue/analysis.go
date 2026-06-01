// analysis.go — Vol. III Ch. 49 types: Concerning the Analysis of the Process
// of Production. The value of the annual product resolves into c + v + s.
// Only v + s is newly created value; the constant capital c is transferred and
// must be reproduced, not consumed as revenue. Adam Smith's error is to resolve
// the whole product into wages + profit + rent, omitting c. The surplus-value s
// is the limit from which profit (of enterprise + interest) and rent are drawn.
package revenue

import (
	"fmt"
	"time"
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

// ComputeTotalValue returns the value of the annual product: c + v + s.
func ComputeTotalValue(constant, variable, surplus int64) int64 {
	return constant + variable + surplus
}

// ComputeNewValue returns the value newly added by living labour: v + s.
func ComputeNewValue(variable, surplus int64) int64 {
	return variable + surplus
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
