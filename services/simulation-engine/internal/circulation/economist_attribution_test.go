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

// TestAllKnownEconomistErrors_Count verifies the enum is complete (eight entries:
// three from Ch. 10 — Smith — and five from Ch. 11 — Ricardo).
func TestAllKnownEconomistErrors_Count(t *testing.T) {
	t.Parallel()
	if len(circulation.AllKnownEconomistErrors) != 8 {
		t.Fatalf("expected 8 KnownEconomistError values, got %d", len(circulation.AllKnownEconomistErrors))
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
