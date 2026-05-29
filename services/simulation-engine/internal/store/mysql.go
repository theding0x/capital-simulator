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
	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
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
		stepID := simulation.NewCentralisationStepID()
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

// CreateColonialLabourMarket persists a Ch. 33 colonial labour market.
// Colony name uniqueness is enforced by the UNIQUE index and mapped to
// ErrAlreadyExists.
func (m *MySQL) CreateColonialLabourMarket(ctx context.Context, market simulation.ColonialLabourMarket) (simulation.ColonialLabourMarket, error) {
	if err := market.Validate(); err != nil {
		return simulation.ColonialLabourMarket{}, err
	}
	if market.ID.IsZero() {
		market.ID = simulation.NewColonialLabourMarketID()
	}
	if market.CreatedAt.IsZero() {
		market.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO colonial_labour_markets
		(id, colony, free_labourers, annual_wage_pence, land_available,
		 wakefield_scheme_applied, independence_years, surplus_labour_extractable, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err := m.db.ExecContext(ctx, q,
		string(market.ID), market.Colony, market.FreeLabourers, int64(market.AnnualWagePence),
		market.LandAvailable, market.WakefieldSchemeApplied,
		market.IndependenceYears, market.SurplusLabourExtractable, market.CreatedAt); err != nil {
		if isDuplicateKey(err) {
			return simulation.ColonialLabourMarket{}, ErrAlreadyExists
		}
		return simulation.ColonialLabourMarket{}, err
	}
	return market, nil
}

// GetColonialLabourMarket returns the market or ErrNotFound.
func (m *MySQL) GetColonialLabourMarket(ctx context.Context, id simulation.ColonialLabourMarketID) (simulation.ColonialLabourMarket, error) {
	const q = `SELECT id, colony, free_labourers, annual_wage_pence, land_available,
			wakefield_scheme_applied, independence_years, surplus_labour_extractable, created_at
		FROM colonial_labour_markets WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	market, err := scanColonialLabourMarket(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return simulation.ColonialLabourMarket{}, ErrNotFound
		}
		return simulation.ColonialLabourMarket{}, err
	}
	return market, nil
}

// ListColonialLabourMarkets returns markets ordered by created_at, colony.
func (m *MySQL) ListColonialLabourMarkets(ctx context.Context) ([]simulation.ColonialLabourMarket, error) {
	const q = `SELECT id, colony, free_labourers, annual_wage_pence, land_available,
			wakefield_scheme_applied, independence_years, surplus_labour_extractable, created_at
		FROM colonial_labour_markets ORDER BY created_at ASC, colony ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]simulation.ColonialLabourMarket, 0)
	for rows.Next() {
		market, err := scanColonialLabourMarket(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, market)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.SliceStable(out, func(i, j int) bool {
		if !out[i].CreatedAt.Equal(out[j].CreatedAt) {
			return out[i].CreatedAt.Before(out[j].CreatedAt)
		}
		return out[i].Colony < out[j].Colony
	})
	return out, nil
}

// UpdateColonialLabourMarket applies the partial regulation update inside
// a transaction so the read-then-write is consistent.
func (m *MySQL) UpdateColonialLabourMarket(ctx context.Context, id simulation.ColonialLabourMarketID, u ColonialLabourMarketUpdate) (simulation.ColonialLabourMarket, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return simulation.ColonialLabourMarket{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const sel = `SELECT id, colony, free_labourers, annual_wage_pence, land_available,
			wakefield_scheme_applied, independence_years, surplus_labour_extractable, created_at
		FROM colonial_labour_markets WHERE id = ? FOR UPDATE`
	row := tx.QueryRowContext(ctx, sel, string(id))
	cur, err := scanColonialLabourMarket(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return simulation.ColonialLabourMarket{}, ErrNotFound
		}
		return simulation.ColonialLabourMarket{}, err
	}
	if u.WakefieldSchemeApplied != nil {
		cur.WakefieldSchemeApplied = *u.WakefieldSchemeApplied
	}
	if u.IndependenceYears != nil {
		cur.IndependenceYears = *u.IndependenceYears
	}
	if u.SurplusLabourExtractable != nil {
		cur.SurplusLabourExtractable = *u.SurplusLabourExtractable
	}
	const upd = `UPDATE colonial_labour_markets
		SET wakefield_scheme_applied = ?, independence_years = ?, surplus_labour_extractable = ?
		WHERE id = ?`
	if _, err := tx.ExecContext(ctx, upd,
		cur.WakefieldSchemeApplied, cur.IndependenceYears, cur.SurplusLabourExtractable,
		string(cur.ID)); err != nil {
		return simulation.ColonialLabourMarket{}, err
	}
	if err := tx.Commit(); err != nil {
		return simulation.ColonialLabourMarket{}, err
	}
	return cur, nil
}

// RegulateColonialLabourMarket reads the market under SELECT ... FOR
// UPDATE, runs simulation.ColonialLabourRegulation against the locked
// row, and writes the regulated state back inside the same
// transaction. Two concurrent /regulate callers cannot both compute
// against the same baseline and clobber each other.
func (m *MySQL) RegulateColonialLabourMarket(ctx context.Context, id simulation.ColonialLabourMarketID, scheme simulation.SystematicColonisation) (simulation.ColonialLabourMarket, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return simulation.ColonialLabourMarket{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const sel = `SELECT id, colony, free_labourers, annual_wage_pence, land_available,
			wakefield_scheme_applied, independence_years, surplus_labour_extractable, created_at
		FROM colonial_labour_markets WHERE id = ? FOR UPDATE`
	row := tx.QueryRowContext(ctx, sel, string(id))
	cur, err := scanColonialLabourMarket(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return simulation.ColonialLabourMarket{}, ErrNotFound
		}
		return simulation.ColonialLabourMarket{}, err
	}
	regulated := simulation.ColonialLabourRegulation(cur, scheme)
	const upd = `UPDATE colonial_labour_markets
		SET wakefield_scheme_applied = ?, independence_years = ?, surplus_labour_extractable = ?
		WHERE id = ?`
	if _, err := tx.ExecContext(ctx, upd,
		regulated.WakefieldSchemeApplied, regulated.IndependenceYears, regulated.SurplusLabourExtractable,
		string(regulated.ID)); err != nil {
		return simulation.ColonialLabourMarket{}, err
	}
	if err := tx.Commit(); err != nil {
		return simulation.ColonialLabourMarket{}, err
	}
	return regulated, nil
}

// rowScanner is the minimal interface satisfied by *sql.Row and *sql.Rows
// so scanColonialLabourMarket can serve both Get and List paths.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanColonialLabourMarket(r rowScanner) (simulation.ColonialLabourMarket, error) {
	var (
		market simulation.ColonialLabourMarket
		rawID  string
		annual int64
	)
	if err := r.Scan(
		&rawID, &market.Colony, &market.FreeLabourers, &annual, &market.LandAvailable,
		&market.WakefieldSchemeApplied, &market.IndependenceYears, &market.SurplusLabourExtractable,
		&market.CreatedAt,
	); err != nil {
		return simulation.ColonialLabourMarket{}, err
	}
	market.ID = simulation.ColonialLabourMarketID(rawID)
	market.AnnualWagePence = simulation.Pence(annual)
	return market, nil
}

// --- Vol. II Ch. 2 — ProductiveCircuitStore ---

