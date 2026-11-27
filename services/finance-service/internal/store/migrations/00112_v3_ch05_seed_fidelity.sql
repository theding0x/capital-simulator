-- +goose Up
-- Vol. III Ch. 5 seed fidelity (issue #450, v3-ch05 item). 00010 seeded three
-- economies. Only one of them carries Marx's figures, and those are Ch. 7's, not
-- Ch. 5 §I's. This restates the seed in Ch. 5's own figures.
--
-- Tags: M = Marx's figure (stated or quoted by him), D = derived from his, I =
-- illustrative (Marx gives no number). Money scale: the Vol. III Part I convention, one stored unit = £1 as
-- the panel prints it (LabourMinutes-typed columns, NOT the pkg/money MELT of
-- £1 = 2,400 LM). The rates are basis points and independent of the scale. Marx's
-- £4 4s. is not a whole pound, so only the £44 is stored here. Both prices are
-- held exactly, at the MELT, in profit's economy_test.go.
--
--   0501 waste_reduction. Was 40 lb of waste cotton at £5 a lb = 200: invented,
--     and about 250 times Marx's rag price per lb. Now the woollen rags
--     Marx cites from a factory report, about £4 4s. a ton thirty years ago and
--     £44 a ton now (M, Ch. 5 §IV). 10 tons
--     reclaimed (I) at £44 = 440 (D). Base 2,000c + 1,000v + 1,000s unchanged (I).
--     p' 3333 bp -> 3906 bp, gain 573.
--   0502 shared_conditions. Was a boiler-house shared by 12 factories, saving 100:
--     in Marx nowhere. Now Ch. 5 §III: a cotton spinner Marx quotes drives 29,000
--     doubling spindles by a single shaft, saving 60 per cent of the shafting, in
--     some cases 80 (M). Shafting other firms need for the same spindles 500 (I), here 200 (D),
--     saving 300 (D). Base 1,200c + 600v + 600s unchanged (I). p' 3333 -> 4000 bp,
--     gain 667.
--   0503 raw_material_saving, left as it is. These are Ch. 7's figures (Supplementary
--     Remarks, reason 2), not Ch. 5 §I: s = v = £1,000, c = £9,000 (p' 10%) or
--     £11,000 (8 1/3%) (M). The economy of 2,000 is their difference (D). It is
--     the same pair as composition_effects row 0701 (00014).
--   0504 workday_extension, new. Ch. 5 §I: the fixed capital is the same whether
--     the labour-process serves 12 or 16 hours, and at 18 instead of 12 hours its
--     value is reproduced in two thirds of the time (M). 3,600 of fixed capital
--     (I, and 3,600 = 12 x 300 = 18 x 200), taken as all of c (I, the upper bound),
--     v = s = 1,200 (I). The charge on the same 12 working hours falls from 3,600
--     to 2,400, a saving of 1,200 (D). p' 2500 -> 3333 bp, gain 833.
--
-- Rows are UPDATEd in place, each guarded by the whole pre-image it changes. The
-- Down restores 00010's values and deletes only the id it inserted.

UPDATE economy_analyses SET saving = 440, new_profit_rate_bp = 3906, profit_rate_gain_bp = 573 WHERE id = '5eed000000000000000501' AND saving = 200 AND new_profit_rate_bp = 3571 AND profit_rate_gain_bp = 238;
UPDATE economy_analyses SET saving = 300, new_profit_rate_bp = 4000, profit_rate_gain_bp = 667 WHERE id = '5eed000000000000000502' AND saving = 100 AND new_profit_rate_bp = 3529 AND profit_rate_gain_bp = 196;

INSERT INTO economy_analyses
    (id, kind, constant_capital, variable_capital, surplus_value, saving,
     old_profit_rate_bp, new_profit_rate_bp, profit_rate_gain_bp, created_at)
VALUES
    ('5eed000000000000000504', 'workday_extension', 3600, 1200, 1200, 1200,
     2500, 3333, 833, '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM economy_analyses WHERE id = '5eed000000000000000504';
UPDATE economy_analyses SET saving = 100, new_profit_rate_bp = 3529, profit_rate_gain_bp = 196 WHERE id = '5eed000000000000000502' AND saving = 300 AND new_profit_rate_bp = 4000 AND profit_rate_gain_bp = 667;
UPDATE economy_analyses SET saving = 200, new_profit_rate_bp = 3571, profit_rate_gain_bp = 238 WHERE id = '5eed000000000000000501' AND saving = 440 AND new_profit_rate_bp = 3906 AND profit_rate_gain_bp = 573;
