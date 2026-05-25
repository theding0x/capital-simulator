package store

import (
	"context"
	"database/sql"
	"time"

	comp "github.com/theding0x/capital-simulator/services/simulation-engine/internal/composition"
)

// CreateComponent inserts a new capital_component row.
func (m *MySQL) CreateComponent(ctx context.Context, c comp.CapitalComponent) (comp.CapitalComponent, error) {
	if err := c.Validate(); err != nil {
		return comp.CapitalComponent{}, err
	}
	if c.ID == "" {
		c.ID = comp.NewCapitalComponentID()
	}
	if c.EnteredAt.IsZero() {
		c.EnteredAt = m.now().UTC()
	}
	const q = `INSERT INTO capital_components
		(id, industrial_capital_id, kind, role, pence, entered_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID), c.IndustrialCapitalID, string(c.Kind), string(c.Role),
		int64(c.Pence), c.EnteredAt,
	)
	if err != nil {
		if isDuplicateKey(err) {
			return comp.CapitalComponent{}, ErrAlreadyExists
		}
		return comp.CapitalComponent{}, err
	}
	return c, nil
}

// GetComponent retrieves a single capital_component by ID.
func (m *MySQL) GetComponent(ctx context.Context, id comp.CapitalComponentID) (comp.CapitalComponent, error) {
	const q = `SELECT id, industrial_capital_id, kind, role, pence, entered_at
		FROM capital_components WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanComponent(row.Scan)
}

