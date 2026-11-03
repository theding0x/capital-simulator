package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/store"
	wp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/workperiod"
)

const (
	day          = int64(86_400_000_000_000)
	workingDay   = int64(600) * 60_000_000_000 // 10-hour working-day in ns
	cornPeriodID = wp.WorkingPeriodID("5eed000000000000001301")
)

// seedCornPeriod stores Marx's corn fixture: 10 connected labour-days inside a
// 365-day production span. The Ch. 12 annual-crop floor (365 days) is recorded
// after the period, mirroring the seeds, which load the labour-sized period by
// SQL; the floor then binds the gap's production span.
func seedCornPeriod(t *testing.T, st *store.Memory) {
	t.Helper()
	ctx := context.Background()
	if _, err := st.CreateWorkingPeriod(ctx, wp.WorkingPeriod{
		ID: cornPeriodID, IndustrialCapitalID: "corn-farm-1871-ic", CommodityID: "corn",
		WorkingDayCount: 10, Mode: wp.ModeConnected, WorkingDayLengthMinutes: 600,
	}); err != nil {
		t.Fatalf("period: %v", err)
	}
	if _, err := st.CreateNaturalConstraint(ctx, wp.NaturalWorkingPeriodConstraint{
		CommodityID: "corn", MinWorkingDays: 365, Reason: "annual crop",
	}); err != nil {
		t.Fatalf("constraint: %v", err)
	}
}

func cornGap(production int64) wp.ProductionLabourGap {
	labour := 10 * workingDay
	return wp.ProductionLabourGap{
		WorkingPeriodID: cornPeriodID, ProductionTimeNanos: production,
		LabourTimeNanos: labour, GapNanos: production - labour,
	}
}

func TestMemory_RecordProductionLabourGap_RoundTrip(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	seedCornPeriod(t, st)
	ctx := context.Background()

	created, err := st.RecordProductionLabourGap(ctx, cornGap(365*day))
	if err != nil {
		t.Fatalf("record: %v", err)
	}
	got, err := st.GetProductionLabourGap(ctx, cornPeriodID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got != created {
		t.Errorf("round-trip: got %+v, want %+v", got, created)
	}
	// Gap time creates no value: the owning period's value-bearing field is
	// still exactly the labour-days.
	period, _ := st.GetWorkingPeriod(ctx, cornPeriodID)
	if period.WorkingDayCount != 10 {
		t.Errorf("WorkingDayCount changed by recording a gap: got %d, want 10", period.WorkingDayCount)
	}
}

func TestMemory_RecordProductionLabourGap_UnknownPeriod(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	_, err := st.RecordProductionLabourGap(context.Background(), cornGap(365*day))
	if !errors.Is(err, store.ErrNotFound) {
		t.Errorf("got %v, want ErrNotFound", err)
	}
}

func TestMemory_RecordProductionLabourGap_LabourMustMatchPeriod(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	seedCornPeriod(t, st)
	// 10 labour-days counted as 24-hour days (the old seed convention) is not
	// the period's 10 ten-hour working-days.
	g := cornGap(365 * day)
	g.LabourTimeNanos = 10 * day
	g.GapNanos = g.ProductionTimeNanos - g.LabourTimeNanos
	_, err := st.RecordProductionLabourGap(context.Background(), g)
	if !errors.Is(err, wp.ErrGapLabourMismatch) {
		t.Errorf("got %v, want ErrGapLabourMismatch", err)
	}
}

func TestMemory_RecordProductionLabourGap_ProductionSpanHonoursFloor(t *testing.T) {
	t.Parallel()
	st := store.NewMemory()
	seedCornPeriod(t, st)
	// The pre-fix seed's 280-day corn span is below the 365-day floor.
	_, err := st.RecordProductionLabourGap(context.Background(), cornGap(280*day))
	if !errors.Is(err, wp.ErrPrematureDelivery) {
		t.Errorf("got %v, want ErrPrematureDelivery", err)
	}
}
