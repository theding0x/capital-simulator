-- +goose Up
-- Vol. II Ch. 20 — Simple Reproduction.
-- Tables for the two-department reproduction scheme: SimpleReproductionScheme
-- (root), DepartmentalCapital (c+v+s per department), InterDepartmentExchange
-- (the three value flows), ReproductionTick (cycle advancement), and
-- MoneyClosedLoop (closure verification).
--
-- Keystone equation: I(v+s) == II(c)
-- Marx's canonical fixture: Dept I 4000c+1000v+1000s=6000 | Dept II 2000c+500v+500s=3000

CREATE TABLE IF NOT EXISTS simple_reproduction_schemes (
    id          VARCHAR(24)  NOT NULL PRIMARY KEY,
    period      VARCHAR(20)  NOT NULL,
    is_balanced TINYINT(1)   NOT NULL DEFAULT 0,
    tick_count  BIGINT       NOT NULL DEFAULT 0,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS departmental_capitals (
    id             VARCHAR(24)  NOT NULL PRIMARY KEY,
    scheme_id      VARCHAR(24)  NOT NULL,
    department     VARCHAR(20)  NOT NULL,
    constant_pence BIGINT       NOT NULL,
    variable_pence BIGINT       NOT NULL,
    surplus_pence  BIGINT       NOT NULL,
    total_pence    BIGINT       NOT NULL,
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uq_scheme_dept (scheme_id, department),
    CONSTRAINT fk_dept_cap_scheme
        FOREIGN KEY (scheme_id) REFERENCES simple_reproduction_schemes (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS interdepartment_exchanges (
    id              VARCHAR(24)   NOT NULL PRIMARY KEY,
    scheme_id       VARCHAR(24)   NOT NULL,
    from_department VARCHAR(20)   NOT NULL,
    to_department   VARCHAR(20)   NOT NULL,
    pence           BIGINT        NOT NULL,
    kind            VARCHAR(20)   NOT NULL,
    description     VARCHAR(255)  NOT NULL DEFAULT '',
    created_at      TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_exchange_scheme
        FOREIGN KEY (scheme_id) REFERENCES simple_reproduction_schemes (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS reproduction_ticks (
    id          VARCHAR(24)  NOT NULL PRIMARY KEY,
    scheme_id   VARCHAR(24)  NOT NULL,
    tick_number BIGINT       NOT NULL,
    period      VARCHAR(20)  NOT NULL,
    is_balanced TINYINT(1)   NOT NULL DEFAULT 0,
    created_at  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_tick_scheme
        FOREIGN KEY (scheme_id) REFERENCES simple_reproduction_schemes (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS money_closed_loops (
    id             VARCHAR(24)  NOT NULL PRIMARY KEY,
    scheme_id      VARCHAR(24)  NOT NULL,
    period         VARCHAR(20)  NOT NULL,
    is_closed      TINYINT(1)   NOT NULL DEFAULT 0,
    net_flow_pence BIGINT       NOT NULL DEFAULT 0,
    created_at     TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_loop_scheme
        FOREIGN KEY (scheme_id) REFERENCES simple_reproduction_schemes (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS money_closed_loops;
DROP TABLE IF EXISTS reproduction_ticks;
DROP TABLE IF EXISTS interdepartment_exchanges;
DROP TABLE IF EXISTS departmental_capitals;
DROP TABLE IF EXISTS simple_reproduction_schemes;
