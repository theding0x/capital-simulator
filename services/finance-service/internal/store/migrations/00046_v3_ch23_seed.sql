-- +goose Up
-- Vol. III Ch. 23 — Interest and Profit of Enterprise seed data.
-- Two Marx-faithful exemplars from Capital Vol. III Ch. 23.
--
-- Seed IDs begin with 5eed000000000000002300 (chapter prefix 23).
--
-- Seed 2301 — borrowed capital: avg profit 20% (2000 bp), interest 5% (500 bp),
--   ownership form 2 (Separated), profit of enterprise = 1500 bp.
--
-- Seed 2302 — own capital: avg profit 20% (2000 bp), interest 0 bp,
--   ownership form 1 (OwnerOperator), profit of enterprise = 2000 bp.
--   (owner and functioning capitalist are the same person; no interest paid)
INSERT INTO `profit_divisions`
    (`id`, `total_profit_bp`, `interest_bp`, `profit_of_enterprise_bp`, `ownership_form`, `wages_labour_minutes`, `mystification_index`, `created_at`)
VALUES
    ('5eed000000000000002301', 2000,  500, 1500, 2, 0,   5000, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000002302', 2000,    0, 2000, 1, 480, 2000, '2026-05-01 00:00:01.000000');

-- +goose Down
DELETE FROM `profit_divisions`
WHERE `id` IN (
    '5eed000000000000002301',
    '5eed000000000000002302'
);