func (m *MySQL) CreateProductiveCircuit(ctx context.Context, pc circulation.ProductiveCircuit) (circulation.ProductiveCircuit, error) {
	if err := pc.Validate(); err != nil {
		return circulation.ProductiveCircuit{}, err
	}
	if pc.ID.IsZero() {
		pc.ID = circulation.NewProductiveCircuitID()
	}
	pc.LatentMoneyCapital.ProductiveCircuitID = pc.ID
	pc.LatentMoneyCapital.Threshold = pc.MinCapitalisationIncrement
	pc.ReserveFund.ProductiveCircuitID = pc.ID
	pc.CreatedAt = m.now().UTC()

	const q = `INSERT INTO productive_circuits
		(id, agent_id, constant_pence, variable_pence, mode, min_capitalisation_increment_pence,
		 latent_accumulated, reserve_balance, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(pc.ID), pc.AgentID, int64(pc.ConstantPence), int64(pc.VariablePence),
		string(pc.Mode), int64(pc.MinCapitalisationIncrement),
		int64(0), int64(0), pc.CreatedAt,
	)
	if err != nil {
		if isDuplicateKey(err) {
			return circulation.ProductiveCircuit{}, ErrAlreadyExists
		}
		return circulation.ProductiveCircuit{}, err
	}
	if pc.RevenueCircuits == nil {
		pc.RevenueCircuits = []circulation.RevenueCircuit{}
	}
	if pc.CapitalisationSteps == nil {
		pc.CapitalisationSteps = []circulation.CapitalisationStep{}
	}
	if pc.ReserveDraws == nil {
		pc.ReserveDraws = []circulation.ReserveDraw{}
	}
	return pc, nil
}

func (m *MySQL) GetProductiveCircuit(ctx context.Context, id circulation.ProductiveCircuitID) (circulation.ProductiveCircuit, error) {
	const q = `SELECT id, agent_id, constant_pence, variable_pence, mode,
		min_capitalisation_increment_pence, latent_accumulated, reserve_balance, created_at
		FROM productive_circuits WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	pc, err := scanProductiveCircuit(row)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return circulation.ProductiveCircuit{}, ErrNotFound
		}
		return circulation.ProductiveCircuit{}, err
	}

	rc, err := m.listRevenueCircuits(ctx, id)
	if err != nil {
		return circulation.ProductiveCircuit{}, err
	}
	pc.RevenueCircuits = rc

	steps, err := m.listCapitalisationSteps(ctx, id)
	if err != nil {
		return circulation.ProductiveCircuit{}, err
	}
	pc.CapitalisationSteps = steps

	draws, err := m.listReserveDraws(ctx, id)
	if err != nil {
		return circulation.ProductiveCircuit{}, err
	}
	pc.ReserveDraws = draws
	return pc, nil
}

func (m *MySQL) ListProductiveCircuits(ctx context.Context) ([]circulation.ProductiveCircuit, error) {
	const q = `SELECT id, agent_id, constant_pence, variable_pence, mode,
		min_capitalisation_increment_pence, latent_accumulated, reserve_balance, created_at
		FROM productive_circuits ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []circulation.ProductiveCircuit
	for rows.Next() {
		pc, err := scanProductiveCircuit(rows)
		if err != nil {
			return nil, err
		}
		rc, _ := m.listRevenueCircuits(ctx, pc.ID)
		pc.RevenueCircuits = rc
		steps, _ := m.listCapitalisationSteps(ctx, pc.ID)
		pc.CapitalisationSteps = steps
		draws, _ := m.listReserveDraws(ctx, pc.ID)
		pc.ReserveDraws = draws
		out = append(out, pc)
	}
	if out == nil {
		out = []circulation.ProductiveCircuit{}
	}
	return out, rows.Err()
}

func (m *MySQL) RecordRevenue(ctx context.Context, id circulation.ProductiveCircuitID, rc circulation.RevenueCircuit) (circulation.RevenueCircuit, error) {
	if _, err := m.GetProductiveCircuit(ctx, id); err != nil {
		return circulation.RevenueCircuit{}, err
	}
	rc.ProductiveCircuitID = id
	if rc.SpentAt.IsZero() {
		rc.SpentAt = m.now().UTC()
	}
	const q = `INSERT INTO revenue_circuits (productive_circuit_id, amount, spent_at) VALUES (?, ?, ?)`
	if _, err := m.db.ExecContext(ctx, q, string(id), int64(rc.Amount), rc.SpentAt); err != nil {
		return circulation.RevenueCircuit{}, err
	}
	return rc, nil
}

func (m *MySQL) Accumulate(ctx context.Context, id circulation.ProductiveCircuitID, amount circulation.Pence) (circulation.LatentMoneyCapital, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return circulation.LatentMoneyCapital{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const sel = `SELECT latent_accumulated, min_capitalisation_increment_pence FROM productive_circuits WHERE id = ? FOR UPDATE`
	var latentAcc, threshold int64
	if err := tx.QueryRowContext(ctx, sel, string(id)).Scan(&latentAcc, &threshold); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return circulation.LatentMoneyCapital{}, ErrNotFound
		}
		return circulation.LatentMoneyCapital{}, err
	}
	latentAcc += int64(amount)
	if _, err := tx.ExecContext(ctx, `UPDATE productive_circuits SET latent_accumulated = ? WHERE id = ?`, latentAcc, string(id)); err != nil {
		return circulation.LatentMoneyCapital{}, err
	}
	if err := tx.Commit(); err != nil {
		return circulation.LatentMoneyCapital{}, err
	}
	return circulation.LatentMoneyCapital{
		ProductiveCircuitID: id,
		Accumulated:         circulation.Pence(latentAcc),
		Threshold:           circulation.Pence(threshold),
	}, nil
}

func (m *MySQL) Capitalise(ctx context.Context, id circulation.ProductiveCircuitID, amount, dc, dv circulation.Pence) (circulation.CapitalisationStep, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return circulation.CapitalisationStep{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const sel = `SELECT constant_pence, variable_pence, latent_accumulated, min_capitalisation_increment_pence
		FROM productive_circuits WHERE id = ? FOR UPDATE`
	var (
		constant, variable, latentAcc, minIncrement int64
	)
	if err := tx.QueryRowContext(ctx, sel, string(id)).Scan(&constant, &variable, &latentAcc, &minIncrement); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return circulation.CapitalisationStep{}, ErrNotFound
		}
		return circulation.CapitalisationStep{}, err
	}
	if int64(amount) > 0 && int64(amount) < minIncrement {
		return circulation.CapitalisationStep{}, circulation.ErrIndivisibleProductionElement
	}
	now := m.now().UTC()
	const ins = `INSERT INTO capitalisation_steps
		(productive_circuit_id, amount_injected, delta_constant_pence, delta_variable_pence, occurred_at)
		VALUES (?, ?, ?, ?, ?)`
	if _, err := tx.ExecContext(ctx, ins, string(id), int64(amount), int64(dc), int64(dv), now); err != nil {
		return circulation.CapitalisationStep{}, err
	}
	const upd = `UPDATE productive_circuits
		SET constant_pence = ?, variable_pence = ?, latent_accumulated = ? WHERE id = ?`
	if _, err := tx.ExecContext(ctx, upd, constant+int64(dc), variable+int64(dv), latentAcc-int64(amount), string(id)); err != nil {
		return circulation.CapitalisationStep{}, err
	}
	if err := tx.Commit(); err != nil {
		return circulation.CapitalisationStep{}, err
	}
	return circulation.CapitalisationStep{
		ProductiveCircuitID: id,
		AmountInjected:      amount,
		DeltaConstantPence:  dc,
		DeltaVariablePence:  dv,
		OccurredAt:          now,
	}, nil
}

func (m *MySQL) DepositReserve(ctx context.Context, id circulation.ProductiveCircuitID, amount circulation.Pence) (circulation.ReserveFund, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return circulation.ReserveFund{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const sel = `SELECT reserve_balance FROM productive_circuits WHERE id = ? FOR UPDATE`
	var balance int64
	if err := tx.QueryRowContext(ctx, sel, string(id)).Scan(&balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return circulation.ReserveFund{}, ErrNotFound
		}
		return circulation.ReserveFund{}, err
	}
	balance += int64(amount)
	if _, err := tx.ExecContext(ctx, `UPDATE productive_circuits SET reserve_balance = ? WHERE id = ?`, balance, string(id)); err != nil {
		return circulation.ReserveFund{}, err
	}
	if err := tx.Commit(); err != nil {
		return circulation.ReserveFund{}, err
	}
	return circulation.ReserveFund{ProductiveCircuitID: id, Balance: circulation.Pence(balance)}, nil
}

func (m *MySQL) WithdrawReserve(ctx context.Context, id circulation.ProductiveCircuitID, amount circulation.Pence, reason circulation.ReserveDrawReason) (circulation.ReserveDraw, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return circulation.ReserveDraw{}, err
	}
	defer func() { _ = tx.Rollback() }()

	const sel = `SELECT reserve_balance FROM productive_circuits WHERE id = ? FOR UPDATE`
	var balance int64
	if err := tx.QueryRowContext(ctx, sel, string(id)).Scan(&balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return circulation.ReserveDraw{}, ErrNotFound
		}
		return circulation.ReserveDraw{}, err
	}
	if int64(amount) > balance {
		return circulation.ReserveDraw{}, circulation.ErrInsufficientReserve
	}
	now := m.now().UTC()
	const ins = `INSERT INTO reserve_draws (productive_circuit_id, drawn_pence, reason, occurred_at) VALUES (?, ?, ?, ?)`
	if _, err := tx.ExecContext(ctx, ins, string(id), int64(amount), string(reason), now); err != nil {
		return circulation.ReserveDraw{}, err
	}
	balance -= int64(amount)
	if _, err := tx.ExecContext(ctx, `UPDATE productive_circuits SET reserve_balance = ? WHERE id = ?`, balance, string(id)); err != nil {
		return circulation.ReserveDraw{}, err
	}
	if err := tx.Commit(); err != nil {
		return circulation.ReserveDraw{}, err
	}
	return circulation.ReserveDraw{
		ProductiveCircuitID: id,
		DrawnPence:          amount,
		Reason:              reason,
		OccurredAt:          now,
	}, nil
}

func (m *MySQL) listRevenueCircuits(ctx context.Context, id circulation.ProductiveCircuitID) ([]circulation.RevenueCircuit, error) {
	const q = `SELECT productive_circuit_id, amount, spent_at FROM revenue_circuits WHERE productive_circuit_id = ? ORDER BY spent_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(id))
	if err != nil {
		return []circulation.RevenueCircuit{}, err
	}
	defer rows.Close()
	var out []circulation.RevenueCircuit
	for rows.Next() {
		var rc circulation.RevenueCircuit
		var rawID string
		var amt int64
		if err := rows.Scan(&rawID, &amt, &rc.SpentAt); err != nil {
			return nil, err
		}
		rc.ProductiveCircuitID = circulation.ProductiveCircuitID(rawID)
		rc.Amount = circulation.Pence(amt)
		out = append(out, rc)
	}
	if out == nil {
		out = []circulation.RevenueCircuit{}
	}
	return out, rows.Err()
}

