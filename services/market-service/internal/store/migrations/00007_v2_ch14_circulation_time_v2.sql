-- +goose Up

-- Detailed sub-span breakdown of a selling phase.
-- Extends selling_phases from Ch. 5 without modifying that table.
CREATE TABLE selling_phase_details (
    id                        VARCHAR(24)  NOT NULL PRIMARY KEY,
    selling_phase_id          VARCHAR(24)  NOT NULL,
    turnover_time_id          VARCHAR(24)  NOT NULL,
    industrial_capital_id     VARCHAR(24)  NOT NULL,
    market_id                 VARCHAR(64)  NOT NULL DEFAULT '',
    market_distance_meters    BIGINT       NOT NULL DEFAULT 0,
    communication_lag_minutes BIGINT       NOT NULL DEFAULT 0,
    buyer_search_minutes      BIGINT       NOT NULL DEFAULT 0,
    bargaining_minutes        BIGINT       NOT NULL DEFAULT 0,
    legal_settlement_minutes  BIGINT       NOT NULL DEFAULT 0,
    total_minutes             BIGINT       NOT NULL DEFAULT 0,
    created_at                DATETIME(6)  NOT NULL,
    INDEX idx_spd_selling_phase (selling_phase_id),
    INDEX idx_spd_turnover (turnover_time_id),
    INDEX idx_spd_market (market_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Detailed sub-span breakdown of a buying phase.
CREATE TABLE buying_phase_details (
    id                        VARCHAR(24)  NOT NULL PRIMARY KEY,
    buying_phase_id           VARCHAR(24)  NOT NULL,
    turnover_time_id          VARCHAR(24)  NOT NULL,
    industrial_capital_id     VARCHAR(24)  NOT NULL,
    market_id                 VARCHAR(64)  NOT NULL DEFAULT '',
    market_distance_meters    BIGINT       NOT NULL DEFAULT 0,
    communication_lag_minutes BIGINT       NOT NULL DEFAULT 0,
    supplier_search_minutes   BIGINT       NOT NULL DEFAULT 0,
    order_processing_minutes  BIGINT       NOT NULL DEFAULT 0,
    legal_settlement_minutes  BIGINT       NOT NULL DEFAULT 0,
    total_minutes             BIGINT       NOT NULL DEFAULT 0,
    created_at                DATETIME(6)  NOT NULL,
    INDEX idx_bpd_buying_phase (buying_phase_id),
    INDEX idx_bpd_turnover (turnover_time_id),
    INDEX idx_bpd_market (market_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Median communication lag per kilometre for a market route.
CREATE TABLE distance_lag_relations (
    id                 VARCHAR(24)  NOT NULL PRIMARY KEY,
    market_id          VARCHAR(64)  NOT NULL,
    distance_meters    BIGINT       NOT NULL DEFAULT 0,
    lag_minutes_per_km BIGINT       NOT NULL DEFAULT 0,
    created_at         DATETIME(6)  NOT NULL,
    INDEX idx_dlr_market (market_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Historical events that changed the lag per km for a market route.
CREATE TABLE circulation_speed_improvements (
    id                     VARCHAR(24)  NOT NULL PRIMARY KEY,
    market_id              VARCHAR(64)  NOT NULL,
    old_lag_minutes_per_km BIGINT       NOT NULL DEFAULT 0,
    new_lag_minutes_per_km BIGINT       NOT NULL DEFAULT 0,
    reason                 VARCHAR(16)  NOT NULL DEFAULT 'other',
    effective_at           DATETIME(6)  NOT NULL,
    created_at             DATETIME(6)  NOT NULL,
    INDEX idx_csi_market (market_id),
    INDEX idx_csi_effective (market_id, effective_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Annual surplus-value penalty from circulation time. Full calc deferred to Ch. 16.
CREATE TABLE annual_surplus_penalties (
    id                              VARCHAR(24)  NOT NULL PRIMARY KEY,
    industrial_capital_id           VARCHAR(24)  NOT NULL,
    period                          VARCHAR(16)  NOT NULL,
    ideal_annual_rate_basis_points  BIGINT       NOT NULL DEFAULT 0,
    actual_annual_rate_basis_points BIGINT       NOT NULL DEFAULT 0,
    penalty_basis_points            BIGINT       NOT NULL DEFAULT 0,
    UNIQUE KEY uq_asp_ic_period (industrial_capital_id, period),
    INDEX idx_asp_ic (industrial_capital_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down

DROP TABLE IF EXISTS annual_surplus_penalties;
DROP TABLE IF EXISTS circulation_speed_improvements;
DROP TABLE IF EXISTS distance_lag_relations;
DROP TABLE IF EXISTS buying_phase_details;
DROP TABLE IF EXISTS selling_phase_details;
