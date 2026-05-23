-- +goose Up
CREATE TABLE commodity_circuits (
    id                   VARCHAR(24)  NOT NULL,
    agent_id             VARCHAR(24)  NOT NULL DEFAULT '',
    constant_pence       BIGINT       NOT NULL,
    variable_pence       BIGINT       NOT NULL,
    surplus_pence        BIGINT       NOT NULL,
    pounds_total         BIGINT       NOT NULL DEFAULT 0,
    mode                 VARCHAR(16)  NOT NULL,
    is_first_investment  TINYINT(1)   NOT NULL DEFAULT 0,
    social_capital_lens  TINYINT(1)   NOT NULL DEFAULT 0,
    terminal_constant    BIGINT       NULL,
    terminal_variable    BIGINT       NULL,
    terminal_surplus     BIGINT       NULL,
    terminal_capitalised BIGINT       NULL,
    terminal_pounds      BIGINT       NULL,
    terminal_closed_at   DATETIME(6)  NULL,
    created_at           DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_commodity_circuits_agent (agent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE successive_partial_sales (
    id                   BIGINT       NOT NULL AUTO_INCREMENT,
    commodity_circuit_id VARCHAR(24)  NOT NULL,
    quantity             BIGINT       NOT NULL,
    realised_pence       BIGINT       NOT NULL,
    constant_share       BIGINT       NOT NULL,
    variable_share       BIGINT       NOT NULL,
    surplus_share        BIGINT       NOT NULL,
    sold_at              DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_sps_circuit_id (commodity_circuit_id),
    CONSTRAINT fk_sps_circuit FOREIGN KEY (commodity_circuit_id)
        REFERENCES commodity_circuits (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE mp_sources (
    id                          BIGINT       NOT NULL AUTO_INCREMENT,
    commodity_circuit_id        VARCHAR(24)  NOT NULL,
    source_commodity_circuit_id VARCHAR(24)  NOT NULL DEFAULT '',
    source_kind                 VARCHAR(32)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_mp_sources_circuit_id (commodity_circuit_id),
    CONSTRAINT fk_mp_sources_circuit FOREIGN KEY (commodity_circuit_id)
        REFERENCES commodity_circuits (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS mp_sources;
DROP TABLE IF EXISTS successive_partial_sales;
DROP TABLE IF EXISTS commodity_circuits;
