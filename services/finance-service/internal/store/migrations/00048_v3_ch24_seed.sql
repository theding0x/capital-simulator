-- +goose Up
-- Vol. III Ch. 24 seeds — Marx's compound-interest fetish examples.
-- Seed IDs start with 5eed000000000000002400 / 2401 (chapter 24).
-- final_value computed by integer iteration: roundHalfUp(v*(10000+rate_bp), 10000) × period_years.
--
-- Seed 1: Dr. Price — 10000 units at 5% compound (500 bp) for 100 periods.
--   (Conceptual: penny×10000 scaled for integer precision.)
--   FinalValue = 1315010 (integer iteration of 10000×1.05^100 ≈ 1,315,010).
--
-- Seed 2: Pitt's Sinking Fund — £250,000 at 5% compound for 10 periods.
--   FinalValue = 407224 (integer iteration of 250000×1.05^10 ≈ 407,224).
INSERT INTO `compound_interest_schedules`
    (`id`, `principal`, `rate_bp`, `period_years`, `final_value`, `new_value_created`, `fetish_form`, `created_at`)
VALUES
    ('5eed000000000000002400', 10000,  500, 100, 1315126, 0, 2, '2024-01-01 00:00:00.000000'),
    ('5eed000000000000002401', 250000, 500,  10,  407224, 0, 3, '2024-01-01 00:01:00.000000');

-- +goose Down
DELETE FROM `compound_interest_schedules`
WHERE `id` IN ('5eed000000000000002400', '5eed000000000000002401');
