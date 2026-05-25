package costs_test

import (
	"testing"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/costs"
)

// --- NatureOf ---------------------------------------------------------------

func TestNatureOf_FauxFrais(t *testing.T) {
	t.Parallel()
	fauxFrais := []costs.CirculationCostKind{
		costs.CostPurchaseSaleTime,
		costs.CostBookKeeping,
		costs.CostMoneyAsFauxFraisKind,
		costs.CostSupplyFormation,
		costs.CostCommoditySupply,
	}
	for _, k := range fauxFrais {
		k := k
		t.Run(string(k), func(t *testing.T) {
			t.Parallel()
			if got := costs.NatureOf(k); got != costs.NatureValuePreserving {
				t.Errorf("NatureOf(%s) = %s, want value_preserving", k, got)
			}
		})
	}
}

func TestNatureOf_Transportation_ValueAdding(t *testing.T) {
	t.Parallel()
	if got := costs.NatureOf(costs.CostTransportation); got != costs.NatureValueAdding {
		t.Errorf("NatureOf(transportation) = %s, want value_adding", got)
	}
}

// --- NewCirculationCost -----------------------------------------------------

func TestNewCirculationCost_Happy(t *testing.T) {
	t.Parallel()
	c, err := costs.NewCirculationCost("ic-001", costs.CostBookKeeping, 1200, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.ID == "" {
		t.Error("ID should be set")
	}
	if c.Nature != costs.NatureValuePreserving {
		t.Errorf("Nature = %s, want value_preserving", c.Nature)
	}
}

func TestNewCirculationCost_Transportation_NatureValueAdding(t *testing.T) {
	t.Parallel()
	c, err := costs.NewCirculationCost("ic-002", costs.CostTransportation, 500, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Nature != costs.NatureValueAdding {
		t.Errorf("Nature = %s, want value_adding", c.Nature)
	}
}

func TestNewCirculationCost_EmptyCapitalID_Error(t *testing.T) {
	t.Parallel()
	_, err := costs.NewCirculationCost("", costs.CostBookKeeping, 100, time.Now())
	if err == nil {
		t.Error("expected error for empty IndustrialCapitalID")
	}
}

func TestNewCirculationCost_NegativePence_Error(t *testing.T) {
	t.Parallel()
	_, err := costs.NewCirculationCost("ic-003", costs.CostBookKeeping, -1, time.Now())
	if err == nil {
		t.Error("expected error for negative pence")
	}
}

// --- CirculationAgent -------------------------------------------------------

// §I(a): "He may receive daily the value of the product of eight working-hours,
// yet functions ten. But the two hours of surplus-labour he performs do not
// produce value anymore than his eight hours of necessary labour."
//
// Fixture: Spinning Mill 1871 seller — 480 necessary + 120 surplus minutes.
func TestCirculationAgent_ValueContribution_AlwaysZero(t *testing.T) {
	t.Parallel()
	agent := costs.CirculationAgent{
		ID:                     costs.NewCirculationAgentID(),
		IndustrialCapitalID:    "ic-spinner-1871",
		Role:                   costs.RoleSeller,
		WageRatePence:          480,
		LabourMinutesNecessary: 480,
		LabourMinutesSurplus:   120,
	}
	if got := agent.ValueContribution(); got != 0 {
		t.Errorf("ValueContribution() = %d, want 0", got)
	}
}

func TestCirculationAgent_WageEffectivePence_ReducedBySurplus(t *testing.T) {
	t.Parallel()
	// 480 necessary + 120 surplus = 600 total; effective wage = 480 × 480/600 = 384
	agent := costs.CirculationAgent{
		WageRatePence:          480,
		LabourMinutesNecessary: 480,
		LabourMinutesSurplus:   120,
	}
	got := agent.WageEffectivePence()
	want := costs.Pence(384)
	if got != want {
		t.Errorf("WageEffectivePence() = %d, want %d", got, want)
	}
}

func TestCirculationAgent_WageEffectivePence_ZeroTotal_NoDiv(t *testing.T) {
	t.Parallel()
	agent := costs.CirculationAgent{WageRatePence: 100}
	if got := agent.WageEffectivePence(); got != 0 {
		t.Errorf("WageEffectivePence() with zero labour = %d, want 0", got)
	}
}

// §I(b): Spinning Mill 1871 book-keeper at 1,000 pence/week labour + 200 stationery = 1,200.
func TestNewCirculationCost_BookKeeping_Nature(t *testing.T) {
	t.Parallel()
	c, err := costs.NewCirculationCost("ic-spinner-1871", costs.CostBookKeeping, 1200, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if c.Nature != costs.NatureValuePreserving {
		t.Error("book-keeping cost must be value_preserving")
	}
}

// --- StorageCost ------------------------------------------------------------

func TestStorageCost_Total(t *testing.T) {
	t.Parallel()
	// §II(b) voluntary supply: building 500 + labour 300 + preservation 200 = 1000.
	sc := costs.StorageCost{
		ID:                      costs.NewStorageCostID(),
		CommoditySupplyID:       "cs-001",
		BuildingPence:           500,
		LabourPence:             300,
		PreservationLabourPence: 200,
	}
	if got := sc.Total(); got != 1000 {
		t.Errorf("Total() = %d, want 1000", got)
	}
}

// --- TransportTariff --------------------------------------------------------

// §III Birmingham glass: historic 10s/ton over 50 miles ≈ 2.4d/ton-mile.
// Modern: 3× for breakage risk (BreakageRiskMultiplierBasisPoints = 30000).
func TestTransportTariff_EffectiveRate_BreakageRisk3x(t *testing.T) {
	t.Parallel()
	tariff := costs.TransportTariff{
		BasePencePerTonMile:               24,
		BreakageRiskMultiplierBasisPoints: 30000,
	}
	got := tariff.EffectiveRatePencePerTonMile()
	want := costs.Pence(72)
	if got != want {
		t.Errorf("EffectiveRatePencePerTonMile() = %d, want %d", got, want)
	}
}

func TestTransportTariff_EffectiveRate_NoMultiplier(t *testing.T) {
	t.Parallel()
	tariff := costs.TransportTariff{BasePencePerTonMile: 10}
	if got := tariff.EffectiveRatePencePerTonMile(); got != 10 {
		t.Errorf("EffectiveRatePencePerTonMile() with no multiplier = %d, want 10", got)
	}
}

func TestTransportTariff_EffectiveRate_FragilityBeatsBreakage(t *testing.T) {
	t.Parallel()
	tariff := costs.TransportTariff{
		BasePencePerTonMile:               10,
		FragilityMultiplierBasisPoints:    20000,
		BreakageRiskMultiplierBasisPoints: 15000,
	}
	if got := tariff.EffectiveRatePencePerTonMile(); got != 20 {
		t.Errorf("EffectiveRatePencePerTonMile() = %d, want 20", got)
	}
}

// --- TransportLeg -----------------------------------------------------------

func TestTransportLeg_AddedValue(t *testing.T) {
	t.Parallel()
	leg := costs.TransportLeg{
		LabourCostPence:       300,
		MeansOfTransportPence: 100,
		SurplusPence:          50,
	}
	if got := leg.AddedValue(); got != 450 {
		t.Errorf("AddedValue() = %d, want 450", got)
	}
}

func TestTransportLeg_AddedValue_ZeroForEmptyLeg(t *testing.T) {
	t.Parallel()
	var leg costs.TransportLeg
	if got := leg.AddedValue(); got != 0 {
		t.Errorf("empty TransportLeg.AddedValue() = %d, want 0", got)
	}
}

// §III distance invariant: doubling the cost components (proportional to distance)
// doubles AddedValue().
func TestTransportLeg_AddedValue_MonotonicInDistance(t *testing.T) {
	t.Parallel()
	short := costs.TransportLeg{LabourCostPence: 100, MeansOfTransportPence: 50, SurplusPence: 25}
	long := costs.TransportLeg{LabourCostPence: 200, MeansOfTransportPence: 100, SurplusPence: 50}
	if long.AddedValue() != 2*short.AddedValue() {
		t.Errorf("doubling distance components should double AddedValue: got %d, want %d",
			long.AddedValue(), 2*short.AddedValue())
	}
}

// §III Birmingham glass: relative transport share rises as commodity value falls.
// £11/crate → transport fraction; £2/crate → transport fraction ≈ 5.5× larger.
func TestTransportLeg_RelativeShareInverselyProportionalToValue(t *testing.T) {
	t.Parallel()
	addedValue := costs.Pence(120)
	historicCrateValue := costs.Pence(2640) // £11 × 240
	modernCrateValue := costs.Pence(480)    // £2  × 240

	historicShare := float64(addedValue) / float64(historicCrateValue)
	modernShare := float64(addedValue) / float64(modernCrateValue)

	if modernShare <= historicShare {
		t.Errorf("modern share (%f) should exceed historic share (%f)", modernShare, historicShare)
	}
	ratio := modernShare / historicShare
	if ratio < 5.0 || ratio > 6.0 {
		t.Errorf("ratio modernShare/historicShare = %f, want ≈5.5", ratio)
	}
}

// §I(c): National gold reserve 1,000,000 pence, 1% annual wear = 10,000 replacement.
func TestMoneyAsFauxFrais_AnnualReplacement(t *testing.T) {
	t.Parallel()
	mf := costs.MoneyAsFauxFrais{
		ID:                     costs.NewMoneyAsFauxFraisID(),
		Pence:                  1_000_000,
		AnnualReplacementPence: 10_000,
	}
	if mf.AnnualReplacementPence <= 0 {
		t.Error("AnnualReplacementPence must be positive")
	}
	if mf.AnnualReplacementPence >= mf.Pence {
		t.Error("AnnualReplacementPence must be less than total reserve")
	}
}
