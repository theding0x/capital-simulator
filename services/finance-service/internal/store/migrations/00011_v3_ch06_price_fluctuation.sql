-- +goose Up
-- Vol. III Ch. 6 — price-fluctuation analyses. Each row records the effect of a
-- re-pricing of the raw-material element of constant capital on the rate of
-- profit. The constant capital is split into a fixed component the price movement
-- does not touch and a circulating raw-material component that scales with the
-- price; price_factor is that price as a percentage of its original level (100 =
-- unchanged, 200 = doubled, 50 = halved). v and s are held fixed, so the whole
-- movement of p' comes from the re-pricing of the material. The derived columns
-- hold the rate of profit before the movement s/(c+v) and after it s/(c'+v).
-- Capital magnitudes are LabourMinutes (BIGINT); the *_bp columns are
-- dimensionless basis points (100% = 10000 bp); kind is 'rise' or 'fall'.

CREATE TABLE IF NOT EXISTS price_fluctuation_analyses (
    id                        VARCHAR(24) NOT NULL PRIMARY KEY,
    kind                      VARCHAR(8)  NOT NULL,
    fixed_capital             BIGINT      NOT NULL,
    original_material_capital BIGINT      NOT NULL,
    price_factor              BIGINT      NOT NULL,
    variable_capital          BIGINT      NOT NULL,
    surplus_value             BIGINT      NOT NULL,
    old_profit_rate_bp        BIGINT      NOT NULL,
    new_profit_rate_bp        BIGINT      NOT NULL,
    created_at                DATETIME(6) NOT NULL,
    INDEX idx_price_fluctuation_analyses_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS price_fluctuation_analyses;
