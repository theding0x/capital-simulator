package simulation

import (
	"testing"
)

// Ch. 33 — The Modern Theory of Colonisation.
// Fixtures drawn from Marx's actual cases: Mr. Peel at Swan River
// (300 working-class persons, £50,000 capital), Wakefield's "sufficient
// price" of land, and the systematic colonisation scheme.

func TestColonialLabourMarket_Validate(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name    string
		market  ColonialLabourMarket
		wantErr error
	}{
		{
			name: "Peel at Swan River is valid",
			market: ColonialLabourMarket{
				Colony:          "Swan River",
				FreeLabourers:   300,
				AnnualWagePence: 4800,
				LandAvailable:   true,
			},
			wantErr: nil,
		},
		{
			name: "colony required",
			market: ColonialLabourMarket{
				FreeLabourers:   300,
				AnnualWagePence: 4800,
			},
			wantErr: ErrColonyRequired,
		},
		{
			name: "free labourers cannot be negative",
			market: ColonialLabourMarket{
				Colony:          "Swan River",
				FreeLabourers:   -1,
				AnnualWagePence: 4800,
			},
			wantErr: ErrFreeLabourersNonNegative,
		},
		{
			name: "annual wage must be positive",
			market: ColonialLabourMarket{
				Colony:          "Swan River",
				FreeLabourers:   300,
				AnnualWagePence: 0,
			},
			wantErr: ErrAnnualWagePositive,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.market.Validate(); got != tc.wantErr {
				t.Fatalf("Validate() = %v, want %v", got, tc.wantErr)
			}
		})
	}
}

