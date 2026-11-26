package store

import (
	"fmt"
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

// 00111 corrects the Vol. III Ch. 7 seed to Marx's text (issue #450). He states
// the proportional case of the Rodbertus refutation as "C, or nC, or C/n, e.g.,
// 1,000, or 2,000, or 500, and the rate of profit being 20%": profit 200, 400
// and 100. 00014 had an invented £3,000 / £600 growing 2.5x instead, and a
// composition comparison (row 0702) that appears nowhere in the chapter.
func TestMigration00111_Ch7SeedFidelity(t *testing.T) {
	t.Parallel()
	raw, err := migrationsFS.ReadFile("migrations/00111_v3_ch07_seed_fidelity.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	up, down, ok := strings.Cut(string(raw), "-- +goose Down")
	if !ok {
		t.Fatal("migration has no Down section")
	}

	// Marx's base capital of 1,000 at 20%, doubled (nC = 2,000, profit 400) and
	// halved (C/n = 500, profit 100). The rates the migration stores are the ones
	// the domain derives, not independently typed numbers.
	const capital, profitMass = 1000, 200
	for _, tc := range []struct {
		id      string
		factor  int64
		wantC   int64
		wantPro int64
	}{
		{"5eed000000000000000704", 200, 2000, 400},
		{"5eed000000000000000706", 50, 500, 100},
	} {
		change := profit.CapitalMagnitudeChange{
			Kind:            profit.KindProportionalChange,
			OriginalCapital: capital,
			OriginalProfit:  profitMass,
			Factor:          tc.factor,
		}
		a := profit.ComputeMagnitudeChangeAnalysis(change)
		if a.OldRate != 2000 || a.NewRate != 2000 || !a.RateUnchanged {
			t.Fatalf("%s: domain gives old=%d new=%d unchanged=%v, want 2000/2000/true",
				tc.id, a.OldRate, a.NewRate, a.RateUnchanged)
		}
		if got := int64(change.OriginalCapital) * tc.factor / 100; got != tc.wantC {
			t.Errorf("%s: re-valued capital = %d, want Marx's %d", tc.id, got, tc.wantC)
		}
		if got := int64(change.OriginalProfit) * tc.factor / 100; got != tc.wantPro {
			t.Errorf("%s: re-valued profit = %d, want Marx's %d", tc.id, got, tc.wantPro)
		}
	}

	// Row 0704 keeps its derived columns and only has its inputs corrected; row 0706
	// is inserted with the rate columns the domain just produced.
	for _, want := range []string{
		"DELETE FROM composition_effects WHERE id = '5eed000000000000000702'",
		"UPDATE magnitude_changes SET original_capital = 1000, original_profit = 200, factor = 200 WHERE id = '5eed000000000000000704'",
		fmt.Sprintf("('5eed000000000000000706', 'proportional_change', 1000, 200, 50, %d, %d, 1,", 2000, 2000),
	} {
		if !strings.Contains(up, want) {
			t.Errorf("Up missing %q", want)
		}
	}
	if strings.Contains(up, "3000") || strings.Contains(up, "4000") {
		t.Error("Up must not reintroduce the invented £3,000 / £4,000 magnitudes")
	}

	// Down restores 00014's rows exactly.
	for _, want := range []string{
		"DELETE FROM magnitude_changes WHERE id = '5eed000000000000000706'",
		"UPDATE magnitude_changes SET original_capital = 3000, original_profit = 600, factor = 250 WHERE id = '5eed000000000000000704'",
		"('5eed000000000000000702', 1000, 1000, 4000, 1000, 1000, 9000, 2000, 1000, 1000,",
	} {
		if !strings.Contains(down, want) {
			t.Errorf("Down missing %q", want)
		}
	}
	if strings.Count(up, "UPDATE") != strings.Count(down, "UPDATE") {
		t.Errorf("Up has %d UPDATEs but Down has %d", strings.Count(up, "UPDATE"), strings.Count(down, "UPDATE"))
	}
}
