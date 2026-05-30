-- +goose Up
-- Vol. III Ch. 17 seed — three commercial-profit exemplars.
-- Seed IDs: 5eed000000000000001701–5eed000000000000001703 (22-char, never collide with 24-char NewID()).
--
-- Row 1: productive=720000, commercial=180000, surplus=144000
--   unadjusted = (144000*10000 + 360000) / 720000 = 2000 bp (20%)
--   adjusted   = (144000*10000 + 450000) / 900000 = 1600 bp (16%)
--
-- Row 2: same productive/commercial split as row 1 but surplus halved (72000)
--   unadjusted = (72000*10000 + 360000) / 720000 = 1000 bp (10%)
--   adjusted   = (72000*10000 + 450000) / 900000 = 800 bp (8%)
--
-- Row 3: productive=800000, commercial=100000, surplus=160000
--   unadjusted = (160000*10000 + 400000) / 800000 = 2000 bp (20%)
--   adjusted   = (160000*10000 + 450000) / 900000 = 1778 bp (17.78%)

INSERT INTO commercial_profits
    (id, productive_capital, commercial_capital, surplus_value, unadjusted_rate_bp, adjusted_rate_bp, new_value_created, created_at)
VALUES
    ('5eed000000000000001701', 720000, 180000, 144000, 2000, 1600, 0, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000001702', 720000, 180000,  72000, 1000,  800, 0, '2026-05-01 00:00:01.000000'),
    ('5eed000000000000001703', 800000, 100000, 160000, 2000, 1778, 0, '2026-05-01 00:00:02.000000');

-- +goose Down
DELETE FROM commercial_profits WHERE id IN (
    '5eed000000000000001701',
    '5eed000000000000001702',
    '5eed000000000000001703'
);
