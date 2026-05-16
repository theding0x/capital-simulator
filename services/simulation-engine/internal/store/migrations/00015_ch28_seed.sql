-- +goose Up
-- Ch. 28 — Bloody Legislation Against the Expropriated.
-- All records are attached to the seeded "England 15th–18th c." historical
-- stage (5eed000000000000002601) from the Ch. 26 seed migration.
-- Magnitudes are from Marx's text; halfpence used as unit for Statute of
-- Labourers (2s = 24 halfpennies), George II tailors (2s 7¾d ≈ 31 halfpennies).

INSERT INTO wage_statutes
    (id, historical_stage_id, period, jurisdiction, max_wage_pence, min_wage_pence, enforcement_penalty, created_at)
VALUES
    ('5eed000000000000002801',
     '5eed000000000000002601',
     '1349',
     'England',
     24,
     0,
     'imprisonment',
     NOW(6)),
    ('5eed000000000000002802',
     '5eed000000000000002601',
     '1740s',
     'London tailors',
     31,
     0,
     'fine and imprisonment',
     NOW(6));

INSERT INTO vagrancy_laws
    (id, historical_stage_id, period, jurisdiction, punishment, target_population, created_at)
VALUES
    ('5eed000000000000002803',
     '5eed000000000000002601',
     '1530',
     'England',
     'whipping+imprisonment',
     'sturdy vagabonds',
     NOW(6)),
    ('5eed000000000000002804',
     '5eed000000000000002601',
     '1597',
     'England',
     'whipping, branding, execution for repeat offenders',
     'rogues, vagabonds, and sturdy beggars',
     NOW(6));

-- +goose Down
DELETE FROM vagrancy_laws WHERE id IN (
    '5eed000000000000002803',
    '5eed000000000000002804'
);
DELETE FROM wage_statutes WHERE id IN (
    '5eed000000000000002801',
    '5eed000000000000002802'
);