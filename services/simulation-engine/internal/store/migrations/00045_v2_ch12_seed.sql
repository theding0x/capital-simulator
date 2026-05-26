-- +goose Up
-- Vol. II Ch. 12 — The Working Period: seed data
-- Seed IDs follow pattern 5eed000000000000001200..1299 (ch12 = 12).

-- Working periods
-- SpinningMill1871: daily discrete production (yarn sold each day, multiplier=1).
-- LocomotiveWorks:  84-day connected production (nothing sold until day 84, multiplier=84).
INSERT INTO working_periods (id, industrial_capital_id, commodity_id, working_day_count, mode, working_day_length_minutes, created_at)
VALUES
  ('5eed000000000000001201', 'spinning-mill-1871-ic', 'yarn',       7,  'discrete',  600, '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001202', 'locomotive-works-ic',   'locomotive', 84, 'connected', 600, '1871-01-01 00:00:00.000000');

-- Natural constraints for commodities with inviolable working-period floors.
INSERT INTO natural_working_period_constraints (id, commodity_id, min_working_days, reason, created_at)
VALUES
  ('5eed000000000000001210', 'corn',           365,  'annual crop — grain cannot be harvested before one full growing season', '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001211', 'cattle',         1095, 'livestock maturation — cattle require approximately three years to reach market weight (pre-Bakewell)', '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001212', 'timber',         3650, 'forest rotation — commercial timber requires 10+ years to reach harvestable size', '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001213', 'merino-wool',    365,  'annual shearing cycle — wool clip requires one full year of fleece growth', '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001214', 'bakewell-cattle',547,  'Bakewell selective breeding reduces cattle maturation from 3 years to ~1.5 years', '1871-01-01 00:00:00.000000');

-- Interruption contrast pair:
-- SpinningMill stop: discrete — yarn already sold; zero deterioration enforced.
-- LocomotiveWorks stop: connected — half-built locomotive deteriorates.
INSERT INTO working_period_interruptions (id, working_period_id, interrupted_at, resumed_at, deterioration_pence, waste_of_means_of_production_pence)
VALUES
  ('5eed000000000000001220', '5eed000000000000001201', '1871-06-01 08:00:00.000000', '1871-06-02 08:00:00.000000', 0,   0),
  ('5eed000000000000001221', '5eed000000000000001202', '1871-06-01 08:00:00.000000', '1871-06-08 08:00:00.000000', 240, 120);

-- Shortening rows: Bakewell breeding, machinery, parallel-labour.
INSERT INTO working_period_shortenings (id, working_period_id, kind, old_working_day_count, new_working_day_count, additional_fixed_capital_pence, additional_labour_count, occurred_at)
VALUES
  ('5eed000000000000001230', '5eed000000000000001202', 'bakewell_breeding', 1095, 547, 0,    0,  '1871-01-01 00:00:00.000000'),
  ('5eed000000000000001231', '5eed000000000000001202', 'machinery',         84,   42,  5000, 0,  '1871-06-01 00:00:00.000000'),
  ('5eed000000000000001232', '5eed000000000000001202', 'parallel_labour',   84,   28,  0,    50, '1871-09-01 00:00:00.000000');

-- London-builder credit-financing:
-- Own capital £100 (2,400 pence), borrowed capital £2,000 (48,000 pence).
INSERT INTO working_period_credit_financings (id, working_period_id, own_capital_pence, borrowed_capital_pence, mortgage_provider_id, created_at)
VALUES
  ('5eed000000000000001240', '5eed000000000000001202', 2400, 48000, 'vol3-credit-placeholder', '1871-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM working_period_credit_financings  WHERE id IN ('5eed000000000000001240');
DELETE FROM working_period_shortenings        WHERE id IN ('5eed000000000000001230','5eed000000000000001231','5eed000000000000001232');
DELETE FROM working_period_interruptions      WHERE id IN ('5eed000000000000001220','5eed000000000000001221');
DELETE FROM natural_working_period_constraints WHERE id IN ('5eed000000000000001210','5eed000000000000001211','5eed000000000000001212','5eed000000000000001213','5eed000000000000001214');
DELETE FROM working_periods                   WHERE id IN ('5eed000000000000001201','5eed000000000000001202');
