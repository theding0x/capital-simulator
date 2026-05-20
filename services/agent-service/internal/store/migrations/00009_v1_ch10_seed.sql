-- +goose Up
-- Seed data for Capital Vol. I, Ch. 10 — The Working-Day.
--
-- working_days rows record (necessary, surplus) labour-minute pairs for each
-- regime Marx examines in §1, §3, and §6. Rate of surplus-value is derived
-- as surplus / necessary at read time.
--
-- relay_schedules rows model §4 Sundays-included relay system: two alternating
-- worker sets covering a 24-hour period. Worker IDs reference the labour_workers
-- seeded by 00007 (Mary 5eed00000000000000060001, Tom 5eed00000000000000060002).

INSERT IGNORE INTO working_days
    (id, necessary_labour_minutes, surplus_labour_minutes, created_at)
VALUES
    -- §1 working-day I: A–B = 6 h necessary, B–C = 1 h surplus.
    -- Rate of surplus-value = 1/6 ≈ 16.67%.
    ('5eed00000000000000100001', 360,  60, '2026-05-01 09:00:00.000000'),

    -- §1 working-day II: A–B = 6 h necessary, B–C = 3 h surplus.
    -- Rate of surplus-value = 3/6 = 50%.
    ('5eed00000000000000100002', 360, 180, '2026-05-01 09:30:00.000000'),

    -- Ch.7 §2 canonical 12-hour mill day: 6 h necessary + 6 h surplus.
    -- Rate of surplus-value = 100% — the rate cited throughout Ch.9.
    ('5eed00000000000000100003', 360, 360, '2026-05-01 10:00:00.000000'),

    -- §6 Factory Act 1847 (Ten Hours' Bill): 10 h working day (5 h + 5 h).
    -- Rate of surplus-value = 100% within the new statutory cap.
    ('5eed00000000000000100004', 300, 300, '2026-05-01 11:00:00.000000'),

    -- §3 maximum-day extreme: 6 h necessary + 12 h surplus = 18 h day.
    -- Rate of surplus-value = 200%; bounded by PhysicalMaxMinutes (1440).
    ('5eed00000000000000100005', 360, 720, '2026-05-01 12:00:00.000000');

-- §4 relay system: two sets alternating, day shift 6+6 and night shift 6+6,
-- jointly absorbing 24 h of machinery use. Mary works the day, Tom the night.
INSERT IGNORE INTO relay_schedules
    (id,
     shift_kind_0, nl_minutes_0, sl_minutes_0, worker_ids_0,
     shift_kind_1, nl_minutes_1, sl_minutes_1, worker_ids_1,
     created_at)
VALUES
    ('5eed00000000000000100101',
     'day',   360, 360, JSON_ARRAY('5eed00000000000000060001'),
     'night', 360, 360, JSON_ARRAY('5eed00000000000000060002'),
     '2026-05-01 13:00:00.000000');

-- +goose Down
DELETE FROM relay_schedules WHERE id = '5eed00000000000000100101';
DELETE FROM working_days WHERE id IN (
    '5eed00000000000000100001',
    '5eed00000000000000100002',
    '5eed00000000000000100003',
    '5eed00000000000000100004',
    '5eed00000000000000100005'
);
