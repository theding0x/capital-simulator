-- +goose Up
-- Vol. III Ch. 12 seed — the two admissible causes of a change in the price of
-- production, so the Ch. 12 panel comes up populated on a fresh MySQL volume.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=12, XX=row.
-- These 22-char IDs never collide with the 24-char hex output of New*ID().
--
--   1201 cotton: the general rate of profit rose (own value unchanged) → price up.
--   1202 linen:  the commodity's own value fell (general rate unchanged) → price down.

INSERT INTO price_of_production_changes
    (id, sphere_name, cause, rate_changed, value_changed, price_changed, old_price, new_price, created_at)
VALUES
    ('5eed000000000000001201', 'cotton', 'general_rate_change',     1, 0, 1, 122, 130, '2026-05-01 00:00:00.000000'),
    ('5eed000000000000001202', 'linen',  'individual_value_change', 0, 1, 1, 122, 110, '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM price_of_production_changes WHERE id IN (
    '5eed000000000000001201',
    '5eed000000000000001202'
);
