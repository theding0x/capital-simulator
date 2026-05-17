-- +goose Up
-- Ch. 32 — The Historical Tendency of Capitalist Accumulation.
-- Two tables: a trajectory header carries the initial/final aggregates;
-- centralisation_steps records each absorption event in step_index order.

CREATE TABLE accumulation_trajectories (
    id                    CHAR(24)     NOT NULL PRIMARY KEY,
    name                  VARCHAR(255) NOT NULL,
    initial_firms         BIGINT       NOT NULL,
    initial_capital_pence BIGINT       NOT NULL,
    final_firms           BIGINT       NOT NULL,
    final_capital_pence   BIGINT       NOT NULL,
    reserve_army_size     BIGINT       NOT NULL DEFAULT 0,
    created_at            DATETIME(6)  NOT NULL,
    UNIQUE KEY uq_accumulation_trajectories_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE centralisation_steps (
    id                         CHAR(24)    NOT NULL PRIMARY KEY,
    trajectory_id              CHAR(24)    NOT NULL,
    step_index                 BIGINT      NOT NULL,
    firms_absorbed             BIGINT      NOT NULL,
    capital_concentrated_pence BIGINT      NOT NULL,
    created_at                 DATETIME(6) NOT NULL,
    CONSTRAINT fk_centralisation_steps_trajectory
        FOREIGN KEY (trajectory_id)
        REFERENCES accumulation_trajectories (id)
        ON DELETE CASCADE,
    KEY idx_centralisation_steps_trajectory (trajectory_id, step_index)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- +goose Down
DROP TABLE IF EXISTS centralisation_steps;
DROP TABLE IF EXISTS accumulation_trajectories;
