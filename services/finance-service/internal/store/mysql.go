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
	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/merchant"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/profit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/tendency"
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
		id, name                              string
		c, v, sRate, cc, sv, ipr, varPct, lpi int64
		createdAt                             time.Time
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

// CreateGeneralProfitRate persists g, assigning an ID and timestamp when absent.
func (m *MySQL) CreateGeneralProfitRate(ctx context.Context, g avgprofit.GeneralProfitRate) (avgprofit.GeneralProfitRate, error) {
	if g.ID.IsZero() {
		g.ID = avgprofit.NewGeneralProfitRateID()
	}
	if g.CreatedAt.IsZero() {
		g.CreatedAt = m.now().UTC()
	}

	spheresJSON, err := json.Marshal(g.Spheres)
	if err != nil {
		return avgprofit.GeneralProfitRate{}, err
	}

	const q = `INSERT INTO general_profit_rates
		(id, rate, sum_surplus_values, sum_total_capitals, average_variable_percent, spheres_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = m.db.ExecContext(ctx, q,
		string(g.ID),
		int64(g.Rate), int64(g.SumSurplusValues), int64(g.SumTotalCapitals),
		g.AverageVariablePercent, string(spheresJSON),
		g.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return avgprofit.GeneralProfitRate{}, ErrAlreadyExists
		}
		return avgprofit.GeneralProfitRate{}, err
	}
	return g, nil
}

// GetGeneralProfitRate returns the general rate with id, or ErrNotFound.
func (m *MySQL) GetGeneralProfitRate(ctx context.Context, id avgprofit.GeneralProfitRateID) (avgprofit.GeneralProfitRate, error) {
	const q = `SELECT id, rate, sum_surplus_values, sum_total_capitals, average_variable_percent, spheres_json, created_at
		FROM general_profit_rates WHERE id = ?`
	return scanGeneralProfitRate(m.db.QueryRowContext(ctx, q, string(id)))
}

// CreatePriceOfProduction persists p, assigning an ID and timestamp when absent.
func (m *MySQL) CreatePriceOfProduction(ctx context.Context, p avgprofit.PriceOfProduction) (avgprofit.PriceOfProduction, error) {
	if p.ID.IsZero() {
		p.ID = avgprofit.NewPriceOfProductionID()
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO prices_of_production
		(id, sphere_name, cost_price, general_rate, commodity_value, average_profit, price, deviation, composition, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(p.ID), p.SphereName,
		int64(p.CostPrice), int64(p.GeneralRate), int64(p.CommodityValue),
		int64(p.AverageProfit), int64(p.Price), int64(p.Deviation),
		string(p.Composition), p.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return avgprofit.PriceOfProduction{}, ErrAlreadyExists
		}
		return avgprofit.PriceOfProduction{}, err
	}
	return p, nil
}

// GetPriceOfProduction returns the record with id, or ErrNotFound.
func (m *MySQL) GetPriceOfProduction(ctx context.Context, id avgprofit.PriceOfProductionID) (avgprofit.PriceOfProduction, error) {
	const q = `SELECT id, sphere_name, cost_price, general_rate, commodity_value, average_profit, price, deviation, composition, created_at
		FROM prices_of_production WHERE id = ?`
	return scanPriceOfProduction(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListPricesOfProduction returns all stored records, newest first.
func (m *MySQL) ListPricesOfProduction(ctx context.Context) ([]avgprofit.PriceOfProduction, error) {
	const q = `SELECT id, sphere_name, cost_price, general_rate, commodity_value, average_profit, price, deviation, composition, created_at
		FROM prices_of_production ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []avgprofit.PriceOfProduction
	for rows.Next() {
		p, err := scanPriceOfProduction(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func scanGeneralProfitRate(s rowScanner) (avgprofit.GeneralProfitRate, error) {
	var (
		id, spheresJSON             string
		rate, sumS, sumC, avgVarPct int64
		createdAt                   time.Time
	)
	err := s.Scan(&id, &rate, &sumS, &sumC, &avgVarPct, &spheresJSON, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return avgprofit.GeneralProfitRate{}, ErrNotFound
		}
		return avgprofit.GeneralProfitRate{}, err
	}
	var spheres []avgprofit.ProductionSphere
	if err := json.Unmarshal([]byte(spheresJSON), &spheres); err != nil {
		return avgprofit.GeneralProfitRate{}, err
	}
	if spheres == nil {
		spheres = []avgprofit.ProductionSphere{}
	}
	return avgprofit.GeneralProfitRate{
		ID:                     avgprofit.GeneralProfitRateID(id),
		Spheres:                spheres,
		Rate:                   avgprofit.SphereProfitRate(rate),
		SumSurplusValues:       avgprofit.SurplusValue(sumS),
		SumTotalCapitals:       avgprofit.TotalCapital(sumC),
		AverageVariablePercent: avgVarPct,
		CreatedAt:              createdAt,
	}, nil
}

func scanPriceOfProduction(s rowScanner) (avgprofit.PriceOfProduction, error) {
	var (
		id, sphereName, composition            string
		costPrice, generalRate, commodityValue int64
		averageProfit, price, deviation        int64
		createdAt                              time.Time
	)
	err := s.Scan(&id, &sphereName, &costPrice, &generalRate, &commodityValue,
		&averageProfit, &price, &deviation, &composition, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return avgprofit.PriceOfProduction{}, ErrNotFound
		}
		return avgprofit.PriceOfProduction{}, err
	}
	return avgprofit.PriceOfProduction{
		ID:             avgprofit.PriceOfProductionID(id),
		SphereName:     sphereName,
		CostPrice:      avgprofit.CostPrice(costPrice),
		GeneralRate:    avgprofit.SphereProfitRate(generalRate),
		CommodityValue: avgprofit.CommodityValue(commodityValue),
		AverageProfit:  avgprofit.Profit(averageProfit),
		Price:          avgprofit.ProductionPrice(price),
		Deviation:      avgprofit.ValueDeviation(deviation),
		Composition:    avgprofit.CompositionKind(composition),
		CreatedAt:      createdAt,
	}, nil
}

// CreateMarketValue persists v, assigning an ID and timestamp when absent.
func (m *MySQL) CreateMarketValue(ctx context.Context, v avgprofit.MarketValue) (avgprofit.MarketValue, error) {
	if v.ID.IsZero() {
		v.ID = avgprofit.NewMarketValueID()
	}
	if v.CreatedAt.IsZero() {
		v.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO market_values
		(id, sphere_name, bulk_condition_value, best_condition_value, worst_condition_value, value, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(v.ID), v.SphereName,
		int64(v.BulkConditionValue), int64(v.BestConditionValue), int64(v.WorstConditionValue),
		int64(v.Value), v.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return avgprofit.MarketValue{}, ErrAlreadyExists
		}
		return avgprofit.MarketValue{}, err
	}
	return v, nil
}

// GetMarketValue returns the market-value record with id, or ErrNotFound.
func (m *MySQL) GetMarketValue(ctx context.Context, id avgprofit.MarketValueID) (avgprofit.MarketValue, error) {
	const q = `SELECT id, sphere_name, bulk_condition_value, best_condition_value, worst_condition_value, value, created_at
		FROM market_values WHERE id = ?`
	return scanMarketValue(m.db.QueryRowContext(ctx, q, string(id)))
}

// CreateSurplusProfit persists s, assigning an ID and timestamp when absent.
func (m *MySQL) CreateSurplusProfit(ctx context.Context, s avgprofit.SurplusProfit) (avgprofit.SurplusProfit, error) {
	if s.ID.IsZero() {
		s.ID = avgprofit.NewSurplusProfitID()
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO surplus_profits
		(id, firm_name, individual_value, market_value, output_qty, general_rate, amount, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(s.ID), s.FirmName,
		int64(s.IndividualValue), int64(s.MarketValue), s.OutputQty,
		int64(s.GeneralRate), int64(s.Amount), s.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return avgprofit.SurplusProfit{}, ErrAlreadyExists
		}
		return avgprofit.SurplusProfit{}, err
	}
	return s, nil
}

// GetSurplusProfit returns the surplus-profit record with id, or ErrNotFound.
func (m *MySQL) GetSurplusProfit(ctx context.Context, id avgprofit.SurplusProfitID) (avgprofit.SurplusProfit, error) {
	const q = `SELECT id, firm_name, individual_value, market_value, output_qty, general_rate, amount, created_at
		FROM surplus_profits WHERE id = ?`
	return scanSurplusProfit(m.db.QueryRowContext(ctx, q, string(id)))
}

// CreateEqualisation persists e, assigning an ID and timestamp when absent.
func (m *MySQL) CreateEqualisation(ctx context.Context, e avgprofit.Equalisation) (avgprofit.Equalisation, error) {
	if e.ID.IsZero() {
		e.ID = avgprofit.NewEqualisationID()
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO equalisations
		(id, sphere_name, initial_rate, target_rate, direction, is_converging, market_price, market_value, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(e.ID), e.SphereName,
		int64(e.InitialRate), int64(e.TargetRate),
		string(e.Direction), e.IsConverging,
		int64(e.Flow.MarketPrice), int64(e.Flow.MarketValue),
		e.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return avgprofit.Equalisation{}, ErrAlreadyExists
		}
		return avgprofit.Equalisation{}, err
	}
	return e, nil
}

// GetEqualisation returns the equalisation record with id, or ErrNotFound.
func (m *MySQL) GetEqualisation(ctx context.Context, id avgprofit.EqualisationID) (avgprofit.Equalisation, error) {
	const q = `SELECT id, sphere_name, initial_rate, target_rate, direction, is_converging, market_price, market_value, created_at
		FROM equalisations WHERE id = ?`
	return scanEqualisation(m.db.QueryRowContext(ctx, q, string(id)))
}

func scanMarketValue(s rowScanner) (avgprofit.MarketValue, error) {
	var (
		id, sphereName           string
		bulk, best, worst, value int64
		createdAt                time.Time
	)
	err := s.Scan(&id, &sphereName, &bulk, &best, &worst, &value, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return avgprofit.MarketValue{}, ErrNotFound
		}
		return avgprofit.MarketValue{}, err
	}
	return avgprofit.MarketValue{
		ID:                  avgprofit.MarketValueID(id),
		SphereName:          sphereName,
		BulkConditionValue:  avgprofit.CommodityValue(bulk),
		BestConditionValue:  avgprofit.CommodityValue(best),
		WorstConditionValue: avgprofit.CommodityValue(worst),
		Value:               avgprofit.CommodityValue(value),
		CreatedAt:           createdAt,
	}, nil
}

