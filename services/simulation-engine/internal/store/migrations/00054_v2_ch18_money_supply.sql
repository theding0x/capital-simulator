-- +goose Up
-- Vol. II Ch. 18 — The Role of Money-Capital in Reproduction
-- Tables for the money-supply apportionment model, department reserves,
-- circulating money mass (M×V), wage-rotation funds, and inter-department
-- settlement flows.

CREATE TABLE IF NOT EXISTS money_supply_apportionments (
    id                     VARCHAR(24)  NOT NULL PRIMARY KEY,
    total_social_money_pence BIGINT     NOT NULL,
    dept_i_reserve_pence   BIGINT       NOT NULL,
    dept_ii_reserve_pence  BIGINT       NOT NULL,
    wage_rotation_pence    BIGINT       NOT NULL,
    idle_hoard_pence       BIGINT       NOT NULL,
    period                 VARCHAR(20)  NOT NULL,
    created_at             TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS department_money_reserves (
    id              VARCHAR(24)  NOT NULL PRIMARY KEY,
    department      VARCHAR(20)  NOT NULL,
    reserve_pence   BIGINT       NOT NULL,
    reserve_purpose VARCHAR(30)  NOT NULL,
    period          VARCHAR(20)  NOT NULL,
    created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS circulating_money_masses (
    id                                VARCHAR(24)  NOT NULL PRIMARY KEY,
    money_stock_pence                 BIGINT       NOT NULL,
    velocity_per_year_basis_points    BIGINT       NOT NULL,
    effective_circulating_value_pence BIGINT       NOT NULL,
    period                            VARCHAR(20)  NOT NULL,
    created_at                        TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS wage_rotation_funds (
    id                   VARCHAR(24)  NOT NULL PRIMARY KEY,
    fund_pence           BIGINT       NOT NULL,
    wage_cycle_frequency BIGINT       NOT NULL,
    department           VARCHAR(20)  NOT NULL,
    period               VARCHAR(20)  NOT NULL,
    created_at           TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS inter_department_settlements (
    id                 VARCHAR(24)  NOT NULL PRIMARY KEY,
    from_department    VARCHAR(20)  NOT NULL,
    to_department      VARCHAR(20)  NOT NULL,
    settled_pence      BIGINT       NOT NULL,
    settlement_purpose VARCHAR(30)  NOT NULL,
    period             VARCHAR(20)  NOT NULL,
    created_at         TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS inter_department_settlements;
DROP TABLE IF EXISTS wage_rotation_funds;
DROP TABLE IF EXISTS circulating_money_masses;
DROP TABLE IF EXISTS department_money_reserves;
DROP TABLE IF EXISTS money_supply_apportionments;
