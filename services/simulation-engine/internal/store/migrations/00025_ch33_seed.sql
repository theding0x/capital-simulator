-- +goose Up
-- Ch. 33 seed: four Marx-faithful colonial labour markets, named after
-- the cases Marx and Wakefield actually invoke in §1-§4 of the chapter.
--
--   Swan River, 1829 — Mr. Peel's £50,000 fiasco; 300 working-class
--     persons brought from England who walked off into independent
--     settlement because no sufficient price held them on the wage
--     market. The canonical anti-capitalist fixture.
--
--   Saint Domingo, c.1700 — the Spanish settlers Marx cites at the
--     end of §1 ("without labourers, their capital must have
--     perished"); free land, small immigrant cohort, no scheme.
--
--   South Australia, 1836 — Wakefield's testbed for systematic
--     colonisation; sufficient price in force, labourers held on the
--     wage market while the Land Fund imports fresh have-nothings.
--
--   Upper Canada, 1830s — the "supernumerary" labour pool Wakefield
--     describes in §4 (the conversation with "capitalists of Canada
--     and the state of New York"). Sufficient price applied with a
--     shorter dependence period.
--
-- Seed IDs follow 5eed00000000000000<CC><XX> where CC=33.

INSERT INTO colonial_labour_markets
    (id, colony, free_labourers, annual_wage_pence, land_available,
     wakefield_scheme_applied, independence_years, surplus_labour_extractable, created_at)
VALUES
    ('5eed000000000000003301',
     'Swan River, 1829',
     300, 4800, 1,
     0, 0, 0,
     '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003302',
     'Saint Domingo, 1700',
     50, 3600, 1,
     0, 0, 0,
     '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003303',
     'South Australia, 1836',
     5000, 4800, 1,
     1, 50, 1,
     '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003304',
     'Upper Canada, 1830s',
     12000, 5200, 1,
     1, 8, 1,
     '1867-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM colonial_labour_markets WHERE id IN (
    '5eed000000000000003301',
    '5eed000000000000003302',
    '5eed000000000000003303',
    '5eed000000000000003304'
);