func scanSurplusProfit(s rowScanner) (avgprofit.SurplusProfit, error) {
	var (
		id, firmName                                      string
		indVal, marketVal, outputQty, generalRate, amount int64
		createdAt                                         time.Time
	)
	err := s.Scan(&id, &firmName, &indVal, &marketVal, &outputQty, &generalRate, &amount, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return avgprofit.SurplusProfit{}, ErrNotFound
		}
		return avgprofit.SurplusProfit{}, err
	}
	return avgprofit.SurplusProfit{
		ID:              avgprofit.SurplusProfitID(id),
		FirmName:        firmName,
		IndividualValue: avgprofit.IndividualValue(indVal),
		MarketValue:     avgprofit.CommodityValue(marketVal),
		OutputQty:       outputQty,
		GeneralRate:     avgprofit.SphereProfitRate(generalRate),
		Amount:          avgprofit.Profit(amount),
		CreatedAt:       createdAt,
	}, nil
}

func scanEqualisation(s rowScanner) (avgprofit.Equalisation, error) {
	var (
		id, sphereName, direction string
		initialRate, targetRate   int64
		isConverging              bool
		marketPrice, marketValue  int64
		createdAt                 time.Time
	)
	err := s.Scan(&id, &sphereName, &initialRate, &targetRate, &direction, &isConverging, &marketPrice, &marketValue, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return avgprofit.Equalisation{}, ErrNotFound
		}
		return avgprofit.Equalisation{}, err
	}
	flow := avgprofit.ComputeCapitalFlow(
		sphereName,
		avgprofit.SphereProfitRate(initialRate),
		avgprofit.SphereProfitRate(targetRate),
		avgprofit.MarketPriceAmount(marketPrice),
		avgprofit.CommodityValue(marketValue),
	)
	return avgprofit.Equalisation{
		ID:           avgprofit.EqualisationID(id),
		SphereName:   sphereName,
		InitialRate:  avgprofit.SphereProfitRate(initialRate),
		TargetRate:   avgprofit.SphereProfitRate(targetRate),
		Direction:    avgprofit.CapitalFlowDirection(direction),
		IsConverging: isConverging,
		Flow:         flow,
		CreatedAt:    createdAt,
	}, nil
}

