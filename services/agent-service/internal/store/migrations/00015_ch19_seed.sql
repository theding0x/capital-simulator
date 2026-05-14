-- +goose Up
-- Seed data for Capital Vol. I, Ch. 19 — The Transformation of the Value of
-- Labour-Power into Wages.
--
-- Core Marx fixture (Ch. 19): "daily value of labour-power is 3s., the value
-- of the product of 6 working-hours; working-day is 12 hours."
--
--   Worker: Thomas Hobson — the seed worker whose wage form embodies the
--     standard 3s/12h working day from the chapter text.
--
--   WageForm:
--     daily_pence        = 36  (3 shillings × 12 pence)
--     working_day_hours  = 12
--     lpv_daily_pence    = 36  (value of labour-power = 3s.)
--     necessary_minutes  = 360 (6 hours necessary labour)
--
-- Under Appearance the wage presents all 12h as paid;
-- under Decompose only 6h is necessary, 6h is surplus.

INSERT INTO labour_workers
    (id, kind, owns_labour_power, owns_commodities_to_sell,
     capacity_minutes_per_day, labour_power_value_minutes, created_at, updated_at)
VALUES
    ('5eed000000000000001900ff', 'worker', 1, 0, 720, 360,
     '2026-05-14 09:00:00.000000', '2026-05-14 09:00:00.000000');

INSERT INTO wage_forms
    (id, agent_id, daily_pence, working_day_hours, lpv_daily_pence, necessary_minutes, created_at)
VALUES
    ('5eed000000000000001900f1',
     '5eed000000000000001900ff',
     36, 12, 36, 360,
     '2026-05-14 09:00:00.000000');

-- +goose Down
DELETE FROM wage_forms WHERE id = '5eed000000000000001900f1';
DELETE FROM labour_workers WHERE id = '5eed000000000000001900ff';
