package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"io/fs"
	"strings"
	"time"

	pkgmysql "github.com/theding0x/capital-simulator/pkg/mysql"
	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)

//go:embed migrations
var migrationsFS embed.FS

// MySQL is a MySQL-backed Store. Construct via NewMySQL.
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

func (m *MySQL) Create(ctx context.Context, c commodity.Commodity) (commodity.Commodity, error) {
	if err := c.Validate(); err != nil {
		return commodity.Commodity{}, err
	}
	if c.ID.IsZero() {
		c.ID = commodity.NewID()
	}
	now := m.now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now

	const q = `INSERT INTO commodities
		(id, name, use_value_desc, use_value_unit, cl_kind, cl_description, snlt_per_unit, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID), c.Name,
		c.UseValue.Description, c.UseValue.Unit,
		c.ConcreteLabour.Kind, c.ConcreteLabour.Description,
		int64(c.SNLTPerUnit), c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return commodity.Commodity{}, ErrAlreadyExists
		}
		return commodity.Commodity{}, err
	}
	return c, nil
}

func (m *MySQL) Get(ctx context.Context, id commodity.ID) (commodity.Commodity, error) {
	const q = `SELECT id, name, use_value_desc, use_value_unit, cl_kind, cl_description,
		snlt_per_unit, created_at, updated_at
		FROM commodities WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanCommodity(row)
}

func (m *MySQL) List(ctx context.Context) ([]commodity.Commodity, error) {
	const q = `SELECT id, name, use_value_desc, use_value_unit, cl_kind, cl_description,
		snlt_per_unit, created_at, updated_at
		FROM commodities ORDER BY name ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []commodity.Commodity
	for rows.Next() {
		c, err := scanCommodityRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (m *MySQL) Update(ctx context.Context, id commodity.ID, u Update) (commodity.Commodity, error) {
	if u.IsEmpty() {
		return m.Get(ctx, id)
	}
	cur, err := m.Get(ctx, id)
	if err != nil {
		return commodity.Commodity{}, err
	}
	next := u.Apply(cur)
	if err := next.Validate(); err != nil {
		return commodity.Commodity{}, err
	}
	next.UpdatedAt = m.now().UTC()

	const q = `UPDATE commodities SET
		name = ?, use_value_desc = ?, use_value_unit = ?,
		cl_kind = ?, cl_description = ?, snlt_per_unit = ?, updated_at = ?
		WHERE id = ?`
	res, err := m.db.ExecContext(ctx, q,
		next.Name,
		next.UseValue.Description, next.UseValue.Unit,
		next.ConcreteLabour.Kind, next.ConcreteLabour.Description,
		int64(next.SNLTPerUnit), next.UpdatedAt,
		string(id),
	)
	if err != nil {
		if isDuplicate(err) {
			return commodity.Commodity{}, ErrAlreadyExists
		}
		return commodity.Commodity{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return commodity.Commodity{}, ErrNotFound
	}
	return next, nil
}

func (m *MySQL) Delete(ctx context.Context, id commodity.ID) error {
	res, err := m.db.ExecContext(ctx, `DELETE FROM commodities WHERE id = ?`, string(id))
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// scanCommodity scans a *sql.Row (single row) into a Commodity.
func scanCommodity(row *sql.Row) (commodity.Commodity, error) {
	var c commodity.Commodity
	var id string
	err := row.Scan(
		&id,
		&c.Name,
		&c.UseValue.Description,
		&c.UseValue.Unit,
		&c.ConcreteLabour.Kind,
		&c.ConcreteLabour.Description,
		&c.SNLTPerUnit,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return commodity.Commodity{}, ErrNotFound
	}
	if err != nil {
		return commodity.Commodity{}, err
	}
	c.ID = commodity.ID(id)
	return c, nil
}

// scanCommodityRow scans a *sql.Rows (cursor) into a Commodity.
func scanCommodityRow(rows *sql.Rows) (commodity.Commodity, error) {
	var c commodity.Commodity
	var id string
	err := rows.Scan(
		&id,
		&c.Name,
		&c.UseValue.Description,
		&c.UseValue.Unit,
		&c.ConcreteLabour.Kind,
		&c.ConcreteLabour.Description,
		&c.SNLTPerUnit,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return commodity.Commodity{}, err
	}
	c.ID = commodity.ID(id)
	return c, nil
}

// isDuplicate reports whether err is a MySQL duplicate-key error (1062).
func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "1062") || strings.Contains(s, "Duplicate entry")
}

func (m *MySQL) CreateProductionAccount(ctx context.Context, a commodity.ProductionAccount) (commodity.ProductionAccount, error) {
	if err := a.Validate(); err != nil {
		return commodity.ProductionAccount{}, err
	}
	if a.ID.IsZero() {
		a.ID = commodity.NewProductionAccountID()
	}
	a.CreatedAt = m.now().UTC()
	const q = `INSERT INTO production_accounts (id, constant, variable, surplus, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID),
		int64(a.Constant),
		int64(a.Variable),
		int64(a.Surplus),
		a.CreatedAt,
	)
	if err != nil {
		return commodity.ProductionAccount{}, err
	}
	return a, nil
}

func (m *MySQL) GetProductionAccount(ctx context.Context, id commodity.ProductionAccountID) (commodity.ProductionAccount, error) {
	const q = `SELECT id, constant, variable, surplus, created_at FROM production_accounts WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanProductionAccount(row)
}

func (m *MySQL) ListProductionAccounts(ctx context.Context) ([]commodity.ProductionAccount, error) {
	const q = `SELECT id, constant, variable, surplus, created_at FROM production_accounts ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []commodity.ProductionAccount
	for rows.Next() {
		a, err := scanProductionAccountRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func scanProductionAccount(row *sql.Row) (commodity.ProductionAccount, error) {
	var a commodity.ProductionAccount
	var id string
	var c, v, s int64
	err := row.Scan(&id, &c, &v, &s, &a.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return commodity.ProductionAccount{}, ErrNotFound
	}
	if err != nil {
		return commodity.ProductionAccount{}, err
	}
	a.ID = commodity.ProductionAccountID(id)
	a.Constant = commodity.LabourMinutes(c)
	a.Variable = commodity.LabourMinutes(v)
	a.Surplus = commodity.SurplusValue(s)
	return a, nil
}

func scanProductionAccountRow(rows *sql.Rows) (commodity.ProductionAccount, error) {
	var a commodity.ProductionAccount
	var id string
	var c, v, s int64
	if err := rows.Scan(&id, &c, &v, &s, &a.CreatedAt); err != nil {
		return commodity.ProductionAccount{}, err
	}
	a.ID = commodity.ProductionAccountID(id)
	a.Constant = commodity.LabourMinutes(c)
	a.Variable = commodity.LabourMinutes(v)
	a.Surplus = commodity.SurplusValue(s)
	return a, nil
}
