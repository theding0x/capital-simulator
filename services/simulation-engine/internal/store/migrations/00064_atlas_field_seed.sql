-- +goose Up
-- Atlas Observatory field seed — extra industrial capitals of differing
-- organic composition (Vol. III Ch. 9 "spheres") so the average rate of
-- profit visibly forms from spread. Token 90 = Atlas field. IDs: 5eed...90XX.

-- Tannery (low composition, high surplus): total £2,000.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, created_at, updated_at)
VALUES
    ('5eed000000000000009001', 200000, 'money', 3, 'active',
     '1871-01-07 08:00:00.000000', '1871-01-07 08:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009002', '5eed000000000000009001',
     '1871-01-21 08:00:00.000000', 40000, 120000, 40000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009003', '5eed000000000000009001', '1871', 150000, 195000, 45000);

-- Steelworks (high composition, low surplus): total £6,000.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, created_at, updated_at)
VALUES
    ('5eed000000000000009011', 600000, 'money', 3, 'active',
     '1871-01-07 08:00:00.000000', '1871-01-07 08:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009012', '5eed000000000000009011',
     '1871-01-21 08:00:00.000000', 120000, 360000, 120000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009013', '5eed000000000000009011', '1871', 480000, 504000, 24000);

-- Textile (mid composition): total £3,500.
INSERT INTO industrial_capitals
    (id, total_pence, economy_mode, stagnation_tolerance_ticks, status, created_at, updated_at)
VALUES
    ('5eed000000000000009021', 350000, 'money', 3, 'active',
     '1871-01-07 08:00:00.000000', '1871-01-07 08:00:00.000000');
INSERT INTO stage_distributions
    (id, industrial_capital_id, at_time, money_pence, production_pence, commodity_pence)
VALUES
    ('5eed000000000000009022', '5eed000000000000009021',
     '1871-01-21 08:00:00.000000', 70000, 210000, 70000);
INSERT INTO supply_demand_imbalances
    (id, industrial_capital_id, period, demand_pence, supply_pence, excess_pence)
VALUES
    ('5eed000000000000009023', '5eed000000000000009021', '1871', 280000, 322000, 42000);

-- +goose Down
DELETE FROM supply_demand_imbalances WHERE id IN ('5eed000000000000009003','5eed000000000000009013','5eed000000000000009023');
DELETE FROM stage_distributions      WHERE id IN ('5eed000000000000009002','5eed000000000000009012','5eed000000000000009022');
DELETE FROM industrial_capitals      WHERE id IN ('5eed000000000000009001','5eed000000000000009011','5eed000000000000009021');
