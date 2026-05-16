-- +goose Up
-- §1 composition unchanged: capital doubles while organic composition stays fixed.
-- Workers absorbed proportionally; reserve army does not swell.
INSERT INTO general_law_scenarios
    (id, name, constant_capital, variable_capital, surplus_rate,
     accumulation_rate, productivity_growth, wage_pence, worker_supply, periods, created_at)
VALUES
    ('5eed000000000000002501',
     '§1 Unchanged composition — demand rises with capital',
     8000, 2000, 1.0, 1.0, 0.0, 200, 50, 5, NOW(6));

-- §2 rising organic composition: machinery displaces labour despite capital growth.
-- Reserve army swells even as total capital grows.
INSERT INTO general_law_scenarios
    (id, name, constant_capital, variable_capital, surplus_rate,
     accumulation_rate, productivity_growth, wage_pence, worker_supply, periods, created_at)
VALUES
    ('5eed000000000000002502',
     '§2 Rising OC — reserve army swells as machinery displaces labour',
     8000, 2000, 1.0, 1.0, 0.05, 200, 50, 10, NOW(6));

-- +goose Down
DELETE FROM general_law_scenarios WHERE id IN (
    '5eed000000000000002501',
    '5eed000000000000002502'
);
