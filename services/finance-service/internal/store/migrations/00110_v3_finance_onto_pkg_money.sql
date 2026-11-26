-- +goose Up
-- Issue #449 / ADR 0001 — finance-service moves onto pkg/money. Its own
-- internal/money package (£1 = 480 labour-minutes) is replaced by the canonical
-- MELT, £1 = 2,400 labour-minutes (Marx Ch. 17: 6s = a 12-hour day), a factor
-- of 5. Two kinds of seeded row carry a £-denominated labour-minute magnitude
-- and are rescaled here by ID (migrations are append-only, so 00016/00076/00109
-- stay intact); every UPDATE is guarded by the old value so it is idempotent.
--
-- Ch. 8 labour-power index, v·LabourMinutesPerPound·100/C (truncated):
--   …0801 Sphere A   (v=100, C=700): 6857  -> 34285
--   …0802 Sphere B   (v=600, C=700): 41142 -> 205714
--   …0803 European   (v=16,  C=100): 7680  -> 38400
--   …0804 Asian      (v=84,  C=100): 40320 -> 201600
--
-- Ch. 38 waterfall: the £10 surplus-profit IS the £10 differential rent, so the
-- rent and the £200 land-price it capitalises at 5% scale with the MELT:
--   …3801 annual rent  4800  -> 24000 LM
--   …3802 annual rent  4800  -> 24000 LM, capitalised price 96000 -> 480000 LM
--
-- Other labour-minute rent seeds (ground rents, building-site, monopoly) are not
-- tied to a £ magnitude by any identity and are left as recorded.

UPDATE production_spheres SET labour_power_index = 34285  WHERE id = '5eed000000000000000801' AND labour_power_index = 6857;
UPDATE production_spheres SET labour_power_index = 205714 WHERE id = '5eed000000000000000802' AND labour_power_index = 41142;
UPDATE production_spheres SET labour_power_index = 38400  WHERE id = '5eed000000000000000803' AND labour_power_index = 7680;
UPDATE production_spheres SET labour_power_index = 201600 WHERE id = '5eed000000000000000804' AND labour_power_index = 40320;

UPDATE differential_rents_ch38 SET annual_rent_labour_minutes = 24000 WHERE id = '5eed000000000000003801' AND annual_rent_labour_minutes = 4800;
UPDATE capitalised_rent_prices SET annual_rent_labour_minutes = 24000, capitalised_price_labour_minutes = 480000 WHERE id = '5eed000000000000003802' AND annual_rent_labour_minutes = 4800 AND capitalised_price_labour_minutes = 96000;

-- +goose Down
UPDATE capitalised_rent_prices SET annual_rent_labour_minutes = 4800, capitalised_price_labour_minutes = 96000 WHERE id = '5eed000000000000003802' AND annual_rent_labour_minutes = 24000 AND capitalised_price_labour_minutes = 480000;
UPDATE differential_rents_ch38 SET annual_rent_labour_minutes = 4800 WHERE id = '5eed000000000000003801' AND annual_rent_labour_minutes = 24000;

UPDATE production_spheres SET labour_power_index = 40320  WHERE id = '5eed000000000000000804' AND labour_power_index = 201600;
UPDATE production_spheres SET labour_power_index = 7680   WHERE id = '5eed000000000000000803' AND labour_power_index = 38400;
UPDATE production_spheres SET labour_power_index = 41142  WHERE id = '5eed000000000000000802' AND labour_power_index = 205714;
UPDATE production_spheres SET labour_power_index = 6857   WHERE id = '5eed000000000000000801' AND labour_power_index = 34285;
