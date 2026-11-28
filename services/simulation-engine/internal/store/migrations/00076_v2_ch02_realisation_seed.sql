-- +goose Up
-- Spinning Mill 1871 (Vol. II Ch. 2 §I): C' is 10,000 lbs of yarn valued at
-- £500, of which £422 (8,440 lbs) replaces the productive capital and £78
-- (1,560 lbs) is surplus-value. Legacy pence: 42,200 + 7,800 = 50,000.
-- Both exemplars from 00028 sold C' in full ("at its value and in its
-- entirety"). A partial sale is rejected by design and cannot be seeded.
INSERT INTO productive_circuit_commodities
    (productive_circuit_id, value_original, value_surplus)
VALUES
    ('5eed000000000000020a', 42200, 7800),
    ('5eed000000000000020b', 42200, 7800);

INSERT INTO productive_circuit_realisations
    (productive_circuit_id, realised_pence, sold_at)
VALUES
    ('5eed000000000000020a', 50000, '1871-01-06 00:00:00.000000'),
    ('5eed000000000000020b', 50000, '1871-01-12 00:00:00.000000');

-- +goose Down
DELETE FROM productive_circuit_realisations WHERE productive_circuit_id IN ('5eed000000000000020a', '5eed000000000000020b');
DELETE FROM productive_circuit_commodities WHERE productive_circuit_id IN ('5eed000000000000020a', '5eed000000000000020b');
