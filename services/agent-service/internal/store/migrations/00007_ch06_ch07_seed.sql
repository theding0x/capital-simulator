-- +goose Up
-- Seed data for Capital Vol. I, Ch. 6–7.
-- Ch.6: labour-power as commodity — a cotton spinner sells their capacity
-- to a mill-owning capitalist.
-- Ch.7: the capitalist sets the spinner to work; surplus-value emerges.
--
-- Shilling conversion (from §2.a): 1 shilling = 120 labour-minutes.
-- Daily value of labour-power = 3 shillings = 360 min (necessary labour).
-- Working day = 12 hours = 720 min; 6 h necessary + 6 h surplus.

INSERT INTO labour_workers
    (id, kind, owns_labour_power, owns_commodities_to_sell,
     capacity_minutes_per_day, labour_power_value_minutes, created_at, updated_at)
VALUES
    -- Mary: free labourer — owns only her capacity for labour (§6.1).
    -- Labour-power value: 360 min/day (3 shillings) to reproduce subsistence.
    ('5eed00000000000000060001', 'worker', 1, 0, 720, 360,
     '2026-05-01 08:00:00.000000', '2026-05-01 08:00:00.000000'),
    -- Tom: simple commodity producer — owns both goods and labour-power.
    ('5eed00000000000000060002', 'worker', 1, 1, 480, 240,
     '2026-05-01 08:00:00.000000', '2026-05-01 08:00:00.000000');

INSERT INTO labour_capitalists
    (id, kind, money_capital, created_at, updated_at)
VALUES
    -- Mill owner: advances 1800 min money-capital (15 shillings):
    -- 1440 min constant capital (cotton + spindle) + 360 min variable capital.
    ('5eed00000000000000060010', 'capitalist', 1800,
     '2026-05-01 08:00:00.000000', '2026-05-01 08:00:00.000000');

-- Mary posts her labour-power for sale: 720 min/day capacity, 5-day contract,
-- asking wage 360 min/day (= daily reproduction cost).
INSERT INTO labour_power_offerings
    (id, owner_id, capacity_minutes_per_day, contract_days, asking_wage, created_at)
VALUES
    ('5eed00000000000000060030', '5eed00000000000000060001', 720, 5, 360,
     '2026-05-01 08:30:00.000000');

-- Mill owner purchases Mary's labour-power for 1 day at 360 min/day wage.
INSERT INTO labour_power_purchases
    (id, seller_id, buyer_id, wage_minutes, contract_days, created_at)
VALUES
    ('5eed00000000000000060020', '5eed00000000000000060001', '5eed00000000000000060010',
     360, 1, '2026-05-01 09:00:00.000000');

-- Ch.7 §2: mill owner sets Mary to work for 12 hours (720 min).
-- Raw material: 10 lbs cotton at 120 min SNLT/lb → 1200 min transferred.
-- Instrument:   spindle, 240 min wear per run → 240 min transferred.
-- TransferredValue = 1440 min.  ValueAdded = 720 min.  ProductValue = 2160 min.
-- NecessaryLabour = 360 min.    SurplusLabour = 360 min.  SurplusValue = 360 min.
-- Rate of surplus-value (s/v) = 360/360 = 100%.
INSERT INTO labour_processes
    (id, worker_id, capitalist_id, duration, necessary_labour_minutes,
     means_json, product_kind, product_quantity, created_at)
VALUES
    ('5eed00000000000000070001',
     '5eed00000000000000060001',
     '5eed00000000000000060010',
     720, 360,
     '{"raw_materials":[{"commodity_id":"cotton","quantity":10,"snlt_per_unit":120}],"instruments":[{"commodity_id":"spindle","wear_per_run":240}]}',
     'yarn', 10,
     '2026-05-01 10:00:00.000000');

-- +goose Down
DELETE FROM labour_processes WHERE id = '5eed00000000000000070001';
DELETE FROM labour_power_purchases WHERE id = '5eed00000000000000060020';
DELETE FROM labour_power_offerings WHERE id = '5eed00000000000000060030';
DELETE FROM labour_capitalists WHERE id = '5eed00000000000000060010';
DELETE FROM labour_workers WHERE id IN (
    '5eed00000000000000060001',
    '5eed00000000000000060002'
);