func TestSufficientPrice_Ransom(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name  string
		price SufficientPrice
		want  Pence
	}{
		{
			name:  "60 acres at 240d = 14400d (Wakefield's £60 ransom)",
			price: SufficientPrice{PricePerAcrePence: 240, DesiredAcres: 60},
			want:  14400,
		},
		{
			name:  "zero acres yields zero",
			price: SufficientPrice{PricePerAcrePence: 240, DesiredAcres: 0},
			want:  0,
		},
		{
			name:  "zero price yields zero",
			price: SufficientPrice{PricePerAcrePence: 0, DesiredAcres: 60},
			want:  0,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.price.Ransom(); got != tc.want {
				t.Fatalf("Ransom() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestColonialLabourRegulation_PeelAtSwanRiverWithoutScheme(t *testing.T) {
	t.Parallel()
	// A scheme with zero sufficient price is no scheme. The market is
	// returned unchanged: WakefieldSchemeApplied stays false, the wage
	// relation cannot extract surplus, labour walks off into
	// independent production. This is the Peel-fiasco case.
	market := ColonialLabourMarket{
		Colony:          "Swan River",
		FreeLabourers:   300,
		AnnualWagePence: 4800,
		LandAvailable:   true,
	}
	scheme := SystematicColonisation{Colony: "Swan River"} // empty sufficient price

	out := ColonialLabourRegulation(market, scheme)

	if out.WakefieldSchemeApplied {
		t.Errorf("WakefieldSchemeApplied = true, want false (no real scheme without a sufficient price)")
	}
	if out.IndependenceYears != 0 {
		t.Errorf("IndependenceYears = %d, want 0 (no ransom)", out.IndependenceYears)
	}
	if out.SurplusLabourExtractable {
		t.Errorf("SurplusLabourExtractable = true, want false (Peel's fiasco)")
	}
	if out.FreeLabourers != market.FreeLabourers {
		t.Errorf("FreeLabourers changed from %d to %d (invariant)", market.FreeLabourers, out.FreeLabourers)
	}
}

func TestColonialLabourRegulation_SufficientPriceExtendsIndependence(t *testing.T) {
	t.Parallel()
	// South Australia 1836: Wakefield's testbed. £20/year wage, 50%
	// savings = £10/year (2400 pence). Sufficient price set so the
	// ransom is 5 years' savings — labourer remains a wage-worker for
	// 5 years.
	market := ColonialLabourMarket{
		Colony:          "South Australia",
		FreeLabourers:   5000,
		AnnualWagePence: 4800, // £20/year
		LandAvailable:   true,
	}
	// Annual savings = 2400; ransom = 12000 -> 5 years.
	scheme := SystematicColonisation{
		Colony: "South Australia",
		SufficientPrice: SufficientPrice{
			PricePerAcrePence: 240,
			DesiredAcres:      50,
		},
	}

	out := ColonialLabourRegulation(market, scheme)

	if out.IndependenceYears != 5 {
		t.Errorf("IndependenceYears = %d, want 5 (ransom 12000 / savings 2400)", out.IndependenceYears)
	}
	if !out.SurplusLabourExtractable {
		t.Errorf("SurplusLabourExtractable = false, want true (scheme creates wage dependence)")
	}
	if !out.WakefieldSchemeApplied {
		t.Errorf("WakefieldSchemeApplied = false, want true")
	}
	if out.FreeLabourers != market.FreeLabourers {
		t.Errorf("FreeLabourers changed from %d to %d (invariant)", market.FreeLabourers, out.FreeLabourers)
	}
}

func TestColonialLabourRegulation_TenfoldPriceTenfoldDependence(t *testing.T) {
	t.Parallel()
	// At normal wages a worker buys land in 5 years; raising the
	// sufficient price tenfold should extend dependence to 50 years —
	// the comparison Marx draws in the chapter between Wakefield's
	// scheme and unregulated colonisation.
	market := ColonialLabourMarket{
		Colony:          "Test",
		FreeLabourers:   100,
		AnnualWagePence: 4800,
		LandAvailable:   true,
	}
	baseline := ColonialLabourRegulation(market, SystematicColonisation{
		SufficientPrice: SufficientPrice{PricePerAcrePence: 240, DesiredAcres: 50},
	})
	tenfold := ColonialLabourRegulation(market, SystematicColonisation{
		SufficientPrice: SufficientPrice{PricePerAcrePence: 2400, DesiredAcres: 50},
	})

	if baseline.IndependenceYears != 5 {
		t.Fatalf("baseline IndependenceYears = %d, want 5", baseline.IndependenceYears)
	}
	if tenfold.IndependenceYears != 50 {
		t.Fatalf("tenfold IndependenceYears = %d, want 50", tenfold.IndependenceYears)
	}
}

func TestComputeWageWorkerIndependence_ReachesIndependence(t *testing.T) {
	t.Parallel()
	market := ColonialLabourMarket{
		Colony:          "South Australia",
		FreeLabourers:   5000,
		AnnualWagePence: 4800,
	}
	price := SufficientPrice{PricePerAcrePence: 240, DesiredAcres: 50}

	indep := ComputeWageWorkerIndependence("settler-001", market, price, 0)

	if indep.WorkerID != "settler-001" {
		t.Errorf("WorkerID = %q", indep.WorkerID)
	}
	if indep.AnnualSavingsPence != 2400 {
		t.Errorf("AnnualSavingsPence = %d, want 2400", indep.AnnualSavingsPence)
	}
	if indep.TargetRansomPence != 12000 {
		t.Errorf("TargetRansomPence = %d, want 12000", indep.TargetRansomPence)
	}
	if indep.YearsWorked != 5 {
		t.Errorf("YearsWorked = %d, want 5", indep.YearsWorked)
	}
	if indep.SavingsPence != 12000 {
		t.Errorf("SavingsPence = %d, want 12000", indep.SavingsPence)
	}
	if !indep.BecameLandowner {
		t.Errorf("BecameLandowner = false, want true (savings >= ransom)")
	}
}

func TestComputeWageWorkerIndependence_YearsLimitPreventsExit(t *testing.T) {
	t.Parallel()
	// The labourer can only work 3 years before the projection ends —
	// not long enough to save the ransom, so BecameLandowner is false.
	market := ColonialLabourMarket{
		Colony:          "South Australia",
		FreeLabourers:   5000,
		AnnualWagePence: 4800,
	}
	price := SufficientPrice{PricePerAcrePence: 240, DesiredAcres: 50}

	indep := ComputeWageWorkerIndependence("indentured-002", market, price, 3)

	if indep.YearsWorked != 3 {
		t.Errorf("YearsWorked = %d, want 3 (limit)", indep.YearsWorked)
	}
	if indep.SavingsPence != 7200 {
		t.Errorf("SavingsPence = %d, want 7200 (3 * 2400)", indep.SavingsPence)
	}
	if indep.BecameLandowner {
		t.Errorf("BecameLandowner = true, want false (savings < ransom)")
	}
}

func TestRecordLandSale_GrowsLandFundAndImports(t *testing.T) {
	t.Parallel()
	// Wakefield's self-financing invariant: every land sale deposits
	// the ransom into the fund and imports one fresh have-nothing.
	price := SufficientPrice{PricePerAcrePence: 240, DesiredAcres: 50}
	scheme := SystematicColonisation{
		Colony:          "South Australia",
		SufficientPrice: price,
	}

	after1 := RecordLandSale(scheme)
	if after1.LandSalesCompleted != 1 {
		t.Errorf("LandSalesCompleted = %d, want 1", after1.LandSalesCompleted)
	}
	if after1.LandFundPence != 12000 {
		t.Errorf("LandFundPence = %d, want 12000 (one ransom)", after1.LandFundPence)
	}
	if after1.ImportedLabourers != 1 {
		t.Errorf("ImportedLabourers = %d, want 1", after1.ImportedLabourers)
	}

	after3 := RecordLandSale(RecordLandSale(after1))
	if after3.LandSalesCompleted != 3 {
		t.Errorf("LandSalesCompleted = %d, want 3", after3.LandSalesCompleted)
	}
	if after3.LandFundPence != 36000 {
		t.Errorf("LandFundPence = %d, want 36000 (3 * ransom)", after3.LandFundPence)
	}
	if after3.ImportedLabourers != 3 {
		t.Errorf("ImportedLabourers = %d, want 3", after3.ImportedLabourers)
	}
}

func TestNewColonialLabourMarketID_Unique(t *testing.T) {
	t.Parallel()
	seen := make(map[ColonialLabourMarketID]struct{}, 256)
	for i := 0; i < 256; i++ {
		id := NewColonialLabourMarketID()
		if id.IsZero() {
			t.Fatal("NewColonialLabourMarketID returned empty id")
		}
		if len(id) != 24 {
			t.Fatalf("id length = %d, want 24 hex chars", len(id))
		}
		if _, dup := seen[id]; dup {
			t.Fatalf("duplicate id generated: %s", id)
		}
		seen[id] = struct{}{}
	}
}
