-- +goose Up
-- Vol. III Ch. 10 — Equalisation of the General Rate of Profit Through Competition.
-- Three tables: market_values, surplus_profits, equalisations.

CREATE TABLE market_values (
    id                   VARCHAR(24)  NOT NULL,
    sphere_name          VARCHAR(120) NOT NULL,
    bulk_condition_value BIGINT       NOT NULL,
    best_condition_value BIGINT       NOT NULL,
    worst_condition_value BIGINT      NOT NULL,
    value                BIGINT       NOT NULL,
    created_at           DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_market_values_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE surplus_profits (
    id               VARCHAR(24)  NOT NULL,
    firm_name        VARCHAR(120) NOT NULL,
    individual_value BIGINT       NOT NULL,
    market_value     BIGINT       NOT NULL,
    output_qty       BIGINT       NOT NULL,
    general_rate     BIGINT       NOT NULL,
    amount           BIGINT       NOT NULL, -- signed: positive = extra profit, negative = shortfall
    created_at       DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_surplus_profits_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE equalisations (
    id            VARCHAR(24)  NOT NULL,
    sphere_name   VARCHAR(120) NOT NULL,
    initial_rate  BIGINT       NOT NULL,
    target_rate   BIGINT       NOT NULL,
    direction     VARCHAR(16)  NOT NULL, -- 'inflow' | 'outflow' | 'stable'
    is_converging TINYINT(1)   NOT NULL,
    market_price  BIGINT       NOT NULL, -- Flow.MarketPrice; 0 when rate-based
    market_value  BIGINT       NOT NULL, -- Flow.MarketValue; 0 when rate-based
    created_at    DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_equalisations_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS equalisations;
DROP TABLE IF EXISTS surplus_profits;
DROP TABLE IF EXISTS market_values;
