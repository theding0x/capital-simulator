-- +goose Up
-- Vol. II Ch. 7 — The Turnover Time and the Number of Turnovers

CREATE TABLE turnovers (
    id                    CHAR(24)     NOT NULL,
    industrial_capital_id VARCHAR(128) NOT NULL DEFAULT '',
    lens                  VARCHAR(16)  NOT NULL,
    turnover_time_minutes BIGINT       NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    INDEX idx_tv_capital (industrial_capital_id),
    INDEX idx_tv_lens    (lens)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE turnover_cycles (
    id                  CHAR(24)    NOT NULL,
    turnover_id         CHAR(24)    NOT NULL,
    started_at          DATETIME(3) NOT NULL,
    ended_at            DATETIME(3) NOT NULL,
    advance_pence       BIGINT      NOT NULL DEFAULT 0,
    returned_pence      BIGINT      NOT NULL DEFAULT 0,
    production_minutes  BIGINT      NOT NULL DEFAULT 0,
    circulation_minutes BIGINT      NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    INDEX idx_tc_turnover (turnover_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE turnover_numbers (
    id                     CHAR(24) NOT NULL,
    turnover_id            CHAR(24) NOT NULL,
    year_reference_minutes BIGINT   NOT NULL DEFAULT 525600,
    turnover_time_minutes  BIGINT   NOT NULL DEFAULT 0,
    basis_points           BIGINT   NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE KEY uq_tn_turnover (turnover_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS turnover_numbers;
DROP TABLE IF EXISTS turnover_cycles;
DROP TABLE IF EXISTS turnovers;
