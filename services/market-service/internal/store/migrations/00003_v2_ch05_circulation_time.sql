-- +goose Up

CREATE TABLE turnover_times (
    id                        VARCHAR(24)  NOT NULL PRIMARY KEY,
    industrial_capital_id     VARCHAR(24)  NOT NULL,
    labour_time_nanos         BIGINT       NOT NULL DEFAULT 0,
    labour_interruption_nanos BIGINT       NOT NULL DEFAULT 0,
    latent_nanos              BIGINT       NOT NULL DEFAULT 0,
    natural_process_nanos     BIGINT       NOT NULL DEFAULT 0,
    selling_time_nanos        BIGINT       NOT NULL DEFAULT 0,
    buying_time_nanos         BIGINT       NOT NULL DEFAULT 0,
    circulation_open          TINYINT(1)   NOT NULL DEFAULT 0,
    created_at                DATETIME(6)  NOT NULL,
    updated_at                DATETIME(6)  NOT NULL,
    INDEX idx_tt_industrial_capital (industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE selling_phases (
    id                    VARCHAR(24)  NOT NULL PRIMARY KEY,
    turnover_time_id      VARCHAR(24)  NOT NULL,
    industrial_capital_id VARCHAR(24)  NOT NULL,
    commodity_circuit_id  VARCHAR(24)  NOT NULL DEFAULT '',
    opened_at             DATETIME(6)  NOT NULL,
    closed_at             DATETIME(6)  NULL,
    outcome               VARCHAR(16)  NOT NULL DEFAULT 'pending',
    INDEX idx_sp_turnover (turnover_time_id),
    INDEX idx_sp_open (turnover_time_id, closed_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE buying_phases (
    id                    VARCHAR(24)  NOT NULL PRIMARY KEY,
    turnover_time_id      VARCHAR(24)  NOT NULL,
    industrial_capital_id VARCHAR(24)  NOT NULL,
    money_circuit_id      VARCHAR(24)  NOT NULL DEFAULT '',
    opened_at             DATETIME(6)  NOT NULL,
    closed_at             DATETIME(6)  NULL,
    market_location       VARCHAR(128) NOT NULL DEFAULT '',
    INDEX idx_bp_turnover (turnover_time_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE natural_process_spans (
    id                    VARCHAR(24)  NOT NULL PRIMARY KEY,
    turnover_time_id      VARCHAR(24)  NOT NULL,
    industrial_capital_id VARCHAR(24)  NOT NULL,
    process               VARCHAR(24)  NOT NULL,
    duration_nanos        BIGINT       NOT NULL,
    started_at            DATETIME(6)  NOT NULL,
    INDEX idx_nps_turnover (turnover_time_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE latent_productive_capital (
    id                      VARCHAR(24)  NOT NULL PRIMARY KEY,
    turnover_time_id        VARCHAR(24)  NOT NULL,
    industrial_capital_id   VARCHAR(24)  NOT NULL,
    pence                   BIGINT       NOT NULL,
    held_at                 DATETIME(6)  NOT NULL,
    entered_production_at   DATETIME(6)  NULL,
    INDEX idx_lpc_turnover (turnover_time_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE perishability (
    id           VARCHAR(24)  NOT NULL PRIMARY KEY,
    commodity_id VARCHAR(24)  NOT NULL,
    window_nanos BIGINT       NOT NULL,
    UNIQUE KEY uq_perishability_commodity (commodity_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE market_separation (
    id                    VARCHAR(24)  NOT NULL PRIMARY KEY,
    industrial_capital_id VARCHAR(24)  NOT NULL,
    selling_market_id     VARCHAR(64)  NOT NULL,
    buying_market_id      VARCHAR(64)  NOT NULL,
    UNIQUE KEY uq_ms_industrial_capital (industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down

DROP TABLE IF EXISTS market_separation;
DROP TABLE IF EXISTS perishability;
DROP TABLE IF EXISTS latent_productive_capital;
DROP TABLE IF EXISTS natural_process_spans;
DROP TABLE IF EXISTS buying_phases;
DROP TABLE IF EXISTS selling_phases;
DROP TABLE IF EXISTS turnover_times;
