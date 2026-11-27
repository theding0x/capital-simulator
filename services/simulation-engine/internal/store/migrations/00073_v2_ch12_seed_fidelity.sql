-- +goose Up
-- Vol. II Ch. 12 seed fidelity (issue #450, v2-ch12 item). 00045 is untouched; this
-- restates the Bakewell and London-builder fixtures and adds the missing reference
-- working periods. Tags: M = Marx's own figure, D = derived from his, I = illustrative.
--
-- Bakewell (M). Marx's example is sheep: formerly "not fit for the butcher until four or
-- five years old"; on the Bakewell system "fattened when only one year old". 00045 had it
-- as cattle, 1095 to 547 (3 y to 1.5 y), hung on the 84-day locomotive period. Now: a
-- sheep period of 1460 days (4 y, the low end of "four or five", D), shortening 1460 to
-- 365, no extra fixed capital or labour. Constraint 1214 'bakewell-cattle' 547 becomes
-- 'sheep' 365, the Bakewell floor. Lavergne's three-for-one is a throughput ratio and is
-- not fixtured. The cattle floor 1095 (3 y) is the spec's figure, not Marx's ("several
-- years", Good: "four- and five-year-old oxen"); the value stays because Ch. 13 spans
-- (00071) are pinned to it, and only its "pre-Bakewell" reason is replaced.
--
-- London builder (D, I). Marx names no sums, only "build from 100 to 200 houses ... an
-- enterprise which exceeds his resources twenty to fifty times. The funds are procured
-- through mortgaging" (Select Committee on Bank Acts 1857, Evidence 5413-18, 5435-36).
-- 00045 stored 2400 / 48000 pence, commented as GBP 100 / GBP 2,000, but at 240 pence to
-- the pound that is GBP 10 / GBP 200. Now 24000 / 480000 (GBP 100 / GBP 2,000): borrowed is
-- 20 x own, so the enterprise is 21 x own; both readings of "exceeds his resources twenty
-- times" (enterprise = 20 x own, or excess = 20 x own) lie inside Marx's 20 to 50. Only the
-- ratio is his; the pound amounts are illustrative. The financing is moved from the
-- locomotive works to a new 'house' period of 168 days (24 weeks x 7, illustrative: Marx
-- says only that a house is "built in less time than a large factory"). The mortgage
-- provider stays the Vol. III placeholder; Marx names no bank.
--
-- New reference periods (all connected). Each row carries its own tag. Ten-hour day, 600
-- min: M ("working-days of equal length, say, of ten hours each"). Man-of-war 365: D ("one
-- year or more", his lower bound). Carpet 7: M rounded ("about a week"). Country road 84:
-- I ("a few months", read as the locomotive's three months, 12 x 7). Railway 730: I ("a job
-- of years", at least two). Gobelin 1095: I ("takes years"). Sheep 1460: D (low end of
-- "four or five years old"). House 168: I (Marx gives no number). The ordering of the road,
-- railway, Gobelin and man-of-war durations among themselves is the seed's choice, not his.
--
-- Grain, big cattle and timber are NOT repeated here: Ch. 13 (00047/00071) already seeds
-- corn, cattle and timber working periods, and 00071 separated their labour-days from the
-- production span. Those Ch. 13 rows now hold labour-days (10, 5 and 30), not Ch. 12
-- working-period spans, and the spec's second timber row (100 years) is not seeded.
--
-- Left as 00045 seeded and declared illustrative (Marx gives no figure): the spinning-mill
-- and locomotive interruptions (0/0 and 240/120 pence) and their dates, the machinery
-- shortening 84 to 42 with 5000 pence, the parallel-labour shortening 84 to 28 with 50
-- labourers (it hangs on the locomotive, 1232, although Marx's example is the railway),
-- the timber floor 3650 (00045 says "10+ years", Marx says "from twelve to one hundred
-- years"), the merino-wool floor 365 (not in Marx), and every created_at. His own:
-- ten-hour day, the spinner's week, the locomotive's three months (12 x the week).
--
-- Note: 00045's header ("multiplier=1" for the 7-day mill) is stale; Multiplier() is the
-- working-day count in both modes.

INSERT INTO working_periods (id, industrial_capital_id, commodity_id, working_day_count, mode, working_day_length_minutes, created_at)
VALUES
  ('5eed000000000000001203', 'armoured-ship-yard-ic',        'man-of-war',       365,  'connected', 600, '1871-01-02 00:00:00.000000'),
  ('5eed000000000000001204', 'country-road-works-ic',        'country-road',      84,  'connected', 600, '1871-01-02 00:00:00.000000'),
  ('5eed000000000000001205', 'railway-contractor-ic',        'railway',          730,  'connected', 600, '1871-01-02 00:00:00.000000'),
  ('5eed000000000000001206', 'carpet-works-ic',              'carpet',             7,  'connected', 600, '1871-01-02 00:00:00.000000'),
  ('5eed000000000000001207', 'tapestry-works-ic',            'gobelin-tapestry', 1095, 'connected', 600, '1871-01-02 00:00:00.000000'),
  ('5eed000000000000001208', 'dishley-grange-ic',            'sheep',           1460,  'connected', 600, '1871-01-02 00:00:00.000000'),
  ('5eed000000000000001209', 'london-speculative-builder-ic','house',            168,  'connected', 600, '1871-01-02 00:00:00.000000');

UPDATE working_period_shortenings SET working_period_id = '5eed000000000000001208', old_working_day_count = 1460, new_working_day_count = 365 WHERE id = '5eed000000000000001230';

UPDATE natural_working_period_constraints SET commodity_id = 'sheep', min_working_days = 365, reason = 'Bakewell sheep may be fattened "when only one year old" (Marx, Ch. 12): the natural floor after breeding is one year (365 days), against four or five years before Bakewell' WHERE id = '5eed000000000000001214';

UPDATE natural_working_period_constraints SET reason = 'livestock maturation - Marx: "several years to raise big cattle" (Good: four- and five-year-old oxen). 1095 days (three years) is an illustrative figure, not his' WHERE id = '5eed000000000000001211';

UPDATE working_period_credit_financings SET working_period_id = '5eed000000000000001209', own_capital_pence = 24000, borrowed_capital_pence = 480000 WHERE id = '5eed000000000000001240';

-- +goose Down
-- Restores 00045 exactly: financing and the Bakewell shortening back on the locomotive
-- period, the old constraint rows, then removes only the seven added periods.
UPDATE working_period_credit_financings SET working_period_id = '5eed000000000000001202', own_capital_pence = 2400, borrowed_capital_pence = 48000 WHERE id = '5eed000000000000001240';

UPDATE natural_working_period_constraints SET reason = 'livestock maturation — cattle require approximately three years to reach market weight (pre-Bakewell)' WHERE id = '5eed000000000000001211';

UPDATE natural_working_period_constraints SET commodity_id = 'bakewell-cattle', min_working_days = 547, reason = 'Bakewell selective breeding reduces cattle maturation from 3 years to ~1.5 years' WHERE id = '5eed000000000000001214';

UPDATE working_period_shortenings SET working_period_id = '5eed000000000000001202', old_working_day_count = 1095, new_working_day_count = 547 WHERE id = '5eed000000000000001230';

DELETE FROM working_periods WHERE id IN ('5eed000000000000001203','5eed000000000000001204','5eed000000000000001205','5eed000000000000001206','5eed000000000000001207','5eed000000000000001208','5eed000000000000001209');
