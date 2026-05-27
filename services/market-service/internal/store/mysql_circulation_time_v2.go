package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/circulation"
)

// --- SellingPhaseDetailed ---------------------------------------------------

func (m *MySQL) CreateSellingPhaseDetailed(ctx context.Context, d circulation.SellingPhaseDetailed) (circulation.SellingPhaseDetailed, error) {
	if d.ID.IsZero() {
		d.ID = circulation.NewSellingPhaseDetailedID()
	}
	d.CreatedAt = m.now().UTC()
	const q = `
		INSERT INTO selling_phase_details
		    (id, selling_phase_id, turnover_time_id, industrial_capital_id,
		     market_id, market_distance_meters, communication_lag_minutes,
		     buyer_search_minutes, bargaining_minutes, legal_settlement_minutes,
		     total_minutes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(d.ID), string(d.SellingPhaseID), string(d.TurnoverTimeID),
		string(d.IndustrialCapitalID), string(d.MarketID),
		d.MarketDistanceMeters, d.CommunicationLagMinutes,
		d.BuyerSearchMinutes, d.BargainingMinutes, d.LegalSettlementMinutes,
		d.TotalMinutes, d.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return circulation.SellingPhaseDetailed{}, ErrAlreadyExists
		}
		return circulation.SellingPhaseDetailed{}, err
	}
	return d, nil
}

func (m *MySQL) GetSellingPhaseDetailed(ctx context.Context, id circulation.SellingPhaseDetailedID) (circulation.SellingPhaseDetailed, error) {
	const q = `
		SELECT id, selling_phase_id, turnover_time_id, industrial_capital_id,
		       market_id, market_distance_meters, communication_lag_minutes,
		       buyer_search_minutes, bargaining_minutes, legal_settlement_minutes,
		       total_minutes, created_at
		FROM selling_phase_details WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanSellingPhaseDetailed(row)
}

func scanSellingPhaseDetailed(row *sql.Row) (circulation.SellingPhaseDetailed, error) {
	var d circulation.SellingPhaseDetailed
	var (
		sellingPhaseID      string
		turnoverTimeID      string
		industrialCapitalID string
		marketID            string
		createdAt           time.Time
	)
	err := row.Scan(
		(*string)(&d.ID), &sellingPhaseID, &turnoverTimeID, &industrialCapitalID,
		&marketID, &d.MarketDistanceMeters, &d.CommunicationLagMinutes,
		&d.BuyerSearchMinutes, &d.BargainingMinutes, &d.LegalSettlementMinutes,
		&d.TotalMinutes, &createdAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return circulation.SellingPhaseDetailed{}, ErrNotFound
	}
	if err != nil {
		return circulation.SellingPhaseDetailed{}, err
	}
	d.SellingPhaseID = circulation.SellingPhaseID(sellingPhaseID)
	d.TurnoverTimeID = circulation.TurnoverTimeID(turnoverTimeID)
	d.IndustrialCapitalID = circulation.IndustrialCapitalID(industrialCapitalID)
	d.MarketID = circulation.MarketID(marketID)
	d.CreatedAt = createdAt
	return d, nil
}

// --- BuyingPhaseDetailed ----------------------------------------------------

func (m *MySQL) CreateBuyingPhaseDetailed(ctx context.Context, d circulation.BuyingPhaseDetailed) (circulation.BuyingPhaseDetailed, error) {
	if d.ID.IsZero() {
		d.ID = circulation.NewBuyingPhaseDetailedID()
	}
	d.CreatedAt = m.now().UTC()
	const q = `
		INSERT INTO buying_phase_details
		    (id, buying_phase_id, turnover_time_id, industrial_capital_id,
		     market_id, market_distance_meters, communication_lag_minutes,
		     supplier_search_minutes, order_processing_minutes, legal_settlement_minutes,
		     total_minutes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(d.ID), string(d.BuyingPhaseID), string(d.TurnoverTimeID),
		string(d.IndustrialCapitalID), string(d.MarketID),
		d.MarketDistanceMeters, d.CommunicationLagMinutes,
		d.SupplierSearchMinutes, d.OrderProcessingMinutes, d.LegalSettlementMinutes,
		d.TotalMinutes, d.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return circulation.BuyingPhaseDetailed{}, ErrAlreadyExists
		}
		return circulation.BuyingPhaseDetailed{}, err
	}
	return d, nil
}

func (m *MySQL) GetBuyingPhaseDetailed(ctx context.Context, id circulation.BuyingPhaseDetailedID) (circulation.BuyingPhaseDetailed, error) {
	const q = `
		SELECT id, buying_phase_id, turnover_time_id, industrial_capital_id,
		       market_id, market_distance_meters, communication_lag_minutes,
		       supplier_search_minutes, order_processing_minutes, legal_settlement_minutes,
		       total_minutes, created_at
		FROM buying_phase_details WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanBuyingPhaseDetailed(row)
}

