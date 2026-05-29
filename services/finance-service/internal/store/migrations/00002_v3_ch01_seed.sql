-- +goose Up
-- Vol. III Ch. 1 seed — Marx's worked example from Capital III, Ch. 1, §1, so
-- the cost-price/profit panel comes up populated on a fresh MySQL volume.
--
-- Seed ID convention: 5eed00000000000000<CC><XX> where CC=01 (chapter) and XX
-- is the row sequence. These 22-char IDs never collide with the 24-char hex
-- profit.NewCostPriceID() output.
--
-- Row 0101 is the plain example: 400c + 100v + 100s = £600, cost-price k = £500.
-- Row 0102 is the same outlay shown with the fixed/circulating split: a £1,200
-- instrument of labour contributes only its £20 wear-and-tear to the cost-price,
-- alongside £380 circulating constant and £100 variable — fixed component £20,
-- circulating component £480, k still £500.

INSERT INTO cost_prices
    (id, constant, variable, fixed_wear_and_tear, fixed_advanced, k, fixed_component, circulating_component, created_at)
VALUES
    ('5eed000000000000000101', 400, 100, 0, 0, 500, 0, 500,
     '2026-05-01 00:00:00.000000'),
    ('5eed000000000000000102', 400, 100, 20, 1200, 500, 20, 480,
     '2026-05-01 00:00:00.000000');

-- +goose Down
DELETE FROM cost_prices WHERE id IN (
    '5eed000000000000000101',
    '5eed000000000000000102'
);
