-- +goose Up
-- Vol. II Ch. 9 seed — Marx's canonical aggregate-turnover fixtures.
--
-- Fixture 1: §3 Cotton Mill 1867 — £80,000 fixed (10yr) + £20,000 circulating (5×/yr).
--   AnnualTurnedOverPence = £800k + £10,000k = £10,800k → aggregate n = 10800 bp (1 + 2/25).
--
-- Fixture 2: §5 Scrope Three-Component Capital — $50,000 over three components,
--   $33,750 annual → mean term ≈ 18 months without surplus-value.
--
-- Seed IDs: 5eed000000000000000900xx (ch09 = 09).

-- Fixture 1: Cotton Mill 1867 aggregate turnover
INSERT INTO aggregate_turnovers (id, industrial_capital_id, lens, advanced_total_pence, annual_turned_over_pence, aggregate_number_bp)
VALUES ('5eed00000000000000090001', '5eed00000000000000090010', 'money', 10000000, 10800000, 10800);

-- Fixture 1 contributions
INSERT INTO component_turnover_contributions (id, aggregate_turnover_id, capital_component_id, kind, advanced_pence, turnover_number_bp, annual_return_pence, difference_kind)
VALUES
  ('5eed00000000000000090002', '5eed00000000000000090001', '', 'fixed',       8000000, 1000,  800000,   'real'),
  ('5eed00000000000000090003', '5eed00000000000000090001', '', 'circulating', 2000000, 50000, 10000000, 'real');

-- Fixture 1 lifetime cycle — 10-year fixed-capital reproduction cycle (10 × 525,600 = 5,256,000 min)
INSERT INTO lifetime_cycles (id, aggregate_turnover_id, duration_minutes, started_at, completed_cycles)
VALUES ('5eed00000000000000090004', '5eed00000000000000090001', 5256000, '1867-01-01 00:00:00.000000', 0);

-- Fixture 1 crisis phases — Marx §4 references 1847, 1857, 1866 as crisis years
INSERT INTO crisis_cycle_phase_records (id, lifetime_cycle_id, phase, observed_at, notes)
VALUES
  ('5eed00000000000000090005', '5eed00000000000000090004', 'crisis',          '1847-10-01 00:00:00.000000', '1847 commercial crisis — railway mania and credit collapse'),
  ('5eed00000000000000090006', '5eed00000000000000090004', 'medium_activity', '1850-01-01 00:00:00.000000', 'post-1847 recovery, medium activity 1850–1856'),
  ('5eed00000000000000090007', '5eed00000000000000090004', 'precipitancy',    '1856-01-01 00:00:00.000000', 'precipitancy phase preceding 1857 crisis'),
  ('5eed00000000000000090008', '5eed00000000000000090004', 'crisis',          '1857-11-01 00:00:00.000000', '1857 world market crisis — beginning of a new cycle');

-- Fixture 2: Scrope $50,000 three-component capital (§5)
INSERT INTO aggregate_turnovers (id, industrial_capital_id, lens, advanced_total_pence, annual_turned_over_pence, aggregate_number_bp)
VALUES ('5eed00000000000000090011', '5eed00000000000000090020', 'money', 5000000, 3375000, 6750);

-- Fixture 2 contributions (three-component split matching §5 narrative)
INSERT INTO component_turnover_contributions (id, aggregate_turnover_id, capital_component_id, kind, advanced_pence, turnover_number_bp, annual_return_pence, difference_kind)
VALUES
  ('5eed00000000000000090012', '5eed00000000000000090011', '', 'fixed',       2000000, 1000,  200000,  'real'),
  ('5eed00000000000000090013', '5eed00000000000000090011', '', 'circulating', 2000000, 10000, 2000000, 'real'),
  ('5eed00000000000000090014', '5eed00000000000000090011', '', 'circulating', 1000000, 11750, 1175000, 'real');

-- +goose Down
DELETE FROM crisis_cycle_phase_records WHERE id IN (
  '5eed00000000000000090005','5eed00000000000000090006',
  '5eed00000000000000090007','5eed00000000000000090008'
);
DELETE FROM lifetime_cycles WHERE id = '5eed00000000000000090004';
DELETE FROM component_turnover_contributions WHERE id IN (
  '5eed00000000000000090002','5eed00000000000000090003',
  '5eed00000000000000090012','5eed00000000000000090013','5eed00000000000000090014'
);
DELETE FROM aggregate_turnovers WHERE id IN (
  '5eed00000000000000090001','5eed00000000000000090011'
);
