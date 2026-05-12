-- +goose Up
-- Seed data for Capital Vol. I, Ch. 13 — Co-operation.
--
-- Two cooperations are seeded so the dashboard demonstrates the §3 fixture
-- (12 men × 12 h = 144 h collective working-day) and the §3 Burke-footnote
-- platoon (5 men, the minimum scale at which individual deviations cancel).
--
-- Co-operation requires (Capital Vol. I, Ch. 13 §1) a number of wage-
-- labourers brought under the direction of a single capitalist. The
-- existing `agents` seed (00003) provides James Watt as capitalist; we
-- add a second capitalist (Arkwright) plus enough worker agents to
-- populate the 12-man and 5-man cooperations.
--
-- Worker working-day per period: 720 min = 12 hours = the canonical
-- working-day from §3.

-- 1. Additional agents required for the cooperations.
INSERT INTO agents (id, name, class, money_balance, hoarding, labour_minutes, created_at, updated_at) VALUES
    -- A second capitalist who runs the larger cotton-spinning floor.
    ('5eed00000000000000130100', 'Richard Arkwright', 'capitalist', 200000, 0, 0, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),

    -- Twelve cotton spinners for the §3 collective-working-day fixture.
    ('5eed00000000000000130001', 'Spinner 01',  'worker', 720, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130002', 'Spinner 02',  'worker', 720, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130003', 'Spinner 03',  'worker', 720, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130004', 'Spinner 04',  'worker', 720, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130005', 'Spinner 05',  'worker', 720, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130006', 'Spinner 06',  'worker', 720, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130007', 'Spinner 07',  'worker', 720, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130008', 'Spinner 08',  'worker', 720, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130009', 'Spinner 09',  'worker', 720, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130010', 'Spinner 10',  'worker', 720, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130011', 'Overlooker A','worker', 720, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130012', 'Overlooker B','worker', 720, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),

    -- Five reapers for the Burke-footnote (§3 fn. 1) minimum platoon.
    ('5eed00000000000000130021', 'Reaper 01',   'worker', 600, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130022', 'Reaper 02',   'worker', 600, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130023', 'Reaper 03',   'worker', 600, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130024', 'Reaper 04',   'worker', 600, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000'),
    ('5eed00000000000000130025', 'Reaper 05',   'worker', 600, 0, 720, '2026-05-01 09:00:00.000000', '2026-05-01 09:00:00.000000');

-- 2. Cooperations.
--
-- Arkwright's Spinning Floor: 12 workers × 12 h = 144 h collective working-day
-- (Capital Vol. I, Ch. 13 §3). Two of the twelve are supervisory ("Overlooker A"
-- and "Overlooker B") — Marx's "officers (managers) and sergeants (foremen,
-- overlookers)" who direct the others on the capitalist's behalf (§4).
INSERT INTO cooperations (id, name, capitalist_id, members_json, created_at) VALUES
    ('5eed00000000000000130c01',
     'Arkwright Spinning Floor',
     '5eed00000000000000130100',
     JSON_ARRAY(
        JSON_OBJECT('worker_id','5eed00000000000000130001','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130002','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130003','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130004','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130005','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130006','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130007','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130008','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130009','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130010','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130011','supervisory',TRUE, 'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130012','supervisory',TRUE, 'working_day_minutes',720)
     ),
     '2026-05-01 10:00:00.000000'),

    -- Burke's reaper platoon: 5 workers under James Watt at harvest. The
    -- minimum scale (Capital Vol. I, Ch. 13 §3 fn. 1) at which individual
    -- variations average out into "average social labour." No supervisor.
    ('5eed00000000000000130c02',
     'Watt Harvest Platoon',
     '5eed00000000000000000001',
     JSON_ARRAY(
        JSON_OBJECT('worker_id','5eed00000000000000130021','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130022','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130023','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130024','supervisory',FALSE,'working_day_minutes',720),
        JSON_OBJECT('worker_id','5eed00000000000000130025','supervisory',FALSE,'working_day_minutes',720)
     ),
     '2026-05-01 10:30:00.000000');

-- +goose Down
DELETE FROM cooperations WHERE id IN (
    '5eed00000000000000130c01',
    '5eed00000000000000130c02'
);
DELETE FROM agents WHERE id IN (
    '5eed00000000000000130100',
    '5eed00000000000000130001',
    '5eed00000000000000130002',
    '5eed00000000000000130003',
    '5eed00000000000000130004',
    '5eed00000000000000130005',
    '5eed00000000000000130006',
    '5eed00000000000000130007',
    '5eed00000000000000130008',
    '5eed00000000000000130009',
    '5eed00000000000000130010',
    '5eed00000000000000130011',
    '5eed00000000000000130012',
    '5eed00000000000000130021',
    '5eed00000000000000130022',
    '5eed00000000000000130023',
    '5eed00000000000000130024',
    '5eed00000000000000130025'
);
