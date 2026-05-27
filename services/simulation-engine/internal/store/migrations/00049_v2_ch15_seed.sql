-- +goose Up
-- Vol. II Ch. 15 — The Effects of a Change of Prices
-- Seed fixtures (CC=15).
--
-- Fixture A — Spinning Mill 1871: raw-cotton price falls 15% (§1 Price Fall).
--   Marx: "If the value (price) of raw and auxiliary materials falls, then in the
--   renewed process of production less money-capital is required to buy the same
--   amount of them; a portion of the money-capital previously employed is set free."
--   Old advance on cotton: £480 (48000p). New advance: 48000 × 8500/10000 = 40800p.
--   Capital freed: 7200p.  Operator elects "same scale" (capital hoarded).
--
-- Fixture B — Manchester Spinner 1857 Crisis: output price collapses 10% (§3).
--   Marx: "If the value of the commodities falls after they have been produced but
--   before they have been sold, the capitalist loses part of his surplus-value."
--   Expected realisation £600 (60000p); actual £540 (54000p). Loss 6000p.
--
-- Fixture C — Cotton stock revaluation (§2 paper loss on inventory).
--   5000 lbs × (10p − 12p) = −10000p unrealised loss.
--
-- Fixture D — Speculative hold (§ "concentration of wealth, not creation of value").

-- Fixture A: event + revalued capital + case.
INSERT INTO value_revolution_events
    (id, industrial_capital_id, affecting, direction_pence, occurred_at,
     basis_points_change, revalued_capital_pence)
VALUES
    ('5eed0000000000001502', '5eed0000000000001501',
     'means_of_production', -7200, '1871-04-01 00:00:00.000000',
     -1500, 40800);

INSERT INTO revalued_capitals
    (id, value_revolution_event_id, industrial_capital_id,
     old_advance_pence, new_advance_pence, delta_pence)
VALUES
    ('5eed0000000000001503', '5eed0000000000001502', '5eed0000000000001501',
     48000, 40800, -7200);

INSERT INTO price_case_records (value_revolution_event_id, fall_case, rise_case)
VALUES ('5eed0000000000001502', 'same_scale', NULL);

-- Fixture B: event + realised value delta.
INSERT INTO value_revolution_events
    (id, industrial_capital_id, affecting, direction_pence, occurred_at,
     basis_points_change, revalued_capital_pence)
VALUES
    ('5eed0000000000001505', '5eed0000000000001504',
     'output', -6000, '1857-10-01 00:00:00.000000',
     0, 0);

INSERT INTO realised_value_deltas
    (id, value_revolution_event_id, industrial_capital_id,
     expected_realisation_pence, actual_realisation_pence, delta_pence)
VALUES
    ('5eed0000000000001506', '5eed0000000000001505', '5eed0000000000001504',
     60000, 54000, -6000);

-- Fixture C: inventory revaluation (paper loss on held cotton stock).
INSERT INTO inventory_revaluations
    (id, industrial_capital_id, commodity_id,
     quantity_held, old_unit_pence, new_unit_pence, delta_pence)
VALUES
    ('5eed0000000000001507', '5eed0000000000001501', 'raw-cotton-1871',
     5000, 12, 10, -10000);

-- Fixture D: speculative hold.
INSERT INTO speculative_holds
    (id, industrial_capital_id, commodity_id, held_since, anticipated_direction)
VALUES
    ('5eed0000000000001508', '5eed0000000000001509',
     'manchester-cotton-1857', '1857-09-01 00:00:00.000000', 1);

-- +goose Down
DELETE FROM speculative_holds WHERE id IN ('5eed0000000000001508');
DELETE FROM inventory_revaluations WHERE id IN ('5eed0000000000001507');
DELETE FROM realised_value_deltas WHERE id IN ('5eed0000000000001506');
DELETE FROM price_case_records WHERE value_revolution_event_id IN
    ('5eed0000000000001502', '5eed0000000000001505');
DELETE FROM revalued_capitals WHERE id IN ('5eed0000000000001503');
DELETE FROM value_revolution_events WHERE id IN
    ('5eed0000000000001502', '5eed0000000000001505');
