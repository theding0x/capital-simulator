package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"io/fs"
	"sort"
	"time"

	pkgmysql "github.com/theding0x/capital-simulator/pkg/mysql"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

//go:embed migrations
var migrationsFS embed.FS

// MySQL implements MachineStore and FactoryStore backed by MySQL.
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

func (m *MySQL) CreateMachine(ctx context.Context, mc machinery.Machine) (machinery.Machine, error) {
	if err := mc.Validate(); err != nil {
		return machinery.Machine{}, err
	}
	if mc.ID.IsZero() {
		mc.ID = machinery.NewMachineID()
	}
	mc.CreatedAt = m.now().UTC()
	const q = `INSERT INTO machines
		(id, name, motor_mechanism, transmitting_mechanism, working_tool,
		 machine_value, lifespan_days, productive_power, hand_labour_per_unit,
		 accumulated_wear, accumulated_depreciation, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(mc.ID), mc.Name, mc.MotorMechanism, mc.TransmittingMechanism, mc.WorkingTool,
		int64(mc.MachineValue), int64(mc.LifespanDays), int64(mc.ProductivePower),
		int64(mc.HandLabourPerUnit),
		int64(mc.AccumulatedWear.Value), int64(mc.AccumulatedDepreciation.Value),
		mc.CreatedAt,
	)
	if err != nil {
		return machinery.Machine{}, err
	}
	return mc, nil
}

func (m *MySQL) GetMachine(ctx context.Context, id machinery.MachineID) (machinery.Machine, error) {
	const q = `SELECT id, name, motor_mechanism, transmitting_mechanism, working_tool,
		machine_value, lifespan_days, productive_power, hand_labour_per_unit,
		accumulated_wear, accumulated_depreciation, created_at
		FROM machines WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanMachine(row.Scan)
}

func (m *MySQL) ListMachines(ctx context.Context) ([]machinery.Machine, error) {
	const q = `SELECT id, name, motor_mechanism, transmitting_mechanism, working_tool,
		machine_value, lifespan_days, productive_power, hand_labour_per_unit,
		accumulated_wear, accumulated_depreciation, created_at
		FROM machines ORDER BY name ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []machinery.Machine
	for rows.Next() {
		mc, err := scanMachine(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, mc)
	}
	return out, rows.Err()
}

func (m *MySQL) UpdateMachine(ctx context.Context, id machinery.MachineID, u MachineUpdate) (machinery.Machine, error) {
	if u.IsEmpty() {
		return m.GetMachine(ctx, id)
	}
	cur, err := m.GetMachine(ctx, id)
	if err != nil {
		return machinery.Machine{}, err
	}
	if u.AccumulatedWear != nil {
		cur.AccumulatedWear = *u.AccumulatedWear
	}
	if u.AccumulatedDepreciation != nil {
		cur.AccumulatedDepreciation = *u.AccumulatedDepreciation
	}
	const q = `UPDATE machines SET accumulated_wear = ?, accumulated_depreciation = ? WHERE id = ?`
	res, err := m.db.ExecContext(ctx, q,
		int64(cur.AccumulatedWear.Value), int64(cur.AccumulatedDepreciation.Value), string(id),
	)
	if err != nil {
		return machinery.Machine{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return machinery.Machine{}, ErrNotFound
	}
	return cur, nil
}

func (m *MySQL) CreateFactory(ctx context.Context, f machinery.Factory) (machinery.Factory, error) {
	if f.ID.IsZero() {
		f.ID = machinery.NewFactoryID()
	}
	f.CreatedAt = m.now().UTC()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return machinery.Factory{}, err
	}
	defer func() { _ = tx.Rollback() }()

	// Resolve embedded vs. referenced machines. Embedded ones get inserted.
	resolved := make([]machinery.Machine, 0, len(f.Machines))
	for _, mc := range f.Machines {
		if !mc.ID.IsZero() {
			existing, err := m.getMachineTx(ctx, tx, mc.ID)
			if err == nil {
				resolved = append(resolved, existing)
				continue
			}
			if !errors.Is(err, ErrNotFound) {
				return machinery.Factory{}, err
			}
		}
		if err := mc.Validate(); err != nil {
			return machinery.Factory{}, err
		}
		if mc.ID.IsZero() {
			mc.ID = machinery.NewMachineID()
		}
		mc.CreatedAt = m.now().UTC()
		const insertM = `INSERT INTO machines
			(id, name, motor_mechanism, transmitting_mechanism, working_tool,
			 machine_value, lifespan_days, productive_power, hand_labour_per_unit,
			 accumulated_wear, accumulated_depreciation, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
		if _, err := tx.ExecContext(ctx, insertM,
			string(mc.ID), mc.Name, mc.MotorMechanism, mc.TransmittingMechanism, mc.WorkingTool,
			int64(mc.MachineValue), int64(mc.LifespanDays), int64(mc.ProductivePower),
			int64(mc.HandLabourPerUnit),
			int64(mc.AccumulatedWear.Value), int64(mc.AccumulatedDepreciation.Value),
			mc.CreatedAt,
		); err != nil {
			return machinery.Factory{}, err
		}
		resolved = append(resolved, mc)
	}
	f.Machines = resolved
	if err := f.Validate(); err != nil {
		return machinery.Factory{}, err
	}

	const insertF = `INSERT INTO factories
		(id, name, prime_mover_kind, prime_mover_horsepower, tick_count, intensity_factor, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	if _, err := tx.ExecContext(ctx, insertF,
		string(f.ID), f.Name, string(f.PrimeMover.Kind), f.PrimeMover.Horsepower,
		f.TickCount, float64(f.IntensityFactor), f.CreatedAt,
	); err != nil {
		return machinery.Factory{}, err
	}
	for idx, mc := range f.Machines {
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO factory_machines (factory_id, machine_id, ordinal) VALUES (?, ?, ?)`,
			string(f.ID), string(mc.ID), idx,
		); err != nil {
			return machinery.Factory{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return machinery.Factory{}, err
	}
	return f, nil
}

