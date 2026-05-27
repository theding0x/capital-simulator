-- +goose Up
-- Vol. II Ch. 13 — The Time of Production

-- production_labour_gaps: the production-time / labour-time split for a working period.
-- One row per working period. gap_nanos = production_time_nanos - labour_time_nanos (≥ 0).
CREATE TABLE production_labour_gaps (
    id                    VARCHAR(24) NOT NULL,
    working_period_id     VARCHAR(24) NOT NULL,
    production_time_nanos BIGINT      NOT NULL,
    labour_time_nanos     BIGINT      NOT NULL,
    gap_nanos             BIGINT      NOT NULL,
    created_at            DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_plg_period (working_period_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- natural_process_activations: a span during which a natural process runs.
-- Multiple activations per working period are allowed.
CREATE TABLE natural_process_activations (
    id                          VARCHAR(24) NOT NULL,
    working_period_id           VARCHAR(24) NOT NULL,
    kind                        VARCHAR(64) NOT NULL,
    started_at                  DATETIME(6) NOT NULL,
    ended_at                    DATETIME(6) NULL,
    requires_labour_maintenance TINYINT(1)  NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    KEY idx_npa_period (working_period_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- natural_subjects: the commodity or material undergoing a natural process.
-- One row per working period.
CREATE TABLE natural_subjects (
    id                VARCHAR(24)  NOT NULL,
    working_period_id VARCHAR(24)  NOT NULL,
    commodity_id      VARCHAR(255) NOT NULL,
    description       TEXT         NOT NULL,
    is_living         TINYINT(1)   NOT NULL DEFAULT 0,
    created_at        DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_ns_period (working_period_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- natural_process_industries: descriptive benchmarks for industries dominated
-- by natural processes. Not linked to specific working periods.
CREATE TABLE natural_process_industries (
    id                            VARCHAR(24)  NOT NULL,
    industry_name                 VARCHAR(255) NOT NULL,
    typical_production_days_count BIGINT       NOT NULL,
    typical_labour_days_count     BIGINT       NOT NULL,
    ratio_basis_points            BIGINT       NOT NULL,
    created_at                    DATETIME(6)  NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS natural_process_industries;
DROP TABLE IF EXISTS natural_subjects;
DROP TABLE IF EXISTS natural_process_activations;
DROP TABLE IF EXISTS production_labour_gaps;
