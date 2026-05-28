package circulation_test

import (
	"strings"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

// Vol. II Ch. 10–11 — Theories of Fixed and Circulating Capital:
// The Physiocrats, Adam Smith, and Ricardo.

func TestNewEconomistAttributionID_Format(t *testing.T) {
	t.Parallel()
	id := circulation.NewEconomistAttributionID()
	if id.IsZero() {
		t.Fatal("expected non-empty ID")
	}
	s := string(id)
	if len(s) != 24 {
		t.Fatalf("expected 24-char hex ID, got %d chars: %q", len(s), s)
	}
	for _, c := range s {
		if !strings.ContainsRune("0123456789abcdef", c) {
			t.Fatalf("non-hex character %q in ID %q", c, s)
		}
	}
}

func TestNewEconomistAttributionID_Unique(t *testing.T) {
	t.Parallel()
	a := circulation.NewEconomistAttributionID()
	b := circulation.NewEconomistAttributionID()
	if a == b {
		t.Fatal("two successive IDs must not collide")
	}
}

// TestAllKnownEconomistErrors_Count verifies the enum is complete (fourteen
// entries: three from Ch. 10 — Smith — five from Ch. 11 — Ricardo — and six
// added in Ch. 19 — Quesnay×2, Smith×3, Ricardo×1).
func TestAllKnownEconomistErrors_Count(t *testing.T) {
	t.Parallel()
	if len(circulation.AllKnownEconomistErrors) != 14 {
		t.Fatalf("expected 14 KnownEconomistError values, got %d", len(circulation.AllKnownEconomistErrors))
	}
}

// TestAllKnownEconomistErrors_NoDuplicates verifies the enum has no repeated values.
func TestAllKnownEconomistErrors_NoDuplicates(t *testing.T) {
	t.Parallel()
	seen := make(map[circulation.KnownEconomistError]bool)
	for _, e := range circulation.AllKnownEconomistErrors {
		if seen[e] {
			t.Fatalf("duplicate KnownEconomistError value: %q", e)
		}
		seen[e] = true
	}
}

// TestEconomistAttribution_Quesnay_AvancesPrimitives exercises the Quesnay avances
// primitives fixture — the Tableau Économique (1758) — which anticipates the fixed
// capital distinction without theorising it.
func TestEconomistAttribution_Quesnay_AvancesPrimitives(t *testing.T) {
	t.Parallel()
	attr := circulation.EconomistAttribution{
		ID:          circulation.EconomistAttributionID("5eed000000000000001001"),
		Concept:     "avances primitives",
		Theorist:    "Quesnay",
		EditionYear: 1758,
		Anticipates: "fixed_capital_item",
		Errors:      nil,
	}
	if attr.Theorist != "Quesnay" {
		t.Errorf("expected theorist Quesnay, got %q", attr.Theorist)
	}
	if attr.EditionYear != 1758 {
		t.Errorf("expected edition year 1758, got %d", attr.EditionYear)
	}
	if attr.Anticipates != "fixed_capital_item" {
		t.Errorf("expected anticipates fixed_capital_item, got %q", attr.Anticipates)
	}
	if len(attr.Errors) != 0 {
		t.Errorf("Quesnay attribution must carry no errors; got %v", attr.Errors)
	}
}

// TestEconomistAttribution_Quesnay_AvancesAnnuelles verifies the second Quesnay
// concept (avances annuelles) anticipates circulating capital.
func TestEconomistAttribution_Quesnay_AvancesAnnuelles(t *testing.T) {
	t.Parallel()
	attr := circulation.EconomistAttribution{
		ID:          circulation.EconomistAttributionID("5eed000000000000001002"),
		Concept:     "avances annuelles",
		Theorist:    "Quesnay",
		EditionYear: 1758,
		Anticipates: "circulating",
		Errors:      nil,
	}
	if attr.Anticipates != "circulating" {
		t.Errorf("expected anticipates circulating, got %q", attr.Anticipates)
	}
}

// TestEconomistAttribution_Smith verifies that Adam Smith's Wealth of Nations
// (1776) carries exactly the three Ch. 10 KnownEconomistError entries.
func TestEconomistAttribution_Smith(t *testing.T) {
	t.Parallel()
	smithErrors := []circulation.KnownEconomistError{
		circulation.ErrorSmithConflation,
		circulation.ErrorSmithCirculationCapitalConflation,
		circulation.ErrorSmithRevenueInCapital,
	}
	attr := circulation.EconomistAttribution{
		ID:          circulation.EconomistAttributionID("5eed000000000000001003"),
		Concept:     "fixed and circulating stock",
		Theorist:    "Smith",
		EditionYear: 1776,
		Anticipates: "fixed_capital_item",
		Errors:      smithErrors,
	}
	if len(attr.Errors) != 3 {
		t.Fatalf("Smith attribution must carry 3 errors, got %d", len(attr.Errors))
	}
	present := make(map[circulation.KnownEconomistError]bool)
	for _, e := range attr.Errors {
		present[e] = true
	}
	for _, e := range smithErrors {
		if !present[e] {
			t.Errorf("Smith attribution missing error %q", e)
		}
	}
}

// TestEconomistAttribution_Ricardo verifies that Ricardo's Principles of
// Political Economy and Taxation (1817) carries all five Ch. 11
// KnownEconomistError entries.
func TestEconomistAttribution_Ricardo(t *testing.T) {
	t.Parallel()
	ricardoErrors := []circulation.KnownEconomistError{
		circulation.ErrorRicardoDurabilityCollapse,
		circulation.ErrorRicardoConflation,
		circulation.ErrorRicardoFixedCapitalPriceExplanation,
		circulation.ErrorRicardoNoAggregateTurnover,
		circulation.ErrorRicardoNoValueRevolution,
	}
	attr := circulation.EconomistAttribution{
		ID:          circulation.EconomistAttributionID("5eed000000000000001101"),
		Concept:     "fixed and circulating capital",
		Theorist:    "Ricardo",
		EditionYear: 1817,
		Anticipates: "prices_of_production",
		Errors:      ricardoErrors,
	}
	if attr.Theorist != "Ricardo" {
		t.Errorf("expected theorist Ricardo, got %q", attr.Theorist)
	}
	if attr.EditionYear != 1817 {
		t.Errorf("expected edition year 1817, got %d", attr.EditionYear)
	}
	if len(attr.Errors) != 5 {
		t.Fatalf("Ricardo attribution must carry 5 errors, got %d", len(attr.Errors))
	}
	present := make(map[circulation.KnownEconomistError]bool)
	for _, e := range attr.Errors {
		present[e] = true
	}
	for _, e := range ricardoErrors {
		if !present[e] {
			t.Errorf("Ricardo attribution missing error %q", e)
		}
	}
}

// TestEconomistAttribution_Quesnay_TableauReproduction verifies the Ch. 19
// Quesnay reproduction-context attribution: the Tableau Économique (1758)
// anticipates the social character of the reproduction circuit.
func TestEconomistAttribution_Quesnay_TableauReproduction(t *testing.T) {
	t.Parallel()
	quesnayErrors := []circulation.KnownEconomistError{
		circulation.ErrorQuesnayProductiveSterileDivision,
		circulation.ErrorQuesnayMissingValueTheory,
	}
	attr := circulation.EconomistAttribution{
		ID:          circulation.EconomistAttributionID("5eed000000000000001901"),
		Concept:     "Tableau Économique",
		Theorist:    "Quesnay",
		EditionYear: 1758,
		Anticipates: "reproduction_as_social",
		Errors:      quesnayErrors,
	}
	if attr.Theorist != "Quesnay" {
		t.Errorf("expected theorist Quesnay, got %q", attr.Theorist)
	}
	if attr.EditionYear != 1758 {
		t.Errorf("expected edition year 1758, got %d", attr.EditionYear)
	}
	if attr.Anticipates != "reproduction_as_social" {
		t.Errorf("expected anticipates reproduction_as_social, got %q", attr.Anticipates)
	}
	if len(attr.Errors) != 2 {
		t.Fatalf("Quesnay reproduction attribution must carry 2 errors, got %d", len(attr.Errors))
	}
	present := make(map[circulation.KnownEconomistError]bool)
	for _, e := range attr.Errors {
		present[e] = true
	}
	for _, e := range quesnayErrors {
		if !present[e] {
			t.Errorf("Quesnay reproduction attribution missing error %q", e)
		}
	}
}

// TestEconomistAttribution_Ricardo_Reproduction verifies the Ch. 19 Ricardo
// reproduction-context attribution: Principles (1817) reproduction entry.
func TestEconomistAttribution_Ricardo_Reproduction(t *testing.T) {
	t.Parallel()
	attr := circulation.EconomistAttribution{
		ID:          circulation.EconomistAttributionID("5eed000000000000001902"),
		Concept:     "reproduction and net revenue",
		Theorist:    "Ricardo",
		EditionYear: 1817,
		Anticipates: "reproduction_incomplete",
		Errors:      []circulation.KnownEconomistError{circulation.ErrorRicardoIncompleteReproduction},
	}
	if attr.Theorist != "Ricardo" {
		t.Errorf("expected theorist Ricardo, got %q", attr.Theorist)
	}
	if attr.EditionYear != 1817 {
		t.Errorf("expected edition year 1817, got %d", attr.EditionYear)
	}
	if len(attr.Errors) != 1 {
		t.Fatalf("Ricardo reproduction attribution must carry 1 error, got %d", len(attr.Errors))
	}
	if attr.Errors[0] != circulation.ErrorRicardoIncompleteReproduction {
		t.Errorf("expected ErrorRicardoIncompleteReproduction, got %q", attr.Errors[0])
	}
}

// TestEconomistAttribution_Smith_Ch19Errors verifies that the three Ch. 19
// additions to Smith's row — revenue dogma and its two consequences — are
// distinct from the three Ch. 10 errors already attributed to Smith.
func TestEconomistAttribution_Smith_Ch19Errors(t *testing.T) {
	t.Parallel()
	ch19SmithErrors := []circulation.KnownEconomistError{
		circulation.ErrorSmithRevenueDogma,
		circulation.ErrorSmithMissingConstantReplacement,
		circulation.ErrorSmithLabourRegress,
	}
	// Verify each Ch. 19 error is distinct from the Ch. 10 errors.
	ch10SmithErrors := map[circulation.KnownEconomistError]bool{
		circulation.ErrorSmithConflation:                   true,
		circulation.ErrorSmithCirculationCapitalConflation: true,
		circulation.ErrorSmithRevenueInCapital:             true,
	}
	for _, e := range ch19SmithErrors {
		if ch10SmithErrors[e] {
			t.Errorf("Ch. 19 Smith error %q duplicates a Ch. 10 error", e)
		}
	}
	// Verify they are all in AllKnownEconomistErrors.
	all := make(map[circulation.KnownEconomistError]bool)
	for _, e := range circulation.AllKnownEconomistErrors {
		all[e] = true
	}
	for _, e := range ch19SmithErrors {
		if !all[e] {
			t.Errorf("Ch. 19 Smith error %q not present in AllKnownEconomistErrors", e)
		}
	}
}

// TestEconomistAttribution_EveryErrorHasAttribution verifies the build-time
// invariant: every KnownEconomistError value appears in at least one seed
// attribution. This is the constructive proof that the enum is grounded in the
// chapter text.
func TestEconomistAttribution_EveryErrorHasAttribution(t *testing.T) {
	t.Parallel()

	seedAttributions := []circulation.EconomistAttribution{
		{
			ID:          "5eed000000000000001001",
			Concept:     "avances primitives",
			Theorist:    "Quesnay",
			EditionYear: 1758,
			Anticipates: "fixed_capital_item",
			Errors:      nil,
		},
		{
			ID:          "5eed000000000000001002",
			Concept:     "avances annuelles",
			Theorist:    "Quesnay",
			EditionYear: 1758,
			Anticipates: "circulating",
			Errors:      nil,
		},
		{
			ID:          "5eed000000000000001003",
			Concept:     "fixed and circulating stock",
			Theorist:    "Smith",
			EditionYear: 1776,
			Anticipates: "fixed_capital_item",
			Errors: []circulation.KnownEconomistError{
				circulation.ErrorSmithConflation,
				circulation.ErrorSmithCirculationCapitalConflation,
				circulation.ErrorSmithRevenueInCapital,
				// Ch. 19 extensions to the Smith row
				circulation.ErrorSmithRevenueDogma,
				circulation.ErrorSmithMissingConstantReplacement,
				circulation.ErrorSmithLabourRegress,
			},
		},
		{
			ID:          "5eed000000000000001101",
			Concept:     "fixed and circulating capital",
			Theorist:    "Ricardo",
			EditionYear: 1817,
			Anticipates: "prices_of_production",
			Errors: []circulation.KnownEconomistError{
				circulation.ErrorRicardoDurabilityCollapse,
				circulation.ErrorRicardoConflation,
				circulation.ErrorRicardoFixedCapitalPriceExplanation,
				circulation.ErrorRicardoNoAggregateTurnover,
				circulation.ErrorRicardoNoValueRevolution,
			},
		},
		// Ch. 19 new rows
		{
			ID:          "5eed000000000000001901",
			Concept:     "Tableau Économique",
			Theorist:    "Quesnay",
			EditionYear: 1758,
			Anticipates: "reproduction_as_social",
			Errors: []circulation.KnownEconomistError{
				circulation.ErrorQuesnayProductiveSterileDivision,
				circulation.ErrorQuesnayMissingValueTheory,
			},
		},
		{
			ID:          "5eed000000000000001902",
			Concept:     "reproduction and net revenue",
			Theorist:    "Ricardo",
			EditionYear: 1817,
			Anticipates: "reproduction_incomplete",
			Errors: []circulation.KnownEconomistError{
				circulation.ErrorRicardoIncompleteReproduction,
			},
		},
	}

	cited := make(map[circulation.KnownEconomistError]bool)
	for _, a := range seedAttributions {
		for _, e := range a.Errors {
			cited[e] = true
		}
	}
	for _, e := range circulation.AllKnownEconomistErrors {
		if !cited[e] {
			t.Errorf("KnownEconomistError %q has no seed attribution row", e)
		}
	}
}
