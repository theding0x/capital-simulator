package store

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
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

func (m *MySQL) CreateWorker(ctx context.Context, w agent.Worker) (agent.Worker, error) {
	if err := w.Validate(); err != nil {
		return agent.Worker{}, err
	}
	if w.ID.IsZero() {
		w.ID = agent.NewAgentID()
	}
	now := m.now().UTC()
	w.CreatedAt = now
	w.UpdatedAt = now
	w.Kind = agent.AgentKindWorker
	const q = `INSERT INTO labour_workers
		(id, kind, owns_labour_power, owns_commodities_to_sell,
		 capacity_minutes_per_day, labour_power_value_minutes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(w.ID), string(w.Kind),
		w.OwnsLabourPower, w.OwnsCommoditiesToSell,
		int64(w.LabourPower.CapacityMinutesPerDay),
		int64(w.LabourPowerValueMinutes),
		w.CreatedAt, w.UpdatedAt,
	)
	if err != nil {
		return agent.Worker{}, err
	}
	return w, nil
}

func (m *MySQL) GetWorker(ctx context.Context, id agent.AgentID) (agent.Worker, error) {
	const q = `SELECT id, kind, owns_labour_power, owns_commodities_to_sell,
		capacity_minutes_per_day, labour_power_value_minutes, created_at, updated_at
		FROM labour_workers WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanWorker(row)
}

func (m *MySQL) ListWorkers(ctx context.Context) ([]agent.Worker, error) {
	const q = `SELECT id, kind, owns_labour_power, owns_commodities_to_sell,
		capacity_minutes_per_day, labour_power_value_minutes, created_at, updated_at
		FROM labour_workers ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []agent.Worker
	for rows.Next() {
		w, err := scanWorkerRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

func (m *MySQL) CreateCapitalist(ctx context.Context, c agent.Capitalist) (agent.Capitalist, error) {
	if err := c.Validate(); err != nil {
		return agent.Capitalist{}, err
	}
	if c.ID.IsZero() {
		c.ID = agent.NewAgentID()
	}
	now := m.now().UTC()
	c.CreatedAt = now
	c.UpdatedAt = now
	c.Kind = agent.AgentKindCapitalist
	const q = `INSERT INTO labour_capitalists
		(id, kind, money_capital, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID), string(c.Kind),
		int64(c.MoneyCapital),
		c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return agent.Capitalist{}, err
	}
	return c, nil
}

