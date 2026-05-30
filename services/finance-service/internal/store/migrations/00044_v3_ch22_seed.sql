-- +goose Up
-- Vol. III Ch. 22 — Division of Profit; Rate of Interest seed data.
-- Three Marx-faithful exemplars from Capital Vol. III Ch. 22.
--
-- Seed IDs begin with 5eed000000000000002200 (chapter prefix 22).
--
-- Seed 2201 — canonical fixture: avg profit 20% (2000 bp), interest ¼ = 500 bp,
--   cycle phase 3 (Prosperity), profit of enterprise = 1500 bp.
--
-- Seed 2202 — Prosperity 1843: interest fell to 150 bp (1.5%), avg 2000 bp (20%).
--   cycle phase 3 (Prosperity).
--
-- Seed 2203 — Crisis 1847: interest rose to 800 bp (8%+), avg 2000 bp (20%).
--   cycle phase 5 (Crisis).
INSERT INTO `rate_of_interests`
    (`id`, `rate_bp`, `average_profit_rate_bp`, `cycle_phase`, `period`, `created_at`)
VALUES
    ('5eed000000000000002201', 500,  2000, 3, 'Marx canonical — avg 20% interest ¼',  '2026-05-01 00:00:00.000000'),
    ('5eed000000000000002202', 150,  2000, 3, '1843 prosperity',                       '2026-05-01 00:00:01.000000'),
    ('5eed000000000000002203', 800,  2000, 5, '1847 crisis',                           '2026-05-01 00:00:02.000000');

-- +goose Down
DELETE FROM `rate_of_interests`
WHERE `id` IN (
    '5eed000000000000002201',
    '5eed000000000000002202',
    '5eed000000000000002203'
);
