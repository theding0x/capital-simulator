-- +goose Up
-- Vol. II Ch. 16 — The Turnover of Variable Capital
-- Seed fixtures (CC=16). Seed IDs use prefix 5eed000000000000160X.
--
-- Fixture A — Capitalist A (5-week turnover, n=10, £500 advance).
--   Marx §2: "Capital A of £500 is turned over 10 times a year.
--   Annual rate of surplus-value = 1000%."
--
-- Fixture B — Capitalist B (50-week turnover, n=1, £5,000 advance).
--   Marx §2: "Capital B of £5,000 is turned over once a year.
--   Annual rate of surplus-value = 100%."
--   Both pay out £5,000 in wages and produce £5,000 surplus per year.
--
-- Fixture C — Spinning Mill 1871 (n=52 weekly turnovers, s/v=156%).
--   Annual rate = 52 × 156% = 8,112%.
--
-- Fixture D — Locomotive Works (n≈4.3, 12-week connected cycle, s/v=80%).

-- Advances
INSERT INTO variable_capital_advances
    (id, industrial_capital_id, pence, turnover_period_minutes,
     turnovers_per_year_basis_points, created_at)
VALUES
    -- A: £500, 5-week period, n=10
    ('5eed0000000000001601', '5eed0000000000001600',
     120000, 50400, 100000, '1871-01-01 00:00:00.000000'),
    -- B: £5000, 50-week period, n=1
    ('5eed0000000000001602', '5eed0000000000001600',
     1200000, 504000, 10000, '1871-01-01 00:00:00.000000'),
    -- Spinning Mill 1871: £100/week advance, 1-week period, n=52
    ('5eed0000000000001603', '5eed0000000000001604',
     24000, 10080, 520000, '1871-01-01 00:00:00.000000'),
    -- Locomotive Works: £800 advance, 12-week period, n≈4.3
    ('5eed0000000000001605', '5eed0000000000001606',
     192000, 120960, 43461, '1871-01-01 00:00:00.000000');

-- Per-circuit surplus rates (s/v = 100% for A and B; 156% for SpinMill; 80% for LocoWorks)
INSERT INTO per_circuit_surplus_rates
    (id, variable_capital_advance_id, surplus_pence, variable_pence, basis_points)
VALUES
    ('5eed0000000000001607', '5eed0000000000001601', 120000, 120000, 10000),
    ('5eed0000000000001608', '5eed0000000000001602', 1200000, 1200000, 10000),
    ('5eed0000000000001609', '5eed0000000000001603', 37440, 24000, 15600),
    ('5eed0000000000001610', '5eed0000000000001605', 153600, 192000, 8000);

-- Annual surplus rates (period = "1871")
INSERT INTO annual_surplus_rates
    (id, variable_capital_advance_id, period,
     n_basis_points, per_circuit_basis_points, annual_basis_points)
VALUES
    -- A: n=10, s/v=100% → annual=1000%
    ('5eed0000000000001611', '5eed0000000000001601', '1871',
     100000, 10000, 100000),
    -- B: n=1, s/v=100% → annual=100%
    ('5eed0000000000001612', '5eed0000000000001602', '1871',
     10000, 10000, 10000),
    -- SpinMill: n=52, s/v=156% → annual=8112%
    ('5eed0000000000001613', '5eed0000000000001603', '1871',
     520000, 15600, 811200),
    -- LocoWorks: n≈4.3, s/v=80% → annual≈347%
    ('5eed0000000000001614', '5eed0000000000001605', '1871',
     43461, 8000, 34769);

-- Canonical A-vs-B contrast (ratio = 10×, or 100000 bp)
INSERT INTO annual_surplus_rate_contrasts
    (id, capital_a_id, capital_b_id,
     capital_a_annual_basis_points, capital_b_annual_basis_points, ratio_basis_points)
VALUES
    ('5eed0000000000001615',
     '5eed0000000000001601', '5eed0000000000001602',
     100000, 10000, 100000);

-- +goose Down
DELETE FROM annual_surplus_rate_contrasts
    WHERE id = '5eed0000000000001615';
DELETE FROM annual_surplus_rates
    WHERE id IN (
        '5eed0000000000001611', '5eed0000000000001612',
        '5eed0000000000001613', '5eed0000000000001614'
    );
DELETE FROM per_circuit_surplus_rates
    WHERE id IN (
        '5eed0000000000001607', '5eed0000000000001608',
        '5eed0000000000001609', '5eed0000000000001610'
    );
DELETE FROM variable_capital_advances
    WHERE id IN (
        '5eed0000000000001601', '5eed0000000000001602',
        '5eed0000000000001603', '5eed0000000000001605'
    );
