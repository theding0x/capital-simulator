package store

import (
	"context"
	"database/sql"
	"time"
)

// MySQL is a MySQL-backed Store for finance-service. Construct via NewMySQL.
//
// finance-service has no migrations yet — the migrations/ directory is
// empty and goose is not invoked here. When the first Vol. III chapter PR
// lands, it will:
//   1. Add a `//go:embed migrations` directive at the top of this file.
//   2. Import `embed`, `io/fs`, and `pkgmysql` (alias for pkg/mysql).
//   3. Call `pkgmysql.Migrate(ctx, db, sub)` inside NewMySQL.
//   4. Drop the underscore prefix from `_ context.Context` once ctx is used.
type MySQL struct {
	db  *sql.DB
	now func() time.Time
}

// NewMySQL returns a Store backed by db. No migrations run because none
// exist yet — see the type doc above for the per-chapter onboarding.
func NewMySQL(_ context.Context, db *sql.DB) (*MySQL, error) {
	return &MySQL{db: db, now: time.Now}, nil
}
