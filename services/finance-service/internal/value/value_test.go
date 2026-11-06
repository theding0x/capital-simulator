package value_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/value"
)

// Compile-time guards: profit and avgprofit alias the shared types, so a
// value built in one is assignable to the other with no conversion. If either
// package re-declares its own c, v or s, this file stops compiling.
var (
	_ value.LabourMinutes   = profit.LabourMinutes(0)
	_ value.LabourMinutes   = avgprofit.LabourMinutes(0)
	_ value.ConstantCapital = avgprofit.ConstantCapital(profit.ConstantCapital(0))
	_ value.VariableCapital = avgprofit.VariableCapital(profit.VariableCapital(0))
	_ profit.SurplusValue   = avgprofit.SurplusValue(0)
)

// Ch. 49 aggregate: c=6000, v=1500, s=1500. The s that Vol. I produces is the
// same type Vol. III's revenue analysis divides into profit, interest and rent.
func TestSharedTypesFlowAcrossPackages(t *testing.T) {
	t.Parallel()
	var (
		c profit.ConstantCapital    = 6000
		v avgprofit.VariableCapital = 1500
		s profit.SurplusValue       = 1500
	)
	if got := revenue.ComputeTotalValue(c, v, s); got != 9000 {
		t.Errorf("total value = %d, want 9000", got)
	}
	if got := revenue.ComputeNewValue(v, s); got != 3000 {
		t.Errorf("new value = %d, want 3000", got)
	}
}
