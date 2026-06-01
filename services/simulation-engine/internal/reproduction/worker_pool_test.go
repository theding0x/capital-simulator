package reproduction_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// TestComputeLabourForceReproduction covers the issue #220 assertion: a period's
// wage bill must cover worker_pool_size × subsistence_basket; a zero pool means
// the check was not requested (vacuously reproduced).
func TestComputeLabourForceReproduction(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name        string
		wageBill    int64
		workers     int64
		basket      int64
		wantCovered bool
		wantStatus  string
	}{
		{"untracked when pool zero", 0, 0, 500, true, reproduction.TickStatusReproduced},
		{"covered exactly", 1500, 3, 500, true, reproduction.TickStatusReproduced},
		{"covered with surplus", 2000, 3, 500, true, reproduction.TickStatusReproduced},
		{"shortfall by one", 1499, 3, 500, false, reproduction.TickStatusLabourShortfall},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			covered, status := reproduction.ComputeLabourForceReproduction(tc.wageBill, tc.workers, tc.basket)
			if covered != tc.wantCovered || status != tc.wantStatus {
				t.Fatalf("ComputeLabourForceReproduction(%d, %d, %d) = (%v, %q); want (%v, %q)",
					tc.wageBill, tc.workers, tc.basket, covered, status, tc.wantCovered, tc.wantStatus)
			}
		})
	}
}

// TestSimpleReproductionScheme_WageBillPence sums the two departments' variable
// capital (the wage bill). Marx's 1865 fixture: I.v=1000 + II.v=500 = 1500.
func TestSimpleReproductionScheme_WageBillPence(t *testing.T) {
	t.Parallel()
	deptI := deptI1865(0)
	deptII := deptII1865(0)
	s := reproduction.SimpleReproductionScheme{DepartmentI: &deptI, DepartmentII: &deptII}
	if got := s.WageBillPence(); got != 1500 {
		t.Fatalf("WageBillPence: got %d, want 1500 (I.v 1000 + II.v 500)", got)
	}

	// A scheme missing a department contributes only what is present.
	partial := reproduction.SimpleReproductionScheme{DepartmentI: &deptI}
	if got := partial.WageBillPence(); got != 1000 {
		t.Fatalf("partial WageBillPence: got %d, want 1000", got)
	}
}
