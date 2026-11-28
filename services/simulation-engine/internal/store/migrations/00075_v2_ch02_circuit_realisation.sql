-- +goose Up
-- Vol. II Ch. 2 §I: C' is the exit form of a productive circuit's P, and "C'
-- must be sold at its value and in its entirety". Both tables are 1:1 with
-- productive_circuits (primary key = foreign key), so a circuit holds at most
-- one C' and at most one full C'—M' sale. Amounts are circulation.Pence, the
-- legacy 100-per-£1 scale (£422 = 42,200).
CREATE TABLE productive_circuit_commodities (
    productive_circuit_id   VARCHAR(24)  NOT NULL,
    value_original          BIGINT       NOT NULL,
    value_surplus           BIGINT       NOT NULL,
    PRIMARY KEY (productive_circuit_id),
    CONSTRAINT fk_pc_commodity_pc FOREIGN KEY (productive_circuit_id)
        REFERENCES productive_circuits (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE productive_circuit_realisations (
    productive_circuit_id   VARCHAR(24)  NOT NULL,
    realised_pence          BIGINT       NOT NULL,
    sold_at                 DATETIME(6)  NOT NULL,
    PRIMARY KEY (productive_circuit_id),
    CONSTRAINT fk_pc_realisation_pc FOREIGN KEY (productive_circuit_id)
        REFERENCES productive_circuits (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS productive_circuit_realisations;
DROP TABLE IF EXISTS productive_circuit_commodities;
