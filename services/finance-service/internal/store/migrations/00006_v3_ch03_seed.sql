-- +goose Up
-- Vol. III Ch. 3 seed — Marx's worked variations from Capital III, Ch. 3, so
-- the profit-rate-relation panel comes up populated on a fresh MySQL volume.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=03 (chapter) and XX
-- is the row sequence. These 22-char IDs never collide with the 24-char hex
-- profit.NewVariationAnalysisID() output.
--
-- All three start from Marx's base capital 80c + 20v at s' = 100% (C = 100,
-- p' = 20%, v/C = 20%):
--   Row 0301 (§I.4a) s' constant, value-composition rises (c 80→100): p' falls
--     20% → 16⅔% — the normal case of modern industry.
--   Row 0302 (§I.4c) s' constant, value-composition falls (c 80→60): p' rises
--     20% → 25%.
--   Row 0303 (§II.1) value-composition held at 20% (80c+20v → 160c+40v), s'
--     halved 100% → 50%: p' falls 20% → 10% in exact proportion to s'.

INSERT INTO variation_analyses
    (id, vcase, initial_c, initial_v, initial_s_rate_bp, changed_c, changed_v, changed_s_rate_bp,
     old_profit_rate_bp, new_profit_rate_bp, old_composition_bp, new_composition_bp, proportional_change, created_at)
VALUES
    ('5eed000000000000000301', 's_constant_vc_variable', 80, 20, 10000, 100, 20, 10000,
     2000, 1666, 2000, 1666, 0, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000302', 's_constant_vc_variable', 80, 20, 10000, 60, 20, 10000,
     2000, 2500, 2000, 2500, 0, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000303', 's_variable_vc_constant', 80, 20, 10000, 160, 40, 5000,
     2000, 1000, 2000, 2000, 1, '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM variation_analyses WHERE id IN (
    '5eed000000000000000301',
    '5eed000000000000000302',
    '5eed000000000000000303'
);
