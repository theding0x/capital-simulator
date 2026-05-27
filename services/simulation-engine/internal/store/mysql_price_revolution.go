package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

// CreatePriceRevolution inserts a ValueRevolutionEvent plus its derived
// RevaluedCapital or RealisedValueDelta row in a single transaction.
func (m *MySQL) CreatePriceRevolution(ctx context.Context, rec circulation.PriceRevolutionRecord) (circulation.PriceRevolutionRecord, error) {
	evt := rec.Event
	if evt.ID == "" {
		evt.ID = circulation.NewValueRevolutionEventID()
	}
	if evt.OccurredAt.IsZero() {
		evt.OccurredAt = m.now().UTC()
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return circulation.PriceRevolutionRecord{}, err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO value_revolution_events
		   (id, industrial_capital_id, affecting, direction_pence, occurred_at,
		    basis_points_change, revalued_capital_pence)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		string(evt.ID), string(evt.IndustrialCapitalID),
		string(evt.Affecting), int64(evt.DirectionPence),
		evt.OccurredAt,
		evt.BasisPointsChange, int64(evt.RevaluedCapitalPence),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return circulation.PriceRevolutionRecord{}, ErrAlreadyExists
		}
		return circulation.PriceRevolutionRecord{}, err
	}
	rec.Event = evt

	if rec.RevaluedCapital != nil {
		rc := *rec.RevaluedCapital
		if rc.ID == "" {
			rc.ID = circulation.NewRevaluedCapitalID()
		}
		_, err = tx.ExecContext(ctx,
			`INSERT INTO revalued_capitals
			   (id, value_revolution_event_id, industrial_capital_id,
			    old_advance_pence, new_advance_pence, delta_pence)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			string(rc.ID), string(rc.ValueRevolutionEventID),
			string(rc.IndustrialCapitalID),
			int64(rc.OldAdvancePence), int64(rc.NewAdvancePence), int64(rc.DeltaPence),
		)
		if err != nil {
			return circulation.PriceRevolutionRecord{}, err
		}
		rec.RevaluedCapital = &rc
	}

	if rec.RealisedDelta != nil {
		rvd := *rec.RealisedDelta
		if rvd.ID == "" {
			rvd.ID = circulation.NewRealisedValueDeltaID()
		}
		_, err = tx.ExecContext(ctx,
			`INSERT INTO realised_value_deltas
			   (id, value_revolution_event_id, industrial_capital_id,
			    expected_realisation_pence, actual_realisation_pence, delta_pence)
			 VALUES (?, ?, ?, ?, ?, ?)`,
			string(rvd.ID), string(rvd.ValueRevolutionEventID),
			string(rvd.IndustrialCapitalID),
			int64(rvd.ExpectedRealisationPence), int64(rvd.ActualRealisationPence),
			int64(rvd.DeltaPence),
		)
		if err != nil {
			return circulation.PriceRevolutionRecord{}, err
		}
		rec.RealisedDelta = &rvd
	}

	if err := tx.Commit(); err != nil {
		return circulation.PriceRevolutionRecord{}, err
	}
	if rec.Revaluations == nil {
		rec.Revaluations = []circulation.InventoryRevaluation{}
	}
	return rec, nil
}

// GetPriceRevolution retrieves an event with all derived rows.
func (m *MySQL) GetPriceRevolution(ctx context.Context, id circulation.ValueRevolutionEventID) (circulation.PriceRevolutionRecord, error) {
	row := m.db.QueryRowContext(ctx,
		`SELECT id, industrial_capital_id, affecting, direction_pence, occurred_at,
		        basis_points_change, revalued_capital_pence
		 FROM value_revolution_events WHERE id = ?`, string(id))
	var evt circulation.ValueRevolutionEvent
	var evtID, icID, affecting string
	var occurredAt time.Time
	if err := row.Scan(&evtID, &icID, &affecting,
		(*int64)(&evt.DirectionPence), &occurredAt,
		&evt.BasisPointsChange, (*int64)(&evt.RevaluedCapitalPence)); err != nil {
		if err == sql.ErrNoRows {
			return circulation.PriceRevolutionRecord{}, ErrNotFound
		}
		return circulation.PriceRevolutionRecord{}, err
	}
	evt.ID = circulation.ValueRevolutionEventID(evtID)
	evt.IndustrialCapitalID = circulation.IndustrialCapitalID(icID)
	evt.Affecting = circulation.ValueRevolutionAffecting(affecting)
	evt.OccurredAt = occurredAt.UTC()

	rec := circulation.PriceRevolutionRecord{
		Event:        evt,
		Revaluations: []circulation.InventoryRevaluation{},
	}

	// RevaluedCapital
	rcRow := m.db.QueryRowContext(ctx,
		`SELECT id, value_revolution_event_id, industrial_capital_id,
		        old_advance_pence, new_advance_pence, delta_pence
		 FROM revalued_capitals WHERE value_revolution_event_id = ?`, string(id))
	var rc circulation.RevaluedCapital
	var rcID, rcEvtID, rcIcID string
	if err := rcRow.Scan(&rcID, &rcEvtID, &rcIcID,
		(*int64)(&rc.OldAdvancePence), (*int64)(&rc.NewAdvancePence),
		(*int64)(&rc.DeltaPence)); err == nil {
		rc.ID = circulation.RevaluedCapitalID(rcID)
		rc.ValueRevolutionEventID = circulation.ValueRevolutionEventID(rcEvtID)
		rc.IndustrialCapitalID = circulation.IndustrialCapitalID(rcIcID)
		rec.RevaluedCapital = &rc
	}

	// RealisedValueDelta
	rvdRow := m.db.QueryRowContext(ctx,
		`SELECT id, value_revolution_event_id, industrial_capital_id,
		        expected_realisation_pence, actual_realisation_pence, delta_pence
		 FROM realised_value_deltas WHERE value_revolution_event_id = ?`, string(id))
	var rvd circulation.RealisedValueDelta
	var rvdID, rvdEvtID, rvdIcID string
	if err := rvdRow.Scan(&rvdID, &rvdEvtID, &rvdIcID,
		(*int64)(&rvd.ExpectedRealisationPence), (*int64)(&rvd.ActualRealisationPence),
		(*int64)(&rvd.DeltaPence)); err == nil {
		rvd.ID = circulation.RealisedValueDeltaID(rvdID)
		rvd.ValueRevolutionEventID = circulation.ValueRevolutionEventID(rvdEvtID)
		rvd.IndustrialCapitalID = circulation.IndustrialCapitalID(rvdIcID)
		rec.RealisedDelta = &rvd
	}

	// InventoryRevaluations
	rows, err := m.db.QueryContext(ctx,
		`SELECT id, industrial_capital_id, commodity_id,
		        quantity_held, old_unit_pence, new_unit_pence, delta_pence
		 FROM inventory_revaluations WHERE industrial_capital_id = ?`,
		string(evt.IndustrialCapitalID))
	if err != nil {
		return circulation.PriceRevolutionRecord{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var inv circulation.InventoryRevaluation
		var invID, invIcID, commodityID string
		if err := rows.Scan(&invID, &invIcID, &commodityID,
			&inv.QuantityHeld,
			(*int64)(&inv.OldUnitPence), (*int64)(&inv.NewUnitPence),
			(*int64)(&inv.DeltaPence)); err != nil {
			return circulation.PriceRevolutionRecord{}, err
		}
		inv.ID = circulation.InventoryRevaluationID(invID)
		inv.IndustrialCapitalID = circulation.IndustrialCapitalID(invIcID)
		inv.CommodityID = commodityID
		rec.Revaluations = append(rec.Revaluations, inv)
	}

	// Price case
	caseRow := m.db.QueryRowContext(ctx,
		`SELECT fall_case, rise_case FROM price_case_records
		 WHERE value_revolution_event_id = ?`, string(id))
	var fallCase, riseCase sql.NullString
	if err := caseRow.Scan(&fallCase, &riseCase); err == nil {
		if fallCase.Valid {
			fc := circulation.PriceFallCase(fallCase.String)
			rec.FallCase = &fc
		}
		if riseCase.Valid {
			rcase := circulation.PriceRiseCase(riseCase.String)
			rec.RiseCase = &rcase
		}
	}

	return rec, nil
}

// RecordPriceCase inserts or replaces the case for an event.
func (m *MySQL) RecordPriceCase(ctx context.Context, pc circulation.PriceCaseRecord) (circulation.PriceCaseRecord, error) {
	if err := pc.Validate(); err != nil {
		return circulation.PriceCaseRecord{}, err
	}
	var fallCase, riseCase interface{}
	if pc.FallCase != nil {
		fallCase = string(*pc.FallCase)
	}
	if pc.RiseCase != nil {
		riseCase = string(*pc.RiseCase)
	}
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO price_case_records (value_revolution_event_id, fall_case, rise_case)
		 VALUES (?, ?, ?)
		 ON DUPLICATE KEY UPDATE fall_case = VALUES(fall_case), rise_case = VALUES(rise_case)`,
		string(pc.ValueRevolutionEventID), fallCase, riseCase,
	)
	if err != nil {
		return circulation.PriceCaseRecord{}, err
	}
	return pc, nil
}

