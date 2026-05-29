-- +goose Up
-- Vol. III Ch. 4 seed — Marx's worked capitals from Capital III, Ch. 4, so the
-- turnover panel comes up populated on a fresh MySQL volume.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=04 (chapter) and XX
-- is the row sequence. These 22-char IDs never collide with the 24-char hex
-- profit.NewTurnoverAnalysisID() output.
--
-- All three illustrate p' = s'·n·(v/C), the annual rate of profit:
--   Row 0401 Capital A: 80c+20v = 100C, n = 2, s' = 100% → annual p' = 40%
--     (single-turnover rate 20% doubled by turning over twice; S' = 200%).
--   Row 0402 Capital B: 160c+40v = 200C, n = 1, s' = 100% → annual p' = 20%.
--     Equal composition and s' to A, but half the annual rate — "the rates of
--     profit are related inversely as their periods of turnover."
--   Row 0403 Capital I: £10,000 fixed + £500 circ. constant + £500 v = £11,000C,
--     n = 10, s' = 100% → annual p' = 5,000/11,000 = 45 5/11% (4545 bp); single
--     turnover 454 bp; S' = 1000% (100000 bp); annual wage bill vn = £5,000.

INSERT INTO turnover_analyses
    (id, total_capital, variable_capital, surplus_value_rate_bp, turnovers,
     annual_profit_rate_bp, single_turnover_rate_bp, annual_surplus_value_rate_bp, annual_wages, created_at)
VALUES
    ('5eed000000000000000401', 100, 20, 10000, 2, 4000, 2000, 20000, 40,
     '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000402', 200, 40, 10000, 1, 2000, 2000, 10000, 40,
     '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000403', 11000, 500, 10000, 10, 4545, 454, 100000, 5000,
     '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM turnover_analyses WHERE id IN (
    '5eed000000000000000401',
    '5eed000000000000000402',
    '5eed000000000000000403'
);
