-- +goose Up
-- Ch. 44 — DR on worst soil. Prices £1=10000bp: £3=30000, £3.50=35000, gap £0.50=5000.
INSERT INTO `rent_on_worst_soils`
    (`id`, `parcel_id`, `mechanism`, `old_price_of_production_bp`, `new_price_of_production_bp`, `rent_per_acre_bp`, `created_at`)
VALUES
    ('5eed000000000000004400', '5eed000000000000003700', 1, 30000, 35000, 5000, '1894-01-01 00:00:00.000000'),
    ('5eed000000000000004401', '5eed000000000000003700', 3, 30000, 35000, 5000, '1894-01-01 00:00:00.000000');

INSERT INTO `soil_improvement_capitals`
    (`id`, `parcel_id`, `capital_invested_bp`, `productivity_gain_bp`, `amortisation_years`, `amortised_years`, `created_at`)
VALUES
    ('5eed000000000000004402', '5eed000000000000003700', 3000, 500, 15, 0, '1894-01-01 00:00:00.000000');

INSERT INTO `regulating_price_shifts`
    (`id`, `from_grade`, `to_grade`, `old_price_bp`, `new_price_bp`, `created_at`)
VALUES
    ('5eed000000000000004403', 1, 2, 30000, 35000, '1894-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM `regulating_price_shifts`   WHERE `id` = '5eed000000000000004403';
DELETE FROM `soil_improvement_capitals` WHERE `id` = '5eed000000000000004402';
DELETE FROM `rent_on_worst_soils`       WHERE `id` IN ('5eed000000000000004400', '5eed000000000000004401');
