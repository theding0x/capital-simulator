-- +goose Up
-- Vol. III Ch. 25 seeds — Marx's credit and fictitious capital examples.
-- Seed IDs start with 5eed000000000000002500..2503 (chapter 25).
--
-- Bill seeds:
--   Seed 2500: 1839 UK bill of exchange — face value £132,123,460 (×100 pence),
--              90-day maturity; represents commercial credit circulating as money.
--   Seed 2501: Clearing House daily settlement — £3,000,000 face settled
--              with only £200,000 actual money (1840s); 1-day maturity (symbolic).
--
-- Fictitious capital seeds:
--   Seed 2502: UK Government Consol (KindPublicDebt=3) — AnnualIncome=500 p (£5)
--              at 5% (500 bp) → CapitalisedValue=10000 p (£100).
--              real_capital_exists=0: the original £100 is consumed.
--   Seed 2503: Railway Share (KindShare=4) — AnnualIncome=1000 p at 5% (500 bp)
--              → CapitalisedValue=20000 p; real_capital_exists=1.

INSERT INTO `bills_of_exchange`
    (`id`, `drawer`, `drawee`, `face_value`, `maturity_days`, `description`, `created_at`)
VALUES
    ('5eed000000000000002500', 'Manchester Spinner', 'London Merchant',   13212346000, 90, '1839 UK bills outstanding — £132m face, circulating as commercial money', '2025-01-01 00:00:00.000000'),
    ('5eed000000000000002501', 'Bank of England',    'Clearing House',          300000,  1, '1840s clearing-house settlement: £3m settled with £200k actual money (symbolic)', '2025-01-01 00:01:00.000000');

INSERT INTO `fictitious_capitals`
    (`id`, `kind`, `name`, `annual_income`, `rate_bp`, `capitalised_value`, `real_capital_exists`, `description`, `created_at`)
VALUES
    ('5eed000000000000002502', 3, 'UK Government Consol', 500,  500, 10000, 0, 'Government bond: £5 annual income @ 5% → £100 capitalised; original capital consumed', '2025-01-01 00:02:00.000000'),
    ('5eed000000000000002503', 4, 'Railway Share',        1000, 500, 20000, 1, 'Railway Co. stock: £10 dividend @ 5% → £200 capitalised; real capital in track/locos', '2025-01-01 00:03:00.000000');

-- +goose Down
DELETE FROM `fictitious_capitals`
WHERE `id` IN ('5eed000000000000002502', '5eed000000000000002503');

DELETE FROM `bills_of_exchange`
WHERE `id` IN ('5eed000000000000002500', '5eed000000000000002501');