// ListComponents returns capital_components filtered by optional industrialCapitalID, kind, and role.
func (m *MySQL) ListComponents(ctx context.Context, industrialCapitalID string, kind comp.CapitalKind, role comp.RoleInProcess) ([]comp.CapitalComponent, error) {
	q := `SELECT id, industrial_capital_id, kind, role, pence, entered_at
		FROM capital_components WHERE 1=1`
	var args []any
	if industrialCapitalID != "" {
		q += " AND industrial_capital_id = ?"
		args = append(args, industrialCapitalID)
	}
	if kind != "" {
		q += " AND kind = ?"
		args = append(args, string(kind))
	}
	if role != "" {
		q += " AND role = ?"
		args = append(args, string(role))
	}
	q += " ORDER BY entered_at ASC"

	rows, err := m.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []comp.CapitalComponent{}
	for rows.Next() {
		c, err := scanComponent(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func scanComponent(scan func(...any) error) (comp.CapitalComponent, error) {
	var c comp.CapitalComponent
	var id, icID, kind, role string
	var pence int64
	var enteredAt time.Time
	if err := scan(&id, &icID, &kind, &role, &pence, &enteredAt); err != nil {
		if err == sql.ErrNoRows {
			return comp.CapitalComponent{}, ErrNotFound
		}
		return comp.CapitalComponent{}, err
	}
	c.ID = comp.CapitalComponentID(id)
	c.IndustrialCapitalID = icID
	c.Kind = comp.CapitalKind(kind)
	c.Role = comp.RoleInProcess(role)
	c.Pence = comp.Pence(pence)
	c.EnteredAt = enteredAt.UTC()
	return c, nil
}

// RegisterFixedItem inserts a fixed_capital_item and its sinking_fund in one transaction.
func (m *MySQL) RegisterFixedItem(ctx context.Context, item comp.FixedCapitalItem) (comp.FixedCapitalItem, error) {
	if err := item.Validate(); err != nil {
		return comp.FixedCapitalItem{}, err
	}
	if item.ID == "" {
		item.ID = comp.NewFixedCapitalItemID()
	}
	if item.PurchasedAt.IsZero() {
		item.PurchasedAt = m.now().UTC()
	}
	if item.WearModel == "" {
		item.WearModel = comp.WearLinear
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return comp.FixedCapitalItem{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const qItem = `INSERT INTO fixed_capital_items
		(id, capital_component_id, description, pence_purchased,
		 service_life_minutes, wear_model, purchased_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, qItem,
		string(item.ID), string(item.CapitalComponentID), item.Description,
		int64(item.PencePurchased), item.ServiceLifeMinutes,
		string(item.WearModel), item.PurchasedAt,
	)
	if err != nil {
		if isDuplicateKey(err) {
			return comp.FixedCapitalItem{}, ErrAlreadyExists
		}
		return comp.FixedCapitalItem{}, err
	}

	sfID := comp.NewSinkingFundForItemID()
	const qSF = `INSERT INTO sinking_funds_for_items
		(id, fixed_capital_item_id, accumulated_pence, target_pence)
		VALUES (?, ?, 0, ?)`
	if _, err = tx.ExecContext(ctx, qSF, string(sfID), string(item.ID), int64(item.PencePurchased)); err != nil {
		return comp.FixedCapitalItem{}, err
	}

	return item, tx.Commit()
}

// GetFixedItem retrieves a single fixed_capital_item by ID.
func (m *MySQL) GetFixedItem(ctx context.Context, id comp.FixedCapitalItemID) (comp.FixedCapitalItem, error) {
	const q = `SELECT id, capital_component_id, description, pence_purchased,
		service_life_minutes, wear_model, purchased_at, retired_at
		FROM fixed_capital_items WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanFixedItem(row.Scan)
}

func scanFixedItem(scan func(...any) error) (comp.FixedCapitalItem, error) {
	var item comp.FixedCapitalItem
	var id, ccID, desc, wearModel string
	var pencePurchased, serviceLifeMinutes int64
	var purchasedAt time.Time
	var retiredAt sql.NullTime
	if err := scan(&id, &ccID, &desc, &pencePurchased, &serviceLifeMinutes, &wearModel, &purchasedAt, &retiredAt); err != nil {
		if err == sql.ErrNoRows {
			return comp.FixedCapitalItem{}, ErrNotFound
		}
		return comp.FixedCapitalItem{}, err
	}
	item.ID = comp.FixedCapitalItemID(id)
	item.CapitalComponentID = comp.CapitalComponentID(ccID)
	item.Description = desc
	item.PencePurchased = comp.Pence(pencePurchased)
	item.ServiceLifeMinutes = serviceLifeMinutes
	item.WearModel = comp.WearModel(wearModel)
	item.PurchasedAt = purchasedAt.UTC()
	if retiredAt.Valid {
		t := retiredAt.Time.UTC()
		item.RetiredAt = &t
	}
	return item, nil
}

// RecordSubcomponent inserts a fixed_capital_subcomponent row.
func (m *MySQL) RecordSubcomponent(ctx context.Context, sub comp.FixedCapitalSubcomponent) (comp.FixedCapitalSubcomponent, error) {
	if sub.ID == "" {
		sub.ID = comp.NewFixedCapitalSubcomponentID()
	}
	if sub.WearModel == "" {
		sub.WearModel = comp.WearLinear
	}
	const q = `INSERT INTO fixed_capital_subcomponents
		(id, fixed_capital_item_id, name, pence, service_life_minutes, wear_model, last_replaced_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(sub.ID), string(sub.FixedCapitalItemID), sub.Name,
		int64(sub.Pence), sub.ServiceLifeMinutes, string(sub.WearModel), sub.LastReplacedAt,
	)
	if err != nil {
		return comp.FixedCapitalSubcomponent{}, err
	}
	return sub, nil
}

// RecordWear inserts a wear_and_tear row and increments the sinking fund in one transaction.
func (m *MySQL) RecordWear(ctx context.Context, w comp.WearAndTear) (comp.WearAndTear, error) {
	if w.ID == "" {
		w.ID = comp.NewWearAndTearID()
	}
	if w.CalculatedAt.IsZero() {
		w.CalculatedAt = m.now().UTC()
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return comp.WearAndTear{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const qWear = `INSERT INTO wear_and_tear
		(id, fixed_capital_item_id, pence, period_covered_minutes, calculated_at, cause)
		VALUES (?, ?, ?, ?, ?, ?)`
	if _, err = tx.ExecContext(ctx, qWear,
		string(w.ID), string(w.FixedCapitalItemID), int64(w.Pence),
		w.PeriodCoveredMinutes, w.CalculatedAt, string(w.Cause),
	); err != nil {
		return comp.WearAndTear{}, err
	}

	const qSF = `UPDATE sinking_funds_for_items
		SET accumulated_pence = accumulated_pence + ?
		WHERE fixed_capital_item_id = ?`
	if _, err = tx.ExecContext(ctx, qSF, int64(w.Pence), string(w.FixedCapitalItemID)); err != nil {
		return comp.WearAndTear{}, err
	}

	return w, tx.Commit()
}

// GetSinkingFund retrieves the sinking fund for a given fixed_capital_item_id.
func (m *MySQL) GetSinkingFund(ctx context.Context, itemID comp.FixedCapitalItemID) (comp.SinkingFundForItem, error) {
	const q = `SELECT id, fixed_capital_item_id, accumulated_pence, target_pence
		FROM sinking_funds_for_items WHERE fixed_capital_item_id = ?`
	row := m.db.QueryRowContext(ctx, q, string(itemID))
	var sf comp.SinkingFundForItem
	var id, ficID string
	var accumulated, target int64
	if err := row.Scan(&id, &ficID, &accumulated, &target); err != nil {
		if err == sql.ErrNoRows {
			return comp.SinkingFundForItem{}, ErrNotFound
		}
		return comp.SinkingFundForItem{}, err
	}
	sf.ID = comp.SinkingFundForItemID(id)
	sf.FixedCapitalItemID = comp.FixedCapitalItemID(ficID)
	sf.AccumulatedPence = comp.Pence(accumulated)
	sf.TargetPence = comp.Pence(target)
	return sf, nil
}

// RecordRepair inserts a repair row; for large repairs also extends service_life_minutes.
func (m *MySQL) RecordRepair(ctx context.Context, r comp.Repair) (comp.Repair, error) {
	if err := r.Validate(); err != nil {
		return comp.Repair{}, err
	}
	if r.ID == "" {
		r.ID = comp.NewRepairID()
	}
	if r.OccurredAt.IsZero() {
		r.OccurredAt = m.now().UTC()
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return comp.Repair{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const qR = `INSERT INTO repairs
		(id, fixed_capital_item_id, kind, pence, occurred_at, service_life_extension_minutes)
		VALUES (?, ?, ?, ?, ?, ?)`
	if _, err = tx.ExecContext(ctx, qR,
		string(r.ID), string(r.FixedCapitalItemID), string(r.Kind),
		int64(r.Pence), r.OccurredAt, r.ServiceLifeExtensionMinutes,
	); err != nil {
		return comp.Repair{}, err
	}

	if r.Kind == comp.RepairLarge {
		const qExt = `UPDATE fixed_capital_items
			SET service_life_minutes = service_life_minutes + ?
			WHERE id = ?`
		if _, err = tx.ExecContext(ctx, qExt, r.ServiceLifeExtensionMinutes, string(r.FixedCapitalItemID)); err != nil {
			return comp.Repair{}, err
		}
	}

	return r, tx.Commit()
}

// RecordCirculatingCycle inserts a circulating_cycles row.
func (m *MySQL) RecordCirculatingCycle(ctx context.Context, c comp.CirculatingCycle) (comp.CirculatingCycle, error) {
	if c.ID == "" {
		c.ID = comp.NewCirculatingCycleID()
	}
	const q = `INSERT INTO circulating_cycles
		(id, capital_component_id, started_at, ended_at, pence)
		VALUES (?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID), string(c.CapitalComponentID),
		c.StartedAt, c.EndedAt, int64(c.Pence),
	)
	if err != nil {
		return comp.CirculatingCycle{}, err
	}
	return c, nil
}

// RecordReinvestment inserts a sinking_fund_reinvestment row.
func (m *MySQL) RecordReinvestment(ctx context.Context, r comp.SinkingFundReinvestment) (comp.SinkingFundReinvestment, error) {
	if r.ID == "" {
		r.ID = comp.NewSinkingFundReinvestmentID()
	}
	if r.OccurredAt.IsZero() {
		r.OccurredAt = m.now().UTC()
	}
	const q = `INSERT INTO sinking_fund_reinvestments
		(id, fixed_capital_item_id, pence, kind, occurred_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(r.ID), string(r.FixedCapitalItemID),
		int64(r.Pence), string(r.Kind), r.OccurredAt,
	)
	if err != nil {
		return comp.SinkingFundReinvestment{}, err
	}
	return r, nil
}

