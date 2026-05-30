-- +goose Up
-- Vol. III Ch. 19 — Money-Dealing Capital seed data.
-- Three Marx-faithful exemplars illustrating the five technical operations.
--
-- Seed IDs begin with 5eed000000000000001901 (chapter prefix 19).
-- Arithmetic:
--   roundHalfUp(500000*1500, 10000) = (750000000 + 5000) / 10000 = 75000
--   roundHalfUp(200000*1500, 10000) = (300000000 + 5000) / 10000 = 30000
--   roundHalfUp(300000*1500, 10000) = (450000000 + 5000) / 10000 = 45000
INSERT INTO money_dealing_capitals
    (id, money_advanced, general_rate_bp, operation_kinds, money_dealing_profit, new_value_created, created_at)
VALUES
    ('5eed000000000000001901', 500000, 1500, '[1,2,3]', 75000, 0, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000001902', 200000, 1500, '[4]',     30000, 0, '2026-05-01 00:00:01.000000'),
    ('5eed000000000000001903', 300000, 1500, '[3,5]',   45000, 0, '2026-05-01 00:00:02.000000');

-- +goose Down
DELETE FROM money_dealing_capitals
WHERE id IN (
    '5eed000000000000001901',
    '5eed000000000000001902',
    '5eed000000000000001903'
);
