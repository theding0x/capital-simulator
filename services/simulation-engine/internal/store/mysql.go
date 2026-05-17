package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"io/fs"
	"sort"
	"time"

	"github.com/go-sql-driver/mysql"

	pkgmysql "github.com/theding0x/capital-simulator/pkg/mysql"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/engine"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/machinery"
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/simulation"
)

// isDuplicateKey reports whether err is a MySQL duplicate-key error (1062),
// which the handler layer maps to a 409 Conflict via ErrAlreadyExists.
func isDuplicateKey(err error) bool {
	var me *mysql.MySQLError
	return errors.As(err, &me) && me.Number == 1062
}

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
		if isDuplicateKey(err) {
			return simulation.HistoricalStage{}, ErrAlreadyExists
		}
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

func (m *MySQL) CreateEnclosureEvent(ctx context.Context, e simulation.EnclosureEvent) (simulation.EnclosureEvent, error) {
	if err := e.Validate(); err != nil {
		return simulation.EnclosureEvent{}, err
	}
	if e.ID.IsZero() {
		e.ID = simulation.NewEnclosureEventID()
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = m.now()
	}
	const q = `INSERT INTO enclosure_events (id, period, acres_enclosed, population_displaced, beneficiary, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q, string(e.ID), e.Period, e.AcresEnclosed, e.PopulationDisplaced, e.Beneficiary, e.CreatedAt)
	if err != nil {
		return simulation.EnclosureEvent{}, err
	}
	return e, nil
}

func (m *MySQL) ListEnclosureEvents(ctx context.Context) ([]simulation.EnclosureEvent, error) {
	const q = `SELECT id, period, acres_enclosed, population_displaced, beneficiary, created_at FROM enclosure_events ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []simulation.EnclosureEvent
	for rows.Next() {
		var e simulation.EnclosureEvent
		var id string
		if err := rows.Scan(&id, &e.Period, &e.AcresEnclosed, &e.PopulationDisplaced, &e.Beneficiary, &e.CreatedAt); err != nil {
			return nil, err
		}
		e.ID = simulation.EnclosureEventID(id)
		out = append(out, e)
	}
	if out == nil {
		out = []simulation.EnclosureEvent{}
	}
	return out, rows.Err()
}

func (m *MySQL) CreateWageStatute(ctx context.Context, w simulation.WageStatute) (simulation.WageStatute, error) {
	if err := w.Validate(); err != nil {
		return simulation.WageStatute{}, err
	}
	if w.ID.IsZero() {
		w.ID = simulation.NewWageStatuteID()
	}
	if w.CreatedAt.IsZero() {
		w.CreatedAt = m.now()
	}
	const q = `INSERT INTO wage_statutes
		(id, historical_stage_id, period, jurisdiction, max_wage_pence, min_wage_pence, enforcement_penalty, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(w.ID), string(w.HistoricalStageID),
		w.Period, w.Jurisdiction, w.MaxWagePence, w.MinWagePence, w.EnforcementPenalty, w.CreatedAt)
	if err != nil {
		return simulation.WageStatute{}, err
	}
	return w, nil
}

func (m *MySQL) ListWageStatutesByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.WageStatute, error) {
	const q = `SELECT id, historical_stage_id, period, jurisdiction, max_wage_pence, min_wage_pence, enforcement_penalty, created_at
		FROM wage_statutes WHERE historical_stage_id = ? ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(stageID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]simulation.WageStatute, 0)
	for rows.Next() {
		var w simulation.WageStatute
		var id, sid string
		if err := rows.Scan(&id, &sid, &w.Period, &w.Jurisdiction, &w.MaxWagePence, &w.MinWagePence, &w.EnforcementPenalty, &w.CreatedAt); err != nil {
			return nil, err
		}
		w.ID = simulation.WageStatuteID(id)
		w.HistoricalStageID = simulation.HistoricalStageID(sid)
		out = append(out, w)
	}
	return out, rows.Err()
}