func (m *MySQL) listCapitalisationSteps(ctx context.Context, id circulation.ProductiveCircuitID) ([]circulation.CapitalisationStep, error) {
	const q = `SELECT productive_circuit_id, amount_injected, delta_constant_pence, delta_variable_pence, occurred_at
		FROM capitalisation_steps WHERE productive_circuit_id = ? ORDER BY occurred_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(id))
	if err != nil {
		return []circulation.CapitalisationStep{}, err
	}
	defer rows.Close()
	var out []circulation.CapitalisationStep
	for rows.Next() {
		var step circulation.CapitalisationStep
		var rawID string
		var injected, dc, dv int64
		if err := rows.Scan(&rawID, &injected, &dc, &dv, &step.OccurredAt); err != nil {
			return nil, err
		}
		step.ProductiveCircuitID = circulation.ProductiveCircuitID(rawID)
		step.AmountInjected = circulation.Pence(injected)
		step.DeltaConstantPence = circulation.Pence(dc)
		step.DeltaVariablePence = circulation.Pence(dv)
		out = append(out, step)
	}
	if out == nil {
		out = []circulation.CapitalisationStep{}
	}
	return out, rows.Err()
}

func (m *MySQL) listReserveDraws(ctx context.Context, id circulation.ProductiveCircuitID) ([]circulation.ReserveDraw, error) {
	const q = `SELECT productive_circuit_id, drawn_pence, reason, occurred_at
		FROM reserve_draws WHERE productive_circuit_id = ? ORDER BY occurred_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(id))
	if err != nil {
		return []circulation.ReserveDraw{}, err
	}
	defer rows.Close()
	var out []circulation.ReserveDraw
	for rows.Next() {
		var draw circulation.ReserveDraw
		var rawID, reason string
		var drawn int64
		if err := rows.Scan(&rawID, &drawn, &reason, &draw.OccurredAt); err != nil {
			return nil, err
		}
		draw.ProductiveCircuitID = circulation.ProductiveCircuitID(rawID)
		draw.DrawnPence = circulation.Pence(drawn)
		draw.Reason = circulation.ReserveDrawReason(reason)
		out = append(out, draw)
	}
	if out == nil {
		out = []circulation.ReserveDraw{}
	}
	return out, rows.Err()
}

// Vol. II Ch. 3 — CommodityCircuitStore implementation.

func (m *MySQL) CreateCommodityCircuit(ctx context.Context, cc circulation.CommodityCircuit) (circulation.CommodityCircuit, error) {
	if err := cc.Validate(); err != nil {
		return circulation.CommodityCircuit{}, err
	}
	if cc.ID.IsZero() {
		cc.ID = circulation.NewCommodityCircuitID()
	}
	cc.CreatedAt = m.now().UTC()
	const q = `INSERT INTO commodity_circuits
		(id, agent_id, constant_pence, variable_pence, surplus_pence, pounds_total,
		 mode, is_first_investment, social_capital_lens, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(cc.ID), cc.AgentID,
		int64(cc.Initial.ConstantPence), int64(cc.Initial.VariablePence),
		int64(cc.Initial.SurplusPence), cc.Initial.PoundsTotal,
		string(cc.Mode), cc.IsFirstInvestment, bool(cc.SocialCapitalLens),
		cc.CreatedAt,
	)
	if err != nil {
		if isDuplicateKey(err) {
			return circulation.CommodityCircuit{}, ErrAlreadyExists
		}
		return circulation.CommodityCircuit{}, err
	}
	if cc.PartialSales == nil {
		cc.PartialSales = []circulation.SuccessivePartialSale{}
	}
	if cc.MPSources == nil {
		cc.MPSources = []circulation.MeansOfProductionSource{}
	}
	return cc, nil
}

func (m *MySQL) GetCommodityCircuit(ctx context.Context, id circulation.CommodityCircuitID) (circulation.CommodityCircuit, error) {
	const q = `SELECT id, agent_id, constant_pence, variable_pence, surplus_pence, pounds_total,
		mode, is_first_investment, social_capital_lens,
		terminal_constant, terminal_variable, terminal_surplus, terminal_capitalised,
		terminal_pounds, terminal_closed_at, created_at
		FROM commodity_circuits WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	cc, err := scanCommodityCircuit(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return circulation.CommodityCircuit{}, ErrNotFound
		}
		return circulation.CommodityCircuit{}, err
	}
	sales, err := m.listPartialSales(ctx, id)
	if err != nil {
		return circulation.CommodityCircuit{}, err
	}
	cc.PartialSales = sales
	sources, err := m.listMPSources(ctx, id)
	if err != nil {
		return circulation.CommodityCircuit{}, err
	}
	cc.MPSources = sources
	return cc, nil
}

