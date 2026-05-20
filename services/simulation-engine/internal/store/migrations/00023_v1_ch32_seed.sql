-- +goose Up
-- Ch. 32 seed: long-run centralisation trajectories.
-- Three Marx-anchored series whose initial conditions approximate the
-- industries Marx and his readers had in view at publication (1867):
--
--   Lancashire cotton, 1820-1880 — many small mills consolidate into
--     a handful of giants, freeing labourers into the reserve army.
--   English banking, 1810-1900 — the country banks are absorbed by
--     the joint-stock banks; capital concentrates while assets grow.
--   American railroads, 1860-1900 — small lines merged by the trusts
--     into the great trunk roads; Marx mentions the railroads as the
--     paradigm of centralisation in §1.
--
-- Seed IDs follow 5eed00000000000000<CC><XX> where CC=32. Trajectory
-- headers use prefix 320, step rows use prefix 321/322/323.

INSERT INTO accumulation_trajectories
    (id, name, initial_firms, initial_capital_pence, final_firms, final_capital_pence, reserve_army_size, created_at)
VALUES
    ('5eed000000000000003201',
     'Lancashire cotton, 1820-1880',
     1000, 480000000, 217, 614400000, 17400,
     '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003202',
     'English banking, 1810-1900',
     650, 240000000, 110, 384000000, 9270,
     '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003203',
     'American railroads, 1860-1900',
     420, 9600000000, 70, 19200000000, 24000,
     '1867-01-01 00:00:00.000000');

-- Lancashire cotton steps: six absorption events over six decades.
INSERT INTO centralisation_steps
    (id, trajectory_id, step_index, firms_absorbed, capital_concentrated_pence, created_at)
VALUES
    ('5eed000000000000003211', '5eed000000000000003201', 1, 200, 96000000,  '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003212', '5eed000000000000003201', 2, 200, 100800000, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003213', '5eed000000000000003201', 3, 200, 105840000, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003214', '5eed000000000000003201', 4, 150, 83349000,  '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003215', '5eed000000000000003201', 5, 100, 58344300,  '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003216', '5eed000000000000003201', 6, 133, 81600000,  '1867-01-01 00:00:00.000000');

-- English banking steps: consolidation of country banks into joint-stocks.
INSERT INTO centralisation_steps
    (id, trajectory_id, step_index, firms_absorbed, capital_concentrated_pence, created_at)
VALUES
    ('5eed000000000000003221', '5eed000000000000003202', 1, 130, 48000000, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003222', '5eed000000000000003202', 2, 130, 50400000, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003223', '5eed000000000000003202', 3, 130, 52920000, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003224', '5eed000000000000003202', 4, 100, 41580000, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003225', '5eed000000000000003202', 5,  50, 21800000, '1867-01-01 00:00:00.000000');

-- American railroads steps: the trunk lines absorb the small carriers.
INSERT INTO centralisation_steps
    (id, trajectory_id, step_index, firms_absorbed, capital_concentrated_pence, created_at)
VALUES
    ('5eed000000000000003231', '5eed000000000000003203', 1,  84, 1920000000, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003232', '5eed000000000000003203', 2,  84, 2016000000, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003233', '5eed000000000000003203', 3,  84, 2116800000, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003234', '5eed000000000000003203', 4,  60, 1620000000, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003235', '5eed000000000000003203', 5,  38, 1056000000, '1867-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM centralisation_steps WHERE id IN (
    '5eed000000000000003211','5eed000000000000003212','5eed000000000000003213',
    '5eed000000000000003214','5eed000000000000003215','5eed000000000000003216',
    '5eed000000000000003221','5eed000000000000003222','5eed000000000000003223',
    '5eed000000000000003224','5eed000000000000003225',
    '5eed000000000000003231','5eed000000000000003232','5eed000000000000003233',
    '5eed000000000000003234','5eed000000000000003235'
);
DELETE FROM accumulation_trajectories WHERE id IN (
    '5eed000000000000003201',
    '5eed000000000000003202',
    '5eed000000000000003203'
);