func (m *MySQL) CreateVagrancyLaw(ctx context.Context, v simulation.VagrancyLaw) (simulation.VagrancyLaw, error) {
	if err := v.Validate(); err != nil {
		return simulation.VagrancyLaw{}, err
	}
	if v.ID.IsZero() {
		v.ID = simulation.NewVagrancyLawID()
	}
	if v.CreatedAt.IsZero() {
		v.CreatedAt = m.now()
	}
	const q = `INSERT INTO vagrancy_laws
		(id, historical_stage_id, period, jurisdiction, punishment, target_population, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(v.ID), string(v.HistoricalStageID),
		v.Period, v.Jurisdiction, v.Punishment, v.TargetPopulation, v.CreatedAt)
	if err != nil {
		return simulation.VagrancyLaw{}, err
	}
	return v, nil
}

func (m *MySQL) ListVagrancyLawsByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.VagrancyLaw, error) {
	const q = `SELECT id, historical_stage_id, period, jurisdiction, punishment, target_population, created_at
		FROM vagrancy_laws WHERE historical_stage_id = ? ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(stageID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]simulation.VagrancyLaw, 0)
	for rows.Next() {
		var v simulation.VagrancyLaw
		var id, sid string
		if err := rows.Scan(&id, &sid, &v.Period, &v.Jurisdiction, &v.Punishment, &v.TargetPopulation, &v.CreatedAt); err != nil {
			return nil, err
		}
		v.ID = simulation.VagrancyLawID(id)
		v.HistoricalStageID = simulation.HistoricalStageID(sid)
		out = append(out, v)
	}
	return out, rows.Err()
}

func (m *MySQL) CreateFarmTenure(ctx context.Context, f simulation.FarmTenure) (simulation.FarmTenure, error) {
	if err := f.Validate(); err != nil {
		return simulation.FarmTenure{}, err
	}
	if f.ID.IsZero() {
		f.ID = simulation.NewFarmTenureID()
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = m.now()
	}
	const q = `INSERT INTO farm_tenures
		(id, historical_stage_id, form, lease_period_years, rent_pence, capital_advanced_pence, revenue_pence, wage_costs_pence, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(f.ID), string(f.HistoricalStageID),
		string(f.Form), f.LeasePeriodYears,
		int64(f.RentPence), int64(f.CapitalAdvancedPence),
		int64(f.RevenuePence), int64(f.WageCostsPence),
		f.CreatedAt)
	if err != nil {
		return simulation.FarmTenure{}, err
	}
	return f, nil
}

func (m *MySQL) ListFarmTenuresByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.FarmTenure, error) {
	const q = `SELECT id, historical_stage_id, form, lease_period_years, rent_pence, capital_advanced_pence, revenue_pence, wage_costs_pence, created_at
		FROM farm_tenures WHERE historical_stage_id = ? ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(stageID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]simulation.FarmTenure, 0)
	for rows.Next() {
		var f simulation.FarmTenure
		var id, sid, form string
		var rent, capital, revenue, wages int64
		if err := rows.Scan(&id, &sid, &form, &f.LeasePeriodYears, &rent, &capital, &revenue, &wages, &f.CreatedAt); err != nil {
			return nil, err
		}
		f.ID = simulation.FarmTenureID(id)
		f.HistoricalStageID = simulation.HistoricalStageID(sid)
		f.Form = simulation.TenantForm(form)
		f.RentPence = simulation.Pence(rent)
		f.CapitalAdvancedPence = simulation.Pence(capital)
		f.RevenuePence = simulation.Pence(revenue)
		f.WageCostsPence = simulation.Pence(wages)
		out = append(out, f)
	}
	return out, rows.Err()
}

