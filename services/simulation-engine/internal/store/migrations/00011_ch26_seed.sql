-- +goose Up
-- England 15th–18th c. — the classic form of primitive accumulation.
-- Three episodes — agrarian expropriation by enclosure, expropriation
-- of church property, and colonial plunder — together build the initial
-- capital and free-proletarian supply that the Ch. 23/24 reproduction
-- loop takes as its starting point. Magnitudes are pedagogical.
INSERT INTO historical_stages (id, name, description, created_at) VALUES
    ('5eed000000000000002601',
     'England 15th–18th c.',
     'Agrarian expropriation creating a free proletariat — the classic form of primitive accumulation, in which enclosure, statutory wage-fixing, and colonial plunder together found capitalist production.',
     NOW(6));

INSERT INTO primitive_accumulations
    (id, stage_id, ordinal, period, method, labourers_expropriated, capital_formed)
VALUES
    ('5eed000000000000002602',
     '5eed000000000000002601', 0,
     '1450–1640', 'enclosure of common lands',
     40000, 80000),
    ('5eed000000000000002603',
     '5eed000000000000002601', 1,
     '1500–1640', 'expropriation of church property',
     15000, 30000),
    ('5eed000000000000002604',
     '5eed000000000000002601', 2,
     '1500–1800', 'colonial plunder and slave trade',
     5000, 90000);

-- +goose Down
DELETE FROM primitive_accumulations WHERE id IN (
    '5eed000000000000002602',
    '5eed000000000000002603',
    '5eed000000000000002604'
);
DELETE FROM historical_stages WHERE id IN (
    '5eed000000000000002601'
);
