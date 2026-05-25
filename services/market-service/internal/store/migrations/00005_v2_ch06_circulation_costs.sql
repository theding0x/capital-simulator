-- +goose Up
-- Vol. II Ch. 6 — Costs of Circulation

CREATE TABLE circulation_costs (
    id                    CHAR(24)     NOT NULL,
    industrial_capital_id CHAR(24)     NOT NULL,
    kind                  VARCHAR(32)  NOT NULL,
    nature                VARCHAR(24)  NOT NULL,
    pence                 BIGINT       NOT NULL,
    incurred_at           DATETIME(3)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_cc_capital (industrial_capital_id),
    INDEX idx_cc_kind    (kind),
    INDEX idx_cc_nature  (nature)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE circulation_agents (
    id                       CHAR(24)     NOT NULL,
    industrial_capital_id    CHAR(24)     NOT NULL,
    agent_id                 VARCHAR(64)  NOT NULL DEFAULT '',
    role                     VARCHAR(24)  NOT NULL,
    wage_rate_pence          BIGINT       NOT NULL,
    labour_minutes_necessary BIGINT       NOT NULL,
    labour_minutes_surplus   BIGINT       NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_ca_capital (industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE money_as_faux_frais (
    id                       CHAR(24)  NOT NULL,
    industrial_capital_id    CHAR(24)  NOT NULL DEFAULT '',
    pence                    BIGINT    NOT NULL,
    annual_replacement_pence BIGINT    NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE commodity_supplies (
    id                    CHAR(24)     NOT NULL,
    industrial_capital_id CHAR(24)     NOT NULL,
    commodity_id          VARCHAR(64)  NOT NULL DEFAULT '',
    pence                 BIGINT       NOT NULL,
    held_since            DATETIME(3)  NOT NULL,
    is_voluntary          TINYINT(1)   NOT NULL DEFAULT 1,
    PRIMARY KEY (id),
    INDEX idx_cs_capital (industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE storage_costs (
    id                        CHAR(24)  NOT NULL,
    commodity_supply_id       CHAR(24)  NOT NULL,
    building_pence            BIGINT    NOT NULL DEFAULT 0,
    labour_pence              BIGINT    NOT NULL DEFAULT 0,
    preservation_labour_pence BIGINT    NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    INDEX idx_sc_supply (commodity_supply_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE transport_tariffs (
    id                                    CHAR(24)     NOT NULL,
    commodity_id                          VARCHAR(64)  NOT NULL,
    base_pence_per_ton_mile               BIGINT       NOT NULL,
    fragility_multiplier_basis_points     BIGINT       NOT NULL DEFAULT 0,
    breakage_risk_multiplier_basis_points BIGINT       NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    UNIQUE KEY uq_tt_commodity (commodity_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE transport_legs (
    id                       CHAR(24)     NOT NULL,
    commodity_id             VARCHAR(64)  NOT NULL DEFAULT '',
    quantity                 BIGINT       NOT NULL DEFAULT 1,
    origin_id                VARCHAR(128) NOT NULL DEFAULT '',
    destination_id           VARCHAR(128) NOT NULL DEFAULT '',
    distance_meters          BIGINT       NOT NULL DEFAULT 0,
    weight_grams             BIGINT       NOT NULL DEFAULT 0,
    volume_cubic_millimeters BIGINT       NOT NULL DEFAULT 0,
    labour_cost_pence        BIGINT       NOT NULL DEFAULT 0,
    means_of_transport_pence BIGINT       NOT NULL DEFAULT 0,
    surplus_pence            BIGINT       NOT NULL DEFAULT 0,
    PRIMARY KEY (id),
    INDEX idx_tl_commodity (commodity_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS transport_legs;
DROP TABLE IF EXISTS transport_tariffs;
DROP TABLE IF EXISTS storage_costs;
DROP TABLE IF EXISTS commodity_supplies;
DROP TABLE IF EXISTS money_as_faux_frais;
DROP TABLE IF EXISTS circulation_agents;
DROP TABLE IF EXISTS circulation_costs;
