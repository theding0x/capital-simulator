package store

import (
	"context"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"time"

	pkgmysql "github.com/theding0x/capital-simulator/pkg/mysql"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/avgprofit"
	"github.com/theding0x/capital-simulator/services/finance-service/internal/credit"
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

// CreateCommercialProfit persists cp, assigning an ID and timestamp when absent.
func (m *MySQL) CreateCommercialProfit(ctx context.Context, cp merchant.CommercialProfit) (merchant.CommercialProfit, error) {
	if cp.ID.IsZero() {
		cp.ID = merchant.NewCommercialProfitID()
	}
	if cp.CreatedAt.IsZero() {
		cp.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO commercial_profits
		(id, productive_capital, commercial_capital, surplus_value, unadjusted_rate_bp, adjusted_rate_bp, new_value_created, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(cp.ID),
		cp.ProductiveCapital, cp.CommercialCapital, cp.SurplusValue,
		cp.UnadjustedRateBP, cp.AdjustedRateBP, cp.NewValueCreated,
		cp.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return merchant.CommercialProfit{}, ErrAlreadyExists
		}
		return merchant.CommercialProfit{}, err
	}
	return cp, nil
}

// GetCommercialProfit returns the commercial-profit record with id, or ErrNotFound.
func (m *MySQL) GetCommercialProfit(ctx context.Context, id merchant.CommercialProfitID) (merchant.CommercialProfit, error) {
	const q = `SELECT id, productive_capital, commercial_capital, surplus_value, unadjusted_rate_bp, adjusted_rate_bp, new_value_created, created_at
		FROM commercial_profits WHERE id = ?`
	return scanCommercialProfit(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListCommercialProfits returns all stored commercial-profit records, newest first.
func (m *MySQL) ListCommercialProfits(ctx context.Context) ([]merchant.CommercialProfit, error) {
	const q = `SELECT id, productive_capital, commercial_capital, surplus_value, unadjusted_rate_bp, adjusted_rate_bp, new_value_created, created_at
		FROM commercial_profits ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]merchant.CommercialProfit, 0)
	for rows.Next() {
		cp, err := scanCommercialProfit(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, cp)
	}
	return out, rows.Err()
}

func scanCommercialProfit(s rowScanner) (merchant.CommercialProfit, error) {
	var (
		id                                                 string
		productiveCapital, commercialCapital, surplusValue int64
		unadjustedRateBP, adjustedRateBP, newValueCreated  int64
		createdAt                                          time.Time
	)
	err := s.Scan(&id, &productiveCapital, &commercialCapital, &surplusValue,
		&unadjustedRateBP, &adjustedRateBP, &newValueCreated, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return merchant.CommercialProfit{}, ErrNotFound
		}
		return merchant.CommercialProfit{}, err
	}
	return merchant.CommercialProfit{
		ID:                merchant.CommercialProfitID(id),
		ProductiveCapital: productiveCapital,
		CommercialCapital: commercialCapital,
		SurplusValue:      surplusValue,
		UnadjustedRateBP:  unadjustedRateBP,
		AdjustedRateBP:    adjustedRateBP,
		NewValueCreated:   newValueCreated,
		CreatedAt:         createdAt,
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

// CreateMerchantTurnover persists mt, assigning an ID and timestamp when absent.
func (m *MySQL) CreateMerchantTurnover(ctx context.Context, mt merchant.MerchantTurnover) (merchant.MerchantTurnover, error) {
	if mt.ID.IsZero() {
		mt.ID = merchant.NewMerchantTurnoverID()
	}
	if mt.CreatedAt.IsZero() {
		mt.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO merchant_turnovers
		(id, money_advanced, general_rate_bp, turnover_count, units_per_turnover, annual_merchant_profit, markup_per_unit, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(mt.ID),
		mt.MoneyAdvanced, mt.GeneralRateBP, mt.TurnoverCount, mt.UnitsPerTurnover,
		mt.AnnualMerchantProfit, mt.MarkupPerUnit,
		mt.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return merchant.MerchantTurnover{}, ErrAlreadyExists
		}
		return merchant.MerchantTurnover{}, err
	}
	return mt, nil
}

// GetMerchantTurnover returns the merchant-turnover record with id, or ErrNotFound.
func (m *MySQL) GetMerchantTurnover(ctx context.Context, id merchant.MerchantTurnoverID) (merchant.MerchantTurnover, error) {
	const q = `SELECT id, money_advanced, general_rate_bp, turnover_count, units_per_turnover, annual_merchant_profit, markup_per_unit, created_at
		FROM merchant_turnovers WHERE id = ?`
	return scanMerchantTurnover(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListMerchantTurnovers returns all stored merchant-turnover records, newest first.
func (m *MySQL) ListMerchantTurnovers(ctx context.Context) ([]merchant.MerchantTurnover, error) {
	const q = `SELECT id, money_advanced, general_rate_bp, turnover_count, units_per_turnover, annual_merchant_profit, markup_per_unit, created_at
		FROM merchant_turnovers ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]merchant.MerchantTurnover, 0)
	for rows.Next() {
		mt, err := scanMerchantTurnover(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, mt)
	}
	return out, rows.Err()
}

func scanMerchantTurnover(s rowScanner) (merchant.MerchantTurnover, error) {
	var (
		id                                                            string
		moneyAdvanced, generalRateBP, turnoverCount, unitsPerTurnover int64
		annualMerchantProfit, markupPerUnit                           int64
		createdAt                                                     time.Time
	)
	err := s.Scan(&id, &moneyAdvanced, &generalRateBP, &turnoverCount, &unitsPerTurnover,
		&annualMerchantProfit, &markupPerUnit, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return merchant.MerchantTurnover{}, ErrNotFound
		}
		return merchant.MerchantTurnover{}, err
	}
	return merchant.MerchantTurnover{
		ID:                   merchant.MerchantTurnoverID(id),
		MoneyAdvanced:        moneyAdvanced,
		GeneralRateBP:        generalRateBP,
		TurnoverCount:        turnoverCount,
		UnitsPerTurnover:     unitsPerTurnover,
		AnnualMerchantProfit: annualMerchantProfit,
		MarkupPerUnit:        markupPerUnit,
		CreatedAt:            createdAt,
	}, nil
}

// ── Vol. III Ch. 19 — Money-Dealing Capital ──────────────────────────────────

// CreateMoneyDealingCapital persists md, assigning an ID and timestamp when absent.
func (m *MySQL) CreateMoneyDealingCapital(ctx context.Context, md merchant.MoneyDealingCapital) (merchant.MoneyDealingCapital, error) {
	if md.ID.IsZero() {
		md.ID = merchant.NewMoneyDealingCapitalID()
	}
	if md.CreatedAt.IsZero() {
		md.CreatedAt = m.now().UTC()
	}
	kindsJSON, err := json.Marshal(md.OperationKinds)
	if err != nil {
		return merchant.MoneyDealingCapital{}, err
	}
	const q = `INSERT INTO money_dealing_capitals
		(id, money_advanced, general_rate_bp, operation_kinds, money_dealing_profit, new_value_created, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = m.db.ExecContext(ctx, q,
		string(md.ID),
		md.MoneyAdvanced,
		md.GeneralRateBP,
		string(kindsJSON),
		md.MoneyDealingProfit,
		md.NewValueCreated,
		md.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return merchant.MoneyDealingCapital{}, ErrAlreadyExists
		}
		return merchant.MoneyDealingCapital{}, err
	}
	return md, nil
}

// GetMoneyDealingCapital returns the money-dealing-capital record with id, or ErrNotFound.
func (m *MySQL) GetMoneyDealingCapital(ctx context.Context, id merchant.MoneyDealingCapitalID) (merchant.MoneyDealingCapital, error) {
	const q = `SELECT id, money_advanced, general_rate_bp, operation_kinds, money_dealing_profit, new_value_created, created_at
		FROM money_dealing_capitals WHERE id = ?`
	return scanMoneyDealingCapital(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListMoneyDealingCapitals returns all stored money-dealing-capital records, newest first.
func (m *MySQL) ListMoneyDealingCapitals(ctx context.Context) ([]merchant.MoneyDealingCapital, error) {
	const q = `SELECT id, money_advanced, general_rate_bp, operation_kinds, money_dealing_profit, new_value_created, created_at
		FROM money_dealing_capitals ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]merchant.MoneyDealingCapital, 0)
	for rows.Next() {
		md, err := scanMoneyDealingCapital(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, md)
	}
	return out, rows.Err()
}

