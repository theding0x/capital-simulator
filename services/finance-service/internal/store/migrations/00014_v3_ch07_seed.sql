-- +goose Up
-- Vol. III Ch. 7 seed — worked composition comparisons and magnitude changes, so
-- the supplementary-remarks panel comes up populated on a fresh MySQL volume.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=07 (chapter) and XX is
-- the row sequence. These 22-char IDs never collide with the 24-char hex
-- profit.New*ID() output.
--
-- composition_effects
--   Row 0701: Marx's illustration — equal s = £1,000 on equal v = £1,000;
--     A carries c = £9,000 (p′ = 10% = 1000 bp), B carries c = £11,000 (8⅓% =
--     833 bp). The £2,000 of constant capital is the whole 167 bp difference.
--   Row 0702: a sharper composition gap — A at c = £4,000 (p′ = 20% = 2000 bp),
--     B at c = £9,000 (10% = 1000 bp); the heavier composition halves the rate.
--
-- magnitude_changes (refutation of Rodbertus)
--   Row 0703 money_value_change: capital £100, profit £20 (20%); gold falls by
--     half → reckoned £200/£40, rate still 20%. rate_unchanged = 1.
--   Row 0704 proportional_change: capital £3,000, profit £600 (20%); the capital
--     grows 2.5× at the same composition, rate still 20%. rate_unchanged = 1.
--   Row 0705 organic_change: capital £100, profit £20 (20%); the capital doubles
--     with the same profit (composition shifts) → rate falls to 10%. rate_unchanged = 0.

INSERT INTO composition_effects
    (id, s_a, v_a, c_a, s_b, v_b, c_b, profit_rate_a_bp, profit_rate_b_bp, rate_difference_bp, created_at)
VALUES
    ('5eed000000000000000701', 1000, 1000, 9000, 1000, 1000, 11000, 1000, 833, 167, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000702', 1000, 1000, 4000, 1000, 1000, 9000, 2000, 1000, 1000, '2026-05-01 00:00:00.000000');

INSERT INTO magnitude_changes
    (id, kind, original_capital, original_profit, factor, old_rate_bp, new_rate_bp, rate_unchanged, created_at)
VALUES
    ('5eed000000000000000703', 'money_value_change', 100, 20, 200, 2000, 2000, 1, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000704', 'proportional_change', 3000, 600, 250, 2000, 2000, 1, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000705', 'organic_change', 100, 20, 200, 2000, 1000, 0, '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM composition_effects WHERE id IN (
    '5eed000000000000000701',
    '5eed000000000000000702'
);
DELETE FROM magnitude_changes WHERE id IN (
    '5eed000000000000000703',
    '5eed000000000000000704',
    '5eed000000000000000705'
);
