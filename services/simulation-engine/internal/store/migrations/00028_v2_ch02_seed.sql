-- +goose Up
-- Spinning Mill 1871 — simple reproduction (§I).
-- £422 productive capital; £78 surplus exits as revenue. Circuit opens again at £422.
-- Seed IDs follow the 5eed00000000000000<CC> prefix (CC=02 for Ch. 2).
INSERT INTO productive_circuits
    (id, agent_id, constant_pence, variable_pence, mode, min_capitalisation_increment_pence,
     latent_accumulated, reserve_balance, created_at)
VALUES
    ('5eed000000000000020a', '', 37200, 5000, 'simple', 100, 0, 0, '1871-01-01 00:00:00.000000');

INSERT INTO revenue_circuits (productive_circuit_id, amount, spent_at)
VALUES ('5eed000000000000020a', 7800, '1871-01-07 00:00:00.000000');

-- Spinning Mill 1871 — extended reproduction (§II).
-- £422 productive capital capitalised up to £500 after a CapitalisationStep of £78.
INSERT INTO productive_circuits
    (id, agent_id, constant_pence, variable_pence, mode, min_capitalisation_increment_pence,
     latent_accumulated, reserve_balance, created_at)
VALUES
    ('5eed000000000000020b', '', 37200, 5000, 'extended', 100, 0, 0, '1871-01-08 00:00:00.000000');

-- The £78 capitalisation step: all injected into constant capital for this exemplar.
-- Marx leaves the Δc/Δv split unspecified at this stage; we use full Δc for illustration.
INSERT INTO capitalisation_steps
    (productive_circuit_id, amount_injected, delta_constant_pence, delta_variable_pence, occurred_at)
VALUES
    ('5eed000000000000020b', 7800, 7800, 0, '1871-01-14 00:00:00.000000');

-- +goose Down
DELETE FROM capitalisation_steps WHERE productive_circuit_id IN ('5eed000000000000020b');
DELETE FROM revenue_circuits WHERE productive_circuit_id IN ('5eed000000000000020a');
DELETE FROM productive_circuits WHERE id IN ('5eed000000000000020a', '5eed000000000000020b');
