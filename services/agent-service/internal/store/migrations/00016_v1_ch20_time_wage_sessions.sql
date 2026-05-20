-- +goose Up
-- Schema for Capital Vol. I, Ch. 20 — Time-Wages.
--
-- time_wage_sessions stores one record per working session: the daily value of
-- labour-power, the length of the working day, any overtime, and the wage period.
-- The hourly price and nominal wage are derived on read — not persisted — so they
-- always reflect the current formula.

CREATE TABLE time_wage_sessions (
    id                          CHAR(24)     NOT NULL,
    agent_id                    CHAR(24)     NOT NULL,
    daily_labour_power_value    BIGINT       NOT NULL,
    working_day_hours           BIGINT       NOT NULL,
    overtime_hours              BIGINT       NOT NULL DEFAULT 0,
    overtime_rate_pence         BIGINT       NOT NULL DEFAULT 0,
    wage_period                 VARCHAR(8)   NOT NULL,
    created_at                  DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    INDEX idx_time_wage_sessions_agent (agent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS time_wage_sessions;
