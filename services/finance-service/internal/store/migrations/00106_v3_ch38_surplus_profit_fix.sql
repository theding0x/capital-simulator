-- +goose Up
-- Ch. 38 surplus-profit fix (#378): the original seed (00076) fed the waterfall
-- factory's individual COST-PRICE (£90 = 9000) into the individual price-of-
-- production field, so surplus-profit came out as general 11500 - 9000 = 2500
-- (£25). Marx computes surplus-profit as general price of production - individual
-- PRICE OF PRODUCTION (cost-price £90 + average profit £15 = £105 = 10500),
-- giving £10 — and that £10 IS the differential rent (4800 LM) and the £200
-- capitalised price (96000 LM) seeded alongside. Correct the individual PoP and
-- both surplus-profit fields so they no longer contradict the rent they equal.
-- Migrations are append-only, so we correct the seeded rows here, not in 00076.
UPDATE pop_surplus_profits
   SET individual_price_of_production_bp = 10500, surplus_profit_bp = 1000
 WHERE id = '5eed000000000000003803';
UPDATE differential_rents_ch38
   SET surplus_profit_bp = 1000
 WHERE id = '5eed000000000000003801';

-- +goose Down
UPDATE differential_rents_ch38
   SET surplus_profit_bp = 2500
 WHERE id = '5eed000000000000003801';
UPDATE pop_surplus_profits
   SET individual_price_of_production_bp = 9000, surplus_profit_bp = 2500
 WHERE id = '5eed000000000000003803';