func (m *MySQL) GetCapitalist(ctx context.Context, id agent.AgentID) (agent.Capitalist, error) {
	const q = `SELECT id, kind, money_capital, created_at, updated_at
		FROM labour_capitalists WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanCapitalist(row)
}

func (m *MySQL) ListCapitalists(ctx context.Context) ([]agent.Capitalist, error) {
	const q = `SELECT id, kind, money_capital, created_at, updated_at
		FROM labour_capitalists ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []agent.Capitalist
	for rows.Next() {
		c, err := scanCapitalistRow(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (m *MySQL) CreateOffering(ctx context.Context, o agent.LabourPowerOffering) (agent.LabourPowerOffering, error) {
	if err := o.Validate(); err != nil {
		return agent.LabourPowerOffering{}, err
	}
	if o.ID.IsZero() {
		o.ID = agent.NewAgentID()
	}
	o.CreatedAt = m.now().UTC()
	const q = `INSERT INTO labour_power_offerings
		(id, owner_id, capacity_minutes_per_day, contract_days, asking_wage, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(o.ID), string(o.OwnerID),
		int64(o.CapacityMinutesPerDay), o.ContractDays,
		int64(o.AskingWage), o.CreatedAt,
	)
	if err != nil {
		return agent.LabourPowerOffering{}, err
	}
	return o, nil
}

func (m *MySQL) ListOfferings(ctx context.Context) ([]agent.LabourPowerOffering, error) {
	const q = `SELECT id, owner_id, capacity_minutes_per_day, contract_days, asking_wage, created_at
		FROM labour_power_offerings ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []agent.LabourPowerOffering
	for rows.Next() {
		var o agent.LabourPowerOffering
		var id, ownerID string
		var cap, wage int64
		if err := rows.Scan(&id, &ownerID, &cap, &o.ContractDays, &wage, &o.CreatedAt); err != nil {
			return nil, err
		}
		o.ID = agent.AgentID(id)
		o.OwnerID = agent.AgentID(ownerID)
		o.CapacityMinutesPerDay = agent.LabourMinutes(cap)
		o.AskingWage = agent.LabourMinutes(wage)
		out = append(out, o)
	}
	return out, rows.Err()
}

func (m *MySQL) CreatePurchase(ctx context.Context, p agent.LabourPowerPurchase) (agent.LabourPowerPurchase, error) {
	if err := p.Validate(); err != nil {
		return agent.LabourPowerPurchase{}, err
	}
	if p.ID.IsZero() {
		p.ID = agent.NewPurchaseID()
	}
	p.CreatedAt = m.now().UTC()
	const q = `INSERT INTO labour_power_purchases
		(id, seller_id, buyer_id, wage_minutes, contract_days, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(p.ID), string(p.SellerID), string(p.BuyerID),
		int64(p.WageMinutes), p.ContractDays, p.CreatedAt,
	)
	if err != nil {
		return agent.LabourPowerPurchase{}, err
	}
	return p, nil
}

func (m *MySQL) GetPurchase(ctx context.Context, id agent.PurchaseID) (agent.LabourPowerPurchase, error) {
	const q = `SELECT id, seller_id, buyer_id, wage_minutes, contract_days, created_at
		FROM labour_power_purchases WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	var p agent.LabourPowerPurchase
	var pid, sellerID, buyerID string
	var wage int64
	err := row.Scan(&pid, &sellerID, &buyerID, &wage, &p.ContractDays, &p.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.LabourPowerPurchase{}, ErrNotFound
	}
	if err != nil {
		return agent.LabourPowerPurchase{}, err
	}
	p.ID = agent.PurchaseID(pid)
	p.SellerID = agent.AgentID(sellerID)
	p.BuyerID = agent.AgentID(buyerID)
	p.WageMinutes = agent.LabourMinutes(wage)
	return p, nil
}

func (m *MySQL) ListPurchases(ctx context.Context) ([]agent.LabourPowerPurchase, error) {
	const q = `SELECT id, seller_id, buyer_id, wage_minutes, contract_days, created_at
		FROM labour_power_purchases ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []agent.LabourPowerPurchase
	for rows.Next() {
		var p agent.LabourPowerPurchase
		var pid, sellerID, buyerID string
		var wage int64
		if err := rows.Scan(&pid, &sellerID, &buyerID, &wage, &p.ContractDays, &p.CreatedAt); err != nil {
			return nil, err
		}
		p.ID = agent.PurchaseID(pid)
		p.SellerID = agent.AgentID(sellerID)
		p.BuyerID = agent.AgentID(buyerID)
		p.WageMinutes = agent.LabourMinutes(wage)
		out = append(out, p)
	}
	return out, rows.Err()
}

func scanWorker(row *sql.Row) (agent.Worker, error) {
	var w agent.Worker
	var id, kind string
	var cap, lpv int64
	err := row.Scan(&id, &kind, &w.OwnsLabourPower, &w.OwnsCommoditiesToSell,
		&cap, &lpv, &w.CreatedAt, &w.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.Worker{}, ErrNotFound
	}
	if err != nil {
		return agent.Worker{}, err
	}
	w.ID = agent.AgentID(id)
	w.Kind = agent.AgentKind(kind)
	w.LabourPower.CapacityMinutesPerDay = agent.LabourMinutes(cap)
	w.LabourPowerValueMinutes = agent.LabourMinutes(lpv)
	return w, nil
}

func scanWorkerRow(rows *sql.Rows) (agent.Worker, error) {
	var w agent.Worker
	var id, kind string
	var cap, lpv int64
	if err := rows.Scan(&id, &kind, &w.OwnsLabourPower, &w.OwnsCommoditiesToSell,
		&cap, &lpv, &w.CreatedAt, &w.UpdatedAt); err != nil {
		return agent.Worker{}, err
	}
	w.ID = agent.AgentID(id)
	w.Kind = agent.AgentKind(kind)
	w.LabourPower.CapacityMinutesPerDay = agent.LabourMinutes(cap)
	w.LabourPowerValueMinutes = agent.LabourMinutes(lpv)
	return w, nil
}

func scanCapitalist(row *sql.Row) (agent.Capitalist, error) {
	var c agent.Capitalist
	var id, kind string
	var mc int64
	err := row.Scan(&id, &kind, &mc, &c.CreatedAt, &c.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.Capitalist{}, ErrNotFound
	}
	if err != nil {
		return agent.Capitalist{}, err
	}
	c.ID = agent.AgentID(id)
	c.Kind = agent.AgentKind(kind)
	c.MoneyCapital = agent.LabourMinutes(mc)
	return c, nil
}

func scanCapitalistRow(rows *sql.Rows) (agent.Capitalist, error) {
	var c agent.Capitalist
	var id, kind string
	var mc int64
	if err := rows.Scan(&id, &kind, &mc, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return agent.Capitalist{}, err
	}
	c.ID = agent.AgentID(id)
	c.Kind = agent.AgentKind(kind)
	c.MoneyCapital = agent.LabourMinutes(mc)
	return c, nil
}

func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "1062") || strings.Contains(s, "Duplicate entry")
}

func (m *MySQL) CreateLabourProcess(ctx context.Context, lp agent.LabourProcess) (agent.LabourProcess, error) {
	if err := lp.Validate(); err != nil {
		return agent.LabourProcess{}, err
	}
	if lp.ID.IsZero() {
		lp.ID = agent.NewLabourProcessID()
	}
	lp.CreatedAt = m.now().UTC()
	meansJSON, err := json.Marshal(lp.Means)
	if err != nil {
		return agent.LabourProcess{}, err
	}
	const q = `INSERT INTO labour_processes
		(id, worker_id, capitalist_id, duration, necessary_labour_minutes,
		 means_json, product_kind, product_quantity, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = m.db.ExecContext(ctx, q,
		string(lp.ID), string(lp.WorkerID), string(lp.CapitalistID),
		int64(lp.Duration), int64(lp.NecessaryLabourMinutes),
		string(meansJSON), lp.ProductKind, lp.ProductQuantity,
		lp.CreatedAt,
	)
	if err != nil {
		return agent.LabourProcess{}, err
	}
	return lp, nil
}

func (m *MySQL) GetLabourProcess(ctx context.Context, id agent.LabourProcessID) (agent.LabourProcess, error) {
	const q = `SELECT id, worker_id, capitalist_id, duration,
		necessary_labour_minutes, means_json, product_kind, product_quantity, created_at
		FROM labour_processes WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanLabourProcess(row)
}

func scanLabourProcess(row *sql.Row) (agent.LabourProcess, error) {
	var lp agent.LabourProcess
	var id, workerID, capitalistID, meansJSON string
	var dur, nl int64
	err := row.Scan(&id, &workerID, &capitalistID, &dur, &nl,
		&meansJSON, &lp.ProductKind, &lp.ProductQuantity, &lp.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.LabourProcess{}, ErrNotFound
	}
	if err != nil {
		return agent.LabourProcess{}, err
	}
	lp.ID = agent.LabourProcessID(id)
	lp.WorkerID = agent.AgentID(workerID)
	lp.CapitalistID = agent.AgentID(capitalistID)
	lp.Duration = agent.LabourMinutes(dur)
	lp.NecessaryLabourMinutes = agent.LabourMinutes(nl)
	if err := json.Unmarshal([]byte(meansJSON), &lp.Means); err != nil {
		return agent.LabourProcess{}, err
	}
	return lp, nil
}

func (m *MySQL) CreateWorkingDay(ctx context.Context, wd agent.WorkingDay) (agent.WorkingDay, error) {
	if err := wd.Validate(); err != nil {
		return agent.WorkingDay{}, err
	}
	if wd.ID.IsZero() {
		wd.ID = agent.NewWorkingDayID()
	}
	wd.CreatedAt = m.now().UTC()
	const q = `INSERT INTO working_days (id, necessary_labour_minutes, surplus_labour_minutes, created_at)
		VALUES (?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(wd.ID),
		int64(wd.NecessaryLabourMinutes),
		int64(wd.SurplusLabourMinutes),
		wd.CreatedAt,
	)
	if err != nil {
		return agent.WorkingDay{}, err
	}
	return wd, nil
}

func (m *MySQL) GetWorkingDay(ctx context.Context, id agent.WorkingDayID) (agent.WorkingDay, error) {
	const q = `SELECT id, necessary_labour_minutes, surplus_labour_minutes, created_at
		FROM working_days WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	var wd agent.WorkingDay
	var wid string
	var nl, sl int64
	err := row.Scan(&wid, &nl, &sl, &wd.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.WorkingDay{}, ErrNotFound
	}
	if err != nil {
		return agent.WorkingDay{}, err
	}
	wd.ID = agent.WorkingDayID(wid)
	wd.NecessaryLabourMinutes = agent.NecessaryLabourMinutes(nl)
	wd.SurplusLabourMinutes = agent.SurplusLabourMinutes(sl)
	return wd, nil
}

func (m *MySQL) CreateRelaySchedule(ctx context.Context, rs agent.RelaySchedule) (agent.RelaySchedule, error) {
	if err := rs.Validate(); err != nil {
		return agent.RelaySchedule{}, err
	}
	if rs.ID.IsZero() {
		rs.ID = agent.NewRelayScheduleID()
	}
	rs.CreatedAt = m.now().UTC()
	wids0, err := json.Marshal(rs.Sets[0].WorkerIDs)
	if err != nil {
		return agent.RelaySchedule{}, err
	}
	wids1, err := json.Marshal(rs.Sets[1].WorkerIDs)
	if err != nil {
		return agent.RelaySchedule{}, err
	}
	const q = `INSERT INTO relay_schedules
		(id, shift_kind_0, nl_minutes_0, sl_minutes_0, worker_ids_0,
		     shift_kind_1, nl_minutes_1, sl_minutes_1, worker_ids_1, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = m.db.ExecContext(ctx, q,
		string(rs.ID),
		string(rs.Sets[0].ShiftKind),
		int64(rs.Sets[0].WorkingDay.NecessaryLabourMinutes),
		int64(rs.Sets[0].WorkingDay.SurplusLabourMinutes),
		string(wids0),
		string(rs.Sets[1].ShiftKind),
		int64(rs.Sets[1].WorkingDay.NecessaryLabourMinutes),
		int64(rs.Sets[1].WorkingDay.SurplusLabourMinutes),
		string(wids1),
		rs.CreatedAt,
	)
	if err != nil {
		return agent.RelaySchedule{}, err
	}
	return rs, nil
}

func (m *MySQL) CreateCooperation(ctx context.Context, c agent.Cooperation) (agent.Cooperation, error) {
	if err := c.Validate(); err != nil {
		return agent.Cooperation{}, err
	}
	if c.ID.IsZero() {
		c.ID = agent.NewCooperationID()
	}
	c.CreatedAt = m.now().UTC()
	membersJSON, err := json.Marshal(c.Members)
	if err != nil {
		return agent.Cooperation{}, err
	}
	const q = `INSERT INTO cooperations
		(id, name, capitalist_id, members_json, created_at)
		VALUES (?, ?, ?, ?, ?)`
	if _, err := m.db.ExecContext(ctx, q,
		string(c.ID), c.Name, string(c.CapitalistID),
		string(membersJSON), c.CreatedAt,
	); err != nil {
		return agent.Cooperation{}, err
	}
	return c, nil
}

func (m *MySQL) GetCooperation(ctx context.Context, id agent.CooperationID) (agent.Cooperation, error) {
	const q = `SELECT id, name, capitalist_id, members_json, created_at
		FROM cooperations WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanCooperation(row)
}

func (m *MySQL) ListCooperations(ctx context.Context) ([]agent.Cooperation, error) {
	const q = `SELECT id, name, capitalist_id, members_json, created_at
		FROM cooperations ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCooperations(rows)
}

func (m *MySQL) ListCooperationsByCapitalist(ctx context.Context, capitalistID agent.AgentID) ([]agent.Cooperation, error) {
	const q = `SELECT id, name, capitalist_id, members_json, created_at
		FROM cooperations WHERE capitalist_id = ? ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(capitalistID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanCooperations(rows)
}

func scanCooperation(row *sql.Row) (agent.Cooperation, error) {
	var c agent.Cooperation
	var id, capID, membersRaw string
	err := row.Scan(&id, &c.Name, &capID, &membersRaw, &c.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.Cooperation{}, ErrNotFound
	}
	if err != nil {
		return agent.Cooperation{}, err
	}
	c.ID = agent.CooperationID(id)
	c.CapitalistID = agent.AgentID(capID)
	if err := json.Unmarshal([]byte(membersRaw), &c.Members); err != nil {
		return agent.Cooperation{}, err
	}
	return c, nil
}

func scanCooperations(rows *sql.Rows) ([]agent.Cooperation, error) {
	var out []agent.Cooperation
	for rows.Next() {
		var c agent.Cooperation
		var id, capID, membersRaw string
		if err := rows.Scan(&id, &c.Name, &capID, &membersRaw, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.ID = agent.CooperationID(id)
		c.CapitalistID = agent.AgentID(capID)
		if err := json.Unmarshal([]byte(membersRaw), &c.Members); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func (m *MySQL) CreateManufacture(ctx context.Context, mf agent.Manufacture) (agent.Manufacture, error) {
	if err := mf.Validate(); err != nil {
		return agent.Manufacture{}, err
	}
	if mf.ID.IsZero() {
		mf.ID = agent.NewManufactureID()
	}
	mf.CreatedAt = m.now().UTC()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return agent.Manufacture{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const insertManufacture = `INSERT INTO manufactures
		(id, name, capitalist_id, form, origin, individual_working_day_minutes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	if _, err := tx.ExecContext(ctx, insertManufacture,
		string(mf.ID), mf.Name, string(mf.CapitalistID),
		string(mf.Form), string(mf.Origin),
		int64(mf.IndividualWorkingDayMinutes),
		mf.CreatedAt,
	); err != nil {
		return agent.Manufacture{}, err
	}

	const insertRole = `INSERT INTO detail_roles
		(manufacture_id, ordinal, role_name, skill_level, output_rate_per_hour, head_count, tool_name)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	for i, r := range mf.Roles {
		if _, err := tx.ExecContext(ctx, insertRole,
			string(mf.ID), i, r.Name, string(r.SkillLevel),
			r.OutputRatePerHour, r.HeadCount, r.ToolName,
		); err != nil {
			return agent.Manufacture{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return agent.Manufacture{}, err
	}
	return mf, nil
}

func (m *MySQL) GetManufacture(ctx context.Context, id agent.ManufactureID) (agent.Manufacture, error) {
	const q = `SELECT id, name, capitalist_id, form, origin, individual_working_day_minutes, created_at
		FROM manufactures WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	mf, err := scanManufactureRow(row)
	if err != nil {
		return agent.Manufacture{}, err
	}
	roles, err := m.loadDetailRoles(ctx, mf.ID)
	if err != nil {
		return agent.Manufacture{}, err
	}
	mf.Roles = roles
	return mf, nil
}

func (m *MySQL) ListManufactures(ctx context.Context) ([]agent.Manufacture, error) {
	const q = `SELECT id, name, capitalist_id, form, origin, individual_working_day_minutes, created_at
		FROM manufactures ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return m.scanManufactures(ctx, rows)
}

func (m *MySQL) ListManufacturesByCapitalist(ctx context.Context, capitalistID agent.AgentID) ([]agent.Manufacture, error) {
	const q = `SELECT id, name, capitalist_id, form, origin, individual_working_day_minutes, created_at
		FROM manufactures WHERE capitalist_id = ? ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(capitalistID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return m.scanManufactures(ctx, rows)
}

func (m *MySQL) ListManufacturesByForm(ctx context.Context, form agent.ManufactureForm) ([]agent.Manufacture, error) {
	const q = `SELECT id, name, capitalist_id, form, origin, individual_working_day_minutes, created_at
		FROM manufactures WHERE form = ? ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(form))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return m.scanManufactures(ctx, rows)
}

func (m *MySQL) loadDetailRoles(ctx context.Context, id agent.ManufactureID) ([]agent.DetailRole, error) {
	const q = `SELECT role_name, skill_level, output_rate_per_hour, head_count, tool_name
		FROM detail_roles WHERE manufacture_id = ? ORDER BY ordinal ASC`
	rows, err := m.db.QueryContext(ctx, q, string(id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []agent.DetailRole
	for rows.Next() {
		var r agent.DetailRole
		var skill string
		if err := rows.Scan(&r.Name, &skill, &r.OutputRatePerHour, &r.HeadCount, &r.ToolName); err != nil {
			return nil, err
		}
		r.SkillLevel = agent.SkillLevel(skill)
		out = append(out, r)
	}
	return out, rows.Err()
}

func (m *MySQL) scanManufactures(ctx context.Context, rows *sql.Rows) ([]agent.Manufacture, error) {
	var out []agent.Manufacture
	for rows.Next() {
		var mf agent.Manufacture
		var id, capID, form, origin string
		var wd int64
		if err := rows.Scan(&id, &mf.Name, &capID, &form, &origin, &wd, &mf.CreatedAt); err != nil {
			return nil, err
		}
		mf.ID = agent.ManufactureID(id)
		mf.CapitalistID = agent.AgentID(capID)
		mf.Form = agent.ManufactureForm(form)
		mf.Origin = agent.ManufactureOrigin(origin)
		mf.IndividualWorkingDayMinutes = agent.LabourMinutes(wd)
		out = append(out, mf)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i, mf := range out {
		roles, err := m.loadDetailRoles(ctx, mf.ID)
		if err != nil {
			return nil, err
		}
		out[i].Roles = roles
	}
	return out, nil
}

func scanManufactureRow(row *sql.Row) (agent.Manufacture, error) {
	var mf agent.Manufacture
	var id, capID, form, origin string
	var wd int64
	err := row.Scan(&id, &mf.Name, &capID, &form, &origin, &wd, &mf.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.Manufacture{}, ErrNotFound
	}
	if err != nil {
		return agent.Manufacture{}, err
	}
	mf.ID = agent.ManufactureID(id)
	mf.CapitalistID = agent.AgentID(capID)
	mf.Form = agent.ManufactureForm(form)
	mf.Origin = agent.ManufactureOrigin(origin)
	mf.IndividualWorkingDayMinutes = agent.LabourMinutes(wd)
	return mf, nil
}

func (m *MySQL) GetRelaySchedule(ctx context.Context, id agent.RelayScheduleID) (agent.RelaySchedule, error) {
	const q = `SELECT id, shift_kind_0, nl_minutes_0, sl_minutes_0, worker_ids_0,
		shift_kind_1, nl_minutes_1, sl_minutes_1, worker_ids_1, created_at
		FROM relay_schedules WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	var rs agent.RelaySchedule
	var rid, sk0, sk1 string
	var nl0, sl0, nl1, sl1 int64
	var wids0raw, wids1raw string
	err := row.Scan(&rid, &sk0, &nl0, &sl0, &wids0raw, &sk1, &nl1, &sl1, &wids1raw, &rs.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.RelaySchedule{}, ErrNotFound
	}
	if err != nil {
		return agent.RelaySchedule{}, err
	}
	rs.ID = agent.RelayScheduleID(rid)
	rs.Sets[0].ShiftKind = agent.ShiftKind(sk0)
	rs.Sets[0].WorkingDay.NecessaryLabourMinutes = agent.NecessaryLabourMinutes(nl0)
	rs.Sets[0].WorkingDay.SurplusLabourMinutes = agent.SurplusLabourMinutes(sl0)
	rs.Sets[1].ShiftKind = agent.ShiftKind(sk1)
	rs.Sets[1].WorkingDay.NecessaryLabourMinutes = agent.NecessaryLabourMinutes(nl1)
	rs.Sets[1].WorkingDay.SurplusLabourMinutes = agent.SurplusLabourMinutes(sl1)
	if err := json.Unmarshal([]byte(wids0raw), &rs.Sets[0].WorkerIDs); err != nil {
		return agent.RelaySchedule{}, err
	}
	if err := json.Unmarshal([]byte(wids1raw), &rs.Sets[1].WorkerIDs); err != nil {
		return agent.RelaySchedule{}, err
	}
	return rs, nil
}

func (m *MySQL) CreateWageForm(ctx context.Context, wf agent.WageForm) (agent.WageForm, error) {
	if err := wf.Validate(); err != nil {
		return agent.WageForm{}, err
	}
	if wf.ID.IsZero() {
		wf.ID = agent.NewWageFormID()
	}
	wf.CreatedAt = m.now().UTC()
	const q = `INSERT INTO wage_forms
		(id, agent_id, daily_pence, working_day_hours, lpv_daily_pence, necessary_minutes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			id = VALUES(id),
			daily_pence = VALUES(daily_pence),
			working_day_hours = VALUES(working_day_hours),
			lpv_daily_pence = VALUES(lpv_daily_pence),
			necessary_minutes = VALUES(necessary_minutes),
			created_at = VALUES(created_at)`
	_, err := m.db.ExecContext(ctx, q,
		string(wf.ID), string(wf.AgentID),
		int64(wf.Wage.DailyPence), wf.Wage.WorkingDayHours,
		int64(wf.LabourPowerValue.DailyPence), int64(wf.LabourPowerValue.NecessaryMinutes),
		wf.CreatedAt,
	)
	if err != nil {
		return agent.WageForm{}, err
	}
	return wf, nil
}

func (m *MySQL) GetWageForm(ctx context.Context, agentID agent.AgentID) (agent.WageForm, error) {
	const q = `SELECT id, agent_id, daily_pence, working_day_hours, lpv_daily_pence, necessary_minutes, created_at
		FROM wage_forms WHERE agent_id = ?`
	row := m.db.QueryRowContext(ctx, q, string(agentID))
	var wf agent.WageForm
	var id, aid string
	var dp, wdh, lpvDp, nm int64
	err := row.Scan(&id, &aid, &dp, &wdh, &lpvDp, &nm, &wf.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.WageForm{}, ErrNotFound
	}
	if err != nil {
		return agent.WageForm{}, err
	}
	wf.ID = agent.WageFormID(id)
	wf.AgentID = agent.AgentID(aid)
	wf.Wage = agent.Wage{DailyPence: agent.Pence(dp), WorkingDayHours: wdh}
	wf.LabourPowerValue = agent.WageLabourValue{
		DailyPence:       agent.Pence(lpvDp),
		NecessaryMinutes: agent.LabourMinutes(nm),
	}
	return wf, nil
}

func (m *MySQL) CreateWorkingSession(ctx context.Context, s agent.WorkingSession) (agent.WorkingSession, error) {
	if err := s.Validate(); err != nil {
		return agent.WorkingSession{}, err
	}
	if s.ID.IsZero() {
		s.ID = agent.NewWorkingSessionID()
	}
	s.CreatedAt = m.now().UTC()
	const q = `INSERT INTO time_wage_sessions
		(id, agent_id, daily_labour_power_value, working_day_hours,
		 overtime_hours, overtime_rate_pence, wage_period, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(s.ID), string(s.AgentID),
		s.DailyLabourPowerValue.Pence,
		s.WorkingDayHours.Hours,
		s.OvertimeHours.Hours,
		s.OvertimeRatePence.Pence,
		string(s.WagePeriod),
		s.CreatedAt,
	)
	if err != nil {
		return agent.WorkingSession{}, err
	}
	return s, nil
}

func (m *MySQL) GetWorkingSession(ctx context.Context, id agent.WorkingSessionID) (agent.WorkingSession, error) {
	const q = `SELECT id, agent_id, daily_labour_power_value, working_day_hours,
		overtime_hours, overtime_rate_pence, wage_period, created_at
		FROM time_wage_sessions WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	var s agent.WorkingSession
	var sid, aid, period string
	var dlpv, wdh, oth, orp int64
	err := row.Scan(&sid, &aid, &dlpv, &wdh, &oth, &orp, &period, &s.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.WorkingSession{}, ErrNotFound
	}
	if err != nil {
		return agent.WorkingSession{}, err
	}
	s.ID = agent.WorkingSessionID(sid)
	s.AgentID = agent.AgentID(aid)
	s.DailyLabourPowerValue = agent.DailyLabourPowerValue{Pence: dlpv}
	s.WorkingDayHours = agent.WorkingDayHours{Hours: wdh}
	s.OvertimeHours = agent.OvertimeHours{Hours: oth}
	s.OvertimeRatePence = agent.OvertimeRatePence{Pence: orp}
	s.WagePeriod = agent.WagePeriod(period)
	return s, nil
}

func (m *MySQL) CreatePieceWage(ctx context.Context, pw agent.PieceWage) (agent.PieceWage, error) {
	if err := pw.Validate(); err != nil {
		return agent.PieceWage{}, err
	}
	if pw.ID.IsZero() {
		pw.ID = agent.NewPieceWageID()
	}
	pw.CreatedAt = m.now().UTC()
	const q = `INSERT INTO piece_wages (id, agent_id, price_pence, normal_output, created_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(pw.ID), string(pw.AgentID),
		pw.PricePence, pw.NormalOutput, pw.CreatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return agent.PieceWage{}, ErrAlreadyExists
		}
		return agent.PieceWage{}, err
	}
	return pw, nil
}

func (m *MySQL) GetPieceWage(ctx context.Context, agentID agent.AgentID) (agent.PieceWage, error) {
	const q = `SELECT id, agent_id, price_pence, normal_output, created_at
		FROM piece_wages WHERE agent_id = ?`
	row := m.db.QueryRowContext(ctx, q, string(agentID))
	var pw agent.PieceWage
	var id, aid string
	err := row.Scan(&id, &aid, &pw.PricePence, &pw.NormalOutput, &pw.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.PieceWage{}, ErrNotFound
	}
	if err != nil {
		return agent.PieceWage{}, err
	}
	pw.ID = agent.PieceWageID(id)
	pw.AgentID = agent.AgentID(aid)
	return pw, nil
}

func (m *MySQL) CreateSubContract(ctx context.Context, sc agent.SubContract) (agent.SubContract, error) {
	if err := sc.Validate(); err != nil {
		return agent.SubContract{}, err
	}
	if sc.ID.IsZero() {
		sc.ID = agent.NewSubContractID()
	}
	sc.CreatedAt = m.now().UTC()
	assistantJSON, err := json.Marshal(sc.AssistantIDs)
	if err != nil {
		return agent.SubContract{}, err
	}
	const q = `INSERT INTO sub_contracts
		(id, head_labourer_id, assistant_ids, piece_rate_pence, assistant_rate_pence, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err = m.db.ExecContext(ctx, q,
		string(sc.ID), string(sc.HeadLabourerID),
		string(assistantJSON),
		sc.PieceRatePence, sc.AssistantRatePence, sc.CreatedAt,
	)
	if err != nil {
		return agent.SubContract{}, err
	}
	return sc, nil
}

func (m *MySQL) GetSubContract(ctx context.Context, id agent.SubContractID) (agent.SubContract, error) {
	const q = `SELECT id, head_labourer_id, assistant_ids, piece_rate_pence, assistant_rate_pence, created_at
		FROM sub_contracts WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	var sc agent.SubContract
	var scID, headID, assistJSON string
	err := row.Scan(&scID, &headID, &assistJSON, &sc.PieceRatePence, &sc.AssistantRatePence, &sc.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return agent.SubContract{}, ErrNotFound
	}
	if err != nil {
		return agent.SubContract{}, err
	}
	sc.ID = agent.SubContractID(scID)
	sc.HeadLabourerID = agent.AgentID(headID)
	var rawIDs []string
	if err := json.Unmarshal([]byte(assistJSON), &rawIDs); err != nil {
		return agent.SubContract{}, err
	}
	sc.AssistantIDs = make([]agent.AgentID, len(rawIDs))
	for i, r := range rawIDs {
		sc.AssistantIDs[i] = agent.AgentID(r)
	}
	return sc, nil
}
