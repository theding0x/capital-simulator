-- +goose Up
-- Ch. 45 — Marx core: C=60,V=40,S=40, value=140, price=120, gap=20, absolute rent=20.
INSERT INTO `agricultural_compositions` (`id`,`constant_capital_bp`,`variable_capital_bp`,`composition_ratio_bp`,`social_average_ratio_bp`,`is_below_average`,`created_at`) VALUES
    ('5eed000000000000004500', 60, 40, 15000, 40000, 1, '1894-01-01 00:00:00.000000');
INSERT INTO `value_price_gaps` (`id`,`product_id`,`value_labour_minutes`,`price_of_production_labour_minutes`,`gap_labour_minutes`,`created_at`) VALUES
    ('5eed000000000000004501', 'agricultural-output', 140, 120, 20, '1894-01-01 00:00:00.000000');
INSERT INTO `absolute_rents` (`id`,`parcel_id`,`value_price_gap_id`,`surplus_value_bp`,`average_profit_bp`,`absolute_rent_bp`,`created_at`) VALUES
    ('5eed000000000000004502', '5eed000000000000003700', '5eed000000000000004501', 40, 20, 20, '1894-01-01 00:00:00.000000');
INSERT INTO `absolute_rent_limits` (`id`,`max_rent_bp`,`actual_rent_bp`,`is_at_limit`,`created_at`) VALUES
    ('5eed000000000000004503', 20, 20, 1, '1894-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM `absolute_rent_limits`      WHERE `id` = '5eed000000000000004503';
DELETE FROM `absolute_rents`            WHERE `id` = '5eed000000000000004502';
DELETE FROM `value_price_gaps`          WHERE `id` = '5eed000000000000004501';
DELETE FROM `agricultural_compositions` WHERE `id` = '5eed000000000000004500';
