-- +goose Up
-- Vol. III Ch. 31 seeds — Marx's two sources of loanable money-capital.
-- Seed IDs start with 5eed000000000000003100..3101 (chapter 31).
--
-- Seed 3100: Ipswich farmers' surplus deposits (source=2) — rural revenue and
--            capital converted into money that bill-brokers rediscount to the
--            manufacturing centres. This ACCOMPANIES real accumulation. The
--            velocity row shows the same £20 lent 5× performing £100 of loan
--            capital (2000 × 5 = 10000 pence).
--
-- Seed 3101: Post-crisis idle industrial capital (source=1) — money lying fallow
--            after a crisis, transformed into loan capital. This is INVERSELY
--            related to real accumulation: it swells precisely when production
--            stagnates. velocity_times_lent=1 (lent once); banking economies in
--            the circulation reserve are recorded as a saving.

INSERT INTO `floating_capitals`
    (`id`, `name`, `amount`, `source`, `velocity_money_amount`, `velocity_times_lent`, `velocity_loan_capital_created`, `reserve_saving_amount`, `reserve_saving_description`, `description`, `created_at`)
VALUES
    ('5eed000000000000003100',
     'Ipswich farmers'' surplus deposits',
     200000000,
     2,
     2000,
     5,
     10000,
     5000000,
     'country-bank till economies pooled by the credit system',
     'rural surplus rediscounted by Lombard Street bill-brokers to the manufacturing centres; accompanies real accumulation',
     '2025-01-01 00:00:00.000000'),
    ('5eed000000000000003101',
     'Post-crisis idle industrial capital',
     300000000,
     1,
     100000,
     1,
     100000,
     2000000,
     'economy in the reserve of means of circulation concentrated by banking',
     'money lying fallow after a crisis, simply transformed into loan capital; inverse to real accumulation',
     '2025-01-01 00:01:00.000000');

-- +goose Down
DELETE FROM `floating_capitals`
WHERE `id` IN ('5eed000000000000003100', '5eed000000000000003101');
