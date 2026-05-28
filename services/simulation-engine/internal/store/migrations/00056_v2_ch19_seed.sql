-- Vol. II Ch. 19 — Former Presentations of the Subject: seed data.
--
-- Seed IDs follow the pattern 5eed000000000000001901, 1902, …
-- Two new EconomistAttribution rows (Quesnay 1758 reproduction context,
-- Ricardo 1817 reproduction entry) and three new error codes on the
-- existing Smith row (5eed000000000000001003, seeded in 00042).

-- +goose Up

-- Quesnay — Tableau Économique (1758): reproduction context.
-- Anticipates the social character of the reproduction circuit;
-- carries two Ch. 19 theoretical errors.
INSERT INTO economist_attributions (id, concept, theorist, edition_year, anticipates)
VALUES ('5eed000000000000001901', 'Tableau Économique', 'Quesnay', 1758, 'reproduction_as_social');

INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001901', 'error_quesnay_productive_sterile_division');

INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001901', 'error_quesnay_missing_value_theory');

-- Ricardo — Principles of Political Economy (1817): reproduction entry.
-- Sees Smith's revenue-dogma error at the individual-capital level but
-- never develops it into a two-department balance for total social capital.
INSERT INTO economist_attributions (id, concept, theorist, edition_year, anticipates)
VALUES ('5eed000000000000001902', 'reproduction and net revenue', 'Ricardo', 1817, 'reproduction_incomplete');

INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001902', 'error_ricardo_incomplete_reproduction');

-- Smith — extensions to the existing row (5eed000000000000001003, from 00042).
-- Three additional Ch. 19 errors: the revenue dogma and its two consequences.
INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001003', 'error_smith_revenue_dogma');

INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001003', 'error_smith_missing_constant_replacement');

INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001003', 'error_smith_labour_regress');

-- +goose Down
DELETE FROM economist_attribution_errors
WHERE attribution_id IN (
    '5eed000000000000001901',
    '5eed000000000000001902'
);
DELETE FROM economist_attributions
WHERE id IN (
    '5eed000000000000001901',
    '5eed000000000000001902'
);
DELETE FROM economist_attribution_errors
WHERE attribution_id = '5eed000000000000001003'
  AND error_code IN (
    'error_smith_revenue_dogma',
    'error_smith_missing_constant_replacement',
    'error_smith_labour_regress'
);
