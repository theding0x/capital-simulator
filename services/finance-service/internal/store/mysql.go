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
	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
)

//go:embed migrations
var migrationsFS embed.FS

// MySQL is a MySQL-backed Store for finance-service. Construct via NewMySQL.
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

// CreateCostPrice persists cp, assigning an ID and timestamp when absent.
func (m *MySQL) CreateCostPrice(ctx context.Context, cp profit.CostPrice) (profit.CostPrice, error) {
	if cp.ID.IsZero() {
		cp.ID = profit.NewCostPriceID()
	}
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO cost_prices
		(id, constant, variable, fixed_wear_and_tear, fixed_advanced, k, fixed_component, circulating_component, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(cp.ID),
		int64(cp.Outlay.Constant), int64(cp.Outlay.Variable),
		int64(cp.Outlay.FixedWearAndTear), int64(cp.Outlay.FixedAdvanced),
		int64(cp.K), int64(cp.FixedComponent), int64(cp.CirculatingComponent),
		cp.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return profit.CostPrice{}, ErrAlreadyExists
		}
		return profit.CostPrice{}, err
	}
	return cp, nil
}

// GetCostPrice returns the cost-price with id, or ErrNotFound.
func (m *MySQL) GetCostPrice(ctx context.Context, id profit.CostPriceID) (profit.CostPrice, error) {
	const q = `SELECT id, constant, variable, fixed_wear_and_tear, fixed_advanced,
		k, fixed_component, circulating_component, created_at
		FROM cost_prices WHERE id = ?`
	return scanCostPrice(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListCostPrices returns all stored cost-prices, newest first.
func (m *MySQL) ListCostPrices(ctx context.Context) ([]profit.CostPrice, error) {
	const q = `SELECT id, constant, variable, fixed_wear_and_tear, fixed_advanced,
		k, fixed_component, circulating_component, created_at
		FROM cost_prices ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []profit.CostPrice
	for rows.Next() {
		cp, err := scanCostPrice(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, cp)
	}
	return out, rows.Err()
}

// rowScanner is satisfied by both *sql.Row and *sql.Rows.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanCostPrice(s rowScanner) (profit.CostPrice, error) {
	var (
		id                                           string
		constant, variable, fixedWear, fixedAdvanced int64
		k, fixedComponent, circulatingComponent      int64
		createdAt                                    time.Time
	)
	err := s.Scan(&id, &constant, &variable, &fixedWear, &fixedAdvanced,
		&k, &fixedComponent, &circulatingComponent, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return profit.CostPrice{}, ErrNotFound
		}
		return profit.CostPrice{}, err
	}
	return profit.CostPrice{
		ID: profit.CostPriceID(id),
		Outlay: profit.CapitalOutlay{
			Constant:         profit.ConstantCapital(constant),
			Variable:         profit.VariableCapital(variable),
			FixedWearAndTear: profit.ConstantCapital(fixedWear),
			FixedAdvanced:    profit.ConstantCapital(fixedAdvanced),
		},
		K:                    profit.LabourMinutes(k),
		FixedComponent:       profit.LabourMinutes(fixedComponent),
		CirculatingComponent: profit.LabourMinutes(circulatingComponent),
		CreatedAt:            createdAt,
	}, nil
}

// isDuplicate reports whether err is a MySQL duplicate-key error (1062).
func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "1062") || strings.Contains(s, "Duplicate entry")
}
