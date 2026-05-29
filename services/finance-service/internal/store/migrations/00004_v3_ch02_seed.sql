-- +goose Up
-- Vol. III Ch. 2 seed — Marx's worked example from Capital III, Ch. 2, so the
-- rate-of-profit panel comes up populated on a fresh MySQL volume.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=02 (chapter) and XX
-- is the row sequence. These 22-char IDs never collide with the 24-char hex
-- profit.NewProfitRateID() output.
--
-- Row 0201 is the canonical example: 400c + 100v + 100s. C = 500, so the rate
-- of profit p' = s/C = 100/500 = 20% (2000 bp) while the rate of surplus-value
-- s' = s/v = 100/100 = 100% (10000 bp); the profit-form conceals four fifths of
-- the rate (mystification 8000 bp).
-- Row 0202 is Marx's limiting case v = C (no constant capital): 0c + 100v +
-- 100s. C = 100, so p' = s' = 100% and nothing is mystified (0 bp) — the only
-- case in which the two rates coincide.

INSERT INTO profit_rates
    (id, constant, variable, surplus_value, total_capital, profit_rate_bp, surplus_value_rate_bp, mystification_bp, created_at)
VALUES
    ('5eed000000000000000201', 400, 100, 100, 500, 2000, 10000, 8000,
     '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000202', 0, 100, 100, 100, 10000, 10000, 0,
     '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM profit_rates WHERE id IN (
    '5eed000000000000000201',
    '5eed000000000000000202'
);