// RecordInventoryRevaluation inserts an InventoryRevaluation row.
func (m *MySQL) RecordInventoryRevaluation(ctx context.Context, inv circulation.InventoryRevaluation) (circulation.InventoryRevaluation, error) {
	if err := inv.Validate(); err != nil {
		return circulation.InventoryRevaluation{}, err
	}
	if inv.ID == "" {
		inv.ID = circulation.NewInventoryRevaluationID()
	}
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO inventory_revaluations
		   (id, industrial_capital_id, commodity_id,
		    quantity_held, old_unit_pence, new_unit_pence, delta_pence)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		string(inv.ID), string(inv.IndustrialCapitalID), inv.CommodityID,
		inv.QuantityHeld,
		int64(inv.OldUnitPence), int64(inv.NewUnitPence), int64(inv.DeltaPence),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return circulation.InventoryRevaluation{}, ErrAlreadyExists
		}
		return circulation.InventoryRevaluation{}, err
	}
	return inv, nil
}

// RecordSpeculativeHold inserts a SpeculativeHold row.
func (m *MySQL) RecordSpeculativeHold(ctx context.Context, sh circulation.SpeculativeHold) (circulation.SpeculativeHold, error) {
	if err := sh.Validate(); err != nil {
		return circulation.SpeculativeHold{}, err
	}
	if sh.ID == "" {
		sh.ID = circulation.NewSpeculativeHoldID()
	}
	if sh.HeldSince.IsZero() {
		sh.HeldSince = m.now().UTC()
	}
	_, err := m.db.ExecContext(ctx,
		`INSERT INTO speculative_holds
		   (id, industrial_capital_id, commodity_id, held_since, anticipated_direction)
		 VALUES (?, ?, ?, ?, ?)`,
		string(sh.ID), string(sh.IndustrialCapitalID), sh.CommodityID,
		sh.HeldSince, int(sh.AnticipatedDirection),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return circulation.SpeculativeHold{}, ErrAlreadyExists
		}
		return circulation.SpeculativeHold{}, err
	}
	return sh, nil
}