func scanBuyingPhaseDetailed(row *sql.Row) (circulation.BuyingPhaseDetailed, error) {
	var d circulation.BuyingPhaseDetailed
	var (
		buyingPhaseID       string
		turnoverTimeID      string
		industrialCapitalID string
		marketID            string
		createdAt           time.Time
	)
	err := row.Scan(
		(*string)(&d.ID), &buyingPhaseID, &turnoverTimeID, &industrialCapitalID,
		&marketID, &d.MarketDistanceMeters, &d.CommunicationLagMinutes,
		&d.SupplierSearchMinutes, &d.OrderProcessingMinutes, &d.LegalSettlementMinutes,
		&d.TotalMinutes, &createdAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return circulation.BuyingPhaseDetailed{}, ErrNotFound
	}
	if err != nil {
		return circulation.BuyingPhaseDetailed{}, err
	}
	d.BuyingPhaseID = circulation.BuyingPhaseID(buyingPhaseID)
	d.TurnoverTimeID = circulation.TurnoverTimeID(turnoverTimeID)
	d.IndustrialCapitalID = circulation.IndustrialCapitalID(industrialCapitalID)
	d.MarketID = circulation.MarketID(marketID)
	d.CreatedAt = createdAt
	return d, nil
}

// --- DistanceLagRelation ----------------------------------------------------

func (m *MySQL) CreateDistanceLagRelation(ctx context.Context, r circulation.DistanceLagRelation) (circulation.DistanceLagRelation, error) {
	if r.ID.IsZero() {
		r.ID = circulation.NewDistanceLagRelationID()
	}
	r.CreatedAt = m.now().UTC()
	const q = `
		INSERT INTO distance_lag_relations
		    (id, market_id, distance_meters, lag_minutes_per_km, created_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(r.ID), string(r.MarketID), r.DistanceMeters, r.LagMinutesPerKm, r.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return circulation.DistanceLagRelation{}, ErrAlreadyExists
		}
		return circulation.DistanceLagRelation{}, err
	}
	return r, nil
}

func (m *MySQL) GetDistanceLagRelation(ctx context.Context, id circulation.DistanceLagRelationID) (circulation.DistanceLagRelation, error) {
	const q = `SELECT id, market_id, distance_meters, lag_minutes_per_km, created_at FROM distance_lag_relations WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	var r circulation.DistanceLagRelation
	var marketID string
	err := row.Scan((*string)(&r.ID), &marketID, &r.DistanceMeters, &r.LagMinutesPerKm, &r.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return circulation.DistanceLagRelation{}, ErrNotFound
	}
	if err != nil {
		return circulation.DistanceLagRelation{}, err
	}
	r.MarketID = circulation.MarketID(marketID)
	return r, nil
}

func (m *MySQL) ListDistanceLagRelations(ctx context.Context) ([]circulation.DistanceLagRelation, error) {
	const q = `SELECT id, market_id, distance_meters, lag_minutes_per_km, created_at FROM distance_lag_relations ORDER BY created_at`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []circulation.DistanceLagRelation{}
	for rows.Next() {
		var r circulation.DistanceLagRelation
		var marketID string
		if err := rows.Scan((*string)(&r.ID), &marketID, &r.DistanceMeters, &r.LagMinutesPerKm, &r.CreatedAt); err != nil {
			return nil, err
		}
		r.MarketID = circulation.MarketID(marketID)
		out = append(out, r)
	}
	return out, rows.Err()
}

// --- CirculationSpeedImprovement --------------------------------------------

