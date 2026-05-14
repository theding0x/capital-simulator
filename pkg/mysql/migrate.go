package mysql

import (
	"context"
	"database/sql"
	"io/fs"

	"github.com/pressly/goose/v3"
)

// Migrate applies all pending SQL migrations from fsys.
// fsys must contain .sql files at its root — use fs.Sub to strip a subdirectory prefix.
func Migrate(ctx context.Context, db *sql.DB, fsys fs.FS) error {
	p, err := goose.NewProvider(goose.DialectMySQL, db, fsys)
	if err != nil {
		return err
	}
	_, err = p.Up(ctx)
	return err
}
