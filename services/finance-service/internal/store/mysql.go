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

// isDuplicate reports whether err is a MySQL duplicate-key error (1062).
func isDuplicate(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "1062") || strings.Contains(s, "Duplicate entry")
}
