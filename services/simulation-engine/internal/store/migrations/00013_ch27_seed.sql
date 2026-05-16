-- +goose Up
-- England — the three paradigm waves of agricultural expropriation that
-- Marx analyses in Ch. 27. Magnitudes are pedagogical; Sutherland figures
-- (794,000 acres, 15,000 displaced) are from Marx's own footnote.
INSERT INTO enclosure_events
    (id, period, acres_enclosed, population_displaced, beneficiary, created_at)
VALUES
    ('5eed000000000000002701',
     '1450–1540',
     500000,
     40000,
     'landed gentry',
     NOW(6)),
    ('5eed000000000000002702',
     '1760–1820',
     3000000,
     150000,
     'improving landlords (Parliamentary enclosures)',
     NOW(6)),
    ('5eed000000000000002703',
     '1814–1820',
     794000,
     15000,
     'Countess of Sutherland',
     NOW(6));

-- +goose Down
DELETE FROM enclosure_events WHERE id IN (
    '5eed000000000000002701',
    '5eed000000000000002702',
    '5eed000000000000002703'
);