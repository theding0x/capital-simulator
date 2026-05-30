-- +goose Up
-- Vol. III Ch. 9 — Formation of a General Rate of Profit (Average Rate of
-- Profit) and Transformation of the Values of Commodities into Prices of
-- Production.
--
-- general_profit_rates: one row per computed general (average) rate. The
-- Spheres slice is stored as a JSON column because it is variable-length
-- and read-only after creation.
--
-- prices_of_production: one row per sphere's price-of-production record.
-- Flat columns for easy querying; deviation is signed (BIGINT, not UNSIGNED).
--
-- Rates are in basis points (10000 = 100%); capital magnitudes are
-- LabourMinutes (BIGINT). ID column width is 24 chars (crypto/rand 96-bit hex).

CREATE TABLE IF NOT EXISTS general_profit_rates (
    id                       VARCHAR(24) NOT NULL PRIMARY KEY,
    rate                     BIGINT      NOT NULL,  -- p̄′ in bp
    sum_surplus_values       BIGINT      NOT NULL,
    sum_total_capitals       BIGINT      NOT NULL,
    average_variable_percent BIGINT      NOT NULL,
    spheres_json             JSON        NOT NULL,
    created_at               DATETIME(6) NOT NULL,
    INDEX idx_general_profit_rates_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS prices_of_production (
    id              VARCHAR(24)  NOT NULL PRIMARY KEY,
    sphere_name     VARCHAR(120) NOT NULL,
    cost_price      BIGINT       NOT NULL,
    general_rate    BIGINT       NOT NULL,
    commodity_value BIGINT       NOT NULL,
    average_profit  BIGINT       NOT NULL,
    price           BIGINT       NOT NULL,
    deviation       BIGINT       NOT NULL,  -- signed: price − value
    composition     VARCHAR(16)  NOT NULL,
    created_at      DATETIME(6)  NOT NULL,
    INDEX idx_prices_of_production_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS prices_of_production;
DROP TABLE IF EXISTS general_profit_rates;
