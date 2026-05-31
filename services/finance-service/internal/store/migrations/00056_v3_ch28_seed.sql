-- +goose Up
-- Vol. III Ch. 28 seeds — Marx's currency-circulation examples.
-- Seed IDs start with 5eed000000000000002800..2801 (chapter 28).
--
-- Seed 2800: Bank of England 1844 — £20m notes outstanding, £8m reserve.
--            Coin function dominant: ~£14m serves retail revenue expenditure.
--            Capital-transfer: ~£4m serves productive-capital movement.
--            Marx cites these figures to show Tooke/Fullarton's conflation.
--
-- Seed 2801: Scottish Banking System — credit velocity reduces coin need.
--            £100 Scotland ≈ £420 England in circulation effect (Marx).
--            Capital-transfer function dominates; minimal metallic reserve.

INSERT INTO `currency_observations`
    (`id`, `name`, `total_currency_outstanding`, `coin_function_amount`, `capital_transfer_amount`,
     `reserve_amount`, `reserve_description`, `description`, `created_at`)
VALUES
    ('5eed000000000000002800',
     'Bank of England 1844',
     2000000000,
     1400000000,
     400000000,
     800000000,
     'Bank of England Banking Department reserve 1844',
     'Bank Charter Act 1844: £20m notes outstanding; £8m reserve; coin function dominates (retail/wage payments); capital transfers handled via bills and book entries',
     '2025-01-01 00:00:00.000000'),
    ('5eed000000000000002801',
     'Scottish Banking System — Credit Velocity',
     1000000000,
     100000000,
     700000000,
     5000000,
     'Scottish bank metallic reserve, minimal due to credit velocity',
     'Scottish credit system: high velocity reduces need for coin; £100 Scottish ≈ £420 English in circulation effect; capital-transfer function dominates',
     '2025-01-01 00:01:00.000000');

-- +goose Down
DELETE FROM `currency_observations`
WHERE `id` IN ('5eed000000000000002800', '5eed000000000000002801');
