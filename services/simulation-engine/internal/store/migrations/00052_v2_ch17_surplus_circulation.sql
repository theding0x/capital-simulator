-- +goose Up
-- Vol. II Ch. 17 — The Circulation of Surplus-Value
-- Schema DDL for the reproduction package tables.

-- Primary table: one row per period's aggregate surplus circulation.
CREATE TABLE surplus_circulations (
    id                       CHAR(24)    NOT NULL,
    period                   VARCHAR(64) NOT NULL,
    total_surplus_pence      BIGINT      NOT NULL,
    realised_surplus_pence   BIGINT      NOT NULL DEFAULT 0,
    unrealised_surplus_pence BIGINT      NOT NULL,
    created_at               DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_sc_period (period)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- One row per realisation-source entry; multiple rows per circulation.
-- Auto-increment id provides insertion-order ordering.
CREATE TABLE realisation_source_entries (
    id                         BIGINT      NOT NULL AUTO_INCREMENT,
    surplus_circulation_id     CHAR(24)    NOT NULL,
    source                     VARCHAR(64) NOT NULL,
    pence                      BIGINT      NOT NULL,
    deferred_realisation_pence BIGINT      NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    INDEX idx_rse_circulation (surplus_circulation_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Social-capital aggregate snapshots, one per period.
-- department_i/ii_share_bp are placeholders until Chs. 20-21.
CREATE TABLE social_capital_aggregates (
    period                     VARCHAR(64) NOT NULL,
    total_advanced_pence       BIGINT      NOT NULL DEFAULT 0,
    total_annual_output_pence  BIGINT      NOT NULL DEFAULT 0,
    total_annual_surplus_pence BIGINT      NOT NULL DEFAULT 0,
    department_i_share_bp      BIGINT      NOT NULL DEFAULT 0,
    department_ii_share_bp     BIGINT      NOT NULL DEFAULT 0,
    PRIMARY KEY (period)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Canonical realisation-puzzle reference rows.
CREATE TABLE realisation_puzzles (
    surplus_circulation_id CHAR(24) NOT NULL,
    puzzle_statement       TEXT     NOT NULL,
    resolved_in_chapter    BIGINT   NOT NULL,
    PRIMARY KEY (surplus_circulation_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS realisation_puzzles;
DROP TABLE IF EXISTS social_capital_aggregates;
DROP TABLE IF EXISTS realisation_source_entries;
DROP TABLE IF EXISTS surplus_circulations;
