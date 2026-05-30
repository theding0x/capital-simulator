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
	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
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

// CreateProfitRate persists a, assigning an ID and timestamp when absent.
func (m *MySQL) CreateProfitRate(ctx context.Context, a profit.ProfitRateAnalysis) (profit.ProfitRateAnalysis, error) {
	if a.ID.IsZero() {
		a.ID = profit.NewProfitRateID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO profit_rates
		(id, constant, variable, surplus_value, total_capital, profit_rate_bp, surplus_value_rate_bp, mystification_bp, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID),
		int64(a.ConstantCapital), int64(a.VariableCapital),
		int64(a.ProfitRate.SurplusValue), int64(a.ProfitRate.TotalCapital),
		a.ProfitRate.BasisPoints, int64(a.SurplusValueRate), int64(a.Mystification),
		a.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return profit.ProfitRateAnalysis{}, ErrAlreadyExists
		}
		return profit.ProfitRateAnalysis{}, err
	}
	return a, nil
}

// GetProfitRate returns the analysis with id, or ErrNotFound.
func (m *MySQL) GetProfitRate(ctx context.Context, id profit.ProfitRateID) (profit.ProfitRateAnalysis, error) {
	const q = `SELECT id, constant, variable, surplus_value, total_capital,
		profit_rate_bp, surplus_value_rate_bp, mystification_bp, created_at
		FROM profit_rates WHERE id = ?`
	return scanProfitRate(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListProfitRates returns all stored analyses, newest first.
func (m *MySQL) ListProfitRates(ctx context.Context) ([]profit.ProfitRateAnalysis, error) {
	const q = `SELECT id, constant, variable, surplus_value, total_capital,
		profit_rate_bp, surplus_value_rate_bp, mystification_bp, created_at
		FROM profit_rates ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []profit.ProfitRateAnalysis
	for rows.Next() {
		a, err := scanProfitRate(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// CreateVariation persists a, assigning an ID and timestamp when absent.
func (m *MySQL) CreateVariation(ctx context.Context, a profit.VariationAnalysis) (profit.VariationAnalysis, error) {
	if a.ID.IsZero() {
		a.ID = profit.NewVariationAnalysisID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO variation_analyses
		(id, vcase, initial_c, initial_v, initial_s_rate_bp, changed_c, changed_v, changed_s_rate_bp,
		 old_profit_rate_bp, new_profit_rate_bp, old_composition_bp, new_composition_bp, proportional_change, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID), string(a.Case),
		int64(a.Initial.C), int64(a.Initial.V), int64(a.Initial.SRate),
		int64(a.Changed.C), int64(a.Changed.V), int64(a.Changed.SRate),
		a.OldProfitRate, a.NewProfitRate, a.OldCompositionBP, a.NewCompositionBP,
		a.ProportionalChange, a.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return profit.VariationAnalysis{}, ErrAlreadyExists
		}
		return profit.VariationAnalysis{}, err
	}
	return a, nil
}

// GetVariation returns the variation analysis with id, or ErrNotFound.
func (m *MySQL) GetVariation(ctx context.Context, id profit.VariationAnalysisID) (profit.VariationAnalysis, error) {
	const q = `SELECT id, vcase, initial_c, initial_v, initial_s_rate_bp, changed_c, changed_v, changed_s_rate_bp,
		old_profit_rate_bp, new_profit_rate_bp, old_composition_bp, new_composition_bp, proportional_change, created_at
		FROM variation_analyses WHERE id = ?`
	return scanVariation(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListVariations returns all stored variation analyses, newest first.
func (m *MySQL) ListVariations(ctx context.Context) ([]profit.VariationAnalysis, error) {
	const q = `SELECT id, vcase, initial_c, initial_v, initial_s_rate_bp, changed_c, changed_v, changed_s_rate_bp,
		old_profit_rate_bp, new_profit_rate_bp, old_composition_bp, new_composition_bp, proportional_change, created_at
		FROM variation_analyses ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []profit.VariationAnalysis
	for rows.Next() {
		a, err := scanVariation(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// CreateTurnoverAnalysis persists a, assigning an ID and timestamp when absent.
func (m *MySQL) CreateTurnoverAnalysis(ctx context.Context, a profit.TurnoverAnalysis) (profit.TurnoverAnalysis, error) {
	if a.ID.IsZero() {
		a.ID = profit.NewTurnoverAnalysisID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO turnover_analyses
		(id, total_capital, variable_capital, surplus_value_rate_bp, turnovers,
		 annual_profit_rate_bp, single_turnover_rate_bp, annual_surplus_value_rate_bp, annual_wages, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID),
		int64(a.C), int64(a.V), int64(a.SRate), int64(a.N),
		a.AnnualProfitRate.BasisPoints, a.SingleTurnoverProfitRate,
		int64(a.AnnualSurplusValueRate), int64(a.AnnualWages),
		a.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return profit.TurnoverAnalysis{}, ErrAlreadyExists
		}
		return profit.TurnoverAnalysis{}, err
	}
	return a, nil
}

// GetTurnoverAnalysis returns the analysis with id, or ErrNotFound.
func (m *MySQL) GetTurnoverAnalysis(ctx context.Context, id profit.TurnoverAnalysisID) (profit.TurnoverAnalysis, error) {
	const q = `SELECT id, total_capital, variable_capital, surplus_value_rate_bp, turnovers,
		annual_profit_rate_bp, single_turnover_rate_bp, annual_surplus_value_rate_bp, annual_wages, created_at
		FROM turnover_analyses WHERE id = ?`
	return scanTurnoverAnalysis(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListTurnoverAnalyses returns all stored turnover analyses, newest first.
func (m *MySQL) ListTurnoverAnalyses(ctx context.Context) ([]profit.TurnoverAnalysis, error) {
	const q = `SELECT id, total_capital, variable_capital, surplus_value_rate_bp, turnovers,
		annual_profit_rate_bp, single_turnover_rate_bp, annual_surplus_value_rate_bp, annual_wages, created_at
		FROM turnover_analyses ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []profit.TurnoverAnalysis
	for rows.Next() {
		a, err := scanTurnoverAnalysis(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// CreateEconomyAnalysis persists a, assigning an ID and timestamp when absent.
func (m *MySQL) CreateEconomyAnalysis(ctx context.Context, a profit.EconomyAnalysis) (profit.EconomyAnalysis, error) {
	if a.ID.IsZero() {
		a.ID = profit.NewEconomyAnalysisID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO economy_analyses
		(id, kind, constant_capital, variable_capital, surplus_value, saving,
		 old_profit_rate_bp, new_profit_rate_bp, profit_rate_gain_bp, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID), string(a.Economy.Kind),
		int64(a.Economy.ConstantCapital), int64(a.Economy.VariableCapital),
		int64(a.Economy.SurplusValue), int64(a.Economy.Saving),
		a.OldProfitRate, a.NewProfitRate, a.ProfitRateGain,
		a.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return profit.EconomyAnalysis{}, ErrAlreadyExists
		}
		return profit.EconomyAnalysis{}, err
	}
	return a, nil
}

// GetEconomyAnalysis returns the analysis with id, or ErrNotFound.
func (m *MySQL) GetEconomyAnalysis(ctx context.Context, id profit.EconomyAnalysisID) (profit.EconomyAnalysis, error) {
	const q = `SELECT id, kind, constant_capital, variable_capital, surplus_value, saving,
		old_profit_rate_bp, new_profit_rate_bp, profit_rate_gain_bp, created_at
		FROM economy_analyses WHERE id = ?`
	return scanEconomyAnalysis(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListEconomyAnalyses returns all stored economy analyses, newest first.
func (m *MySQL) ListEconomyAnalyses(ctx context.Context) ([]profit.EconomyAnalysis, error) {
	const q = `SELECT id, kind, constant_capital, variable_capital, surplus_value, saving,
		old_profit_rate_bp, new_profit_rate_bp, profit_rate_gain_bp, created_at
		FROM economy_analyses ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []profit.EconomyAnalysis
	for rows.Next() {
		a, err := scanEconomyAnalysis(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// CreatePriceFluctuationAnalysis persists a, assigning an ID and timestamp when absent.
func (m *MySQL) CreatePriceFluctuationAnalysis(ctx context.Context, a profit.PriceFluctuationAnalysis) (profit.PriceFluctuationAnalysis, error) {
	if a.ID.IsZero() {
		a.ID = profit.NewPriceFluctuationAnalysisID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO price_fluctuation_analyses
		(id, kind, fixed_capital, original_material_capital, price_factor, variable_capital, surplus_value,
		 old_profit_rate_bp, new_profit_rate_bp, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID), string(a.Kind),
		int64(a.Effect.FixedCapital), int64(a.Effect.OriginalMaterialCapital), a.Effect.PriceFactor,
		int64(a.Effect.VariableCapital), int64(a.Effect.SurplusValue),
		a.OldProfitRate, a.NewProfitRate,
		a.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return profit.PriceFluctuationAnalysis{}, ErrAlreadyExists
		}
		return profit.PriceFluctuationAnalysis{}, err
	}
	return a, nil
}

// GetPriceFluctuationAnalysis returns the analysis with id, or ErrNotFound.
func (m *MySQL) GetPriceFluctuationAnalysis(ctx context.Context, id profit.PriceFluctuationAnalysisID) (profit.PriceFluctuationAnalysis, error) {
	const q = `SELECT id, kind, fixed_capital, original_material_capital, price_factor, variable_capital, surplus_value,
		old_profit_rate_bp, new_profit_rate_bp, created_at
		FROM price_fluctuation_analyses WHERE id = ?`
	return scanPriceFluctuationAnalysis(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListPriceFluctuationAnalyses returns all stored analyses, newest first.
func (m *MySQL) ListPriceFluctuationAnalyses(ctx context.Context) ([]profit.PriceFluctuationAnalysis, error) {
	const q = `SELECT id, kind, fixed_capital, original_material_capital, price_factor, variable_capital, surplus_value,
		old_profit_rate_bp, new_profit_rate_bp, created_at
		FROM price_fluctuation_analyses ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []profit.PriceFluctuationAnalysis
	for rows.Next() {
		a, err := scanPriceFluctuationAnalysis(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// CreateCompositionEffect persists a, assigning an ID and timestamp when absent.
func (m *MySQL) CreateCompositionEffect(ctx context.Context, a profit.CompositionEffectAnalysis) (profit.CompositionEffectAnalysis, error) {
	if a.ID.IsZero() {
		a.ID = profit.NewCompositionEffectID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO composition_effects
		(id, s_a, v_a, c_a, s_b, v_b, c_b, profit_rate_a_bp, profit_rate_b_bp, rate_difference_bp, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID),
		int64(a.Effect.SA), int64(a.Effect.VA), int64(a.Effect.CA),
		int64(a.Effect.SB), int64(a.Effect.VB), int64(a.Effect.CB),
		a.ProfitRateA, a.ProfitRateB, a.RateDifference,
		a.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return profit.CompositionEffectAnalysis{}, ErrAlreadyExists
		}
		return profit.CompositionEffectAnalysis{}, err
	}
	return a, nil
}

// GetCompositionEffect returns the comparison with id, or ErrNotFound.
func (m *MySQL) GetCompositionEffect(ctx context.Context, id profit.CompositionEffectID) (profit.CompositionEffectAnalysis, error) {
	const q = `SELECT id, s_a, v_a, c_a, s_b, v_b, c_b, profit_rate_a_bp, profit_rate_b_bp, rate_difference_bp, created_at
		FROM composition_effects WHERE id = ?`
	return scanCompositionEffect(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListCompositionEffects returns all stored comparisons, newest first.
func (m *MySQL) ListCompositionEffects(ctx context.Context) ([]profit.CompositionEffectAnalysis, error) {
	const q = `SELECT id, s_a, v_a, c_a, s_b, v_b, c_b, profit_rate_a_bp, profit_rate_b_bp, rate_difference_bp, created_at
		FROM composition_effects ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []profit.CompositionEffectAnalysis
	for rows.Next() {
		a, err := scanCompositionEffect(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// CreateMagnitudeChange persists a, assigning an ID and timestamp when absent.
func (m *MySQL) CreateMagnitudeChange(ctx context.Context, a profit.MagnitudeChangeAnalysis) (profit.MagnitudeChangeAnalysis, error) {
	if a.ID.IsZero() {
		a.ID = profit.NewMagnitudeChangeID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO magnitude_changes
		(id, kind, original_capital, original_profit, factor, old_rate_bp, new_rate_bp, rate_unchanged, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID), string(a.Change.Kind),
		int64(a.Change.OriginalCapital), int64(a.Change.OriginalProfit), a.Change.Factor,
		a.OldRate, a.NewRate, a.RateUnchanged,
		a.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return profit.MagnitudeChangeAnalysis{}, ErrAlreadyExists
		}
		return profit.MagnitudeChangeAnalysis{}, err
	}
	return a, nil
}

// GetMagnitudeChange returns the change with id, or ErrNotFound.
func (m *MySQL) GetMagnitudeChange(ctx context.Context, id profit.MagnitudeChangeID) (profit.MagnitudeChangeAnalysis, error) {
	const q = `SELECT id, kind, original_capital, original_profit, factor, old_rate_bp, new_rate_bp, rate_unchanged, created_at
		FROM magnitude_changes WHERE id = ?`
	return scanMagnitudeChange(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListMagnitudeChanges returns all stored changes, newest first.
func (m *MySQL) ListMagnitudeChanges(ctx context.Context) ([]profit.MagnitudeChangeAnalysis, error) {
	const q = `SELECT id, kind, original_capital, original_profit, factor, old_rate_bp, new_rate_bp, rate_unchanged, created_at
		FROM magnitude_changes ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []profit.MagnitudeChangeAnalysis
	for rows.Next() {
		a, err := scanMagnitudeChange(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// CreateProductionSphere persists s, assigning an ID and timestamp when absent.
func (m *MySQL) CreateProductionSphere(ctx context.Context, s avgprofit.ProductionSphere) (avgprofit.ProductionSphere, error) {
	if s.ID.IsZero() {
		s.ID = avgprofit.NewProductionSphereID()
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO production_spheres
		(id, name, c, v, s_rate, constant_capital, surplus_value, individual_profit_rate, variable_percent, labour_power_index, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(s.ID), s.Name,
		int64(s.C), int64(s.V), int64(s.SRate),
		int64(s.ConstantCapital), int64(s.SurplusValue),
		int64(s.IndividualProfitRate), s.VariablePercent,
		int64(s.LabourPowerIndex), s.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return avgprofit.ProductionSphere{}, ErrAlreadyExists
		}
		return avgprofit.ProductionSphere{}, err
	}
	return s, nil
}

// GetProductionSphere returns the sphere with id, or ErrNotFound.
func (m *MySQL) GetProductionSphere(ctx context.Context, id avgprofit.ProductionSphereID) (avgprofit.ProductionSphere, error) {
	const q = `SELECT id, name, c, v, s_rate, constant_capital, surplus_value, individual_profit_rate, variable_percent, labour_power_index, created_at
		FROM production_spheres WHERE id = ?`
	return scanProductionSphere(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListProductionSpheres returns all stored spheres, newest first.
func (m *MySQL) ListProductionSpheres(ctx context.Context) ([]avgprofit.ProductionSphere, error) {
	const q = `SELECT id, name, c, v, s_rate, constant_capital, surplus_value, individual_profit_rate, variable_percent, labour_power_index, created_at
		FROM production_spheres ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []avgprofit.ProductionSphere
	for rows.Next() {
		s, err := scanProductionSphere(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
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

func scanProfitRate(s rowScanner) (profit.ProfitRateAnalysis, error) {
	var (
		id                                 string
		constant, variable, surplus, total int64
		profitBP, surplusBP, mystBP        int64
		createdAt                          time.Time
	)
	err := s.Scan(&id, &constant, &variable, &surplus, &total,
		&profitBP, &surplusBP, &mystBP, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return profit.ProfitRateAnalysis{}, ErrNotFound
		}
		return profit.ProfitRateAnalysis{}, err
	}
	return profit.ProfitRateAnalysis{
		ID:              profit.ProfitRateID(id),
		ConstantCapital: profit.ConstantCapital(constant),
		VariableCapital: profit.VariableCapital(variable),
		ProfitRate: profit.RateOfProfit{
			SurplusValue: profit.SurplusValue(surplus),
			TotalCapital: profit.TotalCapital(total),
			BasisPoints:  profitBP,
		},
		SurplusValueRate: profit.RateOfSurplusValue(surplusBP),
		Mystification:    profit.MystificationDegree(mystBP),
		CreatedAt:        createdAt,
	}, nil
}

func scanVariation(s rowScanner) (profit.VariationAnalysis, error) {
	var (
		id, vcase                             string
		initC, initV, initS, chgC, chgV, chgS int64
		oldP, newP, oldComp, newComp          int64
		proportional                          bool
		createdAt                             time.Time
	)
	err := s.Scan(&id, &vcase, &initC, &initV, &initS, &chgC, &chgV, &chgS,
		&oldP, &newP, &oldComp, &newComp, &proportional, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return profit.VariationAnalysis{}, ErrNotFound
		}
		return profit.VariationAnalysis{}, err
	}
	return profit.VariationAnalysis{
		ID:   profit.VariationAnalysisID(id),
		Case: profit.VariationCase(vcase),
		Initial: profit.ProfitRateFormula{
			C:     profit.ConstantCapital(initC),
			V:     profit.VariableCapital(initV),
			SRate: profit.RateOfSurplusValue(initS),
		},
		Changed: profit.ProfitRateFormula{
			C:     profit.ConstantCapital(chgC),
			V:     profit.VariableCapital(chgV),
			SRate: profit.RateOfSurplusValue(chgS),
		},
		OldProfitRate:      oldP,
		NewProfitRate:      newP,
		OldCompositionBP:   oldComp,
		NewCompositionBP:   newComp,
		ProportionalChange: proportional,
		CreatedAt:          createdAt,
	}, nil
}

func scanTurnoverAnalysis(s rowScanner) (profit.TurnoverAnalysis, error) {
	var (
		id                                         string
		totalCapital, variableCapital, sRateBP, n  int64
		annualBP, singleBP, annualSurplusBP, wages int64
		createdAt                                  time.Time
	)
	err := s.Scan(&id, &totalCapital, &variableCapital, &sRateBP, &n,
		&annualBP, &singleBP, &annualSurplusBP, &wages, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return profit.TurnoverAnalysis{}, ErrNotFound
		}
		return profit.TurnoverAnalysis{}, err
	}
	return profit.TurnoverAnalysis{
		ID:    profit.TurnoverAnalysisID(id),
		C:     profit.TotalCapital(totalCapital),
		V:     profit.VariableCapitalPerTurnover(variableCapital),
		SRate: profit.RateOfSurplusValue(sRateBP),
		N:     profit.TurnoverCount(n),
		AnnualProfitRate: profit.AnnualProfitRate{
			SurplusValueRate: profit.RateOfSurplusValue(sRateBP),
			Turnovers:        profit.TurnoverCount(n),
			VariableCapital:  profit.VariableCapitalPerTurnover(variableCapital),
			TotalCapital:     profit.TotalCapital(totalCapital),
			BasisPoints:      annualBP,
		},
		SingleTurnoverProfitRate: singleBP,
		AnnualSurplusValueRate:   profit.AnnualSurplusValueRate(annualSurplusBP),
		AnnualWages:              profit.LabourMinutes(wages),
		CreatedAt:                createdAt,
	}, nil
}

func scanEconomyAnalysis(s rowScanner) (profit.EconomyAnalysis, error) {
	var (
		id, kind                          string
		constant, variable, surplus, save int64
		oldBP, newBP, gainBP              int64
		createdAt                         time.Time
	)
	err := s.Scan(&id, &kind, &constant, &variable, &surplus, &save,
		&oldBP, &newBP, &gainBP, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return profit.EconomyAnalysis{}, ErrNotFound
		}
		return profit.EconomyAnalysis{}, err
	}
	return profit.EconomyAnalysis{
		ID: profit.EconomyAnalysisID(id),
		Economy: profit.ConstantCapitalEconomy{
			Kind:            profit.EconomyKind(kind),
			ConstantCapital: profit.ConstantCapital(constant),
			VariableCapital: profit.VariableCapital(variable),
			SurplusValue:    profit.SurplusValue(surplus),
			Saving:          profit.ConstantCapital(save),
		},
		OldProfitRate:  oldBP,
		NewProfitRate:  newBP,
		ProfitRateGain: gainBP,
		CreatedAt:      createdAt,
	}, nil
}

func scanPriceFluctuationAnalysis(s rowScanner) (profit.PriceFluctuationAnalysis, error) {
	var (
		id, kind                              string
		fixed, material, factor, variable, sv int64
		oldBP, newBP                          int64
		createdAt                             time.Time
	)
	err := s.Scan(&id, &kind, &fixed, &material, &factor, &variable, &sv,
		&oldBP, &newBP, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return profit.PriceFluctuationAnalysis{}, ErrNotFound
		}
		return profit.PriceFluctuationAnalysis{}, err
	}
	return profit.PriceFluctuationAnalysis{
		ID:   profit.PriceFluctuationAnalysisID(id),
		Kind: profit.PriceFluctuationKind(kind),
		Effect: profit.ConstantCapitalPriceEffect{
			FixedCapital:            profit.ConstantCapital(fixed),
			OriginalMaterialCapital: profit.ConstantCapital(material),
			PriceFactor:             factor,
			VariableCapital:         profit.VariableCapital(variable),
			SurplusValue:            profit.SurplusValue(sv),
		},
		OldProfitRate: oldBP,
		NewProfitRate: newBP,
		CreatedAt:     createdAt,
	}, nil
}

func scanCompositionEffect(s rowScanner) (profit.CompositionEffectAnalysis, error) {
	var (
		id                           string
		sa, va, ca, sb, vb, cb       int64
		profitA, profitB, difference int64
		createdAt                    time.Time
	)
	err := s.Scan(&id, &sa, &va, &ca, &sb, &vb, &cb, &profitA, &profitB, &difference, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return profit.CompositionEffectAnalysis{}, ErrNotFound
		}
		return profit.CompositionEffectAnalysis{}, err
	}
	return profit.CompositionEffectAnalysis{
		ID: profit.CompositionEffectID(id),
		Effect: profit.OrganicCompositionEffect{
			SA: profit.SurplusValue(sa),
			VA: profit.VariableCapital(va),
			CA: profit.ConstantCapital(ca),
			SB: profit.SurplusValue(sb),
			VB: profit.VariableCapital(vb),
			CB: profit.ConstantCapital(cb),
		},
		ProfitRateA:    profitA,
		ProfitRateB:    profitB,
		RateDifference: difference,
		CreatedAt:      createdAt,
	}, nil
}

func scanMagnitudeChange(s rowScanner) (profit.MagnitudeChangeAnalysis, error) {
	var (
		id, kind                        string
		originalCapital, originalProfit int64
		factor, oldRate, newRate        int64
		rateUnchanged                   bool
		createdAt                       time.Time
	)
	err := s.Scan(&id, &kind, &originalCapital, &originalProfit, &factor, &oldRate, &newRate, &rateUnchanged, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return profit.MagnitudeChangeAnalysis{}, ErrNotFound
		}
		return profit.MagnitudeChangeAnalysis{}, err
	}
	return profit.MagnitudeChangeAnalysis{
		ID: profit.MagnitudeChangeID(id),
		Change: profit.CapitalMagnitudeChange{
			Kind:            profit.MagnitudeChangeKind(kind),
			OriginalCapital: profit.TotalCapital(originalCapital),
			OriginalProfit:  profit.SurplusValue(originalProfit),
			Factor:          factor,
		},
		OldRate:       oldRate,
		NewRate:       newRate,
		RateUnchanged: rateUnchanged,
		CreatedAt:     createdAt,
	}, nil
}

func scanProductionSphere(s rowScanner) (avgprofit.ProductionSphere, error) {
	var (
		id, name                                                      string
		c, v, sRate, cc, sv, ipr, varPct, lpi                        int64
		createdAt                                                     time.Time
	)
	err := s.Scan(&id, &name, &c, &v, &sRate, &cc, &sv, &ipr, &varPct, &lpi, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return avgprofit.ProductionSphere{}, ErrNotFound
		}
		return avgprofit.ProductionSphere{}, err
	}
	return avgprofit.ProductionSphere{
		ID:                   avgprofit.ProductionSphereID(id),
		Name:                 name,
		C:                    avgprofit.TotalCapital(c),
		V:                    avgprofit.VariableCapital(v),
		SRate:                avgprofit.RateOfSurplusValue(sRate),
		ConstantCapital:      avgprofit.ConstantCapital(cc),
		SurplusValue:         avgprofit.SurplusValue(sv),
		IndividualProfitRate: avgprofit.SphereProfitRate(ipr),
		VariablePercent:      varPct,
		LabourPowerIndex:     avgprofit.LabourPowerIndex(lpi),
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
