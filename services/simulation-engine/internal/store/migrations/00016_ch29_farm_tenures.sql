-- +goose Up
CREATE TABLE farm_tenures (
    id                     VARCHAR(24)  NOT NULL PRIMARY KEY,
    historical_stage_id    VARCHAR(24)  NOT NULL,
    form                   VARCHAR(32)  NOT NULL,
    lease_period_years     BIGINT       NOT NULL,
    rent_pence             BIGINT       NOT NULL DEFAULT 0,
    capital_advanced_pence BIGINT       NOT NULL DEFAULT 0,
    revenue_pence          BIGINT       NOT NULL DEFAULT 0,
    wage_costs_pence       BIGINT       NOT NULL DEFAULT 0,
    created_at             DATETIME(6)  NOT NULL,
    INDEX idx_farm_tenures_stage (historical_stage_id)
);

-- +goose Down
DROP TABLE IF EXISTS farm_tenures;
