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

// isDuplicate reports whether err is a MySQL duplicate-key error (1062).
func isDuplicate(err error) bool {
	var me *mysql.MySQLError
	return errors.As(err, &me) && me.Number == 1062
}
