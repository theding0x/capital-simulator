-- +goose Up
-- Marx Ch. 41 — DR II constant price. Case II (proportional) Soil D: total rent
-- 1800 (£18), rate 36000 bp. Case III (decreasing rate): total 1350, rate 27000.
-- Intensive (1 acre £18/acre) vs extensive (2 acres £9/acre), same £18 total.
INSERT INTO `dr2_outcomes`
    (`id`, `parcel_id`, `productivity_case`, `additional_capital_bp`, `new_rent_per_acre_bp`, `total_rental_bp`, `rate_of_surplus_profit_bp`, `created_at`)
VALUES
    ('5eed000000000000004100', '5eed000000000000003700', 2, 250, 1800, 1800, 36000, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000004101', '5eed000000000000003700', 3, 250, 1350, 1350, 27000, '1867-01-01 00:00:00.000000');

INSERT INTO `intensive_extensive_comparisons`
    (`id`, `intensive_rent_per_acre_bp`, `extensive_rent_per_acre_bp`, `total_capital_bp`, `total_rental_same_bp`, `created_at`)
VALUES
    ('5eed000000000000004102', 1800, 900, 500, 1800, '1867-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM `dr2_outcomes` WHERE `id` IN ('5eed000000000000004100', '5eed000000000000004101');
DELETE FROM `intensive_extensive_comparisons` WHERE `id` = '5eed000000000000004102';
