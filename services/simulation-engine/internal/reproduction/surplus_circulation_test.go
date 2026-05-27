package reproduction_test

import (
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// Canonical fixture from the issue: 100-spinning-mill aggregate.
// Total surplus: 50_000_000 pence (≈£208,333).
// Realisation mix: 80% consumption, 19% capitalised, 1% gold (sums to 100%).
const (
	fixtureTotalPence    = reproduction.Pence(50_000_000)
	fixtureConsumptionBP = int64(8000) // 80%
	fixtureCapitalisedBP = int64(1900) // 19%
	fixtureGoldBP        = int64(100)  // 1%
)

func fixtureConsumptionPence() reproduction.Pence {
	return fixtureTotalPence * reproduction.Pence(fixtureConsumptionBP) / 10000
}
func fixtureCapitalisedPence() reproduction.Pence {
	return fixtureTotalPence * reproduction.Pence(fixtureCapitalisedBP) / 10000
}
func fixtureGoldPence() reproduction.Pence {
	return fixtureTotalPence * reproduction.Pence(fixtureGoldBP) / 10000
}

// newFixtureSC builds a SurplusCirculation matching the seed fixture.
func newFixtureSC() reproduction.SurplusCirculation {
	realised := fixtureConsumptionPence() + fixtureCapitalisedPence() + fixtureGoldPence()
	return reproduction.SurplusCirculation{
		ID:                     reproduction.NewSurplusCirculationID(),
		Period:                 "1871",
		TotalSurplusPence:      fixtureTotalPence,
		RealisedSurplusPence:   realised,
		UnrealisedSurplusPence: fixtureTotalPence - realised,
		RealisationSources: []reproduction.RealisationSourceEntry{
			{Source: reproduction.SourceCapitalistConsumption, Pence: fixtureConsumptionPence()},
			{Source: reproduction.SourceCapitalisedSurplus, Pence: fixtureCapitalisedPence()},
			{Source: reproduction.SourceGoldProduction, Pence: fixtureGoldPence()},
		},
	}
}

// TestSurplusCirculation_Validate_HappyPath verifies that a well-formed
// SurplusCirculation passes validation.
func TestSurplusCirculation_Validate_HappyPath(t *testing.T) {
	t.Parallel()
	sc := newFixtureSC()
	if err := sc.Validate(); err != nil {
		t.Fatalf("expected valid: %v", err)
	}
}

// TestSurplusCirculation_Validate_NegativeTotal verifies ErrNegativeSurplus.
func TestSurplusCirculation_Validate_NegativeTotal(t *testing.T) {
	t.Parallel()
	sc := reproduction.SurplusCirculation{
		TotalSurplusPence:      -1,
		RealisedSurplusPence:   0,
		UnrealisedSurplusPence: -1,
	}
	if err := sc.Validate(); err != reproduction.ErrNegativeSurplus {
		t.Fatalf("expected ErrNegativeSurplus, got %v", err)
	}
}

// TestSurplusCirculation_Validate_MixMismatch verifies ErrSourceMixMismatch
// when the source entries do not sum to the realised amount.
func TestSurplusCirculation_Validate_MixMismatch(t *testing.T) {
	t.Parallel()
	sc := reproduction.SurplusCirculation{
		TotalSurplusPence:      1000,
		RealisedSurplusPence:   1000,
		UnrealisedSurplusPence: 0,
		RealisationSources: []reproduction.RealisationSourceEntry{
			{Source: reproduction.SourceCapitalistConsumption, Pence: 500}, // only 500, not 1000
		},
	}
	if err := sc.Validate(); err != reproduction.ErrSourceMixMismatch {
		t.Fatalf("expected ErrSourceMixMismatch, got %v", err)
	}
}

// TestSurplusCirculation_Validate_ZeroSurplus verifies that a zero surplus
// with zero realisation and no sources is valid.
func TestSurplusCirculation_Validate_ZeroSurplus(t *testing.T) {
	t.Parallel()
	sc := reproduction.SurplusCirculation{
		ID:     reproduction.NewSurplusCirculationID(),
		Period: "test",
	}
	if err := sc.Validate(); err != nil {
		t.Fatalf("expected valid zero-surplus: %v", err)
	}
}

// TestAddRealisationSource_ConsumptionAndCapitalised verifies the happy path
// for adding consumption and capitalised surplus sources.
func TestAddRealisationSource_ConsumptionAndCapitalised(t *testing.T) {
	t.Parallel()
	sc := reproduction.SurplusCirculation{
		ID:                     reproduction.NewSurplusCirculationID(),
		Period:                 "1871",
		TotalSurplusPence:      fixtureTotalPence,
		UnrealisedSurplusPence: fixtureTotalPence,
		RealisationSources:     []reproduction.RealisationSourceEntry{},
	}

	var err error
	sc, err = reproduction.AddRealisationSource(sc, reproduction.RealisationSourceEntry{
		Source: reproduction.SourceCapitalistConsumption,
		Pence:  fixtureConsumptionPence(),
	})
	if err != nil {
		t.Fatalf("adding consumption: %v", err)
	}
	sc, err = reproduction.AddRealisationSource(sc, reproduction.RealisationSourceEntry{
		Source: reproduction.SourceCapitalisedSurplus,
		Pence:  fixtureCapitalisedPence(),
	})
	if err != nil {
		t.Fatalf("adding capitalised: %v", err)
	}

	wantRealised := fixtureConsumptionPence() + fixtureCapitalisedPence()
	if sc.RealisedSurplusPence != wantRealised {
		t.Fatalf("realised: want %d, got %d", wantRealised, sc.RealisedSurplusPence)
	}
	if sc.UnrealisedSurplusPence != fixtureTotalPence-wantRealised {
		t.Fatalf("unrealised: want %d, got %d", fixtureTotalPence-wantRealised, sc.UnrealisedSurplusPence)
	}
}

// TestAddRealisationSource_GoldExpandsMoneyStock verifies that after adding
// gold-production pence, UnrealisedSurplusPence decreases by that amount.
func TestAddRealisationSource_GoldExpandsMoneyStock(t *testing.T) {
	t.Parallel()
	sc := reproduction.SurplusCirculation{
		ID:                     reproduction.NewSurplusCirculationID(),
		Period:                 "1871",
		TotalSurplusPence:      1200,
		UnrealisedSurplusPence: 1200,
		RealisationSources:     []reproduction.RealisationSourceEntry{},
	}
	goldPence := reproduction.Pence(12) // 1% of 1200
	got, err := reproduction.AddRealisationSource(sc, reproduction.RealisationSourceEntry{
		Source: reproduction.SourceGoldProduction,
		Pence:  goldPence,
	})
	if err != nil {
		t.Fatalf("adding gold: %v", err)
	}
	if got.RealisedSurplusPence != goldPence {
		t.Fatalf("realised: want %d, got %d", goldPence, got.RealisedSurplusPence)
	}
	if got.UnrealisedSurplusPence != sc.TotalSurplusPence-goldPence {
		t.Fatalf("unrealised: want %d, got %d", sc.TotalSurplusPence-goldPence, got.UnrealisedSurplusPence)
	}
}

// TestAddRealisationSource_CreditCreatesDeferred verifies that a credit
// source sets DeferredRealisationPence equal to Pence.
func TestAddRealisationSource_CreditCreatesDeferred(t *testing.T) {
	t.Parallel()
	sc := reproduction.SurplusCirculation{
		ID:                     reproduction.NewSurplusCirculationID(),
		TotalSurplusPence:      1000,
		UnrealisedSurplusPence: 1000,
		RealisationSources:     []reproduction.RealisationSourceEntry{},
	}
	got, err := reproduction.AddRealisationSource(sc, reproduction.RealisationSourceEntry{
		Source: reproduction.SourceCredit,
		Pence:  200,
	})
	if err != nil {
		t.Fatalf("adding credit: %v", err)
	}
	if len(got.RealisationSources) != 1 {
		t.Fatalf("expected 1 source entry, got %d", len(got.RealisationSources))
	}
	entry := got.RealisationSources[0]
	if entry.DeferredRealisationPence != 200 {
		t.Fatalf("deferred: want 200, got %d", entry.DeferredRealisationPence)
	}
}

// TestAddRealisationSource_ForeignTradeZero verifies that foreign trade can
// be recorded with zero pence (Vol. II abstraction).
func TestAddRealisationSource_ForeignTradeZero(t *testing.T) {
	t.Parallel()
	sc := reproduction.SurplusCirculation{
		ID:                     reproduction.NewSurplusCirculationID(),
		TotalSurplusPence:      1000,
		UnrealisedSurplusPence: 1000,
		RealisationSources:     []reproduction.RealisationSourceEntry{},
	}
	got, err := reproduction.AddRealisationSource(sc, reproduction.RealisationSourceEntry{
		Source: reproduction.SourceForeignTrade,
		Pence:  0,
	})
	if err != nil {
		t.Fatalf("foreign trade zero: %v", err)
	}
	if got.RealisedSurplusPence != 0 {
		t.Fatal("realised should remain 0 with zero foreign trade")
	}
}

// TestAddRealisationSource_ExceedsTotal verifies ErrRealisationExceedsTotal.
func TestAddRealisationSource_ExceedsTotal(t *testing.T) {
	t.Parallel()
	sc := reproduction.SurplusCirculation{
		ID:                     reproduction.NewSurplusCirculationID(),
		TotalSurplusPence:      500,
		UnrealisedSurplusPence: 500,
		RealisationSources:     []reproduction.RealisationSourceEntry{},
	}
	_, err := reproduction.AddRealisationSource(sc, reproduction.RealisationSourceEntry{
		Source: reproduction.SourceCapitalistConsumption,
		Pence:  600, // > 500
	})
	if err != reproduction.ErrRealisationExceedsTotal {
		t.Fatalf("expected ErrRealisationExceedsTotal, got %v", err)
	}
}

// TestAddRealisationSource_NegativePence verifies ErrNegativePence.
func TestAddRealisationSource_NegativePence(t *testing.T) {
	t.Parallel()
	sc := reproduction.SurplusCirculation{
		ID:                     reproduction.NewSurplusCirculationID(),
		TotalSurplusPence:      1000,
		UnrealisedSurplusPence: 1000,
		RealisationSources:     []reproduction.RealisationSourceEntry{},
	}
	_, err := reproduction.AddRealisationSource(sc, reproduction.RealisationSourceEntry{
		Source: reproduction.SourceCapitalistConsumption,
		Pence:  -1,
	})
	if err != reproduction.ErrNegativePence {
		t.Fatalf("expected ErrNegativePence, got %v", err)
	}
}

// TestAddRealisationSource_UnknownSource verifies ErrUnknownSource.
func TestAddRealisationSource_UnknownSource(t *testing.T) {
	t.Parallel()
	sc := reproduction.SurplusCirculation{
		ID:                     reproduction.NewSurplusCirculationID(),
		TotalSurplusPence:      1000,
		UnrealisedSurplusPence: 1000,
		RealisationSources:     []reproduction.RealisationSourceEntry{},
	}
	_, err := reproduction.AddRealisationSource(sc, reproduction.RealisationSourceEntry{
		Source: "unknown_source",
		Pence:  100,
	})
	if err != reproduction.ErrUnknownSource {
		t.Fatalf("expected ErrUnknownSource, got %v", err)
	}
}

// TestRealisationSource_Validate verifies all five source constants are valid.
func TestRealisationSource_Validate(t *testing.T) {
	t.Parallel()
	valid := []reproduction.RealisationSource{
		reproduction.SourceCapitalistConsumption,
		reproduction.SourceCapitalisedSurplus,
		reproduction.SourceForeignTrade,
		reproduction.SourceGoldProduction,
		reproduction.SourceCredit,
	}
	for _, s := range valid {
		s := s
		t.Run(string(s), func(t *testing.T) {
			t.Parallel()
			if err := s.Validate(); err != nil {
				t.Fatalf("expected valid source %q: %v", s, err)
			}
		})
	}
	if err := reproduction.RealisationSource("bad").Validate(); err != reproduction.ErrUnknownSource {
		t.Fatalf("expected ErrUnknownSource for bad source, got %v", err)
	}
}

// TestNewSurplusCirculationID verifies that IDs are non-empty and distinct.
func TestNewSurplusCirculationID(t *testing.T) {
	t.Parallel()
	a := reproduction.NewSurplusCirculationID()
	b := reproduction.NewSurplusCirculationID()
	if a == "" {
		t.Fatal("expected non-empty ID")
	}
	if a == b {
		t.Fatal("expected distinct IDs")
	}
	if len(string(a)) != 24 {
		t.Fatalf("expected 24-char hex ID, got %d chars", len(string(a)))
	}
}

// TestSurplusCirculation_FullMix_SeedFixture verifies the full 100-spinning-
// mill seed fixture: 80% consumption + 20% capitalised + 1% gold = 101%
// would exceed total; actually the fixture uses absolute pence that sum to
// TotalSurplusPence leaving no unrealised. With the seed values:
//
//	consumption = 40_000_000
//	capitalised = 10_000_000
//	gold        =    500_000
//	Total       = 50_500_000  — but that exceeds 50_000_000,
//
// so the seed uses 80% + 20% + 0% = 100% with gold as additional over
// realised portion. This test uses the exact fixture percentages
// from the issue (80% + 20% + 1% = 101% of surplus) — per the issue
// they sum to more than total, so the test validates partial realisation
// with UnrealisedSurplusPence accounting for the gap.
//
// Re-reading the issue: "80% consumption, 20% capitalised, 1% gold" sums
// to 101%; with UnrealisedSurplusPence = 0, these cannot all be 101%.
// The correct interpretation: consumption=80%, capitalised=19%, gold=1%
// totalling 100%. We test with these corrected values.
func TestSurplusCirculation_FullMix_SeedFixture(t *testing.T) {
	t.Parallel()
	total := reproduction.Pence(50_000_000)
	cons := reproduction.Pence(40_000_000) // 80%
	cap2 := reproduction.Pence(9_500_000)  // 19%
	gold := reproduction.Pence(500_000)    // 1%

	sc := reproduction.SurplusCirculation{
		ID:                     reproduction.NewSurplusCirculationID(),
		Period:                 "1871",
		TotalSurplusPence:      total,
		RealisedSurplusPence:   cons + cap2 + gold,
		UnrealisedSurplusPence: total - cons - cap2 - gold,
		RealisationSources: []reproduction.RealisationSourceEntry{
			{Source: reproduction.SourceCapitalistConsumption, Pence: cons},
			{Source: reproduction.SourceCapitalisedSurplus, Pence: cap2},
			{Source: reproduction.SourceGoldProduction, Pence: gold},
		},
	}
	if err := sc.Validate(); err != nil {
		t.Fatalf("full-mix fixture: %v", err)
	}
	if sc.TotalSurplusPence != sc.RealisedSurplusPence+sc.UnrealisedSurplusPence {
		t.Fatal("invariant: total != realised + unrealised")
	}
}