func (m *MySQL) CreateCirculationSpeedImprovement(ctx context.Context, csi circulation.CirculationSpeedImprovement) (circulation.CirculationSpeedImprovement, error) {
	if csi.ID.IsZero() {
		csi.ID = circulation.NewCirculationSpeedImprovementID()
	}
	csi.CreatedAt = m.now().UTC()
	const q = `
		INSERT INTO circulation_speed_improvements
		    (id, market_id, old_lag_minutes_per_km, new_lag_minutes_per_km, reason, effective_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(csi.ID), string(csi.MarketID),
		csi.OldLagMinutesPerKm, csi.NewLagMinutesPerKm,
		string(csi.Reason), csi.EffectiveAt, csi.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return circulation.CirculationSpeedImprovement{}, ErrAlreadyExists
		}
		return circulation.CirculationSpeedImprovement{}, err
	}
	return csi, nil
}

func (m *MySQL) ListCirculationSpeedImprovements(ctx context.Context, marketID circulation.MarketID) ([]circulation.CirculationSpeedImprovement, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if marketID == "" {
		const q = `SELECT id, market_id, old_lag_minutes_per_km, new_lag_minutes_per_km, reason, effective_at, created_at FROM circulation_speed_improvements ORDER BY effective_at`
		rows, err = m.db.QueryContext(ctx, q)
	} else {
		const q = `SELECT id, market_id, old_lag_minutes_per_km, new_lag_minutes_per_km, reason, effective_at, created_at FROM circulation_speed_improvements WHERE market_id = ? ORDER BY effective_at`
		rows, err = m.db.QueryContext(ctx, q, string(marketID))
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []circulation.CirculationSpeedImprovement{}
	for rows.Next() {
		var csi circulation.CirculationSpeedImprovement
		var mid, reason string
		if err := rows.Scan((*string)(&csi.ID), &mid, &csi.OldLagMinutesPerKm, &csi.NewLagMinutesPerKm, &reason, &csi.EffectiveAt, &csi.CreatedAt); err != nil {
			return nil, err
		}
		csi.MarketID = circulation.MarketID(mid)
		csi.Reason = circulation.CirculationSpeedReason(reason)
		out = append(out, csi)
	}
	return out, rows.Err()
}

// --- AnnualSurplusPenalty ---------------------------------------------------

func (m *MySQL) CreateAnnualSurplusPenalty(ctx context.Context, asp circulation.AnnualSurplusPenalty) (circulation.AnnualSurplusPenalty, error) {
	if asp.ID.IsZero() {
		asp.ID = circulation.NewAnnualSurplusPenaltyID()
	}
	asp.PenaltyBasisPoints = asp.ComputePenalty()
	const q = `
		INSERT INTO annual_surplus_penalties
		    (id, industrial_capital_id, period, ideal_annual_rate_basis_points,
		     actual_annual_rate_basis_points, penalty_basis_points)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
		    ideal_annual_rate_basis_points  = VALUES(ideal_annual_rate_basis_points),
		    actual_annual_rate_basis_points = VALUES(actual_annual_rate_basis_points),
		    penalty_basis_points            = VALUES(penalty_basis_points)`
	_, err := m.db.ExecContext(ctx, q,
		string(asp.ID), string(asp.IndustrialCapitalID), asp.Period,
		asp.IdealAnnualRateBasisPoints, asp.ActualAnnualRateBasisPoints, asp.PenaltyBasisPoints,
	)
	if err != nil {
		return circulation.AnnualSurplusPenalty{}, err
	}
	return asp, nil
}

func (m *MySQL) GetAnnualSurplusPenalty(ctx context.Context, icID circulation.IndustrialCapitalID, period string) (circulation.AnnualSurplusPenalty, error) {
	const q = `
		SELECT id, industrial_capital_id, period, ideal_annual_rate_basis_points,
		       actual_annual_rate_basis_points, penalty_basis_points
		FROM annual_surplus_penalties
		WHERE industrial_capital_id = ? AND period = ?`
	row := m.db.QueryRowContext(ctx, q, string(icID), period)
	var asp circulation.AnnualSurplusPenalty
	var icIDStr string
	err := row.Scan(
		(*string)(&asp.ID), &icIDStr, &asp.Period,
		&asp.IdealAnnualRateBasisPoints, &asp.ActualAnnualRateBasisPoints, &asp.PenaltyBasisPoints,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return circulation.AnnualSurplusPenalty{}, ErrNotFound
	}
	if err != nil {
		return circulation.AnnualSurplusPenalty{}, err
	}
	asp.IndustrialCapitalID = circulation.IndustrialCapitalID(icIDStr)
	return asp, nil
}
