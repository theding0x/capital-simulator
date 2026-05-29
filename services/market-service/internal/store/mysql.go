package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"time"

	mysql "github.com/go-sql-driver/mysql"
	pkgmysql "github.com/theding0x/capital-simulator/pkg/mysql"
	"github.com/theding0x/capital-simulator/services/market-service/internal/circulation"
	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

//go:embed migrations
var migrationsFS embed.FS

const (
	keyUniversalEquivalent = "universal_equivalent"
	keyMoneyCommodity      = "money_commodity"
)

// MySQL is a MySQL-backed Store for market-service. Construct via NewMySQL.
type MySQL struct {
	db  *sql.DB
	now func() time.Time
}

// NewMySQL returns a Store backed by db and runs any pending migrations.
func NewMySQL(ctx context.Context, db *sql.DB) (*MySQL, error) {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		return nil, err
	}
	if err := pkgmysql.Migrate(ctx, db, sub); err != nil {
		return nil, err
	}
	return &MySQL{db: db, now: time.Now}, nil
}

// --- Owner ------------------------------------------------------------------

func (m *MySQL) CreateOwner(ctx context.Context, o market.Owner) (market.Owner, error) {
	if err := o.Validate(); err != nil {
		return market.Owner{}, err
	}
	if o.ID.IsZero() {
		o.ID = market.NewOwnerID()
	}
	now := m.now().UTC()
	o.CreatedAt = now
	o.UpdatedAt = now
	const q = `INSERT INTO owners (id, name, created_at, updated_at) VALUES (?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q, string(o.ID), o.Name, o.CreatedAt, o.UpdatedAt)
	if err != nil {
		if isDuplicate(err) {
			return market.Owner{}, ErrAlreadyExists
		}
		return market.Owner{}, err
	}
	return o, nil
}

func (m *MySQL) GetOwner(ctx context.Context, id market.OwnerID) (market.Owner, error) {
	const q = `SELECT id, name, created_at, updated_at FROM owners WHERE id = ?`
	var o market.Owner
	var rawID string
	err := m.db.QueryRowContext(ctx, q, string(id)).Scan(&rawID, &o.Name, &o.CreatedAt, &o.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return market.Owner{}, ErrNotFound
	}
	if err != nil {
		return market.Owner{}, err
	}
	o.ID = market.OwnerID(rawID)
	return o, nil
}

func (m *MySQL) ListOwners(ctx context.Context) ([]market.Owner, error) {
	const q = `SELECT id, name, created_at, updated_at FROM owners ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []market.Owner
	for rows.Next() {
		var o market.Owner
		var rawID string
		if err := rows.Scan(&rawID, &o.Name, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		o.ID = market.OwnerID(rawID)
		out = append(out, o)
	}
	return out, rows.Err()
}

// --- Offer ------------------------------------------------------------------

func (m *MySQL) CreateOffer(ctx context.Context, o market.Offer) (market.Offer, error) {
	if err := o.Validate(); err != nil {
		return market.Offer{}, err
	}
	if o.ID.IsZero() {
		o.ID = market.NewOfferID()
	}
	o.CreatedAt = m.now().UTC()
	const q = `INSERT INTO offers
		(id, owner_id, commodity_id, quantity, seeks_kind, seeks_commodity_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(o.ID), string(o.OwnerID), string(o.CommodityID),
		o.Quantity, o.SeeksKind, string(o.SeeksCommodityID), o.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return market.Offer{}, ErrAlreadyExists
		}
		return market.Offer{}, err
	}
	return o, nil
}

func (m *MySQL) GetOffer(ctx context.Context, id market.OfferID) (market.Offer, error) {
	const q = `SELECT id, owner_id, commodity_id, quantity, seeks_kind, seeks_commodity_id, created_at
		FROM offers WHERE id = ?`
	return scanOffer(m.db.QueryRowContext(ctx, q, string(id)))
}

func (m *MySQL) ListOffers(ctx context.Context) ([]market.Offer, error) {
	const q = `SELECT id, owner_id, commodity_id, quantity, seeks_kind, seeks_commodity_id, created_at
		FROM offers ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []market.Offer
	for rows.Next() {
		var o market.Offer
		var id, ownerID, commodityID, seeksCommodityID string
		if err := rows.Scan(&id, &ownerID, &commodityID, &o.Quantity,
			&o.SeeksKind, &seeksCommodityID, &o.CreatedAt); err != nil {
			return nil, err
		}
		o.ID = market.OfferID(id)
		o.OwnerID = market.OwnerID(ownerID)
		o.CommodityID = market.CommodityID(commodityID)
		o.SeeksCommodityID = market.CommodityID(seeksCommodityID)
		out = append(out, o)
	}
	return out, rows.Err()
}

func (m *MySQL) DeleteOffer(ctx context.Context, id market.OfferID) error {
	res, err := m.db.ExecContext(ctx, `DELETE FROM offers WHERE id = ?`, string(id))
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete offer: rows affected: %w", err)
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func scanOffer(row *sql.Row) (market.Offer, error) {
	var o market.Offer
	var id, ownerID, commodityID, seeksCommodityID string
	err := row.Scan(&id, &ownerID, &commodityID, &o.Quantity,
		&o.SeeksKind, &seeksCommodityID, &o.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return market.Offer{}, ErrNotFound
	}
	if err != nil {
		return market.Offer{}, err
	}
	o.ID = market.OfferID(id)
	o.OwnerID = market.OwnerID(ownerID)
	o.CommodityID = market.CommodityID(commodityID)
	o.SeeksCommodityID = market.CommodityID(seeksCommodityID)
	return o, nil
}

// --- Exchange ---------------------------------------------------------------

func (m *MySQL) CreateExchange(ctx context.Context, e market.Exchange) (market.Exchange, error) {
	if err := e.Validate(); err != nil {
		return market.Exchange{}, err
	}
	if e.ID.IsZero() {
		e.ID = market.NewExchangeID()
	}
	e.CreatedAt = m.now().UTC()
	const q = `INSERT INTO exchanges
		(id, giver_id, receiver_id, giver_commodity_id, giver_qty, receiver_commodity_id, receiver_qty, realised_value, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(e.ID),
		string(e.GiverID), string(e.ReceiverID),
		string(e.GiverCommodityID), e.GiverQty,
		string(e.ReceiverCommodityID), e.ReceiverQty,
		int64(e.RealisedValue), e.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return market.Exchange{}, ErrAlreadyExists
		}
		return market.Exchange{}, err
	}
	return e, nil
}

func (m *MySQL) GetExchange(ctx context.Context, id market.ExchangeID) (market.Exchange, error) {
	const q = `SELECT id, giver_id, receiver_id, giver_commodity_id, giver_qty,
		receiver_commodity_id, receiver_qty, realised_value, created_at
		FROM exchanges WHERE id = ?`
	return scanExchange(m.db.QueryRowContext(ctx, q, string(id)))
}

func (m *MySQL) ListExchanges(ctx context.Context) ([]market.Exchange, error) {
	const q = `SELECT id, giver_id, receiver_id, giver_commodity_id, giver_qty,
		receiver_commodity_id, receiver_qty, realised_value, created_at
		FROM exchanges ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []market.Exchange
	for rows.Next() {
		e, err := scanExchangeRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func scanExchange(row *sql.Row) (market.Exchange, error) {
	var e market.Exchange
	var id, giverID, receiverID, giverCID, receiverCID string
	var rv int64
	err := row.Scan(&id, &giverID, &receiverID, &giverCID, &e.GiverQty,
		&receiverCID, &e.ReceiverQty, &rv, &e.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return market.Exchange{}, ErrNotFound
	}
	if err != nil {
		return market.Exchange{}, err
	}
	e.ID = market.ExchangeID(id)
	e.GiverID = market.OwnerID(giverID)
	e.ReceiverID = market.OwnerID(receiverID)
	e.GiverCommodityID = market.CommodityID(giverCID)
	e.ReceiverCommodityID = market.CommodityID(receiverCID)
	e.RealisedValue = market.RealisedValue(rv)
	return e, nil
}

func scanExchangeRow(rows *sql.Rows) (market.Exchange, error) {
	var e market.Exchange
	var id, giverID, receiverID, giverCID, receiverCID string
	var rv int64
	if err := rows.Scan(&id, &giverID, &receiverID, &giverCID, &e.GiverQty,
		&receiverCID, &e.ReceiverQty, &rv, &e.CreatedAt); err != nil {
		return market.Exchange{}, err
	}
	e.ID = market.ExchangeID(id)
	e.GiverID = market.OwnerID(giverID)
	e.ReceiverID = market.OwnerID(receiverID)
	e.GiverCommodityID = market.CommodityID(giverCID)
	e.ReceiverCommodityID = market.CommodityID(receiverCID)
	e.RealisedValue = market.RealisedValue(rv)
	return e, nil
}

// --- UniversalEquivalent (singleton in market_config) -----------------------

func (m *MySQL) SetUniversalEquivalent(ctx context.Context, ue market.UniversalEquivalent) (market.UniversalEquivalent, error) {
	existing, err := m.GetUniversalEquivalent(ctx)
	if err == nil && existing.CommodityID == ue.CommodityID {
		return existing, nil
	}
	if ue.SetAt.IsZero() {
		ue.SetAt = m.now().UTC()
	}
	const q = `INSERT INTO market_config (key_name, commodity_id, ts) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE commodity_id = VALUES(commodity_id), ts = VALUES(ts)`
	_, err = m.db.ExecContext(ctx, q, keyUniversalEquivalent, string(ue.CommodityID), ue.SetAt)
	if err != nil {
		return market.UniversalEquivalent{}, err
	}
	return ue, nil
}

func (m *MySQL) GetUniversalEquivalent(ctx context.Context) (market.UniversalEquivalent, error) {
	const q = `SELECT commodity_id, ts FROM market_config WHERE key_name = ?`
	var commodityID string
	var ue market.UniversalEquivalent
	err := m.db.QueryRowContext(ctx, q, keyUniversalEquivalent).Scan(&commodityID, &ue.SetAt)
	if errors.Is(err, sql.ErrNoRows) {
		return market.UniversalEquivalent{}, ErrNotFound
	}
	if err != nil {
		return market.UniversalEquivalent{}, err
	}
	ue.CommodityID = market.CommodityID(commodityID)
	return ue, nil
}

// --- MoneyCommodity (singleton in market_config) ----------------------------

func (m *MySQL) SetMoneyCommodity(ctx context.Context, mc market.MoneyCommodity) (market.MoneyCommodity, error) {
	if mc.CreatedAt.IsZero() {
		mc.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO market_config (key_name, commodity_id, ts) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE commodity_id = VALUES(commodity_id), ts = VALUES(ts)`
	_, err := m.db.ExecContext(ctx, q, keyMoneyCommodity, string(mc.CommodityID), mc.CreatedAt)
	if err != nil {
		return market.MoneyCommodity{}, err
	}
	return mc, nil
}

func (m *MySQL) GetMoneyCommodity(ctx context.Context) (market.MoneyCommodity, error) {
	const q = `SELECT commodity_id, ts FROM market_config WHERE key_name = ?`
	var commodityID string
	var mc market.MoneyCommodity
	err := m.db.QueryRowContext(ctx, q, keyMoneyCommodity).Scan(&commodityID, &mc.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return market.MoneyCommodity{}, ErrNotFound
	}
	if err != nil {
		return market.MoneyCommodity{}, err
	}
	mc.CommodityID = market.CommodityID(commodityID)
	return mc, nil
}

// --- Price ------------------------------------------------------------------

func (m *MySQL) SetPrice(ctx context.Context, p market.Price) (market.Price, error) {
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = m.now().UTC()
	}
	const q = `INSERT INTO prices (commodity_id, money_commodity_id, amount, updated_at) VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE money_commodity_id = VALUES(money_commodity_id),
		amount = VALUES(amount), updated_at = VALUES(updated_at)`
	_, err := m.db.ExecContext(ctx, q,
		string(p.CommodityID), string(p.MoneyCommodityID), int64(p.Amount), p.UpdatedAt)
	if err != nil {
		return market.Price{}, err
	}
	return p, nil
}

func (m *MySQL) GetPrice(ctx context.Context, commodityID market.CommodityID) (market.Price, error) {
	const q = `SELECT commodity_id, money_commodity_id, amount, updated_at FROM prices WHERE commodity_id = ?`
	var p market.Price
	var cid, mcid string
	var amount int64
	err := m.db.QueryRowContext(ctx, q, string(commodityID)).Scan(&cid, &mcid, &amount, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return market.Price{}, ErrNotFound
	}
	if err != nil {
		return market.Price{}, err
	}
	p.CommodityID = market.CommodityID(cid)
	p.MoneyCommodityID = market.CommodityID(mcid)
	p.Amount = market.PriceAmount(amount)
	return p, nil
}

func (m *MySQL) ListPrices(ctx context.Context) ([]market.Price, error) {
	const q = `SELECT commodity_id, money_commodity_id, amount, updated_at FROM prices ORDER BY commodity_id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []market.Price
	for rows.Next() {
		var p market.Price
		var cid, mcid string
		var amount int64
		if err := rows.Scan(&cid, &mcid, &amount, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.CommodityID = market.CommodityID(cid)
		p.MoneyCommodityID = market.CommodityID(mcid)
		p.Amount = market.PriceAmount(amount)
		out = append(out, p)
	}
	return out, rows.Err()
}

// --- Vol. II Ch. 5 — TurnoverTimeStore -------------------------------------

func (m *MySQL) CreateTurnoverTime(ctx context.Context, t circulation.TurnoverTime) (circulation.TurnoverTime, error) {
	if err := t.Validate(); err != nil {
		return circulation.TurnoverTime{}, err
	}
	if t.ID.IsZero() {
		t.ID = circulation.NewTurnoverTimeID()
	}
	now := m.now()
	t.CreatedAt = now
	t.UpdatedAt = now
	const q = `INSERT INTO turnover_times
		(id, industrial_capital_id, labour_time_nanos, labour_interruption_nanos, latent_nanos,
		 natural_process_nanos, selling_time_nanos, buying_time_nanos, circulation_open, created_at, updated_at)
		VALUES (?,?,?,?,?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(t.ID), string(t.IndustrialCapitalID),
		t.Production.LabourTimeNanos, t.Production.LabourInterruptionNanos,
		t.Production.LatentNanos, t.Production.NaturalProcessNanos,
		t.Circulation.SellingTimeNanos, t.Circulation.BuyingTimeNanos,
		t.CirculationOpen, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		if isDuplicate(err) {
			return circulation.TurnoverTime{}, ErrAlreadyExists
		}
		return circulation.TurnoverTime{}, err
	}
	return t, nil
}

func (m *MySQL) GetTurnoverTime(ctx context.Context, id circulation.TurnoverTimeID) (circulation.TurnoverTime, error) {
	const q = `SELECT id, industrial_capital_id, labour_time_nanos, labour_interruption_nanos,
		latent_nanos, natural_process_nanos, selling_time_nanos, buying_time_nanos,
		circulation_open, created_at, updated_at FROM turnover_times WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanTurnoverTime(row)
}

func scanTurnoverTime(row *sql.Row) (circulation.TurnoverTime, error) {
	var t circulation.TurnoverTime
	var tid, capID string
	err := row.Scan(&tid, &capID,
		&t.Production.LabourTimeNanos, &t.Production.LabourInterruptionNanos,
		&t.Production.LatentNanos, &t.Production.NaturalProcessNanos,
		&t.Circulation.SellingTimeNanos, &t.Circulation.BuyingTimeNanos,
		&t.CirculationOpen, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return circulation.TurnoverTime{}, ErrNotFound
	}
	if err != nil {
		return circulation.TurnoverTime{}, err
	}
	t.ID = circulation.TurnoverTimeID(tid)
	t.IndustrialCapitalID = circulation.IndustrialCapitalID(capID)
	return t, nil
}

func (m *MySQL) ListTurnoverTimes(ctx context.Context) ([]circulation.TurnoverTime, error) {
	const q = `SELECT id, industrial_capital_id, labour_time_nanos, labour_interruption_nanos,
		latent_nanos, natural_process_nanos, selling_time_nanos, buying_time_nanos,
		circulation_open, created_at, updated_at FROM turnover_times ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []circulation.TurnoverTime
	for rows.Next() {
		var t circulation.TurnoverTime
		var tid, capID string
		if err := rows.Scan(&tid, &capID,
			&t.Production.LabourTimeNanos, &t.Production.LabourInterruptionNanos,
			&t.Production.LatentNanos, &t.Production.NaturalProcessNanos,
			&t.Circulation.SellingTimeNanos, &t.Circulation.BuyingTimeNanos,
			&t.CirculationOpen, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		t.ID = circulation.TurnoverTimeID(tid)
		t.IndustrialCapitalID = circulation.IndustrialCapitalID(capID)
		out = append(out, t)
	}
	return out, rows.Err()
}

func (m *MySQL) AddLabourTime(ctx context.Context, id circulation.TurnoverTimeID, nanos int64) (circulation.TurnoverTime, error) {
	const check = `SELECT circulation_open FROM turnover_times WHERE id = ?`
	var open bool
	if err := m.db.QueryRowContext(ctx, check, string(id)).Scan(&open); errors.Is(err, sql.ErrNoRows) {
		return circulation.TurnoverTime{}, ErrNotFound
	} else if err != nil {
		return circulation.TurnoverTime{}, err
	}
	if open {
		return circulation.TurnoverTime{}, circulation.ErrConcurrentProductionAndCirculation
	}
	const q = `UPDATE turnover_times SET labour_time_nanos = labour_time_nanos + ?, updated_at = ? WHERE id = ?`
	if _, err := m.db.ExecContext(ctx, q, nanos, m.now(), string(id)); err != nil {
		return circulation.TurnoverTime{}, err
	}
	return m.GetTurnoverTime(ctx, id)
}

func (m *MySQL) AddLabourInterruption(ctx context.Context, id circulation.TurnoverTimeID, nanos int64) (circulation.TurnoverTime, error) {
	const check = `SELECT circulation_open FROM turnover_times WHERE id = ?`
	var open bool
	if err := m.db.QueryRowContext(ctx, check, string(id)).Scan(&open); errors.Is(err, sql.ErrNoRows) {
		return circulation.TurnoverTime{}, ErrNotFound
	} else if err != nil {
		return circulation.TurnoverTime{}, err
	}
	if open {
		return circulation.TurnoverTime{}, circulation.ErrConcurrentProductionAndCirculation
	}
	const q = `UPDATE turnover_times SET labour_interruption_nanos = labour_interruption_nanos + ?, updated_at = ? WHERE id = ?`
	if _, err := m.db.ExecContext(ctx, q, nanos, m.now(), string(id)); err != nil {
		return circulation.TurnoverTime{}, err
	}
	return m.GetTurnoverTime(ctx, id)
}

func (m *MySQL) RecordLatentMP(ctx context.Context, lpc circulation.LatentProductiveCapital) (circulation.LatentProductiveCapital, error) {
	if lpc.ID.IsZero() {
		lpc.ID = circulation.NewLatentProductiveCapitalID()
	}
	const q = `INSERT INTO latent_productive_capital
		(id, turnover_time_id, industrial_capital_id, pence, held_at, entered_production_at)
		VALUES (?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(lpc.ID), string(lpc.TurnoverTimeID), string(lpc.IndustrialCapitalID),
		int64(lpc.Pence), lpc.HeldAt, lpc.EnteredProductionAt)
	if err != nil {
		return circulation.LatentProductiveCapital{}, err
	}
	const upd = `UPDATE turnover_times SET latent_nanos = latent_nanos + ?, updated_at = ? WHERE id = ?`
	if _, err := m.db.ExecContext(ctx, upd, int64(time.Second), m.now(), string(lpc.TurnoverTimeID)); err != nil {
		return circulation.LatentProductiveCapital{}, err
	}
	return lpc, nil
}

func (m *MySQL) RecordNaturalProcess(ctx context.Context, nps circulation.NaturalProcessSpan) (circulation.NaturalProcessSpan, error) {
	if nps.ID.IsZero() {
		nps.ID = circulation.NewNaturalProcessSpanID()
	}
	const q = `INSERT INTO natural_process_spans
		(id, turnover_time_id, industrial_capital_id, process, duration_nanos, started_at)
		VALUES (?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(nps.ID), string(nps.TurnoverTimeID), string(nps.IndustrialCapitalID),
		string(nps.Process), nps.DurationNanos, nps.StartedAt)
	if err != nil {
		return circulation.NaturalProcessSpan{}, err
	}
	const upd = `UPDATE turnover_times SET natural_process_nanos = natural_process_nanos + ?, updated_at = ? WHERE id = ?`
	if _, err := m.db.ExecContext(ctx, upd, nps.DurationNanos, m.now(), string(nps.TurnoverTimeID)); err != nil {
		return circulation.NaturalProcessSpan{}, err
	}
	return nps, nil
}

func (m *MySQL) OpenSellingPhase(ctx context.Context, sp circulation.SellingPhase) (circulation.SellingPhase, error) {
	if sp.ID.IsZero() {
		sp.ID = circulation.NewSellingPhaseID()
	}
	if sp.Outcome == "" {
		sp.Outcome = circulation.OutcomePending
	}
	const q = `INSERT INTO selling_phases
		(id, turnover_time_id, industrial_capital_id, commodity_circuit_id, opened_at, outcome)
		VALUES (?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(sp.ID), string(sp.TurnoverTimeID), string(sp.IndustrialCapitalID),
		string(sp.CommodityCircuitID), sp.OpenedAt, string(sp.Outcome))
	if err != nil {
		if isDuplicate(err) {
			return circulation.SellingPhase{}, circulation.ErrSellingPhaseAlreadyOpen
		}
		return circulation.SellingPhase{}, err
	}
	const upd = `UPDATE turnover_times SET circulation_open = true, updated_at = ? WHERE id = ?`
	if _, err := m.db.ExecContext(ctx, upd, m.now(), string(sp.TurnoverTimeID)); err != nil {
		return circulation.SellingPhase{}, err
	}
	return sp, nil
}

func (m *MySQL) CloseSellingPhase(ctx context.Context, id circulation.TurnoverTimeID, spID circulation.SellingPhaseID, outcome circulation.SellingOutcome) (circulation.SellingPhase, error) {
	now := m.now()
	const q = `UPDATE selling_phases SET closed_at = ?, outcome = ? WHERE id = ? AND turnover_time_id = ? AND closed_at IS NULL`
	res, err := m.db.ExecContext(ctx, q, now, string(outcome), string(spID), string(id))
	if err != nil {
		return circulation.SellingPhase{}, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return circulation.SellingPhase{}, circulation.ErrNoOpenSellingPhase
	}
	const sel = `SELECT id, turnover_time_id, industrial_capital_id, commodity_circuit_id, opened_at, closed_at, outcome FROM selling_phases WHERE id = ?`
	row := m.db.QueryRowContext(ctx, sel, string(spID))
	var sp circulation.SellingPhase
	var sid, ttid, capID, ccid, oc string
	if err := row.Scan(&sid, &ttid, &capID, &ccid, &sp.OpenedAt, &sp.ClosedAt, &oc); err != nil {
		return circulation.SellingPhase{}, err
	}
	sp.ID = circulation.SellingPhaseID(sid)
	sp.TurnoverTimeID = circulation.TurnoverTimeID(ttid)
	sp.IndustrialCapitalID = circulation.IndustrialCapitalID(capID)
	sp.CommodityCircuitID = circulation.CommodityCircuitID(ccid)
	sp.Outcome = circulation.SellingOutcome(oc)

	const upd = `UPDATE turnover_times SET selling_time_nanos = selling_time_nanos + TIMESTAMPDIFF(second, ?, ?) * 1000000000, updated_at = ? WHERE id = ?`
	if _, err := m.db.ExecContext(ctx, upd, sp.OpenedAt, now, now, string(id)); err != nil {
		return circulation.SellingPhase{}, err
	}
	return sp, nil
}

func (m *MySQL) OpenBuyingPhase(ctx context.Context, bp circulation.BuyingPhase) (circulation.BuyingPhase, error) {
	if bp.ID.IsZero() {
		bp.ID = circulation.NewBuyingPhaseID()
	}
	const q = `INSERT INTO buying_phases
		(id, turnover_time_id, industrial_capital_id, money_circuit_id, opened_at, market_location)
		VALUES (?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(bp.ID), string(bp.TurnoverTimeID), string(bp.IndustrialCapitalID),
		string(bp.MoneyCircuitID), bp.OpenedAt, bp.MarketLocation)
	if err != nil {
		if isDuplicate(err) {
			return circulation.BuyingPhase{}, circulation.ErrBuyingPhaseAlreadyOpen
		}
		return circulation.BuyingPhase{}, err
	}
	const upd = `UPDATE turnover_times SET circulation_open = true, updated_at = ? WHERE id = ?`
	if _, err := m.db.ExecContext(ctx, upd, m.now(), string(bp.TurnoverTimeID)); err != nil {
		return circulation.BuyingPhase{}, err
	}
	return bp, nil
}

func (m *MySQL) CloseBuyingPhase(ctx context.Context, id circulation.TurnoverTimeID, bpID circulation.BuyingPhaseID) (circulation.BuyingPhase, error) {
	now := m.now()
	const q = `UPDATE buying_phases SET closed_at = ? WHERE id = ? AND turnover_time_id = ? AND closed_at IS NULL`
	res, err := m.db.ExecContext(ctx, q, now, string(bpID), string(id))
	if err != nil {
		return circulation.BuyingPhase{}, err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return circulation.BuyingPhase{}, circulation.ErrNoOpenBuyingPhase
	}
	const sel = `SELECT id, turnover_time_id, industrial_capital_id, money_circuit_id, opened_at, closed_at, market_location FROM buying_phases WHERE id = ?`
	row := m.db.QueryRowContext(ctx, sel, string(bpID))
	var bp circulation.BuyingPhase
	var bid, ttid, capID, mcid string
	if err := row.Scan(&bid, &ttid, &capID, &mcid, &bp.OpenedAt, &bp.ClosedAt, &bp.MarketLocation); err != nil {
		return circulation.BuyingPhase{}, err
	}
	bp.ID = circulation.BuyingPhaseID(bid)
	bp.TurnoverTimeID = circulation.TurnoverTimeID(ttid)
	bp.IndustrialCapitalID = circulation.IndustrialCapitalID(capID)
	bp.MoneyCircuitID = circulation.MoneyCircuitID(mcid)

	const upd = `UPDATE turnover_times SET buying_time_nanos = buying_time_nanos + TIMESTAMPDIFF(second, ?, ?) * 1000000000, circulation_open = false, updated_at = ? WHERE id = ?`
	if _, err := m.db.ExecContext(ctx, upd, bp.OpenedAt, now, now, string(id)); err != nil {
		return circulation.BuyingPhase{}, err
	}
	return bp, nil
}

func (m *MySQL) SetPerishability(ctx context.Context, p circulation.Perishability) (circulation.Perishability, error) {
	if p.ID.IsZero() {
		p.ID = circulation.NewPerishabilityID()
	}
	const q = `INSERT INTO perishability (id, commodity_id, window_nanos)
		VALUES (?,?,?) ON DUPLICATE KEY UPDATE window_nanos = VALUES(window_nanos)`
	if _, err := m.db.ExecContext(ctx, q, string(p.ID), string(p.CommodityID), p.WindowNanos); err != nil {
		return circulation.Perishability{}, err
	}
	return p, nil
}

func (m *MySQL) GetPerishability(ctx context.Context, commodityID circulation.CommodityID) (circulation.Perishability, error) {
	const q = `SELECT id, commodity_id, window_nanos FROM perishability WHERE commodity_id = ?`
	row := m.db.QueryRowContext(ctx, q, string(commodityID))
	var p circulation.Perishability
	var pid, cid string
	if err := row.Scan(&pid, &cid, &p.WindowNanos); errors.Is(err, sql.ErrNoRows) {
		return circulation.Perishability{}, ErrNotFound
	} else if err != nil {
		return circulation.Perishability{}, err
	}
	p.ID = circulation.PerishabilityID(pid)
	p.CommodityID = circulation.CommodityID(cid)
	return p, nil
}

func (m *MySQL) SetMarketSeparation(ctx context.Context, ms circulation.MarketSeparation) (circulation.MarketSeparation, error) {
	if ms.ID.IsZero() {
		ms.ID = circulation.NewMarketSeparationID()
	}
	const q = `INSERT INTO market_separation
		(id, industrial_capital_id, selling_market_id, buying_market_id)
		VALUES (?,?,?,?) ON DUPLICATE KEY UPDATE selling_market_id = VALUES(selling_market_id), buying_market_id = VALUES(buying_market_id)`
	if _, err := m.db.ExecContext(ctx, q, string(ms.ID), string(ms.IndustrialCapitalID), ms.SellingMarketID, ms.BuyingMarketID); err != nil {
		return circulation.MarketSeparation{}, err
	}
	return ms, nil
}

// isDuplicate reports whether err is a MySQL duplicate-key error (1062).
func isDuplicate(err error) bool {
	var me *mysql.MySQLError
	return errors.As(err, &me) && me.Number == 1062
}
