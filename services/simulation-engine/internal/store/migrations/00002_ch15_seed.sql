-- +goose Up
-- Seed data for Capital Vol. I, Ch. 15 — Machinery and Modern Industry.
--
-- Three exemplars cover the principal §-references:
--
--   1. Four "needle-machines" wired into a single steam-driven "Needle Works"
--      Factory. Marx, §8A: "a single needle-machine makes 145,000 in a
--      working-day of 11 hours. One woman superintends four such machines
--      producing near upon 600,000 needles a day." TotalProductivePower =
--      580,000 units/day; DailyValueTransfer = 4 × 600 = 2,400 LabourMinutes.
--
--   2. A steam-plough — §2: "a steam-plough does as much work in one hour
--      at a cost of three-pence, as 66 men at a cost of 15 shillings."
--      HandLabourPerUnit = 66 men × 60 min = 3,960 LabourMinutes of hand
--      labour displaced for each hour-unit of ploughing.
--
--   3. A power-loom — §2 spinning/weaving exemplar, fitted with a hand-
--      labour displacement that reflects the 180× saving (machine path: 24
--      LabourMinutes per lb. of yarn; hand path: 4,427 — we round the
--      hand_labour_per_unit to a clean 4,400 lab-min/unit for legibility).

-- 1. Machines.
INSERT INTO machines
    (id, name, motor_mechanism, transmitting_mechanism, working_tool,
     machine_value, lifespan_days, productive_power, hand_labour_per_unit,
     accumulated_wear, accumulated_depreciation, created_at) VALUES
    ('5eed00000000000000150m01',
     'Needle-machine I',
     'steam piston', 'shaft & belt', 'needle die',
     600000, 1000, 145000, 0,
     0, 0, '2026-05-08 09:00:00.000000'),
    ('5eed00000000000000150m02',
     'Needle-machine II',
     'steam piston', 'shaft & belt', 'needle die',
     600000, 1000, 145000, 0,
     0, 0, '2026-05-08 09:00:00.000000'),
    ('5eed00000000000000150m03',
     'Needle-machine III',
     'steam piston', 'shaft & belt', 'needle die',
     600000, 1000, 145000, 0,
     0, 0, '2026-05-08 09:00:00.000000'),
    ('5eed00000000000000150m04',
     'Needle-machine IV',
     'steam piston', 'shaft & belt', 'needle die',
     600000, 1000, 145000, 0,
     0, 0, '2026-05-08 09:00:00.000000'),
    ('5eed00000000000000150m05',
     'Steam-plough',
     'steam piston', 'traction cable', 'plough share',
     720000, 500, 1, 3960,
     0, 0, '2026-05-08 09:30:00.000000'),
    ('5eed00000000000000150m06',
     'Power-loom',
     'steam piston', 'overhead shaft', 'shuttle',
     480000, 2000, 366, 4400,
     0, 0, '2026-05-08 10:00:00.000000');

-- 2. Factories.
INSERT INTO factories
    (id, name, prime_mover_kind, prime_mover_horsepower, tick_count, created_at) VALUES
    ('5eed00000000000000150f01',
     'Needle Works',
     'steam', 30, 0,
     '2026-05-08 11:00:00.000000'),
    ('5eed00000000000000150f02',
     'Manchester Weaving Shed',
     'steam', 60, 0,
     '2026-05-08 11:30:00.000000');

-- 3. Factory machine memberships.
--
-- The Needle Works binds all four needle-machines under one supervisor
-- (Marx's §8A "one woman superintends four such machines"). The
-- Manchester Weaving Shed runs the power-loom on its own.
INSERT INTO factory_machines (factory_id, machine_id, ordinal) VALUES
    ('5eed00000000000000150f01', '5eed00000000000000150m01', 0),
    ('5eed00000000000000150f01', '5eed00000000000000150m02', 1),
    ('5eed00000000000000150f01', '5eed00000000000000150m03', 2),
    ('5eed00000000000000150f01', '5eed00000000000000150m04', 3),
    ('5eed00000000000000150f02', '5eed00000000000000150m06', 0);

-- +goose Down
DELETE FROM factory_machines WHERE factory_id IN (
    '5eed00000000000000150f01',
    '5eed00000000000000150f02'
);
DELETE FROM factories WHERE id IN (
    '5eed00000000000000150f01',
    '5eed00000000000000150f02'
);
DELETE FROM machines WHERE id IN (
    '5eed00000000000000150m01',
    '5eed00000000000000150m02',
    '5eed00000000000000150m03',
    '5eed00000000000000150m04',
    '5eed00000000000000150m05',
    '5eed00000000000000150m06'
);