// AggregateCompound sums all delta streams for a capital+period.
func (m *MySQL) AggregateCompound(ctx context.Context, icID circulation.IndustrialCapitalID, _ string) (circulation.CompoundCapitalAdjustment, error) {
	var revDelta, reaDelta, invDelta int64

	row := m.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(rc.delta_pence), 0)
		 FROM revalued_capitals rc
		 JOIN value_revolution_events vre ON vre.id = rc.value_revolution_event_id
		 WHERE vre.industrial_capital_id = ?`,
		string(icID))
	if err := row.Scan(&revDelta); err != nil {
		return circulation.CompoundCapitalAdjustment{}, err
	}

	row = m.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(rvd.delta_pence), 0)
		 FROM realised_value_deltas rvd
		 JOIN value_revolution_events vre ON vre.id = rvd.value_revolution_event_id
		 WHERE vre.industrial_capital_id = ?`,
		string(icID))
	if err := row.Scan(&reaDelta); err != nil {
		return circulation.CompoundCapitalAdjustment{}, err
	}

	row = m.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(delta_pence), 0)
		 FROM inventory_revaluations
		 WHERE industrial_capital_id = ?`,
		string(icID))
	if err := row.Scan(&invDelta); err != nil {
		return circulation.CompoundCapitalAdjustment{}, err
	}

	net := revDelta + reaDelta + invDelta
	return circulation.CompoundCapitalAdjustment{
		ID:                    circulation.NewCompoundCapitalAdjustmentID(),
		IndustrialCapitalID:   icID,
		RevaluationDeltaPence: circulation.Pence(revDelta),
		RealisationDeltaPence: circulation.Pence(reaDelta),
		InventoryDeltaPence:   circulation.Pence(invDelta),
		NetEffectPence:        circulation.Pence(net),
	}, nil
}
