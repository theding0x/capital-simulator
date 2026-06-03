-- +goose Up
-- Ch. 48 conservation fix (#377): the original trinity seed (00096) double-counted
-- the ground-rent. Capital's actual surplus share was set to the full average
-- profit (20), and the formula's total_surplus_value_bp to 20 profit + 3 rent = 23.
-- Marx's point is that rent is *deducted from* the single social surplus-value,
-- not added on top of it: profit of enterprise 12 + interest 5 = capital's
-- retained share 17, plus ground-rent 3 = one surplus-value s = 20. Apparent
-- revenue 103 then overshoots newly-created value v + s = 100 by exactly the rent
-- 3 — and that overshoot is the fetish, surfaced rather than masked. Migrations
-- are append-only, so we correct the seeded rows here instead of editing 00096.
UPDATE `revenue_streams` SET `actual_source_bp` = 17 WHERE `id` = '5eed000000000000004801';
UPDATE `trinity_formulas` SET `total_surplus_value_bp` = 20 WHERE `id` = '5eed000000000000004800';

-- +goose Down
UPDATE `trinity_formulas` SET `total_surplus_value_bp` = 23 WHERE `id` = '5eed000000000000004800';
UPDATE `revenue_streams` SET `actual_source_bp` = 20 WHERE `id` = '5eed000000000000004801';
