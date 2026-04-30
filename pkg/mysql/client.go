// Package mysql wraps database/sql with the conventions used across
// capital-simulator: env-driven DSN, ping-on-connect, and a single *sql.DB
// handle that callers receive.
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Config holds connection settings for a MySQL client.
type Config struct {
	DSN            string        // e.g. root:capital@tcp(mysql:3306)/commodity?parseTime=true&loc=UTC
	ConnectTimeout time.Duration // default 10s
	AppName        string        // informational only
}

// ConfigFromEnv reads MYSQL_DSN from the environment. If unset, it builds a
// default DSN pointing at mysql:3306 using the service name as the schema.
func ConfigFromEnv(service string) Config {
	dsn := getenvDefault("MYSQL_DSN",
		fmt.Sprintf("root:capital@tcp(mysql:3306)/%s?parseTime=true&loc=UTC", service))
	return Config{
		DSN:            dsn,
		ConnectTimeout: 10 * time.Second,
		AppName:        service,
	}
}

// Validate ensures the config has the minimum required values.
func (c Config) Validate() error {
	if c.DSN == "" {
		return errors.New("mysql: DSN is required")
	}
	return nil
}

// DB wraps a connected *sql.DB.
type DB struct {
	SQL *sql.DB
}

// Connect opens a *sql.DB using cfg.DSN and pings the server to verify the
// connection is live. Callers should defer Close.
func Connect(ctx context.Context, cfg Config) (*DB, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = 10 * time.Second
	}

	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("mysql: open: %w", err)
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	pingCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("mysql: ping: %w", err)
	}
	return &DB{SQL: db}, nil
}

// Close closes the underlying *sql.DB. Safe to call on nil.
func (d *DB) Close() error {
	if d == nil || d.SQL == nil {
		return nil
	}
	return d.SQL.Close()
}

func getenvDefault(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
