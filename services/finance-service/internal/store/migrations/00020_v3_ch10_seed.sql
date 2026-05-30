-- +goose Up
-- Vol. III Ch. 10 seed — Marx's cotton sphere and high-productivity firm
-- so the Ch. 10 panels come up populated on a fresh MySQL volume.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=10 and XX is row.
-- These 22-char IDs never collide with the 24-char hex output of the New*ID()
-- functions.
--
-- …1001  market_value  cotton: bulk=100, best=80, worst=130, value=100
-- …1002  surplus_profit: High-productivity firm; ind=60, market=100, qty=1000,
--         general_rate=2200, amount=(100-60)*1000=40000
-- …1003  equalisation  cotton: initial_rate=3000, target_rate=2200,
--         direction='inflow', is_converging=1, market_price=0, market_value=0

INSERT INTO market_values
    (id, sphere_name, bulk_condition_value, best_condition_value, worst_condition_value, value, created_at)
VALUES
    ('5eed000000000000001001', 'cotton', 100, 80, 130, 100, '2026-05-01 00:00:00.000000');

INSERT INTO surplus_profits
    (id, firm_name, individual_value, market_value, output_qty, general_rate, amount, created_at)
VALUES
    ('5eed000000000000001002', 'High-productivity firm', 60, 100, 1000, 2200, 40000, '2026-05-01 00:00:00.000000');

INSERT INTO equalisations
    (id, sphere_name, initial_rate, target_rate, direction, is_converging, market_price, market_value, created_at)
VALUES
    ('5eed000000000000001003', 'cotton', 3000, 2200, 'inflow', 1, 0, 0, '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM equalisations WHERE id IN (
    '5eed000000000000001003'
);
DELETE FROM surplus_profits WHERE id IN (
    '5eed000000000000001002'
);
DELETE FROM market_values WHERE id IN (
    '5eed000000000000001001'
);
