-- +goose Up
-- Vol. II Ch. 7 seed — Marx's turnover fixtures.
-- Seed IDs: 5eed000000000000 + 07 (chapter) + sequence.

-- Spinning Mill: weekly turnover, n ≈ 52/year.
-- t = 10080 min (7 days). BasisPoints = 10000*525600/10080 = 521428.
INSERT INTO turnovers (id, industrial_capital_id, lens, turnover_time_minutes)
VALUES ('5eed000000000000070101', '5eed000000000000000401', 'productive', 10080);

INSERT INTO turnover_cycles (id, turnover_id, started_at, ended_at, advance_pence, returned_pence, production_minutes, circulation_minutes)
VALUES ('5eed000000000000070201', '5eed000000000000070101', '1871-01-01 00:00:00.000', '1871-01-08 00:00:00.000', 48000, 57600, 8000, 2080),
       ('5eed000000000000070301', '5eed000000000000070101', '1871-01-08 00:00:00.000', '1871-01-15 00:00:00.000', 48000, 57600, 8000, 2080);

INSERT INTO turnover_numbers (id, turnover_id, year_reference_minutes, turnover_time_minutes, basis_points)
VALUES ('5eed000000000000070401', '5eed000000000000070101', 525600, 10080, 521428);

-- Locomotive-maker: 18-month turnover, n = 2/3.
-- t = 18*30*24*60 = 777600 min. BasisPoints = 10000*525600/777600 = 6759.
INSERT INTO turnovers (id, industrial_capital_id, lens, turnover_time_minutes)
VALUES ('5eed000000000000070501', '', 'money', 777600);

INSERT INTO turnover_cycles (id, turnover_id, started_at, ended_at, advance_pence, returned_pence, production_minutes, circulation_minutes)
VALUES ('5eed000000000000070601', '5eed000000000000070501', '1870-01-01 00:00:00.000', '1871-07-01 00:00:00.000', 500000, 600000, 720000, 57600),
       ('5eed000000000000070701', '5eed000000000000070501', '1871-07-01 00:00:00.000', '1873-01-01 00:00:00.000', 500000, 600000, 720000, 57600);

INSERT INTO turnover_numbers (id, turnover_id, year_reference_minutes, turnover_time_minutes, basis_points)
VALUES ('5eed000000000000070801', '5eed000000000000070501', 525600, 777600, 6759);

-- Wine-maturer: 4-year turnover, n = 1/4.
-- t = 4*525600 = 2102400 min. BasisPoints = 10000*525600/2102400 = 2500.
INSERT INTO turnovers (id, industrial_capital_id, lens, turnover_time_minutes)
VALUES ('5eed000000000000070901', '', 'productive', 2102400);

INSERT INTO turnover_cycles (id, turnover_id, started_at, ended_at, advance_pence, returned_pence, production_minutes, circulation_minutes)
VALUES ('5eed000000000000070a01', '5eed000000000000070901', '1867-01-01 00:00:00.000', '1871-01-01 00:00:00.000', 200000, 260000, 2060000, 42400),
       ('5eed000000000000070b01', '5eed000000000000070901', '1871-01-01 00:00:00.000', '1875-01-01 00:00:00.000', 200000, 260000, 2060000, 42400);

INSERT INTO turnover_numbers (id, turnover_id, year_reference_minutes, turnover_time_minutes, basis_points)
VALUES ('5eed000000000000070c01', '5eed000000000000070901', 525600, 2102400, 2500);

-- +goose Down
DELETE FROM turnover_numbers WHERE id IN ('5eed000000000000070401','5eed000000000000070801','5eed000000000000070c01');
DELETE FROM turnover_cycles WHERE id IN ('5eed000000000000070201','5eed000000000000070301','5eed000000000000070601','5eed000000000000070701','5eed000000000000070a01','5eed000000000000070b01');
DELETE FROM turnovers WHERE id IN ('5eed000000000000070101','5eed000000000000070501','5eed000000000000070901');
