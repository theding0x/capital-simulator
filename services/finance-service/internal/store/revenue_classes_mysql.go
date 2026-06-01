package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// CreateSocialClass inserts c, assigning an ID and timestamp when absent (Vol. III Ch. 52).
func (m *MySQL) CreateSocialClass(ctx context.Context, c revenue.SocialClass) (revenue.SocialClass, error) {
	if c.ID == "" {
		c.ID = revenue.NewSocialClassID()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `social_classes`" +
		" (`id`, `name`, `relation_to_means_of_production`, `primary_revenue_form`, `historical_precondition`, `created_at`)" +
		" VALUES (?,?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q, string(c.ID), c.Name, c.RelationToMeansOfProduction, c.PrimaryRevenueForm, c.HistoricalPrecondition, c.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return revenue.SocialClass{}, ErrAlreadyExists
		}
		return revenue.SocialClass{}, err
	}
	return c, nil
}

// GetSocialClass returns the class with the given ID, or ErrNotFound (Vol. III Ch. 52).
func (m *MySQL) GetSocialClass(ctx context.Context, id revenue.SocialClassID) (revenue.SocialClass, error) {
	const q = "SELECT `id`, `name`, `relation_to_means_of_production`, `primary_revenue_form`, `historical_precondition`, `created_at`" +
		" FROM `social_classes` WHERE `id` = ?"
	var (
		cid, name, relation, form, precond string
		createdAt                          time.Time
	)
	err := m.db.QueryRowContext(ctx, q, string(id)).Scan(&cid, &name, &relation, &form, &precond, &createdAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return revenue.SocialClass{}, ErrNotFound
		}
		return revenue.SocialClass{}, err
	}
	return revenue.SocialClass{
		ID:                          revenue.SocialClassID(cid),
		Name:                        name,
		RelationToMeansOfProduction: relation,
		PrimaryRevenueForm:          form,
		HistoricalPrecondition:      precond,
		CreatedAt:                   createdAt,
	}, nil
}

// ListSocialClasses returns all stored classes, newest first (Vol. III Ch. 52).
func (m *MySQL) ListSocialClasses(ctx context.Context) ([]revenue.SocialClass, error) {
	const q = "SELECT `id`, `name`, `relation_to_means_of_production`, `primary_revenue_form`, `historical_precondition`, `created_at`" +
		" FROM `social_classes` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.SocialClass, 0)
	for rows.Next() {
		var (
			id, name, relation, form, precond string
			createdAt                         time.Time
		)
		if err := rows.Scan(&id, &name, &relation, &form, &precond, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, revenue.SocialClass{
			ID:                          revenue.SocialClassID(id),
			Name:                        name,
			RelationToMeansOfProduction: relation,
			PrimaryRevenueForm:          form,
			HistoricalPrecondition:      precond,
			CreatedAt:                   createdAt,
		})
	}
	return out, rows.Err()
}

// CreateClassIncomeSource inserts s, assigning an ID and timestamp when absent (Vol. III Ch. 52).
func (m *MySQL) CreateClassIncomeSource(ctx context.Context, s revenue.ClassIncomeSource) (revenue.ClassIncomeSource, error) {
	if s.ID == "" {
		s.ID = revenue.NewClassIncomeSourceID()
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `class_income_sources`" +
		" (`id`, `class_id`, `revenue_form`, `surplus_value_component_bp`, `created_at`)" +
		" VALUES (?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q, string(s.ID), s.ClassID, s.RevenueForm, s.SurplusValueComponentBP, s.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return revenue.ClassIncomeSource{}, ErrAlreadyExists
		}
		return revenue.ClassIncomeSource{}, err
	}
	return s, nil
}

// ListClassIncomeSources returns all stored income sources, newest first (Vol. III Ch. 52).
func (m *MySQL) ListClassIncomeSources(ctx context.Context) ([]revenue.ClassIncomeSource, error) {
	const q = "SELECT `id`, `class_id`, `revenue_form`, `surplus_value_component_bp`, `created_at`" +
		" FROM `class_income_sources` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.ClassIncomeSource, 0)
	for rows.Next() {
		var (
			id, classID, form string
			component         int64
			createdAt         time.Time
		)
		if err := rows.Scan(&id, &classID, &form, &component, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, revenue.ClassIncomeSource{
			ID:                      revenue.ClassIncomeSourceID(id),
			ClassID:                 classID,
			RevenueForm:             form,
			SurplusValueComponentBP: component,
			CreatedAt:               createdAt,
		})
	}
	return out, rows.Err()
}

// CreateClassTendency inserts t, assigning an ID and timestamp when absent (Vol. III Ch. 52).
func (m *MySQL) CreateClassTendency(ctx context.Context, t revenue.ClassTendency) (revenue.ClassTendency, error) {
	if t.ID == "" {
		t.ID = revenue.NewClassTendencyID()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `class_tendencies`" +
		" (`id`, `class_id`, `description`, `direction`, `created_at`)" +
		" VALUES (?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q, string(t.ID), t.ClassID, t.Description, t.Direction, t.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return revenue.ClassTendency{}, ErrAlreadyExists
		}
		return revenue.ClassTendency{}, err
	}
	return t, nil
}

// ListClassTendencies returns all stored tendencies, newest first (Vol. III Ch. 52).
func (m *MySQL) ListClassTendencies(ctx context.Context) ([]revenue.ClassTendency, error) {
	const q = "SELECT `id`, `class_id`, `description`, `direction`, `created_at`" +
		" FROM `class_tendencies` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.ClassTendency, 0)
	for rows.Next() {
		var (
			id, classID, description, direction string
			createdAt                           time.Time
		)
		if err := rows.Scan(&id, &classID, &description, &direction, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, revenue.ClassTendency{
			ID:          revenue.ClassTendencyID(id),
			ClassID:     classID,
			Description: description,
			Direction:   direction,
			CreatedAt:   createdAt,
		})
	}
	return out, rows.Err()
}

// CreateClassAmbiguity inserts a, assigning an ID and timestamp when absent (Vol. III Ch. 52).
func (m *MySQL) CreateClassAmbiguity(ctx context.Context, a revenue.ClassAmbiguity) (revenue.ClassAmbiguity, error) {
	if a.ID == "" {
		a.ID = revenue.NewClassAmbiguityID()
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `class_ambiguities`" +
		" (`id`, `intermediate_group`, `income_source`, `why_not_a_class`, `created_at`)" +
		" VALUES (?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q, string(a.ID), a.IntermediateGroup, a.IncomeSource, a.WhyNotAClass, a.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return revenue.ClassAmbiguity{}, ErrAlreadyExists
		}
		return revenue.ClassAmbiguity{}, err
	}
	return a, nil
}

// ListClassAmbiguities returns all stored ambiguities, newest first (Vol. III Ch. 52).
func (m *MySQL) ListClassAmbiguities(ctx context.Context) ([]revenue.ClassAmbiguity, error) {
	const q = "SELECT `id`, `intermediate_group`, `income_source`, `why_not_a_class`, `created_at`" +
		" FROM `class_ambiguities` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.ClassAmbiguity, 0)
	for rows.Next() {
		var (
			id, group, income, why string
			createdAt              time.Time
		)
		if err := rows.Scan(&id, &group, &income, &why, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, revenue.ClassAmbiguity{
			ID:                revenue.ClassAmbiguityID(id),
			IntermediateGroup: group,
			IncomeSource:      income,
			WhyNotAClass:      why,
			CreatedAt:         createdAt,
		})
	}
	return out, rows.Err()
}