func (m *MySQL) ListCommodityCircuits(ctx context.Context, agentID string) ([]circulation.CommodityCircuit, error) {
	q := `SELECT id, agent_id, constant_pence, variable_pence, surplus_pence, pounds_total,
		mode, is_first_investment, social_capital_lens,
		terminal_constant, terminal_variable, terminal_surplus, terminal_capitalised,
		terminal_pounds, terminal_closed_at, created_at
		FROM commodity_circuits`
	args := []any{}
	if agentID != "" {
		q += ` WHERE agent_id = ?`
		args = append(args, agentID)
	}
	q += ` ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []circulation.CommodityCircuit
	for rows.Next() {
		cc, err := scanCommodityCircuit(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, cc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if out == nil {
		return []circulation.CommodityCircuit{}, nil
	}
	for i := range out {
		sales, err := m.listPartialSales(ctx, out[i].ID)
		if err != nil {
			return nil, err
		}
		out[i].PartialSales = sales
		sources, err := m.listMPSources(ctx, out[i].ID)
		if err != nil {
			return nil, err
		}
		out[i].MPSources = sources
	}
	return out, nil
}

func (m *MySQL) RecordPartialSale(ctx context.Context, id circulation.CommodityCircuitID, sale circulation.SuccessivePartialSale) (circulation.SuccessivePartialSale, error) {
	if _, err := m.GetCommodityCircuit(ctx, id); err != nil {
		return circulation.SuccessivePartialSale{}, err
	}
	sale.CommodityCircuitID = id
	if sale.SoldAt.IsZero() {
		sale.SoldAt = m.now().UTC()
	}
	const q = `INSERT INTO successive_partial_sales
		(commodity_circuit_id, quantity, realised_pence, constant_share, variable_share, surplus_share, sold_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(id), sale.Quantity, int64(sale.RealisedPence),
		int64(sale.Decomposition.ConstantShare), int64(sale.Decomposition.VariableShare),
		int64(sale.Decomposition.SurplusShare), sale.SoldAt,
	)
	if err != nil {
		return circulation.SuccessivePartialSale{}, err
	}
	return sale, nil
}

func (m *MySQL) LinkMPSource(ctx context.Context, id circulation.CommodityCircuitID, source circulation.MeansOfProductionSource) (circulation.MeansOfProductionSource, error) {
	if err := source.Validate(); err != nil {
		return circulation.MeansOfProductionSource{}, err
	}
	if _, err := m.GetCommodityCircuit(ctx, id); err != nil {
		return circulation.MeansOfProductionSource{}, err
	}
	source.CommodityCircuitID = id
	const q = `INSERT INTO mp_sources
		(commodity_circuit_id, source_commodity_circuit_id, source_kind)
		VALUES (?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q, string(id), source.SourceCommodityCircuitID, string(source.SourceKind))
	if err != nil {
		return circulation.MeansOfProductionSource{}, err
	}
	return source, nil
}

func (m *MySQL) CloseCommodityCircuit(ctx context.Context, id circulation.CommodityCircuitID, aug circulation.CommodityAugmented) (circulation.CommodityCircuit, error) {
	aug.CommodityCircuitID = id
	if aug.ClosedAt.IsZero() {
		aug.ClosedAt = m.now().UTC()
	}
	const q = `UPDATE commodity_circuits
		SET terminal_constant = ?, terminal_variable = ?, terminal_surplus = ?,
		    terminal_capitalised = ?, terminal_pounds = ?, terminal_closed_at = ?
		WHERE id = ?`
	res, err := m.db.ExecContext(ctx, q,
		int64(aug.ConstantPence), int64(aug.VariablePence), int64(aug.SurplusPence),
		int64(aug.CapitalisedPence), aug.PoundsTotal, aug.ClosedAt,
		string(id),
	)
	if err != nil {
		return circulation.CommodityCircuit{}, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return circulation.CommodityCircuit{}, err
	}
	if n == 0 {
		return circulation.CommodityCircuit{}, ErrNotFound
	}
	return m.GetCommodityCircuit(ctx, id)
}

func (m *MySQL) listPartialSales(ctx context.Context, id circulation.CommodityCircuitID) ([]circulation.SuccessivePartialSale, error) {
	const q = `SELECT commodity_circuit_id, quantity, realised_pence,
		constant_share, variable_share, surplus_share, sold_at
		FROM successive_partial_sales WHERE commodity_circuit_id = ? ORDER BY sold_at ASC`
	rows, err := m.db.QueryContext(ctx, q, string(id))
	if err != nil {
		return []circulation.SuccessivePartialSale{}, err
	}
	defer rows.Close()
	var out []circulation.SuccessivePartialSale
	for rows.Next() {
		var sale circulation.SuccessivePartialSale
		var rawID string
		var qty, realised, cs, vs, ss int64
		if err := rows.Scan(&rawID, &qty, &realised, &cs, &vs, &ss, &sale.SoldAt); err != nil {
			return nil, err
		}
		sale.CommodityCircuitID = circulation.CommodityCircuitID(rawID)
		sale.Quantity = qty
		sale.RealisedPence = circulation.Pence(realised)
		sale.Decomposition = circulation.ValueComposition{
			ConstantShare: circulation.Pence(cs),
			VariableShare: circulation.Pence(vs),
			SurplusShare:  circulation.Pence(ss),
		}
		out = append(out, sale)
	}
	if out == nil {
		out = []circulation.SuccessivePartialSale{}
	}
	return out, rows.Err()
}

func (m *MySQL) listMPSources(ctx context.Context, id circulation.CommodityCircuitID) ([]circulation.MeansOfProductionSource, error) {
	const q = `SELECT commodity_circuit_id, source_commodity_circuit_id, source_kind
		FROM mp_sources WHERE commodity_circuit_id = ? ORDER BY id ASC`
	rows, err := m.db.QueryContext(ctx, q, string(id))
	if err != nil {
		return []circulation.MeansOfProductionSource{}, err
	}
	defer rows.Close()
	var out []circulation.MeansOfProductionSource
	for rows.Next() {
		var src circulation.MeansOfProductionSource
		var rawID, sourceID, kind string
		if err := rows.Scan(&rawID, &sourceID, &kind); err != nil {
			return nil, err
		}
		src.CommodityCircuitID = circulation.CommodityCircuitID(rawID)
		src.SourceCommodityCircuitID = sourceID
		src.SourceKind = circulation.MeansOfProductionSourceKind(kind)
		out = append(out, src)
	}
	if out == nil {
		out = []circulation.MeansOfProductionSource{}
	}
	return out, rows.Err()
}

func scanCommodityCircuit(scan func(...any) error) (circulation.CommodityCircuit, error) {
	var (
		cc             circulation.CommodityCircuit
		rawID, agentID string
		constant, variable, surplus, pounds int64
		mode           string
		isFirst, social bool
		tConst, tVar, tSurp, tCap sql.NullInt64
		tPounds        sql.NullInt64
		tClosedAt      sql.NullTime
	)
	if err := scan(
		&rawID, &agentID, &constant, &variable, &surplus, &pounds,
		&mode, &isFirst, &social,
		&tConst, &tVar, &tSurp, &tCap, &tPounds, &tClosedAt,
		&cc.CreatedAt,
	); err != nil {
		return circulation.CommodityCircuit{}, err
	}
	cc.ID = circulation.CommodityCircuitID(rawID)
	cc.AgentID = agentID
	cc.Initial = circulation.OpeningCommodityCapital{
		ConstantPence: circulation.Pence(constant),
		VariablePence: circulation.Pence(variable),
		SurplusPence:  circulation.Pence(surplus),
		PoundsTotal:   pounds,
	}
	cc.Mode = circulation.ReproductionMode(mode)
	cc.IsFirstInvestment = isFirst
	cc.SocialCapitalLens = circulation.SocialCapitalLens(social)
	if tConst.Valid {
		aug := &circulation.CommodityAugmented{
			CommodityCircuitID: cc.ID,
			ConstantPence:      circulation.Pence(tConst.Int64),
			VariablePence:      circulation.Pence(tVar.Int64),
			SurplusPence:       circulation.Pence(tSurp.Int64),
			CapitalisedPence:   circulation.Pence(tCap.Int64),
		}
		if tPounds.Valid {
			aug.PoundsTotal = tPounds.Int64
		}
		if tClosedAt.Valid {
			aug.ClosedAt = tClosedAt.Time
		}
		cc.Terminal = aug
	}
	return cc, nil
}

func scanProductiveCircuit(r rowScanner) (circulation.ProductiveCircuit, error) {
	var (
		pc      circulation.ProductiveCircuit
		rawID   string
		agentID sql.NullString
		constant, variable, minIncrement, latentAcc, reserveBal int64
		mode    string
	)
	if err := r.Scan(&rawID, &agentID, &constant, &variable, &mode, &minIncrement, &latentAcc, &reserveBal, &pc.CreatedAt); err != nil {
		return circulation.ProductiveCircuit{}, err
	}
	pc.ID = circulation.ProductiveCircuitID(rawID)
	if agentID.Valid {
		pc.AgentID = agentID.String
	}
	pc.ConstantPence = circulation.Pence(constant)
	pc.VariablePence = circulation.Pence(variable)
	pc.Mode = circulation.ReproductionMode(mode)
	pc.MinCapitalisationIncrement = circulation.Pence(minIncrement)
	pc.LatentMoneyCapital = circulation.LatentMoneyCapital{
		ProductiveCircuitID: pc.ID,
		Accumulated:         circulation.Pence(latentAcc),
		Threshold:           circulation.Pence(minIncrement),
	}
	pc.ReserveFund = circulation.ReserveFund{
		ProductiveCircuitID: pc.ID,
		Balance:             circulation.Pence(reserveBal),
	}
	return pc, nil
}

// --- Vol. II Ch. 1 — MoneyCircuitStore ---

func (m *MySQL) CreateMoneyCircuit(ctx context.Context, mc circulation.MoneyCircuit) (circulation.MoneyCircuit, error) {
	if mc.ID.IsZero() {
		mc.ID = circulation.NewMoneyCircuitID()
	}
	mc.Moment = circulation.MomentM
	if mc.Advance.AdvancedAt.IsZero() {
		mc.Advance.AdvancedAt = time.Now().UTC()
	}
	mc.CreatedAt = time.Now().UTC()
	const q = `
		INSERT INTO money_circuits
			(id, agent_id, advance_amount, advanced_at, moment, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(mc.ID), mc.Advance.AgentID, int64(mc.Advance.Amount),
		mc.Advance.AdvancedAt, string(mc.Moment), mc.CreatedAt)
	if err != nil {
		return circulation.MoneyCircuit{}, err
	}
	return mc, nil
}

