-- +goose Up
-- Ch.23 and Ch.24 endpoints remain stateless compute endpoints; no Go code
-- reads or writes reproduction_cycles or accumulation_scenarios. Dropping
-- these tables removes the false impression of persistence for those chapters.
DROP TABLE IF EXISTS reproduction_cycles;
DROP TABLE IF EXISTS accumulation_scenarios;

-- +goose Down
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

CREATE TABLE reproduction_cycles (
    id               CHAR(24)     NOT NULL PRIMARY KEY,
    constant_capital BIGINT       NOT NULL,
    variable_capital BIGINT       NOT NULL,
    surplus_rate     DOUBLE       NOT NULL,
    periods          BIGINT       NOT NULL,
    created_at       DATETIME(6)  NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
