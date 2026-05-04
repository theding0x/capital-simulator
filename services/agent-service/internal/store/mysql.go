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
	"github.com/theding0x/capital-simulator/services/agent-service/internal/agent"
)

//go:embed migrations
var migrationsFS embed.FS

// MySQL implements Store and CircuitStore backed by MySQL.
type MySQL struct {
	db  *sql.DB
	now func() time.Time
}

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

func (m *MySQL) Create(ctx context.Context, a agent.Agent) (agent.Agent, error) {
	if err := a.Validate(); err != nil {
		return agent.Agent{}, err
	}
	if a.ID.IsZero() {
		a.ID = agent.NewID()
	}
	now := m.now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now
	const q = `INSERT INTO agents (id, name, class, money_balance, hoarding, labour_minutes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID), a.Name, string(a.Class), int64(a.MoneyBalance), a.Hoarding,
		a.LabourMinutes, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return agent.Agent{}, err
	}
	return a, nil
}

func (m *MySQL) Get(ctx context.Context, id agent.ID) (agent.Agent, error) {
	const q = `SELECT id, name, class, money_balance, hoarding, labour_minutes, created_at, updated_at
		FROM agents WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanAgent(row)
}

func (m *MySQL) List(ctx context.Context) ([]agent.Agent, error) {
	const q = `SELECT id, name, class, money_balance, hoarding, labour_minutes, created_at, updated_at
		FROM agents ORDER BY name ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAgents(rows)
}

func (m *MySQL) ListByClass(ctx context.Context, class agent.Class) ([]agent.Agent, error) {
	const q = `SELECT id, name, class, money_balance, hoarding, labour_minutes, created_at, updated_at
		FROM agents WHERE class = ? ORDER BY name ASC`
	rows, err := m.db.QueryContext(ctx, q, string(class))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAgents(rows)
}

func (m *MySQL) Update(ctx context.Context, id agent.ID, u Update) (agent.Agent, error) {
	if u.IsEmpty() {
		return m.Get(ctx, id)
	}
	cur, err := m.Get(ctx, id)
	if err != nil {
		return agent.Agent{}, err
	}
	next := u.Apply(cur)
	if err := next.Validate(); err != nil {
		return agent.Agent{}, err
	}
	next.UpdatedAt = m.now().UTC()
	const q = `UPDATE agents SET name = ?, class = ?, money_balance = ?, hoarding = ?, labour_minutes = ?, updated_at = ?
		WHERE id = ?`
	res, err := m.db.ExecContext(ctx, q,
		next.Name, string(next.Class), int64(next.MoneyBalance), next.Hoarding,
		next.LabourMinutes, next.UpdatedAt, string(id),
	)
	if err != nil {
		return agent.Agent{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return agent.Agent{}, ErrNotFound
	}
	return next, nil
}

func (m *MySQL) Delete(ctx context.Context, id agent.ID) error {
	res, err := m.db.ExecContext(ctx, `DELETE FROM agents WHERE id = ?`, string(id))
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

// CreateCircuit atomically inserts the circuit and updates the agent's
// money_balance by circuit.SurplusValue in a single transaction.
func (m *MySQL) CreateCircuit(ctx context.Context, c agent.CapitalCircuit) (agent.CapitalCircuit, error) {
	if err := c.Validate(); err != nil {
		return agent.CapitalCircuit{}, err
	}
	if c.ID.IsZero() {
		c.ID = agent.NewID()
	}
	c.CreatedAt = m.now().UTC()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return agent.CapitalCircuit{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var bal int64
	err = tx.QueryRowContext(ctx,
		`SELECT money_balance FROM agents WHERE id = ? FOR UPDATE`, string(c.AgentID),
	).Scan(&bal)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.CapitalCircuit{}, ErrNotFound
	}
	if err != nil {
		return agent.CapitalCircuit{}, err
	}
	if agent.Pence(bal) < c.MAdvanced {
		return agent.CapitalCircuit{}, agent.ErrInsufficientFunds
	}

	const insertQ = `INSERT INTO capital_circuits
		(id, agent_id, m_advanced, commodity_id, m_returned, surplus_value, circuit_type, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, insertQ,
		string(c.ID), string(c.AgentID), int64(c.MAdvanced),
		c.CommodityID, int64(c.MReturned), int64(c.SurplusValue),
		string(c.CircuitType), c.CreatedAt,
	)
	if err != nil {
		return agent.CapitalCircuit{}, err
	}

	now := m.now().UTC()
	_, err = tx.ExecContext(ctx,
		`UPDATE agents SET money_balance = money_balance + ?, updated_at = ? WHERE id = ?`,
		int64(c.SurplusValue), now, string(c.AgentID),
	)
	if err != nil {
		return agent.CapitalCircuit{}, err
	}
	if err := tx.Commit(); err != nil {
		return agent.CapitalCircuit{}, err
	}
	return c, nil
}

