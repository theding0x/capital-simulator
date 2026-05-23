-- +goose Up
CREATE TABLE productive_circuits (
    id                                  VARCHAR(24)  NOT NULL,
    agent_id                            VARCHAR(24)  NOT NULL DEFAULT '',
    constant_pence                      BIGINT       NOT NULL,
    variable_pence                      BIGINT       NOT NULL,
    mode                                VARCHAR(16)  NOT NULL,
    min_capitalisation_increment_pence  BIGINT       NOT NULL,
    latent_accumulated                  BIGINT       NOT NULL DEFAULT 0,
    reserve_balance                     BIGINT       NOT NULL DEFAULT 0,
    created_at                          DATETIME(6)  NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE revenue_circuits (
    id                      BIGINT       NOT NULL AUTO_INCREMENT,
    productive_circuit_id   VARCHAR(24)  NOT NULL,
    amount                  BIGINT       NOT NULL,
    spent_at                DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_revenue_circuits_pc_id (productive_circuit_id),
    CONSTRAINT fk_revenue_circuit_pc FOREIGN KEY (productive_circuit_id)
        REFERENCES productive_circuits (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE capitalisation_steps (
    id                      BIGINT       NOT NULL AUTO_INCREMENT,
    productive_circuit_id   VARCHAR(24)  NOT NULL,
    amount_injected         BIGINT       NOT NULL,
    delta_constant_pence    BIGINT       NOT NULL,
    delta_variable_pence    BIGINT       NOT NULL,
    occurred_at             DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_capitalisation_steps_pc_id (productive_circuit_id),
    CONSTRAINT fk_capitalisation_step_pc FOREIGN KEY (productive_circuit_id)
        REFERENCES productive_circuits (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE reserve_draws (
    id                      BIGINT       NOT NULL AUTO_INCREMENT,
    productive_circuit_id   VARCHAR(24)  NOT NULL,
    drawn_pence             BIGINT       NOT NULL,
    reason                  VARCHAR(32)  NOT NULL,
    occurred_at             DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_reserve_draws_pc_id (productive_circuit_id),
    CONSTRAINT fk_reserve_draw_pc FOREIGN KEY (productive_circuit_id)
        REFERENCES productive_circuits (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS reserve_draws;
DROP TABLE IF EXISTS capitalisation_steps;
DROP TABLE IF EXISTS revenue_circuits;
DROP TABLE IF EXISTS productive_circuits;
