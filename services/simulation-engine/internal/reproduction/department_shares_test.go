package reproduction_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// TestComputeDepartmentSharesBP verifies the basis-point partition of advanced
// capital (c+v) across the two departments (issue #215). The shares always sum
// to exactly 10000 except when the total is zero.
func TestComputeDepartmentSharesBP(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		deptIAdvanced  int64
		deptIIAdvanced int64
		wantI          int64
		wantII         int64
	}{
		{
			// Marx's canonical Ch. 20 fixture (period 1865), pence:
			//   Dept I  4 000 000c + 1 000 000v = 5 000 000 advanced
			//   Dept II 2 000 000c +   500 000v = 2 500 000 advanced
			name:           "marx 1865 fixture two-to-one",
			deptIAdvanced:  5_000_000,
			deptIIAdvanced: 2_500_000,
			wantI:          6667,
			wantII:         3333,
		},
		{
			name:           "equal departments split evenly",
			deptIAdvanced:  1000,
			deptIIAdvanced: 1000,
			wantI:          5000,
			wantII:         5000,
		},
		{
			name:           "dept I only takes the whole share",
			deptIAdvanced:  7500,
			deptIIAdvanced: 0,
			wantI:          10000,
			wantII:         0,
		},
		{
			name:           "zero total yields zero shares",
			deptIAdvanced:  0,
			deptIIAdvanced: 0,
			wantI:          0,
			wantII:         0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			gotI, gotII := reproduction.ComputeDepartmentSharesBP(tc.deptIAdvanced, tc.deptIIAdvanced)
			if gotI != tc.wantI || gotII != tc.wantII {
				t.Fatalf("ComputeDepartmentSharesBP(%d, %d) = (%d, %d); want (%d, %d)",
					tc.deptIAdvanced, tc.deptIIAdvanced, gotI, gotII, tc.wantI, tc.wantII)
			}
			if tc.deptIAdvanced+tc.deptIIAdvanced > 0 && gotI+gotII != 10000 {
				t.Fatalf("shares must sum to 10000, got %d + %d = %d", gotI, gotII, gotI+gotII)
			}
		})
	}
}

// TestDepartmentalCapitalAdvancedPence confirms AdvancedPence is c + v and
// excludes the surplus the department produces.
func TestDepartmentalCapitalAdvancedPence(t *testing.T) {
	t.Parallel()
	d := reproduction.DepartmentalCapital{
		ConstantPence: 4_000_000,
		VariablePence: 1_000_000,
		SurplusPence:  1_000_000,
		TotalPence:    6_000_000,
	}
	if got := d.AdvancedPence(); got != 5_000_000 {
		t.Fatalf("AdvancedPence() = %d; want 5000000 (c+v, excludes surplus)", got)
	}
}
