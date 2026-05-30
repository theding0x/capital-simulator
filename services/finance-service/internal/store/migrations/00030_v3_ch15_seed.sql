-- +goose Up
-- Vol. III Ch. 15 seed — three crises + one contradiction per kind.
-- Seed IDs: 5eed000000000000001501–5eed000000000000001507.
-- These 22-char IDs never collide with the 24-char hex output of New*ID().
-- Crisis formula: post = roundHalfUp(pre*10000 / (10000 - writedown*100))
--   §crisis-A: writedown=30, pre=1000 → numer=10_000_000 denom=7000 → 1429
--   §crisis-B: writedown=20, pre=800  → numer=8_000_000  denom=8000 → 1000
--   §crisis-C: writedown=50, pre=600  → numer=6_000_000  denom=5000 → 1200

INSERT INTO crises
    (id, constant_capital_writedown, pre_crisis_profit_rate, post_crisis_profit_rate, created_at)
VALUES
    ('5eed000000000000001501', 30, 1000, 1429, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000001502', 20,  800, 1000, '2026-05-01 00:00:01.000000'),
    ('5eed000000000000001503', 50,  600, 1200, '2026-05-01 00:00:02.000000');

INSERT INTO internal_contradictions
    (id, kind, is_coexistent, note, created_at)
VALUES
    ('5eed000000000000001504', 'falling_rate_rising_mass',                  1, 'Marx §I: falling rate and rising mass are two sides of the same movement',      '2026-05-01 00:00:00.000000'),
    ('5eed000000000000001505', 'overproduction_underconsumption',           1, 'Marx §II: rising productivity expands commodities beyond consumption limits',   '2026-05-01 00:00:01.000000'),
    ('5eed000000000000001506', 'overaccumulation_relative_overpopulation',  1, 'Marx §III: excess capital cannot be deployed at the average rate of profit',    '2026-05-01 00:00:02.000000'),
    ('5eed000000000000001507', 'concentration_centralisation',              1, 'Marx §IV: falling rate accelerates concentration and centralisation of capital', '2026-05-01 00:00:03.000000');

-- +goose Down
DELETE FROM internal_contradictions WHERE id IN (
    '5eed000000000000001504',
    '5eed000000000000001505',
    '5eed000000000000001506',
    '5eed000000000000001507'
);
DELETE FROM crises WHERE id IN (
    '5eed000000000000001501',
    '5eed000000000000001502',
    '5eed000000000000001503'
);
