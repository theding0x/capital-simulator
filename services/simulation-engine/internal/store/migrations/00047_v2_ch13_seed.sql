-- +goose Up
-- Vol. II Ch. 13 — The Time of Production: seed data
-- Seed IDs follow pattern 5eed000000000000001300..1399 (ch13 = 13).

-- Working periods for the six natural-process fixtures.
-- These are connected working periods representing natural-process-dominated production.
INSERT INTO working_periods (id, industrial_capital_id, commodity_id, working_day_count, mode, working_day_length_minutes, created_at)
VALUES
  ('5eed000000000000001301', 'corn-farm-1871-ic',   'corn',      280, 'connected', 600, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001302', 'winery-1871-ic',      'wine',      365, 'connected', 600, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001303', 'cattle-farm-1871-ic', 'cattle',    730, 'connected', 600, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001304', 'timber-1871-ic',      'timber',  18250, 'connected', 600, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001305', 'porcelain-1871-ic',   'porcelain',  30, 'connected', 600, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001306', 'fallow-farm-1871-ic', 'wheat',     180, 'connected', 600, '1871-01-01 00:00:00.000000');

-- Production-labour gaps (nanoseconds; 1 day = 86400000000000 ns).
-- Corn:      10 labour-days /  280 production-days.
-- Wine:      10 labour-days /  365 production-days.
-- Cattle:     5 labour-days /  730 production-days.
-- Timber:    30 labour-days / 18250 production-days (50-year rotation).
-- Porcelain:  3 labour-days /   30 production-days.
-- Fallow:     5 labour-days /  180 production-days.
INSERT INTO production_labour_gaps (id, working_period_id, production_time_nanos, labour_time_nanos, gap_nanos, created_at)
VALUES
  ('5eed000000000000001310', '5eed000000000000001301',  24192000000000000,   864000000000000,  23328000000000000, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001311', '5eed000000000000001302',  31536000000000000,   864000000000000,  30672000000000000, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001312', '5eed000000000000001303',  63072000000000000,   432000000000000,  62640000000000000, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001313', '5eed000000000000001304',1576800000000000000,  2592000000000000,1574208000000000000, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001314', '5eed000000000000001305',   2592000000000000,   259200000000000,   2332800000000000, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001315', '5eed000000000000001306',  15552000000000000,   432000000000000,  15120000000000000, '1871-01-01 00:00:00.000000');

-- Natural-process activations (one representative activation per working period).
INSERT INTO natural_process_activations (id, working_period_id, kind, started_at, ended_at, requires_labour_maintenance)
VALUES
  ('5eed000000000000001320', '5eed000000000000001301', 'agricultural_growth', '1871-04-01 00:00:00.000000', '1871-09-30 00:00:00.000000', 0),
  ('5eed000000000000001321', '5eed000000000000001302', 'fermentation',        '1871-10-01 00:00:00.000000', '1872-09-30 00:00:00.000000', 0),
  ('5eed000000000000001322', '5eed000000000000001303', 'animal_maturation',   '1871-01-01 00:00:00.000000', NULL,                          1),
  ('5eed000000000000001323', '5eed000000000000001304', 'forest_growth',       '1871-01-01 00:00:00.000000', NULL,                          0),
  ('5eed000000000000001324', '5eed000000000000001305', 'chemical',            '1871-01-01 00:00:00.000000', '1871-01-31 00:00:00.000000',  0),
  ('5eed000000000000001325', '5eed000000000000001306', 'seasonal_rest',       '1871-11-01 00:00:00.000000', '1872-04-30 00:00:00.000000',  0),
  ('5eed000000000000001326', '5eed000000000000001303', 'animal_maturation',   '1872-01-01 00:00:00.000000', NULL,                          1);

-- Natural subjects (one per working period).
INSERT INTO natural_subjects (id, working_period_id, commodity_id, description, is_living, created_at)
VALUES
  ('5eed000000000000001330', '5eed000000000000001301', 'corn',      'Corn crop; seed planted April, harvested late September. Labour required only for sowing and harvest.', 0, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001331', '5eed000000000000001302', 'wine',      'Grape must fermenting in casks; inert subject — no labour required during fermentation.', 0, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001332', '5eed000000000000001303', 'cattle',    'Cattle on Bakewell farm; living subject requiring periodic feeding and care.', 1, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001333', '5eed000000000000001304', 'timber',    'Commercial timber stand (50-year rotation); living subject requiring minimal periodic labour.', 1, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001334', '5eed000000000000001305', 'porcelain', 'Porcelain in kiln firing; inert subject — no labour during the chemical transformation.', 0, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001335', '5eed000000000000001306', 'wheat',     'Winter wheat under fallow; land lies idle between harvest and next sowing.', 0, '1871-01-01 00:00:00.000000');

-- Industry benchmarks (five Marx-faithful examples).
-- ratio_basis_points = 10000 * labour_days / production_days (integer division).
INSERT INTO natural_process_industries (id, industry_name, typical_production_days_count, typical_labour_days_count, ratio_basis_points, created_at)
VALUES
  ('5eed000000000000001340', 'Agriculture (corn)',           280,    10,  357, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001341', 'Forestry (100-year timber)', 36500,    30,    8, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001342', 'Leather tanning',              365,   120, 3287, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001343', 'Brewing (ale)',                 21,    14, 6666, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001344', 'Porcelain (hard-paste)',        30,     3, 1000, '1871-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM natural_process_industries  WHERE id IN ('5eed000000000000001340','5eed000000000000001341','5eed000000000000001342','5eed000000000000001343','5eed000000000000001344');
DELETE FROM natural_subjects            WHERE id IN ('5eed000000000000001330','5eed000000000000001331','5eed000000000000001332','5eed000000000000001333','5eed000000000000001334','5eed000000000000001335');
DELETE FROM natural_process_activations WHERE id IN ('5eed000000000000001320','5eed000000000000001321','5eed000000000000001322','5eed000000000000001323','5eed000000000000001324','5eed000000000000001325','5eed000000000000001326');
DELETE FROM production_labour_gaps      WHERE id IN ('5eed000000000000001310','5eed000000000000001311','5eed000000000000001312','5eed000000000000001313','5eed000000000000001314','5eed000000000000001315');
DELETE FROM working_periods             WHERE id IN ('5eed000000000000001301','5eed000000000000001302','5eed000000000000001303','5eed000000000000001304','5eed000000000000001305','5eed000000000000001306');