func scanMoneyDealingCapital(s rowScanner) (merchant.MoneyDealingCapital, error) {
	var (
		id                                  string
		moneyAdvanced, generalRateBP        int64
		operationKindsStr                   string
		moneyDealingProfit, newValueCreated int64
		createdAt                           time.Time
	)
	err := s.Scan(&id, &moneyAdvanced, &generalRateBP, &operationKindsStr,
		&moneyDealingProfit, &newValueCreated, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return merchant.MoneyDealingCapital{}, ErrNotFound
		}
		return merchant.MoneyDealingCapital{}, err
	}
	var kinds []merchant.MoneyOperationKind
	if err := json.Unmarshal([]byte(operationKindsStr), &kinds); err != nil {
		return merchant.MoneyDealingCapital{}, fmt.Errorf("store: money_dealing_capitals: unmarshal operation_kinds: %w", err)
	}
	for _, k := range kinds {
		if !k.IsValid() {
			return merchant.MoneyDealingCapital{}, fmt.Errorf("store: money_dealing_capitals: invalid operation_kind %d", k)
		}
	}
	return merchant.MoneyDealingCapital{
		ID:                 merchant.MoneyDealingCapitalID(id),
		MoneyAdvanced:      moneyAdvanced,
		GeneralRateBP:      generalRateBP,
		OperationKinds:     kinds,
		MoneyDealingProfit: moneyDealingProfit,
		NewValueCreated:    newValueCreated,
		CreatedAt:          createdAt,
	}, nil
}

// CreateHistoricalMerchantCapital persists hm, assigning an ID and timestamp when absent.
func (m *MySQL) CreateHistoricalMerchantCapital(ctx context.Context, hm merchant.HistoricalMerchantCapital) (merchant.HistoricalMerchantCapital, error) {
	if hm.ID.IsZero() {
		hm.ID = merchant.NewHistoricalMerchantCapitalID()
	}
	if hm.CreatedAt.IsZero() {
		hm.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO historical_merchant_capitals
		(id, name, stage, profit_source, wage_labour, subordination_index, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(hm.ID),
		hm.Name,
		int64(hm.Stage),
		hm.ProfitSource,
		hm.WageLabour,
		hm.SubordinationIndex,
		hm.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return merchant.HistoricalMerchantCapital{}, ErrAlreadyExists
		}
		return merchant.HistoricalMerchantCapital{}, err
	}
	return hm, nil
}

