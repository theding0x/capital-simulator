package store

import (
	"context"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// CreateCompetitionIllusion inserts i, assigning an ID and timestamp when absent (Vol. III Ch. 50).
func (m *MySQL) CreateCompetitionIllusion(ctx context.Context, i revenue.CompetitionIllusion) (revenue.CompetitionIllusion, error) {
	if i.ID == "" {
		i.ID = revenue.NewCompetitionIllusionID()
	}
	if i.CreatedAt.IsZero() {
		i.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `competition_illusions`" +
		" (`id`, `name`, `surface_appearance`, `real_relation`, `produced_by`, `created_at`)" +
		" VALUES (?,?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q, string(i.ID), i.Name, i.SurfaceAppearance, i.RealRelation, i.ProducedBy, i.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return revenue.CompetitionIllusion{}, ErrAlreadyExists
		}
		return revenue.CompetitionIllusion{}, err
	}
	return i, nil
}

// ListCompetitionIllusions returns all stored illusions, newest first (Vol. III Ch. 50).
func (m *MySQL) ListCompetitionIllusions(ctx context.Context) ([]revenue.CompetitionIllusion, error) {
	const q = "SELECT `id`, `name`, `surface_appearance`, `real_relation`, `produced_by`, `created_at`" +
		" FROM `competition_illusions` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.CompetitionIllusion, 0)
	for rows.Next() {
		var (
			id, name, surface, real, producedBy string
			createdAt                           time.Time
		)
		if err := rows.Scan(&id, &name, &surface, &real, &producedBy, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, revenue.CompetitionIllusion{
			ID:                revenue.CompetitionIllusionID(id),
			Name:              name,
			SurfaceAppearance: surface,
			RealRelation:      real,
			ProducedBy:        producedBy,
			CreatedAt:         createdAt,
		})
	}
	return out, rows.Err()
}

// CreateCostPriceFetish inserts f, assigning an ID and timestamp when absent (Vol. III Ch. 50).
func (m *MySQL) CreateCostPriceFetish(ctx context.Context, f revenue.CostPriceFetish) (revenue.CostPriceFetish, error) {
	if f.ID == "" {
		f.ID = revenue.NewCostPriceFetishID()
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `cost_price_fetishes`" +
		" (`id`, `constant_capital_bp`, `variable_capital_bp`, `cost_price_bp`, `profit_bp`, `actual_value_bp`, `created_at`)" +
		" VALUES (?,?,?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q, string(f.ID), f.ConstantCapitalBP, f.VariableCapitalBP, f.CostPriceBP, f.ProfitBP, f.ActualValueBP, f.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return revenue.CostPriceFetish{}, ErrAlreadyExists
		}
		return revenue.CostPriceFetish{}, err
	}
	return f, nil
}

// ListCostPriceFetishes returns all stored cost-price fetishes, newest first (Vol. III Ch. 50).
func (m *MySQL) ListCostPriceFetishes(ctx context.Context) ([]revenue.CostPriceFetish, error) {
	const q = "SELECT `id`, `constant_capital_bp`, `variable_capital_bp`, `cost_price_bp`, `profit_bp`, `actual_value_bp`, `created_at`" +
		" FROM `cost_price_fetishes` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.CostPriceFetish, 0)
	for rows.Next() {
		var (
			id                                       string
			constant, variable, cost, profit, actual int64
			createdAt                                time.Time
		)
		if err := rows.Scan(&id, &constant, &variable, &cost, &profit, &actual, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, revenue.CostPriceFetish{
			ID:                revenue.CostPriceFetishID(id),
			ConstantCapitalBP: constant,
			VariableCapitalBP: variable,
			CostPriceBP:       cost,
			ProfitBP:          profit,
			ActualValueBP:     actual,
			CreatedAt:         createdAt,
		})
	}
	return out, rows.Err()
}

// CreateAnnualRevenueAccount inserts a, assigning an ID and timestamp when absent (Vol. III Ch. 50).
func (m *MySQL) CreateAnnualRevenueAccount(ctx context.Context, a revenue.AnnualRevenueAccount) (revenue.AnnualRevenueAccount, error) {
	if a.ID == "" {
		a.ID = revenue.NewAnnualRevenueAccountID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `annual_revenue_accounts`" +
		" (`id`, `new_value_created_lm`, `wages_share_lm`, `profit_and_rent_share_lm`, `constant_capital_lm`, `total_social_product_lm`, `created_at`)" +
		" VALUES (?,?,?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q, string(a.ID), a.NewValueCreatedLM, a.WagesShareLM, a.ProfitAndRentShareLM, a.ConstantCapitalLM, a.TotalSocialProductLM, a.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return revenue.AnnualRevenueAccount{}, ErrAlreadyExists
		}
		return revenue.AnnualRevenueAccount{}, err
	}
	return a, nil
}

// ListAnnualRevenueAccounts returns all stored accounts, newest first (Vol. III Ch. 50).
func (m *MySQL) ListAnnualRevenueAccounts(ctx context.Context) ([]revenue.AnnualRevenueAccount, error) {
	const q = "SELECT `id`, `new_value_created_lm`, `wages_share_lm`, `profit_and_rent_share_lm`, `constant_capital_lm`, `total_social_product_lm`, `created_at`" +
		" FROM `annual_revenue_accounts` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.AnnualRevenueAccount, 0)
	for rows.Next() {
		var (
			id                                           string
			newValue, wages, profitRent, constant, total int64
			createdAt                                    time.Time
		)
		if err := rows.Scan(&id, &newValue, &wages, &profitRent, &constant, &total, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, revenue.AnnualRevenueAccount{
			ID:                   revenue.AnnualRevenueAccountID(id),
			NewValueCreatedLM:    newValue,
			WagesShareLM:         wages,
			ProfitAndRentShareLM: profitRent,
			ConstantCapitalLM:    constant,
			TotalSocialProductLM: total,
			CreatedAt:            createdAt,
		})
	}
	return out, rows.Err()
}

// CreateValueAppearanceContrast inserts c, assigning an ID and timestamp when absent (Vol. III Ch. 50).
func (m *MySQL) CreateValueAppearanceContrast(ctx context.Context, c revenue.ValueAppearanceContrast) (revenue.ValueAppearanceContrast, error) {
	if c.ID == "" {
		c.ID = revenue.NewValueAppearanceContrastID()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `value_appearance_contrasts`" +
		" (`id`, `illusion_id`, `surface`, `reality`, `created_at`)" +
		" VALUES (?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q, string(c.ID), c.IllusionID, c.Surface, c.Reality, c.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return revenue.ValueAppearanceContrast{}, ErrAlreadyExists
		}
		return revenue.ValueAppearanceContrast{}, err
	}
	return c, nil
}

// ListValueAppearanceContrasts returns all stored contrasts, newest first (Vol. III Ch. 50).
func (m *MySQL) ListValueAppearanceContrasts(ctx context.Context) ([]revenue.ValueAppearanceContrast, error) {
	const q = "SELECT `id`, `illusion_id`, `surface`, `reality`, `created_at`" +
		" FROM `value_appearance_contrasts` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.ValueAppearanceContrast, 0)
	for rows.Next() {
		var (
			id, illusionID, surface, reality string
			createdAt                        time.Time
		)
		if err := rows.Scan(&id, &illusionID, &surface, &reality, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, revenue.ValueAppearanceContrast{
			ID:         revenue.ValueAppearanceContrastID(id),
			IllusionID: illusionID,
			Surface:    surface,
			Reality:    reality,
			CreatedAt:  createdAt,
		})
	}
	return out, rows.Err()
}
