-- +goose Up
-- Vol. III Ch. 14 — Counteracting Influences.
-- Two tables: counteracting_forces (individual force records) and
-- counteracting_scenarios (scenarios applying ≥1 force to a base trajectory).
-- JSON TEXT columns store slice fields; domain layer owns serialisation.

CREATE TABLE counteracting_forces (
    id                         VARCHAR(24)   NOT NULL,
    kind                       VARCHAR(40)   NOT NULL,  -- CounteractingForceKind value
    excluded_from_general_rate TINYINT(1)    NOT NULL DEFAULT 0,
    note                       VARCHAR(255)  NOT NULL DEFAULT '',
    created_at                 DATETIME(6)   NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_counteracting_forces_kind (kind),
    INDEX idx_counteracting_forces_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE counteracting_scenarios (
    id                   VARCHAR(24)   NOT NULL,
    label                VARCHAR(120)  NOT NULL,
    forces_json          TEXT          NOT NULL,  -- JSON []CounteractingForce
    base_trajectory_json TEXT          NOT NULL,  -- JSON CompositionTrajectory
    initial_profit_rate  BIGINT        NOT NULL,  -- first period modified p'; basis points
    final_profit_rate    BIGINT        NOT NULL,  -- last period modified p'; basis points
    modified_rates_json  TEXT          NOT NULL,  -- JSON []int64 basis-point rates per period
    created_at           DATETIME(6)   NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_counteracting_scenarios_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS counteracting_scenarios;
DROP TABLE IF EXISTS counteracting_forces;
