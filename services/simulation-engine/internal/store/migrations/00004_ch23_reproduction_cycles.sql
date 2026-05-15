-- +goose Up
CREATE TABLE reproduction_cycles (
    id               CHAR(24)     NOT NULL PRIMARY KEY,
    constant_capital BIGINT       NOT NULL,
    variable_capital BIGINT       NOT NULL,
    surplus_rate     DOUBLE       NOT NULL,
    periods          BIGINT       NOT NULL,
    created_at       DATETIME(6)  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS reproduction_cycles;
