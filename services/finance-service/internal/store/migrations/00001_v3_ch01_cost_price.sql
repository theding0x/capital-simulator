-- +goose Up
-- Vol. III Ch. 1 — cost-price records. A cost-price is k = c + v, split into
-- the fixed component (depreciation of instruments of labour entering the
-- cost-price) and the circulating component (circulating constant capital plus
-- variable capital). All magnitudes are LabourMinutes (BIGINT).

CREATE TABLE IF NOT EXISTS cost_prices (
    id                    VARCHAR(24) NOT NULL PRIMARY KEY,
    constant              BIGINT      NOT NULL,
    variable              BIGINT      NOT NULL,
    fixed_wear_and_tear   BIGINT      NOT NULL,
    fixed_advanced        BIGINT      NOT NULL,
    k                     BIGINT      NOT NULL,
    fixed_component       BIGINT      NOT NULL,
    circulating_component BIGINT      NOT NULL,
    created_at            DATETIME(6) NOT NULL,
    INDEX idx_cost_prices_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS cost_prices;
