-- +goose Up
-- Vol. III Ch. 4 — turnover analyses. Each row records one capital and the
-- effect of its turnover on the rate of profit: the total capital C, the
-- variable capital advanced v (the v in C), the rate of surplus-value s' (basis
-- points, 100% = 10000 bp) and the number of turnovers n. The derived columns
-- hold the annual rate of profit p' = s'·n·(v/C), the single-turnover rate
-- s'·(v/C) it would be without turnover, the annual rate of surplus-value
-- S' = s'·n, and the year's whole wage bill vn (the figure the capitalist's
-- books record, as against the variable capital v it conceals). Capital
-- magnitudes are LabourMinutes (BIGINT); the *_bp columns are dimensionless.

CREATE TABLE IF NOT EXISTS turnover_analyses (
    id                           VARCHAR(24) NOT NULL PRIMARY KEY,
    total_capital                BIGINT      NOT NULL,
    variable_capital             BIGINT      NOT NULL,
    surplus_value_rate_bp        BIGINT      NOT NULL,
    turnovers                    BIGINT      NOT NULL,
    annual_profit_rate_bp        BIGINT      NOT NULL,
    single_turnover_rate_bp      BIGINT      NOT NULL,
    annual_surplus_value_rate_bp BIGINT      NOT NULL,
    annual_wages                 BIGINT      NOT NULL,
    created_at                   DATETIME(6) NOT NULL,
    INDEX idx_turnover_analyses_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS turnover_analyses;
