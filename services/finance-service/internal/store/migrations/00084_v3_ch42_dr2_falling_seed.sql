-- +goose Up
-- Marx Ch. 42 — DR II Falling Price. Two A→B exclusions illustrate the
-- compensation law. Table IV: grain 250 @ 540p → 1350 (£13.50).
-- Table IVd: grain 500 @ 360p → 1800 (£18, = original Table I rental).
INSERT INTO `soil_exclusion_events`
    (`id`, `excluded_grade`, `new_regulating_grade`, `new_price_of_production_bp`, `created_at`)
VALUES
    ('5eed000000000000004200', 1, 2, 540, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000004205', 1, 2, 360, '1867-01-01 00:00:00.000000');

INSERT INTO `falling_price_dr2_outcomes`
    (`id`, `falling_variant`, `exclusion_event_id`, `total_money_rental_bp`, `total_grain_rent_qtrs`, `total_output_quarters`, `total_capital_bp`, `created_at`)
VALUES
    ('5eed000000000000004201', 1, '5eed000000000000004200', 1350, 250, 3600, 1500, '1867-01-01 00:00:00.000000'),
    ('5eed000000000000004202', 3, '5eed000000000000004205', 1800, 500, 3600, 1500, '1867-01-01 00:00:00.000000');

INSERT INTO `normal_capital_per_acre`
    (`id`, `period_label`, `capital_per_acre_bp`, `created_at`)
VALUES
    ('5eed000000000000004203', 'England pre-1848',  800,  '1867-01-01 00:00:00.000000'),
    ('5eed000000000000004204', 'England post-1848', 1200, '1867-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM `falling_price_dr2_outcomes` WHERE `id` IN ('5eed000000000000004201', '5eed000000000000004202');
DELETE FROM `soil_exclusion_events`      WHERE `id` IN ('5eed000000000000004200', '5eed000000000000004205');
DELETE FROM `normal_capital_per_acre`    WHERE `id` IN ('5eed000000000000004203', '5eed000000000000004204');
