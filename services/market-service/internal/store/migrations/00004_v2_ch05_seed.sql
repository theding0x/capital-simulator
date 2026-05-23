-- +goose Up

-- SpinningMill1871 TurnoverTime: 1-week production + 1-week circulation.
-- §: "total time = time of production + time of circulation."
INSERT INTO turnover_times
    (id, industrial_capital_id, labour_time_nanos, labour_interruption_nanos, latent_nanos,
     natural_process_nanos, selling_time_nanos, buying_time_nanos, circulation_open, created_at, updated_at)
VALUES
    ('5eed000000000000000501',
     '5eed000000000000000502',
     604800000000000, 0, 0, 0,
     432000000000000, 172800000000000,
     0,
     '1871-01-01 00:00:00.000000', '1871-01-01 00:00:00.000000');

-- Natural process spans: grain (90d ripening), wine (14d fermentation), leather (30d tanning).
-- §: "grain sown, wine fermenting in the cellar, tanneries exposed to chemical processes."
INSERT INTO natural_process_spans
    (id, turnover_time_id, industrial_capital_id, process, duration_nanos, started_at)
VALUES
    ('5eed000000000000000503',
     '5eed000000000000000501',
     '5eed000000000000000502',
     'ripening', 7776000000000000, '1871-03-01 00:00:00.000000'),
    ('5eed000000000000000504',
     '5eed000000000000000501',
     '5eed000000000000000502',
     'fermentation', 1209600000000000, '1871-03-01 00:00:00.000000'),
    ('5eed000000000000000505',
     '5eed000000000000000501',
     '5eed000000000000000502',
     'tanning', 2592000000000000, '1871-03-01 00:00:00.000000');

-- Perishability: milk (2 days), beer (7 days), yarn (0 = no constraint).
-- §: "perishable commodity … the absolute limit of its time of circulation."
INSERT INTO perishability (id, commodity_id, window_nanos)
VALUES
    ('5eed000000000000000506', '5eed000000000000000510', 172800000000000),
    ('5eed000000000000000507', '5eed000000000000000511', 604800000000000),
    ('5eed000000000000000508', '5eed000000000000000512', 0);

-- LatentProductiveCapital: £50 cotton stockpile at SpinningMill1871.
-- §: "£50 of cotton stockpiled but not yet spun … fallow capital."
INSERT INTO latent_productive_capital
    (id, turnover_time_id, industrial_capital_id, pence, held_at, entered_production_at)
VALUES
    ('5eed000000000000000509',
     '5eed000000000000000501',
     '5eed000000000000000502',
     12000, '1871-01-01 00:00:00.000000', NULL);

-- +goose Down

DELETE FROM latent_productive_capital WHERE id IN (
    '5eed000000000000000509'
);
DELETE FROM perishability WHERE id IN (
    '5eed000000000000000506',
    '5eed000000000000000507',
    '5eed000000000000000508'
);
DELETE FROM natural_process_spans WHERE id IN (
    '5eed000000000000000503',
    '5eed000000000000000504',
    '5eed000000000000000505'
);
DELETE FROM turnover_times WHERE id IN (
    '5eed000000000000000501'
);
