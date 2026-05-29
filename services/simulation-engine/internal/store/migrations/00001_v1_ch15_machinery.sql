-- +goose Up
-- Ch. 15 — Machinery and Modern Industry.
--
-- First simulation-engine persistence: a Machine is Marx's three-part
-- instrument (motor mechanism, transmitting mechanism, working tool); a
-- Factory aggregates Machine records under a shared PrimeMover. Each
-- factory_ticks row is one period's worth of value transfer and use-value
-- output written by POST /v1/factories/{id}/tick.

CREATE TABLE IF NOT EXISTS machines (
    id                       VARCHAR(24)  NOT NULL PRIMARY KEY,
    name                     VARCHAR(255) NOT NULL DEFAULT '',
    motor_mechanism          VARCHAR(255) NOT NULL DEFAULT '',
    transmitting_mechanism   VARCHAR(255) NOT NULL DEFAULT '',
    working_tool             VARCHAR(255) NOT NULL DEFAULT '',
    machine_value            BIGINT       NOT NULL DEFAULT 0,
    lifespan_days            BIGINT       NOT NULL DEFAULT 0,
    productive_power         BIGINT       NOT NULL DEFAULT 0,
    hand_labour_per_unit     BIGINT       NOT NULL DEFAULT 0,
    accumulated_wear         BIGINT       NOT NULL DEFAULT 0,
    accumulated_depreciation BIGINT       NOT NULL DEFAULT 0,
    created_at               DATETIME(6)  NOT NULL,
    INDEX idx_machines_name (name)
);

CREATE TABLE IF NOT EXISTS factories (
    id                     VARCHAR(24)  NOT NULL PRIMARY KEY,
    name                   VARCHAR(255) NOT NULL DEFAULT '',
    prime_mover_kind       VARCHAR(32)  NOT NULL,
    prime_mover_horsepower BIGINT       NOT NULL DEFAULT 0,
    tick_count             BIGINT       NOT NULL DEFAULT 0,
    created_at             DATETIME(6)  NOT NULL,
    INDEX idx_factories_name (name)
);

CREATE TABLE IF NOT EXISTS factory_machines (
    factory_id VARCHAR(24) NOT NULL,
    machine_id VARCHAR(24) NOT NULL,
    ordinal    INT         NOT NULL,
    PRIMARY KEY (factory_id, machine_id),
    INDEX idx_factory_machines_factory (factory_id),
    INDEX idx_factory_machines_machine (machine_id)
);

CREATE TABLE IF NOT EXISTS factory_ticks (
    factory_id        VARCHAR(24) NOT NULL,
    sequence          BIGINT      NOT NULL,
    value_transferred BIGINT      NOT NULL DEFAULT 0,
    units_produced    BIGINT      NOT NULL DEFAULT 0,
    hand_labour_saved BIGINT      NOT NULL DEFAULT 0,
    occurred_at       DATETIME(6) NOT NULL,
    PRIMARY KEY (factory_id, sequence),
    INDEX idx_factory_ticks_factory (factory_id)
);

-- +goose Down
DROP TABLE IF EXISTS factory_ticks;
DROP TABLE IF EXISTS factory_machines;
DROP TABLE IF EXISTS factories;
DROP TABLE IF EXISTS machines;
