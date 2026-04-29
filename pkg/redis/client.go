// Package redis is the scaffolding equivalent of pkg/mongo: it centralizes how
// services read Redis connection settings from env. The actual go-redis client
// will be wired in once a chapter introduces caching needs.
package redis

import (
	"errors"
	"os"
	"strconv"
	"time"
)

// Config holds connection settings for a Redis client.
type Config struct {
	Addr           string        // e.g. redis:6379
	Password       string        // empty if not required
	DB             int           // logical database number
	DialTimeout    time.Duration // default 5s
	ReadTimeout    time.Duration // default 3s
	WriteTimeout   time.Duration // default 3s
}

// ConfigFromEnv reads REDIS_ADDR, REDIS_PASSWORD, REDIS_DB.
func ConfigFromEnv() Config {
	db, _ := strconv.Atoi(os.Getenv("REDIS_DB"))
	return Config{
		Addr:         getenvDefault("REDIS_ADDR", "redis:6379"),
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           db,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// Validate ensures the config has the minimum required values to dial.
func (c Config) Validate() error {
	if c.Addr == "" {
		return errors.New("redis: Addr is required")
	}
	return nil
}

func getenvDefault(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