func (m *MySQL) GetMoneyCircuit(ctx context.Context, id circulation.MoneyCircuitID) (circulation.MoneyCircuit, error) {
	const q = `
		SELECT id, agent_id, advance_amount, advanced_at, moment,
		       labour_amount, labour_power_hours,
		       means_amount, means_capacity_hours,
		       productive_constant, productive_variable, productive_entered_at,
		       commodity_id, value_original, value_surplus,
		       realised_pence, sold_at,
		       created_at
		FROM money_circuits WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	mc, err := scanMoneyCircuit(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return circulation.MoneyCircuit{}, ErrNotFound
		}
		return circulation.MoneyCircuit{}, err
	}
	return mc, nil
}

func (m *MySQL) ListMoneyCircuits(ctx context.Context, agentID string, moment circulation.CircuitMoment) ([]circulation.MoneyCircuit, error) {
	q := `
		SELECT id, agent_id, advance_amount, advanced_at, moment,
		       labour_amount, labour_power_hours,
		       means_amount, means_capacity_hours,
		       productive_constant, productive_variable, productive_entered_at,
		       commodity_id, value_original, value_surplus,
		       realised_pence, sold_at,
		       created_at
		FROM money_circuits WHERE 1=1`
	args := []any{}
	if agentID != "" {
		q += " AND agent_id = ?"
		args = append(args, agentID)
	}
	if moment != "" {
		q += " AND moment = ?"
		args = append(args, string(moment))
	}
	rows, err := m.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []circulation.MoneyCircuit
	for rows.Next() {
		mc, err := scanMoneyCircuit(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, mc)
	}
	if out == nil {
		return []circulation.MoneyCircuit{}, nil
	}
	return out, rows.Err()
}

func (m *MySQL) RecordPurchase(ctx context.Context, id circulation.MoneyCircuitID, p circulation.PurchasePhase) (circulation.MoneyCircuit, error) {
	if err := p.Validate(); err != nil {
		return circulation.MoneyCircuit{}, err
	}
	const q = `
		UPDATE money_circuits
		SET labour_amount=?, labour_power_hours=?,
		    means_amount=?, means_capacity_hours=?,
		    moment=?
		WHERE id=? AND moment=?`
	res, err := m.db.ExecContext(ctx, q,
		int64(p.LabourLeg.Amount), p.LabourLeg.LabourPowerHours,
		int64(p.MeansLeg.Amount), p.MeansLeg.MeansCapacityHours,
		string(circulation.MomentMC),
		string(id), string(circulation.MomentM))
	if err != nil {
		return circulation.MoneyCircuit{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		if _, getErr := m.GetMoneyCircuit(ctx, id); getErr != nil {
			return circulation.MoneyCircuit{}, getErr
		}
		return circulation.MoneyCircuit{}, circulation.ErrInvalidPhaseTransition
	}
	return m.GetMoneyCircuit(ctx, id)
}

func (m *MySQL) RecordProductive(ctx context.Context, id circulation.MoneyCircuitID, ps circulation.ProductiveState) (circulation.MoneyCircuit, error) {
	mc, err := m.GetMoneyCircuit(ctx, id)
	if err != nil {
		return circulation.MoneyCircuit{}, err
	}
	if mc.Advance.Amount != ps.TotalAdvance() {
		return circulation.MoneyCircuit{}, circulation.ErrMagnitudeNotPreserved
	}
	if ps.EnteredAt.IsZero() {
		ps.EnteredAt = time.Now().UTC()
	}
	const q = `
		UPDATE money_circuits
		SET productive_constant=?, productive_variable=?, productive_entered_at=?,
		    moment=?
		WHERE id=? AND moment=?`
	res, err := m.db.ExecContext(ctx, q,
		int64(ps.ConstantPence), int64(ps.VariablePence), ps.EnteredAt,
		string(circulation.MomentP),
		string(id), string(circulation.MomentMC))
	if err != nil {
		return circulation.MoneyCircuit{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return circulation.MoneyCircuit{}, circulation.ErrInvalidPhaseTransition
	}
	return m.GetMoneyCircuit(ctx, id)
}

func (m *MySQL) RecordCommodity(ctx context.Context, id circulation.MoneyCircuitID, cc circulation.CommodityCapital) (circulation.MoneyCircuit, error) {
	mc, err := m.GetMoneyCircuit(ctx, id)
	if err != nil {
		return circulation.MoneyCircuit{}, err
	}
	if mc.Productive == nil {
		return circulation.MoneyCircuit{}, circulation.ErrNoProductionPhase
	}
	const q = `
		UPDATE money_circuits
		SET commodity_id=?, value_original=?, value_surplus=?,
		    moment=?
		WHERE id=? AND moment=?`
	res, err := m.db.ExecContext(ctx, q,
		cc.CommodityID, int64(cc.ValueOriginal), int64(cc.ValueSurplus),
		string(circulation.MomentCPrime),
		string(id), string(circulation.MomentP))
	if err != nil {
		return circulation.MoneyCircuit{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return circulation.MoneyCircuit{}, circulation.ErrInvalidPhaseTransition
	}
	return m.GetMoneyCircuit(ctx, id)
}

func (m *MySQL) RecordRealisation(ctx context.Context, id circulation.MoneyCircuitID, r circulation.Realisation) (circulation.MoneyCircuit, error) {
	if r.SoldAt.IsZero() {
		r.SoldAt = time.Now().UTC()
	}
	const q = `
		UPDATE money_circuits
		SET realised_pence=?, sold_at=?,
		    moment=?
		WHERE id=? AND moment=?`
	res, err := m.db.ExecContext(ctx, q,
		int64(r.RealisedPence), r.SoldAt,
		string(circulation.MomentMPrime),
		string(id), string(circulation.MomentCPrime))
	if err != nil {
		return circulation.MoneyCircuit{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		if _, getErr := m.GetMoneyCircuit(ctx, id); getErr != nil {
			return circulation.MoneyCircuit{}, getErr
		}
		return circulation.MoneyCircuit{}, circulation.ErrInvalidPhaseTransition
	}
	return m.GetMoneyCircuit(ctx, id)
}

func scanMoneyCircuit(scan func(...any) error) (circulation.MoneyCircuit, error) {
	var (
		mc                  circulation.MoneyCircuit
		rawID               string
		agentID             string
		advanceAmount       int64
		advancedAt          time.Time
		moment              string
		labourAmount        sql.NullInt64
		labourPowerHours    sql.NullInt64
		meansAmount         sql.NullInt64
		meansCapacityHours  sql.NullInt64
		productiveConstant  sql.NullInt64
		productiveVariable  sql.NullInt64
		productiveEnteredAt sql.NullTime
		commodityID         sql.NullString
		valueOriginal       sql.NullInt64
		valueSurplus        sql.NullInt64
		realisedPence       sql.NullInt64
		soldAt              sql.NullTime
		createdAt           time.Time
	)
	if err := scan(
		&rawID, &agentID, &advanceAmount, &advancedAt, &moment,
		&labourAmount, &labourPowerHours,
		&meansAmount, &meansCapacityHours,
		&productiveConstant, &productiveVariable, &productiveEnteredAt,
		&commodityID, &valueOriginal, &valueSurplus,
		&realisedPence, &soldAt,
		&createdAt,
	); err != nil {
		return circulation.MoneyCircuit{}, err
	}
	mc.ID = circulation.MoneyCircuitID(rawID)
	mc.Advance = circulation.MoneyAdvance{
		Amount:     circulation.Pence(advanceAmount),
		AgentID:    agentID,
		AdvancedAt: advancedAt,
	}
	mc.Moment = circulation.CircuitMoment(moment)
	mc.CreatedAt = createdAt

	if labourAmount.Valid && meansAmount.Valid {
		mc.Purchase = &circulation.PurchasePhase{
			MoneyCircuitID: mc.ID,
			LabourLeg: circulation.LabourLeg{
				Amount:           circulation.Pence(labourAmount.Int64),
				LabourPowerHours: labourPowerHours.Int64,
			},
			MeansLeg: circulation.MeansLeg{
				Amount:             circulation.Pence(meansAmount.Int64),
				MeansCapacityHours: meansCapacityHours.Int64,
			},
		}
	}
	if productiveConstant.Valid && productiveVariable.Valid {
		ps := &circulation.ProductiveState{
			MoneyCircuitID: mc.ID,
			ConstantPence:  circulation.Pence(productiveConstant.Int64),
			VariablePence:  circulation.Pence(productiveVariable.Int64),
		}
		if productiveEnteredAt.Valid {
			ps.EnteredAt = productiveEnteredAt.Time
		}
		mc.Productive = ps
	}
	if valueOriginal.Valid {
		cc := &circulation.CommodityCapital{
			MoneyCircuitID: mc.ID,
			ValueOriginal:  circulation.Pence(valueOriginal.Int64),
			ValueSurplus:   circulation.Pence(valueSurplus.Int64),
		}
		if commodityID.Valid {
			cc.CommodityID = commodityID.String
		}
		mc.Commodity = cc
	}
	if realisedPence.Valid {
		r := &circulation.Realisation{
			MoneyCircuitID: mc.ID,
			RealisedPence:  circulation.Pence(realisedPence.Int64),
		}
		if soldAt.Valid {
			r.SoldAt = soldAt.Time
		}
		mc.Realisation = r
	}
	return mc, nil
}

// --- IndustrialCapitalStore (Vol. II Ch. 4) ----------------------------------

func (m *MySQL) CreateIndustrialCapital(ctx context.Context, ic circulation.IndustrialCapital) (circulation.IndustrialCapital, error) {
	if err := ic.Validate(); err != nil {
		return circulation.IndustrialCapital{}, err
	}
	if ic.ID.IsZero() {
		ic.ID = circulation.NewIndustrialCapitalID()
	}
	now := m.now().UTC()
	ic.CreatedAt = now
	ic.UpdatedAt = now
	if ic.Status == "" {
		ic.Status = circulation.StatusActive
	}
	const q = `INSERT INTO industrial_capitals
		(id, agent_id, money_circuit_id, productive_circuit_id, commodity_circuit_id,
		 total_pence, economy_mode, stagnation_tolerance_ticks, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(ic.ID), ic.AgentID,
		string(ic.MoneyCircuitID), string(ic.ProductiveCircuitID), string(ic.CommodityCircuitID),
		int64(ic.TotalPence), string(ic.EconomyMode), ic.StagnationToleranceTicks,
		string(ic.Status), ic.CreatedAt, ic.UpdatedAt,
	)
	if err != nil {
		if isDuplicateKey(err) {
			return circulation.IndustrialCapital{}, ErrAlreadyExists
		}
		return circulation.IndustrialCapital{}, err
	}
	return ic, nil
}

