-- +goose Up
-- Seed data for Capital Vol. I, Ch. 20 — Time-Wages.
--
-- Core Marx fixture (Ch. 20): "daily value of labour-power 3s., working-day 12 hours
-- → price of 1 working-hour = 3s./12 = 3d."
--
-- Session 1 (Thomas Hobson, standard 12-hour day):
--   daily_labour_power_value = 36 (3 shillings in pence)
--   working_day_hours        = 12
--   overtime_hours           = 0
--   overtime_rate_pence      = 0
--   wage_period              = 'daily'
--   → hourly price = 36/12 = 3d. (exact fraction)
--   → nominal wage = 36p
--
-- Session 2 (Thomas Hobson, extended 15-hour day — shows price falls):
--   daily_labour_power_value = 36
--   working_day_hours        = 15
--   overtime_hours           = 3  (at 4d./h — the overtime rate Marx cites)
--   overtime_rate_pence      = 4
--   wage_period              = 'daily'
--   → hourly price = 36/15 = 2.4d. (price falls despite higher gross pay)
--   → nominal wage = (36/15)*15 + 3*4 = 36 + 12 = 48p

INSERT INTO time_wage_sessions
    (id, agent_id, daily_labour_power_value, working_day_hours,
     overtime_hours, overtime_rate_pence, wage_period, created_at)
VALUES
    ('5eed000000000000002000f1',
     '5eed000000000000001900ff',
     36, 12, 0, 0, 'daily',
     '2026-05-14 09:00:00.000000'),
    ('5eed000000000000002000f2',
     '5eed000000000000001900ff',
     36, 15, 3, 4, 'daily',
     '2026-05-14 09:01:00.000000');

-- +goose Down
DELETE FROM time_wage_sessions
WHERE id IN ('5eed000000000000002000f1', '5eed000000000000002000f2');
