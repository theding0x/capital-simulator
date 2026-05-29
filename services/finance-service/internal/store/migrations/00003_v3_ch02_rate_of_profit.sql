-- +goose Up
-- Vol. III Ch. 2 — rate-of-profit analyses. Each row records one capital
-- decomposition c + v + s and the two rates it gives rise to: the rate of
-- profit p' = s/C = s/(c+v) and the rate of surplus-value s' = s/v, both stored
-- in basis points (100% = 10000 bp), together with the mystification degree
-- (s' - p')/s' in basis points. Capital magnitudes are LabourMinutes (BIGINT);
-- the *_bp columns are dimensionless ratios.

CREATE TABLE IF NOT EXISTS profit_rates (
    id                    VARCHAR(24) NOT NULL PRIMARY KEY,
    constant              BIGINT      NOT NULL,
    variable              BIGINT      NOT NULL,
    surplus_value         BIGINT      NOT NULL,
    total_capital         BIGINT      NOT NULL,
    profit_rate_bp        BIGINT      NOT NULL,
    surplus_value_rate_bp BIGINT      NOT NULL,
    mystification_bp      BIGINT      NOT NULL,
    created_at            DATETIME(6) NOT NULL,
    INDEX idx_profit_rates_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS profit_rates;
