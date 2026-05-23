-- +goose Up
-- Spinning Mill 1871 — simple reproduction (§I).
-- Opening C′ = £500: c=£372 (37200p), v=£50 (5000p), s=£78 (7800p). 10,000 lbs.
-- Terminal C′ = opening C′ (simple reproduction; all surplus exits as revenue).
-- Seed IDs follow 5eed00000000000003<XX> (chapter 03).
INSERT INTO commodity_circuits
    (id, agent_id, constant_pence, variable_pence, surplus_pence, pounds_total,
     mode, is_first_investment, social_capital_lens,
     terminal_constant, terminal_variable, terminal_surplus, terminal_capitalised,
     terminal_pounds, terminal_closed_at, created_at)
VALUES
    ('5eed000000000000030a', '', 37200, 5000, 7800, 10000,
     'simple', 0, 0,
     37200, 5000, 7800, 0, 10000, '1871-01-07 00:00:00.000000',
     '1871-01-01 00:00:00.000000');

-- Three successive partial sales for SpinningMill1871Simple (Marx's aliquot table).
-- Decomposition uses floor division: c = opening_c * qty / total; surplus absorbs remainder.
-- Sale 1: 7,440 lbs at £372 (37200p) — c=27676 (floor 27676.8), v=3720, s=5804.
-- Sale 2: 1,000 lbs at £50  (5000p)  — c=3720,  v=500,  s=780.
-- Sale 3: 1,560 lbs at £78  (7800p)  — c=5803 (floor 5803.2), v=780, s=1217.
-- Each row: c+v+s == realised_pence. Grand total realised = 50000p (£500). ✓
INSERT INTO successive_partial_sales
    (commodity_circuit_id, quantity, realised_pence, constant_share, variable_share, surplus_share, sold_at)
VALUES
    ('5eed000000000000030a', 7440, 37200, 27676, 3720, 5804, '1871-01-02 00:00:00.000000'),
    ('5eed000000000000030a', 1000,  5000,  3720,  500,  780, '1871-01-03 00:00:00.000000'),
    ('5eed000000000000030a', 1560,  7800,  5803,  780, 1217, '1871-01-04 00:00:00.000000');

-- MP sources for SpinningMill1871Simple: coal from a producer circuit; raw cotton as import.
INSERT INTO mp_sources
    (commodity_circuit_id, source_commodity_circuit_id, source_kind)
VALUES
    ('5eed000000000000030a', '', 'producer_circuit'),
    ('5eed000000000000030a', '', 'import');

-- Spinning Mill 1871 — extended reproduction (§III).
-- Opening C′ = £500. All surplus capitalised into constant capital.
-- Terminal C′ = £578: c=45000 (£450), v=5000 (£50), s=7800 (£78), capitalised=7800.
INSERT INTO commodity_circuits
    (id, agent_id, constant_pence, variable_pence, surplus_pence, pounds_total,
     mode, is_first_investment, social_capital_lens,
     terminal_constant, terminal_variable, terminal_surplus, terminal_capitalised,
     terminal_pounds, terminal_closed_at, created_at)
VALUES
    ('5eed000000000000030b', '', 37200, 5000, 7800, 10000,
     'extended', 0, 0,
     45000, 5000, 7800, 7800, 11560, '1871-01-14 00:00:00.000000',
     '1871-01-08 00:00:00.000000');

-- +goose Down
DELETE FROM mp_sources WHERE commodity_circuit_id IN ('5eed000000000000030a');
DELETE FROM successive_partial_sales WHERE commodity_circuit_id IN ('5eed000000000000030a');
DELETE FROM commodity_circuits WHERE id IN ('5eed000000000000030a', '5eed000000000000030b');
