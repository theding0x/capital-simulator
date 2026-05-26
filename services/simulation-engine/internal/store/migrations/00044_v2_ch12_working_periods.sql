-- +goose Up
-- Vol. II Ch. 12 — The Working Period

CREATE TABLE working_periods (
    id                         VARCHAR(24)  NOT NULL,
    industrial_capital_id      VARCHAR(255) NOT NULL,
    commodity_id               VARCHAR(255) NOT NULL,
    working_day_count          BIGINT       NOT NULL,
    mode                       ENUM('discrete','connected') NOT NULL,
    working_day_length_minutes BIGINT       NOT NULL,
    created_at                 DATETIME(6)  NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE working_period_interruptions (
    id                                 VARCHAR(24) NOT NULL,
    working_period_id                  VARCHAR(24) NOT NULL,
    interrupted_at                     DATETIME(6) NOT NULL,
    resumed_at                         DATETIME(6) NULL,
    deterioration_pence                BIGINT      NOT NULL DEFAULT 0,
    waste_of_means_of_production_pence BIGINT      NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    KEY idx_wpi_period (working_period_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE working_period_shortenings (
    id                             VARCHAR(24) NOT NULL,
    working_period_id              VARCHAR(24) NOT NULL,
    kind                           ENUM('cooperation','machinery','parallel_labour','bakewell_breeding') NOT NULL,
    old_working_day_count          BIGINT      NOT NULL,
    new_working_day_count          BIGINT      NOT NULL,
    additional_fixed_capital_pence BIGINT      NOT NULL DEFAULT 0,
    additional_labour_count        BIGINT      NOT NULL DEFAULT 0,
    occurred_at                    DATETIME(6) NOT NULL,
    PRIMARY KEY (id),
    KEY idx_wps_period (working_period_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE natural_working_period_constraints (
    id               VARCHAR(24)  NOT NULL,
    commodity_id     VARCHAR(255) NOT NULL,
    min_working_days BIGINT       NOT NULL,
    reason           TEXT         NOT NULL,
    created_at       DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_nwpc_commodity (commodity_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE working_period_credit_financings (
    id                     VARCHAR(24)  NOT NULL,
    working_period_id      VARCHAR(24)  NOT NULL,
    own_capital_pence      BIGINT       NOT NULL DEFAULT 0,
    borrowed_capital_pence BIGINT       NOT NULL DEFAULT 0,
    mortgage_provider_id   VARCHAR(255) NOT NULL DEFAULT '',
    created_at             DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_wpcf_period (working_period_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS working_period_credit_financings;
DROP TABLE IF EXISTS natural_working_period_constraints;
DROP TABLE IF EXISTS working_period_shortenings;
DROP TABLE IF EXISTS working_period_interruptions;
DROP TABLE IF EXISTS working_periods;