func (m *MySQL) CreateDomesticIndustry(ctx context.Context, d simulation.DomesticIndustry) (simulation.DomesticIndustry, error) {
	if err := d.Validate(); err != nil {
		return simulation.DomesticIndustry{}, err
	}
	if d.ID.IsZero() {
		d.ID = simulation.NewDomesticIndustryID()
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = m.now()
	}
	const q = `INSERT INTO domestic_industries
		(id, historical_stage_id, name, households_engaged, annual_output_pence, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(d.ID), string(d.HistoricalStageID),
		d.Name, d.HouseholdsEngaged, int64(d.AnnualOutputPence), d.CreatedAt)
	if err != nil {
		return simulation.DomesticIndustry{}, err
	}
	return d, nil
}

func (m *MySQL) ListDomesticIndustriesByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.DomesticIndustry, error) {
	const q = `SELECT id, historical_stage_id, name, households_engaged, annual_output_pence, created_at
		FROM domestic_industries WHERE historical_stage_id = ? ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(stageID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]simulation.DomesticIndustry, 0)
	for rows.Next() {
		var d simulation.DomesticIndustry
		var id, sid string
		var outputPence int64
		if err := rows.Scan(&id, &sid, &d.Name, &d.HouseholdsEngaged, &outputPence, &d.CreatedAt); err != nil {
			return nil, err
		}
		d.ID = simulation.DomesticIndustryID(id)
		d.HistoricalStageID = simulation.HistoricalStageID(sid)
		d.AnnualOutputPence = simulation.Pence(outputPence)
		out = append(out, d)
	}
	return out, rows.Err()
}

func (m *MySQL) CreateCapitalOrigin(ctx context.Context, c simulation.CapitalOrigin) (simulation.CapitalOrigin, error) {
	if err := c.Validate(); err != nil {
		return simulation.CapitalOrigin{}, err
	}
	if c.ID.IsZero() {
		c.ID = simulation.NewCapitalOriginID()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now()
	}
	const q = `INSERT INTO capital_origins
		(id, historical_stage_id, source, amount_pence, period, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID), string(c.HistoricalStageID),
		c.Source, int64(c.AmountPence), c.Period, c.CreatedAt)
	if err != nil {
		return simulation.CapitalOrigin{}, err
	}
	return c, nil
}

func (m *MySQL) ListCapitalOriginsByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.CapitalOrigin, error) {
	const q = `SELECT id, historical_stage_id, source, amount_pence, period, created_at
		FROM capital_origins WHERE historical_stage_id = ? ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(stageID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]simulation.CapitalOrigin, 0)
	for rows.Next() {
		var c simulation.CapitalOrigin
		var id, sid string
		var amount int64
		if err := rows.Scan(&id, &sid, &c.Source, &amount, &c.Period, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.ID = simulation.CapitalOriginID(id)
		c.HistoricalStageID = simulation.HistoricalStageID(sid)
		c.AmountPence = simulation.Pence(amount)
		out = append(out, c)
	}
	return out, rows.Err()
}

func (m *MySQL) CreateColonialTransfer(ctx context.Context, t simulation.ColonialTransfer) (simulation.ColonialTransfer, error) {
	if err := t.Validate(); err != nil {
		return simulation.ColonialTransfer{}, err
	}
	if t.ID.IsZero() {
		t.ID = simulation.NewColonialTransferID()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = m.now()
	}
	const q = "INSERT INTO colonial_transfers (id, historical_stage_id, `from`, `to`, value_pence, method, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := m.db.ExecContext(ctx, q,
		string(t.ID), string(t.HistoricalStageID),
		t.From, t.To, int64(t.ValuePence), t.Method, t.CreatedAt)
	if err != nil {
		return simulation.ColonialTransfer{}, err
	}
	return t, nil
}

