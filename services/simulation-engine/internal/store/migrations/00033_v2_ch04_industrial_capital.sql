-- +goose Up
-- Vol. II Ch. 4 — The Three Formulas of the Circuit
-- Schema for IndustrialCapital and its satellite tables.

CREATE TABLE industrial_capitals (
    id                         CHAR(24)     NOT NULL,
    agent_id                   VARCHAR(255) NOT NULL DEFAULT '',
    money_circuit_id           CHAR(24)     NOT NULL DEFAULT '',
    productive_circuit_id      CHAR(24)     NOT NULL DEFAULT '',
    commodity_circuit_id       CHAR(24)     NOT NULL DEFAULT '',
    total_pence                BIGINT       NOT NULL,
    economy_mode               VARCHAR(16)  NOT NULL DEFAULT 'money',
    stagnation_tolerance_ticks BIGINT       NOT NULL DEFAULT 3,
    status                     VARCHAR(16)  NOT NULL DEFAULT 'active',
    created_at                 DATETIME(6)  NOT NULL,
    updated_at                 DATETIME(6)  NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE capital_parts (
    id                    CHAR(24)    NOT NULL,
    industrial_capital_id CHAR(24)    NOT NULL,
    pence                 BIGINT      NOT NULL,
    stage                 VARCHAR(16) NOT NULL,
    entered_stage_at      DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_cp_capital (industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE stage_distributions (
    id                    CHAR(24)    NOT NULL,
    industrial_capital_id CHAR(24)    NOT NULL,
    at_time               DATETIME(6) NOT NULL,
    money_pence           BIGINT      NOT NULL,
    production_pence      BIGINT      NOT NULL,
    commodity_pence       BIGINT      NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_sd_capital_time (industrial_capital_id, at_time DESC)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE stage_blocks (
    id                    CHAR(24)    NOT NULL,
    industrial_capital_id CHAR(24)    NOT NULL,
    stage                 VARCHAR(16) NOT NULL,
    reason                VARCHAR(64) NOT NULL,
    opened_at             DATETIME(6) NOT NULL,
    closed_at             DATETIME(6) NULL,
    PRIMARY KEY (id),
    INDEX idx_sb_capital_open (industrial_capital_id, closed_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE value_revolution_events (
    id                    CHAR(24)    NOT NULL,
    industrial_capital_id CHAR(24)    NOT NULL,
    affecting             VARCHAR(32) NOT NULL,
    direction_pence       BIGINT      NOT NULL,
    occurred_at           DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_vre_capital (industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE money_capital_set_free (
    id                        CHAR(24) NOT NULL,
    value_revolution_event_id CHAR(24) NOT NULL,
    industrial_capital_id     CHAR(24) NOT NULL,
    amount_pence              BIGINT   NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_mcsf_event (value_revolution_event_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE money_capital_tied_up (
    id                        CHAR(24) NOT NULL,
    value_revolution_event_id CHAR(24) NOT NULL,
    industrial_capital_id     CHAR(24) NOT NULL,
    amount_pence              BIGINT   NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_mctu_event (value_revolution_event_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE metamorphosis_interlocks (
    id                           CHAR(24)    NOT NULL,
    buyer_industrial_capital_id  CHAR(24)    NOT NULL,
    seller_industrial_capital_id CHAR(24)    NOT NULL,
    seller_mp_source_origin      VARCHAR(32) NOT NULL,
    pence                        BIGINT      NOT NULL,
    occurred_at                  DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_mi_buyer (buyer_industrial_capital_id),
    INDEX idx_mi_seller (seller_industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE supply_demand_imbalances (
    id                    CHAR(24)    NOT NULL,
    industrial_capital_id CHAR(24)    NOT NULL,
    period                VARCHAR(32) NOT NULL,
    demand_pence          BIGINT      NOT NULL,
    supply_pence          BIGINT      NOT NULL,
    excess_pence          BIGINT      NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_sdi_capital_period (industrial_capital_id, period)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sinking_funds (
    id                    CHAR(24) NOT NULL,
    industrial_capital_id CHAR(24) NOT NULL UNIQUE,
    fixed_capital_pence   BIGINT   NOT NULL,
    lifetime_years        BIGINT   NOT NULL,
    annual_payment_pence  BIGINT   NOT NULL,
    accumulated_pence     BIGINT   NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    INDEX idx_sf_capital (industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS sinking_funds;
DROP TABLE IF EXISTS supply_demand_imbalances;
DROP TABLE IF EXISTS metamorphosis_interlocks;
DROP TABLE IF EXISTS money_capital_tied_up;
DROP TABLE IF EXISTS money_capital_set_free;
DROP TABLE IF EXISTS value_revolution_events;
DROP TABLE IF EXISTS stage_blocks;
DROP TABLE IF EXISTS stage_distributions;
DROP TABLE IF EXISTS capital_parts;
DROP TABLE IF EXISTS industrial_capitals;
