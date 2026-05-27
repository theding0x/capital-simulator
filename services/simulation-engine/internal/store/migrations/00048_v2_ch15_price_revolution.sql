-- +goose Up
-- Vol. II Ch. 15 — The Effects of a Change of Prices
-- Extends value_revolution_events with Ch. 15 fields and adds satellite tables.

ALTER TABLE value_revolution_events
    ADD COLUMN basis_points_change    BIGINT NOT NULL DEFAULT 0,
    ADD COLUMN revalued_capital_pence BIGINT NOT NULL DEFAULT 0;

-- RevaluedCapital: new advance required after an MP/LP price revolution.
CREATE TABLE revalued_capitals (
    id                        CHAR(24)  NOT NULL,
    value_revolution_event_id CHAR(24)  NOT NULL,
    industrial_capital_id     CHAR(24)  NOT NULL,
    old_advance_pence         BIGINT    NOT NULL,
    new_advance_pence         BIGINT    NOT NULL,
    delta_pence               BIGINT    NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_rc_event (value_revolution_event_id),
    INDEX idx_rc_capital (industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- RealisedValueDelta: gap between expected and actual output realisation.
CREATE TABLE realised_value_deltas (
    id                         CHAR(24) NOT NULL,
    value_revolution_event_id  CHAR(24) NOT NULL,
    industrial_capital_id      CHAR(24) NOT NULL,
    expected_realisation_pence BIGINT   NOT NULL,
    actual_realisation_pence   BIGINT   NOT NULL,
    delta_pence                BIGINT   NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_rvd_event (value_revolution_event_id),
    INDEX idx_rvd_capital (industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- InventoryRevaluation: paper gain/loss on held stock.
CREATE TABLE inventory_revaluations (
    id                    CHAR(24)     NOT NULL,
    industrial_capital_id CHAR(24)     NOT NULL,
    commodity_id          VARCHAR(255) NOT NULL,
    quantity_held         BIGINT       NOT NULL,
    old_unit_pence        BIGINT       NOT NULL,
    new_unit_pence        BIGINT       NOT NULL,
    delta_pence           BIGINT       NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_ir_capital (industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- SpeculativeHold: commodity held in anticipation of price movement.
CREATE TABLE speculative_holds (
    id                    CHAR(24)     NOT NULL,
    industrial_capital_id CHAR(24)     NOT NULL,
    commodity_id          VARCHAR(255) NOT NULL,
    held_since            DATETIME(6)  NOT NULL,
    anticipated_direction TINYINT      NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_sh_capital (industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- PriceCaseRecord: operator-selected fall/rise case for a revolution event.
CREATE TABLE price_case_records (
    value_revolution_event_id CHAR(24)    NOT NULL,
    fall_case                 VARCHAR(32) NULL,
    rise_case                 VARCHAR(32) NULL,
    PRIMARY KEY (value_revolution_event_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS price_case_records;
DROP TABLE IF EXISTS speculative_holds;
DROP TABLE IF EXISTS inventory_revaluations;
DROP TABLE IF EXISTS realised_value_deltas;
DROP TABLE IF EXISTS revalued_capitals;
ALTER TABLE value_revolution_events
    DROP COLUMN basis_points_change,
    DROP COLUMN revalued_capital_pence;