// CreateWageEffectAnalysis persists a, assigning an ID and timestamp when absent.
func (m *MySQL) CreateWageEffectAnalysis(ctx context.Context, a avgprofit.WageEffectAnalysis) (avgprofit.WageEffectAnalysis, error) {
	if a.ID.IsZero() {
		a.ID = avgprofit.NewWageEffectAnalysisID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	if a.Outcomes == nil {
		a.Outcomes = []avgprofit.SphereWageOutcome{}
	}

	outcomesJSON, err := json.Marshal(a.Outcomes)
	if err != nil {
		return avgprofit.WageEffectAnalysis{}, err
	}
	avgOutcomeJSON, err := json.Marshal(a.AverageOutcome)
	if err != nil {
		return avgprofit.WageEffectAnalysis{}, err
	}

	const q = `INSERT INTO wage_effect_analyses
		(id, base_constant, base_variable, s_rate, wage_factor, old_general_rate, new_general_rate,
		 kind, average_outcome_json, outcomes_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = m.db.ExecContext(ctx, q,
		string(a.ID),
		int64(a.BaseConstant), int64(a.BaseVariable), int64(a.SRate),
		a.WageFactor, int64(a.OldGeneralRate), int64(a.NewGeneralRate),
		string(a.Kind), string(avgOutcomeJSON), string(outcomesJSON),
		a.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return avgprofit.WageEffectAnalysis{}, ErrAlreadyExists
		}
		return avgprofit.WageEffectAnalysis{}, err
	}
	return a, nil
}

// GetWageEffectAnalysis returns the analysis with id, or ErrNotFound.
func (m *MySQL) GetWageEffectAnalysis(ctx context.Context, id avgprofit.WageEffectAnalysisID) (avgprofit.WageEffectAnalysis, error) {
	const q = `SELECT id, base_constant, base_variable, s_rate, wage_factor,
		old_general_rate, new_general_rate, kind, average_outcome_json, outcomes_json, created_at
		FROM wage_effect_analyses WHERE id = ?`
	return scanWageEffectAnalysis(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListWageEffectAnalyses returns all stored analyses, newest first.
func (m *MySQL) ListWageEffectAnalyses(ctx context.Context) ([]avgprofit.WageEffectAnalysis, error) {
	const q = `SELECT id, base_constant, base_variable, s_rate, wage_factor,
		old_general_rate, new_general_rate, kind, average_outcome_json, outcomes_json, created_at
		FROM wage_effect_analyses ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []avgprofit.WageEffectAnalysis
	for rows.Next() {
		a, err := scanWageEffectAnalysis(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func scanWageEffectAnalysis(s rowScanner) (avgprofit.WageEffectAnalysis, error) {
	var (
		id, kind, avgOutcomeJSON, outcomesJSON string
		baseC, baseV, sRate, wageFactor        int64
		oldRate, newRate                       int64
		createdAt                              time.Time
	)
	err := s.Scan(&id, &baseC, &baseV, &sRate, &wageFactor,
		&oldRate, &newRate, &kind, &avgOutcomeJSON, &outcomesJSON, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return avgprofit.WageEffectAnalysis{}, ErrNotFound
		}
		return avgprofit.WageEffectAnalysis{}, err
	}

	var avgOutcome avgprofit.SphereWageOutcome
	if err := json.Unmarshal([]byte(avgOutcomeJSON), &avgOutcome); err != nil {
		return avgprofit.WageEffectAnalysis{}, err
	}

	var outcomes []avgprofit.SphereWageOutcome
	if err := json.Unmarshal([]byte(outcomesJSON), &outcomes); err != nil {
		return avgprofit.WageEffectAnalysis{}, err
	}
	if outcomes == nil {
		outcomes = []avgprofit.SphereWageOutcome{}
	}

	return avgprofit.WageEffectAnalysis{
		ID:             avgprofit.WageEffectAnalysisID(id),
		BaseConstant:   avgprofit.ConstantCapital(baseC),
		BaseVariable:   avgprofit.VariableCapital(baseV),
		SRate:          avgprofit.RateOfSurplusValue(sRate),
		WageFactor:     wageFactor,
		OldGeneralRate: avgprofit.SphereProfitRate(oldRate),
		NewGeneralRate: avgprofit.SphereProfitRate(newRate),
		Kind:           avgprofit.WageFluctuationKind(kind),
		AverageOutcome: avgOutcome,
		Outcomes:       outcomes,
		CreatedAt:      createdAt,
	}, nil
}

// CreatePriceOfProductionChange persists c, assigning an ID and timestamp when absent.
func (m *MySQL) CreatePriceOfProductionChange(ctx context.Context, c avgprofit.PriceOfProductionChange) (avgprofit.PriceOfProductionChange, error) {
	if c.ID.IsZero() {
		c.ID = avgprofit.NewPriceOfProductionChangeID()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO price_of_production_changes
		(id, sphere_name, cause, rate_changed, value_changed, price_changed, old_price, new_price, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID), c.SphereName, string(c.Cause),
		c.RateChanged, c.ValueChanged, c.PriceChanged,
		int64(c.OldPrice), int64(c.NewPrice), c.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return avgprofit.PriceOfProductionChange{}, ErrAlreadyExists
		}
		return avgprofit.PriceOfProductionChange{}, err
	}
	return c, nil
}

// GetPriceOfProductionChange returns the change with id, or ErrNotFound.
func (m *MySQL) GetPriceOfProductionChange(ctx context.Context, id avgprofit.PriceOfProductionChangeID) (avgprofit.PriceOfProductionChange, error) {
	const q = `SELECT id, sphere_name, cause, rate_changed, value_changed, price_changed, old_price, new_price, created_at
		FROM price_of_production_changes WHERE id = ?`
	return scanPriceOfProductionChange(m.db.QueryRowContext(ctx, q, string(id)))
}

func scanPriceOfProductionChange(s rowScanner) (avgprofit.PriceOfProductionChange, error) {
	var (
		id, sphereName, cause                   string
		rateChanged, valueChanged, priceChanged bool
		oldPrice, newPrice                      int64
		createdAt                               time.Time
	)
	err := s.Scan(&id, &sphereName, &cause, &rateChanged, &valueChanged, &priceChanged,
		&oldPrice, &newPrice, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return avgprofit.PriceOfProductionChange{}, ErrNotFound
		}
		return avgprofit.PriceOfProductionChange{}, err
	}
	return avgprofit.PriceOfProductionChange{
		ID:           avgprofit.PriceOfProductionChangeID(id),
		SphereName:   sphereName,
		Cause:        avgprofit.PriceChangeCause(cause),
		RateChanged:  rateChanged,
		ValueChanged: valueChanged,
		PriceChanged: priceChanged,
		OldPrice:     avgprofit.ProductionPrice(oldPrice),
		NewPrice:     avgprofit.ProductionPrice(newPrice),
		CreatedAt:    createdAt,
	}, nil
}

// CreateCompositionTrajectory persists t, assigning an ID and timestamp when
// absent. Periods are stored as JSON TEXT; ProfitRates is derived on read.
func (m *MySQL) CreateCompositionTrajectory(ctx context.Context, t tendency.CompositionTrajectory) (tendency.CompositionTrajectory, error) {
	if t.ID.IsZero() {
		t.ID = tendency.NewCompositionTrajectoryID()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = m.now().UTC()
	}
	if t.Periods == nil {
		t.Periods = []tendency.TrajectoryPeriod{}
	}
	t.ProfitRates = t.DeriveProfitRates()

	periodsJSON, err := json.Marshal(t.Periods)
	if err != nil {
		return tendency.CompositionTrajectory{}, err
	}

	const q = `INSERT INTO composition_trajectories
		(id, label, surplus_value_rate, periods_json, created_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err = m.db.ExecContext(ctx, q,
		string(t.ID), t.Label, t.SurplusValueRate, string(periodsJSON), t.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return tendency.CompositionTrajectory{}, ErrAlreadyExists
		}
		return tendency.CompositionTrajectory{}, err
	}
	return t, nil
}

// GetCompositionTrajectory returns the trajectory with id, or ErrNotFound.
func (m *MySQL) GetCompositionTrajectory(ctx context.Context, id tendency.CompositionTrajectoryID) (tendency.CompositionTrajectory, error) {
	const q = `SELECT id, label, surplus_value_rate, periods_json, created_at
		FROM composition_trajectories WHERE id = ?`
	return scanCompositionTrajectory(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListCompositionTrajectories returns all stored trajectories, newest first.
func (m *MySQL) ListCompositionTrajectories(ctx context.Context) ([]tendency.CompositionTrajectory, error) {
	const q = `SELECT id, label, surplus_value_rate, periods_json, created_at
		FROM composition_trajectories ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]tendency.CompositionTrajectory, 0)
	for rows.Next() {
		t, err := scanCompositionTrajectory(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func scanCompositionTrajectory(s rowScanner) (tendency.CompositionTrajectory, error) {
	var (
		id, label, periodsJSON string
		surplusValueRate       int64
		createdAt              time.Time
	)
	err := s.Scan(&id, &label, &surplusValueRate, &periodsJSON, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tendency.CompositionTrajectory{}, ErrNotFound
		}
		return tendency.CompositionTrajectory{}, err
	}
	var periods []tendency.TrajectoryPeriod
	if err := json.Unmarshal([]byte(periodsJSON), &periods); err != nil {
		return tendency.CompositionTrajectory{}, err
	}
	if periods == nil {
		periods = []tendency.TrajectoryPeriod{}
	}
	t := tendency.CompositionTrajectory{
		ID:               tendency.CompositionTrajectoryID(id),
		Label:            label,
		SurplusValueRate: surplusValueRate,
		Periods:          periods,
		CreatedAt:        createdAt,
	}
	t.ProfitRates = t.DeriveProfitRates()
	return t, nil
}

// CreateRateMassContradiction persists r, assigning an ID and timestamp when absent.
func (m *MySQL) CreateRateMassContradiction(ctx context.Context, r tendency.RateMassContradiction) (tendency.RateMassContradiction, error) {
	if r.ID.IsZero() {
		r.ID = tendency.NewRateMassContradictionID()
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO rate_mass_contradictions
		(id, old_c, old_rate, new_c, new_rate, old_mass, new_mass, mass_change, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(r.ID), r.OldC, r.OldRate, r.NewC, r.NewRate,
		r.OldMass, r.NewMass, r.MassChange, r.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return tendency.RateMassContradiction{}, ErrAlreadyExists
		}
		return tendency.RateMassContradiction{}, err
	}
	return r, nil
}

// GetRateMassContradiction returns the record with id, or ErrNotFound.
func (m *MySQL) GetRateMassContradiction(ctx context.Context, id tendency.RateMassContradictionID) (tendency.RateMassContradiction, error) {
	const q = `SELECT id, old_c, old_rate, new_c, new_rate, old_mass, new_mass, mass_change, created_at
		FROM rate_mass_contradictions WHERE id = ?`
	return scanRateMassContradiction(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListRateMassContradictions returns all stored records, newest first.
func (m *MySQL) ListRateMassContradictions(ctx context.Context) ([]tendency.RateMassContradiction, error) {
	const q = `SELECT id, old_c, old_rate, new_c, new_rate, old_mass, new_mass, mass_change, created_at
		FROM rate_mass_contradictions ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]tendency.RateMassContradiction, 0)
	for rows.Next() {
		r, err := scanRateMassContradiction(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func scanRateMassContradiction(s rowScanner) (tendency.RateMassContradiction, error) {
	var (
		id                                                       string
		oldC, oldRate, newC, newRate, oldMass, newMass, massChng int64
		createdAt                                                time.Time
	)
	err := s.Scan(&id, &oldC, &oldRate, &newC, &newRate, &oldMass, &newMass, &massChng, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tendency.RateMassContradiction{}, ErrNotFound
		}
		return tendency.RateMassContradiction{}, err
	}
	return tendency.RateMassContradiction{
		ID:         tendency.RateMassContradictionID(id),
		OldC:       oldC,
		OldRate:    oldRate,
		NewC:       newC,
		NewRate:    newRate,
		OldMass:    oldMass,
		NewMass:    newMass,
		MassChange: massChng,
		CreatedAt:  createdAt,
	}, nil
}

// CreateCounteractingForce persists f, assigning an ID and timestamp when absent.
func (m *MySQL) CreateCounteractingForce(ctx context.Context, f tendency.CounteractingForce) (tendency.CounteractingForce, error) {
	if f.ID.IsZero() {
		f.ID = tendency.NewCounteractingForceID()
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO counteracting_forces
		(id, kind, excluded_from_general_rate, note, created_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(f.ID), string(f.Kind), f.ExcludedFromGeneralRate, f.Note, f.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return tendency.CounteractingForce{}, ErrAlreadyExists
		}
		return tendency.CounteractingForce{}, err
	}
	return f, nil
}

// GetCounteractingForce returns the force with id, or ErrNotFound.
func (m *MySQL) GetCounteractingForce(ctx context.Context, id tendency.CounteractingForceID) (tendency.CounteractingForce, error) {
	const q = `SELECT id, kind, excluded_from_general_rate, note, created_at
		FROM counteracting_forces WHERE id = ?`
	return scanCounteractingForce(m.db.QueryRowContext(ctx, q, string(id)))
}

func scanCounteractingForce(s rowScanner) (tendency.CounteractingForce, error) {
	var (
		id, kind  string
		excluded  bool
		note      string
		createdAt time.Time
	)
	err := s.Scan(&id, &kind, &excluded, &note, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tendency.CounteractingForce{}, ErrNotFound
		}
		return tendency.CounteractingForce{}, err
	}
	return tendency.CounteractingForce{
		ID:                      tendency.CounteractingForceID(id),
		Kind:                    tendency.CounteractingForceKind(kind),
		ExcludedFromGeneralRate: excluded,
		Note:                    note,
		CreatedAt:               createdAt,
	}, nil
}

// CreateCounteractingScenario persists s, assigning an ID and timestamp when
// absent. Forces, the base trajectory, and the modified rates are stored as JSON
// TEXT; the domain layer owns serialisation.
func (m *MySQL) CreateCounteractingScenario(ctx context.Context, s tendency.CounteractingScenario) (tendency.CounteractingScenario, error) {
	if s.ID.IsZero() {
		s.ID = tendency.NewCounteractingScenarioID()
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	if s.Forces == nil {
		s.Forces = []tendency.CounteractingForce{}
	}
	if s.ModifiedRates == nil {
		s.ModifiedRates = []int64{}
	}

	forcesJSON, err := json.Marshal(s.Forces)
	if err != nil {
		return tendency.CounteractingScenario{}, err
	}
	trajectoryJSON, err := json.Marshal(s.BaseTrajectory)
	if err != nil {
		return tendency.CounteractingScenario{}, err
	}
	ratesJSON, err := json.Marshal(s.ModifiedRates)
	if err != nil {
		return tendency.CounteractingScenario{}, err
	}

	const q = `INSERT INTO counteracting_scenarios
		(id, label, forces_json, base_trajectory_json, initial_profit_rate, final_profit_rate, modified_rates_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err = m.db.ExecContext(ctx, q,
		string(s.ID), s.Label, string(forcesJSON), string(trajectoryJSON),
		s.InitialProfitRate, s.FinalProfitRate, string(ratesJSON), s.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return tendency.CounteractingScenario{}, ErrAlreadyExists
		}
		return tendency.CounteractingScenario{}, err
	}
	return s, nil
}

// GetCounteractingScenario returns the scenario with id, or ErrNotFound.
func (m *MySQL) GetCounteractingScenario(ctx context.Context, id tendency.CounteractingScenarioID) (tendency.CounteractingScenario, error) {
	const q = `SELECT id, label, forces_json, base_trajectory_json, initial_profit_rate, final_profit_rate, modified_rates_json, created_at
		FROM counteracting_scenarios WHERE id = ?`
	return scanCounteractingScenario(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListCounteractingScenarios returns all stored scenarios, newest first.
func (m *MySQL) ListCounteractingScenarios(ctx context.Context) ([]tendency.CounteractingScenario, error) {
	const q = `SELECT id, label, forces_json, base_trajectory_json, initial_profit_rate, final_profit_rate, modified_rates_json, created_at
		FROM counteracting_scenarios ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]tendency.CounteractingScenario, 0)
	for rows.Next() {
		s, err := scanCounteractingScenario(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func scanCounteractingScenario(s rowScanner) (tendency.CounteractingScenario, error) {
	var (
		id, label                             string
		forcesJSON, trajectoryJSON, ratesJSON string
		initialRate, finalRate                int64
		createdAt                             time.Time
	)
	err := s.Scan(&id, &label, &forcesJSON, &trajectoryJSON, &initialRate, &finalRate, &ratesJSON, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tendency.CounteractingScenario{}, ErrNotFound
		}
		return tendency.CounteractingScenario{}, err
	}

	var forces []tendency.CounteractingForce
	if err := json.Unmarshal([]byte(forcesJSON), &forces); err != nil {
		return tendency.CounteractingScenario{}, err
	}
	if forces == nil {
		forces = []tendency.CounteractingForce{}
	}

	var trajectory tendency.CompositionTrajectory
	if err := json.Unmarshal([]byte(trajectoryJSON), &trajectory); err != nil {
		return tendency.CounteractingScenario{}, err
	}

	var rates []int64
	if err := json.Unmarshal([]byte(ratesJSON), &rates); err != nil {
		return tendency.CounteractingScenario{}, err
	}
	if rates == nil {
		rates = []int64{}
	}

	return tendency.CounteractingScenario{
		ID:                tendency.CounteractingScenarioID(id),
		Label:             label,
		Forces:            forces,
		BaseTrajectory:    trajectory,
		InitialProfitRate: initialRate,
		FinalProfitRate:   finalRate,
		ModifiedRates:     rates,
		CreatedAt:         createdAt,
	}, nil
}

// CreateCrisis persists c, assigning an ID and timestamp when absent. The
// derived post-crisis rate is stored as a scalar column for query efficiency.
func (m *MySQL) CreateCrisis(ctx context.Context, c tendency.Crisis) (tendency.Crisis, error) {
	if c.ID.IsZero() {
		c.ID = tendency.NewCrisisID()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO crises
		(id, constant_capital_writedown, pre_crisis_profit_rate, post_crisis_profit_rate, created_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID), c.ConstantCapitalWritedown, c.PreCrisisProfitRate, c.PostCrisisProfitRate, c.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return tendency.Crisis{}, ErrAlreadyExists
		}
		return tendency.Crisis{}, err
	}
	return c, nil
}

// GetCrisis returns the crisis with id, or ErrNotFound.
func (m *MySQL) GetCrisis(ctx context.Context, id tendency.CrisisID) (tendency.Crisis, error) {
	const q = `SELECT id, constant_capital_writedown, pre_crisis_profit_rate, post_crisis_profit_rate, created_at
		FROM crises WHERE id = ?`
	return scanCrisis(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListCrises returns all stored crises, newest first.
func (m *MySQL) ListCrises(ctx context.Context) ([]tendency.Crisis, error) {
	const q = `SELECT id, constant_capital_writedown, pre_crisis_profit_rate, post_crisis_profit_rate, created_at
		FROM crises ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]tendency.Crisis, 0)
	for rows.Next() {
		c, err := scanCrisis(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func scanCrisis(s rowScanner) (tendency.Crisis, error) {
	var (
		id                                       string
		writedown, preCrisisRate, postCrisisRate int64
		createdAt                                time.Time
	)
	err := s.Scan(&id, &writedown, &preCrisisRate, &postCrisisRate, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tendency.Crisis{}, ErrNotFound
		}
		return tendency.Crisis{}, err
	}
	return tendency.Crisis{
		ID:                       tendency.CrisisID(id),
		ConstantCapitalWritedown: writedown,
		PreCrisisProfitRate:      preCrisisRate,
		PostCrisisProfitRate:     postCrisisRate,
		CreatedAt:                createdAt,
	}, nil
}

// CreateInternalContradiction persists c, assigning an ID and timestamp when absent.
func (m *MySQL) CreateInternalContradiction(ctx context.Context, c tendency.InternalContradiction) (tendency.InternalContradiction, error) {
	if c.ID.IsZero() {
		c.ID = tendency.NewInternalContradictionID()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO internal_contradictions
		(id, kind, is_coexistent, note, created_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID), string(c.Kind), c.IsCoexistent, c.Note, c.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return tendency.InternalContradiction{}, ErrAlreadyExists
		}
		return tendency.InternalContradiction{}, err
	}
	return c, nil
}

// GetInternalContradiction returns the contradiction with id, or ErrNotFound.
func (m *MySQL) GetInternalContradiction(ctx context.Context, id tendency.InternalContradictionID) (tendency.InternalContradiction, error) {
	const q = `SELECT id, kind, is_coexistent, note, created_at
		FROM internal_contradictions WHERE id = ?`
	return scanInternalContradiction(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListInternalContradictions returns all stored contradictions, newest first.
func (m *MySQL) ListInternalContradictions(ctx context.Context) ([]tendency.InternalContradiction, error) {
	const q = `SELECT id, kind, is_coexistent, note, created_at
		FROM internal_contradictions ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]tendency.InternalContradiction, 0)
	for rows.Next() {
		c, err := scanInternalContradiction(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func scanInternalContradiction(s rowScanner) (tendency.InternalContradiction, error) {
	var (
		id, kind     string
		isCoexistent bool
		note         string
		createdAt    time.Time
	)
	err := s.Scan(&id, &kind, &isCoexistent, &note, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tendency.InternalContradiction{}, ErrNotFound
		}
		return tendency.InternalContradiction{}, err
	}
	return tendency.InternalContradiction{
		ID:           tendency.InternalContradictionID(id),
		Kind:         tendency.ContradictionKind(kind),
		IsCoexistent: isCoexistent,
		Note:         note,
		CreatedAt:    createdAt,
	}, nil
}

// CreateCommercialCapital persists cc, assigning an ID and timestamp when absent.
func (m *MySQL) CreateCommercialCapital(ctx context.Context, cc merchant.CommercialCapital) (merchant.CommercialCapital, error) {
	if cc.ID.IsZero() {
		cc.ID = merchant.NewCommercialCapitalID()
	}
	if cc.CreatedAt.IsZero() {
		cc.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO commercial_capitals
		(id, money_advanced, commodity_description, ` + "`function`" + `, surplus_value_produced, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(cc.ID), cc.MoneyAdvanced, cc.CommodityDescription,
		int(cc.Function), cc.SurplusValueProduced, cc.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return merchant.CommercialCapital{}, ErrAlreadyExists
		}
		return merchant.CommercialCapital{}, err
	}
	return cc, nil
}

// GetCommercialCapital returns the commercial-capital record with id, or ErrNotFound.
func (m *MySQL) GetCommercialCapital(ctx context.Context, id merchant.CommercialCapitalID) (merchant.CommercialCapital, error) {
	const q = `SELECT id, money_advanced, commodity_description, ` + "`function`" + `, surplus_value_produced, created_at
		FROM commercial_capitals WHERE id = ?`
	return scanCommercialCapital(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListCommercialCapitals returns all stored commercial-capital records, newest first.
func (m *MySQL) ListCommercialCapitals(ctx context.Context) ([]merchant.CommercialCapital, error) {
	const q = `SELECT id, money_advanced, commodity_description, ` + "`function`" + `, surplus_value_produced, created_at
		FROM commercial_capitals ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]merchant.CommercialCapital, 0)
	for rows.Next() {
		cc, err := scanCommercialCapital(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, cc)
	}
	return out, rows.Err()
}

func scanCommercialCapital(s rowScanner) (merchant.CommercialCapital, error) {
	var (
		id, commodityDescription string
		moneyAdvanced            int64
		function                 int
		surplusValueProduced     int64
		createdAt                time.Time
	)
	err := s.Scan(&id, &moneyAdvanced, &commodityDescription, &function, &surplusValueProduced, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return merchant.CommercialCapital{}, ErrNotFound
		}
		return merchant.CommercialCapital{}, err
	}
	return merchant.CommercialCapital{
		ID:                   merchant.CommercialCapitalID(id),
		MoneyAdvanced:        moneyAdvanced,
		CommodityDescription: commodityDescription,
		Function:             merchant.CommercialCapitalFunction(function),
		SurplusValueProduced: surplusValueProduced,
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
