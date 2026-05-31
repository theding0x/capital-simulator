-- +goose Up
-- Engels Ch. 43 Table XI — five soils A–E at 6s/bushel; base=0 increment=12 grades=5
-- → rents [0,12,24,36,48] sum=120s. Capital 50s each (total 250s); output 70 bushels.
INSERT INTO `engels_table_entries`
    (`id`, `table_number`, `soil_grade_label`, `output_bushels`, `rent_shillings`, `capital_shillings`, `created_at`)
VALUES
    ('5eed000000000000004300', 11, 'A', 10, 0,  50, '1894-01-01 00:00:00.000000'),
    ('5eed000000000000004301', 11, 'B', 12, 12, 50, '1894-01-01 00:00:00.000000'),
    ('5eed000000000000004302', 11, 'C', 14, 24, 50, '1894-01-01 00:00:00.000000'),
    ('5eed000000000000004303', 11, 'D', 16, 36, 50, '1894-01-01 00:00:00.000000'),
    ('5eed000000000000004304', 11, 'E', 18, 48, 50, '1894-01-01 00:00:00.000000');

INSERT INTO `rent_sequences`
    (`id`, `base_rent`, `rent_increment`, `num_grades`, `created_at`)
VALUES
    ('5eed000000000000004305', 0, 12, 5, '1894-01-01 00:00:00.000000');

INSERT INTO `rising_price_dr2_outcomes`
    (`id`, `rising_variant`, `total_capital_shillings`, `total_rent_shillings`, `total_output_bushels`, `price_per_bushel_shillings`, `created_at`)
VALUES
    ('5eed000000000000004306', 3, 250, 120, 70, 6, '1894-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM `rising_price_dr2_outcomes` WHERE `id` = '5eed000000000000004306';
DELETE FROM `rent_sequences`            WHERE `id` = '5eed000000000000004305';
DELETE FROM `engels_table_entries`      WHERE `id` IN ('5eed000000000000004300','5eed000000000000004301','5eed000000000000004302','5eed000000000000004303','5eed000000000000004304');
