-- +goose Up
-- Vol. II Ch. 4 seed — Spinning Mill 1871 IndustrialCapital
-- Aggregates the Ch.1/2/3 seed circuits for the Spinning Mill.
-- §intro: "10,000 lbs. of yarn … simultaneously in all three circuits."
-- §Meeting of Demand and Supply: £5,000 capital, 80c+20v+20s, £400/yr sinking fund.
-- Seed IDs: 5eed00000000000000 04 XX

-- IndustrialCapital: Spinning Mill 1871
INSERT INTO industrial_capitals
    (id, agent_id, money_circuit_id, productive_circuit_id, commodity_circuit_id,
     total_pence, economy_mode, stagnation_tolerance_ticks, status, created_at, updated_at)
VALUES
    ('5eed000000000000000401',
     '5eed000000000000000201',
     '5eed000000000000003101',
     '5eed000000000000002701',
     '5eed000000000000002901',
     500000, 'money', 3, 'active',
     '1871-01-07 08:00:00.000000', '1871-01-07 08:00:00.000000');

-- CapitalParts: three cohorts simultaneously present (Nebeneinander)
INSERT INTO capital_parts (id, industrial_capital_id, pence, stage, entered_stage_at)
VALUES ('5eed000000000000000402', '5eed000000000000000401',
        100000, 'money', '1871-01-07 08:00:00.000000');
INSERT INTO capital_parts (id, industrial_capital_id, pence, stage, entered_stage_at)
VALUES ('5eed000000000000000403', '5eed000000000000000401',
        300000, 'production', '1871-01-14 08:00:00.000000');
INSERT INTO capital_parts (id, industrial_capital_id, pence, stage, entered_stage_at)
VALUES ('5eed000000000000000404', '5eed000000000000000401',
        100000, 'commodity', '1871-01-21 08:00:00.000000');

-- StageDistribution: three-week history (total == 500000 each week)
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000000405', '5eed000000000000000401',
     '1871-01-07 08:00:00.000000', 100000, 300000, 100000),
    ('5eed000000000000000406', '5eed000000000000000401',
     '1871-01-14 08:00:00.000000', 100000, 300000, 100000),
    ('5eed000000000000000407', '5eed000000000000000401',
     '1871-01-21 08:00:00.000000', 100000, 300000, 100000);

-- SupplyDemandImbalance: 80c+20v+20s at £5000 scale, 5x turnover
-- demand=400000, supply=480000, excess=80000 (5:6 ratio preserved)
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000000408', '5eed000000000000000401',
     '1871', 400000, 480000, 80000);

-- SinkingFund: £4000 fixed capital, 10-year life → £400/yr
INSERT INTO sinking_funds
    (id, industrial_capital_id, fixed_capital_pence, lifetime_years,
     annual_payment_pence, accumulated_pence)
VALUES
    ('5eed000000000000000409', '5eed000000000000000401',
     400000, 10, 40000, 0);

-- +goose Down
DELETE FROM sinking_funds            WHERE id IN ('5eed000000000000000409');
DELETE FROM supply_demand_imbalances WHERE id IN ('5eed000000000000000408');
DELETE FROM stage_distributions      WHERE id IN ('5eed000000000000000405','5eed000000000000000406','5eed000000000000000407');
DELETE FROM capital_parts            WHERE id IN ('5eed000000000000000402','5eed000000000000000403','5eed000000000000000404');
DELETE FROM industrial_capitals      WHERE id IN ('5eed000000000000000401');
