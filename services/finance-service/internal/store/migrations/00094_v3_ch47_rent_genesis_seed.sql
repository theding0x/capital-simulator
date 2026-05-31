-- +goose Up
-- Ch. 47 — Genesis: four rent forms (labour→product→money→capitalist), two transitions,
-- one small-peasant example (France: rent 2 + profit 1 + wages 3 = income 6, conflated).
INSERT INTO `rent_historical_forms` (`id`,`stage`,`name`,`coercion_kind`,`surplus_form`,`example_region`,`example_period`,`created_at`) VALUES
    ('5eed000000000000004700', 1, 'labour-rent',     'extra-economic',             'direct_labour',    'Normandy, England', '10th-13th century', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000004701', 2, 'product-rent',    'extra-economic (diminished)','product',          'Poland, Russia',    '14th-17th century', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000004702', 3, 'money-rent',      'contractual',                'money',            'England, France',   '15th-17th century', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000004703', 4, 'capitalist-rent', 'economic',                   'profit_deduction', 'England',           '18th century onward','1894-01-01 00:00:00.000000');
INSERT INTO `historical_rent_transitions` (`id`,`from_stage`,`to_stage`,`preconditions`,`driving_force`,`created_at`) VALUES
    ('5eed000000000000004704', 1, 2, 'rising productivity on the peasant plot', 'objectification of surplus-labour in the product', '1894-01-01 00:00:00.000000'),
    ('5eed000000000000004705', 3, 4, 'expropriation of the producer; wage-labour in agriculture; capitalist mode generalised', 'capital accumulation, market integration', '1894-01-01 00:00:00.000000');
INSERT INTO `small_peasant_productions` (`id`,`parcel_id`,`total_income_bp`,`rent_component_bp`,`profit_component_bp`,`wages_component_bp`,`is_conflated`,`created_at`) VALUES
    ('5eed000000000000004706', '5eed000000000000003700', 6, 2, 1, 3, 1, '1894-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM `small_peasant_productions` WHERE `id` = '5eed000000000000004706';
DELETE FROM `historical_rent_transitions` WHERE `id` IN ('5eed000000000000004704','5eed000000000000004705');
DELETE FROM `rent_historical_forms` WHERE `id` IN ('5eed000000000000004700','5eed000000000000004701','5eed000000000000004702','5eed000000000000004703');
