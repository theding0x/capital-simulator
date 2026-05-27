-- +goose Up

-- England-India merchant capital TurnoverTime (12-month circuit in the 1850s).
-- §: "If capital … requires 12 months to complete its circuit, then it can
--    produce a mass of surplus-value only equal to what it produces in one year."
INSERT INTO turnover_times
    (id, industrial_capital_id, labour_time_nanos, labour_interruption_nanos,
     latent_nanos, natural_process_nanos, selling_time_nanos, buying_time_nanos,
     circulation_open, created_at, updated_at)
VALUES
    ('5eed000000000000001401',
     '5eed000000000000001414',
     0, 0, 0, 0,
     525600000000000, 0,
     0,
     '1850-01-01 00:00:00.000000', '1850-01-01 00:00:00.000000');

-- England-India selling phases (pre- and post-telegraph).
-- Pre-telegraph (1850): ship post London–Bombay took ~6 months each way.
-- §: "the communication of news … required from 12 to 18 months."
INSERT INTO selling_phases
    (id, turnover_time_id, industrial_capital_id, commodity_circuit_id,
     opened_at, closed_at, outcome)
VALUES
    ('5eed000000000000001411',
     '5eed000000000000001401', '5eed000000000000001414', '',
     '1850-01-01 00:00:00.000000', '1851-01-01 00:00:00.000000', 'sold'),
    ('5eed000000000000001412',
     '5eed000000000000001401', '5eed000000000000001414', '',
     '1870-06-01 00:00:00.000000', '1870-07-01 00:00:00.000000', 'sold');

-- SellingPhaseDetailed rows for the two England-India selling phases.
-- Pre-telegraph (1850): 365-day phase; communication lag = 250 days (ship post).
-- Post-telegraph (1870): 30-day phase; communication lag = 1 day (telegraph).
-- §: "whereas the communication of a house in Calcutta required from 12 to
--    18 months, now it requires only a few hours."
INSERT INTO selling_phase_details
    (id, selling_phase_id, turnover_time_id, industrial_capital_id,
     market_id, market_distance_meters, communication_lag_minutes,
     buyer_search_minutes, bargaining_minutes, legal_settlement_minutes,
     total_minutes, created_at)
VALUES
    ('5eed000000000000001413',
     '5eed000000000000001411', '5eed000000000000001401', '5eed000000000000001414',
     'london-bombay', 9000000, 360000,
     72000, 72000, 21600,
     525600, '1850-01-01 00:00:00.000000'),
    ('5eed000000000000001415',
     '5eed000000000000001412', '5eed000000000000001401', '5eed000000000000001414',
     'london-bombay', 9000000, 1440,
     14400, 14400, 12960,
     43200, '1870-06-01 00:00:00.000000');

-- DistanceLagRelations for three routes.
-- London-Manchester (~320 km): horse post 3 min/km before railway.
-- London-Bombay (~9 000 km): 58 min/km by sail (6-month mail); 1 min/km post-telegraph.
INSERT INTO distance_lag_relations
    (id, market_id, distance_meters, lag_minutes_per_km, created_at)
VALUES
    ('5eed000000000000001403',
     'london-manchester', 320000, 3, '1840-01-01 00:00:00.000000'),
    ('5eed000000000000001404',
     'london-bombay', 9000000, 58, '1845-01-01 00:00:00.000000'),
    ('5eed000000000000001405',
     'london-bombay', 9000000, 1, '1870-01-01 00:00:00.000000');

-- CirculationSpeedImprovements — four key technological events.
-- §: "railways … steamships … the electric telegraph — means that have
--    revolutionised the time of circulation."
INSERT INTO circulation_speed_improvements
    (id, market_id, old_lag_minutes_per_km, new_lag_minutes_per_km, reason,
     effective_at, created_at)
VALUES
    ('5eed000000000000001406',
     'london-manchester', 10, 3, 'railway',
     '1840-01-01 00:00:00.000000', '1840-01-01 00:00:00.000000'),
    ('5eed000000000000001407',
     'london-bombay', 90, 58, 'steamship',
     '1845-01-01 00:00:00.000000', '1845-01-01 00:00:00.000000'),
    ('5eed000000000000001408',
     'london-bombay', 58, 1, 'telegraph',
     '1870-01-01 00:00:00.000000', '1870-01-01 00:00:00.000000'),
    ('5eed000000000000001409',
     'london-manchester', 3, 0, 'telephone',
     '1876-01-01 00:00:00.000000', '1876-01-01 00:00:00.000000');

-- AnnualSurplusPenalty: England-India merchant capital, period 1850.
-- Ideal = 10 000 bp (100% if no circulation time).
-- Actual = 5 000 bp (50%: 6-month selling phase halves productive fraction).
-- §: "the longer the time of circulation … the less surplus-value is produced."
INSERT INTO annual_surplus_penalties
    (id, industrial_capital_id, period,
     ideal_annual_rate_basis_points, actual_annual_rate_basis_points, penalty_basis_points)
VALUES
    ('5eed000000000000001410',
     '5eed000000000000001414', '1850',
     10000, 5000, 5000);

-- +goose Down

DELETE FROM annual_surplus_penalties WHERE id IN (
    '5eed000000000000001410'
);
DELETE FROM circulation_speed_improvements WHERE id IN (
    '5eed000000000000001406',
    '5eed000000000000001407',
    '5eed000000000000001408',
    '5eed000000000000001409'
);
DELETE FROM distance_lag_relations WHERE id IN (
    '5eed000000000000001403',
    '5eed000000000000001404',
    '5eed000000000000001405'
);
DELETE FROM selling_phase_details WHERE id IN (
    '5eed000000000000001413',
    '5eed000000000000001415'
);
DELETE FROM selling_phases WHERE id IN (
    '5eed000000000000001411',
    '5eed000000000000001412'
);
DELETE FROM turnover_times WHERE id IN (
    '5eed000000000000001401'
);
