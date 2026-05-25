-- +goose Up
-- Vol. II Ch. 8 — Fixed Capital and Circulating Capital
-- Instruments of labour (fixed) vs raw/auxiliary/labour-power (circulating).

CREATE TABLE capital_components (
    id                    VARCHAR(24)  NOT NULL PRIMARY KEY,
    industrial_capital_id VARCHAR(24)  NOT NULL DEFAULT '',
    kind                  VARCHAR(16)  NOT NULL COMMENT 'fixed|circulating',
    role                  VARCHAR(32)  NOT NULL COMMENT 'instrument_of_labour|raw_material|auxiliary_material|labour_power',
    pence                 BIGINT       NOT NULL DEFAULT 0,
    entered_at            DATETIME(3)  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE fixed_capital_items (
    id                    VARCHAR(24)  NOT NULL PRIMARY KEY,
    capital_component_id  VARCHAR(24)  NOT NULL DEFAULT '',
    description           VARCHAR(255) NOT NULL DEFAULT '',
    pence_purchased       BIGINT       NOT NULL DEFAULT 0,
    service_life_minutes  BIGINT       NOT NULL,
    wear_model            VARCHAR(16)  NOT NULL DEFAULT 'linear' COMMENT 'linear|quadratic',
    purchased_at          DATETIME(3)  NOT NULL,
    retired_at            DATETIME(3)  NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE fixed_capital_subcomponents (
    id                    VARCHAR(24)  NOT NULL PRIMARY KEY,
    fixed_capital_item_id VARCHAR(24)  NOT NULL,
    name                  VARCHAR(255) NOT NULL DEFAULT '',
    pence                 BIGINT       NOT NULL DEFAULT 0,
    service_life_minutes  BIGINT       NOT NULL,
    wear_model            VARCHAR(16)  NOT NULL DEFAULT 'linear',
    last_replaced_at      DATETIME(3)  NOT NULL,
    INDEX idx_fcs_item (fixed_capital_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE wear_and_tear (
    id                     VARCHAR(24) NOT NULL PRIMARY KEY,
    fixed_capital_item_id  VARCHAR(24) NOT NULL,
    pence                  BIGINT      NOT NULL DEFAULT 0,
    period_covered_minutes BIGINT      NOT NULL DEFAULT 0,
    calculated_at          DATETIME(3) NOT NULL,
    cause                  VARCHAR(32) NOT NULL COMMENT 'use|atmospheric|moral_depreciation',
    INDEX idx_wat_item (fixed_capital_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sinking_funds_for_items (
    id                    VARCHAR(24) NOT NULL PRIMARY KEY,
    fixed_capital_item_id VARCHAR(24) NOT NULL UNIQUE,
    accumulated_pence     BIGINT      NOT NULL DEFAULT 0,
    target_pence          BIGINT      NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE repairs (
    id                             VARCHAR(24)  NOT NULL PRIMARY KEY,
    fixed_capital_item_id          VARCHAR(24)  NOT NULL,
    kind                           VARCHAR(16)  NOT NULL COMMENT 'running|large',
    pence                          BIGINT       NOT NULL DEFAULT 0,
    occurred_at                    DATETIME(3)  NOT NULL,
    service_life_extension_minutes BIGINT       NOT NULL DEFAULT 0,
    INDEX idx_repairs_item (fixed_capital_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE circulating_cycles (
    id                   VARCHAR(24) NOT NULL PRIMARY KEY,
    capital_component_id VARCHAR(24) NOT NULL,
    started_at           DATETIME(3) NOT NULL,
    ended_at             DATETIME(3) NOT NULL,
    pence                BIGINT      NOT NULL DEFAULT 0,
    INDEX idx_cc_component (capital_component_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE sinking_fund_reinvestments (
    id                    VARCHAR(24) NOT NULL PRIMARY KEY,
    fixed_capital_item_id VARCHAR(24) NOT NULL,
    pence                 BIGINT      NOT NULL DEFAULT 0,
    kind                  VARCHAR(16) NOT NULL COMMENT 'intensive|extensive',
    occurred_at           DATETIME(3) NOT NULL,
    INDEX idx_sfr_item (fixed_capital_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS sinking_fund_reinvestments;
DROP TABLE IF EXISTS circulating_cycles;
DROP TABLE IF EXISTS repairs;
DROP TABLE IF EXISTS sinking_funds_for_items;
DROP TABLE IF EXISTS wear_and_tear;
DROP TABLE IF EXISTS fixed_capital_subcomponents;
DROP TABLE IF EXISTS fixed_capital_items;
DROP TABLE IF EXISTS capital_components;
