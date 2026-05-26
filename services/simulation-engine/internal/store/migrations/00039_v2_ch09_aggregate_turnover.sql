-- +goose Up
-- Vol. II Ch. 9 — The Aggregate Turnover of Advanced Capital. Cycles of Turnover.
CREATE TABLE aggregate_turnovers (
    id                          VARCHAR(24)  NOT NULL PRIMARY KEY,
    industrial_capital_id       VARCHAR(24)  NOT NULL,
    lens                        VARCHAR(16)  NOT NULL DEFAULT 'money',
    advanced_total_pence        BIGINT       NOT NULL DEFAULT 0,
    annual_turned_over_pence    BIGINT       NOT NULL DEFAULT 0,
    aggregate_number_bp         BIGINT       NOT NULL DEFAULT 0,
    created_at                  DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE component_turnover_contributions (
    id                          VARCHAR(24)  NOT NULL PRIMARY KEY,
    aggregate_turnover_id       VARCHAR(24)  NOT NULL,
    capital_component_id        VARCHAR(24)  NOT NULL DEFAULT '',
    kind                        VARCHAR(16)  NOT NULL,
    advanced_pence              BIGINT       NOT NULL DEFAULT 0,
    turnover_number_bp          BIGINT       NOT NULL DEFAULT 0,
    annual_return_pence         BIGINT       NOT NULL DEFAULT 0,
    difference_kind             VARCHAR(16)  NOT NULL DEFAULT 'real',
    recorded_at                 DATETIME(6)  NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    CONSTRAINT fk_ctc_at FOREIGN KEY (aggregate_turnover_id) REFERENCES aggregate_turnovers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE lifetime_cycles (
    id                          VARCHAR(24)  NOT NULL PRIMARY KEY,
    aggregate_turnover_id       VARCHAR(24)  NOT NULL UNIQUE,
    duration_minutes            BIGINT       NOT NULL,
    started_at                  DATETIME(6)  NOT NULL,
    completed_cycles            BIGINT       NOT NULL DEFAULT 0,
    CONSTRAINT fk_lc_at FOREIGN KEY (aggregate_turnover_id) REFERENCES aggregate_turnovers(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE crisis_cycle_phase_records (
    id                          VARCHAR(24)  NOT NULL PRIMARY KEY,
    lifetime_cycle_id           VARCHAR(24)  NOT NULL,
    phase                       VARCHAR(32)  NOT NULL,
    observed_at                 DATETIME(6)  NOT NULL,
    notes                       TEXT         NOT NULL DEFAULT '',
    CONSTRAINT fk_ccpr_lc FOREIGN KEY (lifetime_cycle_id) REFERENCES lifetime_cycles(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS crisis_cycle_phase_records;
DROP TABLE IF EXISTS lifetime_cycles;
DROP TABLE IF EXISTS component_turnover_contributions;
DROP TABLE IF EXISTS aggregate_turnovers;