func (m *MySQL) ListColonialTransfersByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.ColonialTransfer, error) {
	const q = "SELECT id, historical_stage_id, `from`, `to`, value_pence, method, created_at FROM colonial_transfers WHERE historical_stage_id = ? ORDER BY created_at ASC"
	rows, err := m.db.QueryContext(ctx, q, string(stageID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]simulation.ColonialTransfer, 0)
	for rows.Next() {
		var t simulation.ColonialTransfer
		var id, sid string
		var value int64
		if err := rows.Scan(&id, &sid, &t.From, &t.To, &value, &t.Method, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.ID = simulation.ColonialTransferID(id)
		t.HistoricalStageID = simulation.HistoricalStageID(sid)
		t.ValuePence = simulation.Pence(value)
		out = append(out, t)
	}
	return out, rows.Err()
}

func (m *MySQL) CreateNationalDebt(ctx context.Context, d simulation.NationalDebt) (simulation.NationalDebt, error) {
	if err := d.Validate(); err != nil {
		return simulation.NationalDebt{}, err
	}
	if d.ID.IsZero() {
		d.ID = simulation.NewNationalDebtID()
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = m.now()
	}
	const q = `INSERT INTO national_debts
		(id, historical_stage_id, amount_pence, interest_rate_bps, creditor_class, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(d.ID), string(d.HistoricalStageID),
		int64(d.AmountPence), d.InterestRateBps, d.CreditorClass, d.CreatedAt)
	if err != nil {
		return simulation.NationalDebt{}, err
	}
	return d, nil
}

func (m *MySQL) ListNationalDebtsByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.NationalDebt, error) {
	const q = `SELECT id, historical_stage_id, amount_pence, interest_rate_bps, creditor_class, created_at
		FROM national_debts WHERE historical_stage_id = ? ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(stageID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]simulation.NationalDebt, 0)
	for rows.Next() {
		var d simulation.NationalDebt
		var id, sid string
		var amount int64
		if err := rows.Scan(&id, &sid, &amount, &d.InterestRateBps, &d.CreditorClass, &d.CreatedAt); err != nil {
			return nil, err
		}
		d.ID = simulation.NationalDebtID(id)
		d.HistoricalStageID = simulation.HistoricalStageID(sid)
		d.AmountPence = simulation.Pence(amount)
		out = append(out, d)
	}
	return out, rows.Err()
}

func (m *MySQL) ListProtectionSystemsByStage(ctx context.Context, stageID simulation.HistoricalStageID) ([]simulation.ProtectionSystem, error) {
	const q = `SELECT id, historical_stage_id, tariff_rate_bps, beneficiary, period_start, period_end, created_at
		FROM protection_systems WHERE historical_stage_id = ? ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(stageID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]simulation.ProtectionSystem, 0)
	for rows.Next() {
		var s simulation.ProtectionSystem
		var id, sid string
		if err := rows.Scan(&id, &sid, &s.TariffRateBps, &s.Beneficiary, &s.PeriodStart, &s.PeriodEnd, &s.CreatedAt); err != nil {
			return nil, err
		}
		s.ID = simulation.ProtectionSystemID(id)
		s.HistoricalStageID = simulation.HistoricalStageID(sid)
		out = append(out, s)
	}
	return out, rows.Err()
}

// CreateAccumulationTrajectory persists a Ch. 32 long-run centralisation
// trajectory header plus its CentralisationStep children in a single
// transaction. Steps are stored in the order supplied.
func (m *MySQL) CreateAccumulationTrajectory(ctx context.Context, t simulation.AccumulationTrajectory) (simulation.AccumulationTrajectory, error) {
	if err := t.Validate(); err != nil {
		return simulation.AccumulationTrajectory{}, err
	}
	if t.ID.IsZero() {
		t.ID = simulation.NewAccumulationTrajectoryID()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = m.now().UTC()
	}
	if t.Steps == nil {
		t.Steps = []simulation.CentralisationStep{}
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return simulation.AccumulationTrajectory{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const insertHeader = `INSERT INTO accumulation_trajectories
		(id, name, initial_firms, initial_capital_pence, final_firms, final_capital_pence, reserve_army_size, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := tx.ExecContext(ctx, insertHeader,
		string(t.ID), t.Name, t.InitialFirms, int64(t.InitialCapitalPence),
		t.FinalFirms, int64(t.FinalCapitalPence), t.ReserveArmySize, t.CreatedAt); err != nil {
		if isDuplicateKey(err) {
			return simulation.AccumulationTrajectory{}, ErrAlreadyExists
		}
		return simulation.AccumulationTrajectory{}, err
	}

	const insertStep = `INSERT INTO centralisation_steps
		(id, trajectory_id, step_index, firms_absorbed, capital_concentrated_pence, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := tx.PrepareContext(ctx, insertStep)
	if err != nil {
		return simulation.AccumulationTrajectory{}, err
	}
	defer stmt.Close()
	for _, s := range t.Steps {
		stepID := simulation.NewAccumulationTrajectoryID()
		if _, err := stmt.ExecContext(ctx,
			string(stepID), string(t.ID), s.StepIndex,
			s.FirmsAbsorbed, int64(s.CapitalConcentratedPence), t.CreatedAt); err != nil {
			return simulation.AccumulationTrajectory{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return simulation.AccumulationTrajectory{}, err
	}
	return t, nil
}

// GetAccumulationTrajectory returns the trajectory header and its steps,
// ordered by step_index ascending. Returns ErrNotFound if no header
// exists.
func (m *MySQL) GetAccumulationTrajectory(ctx context.Context, id simulation.AccumulationTrajectoryID) (simulation.AccumulationTrajectory, error) {
	const q = `SELECT id, name, initial_firms, initial_capital_pence, final_firms, final_capital_pence, reserve_army_size, created_at
		FROM accumulation_trajectories WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	var t simulation.AccumulationTrajectory
	var rawID string
	var initCap, finalCap int64
	if err := row.Scan(&rawID, &t.Name, &t.InitialFirms, &initCap, &t.FinalFirms, &finalCap, &t.ReserveArmySize, &t.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return simulation.AccumulationTrajectory{}, ErrNotFound
		}
		return simulation.AccumulationTrajectory{}, err
	}
	t.ID = simulation.AccumulationTrajectoryID(rawID)
	t.InitialCapitalPence = simulation.Pence(initCap)
	t.FinalCapitalPence = simulation.Pence(finalCap)

	steps, err := m.loadTrajectorySteps(ctx, t.ID)
	if err != nil {
		return simulation.AccumulationTrajectory{}, err
	}
	t.Steps = steps
	return t, nil
}

