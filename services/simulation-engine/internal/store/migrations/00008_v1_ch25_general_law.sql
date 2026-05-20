-- +goose Up
CREATE TABLE general_law_scenarios (
    id                  CHAR(24)     NOT NULL PRIMARY KEY,
    name                VARCHAR(255) NOT NULL,
    constant_capital    BIGINT       NOT NULL,
    variable_capital    BIGINT       NOT NULL,
    surplus_rate        DOUBLE       NOT NULL,
    accumulation_rate   DOUBLE       NOT NULL,
    productivity_growth DOUBLE       NOT NULL,
    wage_pence          BIGINT       NOT NULL,
    worker_supply       BIGINT       NOT NULL,
    periods             BIGINT       NOT NULL,
    created_at          DATETIME(6)  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS general_law_scenarios;
