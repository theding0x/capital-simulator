-- +goose Up
-- Vol. III Ch. 4 — the §Manchester worked example, restored as a seed so the
-- turnover panel ships the signature fixture the spec names (issue #404).
--
-- Marx's 1871 Manchester cotton-spinnery: C = £12,500, variable advanced
-- v = £318, s' = 153 11/13% (≈ 15385 bp), turning over n ≈ 8½ times a year.
-- The integer TurnoverCount cannot hold a half-turnover, so the model floors
-- Marx's 8½ to 8. The stored magnitudes are therefore the integer-model values
-- p' = s'·n·v/C = 15385·8·318/12500 = 3131 bp, S' = s'·n = 123080 bp,
-- single-turnover p' = 391 bp, annual wage bill vn = 2544 — the integer
-- approximation of Marx's idealised 8½ → ≈ 3327 bp / ≈ 130709 bp.
--
-- Seed ID 5eed00000000000000<CC><XX>: CC=04 (chapter), XX=04 (row), following
-- the three §main capitals seeded as 0401–0403 in 00008_v3_ch04_seed.sql.

INSERT INTO turnover_analyses
    (id, total_capital, variable_capital, surplus_value_rate_bp, turnovers,
     annual_profit_rate_bp, single_turnover_rate_bp, annual_surplus_value_rate_bp, annual_wages, created_at)
VALUES
    ('5eed000000000000000404', 12500, 318, 15385, 8, 3131, 391, 123080, 2544,
     '2026-06-03 00:00:00.000000');

-- +goose Down
DELETE FROM turnover_analyses WHERE id = '5eed000000000000000404';