// ListAccumulationTrajectories returns headers in created_at ascending
// order, each with its steps populated.
func (m *MySQL) ListAccumulationTrajectories(ctx context.Context) ([]simulation.AccumulationTrajectory, error) {
	const q = `SELECT id, name, initial_firms, initial_capital_pence, final_firms, final_capital_pence, reserve_army_size, created_at
		FROM accumulation_trajectories ORDER BY created_at ASC, name ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]simulation.AccumulationTrajectory, 0)
	for rows.Next() {
		var t simulation.AccumulationTrajectory
		var rawID string
		var initCap, finalCap int64
		if err := rows.Scan(&rawID, &t.Name, &t.InitialFirms, &initCap, &t.FinalFirms, &finalCap, &t.ReserveArmySize, &t.CreatedAt); err != nil {
			return nil, err
		}
		t.ID = simulation.AccumulationTrajectoryID(rawID)
		t.InitialCapitalPence = simulation.Pence(initCap)
		t.FinalCapitalPence = simulation.Pence(finalCap)
		out = append(out, t)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for i := range out {
		steps, err := m.loadTrajectorySteps(ctx, out[i].ID)
		if err != nil {
			return nil, err
		}
		out[i].Steps = steps
	}
	return out, nil
}

func (m *MySQL) loadTrajectorySteps(ctx context.Context, id simulation.AccumulationTrajectoryID) ([]simulation.CentralisationStep, error) {
	const q = `SELECT step_index, firms_absorbed, capital_concentrated_pence
		FROM centralisation_steps WHERE trajectory_id = ? ORDER BY step_index ASC`
	rows, err := m.db.QueryContext(ctx, q, string(id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	steps := make([]simulation.CentralisationStep, 0)
	for rows.Next() {
		var s simulation.CentralisationStep
		var concentrated int64
		if err := rows.Scan(&s.StepIndex, &s.FirmsAbsorbed, &concentrated); err != nil {
			return nil, err
		}
		s.CapitalConcentratedPence = simulation.Pence(concentrated)
		steps = append(steps, s)
	}
	return steps, rows.Err()
}
