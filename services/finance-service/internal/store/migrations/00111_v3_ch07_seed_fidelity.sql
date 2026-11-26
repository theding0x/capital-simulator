-- +goose Up
-- Vol. III Ch. 7 seed fidelity (issue #450, v3-ch07 item). 00014 seeded two rows
-- that trace to nothing in Marx's chapter. The chapter gives exactly one
-- composition comparison (row 0701: s = v = £1,000, A carrying c = £9,000 and B
-- c = £11,000), and states the proportional case of the Rodbertus refutation as
-- "C, or nC, or C/n, e.g., 1,000, or 2,000, or 500, and the rate of profit being
-- 20%, the profit = 200 in the first case, = 400 in the second, and = 100 in the
-- third".
--
--   * Row 0702 (A c = £4,000 vs B c = £9,000, "a sharper composition gap") is not
--     in the text: removed from composition_effects.
--   * Row 0704 was an invented £3,000 / £600 growing 2.5x: restated as Marx's
--     £1,000 / £200 doubled (nC = 2,000, profit 400). Its derived columns (old and
--     new rate 2000 bp, rate_unchanged = 1) stay correct, since p/C is 20% either way.
--   * Row 0706 is new: the same £1,000 / £200 halved (C/n = 500, profit 100).
--
-- Rows 0701, 0703 (gold halves: £100/£20 reckoned as £200/£40) and 0705 are left
-- alone. Seed IDs stay in the 5eed...07XX space; the Down restores 00014's values.

DELETE FROM composition_effects WHERE id = '5eed000000000000000702';

UPDATE magnitude_changes SET original_capital = 1000, original_profit = 200, factor = 200 WHERE id = '5eed000000000000000704';

INSERT INTO magnitude_changes
    (id, kind, original_capital, original_profit, factor, old_rate_bp, new_rate_bp, rate_unchanged, created_at)
VALUES
    ('5eed000000000000000706', 'proportional_change', 1000, 200, 50, 2000, 2000, 1, '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM magnitude_changes WHERE id = '5eed000000000000000706';

UPDATE magnitude_changes SET original_capital = 3000, original_profit = 600, factor = 250 WHERE id = '5eed000000000000000704';

INSERT INTO composition_effects
    (id, s_a, v_a, c_a, s_b, v_b, c_b, profit_rate_a_bp, profit_rate_b_bp, rate_difference_bp, created_at)
VALUES
    ('5eed000000000000000702', 1000, 1000, 4000, 1000, 1000, 9000, 2000, 1000, 1000, '2026-05-01 00:00:00.000000');
