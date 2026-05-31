// Package rent — tests for Ch. 39 Differential Rent I domain types and
// ComputeDifferentialRentI. Fixtures drawn from Marx, Vol. III Ch. 39,
// Table I: A=100 qtrs out / £2.50 capital / £3 price; B=200; C=300; D=400.
// All monetary amounts in pence; quantities in 1/100 quarter.
package rent

import (
	"testing"
)

// ── SoilGrade ─────────────────────────────────────────────────────────────────

func TestSoilGradeIsValid(t *testing.T) {
	t.Parallel()
	valid := []SoilGrade{SoilGradeA, SoilGradeB, SoilGradeC, SoilGradeD}
	for _, g := range valid {
		if !g.IsValid() {
			t.Errorf("SoilGrade(%d).IsValid() = false, want true", g)
		}
	}
	invalid := []SoilGrade{0, 5}
	for _, g := range invalid {
		if g.IsValid() {
			t.Errorf("SoilGrade(%d).IsValid() = true, want false", g)
		}
	}
}

func TestSoilGradeString(t *testing.T) {
	t.Parallel()
	cases := []struct {
		grade SoilGrade
		want  string
	}{
		{SoilGradeA, "A"},
		{SoilGradeB, "B"},
		{SoilGradeC, "C"},
		{SoilGradeD, "D"},
	}
	for _, tc := range cases {
		if got := tc.grade.String(); got != tc.want {
			t.Errorf("SoilGrade(%d).String() = %q, want %q", tc.grade, got, tc.want)
		}
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// marxTableIEntries returns the 4 soil entries from Marx Vol. III Ch. 39
// Table I: capital £2.50 (250p), price £3/qtr (300p), output A=1qtr(100),
// B=2qtrs(200), C=3qtrs(300), D=4qtrs(400).
func marxTableIEntries() []DifferentialRentIEntry {
	return []DifferentialRentIEntry{
		{Grade: SoilGradeA, CapitalAdvancedBP: 250, OutputQuarters: 100, PriceOfProductionBP: 300},
		{Grade: SoilGradeB, CapitalAdvancedBP: 250, OutputQuarters: 200, PriceOfProductionBP: 300},
		{Grade: SoilGradeC, CapitalAdvancedBP: 250, OutputQuarters: 300, PriceOfProductionBP: 300},
		{Grade: SoilGradeD, CapitalAdvancedBP: 250, OutputQuarters: 400, PriceOfProductionBP: 300},
	}
}

// ── ComputeDifferentialRentI ──────────────────────────────────────────────────

func TestComputeDifferentialRentI_MarxTableI(t *testing.T) {
	t.Parallel()
	tbl, entries := ComputeDifferentialRentI(marxTableIEntries())

	if tbl.TotalRentalBP != 1800 {
		t.Errorf("TotalRentalBP = %d, want 1800", tbl.TotalRentalBP)
	}
	if tbl.RegulatingGrade != SoilGradeA {
		t.Errorf("RegulatingGrade = %v, want SoilGradeA", tbl.RegulatingGrade)
	}
	if len(entries) != 4 {
		t.Fatalf("len(entries) = %d, want 4", len(entries))
	}

	type wantRow struct {
		grade   SoilGrade
		surplus int64
		rent    int64
	}
	wants := []wantRow{
		{SoilGradeA, 0, 0},
		{SoilGradeB, 100, 300},
		{SoilGradeC, 200, 600},
		{SoilGradeD, 300, 900},
	}
	// Build a map by grade for order-independent lookup.
	byGrade := make(map[SoilGrade]DifferentialRentIEntry, len(entries))
	for _, e := range entries {
		byGrade[e.Grade] = e
	}
	for _, w := range wants {
		e, ok := byGrade[w.grade]
		if !ok {
			t.Errorf("no entry for grade %v in results", w.grade)
			continue
		}
		if e.SurplusProductQtrs != w.surplus {
			t.Errorf("grade %v SurplusProductQtrs = %d, want %d", w.grade, e.SurplusProductQtrs, w.surplus)
		}
		if e.MoneyRentBP != w.rent {
			t.Errorf("grade %v MoneyRentBP = %d, want %d", w.grade, e.MoneyRentBP, w.rent)
		}
	}
}

func TestComputeDifferentialRentI_Deterministic(t *testing.T) {
	t.Parallel()
	// Shuffle: D, B, A, C — different order from Table I.
	shuffled := []DifferentialRentIEntry{
		{Grade: SoilGradeD, CapitalAdvancedBP: 250, OutputQuarters: 400, PriceOfProductionBP: 300},
		{Grade: SoilGradeB, CapitalAdvancedBP: 250, OutputQuarters: 200, PriceOfProductionBP: 300},
		{Grade: SoilGradeA, CapitalAdvancedBP: 250, OutputQuarters: 100, PriceOfProductionBP: 300},
		{Grade: SoilGradeC, CapitalAdvancedBP: 250, OutputQuarters: 300, PriceOfProductionBP: 300},
	}
	tbl, entries := ComputeDifferentialRentI(shuffled)

	if tbl.TotalRentalBP != 1800 {
		t.Errorf("TotalRentalBP = %d, want 1800 (order-independent)", tbl.TotalRentalBP)
	}
	if tbl.RegulatingGrade != SoilGradeA {
		t.Errorf("RegulatingGrade = %v, want SoilGradeA", tbl.RegulatingGrade)
	}
	var dEntry DifferentialRentIEntry
	for _, e := range entries {
		if e.Grade == SoilGradeD {
			dEntry = e
			break
		}
	}
	if dEntry.MoneyRentBP != 900 {
		t.Errorf("grade D MoneyRentBP = %d, want 900", dEntry.MoneyRentBP)
	}
}

func TestComputeDifferentialRentI_AllEqualOutput(t *testing.T) {
	t.Parallel()
	entries := []DifferentialRentIEntry{
		{Grade: SoilGradeA, CapitalAdvancedBP: 250, OutputQuarters: 200, PriceOfProductionBP: 300},
		{Grade: SoilGradeB, CapitalAdvancedBP: 250, OutputQuarters: 200, PriceOfProductionBP: 300},
		{Grade: SoilGradeC, CapitalAdvancedBP: 250, OutputQuarters: 200, PriceOfProductionBP: 300},
	}
	tbl, computed := ComputeDifferentialRentI(entries)

	if tbl.TotalRentalBP != 0 {
		t.Errorf("TotalRentalBP = %d, want 0 (equal outputs → no surplus)", tbl.TotalRentalBP)
	}
	for _, e := range computed {
		if e.SurplusProductQtrs != 0 {
			t.Errorf("grade %v SurplusProductQtrs = %d, want 0", e.Grade, e.SurplusProductQtrs)
		}
		if e.MoneyRentBP != 0 {
			t.Errorf("grade %v MoneyRentBP = %d, want 0", e.Grade, e.MoneyRentBP)
		}
	}
}

func TestComputeDifferentialRentI_Empty(t *testing.T) {
	t.Parallel()
	tbl, entries := ComputeDifferentialRentI(nil)
	if tbl.TotalRentalBP != 0 || tbl.RegulatingGrade != 0 {
		t.Errorf("empty input: table = %+v, want zero value", tbl)
	}
	if entries != nil {
		t.Errorf("empty input: entries = %v, want nil", entries)
	}

	tbl2, entries2 := ComputeDifferentialRentI([]DifferentialRentIEntry{})
	if tbl2.TotalRentalBP != 0 {
		t.Errorf("empty slice: TotalRentalBP = %d, want 0", tbl2.TotalRentalBP)
	}
	if entries2 != nil {
		t.Errorf("empty slice: entries = %v, want nil", entries2)
	}
}

// ── ID constructors ───────────────────────────────────────────────────────────

func TestCh39NewIDsDistinct(t *testing.T) {
	t.Parallel()

	// DifferentialRentITableID: 24-char hex, two calls differ.
	id1 := NewDifferentialRentITableID()
	id2 := NewDifferentialRentITableID()
	if len(string(id1)) != 24 {
		t.Errorf("NewDifferentialRentITableID len = %d, want 24", len(string(id1)))
	}
	if id1 == id2 {
		t.Error("NewDifferentialRentITableID: two calls returned same value")
	}

	// DifferentialRentIEntryID: 24-char hex, two calls differ.
	eid1 := NewDifferentialRentIEntryID()
	eid2 := NewDifferentialRentIEntryID()
	if len(string(eid1)) != 24 {
		t.Errorf("NewDifferentialRentIEntryID len = %d, want 24", len(string(eid1)))
	}
	if eid1 == eid2 {
		t.Error("NewDifferentialRentIEntryID: two calls returned same value")
	}

	// LocationRentFactorID: 24-char hex, two calls differ.
	lid1 := NewLocationRentFactorID()
	lid2 := NewLocationRentFactorID()
	if len(string(lid1)) != 24 {
		t.Errorf("NewLocationRentFactorID len = %d, want 24", len(string(lid1)))
	}
	if lid1 == lid2 {
		t.Error("NewLocationRentFactorID: two calls returned same value")
	}
}
