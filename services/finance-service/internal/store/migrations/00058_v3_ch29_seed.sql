-- +goose Up
-- Vol. III Ch. 29 seeds — Marx's bank-capital exemplars.
-- Seed IDs: 5eed000000000000002900..2901 for bank_capitals,
--           5eed000000000000002910..2911 for fictitious_capital_valuations.
--
-- Bank capital 2900: a Scottish bank — cash £3m (300,000,000 p), securities £24m
--   (2,400,000,000 p), total £27m. Cash is a small fraction of the whole; the
--   greater part is fictitious (bills, bonds, deposits as book entries).
-- Bank capital 2901: Bank of England — larger reserve in coin, deep holdings of
--   public securities.
--
-- Valuation 2910: government bond — £5 annual income (500 p) capitalised at the
--   prevailing 5% (500 bp) → market value £100 (10,000 p).
-- Valuation 2911: joint-stock share — nominal £100 (10,000 p), 10% dividend
--   (1000 bp) → £10 income (1000 p), capitalised at 5% (500 bp) → £200 (20,000 p).

INSERT INTO `bank_capitals`
    (`id`, `name`, `cash_amount`, `securities_amount`, `total_capital`, `components_json`, `description`, `created_at`)
VALUES
    ('5eed000000000000002900', 'Scottish Bank', 300000000, 2400000000, 2700000000,
     '[{"amount":300000000,"kind":1,"description":"gold and notes in the till"},{"amount":1400000000,"kind":2,"description":"discounted commercial paper"},{"amount":1000000000,"kind":3,"description":"government consols"}]',
     'Cash £3m against £24m of securities; the greater part of bank capital is fictitious',
     '2025-01-01 00:00:00.000000'),
    ('5eed000000000000002901', 'Bank of England', 1200000000, 4800000000, 6000000000,
     '[{"amount":1200000000,"kind":1,"description":"coin and bullion reserve"},{"amount":3000000000,"kind":3,"description":"public securities"},{"amount":1800000000,"kind":2,"description":"discounted bills"}]',
     'National bank: deep holdings of public securities over a coin reserve',
     '2025-01-01 00:01:00.000000');

INSERT INTO `fictitious_capital_valuations`
    (`id`, `name`, `nominal_value`, `annual_income`, `interest_rate_bp`, `dividend_rate_bp`, `market_value`, `description`, `created_at`)
VALUES
    ('5eed000000000000002910', 'Government Bond @5%', 10000, 500, 500, 0, 10000,
     'A £5 annual income capitalised at 5% is worth £100; a fall in the rate would inflate it',
     '2025-01-01 00:00:00.000000'),
    ('5eed000000000000002911', 'Joint-Stock Share', 10000, 1000, 500, 1000, 20000,
     'Nominal £100 paying a 10% dividend (£10) is worth £200 when capitalised at the 5% market rate',
     '2025-01-01 00:01:00.000000');

-- +goose Down
DELETE FROM `fictitious_capital_valuations`
WHERE `id` IN ('5eed000000000000002910', '5eed000000000000002911');

DELETE FROM `bank_capitals`
WHERE `id` IN ('5eed000000000000002900', '5eed000000000000002901');
