-- +goose Up
-- Vol. III Ch. 9 seed — Marx's five worked production spheres (all C=100,
-- s′=100%) so the general-rate and prices-of-production panels come up
-- populated on a fresh MySQL volume.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=09 (chapter) and
-- XX is the row sequence. These 22-char IDs never collide with the 24-char hex
-- NewGeneralProfitRateID() / NewPriceOfProductionID() output.
--
-- General rate (seed id …0901):
--   Five spheres, each C=100, s′=100%. v: I=20, II=30, III=40, IV=15, V=5.
--   ΣS = 20+30+40+15+5 = 110; ΣC = 500.
--   p̄′ = round_half_up(110·10000 / 500) = (1100000 + 250) / 500 = 2200.
--   AverageVariablePercent = 110·100 / 500 = 22.
--   spheres_json stored as empty array (the five spheres are in prices_of_production).
--
-- Prices of production (seed ids …0902–…0906), k=100, p̄′=2200, avg profit=22,
-- price=122 for all five spheres:
--   Sphere I   (80c+20v): value=120, deviation=+2,  composition=higher
--   Sphere II  (70c+30v): value=130, deviation=−8,  composition=lower
--   Sphere III (60c+40v): value=140, deviation=−18, composition=lower
--   Sphere IV  (85c+15v): value=115, deviation=+7,  composition=higher
--   Sphere V   (95c+5v):  value=105, deviation=+17, composition=higher

INSERT INTO general_profit_rates
    (id, rate, sum_surplus_values, sum_total_capitals, average_variable_percent, spheres_json, created_at)
VALUES
    ('5eed000000000000000901', 2200, 110, 500, 22, '[]', '2026-05-01 00:00:00.000000');

INSERT INTO prices_of_production
    (id, sphere_name, cost_price, general_rate, commodity_value, average_profit, price, deviation, composition, created_at)
VALUES
    ('5eed000000000000000902', 'Sphere I (80c+20v)',  100, 2200, 120, 22, 122,   2, 'higher', '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000903', 'Sphere II (70c+30v)', 100, 2200, 130, 22, 122,  -8, 'lower',  '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000904', 'Sphere III (60c+40v)',100, 2200, 140, 22, 122, -18, 'lower',  '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000905', 'Sphere IV (85c+15v)', 100, 2200, 115, 22, 122,   7, 'higher', '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000906', 'Sphere V (95c+5v)',   100, 2200, 105, 22, 122,  17, 'higher', '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM prices_of_production WHERE id IN (
    '5eed000000000000000902',
    '5eed000000000000000903',
    '5eed000000000000000904',
    '5eed000000000000000905',
    '5eed000000000000000906'
);
DELETE FROM general_profit_rates WHERE id IN (
    '5eed000000000000000901'
);
