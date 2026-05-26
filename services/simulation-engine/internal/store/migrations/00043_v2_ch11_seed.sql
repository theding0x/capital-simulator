-- Vol. II Ch. 11 — Theories of Fixed and Circulating Capital: Ricardo. Seed data.
--
-- Seed IDs follow the pattern 5eed000000000000011NN where NN = two-digit ordinal.
-- Ricardo: Principles of Political Economy and Taxation (1817).

-- +goose Up

-- Ricardo — fixed and circulating capital (by durability), five KnownEconomistError entries.
INSERT INTO economist_attributions (id, concept, theorist, edition_year, anticipates)
VALUES ('5eed000000000000001101', 'fixed and circulating capital', 'Ricardo', 1817, 'prices_of_production');

INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001101', 'error_ricardo_durability_collapse');

INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001101', 'error_ricardo_conflation');

INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001101', 'error_ricardo_fixed_capital_price_explanation');

INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001101', 'error_ricardo_no_aggregate_turnover');

INSERT INTO economist_attribution_errors (attribution_id, error_code)
VALUES ('5eed000000000000001101', 'error_ricardo_no_value_revolution');

-- +goose Down
DELETE FROM economist_attribution_errors
WHERE attribution_id IN (
    '5eed000000000000001101'
);
DELETE FROM economist_attributions
WHERE id IN (
    '5eed000000000000001101'
);
