// Vol. II Ch. 6 — MySQL implementation of CirculationCostStore.
package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/theding0x/capital-simulator/services/market-service/internal/costs"
)

// --- CirculationCost --------------------------------------------------------

func (m *MySQL) CreateCirculationCost(ctx context.Context, c costs.CirculationCost) (costs.CirculationCost, error) {
	if c.ID == "" {
		c.ID = costs.NewCirculationCostID()
	}
	if c.IncurredAt.IsZero() {
		c.IncurredAt = m.now().UTC()
	}
	c.Nature = costs.NatureOf(c.Kind)
	const q = `INSERT INTO circulation_costs (id, industrial_capital_id, kind, nature, pence, incurred_at)
		VALUES (?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID), string(c.IndustrialCapitalID),
		string(c.Kind), string(c.Nature),
		int64(c.Pence), c.IncurredAt)
	if err != nil {
		return costs.CirculationCost{}, err
	}
	return c, nil
}

func (m *MySQL) GetCirculationCost(ctx context.Context, id costs.CirculationCostID) (costs.CirculationCost, error) {
	const q = `SELECT id, industrial_capital_id, kind, nature, pence, incurred_at
		FROM circulation_costs WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanCirculationCost(row)
}

func (m *MySQL) ListCirculationCosts(ctx context.Context, icID costs.IndustrialCapitalID, kind costs.CirculationCostKind, nature costs.CostNature) ([]costs.CirculationCost, error) {
	q := `SELECT id, industrial_capital_id, kind, nature, pence, incurred_at
		FROM circulation_costs WHERE 1=1`
	args := []any{}
	if icID != "" {
		q += " AND industrial_capital_id = ?"
		args = append(args, string(icID))
	}
	if kind != "" {
		q += " AND kind = ?"
		args = append(args, string(kind))
	}
	if nature != "" {
		q += " AND nature = ?"
		args = append(args, string(nature))
	}
	q += " ORDER BY incurred_at ASC"
	rows, err := m.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []costs.CirculationCost{}
	for rows.Next() {
		c, err := scanCirculationCostRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func scanCirculationCost(row *sql.Row) (costs.CirculationCost, error) {
	var c costs.CirculationCost
	var id, icID, kind, nature string
	var pence int64
	var incurredAt time.Time
	if err := row.Scan(&id, &icID, &kind, &nature, &pence, &incurredAt); err != nil {
		if isNotFound(err) {
			return costs.CirculationCost{}, ErrNotFound
		}
		return costs.CirculationCost{}, err
	}
	c.ID = costs.CirculationCostID(id)
	c.IndustrialCapitalID = costs.IndustrialCapitalID(icID)
	c.Kind = costs.CirculationCostKind(kind)
	c.Nature = costs.CostNature(nature)
	c.Pence = costs.Pence(pence)
	c.IncurredAt = incurredAt
	return c, nil
}

func scanCirculationCostRow(rows *sql.Rows) (costs.CirculationCost, error) {
	var c costs.CirculationCost
	var id, icID, kind, nature string
	var pence int64
	var incurredAt time.Time
	if err := rows.Scan(&id, &icID, &kind, &nature, &pence, &incurredAt); err != nil {
		return costs.CirculationCost{}, err
	}
	c.ID = costs.CirculationCostID(id)
	c.IndustrialCapitalID = costs.IndustrialCapitalID(icID)
	c.Kind = costs.CirculationCostKind(kind)
	c.Nature = costs.CostNature(nature)
	c.Pence = costs.Pence(pence)
	c.IncurredAt = incurredAt
	return c, nil
}

// --- CirculationAgent -------------------------------------------------------

func (m *MySQL) CreateCirculationAgent(ctx context.Context, a costs.CirculationAgent) (costs.CirculationAgent, error) {
	if a.ID == "" {
		a.ID = costs.NewCirculationAgentID()
	}
	const q = `INSERT INTO circulation_agents
		(id, industrial_capital_id, agent_id, role, wage_rate_pence, labour_minutes_necessary, labour_minutes_surplus)
		VALUES (?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID), string(a.IndustrialCapitalID), string(a.AgentID),
		string(a.Role), int64(a.WageRatePence),
		a.LabourMinutesNecessary, a.LabourMinutesSurplus)
	if err != nil {
		return costs.CirculationAgent{}, err
	}
	return a, nil
}

func (m *MySQL) ListCirculationAgents(ctx context.Context, icID costs.IndustrialCapitalID) ([]costs.CirculationAgent, error) {
	q := `SELECT id, industrial_capital_id, agent_id, role, wage_rate_pence, labour_minutes_necessary, labour_minutes_surplus
		FROM circulation_agents WHERE 1=1`
	args := []any{}
	if icID != "" {
		q += " AND industrial_capital_id = ?"
		args = append(args, string(icID))
	}
	rows, err := m.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []costs.CirculationAgent{}
	for rows.Next() {
		var a costs.CirculationAgent
		var id, icid, agid, role string
		var wage, necessary, surplus int64
		if err := rows.Scan(&id, &icid, &agid, &role, &wage, &necessary, &surplus); err != nil {
			return nil, err
		}
		a.ID = costs.CirculationAgentID(id)
		a.IndustrialCapitalID = costs.IndustrialCapitalID(icid)
		a.AgentID = costs.AgentID(agid)
		a.Role = costs.CirculationAgentRole(role)
		a.WageRatePence = costs.Pence(wage)
		a.LabourMinutesNecessary = necessary
		a.LabourMinutesSurplus = surplus
		out = append(out, a)
	}
	return out, rows.Err()
}

// --- MoneyAsFauxFrais -------------------------------------------------------

func (m *MySQL) CreateMoneyAsFauxFrais(ctx context.Context, mf costs.MoneyAsFauxFrais) (costs.MoneyAsFauxFrais, error) {
	if mf.ID == "" {
		mf.ID = costs.NewMoneyAsFauxFraisID()
	}
	const q = `INSERT INTO money_as_faux_frais (id, industrial_capital_id, pence, annual_replacement_pence)
		VALUES (?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(mf.ID), string(mf.IndustrialCapitalID),
		int64(mf.Pence), int64(mf.AnnualReplacementPence))
	if err != nil {
		return costs.MoneyAsFauxFrais{}, err
	}
	return mf, nil
}

func (m *MySQL) ListMoneyAsFauxFrais(ctx context.Context) ([]costs.MoneyAsFauxFrais, error) {
	const q = `SELECT id, industrial_capital_id, pence, annual_replacement_pence FROM money_as_faux_frais`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []costs.MoneyAsFauxFrais{}
	for rows.Next() {
		var mf costs.MoneyAsFauxFrais
		var id, icID string
		var pence, annual int64
		if err := rows.Scan(&id, &icID, &pence, &annual); err != nil {
			return nil, err
		}
		mf.ID = costs.MoneyAsFauxFraisID(id)
		mf.IndustrialCapitalID = costs.IndustrialCapitalID(icID)
		mf.Pence = costs.Pence(pence)
		mf.AnnualReplacementPence = costs.Pence(annual)
		out = append(out, mf)
	}
	return out, rows.Err()
}

// --- CommoditySupply --------------------------------------------------------

func (m *MySQL) CreateCommoditySupply(ctx context.Context, s costs.CommoditySupply) (costs.CommoditySupply, error) {
	if s.ID == "" {
		s.ID = costs.NewCommoditySupplyID()
	}
	if s.HeldSince.IsZero() {
		s.HeldSince = m.now().UTC()
	}
	const q = `INSERT INTO commodity_supplies (id, industrial_capital_id, commodity_id, pence, held_since, is_voluntary)
		VALUES (?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(s.ID), string(s.IndustrialCapitalID), string(s.CommodityID),
		int64(s.Pence), s.HeldSince, s.IsVoluntary)
	if err != nil {
		return costs.CommoditySupply{}, err
	}
	return s, nil
}

func (m *MySQL) GetCommoditySupply(ctx context.Context, id costs.CommoditySupplyID) (costs.CommoditySupply, error) {
	const q = `SELECT id, industrial_capital_id, commodity_id, pence, held_since, is_voluntary
		FROM commodity_supplies WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	var s costs.CommoditySupply
	var sid, icID, comID string
	var pence int64
	var isVol bool
	var heldSince time.Time
	if err := row.Scan(&sid, &icID, &comID, &pence, &heldSince, &isVol); err != nil {
		if isNotFound(err) {
			return costs.CommoditySupply{}, ErrNotFound
		}
		return costs.CommoditySupply{}, err
	}
	s.ID = costs.CommoditySupplyID(sid)
	s.IndustrialCapitalID = costs.IndustrialCapitalID(icID)
	s.CommodityID = costs.CommodityID(comID)
	s.Pence = costs.Pence(pence)
	s.HeldSince = heldSince
	s.IsVoluntary = isVol
	return s, nil
}

func (m *MySQL) ListCommoditySupplies(ctx context.Context, icID costs.IndustrialCapitalID) ([]costs.CommoditySupply, error) {
	q := `SELECT id, industrial_capital_id, commodity_id, pence, held_since, is_voluntary
		FROM commodity_supplies WHERE 1=1`
	args := []any{}
	if icID != "" {
		q += " AND industrial_capital_id = ?"
		args = append(args, string(icID))
	}
	q += " ORDER BY held_since ASC"
	rows, err := m.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []costs.CommoditySupply{}
	for rows.Next() {
		var s costs.CommoditySupply
		var sid, icid, comID string
		var pence int64
		var isVol bool
		var heldSince time.Time
		if err := rows.Scan(&sid, &icid, &comID, &pence, &heldSince, &isVol); err != nil {
			return nil, err
		}
		s.ID = costs.CommoditySupplyID(sid)
		s.IndustrialCapitalID = costs.IndustrialCapitalID(icid)
		s.CommodityID = costs.CommodityID(comID)
		s.Pence = costs.Pence(pence)
		s.HeldSince = heldSince
		s.IsVoluntary = isVol
		out = append(out, s)
	}
	return out, rows.Err()
}

// --- StorageCost ------------------------------------------------------------

func (m *MySQL) CreateStorageCost(ctx context.Context, sc costs.StorageCost) (costs.StorageCost, error) {
	if sc.ID == "" {
		sc.ID = costs.NewStorageCostID()
	}
	const q = `INSERT INTO storage_costs (id, commodity_supply_id, building_pence, labour_pence, preservation_labour_pence)
		VALUES (?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(sc.ID), string(sc.CommoditySupplyID),
		int64(sc.BuildingPence), int64(sc.LabourPence), int64(sc.PreservationLabourPence))
	if err != nil {
		return costs.StorageCost{}, err
	}
	return sc, nil
}

func (m *MySQL) ListStorageCosts(ctx context.Context, supplyID costs.CommoditySupplyID) ([]costs.StorageCost, error) {
	const q = `SELECT id, commodity_supply_id, building_pence, labour_pence, preservation_labour_pence
		FROM storage_costs WHERE commodity_supply_id = ?`
	rows, err := m.db.QueryContext(ctx, q, string(supplyID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []costs.StorageCost{}
	for rows.Next() {
		var sc costs.StorageCost
		var id, csID string
		var building, labour, preservation int64
		if err := rows.Scan(&id, &csID, &building, &labour, &preservation); err != nil {
			return nil, err
		}
		sc.ID = costs.StorageCostID(id)
		sc.CommoditySupplyID = costs.CommoditySupplyID(csID)
		sc.BuildingPence = costs.Pence(building)
		sc.LabourPence = costs.Pence(labour)
		sc.PreservationLabourPence = costs.Pence(preservation)
		out = append(out, sc)
	}
	return out, rows.Err()
}

// --- TransportTariff --------------------------------------------------------

func (m *MySQL) SetTransportTariff(ctx context.Context, t costs.TransportTariff) (costs.TransportTariff, error) {
	if t.ID == "" {
		t.ID = costs.NewTransportTariffID()
	}
	const q = `INSERT INTO transport_tariffs
		(id, commodity_id, base_pence_per_ton_mile, fragility_multiplier_basis_points, breakage_risk_multiplier_basis_points)
		VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE
		base_pence_per_ton_mile = VALUES(base_pence_per_ton_mile),
		fragility_multiplier_basis_points = VALUES(fragility_multiplier_basis_points),
		breakage_risk_multiplier_basis_points = VALUES(breakage_risk_multiplier_basis_points)`
	_, err := m.db.ExecContext(ctx, q,
		string(t.ID), string(t.CommodityID),
		int64(t.BasePencePerTonMile),
		t.FragilityMultiplierBasisPoints,
		t.BreakageRiskMultiplierBasisPoints)
	if err != nil {
		return costs.TransportTariff{}, err
	}
	return t, nil
}

func (m *MySQL) GetTransportTariff(ctx context.Context, commodityID costs.CommodityID) (costs.TransportTariff, error) {
	const q = `SELECT id, commodity_id, base_pence_per_ton_mile, fragility_multiplier_basis_points, breakage_risk_multiplier_basis_points
		FROM transport_tariffs WHERE commodity_id = ?`
	row := m.db.QueryRowContext(ctx, q, string(commodityID))
	var t costs.TransportTariff
	var id, cid string
	var base, fragility, breakage int64
	if err := row.Scan(&id, &cid, &base, &fragility, &breakage); err != nil {
		if isNotFound(err) {
			return costs.TransportTariff{}, ErrNotFound
		}
		return costs.TransportTariff{}, err
	}
	t.ID = costs.TransportTariffID(id)
	t.CommodityID = costs.CommodityID(cid)
	t.BasePencePerTonMile = costs.Pence(base)
	t.FragilityMultiplierBasisPoints = fragility
	t.BreakageRiskMultiplierBasisPoints = breakage
	return t, nil
}

// --- TransportLeg -----------------------------------------------------------

func (m *MySQL) CreateTransportLeg(ctx context.Context, leg costs.TransportLeg) (costs.TransportLeg, error) {
	if leg.ID == "" {
		leg.ID = costs.NewTransportLegID()
	}
	const q = `INSERT INTO transport_legs
		(id, commodity_id, quantity, origin_id, destination_id, distance_meters, weight_grams,
		 volume_cubic_millimeters, labour_cost_pence, means_of_transport_pence, surplus_pence)
		VALUES (?,?,?,?,?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(leg.ID), string(leg.CommodityID), leg.Quantity,
		leg.OriginID, leg.DestinationID,
		leg.DistanceMeters, leg.WeightGrams, leg.VolumeCubicMillimeters,
		int64(leg.LabourCostPence), int64(leg.MeansOfTransportPence), int64(leg.SurplusPence))
	if err != nil {
		return costs.TransportLeg{}, err
	}
	return leg, nil
}

func (m *MySQL) ListTransportLegs(ctx context.Context, commodityID costs.CommodityID) ([]costs.TransportLeg, error) {
	q := `SELECT id, commodity_id, quantity, origin_id, destination_id, distance_meters, weight_grams,
		volume_cubic_millimeters, labour_cost_pence, means_of_transport_pence, surplus_pence
		FROM transport_legs WHERE 1=1`
	args := []any{}
	if commodityID != "" {
		q += " AND commodity_id = ?"
		args = append(args, string(commodityID))
	}
	rows, err := m.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []costs.TransportLeg{}
	for rows.Next() {
		var leg costs.TransportLeg
		var id, cid, origin, dest string
		var qty, dist, weight, vol, labour, means, surplus int64
		if err := rows.Scan(&id, &cid, &qty, &origin, &dest, &dist, &weight, &vol, &labour, &means, &surplus); err != nil {
			return nil, err
		}
		leg.ID = costs.TransportLegID(id)
		leg.CommodityID = costs.CommodityID(cid)
		leg.Quantity = qty
		leg.OriginID = origin
		leg.DestinationID = dest
		leg.DistanceMeters = dist
		leg.WeightGrams = weight
		leg.VolumeCubicMillimeters = vol
		leg.LabourCostPence = costs.Pence(labour)
		leg.MeansOfTransportPence = costs.Pence(means)
		leg.SurplusPence = costs.Pence(surplus)
		out = append(out, leg)
	}
	return out, rows.Err()
}

// --- Aggregate queries ------------------------------------------------------

func (m *MySQL) AggregateForCapital(ctx context.Context, icID costs.IndustrialCapitalID, period string) (costs.AggregateResult, error) {
	items, err := m.ListCirculationCosts(ctx, icID, "", "")
	if err != nil {
		return costs.AggregateResult{}, err
	}
	res := costs.AggregateResult{
		IndustrialCapitalID: icID,
		Period:              period,
		Items:               items,
	}
	for _, c := range items {
		if c.Nature == costs.NatureValueAdding {
			res.TotalTransportPence += c.Pence
		} else {
			res.TotalFauxFraisPence += c.Pence
		}
	}
	return res, nil
}

func (m *MySQL) SystemFauxFrais(ctx context.Context, period string) (costs.SystemFauxFraisResult, error) {
	items, err := m.ListMoneyAsFauxFrais(ctx)
	if err != nil {
		return costs.SystemFauxFraisResult{}, err
	}
	res := costs.SystemFauxFraisResult{Period: period, Items: items}
	for _, mf := range items {
		res.TotalReservePence += mf.Pence
		res.TotalAnnualReplacementPence += mf.AnnualReplacementPence
	}
	return res, nil
}

// isNotFound checks for sql.ErrNoRows.
func isNotFound(err error) bool {
	return err == sql.ErrNoRows
}