func (m *MySQL) GetCircuit(ctx context.Context, id agent.ID) (agent.CapitalCircuit, error) {
	const q = `SELECT id, agent_id, m_advanced, commodity_id, m_returned, surplus_value, circuit_type, created_at
		FROM capital_circuits WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanCircuit(row)
}

func (m *MySQL) ListCircuits(ctx context.Context, agentID agent.ID) ([]agent.CapitalCircuit, error) {
	const q = `SELECT id, agent_id, m_advanced, commodity_id, m_returned, surplus_value, circuit_type, created_at
		FROM capital_circuits WHERE agent_id = ? ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(agentID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []agent.CapitalCircuit
	for rows.Next() {
		c, err := scanCircuitRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func scanAgent(row *sql.Row) (agent.Agent, error) {
	var a agent.Agent
	var id, class string
	var balance int64
	var hoarding bool
	err := row.Scan(&id, &a.Name, &class, &balance, &hoarding, &a.LabourMinutes, &a.CreatedAt, &a.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.Agent{}, ErrNotFound
	}
	if err != nil {
		return agent.Agent{}, err
	}
	a.ID = agent.ID(id)
	a.Class = agent.Class(class)
	a.MoneyBalance = agent.Pence(balance)
	a.Hoarding = hoarding
	return a, nil
}

func scanAgents(rows *sql.Rows) ([]agent.Agent, error) {
	var out []agent.Agent
	for rows.Next() {
		var a agent.Agent
		var id, class string
		var balance int64
		var hoarding bool
		if err := rows.Scan(&id, &a.Name, &class, &balance, &hoarding, &a.LabourMinutes, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, err
		}
		a.ID = agent.ID(id)
		a.Class = agent.Class(class)
		a.MoneyBalance = agent.Pence(balance)
		a.Hoarding = hoarding
		out = append(out, a)
	}
	return out, rows.Err()
}

func scanCircuit(row *sql.Row) (agent.CapitalCircuit, error) {
	var c agent.CapitalCircuit
	var id, agentID, circuitType string
	var mAdv, mRet, sv int64
	err := row.Scan(&id, &agentID, &mAdv, &c.CommodityID, &mRet, &sv, &circuitType, &c.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.CapitalCircuit{}, ErrNotFound
	}
	if err != nil {
		return agent.CapitalCircuit{}, err
	}
	c.ID = agent.ID(id)
	c.AgentID = agent.ID(agentID)
	c.MAdvanced = agent.Pence(mAdv)
	c.MReturned = agent.Pence(mRet)
	c.SurplusValue = agent.Pence(sv)
	c.CircuitType = agent.CircuitType(circuitType)
	return c, nil
}

func scanCircuitRow(rows *sql.Rows) (agent.CapitalCircuit, error) {
	var c agent.CapitalCircuit
	var id, agentID, circuitType string
	var mAdv, mRet, sv int64
	if err := rows.Scan(&id, &agentID, &mAdv, &c.CommodityID, &mRet, &sv, &circuitType, &c.CreatedAt); err != nil {
		return agent.CapitalCircuit{}, err
	}
	c.ID = agent.ID(id)
	c.AgentID = agent.ID(agentID)
	c.MAdvanced = agent.Pence(mAdv)
	c.MReturned = agent.Pence(mRet)
	c.SurplusValue = agent.Pence(sv)
	c.CircuitType = agent.CircuitType(circuitType)
	return c, nil
}

func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "1062") || strings.Contains(s, "Duplicate entry")
}
