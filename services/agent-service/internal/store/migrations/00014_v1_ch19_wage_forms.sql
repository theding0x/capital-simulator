-- +goose Up
-- Schema for Capital Vol. I, Ch. 19 — The Transformation of the Value of
-- Labour-Power into Wages.
--
-- wage_forms stores one wage-form record per agent: the daily money wage
-- (daily_pence), the working day length (working_day_hours), and the
-- labour-power value it conceals (lpv_daily_pence, necessary_minutes).
-- Keyed by agent_id so there is at most one current wage form per agent.

CREATE TABLE wage_forms (
    id                   CHAR(24)     NOT NULL,
    agent_id             CHAR(24)     NOT NULL,
    daily_pence          BIGINT       NOT NULL,
    working_day_hours    BIGINT       NOT NULL,
    lpv_daily_pence      BIGINT       NOT NULL,
    necessary_minutes    BIGINT       NOT NULL,
    created_at           DATETIME(6)  NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_wage_forms_agent (agent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- +goose Down
DROP TABLE IF EXISTS wage_forms;
