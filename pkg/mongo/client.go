// Package mongo is a thin scaffolding wrapper around MongoDB connection
// configuration. The real *mongo.Client is intentionally not wired up yet -
// it will be introduced in a chapter that requires persistence (likely the
// commodity service in Capital ch. 1). This package's job today is to give
// every service a single canonical place to read mongo config from env.
package mongo

import (
	"errors"
	"os"
	"time"
)

// Config holds connection settings for a MongoDB client.
type Config struct {
	URI            string        // e.g. mongodb://mongo:27017
	Database       string        // logical database name for the service
	ConnectTimeout time.Duration // default 10s
}

// ConfigFromEnv reads MONGO_URI and MONGO_DATABASE from the environment.
// `service` is used as a fallback database name when MONGO_DATABASE is unset.
func ConfigFromEnv(service string) Config {
	db := os.Getenv("MONGO_DATABASE")
	if db == "" {
		db = service
	}
	return Config{
		URI:            getenvDefault("MONGO_URI", "mongodb://mongo:27017"),
		Database:       db,
		ConnectTimeout: 10 * time.Second,
	}
}

// Validate ensures the config has the minimum required values to dial.
func (c Config) Validate() error {
	if c.URI == "" {
		return errors.New("mongo: URI is required")
	}
	if c.Database == "" {
		return errors.New("mongo: Database is required")
	}
	return nil
}

func getenvDefault(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
