-- +goose Up
-- Issue #411 — finance-service had labourMinutesPerPound defined twice with
-- divergent values: avgprofit used 3600 (£1 = a 60-hour week), rent used 480
-- (£1 = an 8-hour working day). The two are collapsed onto one canonical scale
-- in the new internal/money package, and 480 is chosen because the whole rent
-- sub-system (Ch. 37–46 LabourMinute rents and land-prices) is already built on
-- it; avgprofit's labour-power index was the lone 3600 outlier.
--
-- avgprofit.labourPowerIndex = v · LabourMinutesPerPound · 100 / C now divides
-- by 480 instead of 3600, so the Ch. 8 seed's hard-coded indices (computed once
-- at 3600 in 00016) are rescaled here by 480/3600 = 2/15. Migrations are
-- append-only, so 00016 stays intact and the derived column is corrected by ID:
--
--   …0801 Sphere A   (v=100, C=700): 100·480·100/700 = 6857  (was 51428)
--   …0802 Sphere B   (v=600, C=700): 600·480·100/700 = 41142 (was 308571)
--   …0803 European   (v=16,  C=100): 16·480·100/100  = 7680  (was 57600)
--   …0804 Asian      (v=84,  C=100): 84·480·100/100  = 40320 (was 302400)
--
-- The index is a per-£100 comparison magnitude; rescaling preserves every
-- cross-sphere ordering the chapter exhibits (Asian still sets the most living
-- labour in motion).

UPDATE production_spheres SET labour_power_index = 6857  WHERE id = '5eed000000000000000801' AND labour_power_index = 51428;
UPDATE production_spheres SET labour_power_index = 41142 WHERE id = '5eed000000000000000802' AND labour_power_index = 308571;
UPDATE production_spheres SET labour_power_index = 7680  WHERE id = '5eed000000000000000803' AND labour_power_index = 57600;
UPDATE production_spheres SET labour_power_index = 40320 WHERE id = '5eed000000000000000804' AND labour_power_index = 302400;

-- +goose Down
-- Restore the 3600-scale indices for the four seeded spheres by ID.
UPDATE production_spheres SET labour_power_index = 51428  WHERE id = '5eed000000000000000801' AND labour_power_index = 6857;
UPDATE production_spheres SET labour_power_index = 308571 WHERE id = '5eed000000000000000802' AND labour_power_index = 41142;
UPDATE production_spheres SET labour_power_index = 57600  WHERE id = '5eed000000000000000803' AND labour_power_index = 7680;
UPDATE production_spheres SET labour_power_index = 302400 WHERE id = '5eed000000000000000804' AND labour_power_index = 40320;
