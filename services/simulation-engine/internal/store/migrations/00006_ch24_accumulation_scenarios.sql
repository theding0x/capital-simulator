-- +goose Up
CREATE TABLE accumulation_scenarios (
    id                CHAR(24)     NOT NULL PRIMARY KEY,
    name              VARCHAR(255) NOT NULL,
    constant_capital  BIGINT       NOT NULL,
    variable_capital  BIGINT       NOT NULL,
    surplus_rate      DOUBLE       NOT NULL,
    accum_rate        DOUBLE       NOT NULL,
    composition_ratio DOUBLE       NOT NULL,
    periods           BIGINT       NOT NULL,
    created_at        DATETIME(6)  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS accumulation_scenarios;
