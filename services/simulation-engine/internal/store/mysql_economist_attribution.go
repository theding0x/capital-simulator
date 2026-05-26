package store

// Vol. II Ch. 10 — Theories of Fixed and Circulating Capital.

import (
	"context"
	"database/sql"
	"strings"

	"github.com/theding0x/capital-simulator/services/simulation-engine/internal/circulation"
)

// ListEconomistAttributions returns economist-attribution reference records,
// optionally filtered by theorist (case-insensitive). An empty theorist string
// returns all records. Records are static reference data seeded at migration
// time and never mutated at runtime.
func (m *MySQL) ListEconomistAttributions(ctx context.Context, theorist string) ([]circulation.EconomistAttribution, error) {
	const baseQ = `
		SELECT ea.id, ea.concept, ea.theorist, ea.edition_year, ea.anticipates,
		       eae.error_code
		FROM economist_attributions ea
		LEFT JOIN economist_attribution_errors eae ON eae.attribution_id = ea.id
		ORDER BY ea.edition_year ASC, ea.id ASC`

	const filteredQ = `
		SELECT ea.id, ea.concept, ea.theorist, ea.edition_year, ea.anticipates,
		       eae.error_code
		FROM economist_attributions ea
		LEFT JOIN economist_attribution_errors eae ON eae.attribution_id = ea.id
		WHERE ea.theorist = ?
		ORDER BY ea.edition_year ASC, ea.id ASC`

	var (
		rows *sql.Rows
		err  error
	)
	if strings.TrimSpace(theorist) == "" {
		rows, err = m.db.QueryContext(ctx, baseQ)
	} else {
		rows, err = m.db.QueryContext(ctx, filteredQ, theorist)
	}
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	// One row per error entry due to the LEFT JOIN; fold into structs.
	type key = circulation.EconomistAttributionID
	var order []key
	byID := map[key]*circulation.EconomistAttribution{}

	for rows.Next() {
		var (
			id          string
			concept     string
			thr         string
			editionYear int64
			anticipates string
			errorCode   *string
		)
		if scanErr := rows.Scan(&id, &concept, &thr, &editionYear, &anticipates, &errorCode); scanErr != nil {
			return nil, scanErr
		}
		k := circulation.EconomistAttributionID(id)
		if _, exists := byID[k]; !exists {
			order = append(order, k)
			byID[k] = &circulation.EconomistAttribution{
				ID:          k,
				Concept:     concept,
				Theorist:    thr,
				EditionYear: editionYear,
				Anticipates: anticipates,
				Errors:      []circulation.KnownEconomistError{},
			}
		}
		if errorCode != nil {
			byID[k].Errors = append(byID[k].Errors, circulation.KnownEconomistError(*errorCode))
		}
	}
	if closeErr := rows.Close(); closeErr != nil {
		return nil, closeErr
	}

	out := make([]circulation.EconomistAttribution, 0, len(order))
	for _, k := range order {
		out = append(out, *byID[k])
	}
	return out, nil
}
