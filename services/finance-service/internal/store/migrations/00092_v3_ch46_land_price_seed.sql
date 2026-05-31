-- +goose Up
-- Ch. 46: London ground lease £300/99yr; Scottish coal worst mine (0 rent);
-- Burgundy monopoly value 200 sale 400 rent 200; land price rent 10 @ 5%→2.5%, price 200→400.
INSERT INTO `building_site_rents` (`id`,`parcel_id`,`location_grade`,`annual_rent_labour_minutes`,`lease_term_years`,`created_at`) VALUES
    ('5eed000000000000004600', '5eed000000000000003700', 5, 300, 99, '1894-01-01 00:00:00.000000');
INSERT INTO `mining_rents` (`id`,`parcel_id`,`ore_grade`,`annual_rent_labour_minutes`,`surplus_profit_bp`,`created_at`) VALUES
    ('5eed000000000000004601', '5eed000000000000003700', 1, 0, 0, '1894-01-01 00:00:00.000000');
INSERT INTO `monopoly_price_rents` (`id`,`parcel_id`,`product_kind`,`value_labour_minutes`,`monopoly_sale_price_labour_minutes`,`monopoly_rent_labour_minutes`,`created_at`) VALUES
    ('5eed000000000000004602', '5eed000000000000003700', 'burgundy_grand_cru', 200, 400, 200, '1894-01-01 00:00:00.000000');
INSERT INTO `land_price_scenarios` (`id`,`kind`,`parcel_id`,`old_annual_rent_bp`,`new_annual_rent_bp`,`old_interest_rate_bp`,`new_interest_rate_bp`,`old_land_price_bp`,`new_land_price_bp`,`created_at`) VALUES
    ('5eed000000000000004603', 1, '5eed000000000000003700', 10, 10, 500, 250, 200, 400, '1894-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM `land_price_scenarios` WHERE `id` = '5eed000000000000004603';
DELETE FROM `monopoly_price_rents` WHERE `id` = '5eed000000000000004602';
DELETE FROM `mining_rents` WHERE `id` = '5eed000000000000004601';
DELETE FROM `building_site_rents` WHERE `id` = '5eed000000000000004600';
