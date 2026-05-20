-- +goose Up
-- §15th-century metayer: "He advances one part of the agricultural stock,
-- the landlord the other. The two divide the total product in proportions
-- determined by contract." Annual lease, no money rent, half-stock advanced.
-- §long-lease effect: "contracts for farms ran for a long time, often for 99
-- years. The progressive fall in the value of the precious metals ... lowered
-- wages. A portion of the latter was now added to profits."
INSERT INTO farm_tenures
    (id, historical_stage_id, form, lease_period_years, rent_pence, capital_advanced_pence, revenue_pence, wage_costs_pence, created_at)
VALUES
    ('5eed000000000000002901',
     '5eed000000000000002601',
     'metayer',
     1, 0, 500, 1000, 200,
     NOW(6)),
    ('5eed000000000000002902',
     '5eed000000000000002601',
     'capitalist-farmer',
     99, 1200, 0, 3000, 600,
     NOW(6));

-- +goose Down
DELETE FROM farm_tenures WHERE id IN (
    '5eed000000000000002901',
    '5eed000000000000002902'
);
