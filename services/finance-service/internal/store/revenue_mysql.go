package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/theding0x/capital-simulator/services/finance-service/internal/revenue"
)

// CreateRevenueStream inserts s, assigning an ID and timestamp when absent (Vol. III Ch. 48).
func (m *MySQL) CreateRevenueStream(ctx context.Context, s revenue.RevenueStream) (revenue.RevenueStream, error) {
	if s.ID == "" {
		s.ID = revenue.NewRevenueStreamID()
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = m.now().UTC()
	}
	fetished := 0
	if s.IsFetishised {
		fetished = 1
	}
	const q = "INSERT INTO `revenue_streams`" +
		" (`id`, `source`, `apparent_revenue_bp`, `actual_source_bp`, `is_fetishised`, `created_at`)" +
		" VALUES (?,?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q,
		string(s.ID), int(s.Source), s.ApparentRevenueBP, s.ActualSourceBP, fetished, s.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return revenue.RevenueStream{}, ErrAlreadyExists
		}
		return revenue.RevenueStream{}, err
	}
	return s, nil
}

// ListRevenueStreams returns all stored revenue streams, newest first (Vol. III Ch. 48).
func (m *MySQL) ListRevenueStreams(ctx context.Context) ([]revenue.RevenueStream, error) {
	const q = "SELECT `id`, `source`, `apparent_revenue_bp`, `actual_source_bp`, `is_fetishised`, `created_at`" +
		" FROM `revenue_streams` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.RevenueStream, 0)
	for rows.Next() {
		s, err := scanRevenueStream(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func scanRevenueStream(s rowScanner) (revenue.RevenueStream, error) {
	var (
		id        string
		source    int
		apparent  int64
		actual    int64
		fetished  int
		createdAt time.Time
	)
	if err := s.Scan(&id, &source, &apparent, &actual, &fetished, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return revenue.RevenueStream{}, ErrNotFound
		}
		return revenue.RevenueStream{}, err
	}
	return revenue.RevenueStream{
		ID:                revenue.RevenueStreamID(id),
		Source:            revenue.RevenueSourceKind(source),
		ApparentRevenueBP: apparent,
		ActualSourceBP:    actual,
		IsFetishised:      fetished != 0,
		CreatedAt:         createdAt,
	}, nil
}

// CreateTrinityFormula inserts tf, assigning an ID and timestamp when absent (Vol. III Ch. 48).
func (m *MySQL) CreateTrinityFormula(ctx context.Context, tf revenue.TrinityFormula) (revenue.TrinityFormula, error) {
	if tf.ID == "" {
		tf.ID = revenue.NewTrinityFormulaID()
	}
	if tf.CreatedAt.IsZero() {
		tf.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `trinity_formulas`" +
		" (`id`, `capital_stream_id`, `land_stream_id`, `labour_stream_id`, `total_surplus_value_bp`, `total_apparent_revenue_bp`, `created_at`)" +
		" VALUES (?,?,?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q,
		string(tf.ID), tf.CapitalStreamID, tf.LandStreamID, tf.LabourStreamID,
		tf.TotalSurplusValueBP, tf.TotalApparentRevenueBP, tf.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return revenue.TrinityFormula{}, ErrAlreadyExists
		}
		return revenue.TrinityFormula{}, err
	}
	return tf, nil
}

// GetTrinityFormula returns the trinity formula with the given ID, or ErrNotFound (Vol. III Ch. 48).
func (m *MySQL) GetTrinityFormula(ctx context.Context, id revenue.TrinityFormulaID) (revenue.TrinityFormula, error) {
	const q = "SELECT `id`, `capital_stream_id`, `land_stream_id`, `labour_stream_id`, `total_surplus_value_bp`, `total_apparent_revenue_bp`, `created_at`" +
		" FROM `trinity_formulas` WHERE `id` = ?"
	return scanTrinityFormula(m.db.QueryRowContext(ctx, q, string(id)))
}

func scanTrinityFormula(s rowScanner) (revenue.TrinityFormula, error) {
	var (
		id        string
		capitalID string
		landID    string
		labourID  string
		totalSV   int64
		totalApp  int64
		createdAt time.Time
	)
	if err := s.Scan(&id, &capitalID, &landID, &labourID, &totalSV, &totalApp, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return revenue.TrinityFormula{}, ErrNotFound
		}
		return revenue.TrinityFormula{}, err
	}
	return revenue.TrinityFormula{
		ID:                     revenue.TrinityFormulaID(id),
		CapitalStreamID:        capitalID,
		LandStreamID:           landID,
		LabourStreamID:         labourID,
		TotalSurplusValueBP:    totalSV,
		TotalApparentRevenueBP: totalApp,
		CreatedAt:              createdAt,
	}, nil
}

// CreateRevenueFetishForm inserts f, assigning an ID and timestamp when absent (Vol. III Ch. 48).
func (m *MySQL) CreateRevenueFetishForm(ctx context.Context, f revenue.RevenueFetishForm) (revenue.RevenueFetishForm, error) {
	if f.ID == "" {
		f.ID = revenue.NewRevenueFetishFormID()
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = m.now().UTC()
	}
	const q = "INSERT INTO `revenue_fetish_forms`" +
		" (`id`, `source`, `surface_formula`, `real_relation`, `mystification_kind`, `created_at`)" +
		" VALUES (?,?,?,?,?,?)"
	_, err := m.db.ExecContext(ctx, q,
		string(f.ID), int(f.Source), f.SurfaceFormula, f.RealRelation, f.MystificationKind, f.CreatedAt,
	)
	if err != nil {
		if isDuplicate(err) {
			return revenue.RevenueFetishForm{}, ErrAlreadyExists
		}
		return revenue.RevenueFetishForm{}, err
	}
	return f, nil
}

// ListRevenueFetishForms returns all stored revenue fetish forms, newest first (Vol. III Ch. 48).
func (m *MySQL) ListRevenueFetishForms(ctx context.Context) ([]revenue.RevenueFetishForm, error) {
	const q = "SELECT `id`, `source`, `surface_formula`, `real_relation`, `mystification_kind`, `created_at`" +
		" FROM `revenue_fetish_forms` ORDER BY `created_at` DESC, `id` ASC"
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]revenue.RevenueFetishForm, 0)
	for rows.Next() {
		f, err := scanRevenueFetishForm(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

func scanRevenueFetishForm(s rowScanner) (revenue.RevenueFetishForm, error) {
	var (
		id        string
		source    int
		surface   string
		realRel   string
		myst      string
		createdAt time.Time
	)
	if err := s.Scan(&id, &source, &surface, &realRel, &myst, &createdAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return revenue.RevenueFetishForm{}, ErrNotFound
		}
		return revenue.RevenueFetishForm{}, err
	}
	return revenue.RevenueFetishForm{
		ID:                revenue.RevenueFetishFormID(id),
		Source:            revenue.RevenueSourceKind(source),
		SurfaceFormula:    surface,
		RealRelation:      realRel,
		MystificationKind: myst,
		CreatedAt:         createdAt,
	}, nil
}