func (m *MySQL) GetIndustrialCapital(ctx context.Context, id circulation.IndustrialCapitalID) (circulation.IndustrialCapital, error) {
	const q = `SELECT id, agent_id, money_circuit_id, productive_circuit_id, commodity_circuit_id,
		total_pence, economy_mode, stagnation_tolerance_ticks, status, created_at, updated_at
		FROM industrial_capitals WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	ic, err := scanIndustrialCapital(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return circulation.IndustrialCapital{}, ErrNotFound
		}
		return circulation.IndustrialCapital{}, err
	}
	// Populate latest StageDistribution.
	const sdQ = `SELECT id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence
		FROM stage_distributions WHERE industrial_capital_id = ?
		ORDER BY at_time DESC LIMIT 1`
	sdRow := m.db.QueryRowContext(ctx, sdQ, string(id))
	sd, sdErr := scanStageDistribution(sdRow.Scan)
	if sdErr == nil {
		ic.Latest = &sd
	}
	// Populate open StageBlocks.
	const sbQ = `SELECT id, industrial_capital_id, stage, reason, opened_at, closed_at
		FROM stage_blocks WHERE industrial_capital_id = ? AND closed_at IS NULL`
	sbRows, err := m.db.QueryContext(ctx, sbQ, string(id))
	if err != nil {
		return circulation.IndustrialCapital{}, err
	}
	defer sbRows.Close()
	ic.OpenBlocks = []circulation.StageBlock{}
	for sbRows.Next() {
		b, err := scanStageBlock(sbRows.Scan)
		if err != nil {
			return circulation.IndustrialCapital{}, err
		}
		ic.OpenBlocks = append(ic.OpenBlocks, b)
	}
	return ic, sbRows.Err()
}

func (m *MySQL) ListIndustrialCapitals(ctx context.Context, agentID, status, economyMode string) ([]circulation.IndustrialCapital, error) {
	q := `SELECT id, agent_id, money_circuit_id, productive_circuit_id, commodity_circuit_id,
		total_pence, economy_mode, stagnation_tolerance_ticks, status, created_at, updated_at
		FROM industrial_capitals WHERE 1=1`
	args := []any{}
	if agentID != "" {
		q += " AND agent_id = ?"
		args = append(args, agentID)
	}
	if status != "" {
		q += " AND status = ?"
		args = append(args, status)
	}
	if economyMode != "" {
		q += " AND economy_mode = ?"
		args = append(args, economyMode)
	}
	q += " ORDER BY created_at ASC"
	rows, err := m.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []circulation.IndustrialCapital{}
	for rows.Next() {
		ic, err := scanIndustrialCapital(rows.Scan)
		if err != nil {
			return nil, err
		}
		out = append(out, ic)
	}
	return out, rows.Err()
}

func (m *MySQL) RecordCapitalPart(ctx context.Context, id circulation.IndustrialCapitalID, part circulation.CapitalPart) (circulation.CapitalPart, error) {
	if _, err := m.GetIndustrialCapital(ctx, id); err != nil {
		return circulation.CapitalPart{}, err
	}
	if part.ID.IsZero() {
		part.ID = circulation.NewCapitalPartID()
	}
	part.IndustrialCapitalID = id
	if part.EnteredStageAt.IsZero() {
		part.EnteredStageAt = m.now().UTC()
	}
	const q = `INSERT INTO capital_parts (id, industrial_capital_id, pence, stage, entered_stage_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(part.ID), string(id), int64(part.Pence),
		string(part.Stage), part.EnteredStageAt,
	)
	if err != nil {
		if isDuplicateKey(err) {
			return circulation.CapitalPart{}, ErrAlreadyExists
		}
		return circulation.CapitalPart{}, err
	}
	return part, nil
}

