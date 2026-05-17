-- +goose Up
CREATE TABLE domestic_industries (
    id                  VARCHAR(24)  NOT NULL,
    historical_stage_id VARCHAR(24)  NOT NULL,
    name                VARCHAR(255) NOT NULL,
    households_engaged  BIGINT       NOT NULL,
    annual_output_pence BIGINT       NOT NULL,
    created_at          DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_domestic_industries_stage
        FOREIGN KEY (historical_stage_id)
        REFERENCES historical_stages (id)
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS domestic_industries;
