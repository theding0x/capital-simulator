-- +goose Up
-- Vol. III Ch. 21 — Interest-Bearing Capital seed data.
-- Two Marx-faithful exemplars from Capital Vol. III Ch. 21.
--
-- Seed IDs begin with 5eed000000000000002100 (chapter prefix 21).
-- Arithmetic:
--   Seed 2100 (£1,000 at 5% — canonical Marx fixture):
--     roundHalfUp(100000*500, 10000) = roundHalfUp(50000000, 10000) = 5000
--     money_returned = 100000 + 5000 = 105000
--
--   Seed 2101 (£5,000 at 15% — generic loanable capital):
--     roundHalfUp(500000*1500, 10000) = roundHalfUp(750000000, 10000) = 75000
--     money_returned = 500000 + 75000 = 575000
INSERT INTO `interest_bearing_capitals`
    (`id`, `money_advanced`, `interest_rate_bp`, `loan_term_days`, `interest_earned`, `money_returned`, `new_value_created`, `created_at`)
VALUES
    ('5eed000000000000002100', 100000, 500,  365, 5000,  105000, 0, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000002101', 500000, 1500, 180, 75000, 575000, 0, '2026-05-01 00:00:01.000000');

-- +goose Down
DELETE FROM `interest_bearing_capitals`
WHERE `id` IN (
    '5eed000000000000002100',
    '5eed000000000000002101'
);
