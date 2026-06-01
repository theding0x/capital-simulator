package store

import (
	"context"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// CreateProductionRelationBasis inserts b, assigning an ID and timestamp when absent (Vol. III Ch. 51).
func (m *MySQL) CreateProductionRelationBasis(ctx context.Context, b revenue.ProductionRelationBasis) (revenue.ProductionRelationBasis, error) {
	if b.ID == "" {
		b.ID = revenue.NewProductionRelationBasisID()
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = m.now().UTC()
	}
	spec := 0
	if b.IsHistoricallySpecific {
		spec = 1
	}
	const q = "INSERT INTO `production_relation_bases`" +
		" (`id`, `name`, `characteristic_feature`, `is_historically_specific`, `created_at`)" +
		" VALUES (?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q, string(b.ID), b.Name, b.CharacteristicFeature, spec, b.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return revenue.ProductionRelationBasis{}, ErrAlreadyExists
		}
		return revenue.ProductionRelationBasis{}, err
	}
	return b, nil
}

// ListProductionRelationBases returns all stored bases, newest first (Vol. III Ch. 51).
func (m *MySQL) ListProductionRelationBases(ctx context.Context) ([]revenue.ProductionRelationBasis, error) {
	const q = "SELECT `id`, `name`, `characteristic_feature`, `is_historically_specific`, `created_at`" +
		" FROM `production_relation_bases` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.ProductionRelationBasis, 0)
	for rows.Next() {
		var (
			id, name, feature string
			spec              int
			createdAt         time.Time
		)
		if err := rows.Scan(&id, &name, &feature, &spec, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, revenue.ProductionRelationBasis{
			ID:                     revenue.ProductionRelationBasisID(id),
			Name:                   name,
			CharacteristicFeature:  feature,
			IsHistoricallySpecific: spec != 0,
			CreatedAt:              createdAt,
		})
	}
	return out, rows.Err()
}

// CreateDistributionRelation inserts d, assigning an ID and timestamp when absent (Vol. III Ch. 51).
func (m *MySQL) CreateDistributionRelation(ctx context.Context, d revenue.DistributionRelation) (revenue.DistributionRelation, error) {
	if d.ID == "" {
		d.ID = revenue.NewDistributionRelationID()
	}
	if d.CreatedAt.IsZero() {
		d.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `distribution_relations`" +
		" (`id`, `form`, `production_basis_id`, `apparent_character`, `real_character`, `corresponds_to_surplus_form`, `created_at`)" +
		" VALUES (?,?,?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q, string(d.ID), d.Form, d.ProductionBasisID, d.ApparentCharacter, d.RealCharacter, d.CorrespondsToSurplusForm, d.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return revenue.DistributionRelation{}, ErrAlreadyExists
		}
		return revenue.DistributionRelation{}, err
	}
	return d, nil
}

// ListDistributionRelations returns all stored relations, newest first (Vol. III Ch. 51).
func (m *MySQL) ListDistributionRelations(ctx context.Context) ([]revenue.DistributionRelation, error) {
	const q = "SELECT `id`, `form`, `production_basis_id`, `apparent_character`, `real_character`, `corresponds_to_surplus_form`, `created_at`" +
		" FROM `distribution_relations` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.DistributionRelation, 0)
	for rows.Next() {
		var (
			id, form, basisID, apparent, real, corresponds string
			createdAt                                      time.Time
		)
		if err := rows.Scan(&id, &form, &basisID, &apparent, &real, &corresponds, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, revenue.DistributionRelation{
			ID:                       revenue.DistributionRelationID(id),
			Form:                     form,
			ProductionBasisID:        basisID,
			ApparentCharacter:        apparent,
			RealCharacter:            real,
			CorrespondsToSurplusForm: corresponds,
			CreatedAt:                createdAt,
		})
	}
	return out, rows.Err()
}

// CreateHistoricalTransience inserts t, assigning an ID and timestamp when absent (Vol. III Ch. 51).
func (m *MySQL) CreateHistoricalTransience(ctx context.Context, t revenue.HistoricalTransience) (revenue.HistoricalTransience, error) {
	if t.ID == "" {
		t.ID = revenue.NewHistoricalTransienceID()
	}
	if t.CreatedAt.IsZero() {
		t.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `historical_transiences`" +
		" (`id`, `distribution_form_id`, `mode_of_production`, `precondition`, `succeeding_form`, `created_at`)" +
		" VALUES (?,?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q, string(t.ID), t.DistributionFormID, t.ModeOfProduction, t.Precondition, t.SucceedingForm, t.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return revenue.HistoricalTransience{}, ErrAlreadyExists
		}
		return revenue.HistoricalTransience{}, err
	}
	return t, nil
}

// ListHistoricalTransiences returns all stored records, newest first (Vol. III Ch. 51).
func (m *MySQL) ListHistoricalTransiences(ctx context.Context) ([]revenue.HistoricalTransience, error) {
	const q = "SELECT `id`, `distribution_form_id`, `mode_of_production`, `precondition`, `succeeding_form`, `created_at`" +
		" FROM `historical_transiences` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.HistoricalTransience, 0)
	for rows.Next() {
		var (
			id, formID, mode, precond, succeeding string
			createdAt                             time.Time
		)
		if err := rows.Scan(&id, &formID, &mode, &precond, &succeeding, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, revenue.HistoricalTransience{
			ID:                 revenue.HistoricalTransienceID(id),
			DistributionFormID: formID,
			ModeOfProduction:   mode,
			Precondition:       precond,
			SucceedingForm:     succeeding,
			CreatedAt:          createdAt,
		})
	}
	return out, rows.Err()
}

// CreateProductiveForceContradiction inserts c, assigning an ID and timestamp when absent (Vol. III Ch. 51).
func (m *MySQL) CreateProductiveForceContradiction(ctx context.Context, c revenue.ProductiveForceContradiction) (revenue.ProductiveForceContradiction, error) {
	if c.ID == "" {
		c.ID = revenue.NewProductiveForceContradictionID()
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `productive_force_contradictions`" +
		" (`id`, `productive_force`, `social_form`, `contradiction_nature`, `crisis_signal`, `created_at`)" +
		" VALUES (?,?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q, string(c.ID), c.ProductiveForce, c.SocialForm, c.ContradictionNature, c.CrisisSignal, c.CreatedAt)
	if err != nil {
		if isDuplicate(err) {
			return revenue.ProductiveForceContradiction{}, ErrAlreadyExists
		}
		return revenue.ProductiveForceContradiction{}, err
	}
	return c, nil
}

// ListProductiveForceContradictions returns all stored records, newest first (Vol. III Ch. 51).
func (m *MySQL) ListProductiveForceContradictions(ctx context.Context) ([]revenue.ProductiveForceContradiction, error) {
	const q = "SELECT `id`, `productive_force`, `social_form`, `contradiction_nature`, `crisis_signal`, `created_at`" +
		" FROM `productive_force_contradictions` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.ProductiveForceContradiction, 0)
	for rows.Next() {
		var (
			id, force, social, nature, signal string
			createdAt                         time.Time
		)
		if err := rows.Scan(&id, &force, &social, &nature, &signal, &createdAt); err != nil {
			return nil, err
		}
		out = append(out, revenue.ProductiveForceContradiction{
			ID:                  revenue.ProductiveForceContradictionID(id),
			ProductiveForce:     force,
			SocialForm:          social,
			ContradictionNature: nature,
			CrisisSignal:        signal,
			CreatedAt:           createdAt,
		})
	}
	return out, rows.Err()
}
