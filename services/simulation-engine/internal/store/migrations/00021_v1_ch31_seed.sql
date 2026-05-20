-- +goose Up
-- Ch. 31 seed: Genesis of the Industrial Capitalist.
-- Fixtures drawn from Marx's text: colonial plunder from the Americas,
-- Liverpool slave trade, Bank of England founding, and English protectionism.
-- Seed IDs follow 5eed00000000000000<CC><XX> where CC=31.
-- All records attach to the England 15th-18th c. historical stage (5eed000000000000002601).

INSERT INTO capital_origins (id, historical_stage_id, source, amount_pence, period, created_at)
VALUES
    ('5eed000000000000003101',
     '5eed000000000000002601',
     'colonial-plunder',
     200000,
     '1500-1800',
     '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003102',
     '5eed000000000000002601',
     'national-debt',
     120000,
     '1694-1800',
     '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003103',
     '5eed000000000000002601',
     'usury',
     40000,
     '15th-17th c.',
     '1867-01-01 00:00:00.000000');

INSERT INTO colonial_transfers (id, historical_stage_id, `from`, `to`, value_pence, method, created_at)
VALUES
    ('5eed000000000000003111',
     '5eed000000000000002601',
     'Americas',
     'England/Spain/Portugal',
     180000,
     'colonial-plunder',
     '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003112',
     '5eed000000000000002601',
     'West Africa',
     'England',
     15000,
     'slave-trade',
     '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003113',
     '5eed000000000000002601',
     'West Africa',
     'England',
     53000,
     'slave-trade',
     '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003114',
     '5eed000000000000002601',
     'West Africa',
     'England',
     74000,
     'slave-trade',
     '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003115',
     '5eed000000000000002601',
     'West Africa',
     'England',
     96000,
     'slave-trade',
     '1867-01-01 00:00:00.000000'),
    ('5eed000000000000003116',
     '5eed000000000000002601',
     'West Africa',
     'England',
     132000,
     'slave-trade',
     '1867-01-01 00:00:00.000000');

INSERT INTO national_debts (id, historical_stage_id, amount_pence, interest_rate_bps, creditor_class, created_at)
VALUES
    ('5eed000000000000003121',
     '5eed000000000000002601',
     2400000,
     800,
     'private-bankers',
     '1867-01-01 00:00:00.000000');

INSERT INTO protection_systems (id, historical_stage_id, tariff_rate_bps, beneficiary, period_start, period_end, created_at)
VALUES
    ('5eed000000000000003131',
     '5eed000000000000002601',
     3000,
     'English manufacturers',
     '17th c',
     '19th c',
     '1867-01-01 00:00:00.000000');

-- +goose Down
DELETE FROM protection_systems WHERE id IN ('5eed000000000000003131');
DELETE FROM national_debts WHERE id IN ('5eed000000000000003121');
DELETE FROM colonial_transfers WHERE id IN (
    '5eed000000000000003111',
    '5eed000000000000003112',
    '5eed000000000000003113',
    '5eed000000000000003114',
    '5eed000000000000003115',
    '5eed000000000000003116'
);
DELETE FROM capital_origins WHERE id IN (
    '5eed000000000000003101',
    '5eed000000000000003102',
    '5eed000000000000003103'
);