func (m *MySQL) GetFactory(ctx context.Context, id machinery.FactoryID) (machinery.Factory, error) {
	const q = `SELECT id, name, prime_mover_kind, prime_mover_horsepower, tick_count, intensity_factor, created_at
		FROM factories WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	f, err := scanFactory(row.Scan)
	if err != nil {
		return machinery.Factory{}, err
	}
	machines, err := m.machinesForFactory(ctx, id)
	if err != nil {
		return machinery.Factory{}, err
	}
	f.Machines = machines
	return f, nil
}

func (m *MySQL) ListFactories(ctx context.Context) ([]machinery.Factory, error) {
	const q = `SELECT id, name, prime_mover_kind, prime_mover_horsepower, tick_count, intensity_factor, created_at
		FROM factories ORDER BY name ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []machinery.Factory
	for rows.Next() {
		f, err := scanFactory(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range out {
		machines, err := m.machinesForFactory(ctx, out[i].ID)
		if err != nil {
			return nil, err
		}
		out[i].Machines = machines
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out, nil
}

func (m *MySQL) AdvanceTick(ctx context.Context, id machinery.FactoryID) (machinery.Factory, engine.Tick, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return machinery.Factory{}, engine.Tick{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const fq = `SELECT id, name, prime_mover_kind, prime_mover_horsepower, tick_count, intensity_factor, created_at
		FROM factories WHERE id = ? FOR UPDATE`
	row := tx.QueryRowContext(ctx, fq, string(id))
	f, err := scanFactory(row.Scan)
	if err != nil {
		return machinery.Factory{}, engine.Tick{}, err
	}
	rows, err := tx.QueryContext(ctx, `SELECT machine_id FROM factory_machines WHERE factory_id = ? ORDER BY ordinal ASC`, string(id))
	if err != nil {
		return machinery.Factory{}, engine.Tick{}, err
	}
	var ids []machinery.MachineID
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			rows.Close()
			return machinery.Factory{}, engine.Tick{}, err
		}
		ids = append(ids, machinery.MachineID(s))
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return machinery.Factory{}, engine.Tick{}, err
	}
	machines := make([]machinery.Machine, 0, len(ids))
	for _, mid := range ids {
		mc, err := m.getMachineTx(ctx, tx, mid)
		if err != nil {
			return machinery.Factory{}, engine.Tick{}, err
		}
		machines = append(machines, mc)
	}
	f.Machines = machines
	result := f.RunTick()
	// Accumulate per-machine wear, persist.
	for i, mc := range f.Machines {
		dwt := machinery.DailyWearAndTear(mc)
		mc.AccumulatedWear.Value += dwt
		f.Machines[i] = mc
		if _, err := tx.ExecContext(ctx,
			`UPDATE machines SET accumulated_wear = ? WHERE id = ?`,
			int64(mc.AccumulatedWear.Value), string(mc.ID),
		); err != nil {
			return machinery.Factory{}, engine.Tick{}, err
		}
	}
	f.TickCount++
	now := m.now().UTC()
	if _, err := tx.ExecContext(ctx,
		`UPDATE factories SET tick_count = ? WHERE id = ?`,
		f.TickCount, string(id),
	); err != nil {
		return machinery.Factory{}, engine.Tick{}, err
	}
	tick := engine.Tick{
		FactoryID:        string(id),
		Sequence:         f.TickCount,
		ValueTransferred: int64(result.ValueTransferred),
		UnitsProduced:    result.UnitsProduced,
		HandLabourSaved:  int64(result.HandLabourSaved),
		OccurredAt:       now,
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO factory_ticks (factory_id, sequence, value_transferred, units_produced, hand_labour_saved, occurred_at) VALUES (?, ?, ?, ?, ?, ?)`,
		tick.FactoryID, tick.Sequence, tick.ValueTransferred, tick.UnitsProduced, tick.HandLabourSaved, tick.OccurredAt,
	); err != nil {
		return machinery.Factory{}, engine.Tick{}, err
	}
	if err := tx.Commit(); err != nil {
		return machinery.Factory{}, engine.Tick{}, err
	}
	return f, tick, nil
}

func (m *MySQL) ListTicks(ctx context.Context, id machinery.FactoryID, limit int) ([]engine.Tick, error) {
	if limit <= 0 {
		limit = 100
	}
	const q = `SELECT factory_id, sequence, value_transferred, units_produced, hand_labour_saved, occurred_at
		FROM factory_ticks WHERE factory_id = ? ORDER BY sequence DESC LIMIT ?`
	rows, err := m.db.QueryContext(ctx, q, string(id), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []engine.Tick
	for rows.Next() {
		var t engine.Tick
		if err := rows.Scan(&t.FactoryID, &t.Sequence, &t.ValueTransferred, &t.UnitsProduced, &t.HandLabourSaved, &t.OccurredAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Reverse to ascending-sequence order so the UI can render left-to-right.
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, nil
}

// machinesForFactory returns the embedded machines for a given factory.
func (m *MySQL) machinesForFactory(ctx context.Context, id machinery.FactoryID) ([]machinery.Machine, error) {
	const q = `SELECT m.id, m.name, m.motor_mechanism, m.transmitting_mechanism, m.working_tool,
		m.machine_value, m.lifespan_days, m.productive_power, m.hand_labour_per_unit,
		m.accumulated_wear, m.accumulated_depreciation, m.created_at
		FROM machines m
		INNER JOIN factory_machines fm ON fm.machine_id = m.id
		WHERE fm.factory_id = ?
		ORDER BY fm.ordinal ASC`
	rows, err := m.db.QueryContext(ctx, q, string(id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []machinery.Machine
	for rows.Next() {
		mc, err := scanMachine(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, mc)
	}
	return out, rows.Err()
}

func (m *MySQL) getMachineTx(ctx context.Context, tx *sql.Tx, id machinery.MachineID) (machinery.Machine, error) {
	const q = `SELECT id, name, motor_mechanism, transmitting_mechanism, working_tool,
		machine_value, lifespan_days, productive_power, hand_labour_per_unit,
		accumulated_wear, accumulated_depreciation, created_at
		FROM machines WHERE id = ?`
	row := tx.QueryRowContext(ctx, q, string(id))
	return scanMachine(row.Scan)
}

func (m *MySQL) CreateGeneralLawScenario(ctx context.Context, s simulation.GeneralLawScenario) (simulation.GeneralLawScenario, error) {
	if s.ID.IsZero() {
		s.ID = simulation.NewGeneralLawScenarioID()
	}
	s.CreatedAt = m.now().UTC()
	const q = `INSERT INTO general_law_scenarios
		(id, name, constant_capital, variable_capital, surplus_rate, accumulation_rate,
		 productivity_growth, wage_pence, worker_supply, periods, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(s.ID), s.Name,
		int64(s.ConstantCapital), int64(s.VariableCapital),
		s.SurplusRate, s.AccumulationRate, s.ProductivityGrowth,
		s.WagePence, s.WorkerSupply, s.Periods,
		s.CreatedAt,
	)
	if err != nil {
		return simulation.GeneralLawScenario{}, err
	}
	return s, nil
}

func (m *MySQL) GetGeneralLawScenario(ctx context.Context, id simulation.GeneralLawScenarioID) (simulation.GeneralLawScenario, error) {
	const q = `SELECT id, name, constant_capital, variable_capital, surplus_rate,
		accumulation_rate, productivity_growth, wage_pence, worker_supply, periods, created_at
		FROM general_law_scenarios WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	return scanGeneralLawScenario(row.Scan)
}

func (m *MySQL) CreateHistoricalStage(ctx context.Context, h simulation.HistoricalStage) (simulation.HistoricalStage, error) {
	if err := h.Validate(); err != nil {
		return simulation.HistoricalStage{}, err
	}
	if h.ID.IsZero() {
		h.ID = simulation.NewHistoricalStageID()
	}
	h.CreatedAt = m.now().UTC()

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return simulation.HistoricalStage{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const insertStage = `INSERT INTO historical_stages (id, name, description, created_at) VALUES (?, ?, ?, ?)`
	if _, err := tx.ExecContext(ctx, insertStage,
		string(h.ID), h.Name, h.Description, h.CreatedAt,
	); err != nil {
		return simulation.HistoricalStage{}, err
	}
	const insertEpisode = `INSERT INTO primitive_accumulations
		(id, stage_id, ordinal, period, method, labourers_expropriated, capital_formed)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	for i, p := range h.PrimitiveAccumulations {
		episodeID := simulation.NewHistoricalStageID()
		if _, err := tx.ExecContext(ctx, insertEpisode,
			string(episodeID), string(h.ID), i,
			p.Period, p.Method, p.LabourersExpropriated, int64(p.CapitalFormed),
		); err != nil {
			return simulation.HistoricalStage{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return simulation.HistoricalStage{}, err
	}
	return h, nil
}

func (m *MySQL) GetHistoricalStage(ctx context.Context, id simulation.HistoricalStageID) (simulation.HistoricalStage, error) {
	const q = `SELECT id, name, description, created_at FROM historical_stages WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	h, err := scanHistoricalStage(row.Scan)
	if err != nil {
		return simulation.HistoricalStage{}, err
	}
	episodes, err := m.episodesForStage(ctx, id)
	if err != nil {
		return simulation.HistoricalStage{}, err
	}
	h.PrimitiveAccumulations = episodes
	return h, nil
}

func (m *MySQL) ListHistoricalStages(ctx context.Context) ([]simulation.HistoricalStage, error) {
	const q = `SELECT id, name, description, created_at FROM historical_stages ORDER BY name ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []simulation.HistoricalStage
	for rows.Next() {
		h, err := scanHistoricalStage(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range out {
		episodes, err := m.episodesForStage(ctx, out[i].ID)
		if err != nil {
			return nil, err
		}
		out[i].PrimitiveAccumulations = episodes
	}
	return out, nil
}

func (m *MySQL) episodesForStage(ctx context.Context, id simulation.HistoricalStageID) ([]simulation.PrimitiveAccumulation, error) {
	const q = `SELECT period, method, labourers_expropriated, capital_formed
		FROM primitive_accumulations WHERE stage_id = ? ORDER BY ordinal ASC`
	rows, err := m.db.QueryContext(ctx, q, string(id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []simulation.PrimitiveAccumulation
	for rows.Next() {
		var p simulation.PrimitiveAccumulation
		var capital int64
		if err := rows.Scan(&p.Period, &p.Method, &p.LabourersExpropriated, &capital); err != nil {
			return nil, err
		}
		p.CapitalFormed = simulation.Pence(capital)
		out = append(out, p)
	}
	return out, rows.Err()
}

func scanHistoricalStage(scan scanFn) (simulation.HistoricalStage, error) {
	var h simulation.HistoricalStage
	var id string
	err := scan(&id, &h.Name, &h.Description, &h.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return simulation.HistoricalStage{}, ErrNotFound
	}
	if err != nil {
		return simulation.HistoricalStage{}, err
	}
	h.ID = simulation.HistoricalStageID(id)
	return h, nil
}

func scanGeneralLawScenario(scan scanFn) (simulation.GeneralLawScenario, error) {
	var s simulation.GeneralLawScenario
	var id string
	var constant, variable int64
	err := scan(&id, &s.Name, &constant, &variable,
		&s.SurplusRate, &s.AccumulationRate, &s.ProductivityGrowth,
		&s.WagePence, &s.WorkerSupply, &s.Periods, &s.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return simulation.GeneralLawScenario{}, ErrNotFound
	}
	if err != nil {
		return simulation.GeneralLawScenario{}, err
	}
	s.ID = simulation.GeneralLawScenarioID(id)
	s.ConstantCapital = simulation.Pence(constant)
	s.VariableCapital = simulation.Pence(variable)
	return s, nil
}

type scanFn func(dest ...any) error

func scanMachine(scan scanFn) (machinery.Machine, error) {
	var mc machinery.Machine
	var id string
	var mv, lifespan, productive, handPerUnit, wear, depr int64
	err := scan(&id, &mc.Name, &mc.MotorMechanism, &mc.TransmittingMechanism, &mc.WorkingTool,
		&mv, &lifespan, &productive, &handPerUnit, &wear, &depr, &mc.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return machinery.Machine{}, ErrNotFound
	}
	if err != nil {
		return machinery.Machine{}, err
	}
	mc.ID = machinery.MachineID(id)
	mc.MachineValue = machinery.MachineValue(mv)
	mc.LifespanDays = machinery.LifespanDays(lifespan)
	mc.ProductivePower = machinery.ProductivePower(productive)
	mc.HandLabourPerUnit = machinery.LabourMinutes(handPerUnit)
	mc.AccumulatedWear.Value = machinery.LabourMinutes(wear)
	mc.AccumulatedDepreciation.Value = machinery.LabourMinutes(depr)
	return mc, nil
}

func scanFactory(scan scanFn) (machinery.Factory, error) {
	var f machinery.Factory
	var id, kind string
	var hp, tickCount int64
	var intensity float64
	err := scan(&id, &f.Name, &kind, &hp, &tickCount, &intensity, &f.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return machinery.Factory{}, ErrNotFound
	}
	if err != nil {
		return machinery.Factory{}, err
	}
	f.ID = machinery.FactoryID(id)
	f.PrimeMover = machinery.PrimeMover{Kind: machinery.PrimeMoverKind(kind), Horsepower: hp}
	f.TickCount = tickCount
	f.IntensityFactor = machinery.IntensityFactor(intensity)
	return f, nil
}