func (m *MySQL) Snapshot(ctx context.Context, id circulation.IndustrialCapitalID, sd circulation.StageDistribution) (circulation.StageDistribution, error) {
	ic, err := m.GetIndustrialCapital(ctx, id)
	if err != nil {
		return circulation.StageDistribution{}, err
	}
	if err := sd.Validate(ic.TotalPence); err != nil {
		return circulation.StageDistribution{}, err
	}
	if sd.ID.IsZero() {
		sd.ID = circulation.NewStageDistributionID()
	}
	sd.IndustrialCapitalID = id
	if sd.At.IsZero() {
		sd.At = m.now().UTC()
	}
	const q = `INSERT INTO stage_distributions
		(id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err = m.db.ExecContext(ctx, q,
		string(sd.ID), string(id), sd.At,
		int64(sd.MoneyPence), int64(sd.ProductionPence), int64(sd.CommodityPence),
	)
	if err != nil {
		return circulation.StageDistribution{}, err
	}
	return sd, nil
}

func (m *MySQL) OpenBlock(ctx context.Context, id circulation.IndustrialCapitalID, b circulation.StageBlock) (circulation.StageBlock, error) {
	ic, err := m.GetIndustrialCapital(ctx, id)
	if err != nil {
		return circulation.StageBlock{}, err
	}
	if b.ID.IsZero() {
		b.ID = circulation.NewStageBlockID()
	}
	b.IndustrialCapitalID = id
	if b.OpenedAt.IsZero() {
		b.OpenedAt = m.now().UTC()
	}
	const q = `INSERT INTO stage_blocks (id, industrial_capital_id, stage, reason, opened_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err = m.db.ExecContext(ctx, q,
		string(b.ID), string(id), string(b.Stage), string(b.Reason), b.OpenedAt,
	)
	if err != nil {
		return circulation.StageBlock{}, err
	}
	// Check if open block count exceeds tolerance → halt the capital.
	const countQ = `SELECT COUNT(*) FROM stage_blocks
		WHERE industrial_capital_id = ? AND closed_at IS NULL`
	var openCount int64
	if err := m.db.QueryRowContext(ctx, countQ, string(id)).Scan(&openCount); err != nil {
		return b, nil // non-fatal — block was inserted
	}
	if ic.StagnationToleranceTicks > 0 && openCount > ic.StagnationToleranceTicks {
		_, _ = m.db.ExecContext(ctx,
			`UPDATE industrial_capitals SET status = ?, updated_at = ? WHERE id = ?`,
			string(circulation.StatusHalted), m.now().UTC(), string(id),
		)
	}
	return b, nil
}

func (m *MySQL) CloseBlock(ctx context.Context, id circulation.IndustrialCapitalID, blockID circulation.StageBlockID) (circulation.StageBlock, error) {
	if _, err := m.GetIndustrialCapital(ctx, id); err != nil {
		return circulation.StageBlock{}, err
	}
	now := m.now().UTC()
	const q = `UPDATE stage_blocks SET closed_at = ? WHERE id = ? AND industrial_capital_id = ? AND closed_at IS NULL`
	res, err := m.db.ExecContext(ctx, q, now, string(blockID), string(id))
	if err != nil {
		return circulation.StageBlock{}, err
	}
	n, _ := res.RowsAffected()
	const selQ = `SELECT id, industrial_capital_id, stage, reason, opened_at, closed_at
		FROM stage_blocks WHERE id = ? AND industrial_capital_id = ?`
	row := m.db.QueryRowContext(ctx, selQ, string(blockID), string(id))
	b, err := scanStageBlock(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || n == 0 {
			return circulation.StageBlock{}, ErrNotFound
		}
		return circulation.StageBlock{}, err
	}
	return b, nil
}

func (m *MySQL) RecordValueRevolution(ctx context.Context, res circulation.ValueRevolutionResult) (circulation.ValueRevolutionResult, error) {
	id := res.Event.IndustrialCapitalID
	if _, err := m.GetIndustrialCapital(ctx, id); err != nil {
		return circulation.ValueRevolutionResult{}, err
	}
	if res.Event.ID.IsZero() {
		res.Event.ID = circulation.NewValueRevolutionEventID()
	}
	if res.Event.OccurredAt.IsZero() {
		res.Event.OccurredAt = m.now().UTC()
	}
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return circulation.ValueRevolutionResult{}, err
	}
	const evtQ = `INSERT INTO value_revolution_events
		(id, industrial_capital_id, affecting, direction_pence, occurred_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, evtQ,
		string(res.Event.ID), string(id),
		string(res.Event.Affecting), int64(res.Event.DirectionPence), res.Event.OccurredAt,
	)
	if err != nil {
		_ = tx.Rollback()
		return circulation.ValueRevolutionResult{}, err
	}
	if res.SetFree != nil {
		if res.SetFree.ID.IsZero() {
			res.SetFree.ID = circulation.NewMoneyCapitalSetFreeID()
		}
		const sfQ = `INSERT INTO money_capital_set_free
			(id, value_revolution_event_id, industrial_capital_id, amount_pence)
			VALUES (?, ?, ?, ?)`
		_, err = tx.ExecContext(ctx, sfQ,
			string(res.SetFree.ID), string(res.Event.ID), string(id), int64(res.SetFree.AmountPence),
		)
		if err != nil {
			_ = tx.Rollback()
			return circulation.ValueRevolutionResult{}, err
		}
	}
	if res.TiedUp != nil {
		if res.TiedUp.ID.IsZero() {
			res.TiedUp.ID = circulation.NewMoneyCapitalTiedUpID()
		}
		const tuQ = `INSERT INTO money_capital_tied_up
			(id, value_revolution_event_id, industrial_capital_id, amount_pence)
			VALUES (?, ?, ?, ?)`
		_, err = tx.ExecContext(ctx, tuQ,
			string(res.TiedUp.ID), string(res.Event.ID), string(id), int64(res.TiedUp.AmountPence),
		)
		if err != nil {
			_ = tx.Rollback()
			return circulation.ValueRevolutionResult{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return circulation.ValueRevolutionResult{}, err
	}
	return res, nil
}

func (m *MySQL) RecordInterlock(ctx context.Context, id circulation.IndustrialCapitalID, mi circulation.MetamorphosisInterlock) (circulation.MetamorphosisInterlock, error) {
	if err := mi.Validate(); err != nil {
		return circulation.MetamorphosisInterlock{}, err
	}
	if _, err := m.GetIndustrialCapital(ctx, id); err != nil {
		return circulation.MetamorphosisInterlock{}, err
	}
	if mi.ID.IsZero() {
		mi.ID = circulation.NewMetamorphosisInterlockID()
	}
	mi.BuyerIndustrialCapitalID = id
	if mi.OccurredAt.IsZero() {
		mi.OccurredAt = m.now().UTC()
	}
	const q = `INSERT INTO metamorphosis_interlocks
		(id, buyer_industrial_capital_id, seller_industrial_capital_id, seller_mp_source_origin, pence, occurred_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(mi.ID), string(id), string(mi.SellerIndustrialCapitalID),
		string(mi.SellerMPSourceOrigin), int64(mi.Pence), mi.OccurredAt,
	)
	if err != nil {
		return circulation.MetamorphosisInterlock{}, err
	}
	return mi, nil
}

