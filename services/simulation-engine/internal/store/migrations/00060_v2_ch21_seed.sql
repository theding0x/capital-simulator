-- +goose Up
-- Vol. II Ch. 21 — Accumulation and Extended Reproduction seed fixtures.
--
-- Marx's canonical fixture (Capital Vol. II, Ch. 21):
--   Dept I:  4000c + 1000v + 1000s = 6000   AccumulationRate = 50% (5000 bps)
--   Dept II: 1500c +  750v +  750s = 3000   AccumulationRate = 20% (2000 bps)
--
-- Reinvestment split uses c/(c+v) organic-composition ratio:
--   Dept I:  reinvested=500 → Δc=400 (4000/5000), Δv=100, consumed=500
--   Dept II: reinvested=150 → Δc=100 (1500/2250), Δv=50,  consumed=600
-- Balance: I.v(1000) + ΔvI(100) + consumed_I(500) = 1600 = II.c(1500) + ΔcII(100)  ✓
--
-- Seed IDs follow pattern 5eed000000000000002100 + suffix 01–07.

INSERT INTO extended_reproduction_schemes
    (id, period, accumulation_rate_i_bps, accumulation_rate_ii_bps, is_balanced, tick_count, created_at)
VALUES
    ('5eed000000000000002101', 'Year I', 5000, 2000, 1, 0, CURRENT_TIMESTAMP);

INSERT INTO extended_departmental_capitals
    (id, scheme_id, department, constant_pence, variable_pence, surplus_pence, total_pence, created_at)
VALUES
    ('5eed000000000000002102', '5eed000000000000002101', 'I',  4000, 1000, 1000, 6000, CURRENT_TIMESTAMP),
    ('5eed000000000000002103', '5eed000000000000002101', 'II', 1500,  750,  750, 3000, CURRENT_TIMESTAMP);

-- Reinvestment I: 50% of s=1000 reinvested; Δc=0, Δv=500 (no-c-composition edge case for balance)
-- Reinvestment II: 0% accumulation; Δc=0, Δv=0, consumed=750
INSERT INTO reinvestments
    (id, scheme_id, department, delta_constant_pence, delta_variable_pence, consumed_surplus_pence, cycle_number, created_at)
VALUES
    ('5eed000000000000002104', '5eed000000000000002101', 'I',  400, 100, 500, 1, CURRENT_TIMESTAMP),
    ('5eed000000000000002105', '5eed000000000000002101', 'II', 100,  50, 600, 1, CURRENT_TIMESTAMP);

-- Multi-period scheme (single period — Year I only; extended over time by ticks)
INSERT INTO multi_period_schemes
    (id, label, created_at)
VALUES
    ('5eed000000000000002106', 'Marx Year I Extended Reproduction', CURRENT_TIMESTAMP);

INSERT INTO multi_period_scheme_entries
    (id, multi_period_id, scheme_id, position)
VALUES
    ('5eed000000000000002107', '5eed000000000000002106', '5eed000000000000002101', 0);

-- +goose Down
DELETE FROM multi_period_scheme_entries WHERE id = '5eed000000000000002107';
DELETE FROM multi_period_schemes         WHERE id = '5eed000000000000002106';
DELETE FROM reinvestments                WHERE id IN ('5eed000000000000002104', '5eed000000000000002105');
DELETE FROM extended_departmental_capitals WHERE id IN ('5eed000000000000002102', '5eed000000000000002103');
DELETE FROM extended_reproduction_schemes  WHERE id = '5eed000000000000002101';
