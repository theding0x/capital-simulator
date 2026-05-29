-- +goose Up
-- Issue #214 — automatic tick scheduler.
--
-- The simulation-engine scheduler advances the economy one period at a time on a
-- fixed interval (SIM_TICK_INTERVAL, default 5s), invoking each registered ticker
-- (factory, reproduction) in order. Each engine_ticks row is one scheduler pass:
-- how many tickers ran, how many domain entities they advanced, how long the pass
-- took, and how many tickers errored. The table is an append-only audit log
-- written by the scheduler at runtime; it is not seeded.

CREATE TABLE IF NOT EXISTS engine_ticks (
    id                VARCHAR(24)  NOT NULL PRIMARY KEY,
    sequence          BIGINT       NOT NULL,
    occurred_at       DATETIME(6)  NOT NULL,
    duration_ms       BIGINT       NOT NULL DEFAULT 0,
    tickers_run       INT          NOT NULL DEFAULT 0,
    entities_advanced INT          NOT NULL DEFAULT 0,
    error_count       INT          NOT NULL DEFAULT 0,
    INDEX idx_engine_ticks_sequence (sequence),
    INDEX idx_engine_ticks_occurred (occurred_at)
);

-- +goose Down
DROP TABLE IF EXISTS engine_ticks;
