-- +goose Up
-- Vol. III Ch. 13 seed — Marx's §law trajectory (s'=100%, v=100, c rises) plus
-- two rate-mass contradiction fixtures, so the Ch. 13 panel comes up populated
-- on a fresh MySQL volume.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=13, XX=row.
-- These 22-char IDs never collide with the 24-char hex output of New*ID().
--
--   1301 trajectory: s'=100%, v=100, c rises 50→400; p' falls 6667→2000 bp.
--   1302 rate-mass equal:  C doubles, rate halves → mass unchanged (200,000).
--   1303 rate-mass rising: C grows at constant rate → mass rises 200,000→300,000.
--
-- profit_rate values are round-half-up basis points: 6667 (≈66.7%), 5000 (50%),
-- 3333 (≈33.3%), 2500 (25%), 2000 (20%). The periods_json field names match the
-- Go json tags on tendency.TrajectoryPeriod.

INSERT INTO composition_trajectories
    (id, label, surplus_value_rate, periods_json, created_at)
VALUES
    (
        '5eed000000000000001301',
        'The Law As Such — rising composition',
        10000,
        '[{"constant_capital":50,"variable_capital":100,"surplus_value":100,"profit_rate":6667},{"constant_capital":100,"variable_capital":100,"surplus_value":100,"profit_rate":5000},{"constant_capital":200,"variable_capital":100,"surplus_value":100,"profit_rate":3333},{"constant_capital":300,"variable_capital":100,"surplus_value":100,"profit_rate":2500},{"constant_capital":400,"variable_capital":100,"surplus_value":100,"profit_rate":2000}]',
        '2026-05-01 00:00:00.000000'
    );

INSERT INTO rate_mass_contradictions
    (id, old_c, old_rate, new_c, new_rate, old_mass, new_mass, mass_change, created_at)
VALUES
    ('5eed000000000000001302', 1000000, 2000, 2000000, 1000, 200000, 200000, 0,      '2026-05-01 00:00:00.000000'),
    ('5eed000000000000001303', 1000000, 2000, 3000000, 1000, 200000, 300000, 100000, '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM rate_mass_contradictions WHERE id IN (
    '5eed000000000000001302',
    '5eed000000000000001303'
);
DELETE FROM composition_trajectories WHERE id IN (
    '5eed000000000000001301'
);
