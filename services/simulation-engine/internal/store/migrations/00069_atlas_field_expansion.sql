-- +goose Up
-- Atlas Observatory field expansion — a fuller working town. In reality far more
-- than four circuits revolve at once, so this seeds the trades of a Lancashire /
-- Manchester mill-town of the 1860s — the very world from which Marx drew his
-- examples. Each is a sphere of differing organic composition, so the orrery
-- shows the spread out of which a general rate of profit forms (Vol. III Ch. 9):
-- labour-heavy trades pump a high s/(c+v); machine- and furnace-heavy trades a low
-- one. Several are literal Capital Vol. I fixtures — the full-priced baker and the
-- lucifer-match works of Ch. 10, the Staffordshire potteries of Ch. 10, the
-- brickfields of Ch. 15. Pence are pounds×100. Token 91 = Atlas expansion;
-- IDs 5eed...91XX. The turnover_number column exists from migration 00065.

-- Full-priced bakery (Vol. I Ch. 10 §3) — almost all living labour, fast daily
-- turnover, the highest rate of profit of the field. Total £800.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, turnover_number, created_at, updated_at)
VALUES
    ('5eed000000000000009101', 80000, 'money', 3, 'active', 12,
     '1865-01-07 04:00:00.000000', '1865-01-07 04:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009102', '5eed000000000000009101',
     '1865-01-21 04:00:00.000000', 20000, 40000, 20000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009103', '5eed000000000000009101', '1865', 60000, 84000, 24000);

-- Lucifer-match works (Vol. I Ch. 10 §3, the phosphorus trade) — the lowest
-- composition and the most savage rate of surplus-value in the town. Total £600.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, turnover_number, created_at, updated_at)
VALUES
    ('5eed000000000000009111', 60000, 'money', 3, 'active', 6,
     '1865-01-07 06:00:00.000000', '1865-01-07 06:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009112', '5eed000000000000009111',
     '1865-01-21 06:00:00.000000', 12000, 36000, 12000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009113', '5eed000000000000009111', '1865', 45000, 65000, 20000);

-- Staffordshire pottery / earthenware (Vol. I Ch. 10 §3, Dr. Greenhow's report) —
-- low-to-mid composition, hand-thrown, a high rate of profit. Total £1,800.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, turnover_number, created_at, updated_at)
VALUES
    ('5eed000000000000009121', 180000, 'money', 3, 'active', 4,
     '1865-01-07 08:00:00.000000', '1865-01-07 08:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009122', '5eed000000000000009121',
     '1865-01-21 08:00:00.000000', 36000, 108000, 36000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009123', '5eed000000000000009121', '1865', 135000, 174000, 39000);

-- Brickfield (Vol. I Ch. 15 §8, the brickmaking gangs) — low composition,
-- seasonal, child labour, a high rate of profit. Total £1,200.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, turnover_number, created_at, updated_at)
VALUES
    ('5eed000000000000009131', 120000, 'money', 3, 'active', 3,
     '1865-01-07 08:00:00.000000', '1865-01-07 08:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009132', '5eed000000000000009131',
     '1865-01-21 08:00:00.000000', 24000, 72000, 24000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009133', '5eed000000000000009131', '1865', 90000, 117000, 27000);

-- Calico print- & dye-works (Lancashire, the Print Works Act trade) — mid
-- composition, a middling rate of profit. Total £3,000.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, turnover_number, created_at, updated_at)
VALUES
    ('5eed000000000000009141', 300000, 'money', 3, 'active', 4,
     '1865-01-07 08:00:00.000000', '1865-01-07 08:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009142', '5eed000000000000009141',
     '1865-01-21 08:00:00.000000', 60000, 180000, 60000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009143', '5eed000000000000009141', '1865', 240000, 276000, 36000);

-- Brewery (Manchester ale) — mid-to-high composition, vats and plant, a modest
-- rate of profit. Total £4,500.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, turnover_number, created_at, updated_at)
VALUES
    ('5eed000000000000009151', 450000, 'money', 3, 'active', 3,
     '1865-01-07 08:00:00.000000', '1865-01-07 08:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009152', '5eed000000000000009151',
     '1865-01-21 08:00:00.000000', 90000, 270000, 90000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009153', '5eed000000000000009151', '1865', 360000, 396000, 36000);

-- Colliery (coal mine) — high composition, heavy fixed capital in shafts and
-- winding engines, slow turnover, a low rate of profit. Total £8,000.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, turnover_number, created_at, updated_at)
VALUES
    ('5eed000000000000009161', 800000, 'money', 3, 'active', 2,
     '1865-01-07 08:00:00.000000', '1865-01-07 08:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009162', '5eed000000000000009161',
     '1865-01-21 08:00:00.000000', 160000, 480000, 160000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009163', '5eed000000000000009161', '1865', 640000, 670000, 30000);

-- Gas-works — the highest composition in the town: retorts, mains, gasometers
-- against a small crew. The lowest rate of profit, turning over once a year.
-- Total £7,000.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, turnover_number, created_at, updated_at)
VALUES
    ('5eed000000000000009171', 700000, 'money', 3, 'active', 1,
     '1865-01-07 08:00:00.000000', '1865-01-07 08:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009172', '5eed000000000000009171',
     '1865-01-21 08:00:00.000000', 140000, 420000, 140000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009173', '5eed000000000000009171', '1865', 560000, 582000, 22000);

-- Locomotive & engineering works (Manchester, the Gorton engine sheds) — high
-- composition, great machine-tools, slow turnover, a low rate of profit.
-- Total £9,000.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, turnover_number, created_at, updated_at)
VALUES
    ('5eed000000000000009181', 900000, 'money', 3, 'active', 1,
     '1865-01-07 08:00:00.000000', '1865-01-07 08:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009182', '5eed000000000000009181',
     '1865-01-21 08:00:00.000000', 180000, 540000, 180000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009183', '5eed000000000000009181', '1865', 720000, 750000, 30000);

-- +goose Down
DELETE FROM supply_demand_imbalances WHERE id IN (
    '5eed000000000000009103','5eed000000000000009113','5eed000000000000009123',
    '5eed000000000000009133','5eed000000000000009143','5eed000000000000009153',
    '5eed000000000000009163','5eed000000000000009173','5eed000000000000009183');
DELETE FROM stage_distributions WHERE id IN (
    '5eed000000000000009102','5eed000000000000009112','5eed000000000000009122',
    '5eed000000000000009132','5eed000000000000009142','5eed000000000000009152',
    '5eed000000000000009162','5eed000000000000009172','5eed000000000000009182');
DELETE FROM industrial_capitals WHERE id IN (
    '5eed000000000000009101','5eed000000000000009111','5eed000000000000009121',
    '5eed000000000000009131','5eed000000000000009141','5eed000000000000009151',
    '5eed000000000000009161','5eed000000000000009171','5eed000000000000009181');