func (m *MySQL) RecordSupplyDemand(ctx context.Context, sdi circulation.SupplyDemandImbalance) (circulation.SupplyDemandImbalance, error) {
	if _, err := m.GetIndustrialCapital(ctx, sdi.IndustrialCapitalID); err != nil {
		return circulation.SupplyDemandImbalance{}, err
	}
	if sdi.ID.IsZero() {
		sdi.ID = circulation.NewSupplyDemandImbalanceID()
	}
	const q = `INSERT INTO supply_demand_imbalances
		(id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(sdi.ID), string(sdi.IndustrialCapitalID), sdi.Period,
		int64(sdi.DemandPence), int64(sdi.SupplyPence), int64(sdi.ExcessPence),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return circulation.SupplyDemandImbalance{}, ErrAlreadyExists
		}
		return circulation.SupplyDemandImbalance{}, err
	}
	return sdi, nil
}

func (m *MySQL) GetSupplyDemand(ctx context.Context, id circulation.IndustrialCapitalID, period string) (circulation.SupplyDemandImbalance, error) {
	const q = `SELECT id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence
		FROM supply_demand_imbalances WHERE industrial_capital_id = ? AND period = ? LIMIT 1`
	row := m.db.QueryRowContext(ctx, q, string(id), period)
	var (
		rawID  string
		rawCap string
		per    string
		demand int64
		supply int64
		excess int64
	)
	err := row.Scan(&rawID, &rawCap, &per, &demand, &supply, &excess)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return circulation.SupplyDemandImbalance{}, ErrNotFound
		}
		return circulation.SupplyDemandImbalance{}, err
	}
	return circulation.SupplyDemandImbalance{
		ID:                  circulation.SupplyDemandImbalanceID(rawID),
		IndustrialCapitalID: circulation.IndustrialCapitalID(rawCap),
		Period:              period,
		DemandPence:         circulation.Pence(demand),
		SupplyPence:         circulation.Pence(supply),
		ExcessPence:         circulation.Pence(excess),
	}, nil
}

func (m *MySQL) AggregateSupplyDemand(ctx context.Context, period string) (circulation.AggregateSupplyDemandImbalance, error) {
	const q = `SELECT COALESCE(SUM(demand_pence),0), COALESCE(SUM(supply_pence),0), COALESCE(SUM(excess_pence),0)
		FROM supply_demand_imbalances WHERE period = ?`
	row := m.db.QueryRowContext(ctx, q, period)
	var demand, supply, excess int64
	if err := row.Scan(&demand, &supply, &excess); err != nil {
		return circulation.AggregateSupplyDemandImbalance{}, err
	}
	return circulation.AggregateSupplyDemandImbalance{
		Period:      period,
		DemandPence: circulation.Pence(demand),
		SupplyPence: circulation.Pence(supply),
		ExcessPence: circulation.Pence(excess),
	}, nil
}

func (m *MySQL) SetSinkingFund(ctx context.Context, id circulation.IndustrialCapitalID, sf circulation.SinkingFund) (circulation.SinkingFund, error) {
	if err := sf.Validate(); err != nil {
		return circulation.SinkingFund{}, err
	}
	if _, err := m.GetIndustrialCapital(ctx, id); err != nil {
		return circulation.SinkingFund{}, err
	}
	if sf.ID.IsZero() {
		sf.ID = circulation.NewSinkingFundID()
	}
	sf.IndustrialCapitalID = id
	const q = `INSERT INTO sinking_funds
		(id, industrial_capital_id, fixed_capital_pence, lifetime_years, annual_payment_pence, accumulated_pence)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			fixed_capital_pence  = VALUES(fixed_capital_pence),
			lifetime_years       = VALUES(lifetime_years),
			annual_payment_pence = VALUES(annual_payment_pence),
			accumulated_pence    = VALUES(accumulated_pence)`
	_, err := m.db.ExecContext(ctx, q,
		string(sf.ID), string(id),
		int64(sf.FixedCapitalPence), sf.LifetimeYears,
		int64(sf.AnnualPaymentPence), int64(sf.AccumulatedPence),
	)
	if err != nil {
		return circulation.SinkingFund{}, err
	}
	return sf, nil
}

func (m *MySQL) TickSinkingFund(ctx context.Context, id circulation.IndustrialCapitalID) (circulation.SinkingFund, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return circulation.SinkingFund{}, err
	}
	const selQ = `SELECT id, industrial_capital_id, fixed_capital_pence, lifetime_years,
		annual_payment_pence, accumulated_pence
		FROM sinking_funds WHERE industrial_capital_id = ? FOR UPDATE`
	row := tx.QueryRowContext(ctx, selQ, string(id))
	var (
		rawID, rawCap              string
		fixed, annual, accumulated int64
		lifetime                   int64
	)
	if err := row.Scan(&rawID, &rawCap, &fixed, &lifetime, &annual, &accumulated); err != nil {
		_ = tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return circulation.SinkingFund{}, ErrNotFound
		}
		return circulation.SinkingFund{}, err
	}
	accumulated += annual
	const upQ = `UPDATE sinking_funds SET accumulated_pence = ? WHERE id = ?`
	if _, err := tx.ExecContext(ctx, upQ, accumulated, rawID); err != nil {
		_ = tx.Rollback()
		return circulation.SinkingFund{}, err
	}
	if err := tx.Commit(); err != nil {
		return circulation.SinkingFund{}, err
	}
	return circulation.SinkingFund{
		ID:                  circulation.SinkingFundID(rawID),
		IndustrialCapitalID: circulation.IndustrialCapitalID(rawCap),
		FixedCapitalPence:   circulation.Pence(fixed),
		LifetimeYears:       lifetime,
		AnnualPaymentPence:  circulation.Pence(annual),
		AccumulatedPence:    circulation.Pence(accumulated),
	}, nil
}

// scan helpers for IndustrialCapital tables

func scanIndustrialCapital(scan func(...any) error) (circulation.IndustrialCapital, error) {
	var (
		rawID, agentID                                            string
		moneyCircuitID, productiveCircuitID, commodityCircuitID string
		totalPence                                                int64
		economyMode, status                                       string
		stagnationTolerance                                       int64
		createdAt, updatedAt                                      time.Time
	)
	if err := scan(
		&rawID, &agentID, &moneyCircuitID, &productiveCircuitID, &commodityCircuitID,
		&totalPence, &economyMode, &stagnationTolerance, &status, &createdAt, &updatedAt,
	); err != nil {
		return circulation.IndustrialCapital{}, err
	}
	return circulation.IndustrialCapital{
		ID:                       circulation.IndustrialCapitalID(rawID),
		AgentID:                  agentID,
		MoneyCircuitID:           circulation.MoneyCircuitID(moneyCircuitID),
		ProductiveCircuitID:      circulation.ProductiveCircuitID(productiveCircuitID),
		CommodityCircuitID:       circulation.CommodityCircuitID(commodityCircuitID),
		TotalPence:               circulation.Pence(totalPence),
		EconomyMode:              circulation.EconomyMode(economyMode),
		StagnationToleranceTicks: stagnationTolerance,
		Status:                   circulation.IndustrialCapitalStatus(status),
		CreatedAt:                createdAt,
		UpdatedAt:                updatedAt,
	}, nil
}

func scanStageDistribution(scan func(...any) error) (circulation.StageDistribution, error) {
	var (
		rawID, rawCap          string
		atTime                 time.Time
		money, prod, commodity int64
	)
	if err := scan(&rawID, &rawCap, &atTime, &money, &prod, &commodity); err != nil {
		return circulation.StageDistribution{}, err
	}
	return circulation.StageDistribution{
		ID:                  circulation.StageDistributionID(rawID),
		IndustrialCapitalID: circulation.IndustrialCapitalID(rawCap),
		At:                  atTime,
		MoneyPence:          circulation.Pence(money),
		ProductionPence:     circulation.Pence(prod),
		CommodityPence:      circulation.Pence(commodity),
	}, nil
}

func scanStageBlock(scan func(...any) error) (circulation.StageBlock, error) {
	var (
		rawID, rawCap string
		stage, reason string
		openedAt      time.Time
		closedAt      sql.NullTime
	)
	if err := scan(&rawID, &rawCap, &stage, &reason, &openedAt, &closedAt); err != nil {
		return circulation.StageBlock{}, err
	}
	b := circulation.StageBlock{
		ID:                  circulation.StageBlockID(rawID),
		IndustrialCapitalID: circulation.IndustrialCapitalID(rawCap),
		Stage:               circulation.CapitalStage(stage),
		Reason:              circulation.StageBlockReason(reason),
		OpenedAt:            openedAt,
	}
	if closedAt.Valid {
		t := closedAt.Time
		b.ClosedAt = &t
	}
	return b, nil
}