// GetHistoricalMerchantCapital returns the historical-merchant-capital record with id, or ErrNotFound.
func (m *MySQL) GetHistoricalMerchantCapital(ctx context.Context, id merchant.HistoricalMerchantCapitalID) (merchant.HistoricalMerchantCapital, error) {
	const q = `SELECT id, name, stage, profit_source, wage_labour, subordination_index, created_at
		FROM historical_merchant_capitals WHERE id = ?`
	return scanHistoricalMerchantCapital(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListHistoricalMerchantCapitals returns all stored historical-merchant-capital records, newest first.
func (m *MySQL) ListHistoricalMerchantCapitals(ctx context.Context) ([]merchant.HistoricalMerchantCapital, error) {
	const q = `SELECT id, name, stage, profit_source, wage_labour, subordination_index, created_at
		FROM historical_merchant_capitals ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]merchant.HistoricalMerchantCapital, 0)
	for rows.Next() {
		hm, err := scanHistoricalMerchantCapital(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, hm)
	}
	return out, rows.Err()
}

func scanHistoricalMerchantCapital(s rowScanner) (merchant.HistoricalMerchantCapital, error) {
	var (
		id                 string
		name               string
		stageInt           int64
		profitSource       string
		wageLabour         bool
		subordinationIndex int64
		createdAt          time.Time
	)
	err := s.Scan(&id, &name, &stageInt, &profitSource, &wageLabour, &subordinationIndex, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return merchant.HistoricalMerchantCapital{}, ErrNotFound
		}
		return merchant.HistoricalMerchantCapital{}, err
	}
	return merchant.HistoricalMerchantCapital{
		ID:                 merchant.HistoricalMerchantCapitalID(id),
		Name:               name,
		Stage:              merchant.DevelopmentStage(stageInt),
		ProfitSource:       profitSource,
		WageLabour:         wageLabour,
		SubordinationIndex: subordinationIndex,
		CreatedAt:          createdAt,
	}, nil
}

// CreateInterestBearingCapital persists ibc, assigning an ID and timestamp when absent.
func (m *MySQL) CreateInterestBearingCapital(ctx context.Context, ibc credit.InterestBearingCapital) (credit.InterestBearingCapital, error) {
	if ibc.ID.IsZero() {
		ibc.ID = credit.NewInterestBearingCapitalID()
	}
	if ibc.CreatedAt.IsZero() {
		ibc.CreatedAt = m.now().UTC()
	}

	const q = `INSERT INTO interest_bearing_capitals
		(id, money_advanced, interest_rate_bp, loan_term_days, interest_earned, money_returned, new_value_created, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(ibc.ID),
		ibc.MoneyAdvanced,
		ibc.InterestRateBP,
		ibc.LoanTermDays,
		ibc.InterestEarned,
		ibc.MoneyReturned,
		ibc.NewValueCreated,
		ibc.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.InterestBearingCapital{}, ErrAlreadyExists
		}
		return credit.InterestBearingCapital{}, err
	}
	return ibc, nil
}

// GetInterestBearingCapital returns the interest-bearing-capital record with id, or ErrNotFound.
func (m *MySQL) GetInterestBearingCapital(ctx context.Context, id credit.InterestBearingCapitalID) (credit.InterestBearingCapital, error) {
	const q = `SELECT id, money_advanced, interest_rate_bp, loan_term_days, interest_earned, money_returned, new_value_created, created_at
		FROM interest_bearing_capitals WHERE id = ?`
	return scanInterestBearingCapital(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListInterestBearingCapitals returns all stored interest-bearing-capital records, newest first.
func (m *MySQL) ListInterestBearingCapitals(ctx context.Context) ([]credit.InterestBearingCapital, error) {
	const q = `SELECT id, money_advanced, interest_rate_bp, loan_term_days, interest_earned, money_returned, new_value_created, created_at
		FROM interest_bearing_capitals ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.InterestBearingCapital, 0)
	for rows.Next() {
		ibc, err := scanInterestBearingCapital(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, ibc)
	}
	return out, rows.Err()
}

func scanInterestBearingCapital(s rowScanner) (credit.InterestBearingCapital, error) {
	var (
		id                                             string
		moneyAdvanced, interestRateBP, loanTermDays    int64
		interestEarned, moneyReturned, newValueCreated int64
		createdAt                                      time.Time
	)
	err := s.Scan(&id, &moneyAdvanced, &interestRateBP, &loanTermDays,
		&interestEarned, &moneyReturned, &newValueCreated, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.InterestBearingCapital{}, ErrNotFound
		}
		return credit.InterestBearingCapital{}, err
	}
	return credit.InterestBearingCapital{
		ID:              credit.InterestBearingCapitalID(id),
		MoneyAdvanced:   moneyAdvanced,
		InterestRateBP:  interestRateBP,
		LoanTermDays:    loanTermDays,
		InterestEarned:  interestEarned,
		MoneyReturned:   moneyReturned,
		NewValueCreated: newValueCreated,
		CreatedAt:       createdAt,
	}, nil
}

// CreateRateOfInterest inserts r, assigning an ID and timestamp when absent.
func (m *MySQL) CreateRateOfInterest(ctx context.Context, r credit.RateOfInterest) (credit.RateOfInterest, error) {
	if r.ID.IsZero() {
		r.ID = credit.NewRateOfInterestID()
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = time.Now().UTC()
	}
	const q = `INSERT INTO rate_of_interests (id, rate_bp, average_profit_rate_bp, cycle_phase, period, created_at) VALUES (?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(r.ID), r.RateBP, r.AverageProfitRateBP, int(r.CyclePhase), r.Period, r.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return credit.RateOfInterest{}, ErrAlreadyExists
		}
		return credit.RateOfInterest{}, err
	}
	return r, nil
}

// GetRateOfInterest returns the rate-of-interest record with id, or ErrNotFound.
func (m *MySQL) GetRateOfInterest(ctx context.Context, id credit.RateOfInterestID) (credit.RateOfInterest, error) {
	const q = `SELECT id, rate_bp, average_profit_rate_bp, cycle_phase, period, created_at FROM rate_of_interests WHERE id = ?`
	return scanRateOfInterest(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListRatesOfInterest returns all stored rate-of-interest records, newest first.
func (m *MySQL) ListRatesOfInterest(ctx context.Context) ([]credit.RateOfInterest, error) {
	const q = `SELECT id, rate_bp, average_profit_rate_bp, cycle_phase, period, created_at FROM rate_of_interests ORDER BY created_at DESC, id ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.RateOfInterest, 0)
	for rows.Next() {
		r, err := scanRateOfInterest(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func scanRateOfInterest(s rowScanner) (credit.RateOfInterest, error) {
	var (
		id            string
		rateBP, avgBP int64
		phase         int
		period        string
		createdAt     time.Time
	)
	err := s.Scan(&id, &rateBP, &avgBP, &phase, &period, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.RateOfInterest{}, ErrNotFound
		}
		return credit.RateOfInterest{}, err
	}
	return credit.RateOfInterest{
		ID:                  credit.RateOfInterestID(id),
		RateBP:              rateBP,
		AverageProfitRateBP: avgBP,
		CyclePhase:          credit.IndustrialCyclePhase(phase),
		Period:              period,
		CreatedAt:           createdAt,
	}, nil
}

// CreateProfitDivision inserts pd, assigning an ID and timestamp when absent.
func (m *MySQL) CreateProfitDivision(ctx context.Context, pd credit.ProfitDivision) (credit.ProfitDivision, error) {
	if pd.ID.IsZero() {
		pd.ID = credit.NewProfitDivisionID()
	}
	if pd.CreatedAt.IsZero() {
		pd.CreatedAt = time.Now().UTC()
	}
	const q = `INSERT INTO profit_divisions (` +
		"`id`, `total_profit_bp`, `interest_bp`, `profit_of_enterprise_bp`, " +
		"`ownership_form`, `wages_labour_minutes`, `mystification_index`, `created_at`" +
		`) VALUES (?,?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(pd.ID),
		pd.TotalProfitBP,
		pd.InterestBP,
		pd.ProfitOfEnterpriseBP,
		int(pd.OwnershipForm),
		pd.WagesOfSuperintendence.LabourMinutes,
		int64(pd.MystificationIndex),
		pd.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.ProfitDivision{}, ErrAlreadyExists
		}
		return credit.ProfitDivision{}, err
	}
	return pd, nil
}

// GetProfitDivision returns the profit-division record with id, or ErrNotFound.
func (m *MySQL) GetProfitDivision(ctx context.Context, id credit.ProfitDivisionID) (credit.ProfitDivision, error) {
	const q = `SELECT ` +
		"`id`, `total_profit_bp`, `interest_bp`, `profit_of_enterprise_bp`, " +
		"`ownership_form`, `wages_labour_minutes`, `mystification_index`, `created_at` " +
		"FROM `profit_divisions` WHERE `id` = ?"
	return scanProfitDivision(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListProfitDivisions returns all stored profit-division records, newest first.
func (m *MySQL) ListProfitDivisions(ctx context.Context) ([]credit.ProfitDivision, error) {
	const q = `SELECT ` +
		"`id`, `total_profit_bp`, `interest_bp`, `profit_of_enterprise_bp`, " +
		"`ownership_form`, `wages_labour_minutes`, `mystification_index`, `created_at` " +
		"FROM `profit_divisions` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.ProfitDivision, 0)
	for rows.Next() {
		pd, err := scanProfitDivision(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, pd)
	}
	return out, rows.Err()
}

func scanProfitDivision(s rowScanner) (credit.ProfitDivision, error) {
	var (
		id                  string
		totalBP, interestBP int64
		enterpriseBP        int64
		ownershipForm       int
		labourMinutes       int64
		mystificationIndex  int64
		createdAt           time.Time
	)
	err := s.Scan(&id, &totalBP, &interestBP, &enterpriseBP,
		&ownershipForm, &labourMinutes, &mystificationIndex, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.ProfitDivision{}, ErrNotFound
		}
		return credit.ProfitDivision{}, err
	}
	return credit.ProfitDivision{
		ID:                   credit.ProfitDivisionID(id),
		TotalProfitBP:        totalBP,
		InterestBP:           interestBP,
		ProfitOfEnterpriseBP: enterpriseBP,
		OwnershipForm:        credit.CapitalOwnershipForm(ownershipForm),
		WagesOfSuperintendence: credit.WagesOfSuperintendence{
			LabourMinutes: labourMinutes,
		},
		MystificationIndex: credit.MystificationIndex(mystificationIndex),
		CreatedAt:          createdAt,
	}, nil
}

// CreateCompoundInterestSchedule inserts s, assigning an ID and timestamp when absent.
func (m *MySQL) CreateCompoundInterestSchedule(ctx context.Context, s credit.CompoundInterestSchedule) (credit.CompoundInterestSchedule, error) {
	if s.ID.IsZero() {
		s.ID = credit.NewCompoundInterestID()
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO compound_interest_schedules (` +
		"`id`, `principal`, `rate_bp`, `period_years`, `final_value`, `new_value_created`, `fetish_form`, `created_at`" +
		`) VALUES (?,?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(s.ID),
		s.Principal,
		s.RateBP,
		s.PeriodYears,
		s.FinalValue,
		s.NewValueCreated,
		int(s.FetishForm),
		s.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.CompoundInterestSchedule{}, ErrAlreadyExists
		}
		return credit.CompoundInterestSchedule{}, err
	}
	return s, nil
}

// GetCompoundInterestSchedule returns the schedule with id, or ErrNotFound.
func (m *MySQL) GetCompoundInterestSchedule(ctx context.Context, id credit.CompoundInterestID) (credit.CompoundInterestSchedule, error) {
	const q = `SELECT ` +
		"`id`, `principal`, `rate_bp`, `period_years`, `final_value`, `new_value_created`, `fetish_form`, `created_at` " +
		"FROM `compound_interest_schedules` WHERE `id` = ?"
	return scanCompoundInterestSchedule(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListCompoundInterestSchedules returns all stored schedules, newest first.
func (m *MySQL) ListCompoundInterestSchedules(ctx context.Context) ([]credit.CompoundInterestSchedule, error) {
	const q = `SELECT ` +
		"`id`, `principal`, `rate_bp`, `period_years`, `final_value`, `new_value_created`, `fetish_form`, `created_at` " +
		"FROM `compound_interest_schedules` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.CompoundInterestSchedule, 0)
	for rows.Next() {
		s, err := scanCompoundInterestSchedule(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func scanCompoundInterestSchedule(s rowScanner) (credit.CompoundInterestSchedule, error) {
	var (
		id              string
		principal       int64
		rateBP          int64
		periodYears     int64
		finalValue      int64
		newValueCreated int64
		fetishForm      int
		createdAt       time.Time
	)
	err := s.Scan(&id, &principal, &rateBP, &periodYears, &finalValue, &newValueCreated, &fetishForm, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.CompoundInterestSchedule{}, ErrNotFound
		}
		return credit.CompoundInterestSchedule{}, err
	}
	return credit.CompoundInterestSchedule{
		ID:              credit.CompoundInterestID(id),
		Principal:       principal,
		RateBP:          rateBP,
		PeriodYears:     periodYears,
		FinalValue:      finalValue,
		NewValueCreated: newValueCreated,
		FetishForm:      credit.FetishForm(fetishForm),
		CreatedAt:       createdAt,
	}, nil
}

// CreateBillOfExchange inserts b, assigning an ID and timestamp when absent (Vol. III Ch. 25).
func (m *MySQL) CreateBillOfExchange(ctx context.Context, b credit.BillOfExchange) (credit.BillOfExchange, error) {
	if b.ID.IsZero() {
		b.ID = credit.NewBillOfExchangeID()
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO bills_of_exchange (` +
		"`id`, `drawer`, `drawee`, `face_value`, `maturity_days`, `description`, `created_at`" +
		`) VALUES (?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(b.ID),
		b.Drawer,
		b.Drawee,
		b.FaceValue,
		b.MaturityDays,
		b.Description,
		b.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.BillOfExchange{}, ErrAlreadyExists
		}
		return credit.BillOfExchange{}, err
	}
	return b, nil
}

// GetBillOfExchange returns the bill with id, or ErrNotFound.
func (m *MySQL) GetBillOfExchange(ctx context.Context, id credit.BillOfExchangeID) (credit.BillOfExchange, error) {
	const q = `SELECT ` +
		"`id`, `drawer`, `drawee`, `face_value`, `maturity_days`, `description`, `created_at` " +
		"FROM `bills_of_exchange` WHERE `id` = ?"
	return scanBillOfExchange(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListBillsOfExchange returns all stored bills, newest first.
func (m *MySQL) ListBillsOfExchange(ctx context.Context) ([]credit.BillOfExchange, error) {
	const q = `SELECT ` +
		"`id`, `drawer`, `drawee`, `face_value`, `maturity_days`, `description`, `created_at` " +
		"FROM `bills_of_exchange` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.BillOfExchange, 0)
	for rows.Next() {
		b, err := scanBillOfExchange(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func scanBillOfExchange(s rowScanner) (credit.BillOfExchange, error) {
	var (
		id           string
		drawer       string
		drawee       string
		faceValue    int64
		maturityDays int64
		description  string
		createdAt    time.Time
	)
	err := s.Scan(&id, &drawer, &drawee, &faceValue, &maturityDays, &description, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.BillOfExchange{}, ErrNotFound
		}
		return credit.BillOfExchange{}, err
	}
	return credit.BillOfExchange{
		ID:           credit.BillOfExchangeID(id),
		Drawer:       drawer,
		Drawee:       drawee,
		FaceValue:    faceValue,
		MaturityDays: maturityDays,
		Description:  description,
		CreatedAt:    createdAt,
	}, nil
}

// CreateFictitiousCapital inserts fc, assigning an ID and timestamp when absent (Vol. III Ch. 25).
func (m *MySQL) CreateFictitiousCapital(ctx context.Context, fc credit.FictitiousCapital) (credit.FictitiousCapital, error) {
	if fc.ID.IsZero() {
		fc.ID = credit.NewFictitiousCapitalID()
	}
	if fc.CreatedAt.IsZero() {
		fc.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO fictitious_capitals (` +
		"`id`, `kind`, `name`, `annual_income`, `rate_bp`, `capitalised_value`, `real_capital_exists`, `description`, `created_at`" +
		`) VALUES (?,?,?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(fc.ID),
		int(fc.Kind),
		fc.Name,
		fc.AnnualIncome,
		fc.RateBP,
		fc.CapitalisedValue,
		fc.RealCapitalExists,
		fc.Description,
		fc.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.FictitiousCapital{}, ErrAlreadyExists
		}
		return credit.FictitiousCapital{}, err
	}
	return fc, nil
}

// GetFictitiousCapital returns the record with id, or ErrNotFound.
func (m *MySQL) GetFictitiousCapital(ctx context.Context, id credit.FictitiousCapitalID) (credit.FictitiousCapital, error) {
	const q = `SELECT ` +
		"`id`, `kind`, `name`, `annual_income`, `rate_bp`, `capitalised_value`, `real_capital_exists`, `description`, `created_at` " +
		"FROM `fictitious_capitals` WHERE `id` = ?"
	return scanFictitiousCapital(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListFictitiousCapitals returns all stored records, newest first.
func (m *MySQL) ListFictitiousCapitals(ctx context.Context) ([]credit.FictitiousCapital, error) {
	const q = `SELECT ` +
		"`id`, `kind`, `name`, `annual_income`, `rate_bp`, `capitalised_value`, `real_capital_exists`, `description`, `created_at` " +
		"FROM `fictitious_capitals` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.FictitiousCapital, 0)
	for rows.Next() {
		fc, err := scanFictitiousCapital(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, fc)
	}
	return out, rows.Err()
}

func scanFictitiousCapital(s rowScanner) (credit.FictitiousCapital, error) {
	var (
		id                string
		kind              int
		name              string
		annualIncome      int64
		rateBP            int64
		capitalisedValue  int64
		realCapitalExists bool
		description       string
		createdAt         time.Time
	)
	err := s.Scan(&id, &kind, &name, &annualIncome, &rateBP, &capitalisedValue, &realCapitalExists, &description, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.FictitiousCapital{}, ErrNotFound
		}
		return credit.FictitiousCapital{}, err
	}
	return credit.FictitiousCapital{
		ID:                credit.FictitiousCapitalID(id),
		Kind:              credit.FictitiousKind(kind),
		Name:              name,
		AnnualIncome:      annualIncome,
		RateBP:            rateBP,
		CapitalisedValue:  capitalisedValue,
		RealCapitalExists: realCapitalExists,
		Description:       description,
		CreatedAt:         createdAt,
	}, nil
}

// CreateMoneyCapitalAccumulation inserts a, assigning an ID and timestamp when absent (Vol. III Ch. 26).
func (m *MySQL) CreateMoneyCapitalAccumulation(ctx context.Context, a credit.MoneyCapitalAccumulation) (credit.MoneyCapitalAccumulation, error) {
	if a.ID.IsZero() {
		a.ID = credit.NewMoneyCapitalAccumulationID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO money_capital_accumulations (` +
		"`id`, `name`, `phase`, `interest_rate_bp`, `period`, `loanable_amount`, `loanable_abundant`, `idle_amount`, `idle_description`, `gap_bp`, `gap_description`, `description`, `created_at`" +
		`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID),
		a.Name,
		int(a.CycleObservation.Phase),
		a.CycleObservation.InterestRateBP,
		a.CycleObservation.Period,
		a.LoanableSupply.Amount,
		a.LoanableSupply.Abundant,
		a.IdleCapital.Amount,
		a.IdleCapital.Description,
		a.Gap.GapBP,
		a.Gap.Description,
		a.Description,
		a.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.MoneyCapitalAccumulation{}, ErrAlreadyExists
		}
		return credit.MoneyCapitalAccumulation{}, err
	}
	return a, nil
}

// GetMoneyCapitalAccumulation returns the record with id, or ErrNotFound.
func (m *MySQL) GetMoneyCapitalAccumulation(ctx context.Context, id credit.MoneyCapitalAccumulationID) (credit.MoneyCapitalAccumulation, error) {
	const q = `SELECT ` +
		"`id`, `name`, `phase`, `interest_rate_bp`, `period`, `loanable_amount`, `loanable_abundant`, `idle_amount`, `idle_description`, `gap_bp`, `gap_description`, `description`, `created_at` " +
		"FROM `money_capital_accumulations` WHERE `id` = ?"
	return scanMoneyCapitalAccumulation(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListMoneyCapitalAccumulations returns all stored records, newest first.
func (m *MySQL) ListMoneyCapitalAccumulations(ctx context.Context) ([]credit.MoneyCapitalAccumulation, error) {
	const q = `SELECT ` +
		"`id`, `name`, `phase`, `interest_rate_bp`, `period`, `loanable_amount`, `loanable_abundant`, `idle_amount`, `idle_description`, `gap_bp`, `gap_description`, `description`, `created_at` " +
		"FROM `money_capital_accumulations` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.MoneyCapitalAccumulation, 0)
	for rows.Next() {
		a, err := scanMoneyCapitalAccumulation(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func scanMoneyCapitalAccumulation(s rowScanner) (credit.MoneyCapitalAccumulation, error) {
	var (
		id               string
		name             string
		phase            int
		interestRateBP   int64
		period           string
		loanableAmount   int64
		loanableAbundant bool
		idleAmount       int64
		idleDescription  string
		gapBP            int64
		gapDescription   string
		description      string
		createdAt        time.Time
	)
	err := s.Scan(
		&id, &name, &phase, &interestRateBP, &period,
		&loanableAmount, &loanableAbundant,
		&idleAmount, &idleDescription,
		&gapBP, &gapDescription,
		&description, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.MoneyCapitalAccumulation{}, ErrNotFound
		}
		return credit.MoneyCapitalAccumulation{}, err
	}
	return credit.MoneyCapitalAccumulation{
		ID:   credit.MoneyCapitalAccumulationID(id),
		Name: name,
		CycleObservation: credit.InterestCycleObservation{
			Phase:          credit.IndustrialCyclePhase(phase),
			InterestRateBP: interestRateBP,
			Period:         period,
		},
		LoanableSupply: credit.LoanableCapitalSupply{
			Amount:   loanableAmount,
			Abundant: loanableAbundant,
		},
		IdleCapital: credit.IdleIndustrialCapital{
			Amount:      idleAmount,
			Description: idleDescription,
		},
		Gap: credit.AccumulationGap{
			GapBP:       gapBP,
			Description: gapDescription,
		},
		Description: description,
		CreatedAt:   createdAt,
	}, nil
}

// ── Vol. III Ch. 27 — The Role of Credit (StockCompany, CooperativeFactory) ──

// CreateStockCompany inserts sc, assigning an ID and timestamp when absent (Vol. III Ch. 27).
func (m *MySQL) CreateStockCompany(ctx context.Context, sc credit.StockCompany) (credit.StockCompany, error) {
	if sc.ID.IsZero() {
		sc.ID = credit.NewStockCompanyID()
	}
	if sc.CreatedAt.IsZero() {
		sc.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO stock_companies (` +
		"`id`, `name`, `fixed_capital`, `floating_capital`, `total_capital`, " +
		"`manager_name`, `manager_title`, `manager_wage_minutes`, `description`, `created_at`" +
		`) VALUES (?,?,?,?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(sc.ID),
		sc.Name,
		sc.FixedCapital,
		sc.FloatingCapital,
		sc.TotalCapital,
		sc.Manager.Name,
		sc.Manager.Title,
		sc.Manager.WageMinutes,
		sc.Description,
		sc.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.StockCompany{}, ErrAlreadyExists
		}
		return credit.StockCompany{}, err
	}
	return sc, nil
}

// GetStockCompany returns the stock company with id, or ErrNotFound.
func (m *MySQL) GetStockCompany(ctx context.Context, id credit.StockCompanyID) (credit.StockCompany, error) {
	const q = `SELECT ` +
		"`id`, `name`, `fixed_capital`, `floating_capital`, `total_capital`, " +
		"`manager_name`, `manager_title`, `manager_wage_minutes`, `description`, `created_at` " +
		"FROM `stock_companies` WHERE `id` = ?"
	return scanStockCompany(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListStockCompanies returns all stored stock companies, newest first.
func (m *MySQL) ListStockCompanies(ctx context.Context) ([]credit.StockCompany, error) {
	const q = `SELECT ` +
		"`id`, `name`, `fixed_capital`, `floating_capital`, `total_capital`, " +
		"`manager_name`, `manager_title`, `manager_wage_minutes`, `description`, `created_at` " +
		"FROM `stock_companies` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.StockCompany, 0)
	for rows.Next() {
		sc, err := scanStockCompany(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, sc)
	}
	return out, rows.Err()
}

func scanStockCompany(s rowScanner) (credit.StockCompany, error) {
	var (
		id              string
		name            string
		fixedCapital    int64
		floatingCapital int64
		totalCapital    int64
		managerName     string
		managerTitle    string
		managerWage     int64
		description     string
		createdAt       time.Time
	)
	err := s.Scan(
		&id, &name, &fixedCapital, &floatingCapital, &totalCapital,
		&managerName, &managerTitle, &managerWage,
		&description, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.StockCompany{}, ErrNotFound
		}
		return credit.StockCompany{}, err
	}
	return credit.StockCompany{
		ID:              credit.StockCompanyID(id),
		Name:            name,
		FixedCapital:    fixedCapital,
		FloatingCapital: floatingCapital,
		TotalCapital:    totalCapital,
		Manager: credit.Manager{
			Name:        managerName,
			Title:       managerTitle,
			WageMinutes: managerWage,
		},
		Description: description,
		CreatedAt:   createdAt,
	}, nil
}

// CreateCooperativeFactory inserts cf, assigning an ID and timestamp when absent (Vol. III Ch. 27).
func (m *MySQL) CreateCooperativeFactory(ctx context.Context, cf credit.CooperativeFactory) (credit.CooperativeFactory, error) {
	if cf.ID.IsZero() {
		cf.ID = credit.NewCooperativeFactoryID()
	}
	if cf.CreatedAt.IsZero() {
		cf.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO cooperative_factories (` +
		"`id`, `name`, `worker_owned`, `worker_count`, `capital_advanced`, `description`, `created_at`" +
		`) VALUES (?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(cf.ID),
		cf.Name,
		cf.WorkerOwned,
		cf.WorkerCount,
		cf.CapitalAdvanced,
		cf.Description,
		cf.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.CooperativeFactory{}, ErrAlreadyExists
		}
		return credit.CooperativeFactory{}, err
	}
	return cf, nil
}

// GetCooperativeFactory returns the cooperative factory with id, or ErrNotFound.
func (m *MySQL) GetCooperativeFactory(ctx context.Context, id credit.CooperativeFactoryID) (credit.CooperativeFactory, error) {
	const q = `SELECT ` +
		"`id`, `name`, `worker_owned`, `worker_count`, `capital_advanced`, `description`, `created_at` " +
		"FROM `cooperative_factories` WHERE `id` = ?"
	return scanCooperativeFactory(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListCooperativeFactories returns all stored cooperative factories, newest first.
func (m *MySQL) ListCooperativeFactories(ctx context.Context) ([]credit.CooperativeFactory, error) {
	const q = `SELECT ` +
		"`id`, `name`, `worker_owned`, `worker_count`, `capital_advanced`, `description`, `created_at` " +
		"FROM `cooperative_factories` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.CooperativeFactory, 0)
	for rows.Next() {
		cf, err := scanCooperativeFactory(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, cf)
	}
	return out, rows.Err()
}

func scanCooperativeFactory(s rowScanner) (credit.CooperativeFactory, error) {
	var (
		id              string
		name            string
		workerOwned     bool
		workerCount     int64
		capitalAdvanced int64
		description     string
		createdAt       time.Time
	)
	err := s.Scan(
		&id, &name, &workerOwned, &workerCount, &capitalAdvanced,
		&description, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.CooperativeFactory{}, ErrNotFound
		}
		return credit.CooperativeFactory{}, err
	}
	return credit.CooperativeFactory{
		ID:              credit.CooperativeFactoryID(id),
		Name:            name,
		WorkerOwned:     workerOwned,
		WorkerCount:     workerCount,
		CapitalAdvanced: capitalAdvanced,
		Description:     description,
		CreatedAt:       createdAt,
	}, nil
}

// ── Vol. III Ch. 28 — Medium of Circulation and Capital (CurrencyObservation) ──

// CreateCurrencyObservation inserts o, assigning an ID and timestamp when absent (Vol. III Ch. 28).
func (m *MySQL) CreateCurrencyObservation(ctx context.Context, o credit.CurrencyObservation) (credit.CurrencyObservation, error) {
	if o.ID.IsZero() {
		o.ID = credit.NewCurrencyObservationID()
	}
	if o.CreatedAt.IsZero() {
		o.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO currency_observations (` +
		"`id`, `name`, `total_currency_outstanding`, `coin_function_amount`, `capital_transfer_amount`, " +
		"`reserve_amount`, `reserve_description`, `description`, `created_at`" +
		`) VALUES (?,?,?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(o.ID),
		o.Name,
		o.TotalCurrencyOutstanding,
		o.CoinFunctionAmount,
		o.CapitalTransferAmount,
		o.ReserveFund.Amount,
		o.ReserveFund.Description,
		o.Description,
		o.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.CurrencyObservation{}, ErrAlreadyExists
		}
		return credit.CurrencyObservation{}, err
	}
	return o, nil
}

// GetCurrencyObservation returns the record with id, or ErrNotFound.
func (m *MySQL) GetCurrencyObservation(ctx context.Context, id credit.CurrencyObservationID) (credit.CurrencyObservation, error) {
	const q = `SELECT ` +
		"`id`, `name`, `total_currency_outstanding`, `coin_function_amount`, `capital_transfer_amount`, " +
		"`reserve_amount`, `reserve_description`, `description`, `created_at` " +
		"FROM `currency_observations` WHERE `id` = ?"
	return scanCurrencyObservation(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListCurrencyObservations returns all stored records, newest first.
func (m *MySQL) ListCurrencyObservations(ctx context.Context) ([]credit.CurrencyObservation, error) {
	const q = `SELECT ` +
		"`id`, `name`, `total_currency_outstanding`, `coin_function_amount`, `capital_transfer_amount`, " +
		"`reserve_amount`, `reserve_description`, `description`, `created_at` " +
		"FROM `currency_observations` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.CurrencyObservation, 0)
	for rows.Next() {
		o, err := scanCurrencyObservation(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

func scanCurrencyObservation(s rowScanner) (credit.CurrencyObservation, error) {
	var (
		id                       string
		name                     string
		totalCurrencyOutstanding int64
		coinFunctionAmount       int64
		capitalTransferAmount    int64
		reserveAmount            int64
		reserveDescription       string
		description              string
		createdAt                time.Time
	)
	err := s.Scan(
		&id, &name,
		&totalCurrencyOutstanding, &coinFunctionAmount, &capitalTransferAmount,
		&reserveAmount, &reserveDescription,
		&description, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.CurrencyObservation{}, ErrNotFound
		}
		return credit.CurrencyObservation{}, err
	}
	return credit.CurrencyObservation{
		ID:                       credit.CurrencyObservationID(id),
		Name:                     name,
		TotalCurrencyOutstanding: totalCurrencyOutstanding,
		CoinFunctionAmount:       coinFunctionAmount,
		CapitalTransferAmount:    capitalTransferAmount,
		ReserveFund: credit.ReserveFund{
			Amount:      reserveAmount,
			Description: reserveDescription,
		},
		Description: description,
		CreatedAt:   createdAt,
	}, nil
}

// ── Vol. III Ch. 29 — Component Parts of Bank Capital ─────────────────────────

// CreateBankCapital persists bc, assigning an ID and timestamp when absent (Vol. III Ch. 29).
func (m *MySQL) CreateBankCapital(ctx context.Context, bc credit.BankCapital) (credit.BankCapital, error) {
	if bc.ID.IsZero() {
		bc.ID = credit.NewBankCapitalID()
	}
	if bc.CreatedAt.IsZero() {
		bc.CreatedAt = m.now().UTC()
	}
	components := bc.Components
	if components == nil {
		components = []credit.BankCapitalComponent{}
	}
	componentsJSON, err := json.Marshal(components)
	if err != nil {
		return credit.BankCapital{}, err
	}
	const q = `INSERT INTO bank_capitals (` +
		"`id`, `name`, `cash_amount`, `securities_amount`, `total_capital`, `components_json`, `description`, `created_at`" +
		`) VALUES (?,?,?,?,?,?,?,?)`
	_, err = m.db.ExecContext(ctx, q,
		string(bc.ID),
		bc.Name,
		bc.CashAmount,
		bc.SecuritiesAmount,
		bc.TotalCapital,
		string(componentsJSON),
		bc.Description,
		bc.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.BankCapital{}, ErrAlreadyExists
		}
		return credit.BankCapital{}, err
	}
	return bc, nil
}

// GetBankCapital returns the bank-capital record with id, or ErrNotFound.
func (m *MySQL) GetBankCapital(ctx context.Context, id credit.BankCapitalID) (credit.BankCapital, error) {
	const q = `SELECT ` +
		"`id`, `name`, `cash_amount`, `securities_amount`, `total_capital`, `components_json`, `description`, `created_at` " +
		"FROM `bank_capitals` WHERE `id` = ?"
	return scanBankCapital(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListBankCapitals returns all stored bank-capital records, newest first.
func (m *MySQL) ListBankCapitals(ctx context.Context) ([]credit.BankCapital, error) {
	const q = `SELECT ` +
		"`id`, `name`, `cash_amount`, `securities_amount`, `total_capital`, `components_json`, `description`, `created_at` " +
		"FROM `bank_capitals` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.BankCapital, 0)
	for rows.Next() {
		bc, err := scanBankCapital(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, bc)
	}
	return out, rows.Err()
}

func scanBankCapital(s rowScanner) (credit.BankCapital, error) {
	var (
		id               string
		name             string
		cashAmount       int64
		securitiesAmount int64
		totalCapital     int64
		componentsStr    string
		description      string
		createdAt        time.Time
	)
	err := s.Scan(&id, &name, &cashAmount, &securitiesAmount, &totalCapital, &componentsStr, &description, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.BankCapital{}, ErrNotFound
		}
		return credit.BankCapital{}, err
	}
	var components []credit.BankCapitalComponent
	if err := json.Unmarshal([]byte(componentsStr), &components); err != nil {
		return credit.BankCapital{}, fmt.Errorf("store: bank_capitals: unmarshal components_json: %w", err)
	}
	if components == nil {
		components = []credit.BankCapitalComponent{}
	}
	for _, c := range components {
		if !c.Kind.IsValid() {
			return credit.BankCapital{}, fmt.Errorf("store: bank_capitals: invalid component kind %d", c.Kind)
		}
	}
	return credit.BankCapital{
		ID:               credit.BankCapitalID(id),
		Name:             name,
		CashAmount:       cashAmount,
		SecuritiesAmount: securitiesAmount,
		TotalCapital:     totalCapital,
		Components:       components,
		Description:      description,
		CreatedAt:        createdAt,
	}, nil
}

// CreateFictitiousCapitalValuation persists v, assigning an ID and timestamp when absent (Vol. III Ch. 29).
func (m *MySQL) CreateFictitiousCapitalValuation(ctx context.Context, v credit.FictitiousCapitalValuation) (credit.FictitiousCapitalValuation, error) {
	if v.ID.IsZero() {
		v.ID = credit.NewFictitiousCapitalValuationID()
	}
	if v.CreatedAt.IsZero() {
		v.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO fictitious_capital_valuations (` +
		"`id`, `name`, `nominal_value`, `annual_income`, `interest_rate_bp`, `dividend_rate_bp`, `market_value`, `description`, `created_at`" +
		`) VALUES (?,?,?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(v.ID),
		v.Name,
		v.NominalValue,
		v.AnnualIncome,
		v.InterestRateBP,
		v.DividendRateBP,
		v.MarketValue,
		v.Description,
		v.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.FictitiousCapitalValuation{}, ErrAlreadyExists
		}
		return credit.FictitiousCapitalValuation{}, err
	}
	return v, nil
}

// GetFictitiousCapitalValuation returns the valuation with id, or ErrNotFound.
func (m *MySQL) GetFictitiousCapitalValuation(ctx context.Context, id credit.FictitiousCapitalValuationID) (credit.FictitiousCapitalValuation, error) {
	const q = `SELECT ` +
		"`id`, `name`, `nominal_value`, `annual_income`, `interest_rate_bp`, `dividend_rate_bp`, `market_value`, `description`, `created_at` " +
		"FROM `fictitious_capital_valuations` WHERE `id` = ?"
	return scanFictitiousCapitalValuation(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListFictitiousCapitalValuations returns all stored valuations, newest first.
func (m *MySQL) ListFictitiousCapitalValuations(ctx context.Context) ([]credit.FictitiousCapitalValuation, error) {
	const q = `SELECT ` +
		"`id`, `name`, `nominal_value`, `annual_income`, `interest_rate_bp`, `dividend_rate_bp`, `market_value`, `description`, `created_at` " +
		"FROM `fictitious_capital_valuations` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.FictitiousCapitalValuation, 0)
	for rows.Next() {
		v, err := scanFictitiousCapitalValuation(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func scanFictitiousCapitalValuation(s rowScanner) (credit.FictitiousCapitalValuation, error) {
	var (
		id             string
		name           string
		nominalValue   int64
		annualIncome   int64
		interestRateBP int64
		dividendRateBP int64
		marketValue    int64
		description    string
		createdAt      time.Time
	)
	err := s.Scan(&id, &name, &nominalValue, &annualIncome, &interestRateBP, &dividendRateBP, &marketValue, &description, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.FictitiousCapitalValuation{}, ErrNotFound
		}
		return credit.FictitiousCapitalValuation{}, err
	}
	return credit.FictitiousCapitalValuation{
		ID:             credit.FictitiousCapitalValuationID(id),
		Name:           name,
		NominalValue:   nominalValue,
		AnnualIncome:   annualIncome,
		InterestRateBP: interestRateBP,
		DividendRateBP: dividendRateBP,
		MarketValue:    marketValue,
		Description:    description,
		CreatedAt:      createdAt,
	}, nil
}

// ── Vol. III Ch. 30 — Money-Capital and Real Capital, I (RealCapitalAccumulation) ──

// CreateRealCapitalAccumulation inserts a, assigning an ID and timestamp when absent (Vol. III Ch. 30).
func (m *MySQL) CreateRealCapitalAccumulation(ctx context.Context, a credit.RealCapitalAccumulation) (credit.RealCapitalAccumulation, error) {
	if a.ID.IsZero() {
		a.ID = credit.NewRealCapitalAccumulationID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO real_capital_accumulations (` +
		"`id`, `name`, `phase`, `money_capital_growth_bp`, `real_capital_growth_bp`, `corresponds`, `correspondence_explanation`, `reserve_capital`, `credit_limit_description`, `commercial_credit_contracts`, `cause_is_money_scarcity`, `description`, `created_at`" +
		`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(a.ID),
		a.Name,
		int(a.Phase),
		a.Correspondence.MoneyCapitalGrowthBP,
		a.Correspondence.RealCapitalGrowthBP,
		a.Correspondence.Corresponds,
		a.Correspondence.Explanation,
		a.CreditLimit.ReserveCapital,
		a.CreditLimit.Description,
		a.CommercialCreditContracts,
		a.CauseIsMoneyScarcity,
		a.Description,
		a.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.RealCapitalAccumulation{}, ErrAlreadyExists
		}
		return credit.RealCapitalAccumulation{}, err
	}
	return a, nil
}

// GetRealCapitalAccumulation returns the record with id, or ErrNotFound.
func (m *MySQL) GetRealCapitalAccumulation(ctx context.Context, id credit.RealCapitalAccumulationID) (credit.RealCapitalAccumulation, error) {
	const q = `SELECT ` +
		"`id`, `name`, `phase`, `money_capital_growth_bp`, `real_capital_growth_bp`, `corresponds`, `correspondence_explanation`, `reserve_capital`, `credit_limit_description`, `commercial_credit_contracts`, `cause_is_money_scarcity`, `description`, `created_at` " +
		"FROM `real_capital_accumulations` WHERE `id` = ?"
	return scanRealCapitalAccumulation(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListRealCapitalAccumulations returns all stored records, newest first.
func (m *MySQL) ListRealCapitalAccumulations(ctx context.Context) ([]credit.RealCapitalAccumulation, error) {
	const q = `SELECT ` +
		"`id`, `name`, `phase`, `money_capital_growth_bp`, `real_capital_growth_bp`, `corresponds`, `correspondence_explanation`, `reserve_capital`, `credit_limit_description`, `commercial_credit_contracts`, `cause_is_money_scarcity`, `description`, `created_at` " +
		"FROM `real_capital_accumulations` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.RealCapitalAccumulation, 0)
	for rows.Next() {
		a, err := scanRealCapitalAccumulation(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func scanRealCapitalAccumulation(s rowScanner) (credit.RealCapitalAccumulation, error) {
	var (
		id                        string
		name                      string
		phase                     int
		moneyCapitalGrowthBP      int64
		realCapitalGrowthBP       int64
		corresponds               bool
		correspondenceExplanation string
		reserveCapital            int64
		creditLimitDescription    string
		commercialCreditContracts bool
		causeIsMoneyScarcity      bool
		description               string
		createdAt                 time.Time
	)
	err := s.Scan(
		&id, &name, &phase,
		&moneyCapitalGrowthBP, &realCapitalGrowthBP, &corresponds, &correspondenceExplanation,
		&reserveCapital, &creditLimitDescription,
		&commercialCreditContracts, &causeIsMoneyScarcity,
		&description, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.RealCapitalAccumulation{}, ErrNotFound
		}
		return credit.RealCapitalAccumulation{}, err
	}
	return credit.RealCapitalAccumulation{
		ID:    credit.RealCapitalAccumulationID(id),
		Name:  name,
		Phase: credit.IndustrialCyclePhase(phase),
		Correspondence: credit.AccumulationCorrespondence{
			MoneyCapitalGrowthBP: moneyCapitalGrowthBP,
			RealCapitalGrowthBP:  realCapitalGrowthBP,
			Corresponds:          corresponds,
			Explanation:          correspondenceExplanation,
		},
		CreditLimit: credit.CreditLimit{
			ReserveCapital: reserveCapital,
			Description:    creditLimitDescription,
		},
		CommercialCreditContracts: commercialCreditContracts,
		CauseIsMoneyScarcity:      causeIsMoneyScarcity,
		Description:               description,
		CreatedAt:                 createdAt,
	}, nil
}

// ── Vol. III Ch. 31 — Money-Capital and Real Capital, II (FloatingCapital) ──

// CreateFloatingCapital inserts f, assigning an ID and timestamp when absent (Vol. III Ch. 31).
func (m *MySQL) CreateFloatingCapital(ctx context.Context, f credit.FloatingCapital) (credit.FloatingCapital, error) {
	if f.ID.IsZero() {
		f.ID = credit.NewFloatingCapitalID()
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO floating_capitals (` +
		"`id`, `name`, `amount`, `source`, `velocity_money_amount`, `velocity_times_lent`, `velocity_loan_capital_created`, `reserve_saving_amount`, `reserve_saving_description`, `description`, `created_at`" +
		`) VALUES (?,?,?,?,?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(f.ID),
		f.Name,
		f.Amount,
		int(f.Source),
		f.Velocity.MoneyAmount,
		f.Velocity.TimesLent,
		f.Velocity.LoanCapitalCreated,
		f.ReserveSaving.Amount,
		f.ReserveSaving.Description,
		f.Description,
		f.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.FloatingCapital{}, ErrAlreadyExists
		}
		return credit.FloatingCapital{}, err
	}
	return f, nil
}

// GetFloatingCapital returns the record with id, or ErrNotFound.
func (m *MySQL) GetFloatingCapital(ctx context.Context, id credit.FloatingCapitalID) (credit.FloatingCapital, error) {
	const q = `SELECT ` +
		"`id`, `name`, `amount`, `source`, `velocity_money_amount`, `velocity_times_lent`, `velocity_loan_capital_created`, `reserve_saving_amount`, `reserve_saving_description`, `description`, `created_at` " +
		"FROM `floating_capitals` WHERE `id` = ?"
	return scanFloatingCapital(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListFloatingCapitals returns all stored records, newest first.
func (m *MySQL) ListFloatingCapitals(ctx context.Context) ([]credit.FloatingCapital, error) {
	const q = `SELECT ` +
		"`id`, `name`, `amount`, `source`, `velocity_money_amount`, `velocity_times_lent`, `velocity_loan_capital_created`, `reserve_saving_amount`, `reserve_saving_description`, `description`, `created_at` " +
		"FROM `floating_capitals` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.FloatingCapital, 0)
	for rows.Next() {
		f, err := scanFloatingCapital(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func scanFloatingCapital(s rowScanner) (credit.FloatingCapital, error) {
	var (
		id                         string
		name                       string
		amount                     int64
		source                     int
		velocityMoneyAmount        int64
		velocityTimesLent          int64
		velocityLoanCapitalCreated int64
		reserveSavingAmount        int64
		reserveSavingDescription   string
		description                string
		createdAt                  time.Time
	)
	err := s.Scan(
		&id, &name, &amount, &source,
		&velocityMoneyAmount, &velocityTimesLent, &velocityLoanCapitalCreated,
		&reserveSavingAmount, &reserveSavingDescription,
		&description, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.FloatingCapital{}, ErrNotFound
		}
		return credit.FloatingCapital{}, err
	}
	return credit.FloatingCapital{
		ID:     credit.FloatingCapitalID(id),
		Name:   name,
		Amount: amount,
		Source: credit.FloatingCapitalSource(source),
		Velocity: credit.LoanCapitalVelocity{
			MoneyAmount:        velocityMoneyAmount,
			TimesLent:          velocityTimesLent,
			LoanCapitalCreated: velocityLoanCapitalCreated,
		},
		ReserveSaving: credit.CirculationReserveSaving{
			Amount:      reserveSavingAmount,
			Description: reserveSavingDescription,
		},
		Description: description,
		CreatedAt:   createdAt,
	}, nil
}

// ── Vol. III Ch. 32 — Money-Capital and Real Capital, III (CapitalRelease) ──

// CreateCapitalRelease inserts c, assigning an ID and timestamp when absent (Vol. III Ch. 32).
func (m *MySQL) CreateCapitalRelease(ctx context.Context, c credit.CapitalRelease) (credit.CapitalRelease, error) {
	if c.ID.IsZero() {
		c.ID = credit.NewCapitalReleaseID()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO capital_releases (` +
		"`id`, `name`, `released_amount`, `cause`, `reinvested`, `overstatement_pct`, `overstatement_description`, `description`, `created_at`" +
		`) VALUES (?,?,?,?,?,?,?,?,?)`
	_, err := m.db.ExecContext(ctx, q,
		string(c.ID),
		c.Name,
		c.ReleasedAmount,
		int(c.Cause),
		c.Reinvested,
		c.Overstatement.OverstatementPct,
		c.Overstatement.Description,
		c.Description,
		c.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.CapitalRelease{}, ErrAlreadyExists
		}
		return credit.CapitalRelease{}, err
	}
	return c, nil
}

// GetCapitalRelease returns the record with id, or ErrNotFound.
func (m *MySQL) GetCapitalRelease(ctx context.Context, id credit.CapitalReleaseID) (credit.CapitalRelease, error) {
	const q = `SELECT ` +
		"`id`, `name`, `released_amount`, `cause`, `reinvested`, `overstatement_pct`, `overstatement_description`, `description`, `created_at` " +
		"FROM `capital_releases` WHERE `id` = ?"
	return scanCapitalRelease(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListCapitalReleases returns all stored records, newest first.
func (m *MySQL) ListCapitalReleases(ctx context.Context) ([]credit.CapitalRelease, error) {
	const q = `SELECT ` +
		"`id`, `name`, `released_amount`, `cause`, `reinvested`, `overstatement_pct`, `overstatement_description`, `description`, `created_at` " +
		"FROM `capital_releases` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.CapitalRelease, 0)
	for rows.Next() {
		c, err := scanCapitalRelease(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

func scanCapitalRelease(s rowScanner) (credit.CapitalRelease, error) {
	var (
		id                       string
		name                     string
		releasedAmount           int64
		cause                    int
		reinvested               bool
		overstatementPct         int64
		overstatementDescription string
		description              string
		createdAt                time.Time
	)
	err := s.Scan(
		&id, &name, &releasedAmount, &cause, &reinvested,
		&overstatementPct, &overstatementDescription,
		&description, &createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.CapitalRelease{}, ErrNotFound
		}
		return credit.CapitalRelease{}, err
	}
	return credit.CapitalRelease{
		ID:             credit.CapitalReleaseID(id),
		Name:           name,
		ReleasedAmount: releasedAmount,
		Cause:          credit.CapitalReleaseCause(cause),
		Reinvested:     reinvested,
		Overstatement: credit.LoanCapitalOverstatement{
			OverstatementPct: overstatementPct,
			Description:      overstatementDescription,
		},
		Description: description,
		CreatedAt:   createdAt,
	}, nil
}

// CreateClearingHouseSettlement inserts s, assigning an ID and timestamp when absent (Vol. III Ch. 33).
func (m *MySQL) CreateClearingHouseSettlement(ctx context.Context, s credit.ClearingHouseSettlement) (credit.ClearingHouseSettlement, error) {
	if s.ID.IsZero() {
		s.ID = credit.NewClearingHouseSettlementID()
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO clearing_house_settlements ` +
		"(`id`, `name`, `description`, `total_claims`, `money_used`, `net_balance`, `created_at`) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := m.db.ExecContext(ctx, q,
		string(s.ID),
		s.Name,
		s.Description,
		s.TotalClaims,
		s.MoneyUsed,
		s.NetBalance,
		s.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.ClearingHouseSettlement{}, ErrAlreadyExists
		}
		return credit.ClearingHouseSettlement{}, err
	}
	return s, nil
}

// GetClearingHouseSettlement returns the record with id, or ErrNotFound.
func (m *MySQL) GetClearingHouseSettlement(ctx context.Context, id credit.ClearingHouseSettlementID) (credit.ClearingHouseSettlement, error) {
	const q = `SELECT ` +
		"`id`, `name`, `description`, `total_claims`, `money_used`, `net_balance`, `created_at` " +
		"FROM `clearing_house_settlements` WHERE `id` = ?"
	return scanClearingHouseSettlement(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListClearingHouseSettlements returns all stored records, newest first.
func (m *MySQL) ListClearingHouseSettlements(ctx context.Context) ([]credit.ClearingHouseSettlement, error) {
	const q = `SELECT ` +
		"`id`, `name`, `description`, `total_claims`, `money_used`, `net_balance`, `created_at` " +
		"FROM `clearing_house_settlements` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.ClearingHouseSettlement, 0)
	for rows.Next() {
		s, err := scanClearingHouseSettlement(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func scanClearingHouseSettlement(s rowScanner) (credit.ClearingHouseSettlement, error) {
	var (
		id          string
		name        string
		description string
		totalClaims int64
		moneyUsed   int64
		netBalance  int64
		createdAt   time.Time
	)
	err := s.Scan(&id, &name, &description, &totalClaims, &moneyUsed, &netBalance, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.ClearingHouseSettlement{}, ErrNotFound
		}
		return credit.ClearingHouseSettlement{}, err
	}
	return credit.ClearingHouseSettlement{
		ID:          credit.ClearingHouseSettlementID(id),
		Name:        name,
		Description: description,
		TotalClaims: totalClaims,
		MoneyUsed:   moneyUsed,
		NetBalance:  netBalance,
		CreatedAt:   createdAt,
	}, nil
}

// CreateNoteIssueConstraint inserts n, assigning an ID and timestamp when absent (Vol. III Ch. 34).
func (m *MySQL) CreateNoteIssueConstraint(ctx context.Context, n credit.NoteIssueConstraint) (credit.NoteIssueConstraint, error) {
	if n.ID.IsZero() {
		n.ID = credit.NewNoteIssueConstraintID()
	}
	if n.CreatedAt.IsZero() {
		n.CreatedAt = m.now().UTC()
	}
	const q = `INSERT INTO note_issue_constraints ` +
		"(`id`, `name`, `regime`, `max_notes`, `gold_backing`, `unbacked_limit`, `notes_beyond_max`, `description`, `created_at`) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := m.db.ExecContext(ctx, q,
		string(n.ID),
		n.Name,
		int(n.Regime),
		n.MaxNotes,
		n.GoldBacking,
		n.UnbackedLimit,
		n.NotesBeyondMax,
		n.Description,
		n.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return credit.NoteIssueConstraint{}, ErrAlreadyExists
		}
		return credit.NoteIssueConstraint{}, err
	}
	return n, nil
}

// GetNoteIssueConstraint returns the record with id, or ErrNotFound.
func (m *MySQL) GetNoteIssueConstraint(ctx context.Context, id credit.NoteIssueConstraintID) (credit.NoteIssueConstraint, error) {
	const q = `SELECT ` +
		"`id`, `name`, `regime`, `max_notes`, `gold_backing`, `unbacked_limit`, `notes_beyond_max`, `description`, `created_at` " +
		"FROM `note_issue_constraints` WHERE `id` = ?"
	return scanNoteIssueConstraint(m.db.QueryRowContext(ctx, q, string(id)))
}

// ListNoteIssueConstraints returns all stored records, newest first.
func (m *MySQL) ListNoteIssueConstraints(ctx context.Context) ([]credit.NoteIssueConstraint, error) {
	const q = `SELECT ` +
		"`id`, `name`, `regime`, `max_notes`, `gold_backing`, `unbacked_limit`, `notes_beyond_max`, `description`, `created_at` " +
		"FROM `note_issue_constraints` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]credit.NoteIssueConstraint, 0)
	for rows.Next() {
		n, err := scanNoteIssueConstraint(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func scanNoteIssueConstraint(s rowScanner) (credit.NoteIssueConstraint, error) {
	var (
		id             string
		name           string
		regime         int
		maxNotes       int64
		goldBacking    int64
		unbackedLimit  int64
		notesBeyondMax int64
		description    string
		createdAt      time.Time
	)
	err := s.Scan(&id, &name, &regime, &maxNotes, &goldBacking, &unbackedLimit, &notesBeyondMax, &description, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return credit.NoteIssueConstraint{}, ErrNotFound
		}
		return credit.NoteIssueConstraint{}, err
	}
	return credit.NoteIssueConstraint{
		ID:             credit.NoteIssueConstraintID(id),
		Name:           name,
		Regime:         credit.BankActRegime(regime),
		MaxNotes:       maxNotes,
		GoldBacking:    goldBacking,
		UnbackedLimit:  unbackedLimit,
		NotesBeyondMax: notesBeyondMax,
		Description:    description,
		CreatedAt:      createdAt,
	}, nil
}
