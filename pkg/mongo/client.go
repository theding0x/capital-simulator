// Package mongo wraps the official MongoDB driver with the conventions used
// across capital-simulator: env-driven config, ping-on-connect, and a single
// resolved Database handle that callers receive together with the *mongo.Client.
package mongo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config holds connection settings for a MongoDB client.
type Config struct {
	URI            string        // e.g. mongodb://mongo:27017
	Database       string        // logical database name for the service
	ConnectTimeout time.Duration // default 10s
	AppName        string        // shows up in mongo logs / currentOp
}

// ConfigFromEnv reads MONGO_URI and MONGO_DATABASE from the environment.
// `service` is used as the AppName and as the fallback database name when
// MONGO_DATABASE is unset.
func ConfigFromEnv(service string) Config {
	db := os.Getenv("MONGO_DATABASE")
	if db == "" {
		db = service
	}
	return Config{
		URI:            getenvDefault("MONGO_URI", "mongodb://mongo:27017"),
		Database:       db,
		ConnectTimeout: 10 * time.Second,
		AppName:        service,
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

// Client bundles a connected *mongo.Client and the resolved database handle.
type Client struct {
	Mongo *mongo.Client
	DB    *mongo.Database
}

// Connect dials MongoDB using cfg and pings the primary to verify the
// connection is live before returning. Callers should defer Close.
func Connect(ctx context.Context, cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = 10 * time.Second
	}

	dialCtx, cancel := context.WithTimeout(ctx, cfg.ConnectTimeout)
	defer cancel()

	opts := options.Client().ApplyURI(cfg.URI).SetAppName(cfg.AppName)
	cli, err := mongo.Connect(dialCtx, opts)
	if err != nil {
		return nil, fmt.Errorf("mongo: connect: %w", err)
	}

	if err := cli.Ping(dialCtx, readpref.Primary()); err != nil {
		_ = cli.Disconnect(context.Background())
		return nil, fmt.Errorf("mongo: ping: %w", err)
	}

	return &Client{
		Mongo: cli,
		DB:    cli.Database(cfg.Database),
	}, nil
}

// Close disconnects from MongoDB. It is safe to call on a nil receiver.
func (c *Client) Close(ctx context.Context) error {
	if c == nil || c.Mongo == nil {
		return nil
	}
	return c.Mongo.Disconnect(ctx)
}

func getenvDefault(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
