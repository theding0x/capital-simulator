package store

import (
	"context"
	"database/sql"

	repro "github.com/theding0x/capital-simulator/services/simulation-engine/internal/reproduction"
)

// CreateSurplusCirculation inserts a new SurplusCirculation row. Source
// entries are ignored on creation; use RecordRealisationSource to append them.
func (m *MySQL) CreateSurplusCirculation(ctx context.Context, sc repro.SurplusCirculation) (repro.SurplusCirculation, error) {
	if sc.TotalSurplusPence < 0 {
		return repro.SurplusCirculation{}, repro.ErrNegativeSurplus
	}
	if sc.ID.IsZero() {
		sc.ID = repro.NewSurplusCirculationID()
	}
	sc.RealisedSurplusPence = 0
	sc.UnrealisedSurplusPence = sc.TotalSurplusPence
	sc.RealisationSources = []repro.RealisationSourceEntry{}

	const q = `INSERT INTO surplus_circulations
		(id, period, total_surplus_pence, realised_surplus_pence,
		 unrealised_surplus_pence, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err := m.db.ExecContext(ctx, q,
		string(sc.ID), sc.Period,
		int64(sc.TotalSurplusPence), int64(sc.RealisedSurplusPence),
		int64(sc.UnrealisedSurplusPence), m.now(),
	)
	if err != nil {
		if isDuplicateKey(err) {
			return repro.SurplusCirculation{}, ErrAlreadyExists
		}
		return repro.SurplusCirculation{}, err
	}
	return sc, nil
}

// GetSurplusCirculation returns a SurplusCirculation with all source entries.
func (m *MySQL) GetSurplusCirculation(ctx context.Context, id repro.SurplusCirculationID) (repro.SurplusCirculation, error) {
	const q = `SELECT id, period, total_surplus_pence, realised_surplus_pence,
		unrealised_surplus_pence
		FROM surplus_circulations WHERE id = ?`
	row := m.db.QueryRowContext(ctx, q, string(id))
	sc, err := scanSurplusCirculation(row.Scan)
	if err != nil {
		if err == sql.ErrNoRows {
			return repro.SurplusCirculation{}, ErrNotFound
		}
		return repro.SurplusCirculation{}, err
	}
	sources, err := m.listRealisationSourcesForID(ctx, id)
	if err != nil {
		return repro.SurplusCirculation{}, err
	}
	sc.RealisationSources = sources
	return sc, nil
}

// ListSurplusCirculations returns all SurplusCirculation records ordered by
// creation time ascending, with source entries assembled.
func (m *MySQL) ListSurplusCirculations(ctx context.Context) ([]repro.SurplusCirculation, error) {
	const q = `SELECT id, period, total_surplus_pence, realised_surplus_pence,
		unrealised_surplus_pence
		FROM surplus_circulations ORDER BY created_at ASC`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []repro.SurplusCirculation
	for rows.Next() {
		sc, err := scanSurplusCirculation(rows.Scan)
		if err != nil {
			return nil, err
		}
		sources, err := m.listRealisationSourcesForID(ctx, sc.ID)
		if err != nil {
			return nil, err
		}
		sc.RealisationSources = sources
		out = append(out, sc)
	}
	if out == nil {
		out = []repro.SurplusCirculation{}
	}
	return out, nil
}

// RecordRealisationSource appends a source entry and updates the parent
// surplus circulation pence counters atomically within a transaction.
func (m *MySQL) RecordRealisationSource(ctx context.Context, id repro.SurplusCirculationID, entry repro.RealisationSourceEntry) (repro.RealisationSourceEntry, error) {
	if entry.Pence < 0 {
		return repro.RealisationSourceEntry{}, repro.ErrNegativePence
	}
	if err := entry.Source.Validate(); err != nil {
		return repro.RealisationSourceEntry{}, err
	}
	entry.SurplusCirculationID = id
	if entry.Source == repro.SourceCredit {
		entry.DeferredRealisationPence = entry.Pence
	}

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return repro.RealisationSourceEntry{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var totalPence, realisedPence int64
	const lockQ = `SELECT total_surplus_pence, realised_surplus_pence
		FROM surplus_circulations WHERE id = ? FOR UPDATE`
	if err := tx.QueryRowContext(ctx, lockQ, string(id)).Scan(&totalPence, &realisedPence); err != nil {
		if err == sql.ErrNoRows {
			return repro.RealisationSourceEntry{}, ErrNotFound
		}
		return repro.RealisationSourceEntry{}, err
	}

	newRealised := realisedPence + int64(entry.Pence)
	if newRealised > totalPence {
		return repro.RealisationSourceEntry{}, repro.ErrRealisationExceedsTotal
	}
	newUnrealised := totalPence - newRealised

	const insQ = `INSERT INTO realisation_source_entries
		(surplus_circulation_id, source, pence, deferred_realisation_pence)
		VALUES (?, ?, ?, ?)`
	if _, err := tx.ExecContext(ctx, insQ,
		string(id), string(entry.Source),
		int64(entry.Pence), int64(entry.DeferredRealisationPence),
	); err != nil {
		return repro.RealisationSourceEntry{}, err
	}

	const updQ = `UPDATE surplus_circulations
		SET realised_surplus_pence = ?, unrealised_surplus_pence = ?
		WHERE id = ?`
	if _, err := tx.ExecContext(ctx, updQ, newRealised, newUnrealised, string(id)); err != nil {
		return repro.RealisationSourceEntry{}, err
	}

	if err := tx.Commit(); err != nil {
		return repro.RealisationSourceEntry{}, err
	}
	return entry, nil
}

// SnapshotAggregate returns a SocialCapitalAggregate for the given period.
// If an explicit aggregate row exists it is returned; otherwise
// TotalAnnualSurplusPence is derived by summing all SurplusCirculation rows
// whose period matches.
func (m *MySQL) SnapshotAggregate(ctx context.Context, period string) (repro.SocialCapitalAggregate, error) {
	const aggQ = `SELECT period, total_advanced_pence, total_annual_output_pence,
		total_annual_surplus_pence, department_i_share_bp, department_ii_share_bp
		FROM social_capital_aggregates WHERE period = ? LIMIT 1`
	var agg repro.SocialCapitalAggregate
	var totalAdv, totalOut, totalSurp int64
	err := m.db.QueryRowContext(ctx, aggQ, period).Scan(
		&agg.Period, &totalAdv, &totalOut, &totalSurp,
		&agg.DepartmentIShareBP, &agg.DepartmentIIShareBP,
	)
	if err == nil {
		agg.TotalAdvancedPence = repro.Pence(totalAdv)
		agg.TotalAnnualOutputPence = repro.Pence(totalOut)
		agg.TotalAnnualSurplusPence = repro.Pence(totalSurp)
		return agg, nil
	}
	if err != sql.ErrNoRows {
		return repro.SocialCapitalAggregate{}, err
	}

	const deriveQ = `SELECT COALESCE(SUM(total_surplus_pence), 0)
		FROM surplus_circulations WHERE period = ?`
	var sumPence int64
	if err := m.db.QueryRowContext(ctx, deriveQ, period).Scan(&sumPence); err != nil {
		return repro.SocialCapitalAggregate{}, err
	}
	return repro.SocialCapitalAggregate{
		Period:                  period,
		TotalAnnualSurplusPence: repro.Pence(sumPence),
	}, nil
}

// ListRealisationPuzzles returns all canonical RealisationPuzzle rows.
func (m *MySQL) ListRealisationPuzzles(ctx context.Context) ([]repro.RealisationPuzzle, error) {
	const q = `SELECT surplus_circulation_id, puzzle_statement, resolved_in_chapter
		FROM realisation_puzzles`
	rows, err := m.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []repro.RealisationPuzzle
	for rows.Next() {
		var p repro.RealisationPuzzle
		var scID string
		if err := rows.Scan(&scID, &p.PuzzleStatement, &p.ResolvedInChapter); err != nil {
			return nil, err
		}
		p.SurplusCirculationID = repro.SurplusCirculationID(scID)
		out = append(out, p)
	}
	if out == nil {
		out = []repro.RealisationPuzzle{}
	}
	return out, nil
}

// --- helpers ----------------------------------------------------------------

func scanSurplusCirculation(scan func(...any) error) (repro.SurplusCirculation, error) {
	var sc repro.SurplusCirculation
	var id string
	var total, realised, unrealised int64
	if err := scan(&id, &sc.Period, &total, &realised, &unrealised); err != nil {
		return repro.SurplusCirculation{}, err
	}
	sc.ID = repro.SurplusCirculationID(id)
	sc.TotalSurplusPence = repro.Pence(total)
	sc.RealisedSurplusPence = repro.Pence(realised)
	sc.UnrealisedSurplusPence = repro.Pence(unrealised)
	return sc, nil
}

func (m *MySQL) listRealisationSourcesForID(ctx context.Context, id repro.SurplusCirculationID) ([]repro.RealisationSourceEntry, error) {
	const q = `SELECT source, pence, deferred_realisation_pence
		FROM realisation_source_entries
		WHERE surplus_circulation_id = ?
		ORDER BY id ASC`
	rows, err := m.db.QueryContext(ctx, q, string(id))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []repro.RealisationSourceEntry
	for rows.Next() {
		var e repro.RealisationSourceEntry
		var src string
		var pence, deferred int64
		if err := rows.Scan(&src, &pence, &deferred); err != nil {
			return nil, err
		}
		e.SurplusCirculationID = id
		e.Source = repro.RealisationSource(src)
		e.Pence = repro.Pence(pence)
		e.DeferredRealisationPence = repro.Pence(deferred)
		out = append(out, e)
	}
	if out == nil {
		out = []repro.RealisationSourceEntry{}
	}
	return out, nil
}
