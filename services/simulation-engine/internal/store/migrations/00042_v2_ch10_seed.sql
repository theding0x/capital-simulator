-- Vol. II Ch. 10 — Theories of Fixed and Circulating Capital: seed data.
--
-- Seed IDs follow the pattern 5eed000000000000010NN where NN = two-digit ordinal.
-- Quesnay: Tableau Économique (1758).
-- Smith:   An Inquiry into the Nature and Causes of the Wealth of Nations (1776).

-- +goose Up

-- Quesnay — avances primitives (anticipates fixed capital)
INSERT INTO economist_attributions (id, concept, theorist, edition_year, anticipates)
VALUES ('5eed000000000000001001', 'avances primitives', 'Quesnay', 1758, 'fixed_capital_item');

-- Quesnay — avances annuelles (anticipates circulating capital)
INSERT INTO economist_attributions (id, concept, theorist, edition_year, anticipates)
VALUES ('5eed000000000000001002', 'avances annuelles', 'Quesnay', 1758, 'circulating');

-- Smith — fixed and circulating stock (carries all three KnownEconomistError entries)
INSERT INTO economist_attributions (id, concept, theorist, edition_year, anticipates)
VALUES ('5eed000000000000001003', 'fixed and circulating stock', 'Smith', 1776, 'fixed_capital_item');

INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001003', 'error_smith_conflation');

INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001003', 'error_smith_circulation_capital_conflation');

INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001003', 'error_smith_revenue_in_capital');

-- +goose Down
DELETE FROM economist_attribution_errors
WHERE attribution_id IN (
    '5eed000000000000001003'
);
DELETE FROM economist_attributions
WHERE id IN (
    '5eed000000000000001001',
    '5eed000000000000001002',
    '5eed000000000000001003'
);
