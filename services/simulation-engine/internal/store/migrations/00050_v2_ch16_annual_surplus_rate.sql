-- +goose Up
-- Vol. II Ch. 16 — The Turnover of Variable Capital
-- Schema DDL for the valorisation package tables.

CREATE TABLE variable_capital_advances (
    id                              CHAR(24)     NOT NULL,
    industrial_capital_id           VARCHAR(255) NOT NULL,
    pence                           BIGINT       NOT NULL,
    turnover_period_minutes         BIGINT       NOT NULL,
    turnovers_per_year_basis_points BIGINT       NOT NULL,
    created_at                      DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_vca_capital (industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- One per-circuit rate per advance (UNIQUE on advance_id enforces the
-- one-rate-per-advance invariant; ON DUPLICATE KEY UPDATE replaces it).
CREATE TABLE per_circuit_surplus_rates (
    id                          CHAR(24) NOT NULL,
    variable_capital_advance_id CHAR(24) NOT NULL,
    surplus_pence               BIGINT   NOT NULL,
    variable_pence              BIGINT   NOT NULL,
    basis_points                BIGINT   NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_pcsr_advance (variable_capital_advance_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Annual surplus rate snapshots (multiple periods per advance allowed).
CREATE TABLE annual_surplus_rates (
    id                          CHAR(24)    NOT NULL,
    variable_capital_advance_id CHAR(24)    NOT NULL,
    period                      VARCHAR(64) NOT NULL,
    n_basis_points              BIGINT      NOT NULL,
    per_circuit_basis_points    BIGINT      NOT NULL,
    annual_basis_points         BIGINT      NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_asr_advance (variable_capital_advance_id),
    INDEX idx_asr_period  (period)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Canonical A-vs-B contrast records.
CREATE TABLE annual_surplus_rate_contrasts (
    id                            CHAR(24) NOT NULL,
    capital_a_id                  CHAR(24) NOT NULL,
    capital_b_id                  CHAR(24) NOT NULL,
    capital_a_annual_basis_points BIGINT   NOT NULL,
    capital_b_annual_basis_points BIGINT   NOT NULL,
    ratio_basis_points            BIGINT   NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS annual_surplus_rate_contrasts;
DROP TABLE IF EXISTS annual_surplus_rates;
DROP TABLE IF EXISTS per_circuit_surplus_rates;
DROP TABLE IF EXISTS variable_capital_advances;
