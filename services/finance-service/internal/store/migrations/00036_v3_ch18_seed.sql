-- +goose Up
-- Seed data for Vol. III Ch. 18 — Merchant Turnover.
-- Marx fixture: capital=£1,000 (100000p), general rate=15% (1500 bp).
-- annual = roundHalfUp(100000 * 1500, 10000) = 15000p (£150).
-- Row 1801: turnover=1, units=300  → markup = 15000 / (1*300)   = 50p per unit.
-- Row 1802: turnover=5, units=300  → markup = 15000 / (5*300)   = 10p per unit.
-- Row 1803: turnover=3, units=300  → markup = 15000 / (3*300)   = 16p per unit
--           (integer division: 15000/900 = 16, remainder 600 discarded).
INSERT INTO merchant_turnovers
    (id, money_advanced, general_rate_bp, turnover_count, units_per_turnover, annual_merchant_profit, markup_per_unit, created_at)
VALUES
    ('5eed000000000000001801', 100000, 1500, 1, 300, 15000, 50, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000001802', 100000, 1500, 5, 300, 15000, 10, '2026-05-01 00:00:01.000000'),
    ('5eed000000000000001803', 100000, 1500, 3, 300, 15000, 16, '2026-05-01 00:00:02.000000');

-- +goose Down
DELETE FROM merchant_turnovers
WHERE id IN (
    '5eed000000000000001801',
    '5eed000000000000001802',
    '5eed000000000000001803'
);
