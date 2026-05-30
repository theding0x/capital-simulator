-- +goose Up
-- Vol. III Ch. 13 — The Law as Such (Tendential Fall in the Rate of Profit).
-- Two tables: composition_trajectories and rate_mass_contradictions.
-- Trajectory periods are stored as JSON TEXT for simplicity; the domain layer
-- derives ProfitRates from them on read.

CREATE TABLE composition_trajectories (
    id                 VARCHAR(24)   NOT NULL,
    label              VARCHAR(120)  NOT NULL,
    surplus_value_rate BIGINT        NOT NULL, -- s' in basis points
    periods_json       TEXT          NOT NULL, -- JSON []TrajectoryPeriod
    created_at         DATETIME(6)   NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_composition_trajectories_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE rate_mass_contradictions (
    id          VARCHAR(24)  NOT NULL,
    old_c       BIGINT       NOT NULL,
    old_rate    BIGINT       NOT NULL,  -- basis points
    new_c       BIGINT       NOT NULL,
    new_rate    BIGINT       NOT NULL,  -- basis points
    old_mass    BIGINT       NOT NULL,
    new_mass    BIGINT       NOT NULL,
    mass_change BIGINT       NOT NULL,  -- signed
    created_at  DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_rate_mass_contradictions_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS rate_mass_contradictions;
DROP TABLE IF EXISTS composition_trajectories;
