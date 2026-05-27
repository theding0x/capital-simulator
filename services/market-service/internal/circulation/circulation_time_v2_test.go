package circulation

import (
	"errors"
	"testing"
	"time"
)

// England–India fixture dimensions (from the seed migration).
// Pre-telegraph (1850): 365-day selling phase; comm lag = 250 days.
// Post-telegraph (1870): 30-day selling phase; comm lag = 1 day.

func TestSellingPhaseDetailed_ValidateSubSpans(t *testing.T) {
	t.Parallel()

	t.Run("valid pre-telegraph England-India phase sums correctly", func(t *testing.T) {
		t.Parallel()
		d := SellingPhaseDetailed{
			BuyerSearchMinutes:      72000,  // 50 days
			BargainingMinutes:       72000,  // 50 days
			CommunicationLagMinutes: 360000, // 250 days (ship post)
			LegalSettlementMinutes:  21600,  // 15 days
			TotalMinutes:            525600, // 365 days
		}
		if err := d.ValidateSubSpans(); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("valid post-telegraph phase sums correctly", func(t *testing.T) {
		t.Parallel()
		d := SellingPhaseDetailed{
			BuyerSearchMinutes:      14400, // 10 days
			BargainingMinutes:       14400, // 10 days
			CommunicationLagMinutes: 1440,  // 1 day (telegraph)
			LegalSettlementMinutes:  12960, // 9 days
			TotalMinutes:            43200, // 30 days
		}
		if err := d.ValidateSubSpans(); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("sub-spans do not sum to total", func(t *testing.T) {
		t.Parallel()
		d := SellingPhaseDetailed{
			BuyerSearchMinutes:      100,
			BargainingMinutes:       100,
			CommunicationLagMinutes: 100,
			LegalSettlementMinutes:  100,
			TotalMinutes:            500, // 400 != 500
		}
		err := d.ValidateSubSpans()
		if err == nil {
			t.Fatal("expected ErrSubSpanMismatch, got nil")
		}
		if !errors.Is(err, ErrSubSpanMismatch) {
			t.Fatalf("expected ErrSubSpanMismatch, got: %v", err)
		}
	})

	t.Run("zero total with zero sub-spans passes", func(t *testing.T) {
		t.Parallel()
		d := SellingPhaseDetailed{}
		if err := d.ValidateSubSpans(); err != nil {
			t.Fatalf("expected no error for all-zero, got: %v", err)
		}
	})
}

func TestSellingPhaseDetailed_ValidateLagBound(t *testing.T) {
	t.Parallel()

	t.Run("pre-telegraph comm lag within 58 min/km * 9000 km bound", func(t *testing.T) {
		t.Parallel()
		d := SellingPhaseDetailed{
			MarketDistanceMeters:    9000000, // 9 000 km
			CommunicationLagMinutes: 360000,  // 250 days
		}
		// 58 min/km * 9000 km = 522 000 min; 360 000 < 522 000: valid
		if err := d.ValidateLagBound(58); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("post-telegraph comm lag within 1 min/km * 9000 km bound", func(t *testing.T) {
		t.Parallel()
		d := SellingPhaseDetailed{
			MarketDistanceMeters:    9000000, // 9 000 km
			CommunicationLagMinutes: 1440,    // 1 day (telegraph)
		}
		// 1 min/km * 9000 km = 9 000 min; 1 440 < 9 000: valid
		if err := d.ValidateLagBound(1); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("lag exceeds max returns ErrLagExceedsMaximum", func(t *testing.T) {
		t.Parallel()
		d := SellingPhaseDetailed{
			MarketDistanceMeters:    320000, // 320 km (London-Manchester)
			CommunicationLagMinutes: 1000,   // 1000 min
		}
		// 1 min/km * 320 km = 320 min; 1000 > 320: exceeds
		err := d.ValidateLagBound(1)
		if err == nil {
			t.Fatal("expected ErrLagExceedsMaximum, got nil")
		}
		if !errors.Is(err, ErrLagExceedsMaximum) {
			t.Fatalf("expected ErrLagExceedsMaximum, got: %v", err)
		}
	})

	t.Run("zero lag always valid", func(t *testing.T) {
		t.Parallel()
		d := SellingPhaseDetailed{
			MarketDistanceMeters:    9000000,
			CommunicationLagMinutes: 0,
		}
		if err := d.ValidateLagBound(58); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})
}

func TestBuyingPhaseDetailed_ValidateSubSpans(t *testing.T) {
	t.Parallel()

	t.Run("valid buying sub-spans sum correctly", func(t *testing.T) {
		t.Parallel()
		d := BuyingPhaseDetailed{
			SupplierSearchMinutes:   14400, // 10 days
			OrderProcessingMinutes:  14400, // 10 days
			CommunicationLagMinutes: 2880,  // 2 days
			LegalSettlementMinutes:  11520, // 8 days
			TotalMinutes:            43200, // 30 days
		}
		if err := d.ValidateSubSpans(); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("mismatched buying sub-spans return ErrSubSpanMismatch", func(t *testing.T) {
		t.Parallel()
		d := BuyingPhaseDetailed{
			SupplierSearchMinutes:   100,
			OrderProcessingMinutes:  100,
			CommunicationLagMinutes: 100,
			LegalSettlementMinutes:  100,
			TotalMinutes:            450, // 400 != 450
		}
		err := d.ValidateSubSpans()
		if !errors.Is(err, ErrSubSpanMismatch) {
			t.Fatalf("expected ErrSubSpanMismatch, got: %v", err)
		}
	})
}

func TestBuyingPhaseDetailed_ValidateLagBound(t *testing.T) {
	t.Parallel()

	t.Run("lag within bound passes", func(t *testing.T) {
		t.Parallel()
		d := BuyingPhaseDetailed{
			MarketDistanceMeters:    9000000,
			CommunicationLagMinutes: 1440,
		}
		if err := d.ValidateLagBound(1); err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}
	})

	t.Run("lag exceeds bound returns error", func(t *testing.T) {
		t.Parallel()
		d := BuyingPhaseDetailed{
			MarketDistanceMeters:    100000, // 100 km
			CommunicationLagMinutes: 500,
		}
		// 2 min/km * 100 km = 200; 500 > 200
		if !errors.Is(d.ValidateLagBound(2), ErrLagExceedsMaximum) {
			t.Fatal("expected ErrLagExceedsMaximum")
		}
	})
}

func TestDistanceLagRelation_MaxLagMinutes(t *testing.T) {
	t.Parallel()

	t.Run("London-Bombay pre-telegraph: 58 min/km * 9000 km = 522000 min", func(t *testing.T) {
		t.Parallel()
		r := DistanceLagRelation{
			DistanceMeters:  9000000,
			LagMinutesPerKm: 58,
		}
		got := r.MaxLagMinutes()
		want := int64(522000)
		if got != want {
			t.Fatalf("MaxLagMinutes = %d, want %d", got, want)
		}
	})

	t.Run("London-Manchester post-telephone: 0 min/km * 320 km = 0", func(t *testing.T) {
		t.Parallel()
		r := DistanceLagRelation{
			DistanceMeters:  320000,
			LagMinutesPerKm: 0,
		}
		if r.MaxLagMinutes() != 0 {
			t.Fatal("expected 0 for zero lag rate")
		}
	})
}

func TestCirculationSpeedImprovement_LagReduction(t *testing.T) {
	t.Parallel()

	t.Run("telegraph reduced London-Bombay from 58 to 1 min/km", func(t *testing.T) {
		t.Parallel()
		csi := CirculationSpeedImprovement{
			OldLagMinutesPerKm: 58,
			NewLagMinutesPerKm: 1,
		}
		got := csi.LagReduction()
		want := int64(57)
		if got != want {
			t.Fatalf("LagReduction = %d, want %d", got, want)
		}
	})

	t.Run("zero reduction when no change", func(t *testing.T) {
		t.Parallel()
		csi := CirculationSpeedImprovement{
			OldLagMinutesPerKm: 10,
			NewLagMinutesPerKm: 10,
		}
		if csi.LagReduction() != 0 {
			t.Fatal("expected 0 for no-change improvement")
		}
	})
}

func TestAnnualSurplusPenalty_ComputePenalty(t *testing.T) {
	t.Parallel()

	t.Run("England-India 1850: ideal 10000 bp, actual 5000 bp -> penalty 5000 bp", func(t *testing.T) {
		t.Parallel()
		asp := AnnualSurplusPenalty{
			IdealAnnualRateBasisPoints:  10000,
			ActualAnnualRateBasisPoints: 5000,
		}
		got := asp.ComputePenalty()
		want := int64(5000)
		if got != want {
			t.Fatalf("ComputePenalty = %d, want %d", got, want)
		}
	})

	t.Run("penalty is zero when actual >= ideal", func(t *testing.T) {
		t.Parallel()
		asp := AnnualSurplusPenalty{
			IdealAnnualRateBasisPoints:  5000,
			ActualAnnualRateBasisPoints: 7000,
		}
		if asp.ComputePenalty() != 0 {
			t.Fatal("expected 0 when actual >= ideal")
		}
	})

	t.Run("zero circulation time -> zero penalty", func(t *testing.T) {
		t.Parallel()
		asp := AnnualSurplusPenalty{
			IdealAnnualRateBasisPoints:  10000,
			ActualAnnualRateBasisPoints: 10000,
		}
		if asp.ComputePenalty() != 0 {
			t.Fatal("expected 0 for equal rates")
		}
	})
}

func TestNewIDs(t *testing.T) {
	t.Parallel()

	t.Run("NewSellingPhaseDetailedID produces 24-char unique IDs", func(t *testing.T) {
		t.Parallel()
		id := NewSellingPhaseDetailedID()
		if len(string(id)) != 24 {
			t.Fatalf("expected 24 hex chars, got %d: %s", len(string(id)), id)
		}
		if id == NewSellingPhaseDetailedID() {
			t.Fatal("two sequential IDs should not be equal")
		}
	})

	t.Run("NewBuyingPhaseDetailedID is unique", func(t *testing.T) {
		t.Parallel()
		if NewBuyingPhaseDetailedID() == NewBuyingPhaseDetailedID() {
			t.Fatal("expected unique IDs")
		}
	})

	t.Run("NewDistanceLagRelationID is unique", func(t *testing.T) {
		t.Parallel()
		if NewDistanceLagRelationID() == NewDistanceLagRelationID() {
			t.Fatal("expected unique IDs")
		}
	})

	t.Run("NewCirculationSpeedImprovementID is unique", func(t *testing.T) {
		t.Parallel()
		if NewCirculationSpeedImprovementID() == NewCirculationSpeedImprovementID() {
			t.Fatal("expected unique IDs")
		}
	})

	t.Run("NewAnnualSurplusPenaltyID is unique", func(t *testing.T) {
		t.Parallel()
		if NewAnnualSurplusPenaltyID() == NewAnnualSurplusPenaltyID() {
			t.Fatal("expected unique IDs")
		}
	})
}

// TestEnglandIndiaFixture validates the key quantitative result from Marx's text:
// the 1870 telegraph collapsed the London-Bombay communication lag from ~250 days
// to 1 day — a ratio of 58:1.
func TestEnglandIndiaFixture(t *testing.T) {
	t.Parallel()

	telegraph := CirculationSpeedImprovement{
		MarketID:           "london-bombay",
		OldLagMinutesPerKm: 58,
		NewLagMinutesPerKm: 1,
		Reason:             ReasonTelegraph,
		EffectiveAt:        time.Date(1870, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	if telegraph.LagReduction() != 57 {
		t.Fatalf("expected telegraph lag reduction of 57 min/km, got %d", telegraph.LagReduction())
	}

	// 9000 km * 58 min/km = 522 000 min pre-telegraph.
	// 9000 km * 1 min/km = 9 000 min post-telegraph.
	pre := DistanceLagRelation{DistanceMeters: 9000000, LagMinutesPerKm: 58}
	post := DistanceLagRelation{DistanceMeters: 9000000, LagMinutesPerKm: 1}

	if pre.MaxLagMinutes() <= post.MaxLagMinutes() {
		t.Fatal("post-telegraph max lag should be less than pre-telegraph")
	}
	ratio := pre.MaxLagMinutes() / post.MaxLagMinutes()
	if ratio != 58 {
		t.Fatalf("expected lag ratio of 58 (pre/post telegraph), got %d", ratio)
	}
}
