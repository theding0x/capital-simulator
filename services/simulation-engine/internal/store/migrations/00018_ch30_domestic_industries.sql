-- +goose Up
CREATE TABLE domestic_industries (
    id                  CHAR(24)     NOT NULL PRIMARY KEY,
    historical_stage_id CHAR(24)     NOT NULL,
    name                VARCHAR(255) NOT NULL,
    households_engaged  BIGINT       NOT NULL,
    annual_output_pence BIGINT       NOT NULL,
    created_at          DATETIME(6)  NOT NULL,
    CONSTRAINT fk_domestic_industries_stage
        FOREIGN KEY (historical_stage_id)
        REFERENCES historical_stages (id)
        ON DELETE CASCADE,
    KEY idx_domestic_industries_stage (historical_stage_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS domestic_industries;
